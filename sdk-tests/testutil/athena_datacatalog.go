package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestCtx) testDataCatalogs() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("athena", "ListDataCatalogs", func() error {
		resp, err := client.ListDataCatalogs(ctx, &athena.ListDataCatalogsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.DataCatalogsSummary == nil {
			return fmt.Errorf("data catalogs summary list is nil")
		}
		var foundAws bool
		for _, dc := range resp.DataCatalogsSummary {
			if aws.ToString(dc.CatalogName) == "AwsDataCatalog" {
				foundAws = true
				if dc.Type != types.DataCatalogTypeGlue {
					return fmt.Errorf("AwsDataCatalog type: expected GLUE, got %s", dc.Type)
				}
			}
		}
		if !foundAws {
			return fmt.Errorf("AwsDataCatalog not found in list")
		}
		return nil
	}))

	customCatalogName := fmt.Sprintf("test-catalog-%d", time.Now().UnixNano()%1000000000)

	results = append(results, tc.runner.RunTest("athena", "CreateDataCatalog", func() error {
		resp, err := client.CreateDataCatalog(ctx, &athena.CreateDataCatalogInput{
			Name:        aws.String(customCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("Test catalog for GetDataCatalog"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetDataCatalog", func() error {
		resp, err := client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err != nil {
			return err
		}
		dc := resp.DataCatalog
		if dc == nil {
			return fmt.Errorf("data catalog is nil")
		}
		if aws.ToString(dc.Name) != customCatalogName {
			return fmt.Errorf("expected name %q, got %q", customCatalogName, aws.ToString(dc.Name))
		}
		if dc.Type != types.DataCatalogTypeGlue {
			return fmt.Errorf("expected type GLUE, got %s", dc.Type)
		}
		if aws.ToString(dc.Description) != "Test catalog for GetDataCatalog" {
			return fmt.Errorf("expected description 'Test catalog for GetDataCatalog', got %q", aws.ToString(dc.Description))
		}
		return nil
	}))

	udcCatalogName := fmt.Sprintf("udc-cat-%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "UpdateDataCatalog_Setup", func() error {
		_, err := client.CreateDataCatalog(ctx, &athena.CreateDataCatalogInput{
			Name:        aws.String(udcCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("Before update"),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "UpdateDataCatalog", func() error {
		_, err := client.UpdateDataCatalog(ctx, &athena.UpdateDataCatalogInput{
			Name:        aws.String(udcCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("After update"),
			Parameters: map[string]string{
				"key1": "value1",
			},
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "UpdateDataCatalog_Verify", func() error {
		resp, err := client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String(udcCatalogName),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.DataCatalog.Description) != "After update" {
			return fmt.Errorf("expected description 'After update', got %q", aws.ToString(resp.DataCatalog.Description))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteDataCatalog_UDCCleanup", func() error {
		_, err := client.DeleteDataCatalog(ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(udcCatalogName),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteDataCatalog", func() error {
		_, err := client.DeleteDataCatalog(ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err == nil {
			return fmt.Errorf("data catalog should be deleted but GetDataCatalog succeeded")
		}
		return nil
	}))

	return results
}
