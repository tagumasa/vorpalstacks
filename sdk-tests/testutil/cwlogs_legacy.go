package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) legacyTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "ListLogGroups_ReturnsCreated", func() error {
		llgName := tc.uniquePrefix("LLGGroup")
		if err := tc.createLogGroup(llgName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(llgName)

		resp, err := tc.client.ListLogGroups(tc.ctx, &cloudwatchlogs.ListLogGroupsInput{
			LogGroupNamePattern: aws.String("^" + llgName + "$"),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(resp.LogGroups) < 1 {
			return fmt.Errorf("expected at least 1 log group, got %d", len(resp.LogGroups))
		}
		found := false
		for _, lg := range resp.LogGroups {
			if lg.LogGroupName != nil && *lg.LogGroupName == llgName {
				found = true
				if lg.LogGroupArn == nil || *lg.LogGroupArn == "" {
					return fmt.Errorf("logGroupArn is nil or empty")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created log group %q not found in ListLogGroups response", llgName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "TagLogGroup_LegacyAPI", func() error {
		tlgName := tc.uniquePrefix("TLGLegGroup")
		if err := tc.createLogGroup(tlgName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(tlgName)

		_, err := tc.client.TagLogGroup(tc.ctx, &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(tlgName),
			Tags: map[string]string{
				"LegacyKey": "LegacyValue",
			},
		})
		if err != nil {
			return fmt.Errorf("tag log group: %v", err)
		}

		ltResp, err := tc.client.ListTagsLogGroup(tc.ctx, &cloudwatchlogs.ListTagsLogGroupInput{
			LogGroupName: aws.String(tlgName),
		})
		if err != nil {
			return fmt.Errorf("list tags log group: %v", err)
		}
		if ltResp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if ltResp.Tags["LegacyKey"] != "LegacyValue" {
			return fmt.Errorf("LegacyKey mismatch: got %q, want %q", ltResp.Tags["LegacyKey"], "LegacyValue")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "UntagLogGroup_LegacyAPI", func() error {
		utlgName := tc.uniquePrefix("UTLGLegGroup")
		if err := tc.createLogGroup(utlgName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(utlgName)

		tc.client.TagLogGroup(tc.ctx, &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(utlgName),
			Tags: map[string]string{
				"Remove": "yes",
				"Keep":   "yes",
			},
		})

		_, err := tc.client.UntagLogGroup(tc.ctx, &cloudwatchlogs.UntagLogGroupInput{
			LogGroupName: aws.String(utlgName),
			Tags:         []string{"Remove"},
		})
		if err != nil {
			return fmt.Errorf("untag log group: %v", err)
		}

		ltResp, err := tc.client.ListTagsLogGroup(tc.ctx, &cloudwatchlogs.ListTagsLogGroupInput{
			LogGroupName: aws.String(utlgName),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if _, ok := ltResp.Tags["Remove"]; ok {
			return fmt.Errorf("Remove tag should have been removed")
		}
		if ltResp.Tags["Keep"] != "yes" {
			return fmt.Errorf("Keep tag should still exist")
		}
		return nil
	}))

	return results
}
