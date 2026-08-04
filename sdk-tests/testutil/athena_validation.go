package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

// waitForQuery polls GetQueryExecution until the query reaches a terminal
// state (SUCCEEDED, FAILED, or CANCELLED) or the max poll count is reached.
func (tc *athenaTestCtx) waitForQuery(queryExecutionId string) (types.QueryExecutionState, error) {
	for i := 0; i < 60; i++ {
		resp, err := tc.client.GetQueryExecution(tc.ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: aws.String(queryExecutionId),
		})
		if err != nil {
			return "", err
		}
		state := resp.QueryExecution.Status.State
		if state == types.QueryExecutionStateSucceeded ||
			state == types.QueryExecutionStateFailed ||
			state == types.QueryExecutionStateCancelled {
			if state != types.QueryExecutionStateSucceeded {
				reason := aws.ToString(resp.QueryExecution.Status.StateChangeReason)
				return state, fmt.Errorf("query %s ended in %s: %s", queryExecutionId, state, reason)
			}
			return state, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return "", fmt.Errorf("query %s did not complete within 30s", queryExecutionId)
}

func (tc *athenaTestCtx) testValidation() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	// MaxResults=0 must be accepted as the "use the service default"
	// sentinel and must not produce a NextToken that points back at the
	// first page, which would otherwise cause callers to loop forever.
	results = append(results, tc.runner.RunTest("athena", "ListNamedQueries_MaxResultsZero", func() error {
		resp, err := client.ListNamedQueries(ctx, &athena.ListNamedQueriesInput{
			MaxResults: aws.Int32(0),
		})
		if err != nil {
			return fmt.Errorf("MaxResults=0 should be accepted (as default), got: %v", err)
		}
		// Should return a valid response (not empty-with-NextToken loop)
		if resp == nil {
			return fmt.Errorf("expected non-nil response for MaxResults=0")
		}
		// NextToken should not point back to offset 0 causing a loop
		if resp.NextToken != nil && aws.ToString(resp.NextToken) == "0" {
			return fmt.Errorf("NextToken=0 indicates infinite loop bug")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListQueryExecutions_MaxResultsZero", func() error {
		resp, err := client.ListQueryExecutions(ctx, &athena.ListQueryExecutionsInput{
			MaxResults: aws.Int32(0),
		})
		if err != nil {
			return fmt.Errorf("MaxResults=0 should be accepted (as default), got: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("expected non-nil response for MaxResults=0")
		}
		if resp.NextToken != nil && aws.ToString(resp.NextToken) == "0" {
			return fmt.Errorf("NextToken=0 indicates infinite loop bug")
		}
		return nil
	}))

	// --- BytesScannedCutoffPerQuery strict min ---
	results = append(results, tc.runner.RunTest("athena", "CreateWorkGroup_BytesScannedCutoff_Invalid", func() error {
		wgName := fmt.Sprintf("bsc-wg-%d", time.Now().UnixNano()%1000000000)
		defer client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wgName)})

		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(wgName),
			Configuration: &types.WorkGroupConfiguration{
				BytesScannedCutoffPerQuery: aws.Int64(5000000),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for BytesScannedCutoffPerQuery=5000000 (< 10000000)")
		}
		if !strings.Contains(err.Error(), "InvalidParameterException") && !strings.Contains(err.Error(), "InvalidRequestException") {
			return fmt.Errorf("expected InvalidParameterException, got: %v", err)
		}
		return nil
	}))

	// --- ClientRequestToken length validation ---
	results = append(results, tc.runner.RunTest("athena", "StartQueryExecution_ClientRequestToken_TooShort", func() error {
		_, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString:        aws.String("SELECT 1"),
			ClientRequestToken: aws.String("short-token"),
		})
		if err == nil {
			return fmt.Errorf("expected error for ClientRequestToken < 32 chars")
		}
		if !strings.Contains(err.Error(), "InvalidParameterException") && !strings.Contains(err.Error(), "InvalidRequestException") {
			return fmt.Errorf("expected validation error, got: %v", err)
		}
		return nil
	}))

	// --- ExecutionRole ARN pattern validation ---
	results = append(results, tc.runner.RunTest("athena", "CreateWorkGroup_InvalidExecutionRole", func() error {
		wgName := fmt.Sprintf("erole-wg-%d", time.Now().UnixNano()%1000000000)
		defer client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wgName)})

		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(wgName),
			Configuration: &types.WorkGroupConfiguration{
				ExecutionRole: aws.String("not-a-valid-arn"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid ExecutionRole ARN")
		}
		if !strings.Contains(err.Error(), "InvalidParameterException") && !strings.Contains(err.Error(), "InvalidRequestException") {
			return fmt.Errorf("expected validation error, got: %v", err)
		}
		return nil
	}))

	// --- AdditionalConfiguration length validation ---
	results = append(results, tc.runner.RunTest("athena", "CreateWorkGroup_AdditionalConfiguration_TooLong", func() error {
		wgName := fmt.Sprintf("ac-wg-%d", time.Now().UnixNano()%1000000000)
		defer client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wgName)})

		longValue := make([]byte, 129)
		for i := range longValue {
			longValue[i] = 'x'
		}

		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(wgName),
			Configuration: &types.WorkGroupConfiguration{
				AdditionalConfiguration: aws.String(string(longValue)),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for AdditionalConfiguration > 128 chars")
		}
		if !strings.Contains(err.Error(), "InvalidParameterException") && !strings.Contains(err.Error(), "InvalidRequestException") {
			return fmt.Errorf("expected validation error, got: %v", err)
		}
		return nil
	}))

	// --- Multi-page pagination: ListNamedQueries ---
	// PaginateSlice is shared by 8+ Athena List operations. ListNamedQueries
	// (MaxResults min=0, max=50 per Smithy) allows MaxResults=2, making real
	// multi-page traversal testable. We create 5 queries and page through them.
	results = append(results, tc.runner.RunTest("athena", "ListNamedQueries_MultiPage", func() error {
		prefix := fmt.Sprintf("mpq-%d-", time.Now().UnixNano()%1000000000)

		// Create 5 named queries in the primary workgroup.
		myIds := make(map[string]bool)
		for i := 0; i < 5; i++ {
			resp, err := client.CreateNamedQuery(ctx, &athena.CreateNamedQueryInput{
				Name:        aws.String(fmt.Sprintf("%s%d", prefix, i)),
				Database:    aws.String("default"),
				QueryString: aws.String("SELECT 1"),
			})
			if err != nil {
				return fmt.Errorf("CreateNamedQuery %d: %w", i, err)
			}
			id := aws.ToString(resp.NamedQueryId)
			if id == "" {
				return fmt.Errorf("CreateNamedQuery %d returned empty id", i)
			}
			myIds[id] = true
			defer client.DeleteNamedQuery(ctx, &athena.DeleteNamedQueryInput{
				NamedQueryId: aws.String(id),
			})
		}

		// Walk all pages with MaxResults=2 (min=0 in Smithy, so 2 is valid).
		collected := make(map[string]int) // id → page number it appeared on
		var nextToken *string
		pagesWithItems := 0
		for page := 0; page < 50; page++ {
			resp, err := client.ListNamedQueries(ctx, &athena.ListNamedQueriesInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				return fmt.Errorf("ListNamedQueries page %d: %w", page, err)
			}
			if len(resp.NamedQueryIds) == 0 {
				if resp.NextToken != nil {
					return fmt.Errorf("page %d: empty page with non-nil NextToken", page)
				}
				break
			}
			pagesWithItems++
			if len(resp.NamedQueryIds) > 2 {
				return fmt.Errorf("page %d returned %d ids, MaxResults=2", page, len(resp.NamedQueryIds))
			}
			for _, id := range resp.NamedQueryIds {
				if _, dup := collected[id]; dup {
					return fmt.Errorf("id %s appeared on multiple pages", id)
				}
				collected[id] = page
			}
			nextToken = resp.NextToken
			if nextToken == nil {
				break
			}
		}

		if pagesWithItems < 3 {
			return fmt.Errorf("expected ≥3 pages (≥5 queries at MaxResults=2), got %d", pagesWithItems)
		}
		for id := range myIds {
			if _, ok := collected[id]; !ok {
				return fmt.Errorf("created query %s not found in any page", id)
			}
		}
		return nil
	}))

	// --- Multi-page pagination: ListQueryExecutions ---
	results = append(results, tc.runner.RunTest("athena", "ListQueryExecutions_MultiPage", func() error {
		// Start 5 trivial queries to guarantee multi-page at MaxResults=2.
		myIds := make(map[string]bool)
		for i := 0; i < 5; i++ {
			resp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
				QueryString: aws.String("SELECT 1"),
			})
			if err != nil {
				return fmt.Errorf("StartQueryExecution %d: %w", i, err)
			}
			id := aws.ToString(resp.QueryExecutionId)
			if id == "" {
				return fmt.Errorf("StartQueryExecution %d returned empty id", i)
			}
			myIds[id] = true
		}

		// Walk pages with MaxResults=2. Query executions accumulate across
		// test runs (they cannot be deleted), so we may need many pages.
		// Exit early once all our IDs are found, but require ≥3 pages to
		// prove real multi-page traversal.
		collected := make(map[string]int)
		var nextToken *string
		pagesWithItems := 0
		for page := 0; page < 500; page++ {
			resp, err := client.ListQueryExecutions(ctx, &athena.ListQueryExecutionsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				return fmt.Errorf("ListQueryExecutions page %d: %w", page, err)
			}
			if len(resp.QueryExecutionIds) == 0 {
				if resp.NextToken != nil {
					return fmt.Errorf("page %d: empty page with non-nil NextToken", page)
				}
				break
			}
			pagesWithItems++
			if len(resp.QueryExecutionIds) > 2 {
				return fmt.Errorf("page %d returned %d ids, MaxResults=2", page, len(resp.QueryExecutionIds))
			}
			for _, id := range resp.QueryExecutionIds {
				if _, dup := collected[id]; dup {
					return fmt.Errorf("id %s appeared on multiple pages", id)
				}
				collected[id] = page
			}
			nextToken = resp.NextToken
			if nextToken == nil {
				break
			}
			// Early exit: all IDs found and we've seen enough pages.
			if pagesWithItems >= 3 {
				allFound := true
				for id := range myIds {
					if _, ok := collected[id]; !ok {
						allFound = false
						break
					}
				}
				if allFound {
					break
				}
			}
		}

		if pagesWithItems < 3 {
			return fmt.Errorf("expected ≥3 pages (≥5 executions at MaxResults=2), got %d", pagesWithItems)
		}
		for id := range myIds {
			if _, ok := collected[id]; !ok {
				return fmt.Errorf("started execution %s not found in any page", id)
			}
		}
		return nil
	}))

	// --- P1: ListTagsForResource accepts MaxResults ---
	// AWS spec: ListTagsForResource MaxResults minimum is 75 (confirmed via
	// Smithy model MaxTagsCount range and AWS docs). The per-resource tag
	// limit is 50, so 50 < 75 means NextToken is never produced in practice.
	// Multi-page pagination for PaginateSlice is verified by
	// ListNamedQueries_MultiPage and ListQueryExecutions_MultiPage above,
	// which use the same pagination code with testable MaxResults values.
	results = append(results, tc.runner.RunTest("athena", "ListTagsForResource_WithMaxResults", func() error {
		wgName := fmt.Sprintf("tagmax-wg-%d", time.Now().UnixNano()%1000000000)
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(wgName),
		})
		if err != nil {
			return err
		}
		defer client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wgName)})

		var tags []types.Tag
		for i := 0; i < 5; i++ {
			tags = append(tags, types.Tag{
				Key:   aws.String(fmt.Sprintf("key%d", i)),
				Value: aws.String(fmt.Sprintf("val%d", i)),
			})
		}
		arn := fmt.Sprintf("arn:aws:athena:%s:%s:workgroup/%s", tc.runner.region, tc.runner.AccountID(), wgName)
		if _, err := client.TagResource(ctx, &athena.TagResourceInput{
			ResourceARN: aws.String(arn),
			Tags:        tags,
		}); err != nil {
			return err
		}

		resp, err := client.ListTagsForResource(ctx, &athena.ListTagsForResourceInput{
			ResourceARN: aws.String(arn),
			MaxResults:  aws.Int32(75),
		})
		if err != nil {
			return fmt.Errorf("ListTagsForResource with MaxResults=75: %w", err)
		}

		// 5 tags < MaxResults=75 → single page is expected.
		if resp.NextToken != nil {
			return fmt.Errorf("expected nil NextToken (5 tags fit in one page at MaxResults=75), got %q", aws.ToString(resp.NextToken))
		}
		if len(resp.Tags) != 5 {
			return fmt.Errorf("expected 5 tags, got %d", len(resp.Tags))
		}

		collected := make(map[string]string)
		for _, t := range resp.Tags {
			collected[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		for i := 0; i < 5; i++ {
			key := fmt.Sprintf("key%d", i)
			expected := fmt.Sprintf("val%d", i)
			if v, ok := collected[key]; !ok {
				return fmt.Errorf("tag %q missing", key)
			} else if v != expected {
				return fmt.Errorf("tag %q = %q, expected %q", key, v, expected)
			}
		}
		return nil
	}))

	return results
}
