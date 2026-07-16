package protocol

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestEncodeJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		response interface{}
		wantCode int
	}{
		{
			name:     "simple map",
			response: map[string]string{"key": "value"},
			wantCode: 200,
		},
		{
			name:     "slice",
			response: []int{1, 2, 3},
			wantCode: 200,
		},
		{
			name:     "string",
			response: "hello",
			wantCode: 200,
		},
		{
			name:     "nil",
			response: nil,
			wantCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := EncodeJSONResponse(w, tt.response)
			if err != nil {
				t.Errorf("EncodeJSONResponse() error = %v", err)
			}
			if w.Code != tt.wantCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantCode)
			}
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type = %v, want application/json", w.Header().Get("Content-Type"))
			}
		})
	}
}

func TestEncodeJSONResponseWithStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   interface{}
		wantCode   int
	}{
		{
			name:       "201 created",
			statusCode: 201,
			response:   map[string]string{"id": "123"},
			wantCode:   201,
		},
		{
			name:       "400 bad request",
			statusCode: 400,
			response:   map[string]string{"error": "bad request"},
			wantCode:   400,
		},
		{
			name:       "500 internal server error",
			statusCode: 500,
			response:   map[string]string{"error": "internal"},
			wantCode:   500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := EncodeJSONResponseWithStatus(w, tt.statusCode, tt.response)
			if err != nil {
				t.Errorf("EncodeJSONResponseWithStatus() error = %v", err)
			}
			if w.Code != tt.wantCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}

func TestEncodeJSONResponse_ValidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	resp := map[string]int{"count": 42}
	err := EncodeJSONResponse(w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if result["count"] != 42 {
		t.Errorf("count = %v, want 42", result["count"])
	}
}

func TestSetProtocolHeaders(t *testing.T) {
	tests := []struct {
		name            string
		contentType     string
		xAmzTarget      string
		smithyProtocol  string
		xAmznQueryMode  bool
		wantContentType string
		wantXAmzTarget  string
		wantSmithy      string
		wantQueryMode   string
	}{
		{
			name:            "all headers",
			contentType:     "application/json",
			xAmzTarget:      "SomeService.Operation",
			smithyProtocol:  "aws.rest_json",
			xAmznQueryMode:  true,
			wantContentType: "application/json",
			wantXAmzTarget:  "SomeService.Operation",
			wantSmithy:      "aws.rest_json",
			wantQueryMode:   "true",
		},
		{
			name:            "partial headers",
			contentType:     "application/xml",
			xAmzTarget:      "",
			smithyProtocol:  "",
			xAmznQueryMode:  false,
			wantContentType: "application/xml",
		},
		{
			name:            "empty all",
			contentType:     "",
			xAmzTarget:      "",
			smithyProtocol:  "",
			xAmznQueryMode:  false,
			wantContentType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SetProtocolHeaders(w, tt.contentType, tt.xAmzTarget, tt.smithyProtocol, tt.xAmznQueryMode)

			if tt.wantContentType != "" && w.Header().Get("Content-Type") != tt.wantContentType {
				t.Errorf("Content-Type = %v, want %v", w.Header().Get("Content-Type"), tt.wantContentType)
			}
			if tt.wantXAmzTarget != "" && w.Header().Get("X-Amz-Target") != tt.wantXAmzTarget {
				t.Errorf("X-Amz-Target = %v, want %v", w.Header().Get("X-Amz-Target"), tt.wantXAmzTarget)
			}
			if tt.wantSmithy != "" && w.Header().Get("smithy-protocol") != tt.wantSmithy {
				t.Errorf("smithy-protocol = %v, want %v", w.Header().Get("smithy-protocol"), tt.wantSmithy)
			}
			if tt.wantQueryMode != "" && w.Header().Get("X-Amzn-Query-Mode") != tt.wantQueryMode {
				t.Errorf("X-Amzn-Query-Mode = %v, want %v", w.Header().Get("X-Amzn-Query-Mode"), tt.wantQueryMode)
			}
		})
	}
}
