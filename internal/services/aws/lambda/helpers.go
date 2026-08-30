// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// maxFunctionRefLength is the FunctionName @length(1,140) maximum from the
// Smithy model. The bound applies to the raw input regardless of form: a
// full or partial ARN reference must fit in 140 characters, while the bare
// name inside any form is separately capped at 64 by validateFunctionName.
const maxFunctionRefLength = 140

// resolveFunctionRef parses every FunctionName form the API accepts into
// the bare function name and a qualifier embedded in the reference:
//   - "my-function"                                  (name only)
//   - "my-function:v1"                               (name with alias or version)
//   - "arn:aws:lambda:us-west-2:123456789012:function:my-function[:v1]" (full ARN)
//   - "123456789012:function:my-function[:v1]"       (partial ARN)
//
// Inputs longer than maxFunctionRefLength resolve to an empty name so the
// downstream validation rejects them, matching the model's raw-input bound.
//
// Function names cannot contain colons, so the first colon after the name
// separates the embedded qualifier. An explicit Qualifier request parameter
// takes precedence over the embedded one (see mergeQualifier).
func resolveFunctionRef(nameOrArn string) (name, qualifier string) {
	if nameOrArn == "" {
		return "", ""
	}
	if len(nameOrArn) > maxFunctionRefLength {
		return "", ""
	}
	if strings.HasPrefix(nameOrArn, "arn:") {
		resource := arnutil.ExtractResourceFromARN(nameOrArn)
		if strings.HasPrefix(resource, "function:") {
			rest := strings.TrimPrefix(resource, "function:")
			return splitNameQualifier(rest)
		}
		return nameOrArn, ""
	}
	// Partial ARN: account-id:function:name[:qualifier].
	if idx := strings.Index(nameOrArn, ":function:"); idx >= 0 {
		return splitNameQualifier(nameOrArn[idx+len(":function:"):])
	}
	return splitNameQualifier(nameOrArn)
}

// splitNameQualifier splits "name[:qualifier]" at the first colon.
func splitNameQualifier(s string) (name, qualifier string) {
	if idx := strings.Index(s, ":"); idx >= 0 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
}

// mergeQualifier returns the effective qualifier: an explicit Qualifier
// request parameter is more specific to the request than one embedded in
// the function reference, so it wins when both are present.
func mergeQualifier(paramQualifier, embeddedQualifier string) string {
	if paramQualifier != "" {
		return paramQualifier
	}
	return embeddedQualifier
}

// extractFunctionName returns the bare function name from any accepted
// FunctionName form, discarding an embedded qualifier. Callers that need
// the qualifier must use resolveFunctionRef directly.
func extractFunctionName(arnOrName string) string {
	name, _ := resolveFunctionRef(arnOrName)
	return name
}

// repositoryType returns the AWS repository type for a function's code
// source: "ECR" for container image packages, "S3" for zip-based
// packages.
func repositoryType(fn *lambdastore.Function) string {
	if fn.ImageUri != "" || fn.PackageType == "Image" {
		return "ECR"
	}
	return "S3"
}

// resolveAliasTargetVersion resolves the concrete Version that an alias
// points to, applying weighted routing when RoutingConfig is present.
// Returns nil when the target is $LATEST or the version cannot be found
// (callers treat nil as $LATEST).
func resolveAliasTargetVersion(function *lambdastore.Function, alias *lambdastore.Alias) *lambdastore.Version {
	if alias == nil {
		return nil
	}

	// No routing config: resolve to the alias's primary version directly.
	if alias.RoutingConfig == nil || len(alias.RoutingConfig.AdditionalVersionWeights) == 0 {
		return findVersion(function, alias.FunctionVersion)
	}

	// Weighted routing: primary weight = 1.0 - sum(additional weights).
	primaryWeight := 1.0
	for _, w := range alias.RoutingConfig.AdditionalVersionWeights {
		primaryWeight -= w
	}

	// Weighted random selection.
	r := rand.Float64()
	cumulative := primaryWeight

	// Primary version.
	if r < cumulative {
		return findVersion(function, alias.FunctionVersion)
	}

	// Additional weighted versions.
	for versionStr, weight := range alias.RoutingConfig.AdditionalVersionWeights {
		cumulative += weight
		if r < cumulative {
			return findVersion(function, versionStr)
		}
	}

	// Floating-point rounding fallback: primary version.
	return findVersion(function, alias.FunctionVersion)
}

