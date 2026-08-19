package dynamodb

import (
	"context"
	"errors"

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

// ---------------------------------------------------------------------------
// Core item functions — single validation + persistence path
//
// These functions encapsulate the full item lifecycle including DynamoDB
// Streams capture, Kinesis Data Stream destinations, and global table
// replication. Both the HTTP API handlers and the admin gRPC handler
// delegate to these functions to ensure identical side-effect behaviour.
// ---------------------------------------------------------------------------

// getItemCore retrieves a single item by primary key. It returns
// dbstore.ItemNotFound when the item does not exist.
func (s *DynamoDBService) getItemCore(store dbstore.DynamoDBStoreInterface, tableName string, key map[string]*dbstore.AttributeValue) (*dbstore.Item, error) {
	return store.Items().Get(tableName, key)
}

// ---------------------------------------------------------------------------
// Scan Core — paginated item listing
// ---------------------------------------------------------------------------

// ScanItemsInput is the service-layer DTO for a paginated item scan. It is
// used by the admin Scan handler and any other caller that needs a simple
// paginated listing without filter expressions, parallel segments, or
// projection — features that the data-plane Scan handler applies on top.
type ScanItemsInput struct {
	TableName string
	Limit     int
	Marker    string // empty string means first page
}

// ScanItemsResult is the service-layer result of a paginated item scan.
type ScanItemsResult struct {
	Items      []*dbstore.Item
	NextMarker string // empty string means no more pages
}

// scanItemsCore returns a single page of items for the specified table.
// It applies the standard default-and-cap limit logic and delegates to
// store.Items().List. The admin Scan handler calls this method so that
// no admin code path touches the store layer directly.
//
// The data-plane Scan handler (item_query.go) uses a richer code path with
// parallel-segment, filter-expression, and projection support; this Core
// method covers the simpler paginated-listing case needed by the admin
// console.
func (s *DynamoDBService) scanItemsCore(store dbstore.DynamoDBStoreInterface, in ScanItemsInput) (*ScanItemsResult, error) {
	lim := in.Limit
	if lim <= 0 {
		lim = dataPlaneQueryDefaultLimit
	}
	if lim > dataPlaneQueryMaxLimit {
		lim = dataPlaneQueryMaxLimit
	}
	items, nextMarker, err := store.Items().List(in.TableName, in.Marker, lim)
	if err != nil {
		return nil, err
	}
	return &ScanItemsResult{
		Items:      items,
		NextMarker: nextMarker,
	}, nil
}

// putItemCore creates or replaces an item, applying all side effects:
// DynamoDB Streams capture (within the transaction), Kinesis Data Stream
// destinations (after commit), and global table replication.
//
// The optional conditionChecker is invoked inside the transaction after
// loading the existing item. Pass nil to skip condition evaluation (admin
// console path).
//
// Returns the stored item and the previous item (nil if new).
func (s *DynamoDBService) putItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	key, item map[string]*dbstore.AttributeValue,
	conditionChecker ConditionChecker,
) (storedItem, oldItem *dbstore.Item, err error) {
	var isNew bool

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, getErr := txn.GetItem(table.Name, key)
		if getErr != nil {
			if dbstore.IsItemNotFound(getErr) {
				isNew = true
			} else {
				return getErr
			}
		} else {
			oldItem = existingItem
			if oldItem != nil {
				if delErr := txn.DeleteIndexEntries(table.Name, oldItem); delErr != nil {
					return delErr
				}
			}
		}

		if conditionChecker != nil {
			if condErr := conditionChecker(existingItem, isNew); condErr != nil {
				return condErr
			}
		}

		if putErr := txn.PutItem(table.Name, key, item); putErr != nil {
			return putErr
		}

		storedItem = &dbstore.Item{
			TableName:  table.Name,
			Key:        key,
			Attributes: item,
		}
		if putIdxErr := txn.PutIndexEntries(table.Name, storedItem); putIdxErr != nil {
			return putIdxErr
		}

		newItemSize := calculateItemSize(item)
		if isNew {
			if upErr := txn.UpdateItemCount(table.Name, 1); upErr != nil {
				return upErr
			}
			if upErr := txn.UpdateTableSize(table.Name, newItemSize); upErr != nil {
				return upErr
			}
		} else if oldItem != nil {
			oldItemSize := calculateItemSize(oldItem.Attributes)
			if newItemSize != oldItemSize {
				if upErr := txn.UpdateTableSize(table.Name, newItemSize-oldItemSize); upErr != nil {
					return upErr
				}
			}
		}

		eventName := dbstore.StreamEventModify
		if isNew {
			eventName = dbstore.StreamEventInsert
		}
		s.captureStreamChangeTxn(txn, store, table, eventName, key, storedItem.Attributes, oldItemAttributes(oldItem))

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	eventName := dbstore.StreamEventModify
	if isNew {
		eventName = dbstore.StreamEventInsert
	}
	s.sendToKinesisDestinations(table, eventName, key, storedItem.Attributes, oldItemAttributes(oldItem))

	s.replicateToGlobalTableReplicas(store, region, table.Name, replicaPutOp(table, key, item))

	return storedItem, oldItem, nil
}

