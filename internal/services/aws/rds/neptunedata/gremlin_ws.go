package neptunedata

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
	"vorpalstacks/pkg/gremlinparser"
)

const (
	wsKeepAliveInterval  = 60 * time.Second
	wsWriteTimeout       = 10 * time.Second
	wsSendBufferSize     = 256
	wsMaxMessageSize     = 1 * 1024 * 1024
	wsSessionIdleTimeout = 5 * time.Minute
	wsSessionCleanupTick = 30 * time.Second
)

// gremlinUpgrader is the WebSocket upgrader for TinkerPop Gremlin Server
// connections. Accepts all origins in test mode; production deployments
// should restrict CheckOrigin.
var gremlinUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"graphson-v3"},
}

// GremlinWSServer handles WebSocket connections for the TinkerPop Gremlin
// Server protocol. Each connection is managed by a pair of goroutines
// (readPump/writePump) following the same pattern as AppSync's EventServer.
type GremlinWSServer struct {
	service *NeptuneDataService
}

// NewGremlinWSServer creates a new WebSocket server bound to the given
// NeptuneDataService for graph engine access.
func NewGremlinWSServer(service *NeptuneDataService) *GremlinWSServer {
	return &GremlinWSServer{service: service}
}

// gremlinWSConn represents a single WebSocket connection to a Gremlin client.
type gremlinWSConn struct {
	id        string
	conn      *websocket.Conn
	sendCh    chan []byte
	service   *NeptuneDataService
	clusterID string
	db        *graphengine.DB
	session   *gremlinSession
	mu        sync.RWMutex
	closed    bool
}

// ServeHTTP upgrades the HTTP connection to WebSocket and starts the
// read/write pumps. The clusterID parameter identifies which graph engine
// to use for query execution.
func (s *GremlinWSServer) ServeHTTP(w http.ResponseWriter, r *http.Request, clusterID string) {
	conn, err := gremlinUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("Gremlin WebSocket upgrade failed", logs.Err(err))
		return
	}

	db := s.service.GetClusterEngine(clusterID)
	if db == nil {
		msg, _ := json.Marshal(newGremlinErrorResponse("", statusServerError, "cluster engine not available"))
		_ = conn.WriteMessage(websocket.TextMessage, msg)
		_ = conn.Close()
		return
	}

	ws := &gremlinWSConn{
		id:        uuid.New().String(),
		conn:      conn,
		sendCh:    make(chan []byte, wsSendBufferSize),
		service:   s.service,
		clusterID: clusterID,
		db:        db,
	}

	logs.Info("Gremlin WebSocket connected",
		logs.String("connId", ws.id),
		logs.String("cluster", clusterID))

	go s.writePump(ws)
	go s.readPump(ws)
}

// readPump reads messages from the WebSocket connection and dispatches them
// to handleMessage. Exits when the connection is closed or a read error occurs.
func (s *GremlinWSServer) readPump(ws *gremlinWSConn) {
	defer func() {
		ws.mu.Lock()
		if !ws.closed {
			ws.closed = true
			close(ws.sendCh)
		}
		ws.mu.Unlock()
		ws.conn.Close()
		logs.Info("Gremlin WebSocket disconnected", logs.String("connId", ws.id))
	}()

	ws.conn.SetReadLimit(wsMaxMessageSize)
	_ = ws.conn.SetReadDeadline(time.Now().Add(wsKeepAliveInterval * 2))

	ws.conn.SetPongHandler(func(string) error {
		_ = ws.conn.SetReadDeadline(time.Now().Add(wsKeepAliveInterval * 2))
		return nil
	})

	for {
		msgType, raw, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logs.Warn("Gremlin WebSocket read error",
					logs.String("connId", ws.id),
					logs.Err(err))
			}
			return
		}

		if msgType == websocket.BinaryMessage || msgType == websocket.TextMessage {
			s.handleMessage(ws, raw)
		}
	}
}

