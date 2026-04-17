package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/services/aws/apigateway/runtime/auth"
	"vorpalstacks/internal/services/aws/apigateway/runtime/integration"
	"vorpalstacks/internal/services/aws/apigateway/runtime/validator"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
	"vorpalstacks/internal/utils/aws/arn"
)

// LambdaInvoker invokes Lambda functions for API Gateway integrations.
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}

// RuntimeServer handles API Gateway runtime requests.
type RuntimeServer struct {
	store            *apigatewaystore.RestApiStore
	usageStore       *apigatewaystore.UsageStore
	router           *RuntimeRouter
	executorFactory  *integration.ExecutorFactory
	lambdaInvoker    LambdaInvoker
	validator        *validator.Validator
	authenticator    *auth.APIKeyAuthenticator
	lambdaAuthorizer *auth.LambdaAuthorizer
	bus              eventbus.Bus
	accountID        string
}

// NewRuntimeServer creates a new API Gateway runtime server.
// Optional SQS/SNS stores for integration targets should be injected via
// setter methods after construction.
func NewRuntimeServer(store *apigatewaystore.RestApiStore, usageStore *apigatewaystore.UsageStore, lambdaInvoker LambdaInvoker) *RuntimeServer {
	return &RuntimeServer{
		store:            store,
		usageStore:       usageStore,
		router:           NewRuntimeRouter(),
		executorFactory:  integration.NewExecutorFactory(lambdaInvoker),
		lambdaInvoker:    lambdaInvoker,
		validator:        validator.NewValidator(),
		authenticator:    auth.NewAPIKeyAuthenticator(usageStore),
		lambdaAuthorizer: auth.NewLambdaAuthorizer(lambdaInvoker, store),
	}
}

// SetSQSStore injects an SQS store into the executor factory for SQS integration targets.
func (s *RuntimeServer) SetSQSStore(store sqsstore.SQSStoreInterface, accountID, region string) {
	s.executorFactory.SetSQSStore(store)
	s.executorFactory.SetAccountAndRegion(accountID, region)
}

// SetSNSStore injects an SNS store into the executor factory for SNS integration targets.
func (s *RuntimeServer) SetSNSStore(store snsstore.SNSStoreInterface, accountID, region string) {
	s.executorFactory.SetSNSStore(store)
	s.executorFactory.SetAccountAndRegion(accountID, region)
}

// SetEventBus injects the event bus for SNS fan-out and cross-service delivery.
func (s *RuntimeServer) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	s.executorFactory.SetEventBus(bus)
}

// SetAccountID stores the AWS account ID used for access log ARN parsing.
func (s *RuntimeServer) SetAccountID(accountID string) {
	s.accountID = accountID
}

// RemoveApiKey cleans up rate limiter state for a deleted API key.
func (s *RuntimeServer) RemoveApiKey(apiKeyId string) {
	if s.authenticator != nil {
		s.authenticator.RemoveApiKey(apiKeyId)
	}
}

// Close stops background goroutines in authentication components and
// waits for in-flight SNS delivery goroutines to finish.
func (s *RuntimeServer) Close() {
	if s.lambdaAuthorizer != nil {
		s.lambdaAuthorizer.Close()
	}
	if s.executorFactory != nil {
		s.executorFactory.WaitDelivery()
	}
}

