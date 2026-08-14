package athena

import (
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

type tableResolver struct {
	catalog  string
	database string
	service  *AthenaService
	ctx      *request.RequestContext
}

func (s *AthenaService) executeAdvancedSelect(reqCtx *request.RequestContext, selectStmt *sqlparser.Select, context *athenastore.QueryExecutionContext) (*athenastore.StoredTable, error) {
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
				match, err := s.evaluateWhere(selectStmt.Where.Expr, row.Values)
				if err != nil {
					return nil, err
				}
				if !match {
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

	hasAggregates := len(selectStmt.GroupBy) > 0 || s.hasAggregateFunctions(selectStmt.SelectExprs)
	if hasAggregates {
		rows = s.executeGroupBy(selectStmt, rows)
		if selectStmt.Having != nil {
			filtered, err := s.applyHaving(rows, selectStmt.Having.Expr)
			if err != nil {
				return nil, err
			}
			rows = filtered
		}
	} else {
		if selectStmt.Having != nil {
			filtered, err := s.applyHaving(rows, selectStmt.Having.Expr)
			if err != nil {
				return nil, err
			}
			rows = filtered
		}
		rows = s.projectColumns(rows, selectStmt.SelectExprs)
	}

	if len(selectStmt.OrderBy) > 0 {
		s.applyOrderBy(rows, selectStmt.OrderBy)
	}

	if selectStmt.Limit != nil {
		limited, err := s.applyLimit(rows, selectStmt.Limit)
		if err != nil {
			return nil, err
		}
		rows = limited
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
				return nil, fmt.Errorf("table %q not found: %w", tableName, err)
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

func (s *AthenaService) executeJoin(selectStmt *sqlparser.Select, tables []resolvedTable) ([]*athenastore.StoredRow, error) {
	if len(tables) < 2 {
		return tables[0].Data.Rows, nil
	}

	joinTree := s.buildJoinTree(selectStmt.From, tables)
	if joinTree == nil {
		return nil, fmt.Errorf("failed to build join tree")
	}

	result, err := s.executeJoinTree(joinTree, selectStmt.Where)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AthenaService) buildJoinTree(from sqlparser.TableExprs, tables []resolvedTable) *joinInfo {
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

func (s *AthenaService) executeJoinTree(ji *joinInfo, where *sqlparser.Where) ([]*athenastore.StoredRow, error) {
	if ji == nil {
		return nil, nil
	}

	if ji.table != nil {
		return ji.table.Data.Rows, nil
	}

	if ji.left == nil || ji.right == nil {
		return nil, nil
	}

	leftRows, err := s.executeJoinTree(ji.left, nil)
	if err != nil {
		return nil, err
	}
	rightRows, err := s.executeJoinTree(ji.right, nil)
	if err != nil {
		return nil, err
	}

	leftName := s.getJoinTableName(ji.left)
	rightName := s.getJoinTableName(ji.right)

	var result []*athenastore.StoredRow

	switch ji.joinType {
	case sqlparser.LeftJoinStr, sqlparser.NaturalLeftJoinStr:
		result, err = s.leftJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
		if err != nil {
			return nil, err
		}
	case sqlparser.RightJoinStr, sqlparser.NaturalRightJoinStr:
		result, err = s.rightJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
		if err != nil {
			return nil, err
		}
	case sqlparser.FullJoinStr, sqlparser.NaturalFullJoinStr:
		result, err = s.fullOuterJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
		if err != nil {
			return nil, err
		}
	default:
		result, err = s.innerJoin(leftRows, rightRows, leftName, rightName, ji.onExpr)
		if err != nil {
			return nil, err
		}
	}

	if where != nil {
		var filtered []*athenastore.StoredRow
		for _, row := range result {
			match, wErr := s.evaluateWhere(where.Expr, row.Values)
			if wErr != nil {
				return nil, wErr
			}
			if match {
				filtered = append(filtered, row)
			}
		}
		result = filtered
	}

	return result, nil
}

func (s *AthenaService) getJoinTableName(ji *joinInfo) string {
	if ji == nil {
		return ""
	}
	if ji.table != nil {
		return ji.table.Name
	}
	return "joined"
}

func (s *AthenaService) innerJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) ([]*athenastore.StoredRow, error) {
	var result []*athenastore.StoredRow

	for _, leftRow := range leftRows {
		for _, rightRow := range rightRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil {
				result = append(result, combined)
			} else {
				match, err := s.evaluateWhere(onExpr, combined.Values)
				if err != nil {
					return nil, err
				}
				if match {
					result = append(result, combined)
				}
			}
		}
	}

	return result, nil
}

func (s *AthenaService) leftJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) ([]*athenastore.StoredRow, error) {
	var result []*athenastore.StoredRow

	var rightCols []string
	if len(rightRows) > 0 {
		for k := range rightRows[0].Values {
			rightCols = append(rightCols, k)
		}
	}

	for _, leftRow := range leftRows {
		matched := false
		for _, rightRow := range rightRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil {
				result = append(result, combined)
				matched = true
			} else {
				match, err := s.evaluateWhere(onExpr, combined.Values)
				if err != nil {
					return nil, err
				}
				if match {
					result = append(result, combined)
					matched = true
				}
			}
		}
		if !matched {
			combined := s.mergeRowsWithNullRight(leftRow, leftName, rightName, rightCols)
			result = append(result, combined)
		}
	}

	return result, nil
}

func (s *AthenaService) rightJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) ([]*athenastore.StoredRow, error) {
	var result []*athenastore.StoredRow

	var leftCols []string
	if len(leftRows) > 0 {
		for k := range leftRows[0].Values {
			leftCols = append(leftCols, k)
		}
	}

	for _, rightRow := range rightRows {
		matched := false
		for _, leftRow := range leftRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil {
				result = append(result, combined)
				matched = true
			} else {
				match, err := s.evaluateWhere(onExpr, combined.Values)
				if err != nil {
					return nil, err
				}
				if match {
					result = append(result, combined)
					matched = true
				}
			}
		}
		if !matched {
			combined := s.mergeRowsWithNullLeft(rightRow, leftName, rightName, leftCols)
			result = append(result, combined)
		}
	}

	return result, nil
}

