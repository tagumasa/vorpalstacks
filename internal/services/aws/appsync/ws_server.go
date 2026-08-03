package appsync

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/fqdnrouter"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

const (
	wsSubprotocol           = "aws-appsync-event-ws"
	keepAliveInterval       = 60 * time.Second
	connectionTimeoutMs     = 300000
	maxEventsPerPublish     = 5
	maxEventSizeBytes       = 240 * 1024
	maxSubscriptionsPerConn = 100
)

var (
	subscriptionIDPattern = regexp.MustCompile(`^[a-zA-Z0-9\-_+]{1,128}$`)
	channelPathPattern    = regexp.MustCompile(`^\/?[A-Za-z0-9](?:[A-Za-z0-9\-]{0,48}[A-Za-z0-9])?(?:\/[A-Za-z0-9](?:[A-Za-z0-9\-]{0,48}[A-Za-z0-9])?){0,4}(?:\/\*)?\/?$`)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		host := r.Host
		if host == "" {
			host = r.URL.Host
		}
		return origin == "http://"+host || origin == "https://"+host
	},
	Subprotocols: []string{wsSubprotocol},
}

// wsConnection represents a single WebSocket connection to the Event API.
type wsConnection struct {
	id            string
	apiId         string
	conn          *websocket.Conn
	sendCh        chan []byte
	subscriptions map[string]*subscription
	mu            sync.RWMutex
	closed        bool
	iamVerified   bool
}

// subscription represents a single channel subscription on a connection.
type subscription struct {
	id      string
	channel string
	auth    map[string]interface{}
}

// channelManager tracks all active subscriptions across connections,
// enabling efficient fan-out when events are published.
type channelManager struct {
	channels map[string]map[string]bool
	connSubs map[string]map[string]string
	mu       sync.RWMutex
}

func newChannelManager() *channelManager {
	return &channelManager{
		channels: make(map[string]map[string]bool),
		connSubs: make(map[string]map[string]string),
	}
}

func (cm *channelManager) subscribe(channel string, connId, subId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	compositeKey := connId + ":" + subId

	if cm.channels[channel] == nil {
		cm.channels[channel] = make(map[string]bool)
	}
	cm.channels[channel][compositeKey] = true

	if cm.connSubs[connId] == nil {
		cm.connSubs[connId] = make(map[string]string)
	}
	cm.connSubs[connId][subId] = channel
}

func (cm *channelManager) unsubscribe(channel string, connId string, subId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	compositeKey := connId + ":" + subId

	if subs, ok := cm.channels[channel]; ok {
		delete(subs, compositeKey)
		if len(subs) == 0 {
			delete(cm.channels, channel)
		}
	}

	if subMap, ok := cm.connSubs[connId]; ok {
		delete(subMap, subId)
		if len(subMap) == 0 {
			delete(cm.connSubs, connId)
		}
	}
}

// removeConnection removes all subscriptions for a given connection ID.
func (cm *channelManager) removeConnection(connId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	subMap, ok := cm.connSubs[connId]
	if !ok {
		return
	}
	for subId, channel := range subMap {
		compositeKey := connId + ":" + subId
		if subs, ok := cm.channels[channel]; ok {
			delete(subs, compositeKey)
			if len(subs) == 0 {
				delete(cm.channels, channel)
			}
		}
	}
	delete(cm.connSubs, connId)
}

// matchSubscriptions returns all (connId, subId) pairs whose subscribed channel
// matches the published channel path. Supports wildcard matching:
// /namespace/* matches any channel under /namespace/.
func (cm *channelManager) matchSubscriptions(publishedChannel string) []struct {
	connId string
	subId  string
} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var matches []struct {
		connId string
		subId  string
	}

	for ch, subs := range cm.channels {
		if channelMatches(ch, publishedChannel) {
			for compositeKey := range subs {
				parts := strings.SplitN(compositeKey, ":", 2)
				if len(parts) != 2 {
					continue
				}
				matches = append(matches, struct {
					connId string
					subId  string
				}{connId: parts[0], subId: parts[1]})
			}
		}
	}

	return matches
}

