package handler

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Handler processes an AWS API request and returns a response or error.
type Handler func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error)

// Registrar is the interface for registering API handlers.
type Registrar interface {
	RegisterHandler(operationName string, handler Handler)
	RegisterHandlerForService(serviceName, operationName string, handler Handler)
}
