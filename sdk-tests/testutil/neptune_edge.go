package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneEdgeTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters_InvalidMarker", func() error {
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{
			Marker: aws.String("invalid-nonexistent-marker-xyz"),
		})
		if err != nil {
			return fmt.Errorf("invalid marker should return empty results, not error: %v", err)
		}
		if len(resp.GlobalClusters) != 0 {
			return fmt.Errorf("expected 0 global clusters for invalid marker, got %d", len(resp.GlobalClusters))
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters_NonExistent", func() error {
		_, err := tc.client.DescribeDBClusters(tc.ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String("nonexistent-cluster"),
		})
		if err := AssertErrorContains(err, "DBClusterNotFoundFault"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances_NonExistent", func() error {
		_, err := tc.client.DescribeDBInstances(tc.ctx, &neptune.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String("nonexistent-instance"),
		})
		if err := AssertErrorContains(err, "DBInstanceNotFoundFault"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "CreateDBCluster_Duplicate", func() error {
		_, err := tc.client.CreateDBCluster(tc.ctx, &neptune.CreateDBClusterInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
			Engine:              aws.String("neptune"),
			MasterUsername:      aws.String("admin"),
			MasterUserPassword:  aws.String("Pass123456"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate cluster creation")
		}
		if err := AssertErrorContains(err, "DBClusterAlreadyExistsFault"); err != nil {
			return fmt.Errorf("expected DBClusterAlreadyExistsFault, got: %v", err)
		}
		return nil
	}))

	return results
}
