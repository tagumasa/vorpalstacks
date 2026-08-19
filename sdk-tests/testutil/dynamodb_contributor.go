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

// dynamoDBContributorInsightsTests pins the Contributor Insights contract:
// enabling a table derives the documented CloudWatch rule names, item
// reads and writes are aggregated per tracked key and served through the
// CloudWatch report path, and disabling removes the rules.
func (r *TestRunner) dynamoDBContributorInsightsTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()
	tableName := fmt.Sprintf("contributor-%d", suffix)

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
		return setupErr("ContributorInsights_EnableReturnsRuleList", fmt.Errorf("create table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
	if err := waitKinesisDestTableActive(ctx, client, tableName); err != nil {
		return setupErr("ContributorInsights_EnableReturnsRuleList", fmt.Errorf("wait active: %v", err))
	}

	if _, err := client.UpdateContributorInsights(ctx, &dynamodb.UpdateContributorInsightsInput{
		TableName:                 aws.String(tableName),
		ContributorInsightsAction: dynamodbtypes.ContributorInsightsActionEnable,
	}); err != nil {
		return setupErr("ContributorInsights_EnableReturnsRuleList", fmt.Errorf("enable: %v", err))
	}

	var pkcRule string
	results = append(results, r.RunTest("dynamodb", "ContributorInsights_EnableReturnsRuleList", func() error {
		desc, err := client.DescribeContributorInsights(ctx, &dynamodb.DescribeContributorInsightsInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if desc.ContributorInsightsStatus != dynamodbtypes.ContributorInsightsStatusEnabled {
			return fmt.Errorf("expected ENABLED status, got %v", desc.ContributorInsightsStatus)
		}
		// A composite-key table in the default mode derives the four
		// documented rules: PKC, SKC, PKT, SKT.
		if len(desc.ContributorInsightsRuleList) != 4 {
			return fmt.Errorf("expected 4 rules, got %v", desc.ContributorInsightsRuleList)
		}
		var layouts []string
		for _, rule := range desc.ContributorInsightsRuleList {
			if !strings.HasPrefix(rule, "DynamoDBContributorInsights-") {
				return fmt.Errorf("rule %q lacks the DynamoDB prefix", rule)
			}
			parts := strings.Split(strings.TrimPrefix(rule, "DynamoDBContributorInsights-"), "-")
			if len(parts) < 3 {
				return fmt.Errorf("rule %q is not in layout-resource-timestamp form", rule)
			}
			if resource := strings.Join(parts[1:len(parts)-1], "-"); resource != tableName {
				return fmt.Errorf("rule %q does not reference the table", rule)
			}
			layouts = append(layouts, parts[0])
			if parts[0] == "PKC" {
				pkcRule = rule
			}
		}
		for _, want := range []string{"PKC", "SKC", "PKT", "SKT"} {
			found := false
			for _, l := range layouts {
				if l == want {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected a %s rule, got %v", want, layouts)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ContributorInsights_AggregatesKeyAccess", func() error {
		if pkcRule == "" {
			return fmt.Errorf("no PKC rule captured from the describe step")
		}
		// Hot partition: five writes and two reads. Cold partition: one write.
		hot := map[string]dynamodbtypes.AttributeValue{
			"pk": &dynamodbtypes.AttributeValueMemberS{Value: "hot-key"},
			"sk": &dynamodbtypes.AttributeValueMemberS{Value: "a"},
		}
		for i := 0; i < 5; i++ {
			if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(tableName), Item: hot}); err != nil {
				return err
			}
		}
		for i := 0; i < 2; i++ {
			if _, err := client.GetItem(ctx, &dynamodb.GetItemInput{TableName: aws.String(tableName), Key: hot}); err != nil {
				return err
			}
		}
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: "cold-key"},
				"sk": &dynamodbtypes.AttributeValueMemberS{Value: "b"},
			},
		}); err != nil {
			return err
		}

		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{Endpoint: r.endpoint, Region: r.region})
		if err != nil {
			return err
		}
		cw := cloudwatch.NewFromConfig(cfg)
		now := time.Now()
		report, err := cw.GetInsightRuleReport(ctx, &cloudwatch.GetInsightRuleReportInput{
			RuleName:  aws.String(pkcRule),
			StartTime: aws.Time(now.Add(-5 * time.Minute)),
			EndTime:   aws.Time(now.Add(time.Minute)),
			Period:    aws.Int32(60),
		})
		if err != nil {
			return err
		}
		if len(report.Contributors) == 0 {
			return fmt.Errorf("expected aggregated contributors, got none")
		}
		top := report.Contributors[0]
		if len(top.Keys) == 0 || top.Keys[0] != "hot-key" {
			return fmt.Errorf("expected the hot partition key first, got %v", top.Keys)
		}
		if len(report.Contributors) < 2 || report.Contributors[1].Keys[0] != "cold-key" {
			return fmt.Errorf("expected the cold partition key second, got %v", report.Contributors)
		}
		if report.AggregateValue == nil || *report.AggregateValue <= 0 {
			return fmt.Errorf("expected a positive aggregate value, got %v", report.AggregateValue)
		}

		listed, err := cw.DescribeInsightRules(ctx, &cloudwatch.DescribeInsightRulesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, rule := range listed.InsightRules {
			if aws.ToString(rule.Name) == pkcRule {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("DynamoDB rule %s not listed by DescribeInsightRules", pkcRule)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ContributorInsights_DisableRemovesRules", func() error {
		if _, err := client.UpdateContributorInsights(ctx, &dynamodb.UpdateContributorInsightsInput{
			TableName:                 aws.String(tableName),
			ContributorInsightsAction: dynamodbtypes.ContributorInsightsActionDisable,
		}); err != nil {
			return err
		}
		desc, err := client.DescribeContributorInsights(ctx, &dynamodb.DescribeContributorInsightsInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if desc.ContributorInsightsStatus != dynamodbtypes.ContributorInsightsStatusDisabled {
			return fmt.Errorf("expected DISABLED status, got %v", desc.ContributorInsightsStatus)
		}
		if len(desc.ContributorInsightsRuleList) != 0 {
			return fmt.Errorf("expected no rules after disable, got %v", desc.ContributorInsightsRuleList)
		}

		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{Endpoint: r.endpoint, Region: r.region})
		if err != nil {
			return err
		}
		cw := cloudwatch.NewFromConfig(cfg)
		listed, err := cw.DescribeInsightRules(ctx, &cloudwatch.DescribeInsightRulesInput{})
		if err != nil {
			return err
		}
		for _, rule := range listed.InsightRules {
			if aws.ToString(rule.Name) == pkcRule {
				return fmt.Errorf("rule %s still listed after disable", pkcRule)
			}
		}
		return nil
	}))

	// ListContributorInsights must expose the summary shape while insights
	// are enabled on at least this table.
	results = append(results, r.RunTest("dynamodb", "ListContributorInsights_ReturnsSummaries", func() error {
		if _, err := client.UpdateContributorInsights(ctx, &dynamodb.UpdateContributorInsightsInput{
			TableName:                 aws.String(tableName),
			ContributorInsightsAction: dynamodbtypes.ContributorInsightsActionEnable,
		}); err != nil {
			return err
		}
		resp, err := client.ListContributorInsights(ctx, &dynamodb.ListContributorInsightsInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if len(resp.ContributorInsightsSummaries) == 0 {
			return fmt.Errorf("expected a summary for %s", tableName)
		}
		found := false
		for _, s := range resp.ContributorInsightsSummaries {
			if aws.ToString(s.TableName) == tableName {
				found = true
				if s.ContributorInsightsStatus != dynamodbtypes.ContributorInsightsStatusEnabled {
					return fmt.Errorf("expected ENABLED summary status, got %v", s.ContributorInsightsStatus)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("table %s missing from summaries", tableName)
		}
		return nil
	}))

	return results
}
