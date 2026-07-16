package sts

import (
	"context"
	"encoding/base64"
	"testing"

	"vorpalstacks/internal/common/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeAuthorizationMessage(t *testing.T) {
	svc := &STSService{}

	tests := []struct {
		name           string
		encodedMessage string
		expectError    bool
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "valid base64 JSON message",
			encodedMessage: base64.StdEncoding.EncodeToString([]byte(`{"Decision":"Allow","Action":"s3:GetObject"}`)),
			expectError:    false,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				decoded, ok := result["DecodedMessage"].(string)
				require.True(t, ok, "DecodedMessage should be a string")
				assert.Contains(t, decoded, "Decision")
				assert.Contains(t, decoded, "Allow")
			},
		},
		{
			name:           "valid base64 plain text message",
			encodedMessage: base64.StdEncoding.EncodeToString([]byte("Access denied for resource")),
			expectError:    false,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				decoded, ok := result["DecodedMessage"].(string)
				require.True(t, ok, "DecodedMessage should be a string")
				assert.Contains(t, decoded, "Access denied")
			},
		},
		{
			name:           "URL-safe base64 encoded message",
			encodedMessage: base64.URLEncoding.EncodeToString([]byte(`{"Reason":"Unauthorized"}`)),
			expectError:    false,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				decoded, ok := result["DecodedMessage"].(string)
				require.True(t, ok, "DecodedMessage should be a string")
				assert.Contains(t, decoded, "Reason")
			},
		},
		{
			name:           "empty encoded message",
			encodedMessage: "",
			expectError:    true,
			validateResult: nil,
		},
		{
			name:           "invalid base64 message",
			encodedMessage: "not-valid-base64!!!",
			expectError:    true,
			validateResult: nil,
		},
		{
			name:           "base64 encoded empty JSON object",
			encodedMessage: base64.StdEncoding.EncodeToString([]byte("{}")),
			expectError:    false,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				decoded, ok := result["DecodedMessage"].(string)
				require.True(t, ok, "DecodedMessage should be a string")
				assert.Equal(t, "{}", decoded)
			},
		},
		{
			name:           "base64 encoded JSON array",
			encodedMessage: base64.StdEncoding.EncodeToString([]byte(`["action1","action2"]`)),
			expectError:    false,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				decoded, ok := result["DecodedMessage"].(string)
				require.True(t, ok, "DecodedMessage should be a string")
				assert.Contains(t, decoded, "action1")
			},
		},
		{
			name:           "complex AWS authorization message",
			encodedMessage: base64.StdEncoding.EncodeToString([]byte(`{"explicitDeny":false,"matchedStatements":{"items":[]},"decision":"Allow"}`)),
			expectError:    false,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				decoded, ok := result["DecodedMessage"].(string)
				require.True(t, ok, "DecodedMessage should be a string")
				assert.Contains(t, decoded, "decision")
				assert.Contains(t, decoded, "Allow")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &request.RequestContext{AccountID: "123456789012"}
			req := &request.ParsedRequest{
				Parameters: map[string]interface{}{
					"EncodedMessage": tt.encodedMessage,
				},
			}

			result, err := svc.DecodeAuthorizationMessage(context.Background(), ctx, req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, ErrInvalidEncodedMessage, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				resp, ok := result.(map[string]interface{})
				require.True(t, ok, "result should be a map")
				tt.validateResult(t, resp)
			}
		})
	}
}

func TestDecodeAuthorizationMessage_Base64Variants(t *testing.T) {
	svc := &STSService{}

	testCases := []struct {
		name    string
		message string
		encode  func([]byte) string
	}{
		{
			name:    "standard base64",
			message: `{"test":"value"}`,
			encode:  func(b []byte) string { return base64.StdEncoding.EncodeToString(b) },
		},
		{
			name:    "URL-safe base64",
			message: `{"test":"value"}`,
			encode:  func(b []byte) string { return base64.URLEncoding.EncodeToString(b) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &request.RequestContext{AccountID: "123456789012"}
			encoded := tc.encode([]byte(tc.message))
			req := &request.ParsedRequest{
				Parameters: map[string]interface{}{
					"EncodedMessage": encoded,
				},
			}

			result, err := svc.DecodeAuthorizationMessage(context.Background(), ctx, req)
			require.NoError(t, err)

			resp := result.(map[string]interface{})
			decoded := resp["DecodedMessage"].(string)
			assert.Contains(t, decoded, "test")
		})
	}
}
