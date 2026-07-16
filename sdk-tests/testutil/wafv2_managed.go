package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
)

func (r *TestRunner) runWAFv2ManagedRulesTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("wafv2", "ListAvailableManagedRuleGroups", func() error {
		resp, err := tc.client.ListAvailableManagedRuleGroups(tc.ctx, &wafv2.ListAvailableManagedRuleGroupsInput{
			Scope: tc.scope,
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.ManagedRuleGroups == nil {
			return fmt.Errorf("managed rule groups list is nil")
		}
		found := false
		for _, g := range resp.ManagedRuleGroups {
			if aws.ToString(g.Name) == "AWSManagedRulesCommonRuleSet" {
				found = true
				if aws.ToString(g.VendorName) != "AWS" {
					return fmt.Errorf("expected vendor 'AWS', got '%s'", aws.ToString(g.VendorName))
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("AWSManagedRulesCommonRuleSet not found in managed rule groups")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DescribeManagedRuleGroup", func() error {
		resp, err := tc.client.DescribeManagedRuleGroup(tc.ctx, &wafv2.DescribeManagedRuleGroupInput{
			Name:       aws.String("AWSManagedRulesCommonRuleSet"),
			VendorName: aws.String("AWS"),
			Scope:      tc.scope,
		})
		if err != nil {
			return err
		}
		if resp.Capacity == nil || *resp.Capacity == 0 {
			return fmt.Errorf("expected non-zero capacity")
		}
		if resp.LabelNamespace == nil {
			return fmt.Errorf("LabelNamespace is nil")
		}
		if resp.Rules == nil || len(resp.Rules) == 0 {
			return fmt.Errorf("expected non-empty rules list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DescribeManagedRuleGroup_NotFound", func() error {
		_, err := tc.client.DescribeManagedRuleGroup(tc.ctx, &wafv2.DescribeManagedRuleGroupInput{
			Name:       aws.String("NonExistentRuleGroup"),
			VendorName: aws.String("AWS"),
			Scope:      tc.scope,
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "ListAvailableManagedRuleGroupVersions", func() error {
		resp, err := tc.client.ListAvailableManagedRuleGroupVersions(tc.ctx, &wafv2.ListAvailableManagedRuleGroupVersionsInput{
			Name:       aws.String("AWSManagedRulesCommonRuleSet"),
			VendorName: aws.String("AWS"),
			Scope:      tc.scope,
			Limit:      aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("Versions is nil")
		}
		return nil
	}))

	return results
}
