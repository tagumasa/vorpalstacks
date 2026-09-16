package scheduler

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

func validTarget() *schedulerstore.Target {
	return &schedulerstore.Target{
		Arn:     "arn:aws:sqs:us-east-1:123456789012:my-queue",
		RoleArn: "arn:aws:iam::123456789012:role/scheduler-role",
	}
}

func validFTW() *schedulerstore.FlexibleTimeWindow {
	return &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
}

func TestValidateScheduleFields_NamePattern(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid-name", false},
		{"valid_name", false},
		{"valid.name", false},
		{"VALID-NAME123", false},
		{"", true},
		{"has space", true},
		{"has/slash", true},
		{strings.Repeat("a", 65), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               tt.name,
				ScheduleExpression: "rate(1 minute)",
				Target:             validTarget(),
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateScheduleFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleFields_TargetRoleArnPattern(t *testing.T) {
	tests := []struct {
		roleArn string
		wantErr bool
	}{
		{"arn:aws:iam::123456789012:role/x", false},
		{"arn:aws:iam::123456789012:role/path/x", false},
		{"arn:aws-cn:iam::123456789012:role/x", false},
		{"arn:aws:iam::123456789012:role/", true},           // empty role name
		{"arn:aws:iam::12345678901:role/x", true},           // 11-digit account
		{"arn:aws:iam::1234567890123:role/x", true},         // 13-digit account
		{"arn:aws:iam::abc:role/x", true},                   // non-digit account
		{"arn:aws:iam:us-east-1:123456789012:role/x", true}, // non-empty region
		{"arn:awz:iam::123456789012:role/x", true},          // partition outside the aws family
		{"arn:aws-us-gov:iam::123456789012:role/x", true},   // two dash segments: outside the model's partition alternation
		{"arn:aws:lambda::123456789012:role/x", true},       // not the iam service
		{"arn:aws:iam::123456789012:role/bad name", true},   // character outside [\w+=,.@/-]
	}
	for _, tt := range tests {
		t.Run(tt.roleArn, func(t *testing.T) {
			target := validTarget()
			target.RoleArn = tt.roleArn
			spec := &ScheduleSpec{
				Name:               "valid-name",
				ScheduleExpression: "rate(1 minute)",
				Target:             target,
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateScheduleFields() RoleArn %q error = %v, wantErr %v", tt.roleArn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleFields_TargetArnLength(t *testing.T) {
	sqsPrefix := "arn:aws:sqs:us-east-1:123456789012:queue/"
	maxArn := sqsPrefix + strings.Repeat("q", MaxTargetArnLength-len(sqsPrefix))
	tooLongArn := sqsPrefix + strings.Repeat("q", MaxTargetArnLength-len(sqsPrefix)+1)
	rolePrefix := "arn:aws:iam::123456789012:role/"
	maxRole := rolePrefix + strings.Repeat("r", MaxTargetArnLength-len(rolePrefix))
	tooLongRole := rolePrefix + strings.Repeat("r", MaxTargetArnLength-len(rolePrefix)+1)

	tests := []struct {
		name    string
		arn     string
		roleArn string
		wantErr bool
	}{
		{"arn at maximum", maxArn, "arn:aws:iam::123456789012:role/x", false},
		{"arn over maximum", tooLongArn, "arn:aws:iam::123456789012:role/x", true},
		{"role arn at maximum", "arn:aws:sqs:us-east-1:123456789012:my-queue", maxRole, false},
		{"role arn over maximum", "arn:aws:sqs:us-east-1:123456789012:my-queue", tooLongRole, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := validTarget()
			target.Arn = tt.arn
			target.RoleArn = tt.roleArn
			spec := &ScheduleSpec{
				Name:               "valid-name",
				ScheduleExpression: "rate(1 minute)",
				Target:             target,
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateScheduleFields() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

// The templated-target Input contract the model documents: "If you are
// configuring a templated Lambda, AWS Step Functions, or Amazon EventBridge
// target, the input must be a well-formed JSON" — any JSON value. The
// object form is a delivery-time concern of the events family (PutEvents
// rejects a non-object Detail), not a creation-time restriction, and
// EventBridgeParameters is optional on the wire (no Target sub-parameter
// is @required in the model).
func TestValidateScheduleFields_EventBridgeTargetContract(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		eb      *schedulerstore.EventBridgeParameters
		wantErr bool
	}{
		{"parameters and object input", `{"k":"v"}`,
			&schedulerstore.EventBridgeParameters{Source: "my.app", DetailType: "MyEvent"}, false},
		{"parameters without input", "",
			&schedulerstore.EventBridgeParameters{Source: "my.app", DetailType: "MyEvent"}, false},
		{"missing parameters", `{"k":"v"}`, nil, false},
		{"missing parameters without input", "", nil, false},
		{"array input", `[1,2]`,
			&schedulerstore.EventBridgeParameters{Source: "my.app", DetailType: "MyEvent"}, false},
		{"scalar input", `"text"`,
			&schedulerstore.EventBridgeParameters{Source: "my.app", DetailType: "MyEvent"}, false},
		{"null input", `null`,
			&schedulerstore.EventBridgeParameters{Source: "my.app", DetailType: "MyEvent"}, false},
		{"malformed input", `{k:v}`,
			&schedulerstore.EventBridgeParameters{Source: "my.app", DetailType: "MyEvent"}, true},
		{"malformed input without parameters", `not-json`, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := validTarget()
			target.Arn = "arn:aws:events:us-east-1:123456789012:event-bus/my-bus"
			target.Input = tt.input
			target.EventBridgeParameters = tt.eb
			spec := &ScheduleSpec{
				Name:               "valid-name",
				ScheduleExpression: "rate(1 minute)",
				Target:             target,
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateScheduleFields() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				var awsErr *awserrors.AWSError
				if !errors.As(err, &awsErr) || !strings.Contains(awsErr.Message, "well-formed JSON") {
					t.Errorf("error must name the well-formed JSON contract, got %v", err)
				}
			}
		})
	}

	// The well-formed JSON requirement covers the other JSON-bound templated
	// families: a templated Lambda or Step Functions target with a non-JSON
	// Input is rejected at creation.
	for _, arn := range []string{
		"arn:aws:lambda:us-east-1:123456789012:function:fn",
		"arn:aws:states:us-east-1:123456789012:stateMachine:sm",
	} {
		target := validTarget()
		target.Arn = arn
		target.Input = "hello"
		spec := &ScheduleSpec{
			Name:               "valid-name",
			ScheduleExpression: "rate(1 minute)",
			Target:             target,
			FlexibleTimeWindow: validFTW(),
		}
		if _, err := validateScheduleFields(spec); err == nil {
			t.Errorf("validateScheduleFields() accepted non-JSON Input on templated target %s", arn)
		}
	}

	// "For all other target types, a JSON is not required": non-JSON-bound
	// delivery families accept free-form text and deliver it as-is.
	target := validTarget()
	target.Input = "free-form text"
	spec := &ScheduleSpec{
		Name:               "valid-name",
		ScheduleExpression: "rate(1 minute)",
		Target:             target,
		FlexibleTimeWindow: validFTW(),
	}
	if _, err := validateScheduleFields(spec); err != nil {
		t.Errorf("validateScheduleFields() rejected free-form Input on a non-JSON-bound target: %v", err)
	}
}

func TestValidateScheduleFields_CronFieldCount(t *testing.T) {
	tests := []struct {
		expr    string
		wantErr bool
	}{
		{"cron(0 12 ? * * *)", false}, // 6 fields, DOM ?, valid
		{"cron(0 12 * * ? *)", false}, // 6 fields, DOW ?, valid
		{"cron(0 12 * * *)", true},    // 5 fields, invalid
		{"cron(0)", true},             // 1 field, invalid
		{"cron(* * * * * *)", true},   // neither day field is ?, invalid
		{"cron(0 12 * * * ?)", true},  // neither day field is ?; the year ? is never examined, invalid
		{"cron(0 12 ? * * ?)", true},  // ? outside the two day fields (year), invalid
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: tt.expr,
				Target:             validTarget(),
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("cron field count: expr=%q error=%v wantErr=%v", tt.expr, err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleFields_AtSemanticDate(t *testing.T) {
	tests := []struct {
		expr    string
		wantErr bool
	}{
		{"at(2025-01-15T10:00:00)", false}, // valid date
		{"at(2025-13-45T99:99:99)", true},  // month 13, hour 99
		{"at(2025-02-30T10:00:00)", true},  // Feb 30 doesn't exist
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: tt.expr,
				Target:             validTarget(),
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("at() semantic: expr=%q error=%v wantErr=%v", tt.expr, err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleFields_StateEnum(t *testing.T) {
	tests := []struct {
		state   string
		wantErr bool
	}{
		{"ENABLED", false},
		{"DISABLED", false},
		{"", false}, // defaults to ENABLED
		{"INVALID", true},
		{"enabled", true},
	}
	for _, tt := range tests {
		t.Run(tt.state, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             validTarget(),
				FlexibleTimeWindow: validFTW(),
				State:              tt.state,
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("State enum: state=%q error=%v wantErr=%v", tt.state, err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleFields_FlexibleTimeWindowModeEnum(t *testing.T) {
	tests := []struct {
		mode    string
		maxWin  *int
		wantErr bool
	}{
		{"OFF", nil, false},
		{"OFF", intPtr(15), false},  // a provided window is meaningless in OFF mode but must be in range
		{"OFF", intPtr(1441), true}, // the @range(1,1440) trait is unconditional on the shape
		{"OFF", intPtr(0), true},    // below the unconditional minimum
		{"FLEXIBLE", intPtr(15), false},
		{"FLEXIBLE", nil, true},          // FLEXIBLE requires MaximumWindowInMinutes
		{"FLEXIBLE", intPtr(0), true},    // below min
		{"FLEXIBLE", intPtr(1441), true}, // above max
		{"INVALID_MODE", nil, true},      // invalid enum
	}
	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			ftw := &schedulerstore.FlexibleTimeWindow{
				Mode: schedulerstore.FlexibleTimeWindowMode(tt.mode),
			}
			if tt.maxWin != nil {
				ftw.MaximumWindowInMinutes = tt.maxWin
			}
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             validTarget(),
				FlexibleTimeWindow: ftw,
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("FTW Mode: mode=%q wantErr=%v, got err=%v", tt.mode, tt.wantErr, err)
			}
		})
	}
}

func TestValidateScheduleFields_RetryPolicyRanges(t *testing.T) {
	tests := []struct {
		name     string
		eventAge *int
		retryAtt *int
		wantErr  bool
	}{
		{"valid", intPtr(MinRetryPolicyEventAgeSeconds), intPtr(0), false},
		{"valid-max", intPtr(MaxRetryPolicyEventAgeSeconds), intPtr(MaxRetryPolicyAttempts), false},
		{"age-below-min", intPtr(MinRetryPolicyEventAgeSeconds - 1), nil, true}, // below minimum age
		{"age-above-max", intPtr(MaxRetryPolicyEventAgeSeconds + 1), nil, true}, // above maximum age
		{"retry-above-max", nil, intPtr(MaxRetryPolicyAttempts + 1), true},      // above maximum retry attempts
		{"retry-below-min", nil, intPtr(-1), true},                              // below minimum retry attempts
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := validTarget()
			target.RetryPolicy = &schedulerstore.RetryPolicy{}
			if tt.eventAge != nil {
				target.RetryPolicy.MaximumEventAgeInSeconds = tt.eventAge
			}
			if tt.retryAtt != nil {
				target.RetryPolicy.MaximumRetryAttempts = tt.retryAtt
			}
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             target,
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetryPolicy: name=%q wantErr=%v, got err=%v", tt.name, tt.wantErr, err)
			}
		})
	}
}

func TestValidateScheduleFields_DescriptionLength(t *testing.T) {
	spec := &ScheduleSpec{
		Name:               "test-schedule",
		ScheduleExpression: "rate(1 minute)",
		Target:             validTarget(),
		FlexibleTimeWindow: validFTW(),
		Description:        strings.Repeat("a", 512),
	}
	if _, err := validateScheduleFields(spec); err != nil {
		t.Errorf("512 chars should pass, got err=%v", err)
	}

	spec.Description = strings.Repeat("a", 513)
	if _, err := validateScheduleFields(spec); err == nil {
		t.Error("513 chars should fail")
	}
}

func TestValidateScheduleFields_KmsKeyArn(t *testing.T) {
	tests := []struct {
		arn     string
		wantErr bool
	}{
		{"arn:aws:kms:us-east-1:123456789012:key/abcd-1234", false},
		{"arn:aws:kms:us-east-1:123456789012:alias/my-alias", false}, // alias form
		{"arn:aws:sqs:us-east-1:123456789012:queue", true},           // wrong service
		{"arn:aws:kms:us-east-1:123456789012:secret/xyz", true},      // wrong resource type
		{"arn:aws:kms:us-east-1:12345678901:key/abcd", true},         // 11-digit account
		{"arn:aws:kms:US-EAST-1:123456789012:key/abcd", true},        // uppercase region outside [a-z0-9-]
		{"arn:awz:kms:us-east-1:123456789012:key/abcd", true},        // partition outside the aws family
		{"arn:aws:kms:us-east-1:123456789012:key/abcd/efgh", true},   // '/' outside the key-id charset
		{"not-an-arn", true}, // invalid format
		{"", false},          // empty is OK (optional)
	}
	for _, tt := range tests {
		t.Run(tt.arn, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             validTarget(),
				FlexibleTimeWindow: validFTW(),
				KmsKeyArn:          tt.arn,
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("KmsKeyArn: arn=%q wantErr=%v, got err=%v", tt.arn, tt.wantErr, err)
			}
		})
	}

	// KmsKeyArn @length(1, 2048): the bound is checked on the provided ARN.
	atMax := "arn:aws:kms:us-east-1:123456789012:key/" + strings.Repeat("k", maxKmsKeyArnLength-len("arn:aws:kms:us-east-1:123456789012:key/"))
	overMax := atMax + "k"
	spec := &ScheduleSpec{
		Name:               "test-schedule",
		ScheduleExpression: "rate(1 minute)",
		Target:             validTarget(),
		FlexibleTimeWindow: validFTW(),
		KmsKeyArn:          atMax,
	}
	if _, err := validateScheduleFields(spec); err != nil {
		t.Errorf("KmsKeyArn at the 2048 maximum rejected: %v", err)
	}
	spec.KmsKeyArn = overMax
	if _, err := validateScheduleFields(spec); err == nil {
		t.Error("KmsKeyArn over the 2048 maximum accepted")
	}
}

// TestValidateScheduleFields_DeadLetterConfigArn pins the member @pattern
// (an SQS queue ARN, applied verbatim) and the ResourceArn @length(1, 1600)
// bound on the dead-letter configuration.
func TestValidateScheduleFields_DeadLetterConfigArn(t *testing.T) {
	tests := []struct {
		arn     string
		wantErr bool
	}{
		{"arn:aws:sqs:us-east-1:123456789012:my-dlq", false},
		{"arn:aws-cn:sqs:cn-north-1:123456789012:my-dlq", false},
		{"arn:aws:sns:us-east-1:123456789012:my-topic", true},   // wrong service
		{"arn:aws:sqs:us-east-1:12345678901:my-dlq", true},      // 11-digit account
		{"arn:aws:sqs:us-east-1:123456789012:", true},           // empty queue name
		{"arn:aws:sqs:us-east-1:123456789012:bad name", true},   // space outside [a-zA-Z0-9-_]
		{"arn:aws:sqs:us-east-1:123456789012:queue.fifo", true}, // '.' outside the queue-name charset
		{"not-an-arn", true}, // no ARN grammar at all
	}
	for _, tt := range tests {
		t.Run(tt.arn, func(t *testing.T) {
			target := validTarget()
			target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{Arn: tt.arn}
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             target,
				FlexibleTimeWindow: validFTW(),
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeadLetterConfig.Arn: arn=%q wantErr=%v, got err=%v", tt.arn, tt.wantErr, err)
			}
		})
	}

	prefix := "arn:aws:sqs:us-east-1:123456789012:"
	atMax := prefix + strings.Repeat("q", MaxTargetArnLength-len(prefix))
	target := validTarget()
	target.DeadLetterConfig = &schedulerstore.DeadLetterConfig{Arn: atMax}
	spec := &ScheduleSpec{
		Name:               "test-schedule",
		ScheduleExpression: "rate(1 minute)",
		Target:             target,
		FlexibleTimeWindow: validFTW(),
	}
	if _, err := validateScheduleFields(spec); err != nil {
		t.Errorf("DeadLetterConfig.Arn at the 1600 maximum rejected: %v", err)
	}
	target.DeadLetterConfig.Arn = atMax + "q"
	if _, err := validateScheduleFields(spec); err == nil {
		t.Error("DeadLetterConfig.Arn over the 1600 maximum accepted")
	}
}

func TestValidateScheduleFields_Timezone(t *testing.T) {
	tests := []struct {
		tz      string
		wantErr bool
	}{
		{"UTC", false},
		{"America/New_York", false},
		{"Asia/Tokyo", false},
		{"Invalid/NotReal", true},       // not in IANA database
		{"Local", true},                 // Go runtime alias, not an IANA zone name
		{strings.Repeat("a", 51), true}, // too long
		{"", false},                     // empty is OK (optional)
	}
	for _, tt := range tests {
		t.Run(tt.tz, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:                       "test-schedule",
				ScheduleExpression:         "rate(1 minute)",
				Target:                     validTarget(),
				FlexibleTimeWindow:         validFTW(),
				ScheduleExpressionTimezone: tt.tz,
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timezone: tz=%q wantErr=%v, got err=%v", tt.tz, tt.wantErr, err)
			}
		})
	}
}

func TestValidateScheduleFields_StartEndDateOrdering(t *testing.T) {
	tests := []struct {
		name    string
		start   string
		end     string
		wantErr bool
	}{
		{"valid-order", "2025-01-01T00:00:00Z", "2025-12-31T23:59:59Z", false},
		{"end-before-start", "2025-12-31T00:00:00Z", "2025-01-01T00:00:00Z", true}, // end date precedes start date
		{"same-time", "2025-01-01T00:00:00Z", "2025-01-01T00:00:00Z", false},
		{"start-only", "2025-01-01T00:00:00Z", "", false},
		{"end-only", "", "2025-01-01T00:00:00Z", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             validTarget(),
				FlexibleTimeWindow: validFTW(),
				StartDate:          tt.start,
				EndDate:            tt.end,
			}
			_, err := validateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date ordering: name=%q wantErr=%v, got err=%v", tt.name, tt.wantErr, err)
			}
		})
	}
}

func TestValidateScheduleFields_ActionAfterCompletionDefault(t *testing.T) {
	// When ActionAfterCompletion is omitted, the validator must default to
	// NONE because AWS UpdateSchedule is a PUT operation that resets
	// unspecified fields to the documented defaults.
	spec := &ScheduleSpec{
		Name:               "test-schedule",
		ScheduleExpression: "rate(1 minute)",
		Target:             validTarget(),
		FlexibleTimeWindow: validFTW(),
		// ActionAfterCompletion intentionally omitted
	}
	validated, err := validateScheduleFields(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if validated.ActionAfterCompletion != schedulerstore.ActionAfterCompletionNone {
		t.Errorf("AAC default: got %q, want NONE", validated.ActionAfterCompletion)
	}

	// When explicitly set, should use the provided value.
	spec.ActionAfterCompletion = "DELETE"
	validated, err = validateScheduleFields(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if validated.ActionAfterCompletion != schedulerstore.ActionAfterCompletionDelete {
		t.Errorf("AAC explicit: got %q, want DELETE", validated.ActionAfterCompletion)
	}
}

func TestValidateClientToken(t *testing.T) {
	tests := []struct {
		token   string
		wantErr bool
	}{
		{"abc-123_DEF", false},
		{"a", false},
		{strings.Repeat("a", 64), false},
		{"", true},                      // too short
		{strings.Repeat("a", 65), true}, // too long
		{"has space", true},             // invalid char
		{"has@symbol", true},            // invalid char
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			err := validateClientToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateClientToken(%q) err=%v, wantErr=%v", tt.token, err, tt.wantErr)
			}
		})
	}
}

func TestClientTokenStore_LookupOrClaim(t *testing.T) {
	bs, cleanup := testTokenBaseStore(t)
	defer cleanup()
	store := schedulerstore.NewClientTokenStore(bs)
	defer store.Stop()

	// First call claims the token.
	entry1, created1 := store.LookupOrClaim("token-1", "arn:schedule-1", "schedule")
	if !created1 {
		t.Fatal("first call should claim the token")
	}
	if entry1.ResourceArn != "arn:schedule-1" {
		t.Fatalf("expected arn:schedule-1, got %s", entry1.ResourceArn)
	}

	// A replay of the same request (same scope) returns the existing entry.
	entry2, created2 := store.LookupOrClaim("token-1", "arn:schedule-1", "schedule")
	if created2 {
		t.Fatal("replay of the same scope should NOT claim (token already exists)")
	}
	if entry2.ResourceArn != "arn:schedule-1" {
		t.Fatalf("expected original arn:schedule-1, got %s", entry2.ResourceArn)
	}

	// The same token on another resource is a different request and claims
	// separately instead of replaying the first request's outcome.
	_, created3 := store.LookupOrClaim("token-1", "arn:schedule-2", "schedule")
	if !created3 {
		t.Fatal("same token on another resource should claim")
	}

	// The same token under another operation is also a different request.
	_, created4 := store.LookupOrClaim("token-1", "arn:schedule-1", "schedule-update")
	if !created4 {
		t.Fatal("same token under another operation should claim")
	}

	// Different token claims separately.
	_, created5 := store.LookupOrClaim("token-2", "arn:schedule-3", "schedule")
	if !created5 {
		t.Fatal("different token should claim")
	}
}

func TestClientTokenStore_Release(t *testing.T) {
	bs, cleanup := testTokenBaseStore(t)
	defer cleanup()
	store := schedulerstore.NewClientTokenStore(bs)
	defer store.Stop()

	store.LookupOrClaim("token-release", "arn:original", "schedule")

	// Releasing a different scope must not touch the original claim.
	store.Release("token-release", "arn:other", "schedule")
	_, stillHeld := store.LookupOrClaim("token-release", "arn:original", "schedule")
	if stillHeld {
		t.Fatal("original claim should still be held after a foreign-scope release")
	}

	store.Release("token-release", "arn:original", "schedule")

	// After release, same scope should claim again.
	_, created := store.LookupOrClaim("token-release", "arn:original", "schedule")
	if !created {
		t.Fatal("token should be claimable after Release")
	}
}

func intPtr(v int) *int { return &v }

// testTokenBaseStore creates a Pebble-backed BaseStore for ClientTokenStore
// tests. The temp directory is managed by the testing.T and cleaned up
// automatically when the test finishes.
func testTokenBaseStore(t *testing.T) (*storecommon.BaseStore, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	s, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatalf("storage.Open(%q) failed: %v", tmpDir, err)
	}
	bs := storecommon.NewBaseStore(s.Bucket("scheduler-tokens-test"), "scheduler-tokens")
	return bs, func() { s.Close() }
}

// TestValidateEcsParametersUnicodeLengths pins that the pattern-less ECS
// members (Group @length(1,255), ReferenceId @length(max 1024)) count
// Unicode characters, so rune-legal multibyte values must not be rejected
// on byte length.
func TestValidateEcsParametersUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes
	base := func() *schedulerstore.EcsParameters {
		return &schedulerstore.EcsParameters{
			TaskDefinitionArn: "arn:aws:ecs:us-east-1:123456789012:task-definition/family:1",
		}
	}

	ecs := base()
	ecs.Group = strings.Repeat(cjk, 255)
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("255-character CJK group rejected: %v", err)
	}
	ecs = base()
	ecs.Group = strings.Repeat(cjk, 256)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("256-character CJK group accepted")
	}

	ecs = base()
	ecs.ReferenceId = strings.Repeat(cjk, 1024)
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("1024-character CJK reference id rejected: %v", err)
	}
	ecs = base()
	ecs.ReferenceId = strings.Repeat(cjk, 1025)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("1025-character CJK reference id accepted")
	}

	// CapacityProviderStrategy item names at the 255-character maximum.
	ecs = base()
	ecs.CapacityProviderStrategy = []schedulerstore.CapacityProviderStrategyItem{{
		CapacityProvider: strings.Repeat(cjk, 255),
	}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("255-character CJK capacity provider rejected: %v", err)
	}
	ecs.CapacityProviderStrategy[0].CapacityProvider = strings.Repeat(cjk, 256)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("256-character CJK capacity provider accepted")
	}

	// Subnet and security-group IDs at the 1000-character maximum.
	ecs = base()
	ecs.NetworkConfiguration = &schedulerstore.NetworkConfiguration{
		AwsVpcConfiguration: &schedulerstore.AwsVpcConfiguration{
			Subnets:        []string{strings.Repeat(cjk, 1000)},
			SecurityGroups: []string{strings.Repeat(cjk, 1000)},
		},
	}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("1000-character CJK subnet/security-group ids rejected: %v", err)
	}
	ecs.NetworkConfiguration.AwsVpcConfiguration.Subnets[0] = strings.Repeat(cjk, 1001)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("1001-character CJK subnet id accepted")
	}

	// TaskDefinitionArn @length(1, 1600): an ARN at the maximum is
	// accepted, one over it rejected, and a multibyte ARN whose UTF-8 byte
	// length exceeds the bound stays accepted while its character count is
	// inside it.
	arnPrefix := "arn:aws:ecs:us-east-1:123456789012:task-definition/"
	ecs = base()
	ecs.TaskDefinitionArn = arnPrefix + strings.Repeat("f", MaxTargetArnLength-utf8.RuneCountInString(arnPrefix))
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("task definition ARN at the 1600 maximum rejected: %v", err)
	}
	ecs.TaskDefinitionArn = arnPrefix + strings.Repeat("f", MaxTargetArnLength-utf8.RuneCountInString(arnPrefix)+1)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("task definition ARN over the 1600 maximum accepted")
	}
	ecs.TaskDefinitionArn = arnPrefix + strings.Repeat(cjk, 1000)
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("1052-character multibyte task definition ARN (3000+ bytes) rejected: %v", err)
	}
}

