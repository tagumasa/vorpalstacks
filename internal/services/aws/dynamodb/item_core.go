package dynamodb

import (
	"context"
	"errors"

	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Condition checker callbacks
// ---------------------------------------------------------------------------

// ConditionChecker evaluates a condition expression against the existing
// item state within a transaction. It returns nil if the condition is
// satisfied, or ErrConditionalCheckFailed (or any other error) otherwise.
// When existing is nil and isNotFound is true, the caller should evaluate
// the condition against a synthetic empty item.
type ConditionChecker func(existing *dbstore.Item, isNotFound bool) error

// conditionSpec carries a condition expression and its substitution maps —
// the parsed form the item write Cores receive from both planes. An empty
// Expr means the write is unconditional.
type conditionSpec struct {
	Expr   string
	Names  map[string]string
	Values map[string]*dbstore.AttributeValue
}

// buildConditionChecker turns a condition specification into the
// ConditionChecker a write Core evaluates inside its transaction. A missing
// item is evaluated against a synthetic empty item carrying only the key, so
// attribute_not_exists conditions behave as documented. A nil Expr yields a
// nil checker — unconditional write.
func buildConditionChecker(tableName string, key map[string]*dbstore.AttributeValue, cond conditionSpec) ConditionChecker {
	if cond.Expr == "" {
		return nil
	}
	return func(existing *dbstore.Item, isNotFound bool) error {
		evalItem := existing
		if isNotFound {
			evalItem = &dbstore.Item{
				TableName:  tableName,
				Key:        key,
				Attributes: make(map[string]*dbstore.AttributeValue),
			}
		}
		met, err := evaluateConditionExpression(evalItem, cond.Expr, cond.Names, cond.Values)
		if err != nil {
			return err
		}
		if !met {
			return ErrConditionalCheckFailed
		}
		return nil
	}
}

// validateReturnValuesNoneAllOld rejects ReturnValues other than NONE and
// ALL_OLD — the only values PutItem and DeleteItem recognise (model
// ReturnValues member: "PutItem does not recognize any values other than
// NONE or ALL_OLD"). An empty value means the parameter was omitted.
func validateReturnValuesNoneAllOld(returnValues string) bool {
	return returnValues == "" || returnValues == "NONE" || returnValues == "ALL_OLD"
}

// ---------------------------------------------------------------------------
// Core item functions — single validation + persistence path
//
// These functions encapsulate the full item lifecycle including DynamoDB
// Streams capture, Kinesis Data Stream destinations, and global table
// replication. Both the HTTP API handlers and the admin gRPC handler
// delegate to these functions to ensure identical side-effect behaviour.
// ---------------------------------------------------------------------------

// getItemCore validates the key and retrieves a single item by primary key.
// It returns dbstore.ItemNotFound when the item does not exist.
func (s *DynamoDBService) getItemCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, key map[string]*dbstore.AttributeValue) (*dbstore.Item, error) {
	if !validateKeyAttributeValue(key) {
		return nil, ErrInvalidParameter
	}
	if err := validateKeyTypes(table, key); err != nil {
		return nil, err
	}
	item, err := store.Items().Get(table.Name, key)
	if err == nil {
		s.recordContributorReads(ctx, store, table.Name, []map[string]*dbstore.AttributeValue{key})
	}
	return item, err
}

// recordContributorReads counts one read event per key in the contributor
// access aggregation of tables with contributor insights enabled. Reads
// count as one unit of ConsumedThroughputUnits. Failures are ignored:
// monitoring must never fail the read it observes.
func (s *DynamoDBService) recordContributorReads(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string, keys []map[string]*dbstore.AttributeValue) {
	if len(keys) == 0 {
		return
	}
	if err := store.RecordContributorReads(ctx, tableName, keys); err != nil {
		logs.Warn("failed to record contributor reads",
			logs.String("table", tableName), logs.Err(err))
	}
}

