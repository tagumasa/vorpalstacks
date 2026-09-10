package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// importFailOnWarnings reads the failonwarnings query parameter: an
// import with warnings is rolled back when it is set to true.
func importFailOnWarnings(req *request.ParsedRequest) bool {
	v := request.GetStringParam(req.Parameters, "failonwarnings")
	return v == "true"
}

// ImportRestApi creates a REST API from an external API definition
// document (POST /restapis?mode=import).
func (s *APIGatewayService) ImportRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params, err := parseImportParameters(req.Parameters)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, _, err := s.importRestApiCore(stores, req.Body, params, importFailOnWarnings(req))
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toRestApiResponse(created), nil
}

// PutRestApi updates an existing API from an external API definition
// document, merging or overwriting (PUT /restapis/{restApiId}).
func (s *APIGatewayService) PutRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params, err := parseImportParameters(req.Parameters)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	updated, _, err := s.putRestApiCore(
		stores,
		getRestApiId(req),
		request.GetStringParam(req.Parameters, "mode"),
		req.Body,
		params,
		importFailOnWarnings(req),
	)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toRestApiResponse(updated), nil
}

// ImportApiKeys imports API keys from a CSV document
// (POST /apikeys?mode=import&format=csv).
func (s *APIGatewayService) ImportApiKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ids, warnings, err := s.importApiKeysCore(
		stores,
		req.Body,
		request.GetStringParam(req.Parameters, "format"),
		importFailOnWarnings(req),
	)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	resp := map[string]interface{}{"ids": ids}
	if len(warnings) > 0 {
		resp["warnings"] = warnings
	}
	return resp, nil
}
