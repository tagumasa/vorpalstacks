package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneSnapshotTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "CreateDBClusterSnapshot", func() error {
		resp, err := tc.client.CreateDBClusterSnapshot(tc.ctx, &neptune.CreateDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
			DBClusterIdentifier:         aws.String(tc.clusterID),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterSnapshot == nil {
			return fmt.Errorf("expected DBClusterSnapshot in response")
		}
		snap := resp.DBClusterSnapshot
		if snap.DBClusterSnapshotIdentifier == nil || *snap.DBClusterSnapshotIdentifier != tc.snapshotID {
			return fmt.Errorf("expected snapshot ID=%s, got %v", tc.snapshotID, snap.DBClusterSnapshotIdentifier)
		}
		if snap.DBClusterIdentifier == nil || *snap.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, snap.DBClusterIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterSnapshots", func() error {
		resp, err := tc.client.DescribeDBClusterSnapshots(tc.ctx, &neptune.DescribeDBClusterSnapshotsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, snap := range resp.DBClusterSnapshots {
			if snap.DBClusterSnapshotIdentifier != nil && *snap.DBClusterSnapshotIdentifier == tc.snapshotID {
				found = true
				if snap.DBClusterIdentifier == nil || *snap.DBClusterIdentifier != tc.clusterID {
					return fmt.Errorf("expected DBClusterIdentifier=%s on snapshot, got %v", tc.clusterID, snap.DBClusterIdentifier)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created snapshot not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterSnapshots_ContentVerify", func() error {
		resp, err := tc.client.DescribeDBClusterSnapshots(tc.ctx, &neptune.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err != nil {
			return err
		}
		snap := resp.DBClusterSnapshots[0]
		if snap.Engine == nil || *snap.Engine != "neptune" {
			return fmt.Errorf("expected engine=neptune, got %v", snap.Engine)
		}
		if snap.DBClusterIdentifier == nil || *snap.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected source cluster=%s, got %v", tc.clusterID, snap.DBClusterIdentifier)
		}
		if snap.SnapshotType == nil || *snap.SnapshotType == "" {
			return fmt.Errorf("expected non-empty SnapshotType")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusterSnapshotAttributes", func() error {
		resp, err := tc.client.DescribeDBClusterSnapshotAttributes(tc.ctx, &neptune.DescribeDBClusterSnapshotAttributesInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterSnapshotAttributesResult == nil {
			return fmt.Errorf("expected DBClusterSnapshotAttributesResult in response")
		}
		if resp.DBClusterSnapshotAttributesResult.DBClusterSnapshotIdentifier == nil || *resp.DBClusterSnapshotAttributesResult.DBClusterSnapshotIdentifier != tc.snapshotID {
			return fmt.Errorf("expected DBClusterSnapshotIdentifier=%s, got %v", tc.snapshotID, resp.DBClusterSnapshotAttributesResult.DBClusterSnapshotIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBClusterSnapshotAttribute", func() error {
		resp, err := tc.client.ModifyDBClusterSnapshotAttribute(tc.ctx, &neptune.ModifyDBClusterSnapshotAttributeInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
			AttributeName:               aws.String("restore"),
			ValuesToAdd:                 []string{"all"},
		})
		if err != nil {
			return err
		}
		if resp.DBClusterSnapshotAttributesResult == nil {
			return fmt.Errorf("expected DBClusterSnapshotAttributesResult in response")
		}
		foundRestore := false
		for _, attr := range resp.DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes {
			if attr.AttributeName != nil && *attr.AttributeName == "restore" {
				foundRestore = true
				hasAll := false
				for _, v := range attr.AttributeValues {
					if v == "all" {
						hasAll = true
					}
				}
				if !hasAll {
					return fmt.Errorf("expected 'all' in restore attribute values after modify")
				}
			}
		}
		if !foundRestore {
			return fmt.Errorf("expected 'restore' attribute in ModifyDBClusterSnapshotAttribute response")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBClusterSnapshotAttribute_Verify", func() error {
		resp, err := tc.client.DescribeDBClusterSnapshotAttributes(tc.ctx, &neptune.DescribeDBClusterSnapshotAttributesInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err != nil {
			return err
		}
		attrs := resp.DBClusterSnapshotAttributesResult.DBClusterSnapshotAttributes
		foundRestore := false
		for _, attr := range attrs {
			if attr.AttributeName != nil && *attr.AttributeName == "restore" {
				foundRestore = true
				hasAll := false
				for _, v := range attr.AttributeValues {
					if v == "all" {
						hasAll = true
					}
				}
				if !hasAll {
					return fmt.Errorf("expected 'all' in restore attribute values")
				}
			}
		}
		if !foundRestore {
			return fmt.Errorf("expected 'restore' attribute in snapshot attributes")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "CopyDBClusterSnapshot", func() error {
		copyID := fmt.Sprintf("test-snap-copy-%d", tc.ts)
		resp, err := tc.client.CopyDBClusterSnapshot(tc.ctx, &neptune.CopyDBClusterSnapshotInput{
			SourceDBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
			TargetDBClusterSnapshotIdentifier: aws.String(copyID),
		})
		if err != nil {
			return err
		}
		if resp.DBClusterSnapshot == nil {
			return fmt.Errorf("expected DBClusterSnapshot in CopyDBClusterSnapshot response")
		}
		if resp.DBClusterSnapshot.DBClusterSnapshotIdentifier == nil || *resp.DBClusterSnapshot.DBClusterSnapshotIdentifier != copyID {
			return fmt.Errorf("expected DBClusterSnapshotIdentifier=%s, got %v", copyID, resp.DBClusterSnapshot.DBClusterSnapshotIdentifier)
		}
		_, _ = tc.client.DeleteDBClusterSnapshot(tc.ctx, &neptune.DeleteDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: aws.String(copyID),
		})
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RestoreDBClusterFromSnapshot", func() error {
		restoreID := fmt.Sprintf("test-restore-%d", tc.ts)
		resp, err := tc.client.RestoreDBClusterFromSnapshot(tc.ctx, &neptune.RestoreDBClusterFromSnapshotInput{
			DBClusterIdentifier: aws.String(restoreID),
			SnapshotIdentifier:  aws.String(tc.snapshotID),
			Engine:              aws.String("neptune"),
		})
		if err != nil {
			return err
		}
		if resp.DBCluster == nil {
			return fmt.Errorf("expected DBCluster in RestoreDBClusterFromSnapshot response")
		}
		if resp.DBCluster.DBClusterIdentifier == nil || *resp.DBCluster.DBClusterIdentifier != restoreID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", restoreID, resp.DBCluster.DBClusterIdentifier)
		}
		_, _ = tc.client.DeleteDBCluster(tc.ctx, &neptune.DeleteDBClusterInput{
			DBClusterIdentifier: aws.String(restoreID),
			SkipFinalSnapshot:   aws.Bool(true),
		})
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RestoreDBClusterToPointInTime", func() error {
		pitrID := fmt.Sprintf("test-pitr-%d", tc.ts)
		resp, err := tc.client.RestoreDBClusterToPointInTime(tc.ctx, &neptune.RestoreDBClusterToPointInTimeInput{
			DBClusterIdentifier:       aws.String(pitrID),
			SourceDBClusterIdentifier: aws.String(tc.clusterID),
		})
		if err != nil {
			return err
		}
		if resp.DBCluster == nil {
			return fmt.Errorf("expected DBCluster in RestoreDBClusterToPointInTime response")
		}
		if resp.DBCluster.DBClusterIdentifier == nil || *resp.DBCluster.DBClusterIdentifier != pitrID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", pitrID, resp.DBCluster.DBClusterIdentifier)
		}
		_, _ = tc.client.DeleteDBCluster(tc.ctx, &neptune.DeleteDBClusterInput{
			DBClusterIdentifier: aws.String(pitrID),
			SkipFinalSnapshot:   aws.Bool(true),
		})
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterSnapshot", func() error {
		_, err := tc.client.DeleteDBClusterSnapshot(tc.ctx, &neptune.DeleteDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBClusterSnapshot_VerifyDeleted", func() error {
		resp, err := tc.client.DescribeDBClusterSnapshots(tc.ctx, &neptune.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err != nil {
			if err := AssertErrorContains(err, "DBClusterSnapshotNotFoundFault"); err != nil {
				return err
			}
			return nil
		}
		if len(resp.DBClusterSnapshots) != 0 {
			return fmt.Errorf("expected 0 snapshots after delete, got %d", len(resp.DBClusterSnapshots))
		}
		return nil
	}))

	return results
}
