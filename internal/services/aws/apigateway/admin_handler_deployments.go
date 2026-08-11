package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// CreateDeployment creates a new deployment for a REST API.
func (h *AdminHandler) CreateDeployment(ctx context.Context, req *connect.Request[pb.CreateDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := &DeploymentInput{
		Description:         req.Msg.Description,
		StageName:           req.Msg.Stagename,
		StageDescription:    req.Msg.Stagedescription,
		CacheClusterSize:    cacheClusterSizeFromPb(req.Msg.Cacheclustersize),
		CacheClusterEnabled: req.Msg.GetCacheclusterenabled(),
		TracingEnabled:      req.Msg.GetTracingenabled(),
		Variables:           req.Msg.Variables,
	}
	if req.Msg.Canarysettings != nil {
		in.CanarySettings = &CanarySettingsInput{
			PercentTraffic:         req.Msg.Canarysettings.Percenttraffic,
			StageVariableOverrides: req.Msg.Canarysettings.Stagevariableoverrides,
			UseStageCache:          req.Msg.Canarysettings.GetUsestagecache(),
		}
	}

	created, err := h.service.createDeploymentCore(stores, req.Msg.Restapiid, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbDeployment(created)), nil
}

// GetDeployments returns all deployments for a REST API.
func (h *AdminHandler) GetDeployments(ctx context.Context, req *connect.Request[pb.GetDeploymentsRequest]) (*connect.Response[pb.Deployments], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	deployments, err := h.service.listDeploymentsCore(stores, req.Msg.Restapiid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	limit := int(req.Msg.GetLimit())
	start, end, nextPos, ok := paginateAdminList(len(deployments), req.Msg.Position, limit)
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid position: %s", req.Msg.Position))
	}

	items := make([]*pb.Deployment, 0, end-start)
	for _, d := range deployments[start:end] {
		items = append(items, toPbDeployment(d))
	}
	resp := &pb.Deployments{Items: items}
	if nextPos != "" {
		resp.Position = nextPos
	}
	return connect.NewResponse(resp), nil
}

// GetDeployment returns a single deployment by ID.
func (h *AdminHandler) GetDeployment(ctx context.Context, req *connect.Request[pb.GetDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	d, err := h.service.getDeploymentCore(stores, req.Msg.Restapiid, req.Msg.Deploymentid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbDeployment(d)), nil
}

// DeleteDeployment removes a deployment from a REST API.
func (h *AdminHandler) DeleteDeployment(ctx context.Context, req *connect.Request[pb.DeleteDeploymentRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteDeploymentCore(stores, req.Msg.Restapiid, req.Msg.Deploymentid); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// CreateStage creates a new stage for a REST API deployment.
func (h *AdminHandler) CreateStage(ctx context.Context, req *connect.Request[pb.CreateStageRequest]) (*connect.Response[pb.Stage], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := &StageInput{
		StageName:            req.Msg.Stagename,
		DeploymentId:         req.Msg.Deploymentid,
		Description:          req.Msg.Description,
		CacheClusterSize:     cacheClusterSizeFromPb(req.Msg.Cacheclustersize),
		CacheClusterEnabled:  req.Msg.GetCacheclusterenabled(),
		DocumentationVersion: req.Msg.Documentationversion,
		TracingEnabled:       req.Msg.GetTracingenabled(),
		Variables:            req.Msg.Variables,
		Tags:                 req.Msg.Tags,
	}
	if req.Msg.Canarysettings != nil {
		in.CanarySettings = &CanarySettingsInput{
			PercentTraffic:         req.Msg.Canarysettings.Percenttraffic,
			StageVariableOverrides: req.Msg.Canarysettings.Stagevariableoverrides,
			UseStageCache:          req.Msg.Canarysettings.GetUsestagecache(),
		}
	}

	created, err := h.service.createStageCore(stores, req.Msg.Restapiid, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbStage(created)), nil
}

// GetStages returns all stages for a REST API.
func (h *AdminHandler) GetStages(ctx context.Context, req *connect.Request[pb.GetStagesRequest]) (*connect.Response[pb.Stages], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	stages, err := h.service.listStagesCore(stores, req.Msg.Restapiid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	items := make([]*pb.Stage, 0, len(stages))
	for _, s := range stages {
		items = append(items, toPbStage(s))
	}
	return connect.NewResponse(&pb.Stages{Item: items}), nil
}

// GetStage returns a single stage by name.
func (h *AdminHandler) GetStage(ctx context.Context, req *connect.Request[pb.GetStageRequest]) (*connect.Response[pb.Stage], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	s, err := h.service.getStageCore(stores, req.Msg.Restapiid, req.Msg.Stagename)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbStage(s)), nil
}

// DeleteStage removes a stage from a REST API.
func (h *AdminHandler) DeleteStage(ctx context.Context, req *connect.Request[pb.DeleteStageRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteStageCore(stores, req.Msg.Restapiid, req.Msg.Stagename); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
