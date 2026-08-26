package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestContext) testSQL() []TestResult {
	var results []TestResult

	suffix := time.Now().UnixNano() % 1000000000

	// --- Setup: CREATE TABLE users ---
	results = append(results, tc.runner.RunTest("athena", "SQL_CreateTable_Users", func() error {
		_, err := tc.startAndWaitForQuery(fmt.Sprintf("CREATE TABLE IF NOT EXISTS users_%d (id INT, name STRING, dept_id INT)", suffix))
		return err
	}))

	// --- Setup: CREATE TABLE departments ---
	results = append(results, tc.runner.RunTest("athena", "SQL_CreateTable_Departments", func() error {
		_, err := tc.startAndWaitForQuery(fmt.Sprintf("CREATE TABLE IF NOT EXISTS depts_%d (id INT, name STRING)", suffix))
		return err
	}))

	// --- Setup: INSERT data into users ---
	results = append(results, tc.runner.RunTest("athena", "SQL_Insert_Users", func() error {
		_, err := tc.startAndWaitForQuery(fmt.Sprintf(
			"INSERT INTO users_%d VALUES (1, 'alice', 10), (2, 'bob', 10), (3, 'charlie', 20), (4, 'dave', NULL)",
			suffix))
		return err
	}))

	// --- Setup: INSERT data into departments ---
	results = append(results, tc.runner.RunTest("athena", "SQL_Insert_Departments", func() error {
		_, err := tc.startAndWaitForQuery(fmt.Sprintf(
			"INSERT INTO depts_%d VALUES (10, 'engineering'), (20, 'sales'), (30, 'marketing')",
			suffix))
		return err
	}))

	// --- INNER JOIN ---
	results = append(results, tc.runner.RunTest("athena", "SQL_InnerJoin", func() error {
		q := fmt.Sprintf(
			"SELECT u.name, d.name FROM users_%d u JOIN depts_%d d ON u.dept_id = d.id",
			suffix, suffix)
		queryId, err := tc.startAndWaitForQuery(q)
		if err != nil {
			return err
		}

		// alice+eng, bob+eng, charlie+sales = 3 matched rows (dave has NULL dept)
		resultsResp, err := tc.client.GetQueryResults(tc.ctx, &athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(queryId),
		})
		if err != nil {
			return err
		}
		// Subtract header row
		dataRows := resultsResp.ResultSet.Rows
		if len(dataRows) > 0 && len(dataRows[0].Data) > 0 && dataRows[0].Data[0].VarCharValue != nil {
			if aws.ToString(dataRows[0].Data[0].VarCharValue) == "u.name" || aws.ToString(dataRows[0].Data[0].VarCharValue) == "name" {
				dataRows = dataRows[1:]
			}
		}
		if len(dataRows) < 3 {
			return fmt.Errorf("INNER JOIN expected at least 3 matched rows, got %d", len(dataRows))
		}
		return nil
	}))

	// --- LEFT JOIN ---
	results = append(results, tc.runner.RunTest("athena", "SQL_LeftJoin", func() error {
		q := fmt.Sprintf(
			"SELECT u.name FROM users_%d u LEFT JOIN depts_%d d ON u.dept_id = d.id",
			suffix, suffix)
		queryId, err := tc.startAndWaitForQuery(q)
		if err != nil {
			return err
		}

		// LEFT JOIN returns all 4 users (dave has no matching dept, but still appears)
		resultsResp, err := tc.client.GetQueryResults(tc.ctx, &athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(queryId),
		})
		if err != nil {
			return err
		}
		dataRows := resultsResp.ResultSet.Rows
		if len(dataRows) > 0 && len(dataRows[0].Data) > 0 && dataRows[0].Data[0].VarCharValue != nil {
			firstVal := aws.ToString(dataRows[0].Data[0].VarCharValue)
			if firstVal == "u.name" || firstVal == "name" {
				dataRows = dataRows[1:]
			}
		}
		if len(dataRows) < 4 {
			return fmt.Errorf("LEFT JOIN expected at least 4 rows (all users), got %d", len(dataRows))
		}
		return nil
	}))

	// --- FULL OUTER JOIN (was unparseable before fix) ---
	results = append(results, tc.runner.RunTest("athena", "SQL_FullOuterJoin", func() error {
		q := fmt.Sprintf(
			"SELECT u.name, d.name FROM users_%d u FULL OUTER JOIN depts_%d d ON u.dept_id = d.id",
			suffix, suffix)
		queryId, err := tc.startAndWaitForQuery(q)
		if err != nil {
			return fmt.Errorf("FULL OUTER JOIN should parse and execute: %v", err)
		}

		// FULL OUTER JOIN returns:
		// - 3 matched (alice+eng, bob+eng, charlie+sales)
		// - 1 unmatched left (dave with NULL dept)
		// - 1 unmatched right (marketing with NULL user)
		// = 5 rows total
		resultsResp, err := tc.client.GetQueryResults(tc.ctx, &athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(queryId),
		})
		if err != nil {
			return err
		}
		dataRows := resultsResp.ResultSet.Rows
		if len(dataRows) > 0 && len(dataRows[0].Data) > 0 && dataRows[0].Data[0].VarCharValue != nil {
			firstVal := aws.ToString(dataRows[0].Data[0].VarCharValue)
			if firstVal == "u.name" || firstVal == "name" {
				dataRows = dataRows[1:]
			}
		}
		// Should have more rows than INNER JOIN (3) because of unmatched rows
		if len(dataRows) <= 3 {
			return fmt.Errorf("FULL OUTER JOIN expected >3 rows (matched + unmatched), got %d — may be falling back to INNER JOIN", len(dataRows))
		}
		return nil
	}))

	// --- WHERE with NOT expression (was silent fail-OPEN before fix) ---
	results = append(results, tc.runner.RunTest("athena", "SQL_WhereNotExpr", func() error {
		q := fmt.Sprintf("SELECT name FROM users_%d WHERE NOT (id = 1)", suffix)
		queryId, err := tc.startAndWaitForQuery(q)
		if err != nil {
			return err
		}

		// NOT (id=1) should return 3 rows (bob, charlie, dave)
		// Before the fix, evaluateWhere returned true for NotExpr → all 4 rows
		resultsResp, err := tc.client.GetQueryResults(tc.ctx, &athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(queryId),
		})
		if err != nil {
			return err
		}
		dataRows := resultsResp.ResultSet.Rows
		if len(dataRows) > 0 && len(dataRows[0].Data) > 0 && dataRows[0].Data[0].VarCharValue != nil {
			if aws.ToString(dataRows[0].Data[0].VarCharValue) == "name" {
				dataRows = dataRows[1:]
			}
		}
		if len(dataRows) != 3 {
			return fmt.Errorf("WHERE NOT (id=1) expected exactly 3 rows, got %d (fail-OPEN bug?)", len(dataRows))
		}
		// Verify alice is NOT in results
		for _, row := range dataRows {
			if len(row.Data) > 0 && row.Data[0].VarCharValue != nil {
				if aws.ToString(row.Data[0].VarCharValue) == "alice" {
					return fmt.Errorf("alice should not appear in WHERE NOT (id=1) results")
				}
			}
		}
		return nil
	}))

	// --- DROP TABLE (tests error propagation for data deletion) ---
	results = append(results, tc.runner.RunTest("athena", "SQL_DropTable", func() error {
		// Drop both tables
		for _, table := range []string{fmt.Sprintf("users_%d", suffix), fmt.Sprintf("depts_%d", suffix)} {
			if _, err := tc.startAndWaitForQuery(fmt.Sprintf("DROP TABLE %s", table)); err != nil {
				return err
			}
		}
		return nil
	}))

	return results
}
