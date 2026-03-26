package events

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryBus(t *testing.T) {
	bus := NewInMemoryBus()

	t.Run("Subscribe and Publish", func(t *testing.T) {
		received := false
		subID, err := bus.Subscribe("test.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			received = true
			assert.Equal(t, "test.event", event.Type)
			return nil
		}))
		require.NoError(t, err)
		assert.NotEmpty(t, subID)

		err = bus.Publish(context.Background(), &Event{Type: "test.event"})
		require.NoError(t, err)
		assert.True(t, received)
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		received := false
		subID, err := bus.Subscribe("unsub.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			received = true
			return nil
		}))
		require.NoError(t, err)

		err = bus.Unsubscribe(subID)
		require.NoError(t, err)

		err = bus.Publish(context.Background(), &Event{Type: "unsub.event"})
		require.NoError(t, err)
		assert.False(t, received)
	})

	t.Run("Multiple Subscribers", func(t *testing.T) {
		var mu sync.Mutex
		count := 0
		bus.Subscribe("multi.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			mu.Lock()
			count++
			mu.Unlock()
			return nil
		}))
		bus.Subscribe("multi.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			mu.Lock()
			count++
			mu.Unlock()
			return nil
		}))

		err := bus.Publish(context.Background(), &Event{Type: "multi.event"})
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("Handler Error", func(t *testing.T) {
		bus.Subscribe("error.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			return errors.New("handler error")
		}))

		err := bus.Publish(context.Background(), &Event{Type: "error.event"})
		assert.Error(t, err)
	})

	t.Run("Nil Event", func(t *testing.T) {
		err := bus.Publish(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("Empty EventType", func(t *testing.T) {
		_, err := bus.Subscribe("", HandlerFunc(func(ctx context.Context, event *Event) error {
			return nil
		}))
		assert.Error(t, err)
	})

	t.Run("Nil Handler", func(t *testing.T) {
		_, err := bus.Subscribe("nil.handler", nil)
		assert.Error(t, err)
	})

	t.Run("No Subscribers", func(t *testing.T) {
		err := bus.Publish(context.Background(), &Event{Type: "no.subscribers"})
		assert.NoError(t, err)
	})

	t.Run("GetSubscriberCount", func(t *testing.T) {
		bus.Subscribe("count.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			return nil
		}))
		bus.Subscribe("count.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			return nil
		}))

		count := bus.GetSubscriberCount("count.event")
		assert.Equal(t, 2, count)
	})

	t.Run("ClearAllSubscribers", func(t *testing.T) {
		bus.Subscribe("clear.event", HandlerFunc(func(ctx context.Context, event *Event) error {
			return nil
		}))
		bus.ClearAllSubscribers()

		count := bus.GetSubscriberCount("clear.event")
		assert.Equal(t, 0, count)
	})
}

func TestEventDefaults(t *testing.T) {
	bus := NewInMemoryBus()
	var receivedEvent *Event

	bus.Subscribe("defaults.event", HandlerFunc(func(ctx context.Context, event *Event) error {
		receivedEvent = event
		return nil
	}))

	bus.Publish(context.Background(), &Event{Type: "defaults.event"})

	require.NotNil(t, receivedEvent)
	assert.NotEmpty(t, receivedEvent.ID)
	assert.False(t, receivedEvent.Timestamp.IsZero())
}