// channelMatches checks whether a subscriber's channel pattern matches
// a published channel path. Supports trailing wildcard segments.
func channelMatches(subscriberChannel, publishedChannel string) bool {
	subscriberChannel = strings.TrimSuffix(subscriberChannel, "/")
	publishedChannel = strings.TrimSuffix(publishedChannel, "/")

	if subscriberChannel == publishedChannel {
		return true
	}

	if strings.HasSuffix(subscriberChannel, "/*") {
		prefix := strings.TrimSuffix(subscriberChannel, "/*")
		return publishedChannel == prefix || strings.HasPrefix(publishedChannel, prefix+"/")
	}

	return false
}

// EventServer handles WebSocket and HTTP publish endpoints for AppSync Events.
// Publishes events to the EventBus for cross-service visibility while maintaining
// an in-memory channelManager for in-process WebSocket fan-out delivery.
type EventServer struct {
	connections map[string]*wsConnection
	connMu      sync.RWMutex
	channels    *channelManager
	bus         eventbus.Bus
	storeLookup StoreLookupFunc
	sigVerifier *auth.SignatureV4Verifier
}

// StoreLookupFunc resolves an API ID to its store for auth-mode enforcement.
type StoreLookupFunc func(apiId string) (*appsyncstore.AppSyncStore, error)

// NewEventServer creates a new EventServer ready to accept connections.
func NewEventServer() *EventServer {
	return &EventServer{
		connections: make(map[string]*wsConnection),
		channels:    newChannelManager(),
	}
}

// SetEventBus injects the global event bus for cross-service event publishing.
// Events published via WebSocket or HTTP are forwarded to the bus so that
// other services (Lambda, EventBridge rules, etc.) can react to them.
func (s *EventServer) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// SetStoreLookup injects the store lookup function used for Event API
// auth-mode enforcement (R2-H1).
func (s *EventServer) SetStoreLookup(fn StoreLookupFunc) {
	s.storeLookup = fn
}

// SetSigVerifier injects the SigV4 verifier used for AWS_IAM auth mode
// enforcement on HTTP publish and WebSocket connection upgrade (P-5).
func (s *EventServer) SetSigVerifier(v *auth.SignatureV4Verifier) {
	s.sigVerifier = v
}

// authorizeEventOperation checks whether the caller is authenticated to
// perform a publish or subscribe operation on the given channel.
//
// Enforcement rules:
//   - If storeLookup is nil or apiId is empty → allow (test mode / no FQDN).
//   - If the API has no EventConfig or no auth modes configured → allow.
//   - Per-namespace auth modes take precedence over EventConfig defaults
//     when the channel falls under a namespace with explicit overrides.
//   - If auth modes are configured → verify credentials for each type:
//     API_KEY (store lookup + expiry), AWS_IAM (connection-level SigV4),
//     AMAZON_COGNITO_USER_POOLS (JWT via CognitoTokenValidator),
//     AWS_LAMBDA (Lambda authorizer invocation),
//     OPENID_CONNECT (fail-closed — out of scope per docs/services.md).
//
// Returns nil when authorised; a non-nil error string when denied.
func (s *EventServer) authorizeEventOperation(ctx context.Context, apiId, channel string, auth map[string]interface{}, isSubscribe bool, iamVerified bool) *string {
	if s.storeLookup == nil || apiId == "" {
		return nil
	}

	store, err := s.storeLookup(apiId)
	if err != nil || store == nil {
		return nil
	}

	api, apiErr := store.GetApiById(apiId)
	if apiErr != nil || api == nil || api.EventConfig == nil {
		return nil
	}

	// Determine the effective auth modes. If the channel falls under a
	// namespace with explicit auth-mode overrides, those take precedence;
	// otherwise fall back to the EventConfig defaults.
	authModes := resolveEffectiveAuthModes(store, apiId, channel, api.EventConfig, isSubscribe)

	if len(authModes) == 0 {
		return nil
	}

	for _, mode := range authModes {
		switch mode.AuthType {
		case "API_KEY":
			if s.verifyAPIKey(store, apiId, auth) {
				return nil
			}
		case "AWS_IAM":
			if iamVerified {
				return nil
			}
		case "AMAZON_COGNITO_USER_POOLS":
			if s.verifyCognito(ctx, api.EventConfig.AuthProviders, auth) {
				return nil
			}
		case "AWS_LAMBDA":
			if s.verifyLambdaAuthorizer(ctx, api.EventConfig.AuthProviders, apiId, auth) {
				return nil
			}
		case "OPENID_CONNECT":
			// OIDC token verification requires external IdP connectivity.
			// Out of scope per docs/services.md "No external IdP".
			// Fail-closed: do not accept unverified credentials.
		}
	}

	msg := "Unauthorized: valid authentication required for this operation"
	return &msg
}

