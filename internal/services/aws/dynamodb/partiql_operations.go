package dynamodb

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

// ExecuteStatement executes a PartiQL statement.
func (s *DynamoDBService) ExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statement := request.GetStringParam(req.Parameters, "Statement")
	if statement == "" {
		return nil, ErrInvalidParameter
	}

	params := parsePartiQLParams(req.Parameters)
	consistentRead := request.GetBoolParam(req.Parameters, "ConsistentRead")
	limit := request.GetIntParam(req.Parameters, "Limit")

	upperStmt := strings.ToUpper(strings.TrimSpace(statement))
	var result interface{}
	var err error

	switch {
	case strings.HasPrefix(upperStmt, "SELECT"):
		result, err = s.executePartiQLSelectEnhanced(ctx, reqCtx, statement, params, consistentRead, limit)
	case strings.HasPrefix(upperStmt, "INSERT"):
		result, err = s.executePartiQLInsert(ctx, reqCtx, statement, params)
	case strings.HasPrefix(upperStmt, "UPDATE"):
		result, err = s.executePartiQLUpdate(ctx, reqCtx, statement, params)
	case strings.HasPrefix(upperStmt, "DELETE"):
		result, err = s.executePartiQLDelete(ctx, reqCtx, statement, params)
	default:
		return nil, ErrInvalidParameter
	}

	if err != nil {
		return nil, err
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableName := extractTableNameFromStatement(statement)
		capacityUnits := 0.5
		if strings.HasPrefix(upperStmt, "INSERT") || strings.HasPrefix(upperStmt, "UPDATE") || strings.HasPrefix(upperStmt, "DELETE") {
			capacityUnits = 1.0
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			resultMap["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, capacityUnits)
		}
	}

	return result, nil
}

