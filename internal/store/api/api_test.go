package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
)

func newTestStorage(t *testing.T) storage.Storage {
	tmpDir := filepath.Join(os.TempDir(), "api-store-test")
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	return s
}

func newTestStorageForSubtest(t *testing.T) storage.Storage {
	tmpDir := filepath.Join(os.TempDir(), "api-store-test-"+t.Name())
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	return s
}

func TestConfigStore(t *testing.T) {
	s := newTestStorage(t)
	defer s.Close()

	store := NewConfigStore(s)

	t.Run("Get returns default config for non-existent key", func(t *testing.T) {
		cfg, err := store.Get("nonexistent-unique")
		require.NoError(t, err)
		assert.Equal(t, "nonexistent-unique", cfg.ServiceName)
		assert.Equal(t, ServiceModeFallback, cfg.Mode)
		assert.True(t, cfg.Enabled)
	})

	t.Run("Put and Get", func(t *testing.T) {
		cfg := NewServiceConfig("test-service-unique", ServiceModeImplemented)
		cfg.Enabled = false
		err := store.Put(cfg)
		require.NoError(t, err)

		got, err := store.Get("test-service-unique")
		require.NoError(t, err)
		assert.Equal(t, "test-service-unique", got.ServiceName)
		assert.Equal(t, ServiceModeImplemented, got.Mode)
		assert.False(t, got.Enabled)
	})

	t.Run("Delete", func(t *testing.T) {
		cfg := NewServiceConfig("to-delete-unique", ServiceModeFallback)
		err := store.Put(cfg)
		require.NoError(t, err)

		err = store.Delete("to-delete-unique")
		require.NoError(t, err)

		got, err := store.Get("to-delete-unique")
		require.NoError(t, err)
		assert.Equal(t, ServiceModeFallback, got.Mode)
	})

	t.Run("SetMode", func(t *testing.T) {
		cfg := NewServiceConfig("mode-test-unique", ServiceModeFallback)
		err := store.Put(cfg)
		require.NoError(t, err)

		err = store.SetMode("mode-test-unique", ServiceModeImplemented)
		require.NoError(t, err)

		got, err := store.Get("mode-test-unique")
		require.NoError(t, err)
		assert.Equal(t, ServiceModeImplemented, got.Mode)
	})

	t.Run("SetError", func(t *testing.T) {
		cfg := NewServiceConfig("error-test-unique", ServiceModeFallback)
		err := store.Put(cfg)
		require.NoError(t, err)

		errCfg := &ErrorConfig{
			HTTPStatusCode: 500,
			Code:           "InternalError",
			Message:        "test error",
		}
		err = store.SetError("error-test-unique", errCfg)
		require.NoError(t, err)

		got, err := store.Get("error-test-unique")
		require.NoError(t, err)
		assert.Equal(t, ServiceModeErrorInjection, got.Mode)
		assert.NotNil(t, got.Error)
		assert.Equal(t, 500, got.Error.HTTPStatusCode)
	})

	t.Run("SetEnabled", func(t *testing.T) {
		cfg := NewServiceConfig("enabled-test-unique", ServiceModeImplemented)
		cfg.Enabled = true
		err := store.Put(cfg)
		require.NoError(t, err)

		err = store.SetEnabled("enabled-test-unique", false)
		require.NoError(t, err)

		got, err := store.Get("enabled-test-unique")
		require.NoError(t, err)
		assert.False(t, got.Enabled)
	})

	t.Run("ForEach", func(t *testing.T) {
		store.Put(NewServiceConfig("foreach1-unique", ServiceModeImplemented))
		store.Put(NewServiceConfig("foreach2-unique", ServiceModeFallback))

		svcs := []string{}
		err := store.ForEach(func(c *ServiceConfig) error {
			svcs = append(svcs, c.ServiceName)
			return nil
		})
		require.NoError(t, err)
		assert.Contains(t, svcs, "foreach1-unique")
		assert.Contains(t, svcs, "foreach2-unique")
	})

	t.Run("Count", func(t *testing.T) {
		store.Put(NewServiceConfig("count1-unique", ServiceModeImplemented))
		store.Put(NewServiceConfig("count2-unique", ServiceModeFallback))

		assert.GreaterOrEqual(t, store.Count(), 2)
	})
}

