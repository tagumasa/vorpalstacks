package scheduler

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/scheduler"
	schedulerconnect "vorpalstacks/internal/pb/aws/scheduler/schedulerconnect"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

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

	maxResults := req.Msg.GetMaxresults()
	if maxResults <= 0 {
		maxResults = 50
	}
	if maxResults > 100 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("max-results must be between 1 and 100"))
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
// Uses the same shared validation layer as the HTTP API (H1).
func (h *AdminHandler) CreateSchedule(ctx context.Context, req *connect.Request[pb.CreateScheduleInput]) (*connect.Response[pb.CreateScheduleOutput], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Convert the protobuf FlexibleTimeWindow to the store type.
	var ftw *schedulerstore.FlexibleTimeWindow
	if req.Msg.Flexibletimewindow != nil {
		ftw = &schedulerstore.FlexibleTimeWindow{
			Mode: schedulerstore.FlexibleTimeWindowMode(req.Msg.Flexibletimewindow.Mode),
		}
	}
	if ftw == nil {
		ftw = &schedulerstore.FlexibleTimeWindow{Mode: schedulerstore.FlexibleTimeWindowModeOff}
	}

	// Convert the protobuf Target to the store type.
	var target *schedulerstore.Target
	if req.Msg.Target != nil {
		target = &schedulerstore.Target{
			Arn:   req.Msg.Target.Arn,
			Input: req.Msg.Target.Input,
		}
		if req.Msg.Target.Rolearn != "" {
			target.RoleArn = req.Msg.Target.Rolearn
		}
	}

	// Build the common spec and run full validation — same path as the
	// HTTP API (H1). This covers namePattern, ScheduleExpression, Target
	// ARN, RoleArn, State enum, ActionAfterCompletion enum,
	// FlexibleTimeWindow Mode enum, KmsKeyArn ARN, Timezone IANA,
	// Description length, and StartDate/EndDate ordering.
	spec := &ScheduleSpec{
		Name:                       req.Msg.Name,
		GroupName:                  req.Msg.Groupname,
		ScheduleExpression:         req.Msg.Scheduleexpression,
		ScheduleExpressionTimezone: req.Msg.Scheduleexpressiontimezone,
		Description:                req.Msg.Description,
		State:                      req.Msg.State,
		KmsKeyArn:                  req.Msg.Kmskeyarn,
		StartDate:                  req.Msg.Startdate,
		EndDate:                    req.Msg.Enddate,
		ActionAfterCompletion:      req.Msg.Actionaftercompletion,
		Target:                     target,
		FlexibleTimeWindow:         ftw,
	}

	validated, err := ValidateScheduleFields(spec)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Validate RoleArn via IAM (same as HTTP API path).
	if target != nil && target.RoleArn != "" {
		rp := h.service.RoleProvider()
		if rp != nil {
			validator := iam.NewIAMValidator(rp, h.service.AccountID())
			if err := validator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
				return nil, svcerrors.StoreErrorToGRPC(err)
			}
		}
	}

	// Validate VPC config (same as HTTP API path) (Minor 3).
	if err := h.service.validateVpcConfig(ctx, store.GetRegion(), target); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Verify schedule group exists (same as HTTP API path).
	groupName := spec.GroupName
	if groupName == "" {
		groupName = "default"
	}
	if groupName != "default" {
		if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
			if err == schedulerstore.ErrScheduleGroupNotFound {
				return nil, svcerrors.StoreErrorToGRPC(ErrScheduleGroupNotFound)
			}
			return nil, svcerrors.StoreErrorToGRPC(ErrInternalServer)
		}
	}

	schedule := &schedulerstore.Schedule{
		Name:                       spec.Name,
		GroupName:                  groupName,
		ScheduleExpression:         spec.ScheduleExpression,
		Target:                     target,
		FlexibleTimeWindow:         ftw,
		State:                      validated.State,
		ScheduleExpressionTimezone: spec.ScheduleExpressionTimezone,
		Description:                spec.Description,
		KmsKeyArn:                  spec.KmsKeyArn,
		StartDate:                  validated.StartDate,
		EndDate:                    validated.EndDate,
		ActionAfterCompletion:      validated.ActionAfterCompletion,
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
