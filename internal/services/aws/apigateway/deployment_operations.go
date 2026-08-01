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

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	// Pre-validate stage parameters before persisting the deployment so
	// that validation failures do not leave an orphaned deployment.
	stageName := request.GetStringParam(req.Parameters, "stageName")
	cacheClusterSize := request.GetStringParam(req.Parameters, "cacheClusterSize")
	if stageName != "" {
		if !validateStageName(stageName) {
			return nil, NewBadRequestException("Invalid stage name: must be alphanumeric, underscore, or hyphen, max 128 characters")
		}
		if !validateCacheClusterSize(cacheClusterSize) {
			return nil, NewBadRequestException("Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237")
		}
	}

	deployment := &store.Deployment{
		Description: request.GetStringParam(req.Parameters, "description"),
	}

	created, err := stores.restApis.CreateDeployment(apiId, deployment)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

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
		stage.CacheClusterSize = cacheClusterSize
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
			// Compensating delete: roll back the deployment so that a
			// failed stage creation does not leave an orphaned resource.
			_ = stores.restApis.DeleteDeployment(apiId, created.Id)
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

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
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

	page, nextPos, found := paginateItems(items, position, limit)
	if !found {
		return nil, NewBadRequestException("Invalid position: " + position)
	}
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
	if !validateStageName(stageName) {
		return nil, NewBadRequestException("Invalid stage name: must be alphanumeric, underscore, or hyphen, max 128 characters")
	}

	cacheClusterSize := request.GetStringParam(req.Parameters, "cacheClusterSize")
	if !validateCacheClusterSize(cacheClusterSize) {
		return nil, NewBadRequestException("Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237")
	}

	stage := &store.Stage{
		StageName:            stageName,
		DeploymentId:         request.GetStringParam(req.Parameters, "deploymentId"),
		Description:          request.GetStringParam(req.Parameters, "description"),
		CacheClusterEnabled:  request.GetBoolParam(req.Parameters, "cacheClusterEnabled"),
		CacheClusterSize:     cacheClusterSize,
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

	if s.runtimeServer != nil {
		s.runtimeServer.CleanupStageThrottlers(stageName)
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
			if !validateCacheClusterSize(po.Value) {
				return nil, NewBadRequestException("Invalid cacheClusterSize: must be one of 0.5, 1.6, 6.1, 13.5, 28.4, 58.2, 118, 237")
			}
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
			if err := applyMethodSettingsPatch(stage, po); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/canarySettings/"):
			if err := applyCanarySettingsPatch(stage, po); err != nil {
				return nil, err
			}
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
func applyMethodSettingsPatch(stage *store.Stage, po PatchOperation) error {
	if stage.MethodSettings == nil {
		stage.MethodSettings = make(map[string]*store.MethodSetting)
	}

	rest := strings.TrimPrefix(po.Path, "/methodSettings/")
	parts := strings.SplitN(rest, "/", 3)

	resourcePath := parts[0]
	httpMethod := ""
	settingName := ""
	if len(parts) >= 2 {
		httpMethod = parts[1]
	}
	if len(parts) >= 3 {
		settingName = parts[2]
	}

	key := resourcePath + "/" + httpMethod

	// "remove" on the entire entry (no settingName) deletes the entry.
	if po.Op == "remove" && settingName == "" {
		delete(stage.MethodSettings, key)
		return nil
	}

	ms, ok := stage.MethodSettings[key]
	if !ok {
		ms = &store.MethodSetting{}
		stage.MethodSettings[key] = ms
	}

	// For remove operations on individual settings, reset the field to its
	// zero value. If all fields end up at zero, remove the entry entirely
	// to avoid cluttering the settings map with empty entries.
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
		if isMethodSettingEmpty(ms) {
			delete(stage.MethodSettings, key)
		}
		return nil
	}

	switch settingName {
	case "metricsEnabled":
		ms.MetricsEnabled = po.Value == "true"
	case "loggingLevel":
		if !validateLoggingLevel(po.Value) {
			return NewBadRequestException("Invalid loggingLevel: must be OFF, ERROR, or INFO")
		}
		ms.LoggingLevel = po.Value
	case "dataTraceEnabled":
		ms.DataTraceEnabled = po.Value == "true"
	case "throttlingBurstLimit":
		v, err := strconv.ParseInt(po.Value, 10, 32)
		if err != nil {
			return NewBadRequestException("Invalid throttlingBurstLimit: not a number")
		}
		if v < 0 || v > 100000 {
			return NewBadRequestException("throttlingBurstLimit must be between 0 and 100000")
		}
		ms.ThrottlingBurstLimit = int32(v)
	case "throttlingRateLimit":
		v, err := strconv.ParseFloat(po.Value, 64)
		if err != nil {
			return NewBadRequestException("Invalid throttlingRateLimit: not a number")
		}
		if v < 0 || v > 100000 {
			return NewBadRequestException("throttlingRateLimit must be between 0 and 100000")
		}
		ms.ThrottlingRateLimit = v
	case "cachingEnabled":
		ms.CachingEnabled = po.Value == "true"
	case "cacheTtlInSeconds":
		v, err := strconv.ParseInt(po.Value, 10, 32)
		if err != nil {
			return NewBadRequestException("Invalid cacheTtlInSeconds: not a number")
		}
		if v < 0 || v > 86400 {
			return NewBadRequestException("cacheTtlInSeconds must be between 0 and 86400")
		}
		ms.CacheTtlInSeconds = int32(v)
	case "cacheDataEncrypted":
		ms.CacheDataEncrypted = po.Value == "true"
	case "requireAuthorizationForCacheControl":
		ms.RequireAuthorizationForCacheControl = po.Value == "true"
	}
	return nil
}

// isMethodSettingEmpty returns true if all fields of the MethodSetting are
// at their zero values, indicating the entry carries no meaningful
// configuration and can be removed from the map.
func isMethodSettingEmpty(ms *store.MethodSetting) bool {
	return !ms.MetricsEnabled &&
		ms.LoggingLevel == "" &&
		!ms.DataTraceEnabled &&
		ms.ThrottlingBurstLimit == 0 &&
		ms.ThrottlingRateLimit == 0 &&
		!ms.CachingEnabled &&
		ms.CacheTtlInSeconds == 0 &&
		!ms.CacheDataEncrypted &&
		!ms.RequireAuthorizationForCacheControl &&
		len(ms.UnreservedCacheParameters) == 0
}

// applyCanarySettingsPatch applies a canarySettings patch operation to the
// stage. The path format is:
//
//	/canarySettings/{settingName}
//
// or
//
//	/canarySettings/stageVariableOverrides/{varName}
func applyCanarySettingsPatch(stage *store.Stage, po PatchOperation) error {
	if stage.CanarySettings == nil {
		stage.CanarySettings = &store.CanarySettings{}
	}

	rest := strings.TrimPrefix(po.Path, "/canarySettings/")
	parts := strings.SplitN(rest, "/", 2)

	switch parts[0] {
	case "percentTraffic":
		if po.Op == "remove" {
			stage.CanarySettings.PercentTraffic = 0
		} else {
			v, err := strconv.ParseFloat(po.Value, 64)
			if err != nil {
				return NewBadRequestException("Invalid percentTraffic: not a number")
			}
			if v < 0 || v > 100 {
				return NewBadRequestException("percentTraffic must be between 0 and 100")
			}
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
			return nil
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
	return nil
}

func (s *APIGatewayService) GetStages(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
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

	page, nextPos, found := paginateItemsWithKey(items, position, limit, "stageName")
	if !found {
		return nil, NewBadRequestException("Invalid position: " + position)
	}
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
