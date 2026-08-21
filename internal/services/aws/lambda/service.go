// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/client/mobyclient"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/naming"
)

// lambdaStore holds the stores for Lambda resources.
type lambdaStore struct {
	Functions    *lambdastore.FunctionStore
	Layers       *lambdastore.LayerStore
	EventSources *lambdastore.EventSourceStore
}

// LambdaService provides Lambda operations.
type LambdaService struct {
	storageManager *storage.RegionStorageManager
	s3Invoker      eventbus.S3Invoker
	logsInvoker    eventbus.LogsInvoker
	dockerClient   mobyclient.ContainerLifecycle
	bus            eventbus.Bus
	storeCache     sync.Map // region → *lambdaStore
	accountID      string
	region         string
	hostEndpoint   string
	dataDir        string
	dataDirOnce    sync.Once
	asyncWg        sync.WaitGroup
	esmPoller      *esmPoller
	inflight       sync.Map // functionArn → *atomic.Int32
	// containerEnsureMu serialises cold container creation per container
	// name, and containerIDs caches the last ensured container per name.
	// Concurrent invocations of the same cold function (for example
	// parallel event source batches) would otherwise race the
	// check-then-create sequence; the loser fails with a name conflict and
	// its batch is re-delivered on the next poll.
	containerEnsureMu sync.Map // container name → *sync.Mutex
	containerIDs      sync.Map // container name → container ID
}

func (s *LambdaService) store(reqCtx *request.RequestContext) (*lambdaStore, error) {
	return storecommon.GetOrCreateStoreE(&s.storeCache, reqCtx.GetRegion(), func() (*lambdaStore, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return &lambdaStore{
			Functions:    lambdastore.NewFunctionStore(storage, s.accountID, reqCtx.GetRegion()),
			Layers:       lambdastore.NewLayerStore(storage, s.accountID, reqCtx.GetRegion()),
			EventSources: lambdastore.NewEventSourceStore(storage, s.accountID, reqCtx.GetRegion()),
		}, nil
	})
}

// NewLambdaService creates a new Lambda service instance.
// Optional dependencies (logs store, S3 object store) should be injected via
// setter methods before registering handlers.
// getInflightCounter returns the atomic counter for in-flight invocations
// of the given function ARN, creating one on first use.
func (s *LambdaService) getInflightCounter(functionArn string) *atomic.Int32 {
	if v, ok := s.inflight.Load(functionArn); ok {
		return v.(*atomic.Int32)
	}
	c := &atomic.Int32{}
	actual, loaded := s.inflight.LoadOrStore(functionArn, c)
	if loaded {
		return actual.(*atomic.Int32)
	}
	return c
}

func NewLambdaService(dockerClient mobyclient.ContainerLifecycle, accountID, region, dataDir string) *LambdaService {
	svc := &LambdaService{
		dockerClient: dockerClient,
		accountID:    accountID,
		region:       region,
		dataDir:      dataDir,
	}
	// Remove orphaned Lambda containers from previous server instances that
	// were killed (SIGKILL) before Shutdown() could clean them up.
	svc.cleanupOrphanedContainers()
	return svc
}

// cleanupOrphanedContainers scans Docker for containers whose name starts
// with "lambda-" and removes them. This prevents resource accumulation when
// the server is killed without graceful shutdown (e.g. pkill -9 during test
// runs). Called once at construction time before any new containers are
// created.
func (s *LambdaService) cleanupOrphanedContainers() {
	if s.dockerClient == nil {
		return
	}
	ctx := context.Background()
	containers, err := s.dockerClient.ListContainers(ctx, true)
	if err != nil {
		logs.Warn("Failed to list containers for orphan cleanup", logs.Err(err))
		return
	}
	removed := 0
	for _, c := range containers {
		// Docker container names returned by the API include a leading "/".
		name := strings.TrimPrefix(c.Name, "/")
		if strings.HasPrefix(name, "lambda-") {
			if rmErr := s.dockerClient.RemoveContainer(ctx, c.ID, true); rmErr != nil {
				logs.Warn("Failed to remove orphaned Lambda container",
					logs.String("containerID", c.ID),
					logs.String("name", name),
					logs.Err(rmErr))
			} else {
				removed++
			}
		}
	}
	if removed > 0 {
		logs.Info("Cleaned up orphaned Lambda containers", logs.Int("count", removed))
	}
}

// SetS3Invoker injects the S3 invoker for reading deployment packages.
func (s *LambdaService) SetS3Invoker(invoker eventbus.S3Invoker) {
	s.s3Invoker = invoker
}

