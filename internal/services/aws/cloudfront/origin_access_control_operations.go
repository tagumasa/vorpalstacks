package cloudfront

import (
	"context"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateOriginAccessControl creates an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateOriginAccessControl.html
func (s *CloudFrontService) CreateOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	oac, err := s.createOriginAccessControlCore(store, CreateOriginAccessControlInput{
		Config: parseOACConfig(req),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControl": formatOACResponse(oac),
		"Location":            oac.ARN,
	}, nil
}

// GetOriginAccessControl returns an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetOriginAccessControl.html
func (s *CloudFrontService) GetOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	oac, err := s.getOriginAccessControlCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControl": formatOACResponse(oac),
	}, nil
}

// GetOriginAccessControlConfig returns the configuration of an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetOriginAccessControlConfig.html
func (s *CloudFrontService) GetOriginAccessControlConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	oac, err := s.getOriginAccessControlCore(store, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControlConfig": buildOACConfigResponse(oac),
		"ETag":                      oac.ETag,
	}, nil
}

// UpdateOriginAccessControl updates an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_UpdateOriginAccessControl.html
func (s *CloudFrontService) UpdateOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	oac, err := s.updateOriginAccessControlCore(store, UpdateOriginAccessControlInput{
		Id:      request.GetStringParam(req.Parameters, "Id"),
		IfMatch: getIfMatch(req),
		Config:  parseOACConfig(req),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControl": formatOACResponse(oac),
	}, nil
}

// DeleteOriginAccessControl deletes an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DeleteOriginAccessControl.html
func (s *CloudFrontService) DeleteOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteOriginAccessControlCore(store,
		request.GetStringParam(req.Parameters, "Id"), getIfMatch(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListOriginAccessControls lists origin access controls.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListOriginAccessControls.html
func (s *CloudFrontService) ListOriginAccessControls(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listOriginAccessControlsCore(store, marker, request.GetIntParam(req.Parameters, "MaxItems"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Controls))
	for _, oac := range result.Controls {
		items = append(items, map[string]interface{}{
			"Id":                            oac.ID,
			"Name":                          oac.Name,
			"Description":                   oac.Description,
			"OriginAccessControlOriginType": oac.OriginAccessControlOriginType,
			"SigningBehavior":               oac.SigningBehavior,
			"SigningProtocol":               oac.SigningProtocol,
		})
	}

	oacList := map[string]interface{}{
		"Marker":      marker,
		"MaxItems":    result.EffectiveMaxItems,
		"IsTruncated": result.IsTruncated,
		"Quantity":    len(items),
		"Items":       protocol.XMLElements{ElementName: "OriginAccessControlSummary", Items: items},
	}
	if result.NextMarker != "" {
		oacList["NextMarker"] = result.NextMarker
	}
	return map[string]interface{}{"OriginAccessControlList": oacList}, nil
}

// parseOACConfig parses the flat origin access control request payload
// into the store configuration type.
func parseOACConfig(req *request.ParsedRequest) *cloudfrontstore.OriginAccessControlConfig {
	return &cloudfrontstore.OriginAccessControlConfig{
		Name:                          request.GetStringParam(req.Parameters, "Name"),
		Description:                   request.GetStringParam(req.Parameters, "Description"),
		OriginAccessControlOriginType: request.GetStringParam(req.Parameters, "OriginAccessControlOriginType"),
		SigningBehavior:               request.GetStringParam(req.Parameters, "SigningBehavior"),
		SigningProtocol:               request.GetStringParam(req.Parameters, "SigningProtocol"),
	}
}

// formatOACResponse renders an origin access control with its metadata.
func formatOACResponse(oac *cloudfrontstore.OriginAccessControl) map[string]interface{} {
	return map[string]interface{}{
		"Id":                        oac.ID,
		"ARN":                       oac.ARN,
		"ETag":                      oac.ETag,
		"OriginAccessControlConfig": buildOACConfigResponse(oac),
		"CreatedTime":               oac.CreatedAt.Format(time.RFC3339),
		"LastModifiedTime":          oac.LastModifiedAt.Format(time.RFC3339),
	}
}

func buildOACConfigResponse(oac *cloudfrontstore.OriginAccessControl) map[string]interface{} {
	return map[string]interface{}{
		"Name":                          oac.Name,
		"Description":                   oac.Description,
		"OriginAccessControlOriginType": oac.OriginAccessControlOriginType,
		"SigningBehavior":               oac.SigningBehavior,
		"SigningProtocol":               oac.SigningProtocol,
	}
}
