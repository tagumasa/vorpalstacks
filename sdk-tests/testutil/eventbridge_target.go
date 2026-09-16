package testutil

import (
	"context"
	"fmt"
	"strings"
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
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:TestFunction", r.region, r.accountID)),
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
		cleanupBus, err := createEventBridgeTestBus(ctx, client, trBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		cleanupRule, err := createEventBridgeTestRule(ctx, client, trBus, trRule)
		if err != nil {
			return err
		}
		defer cleanupRule()

		targetARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:TargetFunc", r.region, r.accountID)
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
		cleanupBus, err := createEventBridgeTestBus(ctx, client, lrntBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		cleanupRule, err := createEventBridgeTestRule(ctx, client, lrntBus, lrntRule)
		if err != nil {
			return err
		}
		defer cleanupRule()

		targetARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:ListRulesFn", r.region, r.accountID)
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
		cleanupBus, err := createEventBridgeTestBus(ctx, client, dtBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		cleanupRule, err := createEventBridgeTestRule(ctx, client, dtBus, dtRule)
		if err != nil {
			return err
		}
		defer cleanupRule()

		// Targets must be removed before the deferred rule deletion can
		// succeed; this defer runs ahead of the helper cleanups.
		defer client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(dtRule), EventBusName: aws.String(dtBus), Ids: []string{dtTarget},
		})

		_, err = client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(dtRule),
			EventBusName: aws.String(dtBus),
			Targets: []types.Target{
				{
					Id:  aws.String(dtTarget),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:F", r.region, r.accountID)),
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

	// The accept/deliver matrix: target ARNs whose service or resource form
	// has no delivery path on this platform are rejected per-entry with a
	// reason, while the AppSync endpoint form is accepted.
	results = append(results, r.RunTest("events", "PutTargets_AcceptDeliverMatrix", func() error {
		mtBus := fmt.Sprintf("MtBus-%d", time.Now().UnixNano())
		mtRule := fmt.Sprintf("MtRule-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, mtBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		cleanupRule, err := createEventBridgeTestRule(ctx, client, mtBus, mtRule)
		if err != nil {
			return err
		}
		defer cleanupRule()

		resp, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(mtRule),
			EventBusName: aws.String(mtBus),
			Targets: []types.Target{
				{Id: aws.String("t-ecs"), Arn: aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/none", r.region, r.accountID))},
				{Id: aws.String("t-ssm"), Arn: aws.String(fmt.Sprintf("arn:aws:ssm:%s:%s:document/AWS-RunShellScript", r.region, r.accountID))},
				{Id: aws.String("t-rule"), Arn: aws.String(fmt.Sprintf("arn:aws:events:%s:%s:rule/%s/%s", r.region, r.accountID, mtBus, mtRule))},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets: %v", err)
		}
		if resp.FailedEntryCount != 3 {
			return fmt.Errorf("expected all three unsupported target forms to fail, got %d", resp.FailedEntryCount)
		}
		failed := map[string]bool{}
		for _, e := range resp.FailedEntries {
			if e.ErrorCode == nil || *e.ErrorCode != "ValidationException" {
				return fmt.Errorf("target %s: expected ValidationException, got %+v", aws.ToString(e.TargetId), e.ErrorCode)
			}
			if e.ErrorMessage == nil || *e.ErrorMessage == "" {
				return fmt.Errorf("target %s: rejection must carry a reason", aws.ToString(e.TargetId))
			}
			failed[aws.ToString(e.TargetId)] = true
		}
		for _, id := range []string{"t-ecs", "t-ssm", "t-rule"} {
			if !failed[id] {
				return fmt.Errorf("target %s missing from the failed entries", id)
			}
		}

		// The events service's api-destination resource form is accepted
		// (delivery invokes the destination's endpoint); only its
		// rule/archive resource forms stay rejected.
		apiDest, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(mtRule),
			EventBusName: aws.String(mtBus),
			Targets: []types.Target{
				{Id: aws.String("t-apidest"), Arn: aws.String(fmt.Sprintf("arn:aws:events:%s:%s:api-destination/none/id-1", r.region, r.accountID))},
			},
		})
		if err != nil {
			return fmt.Errorf("put api-destination target: %v", err)
		}
		if apiDest.FailedEntryCount != 0 {
			return fmt.Errorf("the api-destination form must be accepted (delivery, not acceptance, owns existence): %v", apiDest.FailedEntries)
		}
		defer client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(mtRule), EventBusName: aws.String(mtBus), Ids: []string{"t-apidest"},
		})

		appSync, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(mtRule),
			EventBusName: aws.String(mtBus),
			Targets: []types.Target{
				{
					Id:                aws.String("t-appsync"),
					Arn:               aws.String(fmt.Sprintf("arn:aws:appsync:%s:%s:apis/nonexistent/endpoints/GRAPHQL", r.region, r.accountID)),
					AppSyncParameters: &types.AppSyncParameters{GraphQLOperation: aws.String("mutation { pushEvent { id } }")},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put appsync target: %v", err)
		}
		if appSync.FailedEntryCount != 0 {
			return fmt.Errorf("the AppSync endpoint form must be accepted (delivery, not acceptance, owns existence): %v", appSync.FailedEntries)
		}
		defer client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(mtRule), EventBusName: aws.String(mtBus), Ids: []string{"t-appsync"},
		})

		// The RetryPolicy age member distinguishes "explicitly supplied"
		// from "omitted": an explicitly supplied MaximumEventAgeInSeconds
		// of 0 sits outside the modelled 60-86400 window and fails its
		// entry, while an omitted member keeps the deployment default.
		zeroAge, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(mtRule),
			EventBusName: aws.String(mtBus),
			Targets: []types.Target{{
				Id:          aws.String("t-retry-zero"),
				Arn:         aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:retry-zero-queue", r.region, r.accountID)),
				RetryPolicy: &types.RetryPolicy{MaximumEventAgeInSeconds: aws.Int32(0)},
			}},
		})
		if err != nil {
			return fmt.Errorf("put zero-age retry target: %v", err)
		}
		if zeroAge.FailedEntryCount != 1 {
			return fmt.Errorf("explicit MaximumEventAgeInSeconds 0 must fail its entry, got %+v", zeroAge.FailedEntries)
		}
		if code := aws.ToString(zeroAge.FailedEntries[0].ErrorCode); code != "ValidationException" {
			return fmt.Errorf("zero-age retry target: expected ValidationException, got %s", code)
		}
		return nil
	}))

	// The dead-letter queue contract: SQS standard queues only, in the
	// rule's region; anything else is a per-entry validation failure.
	results = append(results, r.RunTest("events", "PutTargets_DLQValidation", func() error {
		dqBus := fmt.Sprintf("DqBus-%d", time.Now().UnixNano())
		dqRule := fmt.Sprintf("DqRule-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, dqBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		cleanupRule, err := createEventBridgeTestRule(ctx, client, dqBus, dqRule)
		if err != nil {
			return err
		}
		defer cleanupRule()

		resp, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(dqRule),
			EventBusName: aws.String(dqBus),
			Targets: []types.Target{
				{
					Id:  aws.String("t-lambda-dlq"),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:F", r.region, r.accountID)),
					DeadLetterConfig: &types.DeadLetterConfig{
						Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:not-a-dlq", r.region, r.accountID)),
					},
				},
				{
					Id:  aws.String("t-fifo-dlq"),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:F", r.region, r.accountID)),
					DeadLetterConfig: &types.DeadLetterConfig{
						Arn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:events-dlq.fifo", r.region, r.accountID)),
					},
				},
				{
					Id:  aws.String("t-cross-region-dlq"),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:F", r.region, r.accountID)),
					DeadLetterConfig: &types.DeadLetterConfig{
						Arn: aws.String(fmt.Sprintf("arn:aws:sqs:eu-central-1:%s:events-dlq", r.accountID)),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets: %v", err)
		}
		if resp.FailedEntryCount != 3 {
			return fmt.Errorf("expected all three illegal DLQ configurations to fail, got %d", resp.FailedEntryCount)
		}
		for _, e := range resp.FailedEntries {
			if e.ErrorCode == nil || *e.ErrorCode != "ValidationException" {
				return fmt.Errorf("target %s: expected ValidationException, got %+v", aws.ToString(e.TargetId), e.ErrorCode)
			}
		}

		ok, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(dqRule),
			EventBusName: aws.String(dqBus),
			Targets: []types.Target{
				{
					Id:  aws.String("t-sqs-dlq"),
					Arn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:F", r.region, r.accountID)),
					DeadLetterConfig: &types.DeadLetterConfig{
						Arn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:events-dlq", r.region, r.accountID)),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets with sqs dlq: %v", err)
		}
		if ok.FailedEntryCount != 0 {
			return fmt.Errorf("an SQS DLQ ARN in the rule's region must be accepted: %v", ok.FailedEntries)
		}
		defer client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(dqRule), EventBusName: aws.String(dqBus), Ids: []string{"t-sqs-dlq"},
		})
		return nil
	}))

	// The input configuration carries an acceptance contract: the three
	// members are mutually exclusive, Input is valid JSON text within the
	// documented length, and the transformer's template and paths map obey
	// their bounds; a well-formed transformer is accepted.
	results = append(results, r.RunTest("events", "PutTargets_InputConfigurationValidation", func() error {
		ivBus := fmt.Sprintf("IvBus-%d", time.Now().UnixNano())
		ivRule := fmt.Sprintf("IvRule-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, ivBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		cleanupRule, err := createEventBridgeTestRule(ctx, client, ivBus, ivRule)
		if err != nil {
			return err
		}
		defer cleanupRule()

		lambdaARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:IvFunc", r.region, r.accountID)
		// The SDK rejects a missing InputTemplate client-side (required
		// member), so that path is pinned by unit tests; the remaining
		// three misconfigurations reach the server.
		resp, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(ivRule),
			EventBusName: aws.String(ivBus),
			Targets: []types.Target{
				{Id: aws.String("t-both"), Arn: aws.String(lambdaARN), Input: aws.String(`{"k":1}`), InputPath: aws.String("$.detail")},
				{Id: aws.String("t-badjson"), Arn: aws.String(lambdaARN), Input: aws.String(`{not json`)},
				{Id: aws.String("t-longpath"), Arn: aws.String(lambdaARN), InputPath: aws.String("$." + strings.Repeat("a", 257))},
			},
		})
		if err != nil {
			return fmt.Errorf("put targets: %v", err)
		}
		if resp.FailedEntryCount != 3 {
			return fmt.Errorf("expected all three misconfigured input configurations to fail, got %d", resp.FailedEntryCount)
		}
		for _, e := range resp.FailedEntries {
			if e.ErrorCode == nil || *e.ErrorCode != "ValidationException" {
				return fmt.Errorf("target %s: expected ValidationException, got %+v", aws.ToString(e.TargetId), e.ErrorCode)
			}
			if e.ErrorMessage == nil || *e.ErrorMessage == "" {
				return fmt.Errorf("target %s: rejection must carry a reason", aws.ToString(e.TargetId))
			}
		}

		valid, err := client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(ivRule),
			EventBusName: aws.String(ivBus),
			Targets: []types.Target{
				{
					Id:  aws.String("t-transformer"),
					Arn: aws.String(lambdaARN),
					InputTransformer: &types.InputTransformer{
						InputPathsMap: map[string]string{"state": "$.detail.state"},
						InputTemplate: aws.String(`{"state": <state>}`),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put transformer target: %v", err)
		}
		if valid.FailedEntryCount != 0 {
			return fmt.Errorf("a well-formed transformer must be accepted, got %v", valid.FailedEntries)
		}
		defer client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule: aws.String(ivRule), EventBusName: aws.String(ivBus), Ids: []string{"t-transformer"},
		})
		return nil
	}))

	return results
}
