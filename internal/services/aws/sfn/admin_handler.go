package sfn

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/states"
	"vorpalstacks/internal/pb/aws/states/statesconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

type AdminHandler struct {
	statesconnect.UnimplementedSFNServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ statesconnect.SFNServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getStoreByRegion(region string) (*sfnstore.StepFunctionStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return sfnstore.NewStepFunctionStore(regionStorage, h.accountId, region), nil
}

func (h *AdminHandler) ListStateMachines(ctx context.Context, req *connect.Request[pb.ListStateMachinesInput]) (*connect.Response[pb.ListStateMachinesOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := store.ListStateMachines(ctx, req.Msg.Maxresults, req.Msg.Nexttoken)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	stateMachines := make([]*pb.StateMachineListItem, len(result.StateMachines))
	for i, sm := range result.StateMachines {
		stateMachines[i] = &pb.StateMachineListItem{
			Statemachinearn: sm.StateMachineArn,
			Name:            sm.Name,
			Creationdate:    sm.CreationDate.Format("2006-01-02T15:04:05Z"),
			Type:            toPbStateMachineType(sm.Type),
		}
	}

	return connect.NewResponse(&pb.ListStateMachinesOutput{
		Statemachines: stateMachines,
		Nexttoken:     result.NextToken,
	}), nil
}

func toPbStateMachineType(typ string) pb.StateMachineType {
	switch typ {
	case "STANDARD":
		return pb.StateMachineType_STATE_MACHINE_TYPE_STANDARD
	case "EXPRESS":
		return pb.StateMachineType_STATE_MACHINE_TYPE_EXPRESS
	default:
		return pb.StateMachineType_STATE_MACHINE_TYPE_STANDARD
	}
}
