package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func (r *TestRunner) runSSMBatch(tc *ssmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("ssm", "DeleteParameters_Success", func() error {
		n1 := tc.uniqueName("/batch-del1")
		n2 := tc.uniqueName("/batch-del2")
		for _, n := range []string{n1, n2} {
			_, err := tc.putParam(n, "batch", types.ParameterTypeString)
			if err != nil {
				return fmt.Errorf("put %s: %v", n, err)
			}
		}

		resp, err := tc.client.DeleteParameters(tc.ctx, &ssm.DeleteParametersInput{
			Names: []string{n1, n2},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		if len(resp.DeletedParameters) != 2 {
			return fmt.Errorf("expected 2 deleted, got %d", len(resp.DeletedParameters))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameters_MixedValidInvalid", func() error {
		n1 := tc.uniqueName("/batch-mixed")
		_, err := tc.putParam(n1, "batch", types.ParameterTypeString)
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := tc.client.DeleteParameters(tc.ctx, &ssm.DeleteParametersInput{
			Names: []string{n1, "/nonexistent/batch-xyz"},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		if len(resp.DeletedParameters) != 1 {
			return fmt.Errorf("expected 1 deleted, got %d", len(resp.DeletedParameters))
		}
		if resp.DeletedParameters[0] != n1 {
			return fmt.Errorf("expected deleted param %q, got %q", n1, resp.DeletedParameters[0])
		}
		if len(resp.InvalidParameters) != 1 {
			return fmt.Errorf("expected 1 invalid, got %d", len(resp.InvalidParameters))
		}
		return nil
	}))

	return results
}
