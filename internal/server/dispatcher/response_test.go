package dispatcher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// delayedChunkReader holds the response's first byte past the test
// server's WriteTimeout, so a surviving write proves the streaming
// branch cleared the absolute deadline.
type delayedChunkReader struct {
	delay time.Duration
	once  sync.Once
}

func (r *delayedChunkReader) Read(p []byte) (int, error) {
	r.once.Do(func() { time.Sleep(r.delay) })
	return copy(p, "streamed-past-the-write-timeout"), io.EOF
}

// testStreamResponse is the minimal StreamableResponse: the event-stream
// content type and one delayed reader.
type testStreamResponse struct {
	reader io.Reader
}

func (t testStreamResponse) GetStream() io.Reader { return t.reader }

func (t testStreamResponse) GetStreamHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/vnd.amazon.eventstream")
	return headers
}

// TestWriteResponseStreamClearsWriteDeadline pins the streaming branch's
// deadline contract: an event-stream response outlives the server's
// WriteTimeout (Go sets it absolutely at request-head time — an
// HTTP/1.1 connection deadline, an HTTP/2 per-stream reset timer — and
// it does not extend with writes), so the branch clears it through the
// response controller before writing: a first write held past the
// deadline must still reach the client.
func TestWriteResponseStreamClearsWriteDeadline(t *testing.T) {
	d := &Dispatcher{}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d.writeResponse(w, r, nil, "", testStreamResponse{
			reader: &delayedChunkReader{delay: 600 * time.Millisecond},
		})
	}))
	srv.Config.WriteTimeout = 300 * time.Millisecond
	srv.Start()
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Type") != "application/vnd.amazon.eventstream" {
		t.Fatalf("content type = %q, want the stream headers", resp.Header.Get("Content-Type"))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if string(body) != "streamed-past-the-write-timeout" {
		t.Fatalf("body = %q — the write held past the WriteTimeout did not reach the client", body)
	}
}
