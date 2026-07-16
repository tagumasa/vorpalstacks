package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneEngineTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "DescribeDBEngineVersions", func() error {
		resp, err := tc.client.DescribeDBEngineVersions(tc.ctx, &neptune.DescribeDBEngineVersionsInput{
			Engine: aws.String("neptune"),
		})
		if err != nil {
			return err
		}
		if len(resp.DBEngineVersions) == 0 {
			return fmt.Errorf("expected at least one engine version")
		}
		ev := resp.DBEngineVersions[0]
		if ev.Engine == nil || *ev.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune, got %v", ev.Engine)
		}
		if ev.EngineVersion == nil || *ev.EngineVersion == "" {
			return fmt.Errorf("expected non-empty EngineVersion")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBEngineVersions_DefaultEngine", func() error {
		resp, err := tc.client.DescribeDBEngineVersions(tc.ctx, &neptune.DescribeDBEngineVersionsInput{})
		if err != nil {
			return err
		}
		if len(resp.DBEngineVersions) == 0 {
			return fmt.Errorf("expected at least one engine version with default engine")
		}
		foundNeptune := false
		for _, ev := range resp.DBEngineVersions {
			if ev.Engine != nil && *ev.Engine == "neptune" {
				foundNeptune = true
				break
			}
		}
		if !foundNeptune {
			return fmt.Errorf("expected neptune engine in default list")
		}
		return nil
	}))

	return results
}