// recordQueryContributorEvent counts the single read event a Query
// contributes to the partition-key series of tables with contributor
// insights enabled. Failures are ignored: monitoring must never fail the
// read it observes.
func (s *DynamoDBService) recordQueryContributorEvent(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, pkValue *dbstore.AttributeValue) {
	if !table.ContributorInsightsEnabled || pkValue == nil {
		return
	}
	pkName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			pkName = ks.AttributeName
			break
		}
	}
	if pkName == "" {
		return
	}
	if err := store.RecordContributorQuery(ctx, table.Name, map[string]*dbstore.AttributeValue{pkName: pkValue}); err != nil {
		logs.Warn("failed to record contributor query event",
			logs.String("table", table.Name), logs.Err(err))
	}
}

// PutItemCoreInput is the service-layer DTO for PutItem on both the HTTP API
// and the admin console.
type PutItemCoreInput struct {
	Table        *dbstore.Table
	Item         map[string]*dbstore.AttributeValue
	Condition    conditionSpec
	ReturnValues string
}

// PutItemCoreResult holds the stored item and the item it overwrote.
type PutItemCoreResult struct {
	StoredItem *dbstore.Item
	OldItem    *dbstore.Item
}

// putItemCore validates the write, then creates or replaces an item,
// applying all side effects: DynamoDB Streams capture (within the
// transaction), Kinesis Data Stream destinations (after commit), and global
// table replication. The condition expression is evaluated inside the
// transaction after loading the existing item.
func (s *DynamoDBService) putItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	in PutItemCoreInput,
) (*PutItemCoreResult, error) {
	table := in.Table

	if !validateReturnValuesNoneAllOld(in.ReturnValues) {
		return nil, ErrInvalidParameter
	}
	if itemSize := dbstore.CalculateItemSize(in.Item); itemSize > dbstore.MaxItemSizeBytes {
		return nil, ErrInvalidParameter
	}
	key := s.extractKeyFromItem(table, in.Item)
	if key == nil {
		return nil, ErrMissingKey
	}
	if !validateKeyAttributeValue(key) {
		return nil, ErrInvalidParameter
	}
	// Key attribute types must match the table schema and any index key
	// definitions.
	if err := validateItemKeyTypes(table, in.Item); err != nil {
		return nil, err
	}
	item := in.Item
	conditionChecker := buildConditionChecker(table.Name, key, in.Condition)

	var storedItem, oldItem *dbstore.Item
	var isNew bool

	err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, getErr := txn.GetItem(table.Name, key)
		if getErr != nil {
			if dbstore.IsItemNotFound(getErr) {
				isNew = true
			} else {
				return getErr
			}
		} else {
			oldItem = existingItem
		}

		if conditionChecker != nil {
			if condErr := conditionChecker(existingItem, isNew); condErr != nil {
				return condErr
			}
		}

		storedItem = &dbstore.Item{
			TableName:  table.Name,
			Key:        key,
			Attributes: item,
		}
		var oldSize int64
		if oldItem != nil {
			oldSize = dbstore.CalculateItemSize(oldItem.Attributes)
		}
		if err := txn.StoreItemWrite(table.Name, key, item, oldItem, oldItem != nil, oldSize); err != nil {
			return err
		}

		s.captureStreamChangeTxn(txn, store, table, streamEventForWrite(false, isNew), key, storedItem.Attributes, oldItemAttributes(oldItem))

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.emitChangePropagation(store, region, table, streamEventForWrite(false, isNew), key, storedItem.Attributes, oldItemAttributes(oldItem), s.replicaPutOp(table, key, item))

	return &PutItemCoreResult{StoredItem: storedItem, OldItem: oldItem}, nil
}

// DeleteItemCoreInput is the service-layer DTO for DeleteItem on both the
// HTTP API and the admin console.
type DeleteItemCoreInput struct {
	Table        *dbstore.Table
	Key          map[string]*dbstore.AttributeValue
	Condition    conditionSpec
	ReturnValues string
}

// DeleteItemCoreResult holds the item that was removed (nil when the key did
// not exist).
type DeleteItemCoreResult struct {
	OldItem *dbstore.Item
}

