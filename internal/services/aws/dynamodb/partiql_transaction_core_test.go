package dynamodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestSingletonCancellationReasonClassification(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		code    string
		message string
		cancels bool
	}{
		{
			name:    "duplicate-key insert cancels",
			err:     ErrConditionalCheckFailed,
			code:    "ConditionalCheckFailed",
			message: "The conditional request failed",
			cancels: true,
		},
		{
			name:    "clause type mismatch cancels",
			err:     ErrTypeMismatch,
			code:    "ValidationError",
			message: "Type mismatch for attribute to update.",
			cancels: true,
		},
		{
			name:    "wrapped clause mismatch still detected",
			err:     fmt.Errorf("apply clause: %w", ErrTypeMismatch),
			code:    "ValidationError",
			message: "Type mismatch for attribute to update.",
			cancels: true,
		},
		{
			name: "multi-item match cancels",
			err: NewAPIError("com.amazonaws.dynamodb.v20120810#ValidationException",
				"UPDATE statement must match exactly one item", http.StatusBadRequest),
			code:    "ValidationError",
			message: "UPDATE statement must match exactly one item",
			cancels: true,
		},
		{
			name:    "coral validation stays direct",
			err:     ErrInvalidParameter,
			cancels: false,
		},
		{
			name:    "unknown table stays direct",
			err:     ErrTableNotFound,
			cancels: false,
		},
		{
			name:    "storage fault stays direct",
			err:     errors.New("pebble: closed"),
			cancels: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reason, cancels := singletonCancellationReason(tc.err)
			if cancels != tc.cancels {
				t.Fatalf("cancels = %v, want %v", cancels, tc.cancels)
			}
			if !cancels {
				return
			}
			if reason.Code != tc.code {
				t.Fatalf("code = %q, want %q", reason.Code, tc.code)
			}
			if reason.Message != tc.message {
				t.Fatalf("message = %q, want %q", reason.Message, tc.message)
			}
		})
	}
}

// The wire contract for the envelope: unaffected statements carry the
// literal code "None" and omit the message entirely, while the failing
// statement's reason keeps both members.
func TestTransactionCanceledEnvelopeRendersNoneWithoutMessage(t *testing.T) {
	canceled := NewTransactionCanceledError("Transaction canceled", []CancellationReason{
		{Code: "None"},
		{Code: "ConditionalCheckFailed", Message: "The conditional request failed"},
	})

	var payload struct {
		Type                string `json:"__type"`
		CancellationReasons []struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"CancellationReasons"`
	}
	if err := json.Unmarshal([]byte(canceled.ToJSON()), &payload); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if !strings.HasSuffix(payload.Type, "TransactionCanceledException") {
		t.Fatalf("type = %q, want TransactionCanceledException", payload.Type)
	}
	if len(payload.CancellationReasons) != 2 {
		t.Fatalf("reasons = %d, want 2", len(payload.CancellationReasons))
	}
	if payload.CancellationReasons[0].Code != "None" {
		t.Fatalf("unaffected reason code = %q, want None", payload.CancellationReasons[0].Code)
	}
	if payload.CancellationReasons[0].Message != "" {
		t.Fatalf("unaffected reason must omit the message, got %q", payload.CancellationReasons[0].Message)
	}
	if payload.CancellationReasons[1].Code != "ConditionalCheckFailed" ||
		payload.CancellationReasons[1].Message != "The conditional request failed" {
		t.Fatalf("failing reason = %+v, want ConditionalCheckFailed with message", payload.CancellationReasons[1])
	}
	if strings.Contains(canceled.ToJSON(), `"Message":""`) {
		t.Fatalf("None reason must be serialised without a Message member")
	}
}
