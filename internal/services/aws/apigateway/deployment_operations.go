// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateDeployment creates a new deployment for the specified REST API.
func (s *APIGatewayService) CreateDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	deployment := &store.Deployment{
		Description: request.GetStringParam(req.Parameters, "description"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	created, err := stores.restApis.CreateDeployment(apiId, deployment)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName != "" {
		stage := &store.Stage{
			StageName:    stageName,
			DeploymentId: created.Id,
		}

		stageDescription := request.GetStringParam(req.Parameters, "stageDescription")
		if stageDescription != "" {
			stage.Description = stageDescription
		} else {
			stage.Description = "Auto-created stage"
		}

		stage.CacheClusterEnabled = request.GetBoolParam(req.Parameters, "cacheClusterEnabled")
		stage.CacheClusterSize = request.GetStringParam(req.Parameters, "cacheClusterSize")
		stage.TracingEnabled = request.GetBoolParam(req.Parameters, "tracingEnabled")

		if variables, ok := req.Parameters["variables"].(map[string]interface{}); ok {
			stage.Variables = make(map[string]string)
			for k, v := range variables {
				if vs, ok := v.(string); ok {
					stage.Variables[k] = vs
				}
			}
		}

		if canaryMap, ok := req.Parameters["canarySettings"].(map[string]interface{}); ok {
			canary := &store.CanarySettings{DeploymentId: created.Id}
			if v, ok := canaryMap["percentTraffic"].(float64); ok {
				canary.PercentTraffic = v
			}
			if v, ok := canaryMap["useStageCache"].(bool); ok {
				canary.UseStageCache = v
			}
			if overrides, ok := canaryMap["stageVariableOverrides"].(map[string]interface{}); ok {
				canary.StageVariableOverrides = make(map[string]string)
				for k, v := range overrides {
					if vs, ok := v.(string); ok {
						canary.StageVariableOverrides[k] = vs
					}
				}
			}
			stage.CanarySettings = canary
		}

		if _, err := stores.restApis.CreateStage(apiId, stage); err != nil {
			return nil, toApiGatewayError(err)
		}
	}

	return s.toDeploymentResponse(created), nil
}

// GetDeployment retrieves a deployment by its ID for the specified REST API.
func (s *APIGatewayService) GetDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	deploymentId := request.GetStringParam(req.Parameters, "deploymentId")
	if deploymentId == "" {
		deploymentId = getPathParam(req, "deploymentId")
	}
	if deploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	deployment, err := stores.restApis.GetDeployment(apiId, deploymentId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toDeploymentResponse(deployment), nil
}

// DeleteDeployment deletes a deployment from the specified REST API.
func (s *APIGatewayService) DeleteDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	deploymentId := request.GetStringParam(req.Parameters, "deploymentId")
	if deploymentId == "" {
		deploymentId = getPathParam(req, "deploymentId")
	}
	if deploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteDeployment(apiId, deploymentId); err != nil {
		if errors.Is(err, store.ErrDeploymentInUse) {
			return nil, NewConflictException("Deployment is in use by a stage")
		}
		return nil, toApiGatewayError(err)
	}

	return response.EmptyResponse(), nil
}

// GetDeployments lists all deployments for the specified REST API.
func (s *APIGatewayService) GetDeployments(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	limit := request.GetIntParam(req.Parameters, "limit")
	if limit <= 0 {
		limit = 25
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	deployments, err := stores.restApis.ListDeployments(apiId)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(deployments))
	for _, d := range deployments {
		items = append(items, s.toDeploymentResponse(d))
	}

	page, nextPos := paginateItems(items, position, limit)
	result := map[string]interface{}{
		"item": page,
	}
	if nextPos != "" {
		result["position"] = nextPos
	}
	return result, nil
}

func (s *APIGatewayService) toDeploymentResponse(d *store.Deployment) map[string]interface{} {
	response := map[string]interface{}{
		"id":          d.Id,
		"createdDate": timeutils.FormatEpochSeconds(d.CreatedDate),
	}

	if d.Description != "" {
		response["description"] = d.Description
	}

	return response
}

