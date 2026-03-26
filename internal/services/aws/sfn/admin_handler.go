package sfn

import (
	"connectrpc.com/connect"
	"context"
	"fmt"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/states"
	"vorpalstacks/internal/pb/aws/states/statesconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// AdminHandler provides Step Functions (SFN) service administration functionality.
// It implements the SFNServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	statesconnect.UnimplementedSFNServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ statesconnect.SFNServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Step Functions AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getStoreByRegion retrieves the Step Functions store for the specified region.
// It fetches the region-specific storage and creates a new StepFunctionStore instance.
func (h *AdminHandler) getStoreByRegion(region string) (*sfnstore.StepFunctionStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return sfnstore.NewStepFunctionStore(regionStorage, h.accountId, region), nil
}

// ListStateMachines lists state machines in Step Functions.
// It supports pagination via the NextToken parameter.
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

// CreateActivity creates a new activity in Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) CreateActivity(ctx context.Context, req *connect.Request[pb.CreateActivityInput]) (*connect.Response[pb.CreateActivityOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// CreateStateMachine creates a new state machine in Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) CreateStateMachine(ctx context.Context, req *connect.Request[pb.CreateStateMachineInput]) (*connect.Response[pb.CreateStateMachineOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// CreateStateMachineAlias creates a new alias for a state machine.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) CreateStateMachineAlias(ctx context.Context, req *connect.Request[pb.CreateStateMachineAliasInput]) (*connect.Response[pb.CreateStateMachineAliasOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteActivity deletes an activity from Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteActivity(ctx context.Context, req *connect.Request[pb.DeleteActivityInput]) (*connect.Response[pb.DeleteActivityOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteStateMachine deletes a state machine from Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteStateMachine(ctx context.Context, req *connect.Request[pb.DeleteStateMachineInput]) (*connect.Response[pb.DeleteStateMachineOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteStateMachineAlias deletes a state machine alias.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteStateMachineAlias(ctx context.Context, req *connect.Request[pb.DeleteStateMachineAliasInput]) (*connect.Response[pb.DeleteStateMachineAliasOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteStateMachineVersion deletes a state machine version.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) DeleteStateMachineVersion(ctx context.Context, req *connect.Request[pb.DeleteStateMachineVersionInput]) (*connect.Response[pb.DeleteStateMachineVersionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DescribeActivity returns details about an activity.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeActivity(ctx context.Context, req *connect.Request[pb.DescribeActivityInput]) (*connect.Response[pb.DescribeActivityOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeExecution returns details about an execution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeExecution(ctx context.Context, req *connect.Request[pb.DescribeExecutionInput]) (*connect.Response[pb.DescribeExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeMapRun returns details about a Map Run.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeMapRun(ctx context.Context, req *connect.Request[pb.DescribeMapRunInput]) (*connect.Response[pb.DescribeMapRunOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeStateMachine returns details about a state machine.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeStateMachine(ctx context.Context, req *connect.Request[pb.DescribeStateMachineInput]) (*connect.Response[pb.DescribeStateMachineOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeStateMachineAlias returns details about a state machine alias.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeStateMachineAlias(ctx context.Context, req *connect.Request[pb.DescribeStateMachineAliasInput]) (*connect.Response[pb.DescribeStateMachineAliasOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeStateMachineForExecution returns details about a state machine for an execution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeStateMachineForExecution(ctx context.Context, req *connect.Request[pb.DescribeStateMachineForExecutionInput]) (*connect.Response[pb.DescribeStateMachineForExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetActivityTask retrieves a task for an activity worker.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetActivityTask(ctx context.Context, req *connect.Request[pb.GetActivityTaskInput]) (*connect.Response[pb.GetActivityTaskOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetExecutionHistory returns the history of an execution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetExecutionHistory(ctx context.Context, req *connect.Request[pb.GetExecutionHistoryInput]) (*connect.Response[pb.GetExecutionHistoryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListActivities lists activities in Step Functions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListActivities(ctx context.Context, req *connect.Request[pb.ListActivitiesInput]) (*connect.Response[pb.ListActivitiesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListExecutions lists executions for a state machine.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListExecutions(ctx context.Context, req *connect.Request[pb.ListExecutionsInput]) (*connect.Response[pb.ListExecutionsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListMapRuns lists Map Runs for a state machine.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListMapRuns(ctx context.Context, req *connect.Request[pb.ListMapRunsInput]) (*connect.Response[pb.ListMapRunsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListStateMachineAliases lists aliases for a state machine.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListStateMachineAliases(ctx context.Context, req *connect.Request[pb.ListStateMachineAliasesInput]) (*connect.Response[pb.ListStateMachineAliasesOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListStateMachineVersions lists versions of a state machine.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListStateMachineVersions(ctx context.Context, req *connect.Request[pb.ListStateMachineVersionsInput]) (*connect.Response[pb.ListStateMachineVersionsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForResource lists tags for a Step Functions resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceInput]) (*connect.Response[pb.ListTagsForResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PublishStateMachineVersion publishes a version of a state machine.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) PublishStateMachineVersion(ctx context.Context, req *connect.Request[pb.PublishStateMachineVersionInput]) (*connect.Response[pb.PublishStateMachineVersionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// RedriveExecution redrives a Map Run execution.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) RedriveExecution(ctx context.Context, req *connect.Request[pb.RedriveExecutionInput]) (*connect.Response[pb.RedriveExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// SendTaskFailure reports a task failure to Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) SendTaskFailure(ctx context.Context, req *connect.Request[pb.SendTaskFailureInput]) (*connect.Response[pb.SendTaskFailureOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// SendTaskHeartbeat sends a task heartbeat to Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) SendTaskHeartbeat(ctx context.Context, req *connect.Request[pb.SendTaskHeartbeatInput]) (*connect.Response[pb.SendTaskHeartbeatOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// SendTaskSuccess reports a task success to Step Functions.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) SendTaskSuccess(ctx context.Context, req *connect.Request[pb.SendTaskSuccessInput]) (*connect.Response[pb.SendTaskSuccessOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// StartExecution starts a new execution of a state machine.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) StartExecution(ctx context.Context, req *connect.Request[pb.StartExecutionInput]) (*connect.Response[pb.StartExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// StartSyncExecution starts a synchronous execution of a state machine.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) StartSyncExecution(ctx context.Context, req *connect.Request[pb.StartSyncExecutionInput]) (*connect.Response[pb.StartSyncExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// StopExecution stops a running execution.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) StopExecution(ctx context.Context, req *connect.Request[pb.StopExecutionInput]) (*connect.Response[pb.StopExecutionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// TagResource adds tags to a Step Functions resource.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceInput]) (*connect.Response[pb.TagResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// TestState tests a state machine definition.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestState(ctx context.Context, req *connect.Request[pb.TestStateInput]) (*connect.Response[pb.TestStateOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagResource removes tags from a Step Functions resource.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceInput]) (*connect.Response[pb.UntagResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UpdateMapRun updates a Map Run.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) UpdateMapRun(ctx context.Context, req *connect.Request[pb.UpdateMapRunInput]) (*connect.Response[pb.UpdateMapRunOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UpdateStateMachine updates a state machine.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) UpdateStateMachine(ctx context.Context, req *connect.Request[pb.UpdateStateMachineInput]) (*connect.Response[pb.UpdateStateMachineOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UpdateStateMachineAlias updates a state machine alias.
// Use the AWS REST API for write operations as gRPC-Web does not support it.
func (h *AdminHandler) UpdateStateMachineAlias(ctx context.Context, req *connect.Request[pb.UpdateStateMachineAliasInput]) (*connect.Response[pb.UpdateStateMachineAliasOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// ValidateStateMachineDefinition validates a state machine definition.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ValidateStateMachineDefinition(ctx context.Context, req *connect.Request[pb.ValidateStateMachineDefinitionInput]) (*connect.Response[pb.ValidateStateMachineDefinitionOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
