// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"errors"
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
	created, err := stores.restApis.CreateDeployment(apiId, deployment)
	if err != nil {
		return nil, err
	}

	stageName := request.GetStringParam(req.Parameters, "stageName")
	if stageName != "" {
		stage := &store.Stage{
			StageName:    stageName,
			DeploymentId: created.Id,
			Description:  "Auto-created stage",
		}
		if _, err := stores.restApis.CreateStage(apiId, stage); err != nil {
			return nil, err
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
		return nil, ErrNotFoundException
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
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// GetDeployments lists all deployments for the specified REST API.
func (s *APIGatewayService) GetDeployments(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

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

	return map[string]interface{}{
		"item": items,
	}, nil
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
		StageName:           stageName,
		DeploymentId:        request.GetStringParam(req.Parameters, "deploymentId"),
		Description:         request.GetStringParam(req.Parameters, "description"),
		CacheClusterEnabled: request.GetBoolParam(req.Parameters, "cacheClusterEnabled"),
		CacheClusterSize:    request.GetStringParam(req.Parameters, "cacheClusterSize"),
		TracingEnabled:      request.GetBoolParam(req.Parameters, "tracingEnabled"),
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
	stage, err := stores.restApis.GetStage(apiId, stageName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	patchOps, ok := req.Parameters["patchOperations"].([]interface{})
	if ok {
		for _, op := range patchOps {
			if opMap, ok := op.(map[string]interface{}); ok {
				path := ""
				value := ""
				if p, ok := opMap["path"].(string); ok {
					path = p
				}
				if v, ok := opMap["value"].(string); ok {
					value = v
				}

				switch path {
				case "/description":
					stage.Description = value
				case "/deploymentId":
					if value != "" {
						_, err := stores.restApis.GetDeployment(apiId, value)
						if err != nil {
							return nil, NewBadRequestException("deployment not found")
						}
					}
					stage.DeploymentId = value
				case "/cacheClusterEnabled":
					stage.CacheClusterEnabled = value == "true"
				case "/cacheClusterSize":
					stage.CacheClusterSize = value
				case "/tracingEnabled":
					stage.TracingEnabled = value == "true"
				}
				if strings.HasPrefix(path, "/variables/") {
					if stage.Variables == nil {
						stage.Variables = make(map[string]string)
					}
					varName := strings.TrimPrefix(path, "/variables/")
					if opMap["op"] == "remove" {
						delete(stage.Variables, varName)
					} else {
						stage.Variables[varName] = value
					}
				}
				if strings.HasPrefix(path, "/accessLogSettings/") {
					if stage.AccessLogSettings == nil {
						stage.AccessLogSettings = &store.AccessLogSettings{}
					}
					setting := strings.TrimPrefix(path, "/accessLogSettings/")
					switch setting {
					case "destinationArn":
						stage.AccessLogSettings.DestinationArn = value
					case "format":
						stage.AccessLogSettings.Format = value
					}
				}
			}
		}
	}

	if err := stores.restApis.UpdateStage(apiId, stage); err != nil {
		return nil, err
	}

	return s.toStageResponse(stage), nil
}

// GetStages lists all stages for the specified REST API.
func (s *APIGatewayService) GetStages(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

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

	return map[string]interface{}{
		"item": items,
	}, nil
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
	if st.CacheClusterSize != "" {
		response["cacheClusterSize"] = st.CacheClusterSize
	}
	if st.CacheClusterStatus != "" {
		response["cacheClusterStatus"] = st.CacheClusterStatus
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
