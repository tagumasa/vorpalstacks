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
			Name:               aws.String(ruleName),
			EventBusName:       aws.String(busName),
			ScheduleExpression: aws.String("rate(5 minutes)"),
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

	// RuleDescription @length(0,512) counts Unicode characters; the shape
	// carries no pattern, so a 256-character CJK description (768 bytes) is
	// legal AWS input and must survive a round trip.
	results = append(results, r.RunTest("events", "PutRule_DescriptionMultibyteAccepted", func() error {
		desc := strings.Repeat("\u65e5", 256)
		if _, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:               aws.String(ruleName),
			EventBusName:       aws.String(busName),
			Description:        aws.String(desc),
			ScheduleExpression: aws.String("rate(5 minutes)"),
		}); err != nil {
			return fmt.Errorf("PutRule with multibyte description: %v", err)
		}
		got, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return fmt.Errorf("DescribeRule: %v", err)
		}
		if aws.ToString(got.Description) != desc {
			return fmt.Errorf("Description mismatch: got %d characters, want %d",
				len([]rune(aws.ToString(got.Description))), len([]rune(desc)))
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule_InvalidScheduleExpression", func() error {
		invalid := []string{
			// week is not a unit accepted by scheduled rules
			"rate(1 week)",
			"rate(2 weeks)",
			// a value of 1 requires a singular unit, values above 1 a plural one
			"rate(1 minutes)",
			"rate(5 minute)",
			// the value must be a positive number
			"rate(0 minutes)",
			// at() is EventBridge Scheduler syntax, not scheduled-rule syntax
			"at(2026-01-01T12:00:00)",
		}
		for _, expr := range invalid {
			_, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
				Name:               aws.String(fmt.Sprintf("BadRule-%d", time.Now().UnixNano())),
				EventBusName:       aws.String(busName),
				ScheduleExpression: aws.String(expr),
			})
			if err == nil {
				return fmt.Errorf("expected error for ScheduleExpression %q", expr)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule_CronLastFriday", func() error {
		// Documented AWS cron examples using the L, W and # day
		// wildcards (EventBridge cron reference / PutRule examples).
		valid := []string{
			"cron(15 10 ? * 6L 2019-2022)", // last Friday of the month
			"cron(0 9 1W * ? *)",           // weekday nearest the 1st
			"cron(0 9 ? * FRI#3 2027)",     // third Friday of the month
			"cron(0 0 L * ? *)",            // last day of the month
		}
		for _, expr := range valid {
			name := fmt.Sprintf("CronRule-%d", time.Now().UnixNano())
			_, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
				Name:               aws.String(name),
				EventBusName:       aws.String(busName),
				ScheduleExpression: aws.String(expr),
			})
			if err != nil {
				return fmt.Errorf("PutRule(%q): %v", expr, err)
			}
			defer client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(name), EventBusName: aws.String(busName)})

			resp, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
				Name:         aws.String(name),
				EventBusName: aws.String(busName),
			})
			if err != nil {
				return fmt.Errorf("DescribeRule(%q): %v", expr, err)
			}
			if resp.ScheduleExpression == nil || *resp.ScheduleExpression != expr {
				return fmt.Errorf("ScheduleExpression mismatch for %q, got %v", expr, resp.ScheduleExpression)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule_CronFiveFieldsRejected", func() error {
		_, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:               aws.String(fmt.Sprintf("BadRule5-%d", time.Now().UnixNano())),
			EventBusName:       aws.String(busName),
			ScheduleExpression: aws.String("cron(0 12 * * ?)"),
		})
		if err == nil {
			return fmt.Errorf("expected error for five-field cron expression")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutRule_CronDomDowBothSpecified", func() error {
		_, err := client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:               aws.String(fmt.Sprintf("BadRuleDD-%d", time.Now().UnixNano())),
			EventBusName:       aws.String(busName),
			ScheduleExpression: aws.String("cron(0 12 15 * FRI 2027)"),
		})
		if err == nil {
			return fmt.Errorf("expected error when day-of-month and day-of-week are both specified")
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
			Name:               aws.String(rdRule),
			EventBusName:       aws.String(rdBus),
			Description:        aws.String("test rule for disable"),
			ScheduleExpression: aws.String("rate(5 minutes)"),
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
				Name:               aws.String(name),
				ScheduleExpression: aws.String("rate(1 hour)"),
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
