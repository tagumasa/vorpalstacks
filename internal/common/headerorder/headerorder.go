// Package headerorder preserves the wire order of HTTP/1.1 request
// header names across the HTTP server's parsing, which stores headers
// in an orderless map. Consumers that must match on the order — the
// AWS WAF HeaderOrder request component is the driver — read the names
// of the current request from the request context.
//
// The capture wraps the accepted connection and passively tees the
// plaintext bytes the HTTP server reads, so it sees HTTP/1.1 request
// heads as they arrive. Connections that turn out to speak HTTP/2 (the
// h2c preface or an ALPN-negotiated h2 stream) stop capturing: the
// HPACK-encoded header order is not recovered and consumers fall back
// to the header map's unspecified order. A request head that carries a
// chunked body also ends the connection's capture, because the body's
// end cannot be located without full chunked framing.
package headerorder

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// maxPendingHeads bounds the queue of completed request heads kept per
// connection. A pipelining client may deliver heads faster than
// handlers consume them; the oldest entries are dropped beyond this.
const maxPendingHeads = 32

// headerPair is one header line of a captured request head: the name
// as it appeared on the wire and the trimmed value.
type headerPair struct {
	name  string
	value string
}

// ConnState holds the header name and value pairs of the completed
// HTTP/1.1 request heads on one connection, in arrival order. The HTTP
// server's ConnContext hook injects it into the request context.
// Consumers pair each request with its own head through the multiset
// of header name and value pairs, so pipelined requests — even two
// with the same header names but different values — cannot observe
// another request's order.
type ConnState struct {
	mu      sync.Mutex
	pending [][]headerPair
}

// record appends one completed head's header pairs, dropping the
// oldest entry when the queue is full.
func (s *ConnState) record(pairs []headerPair) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pending) >= maxPendingHeads {
		copy(s.pending, s.pending[1:])
		s.pending[len(s.pending)-1] = pairs
		return
	}
	s.pending = append(s.pending, pairs)
}

// take removes and returns the pending head whose header pairs match the
// one described by match, which reports whether a candidate head belongs
// to the calling request. When zero or several queued heads match, no
// head is removed: several matching heads cannot be attributed to one
// request — an earlier head whose consumer never took it is
// indistinguishable from the caller's own — so the caller must fall back.
func (s *ConnState) take(match func(pairs []headerPair) bool) ([]headerPair, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	index := -1
	for i, pairs := range s.pending {
		if !match(pairs) {
			continue
		}
		if index >= 0 {
			return nil, false
		}
		index = i
	}
	if index < 0 {
		return nil, false
	}
	pairs := s.pending[index]
	s.pending = append(s.pending[:index], s.pending[index+1:]...)
	out := make([]headerPair, len(pairs))
	copy(out, pairs)
	return out, true
}

type ctxKey struct{}

// NewContext returns a context carrying the connection's capture
// state. The HTTP server wires this through its ConnContext hook.
func NewContext(ctx context.Context, state *ConnState) context.Context {
	return context.WithValue(ctx, ctxKey{}, state)
}

// FromContext returns the wire-ordered header names of the request the
// context belongs to. The headers argument is the request's header
// map with the Host header present under the canonical "Host" key —
// net/http promotes Host out of the map, so callers inject it with the
// request's host value (the WAF planes do). The request is paired with
// its own captured head by comparing the multiset of (name, value)
// pairs: a queued head with a different multiset never pairs with this
// request, while an unconsumed earlier head whose multiset matches — a
// pipelined duplicate that its own consumer never took — makes the
// pairing ambiguous, and no order is reported. The boolean reports
// whether an order was captured; without it the caller must fall back
// to the header map.
func FromContext(ctx context.Context, headers map[string][]string) ([]string, bool) {
	state, ok := ctx.Value(ctxKey{}).(*ConnState)
	if !ok || state == nil {
		return nil, false
	}
	expected, hasHost := headerPairCounts(headers)
	if !hasHost {
		// The wire head always carries Host, so a map without it cannot
		// be paired with any captured head.
		return nil, false
	}
	pairs, ok := state.take(func(candidate []headerPair) bool {
		return headerPairsEqual(candidate, expected)
	})
	if !ok {
		return nil, false
	}
	names := make([]string, len(pairs))
	for i, pair := range pairs {
		names[i] = pair.name
	}
	return names, true
}

