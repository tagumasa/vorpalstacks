package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func (r *TestRunner) runSSMTag(tc *ssmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("ssm", "AddTagsToResource_ListTagsForResource", func() error {
		name := tc.uniqueName("/tag")
		_, err := tc.putParam(name, "tagged", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.client.AddTagsToResource(tc.ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("platform")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		resp, err := tc.client.ListTagsForResource(tc.ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(resp.TagList) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(resp.TagList))
		}
		tagMap := make(map[string]string)
		for _, t := range resp.TagList {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["env"] != "test" || tagMap["team"] != "platform" {
			return fmt.Errorf("tag values mismatch: %v", tagMap)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "RemoveTagsFromResource", func() error {
		name := tc.uniqueName("/tag-rm")
		_, err := tc.putParam(name, "tagged", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.client.AddTagsToResource(tc.ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("keep"), Value: aws.String("yes")},
				{Key: aws.String("remove"), Value: aws.String("me")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		_, err = tc.client.RemoveTagsFromResource(tc.ctx, &ssm.RemoveTagsFromResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			TagKeys:      []string{"remove"},
		})
		if err != nil {
			return fmt.Errorf("remove tags: %v", err)
		}

		resp, err := tc.client.ListTagsForResource(tc.ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(resp.TagList) != 1 {
			return fmt.Errorf("expected 1 tag after removal, got %d", len(resp.TagList))
		}
		if aws.ToString(resp.TagList[0].Key) != "keep" {
			return fmt.Errorf("expected 'keep' tag, got %q", aws.ToString(resp.TagList[0].Key))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "AddTagsToResource_Merge", func() error {
		name := tc.uniqueName("/tag-merge")
		_, err := tc.putParam(name, "tagged", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		_, err = tc.client.AddTagsToResource(tc.ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("a"), Value: aws.String("1")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags 1: %v", err)
		}

		_, err = tc.client.AddTagsToResource(tc.ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("b"), Value: aws.String("2")},
				{Key: aws.String("a"), Value: aws.String("updated")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags 2: %v", err)
		}

		resp, err := tc.client.ListTagsForResource(tc.ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		tagMap := make(map[string]string)
		for _, t := range resp.TagList {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if len(tagMap) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagMap))
		}
		if tagMap["a"] != "updated" {
			return fmt.Errorf("expected 'a'='updated', got %q", tagMap["a"])
		}
		if tagMap["b"] != "2" {
			return fmt.Errorf("expected 'b'='2', got %q", tagMap["b"])
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "ListTagsForResource_Empty", func() error {
		name := tc.uniqueName("/tag-empty")
		_, err := tc.putParam(name, "notags", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.deleteParam(name)

		resp, err := tc.client.ListTagsForResource(tc.ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(resp.TagList) != 0 {
			return fmt.Errorf("expected 0 tags, got %d", len(resp.TagList))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "AddTagsToResource_NonExistent", func() error {
		_, err := tc.client.AddTagsToResource(tc.ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String("/nonexistent/param-xyz"),
			Tags: []types.Tag{
				{Key: aws.String("k"), Value: aws.String("v")},
			},
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "RemoveTagsFromResource_NonExistent", func() error {
		_, err := tc.client.RemoveTagsFromResource(tc.ctx, &ssm.RemoveTagsFromResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String("/nonexistent/param-xyz"),
			TagKeys:      []string{"k"},
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
