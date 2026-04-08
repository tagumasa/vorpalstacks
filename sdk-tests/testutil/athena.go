package testutil

import (
	"context"
	"fmt"
	"strings"
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
	// Uses /* SLOW */ hint to ensure query is still running when stop is called.
	results = append(results, r.RunTest("athena", "StopQueryExecution", func() error {
		stopResp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("/* SLOW */ SELECT 1"),
			WorkGroup:   aws.String(workGroupName),
		})
		if err != nil {
			return err
		}
		stopQueryId := stopResp.QueryExecutionId

		_, err = client.StopQueryExecution(ctx, &athena.StopQueryExecutionInput{
			QueryExecutionId: stopQueryId,
		})
		if err != nil {
			return fmt.Errorf("StopQueryExecution failed: %v", err)
		}

		getResp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: stopQueryId,
		})
		if err != nil {
			return err
		}
		if getResp.QueryExecution.Status.State != types.QueryExecutionStateCancelled {
			return fmt.Errorf("expected CANCELLED, got %s", getResp.QueryExecution.Status.State)
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeleteWorkGroup_NonExistent", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String("nonexistent_wg_xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetNamedQuery_NonExistent", func() error {
		_, err := client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetDataCatalog_NonExistent", func() error {
		_, err := client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String("nonexistent_catalog_xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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

	// === TAGGING OPERATIONS ===

	tagWorkGroupName := fmt.Sprintf("tag-wg-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "TagResource_CreateWG", func() error {
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(tagWorkGroupName),
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "TagResource", func() error {
		_, err := client.TagResource(ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", r.region, tagWorkGroupName)),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("athena")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", r.region, tagWorkGroupName)),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.Tags))
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["env"] != "test" {
			return fmt.Errorf("expected env=test, got %q", tagMap["env"])
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &athena.UntagResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", r.region, tagWorkGroupName)),
			TagKeys:     []string{"env"},
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "ListTagsForResource_AfterUntag", func() error {
		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:workgroup/%s", r.region, tagWorkGroupName)),
		})
		if err != nil {
			return err
		}
		for _, t := range resp.Tags {
			if aws.ToString(t.Key) == "env" {
				return fmt.Errorf("env tag should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeleteWorkGroup_TagCleanup", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String(tagWorkGroupName),
		})
		return err
	}))

	// === TAGGING ON DATACATALOG ===

	tagCatalogName := fmt.Sprintf("tag-cat-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "TagResource_DataCatalog", func() error {
		_, err := client.CreateDataCatalog(ctx, &athena.CreateDataCatalogInput{
			Name:        aws.String(tagCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("Catalog for tag test"),
		})
		if err != nil {
			return err
		}
		_, err = client.TagResource(ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:datacatalog/%s", r.region, tagCatalogName)),
			Tags: []types.Tag{
				{Key: aws.String("purpose"), Value: aws.String("testing")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(fmt.Sprintf("arn:aws:athena:%s:000000000000:datacatalog/%s", r.region, tagCatalogName)),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 1 {
			return fmt.Errorf("expected at least 1 tag on datacatalog")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeleteDataCatalog_TagCleanup", func() error {
		_, err := client.DeleteDataCatalog(ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(tagCatalogName),
		})
		return err
	}))

	// === PREPARED STATEMENT OPERATIONS ===

	psWorkGroup := fmt.Sprintf("ps-wg-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "PreparedStatement_CreateWG", func() error {
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(psWorkGroup),
		})
		return err
	}))

	psName := fmt.Sprintf("ps-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "CreatePreparedStatement", func() error {
		_, err := client.CreatePreparedStatement(ctx, &athena.CreatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT * FROM users WHERE id = ?"),
			Description:    aws.String("Test prepared statement"),
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "CreatePreparedStatement_Duplicate", func() error {
		_, err := client.CreatePreparedStatement(ctx, &athena.CreatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT 1"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate prepared statement")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetPreparedStatement", func() error {
		resp, err := client.GetPreparedStatement(ctx, &athena.GetPreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err != nil {
			return err
		}
		if resp.PreparedStatement == nil {
			return fmt.Errorf("prepared statement is nil")
		}
		if aws.ToString(resp.PreparedStatement.QueryStatement) != "SELECT * FROM users WHERE id = ?" {
			return fmt.Errorf("unexpected query statement: %q", aws.ToString(resp.PreparedStatement.QueryStatement))
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "ListPreparedStatements", func() error {
		resp, err := client.ListPreparedStatements(ctx, &athena.ListPreparedStatementsInput{
			WorkGroup: aws.String(psWorkGroup),
		})
		if err != nil {
			return err
		}
		if len(resp.PreparedStatements) == 0 {
			return fmt.Errorf("expected at least 1 prepared statement")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "UpdatePreparedStatement", func() error {
		_, err := client.UpdatePreparedStatement(ctx, &athena.UpdatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT * FROM orders WHERE id = ?"),
			Description:    aws.String("Updated prepared statement"),
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "GetPreparedStatement_AfterUpdate", func() error {
		resp, err := client.GetPreparedStatement(ctx, &athena.GetPreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.PreparedStatement.QueryStatement) != "SELECT * FROM orders WHERE id = ?" {
			return fmt.Errorf("expected updated query statement, got %q", aws.ToString(resp.PreparedStatement.QueryStatement))
		}
		return nil
	}))

	psName2 := fmt.Sprintf("ps2-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "CreatePreparedStatement_Second", func() error {
		_, err := client.CreatePreparedStatement(ctx, &athena.CreatePreparedStatementInput{
			StatementName:  aws.String(psName2),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT 1"),
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "BatchGetPreparedStatement", func() error {
		resp, err := client.BatchGetPreparedStatement(ctx, &athena.BatchGetPreparedStatementInput{
			PreparedStatementNames: []string{psName, psName2, "nonexistent_ps"},
			WorkGroup:              aws.String(psWorkGroup),
		})
		if err != nil {
			return err
		}
		if len(resp.PreparedStatements) != 2 {
			return fmt.Errorf("expected 2 prepared statements, got %d", len(resp.PreparedStatements))
		}
		if len(resp.UnprocessedPreparedStatementNames) != 1 {
			return fmt.Errorf("expected 1 unprocessed, got %d", len(resp.UnprocessedPreparedStatementNames))
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeletePreparedStatement", func() error {
		_, err := client.DeletePreparedStatement(ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "GetPreparedStatement_NonExistent", func() error {
		_, err := client.GetPreparedStatement(ctx, &athena.GetPreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err := AssertErrorContains(err, "InvalidRequestException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeletePreparedStatement_NonExistent", func() error {
		_, err := client.DeletePreparedStatement(ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err := AssertErrorContains(err, "InvalidRequestException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "DeleteWorkGroup_PSCleanup", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String(psWorkGroup),
		})
		return err
	}))

	// === BATCH GET NAMED QUERY ===

	batchNQName1 := fmt.Sprintf("batch-nq1-%d", time.Now().UnixNano())
	batchNQName2 := fmt.Sprintf("batch-nq2-%d", time.Now().UnixNano())
	var batchNQId1, batchNQId2 string
	results = append(results, r.RunTest("athena", "BatchGetNamedQuery_Setup", func() error {
		resp1, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(batchNQName1),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT 1"),
		})
		if err != nil {
			return err
		}
		batchNQId1 = aws.ToString(resp1.NamedQueryId)

		resp2, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(batchNQName2),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT 2"),
		})
		if err != nil {
			return err
		}
		batchNQId2 = aws.ToString(resp2.NamedQueryId)
		return nil
	}))

	results = append(results, r.RunTest("athena", "BatchGetNamedQuery", func() error {
		resp, err := client.BatchGetNamedQuery(ctx, &athena.BatchGetNamedQueryInput{
			NamedQueryIds: []string{batchNQId1, batchNQId2, "nonexistent-id"},
		})
		if err != nil {
			return err
		}
		if len(resp.NamedQueries) != 2 {
			return fmt.Errorf("expected 2 named queries, got %d", len(resp.NamedQueries))
		}
		if len(resp.UnprocessedNamedQueryIds) != 1 {
			return fmt.Errorf("expected 1 unprocessed, got %d", len(resp.UnprocessedNamedQueryIds))
		}
		return nil
	}))

	if batchNQId1 != "" {
		client.DeleteNamedQuery(ctx, &athena.DeleteNamedQueryInput{NamedQueryId: aws.String(batchNQId1)})
	}
	if batchNQId2 != "" {
		client.DeleteNamedQuery(ctx, &athena.DeleteNamedQueryInput{NamedQueryId: aws.String(batchNQId2)})
	}

	// === GET QUERY RESULTS & RUNTIME STATISTICS ===

	var resultsQueryId string
	results = append(results, r.RunTest("athena", "GetQueryResults_StartQuery", func() error {
		resp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SHOW DATABASES"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
		})
		if err != nil {
			return err
		}
		resultsQueryId = aws.ToString(resp.QueryExecutionId)
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetQueryResults_WaitForCompletion", func() error {
		for i := 0; i < 30; i++ {
			resp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
				QueryExecutionId: aws.String(resultsQueryId),
			})
			if err != nil {
				return err
			}
			state := resp.QueryExecution.Status.State
			if state == types.QueryExecutionStateSucceeded {
				return nil
			}
			if state == types.QueryExecutionStateFailed || state == types.QueryExecutionStateCancelled {
				return fmt.Errorf("query ended in state %s", state)
			}
			time.Sleep(500 * time.Millisecond)
		}
		return fmt.Errorf("query did not complete within timeout")
	}))

	results = append(results, r.RunTest("athena", "GetQueryResults", func() error {
		resp, err := client.GetQueryResults(ctx, &athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(resultsQueryId),
		})
		if err != nil {
			return err
		}
		if resp.ResultSet == nil {
			return fmt.Errorf("result set is nil")
		}
		if resp.ResultSet.ResultSetMetadata == nil {
			return fmt.Errorf("result set metadata is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("athena", "GetQueryRuntimeStatistics", func() error {
		resp, err := client.GetQueryRuntimeStatistics(ctx, &athena.GetQueryRuntimeStatisticsInput{
			QueryExecutionId: aws.String(resultsQueryId),
		})
		if err != nil {
			return err
		}
		if resp.QueryRuntimeStatistics == nil {
			return fmt.Errorf("query runtime statistics is nil")
		}
		return nil
	}))

	// === UPDATE DATA CATALOG ===

	udcCatalogName := fmt.Sprintf("udc-cat-%d", time.Now().UnixNano()%1000000000)
	results = append(results, r.RunTest("athena", "UpdateDataCatalog_Setup", func() error {
		_, err := client.CreateDataCatalog(ctx, &athena.CreateDataCatalogInput{
			Name:        aws.String(udcCatalogName),
			Type:        types.DataCatalogTypeGlue,
			Description: aws.String("Before update"),
		})
		return err
	}))

	results = append(results, r.RunTest("athena", "UpdateDataCatalog", func() error {
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

	results = append(results, r.RunTest("athena", "UpdateDataCatalog_Verify", func() error {
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

	results = append(results, r.RunTest("athena", "DeleteDataCatalog_UDCCleanup", func() error {
		_, err := client.DeleteDataCatalog(ctx, &athena.DeleteDataCatalogInput{
			Name: aws.String(udcCatalogName),
		})
		return err
	}))

	// === LIST ENGINE VERSIONS ===

	results = append(results, r.RunTest("athena", "ListEngineVersions", func() error {
		resp, err := client.ListEngineVersions(ctx, &athena.ListEngineVersionsInput{})
		if err != nil {
			return err
		}
		if len(resp.EngineVersions) == 0 {
			return fmt.Errorf("expected at least 1 engine version")
		}
		return nil
	}))

	// === BATCH GET QUERY EXECUTION ===

	var batchQEId1, batchQEId2 string
	results = append(results, r.RunTest("athena", "BatchGetQueryExecution_Setup", func() error {
		resp1, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SELECT 1"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
		})
		if err != nil {
			return err
		}
		batchQEId1 = aws.ToString(resp1.QueryExecutionId)

		resp2, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SELECT 2"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
		})
		if err != nil {
			return err
		}
		batchQEId2 = aws.ToString(resp2.QueryExecutionId)

		for i := 0; i < 30; i++ {
			r1, _ := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{QueryExecutionId: aws.String(batchQEId1)})
			r2, _ := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{QueryExecutionId: aws.String(batchQEId2)})
			if r1.QueryExecution.Status.State == types.QueryExecutionStateSucceeded &&
				r2.QueryExecution.Status.State == types.QueryExecutionStateSucceeded {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
		return fmt.Errorf("queries did not complete within timeout")
	}))

	results = append(results, r.RunTest("athena", "BatchGetQueryExecution", func() error {
		resp, err := client.BatchGetQueryExecution(ctx, &athena.BatchGetQueryExecutionInput{
			QueryExecutionIds: []string{batchQEId1, batchQEId2, "nonexistent-qe-id"},
		})
		if err != nil {
			return err
		}
		if len(resp.QueryExecutions) != 2 {
			return fmt.Errorf("expected 2 query executions, got %d", len(resp.QueryExecutions))
		}
		if len(resp.UnprocessedQueryExecutionIds) != 1 {
			return fmt.Errorf("expected 1 unprocessed, got %d", len(resp.UnprocessedQueryExecutionIds))
		}
		return nil
	}))

	// === GET TABLE METADATA (NonExistent — requires real table data to test fully) ===

	results = append(results, r.RunTest("athena", "GetTableMetadata_NonExistent", func() error {
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

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("athena", "ListWorkGroups_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano()%1000000000)
		var pgWorkGroups []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagWG-%s-%d", pgTs, i)
			_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
				Name: aws.String(name),
			})
			if err != nil {
				for _, wn := range pgWorkGroups {
					client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wn)})
				}
				return fmt.Errorf("create work group %s: %v", name, err)
			}
			pgWorkGroups = append(pgWorkGroups, name)
		}

		var allWorkGroups []string
		var nextToken *string
		for {
			resp, err := client.ListWorkGroups(ctx, &athena.ListWorkGroupsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, wn := range pgWorkGroups {
					client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wn)})
				}
				return fmt.Errorf("list work groups page: %v", err)
			}
			for _, wg := range resp.WorkGroups {
				if wg.Name != nil && strings.Contains(*wg.Name, "PagWG-"+pgTs) {
					allWorkGroups = append(allWorkGroups, *wg.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, wn := range pgWorkGroups {
			client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wn)})
		}
		if len(allWorkGroups) != 5 {
			return fmt.Errorf("expected 5 paginated work groups, got %d", len(allWorkGroups))
		}
		return nil
	}))

	return results
}
