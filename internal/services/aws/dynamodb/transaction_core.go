package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// TransactGetItems / TransactWriteItems Cores — single validation and
// persistence paths of the transactional data plane
// ---------------------------------------------------------------------------

// transactGetItemsInput carries the already-typed TransactItems member plus
// the raw wire parameters (consumed for the ReturnConsumedCapacity reporting).
type transactGetItemsInput struct {
	TransactItems []interface{}
	Parameters    map[string]interface{}
}

// transactGetItemsCore resolves every referenced table, validates the key
// types, and reads all items in one snapshot view.
func (s *DynamoDBService) transactGetItemsCore(ctx context.Context, reqCtx *request.RequestContext, in transactGetItemsInput) (map[string]interface{}, error) {
	transactItems := in.TransactItems

	if len(transactItems) > transactMaxItems {
		return nil, ErrInvalidParameter
	}

	if len(transactItems) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, fmt.Errorf("transact get items: get store: %w", err)
	}

	type getItem struct {
		tableName  string
		key        map[string]*dbstore.AttributeValue
		projection []string
	}

	var getItems []getItem
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
	}

	var responses []map[string]interface{}
	foundKeys := make(map[string][]map[string]*dbstore.AttributeValue)
	// Every targeted item consumes read capacity of its own, charged to its
	// table: the ConsumedCapacity response carries one entry per table
	// addressed, reporting that table's units.
	tableReadUnits := make(map[string]float64)

	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		for _, gi := range getItems {
			dbItem, err := txn.GetItem(gi.tableName, gi.key)
			if err != nil {
				if !dbstore.IsItemNotFound(err) {
					return fmt.Errorf("transact get item on %s: %w", gi.tableName, err)
				}
				tableReadUnits[gi.tableName] += transactItemReadUnits(0)
				responses = append(responses, map[string]interface{}{"Item": nil})
				continue
			}

			// The charge follows the item's full size as read, before any
			// projection narrows the returned attributes.
			tableReadUnits[gi.tableName] += transactItemReadUnits(dbstore.CalculateItemSize(dbItem.Attributes))

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

	// Every item the transactional read actually returned counts as one
	// read event per tracked key layout. The counters cannot update inside
	// the read-only view, so they are applied after it succeeds.
	for tableName, keys := range foundKeys {
		s.recordContributorReads(ctx, store, tableName, keys)
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(in.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		consumedCapacities := make([]map[string]interface{}, 0, len(tableReadUnits))
		for tableName, units := range tableReadUnits {
			consumedCapacities = append(consumedCapacities, buildReadConsumedCapacityResponse(tableName, units))
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
func transactItemReadUnits(itemSizeBytes int64) float64 {
	units := (itemSizeBytes + dbstore.ReadCapacityUnitBytes - 1) / dbstore.ReadCapacityUnitBytes
	if units < 1 {
		units = 1
	}
	return float64(units * 2)
}

// transactWriteItemsInput carries the already-typed TransactItems member plus
// the raw wire parameters (consumed for the idempotency payload hash and the
// ReturnConsumedCapacity reporting).
type transactWriteItemsInput struct {
	TransactItems []interface{}
	Parameters    map[string]interface{}
}

// transactWriteItemsCore applies the ClientRequestToken idempotency window,
// parses and validates every write operation, executes them in one two-phase
// transaction, and fires the post-commit Kinesis and global-table replication
// side effects.
func (s *DynamoDBService) transactWriteItemsCore(ctx context.Context, reqCtx *request.RequestContext, in transactWriteItemsInput) (map[string]interface{}, error) {
	transactItems := in.TransactItems

	if len(transactItems) > transactMaxItems {
		return nil, ErrInvalidParameter
	}

	if len(transactItems) == 0 {
		return nil, ErrInvalidParameter
	}

	clientRequestToken := request.GetStringParam(in.Parameters, "ClientRequestToken")
	if !validateClientRequestToken(clientRequestToken) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, fmt.Errorf("transact write items: get store: %w", err)
	}

	// A ClientRequestToken makes the call idempotent for a ten-minute
	// window: a retry with the same token and the same payload replays the
	// recorded outcome without re-executing, while the same token with a
	// different payload is rejected.
	requestHash := ""
	claimedToken := false
	if clientRequestToken != "" {
		requestHash = hashTransactWriteRequest(in.Parameters)

		// The claim section is serialised per token so concurrent retries
		// cannot both execute: the loser observes the in-progress record
		// and fails with the documented in-progress error.
		unlock := s.lockClientRequestToken(clientRequestToken)
		recordedHash, state, recordedReadUnits, found, lookupErr := store.Idempotency().Lookup(clientRequestToken)
		if lookupErr != nil {
			unlock()
			return nil, lookupErr
		}
		if found {
			unlock()
			if recordedHash != requestHash {
				return nil, ErrIdempotentParameterMismatch
			}
			if state != dbstore.IdempotencyStateCompleted {
				return nil, ErrTransactionInProgress
			}
			// Replay: parse for the response shape only; no execution and
			// no stream or replication side effects. The capacity report
			// replays the read units recorded at completion.
			replayReasons := make([]CancellationReason, len(transactItems))
			for i := range replayReasons {
				replayReasons[i] = CancellationReason{Code: "None"}
			}
			replayOps, replayErr := parseTransactWriteItems(s, store, transactItems, replayReasons)
			if replayErr != nil {
				return nil, replayErr
			}
			return buildTransactWriteResponse(in.Parameters, replayOps, true, recordedReadUnits), nil
		}
		if claimErr := store.Idempotency().Record(clientRequestToken, requestHash,
			dbstore.IdempotencyStateInProgress,
			time.Now().Add(idempotencyWindowMinutes*time.Minute), nil); claimErr != nil {
			unlock()
			return nil, claimErr
		}
		claimedToken = true
		unlock()
	}

	// releaseClaimedToken drops the in-progress record after a failed
	// execution so the client can retry the token.
	releaseClaimedToken := func() {
		if !claimedToken {
			return
		}
		if delErr := store.Idempotency().Delete(clientRequestToken); delErr != nil {
			logs.Error("Failed to release idempotency token claim", logs.Err(delErr))
		}
	}

	cancellationReasons := make([]CancellationReason, len(transactItems))
	for i := range cancellationReasons {
		cancellationReasons[i] = CancellationReason{Code: "None"}
	}

	operations, err := parseTransactWriteItems(s, store, transactItems, cancellationReasons)
	if err != nil {
		releaseClaimedToken()
		return nil, err
	}

	if err := executeTransactWriteItems(ctx, s, store, operations, cancellationReasons); err != nil {
		releaseClaimedToken()
		return nil, err
	}

	if claimedToken {
		if recordErr := store.Idempotency().Record(clientRequestToken, requestHash,
			dbstore.IdempotencyStateCompleted,
			time.Now().Add(idempotencyWindowMinutes*time.Minute),
			transactReplayReadUnits(operations)); recordErr != nil {
			return nil, recordErr
		}
	}

	// Dispatch change records to Kinesis destinations asynchronously
	// after the transaction has committed.
	tableCache := make(map[string]*dbstore.Table)
	var metricsWrites []itemCollectionWriteRef
	for _, op := range operations {
		if op.opType == "ConditionCheck" {
			continue
		}
		table, ok := tableCache[op.tableName]
		if !ok {
			var tblErr error
			table, tblErr = store.Tables().Get(op.tableName)
			if tblErr != nil || table == nil {
				continue
			}
			tableCache[op.tableName] = table
		}
		metricsWrites = append(metricsWrites, itemCollectionWriteRef{tableName: op.tableName, table: table, key: op.key})
		eventName := streamEventForWrite(op.opType == "Delete", op.streamWasNew)

		// Committed transaction writes replicate to global table replica
		// regions just like single-item writes do; Update replicates its
		// committed post-image as a put.
		var replicaOp func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error
		switch op.opType {
		case "Put":
			if op.itemData != nil {
				replicaOp = s.replicaPutOp(table, op.key, op.itemData)
			}
		case "Update":
			if op.streamNewImage != nil {
				replicaOp = s.replicaPutOp(table, op.key, op.streamNewImage)
			}
		case "Delete":
			replicaOp = s.replicaDeleteOp(table, op.key)
		}
		s.emitChangePropagation(store, reqCtx.GetRegion(), table, eventName, op.key, op.streamNewImage, op.streamOldImage, replicaOp)
	}

	resp := buildTransactWriteResponse(in.Parameters, operations, false, nil)
	// ReturnItemCollectionMetrics=SIZE asks for one entry per item
	// collection the committed transaction wrote; the idempotent replay
	// path answers with the response shape only, without re-deriving it.
	if request.GetStringParam(in.Parameters, "ReturnItemCollectionMetrics") == "SIZE" {
		if metrics := buildItemCollectionMetricsPerTable(metricsWrites); metrics != nil {
			resp["ItemCollectionMetrics"] = metrics
		}
	}
	return resp, nil
}

// writeOperation represents a single write operation within a transaction.
type writeOperation struct {
	idx                          int
	opType                       string
	tableName                    string
	key                          map[string]*dbstore.AttributeValue
	itemData                     map[string]*dbstore.AttributeValue
	updateReq                    map[string]interface{}
	conditionExpr                string
	exprAttrNames                map[string]string
	exprAttrValues               map[string]*dbstore.AttributeValue
	returnValuesOnConditionCheck string
	// streamOldImage and streamNewImage are populated during execution
	// so that stream records can be captured after the transaction commits.
	streamOldImage map[string]*dbstore.AttributeValue
	streamNewImage map[string]*dbstore.AttributeValue
	streamWasNew   bool
	// itemSize is the size of the targeted item as it stands before the
	// operation, captured while validating the transaction — the charge
	// basis for deletes (the deleted item's size) and condition checks.
	itemSize int64
}

func parseTransactWriteItems(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, transactItems []interface{}, cancellationReasons []CancellationReason) ([]writeOperation, error) {
	usedWriteKeys := make(map[string]bool)
	usedConditionKeys := make(map[string]bool)
	var operations []writeOperation

	for idx, item := range transactItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		op, err := parseWriteOperation(s, store, idx, itemMap, usedWriteKeys, usedConditionKeys, cancellationReasons)
		if err != nil {
			return nil, err
		}
		if op != nil {
			operations = append(operations, *op)
		}
	}

	return operations, nil
}

func parseWriteOperation(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, idx int, itemMap map[string]interface{}, usedWriteKeys map[string]bool, usedConditionKeys map[string]bool, cancellationReasons []CancellationReason) (*writeOperation, error) {
	for _, opType := range []string{"Put", "Update", "Delete", "ConditionCheck"} {
		opMap, ok := itemMap[opType].(map[string]interface{})
		if !ok {
			continue
		}

		// Request-shape failures (malformed members, unknown or missing
		// tables) are request-level validation errors answered with
		// ValidationException or ResourceNotFoundException. Only semantic
		// per-item outcomes — the same item targeted by more than one
		// action, condition failures during execution — cancel the
		// transaction and are reported through CancellationReasons.
		tableName := request.GetStringParam(opMap, "TableName")
		if !validateResourceName(tableName) {
			return nil, ErrInvalidParameter
		}

		// Every table a write transaction references must be ACTIVE: a table
		// mid-restore (CREATING) must be neither mutated nor condition-checked
		// against its partial data, so the whole request is rejected like a
		// single-item write against the same table.
		opTable, opTableErr := store.Tables().Get(tableName)
		if opTableErr != nil {
			return nil, ErrTableNotFound
		}
		if opTable.Status != dbstore.TableStatusActive {
			return nil, ErrTableNotActive
		}

		key, err := extractOperationKey(s, store, opType, opMap, tableName)
		if err != nil {
			var opErr *opParseError
			if errors.As(err, &opErr) && opErr.code == "ResourceNotFound" {
				return nil, ErrTableNotFound
			}
			return nil, ErrInvalidParameter
		}

		keyStr := buildKeyString(tableName, key)
		if opType == "ConditionCheck" {
			if usedWriteKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError", Message: "One or more parameter values were invalid."}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedConditionKeys[keyStr] = true
		} else {
			if usedWriteKeys[keyStr] || usedConditionKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError", Message: "One or more parameter values were invalid."}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedWriteKeys[keyStr] = true
		}

		opNames, namesErr := parseExpressionAttributeNames(opMap)
		if namesErr != nil {
			return nil, ErrInvalidParameter
		}

		opValues, opValsErr := parseExpressionAttributeValues(opMap)
		if opValsErr != nil {
			return nil, ErrInvalidParameter
		}

		op := &writeOperation{
			idx:                          idx,
			opType:                       opType,
			tableName:                    tableName,
			key:                          key,
			conditionExpr:                request.GetStringParam(opMap, "ConditionExpression"),
			exprAttrNames:                opNames,
			exprAttrValues:               opValues,
			returnValuesOnConditionCheck: request.GetStringParam(opMap, "ReturnValuesOnConditionCheckFailure"),
		}

		if opType == "ConditionCheck" && op.conditionExpr == "" {
			// ConditionExpression is required for ConditionCheck (Smithy @required).
			return nil, ErrInvalidParameter
		}

		if opType == "Put" {
			itemData, itemErr := parseItem(opMap["Item"])
			if itemErr != nil || itemData == nil {
				return nil, ErrInvalidParameter
			}
			op.itemData = itemData
		}

		if opType == "Update" {
			// UpdateExpression is required for Update (Smithy @required).
			if request.GetStringParam(opMap, "UpdateExpression") == "" {
				return nil, ErrInvalidParameter
			}
			op.updateReq = opMap
		}

		return op, nil
	}

	return nil, ErrInvalidParameter
}

// opParseError represents a parse error encountered during transaction operation parsing.
type opParseError struct {
	code string
	err  error
}

// Error returns the underlying error message for a transaction operation parse failure.
func (e *opParseError) Error() string {
	return e.err.Error()
}

func extractOperationKey(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, opType string, opMap map[string]interface{}, tableName string) (map[string]*dbstore.AttributeValue, error) {
	table, tblErr := store.Tables().Get(tableName)
	if tblErr != nil {
		return nil, &opParseError{code: "ResourceNotFound", err: tblErr}
	}

	if opType == "Put" {
		itemData, itemErr := parseItem(opMap["Item"])
		if itemErr != nil || itemData == nil {
			return nil, &opParseError{code: "ValidationError", err: fmt.Errorf("invalid item data")}
		}
		key := s.extractKeyFromItem(table, itemData)
		if key == nil {
			return nil, &opParseError{code: "ValidationError", err: fmt.Errorf("failed to extract key")}
		}
		if err := validateItemKeyTypes(table, itemData); err != nil {
			return nil, &opParseError{code: "ValidationError", err: err}
		}
		return key, nil
	}

	key, keyErr := parseKey(opMap["Key"])
	if keyErr != nil || key == nil {
		return nil, &opParseError{code: "ValidationError", err: fmt.Errorf("invalid key")}
	}
	if err := validateKeyTypes(table, key); err != nil {
		return nil, &opParseError{code: "ValidationError", err: err}
	}
	return key, nil
}

func executeTransactWriteItems(ctx context.Context, s *DynamoDBService, store dbstore.DynamoDBStoreInterface, operations []writeOperation, cancellationReasons []CancellationReason) error {
	twoPhase := store.Storage().TwoPhaseTransaction()

	for i := range operations {
		twoPhase.AddValidator(storage.ValidatorFunc(func(ctx context.Context, txn storage.Transaction) error {
			return validateWriteOperation(ctx, txn, store.NewTxn(txn), &operations[i], cancellationReasons)
		}))
	}

	// Contributor write events queue on each executor's transaction wrapper
	// and are applied to the access counters after the commit succeeds; the
	// table metric deltas queue the same way and are applied to the table
	// records. The validators and executors run sequentially inside one
	// storage transaction, so plain collection needs no lock.
	var contributorEvents []dbstore.ContributorWriteEvent
	var metricDeltas map[string]dbstore.TableMetricDelta
	for i := range operations {
		opPtr := &operations[i]
		twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
			dbTxn := store.NewTxn(txn)
			if err := executeWriteOperation(dbTxn, opPtr); err != nil {
				return err
			}
			contributorEvents = append(contributorEvents, dbTxn.TakeContributorWrites()...)
			metricDeltas = dbstore.MergeTableMetricDeltas(metricDeltas, dbTxn.TakeTableMetricDeltas())
			return nil
		}))
	}

	// Add a final executor to capture stream records within the same
	// transaction, ensuring atomicity with the item mutations.
	twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
		dbTxn := store.NewTxn(txn)
		for _, op := range operations {
			if op.opType == "ConditionCheck" {
				continue
			}
			table, tblErr := store.Tables().Get(op.tableName)
			if tblErr != nil || table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
				continue
			}
			s.captureStreamChangeTxn(dbTxn, store, table, streamEventForWrite(op.opType == "Delete", op.streamWasNew), op.key, op.streamNewImage, op.streamOldImage)
		}
		return nil
	}))

	if err := twoPhase.Commit(ctx); err != nil {
		// The two-phase commit wraps validator errors, so the cancellation
		// error must be unwrapped with errors.As to preserve its
		// CancellationReasons on the wire.
		var canceledErr *TransactionCanceledError
		if errors.As(err, &canceledErr) {
			return canceledErr
		}
		return ErrTransactionCanceled
	}

	store.FlushContributorWrites(ctx, contributorEvents)
	store.FlushTableMetrics(metricDeltas)

	return nil
}

