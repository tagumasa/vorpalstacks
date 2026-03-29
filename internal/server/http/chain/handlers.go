package chain

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/server/http/router"
	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/utils/buffer"
)

// AddRequestContext adds the request context to the request context object.
func AddRequestContext(ctx *RequestContext) error {
	if ctx.Request == nil {
		return fmt.Errorf("nil request")
	}
	return nil
}

// ParseServiceName returns a RequestHandler that parses the service name from the request.
func ParseServiceName(serviceRouter *router.ServiceRouter) RequestHandler {
	return func(ctx *RequestContext) error {
		if serviceRouter == nil {
			return nil
		}

		serviceName := serviceRouter.DetermineService(ctx.Request)
		if serviceName != "" {
			ctx.ServiceName = serviceName
			logs.Debug("Parsed service name", logs.String("service", serviceName))
		}

		return nil
	}
}

// ParseOperationName returns a RequestHandler that parses the operation name from the request.
func ParseOperationName(serviceRouter *router.ServiceRouter) RequestHandler {
	return func(ctx *RequestContext) error {
		if serviceRouter == nil || ctx.ServiceName == "" {
			return nil
		}

		opName := serviceRouter.DetermineOperation(ctx.Request, ctx.ServiceName)
		if opName != "" {
			ctx.Operation = opName
			logs.Debug("Parsed operation", logs.String("service", ctx.ServiceName), logs.String("operation", opName))
		}

		return nil
	}
}

// ParseRequestParams parses the request parameters from the request body or query string.
func ParseRequestParams(ctx *RequestContext) error {
	if ctx.Request == nil {
		return nil
	}

	req := ctx.Request
	contentType := req.Header.Get("Content-Type")

	if req.Method == http.MethodGet {
		parseQueryParams(req, ctx.Params)
		return nil
	}

	if strings.Contains(contentType, "application/x-amz-json") {
		return parseJSONBody(req, ctx.Params)
	}

	if strings.Contains(contentType, "application/json") {
		return parseJSONBody(req, ctx.Params)
	}

	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return parseFormData(req, ctx.Params)
	}

	return nil
}

func parseQueryParams(req *http.Request, params map[string]interface{}) {
	for key, values := range req.URL.Query() {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}
}

func parseJSONBody(req *http.Request, params map[string]interface{}) error {
	if req.Body == nil {
		return nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	if len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, &params); err != nil {
		return fmt.Errorf("failed to parse JSON body: %w", err)
	}

	return nil
}

func parseFormData(req *http.Request, params map[string]interface{}) error {
	if req.Body == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	for key, values := range req.Form {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}

	for key, values := range req.PostForm {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}

	req.Form = nil
	req.PostForm = nil
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return nil
}

// AddRegionFromHeader extracts the region from the Authorization header and adds it to the request params.
func AddRegionFromHeader(ctx *RequestContext) error {
	if ctx.Request == nil {
		return nil
	}

	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	if strings.Contains(authHeader, "Credential=") {
		parts := strings.Split(authHeader, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if idx := strings.Index(part, "Credential="); idx >= 0 {
				cred := part[idx+11:]
				credParts := strings.Split(cred, "/")
				if len(credParts) >= 3 {
					if ctx.Params == nil {
						ctx.Params = make(map[string]interface{})
					}
					ctx.Params["Region"] = credParts[2]
				}
			}
		}
	}

	return nil
}

// EnforceCORS adds CORS headers to the response and handles OPTIONS preflight requests.
func EnforceCORS(ctx *RequestContext) error {
	if ctx.Request == nil {
		return nil
	}

	origin := ctx.Request.Header.Get("Origin")
	if origin == "" {
		return nil
	}

	ctx.Response.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Amz-Target, X-Amz-Date, X-Amz-Content-Sha256")

	if ctx.Request.Method == http.MethodOptions {
		ctx.SetHandled()
		ctx.StatusCode = http.StatusOK
	}

	return nil
}

// LogRequest logs the incoming request details.
func LogRequest(ctx *RequestContext) error {
	if ctx.Request == nil {
		return nil
	}

	logs.Info("Request received",
		logs.String("method", ctx.Request.Method),
		logs.String("path", ctx.Request.URL.Path),
		logs.String("service", ctx.ServiceName),
		logs.String("operation", ctx.Operation),
	)

	return nil
}

// LogException logs exceptions that occur during request processing.
func LogException(ctx *RequestContext, err error) error {
	if err == nil {
		return nil
	}

	logs.Error("Request error",
		logs.Err(err),
		logs.String("service", ctx.ServiceName),
		logs.String("operation", ctx.Operation),
		logs.String("path", ctx.Request.URL.Path),
	)

	return nil
}

