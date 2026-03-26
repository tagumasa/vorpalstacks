package sqs

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/sqs"
	sqsconnect "vorpalstacks/internal/pb/aws/sqs/sqsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// AdminHandler provides administrative operations for SQS queues.
// It implements the SQSServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	sqsconnect.UnimplementedSQSServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	baseURL        string
}

var _ sqsconnect.SQSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new AdminHandler instance.
// It initialises the handler with the provided storage manager, account ID, and base URL.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId, baseURL string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
		baseURL:        baseURL,
	}
}

func (h *AdminHandler) getSQSStoreByRegion(region string) (sqsstore.SQSStoreInterface, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return sqsstore.NewSQSStore(regionStorage, h.accountId, region, h.baseURL), nil
}

// ListQueues returns a list of queue URLs for the specified region.
// It supports pagination via the NextToken parameter and filtering by queue name prefix.
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

// GetQueueUrl returns the URL of the specified queue.
// It looks up the queue by name and returns its URL, or an error if the queue does not exist.
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

// AddPermission adds permission to a queue for a specific principal.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddPermission(ctx context.Context, req *connect.Request[pb.AddPermissionRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelMessageMoveTask cancels a message move task that is currently in progress.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelMessageMoveTask(ctx context.Context, req *connect.Request[pb.CancelMessageMoveTaskRequest]) (*connect.Response[pb.CancelMessageMoveTaskResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ChangeMessageVisibility changes the visibility timeout of a single message.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ChangeMessageVisibility(ctx context.Context, req *connect.Request[pb.ChangeMessageVisibilityRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ChangeMessageVisibilityBatch changes the visibility timeout of multiple messages in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ChangeMessageVisibilityBatch(ctx context.Context, req *connect.Request[pb.ChangeMessageVisibilityBatchRequest]) (*connect.Response[pb.ChangeMessageVisibilityBatchResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateQueue creates a new standard or FIFO queue with the specified attributes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateQueue(ctx context.Context, req *connect.Request[pb.CreateQueueRequest]) (*connect.Response[pb.CreateQueueResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteMessage deletes a message from the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteMessage(ctx context.Context, req *connect.Request[pb.DeleteMessageRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteMessageBatch deletes multiple messages from a queue in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteMessageBatch(ctx context.Context, req *connect.Request[pb.DeleteMessageBatchRequest]) (*connect.Response[pb.DeleteMessageBatchResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteQueue deletes the specified queue and all its messages.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteQueue(ctx context.Context, req *connect.Request[pb.DeleteQueueRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetQueueAttributes retrieves the attributes of the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetQueueAttributes(ctx context.Context, req *connect.Request[pb.GetQueueAttributesRequest]) (*connect.Response[pb.GetQueueAttributesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListDeadLetterSourceQueues lists the queues that are associated with a dead-letter queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDeadLetterSourceQueues(ctx context.Context, req *connect.Request[pb.ListDeadLetterSourceQueuesRequest]) (*connect.Response[pb.ListDeadLetterSourceQueuesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListMessageMoveTasks lists the message move tasks for a queue or across all queues.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListMessageMoveTasks(ctx context.Context, req *connect.Request[pb.ListMessageMoveTasksRequest]) (*connect.Response[pb.ListMessageMoveTasksResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListQueueTags lists all tags associated with the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListQueueTags(ctx context.Context, req *connect.Request[pb.ListQueueTagsRequest]) (*connect.Response[pb.ListQueueTagsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PurgeQueue deletes all messages from the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PurgeQueue(ctx context.Context, req *connect.Request[pb.PurgeQueueRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ReceiveMessage retrieves one or more messages from the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ReceiveMessage(ctx context.Context, req *connect.Request[pb.ReceiveMessageRequest]) (*connect.Response[pb.ReceiveMessageResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RemovePermission removes the permission from the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemovePermission(ctx context.Context, req *connect.Request[pb.RemovePermissionRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SendMessage sends a message to the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SendMessage(ctx context.Context, req *connect.Request[pb.SendMessageRequest]) (*connect.Response[pb.SendMessageResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SendMessageBatch sends multiple messages to the specified queue in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SendMessageBatch(ctx context.Context, req *connect.Request[pb.SendMessageBatchRequest]) (*connect.Response[pb.SendMessageBatchResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SetQueueAttributes sets the attributes of the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SetQueueAttributes(ctx context.Context, req *connect.Request[pb.SetQueueAttributesRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// StartMessageMoveTask starts a task that moves messages from one queue to another.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) StartMessageMoveTask(ctx context.Context, req *connect.Request[pb.StartMessageMoveTaskRequest]) (*connect.Response[pb.StartMessageMoveTaskResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagQueue adds tags to the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagQueue(ctx context.Context, req *connect.Request[pb.TagQueueRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagQueue removes tags from the specified queue.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagQueue(ctx context.Context, req *connect.Request[pb.UntagQueueRequest]) (*connect.Response[pbcommon.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
