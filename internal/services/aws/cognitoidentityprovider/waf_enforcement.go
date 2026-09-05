package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/headerorder"
	waf "vorpalstacks/internal/common/invokers/waf"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/waflimits"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// wafEnforcement holds the injected WAF request-inspection entry point
// for the Cognito data plane. The zero value (inspector nil) allows all
// traffic, which is the state before the cross-service wiring runs.
type wafEnforcement struct {
	mu   sync.RWMutex
	insp waf.WebACLInspector
}

func (w *wafEnforcement) setInspector(inspector waf.WebACLInspector) {
	w.mu.Lock()
	w.insp = inspector
	w.mu.Unlock()
}

func (w *wafEnforcement) currentInspector() waf.WebACLInspector {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.insp
}

// SetWebACLInspector injects the WAF request-inspection entry point so
// WebACLs associated with user pools are enforced on the hosted UI and
// the public user pools API operations.
func (s *CognitoService) SetWebACLInspector(inspector waf.WebACLInspector) {
	s.waf.setInspector(inspector)
}

// wafInspectedOperations lists the user pools API operations whose
// requests AWS WAF inspects: those that do not require authentication
// with AWS credentials (unauthenticated, or authorized with a session
// string or access token). Management operations authenticated with AWS
// credentials are not inspected.
var wafInspectedOperations = map[string]bool{
	"SignUp":                           true,
	"ConfirmSignUp":                    true,
	"ResendConfirmationCode":           true,
	"InitiateAuth":                     true,
	"RespondToAuthChallenge":           true,
	"ForgotPassword":                   true,
	"ConfirmForgotPassword":            true,
	"GetUser":                          true,
	"UpdateUserAttributes":             true,
	"DeleteUser":                       true,
	"DeleteUserAttributes":             true,
	"GlobalSignOut":                    true,
	"ChangePassword":                   true,
	"GetUserAttributeVerificationCode": true,
	"VerifyUserAttribute":              true,
	"AssociateSoftwareToken":           true,
	"VerifySoftwareToken":              true,
	"SetUserMFAPreference":             true,
	"RevokeToken":                      true,
}

// cognitoOperationHandler is the dispatcher handler signature for the
// Cognito user pools API.
type cognitoOperationHandler = handler.Handler

// withWAFEnforcement wraps a public user pools API operation handler
// with the WebACL inspection that AWS applies to credential-free
// requests.
func (s *CognitoService) withWAFEnforcement(operation string, h cognitoOperationHandler) cognitoOperationHandler {
	return func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
		if err := s.enforceWAFOnAPIRequest(ctx, reqCtx, req, operation); err != nil {
			return nil, err
		}
		return h(ctx, reqCtx, req)
	}
}

// enforceWAFOnAPIRequest resolves the user pool the request addresses,
// inspects the request against its associated WebACL, and returns a
// ForbiddenException-shaped error for blocked requests. Inspection
// failures fail closed with InternalErrorException: the AWS WAF
// Developer Guide documents that when WAF encounters an internal
// error, Regional services typically deny the request and don't serve
// the content.
func (s *CognitoService) enforceWAFOnAPIRequest(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest, operation string) error {
	inspector := s.waf.currentInspector()
	if inspector == nil || !wafInspectedOperations[operation] {
		return nil
	}
	pool, err := s.poolForWAFRequest(reqCtx, req)
	if err != nil || pool == nil || pool.Arn == "" {
		if err != nil {
			logs.Warn("cognito waf pool resolution failed, serving request", logs.String("operation", operation), logs.Err(err))
		}
		return nil
	}

	body := redactPIIForWAF(req.Body)
	truncated := false
	if int64(len(body)) > waflimits.DefaultBodyInspectionLimit {
		body = body[:waflimits.DefaultBodyInspectionLimit]
		truncated = true
	}

	inspHeaders := waf.RequestHeadersWithHost(req.Headers, req.Host)
	headerOrder, _ := headerorder.FromContext(ctx, inspHeaders)
	result, err := inspector.InspectWebACLRequest(ctx, reqCtx.GetRegion(), pool.Arn, waf.BuildWebACLInspectionRequest(
		http.MethodPost, "/", "", reqCtx.SourceIP, "",
		inspHeaders, headerOrder, body, truncated,
	))
	if err != nil {
		logs.Warn("cognito waf inspection failed, denying request", logs.String("operation", operation), logs.Err(err))
		return awserrors.NewAWSError("InternalErrorException", "The request could not be inspected by the associated AWS WAF Web ACL", http.StatusInternalServerError)
	}
	if !result.Interrupts() {
		return nil
	}
	status := result.ResponseCode
	if status == 0 || status < 400 || status > 499 {
		// The user pools API surfaces Block responses as
		// ForbiddenException; a configured custom code replaces the
		// status within the 400-499 range.
		status = http.StatusForbidden
	}
	// The Cognito WAF integration documents that the CLI and SDKs see a
	// ForbiddenException for requests that produce a Block or CAPTCHA
	// response. The Challenge action is absent from the integration's
	// documented action list and this API plane cannot serve an
	// interstitial, so it joins the same error path; the hosted UI plane
	// serves the real challenge responses.
	message := "Request blocked by WAF"
	if result.ResponseBody != "" {
		message = result.ResponseBody
	}
	return awserrors.NewAWSError("ForbiddenException", message, status)
}

