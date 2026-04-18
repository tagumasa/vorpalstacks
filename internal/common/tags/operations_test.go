package tags

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/utils/aws/types"
)

func TestParseTagsFromMap(t *testing.T) {
	t.Run("array format with Key/Value", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags": []interface{}{
				map[string]interface{}{
					"Key":   "k1",
					"Value": "v1",
				},
			},
		}
		result := ParseTags(params, "Tags")
		require.Len(t, result, 1)
		assert.Equal(t, "k1", result[0].Key)
		assert.Equal(t, "v1", result[0].Value)
	})

	t.Run("multiple tags", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags": []interface{}{
				map[string]interface{}{"Key": "k2", "Value": "v2"},
				map[string]interface{}{"Key": "k3", "Value": "v3"},
			},
		}
		result := ParseTags(params, "Tags")
		require.Len(t, result, 2)
		assert.Equal(t, "k2", result[0].Key)
		assert.Equal(t, "k3", result[1].Key)
	})

	t.Run("fallback to Items with Tag wrapper", func(t *testing.T) {
		params := map[string]interface{}{
			"Items": map[string]interface{}{
				"Items": []interface{}{
					map[string]interface{}{
						"Tag": map[string]interface{}{
							"Key":   "k3",
							"Value": "v3",
						},
					},
				},
			},
		}
		result := ParseTags(params, "Tags")
		require.Len(t, result, 1)
		assert.Equal(t, "k3", result[0].Key)
	})

	t.Run("map format", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags": map[string]interface{}{
				"Env": "Prod",
			},
		}
		result := ParseTags(params, "Tags")
		require.Len(t, result, 1)
		assert.Equal(t, "Env", result[0].Key)
	})
}

func TestParseTagsWithQueryFallback(t *testing.T) {
	t.Run("JSON format", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags": []interface{}{
				map[string]interface{}{"Key": "k1", "Value": "v1"},
			},
		}
		result := ParseTagsWithQueryFallback(params, "Tags")
		require.Len(t, result, 1)
	})

	t.Run("query member format", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags.member.1.Key":   "env",
			"Tags.member.1.Value": "prod",
			"Tags.member.2.Key":   "team",
			"Tags.member.2.Value": "backend",
		}
		result := ParseTagsWithQueryFallback(params, "Tags")
		require.Len(t, result, 2)
		assert.Equal(t, "env", result[0].Key)
		assert.Equal(t, "team", result[1].Key)
	})
}

func TestParseTagKeysWithQueryFallback(t *testing.T) {
	t.Run("slice format", func(t *testing.T) {
		params := map[string]interface{}{
			"TagKeys": []interface{}{"k1", "k2"},
		}
		result := ParseTagKeysWithQueryFallback(params, "TagKeys")
		assert.Equal(t, []string{"k1", "k2"}, result)
	})

	t.Run("query member format", func(t *testing.T) {
		params := map[string]interface{}{
			"TagKeys.member.1": "k1",
			"TagKeys.member.2": "k2",
		}
		result := ParseTagKeysWithQueryFallback(params, "TagKeys")
		assert.Equal(t, []string{"k1", "k2"}, result)
	})

	t.Run("single key as string", func(t *testing.T) {
		params := map[string]interface{}{
			"TagKeys": "single",
		}
		result := ParseTagKeysWithQueryFallback(params, "TagKeys")
		assert.Equal(t, []string{"single"}, result)
	})
}

func TestParseTagKeysAsSlice(t *testing.T) {
	t.Run("[]interface{}", func(t *testing.T) {
		params := map[string]interface{}{
			"Keys": []interface{}{"a", "b"},
		}
		assert.Equal(t, []string{"a", "b"}, ParseTagKeysAsSlice(params, "Keys"))
	})

	t.Run("[]string", func(t *testing.T) {
		params := map[string]interface{}{
			"Keys": []string{"a", "b"},
		}
		assert.Equal(t, []string{"a", "b"}, ParseTagKeysAsSlice(params, "Keys"))
	})

	t.Run("string", func(t *testing.T) {
		params := map[string]interface{}{
			"Keys": "single",
		}
		assert.Equal(t, []string{"single"}, ParseTagKeysAsSlice(params, "Keys"))
	})

	t.Run("missing key", func(t *testing.T) {
		params := map[string]interface{}{}
		assert.Nil(t, ParseTagKeysAsSlice(params, "Keys"))
	})
}

