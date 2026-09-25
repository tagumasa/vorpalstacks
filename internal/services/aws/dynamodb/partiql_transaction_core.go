package dynamodb

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type statementType int

const (
	statementTypeRead statementType = iota
	statementTypeWrite
)

// executeTransactionInput carries the already-typed TransactStatements member
// plus the raw wire parameters (consumed for the ClientRequestToken
// idempotency member and the ReturnConsumedCapacity reporting).
type executeTransactionInput struct {
	TransactStatements []interface{}
	Parameters         map[string]interface{}
}

// txnStatementChange records what one write statement inside an
// ExecuteTransaction did: the stream capture runs within the shared storage
// transaction, and the record joins the post-commit propagation queue so
// Kinesis destinations and global table replicas see the write only once
// the transaction has committed.
type txnStatementChange struct {
	table     *dbstore.Table
	eventName dbstore.StreamEventName
	keys      map[string]*dbstore.AttributeValue
	newImage  map[string]*dbstore.AttributeValue
	oldImage  map[string]*dbstore.AttributeValue
	replicaOp func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error
}

// recordStatementChange captures the statement's stream record inside the
// running transaction and queues its post-commit propagation.
func (s *DynamoDBService) recordStatementChange(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, changes *[]txnStatementChange, ch txnStatementChange) error {
	if err := s.captureStreamChangeTxn(txn, store, ch.table, ch.eventName, ch.keys, ch.newImage, ch.oldImage); err != nil {
		return err
	}
	*changes = append(*changes, ch)
	return nil
}

// transactionalWritePlane returns the write plane of one statement inside
// ExecuteTransaction's store transaction: tables and items resolve through
// the transaction, and the change capture records the stream record plus
// the queued post-commit propagation inside the running transaction.
func (s *DynamoDBService) transactionalWritePlane(store dbstore.DynamoDBStoreInterface, txn *dbstore.DynamoDBTxn, changes *[]txnStatementChange) partiqlWritePlane {
	return partiqlWritePlane{
		txn:      txn,
		getTable: txn.GetTable,
		getItem:  txn.GetItem,
		capture: func(table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) error {
			var replicaOp func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error
			if eventName == dbstore.StreamEventRemove {
				replicaOp = s.replicaDeleteOp(table, keys)
			} else {
				replicaOp = s.replicaPutOp(table, keys, newImage)
			}
			return s.recordStatementChange(txn, store, changes, txnStatementChange{
				table:     table,
				eventName: eventName,
				keys:      keys,
				newImage:  newImage,
				oldImage:  oldImage,
				replicaOp: replicaOp,
			})
		},
	}
}

