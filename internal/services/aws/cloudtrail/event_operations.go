package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// selectorTrailName resolves the trail selector shared by the event-selector
// operations: TrailName first, then the TrailArn fallback.
func selectorTrailName(req *request.ParsedRequest) string {
	trailName := req.GetParam("TrailName")
	if trailName == "" {
		trailName = req.GetParam("TrailArn")
	}
	return trailName
}

// LookupEvents looks up events in CloudTrail based on the specified lookup attributes.
func (s *CloudTrailService) LookupEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.lookupEventsCore(store, LookupEventsInput{
		StartTimeStr:     req.GetParam("StartTime"),
		StartTimeRaw:     req.Parameters["StartTime"],
		EndTimeStr:       req.GetParam("EndTime"),
		EndTimeRaw:       req.Parameters["EndTime"],
		NextToken:        req.GetParam("NextToken"),
		LookupAttributes: req.Parameters["LookupAttributes"],
		EventNames:       req.Parameters["EventNames"],
		Username:         req.GetParam("Username"),
		EventCategory:    req.GetParam("EventCategory"),
		MaxResults:       request.GetIntParam(req.Parameters, "MaxResults"),
	})
}

// ListPublicKeys retrieves the public keys for CloudTrail.
func (s *CloudTrailService) ListPublicKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.listPublicKeysCore(store, ListPublicKeysInput{
		StartTimeStr: req.GetParam("StartTime"),
		StartTimeRaw: req.Parameters["StartTime"],
		EndTimeStr:   req.GetParam("EndTime"),
		EndTimeRaw:   req.Parameters["EndTime"],
		NextToken:    req.GetParam("NextToken"),
	})
}

// GetEventSelectors retrieves the event selectors for a trail.
func (s *CloudTrailService) GetEventSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := s.resolveTrailCore(store, selectorTrailName(req))
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"TrailArn": trail.TrailARN,
	}
	if len(trail.EventSelectors) > 0 {
		resp["EventSelectors"] = formatEventSelectors(trail.EventSelectors)
	}
	if len(trail.AdvancedEventSelectors) > 0 {
		resp["AdvancedEventSelectors"] = formatAdvancedEventSelectors(trail.AdvancedEventSelectors)
	}
	return resp, nil
}

// PutEventSelectors configures event selectors for a trail.
func (s *CloudTrailService) PutEventSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var advancedRaw interface{}
	if aesRaw, ok := req.Parameters["AdvancedEventSelectors"]; ok && aesRaw != nil {
		advancedRaw = aesRaw
	}

	return s.putEventSelectorsCore(store, PutEventSelectorsInput{
		TrailName:                 selectorTrailName(req),
		EventSelectorsRaw:         req.Parameters["EventSelectors"],
		AdvancedEventSelectorsRaw: advancedRaw,
	})
}

// GetInsightSelectors retrieves the insight selectors for a trail.
func (s *CloudTrailService) GetInsightSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trail, err := s.resolveTrailCore(store, selectorTrailName(req))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TrailArn":         trail.TrailARN,
		"InsightSelectors": formatInsightSelectors(trail.InsightSelectors),
	}, nil
}

// PutInsightSelectors configures insight selectors for a trail.
func (s *CloudTrailService) PutInsightSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		TrailName:           selectorTrailName(req),
		InsightSelectorsRaw: req.Parameters["InsightSelectors"],
	})
}