// HandleRequest handles incoming API Gateway requests.
func (s *RuntimeServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		s.sendError(w, http.StatusBadRequest, "Invalid request path format")
		return
	}

	if pathParts[0] != "restapis" || pathParts[3] != "_user_request_" {
		s.sendError(w, http.StatusBadRequest, "Invalid request path format")
		return
	}

	restApiID := pathParts[1]
	stageName := pathParts[2]
	requestPath := "/" + strings.Join(pathParts[4:], "/")

	restAPI, err := s.store.Get(restApiID)
	if err != nil || restAPI == nil {
		s.sendError(w, http.StatusNotFound, fmt.Sprintf("REST API %s not found", restApiID))
		return
	}

	stage, ok := restAPI.Stages[stageName]
	if !ok || stage == nil {
		s.sendError(w, http.StatusNotFound, fmt.Sprintf("Stage %s not found for API %s", stageName, restApiID))
		return
	}

	_, ok = restAPI.Deployments[stage.DeploymentId]
	if !ok {
		s.sendError(w, http.StatusNotFound, fmt.Sprintf("Deployment %s not found", stage.DeploymentId))
		return
	}

	resourceInfos := make([]ResourceInfo, 0, len(restAPI.Resources))
	for _, res := range restAPI.Resources {
		resourceInfos = append(resourceInfos, ResourceInfo{
			Id:              res.Id,
			Path:            res.Path,
			ParentId:        res.ParentId,
			ResourceMethods: res.ResourceMethods,
		})
	}

	match, err := s.router.MatchWithMethods(resourceInfos, requestPath, r.Method)
	if err != nil {
		s.sendError(w, http.StatusNotFound, err.Error())
		return
	}

	if match == nil || match.Method == nil {
		s.sendError(w, http.StatusNotFound, fmt.Sprintf("No matching route for %s %s", r.Method, requestPath))
		return
	}

	headers := s.getHeaders(r)
	queryParams := s.getQueryParams(r)
	pathParams := match.Params

	authReq := &auth.AuthRequest{
		Method:            r.Method,
		Resource:          match.Path,
		Path:              requestPath,
		Headers:           headers,
		QueryStringParams: queryParams,
		PathParameters:    pathParams,
		StageVariables:    stage.Variables,
		RestApiId:         restApiID,
		StageName:         stageName,
	}

	if match.Method.AuthorizationType != "" && match.Method.AuthorizationType != "NONE" {
		authResult, err := s.lambdaAuthorizer.Authorize(r.Context(), match.Method, restAPI, authReq)
		if err != nil {
			if authErr, ok := err.(*auth.AuthError); ok {
				s.sendError(w, authErr.HTTPCode, authErr.Message)
				return
			}
			s.sendError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if !authResult.Allowed {
			s.sendError(w, http.StatusForbidden, "Access denied")
			return
		}
	}

	if match.Method.ApiKeyRequired {
		apiKeyValue := r.Header.Get("x-api-key")
		if err := s.authenticator.Authenticate(r.Context(), apiKeyValue, match.Method, restApiID, stageName); err != nil {
			if authErr, ok := err.(*auth.AuthError); ok {
				s.sendError(w, authErr.HTTPCode, authErr.Message)
				return
			}
			s.sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024*1024))
	if err != nil {
		s.sendError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read request body: %v", err))
		return
	}

	if err := s.validator.ValidateRequest(match.Method, restAPI, headers, queryParams, pathParams, body); err != nil {
		s.sendError(w, err.HTTPCode, err.Message)
		return
	}

	s.executeIntegration(w, r, match, stage, restApiID, body, headers, queryParams, pathParams)
}

func (s *RuntimeServer) executeIntegration(w http.ResponseWriter, r *http.Request, match *RouteMatch, stage *apigatewaystore.Stage, restApiID string, body []byte, headers, queryParams, pathParams map[string]string) {
	methodIntegration := match.Method.MethodIntegration
	if methodIntegration == nil {
		s.sendError(w, http.StatusNotFound, "No integration configured for this method")
		return
	}

	executor, err := s.executorFactory.CreateExecutor(methodIntegration.Type)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create executor: %v", err))
		return
	}

	integrationReq := &integration.IntegrationRequest{
		Method:                      r.Method,
		URI:                         methodIntegration.Uri,
		Headers:                     headers,
		MultiValueHeaders:           s.getMultiValueHeaders(r),
		Body:                        body,
		QueryParams:                 queryParams,
		MultiValueQueryStringParams: s.getMultiValueQueryParams(r),
		PathParams:                  pathParams,
		Path:                        match.Path,
		StageVariables:              stage.Variables,
		IntegrationType:             methodIntegration.Type,
		RequestParameters:           methodIntegration.RequestParameters,
		RequestTemplates:            methodIntegration.RequestTemplates,
		RestApiId:                   restApiID,
		StageName:                   stage.StageName,
		SourceIP:                    clientIP(r),
		IntegrationResponses:        convertIntegrationResponses(methodIntegration.IntegrationResponses),
	}

	startTime := time.Now()
	integrationResp, err := executor.Execute(r.Context(), integrationReq)
	latency := time.Since(startTime)
	statusCode := http.StatusInternalServerError
	if integrationResp != nil {
		statusCode = integrationResp.StatusCode
	}

	if err != nil {
		if intErr, ok := err.(*integration.IntegrationError); ok {
			statusCode = intErr.HTTPCode
			s.sendError(w, intErr.HTTPCode, intErr.Message)
		} else {
			s.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Integration execution failed: %v", err))
		}
		s.writeAccessLog(r, stage, restApiID, match.Path, statusCode, latency)
		return
	}

	s.sendResponse(w, integrationResp)
	s.writeAccessLog(r, stage, restApiID, match.Path, statusCode, latency)
}

func (s *RuntimeServer) getHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = strings.Join(values, ",")
		}
	}
	return headers
}

func (s *RuntimeServer) getQueryParams(r *http.Request) map[string]string {
	params := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params
}

func (s *RuntimeServer) getMultiValueHeaders(r *http.Request) map[string][]string {
	headers := make(map[string][]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values
		}
	}
	return headers
}

func (s *RuntimeServer) getMultiValueQueryParams(r *http.Request) map[string][]string {
	params := make(map[string][]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values
		}
	}
	return params
}

func convertIntegrationResponses(responses map[string]*apigatewaystore.IntegrationResponse) map[string]*integration.IntegrationResponseConfig {
	if responses == nil {
		return nil
	}
	result := make(map[string]*integration.IntegrationResponseConfig)
	for k, v := range responses {
		result[k] = &integration.IntegrationResponseConfig{
			StatusCode:        v.StatusCode,
			SelectionPattern:  v.SelectionPattern,
			ResponseHeaders:   v.ResponseParameters,
			ResponseTemplates: v.ResponseTemplates,
			ContentHandling:   v.ContentHandling,
		}
	}
	return result
}