func validateWriteOperation(_ context.Context, txn storage.Transaction, dbTxn *dbstore.DynamoDBTxn, op *writeOperation, cancellationReasons []CancellationReason) error {

	var item *dbstore.Item
	itemExists := true

	existingItem, err := dbTxn.GetItem(op.tableName, op.key)
	if err != nil {
		if dbstore.IsItemNotFound(err) {
			itemExists = false
			item = &dbstore.Item{
				TableName:  op.tableName,
				Key:        op.key,
				Attributes: make(map[string]*dbstore.AttributeValue),
			}
		} else {
			return fmt.Errorf("validator get item %s: %w", op.tableName, err)
		}
	} else {
		item = existingItem
	}

	// The billed size of the targeted item as it stands before the
	// operation; a stored item's Attributes carry the full item including
	// its key attributes, and an absent item sizes to zero.
	op.itemSize = dbstore.CalculateItemSize(item.Attributes)

	if op.conditionExpr != "" {
		conditionMet, err := evaluateConditionExpression(item, op.conditionExpr, op.exprAttrNames, op.exprAttrValues)
		if err != nil {
			return fmt.Errorf("validator condition check %s: %w", op.tableName, err)
		}
		if !conditionMet {
			reason := CancellationReason{Code: "ConditionalCheckFailed", Message: "The conditional request failed"}
			if op.returnValuesOnConditionCheck == "ALL_OLD" && item != nil && itemExists {
				reason.Item = buildItemResponse(item.Attributes)
			}
			cancellationReasons[op.idx] = reason
			return NewTransactionCanceledError("Transaction canceled", cancellationReasons)
		}
	}

	return nil
}