// TestValidateEcsParametersPlacementLengths pins the placement member
// @length maxima: PlacementConstraintExpression @length(max 2000) and
// PlacementStrategyField @length(max 255).
func TestValidateEcsParametersPlacementLengths(t *testing.T) {
	base := func() *schedulerstore.EcsParameters {
		return &schedulerstore.EcsParameters{
			TaskDefinitionArn: "arn:aws:ecs:us-east-1:123456789012:task-definition/family:1",
		}
	}

	ecs := base()
	ecs.PlacementConstraints = []schedulerstore.PlacementConstraint{{
		Type:       "memberOf",
		Expression: strings.Repeat("a", maxPlacementExpressionLength),
	}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("placement expression at the 2000 maximum rejected: %v", err)
	}
	ecs.PlacementConstraints[0].Expression = strings.Repeat("a", maxPlacementExpressionLength+1)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("placement expression over the 2000 maximum accepted")
	}

	ecs = base()
	ecs.PlacementStrategy = []schedulerstore.PlacementStrategy{{
		Type:  "spread",
		Field: strings.Repeat("a", maxPlacementStrategyFieldLength),
	}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("placement field at the 255 maximum rejected: %v", err)
	}
	ecs.PlacementStrategy[0].Field = strings.Repeat("a", maxPlacementStrategyFieldLength+1)
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("placement field over the 255 maximum accepted")
	}
}

