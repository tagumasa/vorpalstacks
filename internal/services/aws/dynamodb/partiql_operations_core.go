package dynamodb

// This file carries the PartiQL execution engines: the ExecuteStatement Core,
// the shared per-verb write engines on the write plane both the standalone
// statements and the ExecuteTransaction statements run on, and the
// statement-shape helpers (WHERE/key extraction, table-name resolution,
// ordering) those engines consume; the write-clause appliers they drive live
// in partiql_write_apply.go. The handlers in partiql_operations.go and
// partiql_transaction.go parse and dispatch; everything here owns validation
// and persistence.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

// executeStatementInput carries the raw wire parameters of ExecuteStatement.
type executeStatementInput struct {
	Parameters map[string]interface{}
}

// executeStatementCore is the single validation and persistence path of the
// ExecuteStatement operation: statement validation, placeholder
// normalisation with the parameter-count check, the response-shaping enums,
// then dispatch to the per-verb engine and the response's ConsumedCapacity
// assembly.
func (s *DynamoDBService) executeStatementCore(ctx context.Context, reqCtx *request.RequestContext, in executeStatementInput) (interface{}, error) {
	statement := request.GetStringParam(in.Parameters, "Statement")
	if !validatePartiQLStatement(statement) {
		return nil, ErrInvalidParameter
	}

	params := parsePartiQLParams(in.Parameters)
	// Normalise `?` placeholders into their explicit whole-statement :vN
	// form before dispatch: the engines' parses (including the manual
	// UPDATE path's per-segment re-parses) then share one numbering, and
	// the count check rejects a parameter list that does not carry exactly
	// one value per placeholder.
	statement, placeholderCount := preparePartiQLStatement(statement)
	if err := validateParameterCount(placeholderCount, params); err != nil {
		return nil, err
	}
	consistentRead, crErr := validateBoolParam(in.Parameters, "ConsistentRead", false)
	if crErr != nil {
		return nil, crErr
	}
	// Limit targets the model's PositiveIntegerObject (range 1+): an
	// explicitly-present value below 1 is rejected, never read as "no
	// limit" — GetIntParam cannot distinguish an explicit 0 from an
	// absent member, so the presence check carries the range decision.
	limit := request.GetIntParam(in.Parameters, "Limit")
	if _, present := in.Parameters["Limit"]; present && !validateExecuteStatementLimit(limit) {
		return nil, ErrInvalidParameter
	}
	nextToken := request.GetStringParam(in.Parameters, "NextToken")
	returnValuesOnConditionCheckFailure := request.GetStringParam(in.Parameters, "ReturnValuesOnConditionCheckFailure")

	// Response-shaping enums are request validation: an unknown value is
	// rejected before anything executes.
	returnConsumedCapacity, err := getReturnConsumedCapacity(in.Parameters)
	if err != nil {
		return nil, err
	}

	upperStmt := strings.ToUpper(strings.TrimSpace(statement))
	var result interface{}

	switch {
	case strings.HasPrefix(upperStmt, "SELECT"):
		result, err = s.executePartiQLSelectEnhanced(ctx, reqCtx, statement, params, consistentRead, limit, nextToken)
	case strings.HasPrefix(upperStmt, "INSERT"):
		result, err = s.executePartiQLInsert(ctx, reqCtx, statement, params)
	case strings.HasPrefix(upperStmt, "UPDATE"):
		result, err = s.executePartiQLUpdate(ctx, reqCtx, statement, params, returnValuesOnConditionCheckFailure)
	case strings.HasPrefix(upperStmt, "DELETE"):
		result, err = s.executePartiQLDelete(ctx, reqCtx, statement, params, returnValuesOnConditionCheckFailure)
	default:
		return nil, ErrInvalidParameter
	}

	if err != nil {
		return nil, err
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableName := extractTableNameFromStatement(statement)
		// The statement-level charge reports the unit helpers' minimum: an
		// eventually consistent read of a sub-4 KB item, a sub-1 KB write.
		capacityUnits := itemReadUnits(0, false)
		if strings.HasPrefix(upperStmt, "INSERT") || strings.HasPrefix(upperStmt, "UPDATE") || strings.HasPrefix(upperStmt, "DELETE") {
			capacityUnits = itemWriteUnits(0)
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			resultMap["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, capacityUnits)
		}
	}

	return result, nil
}

func (s *DynamoDBService) executePartiQLSelectEnhanced(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, consistentRead bool, limit int, nextToken string) (interface{}, error) {
	tableName, whereExpr, orderBy, selectCols := parseSelectStatementWithOrderBy(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !store.Tables().Exists(tableName) {
		return nil, ErrTableNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}

	pkName := getHashKeyName(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	// Decode pagination offset before scanning so we can calculate
	// how many items to collect before early termination. A malformed
	// token — not base64, not the offset JSON, or a non-positive offset —
	// is a rejected request, not an absent cursor: restarting from the
	// first page would silently re-serve it.
	startOffset := 0
	if nextToken != "" {
		decoded, decErr := base64.StdEncoding.DecodeString(nextToken)
		if decErr != nil {
			return nil, ErrInvalidParameter
		}
		var offset int
		if json.Unmarshal(decoded, &offset) != nil || offset <= 0 {
			return nil, ErrInvalidParameter
		}
		startOffset = offset
	}

	// When no ORDER BY is present, filter inside the scan callback and
	// break early after collecting enough items. This reduces memory
	// from O(N) to O(limit) for filtered queries on large tables.
	needed := 0
	if orderBy == nil && limit > 0 {
		needed = startOffset + limit
	}

	var items []*dbstore.Item
	// sawMore records that a matching item was rejected because the window
	// was full — the proof that a further page exists, which the early
	// termination would otherwise discard along with the item.
	sawMore := false

	scanCallback := func(item *dbstore.Item) error {
		if whereExpr != nil && !evaluateExpr(item.Attributes, whereExpr, params) {
			return nil
		}
		if needed > 0 && len(items) >= needed {
			sawMore = true
			return errScanSufficient
		}
		items = append(items, item)
		return nil
	}

	if pkValue != "" {
		err = store.Items().ScanByPartitionKey(tableName, pkValue, scanCallback)
	} else {
		err = store.Items().Scan(tableName, scanCallback)
	}
	if err != nil && !errors.Is(err, errScanSufficient) {
		return nil, err
	}

	if orderBy != nil {
		items = sortItemsByOrderBy(items, orderBy)
	}
	// A stale NextToken replayed against a shrunken result set carries an
	// offset past the end; that is the last page, not an error.
	if startOffset >= len(items) {
		items = nil
	} else {
		items = items[startOffset:]
	}

	var newNextToken string
	if limit > 0 && (sawMore || len(items) > limit) {
		items = items[:limit]
		offsetBytes, _ := json.Marshal(startOffset + limit)
		newNextToken = base64.StdEncoding.EncodeToString(offsetBytes)
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		attrs := item.Attributes
		if selectCols != nil {
			filtered := make(map[string]*dbstore.AttributeValue)
			for _, col := range selectCols {
				if v, ok := attrs[col]; ok {
					filtered[col] = v
				}
			}
			attrs = filtered
		}
		result = append(result, buildItemResponse(attrs))
	}

	_ = consistentRead

	// ExecuteStatementOutput defines Items, NextToken and ConsumedCapacity
	// (assembled by the Core) alone: the Query/Scan response members Count
	// and ScannedCount have no counterpart on this operation.
	resp := map[string]interface{}{
		"Items": result,
	}
	if newNextToken != "" {
		resp["NextToken"] = newNextToken
	}
	return resp, nil
}

// partiqlWritePlane carries the adapters one shared per-verb write engine
// runs against: the storage transaction every read and write goes through,
// the table and item resolution strategies, and the change-capture
// strategy. The standalone ExecuteStatement engines supply a fresh
// single-statement transaction with the store's committed-state reads and
// an immediate stream capture; the ExecuteTransaction statements supply
// the caller's running transaction with transactional reads and the
// capture that queues the change for the post-commit propagation. One
// engine body per verb then owns the whole write path exactly once.
type partiqlWritePlane struct {
	txn      *dbstore.DynamoDBTxn
	getTable func(name string) (*dbstore.Table, error)
	getItem  func(tableName string, key map[string]*dbstore.AttributeValue) (*dbstore.Item, error)
	capture  func(table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) error
}

// partiqlWriteResult records what a shared write engine body did, for the
// calling plane's post-commit propagation and response shaping. A DELETE
// whose key addresses no stored item commits as a no-op and reports a nil
// result. The RETURNING clause the statement carried (canonical form) and
// the top-level attribute names it wrote feed the write-response builder's
// image choice.
type partiqlWriteResult struct {
	table     *dbstore.Table
	keys      map[string]*dbstore.AttributeValue
	newImage  map[string]*dbstore.AttributeValue
	oldImage  map[string]*dbstore.AttributeValue
	returning string
	modified  []string
}

// standaloneWritePlane returns the write plane of a single-statement
// ExecuteStatement write: tables resolve through the store (a missing table
// answers the service's table-not-found sentinel) and the stream change is
// captured inside the statement's own transaction.
func (s *DynamoDBService) standaloneWritePlane(store dbstore.DynamoDBStoreInterface, txn *dbstore.DynamoDBTxn) partiqlWritePlane {
	return partiqlWritePlane{
		txn: txn,
		getTable: func(name string) (*dbstore.Table, error) {
			if !store.Tables().Exists(name) {
				return nil, ErrTableNotFound
			}
			return store.Tables().Get(name)
		},
		getItem: store.Items().Get,
		capture: func(table *dbstore.Table, eventName dbstore.StreamEventName, keys, newImage, oldImage map[string]*dbstore.AttributeValue) error {
			return s.captureStreamChangeTxn(txn, store, table, eventName, keys, newImage, oldImage)
		},
	}
}

// runInsert executes the shared INSERT engine body: table resolution with
// the write-requires-ACTIVE rule, key and item validation against the
// single-item write plane's own bounds, the duplicate-key rejection, the
// store write and the change capture.
func (p partiqlWritePlane) runInsert(statement string, params *partiQLParams) (partiqlWriteResult, error) {
	tableName, itemData, err := parseInsertStatementWithParams(statement, params)
	if err != nil {
		return partiqlWriteResult{}, err
	}

	table, err := p.getTable(tableName)
	if err != nil {
		return partiqlWriteResult{}, err
	}
	// A table mid-restore (CREATING) must not be mutated: a write statement
	// follows the write-requires-ACTIVE rule of the single-item writes.
	if table.Status != dbstore.TableStatusActive {
		return partiqlWriteResult{}, ErrTableNotActive
	}

	keyAttrs := buildKeyFromSchema(table.KeySchema, itemData)
	if keyAttrs == nil {
		return partiqlWriteResult{}, ErrInvalidParameter
	}
	// The single-item write plane's validations bind here too: the item's
	// key and index-key attribute types must match the schema definitions,
	// and the composed item must stay within the documented item size —
	// the same item cannot be legal on the PutItem plane and illegal here.
	if err := validateItemKeyTypes(table, itemData); err != nil {
		return partiqlWriteResult{}, err
	}
	if dbstore.CalculateItemSize(itemData) > dbstore.MaxItemSizeBytes {
		return partiqlWriteResult{}, ErrInvalidParameter
	}

	if _, err = p.txn.GetItem(tableName, keyAttrs); err != nil {
		if !dbstore.IsItemNotFound(err) {
			return partiqlWriteResult{}, err
		}
	} else {
		return partiqlWriteResult{}, ErrConditionalCheckFailed
	}

	if err := p.txn.StoreItemWrite(tableName, keyAttrs, itemData, nil, false, 0); err != nil {
		return partiqlWriteResult{}, err
	}

	if err := p.capture(table, dbstore.StreamEventInsert, keyAttrs, itemData, nil); err != nil {
		return partiqlWriteResult{}, err
	}

	return partiqlWriteResult{table: table, keys: keyAttrs, newImage: itemData}, nil
}

// derivePartiQLWriteLockKey resolves the item registry key one write
// statement's engine body will address — INSERT from the parsed item's
// schema keys, UPDATE and DELETE from the WHERE clause's primary-key
// equalities — the single derivation shared by the standalone write
// plane and the batch transaction. strict selects the failure strategy:
// the transactional caller answers request-level state directly (a
// missing table answers the documented ResourceNotFoundException, never
// a raw storage error and never a per-statement cancellation; a
// non-ACTIVE table answers the ResourceInUseException shape the
// single-write planes apply), while the standalone caller defers to the
// engine body's own per-verb semantics and locks nothing. A statement
// that is not a write, does not parse, or whose key does not resolve
// cleanly contributes no key under either strategy.
func derivePartiQLWriteLockKey(store dbstore.DynamoDBStoreInterface, region, statement string, params *partiQLParams, strict bool) (string, error) {
	upperStmt := strings.ToUpper(strings.TrimSpace(statement))
	switch {
	case strings.HasPrefix(upperStmt, "INSERT"):
		tableName, itemData, err := parseInsertStatementWithParams(statement, params)
		if err != nil {
			return "", nil
		}
		table, terr := store.Tables().Get(tableName)
		if terr != nil || table == nil {
			if !strict {
				return "", nil
			}
			if tableMissing(terr) {
				return "", ErrResourceNotFound
			}
			return "", terr
		}
		if strict && table.Status != dbstore.TableStatusActive {
			return "", ErrTableNotActive
		}
		keyAttrs := buildKeyFromSchema(table.KeySchema, itemData)
		if keyAttrs == nil {
			return "", nil
		}
		return itemLockKey(region, tableName, keyAttrs), nil
	case strings.HasPrefix(upperStmt, "UPDATE"), strings.HasPrefix(upperStmt, "DELETE"):
		var tableName string
		var whereExpr sqlparser.Expr
		if strings.HasPrefix(upperStmt, "UPDATE") {
			tableName, _, whereExpr, _ = parseUpdateStatement(statement)
		} else {
			tableName, whereExpr, _ = parseDeleteStatement(statement)
		}
		if tableName == "" {
			return "", nil
		}
		table, terr := store.Tables().Get(tableName)
		if terr != nil || table == nil {
			if !strict {
				return "", nil
			}
			if tableMissing(terr) {
				return "", ErrResourceNotFound
			}
			return "", terr
		}
		if strict && table.Status != dbstore.TableStatusActive {
			return "", ErrTableNotActive
		}
		keyAttrs, ok := resolvePrimaryKeyFromWhere(whereExpr, table.KeySchema, params)
		if !ok {
			return "", nil
		}
		return itemLockKey(region, tableName, keyAttrs), nil
	default:
		return "", nil
	}
}

// partiQLWriteLockKey is the standalone plane's face of the shared
// derivation: a statement whose key does not derive locks nothing; the
// engine body answers it with its own per-verb semantics.
func (s *DynamoDBService) partiQLWriteLockKey(store dbstore.DynamoDBStoreInterface, reqCtx *request.RequestContext, statement string, params *partiQLParams) string {
	lockKey, _ := derivePartiQLWriteLockKey(store, reqCtx.Region, statement, params, false)
	return lockKey
}

// lockStandaloneWriteStatement acquires the item registry lock a standalone
// write statement owes before its engine body runs. The read-modify-write
// transaction serialises against in-flight TransactWriteItems on the same
// item — the registry coordination the single-item write plane applies; a
// rejection answers TransactionConflictException.
func (s *DynamoDBService) lockStandaloneWriteStatement(store dbstore.DynamoDBStoreInterface, reqCtx *request.RequestContext, statement string, params *partiQLParams) (func(), error) {
	lockKey := s.partiQLWriteLockKey(store, reqCtx, statement, params)
	if lockKey == "" {
		return func() {}, nil
	}
	if _, free := tryLockItems(itemLockModeItem, []string{lockKey}); !free {
		return nil, ErrTransactionConflict
	}
	return func() { unlockItems(itemLockModeItem, []string{lockKey}) }, nil
}

func (s *DynamoDBService) executePartiQLInsert(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	release, lockErr := s.lockStandaloneWriteStatement(store, reqCtx, statement, params)
	if lockErr != nil {
		return nil, lockErr
	}
	defer release()

	var result partiqlWriteResult
	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		var runErr error
		result, runErr = s.standaloneWritePlane(store, txn).runInsert(statement, params)
		return runErr
	})
	if err != nil {
		// The standalone plane answers the documented
		// DuplicateItemException for the duplicate primary key — the
		// transactional plane keeps the ConditionalCheckFailed cancellation
		// reason, whose code list defines no DuplicateItem code.
		if errors.Is(err, ErrConditionalCheckFailed) {
			return nil, ErrDuplicateItem
		}
		return nil, err
	}

	s.emitChangePropagation(store, reqCtx.GetRegion(), result.table, dbstore.StreamEventInsert, result.keys, result.newImage, result.oldImage, s.replicaPutOp(result.table, result.keys, result.newImage))

	return response.EmptyResponse(), nil
}

// runUpdate executes the shared UPDATE engine body: clause presence, the
// write-requires-ACTIVE rule, the key-attribute guard, the single-item
// WHERE contract with the UPDATE zero-match contract, the clause
// application with the post-clause size bound, the store write and the
// change capture.
func (p partiqlWritePlane) runUpdate(statement string, params *partiQLParams, returnValuesOnConditionCheckFailure string) (partiqlWriteResult, error) {
	tableName, clauses, whereExpr, parseErr := parseUpdateStatement(statement)
	if parseErr != nil {
		return partiqlWriteResult{}, parseErr
	}
	if tableName == "" {
		return partiqlWriteResult{}, ErrInvalidParameter
	}

	// At least one clause must be present.
	if len(clauses.setAssignments) == 0 && len(clauses.removeAttrs) == 0 &&
		len(clauses.addAssignments) == 0 && len(clauses.deleteAssignments) == 0 {
		return partiqlWriteResult{}, ErrInvalidParameter
	}

	table, err := p.getTable(tableName)
	if err != nil {
		return partiqlWriteResult{}, err
	}
	// A table mid-restore (CREATING) must not be mutated: a write statement
	// follows the write-requires-ACTIVE rule of the single-item writes.
	if table.Status != dbstore.TableStatusActive {
		return partiqlWriteResult{}, ErrTableNotActive
	}

	// The key attributes identify the item being updated; no UPDATE clause
	// may write them.
	if err := validateNotKeyAttributes(table, updateClauseTargetNames(clauses)); err != nil {
		return partiqlWriteResult{}, err
	}

	// The AWS single-item contract: the WHERE clause resolves to a single
	// primary-key value and the statement addresses the item under that key.
	item, found, conditionHeld, err := matchKeyedTargetItem(
		func(key map[string]*dbstore.AttributeValue) (*dbstore.Item, error) {
			return p.getItem(tableName, key)
		}, table, whereExpr, params, "UPDATE")
	if err != nil {
		return partiqlWriteResult{}, err
	}

	// The UPDATE zero-match contract is unqualified: a WHERE clause that
	// evaluates true for no item fails the condition check — whether the key
	// addresses no stored item or the remaining predicates reject it.
	if !found || !conditionHeld {
		var existing *dbstore.Item
		if found {
			existing = item
		}
		return partiqlWriteResult{}, conditionalCheckFailedError(returnValuesOnConditionCheckFailure, existing)
	}

	oldSize := dbstore.CalculateItemSize(item.Attributes)

	oldItem := &dbstore.Item{
		TableName:  tableName,
		Key:        copyAttributes(item.Key),
		Attributes: copyAttributes(item.Attributes),
	}

	// One statement is one update of one item: every operand and action
	// base reads the pre-statement image — the documented multi-action
	// evaluation basis — while the writes apply to the item.
	before := copyAttributes(item.Attributes)
	if err := applySetAssignments(item.Attributes, before, clauses.setAssignments, params); err != nil {
		return partiqlWriteResult{}, err
	}
	if err := applyRemoveAttrs(item.Attributes, clauses.removeAttrs); err != nil {
		return partiqlWriteResult{}, err
	}
	if err := applyAddAssignments(item.Attributes, before, clauses.addAssignments, params); err != nil {
		return partiqlWriteResult{}, err
	}
	if err := applyDeleteAssignments(item.Attributes, before, clauses.deleteAssignments, params); err != nil {
		return partiqlWriteResult{}, err
	}

	// The post-clause image must stay within the documented item size, the
	// single-item update plane's own bound.
	if dbstore.CalculateItemSize(item.Attributes) > dbstore.MaxItemSizeBytes {
		return partiqlWriteResult{}, ErrInvalidParameter
	}

	if err := p.txn.StoreItemWrite(tableName, item.Key, item.Attributes, oldItem, true, oldSize); err != nil {
		return partiqlWriteResult{}, err
	}

	if err := p.capture(table, dbstore.StreamEventModify, item.Key, item.Attributes, oldItem.Attributes); err != nil {
		return partiqlWriteResult{}, err
	}

	return partiqlWriteResult{
		table:     table,
		keys:      item.Key,
		newImage:  item.Attributes,
		oldImage:  oldItem.Attributes,
		returning: clauses.returning,
		modified:  updateClauseTargetNames(clauses),
	}, nil
}

// buildPartiqlWriteResponse shapes every standalone PartiQL write
// statement's response from the write plane's result: without a RETURNING
// clause the Items list is empty, with one it carries the clause's chosen
// image of the single written item. One builder serves the UPDATE and
// DELETE planes, so the two response faces cannot drift apart.
func buildPartiqlWriteResponse(result *partiqlWriteResult) map[string]interface{} {
	items := []map[string]interface{}{}
	if result != nil && result.returning != "" {
		items = append(items, buildItemResponse(result.returningImage()))
	}
	return map[string]interface{}{"Items": items}
}

// returningImage resolves the RETURNING clause the write recorded into the
// image it names: the full old or new image of the written item, or that
// image restricted to the attributes the statement wrote. A modified
// attribute the chosen image does not carry (a SET target with no old
// value) simply does not appear.
func (r *partiqlWriteResult) returningImage() map[string]*dbstore.AttributeValue {
	var image map[string]*dbstore.AttributeValue
	switch r.returning {
	case "ALL OLD *":
		image = r.oldImage
	case "MODIFIED OLD *":
		image = restrictAttributes(r.oldImage, r.modified)
	case "ALL NEW *":
		image = r.newImage
	case "MODIFIED NEW *":
		image = restrictAttributes(r.newImage, r.modified)
	}
	if image == nil {
		image = map[string]*dbstore.AttributeValue{}
	}
	return image
}

func restrictAttributes(image map[string]*dbstore.AttributeValue, names []string) map[string]*dbstore.AttributeValue {
	restricted := make(map[string]*dbstore.AttributeValue, len(names))
	for _, name := range names {
		if v, ok := image[name]; ok && v != nil {
			restricted[name] = v
		}
	}
	return restricted
}

func (s *DynamoDBService) executePartiQLUpdate(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, returnValuesOnConditionCheckFailure string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	release, lockErr := s.lockStandaloneWriteStatement(store, reqCtx, statement, params)
	if lockErr != nil {
		return nil, lockErr
	}
	defer release()

	var result partiqlWriteResult
	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		var runErr error
		result, runErr = s.standaloneWritePlane(store, txn).runUpdate(statement, params, returnValuesOnConditionCheckFailure)
		return runErr
	})
	if err != nil {
		return nil, err
	}

	s.emitChangePropagation(store, reqCtx.GetRegion(), result.table, dbstore.StreamEventModify, result.keys, result.newImage, result.oldImage, s.replicaPutOp(result.table, result.keys, result.newImage))

	return buildPartiqlWriteResponse(&result), nil
}

