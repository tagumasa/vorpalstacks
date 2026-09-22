package cloudwatchlogs

import (
	"fmt"
	"strings"

	filterpattern "vorpalstacks/pkg/filterpattern"
)

// The fieldSelectionCriteria members of PutMetricFilter and
// PutSubscriptionFilter: "A filter expression that specifies which log
// events should be processed by this metric filter based on system fields
// such as source account and source region. Uses selection criteria
// syntax with operators like =, !=, AND, OR, IN, NOT IN" — the worked
// examples run @aws.region = "us-east-1", @aws.account IN [...] and
// @aws.region NOT IN [...]. The system fields the criteria can address
// are the two the platform stamps on every ingested event: the region of
// the group it landed in and the owning account.

// fieldSelectionFieldValue resolves one system field of the events being
// processed.
func fieldSelectionFieldValue(field, region, accountID string) (string, bool) {
	switch field {
	case "@aws.region":
		return region, true
	case "@aws.account":
		return accountID, true
	}
	return "", false
}

// evalFieldSelectionCriteria reports whether the events being processed
// satisfy a filter's fieldSelectionCriteria; an absent criteria processes
// everything. The grammar's composition is the documented operator set:
// OR alternations over AND conjunctions (an AND binds tighter than an
// enclosing OR, the operator set's SQL shape).
func evalFieldSelectionCriteria(criteria, region, accountID string) bool {
	criteria = strings.TrimSpace(criteria)
	if criteria == "" {
		return true
	}
	if parts := filterpattern.SplitLogical(criteria, " OR "); len(parts) > 1 {
		for _, p := range parts {
			if evalFieldSelectionCriteria(p, region, accountID) {
				return true
			}
		}
		return false
	}
	if parts := filterpattern.SplitLogical(criteria, " AND "); len(parts) > 1 {
		for _, p := range parts {
			if !evalFieldSelectionCriteria(p, region, accountID) {
				return false
			}
		}
		return true
	}
	return evalFieldSelectionLeaf(criteria, region, accountID)
}

// evalFieldSelectionLeaf evaluates one comparison: a system field, an
// operator (=, !=, IN, NOT IN) and a quoted value or bracketed list of
// quoted values. A leaf naming an unknown system field matches nothing.
func evalFieldSelectionLeaf(leaf, region, accountID string) bool {
	field, op, values, ok := parseFieldSelectionLeaf(leaf)
	if !ok {
		return false
	}
	actual, known := fieldSelectionFieldValue(field, region, accountID)
	if !known {
		return false
	}
	switch op {
	case "=":
		return len(values) == 1 && actual == values[0]
	case "!=":
		return len(values) == 1 && actual != values[0]
	case "IN":
		for _, v := range values {
			if actual == v {
				return true
			}
		}
		return false
	case "NOT IN":
		for _, v := range values {
			if actual == v {
				return false
			}
		}
		return true
	}
	return false
}

// parseFieldSelectionLeaf splits one comparison into its field, operator
// and values; ok reports whether the leaf is grammatical. The value side
// (a quoted string or a bracketed list) begins at the leaf's first quote
// or bracket — the field name and operator carry neither — so every
// operator scan runs over the head before it alone: an " IN " or "!="
// inside a quoted value is value text, never an operator.
func parseFieldSelectionLeaf(leaf string) (field, op string, values []string, ok bool) {
	leaf = strings.TrimSpace(leaf)
	valueStart := strings.IndexAny(leaf, "\"'[")
	if valueStart < 0 {
		return "", "", nil, false
	}
	head := leaf[:valueStart]
	rest := strings.TrimSpace(leaf[valueStart:])
	upper := strings.ToUpper(head)
	cut := func(i, n int) (string, string) {
		return strings.TrimSpace(head[:i]), strings.TrimSpace(head[i+n:])
	}
	var rhs string
	if i := strings.Index(upper, " NOT IN "); i > 0 {
		field, rhs = cut(i, len(" NOT IN "))
		op = "NOT IN"
	} else if i := strings.Index(upper, " IN "); i > 0 {
		field, rhs = cut(i, len(" IN "))
		op = "IN"
	} else if i := strings.Index(head, "!="); i > 0 {
		field, rhs = cut(i, 2)
		op = "!="
	} else if i := strings.Index(head, "="); i > 0 {
		field, rhs = cut(i, 1)
		op = "="
	} else {
		return "", "", nil, false
	}
	if rhs != "" {
		return "", "", nil, false
	}
	values, ok = parseFieldSelectionValues(rest)
	return field, op, values, ok
}

// parseFieldSelectionValues parses the leaf's right-hand side: one quoted
// string or a bracketed list of quoted strings.
func parseFieldSelectionValues(rest string) ([]string, bool) {
	rest = strings.TrimSpace(rest)
	if strings.HasPrefix(rest, "\"") || strings.HasPrefix(rest, "'") {
		q := rest[0:1]
		if len(rest) < 2 || !strings.HasSuffix(rest, q) {
			return nil, false
		}
		return []string{rest[1 : len(rest)-1]}, true
	}
	if !strings.HasPrefix(rest, "[") || !strings.HasSuffix(rest, "]") {
		return nil, false
	}
	inner := strings.TrimSpace(rest[1 : len(rest)-1])
	if inner == "" {
		return nil, false
	}
	var values []string
	for _, part := range strings.Split(inner, ",") {
		part = strings.TrimSpace(part)
		if len(part) < 2 || part[0] != '"' || part[len(part)-1] != '"' {
			return nil, false
		}
		values = append(values, part[1:len(part)-1])
	}
	return values, true
}

// validateFieldSelectionCriteriaSyntax walks the criteria with the
// evaluator's own grammar and rejects what it cannot parse: an unknown
// operator form, a malformed value list, or a field outside the system
// field vocabulary.
func validateFieldSelectionCriteriaSyntax(criteria string) error {
	criteria = strings.TrimSpace(criteria)
	if criteria == "" {
		return nil
	}
	if parts := filterpattern.SplitLogical(criteria, " OR "); len(parts) > 1 {
		for _, p := range parts {
			if err := validateFieldSelectionCriteriaSyntax(p); err != nil {
				return err
			}
		}
		return nil
	}
	if parts := filterpattern.SplitLogical(criteria, " AND "); len(parts) > 1 {
		for _, p := range parts {
			if err := validateFieldSelectionCriteriaSyntax(p); err != nil {
				return err
			}
		}
		return nil
	}
	field, op, values, ok := parseFieldSelectionLeaf(criteria)
	if !ok {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid fieldSelectionCriteria leaf %q: expected a system field compared with =, !=, IN or NOT IN against a quoted value or list", criteria), 400)
	}
	if _, known := fieldSelectionFieldValue(field, "", ""); !known {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("fieldSelectionCriteria addresses system fields by name; %q is not one (valid fields: @aws.region, @aws.account)", field), 400)
	}
	if op == "NOT IN" || op == "IN" {
		if len(values) == 0 {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The %s comparison of fieldSelectionCriteria needs a bracketed list of values", op), 400)
		}
	}
	return nil
}

// validateEmitSystemFields enforces the PutSubscriptionFilter member's
// documented value set: "A list of system fields to include in the log
// events sent to the subscription destination. Valid values are
// @aws.account, @aws.region, and @source.log."
func validateEmitSystemFields(fields []string) error {
	for _, f := range fields {
		if f != "@aws.account" && f != "@aws.region" && f != "@source.log" {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid emitSystemFields entry %s. Valid values: @aws.account, @aws.region, @source.log", f), 400)
		}
	}
	return nil
}