// headerPairCounts canonicalises a header map into (name, value) pair
// counts. The boolean reports whether the map carries a Host header.
func headerPairCounts(headers map[string][]string) (map[headerPair]int, bool) {
	counts := make(map[headerPair]int)
	hasHost := false
	for name, values := range headers {
		canonical := http.CanonicalHeaderKey(name)
		if canonical == "Host" {
			hasHost = true
		}
		// Every value of a repeated header is one wire line, so each
		// contributes its own pair.
		for _, value := range values {
			counts[headerPair{canonical, value}]++
		}
	}
	return counts, hasHost
}

// headerPairsEqual reports whether a captured head's header pairs
// canonicalise to exactly the expected pair counts.
func headerPairsEqual(pairs []headerPair, expected map[headerPair]int) bool {
	seen := make(map[headerPair]int, len(pairs))
	for _, pair := range pairs {
		seen[headerPair{http.CanonicalHeaderKey(pair.name), pair.value}]++
	}
	if len(seen) != len(expected) {
		return false
	}
	for pair, count := range seen {
		if expected[pair] != count {
			return false
		}
	}
	return true
}

// ConnContext returns an http.Server.ConnContext hook that injects the
// connection's capture state into the request context. The optional
// base hook runs first so existing behaviour chains. For captured TLS
// connections the hook also completes the handshake eagerly and stores
// the connection state: the capture wrapper hides the *tls.Conn from
// the HTTP server's concrete-type detection, so the server neither
// populates Request.TLS nor applies its own handshake handling —
// TLSStateMiddleware restores Request.TLS from the stored state.
func ConnContext(base func(ctx context.Context, c net.Conn) context.Context) func(ctx context.Context, c net.Conn) context.Context {
	return func(ctx context.Context, c net.Conn) context.Context {
		if base != nil {
			ctx = base(ctx, c)
		}
		if cc, ok := c.(*captureConn); ok {
			ctx = NewContext(ctx, cc.State())
			if state := cc.completeTLS(ctx); state != nil {
				ctx = withTLSState(ctx, state)
			}
		}
		return ctx
	}
}

type tlsStateKey struct{}

func withTLSState(ctx context.Context, state *tls.ConnectionState) context.Context {
	return context.WithValue(ctx, tlsStateKey{}, state)
}

// TLSStateMiddleware restores Request.TLS from the connection state a
// ConnContext hook captured, keeping TLS-visible behaviour (viewer
// protocol policies, client-certificate authentication) intact on
// capture-wrapped TLS listeners. It is a no-op without stored state.
func TLSStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if state, ok := r.Context().Value(tlsStateKey{}).(*tls.ConnectionState); ok && state != nil {
			r.TLS = state
		}
		next.ServeHTTP(w, r)
	})
}

// tlsHandshakeTimeout bounds the eagerly completed TLS handshake so a
// stalling client cannot hold the connection goroutine indefinitely.
const tlsHandshakeTimeout = 10 * time.Second

// completeTLS drives the TLS handshake of a captured TLS connection and
// returns its connection state, or nil for plaintext connections. A
// handshake failure returns nil; the server's read path surfaces it.
func (c *captureConn) completeTLS(ctx context.Context) *tls.ConnectionState {
	tlsConn, ok := c.Conn.(*tls.Conn)
	if !ok {
		return nil
	}
	handshakeCtx, cancel := context.WithTimeout(ctx, tlsHandshakeTimeout)
	defer cancel()
	if err := tlsConn.HandshakeContext(handshakeCtx); err != nil {
		return nil
	}
	state := tlsConn.ConnectionState()
	return &state
}

// NewListener wraps a listener so every accepted connection captures
// the wire order of its HTTP/1.1 request header names. For TLS the
// caller wraps the decrypted layer: pass the tls.NewListener result so
// the tee observes plaintext.
func NewListener(l net.Listener) net.Listener {
	return &captureListener{Listener: l}
}

type captureListener struct {
	net.Listener
}

func (l *captureListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &captureConn{Conn: conn, parser: newHeadParser()}, nil
}

