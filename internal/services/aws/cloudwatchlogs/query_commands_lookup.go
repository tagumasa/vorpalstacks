package cloudwatchlogs

import (
	"fmt"
	"net"
	"strings"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// parsedLookupTable is the query-time view of a stored lookup table: the
// CSV header and data records with a column index.
type parsedLookupTable struct {
	columns  []string
	rows     [][]string
	colIndex map[string]int
}

// newParsedLookupTable indexes a parsed CSV table for the lookup commands.
// loadParsedLookupTable resolves a lookup table to its parsed CSV form:
// the shared body of the execContext's getLookupTable closure on both
// the interactive and the scheduled query planes.
func (s *LogsService) loadParsedLookupTable(store *logsstore.Store, region, name string) (*parsedLookupTable, error) {
	lt, err := store.GetLookupTable(name)
	if err != nil {
		return nil, fmt.Errorf("lookup table %s not found", name)
	}
	body, err := s.lookupTablePlainBody(lt, region)
	if err != nil {
		return nil, fmt.Errorf("lookup table %s is unavailable: %v", name, err)
	}
	columns, records, err := parseLookupCSV(body)
	if err != nil {
		return nil, fmt.Errorf("lookup table %s is invalid: %v", name, err)
	}
	return newParsedLookupTable(columns, records), nil
}

// lookupColumn resolves one field name to its column index in the
// parsed table, reporting an unknown field through the shared error
// format.
func lookupColumn(table *parsedLookupTable, tableName, field string) (int, error) {
	col, ok := table.column(field)
	if !ok {
		return 0, fmt.Errorf("lookup table %s has no column %s", tableName, field)
	}
	return col, nil
}

// lookupColumns resolves field names to column indexes, reporting the
// first unknown field through the same shared error format.
func lookupColumns(table *parsedLookupTable, tableName string, fields []string) ([]int, error) {
	cols := make([]int, 0, len(fields))
	for _, f := range fields {
		col, err := lookupColumn(table, tableName, f)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func newParsedLookupTable(columns []string, rows [][]string) *parsedLookupTable {
	idx := make(map[string]int, len(columns))
	for i, c := range columns {
		idx[c] = i
	}
	return &parsedLookupTable{columns: columns, rows: rows, colIndex: idx}
}

// column returns the index of a column, reporting whether it exists.
func (t *parsedLookupTable) column(name string) (int, bool) {
	i, ok := t.colIndex[name]
	return i, ok
}

// loadLookupTable fetches and caches a lookup table for the running query.
func (ctx *execContext) loadLookupTable(name string) (*parsedLookupTable, error) {
	if ctx.lookupCache != nil {
		if t, ok := ctx.lookupCache[name]; ok {
			return t, nil
		}
	}
	if ctx.getLookupTable == nil {
		return nil, fmt.Errorf("lookup table %s is not available", name)
	}
	t, err := ctx.getLookupTable(name)
	if err != nil {
		return nil, err
	}
	if ctx.lookupCache == nil {
		ctx.lookupCache = map[string]*parsedLookupTable{}
	}
	ctx.lookupCache[name] = t
	return t, nil
}

// --- lookup ---

// lookupMatchField pairs a lookup table column with the log event field it
// is matched against.
type lookupMatchField struct {
	lookupField string
	logField    string
}

// lookupCommand enriches rows with reference data from a lookup table,
// matching one or more fields with AND logic.
type lookupCommand struct {
	table     string
	matches   []lookupMatchField
	outputNew bool
	outputs   []string
	headTk    token
}

func (c *lookupCommand) name() string { return "lookup" }

// parseLookupCommand parses
// lookup {{table}} {{lookup-field}}[as {{log-field}}][,...] OUTPUT|OUTPUTNEW {{output-field}}[,...]
func parseLookupCommand(args []token, head token) (command, error) {
	if len(args) < 3 {
		return nil, newQueryCompileError("Syntax error: lookup requires a table, match fields, an output mode, and output fields", head.start, head.end)
	}
	if args[0].kind != tokIdent && args[0].kind != tokBacktickIdent {
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected a table name", args[0].raw), args[0].start, args[0].end)
	}
	c := &lookupCommand{table: args[0].text, headTk: head}
	i := 1
	for i < len(args) {
		t := args[i]
		if t.kind == tokIdent && (strings.EqualFold(t.text, "OUTPUT") || strings.EqualFold(t.text, "OUTPUTNEW")) {
			c.outputNew = strings.EqualFold(t.text, "OUTPUTNEW")
			i++
			break
		}
		if t.kind != tokIdent && t.kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected a match field", t.raw), t.start, t.end)
		}
		if len(args) > i+2 && args[i+1].kind == tokIdent && strings.EqualFold(args[i+1].text, "as") {
			lf := args[i+2]
			if lf.kind != tokIdent && lf.kind != tokBacktickIdent {
				return nil, newQueryCompileError("Syntax error: as requires a field name", lf.start, lf.end)
			}
			c.matches = append(c.matches, lookupMatchField{lookupField: t.text, logField: lf.text})
			i += 3
		} else {
			c.matches = append(c.matches, lookupMatchField{lookupField: t.text, logField: t.text})
			i++
		}
		if i < len(args) && args[i].kind == tokComma {
			i++
		}
	}
	if len(c.matches) == 0 {
		return nil, newQueryCompileError("Syntax error: lookup requires match fields", head.start, head.end)
	}
	for ; i < len(args); i++ {
		t := args[i]
		if t.kind == tokComma {
			continue
		}
		if t.kind != tokIdent && t.kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected an output field", t.raw), t.start, t.end)
		}
		c.outputs = append(c.outputs, t.text)
	}
	if len(c.outputs) == 0 {
		return nil, newQueryCompileError("Syntax error: lookup requires output fields after OUTPUT or OUTPUTNEW", head.start, head.end)
	}
	return c, nil
}

func (c *lookupCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	table, err := ctx.loadLookupTable(c.table)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	matchFields := make([]string, 0, len(c.matches))
	for _, m := range c.matches {
		matchFields = append(matchFields, m.lookupField)
	}
	matchCols, err := lookupColumns(table, c.table, matchFields)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	outputCols, err := lookupColumns(table, c.table, c.outputs)
	if err != nil {
		ctx.sourceError = err
		return rows
	}

	// Index the first table row per composite match key; rows must match
	// all match fields (AND logic).
	index := make(map[string]int, len(table.rows))
	for ri, rec := range table.rows {
		key := lookupRowKey(rec, matchCols)
		if _, ok := index[key]; !ok {
			index[key] = ri
		}
	}

	for i := range rows {
		var key strings.Builder
		for _, m := range c.matches {
			key.WriteString(readField(rows[i], m.logField, ctx))
			key.WriteByte(0)
		}
		ri, ok := index[key.String()]
		if !ok {
			if !c.outputNew {
				// OUTPUT sets unmatched output fields to null.
				for _, out := range c.outputs {
					rows[i].set(out, "")
				}
			}
			continue
		}
		for oi, out := range c.outputs {
			v := table.rows[ri][outputCols[oi]]
			if c.outputNew {
				if existing, exists := rows[i].fields[out]; exists && existing != "" {
					continue
				}
			}
			rows[i].set(out, v)
		}
	}
	return rows
}

// lookupRowKey builds the composite match key of a table record.
func lookupRowKey(rec []string, cols []int) string {
	var b strings.Builder
	for _, c := range cols {
		if c < len(rec) {
			b.WriteString(rec[c])
		}
		b.WriteByte(0)
	}
	return b.String()
}

// --- cidrlookup ---

// cidrOutputField pairs a lookup table column with the event field it is
// written to.
type cidrOutputField struct {
	lookupField string
	eventField  string
}

// cidrlookupCommand enriches rows by matching an IP field against CIDR
// ranges stored in a lookup table column.
type cidrlookupCommand struct {
	table      string
	ipField    string
	cidrColumn string // empty selects the first column
	outputNew  bool
	outputs    []cidrOutputField
	headTk     token
}

func (c *cidrlookupCommand) name() string { return "cidrlookup" }

// parseCidrLookupCommand parses
// cidrlookup {{table}} {{ipField}} [as {{cidrColumn}}] OUTPUT|OUTPUTNEW {{lookupField}}[as {{eventField}}][,...]
func parseCidrLookupCommand(args []token, head token) (command, error) {
	if len(args) < 4 {
		return nil, newQueryCompileError("Syntax error: cidrlookup requires a table, an IP field, an output mode, and output fields", head.start, head.end)
	}
	c := &cidrlookupCommand{headTk: head}
	i := 0
	for ; i < 2 && i < len(args); i++ {
		if args[i].kind != tokIdent && args[i].kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected a field", args[i].raw), args[i].start, args[i].end)
		}
		if i == 0 {
			c.table = args[i].text
		} else {
			c.ipField = args[i].text
		}
	}
	if i < len(args) && args[i].kind == tokIdent && strings.EqualFold(args[i].text, "as") {
		if i+1 >= len(args) || (args[i+1].kind != tokIdent && args[i+1].kind != tokBacktickIdent) {
			return nil, newQueryCompileError("Syntax error: as requires a column name", args[i].start, args[i].end)
		}
		c.cidrColumn = args[i+1].text
		i += 2
	}
	if i >= len(args) || args[i].kind != tokIdent || (!strings.EqualFold(args[i].text, "OUTPUT") && !strings.EqualFold(args[i].text, "OUTPUTNEW")) {
		return nil, newQueryCompileError("Syntax error: cidrlookup requires OUTPUT or OUTPUTNEW", head.start, head.end)
	}
	c.outputNew = strings.EqualFold(args[i].text, "OUTPUTNEW")
	i++
	for ; i < len(args); i++ {
		t := args[i]
		if t.kind == tokComma {
			continue
		}
		if t.kind != tokIdent && t.kind != tokBacktickIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected an output field", t.raw), t.start, t.end)
		}
		out := cidrOutputField{lookupField: t.text, eventField: t.text}
		if i+2 < len(args) && args[i+1].kind == tokIdent && strings.EqualFold(args[i+1].text, "as") {
			ef := args[i+2]
			if ef.kind != tokIdent && ef.kind != tokBacktickIdent {
				return nil, newQueryCompileError("Syntax error: as requires a field name", ef.start, ef.end)
			}
			out.eventField = ef.text
			i += 2
		}
		c.outputs = append(c.outputs, out)
	}
	if len(c.outputs) == 0 {
		return nil, newQueryCompileError("Syntax error: cidrlookup requires output fields", head.start, head.end)
	}
	return c, nil
}

