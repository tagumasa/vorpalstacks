package cloudwatchlogs

import (
	"fmt"
	"sort"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The group-listing family: the ListLogGroups and DescribeLogGroups cores,
// their shared input/result shapes and the filter vocabulary every listing
// member contributes (name prefix and pattern, tag filters, data sources,
// field indexes, account scoping).
// TagFilterInput is the ListLogGroups logGroupTags member's element: a
// tag key and the optional values any-of semantics rides on.
type TagFilterInput struct {
	Key    string
	Values []string
}

// ListLogGroupsInput is the transport-agnostic input for listing log groups.
type ListLogGroupsInput struct {
	LogGroupNamePrefix string
	// LogGroupNamePattern is the DescribeLogGroups pattern member: a
	// case-sensitive substring match (its member documentation).
	LogGroupNamePattern string
	// LogGroupIdentifiers restricts DescribeLogGroups to the named groups
	// (name-or-ARN per element); "the only other filter that you can
	// choose to specify is includeLinkedAccounts".
	LogGroupIdentifiers   []string
	AccountIdentifiers    []string
	IncludeLinkedAccounts bool
	LogGroupClass         string
	// LogGroupTags is the ListLogGroups logGroupTags member: "An array of
	// tag filters to return only log groups that have specific tags.
	// Multiple filters are combined with AND logic."
	LogGroupTags []TagFilterInput
	// DataSources is the ListLogGroups dataSources member: filters "by
	// data source name, type, or both. Multiple filters within the same
	// dimension are combined with OR logic, while filters across different
	// dimensions are combined with AND logic."
	DataSources []DataSourceFilterInput
	// FieldIndexNames is the ListLogGroups fieldIndexNames member: "Only
	// log groups containing all specified field indexes are returned."
	FieldIndexNames []string
	// NameFilter is the predicate form of the ListLogGroups
	// logGroupNamePattern member (anchored-or-substring alternatives); the
	// listing filters the full store listing through it.
	NameFilter func(string) bool
	// NameFilterPattern carries the pattern string the matcher was built
	// from — the request identity the listing's page tokens scope to (a
	// function cannot ride in a token).
	NameFilterPattern string
	NextToken         string
	Limit             int32
	Region            string
}

// listingScopeDigest digests the request identity a log-group listing's
// page tokens scope to: prefix, pattern (the Describe substring member or
// the matcher's pattern string), class, the identifier set (whose order
// carries no meaning) and the ListLogGroups filter members (tag filters,
// data source filters and field index names, in request order — a token
// minted under one filter set never serves another).
func (in *ListLogGroupsInput) listingScopeDigest() string {
	ids := make([]string, len(in.LogGroupIdentifiers))
	copy(ids, in.LogGroupIdentifiers)
	sort.Strings(ids)
	pattern := in.LogGroupNamePattern
	if pattern == "" {
		pattern = in.NameFilterPattern
	}
	tagParts := make([]string, 0, len(in.LogGroupTags))
	for _, tf := range in.LogGroupTags {
		tagParts = append(tagParts, tf.Key+"\x1f"+strings.Join(tf.Values, "\x1e"))
	}
	dsParts := make([]string, 0, len(in.DataSources))
	for _, ds := range in.DataSources {
		dsParts = append(dsParts, ds.Name+"\x1f"+ds.Type)
	}
	return listingScope("loggroups", in.LogGroupNamePrefix, pattern, in.LogGroupClass, strings.Join(ids, "\x1f"),
		strings.Join(tagParts, "\x1e"), strings.Join(dsParts, "\x1e"), strings.Join(in.FieldIndexNames, "\x1f"))
}

// ListLogGroupsResult holds the outcome of a successful log group listing.
type ListLogGroupsResult struct {
	LogGroups []*logsstore.LogGroup
	NextToken string
}

// listLogGroupsCore lists log groups with optional class filtering,
// enforcing the operation's ListLimit bound (1..1000, default 50).
func (s *LogsService) listLogGroupsCore(input ListLogGroupsInput) (*ListLogGroupsResult, error) {
	limit, err := validateListLimit(input.Limit, logsstore.DefaultListLogGroupsLimit, logsstore.MaxListLogGroupsLimit)
	if err != nil {
		return nil, err
	}
	input.Limit = limit

	// The class member carries the modelled enum on both operations that
	// share this core ("Valid Values: STANDARD | INFREQUENT_ACCESS |
	// DELIVERY").
	if input.LogGroupClass != "" && !validateLogGroupClass(input.LogGroupClass) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid log group class: %s. Valid values: STANDARD, INFREQUENT_ACCESS, DELIVERY", input.LogGroupClass), 400)
	}
	if err := validateLogGroupTagFilters(input.LogGroupTags); err != nil {
		return nil, err
	}
	if err := validateDataSourceFilters(input.DataSources); err != nil {
		return nil, err
	}
	if err := validateFieldIndexNameFilters(input.FieldIndexNames); err != nil {
		return nil, err
	}
	if err := validateAccountIdentifierList(input.AccountIdentifiers); err != nil {
		return nil, err
	}
	// The accountIdentifiers members of both listing operations target the
	// shared AccountIds shape ("You can specify as many as 20 account IDs
	// in the array"), so the count bound holds on this listing core
	// whichever operation reaches it.
	if len(input.AccountIdentifiers) > maxDescribeAccountIdentifiers {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("accountIdentifiers accepts as many as %d account IDs; the request carries %d",
				maxDescribeAccountIdentifiers, len(input.AccountIdentifiers)), 400)
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	// The listing pages through the scoped token vocabulary: the request
	// identity digests into the token's scope, and the value cursor the
	// store scan or the filtered walk produces never surfaces bare — a
	// token replayed against a different filter rejects.
	scope := input.listingScopeDigest()
	if input.NextToken != "" {
		cursor, decErr := decodeListingToken(input.NextToken, scope)
		if decErr != nil {
			return nil, decErr
		}
		input.NextToken = cursor
	}

	var result *ListLogGroupsResult
	// The identifier, pattern, class and matcher members filter the full
	// listing (the store's own scan is a prefix match, which is not those
	// members' documented semantics); those listings paginate by group
	// name, the marker the name-ordered scan makes exact. The class
	// filter rides the filtered walk too — it "limit[s] the results to
	// only those log groups in the specified log group class"
	// (DescribeLogGroups logGroupClass member documentation), so a page
	// is a page of matching groups, never an empty carrier of a token.
	// The ListLogGroups filter members (tags, data sources, field indexes)
	// ride the same walk: "If you specify more than one filter type, the
	// results include log groups that satisfy all filters."
	if len(input.LogGroupIdentifiers) > 0 || input.LogGroupNamePattern != "" || input.NameFilter != nil || input.LogGroupClass != "" ||
		len(input.LogGroupTags) > 0 || len(input.DataSources) > 0 || len(input.FieldIndexNames) > 0 || input.accountScoped() {
		result, err = listLogGroupsFiltered(s, store, input, s.accountID)
	} else {
		var groups []*logsstore.LogGroup
		var nextToken string
		groups, nextToken, err = store.ListLogGroups(input.LogGroupNamePrefix, input.NextToken, int(input.Limit))
		if err != nil {
			return nil, mapStoreError(err)
		}

		result = &ListLogGroupsResult{LogGroups: groups, NextToken: nextToken}
	}
	if err != nil {
		return nil, err
	}
	wrapped, wrapErr := wrapScopedListingMarker(scope, result.NextToken)
	if wrapErr != nil {
		return nil, wrapErr
	}
	result.NextToken = wrapped
	return result, nil
}

