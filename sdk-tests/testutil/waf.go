package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunWAFTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "waf",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := wafv2.NewFromConfig(cfg)
	ctx := context.Background()

	// Test 1: List Web ACLs
	results = append(results, r.RunTest("waf", "ListWebACLs", func() error {
		resp, err := client.ListWebACLs(ctx, &wafv2.ListWebACLsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.WebACLs == nil {
			return fmt.Errorf("web ACLs list is nil")
		}
		return nil
	}))

	// Test 2: Create IP Set
	ipSetName := fmt.Sprintf("test-ipset-%d", time.Now().UnixNano())
	var ipSetID string
	var ipSetARN string
	var ipSetLockToken string
	results = append(results, r.RunTest("waf", "CreateIPSet", func() error {
		resp, err := client.CreateIPSet(ctx, &wafv2.CreateIPSetInput{
			Name:             aws.String(ipSetName),
			Scope:            types.ScopeCloudfront,
			IPAddressVersion: types.IPAddressVersionIpv4,
			Addresses:        []string{},
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("sdk-tests")},
			},
		})
		if err == nil {
			ipSetID = aws.ToString(resp.Summary.Id)
			ipSetARN = aws.ToString(resp.Summary.ARN)
			ipSetLockToken = aws.ToString(resp.Summary.LockToken)
		}
		return err
	}))

	// Test 3: Get IP Set
	results = append(results, r.RunTest("waf", "GetIPSet", func() error {
		resp, err := client.GetIPSet(ctx, &wafv2.GetIPSetInput{
			Name:  aws.String(ipSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ipSetID),
		})
		if err != nil {
			return err
		}
		if resp.IPSet == nil {
			return fmt.Errorf("IP set is nil")
		}
		return nil
	}))

	// Test 4: List IP Sets
	results = append(results, r.RunTest("waf", "ListIPSets", func() error {
		resp, err := client.ListIPSets(ctx, &wafv2.ListIPSetsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.IPSets == nil {
			return fmt.Errorf("IP sets list is nil")
		}
		return nil
	}))

	// Test 5: List Tags for Resource (before update/delete)
	results = append(results, r.RunTest("waf", "ListTagsForResource", func() error {
		if ipSetARN == "" {
			return fmt.Errorf("no IP Set ARN available")
		}
		resp, err := client.ListTagsForResource(ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(ipSetARN),
		})
		if err != nil {
			return err
		}
		if resp.TagInfoForResource == nil || resp.TagInfoForResource.TagList == nil {
			return fmt.Errorf("expected tags in response")
		}
		if len(resp.TagInfoForResource.TagList) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.TagInfoForResource.TagList))
		}
		return nil
	}))

	// Test 6: Update IP Set
	results = append(results, r.RunTest("waf", "UpdateIPSet", func() error {
		resp, err := client.UpdateIPSet(ctx, &wafv2.UpdateIPSetInput{
			Name:      aws.String(ipSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ipSetID),
			LockToken: aws.String(ipSetLockToken),
			Addresses: []string{
				"192.0.2.0/24",
				"203.0.113.0/24",
			},
		})
		if err == nil && resp.NextLockToken != nil {
			ipSetLockToken = aws.ToString(resp.NextLockToken)
		}
		return err
	}))

	// Test 7: Delete IP Set
	results = append(results, r.RunTest("waf", "DeleteIPSet", func() error {
		resp, err := client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
			Name:      aws.String(ipSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ipSetID),
			LockToken: aws.String(ipSetLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 8: Create Regex Pattern Set
	regexPatternSetName := fmt.Sprintf("test-regex-%d", time.Now().UnixNano())
	var regexPatternSetID string
	var regexPatternSetLockToken string
	results = append(results, r.RunTest("waf", "CreateRegexPatternSet", func() error {
		resp, err := client.CreateRegexPatternSet(ctx, &wafv2.CreateRegexPatternSetInput{
			Name:  aws.String(regexPatternSetName),
			Scope: types.ScopeCloudfront,
			RegularExpressionList: []types.Regex{
				{
					RegexString: aws.String(`[a-z]+@[a-z]+\.[a-z]+`),
				},
			},
		})
		if err == nil {
			regexPatternSetID = aws.ToString(resp.Summary.Id)
			regexPatternSetLockToken = aws.ToString(resp.Summary.LockToken)
		}
		return err
	}))

	// Test 9: Get Regex Pattern Set
	results = append(results, r.RunTest("waf", "GetRegexPatternSet", func() error {
		resp, err := client.GetRegexPatternSet(ctx, &wafv2.GetRegexPatternSetInput{
			Name:  aws.String(regexPatternSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(regexPatternSetID),
		})
		if err != nil {
			return err
		}
		if resp.RegexPatternSet == nil {
			return fmt.Errorf("regex pattern set is nil")
		}
		return nil
	}))

	// Test 10: List Regex Pattern Sets
	results = append(results, r.RunTest("waf", "ListRegexPatternSets", func() error {
		resp, err := client.ListRegexPatternSets(ctx, &wafv2.ListRegexPatternSetsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.RegexPatternSets == nil {
			return fmt.Errorf("regex pattern sets list is nil")
		}
		return nil
	}))

	// Test 11: Delete Regex Pattern Set
	results = append(results, r.RunTest("waf", "DeleteRegexPatternSet", func() error {
		resp, err := client.DeleteRegexPatternSet(ctx, &wafv2.DeleteRegexPatternSetInput{
			Name:      aws.String(regexPatternSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(regexPatternSetID),
			LockToken: aws.String(regexPatternSetLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 12: Create Rule Group
	ruleGroupName := fmt.Sprintf("test-rulegroup-%d", time.Now().UnixNano())
	var ruleGroupID string
	var ruleGroupLockToken string
	results = append(results, r.RunTest("waf", "CreateRuleGroup", func() error {
		resp, err := client.CreateRuleGroup(ctx, &wafv2.CreateRuleGroupInput{
			Name:     aws.String(ruleGroupName),
			Scope:    types.ScopeCloudfront,
			Capacity: aws.Int64(100),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("test-metric"),
			},
		})
		if err == nil {
			ruleGroupID = aws.ToString(resp.Summary.Id)
			ruleGroupLockToken = aws.ToString(resp.Summary.LockToken)
		}
		return err
	}))

	// Test 13: Get Rule Group
	results = append(results, r.RunTest("waf", "GetRuleGroup", func() error {
		resp, err := client.GetRuleGroup(ctx, &wafv2.GetRuleGroupInput{
			Name:  aws.String(ruleGroupName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ruleGroupID),
		})
		if err != nil {
			return err
		}
		if resp.RuleGroup == nil {
			return fmt.Errorf("rule group is nil")
		}
		return nil
	}))

	// Test 14: List Rule Groups
	results = append(results, r.RunTest("waf", "ListRuleGroups", func() error {
		resp, err := client.ListRuleGroups(ctx, &wafv2.ListRuleGroupsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.RuleGroups == nil {
			return fmt.Errorf("rule groups list is nil")
		}
		return nil
	}))

	// Test 15: Delete Rule Group
	results = append(results, r.RunTest("waf", "DeleteRuleGroup", func() error {
		resp, err := client.DeleteRuleGroup(ctx, &wafv2.DeleteRuleGroupInput{
			Name:      aws.String(ruleGroupName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ruleGroupID),
			LockToken: aws.String(ruleGroupLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 16: List Available Managed Rule Groups
	results = append(results, r.RunTest("waf", "ListAvailableManagedRuleGroups", func() error {
		resp, err := client.ListAvailableManagedRuleGroups(ctx, &wafv2.ListAvailableManagedRuleGroupsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.ManagedRuleGroups == nil {
			return fmt.Errorf("managed rule groups list is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("waf", "GetIPSet_NonExistent", func() error {
		if ipSetID == "" {
			return fmt.Errorf("IP Set ID not available")
		}
		_, err := client.GetIPSet(ctx, &wafv2.GetIPSetInput{
			Name:  aws.String(ipSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ipSetID),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted IP Set")
		}
		return nil
	}))

	results = append(results, r.RunTest("waf", "DeleteIPSet_NonExistent", func() error {
		if ipSetID == "" {
			return fmt.Errorf("IP Set ID not available")
		}
		_, err := client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
			Name:      aws.String(ipSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ipSetID),
			LockToken: aws.String("fake-lock-token"),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted IP Set")
		}
		return nil
	}))

	results = append(results, r.RunTest("waf", "GetRegexPatternSet_NonExistent", func() error {
		if regexPatternSetID == "" {
			return fmt.Errorf("Regex Pattern Set ID not available")
		}
		_, err := client.GetRegexPatternSet(ctx, &wafv2.GetRegexPatternSetInput{
			Name:  aws.String(regexPatternSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(regexPatternSetID),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted Regex Pattern Set")
		}
		return nil
	}))

	results = append(results, r.RunTest("waf", "GetRuleGroup_NonExistent", func() error {
		if ruleGroupID == "" {
			return fmt.Errorf("Rule Group ID not available")
		}
		_, err := client.GetRuleGroup(ctx, &wafv2.GetRuleGroupInput{
			Name:  aws.String(ruleGroupName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ruleGroupID),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted Rule Group")
		}
		return nil
	}))

	var verifyIPSetID, verifyIPSetLockToken string
	verifyIPSetName := fmt.Sprintf("verify-ipset-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("waf", "UpdateIPSet_ContentVerify", func() error {
		resp, err := client.CreateIPSet(ctx, &wafv2.CreateIPSetInput{
			Name:             aws.String(verifyIPSetName),
			Scope:            types.ScopeCloudfront,
			IPAddressVersion: types.IPAddressVersionIpv4,
			Addresses:        []string{},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		verifyIPSetID = aws.ToString(resp.Summary.Id)
		verifyIPSetLockToken = aws.ToString(resp.Summary.LockToken)

		newAddresses := []string{"10.0.0.0/8", "172.16.0.0/12"}
		updateResp, err := client.UpdateIPSet(ctx, &wafv2.UpdateIPSetInput{
			Name:      aws.String(verifyIPSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(verifyIPSetID),
			LockToken: aws.String(verifyIPSetLockToken),
			Addresses: newAddresses,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		verifyIPSetLockToken = aws.ToString(updateResp.NextLockToken)

		getResp, err := client.GetIPSet(ctx, &wafv2.GetIPSetInput{
			Name:  aws.String(verifyIPSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(verifyIPSetID),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.IPSet.Addresses) != 2 {
			return fmt.Errorf("expected 2 addresses, got %d", len(getResp.IPSet.Addresses))
		}
		return nil
	}))

	results = append(results, r.RunTest("waf", "ListIPSets_ContainsCreated", func() error {
		if verifyIPSetID == "" {
			return fmt.Errorf("IP Set ID not available")
		}
		defer func() {
			client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
				Name:      aws.String(verifyIPSetName),
				Scope:     types.ScopeCloudfront,
				Id:        aws.String(verifyIPSetID),
				LockToken: aws.String(verifyIPSetLockToken),
			})
		}()
		resp, err := client.ListIPSets(ctx, &wafv2.ListIPSetsInput{
			Scope: types.ScopeCloudfront,
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.IPSets {
			if s.Id != nil && *s.Id == verifyIPSetID {
				found = true
				if s.Name == nil || *s.Name != verifyIPSetName {
					return fmt.Errorf("IP Set name mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("IP Set not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("waf", "UpdateIPSet_StaleLockToken", func() error {
		if verifyIPSetID == "" {
			return fmt.Errorf("IP Set ID not available")
		}
		_, err := client.UpdateIPSet(ctx, &wafv2.UpdateIPSetInput{
			Name:      aws.String(verifyIPSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(verifyIPSetID),
			LockToken: aws.String("stale-token-should-fail"),
			Addresses: []string{"192.168.0.0/16"},
		})
		if err == nil {
			return fmt.Errorf("expected error for stale lock token")
		}
		return nil
	}))

	return results
}
