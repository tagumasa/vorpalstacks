package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
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
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{})
		if err != nil {
			return err
		}
		found := false
		for _, gc := range resp.GlobalClusters {
			if gc.GlobalClusterIdentifier != nil && *gc.GlobalClusterIdentifier == tc.globalClusterID {
				found = true
				if gc.Engine == nil || *gc.Engine != "neptune" {
					return fmt.Errorf("expected Engine=neptune on global cluster, got %v", gc.Engine)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created global cluster not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters_FilterByID", func() error {
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
		})
		if err != nil {
			return err
		}
		if len(resp.GlobalClusters) != 1 {
			return fmt.Errorf("expected 1 global cluster, got %d", len(resp.GlobalClusters))
		}
		gc := resp.GlobalClusters[0]
		if gc.GlobalClusterIdentifier == nil || *gc.GlobalClusterIdentifier != tc.globalClusterID {
			return fmt.Errorf("expected GlobalClusterIdentifier=%s, got %v", tc.globalClusterID, gc.GlobalClusterIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters_ContentVerify", func() error {
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
		})
		if err != nil {
			return err
		}
		if len(resp.GlobalClusters) != 1 {
			return fmt.Errorf("expected 1 global cluster, got %d", len(resp.GlobalClusters))
		}
		gc := resp.GlobalClusters[0]
		if gc.Engine == nil || *gc.Engine != "neptune" {
			return fmt.Errorf("expected engine=neptune, got %v", gc.Engine)
		}
		if gc.GlobalClusterIdentifier == nil || *gc.GlobalClusterIdentifier != tc.globalClusterID {
			return fmt.Errorf("expected globalClusterId=%s, got %v", tc.globalClusterID, gc.GlobalClusterIdentifier)
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
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{
			GlobalClusterIdentifier: aws.String(tc.globalClusterID),
		})
		if err != nil {
			return err
		}
		if len(resp.GlobalClusters) != 1 {
			return fmt.Errorf("expected 1 global cluster, got %d", len(resp.GlobalClusters))
		}
		gc := resp.GlobalClusters[0]
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
		if err != nil {
			if err := AssertErrorContains(err, "GlobalClusterNotFoundFault"); err != nil {
				return err
			}
			return nil
		}
		if len(resp.GlobalClusters) != 0 {
			return fmt.Errorf("expected 0 global clusters after delete, got %d", len(resp.GlobalClusters))
		}
		return nil
	}))

	return results
}
