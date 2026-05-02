package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

func (r *TestRunner) runCloudTrailTagTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "Tags_Lifecycle", func() error {
		name := fmt.Sprintf("tagcycle-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("tagcycle-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = tc.client.AddTags(tc.ctx, &cloudtrail.AddTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("Value1")},
				{Key: aws.String("Key2"), Value: aws.String("Value2")},
				{Key: aws.String("Key3"), Value: aws.String("Value3")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		listResp, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.ResourceTagList) == 0 || len(listResp.ResourceTagList[0].TagsList) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(listResp.ResourceTagList[0].TagsList))
		}

		_, err = tc.client.RemoveTags(tc.ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Key2")},
			},
		})
		if err != nil {
			return fmt.Errorf("remove tags: %v", err)
		}

		listResp2, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags after remove: %v", err)
		}
		if len(listResp2.ResourceTagList[0].TagsList) != 2 {
			return fmt.Errorf("expected 2 tags after remove, got %d", len(listResp2.ResourceTagList[0].TagsList))
		}

		_, err = tc.client.AddTags(tc.ctx, &cloudtrail.AddTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Key1"), Value: aws.String("UpdatedValue1")},
			},
		})
		if err != nil {
			return fmt.Errorf("update tag: %v", err)
		}

		listResp3, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags after update: %v", err)
		}
		found := false
		for _, t := range listResp3.ResourceTagList[0].TagsList {
			if t.Key != nil && *t.Key == "Key1" && t.Value != nil && *t.Value == "UpdatedValue1" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected Key1=UpdatedValue1 after update")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail_WithTags", func() error {
		name := fmt.Sprintf("tagtrail-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("tagtrail-bucket"),
			TagsList: []types.Tag{
				{Key: aws.String("CreatedBy"), Value: aws.String("sdk-test")},
				{Key: aws.String("Env"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return fmt.Errorf("create with tags: %v", err)
		}

		resp, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail: %v", err)
		}

		listResp, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*resp.Trail.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(listResp.ResourceTagList[0].TagsList) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(listResp.ResourceTagList[0].TagsList))
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "AddTags", func() error {
		name := fmt.Sprintf("addtags-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("addtags-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = tc.client.AddTags(tc.ctx, &cloudtrail.AddTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("test-user")},
			},
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags after add: %v", err)
		}
		if len(listResp.ResourceTagList) == 0 {
			return fmt.Errorf("expected resource tag entry")
		}
		tags := listResp.ResourceTagList[0].TagsList
		if len(tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tags))
		}
		tagMap := make(map[string]string)
		for _, t := range tags {
			if t.Key != nil && t.Value != nil {
				tagMap[*t.Key] = *t.Value
			}
		}
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("expected Environment=test, got %s", tagMap["Environment"])
		}
		if tagMap["Owner"] != "test-user" {
			return fmt.Errorf("expected Owner=test-user, got %s", tagMap["Owner"])
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListTags", func() error {
		name := fmt.Sprintf("listtags-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("listtags-bucket"),
			TagsList: []types.Tag{
				{Key: aws.String("ListTest"), Value: aws.String("yes")},
			},
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return err
		}
		if len(resp.ResourceTagList) == 0 {
			return fmt.Errorf("expected resource tag entry")
		}
		tags := resp.ResourceTagList[0].TagsList
		if len(tags) != 1 {
			return fmt.Errorf("expected 1 tag, got %d", len(tags))
		}
		if tags[0].Key == nil || *tags[0].Key != "ListTest" {
			return fmt.Errorf("expected Key=ListTest, got %v", tags[0].Key)
		}
		if tags[0].Value == nil || *tags[0].Value != "yes" {
			return fmt.Errorf("expected Value=yes, got %v", tags[0].Value)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "RemoveTags", func() error {
		name := fmt.Sprintf("rmtags-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("rmtags-bucket"),
			TagsList: []types.Tag{
				{Key: aws.String("RemoveMe"), Value: aws.String("please")},
				{Key: aws.String("KeepMe"), Value: aws.String("yes")},
			},
		})
		if err != nil {
			return err
		}

		_, err = tc.client.RemoveTags(tc.ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: createResp.TrailARN,
			TagsList: []types.Tag{
				{Key: aws.String("RemoveMe")},
			},
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{*createResp.TrailARN},
		})
		if err != nil {
			return fmt.Errorf("list tags after remove: %v", err)
		}
		tags := listResp.ResourceTagList[0].TagsList
		if len(tags) != 1 {
			return fmt.Errorf("expected 1 tag after remove, got %d", len(tags))
		}
		if tags[0].Key == nil || *tags[0].Key != "KeepMe" {
			return fmt.Errorf("expected Key=KeepMe, got %v", tags[0].Key)
		}
		if tags[0].Value == nil || *tags[0].Value != "yes" {
			return fmt.Errorf("expected Value=yes, got %v", tags[0].Value)
		}
		return nil
	}))

	return results
}
