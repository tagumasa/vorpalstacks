package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/dynamodb"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// GetItem retrieves a single DynamoDB item by primary key.
func (h *AdminHandler) GetItem(ctx context.Context, req *connect.Request[pb.GetItemInput]) (*connect.Response[pb.GetItemOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}
	if len(req.Msg.GetKey()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Key is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tableName := req.Msg.GetTablename()
	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	if table.Status != dbstore.TableStatusActive {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("table %s is not active", tableName))
	}

	key := protoAVMapToStore(req.Msg.GetKey())
	if !validateKeyValueNotEmpty(key) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key value must not be empty"))
	}

	item, err := store.Items().Get(tableName, key)
	if err != nil {
		if dbstore.IsItemNotFound(err) {
			return connect.NewResponse(&pb.GetItemOutput{}), nil
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.GetItemOutput{
		Item: storeAVMapToProto(item.Attributes),
	}), nil
}

// Scan returns all items in a DynamoDB table with optional pagination.
func (h *AdminHandler) Scan(ctx context.Context, req *connect.Request[pb.ScanInput]) (*connect.Response[pb.ScanOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	limit := 100
	if req.Msg.Limit > 0 {
		limit = int(req.Msg.Limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	marker := ""
	if len(req.Msg.Exclusivestartkey) > 0 {
		marker, _ = h.buildItemMarker(store, req.Msg.GetTablename(), protoAVMapToStore(req.Msg.GetExclusivestartkey()))
	}

	items, nextMarker, err := store.Items().List(req.Msg.GetTablename(), marker, limit)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	pbItems := make([]*pb.ItemListEntry, len(items))
	for i, item := range items {
		pbItems[i] = &pb.ItemListEntry{
			Value: storeAVMapToProto(item.Attributes),
		}
	}

	output := &pb.ScanOutput{
		Items:        pbItems,
		Count:        int32(len(items)),
		Scannedcount: int32(len(items)),
	}
	if nextMarker != "" && len(items) > 0 {
		lastItem := items[len(items)-1]
		output.Lastevaluatedkey = storeAVMapToProto(lastItem.Key)
	}

	return connect.NewResponse(output), nil
}

// PutItem inserts or replaces a DynamoDB item.
func (h *AdminHandler) PutItem(ctx context.Context, req *connect.Request[pb.PutItemInput]) (*connect.Response[pb.PutItemOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}
	if len(req.Msg.GetItem()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Item is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tableName := req.Msg.GetTablename()
	attrs := protoAVMapToStore(req.Msg.GetItem())

	if itemSize := calculateItemSize(attrs); itemSize > maxItemSizeBytes {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("item size %d exceeds maximum allowed size of %d bytes", itemSize, maxItemSizeBytes))
	}

	var storedItem *dbstore.Item

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		table, err := txn.GetTable(tableName)
		if err != nil {
			return err
		}
		if table.Status != dbstore.TableStatusActive {
			return connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("table %s is not active", tableName))
		}

		key := extractKeyFromAttributes(table, attrs)
		if key == nil {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key attributes not found in item"))
		}

		if !validateKeyValueNotEmpty(key) {
			return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key value must not be empty"))
		}

		isNewItem := false
		existingItem, err := txn.GetItem(tableName, key)
		if err != nil {
			if !dbstore.IsItemNotFound(err) {
				return err
			}
			isNewItem = true
		} else if existingItem != nil {
			if err := txn.DeleteIndexEntries(tableName, existingItem); err != nil {
				return err
			}
		}

		if err := txn.PutItem(tableName, key, attrs); err != nil {
			return err
		}

		storedItem = &dbstore.Item{
			TableName:  tableName,
			Key:        key,
			Attributes: attrs,
		}
		if err := txn.PutIndexEntries(tableName, storedItem); err != nil {
			return err
		}

		newItemSize := calculateItemSize(attrs)
		if isNewItem {
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
				return err
			}
		} else if existingItem != nil {
			oldItemSize := calculateItemSize(existingItem.Attributes)
			if newItemSize != oldItemSize {
				if err := txn.UpdateTableSize(tableName, newItemSize-oldItemSize); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			return nil, connectErr
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.PutItemOutput{
		Attributes: storeAVMapToProto(storedItem.Attributes),
	}), nil
}

// DeleteItem removes a DynamoDB item by primary key.
func (h *AdminHandler) DeleteItem(ctx context.Context, req *connect.Request[pb.DeleteItemInput]) (*connect.Response[pb.DeleteItemOutput], error) {
	if req.Msg.GetTablename() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("TableName is required"))
	}
	if len(req.Msg.GetKey()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Key is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tableName := req.Msg.GetTablename()
	key := protoAVMapToStore(req.Msg.GetKey())

	if !validateKeyValueNotEmpty(key) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key value must not be empty"))
	}

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		table, err := txn.GetTable(tableName)
		if err != nil {
			return err
		}
		if table.Status != dbstore.TableStatusActive {
			return connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("table %s is not active", tableName))
		}

		existingItem, err := txn.GetItem(tableName, key)
		if err != nil {
			if dbstore.IsItemNotFound(err) {
				return nil
			}
			return err
		}

		if existingItem != nil {
			if err := txn.DeleteIndexEntries(tableName, existingItem); err != nil {
				return err
			}
		}

		if err := txn.DeleteItem(tableName, key); err != nil {
			return err
		}

		if existingItem != nil {
			oldItemSize := calculateItemSize(existingItem.Attributes)
			if err := txn.UpdateItemCount(tableName, -1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, -oldItemSize); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			return nil, connectErr
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteItemOutput{}), nil
}

func (h *AdminHandler) buildItemMarker(store dbstore.DynamoDBStoreInterface, tableName string, key map[string]*dbstore.AttributeValue) (string, error) {
	table, err := store.Tables().Get(tableName)
	if err != nil {
		return "", err
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
		return tableName + "#", nil
	}

	if skName != "" {
		if key[skName] != nil {
			skValue := avToString(key[skName])
			if skValue != "" {
				return tableName + "#" + pkValue + "#" + skValue, nil
			}
		}
	}

	return tableName + "#" + pkValue, nil
}

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