// resolveEffectiveAuthModes returns the auth modes that apply to the given
// channel. Channels use /{namespace}/{path} format. If a ChannelNamespace
// with a matching name exists and has non-empty auth-mode overrides for the
// requested operation, those are used instead of the EventConfig defaults.
func resolveEffectiveAuthModes(store *appsyncstore.AppSyncStore, apiId, channel string, ec *appsyncstore.EventConfig, isSubscribe bool) []appsyncstore.AuthMode {
	// Default auth modes from EventConfig.
	var defaults []appsyncstore.AuthMode
	if isSubscribe {
		defaults = ec.DefaultSubscribeAuthModes
	} else {
		defaults = ec.DefaultPublishAuthModes
	}

	nsName := extractNamespaceFromChannel(channel)
	if nsName == "" {
		return defaults
	}

	ns, err := store.GetChannelNamespace(apiId, nsName)
	if err != nil || ns == nil {
		return defaults
	}

	if isSubscribe {
		if len(ns.SubscribeAuthModes) > 0 {
			return ns.SubscribeAuthModes
		}
	} else {
		if len(ns.PublishAuthModes) > 0 {
			return ns.PublishAuthModes
		}
	}

	return defaults
}

// extractNamespaceFromChannel extracts the namespace segment from a channel
// path. Channels use /{namespace}/{path...} format. Returns "" if the path
// is too short or the namespace segment is empty.
func extractNamespaceFromChannel(channel string) string {
	channel = strings.TrimPrefix(channel, "/")
	idx := strings.Index(channel, "/")
	if idx < 0 {
		return channel
	}
	return channel[:idx]
}

// verifyAPIKey looks up the provided API key in the store and checks expiry.
func (s *EventServer) verifyAPIKey(store *appsyncstore.AppSyncStore, apiId string, auth map[string]interface{}) bool {
	keyValue, ok := auth["x-api-key"].(string)
	if !ok || keyValue == "" {
		return false
	}
	apiKey, err := store.GetApiKey(apiId, keyValue)
	if err != nil {
		return false
	}
	if apiKey.Expires > 0 && time.Now().Unix() > apiKey.Expires {
		return false
	}
	return true
}

// verifyCognito validates a Cognito JWT access token using the eventbus
// CognitoTokenValidator. The AuthProvider with matching auth type provides
// the user pool ID and region.
func (s *EventServer) verifyCognito(ctx context.Context, providers []appsyncstore.AuthProvider, auth map[string]interface{}) bool {
	if s.bus == nil || s.bus.CognitoTokenValidator() == nil {
		return false
	}
	token := extractBearerToken(auth)
	if token == "" {
		return false
	}
	for _, p := range providers {
		if p.AuthType != "AMAZON_COGNITO_USER_POOLS" || p.CognitoConfig == nil {
			continue
		}
		if _, err := s.bus.CognitoTokenValidator().ValidateTokenForPool(ctx, p.CognitoConfig.AwsRegion, p.CognitoConfig.UserPoolId, token); err == nil {
			return true
		}
	}
	return false
}

