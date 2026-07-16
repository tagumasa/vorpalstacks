package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

func (r *TestRunner) runSFNVersionTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	verSMName := fmt.Sprintf("VerSM-%d", time.Now().UnixNano())
	_, verRoleARN, verRoleCleanup := tc.createRoleForSM("VerRole")
	defer verRoleCleanup()

	def2 := `{"Comment":"v2","StartAt":"B","States":{"B":{"Type":"Pass","Result":"hello","End":true}}}`
	var verSMARN string
	resp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(verSMName),
		Definition: aws.String(def2),
		RoleArn:    aws.String(verRoleARN),
	})
	if err != nil {
		return append(results, TestResult{Service: "stepfunctions", TestName: "VersionSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)})
	}
	verSMARN = *resp.StateMachineArn
	defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(verSMARN)})

	var versionARN string
	results = append(results, r.RunTest("stepfunctions", "PublishStateMachineVersion", func() error {
		resp, err := tc.client.PublishStateMachineVersion(tc.ctx, &sfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(verSMARN),
			Description:     aws.String("test version 1"),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineVersionArn == nil || *resp.StateMachineVersionArn == "" {
			return fmt.Errorf("version ARN is nil or empty")
		}
		versionARN = *resp.StateMachineVersionArn
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "ListStateMachineVersions", func() error {
		resp, err := tc.client.ListStateMachineVersions(tc.ctx, &sfn.ListStateMachineVersionsInput{
			StateMachineArn: aws.String(verSMARN),
		})
		if err != nil {
			return err
		}
		if len(resp.StateMachineVersions) == 0 {
			return fmt.Errorf("expected at least one version")
		}
		found := false
		for _, v := range resp.StateMachineVersions {
			if v.StateMachineVersionArn != nil && *v.StateMachineVersionArn == versionARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("published version not found in list")
		}
		return nil
	}))

	var secondVersionARN string
	results = append(results, r.RunTest("stepfunctions", "PublishStateMachineVersion_Second", func() error {
		resp, err := tc.client.PublishStateMachineVersion(tc.ctx, &sfn.PublishStateMachineVersionInput{
			StateMachineArn: aws.String(verSMARN),
			Description:     aws.String("test version 2"),
		})
		if err != nil {
			return err
		}
		if resp.StateMachineVersionArn == nil || *resp.StateMachineVersionArn == "" {
			return fmt.Errorf("version ARN is nil or empty")
		}
		if *resp.StateMachineVersionArn == versionARN {
			return fmt.Errorf("second version should have different ARN")
		}
		secondVersionARN = *resp.StateMachineVersionArn
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "DeleteStateMachineVersion", func() error {
		if versionARN == "" {
			return fmt.Errorf("version ARN not available")
		}
		_, err := tc.client.DeleteStateMachineVersion(tc.ctx, &sfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String(versionARN),
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListStateMachineVersions(tc.ctx, &sfn.ListStateMachineVersionsInput{
			StateMachineArn: aws.String(verSMARN),
		})
		if err != nil {
			return err
		}
		for _, v := range listResp.StateMachineVersions {
			if v.StateMachineVersionArn != nil && *v.StateMachineVersionArn == versionARN {
				return fmt.Errorf("deleted version still appears in list")
			}
		}
		return nil
	}))

	if secondVersionARN != "" {
		tc.client.DeleteStateMachineVersion(tc.ctx, &sfn.DeleteStateMachineVersionInput{
			StateMachineVersionArn: aws.String(secondVersionARN),
		})
	}

	return results
}
