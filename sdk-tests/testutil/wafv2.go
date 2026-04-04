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

func (r *TestRunner) RunWAFv2Tests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "wafv2",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := wafv2.NewFromConfig(cfg)
	ctx := context.Background()

	// === IP Set CRUD ===

	results = append(results, r.RunTest("wafv2", "ListIPSets_Empty", func() error {
		_, err := client.ListIPSets(ctx, &wafv2.ListIPSetsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		return err
	}))

	ipSetName := fmt.Sprintf("test-ipset-%d", time.Now().UnixNano())
	var ipSetID string
	var ipSetARN string
	var ipSetLockToken string
	results = append(results, r.RunTest("wafv2", "CreateIPSet", func() error {
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

	results = append(results, r.RunTest("wafv2", "GetIPSet", func() error {
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
		if resp.IPSet.IPAddressVersion != types.IPAddressVersionIpv4 {
			return fmt.Errorf("expected IPV4, got %s", resp.IPSet.IPAddressVersion)
		}
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListIPSets_ContainsCreated", func() error {
		resp, err := client.ListIPSets(ctx, &wafv2.ListIPSetsInput{
			Scope: types.ScopeCloudfront,
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.IPSets {
			if s.Id != nil && *s.Id == ipSetID {
				found = true
				if s.Name == nil || *s.Name != ipSetName {
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

	results = append(results, r.RunTest("wafv2", "ListTagsForResource", func() error {
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

	results = append(results, r.RunTest("wafv2", "UpdateIPSet", func() error {
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

	results = append(results, r.RunTest("wafv2", "UpdateIPSet_ContentVerify", func() error {
		resp, err := client.GetIPSet(ctx, &wafv2.GetIPSetInput{
			Name:  aws.String(ipSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ipSetID),
		})
		if err != nil {
			return err
		}
		if len(resp.IPSet.Addresses) != 2 {
			return fmt.Errorf("expected 2 addresses, got %d", len(resp.IPSet.Addresses))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteIPSet", func() error {
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

	results = append(results, r.RunTest("wafv2", "GetIPSet_NonExistent", func() error {
		_, err := client.GetIPSet(ctx, &wafv2.GetIPSetInput{
			Name:  aws.String(ipSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ipSetID),
		})
		if err := AssertErrorContains(err, "WAFNonexistentItemException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteIPSet_NonExistent", func() error {
		_, err := client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
			Name:      aws.String(ipSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ipSetID),
			LockToken: aws.String("fake-lock-token"),
		})
		if err := AssertErrorContains(err, "WAFNonexistentItemException"); err != nil {
			return err
		}
		return nil
	}))

	// === Regex Pattern Set CRUD ===

	regexPatternSetName := fmt.Sprintf("test-regex-%d", time.Now().UnixNano())
	var regexPatternSetID string
	var regexPatternSetLockToken string
	results = append(results, r.RunTest("wafv2", "CreateRegexPatternSet", func() error {
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

	results = append(results, r.RunTest("wafv2", "GetRegexPatternSet", func() error {
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
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListRegexPatternSets", func() error {
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

	results = append(results, r.RunTest("wafv2", "UpdateRegexPatternSet", func() error {
		resp, err := client.UpdateRegexPatternSet(ctx, &wafv2.UpdateRegexPatternSetInput{
			Name:      aws.String(regexPatternSetName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(regexPatternSetID),
			LockToken: aws.String(regexPatternSetLockToken),
			RegularExpressionList: []types.Regex{
				{RegexString: aws.String(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)},
				{RegexString: aws.String(`[0-9a-f]{2}:[0-9a-f]{2}:[0-9a-f]{2}`)},
			},
		})
		if err == nil && resp.NextLockToken != nil {
			regexPatternSetLockToken = aws.ToString(resp.NextLockToken)
		}
		return err
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRegexPatternSet_ContentVerify", func() error {
		resp, err := client.GetRegexPatternSet(ctx, &wafv2.GetRegexPatternSetInput{
			Name:  aws.String(regexPatternSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(regexPatternSetID),
		})
		if err != nil {
			return err
		}
		if len(resp.RegexPatternSet.RegularExpressionList) != 2 {
			return fmt.Errorf("expected 2 patterns, got %d", len(resp.RegexPatternSet.RegularExpressionList))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteRegexPatternSet", func() error {
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

	results = append(results, r.RunTest("wafv2", "GetRegexPatternSet_NonExistent", func() error {
		_, err := client.GetRegexPatternSet(ctx, &wafv2.GetRegexPatternSetInput{
			Name:  aws.String(regexPatternSetName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(regexPatternSetID),
		})
		if err := AssertErrorContains(err, "WAFNonexistentItemException"); err != nil {
			return err
		}
		return nil
	}))

	// === Rule Group CRUD ===

	ruleGroupName := fmt.Sprintf("test-rulegroup-%d", time.Now().UnixNano())
	var ruleGroupID string
	var ruleGroupLockToken string
	results = append(results, r.RunTest("wafv2", "CreateRuleGroup", func() error {
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

	results = append(results, r.RunTest("wafv2", "GetRuleGroup", func() error {
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
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		if resp.RuleGroup.Capacity == nil || *resp.RuleGroup.Capacity != 100 {
			return fmt.Errorf("expected capacity 100, got %d", resp.RuleGroup.Capacity)
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListRuleGroups", func() error {
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

	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup", func() error {
		resp, err := client.UpdateRuleGroup(ctx, &wafv2.UpdateRuleGroupInput{
			Name:      aws.String(ruleGroupName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ruleGroupID),
			LockToken: aws.String(ruleGroupLockToken),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: true,
				MetricName:               aws.String("updated-metric"),
			},
		})
		if err == nil && resp.NextLockToken != nil {
			ruleGroupLockToken = aws.ToString(resp.NextLockToken)
		}
		return err
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup_ContentVerify", func() error {
		resp, err := client.GetRuleGroup(ctx, &wafv2.GetRuleGroupInput{
			Name:  aws.String(ruleGroupName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ruleGroupID),
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
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup_StaleLockToken", func() error {
		_, err := client.UpdateRuleGroup(ctx, &wafv2.UpdateRuleGroupInput{
			Name:      aws.String(ruleGroupName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(ruleGroupID),
			LockToken: aws.String("stale-token-should-fail"),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("bad-metric"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for stale lock token")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteRuleGroup", func() error {
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

	results = append(results, r.RunTest("wafv2", "GetRuleGroup_NonExistent", func() error {
		_, err := client.GetRuleGroup(ctx, &wafv2.GetRuleGroupInput{
			Name:  aws.String(ruleGroupName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(ruleGroupID),
		})
		if err := AssertErrorContains(err, "WAFNonexistentItemException"); err != nil {
			return err
		}
		return nil
	}))

	// === Web ACL CRUD ===

	results = append(results, r.RunTest("wafv2", "ListWebACLs_Empty", func() error {
		_, err := client.ListWebACLs(ctx, &wafv2.ListWebACLsInput{
			Scope: types.ScopeCloudfront,
			Limit: aws.Int32(10),
		})
		return err
	}))

	webACLName := fmt.Sprintf("test-webacl-%d", time.Now().UnixNano())
	var webACLID string
	var webACLARN string
	var webACLLockToken string
	results = append(results, r.RunTest("wafv2", "CreateWebACL", func() error {
		resp, err := client.CreateWebACL(ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(webACLName),
			Scope: types.ScopeCloudfront,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: true,
				MetricName:               aws.String("webacl-test-metric"),
			},
			Description: aws.String("Test WebACL for SDK tests"),
			Rules: []types.Rule{
				{
					Name:     aws.String("AllowTestRule"),
					Priority: 0,
					Action: &types.RuleAction{
						Allow: &types.AllowAction{},
					},
					VisibilityConfig: &types.VisibilityConfig{
						SampledRequestsEnabled:   true,
						CloudWatchMetricsEnabled: true,
						MetricName:               aws.String("allow-test-rule-metric"),
					},
					Statement: &types.Statement{
						ByteMatchStatement: &types.ByteMatchStatement{
							FieldToMatch: &types.FieldToMatch{
								UriPath: &types.UriPath{},
							},
							PositionalConstraint: types.PositionalConstraintContains,
							SearchString:         []byte("/test-path"),
							TextTransformations: []types.TextTransformation{
								{Priority: 0, Type: types.TextTransformationTypeNone},
							},
						},
					},
				},
			},
			Tags: []types.Tag{
				{Key: aws.String("Name"), Value: aws.String("WebACLTest")},
			},
		})
		if err == nil {
			webACLID = aws.ToString(resp.Summary.Id)
			webACLARN = aws.ToString(resp.Summary.ARN)
			webACLLockToken = aws.ToString(resp.Summary.LockToken)
		}
		return err
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACL", func() error {
		resp, err := client.GetWebACL(ctx, &wafv2.GetWebACLInput{
			Name:  aws.String(webACLName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(webACLID),
		})
		if err != nil {
			return err
		}
		if resp.WebACL == nil {
			return fmt.Errorf("WebACL is nil")
		}
		if aws.ToString(resp.WebACL.Name) != webACLName {
			return fmt.Errorf("expected name '%s', got '%s'", webACLName, aws.ToString(resp.WebACL.Name))
		}
		if resp.LockToken == nil {
			return fmt.Errorf("LockToken is nil")
		}
		if resp.WebACL.DefaultAction == nil || resp.WebACL.DefaultAction.Allow == nil {
			return fmt.Errorf("expected default Allow action")
		}
		if resp.WebACL.VisibilityConfig == nil {
			return fmt.Errorf("VisibilityConfig is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListWebACLs_ContainsCreated", func() error {
		resp, err := client.ListWebACLs(ctx, &wafv2.ListWebACLsInput{
			Scope: types.ScopeCloudfront,
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.WebACLs {
			if s.Id != nil && *s.Id == webACLID {
				found = true
				if s.Name == nil || *s.Name != webACLName {
					return fmt.Errorf("WebACL name mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("WebACL not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateWebACL", func() error {
		resp, err := client.UpdateWebACL(ctx, &wafv2.UpdateWebACLInput{
			Name:      aws.String(webACLName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(webACLID),
			LockToken: aws.String(webACLLockToken),
			DefaultAction: &types.DefaultAction{
				Block: &types.BlockAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: true,
				MetricName:               aws.String("updated-webacl-metric"),
			},
		})
		if err == nil && resp.NextLockToken != nil {
			webACLLockToken = aws.ToString(resp.NextLockToken)
		}
		return err
	}))

	results = append(results, r.RunTest("wafv2", "UpdateWebACL_ContentVerify", func() error {
		resp, err := client.GetWebACL(ctx, &wafv2.GetWebACLInput{
			Name:  aws.String(webACLName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(webACLID),
		})
		if err != nil {
			return err
		}
		if resp.WebACL.DefaultAction == nil || resp.WebACL.DefaultAction.Block == nil {
			return fmt.Errorf("expected default Block action after update")
		}
		if aws.ToString(resp.WebACL.VisibilityConfig.MetricName) != "updated-webacl-metric" {
			return fmt.Errorf("expected updated metric name, got '%s'", aws.ToString(resp.WebACL.VisibilityConfig.MetricName))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateWebACL_StaleLockToken", func() error {
		_, err := client.UpdateWebACL(ctx, &wafv2.UpdateWebACLInput{
			Name:      aws.String(webACLName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(webACLID),
			LockToken: aws.String("stale-lock-token"),
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   false,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("bad-metric"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for stale lock token")
		}
		return nil
	}))

	// === Tag Operations on WebACL ===

	results = append(results, r.RunTest("wafv2", "TagResource_WebACL", func() error {
		_, err := client.TagResource(ctx, &wafv2.TagResourceInput{
			ResourceARN: aws.String(webACLARN),
			Tags: []types.Tag{
				{Key: aws.String("Team"), Value: aws.String("Security")},
				{Key: aws.String("Env"), Value: aws.String("Production")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("wafv2", "ListTagsForResource_WebACL", func() error {
		resp, err := client.ListTagsForResource(ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(webACLARN),
		})
		if err != nil {
			return err
		}
		if len(resp.TagInfoForResource.TagList) < 3 {
			return fmt.Errorf("expected at least 3 tags (1 create + 2 tagged), got %d", len(resp.TagInfoForResource.TagList))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UntagResource_WebACL", func() error {
		_, err := client.UntagResource(ctx, &wafv2.UntagResourceInput{
			ResourceARN: aws.String(webACLARN),
			TagKeys:     []string{"Env"},
		})
		return err
	}))

	results = append(results, r.RunTest("wafv2", "UntagResource_Verify", func() error {
		resp, err := client.ListTagsForResource(ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(webACLARN),
		})
		if err != nil {
			return err
		}
		for _, t := range resp.TagInfoForResource.TagList {
			if aws.ToString(t.Key) == "Env" {
				return fmt.Errorf("tag 'Env' should have been removed")
			}
		}
		return nil
	}))

	// === Logging Configuration ===

	results = append(results, r.RunTest("wafv2", "PutLoggingConfiguration", func() error {
		resp, err := client.PutLoggingConfiguration(ctx, &wafv2.PutLoggingConfigurationInput{
			LoggingConfiguration: &types.LoggingConfiguration{
				ResourceArn: aws.String(webACLARN),
				LogDestinationConfigs: []string{
					"arn:aws:logs:us-east-1:123456789012:log-group:/aws/waf/test-log",
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.LoggingConfiguration == nil {
			return fmt.Errorf("LoggingConfiguration is nil")
		}
		if len(resp.LoggingConfiguration.LogDestinationConfigs) != 1 {
			return fmt.Errorf("expected 1 log destination, got %d", len(resp.LoggingConfiguration.LogDestinationConfigs))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetLoggingConfiguration", func() error {
		resp, err := client.GetLoggingConfiguration(ctx, &wafv2.GetLoggingConfigurationInput{
			ResourceArn: aws.String(webACLARN),
		})
		if err != nil {
			return err
		}
		if resp.LoggingConfiguration == nil {
			return fmt.Errorf("LoggingConfiguration is nil")
		}
		if aws.ToString(resp.LoggingConfiguration.ResourceArn) != webACLARN {
			return fmt.Errorf("ResourceArn mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteLoggingConfiguration", func() error {
		_, err := client.DeleteLoggingConfiguration(ctx, &wafv2.DeleteLoggingConfigurationInput{
			ResourceArn: aws.String(webACLARN),
		})
		return err
	}))

	results = append(results, r.RunTest("wafv2", "GetLoggingConfiguration_AfterDelete", func() error {
		_, err := client.GetLoggingConfiguration(ctx, &wafv2.GetLoggingConfigurationInput{
			ResourceArn: aws.String(webACLARN),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting logging config")
		}
		return nil
	}))

	// === Association Operations ===

	fakeResourceARN := fmt.Sprintf("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/%d", time.Now().UnixNano())

	results = append(results, r.RunTest("wafv2", "AssociateWebACL", func() error {
		_, err := client.AssociateWebACL(ctx, &wafv2.AssociateWebACLInput{
			WebACLArn:   aws.String(webACLARN),
			ResourceArn: aws.String(fakeResourceARN),
		})
		return err
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACLForResource", func() error {
		resp, err := client.GetWebACLForResource(ctx, &wafv2.GetWebACLForResourceInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		if err != nil {
			return err
		}
		if resp.WebACL == nil {
			return fmt.Errorf("WebACL is nil")
		}
		if aws.ToString(resp.WebACL.Name) != webACLName {
			return fmt.Errorf("expected WebACL name '%s', got '%s'", webACLName, aws.ToString(resp.WebACL.Name))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListResourcesForWebACL", func() error {
		resp, err := client.ListResourcesForWebACL(ctx, &wafv2.ListResourcesForWebACLInput{
			WebACLArn: aws.String(webACLARN),
		})
		if err != nil {
			return err
		}
		if len(resp.ResourceArns) == 0 {
			return fmt.Errorf("expected at least 1 resource ARN")
		}
		found := false
		for _, arn := range resp.ResourceArns {
			if arn == fakeResourceARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("resource ARN not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DisassociateWebACL", func() error {
		_, err := client.DisassociateWebACL(ctx, &wafv2.DisassociateWebACLInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		return err
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACLForResource_AfterDisassociate", func() error {
		_, err := client.GetWebACLForResource(ctx, &wafv2.GetWebACLForResourceInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		if err == nil {
			return fmt.Errorf("expected error after disassociation")
		}
		return nil
	}))

	// === Managed Rule Groups ===

	results = append(results, r.RunTest("wafv2", "ListAvailableManagedRuleGroups", func() error {
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

	results = append(results, r.RunTest("wafv2", "DescribeManagedRuleGroup", func() error {
		resp, err := client.DescribeManagedRuleGroup(ctx, &wafv2.DescribeManagedRuleGroupInput{
			Name:       aws.String("AWSManagedRulesCommonRuleSet"),
			VendorName: aws.String("AWS"),
			Scope:      types.ScopeCloudfront,
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
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DescribeManagedRuleGroup_NotFound", func() error {
		_, err := client.DescribeManagedRuleGroup(ctx, &wafv2.DescribeManagedRuleGroupInput{
			Name:       aws.String("NonExistentRuleGroup"),
			VendorName: aws.String("AWS"),
			Scope:      types.ScopeCloudfront,
		})
		if err := AssertErrorContains(err, "WAFNonexistentItemException"); err != nil {
			return err
		}
		return nil
	}))

	// === Cleanup: Delete WebACL ===

	results = append(results, r.RunTest("wafv2", "DeleteWebACL", func() error {
		resp, err := client.DeleteWebACL(ctx, &wafv2.DeleteWebACLInput{
			Name:      aws.String(webACLName),
			Scope:     types.ScopeCloudfront,
			Id:        aws.String(webACLID),
			LockToken: aws.String(webACLLockToken),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACL_NonExistent", func() error {
		_, err := client.GetWebACL(ctx, &wafv2.GetWebACLInput{
			Name:  aws.String(webACLName),
			Scope: types.ScopeCloudfront,
			Id:    aws.String(webACLID),
		})
		if err := AssertErrorContains(err, "WAFNonexistentItemException"); err != nil {
			return err
		}
		return nil
	}))

	// === IP Set Stale Lock Token ===

	var verifyIPSetID, verifyIPSetLockToken string
	verifyIPSetName := fmt.Sprintf("verify-ipset-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("wafv2", "UpdateIPSet_StaleLockToken", func() error {
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

		_, err = client.UpdateIPSet(ctx, &wafv2.UpdateIPSetInput{
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

	client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
		Name:      aws.String(verifyIPSetName),
		Scope:     types.ScopeCloudfront,
		Id:        aws.String(verifyIPSetID),
		LockToken: aws.String(verifyIPSetLockToken),
	})

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("wafv2", "ListIPSets_Pagination", func() error {
		pgTs := fmt.Sprintf("pgwaf%d", time.Now().UnixNano())
		type pgIPSet struct {
			id, name, lockToken string
		}
		var pgSets []pgIPSet
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("%s-%d", pgTs, i)
			resp, err := client.CreateIPSet(ctx, &wafv2.CreateIPSetInput{
				Name:             aws.String(name),
				Scope:            types.ScopeCloudfront,
				IPAddressVersion: types.IPAddressVersionIpv4,
				Addresses:        []string{},
			})
			if err != nil {
				return fmt.Errorf("create ip set %s: %v", name, err)
			}
			pgSets = append(pgSets, pgIPSet{
				id:        aws.ToString(resp.Summary.Id),
				name:      name,
				lockToken: aws.ToString(resp.Summary.LockToken),
			})
		}

		pgIDSet := make(map[string]bool)
		for _, s := range pgSets {
			pgIDSet[s.id] = true
		}
		var foundCount int
		var nextMarker *string
		for {
			resp, err := client.ListIPSets(ctx, &wafv2.ListIPSetsInput{
				Scope:      types.ScopeCloudfront,
				Limit:      aws.Int32(2),
				NextMarker: nextMarker,
			})
			if err != nil {
				for _, s := range pgSets {
					client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
						Name: aws.String(s.name), Scope: types.ScopeCloudfront,
						Id: aws.String(s.id), LockToken: aws.String(s.lockToken),
					})
				}
				return fmt.Errorf("list ip sets page: %v", err)
			}
			for _, s := range resp.IPSets {
				if pgIDSet[aws.ToString(s.Id)] {
					foundCount++
				}
			}
			if resp.NextMarker != nil && *resp.NextMarker != "" {
				nextMarker = resp.NextMarker
			} else {
				break
			}
		}

		for _, s := range pgSets {
			client.DeleteIPSet(ctx, &wafv2.DeleteIPSetInput{
				Name:      aws.String(s.name),
				Scope:     types.ScopeCloudfront,
				Id:        aws.String(s.id),
				LockToken: aws.String(s.lockToken),
			})
		}
		if foundCount != 5 {
			return fmt.Errorf("expected 5 paginated IP sets, got %d", foundCount)
		}
		return nil
	}))

	return results
}
