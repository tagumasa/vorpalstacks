package cloudtrail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// LookupEventsInput carries the LookupEvents members. Time members keep both
// wire forms (RFC3339 string and Unix epoch number) because the JSON 1.1
// protocol serialises timestamps as epochs while query strings arrive as
// RFC3339 text. The model input surface is exactly StartTime, EndTime,
// LookupAttributes, EventCategory, MaxResults and NextToken — event-name and
// username filtering goes through LookupAttributes, never top-level members.
type LookupEventsInput struct {
	StartTimeStr     string
	StartTimeRaw     interface{}
	EndTimeStr       string
	EndTimeRaw       interface{}
	NextToken        string
	LookupAttributes interface{}
	EventCategory    string
	MaxResults       int
}

// ListPublicKeysInput carries the ListPublicKeys members in both wire time
// forms, mirroring LookupEventsInput.
type ListPublicKeysInput struct {
	StartTimeStr string
	StartTimeRaw interface{}
	EndTimeStr   string
	EndTimeRaw   interface{}
	NextToken    string
}

// PutEventSelectorsInput carries the raw selector lists for PutEventSelectors.
// Basic and advanced selectors are mutually exclusive: setting one clears the
// other (AWS spec).
type PutEventSelectorsInput struct {
	TrailName                 string
	EventSelectorsRaw         interface{}
	AdvancedEventSelectorsRaw interface{}
}

// PutInsightSelectorsInput carries the raw insight selector list for
// PutInsightSelectors. The trail form and the event-data-store form are
// mutually exclusive: TrailName on one side, EventDataStore with
// InsightsDestination on the other.
type PutInsightSelectorsInput struct {
	TrailName           string
	EventDataStore      string
	InsightsDestination string
	InsightSelectorsRaw interface{}
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// epochFromRaw converts a JSON-protocol epoch value to a UTC time. The
// awsJson1.1 wire form carries millisecond decimal precision (the SDKs'
// FormatEpochSeconds), so the fraction must be kept: integer truncation
// would floor the bound to the start of its second and exclude events
// recorded in that same partial second.
func epochFromRaw(raw interface{}) (time.Time, bool) {
	ts, ok := raw.(float64)
	if !ok {
		return time.Time{}, false
	}
	ms := int64(math.Round(ts * 1000))
	return time.Unix(0, ms*int64(time.Millisecond)).UTC(), true
}

// parseWireTime parses a request time member that may arrive in either
// wire form — an RFC3339 string (query-string form) or a Unix epoch number
// with millisecond precision (awsJson1.1 serialisation). A present value
// that matches neither form is an error: the caller rejects it instead of
// silently dropping the bound. A nil time with no error means the member
// was not provided.
func parseWireTime(strForm string, rawForm interface{}) (*time.Time, error) {
	if strForm != "" {
		t, err := time.Parse(time.RFC3339, strForm)
		if err != nil {
			return nil, newInvalidParameterException(
				"The time value is neither a valid timestamp nor a Unix epoch")
		}
		return &t, nil
	}
	if rawForm != nil {
		t, ok := epochFromRaw(rawForm)
		if !ok {
			return nil, newInvalidParameterException(
				"The time value is neither a valid timestamp nor a Unix epoch")
		}
		return &t, nil
	}
	return nil, nil
}

// lookupEventsCore is the single entry point for LookupEvents: it validates
// the time range (EndTime must not precede StartTime), the lookup
// attributes (a single item, key enum, value presence), the EventCategory
// enum and the MaxResults window, runs the store query, and returns the
// formatted response.
func (s *CloudTrailService) lookupEventsCore(store cloudtrailstore.CloudTrailStoreInterface, in LookupEventsInput) (map[string]interface{}, error) {
	query := cloudtrailstore.NewEventQuery()

	startTime, err := parseWireTime(in.StartTimeStr, in.StartTimeRaw)
	if err != nil {
		return nil, err
	}
	query.StartTime = startTime
	endTime, err := parseWireTime(in.EndTimeStr, in.EndTimeRaw)
	if err != nil {
		return nil, err
	}
	query.EndTime = endTime

	// "If the specified end time is before the specified start time, an
	// error is returned" (LookupEvents, StartTime/EndTime) — the model
	// declares InvalidTimeRangeException for exactly this.
	if query.StartTime != nil && query.EndTime != nil && query.EndTime.Before(*query.StartTime) {
		return nil, newInvalidTimeRangeException(
			"End time must be later than or equal to start time")
	}

	if in.NextToken != "" {
		query.NextToken = in.NextToken
	}

	if in.LookupAttributes != nil {
		if attrs, ok := in.LookupAttributes.([]interface{}); ok {
			// "Currently the list can contain only one item"
			// (LookupAttributes) — the model declares
			// InvalidLookupAttributesException.
			if len(attrs) > 1 {
				return nil, newInvalidLookupAttributesException(
					"The list of lookup attributes can contain only one item")
			}
			for _, attr := range attrs {
				if attrMap, ok := attr.(map[string]interface{}); ok {
					key, _ := attrMap["AttributeKey"].(string)
					value, _ := attrMap["AttributeValue"].(string)
					if key == "" {
						return nil, newInvalidLookupAttributesException(
							"AttributeKey is required")
					}
					if err := validateLookupAttributeKey(key); err != nil {
						return nil, err
					}
					if value == "" {
						return nil, newInvalidLookupAttributesException(
							"AttributeValue is required for key: " + key)
					}
					if len(value) > cloudtrailstore.MaxLookupAttributeValueLength {
						return nil, newInvalidLookupAttributesException(
							"AttributeValue must be 2000 characters or fewer")
					}
					switch key {
					case "EventName":
						query.EventNames = append(query.EventNames, value)
					case "Username":
						query.Username = value
					case "ResourceName":
						query.ResourceNames = append(query.ResourceNames, value)
					case "ResourceType":
						query.ResourceType = value
					case "EventSource":
						query.EventSource = value
					case "AccessKeyId":
						query.AccessKeyID = value
					case "EventId":
						query.EventID = value
					case "ReadOnly":
						query.ReadOnly = value
					}
				}
			}
		}
	}

	// The model's EventCategory enum for LookupEvents carries the single
	// value "insight"; any other category is rejected with the declared
	// InvalidEventCategoryException.
	if in.EventCategory != "" {
		if in.EventCategory != "insight" {
			return nil, newInvalidEventCategoryException(
				"EventCategory must be: insight")
		}
		query.EventCategory = in.EventCategory
	}

	if in.MaxResults != 0 {
		if in.MaxResults < 1 || in.MaxResults > cloudtrailstore.MaxLookupEventsResults {
			return nil, newInvalidMaxResultsException(
				fmt.Sprintf("MaxResults must be between 1 and %d", cloudtrailstore.MaxLookupEventsResults))
		}
		query.MaxResults = int32(in.MaxResults)
	}

	events, nextToken, err := store.LookupEvents(query)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	formattedEvents := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		formattedEvents = append(formattedEvents, s.formatEvent(e))
	}

	result := map[string]interface{}{
		"Events": formattedEvents,
	}

	if nextToken != "" {
		result["NextToken"] = nextToken
	}

	return result, nil
}

