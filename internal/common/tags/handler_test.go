package tags

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

func TestHandleTag(t *testing.T) {
	t.Run("happy path with TagFunc", func(t *testing.T) {
		var taggedKey string
		var taggedTags []types.Tag
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			TagFunc: func(_ context.Context, resourceKey string, tags []types.Tag) error {
				taggedKey = resourceKey
				taggedTags = tags
				return nil
			},
			EmptyResponse: func() (interface{}, error) { return nil, nil },
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:aws:sns:us-east-1:123:topic/test",
				"Tags": []interface{}{
					map[string]interface{}{"Key": "env", "Value": "prod"},
				},
			},
		}
		resp, err := HandleTag(context.Background(), req, cfg)
		require.NoError(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "arn:aws:sns:us-east-1:123:topic/test", taggedKey)
		assert.Equal(t, []types.Tag{{Key: "env", Value: "prod"}}, taggedTags)
	})

	t.Run("missing resource returns error", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		req := &request.ParsedRequest{Parameters: map[string]interface{}{}}
		_, err := HandleTag(context.Background(), req, cfg)
		require.Error(t, err)
		assert.IsType(t, &MissingResourceError{}, err)
	})

	t.Run("missing tags when RequireTags returns error", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:aws:sns:us-east-1:123:topic/test",
			},
		}
		_, err := HandleTag(context.Background(), req, cfg)
		require.Error(t, err)
		assert.IsType(t, &MissingTagsError{}, err)
	})

	t.Run("no error when tags empty and RequireTags false", func(t *testing.T) {
		custom := StandardConfig
		custom.RequireTags = false
		cfg := TagHandlerConfig{Param: custom}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:aws:sns:us-east-1:123:topic/test",
			},
		}
		_, err := HandleTag(context.Background(), req, cfg)
		require.NoError(t, err)
	})

	t.Run("ResourceKey transform", func(t *testing.T) {
		var receivedKey string
		cfg := TagHandlerConfig{
			Param:       StandardConfig,
			ResourceKey: func(rawKey string) string { return rawKey + "-transformed" },
			TagFunc: func(_ context.Context, resourceKey string, _ []types.Tag) error {
				receivedKey = resourceKey
				return nil
			},
			EmptyResponse: func() (interface{}, error) { return nil, nil },
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"Tags":        []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}},
			},
		}
		_, err := HandleTag(context.Background(), req, cfg)
		require.NoError(t, err)
		assert.Equal(t, "arn:test-transformed", receivedKey)
	})

	t.Run("ValidateResource rejection", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			ValidateResource: func(_ context.Context, _ string) error {
				return &MissingResourceError{Param: "ResourceArn"}
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"Tags":        []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}},
			},
		}
		_, err := HandleTag(context.Background(), req, cfg)
		require.Error(t, err)
	})

	t.Run("TagFunc error mapped by MapError", func(t *testing.T) {
		sentinelErr := errors.New("tag failed")
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			TagFunc: func(_ context.Context, _ string, _ []types.Tag) error {
				return sentinelErr
			},
			MapError: func(err error) error {
				if errors.Is(err, sentinelErr) {
					return errors.New("mapped error")
				}
				return err
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"Tags":        []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}},
			},
		}
		_, err := HandleTag(context.Background(), req, cfg)
		require.Error(t, err)
		assert.Equal(t, "mapped error", err.Error())
	})

	t.Run("TagResponse callback", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param:   StandardConfig,
			TagFunc: func(_ context.Context, _ string, _ []types.Tag) error { return nil },
			TagResponse: func(_ context.Context, _ string) (interface{}, error) {
				return "tag-response", nil
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"Tags":        []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}},
			},
		}
		resp, err := HandleTag(context.Background(), req, cfg)
		require.NoError(t, err)
		assert.Equal(t, "tag-response", resp)
	})

	t.Run("EmptyResponse callback", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param:   StandardConfig,
			TagFunc: func(_ context.Context, _ string, _ []types.Tag) error { return nil },
			EmptyResponse: func() (interface{}, error) {
				return "empty", nil
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"Tags":        []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}},
			},
		}
		resp, err := HandleTag(context.Background(), req, cfg)
		require.NoError(t, err)
		assert.Equal(t, "empty", resp)
	})

	t.Run("nil TagFunc with tags is no-op", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"Tags":        []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}},
			},
		}
		resp, err := HandleTag(context.Background(), req, cfg)
		require.NoError(t, err)
		assert.Nil(t, resp)
	})
}

