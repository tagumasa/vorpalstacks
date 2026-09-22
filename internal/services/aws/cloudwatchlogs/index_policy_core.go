package cloudwatchlogs

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The field-index cores: effective-policy resolution (group level first,
// else the account-level FIELD_INDEX_POLICY), the three operation cores,
// the derived DescribeFieldIndexes listing and the account-policy
// selection validation.

// --- Effective policy resolution ---

// effectiveIndexPolicy is a group's resolved field index configuration.
type effectiveIndexPolicy struct {
	source         string // LOG_GROUP or ACCOUNT
	fields         []logsstore.IndexFieldSpec
	policyName     string
	policyDocument string
	lastUpdateTime int64
}

// effectiveIndexPolicy resolves the group-level record first, else the
// account-level FIELD_INDEX_POLICY that applies to the group (longest
// matching LogGroupNamePrefix, else the account-wide policy). A nil
// return means the group carries no custom field indexes.
func (s *LogsService) effectiveIndexPolicy(store *logsstore.Store, groupName string) *effectiveIndexPolicy {
	rec, err := store.GetIndexPolicy(groupName)
	if err != nil {
		rec = nil
	}
	return s.effectiveIndexPolicyFrom(store, groupName, rec)
}

// effectiveIndexPolicyFrom resolves against a pre-read group record
// (nil when absent): the listing path reads the record and the INACTIVE
// trail as one mutex-held snapshot, so a concurrent policy replacement
// cannot leave the response carrying a replaced field under neither
// category.
func (s *LogsService) effectiveIndexPolicyFrom(store *logsstore.Store, groupName string, rec *logsstore.IndexPolicy) *effectiveIndexPolicy {
	if rec != nil {
		if fields, perr := parseIndexPolicyDocument(rec.PolicyDocument); perr == nil {
			return &effectiveIndexPolicy{
				source:         "LOG_GROUP",
				fields:         fields,
				policyDocument: rec.PolicyDocument,
				lastUpdateTime: rec.LastUpdateTime,
			}
		}
	}
	ap := s.matchingFieldIndexAccountPolicy(store, groupName)
	if ap == nil {
		return nil
	}
	fields, err := parseIndexPolicyDocument(ap.PolicyDocument)
	if err != nil {
		return nil
	}
	return &effectiveIndexPolicy{
		source:         "ACCOUNT",
		fields:         fields,
		policyName:     ap.PolicyName,
		policyDocument: ap.PolicyDocument,
		lastUpdateTime: ap.LastUpdatedTime,
	}
}

// matchingFieldIndexAccountPolicy resolves the account-level field index
// policy for a group: the FIELD_INDEX_POLICY whose prefix is the group's
// longest matching prefix, else the account-wide policy (empty
// selectionCriteria). The no-overlapping-prefixes rule makes the longest
// prefix unambiguous in the valid configuration space.
func (s *LogsService) matchingFieldIndexAccountPolicy(store *logsstore.Store, groupName string) *logsstore.AccountPolicy {
	policies, err := store.ListAccountPolicies("FIELD_INDEX_POLICY", "")
	if err != nil {
		return nil
	}
	var best *logsstore.AccountPolicy
	bestLen := -1
	for _, policy := range policies {
		prefix, ok := transformerPolicyPrefix(policy.SelectionCriteria)
		if !ok {
			continue
		}
		if prefix != "" && !strings.HasPrefix(groupName, prefix) {
			continue
		}
		if len(prefix) > bestLen {
			best = policy
			bestLen = len(prefix)
		}
	}
	return best
}

// --- Cores ---

// mergeIndexInactiveTrail folds replaced fields into the group's INACTIVE
// trail: every previously indexed field absent from the newly effective
// set. The caller resolves the effective set; an account-level policy
// taking over after a deletion subtracts its own fields the same way.
func mergeIndexInactiveTrail(trail []logsstore.IndexFieldSpec, replaced []logsstore.IndexFieldSpec, effectiveNames map[string]bool) []logsstore.IndexFieldSpec {
	merged := make([]logsstore.IndexFieldSpec, 0, len(trail)+len(replaced))
	seen := make(map[string]bool, cap(merged))
	for _, group := range [2][]logsstore.IndexFieldSpec{trail, replaced} {
		for _, spec := range group {
			name := indexedFieldName(spec.Name)
			if effectiveNames[name] || seen[name] {
				continue
			}
			seen[name] = true
			merged = append(merged, logsstore.IndexFieldSpec{Name: name, Type: spec.Type})
		}
	}
	sort.Slice(merged, func(i, j int) bool { return merged[i].Name < merged[j].Name })
	return merged
}

// putIndexPolicyCore validates and stores one log group's field index
// policy ("Log group-level field index policies created with
// PutIndexPolicy override account-level field index policies ... the log
// group uses only that policy for log group-level indexing").
func (s *LogsService) putIndexPolicyCore(identifier, policyDocument, region string) (*logsstore.IndexPolicy, error) {
	if identifier == "" {
		return nil, errRequiredMember("logGroupIdentifier")
	}
	// The member's element traits bind the identifier alphabet (1..2048
	// of "[\w#+=/:,.@-]"): a foreign character rejects here — and the
	// alphabet excludes '*', so the documented "Don't include an * at
	// the end" (PutIndexPolicy logGroupIdentifier) rejects in the same
	// check instead of being trimmed away.
	if err := validateLogGroupIdentifierElements([]string{identifier}); err != nil {
		return nil, err
	}
	if policyDocument == "" {
		return nil, errRequiredMember("policyDocument")
	}
	// The generic policyDocument trait this validator family shares
	// (1..51200 Unicode characters): the one constant, counted the one
	// way, on every policyDocument surface.
	if utf8.RuneCountInString(policyDocument) > logsstore.MaxPolicyDocumentLength {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("The policyDocument must not exceed %d characters", logsstore.MaxPolicyDocumentLength), 400)
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	groupName := resolveLogGroupIdentifier(identifier)
	group, err := store.GetLogGroup(groupName)
	if err != nil {
		return nil, mapStoreError(err)
	}
	// "Only log groups in the Standard log class support field index
	// policies."
	if group.LogGroupClass != "" && group.LogGroupClass != "STANDARD" {
		return nil, NewLogsError("InvalidParameterException",
			"Field index policies are supported only for log groups in the Standard log class", 400)
	}
	fields, err := parseIndexPolicyDocument(policyDocument)
	if err != nil {
		return nil, err
	}
	effectiveNames := make(map[string]bool, len(fields))
	for _, spec := range fields {
		effectiveNames[indexedFieldName(spec.Name)] = true
	}
	// The trail merge and the policy write run as one locked section:
	// concurrent policy changes on one group compute from the current
	// records, never from a shared stale base.
	var policy *logsstore.IndexPolicy
	err = store.ReplaceIndexPolicy(groupName, func(current *logsstore.IndexPolicy, trail *logsstore.IndexInactiveTrail) (*logsstore.IndexPolicy, *logsstore.IndexInactiveTrail, error) {
		var replaced []logsstore.IndexFieldSpec
		if current != nil {
			if oldFields, perr := parseIndexPolicyDocument(current.PolicyDocument); perr == nil {
				replaced = oldFields
			}
		}
		newTrail := &logsstore.IndexInactiveTrail{
			LogGroupName: groupName,
			Fields:       mergeIndexInactiveTrail(trail.Fields, replaced, effectiveNames),
		}
		policy = &logsstore.IndexPolicy{
			LogGroupName:       groupName,
			LogGroupIdentifier: store.ARNBuilder().CloudWatch().LogGroup(groupName),
			PolicyDocument:     policyDocument,
		}
		return policy, newTrail, nil
	})
	if err != nil {
		return nil, mapStoreError(err)
	}
	return policy, nil
}

// indexPolicyView pairs one returned policy with the queried group's
// ARN — the response's logGroupIdentifier member carries the group's
// ARN for both the group-level and the account-level policy case.
type indexPolicyView struct {
	groupARN string
	eff      *effectiveIndexPolicy
}

// describeIndexPoliciesCore returns the one log group's effective field
// index policy: "If a specified log group has a log-group level index
// policy, that policy is returned by this operation. If a specified log
// group doesn't have a log-group level index policy, but an account-wide
// index policy applies to it, that account-wide policy is returned."
// The single-entry listing never pages (the request accepts exactly one
// group), so the response omits nextToken.
func (s *LogsService) describeIndexPoliciesCore(identifiers []string, region string) ([]*indexPolicyView, error) {
	if len(identifiers) == 0 {
		return nil, errRequiredMember("logGroupIdentifiers")
	}
	// DescribeIndexPoliciesLogGroupIdentifiers carries the length trait
	// 1..1 ("Array Members: Fixed number of 1 item"); the element rides
	// the shared identifier alphabet.
	if len(identifiers) > 1 {
		return nil, NewLogsError("InvalidParameterException",
			"logGroupIdentifiers accepts exactly one log group identifier", 400)
	}
	if err := validateLogGroupIdentifierElements(identifiers); err != nil {
		return nil, err
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	groupName := resolveLogGroupIdentifier(identifiers[0])
	if _, err := store.GetLogGroup(groupName); err != nil {
		return nil, mapStoreError(err)
	}
	if eff := s.effectiveIndexPolicy(store, groupName); eff != nil {
		return []*indexPolicyView{{
			groupARN: store.ARNBuilder().CloudWatch().LogGroup(groupName),
			eff:      eff,
		}}, nil
	}
	return []*indexPolicyView{}, nil
}

// deleteIndexPolicyCore removes a log group's group-level policy
// ("You can't use this operation to delete an account-level index
// policy"). The deleted policy's fields that no surviving policy indexes
// join the INACTIVE trail; a group whose account-level policy applies
// falls back to it ("in a few minutes the log group begins using that
// account-wide policy" — the platform's dynamic resolution applies it at
// once).
func (s *LogsService) deleteIndexPolicyCore(identifier, region string) error {
	if identifier == "" {
		return errRequiredMember("logGroupIdentifier")
	}
	// The member's element traits bind the identifier alphabet; a
	// foreign character rejects here rather than surfacing as the
	// family's not-found error.
	if err := validateLogGroupIdentifierElements([]string{identifier}); err != nil {
		return err
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	groupName := resolveLogGroupIdentifier(identifier)
	if _, err := store.GetLogGroup(groupName); err != nil {
		return mapStoreError(err)
	}
	if _, err := store.GetIndexPolicy(groupName); err != nil {
		return NewLogsError("ResourceNotFoundException",
			fmt.Sprintf("No field index policy is associated with log group %s", groupName), 400)
	}
	effectiveNames := make(map[string]bool)
	if account := s.matchingFieldIndexAccountPolicy(store, groupName); account != nil {
		if fields, perr := parseIndexPolicyDocument(account.PolicyDocument); perr == nil {
			for _, spec := range fields {
				effectiveNames[indexedFieldName(spec.Name)] = true
			}
		}
	}
	// The trail merge and the policy removal run as one locked section
	// with the put path's replacement.
	return mapStoreError(store.RemoveIndexPolicy(groupName, func(current *logsstore.IndexPolicy, trail *logsstore.IndexInactiveTrail) (*logsstore.IndexInactiveTrail, error) {
		deleted, derr := parseIndexPolicyDocument(current.PolicyDocument)
		if derr != nil {
			deleted = nil
		}
		return &logsstore.IndexInactiveTrail{
			LogGroupName: groupName,
			Fields:       mergeIndexInactiveTrail(trail.Fields, deleted, effectiveNames),
		}, nil
	}))
}

// fieldIndexEntry is one DescribeFieldIndexes row before serialisation.
type fieldIndexEntry struct {
	logGroupIdentifier string // empty when the entry's policy is not single-group
	fieldIndexName     string
	indexCategory      string
	indexType          string
	firstEventTime     int64
	lastEventTime      int64
	lastScanTime       int64
	// cursor is the page-token key of this row: the group, field name and
	// category together — a field the policy indexes on top of its
	// DEFAULT entry repeats its fieldIndexName, and the same field name
	// recurs across the listed groups, so the name alone keys no row.
	cursor string
}

// describeFieldIndexesCore derives the field-index listing of the named
// groups: the DEFAULT category's fixed fields, the effective policy's
// CUSTOM fields and the INACTIVE trail, each row carrying the bounds the
// group's ingestion scans accumulated.
func (s *LogsService) describeFieldIndexesCore(identifiers, categories []string, region, nextToken string) ([]*fieldIndexEntry, string, error) {
	if len(identifiers) == 0 {
		return nil, "", errRequiredMember("logGroupIdentifiers")
	}
	if len(identifiers) > logsstore.DescribeFieldIndexesGroupsMax {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("logGroupIdentifiers accepts at most %d log groups", logsstore.DescribeFieldIndexesGroupsMax), 400)
	}
	// Each element rides the shared identifier alphabet; a foreign
	// character rejects here rather than surfacing as a lookup miss.
	if err := validateLogGroupIdentifierElements(identifiers); err != nil {
		return nil, "", err
	}
	selected := make(map[string]bool, len(categories))
	if len(categories) > 4 {
		return nil, "", NewLogsError("InvalidParameterException",
			"indexCategories accepts at most 4 categories", 400)
	}
	for _, category := range categories {
		switch category {
		case "DEFAULT", "CUSTOM", "AUTO", "INACTIVE":
			selected[category] = true
		default:
			return nil, "", NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid index category: %s. Allowed values: DEFAULT, CUSTOM, AUTO, INACTIVE", category), 400)
		}
	}
	// "If you omit this parameter, the response includes the DEFAULT,
	// CUSTOM, and INACTIVE categories."
	if len(selected) == 0 {
		selected["DEFAULT"], selected["CUSTOM"], selected["INACTIVE"] = true, true, true
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}
	entries := make([]*fieldIndexEntry, 0)
	for _, identifier := range identifiers {
		groupName := resolveLogGroupIdentifier(identifier)
		if _, err := store.GetLogGroup(groupName); err != nil {
			return nil, "", mapStoreError(err)
		}
		// The policy half and the trail half compose one response: a
		// replacement committing between two unlocked reads could drop a
		// replaced field's row entirely (neither its active nor its
		// inactive category), so the pair reads as one mutex-held
		// snapshot.
		policyRec, trail := store.ReadIndexPolicySnapshot(groupName)
		eff := s.effectiveIndexPolicyFrom(store, groupName, policyRec)
		effectiveNames := make(map[string]bool)
		if eff != nil {
			for _, spec := range eff.fields {
				effectiveNames[indexedFieldName(spec.Name)] = true
			}
		}
		type candidate struct {
			spec        logsstore.IndexFieldSpec
			category    string
			singleGroup bool
		}
		var candidates []candidate
		if selected["DEFAULT"] {
			for _, name := range defaultFieldIndexes {
				candidates = append(candidates, candidate{
					spec:     logsstore.IndexFieldSpec{Name: name, Type: "FIELD_INDEX"},
					category: "DEFAULT",
				})
			}
		}
		if selected["CUSTOM"] && eff != nil {
			for _, spec := range eff.fields {
				candidates = append(candidates, candidate{
					spec:        spec,
					category:    "CUSTOM",
					singleGroup: eff.source == "LOG_GROUP",
				})
			}
		}
		if selected["INACTIVE"] {
			for _, spec := range trail.Fields {
				if effectiveNames[indexedFieldName(spec.Name)] {
					continue
				}
				candidates = append(candidates, candidate{spec: spec, category: "INACTIVE"})
			}
		}
		// The AUTO category reports no fields on this platform: the
		// automatic query-pattern field selection it describes is an
		// AWS-side learning machinery with no substrate here — the
		// category stays a selectable, honestly-empty listing.
		if len(candidates) == 0 {
			continue
		}
		// Each row's matching-event bounds come from the group's
		// accumulated ingestion-scan record: firstEventTime is "the
		// earliest log event that matches this field index, after the
		// index policy that contains it was created", so only events
		// scanned while a policy watched the field carry. The DEFAULT
		// identity fields alone derive theirs from the stream records —
		// they match every event by construction, including events that
		// predate any scan, so the DEFAULT identity row spans the
		// group's whole event history.
		bounds := store.ReadFieldIndexBounds(groupName)
		var groupFirst, groupLast int64
		if selected["DEFAULT"] {
			if streams, serr := fetchAllLogStreams(store, groupName, ""); serr == nil {
				for _, st := range streams {
					if st.FirstEventTs != 0 && (groupFirst == 0 || st.FirstEventTs < groupFirst) {
						groupFirst = st.FirstEventTs
					}
					if st.LastEventTs > groupLast {
						groupLast = st.LastEventTs
					}
				}
			}
		}
		for _, c := range candidates {
			name := indexedFieldName(c.spec.Name)
			entry := &fieldIndexEntry{
				fieldIndexName: name,
				indexCategory:  c.category,
				indexType:      c.spec.Type,
				lastScanTime:   bounds.LastScanTs,
			}
			if c.singleGroup {
				entry.logGroupIdentifier = store.ARNBuilder().CloudWatch().LogGroup(groupName)
			}
			if logsstore.IndexIdentityFields[name] && c.category == "DEFAULT" {
				entry.firstEventTime, entry.lastEventTime = groupFirst, groupLast
			} else if r, ok := bounds.Fields[name]; ok {
				entry.firstEventTime, entry.lastEventTime = r.First, r.Last
			}
			entry.cursor = groupName + "\x00" + name + "\x00" + c.category
			entries = append(entries, entry)
		}
	}
	// The listing pages through the scoped token vocabulary: the
	// identifier set and category selection digest into the scope, so a
	// token minted by one listing repositions no other (the same field
	// name recurs across groups) and a typed token rejects.
	sortedIds := make([]string, len(identifiers))
	copy(sortedIds, identifiers)
	sort.Strings(sortedIds)
	sortedCats := make([]string, 0, len(selected))
	for category := range selected {
		sortedCats = append(sortedCats, category)
	}
	sort.Strings(sortedCats)
	scope := listingScope("fieldindexes", strings.Join(sortedIds, "\x1f"), strings.Join(sortedCats, "\x1f"))
	result, err := paginateScopedListing(scope, nextToken, entries, logsstore.DefaultDescribeLimit,
		func(e *fieldIndexEntry) string { return e.cursor })
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// --- FIELD_INDEX_POLICY account-policy selection ---

// validateFieldIndexPolicySelection enforces the account-level field
// index policy rules at PutAccountPolicy: the document carries the same
// field syntax, the criteria selects by LogGroupNamePrefix (an empty
// criteria is the account-wide policy), no two policies share or nest
// prefixes, the account-wide policy excludes prefix-based ones and vice
// versa, and at most FieldIndexPolicyPrefixQuota policies select by
// prefix.
func (s *LogsService) validateFieldIndexPolicySelection(region, policyName, policyDocument, selectionCriteria string) error {
	if _, err := parseIndexPolicyDocument(policyDocument); err != nil {
		return err
	}
	if strings.Contains(selectionCriteria, "DataSourceName") || strings.Contains(selectionCriteria, "DataSourceType") {
		return NewLogsError("InvalidParameterException",
			"DataSourceName and DataSourceType selection requires vended-logs data sources; this platform supports LogGroupNamePrefix selection only", 400)
	}
	prefix, ok := transformerPolicyPrefix(selectionCriteria)
	if !ok {
		return NewLogsError("InvalidParameterException",
			"The supported selectionCriteria filters for a FIELD_INDEX_POLICY are LogGroupNamePrefix by itself or DataSourceName and DataSourceType together", 400)
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	policies, err := store.ListAccountPolicies("FIELD_INDEX_POLICY", "")
	if err != nil {
		return mapStoreError(err)
	}
	prefixCount := 0
	for i := range policies {
		if policies[i].PolicyName == policyName {
			continue
		}
		existingPrefix, ok := transformerPolicyPrefix(policies[i].SelectionCriteria)
		if !ok {
			continue
		}
		if existingPrefix == "" {
			// "If you have an account-level index policy that has no
			// name prefixes and applies to all log groups, then no other
			// account-level index policy with log group name prefix
			// filters can be created."
			if prefix != "" {
				return NewLogsError("InvalidParameterException",
					"An account-wide field index policy already exists; a policy with a log group name prefix filter cannot be created", 400)
			}
			// "You can have one account-level field index policy that
			// applies to all log groups in the account" (PutAccountPolicy)
			// — a second account-wide policy under another name has no
			// defined effective resolution, and the platform would pick
			// between them by listing order.
			return NewLogsError("InvalidParameterException",
				"An account-wide field index policy already exists; only one account-level field index policy can apply to all log groups", 400)
		}
		prefixCount++
		if prefix == "" {
			return NewLogsError("InvalidParameterException",
				"An account-wide field index policy cannot be created while prefix-filtered field index policies exist", 400)
		}
		if prefix == existingPrefix ||
			strings.HasPrefix(existingPrefix, prefix) ||
			strings.HasPrefix(prefix, existingPrefix) {
			return NewLogsError("InvalidParameterException",
				"A field index policy for prefix "+existingPrefix+" already exists with an overlapping prefix", 400)
		}
	}
	if prefix != "" && prefixCount >= logsstore.FieldIndexPolicyPrefixQuota {
		return NewLogsError("LimitExceededException",
			fmt.Sprintf("As many as %d field index policies can use log group name prefix selection criteria", logsstore.FieldIndexPolicyPrefixQuota), 400)
	}
	return nil
}
