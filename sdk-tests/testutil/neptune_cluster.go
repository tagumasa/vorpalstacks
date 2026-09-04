package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func (r *TestRunner) runNeptuneClusterTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "CreateDBCluster", func() error {
		resp, err := tc.client.CreateDBCluster(tc.ctx, &neptune.CreateDBClusterInput{
			DBClusterIdentifier:         aws.String(tc.clusterID),
			Engine:                      aws.String("neptune"),
			MasterUsername:              aws.String("admin"),
			MasterUserPassword:          aws.String("Pass123456"),
			Port:                        aws.Int32(8182),
			DBClusterParameterGroupName: aws.String(tc.paramGroupName),
			DBSubnetGroupName:           aws.String(tc.subnetGroupName),
			BackupRetentionPeriod:       aws.Int32(7),
			DeletionProtection:          aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if resp.DBCluster == nil {
			return fmt.Errorf("expected DBCluster in response")
		}
		if resp.DBCluster.DBClusterIdentifier == nil || *resp.DBCluster.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, resp.DBCluster.DBClusterIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters", func() error {
		clusters, err := tc.allClusters(nil)
		if err != nil {
			return err
		}
		c := containsID(clusters, func(c *types.DBCluster) bool {
			return c.DBClusterIdentifier != nil && *c.DBClusterIdentifier == tc.clusterID
		})
		if c == nil {
			return fmt.Errorf("created cluster not found in list")
		}
		if c.Engine == nil || *c.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune, got %v", c.Engine)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters_FilterByID", func() error {
		_, err := tc.describeCluster(tc.clusterID)
		return err
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBClusters_ContentVerify", func() error {
		c, err := tc.describeCluster(tc.clusterID)
		if err != nil {
			return err
		}
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
		if c.DBClusterParameterGroup == nil || *c.DBClusterParameterGroup != tc.paramGroupName {
			return fmt.Errorf("expected DBClusterParameterGroup=%s, got %v", tc.paramGroupName, c.DBClusterParameterGroup)
		}
		if c.DBSubnetGroup == nil || *c.DBSubnetGroup != tc.subnetGroupName {
			return fmt.Errorf("expected DBSubnetGroup=%s, got %v", tc.subnetGroupName, c.DBSubnetGroup)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBCluster", func() error {
		resp, err := tc.client.ModifyDBCluster(tc.ctx, &neptune.ModifyDBClusterInput{
			DBClusterIdentifier:   aws.String(tc.clusterID),
			BackupRetentionPeriod: aws.Int32(14),
			Port:                  aws.Int32(8183),
		})
		if err != nil {
			return err
		}
		if resp.DBCluster == nil {
			return fmt.Errorf("expected DBCluster in ModifyDBCluster response")
		}
		if resp.DBCluster.Port == nil || *resp.DBCluster.Port != 8183 {
			return fmt.Errorf("expected Port=8183 in ModifyDBCluster response, got %v", resp.DBCluster.Port)
		}
		if resp.DBCluster.BackupRetentionPeriod == nil || *resp.DBCluster.BackupRetentionPeriod != 14 {
			return fmt.Errorf("expected BackupRetentionPeriod=14 in ModifyDBCluster response, got %v", resp.DBCluster.BackupRetentionPeriod)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBCluster_Verify", func() error {
		c, err := tc.describeCluster(tc.clusterID)
		if err != nil {
			return err
		}
		if c.BackupRetentionPeriod == nil || *c.BackupRetentionPeriod != 14 {
			return fmt.Errorf("expected backup retention=14 after modify, got %v", c.BackupRetentionPeriod)
		}
		if c.Port == nil || *c.Port != 8183 {
			return fmt.Errorf("expected port=8183 after modify, got %v", c.Port)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "AddRoleToDBCluster", func() error {
		_, err := tc.client.AddRoleToDBCluster(tc.ctx, &neptune.AddRoleToDBClusterInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
			RoleArn:             aws.String(tc.clusterRoleArn()),
			FeatureName:         aws.String("s3Import"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "AddRoleToDBCluster_Verify", func() error {
		c, err := tc.describeCluster(tc.clusterID)
		if err != nil {
			return err
		}
		found := false
		for _, role := range c.AssociatedRoles {
			if role.RoleArn != nil && *role.RoleArn == tc.clusterRoleArn() {
				found = true
				if role.FeatureName != nil && *role.FeatureName != "s3Import" {
					return fmt.Errorf("expected FeatureName=s3Import, got %v", role.FeatureName)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("expected role arn:aws:iam::%s:role/test-role in AssociatedRoles after AddRoleToDBCluster", tc.accountID)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RemoveRoleFromDBCluster", func() error {
		_, err := tc.client.RemoveRoleFromDBCluster(tc.ctx, &neptune.RemoveRoleFromDBClusterInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
			RoleArn:             aws.String(tc.clusterRoleArn()),
			FeatureName:         aws.String("s3Import"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "RemoveRoleFromDBCluster_Verify", func() error {
		c, err := tc.describeCluster(tc.clusterID)
		if err != nil {
			return err
		}
		for _, role := range c.AssociatedRoles {
			if role.RoleArn != nil && *role.RoleArn == tc.clusterRoleArn() {
				return fmt.Errorf("expected role to be removed from AssociatedRoles after RemoveRoleFromDBCluster")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "StopDBCluster", func() error {
		resp, err := tc.client.StopDBCluster(tc.ctx, &neptune.StopDBClusterInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
		})
		if err != nil {
			return err
		}
		if resp.DBCluster == nil {
			return fmt.Errorf("expected DBCluster in StopDBCluster response")
		}
		if resp.DBCluster.DBClusterIdentifier == nil || *resp.DBCluster.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, resp.DBCluster.DBClusterIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "StartDBCluster", func() error {
		resp, err := tc.client.StartDBCluster(tc.ctx, &neptune.StartDBClusterInput{
			DBClusterIdentifier: aws.String(tc.clusterID),
		})
		if err != nil {
			return err
		}
		if resp.DBCluster == nil {
			return fmt.Errorf("expected DBCluster in StartDBCluster response")
		}
		if resp.DBCluster.DBClusterIdentifier == nil || *resp.DBCluster.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, resp.DBCluster.DBClusterIdentifier)
		}
		return nil
	}))

	return results
}