// findVersion looks up a Version by its version string within the function's
// Versions slice.  Returns nil for $LATEST, empty, or not-found versions.
func findVersion(function *lambdastore.Function, versionStr string) *lambdastore.Version {
	if versionStr == "" || versionStr == "$LATEST" {
		return nil
	}
	for i := range function.Versions {
		if function.Versions[i].Version == versionStr {
			return &function.Versions[i]
		}
	}
	return nil
}

func parseVpcConfig(params map[string]interface{}) *lambdastore.VpcConfig {
	vpcMap := request.GetMapParam(params, "VpcConfig")
	if vpcMap == nil {
		return nil
	}

	vpcConfig := &lambdastore.VpcConfig{}
	if subnets, ok := vpcMap["SubnetIds"].([]interface{}); ok {
		for _, s := range subnets {
			if str, ok := s.(string); ok {
				vpcConfig.SubnetIds = append(vpcConfig.SubnetIds, str)
			}
		}
	}
	if sgs, ok := vpcMap["SecurityGroupIds"].([]interface{}); ok {
		for _, sg := range sgs {
			if str, ok := sg.(string); ok {
				vpcConfig.SecurityGroupIds = append(vpcConfig.SecurityGroupIds, str)
			}
		}
	}
	return vpcConfig
}

// resolveVpcConfig uses the EC2 invoker to derive the VPC ID from the first
// subnet. AWS Lambda derives the VPC from the subnets automatically.
// Returns an error if the subnet lookup fails so callers can reject the
// request instead of creating a function with an empty VpcId.
// The region parameter ensures subnet lookup targets the request's region,
// not the service default region.
func (s *LambdaService) resolveVpcConfig(ctx context.Context, region string, vpcConfig *lambdastore.VpcConfig) error {
	if s.bus == nil || len(vpcConfig.SubnetIds) == 0 {
		return nil
	}
	ec2 := s.bus.EC2Invoker()
	if ec2 == nil {
		return nil
	}
	vpcId, _, err := ec2.LookupSubnet(ctx, region, vpcConfig.SubnetIds[0])
	if err != nil {
		return fmt.Errorf("subnet %q not found or EC2 service unavailable: %w", vpcConfig.SubnetIds[0], err)
	}
	vpcConfig.VpcId = vpcId
	return nil
}

func parseEnvironment(params map[string]interface{}) *lambdastore.Environment {
	envMap := request.GetMapParam(params, "Environment")
	if envMap == nil {
		return nil
	}

	env := &lambdastore.Environment{}
	if vars, ok := envMap["Variables"].(map[string]interface{}); ok {
		env.Variables = make(map[string]string)
		for k, v := range vars {
			if str, ok := v.(string); ok {
				env.Variables[k] = str
			}
		}
	}
	return env
}

func parseDeadLetterConfig(params map[string]interface{}) (*lambdastore.DeadLetterConfig, error) {
	dlMap := request.GetMapParam(params, "DeadLetterConfig")
	if dlMap == nil {
		return nil, nil
	}

	targetArn := request.GetStringParam(dlMap, "TargetArn")
	if targetArn == "" {
		return nil, nil
	}

	if !isValidDeadLetterTargetArn(targetArn) {
		return nil, ErrInvalidParameterValue
	}

	return &lambdastore.DeadLetterConfig{
		TargetArn: targetArn,
	}, nil
}

func isValidDeadLetterTargetArn(arn string) bool {
	svc := arnutil.GetServiceFromARN(arn)
	return svc == "sqs" || svc == "sns"
}

func parseTracingConfig(params map[string]interface{}) (*lambdastore.TracingConfig, error) {
	traceMap := request.GetMapParam(params, "TracingConfig")
	if traceMap == nil {
		return nil, nil
	}

	mode := request.GetStringParam(traceMap, "Mode")
	if mode == "" {
		return nil, nil
	}

	if mode != "Active" && mode != "PassThrough" {
		return nil, ErrInvalidParameterValue
	}

	return &lambdastore.TracingConfig{
		Mode: mode,
	}, nil
}

