// Shared helper functions for Gremlin step execution.

package gremlinparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/internal/core/storage/graphengine"
)

// matchHasStep evaluates a has() step's arguments against a traverser.
func matchHasStep(ec *ExecContext, t *Traverser, args []Argument) bool {
	if len(args) == 0 {
		return true
	}

	if args[0].Kind == ArgEnum && args[0].Enum != nil {
		key := args[0].Enum.Value
		if key == "id" {
			if len(args) == 1 {
				return true
			}
			elID := elementID(t.Element)
			if args[1].Kind == ArgPredicate {
				if idInt, err := strconv.ParseInt(elID, 10, 64); err == nil {
					return matchPredicate(idInt, args[1].Pred)
				}
				return matchPredicate(elID, args[1].Pred)
			}
			return elID == fmt.Sprintf("%v", resolveArgValue(args[1]))
		}
		if key == "label" {
			if len(args) == 1 {
				return true
			}
			elLabel := elementLabel(t.Element)
			if args[1].Kind == ArgPredicate {
				return matchPredicate(elLabel, args[1].Pred)
			}
			return elLabel == fmt.Sprintf("%v", resolveArgValue(args[1]))
		}
		if len(args) == 1 {
			props := elementProps(t.Element)
			_, ok := props[key]
			return ok
		}
		if len(args) >= 2 {
			if args[1].Kind == ArgPredicate {
				props := elementProps(t.Element)
				val, ok := props[key]
				return ok && matchPredicate(val, args[1].Pred)
			}
			props := elementProps(t.Element)
			val, ok := props[key]
			return ok && valuesEqual(val, resolveArgValue(args[1]))
		}
	}

	if args[0].Kind == ArgString {
		key := args[0].Str
		if len(args) == 1 {
			props := elementProps(t.Element)
			_, ok := props[key]
			return ok
		}
		if len(args) >= 2 {
			if args[1].Kind == ArgPredicate {
				props := elementProps(t.Element)
				val, ok := props[key]
				return ok && matchPredicate(val, args[1].Pred)
			}
			props := elementProps(t.Element)
			val, ok := props[key]
			return ok && valuesEqual(val, resolveArgValue(args[1]))
		}
	}

	return true
}

// matchLabels returns true if all labels in matchLabels are present in nodeLabels.
func matchLabels(nodeLabels []string, matchLabels []string) bool {
	if len(matchLabels) == 0 {
		return true
	}
	labelSet := make(map[string]bool, len(nodeLabels))
	for _, l := range nodeLabels {
		labelSet[l] = true
	}
	for _, ml := range matchLabels {
		if !labelSet[ml] {
			return false
		}
	}
	return true
}

// matchPredicate evaluates a Gremlin predicate against a value.
// Supports eq, neq, gt, gte, lt, lte, between, inside, outside, within, without,
// containing, notContaining, startingWith, notStartingWith, endingWith, notEndingWith,
// regex, notRegex.
func matchPredicate(val any, pred *Predicate) bool {
	if pred == nil {
		return true
	}

	result := evalPredicateImpl(val, pred)
	if pred.Negate {
		return !result
	}
	return result
}

func evalPredicateImpl(val any, pred *Predicate) bool {

	switch pred.Method {
	case "eq":
		if len(pred.Args) > 0 {
			return elementOrValueEqual(val, resolveArgValue(pred.Args[0]))
		}
		return false
	case "neq":
		if len(pred.Args) > 0 {
			return !elementOrValueEqual(val, resolveArgValue(pred.Args[0]))
		}
		return true
	case "gt":
		return compareNumeric(val, resolveArgValue(pred.Args[0])) > 0
	case "gte":
		return compareNumeric(val, resolveArgValue(pred.Args[0])) >= 0
	case "lt":
		return compareNumeric(val, resolveArgValue(pred.Args[0])) < 0
	case "lte":
		return compareNumeric(val, resolveArgValue(pred.Args[0])) <= 0
	case "between":
		if len(pred.Args) < 2 {
			return false
		}
		cmp1 := compareNumeric(val, resolveArgValue(pred.Args[0]))
		cmp2 := compareNumeric(val, resolveArgValue(pred.Args[1]))
		return cmp1 != -2 && cmp2 != -2 && cmp1 >= 0 && cmp2 <= 0
	case "inside":
		if len(pred.Args) < 2 {
			return false
		}
		cmp1 := compareNumeric(val, resolveArgValue(pred.Args[0]))
		cmp2 := compareNumeric(val, resolveArgValue(pred.Args[1]))
		return cmp1 != -2 && cmp2 != -2 && cmp1 > 0 && cmp2 < 0
	case "outside":
		if len(pred.Args) < 2 {
			return false
		}
		cmp1 := compareNumeric(val, resolveArgValue(pred.Args[0]))
		cmp2 := compareNumeric(val, resolveArgValue(pred.Args[1]))
		return cmp1 != -2 && cmp2 != -2 && (cmp1 < 0 || cmp2 > 0)
	case "within":
		if len(pred.Args) == 0 {
			return false
		}
		for _, arg := range pred.Args {
			if elementOrValueEqual(val, resolveArgValue(arg)) {
				return true
			}
		}
		return false
	case "without":
		if len(pred.Args) == 0 {
			return true
		}
		for _, arg := range pred.Args {
			if elementOrValueEqual(val, resolveArgValue(arg)) {
				return false
			}
		}
		return true
	case "containing":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			return strings.Contains(str, argString(pred.Args[0]))
		}
		return false
	case "notContaining":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			return !strings.Contains(str, argString(pred.Args[0]))
		}
		return false
	case "startingWith":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			return strings.HasPrefix(str, argString(pred.Args[0]))
		}
		return false
	case "notStartingWith":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			return !strings.HasPrefix(str, argString(pred.Args[0]))
		}
		return false
	case "endingWith":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			return strings.HasSuffix(str, argString(pred.Args[0]))
		}
		return false
	case "notEndingWith":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			return !strings.HasSuffix(str, argString(pred.Args[0]))
		}
		return false
	case "regex":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			matched, err := regexp.MatchString(argString(pred.Args[0]), str)
			return err == nil && matched
		}
		return false
	case "notRegex":
		if str, ok := toString(val); ok && len(pred.Args) > 0 {
			matched, err := regexp.MatchString(argString(pred.Args[0]), str)
			return err == nil && !matched
		}
		return false
	case "and":
		for _, arg := range pred.Args {
			if arg.Kind == ArgPredicate {
				if !matchPredicate(val, arg.Pred) {
					return false
				}
			}
		}
		return true
	case "or":
		for _, arg := range pred.Args {
			if arg.Kind == ArgPredicate {
				if matchPredicate(val, arg.Pred) {
					return true
				}
			}
		}
		return false
	}

	return false
}

