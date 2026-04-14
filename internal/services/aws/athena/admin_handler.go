package athena

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/athena"
	athenaconnect "vorpalstacks/internal/pb/aws/athena/athenaconnect"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// AdminHandler implements the Athena admin console gRPC-Web handler.
type AdminHandler struct {
	athenaconnect.UnimplementedAthenaServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ athenaconnect.AthenaServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Athena admin handler with the given storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getWorkGroupStore(req *connect.Request[pb.ListWorkGroupsInput]) (*athenastore.WorkGroupStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return athenastore.NewWorkGroupStore(regionStorage, h.accountId, region), nil
}

// ListWorkGroups retrieves all Athena work groups from the regional store.
func (h *AdminHandler) ListWorkGroups(ctx context.Context, req *connect.Request[pb.ListWorkGroupsInput]) (*connect.Response[pb.ListWorkGroupsOutput], error) {
	store, err := h.getWorkGroupStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	workGroups, err := store.ListWorkGroups()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var summaries []*pb.WorkGroupSummary
	for _, wg := range workGroups {
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
		if !wg.CreatedTime.IsZero() {
			summary.Creationtime = wg.CreatedTime.Format("2006-01-02T15:04:05.000Z")
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListWorkGroupsOutput{
		Workgroups: summaries,
	}), nil
}

func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return athenaconnect.NewAthenaServiceHandler(NewAdminHandler(sm, accountID))
}
