package lambda

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"vorpalstacks/internal/common/iam"
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

	// Region of the request; used to persist per-version code snapshots
	// when Publish is requested.
	Region string

	// IAMValidator, when injected, checks that the execution role's trust
	// policy allows the Lambda service principal. Both planes inject it
	// (the HTTP API from the request context, the admin console from the
	// service's role provider).
	IAMValidator *iam.IAMValidator

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
	Region       string
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
// tags, and optionally publishes a version. The returned Version is
// non-nil only when in.Publish requested an initial publish; callers use
// it to answer with the published version's configuration.
func (s *LambdaService) createFunctionCore(ctx context.Context, stores *lambdaStore, in *CreateFunctionInput) (*lambdastore.Function, *lambdastore.Version, error) {
	if in.FunctionName == "" {
		return nil, nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	if err := validateFunctionName(in.FunctionName); err != nil {
		return nil, nil, err
	}

	if in.Runtime == "" && in.PackageType != "Image" {
		return nil, nil, NewInvalidParameter("Runtime", "Runtime is required for Zip package type")
	}
	if in.Runtime != "" {
		canonical, ok := lambdastore.CanonicalRuntime(in.Runtime)
		if !ok {
			return nil, nil, NewInvalidParameter("Runtime", "Runtime '"+in.Runtime+"' is not supported")
		}
		in.Runtime = string(canonical)
	}

	if in.Handler == "" && in.PackageType != "Image" {
		return nil, nil, NewInvalidParameter("Handler", "Handler is required for Zip package type")
	}
	if in.Handler != "" && in.Runtime != "" {
		if err := ValidateHandler(in.Runtime, in.Handler); err != nil {
			return nil, nil, err
		}
	}

	if in.Role == "" {
		return nil, nil, NewInvalidParameter("Role", "Role ARN is required")
	}

	// The execution role must be assumable by the Lambda service principal.
	// TEST_MODE skips the trust-policy lookup so the regression suite can
	// create functions without seeded IAM roles.
	if in.Role != "" && in.IAMValidator != nil && os.Getenv("TEST_MODE") != "true" {
		if err := in.IAMValidator.ValidateRoleForService(ctx, in.Role, iam.ServicePrincipalLambda); err != nil {
			return nil, nil, err
		}
	}

	if err := validateCodeSigningConfigArn(in.CodeSigningConfigArn); err != nil {
		return nil, nil, err
	}

	if err := validatePackageType(in.PackageType); err != nil {
		return nil, nil, err
	}

	if err := validateKMSKeyArn(in.KMSKeyArn); err != nil {
		return nil, nil, err
	}

	if in.EphemeralStorage != nil {
		if err := validateEphemeralStorageSize(in.EphemeralStorage.Size); err != nil {
			return nil, nil, err
		}
	}
	if err := validateFileSystemConfigs(in.FileSystemConfigs); err != nil {
		return nil, nil, err
	}

	if in.SnapStart != nil {
		if err := validateSnapStartApplyOn(in.SnapStart.ApplyOn); err != nil {
			return nil, nil, err
		}
		if err := validateSnapStartForRuntime(in.Runtime, in.SnapStart); err != nil {
			return nil, nil, err
		}
	}
	if err := validateEnvironmentVariables(in.Environment); err != nil {
		return nil, nil, err
	}
	if err := validateLoggingConfig(in.LoggingConfig); err != nil {
		return nil, nil, err
	}
	if err := validateImageConfig(in.ImageConfig); err != nil {
		return nil, nil, err
	}

	for _, arch := range in.Architectures {
		if err := validateArchitecture(arch); err != nil {
			return nil, nil, err
		}
	}

	timeout := in.Timeout
	if timeout == 0 {
		timeout = lambdastore.DefaultFunctionTimeoutSeconds
	}
	if err := validateTimeout(timeout); err != nil {
		return nil, nil, err
	}

	memorySize := in.MemorySize
	if memorySize == 0 {
		memorySize = lambdastore.DefaultFunctionMemorySizeMB
	}
	if err := validateMemorySize(memorySize); err != nil {
		return nil, nil, err
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

	// Reject input-ground tag failures before the record is written so
	// the tag step cannot fail after Create.
	if len(in.Tags) > 0 {
		if err := stores.Functions.TagStore.ValidateTags(in.Tags); err != nil {
			return nil, nil, NewInvalidParameter("Tags", err.Error())
		}
	}

	created, err := stores.Functions.Create(function)
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionAlreadyExists) {
			return nil, nil, NewResourceConflict(fmt.Sprintf("Function already exist: %s", function.FunctionName))
		}
		return nil, nil, err
	}

	if len(in.Tags) > 0 {
		if err := stores.Functions.TagStore.Tag(in.FunctionName, in.Tags); err != nil {
			// A store-ground failure after the record exists rolls the
			// creation back — a failed CreateFunction must leave nothing
			// behind.
			s.rollbackCreatedRecord("function "+function.FunctionName, func() error {
				return s.deleteFunctionAndCode(stores, function.FunctionName, in.Region)
			})
			return nil, nil, err
		}
	}

	var published *lambdastore.Version
	if in.Publish {
		published, err = s.publishVersionWithCode(stores, created, "", "", in.Region)
		if err != nil {
			s.rollbackCreatedRecord("function "+function.FunctionName, func() error {
				return s.deleteFunctionAndCode(stores, function.FunctionName, in.Region)
			})
			return nil, nil, err
		}
	}

	return created, published, nil
}

