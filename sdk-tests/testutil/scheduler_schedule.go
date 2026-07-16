package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func (tc *schedTestContext) runScheduleTests() []TestResult {
	var results []TestResult

	roleName, roleARN := tc.createIAMRole()
	defer tc.deleteIAMRole(roleName)

	scheduleName := tc.uniqueName("TestSchedule")
	target := tc.defaultTarget(roleARN)

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule", func() error {
		resp, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(scheduleName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             target,
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return err
		}
		if resp.ScheduleArn == nil || *resp.ScheduleArn == "" {
			return fmt.Errorf("ScheduleArn is nil or empty")
		}
		expected := fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/default/%s", tc.region, scheduleName)
		if *resp.ScheduleArn != expected {
			return fmt.Errorf("ScheduleArn mismatch: expected %q, got %q", expected, *resp.ScheduleArn)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "GetSchedule", func() error {
		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != scheduleName {
			return fmt.Errorf("Name mismatch: expected %q, got %q", scheduleName, aws.ToString(resp.Name))
		}
		if resp.ScheduleExpression == nil || *resp.ScheduleExpression != "rate(30 minutes)" {
			return fmt.Errorf("ScheduleExpression mismatch: got %q", aws.ToString(resp.ScheduleExpression))
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		if resp.State != types.ScheduleStateEnabled {
			return fmt.Errorf("State mismatch: expected ENABLED, got %q", resp.State)
		}
		if resp.GroupName == nil || *resp.GroupName != "default" {
			return fmt.Errorf("GroupName mismatch: expected 'default', got %q", aws.ToString(resp.GroupName))
		}
		if resp.Target == nil {
			return fmt.Errorf("Target is nil")
		}
		if resp.Target.Arn == nil || *resp.Target.Arn != tc.lambdaARN() {
			return fmt.Errorf("Target.Arn mismatch: got %q", aws.ToString(resp.Target.Arn))
		}
		if resp.Target.RoleArn == nil || *resp.Target.RoleArn != roleARN {
			return fmt.Errorf("Target.RoleArn mismatch: got %q", aws.ToString(resp.Target.RoleArn))
		}
		if resp.Target.Input == nil || *resp.Target.Input != `{"message":"test message"}` {
			return fmt.Errorf("Target.Input mismatch: got %q", aws.ToString(resp.Target.Input))
		}
		if resp.CreationDate == nil {
			return fmt.Errorf("CreationDate is nil")
		}
		if resp.LastModificationDate == nil {
			return fmt.Errorf("LastModificationDate is nil")
		}
		if resp.FlexibleTimeWindow == nil || resp.FlexibleTimeWindow.Mode != types.FlexibleTimeWindowModeOff {
			return fmt.Errorf("FlexibleTimeWindow.Mode mismatch: got %q", resp.FlexibleTimeWindow.Mode)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule", func() error {
		resp, err := tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(scheduleName),
			ScheduleExpression: aws.String("rate(60 minutes)"),
			Target:             target,
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return err
		}
		if resp.ScheduleArn == nil || *resp.ScheduleArn == "" {
			return fmt.Errorf("ScheduleArn is nil or empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "GetSchedule_AfterUpdate", func() error {
		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != scheduleName {
			return fmt.Errorf("Name mismatch: expected %q, got %q", scheduleName, aws.ToString(resp.Name))
		}
		if resp.ScheduleExpression == nil || *resp.ScheduleExpression != "rate(60 minutes)" {
			return fmt.Errorf("ScheduleExpression not updated: got %q", aws.ToString(resp.ScheduleExpression))
		}
		if resp.Arn == nil {
			return fmt.Errorf("Arn is nil")
		}
		if resp.CreationDate == nil {
			return fmt.Errorf("CreationDate is nil")
		}
		if resp.LastModificationDate == nil {
			return fmt.Errorf("LastModificationDate is nil")
		}
		if resp.Target == nil {
			return fmt.Errorf("Target is nil")
		}
		if resp.FlexibleTimeWindow == nil || resp.FlexibleTimeWindow.Mode != types.FlexibleTimeWindowModeOff {
			return fmt.Errorf("FlexibleTimeWindow mode mismatch")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "DeleteSchedule_VerifyGone", func() error {
		_, err := tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err != nil {
			return fmt.Errorf("delete failed: %v", err)
		}
		_, err = tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("schedule not deleted: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_WithDescription", func() error {
		name := tc.uniqueName("DescSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		desc := "Test schedule description"

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Description:        aws.String(desc),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Description == nil || *resp.Description != desc {
			return fmt.Errorf("Description mismatch: expected %q, got %q", desc, aws.ToString(resp.Description))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_WithTimezone", func() error {
		name := tc.uniqueName("TzSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		tz := "America/New_York"

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:                       aws.String(name),
			ScheduleExpression:         aws.String("rate(30 minutes)"),
			ScheduleExpressionTimezone: aws.String(tz),
			Target:                     tc.defaultTarget(rARN),
			FlexibleTimeWindow:         &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.ScheduleExpressionTimezone == nil || *resp.ScheduleExpressionTimezone != tz {
			return fmt.Errorf("Timezone mismatch: expected %q, got %q", tz, aws.ToString(resp.ScheduleExpressionTimezone))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_WithFlexibleWindow", func() error {
		name := tc.uniqueName("FlexSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode:                   types.FlexibleTimeWindowModeFlexible,
				MaximumWindowInMinutes: aws.Int32(10),
			},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.FlexibleTimeWindow == nil {
			return fmt.Errorf("FlexibleTimeWindow is nil")
		}
		if resp.FlexibleTimeWindow.Mode != types.FlexibleTimeWindowModeFlexible {
			return fmt.Errorf("Mode mismatch: expected FLEXIBLE, got %q", resp.FlexibleTimeWindow.Mode)
		}
		if resp.FlexibleTimeWindow.MaximumWindowInMinutes == nil || *resp.FlexibleTimeWindow.MaximumWindowInMinutes != 10 {
			return fmt.Errorf("MaximumWindowInMinutes mismatch: expected 10, got %d", aws.ToInt32(resp.FlexibleTimeWindow.MaximumWindowInMinutes))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_WithActionAfterCompletion", func() error {
		name := tc.uniqueName("AacSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:                  aws.String(name),
			ScheduleExpression:    aws.String("rate(30 minutes)"),
			ActionAfterCompletion: types.ActionAfterCompletionDelete,
			Target:                tc.defaultTarget(rARN),
			FlexibleTimeWindow:    &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.ActionAfterCompletion != types.ActionAfterCompletionDelete {
			return fmt.Errorf("ActionAfterCompletion mismatch: expected DELETE, got %q", resp.ActionAfterCompletion)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "CreateSchedule_CronExpression", func() error {
		name := tc.uniqueName("CronSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)
		cronExpr := "cron(0 12 * * ? *)"

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String(cronExpr),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		resp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.ScheduleExpression == nil || *resp.ScheduleExpression != cronExpr {
			return fmt.Errorf("ScheduleExpression mismatch: expected %q, got %q", cronExpr, aws.ToString(resp.ScheduleExpression))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule_StateToggle", func() error {
		name := tc.uniqueName("StateSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		_, err = tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			State:              types.ScheduleStateDisabled,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		getResp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.State != types.ScheduleStateDisabled {
			return fmt.Errorf("expected DISABLED, got %q", getResp.State)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule_DescriptionUpdate", func() error {
		name := tc.uniqueName("UpdDescSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Description:        aws.String("original description"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		newDesc := "updated description"
		_, err = tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Description:        aws.String(newDesc),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		getResp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.Description == nil || *getResp.Description != newDesc {
			return fmt.Errorf("Description not updated: expected %q, got %q", newDesc, aws.ToString(getResp.Description))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("scheduler", "UpdateSchedule_VerifyExpression", func() error {
		name := tc.uniqueName("ExprSched")
		rn, rARN := tc.createIAMRole()
		defer tc.deleteIAMRole(rn)

		_, err := tc.client.CreateSchedule(tc.ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.client.DeleteSchedule(tc.ctx, &scheduler.DeleteScheduleInput{Name: aws.String(name)})

		newExpr := "rate(60 minutes)"
		_, err = tc.client.UpdateSchedule(tc.ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(name),
			ScheduleExpression: aws.String(newExpr),
			Target:             tc.defaultTarget(rARN),
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		getResp, err := tc.client.GetSchedule(tc.ctx, &scheduler.GetScheduleInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.ScheduleExpression == nil || *getResp.ScheduleExpression != newExpr {
			return fmt.Errorf("expression not updated: got %q", aws.ToString(getResp.ScheduleExpression))
		}
		return nil
	}))

	return results
}
