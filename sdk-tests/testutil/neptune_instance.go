package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func (r *TestRunner) runNeptuneInstanceTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "CreateDBInstance", func() error {
		resp, err := tc.client.CreateDBInstance(tc.ctx, &neptune.CreateDBInstanceInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
			DBClusterIdentifier:  aws.String(tc.clusterID),
			Engine:               aws.String("neptune"),
			DBInstanceClass:      aws.String("db.r5.large"),
		})
		if err != nil {
			return err
		}
		if resp.DBInstance == nil {
			return fmt.Errorf("expected DBInstance in CreateDBInstance response")
		}
		dbi := resp.DBInstance
		if dbi.DBInstanceIdentifier == nil || *dbi.DBInstanceIdentifier != tc.instanceID {
			return fmt.Errorf("expected DBInstanceIdentifier=%s, got %v", tc.instanceID, dbi.DBInstanceIdentifier)
		}
		if dbi.DBClusterIdentifier == nil || *dbi.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, dbi.DBClusterIdentifier)
		}
		if dbi.Engine == nil || *dbi.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune, got %v", dbi.Engine)
		}
		if dbi.DBInstanceClass == nil || *dbi.DBInstanceClass != "db.r5.large" {
			return fmt.Errorf("expected DBInstanceClass=db.r5.large, got %v", dbi.DBInstanceClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances", func() error {
		instances, err := tc.allInstances()
		if err != nil {
			return err
		}
		i := containsID(instances, func(i *types.DBInstance) bool {
			return i.DBInstanceIdentifier != nil && *i.DBInstanceIdentifier == tc.instanceID
		})
		if i == nil {
			return fmt.Errorf("created instance not found in list")
		}
		if i.Engine == nil || *i.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune on instance, got %v", i.Engine)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances_FilterByID", func() error {
		dbi, err := tc.describeInstance(tc.instanceID)
		if err != nil {
			return err
		}
		if dbi.Engine == nil || *dbi.Engine != "neptune" {
			return fmt.Errorf("expected Engine=neptune, got %v", dbi.Engine)
		}
		if dbi.DBInstanceClass == nil || *dbi.DBInstanceClass != "db.r5.large" {
			return fmt.Errorf("expected DBInstanceClass=db.r5.large, got %v", dbi.DBInstanceClass)
		}
		if dbi.DBClusterIdentifier == nil || *dbi.DBClusterIdentifier != tc.clusterID {
			return fmt.Errorf("expected DBClusterIdentifier=%s, got %v", tc.clusterID, dbi.DBClusterIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBInstance", func() error {
		resp, err := tc.client.ModifyDBInstance(tc.ctx, &neptune.ModifyDBInstanceInput{
			DBInstanceIdentifier:       aws.String(tc.instanceID),
			DBInstanceClass:            aws.String("db.r5.xlarge"),
			PreferredMaintenanceWindow: aws.String("tue:04:00-tue:05:00"),
		})
		if err != nil {
			return err
		}
		if resp.DBInstance == nil || resp.DBInstance.PreferredMaintenanceWindow == nil || *resp.DBInstance.PreferredMaintenanceWindow != "tue:04:00-tue:05:00" {
			return fmt.Errorf("expected PreferredMaintenanceWindow=tue:04:00-tue:05:00 in ModifyDBInstance response, got %v", resp.DBInstance.PreferredMaintenanceWindow)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBInstance_Verify", func() error {
		dbi, err := tc.describeInstance(tc.instanceID)
		if err != nil {
			return err
		}
		if dbi.DBInstanceClass == nil || *dbi.DBInstanceClass != "db.r5.xlarge" {
			return fmt.Errorf("expected DBInstanceClass=db.r5.xlarge after modify, got %v", dbi.DBInstanceClass)
		}
		if dbi.PreferredMaintenanceWindow == nil || *dbi.PreferredMaintenanceWindow != "tue:04:00-tue:05:00" {
			return fmt.Errorf("expected PreferredMaintenanceWindow=tue:04:00-tue:05:00 after modify, got %v", dbi.PreferredMaintenanceWindow)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RebootDBInstance", func() error {
		resp, err := tc.client.RebootDBInstance(tc.ctx, &neptune.RebootDBInstanceInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
		})
		if err != nil {
			return err
		}
		if resp.DBInstance == nil {
			return fmt.Errorf("expected DBInstance in RebootDBInstance response")
		}
		if resp.DBInstance.DBInstanceIdentifier == nil || *resp.DBInstance.DBInstanceIdentifier != tc.instanceID {
			return fmt.Errorf("expected DBInstanceIdentifier=%s, got %v", tc.instanceID, resp.DBInstance.DBInstanceIdentifier)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBInstance", func() error {
		_, err := tc.client.DeleteDBInstance(tc.ctx, &neptune.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
			SkipFinalSnapshot:    aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "DeleteDBInstance_VerifyDeleted", func() error {
		_, err := tc.client.DescribeDBInstances(tc.ctx, &neptune.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
		})
		if err := expectAWSErrorCode(err, "DBInstanceNotFoundFault"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