func (s *AthenaService) fullOuterJoin(leftRows, rightRows []*athenastore.StoredRow, leftName, rightName string, onExpr sqlparser.Expr) ([]*athenastore.StoredRow, error) {
	var result []*athenastore.StoredRow

	var leftCols, rightCols []string
	if len(leftRows) > 0 {
		for k := range leftRows[0].Values {
			leftCols = append(leftCols, k)
		}
	}
	if len(rightRows) > 0 {
		for k := range rightRows[0].Values {
			rightCols = append(rightCols, k)
		}
	}

	rightMatched := make([]bool, len(rightRows))

	for _, leftRow := range leftRows {
		matched := false
		for i, rightRow := range rightRows {
			combined := s.mergeRows(leftRow, rightRow, leftName, rightName)
			if onExpr == nil {
				result = append(result, combined)
				matched = true
				rightMatched[i] = true
			} else {
				match, err := s.evaluateWhere(onExpr, combined.Values)
				if err != nil {
					return nil, err
				}
				if match {
					result = append(result, combined)
					matched = true
					rightMatched[i] = true
				}
			}
		}
		if !matched {
			combined := s.mergeRowsWithNullRight(leftRow, leftName, rightName, rightCols)
			result = append(result, combined)
		}
	}

	for i, rightRow := range rightRows {
		if !rightMatched[i] {
			combined := s.mergeRowsWithNullLeft(rightRow, leftName, rightName, leftCols)
			result = append(result, combined)
		}
	}

	return result, nil
}