// HandleServiceError handles service-level errors and sets the appropriate HTTP status code and response body.
func HandleServiceError(ctx *RequestContext, err error) error {
	if err == nil {
		return nil
	}

	statusCode := http.StatusInternalServerError
	errorBody := map[string]interface{}{
		"__type":  "InternalError",
		"message": "An internal error occurred",
	}

	if awsErr, ok := err.(interface {
		Code() string
		Message() string
		HTTPStatusCode() int
	}); ok {
		statusCode = awsErr.HTTPStatusCode()
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
		}
		errorBody["__type"] = awsErr.Code()
		errorBody["message"] = awsErr.Message()
	} else if awsErr, ok := err.(interface {
		GetCode() string
		GetMessage() string
		GetHTTPStatusCode() int
	}); ok {
		statusCode = awsErr.GetHTTPStatusCode()
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
		}
		errorBody["__type"] = awsErr.GetCode()
		errorBody["message"] = awsErr.GetMessage()
	}

	ctx.StatusCode = statusCode
	ctx.ResponseBody = errorBody

	writeJSONResponse(ctx.Response, statusCode, errorBody)
	ctx.SetHandled()

	return nil
}

// HandleInternalFailure handles unhandled errors and returns a generic internal failure response.
func HandleInternalFailure(ctx *RequestContext, err error) error {
	if err == nil {
		return nil
	}

	if ctx.IsHandled() {
		return nil
	}

	statusCode := http.StatusInternalServerError
	errorBody := map[string]interface{}{
		"__type":  "InternalFailure",
		"message": "An internal error occurred",
	}

	ctx.StatusCode = statusCode
	ctx.ResponseBody = errorBody

	writeJSONResponse(ctx.Response, statusCode, errorBody)
	ctx.SetHandled()

	return nil
}

// SerializeResponse serialises the response body to the appropriate format (JSON or XML).
func SerializeResponse(ctx *RequestContext) error {
	if ctx.ResponseBody == nil {
		return nil
	}

	contentType := determineContentType(ctx)

	switch contentType {
	case "application/xml", "text/xml":
		return serializeXML(ctx)
	default:
		return serializeJSON(ctx)
	}
}

func serializeJSON(ctx *RequestContext) error {
	ctx.Response.Header().Set("Content-Type", "application/json")

	if ctx.StatusCode == 0 {
		ctx.StatusCode = http.StatusOK
	}

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(protocol.ConvertTimestampsToSeconds(ctx.ResponseBody)); err != nil {
		return err
	}

	ctx.Response.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	ctx.Response.WriteHeader(ctx.StatusCode)
	_, err := ctx.Response.Write(buf.Bytes())
	return err
}

func serializeXML(ctx *RequestContext) error {
	ctx.Response.Header().Set("Content-Type", "application/xml")

	if ctx.StatusCode == 0 {
		ctx.StatusCode = http.StatusOK
	}

	if ctx.ResponseBody == nil {
		ctx.Response.Header().Set("Content-Length", "0")
		ctx.Response.WriteHeader(ctx.StatusCode)
		return nil
	}

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	switch v := ctx.ResponseBody.(type) {
	case string:
		buf.WriteString(v)
	case []byte:
		buf.Write(v)
	default:
		data, err := xml.Marshal(ctx.ResponseBody)
		if err != nil {
			return err
		}
		buf.Write(data)
	}

	ctx.Response.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	ctx.Response.WriteHeader(ctx.StatusCode)
	_, err := ctx.Response.Write(buf.Bytes())
	return err
}

// AddCORSResponseHeaders adds CORS headers to the response for cross-origin requests.
func AddCORSResponseHeaders(ctx *RequestContext) error {
	origin := ctx.Request.Header.Get("Origin")
	if origin == "" {
		return nil
	}

	ctx.Response.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Amz-Target, X-Amz-Date, X-Amz-Content-Sha256")
	ctx.Response.Header().Set("Access-Control-Expose-Headers", "x-amzn-RequestId, x-amzn-ErrorType, x-amzn-ErrorMessage")

	return nil
}

// LogResponse logs the response details.
func LogResponse(ctx *RequestContext) error {
	logs.Debug("Response sent",
		logs.Int("status", ctx.StatusCode),
		logs.String("service", ctx.ServiceName),
		logs.String("operation", ctx.Operation),
	)

	return nil
}

// SetCloseConnectionHeader adds a Connection: close header to the response if requested by the client.
func SetCloseConnectionHeader(ctx *RequestContext) {
	if ctx.Request == nil {
		return
	}

	connection := ctx.Request.Header.Get("Connection")
	if strings.ToLower(connection) == "close" {
		ctx.Response.Header().Set("Connection", "close")
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(statusCode)
	if _, err := w.Write(buf.Bytes()); err != nil {
		logs.Error("Failed to write response", logs.Err(err))
	}
}

func determineContentType(ctx *RequestContext) string {
	acceptHeader := ctx.Request.Header.Get("Accept")
	if acceptHeader != "" && !strings.Contains(acceptHeader, "*/*") {
		return acceptHeader
	}

	contentType := ctx.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "xml") {
		return "application/xml"
	}

	return "application/json"
}

func extractHost(req *http.Request) string {
	host := req.Host
	if host == "" {
		host = req.Header.Get("Host")
	}

	if i := strings.LastIndex(host, ":"); i >= 0 && !strings.Contains(host[i:], "]") {
		host = host[:i]
	}

	return host
}

func extractPath(req *http.Request) string {
	path := req.URL.Path
	if rawPath := req.URL.RawPath; rawPath != "" {
		path = rawPath
	}
	return path
}

func extractQuery(req *http.Request) url.Values {
	return req.URL.Query()
}
