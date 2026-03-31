package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
	"vorpalstacks-sdk-tests/config"
)

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

	client := neptune.NewFromConfig(cfg)
	ctx := context.Background()

	ts := time.Now().UnixNano()
	clusterID := fmt.Sprintf("test-cluster-%d", ts)
	paramGroupName := fmt.Sprintf("test-cpg-%d", ts)
	subnetGroupName := fmt.Sprintf("test-sng-%d", ts)
	snapshotID := fmt.Sprintf("test-snap-%d", ts)
	globalClusterID := fmt.Sprintf("test-global-%d", ts)
	instanceID := fmt.Sprintf("test-inst-%d", ts)

	// === Engine Versions ===

	results = append(results, r.RunTest("neptune", "DescribeDBEngineVersions", func() error {
		resp, err := client.DescribeDBEngineVersions(ctx, &neptune.DescribeDBEngineVersionsInput{
			Engine: aws.String("neptune"),
		})
		if err != nil {
			return err
		}
		if len(resp.DBEngineVersions) == 0 {
			return fmt.Errorf("expected at least one engine version")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBEngineVersions_DefaultEngine", func() error {
		resp, err := client.DescribeDBEngineVersions(ctx, &neptune.DescribeDBEngineVersionsInput{})
		if err != nil {
			return err
		}
		if len(resp.DBEngineVersions) == 0 {
			return fmt.Errorf("expected at least one engine version with default engine")
		}
		return nil
	}))

	// === Cluster Parameter Groups ===

	results = append(results, r.RunTest("neptune", "CreateDBClusterParameterGroup", func() error {
		_, err := client.CreateDBClusterParameterGroup(ctx, &neptune.CreateDBClusterParameterGroupInput{
			DBClusterParameterGroupName: aws.String(paramGroupName),
			DBParameterGroupFamily:      aws.String("neptune1"),
			Description:                 aws.String("Test cluster parameter group"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterParameterGroups", func() error {
		resp, err := client.DescribeDBClusterParameterGroups(ctx, &neptune.DescribeDBClusterParameterGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, pg := range resp.DBClusterParameterGroups {
			if pg.DBClusterParameterGroupName != nil && *pg.DBClusterParameterGroupName == paramGroupName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created parameter group not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterParameterGroups_FilterByName", func() error {
		resp, err := client.DescribeDBClusterParameterGroups(ctx, &neptune.DescribeDBClusterParameterGroupsInput{
			DBClusterParameterGroupName: aws.String(paramGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.DBClusterParameterGroups) != 1 {
			return fmt.Errorf("expected 1 parameter group, got %d", len(resp.DBClusterParameterGroups))
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterParameters", func() error {
		_, err := client.DescribeDBClusterParameters(ctx, &neptune.DescribeDBClusterParametersInput{
			DBClusterParameterGroupName: aws.String(paramGroupName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeEngineDefaultClusterParameters", func() error {
		_, err := client.DescribeEngineDefaultClusterParameters(ctx, &neptune.DescribeEngineDefaultClusterParametersInput{
			DBParameterGroupFamily: aws.String("neptune1"),
		})
		return err
	}))

	// NOTE: CopyDBClusterParameterGroup test skipped — SDK v1.44.3 has a bug
	// where validation requires TargetDBClusterParameterGroupDescription but the
	// field is not present on the input struct.

	// === DB Subnet Groups ===

	results = append(results, r.RunTest("neptune", "CreateDBSubnetGroup", func() error {
		_, err := client.CreateDBSubnetGroup(ctx, &neptune.CreateDBSubnetGroupInput{
			DBSubnetGroupName:        aws.String(subnetGroupName),
			DBSubnetGroupDescription: aws.String("Test subnet group"),
			SubnetIds:                []string{"subnet-aaa111", "subnet-bbb222"},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBSubnetGroups", func() error {
		resp, err := client.DescribeDBSubnetGroups(ctx, &neptune.DescribeDBSubnetGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, sg := range resp.DBSubnetGroups {
			if sg.DBSubnetGroupName != nil && *sg.DBSubnetGroupName == subnetGroupName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created subnet group not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBSubnetGroup", func() error {
		_, err := client.ModifyDBSubnetGroup(ctx, &neptune.ModifyDBSubnetGroupInput{
			DBSubnetGroupName:        aws.String(subnetGroupName),
			DBSubnetGroupDescription: aws.String("Modified test subnet group"),
			SubnetIds:                []string{"subnet-aaa111", "subnet-bbb222", "subnet-ccc333"},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBSubnetGroups_FilterByName", func() error {
		resp, err := client.DescribeDBSubnetGroups(ctx, &neptune.DescribeDBSubnetGroupsInput{
			DBSubnetGroupName: aws.String(subnetGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.DBSubnetGroups) != 1 {
			return fmt.Errorf("expected 1 subnet group, got %d", len(resp.DBSubnetGroups))
		}
		return nil
	}))

	// === DB Clusters (CRUD lifecycle) ===

	results = append(results, r.RunTest("neptune", "CreateDBCluster", func() error {
		_, err := client.CreateDBCluster(ctx, &neptune.CreateDBClusterInput{
			DBClusterIdentifier:         aws.String(clusterID),
			Engine:                      aws.String("neptune"),
			MasterUsername:              aws.String("admin"),
			MasterUserPassword:          aws.String("Pass123456"),
			Port:                        aws.Int32(8182),
			DBClusterParameterGroupName: aws.String(paramGroupName),
			DBSubnetGroupName:           aws.String(subnetGroupName),
			BackupRetentionPeriod:       aws.Int32(7),
			DeletionProtection:          aws.Bool(false),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters", func() error {
		resp, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{})
		if err != nil {
			return err
		}
		found := false
		for _, c := range resp.DBClusters {
			if c.DBClusterIdentifier != nil && *c.DBClusterIdentifier == clusterID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created cluster not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters_FilterByID", func() error {
		resp, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(clusterID),
		})
		if err != nil {
			return err
		}
		if len(resp.DBClusters) != 1 {
			return fmt.Errorf("expected 1 cluster, got %d", len(resp.DBClusters))
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters_ContentVerify", func() error {
		resp, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(clusterID),
		})
		if err != nil {
			return err
		}
		c := resp.DBClusters[0]
		if c.Engine == nil || *c.Engine != "neptune" {
			return fmt.Errorf("expected engine=neptune, got %v", c.Engine)
		}
		if c.Port == nil || *c.Port != 8182 {
			return fmt.Errorf("expected port=8182, got %v", c.Port)
		}
		if c.MasterUsername == nil || *c.MasterUsername != "admin" {
			return fmt.Errorf("expected master username=admin, got %v", c.MasterUsername)
		}
		if c.BackupRetentionPeriod == nil || *c.BackupRetentionPeriod != 7 {
			return fmt.Errorf("expected backup retention=7, got %v", c.BackupRetentionPeriod)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBCluster", func() error {
		_, err := client.ModifyDBCluster(ctx, &neptune.ModifyDBClusterInput{
			DBClusterIdentifier:   aws.String(clusterID),
			BackupRetentionPeriod: aws.Int32(14),
			Port:                  aws.Int32(8183),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBCluster_Verify", func() error {
		resp, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(clusterID),
		})
		if err != nil {
			return err
		}
		c := resp.DBClusters[0]
		if c.BackupRetentionPeriod == nil || *c.BackupRetentionPeriod != 14 {
			return fmt.Errorf("expected backup retention=14 after modify, got %v", c.BackupRetentionPeriod)
		}
		if c.Port == nil || *c.Port != 8183 {
			return fmt.Errorf("expected port=8183 after modify, got %v", c.Port)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "AddRoleToDBCluster", func() error {
		_, err := client.AddRoleToDBCluster(ctx, &neptune.AddRoleToDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
			RoleArn:             aws.String("arn:aws:iam::000000000000:role/test-role"),
			FeatureName:         aws.String("s3Import"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "RemoveRoleFromDBCluster", func() error {
		_, err := client.RemoveRoleFromDBCluster(ctx, &neptune.RemoveRoleFromDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
			RoleArn:             aws.String("arn:aws:iam::000000000000:role/test-role"),
			FeatureName:         aws.String("s3Import"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "StopDBCluster", func() error {
		_, err := client.StopDBCluster(ctx, &neptune.StopDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "StartDBCluster", func() error {
		_, err := client.StartDBCluster(ctx, &neptune.StartDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
		})
		return err
	}))

	// === Snapshots ===

	results = append(results, r.RunTest("neptune", "CreateDBClusterSnapshot", func() error {
		_, err := client.CreateDBClusterSnapshot(ctx, &neptune.CreateDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: aws.String(snapshotID),
			DBClusterIdentifier:         aws.String(clusterID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterSnapshots", func() error {
		resp, err := client.DescribeDBClusterSnapshots(ctx, &neptune.DescribeDBClusterSnapshotsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, snap := range resp.DBClusterSnapshots {
			if snap.DBClusterSnapshotIdentifier != nil && *snap.DBClusterSnapshotIdentifier == snapshotID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created snapshot not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterSnapshots_ContentVerify", func() error {
		resp, err := client.DescribeDBClusterSnapshots(ctx, &neptune.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: aws.String(snapshotID),
		})
		if err != nil {
			return err
		}
		snap := resp.DBClusterSnapshots[0]
		if snap.Engine == nil || *snap.Engine != "neptune" {
			return fmt.Errorf("expected engine=neptune, got %v", snap.Engine)
		}
		if snap.DBClusterIdentifier == nil || *snap.DBClusterIdentifier != clusterID {
			return fmt.Errorf("expected source cluster=%s, got %v", clusterID, snap.DBClusterIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterSnapshotAttributes", func() error {
		_, err := client.DescribeDBClusterSnapshotAttributes(ctx, &neptune.DescribeDBClusterSnapshotAttributesInput{
			DBClusterSnapshotIdentifier: aws.String(snapshotID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBClusterSnapshotAttribute", func() error {
		_, err := client.ModifyDBClusterSnapshotAttribute(ctx, &neptune.ModifyDBClusterSnapshotAttributeInput{
			DBClusterSnapshotIdentifier: aws.String(snapshotID),
			AttributeName:               aws.String("restore"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "CopyDBClusterSnapshot", func() error {
		copyID := fmt.Sprintf("test-snap-copy-%d", ts)
		_, err := client.CopyDBClusterSnapshot(ctx, &neptune.CopyDBClusterSnapshotInput{
			SourceDBClusterSnapshotIdentifier: aws.String(snapshotID),
			TargetDBClusterSnapshotIdentifier: aws.String(copyID),
		})
		if err == nil {
			_, _ = client.DeleteDBClusterSnapshot(ctx, &neptune.DeleteDBClusterSnapshotInput{
				DBClusterSnapshotIdentifier: aws.String(copyID),
			})
		}
		return err
	}))

	results = append(results, r.RunTest("neptune", "RestoreDBClusterFromSnapshot", func() error {
		restoreID := fmt.Sprintf("test-restore-%d", ts)
		_, err := client.RestoreDBClusterFromSnapshot(ctx, &neptune.RestoreDBClusterFromSnapshotInput{
			DBClusterIdentifier: aws.String(restoreID),
			SnapshotIdentifier:  aws.String(snapshotID),
			Engine:              aws.String("neptune"),
		})
		if err == nil {
			_, _ = client.DeleteDBCluster(ctx, &neptune.DeleteDBClusterInput{
				DBClusterIdentifier: aws.String(restoreID),
				SkipFinalSnapshot:   aws.Bool(true),
			})
		}
		return err
	}))

	results = append(results, r.RunTest("neptune", "RestoreDBClusterToPointInTime", func() error {
		pitrID := fmt.Sprintf("test-pitr-%d", ts)
		_, err := client.RestoreDBClusterToPointInTime(ctx, &neptune.RestoreDBClusterToPointInTimeInput{
			DBClusterIdentifier:       aws.String(pitrID),
			SourceDBClusterIdentifier: aws.String(clusterID),
		})
		if err == nil {
			_, _ = client.DeleteDBCluster(ctx, &neptune.DeleteDBClusterInput{
				DBClusterIdentifier: aws.String(pitrID),
				SkipFinalSnapshot:   aws.Bool(true),
			})
		}
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterSnapshot", func() error {
		_, err := client.DeleteDBClusterSnapshot(ctx, &neptune.DeleteDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: aws.String(snapshotID),
		})
		return err
	}))

	// === DB Instances ===

	results = append(results, r.RunTest("neptune", "CreateDBInstance", func() error {
		_, err := client.CreateDBInstance(ctx, &neptune.CreateDBInstanceInput{
			DBInstanceIdentifier: aws.String(instanceID),
			DBClusterIdentifier:  aws.String(clusterID),
			Engine:               aws.String("neptune"),
			DBInstanceClass:      aws.String("db.r5.large"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances", func() error {
		resp, err := client.DescribeDBInstances(ctx, &neptune.DescribeDBInstancesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, i := range resp.DBInstances {
			if i.DBInstanceIdentifier != nil && *i.DBInstanceIdentifier == instanceID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created instance not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances_FilterByID", func() error {
		resp, err := client.DescribeDBInstances(ctx, &neptune.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(instanceID),
		})
		if err != nil {
			return err
		}
		if len(resp.DBInstances) != 1 {
			return fmt.Errorf("expected 1 instance, got %d", len(resp.DBInstances))
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBInstance", func() error {
		_, err := client.ModifyDBInstance(ctx, &neptune.ModifyDBInstanceInput{
			DBInstanceIdentifier: aws.String(instanceID),
			DBInstanceClass:      aws.String("db.r5.xlarge"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "RebootDBInstance", func() error {
		_, err := client.RebootDBInstance(ctx, &neptune.RebootDBInstanceInput{
			DBInstanceIdentifier: aws.String(instanceID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBInstance", func() error {
		_, err := client.DeleteDBInstance(ctx, &neptune.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(instanceID),
			SkipFinalSnapshot:    aws.Bool(true),
		})
		return err
	}))

	// === Global Clusters ===

	results = append(results, r.RunTest("neptune", "CreateGlobalCluster", func() error {
		_, err := client.CreateGlobalCluster(ctx, &neptune.CreateGlobalClusterInput{
			GlobalClusterIdentifier: aws.String(globalClusterID),
			Engine:                  aws.String("neptune"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters", func() error {
		resp, err := client.DescribeGlobalClusters(ctx, &neptune.DescribeGlobalClustersInput{})
		if err != nil {
			return err
		}
		found := false
		for _, gc := range resp.GlobalClusters {
			if gc.GlobalClusterIdentifier != nil && *gc.GlobalClusterIdentifier == globalClusterID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created global cluster not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeGlobalClusters_FilterByID", func() error {
		resp, err := client.DescribeGlobalClusters(ctx, &neptune.DescribeGlobalClustersInput{
			GlobalClusterIdentifier: aws.String(globalClusterID),
		})
		if err != nil {
			return err
		}
		if len(resp.GlobalClusters) != 1 {
			return fmt.Errorf("expected 1 global cluster, got %d", len(resp.GlobalClusters))
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyGlobalCluster", func() error {
		_, err := client.ModifyGlobalCluster(ctx, &neptune.ModifyGlobalClusterInput{
			GlobalClusterIdentifier: aws.String(globalClusterID),
			EngineVersion:           aws.String("1.3.2.0"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteGlobalCluster", func() error {
		_, err := client.DeleteGlobalCluster(ctx, &neptune.DeleteGlobalClusterInput{
			GlobalClusterIdentifier: aws.String(globalClusterID),
		})
		return err
	}))

	// === Event Subscriptions ===

	subName := fmt.Sprintf("test-sub-%d", ts)
	results = append(results, r.RunTest("neptune", "CreateEventSubscription", func() error {
		_, err := client.CreateEventSubscription(ctx, &neptune.CreateEventSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SnsTopicArn:      aws.String("arn:aws:sns:us-east-1:000000000000:test-topic"),
			SourceType:       aws.String("db-cluster"),
			SourceIds:        []string{clusterID},
			EventCategories:  []string{"creation", "deletion", "failover"},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeEventSubscriptions", func() error {
		resp, err := client.DescribeEventSubscriptions(ctx, &neptune.DescribeEventSubscriptionsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, sub := range resp.EventSubscriptionsList {
			if sub.CustSubscriptionId != nil && *sub.CustSubscriptionId == subName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created event subscription not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "AddSourceIdentifierToSubscription", func() error {
		_, err := client.AddSourceIdentifierToSubscription(ctx, &neptune.AddSourceIdentifierToSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SourceIdentifier: aws.String("cluster-extra"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "RemoveSourceIdentifierFromSubscription", func() error {
		_, err := client.RemoveSourceIdentifierFromSubscription(ctx, &neptune.RemoveSourceIdentifierFromSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SourceIdentifier: aws.String("cluster-extra"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyEventSubscription", func() error {
		_, err := client.ModifyEventSubscription(ctx, &neptune.ModifyEventSubscriptionInput{
			SubscriptionName: aws.String(subName),
			SourceType:       aws.String("db-instance"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteEventSubscription", func() error {
		_, err := client.DeleteEventSubscription(ctx, &neptune.DeleteEventSubscriptionInput{
			SubscriptionName: aws.String(subName),
		})
		return err
	}))

	// === Descriptive Operations ===

	results = append(results, r.RunTest("neptune", "DescribeEventCategories", func() error {
		_, err := client.DescribeEventCategories(ctx, &neptune.DescribeEventCategoriesInput{})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeEvents", func() error {
		_, err := client.DescribeEvents(ctx, &neptune.DescribeEventsInput{})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribePendingMaintenanceActions", func() error {
		_, err := client.DescribePendingMaintenanceActions(ctx, &neptune.DescribePendingMaintenanceActionsInput{})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeOrderableDBInstanceOptions", func() error {
		_, err := client.DescribeOrderableDBInstanceOptions(ctx, &neptune.DescribeOrderableDBInstanceOptionsInput{
			Engine: aws.String("neptune"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeEngineDefaultParameters", func() error {
		_, err := client.DescribeEngineDefaultParameters(ctx, &neptune.DescribeEngineDefaultParametersInput{
			DBParameterGroupFamily: aws.String("neptune1"),
		})
		return err
	}))

	// === Tags ===

	results = append(results, r.RunTest("neptune", "AddTagsToResource", func() error {
		_, err := client.AddTagsToResource(ctx, &neptune.AddTagsToResourceInput{
			ResourceName: aws.String(fmt.Sprintf("arn:aws:rds:us-east-1:000000000000:cluster:%s", clusterID)),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("sdk-test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &neptune.ListTagsForResourceInput{
			ResourceName: aws.String(fmt.Sprintf("arn:aws:rds:us-east-1:000000000000:cluster:%s", clusterID)),
		})
		if err != nil {
			return err
		}
		if resp.TagList == nil {
			return fmt.Errorf("expected TagList to be non-nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RemoveTagsFromResource", func() error {
		_, err := client.RemoveTagsFromResource(ctx, &neptune.RemoveTagsFromResourceInput{
			ResourceName: aws.String(fmt.Sprintf("arn:aws:rds:us-east-1:000000000000:cluster:%s", clusterID)),
			TagKeys:      []string{"Environment"},
		})
		return err
	}))

	// === Cluster Endpoints ===

	endpointID := fmt.Sprintf("test-ep-%d", ts)
	results = append(results, r.RunTest("neptune", "CreateDBClusterEndpoint", func() error {
		_, err := client.CreateDBClusterEndpoint(ctx, &neptune.CreateDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
			DBClusterIdentifier:         aws.String(clusterID),
			EndpointType:                aws.String("READER"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterEndpoints", func() error {
		_, err := client.DescribeDBClusterEndpoints(ctx, &neptune.DescribeDBClusterEndpointsInput{})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBClusterEndpoint", func() error {
		_, err := client.ModifyDBClusterEndpoint(ctx, &neptune.ModifyDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterEndpoint", func() error {
		_, err := client.DeleteDBClusterEndpoint(ctx, &neptune.DeleteDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(endpointID),
		})
		return err
	}))

	// === Instance Parameter Groups ===

	instancePGName := fmt.Sprintf("test-pg-%d", ts)
	results = append(results, r.RunTest("neptune", "CreateDBParameterGroup", func() error {
		_, err := client.CreateDBParameterGroup(ctx, &neptune.CreateDBParameterGroupInput{
			DBParameterGroupName:   aws.String(instancePGName),
			DBParameterGroupFamily: aws.String("neptune1"),
			Description:            aws.String("Test instance parameter group"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBParameterGroups", func() error {
		resp, err := client.DescribeDBParameterGroups(ctx, &neptune.DescribeDBParameterGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, pg := range resp.DBParameterGroups {
			if pg.DBParameterGroupName != nil && *pg.DBParameterGroupName == instancePGName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created parameter group not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBParameters", func() error {
		_, err := client.DescribeDBParameters(ctx, &neptune.DescribeDBParametersInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBParameterGroup", func() error {
		_, err := client.ModifyDBParameterGroup(ctx, &neptune.ModifyDBParameterGroupInput{
			DBParameterGroupName: aws.String(instancePGName),
			Parameters: []types.Parameter{
				{ParameterName: aws.String("neptune_query_timeout"), ParameterValue: aws.String("180000")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ResetDBParameterGroup", func() error {
		_, err := client.ResetDBParameterGroup(ctx, &neptune.ResetDBParameterGroupInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "CopyDBParameterGroup", func() error {
		copyPGName := fmt.Sprintf("test-pg-copy-%d", ts)
		_, err := client.CopyDBParameterGroup(ctx, &neptune.CopyDBParameterGroupInput{
			SourceDBParameterGroupIdentifier:  aws.String(instancePGName),
			TargetDBParameterGroupIdentifier:  aws.String(copyPGName),
			TargetDBParameterGroupDescription: aws.String("Copied instance parameter group"),
		})
		if err == nil {
			_, _ = client.DeleteDBParameterGroup(ctx, &neptune.DeleteDBParameterGroupInput{
				DBParameterGroupName: aws.String(copyPGName),
			})
		}
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBParameterGroup", func() error {
		_, err := client.DeleteDBParameterGroup(ctx, &neptune.DeleteDBParameterGroupInput{
			DBParameterGroupName: aws.String(instancePGName),
		})
		return err
	}))

	// === Error / Edge Case Tests ===

	results = append(results, r.RunTest("neptune", "DescribeDBClusters_NonExistent", func() error {
		_, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String("nonexistent-cluster"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent cluster")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances_NonExistent", func() error {
		_, err := client.DescribeDBInstances(ctx, &neptune.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String("nonexistent-instance"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent instance")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "CreateDBCluster_Duplicate", func() error {
		_, err := client.CreateDBCluster(ctx, &neptune.CreateDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
			Engine:              aws.String("neptune"),
			MasterUsername:      aws.String("admin"),
			MasterUserPassword:  aws.String("Pass123456"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate cluster creation")
		}
		return nil
	}))

	// === Cleanup ===

	results = append(results, r.RunTest("neptune", "DeleteDBClusterParameterGroup", func() error {
		_, err := client.DeleteDBClusterParameterGroup(ctx, &neptune.DeleteDBClusterParameterGroupInput{
			DBClusterParameterGroupName: aws.String(paramGroupName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBSubnetGroup", func() error {
		_, err := client.DeleteDBSubnetGroup(ctx, &neptune.DeleteDBSubnetGroupInput{
			DBSubnetGroupName: aws.String(subnetGroupName),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBCluster", func() error {
		_, err := client.DeleteDBCluster(ctx, &neptune.DeleteDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
			SkipFinalSnapshot:   aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBCluster_VerifyDeleted", func() error {
		_, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(clusterID),
		})
		if err == nil {
			return fmt.Errorf("expected error after cluster deletion")
		}
		return nil
	}))

	return results
}