type captureConn struct {
	net.Conn
	parser *headParser
	once   sync.Once
	state  *ConnState
}

func (c *captureConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if n > 0 {
		c.parser.feed(b[:n], c.connState())
	}
	return n, err
}

// connState returns the connection's shared capture state, created on
// first use.
func (c *captureConn) connState() *ConnState {
	c.once.Do(func() { c.state = &ConnState{} })
	return c.state
}

// State exposes the capture state for the HTTP server's ConnContext
// hook.
func (c *captureConn) State() *ConnState {
	return c.connState()
}

// headParser incrementally extracts the header names of HTTP/1.1
// request heads from the teed byte stream. It walks head, body and
// head again: the body length comes from Content-Length, so the parser
// resynchronises on the next request head once the body is consumed.
type headParser struct {
	state      *ConnState
	head       []byte
	bodyRemain int64
	inBody     bool
	done       bool
	pairs      []headerPair
}

// http2MethodPrefix starts the HTTP/2 connection preface request line;
// no HTTP/1.1 request line begins with it.
const http2MethodPrefix = "PRI "

// maxHeadBytes bounds the buffered request head; a larger head ends
// capture for the connection instead of growing without limit.
const maxHeadBytes = 64 << 10

func newHeadParser() *headParser {
	return &headParser{head: make([]byte, 0, 1024)}
}

func (p *headParser) feed(b []byte, state *ConnState) {
	if p.done {
		return
	}
	p.state = state
	for len(b) > 0 && !p.done {
		if p.inBody {
			consume := int64(len(b))
			if consume > p.bodyRemain {
				consume = p.bodyRemain
			}
			p.bodyRemain -= consume
			b = b[consume:]
			if p.bodyRemain == 0 {
				p.inBody = false
			}
			continue
		}
		p.head = append(p.head, b...)
		b = nil
		if bytes.HasPrefix(p.head, []byte(http2MethodPrefix)) {
			p.done = true
			return
		}
		idx := bytes.Index(p.head, []byte("\r\n\r\n"))
		if idx < 0 {
			if len(p.head) > maxHeadBytes {
				p.done = true
			}
			continue
		}
		p.consumeHead(p.head[:idx+4])
		remaining := make([]byte, len(p.head)-idx-4)
		copy(remaining, p.head[idx+4:])
		p.head = p.head[:0]
		b = remaining
	}
}

// consumeHead parses one complete request head: the request line is
// skipped, every header line contributes its name and trimmed value,
// and an obs-fold continuation line (leading space or tab) belongs to
// the previous header. With the head parsed, the parser counts the
// declared body before looking for the next request head; a chunked
// body ends the connection's capture.
func (p *headParser) consumeHead(head []byte) {
	p.pairs = p.pairs[:0]
	contentLength := int64(0)
	chunked := false
	lines := bufio.NewReader(bytes.NewReader(head))
	first := true
	for {
		line, err := lines.ReadBytes('\n')
		if err != nil {
			break
		}
		line = bytes.TrimSuffix(line, []byte("\n"))
		line = bytes.TrimSuffix(line, []byte("\r"))
		if first {
			first = false
			continue
		}
		if len(line) == 0 {
			break
		}
		if line[0] == ' ' || line[0] == '\t' {
			continue
		}
		name, value, found := bytes.Cut(line, []byte(":"))
		if !found {
			continue
		}
		nameText := strings.TrimSpace(string(name))
		if nameText == "" {
			continue
		}
		valueText := strings.TrimSpace(string(value))
		p.pairs = append(p.pairs, headerPair{name: nameText, value: valueText})
		switch strings.ToLower(nameText) {
		case "content-length":
			if v, err := strconv.ParseInt(valueText, 10, 64); err == nil {
				contentLength = v
			}
		case "transfer-encoding":
			if strings.Contains(strings.ToLower(valueText), "chunked") {
				chunked = true
			}
		}
	}
	pairs := make([]headerPair, len(p.pairs))
	copy(pairs, p.pairs)
	p.state.record(pairs)
	if chunked || contentLength < 0 {
		p.done = true
		return
	}
	if contentLength > 0 {
		p.inBody = true
		p.bodyRemain = contentLength
	}
}
