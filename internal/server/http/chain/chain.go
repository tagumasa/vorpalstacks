// Package chain provides HTTP middleware chain functionality for vorpalstacks.
package chain

import (
	"context"
	"net/http"
)

// RequestContext holds the request and response state for the handler chain.
type RequestContext struct {
	context.Context

	Request      *http.Request
	Response     http.ResponseWriter
	ServiceName  string
	Operation    string
	Params       map[string]interface{}
	StatusCode   int
	ResponseBody interface{}
	Error        error

	handled bool
}

// NewRequestContext creates a new RequestContext with the given context, response writer, and request.
func NewRequestContext(ctx context.Context, w http.ResponseWriter, r *http.Request) *RequestContext {
	return &RequestContext{
		Context:  ctx,
		Request:  r,
		Response: w,
		Params:   make(map[string]interface{}),
	}
}

// SetHandled marks the request as handled.
func (c *RequestContext) SetHandled() {
	c.handled = true
}

// IsHandled returns whether the request has been handled.
func (c *RequestContext) IsHandled() bool {
	return c.handled
}

// SetError sets the error for the request context.
func (c *RequestContext) SetError(err error) {
	c.Error = err
}

// SetResponse sets the status code and response body for the request.
func (c *RequestContext) SetResponse(statusCode int, body interface{}) {
	c.StatusCode = statusCode
	c.ResponseBody = body
}

// RequestHandler is a function that processes a request and returns an error.
type RequestHandler func(ctx *RequestContext) error

// ExceptionHandler is a function that handles errors that occur during request processing.
type ExceptionHandler func(ctx *RequestContext, err error) error

// ResponseHandler is a function that processes the response.
type ResponseHandler func(ctx *RequestContext) error

// Finalizer is a function that runs after request processing is complete.
type Finalizer func(ctx *RequestContext)

// HandlerChain manages the execution of request handlers, exception handlers, response handlers, and finalisers.
type HandlerChain struct {
	requestHandlers   []RequestHandler
	exceptionHandlers []ExceptionHandler
	responseHandlers  []ResponseHandler
	finalizers        []Finalizer
}

// NewHandlerChain creates a new HandlerChain with empty handler lists.
func NewHandlerChain() *HandlerChain {
	return &HandlerChain{
		requestHandlers:   make([]RequestHandler, 0),
		exceptionHandlers: make([]ExceptionHandler, 0),
		responseHandlers:  make([]ResponseHandler, 0),
		finalizers:        make([]Finalizer, 0),
	}
}

// AddRequestHandler adds a request handler to the chain and returns the chain for method chaining.
func (c *HandlerChain) AddRequestHandler(handler RequestHandler) *HandlerChain {
	c.requestHandlers = append(c.requestHandlers, handler)
	return c
}

// AddExceptionHandler adds an exception handler to the chain and returns the chain for method chaining.
func (c *HandlerChain) AddExceptionHandler(handler ExceptionHandler) *HandlerChain {
	c.exceptionHandlers = append(c.exceptionHandlers, handler)
	return c
}

// AddResponseHandler adds a response handler to the chain and returns the chain for method chaining.
func (c *HandlerChain) AddResponseHandler(handler ResponseHandler) *HandlerChain {
	c.responseHandlers = append(c.responseHandlers, handler)
	return c
}

// AddFinalizer adds a finaliser to the chain and returns the chain for method chaining.
func (c *HandlerChain) AddFinalizer(finalizer Finalizer) *HandlerChain {
	c.finalizers = append(c.finalizers, finalizer)
	return c
}

// Execute runs the handler chain for the given request and response.
func (c *HandlerChain) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := NewRequestContext(context.Background(), w, r)

	defer func() {
		for _, finalizer := range c.finalizers {
			finalizer(ctx)
		}
	}()

	for _, handler := range c.requestHandlers {
		if ctx.IsHandled() {
			break
		}
		if err := handler(ctx); err != nil {
			ctx.SetError(err)
			c.handleException(ctx, err)
			return
		}
	}

	if ctx.Error == nil {
		for _, handler := range c.responseHandlers {
			if err := handler(ctx); err != nil {
				ctx.SetError(err)
				c.handleException(ctx, err)
				return
			}
		}
	}
}

func (c *HandlerChain) handleException(ctx *RequestContext, err error) {
	for _, handler := range c.exceptionHandlers {
		if handlerErr := handler(ctx, err); handlerErr != nil {
			ctx.Error = handlerErr
		}
	}
}

// Clone creates a deep copy of the HandlerChain.
func (c *HandlerChain) Clone() *HandlerChain {
	return &HandlerChain{
		requestHandlers:   append([]RequestHandler{}, c.requestHandlers...),
		exceptionHandlers: append([]ExceptionHandler{}, c.exceptionHandlers...),
		responseHandlers:  append([]ResponseHandler{}, c.responseHandlers...),
		finalizers:        append([]Finalizer{}, c.finalizers...),
	}
}