// TestValidateEcsParametersTagsBounds pins the EcsParameters.Tags bounds:
// the list holds at most fifty TagMap entries, and every pair carries a
// TagKey of 1-128 and a TagValue of 1-256 characters (the @length rune
// basis). The model places no single-pair bound on a TagMap, so an entry
// carrying several pairs is valid input and each pair is checked.
func TestValidateEcsParametersTagsBounds(t *testing.T) {
	base := func() *schedulerstore.EcsParameters {
		return &schedulerstore.EcsParameters{
			TaskDefinitionArn: "arn:aws:ecs:us-east-1:123456789012:task-definition/family:1",
		}
	}

	// Absent member and empty list stay accepted (Tags @length(0, 50));
	// an entry with no pairs carries nothing to violate either.
	ecs := base()
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("absent Tags rejected: %v", err)
	}
	ecs.Tags = []map[string]string{}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("empty Tags list rejected: %v", err)
	}
	ecs.Tags = []map[string]string{{}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("empty TagMap entry rejected: %v", err)
	}

	// List bound: fifty entries accepted, fifty-one rejected.
	ecs.Tags = make([]map[string]string, MaxEcsTagsItems)
	for i := range ecs.Tags {
		ecs.Tags[i] = map[string]string{"env": "prod"}
	}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("fifty TagMap entries rejected: %v", err)
	}
	ecs.Tags = append(ecs.Tags, map[string]string{"extra": "entry"})
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("fifty-one TagMap entries accepted")
	}

	// Key bound at the 128-character maximum, counted in characters.
	ecs.Tags = []map[string]string{{strings.Repeat("k", maxTagKeyLength): "v"}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("128-character key rejected: %v", err)
	}
	ecs.Tags = []map[string]string{{strings.Repeat("k", maxTagKeyLength+1): "v"}}
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("129-character key accepted")
	}
	ecs.Tags = []map[string]string{{strings.Repeat("\u65e5", maxTagKeyLength): "v"}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("128-character multibyte key (384 bytes) rejected: %v", err)
	}

	// Value bound at the 256-character maximum, counted in characters.
	ecs.Tags = []map[string]string{{"k": strings.Repeat("v", maxTagValueLength)}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("256-character value rejected: %v", err)
	}
	ecs.Tags = []map[string]string{{"k": strings.Repeat("v", maxTagValueLength+1)}}
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("257-character value accepted")
	}
	ecs.Tags = []map[string]string{{"k": strings.Repeat("\u65e5", maxTagValueLength)}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("256-character multibyte value (768 bytes) rejected: %v", err)
	}

	// Minimum lengths: an empty key or value violates the TagKey and
	// TagValue @length(1, ...) minima.
	ecs.Tags = []map[string]string{{"": "v"}}
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("empty key accepted")
	}
	ecs.Tags = []map[string]string{{"k": ""}}
	if err := validateEcsParameters(ecs); err == nil {
		t.Error("empty value accepted")
	}

	// No single-pair rule: one entry carrying several pairs stays valid.
	ecs.Tags = []map[string]string{{"team": "core", "env": "prod"}}
	if err := validateEcsParameters(ecs); err != nil {
		t.Errorf("multi-pair TagMap entry rejected: %v", err)
	}
}

