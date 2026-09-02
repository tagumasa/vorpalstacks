package headerorder

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
)

// headerMapWithHost copies a request header map and injects the
// promoted Host header, mirroring what the WAF enforcement planes pass
// to FromContext.
func headerMapWithHost(headers map[string][]string, host string) map[string][]string {
	out := make(map[string][]string, len(headers)+1)
	for name, values := range headers {
		out[name] = values
	}
	out["Host"] = []string{host}
	return out
}

// serveCaptured runs an HTTP server over a capture-wrapped listener
// and records the order each request observes; the return value is the
// listener address.
func serveCaptured(t *testing.T) (string, *capturedOrders) {
	t.Helper()
	orders := &capturedOrders{}
	raw, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			names, ok := FromContext(r.Context(), headerMapWithHost(r.Header, r.Host))
			orders.add(names, ok)
			w.WriteHeader(http.StatusOK)
		}),
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			if cc, ok := c.(*captureConn); ok {
				return NewContext(ctx, cc.State())
			}
			return ctx
		},
	}
	go func() { _ = srv.Serve(NewListener(raw)) }()
	t.Cleanup(func() { _ = srv.Close() })
	return raw.Addr().String(), orders
}

type capturedOrders struct {
	mu     sync.Mutex
	orders [][]string
	found  []bool
}

func (c *capturedOrders) add(names []string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders = append(c.orders, names)
	c.found = append(c.found, ok)
}

func (c *capturedOrders) snapshot() ([][]string, []bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	orders := make([][]string, len(c.orders))
	copy(orders, c.orders)
	found := make([]bool, len(c.found))
	copy(found, c.found)
	return orders, found
}

// sendRaw writes a hand-crafted HTTP/1.1 exchange over one connection.
func sendRaw(t *testing.T, addr string, requests []string) {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for _, request := range requests {
		if _, err := conn.Write([]byte(request)); err != nil {
			t.Fatal(err)
		}
		// Consume the response head so keep-alive sequencing works.
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(line) == "" {
				break
			}
		}
	}
}

func TestHeaderOrderCapturedOnTheWire(t *testing.T) {
	addr, orders := serveCaptured(t)

	sendRaw(t, addr, []string{
		"GET /one HTTP/1.1\r\nHost: example.test\r\nZz-Last: 1\r\nAa-First: 2\r\nMm-Middle: 3\r\n\r\n",
	})
	got, found := orders.snapshot()
	if len(got) != 1 || !found[0] {
		t.Fatalf("captured orders = %v (found %v), want one captured request", got, found)
	}
	want := []string{"Host", "Zz-Last", "Aa-First", "Mm-Middle"}
	if strings.Join(got[0], ",") != strings.Join(want, ",") {
		t.Fatalf("order = %v, want %v", got[0], want)
	}
}

func TestHeaderOrderFollowsKeepAliveRequests(t *testing.T) {
	addr, orders := serveCaptured(t)

	sendRaw(t, addr, []string{
		"POST /one HTTP/1.1\r\nHost: example.test\r\nContent-Length: 5\r\nX-Before-Body: a\r\n\r\nhello",
		"GET /two HTTP/1.1\r\nHost: example.test\r\nA-Second-Request: b\r\n\r\n",
	})
	got, found := orders.snapshot()
	if len(got) != 2 || !found[0] || !found[1] {
		t.Fatalf("captured orders = %v (found %v), want two captured requests", got, found)
	}
	if strings.Join(got[0], ",") != "Host,Content-Length,X-Before-Body" {
		t.Fatalf("first order = %v", got[0])
	}
	if strings.Join(got[1], ",") != "Host,A-Second-Request" {
		t.Fatalf("second order = %v", got[1])
	}
}

// sendRawPipelined writes every request before reading any response,
// so the server receives pipelined request heads back to back.
func sendRawPipelined(t *testing.T, addr string, requests []string) {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var wire strings.Builder
	for _, request := range requests {
		wire.WriteString(request)
	}
	if _, err := conn.Write([]byte(wire.String())); err != nil {
		t.Fatal(err)
	}
	reader := bufio.NewReader(conn)
	for range requests {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(line) == "" {
				break
			}
		}
	}
}