// SetLogsInvoker injects the Logs invoker for writing Lambda execution logs.
func (s *LambdaService) SetLogsInvoker(invoker eventbus.LogsInvoker) {
	s.logsInvoker = invoker
}

// SetStorageManager sets the region storage manager for resolving regional storage.
func (s *LambdaService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetHostEndpoint sets the endpoint URL injected into Lambda containers
// so they can reach the vorpalstacks host from inside Docker.
func (s *LambdaService) SetHostEndpoint(endpoint string) {
	s.hostEndpoint = endpoint
}

// SetEventBus injects the event bus for Lambda log delivery. When set,
// writeLambdaLogs publishes a LambdaLogWriteEvent instead of calling
// the logsStore directly, enabling metric filter and subscription filter
// evaluation on Lambda-produced logs.
func (s *LambdaService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// StartESMPoller starts the background ESM polling goroutine. It creates
// an EventSourceStore for each configured region and begins polling all
// enabled SQS event source mappings. Safe to call multiple times.
func (s *LambdaService) StartESMPoller(ctx context.Context) {
	if s.esmPoller == nil {
		s.esmPoller = newESMPoller(0, 0, nil)
	}
	s.esmPoller.esmStore = lambdastore.NewEventSourceStore(s.getRegionalStorage(s.region), s.accountID, s.region)
	s.esmPoller.lambdaSvc = s
	s.esmPoller.region = s.region
	s.esmPoller.storageManager = s.storageManager
	s.esmPoller.bus = s.bus
	s.esmPoller.Start(ctx)
}

// StopESMPoller gracefully shuts down the ESM polling loop, waiting for
// any in-flight Lambda invocations to complete.
func (s *LambdaService) StopESMPoller() {
	if s.esmPoller != nil {
		s.esmPoller.Stop()
	}
}

func (s *LambdaService) getRegionalStorage(region string) storage.BasicStorage {
	if s.storageManager != nil {
		if st, err := s.storageManager.GetStorage(region); err == nil {
			return st
		}
	}
	return nil
}

// RegisterHandlers registers the Lambda service handlers with the dispatcher.
func (s *LambdaService) RegisterHandlers(d handler.Registrar) {
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
	d.RegisterHandlerForService("lambda", "GetLayerVersionByArn", s.GetLayerVersionByArn)
	d.RegisterHandlerForService("lambda", "ListLayers", s.ListLayers)
	d.RegisterHandlerForService("lambda", "ListLayerVersions", s.ListLayerVersions)

	d.RegisterHandlerForService("lambda", "AddLayerVersionPermission", s.AddLayerVersionPermission)
	d.RegisterHandlerForService("lambda", "RemoveLayerVersionPermission", s.RemoveLayerVersionPermission)
	d.RegisterHandlerForService("lambda", "GetLayerVersionPolicy", s.GetLayerVersionPolicy)

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

	d.RegisterHandlerForService("lambda", "PutProvisionedConcurrencyConfig", s.PutProvisionedConcurrencyConfig)
	d.RegisterHandlerForService("lambda", "GetProvisionedConcurrencyConfig", s.GetProvisionedConcurrencyConfig)
	d.RegisterHandlerForService("lambda", "DeleteProvisionedConcurrencyConfig", s.DeleteProvisionedConcurrencyConfig)
	d.RegisterHandlerForService("lambda", "ListProvisionedConcurrencyConfigs", s.ListProvisionedConcurrencyConfigs)

	d.RegisterHandlerForService("lambda", "PutFunctionEventInvokeConfig", s.PutFunctionEventInvokeConfig)
	d.RegisterHandlerForService("lambda", "GetFunctionEventInvokeConfig", s.GetFunctionEventInvokeConfig)
	d.RegisterHandlerForService("lambda", "DeleteFunctionEventInvokeConfig", s.DeleteFunctionEventInvokeConfig)
	d.RegisterHandlerForService("lambda", "UpdateFunctionEventInvokeConfig", s.UpdateFunctionEventInvokeConfig)
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

func (s *LambdaService) initDataDir() string {
	s.dataDirOnce.Do(func() {
		if s.dataDir == "" {
			s.dataDir = "./data"
		}
	})
	return s.dataDir
}

// storeLayerCode persists a layer version's zip archive to disk so it
// can be retrieved later via GetLayerVersion for download by clients.
func (s *LambdaService) storeLayerCode(layerName string, versionNum int64, code []byte, region string) (string, error) {
	dataDir := s.initDataDir()

	safeLayer := naming.SanitizePathComponent(layerName)
	safeRegion := naming.SanitizePathComponent(region)
	codeDir := fmt.Sprintf("%s/%s/layers/%s/%d", dataDir, safeRegion, safeLayer, versionNum)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create layer code directory: %w", err)
	}

	codePath := fmt.Sprintf("%s/code.zip", codeDir)
	if err := os.WriteFile(codePath, code, 0644); err != nil {
		return "", fmt.Errorf("failed to write layer code file: %w", err)
	}

	return codePath, nil
}

func (s *LambdaService) storeCode(functionName, version string, code []byte, region string) (string, int64, error) {
	dataDir := s.initDataDir()

	if version == "" {
		version = "$LATEST"
	}

	safeFunctionName := naming.SanitizePathComponent(functionName)
	safeVersion := naming.SanitizePathComponent(version)
	codeDir := fmt.Sprintf("%s/%s/code/%s/%s", dataDir, naming.SanitizePathComponent(region), safeFunctionName, safeVersion)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create code directory: %w", err)
	}

	codePath := fmt.Sprintf("%s/code.zip", codeDir)
	if err := os.WriteFile(codePath, code, 0644); err != nil {
		return "", 0, fmt.Errorf("failed to write code file: %w", err)
	}

	return codePath, int64(len(code)), nil
}

// fetchCodeFromS3 reads a deployment package from S3. A non-empty
// versionID reads that specific object version ("For versioned objects,
// the version of the deployment package object to use"); an empty one
// reads the latest version.
func (s *LambdaService) fetchCodeFromS3(ctx context.Context, bucket, key, versionID, region string) ([]byte, error) {
	if s.s3Invoker == nil {
		return nil, fmt.Errorf("S3 invoker not configured")
	}
	data, err := s.s3Invoker.GetObjectVersion(ctx, region, bucket, key, versionID, 250*1024*1024)
	if err != nil {
		if versionID != "" {
			return nil, fmt.Errorf("failed to get object version from S3: s3://%s/%s@%s: %w", bucket, key, versionID, err)
		}
		return nil, fmt.Errorf("failed to get object from S3: s3://%s/%s: %w", bucket, key, err)
	}
	return data, nil
}

func (s *LambdaService) loadCode(functionName, version string, region string) ([]byte, error) {
	dataDir := s.initDataDir()

	if version == "" {
		version = "$LATEST"
	}

	codePath := fmt.Sprintf("%s/%s/code/%s/%s/code.zip", dataDir, naming.SanitizePathComponent(region), naming.SanitizePathComponent(functionName), naming.SanitizePathComponent(version))
	code, err := os.ReadFile(codePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read code file: %w", err)
	}
	return code, nil
}

// publishVersionWithCode publishes a new version of the function and
// persists the version's code snapshot under its own version directory, so
// the published version stays executable after the $LATEST code changes or
// all containers are recycled. Container image packages carry no zip
// archive and skip the code persistence step. This mirrors how layer
// versions persist their content at publish time.
func (s *LambdaService) publishVersionWithCode(stores *lambdaStore, function *lambdastore.Function, description, region string) (*lambdastore.Version, error) {
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

	version, err := stores.Functions.PublishVersion(function, description)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if latestCode != nil {
		if _, _, err := s.storeCode(function.FunctionName, version.Version, latestCode, region); err != nil {
			return nil, NewLambdaError("ServiceException",
				fmt.Sprintf("Failed to persist the code of version %s: %v", version.Version, err),
				http.StatusInternalServerError)
		}
	}

	return version, nil
}

func (s *LambdaService) getRuntimeImage(runtime lambdastore.Runtime) string {
	return lambdastore.GetImageForRuntime(runtime)
}

// executionConfig carries the runtime parameters an execution must use.
// Published versions execute their own immutable snapshot; only $LATEST
// follows the live function configuration.
type executionConfig struct {
	Runtime     lambdastore.Runtime
	Handler     string
	Timeout     int32
	MemorySize  int32
	ImageUri    string
	Environment *lambdastore.Environment
}

func executionConfigFor(function *lambdastore.Function, ver *lambdastore.Version) executionConfig {
	if ver != nil {
		return executionConfig{
			Runtime:     ver.Runtime,
			Handler:     ver.Handler,
			Timeout:     ver.Timeout,
			MemorySize:  ver.MemorySize,
			ImageUri:    ver.ImageUri,
			Environment: ver.Environment,
		}
	}
	return executionConfig{
		Runtime:     function.Runtime,
		Handler:     function.Handler,
		Timeout:     function.Timeout,
		MemorySize:  function.MemorySize,
		ImageUri:    function.ImageUri,
		Environment: function.Environment,
	}
}

func (s *LambdaService) ensureFunctionContainer(function *lambdastore.Function, ver *lambdastore.Version, store *lambdastore.FunctionStore, region string) (string, error) {
	ctx := context.Background()

	version := "$LATEST"
	if ver != nil {
		version = ver.Version
	}

	containerName := fmt.Sprintf("lambda-%s-%s-%s", region, function.FunctionName, sanitizeForContainerName(version))

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

	// Serialise creation per container name. Waiters re-check the cache
	// after acquiring: the goroutine that held the lock created the
	// container and recorded its ID.
	muAny, _ := s.containerEnsureMu.LoadOrStore(containerName, &sync.Mutex{})
	mu := muAny.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	if hintAny, ok := s.containerIDs.Load(containerName); ok {
		if hint, _ := hintAny.(string); hint != "" {
			if status, err := s.dockerClient.GetContainerStatus(ctx, hint); err == nil && status == mobyclient.ContainerStatusRunning {
				return hint, nil
			}
		}
	}

	execCfg := executionConfigFor(function, ver)

	image := s.getRuntimeImage(execCfg.Runtime)
	if execCfg.ImageUri != "" {
		image = execCfg.ImageUri
	}

	envVars := map[string]string{
		"AWS_LAMBDA_FUNCTION_TIMEOUT":     fmt.Sprintf("%d", execCfg.Timeout),
		"AWS_LAMBDA_FUNCTION_MEMORY_SIZE": fmt.Sprintf("%d", execCfg.MemorySize),
		"AWS_LAMBDA_FUNCTION_HANDLER":     execCfg.Handler,
		"AWS_LAMBDA_FUNCTION_NAME":        function.FunctionName,
		"AWS_LAMBDA_FUNCTION_VERSION":     version,
		"AWS_REGION":                      region,
	}

	if execCfg.Environment != nil {
		for k, v := range execCfg.Environment.Variables {
			envVars[k] = v
		}
	}

	if _, ok := envVars["AWS_ENDPOINT_URL"]; !ok && s.hostEndpoint != "" {
		envVars["AWS_ENDPOINT_URL"] = s.hostEndpoint
	}

	cfg := mobyclient.AdvancedContainerConfig{
		Name:       containerName,
		Image:      image,
		PullImage:  true,
		Env:        envVars,
		Entrypoint: []string{"/lambda-entrypoint.sh"},
		Cmd:        []string{execCfg.Handler},
		Network:    "bridge",
		AutoRemove: false,
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
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
		// Clean up the previous container to prevent resource leaks when the
		// function is updated or the old container has crashed.
		if originalContainerID != "" && originalContainerID != result.ID {
			if rmErr := s.dockerClient.RemoveContainer(ctx, originalContainerID, true); rmErr != nil {
				logs.Warn("Failed to remove previous container", logs.String("containerID", originalContainerID), logs.Err(rmErr))
			}
		}
	} else if ver != nil {
		originalVerContainerID := ver.ContainerID
		ver.ContainerID = result.ID
		ver.ContainerImageID = result.ID
		if err := store.Update(function); err != nil {
			ver.ContainerID = originalVerContainerID
			ver.ContainerImageID = originalVerContainerID
			if rmErr := s.dockerClient.RemoveContainer(ctx, result.ID, true); rmErr != nil {
				logs.Warn("Failed to remove container during rollback", logs.String("containerID", result.ID), logs.Err(rmErr))
			}
			return "", fmt.Errorf("failed to update function version: %w", err)
		}
		// Clean up the previous version container to prevent resource leaks.
		if originalVerContainerID != "" && originalVerContainerID != result.ID {
			if rmErr := s.dockerClient.RemoveContainer(ctx, originalVerContainerID, true); rmErr != nil {
				logs.Warn("Failed to remove previous version container", logs.String("containerID", originalVerContainerID), logs.Err(rmErr))
			}
		}
	}

	s.containerIDs.Store(containerName, result.ID)

	return result.ID, nil
}

func (s *LambdaService) copyCodeToContainer(containerID string, code []byte, runtime lambdastore.Runtime) error {
	ctx := context.Background()

	reader, err := zip.NewReader(bytes.NewReader(code), int64(len(code)))
	if err != nil {
		return fmt.Errorf("invalid deployment package: failed to read zip: %w", err)
	}

	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open zip entry %s: %w", f.Name, err)
		}

		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("failed to read zip entry %s: %w", f.Name, err)
		}

		destPath := fmt.Sprintf("/var/task/%s", f.Name)
		if err := s.dockerClient.CreateFileInContainer(ctx, containerID, destPath, data); err != nil {
			return fmt.Errorf("failed to copy %s to container: %w", f.Name, err)
		}
	}
	return nil
}

