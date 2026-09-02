package appsync

import (
	"context"
	"net/http"
	"sync"

	"vorpalstacks/internal/common/headerorder"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/waflimits"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/utils/aws/arn"
)

// wafEnforcement holds the injected WAF request-inspection entry point
// for the AppSync data plane. The zero value (inspector nil) allows all
// traffic, which is the state before the cross-service wiring runs.
type wafEnforcement struct {
	mu   sync.RWMutex
	insp eventbus.WebACLInspector
}

func (w *wafEnforcement) setInspector(inspector eventbus.WebACLInspector) {
	w.mu.Lock()
	w.insp = inspector
	w.mu.Unlock()
}

func (w *wafEnforcement) currentInspector() eventbus.WebACLInspector {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.insp
}

// SetWebACLInspector injects the WAF request-inspection entry point so
// WebACLs associated with GraphQL APIs are enforced on GraphQL
// execution traffic.
func (s *AppSyncService) SetWebACLInspector(inspector eventbus.WebACLInspector) {
	s.waf.setInspector(inspector)
}

// enforceWebACL inspects a GraphQL execution request against the WebACL
// associated with the API. It returns a non-nil GraphQL response (the
// blocked response) when the request must not be executed, and applies
// custom request headers to the header set for allowed requests.
// Inspection failures fail closed with 500: the AWS WAF Developer Guide
// documents that when WAF encounters an internal error, Regional
// services typically deny the request and don't serve the content.
func (s *AppSyncService) enforceWebACL(ctx context.Context, reqCtx *request.RequestContext, apiId string, headers http.Header, host string, body []byte) *graphqlResponse {
	inspector := s.waf.currentInspector()
	if inspector == nil || apiId == "" {
		return nil
	}

	// AppSync's body inspection limit is fixed at 8 KB, unlike the
	// configurable 16 KB default the other protected-resource planes
	// share.
	inspBody := body
	bodyTruncated := false
	if int64(len(inspBody)) > waflimits.AppSyncBodyInspectionLimit {
		inspBody = inspBody[:waflimits.AppSyncBodyInspectionLimit]
		bodyTruncated = true
	}

	apiArn := arn.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).AppSync().GraphQLApi(apiId)
	inspHeaders := eventbus.RequestHeadersWithHost(headers, host)
	headerOrder, _ := headerorder.FromContext(ctx, inspHeaders)
	result, err := inspector.InspectWebACLRequest(ctx, reqCtx.GetRegion(), apiArn, eventbus.BuildWebACLInspectionRequest(
		http.MethodPost, "/v1/apis/"+apiId+"/graphql", "", reqCtx.SourceIP, "",
		inspHeaders, headerOrder, inspBody, bodyTruncated,
	))
	if err != nil {
		logs.Warn("appsync waf inspection failed, denying request", logs.String("api", apiId), logs.Err(err))
		return s.graphqlErrorResponse(http.StatusInternalServerError, "WAFUnavailable", "The request could not be inspected by the associated AWS WAF Web ACL")
	}
	if !result.Interrupts() {
		// Inserted header names arrive prefixed with x-amzn-waf-, so
		// adding cannot overwrite a header the client sent.
		for _, h := range result.InsertHeaders {
			headers.Add(h.Name, h.Value)
		}
		return nil
	}
	status := result.ResponseCode
	if status == 0 {
		status = http.StatusForbidden
	}
	response := s.graphqlErrorResponse(status, "WAFBlocked", "The request was blocked by an associated AWS WAF Web ACL")
	if result.ResponseBody != "" {
		response = &graphqlResponse{
			body:    []byte(result.ResponseBody),
			headers: http.Header{"Content-Type": []string{"text/plain"}},
			status:  status,
		}
	}
	// The inspection's response headers (a custom Block response's
	// headers, and the x-amzn-waf-action header an interrupting Captcha
	// or Challenge carries) apply on this plane too.
	if response.headers == nil {
		response.headers = http.Header{}
	}
	for _, h := range result.ResponseHeaders {
		response.headers.Set(h.Name, h.Value)
	}
	return response
}
