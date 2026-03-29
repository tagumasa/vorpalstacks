package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/storage"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

type FunctionURLServer struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	region         string
	lambdaInvoker  LambdaInvoker
}

type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}

func NewFunctionURLServer(storageManager *storage.RegionStorageManager, accountID, region string, invoker LambdaInvoker) *FunctionURLServer {
	return &FunctionURLServer{
		storageManager: storageManager,
		accountID:      accountID,
		region:         region,
		lambdaInvoker:  invoker,
	}
}

func (s *FunctionURLServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	functionName := s.extractFunctionName(r.Host)
	if functionName == "" {
		http.Error(w, `{"message":"Function not found"}`, http.StatusNotFound)
		return
	}

	store := s.getFunctionStore()
	function, err := store.Get(functionName)
	if err != nil || function == nil {
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

	body, err := io.ReadAll(r.Body)
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
		log.Printf("Function URL invocation failed for %s: %v", functionName, err)
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
	w.Write(payload)
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

func (s *FunctionURLServer) getFunctionStore() *lambdastore.FunctionStore {
	store, err := s.storageManager.GetStorage(s.region)
	if err != nil {
		return lambdastore.NewFunctionStore(nil, s.accountID, s.region)
	}
	return lambdastore.NewFunctionStore(store, s.accountID, s.region)
}
