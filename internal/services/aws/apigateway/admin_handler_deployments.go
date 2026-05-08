package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

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
	return connect.NewResponse(toPbDeployment(created)), nil
}

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

func (h *AdminHandler) CreateStage(ctx context.Context, req *connect.Request[pb.CreateStageRequest]) (*connect.Response[pb.Stage], error) {
	if req.Msg.Restapiid == "" || req.Msg.Deploymentid == "" || req.Msg.Stagename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, deployment_id, and stage_name are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	stage := &apigatewaystore.Stage{
		StageName:           req.Msg.Stagename,
		DeploymentId:        req.Msg.Deploymentid,
		Description:         req.Msg.Description,
		CacheClusterEnabled: req.Msg.Cacheclusterenabled,
		Variables:           req.Msg.Variables,
	}

	created, err := stores.restApis.CreateStage(req.Msg.Restapiid, stage)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbStage(created)), nil
}

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
