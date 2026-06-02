package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/scheduler"
	schedulerconnect "vorpalstacks/internal/pb/aws/scheduler/schedulerconnect"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// parseTime attempts to parse a timestamp string using common AWS formats.
func parseTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, timeutils.ISO8601UTCFormat, "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", s)
}

// AdminHandler provides EventBridge Scheduler service administration functionality.
// It implements the SchedulerServiceHandler interface for gRPC-Web communication.
// It delegates to the shared SchedulerService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	schedulerconnect.UnimplementedSchedulerServiceHandler
	service *SchedulerService
}

// NewAdminHandler creates a new EventBridge Scheduler AdminHandler backed by
// the given service instance.
func NewAdminHandler(svc *SchedulerService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*schedulerstore.SchedulerStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListSchedules retrieves schedules from the store with optional filtering and pagination.
func (h *AdminHandler) ListSchedules(ctx context.Context, req *connect.Request[pb.ListSchedulesInput]) (*connect.Response[pb.ListSchedulesOutput], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	maxResults := req.Msg.Maxresults
	if maxResults <= 0 {
		maxResults = 50
	}

	result, err := store.ListSchedules(ctx, req.Msg.Groupname, req.Msg.Nameprefix, schedulerstore.ScheduleState(req.Msg.State), maxResults, req.Msg.Nexttoken)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var summaries []*pb.ScheduleSummary
	for _, s := range result.Schedules {
		summary := &pb.ScheduleSummary{
			Arn:       s.Arn,
			Name:      s.Name,
			Groupname: s.GroupName,
			State:     string(s.State),
		}
		if s.CreationDate != nil {
			summary.Creationdate = s.CreationDate.Format(timeutils.ISO8601UTCFormat)
		}
		if s.LastModificationDate != nil {
			summary.Lastmodificationdate = s.LastModificationDate.Format(timeutils.ISO8601UTCFormat)
		}
		if s.Target != nil {
			summary.Target = &pb.TargetSummary{Arn: s.Target.Arn}
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListSchedulesOutput{
		Schedules: summaries,
		Nexttoken: result.NextToken,
	}), nil
}

// CreateSchedule creates a new schedule via the admin console.
func (h *AdminHandler) CreateSchedule(ctx context.Context, req *connect.Request[pb.CreateScheduleInput]) (*connect.Response[pb.CreateScheduleOutput], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	schedule := &schedulerstore.Schedule{
		Name:                       req.Msg.Name,
		GroupName:                  req.Msg.Groupname,
		ScheduleExpression:         req.Msg.Scheduleexpression,
		ScheduleExpressionTimezone: req.Msg.Scheduleexpressiontimezone,
		State:                      schedulerstore.ScheduleState(req.Msg.State),
		Description:                req.Msg.Description,
		KmsKeyArn:                  req.Msg.Kmskeyarn,
		ActionAfterCompletion:      schedulerstore.ActionAfterCompletion(req.Msg.Actionaftercompletion),
	}

	if req.Msg.Startdate != "" {
		t, parseErr := parseTime(req.Msg.Startdate)
		if parseErr != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid StartDate: %w", parseErr))
		}
		schedule.StartDate = &t
	}
	if req.Msg.Enddate != "" {
		t, parseErr := parseTime(req.Msg.Enddate)
		if parseErr != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid EndDate: %w", parseErr))
		}
		schedule.EndDate = &t
	}

	if req.Msg.Flexibletimewindow != nil {
		schedule.FlexibleTimeWindow = &schedulerstore.FlexibleTimeWindow{
			Mode: schedulerstore.FlexibleTimeWindowMode(req.Msg.Flexibletimewindow.Mode),
		}
	}

	if req.Msg.Target != nil {
		schedule.Target = &schedulerstore.Target{
			Arn:   req.Msg.Target.Arn,
			Input: req.Msg.Target.Input,
		}
		if req.Msg.Target.Rolearn != "" {
			schedule.Target.RoleArn = req.Msg.Target.Rolearn
		}
	}

	if err := store.CreateSchedule(ctx, schedule); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateScheduleOutput{
		Schedulearn: schedule.ARN,
	}), nil
}

// DeleteSchedule deletes a schedule via the admin console.
func (h *AdminHandler) DeleteSchedule(ctx context.Context, req *connect.Request[pb.DeleteScheduleInput]) (*connect.Response[pb.DeleteScheduleOutput], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteSchedule(ctx, req.Msg.Groupname, req.Msg.Name); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteScheduleOutput{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Scheduler admin console.
func NewConnectHandler(svc *SchedulerService) (string, http.Handler) {
	return schedulerconnect.NewSchedulerServiceHandler(NewAdminHandler(svc))
}
