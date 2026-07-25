package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailQueryAITests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// Create a dedicated EDS for query AI tests.
	var edsID string
	results = append(results, r.RunTest("cloudtrail", "QueryAI_Setup_EDS", func() error {
		resp, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:                         aws.String(fmt.Sprintf("queryai-test-%d", time.Now().UnixNano())),
			TerminationProtectionEnabled: aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("create EDS for query AI: %w", err)
		}
		if resp.EventDataStoreArn == nil {
			return fmt.Errorf("EDS ARN is nil")
		}
		edsID = *resp.EventDataStoreArn
		return nil
	}))
	defer func() {
		if edsID != "" {
			_, _ = tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
				EventDataStore: aws.String(edsID),
			})
		}
	}()

	// GenerateQuery.
	results = append(results, r.RunTest("cloudtrail", "GenerateQuery_Success", func() error {
		resp, err := tc.client.GenerateQuery(tc.ctx, &cloudtrail.GenerateQueryInput{
			EventDataStores: []string{edsID},
			Prompt:          aws.String("Show me all S3 events"),
		})
		if err != nil {
			return fmt.Errorf("GenerateQuery failed: %w", err)
		}
		if resp.QueryStatement == nil || *resp.QueryStatement == "" {
			return fmt.Errorf("QueryStatement is empty")
		}
		return nil
	}))

	// GenerateQuery_MissingPrompt.
	results = append(results, r.RunTest("cloudtrail", "GenerateQuery_MissingPrompt", func() error {
		_, err := tc.client.GenerateQuery(tc.ctx, &cloudtrail.GenerateQueryInput{
			EventDataStores: []string{edsID},
		})
		if err == nil {
			return fmt.Errorf("expected error for missing Prompt")
		}
		return nil
	}))

	// SearchSampleQueries.
	results = append(results, r.RunTest("cloudtrail", "SearchSampleQueries_Success", func() error {
		resp, err := tc.client.SearchSampleQueries(tc.ctx, &cloudtrail.SearchSampleQueriesInput{
			SearchPhrase: aws.String("S3"),
		})
		if err != nil {
			return fmt.Errorf("SearchSampleQueries failed: %w", err)
		}
		if len(resp.SearchResults) == 0 {
			return fmt.Errorf("expected at least 1 result for 'S3'")
		}
		return nil
	}))

	// SearchSampleQueries_NoMatch.
	results = append(results, r.RunTest("cloudtrail", "SearchSampleQueries_NoMatch", func() error {
		resp, err := tc.client.SearchSampleQueries(tc.ctx, &cloudtrail.SearchSampleQueriesInput{
			SearchPhrase: aws.String("xyznonexistent"),
		})
		if err != nil {
			return fmt.Errorf("SearchSampleQueries with no match failed: %w", err)
		}
		if len(resp.SearchResults) != 0 {
			return fmt.Errorf("expected 0 results for non-matching phrase, got %d", len(resp.SearchResults))
		}
		return nil
	}))

	return results
}
