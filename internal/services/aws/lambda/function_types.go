package lambda

import (
	"time"

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

// configFields holds the common fields shared between Function and Version
// configurations, used to eliminate copy-paste in response building.
type configFields struct {
	FunctionName               string
	FunctionArn                string
	Role                       string
	CodeSize                   int64
	Description                string
	Timeout                    int32
	MemorySize                 int32
	LastModified               time.Time
	CodeSha256                 string
	Version                    string
	RevisionId                 string
	State                      string
	StateReason                string
	StateReasonCode            string
	LastUpdateStatus           string
	LastUpdateStatusReason     string
	LastUpdateStatusReasonCode string
	PackageType                string
	Runtime                    string
	Handler                    string
	KMSKeyArn                  string
	ImageUri                   string
	EphemeralStorage           *lambdastore.EphemeralStorage
	VpcConfig                  *lambdastore.VpcConfig
	Environment                *lambdastore.Environment
	TracingConfig              *lambdastore.TracingConfig
	DeadLetterConfig           *lambdastore.DeadLetterConfig
	SnapStart                  *lambdastore.SnapStart
	Architectures              []string
	Layers                     []lambdastore.LayerReference
	LoggingConfig              *lambdastore.LoggingConfig
	ImageConfig                *lambdastore.ImageConfig
	FileSystemConfigs          []lambdastore.FileSystemConfig
	SigningProfileVersionArn   string
	SigningJobArn              string
}

func functionToConfigFields(fn *lambdastore.Function) configFields {
	return configFields{
		FunctionName:               fn.FunctionName,
		FunctionArn:                fn.FunctionArn,
		Role:                       fn.Role,
		CodeSize:                   fn.CodeSize,
		Description:                fn.Description,
		Timeout:                    fn.Timeout,
		MemorySize:                 fn.MemorySize,
		LastModified:               fn.LastModified,
		CodeSha256:                 fn.CodeSha256,
		Version:                    fn.CurrentVersion,
		RevisionId:                 fn.RevisionId,
		State:                      string(fn.State),
		StateReason:                fn.StateReason,
		StateReasonCode:            fn.StateReasonCode,
		LastUpdateStatus:           string(fn.LastUpdateStatus),
		LastUpdateStatusReason:     fn.LastUpdateStatusReason,
		LastUpdateStatusReasonCode: fn.LastUpdateStatusReasonCode,
		PackageType:                fn.PackageType,
		Runtime:                    string(fn.Runtime),
		Handler:                    fn.Handler,
		KMSKeyArn:                  fn.KMSKeyArn,
		ImageUri:                   fn.ImageUri,
		EphemeralStorage:           fn.EphemeralStorage,
		VpcConfig:                  fn.VpcConfig,
		Environment:                fn.Environment,
		TracingConfig:              fn.TracingConfig,
		DeadLetterConfig:           fn.DeadLetterConfig,
		SnapStart:                  fn.SnapStart,
		Architectures:              fn.Architectures,
		Layers:                     fn.Layers,
		LoggingConfig:              fn.LoggingConfig,
		ImageConfig:                fn.ImageConfig,
		FileSystemConfigs:          fn.FileSystemConfigs,
		SigningProfileVersionArn:   fn.SigningProfileVersionArn,
		SigningJobArn:              fn.SigningJobArn,
	}
}

func versionToConfigFields(v *lambdastore.Version) configFields {
	return configFields{
		FunctionName:             arn.ExtractFunctionNameFromARN(v.FunctionArn),
		FunctionArn:              v.FunctionArn,
		Role:                     v.Role,
		CodeSize:                 v.CodeSize,
		Description:              v.Description,
		Timeout:                  v.Timeout,
		MemorySize:               v.MemorySize,
		LastModified:             v.LastModified,
		CodeSha256:               v.CodeSha256,
		Version:                  v.Version,
		RevisionId:               v.RevisionId,
		State:                    string(v.State),
		StateReason:              v.StateReason,
		StateReasonCode:          v.StateReasonCode,
		LastUpdateStatus:         string(v.LastUpdateStatus),
		PackageType:              v.PackageType,
		Runtime:                  string(v.Runtime),
		Handler:                  v.Handler,
		KMSKeyArn:                v.KMSKeyArn,
		ImageUri:                 v.ImageUri,
		EphemeralStorage:         v.EphemeralStorage,
		VpcConfig:                v.VpcConfig,
		Environment:              v.Environment,
		TracingConfig:            v.TracingConfig,
		DeadLetterConfig:         v.DeadLetterConfig,
		SnapStart:                v.SnapStart,
		Architectures:            v.Architectures,
		Layers:                   v.Layers,
		LoggingConfig:            v.LoggingConfig,
		ImageConfig:              v.ImageConfig,
		FileSystemConfigs:        v.FileSystemConfigs,
		SigningProfileVersionArn: v.SigningProfileVersionArn,
		SigningJobArn:            v.SigningJobArn,
	}
}

func buildConfigMap(f configFields) map[string]interface{} {
	config := map[string]interface{}{
		"FunctionName":     f.FunctionName,
		"FunctionArn":      f.FunctionArn,
		"Role":             f.Role,
		"CodeSize":         f.CodeSize,
		"Description":      f.Description,
		"Timeout":          f.Timeout,
		"MemorySize":       f.MemorySize,
		"LastModified":     f.LastModified.Format(timeutils.ISO8601UTCFormat),
		"CodeSha256":       f.CodeSha256,
		"Version":          f.Version,
		"RevisionId":       f.RevisionId,
		"State":            f.State,
		"LastUpdateStatus": f.LastUpdateStatus,
		"PackageType":      f.PackageType,
	}

	if f.StateReason != "" {
		config["StateReason"] = f.StateReason
	}
	if f.StateReasonCode != "" {
		config["StateReasonCode"] = f.StateReasonCode
	}
	if f.LastUpdateStatusReason != "" {
		config["LastUpdateStatusReason"] = f.LastUpdateStatusReason
	}
	if f.LastUpdateStatusReasonCode != "" {
		config["LastUpdateStatusReasonCode"] = f.LastUpdateStatusReasonCode
	}

	if f.PackageType != "Image" {
		config["Runtime"] = f.Runtime
		config["Handler"] = f.Handler
	}

	if f.EphemeralStorage != nil {
		config["EphemeralStorage"] = map[string]interface{}{
			"Size": f.EphemeralStorage.Size,
		}
	}

	if f.VpcConfig != nil {
		config["VpcConfig"] = map[string]interface{}{
			"SubnetIds":        f.VpcConfig.SubnetIds,
			"SecurityGroupIds": f.VpcConfig.SecurityGroupIds,
			"VpcId":            f.VpcConfig.VpcId,
		}
	}

	if f.Environment != nil {
		config["Environment"] = map[string]interface{}{
			"Variables": f.Environment.Variables,
		}
	}

	if f.TracingConfig != nil {
		config["TracingConfig"] = map[string]interface{}{
			"Mode": f.TracingConfig.Mode,
		}
	}

	if f.DeadLetterConfig != nil {
		config["DeadLetterConfig"] = map[string]interface{}{
			"TargetArn": f.DeadLetterConfig.TargetArn,
		}
	}

	if f.KMSKeyArn != "" {
		config["KMSKeyArn"] = f.KMSKeyArn
	}

	if f.SnapStart != nil {
		snapStart := map[string]interface{}{
			"ApplyOn": f.SnapStart.ApplyOn,
		}
		if f.SnapStart.OptimizationStatus != "" {
			snapStart["OptimizationStatus"] = f.SnapStart.OptimizationStatus
		}
		config["SnapStart"] = snapStart
	}

	if len(f.Architectures) > 0 {
		config["Architectures"] = f.Architectures
	}

	if len(f.Layers) > 0 {
		layers := make([]map[string]interface{}, 0, len(f.Layers))
		for _, l := range f.Layers {
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

	if f.ImageUri != "" {
		config["ImageUri"] = f.ImageUri
	}

	if f.LoggingConfig != nil {
		lc := map[string]interface{}{}
		if f.LoggingConfig.LogFormat != "" {
			lc["LogFormat"] = f.LoggingConfig.LogFormat
		}
		if f.LoggingConfig.ApplicationLogLevel != "" {
			lc["ApplicationLogLevel"] = f.LoggingConfig.ApplicationLogLevel
		}
		if f.LoggingConfig.SystemLogLevel != "" {
			lc["SystemLogLevel"] = f.LoggingConfig.SystemLogLevel
		}
		if f.LoggingConfig.LogGroup != "" {
			lc["LogGroup"] = f.LoggingConfig.LogGroup
		}
		config["LoggingConfig"] = lc
	}

	if f.ImageConfig != nil {
		ic := map[string]interface{}{}
		if len(f.ImageConfig.EntryPoint) > 0 {
			ic["EntryPoint"] = f.ImageConfig.EntryPoint
		}
		if len(f.ImageConfig.Command) > 0 {
			ic["Command"] = f.ImageConfig.Command
		}
		if f.ImageConfig.WorkingDirectory != "" {
			ic["WorkingDirectory"] = f.ImageConfig.WorkingDirectory
		}
		config["ImageConfig"] = ic
	}

	if len(f.FileSystemConfigs) > 0 {
		fscs := make([]map[string]interface{}, 0, len(f.FileSystemConfigs))
		for _, fsc := range f.FileSystemConfigs {
			entry := map[string]interface{}{}
			if fsc.Arn != "" {
				entry["Arn"] = fsc.Arn
			}
			if fsc.LocalMountPath != "" {
				entry["LocalMountPath"] = fsc.LocalMountPath
			}
			fscs = append(fscs, entry)
		}
		config["FileSystemConfigs"] = fscs
	}

	if f.SigningProfileVersionArn != "" {
		config["SigningProfileVersionArn"] = f.SigningProfileVersionArn
	}
	if f.SigningJobArn != "" {
		config["SigningJobArn"] = f.SigningJobArn
	}

	return config
}

func (s *LambdaService) toFunctionConfiguration(fn *lambdastore.Function) map[string]interface{} {
	return buildConfigMap(functionToConfigFields(fn))
}

func (s *LambdaService) toVersionConfiguration(v *lambdastore.Version) map[string]interface{} {
	return buildConfigMap(versionToConfigFields(v))
}
