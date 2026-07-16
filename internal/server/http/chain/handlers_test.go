package chain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddRequestContext(t *testing.T) {
	rc := &RequestContext{
		Request: &http.Request{},
	}

	err := AddRequestContext(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAddRequestContext_NilRequest(t *testing.T) {
	rc := &RequestContext{}

	err := AddRequestContext(rc)
	if err == nil {
		t.Error("expected error for nil request")
	}
}

func TestParseRequestParams_GetMethod(t *testing.T) {
	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodGet, "/test?key1=value1&key2=value2", nil),
		Params:  make(map[string]interface{}),
	}

	err := ParseRequestParams(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rc.Params["key1"] != "value1" {
		t.Errorf("expected key1=value1, got %v", rc.Params["key1"])
	}
}

func TestParseRequestParams_PostJSON(t *testing.T) {
	body := map[string]string{"key1": "value1"}
	bodyBytes, _ := json.Marshal(body)

	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(string(bodyBytes))),
		Params:  make(map[string]interface{}),
	}
	rc.Request.Header.Set("Content-Type", "application/json")

	err := ParseRequestParams(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParseRequestParams_PostForm(t *testing.T) {
	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("key1=value1")),
		Params:  make(map[string]interface{}),
	}
	rc.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err := ParseRequestParams(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParseRequestParams_EmptyBody(t *testing.T) {
	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodPost, "/test", nil),
		Params:  make(map[string]interface{}),
	}
	rc.Request.Header.Set("Content-Type", "application/json")

	err := ParseRequestParams(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAddRegionFromHeader_NoAuthHeader(t *testing.T) {
	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		Params:  make(map[string]interface{}),
	}

	err := AddRegionFromHeader(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rc.Params["Region"] != nil {
		t.Error("expected no region without auth header")
	}
}

func TestAddRegionFromHeader_WithCredential(t *testing.T) {
	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		Params:  make(map[string]interface{}),
	}
	rc.Request.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20230101/us-east-1/service/aws4_request")

	err := AddRegionFromHeader(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rc.Params["Region"] != "us-east-1" {
		t.Errorf("expected region us-east-1, got %v", rc.Params["Region"])
	}
}

func TestEnforceCORS_NoOrigin(t *testing.T) {
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: httptest.NewRecorder(),
	}

	err := EnforceCORS(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEnforceCORS_WithOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}
	rc.Request.Header.Set("Origin", "http://example.com")

	err := EnforceCORS(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header")
	}
}

func TestEnforceCORS_OptionsMethod(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodOptions, "/test", nil),
		Response: w,
	}
	rc.Request.Header.Set("Origin", "http://example.com")

	err := EnforceCORS(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !rc.IsHandled() {
		t.Error("expected request to be handled")
	}
}

func TestLogRequest(t *testing.T) {
	rc := &RequestContext{
		Request:     httptest.NewRequest(http.MethodGet, "/test", nil),
		ServiceName: "test-service",
		Operation:   "test-op",
	}

	err := LogRequest(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLogException_NilError(t *testing.T) {
	rc := &RequestContext{
		Request: httptest.NewRequest(http.MethodGet, "/test", nil),
	}

	err := LogException(rc, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandleServiceError(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}

	err := HandleServiceError(rc, &testError{"test error"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !rc.IsHandled() {
		t.Error("expected request to be handled")
	}
}

func TestHandleInternalFailure(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}

	err := HandleInternalFailure(rc, &testError{"test error"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !rc.IsHandled() {
		t.Error("expected request to be handled")
	}
}

func TestHandleInternalFailure_AlreadyHandled(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}
	rc.SetHandled()

	err := HandleInternalFailure(rc, &testError{"test error"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSerializeResponse_NilBody(t *testing.T) {
	rc := &RequestContext{
		Response: httptest.NewRecorder(),
	}

	err := SerializeResponse(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSerializeResponse_JSON(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:      httptest.NewRequest(http.MethodGet, "/test", nil),
		Response:     w,
		ResponseBody: map[string]string{"key": "value"},
	}

	err := SerializeResponse(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("expected content-type application/json")
	}
}

func TestSerializeResponse_XML(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:      httptest.NewRequest(http.MethodGet, "/test", nil),
		Response:     w,
		ResponseBody: "<xml></xml>",
	}
	rc.Request.Header.Set("Accept", "application/xml")

	err := SerializeResponse(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Header().Get("Content-Type") != "application/xml" {
		t.Error("expected content-type application/xml")
	}
}

func TestAddCORSResponseHeaders_NoOrigin(t *testing.T) {
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: httptest.NewRecorder(),
	}

	err := AddCORSResponseHeaders(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAddCORSResponseHeaders_WithOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}
	rc.Request.Header.Set("Origin", "http://example.com")

	err := AddCORSResponseHeaders(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header")
	}
}

func TestLogResponse(t *testing.T) {
	rc := &RequestContext{
		StatusCode:  200,
		ServiceName: "test-service",
		Operation:   "test-op",
	}

	err := LogResponse(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSetCloseConnectionHeader_NoConnection(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}

	err := SetCloseConnectionHeader(rc)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Header().Get("Connection") != "" {
		t.Errorf("expected no Connection header, got %q", w.Header().Get("Connection"))
	}
}

func TestSetCloseConnectionHeader_WithClose(t *testing.T) {
	w := httptest.NewRecorder()
	rc := &RequestContext{
		Request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		Response: w,
	}
	rc.Request.Header.Set("Connection", "close")

	SetCloseConnectionHeader(rc)

	if w.Header().Get("Connection") != "close" {
		t.Error("expected connection close header")
	}
}

func TestExtractHost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Host = "example.com:8080"

	host := extractHost(req)
	if host != "example.com" {
		t.Errorf("expected example.com, got %s", host)
	}
}

func TestExtractPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)

	path := extractPath(req)
	if path != "/test/path" {
		t.Errorf("expected /test/path, got %s", path)
	}
}

func TestExtractQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?key=value", nil)

	query := extractQuery(req)
	if query.Get("key") != "value" {
		t.Error("expected key=value")
	}
}
