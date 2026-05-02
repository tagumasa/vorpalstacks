package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestCtx) testTagging() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	tagWorkGroupName := fmt.Sprintf("tag-wg-%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "TagResource_CreateWG", func() error {
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(tagWorkGroupName),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "TagResource", func() error {
		_, err := client.TagResource(ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", tc.runner.region, tagWorkGroupName)),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("athena")},
			},
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", tc.runner.region, tagWorkGroupName)),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["env"] != "test" {
			return fmt.Errorf("expected env=test, got %q", tagMap["env"])
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &athena.UntagResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", tc.runner.region, tagWorkGroupName)),
			TagKeys:     []string{"env"},
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "ListTagsForResource_AfterUntag", func() error {
		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", tc.runner.region, tagWorkGroupName)),
		})
		if err != nil {
			return err
		}
		for _, t := range resp.Tags {
			if aws.ToString(t.Key) == "env" {
				return fmt.Errorf("env tag should have been removed")
			}
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["team"] != "athena" {
			return fmt.Errorf("expected team=athena to remain, got %q", tagMap["team"])
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteWorkGroup_TagCleanup", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String(tagWorkGroupName),
		})
		return err
	}))

	tagCatalogName := fmt.Sprintf("tag-cat-%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "TagResource_DataCatalog", func() error {
		_, err := client.CreateDataCatalog(ctx, &athena.CreateDataCatalogInput{
			Name:        aws.String(tagCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("Catalog for tag test"),
		})
		if err != nil {
			return err
		}
		_, err = client.TagResource(ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:datacatalog/%s", tc.runner.region, tagCatalogName)),
			Tags: []types.Tag{
				{Key: aws.String("purpose"), Value: aws.String("testing")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:datacatalog/%s", tc.runner.region, tagCatalogName)),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 1 {
			return fmt.Errorf("expected at least 1 tag on datacatalog")
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["purpose"] != "testing" {
			return fmt.Errorf("expected purpose=testing, got %q", tagMap["purpose"])
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteDataCatalog_TagCleanup", func() error {
		_, err := client.DeleteDataCatalog(ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(tagCatalogName),
		})
		return err
	}))

	return results
}
