package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"

	"vorpalstacks-sdk-tests/config"
)

func cfWebACLCopyTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	distID, _, err := tc.createDistribution(tc.uniqueCallerRef("webacl-dist"), "WebACL test distribution", "example.net")
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "WebACL_Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}
	copyDistID := ""

	// disableAndDeleteCurrent fetches the distribution's current ETag
	// before disabling and deleting it, so callers need not track ETags.
	disableAndDeleteCurrent := func(id string) error {
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(id)})
		if err != nil {
			return err
		}
		cfgResp.DistributionConfig.Enabled = aws.Bool(false)
		updResp, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
			Id:                 aws.String(id),
			IfMatch:            cfgResp.ETag,
			DistributionConfig: cfgResp.DistributionConfig,
		})
		if err != nil {
			return err
		}
		_, err = client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
			Id:      aws.String(id),
			IfMatch: updResp.ETag,
		})
		return err
	}

	// Sweep staging distributions left behind by earlier runs so repeated
	// executions do not accumulate copies.
	if listResp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{MaxItems: aws.Int32(100)}); err == nil {
		for {
			for _, d := range listResp.DistributionList.Items {
				if d.Staging != nil && *d.Staging {
					disableAndDeleteCurrent(*d.Id)
				}
			}
			if listResp.DistributionList.NextMarker == nil || *listResp.DistributionList.NextMarker == "" {
				break
			}
			listResp, err = client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
				Marker:   listResp.DistributionList.NextMarker,
				MaxItems: aws.Int32(100),
			})
			if err != nil {
				break
			}
		}
	}

	defer func() {
		if copyDistID != "" {
			disableAndDeleteCurrent(copyDistID)
		}
		disableAndDeleteCurrent(distID)
	}()

	// A real Web ACL is needed for a successful association; the WAFv2
	// client is created locally for that purpose.
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: tc.runner.endpoint,
		Region:   tc.runner.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "WebACL_Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}
	wafClient := wafv2.NewFromConfig(cfg)
	waclName := tc.uniquePrefix("cf-wacl")
	waclResp, err := wafClient.CreateWebACL(ctx, &wafv2.CreateWebACLInput{
		Name:  aws.String(waclName),
		Scope: wafv2types.ScopeCloudfront,
		DefaultAction: &wafv2types.DefaultAction{
			Allow: &wafv2types.AllowAction{},
		},
		VisibilityConfig: &wafv2types.VisibilityConfig{
			SampledRequestsEnabled:   false,
			CloudWatchMetricsEnabled: false,
			MetricName:               aws.String(waclName),
		},
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "WebACL_Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}
	webACLArn := aws.ToString(waclResp.Summary.ARN)
	defer func() {
		wafClient.DeleteWebACL(ctx, &wafv2.DeleteWebACLInput{
			Name:      aws.String(waclName),
			Scope:     wafv2types.ScopeCloudfront,
			Id:        waclResp.Summary.Id,
			LockToken: waclResp.Summary.LockToken,
		})
	}()

	results = append(results, tc.runner.RunTest("cloudfront", "AssociateDistributionWebACL_Verify", func() error {
		resp, err := client.AssociateDistributionWebACL(ctx, &cloudfront.AssociateDistributionWebACLInput{
			Id:        aws.String(distID),
			WebACLArn: aws.String(webACLArn),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != distID {
			return fmt.Errorf("response Id mismatch: got %q", aws.ToString(resp.Id))
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("response ETag is nil or empty")
		}
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		if aws.ToString(cfgResp.DistributionConfig.WebACLId) != webACLArn {
			return fmt.Errorf("distribution WebACLId mismatch: got %q, want %q",
				aws.ToString(cfgResp.DistributionConfig.WebACLId), webACLArn)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributionsByWebACLId_Verify", func() error {
		resp, err := client.ListDistributionsByWebACLId(ctx, &cloudfront.ListDistributionsByWebACLIdInput{
			WebACLId: aws.String(webACLArn),
		})
		if err != nil {
			return err
		}
		if resp.DistributionList == nil {
			return fmt.Errorf("distribution list is nil")
		}
		found := false
		for _, d := range resp.DistributionList.Items {
			if aws.ToString(d.Id) == distID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("distribution %q not listed by Web ACL %q", distID, webACLArn)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributionsByWebACLId_MaxItemsCapped", func() error {
		resp, err := client.ListDistributionsByWebACLId(ctx, &cloudfront.ListDistributionsByWebACLIdInput{
			WebACLId: aws.String(webACLArn),
			MaxItems: aws.Int32(200),
		})
		if err != nil {
			return err
		}
		if resp.DistributionList == nil {
			return fmt.Errorf("distribution list is nil")
		}
		if aws.ToInt32(resp.DistributionList.MaxItems) != 100 {
			return fmt.Errorf("MaxItems should be capped at 100, got %d", aws.ToInt32(resp.DistributionList.MaxItems))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "AssociateDistributionWebACL_BogusArn_Rejected", func() error {
		_, err := client.AssociateDistributionWebACL(ctx, &cloudfront.AssociateDistributionWebACLInput{
			Id:        aws.String(distID),
			WebACLArn: aws.String("arn:aws:wafv2:us-east-1:000000000000:global/webacl/nonexistent/33333333-3333-4333-8333-333333333333"),
		})
		return AssertErrorContains(err, "EntityNotFound")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "DisassociateDistributionWebACL_Verify", func() error {
		resp, err := client.DisassociateDistributionWebACL(ctx, &cloudfront.DisassociateDistributionWebACLInput{
			Id: aws.String(distID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != distID {
			return fmt.Errorf("response Id mismatch: got %q", aws.ToString(resp.Id))
		}
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		if aws.ToString(cfgResp.DistributionConfig.WebACLId) != "" {
			return fmt.Errorf("distribution WebACLId should be empty, got %q",
				aws.ToString(cfgResp.DistributionConfig.WebACLId))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CopyDistribution_Verify", func() error {
		copyRef := tc.uniqueCallerRef("copy-dist")
		resp, err := client.CopyDistribution(ctx, &cloudfront.CopyDistributionInput{
			PrimaryDistributionId: aws.String(distID),
			CallerReference:       aws.String(copyRef),
			Staging:               aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Distribution == nil || resp.Distribution.Id == nil {
			return fmt.Errorf("copy distribution is nil")
		}
		copyDistID = *resp.Distribution.Id
		if copyDistID == distID {
			return fmt.Errorf("copy distribution ID equals primary ID")
		}
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(copyDistID)})
		if err != nil {
			return err
		}
		if aws.ToString(cfgResp.DistributionConfig.CallerReference) != copyRef {
			return fmt.Errorf("copy callerReference mismatch: got %q", aws.ToString(cfgResp.DistributionConfig.CallerReference))
		}
		if cfgResp.DistributionConfig.Staging == nil || !*cfgResp.DistributionConfig.Staging {
			return fmt.Errorf("copy should be a staging distribution")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CopyDistribution_StagingRenderedConsistently", func() error {
		if copyDistID == "" {
			return fmt.Errorf("no copy distribution from the previous test")
		}
		var marker *string
		found := false
		for {
			listResp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{Marker: marker})
			if err != nil {
				return err
			}
			if listResp.DistributionList == nil {
				break
			}
			for _, d := range listResp.DistributionList.Items {
				if aws.ToString(d.Id) == copyDistID {
					found = true
					if d.Staging == nil || !*d.Staging {
						return fmt.Errorf("summary Staging flag is false for the staging copy")
					}
				}
			}
			if listResp.DistributionList.NextMarker == nil || !*listResp.DistributionList.IsTruncated {
				break
			}
			marker = listResp.DistributionList.NextMarker
		}
		if !found {
			return fmt.Errorf("staging copy %q not found in ListDistributions", copyDistID)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CopyDistribution_DuplicateCallerReference_Rejected", func() error {
		cfgResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		primaryRef := aws.ToString(cfgResp.DistributionConfig.CallerReference)
		primaryETag := aws.ToString(cfgResp.ETag)

		_, err = client.CopyDistribution(ctx, &cloudfront.CopyDistributionInput{
			PrimaryDistributionId: aws.String(distID),
			CallerReference:       aws.String(primaryRef),
			Staging:               aws.Bool(true),
		})
		if err := AssertErrorContains(err, "DistributionAlreadyExists"); err != nil {
			return err
		}

		// The rejected copy must leave the primary untouched.
		afterResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		if aws.ToString(afterResp.ETag) != primaryETag {
			return fmt.Errorf("primary ETag changed after rejected copy: got %q, want %q",
				aws.ToString(afterResp.ETag), primaryETag)
		}
		if afterResp.DistributionConfig.Staging != nil && *afterResp.DistributionConfig.Staging {
			return fmt.Errorf("primary was converted to a staging distribution by the rejected copy")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CopyDistribution_StagingFalse_Rejected", func() error {
		_, err := client.CopyDistribution(ctx, &cloudfront.CopyDistributionInput{
			PrimaryDistributionId: aws.String(distID),
			CallerReference:       aws.String(tc.uniqueCallerRef("copy-bad")),
			Staging:               aws.Bool(false),
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	return results
}