func parseLoggingConfig(m map[string]interface{}) *lambdastore.LoggingConfig {
	lc := &lambdastore.LoggingConfig{}
	if v, ok := m["LogFormat"].(string); ok {
		lc.LogFormat = v
	}
	if v, ok := m["ApplicationLogLevel"].(string); ok {
		lc.ApplicationLogLevel = v
	}
	if v, ok := m["SystemLogLevel"].(string); ok {
		lc.SystemLogLevel = v
	}
	if v, ok := m["LogGroup"].(string); ok {
		lc.LogGroup = v
	}
	if lc.LogFormat == "" && lc.ApplicationLogLevel == "" && lc.SystemLogLevel == "" && lc.LogGroup == "" {
		return nil
	}
	return lc
}

func parseImageConfig(m map[string]interface{}) *lambdastore.ImageConfig {
	ic := &lambdastore.ImageConfig{}
	if eps, ok := m["EntryPoint"].([]interface{}); ok {
		for _, ep := range eps {
			if s, ok := ep.(string); ok {
				ic.EntryPoint = append(ic.EntryPoint, s)
			}
		}
	}
	if cmds, ok := m["Command"].([]interface{}); ok {
		for _, c := range cmds {
			if s, ok := c.(string); ok {
				ic.Command = append(ic.Command, s)
			}
		}
	}
	if wd, ok := m["WorkingDirectory"].(string); ok {
		ic.WorkingDirectory = wd
	}
	if len(ic.EntryPoint) == 0 && len(ic.Command) == 0 && ic.WorkingDirectory == "" {
		return nil
	}
	return ic
}

