package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestCtx) testDatabases() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("athena", "ListDatabases", func() error {
		resp, err := client.ListDatabases(ctx, &athena.ListDatabasesInput{
			CatalogName: aws.String("AwsDataCatalog"),
		})
		if err != nil {
			return err
		}
		if resp.DatabaseList == nil {
			return fmt.Errorf("database list is nil")
		}
		var foundDefault bool
		for _, db := range resp.DatabaseList {
			if aws.ToString(db.Name) == "default" {
				foundDefault = true
			}
		}
		if !foundDefault {
			return fmt.Errorf("'default' database not found in list")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetDatabase", func() error {
		resp, err := client.GetDatabase(ctx, &athena.GetDatabaseInput{
			CatalogName:  aws.String("AwsDataCatalog"),
			DatabaseName: aws.String("default"),
		})
		if err != nil {
			return err
		}
		if resp.Database == nil {
			return fmt.Errorf("database is nil")
		}
		if aws.ToString(resp.Database.Name) != "default" {
			return fmt.Errorf("expected database name 'default', got %q", aws.ToString(resp.Database.Name))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListTableMetadata", func() error {
		resp, err := client.ListTableMetadata(ctx, &athena.ListTableMetadataInput{
			CatalogName:  aws.String("AwsDataCatalog"),
			DatabaseName: aws.String("default"),
		})
		if err != nil {
			return err
		}
		if resp.TableMetadataList == nil {
			return fmt.Errorf("table metadata list is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetTableMetadata_NonExistent", func() error {
		_, err := client.GetTableMetadata(ctx, &athena.GetTableMetadataInput{
			CatalogName:  aws.String("AwsDataCatalog"),
			DatabaseName: aws.String("default"),
			TableName:    aws.String("nonexistent_table_xyz"),
		})
		if err := AssertErrorContains(err, "MetadataException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
