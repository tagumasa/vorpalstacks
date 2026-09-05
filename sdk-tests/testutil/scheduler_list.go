package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runListTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("scheduler", "ListSchedules", func() error {
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		schedName := tc.uniqueName("ListSched")

		_, err := tc.createSchedule(schedName, "rate(30 minutes)", tc.defaultTarget(rARN))
		defer tc.cleanupSchedule(schedName)

		resp, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{})
		if err != nil {
			return err
		}
		if resp.Schedules == nil {
			return fmt.Errorf("Schedules is nil")
		}
		found := false
		for _, s := range resp.Schedules {
			if s.Name != nil && *s.Name == schedName {
				found = true
				if s.Arn == nil || *s.Arn == "" {
					return fmt.Errorf("found schedule has empty Arn")
				}
				if s.State != types.ScheduleStateEnabled {
					return fmt.Errorf("State mismatch: expected ENABLED, got %q", s.State)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created schedule %q not found in list", schedName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListSchedules_NamePrefix", func() error {
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		prefixName := tc.uniqueName("PrefixSched")

		_, err := tc.createSchedule(prefixName, "rate(30 minutes)", tc.defaultTarget(rARN))
		defer tc.cleanupSchedule(prefixName)

		prefix := prefixName[:len(prefixName)-8]
		resp, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{
			NamePrefix: aws.String(prefix),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, s := range resp.Schedules {
			if s.Name != nil && *s.Name == prefixName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("schedule %q not found with prefix %q", prefixName, prefix)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListSchedules_AllGroups", func() error {
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		groupName := tc.uniqueName("AllGroupsGrp")
		schedName := tc.uniqueName("AllGroupsSched")

		_, err := tc.createScheduleGroup(groupName)
		defer tc.cleanupScheduleGroup(groupName)

		_, err = tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			GroupName:          aws.String(groupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name:      aws.String(schedName),
			GroupName: aws.String(groupName),
		})

		// An unfiltered list must include schedules from non-default
		// groups: GroupName only filters "if specified".
		resp, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, s := range resp.Schedules {
			if s.Name != nil && *s.Name == schedName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("schedule %q in group %q not found in unfiltered list", schedName, groupName)
		}

		// The explicit group filter still scopes to the named group.
		filtered, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{
			GroupName: aws.String(groupName),
		})
		if err != nil {
			return fmt.Errorf("filtered list: %v", err)
		}
		filteredFound := false
		for _, s := range filtered.Schedules {
			if s.Name != nil && *s.Name == schedName {
				filteredFound = true
				break
			}
		}
		if !filteredFound {
			return fmt.Errorf("schedule %q not found under its own group filter", schedName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListSchedules_StateFilter", func() error {
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		schedName := tc.uniqueName("StateFilterSched")

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			State:              types.ScheduleStateDisabled,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.cleanupSchedule(schedName)

		resp, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{
			State: types.ScheduleStateDisabled,
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		found := false
		for _, s := range resp.Schedules {
			if s.Name != nil && *s.Name == schedName {
				found = true
				if s.State != types.ScheduleStateDisabled {
					return fmt.Errorf("filtered schedule state mismatch: expected DISABLED, got %q", s.State)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("disabled schedule %q not found with state filter", schedName)
		}
		return nil
	}))

	// Addressing a schedule group that does not exist is
	// ResourceNotFoundException on every operation that scopes by group.
	results = append(results, tc.runner.RunTest("scheduler", "ScheduleGroup_NonExistent", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"DeleteScheduleGroup", func() error {
				_, err := tc.client.DeleteScheduleGroup(tc.ctx, &scheduler.DeleteScheduleGroupInput{
					Name: aws.String("nonexistent-group-xyz"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"ListSchedules", func() error {
				_, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{
					GroupName: aws.String("no-such-group-xyz"),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
		}
		return nil
	}))

	// Both list operations share the MaxResults range contract: 0 and any
	// value above 100 are ValidationExceptions.
	results = append(results, tc.runner.RunTest("scheduler", "List_MaxResultsOutOfRange", func() error {
		rows := []struct {
			name string
			list func(maxResults int32) error
		}{
			{"ListSchedules", func(maxResults int32) error {
				_, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{
					MaxResults: aws.Int32(maxResults),
				})
				return AssertErrorContains(err, "ValidationException")
			}},
			{"ListScheduleGroups", func(maxResults int32) error {
				_, err := tc.client.ListScheduleGroups(tc.ctx, &scheduler.ListScheduleGroupsInput{
					MaxResults: aws.Int32(maxResults),
				})
				return AssertErrorContains(err, "ValidationException")
			}},
		}
		for _, row := range rows {
			for _, maxResults := range []int32{0, 101} {
				if err := row.list(maxResults); err != nil {
					return fmt.Errorf("%s MaxResults=%d: %w", row.name, maxResults, err)
				}
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "ListSchedules_InvalidStateFilter", func() error {
		_, err := tc.client.ListSchedules(tc.ctx, &scheduler.ListSchedulesInput{
			State: types.ScheduleState("BOGUS"),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
