package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
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

	results = append(results, r.RunTest("wafv2", "DescribeManagedRuleGroup_LabelNamespace", func() error {
		resp, err := tc.client.DescribeManagedRuleGroup(tc.ctx, &wafv2.DescribeManagedRuleGroupInput{
			Name:       aws.String("AWSManagedRulesKnownBadInputsRuleSet"),
			VendorName: aws.String("AWS"),
			Scope:      tc.scope,
		})
		if err != nil {
			return err
		}
		if ns := aws.ToString(resp.LabelNamespace); ns != "awswaf:managed:aws:known-bad-inputs:" {
			return fmt.Errorf("LabelNamespace = %q, want awswaf:managed:aws:known-bad-inputs:", ns)
		}
		if aws.ToInt64(resp.Capacity) != 200 {
			return fmt.Errorf("Capacity = %d, want 200", aws.ToInt64(resp.Capacity))
		}
		var log4j *types.RuleSummary
		for i := range resp.Rules {
			if aws.ToString(resp.Rules[i].Name) == "Log4JRCE_QUERYSTRING" {
				log4j = &resp.Rules[i]
			}
		}
		if log4j == nil {
			return fmt.Errorf("Log4JRCE_QUERYSTRING not among the %d rules", len(resp.Rules))
		}
		if log4j.Action == nil || log4j.Action.Block == nil {
			return fmt.Errorf("Log4JRCE_QUERYSTRING action is not Block: %+v", log4j.Action)
		}
		wantLabel := "awswaf:managed:aws:known-bad-inputs:Log4JRCE_QueryString"
		found := false
		for _, label := range resp.AvailableLabels {
			if aws.ToString(label.Name) == wantLabel {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("label %q not among the available labels", wantLabel)
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "CreateWebACL_ManagedRuleGroup_RoundTrip", func() error {
		name := tc.uniqueName("test-managed-acl")
		create, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(name),
			Scope: tc.scope,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("managed-roundtrip-metric"),
			},
			Rules: []types.Rule{
				{
					Name:     aws.String("KnownBadInputs"),
					Priority: 0,
					Statement: &types.Statement{
						ManagedRuleGroupStatement: &types.ManagedRuleGroupStatement{
							VendorName: aws.String("AWS"),
							Name:       aws.String("AWSManagedRulesKnownBadInputsRuleSet"),
							ExcludedRules: []types.ExcludedRule{
								{Name: aws.String("Host_localhost_HEADER")},
							},
						},
					},
					VisibilityConfig: wafActionRuleVisibility("KnownBadInputs"),
				},
			},
		})
		if err != nil {
			return err
		}
		aclID := aws.ToString(create.Summary.Id)
		lockToken := aws.ToString(create.Summary.LockToken)
		defer tc.deleteWebACL(name, aclID, lockToken)

		resp, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(name), Scope: tc.scope, Id: aws.String(aclID),
		})
		if err != nil {
			return err
		}
		if len(resp.WebACL.Rules) != 1 {
			return fmt.Errorf("expected 1 rule, got %d", len(resp.WebACL.Rules))
		}
		stmt := resp.WebACL.Rules[0].Statement
		if stmt == nil || stmt.ManagedRuleGroupStatement == nil {
			return fmt.Errorf("managed rule group statement did not round-trip: %+v", stmt)
		}
		managed := stmt.ManagedRuleGroupStatement
		if aws.ToString(managed.VendorName) != "AWS" ||
			aws.ToString(managed.Name) != "AWSManagedRulesKnownBadInputsRuleSet" {
			return fmt.Errorf("managed group reference mismatch: %+v", managed)
		}
		if len(managed.ExcludedRules) != 1 || aws.ToString(managed.ExcludedRules[0].Name) != "Host_localhost_HEADER" {
			return fmt.Errorf("excluded rules did not round-trip: %+v", managed.ExcludedRules)
		}
		return nil
	}))

	return results
}
