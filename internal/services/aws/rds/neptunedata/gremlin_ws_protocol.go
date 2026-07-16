package neptunedata

import (
	"encoding/json"
	"fmt"
)

// TinkerPop Gremlin Server status codes used in WebSocket response messages.
const (
	statusSuccess          = 200
	statusPartialContent   = 206
	statusUnauthorized     = 401
	statusAuthenticate     = 407
	statusMalformedRequest = 498
	statusServerError      = 500
	statusScriptError      = 598
	statusTimeout          = 599
)

// GremlinRequest represents an incoming TinkerPop Gremlin Server request
// message received over WebSocket. The op field determines the operation:
// "eval" executes a Gremlin traversal, "authentication" handles SASL auth,
// and "close" terminates a session.
type GremlinRequest struct {
	RequestID string      `json:"requestId"`
	Op        string      `json:"op"`
	Processor string      `json:"processor"`
	Args      GremlinArgs `json:"args"`
}

// GremlinArgs holds the arguments of a TinkerPop request. The Gremlin field
// contains the traversal string to evaluate. Aliases and Bindings support
// session-based evaluations. Session is set when processor is "session".
type GremlinArgs struct {
	Gremlin         string                 `json:"gremlin"`
	Language        string                 `json:"language"`
	Aliases         map[string]string      `json:"aliases"`
	Bindings        map[string]interface{} `json:"bindings"`
	Session         string                 `json:"session"`
	Manage          string                 `json:"manage"`
	EvalTimeout     int64                  `json:"evaluationTimeout"`
	SASLMechanisms  string                 `json:"saslMechanisms"`
	SASL            string                 `json:"sasl"`
	TraversalSource string                 `json:"traversalSource"`
}

// GremlinResponse represents a TinkerPop Gremlin Server response sent over
// WebSocket. Status.Code follows TinkerPop conventions (200=success, etc.).
type GremlinResponse struct {
	RequestID string    `json:"requestId"`
	Status    StatusMsg `json:"status"`
	Result    ResultMsg `json:"result"`
}

// StatusMsg carries the TinkerPop status code, human-readable message, and
// optional attributes map.
type StatusMsg struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Attributes map[string]interface{} `json:"attributes"`
}

// ResultMsg carries the evaluation result data and optional metadata.
type ResultMsg struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

// gremlinSession tracks per-session state for TinkerPop session processor
// connections. Sessions maintain aliases and bindings across eval calls.
type gremlinSession struct {
	id       string
	aliases  map[string]string
	bindings map[string]interface{}
}

// newGremlinSuccessResponse creates a 200 OK response with the given data.
func newGremlinSuccessResponse(requestID string, data interface{}) GremlinResponse {
	return GremlinResponse{
		RequestID: requestID,
		Status: StatusMsg{
			Code:       statusSuccess,
			Message:    "",
			Attributes: map[string]interface{}{},
		},
		Result: ResultMsg{
			Data: data,
			Meta: map[string]interface{}{},
		},
	}
}

// newGremlinAuthResponse creates a 407 Authenticate response requesting the
// client to provide SASL credentials.

// newGremlinErrorResponse creates a response with the given error status code
// and message.
func newGremlinErrorResponse(requestID string, code int, message string) GremlinResponse {
	return GremlinResponse{
		RequestID: requestID,
		Status: StatusMsg{
			Code:       code,
			Message:    message,
			Attributes: map[string]interface{}{},
		},
		Result: ResultMsg{
			Data: nil,
			Meta: map[string]interface{}{},
		},
	}
}

// decodeGremlinRequest parses a TinkerPop WebSocket message. Messages may
// include a binary MIME header (1 byte length + MIME string) followed by the
// JSON payload, or be plain JSON when sent as text frames. Returns the parsed
// request and the detected MIME type.
func decodeGremlinRequest(raw []byte) (*GremlinRequest, string, error) {
	payload, mimeType, err := stripMimeHeader(raw)
	if err != nil {
		return nil, "", fmt.Errorf("failed to strip mime header: %w", err)
	}

	var req GremlinRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, mimeType, fmt.Errorf("invalid JSON: %w", err)
	}
	return &req, mimeType, nil
}

// stripMimeHeader strips the optional binary MIME header from a TinkerPop
// message. If the first byte is 0x00 (no MIME) or the message starts with
// '{', it is treated as plain JSON. Otherwise the first byte is the MIME
// string length, followed by that many bytes of MIME type, then the payload.
func stripMimeHeader(raw []byte) ([]byte, string, error) {
	if len(raw) == 0 {
		return nil, "", fmt.Errorf("empty message")
	}

	// Plain JSON text frame (no binary MIME header)
	if raw[0] == '{' {
		return raw, "", nil
	}

	// 0x00 means no MIME type
	if raw[0] == 0x00 {
		return raw[1:], "", nil
	}

	mimeLen := int(raw[0])
	if len(raw) < 1+mimeLen {
		return nil, "", fmt.Errorf("message too short for MIME header: need %d bytes, have %d", 1+mimeLen, len(raw))
	}

	mimeType := string(raw[1 : 1+mimeLen])
	payload := raw[1+mimeLen:]
	return payload, mimeType, nil
}

// encodeResponseToJSON serialises a GremlinResponse to JSON bytes, ready to
// be sent as a WebSocket text message.
func encodeResponseToJSON(resp GremlinResponse) ([]byte, error) {
	return json.Marshal(resp)
}
