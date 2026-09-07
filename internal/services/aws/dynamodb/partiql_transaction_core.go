package dynamodb

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type statementType int

const (
	statementTypeUnknown statementType = iota
	statementTypeRead
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
func (s *DynamoDBService) recordStatementChange(txn *dbstore.DynamoDBTxn, store dbstore.DynamoDBStoreInterface, changes *[]txnStatementChange, ch txnStatementChange) {
	s.captureStreamChangeTxn(txn, store, ch.table, ch.eventName, ch.keys, ch.newImage, ch.oldImage)
	*changes = append(*changes, ch)
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

	if hasRead && hasWrite {
		return nil, ErrTransactionConflict
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	responses := make([]map[string]interface{}, len(parsedStatements))

	if hasRead {
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				result, err := s.executePartiQLSelectInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				if err != nil {
					return err
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
						item = items[0]
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
			}
			return nil
		})
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
		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				upperStmt := strings.ToUpper(strings.TrimSpace(ps.statement))
				var result interface{}
				var execErr error

				switch {
				case strings.HasPrefix(upperStmt, "INSERT"):
					result, execErr = s.executePartiQLInsertInTxn(ctx, reqCtx, store, txn, ps.statement, ps.params, &changes)
				case strings.HasPrefix(upperStmt, "UPDATE"):
					result, execErr = s.executePartiQLUpdateInTxn(ctx, reqCtx, store, txn, ps.statement, ps.params, &changes)
				case strings.HasPrefix(upperStmt, "DELETE"):
					result, execErr = s.executePartiQLDeleteInTxn(ctx, reqCtx, store, txn, ps.statement, ps.params, &changes)
				}

				if execErr != nil {
					reason, cancels := singletonCancellationReason(execErr)
					if !cancels {
						return execErr
					}
					cancellationReasons[i] = reason
					return NewTransactionCanceledError("Transaction canceled", cancellationReasons)
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
						item = items[0]
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
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

	resp := map[string]interface{}{
		"Responses": responses,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(in.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableNames := make(map[string]bool)
		for _, ps := range parsedStatements {
			if tn := extractTableNameFromStatement(ps.statement); tn != "" {
				tableNames[tn] = true
			}
		}
		var consumedCapacities []map[string]interface{}
		for tableName := range tableNames {
			consumedCapacities = append(consumedCapacities, buildConsumedCapacityResponse(tableName, 2.0))
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
}

// singletonCancellationReason classifies a failed write statement's error
// for the transaction-cancellation envelope. Only errors a singleton
// operation reports against the stored data — a duplicate key, a
// multi-item match, a clause applied onto an incompatible attribute type —
// cancel the transaction and carry a per-statement reason; request-level
// failures (malformed statements, unresolvable parameters, unknown tables,
// storage faults) propagate directly with their own error code.
func singletonCancellationReason(err error) (CancellationReason, bool) {
	switch {
	case errors.Is(err, ErrConditionalCheckFailed):
		return CancellationReason{Code: "ConditionalCheckFailed", Message: "The conditional request failed"}, true
	case errors.Is(err, ErrTypeMismatch):
		return CancellationReason{Code: "ValidationError", Message: ErrTypeMismatch.Message}, true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.Code == "com.amazonaws.dynamodb.v20120810#ValidationException" {
		// The multi-item-match rejection of the single-item contract is the
		// only producer of the DynamoDB-namespaced ValidationException code.
		return CancellationReason{Code: "ValidationError", Message: apiErr.Message}, true
	}
	return CancellationReason{}, false
}

// The per-statement-type transactional execution engines below own the full
// write path of one PartiQL statement inside ExecuteTransaction's store
// transaction — table resolution, the single-item contract, index and
// metric maintenance, and the stream/replication capture the commit fires.

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

	scannedCount := len(items)
	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, buildItemResponse(item.Attributes))
	}

	return map[string]interface{}{
		"Items":        result,
		"Count":        len(result),
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLInsertInTxn(ctx context.Context, reqCtx *request.RequestContext, store dbstore.DynamoDBStoreInterface, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams, changes *[]txnStatementChange) (interface{}, error) {
	tableName, itemData, err := parseInsertStatementWithParams(statement, params)
	if err != nil {
		return nil, err
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	// A table mid-restore (CREATING) must not be mutated: a transactional
	// write statement follows the write-requires-ACTIVE rule of the
	// single-item writes.
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	keyAttrs := buildKeyFromSchema(table.KeySchema, itemData)
	if keyAttrs == nil {
		return nil, ErrInvalidParameter
	}

	_, err = txn.GetItem(tableName, keyAttrs)
	if err != nil {
		if !dbstore.IsItemNotFound(err) {
			return nil, err
		}
	} else {
		return nil, ErrConditionalCheckFailed
	}

	if err := txn.StoreItemWrite(tableName, keyAttrs, itemData, nil, false, 0); err != nil {
		return nil, err
	}

	s.recordStatementChange(txn, store, changes, txnStatementChange{
		table:     table,
		eventName: dbstore.StreamEventInsert,
		keys:      keyAttrs,
		newImage:  itemData,
		replicaOp: s.replicaPutOp(table, keyAttrs, itemData),
	})

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        0,
		"ScannedCount": 0,
	}, nil
}

func (s *DynamoDBService) executePartiQLUpdateInTxn(ctx context.Context, reqCtx *request.RequestContext, store dbstore.DynamoDBStoreInterface, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams, changes *[]txnStatementChange) (interface{}, error) {
	tableName, clauses, whereExpr := parseUpdateStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	// At least one clause must be present.
	if len(clauses.setAssignments) == 0 && len(clauses.removeAttrs) == 0 &&
		len(clauses.addAssignments) == 0 && len(clauses.deleteAssignments) == 0 {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	// A table mid-restore (CREATING) must not be mutated: a transactional
	// write statement follows the write-requires-ACTIVE rule of the
	// single-item writes.
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	// The key attributes identify the item being updated; no UPDATE clause
	// may write them.
	if err := validateNotKeyAttributes(table, updateClauseTargetNames(clauses)); err != nil {
		return nil, err
	}

	// The AWS single-item contract: the WHERE clause carries a partition-key
	// equality and the statement matches at most one item.
	item, _, scannedCount, err := matchSingleTargetItem(
		func(pkValue string, cb func(*dbstore.Item) error) error {
			return txn.ScanByPartitionKey(tableName, pkValue, cb)
		}, table, whereExpr, params, "UPDATE", false)
	if err != nil {
		return nil, err
	}

	updatedCount := 0
	if item != nil {
		oldSize := dbstore.CalculateItemSize(item.Attributes)

		oldItem := &dbstore.Item{
			TableName:  tableName,
			Key:        copyAttributes(item.Key),
			Attributes: copyAttributes(item.Attributes),
		}

		if err := applySetAssignments(item.Attributes, clauses.setAssignments, params); err != nil {
			return nil, err
		}
		applyRemoveAttrs(item.Attributes, clauses.removeAttrs)
		if err := applyAddAssignments(item.Attributes, clauses.addAssignments, params); err != nil {
			return nil, err
		}
		if err := applyDeleteAssignments(item.Attributes, clauses.deleteAssignments, params); err != nil {
			return nil, err
		}

		if err := txn.StoreItemWrite(tableName, item.Key, item.Attributes, oldItem, true, oldSize); err != nil {
			return nil, err
		}

		s.recordStatementChange(txn, store, changes, txnStatementChange{
			table:     table,
			eventName: dbstore.StreamEventModify,
			keys:      item.Key,
			newImage:  item.Attributes,
			oldImage:  oldItem.Attributes,
			replicaOp: s.replicaPutOp(table, item.Key, item.Attributes),
		})

		updatedCount = 1
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        updatedCount,
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLDeleteInTxn(ctx context.Context, reqCtx *request.RequestContext, store dbstore.DynamoDBStoreInterface, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams, changes *[]txnStatementChange) (interface{}, error) {
	tableName, whereExpr := parseDeleteStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	// A table mid-restore (CREATING) must not be mutated: a transactional
	// write statement follows the write-requires-ACTIVE rule of the
	// single-item writes.
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	// The AWS single-item contract: the WHERE clause carries a partition-key
	// equality and the statement matches at most one item.
	item, _, scannedCount, err := matchSingleTargetItem(
		func(pkValue string, cb func(*dbstore.Item) error) error {
			return txn.ScanByPartitionKey(tableName, pkValue, cb)
		}, table, whereExpr, params, "DELETE", false)
	if err != nil {
		return nil, err
	}

	deletedCount := 0
	if item != nil {
		if err := txn.DeleteItemWrite(tableName, item.Key, item, true, dbstore.CalculateItemSize(item.Attributes)); err != nil {
			return nil, err
		}

		s.recordStatementChange(txn, store, changes, txnStatementChange{
			table:     table,
			eventName: dbstore.StreamEventRemove,
			keys:      item.Key,
			oldImage:  item.Attributes,
			replicaOp: s.replicaDeleteOp(table, item.Key),
		})

		deletedCount = 1
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        deletedCount,
		"ScannedCount": scannedCount,
	}, nil
}