// runDelete executes the shared DELETE engine body: the write-requires-
// ACTIVE rule, the single-item WHERE contract with the DELETE zero-match
// contract, the store delete and the change capture. A key no stored item
// carries commits the no-op statement and reports a nil result.
func (p partiqlWritePlane) runDelete(statement string, params *partiQLParams, returnValuesOnConditionCheckFailure string) (*partiqlWriteResult, error) {
	tableName, whereExpr, returning := parseDeleteStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := p.getTable(tableName)
	if err != nil {
		return nil, err
	}
	// A table mid-restore (CREATING) must not be mutated: a write statement
	// follows the write-requires-ACTIVE rule of the single-item writes.
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	// The AWS single-item contract: the WHERE clause resolves to a single
	// primary-key value and the statement addresses the item under that key.
	item, found, conditionHeld, err := matchKeyedTargetItem(
		func(key map[string]*dbstore.AttributeValue) (*dbstore.Item, error) {
			return p.getItem(tableName, key)
		}, table, whereExpr, params, "DELETE")
	if err != nil {
		return nil, err
	}

	// The DELETE zero-match contract is two-tier: a primary key no stored
	// item carries answers SUCCESS with zero items deleted, while a stored
	// item whose further WHERE predicates reject the statement fails the
	// condition check.
	if !found {
		return nil, nil
	}
	if !conditionHeld {
		return nil, conditionalCheckFailedError(returnValuesOnConditionCheckFailure, item)
	}

	if err := p.txn.DeleteItemWrite(tableName, item.Key, item, true, dbstore.CalculateItemSize(item.Attributes)); err != nil {
		return nil, err
	}

	if err := p.capture(table, dbstore.StreamEventRemove, item.Key, nil, item.Attributes); err != nil {
		return nil, err
	}

	return &partiqlWriteResult{table: table, keys: item.Key, oldImage: item.Attributes, returning: returning}, nil
}