func executeWriteOperation(dbTxn *dbstore.DynamoDBTxn, op *writeOperation) error {
	exists, err := dbTxn.ItemExists(op.tableName, op.key)
	if err != nil {
		return fmt.Errorf("check item exists %s: %w", op.tableName, err)
	}

	switch op.opType {
	case "Put":
		return executePutOp(dbTxn, op, exists)
	case "Update":
		return executeUpdateOp(dbTxn, op, exists)
	case "Delete":
		return executeDeleteOp(dbTxn, op, exists)
	case "ConditionCheck":
		return nil
	}

	return nil
}

func executePutOp(dbTxn *dbstore.DynamoDBTxn, op *writeOperation, exists bool) error {
	var oldItem *dbstore.Item
	var oldItemSize int64
	if exists {
		var err error
		oldItem, err = dbTxn.GetItem(op.tableName, op.key)
		if err != nil && !dbstore.IsItemNotFound(err) {
			return fmt.Errorf("put get old item %s: %w", op.tableName, err)
		}
		if oldItem != nil {
			oldItemSize = dbstore.CalculateItemSize(oldItem.Attributes)
		}
	}
	if err := dbTxn.StoreItemWrite(op.tableName, op.key, op.itemData, oldItem, exists, oldItemSize); err != nil {
		return fmt.Errorf("put store item write %s: %w", op.tableName, err)
	}
	// Populate stream capture fields.
	op.streamWasNew = !exists
	op.streamNewImage = op.itemData
	if oldItem != nil {
		op.streamOldImage = oldItem.Attributes
	}
	return nil
}

