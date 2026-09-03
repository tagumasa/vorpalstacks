package lambda

// This file emulates the AWS Lambda Runtime API (version 2018-06-01) for
// the runtimes that enter through /var/runtime/bootstrap — the provided.*
// custom runtimes and the RIC-based managed images. One server instance
// serves exactly one invocation: GET /runtime/invocation/next hands the
// event and the Lambda-Runtime-* headers derived from the invocation
// record to the bootstrap, and the bootstrap's response or error POST
// becomes the invocation payload. Serving the API on the host side keeps
// every invoke path on the single invokeFunction choke point, instead of
// routing events into the container's own runtime interface emulator.

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

// runtimeAPIServer serves the Runtime API for a single invocation.
type runtimeAPIServer struct {
	srv      *http.Server
	listener net.Listener
	rec      invocationRecord
	event    []byte

	mu       sync.Mutex
	served   bool
	body     []byte
	kind     string
	captured bool

	// answered closes once when the first answer lands. The AWS contract
	// settles a synchronous invocation at the response POST, not when the
	// bootstrap process ends, so the invoke path waits on this channel
	// instead of the exec.
	answered chan struct{}

	closed chan struct{}
}

// startRuntimeAPI binds the per-invocation Runtime API listener. The
// listener binds all interfaces: the function container reaches the host
// through the Docker bridge gateway (host.docker.internal → host-gateway),
// and a loopback bind is not reachable through it — the same reachability
// the main API server relies on for the endpoint URL injected into
// containers.
func startRuntimeAPI(rec invocationRecord, event []byte) (*runtimeAPIServer, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("listen runtime api: %w", err)
	}
	api := &runtimeAPIServer{
		listener: ln,
		rec:      rec,
		event:    event,
		answered: make(chan struct{}),
		closed:   make(chan struct{}),
	}
	api.srv = &http.Server{Handler: api}
	go func() { _ = api.srv.Serve(ln) }()
	return api, nil
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

// Answered returns a channel that closes when the bootstrap's first answer
// (a response or an error POST) lands.
func (a *runtimeAPIServer) Answered() <-chan struct{} {
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

// serveNext answers the bootstrap's long-poll with the event and the
// invocation headers. The AWS contract also lists Lambda-Runtime-Trace-Id
// and Lambda-Runtime-Cognito-Identity; this platform has neither an X-Ray
// nor a mobile-SDK identity channel, so those headers stay absent and
// bootstraps must tolerate their absence (identity stays null on the
// wrapper runtimes for the same reason).
func (a *runtimeAPIServer) serveNext(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	if a.served {
		a.mu.Unlock()
		// A well-behaved bootstrap loops back to /next after answering.
		// This server serves exactly one invocation, so the long-poll
		// waits until the completed execution closes the server — the
		// timeout reaps the idle bootstrap process from there.
		select {
		case <-a.closed:
		case <-r.Context().Done():
		}
		http.Error(w, "runtime api closed", http.StatusServiceUnavailable)
		return
	}
	a.served = true
	a.mu.Unlock()

	w.Header().Set("Lambda-Runtime-Aws-Request-Id", a.rec.RequestID)
	w.Header().Set("Lambda-Runtime-Deadline-Ms", strconv.FormatInt(a.rec.DeadlineUnixMS(), 10))
	w.Header().Set("Lambda-Runtime-Invoked-Function-Arn", a.rec.InvokedARN)
	// The invocation-id header is optional in the contract; this platform
	// runs one attempt per execution, so the invocation id is the request
	// id the record minted.
	w.Header().Set("Lambda-Runtime-Invocation-Id", a.rec.RequestID)
	if a.rec.ClientContextJSON != "" {
		// The ClientContext document travels base64-encoded, the same
		// encoding the Invoke parameter carried it in with.
		w.Header().Set("Lambda-Runtime-Client-Context",
			base64.StdEncoding.EncodeToString([]byte(a.rec.ClientContextJSON)))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(a.event)
}

// capture records the bootstrap's answer body. The first POST wins; the
// id in the path is not validated because the bootstrap only ever knows
// the request id this server served it. A body over the synchronous
// payload limit is rejected with 413 and settles the invocation as an
// error: truncating it silently would deliver a corrupt payload as a
// successful answer. The AWS Runtime API documentation does not fix the
// error wording for an oversized response, so the document is this
// platform's own, errorMessage-only per the established convention.
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