func (s *DynamoDBService) executePartiQLDelete(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, returnValuesOnConditionCheckFailure string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	release, lockErr := s.lockStandaloneWriteStatement(store, reqCtx, statement, params)
	if lockErr != nil {
		return nil, lockErr
	}
	defer release()

	var result *partiqlWriteResult
	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		var runErr error
		result, runErr = s.standaloneWritePlane(store, txn).runDelete(statement, params, returnValuesOnConditionCheckFailure)
		return runErr
	})
	if err != nil {
		return nil, err
	}
	if result != nil {
		s.emitChangePropagation(store, reqCtx.GetRegion(), result.table, dbstore.StreamEventRemove, result.keys, nil, result.oldImage, s.replicaDeleteOp(result.table, result.keys))
	}

	return buildPartiqlWriteResponse(result), nil
}

// resolvePrimaryKeyFromWhere extracts the equality bound to every
// primary-key attribute of the table from an UPDATE/DELETE WHERE clause.
// ok=false when the clause does not carry an equality on every key
// attribute: the documented contract requires the condition to "resolve to
// a single primary key value".
func resolvePrimaryKeyFromWhere(expr sqlparser.Expr, keySchema []*dbstore.KeySchemaElement, params *partiQLParams) (map[string]*dbstore.AttributeValue, bool) {
	if expr == nil {
		return nil, false
	}
	equalities := make(map[string]*dbstore.AttributeValue)
	collectWhereEqualities(expr, equalities, params)
	key := make(map[string]*dbstore.AttributeValue, len(keySchema))
	for _, ks := range keySchema {
		attr, ok := equalities[ks.AttributeName]
		if !ok {
			return nil, false
		}
		key[ks.AttributeName] = attr
	}
	return key, true
}

