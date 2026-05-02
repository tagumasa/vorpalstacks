package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestCtx) testPreparedStatements() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	psWorkGroup := fmt.Sprintf("ps-wg-%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "PreparedStatement_CreateWG", func() error {
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(psWorkGroup),
		})
		return err
	}))

	psName := fmt.Sprintf("ps-%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "CreatePreparedStatement", func() error {
		_, err := client.CreatePreparedStatement(ctx, &athena.CreatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT * FROM users WHERE id = ?"),
			Description:    aws.String("Test prepared statement"),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "CreatePreparedStatement_Duplicate", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "GetPreparedStatement", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "ListPreparedStatements", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "UpdatePreparedStatement", func() error {
		_, err := client.UpdatePreparedStatement(ctx, &athena.UpdatePreparedStatementInput{
			StatementName:  aws.String(psName),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT * FROM orders WHERE id = ?"),
			Description:    aws.String("Updated prepared statement"),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "GetPreparedStatement_AfterUpdate", func() error {
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
	results = append(results, tc.runner.RunTest("athena", "CreatePreparedStatement_Second", func() error {
		_, err := client.CreatePreparedStatement(ctx, &athena.CreatePreparedStatementInput{
			StatementName:  aws.String(psName2),
			WorkGroup:      aws.String(psWorkGroup),
			QueryStatement: aws.String("SELECT 1"),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "BatchGetPreparedStatement", func() error {
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

	results = append(results, tc.runner.RunTest("athena", "DeletePreparedStatement", func() error {
		_, err := client.DeletePreparedStatement(ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		return err
	}))

	results = append(results, tc.runner.RunTest("athena", "GetPreparedStatement_NonExistent", func() error {
		_, err := client.GetPreparedStatement(ctx, &athena.GetPreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err := AssertErrorContains(err, "InvalidRequestException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeletePreparedStatement_NonExistent", func() error {
		_, err := client.DeletePreparedStatement(ctx, &athena.DeletePreparedStatementInput{
			StatementName: aws.String(psName),
			WorkGroup:     aws.String(psWorkGroup),
		})
		if err := AssertErrorContains(err, "InvalidRequestException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteWorkGroup_PSCleanup", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String(psWorkGroup),
		})
		return err
	}))

	return results
}