// verifyLambdaAuthorizer invokes the configured Lambda authorizer function
// and evaluates the response for an authorisation decision.
func (s *EventServer) verifyLambdaAuthorizer(ctx context.Context, providers []appsyncstore.AuthProvider, apiId string, auth map[string]interface{}) bool {
	if s.bus == nil || s.bus.LambdaInvoker() == nil || len(auth) == 0 {
		return false
	}
	for _, p := range providers {
		if p.AuthType != "AWS_LAMBDA" || p.LambdaAuthorizerConfig == nil {
			continue
		}
		functionName := extractFunctionNameFromAuthorizerUri(p.LambdaAuthorizerConfig.AuthorizerUri)
		if functionName == "" {
			continue
		}
		payload, _ := json.Marshal(map[string]interface{}{
			"authorizationToken": auth,
			"requestContext":     map[string]string{"apiId": apiId},
		})
		_, resp, err := s.bus.LambdaInvoker().InvokeForGateway(ctx, functionName, payload)
		if err != nil {
			continue
		}
		var result struct {
			IsAuthorized *bool `json:"isAuthorized"`
		}
		if json.Unmarshal(resp, &result) == nil && result.IsAuthorized != nil && *result.IsAuthorized {
			return true
		}
	}
	return false
}

// extractBearerToken extracts the JWT from the Authorization header value,
// stripping the "Bearer " prefix if present.
func extractBearerToken(auth map[string]interface{}) string {
	v, ok := auth["Authorization"].(string)
	if !ok || v == "" {
		return ""
	}
	return strings.TrimPrefix(v, "Bearer ")
}

// extractFunctionNameFromAuthorizerUri extracts the Lambda function name
// from an AppSync LambdaAuthorizerConfig AuthorizerUri.
// Format: arn:aws:appsync:<region>:<account>:aws-lambda:<functionName>
// or: arn:aws:lambda:<region>:<account>:function:<functionName>
func extractFunctionNameFromAuthorizerUri(uri string) string {
	parts := strings.Split(uri, ":")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// lookupRegion resolves the store region for a given API ID.
func (s *EventServer) lookupRegion(apiId string) string {
	if s.storeLookup == nil {
		return ""
	}
	store, err := s.storeLookup(apiId)
	if err != nil || store == nil {
		return ""
	}
	return store.GetRegion()
}

// verifyConnectionIAM attempts SigV4 verification of the WebSocket upgrade
// request. Unlike the previous implementation that gated on
// ConnectionAuthModes, this always attempts verification so that the
// iamVerified flag is available for per-operation auth checks that may
// include namespace-level IAM overrides (which are only resolved inside
// authorizeEventOperation).
func (s *EventServer) verifyConnectionIAM(r *http.Request, apiId string) bool {
	if s.sigVerifier == nil || s.storeLookup == nil {
		return false
	}
	store, err := s.storeLookup(apiId)
	if err != nil || store == nil {
		return false
	}
	if err := s.sigVerifier.VerifyRequest(r, "appsync", store.GetRegion()); err != nil {
		logs.Warn("WebSocket IAM verification failed", logs.String("apiId", apiId), logs.Err(err))
		return false
	}
	return true
}

// given API ID. Call this when an API is deleted to prevent stale connections
// from receiving events for a non-existent API.
func (s *EventServer) DisconnectByApiId(apiId string) {
	s.connMu.RLock()
	var toClose []*wsConnection
	for _, ws := range s.connections {
		if apiId == "" || ws.apiId == apiId {
			toClose = append(toClose, ws)
		}
	}
	s.connMu.RUnlock()

	for _, ws := range toClose {
		ws.mu.Lock()
		if !ws.closed {
			ws.closed = true
			close(ws.sendCh)
		}
		ws.mu.Unlock()
	}
}

// RemoveSubscriptionsByNamespace removes all subscriptions whose channel
// falls under the given namespace. Channels use /{namespace}/{channel} format,
// so we match on segment boundaries to avoid over-matching (e.g. namespace
// "foo" must not affect "/foobar/baz").
func (s *EventServer) RemoveSubscriptionsByNamespace(namespace string) {
	if namespace == "" {
		return
	}
	exact := "/" + namespace
	prefix := exact + "/"
	s.connMu.RLock()
	var conns []*wsConnection
	for _, ws := range s.connections {
		conns = append(conns, ws)
	}
	s.connMu.RUnlock()

	for _, ws := range conns {
		ws.mu.Lock()
		for subId, sub := range ws.subscriptions {
			if sub.channel == exact || strings.HasPrefix(sub.channel, prefix) {
				delete(ws.subscriptions, subId)
				s.channels.unsubscribe(sub.channel, ws.id, subId)
			}
		}
		ws.mu.Unlock()
	}
}

// Shutdown closes all active WebSocket connections and cleans up internal state.
func (s *EventServer) Shutdown() {
	s.DisconnectByApiId("")
}

// ServeHTTP routes requests to either the WebSocket upgrade handler or the HTTP publish endpoint.
func (s *EventServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/event/realtime" && websocket.IsWebSocketUpgrade(r) {
		s.handleWebSocket(w, r)
	} else if r.URL.Path == "/event" && r.Method == http.MethodPost {
		s.handleHTTPPublish(w, r)
	} else {
		http.Error(w, `{"message":"Not found"}`, http.StatusNotFound)
	}
}

