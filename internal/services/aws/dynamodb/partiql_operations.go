package dynamodb

import (
	"context"
	"errors"
	"sort"
	"strconv"
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

func applySetAssignments(attrs map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) {
	for _, asgn := range assignments {
		var attrValue *dbstore.AttributeValue

		if funcExpr, ok := asgn.value.(*sqlparser.FuncExpr); ok && strings.EqualFold(funcExpr.Name.String(), "if_not_exists") {
			if len(funcExpr.Exprs) >= 2 {
				if aliased, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr); ok {
					if colName, ok := aliased.Expr.(*sqlparser.ColName); ok {
						attrName := colName.Name.String()
						if existing, exists := attrs[attrName]; exists && existing != nil {
							continue
						}
					}
				}
				if aliased, ok := funcExpr.Exprs[1].(*sqlparser.AliasedExpr); ok {
					if sqlVal, ok := aliased.Expr.(*sqlparser.SQLVal); ok {
						attrValue = exprToAttributeValueWithParams(sqlVal, params)
					}
				}
			}
		} else if e, ok := asgn.value.(*sqlparser.SQLVal); ok {
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
		} else {
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

// applyRemoveAttrs deletes the named top-level attributes from the item.
func applyRemoveAttrs(attrs map[string]*dbstore.AttributeValue, removeAttrs []string) {
	for _, name := range removeAttrs {
		delete(attrs, name)
	}
}
