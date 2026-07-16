package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func (r *TestRunner) runSFNTagTests(tc *sfnTestContext) []TestResult {
	var results []TestResult

	tagSMName := fmt.Sprintf("TagSM-%d", time.Now().UnixNano())
	_, tagRoleARN, tagRoleCleanup := tc.createRoleForSM("TagRole")
	defer tagRoleCleanup()

	tagDef := `{"Comment":"tag test","StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`
	tagResp, err := tc.client.CreateStateMachine(tc.ctx, &sfn.CreateStateMachineInput{
		Name:       aws.String(tagSMName),
		Definition: aws.String(tagDef),
		RoleArn:    aws.String(tagRoleARN),
		Tags: []types.Tag{
			{Key: aws.String("InitialTag"), Value: aws.String("initial")},
		},
	})
	if err != nil {
		return append(results, TestResult{Service: "stepfunctions", TestName: "TagSetup", Status: "FAIL", Error: fmt.Sprintf("create SM: %v", err)})
	}
	tagSMARN := *tagResp.StateMachineArn
	defer tc.client.DeleteStateMachine(tc.ctx, &sfn.DeleteStateMachineInput{StateMachineArn: aws.String(tagSMARN)})

	results = append(results, r.RunTest("stepfunctions", "CreateStateMachine_WithTags", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &sfn.ListTagsForResourceInput{
			ResourceArn: aws.String(tagSMARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) == 0 {
			return fmt.Errorf("expected at least one tag from creation")
		}
		found := false
		for _, t := range resp.Tags {
			if t.Key != nil && *t.Key == "InitialTag" {
				found = true
				if t.Value == nil || *t.Value != "initial" {
					return fmt.Errorf("InitialTag value mismatch: got %v", t.Value)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("InitialTag not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "TagResource", func() error {
		_, err := tc.client.TagResource(tc.ctx, &sfn.TagResourceInput{
			ResourceArn: aws.String(tagSMARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Project"), Value: aws.String("vorpalstacks")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("stepfunctions", "ListTagsForResource", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &sfn.ListTagsForResourceInput{
			ResourceArn: aws.String(tagSMARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 3 {
			return fmt.Errorf("expected at least 3 tags, got %d", len(resp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			if t.Key != nil && t.Value != nil {
				tagMap[*t.Key] = *t.Value
			}
		}
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("Environment tag mismatch: got %q", tagMap["Environment"])
		}
		if tagMap["Project"] != "vorpalstacks" {
			return fmt.Errorf("Project tag mismatch: got %q", tagMap["Project"])
		}
		return nil
	}))

	results = append(results, r.RunTest("stepfunctions", "UntagResource", func() error {
		_, err := tc.client.UntagResource(tc.ctx, &sfn.UntagResourceInput{
			ResourceArn: aws.String(tagSMARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		verifyResp, err := tc.client.ListTagsForResource(tc.ctx, &sfn.ListTagsForResourceInput{
			ResourceArn: aws.String(tagSMARN),
		})
		if err != nil {
			return err
		}
		for _, t := range verifyResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				return fmt.Errorf("Environment tag should have been removed")
			}
		}
		return nil
	}))

	return results
}