func TestHandleUntag(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var receivedKey string
		var receivedTagKeys []string
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
				receivedKey = resourceKey
				receivedTagKeys = tagKeys
				return nil
			},
			EmptyResponse: func() (interface{}, error) { return nil, nil },
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"TagKeys":     []interface{}{"env"},
			},
		}
		_, err := HandleUntag(context.Background(), req, cfg)
		require.NoError(t, err)
		assert.Equal(t, "arn:test", receivedKey)
		assert.Equal(t, []string{"env"}, receivedTagKeys)
	})

	t.Run("missing resource", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		_, err := HandleUntag(context.Background(), &request.ParsedRequest{Parameters: map[string]interface{}{}}, cfg)
		require.Error(t, err)
		assert.IsType(t, &MissingResourceError{}, err)
	})

	t.Run("missing tag keys when RequireTagKeys returns error", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{"ResourceArn": "arn:test"},
		}
		_, err := HandleUntag(context.Background(), req, cfg)
		require.Error(t, err)
		assert.IsType(t, &MissingTagKeysError{}, err)
	})

	t.Run("no error when tag keys empty and RequireTagKeys false", func(t *testing.T) {
		custom := StandardConfig
		custom.RequireTagKeys = false
		cfg := TagHandlerConfig{Param: custom}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{"ResourceArn": "arn:test"},
		}
		_, err := HandleUntag(context.Background(), req, cfg)
		require.NoError(t, err)
	})

	t.Run("UntagFunc error mapped by MapError", func(t *testing.T) {
		sentinelErr := errors.New("untag failed")
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			UntagFunc: func(_ context.Context, _ string, _ []string) error {
				return sentinelErr
			},
			MapError: func(err error) error {
				if errors.Is(err, sentinelErr) {
					return errors.New("mapped untag error")
				}
				return err
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"ResourceArn": "arn:test",
				"TagKeys":     []interface{}{"k"},
			},
		}
		_, err := HandleUntag(context.Background(), req, cfg)
		require.Error(t, err)
		assert.Equal(t, "mapped untag error", err.Error())
	})
}

