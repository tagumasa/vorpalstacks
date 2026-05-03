package sns

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/sns"
	snsconnect "vorpalstacks/internal/pb/aws/sns/snsconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// AdminHandler implements the SNS gRPC-Web admin console handler. It delegates
// to the shared SNSService to ensure the same per-region cached stores are used
// as the HTTP API handlers.
type AdminHandler struct {
	snsconnect.UnimplementedSNSServiceHandler
	service *SNSService
}

var _ snsconnect.SNSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SNS admin console handler backed by the given
// service instance.
func NewAdminHandler(svc *SNSService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getSNSStoreByRegion(region string) (snsstore.SNSStoreInterface, error) {
	return h.service.getSNSStoreByRegion(region)
}

// ListTopics retrieves all SNS topics from the regional store with pagination.
func (h *AdminHandler) ListTopics(ctx context.Context, req *connect.Request[pb.ListTopicsInput]) (*connect.Response[pb.ListTopicsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSNSStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	opts := storecommon.ListOptions{
		Marker:   req.Msg.Nexttoken,
		MaxItems: 100,
	}

	result, err := store.ListTopics(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	topics := make([]*pb.Topic, len(result.Items))
	for i, t := range result.Items {
		topics[i] = &pb.Topic{
			Topicarn: t.Arn,
		}
	}

	return connect.NewResponse(&pb.ListTopicsResponse{
		Topics:    topics,
		Nexttoken: result.NextMarker,
	}), nil
}

// CreateTopic creates a new SNS topic via the admin console.
func (h *AdminHandler) CreateTopic(ctx context.Context, req *connect.Request[pb.CreateTopicInput]) (*connect.Response[pb.CreateTopicResponse], error) {
	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Name is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSNSStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := store.CreateTopic(&snsstore.Topic{Name: req.Msg.GetName()})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateTopicResponse{
		Topicarn: result.Arn,
	}), nil
}

// DeleteTopic deletes an SNS topic via the admin console.
func (h *AdminHandler) DeleteTopic(ctx context.Context, req *connect.Request[pb.DeleteTopicInput]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.GetTopicarn() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TopicArn is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getSNSStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := store.DeleteTopic(req.Msg.GetTopicarn()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Sns admin console.
func NewConnectHandler(svc *SNSService) (string, http.Handler) {
	return snsconnect.NewSNSServiceHandler(NewAdminHandler(svc))
}
