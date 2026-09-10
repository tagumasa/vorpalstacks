package appsync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannelManager_SubscribeUnsubscribe(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "conn1", matches[0].connId)
	assert.Equal(t, "sub1", matches[0].subId)
}

func TestChannelManager_SameConnDifferentSubIds(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	cm.subscribe("ch1", "conn1", "sub2")

	// A connection that subscribes twice under different subscription
	// ids must be tracked twice so that unsubscribing one id leaves the
	// other intact.
	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 2, len(matches), "both subscriptions should be tracked")

	cm.unsubscribe("ch1", "conn1", "sub1")
	matches = cm.matchSubscriptions("ch1")
	assert.Equal(t, 1, len(matches), "only sub2 should remain after unsubscribing sub1")
	assert.Equal(t, "sub2", matches[0].subId)
}

func TestChannelManager_DifferentChannels(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	cm.subscribe("ch2", "conn1", "sub2")

	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "sub1", matches[0].subId)

	matches = cm.matchSubscriptions("ch2")
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "sub2", matches[0].subId)
}

func TestChannelManager_RemoveConnection(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	cm.subscribe("ch2", "conn1", "sub2")
	cm.subscribe("ch1", "conn2", "sub3")

	cm.removeConnection("conn1")

	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 1, len(matches), "only conn2/sub3 should remain on ch1")
	assert.Equal(t, "conn2", matches[0].connId)

	matches = cm.matchSubscriptions("ch2")
	assert.Equal(t, 0, len(matches), "conn1/sub2 should be gone from ch2")
}

func TestChannelManager_RemoveConnection_CleansUpEmptyChannels(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	cm.removeConnection("conn1")

	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 0, len(matches))
}

func TestChannelManager_UnsubscribeNonExistent(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	cm.unsubscribe("ch1", "conn99", "sub99")
	cm.unsubscribe("ch99", "conn1", "sub1")

	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 1, len(matches), "original subscription should be unaffected")
}

func TestChannelManager_RemoveConnectionNonExistent(t *testing.T) {
	cm := newChannelManager()

	cm.subscribe("ch1", "conn1", "sub1")
	cm.removeConnection("conn99")

	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 1, len(matches))
}

func TestChannelMatches(t *testing.T) {
	tests := []struct {
		name     string
		subCh    string
		pubCh    string
		expected bool
	}{
		{"exact match", "/channel/a", "/channel/a", true},
		{"no match", "/channel/a", "/channel/b", false},
		{"wildcard matches exact prefix", "/channel/*", "/channel/a", true},
		{"wildcard matches prefix path", "/channel/*", "/channel/a/b", true},
		{"wildcard no match different prefix", "/channel/*", "/other/a", false},
		{"trailing slash normalised", "/channel/a/", "/channel/a", true},
		{"trailing slash normalised reverse", "/channel/a", "/channel/a/", true},
		{"empty channels", "", "", true},
		{"root wildcard", "/*", "/anything", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := channelMatches(tt.subCh, tt.pubCh)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAPIKeyExpired(t *testing.T) {
	assert.False(t, apiKeyExpired(0), "zero expiry means the key never expires")
	assert.False(t, apiKeyExpired(time.Now().Unix()+3600), "future expiry is not expired")
	assert.True(t, apiKeyExpired(time.Now().Unix()-3600), "past expiry is expired")
}

func TestExtractSubprotocolAuth(t *testing.T) {
	creds := base64.RawURLEncoding.EncodeToString([]byte(`{"host":"api123.example.com","x-api-key":"da2-test"}`))

	joined := httptest.NewRequest("GET", "/event/realtime", nil)
	joined.Header.Set("Sec-WebSocket-Protocol", "aws-appsync-event-ws, header-"+creds)
	auth := extractSubprotocolAuth(joined)
	assert.NotNil(t, auth)
	assert.Equal(t, "api123.example.com", auth["host"])
	assert.Equal(t, "da2-test", auth["x-api-key"])

	repeated := httptest.NewRequest("GET", "/event/realtime", nil)
	repeated.Header.Add("Sec-WebSocket-Protocol", "aws-appsync-event-ws")
	repeated.Header.Add("Sec-WebSocket-Protocol", "header-"+creds)
	auth = extractSubprotocolAuth(repeated)
	assert.NotNil(t, auth, "repeated header lines carry the same protocol list")
	assert.Equal(t, "da2-test", auth["x-api-key"])

	bare := httptest.NewRequest("GET", "/event/realtime", nil)
	bare.Header.Set("Sec-WebSocket-Protocol", "aws-appsync-event-ws")
	assert.Nil(t, extractSubprotocolAuth(bare), "no header- subprotocol yields no credentials")
}

func TestVerifyConnectionAuthFailClosed(t *testing.T) {
	s := NewEventServer()

	assert.False(t, s.verifyConnectionAuth(context.Background(), "", nil), "missing apiId fails closed")
	assert.False(t, s.verifyConnectionAuth(context.Background(), "api123", nil), "missing subprotocol credentials fail closed")
	assert.False(t, s.verifyConnectionAuth(context.Background(), "api123", map[string]string{"host": "h", "x-api-key": "k"}),
		"credentials without a reachable API configuration fail closed")
}

func TestPublishEventsPerEventDataMessages(t *testing.T) {
	s := NewEventServer()
	ws := &wsConnection{
		id:            "conn1",
		sendCh:        make(chan []byte, 8),
		subscriptions: make(map[string]*subscription),
	}
	s.connections["conn1"] = ws
	s.channels.subscribe("/default/ch", "conn1", "sub1")

	result := s.publishEvents("/default/ch", []string{`{"msg":"first"}`, `{"msg":"second"}`})

	// One data message per event; the event member is the documented array
	// form carrying the published stringified event.
	var got [][]string
	for i := 0; i < 2; i++ {
		select {
		case raw := <-ws.sendCh:
			var m struct {
				Type  string   `json:"type"`
				Id    string   `json:"id"`
				Event []string `json:"event"`
			}
			if err := json.Unmarshal(raw, &m); err != nil {
				t.Fatalf("data message %d is not valid JSON: %v", i, err)
			}
			assert.Equal(t, "data", m.Type)
			assert.Equal(t, "sub1", m.Id)
			got = append(got, m.Event)
		case <-time.After(time.Second):
			t.Fatalf("data message %d not delivered", i)
		}
	}
	assert.Equal(t, [][]string{{`{"msg":"first"}`}, {`{"msg":"second"}`}}, got)
	assert.Len(t, result.Successful, 2, "both events reported successful")

	select {
	case raw := <-ws.sendCh:
		t.Fatalf("unexpected extra message after the per-event deliveries: %s", raw)
	default:
	}
}
