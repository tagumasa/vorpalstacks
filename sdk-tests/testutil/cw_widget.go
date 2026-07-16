package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func (tc *cloudwatchTestCtx) widgetTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("cloudwatch", "GetMetricWidgetImage_Basic", func() error {
		widget := `{"metrics":[["AWS/EC2","CPUUtilization"]]}`
		resp, err := tc.client.GetMetricWidgetImage(tc.ctx, &cloudwatch.GetMetricWidgetImageInput{
			MetricWidget: aws.String(widget),
		})
		if err != nil {
			return fmt.Errorf("get metric widget image: %v", err)
		}
		if len(resp.MetricWidgetImage) == 0 {
			return fmt.Errorf("expected non-empty image")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "GetMetricWidgetImage_InvalidWidget", func() error {
		_, err := tc.client.GetMetricWidgetImage(tc.ctx, &cloudwatch.GetMetricWidgetImageInput{
			MetricWidget: aws.String("invalid-json"),
		})
		if err := AssertErrorContains(err, "InvalidParameter"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
