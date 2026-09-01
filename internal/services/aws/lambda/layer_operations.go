// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/utils/timeutils"
)

// PublishLayerVersion publishes a new version of a Lambda layer.
// Creates the layer if it does not exist, and publishes a new version with the provided content.
func (s *LambdaService) PublishLayerVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &LayerVersionCreateInput{
		LayerName:   request.GetStringParam(req.Parameters, "LayerName"),
		Description: request.GetStringParam(req.Parameters, "Description"),
		LicenseInfo: request.GetStringParam(req.Parameters, "LicenseInfo"),
		Content:     request.GetMapParam(req.Parameters, "Content"),
		Region:      reqCtx.GetRegion(),
	}
	if compats, ok := req.Parameters["CompatibleRuntimes"].([]interface{}); ok {
		for _, c := range compats {
			if str, ok := c.(string); ok {
				in.CompatibleRuntimes = append(in.CompatibleRuntimes, str)
			}
		}
	}
	if compats, ok := req.Parameters["CompatibleArchitectures"].([]interface{}); ok {
		for _, c := range compats {
			if str, ok := c.(string); ok {
				in.CompatibleArchitectures = append(in.CompatibleArchitectures, str)
			}
		}
	}

	layer, created, err := s.publishLayerVersionCore(ctx, reqCtx, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Content": map[string]interface{}{
			"Location":   created.CodeLocation,
			"CodeSha256": created.CodeSha256,
			"CodeSize":   created.CodeSize,
		},
		"LayerArn":                layer.LayerArn,
		"LayerVersionArn":         created.LayerVersionArn,
		"Description":             created.Description,
		"CreatedDate":             created.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		"Version":                 created.Version,
		"CompatibleRuntimes":      created.CompatibleRuntimes,
		"LicenseInfo":             created.LicenseInfo,
		"CompatibleArchitectures": created.CompatibleArchitectures,
	}, nil
}

// DeleteLayerVersion deletes a specific version of a Lambda layer.
func (s *LambdaService) DeleteLayerVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")

	versionNumber := int64(request.GetIntParam(req.Parameters, "VersionNumber"))

	if err := s.deleteLayerVersionCore(reqCtx, layerName, versionNumber); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetLayerVersion retrieves information about a specific version of a Lambda layer.
func (s *LambdaService) GetLayerVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")
	versionNumber := int64(request.GetIntParam(req.Parameters, "VersionNumber"))

	layer, layerVersion, err := s.getLayerVersionCore(reqCtx, layerName, versionNumber)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Content": map[string]interface{}{
			"Location":   layerVersion.CodeLocation,
			"CodeSha256": layerVersion.CodeSha256,
			"CodeSize":   layerVersion.CodeSize,
		},
		"LayerArn":                layer.LayerArn,
		"LayerVersionArn":         layerVersion.LayerVersionArn,
		"Description":             layerVersion.Description,
		"CreatedDate":             layerVersion.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		"Version":                 layerVersion.Version,
		"CompatibleRuntimes":      layerVersion.CompatibleRuntimes,
		"LicenseInfo":             layerVersion.LicenseInfo,
		"CompatibleArchitectures": layerVersion.CompatibleArchitectures,
	}, nil
}

// ListLayers lists the Lambda layers in the account, with optional filtering by runtime.
func (s *LambdaService) ListLayers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := validateMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))
	marker := request.GetStringParam(req.Parameters, "Marker")
	compatibleRuntime := request.GetStringParam(req.Parameters, "CompatibleRuntime")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, isTruncated, err := s.listLayersCore(store, compatibleRuntime, marker, maxItems)
	if err != nil {
		return nil, err
	}

	layers := make([]interface{}, 0)
	for _, l := range items {
		layer := map[string]interface{}{
			"LayerName":   l.LayerName,
			"LayerArn":    l.LayerArn,
			"CreatedDate": l.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		}
		if l.LatestMatchingVersion != nil {
			layer["LatestMatchingVersion"] = map[string]interface{}{
				"Version":            l.LatestMatchingVersion.Version,
				"LayerVersionArn":    l.LatestMatchingVersion.LayerVersionArn,
				"Description":        l.LatestMatchingVersion.Description,
				"CreatedDate":        l.LatestMatchingVersion.CreatedDate.Format(timeutils.ISO8601UTCFormat),
				"CompatibleRuntimes": l.LatestMatchingVersion.CompatibleRuntimes,
				"LicenseInfo":        l.LatestMatchingVersion.LicenseInfo,
			}
		}
		layers = append(layers, layer)
	}

	resp := map[string]interface{}{
		"Layers": layers,
	}
	if isTruncated {
		resp["NextMarker"] = nextMarker
	}

	return resp, nil
}