// CreateStage creates a new stage for the specified REST API.
func (s *APIGatewayService) CreateStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}

	stage := &store.Stage{
		StageName:            stageName,
		DeploymentId:         request.GetStringParam(req.Parameters, "deploymentId"),
		Description:          request.GetStringParam(req.Parameters, "description"),
		CacheClusterEnabled:  request.GetBoolParam(req.Parameters, "cacheClusterEnabled"),
		CacheClusterSize:     request.GetStringParam(req.Parameters, "cacheClusterSize"),
		TracingEnabled:       request.GetBoolParam(req.Parameters, "tracingEnabled"),
		DocumentationVersion: request.GetStringParam(req.Parameters, "documentationVersion"),
	}

	if stage.DeploymentId == "" {
		return nil, NewBadRequestException("deploymentId is required")
	}

	if variables, ok := req.Parameters["variables"].(map[string]interface{}); ok {
		stage.Variables = make(map[string]string)
		for k, v := range variables {
			if vs, ok := v.(string); ok {
				stage.Variables[k] = vs
			}
		}
	}

	if tagsRaw := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags")); len(tagsRaw) > 0 {
		stage.Tags = tagutil.MapToTags(tagsRaw)
	}

	if canaryMap, ok := req.Parameters["canarySettings"].(map[string]interface{}); ok {
		canary := &store.CanarySettings{}
		if v, ok := canaryMap["percentTraffic"].(float64); ok {
			canary.PercentTraffic = v
		}
		if v, ok := canaryMap["deploymentId"].(string); ok {
			canary.DeploymentId = v
		}
		if v, ok := canaryMap["useStageCache"].(bool); ok {
			canary.UseStageCache = v
		}
		if overrides, ok := canaryMap["stageVariableOverrides"].(map[string]interface{}); ok {
			canary.StageVariableOverrides = make(map[string]string)
			for k, v := range overrides {
				if vs, ok := v.(string); ok {
					canary.StageVariableOverrides[k] = vs
				}
			}
		}
		stage.CanarySettings = canary
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.restApis.CreateStage(apiId, stage)
	if err != nil {
		return nil, err
	}

	return s.toStageResponse(created), nil
}

// GetStage retrieves a stage by its name for the specified REST API.
func (s *APIGatewayService) GetStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}
	if stageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stage, err := stores.restApis.GetStage(apiId, stageName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toStageResponse(stage), nil
}

// DeleteStage deletes a stage from the specified REST API.
func (s *APIGatewayService) DeleteStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}
	if stageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteStage(apiId, stageName); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// UpdateStage updates an existing stage for the specified REST API.
func (s *APIGatewayService) UpdateStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}
	if stageName == "" {
		return nil, NewBadRequestException("stageName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	stage, err := stores.restApis.GetStage(apiId, stageName)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch {
		case po.Path == "/description":
			stage.Description = po.Value
		case po.Path == "/deploymentId":
			if po.Value != "" {
				if _, err := stores.restApis.GetDeployment(apiId, po.Value); err != nil {
					return nil, NewBadRequestException("deployment not found")
				}
			}
			stage.DeploymentId = po.Value
		case po.Path == "/cacheClusterEnabled":
			stage.CacheClusterEnabled = po.Value == "true"
		case po.Path == "/cacheClusterSize":
			stage.CacheClusterSize = po.Value
		case po.Path == "/tracingEnabled":
			stage.TracingEnabled = po.Value == "true"
		case po.Path == "/clientCertificateId":
			stage.ClientCertificateId = po.Value
		case po.Path == "/documentationVersion":
			stage.DocumentationVersion = po.Value
		case strings.HasPrefix(po.Path, "/variables/"):
			if stage.Variables == nil {
				stage.Variables = make(map[string]string)
			}
			varName := strings.TrimPrefix(po.Path, "/variables/")
			if po.Op == "remove" {
				delete(stage.Variables, varName)
			} else {
				stage.Variables[varName] = po.Value
			}
		case strings.HasPrefix(po.Path, "/accessLogSettings/"):
			if stage.AccessLogSettings == nil {
				stage.AccessLogSettings = &store.AccessLogSettings{}
			}
			switch strings.TrimPrefix(po.Path, "/accessLogSettings/") {
			case "destinationArn":
				stage.AccessLogSettings.DestinationArn = po.Value
			case "format":
				stage.AccessLogSettings.Format = po.Value
			}
		case strings.HasPrefix(po.Path, "/methodSettings/"):
			applyMethodSettingsPatch(stage, po)
		case strings.HasPrefix(po.Path, "/canarySettings/"):
			applyCanarySettingsPatch(stage, po)
		}
	}

	if err := stores.restApis.UpdateStage(apiId, stage); err != nil {
		return nil, err
	}

	return s.toStageResponse(stage), nil
}

