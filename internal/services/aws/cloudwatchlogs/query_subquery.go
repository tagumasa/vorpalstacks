package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strings"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// subqueryNode is a nested query used as the operand of filter's in
// operator. Subqueries execute independently of the outer query and their
// result sets are materialised before being consumed; results are cached on
// the execution context so evaluation per row does not re-run the subquery.
type subqueryNode struct {
	toks    []token
	startTk token
}

func (n *subqueryNode) eval(row *queryResultRow, ctx *execContext) interface{} {
	return n.evalSet(ctx)
}

// evalSet runs the subquery and returns the values of its first output
// field, which are the values the outer in operator matches against. The
// first field is the first column of the subquery's output rows, resolved
// through the row's column order so the choice is deterministic.
func (n *subqueryNode) evalSet(ctx *execContext) []interface{} {
	key := fmt.Sprintf("%d:%d", n.startTk.start, n.startTk.end)
	if cached, ok := ctx.subqueryCache[key]; ok {
		return cached
	}
	var values []interface{}
	rows, err := ctx.runSubquery(n.toks)
	if err == nil {
		for i := range rows {
			cols := rows[i].ordered()
			if len(cols) > 0 {
				values = append(values, rows[i].fields[cols[0]])
			}
		}
	}
	if values == nil {
		values = []interface{}{}
	}
	ctx.subqueryCache[key] = values
	return values
}

// --- SOURCE ---

// sourceCommand selects the log groups the query runs against. It is only
// valid as the first command of a query.
type sourceCommand struct {
	groupFilters []logGroupFilter
	tagFilters   []tagFilter
	dataSources  []string
	headTk       token
}

type logGroupFilter struct {
	kind   string // namePrefix, accountIdentifier, class
	values []string
}

type tagFilter struct {
	key string
	// Each value test compiles once at parse time: the documented glob
	// vocabulary with a leading "!" negating the match.
	tests []tagValueTest
}

type tagValueTest struct {
	negated bool
	re      *regexp.Regexp
}

func (c *sourceCommand) name() string { return "SOURCE" }

// parseSourceCommand parses the documented SOURCE forms:
//
//	SOURCE logGroups(namePrefix: ['a', 'b'], class: ['STANDARD'])
//	SOURCE logGroupTags([{key: 'team', values: ['t1']}])
//	SOURCE dataSource(['amazon_vpc.flow'])
func parseSourceCommand(args []token, head token) (command, error) {
	c := &sourceCommand{headTk: head}
	i := 0
	for i < len(args) {
		t := args[i]
		if t.kind != tokIdent {
			return nil, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", t.raw), t.start, t.end)
		}
		lower := strings.ToLower(t.text)
		switch lower {
		case "loggroups":
			filters, next, err := parseLogGroupFilters(args, i)
			if err != nil {
				return nil, err
			}
			c.groupFilters = append(c.groupFilters, filters...)
			i = next
		case "loggrouptags":
			filters, next, err := parseTagFilters(args, i)
			if err != nil {
				return nil, err
			}
			c.tagFilters = append(c.tagFilters, filters...)
			i = next
		case "datasource":
			names, next, err := parseParenStringList(args, i)
			if err != nil {
				return nil, err
			}
			if len(c.dataSources)+len(names) > logsstore.MaxQueryDataSources {
				return nil, newQueryCompileError(fmt.Sprintf("A query can include as many as %d data sources", logsstore.MaxQueryDataSources), t.start, args[next-1].end)
			}
			c.dataSources = append(c.dataSources, names...)
			i = next
		default:
			return nil, newQueryCompileError(
				fmt.Sprintf("Syntax error at '%s': expected logGroups, logGroupTags or dataSource", t.raw), t.start, t.end)
		}
	}
	return c, nil
}

// parseLogGroupFilters parses logGroups(prefix: [...], class: [...],
// accountIdentifier: [...]) starting at index i of args; returns the index
// after the closing parenthesis. The documented quotas (5 name prefixes,
// 20 account identifiers) are enforced per parameter list.
func parseLogGroupFilters(args []token, i int) ([]logGroupFilter, int, error) {
	if i+1 >= len(args) || args[i+1].kind != tokLParen {
		return nil, 0, newQueryCompileError("Syntax error: logGroups requires parentheses", args[i].start, args[i].end)
	}
	var filters []logGroupFilter
	j := i + 2
	for j < len(args) && args[j].kind != tokRParen {
		kw := args[j]
		if kw.kind != tokIdent {
			return nil, 0, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", kw.raw), kw.start, kw.end)
		}
		if j+1 >= len(args) || args[j+1].kind != tokColon {
			return nil, 0, newQueryCompileError("Syntax error: expected ':'", kw.start, kw.end)
		}
		values, next, err := parseBracketStringList(args, j+2)
		if err != nil {
			return nil, 0, err
		}
		kind := strings.ToLower(kw.text)
		if kind == "nameprefix" && len(values) > logsstore.MaxQueryNamePrefixes {
			return nil, 0, newQueryCompileError(fmt.Sprintf("A query can include as many as %d name prefixes", logsstore.MaxQueryNamePrefixes), kw.start, args[next-1].end)
		}
		if kind == "accountidentifier" && len(values) > logsstore.MaxQueryAccountIdentifiers {
			return nil, 0, newQueryCompileError(fmt.Sprintf("A query can include as many as %d account identifiers", logsstore.MaxQueryAccountIdentifiers), kw.start, args[next-1].end)
		}
		filters = append(filters, logGroupFilter{kind: kind, values: values})
		j = next
		if j < len(args) && args[j].kind == tokComma {
			j++
		}
	}
	if j >= len(args) {
		return nil, 0, newQueryCompileError("Syntax error: missing ')'", args[i+1].start, args[i+1].end)
	}
	return filters, j + 1, nil
}

