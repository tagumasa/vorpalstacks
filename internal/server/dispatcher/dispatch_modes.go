package dispatcher

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
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
				d.executeHandler(w, r, serviceName, opName, parsedReq, handler)
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

	d.executeHandler(w, r, serviceName, operation.Name, parsedReq, handler)
}
