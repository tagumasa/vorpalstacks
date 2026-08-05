package lambda

import (
	"context"
	"errors"
	"fmt"

	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateFunctionInput carries every field that CreateFunction needs, in a
// format independent of the wire protocol (HTTP Query/JSON vs gRPC-Web).
// Both the HTTP API handler (function_crud.go) and the admin gRPC handler
// (admin_handler.go) build this struct from their respective request
// formats and delegate to createFunctionCore, ensuring that validation
// and persistence follow a single code path.
type CreateFunctionInput struct {
	FunctionName         string
	Runtime              string
	Role                 string
	Handler              string
	Description          string
	PackageType          string
	KMSKeyArn            string
	Publish              bool
	CodeSigningConfigArn string

	// Code metadata — pre-processed by the caller (S3 fetch, base64
	// decode, disk storage). Empty values are valid for the admin
	// handler which creates metadata-only functions.
	CodeLocation string
	CodeSize     int64
	CodeSha256   string
	ImageUri     string

	// Configuration (parsed by caller from request format).
	Timeout           int32 // 0 = default 3s
	MemorySize        int32 // 0 = default 128 MB
	VpcConfig         *lambdastore.VpcConfig
	Environment       *lambdastore.Environment
	DeadLetterConfig  *lambdastore.DeadLetterConfig
	TracingConfig     *lambdastore.TracingConfig
	SnapStart         *lambdastore.SnapStart
	LoggingConfig     *lambdastore.LoggingConfig
	ImageConfig       *lambdastore.ImageConfig
	EphemeralStorage  *lambdastore.EphemeralStorage
	FileSystemConfigs []lambdastore.FileSystemConfig
	Layers            []lambdastore.LayerReference
	Architectures     []string

	Tags map[string]string
}

// DeleteFunctionInput carries the fields needed for DeleteFunction.
type DeleteFunctionInput struct {
	FunctionName string
	Qualifier    string
}

// ListFunctionsInput carries pagination parameters for ListFunctions.
type ListFunctionsInput struct {
	Marker   string
	MaxItems int
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createFunctionCore is the single entry point for function creation logic
// shared by the HTTP API and the admin gRPC handler. It performs all
// field validation, constructs the Function struct, persists it, applies
// tags, and optionally publishes a version.
func (s *LambdaService) createFunctionCore(stores *lambdaStore, in *CreateFunctionInput) (*lambdastore.Function, error) {
	if err := validateFunctionName(in.FunctionName); err != nil {
		return nil, err
	}

	if in.Runtime == "" && in.PackageType != "Image" {
		return nil, NewInvalidParameter("Runtime", "Runtime is required for Zip package type")
	}
	if in.Runtime != "" && !ValidateRuntime(in.Runtime) {
		return nil, NewInvalidParameter("Runtime", "Runtime '"+in.Runtime+"' is not supported")
	}

	if in.Handler == "" && in.PackageType != "Image" {
		return nil, NewInvalidParameter("Handler", "Handler is required for Zip package type")
	}
	if in.Handler != "" && in.Runtime != "" {
		if err := ValidateHandler(in.Runtime, in.Handler); err != nil {
			return nil, err
		}
	}

	if in.Role == "" {
		return nil, NewInvalidParameter("Role", "Role ARN is required")
	}

	if err := validateCodeSigningConfigArn(in.CodeSigningConfigArn); err != nil {
		return nil, err
	}

	timeout := in.Timeout
	if timeout == 0 {
		timeout = 3
	}
	if err := validateTimeout(timeout); err != nil {
		return nil, err
	}

	memorySize := in.MemorySize
	if memorySize == 0 {
		memorySize = 128
	}
	if err := validateMemorySize(memorySize); err != nil {
		return nil, err
	}

	function := &lambdastore.Function{
		FunctionName:         in.FunctionName,
		Runtime:              lambdastore.Runtime(in.Runtime),
		Role:                 in.Role,
		Handler:              in.Handler,
		Description:          in.Description,
		Timeout:              timeout,
		MemorySize:           memorySize,
		Publish:              in.Publish,
		KMSKeyArn:            in.KMSKeyArn,
		PackageType:          in.PackageType,
		CodeSigningConfigArn: in.CodeSigningConfigArn,
		CodeLocation:         in.CodeLocation,
		CodeSize:             in.CodeSize,
		CodeSha256:           in.CodeSha256,
		ImageUri:             in.ImageUri,
		VpcConfig:            in.VpcConfig,
		Environment:          in.Environment,
		DeadLetterConfig:     in.DeadLetterConfig,
		TracingConfig:        in.TracingConfig,
		SnapStart:            in.SnapStart,
		LoggingConfig:        in.LoggingConfig,
		ImageConfig:          in.ImageConfig,
		EphemeralStorage:     in.EphemeralStorage,
		FileSystemConfigs:    in.FileSystemConfigs,
		Layers:               in.Layers,
		Architectures:        in.Architectures,
	}

	if function.ImageUri != "" {
		function.PackageType = "Image"
	}

	created, err := stores.Functions.Create(function)
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("Function already exist: %s", function.FunctionName))
		}
		return nil, err
	}

	if len(in.Tags) > 0 {
		if err := stores.Functions.TagStore.Tag(in.FunctionName, in.Tags); err != nil {
			return nil, err
		}
	}

	if in.Publish {
		if _, err := stores.Functions.PublishVersion(created, ""); err != nil {
			return nil, err
		}
	}

	return created, nil
}

