package cloudwatchlogs

import (
	"fmt"
	"strings"
)

// The query pipeline's compile-time validation: command order, the
// multi-stats field visibility rule, and the field-reference collection
// that rule rides on.

func validateQueryPipeline(queryString string) error {
	toks, err := lexQuery(queryString)
	if err != nil {
		return err
	}
	cmds, err := compilePipeline(toks)
	if err != nil {
		return err
	}
	if len(cmds) == 0 {
		return nil
	}
	if err := validateCommandOrder(cmds); err != nil {
		return err
	}
	return validateStatsFieldVisibility(cmds)
}

// validateCommandOrder enforces the documented command placement rules.
func validateCommandOrder(cmds []compiledCommand) error {
	lastStats := -1
	var statsHeads []token
	var joinHeads []token
	patternIdx := -1
	sortIdx := -1
	dedupIdx := -1
	for i, c := range cmds {
		switch c.cmd.name() {
		case "stats":
			statsHeads = append(statsHeads, c.head)
			lastStats = i
		case "join":
			joinHeads = append(joinHeads, c.head)
		case "pattern":
			if patternIdx < 0 {
				patternIdx = i
			}
		case "sort":
			sortIdx = i
		case "dedup":
			if dedupIdx < 0 {
				dedupIdx = i
			}
		}
	}
	if len(statsHeads) > 10 {
		// The error points at the first command beyond the limit.
		return newQueryCompileError("A query can have a maximum of 10 stats commands",
			statsHeads[10].start, statsHeads[10].end)
	}
	if len(joinHeads) > 1 {
		// The error points at the second join, which is the violation.
		return newQueryCompileError("Only one join command is supported per query",
			joinHeads[1].start, joinHeads[1].end)
	}
	if patternIdx >= 0 && sortIdx >= 0 && sortIdx < patternIdx {
		return newQueryCompileError(
			"A query is not valid if it includes a pattern command after a sort command", cmds[patternIdx].head.start, cmds[patternIdx].head.end)
	}
	// sort and limit must appear after the last stats command.
	for i, c := range cmds {
		n := c.cmd.name()
		if (n == "sort" || n == "limit") && lastStats >= 0 && i < lastStats {
			return newQueryCompileError(
				"If you use a sort or limit command, it must appear after the last stats command", c.head.start, c.head.end)
		}
	}
	// Only limit may follow dedup.
	if dedupIdx >= 0 {
		for i := dedupIdx + 1; i < len(cmds); i++ {
			if cmds[i].cmd.name() != "limit" {
				return newQueryCompileError(
					"The only query command that you can use after the dedup command is limit", cmds[i].head.start, cmds[i].head.end)
			}
		}
	}
	// SOURCE is only valid as the first command.
	for i, c := range cmds {
		if c.cmd.name() == "SOURCE" && i > 0 {
			return newQueryCompileError("SOURCE is only valid as the first command of a query", c.head.start, c.head.end)
		}
	}
	// "The estimate command must be the last command in the query."
	for i, c := range cmds {
		if c.cmd.name() == "estimate" && i != len(cmds)-1 {
			return newQueryCompileError("The estimate command must be the last command in the query", c.head.start, c.head.end)
		}
	}
	return nil
}

// validateStatsFieldVisibility enforces the multi-stats field rule: "In
// subsequent stats commands in a single query, you can refer only to
// fields that are defined in the preceding stats command" and "Any
// fields that you reference after a stats command must be defined in
// that stats command", plus the bin corollary — "The bin function always
// implicitly uses the @timestamp field. This means that you can't use
// bin in a subsequent stats command without using the preceding stats
// command to propagate the timestamp field." The first stats command
// sees every field; each stats replaces the visible set with its own
// output names (aggregation aliases and group-key names).
func validateStatsFieldVisibility(cmds []compiledCommand) error {
	visible := map[string]bool{}
	sawStats := false
	for _, c := range cmds {
		st, ok := c.cmd.(*statsCommand)
		if !ok {
			continue
		}
		if sawStats {
			refs := statsFieldReferences(st, visible)
			if len(refs) > 0 {
				return newQueryCompileError(
					fmt.Sprintf("Field %s is not available here: in subsequent stats commands you can refer only to fields that are defined in the preceding stats command", refs[0]),
					st.headTk.start, st.headTk.end)
			}
			for _, g := range st.groups {
				if g.bin != nil && !visible["@timestamp"] {
					return newQueryCompileError(
						"bin implicitly uses the @timestamp field, so a subsequent stats command can use bin only when the preceding stats command propagates the timestamp field",
						st.headTk.start, st.headTk.end)
				}
			}
		}
		visible = statsOutputNames(st)
		sawStats = true
	}
	return nil
}

