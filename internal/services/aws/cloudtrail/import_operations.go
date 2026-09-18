package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// StartImport starts a CloudTrail import of trail events from an S3 bucket
// to a destination event data store.
func (s *CloudTrailService) StartImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// The time bounds keep both wire forms and parse strictly in the
	// Core: an unparseable present value is an error, never a silently
	// dropped bound.
	in := StartImportInput{
		ImportID:          request.GetStringParam(req.Parameters, "ImportId"),
		Destinations:      parseDestinationsList(req.Parameters["Destinations"]),
		StartEventTimeStr: req.GetParam("StartEventTime"),
		StartEventTimeRaw: req.Parameters["StartEventTime"],
		EndEventTimeStr:   req.GetParam("EndEventTime"),
		EndEventTimeRaw:   req.Parameters["EndEventTime"],
	}
	sourceRaw, hasSource := req.Parameters["ImportSource"]
	in.ImportSourceRaw = sourceRaw
	in.ImportSourceProvided = hasSource

	return s.startImportCore(ctx, store, in)
}

// StopImport stops a specified import.
func (s *CloudTrailService) StopImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.stopImportCore(store, ImportIDInput{
		ImportID: request.GetStringParam(req.Parameters, "ImportId"),
	})
}

// GetImport returns information about a specific import.
func (s *CloudTrailService) GetImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.getImportCore(store, ImportIDInput{
		ImportID: request.GetStringParam(req.Parameters, "ImportId"),
	})
}

// ListImports returns information on imports, optionally filtered by
// Destination or ImportStatus.
func (s *CloudTrailService) ListImports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.listImportsCore(store, ListImportsInput{
		NextToken:    req.GetParam("NextToken"),
		MaxResults:   request.GetIntParam(req.Parameters, "MaxResults"),
		Destination:  request.GetStringParam(req.Parameters, "Destination"),
		ImportStatus: request.GetStringParam(req.Parameters, "ImportStatus"),
	})
}

// ListImportFailures returns a list of failures for the specified import.
func (s *CloudTrailService) ListImportFailures(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.listImportFailuresCore(store, ListImportFailuresInput{
		ImportID:   request.GetStringParam(req.Parameters, "ImportId"),
		NextToken:  req.GetParam("NextToken"),
		MaxResults: request.GetIntParam(req.Parameters, "MaxResults"),
	})
}
