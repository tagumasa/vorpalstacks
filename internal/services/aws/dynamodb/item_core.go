package dynamodb

import (
	"context"

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

	repKey := key
	repAttrs := item
	repItemSize := calculateItemSize(item)
	s.replicateToGlobalTableReplicas(store, region, table.Name, func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(table.Name, repKey)
			isNewRep := false
			if getErr != nil {
				if dbstore.IsItemNotFound(getErr) {
					isNewRep = true
				} else {
					return getErr
				}
			}
			if !isNewRep && existing != nil {
				if delErr := txn.DeleteIndexEntries(table.Name, existing); delErr != nil {
					return delErr
				}
			}
			if putErr := txn.PutItem(table.Name, repKey, repAttrs); putErr != nil {
				return putErr
			}
			repItem := &dbstore.Item{
				TableName:  table.Name,
				Key:        repKey,
				Attributes: repAttrs,
			}
			if putIdxErr := txn.PutIndexEntries(table.Name, repItem); putIdxErr != nil {
				return putIdxErr
			}
			if isNewRep {
				if upErr := txn.UpdateItemCount(table.Name, 1); upErr != nil {
					return upErr
				}
				if upErr := txn.UpdateTableSize(table.Name, repItemSize); upErr != nil {
					return upErr
				}
			} else {
				oldSize := calculateItemSize(existing.Attributes)
				if upErr := txn.UpdateTableSize(table.Name, repItemSize-oldSize); upErr != nil {
					return upErr
				}
			}
			return nil
		})
	})

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

	repKey := key
	s.replicateToGlobalTableReplicas(store, region, table.Name, func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(table.Name, repKey)
			if getErr != nil {
				if dbstore.IsItemNotFound(getErr) {
					return nil
				}
				return getErr
			}
			if delIdxErr := txn.DeleteIndexEntries(table.Name, existing); delIdxErr != nil {
				return delIdxErr
			}
			if delErr := txn.DeleteItem(table.Name, repKey); delErr != nil {
				return delErr
			}
			if upErr := txn.UpdateItemCount(table.Name, -1); upErr != nil {
				return upErr
			}
			existingSize := calculateItemSize(existing.Attributes)
			if upErr := txn.UpdateTableSize(table.Name, -existingSize); upErr != nil {
				return upErr
			}
			return nil
		})
	})

	return oldItem, nil
}
