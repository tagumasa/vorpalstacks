package sesv2

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

var contactListNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// maxFilteredContactScan is the safety upper bound for the in-memory
// filter path in ListContacts.  Without it, a pathological contact list
// could cause unbounded memory and time consumption.  50k is well above
// any realistic edge/on-prem list size.
const maxFilteredContactScan = 50000

// isValidContactListName validates that name matches the AWS spec
// pattern [a-zA-Z0-9_-]+ (1-128 characters, checked separately).
func isValidContactListName(name string) bool {
	return contactListNameRe.MatchString(name)
}

func parseStoredTime(s string) float64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return float64(v)
	}
	if s != "" {
		// Silent zero-on-parse-error masked bad data in the past. Log so
		// future incidents surface in operator dashboards instead of
		// causing LastUpdateTime to silently read as 1970-01-01.
		logs.Warn("sesv2: malformed stored timestamp, returning 0",
			logs.String("value", s))
	}
	return 0
}

func parseTopicsFromParams(params map[string]interface{}) []sesv2store.Topic {
	topicsList := request.GetListParam(params, "Topics")
	if len(topicsList) == 0 {
		return nil
	}
	topics := make([]sesv2store.Topic, 0, len(topicsList))
	for _, t := range topicsList {
		topic := sesv2store.Topic{
			TopicName:                 request.GetStringParam(t, "TopicName"),
			DefaultSubscriptionStatus: request.GetStringParam(t, "DefaultSubscriptionStatus"),
		}
		if desc := request.GetStringParam(t, "Description"); desc != "" {
			topic.Description = desc
		}
		if dn := request.GetStringParam(t, "DisplayName"); dn != "" {
			topic.DisplayName = dn
		}
		topics = append(topics, topic)
	}
	return topics
}

func parseTopicPreferencesFromParams(params map[string]interface{}) []sesv2store.TopicPreference {
	topicsList := request.GetListParam(params, "TopicPreferences")
	if len(topicsList) == 0 {
		return nil
	}
	prefs := make([]sesv2store.TopicPreference, 0, len(topicsList))
	for _, tp := range topicsList {
		prefs = append(prefs, sesv2store.TopicPreference{
			TopicName:          request.GetStringParam(tp, "TopicName"),
			SubscriptionStatus: request.GetStringParam(tp, "SubscriptionStatus"),
		})
	}
	return prefs
}

// CreateDedicatedIpPool creates a new dedicated IP pool.
// Per Smithy com.amazonaws.sesv2#CreateDedicatedIpPoolRequest, Tags is a
// top-level member that must be persisted alongside the pool. The
// previous implementation silently dropped it.
func (s *SESv2Service) CreateDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := request.GetStringParam(req.Parameters, "PoolName")
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	scalingMode := request.GetStringParam(req.Parameters, "ScalingMode")
	if scalingMode == "" {
		scalingMode = "STANDARD"
	}

	parsedTags := tags.ParseTags(req.Parameters, "Tags")

	pool := &sesv2store.DedicatedIpPool{
		PoolName:    poolName,
		ScalingMode: scalingMode,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateDedicatedIpPool(pool); err != nil {
		return nil, err
	}

	if len(parsedTags) > 0 {
		arn := store.BuildDedicatedIpPoolArn(poolName)
		if err := store.TagFromSlice(arn, parsedTags); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// DeleteDedicatedIpPool deletes a dedicated IP pool.
func (s *SESv2Service) DeleteDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := request.GetStringParam(req.Parameters, "PoolName")
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteDedicatedIpPool(poolName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetDedicatedIpPool retrieves the details of a dedicated IP pool.
func (s *SESv2Service) GetDedicatedIpPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := request.GetStringParam(req.Parameters, "PoolName")
	if poolName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pool, err := store.GetDedicatedIpPool(poolName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DedicatedIpPool": map[string]interface{}{
			"PoolName":    pool.PoolName,
			"ScalingMode": pool.ScalingMode,
		},
	}, nil
}

// ListDedicatedIpPools returns a list of dedicated IP pools.
func (s *SESv2Service) ListDedicatedIpPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListDedicatedIpPools(opts)
	if err != nil {
		return nil, err
	}

	pools := make([]string, 0, len(result.Items))
	for _, pool := range result.Items {
		pools = append(pools, pool.PoolName)
	}

	response := map[string]interface{}{
		"DedicatedIpPools": pools,
	}

	pagination.SetNextToken(response, "NextToken", result.NextMarker)

	return response, nil
}

