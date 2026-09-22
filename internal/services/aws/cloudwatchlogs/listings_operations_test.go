package cloudwatchlogs

import (
	"fmt"
	"strings"
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// TestStorageTierPolicyRoundTrip pins the account-level storage tier
// policy: the never-set account answers the getter's
// ResourceNotFoundException, a valid put stores and echoes the tier with
// an update stamp, and the enum outside its two values rejects.
func TestStorageTierPolicyRoundTrip(t *testing.T) {
	svc, _ := newReadTestService(t, "tier-group")

	if _, err := svc.getStorageTierPolicyCore("us-east-1"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("never-set policy = %v, want ResourceNotFoundException", err)
	}
	policy, err := svc.putStorageTierPolicyCore("INTELLIGENT_TIERING", "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if policy.StorageTier != "INTELLIGENT_TIERING" || policy.LastUpdatedTime <= 0 {
		t.Fatalf("put echo = %+v", policy)
	}
	got, err := svc.getStorageTierPolicyCore("us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.StorageTier != "INTELLIGENT_TIERING" || got.LastUpdatedTime != policy.LastUpdatedTime {
		t.Fatalf("get = %+v, want the stored round-trip", got)
	}
	if _, err := svc.putStorageTierPolicyCore("COLD", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("foreign tier = %v, want InvalidParameterException", err)
	}
	if _, err := svc.putStorageTierPolicyCore("", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("absent tier = %v, want InvalidParameterException", err)
	}
}

// TestListLogGroupsForQuery pins the record-backed listing: an
// identifier-selected query serves its explicit set as ARNs, a SOURCE
// query serves the resolved set, the member bounds hold (floor 50,
// ceiling 500), an unknown id is ResourceNotFound, and pagination walks
// pages at the default floor.
func TestListLogGroupsForQuery(t *testing.T) {
	svc, store := newReadTestService(t, "l4q-a")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "l4q-a")); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateLogGroup(logsstore.NewLogGroup("l4q-b", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "l4q-b")); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	putReadEvents(t, store, "l4q-a", "s1", now-60000, 1)
	putReadEvents(t, store, "l4q-b", "s1", now-60000, 1)

	queryId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      now/1000 - 3600,
		EndTime:        now/1000 + 60,
		QueryString:    "fields @message",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"l4q-b", "l4q-a"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	waitTerminalQuery(t, svc, queryId)

	arnA := store.ARNBuilder().CloudWatch().LogGroup("l4q-a")
	arnB := store.ARNBuilder().CloudWatch().LogGroup("l4q-b")
	ids, next, err := svc.listLogGroupsForQueryCore(queryId, "", 0, "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != arnA || ids[1] != arnB || next != "" {
		t.Fatalf("listing = %v (next %q), want both ARNs ascending", ids, next)
	}

	// A SOURCE-selected query serves the resolved group set, not the
	// explicit one (a SOURCE command replaces any explicitly named groups
	// at execution — the dynamically selected case the listing exists
	// for).
	sourceId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      now/1000 - 3600,
		EndTime:        now/1000 + 60,
		QueryString:    "SOURCE logGroups(namePrefix: ['l4q-']) | fields @message",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"l4q-a"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	waitTerminalQuery(t, svc, sourceId)
	ids, _, err = svc.listLogGroupsForQueryCore(sourceId, "", 0, "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != arnA || ids[1] != arnB {
		t.Fatalf("SOURCE listing = %v, want the resolved set as ARNs", ids)
	}

	if _, _, err := svc.listLogGroupsForQueryCore("query-nope", "", 0, "us-east-1"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("unknown query = %v, want ResourceNotFoundException", err)
	}
	if _, _, err := svc.listLogGroupsForQueryCore(queryId, "", 10, "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("below-floor maxResults = %v, want InvalidParameterException", err)
	}
	if _, _, err := svc.listLogGroupsForQueryCore(queryId, "", 501, "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("above-ceiling maxResults = %v, want InvalidParameterException", err)
	}

	// Pagination: a SOURCE prefix wide enough to select more groups than
	// the default page serves the first page plus a token that resumes.
	for i := 0; i < 60; i++ {
		name := fmt.Sprintf("l4q-p%02d", i)
		if err := store.CreateLogGroup(logsstore.NewLogGroup(name, "us-east-1", "000000000000")); err != nil {
			t.Fatal(err)
		}
	}
	pageId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      now/1000 - 3600,
		EndTime:        now/1000 + 60,
		QueryString:    "SOURCE logGroups(namePrefix: ['l4q-p']) | fields @message",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"l4q-a"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	waitTerminalQuery(t, svc, pageId)
	first, nextToken, err := svc.listLogGroupsForQueryCore(pageId, "", 0, "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 50 || nextToken == "" {
		t.Fatalf("first page = %d items (next %q), want the 50 floor with a token", len(first), nextToken)
	}
	second, nextToken, err := svc.listLogGroupsForQueryCore(pageId, nextToken, 0, "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(second) != 10 || nextToken != "" {
		t.Fatalf("second page = %d items (next %q), want the 10 remainder", len(second), nextToken)
	}
	if second[0] != store.ARNBuilder().CloudWatch().LogGroup("l4q-p50") {
		t.Fatalf("second page resumes at %s, want l4q-p50", second[0])
	}
}