func TestHeaderOrderPipelinedRequestsKeepOwnOrder(t *testing.T) {
	addr, orders := serveCaptured(t)

	sendRawPipelined(t, addr, []string{
		"GET /a HTTP/1.1\r\nHost: example.test\r\nX-A: 1\r\nX-B: 2\r\n\r\n",
		"GET /b HTTP/1.1\r\nHost: example.test\r\nX-C: 3\r\n\r\n",
	})
	got, found := orders.snapshot()
	if len(got) != 2 || !found[0] || !found[1] {
		t.Fatalf("captured orders = %v (found %v), want two captured requests", got, found)
	}
	if strings.Join(got[0], ",") != "Host,X-A,X-B" {
		t.Fatalf("first request order = %v, want Host,X-A,X-B", got[0])
	}
	if strings.Join(got[1], ",") != "Host,X-C" {
		t.Fatalf("second request order = %v, want Host,X-C", got[1])
	}
}

func TestHeaderOrderDisabledOnHTTP2Preface(t *testing.T) {
	// The HTTP/2 connection preface disables capture; the server cannot
	// serve the stream, so only the parser's behaviour matters here —
	// feeding the preface must leave the connection state empty.
	parser := newHeadParser()
	state := &ConnState{}
	parser.feed([]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"), state)
	if _, ok := FromContext(NewContext(context.Background(), state), nil); ok {
		t.Fatal("an HTTP/2 preface must not produce a captured order")
	}
}

func TestHeaderOrderFallsBackWithoutCapture(t *testing.T) {
	if _, ok := FromContext(context.Background(), nil); ok {
		t.Fatal("a context without capture state must report no order")
	}
	if _, ok := FromContext(NewContext(context.Background(), &ConnState{}), nil); ok {
		t.Fatal("a connection with no completed head must report no order")
	}
}

func TestFromContextPairsRequestByHeaderPairs(t *testing.T) {
	state := &ConnState{}
	state.record([]headerPair{{"Host", "a.test"}, {"X-A", "1"}, {"X-B", "2"}})
	state.record([]headerPair{{"Host", "a.test"}, {"X-C", "3"}})

	// The first request's map pairs with the first head even though the
	// second head has already been recorded (pipelining).
	names, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"a.test"}, "X-A": {"1"}, "X-B": {"2"},
	})
	if !ok || strings.Join(names, ",") != "Host,X-A,X-B" {
		t.Fatalf("first request order = %v (ok %v), want Host,X-A,X-B", names, ok)
	}
	names, ok = FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"a.test"}, "X-C": {"3"},
	})
	if !ok || strings.Join(names, ",") != "Host,X-C" {
		t.Fatalf("second request order = %v (ok %v), want Host,X-C", names, ok)
	}
	// A header map that matches no queued head reports no order and
	// leaves the queue intact for its real consumer.
	state.record([]headerPair{{"Host", "a.test"}, {"X-D", "4"}})
	if _, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"a.test"}, "X-Missing": {"1"},
	}); ok {
		t.Fatal("a header map matching no head must report no order")
	}
	names, ok = FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"a.test"}, "X-D": {"4"},
	})
	if !ok || strings.Join(names, ",") != "Host,X-D" {
		t.Fatalf("third request order = %v (ok %v), want Host,X-D", names, ok)
	}
	// A map without a Host header cannot pair with any wire head.
	if _, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"X-D": {"4"},
	}); ok {
		t.Fatal("a header map without Host must report no order")
	}
}

// An earlier head whose consumer never takes it — a request an
// unprotected virtual host served without reading the order — must not
// pair with a later request that carries the same header names in a
// different order: the Host value distinguishes the two heads.
func TestFromContextSkipsUnconsumedHeadOfOtherHost(t *testing.T) {
	state := &ConnState{}
	state.record([]headerPair{{"Host", "other.test"}, {"X-A", "1"}, {"X-B", "2"}})
	state.record([]headerPair{{"Host", "example.test"}, {"X-B", "2"}, {"X-A", "1"}})

	names, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"example.test"}, "X-A": {"1"}, "X-B": {"2"},
	})
	if !ok || strings.Join(names, ",") != "Host,X-B,X-A" {
		t.Fatalf("order = %v (ok %v), want the second head's Host,X-B,X-A", names, ok)
	}
}

