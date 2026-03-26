package athena

import (
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

type tableResolver struct {
	catalog  string
	database string
	service  *Service
	ctx      *request.RequestContext
}

func (s *Service) executeAdvancedSelect(reqCtx *request.RequestContext, selectStmt *sqlparser.Select, context *athenastore.QueryExecutionContext) (*athenastore.StoredTable, error) {
	catalog := "AwsDataCatalog"
	database := "default"
	if context != nil {
		if context.Catalog != "" {
			catalog = context.Catalog
		}
		if context.Database != "" {
			database = context.Database
		}
	}

	resolver := &tableResolver{
		catalog:  catalog,
		database: database,
		service:  s,
		ctx:      reqCtx,
	}

	tables, err := resolver.resolveTables(selectStmt.From)
	if err != nil {
		return nil, err
	}

	var rows []*athenastore.StoredRow
	if len(tables) == 1 {
		for _, row := range tables[0].Data.Rows {
			if selectStmt.Where != nil {
				if !s.evaluateWhere(selectStmt.Where.Expr, row.Values) {
					continue
				}
			}
			rows = append(rows, row)
		}
	} else if len(tables) > 1 {
		rows, err = s.executeJoin(selectStmt, tables)
		if err != nil {
			return nil, err
		}
	}

	if len(selectStmt.GroupBy) > 0 || s.hasAggregateFunctions(selectStmt.SelectExprs) {
		rows = s.executeGroupBy(selectStmt, rows)
	} else {
		rows = s.projectColumns(rows, selectStmt.SelectExprs)
	}

	if selectStmt.Having != nil {
		rows = s.applyHaving(rows, selectStmt.Having.Expr)
	}

	if len(selectStmt.OrderBy) > 0 {
		s.applyOrderBy(rows, selectStmt.OrderBy)
	}

	if selectStmt.Limit != nil {
		rows = s.applyLimit(rows, selectStmt.Limit)
	}

	var columns []athenastore.Column
	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			colName := s.extractColumnName(aliased.Expr)
			outputName := colName
			if !aliased.As.IsEmpty() {
				outputName = aliased.As.String()
			}
			columns = append(columns, athenastore.Column{
				Name: outputName,
				Type: "string",
			})
		}
	}

	return &athenastore.StoredTable{
		Columns: columns,
		Rows:    rows,
	}, nil
}

type resolvedTable struct {
	Name string
	Data *athenastore.StoredTable
}

func (r *tableResolver) resolveTables(from sqlparser.TableExprs) ([]resolvedTable, error) {
	var tables []resolvedTable

	for _, tableExpr := range from {
		switch expr := tableExpr.(type) {
		case *sqlparser.AliasedTableExpr:
			tableName, err := r.resolveTableName(expr.Expr)
			if err != nil {
				return nil, err
			}

			data, err := r.service.getTableData(r.ctx, r.catalog, r.database, tableName)
			if err != nil {
				data = &athenastore.StoredTable{Columns: []athenastore.Column{}, Rows: []*athenastore.StoredRow{}}
			}

			alias := tableName
			if !expr.As.IsEmpty() {
				alias = expr.As.String()
			}

			tables = append(tables, resolvedTable{
				Name: alias,
				Data: data,
			})

		case *sqlparser.JoinTableExpr:
			leftTables, err := r.resolveTables([]sqlparser.TableExpr{expr.LeftExpr})
			if err != nil {
				return nil, err
			}
			rightTables, err := r.resolveTables([]sqlparser.TableExpr{expr.RightExpr})
			if err != nil {
				return nil, err
			}
			tables = append(tables, leftTables...)
			tables = append(tables, rightTables...)
		}
	}

	return tables, nil
}

func (r *tableResolver) resolveTableName(expr sqlparser.SimpleTableExpr) (string, error) {
	switch e := expr.(type) {
	case sqlparser.TableName:
		return e.Name.String(), nil
	default:
		return "", fmt.Errorf("unsupported table expression type: %T", expr)
	}
}

type joinInfo struct {
	left     *joinInfo
	right    *joinInfo
	table    *resolvedTable
	joinType string
	onExpr   sqlparser.Expr
}

func (s *Service) executeJoin(selectStmt *sqlparser.Select, tables []resolvedTable) ([]*athenastore.StoredRow, error) {
	if len(tables) < 2 {
		return tables[0].Data.Rows, nil
	}

	joinTree := s.buildJoinTree(selectStmt.From, tables)
	if joinTree == nil {
		return nil, fmt.Errorf("failed to build join tree")
	}

	result := s.executeJoinTree(joinTree, selectStmt.Where)

	return result, nil
}

