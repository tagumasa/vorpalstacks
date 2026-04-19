package eventbridge

import (
	"testing"
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

func TestEventBus(t *testing.T) {
	bus := &EventBus{
		Name:      "my-event-bus",
		ARN:       "arn:aws:events:us-east-1:123456789012:event-bus/my-event-bus",
		Region:    "us-east-1",
		AccountID: "123456789012",
		Tags: []types.Tag{
			{Key: "env", Value: "prod"},
		},
	}

	if bus.Name != "my-event-bus" {
		t.Errorf("Name = %v, want 'my-event-bus'", bus.Name)
	}
	if bus.ARN != "arn:aws:events:us-east-1:123456789012:event-bus/my-event-bus" {
		t.Errorf("ARN = %v, want 'arn:aws:events:us-east-1:123456789012:event-bus/my-event-bus'", bus.ARN)
	}
	if len(bus.Tags) != 1 {
		t.Errorf("len(Tags) = %v, want 1", len(bus.Tags))
	}
}

func TestRule(t *testing.T) {
	rule := &Rule{
		Name:               "my-rule",
		ARN:                "arn:aws:events:us-east-1:123456789012:rule/my-bus/my-rule",
		Region:             "us-east-1",
		AccountID:          "123456789012",
		EventBusName:       "my-bus",
		EventPattern:       `{"source":["my.source"]}`,
		ScheduleExpression: "rate(5 minutes)",
		State:              RuleStateEnabled,
	}

	if rule.Name != "my-rule" {
		t.Errorf("Name = %v, want 'my-rule'", rule.Name)
	}
	if rule.EventBusName != "my-bus" {
		t.Errorf("EventBusName = %v, want 'my-bus'", rule.EventBusName)
	}
	if rule.State != RuleStateEnabled {
		t.Errorf("State = %v, want 'ENABLED'", rule.State)
	}
}

func TestArchive(t *testing.T) {
	archive := &Archive{
		Name:           "my-archive",
		ARN:            "arn:aws:events:us-east-1:123456789012:archive/my-archive",
		Region:         "us-east-1",
		AccountID:      "123456789012",
		EventBusName:   "my-bus",
		EventSourceARN: "arn:aws:events:us-east-1:123456789012:event-source/my-source",
		RetentionDays:  30,
		State:          ArchiveStateEnabled,
		EventCount:     1000,
		SizeBytes:      50000,
	}

	if archive.Name != "my-archive" {
		t.Errorf("Name = %v, want 'my-archive'", archive.Name)
	}
	if archive.RetentionDays != 30 {
		t.Errorf("RetentionDays = %v, want 30", archive.RetentionDays)
	}
	if archive.State != ArchiveStateEnabled {
		t.Errorf("State = %v, want 'ENABLED'", archive.State)
	}
}

func TestConnection(t *testing.T) {
	conn := &Connection{
		Name:              "my-connection",
		ARN:               "arn:aws:events:us-east-1:123456789012:connection/my-connection",
		Region:            "us-east-1",
		AccountID:         "123456789012",
		AuthorizationType: "API_KEY",
		State:             ConnectionStateAuthorized,
	}

	if conn.Name != "my-connection" {
		t.Errorf("Name = %v, want 'my-connection'", conn.Name)
	}
	if conn.AuthorizationType != "API_KEY" {
		t.Errorf("AuthorizationType = %v, want 'API_KEY'", conn.AuthorizationType)
	}
	if conn.State != ConnectionStateAuthorized {
		t.Errorf("State = %v, want 'AUTHORIZED'", conn.State)
	}
}

func TestApiDestination(t *testing.T) {
	apiDest := &ApiDestination{
		Name:                         "my-api-dest",
		ARN:                          "arn:aws:events:us-east-1:123456789012:api-destination/my-api-dest",
		Region:                       "us-east-1",
		AccountID:                    "123456789012",
		ConnectionARN:                "arn:aws:events:us-east-1:123456789012:connection/my-connection",
		HttpMethod:                   "POST",
		InvocationEndpoint:           "https://example.com/webhook",
		InvocationRateLimitPerSecond: 100,
		State:                        ApiDestinationStateActive,
	}

	if apiDest.Name != "my-api-dest" {
		t.Errorf("Name = %v, want 'my-api-dest'", apiDest.Name)
	}
	if apiDest.HttpMethod != "POST" {
		t.Errorf("HttpMethod = %v, want 'POST'", apiDest.HttpMethod)
	}
	if apiDest.InvocationRateLimitPerSecond != 100 {
		t.Errorf("InvocationRateLimitPerSecond = %v, want 100", apiDest.InvocationRateLimitPerSecond)
	}
}

func TestEvent(t *testing.T) {
	event := &Event{
		ID:           "event-123",
		Version:      "0",
		DetailType:   "my.event.type",
		Source:       "my.source",
		Account:      "123456789012",
		Time:         time.Now(),
		Region:       "us-east-1",
		Resources:    []string{"resource1", "resource2"},
		Detail:       map[string]interface{}{"key": "value"},
		EventBusName: "my-bus",
	}

	if event.ID != "event-123" {
		t.Errorf("ID = %v, want 'event-123'", event.ID)
	}
	if len(event.Resources) != 2 {
		t.Errorf("len(Resources) = %v, want 2", len(event.Resources))
	}
	if event.Detail["key"] != "value" {
		t.Errorf("Detail[key] = %v, want 'value'", event.Detail["key"])
	}
}

func TestReplay(t *testing.T) {
	replay := &Replay{
		Name:           "my-replay",
		ARN:            "arn:aws:events:us-east-1:123456789012:replay/my-replay",
		Region:         "us-east-1",
		AccountID:      "123456789012",
		State:          ReplayStateRunning,
		EventSourceARN: "arn:aws:events:us-east-1:123456789012:archive/my-archive",
	}

	if replay.Name != "my-replay" {
		t.Errorf("Name = %v, want 'my-replay'", replay.Name)
	}
	if replay.State != ReplayStateRunning {
		t.Errorf("State = %v, want 'RUNNING'", replay.State)
	}
}

func TestTarget(t *testing.T) {
	target := &Target{
		ID:           "my-target",
		RuleName:     "my-rule",
		EventBusName: "my-bus",
		ARN:          "arn:aws:lambda:us-east-1:123456789012:function:my-function",
		Input:        `{"key":"value"}`,
	}

	if target.ID != "my-target" {
		t.Errorf("ID = %v, want 'my-target'", target.ID)
	}
	if target.RuleName != "my-rule" {
		t.Errorf("RuleName = %v, want 'my-rule'", target.RuleName)
	}
	if target.Input != `{"key":"value"}` {
		t.Errorf("Input = %v, want '{\"key\":\"value\"}'", target.Input)
	}
}

func TestEventBusListResult(t *testing.T) {
	result := &EventBusListResult{
		EventBuses: []*EventBus{
			{Name: "bus1"},
			{Name: "bus2"},
		},
		NextToken: "token",
	}

	if len(result.EventBuses) != 2 {
		t.Errorf("len(EventBuses) = %v, want 2", len(result.EventBuses))
	}
	if result.NextToken != "token" {
		t.Errorf("NextToken = %v, want 'token'", result.NextToken)
	}
}

func TestRuleListResult(t *testing.T) {
	result := &RuleListResult{
		Rules: []*Rule{
			{Name: "rule1"},
			{Name: "rule2"},
		},
		NextToken: "token",
	}

	if len(result.Rules) != 2 {
		t.Errorf("len(Rules) = %v, want 2", len(result.Rules))
	}
}
