package cloudtrail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// LookupEvents looks up events in CloudTrail based on the specified lookup attributes.
func (s *CloudTrailService) LookupEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	query := cloudtrailstore.NewEventQuery()

	if startTime := req.GetParam("StartTime"); startTime != "" {
		t, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		query.StartTime = &t
	} else if stRaw := req.Parameters["StartTime"]; stRaw != nil {
		if ts, ok := stRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			query.StartTime = &t
		}
	}
	if endTime := req.GetParam("EndTime"); endTime != "" {
		t, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		query.EndTime = &t
	} else if etRaw := req.Parameters["EndTime"]; etRaw != nil {
		if ts, ok := etRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			query.EndTime = &t
		}
	}

	if nextToken := req.GetParam("NextToken"); nextToken != "" {
		query.NextToken = nextToken
	}

	if lookupAttrsRaw := req.Parameters["LookupAttributes"]; lookupAttrsRaw != nil {
		if attrs, ok := lookupAttrsRaw.([]interface{}); ok {
			for _, attr := range attrs {
				if attrMap, ok := attr.(map[string]interface{}); ok {
					key, _ := attrMap["AttributeKey"].(string)
					value, _ := attrMap["AttributeValue"].(string)
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

	if eventNamesRaw := req.Parameters["EventNames"]; eventNamesRaw != nil {
		if arr, ok := eventNamesRaw.([]interface{}); ok {
			for _, name := range arr {
				if s, ok := name.(string); ok {
					query.EventNames = append(query.EventNames, s)
				}
			}
		}
	}

	if username := req.GetParam("Username"); username != "" {
		query.Username = username
	}

	if maxResults := request.GetIntParam(req.Parameters, "MaxResults"); maxResults > 0 {
		query.MaxResults = int32(maxResults)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
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

// ListPublicKeys retrieves the public keys for CloudTrail.
func (s *CloudTrailService) ListPublicKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var startTime, endTime *time.Time

	// Accept both RFC3339 string and Unix timestamp (float64) formats
	// for StartTime/EndTime, as the AWS SDK serialises time.Time as
	// a Unix timestamp number in JSON-RPC 1.1.
	if st := req.GetParam("StartTime"); st != "" {
		t, err := time.Parse(time.RFC3339, st)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		startTime = &t
	} else if stRaw := req.Parameters["StartTime"]; stRaw != nil {
		if ts, ok := stRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			startTime = &t
		}
	}
	if et := req.GetParam("EndTime"); et != "" {
		t, err := time.Parse(time.RFC3339, et)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		endTime = &t
	} else if etRaw := req.Parameters["EndTime"]; etRaw != nil {
		if ts, ok := etRaw.(float64); ok {
			t := time.Unix(int64(ts), 0).UTC()
			endTime = &t
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
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

	if nextToken := req.GetParam("NextToken"); nextToken != "" {
		result["NextToken"] = nextToken
	}

	return result, nil
}

// GetEventSelectors retrieves the event selectors for a trail.
func (s *CloudTrailService) GetEventSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	trailName := req.GetParam("TrailName")
	if trailName == "" {
		trailName = req.GetParam("TrailArn")
	}

	if trailName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trail *cloudtrailstore.Trail
	if strings.HasPrefix(trailName, "arn:") || strings.Contains(trailName, "/") {
		trail, err = store.GetTrailByARN(trailName)
	} else {
		trail, err = store.GetTrail(trailName)
	}
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	formattedSelectors := make([]map[string]interface{}, 0)
	for _, es := range trail.EventSelectors {
		formattedSelectors = append(formattedSelectors, map[string]interface{}{
			"ReadWriteType":                 es.ReadWriteType,
			"IncludeManagementEvents":       es.IncludeManagementEvents,
			"DataResources":                 es.DataResources,
			"ExcludeManagementEventSources": es.ExcludeManagementEventSources,
		})
	}

	return map[string]interface{}{
		"TrailArn":       trail.TrailARN,
		"EventSelectors": formattedSelectors,
	}, nil
}

// PutEventSelectors configures event selectors for a trail.
func (s *CloudTrailService) PutEventSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	trailName := req.GetParam("TrailName")
	if trailName == "" {
		trailName = req.GetParam("TrailArn")
	}

	if trailName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trail *cloudtrailstore.Trail
	if strings.HasPrefix(trailName, "arn:") || strings.Contains(trailName, "/") {
		trail, err = store.GetTrailByARN(trailName)
	} else {
		trail, err = store.GetTrail(trailName)
	}
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	selectorsRaw := req.Parameters["EventSelectors"]
	if selectorsRaw == nil {
		return nil, ErrInvalidParameter
	}

	var selectors []cloudtrailstore.EventSelector
	switch v := selectorsRaw.(type) {
	case []interface{}:
		for _, sel := range v {
			if sm, ok := sel.(map[string]interface{}); ok {
				es := cloudtrailstore.EventSelector{}
				if rwt, ok := sm["ReadWriteType"].(string); ok {
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

	if err := store.PutEventSelector(trail.Name, selectors); err != nil {
		return nil, s.mapStoreError(err)
	}

	var formattedSelectors []map[string]interface{}
	for _, es := range selectors {
		formattedSelectors = append(formattedSelectors, map[string]interface{}{
			"ReadWriteType":                 es.ReadWriteType,
			"IncludeManagementEvents":       es.IncludeManagementEvents,
			"DataResources":                 es.DataResources,
			"ExcludeManagementEventSources": es.ExcludeManagementEventSources,
		})
	}

	return map[string]interface{}{
		"TrailArn":       trail.TrailARN,
		"EventSelectors": formattedSelectors,
	}, nil
}

// GetInsightSelectors retrieves the insight selectors for a trail.
func (s *CloudTrailService) GetInsightSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	trailName := req.GetParam("TrailName")
	if trailName == "" {
		trailName = req.GetParam("TrailArn")
	}

	if trailName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trail *cloudtrailstore.Trail
	if strings.HasPrefix(trailName, "arn:") || strings.Contains(trailName, "/") {
		trail, err = store.GetTrailByARN(trailName)
	} else {
		trail, err = store.GetTrail(trailName)
	}
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	formattedSelectors := make([]map[string]interface{}, 0)
	for _, is := range trail.InsightSelectors {
		formattedSelectors = append(formattedSelectors, map[string]interface{}{
			"InsightType": is.InsightType,
		})
	}

	return map[string]interface{}{
		"TrailArn":         trail.TrailARN,
		"InsightSelectors": formattedSelectors,
	}, nil
}

// PutInsightSelectors configures insight selectors for a trail.
func (s *CloudTrailService) PutInsightSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	trailName := req.GetParam("TrailName")
	if trailName == "" {
		trailName = req.GetParam("TrailArn")
	}

	if trailName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var trail *cloudtrailstore.Trail
	if strings.HasPrefix(trailName, "arn:") || strings.Contains(trailName, "/") {
		trail, err = store.GetTrailByARN(trailName)
	} else {
		trail, err = store.GetTrail(trailName)
	}
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	selectorsRaw := req.Parameters["InsightSelectors"]
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
					is.InsightType = it
				}
				selectors = append(selectors, is)
			}
		}
	}

	if err := store.PutInsightSelectors(trail.Name, selectors); err != nil {
		return nil, s.mapStoreError(err)
	}

	var formattedSelectors []map[string]interface{}
	for _, is := range selectors {
		formattedSelectors = append(formattedSelectors, map[string]interface{}{
			"InsightType": is.InsightType,
		})
	}

	return map[string]interface{}{
		"TrailArn":         trail.TrailARN,
		"InsightSelectors": formattedSelectors,
	}, nil
}

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
		var parsed map[string]interface{}
		if json.Unmarshal([]byte(e.CloudTrailEvent), &parsed) == nil {
			if ec, ok := parsed["eventCategory"].(string); ok && ec != "" {
				result["EventCategory"] = ec
			}
		}
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
