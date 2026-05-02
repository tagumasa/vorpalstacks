package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
)

func (r *TestRunner) runNeptunedataEngineTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus", func() error {
		resp, err := tc.client.GetEngineStatus(tc.ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Status == nil || *resp.Status != "healthy" {
			return fmt.Errorf("expected status=healthy, got %v", resp.Status)
		}
		if resp.StartTime == nil || *resp.StartTime == "" {
			return fmt.Errorf("expected non-empty startTime")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus_GremlinVersion", func() error {
		resp, err := tc.client.GetEngineStatus(tc.ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Gremlin == nil || resp.Gremlin.Version == nil || *resp.Gremlin.Version == "" {
			return fmt.Errorf("expected non-empty gremlin version")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus_OpenCypherVersion", func() error {
		resp, err := tc.client.GetEngineStatus(tc.ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Opencypher == nil || resp.Opencypher.Version == nil || *resp.Opencypher.Version == "" {
			return fmt.Errorf("expected non-empty opencypher version")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus_Role", func() error {
		resp, err := tc.client.GetEngineStatus(tc.ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Role == nil || (*resp.Role != "writer" && *resp.Role != "reader") {
			return fmt.Errorf("expected role=writer|reader, got %v", resp.Role)
		}
		return nil
	}))

	return results
}