// deleteItemCore validates the key, then removes the item by primary key,
// applying all side effects: DynamoDB Streams capture (within the
// transaction), Kinesis Data Stream destinations (after commit), and global
// table replication. The condition expression is evaluated inside the
// transaction after loading the existing item.
func (s *DynamoDBService) deleteItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	in DeleteItemCoreInput,
) (*DeleteItemCoreResult, error) {
	table := in.Table

	if !validateReturnValuesNoneAllOld(in.ReturnValues) {
		return nil, ErrInvalidParameter
	}
	if !validateKeyAttributeValue(in.Key) {
		return nil, ErrInvalidParameter
	}
	if err := validateKeyTypes(table, in.Key); err != nil {
		return nil, err
	}
	key := in.Key
	conditionChecker := buildConditionChecker(table.Name, key, in.Condition)

	var oldItem *dbstore.Item
	err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, getErr := txn.GetItem(table.Name, key)
		if getErr != nil {
			if dbstore.IsItemNotFound(getErr) {
				if conditionChecker != nil {
					if condErr := conditionChecker(nil, true); condErr != nil {
						return condErr
					}
				}
				return nil
			}
			return getErr
		}

		if conditionChecker != nil {
			if condErr := conditionChecker(existingItem, false); condErr != nil {
				return condErr
			}
		}

		oldItem = existingItem

		if oldItem != nil {
			if err := txn.DeleteItemWrite(table.Name, key, oldItem, true, dbstore.CalculateItemSize(oldItem.Attributes)); err != nil {
				return err
			}
			s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventRemove, key, nil, oldItem.Attributes)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if oldItem != nil {
		s.sendToKinesisDestinations(table, dbstore.StreamEventRemove, key, nil, oldItem.Attributes)
	}

	s.replicateToGlobalTableReplicas(store, region, table.Name, s.replicaDeleteOp(table, key))

	return &DeleteItemCoreResult{OldItem: oldItem}, nil
}

// isItemNotFound reports whether the error is the store's item-not-found
// sentinel, so transport layers can branch on it without importing the store
// package.
func isItemNotFound(err error) bool {
	return dbstore.IsItemNotFound(err)
}

// itemCollectionWriteRef names one item collection a write touched: the
// table record plus the written item's primary key.
type itemCollectionWriteRef struct {
	tableName string
	table     *dbstore.Table
	key       map[string]*dbstore.AttributeValue
}

// itemCollectionKey resolves the partition key name and value of a written
// item from the table's key schema. An empty name or a nil value means no
// item collection key could be derived.
func itemCollectionKey(table *dbstore.Table, key map[string]*dbstore.AttributeValue) (string, *dbstore.AttributeValue) {
	if table == nil {
		return "", nil
	}
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			return ks.AttributeName, key[ks.AttributeName]
		}
	}
	return "", nil
}

// tableHasItemCollections reports whether the table tracks item
// collections: the Developer Guide scopes item collections to tables that
// have one or more local secondary indexes, and the write operations'
// documentation states ItemCollectionMetrics is not returned for a table
// without them.
func tableHasItemCollections(table *dbstore.Table) bool {
	return table != nil && len(table.LocalSecondaryIndexes) > 0
}

// itemCollectionMetricsEntry renders one ItemCollectionMetrics entry: the
// written item's partition key value as the ItemCollectionKey and the
// table's size estimate as the two-element SizeEstimateRangeGB.
func itemCollectionMetricsEntry(table *dbstore.Table, pkName string, pkValue *dbstore.AttributeValue) map[string]interface{} {
	sizeGB := float64(table.TableSizeBytes) / (1024.0 * 1024.0 * 1024.0)
	return map[string]interface{}{
		"ItemCollectionKey":   buildItemResponse(map[string]*dbstore.AttributeValue{pkName: pkValue}),
		"SizeEstimateRangeGB": []float64{sizeGB, sizeGB},
	}
}

