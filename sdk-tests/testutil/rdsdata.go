package testutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"vorpalstacks-sdk-tests/config"
)

type rdsDataTestContext struct {
	client      *rdsdata.Client
	ctx         context.Context
	resourceArn string
	secretArn   string
	database    string
}

func (r *TestRunner) initRDSData() (*rdsDataTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &rdsDataTestContext{
		client:      rdsdata.NewFromConfig(cfg),
		ctx:         context.Background(),
		resourceArn: "arn:aws:rds:us-east-1:000000000000:db:test-instance",
		secretArn:   "arn:aws:secretsmanager:us-east-1:000000000000:secret:test",
		database:    "rdsdata_sdk_test",
	}, nil
}

func ptrStr(v string) *string { return &v }

func (r *TestRunner) RunRDSDataTests() []TestResult {
	tc, err := r.initRDSData()
	if err != nil {
		return []TestResult{{
			Service:  "rdsdata",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	var results []TestResult

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_CreateDatabase", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("CREATE DATABASE " + tc.database),
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_CreateTable", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("CREATE TABLE sdk_test (id INT PRIMARY KEY, name VARCHAR(50), value DOUBLE)"),
			Database:    &tc.database,
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Insert", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("INSERT INTO sdk_test VALUES (1, 'alice', 3.14)"),
			Database:    &tc.database,
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_InsertSecondRow", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("INSERT INTO sdk_test VALUES (2, 'bob', 2.72)"),
			Database:    &tc.database,
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Select", func() error {
		resp, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("SELECT id, name, value FROM sdk_test ORDER BY id"),
			Database:    &tc.database,
		})
		if err != nil {
			return err
		}
		if len(resp.Records) != 2 {
			return fmt.Errorf("expected 2 records, got %d", len(resp.Records))
		}
		if len(resp.Records[0]) != 3 {
			return fmt.Errorf("expected 3 fields per record, got %d", len(resp.Records[0]))
		}
		f0, ok := resp.Records[0][0].(*types.FieldMemberLongValue)
		if !ok || f0.Value != 1 {
			return fmt.Errorf("expected first row id=1, got %T %v", resp.Records[0][0], resp.Records[0][0])
		}
		f1, ok := resp.Records[0][1].(*types.FieldMemberStringValue)
		if !ok || f1.Value != "alice" {
			return fmt.Errorf("expected first row name='alice', got %T %v", resp.Records[0][1], resp.Records[0][1])
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_IncludeResultMetadata", func() error {
		resp, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn:           &tc.resourceArn,
			SecretArn:             &tc.secretArn,
			Sql:                   ptrStr("SELECT id, name FROM sdk_test LIMIT 1"),
			Database:              &tc.database,
			IncludeResultMetadata: true,
		})
		if err != nil {
			return err
		}
		if len(resp.ColumnMetadata) != 2 {
			return fmt.Errorf("expected 2 columns, got %d", len(resp.ColumnMetadata))
		}
		if resp.ColumnMetadata[0].Name == nil || *resp.ColumnMetadata[0].Name != "id" {
			return fmt.Errorf("expected first column name='id', got %v", resp.ColumnMetadata[0].Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_FormatRecordsAsJSON", func() error {
		resp, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn:     &tc.resourceArn,
			SecretArn:       &tc.secretArn,
			Sql:             ptrStr("SELECT id, name FROM sdk_test WHERE id = 1"),
			Database:        &tc.database,
			FormatRecordsAs: types.RecordsFormatTypeJson,
		})
		if err != nil {
			return err
		}
		if resp.FormattedRecords == nil || !strings.Contains(*resp.FormattedRecords, "alice") {
			return fmt.Errorf("expected formatted records containing 'alice', got %v", resp.FormattedRecords)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Update", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("UPDATE sdk_test SET name = 'updated' WHERE id = 1"),
			Database:    &tc.database,
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Delete", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("DELETE FROM sdk_test WHERE id = 2"),
			Database:    &tc.database,
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "BatchExecuteStatement", func() error {
		_, err := tc.client.BatchExecuteStatement(tc.ctx, &rdsdata.BatchExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("INSERT INTO sdk_test VALUES (:id, :name, :val)"),
			Database:    &tc.database,
			ParameterSets: [][]types.SqlParameter{
				{
					{Name: ptrStr("id"), Value: &types.FieldMemberLongValue{Value: 10}},
					{Name: ptrStr("name"), Value: &types.FieldMemberStringValue{Value: "batch1"}},
					{Name: ptrStr("val"), Value: &types.FieldMemberDoubleValue{Value: 1.1}},
				},
				{
					{Name: ptrStr("id"), Value: &types.FieldMemberLongValue{Value: 11}},
					{Name: ptrStr("name"), Value: &types.FieldMemberStringValue{Value: "batch2"}},
					{Name: ptrStr("val"), Value: &types.FieldMemberDoubleValue{Value: 2.2}},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteSql_MultiStatement", func() error {
		resp, err := tc.client.ExecuteSql(tc.ctx, &rdsdata.ExecuteSqlInput{
			DbClusterOrInstanceArn: &tc.resourceArn,
			AwsSecretStoreArn:      &tc.secretArn,
			Database:               &tc.database,
			SqlStatements:          ptrStr("SELECT COUNT(*) as cnt FROM sdk_test; SELECT 42 as val"),
		})
		if err != nil {
			return err
		}
		if len(resp.SqlStatementResults) != 2 {
			return fmt.Errorf("expected 2 statement results, got %d", len(resp.SqlStatementResults))
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "BeginTransaction_CommitTransaction", func() error {
		beginResp, err := tc.client.BeginTransaction(tc.ctx, &rdsdata.BeginTransactionInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Database:    &tc.database,
		})
		if err != nil {
			return fmt.Errorf("BeginTransaction failed: %w", err)
		}
		if beginResp.TransactionId == nil || *beginResp.TransactionId == "" {
			return fmt.Errorf("expected non-empty transactionId")
		}

		_, execErr := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn:   &tc.resourceArn,
			SecretArn:     &tc.secretArn,
			Sql:           ptrStr("INSERT INTO sdk_test VALUES (100, 'tx_commit', 0.0)"),
			Database:      &tc.database,
			TransactionId: beginResp.TransactionId,
		})
		if execErr != nil {
			return fmt.Errorf("ExecuteStatement with transactionId failed: %w", execErr)
		}

		commitResp, commitErr := tc.client.CommitTransaction(tc.ctx, &rdsdata.CommitTransactionInput{
			ResourceArn:   &tc.resourceArn,
			SecretArn:     &tc.secretArn,
			TransactionId: beginResp.TransactionId,
		})
		if commitErr != nil {
			return fmt.Errorf("CommitTransaction failed: %w", commitErr)
		}
		if commitResp.TransactionStatus == nil || *commitResp.TransactionStatus != "COMMIT" {
			return fmt.Errorf("expected COMMIT status, got %v", commitResp.TransactionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "BeginTransaction_RollbackTransaction", func() error {
		beginResp, err := tc.client.BeginTransaction(tc.ctx, &rdsdata.BeginTransactionInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Database:    &tc.database,
		})
		if err != nil {
			return fmt.Errorf("BeginTransaction failed: %w", err)
		}

		rbResp, rbErr := tc.client.RollbackTransaction(tc.ctx, &rdsdata.RollbackTransactionInput{
			ResourceArn:   &tc.resourceArn,
			SecretArn:     &tc.secretArn,
			TransactionId: beginResp.TransactionId,
		})
		if rbErr != nil {
			return fmt.Errorf("RollbackTransaction failed: %w", rbErr)
		}
		if rbResp.TransactionStatus == nil || *rbResp.TransactionStatus != "ROLLBACK" {
			return fmt.Errorf("expected ROLLBACK status, got %v", rbResp.TransactionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "Error_InvalidSQL", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("INVALID SQL SYNTAX HERE"),
			Database:    &tc.database,
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid SQL")
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "Error_MissingResourceArn", func() error {
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			SecretArn: &tc.secretArn,
			Sql:       ptrStr("SELECT 1"),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing resourceArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "Error_NonexistentInstance", func() error {
		arn := "arn:aws:rds:us-east-1:000000000000:db:nonexistent-instance"
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn: &arn,
			SecretArn:   &tc.secretArn,
			Sql:         ptrStr("SELECT 1"),
		})
		if err == nil {
			return fmt.Errorf("expected error for nonexistent instance")
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "Error_ReuseTransactionId", func() error {
		beginResp, err := tc.client.BeginTransaction(tc.ctx, &rdsdata.BeginTransactionInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Database:    &tc.database,
		})
		if err != nil {
			return err
		}
		_, commitErr := tc.client.CommitTransaction(tc.ctx, &rdsdata.CommitTransactionInput{
			ResourceArn:   &tc.resourceArn,
			SecretArn:     &tc.secretArn,
			TransactionId: beginResp.TransactionId,
		})
		if commitErr != nil {
			return commitErr
		}
		_, reuseErr := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			ResourceArn:   &tc.resourceArn,
			SecretArn:     &tc.secretArn,
			Sql:           ptrStr("SELECT 1"),
			Database:      &tc.database,
			TransactionId: beginResp.TransactionId,
		})
		if reuseErr == nil {
			return fmt.Errorf("expected error when reusing committed transactionId")
		}
		return nil
	}))

	return results
}
