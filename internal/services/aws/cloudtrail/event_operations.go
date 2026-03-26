package cloudtrail

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// LookupEvents looks up events in CloudTrail based on the specified lookup attributes.
func (s *CloudTrailService) LookupEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	query := cloudtrailstore.NewEventQuery()

	if startTime := getParam(req, "StartTime"); startTime != "" {
		t, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		query.StartTime = &t
	}
	if endTime := getParam(req, "EndTime"); endTime != "" {
		t, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		query.EndTime = &t
	}

	if nextToken := getParam(req, "NextToken"); nextToken != "" {
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

	if username := getParam(req, "Username"); username != "" {
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
	now := time.Now()
	return map[string]interface{}{
		"PublicKeyList": []map[string]interface{}{
			{
				"Value":      "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0EXAMPLE",
				"ValidFrom":  now.Add(-24 * time.Hour).Format(time.RFC3339),
				"ValidUntil": now.Add(24 * time.Hour).Format(time.RFC3339),
			},
		},
	}, nil
}

// GetEventSelectors retrieves the event selectors for a trail.
func (s *CloudTrailService) GetEventSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	trailName := getParam(req, "TrailName")
	if trailName == "" {
		trailName = getParam(req, "TrailArn")
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
	trailName := getParam(req, "TrailName")
	if trailName == "" {
		trailName = getParam(req, "TrailArn")
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
	trailName := getParam(req, "TrailName")
	if trailName == "" {
		trailName = getParam(req, "TrailArn")
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
	trailName := getParam(req, "TrailName")
	if trailName == "" {
		trailName = getParam(req, "TrailArn")
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
		"EventId":     e.EventID,
		"EventName":   e.EventName,
		"EventSource": e.EventSource,
		"EventTime":   e.EventTime.Unix(),
		"ReadOnly":    e.ReadOnly,
	}

	if e.UserIdentity != nil {
		result["Username"] = e.UserIdentity.UserName
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
