package chain

import (
	"net/http"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/http/classifier"
)

// Gateway is the chain-based HTTP request gateway that classifies incoming
// requests and routes them to the dispatcher.
type Gateway struct {
	chain      *HandlerChain
	dispatcher *dispatcher.Dispatcher
	classifier *classifier.Classifier
}

// NewGateway creates a new Gateway with the given classifier and dispatcher.
func NewGateway(c *classifier.Classifier, disp *dispatcher.Dispatcher) *Gateway {
	g := &Gateway{
		dispatcher: disp,
		classifier: c,
	}

	g.chain = g.buildChain()
	return g
}

// buildChain constructs the default handler chain with request, exception,
// response, and finaliser handlers.
func (g *Gateway) buildChain() *HandlerChain {
	c := NewHandlerChain()

	c.AddRequestHandler(AddRequestContext)
	c.AddRequestHandler(g.parseRequest)
	c.AddRequestHandler(AddRegionFromHeader)
	c.AddRequestHandler(EnforceCORS)
	c.AddRequestHandler(LogRequest)
	c.AddRequestHandler(g.routeRequest)

	c.AddExceptionHandler(LogException)
	c.AddExceptionHandler(HandleServiceError)
	c.AddExceptionHandler(HandleInternalFailure)

	c.AddResponseHandler(AddCORSResponseHeaders)
	c.AddResponseHandler(SerializeResponse)
	c.AddResponseHandler(LogResponse)

	c.AddFinalizer(SetCloseConnectionHeader)
	c.AddFinalizer(g.finalizer)

	return c
}

// parseRequest uses the classifier to extract service name, operation, and
// parameters from the incoming request, populating the RequestContext.
func (g *Gateway) parseRequest(ctx *RequestContext) error {
	if ctx.Request == nil {
		return nil
	}

	cr, err := g.classifier.Classify(ctx.Request)
	if err != nil || cr == nil {
		return nil
	}

	ctx.Classified = cr
	ctx.ServiceName = cr.ServiceName
	ctx.Operation = cr.Operation
	if ctx.Params == nil {
		ctx.Params = make(map[string]interface{})
	}
	for k, v := range cr.Parameters {
		ctx.Params[k] = v
	}

	return nil
}

// routeRequest dispatches the classified request to the appropriate handler
// via the dispatcher. If no service was identified, it returns a 404 error.
func (g *Gateway) routeRequest(ctx *RequestContext) error {
	if ctx.ServiceName == "" {
		ctx.StatusCode = http.StatusNotFound
		ctx.ResponseBody = map[string]interface{}{
			"__type":  "UnknownService",
			"message": "Unable to determine service from request",
		}
		ctx.SetHandled()
		return nil
	}

	if ctx.Classified != nil {
		g.dispatcher.DispatchClassified(ctx.Response, ctx.Request, ctx.Classified)
		ctx.SetHandled()
		return nil
	}

	dispatchCtx := &dispatcher.DispatchContext{
		ServiceName: ctx.ServiceName,
		Operation:   ctx.Operation,
		Params:      ctx.Params,
	}

	g.dispatcher.DispatchWithContext(ctx.Response, ctx.Request, dispatchCtx)
	ctx.SetHandled()

	return nil
}

// finalizer is a no-op chain finaliser reserved for future use.
func (g *Gateway) finalizer(ctx *RequestContext) {
}

// ServeHTTP implements http.Handler, executing the full handler chain for
// each incoming request.
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.chain.Execute(w, r)
}

// GetChain returns the underlying handler chain for inspection or extension.
func (g *Gateway) GetChain() *HandlerChain {
	return g.chain
}

// AddRequestHandler appends a request handler to the chain.
func (g *Gateway) AddRequestHandler(handler RequestHandler) {
	g.chain.AddRequestHandler(handler)
}

// AddExceptionHandler appends an exception handler to the chain.
func (g *Gateway) AddExceptionHandler(handler ExceptionHandler) {
	g.chain.AddExceptionHandler(handler)
}

// AddResponseHandler appends a response handler to the chain.
func (g *Gateway) AddResponseHandler(handler ResponseHandler) {
	g.chain.AddResponseHandler(handler)
}

// AddFinalizer appends a finaliser to the chain.
func (g *Gateway) AddFinalizer(finalizer Finalizer) {
	g.chain.AddFinalizer(finalizer)
}
