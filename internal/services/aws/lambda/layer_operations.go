// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/store/aws/common"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// PublishLayerVersion publishes a new version of a Lambda layer.
// Creates the layer if it does not exist, and publishes a new version with the provided content.
func (s *LambdaService) PublishLayerVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")
	if layerName == "" {
		return nil, NewInvalidParameter("LayerName", "Layer name is required")
	}

	codeMap := request.GetMapParam(req.Parameters, "Content")
	if codeMap == nil {
		return nil, NewInvalidParameter("Content", "Content is required")
	}

	var layer *lambdastore.Layer
	var err error

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	layers := store.Layers
	layer, err = layers.Get(layerName)
	if err != nil {
		layer = &lambdastore.Layer{
			LayerName:   layerName,
			CreatedDate: time.Now().UTC(),
		}
		layer, err = layers.Create(layer)
		if err != nil {
			if err == lambdastore.ErrResourceConflict {
				layer, err = layers.Get(layerName)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	version := &lambdastore.LayerVersion{
		Description: request.GetStringParam(req.Parameters, "Description"),
		LicenseInfo: request.GetStringParam(req.Parameters, "LicenseInfo"),
	}

	if compats, ok := req.Parameters["CompatibleRuntimes"].([]interface{}); ok {
		for _, c := range compats {
			if str, ok := c.(string); ok {
				version.CompatibleRuntimes = append(version.CompatibleRuntimes, lambdastore.Runtime(str))
			}
		}
	}

	if compats, ok := req.Parameters["CompatibleArchitectures"].([]interface{}); ok {
		for _, c := range compats {
			if str, ok := c.(string); ok {
				version.CompatibleArchitectures = append(version.CompatibleArchitectures, str)
			}
		}
	}

	if zipFileStr, ok := codeMap["ZipFile"].(string); ok && zipFileStr != "" {
		zipFile, err := base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ZipFile encoding: %w", err)
		}
		version.CodeSize = int64(len(zipFile))
		version.CodeSha256 = lambdastore.GenerateCodeHash(zipFile)
	}

	created, err := layers.PublishVersion(layer, version)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Content": map[string]interface{}{
			"Location":   "",
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
	if layerName == "" {
		return nil, NewInvalidParameter("LayerName", "Layer name is required")
	}

	version := request.GetIntParam(req.Parameters, "VersionNumber")
	if version <= 0 {
		return nil, NewInvalidParameter("VersionNumber", "Version number is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Layers.DeleteVersion(layerName, int64(version)); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetLayerVersion retrieves information about a specific version of a Lambda layer.
func (s *LambdaService) GetLayerVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")
	if layerName == "" {
		return nil, NewInvalidParameter("LayerName", "Layer name is required")
	}

	version := request.GetIntParam(req.Parameters, "VersionNumber")
	if version <= 0 {
		return nil, NewInvalidParameter("VersionNumber", "Version number is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	layer, err := store.Layers.Get(layerName)
	if err != nil {
		return nil, NewResourceNotFound("LayerVersion", layerName)
	}
	layerVersion, err := store.Layers.GetVersion(layerName, int64(version))
	if err != nil {
		return nil, NewResourceNotFound("LayerVersion", layerName)
	}

	return map[string]interface{}{
		"Content": map[string]interface{}{
			"Location":   "",
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
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Layers.List(common.ListOptions{MaxItems: maxItems})
	if err != nil {
		return nil, err
	}

	layers := make([]interface{}, 0)
	for _, l := range result.Items {
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

	response := map[string]interface{}{
		"Layers": layers,
	}

	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}

// ListLayerVersions lists all versions of a Lambda layer.
func (s *LambdaService) ListLayerVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	layerName := request.GetStringParam(req.Parameters, "LayerName")
	if layerName == "" {
		return nil, NewInvalidParameter("LayerName", "Layer name is required")
	}

	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Layers.ListVersions(layerName, common.ListOptions{MaxItems: maxItems})
	if err != nil {
		return nil, err
	}

	versions := make([]interface{}, 0)
	for _, v := range result.Items {
		versions = append(versions, map[string]interface{}{
			"LayerVersionArn":         v.LayerVersionArn,
			"Version":                 v.Version,
			"Description":             v.Description,
			"CreatedDate":             v.CreatedDate.Format(timeutils.ISO8601UTCFormat),
			"CompatibleRuntimes":      v.CompatibleRuntimes,
			"LicenseInfo":             v.LicenseInfo,
			"CompatibleArchitectures": v.CompatibleArchitectures,
		})
	}

	response := map[string]interface{}{
		"LayerVersions": versions,
	}

	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}