// listPublicKeysCore is the single entry point for ListPublicKeys.
func (s *CloudTrailService) listPublicKeysCore(store cloudtrailstore.CloudTrailStoreInterface, in ListPublicKeysInput) (map[string]interface{}, error) {
	// Both members take the dual wire form — RFC3339 string or Unix
	// timestamp — the shared parser normalises (strict: a value present in
	// neither form is an invalid parameter).
	startTime, err := parseWireTime(in.StartTimeStr, in.StartTimeRaw)
	if err != nil {
		return nil, err
	}
	endTime, err := parseWireTime(in.EndTimeStr, in.EndTimeRaw)
	if err != nil {
		return nil, err
	}

	// Both bounds default to the current time — StartTime: "If not
	// specified, the current time is used, and the current public key is
	// returned"; EndTime: "If not specified, the current time is used" —
	// so the default listing returns exactly the keys valid now, and the
	// overlap filter still honours an explicit past-or-future window.
	now := time.Now().UTC()
	if startTime == nil {
		startTime = &now
	}
	if endTime == nil {
		endTime = &now
	}
	if startTime.After(*endTime) {
		return nil, newInvalidTimeRangeException(
			"Start time must be earlier than or equal to end time")
	}

	keys, err := store.ListPublicKeys(startTime, endTime)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	publicKeyList := make([]map[string]interface{}, 0, len(keys))
	for _, pk := range keys {
		// Map the key's fingerprint to the SDK's Fingerprint field; the SDK
		// type uses "Fingerprint" as the identifier, not "PublicKeyId".
		publicKeyList = append(publicKeyList, map[string]interface{}{
			"Fingerprint":       pk.Fingerprint(),
			"Value":             base64.StdEncoding.EncodeToString(pk.Value),
			"ValidityStartTime": pk.ValidityStartTime.Unix(),
			"ValidityEndTime":   pk.ValidityEndTime.Unix(),
		})
	}

	// The full key list is returned in one response, so no continuation
	// token is produced; echoing the caller's token back would keep a
	// paginating client from ever terminating.
	result := map[string]interface{}{
		"PublicKeyList": publicKeyList,
	}

	return result, nil
}

