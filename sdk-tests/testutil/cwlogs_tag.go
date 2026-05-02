package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) tagTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "TagResource_Basic", func() error {
		tgName := tc.uniquePrefix("TagGroup")
		if err := tc.createLogGroup(tgName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(tgName)

		arn, err := tc.findLogGroupARN(tgName)
		if err != nil {
			return err
		}

		_, err = tc.client.TagResource(tc.ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: arn,
			Tags: map[string]string{
				"Environment": "test",
				"Team":        "vorpalstacks",
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "ListTagsForResource_Basic", func() error {
		ltName := tc.uniquePrefix("ListTagGroup")
		if err := tc.createLogGroup(ltName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(ltName)

		arn, err := tc.findLogGroupARN(ltName)
		if err != nil {
			return err
		}

		tc.client.TagResource(tc.ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: arn,
			Tags: map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
			},
		})

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &cloudwatchlogs.ListTagsForResourceInput{
			ResourceArn: arn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if tagResp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if tagResp.Tags["Key1"] != "Value1" {
			return fmt.Errorf("Key1 mismatch: got %q", tagResp.Tags["Key1"])
		}
		if tagResp.Tags["Key2"] != "Value2" {
			return fmt.Errorf("Key2 mismatch: got %q", tagResp.Tags["Key2"])
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "UntagResource_Basic", func() error {
		utName := tc.uniquePrefix("UntagGroup")
		if err := tc.createLogGroup(utName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(utName)

		arn, err := tc.findLogGroupARN(utName)
		if err != nil {
			return err
		}

		tc.client.TagResource(tc.ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: arn,
			Tags: map[string]string{
				"RemoveMe":  "yes",
				"KeepMe":    "no",
				"KeepMeToo": "also-no",
			},
		})

		_, err = tc.client.UntagResource(tc.ctx, &cloudwatchlogs.UntagResourceInput{
			ResourceArn: arn,
			TagKeys:     []string{"RemoveMe"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &cloudwatchlogs.ListTagsForResourceInput{
			ResourceArn: arn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if _, ok := tagResp.Tags["RemoveMe"]; ok {
			return fmt.Errorf("RemoveMe tag should have been removed")
		}
		if tagResp.Tags["KeepMe"] != "no" {
			return fmt.Errorf("KeepMe tag should still exist")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "TagLogGroup_Basic", func() error {
		tlgName := tc.uniquePrefix("TagLGGroup")
		if err := tc.createLogGroup(tlgName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(tlgName)

		resourceARN := fmt.Sprintf("arn:aws:logs:%s:000000000000:log-group:%s", tc.region, tlgName)
		_, err := tc.client.TagResource(tc.ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: &resourceARN,
			Tags: map[string]string{
				"TestTag": "yes",
			},
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &cloudwatchlogs.ListTagsForResourceInput{
			ResourceArn: &resourceARN,
		})
		if err != nil {
			return fmt.Errorf("list tags resource: %v", err)
		}
		if tagResp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if tagResp.Tags["TestTag"] != "yes" {
			return fmt.Errorf("TestTag mismatch: got %q", tagResp.Tags["TestTag"])
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "CreateLogGroup_WithTags", func() error {
		cwtName := tc.uniquePrefix("CWTagGroup")
		_, err := tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(cwtName),
			Tags: map[string]string{
				"CreatedBy": "sdk-test",
				"CreatedAt": "2026-04-02",
				"Automated": "true",
			},
		})
		if err != nil {
			return fmt.Errorf("create with tags: %v", err)
		}
		defer tc.deleteLogGroup(cwtName)

		resourceARN := fmt.Sprintf("arn:aws:logs:%s:000000000000:log-group:%s", tc.region, cwtName)
		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &cloudwatchlogs.ListTagsForResourceInput{
			ResourceArn: &resourceARN,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tagResp.Tags) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(tagResp.Tags))
		}
		if tagResp.Tags["CreatedBy"] != "sdk-test" {
			return fmt.Errorf("CreatedBy mismatch")
		}
		return nil
	}))

	return results
}
