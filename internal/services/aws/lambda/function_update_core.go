package lambda

import (
	"context"
	"encoding/base64"
	"errors"

	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// UpdateFunctionCodeInput carries every field that UpdateFunctionCode needs,
// in a format independent of the wire protocol. Code metadata (S3 fetch,
// base64 decode, disk storage) is pre-processed by the caller.
type UpdateFunctionCodeInput struct {
	FunctionName  string
	CodeLocation  string
	CodeSize      int64
	CodeSha256    string
	ImageUri      string
	Architectures []string
	Publish       bool

	// Region of the request; used to persist per-version code snapshots
	// when Publish is requested.
	Region string

	// RevisionId is an optional precondition: when set, the update fails
	// with ResourceConflictException if it does not match the function's
	// current revision.
	RevisionId string
}

// UpdateFunctionConfigurationInput carries every field that
// UpdateFunctionConfiguration needs. Optional configuration structs are
// nil when not provided by the caller.
type UpdateFunctionConfigurationInput struct {
	FunctionName         string
	Runtime              string
	Role                 string
	Handler              string
	Description          string
	Timeout              int32
	MemorySize           int32
	KMSKeyArn            string
	CodeSigningConfigArn string
	VpcConfig            *lambdastore.VpcConfig
	Environment          *lambdastore.Environment
	DeadLetterConfig     *lambdastore.DeadLetterConfig
	TracingConfig        *lambdastore.TracingConfig
	LoggingConfig        *lambdastore.LoggingConfig
	ImageConfig          *lambdastore.ImageConfig
	EphemeralStorage     *lambdastore.EphemeralStorage
	SnapStart            *lambdastore.SnapStart
	FileSystemConfigs    []lambdastore.FileSystemConfig
	Layers               []lambdastore.LayerReference
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// prepareFunctionCodeUpdateCore resolves the wire code members of an
// UpdateFunctionCode request into persisted code metadata. UpdateFunctionCode
// carries the code members at the top level of the request (only
// CreateFunction nests them under Code); a ZipFile member is decoded in
// place while an S3 bucket reference fetches the archive. Both are
// persisted under the function's $LATEST code directory with their hash
// recorded.
func (s *LambdaService) prepareFunctionCodeUpdateCore(ctx context.Context, region, functionName, zipFileStr, imageUri, s3Bucket, s3Key, s3Version string) (*functionCodeMetadata, error) {
	if zipFileStr == "" && imageUri == "" && s3Bucket == "" {
		return nil, NewInvalidParameter("Code", "Either ZipFile, ImageUri, or S3Bucket/S3Key must be provided")
	}

	meta := &functionCodeMetadata{}
	if zipFileStr != "" {
		zipFile, err := base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, NewInvalidParameter("ZipFile", "Invalid base64 encoding")
		}
		codeLocation, codeSize, err := s.storeCode(functionName, "$LATEST", zipFile, region)
		if err != nil {
			return nil, err
		}
		meta.CodeLocation, meta.CodeSize = codeLocation, codeSize
		meta.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
	} else if s3Bucket != "" {
		if s3Key == "" {
			return nil, NewInvalidParameter("Code.S3Key", "S3Key is required when S3Bucket is specified")
		}
		zipFile, err := s.fetchCodeFromS3(ctx, s3Bucket, s3Key, s3Version, region)
		if err != nil {
			return nil, NewInvalidParameter("Code", err.Error())
		}
		codeLocation, codeSize, err := s.storeCode(functionName, "$LATEST", zipFile, region)
		if err != nil {
			return nil, err
		}
		meta.CodeLocation, meta.CodeSize = codeLocation, codeSize
		meta.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
	}

	return meta, nil
}

// getFunctionForDryRunCore fetches a function for the UpdateFunctionCode
// DryRun branch, which validates the request without modifying the code.
func (s *LambdaService) getFunctionForDryRunCore(stores *lambdaStore, functionName string) (*lambdastore.Function, error) {
	current, err := stores.Functions.Get(functionName)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return current, nil
}

// updateFunctionCodeCore is the single entry point for function code update
// logic shared by the HTTP API and the admin gRPC handler. It performs
// architecture validation, applies the code update atomically, and
// optionally publishes a new version. The returned Version is non-nil only
// when in.Publish requested a publish; callers use it to answer with the
// published version's configuration.
func (s *LambdaService) updateFunctionCodeCore(stores *lambdaStore, in *UpdateFunctionCodeInput) (*lambdastore.Function, *lambdastore.Version, error) {
	functionName := in.FunctionName

	for _, arch := range in.Architectures {
		if err := validateArchitecture(arch); err != nil {
			return nil, nil, err
		}
	}

	function, err := stores.Functions.UpdateAtomically(functionName, func(fn *lambdastore.Function) error {
		if in.RevisionId != "" && fn.RevisionId != in.RevisionId {
			return NewResourceConflict("The RevisionId provided does not match the current revision of the function")
		}
		if in.CodeLocation != "" {
			fn.CodeLocation = in.CodeLocation
			fn.CodeSize = in.CodeSize
			fn.CodeSha256 = in.CodeSha256
		}

		if in.ImageUri != "" {
			fn.ImageUri = in.ImageUri
			fn.PackageType = "Image"
		}

		if len(in.Architectures) > 0 {
			fn.Architectures = make([]string, 0, len(in.Architectures))
			fn.Architectures = append(fn.Architectures, in.Architectures...)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return nil, nil, NewResourceNotFound("Function", functionName)
		}
		return nil, nil, err
	}

	var published *lambdastore.Version
	if in.Publish {
		published, err = s.publishVersionWithCode(stores, function, "", in.Region)
		if err != nil {
			return nil, nil, err
		}
	}

	return function, published, nil
}

// updateFunctionConfigurationCore is the single entry point for function
// configuration update logic shared by the HTTP API and the admin gRPC
// handler. It performs all field validation, applies the update atomically,
// and cleans up any running container so the next invoke creates a fresh one.
func (s *LambdaService) updateFunctionConfigurationCore(ctx context.Context, stores *lambdaStore, in *UpdateFunctionConfigurationInput) (*lambdastore.Function, error) {
	functionName := in.FunctionName

	if in.Runtime != "" && !ValidateRuntime(in.Runtime) {
		return nil, NewInvalidParameter("Runtime", "Runtime '"+in.Runtime+"' is not supported")
	}

	if in.Timeout > 0 {
		if err := validateTimeout(in.Timeout); err != nil {
			return nil, err
		}
	}
	if in.MemorySize > 0 {
		if err := validateMemorySize(in.MemorySize); err != nil {
			return nil, err
		}
	}

	if err := validateCodeSigningConfigArn(in.CodeSigningConfigArn); err != nil {
		return nil, err
	}

	if err := validateKMSKeyArn(in.KMSKeyArn); err != nil {
		return nil, err
	}

	if in.EphemeralStorage != nil {
		if err := validateEphemeralStorageSize(in.EphemeralStorage.Size); err != nil {
			return nil, err
		}
	}
	if in.SnapStart != nil {
		if err := validateSnapStartApplyOn(in.SnapStart.ApplyOn); err != nil {
			return nil, err
		}
	}
	if err := validateEnvironmentVariables(in.Environment); err != nil {
		return nil, err
	}

	for _, lr := range in.Layers {
		if !isValidLayerARN(lr.Arn) {
			return nil, NewInvalidParameter("Layers", "Invalid layer ARN format: "+lr.Arn)
		}
	}

	var oldContainerID string

	function, err := stores.Functions.UpdateAtomically(functionName, func(fn *lambdastore.Function) error {
		// SnapStart support depends on the effective runtime after this
		// update, so the guard runs where the target state is known.
		if in.SnapStart != nil {
			effectiveRuntime := in.Runtime
			if effectiveRuntime == "" {
				effectiveRuntime = string(fn.Runtime)
			}
			if err := validateSnapStartForRuntime(effectiveRuntime, in.SnapStart); err != nil {
				return err
			}
		}
		if in.Runtime != "" {
			fn.Runtime = lambdastore.Runtime(in.Runtime)
		}
		if in.Role != "" {
			fn.Role = in.Role
		}
		if in.Handler != "" {
			fn.Handler = in.Handler
		}
		if in.Description != "" {
			fn.Description = in.Description
		}
		if in.Timeout > 0 {
			fn.Timeout = in.Timeout
		}
		if in.MemorySize > 0 {
			fn.MemorySize = in.MemorySize
		}
		if in.KMSKeyArn != "" {
			fn.KMSKeyArn = in.KMSKeyArn
		}
		if in.CodeSigningConfigArn != "" {
			fn.CodeSigningConfigArn = in.CodeSigningConfigArn
		}
		if in.VpcConfig != nil {
			fn.VpcConfig = in.VpcConfig
		}
		if in.Environment != nil {
			fn.Environment = in.Environment
		}
		if in.DeadLetterConfig != nil {
			fn.DeadLetterConfig = in.DeadLetterConfig
		}
		if in.TracingConfig != nil {
			fn.TracingConfig = in.TracingConfig
		}
		if in.LoggingConfig != nil {
			fn.LoggingConfig = in.LoggingConfig
		}
		if in.ImageConfig != nil {
			fn.ImageConfig = in.ImageConfig
		}
		if in.EphemeralStorage != nil {
			fn.EphemeralStorage = in.EphemeralStorage
		}
		if in.SnapStart != nil {
			fn.SnapStart = in.SnapStart
		}
		if in.FileSystemConfigs != nil {
			fn.FileSystemConfigs = in.FileSystemConfigs
		}
		if in.Layers != nil {
			fn.Layers = in.Layers
		}

		// Invalidate any running container so the next invoke creates a fresh
		// one with the updated runtime, handler, memory, or environment.
		oldContainerID = fn.ContainerID
		fn.ContainerID = ""
		fn.ContainerImageID = ""
		return nil
	})
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return nil, NewResourceNotFound("Function", functionName)
		}
		return nil, err
	}

	// Remove the previous container to prevent orphaned Docker containers.
	// The ID was captured inside the callback before it was cleared.
	if oldContainerID != "" {
		if rmErr := s.dockerClient.RemoveContainer(ctx, oldContainerID, true); rmErr != nil {
			logs.Warn("Failed to remove container after configuration update",
				logs.String("containerID", oldContainerID),
				logs.String("function", functionName),
				logs.Err(rmErr))
		}
	}

	return function, nil
}
