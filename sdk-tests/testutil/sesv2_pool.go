package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2PoolTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	poolName := fmt.Sprintf("test-pool-%d", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateDedicatedIpPool", func() error {
		_, err := tc.client.CreateDedicatedIpPool(tc.ctx, &sesv2.CreateDedicatedIpPoolInput{
			PoolName:    aws.String(poolName),
			ScalingMode: types.ScalingModeStandard,
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetDedicatedIpPool", func() error {
		resp, err := tc.client.GetDedicatedIpPool(tc.ctx, &sesv2.GetDedicatedIpPoolInput{
			PoolName: aws.String(poolName),
		})
		if err != nil {
			return err
		}
		if resp.DedicatedIpPool == nil {
			return fmt.Errorf("DedicatedIpPool is nil")
		}
		if resp.DedicatedIpPool.PoolName == nil || *resp.DedicatedIpPool.PoolName != poolName {
			return fmt.Errorf("expected pool name %s, got %v", poolName, resp.DedicatedIpPool.PoolName)
		}
		if resp.DedicatedIpPool.ScalingMode != types.ScalingModeStandard {
			return fmt.Errorf("expected ScalingMode=Standard, got %s", resp.DedicatedIpPool.ScalingMode)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListDedicatedIpPools", func() error {
		all, err := tc.listAllDedicatedIpPools()
		if err != nil {
			return err
		}
		if !containsPoolName(all, poolName) {
			return fmt.Errorf("pool %s not found in list", poolName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteDedicatedIpPool", func() error {
		_, err := tc.client.DeleteDedicatedIpPool(tc.ctx, &sesv2.DeleteDedicatedIpPoolInput{
			PoolName: aws.String(poolName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetDedicatedIpPool(tc.ctx, &sesv2.GetDedicatedIpPoolInput{
			PoolName: aws.String(poolName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting pool")
		}
		return nil
	}))

	return results
}