// deleteFunctionAndCode is the compensating delete for a failed create:
// the record goes first, then the seeded code archives.
func (s *LambdaService) deleteFunctionAndCode(stores *lambdaStore, functionName, region string) error {
	if err := stores.Functions.Delete(functionName); err != nil {
		return err
	}
	s.removeFunctionCodeDir(functionName, "", region)
	return nil
}

// rollbackCreatedRecord removes a record whose creating request failed
// after the store write, so a failed Create/Publish leaves no orphan.
// The secondary error is logged rather than returned: the caller reports
// the primary failure that triggered the rollback.
func (s *LambdaService) rollbackCreatedRecord(what string, del func() error) {
	if err := del(); err != nil {
		logs.Warn("Failed to roll back partially created resource",
			logs.String("resource", what), logs.Err(err))
	}
}

// deleteFunctionCore is the single entry point for function deletion. It
// handles both full-function and version-specific deletion, including
// event-source-mapping existence checks and container cleanup.
func (s *LambdaService) deleteFunctionCore(ctx context.Context, stores *lambdaStore, in *DeleteFunctionInput) error {
	if in.FunctionName == "" {
		return NewInvalidParameter("FunctionName", "Function name is required")
	}
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
		// The authoritative delete runs first: it is what rejects an
		// alias-referenced version, so a version that survives the
		// request keeps its container and sandboxes untouched.
		if err := stores.Functions.DeleteVersion(function.FunctionName, qualifier); err != nil {
			return mapStoreError(err)
		}
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
		s.sandboxes.drainVersion(function.FunctionArn, qualifier)
		// The authoritative store delete succeeded; the version's
		// on-disk archive goes with it.
		s.removeFunctionCodeDir(function.FunctionName, qualifier, in.Region)
		return nil
	}

	mappings, err := stores.EventSources.ListByFunction(function.FunctionArn)
	if err != nil {
		return err
	}
	if len(mappings) > 0 {
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

	s.sandboxes.drainFunction(function.FunctionArn)
	if err := stores.Functions.Delete(function.FunctionName); err != nil {
		return mapStoreError(err)
	}
	// The authoritative store delete succeeded; the function's whole
	// on-disk code root ($LATEST plus every published version) goes with
	// it.
	s.removeFunctionCodeDir(function.FunctionName, "", in.Region)
	return nil
}