// describeLogGroupsCore serves DescribeLogGroups on the shared listing
// core under that operation's own DescribeLimit bound (1..50, default 50)
// — the tighter of the two operations sharing the core. The documented
// cross-member rules (member documentation): logGroupNamePattern and
// logGroupNamePrefix are mutually exclusive, and logGroupIdentifiers
// admits includeLinkedAccounts as its only fellow filter.
func (s *LogsService) describeLogGroupsCore(input ListLogGroupsInput) (*ListLogGroupsResult, error) {
	limit, err := validateListLimit(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, err
	}
	input.Limit = limit

	if input.LogGroupNamePattern != "" && input.LogGroupNamePrefix != "" {
		return nil, NewLogsError("InvalidParameterException",
			"logGroupNamePattern and logGroupNamePrefix are mutually exclusive; only one of these parameters can be passed", 400)
	}
	// The list members carry their documented bounds: "You can specify
	// as many as 50 log groups in the array" and "You can specify as
	// many as 20 account IDs in the array" (DescribeLogGroups
	// logGroupIdentifiers/accountIdentifiers; the Smithy length traits
	// 1..50 and 0..20 — the accountIdentifiers count holds on the shared
	// listing core, so it is checked there once).
	if len(input.LogGroupIdentifiers) > maxLogGroupIdentifiers {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("logGroupIdentifiers accepts as many as %d log groups; the request carries %d",
				maxLogGroupIdentifiers, len(input.LogGroupIdentifiers)), 400)
	}
	if len(input.LogGroupIdentifiers) > 0 &&
		(input.LogGroupNamePrefix != "" || input.LogGroupNamePattern != "" || input.LogGroupClass != "") {
		return nil, NewLogsError("InvalidParameterException",
			"when logGroupIdentifiers is specified, the only other filter that can be chosen is includeLinkedAccounts", 400)
	}
	// The filtered listing paths (pattern, identifiers, matcher, account
	// scoping) route inside the shared listing core below.
	return s.listLogGroupsCore(input)
}

