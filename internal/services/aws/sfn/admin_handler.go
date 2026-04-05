package sfn

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/sfn"
	"vorpalstacks/internal/pb/aws/sfn/sfnconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// AdminHandler implements the Step Functions gRPC-Web admin console handler. It
// exposes list operations for state machines for the Flutter management UI,
// delegating to the shared StepFunctionService cache.
type AdminHandler struct {
	sfnconnect.UnimplementedSFNServiceHandler
	service *StepFunctionService
}

var _ sfnconnect.SFNServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Step Functions admin console handler backed by
// the given service instance, ensuring the same per-region cached stores are
// used as the HTTP API handlers.
func NewAdminHandler(svc *StepFunctionService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreByRegion(region string) (*sfnstore.StepFunctionStore, error) {
	return h.service.getStoreForRegion(region)
}

// ListStateMachines returns a paginated list of state machines in the
// requested region.
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
