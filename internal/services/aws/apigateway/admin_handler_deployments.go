package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	tagutil "vorpalstacks/internal/common/tags"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// CreateDeployment creates a new deployment for a REST API.
func (h *AdminHandler) CreateDeployment(ctx context.Context, req *connect.Request[pb.CreateDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	deployment := &apigatewaystore.Deployment{
		Description: req.Msg.Description,
	}

	created, err := stores.restApis.CreateDeployment(req.Msg.Restapiid, deployment)
	if err != nil {
		return nil, storeErr(err)
	}

	if req.Msg.Stagename != "" {
		stage := &apigatewaystore.Stage{
			StageName:           req.Msg.Stagename,
			DeploymentId:        created.Id,
			CacheClusterEnabled: req.Msg.GetCacheclusterenabled(),
			CacheClusterSize:    cacheClusterSizeFromPb(req.Msg.Cacheclustersize),
			TracingEnabled:      req.Msg.GetTracingenabled(),
			Variables:           req.Msg.Variables,
		}
		if req.Msg.Stagedescription != "" {
			stage.Description = req.Msg.Stagedescription
		} else {
			stage.Description = "Auto-created stage"
		}
		if req.Msg.Canarysettings != nil {
			stage.CanarySettings = &apigatewaystore.CanarySettings{
				PercentTraffic:         req.Msg.Canarysettings.Percenttraffic,
				DeploymentId:           created.Id,
				StageVariableOverrides: req.Msg.Canarysettings.Stagevariableoverrides,
				UseStageCache:          req.Msg.Canarysettings.GetUsestagecache(),
			}
		}
		if _, err := stores.restApis.CreateStage(req.Msg.Restapiid, stage); err != nil {
			return nil, storeErr(err)
		}
	}

	return connect.NewResponse(toPbDeployment(created)), nil
}

// GetDeployments returns all deployments for a REST API.
func (h *AdminHandler) GetDeployments(ctx context.Context, req *connect.Request[pb.GetDeploymentsRequest]) (*connect.Response[pb.Deployments], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	deployments, err := stores.restApis.ListDeployments(req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.Deployment, 0, len(deployments))
	for _, d := range deployments {
		items = append(items, toPbDeployment(d))
	}
	return connect.NewResponse(&pb.Deployments{Items: items}), nil
}

// GetDeployment returns a single deployment by ID.
func (h *AdminHandler) GetDeployment(ctx context.Context, req *connect.Request[pb.GetDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	if req.Msg.Restapiid == "" || req.Msg.Deploymentid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and deployment_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	d, err := stores.restApis.GetDeployment(req.Msg.Restapiid, req.Msg.Deploymentid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbDeployment(d)), nil
}

// DeleteDeployment removes a deployment from a REST API.
func (h *AdminHandler) DeleteDeployment(ctx context.Context, req *connect.Request[pb.DeleteDeploymentRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" || req.Msg.Deploymentid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and deployment_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.DeleteDeployment(req.Msg.Restapiid, req.Msg.Deploymentid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateStage creates a new stage for a REST API deployment.
func (h *AdminHandler) CreateStage(ctx context.Context, req *connect.Request[pb.CreateStageRequest]) (*connect.Response[pb.Stage], error) {
	if req.Msg.Restapiid == "" || req.Msg.Deploymentid == "" || req.Msg.Stagename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, deployment_id, and stage_name are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	stage := &apigatewaystore.Stage{
		StageName:            req.Msg.Stagename,
		DeploymentId:         req.Msg.Deploymentid,
		Description:          req.Msg.Description,
		CacheClusterEnabled:  req.Msg.GetCacheclusterenabled(),
		CacheClusterSize:     cacheClusterSizeFromPb(req.Msg.Cacheclustersize),
		DocumentationVersion: req.Msg.Documentationversion,
		TracingEnabled:       req.Msg.GetTracingenabled(),
		Variables:            req.Msg.Variables,
	}
	if len(req.Msg.Tags) > 0 {
		stage.Tags = tagutil.MapToTags(req.Msg.Tags)
	}
	if req.Msg.Canarysettings != nil {
		stage.CanarySettings = &apigatewaystore.CanarySettings{
			PercentTraffic:         req.Msg.Canarysettings.Percenttraffic,
			StageVariableOverrides: req.Msg.Canarysettings.Stagevariableoverrides,
			UseStageCache:          req.Msg.Canarysettings.GetUsestagecache(),
		}
	}

	created, err := stores.restApis.CreateStage(req.Msg.Restapiid, stage)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbStage(created)), nil
}

// GetStages returns all stages for a REST API.
func (h *AdminHandler) GetStages(ctx context.Context, req *connect.Request[pb.GetStagesRequest]) (*connect.Response[pb.Stages], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	stages, err := stores.restApis.ListStages(req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.Stage, 0, len(stages))
	for _, s := range stages {
		items = append(items, toPbStage(s))
	}
	return connect.NewResponse(&pb.Stages{Item: items}), nil
}

// GetStage returns a single stage by name.
func (h *AdminHandler) GetStage(ctx context.Context, req *connect.Request[pb.GetStageRequest]) (*connect.Response[pb.Stage], error) {
	if req.Msg.Restapiid == "" || req.Msg.Stagename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and stage_name are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	s, err := stores.restApis.GetStage(req.Msg.Restapiid, req.Msg.Stagename)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbStage(s)), nil
}

// DeleteStage removes a stage from a REST API.
func (h *AdminHandler) DeleteStage(ctx context.Context, req *connect.Request[pb.DeleteStageRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" || req.Msg.Stagename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and stage_name are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.DeleteStage(req.Msg.Restapiid, req.Msg.Stagename); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
