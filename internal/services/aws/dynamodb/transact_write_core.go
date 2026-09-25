package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// TransactWriteItems Core — single validation and persistence path of the
// transactional write plane: the ClientRequestToken idempotency window,
// operation parsing, the two-phase executor with its validators and stream
// capture, and the response / capacity report. The TransactGetItems half
// lives in transact_get_core.go.
// ---------------------------------------------------------------------------

// recordIdempotencyCompletion promotes a claimed token to completed after
// the transaction has committed. A variable so the post-commit
// failure-safe path (log-and-continue) can be fault-injected in tests.
var recordIdempotencyCompletion = func(store dbstore.DynamoDBStoreInterface, token, requestHash string, expiresAt time.Time, readUnits map[string]float64) error {
	return store.Idempotency().Record(token, requestHash, dbstore.IdempotencyStateCompleted, expiresAt, readUnits)
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

	clientRequestToken := request.GetStringParam(in.Parameters, "ClientRequestToken")
	if !validateClientRequestToken(clientRequestToken) {
		return nil, ErrInvalidParameter
	}

	// Response-shaping enums are request validation: an unknown value is
	// rejected before anything executes.
	if _, err := getReturnConsumedCapacity(in.Parameters); err != nil {
		return nil, err
	}
	collectionMetrics, err := getItemCollectionMetricsSetting(in.Parameters)
	if err != nil {
		return nil, err
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
	// The deferred release covers every post-claim failure path; the
	// committed path marks the guard before the completion record runs.
	tokenGuard := &tokenClaimGuard{store: store}
	defer tokenGuard.release()
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
		tokenGuard.claim(clientRequestToken)
		unlock()
	}

	cancellationReasons := make([]CancellationReason, len(transactItems))
	for i := range cancellationReasons {
		cancellationReasons[i] = CancellationReason{Code: "None"}
	}

	operations, err := parseTransactWriteItems(s, store, transactItems, cancellationReasons)
	if err != nil {
		return nil, err
	}

	// The aggregate of the supplied items cannot exceed the transaction
	// limit — the request is rejected before any execution. Updates and
	// deletes carry no supplied item; their per-item outcomes stay under
	// the executor's own item-size rejection.
	var aggregateSize int64
	for _, op := range operations {
		if op.itemData != nil {
			aggregateSize += dbstore.CalculateItemSize(op.itemData)
		}
	}
	if aggregateSize > dbstore.MaxTransactionAggregateBytes {
		return nil, errTransactWriteAggregateExceeded()
	}

	// Serialise against other in-flight transactions and single-item
	// writes on the same items: without cross-transaction coordination the
	// storage engine's committed-state reads and last-writer-wins commits
	// would let two transactions addressing one item both commit, silently
	// losing the earlier write. A rejected item-level request fails the
	// whole transaction with a TransactionConflict cancellation reason at
	// that item's slot.
	txnLockKeys := make([]string, len(operations))
	lockKeySlot := make(map[string]int, len(operations))
	for i := range operations {
		k := itemLockKey(reqCtx.Region, operations[i].tableName, operations[i].key)
		txnLockKeys[i] = k
		lockKeySlot[k] = operations[i].idx
	}
	if conflictKey, free := tryLockItems(itemLockModeTxn, txnLockKeys); !free {
		conflictReasons := make([]CancellationReason, len(transactItems))
		for i := range conflictReasons {
			conflictReasons[i] = CancellationReason{Code: "None"}
		}
		conflictReasons[lockKeySlot[conflictKey]] = CancellationReason{
			Code:    "TransactionConflict",
			Message: transactionConflictReasonMessage,
		}
		// The transaction never started, so the deferred claim guard
		// releases the token — a client retrying it after the transient
		// conflict must not meet the in-progress record.
		return nil, NewTransactionCanceledError("Transaction canceled", conflictReasons)
	}
	defer unlockItems(itemLockModeTxn, txnLockKeys)

	if err := executeTransactWriteItems(ctx, s, store, operations, cancellationReasons); err != nil {
		return nil, err
	}
	// The writes are durable: the claim survives as the completion path's
	// record.
	tokenGuard.markCommitted()

	if tokenGuard.claimed {
		// The transaction is already durable at this point: a failure to
		// promote the token record must not turn the committed success
		// into an error response — the client would retry into a token
		// still marked in-progress, the Kinesis and global-table
		// propagation below would be skipped, and the completion would
		// never be recorded. The failure is logged and the completion
		// path continues; in-window retries keep answering the documented
		// TransactionInProgress for the rest of the window, which is the
		// bounded outcome, and no code path re-executes the commit.
		if recordErr := recordIdempotencyCompletion(store, clientRequestToken, requestHash,
			time.Now().Add(idempotencyWindowMinutes*time.Minute),
			transactReplayReadUnits(operations)); recordErr != nil {
			logs.Error("Failed to record idempotency completion after committed transaction", logs.Err(recordErr))
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
	if collectionMetrics == "SIZE" {
		if metrics := buildItemCollectionMetricsPerTable(metricsWrites); metrics != nil {
			resp["ItemCollectionMetrics"] = metrics
		}
	}
	return resp, nil
}

// writeOperation represents a single write operation within a transaction.
type writeOperation struct {
	idx           int
	opType        string
	tableName     string
	key           map[string]*dbstore.AttributeValue
	itemData      map[string]*dbstore.AttributeValue
	updateReq     map[string]interface{}
	conditionExpr string
	// condition is the operation's ConditionExpression, parsed and
	// validated when the operation was built — a malformed condition or
	// undefined substitution rejects the whole request before the
	// transaction starts.
	condition                    *compiledCondition
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
	usedItemKeys := make(map[string]bool)
	var operations []writeOperation

	for idx, item := range transactItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		op, err := parseWriteOperation(s, store, idx, itemMap, usedItemKeys, cancellationReasons)
		if err != nil {
			return nil, err
		}
		operations = append(operations, *op)
	}

	return operations, nil
}

// errTransactionDuplicateItem is the read plane's request-validation
// rejection: the TransactGetItems documentation counts the TransactItems
// members as targeting "distinct items" and enumerates no duplicate-item
// circumstance among the transaction's cancellation reasons, so a Get
// action repeated on one item is rejected with the request, before any
// execution. The write plane answers its duplicates as a cancellation —
// see parseWriteOperation.
func errTransactionDuplicateItem() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"Transaction request cannot include multiple operations on one item", http.StatusBadRequest)
}

