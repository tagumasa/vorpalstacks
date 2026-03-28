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

	// Test 10: List Origin Access Controls
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

	// Test 11: Create Origin Access Control
	oacName := fmt.Sprintf("test-oac-%d", time.Now().UnixNano())
	var oacID string
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
		}
		return err
	}))

	// Test 12: Get Origin Access Control
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

		// Test 13: Delete Origin Access Control
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

	// Test 14: List Key Groups
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

	// Test 15: List Cache Policies
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

	// Test 16: Get Cache Policy
	results = append(results, r.RunTest("cloudfront", "GetCachePolicy", func() error {
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

	// Test 17: List Origin Request Policies
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

	// Test 18: List Response Headers Policies
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

	return results
}
