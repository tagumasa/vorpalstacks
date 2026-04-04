package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
		resp, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Targets: []types.Target{
				{
					Id:  aws.String(targetID),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFunction", r.region)),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("TestEvent"),
					Detail:       aws.String(string(event)),
					EventBusName: aws.String(busName),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "RemoveTargets", func() error {
		resp, err := client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Ids:          []string{targetID},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DisableRule", func() error {
		resp, err := client.DisableRule(ctx, &eventbridge.DisableRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "EnableRule", func() error {
		resp, err := client.EnableRule(ctx, &eventbridge.EnableRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "TagResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:000000000000:rule/%s/%s", r.region, busName, ruleName)
		resp, err := client.TagResource(ctx, &eventbridge.TagResourceInput{
			ResourceARN: aws.String(ruleARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.UntagResource(ctx, &eventbridge.UntagResourceInput{
			ResourceARN: aws.String(ruleARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteRule", func() error {
		resp, err := client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteEventBus", func() error {
		resp, err := client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
			Name: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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

	results = append(results, r.RunTest("events", "UpdateEventBus", func() error {
		ueBus := fmt.Sprintf("UeBus-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(ueBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(ueBus)})

		resp, err := client.UpdateEventBus(ctx, &eventbridge.UpdateEventBusInput{
			Name:        aws.String(ueBus),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return fmt.Errorf("update event bus: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Description == nil || *resp.Description != "updated description" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateEventBus_VerifyDescription", func() error {
		uvBus := fmt.Sprintf("UvBus-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name:        aws.String(uvBus),
			Description: aws.String("original"),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(uvBus)})

		_, err = client.UpdateEventBus(ctx, &eventbridge.UpdateEventBusInput{
			Name:        aws.String(uvBus),
			Description: aws.String("updated"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		desc, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(uvBus),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if desc.Description == nil || *desc.Description != "updated" {
			return fmt.Errorf("description not updated, got %v", desc.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "CreateArchive", func() error {
		caBus := fmt.Sprintf("CaBus-%d", time.Now().UnixNano())
		caArchive := fmt.Sprintf("CaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(caBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(caBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, caBus)
		resp, err := client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(caArchive),
			EventSourceArn: aws.String(busARN),
			Description:    aws.String("test archive"),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ArchiveArn == nil || *resp.ArchiveArn == "" {
			return fmt.Errorf("archive ARN is nil or empty")
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(caArchive)})
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeArchive", func() error {
		daBus := fmt.Sprintf("DaBus-%d", time.Now().UnixNano())
		daArchive := fmt.Sprintf("DaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(daBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(daBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, daBus)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(daArchive),
			EventSourceArn: aws.String(busARN),
			Description:    aws.String("test archive for describe"),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(daArchive)})

		resp, err := client.DescribeArchive(ctx, &eventbridge.DescribeArchiveInput{
			ArchiveName: aws.String(daArchive),
		})
		if err != nil {
			return fmt.Errorf("describe archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ArchiveName == nil || *resp.ArchiveName != daArchive {
			return fmt.Errorf("archive name mismatch, got %v", resp.ArchiveName)
		}
		if resp.EventSourceArn == nil || *resp.EventSourceArn != busARN {
			return fmt.Errorf("event source ARN mismatch, got %v", resp.EventSourceArn)
		}
		if resp.Description == nil || *resp.Description != "test archive for describe" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeArchive_NonExistent", func() error {
		_, err := client.DescribeArchive(ctx, &eventbridge.DescribeArchiveInput{
			ArchiveName: aws.String("nonexistent-archive-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent archive")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteArchive_NonExistent", func() error {
		_, err := client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{
			ArchiveName: aws.String("nonexistent-archive-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent archive")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListArchives", func() error {
		resp, err := client.ListArchives(ctx, &eventbridge.ListArchivesInput{})
		if err != nil {
			return fmt.Errorf("list archives: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Archives == nil {
			return fmt.Errorf("archives list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateArchive", func() error {
		uaBus := fmt.Sprintf("UaBus-%d", time.Now().UnixNano())
		uaArchive := fmt.Sprintf("UaArchive-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(uaBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(uaBus)})

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, uaBus)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(uaArchive),
			EventSourceArn: aws.String(busARN),
			Description:    aws.String("original description"),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}
		defer client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(uaArchive)})

		resp, err := client.UpdateArchive(ctx, &eventbridge.UpdateArchiveInput{
			ArchiveName: aws.String(uaArchive),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return fmt.Errorf("update archive: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ArchiveArn == nil || *resp.ArchiveArn == "" {
			return fmt.Errorf("archive ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "CreateConnection", func() error {
		ccName := fmt.Sprintf("CcConn-%d", time.Now().UnixNano())
		resp, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(ccName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("testuser"),
					Password: aws.String("testpass"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn == "" {
			return fmt.Errorf("connection ARN is nil or empty")
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(ccName)})
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeConnection", func() error {
		dcName := fmt.Sprintf("DcConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(dcName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("testuser"),
					Password: aws.String("testpass"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(dcName)})

		resp, err := client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
			Name: aws.String(dcName),
		})
		if err != nil {
			return fmt.Errorf("describe connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Name == nil || *resp.Name != dcName {
			return fmt.Errorf("connection name mismatch, got %v", resp.Name)
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn == "" {
			return fmt.Errorf("connection ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeConnection_NonExistent", func() error {
		_, err := client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
			Name: aws.String("nonexistent-conn-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent connection")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteConnection_NonExistent", func() error {
		_, err := client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{
			Name: aws.String("nonexistent-conn-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent connection")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListConnections", func() error {
		resp, err := client.ListConnections(ctx, &eventbridge.ListConnectionsInput{})
		if err != nil {
			return fmt.Errorf("list connections: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Connections == nil {
			return fmt.Errorf("connections list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateConnection", func() error {
		ucName := fmt.Sprintf("UcConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(ucName),
			AuthorizationType: types.ConnectionAuthorizationTypeApiKey,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				ApiKeyAuthParameters: &types.CreateConnectionApiKeyAuthRequestParameters{
					ApiKeyName:  aws.String("key"),
					ApiKeyValue: aws.String("value"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(ucName)})

		resp, err := client.UpdateConnection(ctx, &eventbridge.UpdateConnectionInput{
			Name:        aws.String(ucName),
			Description: aws.String("updated connection description"),
		})
		if err != nil {
			return fmt.Errorf("update connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn == "" {
			return fmt.Errorf("connection ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "CreateApiDestination", func() error {
		cadName := fmt.Sprintf("CadDest-%d", time.Now().UnixNano())
		connName := fmt.Sprintf("CadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(connName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(connName)})

		resp, err := client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(cadName),
			ConnectionArn:      aws.String(fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, connName)),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/webhook"),
			Description:        aws.String("test api destination"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ApiDestinationArn == nil || *resp.ApiDestinationArn == "" {
			return fmt.Errorf("api destination ARN is nil or empty")
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(cadName)})
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeApiDestination", func() error {
		dadName := fmt.Sprintf("DadDest-%d", time.Now().UnixNano())
		dadConn := fmt.Sprintf("DadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(dadConn),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(dadConn)})

		connARN := fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, dadConn)
		_, err = client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(dadName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/webhook"),
			Description:        aws.String("test api destination for describe"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(dadName)})

		resp, err := client.DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
			Name: aws.String(dadName),
		})
		if err != nil {
			return fmt.Errorf("describe api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Name == nil || *resp.Name != dadName {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn != connARN {
			return fmt.Errorf("connection ARN mismatch, got %v", resp.ConnectionArn)
		}
		if resp.HttpMethod != types.ApiDestinationHttpMethodPost {
			return fmt.Errorf("http method mismatch, got %v", resp.HttpMethod)
		}
		if resp.InvocationEndpoint == nil || *resp.InvocationEndpoint != "https://example.com/webhook" {
			return fmt.Errorf("invocation endpoint mismatch, got %v", resp.InvocationEndpoint)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeApiDestination_NonExistent", func() error {
		_, err := client.DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
			Name: aws.String("nonexistent-apidest-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent api destination")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteApiDestination_NonExistent", func() error {
		_, err := client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{
			Name: aws.String("nonexistent-apidest-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent api destination")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListApiDestinations", func() error {
		resp, err := client.ListApiDestinations(ctx, &eventbridge.ListApiDestinationsInput{})
		if err != nil {
			return fmt.Errorf("list api destinations: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ApiDestinations == nil {
			return fmt.Errorf("api destinations list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateApiDestination", func() error {
		uadName := fmt.Sprintf("UadDest-%d", time.Now().UnixNano())
		uadConn := fmt.Sprintf("UadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(uadConn),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(uadConn)})

		connARN := fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, uadConn)
		_, err = client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(uadName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/original"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(uadName)})

		resp, err := client.UpdateApiDestination(ctx, &eventbridge.UpdateApiDestinationInput{
			Name:               aws.String(uadName),
			Description:        aws.String("updated description"),
			InvocationEndpoint: aws.String("https://example.com/updated"),
		})
		if err != nil {
			return fmt.Errorf("update api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ApiDestinationArn == nil || *resp.ApiDestinationArn == "" {
			return fmt.Errorf("api destination ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "TestEventPattern_Match", func() error {
		pattern, _ := json.Marshal(map[string]interface{}{
			"source": []string{"com.example.custom"},
		})
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.example.custom",
			"detail-type": "TestEvent",
		})
		resp, err := client.TestEventPattern(ctx, &eventbridge.TestEventPatternInput{
			EventPattern: aws.String(string(pattern)),
			Event:        aws.String(string(event)),
		})
		if err != nil {
			return fmt.Errorf("test event pattern: %v", err)
		}
		if !resp.Result {
			return fmt.Errorf("expected pattern to match, got false")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "TestEventPattern_NoMatch", func() error {
		pattern, _ := json.Marshal(map[string]interface{}{
			"source": []string{"com.example.other"},
		})
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.example.custom",
			"detail-type": "TestEvent",
		})
		resp, err := client.TestEventPattern(ctx, &eventbridge.TestEventPatternInput{
			EventPattern: aws.String(string(pattern)),
			Event:        aws.String(string(event)),
		})
		if err != nil {
			return fmt.Errorf("test event pattern: %v", err)
		}
		if resp.Result {
			return fmt.Errorf("expected pattern not to match, got true")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "StartReplay_DescribeReplay", func() error {
		srBus := fmt.Sprintf("SrBus-%d", time.Now().UnixNano())
		srArchive := fmt.Sprintf("SrArchive-%d", time.Now().UnixNano())
		srReplay := fmt.Sprintf("SrReplay-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(srBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}

		busARN := fmt.Sprintf("arn:aws:events:%s:000000000000:event-bus/%s", r.region, srBus)
		archiveARN := fmt.Sprintf("arn:aws:events:%s:000000000000:archive/%s", r.region, srArchive)
		_, err = client.CreateArchive(ctx, &eventbridge.CreateArchiveInput{
			ArchiveName:    aws.String(srArchive),
			EventSourceArn: aws.String(busARN),
		})
		if err != nil {
			return fmt.Errorf("create archive: %v", err)
		}

		now := time.Now()
		startResp, err := client.StartReplay(ctx, &eventbridge.StartReplayInput{
			ReplayName:     aws.String(srReplay),
			EventSourceArn: aws.String(archiveARN),
			Destination: &types.ReplayDestination{
				Arn: aws.String(busARN),
			},
			EventStartTime: aws.Time(now.Add(-1 * time.Hour)),
			EventEndTime:   aws.Time(now),
		})
		if err != nil {
			return fmt.Errorf("start replay: %v", err)
		}
		if startResp == nil {
			return fmt.Errorf("start replay response is nil")
		}
		if startResp.ReplayArn == nil || *startResp.ReplayArn == "" {
			return fmt.Errorf("replay ARN is nil or empty")
		}

		describeResp, err := client.DescribeReplay(ctx, &eventbridge.DescribeReplayInput{
			ReplayName: aws.String(srReplay),
		})
		if err != nil {
			return fmt.Errorf("describe replay: %v", err)
		}
		if describeResp == nil {
			return fmt.Errorf("describe replay response is nil")
		}
		if describeResp.ReplayName == nil || *describeResp.ReplayName != srReplay {
			return fmt.Errorf("replay name mismatch, got %v", describeResp.ReplayName)
		}
		if describeResp.EventSourceArn == nil || *describeResp.EventSourceArn != archiveARN {
			return fmt.Errorf("event source ARN mismatch, got %v", describeResp.EventSourceArn)
		}
		client.CancelReplay(ctx, &eventbridge.CancelReplayInput{ReplayName: aws.String(srReplay)})
		client.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{ArchiveName: aws.String(srArchive)})
		client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(srBus)})
		return nil
	}))

	results = append(results, r.RunTest("events", "CancelReplay_NonExistent", func() error {
		_, err := client.CancelReplay(ctx, &eventbridge.CancelReplayInput{
			ReplayName: aws.String("nonexistent-replay-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent replay")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListReplays", func() error {
		resp, err := client.ListReplays(ctx, &eventbridge.ListReplaysInput{})
		if err != nil {
			return fmt.Errorf("list replays: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Replays == nil {
			return fmt.Errorf("replays list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListRuleNamesByTarget", func() error {
		lrntBus := fmt.Sprintf("LrntBus-%d", time.Now().UnixNano())
		lrntRule := fmt.Sprintf("LrntRule-%d", time.Now().UnixNano())
		lrntTargetID := fmt.Sprintf("LrntTarget-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(lrntBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(lrntBus)})

		_, err = client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(lrntRule),
			EventBusName: aws.String(lrntBus),
		})
		if err != nil {
			return fmt.Errorf("put rule: %v", err)
		}
		defer client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(lrntRule), EventBusName: aws.String(lrntBus)})

		targetARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:ListRulesFn", r.region)
		_, err = client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(lrntRule),
			EventBusName: aws.String(lrntBus),
			Targets: []types.Target{
				{
					Id:  aws.String(lrntTargetID),
					Arn: aws.String(targetARN),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets: %v", err)
		}
		defer client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(lrntRule), EventBusName: aws.String(lrntBus), Ids: []string{lrntTargetID},
		})

		resp, err := client.ListRuleNamesByTarget(ctx, &eventbridge.ListRuleNamesByTargetInput{
			TargetArn:    aws.String(targetARN),
			EventBusName: aws.String(lrntBus),
		})
		if err != nil {
			return fmt.Errorf("list rule names by target: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.RuleNames == nil {
			return fmt.Errorf("rule names list is nil")
		}
		found := false
		for _, name := range resp.RuleNames {
			if name == lrntRule {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected rule %s in list, got %v", lrntRule, resp.RuleNames)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListEventBuses_NamePrefix", func() error {
		lnpBus := fmt.Sprintf("LnpPrefixBus-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(lnpBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(lnpBus)})

		resp, err := client.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{
			NamePrefix: aws.String("LnpPrefixBus"),
		})
		if err != nil {
			return fmt.Errorf("list event buses with prefix: %v", err)
		}
		if resp == nil || resp.EventBuses == nil {
			return fmt.Errorf("response or event buses is nil")
		}
		found := false
		for _, bus := range resp.EventBuses {
			if bus.Name != nil && *bus.Name == lnpBus {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected bus %s in filtered list", lnpBus)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_MultipleEntries", func() error {
		event1, _ := json.Marshal(map[string]interface{}{
			"source": "com.test.multi", "detail-type": "Event1",
			"detail": map[string]string{"id": "1"},
		})
		event2, _ := json.Marshal(map[string]interface{}{
			"source": "com.test.multi", "detail-type": "Event2",
			"detail": map[string]string{"id": "2"},
		})
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:     aws.String("com.test.multi"),
					DetailType: aws.String("Event1"),
					Detail:     aws.String(string(event1)),
				},
				{
					Source:     aws.String("com.test.multi"),
					DetailType: aws.String("Event2"),
					Detail:     aws.String(string(event2)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events: %v", err)
		}
		if resp.FailedEntryCount != 0 {
			return fmt.Errorf("expected 0 failed entries, got %d", resp.FailedEntryCount)
		}
		if len(resp.Entries) != 2 {
			return fmt.Errorf("expected 2 entry results, got %d", len(resp.Entries))
		}
		return nil
	}))

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("events", "ListRules_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgRules []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagRule-%s-%d", pgTs, i)
			_, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
				Name: aws.String(name),
			})
			if err != nil {
				return fmt.Errorf("put rule %s: %v", name, err)
			}
			pgRules = append(pgRules, name)
		}

		var allRules []string
		var nextToken *string
		for {
			resp, err := client.ListRules(ctx, &eventbridge.ListRulesInput{
				NextToken: nextToken,
				Limit:     aws.Int32(2),
			})
			if err != nil {
				for _, name := range pgRules {
					client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(name)})
				}
				return fmt.Errorf("list rules page: %v", err)
			}
			for _, r := range resp.Rules {
				if strings.HasPrefix(aws.ToString(r.Name), "PagRule-"+pgTs) {
					allRules = append(allRules, aws.ToString(r.Name))
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, name := range pgRules {
			client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(name)})
		}
		if len(allRules) != 5 {
			return fmt.Errorf("expected 5 paginated rules, got %d", len(allRules))
		}
		return nil
	}))

	return results
}