// TestListAggregateLogGroupSummaries pins the aggregation's platform
// contract: every matching group aggregates into the one no-identifier
// bucket (no log group carries a data source association), the class and
// name-pattern filters select, a dataSources filter matches nothing,
// foreign account scoping owns nothing, and groupBy is required.
func TestListAggregateLogGroupSummaries(t *testing.T) {
	svc, store := newReadTestService(t, "agg-a")
	if err := store.CreateLogGroup(logsstore.NewLogGroup("agg-b", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	ia := logsstore.NewLogGroup("agg-ia", "us-east-1", "000000000000")
	ia.LogGroupClass = "INFREQUENT_ACCESS"
	if err := store.CreateLogGroup(ia); err != nil {
		t.Fatal(err)
	}

	summaryCount := func(t *testing.T, input ListAggregateLogGroupSummariesInput) int {
		t.Helper()
		summaries, err := svc.listAggregateLogGroupSummariesCore(input)
		if err != nil {
			t.Fatalf("aggregate: %v", err)
		}
		if len(summaries) > 1 {
			t.Fatalf("summaries = %d, want the single no-association bucket", len(summaries))
		}
		if len(summaries) == 0 {
			return 0
		}
		if count, ok := summaries[0]["logGroupCount"].(int); !ok || count < 0 {
			t.Fatalf("logGroupCount = %v, want an int", summaries[0]["logGroupCount"])
		}
		idents := summaries[0]["groupingIdentifiers"].([]interface{})
		if len(idents) != 0 {
			t.Fatalf("groupingIdentifiers = %v, want the empty set", idents)
		}
		return summaries[0]["logGroupCount"].(int)
	}
	base := ListAggregateLogGroupSummariesInput{
		GroupBy: "DATA_SOURCE_NAME_AND_TYPE",
		Region:  "us-east-1",
	}

	if got := summaryCount(t, base); got != 3 {
		t.Fatalf("unfiltered count = %d, want 3", got)
	}
	classFilter := base
	classFilter.LogGroupClass = "STANDARD"
	if got := summaryCount(t, classFilter); got != 2 {
		t.Fatalf("STANDARD count = %d, want 2", got)
	}
	iaFilter := base
	iaFilter.LogGroupClass = "INFREQUENT_ACCESS"
	if got := summaryCount(t, iaFilter); got != 1 {
		t.Fatalf("INFREQUENT_ACCESS count = %d, want 1", got)
	}
	pattern := base
	pattern.LogGroupNamePattern = "^agg-a"
	if got := summaryCount(t, pattern); got != 1 {
		t.Fatalf("^agg-a count = %d, want 1", got)
	}

	filtered := base
	filtered.DataSources = []DataSourceFilterInput{{Name: "amazon_vpc"}}
	summaries, err := svc.listAggregateLogGroupSummariesCore(filtered)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("dataSources filter = %v, want the empty list", summaries)
	}

	foreign := base
	foreign.IncludeLinkedAccounts = true
	foreign.AccountIdentifiers = []string{"111122223333"}
	summaries, err = svc.listAggregateLogGroupSummariesCore(foreign)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("foreign account scoping = %v, want the empty list", summaries)
	}

	noGroupBy := base
	noGroupBy.GroupBy = ""
	if _, err := svc.listAggregateLogGroupSummariesCore(noGroupBy); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("absent groupBy = %v, want InvalidParameterException", err)
	}
	badGroupBy := base
	badGroupBy.GroupBy = "LOG_GROUP_NAME"
	if _, err := svc.listAggregateLogGroupSummariesCore(badGroupBy); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("foreign groupBy = %v, want InvalidParameterException", err)
	}
	// The paging members hold the contract instead of being ignored: a
	// garbage token rejects, a valid one names the page past the single
	// bucket and serves the empty remainder, and the limit bounds the
	// page (a limit of 1 still serves the bucket).
	paged := base
	paged.NextToken = "garbage-token"
	if _, err := svc.listAggregateLogGroupSummariesCore(paged); err == nil {
		t.Fatalf("garbage nextToken accepted")
	}
	aggregateScope := listingScope("aggregateloggroupsummaries", base.GroupBy, base.LogGroupClass,
		base.LogGroupNamePattern, strings.Join(base.AccountIdentifiers, "\x1f"),
		fmt.Sprintf("%t", base.IncludeLinkedAccounts))
	validToken, err := encodeListingToken(aggregateScope, "past")
	if err != nil {
		t.Fatal(err)
	}
	paged.NextToken = validToken
	summaries, err = svc.listAggregateLogGroupSummariesCore(paged)
	if err != nil {
		t.Fatalf("valid nextToken: %v", err)
	}
	if len(summaries) != 0 {
		t.Fatalf("valid nextToken remainder = %v, want the empty page", summaries)
	}
	// The scope carries the filter members: a token minted under another
	// request's filters repositions no other listing.
	otherFilters := base
	otherFilters.LogGroupClass = "STANDARD"
	otherScope := listingScope("aggregateloggroupsummaries", otherFilters.GroupBy, otherFilters.LogGroupClass,
		otherFilters.LogGroupNamePattern, strings.Join(otherFilters.AccountIdentifiers, "\x1f"),
		fmt.Sprintf("%t", otherFilters.IncludeLinkedAccounts))
	foreignToken, err := encodeListingToken(otherScope, "past")
	if err != nil {
		t.Fatal(err)
	}
	replay := base
	replay.NextToken = foreignToken
	if _, err := svc.listAggregateLogGroupSummariesCore(replay); err == nil {
		t.Fatalf("token minted under another filter set must reject")
	}
	limited := base
	limited.Limit = 1
	if got := summaryCount(t, limited); got != 3 {
		t.Fatalf("limit-1 page count = %d, want the single bucket's 3", got)
	}
	badClass := base
	badClass.LogGroupClass = "COLD"
	if _, err := svc.listAggregateLogGroupSummariesCore(badClass); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("foreign class = %v, want InvalidParameterException", err)
	}
}

