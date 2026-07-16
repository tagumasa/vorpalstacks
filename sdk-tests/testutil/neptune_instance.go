package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
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
		resp, err := tc.client.DescribeDBInstances(tc.ctx, &neptune.DescribeDBInstancesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, i := range resp.DBInstances {
			if i.DBInstanceIdentifier != nil && *i.DBInstanceIdentifier == tc.instanceID {
				found = true
				if i.Engine == nil || *i.Engine != "neptune" {
					return fmt.Errorf("expected Engine=neptune on instance, got %v", i.Engine)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created instance not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBInstances_FilterByID", func() error {
		resp, err := tc.client.DescribeDBInstances(tc.ctx, &neptune.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
		})
		if err != nil {
			return err
		}
		if len(resp.DBInstances) != 1 {
			return fmt.Errorf("expected 1 instance, got %d", len(resp.DBInstances))
		}
		dbi := resp.DBInstances[0]
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
		_, err := tc.client.ModifyDBInstance(tc.ctx, &neptune.ModifyDBInstanceInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
			DBInstanceClass:      aws.String("db.r5.xlarge"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBInstance_Verify", func() error {
		resp, err := tc.client.DescribeDBInstances(tc.ctx, &neptune.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(tc.instanceID),
		})
		if err != nil {
			return err
		}
		if len(resp.DBInstances) != 1 {
			return fmt.Errorf("expected 1 instance, got %d", len(resp.DBInstances))
		}
		dbi := resp.DBInstances[0]
		if dbi.DBInstanceClass == nil || *dbi.DBInstanceClass != "db.r5.xlarge" {
			return fmt.Errorf("expected DBInstanceClass=db.r5.xlarge after modify, got %v", dbi.DBInstanceClass)
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
		if err := AssertErrorContains(err, "DBInstanceNotFoundFault"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
