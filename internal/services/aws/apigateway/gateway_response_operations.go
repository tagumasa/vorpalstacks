package apigateway

import (
	"context"
	"net/http"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/apigateway"
)

// acceptedResponse is the 202-with-empty-body reply of the cache flush
// operations, per their modelled http traits.
type acceptedResponse struct{}

// GetStreamStatusCode reports 202 Accepted.
func (acceptedResponse) GetStreamStatusCode() int { return http.StatusAccepted }

// renderGatewayResponse renders a GatewayResponse.
func (s *APIGatewayService) renderGatewayResponse(r *apigateway.GatewayResponse) map[string]interface{} {
	resp := map[string]interface{}{
		"responseType":    r.ResponseType,
		"defaultResponse": r.DefaultResponse,
	}
	if r.StatusCode != "" {
		resp["statusCode"] = r.StatusCode
	}
	if len(r.ResponseParameters) > 0 {
		resp["responseParameters"] = r.ResponseParameters
	}
	if len(r.ResponseTemplates) > 0 {
		resp["responseTemplates"] = r.ResponseTemplates
	}
	return resp
}

// GetGatewayResponse retrieves a gateway response by its response type.
func (s *APIGatewayService) GetGatewayResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	responseType := request.GetStringParam(req.Parameters, "responseType")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resp, err := s.getGatewayResponseCore(stores, apiId, responseType)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.renderGatewayResponse(resp), nil
}

// GetGatewayResponses lists the API's gateway responses.
func (s *APIGatewayService) GetGatewayResponses(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listGatewayResponsesCore(stores, limit, position, apiId)
	if err != nil {
		return nil, err
	}
	items := make([]interface{}, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, s.renderGatewayResponse(item))
	}
	resp := map[string]interface{}{"item": items}
	if result.IsTruncated {
		resp["position"] = result.NextMarker
	}
	return resp, nil
}

// PutGatewayResponse creates or replaces a gateway response.
func (s *APIGatewayService) PutGatewayResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	responseType := request.GetStringParam(req.Parameters, "responseType")
	in := &GatewayResponseInput{
		StatusCode: request.GetStringParam(req.Parameters, "statusCode"),
	}
	if params, ok := req.Parameters["responseParameters"].(map[string]interface{}); ok {
		in.ResponseParameters = make(map[string]string, len(params))
		for k, v := range params {
			if vs, ok := v.(string); ok {
				in.ResponseParameters[k] = vs
			}
		}
	}
	if templates, ok := req.Parameters["responseTemplates"].(map[string]interface{}); ok {
		in.ResponseTemplates = make(map[string]string, len(templates))
		for k, v := range templates {
			if vs, ok := v.(string); ok {
				in.ResponseTemplates[k] = vs
			}
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.putGatewayResponseCore(stores, apiId, responseType, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.renderGatewayResponse(created), nil
}

// DeleteGatewayResponse removes a gateway response.
func (s *APIGatewayService) DeleteGatewayResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	responseType := request.GetStringParam(req.Parameters, "responseType")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteGatewayResponseCore(stores, apiId, responseType); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// UpdateGatewayResponse patches a gateway response.
func (s *APIGatewayService) UpdateGatewayResponse(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	responseType := request.GetStringParam(req.Parameters, "responseType")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	updated, err := s.updateGatewayResponseCore(stores, apiId, responseType, ops)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.renderGatewayResponse(updated), nil
}

// flushStageCacheCore resolves the flush precondition: the stage must
// exist. The API reference documents no cache-cluster precondition for the
// flush operations, so existence is the whole check.
func (s *APIGatewayService) flushStageCacheCore(stores *apiGatewayStores, apiId, stageName string) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if _, err := s.getStageCore(stores, apiId, stageName); err != nil {
		return err
	}
	return nil
}

// FlushStageCache flushes a stage's cache; the modelled reply is 202 with
// an empty body.
func (s *APIGatewayService) FlushStageCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.flushStageCacheCore(stores, getRestApiId(req), stageNameOf(req)); err != nil {
		return nil, toApiGatewayError(err)
	}
	return acceptedResponse{}, nil
}

// FlushStageAuthorizersCache flushes a stage's authorizer cache; the
// modelled reply is 202 with an empty body.
func (s *APIGatewayService) FlushStageAuthorizersCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.flushStageCacheCore(stores, getRestApiId(req), stageNameOf(req)); err != nil {
		return nil, toApiGatewayError(err)
	}
	return acceptedResponse{}, nil
}

// stageNameOf reads the stage name from the query parameters or the path.
func stageNameOf(req *request.ParsedRequest) string {
	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}
	return stageName
}
