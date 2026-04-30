package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeTargetTests(ctx context.Context, client *eventbridge.Client, busName, ruleName, targetID string) []TestResult {
	var results []TestResult

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
		if resp.FailedEntryCount > 0 {
			return fmt.Errorf("expected 0 failed entries, got %d", resp.FailedEntryCount)
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
		if resp.Targets == nil || len(resp.Targets) == 0 {
			return fmt.Errorf("expected at least 1 target")
		}
		found := false
		for _, t := range resp.Targets {
			if t.Id != nil && *t.Id == targetID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected target %s in list", targetID)
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
