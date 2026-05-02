package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestCtx) testNamedQueries() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	namedQueryName := fmt.Sprintf("test-query-%d", time.Now().UnixNano())
	var namedQueryId string
	results = append(results, tc.runner.RunTest("athena", "CreateNamedQuery", func() error {
		resp, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(namedQueryName),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT 1"),
			Description: aws.String("Test query"),
		})
		if err != nil {
			return err
		}
		namedQueryId = aws.ToString(resp.NamedQueryId)
		if namedQueryId == "" {
			return fmt.Errorf("NamedQueryId is empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetNamedQuery", func() error {
		resp, err := client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		nq := resp.NamedQuery
		if nq == nil {
			return fmt.Errorf("named query is nil")
		}
		if aws.ToString(nq.Name) != namedQueryName {
			return fmt.Errorf("expected name %q, got %q", namedQueryName, aws.ToString(nq.Name))
		}
		if aws.ToString(nq.Database) != "default" {
			return fmt.Errorf("expected database 'default', got %q", aws.ToString(nq.Database))
		}
		if aws.ToString(nq.QueryString) != "SELECT 1" {
			return fmt.Errorf("expected query 'SELECT 1', got %q", aws.ToString(nq.QueryString))
		}
		if aws.ToString(nq.Description) != "Test query" {
			return fmt.Errorf("expected description 'Test query', got %q", aws.ToString(nq.Description))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListNamedQueries", func() error {
		resp, err := client.ListNamedQueries(ctx, &athena.ListNamedQueriesInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.NamedQueryIds == nil {
			return fmt.Errorf("named query IDs list is nil")
		}
		var found bool
		for _, id := range resp.NamedQueryIds {
			if id == namedQueryId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created named query ID %q not found in list", namedQueryId)
		}
		return nil
	}))

	updatedQueryName := fmt.Sprintf("updated-query-%d", time.Now().UnixNano())
	results = append(results, tc.runner.RunTest("athena", "UpdateNamedQuery", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "GetNamedQuery_AfterUpdate", func() error {
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

	oldNameReusable := fmt.Sprintf("oldname-reuse-%d", time.Now().UnixNano())
	var reusableQueryId string
	results = append(results, tc.runner.RunTest("athena", "UpdateNamedQuery_OldNameReusable", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "UpdateNamedQuery_NewNameNotReusable", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "DeleteNamedQuery", func() error {
		_, err := client.DeleteNamedQuery(ctx, &athena.DeleteNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err == nil {
			return fmt.Errorf("named query should be deleted but GetNamedQuery succeeded")
		}
		return nil
	}))

	batchNQName1 := fmt.Sprintf("batch-nq1-%d", time.Now().UnixNano())
	batchNQName2 := fmt.Sprintf("batch-nq2-%d", time.Now().UnixNano())
	var batchNQId1, batchNQId2 string
	results = append(results, tc.runner.RunTest("athena", "BatchGetNamedQuery_Setup", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "BatchGetNamedQuery", func() error {
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

	return results
}
