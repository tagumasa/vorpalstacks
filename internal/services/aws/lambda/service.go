// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/client/mobyclient"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// lambdaStore holds the stores for Lambda resources.
type lambdaStore struct {
	Functions    *lambdastore.FunctionStore
	Layers       *lambdastore.LayerStore
	EventSources *lambdastore.EventSourceStore
}

// LambdaService provides Lambda operations.
type LambdaService struct {
	storage      storage.BasicStorage
	s3Objects    map[string]s3store.ObjectStoreInterface
	dockerClient *mobyclient.Client
	logsStores   sync.Map // region → *logsstore.Store
	storeCache   sync.Map // region → *lambdaStore
	accountID    string
	region       string
	dataDir      string
	dataDirOnce  sync.Once
	asyncWg      sync.WaitGroup // goroutine tracking for InvokeAsync
	s3ObjectsMu  sync.RWMutex
}

func (s *LambdaService) store(reqCtx *request.RequestContext) (*lambdaStore, error) {
	region := reqCtx.GetRegion()
	if cached, ok := s.storeCache.Load(region); ok {
		if typed, ok := cached.(*lambdaStore); ok {
			return typed, nil
		}
	}

	if injectedStore := reqCtx.GetLambdaStore(); injectedStore != nil {
		if ls, ok := injectedStore.(*lambdastore.LambdaStore); ok {
			newStore := &lambdaStore{
				Functions:    ls.FunctionsRaw(),
				Layers:       ls.LayersRaw(),
				EventSources: ls.EventSourcesRaw(),
			}
			if actual, loaded := s.storeCache.LoadOrStore(region, newStore); loaded {
				if typed, ok := actual.(*lambdaStore); ok {
					return typed, nil
				}
			}
			return newStore, nil
		}
	}

	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}

	newStore := &lambdaStore{
		Functions:    lambdastore.NewFunctionStore(storage, s.accountID, region),
		Layers:       lambdastore.NewLayerStore(storage, s.accountID, region),
		EventSources: lambdastore.NewEventSourceStore(storage, s.accountID, region),
	}

	if actual, loaded := s.storeCache.LoadOrStore(region, newStore); loaded {
		if typed, ok := actual.(*lambdaStore); ok {
			return typed, nil
		}
	}
	return newStore, nil
}

// NewLambdaService creates a new Lambda service instance.
func NewLambdaService(store storage.BasicStorage, dockerClient *mobyclient.Client, accountID, region, dataDir string) *LambdaService {
	return &LambdaService{
		storage:      store,
		s3Objects:    make(map[string]s3store.ObjectStoreInterface),
		dockerClient: dockerClient,
		accountID:    accountID,
		region:       region,
		dataDir:      dataDir,
	}
}

// NewLambdaServiceWithLogs creates a new Lambda service instance with a pre-configured CloudWatch Logs store.
func NewLambdaServiceWithLogs(store storage.BasicStorage, dockerClient *mobyclient.Client, logsStore *logsstore.Store, accountID, region, dataDir string) *LambdaService {
	s := &LambdaService{
		storage:      store,
		s3Objects:    make(map[string]s3store.ObjectStoreInterface),
		dockerClient: dockerClient,
		accountID:    accountID,
		region:       region,
		dataDir:      dataDir,
	}
	if logsStore != nil {
		s.logsStores.Store(region, logsStore)
	}
	return s
}

// SetS3ObjectStore registers an S3 object store for the specified region.
func (s *LambdaService) SetS3ObjectStore(region string, store s3store.ObjectStoreInterface) {
	s.s3ObjectsMu.Lock()
	defer s.s3ObjectsMu.Unlock()
	if s.s3Objects == nil {
		s.s3Objects = make(map[string]s3store.ObjectStoreInterface)
	}
	s.s3Objects[region] = store
}

func (s *LambdaService) getS3ObjectStore(region string) s3store.ObjectStoreInterface {
	s.s3ObjectsMu.RLock()
	defer s.s3ObjectsMu.RUnlock()
	return s.s3Objects[region]
}

func (s *LambdaService) getLogsStore(region string) *logsstore.Store {
	if cached, ok := s.logsStores.Load(region); ok {
		if typed, ok := cached.(*logsstore.Store); ok {
			return typed
		}
	}
	store := logsstore.NewStore(s.storage.Bucket("logs-"+region), s.accountID, region, s.dataDir)
	if actual, loaded := s.logsStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(*logsstore.Store); ok {
			return typed
		}
	}
	return store
}

