package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfTagTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	callerRef := tc.uniqueCallerRef("test-cftags")
	distID, distETag, err := tc.createDistributionWithTags(callerRef, "Tagged distribution", "example.org",
		[]types.Tag{
			{Key: aws.String("Environment"), Value: aws.String("test")},
			{Key: aws.String("Owner"), Value: aws.String("sdk-tests")},
		})
	if err != nil {
		results = append(results, TestResult{
			Service:  "cloudfront",
			TestName: "CreateDistributionWithTags",
			Status:   "FAIL",
			Error:    err.Error(),
		})
		return results
	}
	results = append(results, TestResult{
		Service:  "cloudfront",
		TestName: "CreateDistributionWithTags",
		Status:   "PASS",
	})

	var distARN string
	getResp, _ := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: aws.String(distID),
	})
	if getResp != nil && getResp.Distribution != nil {
		distARN = aws.ToString(getResp.Distribution.ARN)
	}

	if distARN != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "ListTagsForResource_AfterCreate", func() error {
			resp, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: aws.String(distARN),
			})
			if err != nil {
				return err
			}
			if resp.Tags == nil {
				return fmt.Errorf("tags is nil")
			}
			if len(resp.Tags.Items) < 2 {
				return fmt.Errorf("expected at least 2 tags, got %d", len(resp.Tags.Items))
			}
			tagMap := make(map[string]string)
			for _, t := range resp.Tags.Items {
				tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
			}
			if tagMap["Environment"] != "test" {
				return fmt.Errorf("Environment tag mismatch: got %q", tagMap["Environment"])
			}
			if tagMap["Owner"] != "sdk-tests" {
				return fmt.Errorf("Owner tag mismatch: got %q", tagMap["Owner"])
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "TagResource_Verify", func() error {
			_, err := client.TagResource(ctx, &cloudfront.TagResourceInput{
				Resource: aws.String(distARN),
				Tags: &types.Tags{
					Items: []types.Tag{
						{Key: aws.String("ExtraTag"), Value: aws.String("extra-value")},
					},
				},
			})
			if err != nil {
				return err
			}
			resp, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: aws.String(distARN),
			})
			if err != nil {
				return err
			}
			tagMap := make(map[string]string)
			for _, t := range resp.Tags.Items {
				tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
			}
			if tagMap["ExtraTag"] != "extra-value" {
				return fmt.Errorf("ExtraTag not found or mismatch: got %q", tagMap["ExtraTag"])
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "UntagResource_Verify", func() error {
			_, err := client.UntagResource(ctx, &cloudfront.UntagResourceInput{
				Resource: aws.String(distARN),
				TagKeys: &types.TagKeys{
					Items: []string{"ExtraTag"},
				},
			})
			if err != nil {
				return err
			}
			resp, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: aws.String(distARN),
			})
			if err != nil {
				return err
			}
			for _, t := range resp.Tags.Items {
				if aws.ToString(t.Key) == "ExtraTag" {
					return fmt.Errorf("ExtraTag should have been removed")
				}
			}
			if len(resp.Tags.Items) < 2 {
				return fmt.Errorf("expected at least 2 remaining tags, got %d", len(resp.Tags.Items))
			}
			return nil
		}))
	}

	results = append(results, tc.runner.RunTest("cloudfront", "Cleanup_TaggedDistribution", func() error {
		return tc.disableAndDeleteDistribution(distID, distETag)
	}))

	return results
}

func (tc *cfTestContext) createDistributionWithTags(callerRef, comment, originDomain string, tags []types.Tag) (distID, distETag string, err error) {
	originID := tc.uniquePrefix("origin")
	resp, rerr := tc.client.CreateDistributionWithTags(tc.ctx, &cloudfront.CreateDistributionWithTagsInput{
		DistributionConfigWithTags: &types.DistributionConfigWithTags{
			DistributionConfig: &types.DistributionConfig{
				CallerReference: aws.String(callerRef),
				Enabled:         aws.Bool(true),
				Comment:         aws.String(comment),
				Origins: &types.Origins{
					Quantity: aws.Int32(1),
					Items: []types.Origin{
						{
							Id:         aws.String(originID),
							DomainName: aws.String(originDomain),
							CustomOriginConfig: &types.CustomOriginConfig{
								HTTPPort:             aws.Int32(80),
								HTTPSPort:            aws.Int32(443),
								OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
							},
						},
					},
				},
				DefaultCacheBehavior: &types.DefaultCacheBehavior{
					TargetOriginId:       aws.String(originID),
					ViewerProtocolPolicy: types.ViewerProtocolPolicyAllowAll,
					ForwardedValues: &types.ForwardedValues{
						QueryString: aws.Bool(false),
						Cookies: &types.CookiePreference{
							Forward: types.ItemSelectionNone,
						},
					},
				},
				ViewerCertificate: &types.ViewerCertificate{
					CloudFrontDefaultCertificate: aws.Bool(true),
				},
				Restrictions: &types.Restrictions{
					GeoRestriction: &types.GeoRestriction{
						RestrictionType: types.GeoRestrictionTypeNone,
						Quantity:        aws.Int32(0),
					},
				},
			},
			Tags: &types.Tags{
				Items: tags,
			},
		},
	})
	if rerr != nil {
		return "", "", rerr
	}
	return aws.ToString(resp.Distribution.Id), aws.ToString(resp.ETag), nil
}
