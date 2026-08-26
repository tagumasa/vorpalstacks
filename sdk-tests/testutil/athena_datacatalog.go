package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestContext) testDataCatalogs() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("athena", "ListDataCatalogs", func() error {
		resp, err := tc.client.ListDataCatalogs(tc.ctx, &athena.ListDataCatalogsInput{
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

	customCatalogName := tc.uniqueName("test-catalog")
	// Created here and kept alive for the Get and Delete scenarios below;
	// the DeleteDataCatalog scenario performs the actual deletion.
	results = append(results, tc.runner.RunTest("athena", "CreateDataCatalog", func() error {
		return tc.createDataCatalog(customCatalogName, "Test catalog for GetDataCatalog")
	}))

	results = append(results, tc.runner.RunTest("athena", "GetDataCatalog", func() error {
		resp, err := tc.client.GetDataCatalog(tc.ctx, &athena.GetDataCatalogInput{
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

	// Update then verify in one scenario: the description must change and
	// the update must be observable via GetDataCatalog.
	udcCatalogName := tc.uniqueName("udc-cat")
	results = append(results, tc.runner.RunTest("athena", "UpdateDataCatalog", func() error {
		if err := tc.createDataCatalog(udcCatalogName, "Before update"); err != nil {
			return fmt.Errorf("setup create failed: %w", err)
		}
		defer tc.deleteDataCatalog(udcCatalogName)

		_, err := tc.client.UpdateDataCatalog(tc.ctx, &athena.UpdateDataCatalogInput{
			Name:        aws.String(udcCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("After update"),
			Parameters: map[string]string{
				"key1": "value1",
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetDataCatalog(tc.ctx, &athena.GetDataCatalogInput{
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

	results = append(results, tc.runner.RunTest("athena", "DeleteDataCatalog", func() error {
		_, err := tc.client.DeleteDataCatalog(tc.ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetDataCatalog(tc.ctx, &athena.GetDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err == nil {
			return fmt.Errorf("data catalog should be deleted but GetDataCatalog succeeded")
		}
		return nil
	}))

	return results
}
