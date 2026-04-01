package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	roleName := fmt.Sprintf("TestSchedRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)

	trustPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"scheduler.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	form := url.Values{}
	form.Set("Action", "CreateRole")
	form.Set("Version", "2010-05-08")
	form.Set("RoleName", roleName)
	form.Set("AssumeRolePolicyDocument", trustPolicy)
	req, err := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		results = append(results, TestResult{Service: "scheduler", TestName: "Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role request: %v", err)})
		return results
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		results = append(results, TestResult{Service: "scheduler", TestName: "Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
		return results
	}
	resp.Body.Close()
	defer func() {
		cleanupForm := url.Values{}
		cleanupForm.Set("Action", "DeleteRole")
		cleanupForm.Set("Version", "2010-05-08")
		cleanupForm.Set("RoleName", roleName)
		cleanupReq, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(cleanupForm.Encode()))
		if cleanupReq != nil {
			cleanupReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			cleanupResp, _ := http.DefaultClient.Do(cleanupReq)
			if cleanupResp != nil {
				cleanupResp.Body.Close()
			}
		}
	}()

	createRole := func(rn string) {
		f := url.Values{}
		f.Set("Action", "CreateRole")
		f.Set("Version", "2010-05-08")
		f.Set("RoleName", rn)
		f.Set("AssumeRolePolicyDocument", trustPolicy)
		r, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(f.Encode()))
		if r != nil {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if rp, err := http.DefaultClient.Do(r); err == nil {
				rp.Body.Close()
			}
		}
	}
	deleteRole := func(rn string) {
		f := url.Values{}
		f.Set("Action", "DeleteRole")
		f.Set("Version", "2010-05-08")
		f.Set("RoleName", rn)
		r, _ := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewBufferString(f.Encode()))
		if r != nil {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if rp, _ := http.DefaultClient.Do(r); rp != nil {
				rp.Body.Close()
			}
		}
	}

	targetInput := map[string]string{
		"message": "test message",
	}
	targetInputJSON, _ := json.Marshal(targetInput)

	results = append(results, r.RunTest("scheduler", "CreateSchedule", func() error {
		resp, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
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
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.UpdateSchedule(ctx, &scheduler.UpdateScheduleInput{
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
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "TagResource", func() error {
		scheduleARN := fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/%s", r.region, scheduleName)
		resp, err := client.TagResource(ctx, &scheduler.TagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		resp, err := client.UntagResource(ctx, &scheduler.UntagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "GetSchedule_ContentVerify", func() error {
		getResp, err := client.GetSchedule(ctx, &scheduler.GetScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err != nil {
			return err
		}
		if getResp.Name == nil || *getResp.Name != scheduleName {
			return fmt.Errorf("name mismatch")
		}
		if getResp.ScheduleExpression == nil || *getResp.ScheduleExpression != "rate(60 minutes)" {
			return fmt.Errorf("expression mismatch: %q", aws.ToString(getResp.ScheduleExpression))
		}
		if getResp.Arn == nil {
			return fmt.Errorf("ARN is nil")
		}
		if getResp.CreationDate == nil {
			return fmt.Errorf("CreationDate is nil")
		}
		if getResp.LastModificationDate == nil {
			return fmt.Errorf("LastModificationDate is nil")
		}
		if getResp.Target == nil {
			return fmt.Errorf("Target is nil")
		}
		if getResp.FlexibleTimeWindow == nil || getResp.FlexibleTimeWindow.Mode != types.FlexibleTimeWindowModeOff {
			return fmt.Errorf("FlexibleTimeWindow mode mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "DeleteSchedule", func() error {
		resp, err := client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{
			Name: aws.String(scheduleName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		dupRoleName := fmt.Sprintf("DupSchedRole-%d", time.Now().UnixNano())
		createRole(dupRoleName)
		defer deleteRole(dupRoleName)
		targetInput := map[string]string{"msg": "test"}
		targetInputJSON, _ := json.Marshal(targetInput)
		dupRoleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", dupRoleName)

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(dupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(dupRoleARN),
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
				RoleArn: aws.String(dupRoleARN),
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
		updRoleName := fmt.Sprintf("UpdSchedRole-%d", time.Now().UnixNano())
		createRole(updRoleName)
		defer deleteRole(updRoleName)
		targetInput := map[string]string{"msg": "test"}
		targetInputJSON, _ := json.Marshal(targetInput)
		updRoleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", updRoleName)

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(updName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(updRoleARN),
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
				RoleArn: aws.String(updRoleARN),
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

	// === SCHEDULE GROUP TESTS ===

	groupName := fmt.Sprintf("TestGroup-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("scheduler", "CreateScheduleGroup", func() error {
		resp, err := client.CreateScheduleGroup(ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp == nil || resp.ScheduleGroupArn == nil {
			return fmt.Errorf("response or ARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "GetScheduleGroup", func() error {
		resp, err := client.GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != groupName {
			return fmt.Errorf("expected group name %q, got %q", groupName, aws.ToString(resp.Name))
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "ListScheduleGroups", func() error {
		resp, err := client.ListScheduleGroups(ctx, &scheduler.ListScheduleGroupsInput{})
		if err != nil {
			return err
		}
		if resp.ScheduleGroups == nil {
			return fmt.Errorf("schedule groups list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "GetScheduleGroup_NonExistent", func() error {
		_, err := client.GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String("nonexistent-group-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent group")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "CreateScheduleGroup_DuplicateName", func() error {
		dupGroupName := fmt.Sprintf("DupGroup-%d", time.Now().UnixNano())
		_, err := client.CreateScheduleGroup(ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(dupGroupName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteScheduleGroup(ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(dupGroupName)})

		_, err = client.CreateScheduleGroup(ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(dupGroupName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate group name")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "ListScheduleGroups_ContainsCreated", func() error {
		resp, err := client.ListScheduleGroups(ctx, &scheduler.ListScheduleGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, g := range resp.ScheduleGroups {
			if g.Name != nil && *g.Name == groupName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created group %q not found in list", groupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "CreateSchedule_WithGroupName", func() error {
		schedName := fmt.Sprintf("GroupSched-%d", time.Now().UnixNano())
		groupRoleName := fmt.Sprintf("GroupSchedRole-%d", time.Now().UnixNano())
		createRole(groupRoleName)
		defer deleteRole(groupRoleName)
		groupRoleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", groupRoleName)

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(schedName),
			GroupName:          aws.String(groupName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(groupRoleARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return err
		}
		defer client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{Name: aws.String(schedName), GroupName: aws.String(groupName)})

		getResp, err := client.GetSchedule(ctx, &scheduler.GetScheduleInput{
			Name:      aws.String(schedName),
			GroupName: aws.String(groupName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.GroupName == nil || *getResp.GroupName != groupName {
			return fmt.Errorf("expected group name %q, got %q", groupName, aws.ToString(getResp.GroupName))
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "UpdateSchedule_NonExistent", func() error {
		_, err := client.UpdateSchedule(ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String("nonexistent-schedule-xyz"),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(roleARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent schedule")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "UpdateSchedule_StateToggle", func() error {
		stateName := fmt.Sprintf("StateSchedule-%d", time.Now().UnixNano())
		stateRoleName := fmt.Sprintf("StateSchedRole-%d", time.Now().UnixNano())
		createRole(stateRoleName)
		defer deleteRole(stateRoleName)
		stateRoleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", stateRoleName)

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(stateName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(stateRoleARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{Name: aws.String(stateName)})

		_, err = client.UpdateSchedule(ctx, &scheduler.UpdateScheduleInput{
			Name:               aws.String(stateName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(stateRoleARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
			State: types.ScheduleStateDisabled,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		getResp, err := client.GetSchedule(ctx, &scheduler.GetScheduleInput{
			Name: aws.String(stateName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if getResp.State != types.ScheduleStateDisabled {
			return fmt.Errorf("expected DISABLED, got %q", getResp.State)
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "ListSchedules_NamePrefix", func() error {
		prefixName := fmt.Sprintf("PrefixSched-%d", time.Now().UnixNano())
		prefixRoleName := fmt.Sprintf("PrefixSchedRole-%d", time.Now().UnixNano())
		createRole(prefixRoleName)
		defer deleteRole(prefixRoleName)
		prefixRoleARN := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", prefixRoleName)

		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(prefixName),
			ScheduleExpression: aws.String("rate(30 minutes)"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(prefixRoleARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err != nil {
			return err
		}
		defer client.DeleteSchedule(ctx, &scheduler.DeleteScheduleInput{Name: aws.String(prefixName)})

		prefix := prefixName[:len(prefixName)-8]
		resp, err := client.ListSchedules(ctx, &scheduler.ListSchedulesInput{
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

	results = append(results, r.RunTest("scheduler", "TagResource_ScheduleGroup", func() error {
		groupResp, err := client.CreateScheduleGroup(ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(fmt.Sprintf("TagGroup-%d", time.Now().UnixNano())),
		})
		if err != nil {
			return err
		}
		defer client.DeleteScheduleGroup(ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(fmt.Sprintf("TagGroup-%d", time.Now().UnixNano()))})

		_, err = client.TagResource(ctx, &scheduler.TagResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("prod")},
			},
		})
		if err != nil {
			return err
		}

		tagResp, err := client.ListTagsForResource(ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: groupResp.ScheduleGroupArn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		found := false
		for _, t := range tagResp.Tags {
			if t.Key != nil && *t.Key == "Env" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tag Env not found on schedule group")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "ListTagsForResource_ScheduleGroup", func() error {
		groupResp, err := client.GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String(groupName),
		})
		if err != nil {
			return fmt.Errorf("get group: %v", err)
		}

		_, err = client.ListTagsForResource(ctx, &scheduler.ListTagsForResourceInput{
			ResourceArn: groupResp.Arn,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "UntagResource_NonExistentKey", func() error {
		scheduleARN := fmt.Sprintf("arn:aws:scheduler:%s:000000000000:schedule/%s", r.region, scheduleName)
		_, err := client.UntagResource(ctx, &scheduler.UntagResourceInput{
			ResourceArn: aws.String(scheduleARN),
			TagKeys:     []string{"NonExistentKey"},
		})
		if err != nil {
			return fmt.Errorf("untag non-existent key should not error: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "CreateSchedule_InvalidExpression", func() error {
		invName := fmt.Sprintf("InvExprSched-%d", time.Now().UnixNano())
		_, err := client.CreateSchedule(ctx, &scheduler.CreateScheduleInput{
			Name:               aws.String(invName),
			ScheduleExpression: aws.String("not-a-valid-expression"),
			Target: &types.Target{
				Arn:     aws.String(fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:TestFn", r.region)),
				RoleArn: aws.String(roleARN),
			},
			FlexibleTimeWindow: &types.FlexibleTimeWindow{
				Mode: types.FlexibleTimeWindowModeOff,
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid schedule expression")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "DeleteScheduleGroup_NonExistent", func() error {
		_, err := client.DeleteScheduleGroup(ctx, &scheduler.DeleteScheduleGroupInput{
			Name: aws.String("nonexistent-group-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent group")
		}
		return nil
	}))

	results = append(results, r.RunTest("scheduler", "DeleteScheduleGroup", func() error {
		delGroupName := fmt.Sprintf("DelGroup-%d", time.Now().UnixNano())
		_, err := client.CreateScheduleGroup(ctx, &scheduler.CreateScheduleGroupInput{
			Name: aws.String(delGroupName),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		_, err = client.DeleteScheduleGroup(ctx, &scheduler.DeleteScheduleGroupInput{
			Name: aws.String(delGroupName),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		_, err = client.GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String(delGroupName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting group")
		}
		return nil
	}))

	client.DeleteScheduleGroup(ctx, &scheduler.DeleteScheduleGroupInput{Name: aws.String(groupName)})

	return results
}
