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
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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
		// The deployment zip's permission bits survive the copy: the
		// custom-runtime contract requires /var/task/bootstrap to arrive
		// executable (the provided base image checks -x before exec'ing
		// it). Entries without any exec bit keep the historical 0644 so
		// wrapper-runtime packages land exactly as before.
		perm := f.Mode().Perm()
		if perm&0111 == 0 {
			perm = 0644
		}
		if err := s.dockerClient.CreateFileInContainer(ctx, containerID, destPath, data, perm); err != nil {
			return fmt.Errorf("failed to copy %s to container: %w", f.Name, err)
		}
	}
	return nil
}

func (s *LambdaService) invokeFunction(req invokeRequest) (*lambdastore.InvocationResult, error) {
	function, ver, store, region := req.Function, req.Version, req.Store, req.Region
	payload, logType := req.Payload, req.LogType
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

	// ClientContext arrives base64-encoded on the synchronous Invoke plane
	// (the Smithy model passes it "for synchronous invocations only"); the
	// handler's context object carries the decoded document.
	clientContextJSON := ""
	if req.ClientContextRaw != "" {
		decoded, decErr := base64.StdEncoding.DecodeString(req.ClientContextRaw)
		if decErr != nil {
			return nil, NewInvalidParameter("ClientContext", "must be base64-encoded data")
		}
		clientContextJSON = string(decoded)
	}

	rec := newInvocationRecord(function.FunctionName, version, req.InvokedARN, execCfg.MemorySize, execCfg.Timeout, clientContextJSON)

	invokeCmd, resultMarker := s.buildInvokeCommand(execCfg.Runtime, moduleFile, handlerFunc, eventJSON, rec)

	// Bootstrap runtimes (the provided.* custom runtimes and the RIC-based
	// managed images) receive the event over the Runtime API instead of a
	// wrapper script: a per-invocation host-side HTTP server answers
	// /invocation/next with the event and the Lambda-Runtime-* headers, and
	// the bootstrap POSTs its answer back. The exec environment carries the
	// server address; the image entrypoint's own emulator is bypassed
	// because the exec goes straight to /var/runtime/bootstrap.
	var execEnv []string
	var apiServer *runtimeAPIServer
	if !usesRuntimeWrapper(execCfg.Runtime) {
		apiServer, err = startRuntimeAPI(rec, []byte(eventJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to start the runtime API: %w", err)
		}
		defer apiServer.Close()
		execEnv = append(execEnv, "AWS_LAMBDA_RUNTIME_API="+apiServer.Addr())
	}

	execConfig := mobyclient.ExecConfig{
		Cmd:          invokeCmd,
		AttachStdout: true,
		AttachStderr: true,
		Env:          execEnv,
	}

	execStart := time.Now()
	var execResult *mobyclient.ExecResult
	timedOut := false
	if apiServer != nil {
		// A bootstrap execution settles at the captured answer, its own
		// process exit, or the host-enforced deadline — whichever comes
		// first; see invocation_exec.go.
		outcome := s.runBootstrapExecution(ctx, containerID,
			containerNameFor(region, function.FunctionName, version), execConfig, apiServer, rec.Deadline)
		execResult, timedOut, err = outcome.result, outcome.timedOut, outcome.err
	} else {
		execResult, err = s.dockerClient.Exec(ctx, containerID, execConfig)
		if err == nil {
			timedOut = execResult.ExitCode == execExitTimedOut
		}
	}
	execDuration := time.Since(execStart)
	if err != nil {
		return nil, fmt.Errorf("failed to exec in container: %w", err)
	}

	// The framed region between the markers is exactly the return value and
	// the output before the opening marker is the handler's console output;
	// only that part belongs in the CloudWatch logs and the tailed
	// LogResult. A custom runtime writes no marker; its payload arrives
	// over the Runtime API below, and its stdout is console output.
	payloadOut := strings.TrimSpace(execResult.Stdout)
	logOut := payloadOut
	if resultMarker != "" {
		payloadOut, logOut = splitResultPayload(execResult.Stdout, resultMarker)
	}

	apiAnswered := false
	var apiKind string
	if apiServer != nil {
		// The bootstrap answers over the Runtime API, not stdout: its
		// stdout is console output only. A captured answer outranks the
		// process exit status — a bootstrap that loops back to /next after
		// answering sits idle until the timeout reaps it, which is the
		// platform ending the execution, not a function failure.
		apiBody, kind, captured := apiServer.Captured()
		apiAnswered, apiKind = captured, kind
		if captured {
			payloadOut = string(apiBody)
			timedOut = false
		} else {
			payloadOut = ""
			logOut = strings.TrimSpace(execResult.Stdout)
		}
	}

	if timedOut {
		payloadOut = timeoutEnvelope(rec.TimeoutSeconds)
	}

	logEvents := invocationLogEvents(rec, logOut, execResult.Stderr, execDuration, timedOut)
	s.writeLambdaLogs(rec, region, logEvents)

	functionError := classifyFunctionError(execResult.ExitCode)
	if apiAnswered {
		if apiKind != runtimeAPIResponseKind {
			// The bootstrap reported the failure through an error endpoint;
			// the contract marks such invocations Unhandled.
			functionError = "Unhandled"
		} else {
			// A delivered response is a success whatever the payload
			// contains and whatever the process did after answering — the
			// header may only be present when an error actually occurred.
			functionError = ""
		}
	}

	return &lambdastore.InvocationResult{
		StatusCode:      http.StatusOK,
		ExecutedVersion: version,
		Payload:         []byte(payloadOut),
		FunctionError:   functionError,
		LogResult:       captureLogResult(logType, logEvents),
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

// classifyFunctionError maps the execution outcome onto the AWS wire
// contract for X-Amz-Function-Error: "If present, indicates that an error
// occurred during function execution." A successful exit is therefore
// never an error, whatever the returned payload looks like — on AWS a
// handler returning an error-shaped document is a plain success. The
// classification reads the exit status only: the wrapper runtimes report a
// callback-signalled failure with execExitHandledError (Handled) and an
// uncaught failure, or a timeout kill, with any other non-zero status
// (Unhandled). The Runtime API path overrides this at the call site
// through the captured answer kind.
func classifyFunctionError(exitCode int) string {
	if exitCode == 0 {
		return ""
	}
	if exitCode == execExitHandledError {
		return "Handled"
	}
	return "Unhandled"
}

// timeoutEnvelope is the payload a timed-out invocation answers with: a 200
// response carrying an Unhandled function error whose errorMessage names
// the timeout. AWS documents the "Task timed out after" wording; the
// errorType member varies by runtime generation across AWS services and is
// deliberately omitted here.
func timeoutEnvelope(timeoutSeconds int32) string {
	return fmt.Sprintf(`{"errorMessage":"Task timed out after %.2f seconds"}`, float64(timeoutSeconds))
}

// captureLogResult captures the last 4 KB of execution logs as base64 when
// LogType=Tail is requested, matching the AWS Lambda Invoke API behaviour.
// The tailed text is the full invocation log sequence — START/END/REPORT
// framing included — because the AWS runtime's own lines are part of what
// the tail exposes.
func captureLogResult(logType string, events []eventbus.LogEntry) string {
	if logType != "Tail" {
		return ""
	}
	var b strings.Builder
	for _, e := range events {
		b.WriteString(e.Message)
		b.WriteByte('\n')
	}
	logContent := b.String()
	if len(logContent) > 4096 {
		logContent = logContent[len(logContent)-4096:]
	}
	return base64.StdEncoding.EncodeToString([]byte(logContent))
}

// invocationLogEvents composes the execution's CloudWatch log lines: the
// START framing, the handler's console output, the timeout marker, and the
// END/REPORT framing with the measured duration. The tailed LogResult and
// the written log stream derive from this single sequence, so what the
// caller tails is exactly what lands in the log group.
func invocationLogEvents(rec invocationRecord, stdout, stderr string, duration time.Duration, timedOut bool) []eventbus.LogEntry {
	ts := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	events := []eventbus.LogEntry{
		// AWS log examples frame the start as
		// "START RequestId: <id> Version: <version>".
		{Timestamp: ts, Message: fmt.Sprintf("START RequestId: %s Version: %s", rec.RequestID, rec.Version)},
	}

	for _, line := range strings.Split(stdout, "\n") {
		if line != "" {
			events = append(events, eventbus.LogEntry{Timestamp: ts, Message: line})
		}
	}
	for _, line := range strings.Split(stderr, "\n") {
		if line != "" {
			events = append(events, eventbus.LogEntry{Timestamp: ts, Message: line})
		}
	}

	if timedOut {
		events = append(events, eventbus.LogEntry{Timestamp: ts, Message: "TASK TIMED OUT"})
	}

	// Max Memory Used stays 0 because the exec plane exposes no
	// per-process RSS.
	billedMS := int64(math.Ceil(duration.Seconds() * 1000))
	if billedMS < 1 {
		billedMS = 1
	}
	events = append(events, eventbus.LogEntry{
		Timestamp: ts,
		Message:   fmt.Sprintf("END RequestId: %s", rec.RequestID),
	})
	events = append(events, eventbus.LogEntry{
		Timestamp: ts,
		Message: fmt.Sprintf("REPORT RequestId: %s\tDuration: %.2f ms\tBilled Duration: %d ms\tMemory Size: %d MB\tMax Memory Used: 0 MB",
			rec.RequestID, float64(duration.Microseconds())/1000, billedMS, rec.MemorySize),
	})
	return events
}

func (s *LambdaService) writeLambdaLogs(rec invocationRecord, region string, events []eventbus.LogEntry) {
	// The request id, log group and log stream are the invocation record's:
	// the handler saw the same values in its context object, so
	// context.awsRequestId matches the START/END/REPORT lines.
	if s.bus != nil {
		logEvt := &eventbus.LambdaLogWriteEvent{
			FunctionName: rec.FunctionName,
			Version:      rec.Version,
			LogGroup:     rec.LogGroupName,
			LogStream:    rec.LogStreamName,
			LogEvents:    events,
		}
		logEvt.Region = region
		if err := s.bus.Publish(context.Background(), logEvt); err != nil {
			logs.Warn("failed to publish lambda log event", logs.Err(err))
		}
		return
	}

	s.writeLambdaLogsDirect(rec.LogGroupName, rec.LogStreamName, events, rec.FunctionName, region)
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

// InvokeForTrigger invokes a Lambda function for synchronous trigger
// consumers that must distinguish a failed function execution
// (LambdaInvocation.FunctionError) from an invocation-transport failure
// (the returned error).
func (s *LambdaService) InvokeForTrigger(ctx context.Context, functionRef string, payload []byte) (eventbus.LambdaInvocation, error) {
	result, err := s.InvokeForEventSource(ctx, functionRef, payload)
	if err != nil {
		return eventbus.LambdaInvocation{}, err
	}
	if result == nil {
		return eventbus.LambdaInvocation{}, fmt.Errorf("invocation returned nil result")
	}
	return eventbus.LambdaInvocation{
		StatusCode:    result.StatusCode,
		Payload:       result.Payload,
		FunctionError: result.FunctionError,
	}, nil
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
	return s.invokeFunction(invokeRequest{
		Function:   function,
		Version:    ver,
		Store:      store,
		Region:     region,
		Payload:    payload,
		InvokedARN: qualifiedInvokeARN(function.FunctionArn, embeddedQualifier),
	})
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
	return s.getAccountSettingsCore(store)
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