// buildItemCollectionMetricsEntry renders the single-object
// ItemCollectionMetrics response member of the single-item writes
// (PutItem, UpdateItem, DeleteItem) for ReturnItemCollectionMetrics=SIZE.
// It returns nil when no entry applies: the table has no local secondary
// indexes or the written key carries no partition key value.
func buildItemCollectionMetricsEntry(table *dbstore.Table, key map[string]*dbstore.AttributeValue) map[string]interface{} {
	if !tableHasItemCollections(table) {
		return nil
	}
	pkName, pkValue := itemCollectionKey(table, key)
	if pkName == "" || pkValue == nil {
		return nil
	}
	return itemCollectionMetricsEntry(table, pkName, pkValue)
}

// buildItemCollectionMetricsPerTable renders the ItemCollectionMetrics
// response member for the SIZE setting of ReturnItemCollectionMetrics on
// the batched write planes (BatchWriteItem, TransactWriteItems): one entry
// per distinct item collection (table plus partition key value), carrying
// the same entry the single-item writes report. A delete of a non-existent
// key is a successful no-op here, and the documented contract — item
// collections "affected by individual DeleteItem or PutItem operations" —
// does not determine whether such a delete affects its collection, so it
// registers like any other successful write.
func buildItemCollectionMetricsPerTable(writes []itemCollectionWriteRef) map[string]interface{} {
	perTable := make(map[string]interface{})
	seen := make(map[string]bool)
	for _, w := range writes {
		if !tableHasItemCollections(w.table) {
			continue
		}
		pkName, pkValue := itemCollectionKey(w.table, w.key)
		if pkName == "" || pkValue == nil {
			continue
		}
		identity := w.tableName + dbstore.EncodeKeyValue(pkValue)
		if seen[identity] {
			continue
		}
		seen[identity] = true
		entry := itemCollectionMetricsEntry(w.table, pkName, pkValue)
		existing, _ := perTable[w.tableName].([]map[string]interface{})
		perTable[w.tableName] = append(existing, entry)
	}
	if len(perTable) == 0 {
		return nil
	}
	return perTable
}

// ---------------------------------------------------------------------------
// Transport-agnostic DTO for UpdateItem
// ---------------------------------------------------------------------------

// UpdateItemInput carries every field that UpdateItem needs. Both the HTTP
// API handler and admin gRPC handler build this struct and delegate to
// updateItemCore.
type UpdateItemInput struct {
	Key            map[string]*dbstore.AttributeValue
	UpdateExpr     string
	AttrUpdates    interface{} // legacy AttributeUpdates map (mutually exclusive with UpdateExpr)
	ConditionExpr  string
	ExprAttrNames  map[string]string
	ExprAttrValues map[string]*dbstore.AttributeValue
	ReturnValues   string
}

// UpdateItemResult holds the output of updateItemCore for response formatting.
type UpdateItemResult struct {
	StoredItem       *dbstore.Item
	OldItem          *dbstore.Item
	UpdatedAttrNames []string
	WasNewItem       bool
}