func (s *Service) buildJoinTree(from sqlparser.TableExprs, tables []resolvedTable) *joinInfo {
	if len(from) == 0 || len(tables) == 0 {
		return nil
	}

	tableIdx := 0

	var buildTree func(sqlparser.TableExpr) *joinInfo
	buildTree = func(expr sqlparser.TableExpr) *joinInfo {
		switch e := expr.(type) {
		case *sqlparser.AliasedTableExpr:
			if tableIdx < len(tables) {
				info := &joinInfo{table: &tables[tableIdx]}
				tableIdx++
				return info
			}
		case *sqlparser.JoinTableExpr:
			left := buildTree(e.LeftExpr)
			right := buildTree(e.RightExpr)
			if left != nil && right != nil {
				return &joinInfo{
					left:     left,
					right:    right,
					joinType: e.Join,
					onExpr:   e.Condition.On,
				}
			}
		case *sqlparser.ParenTableExpr:
			if len(e.Exprs) > 0 {
				return buildTree(e.Exprs[0])
			}
		}
		return nil
	}

	if len(from) == 1 {
		return buildTree(from[0])
	}

	var result *joinInfo
	for _, expr := range from {
		next := buildTree(expr)
		if result == nil {
			result = next
		} else if next != nil {
			result = &joinInfo{
				left:     result,
				right:    next,
				joinType: sqlparser.JoinStr,
			}
		}
	}
	return result
}

func (s *Service) executeJoinTree(ji *joinInfo, where *sqlparser.Where) []*athenastore.StoredRow {
	if ji == nil {
		return nil
	}

	if ji.table != nil {
		return ji.table.Data.Rows
	}

	if ji.left == nil || ji.right == nil {
		return nil
	}

	leftRows := s.executeJoinTree(ji.left, nil)
	rightRows := s.executeJoinTree(ji.right, nil)

	leftName := s.getJoinTableName(ji.left)
	rightName := s.getJoinTableName(ji.right)

	var result []*athenastore.StoredRow

	switch ji.joinType {
	case sqlparser.LeftJoinStr, sqlparser.NaturalLeftJoinStr:
		result = s.leftJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
	case sqlparser.RightJoinStr, sqlparser.NaturalRightJoinStr:
		result = s.rightJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
	default:
		result = s.innerJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
	}

	if where != nil {
		var filtered []*athenastore.StoredRow
		for _, row := range result {
			if s.evaluateWhere(where.Expr, row.Values) {
				filtered = append(filtered, row)
			}
		}
		result = filtered
	}

	return result
}

func (s *Service) getJoinTableName(ji *joinInfo) string {
	if ji == nil {
		return ""
	}
	if ji.table != nil {
		return ji.table.Name
	}
	return "joined"
}

func (s *Service) innerJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) []*athenastore.StoredRow {
	var result []*athenastore.StoredRow

	for _, leftRow := range leftRows {
		for _, rightRow := range rightRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil || s.evaluateWhere(onExpr, combined.Values) {
				result = append(result, combined)
			}
		}
	}

	return result
}

func (s *Service) leftJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) []*athenastore.StoredRow {
	var result []*athenastore.StoredRow

	for _, leftRow := range leftRows {
		matched := false
		for _, rightRow := range rightRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil || s.evaluateWhere(onExpr, combined.Values) {
				result = append(result, combined)
				matched = true
			}
		}
		if !matched {
			combined := s.mergeRowsWithNullRight(leftRow, rightName)
			result = append(result, combined)
		}
	}

	return result
}

func (s *Service) rightJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) []*athenastore.StoredRow {
	var result []*athenastore.StoredRow

	for _, rightRow := range rightRows {
		matched := false
		for _, leftRow := range leftRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil || s.evaluateWhere(onExpr, combined.Values) {
				result = append(result, combined)
				matched = true
			}
		}
		if !matched {
			combined := s.mergeRowsWithNullLeft(rightRow, leftName)
			result = append(result, combined)
		}
	}

	return result
}

func (s *Service) mergeRows(leftRow, rightRow *athenastore.StoredRow, leftName, rightName string) *athenastore.StoredRow {
	combined := &athenastore.StoredRow{Values: make(map[string]interface{})}

	for k, v := range leftRow.Values {
		combined.Values[leftName+"."+k] = v
		combined.Values[k] = v
	}
	for k, v := range rightRow.Values {
		combined.Values[rightName+"."+k] = v
		if _, exists := combined.Values[k]; !exists {
			combined.Values[k] = v
		}
	}

	return combined
}

