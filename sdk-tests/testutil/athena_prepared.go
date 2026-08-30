package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestContext) testPreparedStatements() []TestResult {
	var results []TestResult

	// Shared work group fixture: one failure surfaces as a single FAIL row
	// named after the setup step it replaced, and the deferred delete runs
	// once every statement scenario below has completed.
	psWorkGroup := tc.uniqueName("ps-wg")
	if err := tc.createWorkGroup(psWorkGroup, nil); err != nil {
		return append(results, TestResult{
			Service:  "athena",
			TestName: "PreparedStatement_CreateWG",
			Status:   "FAIL",
			Error:    fmt.Sprintf("work group setup failed: %v", err),
		})
	}
	defer tc.deleteWorkGroup(psWorkGroup)

	// Prepared statement names must match ^[a-zA-Z_][a-zA-Z0-9_@:]{0,255}$
	// — no hyphens — so the unique suffix is joined with an underscore.
	psName := fmt.Sprintf("ps_%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "CreatePreparedStatement", func() error {
		if err := tc.createPreparedStatement(psWorkGroup, psName, "SELECT * FROM users WHERE id = ?"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "CreatePreparedStatement_Duplicate", func() error {
		_, err := tc.client.CreatePreparedStatement(tc.ctx, &athena.CreatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT 1"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate prepared statement")
		}
		// Per the Smithy model, CreatePreparedStatement declares only
		// InternalServerException and InvalidRequestException.
		if err := AssertErrorContains(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("expected InvalidRequestException for duplicate prepared statement, got: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetPreparedStatement", func() error {
		resp, err := tc.client.GetPreparedStatement(tc.ctx, &athena.GetPreparedStatementInput{
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

	results = append(results, tc.runner.RunTest("athena", "ListPreparedStatements", func() error {
		resp, err := tc.client.ListPreparedStatements(tc.ctx, &athena.ListPreparedStatementsInput{
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

	// Update then verify: the new query statement must be observable via
	// GetPreparedStatement.
	results = append(results, tc.runner.RunTest("athena", "UpdatePreparedStatement", func() error {
		_, err := tc.client.UpdatePreparedStatement(tc.ctx, &athena.UpdatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT * FROM orders WHERE id = ?"),
			Description:    aws.String("Updated prepared statement"),
		})
		if err != nil {
			return err
		}

		verifyResp, err := tc.client.GetPreparedStatement(tc.ctx, &athena.GetPreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err != nil {
			return err
		}
		if aws.ToString(verifyResp.PreparedStatement.QueryStatement) != "SELECT * FROM orders WHERE id = ?" {
			return fmt.Errorf("expected updated query statement, got %q", aws.ToString(verifyResp.PreparedStatement.QueryStatement))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "BatchGetPreparedStatement", func() error {
		psName2 := fmt.Sprintf("ps2_%d", time.Now().UnixNano()%1000000000)
		if err := tc.createPreparedStatement(psWorkGroup, psName2, "SELECT 1"); err != nil {
			return fmt.Errorf("second statement setup failed: %w", err)
		}

		resp, err := tc.client.BatchGetPreparedStatement(tc.ctx, &athena.BatchGetPreparedStatementInput{
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

		_, err = tc.client.DeletePreparedStatement(tc.ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName2),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err != nil {
			return fmt.Errorf("cleanup of second statement failed: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeletePreparedStatement", func() error {
		_, err := tc.client.DeletePreparedStatement(tc.ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "GetPreparedStatement_NonExistent", func() error {
		_, err := tc.client.GetPreparedStatement(tc.ctx, &athena.GetPreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeletePreparedStatement_NonExistent", func() error {
		_, err := tc.client.DeletePreparedStatement(tc.ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
