package eventbridge

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/events"
	eventsconnect "vorpalstacks/internal/pb/aws/events/eventsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	eventbridgestore "vorpalstacks/internal/store/aws/eventbridge"
)

// AdminHandler provides EventBridge service administration functionality.
// It implements the CloudWatchEventsServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	eventsconnect.UnimplementedCloudWatchEventsServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ eventsconnect.CloudWatchEventsServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new EventBridge AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getStoreByRegion returns the EventBridge store for the specified region.
func (h *AdminHandler) getStoreByRegion(region string) (*eventbridgestore.EventsStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return eventbridgestore.NewEventsStore(regionStorage, h.accountId, region), nil
}

// ListEventBuses lists EventBridge event buses.
// It returns event buses matching the specified name prefix with pagination support.
func (h *AdminHandler) ListEventBuses(ctx context.Context, req *connect.Request[pb.ListEventBusesRequest]) (*connect.Response[pb.ListEventBusesResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := store.ListEventBuses(ctx, req.Msg.GetNameprefix(), req.Msg.GetLimit(), req.Msg.GetNexttoken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	eventBuses := make([]*pb.EventBus, len(result.EventBuses))
	for i, eb := range result.EventBuses {
		eventBuses[i] = toPbEventBus(eb)
	}

	return connect.NewResponse(&pb.ListEventBusesResponse{
		Eventbuses: eventBuses,
		Nexttoken:  result.NextToken,
	}), nil
}

// ListRules lists EventBridge rules.
// It returns rules matching the specified name prefix for an event bus with pagination support.
func (h *AdminHandler) ListRules(ctx context.Context, req *connect.Request[pb.ListRulesRequest]) (*connect.Response[pb.ListRulesResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	eventBusName := req.Msg.GetEventbusname()
	if eventBusName == "" {
		eventBusName = "default"
	}

	result, err := store.ListRules(ctx, eventBusName, req.Msg.GetNameprefix(), req.Msg.GetLimit(), req.Msg.GetNexttoken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	rules := make([]*pb.Rule, len(result.Rules))
	for i, r := range result.Rules {
		rules[i] = toPbRule(r)
	}

	return connect.NewResponse(&pb.ListRulesResponse{
		Rules:     rules,
		Nexttoken: result.NextToken,
	}), nil
}

// toPbEventBus converts an internal EventBus to a protobuf EventBus.
func toPbEventBus(eb *eventbridgestore.EventBus) *pb.EventBus {
	return &pb.EventBus{
		Name:   eb.Name,
		Arn:    eb.ARN,
		Policy: eb.Policy,
	}
}

// toPbRule converts an internal Rule to a protobuf Rule.
func toPbRule(r *eventbridgestore.Rule) *pb.Rule {
	return &pb.Rule{
		Name:               r.Name,
		Arn:                r.ARN,
		Eventbusname:       r.EventBusName,
		Description:        r.Description,
		Eventpattern:       r.EventPattern,
		Scheduleexpression: r.ScheduleExpression,
		State:              toPbRuleState(r.State),
		Managedby:          r.ManagedBy,
		Rolearn:            r.RoleARN,
	}
}

func toPbRuleState(state eventbridgestore.RuleState) pb.RuleState {
	switch state {
	case eventbridgestore.RuleStateEnabled:
		return pb.RuleState_RULE_STATE_ENABLED
	default:
		return pb.RuleState_RULE_STATE_DISABLED
	}
}

// ActivateEventSource activates an event source.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ActivateEventSource(ctx context.Context, req *connect.Request[pb.ActivateEventSourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelReplay cancels a replay.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelReplay(ctx context.Context, req *connect.Request[pb.CancelReplayRequest]) (*connect.Response[pb.CancelReplayResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateApiDestination creates an API destination.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateApiDestination(ctx context.Context, req *connect.Request[pb.CreateApiDestinationRequest]) (*connect.Response[pb.CreateApiDestinationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// CreateArchive creates an archive.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateArchive(ctx context.Context, req *connect.Request[pb.CreateArchiveRequest]) (*connect.Response[pb.CreateArchiveResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// CreateConnection creates a connection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateConnection(ctx context.Context, req *connect.Request[pb.CreateConnectionRequest]) (*connect.Response[pb.CreateConnectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// CreateEventBus creates an event bus.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateEventBus(ctx context.Context, req *connect.Request[pb.CreateEventBusRequest]) (*connect.Response[pb.CreateEventBusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// CreatePartnerEventSource creates a partner event source.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePartnerEventSource(ctx context.Context, req *connect.Request[pb.CreatePartnerEventSourceRequest]) (*connect.Response[pb.CreatePartnerEventSourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeactivateEventSource deactivates an event source.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeactivateEventSource(ctx context.Context, req *connect.Request[pb.DeactivateEventSourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeauthorizeConnection deauthorises a connection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeauthorizeConnection(ctx context.Context, req *connect.Request[pb.DeauthorizeConnectionRequest]) (*connect.Response[pb.DeauthorizeConnectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteApiDestination deletes an API destination.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteApiDestination(ctx context.Context, req *connect.Request[pb.DeleteApiDestinationRequest]) (*connect.Response[pb.DeleteApiDestinationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteArchive deletes an archive.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteArchive(ctx context.Context, req *connect.Request[pb.DeleteArchiveRequest]) (*connect.Response[pb.DeleteArchiveResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteConnection deletes a connection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteConnection(ctx context.Context, req *connect.Request[pb.DeleteConnectionRequest]) (*connect.Response[pb.DeleteConnectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteEventBus deletes an event bus.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteEventBus(ctx context.Context, req *connect.Request[pb.DeleteEventBusRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeletePartnerEventSource deletes a partner event source.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePartnerEventSource(ctx context.Context, req *connect.Request[pb.DeletePartnerEventSourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DeleteRule deletes a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRule(ctx context.Context, req *connect.Request[pb.DeleteRuleRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// DescribeApiDestination describes an API destination.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeApiDestination(ctx context.Context, req *connect.Request[pb.DescribeApiDestinationRequest]) (*connect.Response[pb.DescribeApiDestinationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeArchive describes an archive.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeArchive(ctx context.Context, req *connect.Request[pb.DescribeArchiveRequest]) (*connect.Response[pb.DescribeArchiveResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeConnection describes a connection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeConnection(ctx context.Context, req *connect.Request[pb.DescribeConnectionRequest]) (*connect.Response[pb.DescribeConnectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeEventBus describes an event bus.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeEventBus(ctx context.Context, req *connect.Request[pb.DescribeEventBusRequest]) (*connect.Response[pb.DescribeEventBusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeEventSource describes an event source.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeEventSource(ctx context.Context, req *connect.Request[pb.DescribeEventSourceRequest]) (*connect.Response[pb.DescribeEventSourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribePartnerEventSource describes a partner event source.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribePartnerEventSource(ctx context.Context, req *connect.Request[pb.DescribePartnerEventSourceRequest]) (*connect.Response[pb.DescribePartnerEventSourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeReplay describes a replay.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeReplay(ctx context.Context, req *connect.Request[pb.DescribeReplayRequest]) (*connect.Response[pb.DescribeReplayResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeRule describes a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeRule(ctx context.Context, req *connect.Request[pb.DescribeRuleRequest]) (*connect.Response[pb.DescribeRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableRule disables a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableRule(ctx context.Context, req *connect.Request[pb.DisableRuleRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableRule enables a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableRule(ctx context.Context, req *connect.Request[pb.EnableRuleRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListApiDestinations lists API destinations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListApiDestinations(ctx context.Context, req *connect.Request[pb.ListApiDestinationsRequest]) (*connect.Response[pb.ListApiDestinationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListArchives lists archives.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListArchives(ctx context.Context, req *connect.Request[pb.ListArchivesRequest]) (*connect.Response[pb.ListArchivesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListConnections lists connections.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListConnections(ctx context.Context, req *connect.Request[pb.ListConnectionsRequest]) (*connect.Response[pb.ListConnectionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListEventSources lists event sources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListEventSources(ctx context.Context, req *connect.Request[pb.ListEventSourcesRequest]) (*connect.Response[pb.ListEventSourcesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPartnerEventSourceAccounts lists partner event source accounts.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPartnerEventSourceAccounts(ctx context.Context, req *connect.Request[pb.ListPartnerEventSourceAccountsRequest]) (*connect.Response[pb.ListPartnerEventSourceAccountsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListPartnerEventSources lists partner event sources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPartnerEventSources(ctx context.Context, req *connect.Request[pb.ListPartnerEventSourcesRequest]) (*connect.Response[pb.ListPartnerEventSourcesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListReplays lists replays.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListReplays(ctx context.Context, req *connect.Request[pb.ListReplaysRequest]) (*connect.Response[pb.ListReplaysResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListRuleNamesByTarget lists rule names by target.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRuleNamesByTarget(ctx context.Context, req *connect.Request[pb.ListRuleNamesByTargetRequest]) (*connect.Response[pb.ListRuleNamesByTargetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForResource lists tags for a resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTargetsByRule lists targets by rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTargetsByRule(ctx context.Context, req *connect.Request[pb.ListTargetsByRuleRequest]) (*connect.Response[pb.ListTargetsByRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutEvents puts events to EventBridge.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutEvents(ctx context.Context, req *connect.Request[pb.PutEventsRequest]) (*connect.Response[pb.PutEventsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutPartnerEvents puts partner events.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutPartnerEvents(ctx context.Context, req *connect.Request[pb.PutPartnerEventsRequest]) (*connect.Response[pb.PutPartnerEventsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutPermission adds permissions to an event bus.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutPermission(ctx context.Context, req *connect.Request[pb.PutPermissionRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutRule creates or updates a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutRule(ctx context.Context, req *connect.Request[pb.PutRuleRequest]) (*connect.Response[pb.PutRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// PutTargets adds targets to a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutTargets(ctx context.Context, req *connect.Request[pb.PutTargetsRequest]) (*connect.Response[pb.PutTargetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// RemovePermission removes permissions from an event bus.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemovePermission(ctx context.Context, req *connect.Request[pb.RemovePermissionRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// RemoveTargets removes targets from a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveTargets(ctx context.Context, req *connect.Request[pb.RemoveTargetsRequest]) (*connect.Response[pb.RemoveTargetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// StartReplay starts a replay.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) StartReplay(ctx context.Context, req *connect.Request[pb.StartReplayRequest]) (*connect.Response[pb.StartReplayResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// TagResource adds tags to a resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pb.TagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// TestEventPattern tests an event pattern.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestEventPattern(ctx context.Context, req *connect.Request[pb.TestEventPatternRequest]) (*connect.Response[pb.TestEventPatternResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagResource removes tags from a resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pb.UntagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UpdateApiDestination updates an API destination.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateApiDestination(ctx context.Context, req *connect.Request[pb.UpdateApiDestinationRequest]) (*connect.Response[pb.UpdateApiDestinationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UpdateArchive updates an archive.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateArchive(ctx context.Context, req *connect.Request[pb.UpdateArchiveRequest]) (*connect.Response[pb.UpdateArchiveResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}

// UpdateConnection updates a connection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateConnection(ctx context.Context, req *connect.Request[pb.UpdateConnectionRequest]) (*connect.Response[pb.UpdateConnectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for write operations"))
}
