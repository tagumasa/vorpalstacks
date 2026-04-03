package eventbridge

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/events"
	eventsconnect "vorpalstacks/internal/pb/aws/events/eventsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	eventbridgestore "vorpalstacks/internal/store/aws/eventbridge"
)

type AdminHandler struct {
	eventsconnect.UnimplementedCloudWatchEventsServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ eventsconnect.CloudWatchEventsServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getStoreByRegion(region string) (*eventbridgestore.EventsStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return eventbridgestore.NewEventsStore(regionStorage, h.accountId, region), nil
}

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

func toPbEventBus(eb *eventbridgestore.EventBus) *pb.EventBus {
	return &pb.EventBus{
		Name:   eb.Name,
		Arn:    eb.ARN,
		Policy: eb.Policy,
	}
}

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
