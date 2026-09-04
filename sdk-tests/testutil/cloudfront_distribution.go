package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfDistributionTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributions_VerifyFields", func() error {
		resp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.DistributionList == nil {
			return fmt.Errorf("distribution list is nil")
		}
		if resp.DistributionList.Quantity == nil {
			return fmt.Errorf("quantity is nil")
		}
		return nil
	}))

	var distID string
	var distETag string
	callerRef := tc.uniquePrefix("test-cf")
	originID := tc.uniquePrefix("test-origin")
	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_VerifyFields", func() error {
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
		if err != nil {
			return err
		}
		if resp.Distribution == nil {
			return fmt.Errorf("distribution is nil")
		}
		if resp.Distribution.Id == nil || *resp.Distribution.Id == "" {
			return fmt.Errorf("distribution ID is nil or empty")
		}
		if resp.Distribution.ARN == nil || *resp.Distribution.ARN == "" {
			return fmt.Errorf("distribution ARN is nil or empty")
		}
		if resp.Distribution.Status == nil || *resp.Distribution.Status == "" {
			return fmt.Errorf("distribution Status is nil or empty")
		}
		if resp.Distribution.DistributionConfig == nil {
			return fmt.Errorf("distribution config is nil")
		}
		if resp.Distribution.DistributionConfig.Comment == nil || *resp.Distribution.DistributionConfig.Comment != "SDK test distribution" {
			return fmt.Errorf("comment mismatch: got %q", aws.ToString(resp.Distribution.DistributionConfig.Comment))
		}
		if resp.Distribution.DistributionConfig.Enabled == nil || !*resp.Distribution.DistributionConfig.Enabled {
			return fmt.Errorf("enabled should be true")
		}
		if resp.Distribution.DistributionConfig.DefaultRootObject == nil || *resp.Distribution.DistributionConfig.DefaultRootObject != "index.html" {
			return fmt.Errorf("defaultRootObject mismatch: got %q", aws.ToString(resp.Distribution.DistributionConfig.DefaultRootObject))
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		distID = *resp.Distribution.Id
		distETag = *resp.ETag
		return nil
	}))

	if distID == "" {
		return results
	}

	// Negative validation tests: the server must reject structurally
	// invalid distribution configurations on both the create and update
	// paths.
	baseCfg := func() *types.DistributionConfig {
		return tc.baseDistributionConfig(tc.uniquePrefix("neg-cf"), "negative validation", "example.com")
	}

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_MissingOriginDomain_Rejected", func() error {
		cfg := baseCfg()
		cfg.Origins.Items[0].DomainName = aws.String("")
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_DuplicateOriginId_Rejected", func() error {
		cfg := baseCfg()
		first := cfg.Origins.Items[0]
		cfg.Origins.Quantity = aws.Int32(2)
		cfg.Origins.Items = append(cfg.Origins.Items, types.Origin{
			Id:         first.Id,
			DomainName: aws.String("other.example.com"),
			CustomOriginConfig: &types.CustomOriginConfig{
				HTTPPort:             aws.Int32(80),
				HTTPSPort:            aws.Int32(443),
				OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
			},
		})
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_AliasTooLong_Rejected", func() error {
		cfg := baseCfg()
		longAlias := fmt.Sprintf("%s.example.com", tc.uniquePrefix("a"))
		for len(longAlias) < 260 {
			longAlias = "a" + longAlias
		}
		cfg.Aliases = &types.Aliases{
			Quantity: aws.Int32(1),
			Items:    []string{longAlias},
		}
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_BogusWebACLId_Rejected", func() error {
		cfg := baseCfg()
		cfg.WebACLId = aws.String("arn:aws:wafv2:us-east-1:000000000000:global/webacl/nonexistent/11111111-1111-1111-1111-111111111111")
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		return AssertErrorContains(err, "InvalidWebACLId")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_UnknownCachePolicy_Rejected", func() error {
		cfg := baseCfg()
		cfg.DefaultCacheBehavior.CachePolicyId = aws.String("11111111-2222-3333-4444-555555555555")
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		return AssertErrorContains(err, "NoSuchCachePolicy")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_CNAMEConflict_Rejected", func() error {
		alias := fmt.Sprintf("conflict-%d.example.com", time.Now().UnixNano())
		cfgA := baseCfg()
		cfgA.Aliases = &types.Aliases{Quantity: aws.Int32(1), Items: []string{alias}}
		respA, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfgA})
		if err != nil {
			return err
		}
		idA := aws.ToString(respA.Distribution.Id)
		defer tc.disableAndDeleteDistribution(idA, aws.ToString(respA.ETag))

		cfgB := baseCfg()
		cfgB.Aliases = &types.Aliases{Quantity: aws.Int32(1), Items: []string{alias}}
		_, err = client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfgB})
		if err := AssertErrorContains(err, "CNAMEAlreadyExists"); err != nil {
			return err
		}

		// Re-submitting the distribution's own configuration must not
		// count as a conflict with itself.
		getResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(idA)})
		if err != nil {
			return err
		}
		_, err = client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
			Id:                 aws.String(idA),
			IfMatch:            aws.String(aws.ToString(getResp.ETag)),
			DistributionConfig: getResp.DistributionConfig,
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetDistribution_VerifyFields", func() error {
		resp, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
			Id: aws.String(distID),
		})
		if err != nil {
			return err
		}
		if resp.Distribution == nil {
			return fmt.Errorf("distribution is nil")
		}
		if resp.Distribution.Id == nil || *resp.Distribution.Id != distID {
			return fmt.Errorf("distribution ID mismatch: got %q, want %q", aws.ToString(resp.Distribution.Id), distID)
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		if resp.Distribution.DistributionConfig == nil {
			return fmt.Errorf("distribution config is nil")
		}
		if resp.Distribution.LastModifiedTime == nil {
			return fmt.Errorf("lastModifiedTime is nil")
		}
		if resp.Distribution.InProgressInvalidationBatches != nil && *resp.Distribution.InProgressInvalidationBatches != 0 {
			return fmt.Errorf("unexpected in-progress invalidation batches: %d", *resp.Distribution.InProgressInvalidationBatches)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetDistributionConfig_VerifyFields", func() error {
		resp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
			Id: aws.String(distID),
		})
		if err != nil {
			return err
		}
		if resp.DistributionConfig == nil {
			return fmt.Errorf("distribution config is nil")
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		cfg := resp.DistributionConfig
		if cfg.CallerReference == nil || *cfg.CallerReference != callerRef {
			return fmt.Errorf("callerReference mismatch: got %q, want %q", aws.ToString(cfg.CallerReference), callerRef)
		}
		if cfg.Comment == nil || *cfg.Comment != "SDK test distribution" {
			return fmt.Errorf("comment mismatch: got %q", aws.ToString(cfg.Comment))
		}
		if cfg.DefaultRootObject == nil || *cfg.DefaultRootObject != "index.html" {
			return fmt.Errorf("defaultRootObject mismatch: got %q", aws.ToString(cfg.DefaultRootObject))
		}
		if cfg.Enabled == nil || !*cfg.Enabled {
			return fmt.Errorf("enabled should be true")
		}
		if cfg.Origins == nil || cfg.Origins.Quantity == nil || *cfg.Origins.Quantity != 1 {
			return fmt.Errorf("origins quantity mismatch")
		}
		if cfg.DefaultCacheBehavior == nil {
			return fmt.Errorf("defaultCacheBehavior is nil")
		}
		if cfg.DefaultCacheBehavior.TargetOriginId == nil || *cfg.DefaultCacheBehavior.TargetOriginId != originID {
			return fmt.Errorf("targetOriginId mismatch: got %q, want %q", aws.ToString(cfg.DefaultCacheBehavior.TargetOriginId), originID)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributions_ContainsCreated", func() error {
		items, err := paginate(func(next *string) ([]types.DistributionSummary, *string, error) {
			resp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
				MaxItems: aws.Int32(100),
				Marker:   next,
			})
			if err != nil {
				return nil, nil, err
			}
			if resp.DistributionList == nil {
				return nil, nil, nil
			}
			return resp.DistributionList.Items, resp.DistributionList.NextMarker, nil
		})
		if err != nil {
			return err
		}
		if len(items) < 1 {
			return fmt.Errorf("expected at least 1 distribution, got 0")
		}
		for _, d := range items {
			if d.Id != nil && *d.Id == distID {
				if d.ARN == nil || *d.ARN == "" {
					return fmt.Errorf("distribution ARN is nil or empty in list")
				}
				return nil
			}
		}
		return fmt.Errorf("created distribution %q not found in list", distID)
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributionsByWebACLId_UnknownACL_Rejected", func() error {
		_, err := client.ListDistributionsByWebACLId(ctx, &cloudfront.ListDistributionsByWebACLIdInput{
			WebACLId: aws.String("12345678-1234-1234-1234-123456789012"),
		})
		return AssertErrorContains(err, "InvalidWebACLId")
	}))

	var updateETag string
	results = append(results, tc.runner.RunTest("cloudfront", "UpdateDistribution_Disable", func() error {
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
		if err != nil {
			return err
		}
		if resp.Distribution == nil {
			return fmt.Errorf("distribution is nil")
		}
		if resp.Distribution.DistributionConfig == nil {
			return fmt.Errorf("distribution config is nil")
		}
		if resp.Distribution.DistributionConfig.Enabled != nil && *resp.Distribution.DistributionConfig.Enabled {
			return fmt.Errorf("distribution should be disabled after update")
		}
		updateETag = aws.ToString(resp.ETag)
		return nil
	}))

	if updateETag != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "UpdateDistribution_DuplicateOriginId_Rejected", func() error {
			getResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
				Id: aws.String(distID),
			})
			if err != nil {
				return err
			}
			cfg := getResp.DistributionConfig
			firstID := cfg.Origins.Items[0].Id
			cfg.Origins.Quantity = aws.Int32(2)
			cfg.Origins.Items = append(cfg.Origins.Items, types.Origin{
				Id:         firstID,
				DomainName: aws.String("duplicate.example.com"),
				CustomOriginConfig: &types.CustomOriginConfig{
					HTTPPort:             aws.Int32(80),
					HTTPSPort:            aws.Int32(443),
					OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
				},
			})
			_, err = client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
				Id:                 aws.String(distID),
				IfMatch:            aws.String(updateETag),
				DistributionConfig: cfg,
			})
			return AssertErrorContains(err, "InvalidArgument")
		}))
	}

	if updateETag != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "DeleteDistribution", func() error {
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

	results = append(results, tc.runner.RunTest("cloudfront", "GetDistribution_NonExistent", func() error {
		_, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
			Id: aws.String(distID),
		})
		if aerr := AssertErrorContains(err, "NoSuchDistribution"); aerr != nil {
			return aerr
		}
		// An ID that never existed must fail identically to the deleted
		// one.
		_, err = client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
			Id: aws.String("nonexistent-dist-id"),
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "DeleteDistribution_NonExistent", func() error {
		_, err := client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
			Id:      aws.String(distID),
			IfMatch: aws.String("invalid-etag"),
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	return results
}

func (tc *cfTestContext) createDistribution(callerRef, comment, originDomain string) (distID, distETag string, err error) {
	resp, rerr := tc.client.CreateDistribution(tc.ctx, &cloudfront.CreateDistributionInput{
		DistributionConfig: tc.baseDistributionConfig(callerRef, comment, originDomain),
	})
	if rerr != nil {
		return "", "", rerr
	}
	return aws.ToString(resp.Distribution.Id), aws.ToString(resp.ETag), nil
}

func (tc *cfTestContext) disableAndDeleteDistribution(distID, etag string) error {
	getResp, err := tc.client.GetDistributionConfig(tc.ctx, &cloudfront.GetDistributionConfigInput{
		Id: aws.String(distID),
	})
	if err != nil {
		return err
	}
	getResp.DistributionConfig.Enabled = aws.Bool(false)
	updateResp, err := tc.client.UpdateDistribution(tc.ctx, &cloudfront.UpdateDistributionInput{
		Id:                 aws.String(distID),
		IfMatch:            aws.String(etag),
		DistributionConfig: getResp.DistributionConfig,
	})
	if err != nil {
		return err
	}
	_, err = tc.client.DeleteDistribution(tc.ctx, &cloudfront.DeleteDistributionInput{
		Id:      aws.String(distID),
		IfMatch: updateResp.ETag,
	})
	return err
}
