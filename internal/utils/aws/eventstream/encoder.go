package eventstream

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	// HeaderExceptionType is the header name carrying the modelled
	// exception shape name on an exception-classified message.
	HeaderExceptionType = ":exception-type"
	// ContentTypeOctetStream is the content type for binary data.
	ContentTypeOctetStream = "application/octet-stream"
	// MessageTypeEvent is the message type for events.
	MessageTypeEvent = "event"
	// MessageTypeException is the message type for modelled operation
	// exceptions. The AWS SDKs' typed exception dispatch runs only on
	// this classification; an "error"-classified frame is read through
	// the transport error branch, which consults only :error-code and
	// :error-message.
	MessageTypeException = "exception"
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

// Encode writes a frame with the given payload and headers. The event or
// exception type travels in the headers the caller builds — the Write*
// helpers construct them — so the frame carries no type outside them.
func (e *Encoder) Encode(payload []byte, headers []Header) error {
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
	return e.Encode(payload, headers)
}

// WriteEndEvent writes an end event to the stream.
func (e *Encoder) WriteEndEvent() error {
	headers := []Header{
		{HeaderMessageType, MessageTypeEvent},
		{HeaderEventType, "End"},
	}
	return e.Encode(nil, headers)
}

// WriteInitialResponse writes an initial-response event to the stream.
func (e *Encoder) WriteInitialResponse(payload []byte) error {
	headers := []Header{
		{HeaderMessageType, MessageTypeEvent},
		{HeaderEventType, "initial-response"},
	}
	return e.Encode(payload, headers)
}

// WriteErrorEvent writes a modelled-exception frame. The AWS SDKs' typed
// exception dispatch runs only on a message-type of "exception": the frame
// names the modelled shape in :exception-type and carries a marshalled
// {"message": ...} payload the shape's document decodes — the Kinesis
// exception shapes carry a message member. An error-classified frame is
// never read this way: that branch consults only :error-code and
// :error-message and reports an error named under :exception-type as an
// unclassified UnknownError. A frame that identifies the error under
// :event-type, or omits :exception-type, is a deserialisation error at the
// client, not a typed exception. The payload is marshalled, never
// formatted: an error message carrying quotes, backslashes or control
// bytes must still decode as JSON at the client.
func (e *Encoder) WriteErrorEvent(errorCode, errorMessage string) error {
	payload, err := json.Marshal(struct {
		Message string `json:"message"`
	}{errorMessage})
	if err != nil {
		return err
	}
	headers := []Header{
		{HeaderMessageType, MessageTypeException},
		{HeaderExceptionType, errorCode},
		{HeaderContentType, "application/json"},
	}
	return e.Encode(payload, headers)
}