// handleWebSocket upgrades an HTTP connection to WebSocket and manages the connection lifecycle.
func (s *EventServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	apiId := s.extractApiId(r)

	iamVerified := false
	if s.sigVerifier != nil && apiId != "" && r.Header.Get("Authorization") != "" {
		iamVerified = s.verifyConnectionIAM(r, apiId)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("WebSocket upgrade failed", logs.Err(err))
		return
	}

	connId := uuid.New().String()

	ws := &wsConnection{
		id:            connId,
		apiId:         apiId,
		conn:          conn,
		sendCh:        make(chan []byte, 256),
		subscriptions: make(map[string]*subscription),
		iamVerified:   iamVerified,
	}

	s.connMu.Lock()
	s.connections[connId] = ws
	s.connMu.Unlock()

	logs.Info("AppSync WebSocket connected", logs.String("connId", connId), logs.String("apiId", ws.apiId))

	go s.writePump(ws)
	go s.readPump(ws)
}

// handleHTTPPublish processes a POST /event request, validates the payload,
// and broadcasts events to matching subscribers.
func (s *EventServer) handleHTTPPublish(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 10*maxEventSizeBytes))
	if err != nil {
		s.writeHTTPError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	// Restore r.Body so that SigV4 VerifyRequest can re-read the payload
	// to compute the correct body hash. Without this, VerifyRequest sees an
	// empty body (EOF) and the signature can never match.
	r.Body = io.NopCloser(bytes.NewReader(body))

	var req struct {
		Channel string   `json:"channel"`
		Events  []string `json:"events"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeHTTPError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.Channel == "" {
		s.writeHTTPError(w, http.StatusBadRequest, "channel is required")
		return
	}
	if !channelPathPattern.MatchString(req.Channel) {
		s.writeHTTPError(w, http.StatusBadRequest, "Invalid channel format")
		return
	}
	if len(req.Events) == 0 || len(req.Events) > maxEventsPerPublish {
		s.writeHTTPError(w, http.StatusBadRequest, "events must contain 1-5 items")
		return
	}

	for i, ev := range req.Events {
		if len(ev) > maxEventSizeBytes {
			s.writeHTTPError(w, http.StatusBadRequest, "Event at index %d exceeds 240KB", i)
			return
		}
	}

	// Enforce configured auth modes for HTTP publish.
	apiId := s.extractApiId(r)
	auth := map[string]interface{}{}
	if apiKey := r.Header.Get("x-api-key"); apiKey != "" {
		auth["x-api-key"] = apiKey
	}
	if authorization := r.Header.Get("Authorization"); authorization != "" {
		auth["Authorization"] = authorization
	}

	// Verify SigV4 signature when an Authorization header is present.
	// We always attempt verification rather than gating on a specific auth
	// mode, because namespace-level IAM overrides are only resolved later
	// inside authorizeEventOperation. If the client sent a valid SigV4
	// signature, iamVerified will be true and any IAM auth mode (default or
	// namespace-scoped) will accept it. If IAM is not among the effective
	// auth modes, the flag is simply ignored.
	iamVerified := false
	if s.sigVerifier != nil && apiId != "" && r.Header.Get("Authorization") != "" {
		if err := s.sigVerifier.VerifyRequest(r, "appsync", s.lookupRegion(apiId)); err == nil {
			iamVerified = true
		}
	}

	if msg := s.authorizeEventOperation(r.Context(), apiId, req.Channel, auth, false, iamVerified); msg != nil {
		s.writeHTTPError(w, http.StatusUnauthorized, "%s", *msg)
		return
	}

	result := s.publishEvents(req.Channel, req.Events)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// writeHTTPError writes a JSON error response to the HTTP client.
func (s *EventServer) writeHTTPError(w http.ResponseWriter, status int, format string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf(format, args...),
	})
}

// extractApiId attempts to extract the API ID from the WebSocket handshake.
// In production, the API ID comes from the Host header or subprotocol auth.
// For TEST_MODE, this is best-effort.
func (s *EventServer) extractApiId(r *http.Request) string {
	if id := fqdnrouter.ResourceIDFromContext(r.Context()); id != "" {
		return id
	}

	host := r.Host
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if len(parts) > 0 && parts[0] != "" && parts[0] != "localhost" && parts[0] != "127" {
		return parts[0]
	}

	// Try parsing from the subprotocol header
	for _, proto := range websocket.Subprotocols(r) {
		if strings.HasPrefix(proto, "header-") {
			encoded := strings.TrimPrefix(proto, "header-")
			if decoded, err := base64.RawURLEncoding.DecodeString(encoded); err == nil {
				var auth map[string]string
				if json.Unmarshal(decoded, &auth) == nil {
					if h, ok := auth["host"]; ok {
						hostParts := strings.Split(h, ".")
						if len(hostParts) > 0 {
							return hostParts[0]
						}
					}
				}
			}
		}
	}

	return ""
}

// readPump reads messages from the WebSocket connection and dispatches them.
// This goroutine exits when the connection is closed.
func (s *EventServer) readPump(ws *wsConnection) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		ws.mu.Lock()
		if !ws.closed {
			ws.closed = true
			close(ws.sendCh)
		}
		ws.mu.Unlock()

		s.channels.removeConnection(ws.id)

		s.connMu.Lock()
		delete(s.connections, ws.id)
		s.connMu.Unlock()

		ws.conn.Close()
		logs.Info("AppSync WebSocket disconnected", logs.String("connId", ws.id))
	}()

	ws.conn.SetReadLimit(maxEventSizeBytes * maxEventsPerPublish)

	// Send connection_ack immediately upon upgrade
	ack, _ := json.Marshal(map[string]interface{}{
		"type":                "connection_ack",
		"connectionTimeoutMs": connectionTimeoutMs,
	})
	ws.sendCh <- ack
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logs.Warn("WebSocket read error", logs.String("connId", ws.id), logs.Err(err))
			}
			return
		}

		s.handleMessage(ctx, ws, message)
	}
}

// writePump pumps messages from the send channel to the WebSocket connection.
// Also sends periodic keep-alive messages.
func (s *EventServer) writePump(ws *wsConnection) {
	keepAlive := time.NewTicker(keepAliveInterval)
	defer keepAlive.Stop()

	// Close the connection on exit so that readPump's ReadMessage unblocks and
	// enters its defer cleanup. Without this, a write error would leave readPump
	// running indefinitely on a half-dead connection.
	defer func() {
		ws.mu.Lock()
		if !ws.closed {
			ws.closed = true
			close(ws.sendCh)
		}
		ws.mu.Unlock()
		ws.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-ws.sendCh:
			if !ok {
				return
			}
			_ = ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-keepAlive.C:
			ws.mu.RLock()
			if ws.closed {
				ws.mu.RUnlock()
				return
			}
			ws.mu.RUnlock()

			ka, _ := json.Marshal(map[string]string{"type": "ka"})
			_ = ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.conn.WriteMessage(websocket.TextMessage, ka); err != nil {
				return
			}
		}
	}
}

// handleMessage dispatches an incoming WebSocket message based on its type.
func (s *EventServer) handleMessage(ctx context.Context, ws *wsConnection, raw []byte) {
	var msg struct {
		Type          string                 `json:"type"`
		Id            string                 `json:"id"`
		Channel       string                 `json:"channel"`
		Events        []string               `json:"events"`
		Authorization map[string]interface{} `json:"authorization"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		s.sendError(ws, "", "InvalidMessageFormat", "Invalid JSON message")
		return
	}

	switch msg.Type {
	case "connection_init":
		// Silently accepted — connection_ack already sent on upgrade

	case "subscribe":
		s.handleSubscribe(ctx, ws, msg.Id, msg.Channel, msg.Authorization)

	case "publish":
		s.handlePublish(ctx, ws, msg.Id, msg.Channel, msg.Events, msg.Authorization)

	case "unsubscribe":
		s.handleUnsubscribe(ws, msg.Id)

	default:
		s.sendError(ws, msg.Id, "UnknownOperationError", "Unknown message type: "+msg.Type)
	}
}