// statsOutputNames names one stats command's output columns: every
// aggregation's alias (explicit or the auto-generated expression label)
// and every group key's name.
func statsOutputNames(st *statsCommand) map[string]bool {
	out := make(map[string]bool, len(st.aggs)+len(st.groups))
	for _, a := range st.aggs {
		out[a.alias] = true
	}
	for _, g := range st.groups {
		out[g.name] = true
	}
	return out
}

// statsFieldReferences collects the field references a stats command
// makes that the visible set does not carry. A function-call subtree
// whose rendered text names a visible output is a reference to that
// auto-named column ("You can only reference `SUM(Fault)` or Operation
// at this point") and is not descended into.
func statsFieldReferences(st *statsCommand, visible map[string]bool) []string {
	var missing []string
	seen := map[string]bool{}
	collect := func(e exprNode) {
		collectMissingFieldRefs(e, visible, seen, &missing)
	}
	for _, a := range st.aggs {
		collect(a.expr)
	}
	for _, g := range st.groups {
		collect(g.expr)
	}
	return missing
}

// collectMissingFieldRefs walks one expression tree gathering field-node
// leaf names the visible set lacks.
func collectMissingFieldRefs(e exprNode, visible, seen map[string]bool, missing *[]string) {
	switch n := e.(type) {
	case nil:
		return
	case *fieldNode:
		name := exprText(n)
		if !visible[name] && !seen[name] {
			seen[name] = true
			*missing = append(*missing, name)
		}
	case *funcNode:
		if visibleNameMatches(visible, exprText(n)) {
			return
		}
		for _, a := range n.args {
			collectMissingFieldRefs(a, visible, seen, missing)
		}
	case *binOpNode:
		collectMissingFieldRefs(n.l, visible, seen, missing)
		collectMissingFieldRefs(n.r, visible, seen, missing)
	case *unaryNode:
		collectMissingFieldRefs(n.x, visible, seen, missing)
	case *inNode:
		collectMissingFieldRefs(n.x, visible, seen, missing)
		for _, item := range n.list {
			collectMissingFieldRefs(item, visible, seen, missing)
		}
		collectMissingFieldRefs(n.operand, visible, seen, missing)
	case *likeNode:
		collectMissingFieldRefs(n.x, visible, seen, missing)
	case *accessNode:
		collectMissingFieldRefs(n.base, visible, seen, missing)
	case *arrayLiteralNode:
		for _, item := range n.items {
			collectMissingFieldRefs(item, visible, seen, missing)
		}
	}
}

// visibleNameMatches reports whether the rendered expression names a
// visible output column, case-insensitively (auto labels keep the query
// writer's casing).
func visibleNameMatches(visible map[string]bool, text string) bool {
	for name := range visible {
		if strings.EqualFold(name, text) {
			return true
		}
	}
	return false
}

// buildRows converts raw log events into the initial row set with the
// discoverable fields. JSON log messages contribute their top-level keys as
// discovered fields, matching the automatic field discovery of Logs
// Insights; nested structures are kept as canonical JSON strings.
// maxDiscoveredJSONFields is the documented ceiling on the number of fields
// Logs Insights extracts from a JSON log event; further fields are ignored
// and must be extracted with parse.
const maxDiscoveredJSONFields = 200
