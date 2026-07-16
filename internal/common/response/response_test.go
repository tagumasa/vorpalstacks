package response

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyResponse(t *testing.T) {
	resp := EmptyResponse()
	assert.NotNil(t, resp)
	assert.Empty(t, resp)
}

func TestStreamableResponse(t *testing.T) {
	data := []byte("test payload")
	s := &mockReader{data: data}

	assert.Implements(t, (*io.Reader)(nil), s.GetStream())

	h := s.GetStreamHeaders()
	assert.Equal(t, "application/octet-stream", h.Get("Content-Type"))
}

func TestStatusCodeResponse(t *testing.T) {
	s := &mockReader{data: []byte("x"), statusCode: 202}
	assert.Equal(t, 202, s.GetStreamStatusCode())

	s2 := &mockReader{data: []byte("x")}
	assert.Equal(t, 200, s2.GetStreamStatusCode())
}

type mockReader struct {
	data       []byte
	statusCode int
}

func (m *mockReader) GetStream() io.Reader {
	return bytes.NewReader(m.data)
}

func (m *mockReader) GetStreamHeaders() http.Header {
	h := make(http.Header)
	h.Set("Content-Type", "application/octet-stream")
	return h
}

func (m *mockReader) GetStreamStatusCode() int {
	if m.statusCode != 0 {
		return m.statusCode
	}
	return 200
}
