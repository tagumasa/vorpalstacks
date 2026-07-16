package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
)

func (r *TestRunner) runNeptunedataStatisticsTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphStatistics", func() error {
		resp, err := tc.client.GetPropertygraphStatistics(tc.ctx, &neptunedata.GetPropertygraphStatisticsInput{})
		if err != nil {
			return err
		}
		if resp.Status == nil {
			return fmt.Errorf("expected non-nil status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphSummary", func() error {
		resp, err := tc.client.GetPropertygraphSummary(tc.ctx, &neptunedata.GetPropertygraphSummaryInput{
			Mode: types.GraphSummaryTypeBasic,
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil {
			return fmt.Errorf("expected non-nil statusCode")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ManagePropertygraphStatistics_Disable", func() error {
		resp, err := tc.client.ManagePropertygraphStatistics(tc.ctx, &neptunedata.ManagePropertygraphStatisticsInput{
			Mode: types.StatisticsAutoGenerationModeDisableAutocompute,
		})
		if err != nil {
			return err
		}
		if resp.Status == nil || *resp.Status == "" {
			return fmt.Errorf("expected non-empty status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ManagePropertygraphStatistics_Enable", func() error {
		resp, err := tc.client.ManagePropertygraphStatistics(tc.ctx, &neptunedata.ManagePropertygraphStatisticsInput{
			Mode: types.StatisticsAutoGenerationModeEnableAutocompute,
		})
		if err != nil {
			return err
		}
		if resp.Status == nil || *resp.Status == "" {
			return fmt.Errorf("expected non-empty status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ManagePropertygraphStatistics_Refresh", func() error {
		resp, err := tc.client.ManagePropertygraphStatistics(tc.ctx, &neptunedata.ManagePropertygraphStatisticsInput{
			Mode: types.StatisticsAutoGenerationModeRefresh,
		})
		if err != nil {
			return err
		}
		if resp.Status == nil {
			return fmt.Errorf("expected non-nil status from refresh")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "DeletePropertygraphStatistics", func() error {
		resp, err := tc.client.DeletePropertygraphStatistics(tc.ctx, &neptunedata.DeletePropertygraphStatisticsInput{})
		if err != nil {
			return err
		}
		if resp.Status == nil {
			return fmt.Errorf("expected non-nil status from delete statistics")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphStream", func() error {
		resp, err := tc.client.GetPropertygraphStream(tc.ctx, &neptunedata.GetPropertygraphStreamInput{})
		if err != nil {
			return err
		}
		if resp.Format == nil {
			return fmt.Errorf("expected non-nil format from stream")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphSummary_Detailed", func() error {
		resp, err := tc.client.GetPropertygraphSummary(tc.ctx, &neptunedata.GetPropertygraphSummaryInput{
			Mode: types.GraphSummaryTypeDetailed,
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil {
			return fmt.Errorf("expected non-nil statusCode from detailed summary")
		}
		return nil
	}))

	return results
}
