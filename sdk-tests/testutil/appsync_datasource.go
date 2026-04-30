package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncDataSourceTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "CreateDataSource", func() error {
		resp, err := client.CreateDataSource(ctx, &appsync.CreateDataSourceInput{
			ApiId: aws.String(res.gqlApiId),
			Name:  aws.String("testDS"),
			Type:  types.DataSourceTypeNone,
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil {
			return fmt.Errorf("DataSource is nil")
		}
		if *resp.DataSource.Name != "testDS" {
			return fmt.Errorf("expected name testDS, got %s", *resp.DataSource.Name)
		}
		if resp.DataSource.Type != types.DataSourceTypeNone {
			return fmt.Errorf("expected NONE type, got %s", resp.DataSource.Type)
		}
		if resp.DataSource.DataSourceArn == nil {
			return fmt.Errorf("DataSourceArn is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateDataSource_WithDescription", func() error {
		resp, err := client.CreateDataSource(ctx, &appsync.CreateDataSourceInput{
			ApiId:       aws.String(res.gqlApiId),
			Name:        aws.String("testDS2"),
			Type:        types.DataSourceTypeNone,
			Description: aws.String("test description"),
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil {
			return fmt.Errorf("DataSource is nil")
		}
		if resp.DataSource.Description == nil || *resp.DataSource.Description != "test description" {
			return fmt.Errorf("description not set: %v", resp.DataSource.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDataSource", func() error {
		resp, err := client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(res.gqlApiId),
			Name:  aws.String("testDS"),
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil || *resp.DataSource.Name != "testDS" {
			return fmt.Errorf("expected testDS, got %v", resp.DataSource)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDataSource_NonExistent", func() error {
		_, err := client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(res.gqlApiId),
			Name:  aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListDataSources", func() error {
		resp, err := client.ListDataSources(ctx, &appsync.ListDataSourcesInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.DataSources) < 2 {
			return fmt.Errorf("expected at least 2 data sources, got %d", len(resp.DataSources))
		}
		found := false
		for _, ds := range resp.DataSources {
			if ds.Name != nil && *ds.Name == "testDS" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created data source testDS not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateDataSource", func() error {
		resp, err := client.UpdateDataSource(ctx, &appsync.UpdateDataSourceInput{
			ApiId:       aws.String(res.gqlApiId),
			Name:        aws.String("testDS"),
			Type:        types.DataSourceTypeNone,
			Description: aws.String("updated description"),
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil {
			return fmt.Errorf("DataSource is nil")
		}
		if resp.DataSource.Description == nil || *resp.DataSource.Description != "updated description" {
			return fmt.Errorf("description not updated: %v", resp.DataSource.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteDataSource", func() error {
		_, err := client.DeleteDataSource(ctx, &appsync.DeleteDataSourceInput{
			ApiId: aws.String(res.gqlApiId),
			Name:  aws.String("testDS2"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(res.gqlApiId),
			Name:  aws.String("testDS2"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteDataSource_NonExistent", func() error {
		_, err := client.DeleteDataSource(ctx, &appsync.DeleteDataSourceInput{
			ApiId: aws.String(res.gqlApiId),
			Name:  aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}