// handleSubscribe validates and registers a channel subscription.
func (s *EventServer) handleSubscribe(ctx context.Context, ws *wsConnection, subId, channel string, auth map[string]interface{}) {
	if subId == "" || !subscriptionIDPattern.MatchString(subId) {
		s.sendSubscribeError(ws, subId, "InvalidInput", "Invalid subscription ID format")
		return
	}

	if channel == "" || !channelPathPattern.MatchString(channel) {
		s.sendSubscribeError(ws, subId, "InvalidInput", "Invalid channel format")
		return
	}

	// Enforce configured auth modes.
	if msg := s.authorizeEventOperation(ctx, ws.apiId, channel, auth, true, ws.iamVerified); msg != nil {
		s.sendSubscribeError(ws, subId, "UnauthorizedException", *msg)
		return
	}

	ws.mu.Lock()
	if len(ws.subscriptions) >= maxSubscriptionsPerConn {
		ws.mu.Unlock()
		s.sendSubscribeError(ws, subId, "LimitExceededException", "Maximum subscriptions per connection exceeded")
		return
	}
	if _, exists := ws.subscriptions[subId]; exists {
		ws.mu.Unlock()
		s.sendSubscribeError(ws, subId, "ConflictException", "Subscription ID already exists on this connection")
		return
	}

	sub := &subscription{
		id:      subId,
		channel: channel,
		auth:    auth,
	}
	ws.subscriptions[subId] = sub
	ws.mu.Unlock()

	s.channels.subscribe(channel, ws.id, subId)

	resp, _ := json.Marshal(map[string]string{
		"type": "subscribe_success",
		"id":   subId,
	})
	s.sendControlMessage(ws, resp)
}

