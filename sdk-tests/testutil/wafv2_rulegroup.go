package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2RuleGroupTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	rgName := tc.uniqueName("test-rulegroup")
	var rgID, rgLockToken string

	results = append(results, r.RunTest("wafv2", "CreateRuleGroup", func() error {
		resp, err := tc.client.CreateRuleGroup(tc.ctx, &wafv2.CreateRuleGroupInput{
			Name:     aws.String(rgName),
			Scope:    tc.scope,
			Capacity: aws.Int64(100),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("test-metric"),
			},
		})
		if err != nil {
			return err
		}
		if resp.Summary == nil {
			return fmt.Errorf("Summary is nil")
		}
		if aws.ToString(resp.Summary.Id) == "" {
			return fmt.Errorf("Summary.Id is empty")
		}
		if aws.ToString(resp.Summary.LockToken) == "" {
			return fmt.Errorf("Summary.LockToken is empty")
		}
		if aws.ToString(resp.Summary.Name) != rgName {
			return fmt.Errorf("Summary.Name mismatch: expected %s, got %s", rgName, aws.ToString(resp.Summary.Name))
		}
		rgID = aws.ToString(resp.Summary.Id)
		rgLockToken = aws.ToString(resp.Summary.LockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetRuleGroup", func() error {
		resp, err := tc.client.GetRuleGroup(tc.ctx, &wafv2.GetRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope, Id: aws.String(rgID),
		})
		if err != nil {
			return err
		}
		if resp.RuleGroup == nil {
			return fmt.Errorf("RuleGroup is nil")
		}
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		if resp.RuleGroup.Capacity == nil || *resp.RuleGroup.Capacity != 100 {
			return fmt.Errorf("expected capacity 100, got %v", resp.RuleGroup.Capacity)
		}
		if aws.ToString(resp.RuleGroup.Name) != rgName {
			return fmt.Errorf("Name mismatch: expected %s, got %s", rgName, aws.ToString(resp.RuleGroup.Name))
		}
		if aws.ToString(resp.RuleGroup.Id) != rgID {
			return fmt.Errorf("Id mismatch")
		}
		if resp.RuleGroup.VisibilityConfig == nil {
			return fmt.Errorf("VisibilityConfig is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListRuleGroups", func() error {
		resp, err := tc.client.ListRuleGroups(tc.ctx, &wafv2.ListRuleGroupsInput{
			Scope: tc.scope,
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.RuleGroups == nil {
			return fmt.Errorf("RuleGroups list is nil")
		}
		found := false
		for _, s := range resp.RuleGroups {
			if s.Id != nil && *s.Id == rgID {
				found = true
				if aws.ToString(s.Name) != rgName {
					return fmt.Errorf("RuleGroup name mismatch in list")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("RuleGroup not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup", func() error {
		resp, err := tc.client.UpdateRuleGroup(tc.ctx, &wafv2.UpdateRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope,
			Id: aws.String(rgID), LockToken: aws.String(rgLockToken),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: true,
				MetricName:               aws.String("updated-metric"),
			},
		})
		if err != nil {
			return err
		}
		if resp.NextLockToken == nil || aws.ToString(resp.NextLockToken) == "" {
			return fmt.Errorf("NextLockToken is nil or empty")
		}
		rgLockToken = aws.ToString(resp.NextLockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup_ContentVerify", func() error {
		resp, err := tc.client.GetRuleGroup(tc.ctx, &wafv2.GetRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope, Id: aws.String(rgID),
		})
		if err != nil {
			return err
		}
		if resp.RuleGroup.VisibilityConfig == nil {
			return fmt.Errorf("VisibilityConfig is nil")
		}
		if aws.ToString(resp.RuleGroup.VisibilityConfig.MetricName) != "updated-metric" {
			return fmt.Errorf("expected metric name 'updated-metric', got '%s'", aws.ToString(resp.RuleGroup.VisibilityConfig.MetricName))
		}
		if !resp.RuleGroup.VisibilityConfig.SampledRequestsEnabled {
			return fmt.Errorf("expected SampledRequestsEnabled true")
		}
		if !resp.RuleGroup.VisibilityConfig.CloudWatchMetricsEnabled {
			return fmt.Errorf("expected CloudWatchMetricsEnabled true")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup_StaleLockToken", func() error {
		_, err := tc.client.UpdateRuleGroup(tc.ctx, &wafv2.UpdateRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope,
			Id: aws.String(rgID), LockToken: aws.String("stale-token-should-fail"),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("bad-metric"),
			},
		})
		return AssertErrorContains(err, "WAFInvalidLockTokenException")
	}))

	results = append(results, r.RunTest("wafv2", "DeleteRuleGroup", func() error {
		resp, err := tc.client.DeleteRuleGroup(tc.ctx, &wafv2.DeleteRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope,
			Id: aws.String(rgID), LockToken: aws.String(rgLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetRuleGroup_NonExistent", func() error {
		_, err := tc.client.GetRuleGroup(tc.ctx, &wafv2.GetRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope, Id: aws.String(rgID),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "CheckCapacity", func() error {
		capName := tc.uniqueName("cap-rg")
		capID, capLock, err := tc.createRuleGroup(capName, 50, &types.VisibilityConfig{
			SampledRequestsEnabled:   false,
			CloudWatchMetricsEnabled: false,
			MetricName:               aws.String("cap-metric"),
		})
		if err != nil {
			return fmt.Errorf("create rule group: %v", err)
		}
		defer tc.deleteRuleGroup(capName, capID, capLock)

		resp, err := tc.client.CheckCapacity(tc.ctx, &wafv2.CheckCapacityInput{
			Scope: tc.scope,
			Rules: []types.Rule{
				{
					Name:     aws.String("CapTestRule"),
					Priority: 0,
					Action: &types.RuleAction{
						Allow: &types.AllowAction{},
					},
					VisibilityConfig: &types.VisibilityConfig{
						SampledRequestsEnabled:   false,
						CloudWatchMetricsEnabled: false,
						MetricName:               aws.String("cap-rule-metric"),
					},
					Statement: &types.Statement{
						ByteMatchStatement: &types.ByteMatchStatement{
							FieldToMatch: &types.FieldToMatch{
								UriPath: &types.UriPath{},
							},
							PositionalConstraint: types.PositionalConstraintContains,
							SearchString:         []byte("/test"),
							TextTransformations: []types.TextTransformation{
								{Priority: 0, Type: types.TextTransformationTypeNone},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Capacity == 0 {
			return fmt.Errorf("expected non-zero capacity")
		}
		return nil
	}))

	return results
}
