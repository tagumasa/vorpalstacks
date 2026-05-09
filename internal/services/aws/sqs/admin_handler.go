package sqs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/sqs"
	sqsconnect "vorpalstacks/internal/pb/aws/sqs/sqsconnect"
	storecommon "vorpalstacks/internal/store/aws/common"
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	opts := storecommon.ListOptions{
		MaxItems: int(req.Msg.Maxresults),
		Marker:   req.Msg.Nexttoken,
	}
	if opts.MaxItems <= 0 {
		opts.MaxItems = 100
	}

	result, err := store.ListQueues(opts)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	queueUrls := make([]string, 0, len(result.Items))
	for _, queue := range result.Items {
		if req.Msg.Queuenameprefix != "" && !strings.HasPrefix(queue.Name, req.Msg.Queuenameprefix) {
			continue
		}
		queueUrls = append(queueUrls, queue.URL)
	}

	return connect.NewResponse(&pb.ListQueuesResult{
		Queueurls: queueUrls,
		Nexttoken: result.NextMarker,
	}), nil
}

// GetQueueUrl returns the URL for the specified queue via the admin console gRPC-Web interface.
func (h *AdminHandler) GetQueueUrl(ctx context.Context, req *connect.Request[pb.GetQueueUrlRequest]) (*connect.Response[pb.GetQueueUrlResult], error) {
	if req.Msg.Queuename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("QueueName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSQSStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	queue, err := store.GetQueueByName(req.Msg.Queuename)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.GetQueueUrlResult{
		Queueurl: queue.URL,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Sqs admin console.
func NewConnectHandler(svc *SQSService) (string, http.Handler) {
	return sqsconnect.NewSQSServiceHandler(NewAdminHandler(svc))
}

func (h *AdminHandler) CreateQueue(ctx context.Context, req *connect.Request[pb.CreateQueueRequest]) (*connect.Response[pb.CreateQueueResult], error) {
	if req.Msg.Queuename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("QueueName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSQSStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	queue := sqsstore.NewQueue(req.Msg.Queuename, region, store.GetAccountID())

	if err := applyQueueAttributes(req.Msg.Attributes, queue); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if len(req.Msg.Tags) > 0 {
		queue.Tags = req.Msg.Tags
	}

	created, err := store.CreateQueue(queue)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateQueueResult{
		Queueurl: created.URL,
	}), nil
}

func (h *AdminHandler) DeleteQueue(ctx context.Context, req *connect.Request[pb.DeleteQueueRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Queueurl == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("QueueUrl is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSQSStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteQueue(req.Msg.Queueurl); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}
