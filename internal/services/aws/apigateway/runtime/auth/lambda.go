package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common"
	commonauth "vorpalstacks/internal/common/auth"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	arnutil "vorpalstacks/internal/utils/aws/arn"

	"vorpalstacks/internal/config"
)

// LambdaAuthorizer handles Lambda-based authorizer for API Gateway.
type LambdaAuthorizer struct {
	lambdaInvoker common.LambdaInvoker
	store         *apigatewaystore.RestApiStore
	cache         *authCache
	accountID     string
	region        string
	sigVerifier   *commonauth.SignatureV4Verifier
}

// NewLambdaAuthorizer creates a new Lambda authorizer instance.
func NewLambdaAuthorizer(lambdaInvoker common.LambdaInvoker, store *apigatewaystore.RestApiStore) *LambdaAuthorizer {
	return &LambdaAuthorizer{
		lambdaInvoker: lambdaInvoker,
		store:         store,
		cache:         &authCache{stopCh: make(chan struct{})},
	}
}

type authCache struct {
	results     sync.Map
	cleanupOnce sync.Once
	stopCh      chan struct{}
	stopOnce    sync.Once
}

// Close stops the background cleanup goroutine. Safe to call multiple times.
func (ac *authCache) Close() {
	ac.stopOnce.Do(func() { close(ac.stopCh) })
}

// NewLambdaAuthorizerWithConfig creates a new Lambda authorizer with account ID and region configuration.
func NewLambdaAuthorizerWithConfig(lambdaInvoker common.LambdaInvoker, store *apigatewaystore.RestApiStore, accountID, region string, credProvider commonauth.CredentialsProvider) *LambdaAuthorizer {
	return &LambdaAuthorizer{
		lambdaInvoker: lambdaInvoker,
		store:         store,
		cache:         &authCache{stopCh: make(chan struct{})},
		accountID:     accountID,
		region:        region,
		sigVerifier:   commonauth.NewSignatureV4Verifier(credProvider),
	}
}

// Close stops the background authorisation cache cleanup goroutine.
func (la *LambdaAuthorizer) Close() {
	la.cache.Close()
}

// LambdaAuthEvent represents the event sent to a Lambda authorizer.
type LambdaAuthEvent struct {
	Type               string                 `json:"type"`
	AuthorizationToken string                 `json:"authorizationToken"`
	MethodArn          string                 `json:"methodArn"`
	Resource           string                 `json:"resource"`
	Path               string                 `json:"path"`
	HttpMethod         string                 `json:"httpMethod"`
	Headers            map[string]string      `json:"headers"`
	QueryStringParams  map[string]string      `json:"queryStringParameters"`
	PathParameters     map[string]string      `json:"pathParameters"`
	StageVariables     map[string]string      `json:"stageVariables"`
	RequestContext     map[string]interface{} `json:"requestContext"`
}

// LambdaAuthResponse represents the response from a Lambda authorizer.
type LambdaAuthResponse struct {
	PrincipalID        string                 `json:"principalId"`
	PolicyDocument     PolicyDocument         `json:"policyDocument"`
	Context            map[string]interface{} `json:"context,omitempty"`
	UsageIdentifierKey string                 `json:"usageIdentifierKey,omitempty"`
}

// PolicyDocument defines the IAM policy document for authorization.
type PolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// Statement represents a single statement in an IAM policy document.
type Statement struct {
	Effect    string                 `json:"Effect"`
	Action    string                 `json:"Action"`
	Resource  string                 `json:"Resource"`
	Condition map[string]interface{} `json:"Condition,omitempty"`
}