// putEventSelectorsCore is the single entry point for PutEventSelectors: it
// resolves the trail, validates and parses the selector lists, persists the
// provided selector type, and returns the formatted response. Basic and
// advanced selectors are mutually exclusive in one request — "You can use
// either AdvancedEventSelectors or EventSelectors, but not both"
// (PutEventSelectors) — and applying one form on a later call overwrites the
// stored form of the other.
func (s *CloudTrailService) putEventSelectorsCore(store cloudtrailstore.CloudTrailStoreInterface, in PutEventSelectorsInput) (map[string]interface{}, error) {
	trail, err := s.resolveTrailCore(store, in.TrailName)
	if err != nil {
		return nil, err
	}

	// "You can use either AdvancedEventSelectors or EventSelectors, but not
	// both" — the combination is rejected, and the error shape the operation
	// declares for an invalid selector configuration carries it.
	basicProvided := in.EventSelectorsRaw != nil
	advancedProvided := in.AdvancedEventSelectorsRaw != nil
	if basicProvided && advancedProvided {
		return nil, newInvalidEventSelectorsException(
			"You can use either EventSelectors or AdvancedEventSelectors, but not both")
	}
	if !basicProvided && !advancedProvided {
		return nil, newInvalidEventSelectorsException(
			"Specify a valid number of event selectors (1 to 5) for a trail")
	}

	var selectors []cloudtrailstore.EventSelector
	if basicProvided {
		var parseErr error
		selectors, parseErr = parseEventSelectors(in.EventSelectorsRaw)
		if parseErr != nil {
			return nil, parseErr
		}
		if len(selectors) < 1 || len(selectors) > cloudtrailstore.MaxEventSelectorsPerTrail {
			return nil, newInvalidEventSelectorsException(
				"Specify a valid number of event selectors (1 to 5) for a trail")
		}
		totalDataResources := 0
		for _, es := range selectors {
			totalDataResources += len(es.DataResources)
			for _, dr := range es.DataResources {
				if !basicDataResourceTypes[dr.Type] {
					return nil, newInvalidEventSelectorsException(
						"Event selectors can log data events only for AWS::DynamoDB::Table, AWS::Lambda::Function, and AWS::S3::Object resources")
				}
			}
		}
		if totalDataResources > cloudtrailstore.MaxDataResourcesPerTrail {
			return nil, newInvalidEventSelectorsException(
				"The total number of data resources cannot exceed 250 across all event selectors for a trail")
		}
	}

	var advancedSelectors []cloudtrailstore.AdvancedEventSelector
	if advancedProvided {
		var parseErr error
		advancedSelectors, parseErr = parseAdvancedEventSelectors(in.AdvancedEventSelectorsRaw)
		if parseErr != nil {
			return nil, parseErr
		}
		if err := validateAdvancedEventSelectors(advancedSelectors); err != nil {
			return nil, err
		}
	}

	// Persist selectors. Only the provided type is written, so a later call
	// of the other form overwrites the stored one exactly as the reference
	// describes.
	if advancedProvided {
		if err := store.PutAdvancedEventSelectors(trail.Name, advancedSelectors); err != nil {
			return nil, s.mapStoreError(err)
		}
	} else {
		if err := store.PutEventSelector(trail.Name, selectors); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	resp := map[string]interface{}{
		"TrailARN": trail.TrailARN,
	}
	if len(selectors) > 0 {
		resp["EventSelectors"] = formatEventSelectors(selectors)
	}
	if len(advancedSelectors) > 0 {
		resp["AdvancedEventSelectors"] = formatAdvancedEventSelectors(advancedSelectors)
	}
	return resp, nil
}

// validateAdvancedEventSelectors enforces the documented bounds for
// advanced selectors on a trail or an event data store: at most five
// selectors, at most 500 values across all conditions in all selectors
// — "Advanced event selectors — 500 conditions across all advanced event
// selectors. If a trail or event data store uses advanced event selectors,
// a maximum of 500 total values for all conditions in all selectors is
// allowed" (Quotas in AWS CloudTrail) — and, per field selector, the
// documented field vocabulary and each field's operator restrictions.
func validateAdvancedEventSelectors(selectors []cloudtrailstore.AdvancedEventSelector) error {
	if len(selectors) < 1 || len(selectors) > cloudtrailstore.MaxAdvancedEventSelectors {
		return newInvalidEventSelectorsException(
			"Specify a valid number of advanced event selectors (1 to 5)")
	}
	totalValues := 0
	for _, sel := range selectors {
		for _, fs := range sel.FieldSelectors {
			totalValues += len(fs.Equals) + len(fs.StartsWith) + len(fs.EndsWith) +
				len(fs.NotEquals) + len(fs.NotStartsWith) + len(fs.NotEndsWith)
			if err := validateAdvancedFieldSelector(fs); err != nil {
				return err
			}
		}
		if cloudtrailstore.AdvancedSelectorSelectsNetworkActivity(sel) {
			if err := validateNetworkActivityEventSource(sel); err != nil {
				return err
			}
		}
	}
	if totalValues > cloudtrailstore.MaxAdvancedSelectorValues {
		return newInvalidEventSelectorsException(
			"A maximum of 500 values for all conditions in all advanced event selectors is allowed")
	}
	return nil
}

// validateNetworkActivityEventSource enforces the network-activity
// eventSource rule the AdvancedFieldSelector reference states: within a
// selector whose eventCategory is NetworkActivity, eventSource is required
// and accepts the Equals operator only ("For network activity events,
// this is a required field that only uses the Equals operator"). Outside
// network activity the field is unrestricted, so the rule is
// selector-scoped rather than a vocabulary entry.
func validateNetworkActivityEventSource(sel cloudtrailstore.AdvancedEventSelector) error {
	for _, fs := range sel.FieldSelectors {
		if !strings.EqualFold(fs.Field, "eventSource") {
			continue
		}
		if len(fs.NotEquals) > 0 || len(fs.StartsWith) > 0 || len(fs.NotStartsWith) > 0 ||
			len(fs.EndsWith) > 0 || len(fs.NotEndsWith) > 0 {
			return newInvalidEventSelectorsException(
				`Advanced event selector field "eventSource" accepts only the Equals operators for network activity events`)
		}
		return nil
	}
	return newInvalidEventSelectorsException(
		`Advanced event selector field "eventSource" is required for network activity events`)
}

// validateAdvancedFieldSelector enforces the field vocabulary and each
// field's documented operator restrictions (AdvancedFieldSelector API
// reference): eventCategory, readOnly, resources.type and errorCode
// accept Equals only, sessionCredentialFromConsole accepts Equals and
// NotEquals only. Any other operator on a restricted field, or a field
// outside the vocabulary, is an invalid selector.
func validateAdvancedFieldSelector(fs cloudtrailstore.AdvancedFieldSelector) error {
	restrictedTo, known := cloudtrailstore.AdvancedSelectorFieldRules(fs.Field)
	if !known {
		return newInvalidEventSelectorsException(
			fmt.Sprintf("Advanced event selector field %q is not a valid field", fs.Field))
	}
	// Every operator list is bounded at one item minimum in the model, and
	// a field selector carrying no operator value at all expresses no
	// filter: the evaluation semantics would exclude every event while
	// the API answered valid. Both forms are rejected at write.
	if len(fs.Equals) == 0 && len(fs.NotEquals) == 0 &&
		len(fs.StartsWith) == 0 && len(fs.NotStartsWith) == 0 &&
		len(fs.EndsWith) == 0 && len(fs.NotEndsWith) == 0 {
		return newInvalidEventSelectorsException(
			fmt.Sprintf("Advanced event selector field %q must carry at least one operator value", fs.Field))
	}
	if restrictedTo == nil {
		return nil
	}
	for _, list := range []struct {
		name   string
		values []string
	}{
		{"Equals", fs.Equals},
		{"NotEquals", fs.NotEquals},
		{"StartsWith", fs.StartsWith},
		{"NotStartsWith", fs.NotStartsWith},
		{"EndsWith", fs.EndsWith},
		{"NotEndsWith", fs.NotEndsWith},
	} {
		if len(list.values) > 0 && !restrictedTo[list.name] {
			return newInvalidEventSelectorsException(
				fmt.Sprintf("Advanced event selector field %q accepts only the %s operators",
					fs.Field, permittedOperatorList(restrictedTo)))
		}
	}
	return nil
}

// permittedOperatorList renders the operator-list names a restricted field
// allows, in the fixed evaluation order.
func permittedOperatorList(restricted map[string]bool) string {
	all := []string{"Equals", "NotEquals", "StartsWith", "NotStartsWith", "EndsWith", "NotEndsWith"}
	permitted := make([]string, 0, len(restricted))
	for _, name := range all {
		if restricted[name] {
			permitted = append(permitted, name)
		}
	}
	return strings.Join(permitted, ", ")
}

// parseEventSelectors parses and validates the basic selector list from its
// raw wire value.
func parseEventSelectors(raw interface{}) ([]cloudtrailstore.EventSelector, error) {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil, newInvalidEventSelectorsException(
			"EventSelectors must be a list")
	}
	selectors := make([]cloudtrailstore.EventSelector, 0, len(arr))
	for _, sel := range arr {
		sm, ok := sel.(map[string]interface{})
		if !ok {
			return nil, newInvalidEventSelectorsException(
				"Each event selector must be a map")
		}
		es := cloudtrailstore.EventSelector{}
		if rwt, ok := sm["ReadWriteType"].(string); ok {
			if err := validateReadWriteType(rwt); err != nil {
				return nil, err
			}
			es.ReadWriteType = rwt
		}
		if ime, ok := sm["IncludeManagementEvents"].(bool); ok {
			es.IncludeManagementEvents = ime
		}
		if drRaw, ok := sm["DataResources"].([]interface{}); ok {
			for _, dr := range drRaw {
				if drm, ok := dr.(map[string]interface{}); ok {
					drItem := cloudtrailstore.DataResource{}
					if t, ok := drm["Type"].(string); ok {
						drItem.Type = t
					}
					if vals, ok := drm["Values"].([]interface{}); ok {
						for _, val := range vals {
							if s, ok := val.(string); ok {
								drItem.Values = append(drItem.Values, s)
							}
						}
					}
					es.DataResources = append(es.DataResources, drItem)
				}
			}
		}
		if emesRaw, ok := sm["ExcludeManagementEventSources"].([]interface{}); ok {
			for _, emes := range emesRaw {
				if s, ok := emes.(string); ok {
					es.ExcludeManagementEventSources = append(es.ExcludeManagementEventSources, s)
				}
			}
		}
		selectors = append(selectors, es)
	}
	return selectors, nil
}