// TestValidateEcsParametersSecurityGroupsMinimum pins that an explicitly
// provided empty securityGroups list is rejected — the model constrains
// the provided list to 1-5 entries — while an omitted member (nil)
// selects the VPC default security group and stays accepted.
func TestValidateEcsParametersSecurityGroupsMinimum(t *testing.T) {
	base := func(securityGroups []string) *schedulerstore.EcsParameters {
		return &schedulerstore.EcsParameters{
			TaskDefinitionArn: "arn:aws:ecs:us-east-1:123456789012:task-definition/family:1",
			NetworkConfiguration: &schedulerstore.NetworkConfiguration{
				AwsVpcConfiguration: &schedulerstore.AwsVpcConfiguration{
					Subnets:        []string{"subnet-0123456789abcdef0"},
					SecurityGroups: securityGroups,
				},
			},
		}
	}

	if err := validateEcsParameters(base(nil)); err != nil {
		t.Errorf("omitted securityGroups rejected: %v", err)
	}
	if err := validateEcsParameters(base([]string{"sg-0123456789abcdef0"})); err != nil {
		t.Errorf("single security group rejected: %v", err)
	}
	if err := validateEcsParameters(base([]string{})); err == nil {
		t.Error("explicitly empty securityGroups accepted")
	}
}

