package dynamodb

import (
	"context"

	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Admin-facing service methods
//
// These methods encapsulate store resolution, proto↔store type conversion,
// and core function invocation. The admin gRPC handler delegates exclusively
// to these methods, keeping admin_handler.go free of any store-layer
// dependency — following the same architectural pattern established by ACM.
// ---------------------------------------------------------------------------

// adminListTables resolves the store for the given region and delegates to
// listTablesCore, returning proto-ready table names.
func (s *DynamoDBService) adminListTables(region, marker string, limit int32) ([]string, string, error) {
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, "", err
	}

	// The proto limit is optional int32: zero means absent, which both
	// planes normalise to the documented maximum before the Core's range
	// validation; any other value is the Core's to accept or reject.
	lim := listTablesMaxLimit
	if limit != 0 {
		lim = int(limit)
	}

	tables, nextMarker, err := s.listTablesCore(store, marker, lim)
	if err != nil {
		return nil, "", err
	}

	names := make([]string, len(tables))
	for i, t := range tables {
		names[i] = t.Name
	}
	return names, nextMarker, nil
}

// adminDescribeTable resolves the store for the given region and delegates to
// describeTableCore, returning a proto TableDescription.
func (s *DynamoDBService) adminDescribeTable(region, tableName string) (*pb.TableDescription, error) {
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	table, err := s.describeTableCore(store, tableName)
	if err != nil {
		return nil, err
	}

	return storeTableToProtoDescription(table), nil
}

// adminCreateTable resolves the store for the given region, converts the
// proto request to the transport-agnostic CreateTableInput, delegates to
// createTableCore, and returns a proto TableDescription.
func (s *DynamoDBService) adminCreateTable(ctx context.Context, region string, req *pb.CreateTableInput) (*pb.TableDescription, error) {
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	keySchema := protoKeySchemaToStore(req.GetKeyschema())
	attrDefs := protoAttrDefsToStore(req.GetAttributedefinitions())

	// The proto enum's values map onto the store values directly. The
	// member carries explicit presence (a non-required enum is emitted
	// optional), so an omitted mode arrives as a nil member: it routes
	// through the same empty value the HTTP plane's omission produces,
	// and the core's shared creation contract applies the documented
	// default (PROVISIONED) and its mode/throughput pairing rules — the
	// same contract the HTTP plane's CreateTable resolves through. An
	// explicit PAY_PER_REQUEST is the zero-value enum constant and is
	// told apart from the omission by the member's presence alone.
	var billingMode dbstore.BillingMode
	switch {
	case req.Billingmode == nil:
		// Omitted member: the shared default applies in the core.
	case req.GetBillingmode() == pb.BillingMode_BILLING_MODE_PROVISIONED:
		billingMode = dbstore.BillingModeProvisioned
	case req.GetBillingmode() == pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST:
		billingMode = dbstore.BillingModePayPerRequest
	}
	var provThroughput *dbstore.ProvisionedThroughput
	if pt := req.GetProvisionedthroughput(); pt != nil {
		provThroughput = &dbstore.ProvisionedThroughput{
			ReadCapacityUnits:  pt.GetReadcapacityunits(),
			WriteCapacityUnits: pt.GetWritecapacityunits(),
		}
	}

	table, err := s.createTableCore(ctx, nil, store, CreateTableInput{
		TableName:             req.GetTablename(),
		KeySchema:             keySchema,
		AttributeDefinitions:  attrDefs,
		BillingMode:           billingMode,
		ProvisionedThroughput: provThroughput,
	})
	if err != nil {
		return nil, err
	}

	return storeTableToProtoDescription(table), nil
}

// adminDeleteTable resolves the store for the given region and delegates to
// deleteTableCore, returning a proto TableDescription.
func (s *DynamoDBService) adminDeleteTable(ctx context.Context, region, tableName string) (*pb.TableDescription, error) {
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	deletedTable, err := s.deleteTableCore(ctx, store, region, tableName)
	if err != nil {
		return nil, err
	}

	return storeTableToProtoDescription(deletedTable), nil
}

// ---------------------------------------------------------------------------
// Proto → store type conversion helpers (used only by admin service methods)
// ---------------------------------------------------------------------------

func protoKeySchemaToStore(pbKS []*pb.KeySchemaElement) []*dbstore.KeySchemaElement {
	if len(pbKS) == 0 {
		return nil
	}
	result := make([]*dbstore.KeySchemaElement, len(pbKS))
	for i, ks := range pbKS {
		kt := "HASH"
		if ks.GetKeytype() == pb.KeyType_KEY_TYPE_RANGE {
			kt = "RANGE"
		}
		result[i] = &dbstore.KeySchemaElement{
			AttributeName: ks.GetAttributename(),
			KeyType:       dbstore.KeyType(kt),
		}
	}
	return result
}

