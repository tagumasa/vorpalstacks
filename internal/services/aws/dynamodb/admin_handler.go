package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
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

func (h *AdminHandler) getStore(headers http.Header) (dynamodbstore.DynamoDBStoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	txnStorage, ok := regionStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("storage does not support transactions")
	}
	return dynamodbstore.NewDynamoDBStore(txnStorage, h.accountId, region), nil
}

func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesInput]) (*connect.Response[pb.ListTablesOutput], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	marker := req.Msg.Exclusivestarttablename
	limit := 100
	if req.Msg.Limit > 0 {
		limit = int(req.Msg.Limit)
	}

	tables, nextMarker, err := store.Tables().List(marker, limit)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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

func (h *AdminHandler) DescribeTable(ctx context.Context, req *connect.Request[pb.DescribeTableInput]) (*connect.Response[pb.DescribeTableOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	table, err := store.Tables().Get(req.Msg.GetTablename())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeTableOutput{
		Table: storeTableToProtoDescription(table),
	}), nil
}

func (h *AdminHandler) CreateTable(ctx context.Context, req *connect.Request[pb.CreateTableInput]) (*connect.Response[pb.CreateTableOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var keySchema []*dynamodbstore.KeySchemaElement
	for _, ks := range req.Msg.GetKeyschema() {
		kt := "HASH"
		if ks.GetKeytype() == pb.KeyType_KEY_TYPE_RANGE {
			kt = "RANGE"
		}
		keySchema = append(keySchema, &dynamodbstore.KeySchemaElement{
			AttributeName: ks.GetAttributename(),
			KeyType:       dynamodbstore.KeyType(kt),
		})
	}

	hasHash := false
	for _, ks := range keySchema {
		if ks.KeyType == dynamodbstore.KeyTypeHash {
			hasHash = true
			break
		}
	}
	if !hasHash {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("KeySchema must contain at least one HASH key"))
	}

	var attrDefs []*dynamodbstore.AttributeDefinition
	for _, ad := range req.Msg.GetAttributedefinitions() {
		at := "S"
		if ad.GetAttributetype() == pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_N {
			at = "N"
		} else if ad.GetAttributetype() == pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_B {
			at = "B"
		}
		attrDefs = append(attrDefs, &dynamodbstore.AttributeDefinition{
			AttributeName: ad.GetAttributename(),
			AttributeType: dynamodbstore.ScalarAttributeType(at),
		})
	}

	billingMode := dynamodbstore.BillingModePayPerRequest
	if req.Msg.GetBillingmode() == pb.BillingMode_BILLING_MODE_PROVISIONED {
		billingMode = dynamodbstore.BillingModeProvisioned
	}

	table, err := store.Tables().Create(req.Msg.GetTablename(), keySchema, attrDefs, billingMode, nil, nil, nil, nil, nil, false)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateTableOutput{
		Tabledescription: storeTableToProtoDescription(table),
	}), nil
}

func (h *AdminHandler) DeleteTable(ctx context.Context, req *connect.Request[pb.DeleteTableInput]) (*connect.Response[pb.DeleteTableOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tableName := req.Msg.GetTablename()
	var deletedTable *dynamodbstore.Table

	err = store.Update(ctx, func(txn *dynamodbstore.DynamoDBTxn) error {
		table, err := txn.GetTable(tableName)
		if err != nil {
			return err
		}
		if table.DeletionProtectionEnabled {
			return connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("table %s is protected from deletion", tableName))
		}
		deletedTable = table
		return txn.DeleteTableCascade(tableName)
	})
	if err != nil {
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			return nil, connectErr
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	deletedTable.Status = dynamodbstore.TableStatusArchived

	return connect.NewResponse(&pb.DeleteTableOutput{
		Tabledescription: storeTableToProtoDescription(deletedTable),
	}), nil
}

func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return dynamodbconnect.NewDynamoDBServiceHandler(NewAdminHandler(sm, accountID))
}
