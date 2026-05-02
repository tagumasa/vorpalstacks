package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (tc *cloudwatchTestCtx) tagTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("cloudwatch", "TagResource_VerifyTags", func() error {
		alarmName := tc.uniquePrefix("TagAlarm")
		testNS := tc.uniquePrefix("TagNS")
		if err := tc.createAlarm(alarmName, testNS, "TestMetric", 50.0); err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.MetricAlarms) != 1 || descResp.MetricAlarms[0].AlarmArn == nil {
			return fmt.Errorf("failed to get alarm ARN")
		}
		alarmARN := descResp.MetricAlarms[0].AlarmArn

		_, err = tc.client.TagResource(tc.ctx, &cloudwatch.TagResourceInput{
			ResourceARN: alarmARN,
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Team"), Value: aws.String("platform")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &cloudwatch.ListTagsForResourceInput{
			ResourceARN: alarmARN,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tagResp.Tags) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagResp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range tagResp.Tags {
			if t.Key != nil && t.Value != nil {
				tagMap[*t.Key] = *t.Value
			}
		}
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("Environment tag mismatch: got %q", tagMap["Environment"])
		}
		if tagMap["Team"] != "platform" {
			return fmt.Errorf("Team tag mismatch: got %q", tagMap["Team"])
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "UntagResource_VerifyRemoval", func() error {
		alarmName := tc.uniquePrefix("UntagAlarm")
		testNS := tc.uniquePrefix("UntagNS")
		if err := tc.createAlarm(alarmName, testNS, "TestMetric", 50.0); err != nil {
			return fmt.Errorf("put alarm: %v", err)
		}
		defer tc.deleteAlarms(alarmName)

		descResp, err := tc.client.DescribeAlarms(tc.ctx, &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []string{alarmName},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.MetricAlarms) != 1 || descResp.MetricAlarms[0].AlarmArn == nil {
			return fmt.Errorf("failed to get alarm ARN")
		}
		alarmARN := descResp.MetricAlarms[0].AlarmArn

		_, err = tc.client.TagResource(tc.ctx, &cloudwatch.TagResourceInput{
			ResourceARN: alarmARN,
			Tags: []types.Tag{
				{Key: aws.String("Keep"), Value: aws.String("yes")},
				{Key: aws.String("Remove"), Value: aws.String("yes")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}

		_, err = tc.client.UntagResource(tc.ctx, &cloudwatch.UntagResourceInput{
			ResourceARN: alarmARN,
			TagKeys:     []string{"Remove"},
		})
		if err != nil {
			return fmt.Errorf("untag resource: %v", err)
		}

		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &cloudwatch.ListTagsForResourceInput{
			ResourceARN: alarmARN,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(tagResp.Tags) != 1 {
			return fmt.Errorf("expected 1 tag after untag, got %d", len(tagResp.Tags))
		}
		if tagResp.Tags[0].Key == nil || *tagResp.Tags[0].Key != "Keep" {
			return fmt.Errorf("expected Keep tag, got %v", tagResp.Tags[0].Key)
		}
		if tagResp.Tags[0].Value == nil || *tagResp.Tags[0].Value != "yes" {
			return fmt.Errorf("expected Keep=yes, got %v", tagResp.Tags[0].Value)
		}
		return nil
	}))

	return results
}
