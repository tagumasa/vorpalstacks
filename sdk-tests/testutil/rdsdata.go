package testutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"vorpalstacks-sdk-tests/config"
)

type rdsDataTestContext struct {
	client      *rdsdata.Client
	ctx         context.Context
	resourceArn string
	secretArn   string
	database    string
	region      string
	accountID   string
}

func (r *TestRunner) initRDSData() (*rdsDataTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	secretArn := fmt.Sprintf("arn:aws:secretsmanager:%s:%s:secret:test", r.region, r.accountID)
	// The RDS Data service validates secretArn against Secrets Manager, so
	// the test must seed a real secret carrying username/password fields
	// before any Data API call references it. CreateSecret is idempotent:
	// the surrounding dropTestDB logic likewise tolerates prior state.
	smClient := secretsmanager.NewFromConfig(cfg)
	if _, err := smClient.CreateSecret(context.Background(), &secretsmanager.CreateSecretInput{
		Name:         aws.String("test"),
		SecretString: aws.String(`{"username":"admin","password":"admin","engine":"mysql","host":"localhost","port":3306}`),
	}); err != nil {
		// If the secret already exists, that is fine — the prior run
		// seeded it. Any other error is fatal because subsequent Data API
		// calls would all fail with SecretsErrorException.
		if !strings.Contains(err.Error(), "already exists") &&
			!strings.Contains(err.Error(), "ResourceExistsException") {
			return nil, fmt.Errorf("failed to seed test secret: %w", err)
		}
	}

	return &rdsDataTestContext{
		client:      rdsdata.NewFromConfig(cfg),
		ctx:         context.Background(),
		resourceArn: fmt.Sprintf("arn:aws:rds:%s:%s:db:test-instance", r.region, r.accountID),
		secretArn:   secretArn,
		database:    "rdsdata_sdk_test",
		region:      r.region,
		accountID:   r.accountID,
	}, nil
}

func (tc *rdsDataTestContext) dropTestDB() error {
	tc.exec("DROP TABLE IF EXISTS sdk_test")
	_, err := tc.execAny("DROP DATABASE IF EXISTS " + tc.database)
	return err
}

// execInput builds an ExecuteStatement request against the suite's default
// database; the exec helpers below send it and the tests that need extra
// members set them on the returned input.
func (tc *rdsDataTestContext) execInput(sql string) *rdsdata.ExecuteStatementInput {
	return &rdsdata.ExecuteStatementInput{
		ResourceArn: &tc.resourceArn,
		SecretArn:   &tc.secretArn,
		Sql:         aws.String(sql),
		Database:    &tc.database,
	}
}

// exec runs one statement against the suite's default database and returns
// the raw result so field-level assertions can inspect it.
func (tc *rdsDataTestContext) exec(sql string) (*rdsdata.ExecuteStatementOutput, error) {
	return tc.client.ExecuteStatement(tc.ctx, tc.execInput(sql))
}

// execAny runs one statement without a database qualifier, for statements
// that must execute outside any database such as CREATE DATABASE and
// DROP DATABASE.
func (tc *rdsDataTestContext) execAny(sql string) (*rdsdata.ExecuteStatementOutput, error) {
	in := tc.execInput(sql)
	in.Database = nil
	return tc.client.ExecuteStatement(tc.ctx, in)
}

// execTx runs one statement inside an open transaction.
func (tc *rdsDataTestContext) execTx(sql string, txID *string) (*rdsdata.ExecuteStatementOutput, error) {
	in := tc.execInput(sql)
	in.TransactionId = txID
	return tc.client.ExecuteStatement(tc.ctx, in)
}

// beginTx opens a transaction on the suite's default database.
func (tc *rdsDataTestContext) beginTx() (*rdsdata.BeginTransactionOutput, error) {
	return tc.client.BeginTransaction(tc.ctx, &rdsdata.BeginTransactionInput{
		ResourceArn: &tc.resourceArn,
		SecretArn:   &tc.secretArn,
		Database:    &tc.database,
	})
}

