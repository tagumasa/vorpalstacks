package sns

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/sns"
	snsconnect "vorpalstacks/internal/pb/aws/sns/snsconnect"
)

// AdminHandler implements the SNS gRPC-Web admin console handler. It is a
// thin adapter: every operation delegates to the service-layer Core methods,
// ensuring identical validation to the HTTP API. No store packages are
// imported directly (store-import prohibition).
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

// ListTopics retrieves all SNS topics from the regional store with pagination.
func (h *AdminHandler) ListTopics(ctx context.Context, req *connect.Request[pb.ListTopicsInput]) (*connect.Response[pb.ListTopicsResponse], error) {
	store, err := h.getSNSStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.listTopicsCore(store, ListTopicsInput{NextToken: req.Msg.GetNexttoken()})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbListTopicsResponse(result)), nil
}

// CreateTopic creates a new SNS topic via the admin console.
// L6: FifoTopic, Tags, and KmsMasterKeyId attributes are now fully
// supported because the handler delegates to createTopicCore, which
// performs the same validation and attribute processing as the HTTP API.
func (h *AdminHandler) CreateTopic(ctx context.Context, req *connect.Request[pb.CreateTopicInput]) (*connect.Response[pb.CreateTopicResponse], error) {
	store, err := h.getSNSStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tags := make(map[string]string, len(req.Msg.Tags))
	for _, tag := range req.Msg.Tags {
		tags[tag.Key] = tag.Value
	}

	result, err := h.service.createTopicCore(store, CreateTopicInput{
		Name:       req.Msg.GetName(),
		Attributes: req.Msg.GetAttributes(),
		Tags:       tags,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbCreateTopicResponse(result)), nil
}

// DeleteTopic deletes an SNS topic via the admin console.
func (h *AdminHandler) DeleteTopic(ctx context.Context, req *connect.Request[pb.DeleteTopicInput]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getSNSStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteTopicCore(store, req.Msg.GetTopicarn()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the SNS admin console.
func NewConnectHandler(svc *SNSService) (string, http.Handler) {
	return snsconnect.NewSNSServiceHandler(NewAdminHandler(svc))
}
