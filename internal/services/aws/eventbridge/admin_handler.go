package eventbridge

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cloudwatchevents"
	cloudwatcheventsconnect "vorpalstacks/internal/pb/aws/cloudwatchevents/cloudwatcheventsconnect"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// AdminHandler implements the EventBridge (CloudWatch Events) gRPC-Web admin
// console handler. It exposes list operations for event buses and rules for
// the Flutter management UI, delegating to the shared EventsService store.
type AdminHandler struct {
	cloudwatcheventsconnect.UnimplementedCloudWatchEventsServiceHandler
	service        *EventsService
	storageManager *storage.RegionStorageManager
}

var _ cloudwatcheventsconnect.CloudWatchEventsServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new EventBridge admin console handler backed by
// the given service instance, ensuring the same per-region cached stores are used
// as the HTTP API handlers.
func NewAdminHandler(svc *EventsService, sm *storage.RegionStorageManager) *AdminHandler {
	return &AdminHandler{service: svc, storageManager: sm}
}

func (h *AdminHandler) getStore(header http.Header) (*eventsstore.EventsStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	if cached, ok := h.service.eventsStores.Load(region); ok {
		if typed, ok := cached.(*eventsstore.EventsStore); ok {
			return typed, nil
		}
	}
	if h.storageManager == nil {
		return nil, fmt.Errorf("storage manager not initialised")
	}
	st, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	es := eventsstore.NewEventsStore(st, h.service.accountID, region)
	if actual, loaded := h.service.eventsStores.LoadOrStore(region, es); loaded {
		if typed, ok := actual.(*eventsstore.EventsStore); ok {
			return typed, nil
		}
	}
	return es, nil
}

// ListEventBuses returns a paginated list of event buses in the requested region.
func (h *AdminHandler) ListEventBuses(ctx context.Context, req *connect.Request[pb.ListEventBusesRequest]) (*connect.Response[pb.ListEventBusesResponse], error) {
	store, err := h.getStore(req.Header())
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

// ListRules returns a paginated list of rules in the specified event bus.
func (h *AdminHandler) ListRules(ctx context.Context, req *connect.Request[pb.ListRulesRequest]) (*connect.Response[pb.ListRulesResponse], error) {
	store, err := h.getStore(req.Header())
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

func toPbEventBus(eb *eventsstore.EventBus) *pb.EventBus {
	return &pb.EventBus{
		Name:   eb.Name,
		Arn:    eb.ARN,
		Policy: eb.Policy,
	}
}

func toPbRule(r *eventsstore.Rule) *pb.Rule {
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

func toPbRuleState(state eventsstore.RuleState) pb.RuleState {
	switch state {
	case eventsstore.RuleStateEnabled:
		return pb.RuleState_RULE_STATE_ENABLED
	default:
		return pb.RuleState_RULE_STATE_DISABLED
	}
}

func NewConnectHandler(svc *EventsService, sm *storage.RegionStorageManager) (string, http.Handler) {
	return cloudwatcheventsconnect.NewCloudWatchEventsServiceHandler(NewAdminHandler(svc, sm))
}
