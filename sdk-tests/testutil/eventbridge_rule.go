package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeRuleTests(ctx context.Context, client *eventbridge.Client, busName, ruleName string) []TestResult {
	var results []TestResult

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
		if resp.RuleArn == nil || *resp.RuleArn == "" {
			return fmt.Errorf("rule ARN is nil or empty")
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
		if resp.Name == nil || *resp.Name != ruleName {
			return fmt.Errorf("rule name mismatch, got %v", resp.Name)
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("rule ARN is nil or empty")
		}
		if resp.State != types.RuleStateEnabled {
			return fmt.Errorf("expected ENABLED state, got %v", resp.State)
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
		found := false
		for _, rule := range resp.Rules {
			if rule.Name != nil && *rule.Name == ruleName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected rule %s in list", ruleName)
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
