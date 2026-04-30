package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

func (r *TestRunner) runAppSyncTagTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

	results = append(results, r.RunTest("appsync", "CreateApi_ForTagging", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:        aws.String(fmt.Sprintf("test-api-tag-%d", uid)),
			EventConfig: minEventConfig(),
			Tags: map[string]string{
				"key1": "value1",
			},
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || resp.Api.ApiArn == nil || resp.Api.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		if resp.Api.Tags == nil || resp.Api.Tags["key1"] != "value1" {
			return fmt.Errorf("initial tag not persisted: %v", resp.Api.Tags)
		}
		res.taggedApiArn = *resp.Api.ApiArn
		res.taggedApiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(res.taggedApiArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["key1"] != "value1" {
			return fmt.Errorf("expected key1=value1, got: %v", resp.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "TagResource", func() error {
		_, err := client.TagResource(ctx, &appsync.TagResourceInput{
			ResourceArn: aws.String(res.taggedApiArn),
			Tags: map[string]string{
				"key2": "value2",
				"key3": "value3",
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(res.taggedApiArn),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(resp.Tags))
		}
		if resp.Tags["key2"] != "value2" || resp.Tags["key3"] != "value3" {
			return fmt.Errorf("new tags not found: %v", resp.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &appsync.UntagResourceInput{
			ResourceArn: aws.String(res.taggedApiArn),
			TagKeys:     []string{"key2"},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(res.taggedApiArn),
		})
		if err != nil {
			return err
		}
		if _, exists := resp.Tags["key2"]; exists {
			return fmt.Errorf("key2 should have been removed: %v", resp.Tags)
		}
		if resp.Tags["key1"] != "value1" || resp.Tags["key3"] != "value3" {
			return fmt.Errorf("remaining tags incorrect: %v", resp.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateChannelNamespace_ForTagging", func() error {
		res.tagNsName = fmt.Sprintf("tag-ns-%d", uid)
		resp, err := client.CreateChannelNamespace(ctx, &appsync.CreateChannelNamespaceInput{
			ApiId: aws.String(res.apiId),
			Name:  aws.String(res.tagNsName),
			Tags: map[string]string{
				"nsKey": "nsValue",
			},
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil || resp.ChannelNamespace.ChannelNamespaceArn == nil {
			return fmt.Errorf("invalid response")
		}
		if resp.ChannelNamespace.Name == nil || *resp.ChannelNamespace.Name != res.tagNsName {
			return fmt.Errorf("expected name %s, got %v", res.tagNsName, resp.ChannelNamespace.Name)
		}
		res.nsArn = *resp.ChannelNamespace.ChannelNamespaceArn
		return nil
	}))

	results = append(results, r.RunTest("appsync", "TagResource_ChannelNamespace", func() error {
		_, err := client.TagResource(ctx, &appsync.TagResourceInput{
			ResourceArn: aws.String(res.nsArn),
			Tags: map[string]string{
				"added": "yes",
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(res.nsArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags["nsKey"] != "nsValue" || resp.Tags["added"] != "yes" {
			return fmt.Errorf("namespace tags incorrect: %v", resp.Tags)
		}
		return nil
	}))

	return results
}
