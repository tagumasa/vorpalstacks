package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func (tc *cloudwatchTestCtx) dashboardTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("cloudwatch", "PutDashboard_GetDashboard_Roundtrip", func() error {
		dbName := tc.uniquePrefix("TestDashboard")
		dashboardBody := `{"widgets":[{"type":"metric","x":0,"y":0,"width":12,"height":6,"properties":{"metrics":[["AWS/EC2","CPUUtilization"]],"period":300,"stat":"Average","region":"us-east-1","title":"EC2 CPU"}}]}`
		_, err := tc.client.PutDashboard(tc.ctx, &cloudwatch.PutDashboardInput{
			DashboardName: aws.String(dbName),
			DashboardBody: aws.String(dashboardBody),
		})
		if err != nil {
			return fmt.Errorf("put dashboard: %v", err)
		}
		defer tc.client.DeleteDashboards(tc.ctx, &cloudwatch.DeleteDashboardsInput{
			DashboardNames: []string{dbName},
		})

		getResp, err := tc.client.GetDashboard(tc.ctx, &cloudwatch.GetDashboardInput{
			DashboardName: aws.String(dbName),
		})
		if err != nil {
			return fmt.Errorf("get dashboard: %v", err)
		}
		if getResp.DashboardName == nil || *getResp.DashboardName != dbName {
			return fmt.Errorf("dashboard name mismatch: got %v", getResp.DashboardName)
		}
		if getResp.DashboardBody == nil || *getResp.DashboardBody == "" {
			return fmt.Errorf("dashboard body is empty")
		}
		if getResp.DashboardArn == nil || *getResp.DashboardArn == "" {
			return fmt.Errorf("dashboard ARN is empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "ListDashboards_AfterCreate", func() error {
		db1 := tc.uniquePrefix("ListDB1")
		db2 := tc.uniquePrefix("ListDB2")
		body := `{"widgets":[]}`
		tc.client.PutDashboard(tc.ctx, &cloudwatch.PutDashboardInput{DashboardName: aws.String(db1), DashboardBody: aws.String(body)})
		tc.client.PutDashboard(tc.ctx, &cloudwatch.PutDashboardInput{DashboardName: aws.String(db2), DashboardBody: aws.String(body)})
		defer tc.client.DeleteDashboards(tc.ctx, &cloudwatch.DeleteDashboardsInput{DashboardNames: []string{db1, db2}})

		resp, err := tc.client.ListDashboards(tc.ctx, &cloudwatch.ListDashboardsInput{})
		if err != nil {
			return fmt.Errorf("list dashboards: %v", err)
		}
		if len(resp.DashboardEntries) < 2 {
			return fmt.Errorf("expected >= 2 dashboards, got %d", len(resp.DashboardEntries))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "DeleteDashboards_Verify", func() error {
		dbName := tc.uniquePrefix("DelDB")
		body := `{"widgets":[]}`
		_, err := tc.client.PutDashboard(tc.ctx, &cloudwatch.PutDashboardInput{
			DashboardName: aws.String(dbName),
			DashboardBody: aws.String(body),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = tc.client.DeleteDashboards(tc.ctx, &cloudwatch.DeleteDashboardsInput{
			DashboardNames: []string{dbName},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		_, err = tc.client.GetDashboard(tc.ctx, &cloudwatch.GetDashboardInput{
			DashboardName: aws.String(dbName),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted dashboard")
		}
		return nil
	}))

	return results
}
