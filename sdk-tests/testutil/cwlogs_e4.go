package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// The E4 pins: the account-level storage tier policy round-trip, the
// record-backed ListLogGroupsForQuery listing and the
// ListAggregateLogGroupSummaries aggregation.

func (tc *cwlogsTestCtx) e4Tests() []TestResult {
	var results []TestResult

	// The storage tier policy: an account that never set one answers the
	// getter's ResourceNotFoundException; a put stores the tier with an
	// update stamp and the getter round-trips it.
	results = append(results, tc.runner.RunTest("logs", "StorageTierPolicy_PutGetRoundTrip", func() error {
		// The never-set account's ResourceNotFoundException is
		// unit-pinned: against a persistent server a previous run may
		// already have set the policy, so the SDK pin asserts the
		// round-trip alone.
		put, err := tc.client.PutStorageTierPolicy(tc.ctx, &cloudwatchlogs.PutStorageTierPolicyInput{
			StorageTier: types.StorageTierIntelligentTiering,
		})
		if err != nil {
			return fmt.Errorf("put storage tier policy: %v", err)
		}
		if put.StorageTier != types.StorageTierIntelligentTiering {
			return fmt.Errorf("put echo tier = %v", put.StorageTier)
		}
		if put.LastUpdatedTime == nil || *put.LastUpdatedTime <= 0 {
			return fmt.Errorf("put echo lastUpdatedTime = %v", put.LastUpdatedTime)
		}

		got, err := tc.client.GetStorageTierPolicy(tc.ctx, &cloudwatchlogs.GetStorageTierPolicyInput{})
		if err != nil {
			return fmt.Errorf("get storage tier policy: %v", err)
		}
		if got.StorageTier != types.StorageTierIntelligentTiering || got.LastUpdatedTime == nil || *got.LastUpdatedTime != *put.LastUpdatedTime {
			return fmt.Errorf("round trip = %v/%v, want the stored record", got.StorageTier, got.LastUpdatedTime)
		}

		// Reset the account to STANDARD so later runs observe the same
		// never-set-free starting point.
		if _, err := tc.client.PutStorageTierPolicy(tc.ctx, &cloudwatchlogs.PutStorageTierPolicyInput{
			StorageTier: types.StorageTierStandard,
		}); err != nil {
			return fmt.Errorf("reset storage tier policy: %v", err)
		}
		return nil
	}))

	// ListLogGroupsForQuery serves the groups a query analysed, as ARNs;
	// an unknown queryId is ResourceNotFound and the maxResults floor of
	// 50 rejects below-range values.
	results = append(results, tc.runner.RunTest("logs", "ListLogGroupsForQuery_RecordBacked", func() error {
		base := tc.uniquePrefix("l4q")
		groups := []string{base + "-a", base + "-b"}
		for _, g := range groups {
			if err := tc.createLogGroup(g); err != nil {
				return err
			}
			defer tc.deleteLogGroup(g)
			if err := tc.createLogStream(g, "s1"); err != nil {
				return err
			}
			now := time.Now().UnixMilli()
			if err := tc.putLogEvent(g, "s1", "l4q event", now); err != nil {
				return err
			}
		}

		now := time.Now().UnixMilli()
		start, err := tc.client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: groups,
			QueryString:   aws.String("fields @message | limit 2"),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		if _, err := tc.waitForQuery(start.QueryId, 30); err != nil {
			return err
		}

		list, err := tc.client.ListLogGroupsForQuery(tc.ctx, &cloudwatchlogs.ListLogGroupsForQueryInput{
			QueryId: start.QueryId,
		})
		if err != nil {
			return fmt.Errorf("list log groups for query: %v", err)
		}
		if len(list.LogGroupIdentifiers) != 2 {
			return fmt.Errorf("identifiers = %v, want both groups", list.LogGroupIdentifiers)
		}
		for _, id := range list.LogGroupIdentifiers {
			if !strings.Contains(id, ":log-group:") || !strings.HasSuffix(id, groups[0]) && !strings.HasSuffix(id, groups[1]) {
				return fmt.Errorf("identifier %q is not one of the groups' ARNs", id)
			}
		}

		_, err = tc.client.ListLogGroupsForQuery(tc.ctx, &cloudwatchlogs.ListLogGroupsForQueryInput{
			QueryId: aws.String("no-such-query"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("unknown query: %v", err)
		}

		_, err = tc.client.ListLogGroupsForQuery(tc.ctx, &cloudwatchlogs.ListLogGroupsForQueryInput{
			QueryId:    start.QueryId,
			MaxResults: aws.Int32(10),
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	// ListAggregateLogGroupSummaries: every matching group aggregates
	// into the one no-identifier bucket (no log group on the platform
	// carries a data source association); the class and name-pattern
	// filters select, a dataSources filter matches nothing, and the
	// absent groupBy rejects server-side (the model marks no required
	// member; the API reference carries it as Required: Yes).
	results = append(results, tc.runner.RunTest("logs", "ListAggregateLogGroupSummaries_SingleBucket", func() error {
		base := tc.uniquePrefix("agg")
		groups := []string{base + "-a", base + "-b"}
		for _, g := range groups {
			if err := tc.createLogGroup(g); err != nil {
				return err
			}
			defer tc.deleteLogGroup(g)
		}

		summarize := func(input *cloudwatchlogs.ListAggregateLogGroupSummariesInput) ([]types.AggregateLogGroupSummary, error) {
			resp, err := tc.client.ListAggregateLogGroupSummaries(tc.ctx, input)
			if err != nil {
				return nil, err
			}
			return resp.AggregateLogGroupSummaries, nil
		}

		rows, err := summarize(&cloudwatchlogs.ListAggregateLogGroupSummariesInput{
			GroupBy:             types.ListAggregateLogGroupSummariesGroupByDataSourceNameAndType,
			LogGroupClass:       types.LogGroupClassStandard,
			LogGroupNamePattern: aws.String("^" + base),
		})
		if err != nil {
			return fmt.Errorf("aggregate: %v", err)
		}
		if len(rows) != 1 {
			return fmt.Errorf("summaries = %d rows, want the single bucket", len(rows))
		}
		if rows[0].LogGroupCount == nil || *rows[0].LogGroupCount != 2 {
			return fmt.Errorf("logGroupCount = %v, want 2", rows[0].LogGroupCount)
		}
		if len(rows[0].GroupingIdentifiers) != 0 {
			return fmt.Errorf("groupingIdentifiers = %v, want the empty set", rows[0].GroupingIdentifiers)
		}

		filtered, err := summarize(&cloudwatchlogs.ListAggregateLogGroupSummariesInput{
			GroupBy: types.ListAggregateLogGroupSummariesGroupByDataSourceNameAndType,
			DataSources: []types.DataSourceFilter{{
				Name: aws.String("amazon_vpc"),
			}},
		})
		if err != nil {
			return fmt.Errorf("data-sources aggregate: %v", err)
		}
		if len(filtered) != 0 {
			return fmt.Errorf("dataSources filter = %d rows, want none (no group carries a data source)", len(filtered))
		}

		// The absent groupBy never reaches the server: the SDK's own
		// model marks the member required and rejects the request
		// client-side. The server-side rejection is unit-pinned
		// (TestListAggregateLogGroupSummaries).
		return nil
	}))

	return results
}