// deleteItemCore removes an item by primary key, applying all side effects:
// DynamoDB Streams capture (within the transaction), Kinesis Data Stream
// destinations (after commit), and global table replication.
//
// The optional conditionChecker is invoked inside the transaction after
// loading the existing item. Pass nil to skip condition evaluation (admin
// console path).
//
// Returns the previous item (nil if it did not exist).
func (s *DynamoDBService) deleteItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	key map[string]*dbstore.AttributeValue,
	conditionChecker ConditionChecker,
) (oldItem *dbstore.Item, err error) {
	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
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
			if delIdxErr := txn.DeleteIndexEntries(table.Name, oldItem); delIdxErr != nil {
				return delIdxErr
			}
			if delErr := txn.DeleteItem(table.Name, key); delErr != nil {
				return delErr
			}
			if upErr := txn.UpdateItemCount(table.Name, -1); upErr != nil {
				return upErr
			}
			oldItemSize := calculateItemSize(oldItem.Attributes)
			if upErr := txn.UpdateTableSize(table.Name, -oldItemSize); upErr != nil {
				return upErr
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

	s.replicateToGlobalTableReplicas(store, region, table.Name, replicaDeleteOp(table, key))

	return oldItem, nil
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
// HTTP API and admin gRPC handler. It evaluates conditions, applies update
// expressions, persists within a transaction, and fires all side effects
// (Streams, Kinesis, global table replication).
func (s *DynamoDBService) updateItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	in UpdateItemInput,
) (*UpdateItemResult, error) {
	tableName := table.Name

	var oldItem *dbstore.Item
	var storedItem *dbstore.Item
	var updatedAttrNames []string
	var oldItemSize int64
	var wasNewItem bool

	err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, err := txn.GetItem(tableName, in.Key)
		isNewItem := false
		var item *dbstore.Item

		if err != nil {
			if dbstore.IsItemNotFound(err) {
				isNewItem = true
				if in.ConditionExpr != "" {
					syntheticItem := &dbstore.Item{
						TableName:  tableName,
						Key:        in.Key,
						Attributes: make(map[string]*dbstore.AttributeValue),
					}
					conditionMet, evalErr := evaluateConditionExpression(syntheticItem, in.ConditionExpr, in.ExprAttrNames, in.ExprAttrValues)
					if evalErr != nil {
						return evalErr
					}
					if !conditionMet {
						return ErrConditionalCheckFailed
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
			oldItemSize = calculateItemSize(item.Attributes)
			if in.ConditionExpr != "" {
				conditionMet, err := evaluateConditionExpression(item, in.ConditionExpr, in.ExprAttrNames, in.ExprAttrValues)
				if err != nil {
					return err
				}
				if !conditionMet {
					return ErrConditionalCheckFailed
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

		if itemSize := calculateItemSize(item.Attributes); itemSize > maxItemSizeBytes {
			return ErrInvalidParameter
		}

		if err := txn.PutItem(tableName, in.Key, item.Attributes); err != nil {
			return err
		}

		storedItem = &dbstore.Item{
			TableName:  tableName,
			Key:        in.Key,
			Attributes: item.Attributes,
		}
		if err := txn.PutIndexEntries(tableName, storedItem); err != nil {
			return err
		}

		newItemSize := calculateItemSize(item.Attributes)
		wasNewItem = isNewItem
		if isNewItem {
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
				return err
			}
		} else {
			if newItemSize != oldItemSize {
				if err := txn.UpdateTableSize(tableName, newItemSize-oldItemSize); err != nil {
					return err
				}
			}
		}

		eventName := dbstore.StreamEventModify
		if wasNewItem {
			eventName = dbstore.StreamEventInsert
		}
		s.captureStreamChangeTxn(txn, store, table, eventName, in.Key, storedItem.Attributes, oldItemAttributes(oldItem))

		return nil
	})
	if err != nil {
		return nil, err
	}

	{
		eventName := dbstore.StreamEventModify
		if wasNewItem {
			eventName = dbstore.StreamEventInsert
		}
		s.sendToKinesisDestinations(table, eventName, in.Key, storedItem.Attributes, oldItemAttributes(oldItem))
	}

	repKey := in.Key
	repAttrs := storedItem.Attributes
	repItemSize := calculateItemSize(repAttrs)
	s.replicateToGlobalTableReplicas(store, region, tableName, func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(tableName, repKey)
			isNew := false
			if getErr != nil {
				if dbstore.IsItemNotFound(getErr) {
					isNew = true
				} else {
					return getErr
				}
			}
			if !isNew && existing != nil {
				if err := txn.DeleteIndexEntries(tableName, existing); err != nil {
					return err
				}
			}
			if err := txn.PutItem(tableName, repKey, repAttrs); err != nil {
				return err
			}
			repItem := &dbstore.Item{
				TableName:  tableName,
				Key:        repKey,
				Attributes: repAttrs,
			}
			if err := txn.PutIndexEntries(tableName, repItem); err != nil {
				return err
			}
			if isNew {
				if err := txn.UpdateItemCount(tableName, 1); err != nil {
					return err
				}
				if err := txn.UpdateTableSize(tableName, repItemSize); err != nil {
					return err
				}
			} else {
				oldSize := calculateItemSize(existing.Attributes)
				if err := txn.UpdateTableSize(tableName, repItemSize-oldSize); err != nil {
					return err
				}
			}
			return nil
		})
	})

	return &UpdateItemResult{
		StoredItem:       storedItem,
		OldItem:          oldItem,
		UpdatedAttrNames: updatedAttrNames,
		WasNewItem:       wasNewItem,
	}, nil
}
