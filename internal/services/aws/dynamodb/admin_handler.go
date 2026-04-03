package dynamodb

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type AdminHandler struct {
	dynamodbconnect.UnimplementedDynamoDBServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ dynamodbconnect.DynamoDBServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getTableStore(req *connect.Request[pb.ListTablesInput]) (*dynamodbstore.TableStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return dynamodbstore.NewTableStore(regionStorage, h.accountId, region), nil
}

func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesInput]) (*connect.Response[pb.ListTablesOutput], error) {
	tableStore, err := h.getTableStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	marker := ""
	limit := 100
	if req.Msg.Limit > 0 {
		limit = int(req.Msg.Limit)
	}

	tables, nextMarker, err := tableStore.List(marker, limit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var tableNames []string
	for _, t := range tables {
		tableNames = append(tableNames, t.Name)
	}

	return connect.NewResponse(&pb.ListTablesOutput{
		Tablenames:             tableNames,
		Lastevaluatedtablename: nextMarker,
	}), nil
}