func TestHandleList(t *testing.T) {
	t.Run("happy path with default FormatResponse", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
				return []types.Tag{{Key: "env", Value: "prod"}}, nil
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{"ResourceArn": "arn:test"},
		}
		resp, err := HandleList(context.Background(), req, cfg)
		require.NoError(t, err)
		m, ok := resp.(map[string]interface{})
		require.True(t, ok)
		tags, ok := m["Tags"].([]map[string]interface{})
		require.True(t, ok)
		assert.Len(t, tags, 1)
		assert.Equal(t, "env", tags[0]["Key"])
		assert.Equal(t, "prod", tags[0]["Value"])
	})

	t.Run("missing resource", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		_, err := HandleList(context.Background(), &request.ParsedRequest{Parameters: map[string]interface{}{}}, cfg)
		require.Error(t, err)
		assert.IsType(t, &MissingResourceError{}, err)
	})

	t.Run("nil ListFunc returns error", func(t *testing.T) {
		cfg := TagHandlerConfig{Param: StandardConfig}
		_, err := HandleList(context.Background(), &request.ParsedRequest{Parameters: map[string]interface{}{"ResourceArn": "arn:test"}}, cfg)
		require.Error(t, err)
		assert.IsType(t, &MissingResourceError{}, err)
	})

	t.Run("custom FormatResponse", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			ListFunc: func(_ context.Context, _ string) ([]types.Tag, error) {
				return nil, nil
			},
			FormatResponse: func(tags []types.Tag, rawKey string) (interface{}, error) {
				return map[string]interface{}{"CustomTags": tags, "Resource": rawKey}, nil
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{"ResourceArn": "arn:test"},
		}
		resp, err := HandleList(context.Background(), req, cfg)
		require.NoError(t, err)
		m, ok := resp.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "arn:test", m["Resource"])
	})

	t.Run("default FormatResponse uses custom key names", func(t *testing.T) {
		kmsCfg := KMSConfig
		cfg := TagHandlerConfig{
			Param: kmsCfg,
			ListFunc: func(_ context.Context, _ string) ([]types.Tag, error) {
				return []types.Tag{{Key: "env", Value: "prod"}}, nil
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{"KeyId": "k1"},
		}
		resp, err := HandleList(context.Background(), req, cfg)
		require.NoError(t, err)
		m, ok := resp.(map[string]interface{})
		require.True(t, ok)
		tags, ok := m["Tags"].([]map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "env", tags[0]["TagKey"])
		assert.Equal(t, "prod", tags[0]["TagValue"])
	})

	t.Run("empty tags returns empty array", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			ListFunc: func(_ context.Context, _ string) ([]types.Tag, error) {
				return nil, nil
			},
		}
		req := &request.ParsedRequest{
			Parameters: map[string]interface{}{"ResourceArn": "arn:test"},
		}
		resp, err := HandleList(context.Background(), req, cfg)
		require.NoError(t, err)
		m, ok := resp.(map[string]interface{})
		require.True(t, ok)
		tags, ok := m["Tags"].([]map[string]interface{})
		require.True(t, ok)
		assert.Empty(t, tags)
	})

	t.Run("ValidateResource rejection", func(t *testing.T) {
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			ValidateResource: func(_ context.Context, _ string) error {
				return &MissingResourceError{Param: "ResourceArn"}
			},
		}
		_, err := HandleList(context.Background(), &request.ParsedRequest{Parameters: map[string]interface{}{"ResourceArn": "arn:test"}}, cfg)
		require.Error(t, err)
	})

	t.Run("MapError applied to ListFunc error", func(t *testing.T) {
		sentinelErr := errors.New("list failed")
		cfg := TagHandlerConfig{
			Param: StandardConfig,
			ListFunc: func(_ context.Context, _ string) ([]types.Tag, error) {
				return nil, sentinelErr
			},
			MapError: func(err error) error {
				if errors.Is(err, sentinelErr) {
					return errors.New("mapped list error")
				}
				return err
			},
		}
		_, err := HandleList(context.Background(), &request.ParsedRequest{Parameters: map[string]interface{}{"ResourceArn": "arn:test"}}, cfg)
		require.Error(t, err)
		assert.Equal(t, "mapped list error", err.Error())
	})
}

func TestSliceTagsToKeys(t *testing.T) {
	t.Run("multiple tags", func(t *testing.T) {
		tags := []types.Tag{{Key: "a", Value: "1"}, {Key: "b", Value: "2"}}
		keys := SliceTagsToKeys(tags)
		assert.Equal(t, []string{"a", "b"}, keys)
	})

	t.Run("empty tags", func(t *testing.T) {
		keys := SliceTagsToKeys(nil)
		assert.NotNil(t, keys)
		assert.Empty(t, keys)
	})
}

func TestMissingResourceError(t *testing.T) {
	err := &MissingResourceError{Param: "ResourceId"}
	assert.Equal(t, "ResourceId is required", err.Error())
}

func TestMissingTagsError(t *testing.T) {
	err := &MissingTagsError{Param: "Tags"}
	assert.Equal(t, "Tags is required", err.Error())
}

func TestMissingTagKeysError(t *testing.T) {
	err := &MissingTagKeysError{Param: "TagKeys"}
	assert.Equal(t, "TagKeys is required", err.Error())
}
