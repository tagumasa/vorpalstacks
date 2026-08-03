package dynamodb

import (
	"strings"

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

// updateClauses holds parsed SET/REMOVE/ADD/DELETE clauses from a
// DynamoDB PartiQL UPDATE statement.
type updateClauses struct {
	setAssignments    []setAssignment
	removeAttrs       []string
	addAssignments    []setAssignment
	deleteAssignments []setAssignment
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
	return parseInsertStatementWithParams(statement, nil)
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

// findClauseKeywordPositions scans s for top-level (non-quoted) occurrences
// of PartiQL UPDATE clause keywords: SET, REMOVE, ADD, DELETE, WHERE.
// Keywords inside single-quoted string literals or double-quoted identifiers
// are skipped. The result is a slice of [start, end] byte offsets, matching
// the semantics of regexp FindAllStringIndex so callers can use the same
// indexing logic.
func findClauseKeywordPositions(s string) [][2]int {
	keywords := []string{"SET", "REMOVE", "ADD", "DELETE", "WHERE"}

	upperS := strings.ToUpper(s)
	n := len(s)

	isWordChar := func(b byte) bool {
		return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
	}

	var result [][2]int
	i := 0
	for i < n {
		c := s[i]

		// Skip quoted sections — identifiers ("...") and string literals ('...').
		if c == '"' || c == '\'' {
			quote := c
			i++
			for i < n {
				if s[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if s[i] == quote {
					i++
					break
				}
				i++
			}
			continue
		}

		// Attempt keyword match at the current position.
		matched := false
		for _, kw := range keywords {
			kwLen := len(kw)
			if i+kwLen > n {
				continue
			}
			if upperS[i:i+kwLen] != kw {
				continue
			}
			// Word-boundary check: preceding character must not be a word char.
			if i > 0 && isWordChar(s[i-1]) {
				continue
			}
			// Word-boundary check: following character must not be a word char.
			after := i + kwLen
			if after < n && isWordChar(s[after]) {
				continue
			}
			result = append(result, [2]int{i, after})
			i = after
			matched = true
			break
		}
		if !matched {
			i++
		}
	}
	return result
}

func parseUpdateStatement(statement string) (tableName string, clauses updateClauses, whereExpr sqlparser.Expr) {
	// First try the standard sqlparser — handles SET-only UPDATE.
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err == nil {
		if upd, ok := stmt.(*sqlparser.Update); ok {
			tableName = extractTableNameFromExprs(upd.TableExprs)
			for _, expr := range upd.Exprs {
				name := trimQuotes(sqlparser.String(expr.Name))
				clauses.setAssignments = append(clauses.setAssignments, setAssignment{
					attrName: name,
					value:    expr.Expr,
				})
			}
			if upd.Where != nil {
				whereExpr = upd.Where.Expr
			}
			return tableName, clauses, whereExpr
		}
	}

	// sqlparser fails on REMOVE/ADD/DELETE — parse manually.
	return parseUpdateStatementManual(statement)
}

// parseUpdateStatementManual splits the statement on top-level keywords
// and extracts table name, clause segments, and WHERE expression.
func parseUpdateStatementManual(statement string) (tableName string, clauses updateClauses, whereExpr sqlparser.Expr) {
	upper := strings.ToUpper(statement)

	// Find the UPDATE keyword.
	updateIdx := strings.Index(upper, "UPDATE")
	if updateIdx != 0 {
		return "", clauses, nil
	}
	rest := statement[updateIdx+6:]

	// Find all clause keyword positions in rest, skipping quoted sections.
	matches := findClauseKeywordPositions(rest)

	if len(matches) == 0 {
		return "", clauses, nil
	}

	// Table name is between UPDATE and the first keyword.
	firstKW := matches[0][0]
	tableName = trimQuotes(strings.TrimSpace(rest[:firstKW]))

	// Extract WHERE expression separately.
	whereText := ""
	clauseTexts := make(map[string]string)
	for i, m := range matches {
		kwStart := m[0]
		kwEnd := m[1]
		kw := strings.ToUpper(rest[kwStart:kwEnd])

		// Segment runs until the next keyword or end of string.
		segEnd := len(rest)
		if i+1 < len(matches) {
			segEnd = matches[i+1][0]
		}
		segment := strings.TrimSpace(rest[kwEnd:segEnd])

		if kw == "WHERE" {
			whereText = segment
		} else {
			clauseTexts[kw] = segment
		}
	}

	// Parse the WHERE clause with the sqlparser expression engine.
	if whereText != "" {
		// Wrap in a fake SELECT to parse the WHERE expression.
		fakeStmt := "SELECT * FROM t WHERE " + whereText
		if parsed, pErr := sqlparser.ParseWithOptions(fakeStmt, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL}); pErr == nil {
			if sel, ok := parsed.(*sqlparser.Select); ok && sel.Where != nil {
				whereExpr = sel.Where.Expr
			}
		}
	}

	// Parse each clause segment.
	if seg, ok := clauseTexts["SET"]; ok {
		clauses.setAssignments = parseSetSegment(seg)
	}
	if seg, ok := clauseTexts["REMOVE"]; ok {
		clauses.removeAttrs = parseRemoveSegment(seg)
	}
	if seg, ok := clauseTexts["ADD"]; ok {
		clauses.addAssignments = parseAddDeleteSegment(seg)
	}
	if seg, ok := clauseTexts["DELETE"]; ok {
		clauses.deleteAssignments = parseAddDeleteSegment(seg)
	}

	return tableName, clauses, whereExpr
}

// parseSetSegment parses "col = expr, col = expr" into setAssignments.
func parseSetSegment(seg string) []setAssignment {
	// Use sqlparser to parse a fake UPDATE to extract expressions.
	fakeStmt := `UPDATE "t" SET ` + seg
	stmt, err := sqlparser.ParseWithOptions(fakeStmt, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return nil
	}
	upd, ok := stmt.(*sqlparser.Update)
	if !ok {
		return nil
	}
	var result []setAssignment
	for _, expr := range upd.Exprs {
		name := trimQuotes(sqlparser.String(expr.Name))
		result = append(result, setAssignment{attrName: name, value: expr.Expr})
	}
	return result
}

// parseRemoveSegment parses "col, col, col" into a list of attribute names.
func parseRemoveSegment(seg string) []string {
	parts := strings.Split(seg, ",")
	var result []string
	for _, p := range parts {
		name := trimQuotes(strings.TrimSpace(p))
		if name != "" {
			result = append(result, name)
		}
	}
	return result
}

// parseAddDeleteSegment parses "col expr, col expr" into setAssignments.
// ADD and DELETE both use the same col + value syntax.
func parseAddDeleteSegment(seg string) []setAssignment {
	// Parse by splitting on commas that are NOT inside parentheses.
	var pairs []string
	depth := 0
	start := 0
	for i, ch := range seg {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		} else if ch == ',' && depth == 0 {
			pairs = append(pairs, seg[start:i])
			start = i + 1
		}
	}
	pairs = append(pairs, seg[start:])

	var result []setAssignment
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		// Split into name and value on the first space.
		spaceIdx := strings.IndexAny(pair, " \t")
		if spaceIdx < 0 {
			continue
		}
		name := trimQuotes(strings.TrimSpace(pair[:spaceIdx]))
		valStr := strings.TrimSpace(pair[spaceIdx+1:])

		// Parse the value via a fake SET UPDATE.
		fakeStmt := `UPDATE "t" SET "fake" = ` + valStr
		stmt, err := sqlparser.ParseWithOptions(fakeStmt, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
		if err != nil {
			continue
		}
		upd, ok := stmt.(*sqlparser.Update)
		if !ok || len(upd.Exprs) == 0 {
			continue
		}
		result = append(result, setAssignment{attrName: name, value: upd.Exprs[0].Expr})
	}
	return result
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
