package scheduler

import (
	"context"
	"google.golang.org/protobuf/proto"
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
// This file has ZERO store package imports (store-import prohibition). Store type
// conversions live in schedule_admin_core.go (proto→store) and
// admin_handler_convert.go (store→proto, getStore).
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

	result, err := h.service.listSchedulesCore(ctx, store, &ListSchedulesInput{
		GroupName:  req.Msg.GetGroupname(),
		NamePrefix: req.Msg.GetNameprefix(),
		State:      req.Msg.GetState(),
		MaxResults: req.Msg.Maxresults,
		NextToken:  req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var summaries []*pb.ScheduleSummary
	for _, s := range result.Schedules {
		summary := &pb.ScheduleSummary{
			Arn:       proto.String(s.Arn),
			Name:      proto.String(s.Name),
			Groupname: proto.String(s.GroupName),
			State:     proto.String(string(s.State)),
		}
		if s.CreationDate != nil {
			summary.Creationdate = proto.String(s.CreationDate.Format(timeutils.ISO8601UTCFormat))
		}
		if s.LastModificationDate != nil {
			summary.Lastmodificationdate = proto.String(s.LastModificationDate.Format(timeutils.ISO8601UTCFormat))
		}
		if s.Target != nil {
			summary.Target = &pb.TargetSummary{Arn: s.Target.Arn}
		}
		summaries = append(summaries, summary)
	}

	out := &pb.ListSchedulesOutput{
		Schedules: summaries,
	}
	// The model documents NextToken as "If the value is null, there are no
	// more results": an exhausted listing omits the member instead of
	// emitting an empty token.
	if result.NextToken != "" {
		out.Nexttoken = proto.String(result.NextToken)
	}
	return connect.NewResponse(out), nil
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

	var iamValidator *iam.IAMValidator
	rp := h.service.RoleProvider()
	if rp != nil {
		iamValidator = iam.NewIAMValidator(rp, h.service.AccountID())
	}

	result, err := h.service.createScheduleFromAdmin(ctx, store, AdminCreateScheduleInput{
		Name:                       req.Msg.Name,
		GroupName:                  req.Msg.GetGroupname(),
		ScheduleExpression:         req.Msg.Scheduleexpression,
		ScheduleExpressionTimezone: req.Msg.GetScheduleexpressiontimezone(),
		Description:                req.Msg.GetDescription(),
		State:                      req.Msg.GetState(),
		KmsKeyArn:                  req.Msg.GetKmskeyarn(),
		StartDate:                  req.Msg.GetStartdate(),
		EndDate:                    req.Msg.GetEnddate(),
		ActionAfterCompletion:      req.Msg.GetActionaftercompletion(),
		Target:                     req.Msg.Target,
		FlexibleTimeWindow:         req.Msg.Flexibletimewindow,
		ClientToken:                req.Msg.GetClienttoken(),
		Region:                     store.GetRegion(),
		IAMValidator:               iamValidator,
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
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteScheduleCore(ctx, store, &DeleteScheduleInput{
		Name:      req.Msg.Name,
		GroupName: req.Msg.GetGroupname(),
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
