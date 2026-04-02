package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCloudFrontTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cloudfront.NewFromConfig(cfg)
	ctx := context.Background()

	// Test 1: List Distributions (empty)
	results = append(results, r.RunTest("cloudfront", "ListDistributions", func() error {
		resp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.DistributionList == nil {
			return fmt.Errorf("distribution list is nil")
		}
		return nil
	}))

	// Test 2: Create Distribution
	var distID string
	var distETag string
	callerRef := fmt.Sprintf("test-cf-%d", time.Now().UnixNano())
	originID := fmt.Sprintf("test-origin-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("cloudfront", "CreateDistribution", func() error {
		resp, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{
			DistributionConfig: &types.DistributionConfig{
				CallerReference:   aws.String(callerRef),
				Enabled:           aws.Bool(true),
				Comment:           aws.String("SDK test distribution"),
				DefaultRootObject: aws.String("index.html"),
				Origins: &types.Origins{
					Quantity: aws.Int32(1),
					Items: []types.Origin{
						{
							Id:         aws.String(originID),
							DomainName: aws.String("example.com"),
							CustomOriginConfig: &types.CustomOriginConfig{
								HTTPPort:               aws.Int32(80),
								HTTPSPort:              aws.Int32(443),
								OriginProtocolPolicy:   types.OriginProtocolPolicyHttpOnly,
								OriginReadTimeout:      aws.Int32(30),
								OriginKeepaliveTimeout: aws.Int32(5),
								OriginSslProtocols: &types.OriginSslProtocols{
									Quantity: aws.Int32(1),
									Items: []types.SslProtocol{
										types.SslProtocolTLSv12,
									},
								},
							},
						},
					},
				},
				DefaultCacheBehavior: &types.DefaultCacheBehavior{
					TargetOriginId:       aws.String(originID),
					ViewerProtocolPolicy: types.ViewerProtocolPolicyAllowAll,
					AllowedMethods: &types.AllowedMethods{
						Quantity: aws.Int32(2),
						Items: []types.Method{
							types.MethodHead,
							types.MethodGet,
						},
						CachedMethods: &types.CachedMethods{
							Quantity: aws.Int32(2),
							Items: []types.Method{
								types.MethodHead,
								types.MethodGet,
							},
						},
					},
					ForwardedValues: &types.ForwardedValues{
						QueryString: aws.Bool(false),
						Cookies: &types.CookiePreference{
							Forward: types.ItemSelectionNone,
						},
					},
					MinTTL:     aws.Int64(0),
					DefaultTTL: aws.Int64(3600),
					MaxTTL:     aws.Int64(86400),
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
		})
		if err == nil {
			distID = aws.ToString(resp.Distribution.Id)
			distETag = aws.ToString(resp.ETag)
		}
		return err
	}))

	// Test 3: Get Distribution
	if distID != "" {
		results = append(results, r.RunTest("cloudfront", "GetDistribution", func() error {
			resp, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
				Id: aws.String(distID),
			})
			if err != nil {
				return err
			}
			if resp.Distribution == nil {
				return fmt.Errorf("distribution is nil")
			}
			if resp.ETag == nil || *resp.ETag == "" {
				return fmt.Errorf("ETag header is missing")
			}
			return nil
		}))

		// Test 4: Get Distribution Config
		results = append(results, r.RunTest("cloudfront", "GetDistributionConfig", func() error {
			resp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
				Id: aws.String(distID),
			})
			if err != nil {
				return err
			}
			if resp.DistributionConfig == nil {
				return fmt.Errorf("distribution config is nil")
			}
			return nil
		}))

		// Test 5: List Distributions (should contain the created one)
		results = append(results, r.RunTest("cloudfront", "ListDistributionsAfterCreate", func() error {
			resp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
				MaxItems: aws.Int32(10),
			})
			if err != nil {
				return err
			}
			if resp.DistributionList == nil || resp.DistributionList.Quantity == nil || *resp.DistributionList.Quantity < 1 {
				return fmt.Errorf("expected at least 1 distribution, got 0")
			}
			return nil
		}))

		// Test 6: Update Distribution (disable before delete)
		var updateETag string
		results = append(results, r.RunTest("cloudfront", "UpdateDistribution", func() error {
			getResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
				Id: aws.String(distID),
			})
			if err != nil {
				return err
			}
			getResp.DistributionConfig.Enabled = aws.Bool(false)
			resp, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
				Id:                 aws.String(distID),
				IfMatch:            aws.String(distETag),
				DistributionConfig: getResp.DistributionConfig,
			})
			if err == nil {
				updateETag = aws.ToString(resp.ETag)
			}
			return err
		}))

		// Test 7: Delete Distribution
		if updateETag != "" {
			results = append(results, r.RunTest("cloudfront", "DeleteDistribution", func() error {
				resp, err := client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
					Id:      aws.String(distID),
					IfMatch: aws.String(updateETag),
				})
				if err != nil {
					return err
				}
				if resp == nil {
					return fmt.Errorf("response is nil")
				}
				return nil
			}))
		}

		// Test 8: Get Distribution after delete (should fail with 404)
		results = append(results, r.RunTest("cloudfront", "GetDistributionAfterDelete", func() error {
			_, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
				Id: aws.String(distID),
			})
			if err == nil {
				return fmt.Errorf("expected error for deleted distribution")
			}
			return nil
		}))
	}

	// Test 9: List Distributions By WebACLId
	results = append(results, r.RunTest("cloudfront", "ListDistributionsByWebACLId", func() error {
		resp, err := client.ListDistributionsByWebACLId(ctx, &cloudfront.ListDistributionsByWebACLIdInput{
			WebACLId: aws.String("12345678-1234-1234-1234-123456789012"),
		})
		if err != nil {
			return err
		}
		if resp.DistributionList == nil {
			return fmt.Errorf("distribution list is nil")
		}
		return nil
	}))

	// --- CreateDistributionWithTags ---
	var taggedDistID string
	var taggedDistETag string
	var taggedDistARN string
	callerRef2 := fmt.Sprintf("test-cftags-%d", time.Now().UnixNano())
	originID2 := fmt.Sprintf("test-origin2-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("cloudfront", "CreateDistributionWithTags", func() error {
		resp, err := client.CreateDistributionWithTags(ctx, &cloudfront.CreateDistributionWithTagsInput{
			DistributionConfigWithTags: &types.DistributionConfigWithTags{
				DistributionConfig: &types.DistributionConfig{
					CallerReference: aws.String(callerRef2),
					Enabled:         aws.Bool(true),
					Comment:         aws.String("Tagged distribution"),
					Origins: &types.Origins{
						Quantity: aws.Int32(1),
						Items: []types.Origin{
							{
								Id:         aws.String(originID2),
								DomainName: aws.String("example.org"),
								CustomOriginConfig: &types.CustomOriginConfig{
									HTTPPort:             aws.Int32(80),
									HTTPSPort:            aws.Int32(443),
									OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
								},
							},
						},
					},
					DefaultCacheBehavior: &types.DefaultCacheBehavior{
						TargetOriginId:       aws.String(originID2),
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
					Items: []types.Tag{
						{Key: aws.String("Environment"), Value: aws.String("test")},
						{Key: aws.String("Owner"), Value: aws.String("sdk-tests")},
					},
				},
			},
		})
		if err == nil {
			taggedDistID = aws.ToString(resp.Distribution.Id)
			taggedDistETag = aws.ToString(resp.ETag)
			taggedDistARN = aws.ToString(resp.Distribution.ARN)
		}
		return err
	}))

	// --- Tag operations on the tagged distribution ---
	if taggedDistARN != "" {
		results = append(results, r.RunTest("cloudfront", "ListTagsForResource_Distribution", func() error {
			resp, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: aws.String(taggedDistARN),
			})
			if err != nil {
				return err
			}
			if resp.Tags == nil {
				return fmt.Errorf("tags is nil")
			}
			if resp.Tags.Items == nil || len(resp.Tags.Items) < 2 {
				return fmt.Errorf("expected at least 2 tags, got %d", len(resp.Tags.Items))
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "TagResource_Distribution", func() error {
			_, err := client.TagResource(ctx, &cloudfront.TagResourceInput{
				Resource: aws.String(taggedDistARN),
				Tags: &types.Tags{
					Items: []types.Tag{
						{Key: aws.String("ExtraTag"), Value: aws.String("extra-value")},
					},
				},
			})
			return err
		}))

		results = append(results, r.RunTest("cloudfront", "ListTagsForResource_AfterTag", func() error {
			resp, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: aws.String(taggedDistARN),
			})
			if err != nil {
				return err
			}
			if resp.Tags == nil || len(resp.Tags.Items) < 3 {
				return fmt.Errorf("expected at least 3 tags after TagResource, got %d", len(resp.Tags.Items))
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "UntagResource_Distribution", func() error {
			_, err := client.UntagResource(ctx, &cloudfront.UntagResourceInput{
				Resource: aws.String(taggedDistARN),
				TagKeys: &types.TagKeys{
					Items: []string{"ExtraTag"},
				},
			})
			return err
		}))

		results = append(results, r.RunTest("cloudfront", "ListTagsForResource_AfterUntag", func() error {
			resp, err := client.ListTagsForResource(ctx, &cloudfront.ListTagsForResourceInput{
				Resource: aws.String(taggedDistARN),
			})
			if err != nil {
				return err
			}
			if resp.Tags == nil || len(resp.Tags.Items) < 2 {
				return fmt.Errorf("expected at least 2 tags after UntagResource, got %d", len(resp.Tags.Items))
			}
			for _, t := range resp.Tags.Items {
				if *t.Key == "ExtraTag" {
					return fmt.Errorf("ExtraTag should have been removed")
				}
			}
			return nil
		}))
	}

	// Clean up tagged distribution
	if taggedDistID != "" && taggedDistETag != "" {
		results = append(results, r.RunTest("cloudfront", "Cleanup_TaggedDistribution", func() error {
			getResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
				Id: aws.String(taggedDistID),
			})
			if err != nil {
				return err
			}
			getResp.DistributionConfig.Enabled = aws.Bool(false)
			updateResp, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
				Id:                 aws.String(taggedDistID),
				IfMatch:            aws.String(taggedDistETag),
				DistributionConfig: getResp.DistributionConfig,
			})
			if err != nil {
				return err
			}
			_, err = client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
				Id:      aws.String(taggedDistID),
				IfMatch: updateResp.ETag,
			})
			return err
		}))
	}

	// --- Invalidation tests ---
	// Create a fresh distribution for invalidation tests
	var invDistID string
	callerRef3 := fmt.Sprintf("test-inv-%d", time.Now().UnixNano())
	originID3 := fmt.Sprintf("test-origin3-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("cloudfront", "CreateDistribution_ForInvalidation", func() error {
		resp, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{
			DistributionConfig: &types.DistributionConfig{
				CallerReference: aws.String(callerRef3),
				Enabled:         aws.Bool(true),
				Comment:         aws.String("For invalidation tests"),
				Origins: &types.Origins{
					Quantity: aws.Int32(1),
					Items: []types.Origin{
						{
							Id:         aws.String(originID3),
							DomainName: aws.String("inv-test.example.com"),
							CustomOriginConfig: &types.CustomOriginConfig{
								HTTPPort:             aws.Int32(80),
								HTTPSPort:            aws.Int32(443),
								OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
							},
						},
					},
				},
				DefaultCacheBehavior: &types.DefaultCacheBehavior{
					TargetOriginId:       aws.String(originID3),
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
		})
		if err == nil {
			invDistID = aws.ToString(resp.Distribution.Id)
		}
		return err
	}))

	var invID string
	if invDistID != "" {
		results = append(results, r.RunTest("cloudfront", "CreateInvalidation", func() error {
			resp, err := client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
				DistributionId: aws.String(invDistID),
				InvalidationBatch: &types.InvalidationBatch{
					CallerReference: aws.String(fmt.Sprintf("inv-ref-%d", time.Now().UnixNano())),
					Paths: &types.Paths{
						Quantity: aws.Int32(2),
						Items: []string{
							"/index.html",
							"/images/*",
						},
					},
				},
			})
			if err == nil {
				invID = aws.ToString(resp.Invalidation.Id)
			}
			return err
		}))

		if invID != "" {
			results = append(results, r.RunTest("cloudfront", "GetInvalidation", func() error {
				resp, err := client.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
					DistributionId: aws.String(invDistID),
					Id:             aws.String(invID),
				})
				if err != nil {
					return err
				}
				if resp.Invalidation == nil {
					return fmt.Errorf("invalidation is nil")
				}
				if resp.Invalidation.Status == nil {
					return fmt.Errorf("invalidation status is nil")
				}
				return nil
			}))
		}

		results = append(results, r.RunTest("cloudfront", "ListInvalidations", func() error {
			resp, err := client.ListInvalidations(ctx, &cloudfront.ListInvalidationsInput{
				DistributionId: aws.String(invDistID),
			})
			if err != nil {
				return err
			}
			if resp.InvalidationList == nil {
				return fmt.Errorf("invalidation list is nil")
			}
			if resp.InvalidationList.Quantity == nil || *resp.InvalidationList.Quantity < 1 {
				return fmt.Errorf("expected at least 1 invalidation, got 0")
			}
			return nil
		}))

		// Clean up
		results = append(results, r.RunTest("cloudfront", "Cleanup_InvalidationDist", func() error {
			getResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
				Id: aws.String(invDistID),
			})
			if err != nil {
				return err
			}
			getResp.DistributionConfig.Enabled = aws.Bool(false)
			updateResp, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
				Id:                 aws.String(invDistID),
				IfMatch:            aws.String(aws.ToString(getResp.ETag)),
				DistributionConfig: getResp.DistributionConfig,
			})
			if err != nil {
				return err
			}
			_, err = client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
				Id:      aws.String(invDistID),
				IfMatch: updateResp.ETag,
			})
			return err
		}))
	}

	// --- Origin Access Control tests ---
	results = append(results, r.RunTest("cloudfront", "ListOriginAccessControls", func() error {
		resp, err := client.ListOriginAccessControls(ctx, &cloudfront.ListOriginAccessControlsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.OriginAccessControlList == nil {
			return fmt.Errorf("OAC list is nil")
		}
		return nil
	}))

	oacName := fmt.Sprintf("test-oac-%d", time.Now().UnixNano())
	var oacID string
	var oacETag string
	results = append(results, r.RunTest("cloudfront", "CreateOriginAccessControl", func() error {
		resp, err := client.CreateOriginAccessControl(ctx, &cloudfront.CreateOriginAccessControlInput{
			OriginAccessControlConfig: &types.OriginAccessControlConfig{
				Name:                          aws.String(oacName),
				OriginAccessControlOriginType: types.OriginAccessControlOriginTypesS3,
				SigningBehavior:               types.OriginAccessControlSigningBehaviorsNever,
				SigningProtocol:               types.OriginAccessControlSigningProtocolsSigv4,
			},
		})
		if err == nil {
			oacID = aws.ToString(resp.OriginAccessControl.Id)
			oacETag = aws.ToString(resp.ETag)
		}
		return err
	}))

	if oacID != "" {
		results = append(results, r.RunTest("cloudfront", "GetOriginAccessControl", func() error {
			resp, err := client.GetOriginAccessControl(ctx, &cloudfront.GetOriginAccessControlInput{
				Id: aws.String(oacID),
			})
			if err != nil {
				return err
			}
			if resp.OriginAccessControl == nil {
				return fmt.Errorf("OAC is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "GetOriginAccessControlConfig", func() error {
			resp, err := client.GetOriginAccessControlConfig(ctx, &cloudfront.GetOriginAccessControlConfigInput{
				Id: aws.String(oacID),
			})
			if err != nil {
				return err
			}
			if resp.OriginAccessControlConfig == nil {
				return fmt.Errorf("OAC config is nil")
			}
			if resp.OriginAccessControlConfig.Name == nil || *resp.OriginAccessControlConfig.Name != oacName {
				return fmt.Errorf("OAC config name mismatch")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "UpdateOriginAccessControl", func() error {
			updatedName := oacName + "-updated"
			resp, err := client.UpdateOriginAccessControl(ctx, &cloudfront.UpdateOriginAccessControlInput{
				Id:      aws.String(oacID),
				IfMatch: aws.String(oacETag),
				OriginAccessControlConfig: &types.OriginAccessControlConfig{
					Name:                          aws.String(updatedName),
					OriginAccessControlOriginType: types.OriginAccessControlOriginTypesS3,
					SigningBehavior:               types.OriginAccessControlSigningBehaviorsAlways,
					SigningProtocol:               types.OriginAccessControlSigningProtocolsSigv4,
					Description:                   aws.String("updated description"),
				},
			})
			if err != nil {
				return err
			}
			if resp.OriginAccessControl == nil {
				return fmt.Errorf("updated OAC is nil")
			}
			if resp.OriginAccessControl.OriginAccessControlConfig.Name == nil || *resp.OriginAccessControl.OriginAccessControlConfig.Name != updatedName {
				return fmt.Errorf("OAC name not updated")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "DeleteOriginAccessControl", func() error {
			resp, err := client.DeleteOriginAccessControl(ctx, &cloudfront.DeleteOriginAccessControlInput{
				Id: aws.String(oacID),
			})
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			return nil
		}))
	}

	// --- Cache Policy tests ---
	results = append(results, r.RunTest("cloudfront", "ListCachePolicies", func() error {
		resp, err := client.ListCachePolicies(ctx, &cloudfront.ListCachePoliciesInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.CachePolicyList == nil {
			return fmt.Errorf("cache policy list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudfront", "GetCachePolicy_Managed", func() error {
		resp, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
			Id: aws.String("658327ea-f89d-4fab-a63d-7e88639e58f6"),
		})
		if err != nil {
			return err
		}
		if resp.CachePolicy == nil {
			return fmt.Errorf("cache policy is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudfront", "GetCachePolicyConfig_Managed", func() error {
		resp, err := client.GetCachePolicyConfig(ctx, &cloudfront.GetCachePolicyConfigInput{
			Id: aws.String("658327ea-f89d-4fab-a63d-7e88639e58f6"),
		})
		if err != nil {
			return err
		}
		if resp.CachePolicyConfig == nil {
			return fmt.Errorf("cache policy config is nil")
		}
		return nil
	}))

	var cachePolicyID string
	var cachePolicyETag string
	results = append(results, r.RunTest("cloudfront", "CreateCachePolicy", func() error {
		resp, err := client.CreateCachePolicy(ctx, &cloudfront.CreateCachePolicyInput{
			CachePolicyConfig: &types.CachePolicyConfig{
				Name:       aws.String(fmt.Sprintf("test-cp-%d", time.Now().UnixNano())),
				Comment:    aws.String("Test cache policy"),
				DefaultTTL: aws.Int64(3600),
				MaxTTL:     aws.Int64(86400),
				MinTTL:     aws.Int64(0),
			},
		})
		if err == nil {
			cachePolicyID = aws.ToString(resp.CachePolicy.Id)
			cachePolicyETag = aws.ToString(resp.ETag)
		}
		return err
	}))

	if cachePolicyID != "" {
		results = append(results, r.RunTest("cloudfront", "GetCachePolicy_Custom", func() error {
			resp, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
				Id: aws.String(cachePolicyID),
			})
			if err != nil {
				return err
			}
			if resp.CachePolicy == nil {
				return fmt.Errorf("cache policy is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "GetCachePolicyConfig_Custom", func() error {
			resp, err := client.GetCachePolicyConfig(ctx, &cloudfront.GetCachePolicyConfigInput{
				Id: aws.String(cachePolicyID),
			})
			if err != nil {
				return err
			}
			if resp.CachePolicyConfig == nil {
				return fmt.Errorf("cache policy config is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "UpdateCachePolicy", func() error {
			resp, err := client.UpdateCachePolicy(ctx, &cloudfront.UpdateCachePolicyInput{
				Id:      aws.String(cachePolicyID),
				IfMatch: aws.String(cachePolicyETag),
				CachePolicyConfig: &types.CachePolicyConfig{
					Name:       aws.String("updated-cache-policy"),
					Comment:    aws.String("Updated test cache policy"),
					DefaultTTL: aws.Int64(1800),
					MaxTTL:     aws.Int64(7200),
					MinTTL:     aws.Int64(60),
				},
			})
			if err != nil {
				return err
			}
			if resp.CachePolicy == nil {
				return fmt.Errorf("updated cache policy is nil")
			}
			if resp.CachePolicy.CachePolicyConfig == nil || *resp.CachePolicy.CachePolicyConfig.DefaultTTL != 1800 {
				return fmt.Errorf("cache policy DefaultTTL not updated")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "DeleteCachePolicy", func() error {
			_, err := client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{
				Id: aws.String(cachePolicyID),
			})
			return err
		}))

		results = append(results, r.RunTest("cloudfront", "GetCachePolicy_AfterDelete", func() error {
			_, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
				Id: aws.String(cachePolicyID),
			})
			if err == nil {
				return fmt.Errorf("expected error for deleted cache policy")
			}
			return nil
		}))
	}

	// --- Origin Request Policy tests ---
	results = append(results, r.RunTest("cloudfront", "ListOriginRequestPolicies", func() error {
		resp, err := client.ListOriginRequestPolicies(ctx, &cloudfront.ListOriginRequestPoliciesInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.OriginRequestPolicyList == nil {
			return fmt.Errorf("origin request policy list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudfront", "GetOriginRequestPolicy_Managed", func() error {
		resp, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
			Id: aws.String("88a5eaf4-2fd4-4709-b370-b4c650ea3fcf"),
		})
		if err != nil {
			return err
		}
		if resp.OriginRequestPolicy == nil {
			return fmt.Errorf("origin request policy is nil")
		}
		return nil
	}))

	var orpID string
	var orpETag string
	results = append(results, r.RunTest("cloudfront", "CreateOriginRequestPolicy", func() error {
		resp, err := client.CreateOriginRequestPolicy(ctx, &cloudfront.CreateOriginRequestPolicyInput{
			OriginRequestPolicyConfig: &types.OriginRequestPolicyConfig{
				Name:    aws.String(fmt.Sprintf("test-orp-%d", time.Now().UnixNano())),
				Comment: aws.String("Test ORP"),
				CookiesConfig: &types.OriginRequestPolicyCookiesConfig{
					CookieBehavior: types.OriginRequestPolicyCookieBehaviorNone,
				},
				HeadersConfig: &types.OriginRequestPolicyHeadersConfig{
					HeaderBehavior: types.OriginRequestPolicyHeaderBehaviorNone,
				},
				QueryStringsConfig: &types.OriginRequestPolicyQueryStringsConfig{
					QueryStringBehavior: types.OriginRequestPolicyQueryStringBehaviorNone,
				},
			},
		})
		if err == nil {
			orpID = aws.ToString(resp.OriginRequestPolicy.Id)
			orpETag = aws.ToString(resp.ETag)
		}
		return err
	}))

	if orpID != "" {
		results = append(results, r.RunTest("cloudfront", "GetOriginRequestPolicy_Custom", func() error {
			resp, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
				Id: aws.String(orpID),
			})
			if err != nil {
				return err
			}
			if resp.OriginRequestPolicy == nil {
				return fmt.Errorf("origin request policy is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "GetOriginRequestPolicyConfig", func() error {
			resp, err := client.GetOriginRequestPolicyConfig(ctx, &cloudfront.GetOriginRequestPolicyConfigInput{
				Id: aws.String(orpID),
			})
			if err != nil {
				return err
			}
			if resp.OriginRequestPolicyConfig == nil {
				return fmt.Errorf("origin request policy config is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "UpdateOriginRequestPolicy", func() error {
			resp, err := client.UpdateOriginRequestPolicy(ctx, &cloudfront.UpdateOriginRequestPolicyInput{
				Id:      aws.String(orpID),
				IfMatch: aws.String(orpETag),
				OriginRequestPolicyConfig: &types.OriginRequestPolicyConfig{
					Name:    aws.String("updated-orp"),
					Comment: aws.String("Updated test ORP"),
					CookiesConfig: &types.OriginRequestPolicyCookiesConfig{
						CookieBehavior: types.OriginRequestPolicyCookieBehaviorAll,
					},
					HeadersConfig: &types.OriginRequestPolicyHeadersConfig{
						HeaderBehavior: types.OriginRequestPolicyHeaderBehaviorAllViewer,
					},
					QueryStringsConfig: &types.OriginRequestPolicyQueryStringsConfig{
						QueryStringBehavior: types.OriginRequestPolicyQueryStringBehaviorAll,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp.OriginRequestPolicy == nil {
				return fmt.Errorf("updated ORP is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "DeleteOriginRequestPolicy", func() error {
			_, err := client.DeleteOriginRequestPolicy(ctx, &cloudfront.DeleteOriginRequestPolicyInput{
				Id: aws.String(orpID),
			})
			return err
		}))

		results = append(results, r.RunTest("cloudfront", "GetOriginRequestPolicy_AfterDelete", func() error {
			_, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
				Id: aws.String(orpID),
			})
			if err == nil {
				return fmt.Errorf("expected error for deleted ORP")
			}
			return nil
		}))
	}

	// --- Response Headers Policy tests ---
	results = append(results, r.RunTest("cloudfront", "ListResponseHeadersPolicies", func() error {
		resp, err := client.ListResponseHeadersPolicies(ctx, &cloudfront.ListResponseHeadersPoliciesInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.ResponseHeadersPolicyList == nil {
			return fmt.Errorf("response headers policy list is nil")
		}
		return nil
	}))

	rhpName := fmt.Sprintf("test-rhp-%d", time.Now().UnixNano())
	var rhpID string
	var rhpETag string
	results = append(results, r.RunTest("cloudfront", "CreateResponseHeadersPolicy", func() error {
		resp, err := client.CreateResponseHeadersPolicy(ctx, &cloudfront.CreateResponseHeadersPolicyInput{
			ResponseHeadersPolicyConfig: &types.ResponseHeadersPolicyConfig{
				Name:    aws.String(rhpName),
				Comment: aws.String("Test RHP"),
				ServerTimingHeadersConfig: &types.ResponseHeadersPolicyServerTimingHeadersConfig{
					Enabled:      aws.Bool(true),
					SamplingRate: aws.Float64(0.5),
				},
				SecurityHeadersConfig: &types.ResponseHeadersPolicySecurityHeadersConfig{
					XSSProtection: &types.ResponseHeadersPolicyXSSProtection{
						Override:   aws.Bool(true),
						Protection: aws.Bool(true),
						ModeBlock:  aws.Bool(true),
					},
					ContentTypeOptions: &types.ResponseHeadersPolicyContentTypeOptions{
						Override: aws.Bool(true),
					},
				},
			},
		})
		if err == nil {
			rhpID = aws.ToString(resp.ResponseHeadersPolicy.Id)
			rhpETag = aws.ToString(resp.ETag)
		}
		return err
	}))

	if rhpID != "" {
		results = append(results, r.RunTest("cloudfront", "GetResponseHeadersPolicy", func() error {
			resp, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
				Id: aws.String(rhpID),
			})
			if err != nil {
				return err
			}
			if resp.ResponseHeadersPolicy == nil {
				return fmt.Errorf("response headers policy is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "GetResponseHeadersPolicyConfig", func() error {
			resp, err := client.GetResponseHeadersPolicyConfig(ctx, &cloudfront.GetResponseHeadersPolicyConfigInput{
				Id: aws.String(rhpID),
			})
			if err != nil {
				return err
			}
			if resp.ResponseHeadersPolicyConfig == nil {
				return fmt.Errorf("response headers policy config is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "UpdateResponseHeadersPolicy", func() error {
			resp, err := client.UpdateResponseHeadersPolicy(ctx, &cloudfront.UpdateResponseHeadersPolicyInput{
				Id:      aws.String(rhpID),
				IfMatch: aws.String(rhpETag),
				ResponseHeadersPolicyConfig: &types.ResponseHeadersPolicyConfig{
					Name:    aws.String(rhpName + "-updated"),
					Comment: aws.String("Updated RHP"),
					CustomHeadersConfig: &types.ResponseHeadersPolicyCustomHeadersConfig{
						Quantity: aws.Int32(1),
						Items: []types.ResponseHeadersPolicyCustomHeader{
							{
								Header:   aws.String("X-Custom-Header"),
								Value:    aws.String("custom-value"),
								Override: aws.Bool(true),
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
			if resp.ResponseHeadersPolicy == nil {
				return fmt.Errorf("updated RHP is nil")
			}
			return nil
		}))

		results = append(results, r.RunTest("cloudfront", "DeleteResponseHeadersPolicy", func() error {
			_, err := client.DeleteResponseHeadersPolicy(ctx, &cloudfront.DeleteResponseHeadersPolicyInput{
				Id: aws.String(rhpID),
			})
			return err
		}))

		results = append(results, r.RunTest("cloudfront", "GetResponseHeadersPolicy_AfterDelete", func() error {
			_, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
				Id: aws.String(rhpID),
			})
			if err == nil {
				return fmt.Errorf("expected error for deleted RHP")
			}
			return nil
		}))
	}

	// --- List Key Groups (unchanged) ---
	results = append(results, r.RunTest("cloudfront", "ListKeyGroups", func() error {
		resp, err := client.ListKeyGroups(ctx, &cloudfront.ListKeyGroupsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.KeyGroupList == nil {
			return fmt.Errorf("key group list is nil")
		}
		return nil
	}))

	return results
}
