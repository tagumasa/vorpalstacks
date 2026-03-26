package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSchedulerTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "scheduler",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := scheduler.NewFromConfig(cfg)
	ctx := context.Background()

	scheduleName := fmt.Sprintf("TestSchedule-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/TestRole")

	targetInput := map[string]string{
		"message": "test message",
	}
	targetInputJSON, _ := json.Marshal(targetInput)

	results = append(results, r.RunTest("scheduler", "CreateSchedule", func() error {
		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(scheduleName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFunction", r.region)),
				RoleArn: aws.String(roleARN),
				Input:   aws.String(string(targetInputJSON)),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("scheduler", "GetSchedule", func() error {
		resp, err := client.GetSchedule(ctx, &scheduler.GetScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("schedule name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "ListSchedules", func() error {
		resp, err := client.ListSchedules(ctx, &scheduler.ListSchedulesInput{})
		if err != nil {
			return err
		}
		if resp.Schedules == nil {
			return fmt.Errorf("schedules list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "UpdateSchedule", func() error {
		_, err := client.UpdateSchedule(ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(scheduleName),
			ScheduleExpression: aws.String("rate(60 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFunction", r.region)),
				RoleArn: aws.String(roleARN),
				Input:   aws.String(string(targetInputJSON)),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("scheduler", "TagResource", func() error {
		scheduleARN := fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/%s", r.region, scheduleName)
		_, err := client.TagResource(ctx, &scheduler.TagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("scheduler", "ListTagsForResource", func() error {
		scheduleARN := fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/%s", r.region, scheduleName)
		resp, err := client.ListTagsForResource(ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: aws.String(scheduleARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "UntagResource", func() error {
		scheduleARN := fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/%s", r.region, scheduleName)
		_, err := client.UntagResource(ctx, &scheduler.UntagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			TagKeys:     []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("scheduler", "DeleteSchedule", func() error {
		_, err := client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("scheduler", "GetSchedule_NonExistent", func() error {
		_, err := client.GetSchedule(ctx, &scheduler.GetScheduleInput{
			Name: aws.String("nonexistent-schedule-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent schedule")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "DeleteSchedule_NonExistent", func() error {
		_, err := client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String("nonexistent-schedule-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent schedule")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "CreateSchedule_DuplicateName", func() error {
		dupName := fmt.Sprintf("DupSchedule-%d", time.Now().UnixNano())
		targetInput := map[string]string{"msg": "test"}
		targetInputJSON, _ := json.Marshal(targetInput)
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/TestRole")

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(roleARN),
				Input:   aws.String(string(targetInputJSON)),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{Name: aws.String(dupName)})

		_, err = client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(60 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(roleARN),
				Input:   aws.String(string(targetInputJSON)),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate schedule name")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "UpdateSchedule_VerifyExpression", func() error {
		updName := fmt.Sprintf("UpdSchedule-%d", time.Now().UnixNano())
		targetInput := map[string]string{"msg": "test"}
		targetInputJSON, _ := json.Marshal(targetInput)
		roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/TestRole")

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(updName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(roleARN),
				Input:   aws.String(string(targetInputJSON)),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{Name: aws.String(updName)})

		newExpr := "rate(60 minutes)"
		_, err = client.UpdateSchedule(ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(updName),
			ScheduleExpression: aws.String(newExpr),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(roleARN),
				Input:   aws.String(string(targetInputJSON)),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		getResp, err := client.GetSchedule(ctx, &scheduler.GetScheduleInput{
			Name: aws.String(updName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.ScheduleExpression == nil || *getResp.ScheduleExpression != newExpr {
			return fmt.Errorf("schedule expression not updated, got %q", aws.ToString(getResp.ScheduleExpression))
		}
		return nil
	}))

	return results
}
