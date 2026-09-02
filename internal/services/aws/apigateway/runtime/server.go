package runtime

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/headerorder"
	"vorpalstacks/internal/common/waflimits"
	"vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/fqdnrouter"
	"vorpalstacks/internal/services/aws/apigateway/runtime/auth"
	"vorpalstacks/internal/services/aws/apigateway/runtime/integration"
	"vorpalstacks/internal/services/aws/apigateway/runtime/validator"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/aws/arn"
)

// RuntimeServer handles API Gateway runtime requests.
type RuntimeServer struct {
	store            *apigatewaystore.RestApiStore
	usageStore       *apigatewaystore.UsageStore
	router           *RuntimeRouter
	executorFactory  *integration.ExecutorFactory
	validator        *validator.Validator
	authenticator    *auth.APIKeyAuthenticator
	lambdaAuthorizer *auth.LambdaAuthorizer
	bus              eventbus.Bus
	accountID        string
	region           string
	stageThrottlers  sync.Map
	inspectorMu      sync.RWMutex
	webACLInspector  eventbus.WebACLInspector
}

// NewRuntimeServer creates a new API Gateway runtime server.
func NewRuntimeServer(store *apigatewaystore.RestApiStore, usageStore *apigatewaystore.UsageStore, bus eventbus.Bus) *RuntimeServer {
	return &RuntimeServer{
		store:            store,
		usageStore:       usageStore,
		router:           NewRuntimeRouter(),
		executorFactory:  integration.NewExecutorFactory(bus),
		validator:        validator.NewValidator(),
		authenticator:    auth.NewAPIKeyAuthenticator(usageStore),
		lambdaAuthorizer: auth.NewLambdaAuthorizer(bus, store),
	}
}

// SetEventBus injects the event bus for cross-service delivery.
func (s *RuntimeServer) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	s.executorFactory.SetEventBus(bus)
}

// SetAccountID stores the AWS account ID used for access log ARN parsing.
func (s *RuntimeServer) SetAccountID(accountID string) {
	s.accountID = accountID
}

// SetRegion stores the region the runtime serves, used to address the
// stage ARNs that WAFv2 WebACL associations resolve.
func (s *RuntimeServer) SetRegion(region string) {
	s.region = region
}

// SetWebACLInspector injects the WAF request-inspection entry point so
// WebACLs associated with the served stages are enforced on requests.
func (s *RuntimeServer) SetWebACLInspector(inspector eventbus.WebACLInspector) {
	s.inspectorMu.Lock()
	s.webACLInspector = inspector
	s.inspectorMu.Unlock()
}

// enforceWebACL inspects the request against the WebACL associated with
// the stage and answers blocked requests itself, returning false.
// Inspection failures fail closed: the AWS WAF Developer Guide
// documents that when WAF encounters an internal error, Regional
// services typically deny the request and don't serve the content, so
// the failure is answered with 500 Internal Server Error (the response
// the Application Load Balancer integration documents for an
// unreachable WAF). The body is buffered up to the inspection limit and
// recombined with the unread tail so integration execution still reads
// the full request body.
func (s *RuntimeServer) enforceWebACL(w http.ResponseWriter, r *http.Request, stageARN string) bool {
	s.inspectorMu.RLock()
	inspector := s.webACLInspector
	s.inspectorMu.RUnlock()
	if inspector == nil {
		return true
	}

	var body []byte
	bodyTruncated := false
	if r.Body != nil {
		buffered, err := io.ReadAll(io.LimitReader(r.Body, waflimits.DefaultBodyInspectionLimit+1))
		if err != nil {
			logs.Warn("apigateway waf inspection body read failed, denying request", logs.String("stage", stageARN), logs.Err(err))
			http.Error(w, "The request could not be inspected", http.StatusInternalServerError)
			return false
		}
		body = buffered
		if int64(len(body)) > waflimits.DefaultBodyInspectionLimit {
			body = body[:waflimits.DefaultBodyInspectionLimit]
			bodyTruncated = true
			r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buffered), r.Body))
		} else {
			r.Body = io.NopCloser(bytes.NewReader(buffered))
		}
	}

	inspHeaders := eventbus.RequestHeadersWithHost(r.Header, r.Host)
	headerOrder, _ := headerorder.FromContext(r.Context(), inspHeaders)
	result, err := inspector.InspectWebACLRequest(r.Context(), s.region, stageARN, eventbus.BuildWebACLInspectionRequest(
		r.Method, r.URL.Path, r.URL.RawQuery, remoteAddrHost(r.RemoteAddr), r.Proto,
		inspHeaders, headerOrder, body, bodyTruncated,
	))
	if err != nil {
		logs.Warn("apigateway waf inspection failed, denying request", logs.String("stage", stageARN), logs.Err(err))
		http.Error(w, "The request could not be inspected", http.StatusInternalServerError)
		return false
	}
	if result.Interrupts() {
		status := result.ResponseCode
		if status == 0 {
			status = http.StatusForbidden
		}
		for _, h := range result.ResponseHeaders {
			w.Header().Set(h.Name, h.Value)
		}
		w.WriteHeader(status)
		if result.ResponseBody != "" {
			_, _ = io.WriteString(w, result.ResponseBody)
		}
		return false
	}
	for _, h := range result.InsertHeaders {
		// Inserted header names arrive prefixed with x-amzn-waf-, so
		// adding cannot overwrite a header the client sent.
		r.Header.Add(h.Name, h.Value)
	}
	return true
}

