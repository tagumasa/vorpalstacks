package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestContext) testNamedQueries() []TestResult {
	var results []TestResult

	namedQueryName := tc.uniqueName("test-query")
	var namedQueryId string
	results = append(results, tc.runner.RunTest("athena", "CreateNamedQuery", func() error {
		id, err := tc.createNamedQuery(namedQueryName, "SELECT 1", "Test query")
		if err != nil {
			return err
		}
		namedQueryId = id
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetNamedQuery", func() error {
		resp, err := tc.client.GetNamedQuery(tc.ctx, &athena.GetNamedQueryInput{
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
		allIDs, err := tc.allNamedQueryIDs()
		if err != nil {
			return err
		}
		for _, id := range allIDs {
			if id == namedQueryId {
				return nil
			}
		}
		return fmt.Errorf("created named query ID %q not found in list", namedQueryId)
	}))

	// Update then verify: renamed fields must be observable via GetNamedQuery.
	updatedQueryName := tc.uniqueName("updated-query")
	results = append(results, tc.runner.RunTest("athena", "UpdateNamedQuery", func() error {
		resp, err := tc.client.UpdateNamedQuery(tc.ctx, &athena.UpdateNamedQueryInput{
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

		verifyResp, err := tc.client.GetNamedQuery(tc.ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		nq := verifyResp.NamedQuery
		if nq == nil {
			return fmt.Errorf("named query is nil after update")
		}
		if aws.ToString(nq.Name) != updatedQueryName {
			return fmt.Errorf("expected name %q, got %q", updatedQueryName, aws.ToString(nq.Name))
		}
		if aws.ToString(nq.QueryString) != "SELECT 2" {
			return fmt.Errorf("expected query 'SELECT 2', got %q", aws.ToString(nq.QueryString))
		}
		return nil
	}))

	oldNameReusable := tc.uniqueName("oldname-reuse")
	results = append(results, tc.runner.RunTest("athena", "UpdateNamedQuery_OldNameReusable", func() error {
		reusableQueryId, err := tc.createNamedQuery(oldNameReusable, "SELECT 3", "")
		if err != nil {
			return err
		}

		renamedName := tc.uniqueName("renamed-query")
		_, err = tc.client.UpdateNamedQuery(tc.ctx, &athena.UpdateNamedQueryInput{
			NamedQueryId: aws.String(reusableQueryId),
			Name:         aws.String(renamedName),
			Description:  aws.String("Renamed"),
			QueryString:  aws.String("SELECT 4"),
		})
		if err != nil {
			tc.deleteNamedQuery(reusableQueryId)
			return fmt.Errorf("update failed: %w", err)
		}

		secondId, err := tc.createNamedQuery(oldNameReusable, "SELECT 5", "")
		if err != nil {
			tc.deleteNamedQuery(reusableQueryId)
			return fmt.Errorf("creating query with old name should succeed after rename: %w", err)
		}
		defer tc.deleteNamedQuery(secondId)
		defer tc.deleteNamedQuery(reusableQueryId)
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "UpdateNamedQuery_NewNameNotReusable", func() error {
		_, err := tc.client.CreateNamedQuery(tc.ctx, &athena.CreateNamedQueryInput{
			Name:        aws.String(updatedQueryName),
			Database:    aws.String("default"),
			QueryString: aws.String("SELECT duplicate"),
		})
		if err == nil {
			return fmt.Errorf("creating query with renamed name should fail with ResourceAlreadyExistsException")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteNamedQuery", func() error {
		_, err := tc.client.DeleteNamedQuery(tc.ctx, &athena.DeleteNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetNamedQuery(tc.ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String(namedQueryId),
		})
		if err == nil {
			return fmt.Errorf("named query should be deleted but GetNamedQuery succeeded")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "BatchGetNamedQuery", func() error {
		batchNQId1, err := tc.createNamedQuery(tc.uniqueName("batch-nq1"), "SELECT 1", "")
		if err != nil {
			return fmt.Errorf("setup create 1 failed: %w", err)
		}
		batchNQId2, err := tc.createNamedQuery(tc.uniqueName("batch-nq2"), "SELECT 2", "")
		if err != nil {
			tc.deleteNamedQuery(batchNQId1)
			return fmt.Errorf("setup create 2 failed: %w", err)
		}
		defer tc.deleteNamedQuery(batchNQId1)
		defer tc.deleteNamedQuery(batchNQId2)

		resp, err := tc.client.BatchGetNamedQuery(tc.ctx, &athena.BatchGetNamedQueryInput{
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

	return results
}
