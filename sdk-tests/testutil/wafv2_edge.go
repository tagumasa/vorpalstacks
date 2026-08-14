package testutil

import (
	"fmt"
	"time"

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

	// VisibilityConfig.MetricName is validated against the Smithy
	// pattern; an invalid name must be rejected.
	results = append(results, r.RunTest("wafv2", "CreateWebACL_InvalidMetricName_Rejected", func() error {
		_, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:          aws.String(tc.uniqueName("bad-metric-acl")),
			Scope:         tc.scope,
			DefaultAction: &types.DefaultAction{Allow: &types.AllowAction{}},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String("bad metric!"),
			},
		})
		return AssertErrorContains(err, "WAFInvalidParameterException")
	}))

	// An unknown ResourceType enum value must be rejected instead of
	// silently returning an empty list.
	results = append(results, r.RunTest("wafv2", "ListResourcesForWebACL_UnknownResourceType_Rejected", func() error {
		_, err := tc.client.ListResourcesForWebACL(tc.ctx, &wafv2.ListResourcesForWebACLInput{
			WebACLArn:    aws.String("arn:aws:wafv2:" + tc.region + ":" + tc.accountID + ":regional/webacl/nosuch/" + tc.uniqueName("nosuch")),
			ResourceType: types.ResourceType("EC2_INSTANCE"),
		})
		return AssertErrorContains(err, "WAFInvalidParameterException")
	}))

	// A rule group with an explicit Capacity above the WCU quota must
	// be rejected with WAFLimitsExceededException.
	results = append(results, r.RunTest("wafv2", "CreateRuleGroup_CapacityOverLimit_Rejected", func() error {
		_, err := tc.client.CreateRuleGroup(tc.ctx, &wafv2.CreateRuleGroupInput{
			Name:     aws.String(tc.uniqueName("over-capacity-rg")),
			Scope:    tc.scope,
			Capacity: aws.Int64(5001),
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("cap-metric")),
			},
		})
		return AssertErrorContains(err, "WAFLimitsExceededException")
	}))

	// A WebACL that is still associated with a resource must not be
	// deletable (WAFAssociatedItemException).
	results = append(results, r.RunTest("wafv2", "DeleteWebACL_WhileAssociated_Rejected", func() error {
		name := tc.uniqueName("assoc-del-acl")
		id, arn, lock, err := tc.createWebACL(name,
			&types.DefaultAction{Allow: &types.AllowAction{}},
			&types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("ad-metric")),
			}, nil, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		resourceArn := fmt.Sprintf("arn:aws:elasticloadbalancing:%s:%s:loadbalancer/app/assoc-del-lb/%d",
			tc.region, tc.accountID, time.Now().UnixNano())
		if _, err := tc.client.AssociateWebACL(tc.ctx, &wafv2.AssociateWebACLInput{
			WebACLArn: aws.String(arn), ResourceArn: aws.String(resourceArn),
		}); err != nil {
			return fmt.Errorf("associate: %v", err)
		}

		_, delErr := tc.client.DeleteWebACL(tc.ctx, &wafv2.DeleteWebACLInput{
			Name: aws.String(name), Scope: tc.scope, Id: aws.String(id), LockToken: aws.String(lock),
		})
		assertErr := AssertErrorContains(delErr, "WAFAssociatedItemException")
		if delErr == nil {
			return fmt.Errorf("DeleteWebACL unexpectedly succeeded while the web ACL was still associated with a resource")
		}

		// Clean up regardless of the assertion outcome.
		if _, err := tc.client.DisassociateWebACL(tc.ctx, &wafv2.DisassociateWebACLInput{
			ResourceArn: aws.String(resourceArn),
		}); err != nil && assertErr == nil {
			return fmt.Errorf("cleanup disassociate: %v", err)
		}
		if _, err := tc.client.DeleteWebACL(tc.ctx, &wafv2.DeleteWebACLInput{
			Name: aws.String(name), Scope: tc.scope, Id: aws.String(id), LockToken: aws.String(lock),
		}); err != nil && assertErr == nil {
			return fmt.Errorf("cleanup delete: %v", err)
		}
		return assertErr
	}))

	// A rule group still referenced by a web ACL rule must not be
	// deletable (WAFAssociatedItemException).
	results = append(results, r.RunTest("wafv2", "DeleteRuleGroup_WhileReferenced_Rejected", func() error {
		rgName := tc.uniqueName("ref-rg")
		rgID, rgLock, err := tc.createRuleGroup(rgName, 10, &types.VisibilityConfig{
			SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
			MetricName: aws.String(tc.uniqueName("ref-rg-metric")),
		})
		if err != nil {
			return fmt.Errorf("create rule group: %v", err)
		}
		rgResp, err := tc.client.GetRuleGroup(tc.ctx, &wafv2.GetRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope, Id: aws.String(rgID),
		})
		if err != nil {
			return fmt.Errorf("get rule group: %v", err)
		}
		rgArn := aws.ToString(rgResp.RuleGroup.ARN)

		aclName := tc.uniqueName("ref-acl")
		aclID, _, aclLock, err := tc.createWebACL(aclName,
			&types.DefaultAction{Allow: &types.AllowAction{}},
			&types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("ref-acl-metric")),
			},
			[]types.Rule{{
				Name:     aws.String("use-rg"),
				Priority: 1,
				Action:   &types.RuleAction{Count: &types.CountAction{}},
				Statement: &types.Statement{
					RuleGroupReferenceStatement: &types.RuleGroupReferenceStatement{ARN: aws.String(rgArn)},
				},
				VisibilityConfig: &types.VisibilityConfig{
					SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
					MetricName: aws.String(tc.uniqueName("ref-rule-metric")),
				},
			}}, nil)
		if err != nil {
			return fmt.Errorf("create web acl: %v", err)
		}
		// The rule group can only be deleted after the referencing web
		// ACL is gone, so cleanup runs in reverse creation order.
		defer tc.deleteRuleGroup(rgName, rgID, rgLock)
		defer tc.deleteWebACL(aclName, aclID, aclLock)

		_, delErr := tc.client.DeleteRuleGroup(tc.ctx, &wafv2.DeleteRuleGroupInput{
			Name: aws.String(rgName), Scope: tc.scope, Id: aws.String(rgID), LockToken: aws.String(rgLock),
		})
		return AssertErrorContains(delErr, "WAFAssociatedItemException")
	}))

	// UpdateWebACL is a full-replace operation: omitting Rules clears
	// the rule list and the consumed capacity drops accordingly.
	results = append(results, r.RunTest("wafv2", "UpdateWebACL_RulesOmitted_ClearsRulesAndCapacity", func() error {
		name := tc.uniqueName("fullreplace-acl")
		id, _, lock, err := tc.createWebACL(name,
			&types.DefaultAction{Allow: &types.AllowAction{}},
			&types.VisibilityConfig{SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("fr-metric"))},
			[]types.Rule{{
				Name:     aws.String("base-rule"),
				Priority: 1,
				Action:   &types.RuleAction{Count: &types.CountAction{}},
				Statement: &types.Statement{
					ByteMatchStatement: &types.ByteMatchStatement{
						PositionalConstraint: types.PositionalConstraintContains,
						SearchString:         []byte("/"),
						FieldToMatch:         &types.FieldToMatch{UriPath: &types.UriPath{}},
						TextTransformations: []types.TextTransformation{{
							Priority: 1, Type: types.TextTransformationTypeNone,
						}},
					},
				},
				VisibilityConfig: &types.VisibilityConfig{
					SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
					MetricName: aws.String(tc.uniqueName("fr-rule-metric")),
				},
			}}, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		upd, err := tc.client.UpdateWebACL(tc.ctx, &wafv2.UpdateWebACLInput{
			Name: aws.String(name), Scope: tc.scope,
			Id: aws.String(id), LockToken: aws.String(lock),
			DefaultAction: &types.DefaultAction{Allow: &types.AllowAction{}},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("fr-metric")),
			},
		})
		if err != nil {
			tc.deleteWebACL(name, id, lock)
			return fmt.Errorf("update without rules: %v", err)
		}
		nextLock := aws.ToString(upd.NextLockToken)
		defer tc.deleteWebACL(name, id, nextLock)

		got, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(name), Scope: tc.scope, Id: aws.String(id),
		})
		if err != nil {
			return fmt.Errorf("get after update: %v", err)
		}
		if len(got.WebACL.Rules) != 0 {
			return fmt.Errorf("expected rules to be cleared by full-replace update, got %d", len(got.WebACL.Rules))
		}
		if got.WebACL.Capacity != 0 {
			return fmt.Errorf("expected consumed capacity 0 after clearing rules, got %d", got.WebACL.Capacity)
		}
		return nil
	}))

	// A rule group's declared capacity is fixed at creation and WAF
	// enforces it whenever rules are added: creating rules that
	// already exceed the declared capacity must be rejected.
	results = append(results, r.RunTest("wafv2", "CreateRuleGroup_RulesExceedCapacity_Rejected", func() error {
		rgName := tc.uniqueName("over-rules-rg")
		rule := types.Rule{
			Name:     aws.String("costly-rule"),
			Priority: 1,
			Action:   &types.RuleAction{Count: &types.CountAction{}},
			Statement: &types.Statement{
				ByteMatchStatement: &types.ByteMatchStatement{
					PositionalConstraint: types.PositionalConstraintContains,
					SearchString:         []byte("/"),
					FieldToMatch:         &types.FieldToMatch{UriPath: &types.UriPath{}},
					TextTransformations: []types.TextTransformation{{
						Priority: 1, Type: types.TextTransformationTypeNone,
					}},
				},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("or-rule-metric")),
			},
		}
		_, err := tc.client.CreateRuleGroup(tc.ctx, &wafv2.CreateRuleGroupInput{
			Name:     aws.String(rgName),
			Scope:    tc.scope,
			Capacity: aws.Int64(1),
			Rules:    []types.Rule{rule, rule},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("or-rg-metric")),
			},
		})
		return AssertErrorContains(err, "WAFLimitsExceededException")
	}))

	// Updating a rule group with rules that exceed the declared
	// capacity (fixed at creation) must be rejected.
	results = append(results, r.RunTest("wafv2", "UpdateRuleGroup_RulesExceedDeclaredCapacity_Rejected", func() error {
		rgName := tc.uniqueName("over-upd-rg")
		rgID, rgLock, err := tc.createRuleGroup(rgName, 1, &types.VisibilityConfig{
			SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
			MetricName: aws.String(tc.uniqueName("ou-rg-metric")),
		})
		if err != nil {
			return fmt.Errorf("create rule group: %v", err)
		}
		defer tc.deleteRuleGroup(rgName, rgID, rgLock)

		rule := types.Rule{
			Name:     aws.String("costly-rule"),
			Priority: 1,
			Action:   &types.RuleAction{Count: &types.CountAction{}},
			Statement: &types.Statement{
				ByteMatchStatement: &types.ByteMatchStatement{
					PositionalConstraint: types.PositionalConstraintContains,
					SearchString:         []byte("/"),
					FieldToMatch:         &types.FieldToMatch{UriPath: &types.UriPath{}},
					TextTransformations: []types.TextTransformation{{
						Priority: 1, Type: types.TextTransformationTypeNone,
					}},
				},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("ou-rule-metric")),
			},
		}
		_, err = tc.client.UpdateRuleGroup(tc.ctx, &wafv2.UpdateRuleGroupInput{
			Name:      aws.String(rgName),
			Scope:     tc.scope,
			Id:        aws.String(rgID),
			LockToken: aws.String(rgLock),
			Rules:     []types.Rule{rule, rule},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String(tc.uniqueName("ou-rg-metric")),
			},
		})
		return AssertErrorContains(err, "WAFLimitsExceededException")
	}))

	return results
}
