package eventbridge

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/cloudwatchevents"
	cloudwatcheventsconnect "vorpalstacks/internal/pb/aws/cloudwatchevents/cloudwatcheventsconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// AdminHandler implements the EventBridge (CloudWatch Events) gRPC-Web admin
// console handler. It exposes list operations for event buses and rules for
// the Flutter management UI, delegating to the shared EventsService store.
type AdminHandler struct {
	cloudwatcheventsconnect.UnimplementedCloudWatchEventsServiceHandler
	service *EventsService
}

var _ cloudwatcheventsconnect.CloudWatchEventsServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new EventBridge admin console handler backed by
// the given service instance, ensuring the same per-region cached stores are used
// as the HTTP API handlers.
func NewAdminHandler(svc *EventsService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStore(header http.Header) (*eventsstore.EventsStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// ListEventBuses returns a paginated list of event buses in the requested region.
func (h *AdminHandler) ListEventBuses(ctx context.Context, req *connect.Request[pb.ListEventBusesRequest]) (*connect.Response[pb.ListEventBusesResponse], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := store.ListEventBuses(ctx, req.Msg.GetNameprefix(), req.Msg.GetLimit(), req.Msg.GetNexttoken())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	eventBusName := req.Msg.GetEventbusname()
	if eventBusName == "" {
		eventBusName = "default"
	}

	result, err := store.ListRules(ctx, eventBusName, req.Msg.GetNameprefix(), req.Msg.GetLimit(), req.Msg.GetNexttoken())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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

// CreateEventBus creates a new custom event bus via the admin console.
func (h *AdminHandler) CreateEventBus(ctx context.Context, req *connect.Request[pb.CreateEventBusRequest]) (*connect.Response[pb.CreateEventBusResponse], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{
		Name:      req.Msg.Name,
		Region:    region,
		AccountID: h.service.accountID,
	}); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	arn := fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", region, h.service.accountID, req.Msg.Name)
	return connect.NewResponse(&pb.CreateEventBusResponse{Eventbusarn: arn}), nil
}

// DeleteEventBus deletes a custom event bus by name via the admin console.
// It cascades rules, targets, archives, and tags — mirroring the HTTP API handler.
func (h *AdminHandler) DeleteEventBus(ctx context.Context, req *connect.Request[pb.DeleteEventBusRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if _, err := store.GetEventBus(ctx, req.Msg.Name); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Cascade-delete: rules -> targets (paginated), then archives (mirrors HTTP API handler)
	ruleToken := ""
	for {
		rulesResult, err := store.ListRules(ctx, req.Msg.Name, "", 1000, ruleToken)
		if err != nil {
			break
		}
		for _, rule := range rulesResult.Rules {
			targetToken := ""
			for {
				targetsResult, tErr := store.ListTargetsByRule(ctx, req.Msg.Name, rule.Name, 1000, targetToken)
				if tErr != nil {
					break
				}
				for _, t := range targetsResult.Targets {
					_ = store.DeleteTarget(ctx, req.Msg.Name, rule.Name, t.ID)
				}
				if targetsResult.NextToken == "" {
					break
				}
				targetToken = targetsResult.NextToken
			}
			_ = store.DeleteRule(ctx, req.Msg.Name, rule.Name)
			lastFireTimes.Delete(rule.ARN)
			_ = store.TagStore.Delete(rule.ARN)
		}
		if rulesResult.NextToken == "" {
			break
		}
		ruleToken = rulesResult.NextToken
	}

	archives, err := store.ListArchivesForEventBus(ctx, req.Msg.Name)
	if err == nil {
		for _, a := range archives {
			_ = store.DeleteArchiveEvents(ctx, a.Name)
			_ = store.DeleteArchive(ctx, a.Name)
		}
	}

	if err := store.DeleteEventBus(ctx, req.Msg.Name); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the EventBridge admin console.
func NewConnectHandler(svc *EventsService) (string, http.Handler) {
	return cloudwatcheventsconnect.NewCloudWatchEventsServiceHandler(NewAdminHandler(svc))
}