// collectWhereEqualities gathers the attribute-equality comparisons of a
// WHERE clause tree. Only top-level AND combinations contribute: an
// equality under OR does not identify the key on its own.
func collectWhereEqualities(expr sqlparser.Expr, out map[string]*dbstore.AttributeValue, params *partiQLParams) {
	switch e := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if e.Operator != sqlparser.EqualStr {
			return
		}
		if col, ok := e.Left.(*sqlparser.ColName); ok {
			if v := extractValueAttr(e.Right, params); v != nil {
				out[col.Name.String()] = v
			}
		}
	case *sqlparser.AndExpr:
		collectWhereEqualities(e.Left, out, params)
		collectWhereEqualities(e.Right, out, params)
	}
}

// matchKeyedTargetItem enforces the AWS single-item contract shared by every
// UPDATE and DELETE engine: the WHERE clause must resolve to a single
// primary-key value (an equality on every key attribute), and the statement
// addresses the item under that key. The getItem adapter runs the key
// lookup of the caller's choosing (store or transaction). The flags report
// whether the key resolved to a stored item and whether the WHERE clause
// evaluated true on it — the verb-specific engines turn the flag
// combinations into the documented outcomes (an absent key and a rejected
// condition are different results for DELETE).
func matchKeyedTargetItem(
	getItem func(key map[string]*dbstore.AttributeValue) (*dbstore.Item, error),
	table *dbstore.Table,
	whereExpr sqlparser.Expr,
	params *partiQLParams,
	stmtVerb string,
) (item *dbstore.Item, found, conditionHeld bool, err error) {
	key, ok := resolvePrimaryKeyFromWhere(whereExpr, table.KeySchema, params)
	if !ok {
		return nil, false, false, NewAPIError("com.amazon.coral.validate#ValidationException",
			stmtVerb+" statement WHERE clause must resolve to a single primary key value", http.StatusBadRequest)
	}
	if typeErr := validateKeyTypes(table, key); typeErr != nil {
		return nil, false, false, typeErr
	}

	item, getErr := getItem(key)
	if getErr != nil {
		if dbstore.IsItemNotFound(getErr) {
			return nil, false, false, nil
		}
		return nil, false, false, getErr
	}
	if !evaluateExpr(item.Attributes, whereExpr, params) {
		return item, true, false, nil
	}
	return item, true, true, nil
}