// accountScoped reports whether the account members restrict the listing to
// accounts other than the caller's own: accountIdentifiers applies only
// when includeLinkedAccounts is true (member documentation), and on this
// single-account platform a list that names other accounts owns nothing.
// A null entry (an empty element) means all accounts, per the member
// documentation, which on this platform is the local account alone.
func (in *ListLogGroupsInput) accountScoped() bool {
	if !in.IncludeLinkedAccounts || len(in.AccountIdentifiers) == 0 {
		return false
	}
	for _, id := range in.AccountIdentifiers {
		if id == "" {
			return false
		}
	}
	return true
}

// listLogGroupsFiltered serves the pattern / identifier / account-scoped /
// filter-member listings over the full store listing, paginating by group
// name (the listing is name-ascending, so the name marker is exact).
func listLogGroupsFiltered(svc *LogsService, store *logsstore.Store, input ListLogGroupsInput, localAccountID string) (*ListLogGroupsResult, error) {
	if input.accountScoped() {
		local := false
		for _, id := range input.AccountIdentifiers {
			if id == "" || id == localAccountID {
				local = true
				break
			}
		}
		if !local {
			return &ListLogGroupsResult{LogGroups: []*logsstore.LogGroup{}}, nil
		}
	}

	all, err := fetchAllLogGroups(store)
	if err != nil {
		return nil, mapStoreError(err)
	}

	var selected []*logsstore.LogGroup
	if len(input.LogGroupIdentifiers) > 0 {
		// Name-or-ARN per element; identifiers that address no stored
		// group match nothing (the member filters the returned list).
		wanted := make(map[string]bool, len(input.LogGroupIdentifiers))
		for _, id := range input.LogGroupIdentifiers {
			wanted[resolveLogGroupIdentifier(id)] = true
		}
		for _, lg := range all {
			if wanted[lg.Name] && svc.listingFilterMembersMatch(store, input, lg) {
				selected = append(selected, lg)
			}
		}
	} else {
		for _, lg := range all {
			if input.NameFilter != nil {
				if !input.NameFilter(lg.Name) {
					continue
				}
			} else if input.LogGroupNamePattern != "" && !strings.Contains(lg.Name, input.LogGroupNamePattern) {
				continue
			} else if input.LogGroupNamePrefix != "" && !strings.HasPrefix(lg.Name, input.LogGroupNamePrefix) {
				continue
			}
			if input.LogGroupClass != "" && lg.LogGroupClass != input.LogGroupClass {
				continue
			}
			if !svc.listingFilterMembersMatch(store, input, lg) {
				continue
			}
			selected = append(selected, lg)
		}
	}

	result := &ListLogGroupsResult{}
	for i, lg := range selected {
		if input.NextToken != "" && lg.Name <= input.NextToken {
			continue
		}
		result.LogGroups = append(result.LogGroups, lg)
		if len(result.LogGroups) == int(input.Limit) {
			if i+1 < len(selected) {
				result.NextToken = lg.Name
			}
			break
		}
	}
	return result, nil
}

