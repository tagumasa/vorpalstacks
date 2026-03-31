package dispatcher

import (
	"context"
	"net/http"
	"os"
	"sync"

	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/services/aws/common/audit"
	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/interfaces"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/mock"
	"vorpalstacks/internal/store/api"
	"vorpalstacks/pkg/graphengine"
)

// Dispatcher routes incoming AWS API requests to appropriate handlers.
// It manages service operations, authentication, and response formatting.
type Dispatcher struct {
	serviceStore         *api.ServiceStore
	operationStore       *api.OperationStore
	shapeStore           *api.ShapeStore
	memberStore          *api.MemberStore
	configStore          *api.ConfigStore
	builder              *resilience.OperationResilienceBuilder
	handlers             map[string]Handler
	storageManager       *storage.RegionStorageManager
	storeProvider        interfaces.StoreProvider
	auditConfig          *audit.AuditConfig
	serviceModeCache     map[string]api.ServiceMode
	serviceModeCacheMu   sync.RWMutex
	handlersMu           sync.RWMutex
	handlerExists        map[string]bool
	mockGenerator        *mock.Generator
	authorizer           *Authorizer
	authorizationEnabled bool
	accountID            string
	graphDB              *graphengine.DB
}

// Handler processes an AWS API request and returns a response or error.
type Handler func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error)

// NewDispatcher creates a new Dispatcher instance with the provided configuration.
// It initializes service stores, handlers, and sets up authorization if enabled.
func NewDispatcher(
	serviceStore *api.ServiceStore,
	operationStore *api.OperationStore,
	shapeStore *api.ShapeStore,
	memberStore *api.MemberStore,
	configStore *api.ConfigStore,
	resilienceConfig *resilience.ServiceResilienceConfig,
	storageMgr *storage.RegionStorageManager,
	storeProvider interfaces.StoreProvider,
	authorizer *Authorizer,
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
		storeProvider:        storeProvider,
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

// Dispatch handles an incoming HTTP request by routing it to the appropriate service handler.
// It manages service mode (fallback, error injection, implemented), authentication, and response writing.
func (d *Dispatcher) Dispatch(w http.ResponseWriter, r *http.Request, serviceName string, operation *api.Operation) {
	parsedReq, parseErr := request.ParseAWSRequest(r)
	opName := ""
	if parsedReq != nil {
		opName = parsedReq.Operation
	}

	if parseErr == nil && opName != "" {
		if handler, exists := d.getHandler(serviceName, opName); exists && handler != nil {
			d.executeWithRecovery(w, r, serviceName, opName, func() {
				httpCtx := r.Context()
				reqCtx := request.NewRequestContext(httpCtx, d.storageManager, d.accountID, parsedReq.GetRegion())
				reqCtx.SetStoreProvider(d.storeProvider)
				reqCtx.SetGraphDBManager(d.graphDB)

				if d.authorizationEnabled && d.authorizer != nil {
					authzResult, err := d.authorizer.Authorize(httpCtx, reqCtx, parsedReq, serviceName, r)
					if err != nil {
						d.handleErrorForRequest(w, r, err)
						return
					}
					if !authzResult.Allowed {
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
			})
			return
		}
	}

	if serviceName == "" {
		contentType := d.getErrorContentTypeForRequest(r)
		awserrors.WriteAWSError(w, awserrors.ErrInvalidAction, contentType)
		return
	}

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

	if dispatchCtx.Operation != "" {
		if handler, exists := d.getHandler(dispatchCtx.ServiceName, dispatchCtx.Operation); exists && handler != nil {
			d.executeWithRecovery(w, r, dispatchCtx.ServiceName, dispatchCtx.Operation, func() {
				httpCtx := r.Context()
				parsedReq := d.buildParsedRequest(r, dispatchCtx)
				reqCtx := request.NewRequestContext(httpCtx, d.storageManager, d.accountID, parsedReq.GetRegion())
				reqCtx.SetStoreProvider(d.storeProvider)
				reqCtx.SetGraphDBManager(d.graphDB)

				if d.authorizationEnabled && d.authorizer != nil {
					authzResult, err := d.authorizer.Authorize(httpCtx, reqCtx, parsedReq, dispatchCtx.ServiceName, r)
					if err != nil {
						d.handleErrorForRequest(w, r, err)
						return
					}
					if !authzResult.Allowed {
						d.handleErrorForRequest(w, r, awserrors.ErrAccessDenied)
						return
					}
				}

				wrapper := d.builder.BuildWrapper(dispatchCtx.ServiceName, dispatchCtx.Operation)
				result, err := wrapper.ExecuteWithResult(httpCtx, func() (interface{}, error) {
					return handler(httpCtx, reqCtx, parsedReq)
				})
				d.recordAudit(dispatchCtx.ServiceName, dispatchCtx.Operation, reqCtx, parsedReq, result, err)
				if err != nil {
					d.handleErrorForRequest(w, r, err)
					return
				}
				d.writeResponseWithOpName(w, r, dispatchCtx.ServiceName, dispatchCtx.Operation, result)
			})
			return
		}
	}

	d.Dispatch(w, r, dispatchCtx.ServiceName, nil)
}

func (d *Dispatcher) SetGraphDB(db *graphengine.DB) {
	d.graphDB = db
}
