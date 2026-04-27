// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

var functionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

func extractFunctionName(arnOrName string) string {
	if strings.HasPrefix(arnOrName, "arn:") {
		if name := arnutil.ExtractFunctionNameFromARN(arnOrName); name != "" {
			return name
		}
	}
	return arnOrName
}

func validateFunctionName(name string) error {
	if len(name) == 0 || len(name) > 64 {
		return NewInvalidParameter("FunctionName", "Function name must be between 1 and 64 characters")
	}
	if !functionNameRegex.MatchString(name) {
		return NewInvalidParameter("FunctionName", "Function name can only contain alphanumeric characters, hyphens, and underscores")
	}
	return nil
}

func isValidLayerARN(arnStr string) bool {
	if arnStr == "" {
		return false
	}
	if _, service, _, _, _ := arnutil.SplitARN(arnStr); service != "lambda" {
		return false
	}
	resource := arnutil.ExtractResourceFromARN(arnStr)
	if resource == "" {
		return false
	}
	return strings.HasPrefix(resource, "layer:")
}

func validateTimeout(timeout int32) error {
	if timeout < 1 || timeout > 900 {
		return NewInvalidParameter("Timeout", "Timeout must be between 1 and 900 seconds")
	}
	return nil
}

func validateMemorySize(memorySize int32) error {
	if memorySize < 128 || memorySize > 10240 {
		return NewInvalidParameter("MemorySize", "MemorySize must be between 128 and 10240 MB")
	}
	return nil
}

func (s *LambdaService) validateAndGetFunction(ctx *request.RequestContext, params map[string]interface{}) (*lambdastore.Function, error) {
	functionName := request.GetStringParam(params, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)

	store, err := s.store(ctx)
	if err != nil {
		return nil, err
	}
	function, err := store.Functions.Get(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	return function, nil
}

func (s *LambdaService) validateAndGetFunctionWithQualifier(ctx *request.RequestContext, params map[string]interface{}) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, error) {
	functionName := request.GetStringParam(params, "FunctionName")
	if functionName == "" {
		return nil, nil, nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)

	qualifier := request.GetStringParam(params, "Qualifier")
	store, err := s.store(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	return s.resolveQualifier(store.Functions, functionName, qualifier)
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
func (s *LambdaService) resolveVpcConfig(ctx context.Context, vpcConfig *lambdastore.VpcConfig) {
	if s.bus == nil || len(vpcConfig.SubnetIds) == 0 {
		return
	}
	ec2 := s.bus.EC2Invoker()
	if ec2 == nil {
		return
	}
	vpcId, _, err := ec2.LookupSubnet(ctx, s.region, vpcConfig.SubnetIds[0])
	if err != nil {
		return
	}
	vpcConfig.VpcId = vpcId
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

func deepCopyFunction(fn *lambdastore.Function) *lambdastore.Function {
	if fn == nil {
		return nil
	}

	result := &lambdastore.Function{
		FunctionName:             fn.FunctionName,
		FunctionArn:              fn.FunctionArn,
		Runtime:                  fn.Runtime,
		Role:                     fn.Role,
		Handler:                  fn.Handler,
		CodeSize:                 fn.CodeSize,
		CodeSha256:               fn.CodeSha256,
		CodeLocation:             fn.CodeLocation,
		ImageUri:                 fn.ImageUri,
		SourceCodeHash:           fn.SourceCodeHash,
		Description:              fn.Description,
		Timeout:                  fn.Timeout,
		MemorySize:               fn.MemorySize,
		Publish:                  fn.Publish,
		KMSKeyArn:                fn.KMSKeyArn,
		RevisionId:               fn.RevisionId,
		State:                    fn.State,
		StateReason:              fn.StateReason,
		StateReasonCode:          fn.StateReasonCode,
		LastUpdateStatus:         fn.LastUpdateStatus,
		LastUpdateReason:         fn.LastUpdateReason,
		LastUpdateStatusReason:   fn.LastUpdateStatusReason,
		LastModified:             fn.LastModified,
		LastModifiedUser:         fn.LastModifiedUser,
		PackageType:              fn.PackageType,
		SigningProfileVersionArn: fn.SigningProfileVersionArn,
		SigningJobArn:            fn.SigningJobArn,
		CodeSigningConfigArn:     fn.CodeSigningConfigArn,
		CurrentVersion:           fn.CurrentVersion,
		ReservedConcurrency:      fn.ReservedConcurrency,
		ContainerID:              fn.ContainerID,
		ContainerImageID:         fn.ContainerImageID,
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
			VpcId: fn.VpcConfig.VpcId,
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
