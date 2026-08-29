package eventbridge

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cloudwatchevents"
	cloudwatcheventsconnect "vorpalstacks/internal/pb/aws/cloudwatchevents/cloudwatcheventsconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// AdminHandler implements the EventBridge (CloudWatch Events) gRPC-Web admin
// console handler. It is a thin adapter that delegates to service-layer core
// functions, converting between proto messages and transport-agnostic Input
// structs. No store package is imported directly (store-import prohibition).
type AdminHandler struct {
	cloudwatcheventsconnect.UnimplementedCloudWatchEventsServiceHandler
	service *EventsService
}

var _ cloudwatcheventsconnect.CloudWatchEventsServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new EventBridge admin console handler backed by
// the given service instance, ensuring the same per-region cached stores are
// used as the HTTP API handlers.
func NewAdminHandler(svc *EventsService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListEventBuses returns a paginated list of event buses in the requested region.
func (h *AdminHandler) ListEventBuses(ctx context.Context, req *connect.Request[pb.ListEventBusesRequest]) (*connect.Response[pb.ListEventBusesResponse], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.listEventBusesCore(ctx, store, ListEventBusesInput{
		NamePrefix: req.Msg.GetNameprefix(),
		Limit:      req.Msg.GetLimit(),
		NextToken:  req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
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

	// The console cannot distinguish an explicitly empty event bus name from
	// an omitted one, so an empty field addresses the default bus.
	result, err := h.service.listRulesCore(ctx, store, ListRulesInput{
		EventBusName: req.Msg.GetEventbusname(),
		NamePrefix:   req.Msg.GetNameprefix(),
		Limit:        req.Msg.GetLimit(),
		NextToken:    req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
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

// CreateEventBus creates a new custom event bus via the admin console.
func (h *AdminHandler) CreateEventBus(ctx context.Context, req *connect.Request[pb.CreateEventBusRequest]) (*connect.Response[pb.CreateEventBusResponse], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.createEventBusCore(ctx, store, CreateEventBusInput{
		Name: req.Msg.Name,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateEventBusResponse{
		Eventbusarn: result.EventBus.ARN,
	}), nil
}

// DeleteEventBus deletes a custom event bus by name via the admin console.
// It cascades rules, targets, archives, and tags with strict error handling,
// mirroring the HTTP API contract.
func (h *AdminHandler) DeleteEventBus(ctx context.Context, req *connect.Request[pb.DeleteEventBusRequest]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteEventBusCore(ctx, store, DeleteEventBusInput{
		Name: req.Msg.Name,
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the EventBridge admin console.
func NewConnectHandler(svc *EventsService) (string, http.Handler) {
	return cloudwatcheventsconnect.NewCloudWatchEventsServiceHandler(NewAdminHandler(svc))
}