func (s *LambdaService) getLogsStoreFromContext(reqCtx *request.RequestContext) *logsstore.Store {
	if store := reqCtx.GetCloudWatchLogsStore(); store != nil {
		if concrete, ok := store.(*logsstore.Store); ok {
			return concrete
		}
	}
	return s.getLogsStore(reqCtx.GetRegion())
}

// RegisterHandlers registers the Lambda service handlers with the dispatcher.
func (s *LambdaService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("lambda", "CreateFunction", s.CreateFunction)
	d.RegisterHandlerForService("lambda", "DeleteFunction", s.DeleteFunction)
	d.RegisterHandlerForService("lambda", "GetFunction", s.GetFunction)
	d.RegisterHandlerForService("lambda", "GetFunctionConfiguration", s.GetFunctionConfiguration)
	d.RegisterHandlerForService("lambda", "ListFunctions", s.ListFunctions)
	d.RegisterHandlerForService("lambda", "UpdateFunctionCode", s.UpdateFunctionCode)
	d.RegisterHandlerForService("lambda", "UpdateFunctionConfiguration", s.UpdateFunctionConfiguration)

	d.RegisterHandlerForService("lambda", "Invoke", s.Invoke)
	d.RegisterHandlerForService("lambda", "InvokeAsync", s.InvokeAsync)
	d.RegisterHandlerForService("lambda", "InvokeWithResponseStream", s.InvokeWithResponseStream)

	d.RegisterHandlerForService("lambda", "PublishVersion", s.PublishVersion)
	d.RegisterHandlerForService("lambda", "ListVersionsByFunction", s.ListVersionsByFunction)

	d.RegisterHandlerForService("lambda", "CreateAlias", s.CreateAlias)
	d.RegisterHandlerForService("lambda", "DeleteAlias", s.DeleteAlias)
	d.RegisterHandlerForService("lambda", "GetAlias", s.GetAlias)
	d.RegisterHandlerForService("lambda", "UpdateAlias", s.UpdateAlias)
	d.RegisterHandlerForService("lambda", "ListAliases", s.ListAliases)

	d.RegisterHandlerForService("lambda", "PublishLayerVersion", s.PublishLayerVersion)
	d.RegisterHandlerForService("lambda", "DeleteLayerVersion", s.DeleteLayerVersion)
	d.RegisterHandlerForService("lambda", "GetLayerVersion", s.GetLayerVersion)
	d.RegisterHandlerForService("lambda", "ListLayers", s.ListLayers)
	d.RegisterHandlerForService("lambda", "ListLayerVersions", s.ListLayerVersions)

	d.RegisterHandlerForService("lambda", "CreateEventSourceMapping", s.CreateEventSourceMapping)
	d.RegisterHandlerForService("lambda", "DeleteEventSourceMapping", s.DeleteEventSourceMapping)
	d.RegisterHandlerForService("lambda", "GetEventSourceMapping", s.GetEventSourceMapping)
	d.RegisterHandlerForService("lambda", "UpdateEventSourceMapping", s.UpdateEventSourceMapping)
	d.RegisterHandlerForService("lambda", "ListEventSourceMappings", s.ListEventSourceMappings)

	d.RegisterHandlerForService("lambda", "AddPermission", s.AddPermission)
	d.RegisterHandlerForService("lambda", "RemovePermission", s.RemovePermission)
	d.RegisterHandlerForService("lambda", "GetPolicy", s.GetPolicy)

	d.RegisterHandlerForService("lambda", "TagResource", s.TagResource)
	d.RegisterHandlerForService("lambda", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("lambda", "ListTags", s.ListTags)

	d.RegisterHandlerForService("lambda", "PutFunctionConcurrency", s.PutFunctionConcurrency)
	d.RegisterHandlerForService("lambda", "GetFunctionConcurrency", s.GetFunctionConcurrency)
	d.RegisterHandlerForService("lambda", "DeleteFunctionConcurrency", s.DeleteFunctionConcurrency)

	d.RegisterHandlerForService("lambda", "PutProvisionedConcurrencyConcurrency", s.PutProvisionedConcurrencyConcurrency)
	d.RegisterHandlerForService("lambda", "GetProvisionedConcurrencyConfig", s.GetProvisionedConcurrencyConfig)
	d.RegisterHandlerForService("lambda", "DeleteProvisionedConcurrencyConcurrency", s.DeleteProvisionedConcurrencyConcurrency)
	d.RegisterHandlerForService("lambda", "ListProvisionedConcurrencyConfigs", s.ListProvisionedConcurrencyConfigs)

	d.RegisterHandlerForService("lambda", "PutFunctionEventInvokeConfig", s.PutFunctionEventInvokeConfig)
	d.RegisterHandlerForService("lambda", "GetFunctionEventInvokeConfig", s.GetFunctionEventInvokeConfig)
	d.RegisterHandlerForService("lambda", "DeleteFunctionEventInvokeConfig", s.DeleteFunctionEventInvokeConfig)
	d.RegisterHandlerForService("lambda", "ListFunctionEventInvokeConfigs", s.ListFunctionEventInvokeConfigs)

	d.RegisterHandlerForService("lambda", "CreateFunctionUrlConfig", s.CreateFunctionUrlConfig)
	d.RegisterHandlerForService("lambda", "DeleteFunctionUrlConfig", s.DeleteFunctionUrlConfig)
	d.RegisterHandlerForService("lambda", "GetFunctionUrlConfig", s.GetFunctionUrlConfig)
	d.RegisterHandlerForService("lambda", "UpdateFunctionUrlConfig", s.UpdateFunctionUrlConfig)
	d.RegisterHandlerForService("lambda", "ListFunctionUrlConfigs", s.ListFunctionUrlConfigs)

	d.RegisterHandlerForService("lambda", "GetAccountSettings", s.GetAccountSettings)
}

func (s *LambdaService) resolveQualifier(store *lambdastore.FunctionStore, functionName, qualifier string) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, error) {
	function, version, alias, err := store.ResolveQualifier(functionName, qualifier)
	if err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, nil, nil, ErrResourceNotFound
		}
		if err == lambdastore.ErrVersionNotFound {
			return nil, nil, nil, NewLambdaError("ResourceNotFoundException",
				fmt.Sprintf("Qualifier '%s' not found for function '%s'.", qualifier, functionName),
				http.StatusNotFound)
		}
		return nil, nil, nil, err
	}
	return function, version, alias, nil
}

