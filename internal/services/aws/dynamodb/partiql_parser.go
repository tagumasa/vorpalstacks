package dynamodb

import (
	"fmt"
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
// DynamoDB PartiQL UPDATE statement, plus the statement's RETURNING
// clause value in its canonical form when one is present.
type updateClauses struct {
	setAssignments    []setAssignment
	removeAttrs       []string
	addAssignments    []setAssignment
	deleteAssignments []setAssignment
	returning         string
}

// updateClauseSegment is one keyword-delimited action segment of a manual
// UPDATE parse, kept in statement order: an UPDATE may repeat an action
// keyword (the reference documents a two-SET statement), so collecting the
// segments in a keyword-keyed map would let a later clause overwrite the
// earlier one's intent.
type updateClauseSegment struct {
	keyword string
	text    string
}

func parsePartiQLParams(params map[string]interface{}) *partiQLParams {
	p := &partiQLParams{}
	if parameters, ok := params["Parameters"].([]interface{}); ok {
		p.Parameters = parameters
	}
	return p
}

// preparePartiQLStatement rewrites every top-level `?` placeholder into its
// explicit :vN form — N is the 1-based position of the placeholder within
// the statement — and reports how many placeholders it rewrote. The
// explicit indices survive every later parse, including the per-segment
// re-parses of the manual UPDATE path, whose segment-local tokenisation
// restarts numbering at :v1 and would otherwise bind two placeholders from
// different clauses to the same parameter. The rewrite skips string
// literals and quoted identifiers, and is idempotent: an already-numbered
// statement carries no `?`.
func preparePartiQLStatement(statement string) (string, int) {
	var b strings.Builder
	count := 0
	n := len(statement)
	for i := 0; i < n; {
		c := statement[i]
		if c == '\'' || c == '"' {
			quote := c
			b.WriteByte(c)
			i++
			for i < n {
				if statement[i] == '\\' && i+1 < n {
					b.WriteByte(statement[i])
					b.WriteByte(statement[i+1])
					i += 2
					continue
				}
				ch := statement[i]
				b.WriteByte(ch)
				i++
				if ch == quote {
					break
				}
			}
			continue
		}
		if c == '?' {
			count++
			fmt.Fprintf(&b, ":v%d", count)
			i++
			continue
		}
		b.WriteByte(c)
		i++
	}
	return b.String(), count
}

// validateParameterCount enforces the rule every statement engine shares:
// the request carries exactly one parameter value per placeholder. The
// under-count direction is additionally rejected when a placeholder fails
// to resolve; this check closes the over-count direction, which AWS
// answers with a ValidationException.
func validateParameterCount(placeholders int, params *partiQLParams) error {
	if placeholders != len(params.Parameters) {
		return ErrInvalidParameter
	}
	return nil
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

func parseInsertStatementWithParams(statement string, params *partiQLParams) (tableName string, itemData map[string]*dbstore.AttributeValue, err error) {
	stmt, err := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil, ErrInvalidParameter
	}

	ins, ok := stmt.(*sqlparser.Insert)
	if !ok {
		return "", nil, ErrInvalidParameter
	}

	tableName = sqlparser.String(ins.Table)
	tableName = trimQuotes(tableName)

	values, ok := ins.Rows.(sqlparser.Values)
	if !ok || len(values) != 1 || len(values[0]) != 1 {
		return "", nil, ErrInvalidParameter
	}

	obj, ok := values[0][0].(*sqlparser.ObjectLiteral)
	if !ok {
		return "", nil, ErrInvalidParameter
	}

	itemData, err = objectLiteralToAttributesWithParams(obj, params)
	if err != nil {
		return "", nil, err
	}
	return tableName, itemData, nil
}

// findClauseKeywordPositions scans s for top-level (non-quoted) occurrences
// of PartiQL UPDATE clause keywords: SET, REMOVE, ADD, DELETE, WHERE,
// RETURNING. Keywords inside single-quoted string literals or double-quoted
// identifiers are skipped. The result is a slice of [start, end] byte
// offsets, matching the semantics of regexp FindAllStringIndex so callers
// can use the same indexing logic.
func findClauseKeywordPositions(s string) [][2]int {
	keywords := []string{"SET", "REMOVE", "ADD", "DELETE", "WHERE", "RETURNING"}

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

func parseUpdateStatement(statement string) (tableName string, clauses updateClauses, whereExpr sqlparser.Expr, err error) {
	// First try the standard sqlparser — handles SET-only UPDATE.
	stmt, parseErr := sqlparser.ParseWithOptions(statement, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if parseErr == nil {
		if upd, ok := stmt.(*sqlparser.Update); ok {
			tableName = extractTableNameFromExprs(upd.TableExprs)
			for _, expr := range upd.Exprs {
				name := colNameDocumentPath(expr.Name)
				clauses.setAssignments = append(clauses.setAssignments, setAssignment{
					attrName: name,
					value:    expr.Expr,
				})
			}
			if upd.Where != nil {
				whereExpr = upd.Where.Expr
			}
			return tableName, clauses, whereExpr, nil
		}
	}

	// sqlparser fails on REMOVE/ADD/DELETE — parse manually.
	return parseUpdateStatementManual(statement)
}

// parseUpdateStatementManual splits the statement on top-level keywords
// and extracts table name, clause segments, and WHERE expression. Any
// segment value or WHERE that fails its reparse fails the whole statement
// — a partially-parsed UPDATE must never execute with dropped intent —
// and every failure returns an empty table name, so lock-key derivation
// locks nothing for a statement the engine is about to reject.
func parseUpdateStatementManual(statement string) (tableName string, clauses updateClauses, whereExpr sqlparser.Expr, err error) {
	upper := strings.ToUpper(statement)

	// Find the UPDATE keyword.
	updateIdx := strings.Index(upper, "UPDATE")
	if updateIdx != 0 {
		return "", clauses, nil, ErrInvalidParameter
	}
	rest := statement[updateIdx+6:]

	// Find all clause keyword positions in rest, skipping quoted sections.
	matches := findClauseKeywordPositions(rest)

	if len(matches) == 0 {
		return "", clauses, nil, ErrInvalidParameter
	}

	// Table name is between UPDATE and the first keyword.
	firstKW := matches[0][0]
	tableName = trimQuotes(strings.TrimSpace(rest[:firstKW]))

	// Extract WHERE and RETURNING expressions separately; action segments
	// keep their statement order.
	whereText := ""
	returningText := ""
	var actionSegments []updateClauseSegment
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

		switch kw {
		case "WHERE":
			whereText = segment
		case "RETURNING":
			returningText = segment
		default:
			actionSegments = append(actionSegments, updateClauseSegment{keyword: kw, text: segment})
		}
	}

	// The RETURNING clause value is part of the statement grammar: a value
	// outside the documented set fails the statement before anything
	// executes, never executes with the requested return dropped.
	if returningText != "" {
		returning, rErr := parseReturningClause(returningText)
		if rErr != nil {
			return "", clauses, nil, rErr
		}
		clauses.returning = returning
	}

	// Parse the WHERE clause with the sqlparser expression engine. A WHERE
	// that fails its reparse fails the statement: executing the update with
	// no filter would be the opposite of the request's intent.
	if whereText != "" {
		// Wrap in a fake SELECT to parse the WHERE expression.
		fakeStmt := "SELECT * FROM t WHERE " + whereText
		parsed, pErr := sqlparser.ParseWithOptions(fakeStmt, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
		if pErr != nil {
			return "", clauses, nil, ErrInvalidParameter
		}
		sel, ok := parsed.(*sqlparser.Select)
		if !ok || sel.Where == nil {
			return "", clauses, nil, ErrInvalidParameter
		}
		whereExpr = sel.Where.Expr
	}

	// Parse each action segment in statement order: a repeated keyword
	// appends its assignments to the same clause list, so a two-SET
	// statement keeps both clauses' writes.
	for _, seg := range actionSegments {
		switch seg.keyword {
		case "SET":
			assignments, segErr := parseSetSegment(seg.text)
			if segErr != nil {
				return "", clauses, nil, segErr
			}
			clauses.setAssignments = append(clauses.setAssignments, assignments...)
		case "REMOVE":
			clauses.removeAttrs = append(clauses.removeAttrs, parseRemoveSegment(seg.text)...)
		case "ADD":
			assignments, segErr := parseAddDeleteSegment(seg.text)
			if segErr != nil {
				return "", clauses, nil, segErr
			}
			clauses.addAssignments = append(clauses.addAssignments, assignments...)
		case "DELETE":
			assignments, segErr := parseAddDeleteSegment(seg.text)
			if segErr != nil {
				return "", clauses, nil, segErr
			}
			clauses.deleteAssignments = append(clauses.deleteAssignments, assignments...)
		}
	}

	return tableName, clauses, whereExpr, nil
}

// parseSetSegment parses "col = expr, col = expr" into setAssignments.
// A segment the expression engine cannot parse fails the statement —
// returning no assignments would let the update execute without the
// requested writes.
func parseSetSegment(seg string) ([]setAssignment, error) {
	// Use sqlparser to parse a fake UPDATE to extract expressions.
	fakeStmt := `UPDATE "t" SET ` + seg
	stmt, err := sqlparser.ParseWithOptions(fakeStmt, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return nil, ErrInvalidParameter
	}
	upd, ok := stmt.(*sqlparser.Update)
	if !ok {
		return nil, ErrInvalidParameter
	}
	var result []setAssignment
	for _, expr := range upd.Exprs {
		name := colNameDocumentPath(expr.Name)
		result = append(result, setAssignment{attrName: name, value: expr.Expr})
	}
	return result, nil
}

// parseRemoveSegment parses "col, col, col" into a list of attribute names.
func parseRemoveSegment(seg string) []string {
	parts := strings.Split(seg, ",")
	var result []string
	for _, p := range parts {
		name := renderPathText(strings.TrimSpace(p))
		if name != "" {
			result = append(result, name)
		}
	}
	return result
}

// parseAddDeleteSegment parses "col expr, col expr" into setAssignments.
// ADD and DELETE both use the same col + value syntax. A value that fails
// its reparse — or a target carrying no value at all — fails the segment:
// silently dropping it would execute the statement without the requested
// action.
func parseAddDeleteSegment(seg string) ([]setAssignment, error) {
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
			return nil, ErrInvalidParameter
		}
		name := renderPathText(strings.TrimSpace(pair[:spaceIdx]))
		valStr := strings.TrimSpace(pair[spaceIdx+1:])

		// Parse the value via a fake SET UPDATE.
		fakeStmt := `UPDATE "t" SET "fake" = ` + valStr
		stmt, err := sqlparser.ParseWithOptions(fakeStmt, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
		if err != nil {
			return nil, ErrInvalidParameter
		}
		upd, ok := stmt.(*sqlparser.Update)
		if !ok || len(upd.Exprs) == 0 {
			return nil, ErrInvalidParameter
		}
		result = append(result, setAssignment{attrName: name, value: upd.Exprs[0].Expr})
	}
	return result, nil
}

// parseReturningClause validates one RETURNING clause value against the
// documented grammar and returns its canonical form: the four UPDATE values
// built from ALL/MODIFIED and OLD/NEW plus the closing star. The value is
// case-insensitive like every PartiQL keyword; anything outside the set
// fails the statement.
func parseReturningClause(text string) (string, error) {
	canonical := strings.ToUpper(strings.Join(strings.Fields(text), " "))
	switch canonical {
	case "ALL OLD *", "MODIFIED OLD *", "ALL NEW *", "MODIFIED NEW *":
		return canonical, nil
	}
	return "", ErrInvalidParameter
}

// splitReturningSuffix splits a statement into its body and the text of a
// trailing top-level RETURNING clause, if one is present. The scan skips
// single-quoted string literals and double-quoted identifiers, so a literal
// 'RETURNING' inside a value never splits the statement.
func splitReturningSuffix(statement string) (body, returning string) {
	isWordChar := func(b byte) bool {
		return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
	}
	n := len(statement)
	i := 0
	for i < n {
		c := statement[i]
		if c == '\'' || c == '"' {
			quote := c
			i++
			for i < n {
				if statement[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if statement[i] == quote {
					i++
					break
				}
				i++
			}
			continue
		}
		if i+9 <= n && strings.EqualFold(statement[i:i+9], "RETURNING") &&
			(i == 0 || !isWordChar(statement[i-1])) &&
			(i+9 >= n || !isWordChar(statement[i+9])) {
			return strings.TrimSpace(statement[:i]), strings.TrimSpace(statement[i+9:])
		}
		i++
	}
	return statement, ""
}

// parseDeleteStatement parses a PartiQL DELETE FROM statement. The DELETE
// grammar carries one optional RETURNING clause whose only documented value
// is ALL OLD * — a returnvalues value outside that set fails the statement,
// which then yields an empty table name so lock-key derivation locks
// nothing for it.
func parseDeleteStatement(statement string) (tableName string, whereExpr sqlparser.Expr, returning string) {
	body, returningText := splitReturningSuffix(statement)
	if returningText != "" {
		parsed, rErr := parseReturningClause(returningText)
		if rErr != nil || parsed != "ALL OLD *" {
			return "", nil, ""
		}
		returning = parsed
	}
	stmt, err := sqlparser.ParseWithOptions(body, sqlparser.ParserOptions{Dialect: sqlparser.DialectPartiQL})
	if err != nil {
		return "", nil, ""
	}

	del, ok := stmt.(*sqlparser.Delete)
	if !ok {
		return "", nil, ""
	}

	tableName = extractTableNameFromExprs(del.TableExprs)

	if del.Where != nil {
		whereExpr = del.Where.Expr
	}

	return tableName, whereExpr, returning
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

// selectExprName renders a SELECT-plane name expression — an order key or a
// projection — as the attribute name it addresses: a column reference
// renders its raw segment values (an identifier the printer would escape, a
// keyword or a non-ASCII name, keeps its name), and any other expression
// keeps its printed form.
func selectExprName(expr sqlparser.Expr) string {
	if col, ok := expr.(*sqlparser.ColName); ok {
		return colNameDocumentPath(col)
	}
	return trimQuotes(sqlparser.String(expr))
}

// renderPathText renders a clause target's statement text as the document
// path it names: the strip runs per segment after a quote-aware split,
// because a qualified target carries one quote pair per segment — stripping
// the whole text's outer pair alone would leave the inner segments quoted.
// A segment keeps any bracket index it carries; the doc-path grammar parses
// the suffix on its own.
func renderPathText(text string) string {
	segments := make([]string, 0, strings.Count(text, ".")+1)
	var current strings.Builder
	var quote rune
	for _, r := range text {
		switch {
		case quote != 0:
			current.WriteRune(r)
			if r == quote {
				quote = 0
			}
		case r == '"' || r == '`':
			quote = r
			current.WriteRune(r)
		case r == '.':
			segments = append(segments, trimQuotes(current.String()))
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	segments = append(segments, trimQuotes(current.String()))
	return strings.Join(segments, ".")
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
			column:    selectExprName(sel.OrderBy[0].Expr),
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
			selectCols = append(selectCols, selectExprName(e.Expr))
		}
	}
	if len(selectCols) == 0 {
		selectCols = nil
	}

	return tableName, whereExpr, orderBy, selectCols
}