// poolForWAFRequest resolves the user pool an API request addresses:
// directly by UserPoolId, by ClientId for the auth-flow operations, or
// through the stored access token record for token-authorised
// operations. A nil pool means no WebACL can apply.
func (s *CognitoService) poolForWAFRequest(reqCtx *request.RequestContext, req *request.ParsedRequest) (*cognitostore.UserPool, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if id := request.GetStringParam(req.Parameters, "UserPoolId"); id != "" {
		return store.GetUserPool(id)
	}
	if clientID := request.GetStringParam(req.Parameters, "ClientId"); clientID != "" {
		return store.GetUserPoolByClientID(clientID)
	}
	if token := request.GetStringParam(req.Parameters, "AccessToken"); token != "" {
		record, err := store.GetAccessTokenByValue(token)
		if err != nil {
			return nil, err
		}
		return store.GetUserPool(record.UserPoolID)
	}
	return nil, nil
}

// redactPIIForWAF strips the PII-bearing members from an API request
// body before WAF inspection: Amazon Cognito forwards only
// non-confidential request contents to AWS WAF, and its user-pool WAF
// documentation states rules cannot match on personally identifiable
// information — usernames, passwords, phone numbers and email
// addresses are not available to AWS WAF. The redaction drops those
// members; bodies that fail to parse are forwarded as-is after
// truncation.
func redactPIIForWAF(body []byte) []byte {
	if len(body) == 0 {
		return body
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return body
	}
	redacted := false
	for _, key := range []string{"Username", "Password", "PreviousPassword", "ProposedPassword"} {
		if _, ok := parsed[key]; ok {
			delete(parsed, key)
			redacted = true
		}
	}
	if attributes, ok := parsed["UserAttributes"].([]interface{}); ok {
		kept := make([]interface{}, 0, len(attributes))
		for _, entry := range attributes {
			if attribute, ok := entry.(map[string]interface{}); ok {
				if name, _ := attribute["Name"].(string); name == "email" || name == "phone_number" {
					redacted = true
					continue
				}
			}
			kept = append(kept, entry)
		}
		if len(kept) != len(attributes) {
			parsed["UserAttributes"] = kept
		}
	}
	if authParams, ok := parsed["AuthParameters"].(map[string]interface{}); ok {
		for _, key := range []string{"PASSWORD", "SECRET_HASH", "NEW_PASSWORD"} {
			if _, ok := authParams[key]; ok {
				delete(authParams, key)
				redacted = true
			}
		}
	}
	if !redacted {
		return body
	}
	out, err := json.Marshal(parsed)
	if err != nil {
		return body
	}
	return out
}

// enforceWAFOnHostedUI inspects a hosted UI request against the pool's
// associated WebACL. The hosted UI forwards no request body to AWS WAF,
// so only the headers and path are inspected. Inspection failures fail
// closed with 500, matching the Regional-service behaviour the AWS WAF
// Developer Guide documents for WAF internal errors. The return value
// reports whether the request was answered (blocked or failed).
func (s *CognitoService) enforceWAFOnHostedUI(w http.ResponseWriter, r *http.Request, poolARN string) bool {
	inspector := s.waf.currentInspector()
	if inspector == nil || poolARN == "" {
		return false
	}
	inspHeaders := waf.RequestHeadersWithHost(r.Header, r.Host)
	headerOrder, _ := headerorder.FromContext(r.Context(), inspHeaders)
	result, err := inspector.InspectWebACLRequest(r.Context(), s.region, poolARN, waf.BuildWebACLInspectionRequest(
		r.Method, r.URL.Path, r.URL.RawQuery, remoteAddrHostOf(r.RemoteAddr), r.Proto,
		inspHeaders, headerOrder, nil, false,
	))
	if err != nil {
		logs.Warn("cognito hosted ui waf inspection failed, denying request", logs.String("pool", poolARN), logs.Err(err))
		http.Error(w, "The request could not be inspected", http.StatusInternalServerError)
		return true
	}
	if !result.Interrupts() {
		// Inserted header names arrive prefixed with x-amzn-waf-, so
		// adding cannot overwrite a header the client sent.
		for _, h := range result.InsertHeaders {
			r.Header.Add(h.Name, h.Value)
		}
		return false
	}
	status := result.ResponseCode
	if status == 0 {
		status = http.StatusForbidden
	}
	for _, h := range result.ResponseHeaders {
		w.Header().Set(h.Name, h.Value)
	}
	w.WriteHeader(status)
	if result.ResponseBody != "" {
		_, _ = w.Write([]byte(result.ResponseBody))
	} else {
		_, _ = w.Write([]byte("Request blocked by WAF"))
	}
	return true
}

// poolARNByID resolves a user pool's ARN for enforcement lookups; an
// empty string means the pool could not be resolved and no WebACL can
// apply.
func (s *CognitoService) poolARNByID(poolID string) string {
	reqCtx := request.NewRequestContext(context.Background(), s.storageManager, s.accountID, s.region)
	store, err := s.store(reqCtx)
	if err != nil {
		return ""
	}
	pool, err := store.GetUserPool(poolID)
	if err != nil || pool == nil {
		return ""
	}
	return pool.Arn
}

func remoteAddrHostOf(remoteAddr string) string {
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	return remoteAddr
}
