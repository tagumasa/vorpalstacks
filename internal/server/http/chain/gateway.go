package chain

import (
	"net/http"
	"strings"

	"vorpalstacks/internal/server/actionregistry"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/http/router"
)

// Gateway routes HTTP requests to the appropriate service and operation based on the request content.
type Gateway struct {
	chain      *HandlerChain
	dispatcher *dispatcher.Dispatcher
	router     *router.ServiceRouter
}

// NewGateway creates a new Gateway with the given service router and dispatcher.
func NewGateway(serviceRouter *router.ServiceRouter, disp *dispatcher.Dispatcher) *Gateway {
	g := &Gateway{
		dispatcher: disp,
		router:     serviceRouter,
	}

	g.chain = g.buildChain()
	return g
}

func (g *Gateway) buildChain() *HandlerChain {
	c := NewHandlerChain()

	c.AddRequestHandler(AddRequestContext)
	c.AddRequestHandler(ParseServiceName(g.router))
	c.AddRequestHandler(ParseOperationName(g.router))
	c.AddRequestHandler(ParseRequestParams)
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

func (g *Gateway) routeRequest(ctx *RequestContext) error {
	if ctx.ServiceName == "" {
		ctx.ServiceName = g.fallbackServiceName(ctx.Request)
	}

	if ctx.ServiceName == "" {
		ctx.StatusCode = http.StatusNotFound
		ctx.ResponseBody = map[string]interface{}{
			"__type":  "UnknownService",
			"message": "Unable to determine service from request",
		}
		ctx.SetHandled()
		return nil
	}

	if ctx.Operation == "" {
		ctx.Operation = g.fallbackOperation(ctx.Request, ctx.ServiceName)
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

func (g *Gateway) fallbackServiceName(r *http.Request) string {
	if router.IsLambdaRestPath(r.URL.Path) {
		return "lambda"
	}

	if isApiGatewayPath(r.URL.Path) {
		return "apigateway"
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.Contains(authHeader, "Credential=") {
			parts := strings.Split(authHeader, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if idx := strings.Index(part, "Credential="); idx >= 0 {
					cred := part[idx+11:]
					credParts := strings.Split(cred, "/")
					if len(credParts) >= 4 {
						return credParts[3]
					}
				}
			}
		}
	}

	xAmzTarget := r.Header.Get("X-Amz-Target")
	if xAmzTarget != "" {
		parts := strings.Split(xAmzTarget, ".")
		if len(parts) > 0 {
			target := parts[0]
			switch target {
			case "AmazonSSM":
				return "ssm"
			case "AWSCognitoIdentityProviderService":
				return "cognito-idp"
			case "AWSCognitoIdentityService":
				return "cognito-identity"
			default:
				return strings.ToLower(target)
			}
		}
	}

	action := r.URL.Query().Get("Action")
	if action != "" {
		return g.inferServiceFromAction(action)
	}

	if r.Method == "POST" && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		if actionParam := ctxParams(r, "Action"); actionParam != "" {
			return g.inferServiceFromAction(actionParam)
		}
	}

	return ""
}

func (g *Gateway) fallbackOperation(r *http.Request, serviceName string) string {
	xAmzTarget := r.Header.Get("X-Amz-Target")
	if xAmzTarget != "" {
		parts := strings.Split(xAmzTarget, ".")
		if len(parts) > 1 {
			return parts[1]
		}
		return xAmzTarget
	}

	action := r.URL.Query().Get("Action")
	if action != "" {
		return action
	}

	if r.Method == "POST" && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		if actionParam := ctxParams(r, "Action"); actionParam != "" {
			return actionParam
		}
	}

	return ""
}

func (g *Gateway) inferServiceFromAction(action string) string {
	return actionregistry.LookupServiceByAction(action)
}

func ctxParams(r *http.Request, key string) string {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			return ""
		}
	}
	if vals, ok := r.Form[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (g *Gateway) finalizer(ctx *RequestContext) {
}

func isApiGatewayPath(path string) bool {
	return strings.HasPrefix(path, "/restapis") ||
		strings.HasPrefix(path, "/apikeys") ||
		strings.HasPrefix(path, "/usageplans") ||
		strings.HasPrefix(path, "/domainnames") ||
		strings.HasPrefix(path, "/vpclinks") ||
		strings.HasPrefix(path, "/apis") ||
		strings.HasPrefix(path, "/authorizers")
}

// ServeHTTP implements the http.Handler interface to handle HTTP requests.
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.chain.Execute(w, r)
}

// GetChain returns the handler chain used by the gateway.
func (g *Gateway) GetChain() *HandlerChain {
	return g.chain
}

// AddRequestHandler adds a request handler to the gateway's handler chain.
func (g *Gateway) AddRequestHandler(handler RequestHandler) {
	g.chain.AddRequestHandler(handler)
}

// AddExceptionHandler adds an exception handler to the gateway's handler chain.
func (g *Gateway) AddExceptionHandler(handler ExceptionHandler) {
	g.chain.AddExceptionHandler(handler)
}

// AddResponseHandler adds a response handler to the gateway's handler chain.
func (g *Gateway) AddResponseHandler(handler ResponseHandler) {
	g.chain.AddResponseHandler(handler)
}

// AddFinalizer adds a finaliser to the gateway's handler chain.
func (g *Gateway) AddFinalizer(finalizer Finalizer) {
	g.chain.AddFinalizer(finalizer)
}
