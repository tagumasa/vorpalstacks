package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// FunctionStoreProvider abstracts access to Lambda function stores by region.
type FunctionStoreProvider interface {
	GetFunctionStoreForRegion(region string) *lambdastore.FunctionStore
}

// FunctionURLServer handles incoming requests to Lambda function URLs,
// invoking the target function and proxying the response.
type FunctionURLServer struct {
	accountID     string
	region        string
	storeProvider FunctionStoreProvider
	lambdaInvoker LambdaInvoker
}

// LambdaInvoker abstracts the ability to invoke a Lambda function for gateway use.
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}

// NewFunctionURLServer creates a new FunctionURLServer backed by the shared
// Lambda service store cache.
func NewFunctionURLServer(storeProvider FunctionStoreProvider, accountID, region string, invoker LambdaInvoker) *FunctionURLServer {
	return &FunctionURLServer{
		storeProvider: storeProvider,
		accountID:     accountID,
		region:        region,
		lambdaInvoker: invoker,
	}
}

// HandleRequest routes an incoming request to the target Lambda function via its function URL configuration.
func (s *FunctionURLServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	functionName := s.extractFunctionName(r.Host)
	if functionName == "" {
		http.Error(w, `{"message":"Function not found"}`, http.StatusNotFound)
		return
	}

	function, _ := s.findFunction(functionName)
	if function == nil {
		http.Error(w, `{"message":"Function not found"}`, http.StatusNotFound)
		return
	}

	urlConfig := function.UrlConfig
	if urlConfig == nil {
		http.Error(w, `{"message":"Function URL not configured"}`, http.StatusNotFound)
		return
	}

	if urlConfig.AuthType == "AWS_IAM" {
		http.Error(w, `{"message":"AWS_IAM auth not supported in this mode"}`, http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodOptions {
		s.handleCORS(w, r, urlConfig)
		return
	}

	s.setCORSHeaders(w, r, urlConfig)

	body, err := io.ReadAll(io.LimitReader(r.Body, 6*1024*1024))
	if err != nil {
		http.Error(w, `{"message":"Failed to read request body"}`, http.StatusBadRequest)
		return
	}

	event := s.buildFunctionURLEvent(r, body, functionName, urlConfig)

	eventJSON, err := json.Marshal(event)
	if err != nil {
		http.Error(w, `{"message":"Failed to marshal event"}`, http.StatusInternalServerError)
		return
	}

	statusCode, payload, err := s.lambdaInvoker.InvokeForGateway(r.Context(), functionName, eventJSON)
	if err != nil {
		logs.Error("Function URL invocation failed", logs.String("function", functionName), logs.Err(err))
		http.Error(w, fmt.Sprintf(`{"message":"Invocation failed: %s"}`, err.Error()), http.StatusBadGateway)
		return
	}

	if payload != nil {
		var resp struct {
			StatusCode        int                 `json:"statusCode"`
			Headers           map[string]string   `json:"headers"`
			MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
			Body              string              `json:"body"`
			IsBase64Encoded   bool                `json:"isBase64Encoded"`
		}
		if err := json.Unmarshal(payload, &resp); err == nil && resp.StatusCode > 0 {
			for k, v := range resp.Headers {
				w.Header().Set(k, v)
			}
			for k, vals := range resp.MultiValueHeaders {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(resp.StatusCode)
			fmt.Fprint(w, resp.Body)
			return
		}
	}

	code := int(statusCode)
	if code == 0 {
		code = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(payload)
}

func (s *FunctionURLServer) buildFunctionURLEvent(r *http.Request, body []byte, functionName string, urlConfig *lambdastore.FunctionUrlConfig) map[string]interface{} {
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	rawQuery := r.URL.RawQuery
	path := r.URL.Path

	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return map[string]interface{}{
		"version":               "2.0",
		"rawQueryString":        rawQuery,
		"queryStringParameters": queryParams,
		"headers":               headers,
		"requestContext": map[string]interface{}{
			"accountId":    s.accountID,
			"apiId":        "",
			"domainName":   r.Host,
			"domainPrefix": strings.Split(r.Host, ".")[0],
			"http": map[string]interface{}{
				"method":    r.Method,
				"path":      path,
				"protocol":  r.Proto,
				"sourceIp":  r.RemoteAddr,
				"userAgent": r.UserAgent(),
			},
			"requestId": uuid.New().String(),
			"routeKey":  "$default",
			"stage":     "$default",
			"time":      r.Header.Get("Date"),
			"timeEpoch": 0,
		},
		"body":            string(body),
		"pathParameters":  map[string]string{},
		"isBase64Encoded": false,
		"stageVariables":  map[string]string{},
	}
}

func (s *FunctionURLServer) extractFunctionName(host string) string {
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if len(parts) >= 1 {
		candidate := parts[0]
		if candidate != "" && candidate != "localhost" {
			return candidate
		}
	}
	return ""
}

func (s *FunctionURLServer) handleCORS(w http.ResponseWriter, r *http.Request, urlConfig *lambdastore.FunctionUrlConfig) {
	s.setCORSHeaders(w, r, urlConfig)
	w.WriteHeader(http.StatusNoContent)
}

func (s *FunctionURLServer) setCORSHeaders(w http.ResponseWriter, r *http.Request, urlConfig *lambdastore.FunctionUrlConfig) {
	if urlConfig.Cors == nil {
		return
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	allowed := false
	for _, o := range urlConfig.Cors.AllowOrigins {
		if o == "*" || o == origin {
			allowed = true
			break
		}
	}
	if !allowed {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	if len(urlConfig.Cors.AllowMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(urlConfig.Cors.AllowMethods, ", "))
	}
	if len(urlConfig.Cors.AllowHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(urlConfig.Cors.AllowHeaders, ", "))
	}
	if len(urlConfig.Cors.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(urlConfig.Cors.ExposeHeaders, ", "))
	}
	if urlConfig.Cors.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if urlConfig.Cors.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", urlConfig.Cors.MaxAge))
	}
}

// findFunction looks up a function by name in the service region.
func (s *FunctionURLServer) findFunction(functionName string) (*lambdastore.Function, string) {
	fs := s.storeProvider.GetFunctionStoreForRegion(s.region)
	if fs == nil {
		return nil, ""
	}
	fn, err := fs.Get(functionName)
	if err == nil && fn != nil {
		return fn, s.region
	}
	return nil, ""
}
