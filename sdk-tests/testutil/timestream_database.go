package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
)

func (r *TestRunner) runTimestreamDatabaseTests(tc *tsTestContext) []TestResult {
	var results []TestResult
	dbName := tc.uniqueName("db")

	results = append(results, r.RunTest("timestream", "CreateDatabase_ContentVerify", func() error {
		resp, err := tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dbName),
		})
		if err != nil {
			return err
		}
		if resp.Database == nil {
			return fmt.Errorf("Database is nil")
		}
		if resp.Database.DatabaseName == nil || *resp.Database.DatabaseName != dbName {
			return fmt.Errorf("expected DatabaseName=%s, got %v", dbName, resp.Database.DatabaseName)
		}
		if resp.Database.Arn == nil || *resp.Database.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		if resp.Database.TableCount != 0 {
			return fmt.Errorf("expected TableCount=0 for new DB, got %d", resp.Database.TableCount)
		}
		if resp.Database.CreationTime == nil {
			return fmt.Errorf("CreationTime is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeDatabase_Verify", func() error {
		resp, err := tc.writeClient.DescribeDatabase(tc.ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String(dbName),
		})
		if err != nil {
			return err
		}
		if resp.Database == nil {
			return fmt.Errorf("Database is nil")
		}
		if resp.Database.DatabaseName == nil || *resp.Database.DatabaseName != dbName {
			return fmt.Errorf("expected DatabaseName=%s, got %v", dbName, resp.Database.DatabaseName)
		}
		if resp.Database.Arn == nil || *resp.Database.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		if resp.Database.CreationTime == nil {
			return fmt.Errorf("CreationTime is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListDatabases_ContainsCreated", func() error {
		resp, err := tc.writeClient.ListDatabases(tc.ctx, &timestreamwrite.ListDatabasesInput{
			MaxResults: aws.Int32(20),
		})
		if err != nil {
			return err
		}
		found := false
		for _, db := range resp.Databases {
			if db.DatabaseName != nil && *db.DatabaseName == dbName {
				found = true
				if db.Arn == nil || *db.Arn == "" {
					return fmt.Errorf("listed DB has nil/empty Arn")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("database %s not found in list", dbName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UpdateDatabase_KmsKey", func() error {
		kmsKey := "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
		resp, err := tc.writeClient.UpdateDatabase(tc.ctx, &timestreamwrite.UpdateDatabaseInput{
			DatabaseName: aws.String(dbName),
			KmsKeyId:     aws.String(kmsKey),
		})
		if err != nil {
			return err
		}
		if resp.Database == nil {
			return fmt.Errorf("Database is nil")
		}
		if resp.Database.DatabaseName == nil || *resp.Database.DatabaseName != dbName {
			return fmt.Errorf("expected DatabaseName=%s, got %v", dbName, resp.Database.DatabaseName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DeleteDatabase_Basic", func() error {
		_, err := tc.writeClient.DeleteDatabase(tc.ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(dbName),
		})
		return err
	}))

	tagDBName := tc.uniqueName("tag-db")

	results = append(results, r.RunTest("timestream", "CreateDatabase_WithTags", func() error {
		resp, err := tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(tagDBName),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("dev")},
			},
		})
		if err != nil {
			return err
		}
		if resp.Database == nil {
			return fmt.Errorf("Database is nil")
		}
		if resp.Database.DatabaseName == nil || *resp.Database.DatabaseName != tagDBName {
			return fmt.Errorf("expected DatabaseName=%s, got %v", tagDBName, resp.Database.DatabaseName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeDatabase_TagsVerify", func() error {
		resp, err := tc.writeClient.DescribeDatabase(tc.ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String(tagDBName),
		})
		if err != nil {
			return err
		}
		if resp.Database == nil {
			return fmt.Errorf("Database is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "TagResource_Database", func() error {
		dbARN := fmt.Sprintf("arn:aws:timestream:us-east-1:000000000000:database/%s", tagDBName)
		_, err := tc.writeClient.TagResource(tc.ctx, &timestreamwrite.TagResourceInput{
			ResourceARN: aws.String(dbARN),
			Tags: []types.Tag{
				{Key: aws.String("extra"), Value: aws.String("tag")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "ListTagsForResource_Database", func() error {
		dbARN := fmt.Sprintf("arn:aws:timestream:us-east-1:000000000000:database/%s", tagDBName)
		resp, err := tc.writeClient.ListTagsForResource(tc.ctx, &timestreamwrite.ListTagsForResourceInput{
			ResourceARN: aws.String(dbARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 3 {
			return fmt.Errorf("expected at least 3 tags (2 initial + 1 extra), got %d", len(resp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UntagResource_Database", func() error {
		dbARN := fmt.Sprintf("arn:aws:timestream:us-east-1:000000000000:database/%s", tagDBName)
		_, err := tc.writeClient.UntagResource(tc.ctx, &timestreamwrite.UntagResourceInput{
			ResourceARN: aws.String(dbARN),
			TagKeys:     []string{"extra"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.writeClient.ListTagsForResource(tc.ctx, &timestreamwrite.ListTagsForResourceInput{
			ResourceARN: aws.String(dbARN),
		})
		if err != nil {
			return err
		}
		for _, t := range resp.Tags {
			if t.Key != nil && *t.Key == "extra" {
				return fmt.Errorf("tag 'extra' should have been removed")
			}
		}
		return nil
	}))

	tc.deleteDatabase(tagDBName)

	results = append(results, r.RunTest("timestream", "DescribeDatabase_NonExistent", func() error {
		_, err := tc.writeClient.DescribeDatabase(tc.ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String("nonexistent-db-xyz"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("timestream", "CreateDatabase_Duplicate", func() error {
		dupDBName := tc.uniqueName("dup-db")
		_, err := tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.deleteDatabase(dupDBName)

		_, err = tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})
		return AssertErrorContains(err, "ResourceAlreadyExistsException")
	}))

	return results
}
