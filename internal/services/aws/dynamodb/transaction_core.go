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

	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		for _, gi := range getItems {
			dbItem, err := txn.GetItem(gi.tableName, gi.key)
			if err != nil {
				if !dbstore.IsItemNotFound(err) {
					return fmt.Errorf("transact get item on %s: %w", gi.tableName, err)
				}
				responses = append(responses, map[string]interface{}{"Item": nil})
				continue
			}

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
		// TransactGetItems is always strongly consistent; AWS charges
		// 2x RCU (1.0 per read item) compared to eventually consistent.
		const rcuPerItem = 1.0
		tableNames := make(map[string]bool)
		for _, gi := range getItems {
			tableNames[gi.tableName] = true
		}
		var consumedCapacities []map[string]interface{}
		for tableName := range tableNames {
			consumedCapacities = append(consumedCapacities, buildConsumedCapacityResponse(tableName, float64(len(responses))*rcuPerItem))
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
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
		recordedHash, state, found, lookupErr := store.Idempotency().Lookup(clientRequestToken)
		if lookupErr != nil {
			unlock()
			return nil, lookupErr
		}
		if found {
			unlock()
			if recordedHash != requestHash {
				return nil, ErrIdempotentParameterMismatch
			}
			// A record written before the state field existed was only
			// stored after success, so an empty state counts as completed.
			if state != dbstore.IdempotencyStateCompleted && state != "" {
				return nil, ErrTransactionInProgress
			}
			// Replay: parse for the response shape only; no execution, no
			// stream, replication, or capacity side effects.
			replayReasons := make([]CancellationReason, len(transactItems))
			for i := range replayReasons {
				replayReasons[i] = CancellationReason{Code: "None"}
			}
			replayOps, replayErr := parseTransactWriteItems(s, store, transactItems, replayReasons)
			if replayErr != nil {
				return nil, replayErr
			}
			return buildTransactWriteResponse(in.Parameters, replayOps, true), nil
		}
		if claimErr := store.Idempotency().Record(clientRequestToken, requestHash,
			dbstore.IdempotencyStateInProgress,
			time.Now().Add(idempotencyWindowMinutes*time.Minute)); claimErr != nil {
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
			time.Now().Add(idempotencyWindowMinutes*time.Minute)); recordErr != nil {
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
		eventName := dbstore.StreamEventModify
		if op.opType == "Delete" {
			eventName = dbstore.StreamEventRemove
		} else if op.streamWasNew {
			eventName = dbstore.StreamEventInsert
		}
		s.sendToKinesisDestinations(table, eventName, op.key, op.streamNewImage, op.streamOldImage)

		// Committed transaction writes replicate to global table replica
		// regions just like single-item writes do; Update replicates its
		// committed post-image as a put.
		switch op.opType {
		case "Put":
			if op.itemData != nil {
				s.replicateToGlobalTableReplicas(store, reqCtx.GetRegion(), op.tableName, replicaPutOp(table, op.key, op.itemData))
			}
		case "Update":
			if op.streamNewImage != nil {
				s.replicateToGlobalTableReplicas(store, reqCtx.GetRegion(), op.tableName, replicaPutOp(table, op.key, op.streamNewImage))
			}
		case "Delete":
			s.replicateToGlobalTableReplicas(store, reqCtx.GetRegion(), op.tableName, replicaDeleteOp(table, op.key))
		}
	}

	resp := buildTransactWriteResponse(in.Parameters, operations, false)
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

		if _, err := store.Tables().Get(tableName); err != nil {
			return nil, ErrTableNotFound
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

	for _, op := range operations {
		op := op
		twoPhase.AddValidator(storage.ValidatorFunc(func(ctx context.Context, txn storage.Transaction) error {
			return validateWriteOperation(ctx, txn, store.NewTxn(txn), op, cancellationReasons)
		}))
	}

	// Contributor write events queue on each executor's transaction wrapper
	// and are applied to the access counters after the commit succeeds; the
	// validators and executors run sequentially inside one storage
	// transaction, so plain collection needs no lock.
	var contributorEvents []dbstore.ContributorWriteEvent
	for i := range operations {
		opPtr := &operations[i]
		twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
			dbTxn := store.NewTxn(txn)
			if err := executeWriteOperation(dbTxn, opPtr); err != nil {
				return err
			}
			contributorEvents = append(contributorEvents, dbTxn.TakeContributorWrites()...)
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
			eventName := dbstore.StreamEventModify
			if op.opType == "Delete" {
				eventName = dbstore.StreamEventRemove
			} else if op.streamWasNew {
				eventName = dbstore.StreamEventInsert
			}
			s.captureStreamChangeTxn(dbTxn, store, table, eventName, op.key, op.streamNewImage, op.streamOldImage)
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

	return nil
}

func validateWriteOperation(_ context.Context, txn storage.Transaction, dbTxn *dbstore.DynamoDBTxn, op writeOperation, cancellationReasons []CancellationReason) error {

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
			oldItemSize = calculateItemSize(oldItem.Attributes)
			if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
				return fmt.Errorf("put delete index entries %s: %w", op.tableName, err)
			}
		}
	}
	if err := dbTxn.PutItem(op.tableName, op.key, op.itemData); err != nil {
		return fmt.Errorf("put item %s: %w", op.tableName, err)
	}
	newItem := &dbstore.Item{
		TableName:  op.tableName,
		Key:        op.key,
		Attributes: op.itemData,
	}
	if err := dbTxn.PutIndexEntries(op.tableName, newItem); err != nil {
		return fmt.Errorf("put index entries %s: %w", op.tableName, err)
	}
	// Populate stream capture fields.
	op.streamWasNew = !exists
	op.streamNewImage = op.itemData
	if oldItem != nil {
		op.streamOldImage = oldItem.Attributes
	}
	return updateTableMetrics(dbTxn, op.tableName, exists, oldItemSize, calculateItemSize(op.itemData))
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
			oldItemSize = calculateItemSize(oldItem.Attributes)
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
		if tableErr == nil {
			names := op.exprAttrNames
			paths := extractUpdatedPaths(updateExpr, names)
			if err := validateNotKeyAttributes(table, paths); err != nil {
				return err
			}
		}
		if err := applyUpdateExpression(attrs, updateExpr, op.exprAttrNames, op.exprAttrValues); err != nil {
			return fmt.Errorf("apply update expression %s: %w", op.tableName, err)
		}
	}

	if err := dbTxn.PutItem(op.tableName, op.key, attrs); err != nil {
		return fmt.Errorf("update put item %s: %w", op.tableName, err)
	}
	newItem := &dbstore.Item{
		TableName:  op.tableName,
		Key:        op.key,
		Attributes: attrs,
	}
	if err := dbTxn.PutIndexEntries(op.tableName, newItem); err != nil {
		return fmt.Errorf("update put index entries %s: %w", op.tableName, err)
	}
	// Populate stream capture fields.
	op.streamWasNew = !exists
	op.streamNewImage = attrs
	if oldItem != nil {
		op.streamOldImage = oldItem.Attributes
	}
	return updateTableMetrics(dbTxn, op.tableName, exists, oldItemSize, calculateItemSize(attrs))
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
			oldItemSize = calculateItemSize(oldItem.Attributes)
			if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
				return fmt.Errorf("delete index entries %s: %w", op.tableName, err)
			}
		}
	}
	if err := dbTxn.DeleteItem(op.tableName, op.key); err != nil {
		return fmt.Errorf("delete item %s: %w", op.tableName, err)
	}
	if exists {
		if err := dbTxn.UpdateItemCount(op.tableName, -1); err != nil {
			return fmt.Errorf("delete decrement item count %s: %w", op.tableName, err)
		}
		if oldItemSize > 0 {
			if err := dbTxn.UpdateTableSize(op.tableName, -oldItemSize); err != nil {
				return fmt.Errorf("delete adjust table size %s: %w", op.tableName, err)
			}
		}
	}
	// Populate stream capture fields.
	if oldItem != nil {
		op.streamOldImage = oldItem.Attributes
	}
	return nil
}

