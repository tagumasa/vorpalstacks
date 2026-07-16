package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"vorpalstacks-sdk-tests/config"
)

type neptuneContext struct {
	client  *neptune.Client
	ec2User *ec2.Client
	ctx     context.Context
	region  string
	ts      int64

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