// executeTransactionCore is the single validation and persistence path of the
// PartiQL ExecuteTransaction operation: statement classification with the
// read/write exclusivity rule, then either one snapshot view or one
// transactional update across every statement.
func (s *DynamoDBService) executeTransactionCore(ctx context.Context, reqCtx *request.RequestContext, in executeTransactionInput) (map[string]interface{}, error) {
	statements := in.TransactStatements

	if len(statements) == 0 {
		return nil, ErrInvalidParameter
	}

	if len(statements) > transactMaxItems {
		return nil, ErrInvalidParameter
	}

	clientRequestToken := request.GetStringParam(in.Parameters, "ClientRequestToken")
	if !validateClientRequestToken(clientRequestToken) {
		return nil, ErrInvalidParameter
	}

	// Response-shaping enums are request validation: an unknown value is
	// rejected before anything executes.
	returnConsumedCapacity, capErr := getReturnConsumedCapacity(in.Parameters)
	if capErr != nil {
		return nil, capErr
	}

	parsedStatements := make([]struct {
		statement string
		params    *partiQLParams
		stmtType  statementType
	}, 0, len(statements))

	var hasRead, hasWrite bool

	for _, stmt := range statements {
		stmtMap, ok := stmt.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		statement, _ := stmtMap["Statement"].(string)
		if !validatePartiQLStatement(statement) {
			return nil, ErrInvalidParameter
		}

		params := parsePartiQLParams(stmtMap)
		// Normalise `?` placeholders into their explicit whole-statement
		// :vN form before any engine parse, and reject a parameter list
		// that does not carry exactly one value per placeholder.
		normalised, placeholderCount := preparePartiQLStatement(statement)
		if err := validateParameterCount(placeholderCount, params); err != nil {
			return nil, err
		}
		statement = normalised

		upperStmt := strings.ToUpper(strings.TrimSpace(statement))
		var stmtType statementType
		switch {
		case strings.HasPrefix(upperStmt, "SELECT"):
			stmtType = statementTypeRead
			hasRead = true
		case strings.HasPrefix(upperStmt, "INSERT"),
			strings.HasPrefix(upperStmt, "UPDATE"),
			strings.HasPrefix(upperStmt, "DELETE"):
			// A transactional write statement returns no values ("This
			// statement doesn't return any values for Write operations
			// (INSERT, UPDATE, or DELETE)"), so the RETURNING clause the
			// standalone plane serves has no transactional surface: a
			// statement carrying one is rejected, never executed with the
			// requested return silently dropped.
			if _, returning := splitReturningSuffix(statement); returning != "" {
				return nil, ErrInvalidParameter
			}
			stmtType = statementTypeWrite
			hasWrite = true
		default:
			return nil, ErrInvalidParameter
		}

		parsedStatements = append(parsedStatements, struct {
			statement string
			params    *partiQLParams
			stmtType  statementType
		}{statement, params, stmtType})
	}

	// The model requires the entire transaction to consist of either read
	// statements or write statements — the mix is a request-shape
	// violation answered as a validation error naming the rule, never as
	// a conflict with another transaction on an item.
	if hasRead && hasWrite {
		return nil, errMixedStatementPlanes("transaction")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// A ClientRequestToken makes the call idempotent: the member carries
	// the idempotencyToken trait, and no ExecuteTransaction-specific window
	// is documented anywhere (the member's documentation text is a
	// NextToken copy-paste artifact), so the window semantics follow
	// TransactWriteItems' documented behaviour — the token is valid for
	// ten minutes after the first request completes. The token keyspace is
	// namespaced per API, so a token reused across TransactWriteItems and
	// ExecuteTransaction addresses two independent records.
	const executeTransactionTokenNamespace = "ExecuteTransaction|"
	tokenKey := executeTransactionTokenNamespace + clientRequestToken
	requestHash := ""
	replayOnly := false
	// The deferred release covers every post-claim failure path; the
	// committed path marks the guard before the completion record runs.
	tokenGuard := &tokenClaimGuard{store: store}
	defer tokenGuard.release()
	if clientRequestToken != "" {
		requestHash = hashTransactWriteRequest(in.Parameters)

		// The claim section is serialised per token so concurrent retries
		// cannot both execute: the loser observes the in-progress record
		// and fails with the documented in-progress error.
		unlock := s.lockClientRequestToken(tokenKey)
		recordedHash, state, _, found, lookupErr := store.Idempotency().Lookup(tokenKey)
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
			// Replay: reads re-run side-effect-free; writes must not
			// re-execute.
			replayOnly = true
		} else {
			if claimErr := store.Idempotency().Record(tokenKey, requestHash,
				dbstore.IdempotencyStateInProgress,
				time.Now().Add(idempotencyWindowMinutes*time.Minute), nil); claimErr != nil {
				unlock()
				return nil, claimErr
			}
			tokenGuard.claim(tokenKey)
			unlock()
		}
	}

	responses := make([]map[string]interface{}, len(parsedStatements))

	if hasRead {
		// Reads carry no side effects, so a replayed transaction re-runs
		// them against the current state ("the responses to the calls
		// might not be the same").
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				result, err := s.executePartiQLSelectInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				if err != nil {
					// Table existence is request-level state: a statement
					// naming a table that does not exist answers the
					// documented ResourceNotFoundException directly, the
					// shape the model's error list carries — never the raw
					// storage error the resolution produced.
					if tableMissing(err) {
						return ErrResourceNotFound
					}
					return err
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok {
						// The response shape carries one Item per statement.
						// A WHERE clause that matches several items answers
						// data the shape cannot represent: the documentation
						// defines no multi-match behaviour for transactional
						// reads, and a rejected read discloses the limit
						// where a silent truncation would hand the client
						// one arbitrary match and drop the rest.
						if len(items) > 1 {
							return NewAPIError("com.amazon.coral.validate#ValidationException",
								"A read statement in a transaction must address a single item",
								http.StatusBadRequest)
						}
						if len(items) == 1 {
							item = items[0]
						}
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
			}
			return nil
		})
	} else if replayOnly {
		// A completed token replays the write transaction without
		// re-executing it: executing write statements answer an empty
		// response slot (the protocol omits unset members), which the
		// replay reproduces directly.
		for i := range parsedStatements {
			responses[i] = map[string]interface{}{}
		}
	} else {
		var changes []txnStatementChange
		// A failing singleton statement cancels the whole transaction: its
		// error travels in the ordered CancellationReasons envelope, and
		// every unaffected statement reports the literal code "None" with no
		// message. Request-level failures stay direct (below).
		cancellationReasons := make([]CancellationReason, len(parsedStatements))
		for i := range cancellationReasons {
			cancellationReasons[i] = CancellationReason{Code: "None"}
		}

		// Serialise against the other in-flight write planes on the same
		// items — the registry coordination TransactWriteItems applies: the
		// storage engine's transactions read committed state and commit
		// last-writer-wins, so without this acquisition a concurrent
		// transaction and this statement batch would both commit and one
		// write would be silently lost. Each statement's verb derives its
		// target key through the shared derivation (INSERT from the parsed
		// item, UPDATE and DELETE from the WHERE clause's primary-key
		// equalities); a statement whose key does not resolve cleanly
		// contributes nothing — it fails in-engine with its own
		// per-statement semantics and locks nothing.
		lockKeys := make([]string, 0, len(parsedStatements))
		lockKeySlot := make(map[string]int, len(parsedStatements))
		for i, ps := range parsedStatements {
			if ps.stmtType != statementTypeWrite {
				continue
			}
			k, derr := derivePartiQLWriteLockKey(store, reqCtx.Region, ps.statement, ps.params, true)
			if derr != nil {
				return nil, derr
			}
			if k == "" {
				continue
			}
			lockKeys = append(lockKeys, k)
			lockKeySlot[k] = i
		}
		if len(lockKeys) > 0 {
			if conflictKey, free := tryLockItems(itemLockModeTxn, lockKeys); !free {
				// The operation's declared error set carries
				// TransactionCanceledException and no
				// TransactionConflictException, so the conflict answers
				// through the same cancellation envelope the engine path
				// uses: the conflicting statement's slot names the
				// TransactionConflict reason.
				cancellationReasons[lockKeySlot[conflictKey]] = CancellationReason{
					Code:    "TransactionConflict",
					Message: transactionConflictReasonMessage,
				}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			defer unlockItems(itemLockModeTxn, lockKeys)
		}

		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			plane := s.transactionalWritePlane(store, txn, &changes)
			for i, ps := range parsedStatements {
				upperStmt := strings.ToUpper(strings.TrimSpace(ps.statement))
				var execErr error

				switch {
				case strings.HasPrefix(upperStmt, "INSERT"):
					_, execErr = plane.runInsert(ps.statement, ps.params)
				case strings.HasPrefix(upperStmt, "UPDATE"):
					_, execErr = plane.runUpdate(ps.statement, ps.params, "")
				case strings.HasPrefix(upperStmt, "DELETE"):
					_, execErr = plane.runDelete(ps.statement, ps.params, "")
				}

				if execErr != nil {
					// Every error a singleton write statement returns cancels
					// the transaction ("If any of the singleton INSERT,
					// UPDATE, or DELETE operations return an error, the
					// transactions are canceled with the
					// TransactionCanceledException exception"), with the
					// reason carrying the statement's own error — no error
					// class escapes the envelope raw.
					cancellationReasons[i] = singletonCancellationReason(execErr)
					return NewTransactionCanceledError("Transaction canceled", cancellationReasons)
				}

				// A singleton write statement's response slot carries no
				// Item member: the write shapes expose no readable item,
				// and the protocol omits unset members rather than
				// carrying an explicit null.
				responses[i] = map[string]interface{}{}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		// The transaction committed: fire the asynchronous propagation every
		// write statement owes, in statement order.
		for _, ch := range changes {
			s.emitChangePropagation(store, reqCtx.GetRegion(), ch.table, ch.eventName, ch.keys, ch.newImage, ch.oldImage, ch.replicaOp)
		}
	}

	if err != nil {
		return nil, err
	}
	// Past this point the write plane (when present) has committed, so
	// the claim survives as the completion path's record.
	tokenGuard.markCommitted()

	if tokenGuard.claimed {
		// The transaction is already durable: a failure to promote the
		// token record must not turn the committed success into an error
		// response — the failure is logged and the completion path
		// continues, with in-window retries answering the documented
		// TransactionInProgress for the rest of the window.
		if recordErr := recordIdempotencyCompletion(store, tokenKey, requestHash,
			time.Now().Add(idempotencyWindowMinutes*time.Minute), nil); recordErr != nil {
			logs.Error("Failed to record idempotency completion after committed transaction", logs.Err(recordErr))
		}
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		// Per-statement accounting in first-appearance order — the same
		// per-item contract TransactWriteItems applies: every write
		// statement costs two underlying write operations at the minimum
		// charge, every read statement its transactional read charge.
		writeUnits := make(map[string]float64)
		readUnits := make(map[string]float64)
		var tableOrder []string
		seenTables := make(map[string]bool)
		for _, ps := range parsedStatements {
			tn := extractTableNameFromStatement(ps.statement)
			if tn == "" {
				continue
			}
			if !seenTables[tn] {
				seenTables[tn] = true
				tableOrder = append(tableOrder, tn)
			}
			if ps.stmtType == statementTypeWrite {
				writeUnits[tn] += transactItemWriteUnits(0)
			} else {
				readUnits[tn] += transactItemReadUnits(0)
			}
		}
		var consumedCapacities []map[string]interface{}
		for _, tn := range tableOrder {
			reads := readUnits[tn]
			writes := writeUnits[tn]
			capacity := map[string]interface{}{
				"TableName":     tn,
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

	return resp, nil
}

// partiqlErrorPlane selects how the shared statement-error classifier
// renders the one class its two consuming shapes disagree on: an APIError
// outside the documented code list.
type partiqlErrorPlane int

const (
	// partiqlErrorBatch renders onto the BatchStatementErrorCodeEnum, which
	// defines no code for such exceptions — they report InternalServerError.
	partiqlErrorBatch partiqlErrorPlane = iota
	// partiqlErrorCancellation renders onto a transaction cancellation
	// reason, whose Code member is a free string — the exception's own
	// short code carries through.
	partiqlErrorCancellation
)

// partiqlStatementErrorPair projects a failed statement's error onto the
// short code/message pair the batch response slots and the transaction
// cancellation reasons share: known exception classes render their
// documented reason codes with the exception's own message, a raw
// not-found storage error renders ResourceNotFound, and anything else —
// storage faults included — renders InternalServerError rather than
// escaping raw. The plane selects the default branch's rendering.
func partiqlStatementErrorPair(err error, plane partiqlErrorPlane) (string, string) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case "com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException":
			return "ConditionalCheckFailed", "The conditional request failed"
		case "com.amazon.coral.validate#ValidationException":
			return "ValidationError", apiErr.Message
		case "com.amazonaws.dynamodb.v20120810#ResourceNotFoundException":
			return "ResourceNotFound", apiErr.Message
		case "com.amazonaws.dynamodb.v20120810#TransactionConflictException":
			return "TransactionConflict", apiErr.Message
		case "com.amazonaws.dynamodb.v20120810#DuplicateItemException":
			return "DuplicateItem", apiErr.Message
		}
		if plane == partiqlErrorCancellation {
			return shortCancellationCode(apiErr.Code), apiErr.Message
		}
		// Every other exception — ResourceInUseException included, whose
		// code the batch error enum does not define — reports through
		// InternalServerError with the exception's own message.
		return "InternalServerError", apiErr.Message
	}
	if tableMissing(err) {
		return "ResourceNotFound", "Requested resource not found"
	}
	return "InternalServerError", err.Error()
}

// singletonCancellationReason classifies a failed write statement's error
// for the transaction-cancellation envelope: every error a singleton
// operation reports cancels the transaction and carries a per-statement
// reason ("If any of the singleton INSERT, UPDATE, or DELETE operations
// return an error, the transactions are canceled ... and the cancellation
// reason code includes the errors from the individual singleton
// operations"). Known classes render their documented reason codes; an
// APIError outside the documented code list renders its exception name
// (the list is illustrative, and the reason's Code member is a free
// string in the model).
func singletonCancellationReason(err error) CancellationReason {
	code, message := partiqlStatementErrorPair(err, partiqlErrorCancellation)
	return CancellationReason{Code: code, Message: message}
}

// tableMissing reports whether an error is a store-side table resolution
// failure of the not-found kind, in either sentinel family the store
// layers use: the common not-found class and the dynamodb store's own
// table sentinel.
func tableMissing(err error) bool {
	return errors.Is(err, dbstore.ErrTableNotFound) || commonstore.IsNotFound(err)
}

// shortCancellationCode renders a wire error code as a cancellation-reason
// code: the exception name without its shape namespace and Exception
// suffix, the form the documented reason codes use (ConditionalCheckFailed
// for ConditionalCheckFailedException).
func shortCancellationCode(wireCode string) string {
	if idx := strings.LastIndex(wireCode, "#"); idx != -1 {
		wireCode = wireCode[idx+1:]
	}
	return strings.TrimSuffix(wireCode, "Exception")
}

// The transactional read engine below owns the read path of one PartiQL
// SELECT inside ExecuteTransaction's snapshot view; the write statements
// run through the shared per-verb engines on the transactional write
// plane, which carry the table resolution, the single-item contract, index
// and metric maintenance, and the stream/replication capture the commit
// fires.

func (s *DynamoDBService) executePartiQLSelectInTxn(ctx context.Context, reqCtx *request.RequestContext, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams) (interface{}, error) {
	tableName, whereExpr := parseSelectStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	var items []*dbstore.Item

	pkName := getHashKeyName(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	if pkValue != "" {
		err = txn.ScanByPartitionKey(tableName, pkValue, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		err = txn.Scan(tableName, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, buildItemResponse(item.Attributes))
	}

	return map[string]interface{}{
		"Items": result,
	}, nil
}

// batchExecuteStatementInput carries the already-typed Statements member
// plus the raw wire parameters (consumed for the ReturnConsumedCapacity
// reporting).
type batchExecuteStatementInput struct {
	Statements []interface{}
	Parameters map[string]interface{}
}

// batchReadKeyRuleMessage states the model's batch read rule verbatim: the
// BatchExecuteStatement documentation requires each read statement to
// "specify an equality condition on all key attributes", enforcing that a
// batch SELECT returns at most a single item.
const batchReadKeyRuleMessage = "Each read statement in a BatchExecuteStatement must specify an equality condition on all key attributes"

// batchReadStatementMeetsKeyRule checks the batch read rule against the
// target table's key schema. A statement whose table cannot be resolved
// passes here and reports the missing table through its own execution
// error, as does a statement that does not parse.
func batchReadStatementMeetsKeyRule(store dbstore.DynamoDBStoreInterface, statement string, params *partiQLParams) bool {
	tableName, whereExpr := parseSelectStatement(statement)
	if tableName == "" {
		return true
	}
	if !store.Tables().Exists(tableName) {
		return true
	}
	table, err := store.Tables().Get(tableName)
	if err != nil || table == nil {
		return true
	}
	_, ok := resolvePrimaryKeyFromWhere(whereExpr, table.KeySchema, params)
	return ok
}

// batchExecuteStatementCore is the single validation and persistence path of
// the BatchExecuteStatement operation: batch-size validation, the
// model-documented batch rules (a batch is all reads or all writes; every
// read statement carries an equality on all key attributes), then the
// per-statement loop delegating each statement to the ExecuteStatement Core
// and mapping a failing statement onto its per-statement Error slot. A
// statement's failure never fails the batch response itself; the batch
// rules are request-level and do.
func (s *DynamoDBService) batchExecuteStatementCore(ctx context.Context, reqCtx *request.RequestContext, in batchExecuteStatementInput) (map[string]interface{}, error) {
	if !validatePartiQLBatchCount(len(in.Statements)) {
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
		return nil, err
	}

	// Classification pass: every well-formed statement is a read (SELECT)
	// or a write (INSERT/UPDATE/DELETE). The model forbids mixing the two
	// in one batch, so the classification runs before anything executes.
	type batchStatement struct {
		raw       interface{}
		stmtMap   map[string]interface{}
		statement string
		params    *partiQLParams
		isRead    bool
	}
	validationError := func(message string) map[string]interface{} {
		return map[string]interface{}{
			"Error": map[string]interface{}{
				"Code":    "ValidationError",
				"Message": message,
			},
		}
	}

	var hasRead, hasWrite bool
	statements := make([]batchStatement, 0, len(in.Statements))
	for _, raw := range in.Statements {
		stmtMap, ok := raw.(map[string]interface{})
		if !ok {
			statements = append(statements, batchStatement{raw: raw})
			continue
		}
		statement, _ := stmtMap["Statement"].(string)
		if statement == "" {
			statements = append(statements, batchStatement{raw: raw, stmtMap: stmtMap})
			continue
		}
		bs := batchStatement{raw: raw, stmtMap: stmtMap, statement: statement, params: parsePartiQLParams(stmtMap)}
		switch {
		case strings.HasPrefix(strings.ToUpper(strings.TrimSpace(statement)), "SELECT"):
			bs.isRead = true
			hasRead = true
		case strings.HasPrefix(strings.ToUpper(strings.TrimSpace(statement)), "INSERT"),
			strings.HasPrefix(strings.ToUpper(strings.TrimSpace(statement)), "UPDATE"),
			strings.HasPrefix(strings.ToUpper(strings.TrimSpace(statement)), "DELETE"):
			hasWrite = true
		}
		statements = append(statements, bs)
	}

	if hasRead && hasWrite {
		return nil, errMixedStatementPlanes("batch")
	}

	var responses []map[string]interface{}
	// The batch response aggregates the per-statement capacity entries per
	// table: one summed ConsumedCapacity per distinct table, in first-seen
	// order — the batch plane's list shape, the BatchGetItem convention.
	var tableOrder []string
	tableTotals := make(map[string]float64)

	for _, bs := range statements {
		if bs.stmtMap == nil {
			responses = append(responses, validationError("Invalid statement format"))
			continue
		}
		if bs.statement == "" {
			responses = append(responses, validationError("Statement is required"))
			continue
		}

		// The batch read rule: a SELECT without an equality on every key
		// attribute of its table would scan and silently truncate, so it is
		// rejected before execution in its own response slot.
		if bs.isRead && !batchReadStatementMeetsKeyRule(store, bs.statement, bs.params) {
			responses = append(responses, validationError(batchReadKeyRuleMessage))
			continue
		}

		tableName := extractTableNameFromStatement(bs.statement)
		consistentRead, crErr := validateBoolParam(bs.stmtMap, "ConsistentRead", false)
		if crErr != nil {
			responses = append(responses, validationError("ConsistentRead must be a boolean"))
			continue
		}

		// The batch-level ReturnConsumedCapacity shapes every per-statement
		// execution: each statement's result then carries its own entry for
		// the per-table aggregation below.
		stmtParams := map[string]interface{}{
			"Statement":      bs.statement,
			"Parameters":     bs.params.Parameters,
			"ConsistentRead": consistentRead,
		}
		if returnConsumedCapacity != "" {
			stmtParams["ReturnConsumedCapacity"] = returnConsumedCapacity
		}
		result, err := s.executeStatementCore(ctx, reqCtx, executeStatementInput{
			Parameters: stmtParams,
		})

		response := map[string]interface{}{
			"TableName": tableName,
		}

		if err != nil {
			errorCode, errorMessage := partiqlStatementErrorPair(err, partiqlErrorBatch)
			response["Error"] = map[string]interface{}{
				"Code":    errorCode,
				"Message": errorMessage,
			}
		} else {
			var item interface{} = nil
			if resultMap, ok := result.(map[string]interface{}); ok {
				if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
					item = items[0]
				}
				if cc, ok := resultMap["ConsumedCapacity"].(map[string]interface{}); ok {
					ccTable, _ := cc["TableName"].(string)
					units, _ := cc["CapacityUnits"].(float64)
					if _, seen := tableTotals[ccTable]; !seen {
						tableOrder = append(tableOrder, ccTable)
					}
					tableTotals[ccTable] += units
				}
			}
			// A statement that matched no item carries no Item member —
			// the protocol omits unset members rather than carrying an
			// explicit null.
			if item != nil {
				response["Item"] = item
			}
		}

		responses = append(responses, response)
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	if len(tableOrder) > 0 {
		aggregated := make([]map[string]interface{}, 0, len(tableOrder))
		for _, tn := range tableOrder {
			aggregated = append(aggregated, buildConsumedCapacityResponse(tn, tableTotals[tn]))
		}
		resp["ConsumedCapacity"] = aggregated
	}

	return resp, nil
}
