package sqs

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/sqs"
	sqsconnect "vorpalstacks/internal/pb/aws/sqs/sqsconnect"
)

// AdminHandler implements the SQS gRPC-Web admin console handler. It delegates
// to the shared SQSService core methods to ensure the same validation and
// per-region cached stores are used as the HTTP API handlers. Per the
// store-import prohibition, this file has ZERO store package imports — all store access is
// through the Core methods in queue_core.go.
type AdminHandler struct {
	sqsconnect.UnimplementedSQSServiceHandler
	service *SQSService
}

var _ sqsconnect.SQSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SQS admin console handler backed by the given
// service instance.
func NewAdminHandler(svc *SQSService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListQueues returns a paginated list of SQS queue URLs via the admin console.
func (h *AdminHandler) ListQueues(ctx context.Context, req *connect.Request[pb.ListQueuesRequest]) (*connect.Response[pb.ListQueuesResult], error) {
	store, err := h.getQueueStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.listQueuesCore(store, ListQueuesInput{
		QueueNamePrefix: req.Msg.Queuenameprefix,
		MaxResults:      int(req.Msg.GetMaxresults()),
		NextToken:       req.Msg.Nexttoken,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbListQueuesResult(result)), nil
}

// GetQueueUrl returns the URL for the specified queue via the admin console.
func (h *AdminHandler) GetQueueUrl(ctx context.Context, req *connect.Request[pb.GetQueueUrlRequest]) (*connect.Response[pb.GetQueueUrlResult], error) {
	store, err := h.getQueueStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.getQueueUrlCore(store, GetQueueUrlInput{
		QueueName: req.Msg.Queuename,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbGetQueueUrlResult(result)), nil
}

// CreateQueue creates an SQS queue via the admin console.
func (h *AdminHandler) CreateQueue(ctx context.Context, req *connect.Request[pb.CreateQueueRequest]) (*connect.Response[pb.CreateQueueResult], error) {
	store, err := h.getQueueStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.createQueueCore(ctx, store, CreateQueueInput{
		QueueName: req.Msg.Queuename,
		Region:    h.service.regionFromHeader(req.Header()),
		Attrs:     req.Msg.Attributes,
		Tags:      req.Msg.Tags,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbCreateQueueResult(result)), nil
}

// DeleteQueue removes an SQS queue via the admin console.
func (h *AdminHandler) DeleteQueue(ctx context.Context, req *connect.Request[pb.DeleteQueueRequest]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getQueueStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteQueueCore(store, DeleteQueueInput{
		QueueURL: req.Msg.Queueurl,
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the SQS admin console.
func NewConnectHandler(svc *SQSService) (string, http.Handler) {
	return sqsconnect.NewSQSServiceHandler(NewAdminHandler(svc))
}
