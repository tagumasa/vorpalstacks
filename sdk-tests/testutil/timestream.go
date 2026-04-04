package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	tqtypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunTimestreamTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "timestream",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := timestreamwrite.NewFromConfig(cfg)
	queryClient := timestreamquery.NewFromConfig(cfg)
	ctx := context.Background()

	databaseName := fmt.Sprintf("test-db-%d", time.Now().UnixNano())
	tableName := fmt.Sprintf("test-table-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("timestream", "CreateDatabase", func() error {
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(databaseName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "ListDatabases", func() error {
		_, err := client.ListDatabases(ctx, &timestreamwrite.ListDatabasesInput{
			MaxResults: aws.Int32(10),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeDatabase", func() error {
		_, err := client.DescribeDatabase(ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String(databaseName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "CreateTable", func() error {
		_, err := client.CreateTable(ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "ListTables", func() error {
		_, err := client.ListTables(ctx, &timestreamwrite.ListTablesInput{
			DatabaseName: aws.String(databaseName),
			MaxResults:   aws.Int32(10),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeTable", func() error {
		_, err := client.DescribeTable(ctx, &timestreamwrite.DescribeTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "WriteRecords", func() error {
		_, err := client.WriteRecords(ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
			Records: []types.Record{
				{
					Dimensions: []types.Dimension{
						{Name: aws.String("dim1"), Value: aws.String("value1")},
					},
					MeasureName:      aws.String("metric1"),
					MeasureValue:     aws.String("100"),
					MeasureValueType: types.MeasureValueTypeDouble,
					Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
					TimeUnit:         types.TimeUnitMilliseconds,
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "UpdateTable", func() error {
		_, err := client.UpdateTable(ctx, &timestreamwrite.UpdateTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
			RetentionProperties: &types.RetentionProperties{
				MemoryStoreRetentionPeriodInHours:  aws.Int64(24),
				MagneticStoreRetentionPeriodInDays: aws.Int64(7),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeEndpoints_Write", func() error {
		_, err := client.DescribeEndpoints(ctx, &timestreamwrite.DescribeEndpointsInput{})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DeleteTable", func() error {
		_, err := client.DeleteTable(ctx, &timestreamwrite.DeleteTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "UpdateDatabase", func() error {
		_, err := client.UpdateDatabase(ctx, &timestreamwrite.UpdateDatabaseInput{
			DatabaseName: aws.String(databaseName),
			KmsKeyId:     aws.String("arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DeleteDatabase", func() error {
		_, err := client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(databaseName),
		})
		return err
	}))

	// === Write-side Tagging Tests ===

	tagDBName := fmt.Sprintf("tag-db-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("timestream", "CreateDatabase_WithTags", func() error {
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(tagDBName),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("dev")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeDatabase_Tags", func() error {
		_, err := client.DescribeDatabase(ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String(tagDBName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "TagResource_Database_Cleanup", func() error {
		_, err := client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(tagDBName),
		})
		return err
	}))

	// === Batch Load Task Tests ===

	blDBName := fmt.Sprintf("bl-db-%d", time.Now().UnixNano())
	blTableName := fmt.Sprintf("bl-table-%d", time.Now().UnixNano())
	blTaskID := ""

	results = append(results, r.RunTest("timestream", "BatchLoad_Setup", func() error {
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(blDBName),
		})
		if err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		_, err = client.CreateTable(ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(blDBName),
			TableName:    aws.String(blTableName),
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CreateBatchLoadTask", func() error {
		resp, err := client.CreateBatchLoadTask(ctx, &timestreamwrite.CreateBatchLoadTaskInput{
			TargetDatabaseName: aws.String(blDBName),
			TargetTableName:    aws.String(blTableName),
			DataSourceConfiguration: &types.DataSourceConfiguration{
				DataFormat: types.BatchLoadDataFormatCsv,
				DataSourceS3Configuration: &types.DataSourceS3Configuration{
					BucketName: aws.String("test-bucket"),
				},
			},
			ReportConfiguration: &types.ReportConfiguration{
				ReportS3Configuration: &types.ReportS3Configuration{
					BucketName: aws.String("report-bucket"),
				},
			},
			DataModelConfiguration: &types.DataModelConfiguration{
				DataModel: &types.DataModel{
					DimensionMappings: []types.DimensionMapping{
						{
							SourceColumn:      aws.String("col1"),
							DestinationColumn: aws.String("dim1"),
						},
					},
					MeasureNameColumn: aws.String("measure"),
					TimeColumn:        aws.String("ts"),
					TimeUnit:          types.TimeUnitMilliseconds,
				},
			},
		})
		if err != nil {
			return err
		}
		blTaskID = *resp.TaskId
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeBatchLoadTask", func() error {
		resp, err := client.DescribeBatchLoadTask(ctx, &timestreamwrite.DescribeBatchLoadTaskInput{
			TaskId: aws.String(blTaskID),
		})
		if err != nil {
			return err
		}
		if resp.BatchLoadTaskDescription == nil {
			return fmt.Errorf("expected BatchLoadTaskDescription")
		}
		if resp.BatchLoadTaskDescription.CreationTime == nil {
			return fmt.Errorf("expected CreationTime to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListBatchLoadTasks", func() error {
		resp, err := client.ListBatchLoadTasks(ctx, &timestreamwrite.ListBatchLoadTasksInput{})
		if err != nil {
			return err
		}
		if len(resp.BatchLoadTasks) == 0 {
			return fmt.Errorf("expected at least 1 task, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "BatchLoad_Cleanup", func() error {
		_, _ = client.DeleteTable(ctx, &timestreamwrite.DeleteTableInput{
			DatabaseName: aws.String(blDBName),
			TableName:    aws.String(blTableName),
		})
		_, _ = client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(blDBName),
		})
		return nil
	}))

	// === Query-side Tests ===

	results = append(results, r.RunTest("timestream", "DescribeEndpoints_Query", func() error {
		_, err := queryClient.DescribeEndpoints(ctx, &timestreamquery.DescribeEndpointsInput{})
		return err
	}))

	results = append(results, r.RunTest("timestream", "PrepareQuery", func() error {
		resp, err := queryClient.PrepareQuery(ctx, &timestreamquery.PrepareQueryInput{
			QueryString: aws.String("SELECT 1"),
		})
		if err != nil {
			return err
		}
		if resp.QueryString == nil || *resp.QueryString != "SELECT 1" {
			return fmt.Errorf("expected QueryString to be returned")
		}
		return nil
	}))

	// === Scheduled Query Tests ===

	sqName := fmt.Sprintf("sq-test-%d", time.Now().UnixNano())
	sqDBName := fmt.Sprintf("sq-db-%d", time.Now().UnixNano())
	sqTableName := fmt.Sprintf("sq-table-%d", time.Now().UnixNano())
	sqRoleName := fmt.Sprintf("ts-sq-role-%d", time.Now().UnixNano())
	var sqARN string

	createIAMRole := func(roleName string) error {
		return IAMCreateRole(r.endpoint, roleName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"timestream.amazonaws.com"},"Action":"sts:AssumeRole"}]}`)
	}

	deleteIAMRole := func(roleName string) {
		IAMDeleteRole(r.endpoint, roleName)
	}

	results = append(results, r.RunTest("timestream", "ScheduledQuery_Setup", func() error {
		if err := createIAMRole(sqRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(sqDBName),
		})
		if err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		_, err = client.CreateTable(ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(sqDBName),
			TableName:    aws.String(sqTableName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "CreateScheduledQuery", func() error {
		resp, err := queryClient.CreateScheduledQuery(ctx, &timestreamquery.CreateScheduledQueryInput{
			Name:        aws.String(sqName),
			QueryString: aws.String(fmt.Sprintf(`SELECT * FROM "%s"."%s"`, sqDBName, sqTableName)),
			ScheduleConfiguration: &tqtypes.ScheduleConfiguration{
				ScheduleExpression: aws.String("cron(0 0 * * ? *)"),
			},
			NotificationConfiguration: &tqtypes.NotificationConfiguration{
				SnsConfiguration: &tqtypes.SnsConfiguration{
					TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:test-topic"),
				},
			},
			ErrorReportConfiguration: &tqtypes.ErrorReportConfiguration{
				S3Configuration: &tqtypes.S3Configuration{
					BucketName: aws.String("error-report-bucket"),
				},
			},
			ScheduledQueryExecutionRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::000000000000:role/%s", sqRoleName)),
			Tags: []tqtypes.Tag{
				{Key: aws.String("env"), Value: aws.String("scheduled-test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("expected Arn to be set")
		}
		sqARN = *resp.Arn
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeScheduledQuery", func() error {
		resp, err := queryClient.DescribeScheduledQuery(ctx, &timestreamquery.DescribeScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		if resp.ScheduledQuery == nil {
			return fmt.Errorf("expected ScheduledQuery")
		}
		if resp.ScheduledQuery.Name == nil || *resp.ScheduledQuery.Name != sqName {
			return fmt.Errorf("expected Name=%s, got %v", sqName, resp.ScheduledQuery.Name)
		}
		if resp.ScheduledQuery.CreationTime == nil {
			return fmt.Errorf("expected CreationTime to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListScheduledQueries", func() error {
		resp, err := queryClient.ListScheduledQueries(ctx, &timestreamquery.ListScheduledQueriesInput{})
		if err != nil {
			return err
		}
		if len(resp.ScheduledQueries) == 0 {
			return fmt.Errorf("expected at least 1 scheduled query")
		}
		found := false
		for _, sq := range resp.ScheduledQueries {
			if sq.Name != nil && *sq.Name == sqName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("scheduled query %s not found in list", sqName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UpdateScheduledQuery", func() error {
		_, err := queryClient.UpdateScheduledQuery(ctx, &timestreamquery.UpdateScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
			State:             tqtypes.ScheduledQueryStateDisabled,
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "UpdateScheduledQuery_Verify", func() error {
		resp, err := queryClient.DescribeScheduledQuery(ctx, &timestreamquery.DescribeScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		if resp.ScheduledQuery.State != tqtypes.ScheduledQueryStateDisabled {
			return fmt.Errorf("expected state DISABLED, got %s", resp.ScheduledQuery.State)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ExecuteScheduledQuery", func() error {
		_, err := queryClient.ExecuteScheduledQuery(ctx, &timestreamquery.ExecuteScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
			InvocationTime:    aws.Time(time.Now().UTC()),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "TagResource_ScheduledQuery", func() error {
		_, err := queryClient.TagResource(ctx, &timestreamquery.TagResourceInput{
			ResourceARN: aws.String(sqARN),
			Tags: []tqtypes.Tag{
				{Key: aws.String("extra"), Value: aws.String("tag")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "ListTagsForResource_ScheduledQuery", func() error {
		resp, err := queryClient.ListTagsForResource(ctx, &timestreamquery.ListTagsForResourceInput{
			ResourceARN: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UntagResource_ScheduledQuery", func() error {
		_, err := queryClient.UntagResource(ctx, &timestreamquery.UntagResourceInput{
			ResourceARN: aws.String(sqARN),
			TagKeys:     []string{"extra"},
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DeleteScheduledQuery", func() error {
		_, err := queryClient.DeleteScheduledQuery(ctx, &timestreamquery.DeleteScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeScheduledQuery_NonExistent", func() error {
		fakeARN := "arn:aws:timestream:us-east-1:000000000000:scheduled-query/nonexistent"
		_, err := queryClient.DescribeScheduledQuery(ctx, &timestreamquery.DescribeScheduledQueryInput{
			ScheduledQueryArn: aws.String(fakeARN),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CreateScheduledQuery_Duplicate", func() error {
		dupSQName := fmt.Sprintf("dup-sq-%d", time.Now().UnixNano())
		_, err := queryClient.CreateScheduledQuery(ctx, &timestreamquery.CreateScheduledQueryInput{
			Name:        aws.String(dupSQName),
			QueryString: aws.String("SELECT 1"),
			ScheduleConfiguration: &tqtypes.ScheduleConfiguration{
				ScheduleExpression: aws.String("cron(0 0 * * ? *)"),
			},
			NotificationConfiguration: &tqtypes.NotificationConfiguration{
				SnsConfiguration: &tqtypes.SnsConfiguration{
					TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:test-topic"),
				},
			},
			ErrorReportConfiguration: &tqtypes.ErrorReportConfiguration{
				S3Configuration: &tqtypes.S3Configuration{
					BucketName: aws.String("error-report-bucket"),
				},
			},
			ScheduledQueryExecutionRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::000000000000:role/%s", sqRoleName)),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer queryClient.DeleteScheduledQuery(ctx, &timestreamquery.DeleteScheduledQueryInput{
			ScheduledQueryArn: aws.String(fmt.Sprintf("arn:aws:timestream:us-east-1:000000000000:scheduled-query/%s", dupSQName)),
		})
		_, err = queryClient.CreateScheduledQuery(ctx, &timestreamquery.CreateScheduledQueryInput{
			Name:        aws.String(dupSQName),
			QueryString: aws.String("SELECT 1"),
			ScheduleConfiguration: &tqtypes.ScheduleConfiguration{
				ScheduleExpression: aws.String("cron(0 0 * * ? *)"),
			},
			NotificationConfiguration: &tqtypes.NotificationConfiguration{
				SnsConfiguration: &tqtypes.SnsConfiguration{
					TopicArn: aws.String("arn:aws:sns:us-east-1:000000000000:test-topic"),
				},
			},
			ErrorReportConfiguration: &tqtypes.ErrorReportConfiguration{
				S3Configuration: &tqtypes.S3Configuration{
					BucketName: aws.String("error-report-bucket"),
				},
			},
			ScheduledQueryExecutionRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::000000000000:role/%s", sqRoleName)),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate scheduled query")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ScheduledQuery_Cleanup", func() error {
		_, _ = client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(sqDBName),
		})
		deleteIAMRole(sqRoleName)
		return nil
	}))

	// === Account Settings Tests ===

	results = append(results, r.RunTest("timestream", "DescribeAccountSettings", func() error {
		resp, err := queryClient.DescribeAccountSettings(ctx, &timestreamquery.DescribeAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.MaxQueryTCU == nil {
			return fmt.Errorf("expected MaxQueryTCU to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UpdateAccountSettings", func() error {
		_, err := queryClient.UpdateAccountSettings(ctx, &timestreamquery.UpdateAccountSettingsInput{
			MaxQueryTCU:       aws.Int32(8),
			QueryPricingModel: tqtypes.QueryPricingModelComputeUnits,
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeAccountSettings_AfterUpdate", func() error {
		resp, err := queryClient.DescribeAccountSettings(ctx, &timestreamquery.DescribeAccountSettingsInput{})
		if err != nil {
			return err
		}
		if resp.MaxQueryTCU == nil || *resp.MaxQueryTCU != 8 {
			return fmt.Errorf("expected MaxQueryTCU=8, got %v", resp.MaxQueryTCU)
		}
		return nil
	}))

	// === Error / Edge Case Tests ===

	results = append(results, r.RunTest("timestream", "DescribeDatabase_NonExistent", func() error {
		_, err := client.DescribeDatabase(ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String("nonexistent-db-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeTable_NonExistent", func() error {
		_, err := client.DescribeTable(ctx, &timestreamwrite.DescribeTableInput{
			DatabaseName: aws.String("nonexistent-db-xyz"),
			TableName:    aws.String("nonexistent-table-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CreateDatabase_Duplicate", func() error {
		dupDBName := fmt.Sprintf("dup-db-%d", time.Now().UnixNano())
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})

		_, err = client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate database")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "WriteRecords_GetRecords_Roundtrip", func() error {
		rtDBName := fmt.Sprintf("rt-db-%d", time.Now().UnixNano())
		rtTableName := fmt.Sprintf("rt-table-%d", time.Now().UnixNano())
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(rtDBName),
		})
		if err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		defer client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(rtDBName),
		})

		_, err = client.CreateTable(ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(rtDBName),
			TableName:    aws.String(rtTableName),
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}

		measureValue := fmt.Sprintf("verify-%d", time.Now().UnixNano())
		_, err = client.WriteRecords(ctx, &timestreamwrite.WriteRecordsInput{
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

		queryResp, err := queryClient.Query(ctx, &timestreamquery.QueryInput{
			QueryString: aws.String(fmt.Sprintf(`SELECT * FROM "%s"."%s"`, rtDBName, rtTableName)),
		})
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}
		if len(queryResp.Rows) == 0 {
			return fmt.Errorf("query returned zero rows, expected at least 1")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeBatchLoadTask_NonExistent", func() error {
		_, err := client.DescribeBatchLoadTask(ctx, &timestreamwrite.DescribeBatchLoadTaskInput{
			TaskId: aws.String("nonexistent-task-id"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("timestream", "ListDatabases_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgDatabases []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagDB-%s-%d", pgTs, i)
			_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
				DatabaseName: aws.String(name),
			})
			if err != nil {
				for _, dn := range pgDatabases {
					client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{DatabaseName: aws.String(dn)})
				}
				return fmt.Errorf("create database %s: %v", name, err)
			}
			pgDatabases = append(pgDatabases, name)
		}

		var allDatabases []string
		var nextToken *string
		for {
			resp, err := client.ListDatabases(ctx, &timestreamwrite.ListDatabasesInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, dn := range pgDatabases {
					client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{DatabaseName: aws.String(dn)})
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
			client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{DatabaseName: aws.String(dn)})
		}
		if len(allDatabases) != 5 {
			return fmt.Errorf("expected 5 paginated databases, got %d", len(allDatabases))
		}
		return nil
	}))

	return results
}