func TestServiceStore(t *testing.T) {
	s := newTestStorage(t)
	defer s.Close()

	store := NewServiceStore(s)

	t.Run("Get returns nil for non-existent", func(t *testing.T) {
		svc, err := store.Get("nonexistent-unique")
		require.NoError(t, err)
		assert.Nil(t, svc)
	})

	t.Run("Put and Get", func(t *testing.T) {
		svc := &Service{
			ID:                    1,
			ShapeID:               "test_shape_unique",
			Name:                  "test-service-unique",
			Version:               "2023-01-01",
			SDKID:                 "TestService",
			ARNNamespace:          "test",
			CloudFormationName:    "TestService",
			CloudTrailEventSource: "testservice.amazonaws.com",
			EndpointPrefix:        "test",
			Protocol:              "json",
		}
		err := store.Put(svc)
		require.NoError(t, err)

		got, err := store.Get("test-service-unique")
		require.NoError(t, err)
		assert.Equal(t, "test-service-unique", got.Name)
		assert.Equal(t, "test_shape_unique", got.ShapeID)
	})

	t.Run("GetByShapeID", func(t *testing.T) {
		svc := &Service{
			ID:      2,
			ShapeID: "shape-123-unique",
			Name:    "shape-service-unique",
		}
		err := store.Put(svc)
		require.NoError(t, err)

		got, err := store.GetByShapeID("shape-123-unique")
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "shape-service-unique", got.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		svc := &Service{
			ID:   3,
			Name: "to-delete-unique",
		}
		err := store.Put(svc)
		require.NoError(t, err)

		err = store.Delete("to-delete-unique")
		require.NoError(t, err)

		got, err := store.Get("to-delete-unique")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("ForEach", func(t *testing.T) {
		store.Put(&Service{ID: 1, Name: "foreach1-unique", ShapeID: "s1"})
		store.Put(&Service{ID: 2, Name: "foreach2-unique", ShapeID: "s2"})

		names := []string{}
		err := store.ForEach(func(svc *Service) error {
			names = append(names, svc.Name)
			return nil
		})
		require.NoError(t, err)
		assert.Contains(t, names, "foreach1-unique")
		assert.Contains(t, names, "foreach2-unique")
	})

	t.Run("Count", func(t *testing.T) {
		store.Put(&Service{ID: 1, Name: "count1-unique", ShapeID: "s1"})
		store.Put(&Service{ID: 2, Name: "count2-unique", ShapeID: "s2"})

		assert.GreaterOrEqual(t, store.Count(), 2)
	})
}

func TestOperationStore(t *testing.T) {
	s := newTestStorage(t)
	defer s.Close()

	store := NewOperationStore(s)

	t.Run("Get returns nil for non-existent", func(t *testing.T) {
		op, err := store.Get("svc-unique", "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, op)
	})

	t.Run("Put and Get", func(t *testing.T) {
		op := &Operation{
			ID:             1,
			ShapeID:        "op_shape_unique",
			ServiceID:      100,
			Name:           "TestOperationUnique",
			IsReadonly:     true,
			XAmznQueryMode: false,
		}
		err := store.Put("test-service-unique", op)
		require.NoError(t, err)

		got, err := store.Get("test-service-unique", "TestOperationUnique")
		require.NoError(t, err)
		assert.Equal(t, "TestOperationUnique", got.Name)
		assert.True(t, got.IsReadonly)
	})

	t.Run("GetByShapeID", func(t *testing.T) {
		op := &Operation{
			ID:      2,
			ShapeID: "op-shape-123-unique",
			Name:    "GetShapeUnique",
		}
		err := store.Put("shape-svc-unique", op)
		require.NoError(t, err)

		got, err := store.GetByShapeID("op-shape-123-unique")
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "GetShapeUnique", got.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		op := &Operation{
			ID:   3,
			Name: "ToDeleteUnique",
		}
		err := store.Put("del-svc-unique", op)
		require.NoError(t, err)

		err = store.Delete("del-svc-unique", "ToDeleteUnique")
		require.NoError(t, err)

		got, err := store.Get("del-svc-unique", "ToDeleteUnique")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("ForEachByService", func(t *testing.T) {
		s := newTestStorageForSubtest(t)
		defer s.Close()
		localStore := NewOperationStore(s)

		localStore.Put("forEachByServiceSVC", &Operation{ID: 1, Name: "A1"})
		localStore.Put("forEachByServiceSVC", &Operation{ID: 2, Name: "A2"})
		localStore.Put("otherSVC", &Operation{ID: 3, Name: "A3"})

		count := 0
		err := localStore.ForEachByService("forEachByServiceSVC", func(op *Operation) error {
			count++
			return nil
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 2)
	})

	t.Run("Count", func(t *testing.T) {
		store.Put("count-svc1-unique", &Operation{ID: 1, Name: "C1"})
		store.Put("count-svc2-unique", &Operation{ID: 2, Name: "C2"})

		assert.GreaterOrEqual(t, store.Count(), 2)
	})
}

func TestShapeStore(t *testing.T) {
	s := newTestStorage(t)
	defer s.Close()

	store := NewShapeStore(s)

	t.Run("Get returns nil for non-existent", func(t *testing.T) {
		shape, err := store.Get("nonexistent-unique")
		require.NoError(t, err)
		assert.Nil(t, shape)
	})

	t.Run("Put and Get", func(t *testing.T) {
		shape := &Shape{
			ID:      1,
			ShapeID: "TestShapeUnique",
			Type:    "structure",
		}
		err := store.Put(shape)
		require.NoError(t, err)

		got, err := store.Get("TestShapeUnique")
		require.NoError(t, err)
		assert.Equal(t, "TestShapeUnique", got.ShapeID)
		assert.Equal(t, "structure", got.Type)
	})

	t.Run("GetByID", func(t *testing.T) {
		shape := &Shape{
			ID:      42,
			ShapeID: "IdShapeUnique",
			Type:    "string",
		}
		err := store.Put(shape)
		require.NoError(t, err)

		got, err := store.GetByID(42)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "IdShapeUnique", got.ShapeID)
	})

	t.Run("Delete", func(t *testing.T) {
		shape := &Shape{
			ID:      3,
			ShapeID: "ToDeleteUnique",
		}
		err := store.Put(shape)
		require.NoError(t, err)

		err = store.Delete("ToDeleteUnique")
		require.NoError(t, err)

		got, err := store.Get("ToDeleteUnique")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("ForEach", func(t *testing.T) {
		store.Put(&Shape{ID: 1, ShapeID: "foreach1-unique", Type: "t1"})
		store.Put(&Shape{ID: 2, ShapeID: "foreach2-unique", Type: "t2"})

		ids := []string{}
		err := store.ForEach(func(s *Shape) error {
			ids = append(ids, s.ShapeID)
			return nil
		})
		require.NoError(t, err)
		assert.Contains(t, ids, "foreach1-unique")
		assert.Contains(t, ids, "foreach2-unique")
	})

	t.Run("Count", func(t *testing.T) {
		store.Put(&Shape{ID: 1, ShapeID: "count1-unique", Type: "t1"})
		store.Put(&Shape{ID: 2, ShapeID: "count2-unique", Type: "t2"})

		assert.GreaterOrEqual(t, store.Count(), 2)
	})
}

func TestMemberStore(t *testing.T) {
	s := newTestStorage(t)
	defer s.Close()

	store := NewMemberStore(s)

	t.Run("Get returns nil for non-existent", func(t *testing.T) {
		member, err := store.Get(1, "nonexistent-unique")
		require.NoError(t, err)
		assert.Nil(t, member)
	})

	t.Run("Put and Get", func(t *testing.T) {
		member := &Member{
			ID:            1,
			ShapeID:       100,
			Name:          "TestMemberUnique",
			TargetShapeID: 200,
			IsRequired:    true,
		}
		err := store.Put(member)
		require.NoError(t, err)

		got, err := store.Get(100, "TestMemberUnique")
		require.NoError(t, err)
		assert.Equal(t, "TestMemberUnique", got.Name)
		assert.True(t, got.IsRequired)
	})

	t.Run("Delete", func(t *testing.T) {
		member := &Member{
			ID:      2,
			ShapeID: 50,
			Name:    "ToDeleteUnique",
		}
		err := store.Put(member)
		require.NoError(t, err)

		err = store.Delete(50, "ToDeleteUnique")
		require.NoError(t, err)

		got, err := store.Get(50, "ToDeleteUnique")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("ForEachByShape", func(t *testing.T) {
		s := newTestStorageForSubtest(t)
		defer s.Close()
		localStore := NewMemberStore(s)

		localStore.Put(&Member{ID: 1, ShapeID: 7771, Name: "Z1"})
		localStore.Put(&Member{ID: 2, ShapeID: 7771, Name: "Z2"})
		localStore.Put(&Member{ID: 3, ShapeID: 7772, Name: "Z3"})

		count := 0
		err := localStore.ForEachByShape(7771, func(m *Member) error {
			count++
			return nil
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 2)
	})

	t.Run("Count", func(t *testing.T) {
		store.Put(&Member{ID: 1, ShapeID: 1, Name: "N1Unique"})
		store.Put(&Member{ID: 2, ShapeID: 2, Name: "N2Unique"})

		assert.GreaterOrEqual(t, store.Count(), 2)
	})
}

func TestServiceMode(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "IMPLEMENTED", ServiceModeImplemented.String())
		assert.Equal(t, "FALLBACK", ServiceModeFallback.String())
		assert.Equal(t, "ERROR_INJECTION", ServiceModeErrorInjection.String())
		assert.Equal(t, "UNKNOWN", ServiceMode(99).String())
	})

	t.Run("ParseServiceMode", func(t *testing.T) {
		assert.Equal(t, ServiceModeImplemented, ParseServiceMode("IMPLEMENTED"))
		assert.Equal(t, ServiceModeFallback, ParseServiceMode("FALLBACK"))
		assert.Equal(t, ServiceModeErrorInjection, ParseServiceMode("ERROR_INJECTION"))
		assert.Equal(t, ServiceModeFallback, ParseServiceMode("UNKNOWN"))
	})
}

func TestInt64ToString(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{123, "123"},
		{-456, "-456"},
		{9999999999, "9999999999"},
	}

	for _, tc := range tests {
		result := int64ToString(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
