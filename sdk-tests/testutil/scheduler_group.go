package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runGroupTests() []TestResult {
	var results []TestResult

	// Verify the default schedule group exists after server start.
	results = append(results, tc.runner.RunTest("scheduler", "GetScheduleGroup_Default", func() error {
		resp, err := tc.client.GetScheduleGroup(tc.ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String("default"),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "default" {
			return fmt.Errorf("expected default group name %q, got %q", "default", aws.ToString(resp.Name))
		}
		if resp.State != types.ScheduleGroupStateActive {
			return fmt.Errorf("expected default group state ACTIVE, got %q", resp.State)
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("default group Arn is nil or empty")
		}
		return nil
	}))

	groupName := tc.uniqueName("TestGroup")

	results = append(results, tc.runner.RunTest("scheduler", "CreateScheduleGroup", func() error {
		resp, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp.ScheduleGroupArn == nil || *resp.ScheduleGroupArn == "" {
			return fmt.Errorf("ScheduleGroupArn is nil or empty")
		}
		expected := fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule-group/%s", tc.region, tc.accountID, groupName)
		if *resp.ScheduleGroupArn != expected {
			return fmt.Errorf("ScheduleGroupArn mismatch: expected %q, got %q", expected, *resp.ScheduleGroupArn)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "GetScheduleGroup", func() error {
		resp, err := tc.client.GetScheduleGroup(tc.ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != groupName {
			return fmt.Errorf("Name mismatch: expected %q, got %q", groupName, aws.ToString(resp.Name))
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		if resp.State != types.ScheduleGroupStateActive {
			return fmt.Errorf("State mismatch: expected ACTIVE, got %q", resp.State)
		}
		if resp.CreationDate == nil {
			return fmt.Errorf("CreationDate is nil")
		}
		if resp.LastModificationDate == nil {
			return fmt.Errorf("LastModificationDate is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListScheduleGroups", func() error {
		resp, err := tc.client.ListScheduleGroups(tc.ctx, &scheduler.ListScheduleGroupsInput{})
		if err != nil {
			return err
		}
		if resp.ScheduleGroups == nil {
			return fmt.Errorf("ScheduleGroups is nil")
		}
		if len(resp.ScheduleGroups) == 0 {
			return fmt.Errorf("ScheduleGroups is empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListScheduleGroups_ContainsCreated", func() error {
		resp, err := tc.client.ListScheduleGroups(tc.ctx, &scheduler.ListScheduleGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, g := range resp.ScheduleGroups {
			if g.Name != nil && *g.Name == groupName {
				found = true
				if g.Arn == nil || *g.Arn == "" {
					return fmt.Errorf("found group has empty Arn")
				}
				if g.State != types.ScheduleGroupStateActive {
					return fmt.Errorf("State mismatch: expected ACTIVE, got %q", g.State)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created group %q not found in list", groupName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_WithGroupName", func() error {
		schedName := tc.uniqueName("GroupSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			GroupName:          aws.String(groupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(schedName),
			GroupName: aws.String(groupName),
		})

		getResp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{
			Name:      aws.String(schedName),
			GroupName: aws.String(groupName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.GroupName == nil || *getResp.GroupName != groupName {
			return fmt.Errorf("GroupName mismatch: expected %q, got %q", groupName, aws.ToString(getResp.GroupName))
		}
		if getResp.Name == nil || *getResp.Name != schedName {
			return fmt.Errorf("Name mismatch: expected %q, got %q", schedName, aws.ToString(getResp.Name))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateScheduleGroup_DuplicateName", func() error {
		dupGroupName := tc.uniqueName("DupGroup")
		_, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(dupGroupName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(dupGroupName)})

		_, err = tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(dupGroupName),
		})
		if err := AssertErrorContains(err, "ConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "DeleteScheduleGroup_NonExistent", func() error {
		_, err := tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{
			Name: aws.String("nonexistent-group-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// Deleting is eventually consistent: the delete is acknowledged
	// immediately and the group disappears once the cascade completes.
	results = append(results, tc.runner.RunTest("scheduler", "DeleteScheduleGroup", func() error {
		delGroupName := tc.uniqueName("DelGroup")
		_, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(delGroupName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		_, err = tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{
			Name: aws.String(delGroupName),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		deadline := time.Now().Add(15 * time.Second)
		for {
			_, err := tc.client.GetScheduleGroup(tc.ctx, &scheduler.GetScheduleGroupInput{
				Name: aws.String(delGroupName),
			})
			if err != nil {
				if err := AssertErrorContains(err, "ResourceNotFoundException"); err == nil {
					return nil
				}
				return fmt.Errorf("get after delete: %v", err)
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("group %q still exists after the cascade window", delGroupName)
			}
			time.Sleep(200 * time.Millisecond)
		}
	}))

	// Deleting a group cascades: the group goes to DELETING, its schedules
	// are deleted, and only then does the group disappear.
	results = append(results, tc.runner.RunTest("scheduler", "DeleteScheduleGroup_CascadesToSchedules", func() error {
		notEmptyGroup := tc.uniqueName("NotEmptyGrp")
		_, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(notEmptyGroup),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}

		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		schedName := tc.uniqueName("InsideSched")
		_, err = tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			GroupName:          aws.String(notEmptyGroup),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(tc.lambdaARN()),
				RoleArn: aws.String(rARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create schedule: %v", err)
		}

		if _, err := tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{
			Name: aws.String(notEmptyGroup),
		}); err != nil {
			return fmt.Errorf("delete non-empty group: %v", err)
		}

		// The group (and its schedules) disappear once the cascade
		// completes; the delete is eventually consistent.
		deadline := time.Now().Add(15 * time.Second)
		for {
			_, err := tc.client.GetScheduleGroup(tc.ctx, &scheduler.GetScheduleGroupInput{
				Name: aws.String(notEmptyGroup),
			})
			if err != nil {
				if err := AssertErrorContains(err, "ResourceNotFoundException"); err == nil {
					break
				}
				return fmt.Errorf("get deleting group: %v", err)
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("group %q still exists after the cascade window", notEmptyGroup)
			}
			time.Sleep(200 * time.Millisecond)
		}

		_, err = tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(schedName),
			GroupName: aws.String(notEmptyGroup),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("cascade must delete the member schedule: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "DeleteScheduleGroup_DefaultRejected", func() error {
		_, err := tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{
			Name: aws.String("default"),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		// The default group must survive the rejected delete.
		resp, err := tc.client.GetScheduleGroup(tc.ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String("default"),
		})
		if err != nil {
			return fmt.Errorf("get default after rejected delete: %v", err)
		}
		if resp.State != types.ScheduleGroupStateActive {
			return fmt.Errorf("default group state = %q, want ACTIVE", resp.State)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListScheduleGroups_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgGroups []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagGroup-%s-%d", pgTs, i)
			_, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
				Name: aws.String(name),
			})
			if err != nil {
				for _, gn := range pgGroups {
					tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(gn)})
				}
				return fmt.Errorf("create schedule group %s: %v", name, err)
			}
			pgGroups = append(pgGroups, name)
		}

		var allGroups []string
		var nextToken *string
		for {
			resp, err := tc.client.ListScheduleGroups(tc.ctx, &scheduler.ListScheduleGroupsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, gn := range pgGroups {
					tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(gn)})
				}
				return fmt.Errorf("list schedule groups page: %v", err)
			}
			for _, g := range resp.ScheduleGroups {
				if g.Name != nil && strings.Contains(*g.Name, "PagGroup-"+pgTs) {
					allGroups = append(allGroups, *g.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, gn := range pgGroups {
			tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(gn)})
		}
		if len(allGroups) != 5 {
			return fmt.Errorf("expected 5 paginated schedule groups, got %d", len(allGroups))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListScheduleGroups_NamePrefix", func() error {
		prefixTs := fmt.Sprintf("%d", time.Now().UnixNano())
		prefix := fmt.Sprintf("PrefixGrp-%s-", prefixTs)
		var created []string
		for i := 0; i < 3; i++ {
			name := fmt.Sprintf("%s%d", prefix, i)
			_, err := tc.client.CreateScheduleGroup(tc.ctx, &scheduler.CreateScheduleGroupInput{
				Name: aws.String(name),
			})
			if err != nil {
				for _, gn := range created {
					tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(gn)})
				}
				return fmt.Errorf("create group %s: %v", name, err)
			}
			created = append(created, name)
		}

		resp, err := tc.client.ListScheduleGroups(tc.ctx, &scheduler.ListScheduleGroupsInput{
			NamePrefix: aws.String(prefix),
		})
		if err != nil {
			for _, gn := range created {
				tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(gn)})
			}
			return fmt.Errorf("list with prefix: %v", err)
		}

		for _, gn := range created {
			tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(gn)})
		}

		matched := 0
		for _, g := range resp.ScheduleGroups {
			if g.Name != nil && strings.HasPrefix(*g.Name, prefix) {
				matched++
			}
		}
		if matched != 3 {
			return fmt.Errorf("expected 3 groups with prefix %q, found %d", prefix, matched)
		}
		return nil
	}))

	defer tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})

	return results
}
