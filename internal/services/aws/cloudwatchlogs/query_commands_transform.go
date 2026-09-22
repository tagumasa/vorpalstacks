package cloudwatchlogs

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"
)

// The extraction commands: parse (glob, regular expression, logfmt and
// CSV extraction) and json (chained JSON key extraction) — each adds
// fields derived from the row's own content. The row-window value and
// statistical commands live in their own files.

// --- parse (glob, regular expression, logfmt and CSV extraction) ---

// parseMode names one of the parse command's four documented modes: glob
// expressions, regular expressions, logfmt, and CSV.
type parseMode int

const (
	parseModeGlob parseMode = iota
	parseModeRegex
	parseModeLogfmt
	parseModeCSV
)

type parseExtractCommand struct {
	source exprNode
	re     *regexp.Regexp // the glob and regex modes' compiled pattern —
	// parse-time
	mode          parseMode // which of the four documented modes runs
	globWildcards int       // the glob mode's capture count, for the
	// parse-time name-arity check
	extracts []parseExtract // field names bound to capture group indexes
	multi    bool           // the regex mode's multi-match keyword
	headTk   token
}

// parseExtract binds one extracted field name to its position: the capture
// group index for the glob and regex modes' as clause, the subexpression
// index for the regex mode's named capture groups, and the one-based
// column position for the CSV mode's aliases.
type parseExtract struct {
	name string
	idx  int
}

func (c *parseExtractCommand) name() string { return "parse" }

// parseLeadsWithPattern reports whether a token opens the parse command's
// pattern or mode-keyword position directly — the omitted-fieldName form
// whose source defaults to @message.
func parseLeadsWithPattern(t token) bool {
	if t.kind == tokString || t.kind == tokRegex {
		return true
	}
	return t.kind == tokIdent && (strings.EqualFold(t.text, "logfmt") || strings.EqualFold(t.text, "csv"))
}

func parseParseCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: parse requires a pattern", head.start, head.end)
	}
	c := &parseExtractCommand{headTk: head}
	p := &exprParser{toks: args}
	// A leading pattern omits the source field: "If fieldName is omitted,
	// @message is used by default."
	if t, ok := p.peek(); ok && parseLeadsWithPattern(t) {
		c.source = &fieldNode{segs: []pathSeg{{name: "@message"}}}
	} else {
		src, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		c.source = src
	}
	pt, ok := p.next()
	if !ok {
		return nil, newQueryCompileError("Syntax error: parse requires a pattern", head.start, head.end)
	}
	switch pt.kind {
	case tokString:
		c.re = compileGlobCapture(pt.text)
		c.mode = parseModeGlob
		c.globWildcards = strings.Count(pt.text, "*")
	case tokRegex:
		re, err := regexp.Compile(pt.text)
		if err != nil {
			return nil, newQueryCompileError(fmt.Sprintf("Invalid regular expression: %v", err), pt.start, pt.end)
		}
		c.re = re
		c.mode = parseModeRegex
	case tokIdent:
		switch {
		case strings.EqualFold(pt.text, "logfmt"):
			c.mode = parseModeLogfmt
		case strings.EqualFold(pt.text, "csv"):
			c.mode = parseModeCSV
		default:
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", pt.raw), pt.start, pt.end)
		}
	default:
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", pt.raw), pt.start, pt.end)
	}
	hasAs := p.acceptKeyword("as")
	if c.mode != parseModeRegex && !hasAs {
		// The glob, logfmt and CSV modes have no capture names of their
		// own: their fields are assigned by the as clause of their
		// documented forms, parse {{fieldName}} "{{pattern}}" as
		// {{alias...}} / parse {{fieldName}} logfmt as {{alias}} /
		// parse {{fieldName}} csv as {{alias...}}.
		return nil, newQueryCompileError("Syntax error: parse requires as", head.start, head.end)
	}
	if hasAs {
		for {
			n, ok := p.next()
			if !ok || (n.kind != tokIdent && n.kind != tokBacktickIdent) {
				return nil, newQueryCompileError("Syntax error: parse requires field names after as", head.start, head.end)
			}
			c.extracts = append(c.extracts, parseExtract{name: n.text, idx: len(c.extracts) + 1})
			t, ok := p.next()
			if !ok {
				break
			}
			if t.kind == tokIdent && strings.EqualFold(t.text, "multi") {
				// The multi-match keyword closes the name list; the mode
				// switch below consumes it.
				p.pos--
				break
			}
			if t.kind != tokComma {
				return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
			}
		}
	} else {
		// The regex mode's documented form carries no as clause: named
		// capture groups (?<{{name}}>...) inside the pattern define the
		// extracted fields, in subexpression order.
		for i, name := range c.re.SubexpNames() {
			if i > 0 && name != "" {
				c.extracts = append(c.extracts, parseExtract{name: name, idx: i})
			}
		}
	}
	// The multi-match keyword follows the regex pattern (with or without
	// the as clause): "Add the keyword multi after the regex pattern",
	// producing multiple rows per log event. The other modes leave the
	// keyword to the trailing-token check.
	if c.mode == parseModeRegex {
		c.multi = p.acceptKeyword("multi")
	}
	if t, ok := p.next(); ok {
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
	}
	if c.mode == parseModeGlob && c.globWildcards != len(c.extracts) {
		// Glob mode: each * captures one field.
		return nil, newQueryCompileError(
			fmt.Sprintf("parse pattern has %d wildcards but %d field names", c.globWildcards, len(c.extracts)), head.start, head.end)
	}
	if c.mode == parseModeLogfmt && len(c.extracts) != 1 {
		return nil, newQueryCompileError("Syntax error: parse logfmt requires a single field name after as", head.start, head.end)
	}
	return c, nil
}

