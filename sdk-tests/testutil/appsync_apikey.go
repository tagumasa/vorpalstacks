package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncApiKeyCacheTests(res *appsyncResources) []TestResult {
	var results []TestResult

	results = append(results, r.runAppSyncApiKeyTests(res)...)
	results = append(results, r.runAppSyncCacheTests(res)...)

	return results
}

func (r *TestRunner) runAppSyncApiKeyTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "CreateApiKey", func() error {
		resp, err := client.CreateApiKey(ctx, &appsync.CreateApiKeyInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.ApiKey == nil {
			return fmt.Errorf("ApiKey is nil")
		}
		if resp.ApiKey.Id == nil || *resp.ApiKey.Id == "" {
			return fmt.Errorf("id is empty")
		}
		if resp.ApiKey.Expires == 0 {
			return fmt.Errorf("expires should be set (default 365 days)")
		}
		if resp.ApiKey.Deletes == 0 {
			return fmt.Errorf("deletes should be set (same as Expires)")
		}
		res.apiKeyId = *resp.ApiKey.Id
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateApiKey_WithDescription", func() error {
		resp, err := client.CreateApiKey(ctx, &appsync.CreateApiKeyInput{
			ApiId:       aws.String(res.gqlApiId),
			Description: aws.String("test key"),
		})
		if err != nil {
			return err
		}
		if resp.ApiKey == nil {
			return fmt.Errorf("ApiKey is nil")
		}
		if resp.ApiKey.Description == nil || *resp.ApiKey.Description != "test key" {
			return fmt.Errorf("description not set: %v", resp.ApiKey.Description)
		}
		res.descApiKeyId = *resp.ApiKey.Id
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApiKeys", func() error {
		resp, err := client.ListApiKeys(ctx, &appsync.ListApiKeysInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.ApiKeys) < 2 {
			return fmt.Errorf("expected at least 2 API keys, got %d", len(resp.ApiKeys))
		}
		found := false
		for _, key := range resp.ApiKeys {
			if key.Id != nil && *key.Id == res.apiKeyId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created API key not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApiKeys_WithPagination", func() error {
		resp, err := client.ListApiKeys(ctx, &appsync.ListApiKeysInput{
			ApiId:      aws.String(res.gqlApiId),
			MaxResults: 1,
		})
		if err != nil {
			return err
		}
		if len(resp.ApiKeys) != 1 {
			return fmt.Errorf("expected 1 API key with maxResults=1, got %d", len(resp.ApiKeys))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApiKey", func() error {
		newExpiry := time.Now().Add(180 * 24 * time.Hour).Unix()
		resp, err := client.UpdateApiKey(ctx, &appsync.UpdateApiKeyInput{
			ApiId:       aws.String(res.gqlApiId),
			Id:          aws.String(res.apiKeyId),
			Description: aws.String("updated key"),
			Expires:     newExpiry,
		})
		if err != nil {
			return err
		}
		if resp.ApiKey == nil {
			return fmt.Errorf("ApiKey is nil")
		}
		if resp.ApiKey.Description == nil || *resp.ApiKey.Description != "updated key" {
			return fmt.Errorf("description not updated: %v", resp.ApiKey.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApiKey_NonExistent", func() error {
		_, err := client.UpdateApiKey(ctx, &appsync.UpdateApiKeyInput{
			ApiId:       aws.String(res.gqlApiId),
			Id:          aws.String("does-not-exist"),
			Description: aws.String("noop"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiKey", func() error {
		_, err := client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
			ApiId: aws.String(res.gqlApiId),
			Id:    aws.String(res.apiKeyId),
		})
		if err != nil {
			return err
		}
		listResp, err := client.ListApiKeys(ctx, &appsync.ListApiKeysInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		for _, key := range listResp.ApiKeys {
			if key.Id != nil && *key.Id == res.apiKeyId {
				return fmt.Errorf("deleted API key still appears in list")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiKey_NonExistent", func() error {
		_, err := client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
			ApiId: aws.String(res.gqlApiId),
			Id:    aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}

func (r *TestRunner) runAppSyncCacheTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "CreateApiCache", func() error {
		resp, err := client.CreateApiCache(ctx, &appsync.CreateApiCacheInput{
			ApiId:              aws.String(res.gqlApiId),
			Type:               types.ApiCacheTypeSmall,
			Ttl:                300,
			ApiCachingBehavior: types.ApiCachingBehaviorFullRequestCaching,
		})
		if err != nil {
			return err
		}
		if resp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil")
		}
		if resp.ApiCache.Type != types.ApiCacheTypeSmall {
			return fmt.Errorf("expected SMALL type, got %s", resp.ApiCache.Type)
		}
		if resp.ApiCache.Ttl != 300 {
			return fmt.Errorf("expected TTL 300, got %d", resp.ApiCache.Ttl)
		}
		if resp.ApiCache.ApiCachingBehavior != types.ApiCachingBehaviorFullRequestCaching {
			return fmt.Errorf("expected FULL_REQUEST_CACHING, got %s", resp.ApiCache.ApiCachingBehavior)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApiCache", func() error {
		resp, err := client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil")
		}
		if resp.ApiCache.Ttl != 300 {
			return fmt.Errorf("expected TTL 300, got %d", resp.ApiCache.Ttl)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApiCache_NonExistent", func() error {
		_, err := client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "UpdateApiCache", func() error {
		resp, err := client.UpdateApiCache(ctx, &appsync.UpdateApiCacheInput{
			ApiId:              aws.String(res.gqlApiId),
			Type:               types.ApiCacheTypeMedium,
			Ttl:                600,
			ApiCachingBehavior: types.ApiCachingBehaviorPerResolverCaching,
		})
		if err != nil {
			return err
		}
		if resp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil")
		}
		if resp.ApiCache.Type != types.ApiCacheTypeMedium {
			return fmt.Errorf("expected MEDIUM type, got %s", resp.ApiCache.Type)
		}
		if resp.ApiCache.Ttl != 600 {
			return fmt.Errorf("expected TTL 600, got %d", resp.ApiCache.Ttl)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "FlushApiCache", func() error {
		_, err := client.FlushApiCache(ctx, &appsync.FlushApiCacheInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		verifyResp, err := client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return fmt.Errorf("cache should still exist after flush: %w", err)
		}
		if verifyResp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil after flush")
		}
		if verifyResp.ApiCache.Ttl != 600 {
			return fmt.Errorf("expected TTL 600 after flush, got %d", verifyResp.ApiCache.Ttl)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "FlushApiCache_NonExistent", func() error {
		_, err := client.FlushApiCache(ctx, &appsync.FlushApiCacheInput{
			ApiId: aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiCache", func() error {
		_, err := client.DeleteApiCache(ctx, &appsync.DeleteApiCacheInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String(res.gqlApiId),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiCache_NonExistent", func() error {
		_, err := client.DeleteApiCache(ctx, &appsync.DeleteApiCacheInput{
			ApiId: aws.String(res.gqlApiId),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}
