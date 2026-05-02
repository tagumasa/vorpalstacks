package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2EdgeTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("wafv2", "CreateIPSet_Duplicate", func() error {
		dupName := tc.uniqueName("dup-ipset")
		id, _, lock, err := tc.createIPSet(dupName, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteIPSet(dupName, id, lock)

		_, err = tc.client.CreateIPSet(tc.ctx, &wafv2.CreateIPSetInput{
			Name:             aws.String(dupName),
			Scope:            tc.scope,
			IPAddressVersion: types.IPAddressVersionIpv4,
			Addresses:        []string{},
		})
		return AssertErrorContains(err, "WAFDuplicateItemException")
	}))

	results = append(results, r.RunTest("wafv2", "CreateWebACL_Duplicate", func() error {
		dupName := tc.uniqueName("dup-webacl")
		id, _, lock, err := tc.createWebACL(dupName,
			&types.DefaultAction{Allow: &types.AllowAction{}},
			&types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String("dup-metric"),
			}, nil, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteWebACL(dupName, id, lock)

		_, err = tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:          aws.String(dupName),
			Scope:         tc.scope,
			DefaultAction: &types.DefaultAction{Allow: &types.AllowAction{}},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String("dup-metric-2"),
			},
		})
		return AssertErrorContains(err, "WAFDuplicateItemException")
	}))

	results = append(results, r.RunTest("wafv2", "CreateRuleGroup_Duplicate", func() error {
		dupName := tc.uniqueName("dup-rg")
		id, lock, err := tc.createRuleGroup(dupName, 100, &types.VisibilityConfig{
			SampledRequestsEnabled: false, CloudWatchMetricsEnabled: false,
			MetricName: aws.String("dup-rg-metric"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteRuleGroup(dupName, id, lock)

		_, err = tc.client.CreateRuleGroup(tc.ctx, &wafv2.CreateRuleGroupInput{
			Name:     aws.String(dupName),
			Scope:    tc.scope,
			Capacity: aws.Int64(200),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: false, CloudWatchMetricsEnabled: false,
				MetricName: aws.String("dup-rg-metric-2"),
			},
		})
		return AssertErrorContains(err, "WAFDuplicateItemException")
	}))

	results = append(results, r.RunTest("wafv2", "CreateRegexPatternSet_Duplicate", func() error {
		dupName := tc.uniqueName("dup-regex")
		id, lock, err := tc.createRegexPatternSet(dupName, []types.Regex{
			{RegexString: aws.String("test")},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteRegexPatternSet(dupName, id, lock)

		_, err = tc.client.CreateRegexPatternSet(tc.ctx, &wafv2.CreateRegexPatternSetInput{
			Name:  aws.String(dupName),
			Scope: tc.scope,
			RegularExpressionList: []types.Regex{
				{RegexString: aws.String("test2")},
			},
		})
		return AssertErrorContains(err, "WAFDuplicateItemException")
	}))

	return results
}