func sanitizePathComponent(name string) string {
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, name)
	if safe == "" {
		safe = "unnamed"
	}
	return safe
}

func (s *LambdaService) storeCode(functionName, version string, code []byte, region string) (string, int64, error) {
	s.dataDirOnce.Do(func() {
		if s.dataDir == "" {
			s.dataDir = "./data"
		}
	})
	dataDir := s.dataDir

	if version == "" {
		version = "$LATEST"
	}

	safeFunctionName := sanitizePathComponent(functionName)
	safeVersion := sanitizePathComponent(version)
	codeDir := fmt.Sprintf("%s/%s/code/%s/%s", dataDir, sanitizePathComponent(region), safeFunctionName, safeVersion)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create code directory: %w", err)
	}

	codePath := fmt.Sprintf("%s/code.zip", codeDir)
	if err := os.WriteFile(codePath, code, 0644); err != nil {
		return "", 0, fmt.Errorf("failed to write code file: %w", err)
	}

	return codePath, int64(len(code)), nil
}

func (s *LambdaService) fetchCodeFromS3(ctx context.Context, bucket, key, region string) ([]byte, error) {
	s3ObjStore := s.getS3ObjectStore(region)
	if s3ObjStore == nil {
		return nil, fmt.Errorf("S3 object store not configured for region %s", region)
	}
	reader, _, err := s3ObjStore.Get(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: s3://%s/%s: %w", bucket, key, err)
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (s *LambdaService) loadCode(functionName, version string, region string) ([]byte, error) {
	s.dataDirOnce.Do(func() {
		if s.dataDir == "" {
			s.dataDir = "./data"
		}
	})
	dataDir := s.dataDir

	if version == "" {
		version = "$LATEST"
	}

	codePath := fmt.Sprintf("%s/%s/code/%s/%s/code.zip", dataDir, sanitizePathComponent(region), sanitizePathComponent(functionName), sanitizePathComponent(version))
	code, err := os.ReadFile(codePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read code file: %w", err)
	}
	return code, nil
}

func (s *LambdaService) getRuntimeImage(runtime lambdastore.Runtime) string {
	return lambdastore.GetImageForRuntime(runtime)
}

func (s *LambdaService) ensureFunctionContainer(function *lambdastore.Function, ver *lambdastore.Version, store *lambdastore.FunctionStore, region string) (string, error) {
	ctx := context.Background()

	version := "$LATEST"
	if ver != nil {
		version = ver.Version
	}

	containerName := fmt.Sprintf("lambda-%s-%s", function.FunctionName, sanitizeForContainerName(version))

	containerID := function.ContainerID
	if ver != nil && ver.ContainerID != "" {
		containerID = ver.ContainerID
	}

	if containerID != "" {
		status, err := s.dockerClient.GetContainerStatus(ctx, containerID)
		if err == nil && status == mobyclient.ContainerStatusRunning {
			return containerID, nil
		}
	}

	image := s.getRuntimeImage(function.Runtime)
	if function.ImageUri != "" {
		image = function.ImageUri
	}

	envVars := map[string]string{
		"AWS_LAMBDA_FUNCTION_TIMEOUT":     fmt.Sprintf("%d", function.Timeout),
		"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": fmt.Sprintf("%d", function.MemorySize),
		"AWS_LAMBDA_FUNCTION_HANDLER":     function.Handler,
		"AWS_LAMBDA_FUNCTION_NAME":        function.FunctionName,
		"AWS_LAMBDA_FUNCTION_VERSION":     version,
		"AWS_REGION":                      region,
	}

	if function.Environment != nil {
		for k, v := range function.Environment.Variables {
			envVars[k] = v
		}
	}

	cfg := mobyclient.AdvancedContainerConfig{
		Name:       containerName,
		Image:      image,
		PullImage:  true,
		Env:        envVars,
		Entrypoint: []string{"/lambda-entrypoint.sh"},
		Cmd:        []string{function.Handler},
		Network:    "bridge",
		AutoRemove: false,
	}

	result, err := s.dockerClient.CreateContainerFromConfig(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := s.dockerClient.StartContainer(ctx, result.ID); err != nil {
		if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
			logs.Warn("Failed to remove container after start failure", logs.String("containerID", result.ID), logs.Err(rmErr))
		}
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	if version == "$LATEST" {
		originalContainerID := function.ContainerID
		function.ContainerID = result.ID
		function.ContainerImageID = result.ID
		if err := store.Update(function); err != nil {
			function.ContainerID = originalContainerID
			function.ContainerImageID = originalContainerID
			if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
				logs.Warn("Failed to remove container during rollback", logs.String("containerID", result.ID), logs.Err(rmErr))
			}
			return "", fmt.Errorf("failed to update function: %w", err)
		}
	} else if ver != nil {
		ver.ContainerID = result.ID
		ver.ContainerImageID = result.ID
		if err := store.Update(function); err != nil {
			ver.ContainerID = ""
			ver.ContainerImageID = ""
			if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
				logs.Warn("Failed to remove container during rollback", logs.String("containerID", result.ID), logs.Err(rmErr))
			}
			return "", fmt.Errorf("failed to update function version: %w", err)
		}
	}

	return result.ID, nil
}

func (s *LambdaService) copyCodeToContainer(containerID string, code []byte) error {
	ctx := context.Background()
	return s.dockerClient.CreateFileInContainer(ctx, containerID, "/var/task/code.zip", code)
}

func (s *LambdaService) invokeFunction(function *lambdastore.Function, ver *lambdastore.Version, store *lambdastore.FunctionStore, region string, payload []byte) (*lambdastore.InvocationResult, error) {
	ctx := context.Background()

	containerID, err := s.ensureFunctionContainer(function, ver, store, region)
	if err != nil {
		return nil, err
	}

	version := "$LATEST"
	if ver != nil {
		version = ver.Version
	}

	code, err := s.loadCode(function.FunctionName, version, region)
	if err == nil && len(code) > 0 {
		if err := s.copyCodeToContainer(containerID, code); err != nil {
			return nil, fmt.Errorf("failed to copy code to container: %w", err)
		}
	}

	execResult, err := s.dockerClient.ExecWithStdin(ctx, containerID, mobyclient.ExecConfig{
		Cmd:          []string{"/var/runtime/bootstrap"},
		AttachStdout: true,
		AttachStderr: true,
	}, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to exec in container: %w", err)
	}

	s.writeLambdaLogs(function.FunctionName, version, execResult.Stdout, execResult.Stderr, region)

	return &lambdastore.InvocationResult{
		StatusCode:      http.StatusOK,
		ExecutedVersion: version,
		Payload:         []byte(execResult.Stdout),
	}, nil
}

func (s *LambdaService) writeLambdaLogs(functionName, version, stdout, stderr, region string) {
	logsStore := s.getLogsStore(region)

	logGroupName := "/aws/lambda/" + functionName
	now := time.Now().UTC()
	requestID := uuid.New().String()
	streamName := fmt.Sprintf("%d/%02d/%02d/[%s]%s",
		now.Year(), now.Month(), now.Day(), version, requestID[:8])

	_, err := logsStore.GetLogGroup(logGroupName)
	if err != nil {
		lg := logsstore.NewLogGroup(logGroupName, region, s.accountID)
		if createErr := logsStore.CreateLogGroup(lg); createErr != nil {
			return
		}
	}

	ls := logsstore.NewLogStream(streamName, logGroupName)
	if createErr := logsStore.CreateLogStream(ls); createErr != nil {
		return
	}

	ts := now.UnixNano() / int64(time.Millisecond)
	events := []logsstore.LogEntry{
		{Timestamp: ts, Message: fmt.Sprintf("START RequestId: %s", requestID)},
	}

	for _, line := range strings.Split(stdout, "\n") {
		if line != "" {
			events = append(events, logsstore.LogEntry{Timestamp: ts, Message: line})
		}
	}
	for _, line := range strings.Split(stderr, "\n") {
		if line != "" {
			events = append(events, logsstore.LogEntry{Timestamp: ts, Message: line})
		}
	}

	events = append(events, logsstore.LogEntry{
		Timestamp: ts,
		Message:   fmt.Sprintf("END RequestId: %s", requestID),
	})
	events = append(events, logsstore.LogEntry{
		Timestamp: ts,
		Message:   fmt.Sprintf("REPORT RequestId: %s\tDuration: 0.00 ms\tBilled Duration: 0 ms\tMemory Size: 128 MB\tMax Memory Used: 0 MB", requestID),
	})

	if _, err := logsStore.PutLogEvents(logGroupName, streamName, events); err != nil {
		logs.Warn("Failed to write Lambda logs", logs.String("function", functionName), logs.Err(err))
	}
}

// InvokeForGateway invokes a Lambda function by name for API Gateway integration.
// Returns the HTTP status code, response payload, and any error.
func (s *LambdaService) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	store := lambdastore.NewFunctionStore(s.storage, s.accountID, s.region)
	function, ver, _, err := s.resolveQualifier(store, functionName, "")
	if err != nil {
		return 0, nil, err
	}
	result, err := s.invokeFunction(function, ver, store, s.region, payload)
	if err != nil {
		return 0, nil, err
	}
	if result == nil {
		return 0, nil, fmt.Errorf("invocation returned nil result")
	}
	return result.StatusCode, result.Payload, nil
}