func executeUpdateOp(dbTxn *dbstore.DynamoDBTxn, op *writeOperation, exists bool) error {
	var oldItem *dbstore.Item
	var oldItemSize int64
	if exists {
		var err error
		oldItem, err = dbTxn.GetItem(op.tableName, op.key)
		if err != nil && !dbstore.IsItemNotFound(err) {
			return fmt.Errorf("update get old item %s: %w", op.tableName, err)
		}
		if oldItem != nil {
			oldItemSize = dbstore.CalculateItemSize(oldItem.Attributes)
			if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
				return fmt.Errorf("update delete index entries %s: %w", op.tableName, err)
			}
		}
	}

	attrs := make(map[string]*dbstore.AttributeValue)
	if oldItem != nil {
		attrs = copyAttributes(oldItem.Attributes)
	} else {
		for k, v := range op.key {
			attrs[k] = v
		}
	}

	updateExpr := request.GetStringParam(op.updateReq, "UpdateExpression")
	if updateExpr != "" {
		table, tableErr := dbTxn.GetTable(op.tableName)
		if tableErr != nil {
			return fmt.Errorf("update get table %s: %w", op.tableName, tableErr)
		}
		names := op.exprAttrNames
		paths := extractUpdatedPaths(updateExpr, names)
		if err := validateNotKeyAttributes(table, paths); err != nil {
			return err
		}
		if err := applyUpdateExpression(attrs, updateExpr, op.exprAttrNames, op.exprAttrValues); err != nil {
			return fmt.Errorf("apply update expression %s: %w", op.tableName, err)
		}
	}

	if err := dbTxn.StoreItemWrite(op.tableName, op.key, attrs, nil, exists, oldItemSize); err != nil {
		return fmt.Errorf("update store item write %s: %w", op.tableName, err)
	}
	// Populate stream capture fields.
	op.streamWasNew = !exists
	op.streamNewImage = attrs
	if oldItem != nil {
		op.streamOldImage = oldItem.Attributes
	}
	return nil
}

