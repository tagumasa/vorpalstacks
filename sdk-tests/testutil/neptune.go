package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
	"vorpalstacks-sdk-tests/config"
)

type neptuneContext struct {
	client    *neptune.Client
	ec2User   *ec2.Client
	ctx       context.Context
	region    string
	accountID string
	ts        int64

	clusterID       string
	paramGroupName  string
	subnetGroupName string
	snapshotID      string
	globalClusterID string
	instanceID      string
	vpcID           *string
	subnetIds       []string
	allSubnetIds    []string
}

// clusterRoleArn is the test role ARN used by the
// AddRoleToDBCluster/RemoveRoleFromDBCluster round-trip.
func (tc *neptuneContext) clusterRoleArn() string {
	return fmt.Sprintf("arn:aws:iam::%s:role/test-role", tc.accountID)
}

// The all* walkers below traverse every page of a describe API so the list
// assertions hold during full regression when other suites create resources
// in parallel and a single page may not contain the suite's own resources.

func (tc *neptuneContext) allClusters(filter *string) ([]types.DBCluster, error) {
	return paginate(func(next *string) ([]types.DBCluster, *string, error) {
		resp, err := tc.client.DescribeDBClusters(tc.ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: filter,
			Marker:              next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBClusters, resp.Marker, nil
	})
}

func (tc *neptuneContext) allSnapshots(filter *string) ([]types.DBClusterSnapshot, error) {
	return paginate(func(next *string) ([]types.DBClusterSnapshot, *string, error) {
		resp, err := tc.client.DescribeDBClusterSnapshots(tc.ctx, &neptune.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: filter,
			Marker:                      next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBClusterSnapshots, resp.Marker, nil
	})
}

func (tc *neptuneContext) allInstances() ([]types.DBInstance, error) {
	return paginate(func(next *string) ([]types.DBInstance, *string, error) {
		resp, err := tc.client.DescribeDBInstances(tc.ctx, &neptune.DescribeDBInstancesInput{Marker: next})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBInstances, resp.Marker, nil
	})
}

func (tc *neptuneContext) allGlobalClusters(filter *string) ([]types.GlobalCluster, error) {
	return paginate(func(next *string) ([]types.GlobalCluster, *string, error) {
		resp, err := tc.client.DescribeGlobalClusters(tc.ctx, &neptune.DescribeGlobalClustersInput{
			GlobalClusterIdentifier: filter,
			Marker:                  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.GlobalClusters, resp.Marker, nil
	})
}

func (tc *neptuneContext) allEventSubscriptions() ([]types.EventSubscription, error) {
	return paginate(func(next *string) ([]types.EventSubscription, *string, error) {
		resp, err := tc.client.DescribeEventSubscriptions(tc.ctx, &neptune.DescribeEventSubscriptionsInput{Marker: next})
		if err != nil {
			return nil, nil, err
		}
		return resp.EventSubscriptionsList, resp.Marker, nil
	})
}

func (tc *neptuneContext) allEndpoints(clusterFilter, idFilter *string) ([]types.DBClusterEndpoint, error) {
	return paginate(func(next *string) ([]types.DBClusterEndpoint, *string, error) {
		resp, err := tc.client.DescribeDBClusterEndpoints(tc.ctx, &neptune.DescribeDBClusterEndpointsInput{
			DBClusterIdentifier:         clusterFilter,
			DBClusterEndpointIdentifier: idFilter,
			Marker:                      next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBClusterEndpoints, resp.Marker, nil
	})
}

func (tc *neptuneContext) allSubnetGroups(filter *string) ([]types.DBSubnetGroup, error) {
	return paginate(func(next *string) ([]types.DBSubnetGroup, *string, error) {
		resp, err := tc.client.DescribeDBSubnetGroups(tc.ctx, &neptune.DescribeDBSubnetGroupsInput{
			DBSubnetGroupName: filter,
			Marker:            next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBSubnetGroups, resp.Marker, nil
	})
}

func (tc *neptuneContext) allClusterParameterGroups(filter *string) ([]types.DBClusterParameterGroup, error) {
	return paginate(func(next *string) ([]types.DBClusterParameterGroup, *string, error) {
		resp, err := tc.client.DescribeDBClusterParameterGroups(tc.ctx, &neptune.DescribeDBClusterParameterGroupsInput{
			DBClusterParameterGroupName: filter,
			Marker:                      next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBClusterParameterGroups, resp.Marker, nil
	})
}

func (tc *neptuneContext) allParameterGroups(filter *string) ([]types.DBParameterGroup, error) {
	return paginate(func(next *string) ([]types.DBParameterGroup, *string, error) {
		resp, err := tc.client.DescribeDBParameterGroups(tc.ctx, &neptune.DescribeDBParameterGroupsInput{
			DBParameterGroupName: filter,
			Marker:               next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBParameterGroups, resp.Marker, nil
	})
}

func (r *TestRunner) RunNeptuneTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "neptune",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	ts := time.Now().UnixNano()
	tc := &neptuneContext{
		client:          neptune.NewFromConfig(cfg),
		ec2User:         ec2.NewFromConfig(cfg),
		ctx:             context.Background(),
		region:          r.region,
		accountID:       r.accountID,
		ts:              ts,
		clusterID:       fmt.Sprintf("test-cluster-%d", ts),
		paramGroupName:  fmt.Sprintf("test-cpg-%d", ts),
		subnetGroupName: fmt.Sprintf("test-sng-%d", ts),
		snapshotID:      fmt.Sprintf("test-snap-%d", ts),
		globalClusterID: fmt.Sprintf("test-global-%d", ts),
		instanceID:      fmt.Sprintf("test-inst-%d", ts),
	}

	results = append(results, r.runNeptuneEngineTests(tc)...)
	results = append(results, r.runNeptuneClusterParamGroupTests(tc)...)
	results = append(results, r.runNeptuneSubnetGroupTests(tc)...)
	results = append(results, r.runNeptuneClusterTests(tc)...)
	results = append(results, r.runNeptuneSnapshotTests(tc)...)
	results = append(results, r.runNeptuneInstanceTests(tc)...)
	results = append(results, r.runNeptuneGlobalClusterTests(tc)...)
	results = append(results, r.runNeptuneEventSubscriptionTests(tc)...)
	results = append(results, r.runNeptuneDescriptiveTests(tc)...)
	results = append(results, r.runNeptuneTagTests(tc)...)
	results = append(results, r.runNeptuneClusterEndpointTests(tc)...)
	results = append(results, r.runNeptuneInstanceParamGroupTests(tc)...)
	results = append(results, r.runNeptuneEdgeTests(tc)...)
	results = append(results, r.runNeptuneCleanup(tc)...)

	return results
}

func (r *TestRunner) runNeptuneCleanup(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "DeleteDBClusterParameterGroup", func() error {
		_, err := tc.client.DeleteDBClusterParameterGroup(tc.ctx, &neptune.DeleteDBClusterParameterGroupInput{
			DBClusterParameterGroupName: aws.String(tc.paramGroupName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBSubnetGroup", func() error {
		_, err := tc.client.DeleteDBSubnetGroup(tc.ctx, &neptune.DeleteDBSubnetGroupInput{
			DBSubnetGroupName: aws.String(tc.subnetGroupName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBCluster", func() error {
		_, err := tc.client.DeleteDBCluster(tc.ctx, &neptune.DeleteDBClusterInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
			SkipFinalSnapshot:   aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBCluster_VerifyDeleted", func() error {
		_, err := tc.client.DescribeDBClusters(tc.ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
		})
		if err := AssertErrorContains(err, "DBClusterNotFoundFault"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