// getInsightSelectorsCore is the single entry point for
// GetInsightSelectors: a trail without insight selectors has Insights
// events not enabled, and the operation answers with its declared
// InsightNotEnabledException rather than an empty selector list.
func (s *CloudTrailService) getInsightSelectorsCore(store cloudtrailstore.CloudTrailStoreInterface, trailName string) (map[string]interface{}, error) {
	trail, err := s.resolveTrailCore(store, trailName)
	if err != nil {
		return nil, err
	}

	if len(trail.InsightSelectors) == 0 {
		return nil, newInsightNotEnabledException(
			"Insights are not enabled on the trail")
	}

	return map[string]interface{}{
		"TrailARN":         trail.TrailARN,
		"InsightSelectors": formatInsightSelectors(trail.InsightSelectors),
	}, nil
}

// putInsightSelectorsCore is the single entry point for PutInsightSelectors.
// The parameter combination is validated first — TrailName on one side,
// EventDataStore with InsightsDestination on the other — then every
// selector must carry a valid InsightType.
func (s *CloudTrailService) putInsightSelectorsCore(store cloudtrailstore.CloudTrailStoreInterface, in PutInsightSelectorsInput) (map[string]interface{}, error) {
	// "You cannot use this parameter with the EventDataStore and
	// InsightsDestination parameters" (TrailName); "To enable Insights on an
	// event data store, you must provide both the EventDataStore and
	// InsightsDestination parameters" — the model declares
	// InvalidParameterCombinationException for the broken combinations.
	hasTrailForm := in.TrailName != ""
	hasEDSForm := in.EventDataStore != "" || in.InsightsDestination != ""
	if hasTrailForm && hasEDSForm {
		return nil, newInvalidParameterCombinationException(
			"TrailName cannot be used with EventDataStore and InsightsDestination")
	}
	if hasEDSForm && (in.EventDataStore == "" || in.InsightsDestination == "") {
		return nil, newInvalidParameterCombinationException(
			"EventDataStore and InsightsDestination must be provided together")
	}
	if !hasTrailForm && !hasEDSForm {
		return nil, newInvalidParameterCombinationException(
			"Either TrailName or EventDataStore with InsightsDestination is required")
	}
	// Insight selectors on an event data store (the source/destination
	// form) belong to the event-data-store insight machinery and are not
	// available until that machinery exists.
	if hasEDSForm {
		return nil, newUnsupportedOperationException(
			"Insight selectors on an event data store are not supported")
	}

	trail, err := s.resolveTrailCore(store, in.TrailName)
	if err != nil {
		return nil, err
	}

	selectorsRaw := in.InsightSelectorsRaw
	if selectorsRaw == nil {
		return nil, newInvalidInsightSelectorsException(
			"InsightSelectors is required")
	}

	var selectors []cloudtrailstore.InsightSelector
	switch v := selectorsRaw.(type) {
	case []interface{}:
		for _, sel := range v {
			sm, ok := sel.(map[string]interface{})
			if !ok {
				return nil, newInvalidInsightSelectorsException(
					"Each insight selector must be a map")
			}
			is := cloudtrailstore.InsightSelector{}
			it, hasType := sm["InsightType"].(string)
			if !hasType || it == "" {
				return nil, newInvalidInsightSelectorsException(
					"InsightType is required for each insight selector")
			}
			if err := validateInsightType(it); err != nil {
				return nil, err
			}
			is.InsightType = it
			selectors = append(selectors, is)
		}
	default:
		return nil, newInvalidInsightSelectorsException(
			"InsightSelectors must be a list")
	}
	if len(selectors) < 1 {
		return nil, newInvalidInsightSelectorsException(
			"InsightSelectors must contain at least one selector")
	}

	if err := store.PutInsightSelectors(trail.Name, selectors); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"TrailARN":         trail.TrailARN,
		"InsightSelectors": formatInsightSelectors(selectors),
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// formatEventSelectors converts event selectors to API response format.
func formatEventSelectors(selectors []cloudtrailstore.EventSelector) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(selectors))
	for _, es := range selectors {
		result = append(result, map[string]interface{}{
			"ReadWriteType":                 es.ReadWriteType,
			"IncludeManagementEvents":       es.IncludeManagementEvents,
			"DataResources":                 es.DataResources,
			"ExcludeManagementEventSources": es.ExcludeManagementEventSources,
		})
	}
	return result
}