func (c *cidrlookupCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	table, err := ctx.loadLookupTable(c.table)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	cidrCol := 0
	if c.cidrColumn != "" {
		col, lerr := lookupColumn(table, c.table, c.cidrColumn)
		if lerr != nil {
			ctx.sourceError = lerr
			return rows
		}
		cidrCol = col
	}
	outputFields := make([]string, 0, len(c.outputs))
	for _, out := range c.outputs {
		outputFields = append(outputFields, out.lookupField)
	}
	outputCols, err := lookupColumns(table, c.table, outputFields)
	if err != nil {
		ctx.sourceError = err
		return rows
	}

	for i := range rows {
		ip := net.ParseIP(readField(rows[i], c.ipField, ctx))
		if ip == nil {
			if !c.outputNew {
				for _, out := range c.outputs {
					rows[i].set(out.eventField, "")
				}
			}
			continue
		}
		matched := false
		for _, rec := range table.rows {
			if cidrCol >= len(rec) {
				continue
			}
			_, network, err := net.ParseCIDR(rec[cidrCol])
			if err != nil || !network.Contains(ip) {
				continue
			}
			for oi, out := range c.outputs {
				v := rec[outputCols[oi]]
				if c.outputNew {
					if existing, exists := rows[i].fields[out.eventField]; exists && existing != "" {
						continue
					}
				}
				rows[i].set(out.eventField, v)
			}
			matched = true
			break
		}
		if !matched && !c.outputNew {
			for _, out := range c.outputs {
				rows[i].set(out.eventField, "")
			}
		}
	}
	return rows
}
