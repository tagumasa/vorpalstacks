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

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})

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

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(prefixName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(prefixName)})

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
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName)})

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

	return results
}
