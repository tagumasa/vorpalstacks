package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2WebACLTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	webACLName := tc.uniqueName("test-webacl")
	var webACLID, webACLARN, webACLLockToken string

	results = append(results, r.RunTest("wafv2", "ListWebACLs_Empty", func() error {
		resp, err := tc.client.ListWebACLs(tc.ctx, &wafv2.ListWebACLsInput{
			Scope: tc.scope,
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.WebACLs == nil {
			return fmt.Errorf("WebACLs is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "CreateWebACL", func() error {
		resp, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(webACLName),
			Scope: tc.scope,
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
		if err != nil {
			return err
		}
		if resp.Summary == nil {
			return fmt.Errorf("Summary is nil")
		}
		if aws.ToString(resp.Summary.Id) == "" {
			return fmt.Errorf("Summary.Id is empty")
		}
		if aws.ToString(resp.Summary.ARN) == "" {
			return fmt.Errorf("Summary.ARN is empty")
		}
		if aws.ToString(resp.Summary.LockToken) == "" {
			return fmt.Errorf("Summary.LockToken is empty")
		}
		if aws.ToString(resp.Summary.Name) != webACLName {
			return fmt.Errorf("Summary.Name mismatch: expected %s, got %s", webACLName, aws.ToString(resp.Summary.Name))
		}
		webACLID = aws.ToString(resp.Summary.Id)
		webACLARN = aws.ToString(resp.Summary.ARN)
		webACLLockToken = aws.ToString(resp.Summary.LockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACL", func() error {
		resp, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(webACLName), Scope: tc.scope, Id: aws.String(webACLID),
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
		if aws.ToString(resp.WebACL.Id) != webACLID {
			return fmt.Errorf("Id mismatch")
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
		if aws.ToString(resp.WebACL.Description) != "Test WebACL for SDK tests" {
			return fmt.Errorf("Description mismatch: got '%s'", aws.ToString(resp.WebACL.Description))
		}
		if len(resp.WebACL.Rules) != 1 {
			return fmt.Errorf("expected 1 rule, got %d", len(resp.WebACL.Rules))
		}
		if aws.ToString(resp.WebACL.Rules[0].Name) != "AllowTestRule" {
			return fmt.Errorf("expected rule name 'AllowTestRule', got '%s'", aws.ToString(resp.WebACL.Rules[0].Name))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListWebACLs_ContainsCreated", func() error {
		resp, err := tc.client.ListWebACLs(tc.ctx, &wafv2.ListWebACLsInput{
			Scope: tc.scope,
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.WebACLs {
			if s.Id != nil && *s.Id == webACLID {
				found = true
				if aws.ToString(s.Name) != webACLName {
					return fmt.Errorf("WebACL name mismatch")
				}
				if aws.ToString(s.ARN) == "" {
					return fmt.Errorf("WebACL ARN is empty")
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
		resp, err := tc.client.UpdateWebACL(tc.ctx, &wafv2.UpdateWebACLInput{
			Name: aws.String(webACLName), Scope: tc.scope,
			Id: aws.String(webACLID), LockToken: aws.String(webACLLockToken),
			DefaultAction: &types.DefaultAction{Block: &types.BlockAction{}},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String("updated-webacl-metric"),
			},
		})
		if err != nil {
			return err
		}
		if resp.NextLockToken == nil || aws.ToString(resp.NextLockToken) == "" {
			return fmt.Errorf("NextLockToken is nil or empty")
		}
		webACLLockToken = aws.ToString(resp.NextLockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UpdateWebACL_ContentVerify", func() error {
		resp, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(webACLName), Scope: tc.scope, Id: aws.String(webACLID),
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
		_, err := tc.client.UpdateWebACL(tc.ctx, &wafv2.UpdateWebACLInput{
			Name: aws.String(webACLName), Scope: tc.scope,
			Id: aws.String(webACLID), LockToken: aws.String("stale-lock-token"),
			DefaultAction: &types.DefaultAction{Allow: &types.AllowAction{}},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: false, CloudWatchMetricsEnabled: false,
				MetricName: aws.String("bad-metric"),
			},
		})
		return AssertErrorContains(err, "WAFInvalidLockTokenException")
	}))

	results = append(results, r.RunTest("wafv2", "TagResource_WebACL", func() error {
		_, err := tc.client.TagResource(tc.ctx, &wafv2.TagResourceInput{
			ResourceARN: aws.String(webACLARN),
			Tags: []types.Tag{
				{Key: aws.String("Team"), Value: aws.String("Security")},
				{Key: aws.String("Env"), Value: aws.String("Production")},
			},
		})
		if err != nil {
			return err
		}
		verifyResp, err := tc.client.ListTagsForResource(tc.ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(webACLARN),
		})
		if err != nil {
			return fmt.Errorf("ListTagsForResource after tag: %v", err)
		}
		tagMap := make(map[string]string)
		for _, t := range verifyResp.TagInfoForResource.TagList {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["Team"] != "Security" {
			return fmt.Errorf("tag 'Team' not found or wrong value")
		}
		if tagMap["Env"] != "Production" {
			return fmt.Errorf("tag 'Env' not found or wrong value")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListTagsForResource_WebACL", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &wafv2.ListTagsForResourceInput{
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
		_, err := tc.client.UntagResource(tc.ctx, &wafv2.UntagResourceInput{
			ResourceARN: aws.String(webACLARN),
			TagKeys:     []string{"Env"},
		})
		if err != nil {
			return err
		}
		verifyResp, err := tc.client.ListTagsForResource(tc.ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(webACLARN),
		})
		if err != nil {
			return fmt.Errorf("ListTagsForResource after untag: %v", err)
		}
		for _, t := range verifyResp.TagInfoForResource.TagList {
			if aws.ToString(t.Key) == "Env" {
				return fmt.Errorf("tag 'Env' should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "UntagResource_Verify", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &wafv2.ListTagsForResourceInput{
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

	fakeResourceARN := fmt.Sprintf("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/%d", time.Now().UnixNano())

	results = append(results, r.RunTest("wafv2", "AssociateWebACL", func() error {
		_, err := tc.client.AssociateWebACL(tc.ctx, &wafv2.AssociateWebACLInput{
			WebACLArn:   aws.String(webACLARN),
			ResourceArn: aws.String(fakeResourceARN),
		})
		if err != nil {
			return err
		}
		verifyResp, err := tc.client.GetWebACLForResource(tc.ctx, &wafv2.GetWebACLForResourceInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		if err != nil {
			return fmt.Errorf("GetWebACLForResource after associate: %v", err)
		}
		if verifyResp.WebACL == nil || aws.ToString(verifyResp.WebACL.Name) != webACLName {
			return fmt.Errorf("associated WebACL name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACLForResource", func() error {
		resp, err := tc.client.GetWebACLForResource(tc.ctx, &wafv2.GetWebACLForResourceInput{
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
		resp, err := tc.client.ListResourcesForWebACL(tc.ctx, &wafv2.ListResourcesForWebACLInput{
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
		_, err := tc.client.DisassociateWebACL(tc.ctx, &wafv2.DisassociateWebACLInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		if err != nil {
			return err
		}
		_, verifyErr := tc.client.GetWebACLForResource(tc.ctx, &wafv2.GetWebACLForResourceInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		if verifyErr == nil {
			return fmt.Errorf("expected error after disassociate, got nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetWebACLForResource_AfterDisassociate", func() error {
		_, err := tc.client.GetWebACLForResource(tc.ctx, &wafv2.GetWebACLForResourceInput{
			ResourceArn: aws.String(fakeResourceARN),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "DeleteWebACL", func() error {
		resp, err := tc.client.DeleteWebACL(tc.ctx, &wafv2.DeleteWebACLInput{
			Name: aws.String(webACLName), Scope: tc.scope,
			Id: aws.String(webACLID), LockToken: aws.String(webACLLockToken),
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
		_, err := tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(webACLName), Scope: tc.scope, Id: aws.String(webACLID),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	return results
}