// parseTagFilters parses logGroupTags([{key: 'k', values: [...]}]). The
// documented quotas allow at most 5 tag filters with up to 5 values each.
func parseTagFilters(args []token, i int) ([]tagFilter, int, error) {
	if i+1 >= len(args) || args[i+1].kind != tokLParen {
		return nil, 0, newQueryCompileError("Syntax error: logGroupTags requires parentheses", args[i].start, args[i].end)
	}
	if i+2 >= len(args) || args[i+2].kind != tokLBracket {
		return nil, 0, newQueryCompileError("Syntax error: logGroupTags requires a list", args[i+1].start, args[i+1].end)
	}
	var filters []tagFilter
	j := i + 3
	for j < len(args) && args[j].kind != tokRBracket {
		if args[j].kind == tokComma {
			j++
			continue
		}
		if args[j].kind != tokLBrace {
			return nil, 0, newQueryCompileError("Syntax error: expected '{'", args[j].start, args[j].end)
		}
		if len(filters) == 5 {
			return nil, 0, newQueryCompileError("A query can include as many as 5 tag filters", args[j].start, args[j].end)
		}
		j++
		var tf tagFilter
		for j < len(args) && args[j].kind != tokRBrace {
			if args[j].kind == tokComma {
				j++
				continue
			}
			kw := args[j]
			if kw.kind != tokIdent {
				return nil, 0, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", kw.raw), kw.start, kw.end)
			}
			if j+1 >= len(args) || args[j+1].kind != tokColon {
				return nil, 0, newQueryCompileError("Syntax error: expected ':'", kw.start, kw.end)
			}
			if strings.EqualFold(kw.text, "key") {
				if j+2 >= len(args) || args[j+2].kind != tokString {
					return nil, 0, newQueryCompileError("Syntax error: key requires a string", kw.start, kw.end)
				}
				tf.key = args[j+2].text
				j += 3
				continue
			}
			if strings.EqualFold(kw.text, "values") {
				values, next, err := parseBracketStringList(args, j+2)
				if err != nil {
					return nil, 0, err
				}
				if len(values) > 5 {
					return nil, 0, newQueryCompileError("A tag filter can include as many as 5 values", kw.start, args[next-1].end)
				}
				for _, v := range values {
					negated := strings.HasPrefix(v, "!")
					pattern := strings.TrimPrefix(v, "!")
					tf.tests = append(tf.tests, tagValueTest{negated: negated, re: compileGlob(pattern)})
				}
				j = next
				continue
			}
			return nil, 0, newQueryCompileError(fmt.Sprintf("Syntax error at '%s'", kw.raw), kw.start, kw.end)
		}
		if j >= len(args) {
			return nil, 0, newQueryCompileError("Syntax error: missing '}'", args[i+1].start, args[i+1].end)
		}
		j++ // consume }
		filters = append(filters, tf)
	}
	if j >= len(args) {
		return nil, 0, newQueryCompileError("Syntax error: missing ']'", args[i+2].start, args[i+2].end)
	}
	j++ // consume ]
	if j >= len(args) || args[j].kind != tokRParen {
		return nil, 0, newQueryCompileError("Syntax error: missing ')'", args[i+1].start, args[i+1].end)
	}
	return filters, j + 1, nil
}

// parseParenStringList parses dataSource(['a', 'b']).
func parseParenStringList(args []token, i int) ([]string, int, error) {
	if i+1 >= len(args) || args[i+1].kind != tokLParen {
		return nil, 0, newQueryCompileError("Syntax error: dataSource requires parentheses", args[i].start, args[i].end)
	}
	values, next, err := parseBracketStringList(args, i+2)
	if err != nil {
		return nil, 0, err
	}
	if next >= len(args) || args[next].kind != tokRParen {
		return nil, 0, newQueryCompileError("Syntax error: missing ')'", args[i+1].start, args[i+1].end)
	}
	return values, next + 1, nil
}