// deleteFunctionCore is the single entry point for function deletion. It
// handles both full-function and version-specific deletion, including
// event-source-mapping existence checks and container cleanup.
func (s *LambdaService) deleteFunctionCore(ctx context.Context, stores *lambdaStore, in *DeleteFunctionInput) error {
	functionName := in.FunctionName
	qualifier := in.Qualifier

	if qualifier == "$LATEST" {
		return NewInvalidParameter("Qualifier", "Cannot delete $LATEST version of a function")
	}

	function, err := stores.Functions.Get(functionName)
	if err != nil {
		return ErrResourceNotFound
	}

	if qualifier != "" {
		for _, v := range function.Versions {
			if v.Version == qualifier && v.ContainerID != "" {
				if rmErr := s.dockerClient.RemoveContainer(ctx, v.ContainerID, true); rmErr != nil {
					logs.Warn("Failed to remove container for version",
						logs.String("containerID", v.ContainerID),
						logs.String("function", function.FunctionName),
						logs.String("version", qualifier), logs.Err(rmErr))
				}
			}
		}
		return stores.Functions.DeleteVersion(function.FunctionName, qualifier)
	}

	mappings, err := stores.EventSources.ListByFunction(function.FunctionArn)
	if err == nil && len(mappings) > 0 {
		return ErrResourceInUse
	}

	if function.ContainerID != "" {
		if err := s.dockerClient.RemoveContainer(ctx, function.ContainerID, true); err != nil {
			logs.Warn("Failed to remove container for function",
				logs.String("containerID", function.ContainerID),
				logs.String("function", function.FunctionName), logs.Err(err))
		}
	}

	for _, v := range function.Versions {
		if v.ContainerID != "" && v.ContainerID != function.ContainerID {
			if rmErr := s.dockerClient.RemoveContainer(ctx, v.ContainerID, true); rmErr != nil {
				logs.Warn("Failed to remove version container",
					logs.String("containerID", v.ContainerID),
					logs.String("function", function.FunctionName),
					logs.String("version", v.Version), logs.Err(rmErr))
			}
		}
	}

	return stores.Functions.Delete(function.FunctionName)
}

// listFunctionsCore returns a paginated list of functions. The caller
// converts the result to the appropriate response format (HTTP JSON or
// proto).
func (s *LambdaService) listFunctionsCore(stores *lambdaStore, in *ListFunctionsInput) ([]*lambdastore.Function, string, error) {
	maxItems := in.MaxItems
	if maxItems <= 0 {
		maxItems = 50
	}

	opts := storecommon.ListOptions{
		Marker:   in.Marker,
		MaxItems: maxItems,
	}
	result, err := stores.Functions.List(opts)
	if err != nil {
		return nil, "", err
	}

	return result.Items, result.NextMarker, nil
}

// getOrCreateLambdaStore returns the full lambdaStore for the given region,
// creating it if necessary. Used by the admin handler and core functions.
func (s *LambdaService) getOrCreateLambdaStore(region string) *lambdaStore {
	if cached, ok := s.storeCache.Load(region); ok {
		if typed, ok := cached.(*lambdaStore); ok {
			return typed
		}
	}
	storage := s.getRegionalStorage(region)
	newStore := &lambdaStore{
		Functions:    lambdastore.NewFunctionStore(storage, s.accountID, region),
		Layers:       lambdastore.NewLayerStore(storage, s.accountID, region),
		EventSources: lambdastore.NewEventSourceStore(storage, s.accountID, region),
	}
	if actual, loaded := s.storeCache.LoadOrStore(region, newStore); loaded {
		if typed, ok := actual.(*lambdaStore); ok {
			return typed
		}
	}
	return newStore
}
