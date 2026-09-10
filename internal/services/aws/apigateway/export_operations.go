package apigateway

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"vorpalstacks/internal/common/request"
)

// exportResponse streams a generated export document: the modelled reply
// carries the document as the payload and the content type in the
// Content-Type header.
type exportResponse struct {
	body        []byte
	contentType string
	stage       string
	exportType  string
}

func (r *exportResponse) GetStream() io.Reader { return bytes.NewReader(r.body) }

func (r *exportResponse) GetStreamHeaders() http.Header {
	h := http.Header{}
	h.Set("Content-Type", r.contentType)
	extension := "json"
	if r.contentType == "application/yaml" {
		extension = "yaml"
	}
	h.Set("Content-Disposition", "attachment; filename=\""+r.stage+"_"+r.exportType+"."+extension+"\"")
	return h
}

// GetExport exports a REST API for a stage as a Swagger or OpenAPI
// document.
func (s *APIGatewayService) GetExport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	stageName := stageNameOf(req)
	exportType := request.GetStringParam(req.Parameters, "exportType")
	if exportType == "" {
		exportType = getPathParam(req, "exportType")
	}
	// Accepts is bound to the Accept header on the wire, with the header
	// value carrying any media-type parameters verbatim.
	accepts := req.Headers.Get("Accept")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	body, contentType, err := s.getExportCore(stores, apiId, stageName, exportType, accepts)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return &exportResponse{body: body, contentType: contentType, stage: stageName, exportType: exportType}, nil
}

// GetModelTemplate generates a model's sample mapping template.
func (s *APIGatewayService) GetModelTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	modelName := request.GetStringParam(req.Parameters, "modelName")
	if modelName == "" {
		modelName = getPathParam(req, "modelName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	value, err := s.getModelTemplateCore(stores, apiId, modelName)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return map[string]interface{}{"value": value}, nil
}

// GetSdkType retrieves a static SDK type by identifier.
func (s *APIGatewayService) GetSdkType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "id")
	if id == "" {
		id = getPathParam(req, "id")
	}
	t, err := s.getSdkTypeCore(id)
	if err != nil {
		return nil, err
	}
	return s.toSdkTypeResponse(t), nil
}

// GetSdkTypes lists the static SDK type table.
func (s *APIGatewayService) GetSdkTypes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	page, next, err := s.listSdkTypesCore(limit, request.GetStringParam(req.Parameters, "position"))
	if err != nil {
		return nil, err
	}
	items := make([]interface{}, 0, len(page))
	for _, t := range page {
		items = append(items, s.toSdkTypeResponse(t))
	}
	resp := map[string]interface{}{"item": items}
	if next != "" {
		resp["position"] = next
	}
	return resp, nil
}

func (s *APIGatewayService) toSdkTypeResponse(t *sdkType) map[string]interface{} {
	resp := map[string]interface{}{
		"id":           t.Id,
		"friendlyName": t.FriendlyName,
		"description":  t.Description,
	}
	if len(t.Properties) > 0 {
		props := make([]interface{}, 0, len(t.Properties))
		for _, p := range t.Properties {
			prop := map[string]interface{}{
				"name":         p.Name,
				"friendlyName": p.FriendlyName,
				"description":  p.Description,
				"required":     p.Required,
			}
			if p.DefaultValue != "" {
				prop["defaultValue"] = p.DefaultValue
			}
			props = append(props, prop)
		}
		resp["configurationProperties"] = props
	}
	return resp
}
