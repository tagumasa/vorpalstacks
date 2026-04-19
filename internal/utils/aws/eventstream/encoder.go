package eventstream

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

const (
	// HeaderContentType is the header name for content type.
	HeaderContentType = ":content-type"
	// HeaderMessageType is the header name for message type.
	HeaderMessageType = ":message-type"
	// HeaderEventType is the header name for event type.
	HeaderEventType = ":event-type"
	// ContentTypeOctetStream is the content type for binary data.
	ContentTypeOctetStream = "application/octet-stream"
	// MessageTypeEvent is the message type for events.
	MessageTypeEvent = "event"
	// preludeSize is the size of the event stream prelude header (total length + headers length).
	preludeSize = 8
)

// Header represents a header in an event stream message.
type Header struct {
	Name  string
	Value string
}

// Encoder encodes event stream messages.
type Encoder struct {
	w io.Writer
}

// NewEncoder creates a new Encoder that writes to the given writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode encodes an event with the given type, payload and headers.
func (e *Encoder) Encode(eventType string, payload []byte, headers []Header) error {
	var headerBuf bytes.Buffer
	for _, h := range headers {
		if err := e.encodeHeader(&headerBuf, h.Name, h.Value); err != nil {
			return err
		}
	}

	headersLen := headerBuf.Len()
	payloadLen := len(payload)

	totalLen := 4 + 4 + 4 + headersLen + payloadLen + 4

	prelude := make([]byte, preludeSize)
	binary.BigEndian.PutUint32(prelude[0:4], uint32(totalLen))
	binary.BigEndian.PutUint32(prelude[4:8], uint32(headersLen))
	preludeCrc := crc32.ChecksumIEEE(prelude)

	var message bytes.Buffer
	if err := binary.Write(&message, binary.BigEndian, uint32(totalLen)); err != nil {
		return fmt.Errorf("failed to write total length: %w", err)
	}
	if err := binary.Write(&message, binary.BigEndian, uint32(headersLen)); err != nil {
		return fmt.Errorf("failed to write headers length: %w", err)
	}
	if err := binary.Write(&message, binary.BigEndian, preludeCrc); err != nil {
		return fmt.Errorf("failed to write prelude CRC: %w", err)
	}
	message.Write(headerBuf.Bytes())
	message.Write(payload)

	messageCrc := crc32.ChecksumIEEE(message.Bytes())
	if err := binary.Write(&message, binary.BigEndian, messageCrc); err != nil {
		return fmt.Errorf("failed to write message CRC: %w", err)
	}

	_, err := e.w.Write(message.Bytes())
	return err
}

func (e *Encoder) encodeHeader(w io.Writer, name, value string) error {
	nameBytes := []byte(name)
	valueBytes := []byte(value)

	if _, err := w.Write([]byte{byte(len(nameBytes))}); err != nil {
		return fmt.Errorf("failed to write header name length: %w", err)
	}
	if _, err := w.Write(nameBytes); err != nil {
		return fmt.Errorf("failed to write header name: %w", err)
	}

	if _, err := w.Write([]byte{0x07}); err != nil {
		return fmt.Errorf("failed to write header type: %w", err)
	}

	if err := binary.Write(w, binary.BigEndian, uint16(len(valueBytes))); err != nil {
		return fmt.Errorf("failed to write header value length: %w", err)
	}
	if _, err := w.Write(valueBytes); err != nil {
		return fmt.Errorf("failed to write header value: %w", err)
	}

	return nil
}

// WriteEvent writes an event to the stream.
func (e *Encoder) WriteEvent(eventType, contentType string, payload []byte) error {
	headers := []Header{
		{HeaderMessageType, MessageTypeEvent},
		{HeaderEventType, eventType},
	}
	if contentType != "" {
		headers = append(headers, Header{HeaderContentType, contentType})
	}
	return e.Encode(eventType, payload, headers)
}

// WriteEndEvent writes an end event to the stream.
func (e *Encoder) WriteEndEvent() error {
	headers := []Header{
		{HeaderMessageType, MessageTypeEvent},
		{HeaderEventType, "End"},
	}
	return e.Encode("End", nil, headers)
}

// WriteInitialResponse writes an initial-response event to the stream.
func (e *Encoder) WriteInitialResponse(payload []byte) error {
	headers := []Header{
		{HeaderMessageType, MessageTypeEvent},
		{HeaderEventType, "initial-response"},
	}
	return e.Encode("initial-response", payload, headers)
}

// WriteErrorEvent writes an error event to the stream.
func (e *Encoder) WriteErrorEvent(errorCode, errorMessage string) error {
	payload := fmt.Sprintf(`{"errorCode":"%s","errorMessage":"%s"}`, errorCode, errorMessage)
	headers := []Header{
		{HeaderMessageType, "error"},
		{HeaderEventType, "error"},
		{HeaderContentType, "application/json"},
	}
	return e.Encode("error", []byte(payload), headers)
}
