package audit

import (
	"errors"
	"testing"

	"vorpalstacks/internal/services/aws/common/request"

	"github.com/stretchr/testify/assert"
)

func TestAuditConfig(t *testing.T) {
	t.Run("DefaultConfig returns config with Enabled false", func(t *testing.T) {
		cfg := DefaultConfig()
		assert.NotNil(t, cfg)
		assert.False(t, cfg.Enabled)
	})

	t.Run("LoadConfig reads environment variable", func(t *testing.T) {
		t.Setenv("VS_AUDIT_ENABLED", "true")
		cfg := LoadConfig()
		assert.True(t, cfg.Enabled)
	})

	t.Run("LoadConfig handles invalid env value", func(t *testing.T) {
		t.Setenv("VS_AUDIT_ENABLED", "invalid")
		cfg := LoadConfig()
		assert.False(t, cfg.Enabled)
	})
}

func TestNoOpRecorder(t *testing.T) {
	recorder := NewNoOpRecorder()

	t.Run("RecordEvent returns nil", func(t *testing.T) {
		event := &AuditEvent{
			EventName:   "CreateBucket",
			EventSource: "s3.amazonaws.com",
		}
		err := recorder.RecordEvent(event)
		assert.NoError(t, err)
	})

	t.Run("Record returns nil", func(t *testing.T) {
		err := recorder.Record(&AuditEvent{})
		assert.NoError(t, err)
	})
}

func TestEventBuilder(t *testing.T) {
	t.Run("GetEventSource returns correct source", func(t *testing.T) {
		tests := []struct {
			service  string
			expected string
		}{
			{"s3", "s3.amazonaws.com"},
			{"lambda", "lambda.amazonaws.com"},
			{"dynamodb", "dynamodb.amazonaws.com"},
			{"unknown-service", "unknown-service.amazonaws.com"},
		}

		for _, tt := range tests {
			assert.Equal(t, tt.expected, GetEventSource(tt.service))
		}
	})

	t.Run("Build creates event with correct fields", func(t *testing.T) {
		builder := NewEventBuilder("s3", "CreateBucket", nil, nil)
		response := map[string]interface{}{"Location": "test-bucket"}
		event := builder.Build(response, nil)

		assert.Equal(t, "CreateBucket", event.EventName)
		assert.Equal(t, "s3.amazonaws.com", event.EventSource)
		assert.NotNil(t, event.ResponseElements)
		assert.Empty(t, event.ErrorCode)
	})

	t.Run("Build with error populates error fields", func(t *testing.T) {
		builder := NewEventBuilder("lambda", "CreateFunction", nil, nil)
		err := errors.New("ResourceAlreadyExistsException: function already exists")
		event := builder.Build(nil, err)

		assert.Equal(t, "CreateFunction", event.EventName)
		assert.Equal(t, "ResourceAlreadyExists", event.ErrorCode)
		assert.Contains(t, event.ErrorMessage, "already exists")
	})

	t.Run("isReadOnlyOperation detects read operations", func(t *testing.T) {
		readOps := []string{"GetBucket", "ListObjects", "DescribeTrails", "LookupEvents", "Query"}
		writeOps := []string{"CreateBucket", "DeleteObject", "PutObject", "UpdateFunction"}

		for _, op := range readOps {
			builder := NewEventBuilder("", op, nil, nil)
			event := builder.Build(nil, nil)
			assert.True(t, event.ReadOnly, "Expected %s to be read-only", op)
		}

		for _, op := range writeOps {
			builder := NewEventBuilder("", op, nil, nil)
			event := builder.Build(nil, nil)
			assert.False(t, event.ReadOnly, "Expected %s to be write operation", op)
		}
	})

	t.Run("sanitizeParameters masks sensitive fields", func(t *testing.T) {
		params := map[string]interface{}{
			"Bucket":    "my-bucket",
			"Password":  "secret123",
			"APIKey":    "key123",
			"SecretKey": "sk123",
			"Token":     "tok123",
		}

		builder := NewEventBuilder("", "", nil, &request.ParsedRequest{Parameters: params})
		event := builder.Build(nil, nil)

		assert.Equal(t, "my-bucket", event.RequestParameters["Bucket"])
		assert.Equal(t, "****", event.RequestParameters["Password"])
		assert.Equal(t, "****", event.RequestParameters["APIKey"])
		assert.Equal(t, "****", event.RequestParameters["SecretKey"])
		assert.Equal(t, "****", event.RequestParameters["Token"])
	})
}