// TestParseAwsVpcConfigurationSecurityGroupsPresence pins that an
// explicitly empty securityGroups array survives parsing as a non-nil
// empty list while an omitted member stays nil, so the validator can
// enforce the model minimum only for provided lists.
func TestParseAwsVpcConfigurationSecurityGroupsPresence(t *testing.T) {
	present, err := parseAwsVpcConfiguration(map[string]interface{}{
		"Subnets":        []interface{}{"subnet-0123456789abcdef0"},
		"SecurityGroups": []interface{}{},
	})
	if err != nil {
		t.Fatalf("parseAwsVpcConfiguration: %v", err)
	}
	if present.SecurityGroups == nil || len(present.SecurityGroups) != 0 {
		t.Errorf("explicitly empty securityGroups = %#v, want non-nil empty", present.SecurityGroups)
	}

	absent, err := parseAwsVpcConfiguration(map[string]interface{}{
		"Subnets": []interface{}{"subnet-0123456789abcdef0"},
	})
	if err != nil {
		t.Fatalf("parseAwsVpcConfiguration: %v", err)
	}
	if absent.SecurityGroups != nil {
		t.Errorf("omitted securityGroups = %#v, want nil", absent.SecurityGroups)
	}
}

// TestValidateScheduleGroupTagsUnicodeLengths pins that the TagKey and
// TagValue @length bounds count Unicode characters: a multibyte key at the
// 128-character maximum (384 bytes) is accepted.
func TestValidateScheduleGroupTagsUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"
	if err := ValidateScheduleGroupTags([]tagutil.Tag{
		{Key: strings.Repeat(cjk, 128), Value: "v"},
	}); err != nil {
		t.Errorf("128-character CJK tag key rejected: %v", err)
	}
	if err := ValidateScheduleGroupTags([]tagutil.Tag{
		{Key: strings.Repeat(cjk, 129), Value: "v"},
	}); err == nil {
		t.Error("129-character CJK tag key accepted")
	}
	if err := ValidateScheduleGroupTags([]tagutil.Tag{
		{Key: "k", Value: strings.Repeat(cjk, 256)},
	}); err != nil {
		t.Errorf("256-character CJK tag value rejected: %v", err)
	}
	if err := ValidateScheduleGroupTags([]tagutil.Tag{
		{Key: "k", Value: strings.Repeat(cjk, 257)},
	}); err == nil {
		t.Error("257-character CJK tag value accepted")
	}
}