// listFunctionsCore returns a paginated list of functions. The caller
// converts the result to the appropriate response format (HTTP JSON or
// proto).
func (s *LambdaService) listFunctionsCore(stores *lambdaStore, in *ListFunctionsInput) ([]*lambdastore.Function, string, error) {
	maxItems := validateMaxItems(in.MaxItems)

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

// functionCodeMetadata carries the code metadata a create or update
// request resolved from its wire code members.
type functionCodeMetadata struct {
	CodeLocation string
	CodeSize     int64
	CodeSha256   string
}

// prepareCreateFunctionCodeCore resolves the wire Code map of a
// CreateFunction request into persisted code metadata: an ImageUri member
// switches the package type to Image; otherwise the archive comes from the
// ZipFile member or an S3 bucket reference (resolveCodeContent owns the
// decode order and rejects a map naming neither). Decoded archives are
// persisted under the function's $LATEST code directory with their hash
// recorded.
func (s *LambdaService) prepareCreateFunctionCodeCore(ctx context.Context, region, functionName string, codeMap map[string]interface{}, packageType string) (*functionCodeMetadata, string, string, error) {
	if codeMap == nil {
		return nil, "", "", NewInvalidParameter("Code", "Code is required")
	}

	meta := &functionCodeMetadata{}
	var imageUri string

	if uri, ok := codeMap["ImageUri"].(string); ok && uri != "" {
		imageUri = uri
		packageType = "Image"
	}

	if imageUri == "" {
		zipFile, err := s.resolveCodeContent(ctx, region, "Code", codeMap)
		if err != nil {
			return nil, "", "", err
		}
		codeLocation, codeSize, err := s.storeCode(functionName, "$LATEST", zipFile, region)
		if err != nil {
			return nil, "", "", err
		}
		meta.CodeLocation, meta.CodeSize = codeLocation, codeSize
		meta.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
	}

	return meta, imageUri, packageType, nil
}

// getAccountSettingsCore computes the account limits and usage summary
// over every function in the region.
func (s *LambdaService) getAccountSettingsCore(stores *lambdaStore) (map[string]interface{}, error) {
	result, err := stores.Functions.ListAllFunctions()
	if err != nil {
		return nil, fmt.Errorf("failed to list functions: %w", err)
	}

	var totalCodeSize int64
	var reservedSum int64
	for _, fn := range result {
		totalCodeSize += fn.CodeSize
		if fn.ReservedConcurrency != nil {
			reservedSum += *fn.ReservedConcurrency
		}
	}
	// The unreserved concurrency is the regional limit minus the reserved
	// amounts of every function in the region.
	unreserved := lambdastore.AccountLimitConcurrentExecutions - reservedSum
	if unreserved < 0 {
		unreserved = 0
	}

	return map[string]interface{}{
		"AccountLimit": map[string]interface{}{
			"TotalCodeSize":                  lambdastore.AccountLimitTotalCodeSize,
			"CodeSizeUnzipped":               lambdastore.AccountLimitCodeSizeUnzipped,
			"CodeSizeZipped":                 lambdastore.AccountLimitCodeSizeZipped,
			"ConcurrentExecutions":           lambdastore.AccountLimitConcurrentExecutions,
			"UnreservedConcurrentExecutions": unreserved,
		},
		"AccountUsage": map[string]interface{}{
			"TotalCodeSize": totalCodeSize,
			"FunctionCount": len(result),
		},
	}, nil
}

// publishVersionWithCode publishes a new version of the function and
// persists the version's code snapshot under its own version directory, so
// the published version stays executable after the $LATEST code changes or
// all containers are recycled. Container image packages carry no zip
// archive and skip the code persistence step. This mirrors how layer
// versions persist their content at publish time. A non-empty revisionId is
// the modelled optimistic-locking precondition: publishing fails with
// PreconditionFailedException unless it matches the function's current
// revision.
func (s *LambdaService) publishVersionWithCode(stores *lambdaStore, function *lambdastore.Function, description, revisionId, region string) (*lambdastore.Version, error) {
	if revisionId != "" && function.RevisionId != revisionId {
		return nil, NewPreconditionFailed(revisionMismatchMessage)
	}

	var latestCode []byte
	if function.PackageType != "Image" && function.ImageUri == "" {
		var err error
		latestCode, err = s.loadCode(function.FunctionName, "$LATEST", region)
		if err != nil {
			return nil, NewLambdaError("ServiceException",
				fmt.Sprintf("The $LATEST code of function %s is not available for publishing.", function.FunctionName),
				http.StatusInternalServerError)
		}
	}

	version, err := stores.Functions.PublishVersionAtomically(function.FunctionName, description, revisionId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if latestCode != nil {
		if _, _, err := s.storeCode(function.FunctionName, version.Version, latestCode, region); err != nil {
			// The version record exists but its archive does not — the
			// version would be permanently uninvocable, so the failed
			// publish rolls the version back.
			s.rollbackCreatedRecord(fmt.Sprintf("version %s of %s", version.Version, function.FunctionName), func() error {
				return stores.Functions.DeleteVersion(function.FunctionName, version.Version)
			})
			return nil, NewLambdaError("ServiceException",
				fmt.Sprintf("Failed to persist the code of version %s: %v", version.Version, err),
				http.StatusInternalServerError)
		}
	}

	return version, nil
}