// GetStages lists all stages for the specified REST API.
// applyMethodSettingsPatch applies a methodSettings patch operation to the
// stage. The path format is:
//
//	/methodSettings/{resourcePath}/{httpMethod}/{settingName}
//
// where resourcePath uses ~1 for / (e.g. ~1 for root, ~1users for /users).
// The wildcard * is supported for both resourcePath and httpMethod.
func applyMethodSettingsPatch(stage *store.Stage, po PatchOperation) {
	if stage.MethodSettings == nil {
		stage.MethodSettings = make(map[string]*store.MethodSetting)
	}

	rest := strings.TrimPrefix(po.Path, "/methodSettings/")
	parts := strings.SplitN(rest, "/", 3)
	if len(parts) < 3 {
		return
	}

	resourcePath := parts[0]
	httpMethod := parts[1]
	settingName := parts[2]

	key := resourcePath + "/" + httpMethod
	ms, ok := stage.MethodSettings[key]
	if !ok {
		ms = &store.MethodSetting{}
		stage.MethodSettings[key] = ms
	}

	// For remove operations, reset the field to its zero value.
	if po.Op == "remove" {
		switch settingName {
		case "metricsEnabled":
			ms.MetricsEnabled = false
		case "loggingLevel":
			ms.LoggingLevel = ""
		case "dataTraceEnabled":
			ms.DataTraceEnabled = false
		case "throttlingBurstLimit":
			ms.ThrottlingBurstLimit = 0
		case "throttlingRateLimit":
			ms.ThrottlingRateLimit = 0
		case "cachingEnabled":
			ms.CachingEnabled = false
		case "cacheTtlInSeconds":
			ms.CacheTtlInSeconds = 0
		case "cacheDataEncrypted":
			ms.CacheDataEncrypted = false
		case "requireAuthorizationForCacheControl":
			ms.RequireAuthorizationForCacheControl = false
		}
		return
	}

	switch settingName {
	case "metricsEnabled":
		ms.MetricsEnabled = po.Value == "true"
	case "loggingLevel":
		ms.LoggingLevel = po.Value
	case "dataTraceEnabled":
		ms.DataTraceEnabled = po.Value == "true"
	case "throttlingBurstLimit":
		if v, err := strconv.ParseInt(po.Value, 10, 32); err == nil {
			ms.ThrottlingBurstLimit = int32(v)
		}
	case "throttlingRateLimit":
		if v, err := strconv.ParseFloat(po.Value, 64); err == nil {
			ms.ThrottlingRateLimit = v
		}
	case "cachingEnabled":
		ms.CachingEnabled = po.Value == "true"
	case "cacheTtlInSeconds":
		if v, err := strconv.ParseInt(po.Value, 10, 32); err == nil {
			ms.CacheTtlInSeconds = int32(v)
		}
	case "cacheDataEncrypted":
		ms.CacheDataEncrypted = po.Value == "true"
	case "requireAuthorizationForCacheControl":
		ms.RequireAuthorizationForCacheControl = po.Value == "true"
	}
}

// applyCanarySettingsPatch applies a canarySettings patch operation to the
// stage. The path format is:
//
//	/canarySettings/{settingName}
//
// or
//
//	/canarySettings/stageVariableOverrides/{varName}
func applyCanarySettingsPatch(stage *store.Stage, po PatchOperation) {
	if stage.CanarySettings == nil {
		stage.CanarySettings = &store.CanarySettings{}
	}

	rest := strings.TrimPrefix(po.Path, "/canarySettings/")
	parts := strings.SplitN(rest, "/", 2)

	switch parts[0] {
	case "percentTraffic":
		if po.Op == "remove" {
			stage.CanarySettings.PercentTraffic = 0
		} else if v, err := strconv.ParseFloat(po.Value, 64); err == nil {
			stage.CanarySettings.PercentTraffic = v
		}
	case "deploymentId":
		if po.Op == "remove" {
			stage.CanarySettings.DeploymentId = ""
		} else {
			stage.CanarySettings.DeploymentId = po.Value
		}
	case "useStageCache":
		if po.Op == "remove" {
			stage.CanarySettings.UseStageCache = false
		} else {
			stage.CanarySettings.UseStageCache = po.Value == "true"
		}
	case "stageVariableOverrides":
		if len(parts) < 2 {
			return
		}
		if stage.CanarySettings.StageVariableOverrides == nil {
			stage.CanarySettings.StageVariableOverrides = make(map[string]string)
		}
		varName := parts[1]
		if po.Op == "remove" {
			delete(stage.CanarySettings.StageVariableOverrides, varName)
		} else {
			stage.CanarySettings.StageVariableOverrides[varName] = po.Value
		}
	}
}