func TestToResponseWithKeyNames(t *testing.T) {
	tags := []types.Tag{{Key: "k", Value: "v"}}
	result := ToResponseWithKeyNames(tags, "TagKey", "TagValue")
	require.Len(t, result, 1)
	assert.Equal(t, "k", result[0]["TagKey"])
	assert.Equal(t, "v", result[0]["TagValue"])
}

func TestToMap(t *testing.T) {
	tags := []types.Tag{{Key: "a", Value: "1"}, {Key: "b", Value: "2"}}
	m := ToMap(tags)
	assert.Equal(t, "1", m["a"])
	assert.Equal(t, "2", m["b"])
}

func TestMapToTags(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		result := MapToTags(map[string]string{"k": "v"})
		require.Len(t, result, 1)
		assert.Equal(t, "k", result[0].Key)
	})

	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, MapToTags(nil))
	})
}

func TestMapInterfaceToTags(t *testing.T) {
	t.Run("string values", func(t *testing.T) {
		result := MapInterfaceToTags(map[string]interface{}{"k": "v"})
		require.Len(t, result, 1)
	})

	t.Run("non-string values filtered", func(t *testing.T) {
		result := MapInterfaceToTags(map[string]interface{}{"k": 123})
		assert.Len(t, result, 0)
	})

	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, MapInterfaceToTags(nil))
	})
}

func TestMapToResponse(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		result := MapToResponse(map[string]string{"k": "v"})
		require.Len(t, result, 1)
		assert.Equal(t, "k", result[0]["Key"])
	})

	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, MapToResponse(map[string]string{}))
	})
}

func TestParseTagsAsMap(t *testing.T) {
	t.Run("map format", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags": map[string]interface{}{"Env": "Prod"},
		}
		result := ParseTagsAsMap(params, "Tags")
		assert.Equal(t, "Prod", result["Env"])
	})

	t.Run("missing", func(t *testing.T) {
		result := ParseTagsAsMap(map[string]interface{}{}, "Tags")
		assert.Empty(t, result)
	})
}

func TestConvertTags(t *testing.T) {
	type input struct {
		Name  string
		Value string
	}
	items := []input{{Name: "n1", Value: "v1"}}
	result := ConvertTags(items, func(i input) types.Tag {
		return types.Tag{Key: i.Name, Value: i.Value}
	})
	require.Len(t, result, 1)
	assert.Equal(t, "n1", result[0].Key)

	assert.Nil(t, ConvertTags(nil, func(i input) types.Tag { return types.Tag{} }))
}

func TestParseTagsWithKeyNames(t *testing.T) {
	params := map[string]interface{}{
		"Tags": []interface{}{
			map[string]interface{}{"TagKey": "k", "TagValue": "v"},
		},
	}
	result := ParseTagsWithKeyNames(params, "Tags", "TagKey", "TagValue")
	require.Len(t, result, 1)
	assert.Equal(t, "k", result[0].Key)
}

func TestParseTagsWithQueryFallbackAndKeyNames(t *testing.T) {
	params := map[string]interface{}{
		"Tags.member.1.TagKey":   "k",
		"Tags.member.1.TagValue": "v",
	}
	result := ParseTagsWithQueryFallbackAndKeyNames(params, "Tags", "TagKey", "TagValue")
	require.Len(t, result, 1)
	assert.Equal(t, "k", result[0].Key)
}

func TestParseTagsWithPrefix(t *testing.T) {
	params := map[string]interface{}{
		"Tag.1.Key":   "k1",
		"Tag.1.Value": "v1",
		"Tag.2.Key":   "k2",
		"Tag.2.Value": "v2",
	}
	result := ParseTagsWithPrefix(params, "Tag")
	require.Len(t, result, 2)
	assert.Equal(t, "k1", result[0].Key)
	assert.Equal(t, "k2", result[1].Key)
}

