package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) retentionTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "DeleteRetentionPolicy_Basic", func() error {
		drName := tc.uniquePrefix("DelRetGroup")
		if err := tc.createLogGroup(drName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(drName)

		_, err := tc.client.PutRetentionPolicy(tc.ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(drName),
			RetentionInDays: aws.Int32(14),
		})
		if err != nil {
			return fmt.Errorf("put retention: %v", err)
		}

		_, err = tc.client.DeleteRetentionPolicy(tc.ctx, &cloudwatchlogs.DeleteRetentionPolicyInput{
			LogGroupName: aws.String(drName),
		})
		if err != nil {
			return fmt.Errorf("delete retention: %v", err)
		}

		descResp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(drName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.LogGroups) == 0 {
			return fmt.Errorf("log group not found")
		}
		if descResp.LogGroups[0].RetentionInDays != nil && *descResp.LogGroups[0].RetentionInDays != 0 {
			return fmt.Errorf("expected retention 0 after delete, got %d", *descResp.LogGroups[0].RetentionInDays)
		}
		return nil
	}))

	return results
}
