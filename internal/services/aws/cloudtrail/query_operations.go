package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// StartQuery starts a CloudTrail Lake query.
func (s *CloudTrailService) StartQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.startQueryCore(store, StartQueryInput{
		QueryStatement: request.GetStringParam(req.Parameters, "QueryStatement"),
	})
}

// GetQueryResults retrieves the results of a CloudTrail Lake query.
func (s *CloudTrailService) GetQueryResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.getQueryResultsCore(store, GetQueryResultsInput{
		QueryID:         request.GetStringParam(req.Parameters, "QueryId"),
		MaxQueryResults: request.GetIntParam(req.Parameters, "MaxQueryResults"),
		NextToken:       request.GetStringParam(req.Parameters, "NextToken"),
	})
}

// DescribeQuery retrieves metadata about a CloudTrail Lake query.
func (s *CloudTrailService) DescribeQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.describeQueryCore(store, DescribeQueryInput{
		QueryID: request.GetStringParam(req.Parameters, "QueryId"),
	})
}

// CancelQuery cancels a running CloudTrail Lake query.
func (s *CloudTrailService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.cancelQueryCore(store, CancelQueryInput{
		QueryID: request.GetStringParam(req.Parameters, "QueryId"),
	})
}

// ListQueries lists queries for an event data store.
func (s *CloudTrailService) ListQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.listQueriesCore(store, ListQueriesInput{
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
		QueryStatus:    request.GetStringParam(req.Parameters, "QueryStatus"),
		MaxResults:     request.GetIntParam(req.Parameters, "MaxResults"),
		NextToken:      request.GetStringParam(req.Parameters, "NextToken"),
	})
}