func TestParseTagKeysWithKeyName(t *testing.T) {
	t.Run("string items", func(t *testing.T) {
		params := map[string]interface{}{
			"Keys": []interface{}{"a", "b"},
		}
		result := ParseTagKeysWithKeyName(params, "Keys", "KeyName")
		assert.Equal(t, []string{"a", "b"}, result)
	})

	t.Run("map items", func(t *testing.T) {
		params := map[string]interface{}{
			"Keys": []interface{}{
				map[string]interface{}{"TagName": "a"},
			},
		}
		result := ParseTagKeysWithKeyName(params, "Keys", "TagName")
		assert.Equal(t, []string{"a"}, result)
	})
}

func TestConvertToMapSlice(t *testing.T) {
	tags := []types.Tag{{Key: "k", Value: "v"}}
	result := ConvertToMapSlice(tags)
	require.Len(t, result, 1)
	assert.Equal(t, "k", result[0]["Key"])
	assert.Equal(t, "v", result[0]["Value"])

	assert.Nil(t, ConvertToMapSlice(nil))
}

func TestParseMessageTags(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		params := map[string]interface{}{
			"MessageTags": []interface{}{
				map[string]interface{}{"Name": "n1", "Value": "v1"},
			},
		}
		result := ParseMessageTags(params, "MessageTags")
		require.Len(t, result, 1)
		assert.Equal(t, "n1", result[0].Name)
	})

	t.Run("missing", func(t *testing.T) {
		assert.Nil(t, ParseMessageTags(map[string]interface{}{}, "MessageTags"))
	})

	t.Run("wrong type", func(t *testing.T) {
		params := map[string]interface{}{"MessageTags": "not a list"}
		assert.Nil(t, ParseMessageTags(params, "MessageTags"))
	})
}

func TestParseMessageTagsFromList(t *testing.T) {
	t.Run("mixed types", func(t *testing.T) {
		list := []interface{}{
			map[string]interface{}{"Name": "n1", "Value": "v1"},
			"invalid",
		}
		result := ParseMessageTagsFromList(list)
		require.Len(t, result, 1)
	})
}

func TestParseEcsTags(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		data := []interface{}{
			map[string]interface{}{"key": "k", "value": "v"},
		}
		result := ParseEcsTags(data)
		require.Len(t, result, 1)
		assert.Equal(t, "v", result[0]["value"])
	})

	t.Run("nil input", func(t *testing.T) {
		assert.Nil(t, ParseEcsTags(nil))
	})
}

func TestConvertMessageTagsToSESv2(t *testing.T) {
	tags := []MessageTag{{Name: "n1", Value: "v1"}}
	result := ConvertMessageTagsToSESv2(tags)
	require.Len(t, result, 1)
	assert.Equal(t, "n1", result[0].Name)
	assert.Equal(t, "v1", result[0].Value)

	assert.Nil(t, ConvertMessageTagsToSESv2(nil))
}

func TestGetResourceKey(t *testing.T) {
	t.Run("case insensitive", func(t *testing.T) {
		params := map[string]interface{}{
			"QueueUrl": "http://localhost:8080/123/my-queue",
		}
		key := GetResourceKey(params, SQSConfig)
		assert.Equal(t, "http://localhost:8080/123/my-queue", key)
	})

	t.Run("alt param", func(t *testing.T) {
		params := map[string]interface{}{
			"ResourceIdList": []interface{}{"trail-1"},
		}
		key := GetResourceKey(params, CloudTrailConfig)
		assert.Equal(t, "trail-1", key)
	})

	t.Run("missing required", func(t *testing.T) {
		key := GetResourceKey(map[string]interface{}{}, StandardConfig)
		assert.Equal(t, "", key)
	})
}

func TestHandleListSimple(t *testing.T) {
	f := func(key string) ([]types.Tag, error) {
		return []types.Tag{{Key: "k", Value: "v"}}, nil
	}
	result, err := HandleListSimple(context.Background(), StandardConfig, "res", f)
	require.NoError(t, err)
	m := result.(map[string]interface{})
	tags := m["Tags"].([]map[string]interface{})
	assert.Len(t, tags, 1)
}
