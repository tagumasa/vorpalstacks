package neptune

import (
	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/rds"
	rdsconnect "vorpalstacks/internal/pb/aws/rds/rdsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	storeneptune "vorpalstacks/internal/store/aws/neptune"
)

type AdminHandler struct {
	rdsconnect.UnimplementedNeptuneServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ rdsconnect.NeptuneServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(sm *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{storageManager: sm, accountId: accountId}
}

func (h *AdminHandler) getStore(req connect.Request[pb.CreateDBClusterMessage]) (*storeneptune.NeptuneStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	rs, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return storeneptune.NewNeptuneStore(rs), nil
}
