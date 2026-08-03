package athena

import (
	"context"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/athena"
	athenaconnect "vorpalstacks/internal/pb/aws/athena/athenaconnect"
)

// AdminHandler implements the Athena admin console gRPC-Web handler.
// It delegates to core functions in workgroup_core.go so that validation,
// cascade cleanup, and error handling are shared with the HTTP API path.
type AdminHandler struct {
	athenaconnect.UnimplementedAthenaServiceHandler
	service *AthenaService
}

var _ athenaconnect.AthenaServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(svc *AthenaService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// getStores extracts the region from request headers and returns the full
// athenaStores for that region.
func (h *AdminHandler) getStores(headers http.Header) (*athenaStores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.getStoresForRegion(region)
}

// CreateWorkGroup creates a new Athena work group via the admin console.
func (h *AdminHandler) CreateWorkGroup(ctx context.Context, req *connect.Request[pb.CreateWorkGroupInput]) (*connect.Response[pb.CreateWorkGroupOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	input := WorkGroupCreateInput{
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
	}
	if req.Msg.Configuration != nil {
		protoCfg := req.Msg.Configuration
		cfg := &WorkGroupConfigInput{
			EnforceConfig:           protoCfg.GetEnforceworkgroupconfiguration(),
			PublishMetrics:          protoCfg.GetPublishcloudwatchmetricsenabled(),
			BytesScannedCutoff:      int64(protoCfg.GetBytesscannedcutoffperquery()),
			RequesterPaysEnabled:    protoCfg.GetRequesterpaysenabled(),
			AdditionalConfiguration: protoCfg.GetAdditionalconfiguration(),
			ExecutionRole:           protoCfg.GetExecutionrole(),
		}
		if protoCfg.Resultconfiguration != nil {
			cfg.OutputLocation = protoCfg.Resultconfiguration.Outputlocation
		}
		if protoCfg.Engineversion != nil {
			cfg.EngineVersionSelected = protoCfg.Engineversion.Selectedengineversion
			cfg.EngineVersionEffective = protoCfg.Engineversion.Effectiveengineversion
		}
		input.Config = cfg
	}

	for _, tag := range req.Msg.Tags {
		if input.Tags == nil {
			input.Tags = make(map[string]string)
		}
		input.Tags[tag.Key] = tag.Value
	}

	if err := createWorkGroupCore(stores, input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateWorkGroupOutput{}), nil
}

// DeleteWorkGroup deletes an Athena work group via the admin console.
// Delegates to deleteWorkGroupCore for cascade cleanup of dependent resources.
func (h *AdminHandler) DeleteWorkGroup(ctx context.Context, req *connect.Request[pb.DeleteWorkGroupInput]) (*connect.Response[pb.DeleteWorkGroupOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := deleteWorkGroupCore(stores, req.Msg.Workgroup); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteWorkGroupOutput{}), nil
}

// ListWorkGroups retrieves Athena work groups with pagination support.
func (h *AdminHandler) ListWorkGroups(ctx context.Context, req *connect.Request[pb.ListWorkGroupsInput]) (*connect.Response[pb.ListWorkGroupsOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxResults := clampMaxResults(int(req.Msg.GetMaxresults()), athenaMaxWorkGroupsResults, athenaMaxWorkGroupsResults)

	result, err := listWorkGroupsCore(stores, maxResults, req.Msg.Nexttoken)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var summaries []*pb.WorkGroupSummary
	for _, wg := range result.Items {
		state := pb.WorkGroupState_WORK_GROUP_STATE_DISABLED
		if wg.State == "ENABLED" {
			state = pb.WorkGroupState_WORK_GROUP_STATE_ENABLED
		}
		summary := &pb.WorkGroupSummary{
			Name:  wg.Name,
			State: state,
		}
		if wg.Description != "" {
			summary.Description = wg.Description
		}
		if !wg.CreationTime.IsZero() {
			summary.Creationtime = wg.CreationTime.Format(timeutils.ISO8601UTCFormat)
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListWorkGroupsOutput{
		Workgroups: summaries,
		Nexttoken:  result.NextMarker,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Athena admin console.
func NewConnectHandler(svc *AthenaService) (string, http.Handler) {
	return athenaconnect.NewAthenaServiceHandler(NewAdminHandler(svc))
}
