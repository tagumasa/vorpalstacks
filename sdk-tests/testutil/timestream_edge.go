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

	return results
}