// TestValidateTargetUnicodeLengths pins the character basis of the
// target-member @length bounds whose values can legally carry multibyte
// text: Target.Arn (pattern-less), KinesisParameters.PartitionKey
// (pattern-less), and the EventBridge DetailType / Source pair (the Source
// pattern restricts only the first character, leaving the tail free for
// multibyte text).
func TestValidateTargetUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	// 1000 CJK characters = 3000 bytes, inside the 1600-character maximum.
	target := validTarget()
	target.Arn = "arn:aws:sqs:us-east-1:123456789012:" + strings.Repeat(cjk, 1000)
	if err := validateTarget(target); err != nil {
		t.Errorf("1000-character CJK target ARN rejected: %v", err)
	}

	kin := validTarget()
	kin.Arn = "arn:aws:kinesis:us-east-1:123456789012:stream/pin"
	kin.KinesisParameters = &schedulerstore.KinesisParameters{
		PartitionKey: strings.Repeat(cjk, 256),
	}
	if err := validateTarget(kin); err != nil {
		t.Errorf("256-character CJK partition key rejected: %v", err)
	}
	kin.KinesisParameters.PartitionKey = strings.Repeat(cjk, 257)
	if err := validateTarget(kin); err == nil {
		t.Error("257-character CJK partition key accepted")
	}

	eb := validTarget()
	eb.Arn = "arn:aws:events:us-east-1:123456789012:event-bus/default"
	eb.EventBridgeParameters = &schedulerstore.EventBridgeParameters{
		Source:     "a" + strings.Repeat(cjk, 255),
		DetailType: strings.Repeat(cjk, 128),
	}
	if err := validateTarget(eb); err != nil {
		t.Errorf("at-maximum CJK DetailType/Source rejected: %v", err)
	}
	eb.EventBridgeParameters.DetailType = strings.Repeat(cjk, 129)
	if err := validateTarget(eb); err == nil {
		t.Error("129-character CJK DetailType accepted")
	}
	eb.EventBridgeParameters.DetailType = strings.Repeat(cjk, 128)
	eb.EventBridgeParameters.Source = "a" + strings.Repeat(cjk, 256)
	if err := validateTarget(eb); err == nil {
		t.Error("257-character CJK Source accepted")
	}
}

