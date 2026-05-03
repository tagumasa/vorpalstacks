package sfn

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/sfn"
	"vorpalstacks/internal/pb/aws/sfn/sfnconnect"
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

// CreateStateMachine creates a new Step Functions state machine via the admin
// console gRPC-Web interface.
func (h *AdminHandler) CreateStateMachine(ctx context.Context, req *connect.Request[pb.CreateStateMachineInput]) (*connect.Response[pb.CreateStateMachineOutput], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}
	if req.Msg.Definition == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("definition is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	smType := "STANDARD"
	if req.Msg.Type == pb.StateMachineType_STATE_MACHINE_TYPE_EXPRESS {
		smType = "EXPRESS"
	}

	tags := make(map[string]string)
	for _, t := range req.Msg.Tags {
		tags[t.Key] = t.Value
	}

	now := time.Now().UTC()
	sm := &sfnstore.StateMachine{
		Name:        req.Msg.Name,
		Definition:  req.Msg.Definition,
		RoleArn:     req.Msg.Rolearn,
		Type:        smType,
		CreationDate: now,
		UpdateDate:  now,
		Status:      "ACTIVE",
		Tags:        tags,
	}

	if err := store.CreateStateMachine(ctx, sm); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateStateMachineOutput{
		Statemachinearn: sm.StateMachineArn,
		Creationdate:    sm.CreationDate.Format("2006-01-02T15:04:05Z"),
	}), nil
}

// DeleteStateMachine deletes a Step Functions state machine via the admin
// console gRPC-Web interface.
func (h *AdminHandler) DeleteStateMachine(ctx context.Context, req *connect.Request[pb.DeleteStateMachineInput]) (*connect.Response[pb.DeleteStateMachineOutput], error) {
	if req.Msg.Statemachinearn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("stateMachineArn is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := store.DeleteStateMachine(ctx, req.Msg.Statemachinearn); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteStateMachineOutput{}), nil
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

// NewConnectHandler creates a gRPC-Web connect handler for the Sfn admin console.
func NewConnectHandler(svc *StepFunctionService) (string, http.Handler) {
	return sfnconnect.NewSFNServiceHandler(NewAdminHandler(svc))
}