// updateItemCore is the single entry point for item updates shared by the
// HTTP API and admin gRPC handler. It validates the key, the
// UpdateExpression/AttributeUpdates exclusivity, and the ReturnValues enum,
// then evaluates conditions, applies update expressions, persists within a
// transaction, and fires all side effects (Streams, Kinesis, global table
// replication).
func (s *DynamoDBService) updateItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	in UpdateItemInput,
) (*UpdateItemResult, error) {
	tableName := table.Name

	if !validateKeyAttributeValue(in.Key) {
		return nil, ErrInvalidParameter
	}
	if err := validateKeyTypes(table, in.Key); err != nil {
		return nil, err
	}
	// UpdateExpression and the legacy AttributeUpdates member are mutually
	// exclusive.
	if in.UpdateExpr != "" && in.AttrUpdates != nil {
		return nil, ErrInvalidParameter
	}
	// UpdateItem recognises all five ReturnValue settings (model
	// ReturnValues member); any other value is rejected rather than ignored.
	switch in.ReturnValues {
	case "", "NONE", "ALL_OLD", "UPDATED_OLD", "ALL_NEW", "UPDATED_NEW":
	default:
		return nil, ErrInvalidParameter
	}

	var oldItem *dbstore.Item
	var storedItem *dbstore.Item
	var updatedAttrNames []string
	var oldItemSize int64
	var wasNewItem bool
	conditionChecker := buildConditionChecker(tableName, in.Key, conditionSpec{
		Expr:   in.ConditionExpr,
		Names:  in.ExprAttrNames,
		Values: in.ExprAttrValues,
	})

	err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, err := txn.GetItem(tableName, in.Key)
		isNewItem := false
		var item *dbstore.Item

		if err != nil {
			if dbstore.IsItemNotFound(err) {
				isNewItem = true
				if conditionChecker != nil {
					if condErr := conditionChecker(nil, true); condErr != nil {
						return condErr
					}
				}
				item = &dbstore.Item{
					TableName:  tableName,
					Key:        in.Key,
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
				for k, v := range in.Key {
					item.Attributes[k] = v
				}
			} else {
				return err
			}
		} else {
			item = existingItem
			oldItemSize = dbstore.CalculateItemSize(item.Attributes)
			if conditionChecker != nil {
				if condErr := conditionChecker(item, false); condErr != nil {
					return condErr
				}
			}

			streamNeedsOld := table.StreamSpecification != nil &&
				table.StreamSpecification.StreamEnabled &&
				(table.StreamSpecification.StreamViewType == dbstore.StreamViewTypeOldImage ||
					table.StreamSpecification.StreamViewType == dbstore.StreamViewTypeNewAndOldImages)
			if in.ReturnValues == "ALL_OLD" || in.ReturnValues == "UPDATED_OLD" || streamNeedsOld {
				oldItem = &dbstore.Item{
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
				for k, v := range item.Attributes {
					oldItem.Attributes[k] = deepCopyAttributeValue(v)
				}
			}

			if err := txn.DeleteIndexEntries(tableName, item); err != nil {
				return err
			}
		}

		if in.UpdateExpr != "" {
			paths := extractUpdatedPaths(in.UpdateExpr, in.ExprAttrNames)
			if err := validateNotKeyAttributes(table, paths); err != nil {
				return err
			}
			var err error
			updatedAttrNames, err = applyUpdateExpressionWithTracking(item.Attributes, in.UpdateExpr, in.ExprAttrNames, in.ExprAttrValues)
			if err != nil {
				if errors.Is(err, ErrTypeMismatch) {
					return ErrInvalidParameter
				}
				return err
			}
		} else if in.AttrUpdates != nil {
			var attrNames []string
			if attrMap, ok := in.AttrUpdates.(map[string]interface{}); ok {
				for k := range attrMap {
					attrNames = append(attrNames, k)
				}
			}
			if err := validateNotKeyAttributes(table, attrNames); err != nil {
				return err
			}
			updatedAttrNames, err = applyAttributeUpdatesWithTracking(item.Attributes, in.AttrUpdates)
			if err != nil {
				return err
			}
		}

		if itemSize := dbstore.CalculateItemSize(item.Attributes); itemSize > dbstore.MaxItemSizeBytes {
			return ErrInvalidParameter
		}

		wasNewItem = isNewItem
		storedItem = &dbstore.Item{
			TableName:  tableName,
			Key:        in.Key,
			Attributes: item.Attributes,
		}
		if err := txn.StoreItemWrite(tableName, in.Key, item.Attributes, nil, !isNewItem, oldItemSize); err != nil {
			return err
		}

		eventName := streamEventForWrite(false, wasNewItem)
		s.captureStreamChangeTxn(txn, store, table, eventName, in.Key, storedItem.Attributes, oldItemAttributes(oldItem))

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.emitChangePropagation(store, region, table, streamEventForWrite(false, wasNewItem), in.Key, storedItem.Attributes, oldItemAttributes(oldItem), s.replicaPutOp(table, in.Key, storedItem.Attributes))

	return &UpdateItemResult{
		StoredItem:       storedItem,
		OldItem:          oldItem,
		UpdatedAttrNames: updatedAttrNames,
		WasNewItem:       wasNewItem,
	}, nil
}
