package eventstream

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	t.Run("simple event", func(t *testing.T) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)

		err := enc.Encode([]byte("hello"), []Header{
			{Name: HeaderContentType, Value: ContentTypeOctetStream},
		})
		require.NoError(t, err)

		data := buf.Bytes()
		assert.GreaterOrEqual(t, len(data), 20)

		totalLen := binary.BigEndian.Uint32(data[0:4])
		assert.Equal(t, uint32(len(data)), totalLen)

		headersLen := binary.BigEndian.Uint32(data[4:8])
		assert.Greater(t, headersLen, uint32(0))

		preludeCrc := binary.BigEndian.Uint32(data[8:12])
		assert.NotZero(t, preludeCrc)

		trailingCrc := binary.BigEndian.Uint32(data[len(data)-4:])
		assert.NotZero(t, trailingCrc)
	})

	t.Run("empty payload", func(t *testing.T) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)

		err := enc.Encode(nil, []Header{
			{Name: HeaderMessageType, Value: MessageTypeEvent},
			{Name: HeaderEventType, Value: "End"},
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, buf.Len(), 20)
	})
}

func TestWriteEvent(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.WriteEvent("PayloadChunk", ContentTypeOctetStream, []byte("data"))
	require.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestWriteEndEvent(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.WriteEndEvent()
	require.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

func TestWriteInitialResponse(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.WriteInitialResponse([]byte(`{"status":"ok"}`))
	require.NoError(t, err)
	assert.Greater(t, buf.Len(), 0)
}

// decodeFrameHeaders parses the string-valued headers of one encoded
// frame so tests can assert the classification the AWS SDKs dispatch on
// (the message-type value itself, not merely the header names).
func decodeFrameHeaders(t *testing.T, data []byte) map[string]string {
	t.Helper()
	headersLen := int(binary.BigEndian.Uint32(data[4:8]))
	pos := 12 // prelude (8 bytes) + prelude CRC (4 bytes)
	end := pos + headersLen
	out := make(map[string]string)
	for pos < end {
		nameLen := int(data[pos])
		pos++
		name := string(data[pos : pos+nameLen])
		pos += nameLen
		typ := data[pos]
		pos++
		require.Equal(t, byte(0x07), typ, "header %s: expected string value type", name)
		valLen := int(binary.BigEndian.Uint16(data[pos : pos+2]))
		pos += 2
		out[name] = string(data[pos : pos+valLen])
		pos += valLen
	}
	return out
}

func TestWriteErrorEvent(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.WriteErrorEvent("InvalidArgumentException", "bad param")
	require.NoError(t, err)

	// The AWS SDKs' typed exception dispatch runs only on a message-type
	// of "exception": the frame must carry that classification, name the
	// modelled shape in :exception-type, and never identify the error as
	// an event.
	headers := decodeFrameHeaders(t, buf.Bytes())
	assert.Equal(t, "exception", headers[HeaderMessageType])
	assert.Equal(t, "InvalidArgumentException", headers[HeaderExceptionType])
	_, hasEventType := headers[HeaderEventType]
	assert.False(t, hasEventType, "error frame must not carry :event-type")
}

func TestConstants(t *testing.T) {
	assert.Equal(t, ":content-type", HeaderContentType)
	assert.Equal(t, ":message-type", HeaderMessageType)
	assert.Equal(t, ":event-type", HeaderEventType)
	assert.Equal(t, ":exception-type", HeaderExceptionType)
	assert.Equal(t, "application/octet-stream", ContentTypeOctetStream)
	assert.Equal(t, "event", MessageTypeEvent)
	assert.Equal(t, "exception", MessageTypeException)
}

func TestWriteErrorEventJSONSafety(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	// The message carries quotes, a backslash and a control byte — any
	// format-string assembly would emit invalid JSON.
	message := `stream "shard-1" \ ` + "\x01"
	err := enc.WriteErrorEvent("ResourceNotFoundException", message)
	require.NoError(t, err)

	want, err := json.Marshal(struct {
		Message string `json:"message"`
	}{message})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), string(want))
}