func (s *DynamoDBService) executePartiQLSelectEnhanced(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, consistentRead bool, limit int) (interface{}, error) {
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

	var items []*dbstore.Item

	pkName := getPartitionKeyName(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	if pkValue != "" {
		err = store.Items().ScanByPartitionKey(tableName, pkValue, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
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

	if orderBy != nil {
		items = sortItemsByOrderBy(items, orderBy)
	}

	if limit > 0 && len(items) > limit {
		items = items[:limit]
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

	return map[string]interface{}{
		"Items":        result,
		"Count":        len(result),
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLInsert(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams) (interface{}, error) {
	tableName, itemData := parseInsertStatementWithParams(statement, params)
	if tableName == "" || itemData == nil {
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

	keyAttrs := buildKeyFromSchema(table.KeySchema, itemData)
	if keyAttrs == nil {
		return nil, ErrInvalidParameter
	}

	var isNewItem bool

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		_, err := txn.GetItem(tableName, keyAttrs)
		if err != nil {
			if !dbstore.IsItemNotFound(err) {
				return err
			}
			isNewItem = true
		} else {
			return ErrConditionalCheckFailed
		}

		if err := txn.PutItem(tableName, keyAttrs, itemData); err != nil {
			return err
		}

		newItem := &dbstore.Item{
			TableName:  tableName,
			Key:        keyAttrs,
			Attributes: itemData,
		}
		if err := txn.PutIndexEntries(tableName, newItem); err != nil {
			return err
		}

		if isNewItem {
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			newItemSize := calculateItemSize(itemData)
			if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *DynamoDBService) executePartiQLUpdate(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams) (ret interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in executePartiQLUpdate: %v", r)
		}
	}()

	tableName, assignments, whereExpr := parseUpdateStatement(statement)
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

	var items []*dbstore.Item
	err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, err
	}

	scannedCount := len(items)
	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	updatedCount := 0
	for _, item := range items {
		oldSize := calculateItemSize(item.Attributes)

		oldItem := &dbstore.Item{
			TableName:  tableName,
			Key:        copyAttributes(item.Key),
			Attributes: copyAttributes(item.Attributes),
		}

		applySetAssignments(item.Attributes, assignments, params)

		newSize := calculateItemSize(item.Attributes)
		sizeDelta := newSize - oldSize

		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.DeleteIndexEntries(tableName, oldItem); err != nil {
				return err
			}

			if err := txn.PutItem(tableName, item.Key, item.Attributes); err != nil {
				return err
			}

			updatedItem := &dbstore.Item{
				TableName:  tableName,
				Key:        item.Key,
				Attributes: item.Attributes,
			}
			if err := txn.PutIndexEntries(tableName, updatedItem); err != nil {
				return err
			}

			if sizeDelta != 0 {
				if err := txn.UpdateTableSize(tableName, sizeDelta); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
		updatedCount++
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        updatedCount,
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLDelete(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams) (ret interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in executePartiQLDelete: %v", r)
		}
	}()

	tableName, whereExpr := parseDeleteStatement(statement)
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

	var items []*dbstore.Item
	err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, err
	}

	scannedCount := len(items)
	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	deletedCount := 0
	for _, item := range items {
		itemSize := calculateItemSize(item.Attributes)
		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.DeleteIndexEntries(tableName, item); err != nil {
				return err
			}

			if err := txn.DeleteItem(tableName, item.Key); err != nil {
				return err
			}

			if err := txn.UpdateItemCount(tableName, -1); err != nil {
				return err
			}

			if err := txn.UpdateTableSize(tableName, -itemSize); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
		deletedCount++
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        deletedCount,
		"ScannedCount": scannedCount,
	}, nil
}

func applySetAssignments(attrs map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) {
	for _, asgn := range assignments {
		var attrValue *dbstore.AttributeValue
		switch e := asgn.value.(type) {
		case *sqlparser.SQLVal:
			if e.Type == sqlparser.ValArg {
				if strings.HasPrefix(string(e.Val), ":") {
					idxStr := strings.TrimPrefix(string(e.Val), ":v")
					if idx, err := strconv.Atoi(idxStr); err == nil && params != nil && idx > 0 && idx <= len(params.Parameters) {
						attrValue = paramToAttributeValue(params.Parameters[idx-1])
					}
				}
			} else {
				attrValue = exprToAttributeValue(e)
			}
		default:
			attrValue = exprToAttributeValue(asgn.value)
		}

		if attrValue != nil {
			attrs[asgn.attrName] = attrValue
		}
	}
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

func getPartitionKeyName(table *dbstore.Table) string {
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			return ks.AttributeName
		}
	}
	return ""
}

func extractPartitionKeyFromWhere(expr sqlparser.Expr, pkName string, params *partiQLParams) string {
	if expr == nil || pkName == "" {
		return ""
	}

	if cmp, ok := expr.(*sqlparser.ComparisonExpr); ok {
		if col, ok := cmp.Left.(*sqlparser.ColName); ok {
			if col.Name.String() == pkName && cmp.Operator == sqlparser.EqualStr {
				return extractValueString(cmp.Right, params)
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

func extractValueString(expr sqlparser.Expr, params *partiQLParams) string {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		switch e.Type {
		case sqlparser.StrVal:
			return string(e.Val)
		case sqlparser.IntVal, sqlparser.FloatVal:
			return string(e.Val)
		case sqlparser.ValArg:
			if strings.HasPrefix(string(e.Val), ":") {
				idxStr := strings.TrimPrefix(string(e.Val), ":v")
				if idx, err := strconv.Atoi(idxStr); err == nil && params != nil && idx > 0 && idx <= len(params.Parameters) {
					return paramToString(params.Parameters[idx-1])
				}
			}
		}
	}
	return ""
}

func ptrBool(b bool) *bool {
	return &b
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

	sort.Slice(sorted, func(i, j int) bool {
		return compareItemsByAttr(sorted[i], sorted[j], orderBy.column, orderBy.direction) < 0
	})
	return sorted
}

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

	aStr := attrValueToCompareString(aVal)
	bStr := attrValueToCompareString(bVal)

	cmp := compareValues(aStr, bStr)

	if direction == "DESC" {
		return -cmp
	}
	return cmp
}

func attrValueToCompareString(attr *dbstore.AttributeValue) string {
	if attr == nil {
		return ""
	}
	if attr.S != nil {
		return *attr.S
	}
	if attr.N != nil {
		return *attr.N
	}
	if attr.BOOL != nil {
		if *attr.BOOL {
			return "true"
		}
		return "false"
	}
	return ""
}
