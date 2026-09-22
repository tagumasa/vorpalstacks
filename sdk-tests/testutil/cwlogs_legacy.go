package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) legacyTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "ListLogGroups_ReturnsCreated", func() error {
		// The pattern member's shape caps each pipe-separated
		// alternative at 3-24 restricted characters with an optional
		// leading ^ anchor (no trailing $ in the vocabulary), so the
		// fixture keeps the group name short enough for one anchored
		// alternative to cover it.
		llgName := fmt.Sprintf("lg-%d", time.Now().UnixNano()%1e12)
		if err := tc.createLogGroup(llgName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(llgName)

		resp, err := tc.client.ListLogGroups(tc.ctx, &cloudwatchlogs.ListLogGroupsInput{
			LogGroupNamePattern: aws.String("^" + llgName),
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

		// A bare (unanchored) pattern matches by substring, never by
		// prefix: the fixture below carries the pattern in suffix
		// position behind a leading character, so a prefix matcher would
		// miss it. The pattern stays inside the member's 3-24 character
		// per-alternative bound.
		substr := llgName[3:] + "DataLogs" // up to 12 digits + 8 characters
		suffixGroup := "x" + substr
		if err := tc.createLogGroup(suffixGroup); err != nil {
			return fmt.Errorf("create suffix group: %v", err)
		}
		defer tc.deleteLogGroup(suffixGroup)
		resp, err = tc.client.ListLogGroups(tc.ctx, &cloudwatchlogs.ListLogGroupsInput{
			LogGroupNamePattern: aws.String(substr),
		})
		if err != nil {
			return fmt.Errorf("list with non-prefix substring pattern: %v", err)
		}
		found = false
		for _, lg := range resp.LogGroups {
			if lg.LogGroupName != nil && *lg.LogGroupName == suffixGroup {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("substring in non-prefix position: group %q not found via bare pattern", suffixGroup)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "TagLogGroup_LegacyAPI", func() error {
		tlgName, cleanupGroup, err := tc.newLogGroupFixture("TLGLegGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		_, err = tc.client.TagLogGroup(tc.ctx, &cloudwatchlogs.TagLogGroupInput{
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
		utlgName, cleanupGroup, err := tc.newLogGroupFixture("UTLGLegGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		if _, err := tc.client.TagLogGroup(tc.ctx, &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(utlgName),
			Tags: map[string]string{
				"Remove": "yes",
				"Keep":   "yes",
			},
		}); err != nil {
			return fmt.Errorf("tag log group: %v", err)
		}

		_, err = tc.client.UntagLogGroup(tc.ctx, &cloudwatchlogs.UntagLogGroupInput{
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
