package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestContext) testTagging() []TestResult {
	var results []TestResult

	// Shared work group fixture for the tag scenarios; the deferred delete
	// runs once every tagging scenario below has completed. A setup failure
	// surfaces as one FAIL row named after the step it replaced.
	tagWorkGroupName := tc.uniqueName("tag-wg")
	if err := tc.createWorkGroup(tagWorkGroupName, nil); err != nil {
		return append(results, TestResult{
			Service:  "athena",
			TestName: "TagResource_CreateWG",
			Status:   "FAIL",
			Error:    fmt.Sprintf("work group setup failed: %v", err),
		})
	}
	defer tc.deleteWorkGroup(tagWorkGroupName)

	results = append(results, tc.runner.RunTest("athena", "TagResource", func() error {
		_, err := tc.client.TagResource(tc.ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(tc.workgroupARN(tagWorkGroupName)),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("athena")},
			},
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "TagResource_ReservedPrefixRejected", func() error {
		_, err := tc.client.TagResource(tc.ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(tc.workgroupARN(tagWorkGroupName)),
			Tags: []types.Tag{
				{Key: aws.String("aws:reserved"), Value: aws.String("v")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag key")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListTagsForResource", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(tc.workgroupARN(tagWorkGroupName)),
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

	// Untag then verify: env must be gone and team must remain.
	results = append(results, tc.runner.RunTest("athena", "UntagResource", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &athena.UntagResourceInput{
			ResourceARN: aws.String(tc.workgroupARN(tagWorkGroupName)),
			TagKeys:     []string{"env"},
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.ListTagsForResource(tc.ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(tc.workgroupARN(tagWorkGroupName)),
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

	// Data catalog tagging: create, tag, verify in one scenario with
	// best-effort cleanup.
	tagCatalogName := tc.uniqueName("tag-cat")
	results = append(results, tc.runner.RunTest("athena", "TagResource_DataCatalog", func() error {
		if err := tc.createDataCatalog(tagCatalogName, "Catalog for tag test"); err != nil {
			return err
		}
		defer tc.deleteDataCatalog(tagCatalogName)

		_, err := tc.client.TagResource(tc.ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(tc.datacatalogARN(tagCatalogName)),
			Tags: []types.Tag{
				{Key: aws.String("purpose"), Value: aws.String("testing")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(tc.datacatalogARN(tagCatalogName)),
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

	return results
}
