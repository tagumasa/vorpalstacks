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

type wafv2TestContext struct {
	client *wafv2.Client
	ctx    context.Context
	scope  types.Scope
}

func newWAFv2TestContext(endpoint, region string) (*wafv2TestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: endpoint,
		Region:   region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &wafv2TestContext{
		client: wafv2.NewFromConfig(cfg),
		ctx:    context.Background(),
		scope:  types.ScopeCloudfront,
	}, nil
}

func (tc *wafv2TestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *wafv2TestContext) createIPSet(name string, tags []types.Tag) (id, arn, lockToken string, err error) {
	resp, err := tc.client.CreateIPSet(tc.ctx, &wafv2.CreateIPSetInput{
		Name:             aws.String(name),
		Scope:            tc.scope,
		IPAddressVersion: types.IPAddressVersionIpv4,
		Addresses:        []string{},
		Tags:             tags,
	})
	if err != nil {
		return "", "", "", err
	}
	return aws.ToString(resp.Summary.Id), aws.ToString(resp.Summary.ARN), aws.ToString(resp.Summary.LockToken), nil
}

func (tc *wafv2TestContext) deleteIPSet(name, id, lockToken string) error {
	_, err := tc.client.DeleteIPSet(tc.ctx, &wafv2.DeleteIPSetInput{
		Name: aws.String(name), Scope: tc.scope,
		Id: aws.String(id), LockToken: aws.String(lockToken),
	})
	return err
}

func (tc *wafv2TestContext) createRegexPatternSet(name string, patterns []types.Regex) (id, lockToken string, err error) {
	resp, err := tc.client.CreateRegexPatternSet(tc.ctx, &wafv2.CreateRegexPatternSetInput{
		Name:                  aws.String(name),
		Scope:                 tc.scope,
		RegularExpressionList: patterns,
	})
	if err != nil {
		return "", "", err
	}
	return aws.ToString(resp.Summary.Id), aws.ToString(resp.Summary.LockToken), nil
}

func (tc *wafv2TestContext) deleteRegexPatternSet(name, id, lockToken string) error {
	_, err := tc.client.DeleteRegexPatternSet(tc.ctx, &wafv2.DeleteRegexPatternSetInput{
		Name: aws.String(name), Scope: tc.scope,
		Id: aws.String(id), LockToken: aws.String(lockToken),
	})
	return err
}

func (tc *wafv2TestContext) createRuleGroup(name string, capacity int64, vc *types.VisibilityConfig) (id, lockToken string, err error) {
	resp, err := tc.client.CreateRuleGroup(tc.ctx, &wafv2.CreateRuleGroupInput{
		Name:             aws.String(name),
		Scope:            tc.scope,
		Capacity:         aws.Int64(capacity),
		VisibilityConfig: vc,
	})
	if err != nil {
		return "", "", err
	}
	return aws.ToString(resp.Summary.Id), aws.ToString(resp.Summary.LockToken), nil
}

func (tc *wafv2TestContext) deleteRuleGroup(name, id, lockToken string) error {
	_, err := tc.client.DeleteRuleGroup(tc.ctx, &wafv2.DeleteRuleGroupInput{
		Name: aws.String(name), Scope: tc.scope,
		Id: aws.String(id), LockToken: aws.String(lockToken),
	})
	return err
}

func (tc *wafv2TestContext) createWebACL(name string, defaultAction *types.DefaultAction, vc *types.VisibilityConfig, rules []types.Rule, tags []types.Tag) (id, arn, lockToken string, err error) {
	resp, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
		Name:             aws.String(name),
		Scope:            tc.scope,
		DefaultAction:    defaultAction,
		VisibilityConfig: vc,
		Rules:            rules,
		Tags:             tags,
	})
	if err != nil {
		return "", "", "", err
	}
	return aws.ToString(resp.Summary.Id), aws.ToString(resp.Summary.ARN), aws.ToString(resp.Summary.LockToken), nil
}

func (tc *wafv2TestContext) deleteWebACL(name, id, lockToken string) error {
	_, err := tc.client.DeleteWebACL(tc.ctx, &wafv2.DeleteWebACLInput{
		Name: aws.String(name), Scope: tc.scope,
		Id: aws.String(id), LockToken: aws.String(lockToken),
	})
	return err
}

func (r *TestRunner) RunWAFv2Tests() []TestResult {
	var results []TestResult

	tc, err := newWAFv2TestContext(r.endpoint, r.region)
	if err != nil {
		return append(results, TestResult{
			Service:  "wafv2",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}

	results = append(results, r.runWAFv2IPSetTests(tc)...)
	results = append(results, r.runWAFv2RegexPatternSetTests(tc)...)
	results = append(results, r.runWAFv2RuleGroupTests(tc)...)
	results = append(results, r.runWAFv2WebACLTests(tc)...)
	results = append(results, r.runWAFv2LoggingTests(tc)...)
	results = append(results, r.runWAFv2ManagedRulesTests(tc)...)
	results = append(results, r.runWAFv2EdgeTests(tc)...)

	return results
}
