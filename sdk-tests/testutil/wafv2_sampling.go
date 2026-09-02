package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	waftypes "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2SamplingTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	// A WebACL whose rule has sampling enabled but receives no traffic:
	// the sample query answers an empty list, not an error.
	results = append(results, r.RunTest("wafv2", "GetSampledRequests_Empty", func() error {
		name := tc.uniqueName("sample-empty-acl")
		_, arn, lock, err := tc.createWebACL(name,
			&waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
			&waftypes.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
				MetricName: aws.String(tc.uniqueName("se-acl-metric")),
			},
			[]waftypes.Rule{{
				Name:     aws.String("BlockAll"),
				Priority: 1,
				Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
				Statement: &waftypes.Statement{
					ByteMatchStatement: &waftypes.ByteMatchStatement{
						FieldToMatch:         &waftypes.FieldToMatch{UriPath: &waftypes.UriPath{}},
						PositionalConstraint: waftypes.PositionalConstraintContains,
						SearchString:         []byte("/"),
						TextTransformations: []waftypes.TextTransformation{
							{Priority: 0, Type: waftypes.TextTransformationTypeNone},
						},
					},
				},
				VisibilityConfig: &waftypes.VisibilityConfig{
					SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
					MetricName: aws.String("block-all-metric"),
				},
			}}, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteWebACL(name, arn[strings.LastIndex(arn, "/")+1:], lock)

		now := time.Now()
		resp, err := tc.client.GetSampledRequests(tc.ctx, &wafv2.GetSampledRequestsInput{
			WebAclArn:      aws.String(arn),
			RuleMetricName: aws.String("block-all-metric"),
			Scope:          tc.scope,
			TimeWindow: &waftypes.TimeWindow{
				StartTime: aws.Time(now.Add(-time.Hour)),
				EndTime:   aws.Time(now),
			},
			MaxItems: aws.Int64(100),
		})
		if err != nil {
			return err
		}
		if len(resp.SampledRequests) != 0 {
			return fmt.Errorf("expected no sampled requests, got %d", len(resp.SampledRequests))
		}
		if resp.PopulationSize != 0 {
			return fmt.Errorf("expected PopulationSize 0, got %d", resp.PopulationSize)
		}
		if resp.TimeWindow == nil || resp.TimeWindow.StartTime == nil || resp.TimeWindow.EndTime == nil {
			return fmt.Errorf("expected the echoed TimeWindow")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetSampledRequests_NonexistentWebACL", func() error {
		now := time.Now()
		_, err := tc.client.GetSampledRequests(tc.ctx, &wafv2.GetSampledRequestsInput{
			WebAclArn:      aws.String("arn:aws:wafv2:" + tc.region + ":" + tc.accountID + ":regional/webacl/nosuch/" + tc.uniqueName("nosuch")),
			RuleMetricName: aws.String("metric"),
			Scope:          tc.scope,
			TimeWindow: &waftypes.TimeWindow{
				StartTime: aws.Time(now.Add(-time.Minute)),
				EndTime:   aws.Time(now),
			},
			MaxItems: aws.Int64(100),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "GetSampledRequests_ReversedTimeWindow", func() error {
		name := tc.uniqueName("sample-window-acl")
		_, arn, lock, err := tc.createWebACL(name,
			&waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
			&waftypes.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
				MetricName: aws.String(tc.uniqueName("sw-acl-metric")),
			}, nil, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteWebACL(name, arn[strings.LastIndex(arn, "/")+1:], lock)

		now := time.Now()
		_, err = tc.client.GetSampledRequests(tc.ctx, &wafv2.GetSampledRequestsInput{
			WebAclArn:      aws.String(arn),
			RuleMetricName: aws.String("metric"),
			Scope:          tc.scope,
			TimeWindow: &waftypes.TimeWindow{
				StartTime: aws.Time(now),
				EndTime:   aws.Time(now.Add(-time.Minute)),
			},
			MaxItems: aws.Int64(100),
		})
		return AssertErrorContains(err, "WAFInvalidParameterException")
	}))

	// Managed keys on a rate rule that aggregates on a constant key is
	// the documented unsupported-aggregate rejection.
	results = append(results, r.RunTest("wafv2", "GetRateBasedStatementManagedKeys_UnsupportedAggregateKey", func() error {
		name := tc.uniqueName("mk-constant-acl")
		id, _, lock, err := tc.createWebACL(name,
			&waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
			&waftypes.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
				MetricName: aws.String(tc.uniqueName("mkc-acl-metric")),
			},
			[]waftypes.Rule{{
				Name:     aws.String("ConstantRate"),
				Priority: 1,
				Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
				Statement: &waftypes.Statement{
					RateBasedStatement: &waftypes.RateBasedStatement{
						Limit:            aws.Int64(100),
						AggregateKeyType: waftypes.RateBasedStatementAggregateKeyTypeConstant,
					},
				},
				VisibilityConfig: &waftypes.VisibilityConfig{
					SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
					MetricName: aws.String("constant-rate-metric"),
				},
			}}, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteWebACL(name, id, lock)

		_, err = tc.client.GetRateBasedStatementManagedKeys(tc.ctx, &wafv2.GetRateBasedStatementManagedKeysInput{
			Scope:      tc.scope,
			WebACLName: aws.String(name),
			WebACLId:   aws.String(id),
			RuleName:   aws.String("ConstantRate"),
		})
		return AssertErrorContains(err, "WAFUnsupportedAggregateKeyTypeException")
	}))

	results = append(results, r.RunTest("wafv2", "GetRateBasedStatementManagedKeys_NonexistentRule", func() error {
		name := tc.uniqueName("mk-noRule-acl")
		id, _, lock, err := tc.createWebACL(name,
			&waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
			&waftypes.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
				MetricName: aws.String(tc.uniqueName("mkn-acl-metric")),
			}, nil, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteWebACL(name, id, lock)

		_, err = tc.client.GetRateBasedStatementManagedKeys(tc.ctx, &wafv2.GetRateBasedStatementManagedKeysInput{
			Scope:      tc.scope,
			WebACLName: aws.String(name),
			WebACLId:   aws.String(id),
			RuleName:   aws.String("NoSuchRateRule"),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	// An IP-aggregating rate rule with no traffic has no tracked keys.
	results = append(results, r.RunTest("wafv2", "GetRateBasedStatementManagedKeys_NoTrackedKeys", func() error {
		name := tc.uniqueName("mk-empty-acl")
		id, _, lock, err := tc.createWebACL(name,
			&waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
			&waftypes.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
				MetricName: aws.String(tc.uniqueName("mke-acl-metric")),
			},
			[]waftypes.Rule{{
				Name:     aws.String("IpRate"),
				Priority: 1,
				Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
				Statement: &waftypes.Statement{
					RateBasedStatement: &waftypes.RateBasedStatement{
						Limit:            aws.Int64(100),
						AggregateKeyType: waftypes.RateBasedStatementAggregateKeyTypeIp,
					},
				},
				VisibilityConfig: &waftypes.VisibilityConfig{
					SampledRequestsEnabled: true, CloudWatchMetricsEnabled: false,
					MetricName: aws.String("ip-rate-metric"),
				},
			}}, nil)
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteWebACL(name, id, lock)

		resp, err := tc.client.GetRateBasedStatementManagedKeys(tc.ctx, &wafv2.GetRateBasedStatementManagedKeysInput{
			Scope:      tc.scope,
			WebACLName: aws.String(name),
			WebACLId:   aws.String(id),
			RuleName:   aws.String("IpRate"),
		})
		if err != nil {
			return err
		}
		if resp.ManagedKeysIPV4 != nil && len(resp.ManagedKeysIPV4.Addresses) != 0 {
			return fmt.Errorf("expected no tracked IPv4 keys, got %v", resp.ManagedKeysIPV4.Addresses)
		}
		if resp.ManagedKeysIPV6 != nil && len(resp.ManagedKeysIPV6.Addresses) != 0 {
			return fmt.Errorf("expected no tracked IPv6 keys, got %v", resp.ManagedKeysIPV6.Addresses)
		}
		return nil
	}))

	return results
}
