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

	// Names is a required member with a modelled length range of 1-10; an
	// empty (non-nil) list passes the SDK's required check but AWS rejects
	// it with ValidationException.
	results = append(results, r.RunTest("ssm", "GetParameters_EmptyNamesRejected", func() error {
		_, err := tc.client.GetParameters(tc.ctx, &ssm.GetParametersInput{
			Names: []string{},
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	// A list above the modelled cap of 10 names is rejected with
	// ValidationException (the SDK does not enforce list length client-side).
	results = append(results, r.RunTest("ssm", "GetParameters_NameListLimit", func() error {
		names := make([]string, 11)
		for i := range names {
			names[i] = tc.uniqueName(fmt.Sprintf("/batch-limit-%d", i))
		}
		_, err := tc.client.GetParameters(tc.ctx, &ssm.GetParametersInput{
			Names: names,
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