func (c *parseExtractCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	if c.multi {
		return c.applyMulti(ctx, rows)
	}
	switch c.mode {
	case parseModeLogfmt:
		return c.applyLogfmt(ctx, rows)
	case parseModeCSV:
		return c.applyCSV(ctx, rows)
	}
	return c.applyPattern(ctx, rows)
}

// applyPattern is the glob and regex modes: the first pattern match's
// capture groups are assigned to the extracted fields. An event that does
// not match keeps its row without the extracted fields.
func (c *parseExtractCommand) applyPattern(ctx *execContext, rows []queryResultRow) []queryResultRow {
	for i := range rows {
		row := rows[i]
		src := asString(c.source.eval(&row, ctx))
		if c.re == nil {
			continue
		}
		m := c.re.FindStringSubmatch(src)
		if m == nil {
			continue
		}
		for _, e := range c.extracts {
			if e.idx < len(m) {
				rows[i].set(e.name, m[e.idx])
			}
		}
	}
	return rows
}

// applyLogfmt is the logfmt mode: the source's space-separated key=value
// pairs become a map stored under the single alias — "The result is a map
// that you access with dot notation (for example, lf.level, lf.msg)". The
// map is created structure: it lands in the row's structural view, so dot
// notation traverses it while its canonical JSON string serves display.
// An event with no pairs keeps its row without the field.
func (c *parseExtractCommand) applyLogfmt(ctx *execContext, rows []queryResultRow) []queryResultRow {
	for i := range rows {
		row := rows[i]
		src := asString(c.source.eval(&row, ctx))
		pairs := splitLogfmtPairs(src)
		if len(pairs) == 0 {
			continue
		}
		rows[i].setStruct(c.extracts[0].name, pairs)
	}
	return rows
}

// applyCSV is the CSV mode: "Each comma-separated value is assigned to the
// corresponding alias." A record that does not parse keeps its row without
// the extracted fields.
func (c *parseExtractCommand) applyCSV(ctx *execContext, rows []queryResultRow) []queryResultRow {
	for i := range rows {
		row := rows[i]
		src := asString(c.source.eval(&row, ctx))
		rec, err := parseCSVRecord(src)
		if err != nil {
			continue
		}
		for _, e := range c.extracts {
			if e.idx-1 < len(rec) {
				rows[i].set(e.name, rec[e.idx-1])
			}
		}
	}
	return rows
}

// applyMulti is the multi-match mode: every match of the pattern emits its
// own row ("extract all matches of a regular expression from a field,
// producing multiple rows per log event"), and an event with no match
// keeps one field-less row — the parse command's documented unmatched-row
// behaviour.
func (c *parseExtractCommand) applyMulti(ctx *execContext, rows []queryResultRow) []queryResultRow {
	out := make([]queryResultRow, 0, len(rows))
	for i := range rows {
		row := rows[i]
		src := asString(c.source.eval(&row, ctx))
		if c.re == nil {
			out = append(out, row)
			continue
		}
		matches := c.re.FindAllStringSubmatch(src, -1)
		if len(matches) == 0 {
			out = append(out, row)
			continue
		}
		for _, m := range matches {
			r := cloneRow(row)
			for _, e := range c.extracts {
				if e.idx < len(m) {
					r.set(e.name, m[e.idx])
				}
			}
			out = append(out, r)
		}
	}
	return out
}