// The name-cursor listings page through scoped tokens: the walk still
// pages to completion, a token minted under one filter rejects under
// another, and a typed bare name is not a token.
func TestListingsPageThroughScopedTokens(t *testing.T) {
	svc, store := newReadTestService(t, "seed-group")
	for _, name := range []string{"scoped-app", "scoped-beta", "scoped-gamma"} {
		if err := store.CreateLogGroup(logsstore.NewLogGroup(name, "us-east-1", "000000000000")); err != nil {
			t.Fatal(err)
		}
	}

	// The prefix walk pages through the scoped vocabulary to completion.
	page1, err := svc.listLogGroupsCore(ListLogGroupsInput{
		LogGroupNamePrefix: "scoped-", Limit: 2, Region: "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1.LogGroups) != 2 || page1.NextToken == "" {
		t.Fatalf("first page = %d groups, token %q", len(page1.LogGroups), page1.NextToken)
	}
	page2, err := svc.listLogGroupsCore(ListLogGroupsInput{
		LogGroupNamePrefix: "scoped-", NextToken: page1.NextToken, Limit: 2, Region: "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.LogGroups) != 1 || page2.NextToken != "" {
		t.Fatalf("second page = %d groups, token %q", len(page2.LogGroups), page2.NextToken)
	}

	// The same token repositions no other walk: a narrower prefix and a
	// pattern walk both reject it, and a typed bare name is not a token.
	reject := func(label string, input ListLogGroupsInput) {
		t.Helper()
		if _, err := svc.listLogGroupsCore(input); logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("%s: %v, want InvalidParameterException", label, err)
		}
	}
	reject("token under a different prefix", ListLogGroupsInput{
		LogGroupNamePrefix: "scoped-app", NextToken: page1.NextToken, Limit: 2, Region: "us-east-1",
	})
	reject("token under the pattern walk", ListLogGroupsInput{
		LogGroupNamePattern: "beta", NextToken: page1.NextToken, Limit: 2, Region: "us-east-1",
	})
	reject("a typed bare name", ListLogGroupsInput{
		LogGroupNamePrefix: "scoped-", NextToken: "scoped-beta", Limit: 2, Region: "us-east-1",
	})

	// The metric-filters listing scopes to its group: a group-A token
	// replayed against group B rejects instead of skipping its keys.
	seed := func(group string, names ...string) {
		if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
			t.Fatal(err)
		}
		for _, name := range names {
			if err := store.PutMetricFilter(&logsstore.MetricFilter{
				Name: name, LogGroupName: group, FilterPattern: "ERROR",
				MetricTransformations: []logsstore.MetricTransformation{{
					MetricName: "Count", MetricNamespace: "App", MetricValue: "1",
				}},
			}); err != nil {
				t.Fatal(err)
			}
		}
	}
	seed("scoped-filters-a", "f1", "f2", "f3")
	seed("scoped-filters-b", "g1", "g2", "g3")
	filtersA, tokenA, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "scoped-filters-a", Limit: 2, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtersA) != 2 || tokenA == "" {
		t.Fatalf("metric filter page = %d filters, token %q", len(filtersA), tokenA)
	}
	if _, _, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "scoped-filters-b", NextToken: tokenA, Limit: 2, Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("cross-group token replay: %v, want InvalidParameterException", err)
	}
}

