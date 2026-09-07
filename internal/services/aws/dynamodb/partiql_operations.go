package dynamodb

import (
	"context"
	"errors"
	"sort"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

// errScanSufficient is returned from a scan callback to signal that
// enough items have been collected and the scan can stop early.
var errScanSufficient = errors.New("scan sufficient items collected")

// ExecuteStatement executes a PartiQL statement.
func (s *DynamoDBService) ExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statement := request.GetStringParam(req.Parameters, "Statement")
	if !validatePartiQLStatement(statement) {
		return nil, ErrInvalidParameter
	}

	params := parsePartiQLParams(req.Parameters)
	// Normalise `?` placeholders into their explicit whole-statement :vN
	// form before dispatch: the engines' parses (including the manual
	// UPDATE path's per-segment re-parses) then share one numbering, and
	// the count check rejects a parameter list that does not carry exactly
	// one value per placeholder.
	statement, placeholderCount := preparePartiQLStatement(statement)
	if err := validateParameterCount(placeholderCount, params); err != nil {
		return nil, err
	}
	consistentRead := request.GetBoolParam(req.Parameters, "ConsistentRead")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit > 0 {
		if !validateExecuteStatementLimit(limit) {
			return nil, ErrInvalidParameter
		}
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	returnValuesOnConditionCheckFailure := request.GetStringParam(req.Parameters, "ReturnValuesOnConditionCheckFailure")

	upperStmt := strings.ToUpper(strings.TrimSpace(statement))
	var result interface{}
	var err error

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

// applySetAssignments writes SET clause assignments into attrs. Values are
// literals or bound parameters resolved through the shared materialiser;
// an unresolvable placeholder fails the statement.
func applySetAssignments(attrs map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) error {
	for _, asgn := range assignments {
		value := asgn.value

		// if_not_exists(attr, fallback) keeps the existing attribute when
		// the referenced one is already present.
		if funcExpr, ok := value.(*sqlparser.FuncExpr); ok && strings.EqualFold(funcExpr.Name.String(), "if_not_exists") && len(funcExpr.Exprs) >= 2 {
			if aliased, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr); ok {
				if colName, ok := aliased.Expr.(*sqlparser.ColName); ok {
					if existing, exists := attrs[colName.Name.String()]; exists && existing != nil {
						continue
					}
				}
			}
			if aliased, ok := funcExpr.Exprs[1].(*sqlparser.AliasedExpr); ok {
				value = aliased.Expr
			}
		}

		attrValue, err := exprToAttributeValueWithParams(value, params)
		if err != nil {
			return err
		}
		attrs[asgn.attrName] = attrValue
	}
	return nil
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

	sort.Slice(sorted, func(i, j int) bool {
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

// applyRemoveAttrs deletes the named top-level attributes from the item.
func applyRemoveAttrs(attrs map[string]*dbstore.AttributeValue, removeAttrs []string) {
	for _, name := range removeAttrs {
		delete(attrs, name)
	}
}

// updateClauseTargetNames collects every top-level attribute name an UPDATE
// statement writes across its SET, REMOVE, ADD and DELETE clauses. The key
// attributes identify the item being updated, so any clause touching one of
// them is rejected before the statement applies.
func updateClauseTargetNames(clauses updateClauses) []string {
	names := make([]string, 0, len(clauses.setAssignments)+len(clauses.removeAttrs)+
		len(clauses.addAssignments)+len(clauses.deleteAssignments))
	for _, asgn := range clauses.setAssignments {
		names = append(names, asgn.attrName)
	}
	names = append(names, clauses.removeAttrs...)
	for _, asgn := range clauses.addAssignments {
		names = append(names, asgn.attrName)
	}
	for _, asgn := range clauses.deleteAssignments {
		names = append(names, asgn.attrName)
	}
	return names
}