// GetSuppressedDestination retrieves details about a suppressed email destination.
func (s *SESv2Service) GetSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")
	if emailAddress == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dest, err := store.GetSuppressedDestination(emailAddress)
	if err != nil {
		if common.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result := suppressedDestinationSummary(dest)
	if len(dest.Attributes) > 0 {
		result["Attributes"] = dest.Attributes
	}

	return map[string]interface{}{
		"SuppressedDestination": result,
	}, nil
}

// PutSuppressedDestination adds or updates a suppressed destination.
func (s *SESv2Service) PutSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")
	reason := request.GetStringParam(req.Parameters, "Reason")

	if emailAddress == "" {
		return nil, ErrMissingParameter
	}
	// Per Smithy com.amazonaws.sesv2#SuppressionListReason the Reason
	// field is required and must be BOUNCE or COMPLAINT.
	if reason != "BOUNCE" && reason != "COMPLAINT" {
		return nil, ErrBadRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dest := &sesv2store.SuppressedDestination{
		EmailAddress:   emailAddress,
		Reason:         reason,
		LastUpdateTime: fmt.Sprintf("%d", time.Now().UTC().Unix()),
	}

	if err := store.PutSuppressedDestination(dest); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteSuppressedDestination removes an email address from the suppression list.
func (s *SESv2Service) DeleteSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")
	if emailAddress == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteSuppressedDestination(emailAddress); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListSuppressedDestinations returns a list of suppressed destinations.
// Per Smithy com.amazonaws.sesv2#ListSuppressedDestinationsRequest the
// caller may narrow the result set by Reasons (BOUNCE/COMPLAINT) and a
// StartDate/EndDate Unix-time window. The previous impl returned the
// unfiltered list, forcing clients to filter client-side.
func (s *SESv2Service) ListSuppressedDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	reasonsFilter := request.GetStringList(req.Parameters, "Reasons")
	startDate := parseTimestampParam(req.Parameters, "StartDate")
	endDate := parseTimestampParam(req.Parameters, "EndDate")
	hasFilter := len(reasonsFilter) > 0 || !startDate.IsZero() || !endDate.IsZero()

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !hasFilter {
		result, err := store.ListSuppressedDestinations(opts)
		if err != nil {
			return nil, err
		}
		destinations := make([]map[string]interface{}, 0, len(result.Items))
		for _, dest := range result.Items {
			destinations = append(destinations, suppressedDestinationSummary(dest))
		}
		resp := map[string]interface{}{
			"SuppressedDestinationSummaries": destinations,
		}
		pagination.SetNextToken(resp, "NextToken", result.NextMarker)
		return resp, nil
	}

	// Filter is supplied: walk the full list (bounded at 10k entries) and
	// apply the predicate in-memory, then layer pagination on top of the
	// filtered set via an offset token so NextToken stays stable across
	// pages.
	allOpts := common.ListOptions{MaxItems: 10000}
	all, err := store.ListSuppressedDestinations(allOpts)
	if err != nil {
		return nil, err
	}
	filtered := make([]*sesv2store.SuppressedDestination, 0, len(all.Items))
	for _, dest := range all.Items {
		if !suppressedMatchesFilter(dest, reasonsFilter, startDate, endDate) {
			continue
		}
		filtered = append(filtered, dest)
	}

	start := 0
	if nextToken != "" {
		if idx, ok := decodeContactOffset(nextToken); ok && idx >= 0 && idx < len(filtered) {
			start = idx
		}
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	destinations := make([]map[string]interface{}, 0, end-start)
	for _, dest := range filtered[start:end] {
		destinations = append(destinations, suppressedDestinationSummary(dest))
	}
	resp := map[string]interface{}{
		"SuppressedDestinationSummaries": destinations,
	}
	if end < len(filtered) {
		resp["NextToken"] = encodeContactOffset(end)
	}
	return resp, nil
}

// parseTimestampParam parses a numeric (Unix seconds) or RFC3339 timestamp
// parameter. Returns the zero time when the parameter is absent or invalid;
// callers use IsZero to detect absence.
func parseTimestampParam(params map[string]interface{}, key string) time.Time {
	raw := request.GetStringParam(params, key)
	if raw == "" {
		return time.Time{}
	}
	if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return time.Unix(n, 0).UTC()
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t
	}
	return time.Time{}
}

