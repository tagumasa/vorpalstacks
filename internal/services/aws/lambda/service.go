// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"archive/zip"
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
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	s3store "vorpalstacks/internal/store/aws/s3"
	storesqs "vorpalstacks/internal/store/aws/sqs"
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
	s3Objects      map[string]s3store.ObjectStoreInterface
	dockerClient   *mobyclient.Client
	bus            eventbus.Bus
	logsStores     sync.Map // region → *logsstore.Store
	storeCache     sync.Map // region → *lambdaStore
	accountID      string
	region         string
	hostEndpoint   string
	dataDir        string
	dataDirOnce    sync.Once
	asyncWg        sync.WaitGroup // goroutine tracking for InvokeAsync
	s3ObjectsMu    sync.RWMutex
	esmPoller      *esmPoller
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
func NewLambdaService(dockerClient *mobyclient.Client, accountID, region, dataDir string) *LambdaService {
	return &LambdaService{
		s3Objects:    make(map[string]s3store.ObjectStoreInterface),
		dockerClient: dockerClient,
		accountID:    accountID,
		region:       region,
		dataDir:      dataDir,
	}
}

// SetLogsStore registers a CloudWatch Logs store for a given region for Lambda log delivery.
func (s *LambdaService) SetLogsStore(region string, store *logsstore.Store) {
	if store != nil {
		s.logsStores.Store(region, store)
	}
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

// SetSQSStore injects the SQS store used by the ESM poller to receive
// messages from source queues. Must be called before StartESMPoller.
func (s *LambdaService) SetSQSStore(sqss *storesqs.SQSStore) {
	if s.esmPoller == nil {
		s.esmPoller = newESMPoller(0, 0, nil)
	}
	s.esmPoller.sqsStore = sqss
}

// SetKinesisStore injects a Kinesis store used by the ESM poller to read
// records from source streams. Must be called before StartESMPoller.
func (s *LambdaService) SetKinesisStore(ks kinesisstore.KinesisStoreInterface) {
	if s.esmPoller == nil {
		s.esmPoller = newESMPoller(0, 0, nil)
	}
	s.esmPoller.kinesisStore = ks
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
	s.esmPoller.accountID = s.accountID
	s.esmPoller.region = s.region
	s.esmPoller.storageManager = s.storageManager
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

func (s *LambdaService) getS3ObjectStore(region string) s3store.ObjectStoreInterface {
	s.s3ObjectsMu.RLock()
	defer s.s3ObjectsMu.RUnlock()
	return s.s3Objects[region]
}

func (s *LambdaService) getLogsStore(region string) (*logsstore.Store, error) {
	if cached, ok := s.logsStores.Load(region); ok {
		if typed, ok := cached.(*logsstore.Store); ok {
			return typed, nil
		}
	}
	regionalStore := s.getRegionalStorage(region)
	if regionalStore == nil {
		return nil, fmt.Errorf("lambda: no regional storage available for %s", region)
	}
	store, err := logsstore.NewStore(regionalStore, regionalStore.Bucket("logs-"+region), s.accountID, region, s.dataDir)
	if err != nil {
		return nil, err
	}
	if actual, loaded := s.logsStores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(*logsstore.Store); ok {
			return typed, nil
		}
	}
	return store, nil
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

	d.RegisterHandlerForService("lambda", "PutProvisionedConcurrencyConfig", s.PutProvisionedConcurrencyConfig)
	d.RegisterHandlerForService("lambda", "GetProvisionedConcurrencyConfig", s.GetProvisionedConcurrencyConfig)
	d.RegisterHandlerForService("lambda", "DeleteProvisionedConcurrencyConfig", s.DeleteProvisionedConcurrencyConfig)
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
	return io.ReadAll(io.LimitReader(reader, 250*1024*1024))
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

	codePath := fmt.Sprintf("%s/%s/code/%s/%s/code.zip", dataDir, naming.SanitizePathComponent(region), naming.SanitizePathComponent(functionName), naming.SanitizePathComponent(version))
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

	if _, ok := envVars["AWS_ENDPOINT_URL"]; !ok && s.hostEndpoint != "" {
		envVars["AWS_ENDPOINT_URL"] = s.hostEndpoint
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

	reader, err := zip.NewReader(bytes.NewReader(code), int64(len(code)))
	if err == nil {
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

	return s.dockerClient.CreateFileInContainer(ctx, containerID, "/var/task/index.js", code)
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

	handlerParts := strings.Split(function.Handler, ".")
	moduleFile := handlerParts[0]
	handlerFunc := "handler"
	if len(handlerParts) > 1 {
		handlerFunc = handlerParts[1]
	}

	eventJSON := "{}"
	if len(payload) > 0 {
		eventJSON = string(payload)
	}

	invokeCmd := s.buildInvokeCommand(function.Runtime, moduleFile, handlerFunc, eventJSON)

	execResult, err := s.dockerClient.Exec(ctx, containerID, mobyclient.ExecConfig{
		Cmd:          invokeCmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exec in container: %w", err)
	}

	stdout := strings.TrimSpace(execResult.Stdout)

	s.writeLambdaLogs(function.FunctionName, version, stdout, execResult.Stderr, region)

	var functionError string
	if execResult.ExitCode != 0 {
		functionError = execResult.Stderr
	}

	return &lambdastore.InvocationResult{
		StatusCode:      http.StatusOK,
		ExecutedVersion: version,
		Payload:         []byte(stdout),
		FunctionError:   functionError,
	}, nil
}

func (s *LambdaService) buildInvokeCommand(runtime lambdastore.Runtime, moduleFile, handlerFunc, eventJSON string) []string {
	if strings.HasPrefix(string(runtime), "nodejs") {
		escaped := strings.ReplaceAll(strings.ReplaceAll(eventJSON, `\`, `\\`), "'", `\'`)
		script := fmt.Sprintf(
			"const m=require('/var/task/%s');const h=typeof m==='function'?m:m['%s'];const p=Promise.resolve(h(JSON.parse('%s')));p.then(r=>{if(r&&typeof r==='object')process.stdout.write(JSON.stringify(r));else if(r!==undefined)process.stdout.write(String(r));}).catch(e=>{process.stderr.write(e.message||String(e));process.exit(1);});",
			moduleFile, handlerFunc, escaped,
		)
		return []string{"node", "-e", script}
	}
	if strings.HasPrefix(string(runtime), "python") {
		escaped := strings.ReplaceAll(eventJSON, "'", `\'`)
		return []string{"python3", "-c", fmt.Sprintf(
			"import json,sys;mod=__import__('%s');h=getattr(mod,'%s',mod);r=h(json.loads('%s'));print(json.dumps(r) if isinstance(r,dict) else str(r))",
			moduleFile, handlerFunc, escaped,
		)}
	}
	return []string{"/var/runtime/bootstrap"}
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
	logsStore, err := s.getLogsStore(region)
	if err != nil {
		return
	}

	_, err = logsStore.GetLogGroup(logGroupName)
	if err != nil {
		lg := logsstore.NewLogGroup(logGroupName, region, s.accountID)
		if createErr := logsStore.CreateLogGroup(lg); createErr != nil {
			return
		}
	}

	ls := logsstore.NewLogStream(logStreamName, logGroupName)
	if createErr := logsStore.CreateLogStream(ls); createErr != nil {
		return
	}

	storeEvents := make([]logsstore.LogEntry, len(events))
	for i, e := range events {
		storeEvents[i] = logsstore.LogEntry{Timestamp: e.Timestamp, Message: e.Message}
	}

	if _, err := logsStore.PutLogEvents(logGroupName, logStreamName, storeEvents); err != nil {
		logs.Warn("Failed to write Lambda logs", logs.String("function", functionName), logs.Err(err))
	}
}

// InvokeForGateway invokes a Lambda function for cross-service integration.
// Accepts either a bare function name or a full Lambda ARN. When an ARN is
// provided, both the region and function name are extracted from it;
// otherwise the constructor region is used as a fallback.
func (s *LambdaService) InvokeForGateway(ctx context.Context, functionRef string, payload []byte) (int64, []byte, error) {
	region := s.region
	functionName := functionRef
	if strings.HasPrefix(functionRef, "arn:") {
		if _, _, arnRegion, _, _ := svcarn.SplitARN(functionRef); arnRegion != "" {
			region = arnRegion
		}
		functionName = svcarn.ExtractFunctionNameFromARN(functionRef)
	}
	store := s.getOrCreateFunctionStore(region)
	function, ver, _, err := s.resolveQualifier(store, functionName, "")
	if err != nil {
		return 0, nil, err
	}
	result, err := s.invokeFunction(function, ver, store, region, payload)
	if err != nil {
		return 0, nil, err
	}
	if result == nil {
		return 0, nil, fmt.Errorf("invocation returned nil result")
	}
	return result.StatusCode, result.Payload, nil
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

func (s *LambdaService) getOrCreateFunctionStore(region string) *lambdastore.FunctionStore {
	if cached, ok := s.storeCache.Load(region); ok {
		if typed, ok := cached.(*lambdaStore); ok {
			return typed.Functions
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
			return typed.Functions
		}
	}
	return newStore.Functions
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
	result, err := store.Functions.List(storecommon.ListOptions{MaxItems: 1000})
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