func buildKeyFromSchema(keySchema []*dbstore.KeySchemaElement, itemData map[string]*dbstore.AttributeValue) map[string]*dbstore.AttributeValue {
	key := make(map[string]*dbstore.AttributeValue)
	for _, ks := range keySchema {
		if attr, exists := itemData[ks.AttributeName]; exists {
			key[ks.AttributeName] = attr
		}
	}
	if len(key) < len(keySchema) {
		return nil
	}
	return key
}

// extractPartitionKeyFromWhere extracts the partition-key equality value
// from a WHERE clause and renders it with the store key encoding, so the
// partition scan prefix matches stored keys of every key type (a raw number
// literal or stringified parameter matches nothing on encoded keys).
func extractPartitionKeyFromWhere(expr sqlparser.Expr, pkName string, params *partiQLParams) string {
	if expr == nil || pkName == "" {
		return ""
	}

	if cmp, ok := expr.(*sqlparser.ComparisonExpr); ok {
		if col, ok := cmp.Left.(*sqlparser.ColName); ok {
			if col.Name.String() == pkName && cmp.Operator == sqlparser.EqualStr {
				return dbstore.EncodeKeyValue(extractValueAttr(cmp.Right, params))
			}
		}
	}

	if and, ok := expr.(*sqlparser.AndExpr); ok {
		if val := extractPartitionKeyFromWhere(and.Left, pkName, params); val != "" {
			return val
		}
		return extractPartitionKeyFromWhere(and.Right, pkName, params)
	}

	return ""
}

