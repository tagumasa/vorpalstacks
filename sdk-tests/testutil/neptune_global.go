package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func (r *TestRunner) runNeptuneGlobalClusterTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "CreateGlobalCluster", func() error {
		resp, err := tc.client.CreateGlobalCluster(tc.ctx, &neptune.CreateGlobalClusterInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
			Engine:                  aws.String("neptune"),
		})
		if err != nil {
			return err
		}
		if resp.GlobalCluster == nil {
			return fmt.Errorf("expected GlobalCluster in response")
		}
		gc := resp.GlobalCluster
		if gc.GlobalClusterIdentifier == nil || *gc.GlobalClusterIdentifier != tc.globalClusterID {
			return fmt.Errorf("expected GlobalClusterIdentifier=%s, got %v", tc.globalClusterID, gc.GlobalClusterIdentifier)
		}
		if gc.Engine == nil || *gc.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune, got %v", gc.Engine)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters", func() error {
		globals, err := tc.allGlobalClusters(nil)
		if err != nil {
			return err
		}
		gc := containsID(globals, func(gc *types.GlobalCluster) bool {
			return gc.GlobalClusterIdentifier != nil && *gc.GlobalClusterIdentifier == tc.globalClusterID
		})
		if gc == nil {
			return fmt.Errorf("created global cluster not found in list")
		}
		if gc.Engine == nil || *gc.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune on global cluster, got %v", gc.Engine)
		}
		return nil
	}))

	// The filtered describe pins count and identifier via describeGlobalCluster;
	// the content check below adds the engine field on the same call.
	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters_ContentVerify", func() error {
		gc, err := tc.describeGlobalCluster(tc.globalClusterID)
		if err != nil {
			return err
		}
		if gc.Engine == nil || *gc.Engine != "neptune" {
			return fmt.Errorf("expected engine=neptune, got %v", gc.Engine)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyGlobalCluster", func() error {
		_, err := tc.client.ModifyGlobalCluster(tc.ctx, &neptune.ModifyGlobalClusterInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
			EngineVersion:           aws.String("1.3.2.0"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyGlobalCluster_Verify", func() error {
		gc, err := tc.describeGlobalCluster(tc.globalClusterID)
		if err != nil {
			return err
		}
		if gc.EngineVersion == nil || *gc.EngineVersion != "1.3.2.0" {
			return fmt.Errorf("expected engineVersion=1.3.2.0 after modify, got %v", gc.EngineVersion)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteGlobalCluster", func() error {
		_, err := tc.client.DeleteGlobalCluster(tc.ctx, &neptune.DeleteGlobalClusterInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteGlobalCluster_VerifyDeleted", func() error {
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
		})
		got := 0
		if resp != nil {
			got = len(resp.GlobalClusters)
		}
		return assertDescribeGone(err, got, "GlobalClusterNotFoundFault", "global clusters")
	}))

	return results
}
