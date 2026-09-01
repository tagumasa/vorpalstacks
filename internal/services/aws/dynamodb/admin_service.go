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

	lim := listTablesMaxLimit
	if limit > 0 {
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
func (s *DynamoDBService) adminCreateTable(region string, req *pb.CreateTableInput) (*pb.TableDescription, error) {
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	keySchema := protoKeySchemaToStore(req.GetKeyschema())
	attrDefs := protoAttrDefsToStore(req.GetAttributedefinitions())

	billingMode := dbstore.BillingModePayPerRequest
	var provThroughput *dbstore.ProvisionedThroughput
	if req.GetBillingmode() == pb.BillingMode_BILLING_MODE_PROVISIONED {
		billingMode = dbstore.BillingModeProvisioned
		if pt := req.GetProvisionedthroughput(); pt != nil {
			provThroughput = &dbstore.ProvisionedThroughput{
				ReadCapacityUnits:  pt.GetReadcapacityunits(),
				WriteCapacityUnits: pt.GetWritecapacityunits(),
			}
		}
	}

	table, err := s.createTableCore(store, CreateTableInput{
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

	deletedTable, err := s.deleteTableCore(ctx, store, tableName)
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
// These methods resolve the store, convert proto types, validate the table
// is active, delegate to item core functions (which include all side
// effects: stream capture, Kinesis destinations, global table replication),
// and return proto-ready results.
// ---------------------------------------------------------------------------

// adminGetItem retrieves a single item by primary key for the admin console.
func (s *DynamoDBService) adminGetItem(ctx context.Context, region, tableName string, pbKey map[string]*pb.AttributeValue) (map[string]*pb.AttributeValue, error) {
	if tableName == "" {
		return nil, ErrInvalidParameter
	}
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	table, err := s.describeTableCore(store, tableName)
	if err != nil {
		return nil, err
	}
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	key := protoAVMapToStore(pbKey)
	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	item, err := s.getItemCore(ctx, store, tableName, key)
	if err != nil {
		if dbstore.IsItemNotFound(err) {
			return map[string]*pb.AttributeValue{}, nil
		}
		return nil, err
	}

	return storeAVMapToProto(item.Attributes), nil
}

// adminScanResult holds the proto-ready result of an admin scan operation.
type adminScanResult struct {
	Items            []*pb.ItemListEntry
	Count            int32
	LastEvaluatedKey map[string]*pb.AttributeValue
}

// adminScan retrieves a paginated list of items for the admin console.
// It delegates to scanItemsCore so that the limit-cap logic and store
// access are shared with any other Core caller; no admin code path
// touches the store layer directly.
func (s *DynamoDBService) adminScan(region, tableName string, limit int32, pbStartKey map[string]*pb.AttributeValue) (*adminScanResult, error) {
	if tableName == "" {
		return nil, ErrInvalidParameter
	}
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	// Validate table existence.
	if _, err := s.describeTableCore(store, tableName); err != nil {
		return nil, err
	}

	marker := ""
	if len(pbStartKey) > 0 {
		marker = s.buildItemMarkerFromKey(store, tableName, protoAVMapToStore(pbStartKey))
	}

	coreResult, err := s.scanItemsCore(store, ScanItemsInput{
		TableName: tableName,
		Limit:     int(limit),
		Marker:    marker,
	})
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
	if coreResult.NextMarker != "" && len(coreResult.Items) > 0 {
		lastItem := coreResult.Items[len(coreResult.Items)-1]
		result.LastEvaluatedKey = storeAVMapToProto(lastItem.Key)
	}

	return result, nil
}

// adminPutItem creates or replaces an item for the admin console.
// It applies all side effects: stream capture, Kinesis destinations, and
// global table replication — identical to the HTTP API path.
func (s *DynamoDBService) adminPutItem(ctx context.Context, region, tableName string, pbItem map[string]*pb.AttributeValue) (map[string]*pb.AttributeValue, error) {
	if tableName == "" {
		return nil, ErrInvalidParameter
	}
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return nil, err
	}

	table, err := s.describeTableCore(store, tableName)
	if err != nil {
		return nil, err
	}
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	attrs := protoAVMapToStore(pbItem)

	if itemSize := calculateItemSize(attrs); itemSize > maxItemSizeBytes {
		return nil, ErrInvalidParameter
	}

	key := s.extractKeyFromItem(table, attrs)
	if key == nil {
		return nil, ErrMissingKey
	}
	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	storedItem, _, err := s.putItemCore(ctx, store, region, table, key, attrs, nil)
	if err != nil {
		return nil, err
	}

	return storeAVMapToProto(storedItem.Attributes), nil
}

// adminDeleteItem removes an item for the admin console.
// It applies all side effects: stream capture, Kinesis destinations, and
// global table replication — identical to the HTTP API path.
func (s *DynamoDBService) adminDeleteItem(ctx context.Context, region, tableName string, pbKey map[string]*pb.AttributeValue) error {
	if tableName == "" {
		return ErrInvalidParameter
	}
	store, err := s.GetCachedStoreForRegion(region)
	if err != nil {
		return err
	}

	table, err := s.describeTableCore(store, tableName)
	if err != nil {
		return err
	}
	if table.Status != dbstore.TableStatusActive {
		return ErrTableNotActive
	}

	key := protoAVMapToStore(pbKey)
	if !validateKeyValueNotEmpty(key) {
		return ErrInvalidParameter
	}

	_, err = s.deleteItemCore(ctx, store, region, table, key, nil)
	return err
}

// buildItemMarkerFromKey constructs a pagination marker from a key map.
// This replaces the old admin_handler_items.go buildItemMarker method,
// moving store knowledge into the service layer.
func (s *DynamoDBService) buildItemMarkerFromKey(store dbstore.DynamoDBStoreInterface, tableName string, key map[string]*dbstore.AttributeValue) string {
	table, err := store.Tables().Get(tableName)
	if err != nil {
		return tableName + dbstore.KeySep
	}

	pkName := ""
	skName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			pkName = ks.AttributeName
		} else if ks.KeyType == dbstore.KeyTypeRange {
			skName = ks.AttributeName
		}
	}

	pkValue := avToString(key[pkName])
	if pkValue == "" {
		return tableName + dbstore.KeySep
	}

	if skName != "" {
		if key[skName] != nil {
			skValue := avToString(key[skName])
			if skValue != "" {
				return tableName + dbstore.KeySep + pkValue + dbstore.KeySep + skValue
			}
		}
	}

	return tableName + dbstore.KeySep + pkValue
}

// avToString extracts a string representation of an AttributeValue for
// use in pagination marker construction.
func avToString(av *dbstore.AttributeValue) string {
	if av == nil {
		return ""
	}
	if av.S != nil {
		return *av.S
	}
	if av.N != nil {
		return *av.N
	}
	if av.B != nil {
		return string(av.B)
	}
	return ""
}
