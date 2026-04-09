package sqs

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/sqs"
	sqsconnect "vorpalstacks/internal/pb/aws/sqs/sqsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// AdminHandler implements the SQS gRPC-Web admin console handler. It delegates
// to the shared SQSService to ensure the same per-region cached stores are used
// as the HTTP API handlers.
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

func (h *AdminHandler) getSQSStoreByRegion(region string) (sqsstore.SQSStoreInterface, error) {
	return h.service.GetStoreForRegion(region)
}

// ListQueues returns a paginated list of SQS queue URLs via the admin console gRPC-Web interface.
func (h *AdminHandler) ListQueues(ctx context.Context, req *connect.Request[pb.ListQueuesRequest]) (*connect.Response[pb.ListQueuesResult], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSQSStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	opts := common.ListOptions{
		MaxItems: int(req.Msg.Maxresults),
		Marker:   req.Msg.Nexttoken,
	}

	if req.Msg.Queuenameprefix != "" {
		opts.Prefix = req.Msg.Queuenameprefix
	}

	result, err := store.ListQueues(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	queueUrls := make([]string, 0, len(result.Items))
	for _, queue := range result.Items {
		if req.Msg.Queuenameprefix == "" || strings.HasPrefix(queue.Name, req.Msg.Queuenameprefix) {
			queueUrls = append(queueUrls, queue.URL)
		}
	}

	return connect.NewResponse(&pb.ListQueuesResult{
		Queueurls: queueUrls,
		Nexttoken: result.NextMarker,
	}), nil
}

// GetQueueUrl returns the URL for the specified queue via the admin console gRPC-Web interface.
func (h *AdminHandler) GetQueueUrl(ctx context.Context, req *connect.Request[pb.GetQueueUrlRequest]) (*connect.Response[pb.GetQueueUrlResult], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSQSStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if req.Msg.Queuename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("QueueName is required"))
	}

	queue, err := store.GetQueueByName(req.Msg.Queuename)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("queue not found"))
	}

	return connect.NewResponse(&pb.GetQueueUrlResult{
		Queueurl: queue.URL,
	}), nil
}
