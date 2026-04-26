package appsync

import (
	"testing"

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

	matches := cm.matchSubscriptions("ch1")
	assert.Equal(t, 2, len(matches), "both subscriptions should be tracked (Bug A-1)")

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
