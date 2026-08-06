package sfn

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/sfn"
	"vorpalstacks/internal/pb/aws/sfn/sfnconnect"
)

// AdminHandler implements the Step Functions gRPC-Web admin console handler.
// It delegates all operations to the shared service-layer Core functions,
// ensuring that validation and persistence follow the same single code path
// as the HTTP API (AGENTS.md #29: zero store imports).
type AdminHandler struct {
	sfnconnect.UnimplementedSFNServiceHandler
	service *StepFunctionService
}

var _ sfnconnect.SFNServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Step Functions admin console handler backed by
// the given service instance.
func NewAdminHandler(svc *StepFunctionService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListStateMachines returns a paginated list of state machines in the
// requested region.
func (h *AdminHandler) ListStateMachines(ctx context.Context, req *connect.Request[pb.ListStateMachinesInput]) (*connect.Response[pb.ListStateMachinesOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.listStateMachinesCore(ctx, store, ListStateMachinesInput{
		MaxResults: req.Msg.GetMaxresults(),
		NextToken:  req.Msg.Nexttoken,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	stateMachines := make([]*pb.StateMachineListItem, len(result.StateMachines))
	for i, sm := range result.StateMachines {
		item := stateMachineListItemToProto(sm)
		item.Creationdate = sm.CreationDate.Format(timeutils.ISO8601UTCFormat)
		stateMachines[i] = item
	}

	return connect.NewResponse(&pb.ListStateMachinesOutput{
		Statemachines: stateMachines,
		Nexttoken:     result.NextToken,
	}), nil
}

// CreateStateMachine creates a new Step Functions state machine via the admin
// console gRPC-Web interface.
func (h *AdminHandler) CreateStateMachine(ctx context.Context, req *connect.Request[pb.CreateStateMachineInput]) (*connect.Response[pb.CreateStateMachineOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tags := make(map[string]string)
	for _, t := range req.Msg.Tags {
		tags[t.Key] = t.Value
	}

	result, err := h.service.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:       req.Msg.Name,
		Definition: req.Msg.Definition,
		RoleArn:    req.Msg.Rolearn,
		Type:       smTypeFromProto(req.Msg.Type),
		Tags:       tags,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateStateMachineOutput{
		Statemachinearn: result.StateMachineArn,
		Creationdate:    result.CreationDate.Format(timeutils.ISO8601UTCFormat),
	}), nil
}

// DeleteStateMachine deletes a Step Functions state machine via the admin
// console gRPC-Web interface.
func (h *AdminHandler) DeleteStateMachine(ctx context.Context, req *connect.Request[pb.DeleteStateMachineInput]) (*connect.Response[pb.DeleteStateMachineOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteStateMachineCore(ctx, store, DeleteStateMachineInput{
		StateMachineArn: req.Msg.Statemachinearn,
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteStateMachineOutput{}), nil
}

// smTypeFromProto converts the proto enum to the Smithy string type.
func smTypeFromProto(t pb.StateMachineType) string {
	if t == pb.StateMachineType_STATE_MACHINE_TYPE_EXPRESS {
		return "EXPRESS"
	}
	return "STANDARD"
}

// NewConnectHandler creates a gRPC-Web connect handler for the SFN admin console.
func NewConnectHandler(svc *StepFunctionService) (string, http.Handler) {
	return sfnconnect.NewSFNServiceHandler(NewAdminHandler(svc))
}
