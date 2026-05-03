package dynamodb

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dynamodbconnect "vorpalstacks/internal/pb/aws/dynamodb/dynamodbconnect"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// AdminHandler implements the DynamoDB admin console gRPC-Web handler.
type AdminHandler struct {
	dynamodbconnect.UnimplementedDynamoDBServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ dynamodbconnect.DynamoDBServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new DynamoDB admin handler with the given storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getStoreFromHeader(header http.Header) (*dynamodbstore.TableStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return dynamodbstore.NewTableStore(regionStorage, h.accountId, region), nil
}

// ListTables retrieves all DynamoDB tables from the regional store with pagination.
func (h *AdminHandler) ListTables(ctx context.Context, req *connect.Request[pb.ListTablesInput]) (*connect.Response[pb.ListTablesOutput], error) {
	tableStore, err := h.getStoreFromHeader(req.Header())
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

// CreateTable creates a new DynamoDB table via the admin console.
func (h *AdminHandler) CreateTable(ctx context.Context, req *connect.Request[pb.CreateTableInput]) (*connect.Response[pb.CreateTableOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

	table, err := store.Create(req.Msg.GetTablename(), keySchema, attrDefs, billingMode, nil, nil, nil, nil, nil, false)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateTableOutput{
		Tabledescription: &pb.TableDescription{
			Tablename: table.Name,
			Tablearn:  table.ARN,
		},
	}), nil
}

// DeleteTable deletes a DynamoDB table via the admin console.
func (h *AdminHandler) DeleteTable(ctx context.Context, req *connect.Request[pb.DeleteTableInput]) (*connect.Response[pb.DeleteTableOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := store.Delete(req.Msg.GetTablename()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteTableOutput{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Dynamodb admin console.
func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return dynamodbconnect.NewDynamoDBServiceHandler(NewAdminHandler(sm, accountID))
}
