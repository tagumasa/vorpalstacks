package cloudfront

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateOriginAccessControl creates an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_CreateOriginAccessControl.html
func (s *CloudFrontService) CreateOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
	}

	originType := request.GetStringParam(req.Parameters, "OriginAccessControlOriginType")
	signingBehavior := request.GetStringParam(req.Parameters, "SigningBehavior")
	signingProtocol := request.GetStringParam(req.Parameters, "SigningProtocol")

	if originType == "" {
		return nil, NewAPIError("InvalidArgument", "OriginAccessControlOriginType is required", 400)
	}
	if signingBehavior == "" {
		return nil, NewAPIError("InvalidArgument", "SigningBehavior is required", 400)
	}
	if signingProtocol == "" {
		return nil, NewAPIError("InvalidArgument", "SigningProtocol is required", 400)
	}

	if !isValidOriginAccessControlOriginType(originType) {
		return nil, NewAPIError("InvalidArgument", "Invalid OriginAccessControlOriginType. Must be one of: s3, mediastore, mediapackagev2, lambda", 400)
	}
	if !isValidSigningBehavior(signingBehavior) {
		return nil, NewAPIError("InvalidArgument", "Invalid SigningBehavior. Must be one of: always, never, no-override", 400)
	}
	if signingProtocol != "sigv4" {
		return nil, NewAPIError("InvalidArgument", "Invalid SigningProtocol. Must be: sigv4", 400)
	}

	config := &cloudfrontstore.OriginAccessControlConfig{
		Name:                          name,
		Description:                   request.GetStringParam(req.Parameters, "Description"),
		OriginAccessControlOriginType: originType,
		SigningBehavior:               signingBehavior,
		SigningProtocol:               signingProtocol,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.originAccessControls.GetByName(name)
	if existing != nil {
		return nil, NewAPIError("OriginAccessControlAlreadyExists", "Origin access control with this name already exists", 409)
	}

	oac, err := store.originAccessControls.Create(config)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControl": map[string]interface{}{
			"Id":                        oac.ID,
			"ARN":                       oac.ARN,
			"ETag":                      oac.ETag,
			"OriginAccessControlConfig": buildOACConfigResponse(oac),
			"CreatedTime":               oac.CreatedAt.Format(time.RFC3339),
			"LastModifiedTime":          oac.LastModifiedAt.Format(time.RFC3339),
		},
		"Location": oac.ARN,
	}, nil
}

// GetOriginAccessControl returns an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetOriginAccessControl.html
func (s *CloudFrontService) GetOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	oac, err := store.originAccessControls.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControl": map[string]interface{}{
			"Id":                        oac.ID,
			"ARN":                       oac.ARN,
			"ETag":                      oac.ETag,
			"OriginAccessControlConfig": buildOACConfigResponse(oac),
			"CreatedTime":               oac.CreatedAt.Format(time.RFC3339),
			"LastModifiedTime":          oac.LastModifiedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetOriginAccessControlConfig returns the configuration of an origin access control.
func (s *CloudFrontService) GetOriginAccessControlConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	oac, err := store.originAccessControls.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
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
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidArgument", "Name is required", 400)
	}

	originType := request.GetStringParam(req.Parameters, "OriginAccessControlOriginType")
	signingBehavior := request.GetStringParam(req.Parameters, "SigningBehavior")
	signingProtocol := request.GetStringParam(req.Parameters, "SigningProtocol")

	if originType == "" {
		return nil, NewAPIError("InvalidArgument", "OriginAccessControlOriginType is required", 400)
	}
	if signingBehavior == "" {
		return nil, NewAPIError("InvalidArgument", "SigningBehavior is required", 400)
	}
	if signingProtocol == "" {
		return nil, NewAPIError("InvalidArgument", "SigningProtocol is required", 400)
	}

	if !isValidOriginAccessControlOriginType(originType) {
		return nil, NewAPIError("InvalidArgument", "Invalid OriginAccessControlOriginType. Must be one of: s3, mediastore, mediapackagev2, lambda", 400)
	}
	if !isValidSigningBehavior(signingBehavior) {
		return nil, NewAPIError("InvalidArgument", "Invalid SigningBehavior. Must be one of: always, never, no-override", 400)
	}
	if signingProtocol != "sigv4" {
		return nil, NewAPIError("InvalidArgument", "Invalid SigningProtocol. Must be: sigv4", 400)
	}

	config := &cloudfrontstore.OriginAccessControlConfig{
		Name:                          name,
		Description:                   request.GetStringParam(req.Parameters, "Description"),
		OriginAccessControlOriginType: originType,
		SigningBehavior:               signingBehavior,
		SigningProtocol:               signingProtocol,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	oac, err := store.originAccessControls.Update(id, config)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return nil, err
	}

	return map[string]interface{}{
		"OriginAccessControl": map[string]interface{}{
			"Id":                        oac.ID,
			"ARN":                       oac.ARN,
			"ETag":                      oac.ETag,
			"OriginAccessControlConfig": buildOACConfigResponse(oac),
			"CreatedTime":               oac.CreatedAt.Format(time.RFC3339),
			"LastModifiedTime":          oac.LastModifiedAt.Format(time.RFC3339),
		},
	}, nil
}

// DeleteOriginAccessControl deletes an origin access control.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_DeleteOriginAccessControl.html
func (s *CloudFrontService) DeleteOriginAccessControl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidArgument", "Id is required", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	err = store.originAccessControls.Delete(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, NewAPIError("NoSuchOriginAccessControl", "Origin access control not found", 404)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListOriginAccessControls lists origin access controls.
// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_ListOriginAccessControls.html
func (s *CloudFrontService) ListOriginAccessControls(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems == 0 {
		maxItems = 100
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.originAccessControls.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.OriginAccessControls))
	for _, oac := range result.OriginAccessControls {
		items = append(items, map[string]interface{}{
			"Id":               oac.ID,
			"Name":             oac.Name,
			"Description":      oac.Description,
			"OriginType":       oac.OriginAccessControlOriginType,
			"SigningBehavior":  oac.SigningBehavior,
			"SigningProtocol":  oac.SigningProtocol,
			"CreatedTime":      oac.CreatedAt.Format(time.RFC3339),
			"LastModifiedTime": oac.LastModifiedAt.Format(time.RFC3339),
		})
	}

	return map[string]interface{}{
		"OriginAccessControlList": map[string]interface{}{
			"Marker":      marker,
			"MaxItems":    maxItems,
			"IsTruncated": result.IsTruncated,
			"Quantity":    len(items),
			"NextMarker":  result.NextMarker,
			"Items":       protocol.XMLElements{ElementName: "OriginAccessControlSummary", Items: items},
		},
	}, nil
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
func isValidOriginAccessControlOriginType(t string) bool {
	switch t {
	case "s3", "mediastore", "mediapackagev2", "lambda":
		return true
	default:
		return false
	}
}

func isValidSigningBehavior(b string) bool {
	switch b {
	case "always", "never", "no-override":
		return true
	default:
		return false
	}
}
