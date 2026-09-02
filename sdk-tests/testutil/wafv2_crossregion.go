package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

// runWAFv2CrossRegionTests pins the storage partition of the CloudFront
// scope: the AWS WAFv2 API reference requires the us-east-1 endpoint for
// that scope and its ARNs carry us-east-1 whatever region the call was
// signed for, so a CLOUDFRONT-scope web ACL created from another region
// must live in the partition its ARN names — reads, listings and tags
// from any region resolve the same resource.
func (r *TestRunner) runWAFv2CrossRegionTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("wafv2", "CrossRegion_CloudFrontScopePartition", func() error {
		// A second client signed for a non-us-east-1 region against the
		// same endpoint; the server derives the request region from the
		// credential scope.
		euWest, err := newWAFv2TestContext(r.endpoint, "eu-west-1", r.accountID)
		if err != nil {
			return err
		}
		name := euWest.uniqueName("test-webacl-xregion")
		created, err := euWest.client.CreateWebACL(euWest.ctx, &wafv2.CreateWebACLInput{
			Name:  aws.String(name),
			Scope: types.ScopeCloudfront,
			DefaultAction: &types.DefaultAction{
				Allow: &types.AllowAction{},
			},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: true,
				MetricName:               aws.String("xregion-webacl-metric"),
			},
		})
		if err != nil {
			return err
		}
		id := aws.ToString(created.Summary.Id)
		arn := aws.ToString(created.Summary.ARN)
		lockToken := aws.ToString(created.Summary.LockToken)
		deleted := false
		defer func() {
			if !deleted {
				_, _ = euWest.client.DeleteWebACL(euWest.ctx, &wafv2.DeleteWebACLInput{
					Name: aws.String(name), Scope: types.ScopeCloudfront,
					Id: aws.String(id), LockToken: aws.String(lockToken),
				})
			}
		}()

		// The ARN of the global scope is pinned to us-east-1.
		wantPrefix := fmt.Sprintf("arn:aws:wafv2:us-east-1:%s:global/webacl/%s/", r.accountID, name)
		if !strings.HasPrefix(arn, wantPrefix) {
			return fmt.Errorf("ARN = %s, want prefix %s", arn, wantPrefix)
		}

		// Reading by Id and Scope from the creating region routes to the
		// us-east-1 partition.
		got, err := euWest.client.GetWebACL(euWest.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(name), Scope: types.ScopeCloudfront, Id: aws.String(id),
		})
		if err != nil {
			return fmt.Errorf("GetWebACL from eu-west-1: %w", err)
		}
		if aws.ToString(got.WebACL.ARN) != arn {
			return fmt.Errorf("GetWebACL ARN = %s, want %s", aws.ToString(got.WebACL.ARN), arn)
		}

		// The listing from the runner's own region sees the same
		// partition: the resource is stored where its ARN points, not
		// where the create call was signed. Traverse every page, as
		// parallel tests add resources of their own.
		found := false
		nextMarker := ""
		for {
			listInput := &wafv2.ListWebACLsInput{
				Scope: types.ScopeCloudfront, Limit: aws.Int32(100),
			}
			if nextMarker != "" {
				listInput.NextMarker = aws.String(nextMarker)
			}
			listed, err := tc.client.ListWebACLs(tc.ctx, listInput)
			if err != nil {
				return fmt.Errorf("ListWebACLs from %s: %w", tc.region, err)
			}
			for _, summary := range listed.WebACLs {
				if aws.ToString(summary.ARN) == arn {
					found = true
				}
			}
			if found || listed.NextMarker == nil || *listed.NextMarker == "" {
				break
			}
			nextMarker = *listed.NextMarker
		}
		if !found {
			return fmt.Errorf("web ACL created from eu-west-1 not listed from %s", tc.region)
		}

		// Tags route by the resource ARN's region too: written from
		// eu-west-1, listed from the runner's region.
		if _, err := euWest.client.TagResource(euWest.ctx, &wafv2.TagResourceInput{
			ResourceARN: aws.String(arn),
			Tags:        []types.Tag{{Key: aws.String("CrossRegion"), Value: aws.String("pinned")}},
		}); err != nil {
			return fmt.Errorf("TagResource from eu-west-1: %w", err)
		}
		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &wafv2.ListTagsForResourceInput{
			ResourceARN: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("ListTagsForResource from %s: %w", tc.region, err)
		}
		tagFound := false
		for _, tag := range tagResp.TagInfoForResource.TagList {
			if aws.ToString(tag.Key) == "CrossRegion" && aws.ToString(tag.Value) == "pinned" {
				tagFound = true
			}
		}
		if !tagFound {
			return fmt.Errorf("tag written from eu-west-1 not visible from %s", tc.region)
		}

		// Deleting from the creating region removes it from the shared
		// partition for every region.
		if _, err := euWest.client.DeleteWebACL(euWest.ctx, &wafv2.DeleteWebACLInput{
			Name: aws.String(name), Scope: types.ScopeCloudfront,
			Id: aws.String(id), LockToken: aws.String(lockToken),
		}); err != nil {
			return fmt.Errorf("DeleteWebACL from eu-west-1: %w", err)
		}
		deleted = true
		_, err = tc.client.GetWebACL(tc.ctx, &wafv2.GetWebACLInput{
			Name: aws.String(name), Scope: types.ScopeCloudfront, Id: aws.String(id),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	return results
}