// duplicateItemCancellation is the write plane's answer to a duplicate
// target: a TransactionCanceledException whose reasons carry the offending
// action's slot as the documented ValidationError pair, with every earlier
// slot at its None default — the detection runs before any execution, so
// nothing is written.
func duplicateItemCancellation(idx int, cancellationReasons []CancellationReason) error {
	cancellationReasons[idx] = CancellationReason{
		Code:    "ValidationError",
		Message: "One or more parameter values were invalid.",
	}
	return NewTransactionCanceledError("Transaction canceled", cancellationReasons)
}

// errTransactWriteAggregateExceeded is the write plane's aggregate-size
// rejection, verbatim from the TransactWriteItems documentation: "The
// aggregate size of the items in the transaction exceeds 4 MB."
func errTransactWriteAggregateExceeded() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"The aggregate size of the items in the transaction exceeds 4 MB", http.StatusBadRequest)
}

func parseWriteOperation(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, idx int, itemMap map[string]interface{}, usedItemKeys map[string]bool, cancellationReasons []CancellationReason) (*writeOperation, error) {
	// TransactWriteItem is a union: exactly one action member may be
	// carried, and one of them must be.
	present := 0
	for _, opType := range []string{"Put", "Update", "Delete", "ConditionCheck"} {
		if v, ok := itemMap[opType]; ok && v != nil {
			present++
		}
	}
	if present != 1 {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
			"TransactWriteItem must contain exactly one of Put, Update, Delete or ConditionCheck", http.StatusBadRequest)
	}

	for _, opType := range []string{"Put", "Update", "Delete", "ConditionCheck"} {
		opMap, ok := itemMap[opType].(map[string]interface{})
		if !ok {
			continue
		}

		// Request-shape failures (malformed members, unknown or missing
		// tables) are request-level validation errors answered with
		// ValidationException or ResourceNotFoundException. Semantic
		// per-item outcomes — two actions targeting one item, condition
		// failures during execution — cancel the transaction and are
		// reported through CancellationReasons: the operation
		// documentation enumerates "more than one action ... targets the
		// same item" among the circumstances under which the request is
		// canceled, and lists the reason as Code "ValidationError" with
		// the message below.
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
			if errors.As(err, &opErr) {
				if opErr.code == "ResourceNotFound" {
					return nil, ErrTableNotFound
				}
				// A key that failed a schema validation already carries its
				// own ValidationException (type mismatch, schema mismatch,
				// value length); flattening it to the generic invalid-
				// parameter form would hide which constraint failed.
				var apiErr *APIError
				if errors.As(opErr.err, &apiErr) {
					return nil, apiErr
				}
			}
			return nil, ErrInvalidParameter
		}

		keyStr := buildKeyString(tableName, key)
		// "No two of them can operate on the same item" — the rule is
		// uniform across all four actions, so one map records every
		// operation's item and any second entry is the duplicate
		// cancellation.
		if usedItemKeys[keyStr] {
			return nil, duplicateItemCancellation(idx, cancellationReasons)
		}
		usedItemKeys[keyStr] = true

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

		if op.conditionExpr != "" {
			compiled, condErr := compileConditionExpression(op.conditionExpr, opNames, opValues)
			if condErr != nil {
				return nil, condErr
			}
			op.condition = compiled
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
		if !validateKeyAttributeValue(key) {
			return nil, &opParseError{code: "ValidationError", err: ErrInvalidParameter}
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
	if err := validateKeySchemaMembership(table, key); err != nil {
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
				// An item grown past the size cap by the operation's
				// changes rejects the whole transaction with the
				// documented ValidationError cancellation reason at the
				// offending operation's slot.
				if errors.Is(err, errTransactionItemSizeExceeded) {
					cancellationReasons[opPtr.idx] = CancellationReason{
						Code:    "ValidationError",
						Message: errTransactionItemSizeExceeded.Error(),
					}
					return NewTransactionCanceledError("Transaction canceled", cancellationReasons)
				}
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
			if err := s.captureStreamChangeTxn(dbTxn, store, table, streamEventForWrite(op.opType == "Delete", op.streamWasNew), op.key, op.streamNewImage, op.streamOldImage); err != nil {
				return err
			}
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
		// A validation error surfacing from an executor — an update
		// expression that does not parse, a key-attribute guard hit, a
		// type mismatch — is a user error, and the operation documents
		// user errors among the conditions that reject the request: the
		// same ValidationException the single-item plane answers, never
		// a bare cancellation. TransactionCanceledError is checked first
		// because it embeds APIError. Storage faults carry no APIError
		// and keep the bare cancellation shape.
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			return apiErr
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

	if op.condition != nil {
		if !op.condition.matches(item) {
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

// errTransactionItemSizeExceeded is the executor-phase sentinel for an
// item the transaction's own changes grew past the documented single-item
// cap. The TransactWriteItems contract rejects the entire request when "An
// item size becomes too large (bigger than 400 KB) ... because of changes
// made by the transaction"; the message text is the service's own form for
// an oversized item.
var errTransactionItemSizeExceeded = errors.New("Item size has exceeded the maximum allowed size")

func executePutOp(dbTxn *dbstore.DynamoDBTxn, op *writeOperation, exists bool) error {
	if dbstore.CalculateItemSize(op.itemData) > dbstore.MaxItemSizeBytes {
		return errTransactionItemSizeExceeded
	}
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
		paths, exprErr := extractUpdatedPaths(updateExpr, names)
		if exprErr != nil {
			return exprErr
		}
		if err := validateNotKeyAttributes(table, paths); err != nil {
			return err
		}
		if err := applyUpdateExpression(attrs, updateExpr, op.exprAttrNames, op.exprAttrValues); err != nil {
			return fmt.Errorf("apply update expression %s: %w", op.tableName, err)
		}
	}

	// The post-update size is checked before the write: an expression can
	// grow an item past the cap the whole-item arrivals never exceed.
	if dbstore.CalculateItemSize(attrs) > dbstore.MaxItemSizeBytes {
		return errTransactionItemSizeExceeded
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

	// Both callers validate the enum at their Core entry, so the error path
	// here is defensive: an (impossible) unknown value shapes no capacity.
	returnConsumedCapacity, capErr := getReturnConsumedCapacity(params)
	if capErr != nil {
		return resp
	}
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