// Two pipelined requests with the same header names but different
// Host values must each pair with their own head: an earlier head
// whose consumer never takes it (an unprotected virtual host) must not
// pair with the later request.
func TestFromContextPairsSameNamesByValue(t *testing.T) {
	addr, orders := serveCaptured(t)

	// The first request targets a different virtual host and order of
	// the same two custom headers; the pipelined second request must
	// observe its own order, not the first head's.
	sendRawPipelined(t, addr, []string{
		"GET /unprotected HTTP/1.1\r\nHost: other.test\r\nX-A: 1\r\nX-B: 2\r\n\r\n",
		"GET /protected HTTP/1.1\r\nHost: example.test\r\nX-B: 2\r\nX-A: 1\r\n\r\n",
	})
	got, found := orders.snapshot()
	if len(got) != 2 || !found[0] || !found[1] {
		t.Fatalf("captured orders = %v (found %v), want two captured requests", got, found)
	}
	if strings.Join(got[0], ",") != "Host,X-A,X-B" {
		t.Fatalf("first request order = %v, want Host,X-A,X-B", got[0])
	}
	if strings.Join(got[1], ",") != "Host,X-B,X-A" {
		t.Fatalf("second request order = %v, want Host,X-B,X-A", got[1])
	}
}

// Two queued heads whose (name, value) pair multisets are identical —
// only the header order differs — cannot be attributed to one request
// when the earlier head's consumer never took it. The pairing must fall
// back (report no order) instead of guessing, and the queue must stay
// intact so later requests with distinct pairs still pair.
func TestFromContextFallsBackWhenPairAmbiguous(t *testing.T) {
	state := &ConnState{}
	state.record([]headerPair{{"Host", "example.test"}, {"X-A", "1"}, {"X-B", "2"}})
	state.record([]headerPair{{"Host", "example.test"}, {"X-B", "2"}, {"X-A", "1"}})

	if _, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"example.test"}, "X-A": {"1"}, "X-B": {"2"},
	}); ok {
		t.Fatal("a header map matching two queued heads must report no order")
	}
	if _, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"example.test"}, "X-B": {"2"}, "X-A": {"1"},
	}); ok {
		t.Fatal("a duplicate head must not be attributed to the second map either")
	}
	// The duplicate heads stay queued but do not block a later request
	// whose pairs differ from both.
	state.record([]headerPair{{"Host", "example.test"}, {"X-C", "3"}})
	names, ok := FromContext(NewContext(context.Background(), state), map[string][]string{
		"Host": {"example.test"}, "X-C": {"3"},
	})
	if !ok || strings.Join(names, ",") != "Host,X-C" {
		t.Fatalf("distinct-pair order = %v (ok %v), want Host,X-C", names, ok)
	}
}

// A pipelined request whose consumer never reads the order (an
// unprotected route) followed by a protected request carrying the same
// (name, value) pairs in a different order: the two queued heads are
// indistinguishable, so the protected request must fall back instead of
// observing the earlier request's order.
func TestFromContextUnconsumedDuplicateHeadFallsBack(t *testing.T) {
	orders := &capturedOrders{}
	raw, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/protected" {
				names, ok := FromContext(r.Context(), headerMapWithHost(r.Header, r.Host))
				orders.add(names, ok)
			}
			w.WriteHeader(http.StatusOK)
		}),
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			if cc, ok := c.(*captureConn); ok {
				return NewContext(ctx, cc.State())
			}
			return ctx
		},
	}
	go func() { _ = srv.Serve(NewListener(raw)) }()
	t.Cleanup(func() { _ = srv.Close() })

	sendRawPipelined(t, raw.Addr().String(), []string{
		"GET /plain HTTP/1.1\r\nHost: example.test\r\nX-A: 1\r\nX-B: 2\r\n\r\n",
		"GET /protected HTTP/1.1\r\nHost: example.test\r\nX-B: 2\r\nX-A: 1\r\n\r\n",
	})
	got, found := orders.snapshot()
	if len(got) != 1 || found[0] {
		t.Fatalf("protected request order = %v (found %v), want fallback with no order", got, found)
	}
}

// A request carrying repeated header lines (two Cookie headers) must
// still pair with its captured head: every value of the map slice is
// one wire line.
func TestFromContextPairsRepeatedHeaderValues(t *testing.T) {
	addr, orders := serveCaptured(t)

	sendRaw(t, addr, []string{
		"GET /cookies HTTP/1.1\r\nHost: example.test\r\nCookie: a=1\r\nCookie: b=2\r\nX-After: 3\r\n\r\n",
	})
	got, found := orders.snapshot()
	if len(got) != 1 || !found[0] {
		t.Fatalf("captured orders = %v (found %v), want one captured request", got, found)
	}
	if strings.Join(got[0], ",") != "Host,Cookie,Cookie,X-After" {
		t.Fatalf("order = %v, want Host,Cookie,Cookie,X-After", got[0])
	}
}