// listingFilterMembersMatch applies the ListLogGroups filter members to
// one group: "If you specify more than one filter type, the results
// include log groups that satisfy all filters."
func (s *LogsService) listingFilterMembersMatch(store *logsstore.Store, input ListLogGroupsInput, lg *logsstore.LogGroup) bool {
	if len(input.LogGroupTags) > 0 {
		tags, err := store.Tags().List(lg.ARN)
		if err != nil {
			tags = nil
		}
		if !logGroupTagFiltersMatch(input.LogGroupTags, tags) {
			return false
		}
	}
	if len(input.DataSources) > 0 && !s.groupMatchesDataSourceFilters(store, input.DataSources, lg) {
		return false
	}
	if len(input.FieldIndexNames) > 0 && !s.groupContainsAllFieldIndexes(store, lg.Name, input.FieldIndexNames) {
		return false
	}
	return true
}

// logGroupTagFiltersMatch applies the logGroupTags member: AND across
// filters, any-of within one key's values. "If you don't specify values,
// the response returns all log groups that are tagged with that key, with
// any or no value."
func logGroupTagFiltersMatch(filters []TagFilterInput, tags map[string]string) bool {
	for _, tf := range filters {
		value, tagged := tags[tf.Key]
		if !tagged {
			return false
		}
		if len(tf.Values) == 0 {
			continue
		}
		any := false
		for _, want := range tf.Values {
			if tagFilterValueMatches(want, value) {
				any = true
				break
			}
		}
		if !any {
			return false
		}
	}
	return true
}

// tagFilterValueMatches one value string against the group's tag value:
// "Use * for wildcard matching... Use ! as a prefix for negation...
// Exact matching and negation are case-sensitive. Wildcard matching is
// case-insensitive." The member's pattern allows only optional leading
// and trailing wildcards.
func tagFilterValueMatches(filterValue, tagValue string) bool {
	if strings.HasPrefix(filterValue, "!") {
		return tagValue != strings.TrimPrefix(filterValue, "!")
	}
	if !strings.Contains(filterValue, "*") {
		return filterValue == tagValue
	}
	needle := strings.ToLower(strings.Trim(filterValue, "*"))
	lower := strings.ToLower(tagValue)
	if needle == "" {
		return true
	}
	leadStar := strings.HasPrefix(filterValue, "*")
	trailStar := strings.HasSuffix(filterValue, "*")
	if leadStar && trailStar {
		return strings.Contains(lower, needle)
	}
	if leadStar {
		return strings.HasSuffix(lower, needle)
	}
	return strings.HasPrefix(lower, needle)
}

// groupMatchesDataSourceFilters applies the dataSources member: "Multiple
// filters within the same dimension are combined with OR logic, while
// filters across different dimensions are combined with AND logic" — the
// name dimension and the type dimension each reduce to a set, and the
// group must satisfy both through its associated data sources. The
// platform's log-group data sources are the delivery sources whose
// resource ARN addresses the group.
func (s *LogsService) groupMatchesDataSourceFilters(store *logsstore.Store, filters []DataSourceFilterInput, lg *logsstore.LogGroup) bool {
	names := make(map[string]bool)
	kinds := make(map[string]bool)
	for _, f := range filters {
		if f.Name != "" {
			names[f.Name] = true
		}
		if f.Type != "" {
			kinds[f.Type] = true
		}
	}
	sources, err := store.ListDeliverySources()
	if err != nil {
		return false
	}
	nameOK := len(names) == 0
	typeOK := len(kinds) == 0
	for _, src := range sources {
		if src.ResourceArn != lg.ARN {
			continue
		}
		if !nameOK && names[src.Name] {
			nameOK = true
		}
		if !typeOK && kinds[src.LogType] {
			typeOK = true
		}
		if nameOK && typeOK {
			return true
		}
	}
	return false
}

// groupContainsAllFieldIndexes applies the fieldIndexNames member: "Only
// log groups containing all specified field indexes are returned." The
// membership basis is the field-index machinery's own: the DEFAULT
// category's fixed set plus the effective custom policy's fields
// (group-level record, else the applying account-level policy).
func (s *LogsService) groupContainsAllFieldIndexes(store *logsstore.Store, groupName string, names []string) bool {
	set := make(map[string]bool, len(defaultFieldIndexes))
	for _, name := range defaultFieldIndexes {
		set[name] = true
	}
	if eff := s.effectiveIndexPolicy(store, groupName); eff != nil {
		for _, spec := range eff.fields {
			set[indexedFieldName(spec.Name)] = true
		}
	}
	for _, want := range names {
		if !set[want] {
			return false
		}
	}
	return true
}