// compileGlobCapture compiles the parse command's glob pattern: the
// wildcards become capture groups, everything else is literal. Compiled
// once at parse time — the pattern is evaluated per row.
func compileGlobCapture(pattern string) *regexp.Regexp {
	parts := strings.Split(pattern, "*")
	if len(parts) < 2 {
		return nil
	}
	var re strings.Builder
	for i, p := range parts {
		if i > 0 {
			re.WriteString("(.*?)")
		}
		re.WriteString(regexp.QuoteMeta(p))
	}
	compiled, err := regexp.Compile("^" + re.String() + "$")
	if err != nil {
		return nil
	}
	return compiled
}

// splitLogfmtPairs parses one logfmt-formatted line — space-separated
// key=value pairs — into the logfmt mode's result map, following the
// reference format: values may be double-quoted (with \" and \\ escapes)
// to span spaces, and a bare token without '=' contributes its key with an
// empty value.
func splitLogfmtPairs(line string) map[string]interface{} {
	pairs := make(map[string]interface{})
	i, n := 0, len(line)
	for i < n {
		for i < n && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		if i >= n {
			break
		}
		start := i
		for i < n && line[i] != '=' && line[i] != ' ' && line[i] != '\t' {
			i++
		}
		key := line[start:i]
		if key == "" {
			// A stray separator with no key carries no pair.
			i++
			continue
		}
		if i >= n || line[i] != '=' {
			pairs[key] = ""
			continue
		}
		i++ // consume '='
		if i < n && line[i] == '"' {
			i++
			var val strings.Builder
			for i < n && line[i] != '"' {
				if line[i] == '\\' && i+1 < n && (line[i+1] == '"' || line[i+1] == '\\') {
					i++
				}
				val.WriteByte(line[i])
				i++
			}
			if i < n {
				i++ // closing quote
			}
			pairs[key] = val.String()
			continue
		}
		start = i
		for i < n && line[i] != ' ' && line[i] != '\t' {
			i++
		}
		pairs[key] = line[start:i]
	}
	return pairs
}

// parseCSVRecord parses one CSV-formatted line into its field values per
// RFC 4180 quoting; the CSV mode assigns each value to its corresponding
// alias positionally.
func parseCSVRecord(line string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(line))
	r.FieldsPerRecord = -1
	return r.Read()
}

// --- json (chained JSON key extraction) ---

type jsonExtractCommand struct {
	source exprNode
	segs   []pathSeg // the quoted key path, split like a field path
	alias  string
	headTk token
}

func (c *jsonExtractCommand) name() string { return "json" }

// parseJSONCommand parses json field={{fieldName}} "{{key.subkey}}" as
// {{alias}} — "explicit chained JSON extraction from a previously parsed
// object field. This enables you to extract nested keys from a structured
// field without re-parsing the raw message."
func parseJSONCommand(args []token, head token) (command, error) {
	if len(args) == 0 {
		return nil, newQueryCompileError("Syntax error: json requires field=<fieldName>", head.start, head.end)
	}
	c := &jsonExtractCommand{headTk: head}
	p := &exprParser{toks: args}
	if !p.acceptKeyword("field") {
		return nil, newQueryCompileError("Syntax error: json requires field=<fieldName>", head.start, head.end)
	}
	eq, ok := p.next()
	if !ok || eq.kind != tokOp || eq.text != "=" {
		return nil, newQueryCompileError("Syntax error: json requires field=<fieldName>", head.start, head.end)
	}
	src, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	c.source = src
	key, ok := p.next()
	if !ok || key.kind != tokString {
		return nil, newQueryCompileError("Syntax error: json requires a quoted key path", head.start, head.end)
	}
	c.segs = splitPathSegs(key.text)
	if len(c.segs) == 0 {
		return nil, newQueryCompileError("Syntax error: json requires a key path", key.start, key.end)
	}
	if !p.acceptKeyword("as") {
		return nil, newQueryCompileError("Syntax error: json requires as", head.start, head.end)
	}
	alias, ok := p.next()
	if !ok || (alias.kind != tokIdent && alias.kind != tokBacktickIdent) {
		return nil, newQueryCompileError("Syntax error: json requires a field name after as", head.start, head.end)
	}
	c.alias = alias.text
	if t, ok := p.next(); ok {
		return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
	}
	return c, nil
}

// apply extracts the nested key from the structured source field; a row
// whose source or key is absent keeps its fields — the alias stays unset.
// The command is the documented explicit parser for JSON-encoded string
// fields, so a string source decodes at entry.
func (c *jsonExtractCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	for i := range rows {
		row := rows[i]
		v := c.source.eval(&row, ctx)
		if s, ok := v.(string); ok {
			if decoded, ok := parseJSONValue(s); ok {
				v = decoded
			}
		}
		v = descendPath(v, c.segs)
		if v == nil {
			continue
		}
		rows[i].set(c.alias, storeValue(v))
	}
	return rows
}
