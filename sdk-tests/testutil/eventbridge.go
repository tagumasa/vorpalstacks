package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	awsSmithy "github.com/aws/smithy-go"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunEventBridgeTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "events",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := eventbridge.NewFromConfig(cfg)
	ctx := context.Background()

	busName := fmt.Sprintf("TestBus-%d", time.Now().UnixNano())
	ruleName := fmt.Sprintf("TestRule-%d", time.Now().UnixNano())
	targetID := fmt.Sprintf("TestTarget-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("events", "CreateEventBus", func() error {
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(busName),
		})
		return err
	}))

	results = append(results, r.RunTest("events", "DescribeEventBus", func() error {
		resp, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("event bus name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListEventBuses", func() error {
		resp, err := client.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{})
		if err != nil {
			return err
		}
		if resp.EventBuses == nil {
			return fmt.Errorf("event buses list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule", func() error {
		_, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		return err
	}))

	results = append(results, r.RunTest("events", "DescribeRule", func() error {
		resp, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("rule name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListRules", func() error {
		resp, err := client.ListRules(ctx, &eventbridge.ListRulesInput{
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp.Rules == nil {
			return fmt.Errorf("rules list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutTargets", func() error {
		_, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Targets: []types.Target{
				{
					Id:  aws.String(targetID),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFunction", r.region)),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("events", "ListTargetsByRule", func() error {
		resp, err := client.ListTargetsByRule(ctx, &eventbridge.ListTargetsByRuleInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp.Targets == nil {
			return fmt.Errorf("targets list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents", func() error {
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.example.test",
			"detail-type": "TestEvent",
			"detail":      map[string]string{"message": "test"},
		})
		_, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("TestEvent"),
					Detail:       aws.String(string(event)),
					EventBusName: aws.String(busName),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("events", "RemoveTargets", func() error {
		_, err := client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Ids:          []string{targetID},
		})
		return err
	}))

	results = append(results, r.RunTest("events", "DisableRule", func() error {
		_, err := client.DisableRule(ctx, &eventbridge.DisableRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		return err
	}))

	results = append(results, r.RunTest("events", "EnableRule", func() error {
		_, err := client.EnableRule(ctx, &eventbridge.EnableRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		return err
	}))

	results = append(results, r.RunTest("events", "TagResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:000000000000:rule/%s/%s", r.region, busName, ruleName)
		_, err := client.TagResource(ctx, &eventbridge.TagResourceInput{
			ResourceARN: aws.String(ruleARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("events", "ListTagsForResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:000000000000:rule/%s/%s", r.region, busName, ruleName)
		resp, err := client.ListTagsForResource(ctx, &eventbridge.ListTagsForResourceInput{
			ResourceARN: aws.String(ruleARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UntagResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:000000000000:rule/%s/%s", r.region, busName, ruleName)
		_, err := client.UntagResource(ctx, &eventbridge.UntagResourceInput{
			ResourceARN: aws.String(ruleARN),
			TagKeys:     []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("events", "DeleteRule", func() error {
		_, err := client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		return err
	}))

	results = append(results, r.RunTest("events", "DeleteEventBus", func() error {
		_, err := client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
			Name: aws.String(busName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("events", "DescribeEventBus_NonExistent", func() error {
		_, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String("nonexistent-bus-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent event bus")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteEventBus_NonExistent", func() error {
		_, err := client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
			Name: aws.String("nonexistent-bus-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent event bus")
		}
		var apiErr awsSmithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "ResourceNotFoundException" {
			return fmt.Errorf("expected ResourceNotFoundException, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeRule_NonExistent", func() error {
		_, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name: aws.String("nonexistent-rule-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent rule")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteRule_NonExistent", func() error {
		_, err := client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
			Name: aws.String("nonexistent-rule-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent rule")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "CreateEventBus_DuplicateName", func() error {
		dupBus := fmt.Sprintf("DupBus-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(dupBus),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(dupBus)})

		_, err = client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(dupBus),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate event bus name")
		}
		var riu *types.ResourceAlreadyExistsException
		if !errors.As(err, &riu) {
			return fmt.Errorf("expected ResourceAlreadyExistsException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule_DisableAndVerify", func() error {
		rdBus := fmt.Sprintf("RdBus-%d", time.Now().UnixNano())
		rdRule := fmt.Sprintf("RdRule-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(rdBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(rdBus)})

		_, err = client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(rdRule),
			EventBusName: aws.String(rdBus),
			Description:  aws.String("test rule for disable"),
		})
		if err != nil {
			return fmt.Errorf("put rule: %v", err)
		}
		defer client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(rdRule), EventBusName: aws.String(rdBus)})

		_, err = client.DisableRule(ctx, &eventbridge.DisableRuleInput{
			Name:         aws.String(rdRule),
			EventBusName: aws.String(rdBus),
		})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		resp, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(rdRule),
			EventBusName: aws.String(rdBus),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.State != types.RuleStateDisabled {
			return fmt.Errorf("expected state DISABLED, got %v", resp.State)
		}

		_, err = client.EnableRule(ctx, &eventbridge.EnableRuleInput{
			Name:         aws.String(rdRule),
			EventBusName: aws.String(rdBus),
		})
		if err != nil {
			return fmt.Errorf("enable: %v", err)
		}

		resp2, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(rdRule),
			EventBusName: aws.String(rdBus),
		})
		if err != nil {
			return fmt.Errorf("describe 2: %v", err)
		}
		if resp2.State != types.RuleStateEnabled {
			return fmt.Errorf("expected state ENABLED, got %v", resp2.State)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule_WithEventPattern", func() error {
		epBus := fmt.Sprintf("EpBus-%d", time.Now().UnixNano())
		epRule := fmt.Sprintf("EpRule-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(epBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(epBus)})

		pattern := map[string]interface{}{
			"source":      []string{"com.example.test"},
			"detail-type": []string{"OrderCreated"},
		}
		patternJSON, _ := json.Marshal(pattern)

		_, err = client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(epRule),
			EventBusName: aws.String(epBus),
			EventPattern: aws.String(string(patternJSON)),
		})
		if err != nil {
			return fmt.Errorf("put rule: %v", err)
		}
		defer client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(epRule), EventBusName: aws.String(epBus)})

		resp, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(epRule),
			EventBusName: aws.String(epBus),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.EventPattern == nil {
			return fmt.Errorf("event pattern is nil")
		}

		var gotPattern map[string]interface{}
		if err := json.Unmarshal([]byte(*resp.EventPattern), &gotPattern); err != nil {
			return fmt.Errorf("unmarshal pattern: %v", err)
		}
		gotSource, ok := gotPattern["source"].([]interface{})
		if !ok || len(gotSource) != 1 || gotSource[0] != "com.example.test" {
			return fmt.Errorf("source mismatch in pattern, got %v", gotSource)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_DefaultBus", func() error {
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.test.default",
			"detail-type": "DefaultBusEvent",
			"detail":      map[string]string{"key": "value"},
		})
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:     aws.String("com.test.default"),
					DetailType: aws.String("DefaultBusEvent"),
					Detail:     aws.String(string(event)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events: %v", err)
		}
		if resp.FailedEntryCount != 0 {
			return fmt.Errorf("expected 0 failed entries, got %d", resp.FailedEntryCount)
		}
		if len(resp.Entries) != 1 {
			return fmt.Errorf("expected 1 entry result, got %d", len(resp.Entries))
		}
		if resp.Entries[0].EventId == nil || *resp.Entries[0].EventId == "" {
			return fmt.Errorf("expected non-empty event ID")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutTargets_RemoveTargets_Verify", func() error {
		trBus := fmt.Sprintf("TrBus-%d", time.Now().UnixNano())
		trRule := fmt.Sprintf("TrRule-%d", time.Now().UnixNano())
		trTargetID := fmt.Sprintf("TrTarget-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(trBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(trBus)})

		_, err = client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(trRule),
			EventBusName: aws.String(trBus),
		})
		if err != nil {
			return fmt.Errorf("put rule: %v", err)
		}
		defer client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(trRule), EventBusName: aws.String(trBus)})

		targetARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TargetFunc", r.region)
		_, err = client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(trRule),
			EventBusName: aws.String(trBus),
			Targets: []types.Target{
				{
					Id:    aws.String(trTargetID),
					Arn:   aws.String(targetARN),
					Input: aws.String(`{"action": "test"}`),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets: %v", err)
		}

		listResp, err := client.ListTargetsByRule(ctx, &eventbridge.ListTargetsByRuleInput{
			Rule:         aws.String(trRule),
			EventBusName: aws.String(trBus),
		})
		if err != nil {
			return fmt.Errorf("list targets: %v", err)
		}
		if len(listResp.Targets) != 1 {
			return fmt.Errorf("expected 1 target, got %d", len(listResp.Targets))
		}
		if listResp.Targets[0].Arn == nil || *listResp.Targets[0].Arn != targetARN {
			return fmt.Errorf("target ARN mismatch, got %v", listResp.Targets[0].Arn)
		}
		if listResp.Targets[0].Input == nil || *listResp.Targets[0].Input != `{"action": "test"}` {
			return fmt.Errorf("target input mismatch, got %v", listResp.Targets[0].Input)
		}

		_, err = client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule:         aws.String(trRule),
			EventBusName: aws.String(trBus),
			Ids:          []string{trTargetID},
		})
		if err != nil {
			return fmt.Errorf("remove targets: %v", err)
		}

		listResp2, err := client.ListTargetsByRule(ctx, &eventbridge.ListTargetsByRuleInput{
			Rule:         aws.String(trRule),
			EventBusName: aws.String(trBus),
		})
		if err != nil {
			return fmt.Errorf("list targets after remove: %v", err)
		}
		if len(listResp2.Targets) != 0 {
			return fmt.Errorf("expected 0 targets after removal, got %d", len(listResp2.Targets))
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteRule_WithTargetsFails", func() error {
		dtBus := fmt.Sprintf("DtBus-%d", time.Now().UnixNano())
		dtRule := fmt.Sprintf("DtRule-%d", time.Now().UnixNano())
		dtTarget := fmt.Sprintf("DtTarget-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(dtBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(dtBus)})

		_, err = client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(dtRule),
			EventBusName: aws.String(dtBus),
		})
		if err != nil {
			return fmt.Errorf("put rule: %v", err)
		}
		defer func() {
			client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
				Rule: aws.String(dtRule), EventBusName: aws.String(dtBus), Ids: []string{dtTarget},
			})
			client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(dtRule), EventBusName: aws.String(dtBus)})
		}()

		_, err = client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(dtRule),
			EventBusName: aws.String(dtBus),
			Targets: []types.Target{
				{
					Id:  aws.String(dtTarget),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:F", r.region)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets: %v", err)
		}

		_, err = client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
			Name:         aws.String(dtRule),
			EventBusName: aws.String(dtBus),
		})
		if err == nil {
			return fmt.Errorf("expected error when deleting rule with targets")
		}
		return nil
	}))

	return results
}