// formatInsightSelectors converts insight selectors to API response format.
func formatInsightSelectors(selectors []cloudtrailstore.InsightSelector) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(selectors))
	for _, is := range selectors {
		result = append(result, map[string]interface{}{
			"InsightType": is.InsightType,
		})
	}
	return result
}

// formatAdvancedEventSelectors converts advanced event selectors to API
// response format.
func formatAdvancedEventSelectors(selectors []cloudtrailstore.AdvancedEventSelector) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(selectors))
	for _, sel := range selectors {
		selMap := map[string]interface{}{}
		if sel.Name != "" {
			selMap["Name"] = sel.Name
		}
		fields := make([]map[string]interface{}, 0, len(sel.FieldSelectors))
		for _, fs := range sel.FieldSelectors {
			fm := map[string]interface{}{"Field": fs.Field}
			if len(fs.Equals) > 0 {
				fm["Equals"] = fs.Equals
			}
			if len(fs.StartsWith) > 0 {
				fm["StartsWith"] = fs.StartsWith
			}
			if len(fs.EndsWith) > 0 {
				fm["EndsWith"] = fs.EndsWith
			}
			if len(fs.NotEquals) > 0 {
				fm["NotEquals"] = fs.NotEquals
			}
			if len(fs.NotStartsWith) > 0 {
				fm["NotStartsWith"] = fs.NotStartsWith
			}
			if len(fs.NotEndsWith) > 0 {
				fm["NotEndsWith"] = fs.NotEndsWith
			}
			fields = append(fields, fm)
		}
		selMap["FieldSelectors"] = fields
		result = append(result, selMap)
	}
	return result
}

