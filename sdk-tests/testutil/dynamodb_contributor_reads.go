package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"vorpalstacks-sdk-tests/config"
)

// dynamoDBContributorReadPathsTests pins the contributor insights coverage
// of every item read path: Scan, Query, BatchGetItem and TransactGetItems
// all count towards the ConsumedThroughputUnits aggregation, and a Query
// contributes exactly one event on the partition-key series regardless of
// how many items it returns.
func (r *TestRunner) dynamoDBContributorReadPathsTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()
	tableName := fmt.Sprintf("contributor-reads-%d", suffix)

	setupErr := func(name string, err error) []TestResult {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: name,
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}

	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("sk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
			{AttributeName: aws.String("sk"), KeyType: dynamodbtypes.KeyTypeRange},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("create table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
	if err := waitKinesisDestTableActive(ctx, client, tableName); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("wait active: %v", err))
	}
	if _, err := client.UpdateContributorInsights(ctx, &dynamodb.UpdateContributorInsightsInput{
		TableName:                 aws.String(tableName),
		ContributorInsightsAction: dynamodbtypes.ContributorInsightsActionEnable,
	}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("enable insights: %v", err))
	}

	item := func(pk, sk string) map[string]dynamodbtypes.AttributeValue {
		return map[string]dynamodbtypes.AttributeValue{
			"pk": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"sk": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		}
	}

	windowStart := time.Now().Add(-time.Minute)

	// Writes: A partition receives two items, B partition one. A write
	// counts as three units on both tracked series of its key.
	for _, key := range []struct{ pk, sk string }{
		{"A", "a1"}, {"A", "a2"}, {"B", "b1"},
	} {
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(tableName), Item: item(key.pk, key.sk)}); err != nil {
			return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("put %s/%s: %v", key.pk, key.sk, err))
		}
	}

	// Scan reads all three items: PKC A += 2, B += 1. Both A items share
	// the partition counter, pinning the aggregation of same-key events
	// inside one recording transaction.
	if _, err := client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(tableName)}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("scan: %v", err))
	}
	// A Query is one event on the partition series: PKC A += 1.
	if _, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]dynamodbtypes.AttributeValue{":pk": &dynamodbtypes.AttributeValueMemberS{Value: "A"}},
	}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("query: %v", err))
	}
	// BatchGetItem reads two items: PKC A += 1, B += 1.
	if _, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]dynamodbtypes.KeysAndAttributes{
			tableName: {
				Keys: []map[string]dynamodbtypes.AttributeValue{item("A", "a1"), item("B", "b1")},
			},
		},
	}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("batch get: %v", err))
	}
	// TransactGetItems reads one item: PKC A += 1.
	if _, err := client.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
		TransactItems: []dynamodbtypes.TransactGetItem{{
			Get: &dynamodbtypes.Get{TableName: aws.String(tableName), Key: item("A", "a2")},
		}},
	}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("transact get: %v", err))
	}

	windowEnd := time.Now().Add(time.Minute)

	var pkcRule, skcRule string
	if desc, err := client.DescribeContributorInsights(ctx, &dynamodb.DescribeContributorInsightsInput{TableName: aws.String(tableName)}); err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("describe insights: %v", err))
	} else {
		for _, rule := range desc.ContributorInsightsRuleList {
			if strings.Contains(rule, "-PKC-") {
				pkcRule = rule
			}
			if strings.Contains(rule, "-SKC-") {
				skcRule = rule
			}
		}
	}

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{Endpoint: r.endpoint, Region: r.region})
	if err != nil {
		return setupErr("ContributorInsights_CountsReadsAcrossReadPaths", fmt.Errorf("cloudwatch config: %v", err))
	}
	cw := cloudwatch.NewFromConfig(cfg)

	// Partition series: writes give A 6 and B 3; reads give A 5 (two scan
	// items, one query event, one batch item, one transact item) and B 2
	// (one scan item, one batch item). Total 16 units.
	results = append(results, r.RunTest("dynamodb", "ContributorInsights_CountsReadsAcrossReadPaths", func() error {
		if pkcRule == "" {
			return fmt.Errorf("no PKC rule captured from the describe step")
		}
		report, err := cw.GetInsightRuleReport(ctx, &cloudwatch.GetInsightRuleReportInput{
			RuleName:  aws.String(pkcRule),
			StartTime: aws.Time(windowStart),
			EndTime:   aws.Time(windowEnd),
			Period:    aws.Int32(60),
		})
		if err != nil {
			return err
		}
		if len(report.Contributors) != 2 {
			return fmt.Errorf("expected the A and B partitions, got %v", report.Contributors)
		}
		if keys := report.Contributors[0].Keys; len(keys) == 0 || keys[0] != "A" {
			return fmt.Errorf("expected partition A first, got %v", keys)
		}
		dump := ""
		for _, c := range report.Contributors {
			dump += fmt.Sprintf(" %v=%.0f", c.Keys, aws.ToFloat64(c.ApproximateAggregateValue))
		}
		if got := aws.ToFloat64(report.Contributors[0].ApproximateAggregateValue); got != 11 {
			return fmt.Errorf("expected 11 units on partition A (6 written + 5 read), got %.0f (all:%s)", got, dump)
		}
		if got := aws.ToFloat64(report.Contributors[1].ApproximateAggregateValue); got != 5 {
			return fmt.Errorf("expected 5 units on partition B (3 written + 2 read), got %.0f", got)
		}
		if report.AggregateValue == nil || *report.AggregateValue != 16 {
			return fmt.Errorf("expected a 16 unit aggregate, got %v", report.AggregateValue)
		}
		return nil
	}))

	// Full-key series: every item was written once (3 units) and read once
	// by the scan (1 unit); (A,a1) and (B,b1) were additionally batch-read
	// and (A,a2) transact-read (1 unit each). Every key totals 5 units.
	results = append(results, r.RunTest("dynamodb", "ContributorInsights_ReadPathsFeedSortKeySeries", func() error {
		if skcRule == "" {
			return fmt.Errorf("no SKC rule captured from the describe step")
		}
		report, err := cw.GetInsightRuleReport(ctx, &cloudwatch.GetInsightRuleReportInput{
			RuleName:  aws.String(skcRule),
			StartTime: aws.Time(windowStart),
			EndTime:   aws.Time(windowEnd),
			Period:    aws.Int32(60),
		})
		if err != nil {
			return err
		}
		if len(report.Contributors) != 3 {
			return fmt.Errorf("expected the three item keys, got %v", report.Contributors)
		}
		skcDump := ""
		for _, contributor := range report.Contributors {
			if len(contributor.Keys) != 2 {
				return fmt.Errorf("expected partition and sort key labels, got %v", contributor.Keys)
			}
			skcDump += fmt.Sprintf(" %v=%.0f", contributor.Keys, aws.ToFloat64(contributor.ApproximateAggregateValue))
			if got := aws.ToFloat64(contributor.ApproximateAggregateValue); got != 5 {
				return fmt.Errorf("expected 5 units on key %v (3 written + 1 scanned + 1 read), got %.0f (all:%s)", contributor.Keys, got, skcDump)
			}
		}
		if report.AggregateValue == nil || *report.AggregateValue != 15 {
			return fmt.Errorf("expected a 15 unit aggregate, got %v", report.AggregateValue)
		}
		return nil
	}))

	// A table recreated under the same name must not inherit the dropped
	// table's access counters: the counters are deleted with the table.
	results = append(results, r.RunTest("dynamodb", "ContributorInsights_CountersDroppedWithTable", func() error {
		if _, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}); err != nil {
			return err
		}
		for i := 0; i < 50; i++ {
			if _, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}); err != nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			},
			KeySchema: []dynamodbtypes.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: dynamodbtypes.KeyTypeRange},
			},
			BillingMode: dynamodbtypes.BillingModePayPerRequest,
		}); err != nil {
			return fmt.Errorf("recreate: %v", err)
		}
		if err := waitKinesisDestTableActive(ctx, client, tableName); err != nil {
			return fmt.Errorf("wait active: %v", err)
		}
		if _, err := client.UpdateContributorInsights(ctx, &dynamodb.UpdateContributorInsightsInput{
			TableName:                 aws.String(tableName),
			ContributorInsightsAction: dynamodbtypes.ContributorInsightsActionEnable,
		}); err != nil {
			return err
		}
		desc, err := client.DescribeContributorInsights(ctx, &dynamodb.DescribeContributorInsightsInput{TableName: aws.String(tableName)})
		if err != nil {
			return err
		}
		var freshRule string
		for _, rule := range desc.ContributorInsightsRuleList {
			if strings.Contains(rule, "-PKC-") {
				freshRule = rule
			}
		}
		if freshRule == "" {
			return fmt.Errorf("no PKC rule derived for the recreated table")
		}
		now := time.Now()
		report, err := cw.GetInsightRuleReport(ctx, &cloudwatch.GetInsightRuleReportInput{
			RuleName:  aws.String(freshRule),
			StartTime: aws.Time(now.Add(-10 * time.Minute)),
			EndTime:   aws.Time(now.Add(time.Minute)),
			Period:    aws.Int32(60),
		})
		if err != nil {
			return err
		}
		if len(report.Contributors) != 0 {
			return fmt.Errorf("expected a clean aggregation for the recreated table, got %v", report.Contributors)
		}
		return nil
	}))

	return results
}
