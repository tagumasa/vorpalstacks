package audit

import (
	"errors"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditConfig(t *testing.T) {
	t.Run("DefaultConfig returns config with Enabled false", func(t *testing.T) {
		cfg := DefaultConfig()
		assert.NotNil(t, cfg)
		assert.False(t, cfg.Enabled)
	})

	t.Run("LoadConfig reads environment variable", func(t *testing.T) {
		t.Setenv("CLOUDTRAIL_ENABLED", "true")
		cfg := LoadConfig()
		assert.True(t, cfg.Enabled)
	})

	t.Run("LoadConfig handles invalid env value", func(t *testing.T) {
		t.Setenv("CLOUDTRAIL_ENABLED", "invalid")
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

	t.Run("Build with a wire error records its declared code", func(t *testing.T) {
		builder := NewEventBuilder("lambda", "CreateFunction", nil, nil)
		err := awserrors.NewAWSError("ResourceAlreadyExistsException", "function already exists", 409)
		event := builder.Build(nil, err)

		assert.Equal(t, "CreateFunction", event.EventName)
		assert.Equal(t, "ResourceAlreadyExistsException", event.ErrorCode)
		assert.Contains(t, event.ErrorMessage, "already exists")
	})

	t.Run("Build with a non-wire error records the generic internal code", func(t *testing.T) {
		builder := NewEventBuilder("lambda", "CreateFunction", nil, nil)
		event := builder.Build(nil, errors.New("connection reset by peer"))

		assert.Equal(t, "InternalError", event.ErrorCode)
		assert.Contains(t, event.ErrorMessage, "connection reset")
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

func TestBuildUserIdentity(t *testing.T) {
	newCtx := func(accountID string, mutate func(*request.RequestContext)) *request.RequestContext {
		ctx := request.NewRequestContext(nil, nil, accountID, "us-east-1")
		if mutate != nil {
			mutate(ctx)
		}
		return ctx
	}

	t.Run("root principal records as Root", func(t *testing.T) {
		builder := NewEventBuilder("s3", "PutObject", newCtx("123456789012", func(c *request.RequestContext) {
			c.Principal = "123456789012"
			c.PrincipalID = "123456789012"
			c.PrincipalType = request.PrincipalTypeRoot
		}), nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "Root", id.Type)
		assert.Equal(t, "123456789012", id.PrincipalID)
		assert.Equal(t, "arn:aws:iam::123456789012:root", id.ARN)
		assert.Equal(t, "123456789012", id.AccountID)
	})

	t.Run("resolved IAM user records as IAMUser with its real ID", func(t *testing.T) {
		builder := NewEventBuilder("lambda", "CreateFunction", newCtx("123456789012", func(c *request.RequestContext) {
			c.Principal = "Alice"
			c.PrincipalID = "AIDAEXAMPLE"
			c.PrincipalType = request.PrincipalTypeUser
		}), nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "IAMUser", id.Type)
		assert.Equal(t, "AIDAEXAMPLE", id.PrincipalID)
		assert.Equal(t, "arn:aws:iam::123456789012:user/Alice", id.ARN)
		assert.Equal(t, "Alice", id.UserName)
	})

	t.Run("role session records as AssumedRole with session context", func(t *testing.T) {
		created := time.Date(2026, 9, 17, 1, 2, 3, 456789000, time.UTC)
		builder := NewEventBuilder("s3", "PutObject", newCtx("123456789012", func(c *request.RequestContext) {
			c.Principal = "arn:aws:sts::123456789012:assumed-role/DevRole/Dev1"
			c.PrincipalID = "arn:aws:sts::123456789012:assumed-role/DevRole/Dev1"
			c.PrincipalType = request.PrincipalTypeRole
			c.Session = &request.SessionInfo{
				CredentialPrincipalType: "AssumedRole",
				SessionName:             "Dev1",
				IssuerType:              "Role",
				IssuerPrincipalID:       "AROAEXAMPLE",
				IssuerARN:               "arn:aws:iam::123456789012:role/DevRole",
				IssuerAccountID:         "123456789012",
				IssuerUserName:          "DevRole",
				MFAAuthenticated:        true,
				CreationDate:            created,
			}
		}), nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "AssumedRole", id.Type)
		assert.Equal(t, "AROAEXAMPLE:Dev1", id.PrincipalID)
		assert.Equal(t, "arn:aws:sts::123456789012:assumed-role/DevRole/Dev1", id.ARN)
		assert.Equal(t, "Dev1", id.UserName)
		require.NotNil(t, id.SessionContext)
		assert.Equal(t, "Role", id.SessionContext.SessionIssuer.Type)
		assert.Equal(t, "arn:aws:iam::123456789012:role/DevRole", id.SessionContext.SessionIssuer.ARN)
		assert.Equal(t, "DevRole", id.SessionContext.SessionIssuer.UserName)
		assert.Equal(t, "true", id.SessionContext.Attributes.MFAAuthenticated)
		assert.Equal(t, "2026-09-17T01:02:03Z", id.SessionContext.Attributes.CreationDate.Format(time.RFC3339))
	})

	t.Run("federated session records as FederatedUser", func(t *testing.T) {
		builder := NewEventBuilder("s3", "GetBucket", newCtx("123456789012", func(c *request.RequestContext) {
			c.Principal = "Bob"
			c.PrincipalID = "AIDABOB"
			c.PrincipalType = request.PrincipalTypeUser
			c.Session = &request.SessionInfo{
				CredentialPrincipalType: "FederatedUser",
				SessionName:             "Bob",
				IssuerType:              "IAMUser",
				IssuerPrincipalID:       "AIDABOB",
				IssuerARN:               "arn:aws:iam::123456789012:user/Bob",
				IssuerUserName:          "Bob",
			}
		}), nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "FederatedUser", id.Type)
		assert.Equal(t, "AIDABOB:Bob", id.PrincipalID)
		assert.Equal(t, "arn:aws:sts::123456789012:federated-user/Bob", id.ARN)
		require.NotNil(t, id.SessionContext)
		assert.Equal(t, "IAMUser", id.SessionContext.SessionIssuer.Type)
	})

	t.Run("unresolvable caller records as Unknown", func(t *testing.T) {
		builder := NewEventBuilder("cloudtrail", "ListTrails", newCtx("123456789012", func(c *request.RequestContext) {
			c.Principal = "test"
		}), nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "Unknown", id.Type)
		assert.Equal(t, "123456789012", id.AccountID)
		assert.Empty(t, id.PrincipalID)
	})

	t.Run("root user name without type records as Root", func(t *testing.T) {
		builder := NewEventBuilder("iam", "GetUser", newCtx("123456789012", func(c *request.RequestContext) {
			c.Principal = iam.RootUserName
		}), nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "Root", id.Type)
	})

	t.Run("nil request context records as Unknown", func(t *testing.T) {
		builder := NewEventBuilder("s3", "PutObject", nil, nil)
		id := builder.Build(nil, nil).UserIdentity
		assert.Equal(t, "Unknown", id.Type)
	})
}

func TestExtractResources(t *testing.T) {
	newCtx := func() *request.RequestContext {
		return request.NewRequestContext(nil, nil, "123456789012", "us-east-1")
	}

	t.Run("ARNs across services classify by resource word", func(t *testing.T) {
		builder := NewEventBuilder("cloudtrail", "PutResourcePolicy", newCtx(), nil)
		response := map[string]interface{}{
			"ResourceArn": "arn:aws:cloudtrail:us-east-1:123456789012:trail/example-trail",
			"Nested": map[string]interface{}{
				"Key": "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			},
			"Ignored": "not-an-arn",
		}
		entries := builder.Build(response, nil).Resources
		require.Len(t, entries, 2)
		assert.Equal(t, ResourceEntry{
			ResourceType: "AWS::CloudTrail::Trail",
			ResourceName: "arn:aws:cloudtrail:us-east-1:123456789012:trail/example-trail",
		}, entries[0])
		assert.Equal(t, ResourceEntry{
			ResourceType: "AWS::KMS::Key",
			ResourceName: "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
		}, entries[1])
	})

	t.Run("request-parameter ARNs are swept alongside the response", func(t *testing.T) {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{
			"ResourceArn": "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/example",
		}}
		builder := NewEventBuilder("logs", "DeleteLogGroup", newCtx(), req)
		entries := builder.Build(nil, nil).Resources
		require.Len(t, entries, 1)
		assert.Equal(t, "AWS::Logs::LogGroup", entries[0].ResourceType)
	})

	t.Run("unmapped ARN records without a type identifier", func(t *testing.T) {
		builder := NewEventBuilder("firehose", "PutRecord", newCtx(), nil)
		response := map[string]interface{}{
			"Arn": "arn:aws:firehose:us-east-1:123456789012:deliverystream/example",
		}
		entries := builder.Build(response, nil).Resources
		require.Len(t, entries, 1)
		assert.Empty(t, entries[0].ResourceType)
		assert.Equal(t, "arn:aws:firehose:us-east-1:123456789012:deliverystream/example", entries[0].ResourceName)
	})

	t.Run("S3 bucket location names the bucket resource", func(t *testing.T) {
		builder := NewEventBuilder("s3", "CreateBucket", newCtx(), nil)
		entries := builder.Build(map[string]interface{}{"Location": "example-bucket"}, nil).Resources
		require.Len(t, entries, 1)
		assert.Equal(t, ResourceEntry{ResourceType: "AWS::S3::Bucket", ResourceName: "example-bucket"}, entries[0])
	})

	t.Run("S3 object ARN classifies as Object", func(t *testing.T) {
		builder := NewEventBuilder("s3", "PutObject", newCtx(), nil)
		response := map[string]interface{}{
			"Bucket": "arn:aws:s3:::example-bucket/key.txt",
		}
		entries := builder.Build(response, nil).Resources
		require.Len(t, entries, 1)
		assert.Equal(t, "AWS::S3::Object", entries[0].ResourceType)
	})
}

func TestResponseElementsPolicy(t *testing.T) {
	ctx := request.NewRequestContext(nil, nil, "123456789012", "us-east-1")
	response := map[string]interface{}{"TrailARN": "arn:aws:cloudtrail:us-east-1:123456789012:trail/example-trail"}

	write := NewEventBuilder("cloudtrail", "CreateTrail", ctx, nil).Build(response, nil)
	assert.NotNil(t, write.ResponseElements, "a write records its response elements")

	read := NewEventBuilder("cloudtrail", "GetTrail", ctx, nil).Build(response, nil)
	assert.Nil(t, read.ResponseElements, "a read-only API records null response elements")
}

func TestSanitizeParametersNested(t *testing.T) {
	params := map[string]interface{}{
		"PolicyDocument": map[string]interface{}{
			"Statement": []interface{}{
				map[string]interface{}{
					"Action":   "s3:PutObject",
					"Password": "nested-secret",
				},
			},
		},
		"Tags": []interface{}{
			map[string]interface{}{"Key": "token", "Value": "tok"},
		},
		"Name": "plain",
	}

	builder := NewEventBuilder("", "", nil, &request.ParsedRequest{Parameters: params})
	event := builder.Build(nil, nil)

	policy := event.RequestParameters["PolicyDocument"].(map[string]interface{})
	statement := policy["Statement"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "****", statement["Password"], "nested password redacted")
	assert.Equal(t, "s3:PutObject", statement["Action"])
	tag := event.RequestParameters["Tags"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "tok", tag["Value"], "a Value under a lowercased sensitive key redacts only the key's own value")
	assert.Equal(t, "plain", event.RequestParameters["Name"])
}

// The bare-name ARN resources (an SQS queue, an SNS topic) classify
// through the empty word — their resource carries no type-word segment —
// while separator-carrying resources keep classifying by their word.
func TestResourceTypeFromARN(t *testing.T) {
	tests := []struct {
		arn      string
		expected string
	}{
		{"arn:aws:sns:us-east-1:123456789012:MyTopic", "AWS::SNS::Topic"},
		{"arn:aws:sqs:us-east-1:123456789012:MyQueue", "AWS::SQS::Queue"},
		{"arn:aws:cloudtrail:us-east-1:123456789012:trail/my-trail", "AWS::CloudTrail::Trail"},
		{"arn:aws:states:us-east-1:123456789012:stateMachine:MySM", "AWS::StepFunctions::StateMachine"},
		{"arn:aws:logs:us-east-1:123456789012:log-group:my-group", "AWS::Logs::LogGroup"},
		{"arn:aws:s3:::my-bucket", "AWS::S3::Bucket"},
		{"arn:aws:s3:::my-bucket/key", "AWS::S3::Object"},
		{"arn:aws:iam::123456789012:user/alice", "AWS::IAM::User"},
		// An unmapped bare-name word under a known service stays empty:
		// the record format permits a resources entry carrying only the
		// ARN.
		{"arn:aws:kms:us-east-1:123456789012:unmapped-bare", ""},
		{"not-an-arn", ""},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, resourceTypeFromARN(tt.arn), tt.arn)
	}
}