func (s *RuntimeServer) sendResponse(w http.ResponseWriter, resp *integration.IntegrationResponse) {
	for key, value := range resp.Headers {
		w.Header().Set(key, value)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(resp.StatusCode)

	if len(resp.Body) > 0 {
		if _, err := w.Write(resp.Body); err != nil {
			logs.Error("Failed to write response body", logs.Err(err))
		}
	}
}

func (s *RuntimeServer) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := map[string]interface{}{
		"message": message,
	}

	switch statusCode {
	case http.StatusBadRequest:
		errorResp["__type"] = "BadRequestException"
	case http.StatusUnauthorized:
		errorResp["__type"] = "UnauthorizedException"
	case http.StatusForbidden:
		errorResp["__type"] = "ForbiddenException"
	case http.StatusNotFound:
		errorResp["__type"] = "NotFoundException"
	case http.StatusTooManyRequests:
		errorResp["__type"] = "TooManyRequestsException"
	case http.StatusInternalServerError:
		errorResp["__type"] = "InternalServerError"
	case http.StatusBadGateway:
		errorResp["__type"] = "BadGatewayException"
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		logs.Error("Failed to encode error response", logs.Err(err))
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx >= 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (s *RuntimeServer) writeAccessLog(r *http.Request, stage *apigatewaystore.Stage, restApiID, resourcePath string, statusCode int, latency time.Duration) {
	if stage == nil || stage.AccessLogSettings == nil || stage.AccessLogSettings.DestinationArn == "" {
		return
	}

	if s.bus == nil {
		return
	}

	_, _, region, _, _ := arn.SplitARN(stage.AccessLogSettings.DestinationArn)
	logGroup := arn.ExtractLogGroupNameFromARN(stage.AccessLogSettings.DestinationArn)
	if logGroup == "" || region == "" {
		return
	}

	logStream := fmt.Sprintf("%s/%s", restApiID, stage.StageName)

	formatted := s.formatAccessLog(stage.AccessLogSettings.Format, r, restApiID, stage.StageName, resourcePath, statusCode, latency)

	evt := &eventbus.APIGatewayAccessLogEvent{
		RestApiId:      restApiID,
		StageName:      stage.StageName,
		DestinationArn: stage.AccessLogSettings.DestinationArn,
		LogGroup:       logGroup,
		LogStream:      logStream,
		FormattedLog:   formatted,
	}
	evt.Region = region
	evt.AccountID = s.accountID

	if pubErr := s.bus.Publish(context.Background(), evt); pubErr != nil {
		logs.Warn("failed to publish API Gateway access log event", logs.Err(pubErr))
	}
}

func (s *RuntimeServer) formatAccessLog(format string, r *http.Request, restApiID, stageName, resourcePath string, statusCode int, latency time.Duration) string {
	if format == "" {
		return fmt.Sprintf("%s %s %s %d %dms", r.Method, r.URL.Path, r.Proto, statusCode, latency.Milliseconds())
	}

	requestID := r.Header.Get("X-Amzn-RequestId")
	if requestID == "" {
		requestID = r.Header.Get("X-Amz-Request-Id")
	}
	if requestID == "" {
		requestID = r.Header.Get("X-Amzn-Trace-Id")
	}
	if requestID == "" {
		requestID = "req-unknown"
	}

	result := format
	result = strings.ReplaceAll(result, "$context.requestId", requestID)
	result = strings.ReplaceAll(result, "$context.requestTime", time.Now().UTC().Format("02/Jan/2006:15:04:05 +0000"))
	result = strings.ReplaceAll(result, "$context.httpMethod", r.Method)
	result = strings.ReplaceAll(result, "$context.resourcePath", resourcePath)
	result = strings.ReplaceAll(result, "$context.path", r.URL.Path)
	result = strings.ReplaceAll(result, "$context.protocol", r.Proto)
	result = strings.ReplaceAll(result, "$context.status", fmt.Sprintf("%d", statusCode))
	result = strings.ReplaceAll(result, "$context.responseLatency", fmt.Sprintf("%d", latency.Milliseconds()))
	result = strings.ReplaceAll(result, "$context.integrationLatency", fmt.Sprintf("%d", latency.Milliseconds()))
	result = strings.ReplaceAll(result, "$context.sourceIp", clientIP(r))
	result = strings.ReplaceAll(result, "$context.accountId", s.accountID)
	result = strings.ReplaceAll(result, "$context.apiId", restApiID)
	result = strings.ReplaceAll(result, "$context.stage", stageName)
	result = strings.ReplaceAll(result, "$context.error.message", "")
	result = strings.ReplaceAll(result, "$context.error.messageString", "")

	if strings.Contains(result, "$context.request.header.") {
		for key, values := range r.Header {
			varName := "$context.request.header." + strings.ToLower(strings.ReplaceAll(key, "-", ""))
			result = strings.ReplaceAll(result, varName, strings.Join(values, ","))
		}
	}

	return result
}