// commitTx commits an open transaction.
func (tc *rdsDataTestContext) commitTx(id *string) (*rdsdata.CommitTransactionOutput, error) {
	return tc.client.CommitTransaction(tc.ctx, &rdsdata.CommitTransactionInput{
		ResourceArn:   &tc.resourceArn,
		SecretArn:     &tc.secretArn,
		TransactionId: id,
	})
}

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

	// Idempotent cleanup: drop stale data from prior runs so that
	// CREATE DATABASE / CREATE TABLE never collide with leftovers.
	_ = tc.dropTestDB()

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_CreateDatabase", func() error {
		_, err := tc.execAny("CREATE DATABASE " + tc.database)
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_CreateTable", func() error {
		_, err := tc.exec("CREATE TABLE sdk_test (id INT PRIMARY KEY, name VARCHAR(50), value DOUBLE)")
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Insert", func() error {
		_, err := tc.exec("INSERT INTO sdk_test VALUES (1, 'alice', 3.14), (2, 'bob', 2.72)")
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Select", func() error {
		resp, err := tc.exec("SELECT id, name, value FROM sdk_test ORDER BY id")
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
		in := tc.execInput("SELECT id, name FROM sdk_test LIMIT 1")
		in.IncludeResultMetadata = true
		resp, err := tc.client.ExecuteStatement(tc.ctx, in)
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
		in := tc.execInput("SELECT id, name FROM sdk_test WHERE id = 1")
		in.FormatRecordsAs = types.RecordsFormatTypeJson
		resp, err := tc.client.ExecuteStatement(tc.ctx, in)
		if err != nil {
			return err
		}
		if resp.FormattedRecords == nil || !strings.Contains(*resp.FormattedRecords, "alice") {
			return fmt.Errorf("expected formatted records containing 'alice', got %v", resp.FormattedRecords)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Update", func() error {
		_, err := tc.exec("UPDATE sdk_test SET name = 'updated' WHERE id = 1")
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "ExecuteStatement_Delete", func() error {
		_, err := tc.exec("DELETE FROM sdk_test WHERE id = 2")
		return err
	}))

	results = append(results, r.RunTest("rdsdata", "BatchExecuteStatement", func() error {
		_, err := tc.client.BatchExecuteStatement(tc.ctx, &rdsdata.BatchExecuteStatementInput{
			ResourceArn: &tc.resourceArn,
			SecretArn:   &tc.secretArn,
			Sql:         aws.String("INSERT INTO sdk_test VALUES (:id, :name, :val)"),
			Database:    &tc.database,
			ParameterSets: [][]types.SqlParameter{
				{
					{Name: aws.String("id"), Value: &types.FieldMemberLongValue{Value: 10}},
					{Name: aws.String("name"), Value: &types.FieldMemberStringValue{Value: "batch1"}},
					{Name: aws.String("val"), Value: &types.FieldMemberDoubleValue{Value: 1.1}},
				},
				{
					{Name: aws.String("id"), Value: &types.FieldMemberLongValue{Value: 11}},
					{Name: aws.String("name"), Value: &types.FieldMemberStringValue{Value: "batch2"}},
					{Name: aws.String("val"), Value: &types.FieldMemberDoubleValue{Value: 2.2}},
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
			SqlStatements:          aws.String("SELECT COUNT(*) as cnt FROM sdk_test; SELECT 42 as val"),
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
		beginResp, err := tc.beginTx()
		if err != nil {
			return fmt.Errorf("BeginTransaction failed: %w", err)
		}
		if beginResp.TransactionId == nil || *beginResp.TransactionId == "" {
			return fmt.Errorf("expected non-empty transactionId")
		}

		if _, err := tc.execTx("INSERT INTO sdk_test VALUES (100, 'tx_commit', 0.0)", beginResp.TransactionId); err != nil {
			return fmt.Errorf("ExecuteStatement with transactionId failed: %w", err)
		}

		commitResp, err := tc.commitTx(beginResp.TransactionId)
		if err != nil {
			return fmt.Errorf("CommitTransaction failed: %w", err)
		}
		if commitResp.TransactionStatus == nil || *commitResp.TransactionStatus != "COMMIT" {
			return fmt.Errorf("expected COMMIT status, got %v", commitResp.TransactionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "BeginTransaction_RollbackTransaction", func() error {
		beginResp, err := tc.beginTx()
		if err != nil {
			return fmt.Errorf("BeginTransaction failed: %w", err)
		}

		rbResp, err := tc.client.RollbackTransaction(tc.ctx, &rdsdata.RollbackTransactionInput{
			ResourceArn:   &tc.resourceArn,
			SecretArn:     &tc.secretArn,
			TransactionId: beginResp.TransactionId,
		})
		if err != nil {
			return fmt.Errorf("RollbackTransaction failed: %w", err)
		}
		if rbResp.TransactionStatus == nil || *rbResp.TransactionStatus != "ROLLBACK" {
			return fmt.Errorf("expected ROLLBACK status, got %v", rbResp.TransactionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "ErrorCases", func() error {
		for _, c := range []struct {
			name string
			call func() error
			code string
		}{
			{
				name: "invalid SQL is rejected as DatabaseErrorException",
				call: func() error {
					_, err := tc.exec("INVALID SQL SYNTAX HERE")
					return err
				},
				code: "DatabaseErrorException",
			},
			{
				name: "db-type resourceArn with no engine surfaces DatabaseUnavailableException",
				call: func() error {
					arn := fmt.Sprintf("arn:aws:rds:%s:%s:db:nonexistent-instance", tc.region, tc.accountID)
					_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
						ResourceArn: aws.String(arn),
						SecretArn:   &tc.secretArn,
						Sql:         aws.String("SELECT 1"),
					})
					return err
				},
				// Instance-level ARNs resolve straight to the engine lookup,
				// and an absent engine is deliberately surfaced as
				// DatabaseUnavailableException rather than a not-found fault;
				// the NotFoundException mapping applies to cluster ARNs.
				code: "DatabaseUnavailableException",
			},
		} {
			if err := expectAWSErrorCode(c.call(), c.code); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "Error_MissingResourceArn", func() error {
		// A nil resourceArn never reaches the server: the SDK's generated
		// OperationInputValidation middleware rejects the required member
		// client-side with smithy.InvalidParamsError, which is not a
		// smithy.APIError and therefore carries no AWS error code to pin.
		_, err := tc.client.ExecuteStatement(tc.ctx, &rdsdata.ExecuteStatementInput{
			SecretArn: &tc.secretArn,
			Sql:       aws.String("SELECT 1"),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing resourceArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("rdsdata", "Error_ReuseTransactionId", func() error {
		beginResp, err := tc.beginTx()
		if err != nil {
			return err
		}
		if _, err := tc.commitTx(beginResp.TransactionId); err != nil {
			return err
		}
		_, reuseErr := tc.execTx("SELECT 1", beginResp.TransactionId)
		if reuseErr == nil {
			return fmt.Errorf("expected error when reusing committed transactionId")
		}
		return nil
	}))

	return results
}
