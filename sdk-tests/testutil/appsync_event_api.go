package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncEventApiTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

	results = append(results, r.RunTest("appsync", "CreateApi", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:        aws.String(fmt.Sprintf("test-api-%d", uid)),
			EventConfig: minEventConfig(),
		})
		if err != nil {
			return err
		}
		if resp.Api == nil {
			return fmt.Errorf("api is nil")
		}
		if resp.Api.ApiId == nil || *resp.Api.ApiId == "" {
			return fmt.Errorf("ApiId is empty")
		}
		res.apiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateApi_WithTags", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:        aws.String(fmt.Sprintf("test-api-tags-%d", uid)),
			EventConfig: minEventConfig(),
			Tags: map[string]string{
				"env":  "test",
				"team": "platform",
			},
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || resp.Api.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		if len(resp.Api.Tags) != 2 || resp.Api.Tags["env"] != "test" || resp.Api.Tags["team"] != "platform" {
			return fmt.Errorf("tags not persisted: %v", resp.Api.Tags)
		}
		res.tagsApiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateApi_WithOwnerContact", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:         aws.String(fmt.Sprintf("test-api-owner-%d", uid)),
			EventConfig:  minEventConfig(),
			OwnerContact: aws.String("test@example.com"),
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || resp.Api.OwnerContact == nil || *resp.Api.OwnerContact != "test@example.com" {
			return fmt.Errorf("ownerContact not set: %v", resp.Api.OwnerContact)
		}
		res.ownerApiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApi", func() error {
		resp, err := client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(res.apiId)})
		if err != nil {
			return err
		}
		if resp.Api == nil {
			return fmt.Errorf("api is nil")
		}
		if *resp.Api.ApiId != res.apiId {
			return fmt.Errorf("expected ApiId %s, got %s", res.apiId, *resp.Api.ApiId)
		}
		if resp.Api.EventConfig == nil {
			return fmt.Errorf("EventConfig is nil")
		}
		if len(resp.Api.EventConfig.AuthProviders) != 1 {
			return fmt.Errorf("expected 1 auth provider, got %d", len(resp.Api.EventConfig.AuthProviders))
		}
		if resp.Api.ApiArn == nil {
			return fmt.Errorf("ApiArn is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApi_NonExistent", func() error {
		_, err := client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String("does-not-exist")})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListApis", func() error {
		apis, err := res.allApis()
		if err != nil {
			return err
		}
		if len(apis) < 3 {
			return fmt.Errorf("expected at least 3 APIs, got %d", len(apis))
		}
		if containsID(apis, func(api *types.Api) bool {
			return api.ApiId != nil && *api.ApiId == res.apiId
		}) == nil {
			return fmt.Errorf("created API not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApis_WithPagination", func() error {
		resp, err := client.ListApis(ctx, &appsync.ListApisInput{MaxResults: 1})
		if err != nil {
			return err
		}
		if len(resp.Apis) != 1 {
			return fmt.Errorf("expected 1 API with maxResults=1, got %d", len(resp.Apis))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApi", func() error {
		newName := fmt.Sprintf("updated-api-%d", uid)
		resp, err := client.UpdateApi(ctx, &appsync.UpdateApiInput{
			ApiId:       aws.String(res.apiId),
			Name:        aws.String(newName),
			EventConfig: minEventConfig(),
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || *resp.Api.Name != newName {
			return fmt.Errorf("expected name %s, got %v", newName, resp.Api.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApi_NonExistent", func() error {
		_, err := client.UpdateApi(ctx, &appsync.UpdateApiInput{
			ApiId:       aws.String("does-not-exist"),
			Name:        aws.String("noop"),
			EventConfig: minEventConfig(),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}

func (r *TestRunner) runAppSyncEventApiPaginationTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "ListApis_NextTokenFollowUp", func() error {
		items, err := paginate(func(next *string) ([]types.Api, *string, error) {
			resp, err := client.ListApis(ctx, &appsync.ListApisInput{
				MaxResults: 1,
				NextToken:  next,
			})
			if err != nil {
				return nil, nil, fmt.Errorf("list apis page: %v", err)
			}
			return resp.Apis, resp.NextToken, nil
		})
		if err != nil {
			return err
		}
		// MaxResults=1 forces one item per page, so at least two items
		// prove the NextToken follow-up crossed a page boundary.
		if len(items) < 2 {
			return fmt.Errorf("expected at least 2 pages for ListApis with MaxResults=1, got %d", len(items))
		}
		return nil
	}))

	return results
}