func (s *APIGatewayService) GetStages(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	limit := request.GetIntParam(req.Parameters, "limit")
	if limit <= 0 {
		limit = 25
	}
	position := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stages, err := stores.restApis.ListStages(apiId)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(stages))
	for _, st := range stages {
		items = append(items, s.toStageResponse(st))
	}

	page, nextPos := paginateItemsWithKey(items, position, limit, "stageName")
	result := map[string]interface{}{
		"item": page,
	}
	if nextPos != "" {
		result["position"] = nextPos
	}
	return result, nil
}

func (s *APIGatewayService) toStageResponse(st *store.Stage) map[string]interface{} {
	response := map[string]interface{}{
		"stageName":           st.StageName,
		"deploymentId":        st.DeploymentId,
		"cacheClusterEnabled": st.CacheClusterEnabled,
		"tracingEnabled":      st.TracingEnabled,
		"createdDate":         timeutils.FormatEpochSeconds(st.CreatedDate),
		"lastUpdatedDate":     timeutils.FormatEpochSeconds(st.LastUpdatedDate),
	}

	if st.Description != "" {
		response["description"] = st.Description
	}
	if st.ClientCertificateId != "" {
		response["clientCertificateId"] = st.ClientCertificateId
	}
	if st.CacheClusterSize != "" {
		response["cacheClusterSize"] = st.CacheClusterSize
	}
	if st.CacheClusterStatus != "" {
		response["cacheClusterStatus"] = st.CacheClusterStatus
	} else {
		if st.CacheClusterEnabled {
			response["cacheClusterStatus"] = "AVAILABLE"
		} else {
			response["cacheClusterStatus"] = "NOT_AVAILABLE"
		}
	}
	if len(st.Variables) > 0 {
		response["variables"] = st.Variables
	}
	if st.DocumentationVersion != "" {
		response["documentationVersion"] = st.DocumentationVersion
	}
	if st.WebAclArn != "" {
		response["webAclArn"] = st.WebAclArn
	}
	if len(st.Tags) > 0 {
		response["tags"] = tagutil.ToMap(st.Tags)
	}

	if st.AccessLogSettings != nil {
		response["accessLogSettings"] = map[string]interface{}{
			"destinationArn": st.AccessLogSettings.DestinationArn,
			"format":         st.AccessLogSettings.Format,
		}
	}

	if st.CanarySettings != nil {
		response["canarySettings"] = map[string]interface{}{
			"percentTraffic":         st.CanarySettings.PercentTraffic,
			"deploymentId":           st.CanarySettings.DeploymentId,
			"stageVariableOverrides": st.CanarySettings.StageVariableOverrides,
			"useStageCache":          st.CanarySettings.UseStageCache,
		}
	}

	if len(st.MethodSettings) > 0 {
		methodSettings := make(map[string]interface{})
		for key, ms := range st.MethodSettings {
			methodSettings[key] = map[string]interface{}{
				"metricsEnabled":                      ms.MetricsEnabled,
				"loggingLevel":                        ms.LoggingLevel,
				"dataTraceEnabled":                    ms.DataTraceEnabled,
				"throttlingBurstLimit":                ms.ThrottlingBurstLimit,
				"throttlingRateLimit":                 ms.ThrottlingRateLimit,
				"cachingEnabled":                      ms.CachingEnabled,
				"cacheTtlInSeconds":                   ms.CacheTtlInSeconds,
				"cacheDataEncrypted":                  ms.CacheDataEncrypted,
				"requireAuthorizationForCacheControl": ms.RequireAuthorizationForCacheControl,
				"unreservedCacheParameters":           ms.UnreservedCacheParameters,
			}
		}
		response["methodSettings"] = methodSettings
	}

	return response
}
