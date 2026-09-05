package lambda

// This file emulates the AWS Lambda Runtime API (version 2018-06-01) for
// the runtimes that enter through /var/runtime/bootstrap — the provided.*
// custom runtimes, the RIC-based managed images, and the image-package
// functions' own ENTRYPOINT. A server instance carries any number of
// invocation rounds: the host pushes a round with BeginRound, the
// runtime's GET /runtime/invocation/next delivers it together with the
// Lambda-Runtime-* headers derived from the invocation record, and the
// runtime's response or error POST becomes that round's answer. The
// single-shot bootstrap path uses exactly the constructor's initial round;
// the sandbox pool reuses one server across the consecutive rounds of a
// persistent container. Serving the API on the host side keeps every
// invoke path on the single invokeFunction choke point, instead of routing
// events into the container's own runtime interface emulator.

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	runtimeAPIVersion = "2018-06-01"

	runtimeAPINextPath    = "/" + runtimeAPIVersion + "/runtime/invocation/next"
	runtimeAPIInitErrPath = "/" + runtimeAPIVersion + "/runtime/init/error"
	runtimeAPIInvokePath  = "/" + runtimeAPIVersion + "/runtime/invocation/"

	// The captured-answer kinds: a posted response, an invocation error,
	// or an initialization error. The AWS contract marks invocations that
	// answered through an error endpoint Unhandled.
	runtimeAPIResponseKind = "response"
	runtimeAPIErrorKind    = "error"
	runtimeAPIInitErrKind  = "init-error"
)

// runtimeAPIServer serves the Runtime API for one runtime process across
// any number of invocation rounds. When no round is pending, /next parks
// until one arrives or the server closes — the documented iterator-style
// blocking call that a sandbox's idle runtime sits in between invokes.
type runtimeAPIServer struct {
	srv      *http.Server
	listener net.Listener

	mu    sync.Mutex
	rec   invocationRecord // the pending round, or the round served last
	event []byte

	// pending is set while a pushed round has not been delivered to a /next
	// long-poll; delivering it consumes the flag under the mutex, so only
	// one long-poll ever serves a round.
	pending bool

	body     []byte
	kind     string
	captured bool

	// parked holds the /next long-poll currently waiting for an event. Its
	// presence is the sandbox pool's idle-eligibility signal: the runtime
	// has finished its previous round and asked for the next event.
	parked *runtimeAPIPark

	// answered closes when the current round's answer lands. The AWS
	// contract settles a synchronous invocation at the response POST, not
	// when the runtime process ends, so the invoke path waits on this
	// channel instead of the process. BeginRound re-arms it for the next
	// round; rounds on one server are strictly sequential from the host
	// side, one owning invoke at a time.
	answered chan struct{}

	closed chan struct{}
}

// runtimeAPIPark is one parked /next long-poll; BeginRound closes wake to
// hand the pending round to the waiting runtime.
type runtimeAPIPark struct {
	wake chan struct{}
}

// startRuntimeAPI binds a single-round Runtime API listener and pushes the
// initial round. The listener binds all interfaces: the function container
// reaches the host through the Docker bridge gateway
// (host.docker.internal → host-gateway), and a loopback bind is not
// reachable through it — the same reachability the main API server relies
// on for the endpoint URL injected into containers.
func startRuntimeAPI(rec invocationRecord, event []byte) (*runtimeAPIServer, error) {
	api, err := newRuntimeAPIServer()
	if err != nil {
		return nil, err
	}
	api.BeginRound(rec, event)
	return api, nil
}

// newRuntimeAPIServer binds a round-less Runtime API listener for a
// persistent sandbox: the container's runtime parks its first /next as it
// boots, and each owning invoke pushes its round with BeginRound.
func newRuntimeAPIServer() (*runtimeAPIServer, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("listen runtime api: %w", err)
	}
	api := &runtimeAPIServer{
		listener: ln,
		answered: make(chan struct{}),
		closed:   make(chan struct{}),
	}
	api.srv = &http.Server{Handler: api}
	go func() { _ = api.srv.Serve(ln) }()
	return api, nil
}

// BeginRound pushes the next invocation round: the event and the
// Lambda-Runtime-* headers derived from the record. A parked /next wakes
// and serves the round at once; otherwise the round waits for the runtime's
// next /next. The previous round's answer is consumed by the time the host
// pushes the next one.
func (a *runtimeAPIServer) BeginRound(rec invocationRecord, event []byte) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rec = rec
	a.event = event
	a.pending = true
	a.body = nil
	a.kind = ""
	a.captured = false
	a.answered = make(chan struct{})
	if a.parked != nil {
		close(a.parked.wake)
		a.parked = nil
	}
}

// Idle reports whether a /next long-poll is parked waiting for an event:
// the runtime has finished its previous round and is ready for the next.
func (a *runtimeAPIServer) Idle() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.parked != nil
}

// Port returns the port the listener bound.
func (a *runtimeAPIServer) Port() int {
	return a.listener.Addr().(*net.TCPAddr).Port
}

// Addr returns the host:port value for the AWS_LAMBDA_RUNTIME_API exec
// environment variable.
func (a *runtimeAPIServer) Addr() string {
	return net.JoinHostPort("host.docker.internal", strconv.Itoa(a.Port()))
}

// Close stops the listener, aborts any pending /next long-poll, and
// releases the port. The invoke path always calls it once the execution
// finishes, captured answer or not.
func (a *runtimeAPIServer) Close() {
	select {
	case <-a.closed:
		return
	default:
	}
	close(a.closed)
	_ = a.srv.Close()
}

