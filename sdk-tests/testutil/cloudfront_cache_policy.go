package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfCachePolicyTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListCachePolicies_VerifyFields", func() error {
		resp, err := client.ListCachePolicies(ctx, &cloudfront.ListCachePoliciesInput{
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.CachePolicyList == nil {
			return fmt.Errorf("cache policy list is nil")
		}
		if resp.CachePolicyList.MaxItems == nil {
			return fmt.Errorf("maxItems is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListCachePolicies_ManagedFilter", func() error {
		resp, err := client.ListCachePolicies(ctx, &cloudfront.ListCachePoliciesInput{
			MaxItems: aws.Int32(100),
			Type:     types.CachePolicyTypeManaged,
		})
		if err != nil {
			return err
		}
		if resp.CachePolicyList == nil {
			return fmt.Errorf("cache policy list is nil")
		}
		for _, item := range resp.CachePolicyList.Items {
			if item.Type != types.CachePolicyTypeManaged {
				return fmt.Errorf("managed filter returned a %q item", item.Type)
			}
		}
		if len(resp.CachePolicyList.Items) == 0 {
			return fmt.Errorf("managed filter returned no items")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListCachePolicies_CustomFilter", func() error {
		resp, err := client.ListCachePolicies(ctx, &cloudfront.ListCachePoliciesInput{
			MaxItems: aws.Int32(100),
			Type:     types.CachePolicyTypeCustom,
		})
		if err != nil {
			return err
		}
		for _, item := range resp.CachePolicyList.Items {
			if item.Type != types.CachePolicyTypeCustom {
				return fmt.Errorf("custom filter returned a %q item", item.Type)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetCachePolicy_Managed", func() error {
		resp, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
			Id: aws.String("658327ea-f89d-4fab-a63d-7e88639e58f6"),
		})
		if err != nil {
			return err
		}
		if resp.CachePolicy == nil {
			return fmt.Errorf("cache policy is nil")
		}
		if resp.CachePolicy.Id == nil || *resp.CachePolicy.Id != "658327ea-f89d-4fab-a63d-7e88639e58f6" {
			return fmt.Errorf("cache policy ID mismatch")
		}
		if resp.CachePolicy.CachePolicyConfig == nil {
			return fmt.Errorf("cache policy config is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetCachePolicyConfig_Managed", func() error {
		resp, err := client.GetCachePolicyConfig(ctx, &cloudfront.GetCachePolicyConfigInput{
			Id: aws.String("658327ea-f89d-4fab-a63d-7e88639e58f6"),
		})
		if err != nil {
			return err
		}
		if resp.CachePolicyConfig == nil {
			return fmt.Errorf("cache policy config is nil")
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		return nil
	}))

	cpName := tc.uniquePrefix("test-cp")
	var cpID string
	var cpETag string
	results = append(results, tc.runner.RunTest("cloudfront", "CreateCachePolicy_VerifyFields", func() error {
		resp, err := client.CreateCachePolicy(ctx, &cloudfront.CreateCachePolicyInput{
			CachePolicyConfig: &types.CachePolicyConfig{
				Name:       aws.String(cpName),
				Comment:    aws.String("Test cache policy"),
				DefaultTTL: aws.Int64(3600),
				MaxTTL:     aws.Int64(86400),
				MinTTL:     aws.Int64(0),
			},
		})
		if err != nil {
			return err
		}
		if resp.CachePolicy == nil {
			return fmt.Errorf("cache policy is nil")
		}
		if resp.CachePolicy.Id == nil || *resp.CachePolicy.Id == "" {
			return fmt.Errorf("cache policy ID is nil or empty")
		}
		if resp.CachePolicy.CachePolicyConfig == nil {
			return fmt.Errorf("cache policy config is nil")
		}
		cfg := resp.CachePolicy.CachePolicyConfig
		if cfg.Name == nil || *cfg.Name != cpName {
			return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(cfg.Name), cpName)
		}
		if cfg.DefaultTTL == nil || *cfg.DefaultTTL != 3600 {
			return fmt.Errorf("defaultTTL mismatch: got %d", aws.ToInt64(cfg.DefaultTTL))
		}
		if cfg.MaxTTL == nil || *cfg.MaxTTL != 86400 {
			return fmt.Errorf("maxTTL mismatch: got %d", aws.ToInt64(cfg.MaxTTL))
		}
		if cfg.MinTTL == nil || *cfg.MinTTL != 0 {
			return fmt.Errorf("minTTL mismatch: got %d", aws.ToInt64(cfg.MinTTL))
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		cpID = *resp.CachePolicy.Id
		cpETag = *resp.ETag
		return nil
	}))

	if cpID != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "GetCachePolicy_Custom_VerifyFields", func() error {
			resp, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
				Id: aws.String(cpID),
			})
			if err != nil {
				return err
			}
			if resp.CachePolicy == nil {
				return fmt.Errorf("cache policy is nil")
			}
			if resp.CachePolicy.Id == nil || *resp.CachePolicy.Id != cpID {
				return fmt.Errorf("cache policy ID mismatch: got %q", aws.ToString(resp.CachePolicy.Id))
			}
			if resp.CachePolicy.LastModifiedTime == nil {
				return fmt.Errorf("lastModifiedTime is nil")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetCachePolicyConfig_Custom_VerifyFields", func() error {
			resp, err := client.GetCachePolicyConfig(ctx, &cloudfront.GetCachePolicyConfigInput{
				Id: aws.String(cpID),
			})
			if err != nil {
				return err
			}
			if resp.CachePolicyConfig == nil {
				return fmt.Errorf("cache policy config is nil")
			}
			if resp.CachePolicyConfig.Name == nil || *resp.CachePolicyConfig.Name != cpName {
				return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(resp.CachePolicyConfig.Name), cpName)
			}
			if resp.ETag == nil || *resp.ETag == "" {
				return fmt.Errorf("ETag is nil or empty")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "UpdateCachePolicy_VerifyFields", func() error {
			resp, err := client.UpdateCachePolicy(ctx, &cloudfront.UpdateCachePolicyInput{
				Id:      aws.String(cpID),
				IfMatch: aws.String(cpETag),
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
			cfg := resp.CachePolicy.CachePolicyConfig
			if cfg == nil {
				return fmt.Errorf("cache policy config is nil")
			}
			if cfg.Name == nil || *cfg.Name != "updated-cache-policy" {
				return fmt.Errorf("name not updated: got %q", aws.ToString(cfg.Name))
			}
			if cfg.DefaultTTL == nil || *cfg.DefaultTTL != 1800 {
				return fmt.Errorf("defaultTTL not updated: got %d", aws.ToInt64(cfg.DefaultTTL))
			}
			if cfg.MaxTTL == nil || *cfg.MaxTTL != 7200 {
				return fmt.Errorf("maxTTL not updated: got %d", aws.ToInt64(cfg.MaxTTL))
			}
			if cfg.MinTTL == nil || *cfg.MinTTL != 60 {
				return fmt.Errorf("minTTL not updated: got %d", aws.ToInt64(cfg.MinTTL))
			}
			if resp.ETag != nil {
				cpETag = *resp.ETag
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "DeleteCachePolicy", func() error {
			_, err := client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{
				Id:      aws.String(cpID),
				IfMatch: aws.String(cpETag),
			})
			return err
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetCachePolicy_NonExistent", func() error {
			_, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
				Id: aws.String(cpID),
			})
			return AssertErrorContains(err, "NoSuchCachePolicy")
		}))
	}

	results = append(results, tc.runner.RunTest("cloudfront", "ListCachePolicies_Pagination", func() error {
		pgTs := time.Now().UnixNano()
		var pgIDs []string
		var pgETags []string
		for i := 0; i < 5; i++ {
			resp, err := client.CreateCachePolicy(ctx, &cloudfront.CreateCachePolicyInput{
				CachePolicyConfig: &types.CachePolicyConfig{
					Name:       aws.String(fmt.Sprintf("pagcp-%d-%d", pgTs, i)),
					Comment:    aws.String("pagination test"),
					DefaultTTL: aws.Int64(3600),
					MaxTTL:     aws.Int64(86400),
					MinTTL:     aws.Int64(0),
				},
			})
			if err != nil {
				return fmt.Errorf("create cache policy: %v", err)
			}
			pgIDs = append(pgIDs, aws.ToString(resp.CachePolicy.Id))
			pgETags = append(pgETags, aws.ToString(resp.ETag))
		}

		pageCount := 0
		totalCount := 0
		var marker *string
		for {
			resp, err := client.ListCachePolicies(ctx, &cloudfront.ListCachePoliciesInput{
				Marker:   marker,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				for i, id := range pgIDs {
					client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{Id: aws.String(id), IfMatch: aws.String(pgETags[i])})
				}
				return fmt.Errorf("list cache policies page: %v", err)
			}
			pageCount++
			if resp.CachePolicyList != nil {
				totalCount += len(resp.CachePolicyList.Items)
			}
			if resp.CachePolicyList != nil && resp.CachePolicyList.NextMarker != nil && *resp.CachePolicyList.NextMarker != "" {
				marker = resp.CachePolicyList.NextMarker
			} else {
				break
			}
		}

		for i, id := range pgIDs {
			client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{
				Id:      aws.String(id),
				IfMatch: aws.String(pgETags[i]),
			})
		}

		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d (total items: %d)", pageCount, totalCount)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributionsByX_RoundTrip", func() error {
		byRefTs := time.Now().UnixNano()
		cpResp, err := client.CreateCachePolicy(ctx, &cloudfront.CreateCachePolicyInput{
			CachePolicyConfig: &types.CachePolicyConfig{
				Name:       aws.String(fmt.Sprintf("byref-cp-%d", byRefTs)),
				Comment:    aws.String("list-by-reference test"),
				DefaultTTL: aws.Int64(3600),
				MaxTTL:     aws.Int64(86400),
				MinTTL:     aws.Int64(0),
			},
		})
		if err != nil {
			return err
		}
		byRefCpID := aws.ToString(cpResp.CachePolicy.Id)
		byRefCpETag := aws.ToString(cpResp.ETag)

		distID, distETag, err := tc.createDistribution(tc.uniqueCallerRef("byref-dist"), "list-by-reference distribution", "example.net")
		if err != nil {
			client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{Id: aws.String(byRefCpID), IfMatch: aws.String(byRefCpETag)})
			return err
		}

		// Attach the cache policy to the distribution's default behaviour.
		getResp, err := client.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err == nil {
			getResp.DistributionConfig.DefaultCacheBehavior.CachePolicyId = aws.String(byRefCpID)
			updResp, err := client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
				Id:                 aws.String(distID),
				IfMatch:            aws.String(distETag),
				DistributionConfig: getResp.DistributionConfig,
			})
			if err == nil && updResp.ETag != nil {
				distETag = *updResp.ETag
			}
		}
		if err != nil {
			tc.disableAndDeleteDistribution(distID, distETag)
			client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{Id: aws.String(byRefCpID), IfMatch: aws.String(byRefCpETag)})
			return err
		}

		defer func() {
			tc.disableAndDeleteDistribution(distID, distETag)
			client.DeleteCachePolicy(ctx, &cloudfront.DeleteCachePolicyInput{Id: aws.String(byRefCpID), IfMatch: aws.String(byRefCpETag)})
		}()

		byCp, err := client.ListDistributionsByCachePolicyId(ctx, &cloudfront.ListDistributionsByCachePolicyIdInput{
			CachePolicyId: aws.String(byRefCpID),
		})
		if err != nil {
			return err
		}
		if byCp.DistributionIdList == nil {
			return fmt.Errorf("distribution id list is nil")
		}
		found := false
		for _, id := range byCp.DistributionIdList.Items {
			if id == distID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("distribution %q not listed by cache policy %q", distID, byRefCpID)
		}

		// Unknown target IDs must produce the modelled 404 errors rather
		// than an empty successful listing.
		randomID := "22222222-2222-4222-8222-222222222222"

		_, err = client.ListDistributionsByCachePolicyId(ctx, &cloudfront.ListDistributionsByCachePolicyIdInput{
			CachePolicyId: aws.String(randomID),
		})
		if err := AssertErrorContains(err, "NoSuchCachePolicy"); err != nil {
			return err
		}

		_, err = client.ListDistributionsByResponseHeadersPolicyId(ctx, &cloudfront.ListDistributionsByResponseHeadersPolicyIdInput{
			ResponseHeadersPolicyId: aws.String(randomID),
		})
		if err := AssertErrorContains(err, "NoSuchResponseHeadersPolicy"); err != nil {
			return err
		}

		_, err = client.ListDistributionsByOriginRequestPolicyId(ctx, &cloudfront.ListDistributionsByOriginRequestPolicyIdInput{
			OriginRequestPolicyId: aws.String(randomID),
		})
		if err := AssertErrorContains(err, "NoSuchOriginRequestPolicy"); err != nil {
			return err
		}

		_, err = client.ListDistributionsByKeyGroup(ctx, &cloudfront.ListDistributionsByKeyGroupInput{
			KeyGroupId: aws.String(randomID),
		})
		if err := AssertErrorContains(err, "NoSuchResource"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
