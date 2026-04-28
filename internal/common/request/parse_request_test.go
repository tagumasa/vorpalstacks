package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/common/defaults"
)

func TestParseAWSRequest(t *testing.T) {
	t.Run("JSON body", func(t *testing.T) {
		body := `{"Key1": "value1", "Key2": 123}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)
		assert.Equal(t, "value1", parsed.Parameters["Key1"])
		assert.Equal(t, float64(123), parsed.Parameters["Key2"])
	})

	t.Run("Form urlencoded body", func(t *testing.T) {
		body := "Key1=value1&Key2=value2"
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)
		assert.Equal(t, "value1", parsed.Parameters["Key1"])
		assert.Equal(t, "value2", parsed.Parameters["Key2"])
	})

	t.Run("Query parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?Param1=val1&Param2=val2", nil)

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)
		assert.Equal(t, "val1", parsed.Parameters["Param1"])
		assert.Equal(t, "val2", parsed.Parameters["Param2"])
	})

	t.Run("X-Amz-Target header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("X-Amz-Target", "SomeService.SomeOperation")

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)
		assert.Equal(t, "SomeOperation", parsed.Operation)
	})

	t.Run("Query Action parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?Action=ListBuckets", nil)

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)
		assert.Equal(t, "ListBuckets", parsed.Operation)
	})

	t.Run("Empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)
		assert.Empty(t, parsed.Parameters)
		assert.Empty(t, parsed.Body)
	})
}

func TestGetStringParam(t *testing.T) {
	params := map[string]interface{}{
		"StringKey": "stringValue",
		"IntKey":    42,
		"BoolKey":   true,
	}

	t.Run("Existing string key", func(t *testing.T) {
		result := GetStringParam(params, "StringKey")
		assert.Equal(t, "stringValue", result)
	})

	t.Run("Missing key", func(t *testing.T) {
		result := GetStringParam(params, "NonExistent")
		assert.Equal(t, "", result)
	})

	t.Run("Non-string key returns empty", func(t *testing.T) {
		result := GetStringParam(params, "IntKey")
		assert.Equal(t, "", result)
	})
}

func TestGetIntParam(t *testing.T) {
	params := map[string]interface{}{
		"IntKey":    42,
		"Int64Key":  int64(100),
		"FloatKey":  float64(3.14),
		"StringKey": "123",
	}

	t.Run("Existing int", func(t *testing.T) {
		result := GetIntParam(params, "IntKey")
		assert.Equal(t, 42, result)
	})

	t.Run("Existing int64", func(t *testing.T) {
		result := GetIntParam(params, "Int64Key")
		assert.Equal(t, 100, result)
	})

	t.Run("Existing float64", func(t *testing.T) {
		result := GetIntParam(params, "FloatKey")
		assert.Equal(t, 3, result)
	})

	t.Run("String number", func(t *testing.T) {
		result := GetIntParam(params, "StringKey")
		assert.Equal(t, 123, result)
	})

	t.Run("Missing key", func(t *testing.T) {
		result := GetIntParam(params, "NonExistent")
		assert.Equal(t, 0, result)
	})
}

func TestGetBoolParam(t *testing.T) {
	params := map[string]interface{}{
		"BoolTrue":    true,
		"BoolFalse":   false,
		"StringTrue":  "true",
		"StringFalse": "false",
		"Other":       123,
	}

	t.Run("Existing bool true", func(t *testing.T) {
		result := GetBoolParam(params, "BoolTrue")
		assert.True(t, result)
	})

	t.Run("Existing bool false", func(t *testing.T) {
		result := GetBoolParam(params, "BoolFalse")
		assert.False(t, result)
	})

	t.Run("String true", func(t *testing.T) {
		result := GetBoolParam(params, "StringTrue")
		assert.True(t, result)
	})

	t.Run("String false", func(t *testing.T) {
		result := GetBoolParam(params, "StringFalse")
		assert.False(t, result)
	})

	t.Run("String TRUE uppercase", func(t *testing.T) {
		params["UpperTrue"] = "TRUE"
		result := GetBoolParam(params, "UpperTrue")
		assert.True(t, result)
	})

	t.Run("Missing key", func(t *testing.T) {
		result := GetBoolParam(params, "NonExistent")
		assert.False(t, result)
	})

	t.Run("Non-bool returns false", func(t *testing.T) {
		result := GetBoolParam(params, "Other")
		assert.False(t, result)
	})
}

func TestParsedRequest(t *testing.T) {
	t.Run("Full request", func(t *testing.T) {
		body := `{"BucketName": "test-bucket"}`
		req := httptest.NewRequest(http.MethodPost, "/?Action=CreateBucket", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Amz-Target", "S3.CreateBucket")

		parsed, err := ParseAWSRequest(req)
		require.NoError(t, err)

		assert.Equal(t, "CreateBucket", parsed.Operation)
		assert.Equal(t, "test-bucket", parsed.Parameters["BucketName"])
		assert.NotEmpty(t, parsed.Headers)
		assert.NotNil(t, parsed.QueryParams)
	})
}

func TestExtractRegionFromAuth(t *testing.T) {
	tests := []struct {
		name     string
		auth     string
		expected string
	}{
		{
			name:     "valid auth header",
			auth:     "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-west-2/s3/aws4_request, SignedHeaders=host;range;x-amz-date, Signature=...",
			expected: "us-west-2",
		},
		{
			name:     "empty auth header",
			auth:     "",
			expected: "",
		},
		{
			name:     "wrong algorithm",
			auth:     "AWS4-HMAC-SHA512 Credential=...",
			expected: "",
		},
		{
			name:     "missing credential",
			auth:     "AWS4-HMAC-SHA256 SignedHeaders=host, Signature=abc",
			expected: "",
		},
		{
			name:     "short credential",
			auth:     "AWS4-HMAC-SHA256 Credential=short, Signature=abc",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRegionFromAuth(tt.auth)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractAccountIDFromAuth(t *testing.T) {
	tests := []struct {
		name     string
		auth     string
		expected string
	}{
		{
			name:     "valid auth header",
			auth:     "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-west-2/s3/aws4_request, SignedHeaders=host, Signature=...",
			expected: "AKIAIOSFODNN7EXAMPLE",
		},
		{
			name:     "empty auth header",
			auth:     "",
			expected: "",
		},
		{
			name:     "wrong algorithm",
			auth:     "AWS4-HMAC-SHA512 Credential=...",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAccountIDFromAuth(tt.auth)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMapParam(t *testing.T) {
	params := map[string]interface{}{
		"MapKey": map[string]interface{}{"key": "value"},
		"StrKey": "string",
	}

	t.Run("existing map", func(t *testing.T) {
		result := GetMapParam(params, "MapKey")
		assert.NotNil(t, result)
		assert.Equal(t, "value", result["key"])
	})

	t.Run("missing key", func(t *testing.T) {
		result := GetMapParam(params, "NonExistent")
		assert.Nil(t, result)
	})

	t.Run("non-map key", func(t *testing.T) {
		result := GetMapParam(params, "StrKey")
		assert.Nil(t, result)
	})
}

func TestParseAttributes(t *testing.T) {
	t.Run("valid attributes", func(t *testing.T) {
		params := map[string]interface{}{
			"Attributes": map[string]interface{}{
				"DelaySeconds": "60",
			},
		}
		result := ParseAttributes(params, "Attributes")
		assert.Equal(t, "60", result["DelaySeconds"])
	})

	t.Run("missing attributes", func(t *testing.T) {
		params := map[string]interface{}{}
		result := ParseAttributes(params, "Attributes")
		assert.Empty(t, result)
	})
}

func TestGetStringList(t *testing.T) {
	t.Run("valid list", func(t *testing.T) {
		params := map[string]interface{}{
			"QueueArns": []interface{}{"arn1", "arn2"},
		}
		result := GetStringList(params, "QueueArns")
		assert.Equal(t, []string{"arn1", "arn2"}, result)
	})

	t.Run("missing list", func(t *testing.T) {
		params := map[string]interface{}{}
		result := GetStringList(params, "QueueArns")
		assert.Empty(t, result)
	})

	t.Run("invalid type", func(t *testing.T) {
		params := map[string]interface{}{
			"QueueArns": "invalid",
		}
		result := GetStringList(params, "QueueArns")
		assert.Empty(t, result)
	})
}

func TestGetArrayParam(t *testing.T) {
	t.Run("valid array", func(t *testing.T) {
		params := map[string]interface{}{
			"Items": []interface{}{1, 2, 3},
		}
		result := GetArrayParam(params, "Items")
		assert.Len(t, result, 3)
	})

	t.Run("lowercase fallback", func(t *testing.T) {
		params := map[string]interface{}{
			"items": []interface{}{"a", "b"},
		}
		result := GetArrayParam(params, "Items")
		assert.Len(t, result, 2)
	})

	t.Run("missing", func(t *testing.T) {
		params := map[string]interface{}{}
		result := GetArrayParam(params, "Items")
		assert.Nil(t, result)
	})
}

func TestParseQueryTags(t *testing.T) {
	t.Run("valid query tags", func(t *testing.T) {
		params := map[string]interface{}{
			"Tag.1.Key":   "Name",
			"Tag.1.Value": "Test",
			"Tag.2.Key":   "Env",
			"Tag.2.Value": "Prod",
		}
		result := ParseQueryTags(params)
		assert.Equal(t, "Test", result["Name"])
		assert.Equal(t, "Prod", result["Env"])
	})

	t.Run("empty params", func(t *testing.T) {
		params := map[string]interface{}{}
		result := ParseQueryTags(params)
		assert.Empty(t, result)
	})
}

func TestParseQueryAttributes(t *testing.T) {
	t.Run("valid query attributes", func(t *testing.T) {
		params := map[string]interface{}{
			"Attribute.1.Name":  "DelaySeconds",
			"Attribute.1.Value": "60",
		}
		result := ParseQueryAttributes(params, "Attribute")
		assert.Equal(t, "60", result["DelaySeconds"])
	})
}

func TestGetParamCaseInsensitive(t *testing.T) {
	params := map[string]interface{}{
		"QueueName": "my-queue",
	}

	t.Run("PascalCase match", func(t *testing.T) {
		result := GetParamCaseInsensitive(params, "QueueName")
		assert.Equal(t, "my-queue", result)
	})

	t.Run("lowercase fallback", func(t *testing.T) {
		params := map[string]interface{}{
			"queuename": "my-queue",
		}
		result := GetParamCaseInsensitive(params, "QueueName")
		assert.Equal(t, "my-queue", result)
	})

	t.Run("no match", func(t *testing.T) {
		result := GetParamCaseInsensitive(params, "NonExistent")
		assert.Empty(t, result)
	})
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"QueueName", "queueName"},
		{"ABC", "aBC"},
		{"a", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := LowerFirst(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParsedRequest_GetRegion(t *testing.T) {
	t.Run("explicit region", func(t *testing.T) {
		req := &ParsedRequest{Region: "us-west-2"}
		assert.Equal(t, "us-west-2", req.GetRegion())
	})

	t.Run("default region", func(t *testing.T) {
		req := &ParsedRequest{}
		assert.Equal(t, defaults.DefaultRegion, req.GetRegion())
	})
}

func TestGetMapParamCaseInsensitive(t *testing.T) {
	params := map[string]interface{}{
		"Attributes": map[string]interface{}{"key": "value"},
	}

	t.Run("PascalCase match", func(t *testing.T) {
		result := GetMapParamCaseInsensitive(params, "Attributes")
		assert.NotNil(t, result)
	})

	t.Run("lowercase fallback", func(t *testing.T) {
		params := map[string]interface{}{
			"attributes": map[string]interface{}{"key": "value"},
		}
		result := GetMapParamCaseInsensitive(params, "Attributes")
		assert.NotNil(t, result)
	})

	t.Run("no match", func(t *testing.T) {
		result := GetMapParamCaseInsensitive(params, "NonExistent")
		assert.Nil(t, result)
	})
}

func TestGetArrayParamLowerFirst(t *testing.T) {
	params := map[string]interface{}{
		"Items": []interface{}{1, 2, 3},
	}

	t.Run("PascalCase match", func(t *testing.T) {
		result := GetArrayParamLowerFirst(params, "Items")
		assert.Len(t, result, 3)
	})

	t.Run("camelCase fallback", func(t *testing.T) {
		params := map[string]interface{}{
			"items": []interface{}{"a", "b"},
		}
		result := GetArrayParamLowerFirst(params, "Items")
		assert.Len(t, result, 2)
	})
}
