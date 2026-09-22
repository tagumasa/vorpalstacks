package cloudwatchlogs

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The E4 record-backed and aggregate listings: ListLogGroupsForQuery
// (the group set a query actually analysed, served from the query state
// — the record persists it, so the listing survives the restart the way
// GetQueryResults does) and ListAggregateLogGroupSummaries (the
// account-level aggregation of log groups by data source
// characteristics).

// listLogGroupsForQueryCore serves the log groups one query analysed.
// The worked example emits the groups' ARNs, so the listing renders the
// scanned set in the ARN form the example pins. An unknown queryId is
// the operation's ResourceNotFoundException.
func (s *LogsService) listLogGroupsForQueryCore(queryId, nextToken string, maxResults int32, region string) ([]string, string, error) {
	if queryId == "" {
		return nil, "", errRequiredMember("queryId")
	}
	if len(queryId) > logsstore.MaxQueryIdLength {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("queryId must be 1 to %d characters; the given id carries %d", logsstore.MaxQueryIdLength, len(queryId)), 400)
	}
	// The member's range trait is 50..500 — a floor as well as a ceiling;
	// an absent member serves the family's default page, which equals the
	// floor.
	limit := int32(logsstore.ListLogGroupsForQueryMinResults)
	if maxResults != 0 {
		if maxResults < logsstore.ListLogGroupsForQueryMinResults || maxResults > logsstore.ListLogGroupsForQueryMaxResults {
			return nil, "", NewLogsError("InvalidParameterException",
				fmt.Sprintf("maxResults must be between %d and %d", logsstore.ListLogGroupsForQueryMinResults, logsstore.ListLogGroupsForQueryMaxResults), 400)
		}
		limit = maxResults
	}

	var scanned []string
	if val, ok := s.queries.Load(queryId); ok {
		qs := val.(*queryState)
		qs.mu.RLock()
		scanned = qs.scannedGroups
		if len(scanned) == 0 {
			scanned = qs.logGroupNames
		}
		qs.mu.RUnlock()
	} else {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			return nil, "", err
		}
		rec, err := store.GetQueryRecord(queryId)
		if err != nil {
			return nil, "", NewLogsError("ResourceNotFoundException",
				fmt.Sprintf("The specified query does not exist: %s", queryId), 400)
		}
		scanned = rec.ScannedGroups
		if len(scanned) == 0 {
			scanned = rec.LogGroupNames
		}
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}
	identifiers := make([]string, 0, len(scanned))
	for _, name := range scanned {
		identifiers = append(identifiers, store.ARNBuilder().CloudWatch().LogGroup(name))
	}
	// The listing walks the identifiers in their ascending order — the
	// page marker keys on the identifier, so the order must be the
	// deterministic one regardless of the query's input order. The token
	// carries the query's identity, so a token from one query's listing
	// repositions no other.
	sort.Strings(identifiers)
	scope := listingScope("querygroups", queryId)
	result, err := paginateScopedListing(scope, nextToken, identifiers, int(limit),
		func(id string) string { return id })
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// ListAggregateLogGroupSummariesInput is the parsed request of
// ListAggregateLogGroupSummaries.
type ListAggregateLogGroupSummariesInput struct {
	AccountIdentifiers    []string
	IncludeLinkedAccounts bool
	LogGroupClass         string
	LogGroupNamePattern   string
	DataSources           []DataSourceFilterInput
	GroupBy               string
	NextToken             string
	Limit                 int32
	Region                string
}

// DataSourceFilterInput is one parsed dataSources entry: a name pattern
// (required member) and an optional type pattern.
type DataSourceFilterInput struct {
	Name string
	Type string
}

