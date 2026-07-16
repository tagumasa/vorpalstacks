package hooks

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	m := NewManager()
	defer m.Clear()

	t.Run("Register and Run", func(t *testing.T) {
		executed := false
		m.Register("test.event", func(ctx context.Context) error {
			executed = true
			return nil
		})

		err := m.Run(context.Background(), "test.event")
		require.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("Priority Order", func(t *testing.T) {
		order := []int{}
		m.Register("priority.test", func(ctx context.Context) error {
			order = append(order, 1)
			return nil
		}, WithPriority(1))

		m.Register("priority.test", func(ctx context.Context) error {
			order = append(order, 3)
			return nil
		}, WithPriority(3))

		m.Register("priority.test", func(ctx context.Context) error {
			order = append(order, 2)
			return nil
		}, WithPriority(2))

		err := m.Run(context.Background(), "priority.test")
		require.NoError(t, err)
		assert.Equal(t, []int{3, 2, 1}, order)
	})

	t.Run("Error Propagation", func(t *testing.T) {
		expectedErr := errors.New("hook error")
		m.Register("error.test", func(ctx context.Context) error {
			return expectedErr
		})

		err := m.Run(context.Background(), "error.test")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Empty Event", func(t *testing.T) {
		err := m.Run(context.Background(), "nonexistent.event")
		assert.NoError(t, err)
	})

	t.Run("Clear Event", func(t *testing.T) {
		m.Register("clear.test", func(ctx context.Context) error {
			return nil
		})
		m.ClearEvent("clear.test")
		hooks := m.GetHooks("clear.test")
		assert.Empty(t, hooks)
	})
}

func TestGlobalManager(t *testing.T) {
	Clear()
	defer Clear()

	executed := false
	Register("global.test", func(ctx context.Context) error {
		executed = true
		return nil
	})

	err := Run(context.Background(), "global.test")
	require.NoError(t, err)
	assert.True(t, executed)
}

func TestOn(t *testing.T) {
	Clear()
	defer Clear()

	executed := false
	On("on.test", func(ctx context.Context) error {
		executed = true
		return nil
	})

	err := Run(context.Background(), "on.test")
	require.NoError(t, err)
	assert.True(t, executed)
}

func TestRunAsync(t *testing.T) {
	m := NewManager()
	defer m.Clear()

	executed := false
	m.Register("async.test", func(ctx context.Context) error {
		executed = true
		return nil
	})

	errCh := m.RunAsync(context.Background(), "async.test")
	err := <-errCh
	require.NoError(t, err)
	assert.True(t, executed)
}

func TestEventNames(t *testing.T) {
	assert.NotEmpty(t, EventNames)
	assert.Contains(t, EventNames, EventBeforeInit)
	assert.Contains(t, EventNames, EventAfterStart)
	assert.Contains(t, EventNames, EventBeforeStop)
}