// parseAdvancedEventSelectors parses advanced event selectors from the raw
// wire value. The value may arrive as the JSON-protocol list form or as a
// JSON string holding the same list; both trails and event data stores share
// this single parser, so the accepted forms cannot drift between the two
// selector-bearing resource families. Each selector must carry at least one
// field selector (FieldSelectors is model-required), and every field
// selector must carry a non-empty Field matching the model's SelectorField
// pattern.
func parseAdvancedEventSelectors(raw interface{}) ([]cloudtrailstore.AdvancedEventSelector, error) {
	var rawList []interface{}
	switch v := raw.(type) {
	case []interface{}:
		rawList = v
	case string:
		if err := json.Unmarshal([]byte(v), &rawList); err != nil {
			return nil, newInvalidEventSelectorsException(
				"AdvancedEventSelectors must be a list")
		}
	default:
		return nil, newInvalidEventSelectorsException(
			"AdvancedEventSelectors must be a list")
	}

	result := make([]cloudtrailstore.AdvancedEventSelector, 0, len(rawList))
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, newInvalidEventSelectorsException(
				"Each advanced event selector must be a map")
		}
		sel := cloudtrailstore.AdvancedEventSelector{}
		if name, ok := m["Name"].(string); ok {
			sel.Name = name
		}
		fieldsRaw, hasFields := m["FieldSelectors"].([]interface{})
		if !hasFields || len(fieldsRaw) < 1 {
			return nil, newInvalidEventSelectorsException(
				"Each advanced event selector must contain at least one field selector")
		}
		for _, fRaw := range fieldsRaw {
			fm, ok := fRaw.(map[string]interface{})
			if !ok {
				return nil, newInvalidEventSelectorsException(
					"Each field selector must be a map")
			}
			fs := cloudtrailstore.AdvancedFieldSelector{}
			f, hasField := fm["Field"].(string)
			if !hasField || f == "" || !advancedFieldPattern.MatchString(f) {
				return nil, newInvalidEventSelectorsException(
					"Each field selector must carry a valid Field")
			}
			fs.Field = f
			fs.Equals = toStringSlice(fm["Equals"])
			fs.StartsWith = toStringSlice(fm["StartsWith"])
			fs.EndsWith = toStringSlice(fm["EndsWith"])
			fs.NotEquals = toStringSlice(fm["NotEquals"])
			fs.NotStartsWith = toStringSlice(fm["NotStartsWith"])
			fs.NotEndsWith = toStringSlice(fm["NotEndsWith"])
			sel.FieldSelectors = append(sel.FieldSelectors, fs)
		}
		result = append(result, sel)
	}
	return result, nil
}