func protoAttrDefsToStore(pbADs []*pb.AttributeDefinition) []*dbstore.AttributeDefinition {
	if len(pbADs) == 0 {
		return nil
	}
	result := make([]*dbstore.AttributeDefinition, len(pbADs))
	for i, ad := range pbADs {
		at := "S"
		if ad.GetAttributetype() == pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_N {
			at = "N"
		} else if ad.GetAttributetype() == pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_B {
			at = "B"
		}
		result[i] = &dbstore.AttributeDefinition{
			AttributeName: ad.GetAttributename(),
			AttributeType: dbstore.ScalarAttributeType(at),
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Admin item service methods
//
// These methods resolve the store, convert proto types, resolve the table
// with the same per-operation semantics the HTTP plane enforces (reads
// accept any table status, writes require ACTIVE), delegate to item core
// functions (which include all side effects: stream capture, Kinesis
// destinations, global table replication), and return proto-ready results.
// ---------------------------------------------------------------------------

// adminResolveTable resolves the store and table for an admin item
// operation. Name validation and table lookup live in describeTableCore;
// requireActive mirrors the HTTP plane's split — write handlers resolve
// through validateAndGetActiveTable, read handlers through
// validateAndGetTable — so both planes reject and accept the same table
// states per operation.
func (s *DynamoDBService) adminResolveTable(region, tableName string, requireActive bool) (dbstore.DynamoDBStoreInterface, *dbstore.Table, error) {
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, nil, err
	}
	table, err := s.describeTableCore(store, tableName)
	if err != nil {
		return nil, nil, err
	}
	if requireActive && table.Status != dbstore.TableStatusActive {
		return nil, nil, ErrTableNotActive
	}
	return store, table, nil
}

// adminGetItem retrieves a single item by primary key for the admin console.
func (s *DynamoDBService) adminGetItem(ctx context.Context, region, tableName string, pbKey map[string]*pb.AttributeValue) (map[string]*pb.AttributeValue, error) {
	store, table, err := s.adminResolveTable(region, tableName, false)
	if err != nil {
		return nil, err
	}

	result, err := s.getItemCore(ctx, store, table, protoAVMapToStore(pbKey), false, "")
	if err != nil {
		if dbstore.IsItemNotFound(err) {
			return map[string]*pb.AttributeValue{}, nil
		}
		return nil, err
	}

	return storeAVMapToProto(result.Item.Attributes), nil
}

// adminScanResult holds the proto-ready result of an admin scan operation.
type adminScanResult struct {
	Items            []*pb.ItemListEntry
	Count            int32
	LastEvaluatedKey map[string]*pb.AttributeValue
}

// adminScan retrieves a paginated list of items for the admin console. It
// delegates to scanCore — the same single Scan implementation the data plane
// uses — supplying a minimal parameter set (table, limit, exclusive start
// key) with filters, projections and parallel segments left unset.
func (s *DynamoDBService) adminScan(ctx context.Context, region, tableName string, limit int32, pbStartKey map[string]*pb.AttributeValue) (*adminScanResult, error) {
	store, table, err := s.adminResolveTable(region, tableName, false)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{"TableName": tableName}
	if limit > 0 {
		params["Limit"] = limit
	}
	if len(pbStartKey) > 0 {
		params["ExclusiveStartKey"] = buildItemResponse(protoAVMapToStore(pbStartKey))
	}

	coreResult, err := s.scanCore(ctx, store, table, params)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.ItemListEntry, len(coreResult.Items))
	for i, item := range coreResult.Items {
		pbItems[i] = &pb.ItemListEntry{
			Value: storeAVMapToProto(item.Attributes),
		}
	}

	result := &adminScanResult{
		Items: pbItems,
		Count: int32(len(coreResult.Items)),
	}
	if coreResult.LastEvaluatedKey != nil {
		result.LastEvaluatedKey = storeAVMapToProto(coreResult.LastEvaluatedKey)
	}

	return result, nil
}

// adminPutItem creates or replaces an item for the admin console.
// It applies all side effects: stream capture, Kinesis destinations, and
// global table replication — identical to the HTTP API path, whose
// validation (item size, key completeness, key types) lives in the shared
// putItemCore.
func (s *DynamoDBService) adminPutItem(ctx context.Context, region, tableName string, pbItem map[string]*pb.AttributeValue) (map[string]*pb.AttributeValue, error) {
	store, table, err := s.adminResolveTable(region, tableName, true)
	if err != nil {
		return nil, err
	}

	result, err := s.putItemCore(ctx, store, region, PutItemCoreInput{
		Table: table,
		Item:  protoAVMapToStore(pbItem),
	})
	if err != nil {
		return nil, err
	}

	return storeAVMapToProto(result.StoredItem.Attributes), nil
}

// adminDeleteItem removes an item for the admin console.
// It applies all side effects: stream capture, Kinesis destinations, and
// global table replication — identical to the HTTP API path, whose key
// validation lives in the shared deleteItemCore.
func (s *DynamoDBService) adminDeleteItem(ctx context.Context, region, tableName string, pbKey map[string]*pb.AttributeValue) error {
	store, table, err := s.adminResolveTable(region, tableName, true)
	if err != nil {
		return err
	}

	_, err = s.deleteItemCore(ctx, store, region, DeleteItemCoreInput{
		Table: table,
		Key:   protoAVMapToStore(pbKey),
	})
	return err
}
