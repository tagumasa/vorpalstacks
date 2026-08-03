package scheduler

import (
	"strings"
	"testing"

	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

func validTarget() *schedulerstore.Target {
	return &schedulerstore.Target{
		Arn:     "arn:aws:sqs:us-east-1:123456789012:my-queue",
		RoleArn: "arn:aws:iam::123456789012:role/scheduler-role",
	}
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
			}
			_, err := ValidateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScheduleFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateScheduleFields_CronFieldCount(t *testing.T) {
	tests := []struct {
		expr    string
		wantErr bool
	}{
		{"cron(0 12 * * * ?)", false}, // 6 fields, valid
		{"cron(0 12 * * ? *)", false}, // 6 fields, valid
		{"cron(0 12 * * *)", true},    // 5 fields, invalid (H5)
		{"cron(0)", true},             // 1 field, invalid (H5)
		{"cron(* * * * * *)", false},  // 6 wildcards, valid
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: tt.expr,
				Target:             validTarget(),
			}
			_, err := ValidateScheduleFields(spec)
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
		{"at(2025-13-45T99:99:99)", true},  // month 13, hour 99 (H6)
		{"at(2025-02-30T10:00:00)", true},  // Feb 30 doesn't exist (H6)
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: tt.expr,
				Target:             validTarget(),
			}
			_, err := ValidateScheduleFields(spec)
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
				State:              tt.state,
			}
			_, err := ValidateScheduleFields(spec)
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
		{"FLEXIBLE", intPtr(15), false},
		{"FLEXIBLE", nil, true},          // FLEXIBLE requires MaximumWindowInMinutes
		{"FLEXIBLE", intPtr(0), true},    // below min
		{"FLEXIBLE", intPtr(1441), true}, // above max
		{"INVALID_MODE", nil, true},      // invalid enum (H3)
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
			_, err := ValidateScheduleFields(spec)
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
		{"valid", intPtr(60), intPtr(0), false},
		{"valid-max", intPtr(86400), intPtr(185), false},
		{"age-below-min", intPtr(59), nil, true},    // H2
		{"age-above-max", intPtr(86401), nil, true}, // H2
		{"retry-above-max", nil, intPtr(186), true}, // H2
		{"retry-below-min", nil, intPtr(-1), true},  // H2
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
			}
			_, err := ValidateScheduleFields(spec)
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
		Description:        strings.Repeat("a", 512),
	}
	if _, err := ValidateScheduleFields(spec); err != nil {
		t.Errorf("512 chars should pass, got err=%v", err)
	}

	spec.Description = strings.Repeat("a", 513)
	if _, err := ValidateScheduleFields(spec); err == nil {
		t.Error("513 chars should fail (M6)")
	}
}

func TestValidateScheduleFields_KmsKeyArn(t *testing.T) {
	tests := []struct {
		arn     string
		wantErr bool
	}{
		{"arn:aws:kms:us-east-1:123456789012:key/abcd-1234", false},
		{"arn:aws:sqs:us-east-1:123456789012:queue", true}, // wrong service
		{"not-an-arn", true}, // invalid format
		{"", false},          // empty is OK (optional)
	}
	for _, tt := range tests {
		t.Run(tt.arn, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:               "test-schedule",
				ScheduleExpression: "rate(1 minute)",
				Target:             validTarget(),
				KmsKeyArn:          tt.arn,
			}
			_, err := ValidateScheduleFields(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("KmsKeyArn: arn=%q wantErr=%v, got err=%v", tt.arn, tt.wantErr, err)
			}
		})
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
		{strings.Repeat("a", 51), true}, // too long (M3)
		{"", false},                     // empty is OK (optional)
	}
	for _, tt := range tests {
		t.Run(tt.tz, func(t *testing.T) {
			spec := &ScheduleSpec{
				Name:                       "test-schedule",
				ScheduleExpression:         "rate(1 minute)",
				Target:                     validTarget(),
				ScheduleExpressionTimezone: tt.tz,
			}
			_, err := ValidateScheduleFields(spec)
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
		{"end-before-start", "2025-12-31T00:00:00Z", "2025-01-01T00:00:00Z", true}, // M4
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
				StartDate:          tt.start,
				EndDate:            tt.end,
			}
			_, err := ValidateScheduleFields(spec)
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
		// ActionAfterCompletion intentionally omitted
	}
	validated, err := ValidateScheduleFields(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if validated.ActionAfterCompletion != schedulerstore.ActionAfterCompletionNone {
		t.Errorf("AAC default: got %q, want NONE", validated.ActionAfterCompletion)
	}

	// When explicitly set, should use the provided value.
	spec.ActionAfterCompletion = "DELETE"
	validated, err = ValidateScheduleFields(spec)
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
	store := schedulerstore.NewClientTokenStore()
	defer store.Stop()

	// First call claims the token.
	entry1, created1 := store.LookupOrClaim("token-1", "arn:schedule-1", "schedule")
	if !created1 {
		t.Fatal("first call should claim the token")
	}
	if entry1.ResourceArn != "arn:schedule-1" {
		t.Fatalf("expected arn:schedule-1, got %s", entry1.ResourceArn)
	}

	// Second call with same token returns the existing entry.
	entry2, created2 := store.LookupOrClaim("token-1", "arn:schedule-2", "schedule")
	if created2 {
		t.Fatal("second call should NOT claim (token already exists)")
	}
	if entry2.ResourceArn != "arn:schedule-1" {
		t.Fatalf("expected original arn:schedule-1, got %s", entry2.ResourceArn)
	}

	// Different token claims separately.
	_, created3 := store.LookupOrClaim("token-2", "arn:schedule-3", "schedule")
	if !created3 {
		t.Fatal("different token should claim")
	}
}

func TestClientTokenStore_Release(t *testing.T) {
	store := schedulerstore.NewClientTokenStore()
	defer store.Stop()

	store.LookupOrClaim("token-release", "arn:original", "schedule")
	store.Release("token-release")

	// After release, same token should claim again.
	_, created := store.LookupOrClaim("token-release", "arn:new", "schedule")
	if !created {
		t.Fatal("token should be claimable after Release")
	}
}

func intPtr(v int) *int { return &v }