// writePump pumps messages from the send channel to the WebSocket connection.
// Sends periodic ping frames for keep-alive. Exits when sendCh is closed.
func (s *GremlinWSServer) writePump(ws *gremlinWSConn) {
	ticker := time.NewTicker(wsKeepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-ws.sendCh:
			_ = ws.conn.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
			if !ok {
				_ = ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := ws.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = ws.conn.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage dispatches an incoming TinkerPop message based on its op field.
func (s *GremlinWSServer) handleMessage(ws *gremlinWSConn, raw []byte) {
	req, mimeType, err := decodeGremlinRequest(raw)
	if err != nil {
		s.sendResponse(ws, newGremlinErrorResponse("", statusMalformedRequest, err.Error()))
		return
	}

	switch req.Op {
	case "authentication":
		s.handleAuthentication(ws, req)
	case "eval":
		s.handleEval(ws, req, mimeType)
	case "close":
		s.handleClose(ws, req)
	default:
		s.sendResponse(ws, newGremlinErrorResponse(req.RequestID, statusMalformedRequest,
			"unknown op: "+req.Op))
	}
}

// handleAuthentication processes an authentication request. In TEST_MODE or
// without an external auth provider, authentication is accepted immediately.
func (s *GremlinWSServer) handleAuthentication(ws *gremlinWSConn, req *GremlinRequest) {
	// Accept all authentication without challenge for self-hosted mode.
	// Production deployments would validate SASL tokens here.
	s.sendResponse(ws, newGremlinSuccessResponse(req.RequestID, nil))
}

// handleEval executes a Gremlin traversal from a TinkerPop eval request.
// Parses the query string through gremlinparser, executes it against the
// cluster's graph engine, and returns the result in GraphSON v3 format.
func (s *GremlinWSServer) handleEval(ws *gremlinWSConn, req *GremlinRequest, mimeType string) {
	if req.Args.Gremlin == "" {
		s.sendResponse(ws, newGremlinErrorResponse(req.RequestID, statusMalformedRequest,
			"gremlin query string is required"))
		return
	}

	// Merge session bindings if this connection has an active session
	params := mergeParams(ws, req)

	reader := graphengine.GraphReader(ws.db)
	writer := graphengine.GraphWriter(ws.db)

	parsed, err := gremlinparser.Parse(req.Args.Gremlin)
	if err != nil {
		s.sendResponse(ws, newGremlinErrorResponse(req.RequestID, statusScriptError, err.Error()))
		return
	}

	ctx, cancel := evalContext(ws, req)
	defer cancel()

	result, execErr := gremlinparser.ExecuteQuery(ctx, reader, writer, parsed, params)
	if execErr != nil {
		s.sendResponse(ws, newGremlinErrorResponse(req.RequestID, statusScriptError, execErr.Error()))
		return
	}

	encoded := EncodeGraphSON3Result(result)
	s.sendResponse(ws, newGremlinSuccessResponse(req.RequestID, encoded))

	// Update session state
	updateSession(ws, req, result)
}

// handleClose terminates a session or closes the connection. When a session
// processor sends "close", the session is cleaned up.
func (s *GremlinWSServer) handleClose(ws *gremlinWSConn, req *GremlinRequest) {
	if ws.session != nil {
		ws.mu.Lock()
		ws.session = nil
		ws.mu.Unlock()
	}
	s.sendResponse(ws, newGremlinSuccessResponse(req.RequestID, nil))
}

// mergeParams builds the parameter map for query execution by combining
// session bindings with request-level bindings.
func mergeParams(ws *gremlinWSConn, req *GremlinRequest) map[string]any {
	params := make(map[string]any)

	ws.mu.RLock()
	if ws.session != nil {
		for k, v := range ws.session.bindings {
			params[k] = v
		}
	}
	ws.mu.RUnlock()

	for k, v := range req.Args.Bindings {
		params[k] = DecodeGraphSON3Args(v)
	}
	return params
}

// evalContext creates a context for query execution with an optional timeout.
func evalContext(ws *gremlinWSConn, req *GremlinRequest) (context.Context, context.CancelFunc) {
	if req.Args.EvalTimeout > 0 {
		return context.WithTimeout(context.Background(), time.Duration(req.Args.EvalTimeout)*time.Millisecond)
	}
	return context.WithCancel(context.Background())
}

// updateSession updates the session state after a successful eval. Creates a
// new session if the processor is "session" and none exists, and merges
// aliases and bindings.
func updateSession(ws *gremlinWSConn, req *GremlinRequest, result interface{}) {
	if req.Processor != "session" {
		return
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.session == nil {
		ws.session = &gremlinSession{
			id:       req.Args.Session,
			aliases:  make(map[string]string),
			bindings: make(map[string]interface{}),
		}
	}

	for k, v := range req.Args.Aliases {
		ws.session.aliases[k] = v
	}
	for k, v := range req.Args.Bindings {
		ws.session.bindings[k] = v
	}
}

// sendResponse serialises and queues a GremlinResponse for sending on the
// WebSocket connection.
func (s *GremlinWSServer) sendResponse(ws *gremlinWSConn, resp GremlinResponse) {
	msg, err := encodeResponseToJSON(resp)
	if err != nil {
		logs.Error("failed to encode Gremlin WS response", logs.Err(err))
		return
	}
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ws.closed {
		return
	}
	select {
	case ws.sendCh <- msg:
	default:
		logs.Warn("Gremlin WS send buffer full, dropping response",
			logs.String("connId", ws.id))
	}
}