// filterTraversers filters a slice of traversers by the given predicate function.
func filterTraversers(traversers []*Traverser, fn func(*Traverser) bool) []*Traverser {
	result := make([]*Traverser, 0, len(traversers))
	for _, t := range traversers {
		if fn(t) {
			result = append(result, t)
		}
	}
	return result
}

// elementID returns the string representation of an element's ID.
func elementID(el any) string {
	switch v := el.(type) {
	case *graphengine.Node:
		return fmt.Sprintf("%d", v.ID)
	case *graphengine.Edge:
		return fmt.Sprintf("%d", v.ID)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// elementLabel returns the label of a node (first label) or edge.
func elementLabel(el any) string {
	switch v := el.(type) {
	case *graphengine.Node:
		if len(v.Labels) > 0 {
			return v.Labels[0]
		}
		return ""
	case *graphengine.Edge:
		return v.Label
	default:
		return ""
	}
}

// elementProps returns the property map of a node or edge, or nil for other types.
func elementProps(el any) map[string]any {
	switch v := el.(type) {
	case *graphengine.Node:
		return v.Props
	case *graphengine.Edge:
		return v.Props
	default:
		return nil
	}
}

// resolveArgValue converts an Argument AST node to its runtime Go value.
func resolveArgValue(arg Argument) any {
	switch arg.Kind {
	case ArgString:
		return arg.Str
	case ArgInt:
		return arg.Int
	case ArgFloat:
		return arg.Float
	case ArgBool:
		return arg.Bool
	case ArgNull:
		return nil
	case ArgList:
		result := make([]any, len(arg.List))
		for i, a := range arg.List {
			result[i] = resolveArgValue(a)
		}
		return result
	case ArgMap:
		m := make(map[string]any)
		for _, entry := range arg.Map {
			keyVal := resolveArgValue(entry.Key)
			keyStr, ok := keyVal.(string)
			if !ok {
				return fmt.Errorf("gremlin: map key must be a string, got %T", keyVal)
			}
			m[keyStr] = resolveArgValue(entry.Value)
		}
		return m
	case ArgEnum:
		return arg.Enum.Value
	default:
		return arg.Str
	}
}

// toNodeID converts an Argument to a graphengine.NodeID.
func toNodeID(arg Argument) (graphengine.NodeID, error) {
	switch arg.Kind {
	case ArgInt:
		return graphengine.NodeID(arg.Int), nil
	case ArgFloat:
		return graphengine.NodeID(int64(arg.Float)), nil
	case ArgString:
		s := strings.TrimPrefix(arg.Str, "n")
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("gremlin: toNodeID parse error: %q: %w", arg.Str, err)
		}
		return graphengine.NodeID(id), nil
	default:
		return 0, fmt.Errorf("gremlin: toNodeID unsupported arg kind: %v", arg.Kind)
	}
}

// toEdgeID converts an Argument to a graphengine.EdgeID.
func toEdgeID(arg Argument) (graphengine.EdgeID, error) {
	switch arg.Kind {
	case ArgInt:
		return graphengine.EdgeID(arg.Int), nil
	case ArgFloat:
		return graphengine.EdgeID(int64(arg.Float)), nil
	case ArgString:
		s := strings.TrimPrefix(arg.Str, "n")
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("gremlin: toEdgeID parse error: %q: %w", arg.Str, err)
		}
		return graphengine.EdgeID(id), nil
	default:
		return 0, fmt.Errorf("gremlin: toEdgeID unsupported arg kind: %v", arg.Kind)
	}
}

// toFloat64 attempts to convert a value to float64 for numeric comparisons.
// Returns (0, false) for non-numeric types.
func toFloat64(val any) (float64, bool) {
	switch v := val.(type) {
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	case float64:
		return v, true
	case int32:
		return float64(v), true
	default:
		return 0, false
	}
}

// compareNumeric compares two values numerically. Returns -2 if either value is non-numeric.
func compareNumeric(a, b any) int {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if !aok || !bok {
		return -2
	}
	if af < bf {
		return -1
	}
	if af > bf {
		return 1
	}
	return 0
}

// toString converts a value to its string representation for text predicates.
func toString(val any) (string, bool) {
	switch v := val.(type) {
	case string:
		return v, true
	case int64:
		return strconv.FormatInt(v, 10), true
	case int:
		return strconv.Itoa(v), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		if v {
			return "true", true
		}
		return "false", true
	default:
		return "", false
	}
}

// containsVal checks if a value exists in a slice using elementOrValueEqual for identity comparison.
func containsVal(vals []any, val any) bool {
	for _, v := range vals {
		if elementOrValueEqual(v, val) {
			return true
		}
	}
	return false
}
