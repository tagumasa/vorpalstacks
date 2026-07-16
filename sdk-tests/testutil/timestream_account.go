package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	tqtypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
)

func (r *TestRunner) runTimestreamAccountTests(tc *tsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("timestream", "DescribeAccountSettings_Verify", func() error {
		resp, err := tc.queryClient.DescribeAccountSettings(tc.ctx, &timestreamquery.DescribeAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.MaxQueryTCU == nil {
			return fmt.Errorf("expected MaxQueryTCU to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UpdateAccountSettings_Verify", func() error {
		_, err := tc.queryClient.UpdateAccountSettings(tc.ctx, &timestreamquery.UpdateAccountSettingsInput{
			MaxQueryTCU:       aws.Int32(8),
			QueryPricingModel: tqtypes.QueryPricingModelComputeUnits,
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeAccountSettings_AfterUpdate", func() error {
		resp, err := tc.queryClient.DescribeAccountSettings(tc.ctx, &timestreamquery.DescribeAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.MaxQueryTCU == nil || *resp.MaxQueryTCU != 8 {
			return fmt.Errorf("expected MaxQueryTCU=8, got %v", resp.MaxQueryTCU)
		}
		return nil
	}))

	return results
}
