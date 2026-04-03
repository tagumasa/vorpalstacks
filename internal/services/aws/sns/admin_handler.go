package sns

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/sns"
	snsconnect "vorpalstacks/internal/pb/aws/sns/snsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

type AdminHandler struct {
	snsconnect.UnimplementedSNSServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ snsconnect.SNSServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getSNSStoreByRegion(region string) (*snsstore.SNSStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return snsstore.NewSNSStore(regionStorage, h.accountId, region), nil
}

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