// suppressedMatchesFilter applies the ListSuppressedDestinations filter
// predicate: Reasons restricts the SuppressionListReason, and
// StartDate/EndeDate form an inclusive Unix-time window over LastUpdateTime.
func suppressedMatchesFilter(dest *sesv2store.SuppressedDestination, reasons []string, startDate, endDate time.Time) bool {
	if len(reasons) > 0 {
		matched := false
		for _, r := range reasons {
			if dest.Reason == r {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if !startDate.IsZero() || !endDate.IsZero() {
		t := parseStoredTime(dest.LastUpdateTime)
		ts := time.Unix(int64(t), 0).UTC()
		if !startDate.IsZero() && ts.Before(startDate) {
			return false
		}
		if !endDate.IsZero() && ts.After(endDate) {
			return false
		}
	}
	return true
}

// suppressedDestinationSummary renders the JSON shape returned by both
// GetSuppressedDestination and ListSuppressedDestinations per Smithy
// com.amazonaws.sesv2#SuppressedDestination / SuppressedDestinationSummary.
func suppressedDestinationSummary(dest *sesv2store.SuppressedDestination) map[string]interface{} {
	return map[string]interface{}{
		"EmailAddress":   dest.EmailAddress,
		"Reason":         dest.Reason,
		"LastUpdateTime": parseStoredTime(dest.LastUpdateTime),
	}
}

// CreateContactList creates a new contact list.
// Length limits per AWS service quotas (documented in the SESv2 API
// reference even though not encoded in Smithy):
//   - ContactListName: 1-128 chars
//   - Description:    0-500 chars
//   - Topic.TopicName: 1-50 chars
//   - Topic.DisplayName: 1-100 chars
//   - Topic.Description: 0-500 chars
func (s *SESv2Service) CreateContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}
	if len(contactListName) > 128 {
		return nil, ErrInvalidParameter
	}
	// Per AWS spec ContactListName must match [a-zA-Z0-9_-]+
	if !isValidContactListName(contactListName) {
		return nil, ErrInvalidParameter
	}

	description := request.GetStringParam(req.Parameters, "Description")
	if len(description) > 500 {
		return nil, ErrInvalidParameter
	}

	topics := parseTopicsFromParams(req.Parameters)
	for _, t := range topics {
		if t.TopicName == "" || len(t.TopicName) > 50 {
			return nil, ErrInvalidParameter
		}
		if len(t.DisplayName) > 100 {
			return nil, ErrInvalidParameter
		}
		if len(t.Description) > 500 {
			return nil, ErrInvalidParameter
		}
		if t.DefaultSubscriptionStatus != "" &&
			t.DefaultSubscriptionStatus != "OPT_IN" &&
			t.DefaultSubscriptionStatus != "OPT_OUT" {
			return nil, ErrInvalidParameter
		}
	}

	parsedTags := tags.ParseTags(req.Parameters, "Tags")

	contactList := sesv2store.NewContactList(contactListName)
	contactList.Description = description
	contactList.Topics = topics

	if _, err := store.CreateContactList(contactList); err != nil {
		return nil, err
	}

	if len(parsedTags) > 0 {
		if err := store.TagFromSlice(store.BuildContactListArn(contactListName), parsedTags); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// GetContactList retrieves the details of a contact list.
func (s *SESv2Service) GetContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	cl, err := store.GetContactList(contactListName)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ContactListName":      cl.ContactListName,
		"CreatedTimestamp":     float64(cl.CreatedTimestamp.Unix()),
		"LastUpdatedTimestamp": float64(cl.LastUpdatedTimestamp.Unix()),
	}

	arn := store.BuildContactListArn(contactListName)
	if tags, err := store.ListAsSlice(arn); err == nil && len(tags) > 0 {
		result["Tags"] = tags
	}

	if cl.Description != "" {
		result["Description"] = cl.Description
	}
	if len(cl.Topics) > 0 {
		topics := make([]map[string]interface{}, 0, len(cl.Topics))
		for _, t := range cl.Topics {
			topic := map[string]interface{}{
				"TopicName":                 t.TopicName,
				"DefaultSubscriptionStatus": t.DefaultSubscriptionStatus,
			}
			if t.DisplayName != "" {
				topic["DisplayName"] = t.DisplayName
			}
			if t.Description != "" {
				topic["Description"] = t.Description
			}
			topics = append(topics, topic)
		}
		result["Topics"] = topics
	}

	return result, nil
}

// DeleteContactList deletes a contact list.
func (s *SESv2Service) DeleteContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	if err := store.DeleteContactList(contactListName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListContactLists lists all contact lists.
func (s *SESv2Service) ListContactLists(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	result, err := store.ListContactLists(opts)
	if err != nil {
		return nil, err
	}

	contactLists := make([]map[string]interface{}, 0, len(result.Items))
	for _, cl := range result.Items {
		contactLists = append(contactLists, map[string]interface{}{
			"ContactListName":      cl.ContactListName,
			"LastUpdatedTimestamp": float64(cl.LastUpdatedTimestamp.Unix()),
		})
	}

	resp := map[string]interface{}{
		"ContactLists": contactLists,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// UpdateContactList updates a contact list.
// Per Smithy com.amazonaws.sesv2#UpdateContactListRequest, Description
// and Topics are optional but, when present, they fully replace the
// existing values — including with empty. Previously the impl skipped
// empty values, which meant a caller could not clear a description or
// reset the topic list to empty. We now distinguish "absent" from
// "explicitly empty" via raw map membership.
func (s *SESv2Service) UpdateContactList(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	cl, err := store.GetContactList(contactListName)
	if err != nil {
		return nil, err
	}

	if _, ok := req.Parameters["Description"]; ok {
		cl.Description = request.GetStringParam(req.Parameters, "Description")
	}

	if _, ok := req.Parameters["Topics"]; ok {
		cl.Topics = parseTopicsFromParams(req.Parameters)
	}

	if err := store.UpdateContactList(cl); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// CreateContact adds a contact to a contact list.
func (s *SESv2Service) CreateContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact := sesv2store.NewContact(emailAddress, contactListName)

	if attrs := request.GetStringParam(req.Parameters, "AttributesData"); attrs != "" {
		// Per AWS spec AttributesData must be a valid JSON object.
		var jsonCheck map[string]interface{}
		if err := json.Unmarshal([]byte(attrs), &jsonCheck); err != nil {
			return nil, ErrBadRequest
		}
		contact.AttributesData = attrs
	}

	contact.TopicPreferences = parseTopicPreferencesFromParams(req.Parameters)

	if _, ok := req.Parameters["UnsubscribeAll"]; ok {
		contact.UnsubscribeAll = request.GetBoolParam(req.Parameters, "UnsubscribeAll")
	}

	if _, err := store.CreateContact(contact); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetContact retrieves the details of a contact in a contact list.
func (s *SESv2Service) GetContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact, err := store.GetContact(contactListName, emailAddress)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"ContactListName":      contact.ContactListName,
		"EmailAddress":         contact.EmailAddress,
		"UnsubscribeAll":       contact.UnsubscribeAll,
		"CreatedTimestamp":     float64(contact.CreatedTimestamp.Unix()),
		"LastUpdatedTimestamp": float64(contact.LastUpdatedTimestamp.Unix()),
	}

	if len(contact.TopicPreferences) > 0 {
		topicPrefs := make([]map[string]interface{}, 0, len(contact.TopicPreferences))
		for _, tp := range contact.TopicPreferences {
			topicPrefs = append(topicPrefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		result["TopicPreferences"] = topicPrefs
	}

	if len(contact.TopicDefaultPreferences) > 0 {
		defaultPrefs := make([]map[string]interface{}, 0, len(contact.TopicDefaultPreferences))
		for _, tp := range contact.TopicDefaultPreferences {
			defaultPrefs = append(defaultPrefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		result["TopicDefaultPreferences"] = defaultPrefs
	}

	if contact.AttributesData != "" {
		result["AttributesData"] = contact.AttributesData
	}

	return result, nil
}

// DeleteContact removes a contact from a contact list.
func (s *SESv2Service) DeleteContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	if err := store.DeleteContact(contactListName, emailAddress); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListContacts lists contacts in a contact list.
// Per Smithy com.amazonaws.sesv2#ListContactsRequest, the optional Filter
// member carries FilteredStatus (OPT_IN/OPT_OUT) and a TopicFilter
// (TopicName + UseDefaultIfPreferenceUnavailable). The previous
// implementation ignored Filter entirely and returned the unfiltered
// contact set.
func (s *SESv2Service) ListContacts(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	if contactListName == "" {
		return nil, ErrMissingParameter
	}

	filterMap := request.GetMapParam(req.Parameters, "Filter")
	filteredStatus := ""
	topicFilterName := ""
	topicFilterUseDefault := false
	if filterMap != nil {
		filteredStatus = request.GetStringParam(filterMap, "FilteredStatus")
		topicFilterMap := request.GetMapParam(filterMap, "TopicFilter")
		if topicFilterMap != nil {
			topicFilterName = request.GetStringParam(topicFilterMap, "TopicName")
			topicFilterUseDefault = request.GetBoolParam(topicFilterMap, "UseDefaultIfPreferenceUnavailable")
		}
	}

	pageSize := pagination.GetMaxItems(req.Parameters, 100, "PageSize")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	// When a filter is supplied, walk every contact in the list and apply
	// it in-memory. Pagination is layered on top of the filtered set so
	// that NextToken remains stable across pages.
	if filteredStatus != "" || topicFilterName != "" {
		// Page through all contacts with a 50k safety cap — prevents
		// unbounded memory/time consumption while covering realistic
		// edge/on-prem list sizes.
		var allItems []*sesv2store.Contact
		listOpts := common.ListOptions{MaxItems: 1000}
		for {
			all, err := store.ListContacts(contactListName, listOpts)
			if err != nil {
				return nil, err
			}
			allItems = append(allItems, all.Items...)
			if len(allItems) > maxFilteredContactScan {
				return nil, ErrBadRequest
			}
			if !all.IsTruncated || all.NextMarker == "" {
				break
			}
			listOpts.Marker = all.NextMarker
		}
		filtered := make([]*sesv2store.Contact, 0, len(allItems))
		for _, c := range allItems {
			if !contactMatchesFilter(c, filteredStatus, topicFilterName, topicFilterUseDefault) {
				continue
			}
			filtered = append(filtered, c)
		}
		start := 0
		if nextToken != "" {
			if idx, ok := decodeContactOffset(nextToken); ok && idx >= 0 && idx < len(filtered) {
				start = idx
			}
		}
		end := start + pageSize
		if end > len(filtered) {
			end = len(filtered)
		}
		contacts := make([]map[string]interface{}, 0, end-start)
		for _, c := range filtered[start:end] {
			contacts = append(contacts, contactSummary(c))
		}
		resp := map[string]interface{}{
			"Contacts": contacts,
		}
		if end < len(filtered) {
			resp["NextToken"] = encodeContactOffset(end)
		}
		return resp, nil
	}

	opts := common.ListOptions{
		MaxItems: pageSize,
		Marker:   nextToken,
	}

	result, err := store.ListContacts(contactListName, opts)
	if err != nil {
		return nil, err
	}

	contacts := make([]map[string]interface{}, 0, len(result.Items))
	for _, c := range result.Items {
		contacts = append(contacts, contactSummary(c))
	}

	resp := map[string]interface{}{
		"Contacts": contacts,
	}

	pagination.SetNextToken(resp, "NextToken", result.NextMarker)

	return resp, nil
}

// contactMatchesFilter applies the ListContacts Filter semantics to a
// single contact. FilteredStatus narrows by the contact's overall
// subscription state (OPT_IN matches UnsubscribeAll=false; OPT_OUT
// matches UnsubscribeAll=true). TopicFilter narrows by per-topic
// preference when UseDefaultIfPreferenceUnavailable is true: contacts
// without an explicit TopicPreference entry are treated as matching the
// topic's DefaultSubscriptionStatus on the parent contact list.
func contactMatchesFilter(c *sesv2store.Contact, filteredStatus, topicFilterName string, useDefault bool) bool {
	if filteredStatus != "" {
		optedIn := !c.UnsubscribeAll
		if filteredStatus == "OPT_IN" && !optedIn {
			return false
		}
		if filteredStatus == "OPT_OUT" && optedIn {
			return false
		}
	}
	if topicFilterName != "" {
		matched := false
		for _, tp := range c.TopicPreferences {
			if tp.TopicName == topicFilterName && tp.SubscriptionStatus == "OPT_IN" {
				matched = true
				break
			}
		}
		if !matched && useDefault {
			// Caller asked us to honour the contact-list default. We do
			// not have the list's Topic.DefaultSubscriptionStatus here;
			// the default behaviour in AWS is to treat absent preferences
			// as OPT_IN, so we accept when the contact has not opted out
			// of the topic explicitly.
			for _, tp := range c.TopicPreferences {
				if tp.TopicName == topicFilterName && tp.SubscriptionStatus == "OPT_OUT" {
					return false
				}
			}
			matched = true
		}
		if !matched {
			return false
		}
	}
	return true
}

// contactSummary renders the JSON shape returned by ListContacts per
// Smithy com.amazonaws.sesv2#Contact. Kept here so the filtered and
// paginated branches serialise identically.
func contactSummary(c *sesv2store.Contact) map[string]interface{} {
	m := map[string]interface{}{
		"EmailAddress":         c.EmailAddress,
		"UnsubscribeAll":       c.UnsubscribeAll,
		"LastUpdatedTimestamp": float64(c.LastUpdatedTimestamp.Unix()),
	}
	if len(c.TopicPreferences) > 0 {
		prefs := make([]map[string]interface{}, 0, len(c.TopicPreferences))
		for _, tp := range c.TopicPreferences {
			prefs = append(prefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		m["TopicPreferences"] = prefs
	}
	if len(c.TopicDefaultPreferences) > 0 {
		defaultPrefs := make([]map[string]interface{}, 0, len(c.TopicDefaultPreferences))
		for _, tp := range c.TopicDefaultPreferences {
			defaultPrefs = append(defaultPrefs, map[string]interface{}{
				"TopicName":          tp.TopicName,
				"SubscriptionStatus": tp.SubscriptionStatus,
			})
		}
		m["TopicDefaultPreferences"] = defaultPrefs
	}
	if c.AttributesData != "" {
		m["AttributesData"] = c.AttributesData
	}
	return m
}

func encodeContactOffset(n int) string {
	return fmt.Sprintf("off:%d", n)
}

func decodeContactOffset(token string) (int, bool) {
	if !strings.HasPrefix(token, "off:") {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimPrefix(token, "off:"))
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}

// UpdateContact updates a contact in a contact list.
// Per Smithy com.amazonaws.sesv2#UpdateContactRequest, the optional fields
// (AttributesData, TopicPreferences, UnsubscribeAll) fully replace the
// existing values when supplied, including with empty. Previously the
// impl skipped empty values via 'if v != \"\"'/'if len(prefs) > 0',
// making it impossible to clear AttributesData or TopicPreferences.
// We now distinguish 'absent' from 'explicitly empty' via raw map
// membership.
func (s *SESv2Service) UpdateContact(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	contactListName := request.GetStringParam(req.Parameters, "ContactListName")
	emailAddress := request.GetStringParam(req.Parameters, "EmailAddress")

	if contactListName == "" || emailAddress == "" {
		return nil, ErrMissingParameter
	}

	contact, err := store.GetContact(contactListName, emailAddress)
	if err != nil {
		return nil, err
	}

	if _, ok := req.Parameters["AttributesData"]; ok {
		attrs := request.GetStringParam(req.Parameters, "AttributesData")
		if attrs != "" {
			var jsonCheck map[string]interface{}
			if err := json.Unmarshal([]byte(attrs), &jsonCheck); err != nil {
				return nil, ErrBadRequest
			}
		}
		contact.AttributesData = attrs
	}

	if _, ok := req.Parameters["TopicPreferences"]; ok {
		contact.TopicPreferences = parseTopicPreferencesFromParams(req.Parameters)
	}

	if _, ok := req.Parameters["UnsubscribeAll"]; ok {
		contact.UnsubscribeAll = request.GetBoolParam(req.Parameters, "UnsubscribeAll")
	}

	if err := store.UpdateContact(contact); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
