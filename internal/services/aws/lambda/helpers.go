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

// maxFunctionNameLength is the create-time bare-name cap: "The name of
// the Lambda function... up to 64 characters" (AWS API reference,
// CreateFunction FunctionName). The model's @length(1,140) admits any
// reference form; this cap binds the bare name a created function can
// carry, and every reference form must resolve against it.
const maxFunctionNameLength = 64

// maxAliasNameLength is the Alias shape's @length(1,128) maximum from the
// Smithy model.
const maxAliasNameLength = 128

// maxNamespacedFunctionRefLength is the NamespacedFunctionName
// @length(1,256) maximum from the Smithy model. Members typed by that
// shape admit references up to 256 characters; the bare name inside any
// form still has to resolve against functions created under the
// 64-character create-time bound, so a longer reference resolves to a
// name no function carries (not found) rather than being rejected.
const maxNamespacedFunctionRefLength = 256

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
	return resolveFunctionRefWithin(nameOrArn, maxFunctionRefLength)
}

// resolveNamespacedFunctionRef is resolveFunctionRef for members the
// Smithy model types as NamespacedFunctionName: identical forms with the
// wider 256-character raw-input bound, so a long but pattern-valid
// reference resolves (and typically misses) instead of being rejected as
// over-long.
func resolveNamespacedFunctionRef(nameOrArn string) (name, qualifier string) {
	return resolveFunctionRefWithin(nameOrArn, maxNamespacedFunctionRefLength)
}

func resolveFunctionRefWithin(nameOrArn string, maxLen int) (name, qualifier string) {
	if nameOrArn == "" {
		return "", ""
	}
	if len(nameOrArn) > maxLen {
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

// parseFileSystemConfigs converts the FileSystemConfigs wire member into
// store values. Presence is preserved: an entry without S3FilesConfig keeps
// a nil pointer, and an empty input list stays empty (an update-provided
// empty list clears the configuration).
func parseFileSystemConfigs(fscs []interface{}) []lambdastore.FileSystemConfig {
	var configs []lambdastore.FileSystemConfig
	for _, fsc := range fscs {
		m, ok := fsc.(map[string]interface{})
		if !ok {
			continue
		}
		config := lambdastore.FileSystemConfig{
			Arn:            request.GetStringParam(m, "Arn"),
			LocalMountPath: request.GetStringParam(m, "LocalMountPath"),
		}
		if s3m := request.GetMapParam(m, "S3FilesConfig"); s3m != nil {
			config.S3FilesConfig = &lambdastore.S3FilesConfig{
				DirectS3Read: request.GetStringParam(s3m, "DirectS3Read"),
			}
		}
		configs = append(configs, config)
	}
	return configs
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
