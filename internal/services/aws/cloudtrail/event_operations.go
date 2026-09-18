package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// selectorTrailName resolves the selector operations' single input member.
// The four selector request shapes declare TrailName only; the value may
// still be a trail name or a trail ARN (ResolveTrail accepts both).
func selectorTrailName(req *request.ParsedRequest) string {
	return req.GetParam("TrailName")
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
		"TrailARN": trail.TrailARN,
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

	return s.getInsightSelectorsCore(store, selectorTrailName(req))
}

// PutInsightSelectors configures insight selectors for a trail.
func (s *CloudTrailService) PutInsightSelectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.putInsightSelectorsCore(store, PutInsightSelectorsInput{
		TrailName:           selectorTrailName(req),
		EventDataStore:      req.GetParam("EventDataStore"),
		InsightsDestination: req.GetParam("InsightsDestination"),
		InsightSelectorsRaw: req.Parameters["InsightSelectors"],
	})
}
