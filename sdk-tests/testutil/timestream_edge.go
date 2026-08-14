package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
)

func (r *TestRunner) runTimestreamEdgeTests(tc *tsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("timestream", "WriteRecords_QueryRoundtrip", func() error {
		rtDBName := tc.uniqueName("rt-db")
		rtTableName := tc.uniqueName("rt-table")
		_, err := tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(rtDBName),
		})
		if err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		defer tc.deleteDatabase(rtDBName)
		defer tc.deleteTable(rtDBName, rtTableName)

		_, err = tc.writeClient.CreateTable(tc.ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(rtDBName),
			TableName:    aws.String(rtTableName),
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}

		measureValue := fmt.Sprintf("verify-%d", time.Now().UnixNano())
		_, err = tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(rtDBName),
			TableName:    aws.String(rtTableName),
			Records: []types.Record{
				{
					MeasureName:      aws.String("cpu_utilization"),
					MeasureValue:     aws.String(measureValue),
					MeasureValueType: types.MeasureValueTypeDouble,
					Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
					TimeUnit:         types.TimeUnitMilliseconds,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("write: %v", err)
		}

		queryResp, err := tc.queryClient.Query(tc.ctx, &timestreamquery.QueryInput{
			QueryString: aws.String(fmt.Sprintf(`SELECT * FROM "%s"."%s"`, rtDBName, rtTableName)),
		})
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}
		if queryResp.QueryId == nil || *queryResp.QueryId == "" {
			return fmt.Errorf("QueryId is nil or empty")
		}
		if len(queryResp.ColumnInfo) == 0 {
			return fmt.Errorf("ColumnInfo is empty")
		}
		if len(queryResp.Rows) == 0 {
			return fmt.Errorf("query returned zero rows, expected at least 1")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListDatabases_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgDatabases []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagDB-%s-%d", pgTs, i)
			_, err := tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
				DatabaseName: aws.String(name),
			})
			if err != nil {
				for _, dn := range pgDatabases {
					tc.deleteDatabase(dn)
				}
				return fmt.Errorf("create database %s: %v", name, err)
			}
			pgDatabases = append(pgDatabases, name)
		}

		var allDatabases []string
		var nextToken *string
		for {
			resp, err := tc.writeClient.ListDatabases(tc.ctx, &timestreamwrite.ListDatabasesInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, dn := range pgDatabases {
					tc.deleteDatabase(dn)
				}
				return fmt.Errorf("list databases page: %v", err)
			}
			for _, d := range resp.Databases {
				if d.DatabaseName != nil && strings.Contains(*d.DatabaseName, "PagDB-"+pgTs) {
					allDatabases = append(allDatabases, *d.DatabaseName)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, dn := range pgDatabases {
			tc.deleteDatabase(dn)
		}
		if len(allDatabases) != 5 {
			return fmt.Errorf("expected 5 paginated databases, got %d", len(allDatabases))
		}
		return nil
	}))

	// MeasureValueType, TimeUnit and DimensionValueType are Smithy enums;
	// values outside the enum must be rejected with ValidationException.
	results = append(results, r.RunTest("timestream", "WriteRecords_InvalidEnums_Rejected", func() error {
		ieDBName := tc.uniqueName("ie-db")
		ieTableName := tc.uniqueName("ie-table")
		if err := tc.createDatabase(ieDBName); err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		defer tc.deleteDatabase(ieDBName)
		defer tc.deleteTable(ieDBName, ieTableName)
		if err := tc.createTable(ieDBName, ieTableName); err != nil {
			return fmt.Errorf("create table: %v", err)
		}

		_, err := tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(ieDBName),
			TableName:    aws.String(ieTableName),
			Records: []types.Record{{
				MeasureName:      aws.String("cpu"),
				MeasureValue:     aws.String("1.5"),
				MeasureValueType: types.MeasureValueType("DECIMAL"),
				Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
			}},
		})
		if !strings.Contains(fmt.Sprint(err), "ValidationException") {
			return fmt.Errorf("expected ValidationException for invalid MeasureValueType, got: %v", err)
		}

		_, err = tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(ieDBName),
			TableName:    aws.String(ieTableName),
			Records: []types.Record{{
				MeasureName:      aws.String("cpu"),
				MeasureValue:     aws.String("1.5"),
				MeasureValueType: types.MeasureValueTypeDouble,
				Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
				TimeUnit:         types.TimeUnit("CENTURIES"),
			}},
		})
		if !strings.Contains(fmt.Sprint(err), "ValidationException") {
			return fmt.Errorf("expected ValidationException for invalid TimeUnit, got: %v", err)
		}

		_, err = tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(ieDBName),
			TableName:    aws.String(ieTableName),
			Records: []types.Record{{
				MeasureName:      aws.String("cpu"),
				MeasureValue:     aws.String("1.5"),
				MeasureValueType: types.MeasureValueTypeDouble,
				Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
				Dimensions: []types.Dimension{{
					Name:               aws.String("host"),
					Value:              aws.String("a1"),
					DimensionValueType: types.DimensionValueType("INT"),
				}},
			}},
		})
		if !strings.Contains(fmt.Sprint(err), "ValidationException") {
			return fmt.Errorf("expected ValidationException for invalid DimensionValueType, got: %v", err)
		}

		// MeasureValues is only allowed when MeasureValueType is MULTI
		// (Smithy Record.MeasureValues documentation); the reverse
		// pairing (MULTI without MeasureValues) is also rejected.
		_, err = tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(ieDBName),
			TableName:    aws.String(ieTableName),
			Records: []types.Record{{
				MeasureName:      aws.String("cpu"),
				MeasureValue:     aws.String("1.5"),
				MeasureValueType: types.MeasureValueTypeDouble,
				Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
				MeasureValues: []types.MeasureValue{{
					Name:  aws.String("core"),
					Value: aws.String("3"),
					Type:  types.MeasureValueTypeBigint,
				}},
			}},
		})
		if !strings.Contains(fmt.Sprint(err), "ValidationException") {
			return fmt.Errorf("expected ValidationException for MeasureValues with scalar MeasureValueType, got: %v", err)
		}

		_, err = tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(ieDBName),
			TableName:    aws.String(ieTableName),
			Records: []types.Record{{
				MeasureName:      aws.String("cpu"),
				MeasureValueType: types.MeasureValueTypeMulti,
				Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
			}},
		})
		if !strings.Contains(fmt.Sprint(err), "ValidationException") {
			return fmt.Errorf("expected ValidationException for MULTI without MeasureValues, got: %v", err)
		}
		return nil
	}))

	// DataModelS3Configuration follows the Smithy S3BucketName and
	// S3ObjectKey shapes; invalid values must be rejected.
	results = append(results, r.RunTest("timestream", "CreateBatchLoadTask_InvalidS3Config_Rejected", func() error {
		_, err := tc.writeClient.CreateBatchLoadTask(tc.ctx, &timestreamwrite.CreateBatchLoadTaskInput{
			ClientToken:        aws.String(tc.uniqueName("bl")),
			TargetDatabaseName: aws.String("dummy-db"),
			TargetTableName:    aws.String("dummy-table"),
			DataSourceConfiguration: &types.DataSourceConfiguration{
				DataFormat: types.BatchLoadDataFormatCsv,
				DataSourceS3Configuration: &types.DataSourceS3Configuration{
					BucketName: aws.String("INVALID_BUCKET"),
				},
			},
			DataModelConfiguration: &types.DataModelConfiguration{
				DataModelS3Configuration: &types.DataModelS3Configuration{
					BucketName: aws.String("also_bad"),
					ObjectKey:  aws.String("model.json"),
				},
			},
			ReportConfiguration: &types.ReportConfiguration{
				ReportS3Configuration: &types.ReportS3Configuration{
					BucketName: aws.String("INVALID_BUCKET"),
				},
			},
		})
		if !strings.Contains(fmt.Sprint(err), "ValidationException") {
			return fmt.Errorf("expected ValidationException for invalid S3 bucket names, got: %v", err)
		}
		return nil
	}))

	return results
}