func (s *LambdaService) invokeFunction(function *lambdastore.Function, ver *lambdastore.Version, store *lambdastore.FunctionStore, region string, payload []byte, logType string) (*lambdastore.InvocationResult, error) {
	// Enforce reserved concurrency limit if configured.
	// ReservedConcurrency=0 means the function is effectively paused
	// (zero concurrent executions allowed); any non-nil value must be
	// enforced.
	if function.ReservedConcurrency != nil {
		if *function.ReservedConcurrency == 0 {
			return nil, ErrTooManyRequests
		}
		counter := s.getInflightCounter(function.FunctionArn)
		for {
			current := counter.Load()
			if current >= int32(*function.ReservedConcurrency) {
				return nil, ErrTooManyRequests
			}
			if counter.CompareAndSwap(current, current+1) {
				break
			}
		}
		defer counter.Add(-1)
	}

	ctx := context.Background()

	containerID, err := s.ensureFunctionContainer(function, ver, store, region)
	if err != nil {
		return nil, err
	}

	version := "$LATEST"
	if ver != nil {
		version = ver.Version
	}

	execCfg := executionConfigFor(function, ver)

	code, err := s.loadCode(function.FunctionName, version, region)
	if err != nil {
		// A zip-packaged function must have its code on disk; executing an
		// empty container would silently run stale or no code. Container
		// image packages carry no zip archive.
		if execCfg.ImageUri == "" && function.PackageType != "Image" {
			logs.Error("Function code unavailable for invocation",
				logs.String("function", function.FunctionName),
				logs.String("version", version),
				logs.Err(err))
			return nil, NewLambdaError("ServiceException",
				fmt.Sprintf("The code of function %s version %s is not available.", function.FunctionName, version),
				http.StatusInternalServerError)
		}
	}
	if len(code) > 0 {
		if err := s.copyCodeToContainer(containerID, code, execCfg.Runtime); err != nil {
			return nil, fmt.Errorf("failed to copy code to container: %w", err)
		}
	}

	handlerParts := strings.Split(execCfg.Handler, ".")
	moduleFile := handlerParts[0]
	handlerFunc := "handler"
	if len(handlerParts) > 1 {
		handlerFunc = handlerParts[1]
	}

	eventJSON := "{}"
	if len(payload) > 0 {
		eventJSON = string(payload)
	}

	invokeCmd, resultMarker := s.buildInvokeCommand(execCfg.Runtime, moduleFile, handlerFunc, eventJSON)

	execResult, err := s.dockerClient.Exec(ctx, containerID, mobyclient.ExecConfig{
		Cmd:          invokeCmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exec in container: %w", err)
	}

	// The framed region between the markers is exactly the return value and
	// the output before the opening marker is the handler's console output;
	// only that part belongs in the CloudWatch logs and the tailed
	// LogResult. A custom runtime writes no marker, so its whole stdout
	// remains the payload.
	payloadOut := strings.TrimSpace(execResult.Stdout)
	logOut := payloadOut
	if resultMarker != "" {
		payloadOut, logOut = splitResultPayload(execResult.Stdout, resultMarker)
	}

	s.writeLambdaLogs(function.FunctionName, version, logOut, execResult.Stderr, region)

	functionError := classifyFunctionError(execResult.ExitCode, payloadOut)

	return &lambdastore.InvocationResult{
		StatusCode:      http.StatusOK,
		ExecutedVersion: version,
		Payload:         []byte(payloadOut),
		FunctionError:   functionError,
		LogResult:       captureLogResult(logType, logOut, execResult.Stderr),
	}, nil
}

// finalJSONDocument returns the JSON document the runtime wrapper appended
// to stdout. Handler console output lands in stdout first and the wrapper
// writes the returned payload afterwards, so when the whole output does not
// parse as JSON the payload is the final line. The wrappers serialise the
// return value without indentation, so the payload document is always a
// single line.
func finalJSONDocument(stdout []byte) (json.RawMessage, bool) {
	trimmed := bytes.TrimSpace(stdout)
	if json.Valid(trimmed) {
		return json.RawMessage(trimmed), true
	}
	if idx := bytes.LastIndexByte(trimmed, '\n'); idx >= 0 {
		if line := bytes.TrimSpace(trimmed[idx+1:]); json.Valid(line) {
			return json.RawMessage(line), true
		}
	}
	return nil, false
}

// splitResultPayload separates the wrapper-framed return value from the
// handler's console output in stdout. The wrapper writes the marker, the
// serialised return value, and the marker again, so the region between the
// two marker occurrences is exactly the return value — even when a string
// return spans multiple lines — and everything before the first marker is
// log output. Stdout without a complete marker pair carries no return
// value: the execution was killed before the wrapper could write one.
func splitResultPayload(stdout, marker string) (payload, logs string) {
	open := strings.Index(stdout, marker)
	if open < 0 {
		return "", stdout
	}
	rest := stdout[open+len(marker):]
	closeIdx := strings.Index(rest, marker)
	if closeIdx < 0 {
		return "", stdout
	}
	return rest[:closeIdx], stdout[:open]
}

// classifyFunctionError maps a failed execution onto the AWS wire contract
// for X-Amz-Function-Error. A non-zero exit means the runtime intercepted
// the failure ("Unhandled"); a successful exit whose payload is the error
// document the function returned — a JSON envelope carrying errorMessage —
// is a "Handled" error. A payload without that envelope is a normal
// success and reports no function error. The envelope is read from the
// final JSON document of stdout so handler log lines cannot mask it.
func classifyFunctionError(exitCode int, stdout string) string {
	if exitCode != 0 {
		return "Unhandled"
	}
	doc, ok := finalJSONDocument([]byte(stdout))
	if !ok {
		return ""
	}
	var probe struct {
		ErrorMessage *string `json:"errorMessage"`
	}
	if err := json.Unmarshal(doc, &probe); err == nil && probe.ErrorMessage != nil {
		return "Handled"
	}
	return ""
}

// captureLogResult captures the last 4 KB of execution logs as base64 when
// LogType=Tail is requested, matching the AWS Lambda Invoke API behaviour.
func captureLogResult(logType, stdout, stderr string) string {
	if logType != "Tail" {
		return ""
	}
	logContent := stdout
	if stderr != "" {
		logContent += "\n" + stderr
	}
	if len(logContent) > 4096 {
		logContent = logContent[len(logContent)-4096:]
	}
	return base64.StdEncoding.EncodeToString([]byte(logContent))
}

func (s *LambdaService) buildInvokeCommand(runtime lambdastore.Runtime, moduleFile, handlerFunc, eventJSON string) ([]string, string) {
	// The node/python wrappers frame the serialised return value with a
	// per-invocation marker. Handler code cannot emit the marker (a fresh
	// one is generated per invocation), so the framed region is exactly the
	// return value and the preceding output is console logging — mirroring
	// the AWS runtime contract where the response payload and the logs are
	// separate channels. A custom runtime keeps writing its response as the
	// whole stdout, so it gets no marker.
	marker := ""
	if strings.HasPrefix(string(runtime), "nodejs") || strings.HasPrefix(string(runtime), "python") {
		marker = fmt.Sprintf("__VORPALSTACKS_RESULT_%s__", uuid.New().String())
	}
	if strings.HasPrefix(string(runtime), "nodejs") {
		escaped := strings.ReplaceAll(strings.ReplaceAll(eventJSON, `\`, `\\`), "'", `\'`)
		script := fmt.Sprintf(
			"const m=require('/var/task/%s');const h=typeof m==='function'?m:m['%s'];const p=Promise.resolve(h(JSON.parse('%s')));p.then(r=>{if(r&&typeof r==='object')process.stdout.write('%s'+JSON.stringify(r)+'%s');else if(r!==undefined)process.stdout.write('%s'+String(r)+'%s');}).catch(e=>{process.stderr.write(e.message||String(e));process.exit(1);});",
			moduleFile, handlerFunc, escaped, marker, marker, marker, marker,
		)
		return []string{"node", "-e", script}, marker
	}
	if strings.HasPrefix(string(runtime), "python") {
		escaped := strings.ReplaceAll(eventJSON, "'", `\'`)
		return []string{"python3", "-c", fmt.Sprintf(
			"import json,sys;mod=__import__('%s');h=getattr(mod,'%s',mod);r=h(json.loads('%s'));print('%s'+(json.dumps(r) if isinstance(r,dict) else str(r))+'%s')",
			moduleFile, handlerFunc, escaped, marker, marker,
		)}, marker
	}
	return []string{"/var/runtime/bootstrap"}, marker
}

func (s *LambdaService) writeLambdaLogs(functionName, version, stdout, stderr, region string) {
	logGroupName := "/aws/lambda/" + functionName
	now := time.Now().UTC()
	requestID := uuid.New().String()
	streamName := fmt.Sprintf("%d/%02d/%02d/[%s]%s",
		now.Year(), now.Month(), now.Day(), version, requestID[:8])

	ts := now.UnixNano() / int64(time.Millisecond)
	busEvents := []eventbus.LogEntry{
		{Timestamp: ts, Message: fmt.Sprintf("START RequestId: %s", requestID)},
	}

	for _, line := range strings.Split(stdout, "\n") {
		if line != "" {
			busEvents = append(busEvents, eventbus.LogEntry{Timestamp: ts, Message: line})
		}
	}
	for _, line := range strings.Split(stderr, "\n") {
		if line != "" {
			busEvents = append(busEvents, eventbus.LogEntry{Timestamp: ts, Message: line})
		}
	}

	busEvents = append(busEvents, eventbus.LogEntry{
		Timestamp: ts,
		Message:   fmt.Sprintf("END RequestId: %s", requestID),
	})
	busEvents = append(busEvents, eventbus.LogEntry{
		Timestamp: ts,
		Message:   fmt.Sprintf("REPORT RequestId: %s\tDuration: 0.00 ms\tBilled Duration: 0 ms\tMemory Size: 128 MB\tMax Memory Used: 0 MB", requestID),
	})

	if s.bus != nil {
		logEvt := &eventbus.LambdaLogWriteEvent{
			FunctionName: functionName,
			Version:      version,
			LogGroup:     logGroupName,
			LogStream:    streamName,
			LogEvents:    busEvents,
		}
		logEvt.Region = region
		if err := s.bus.Publish(context.Background(), logEvt); err != nil {
			logs.Warn("failed to publish lambda log event", logs.Err(err))
		}
		return
	}

	s.writeLambdaLogsDirect(logGroupName, streamName, busEvents, functionName, region)
}

// writeLambdaLogsDirect is the fallback path when the event bus is not
// available. It writes Lambda execution logs directly to the CloudWatch
// Logs store without applying metric or subscription filters.
func (s *LambdaService) writeLambdaLogsDirect(logGroupName, logStreamName string, events []eventbus.LogEntry, functionName, region string) {
	if s.logsInvoker == nil {
		return
	}

	ctx := context.Background()

	if err := s.logsInvoker.EnsureLogGroup(ctx, region, logGroupName, s.accountID); err != nil {
		return
	}

	if err := s.logsInvoker.EnsureLogStream(ctx, region, logGroupName, logStreamName); err != nil {
		return
	}

	entries := make([]eventbus.LogsLogEntry, len(events))
	for i, e := range events {
		entries[i] = eventbus.LogsLogEntry(e)
	}

	if err := s.logsInvoker.PutLogEvents(ctx, region, logGroupName, logStreamName, entries); err != nil {
		logs.Warn("Failed to write Lambda logs", logs.String("function", functionName), logs.Err(err))
	}
}

// GetFunctionARN resolves a function name or ARN to its canonical ARN.
// Returns an error if the function does not exist.
func (s *LambdaService) GetFunctionARN(ctx context.Context, functionRef string) (string, error) {
	region := s.region
	functionName := functionRef
	if strings.HasPrefix(functionRef, "arn:") {
		if _, _, arnRegion, _, _ := svcarn.SplitARN(functionRef); arnRegion != "" {
			region = arnRegion
		}
		functionName = svcarn.ExtractFunctionNameFromARN(functionRef)
	}
	store := s.getOrCreateFunctionStore(region)
	fn, err := store.Get(functionName)
	if err != nil {
		return "", err
	}
	return fn.FunctionArn, nil
}

// InvokeForGateway invokes a Lambda function for cross-service integration.
// Accepts either a bare function name or a full Lambda ARN. When an ARN is
// provided, both the region and function name are extracted from it;
// otherwise the constructor region is used as a fallback.
func (s *LambdaService) InvokeForGateway(ctx context.Context, functionRef string, payload []byte) (int64, []byte, error) {
	result, err := s.InvokeForEventSource(ctx, functionRef, payload)
	if err != nil {
		return 0, nil, err
	}
	if result == nil {
		return 0, nil, fmt.Errorf("invocation returned nil result")
	}
	return result.StatusCode, result.Payload, nil
}

// InvokeForEventSource invokes a Lambda function and returns the complete
// invocation result. Event source consumers must observe FunctionError: the
// invoke transport succeeds (HTTP 200) even when the function itself fails,
// and acknowledging a batch in that state would silently drop records that
// AWS semantics require to be retried, bisected or sent to a failure
// destination.
func (s *LambdaService) InvokeForEventSource(ctx context.Context, functionRef string, payload []byte) (*lambdastore.InvocationResult, error) {
	region := s.region
	if strings.HasPrefix(functionRef, "arn:") {
		if _, _, arnRegion, _, _ := svcarn.SplitARN(functionRef); arnRegion != "" {
			region = arnRegion
		}
	}
	functionName, embeddedQualifier := resolveFunctionRef(functionRef)
	store := s.getOrCreateFunctionStore(region)
	function, ver, alias, err := s.resolveQualifier(store, functionName, embeddedQualifier)
	if err != nil {
		return nil, err
	}
	if alias != nil {
		ver = resolveAliasTargetVersion(function, alias)
	}
	return s.invokeFunction(function, ver, store, region, payload, "")
}

// GetFunctionStore returns a new FunctionStore for the Lambda service
// using the constructor region.
func (s *LambdaService) GetFunctionStore() *lambdastore.FunctionStore {
	return s.getOrCreateFunctionStore(s.region)
}

// GetFunctionStoreForRegion returns a FunctionStore for the specified region.
func (s *LambdaService) GetFunctionStoreForRegion(region string) *lambdastore.FunctionStore {
	return s.getOrCreateFunctionStore(region)
}

// IsSubnetInUse implements eventbus.SubnetUsageChecker. It scans all Lambda
// functions in the given region and returns true if any function's
// VpcConfig references the specified subnet ID.
func (s *LambdaService) IsSubnetInUse(ctx context.Context, region, subnetId string) bool {
	store := s.getOrCreateFunctionStore(region)
	functions, err := store.ListAllFunctions()
	if err != nil {
		return false
	}
	for _, fn := range functions {
		if fn.VpcConfig != nil {
			for _, sid := range fn.VpcConfig.SubnetIds {
				if sid == subnetId {
					return true
				}
			}
		}
	}
	return false
}

// IsSecurityGroupInUse implements eventbus.SecurityGroupUsageChecker. It
// scans all Lambda functions in the given region and returns true if any
// function's VpcConfig references the specified security group ID.
func (s *LambdaService) IsSecurityGroupInUse(ctx context.Context, region, sgId string) bool {
	store := s.getOrCreateFunctionStore(region)
	functions, err := store.ListAllFunctions()
	if err != nil {
		return false
	}
	for _, fn := range functions {
		if fn.VpcConfig != nil {
			for _, gid := range fn.VpcConfig.SecurityGroupIds {
				if gid == sgId {
					return true
				}
			}
		}
	}
	return false
}

func (s *LambdaService) getOrCreateFunctionStore(region string) *lambdastore.FunctionStore {
	return s.getOrCreateLambdaStore(region).Functions
}

// GetFunctionPolicy retrieves the resource-based policy for a Lambda function.
func (s *LambdaService) GetFunctionPolicy(functionName string) ([]lambdastore.FunctionPolicy, error) {
	store := s.GetFunctionStore()
	return store.GetPolicy(functionName)
}

// GetAccountSettings returns account limits and usage for Lambda functions.
func (s *LambdaService) GetAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Functions.ListAllFunctions()
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

// Shutdown gracefully shuts down the Lambda service by stopping the ESM
// poller, removing all running Docker containers, and waiting for all
// asynchronous operations to complete.
func (s *LambdaService) Shutdown() {
	s.StopESMPoller()
	s.cleanupAllContainers()
	s.asyncWg.Wait()
}

// cleanupAllContainers iterates every regional store and removes all
// Docker containers associated with Lambda functions ($LATEST and
// published versions). This prevents orphaned containers surviving
// server shutdown.
func (s *LambdaService) cleanupAllContainers() {
	ctx := context.Background()
	s.storeCache.Range(func(key, value any) bool {
		ls, ok := value.(*lambdaStore)
		if !ok || ls == nil || ls.Functions == nil {
			return true
		}
		functions, err := ls.Functions.ListAllFunctions()
		if err != nil {
			logs.Warn("Failed to list functions for container cleanup",
				logs.String("region", fmt.Sprintf("%v", key)), logs.Err(err))
			return true
		}
		for _, fn := range functions {
			if fn.ContainerID != "" {
				if rmErr := s.dockerClient.RemoveContainer(ctx, fn.ContainerID, true); rmErr != nil {
					logs.Warn("Failed to remove container during shutdown",
						logs.String("containerID", fn.ContainerID),
						logs.String("function", fn.FunctionName), logs.Err(rmErr))
				}
			}
			for _, v := range fn.Versions {
				if v.ContainerID != "" && v.ContainerID != fn.ContainerID {
					if rmErr := s.dockerClient.RemoveContainer(ctx, v.ContainerID, true); rmErr != nil {
						logs.Warn("Failed to remove version container during shutdown",
							logs.String("containerID", v.ContainerID),
							logs.String("function", fn.FunctionName),
							logs.String("version", v.Version), logs.Err(rmErr))
					}
				}
			}
		}
		return true
	})
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
