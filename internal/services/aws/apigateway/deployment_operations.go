// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateDeployment creates a new deployment for the specified REST API.
func (s *APIGatewayService) CreateDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	in := &DeploymentInput{
		Description:         request.GetStringParam(req.Parameters, "description"),
		StageName:           request.GetStringParam(req.Parameters, "stageName"),
		StageDescription:    request.GetStringParam(req.Parameters, "stageDescription"),
		CacheClusterSize:    request.GetStringParam(req.Parameters, "cacheClusterSize"),
		CacheClusterEnabled: request.GetBoolParam(req.Parameters, "cacheClusterEnabled"),
		TracingEnabled:      request.GetBoolParam(req.Parameters, "tracingEnabled"),
	}
	if variables, ok := req.Parameters["variables"].(map[string]interface{}); ok {
		in.Variables = make(map[string]string)
		for k, v := range variables {
			if vs, ok := v.(string); ok {
				in.Variables[k] = vs
			}
		}
	}
	if canaryMap, ok := req.Parameters["canarySettings"].(map[string]interface{}); ok {
		in.CanarySettings = parseCanarySettingsInput(canaryMap)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createDeploymentCore(stores, apiId, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toDeploymentResponse(created), nil
}

// GetDeployment retrieves a deployment by its ID for the specified REST API.
func (s *APIGatewayService) GetDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	deploymentId := request.GetStringParam(req.Parameters, "deploymentId")
	if deploymentId == "" {
		deploymentId = getPathParam(req, "deploymentId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	deployment, err := s.getDeploymentCore(stores, apiId, deploymentId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toDeploymentResponse(deployment), nil
}

// DeleteDeployment deletes a deployment from the specified REST API.
func (s *APIGatewayService) DeleteDeployment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	deploymentId := request.GetStringParam(req.Parameters, "deploymentId")
	if deploymentId == "" {
		deploymentId = getPathParam(req, "deploymentId")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteDeploymentCore(stores, apiId, deploymentId); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// GetDeployments lists all deployments for the specified REST API.
func (s *APIGatewayService) GetDeployments(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	deployments, err := s.listDeploymentsCore(stores, apiId)
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

func (s *APIGatewayService) toDeploymentResponse(d *apigateway.Deployment) map[string]interface{} {
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
	in := &StageInput{
		StageName:            request.GetStringParam(req.Parameters, "stageName"),
		DeploymentId:         request.GetStringParam(req.Parameters, "deploymentId"),
		Description:          request.GetStringParam(req.Parameters, "description"),
		CacheClusterEnabled:  request.GetBoolParam(req.Parameters, "cacheClusterEnabled"),
		CacheClusterSize:     request.GetStringParam(req.Parameters, "cacheClusterSize"),
		DocumentationVersion: request.GetStringParam(req.Parameters, "documentationVersion"),
		TracingEnabled:       request.GetBoolParam(req.Parameters, "tracingEnabled"),
	}
	if variables, ok := req.Parameters["variables"].(map[string]interface{}); ok {
		in.Variables = make(map[string]string)
		for k, v := range variables {
			if vs, ok := v.(string); ok {
				in.Variables[k] = vs
			}
		}
	}
	if tagsRaw := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags")); len(tagsRaw) > 0 {
		in.Tags = tagsRaw
	}
	if canaryMap, ok := req.Parameters["canarySettings"].(map[string]interface{}); ok {
		in.CanarySettings = parseCanarySettingsInput(canaryMap)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createStageCore(stores, apiId, in)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toStageResponse(created), nil
}

// GetStage retrieves a stage by its name for the specified REST API.
func (s *APIGatewayService) GetStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stage, err := s.getStageCore(stores, apiId, stageName)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toStageResponse(stage), nil
}

// DeleteStage deletes a stage from the specified REST API.
func (s *APIGatewayService) DeleteStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteStageCore(stores, apiId, stageName); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// UpdateStage updates an existing stage for the specified REST API.
func (s *APIGatewayService) UpdateStage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName == "" {
		stageName = getPathParam(req, "stageName")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	stage, err := s.updateStageCore(stores, apiId, stageName, ops)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toStageResponse(stage), nil
}

// GetStages lists all stages for the specified REST API.
func (s *APIGatewayService) GetStages(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	stages, err := s.listStagesCore(stores, apiId)
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

// parseCanarySettingsInput extracts the canary settings from a raw
// map[string]interface{} (HTTP query param format) into the
// transport-agnostic input struct.
func parseCanarySettingsInput(m map[string]interface{}) *CanarySettingsInput {
	in := &CanarySettingsInput{}
	if v, ok := m["percentTraffic"].(float64); ok {
		in.PercentTraffic = v
	}
	if v, ok := m["deploymentId"].(string); ok {
		in.DeploymentId = v
	}
	if v, ok := m["useStageCache"].(bool); ok {
		in.UseStageCache = v
	}
	if overrides, ok := m["stageVariableOverrides"].(map[string]interface{}); ok {
		in.StageVariableOverrides = make(map[string]string)
		for k, v := range overrides {
			if vs, ok := v.(string); ok {
				in.StageVariableOverrides[k] = vs
			}
		}
	}
	return in
}

func (s *APIGatewayService) toStageResponse(st *apigateway.Stage) map[string]interface{} {
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