// TestValidateSqsParametersEmptyStruct pins that SqsParameters without a
// MessageGroupId is valid input — the member carries no @required in the
// model, so {"sqsParameters": {}} is accepted, with the FIFO fallback to
// the schedule name applied at delivery — while a present MessageGroupId
// is bounded at 128 characters.
func TestValidateSqsParametersEmptyStruct(t *testing.T) {
	if err := validateSqsParameters(&schedulerstore.SqsParameters{}); err != nil {
		t.Errorf("empty SqsParameters rejected: %v", err)
	}

	sqs := validTarget()
	sqs.SqsParameters = &schedulerstore.SqsParameters{}
	if err := validateTarget(sqs); err != nil {
		t.Errorf("SQS target with empty SqsParameters rejected: %v", err)
	}

	cjk := "\u65e5"
	withValue := &schedulerstore.SqsParameters{MessageGroupId: strings.Repeat(cjk, 128)}
	if err := validateSqsParameters(withValue); err != nil {
		t.Errorf("128-character CJK MessageGroupId rejected: %v", err)
	}
	withValue.MessageGroupId = strings.Repeat(cjk, 129)
	if err := validateSqsParameters(withValue); err == nil {
		t.Error("129-character CJK MessageGroupId accepted")
	}
}

// TestValidateNextTokenUnicodeLength pins the character basis of the
// NextToken @length(1, 2048) bound.
func TestValidateNextTokenUnicodeLength(t *testing.T) {
	cjk := "\u65e5"
	if err := validateNextToken(strings.Repeat(cjk, 2048)); err != nil {
		t.Errorf("2048-character token rejected: %v", err)
	}
	if err := validateNextToken(strings.Repeat(cjk, 2049)); err == nil {
		t.Error("2049-character token accepted")
	}
}

