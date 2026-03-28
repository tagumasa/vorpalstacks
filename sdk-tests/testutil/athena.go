package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunAthenaTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "athena",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := athena.NewFromConfig(cfg)
	ctx := context.Background()

	// Test 1: List Work Groups
	results = append(results, r.RunTest("athena", "ListWorkGroups", func() error {
		resp, err := client.ListWorkGroups(ctx, &athena.ListWorkGroupsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.WorkGroups == nil {
			return fmt.Errorf("work groups list is nil")
		}
		return nil
	}))

	// Test 2: Create Work Group
	workGroupName := fmt.Sprintf("test-wg-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "CreateWorkGroup", func() error {
		resp, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(workGroupName),
			Configuration: &types.WorkGroupConfiguration{
				ResultConfiguration: &types.ResultConfiguration{
					OutputLocation: aws.String("s3://test-bucket/athena/"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 3: Get Work Group
	results = append(results, r.RunTest("athena", "GetWorkGroup", func() error {
		resp, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String(workGroupName),
		})
		if err != nil {
			return err
		}
		if resp.WorkGroup == nil {
			return fmt.Errorf("work group is nil")
		}
		return nil
	}))

	// Test 4: List Data Catalogs
	results = append(results, r.RunTest("athena", "ListDataCatalogs", func() error {
		resp, err := client.ListDataCatalogs(ctx, &athena.ListDataCatalogsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.DataCatalogsSummary == nil {
			return fmt.Errorf("data catalogs summary list is nil")
		}
		return nil
	}))

	customCatalogName := fmt.Sprintf("test-catalog-%d", time.Now().UnixNano()%1000000000)

	// Test 5: Create Data Catalog (for GetDataCatalog test)
	results = append(results, r.RunTest("athena", "CreateDataCatalog", func() error {
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

	// Test 6: Get Data Catalog (uses custom catalog - AwsDataCatalog returns error per AWS spec)
	results = append(results, r.RunTest("athena", "GetDataCatalog", func() error {
		resp, err := client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err != nil {
			return err
		}
		if resp.DataCatalog == nil {
			return fmt.Errorf("data catalog is nil")
		}
		return nil
	}))

	// Test 7: List Databases
	results = append(results, r.RunTest("athena", "ListDatabases", func() error {
		resp, err := client.ListDatabases(ctx, &athena.ListDatabasesInput{
			CatalogName: aws.String("AwsDataCatalog"),
		})
		if err != nil {
			return err
		}
		if resp.DatabaseList == nil {
			return fmt.Errorf("database list is nil")
		}
		return nil
	}))

	// Test 8: Get Database
	results = append(results, r.RunTest("athena", "GetDatabase", func() error {
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
		return nil
	}))

	// Test 9: List Table Metadata
	results = append(results, r.RunTest("athena", "ListTableMetadata", func() error {
		resp, err := client.ListTableMetadata(ctx, &athena.ListTableMetadataInput{
			CatalogName:  aws.String("AwsDataCatalog"),
			DatabaseName: aws.String("default"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 10: Create Named Query
	namedQueryName := fmt.Sprintf("test-query-%d", time.Now().UnixNano())
	var namedQueryId string
	results = append(results, r.RunTest("athena", "CreateNamedQuery", func() error {
		resp, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(namedQueryName),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT 1"),
			Description: aws.String("Test query"),
		})
		if err == nil {
			namedQueryId = aws.ToString(resp.NamedQueryId)
		}
		return err
	}))

	// Test 10: Get Named Query
	results = append(results, r.RunTest("athena", "GetNamedQuery", func() error {
		resp, err := client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		if resp.NamedQuery == nil {
			return fmt.Errorf("named query is nil")
		}
		return nil
	}))

	// Test 11: List Named Queries
	results = append(results, r.RunTest("athena", "ListNamedQueries", func() error {
		resp, err := client.ListNamedQueries(ctx, &athena.ListNamedQueriesInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.NamedQueryIds == nil {
			return fmt.Errorf("named query IDs list is nil")
		}
		return nil
	}))

	// Test 12: Update Named Query
	updatedQueryName := fmt.Sprintf("updated-query-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("athena", "UpdateNamedQuery", func() error {
		resp, err := client.UpdateNamedQuery(ctx, &athena.UpdateNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
			Name:         aws.String(updatedQueryName),
			Description:  aws.String("Updated test query"),
			QueryString:  aws.String("SELECT 2"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 12a: Verify updated name via GetNamedQuery
	results = append(results, r.RunTest("athena", "GetNamedQuery_AfterUpdate", func() error {
		resp, err := client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		nq := resp.NamedQuery
		if aws.ToString(nq.Name) != updatedQueryName {
			return fmt.Errorf("expected name %q, got %q", updatedQueryName, aws.ToString(nq.Name))
		}
		if aws.ToString(nq.QueryString) != "SELECT 2" {
			return fmt.Errorf("expected query 'SELECT 2', got %q", aws.ToString(nq.QueryString))
		}
		return nil
	}))

	// Test 12b: Old name should be reusable after rename
	oldNameReusable := fmt.Sprintf("oldname-reuse-%d", time.Now().UnixNano())
	var reusableQueryId string
	results = append(results, r.RunTest("athena", "UpdateNamedQuery_OldNameReusable", func() error {
		createResp, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(oldNameReusable),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT 3"),
		})
		if err != nil {
			return err
		}
		reusableQueryId = aws.ToString(createResp.NamedQueryId)

		renamedName := fmt.Sprintf("renamed-query-%d", time.Now().UnixNano())
		_, err = client.UpdateNamedQuery(ctx, &athena.UpdateNamedQueryInput{
			NamedQueryId: aws.String(reusableQueryId),
			Name:         aws.String(renamedName),
			Description:  aws.String("Renamed"),
			QueryString:  aws.String("SELECT 4"),
		})
		if err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		_, err = client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(oldNameReusable),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT 5"),
		})
		if err != nil {
			return fmt.Errorf("creating query with old name should succeed after rename: %w", err)
		}
		return nil
	}))

	// Test 12c: New name should not be reusable (already exists)
	results = append(results, r.RunTest("athena", "UpdateNamedQuery_NewNameNotReusable", func() error {
		_, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(updatedQueryName),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT duplicate"),
		})
		if err == nil {
			return fmt.Errorf("creating query with renamed name should fail with ResourceAlreadyExistsException")
		}
		return nil
	}))

	if reusableQueryId != "" {
		_, _ = client.DeleteNamedQuery(ctx, &athena.DeleteNamedQueryInput{NamedQueryId: aws.String(reusableQueryId)})
	}

	// Test 13: Delete Named Query
	results = append(results, r.RunTest("athena", "DeleteNamedQuery", func() error {
		resp, err := client.DeleteNamedQuery(ctx, &athena.DeleteNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 13: Start Query Execution
	var queryExecutionId string
	results = append(results, r.RunTest("athena", "StartQueryExecution", func() error {
		resp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SELECT 1"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
			ResultConfiguration: &types.ResultConfiguration{
				OutputLocation: aws.String("s3://test-bucket/athena/"),
			},
		})
		if err == nil {
			queryExecutionId = aws.ToString(resp.QueryExecutionId)
		}
		return err
	}))

	// Test 14: Get Query Execution
	results = append(results, r.RunTest("athena", "GetQueryExecution", func() error {
		resp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: aws.String(queryExecutionId),
		})
		if err != nil {
			return err
		}
		if resp.QueryExecution == nil {
			return fmt.Errorf("query execution is nil")
		}
		return nil
	}))

	// Test 15: List Query Executions
	results = append(results, r.RunTest("athena", "ListQueryExecutions", func() error {
		resp, err := client.ListQueryExecutions(ctx, &athena.ListQueryExecutionsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.QueryExecutionIds == nil {
			return fmt.Errorf("query execution IDs list is nil")
		}
		return nil
	}))

	// Test 16: Stop Query Execution
	// Note: Query may complete very fast after optimization, so check state first
	results = append(results, r.RunTest("athena", "StopQueryExecution", func() error {
		getResp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: aws.String(queryExecutionId),
		})
		if err != nil {
			return err
		}
		state := getResp.QueryExecution.Status.State
		if state == types.QueryExecutionStateQueued || state == types.QueryExecutionStateRunning {
			_, err = client.StopQueryExecution(ctx, &athena.StopQueryExecutionInput{
				QueryExecutionId: aws.String(queryExecutionId),
			})
			return err
		}
		return nil
	}))

	// Test 17: Update Work Group
	results = append(results, r.RunTest("athena", "UpdateWorkGroup", func() error {
		resp, err := client.UpdateWorkGroup(ctx, &athena.UpdateWorkGroupInput{
			WorkGroup:   aws.String(workGroupName),
			Description: aws.String("Updated work group"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 18: Delete Work Group
	results = append(results, r.RunTest("athena", "DeleteWorkGroup", func() error {
		resp, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup:             aws.String(workGroupName),
			RecursiveDeleteOption: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 19: Delete Data Catalog (cleanup)
	results = append(results, r.RunTest("athena", "DeleteDataCatalog", func() error {
		resp, err := client.DeleteDataCatalog(ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(customCatalogName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("athena", "GetWorkGroup_NonExistent", func() error {
		_, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String("nonexistent_wg_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent work group")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeleteWorkGroup_NonExistent", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String("nonexistent_wg_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent work group")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetNamedQuery_NonExistent", func() error {
		_, err := client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent named query")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetDataCatalog_NonExistent", func() error {
		_, err := client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String("nonexistent_catalog_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent catalog")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "CreateWorkGroup_Duplicate", func() error {
		dupWGName := fmt.Sprintf("dup-wg-%d", time.Now().UnixNano()%1000000000)
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(dupWGName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(dupWGName)})

		_, err = client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(dupWGName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate work group")
		}
		return nil
	}))

	return results
}
