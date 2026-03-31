package dispatcher

import (
	"context"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/store/api"
)

func (d *Dispatcher) handleFallback(w http.ResponseWriter, r *http.Request, serviceName string, operation *api.Operation, parsedReq *request.ParsedRequest) {
	if operation == nil {
		operation = d.detectOperationFromParsed(parsedReq)
	}

	response := d.generateFallbackResponse(operation)
	d.writeResponse(w, r, operation, "", response)
}

func (d *Dispatcher) handleErrorInjection(w http.ResponseWriter, r *http.Request, cfg *api.ServiceConfig, parsedReq *request.ParsedRequest) {
	errResp := awserrors.ErrInternal
	if cfg != nil && cfg.Error != nil {
		errResp = awserrors.NewAWSError(cfg.Error.Code, cfg.Error.Message, cfg.Error.HTTPStatusCode)
	}
	operation := d.detectOperationFromParsed(parsedReq)
	contentType := d.getErrorContentType(r, operation)
	awserrors.WriteAWSError(w, errResp, contentType)
}

func (d *Dispatcher) handleImplemented(w http.ResponseWriter, r *http.Request, serviceName string, operation *api.Operation, parsedReq *request.ParsedRequest) {
	if operation == nil {
		operation = d.detectOperationFromParsed(parsedReq)
	}

	if operation == nil {
		opName := ""
		if parsedReq != nil {
			opName = parsedReq.Operation
		}
		if opName != "" {
			if handler, exists := d.getHandler(serviceName, opName); exists && handler != nil {
				httpCtx := r.Context()
				reqCtx := request.NewRequestContext(httpCtx, d.storageManager, d.accountID, parsedReq.GetRegion())
				reqCtx.SetStoreProvider(d.storeProvider)
				reqCtx.SetGraphDBManager(d.graphDB)
				if d.checkAuthorization(w, r, httpCtx, reqCtx, parsedReq, serviceName) {
					return
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
				return
			}
		}
		contentType := d.getErrorContentTypeForRequest(r)
		awserrors.WriteAWSError(w, awserrors.ErrInvalidAction, contentType)
		return
	}

	handler, exists := d.getHandler(serviceName, operation.Name)
	if !exists || handler == nil {
		contentType := d.getErrorContentType(r, operation)
		awserrors.WriteAWSError(w, awserrors.ErrNotImplemented, contentType)
		return
	}

	if parsedReq == nil {
		contentType := d.getErrorContentType(r, operation)
		awserrors.WriteAWSError(w, awserrors.ErrInternal, contentType)
		return
	}

	httpCtx := r.Context()
	reqCtx := request.NewRequestContext(httpCtx, d.storageManager, d.accountID, parsedReq.GetRegion())
	reqCtx.SetStoreProvider(d.storeProvider)
	reqCtx.SetGraphDBManager(d.graphDB)
	if d.checkAuthorization(w, r, httpCtx, reqCtx, parsedReq, serviceName) {
		return
	}
	wrapper := d.builder.BuildWrapper(serviceName, operation.Name)
	result, err := wrapper.ExecuteWithResult(httpCtx, func() (interface{}, error) {
		return handler(httpCtx, reqCtx, parsedReq)
	})
	d.recordAudit(serviceName, operation.Name, reqCtx, parsedReq, result, err)
	if err != nil {
		d.handleError(w, r, operation, err)
		return
	}

	d.writeResponse(w, r, operation, operation.Name, result)
}

func (d *Dispatcher) checkAuthorization(w http.ResponseWriter, r *http.Request, httpCtx context.Context, reqCtx *request.RequestContext, parsedReq *request.ParsedRequest, serviceName string) bool {
	if d.authorizationEnabled && d.authorizer != nil {
		authzResult, err := d.authorizer.Authorize(httpCtx, reqCtx, parsedReq, serviceName, r)
		if err != nil {
			d.handleErrorForRequest(w, r, err)
			return true
		}
		if !authzResult.Allowed {
			d.handleErrorForRequest(w, r, awserrors.ErrAccessDenied)
			return true
		}
	}
	return false
}
