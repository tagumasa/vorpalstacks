package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runTagTests() []TestResult {
	var results []TestResult

	// Only schedule groups carry tags: the TagResourceArn pattern accepts
	// schedule-group ARNs only, so every tag test targets a group.
	groupName := tc.uniqueName("TagGroup")

	groupResp, err := tc.createScheduleGroup(groupName)
	if err != nil {
		results = append(results, TestResult{Service: "scheduler", TestName: "TagSetup", Status: "FAIL", Error: fmt.Sprintf("create group: %v", err)})
		return results
	}
	defer tc.cleanupScheduleGroup(groupName)

	scheduleGroupARN := groupResp.ScheduleGroupArn

	results = append(results, tc.runner.RunTest("scheduler", "TagResource", func() error {
		_, err := tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
			ResourceArn: scheduleGroupARN,
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Project"), Value: aws.String("vorpalstacks")},
			},
		})
		if err != nil {
			return err
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: scheduleGroupARN,
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
			ResourceArn: scheduleGroupARN,
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
			ResourceArn: scheduleGroupARN,
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: scheduleGroupARN,
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
		groupResp, err := tc.createScheduleGroup(tagGroupName)
		if err != nil {
			return err
		}
		defer tc.cleanupScheduleGroup(tagGroupName)

		_, err = tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("prod")},
			},
		})
		if err != nil {
			return err
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		found := false
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Env" {
				found = true
				if t.Value == nil || *t.Value != "prod" {
					return fmt.Errorf("Env value mismatch: got %q", aws.ToString(t.Value))
				}
				break
			}
		}
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
		defer tc.cleanupScheduleGroup(groupName)

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

	// The tag-set constraints (at most 200 tags, tag-key length at most
	// 128) share one ValidationException contract on TagResource.
	results = append(results, tc.runner.RunTest("scheduler", "TagResource_InvalidTagSetRejected", func() error {
		rows := []struct {
			name string
			tags func() []types.Tag
		}{
			{"TooManyTags", func() []types.Tag {
				tooMany := make([]types.Tag, 201)
				for i := range tooMany {
					tooMany[i] = types.Tag{Key: aws.String(fmt.Sprintf("k%03d", i)), Value: aws.String("v")}
				}
				return tooMany
			}},
			{"KeyTooLong", func() []types.Tag {
				longKey := make([]byte, 129)
				for i := range longKey {
					longKey[i] = 'k'
				}
				return []types.Tag{{Key: aws.String(string(longKey)), Value: aws.String("v")}}
			}},
		}
		for _, row := range rows {
			groupName := tc.uniqueName(row.name)
			groupResp, err := tc.createScheduleGroup(groupName)
			if err != nil {
				return fmt.Errorf("%s: create group: %v", row.name, err)
			}
			defer tc.cleanupScheduleGroup(groupName)

			_, err = tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
				ResourceArn: groupResp.ScheduleGroupArn,
				Tags:        row.tags(),
			})
			if err := AssertErrorContains(err, "ValidationException"); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "TagResource_ScheduleArnRejected", func() error {
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		schedName := tc.uniqueName("TagSchedArn")

		_, err := tc.createSchedule(schedName, "rate(30 minutes)", tc.defaultTarget(rARN))
		defer tc.cleanupSchedule(schedName)

		// Schedule ARNs are outside the TagResourceArn pattern: only
		// schedule groups can be tagged.
		_, err = tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
			ResourceArn: aws.String(tc.scheduleARN(schedName)),
			Tags:        []types.Tag{{Key: aws.String("Invalid"), Value: aws.String("target")}},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "TagResource_InvalidArnCharsRejected", func() error {
		// The TagResourceArn pattern constrains the group-name portion to
		// [0-9a-zA-Z-_.]+; a malformed ARN is a validation failure, not
		// a missing group.
		_, err := tc.client.TagResource(tc.ctx, &scheduler.TagResourceInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule-group/bad@name!", tc.region, tc.accountID)),
			Tags:        []types.Tag{{Key: aws.String("Invalid"), Value: aws.String("chars")}},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UntagResource_NonExistentKey", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &scheduler.UntagResourceInput{
			ResourceArn: scheduleGroupARN,
			TagKeys:     []string{"NonExistentKey"},
		})
		if err != nil {
			return fmt.Errorf("untag non-existent key should not error: %v", err)
		}
		return nil
	}))

	return results
}