// GetFunctionStore returns a new FunctionStore for the Lambda service.
func (s *LambdaService) GetFunctionStore() *lambdastore.FunctionStore {
	return lambdastore.NewFunctionStore(s.storage, s.accountID, s.region)
}

// GetAccountSettings returns account limits and usage for Lambda functions.
func (s *LambdaService) GetAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Functions.List(common.ListOptions{MaxItems: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to list functions: %w", err)
	}
	return map[string]interface{}{
		"AccountLimit": map[string]interface{}{
			"TotalCodeSizeZipped":            80530636800,
			"CodeSizeUnzipped":               262144000,
			"ConcurrentExecutions":           1000,
			"UnreservedConcurrentExecutions": 1000,
		},
		"AccountUsage": map[string]interface{}{
			"TotalCodeSizeZipped": 0,
			"FunctionCount":       len(result.Items),
		},
	}, nil
}

// Shutdown gracefully shuts down the Lambda service by waiting for all asynchronous operations to complete.
func (s *LambdaService) Shutdown() {
	s.asyncWg.Wait()
}

func sanitizeForContainerName(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == '_' || r == '.' || r == '-' {
			b.WriteRune(r)
		} else if i > 0 {
			b.WriteRune('-')
		}
	}
	result := b.String()
	if len(result) == 0 {
		return "unknown"
	}
	if c := result[0]; !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
		result = "x-" + result
	}
	return result
}
