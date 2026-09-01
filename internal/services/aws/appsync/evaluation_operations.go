package appsync

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
)

// ListTypesByAssociation returns types from the source API of a merged API association.
// GET /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}/types
func (s *AppSyncService) ListTypesByAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	types, nextToken, err := s.listTypesByAssociationCore(store,
		request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		request.GetStringParam(req.Parameters, "associationId"),
		request.GetStringParam(req.Parameters, "format"),
		request.GetIntParam(req.Parameters, "maxResults"),
		request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(types))
	for _, t := range types {
		items = append(items, typeToMap(t))
	}

	response := map[string]interface{}{
		"types": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// EvaluateCode executes an AppSync function resolver code snippet and returns the result.
func (s *AppSyncService) EvaluateCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var body struct {
		Code     string          `json:"code"`
		Context  json.RawMessage `json:"context"`
		Function string          `json:"function"`
		Runtime  struct {
			Name           string `json:"runtimeName"`
			RuntimeVersion string `json:"runtimeVersion"`
		} `json:"runtime"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return nil, NewBadRequestException("invalid request body")
	}
	return evaluateCodeCore(&EvaluateCodeInput{
		Code:        body.Code,
		Context:     body.Context,
		Function:    body.Function,
		RuntimeName: body.Runtime.Name,
	})
}

// EvaluateMappingTemplate evaluates a VTL mapping template.
// POST /v1/dataplane-evaluatetemplate
func (s *AppSyncService) EvaluateMappingTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var body struct {
		Context  string `json:"context"`
		Template string `json:"template"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return nil, NewBadRequestException("invalid request body")
	}
	return evaluateMappingTemplateCore(&EvaluateMappingTemplateInput{
		Context:  body.Context,
		Template: body.Template,
	})
}
