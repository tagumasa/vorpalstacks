package eventstream

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	t.Run("simple event", func(t *testing.T) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)

		err := enc.Encode("PayloadChunk", []byte("hello"), []Header{
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

		err := enc.Encode("End", nil, []Header{
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

func TestWriteErrorEvent(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	err := enc.WriteErrorEvent("InvalidArgument", "bad param")
	require.NoError(t, err)

	data := buf.Bytes()
	assert.Contains(t, string(data), "InvalidArgument")
	assert.Contains(t, string(data), "bad param")
}

func TestConstants(t *testing.T) {
	assert.Equal(t, ":content-type", HeaderContentType)
	assert.Equal(t, ":message-type", HeaderMessageType)
	assert.Equal(t, ":event-type", HeaderEventType)
	assert.Equal(t, "application/octet-stream", ContentTypeOctetStream)
	assert.Equal(t, "event", MessageTypeEvent)
}
