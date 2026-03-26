package lambda

import (
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateFunctionInput is the request structure for creating a new Lambda function.
type CreateFunctionInput struct {
	FunctionName     string            `json:"FunctionName"`
	Runtime          string            `json:"Runtime"`
	Role             string            `json:"Role"`
	Handler          string            `json:"Handler"`
	Code             *FunctionCode     `json:"Code"`
	Description      string            `json:"Description,omitempty"`
	Timeout          int32             `json:"Timeout,omitempty"`
	MemorySize       int32             `json:"MemorySize,omitempty"`
	Publish          bool              `json:"Publish,omitempty"`
	VpcConfig        *VpcConfig        `json:"VpcConfig,omitempty"`
	Environment      *Environment      `json:"Environment,omitempty"`
	DeadLetterConfig *DeadLetterConfig `json:"DeadLetterConfig,omitempty"`
	KMSKeyArn        string            `json:"KMSKeyArn,omitempty"`
	TracingConfig    *TracingConfig    `json:"TracingConfig,omitempty"`
	Tags             map[string]string `json:"Tags,omitempty"`
	Layers           []string          `json:"Layers,omitempty"`
	Architectures    []string          `json:"Architectures,omitempty"`
	EphemeralStorage *EphemeralStorage `json:"EphemeralStorage,omitempty"`
	SnapStart        *SnapStart        `json:"SnapStart,omitempty"`
	PackageType      string            `json:"PackageType,omitempty"`
}

// FunctionCode contains the location of the function's deployment package.
type FunctionCode struct {
	ZipFile         []byte `json:"ZipFile,omitempty"`
	S3Bucket        string `json:"S3Bucket,omitempty"`
	S3Key           string `json:"S3Key,omitempty"`
	S3ObjectVersion string `json:"S3ObjectVersion,omitempty"`
	ImageUri        string `json:"ImageUri,omitempty"`
	SourceCodeHash  string `json:"SourceCodeHash,omitempty"`
}

// VpcConfig configures the VPC settings for a Lambda function.
type VpcConfig struct {
	SubnetIds               []string `json:"SubnetIds,omitempty"`
	SecurityGroupIds        []string `json:"SecurityGroupIds,omitempty"`
	Ipv6AllowedForDualStack bool     `json:"Ipv6AllowedForDualStack,omitempty"`
}

// Environment defines the environment variables for a Lambda function.
type Environment struct {
	Variables map[string]string `json:"Variables,omitempty"`
}

// DeadLetterConfig defines the dead letter queue configuration for a Lambda function.
type DeadLetterConfig struct {
	TargetArn string `json:"TargetArn,omitempty"`
}

// TracingConfig defines the AWS X-Ray tracing configuration for a Lambda function.
type TracingConfig struct {
	Mode string `json:"Mode,omitempty"`
}

// EphemeralStorage defines the ephemeral storage configuration for a Lambda function.
type EphemeralStorage struct {
	Size int32 `json:"Size"`
}

// SnapStart defines the SnapStart configuration for a Lambda function.
type SnapStart struct {
	ApplyOn string `json:"ApplyOn,omitempty"`
}

func (s *LambdaService) toFunctionConfiguration(fn *lambdastore.Function) map[string]interface{} {
	config := map[string]interface{}{
		"FunctionName":     fn.FunctionName,
		"FunctionArn":      fn.FunctionArn,
		"Role":             fn.Role,
		"CodeSize":         fn.CodeSize,
		"Description":      fn.Description,
		"Timeout":          fn.Timeout,
		"MemorySize":       fn.MemorySize,
		"LastModified":     fn.LastModified.Format(timeutils.ISO8601UTCFormat),
		"CodeSha256":       fn.CodeSha256,
		"Version":          fn.CurrentVersion,
		"RevisionId":       fn.RevisionId,
		"State":            string(fn.State),
		"LastUpdateStatus": string(fn.LastUpdateStatus),
		"PackageType":      fn.PackageType,
	}

	if fn.StateReason != "" {
		config["StateReason"] = fn.StateReason
	}
	if fn.StateReasonCode != "" {
		config["StateReasonCode"] = fn.StateReasonCode
	}

	if fn.PackageType != "Image" {
		config["Runtime"] = string(fn.Runtime)
		config["Handler"] = fn.Handler
	}

	if fn.EphemeralStorage != nil {
		config["EphemeralStorage"] = map[string]interface{}{
			"Size": fn.EphemeralStorage.Size,
		}
	}

	if fn.VpcConfig != nil {
		config["VpcConfig"] = map[string]interface{}{
			"SubnetIds":        fn.VpcConfig.SubnetIds,
			"SecurityGroupIds": fn.VpcConfig.SecurityGroupIds,
			"VpcId":            fn.VpcConfig.VpcId,
		}
	}

	if fn.Environment != nil {
		config["Environment"] = map[string]interface{}{
			"Variables": fn.Environment.Variables,
		}
	}

	if fn.TracingConfig != nil {
		config["TracingConfig"] = map[string]interface{}{
			"Mode": fn.TracingConfig.Mode,
		}
	}

	if fn.DeadLetterConfig != nil {
		config["DeadLetterConfig"] = map[string]interface{}{
			"TargetArn": fn.DeadLetterConfig.TargetArn,
		}
	}

	if fn.KMSKeyArn != "" {
		config["KMSKeyArn"] = fn.KMSKeyArn
	}

	if fn.SnapStart != nil {
		snapStart := map[string]interface{}{
			"ApplyOn": fn.SnapStart.ApplyOn,
		}
		if fn.SnapStart.OptimizationStatus != "" {
			snapStart["OptimizationStatus"] = fn.SnapStart.OptimizationStatus
		}
		config["SnapStart"] = snapStart
	}

	if len(fn.Architectures) > 0 {
		config["Architectures"] = fn.Architectures
	}

	if len(fn.Layers) > 0 {
		layers := make([]map[string]interface{}, 0, len(fn.Layers))
		for _, l := range fn.Layers {
			layer := map[string]interface{}{
				"Arn": l.Arn,
			}
			if l.CodeSize > 0 {
				layer["CodeSize"] = l.CodeSize
			}
			layers = append(layers, layer)
		}
		config["Layers"] = layers
	}

	if fn.ImageUri != "" {
		config["ImageUri"] = fn.ImageUri
	}

	return config
}

func (s *LambdaService) toVersionConfiguration(v *lambdastore.Version) map[string]interface{} {
	config := map[string]interface{}{
		"FunctionName":     arn.ExtractFunctionNameFromARN(v.FunctionArn),
		"FunctionArn":      v.FunctionArn,
		"Role":             v.Role,
		"CodeSize":         v.CodeSize,
		"Description":      v.Description,
		"Timeout":          v.Timeout,
		"MemorySize":       v.MemorySize,
		"LastModified":     v.LastModified.Format(timeutils.ISO8601UTCFormat),
		"CodeSha256":       v.CodeSha256,
		"Version":          v.Version,
		"RevisionId":       v.RevisionId,
		"State":            string(v.State),
		"LastUpdateStatus": string(v.LastUpdateStatus),
		"PackageType":      v.PackageType,
	}

	if v.StateReason != "" {
		config["StateReason"] = v.StateReason
	}
	if v.StateReasonCode != "" {
		config["StateReasonCode"] = v.StateReasonCode
	}

	if v.PackageType != "Image" {
		config["Runtime"] = string(v.Runtime)
		config["Handler"] = v.Handler
	}

	if v.EphemeralStorage != nil {
		config["EphemeralStorage"] = map[string]interface{}{
			"Size": v.EphemeralStorage.Size,
		}
	}

	if v.VpcConfig != nil {
		config["VpcConfig"] = map[string]interface{}{
			"SubnetIds":        v.VpcConfig.SubnetIds,
			"SecurityGroupIds": v.VpcConfig.SecurityGroupIds,
			"VpcId":            v.VpcConfig.VpcId,
		}
	}

	if v.Environment != nil {
		config["Environment"] = map[string]interface{}{
			"Variables": v.Environment.Variables,
		}
	}

	if v.TracingConfig != nil {
		config["TracingConfig"] = map[string]interface{}{
			"Mode": v.TracingConfig.Mode,
		}
	}

	if v.DeadLetterConfig != nil {
		config["DeadLetterConfig"] = map[string]interface{}{
			"TargetArn": v.DeadLetterConfig.TargetArn,
		}
	}

	if v.KMSKeyArn != "" {
		config["KMSKeyArn"] = v.KMSKeyArn
	}

	if v.SnapStart != nil {
		snapStart := map[string]interface{}{
			"ApplyOn": v.SnapStart.ApplyOn,
		}
		if v.SnapStart.OptimizationStatus != "" {
			snapStart["OptimizationStatus"] = v.SnapStart.OptimizationStatus
		}
		config["SnapStart"] = snapStart
	}

	if len(v.Architectures) > 0 {
		config["Architectures"] = v.Architectures
	}

	if len(v.Layers) > 0 {
		layers := make([]map[string]interface{}, 0, len(v.Layers))
		for _, l := range v.Layers {
			layer := map[string]interface{}{
				"Arn": l.Arn,
			}
			if l.CodeSize > 0 {
				layer["CodeSize"] = l.CodeSize
			}
			layers = append(layers, layer)
		}
		config["Layers"] = layers
	}

	if v.ImageUri != "" {
		config["ImageUri"] = v.ImageUri
	}

	return config
}