// listAggregateLogGroupSummariesCore aggregates the region's log groups
// by data source characteristics. No log group on this platform carries
// a data source association — the recorded position of the query plane's
// SOURCE resolution, of GetLogFields' AWS::Logs::LogGroup-only typing
// and of the FIELD_INDEX_POLICY data-source rejection — so every group
// that passes the class, name-pattern and account filters shares the
// same (absent) characteristics and aggregates into the one summary
// bucket with no grouping identifiers, under either groupBy value. A
// dataSources filter selects only "log groups associated with the
// specified data sources" and therefore matches nothing.
func (s *LogsService) listAggregateLogGroupSummariesCore(input ListAggregateLogGroupSummariesInput) ([]map[string]interface{}, error) {
	// groupBy is the model's required member, as the API reference also
	// documents ("Required: Yes").
	if input.GroupBy == "" {
		return nil, errRequiredMember("groupBy")
	}
	if input.GroupBy != "DATA_SOURCE_NAME_TYPE_AND_FORMAT" && input.GroupBy != "DATA_SOURCE_NAME_AND_TYPE" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid groupBy: %s. Valid values: DATA_SOURCE_NAME_TYPE_AND_FORMAT, DATA_SOURCE_NAME_AND_TYPE", input.GroupBy), 400)
	}
	limit, err := validateListLimit(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, err
	}
	if input.LogGroupClass != "" &&
		input.LogGroupClass != "STANDARD" && input.LogGroupClass != "INFREQUENT_ACCESS" && input.LogGroupClass != "DELIVERY" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid logGroupClass: %s. Valid values: STANDARD, INFREQUENT_ACCESS, DELIVERY", input.LogGroupClass), 400)
	}
	if len(input.AccountIdentifiers) > logsstore.AggregateAccountIdentifiersMax {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("accountIdentifiers can include as many as %d account IDs", logsstore.AggregateAccountIdentifiersMax), 400)
	}
	if len(input.DataSources) > logsstore.AggregateDataSourceFiltersMax {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("dataSources can include as many as %d filters", logsstore.AggregateDataSourceFiltersMax), 400)
	}
	for _, ds := range input.DataSources {
		if ds.Name == "" {
			return nil, errRequiredMember("dataSources member name")
		}
	}
	var nameFilter func(string) bool
	if input.LogGroupNamePattern != "" {
		matcher, err := parseLogGroupNameMatcher(input.LogGroupNamePattern)
		if err != nil {
			return nil, err
		}
		nameFilter = matcher.matches
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	// A data-sources filter matches no log group: no association exists
	// to match against.
	if len(input.DataSources) > 0 {
		return []map[string]interface{}{}, nil
	}

	// Account scoping follows the DescribeLogGroups precedent:
	// accountIdentifiers applies only with includeLinkedAccounts, and a
	// list naming only other accounts owns nothing on this single-account
	// platform.
	if aggregateAccountScoped(input) {
		local := false
		for _, id := range input.AccountIdentifiers {
			if id == "" || id == s.accountID {
				local = true
				break
			}
		}
		if !local {
			return []map[string]interface{}{}, nil
		}
	}

	all, err := fetchAllLogGroups(store)
	if err != nil {
		return nil, mapStoreError(err)
	}
	count := 0
	for _, lg := range all {
		// A group with no recorded class is Standard — the read plane's
		// standing default for the class member.
		class := lg.LogGroupClass
		if class == "" {
			class = "STANDARD"
		}
		if input.LogGroupClass != "" && class != input.LogGroupClass {
			continue
		}
		if nameFilter != nil && !nameFilter(lg.Name) {
			continue
		}
		count++
	}
	if count == 0 {
		return []map[string]interface{}{}, nil
	}
	buckets := []map[string]interface{}{
		{
			"logGroupCount":       count,
			"groupingIdentifiers": []interface{}{},
		},
	}
	// The paging members hold the documented contract over the summary
	// list instead of being ignored: the platform's aggregation yields at
	// most one bucket, so the first (and only) page carries no next
	// token, and a valid token — which can only name the page past it —
	// serves the empty remainder. The limit bounds the page the same way,
	// keeping the member load-bearing if a second bucket shape (a future
	// data-source association) ever appears.
	if input.NextToken != "" {
		// The token's scope carries every request-identity member the
		// listing answers from (the page-token discipline: a token minted
		// under one filter set repositions no other) — groupBy plus the
		// class, name-pattern and account members.
		if _, err := decodeListingToken(input.NextToken, listingScope("aggregateloggroupsummaries",
			input.GroupBy, input.LogGroupClass, input.LogGroupNamePattern,
			strings.Join(input.AccountIdentifiers, "\x1f"),
			fmt.Sprintf("%t", input.IncludeLinkedAccounts))); err != nil {
			return nil, err
		}
		return []map[string]interface{}{}, nil
	}
	if int32(len(buckets)) > limit {
		buckets = buckets[:limit]
	}
	return buckets, nil
}

// aggregateAccountScoped reports whether the account members restrict
// the aggregation to accounts other than the caller's own, on the
// DescribeLogGroups accountScoped pattern.
func aggregateAccountScoped(input ListAggregateLogGroupSummariesInput) bool {
	if !input.IncludeLinkedAccounts || len(input.AccountIdentifiers) == 0 {
		return false
	}
	for _, id := range input.AccountIdentifiers {
		if id == "" {
			return false
		}
	}
	return true
}

// --- HTTP handlers ---

func (s *LogsService) ListLogGroupsForQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifiers, nextToken, err := s.listLogGroupsForQueryCore(
		request.GetParamLowerFirst(req.Parameters, "QueryId"),
		request.GetParamLowerFirst(req.Parameters, "NextToken"),
		int32(request.GetIntParam(req.Parameters, "MaxResults")),
		reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{
		"logGroupIdentifiers": identifiers,
	}
	if nextToken != "" {
		resp["nextToken"] = nextToken
	}
	return resp, nil
}

func (s *LogsService) ListAggregateLogGroupSummaries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var dataSources []DataSourceFilterInput
	for _, entry := range request.GetListParamLowerFirst(req.Parameters, "DataSources") {
		dataSources = append(dataSources, DataSourceFilterInput{
			Name: request.GetStringParam(entry, "name"),
			Type: request.GetStringParam(entry, "type"),
		})
	}
	summaries, err := s.listAggregateLogGroupSummariesCore(ListAggregateLogGroupSummariesInput{
		AccountIdentifiers:    request.GetStringList(req.Parameters, "AccountIdentifiers"),
		IncludeLinkedAccounts: request.GetBoolParam(req.Parameters, "IncludeLinkedAccounts"),
		LogGroupClass:         request.GetParamLowerFirst(req.Parameters, "LogGroupClass"),
		LogGroupNamePattern:   request.GetParamLowerFirst(req.Parameters, "LogGroupNamePattern"),
		DataSources:           dataSources,
		GroupBy:               request.GetParamLowerFirst(req.Parameters, "GroupBy"),
		NextToken:             request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:                 int32(request.GetIntParam(req.Parameters, "Limit")),
		Region:                reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"aggregateLogGroupSummaries": summaries,
	}, nil
}
