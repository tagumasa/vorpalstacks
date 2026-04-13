package appsync

import (
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

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
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
	CheckOrigin:     func(r *http.Request) bool { return true },
	Subprotocols:    []string{wsSubprotocol},
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
	// channel -> set of (connId, subId) pairs
	channels map[string]map[string]string
	mu       sync.RWMutex
}

func newChannelManager() *channelManager {
	return &channelManager{
		channels: make(map[string]map[string]string),
	}
}

// subscribe registers a subscription for broadcasting.
// Returns the set of concrete channel paths that match (expanding wildcards).
func (cm *channelManager) subscribe(channel string, connId, subId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.channels[channel] == nil {
		cm.channels[channel] = make(map[string]string)
	}
	cm.channels[channel][connId] = subId
}

// unsubscribe removes a subscription from broadcasting.
func (cm *channelManager) unsubscribe(channel string, connId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if subs, ok := cm.channels[channel]; ok {
		delete(subs, connId)
		if len(subs) == 0 {
			delete(cm.channels, channel)
		}
	}
}

// removeConnection removes all subscriptions for a given connection ID.
func (cm *channelManager) removeConnection(connId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for ch, subs := range cm.channels {
		delete(subs, connId)
		if len(subs) == 0 {
			delete(cm.channels, ch)
		}
	}
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
			for connId, subId := range subs {
				matches = append(matches, struct {
					connId string
					subId  string
				}{connId: connId, subId: subId})
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
}

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

// DisconnectByApiId closes all WebSocket connections associated with the
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
		ws.conn.Close()
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("WebSocket upgrade failed", logs.Err(err))
		return
	}

	connId := uuid.New().String()

	ws := &wsConnection{
		id:            connId,
		apiId:         s.extractApiId(r),
		conn:          conn,
		sendCh:        make(chan []byte, 256),
		subscriptions: make(map[string]*subscription),
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

	result := s.publishEvents(req.Channel, req.Events)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// writeHTTPError writes a JSON error response to the HTTP client.
func (s *EventServer) writeHTTPError(w http.ResponseWriter, status int, format string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf(format, args...),
	})
}

// extractApiId attempts to extract the API ID from the WebSocket handshake.
// In production, the API ID comes from the Host header or subprotocol auth.
// For TEST_MODE, this is best-effort.
func (s *EventServer) extractApiId(r *http.Request) string {
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
	defer func() {
		ws.mu.Lock()
		ws.closed = true
		close(ws.sendCh)
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

		s.handleMessage(ws, message)
	}
}

// writePump pumps messages from the send channel to the WebSocket connection.
// Also sends periodic keep-alive messages.
func (s *EventServer) writePump(ws *wsConnection) {
	keepAlive := time.NewTicker(keepAliveInterval)
	defer keepAlive.Stop()

	for {
		select {
		case msg, ok := <-ws.sendCh:
			if !ok {
				ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
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
			ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.conn.WriteMessage(websocket.TextMessage, ka); err != nil {
				return
			}
		}
	}
}

// handleMessage dispatches an incoming WebSocket message based on its type.
func (s *EventServer) handleMessage(ws *wsConnection, raw []byte) {
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
		s.handleSubscribe(ws, msg.Id, msg.Channel, msg.Authorization)

	case "publish":
		s.handlePublish(ws, msg.Id, msg.Channel, msg.Events, msg.Authorization)

	case "unsubscribe":
		s.handleUnsubscribe(ws, msg.Id)

	default:
		s.sendError(ws, msg.Id, "UnknownOperationError", "Unknown message type: "+msg.Type)
	}
}

// handleSubscribe validates and registers a channel subscription.
func (s *EventServer) handleSubscribe(ws *wsConnection, subId, channel string, auth map[string]interface{}) {
	if subId == "" || !subscriptionIDPattern.MatchString(subId) {
		s.sendSubscribeError(ws, subId, "InvalidInput", "Invalid subscription ID format")
		return
	}

	if channel == "" || !channelPathPattern.MatchString(channel) {
		s.sendSubscribeError(ws, subId, "InvalidInput", "Invalid channel format")
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
	s.sendMessage(ws, resp)
}

// handlePublish validates events and broadcasts them to matching subscribers.
func (s *EventServer) handlePublish(ws *wsConnection, pubId, channel string, events []string, auth map[string]interface{}) {
	if pubId == "" || !subscriptionIDPattern.MatchString(pubId) {
		s.sendPublishError(ws, pubId, "InvalidInput", "Invalid publish ID format")
		return
	}

	if channel == "" || !channelPathPattern.MatchString(channel) {
		s.sendPublishError(ws, pubId, "InvalidInput", "Invalid channel format")
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
	s.sendMessage(ws, resp)
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

	s.channels.unsubscribe(sub.channel, ws.id)

	resp, _ := json.Marshal(map[string]string{
		"type": "unsubscribe_success",
		"id":   subId,
	})
	s.sendMessage(ws, resp)
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
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ws.closed {
		return
	}
	select {
	case ws.sendCh <- msg:
	default:
		logs.Warn("WebSocket send buffer full, dropping message", logs.String("connId", ws.id))
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
	s.sendMessage(ws, resp)
}

func (s *EventServer) sendSubscribeError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "subscribe_error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendMessage(ws, resp)
}

func (s *EventServer) sendPublishError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "publish_error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendMessage(ws, resp)
}

func (s *EventServer) sendUnsubscribeError(ws *wsConnection, id, errorType, message string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"type": "unsubscribe_error",
		"id":   id,
		"errors": []map[string]string{
			{"errorType": errorType, "message": message},
		},
	})
	s.sendMessage(ws, resp)
}
