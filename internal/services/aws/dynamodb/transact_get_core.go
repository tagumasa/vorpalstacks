package dynamodb

import (
	"context"
	"fmt"
	"net/http"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// TransactGetItems / TransactWriteItems Cores — single validation and
// persistence paths of the transactional data plane
// ---------------------------------------------------------------------------
//
// The two Cores this banner named live one per file: this one holds the
// TransactGetItems half, and the write half (operation parsing, the
// two-phase executor, the idempotency window, the capacity report) lives
// in transact_write_core.go. The transactional read-unit helper below is
// shared by both planes — the write side's condition-check reads and
// client-token replay report both charge through it.

// transactGetItemsInput carries the already-typed TransactItems member plus
// the raw wire parameters (consumed for the ReturnConsumedCapacity reporting).
type transactGetItemsInput struct {
	TransactItems []interface{}
	Parameters    map[string]interface{}
}

// transactGetItemsCore resolves every referenced table, validates the key
// types, and reads all items in one snapshot view.
func (s *DynamoDBService) transactGetItemsCore(ctx context.Context, reqCtx *request.RequestContext, in transactGetItemsInput) (map[string]interface{}, error) {
	if in.TransactItems == nil {
		return nil, ErrInvalidParameter
	}

	transactItems := in.TransactItems

	if len(transactItems) > transactMaxItems {
		return nil, ErrInvalidParameter
	}

	if len(transactItems) == 0 {
		return nil, ErrInvalidParameter
	}

	// Response-shaping enums are request validation: an unknown value is
	// rejected before anything executes.
	returnConsumedCapacity, err := getReturnConsumedCapacity(in.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, fmt.Errorf("transact get items: get store: %w", err)
	}

	type getItem struct {
		tableName  string
		key        map[string]*dbstore.AttributeValue
		projection [][]docPathPart
	}

	var getItems []getItem
	seenKeys := make(map[string]bool)
	for _, item := range transactItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		getMap, ok := itemMap["Get"].(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		tableName := request.GetStringParam(getMap, "TableName")
		if !validateResourceName(tableName) {
			return nil, ErrInvalidParameter
		}

		key, keyErr := parseKey(getMap["Key"])
		if keyErr != nil || key == nil {
			return nil, ErrInvalidParameter
		}

		// The Get actions target distinct items ("up to 100 distinct
		// items"; "you can't target the same item with multiple
		// operations within the same transaction"): the same item
		// requested twice is request validation on this plane too.
		keyStr := buildKeyString(tableName, key)
		if seenKeys[keyStr] {
			return nil, errTransactionDuplicateItem()
		}
		seenKeys[keyStr] = true

		projection, projErr := parseProjectionExpression(getMap)
		if projErr != nil {
			return nil, ErrInvalidParameter
		}

		getItems = append(getItems, getItem{
			tableName:  tableName,
			key:        key,
			projection: projection,
		})
	}

	// Resolve every referenced table and validate its key types before
	// the snapshot view: a missing table is a ResourceNotFoundException,
	// not an empty Item slot, and key attribute types must match the
	// table's attribute definitions.
	for _, gi := range getItems {
		table, tableErr := store.Tables().Get(gi.tableName)
		if tableErr != nil || table == nil {
			return nil, ErrTableNotFound
		}
		if keyErr := validateKeyTypes(table, gi.key); keyErr != nil {
			return nil, keyErr
		}
		// The Key must carry exactly the key schema's attributes, matching
		// the single-op read plane.
		if keyErr := validateKeySchemaMembership(table, gi.key); keyErr != nil {
			return nil, keyErr
		}
	}

	var responses []map[string]interface{}
	foundKeys := make(map[string][]map[string]*dbstore.AttributeValue)
	// Every targeted item consumes read capacity of its own, charged to its
	// table: the ConsumedCapacity response carries one entry per table
	// addressed, reporting that table's units.
	tableReadUnits := make(map[string]float64)
	// The entries answer in first-appearance order: TransactItems is an
	// ordered list, so the response follows the request rather than the
	// map's random iteration order.
	var tableOrder []string
	addReadUnits := func(tableName string, units float64) {
		if _, seen := tableReadUnits[tableName]; !seen {
			tableOrder = append(tableOrder, tableName)
		}
		tableReadUnits[tableName] += units
	}
	// The aggregate of the retrieved items is measured on the full items as
	// read, before any projection narrows the returned attributes — the
	// read plane's aggregate is over what the transaction read, and a
	// request past the limit is rejected whole.
	var readAggregate int64

	// The read transaction observes the same serialisation: an item that is
	// part of an ongoing TransactWriteItems or single-item write rejects
	// the whole read with a TransactionConflict cancellation reason at that
	// item's slot. The check marks nothing — an ongoing read never rejects
	// another request.
	readKeys := make([]string, len(getItems))
	readKeySlot := make(map[string]int, len(getItems))
	for i, gi := range getItems {
		k := itemLockKey(reqCtx.Region, gi.tableName, gi.key)
		readKeys[i] = k
		readKeySlot[k] = i
	}
	if conflictKey, locked := firstLockedItem(readKeys); locked {
		conflictReasons := make([]CancellationReason, len(getItems))
		for i := range conflictReasons {
			conflictReasons[i] = CancellationReason{Code: "None"}
		}
		conflictReasons[readKeySlot[conflictKey]] = CancellationReason{
			Code:    "TransactionConflict",
			Message: transactionConflictReasonMessage,
		}
		return nil, NewTransactionCanceledError("Transaction canceled", conflictReasons)
	}

	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		for _, gi := range getItems {
			dbItem, err := txn.GetItem(gi.tableName, gi.key)
			if err != nil {
				if !dbstore.IsItemNotFound(err) {
					return fmt.Errorf("transact get item on %s: %w", gi.tableName, err)
				}
				addReadUnits(gi.tableName, transactItemReadUnits(0))
				// The missing item answers an empty response slot: the
				// protocol omits unset members rather than carrying an
				// explicit null.
				responses = append(responses, map[string]interface{}{})
				continue
			}

			// The charge follows the item's full size as read, before any
			// projection narrows the returned attributes.
			itemSize := dbstore.CalculateItemSize(dbItem.Attributes)
			addReadUnits(gi.tableName, transactItemReadUnits(itemSize))
			readAggregate += itemSize

			attrs := dbItem.Attributes
			if gi.projection != nil {
				attrs = applyProjection(attrs, gi.projection)
			}

			responses = append(responses, map[string]interface{}{
				"Item": buildItemResponse(attrs),
			})
			foundKeys[gi.tableName] = append(foundKeys[gi.tableName], gi.key)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transact get items: snapshot view: %w", err)
	}
	if readAggregate > dbstore.MaxTransactionAggregateBytes {
		return nil, errTransactGetAggregateExceeded()
	}

	// Every item the transactional read actually returned counts as one
	// read event per tracked key layout. The counters cannot update inside
	// the read-only view, so they are applied after it succeeds.
	for tableName, keys := range foundKeys {
		s.recordContributorReads(ctx, store, tableName, keys)
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		// TOTAL reports only the aggregates; INDEXES adds the per-table
		// detail.
		consumedCapacities := make([]map[string]interface{}, 0, len(tableOrder))
		for _, tableName := range tableOrder {
			consumedCapacities = append(consumedCapacities, buildReadConsumedCapacityResponse(tableName, tableReadUnits[tableName], returnConsumedCapacity == "INDEXES"))
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
}

// transactItemReadUnits returns the read capacity a single targeted item
// consumes in a transactional read: DynamoDB performs two underlying reads
// per item — one to prepare the transaction, one to commit it — over the
// item's size rounded up to 4 KB multiples, and consumes that capacity even
// when the item turns out to be absent, because the two reads still happen.
// Transactional reads are strongly consistent.
func transactItemReadUnits(itemSizeBytes int64) float64 {
	return itemReadUnits(itemSizeBytes, true) * 2
}

// errTransactGetAggregateExceeded is the read plane's aggregate-size
// rejection, verbatim from the TransactGetItems documentation: "The
// aggregate size of the items in the transaction exceeded 4 MB."
func errTransactGetAggregateExceeded() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"The aggregate size of the items in the transaction exceeded 4 MB", http.StatusBadRequest)
}
