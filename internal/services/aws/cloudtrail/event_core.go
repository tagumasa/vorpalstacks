package cloudtrail

import (
	"encoding/base64"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// LookupEventsInput carries the LookupEvents members. Time members keep both
// wire forms (RFC3339 string and Unix epoch number) because the JSON 1.1
// protocol serialises timestamps as epochs while query strings arrive as
// RFC3339 text.
type LookupEventsInput struct {
	StartTimeStr     string
	StartTimeRaw     interface{}
	EndTimeStr       string
	EndTimeRaw       interface{}
	NextToken        string
	LookupAttributes interface{}
	EventNames       interface{}
	Username         string
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
// PutInsightSelectors.
type PutInsightSelectorsInput struct {
	TrailName           string
	InsightSelectorsRaw interface{}
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// lookupEventsCore is the single entry point for LookupEvents: it validates
// the lookup attributes (count limit, key enum, value presence) and the
// MaxResults window, runs the store query, and returns the formatted
// response.
func (s *CloudTrailService) lookupEventsCore(store cloudtrailstore.CloudTrailStoreInterface, in LookupEventsInput) (map[string]interface{}, error) {
	query := cloudtrailstore.NewEventQuery()

	if in.StartTimeStr != "" {
		t, err := time.Parse(time.RFC3339, in.StartTimeStr)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		query.StartTime = &t
	} else if in.StartTimeRaw != nil {
		if ts, ok := in.StartTimeRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			query.StartTime = &t
		}
	}
	if in.EndTimeStr != "" {
		t, err := time.Parse(time.RFC3339, in.EndTimeStr)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		query.EndTime = &t
	} else if in.EndTimeRaw != nil {
		if ts, ok := in.EndTimeRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			query.EndTime = &t
		}
	}

	if in.NextToken != "" {
		query.NextToken = in.NextToken
	}

	if in.LookupAttributes != nil {
		if attrs, ok := in.LookupAttributes.([]interface{}); ok {
			if len(attrs) > 50 {
				return nil, awserrors.NewAWSError("InvalidLookupAttributesException",
					"Number of lookup attributes exceeds the limit of 50", 400)
			}
			for _, attr := range attrs {
				if attrMap, ok := attr.(map[string]interface{}); ok {
					key, _ := attrMap["AttributeKey"].(string)
					value, _ := attrMap["AttributeValue"].(string)
					if key == "" {
						return nil, awserrors.NewAWSError("InvalidLookupAttributesException",
							"AttributeKey is required", 400)
					}
					if err := validateLookupAttributeKey(key); err != nil {
						return nil, err
					}
					if value == "" {
						return nil, awserrors.NewAWSError("InvalidLookupAttributesException",
							"AttributeValue is required for key: "+key, 400)
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

	if in.EventNames != nil {
		if arr, ok := in.EventNames.([]interface{}); ok {
			for _, name := range arr {
				if s, ok := name.(string); ok {
					query.EventNames = append(query.EventNames, s)
				}
			}
		}
	}

	if in.Username != "" {
		query.Username = in.Username
	}

	if in.EventCategory != "" {
		query.EventCategory = in.EventCategory
	}

	if in.MaxResults > 0 {
		if in.MaxResults > 50 {
			return nil, awserrors.NewAWSError("InvalidMaxResultsException",
				"MaxResults exceeds the maximum of 50", 400)
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
	var startTime, endTime *time.Time

	// Accept both RFC3339 string and Unix timestamp (float64) formats
	// for StartTime/EndTime, as the AWS SDK serialises time.Time as
	// a Unix timestamp number in JSON-RPC 1.1.
	if in.StartTimeStr != "" {
		t, err := time.Parse(time.RFC3339, in.StartTimeStr)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		startTime = &t
	} else if in.StartTimeRaw != nil {
		if ts, ok := in.StartTimeRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			startTime = &t
		}
	}
	if in.EndTimeStr != "" {
		t, err := time.Parse(time.RFC3339, in.EndTimeStr)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		endTime = &t
	} else if in.EndTimeRaw != nil {
		if ts, ok := in.EndTimeRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			endTime = &t
		}
	}

	keys, err := store.ListPublicKeys(startTime, endTime)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	publicKeyList := make([]map[string]interface{}, 0, len(keys))
	for _, pk := range keys {
		// Map internal PublicKeyID to SDK's Fingerprint field; the SDK type
		// uses "Fingerprint" as the identifier, not "PublicKeyId".
		publicKeyList = append(publicKeyList, map[string]interface{}{
			"Fingerprint":       pk.PublicKeyID,
			"Value":             base64.StdEncoding.EncodeToString(pk.Value),
			"ValidityStartTime": pk.ValidityStartTime.Unix(),
			"ValidityEndTime":   pk.ValidityEndTime.Unix(),
		})
	}

	result := map[string]interface{}{
		"PublicKeyList": publicKeyList,
	}

	if in.NextToken != "" {
		result["NextToken"] = in.NextToken
	}

	return result, nil
}

// putEventSelectorsCore is the single entry point for PutEventSelectors: it
// resolves the trail, validates and parses the selector lists, persists the
// provided selector type (basic and advanced are mutually exclusive — setting
// one clears the other), and returns the formatted response.
func (s *CloudTrailService) putEventSelectorsCore(store cloudtrailstore.CloudTrailStoreInterface, in PutEventSelectorsInput) (map[string]interface{}, error) {
	trail, err := s.resolveTrailCore(store, in.TrailName)
	if err != nil {
		return nil, err
	}

	// Parse AdvancedEventSelectors if provided.
	var advancedSelectors []cloudtrailstore.AdvancedEventSelector
	if in.AdvancedEventSelectorsRaw != nil {
		advancedSelectors = parseTrailAdvancedEventSelectors(in.AdvancedEventSelectorsRaw)
	}

	// Parse basic EventSelectors if provided.
	selectorsRaw := in.EventSelectorsRaw
	if selectorsRaw == nil && len(advancedSelectors) == 0 {
		return nil, ErrInvalidParameter
	}

	var selectors []cloudtrailstore.EventSelector
	if selectorsRaw != nil {
		switch v := selectorsRaw.(type) {
		case []interface{}:
			for _, sel := range v {
				if sm, ok := sel.(map[string]interface{}); ok {
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
			}
		}
	}

	// Persist selectors. Basic and advanced are mutually exclusive: setting
	// one clears the other (AWS spec). Only write the type that was provided.
	if len(advancedSelectors) > 0 {
		if err := store.PutAdvancedEventSelectors(trail.Name, advancedSelectors); err != nil {
			return nil, s.mapStoreError(err)
		}
	} else if selectorsRaw != nil {
		if err := store.PutEventSelector(trail.Name, selectors); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	resp := map[string]interface{}{
		"TrailArn": trail.TrailARN,
	}
	if len(selectors) > 0 {
		resp["EventSelectors"] = formatEventSelectors(selectors)
	}
	if len(advancedSelectors) > 0 {
		resp["AdvancedEventSelectors"] = formatAdvancedEventSelectors(advancedSelectors)
	}
	return resp, nil
}

// putInsightSelectorsCore is the single entry point for PutInsightSelectors.
func (s *CloudTrailService) putInsightSelectorsCore(store cloudtrailstore.CloudTrailStoreInterface, in PutInsightSelectorsInput) (map[string]interface{}, error) {
	trail, err := s.resolveTrailCore(store, in.TrailName)
	if err != nil {
		return nil, err
	}

	selectorsRaw := in.InsightSelectorsRaw
	if selectorsRaw == nil {
		return nil, ErrInvalidParameter
	}

	var selectors []cloudtrailstore.InsightSelector
	switch v := selectorsRaw.(type) {
	case []interface{}:
		for _, sel := range v {
			if sm, ok := sel.(map[string]interface{}); ok {
				is := cloudtrailstore.InsightSelector{}
				if it, ok := sm["InsightType"].(string); ok {
					if err := validateInsightType(it); err != nil {
						return nil, err
					}
					is.InsightType = it
				}
				selectors = append(selectors, is)
			}
		}
	}

	if err := store.PutInsightSelectors(trail.Name, selectors); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"TrailArn":         trail.TrailARN,
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

// parseTrailAdvancedEventSelectors parses advanced event selectors from the
// raw wire value for trail event selector configuration.
func parseTrailAdvancedEventSelectors(raw interface{}) []cloudtrailstore.AdvancedEventSelector {
	var rawList []interface{}
	switch v := raw.(type) {
	case []interface{}:
		rawList = v
	default:
		return nil
	}

	result := make([]cloudtrailstore.AdvancedEventSelector, 0, len(rawList))
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		sel := cloudtrailstore.AdvancedEventSelector{}
		if name, ok := m["Name"].(string); ok {
			sel.Name = name
		}
		if fieldsRaw, ok := m["FieldSelectors"].([]interface{}); ok {
			for _, fRaw := range fieldsRaw {
				fm, ok := fRaw.(map[string]interface{})
				if !ok {
					continue
				}
				fs := cloudtrailstore.AdvancedFieldSelector{}
				if f, ok := fm["Field"].(string); ok {
					fs.Field = f
				}
				fs.Equals = toStringSlice(fm["Equals"])
				fs.StartsWith = toStringSlice(fm["StartsWith"])
				fs.EndsWith = toStringSlice(fm["EndsWith"])
				fs.NotEquals = toStringSlice(fm["NotEquals"])
				fs.NotStartsWith = toStringSlice(fm["NotStartsWith"])
				fs.NotEndsWith = toStringSlice(fm["NotEndsWith"])
				sel.FieldSelectors = append(sel.FieldSelectors, fs)
			}
		}
		result = append(result, sel)
	}
	return result
}

// formatEvent converts a store event to the LookupEvents response format.
func (s *CloudTrailService) formatEvent(e *cloudtrailstore.Event) map[string]interface{} {
	result := map[string]interface{}{
		"EventId":       e.EventID,
		"EventName":     e.EventName,
		"EventSource":   e.EventSource,
		"EventTime":     e.EventTime.Unix(),
		"ReadOnly":      e.ReadOnly,
		"EventVersion":  e.EventVersion,
		"EventType":     e.EventType,
		"EventCategory": e.EventCategory,
	}

	if e.AccessKeyId != "" {
		result["AccessKeyId"] = e.AccessKeyId
	}
	if e.UserIdentity != nil {
		result["Username"] = e.UserIdentity.UserName
		if e.UserIdentity.AccountID != "" {
			result["AccountId"] = e.UserIdentity.AccountID
		}
	}
	if e.RequestID != "" {
		result["RequestId"] = e.RequestID
	}
	if e.SourceIPAddress != "" {
		result["SourceIpAddress"] = e.SourceIPAddress
	}
	if e.UserAgent != "" {
		result["UserAgent"] = e.UserAgent
	}
	if e.ErrorCode != "" {
		result["ErrorCode"] = e.ErrorCode
	}
	if e.ErrorMessage != "" {
		result["ErrorMessage"] = e.ErrorMessage
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