// extractValueAttr materialises a WHERE-clause literal or bound parameter as
// an AttributeValue via the shared value materialiser. An unresolvable
// expression yields nil, which never matches a partition key.
func extractValueAttr(expr sqlparser.Expr, params *partiQLParams) *dbstore.AttributeValue {
	switch expr.(type) {
	case *sqlparser.SQLVal, *sqlparser.ObjectLiteral, *sqlparser.ValTuple, *sqlparser.NullVal, *sqlparser.BoolVal:
		v, err := exprToAttributeValueWithParams(expr, params)
		if err != nil {
			return nil
		}
		return v
	}
	return nil
}

func extractTableNameFromStatement(statement string) string {
	upper := strings.ToUpper(statement)
	var rest string

	switch {
	case strings.HasPrefix(upper, "SELECT"):
		fromIdx := strings.Index(upper, " FROM ")
		if fromIdx == -1 {
			return ""
		}
		rest = strings.TrimSpace(statement[fromIdx+6:])
	case strings.HasPrefix(upper, "INSERT"):
		intoIdx := strings.Index(upper, " INTO ")
		if intoIdx == -1 {
			return ""
		}
		rest = strings.TrimSpace(statement[intoIdx+6:])
	case strings.HasPrefix(upper, "UPDATE"):
		rest = strings.TrimSpace(statement[6:])
	case strings.HasPrefix(upper, "DELETE"):
		fromIdx := strings.Index(upper, " FROM ")
		if fromIdx == -1 {
			return ""
		}
		rest = strings.TrimSpace(statement[fromIdx+6:])
	default:
		return ""
	}

	if strings.HasPrefix(rest, "\"") {
		endQuote := strings.Index(rest[1:], "\"")
		if endQuote != -1 {
			return rest[1 : endQuote+1]
		}
	}
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func sortItemsByOrderBy(items []*dbstore.Item, orderBy *orderByClause) []*dbstore.Item {
	if orderBy == nil || orderBy.column == "" || len(items) <= 1 {
		return items
	}

	sorted := make([]*dbstore.Item, len(items))
	copy(sorted, items)

	// The comparator counts type mismatches and equal values as ties that
	// "keep the stable order" (compareItemsByAttr's contract) — a promise
	// only a stable sort keeps.
	sort.SliceStable(sorted, func(i, j int) bool {
		return compareItemsByAttr(sorted[i], sorted[j], orderBy.column, orderBy.direction) < 0
	})
	return sorted
}

// compareItemsByAttr orders two items by one attribute using DynamoDB
// ordering: same-type S/N/B compare (numbers numerically); any other pair
// — including type mismatches — is unordered and keeps the stable order.
// Items missing the attribute sort first.
func compareItemsByAttr(a, b *dbstore.Item, attrName, direction string) int {
	aVal, aOk := a.Attributes[attrName]
	bVal, bOk := b.Attributes[attrName]

	if !aOk && !bOk {
		return 0
	}
	if !aOk {
		return -1
	}
	if !bOk {
		return 1
	}

	cmp, ok := compareOrderedValues(aVal, bVal)
	if !ok {
		return 0
	}

	if direction == "DESC" {
		return -cmp
	}
	return cmp
}
