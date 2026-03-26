package dynamodb

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
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
			continue
		}

		if !store.Tables().Exists(tableName) {
			return nil, ErrTableNotFound
		}

		keys, ok := tr["Keys"].([]interface{})
		if !ok {
			continue
		}

		_ = request.GetBoolParam(tr, "ConsistentRead")

		var tableItems []map[string]interface{}
		var unprocessedKeys []interface{}

		for _, k := range keys {
			key := parseKey(k)
			if key == nil {
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

			projection := parseProjectionExpression(tr)
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
			continue
		}

		table, tableErr := store.Tables().Get(tableName)
		if tableErr != nil {
			unprocessed[tableName] = writes
			continue
		}
		tableCache[tableName] = table

		for _, w := range writes {
			writeReq, ok := w.(map[string]interface{})
			if !ok {
				continue
			}

			if putReq, ok := writeReq["PutRequest"].(map[string]interface{}); ok {
				item := parseItem(putReq["Item"])
				if item == nil {
					continue
				}

				key := s.extractKeyFromItem(table, item)
				if key == nil {
					continue
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
				key := parseKey(delReq["Key"])
				if key == nil {
					continue
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
		opErr := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			switch op.opType {
			case "Put":
				existingItem, err := txn.GetItem(op.tableName, op.key)
				isNewItem := dbstore.IsItemNotFound(err)
				if err != nil && !isNewItem {
					return err
				}
				if existingItem != nil {
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
			return nil
		})

		if opErr != nil {
			var unprocessedItems []interface{}
			if existing, ok := unprocessed[op.tableName].([]interface{}); ok {
				unprocessedItems = existing
			}
			unprocessed[op.tableName] = append(unprocessedItems, op.rawReq)
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