func updateTableMetrics(dbTxn *dbstore.DynamoDBTxn, tableName string, existed bool, oldSize, newSize int64) error {
	if !existed {
		if err := dbTxn.UpdateItemCount(tableName, 1); err != nil {
			return fmt.Errorf("increment item count %s: %w", tableName, err)
		}
		if err := dbTxn.UpdateTableSize(tableName, newSize); err != nil {
			return fmt.Errorf("update table size %s: %w", tableName, err)
		}
	} else if newSize != oldSize {
		if err := dbTxn.UpdateTableSize(tableName, newSize-oldSize); err != nil {
			return fmt.Errorf("adjust table size %s: %w", tableName, err)
		}
	}
	return nil
}

// transactCapacityUnitsPerTable is the flat per-table capacity charge of a
// transactional write: transactional operations are charged at twice the
// normal rate.
const transactCapacityUnitsPerTable = 2.0

// buildTransactWriteResponse assembles the TransactWriteItems response. The
// initial execution reports write capacity units; a replay with the same
// client token reports read capacity units, as the idempotency contract
// documents.
func buildTransactWriteResponse(params map[string]interface{}, operations []writeOperation, replay bool) map[string]interface{} {
	resp := map[string]interface{}{}

	returnConsumedCapacity := getReturnConsumedCapacity(params)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableNames := make(map[string]bool)
		for _, op := range operations {
			tableNames[op.tableName] = true
		}
		var consumedCapacities []map[string]interface{}
		for tableName := range tableNames {
			capacity := map[string]interface{}{
				"TableName":     tableName,
				"CapacityUnits": transactCapacityUnitsPerTable,
			}
			if replay {
				capacity["ReadCapacityUnits"] = transactCapacityUnitsPerTable
			} else {
				capacity["WriteCapacityUnits"] = transactCapacityUnitsPerTable
			}
			consumedCapacities = append(consumedCapacities, capacity)
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp
}
