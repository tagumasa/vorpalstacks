package dynamodb

import (
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

type partiQLParams struct {
	Parameters []interface{}
}

type setAssignment struct {
	attrName string
	value    sqlparser.Expr
}

func parsePartiQLParams(params map[string]interface{}) *partiQLParams {
	p := &partiQLParams{}
	if parameters, ok := params["Parameters"].([]interface{}); ok {
		p.Parameters = parameters
	}
	return p
}

func parseSelectStatement(statement string) (tableName string, whereExpr sqlparser.Expr) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil
	}

	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return "", nil
	}

	tableName = sqlparser.String(sel.From)
	tableName = trimQuotes(tableName)

	if sel.Where != nil {
		whereExpr = sel.Where.Expr
	}

	return tableName, whereExpr
}

func parseInsertStatement(statement string) (tableName string, itemData map[string]*dbstore.AttributeValue) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil
	}

	ins, ok := stmt.(*sqlparser.Insert)
	if !ok {
		return "", nil
	}

	tableName = sqlparser.String(ins.Table)
	tableName = trimQuotes(tableName)

	values, ok := ins.Rows.(sqlparser.Values)
	if !ok || len(values) != 1 || len(values[0]) != 1 {
		return "", nil
	}

	obj, ok := values[0][0].(*sqlparser.ObjectLiteral)
	if !ok {
		return "", nil
	}

	itemData = objectLiteralToAttributesWithParams(obj, nil)
	return tableName, itemData
}

func parseInsertStatementWithParams(statement string, params *partiQLParams) (tableName string, itemData map[string]*dbstore.AttributeValue) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil
	}

	ins, ok := stmt.(*sqlparser.Insert)
	if !ok {
		return "", nil
	}

	tableName = sqlparser.String(ins.Table)
	tableName = trimQuotes(tableName)

	values, ok := ins.Rows.(sqlparser.Values)
	if !ok || len(values) != 1 || len(values[0]) != 1 {
		return "", nil
	}

	obj, ok := values[0][0].(*sqlparser.ObjectLiteral)
	if !ok {
		return "", nil
	}

	itemData = objectLiteralToAttributesWithParams(obj, params)
	return tableName, itemData
}

func parseUpdateStatement(statement string) (tableName string, assignments []setAssignment, whereExpr sqlparser.Expr) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil, nil
	}

	upd, ok := stmt.(*sqlparser.Update)
	if !ok {
		return "", nil, nil
	}

	tableName = extractTableNameFromExprs(upd.TableExprs)

	for _, expr := range upd.Exprs {
		name := sqlparser.String(expr.Name)
		name = trimQuotes(name)
		assignments = append(assignments, setAssignment{
			attrName: name,
			value:    expr.Expr,
		})
	}

	if upd.Where != nil {
		whereExpr = upd.Where.Expr
	}

	return tableName, assignments, whereExpr
}

func parseDeleteStatement(statement string) (tableName string, whereExpr sqlparser.Expr) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil
	}

	del, ok := stmt.(*sqlparser.Delete)
	if !ok {
		return "", nil
	}

	tableName = extractTableNameFromExprs(del.TableExprs)

	if del.Where != nil {
		whereExpr = del.Where.Expr
	}

	return tableName, whereExpr
}

func extractTableNameFromExprs(tableExprs sqlparser.TableExprs) string {
	if len(tableExprs) == 0 {
		return ""
	}
	aliased, ok := tableExprs[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return ""
	}

	var name string
	switch t := aliased.Expr.(type) {
	case *sqlparser.TableName:
		name = sqlparser.String(t)
	case sqlparser.TableName:
		name = sqlparser.String(&t)
	default:
		return ""
	}

	return trimQuotes(name)
}

func trimQuotes(name string) string {
	if len(name) >= 2 && name[0] == '"' && name[len(name)-1] == '"' {
		return name[1 : len(name)-1]
	}
	if len(name) >= 2 && name[0] == '`' && name[len(name)-1] == '`' {
		return name[1 : len(name)-1]
	}
	return name
}

type orderByClause struct {
	column    string
	direction string
}

func parseSelectStatementWithOrderBy(statement string) (tableName string, whereExpr sqlparser.Expr, orderBy *orderByClause, selectCols []string) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil, nil, nil
	}

	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return "", nil, nil, nil
	}

	tableName = sqlparser.String(sel.From)
	tableName = trimQuotes(tableName)

	if sel.Where != nil {
		whereExpr = sel.Where.Expr
	}

	if len(sel.OrderBy) > 0 {
		orderBy = &orderByClause{
			column:    sqlparser.String(sel.OrderBy[0].Expr),
			direction: "ASC",
		}
		if sel.OrderBy[0].Direction == sqlparser.DescScr {
			orderBy.direction = "DESC"
		}
	}

	for _, expr := range sel.SelectExprs {
		switch e := expr.(type) {
		case *sqlparser.StarExpr:
			selectCols = nil
		case *sqlparser.AliasedExpr:
			colName := sqlparser.String(e.Expr)
			colName = trimQuotes(colName)
			selectCols = append(selectCols, colName)
		}
	}
	if len(selectCols) == 0 {
		selectCols = nil
	}

	return tableName, whereExpr, orderBy, selectCols
}
