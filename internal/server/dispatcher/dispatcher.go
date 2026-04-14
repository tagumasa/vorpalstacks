package dispatcher

import (
	"context"
	"net/http"
	"os"
	"sync"

	"vorpalstacks/internal/common/audit"
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/mock"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/http/classifier"
	"vorpalstacks/internal/store/api"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/pkg/graphengine"
)

// Authorizer evaluates IAM policies to determine whether a request should be allowed.
type Authorizer interface {
	Authorize(ctx context.Context, reqCtx *request.RequestContext, parsedReq *request.ParsedRequest, serviceName string, r *http.Request) (bool, error)
}

// Dispatcher routes incoming AWS API requests to appropriate handlers.
// It manages service operations, authentication, and response formatting.
type Dispatcher struct {
	serviceStore              *api.ServiceStore
	operationStore            *api.OperationStore
	shapeStore                *api.ShapeStore
	memberStore               *api.MemberStore
	configStore               *api.ConfigStore
	builder                   *resilience.OperationResilienceBuilder
	handlers                  map[string]Handler
	storageManager            *storage.RegionStorageManager
	iamStore                  iamstore.IAMStoreInterface
	auditConfig               *audit.AuditConfig
	serviceModeCache          map[string]api.ServiceMode
	serviceModeCacheMu        sync.RWMutex
	handlersMu                sync.RWMutex
	handlerExists             map[string]bool
	mockGenerator             *mock.Generator
	authorizer                Authorizer
	authorizationEnabled      bool
	accountID                 string
	graphDB                   *graphengine.DB
	cloudTrailRecorderFactory CloudTrailRecorderFactory
}

// Handler processes an AWS API request and returns a response or error.
type Handler = handler.Handler

// NewDispatcher creates a new Dispatcher instance with the provided configuration.
// It initialises service stores, handlers, and sets up authorisation if enabled.
func NewDispatcher(
	serviceStore *api.ServiceStore,
	operationStore *api.OperationStore,
	shapeStore *api.ShapeStore,
	memberStore *api.MemberStore,
	configStore *api.ConfigStore,
	resilienceConfig *resilience.ServiceResilienceConfig,
	storageMgr *storage.RegionStorageManager,
	iamStore iamstore.IAMStoreInterface,
	authorizer Authorizer,
	accountID string,
) *Dispatcher {
	authEnabled := os.Getenv("AUTHORIZATION_ENABLED") == "true"

	d := &Dispatcher{
		serviceStore:         serviceStore,
		operationStore:       operationStore,
		shapeStore:           shapeStore,
		memberStore:          memberStore,
		configStore:          configStore,
		builder:              resilience.NewOperationResilienceBuilder(resilienceConfig),
		handlers:             make(map[string]Handler),
		storageManager:       storageMgr,
		iamStore:             iamStore,
		auditConfig:          audit.DefaultConfig(),
		serviceModeCache:     make(map[string]api.ServiceMode),
		handlerExists:        make(map[string]bool),
		mockGenerator:        mock.NewGenerator(shapeStore, memberStore),
		authorizer:           authorizer,
		authorizationEnabled: authEnabled,
		accountID:            accountID,
	}
	d.warmupCache()
	return d
}

// executeHandler runs the shared handler execution pipeline:
// build request context, check authorisation, execute with resilience wrapper,
// record audit, and write the response or error.
func (d *Dispatcher) executeHandler(w http.ResponseWriter, r *http.Request, serviceName, opName string, parsedReq *request.ParsedRequest, handler Handler) {
	httpCtx := r.Context()
	reqCtx := request.NewRequestContext(httpCtx, d.storageManager, d.accountID, parsedReq.GetRegion())
	reqCtx.SetIAMStore(d.iamStore)
	reqCtx.SetGraphDBManager(d.graphDB)

	if d.authorizationEnabled && d.authorizer != nil {
		authzResult, err := d.authorizer.Authorize(httpCtx, reqCtx, parsedReq, serviceName, r)
		if err != nil {
			d.handleErrorForRequest(w, r, err)
			return
		}
		if !authzResult {
			d.handleErrorForRequest(w, r, awserrors.ErrAccessDenied)
			return
		}
	}

	wrapper := d.builder.BuildWrapper(serviceName, opName)
	result, err := wrapper.ExecuteWithResult(httpCtx, func() (interface{}, error) {
		return handler(httpCtx, reqCtx, parsedReq)
	})
	d.recordAudit(serviceName, opName, reqCtx, parsedReq, result, err)
	if err != nil {
		d.handleErrorForRequest(w, r, err)
		return
	}
	d.writeResponseWithOpName(w, r, serviceName, opName, result)
}

// tryExecuteHandler looks up the handler for the given service/operation and
// executes it with recovery if found. Returns true if the handler was executed.
func (d *Dispatcher) tryExecuteHandler(w http.ResponseWriter, r *http.Request, serviceName, opName string, parsedReq *request.ParsedRequest) bool {
	if opName == "" {
		return false
	}
	h, exists := d.getHandler(serviceName, opName)
	if !exists || h == nil {
		return false
	}
	d.executeWithRecovery(w, r, serviceName, opName, func() {
		d.executeHandler(w, r, serviceName, opName, parsedReq, h)
	})
	return true
}