// Authorize performs authorization using the configured Lambda authorizer.
func (la *LambdaAuthorizer) Authorize(ctx context.Context, method *apigatewaystore.Method, restAPI *apigatewaystore.RestApi, req *AuthRequest) (*AuthResult, error) {
	if method == nil || method.AuthorizationType == "" || method.AuthorizationType == "NONE" {
		return &AuthResult{Allowed: true}, nil
	}

	if method.AuthorizationType == "AWS_IAM" {
		return la.authorizeIAM(ctx, req)
	}

	if method.AuthorizerId == "" {
		return nil, &AuthError{
			Message:  "No authorizer configured",
			Type:     "UnauthorizedException",
			HTTPCode: http.StatusUnauthorized,
		}
	}

	authorizer, ok := restAPI.Authorizers[method.AuthorizerId]
	if !ok {
		return nil, &AuthError{
			Message:  "Authorizer not found",
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	switch authorizer.Type {
	case "TOKEN":
		return la.authorizeToken(ctx, authorizer, method, req)
	case "REQUEST":
		return la.authorizeRequest(ctx, authorizer, method, req)
	default:
		return nil, &AuthError{
			Message:  fmt.Sprintf("Unsupported authorizer type: %s", authorizer.Type),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}
}

func (la *LambdaAuthorizer) authorizeIAM(ctx context.Context, req *AuthRequest) (*AuthResult, error) {
	if la.sigVerifier == nil {
		return nil, &AuthError{
			Message:  "IAM authentication not configured",
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	authHeader := req.Headers["Authorization"]
	if authHeader == "" {
		return nil, &AuthError{
			Message:  "Missing Authorization header",
			Type:     "UnauthorizedException",
			HTTPCode: http.StatusUnauthorized,
		}
	}

	httpReq, err := la.buildHTTPRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := la.sigVerifier.VerifyRequest(httpReq, "execute-api", la.region); err != nil {
		return nil, &AuthError{
			Message:  fmt.Sprintf("IAM signature verification failed: %v", err),
			Type:     "UnauthorizedException",
			HTTPCode: http.StatusUnauthorized,
		}
	}

	return &AuthResult{Allowed: true}, nil
}

func (la *LambdaAuthorizer) buildHTTPRequest(ctx context.Context, req *AuthRequest) (*http.Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, "https://apigateway."+la.region+".amazonaws.com/"+req.RestApiId+"/"+req.StageName+req.Resource, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header = http.Header{}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	return httpReq, nil
}

func (la *LambdaAuthorizer) buildMethodArn(req *AuthRequest) string {
	region := la.region
	if region == "" {
		region = "us-east-1"
	}
	accountID := la.accountID
	if accountID == "" {
		accountID = config.AWSAccountID()
	}
	return arnutil.NewARNBuilder(accountID, region).APIGateway().ExecuteApi(req.RestApiId, req.StageName, req.Method, req.Resource)
}

func (la *LambdaAuthorizer) authorizeToken(ctx context.Context, authorizer *apigatewaystore.Authorizer, method *apigatewaystore.Method, req *AuthRequest) (*AuthResult, error) {
	authToken := la.extractToken(authorizer.IdentitySource, req.Headers)
	cacheKey := fmt.Sprintf("%s:%s:%s", authorizer.Id, authToken, req.Resource)
	if cached, ok := la.getCachedResult(cacheKey); ok {
		return cached, nil
	}

	functionName, err := extractFunctionNameFromURI(authorizer.AuthorizerUri)
	if err != nil {
		return nil, &AuthError{
			Message:  fmt.Sprintf("Invalid authorizer URI: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	methodArn := la.buildMethodArn(req)

	event := LambdaAuthEvent{
		Type:               "TOKEN",
		AuthorizationToken: authToken,
		MethodArn:          methodArn,
		Resource:           req.Resource,
		Path:               req.Path,
		HttpMethod:         req.Method,
		Headers:            req.Headers,
		QueryStringParams:  req.QueryStringParams,
		PathParameters:     req.PathParameters,
		StageVariables:     req.StageVariables,
		RequestContext: map[string]interface{}{
			"apiId":      req.RestApiId,
			"stage":      req.StageName,
			"resource":   req.Resource,
			"httpMethod": req.Method,
			"requestId":  fmt.Sprintf("%x", time.Now().UnixNano()),
		},
	}

	result, err := la.invokeAuthorizer(ctx, functionName, event)
	if err != nil {
		return nil, err
	}

	ttl := authorizer.AuthorizerResultTtlInSeconds
	if ttl > 0 {
		la.setCachedResult(cacheKey, result, time.Duration(ttl)*time.Second)
	}

	return result, nil
}

func (la *LambdaAuthorizer) authorizeRequest(ctx context.Context, authorizer *apigatewaystore.Authorizer, method *apigatewaystore.Method, req *AuthRequest) (*AuthResult, error) {
	identityValues := la.extractIdentityValues(authorizer.IdentitySource, req)
	cacheKey := fmt.Sprintf("%s:%s:%s:%s", authorizer.Id, req.Method, req.Resource, identityValues)
	if cached, ok := la.getCachedResult(cacheKey); ok {
		return cached, nil
	}

	functionName, err := extractFunctionNameFromURI(authorizer.AuthorizerUri)
	if err != nil {
		return nil, &AuthError{
			Message:  fmt.Sprintf("Invalid authorizer URI: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	methodArn := la.buildMethodArn(req)

	event := LambdaAuthEvent{
		Type:              "REQUEST",
		MethodArn:         methodArn,
		Resource:          req.Resource,
		Path:              req.Path,
		HttpMethod:        req.Method,
		Headers:           req.Headers,
		QueryStringParams: req.QueryStringParams,
		PathParameters:    req.PathParameters,
		StageVariables:    req.StageVariables,
		RequestContext: map[string]interface{}{
			"apiId":      req.RestApiId,
			"stage":      req.StageName,
			"resource":   req.Resource,
			"httpMethod": req.Method,
			"requestId":  fmt.Sprintf("%x", time.Now().UnixNano()),
		},
	}

	result, err := la.invokeAuthorizer(ctx, functionName, event)
	if err != nil {
		return nil, err
	}

	ttl := authorizer.AuthorizerResultTtlInSeconds
	if ttl > 0 {
		la.setCachedResult(cacheKey, result, time.Duration(ttl)*time.Second)
	}

	return result, nil
}

func (la *LambdaAuthorizer) invokeAuthorizer(ctx context.Context, functionName string, event LambdaAuthEvent) (*AuthResult, error) {
	if la.lambdaInvoker == nil {
		return nil, &AuthError{
			Message:  "Lambda invoker not configured",
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, &AuthError{
			Message:  fmt.Sprintf("Failed to marshal authorizer event: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	_, payload, err := la.lambdaInvoker.InvokeForGateway(ctx, functionName, eventJSON)
	if err != nil {
		return nil, &AuthError{
			Message:  fmt.Sprintf("Authorizer invocation failed: %v", err),
			Type:     "UnauthorizedException",
			HTTPCode: http.StatusUnauthorized,
		}
	}

	var authResp LambdaAuthResponse
	if err := json.Unmarshal(payload, &authResp); err != nil {
		return nil, &AuthError{
			Message:  fmt.Sprintf("Failed to parse authorizer response: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	allowed := la.evaluatePolicy(&authResp.PolicyDocument, event.MethodArn)

	return &AuthResult{
		Allowed:     allowed,
		PrincipalID: authResp.PrincipalID,
		Context:     authResp.Context,
	}, nil
}

func (la *LambdaAuthorizer) evaluatePolicy(policy *PolicyDocument, methodArn string) bool {
	if policy == nil || len(policy.Statement) == 0 {
		return false
	}

	for _, stmt := range policy.Statement {
		if stmt.Effect == "Allow" {
			if la.resourceMatches(stmt.Resource, methodArn) {
				return true
			}
		}
	}
	return false
}

func (la *LambdaAuthorizer) resourceMatches(policyResource, methodArn string) bool {
	if policyResource == "*" {
		return true
	}
	return strings.HasPrefix(methodArn, policyResource) || methodArn == policyResource
}

func (la *LambdaAuthorizer) extractToken(identitySource string, headers map[string]string) string {
	if identitySource == "" {
		return ""
	}

	parts := strings.Split(identitySource, " ")
	if len(parts) < 2 {
		if token, ok := headers[identitySource]; ok {
			return token
		}
		if token, ok := headers[strings.ToLower(identitySource)]; ok {
			return token
		}
		return ""
	}

	headerName := parts[1]
	if token, ok := headers[headerName]; ok {
		return parts[0] + " " + token
	}
	if token, ok := headers[strings.ToLower(headerName)]; ok {
		return parts[0] + " " + token
	}
	return ""
}

func (la *LambdaAuthorizer) extractIdentityValues(identitySource string, req *AuthRequest) string {
	if identitySource == "" {
		return ""
	}

	sources := strings.Split(identitySource, ",")
	var values []string
	for _, source := range sources {
		source = strings.TrimSpace(source)
		value := la.extractIdentityValue(source, req)
		if value != "" {
			values = append(values, value)
		}
	}

	return strings.Join(values, "|")
}

func (la *LambdaAuthorizer) extractIdentityValue(source string, req *AuthRequest) string {
	source = strings.TrimSpace(source)

	if strings.HasPrefix(source, "method.request.header.") {
		headerName := strings.TrimPrefix(source, "method.request.header.")
		if val, ok := req.Headers[headerName]; ok {
			return val
		}
		if val, ok := req.Headers[strings.ToLower(headerName)]; ok {
			return val
		}
	} else if strings.HasPrefix(source, "method.request.querystring.") {
		paramName := strings.TrimPrefix(source, "method.request.querystring.")
		if val, ok := req.QueryStringParams[paramName]; ok {
			return val
		}
	} else if strings.HasPrefix(source, "method.request.path.") {
		paramName := strings.TrimPrefix(source, "method.request.path.")
		if val, ok := req.PathParameters[paramName]; ok {
			return val
		}
	} else if strings.HasPrefix(source, "stage.") {
		varName := strings.TrimPrefix(source, "stage.")
		if val, ok := req.StageVariables[varName]; ok {
			return val
		}
	} else {
		if val, ok := req.Headers[source]; ok {
			return val
		}
		if val, ok := req.Headers[strings.ToLower(source)]; ok {
			return val
		}
	}

	return ""
}

func (la *LambdaAuthorizer) getCachedResult(key string) (*AuthResult, bool) {
	if cached, ok := la.cache.results.Load(key); ok {
		cachedResult := cached.(*cachedAuthResult)
		if time.Now().Before(cachedResult.expires) {
			return cachedResult.result, true
		}
		la.cache.results.Delete(key)
	}
	return nil, false
}

func (la *LambdaAuthorizer) setCachedResult(key string, result *AuthResult, ttl time.Duration) {
	la.cache.results.Store(key, &cachedAuthResult{
		result:  result,
		expires: time.Now().Add(ttl),
	})
	la.cache.cleanupIfNeeded()
}

type cachedAuthResult struct {
	result  *AuthResult
	expires time.Time
}

func (ac *authCache) cleanupIfNeeded() {
	ac.cleanupOnce.Do(func() {
		go ac.cleanupLoop()
	})
}

func (ac *authCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ac.stopCh:
			return
		case <-ticker.C:
			now := time.Now()
			var toDelete []string
			ac.results.Range(func(key, value interface{}) bool {
				if cr, ok := value.(*cachedAuthResult); ok && now.After(cr.expires) {
					toDelete = append(toDelete, key.(string))
				}
				return true
			})
			for _, k := range toDelete {
				ac.results.Delete(k)
			}
		}
	}
}

// AuthRequest contains the request context for authorization.
type AuthRequest struct {
	Method            string
	Resource          string
	Path              string
	Headers           map[string]string
	QueryStringParams map[string]string
	PathParameters    map[string]string
	StageVariables    map[string]string
	RestApiId         string
	StageName         string
}

// AuthResult contains the result of an authorization decision.
type AuthResult struct {
	Allowed     bool
	PrincipalID string
	Context     map[string]interface{}
}

// extractFunctionNameFromURI extracts the Lambda function name from an authorizer URI.
func extractFunctionNameFromURI(uri string) (string, error) {
	uri = strings.TrimPrefix(uri, "arn:aws:apigateway:")
	parts := strings.SplitN(uri, ":", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid URI format")
	}

	functionPath := parts[1]
	if !strings.HasPrefix(functionPath, "lambda:path/2015-03-31/functions/") {
		return "", fmt.Errorf("not a Lambda function URI")
	}

	functionPart := strings.TrimPrefix(functionPath, "lambda:path/2015-03-31/functions/")
	functionPart = strings.TrimSuffix(functionPart, "/invocations")
	if idx := strings.LastIndex(functionPart, ":"); idx > 0 {
		return functionPart[idx+1:], nil
	}

	return functionPart, nil
}
