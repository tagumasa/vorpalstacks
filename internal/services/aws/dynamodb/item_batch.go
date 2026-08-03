package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// BatchGetItem retrieves multiple items from one or more tables in a single request.
func (s *DynamoDBService) BatchGetItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	requestItems, ok := req.Parameters["RequestItems"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	totalKeys := 0
	for _, tableReq := range requestItems {
		if tr, ok := tableReq.(map[string]interface{}); ok {
			if keys, ok := tr["Keys"].([]interface{}); ok {
				totalKeys += len(keys)
			}
		}
	}
	if totalKeys > 100 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	responses := make(map[string]interface{})
	unprocessed := make(map[string]interface{})

	for tableName, tableReq := range requestItems {
		tr, ok := tableReq.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		if !store.Tables().Exists(tableName) {
			return nil, ErrTableNotFound
		}

		keys, ok := tr["Keys"].([]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		// Single-instance Pebble provides strong consistency by default.
		// ConsistentRead is accepted for API compatibility.
		_ = request.GetBoolParam(tr, "ConsistentRead")

		var tableItems []map[string]interface{}
		var unprocessedKeys []interface{}

		for _, k := range keys {
			key, keyErr := parseKey(k)
			if keyErr != nil || key == nil {
				unprocessedKeys = append(unprocessedKeys, k)
				continue
			}

			item, err := store.Items().Get(tableName, key)
			if err != nil {
				if dbstore.IsItemNotFound(err) {
					continue
				}
				unprocessedKeys = append(unprocessedKeys, k)
				continue
			}

			projection, projErr := parseProjectionExpression(tr)
			if projErr != nil {
				unprocessedKeys = append(unprocessedKeys, k)
				continue
			}
			if projection != nil {
				item.Attributes = applyProjection(item.Attributes, projection)
			}

			tableItems = append(tableItems, buildItemResponse(item.Attributes))
		}

		if len(tableItems) > 0 {
			responses[tableName] = tableItems
		}
		if len(unprocessedKeys) > 0 {
			unprocessed[tableName] = map[string]interface{}{"Keys": unprocessedKeys}
		}
	}

	resp := map[string]interface{}{
		"Responses":       responses,
		"UnprocessedKeys": unprocessed,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		var consumedCapacity []interface{}
		for tableName := range responses {
			consumedCapacity = append(consumedCapacity, buildConsumedCapacityResponse(tableName, 0.5))
		}
		if len(consumedCapacity) > 0 {
			resp["ConsumedCapacity"] = consumedCapacity
		}
	}

	return resp, nil
}

// BatchWriteItem inserts, updates, or deletes multiple items across one or more tables.
func (s *DynamoDBService) BatchWriteItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	requestItems, ok := req.Parameters["RequestItems"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	totalItems := 0
	for _, tableReq := range requestItems {
		if writes, ok := tableReq.([]interface{}); ok {
			totalItems += len(writes)
		}
	}
	if totalItems > 25 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	unprocessed := make(map[string]interface{})

	type writeOp struct {
		tableName string
		opType    string
		key       map[string]*dbstore.AttributeValue
		item      map[string]*dbstore.AttributeValue
		rawReq    map[string]interface{}
	}

	var allWrites []writeOp
	tableCache := make(map[string]*dbstore.Table)

	for tableName, tableReq := range requestItems {
		writes, ok := tableReq.([]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		table, tableErr := store.Tables().Get(tableName)
		if tableErr != nil {
			if dbstore.IsTableNotFound(tableErr) {
				return nil, ErrTableNotFound
			}
			unprocessed[tableName] = writes
			continue
		}
		tableCache[tableName] = table

		for _, w := range writes {
			writeReq, ok := w.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameter
			}

			if putReq, ok := writeReq["PutRequest"].(map[string]interface{}); ok {
				item, itemErr := parseItem(putReq["Item"])
				if itemErr != nil || item == nil {
					return nil, ErrInvalidParameter
				}

				// Item size must not exceed 400 KB (same as PutItem).
				if calculateItemSize(item) > maxItemSizeBytes {
					return nil, ErrInvalidParameter
				}

				key := s.extractKeyFromItem(table, item)
				if key == nil {
					return nil, ErrInvalidParameter
				}

				// Key attribute values must not be empty.
				if !validateKeyValueNotEmpty(key) {
					return nil, ErrInvalidParameter
				}

				allWrites = append(allWrites, writeOp{
					tableName: tableName,
					opType:    "Put",
					key:       key,
					item:      item,
					rawReq:    writeReq,
				})
			}

			if delReq, ok := writeReq["DeleteRequest"].(map[string]interface{}); ok {
				key, keyErr := parseKey(delReq["Key"])
				if keyErr != nil || key == nil {
					return nil, ErrInvalidParameter
				}

				// Key attribute values must not be empty.
				if !validateKeyValueNotEmpty(key) {
					return nil, ErrInvalidParameter
				}

				allWrites = append(allWrites, writeOp{
					tableName: tableName,
					opType:    "Delete",
					key:       key,
					rawReq:    writeReq,
				})
			}
		}
	}

	for _, op := range allWrites {
		var batchIsNew bool
		var batchOldAttrs map[string]*dbstore.AttributeValue

		opErr := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			switch op.opType {
			case "Put":
				existingItem, err := txn.GetItem(op.tableName, op.key)
				isNewItem := dbstore.IsItemNotFound(err)
				batchIsNew = isNewItem
				if err != nil && !isNewItem {
					return err
				}
				if existingItem != nil {
					batchOldAttrs = existingItem.Attributes
					if err := txn.DeleteIndexEntries(op.tableName, existingItem); err != nil {
						return err
					}
				}
				if err := txn.PutItem(op.tableName, op.key, op.item); err != nil {
					return err
				}
				newItem := &dbstore.Item{
					TableName:  op.tableName,
					Key:        op.key,
					Attributes: op.item,
				}
				if err := txn.PutIndexEntries(op.tableName, newItem); err != nil {
					return err
				}
				newItemSize := calculateItemSize(op.item)
				if isNewItem {
					if err := txn.UpdateItemCount(op.tableName, 1); err != nil {
						return err
					}
					if err := txn.UpdateTableSize(op.tableName, newItemSize); err != nil {
						return err
					}
				} else if existingItem != nil {
					oldItemSize := calculateItemSize(existingItem.Attributes)
					if err := txn.UpdateTableSize(op.tableName, newItemSize-oldItemSize); err != nil {
						return err
					}
				}

			case "Delete":
				existingItem, err := txn.GetItem(op.tableName, op.key)
				if dbstore.IsItemNotFound(err) {
					return nil
				}
				if err != nil {
					return err
				}
				if existingItem != nil {
					batchOldAttrs = existingItem.Attributes
					if err := txn.DeleteIndexEntries(op.tableName, existingItem); err != nil {
						return err
					}
				}
				if err := txn.DeleteItem(op.tableName, op.key); err != nil {
					return err
				}
				if err := txn.UpdateItemCount(op.tableName, -1); err != nil {
					return err
				}
				if existingItem != nil {
					oldItemSize := calculateItemSize(existingItem.Attributes)
					if err := txn.UpdateTableSize(op.tableName, -oldItemSize); err != nil {
						return err
					}
				}
			}

			table := tableCache[op.tableName]
			if op.opType == "Put" {
				eventName := dbstore.StreamEventModify
				if batchIsNew {
					eventName = dbstore.StreamEventInsert
				}
				s.captureStreamChangeTxn(txn, store, table, eventName, op.key, op.item, batchOldAttrs)
			} else if op.opType == "Delete" && batchOldAttrs != nil {
				s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventRemove, op.key, nil, batchOldAttrs)
			}

			return nil
		})

		if opErr != nil {
			var unprocessedItems []interface{}
			if existing, ok := unprocessed[op.tableName].([]interface{}); ok {
				unprocessedItems = existing
			}
			unprocessed[op.tableName] = append(unprocessedItems, op.rawReq)
		} else {
			table := tableCache[op.tableName]
			if op.opType == "Put" {
				eventName := dbstore.StreamEventModify
				if batchIsNew {
					eventName = dbstore.StreamEventInsert
				}
				s.sendToKinesisDestinations(table, eventName, op.key, op.item, batchOldAttrs)
			} else if op.opType == "Delete" && batchOldAttrs != nil {
				s.sendToKinesisDestinations(table, dbstore.StreamEventRemove, op.key, nil, batchOldAttrs)
			}
		}
	}

	resp := map[string]interface{}{
		"UnprocessedItems": unprocessed,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		var consumedCapacity []interface{}
		for tableName := range requestItems {
			consumedCapacity = append(consumedCapacity, buildConsumedCapacityResponse(tableName, 1.0))
		}
		if len(consumedCapacity) > 0 {
			resp["ConsumedCapacity"] = consumedCapacity
		}
	}

	return resp, nil
}