// toStringSlice converts an interface to a string slice.
func toStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return val
	default:
		return nil
	}
}

// formatEvent converts a store event to the LookupEvents response format.
// The wire Event shape carries exactly the nine model members; the deeper
// record fields (user agent, error codes, source address, ...) live inside
// the CloudTrailEvent record JSON, never as top-level response keys.
func (s *CloudTrailService) formatEvent(e *cloudtrailstore.Event) map[string]interface{} {
	result := map[string]interface{}{
		"EventId":     e.EventID,
		"EventName":   e.EventName,
		"EventSource": e.EventSource,
		"EventTime":   e.EventTime.Unix(),
		"ReadOnly":    e.ReadOnly,
	}

	if e.AccessKeyId != "" {
		result["AccessKeyId"] = e.AccessKeyId
	}
	if e.UserIdentity != nil && e.UserIdentity.UserName != "" {
		result["Username"] = e.UserIdentity.UserName
	}
	if e.CloudTrailEvent != "" {
		result["CloudTrailEvent"] = e.CloudTrailEvent
	}

	if len(e.Resources) > 0 {
		var formattedResources []map[string]interface{}
		for _, r := range e.Resources {
			formattedResources = append(formattedResources, map[string]interface{}{
				"ResourceType": r.ResourceType,
				"ResourceName": r.ResourceName,
			})
		}
		result["Resources"] = formattedResources
	}

	return result
}