// handlePublish validates events and broadcasts them to matching subscribers.
func (s *EventServer) handlePublish(ctx context.Context, ws *wsConnection, pubId, channel string, events []string, auth map[string]interface{}) {
	if pubId == "" || !subscriptionIDPattern.MatchString(pubId) {
		s.sendPublishError(ws, pubId, "InvalidInput", "Invalid publish ID format")
		return
	}

	if channel == "" || !channelPathPattern.MatchString(channel) {
		s.sendPublishError(ws, pubId, "InvalidInput", "Invalid channel format")
		return
	}

	// Enforce configured auth modes.
	if msg := s.authorizeEventOperation(ctx, ws.apiId, channel, auth, false, ws.iamVerified); msg != nil {
		s.sendPublishError(ws, pubId, "UnauthorizedException", *msg)
		return
	}

	if len(events) == 0 || len(events) > maxEventsPerPublish {
		s.sendPublishError(ws, pubId, "InvalidInput", "events must contain 1-5 items")
		return
	}

	for i, ev := range events {
		if len(ev) > maxEventSizeBytes {
			s.sendPublishError(ws, pubId, "LimitExceededException", fmt.Sprintf("Event at index %d exceeds 240KB", i))
			return
		}
	}

	result := s.publishEvents(channel, events)

	resp, _ := json.Marshal(map[string]interface{}{
		"type":       "publish_success",
		"id":         pubId,
		"successful": result.Successful,
		"failed":     result.Failed,
	})
	s.sendControlMessage(ws, resp)
}

// handleUnsubscribe removes a subscription and stops receiving events on that channel.
func (s *EventServer) handleUnsubscribe(ws *wsConnection, subId string) {
	if subId == "" {
		s.sendUnsubscribeError(ws, subId, "UnknownOperationError", "Missing subscription ID")
		return
	}

	ws.mu.Lock()
	sub, exists := ws.subscriptions[subId]
	if !exists {
		ws.mu.Unlock()
		s.sendUnsubscribeError(ws, subId, "UnknownOperationError", "Unknown operation id "+subId)
		return
	}

	delete(ws.subscriptions, subId)
	ws.mu.Unlock()

	s.channels.unsubscribe(sub.channel, ws.id, subId)

	resp, _ := json.Marshal(map[string]string{
		"type": "unsubscribe_success",
		"id":   subId,
	})
	s.sendControlMessage(ws, resp)
}

