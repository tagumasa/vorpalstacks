// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
)

// eventstreamMessage represents a single message in the AWS event stream binary encoding.
// Wire format (all integers big-endian):
//
//	[4] total_length  — prelude (8) + headers + payload + message_crc (4)
//	[4] headers_length
//	[4] prelude_crc    — CRC32 of the 8 prelude bytes
//	[n] headers
//	[m] payload
//	[4] message_crc   — CRC32 of everything from byte 0 to end of payload
//
// See: https://smithy.io/2.0/aws/amazon-eventstream.html
type eventstreamMessage struct {
	headers map[string]string
	payload []byte
}

// encode serialises the message into the AWS event stream binary format and writes it to w.
func (m *eventstreamMessage) encode(w io.Writer) error {
	headerBytes := encodeHeaders(m.headers)
	headersLen := uint32(len(headerBytes))
	payloadLen := uint32(len(m.payload))
	totalLen := uint32(16) + headersLen + payloadLen // 8 prelude + 4 prelude CRC + headers + payload + 4 message CRC

	buf := make([]byte, 0, totalLen)
	buf = binary.BigEndian.AppendUint32(buf, totalLen)
	buf = binary.BigEndian.AppendUint32(buf, headersLen)

	preludeCRC := crc32.ChecksumIEEE(buf[:8])
	buf = binary.BigEndian.AppendUint32(buf, preludeCRC)

	buf = append(buf, headerBytes...)
	buf = append(buf, m.payload...)

	msgCRC := crc32.ChecksumIEEE(buf)
	buf = binary.BigEndian.AppendUint32(buf, msgCRC)

	_, err := w.Write(buf)
	return err
}

// encodeHeaders serialises a set of string-valued headers into the event stream binary header format.
// Each header is encoded as:
//
//	[1] header_name_length
//	[n] header_name (UTF-8)
//	[1] header_value_type = 7 (string)
//	[2] string_length (big-endian)
//	[n] string_value (UTF-8)
func encodeHeaders(headers map[string]string) []byte {
	// Pre-allocate a reasonable size: average ~40 bytes per header.
	var buf []byte
	for k, v := range headers {
		buf = append(buf, byte(len(k)))
		buf = append(buf, k...)
		buf = append(buf, byte(7)) // string value type indicator
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(v)))
		buf = append(buf, v...)
	}
	return buf
}

// InvokeResponseStreamWriter provides methods for writing Lambda response stream events
// in the AWS event stream binary encoding format.
// It is used by InvokeWithResponseStream to stream payload chunks and completion events
// back to the caller over an HTTP response body.
type InvokeResponseStreamWriter struct {
	w io.Writer
}

// NewInvokeResponseStreamWriter creates a new stream writer targeting the provided io.Writer.
func NewInvokeResponseStreamWriter(w io.Writer) *InvokeResponseStreamWriter {
	return &InvokeResponseStreamWriter{w: w}
}

// WritePayloadChunk writes a PayloadChunk event containing the given payload bytes.
// The TS SDK deserialiser reads this as a PayloadChunk event with the raw bytes in the Payload field.
func (s *InvokeResponseStreamWriter) WritePayloadChunk(payload []byte) error {
	return (&eventstreamMessage{
		headers: map[string]string{
			":event-type":   "PayloadChunk",
			":content-type": "application/octet-stream",
			":message-type": "event",
		},
		payload: payload,
	}).encode(s.w)
}

// WriteInvokeComplete writes an InvokeComplete event signalling the end of the response stream.
// The payload is a JSON object containing LogResult and, for error responses, ErrorCode and ErrorDetails.
func (s *InvokeResponseStreamWriter) WriteInvokeComplete(statusCode int, executedVersion string, logResult string) error {
	body := map[string]string{}
	if logResult != "" {
		body["LogResult"] = logResult
	}
	if statusCode >= 400 {
		body["ErrorCode"] = fmt.Sprintf("%d", statusCode)
		body["ErrorDetails"] = fmt.Sprintf("Invocation failed with status %d", statusCode)
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal invoke complete: %w", err)
	}
	return (&eventstreamMessage{
		headers: map[string]string{
			":event-type":   "InvokeComplete",
			":content-type": "application/json",
			":message-type": "event",
		},
		payload: payload,
	}).encode(s.w)
}

// WriteInvokeCompleteError writes an InvokeComplete event with an :exception-type header.
// The TS SDK deserialiser treats messages with :message-type=event and :exception-type as
// exception events, and throws an error with the given type and message.
func (s *InvokeResponseStreamWriter) WriteInvokeCompleteError(errorType string, errorMessage string) error {
	body := map[string]string{
		"ErrorCode":    errorType,
		"ErrorDetails": errorMessage,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal invoke complete error: %w", err)
	}
	return (&eventstreamMessage{
		headers: map[string]string{
			":event-type":     "InvokeComplete",
			":content-type":   "application/json",
			":message-type":   "event",
			":exception-type": errorType,
		},
		payload: payload,
	}).encode(s.w)
}