// dispatchByServiceMode resolves the service mode and dispatches to the
// appropriate handler (fallback, error injection, or implemented).
func (d *Dispatcher) dispatchByServiceMode(w http.ResponseWriter, r *http.Request, serviceName, opName string, operation *api.Operation, parsedReq *request.ParsedRequest) {
	mode, hasMode := d.getServiceMode(serviceName)
	if !hasMode {
		cfg, err := d.configStore.Get(serviceName)
		if err != nil || cfg == nil {
			if opName != "" && d.hasHandlerForService(serviceName, opName) {
				d.executeWithRecovery(w, r, serviceName, opName, func() {
					d.handleImplemented(w, r, serviceName, operation, parsedReq)
				})
				return
			}
			contentType := d.getErrorContentTypeForRequest(r)
			awserrors.WriteAWSError(w, awserrors.ErrInternal, contentType)
			return
		}
		mode = cfg.Mode
		d.setServiceMode(serviceName, mode)
	}

	if mode == api.ServiceModeFallback && opName != "" && d.hasHandlerForService(serviceName, opName) {
		mode = api.ServiceModeImplemented
	}

	switch mode {
	case api.ServiceModeFallback:
		d.executeWithRecovery(w, r, serviceName, opName, func() {
			d.handleFallback(w, r, serviceName, operation, parsedReq)
		})
	case api.ServiceModeErrorInjection:
		d.executeWithRecovery(w, r, serviceName, opName, func() {
			cfg, _ := d.configStore.Get(serviceName)
			d.handleErrorInjection(w, r, cfg, parsedReq)
		})
	case api.ServiceModeImplemented:
		d.executeWithRecovery(w, r, serviceName, opName, func() {
			d.handleImplemented(w, r, serviceName, operation, parsedReq)
		})
	default:
		contentType := d.getErrorContentTypeForRequest(r)
		awserrors.WriteAWSError(w, awserrors.ErrInternal, contentType)
	}
}

// Dispatch handles an incoming HTTP request by routing it to the appropriate service handler.
// It manages service mode (fallback, error injection, implemented), authentication, and response writing.
func (d *Dispatcher) Dispatch(w http.ResponseWriter, r *http.Request, serviceName string, operation *api.Operation) {
	parsedReq, parseErr := request.ParseAWSRequest(r)
	opName := ""
	if parsedReq != nil {
		opName = parsedReq.Operation
	}

	if parseErr == nil && d.tryExecuteHandler(w, r, serviceName, opName, parsedReq) {
		return
	}

	if serviceName == "" {
		contentType := d.getErrorContentTypeForRequest(r)
		awserrors.WriteAWSError(w, awserrors.ErrInvalidAction, contentType)
		return
	}

	d.dispatchByServiceMode(w, r, serviceName, opName, operation, parsedReq)
}

// DispatchContext holds the context information for a dispatched request.
type DispatchContext struct {
	ServiceName string
	Operation   string
	Params      map[string]interface{}
}

// DispatchWithContext handles a request using the provided dispatch context.
// It allows for programmatic dispatching with pre-defined service name and operation.
func (d *Dispatcher) DispatchWithContext(w http.ResponseWriter, r *http.Request, dispatchCtx *DispatchContext) {
	if dispatchCtx == nil {
		dispatchCtx = &DispatchContext{}
	}

	parsedReq := d.buildParsedRequest(r, dispatchCtx)
	if d.tryExecuteHandler(w, r, dispatchCtx.ServiceName, dispatchCtx.Operation, parsedReq) {
		return
	}

	d.Dispatch(w, r, dispatchCtx.ServiceName, nil)
}

// SetGraphDB injects the graph database instance into the dispatcher for query support.
func (d *Dispatcher) SetGraphDB(db *graphengine.DB) {
	d.graphDB = db
}

// DispatchClassified handles an incoming HTTP request that has already been
// classified by the request classifier. It converts the ClassifiedRequest to
// a ParsedRequest and routes it through the standard dispatch pipeline,
// including authorisation, resilience wrapping, and audit recording.
func (d *Dispatcher) DispatchClassified(w http.ResponseWriter, r *http.Request, cr *classifier.ClassifiedRequest) {
	parsedReq := classifiedToParsed(cr)

	if d.tryExecuteHandler(w, r, cr.ServiceName, parsedReq.Operation, parsedReq) {
		return
	}

	if cr.ServiceName == "" {
		contentType := d.getErrorContentTypeForRequest(r)
		awserrors.WriteAWSError(w, awserrors.ErrInvalidAction, contentType)
		return
	}

	d.dispatchByServiceMode(w, r, cr.ServiceName, parsedReq.Operation, nil, parsedReq)
}

// classifiedToParsed converts a ClassifiedRequest into a ParsedRequest
// suitable for the existing handler pipeline.
func classifiedToParsed(cr *classifier.ClassifiedRequest) *request.ParsedRequest {
	parsed := &request.ParsedRequest{
		Operation:   cr.Operation,
		Parameters:  cr.Parameters,
		Headers:     cr.Headers,
		QueryParams: cr.QueryParams,
		PathParams:  cr.PathParams,
		Body:        cr.Body,
		Region:      cr.Region,
		AccessKeyID: cr.AccessKeyID,
	}

	if graphID := cr.Headers.Get("Graphidentifier"); graphID != "" {
		parsed.Parameters["graphIdentifier"] = graphID
	}

	lambdaHeaderMappings := map[string]string{
		"X-Amz-Invocation-Type": "InvocationType",
		"X-Amz-Log-Type":        "LogType",
		"X-Amz-Client-Context":  "ClientContext",
		"X-Amz-Function-Error":  "FunctionError",
		"X-Amz-Qualifier":       "Qualifier",
	}
	for header, param := range lambdaHeaderMappings {
		if _, exists := parsed.Parameters[param]; !exists {
			if v := cr.Headers.Get(header); v != "" {
				parsed.Parameters[param] = v
			}
		}
	}

	return parsed
}

// SetCloudTrailRecorderFactory sets the factory function for creating CloudTrail audit recorders.
func (d *Dispatcher) SetCloudTrailRecorderFactory(factory CloudTrailRecorderFactory) {
	d.cloudTrailRecorderFactory = factory
}