// publishEvents broadcasts events to all subscribers matching the channel path.
// Returns the publish result with per-event identifiers.
func (s *EventServer) publishEvents(channel string, events []string) *publishResult {
	result := &publishResult{}

	// Double-encode events as per the protocol specification:
	// The "event" field in data messages is a JSON string containing a JSON array.
	eventsJSON, _ := json.Marshal(events)
	eventString := string(eventsJSON)

	matches := s.channels.matchSubscriptions(channel)
	for _, match := range matches {
		s.connMu.RLock()
		ws, ok := s.connections[match.connId]
		s.connMu.RUnlock()
		if !ok {
			continue
		}

		dataMsg, _ := json.Marshal(map[string]string{
			"type":  "data",
			"id":    match.subId,
			"event": eventString,
		})

		s.sendMessage(ws, dataMsg)
	}

	for i := range events {
		result.Successful = append(result.Successful, struct {
			Identifier string `json:"identifier"`
			Index      int    `json:"index"`
		}{
			Identifier: uuid.New().String(),
			Index:      i,
		})
	}

	s.publishToBus(channel, events)

	return result
}

// publishToBus forwards published events to the EventBus for cross-service
// consumption. Fire-and-forget — failures are logged but do not affect the
// publish response to the WebSocket/HTTP client.
func (s *EventServer) publishToBus(channel string, events []string) {
	if s.bus == nil {
		return
	}

	ctx := context.Background()
	for _, ev := range events {
		publishEvent := &AppSyncEventPublished{
			EventBase: eventbus.EventBase{
				ID:        uuid.New().String(),
				Timestamp: time.Now().UTC(),
				Source:    "aws.appsync",
				Region:    "",
				AccountID: "",
			},
			Channel: channel,
			Event:   ev,
		}
		if err := s.bus.Publish(ctx, publishEvent); err != nil {
			logs.Warn("Failed to publish AppSync event to bus",
				logs.String("channel", channel),
				logs.Err(err))
		}
	}
}

// AppSyncEventPublished is published to the EventBus when an event is published
// via the AppSync Events WebSocket or HTTP publish endpoint. Other services
// (Lambda, EventBridge rules, etc.) can subscribe to react to these events.
type AppSyncEventPublished struct {
	eventbus.EventBase
	Channel string `json:"channel"`
	Event   string `json:"event"`
}

// EventType returns the event bus identifier for AppSync publish events.
func (e *AppSyncEventPublished) EventType() string { return "appsync:event-published" }

// publishResult holds the outcome of a publish operation.
type publishResult struct {
	Successful []struct {
		Identifier string `json:"identifier"`
		Index      int    `json:"index"`
	} `json:"successful"`
	Failed []interface{} `json:"failed"`
}

func (s *EventServer) sendMessage(ws *wsConnection, msg []byte) {
	s.sendMessageInternal(ws, msg, false)
}

func (s *EventServer) sendControlMessage(ws *wsConnection, msg []byte) {
	s.sendMessageInternal(ws, msg, true)
}

func (s *EventServer) sendMessageInternal(ws *wsConnection, msg []byte, block bool) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ws.closed {
		return
	}
	// Always use non-blocking send. The previous blocking path (when block=true)
	// could deadlock: if writePump exits after a write error and sendCh is full,
	// readPump would block forever inside handleMessage -> sendControlMessage,
	// preventing the defer cleanup from ever running.
	select {
	case ws.sendCh <- msg:
	default:
		logs.Warn("WebSocket send buffer full, dropping message",
			logs.String("connId", ws.id))
	}
}

func (s *EventServer) sendError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendControlMessage(ws, resp)
}

func (s *EventServer) sendSubscribeError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "subscribe_error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendControlMessage(ws, resp)
}

func (s *EventServer) sendPublishError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "publish_error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendControlMessage(ws, resp)
}

func (s *EventServer) sendUnsubscribeError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "unsubscribe_error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendControlMessage(ws, resp)
}