// Captured returns the bootstrap's answer: the posted body, which endpoint
// it went to, and whether anything was posted at all.
func (a *runtimeAPIServer) Captured() (body []byte, kind string, ok bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.body, a.kind, a.captured
}

// Answered returns a channel that closes when the current round's answer
// (a response or an error POST) lands. The channel is per round: an answer
// that landed before the waiter picked the channel up is already closed.
func (a *runtimeAPIServer) Answered() <-chan struct{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.answered
}

func (a *runtimeAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == runtimeAPINextPath:
		a.serveNext(w, r)
	case r.Method == http.MethodPost && r.URL.Path == runtimeAPIInitErrPath:
		a.capture(w, r, runtimeAPIInitErrKind)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, runtimeAPIInvokePath) && strings.HasSuffix(r.URL.Path, "/response"):
		a.capture(w, r, runtimeAPIResponseKind)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, runtimeAPIInvokePath) && strings.HasSuffix(r.URL.Path, "/error"):
		a.capture(w, r, runtimeAPIErrorKind)
	default:
		http.NotFound(w, r)
	}
}

// serveNext answers the runtime's long-poll with the pending round, or
// parks the long-poll until one arrives. A well-behaved runtime loops back
// to /next after answering; between rounds the parked long-poll is exactly
// the idle signal the sandbox pool routes on. When the server closes
// without an event — the single-shot path's completed execution, or a
// sandbox being destroyed — the parked long-poll aborts with 503 and the
// timeout reaps the idle runtime process from there.
func (a *runtimeAPIServer) serveNext(w http.ResponseWriter, r *http.Request) {
	for {
		a.mu.Lock()
		if a.pending {
			a.pending = false
			rec, event := a.rec, a.event
			a.mu.Unlock()
			writeInvocationHeaders(w, rec)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(event)
			return
		}
		// No round to serve: park this long-poll. A runtime that opens a
		// further /next while one is parked misbehaves; the older park is
		// woken to re-park behind the newer one, and no round is lost
		// because delivery is decided under the mutex.
		if a.parked != nil {
			close(a.parked.wake)
		}
		park := &runtimeAPIPark{wake: make(chan struct{})}
		a.parked = park
		a.mu.Unlock()
		select {
		case <-park.wake:
			// A round arrived or a newer long-poll took the slot; loop to
			// decide under the lock.
		case <-a.closed:
			http.Error(w, "runtime api closed", http.StatusServiceUnavailable)
			return
		case <-r.Context().Done():
			// The long-poll went away with its runtime process; only this
			// request may clear the park slot.
			a.mu.Lock()
			if a.parked == park {
				a.parked = nil
			}
			a.mu.Unlock()
			return
		}
	}
}

// writeInvocationHeaders emits the Lambda-Runtime-* headers of one round.
// The AWS contract also lists Lambda-Runtime-Trace-Id and
// Lambda-Runtime-Cognito-Identity; this platform has neither an X-Ray nor a
// mobile-SDK identity channel, so those headers stay absent and bootstraps
// must tolerate their absence (identity stays null on the wrapper runtimes
// for the same reason).
func writeInvocationHeaders(w http.ResponseWriter, rec invocationRecord) {
	w.Header().Set("Lambda-Runtime-Aws-Request-Id", rec.RequestID)
	w.Header().Set("Lambda-Runtime-Deadline-Ms", strconv.FormatInt(rec.DeadlineUnixMS(), 10))
	w.Header().Set("Lambda-Runtime-Invoked-Function-Arn", rec.InvokedARN)
	// The invocation-id header is optional in the contract; this platform
	// runs one attempt per execution, so the invocation id is the request
	// id the record minted.
	w.Header().Set("Lambda-Runtime-Invocation-Id", rec.RequestID)
	if rec.ClientContextJSON != "" {
		// The ClientContext document travels base64-encoded, the same
		// encoding the Invoke parameter carried it in with.
		w.Header().Set("Lambda-Runtime-Client-Context",
			base64.StdEncoding.EncodeToString([]byte(rec.ClientContextJSON)))
	}
}

// capture records the runtime's answer body for the current round. The
// first POST wins; the id in the path is not validated because the runtime
// only ever knows the request id this server served it. A body over the
// synchronous payload limit is rejected with 413 and settles the
// invocation as an error: truncating it silently would deliver a corrupt
// payload as a successful answer. The AWS Runtime API documentation does
// not fix the error wording for an oversized response, so the document is
// this platform's own, errorMessage-only per the established convention.
func (a *runtimeAPIServer) capture(w http.ResponseWriter, r *http.Request, kind string) {
	body, err := io.ReadAll(io.LimitReader(r.Body, syncMaxPayloadSize+1))
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}
	if len(body) > syncMaxPayloadSize {
		tooLarge := fmt.Sprintf(`{"errorMessage":"Response payload exceeded the %d MB synchronous invocation limit"}`,
			syncMaxPayloadSize/(1024*1024))
		a.mu.Lock()
		if !a.captured {
			a.body = []byte(tooLarge)
			a.kind = runtimeAPIErrorKind
			a.captured = true
			close(a.answered)
		}
		a.mu.Unlock()
		http.Error(w, "response payload too large", http.StatusRequestEntityTooLarge)
		return
	}
	a.mu.Lock()
	if !a.captured {
		a.body = body
		a.kind = kind
		a.captured = true
		close(a.answered)
	}
	a.mu.Unlock()
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte("OK"))
}