func executeDeleteOp(dbTxn *dbstore.DynamoDBTxn, op *writeOperation, exists bool) error {
	var oldItem *dbstore.Item
	var oldItemSize int64
	if exists {
		var err error
		oldItem, err = dbTxn.GetItem(op.tableName, op.key)
		if err != nil && !dbstore.IsItemNotFound(err) {
			return fmt.Errorf("delete get old item %s: %w", op.tableName, err)
		}
		if oldItem != nil {
			oldItemSize = dbstore.CalculateItemSize(oldItem.Attributes)
		}
	}
	if err := dbTxn.DeleteItemWrite(op.tableName, op.key, oldItem, exists, oldItemSize); err != nil {
		return fmt.Errorf("delete item write %s: %w", op.tableName, err)
	}
	// Populate stream capture fields.
	if oldItem != nil {
		op.streamOldImage = oldItem.Attributes
	}
	return nil
}

// transactReplayReadUnits sizes the read capacity a client-token replay
// reports for each table: the idempotency contract documents that a replay
// returns the number of read capacity units consumed in reading the item,
// and reading each item of a transaction costs two RCUs at 4 KB granularity
// (one to prepare, one to commit). The basis is the item each operation
// carried — the same size bases the write report uses — captured at
// execution time, because a later replay cannot reconstruct pre-images: a
// delete's item is already gone by then.
func transactReplayReadUnits(operations []writeOperation) map[string]float64 {
	readUnits := make(map[string]float64)
	for _, op := range operations {
		var size int64
		switch op.opType {
		case "Put":
			size = dbstore.CalculateItemSize(op.itemData)
		case "Update":
			size = dbstore.CalculateItemSize(op.streamNewImage)
		default: // Delete and ConditionCheck charge by the targeted item.
			size = op.itemSize
		}
		readUnits[op.tableName] += transactItemReadUnits(size)
	}
	return readUnits
}