// TestClassFilteredPrefixWalkFillsPages pins the class filter's pagination
// contract: the filter "limit[s] the results to only those log groups in
// the specified log group class" (DescribeLogGroups logGroupClass member
// documentation), so a page is a page of matching groups — a scan page
// whose groups are all of another class is stepped over, never served as
// an empty page still carrying a token.
func TestClassFilteredPrefixWalkFillsPages(t *testing.T) {
	svc, store := newReadTestService(t, "cfw-seed")
	ia := func(name string) {
		lg := logsstore.NewLogGroup(name, "us-east-1", "000000000000")
		lg.LogGroupClass = "INFREQUENT_ACCESS"
		if err := store.CreateLogGroup(lg); err != nil {
			t.Fatal(err)
		}
	}
	ia("cfw-alpha")
	ia("cfw-beta")
	if err := store.CreateLogGroup(logsstore.NewLogGroup("cfw-gamma", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateLogGroup(logsstore.NewLogGroup("cfw-other", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}

	walk := ListLogGroupsInput{
		LogGroupNamePrefix: "cfw-", LogGroupClass: "INFREQUENT_ACCESS",
		Limit: 1, Region: "us-east-1",
	}
	page1, err := svc.listLogGroupsCore(walk)
	if err != nil {
		t.Fatal(err)
	}
	if len(page1.LogGroups) != 1 || page1.LogGroups[0].Name != "cfw-alpha" || page1.NextToken == "" {
		t.Fatalf("page 1 = %+v token %q, want [cfw-alpha] with a token", page1.LogGroups, page1.NextToken)
	}
	walk.NextToken = page1.NextToken
	page2, err := svc.listLogGroupsCore(walk)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.LogGroups) != 1 || page2.LogGroups[0].Name != "cfw-beta" || page2.NextToken != "" {
		t.Fatalf("page 2 = %+v token %q, want [cfw-beta] terminal", page2.LogGroups, page2.NextToken)
	}

	// A class nothing carries walks to a token-free empty page, not a
	// chain of tokened empty pages.
	walk.NextToken = ""
	walk.LogGroupClass = "DELIVERY"
	empty, err := svc.listLogGroupsCore(walk)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.LogGroups) != 0 || empty.NextToken != "" {
		t.Fatalf("DELIVERY walk = %d groups token %q, want empty and terminal", len(empty.LogGroups), empty.NextToken)
	}
}