func deepCopyFunction(fn *lambdastore.Function) *lambdastore.Function {
	if fn == nil {
		return nil
	}

	result := &lambdastore.Function{
		FunctionName:               fn.FunctionName,
		FunctionArn:                fn.FunctionArn,
		Runtime:                    fn.Runtime,
		Role:                       fn.Role,
		Handler:                    fn.Handler,
		CodeSize:                   fn.CodeSize,
		CodeSha256:                 fn.CodeSha256,
		CodeLocation:               fn.CodeLocation,
		ImageUri:                   fn.ImageUri,
		SourceCodeHash:             fn.SourceCodeHash,
		Description:                fn.Description,
		Timeout:                    fn.Timeout,
		MemorySize:                 fn.MemorySize,
		Publish:                    fn.Publish,
		KMSKeyArn:                  fn.KMSKeyArn,
		RevisionId:                 fn.RevisionId,
		State:                      fn.State,
		StateReason:                fn.StateReason,
		StateReasonCode:            fn.StateReasonCode,
		LastUpdateStatus:           fn.LastUpdateStatus,
		LastUpdateReason:           fn.LastUpdateReason,
		LastUpdateStatusReason:     fn.LastUpdateStatusReason,
		LastUpdateStatusReasonCode: fn.LastUpdateStatusReasonCode,
		LastModified:               fn.LastModified,
		LastModifiedUser:           fn.LastModifiedUser,
		PackageType:                fn.PackageType,
		SigningProfileVersionArn:   fn.SigningProfileVersionArn,
		SigningJobArn:              fn.SigningJobArn,
		CodeSigningConfigArn:       fn.CodeSigningConfigArn,
		CurrentVersion:             fn.CurrentVersion,
		ReservedConcurrency:        fn.ReservedConcurrency,
		ContainerID:                fn.ContainerID,
		ContainerImageID:           fn.ContainerImageID,
	}

	if fn.EphemeralStorage != nil {
		result.EphemeralStorage = &lambdastore.EphemeralStorage{Size: fn.EphemeralStorage.Size}
	}

	if len(fn.Architectures) > 0 {
		result.Architectures = make([]string, len(fn.Architectures))
		copy(result.Architectures, fn.Architectures)
	}

	if fn.VpcConfig != nil {
		result.VpcConfig = &lambdastore.VpcConfig{
			VpcId:                   fn.VpcConfig.VpcId,
			Ipv6AllowedForDualStack: fn.VpcConfig.Ipv6AllowedForDualStack,
		}
		if len(fn.VpcConfig.SubnetIds) > 0 {
			result.VpcConfig.SubnetIds = make([]string, len(fn.VpcConfig.SubnetIds))
			copy(result.VpcConfig.SubnetIds, fn.VpcConfig.SubnetIds)
		}
		if len(fn.VpcConfig.SecurityGroupIds) > 0 {
			result.VpcConfig.SecurityGroupIds = make([]string, len(fn.VpcConfig.SecurityGroupIds))
			copy(result.VpcConfig.SecurityGroupIds, fn.VpcConfig.SecurityGroupIds)
		}
	}

	if fn.Environment != nil {
		result.Environment = &lambdastore.Environment{}
		if fn.Environment.Variables != nil {
			result.Environment.Variables = make(map[string]string, len(fn.Environment.Variables))
			for k, v := range fn.Environment.Variables {
				result.Environment.Variables[k] = v
			}
		}
	}

	if fn.DeadLetterConfig != nil {
		result.DeadLetterConfig = &lambdastore.DeadLetterConfig{TargetArn: fn.DeadLetterConfig.TargetArn}
	}

	if fn.TracingConfig != nil {
		result.TracingConfig = &lambdastore.TracingConfig{Mode: fn.TracingConfig.Mode}
	}

	if len(fn.Layers) > 0 {
		result.Layers = make([]lambdastore.LayerReference, len(fn.Layers))
		copy(result.Layers, fn.Layers)
	}

	if fn.SnapStart != nil {
		result.SnapStart = &lambdastore.SnapStart{ApplyOn: fn.SnapStart.ApplyOn}
	}

	if fn.LoggingConfig != nil {
		result.LoggingConfig = &lambdastore.LoggingConfig{
			LogFormat:           fn.LoggingConfig.LogFormat,
			ApplicationLogLevel: fn.LoggingConfig.ApplicationLogLevel,
			SystemLogLevel:      fn.LoggingConfig.SystemLogLevel,
			LogGroup:            fn.LoggingConfig.LogGroup,
		}
	}

	if fn.ImageConfig != nil {
		result.ImageConfig = &lambdastore.ImageConfig{
			WorkingDirectory: fn.ImageConfig.WorkingDirectory,
		}
		if len(fn.ImageConfig.EntryPoint) > 0 {
			result.ImageConfig.EntryPoint = make([]string, len(fn.ImageConfig.EntryPoint))
			copy(result.ImageConfig.EntryPoint, fn.ImageConfig.EntryPoint)
		}
		if len(fn.ImageConfig.Command) > 0 {
			result.ImageConfig.Command = make([]string, len(fn.ImageConfig.Command))
			copy(result.ImageConfig.Command, fn.ImageConfig.Command)
		}
	}

	if len(fn.FileSystemConfigs) > 0 {
		result.FileSystemConfigs = make([]lambdastore.FileSystemConfig, len(fn.FileSystemConfigs))
		copy(result.FileSystemConfigs, fn.FileSystemConfigs)
	}

	if fn.UrlConfig != nil {
		result.UrlConfig = deepCopyFunctionUrlConfig(fn.UrlConfig)
	}

	if len(fn.Versions) > 0 {
		result.Versions = make([]lambdastore.Version, len(fn.Versions))
		for i, v := range fn.Versions {
			result.Versions[i] = *deepCopyVersion(&v)
		}
	}

	if len(fn.Aliases) > 0 {
		result.Aliases = make([]lambdastore.Alias, len(fn.Aliases))
		copy(result.Aliases, fn.Aliases)
	}

	if len(fn.Policies) > 0 {
		result.Policies = make([]lambdastore.FunctionPolicy, len(fn.Policies))
		copy(result.Policies, fn.Policies)
	}

	if len(fn.ProvisionedConcurrency) > 0 {
		result.ProvisionedConcurrency = make([]lambdastore.ProvisionedConcurrencyConfig, len(fn.ProvisionedConcurrency))
		copy(result.ProvisionedConcurrency, fn.ProvisionedConcurrency)
	}

	if len(fn.EventInvokeConfigs) > 0 {
		result.EventInvokeConfigs = make([]lambdastore.EventInvokeConfig, len(fn.EventInvokeConfigs))
		copy(result.EventInvokeConfigs, fn.EventInvokeConfigs)
	}

	return result
}