func (s *AthenaService) mergeRows(leftRow, rightRow *athenastore.StoredRow, leftName, rightName string) *athenastore.StoredRow {
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

func (s *AthenaService) mergeRowsWithNullRight(leftRow *athenastore.StoredRow, leftName, rightName string, rightCols []string) *athenastore.StoredRow {
	combined := &athenastore.StoredRow{Values: make(map[string]interface{})}

	for k, v := range leftRow.Values {
		combined.Values[leftName+"."+k] = v
		combined.Values[k] = v
	}
	for _, col := range rightCols {
		combined.Values[rightName+"."+col] = nil
		if _, exists := combined.Values[col]; !exists {
			combined.Values[col] = nil
		}
	}

	return combined
}

func (s *AthenaService) mergeRowsWithNullLeft(rightRow *athenastore.StoredRow, leftName, rightName string, leftCols []string) *athenastore.StoredRow {
	combined := &athenastore.StoredRow{Values: make(map[string]interface{})}

	for _, col := range leftCols {
		combined.Values[leftName+"."+col] = nil
		if _, exists := combined.Values[col]; !exists {
			combined.Values[col] = nil
		}
	}
	for k, v := range rightRow.Values {
		combined.Values[rightName+"."+k] = v
		if _, exists := combined.Values[k]; !exists {
			combined.Values[k] = v
		}
	}

	return combined
}

func (s *AthenaService) hasAggregateFunctions(exprs sqlparser.SelectExprs) bool {
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

func (s *AthenaService) executeGroupBy(selectStmt *sqlparser.Select, rows []*athenastore.StoredRow) []*athenastore.StoredRow {
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

func (s *AthenaService) buildGroupKey(groupBy sqlparser.GroupBy, row map[string]interface{}) string {
	var sb strings.Builder
	for _, expr := range groupBy {
		colName := s.extractColumnName(expr)
		var val string
		if v, ok := row[colName]; ok {
			val = fmt.Sprintf("%v", v)
		}
		// Length-prefix encoding prevents separator collision: if a value
		// contains the delimiter character, the length prefix disambiguates
		// the field boundaries.
		sb.WriteString(strconv.Itoa(len(val)))
		sb.WriteByte(':')
		sb.WriteString(val)
	}
	return sb.String()
}

func (s *AthenaService) executeAggregateWithoutGroup(selectStmt *sqlparser.Select, rows []*athenastore.StoredRow) []*athenastore.StoredRow {
	aggregatedRow := s.aggregateGroup(selectStmt, rows)
	return []*athenastore.StoredRow{aggregatedRow}
}

func (s *AthenaService) aggregateGroup(selectStmt *sqlparser.Select, rows []*athenastore.StoredRow) *athenastore.StoredRow {
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

func (s *AthenaService) computeAggregate(fn *sqlparser.FuncExpr, rows []*athenastore.StoredRow) interface{} {
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
		return s.aggregateMinMax(fn, rows, true)

	case "MAX":
		return s.aggregateMinMax(fn, rows, false)

	default:
		return nil
	}
}

func (s *AthenaService) aggregateNumeric(fn *sqlparser.FuncExpr, rows []*athenastore.StoredRow, op string) interface{} {
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

// aggregateMinMax computes MIN or MAX over a column supporting numeric,
// string, and mixed-type values. Numeric values are compared numerically;
// non-numeric values fall back to lexicographic comparison via compareValues.
func (s *AthenaService) aggregateMinMax(fn *sqlparser.FuncExpr, rows []*athenastore.StoredRow, isMin bool) interface{} {
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

	var best interface{}
	for _, row := range rows {
		val, ok := row.Values[colName]
		if !ok || val == nil {
			continue
		}
		if best == nil {
			best = val
			continue
		}
		cmp := s.compareValues(val, best)
		if (isMin && cmp < 0) || (!isMin && cmp > 0) {
			best = val
		}
	}
	return best
}

func (s *AthenaService) toFloat(val interface{}) (float64, error) {
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

func (s *AthenaService) applyHaving(rows []*athenastore.StoredRow, expr sqlparser.Expr) ([]*athenastore.StoredRow, error) {
	var result []*athenastore.StoredRow

	for _, row := range rows {
		match, err := s.evaluateWhere(expr, row.Values)
		if err != nil {
			return nil, err
		}
		if match {
			result = append(result, row)
		}
	}

	return result, nil
}