// ListLayerVersions lists all versions of a Lambda layer.
func (s *LambdaService) ListLayerVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")

	maxItems := validateMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))
	marker := request.GetStringParam(req.Parameters, "Marker")

	versions, nextMarker, isTruncated, err := s.listLayerVersionsCore(reqCtx, layerName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0)
	for _, v := range versions {
		items = append(items, map[string]interface{}{
			"LayerVersionArn":         v.LayerVersionArn,
			"Version":                 v.Version,
			"Description":             v.Description,
			"CreatedDate":             v.CreatedDate.Format(timeutils.ISO8601UTCFormat),
			"CompatibleRuntimes":      v.CompatibleRuntimes,
			"LicenseInfo":             v.LicenseInfo,
			"CompatibleArchitectures": v.CompatibleArchitectures,
		})
	}

	resp := map[string]interface{}{
		"LayerVersions": items,
	}
	if isTruncated {
		resp["NextMarker"] = nextMarker
	}

	return resp, nil
}

// GetLayerVersionByArn retrieves a layer version by its full ARN.
func (s *LambdaService) GetLayerVersionByArn(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerVersionArn := request.GetStringParam(req.Parameters, "Arn")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	layerVersion, layerArn, err := s.getLayerVersionByArnCore(store, layerVersionArn)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Content": map[string]interface{}{
			"Location":   layerVersion.CodeLocation,
			"CodeSha256": layerVersion.CodeSha256,
			"CodeSize":   layerVersion.CodeSize,
		},
		"LayerArn":                layerArn,
		"LayerVersionArn":         layerVersion.LayerVersionArn,
		"Description":             layerVersion.Description,
		"CreatedDate":             layerVersion.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		"Version":                 layerVersion.Version,
		"CompatibleRuntimes":      layerVersion.CompatibleRuntimes,
		"LicenseInfo":             layerVersion.LicenseInfo,
		"CompatibleArchitectures": layerVersion.CompatibleArchitectures,
	}, nil
}

// AddLayerVersionPermission adds a permission to a layer version's
// resource-based policy, allowing other accounts to use the layer.
func (s *LambdaService) AddLayerVersionPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &LayerPermissionInput{
		LayerName:     request.GetStringParam(req.Parameters, "LayerName"),
		VersionNumber: int64(request.GetIntParam(req.Parameters, "VersionNumber")),
		StatementId:   request.GetStringParam(req.Parameters, "StatementId"),
		Action:        request.GetStringParam(req.Parameters, "Action"),
		Principal:     request.GetStringParam(req.Parameters, "Principal"),
	}

	targetVersion, err := s.addLayerVersionPermissionCore(reqCtx, in)
	if err != nil {
		return nil, err
	}

	statement := map[string]interface{}{
		"Sid":       in.StatementId,
		"Effect":    "Allow",
		"Principal": in.Principal,
		"Action":    in.Action,
		"Resource":  targetVersion.LayerVersionArn,
	}

	statementJSON, err := json.Marshal(statement)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Statement":  string(statementJSON),
		"RevisionId": targetVersion.RevisionId,
	}, nil
}

// RemoveLayerVersionPermission removes a permission from a layer version's
// resource-based policy.
func (s *LambdaService) RemoveLayerVersionPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")
	versionNumber := int64(request.GetIntParam(req.Parameters, "VersionNumber"))
	statementId := request.GetStringParam(req.Parameters, "StatementId")

	if err := s.removeLayerVersionPermissionCore(reqCtx, layerName, versionNumber, statementId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetLayerVersionPolicy returns the resource-based policy for a layer version.
func (s *LambdaService) GetLayerVersionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")
	versionNumber := int64(request.GetIntParam(req.Parameters, "VersionNumber"))

	layerVersion, err := s.getLayerVersionPolicyCore(reqCtx, layerName, versionNumber)
	if err != nil {
		return nil, err
	}

	statements := make([]map[string]interface{}, 0, len(layerVersion.Policies))
	for _, p := range layerVersion.Policies {
		stmt := map[string]interface{}{
			"Sid":       p.Id,
			"Effect":    "Allow",
			"Principal": p.Principal,
			"Action":    p.Action,
			"Resource":  layerVersion.LayerVersionArn,
		}
		statements = append(statements, stmt)
	}

	policyDoc := map[string]interface{}{
		"Version":   "2012-10-17",
		"Id":        "default",
		"Statement": statements,
	}

	policyJSON, err := json.Marshal(policyDoc)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy":     string(policyJSON),
		"RevisionId": layerVersion.RevisionId,
	}, nil
}
