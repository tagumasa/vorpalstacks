package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
)

func (r *TestRunner) runNeptunedataResetTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "ExecuteFastReset_Initiate", func() error {
		resp, err := tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if resp.Payload == nil || resp.Payload.Token == nil || *resp.Payload.Token == "" {
			return fmt.Errorf("expected non-empty fast reset token")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteFastReset_Perform", func() error {
		initResp, err := tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if initResp.Payload == nil || initResp.Payload.Token == nil {
			return fmt.Errorf("expected non-empty token from initiateDatabaseReset")
		}

		performResp, err := tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  initResp.Payload.Token,
		})
		if err != nil {
			return err
		}
		if performResp.Status == nil || *performResp.Status == "" {
			return fmt.Errorf("expected non-empty status from performDatabaseReset")
		}
		return nil
	}))

	return results
}

func (tc *neptunedataContext) fastReset() error {
	initResp, err := tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
		Action: types.ActionInitializeReset,
	})
	if err != nil {
		return err
	}
	if initResp.Payload == nil || initResp.Payload.Token == nil {
		return fmt.Errorf("expected token")
	}
	_, err = tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
		Action: types.ActionPerformReset,
		Token:  initResp.Payload.Token,
	})
	return err
}