// transactItemWriteUnits returns the write capacity a single written item
// consumes in a transactional write: DynamoDB performs two underlying
// writes per item — one to prepare the transaction, one to commit it —
// over the item's size rounded up to 1 KB multiples, and consumes that
// capacity even when the transaction does not succeed. A delete is charged
// by the deleted item's size; an item that did not exist still costs the
// minimum two units.
func transactItemWriteUnits(itemSizeBytes int64) float64 {
	units := (itemSizeBytes + dbstore.WriteCapacityUnitBytes - 1) / dbstore.WriteCapacityUnitBytes
	if units < 1 {
		units = 1
	}
	return float64(units * 2)
}

// buildTransactWriteResponse assembles the TransactWriteItems response. The
// capacity report follows the transaction contract: every item costs two
// underlying operations sized to the item it carries — writes (Put by the
// new item, Update by the post-image, Delete by the pre-image) at 1 KB
// granularity, condition checks as reads of the checked item at 4 KB
// granularity — aggregated per table in TransactItems first-appearance
// order, the ordering the API reference states for the ConsumedCapacity
// list. A replay with the same client token reports the recorded per-table
// read capacity units instead of the write units, per the idempotency
// contract: the replay reads the items, it does not write them again.
func buildTransactWriteResponse(params map[string]interface{}, operations []writeOperation, replay bool, replayReadUnits map[string]float64) map[string]interface{} {
	resp := map[string]interface{}{}

	returnConsumedCapacity := getReturnConsumedCapacity(params)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		writeUnits := make(map[string]float64)
		readUnits := make(map[string]float64)
		var tableOrder []string
		seenTables := make(map[string]bool)
		for _, op := range operations {
			if !seenTables[op.tableName] {
				seenTables[op.tableName] = true
				tableOrder = append(tableOrder, op.tableName)
			}
			switch op.opType {
			case "Put":
				writeUnits[op.tableName] += transactItemWriteUnits(dbstore.CalculateItemSize(op.itemData))
			case "Update":
				writeUnits[op.tableName] += transactItemWriteUnits(dbstore.CalculateItemSize(op.streamNewImage))
			case "Delete":
				writeUnits[op.tableName] += transactItemWriteUnits(op.itemSize)
			case "ConditionCheck":
				readUnits[op.tableName] += transactItemReadUnits(op.itemSize)
			}
		}
		var consumedCapacities []map[string]interface{}
		for _, tableName := range tableOrder {
			if replay {
				rcu := replayReadUnits[tableName]
				consumedCapacities = append(consumedCapacities, map[string]interface{}{
					"TableName":         tableName,
					"CapacityUnits":     rcu,
					"ReadCapacityUnits": rcu,
				})
				continue
			}
			reads := readUnits[tableName]
			writes := writeUnits[tableName]
			capacity := map[string]interface{}{
				"TableName":     tableName,
				"CapacityUnits": reads + writes,
			}
			if reads > 0 {
				capacity["ReadCapacityUnits"] = reads
			}
			if writes > 0 {
				capacity["WriteCapacityUnits"] = writes
			}
			consumedCapacities = append(consumedCapacities, capacity)
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp
}