func remoteAddrHost(remoteAddr string) string {
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	return remoteAddr
}

// RemoveApiKey cleans up rate limiter state for a deleted API key.
func (s *RuntimeServer) RemoveApiKey(apiKeyId string) {
	if s.authenticator != nil {
		s.authenticator.RemoveApiKey(apiKeyId)
	}
}

// CleanupStageThrottlers removes cached rate limiters for the given stage,
// preventing unbounded growth of the stageThrottlers map when stages are
// deleted or recreated.
func (s *RuntimeServer) CleanupStageThrottlers(stageName string) {
	prefix := stageName + ":"
	s.stageThrottlers.Range(func(key, _ interface{}) bool {
		if k, ok := key.(string); ok && strings.HasPrefix(k, prefix) {
			s.stageThrottlers.Delete(k)
		}
		return true
	})
}

// Close stops background goroutines in authentication components.
func (s *RuntimeServer) Close() {
	if s.lambdaAuthorizer != nil {
		s.lambdaAuthorizer.Close()
	}
}

// HandleRequest handles incoming API Gateway requests.
func (s *RuntimeServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	s.inspectorMu.RLock()
	inspector := s.webACLInspector
	s.inspectorMu.RUnlock()
	if eventbus.ServeWAFTokenExchange(r.Context(), inspector, w, r) {
		return
	}

	var restApiID, stageName, requestPath string

	fqdnApiID := fqdnrouter.ResourceIDFromContext(r.Context())
	if fqdnApiID != "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 2 || pathParts[1] != "_user_request_" {
			s.sendError(w, http.StatusBadRequest, "Invalid request path format")
			return
		}
		restApiID = fqdnApiID
		stageName = pathParts[0]
		requestPath = "/" + strings.Join(pathParts[2:], "/")
	} else {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 4 {
			s.sendError(w, http.StatusBadRequest, "Invalid request path format")
			return
		}
		if pathParts[0] != "restapis" || pathParts[3] != "_user_request_" {
			s.sendError(w, http.StatusBadRequest, "Invalid request path format")
			return
		}
		restApiID = pathParts[1]
		stageName = pathParts[2]
		requestPath = "/" + strings.Join(pathParts[4:], "/")
	}

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

	if !s.enforceWebACL(w, r, arn.NewARNBuilder(s.accountID, s.region).APIGateway().Stage(restApiID, stageName)) {
		return
	}

	_, ok = restAPI.Deployments[stage.DeploymentId]
	if !ok {
		s.sendError(w, http.StatusNotFound, fmt.Sprintf("Deployment %s not found", stage.DeploymentId))
		return
	}

	var activeResources map[string]*apigatewaystore.Resource
	deployment := restAPI.Deployments[stage.DeploymentId]
	if deployment != nil && deployment.Snapshot != nil && len(deployment.Snapshot.Resources) > 0 {
		activeResources = deployment.Snapshot.Resources
	} else {
		activeResources = restAPI.Resources
	}

	resourceInfos := make([]ResourceInfo, 0, len(activeResources))
	for _, res := range activeResources {
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

	if err := s.checkStageThrottling(stage, match); err != nil {
		s.sendError(w, http.StatusTooManyRequests, err.Error())
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

	// AWS default payload limit is 10 MiB; allow operator override.
	maxBodyBytes := int64(10 * 1024 * 1024)
	if v := config.GetInt("apigateway.max_body_size_bytes"); v > 0 {
		maxBodyBytes = int64(v)
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		s.sendError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read request body: %v", err))
		return
	}
	if int64(len(body)) > maxBodyBytes {
		s.sendError(w, http.StatusRequestEntityTooLarge, "Request body exceeds the maximum allowed payload size")
		return
	}

	if err := s.validator.ValidateRequest(match.Method, restAPI, headers, queryParams, pathParams, body); err != nil {
		s.sendError(w, err.HTTPCode, err.Message)
		return
	}

	s.executeIntegration(w, r, match, stage, restApiID, restAPI.BinaryMediaTypes, body, headers, queryParams, pathParams)
}

func (s *RuntimeServer) executeIntegration(w http.ResponseWriter, r *http.Request, match *RouteMatch, stage *apigatewaystore.Stage, restApiID string, binaryMediaTypes []string, body []byte, headers, queryParams, pathParams map[string]string) {
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
		PassthroughBehavior:         methodIntegration.PassthroughBehavior,
		TimeoutInMillis:             methodIntegration.TimeoutInMillis,
		BinaryMediaTypes:            binaryMediaTypes,
		IntegrationContentHandling:  methodIntegration.ContentHandling,
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
	for key, values := range resp.MultiValueHeaders {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// An integration result can never be handed to WriteHeader with an
	// invalid code; answer configuration-shaped failures with 502.
	if resp.StatusCode < 100 || resp.StatusCode > 599 {
		resp.StatusCode = http.StatusBadGateway
	}

	w.WriteHeader(resp.StatusCode)

	if len(resp.Body) > 0 {
		if resp.IsBase64Encoded {
			if decoded, err := base64.StdEncoding.DecodeString(string(resp.Body)); err == nil {
				w.Write(decoded)
				return
			}
		}
		if _, err := w.Write(resp.Body); err != nil {
			logs.Error("Failed to write response body", logs.Err(err))
		}
	}
}

// checkStageThrottling evaluates stage MethodSettings throttling limits for
// the matched route. If throttling is configured and the rate limit is
// exceeded, an error is returned resulting in HTTP 429.
//
// AWS API Gateway supports wildcard method settings keys. Settings are
// merged from most general (*/*) to most specific (exact path/method),
// with non-zero values at each level overriding less specific ones.
func (s *RuntimeServer) checkStageThrottling(stage *apigatewaystore.Stage, match *RouteMatch) error {
	if stage == nil || len(stage.MethodSettings) == 0 {
		return nil
	}

	httpMethod := match.HttpMethod
	encodedPath := strings.ReplaceAll(match.Path, "/", "~1")
	exactKey := encodedPath + "/" + httpMethod

	// Merge throttling settings from most general to most specific.
	// Non-zero values override inherited values at each level.
	var rate, burst float64
	candidates := []string{
		"*/*",
		"*/" + httpMethod,
		encodedPath + "/*",
		exactKey,
	}
	for _, key := range candidates {
		if ms, ok := stage.MethodSettings[key]; ok {
			if ms.ThrottlingRateLimit > 0 {
				rate = ms.ThrottlingRateLimit
			}
			if ms.ThrottlingBurstLimit > 0 {
				burst = float64(ms.ThrottlingBurstLimit)
			}
		}
	}

	if rate <= 0 && burst <= 0 {
		return nil
	}

	// Apply AWS defaults for unset values
	if rate <= 0 {
		rate = 1000
	}
	if burst <= 0 {
		burst = 2000
	}

	limiterKey := stage.StageName + ":" + exactKey
	val, loaded := s.stageThrottlers.LoadOrStore(limiterKey, newStageRateLimiter(rate, burst))
	limiter, ok := val.(*stageRateLimiter)
	if !ok {
		return nil
	}
	if loaded {
		// Refresh limiter with latest settings so that UpdateStage
		// changes take effect without a server restart.
		limiter.update(rate, burst)
	}
	if !limiter.allow() {
		return fmt.Errorf("Rate limit exceeded")
	}
	return nil
}

type stageRateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
}

func newStageRateLimiter(rateLimit, burstLimit float64) *stageRateLimiter {
	return &stageRateLimiter{
		tokens:     burstLimit,
		maxTokens:  burstLimit,
		refillRate: rateLimit,
		lastRefill: time.Now(),
	}
}

// update adjusts the rate and burst limits on an existing limiter,
// allowing UpdateStage changes to take effect without a restart.
func (r *stageRateLimiter) update(rateLimit, burstLimit float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	oldMax := r.maxTokens
	r.refillRate = rateLimit
	r.maxTokens = burstLimit
	if burstLimit > oldMax {
		// Burst capacity increased — replenish immediately so the new
		// capacity is available without waiting for natural refill.
		r.tokens = burstLimit
	} else if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
}

func (r *stageRateLimiter) allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens = r.tokens + elapsed*r.refillRate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
	r.lastRefill = now

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
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
	case http.StatusMethodNotAllowed:
		errorResp["__type"] = "MethodNotAllowedException"
	case http.StatusConflict:
		errorResp["__type"] = "ConflictException"
	case http.StatusRequestEntityTooLarge:
		errorResp["__type"] = "RequestTooLargeException"
	case http.StatusUnsupportedMediaType:
		errorResp["__type"] = "UnsupportedMediaTypeException"
	case http.StatusUnprocessableEntity:
		errorResp["__type"] = "UnprocessableEntityException"
	case http.StatusTooManyRequests:
		errorResp["__type"] = "TooManyRequestsException"
	case http.StatusInternalServerError:
		errorResp["__type"] = "InternalServerError"
	case http.StatusBadGateway:
		errorResp["__type"] = "BadGatewayException"
	case http.StatusServiceUnavailable:
		errorResp["__type"] = "ServiceUnavailableException"
	case http.StatusGatewayTimeout:
		errorResp["__type"] = "GatewayTimeoutException"
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		logs.Error("Failed to encode error response", logs.Err(err))
	}
}