// parseBracketStringList parses ['a', 'b'] starting at index i.
func parseBracketStringList(args []token, i int) ([]string, int, error) {
	if i >= len(args) || args[i].kind != tokLBracket {
		return nil, 0, newQueryCompileError("Syntax error: expected a list", args[i].start, args[i].end)
	}
	var values []string
	j := i + 1
	for j < len(args) && args[j].kind != tokRBracket {
		if args[j].kind == tokComma {
			j++
			continue
		}
		if args[j].kind != tokString {
			return nil, 0, newQueryCompileError(fmt.Sprintf("Syntax error at '%s': expected a string", args[j].raw), args[j].start, args[j].end)
		}
		values = append(values, args[j].text)
		j++
	}
	if j >= len(args) {
		return nil, 0, newQueryCompileError("Syntax error: missing ']'", args[i].start, args[i].end)
	}
	return values, j + 1, nil
}

// apply resolves the selected log groups and refetches the query window.
// The refetched events replace the input row set, because the rows built
// from the default log groups are not part of the SOURCE selection.
func (c *sourceCommand) apply(ctx *execContext, rows []queryResultRow) []queryResultRow {
	groups, err := c.resolve(ctx)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	// An empty selection records as an empty, non-nil slice: the field's
	// nil state means "no SOURCE command ran" (the executed-group set the
	// query record serves distinguishes the two).
	if groups == nil {
		groups = []string{}
	}
	ctx.effectiveGroups = groups
	events, err := ctx.fetchEvents(groups, ctx.startTime, ctx.endTime)
	if err != nil {
		ctx.sourceError = err
		return rows
	}
	ctx.events = events
	return buildRows(events, ctx.accountID)
}

// maxSourceLogGroups is the documented ceiling on the number of log groups
// one query can select; exceeding it fails the query with an error
// prompting the user to narrow the selection.
const maxSourceLogGroups = 10000

// resolve computes the log group names matching the SOURCE selection.
func (c *sourceCommand) resolve(ctx *execContext) ([]string, error) {
	if ctx.listLogGroups == nil {
		return nil, nil
	}
	all, err := ctx.listLogGroups()
	if err != nil {
		return nil, nil
	}
	var prefixes, classes []string
	accountsPresent := false
	accountMatches := false
	for _, f := range c.groupFilters {
		switch f.kind {
		case "nameprefix":
			prefixes = append(prefixes, f.values...)
		case "class":
			classes = append(classes, f.values...)
		case "accountidentifier":
			// Single-account deployment: selecting the current account keeps
			// the default scope; any other identifier selects nothing.
			accountsPresent = true
			for _, v := range f.values {
				if v == ctx.accountID {
					accountMatches = true
				}
			}
		}
	}
	if accountsPresent && !accountMatches {
		return nil, nil
	}
	if len(classes) == 0 {
		classes = []string{"STANDARD"}
	}
	var out []string
	for _, g := range all {
		if !classMatches(classes, g.Class) {
			continue
		}
		if len(prefixes) > 0 && !hasAnyPrefix(g.Name, prefixes) {
			continue
		}
		if len(c.tagFilters) > 0 && !tagFiltersMatch(c.tagFilters, g.Tags) {
			continue
		}
		if len(c.dataSources) > 0 {
			// No data sources are ingested in this platform, so data source
			// selections match no log groups.
			continue
		}
		out = append(out, g.Name)
	}
	if len(out) > maxSourceLogGroups {
		return nil, NewLogsError("LimitExceededException",
			fmt.Sprintf("The query selects %d log groups, which exceeds the maximum of %d. Narrow the selection with name prefixes, tags, or classes.",
				len(out), maxSourceLogGroups), 400)
	}
	return out, nil
}

func classMatches(classes []string, groupClass string) bool {
	if groupClass == "" {
		groupClass = "STANDARD"
	}
	for _, c := range classes {
		if strings.EqualFold(c, groupClass) {
			return true
		}
	}
	return false
}

func hasAnyPrefix(name string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// tagFiltersMatch applies the documented tag selection semantics: filters
// with different keys combine with AND, values within one filter with OR,
// a filter without values selects every group carrying the tag, * is a
// wildcard prefix, and ! negates.
func tagFiltersMatch(filters []tagFilter, tags map[string]string) bool {
	for _, f := range filters {
		value, ok := tags[f.key]
		if !ok {
			// A tag filter selects groups that carry the tag; negated
			// values exclude values, not the absence of the tag itself.
			return false
		}
		if len(f.tests) == 0 {
			// The values array is optional: a key-only filter matches any
			// group that carries the tag.
			continue
		}
		matched := false
		for _, test := range f.tests {
			globHit := test.re != nil && test.re.MatchString(value)
			if test.negated {
				if !globHit {
					matched = true
				}
				continue
			}
			if globHit {
				matched = true
			}
		}
		if !matched {
			return false
		}
	}
	return true
}