func (s *Service) mergeRowsWithNullRight(leftRow *athenastore.StoredRow, rightName string) *athenastore.StoredRow {
	combined := &athenastore.StoredRow{Values: make(map[string]interface{})}

	for k, v := range leftRow.Values {
		combined.Values[k] = v
	}

	return combined
}

func (s *Service) mergeRowsWithNullLeft(rightRow *athenastore.StoredRow, leftName string) *athenastore.StoredRow {
	combined := &athenastore.StoredRow{Values: make(map[string]interface{})}

	for k, v := range rightRow.Values {
		combined.Values[k] = v
	}

	return combined
}

func (s *Service) hasAggregateFunctions(exprs sqlparser.SelectExprs) bool {
	for _, expr := range exprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			if fn, ok := aliased.Expr.(*sqlparser.FuncExpr); ok {
				if fn.IsAggregate() {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) executeGroupBy(selectStmt *sqlparser.Select, rows []*athenastore.StoredRow) []*athenastore.StoredRow {
	if len(selectStmt.GroupBy) == 0 {
		return s.executeAggregateWithoutGroup(selectStmt, rows)
	}

	groups := make(map[string][]*athenastore.StoredRow)

	for _, row := range rows {
		key := s.buildGroupKey(selectStmt.GroupBy, row.Values)
		groups[key] = append(groups[key], row)
	}

	var result []*athenastore.StoredRow

	for _, groupRows := range groups {
		aggregatedRow := s.aggregateGroup(selectStmt, groupRows)
		result = append(result, aggregatedRow)
	}

	return result
}

func (s *Service) buildGroupKey(groupBy sqlparser.GroupBy, row map[string]interface{}) string {
	var keyParts []string
	for _, expr := range groupBy {
		colName := s.extractColumnName(expr)
		if val, ok := row[colName]; ok {
			keyParts = append(keyParts, fmt.Sprintf("%v", val))
		} else {
			keyParts = append(keyParts, "")
		}
	}
	return strings.Join(keyParts, "|")
}

func (s *Service) executeAggregateWithoutGroup(selectStmt *sqlparser.Select, rows []*athenastore.StoredRow) []*athenastore.StoredRow {
	aggregatedRow := s.aggregateGroup(selectStmt, rows)
	return []*athenastore.StoredRow{aggregatedRow}
}

func (s *Service) aggregateGroup(selectStmt *sqlparser.Select, rows []*athenastore.StoredRow) *athenastore.StoredRow {
	result := &athenastore.StoredRow{Values: make(map[string]interface{})}

	if len(rows) > 0 {
		for k, v := range rows[0].Values {
			result.Values[k] = v
		}
	}

	for _, expr := range selectStmt.SelectExprs {
		if aliased, ok := expr.(*sqlparser.AliasedExpr); ok {
			outputName := s.extractColumnName(aliased.Expr)
			if !aliased.As.IsEmpty() {
				outputName = aliased.As.String()
			}

			if fn, ok := aliased.Expr.(*sqlparser.FuncExpr); ok {
				if fn.IsAggregate() {
					result.Values[outputName] = s.computeAggregate(fn, rows)
				}
			}
		}
	}

	return result
}

func (s *Service) computeAggregate(fn *sqlparser.FuncExpr, rows []*athenastore.StoredRow) interface{} {
	funcName := strings.ToUpper(fn.Name.String())

	switch funcName {
	case "COUNT":
		if len(fn.Exprs) == 1 {
			if _, isStar := fn.Exprs[0].(*sqlparser.StarExpr); isStar {
				return len(rows)
			}
		}

		colName := ""
		if len(fn.Exprs) > 0 {
			if aliased, ok := fn.Exprs[0].(*sqlparser.AliasedExpr); ok {
				colName = s.extractColumnName(aliased.Expr)
			}
		}

		if colName != "" {
			count := 0
			for _, row := range rows {
				if val, ok := row.Values[colName]; ok && val != nil {
					count++
				}
			}
			return count
		}
		return len(rows)

	case "SUM":
		return s.aggregateNumeric(fn, rows, "sum")

	case "AVG":
		return s.aggregateNumeric(fn, rows, "avg")

	case "MIN":
		return s.aggregateNumeric(fn, rows, "min")

	case "MAX":
		return s.aggregateNumeric(fn, rows, "max")

	default:
		return nil
	}
}

func (s *Service) aggregateNumeric(fn *sqlparser.FuncExpr, rows []*athenastore.StoredRow, op string) interface{} {
	if len(fn.Exprs) == 0 {
		return nil
	}

	colName := ""
	if aliased, ok := fn.Exprs[0].(*sqlparser.AliasedExpr); ok {
		colName = s.extractColumnName(aliased.Expr)
	}

	if colName == "" {
		return nil
	}

	var values []float64
	for _, row := range rows {
		if val, ok := row.Values[colName]; ok && val != nil {
			if num, err := s.toFloat(val); err == nil {
				values = append(values, num)
			}
		}
	}

	if len(values) == 0 {
		return nil
	}

	switch op {
	case "sum":
		var sum float64
		for _, v := range values {
			sum += v
		}
		return sum
	case "avg":
		var sum float64
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values))
	case "min":
		min := values[0]
		for _, v := range values[1:] {
			if v < min {
				min = v
			}
		}
		return min
	case "max":
		max := values[0]
		for _, v := range values[1:] {
			if v > max {
				max = v
			}
		}
		return max
	}

	return nil
}