func deepCopyVersion(v *lambdastore.Version) *lambdastore.Version {
	if v == nil {
		return nil
	}

	result := &lambdastore.Version{
		Version:                  v.Version,
		FunctionArn:              v.FunctionArn,
		Runtime:                  v.Runtime,
		Role:                     v.Role,
		Handler:                  v.Handler,
		CodeSize:                 v.CodeSize,
		CodeSha256:               v.CodeSha256,
		CodeLocation:             v.CodeLocation,
		ImageUri:                 v.ImageUri,
		Description:              v.Description,
		Timeout:                  v.Timeout,
		MemorySize:               v.MemorySize,
		KMSKeyArn:                v.KMSKeyArn,
		RevisionId:               v.RevisionId,
		State:                    v.State,
		StateReason:              v.StateReason,
		StateReasonCode:          v.StateReasonCode,
		LastUpdateStatus:         v.LastUpdateStatus,
		LastModified:             v.LastModified,
		PackageType:              v.PackageType,
		SigningProfileVersionArn: v.SigningProfileVersionArn,
		SigningJobArn:            v.SigningJobArn,
	}

	if v.EphemeralStorage != nil {
		result.EphemeralStorage = &lambdastore.EphemeralStorage{Size: v.EphemeralStorage.Size}
	}

	if len(v.Architectures) > 0 {
		result.Architectures = make([]string, len(v.Architectures))
		copy(result.Architectures, v.Architectures)
	}

	if v.VpcConfig != nil {
		result.VpcConfig = &lambdastore.VpcConfig{
			VpcId:                   v.VpcConfig.VpcId,
			Ipv6AllowedForDualStack: v.VpcConfig.Ipv6AllowedForDualStack,
		}
		if len(v.VpcConfig.SubnetIds) > 0 {
			result.VpcConfig.SubnetIds = make([]string, len(v.VpcConfig.SubnetIds))
			copy(result.VpcConfig.SubnetIds, v.VpcConfig.SubnetIds)
		}
		if len(v.VpcConfig.SecurityGroupIds) > 0 {
			result.VpcConfig.SecurityGroupIds = make([]string, len(v.VpcConfig.SecurityGroupIds))
			copy(result.VpcConfig.SecurityGroupIds, v.VpcConfig.SecurityGroupIds)
		}
	}

	if v.Environment != nil && v.Environment.Variables != nil {
		result.Environment = &lambdastore.Environment{
			Variables: make(map[string]string, len(v.Environment.Variables)),
		}
		for k, val := range v.Environment.Variables {
			result.Environment.Variables[k] = val
		}
	}

	if len(v.Layers) > 0 {
		result.Layers = make([]lambdastore.LayerReference, len(v.Layers))
		copy(result.Layers, v.Layers)
	}

	if v.LoggingConfig != nil {
		result.LoggingConfig = &lambdastore.LoggingConfig{
			LogFormat:           v.LoggingConfig.LogFormat,
			ApplicationLogLevel: v.LoggingConfig.ApplicationLogLevel,
			SystemLogLevel:      v.LoggingConfig.SystemLogLevel,
			LogGroup:            v.LoggingConfig.LogGroup,
		}
	}

	if v.ImageConfig != nil {
		result.ImageConfig = &lambdastore.ImageConfig{
			WorkingDirectory: v.ImageConfig.WorkingDirectory,
		}
		if len(v.ImageConfig.EntryPoint) > 0 {
			result.ImageConfig.EntryPoint = make([]string, len(v.ImageConfig.EntryPoint))
			copy(result.ImageConfig.EntryPoint, v.ImageConfig.EntryPoint)
		}
		if len(v.ImageConfig.Command) > 0 {
			result.ImageConfig.Command = make([]string, len(v.ImageConfig.Command))
			copy(result.ImageConfig.Command, v.ImageConfig.Command)
		}
	}

	if len(v.FileSystemConfigs) > 0 {
		result.FileSystemConfigs = make([]lambdastore.FileSystemConfig, len(v.FileSystemConfigs))
		copy(result.FileSystemConfigs, v.FileSystemConfigs)
	}

	return result
}

func deepCopyFunctionUrlConfig(cfg *lambdastore.FunctionUrlConfig) *lambdastore.FunctionUrlConfig {
	if cfg == nil {
		return nil
	}

	result := &lambdastore.FunctionUrlConfig{
		FunctionUrl:      cfg.FunctionUrl,
		FunctionArn:      cfg.FunctionArn,
		AuthType:         cfg.AuthType,
		CreationTime:     cfg.CreationTime,
		LastModifiedTime: cfg.LastModifiedTime,
		Cors:             cfg.Cors,
		InvokeMode:       cfg.InvokeMode,
	}

	return result
}
