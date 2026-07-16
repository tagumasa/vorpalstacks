package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runTagTests() []TestResult {
	var results []TestResult

	rn, rARN := tc.createIAMRole()
	defer tc.deleteIAMRole(rn)
	schedName := tc.uniqueName("TagSched")

	_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
		Name:               aws.String(schedName),
		ScheduleExpression: aws.String("rate(30 minutes)"),
		Target:             tc.defaultTarget(rARN),
		FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
	})
	if err != nil {
		results = append(results, TestResult{Service: "scheduler", TestName: "TagSetup", Status: "FAIL", Error: fmt.Sprintf("create: %v", err)})
		return results
	}
	defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})

	scheduleARN := tc.scheduleARN(schedName)

	results = append(results, tc.runner.RunTest("scheduler", "TagResource", func() error {
		_, err := tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Project"), Value: aws.String("vorpalstacks")},
			},
		})
		if err != nil {
			return err
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: aws.String(scheduleARN),
		})
		if err != nil {
			return fmt.Errorf("list tags after tag: %v", err)
		}
		envFound := false
		projFound := false
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				envFound = true
				if t.Value == nil || *t.Value != "test" {
					return fmt.Errorf("Environment value mismatch: got %q", aws.ToString(t.Value))
				}
			}
			if t.Key != nil && *t.Key == "Project" {
				projFound = true
				if t.Value == nil || *t.Value != "vorpalstacks" {
					return fmt.Errorf("Project value mismatch: got %q", aws.ToString(t.Value))
				}
			}
		}
		if !envFound || !projFound {
			return fmt.Errorf("not all tags found: env=%v proj=%v", envFound, projFound)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListTagsForResource", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: aws.String(scheduleARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("Tags is nil")
		}
		if len(resp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.Tags))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UntagResource", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &scheduler.UntagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: aws.String(scheduleARN),
		})
		if err != nil {
			return fmt.Errorf("list tags after untag: %v", err)
		}
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				return fmt.Errorf("Environment tag should have been removed")
			}
		}
		projFound := false
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Project" {
				projFound = true
			}
		}
		if !projFound {
			return fmt.Errorf("Project tag should still exist")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "TagResource_ScheduleGroup", func() error {
		tagGroupName := tc.uniqueName("TagGroup")
		groupResp, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(tagGroupName),
		})
		if err != nil {
			return err
		}

		_, err = tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("prod")},
			},
		})
		if err != nil {
			tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(tagGroupName)})
			return err
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
		})
		if err != nil {
			tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(tagGroupName)})
			return fmt.Errorf("list tags: %v", err)
		}
		found := false
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Env" {
				found = true
				if t.Value == nil || *t.Value != "prod" {
					tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(tagGroupName)})
					return fmt.Errorf("Env value mismatch: got %q", aws.ToString(t.Value))
				}
				break
			}
		}
		tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(tagGroupName)})
		if !found {
			return fmt.Errorf("tag Env not found on schedule group")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListTagsForResource_ScheduleGroup", func() error {
		groupName := tc.uniqueName("TagListGroup")
		groupResp, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(groupName),
			Tags: []types.Tag{
				{Key: aws.String("Team"), Value: aws.String("platform")},
			},
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		found := false
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Team" {
				found = true
				if t.Value == nil || *t.Value != "platform" {
					return fmt.Errorf("Team value mismatch: got %q", aws.ToString(t.Value))
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("Team tag not found on group created with tags")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UntagResource_NonExistentKey", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &scheduler.UntagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			TagKeys:     []string{"NonExistentKey"},
		})
		if err != nil {
			return fmt.Errorf("untag non-existent key should not error: %v", err)
		}
		return nil
	}))

	return results
}