// clientIP extracts the client IP address from the request. The
// X-Forwarded-For header is trusted only when the direct connection
// originates from a loopback address (i.e. a local reverse proxy such
// as nginx). This prevents spoofing by remote clients.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// RemoteAddr may be a bare IP without a port.
		if ip := net.ParseIP(r.RemoteAddr); ip != nil {
			host = ip.String()
		} else {
			host = r.RemoteAddr
		}
	}

	// Trust X-Forwarded-For only from loopback (local reverse proxy).
	if isLoopback(host) {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if idx := strings.Index(xff, ","); idx >= 0 {
				return strings.TrimSpace(xff[:idx])
			}
			return strings.TrimSpace(xff)
		}
	}

	return host
}

// isLoopback returns true if the given IP string is a loopback address.
func isLoopback(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
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
		return fmt.Sprintf("%s %s %s %d %dms",
			sanitizeLogValue(r.Method), sanitizeLogValue(r.URL.Path),
			sanitizeLogValue(r.Proto), statusCode, latency.Milliseconds())
	}

	requestID := r.Header.Get("X-Amzn-RequestId")
	if requestID == "" {
		requestID = r.Header.Get("X-Amz-Request-Id")
	}
	if requestID == "" {
		requestID = r.Header.Get("X-Amzn-Trace-Id")
	}
	if requestID == "" {
		requestID = generateRequestID()
	}

	sourceIP := clientIP(r)

	result := format
	result = strings.ReplaceAll(result, "$context.requestId", sanitizeLogValue(requestID))
	result = strings.ReplaceAll(result, "$context.requestTime", time.Now().UTC().Format("02/Jan/2006:15:04:05 +0000"))
	result = strings.ReplaceAll(result, "$context.httpMethod", sanitizeLogValue(r.Method))
	result = strings.ReplaceAll(result, "$context.resourcePath", sanitizeLogValue(resourcePath))
	result = strings.ReplaceAll(result, "$context.path", sanitizeLogValue(r.URL.Path))
	result = strings.ReplaceAll(result, "$context.protocol", sanitizeLogValue(r.Proto))
	result = strings.ReplaceAll(result, "$context.status", fmt.Sprintf("%d", statusCode))
	result = strings.ReplaceAll(result, "$context.responseLatency", fmt.Sprintf("%d", latency.Milliseconds()))
	result = strings.ReplaceAll(result, "$context.integrationLatency", fmt.Sprintf("%d", latency.Milliseconds()))
	result = strings.ReplaceAll(result, "$context.sourceIp", sanitizeLogValue(sourceIP))
	result = strings.ReplaceAll(result, "$context.accountId", s.accountID)
	result = strings.ReplaceAll(result, "$context.apiId", sanitizeLogValue(restApiID))
	result = strings.ReplaceAll(result, "$context.stage", sanitizeLogValue(stageName))
	result = strings.ReplaceAll(result, "$context.error.message", "")
	result = strings.ReplaceAll(result, "$context.error.messageString", "")

	if strings.Contains(result, "$context.request.header.") {
		for key, values := range r.Header {
			varName := "$context.request.header." + strings.ToLower(strings.ReplaceAll(key, "-", ""))
			result = strings.ReplaceAll(result, varName, sanitizeLogValue(strings.Join(values, ",")))
		}
	}

	return result
}

// sanitizeLogValue strips control characters (CR, LF, TAB, NUL) from
// user-supplied values before they are interpolated into access log
// lines, preventing log injection attacks.
func sanitizeLogValue(s string) string {
	if !strings.ContainsAny(s, "\r\n\t\x00") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '\r', '\n', '\t', '\x00':
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// generateRequestID produces a random hex string for log correlation when
// no request ID header is present, avoiding the log collisions caused by
// a static "req-unknown" placeholder.
func generateRequestID() string {
	var buf [16]byte
	if _, err := cryptorand.Read(buf[:]); err != nil {
		return fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", buf[:])
}