// TestValidateClientTokenMultibyteRejected pins the adjudication for the
// pattern-gated @length members: the bound counts characters now, but the
// ClientToken pattern's ASCII charset still rejects multibyte values.
func TestValidateClientTokenMultibyteRejected(t *testing.T) {
	if err := validateClientToken(strings.Repeat("\u65e5", 64)); err == nil {
		t.Error("64-character CJK ClientToken accepted")
	}
}

// TestValidateScheduleFields_FlexibleTimeWindowRequired pins that an absent
// FlexibleTimeWindow (a required member of CreateScheduleInput and
// UpdateScheduleInput) and an absent Mode (a required member of the
// FlexibleTimeWindow shape) are rejected as ValidationException instead of
// silently defaulting to OFF.
func TestValidateScheduleFields_FlexibleTimeWindowRequired(t *testing.T) {
	spec := &ScheduleSpec{
		Name:               "ftw-required",
		ScheduleExpression: "rate(5 minutes)",
		Target:             validTarget(),
	}
	_, err := validateScheduleFields(spec)
	if err == nil || !strings.Contains(err.Error(), "FlexibleTimeWindow is required") {
		t.Errorf("nil FlexibleTimeWindow: got error %v, want 'FlexibleTimeWindow is required'", err)
	}

	spec.FlexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{}
	_, err = validateScheduleFields(spec)
	if err == nil || !strings.Contains(err.Error(), "FlexibleTimeWindow.Mode is required") {
		t.Errorf("empty Mode: got error %v, want 'FlexibleTimeWindow.Mode is required'", err)
	}

	spec.FlexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	if _, err = validateScheduleFields(spec); err != nil {
		t.Errorf("explicit OFF window rejected: %v", err)
	}
}

// TestInternalServerErrorShape pins the wire shape of the scheduler's
// server-fault error: the Smithy model names it InternalServerException
// (smithy.api#error "server", httpError 500), so SDK clients matching
// *types.InternalServerException must recognise it.
func TestInternalServerErrorShape(t *testing.T) {
	if ErrInternalServer.Code != "InternalServerException" {
		t.Errorf("error code = %q, want InternalServerException", ErrInternalServer.Code)
	}
	if ErrInternalServer.HTTPStatus != 500 {
		t.Errorf("HTTP status = %d, want 500", ErrInternalServer.HTTPStatus)
	}
}