func (s *Service) toFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return strconv.ParseFloat(fmt.Sprintf("%v", val), 64)
	}
}

func (s *Service) applyHaving(rows []*athenastore.StoredRow, expr sqlparser.Expr) []*athenastore.StoredRow {
	var result []*athenastore.StoredRow

	for _, row := range rows {
		if s.evaluateHaving(expr, row.Values) {
			result = append(result, row)
		}
	}

	return result
}

func (s *Service) evaluateHaving(expr sqlparser.Expr, row map[string]interface{}) bool {
	switch e := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return s.evaluateHavingComparison(e, row)
	case *sqlparser.AndExpr:
		return s.evaluateHaving(e.Left, row) && s.evaluateHaving(e.Right, row)
	case *sqlparser.OrExpr:
		return s.evaluateHaving(e.Left, row) || s.evaluateHaving(e.Right, row)
	case *sqlparser.ParenExpr:
		return s.evaluateHaving(e.Expr, row)
	}
	return true
}

func (s *Service) evaluateHavingComparison(expr *sqlparser.ComparisonExpr, row map[string]interface{}) bool {
	leftVal := s.getHavingExprValue(expr.Left, row)
	rightVal := s.getHavingExprValue(expr.Right, row)

	leftStr := fmt.Sprintf("%v", leftVal)
	rightStr := fmt.Sprintf("%v", rightVal)

	switch expr.Operator {
	case sqlparser.EqualStr:
		return leftStr == rightStr
	case sqlparser.NotEqualStr:
		return leftStr != rightStr
	case sqlparser.LessThanStr:
		return s.compareNumeric(leftVal, rightVal) < 0
	case sqlparser.LessEqualStr:
		return s.compareNumeric(leftVal, rightVal) <= 0
	case sqlparser.GreaterThanStr:
		return s.compareNumeric(leftVal, rightVal) > 0
	case sqlparser.GreaterEqualStr:
		return s.compareNumeric(leftVal, rightVal) >= 0
	}
	return false
}

func (s *Service) getHavingExprValue(expr sqlparser.Expr, row map[string]interface{}) interface{} {
	switch e := expr.(type) {
	case *sqlparser.ColName:
		colName := e.Name.String()
		return row[colName]
	case *sqlparser.SQLVal:
		if e.Type == sqlparser.StrVal {
			return string(e.Val)
		} else if e.Type == sqlparser.IntVal {
			if val, err := strconv.ParseInt(string(e.Val), 10, 64); err == nil {
				return val
			}
		} else if e.Type == sqlparser.FloatVal {
			if val, err := strconv.ParseFloat(string(e.Val), 64); err == nil {
				return val
			}
		}
		return string(e.Val)
	case *sqlparser.NullVal:
		return nil
	default:
		return nil
	}
}

func (s *Service) compareNumeric(left, right interface{}) int {
	leftFloat, leftErr := s.toFloat(left)
	rightFloat, rightErr := s.toFloat(right)

	if leftErr != nil || rightErr != nil {
		leftStr := fmt.Sprintf("%v", left)
		rightStr := fmt.Sprintf("%v", right)
		if leftStr < rightStr {
			return -1
		} else if leftStr > rightStr {
			return 1
		}
		return 0
	}

	if leftFloat < rightFloat {
		return -1
	} else if leftFloat > rightFloat {
		return 1
	}
	return 0
}
