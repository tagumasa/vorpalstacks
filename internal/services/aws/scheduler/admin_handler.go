package scheduler

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/scheduler"
	schedulerconnect "vorpalstacks/internal/pb/aws/scheduler/schedulerconnect"
)

// AdminHandler provides EventBridge Scheduler service administration
// functionality via gRPC-Web. It is a thin adapter that converts protobuf
// requests to service-layer Input structs, delegates to the Core methods,
// and converts results back to protobuf responses.
//
// This file has ZERO store package imports (AGENTS.md #29). All store
// type conversions are in admin_handler_convert.go.
type AdminHandler struct {
	schedulerconnect.UnimplementedSchedulerServiceHandler
	service *SchedulerService
}

func NewAdminHandler(svc *SchedulerService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListSchedules retrieves schedules with optional filtering and pagination.
func (h *AdminHandler) ListSchedules(ctx context.Context, req *connect.Request[pb.ListSchedulesInput]) (*connect.Response[pb.ListSchedulesOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxResults := req.Msg.GetMaxresults()
	if maxResults <= 0 {
		maxResults = 50
	}
	if maxResults > 100 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("max-results must be between 1 and 100"))
	}

	result, err := h.service.listSchedulesCore(ctx, store, &ListSchedulesInput{
		GroupName:  req.Msg.Groupname,
		NamePrefix: req.Msg.Nameprefix,
		State:      req.Msg.State,
		MaxResults: maxResults,
		NextToken:  req.Msg.Nexttoken,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
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
// Delegates to createScheduleCore, which performs validation, IAM
// validation, VPC validation, group existence check, and ClientToken
// idempotency — identical to the HTTP API path.
func (h *AdminHandler) CreateSchedule(ctx context.Context, req *connect.Request[pb.CreateScheduleInput]) (*connect.Response[pb.CreateScheduleOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

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
		Target:                     protoTargetToStore(req.Msg.Target),
		FlexibleTimeWindow:         protoFTWToStore(req.Msg.Flexibletimewindow),
	}

	var iamValidator *iam.IAMValidator
	rp := h.service.RoleProvider()
	if rp != nil {
		iamValidator = iam.NewIAMValidator(rp, h.service.AccountID())
	}

	result, err := h.service.createScheduleCore(ctx, store, &CreateScheduleInput{
		Spec:         spec,
		ClientToken:  req.Msg.Clienttoken,
		Region:       store.GetRegion(),
		IAMValidator: iamValidator,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateScheduleOutput{
		Schedulearn: result.ScheduleArn,
	}), nil
}

// DeleteSchedule deletes a schedule via the admin console.
func (h *AdminHandler) DeleteSchedule(ctx context.Context, req *connect.Request[pb.DeleteScheduleInput]) (*connect.Response[pb.DeleteScheduleOutput], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	groupName := req.Msg.Groupname
	if groupName == "" {
		groupName = "default"
	}

	if err := h.service.deleteScheduleCore(ctx, store, &DeleteScheduleInput{
		Name:      req.Msg.Name,
		GroupName: groupName,
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteScheduleOutput{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Scheduler
// admin console.
func NewConnectHandler(svc *SchedulerService) (string, http.Handler) {
	return schedulerconnect.NewSchedulerServiceHandler(NewAdminHandler(svc))
}
