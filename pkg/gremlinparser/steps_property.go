// Property and mutation step executors: has, values, valueMap, properties, propertyMap, id, label, property, drop, count, identity, keys, sample, sack.

package gremlinparser

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	"vorpalstacks/pkg/graphengine"
)

// execHas filters traversers by property value, id, label, or predicate.
// Supports has(T.id, ...), has(T.label, ...), has('key', predicate), has('key', value).
func execHas(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return filterTraversers(traversers, func(t *Traverser) bool {
		return matchHasStep(ec, t, step.Args)
	}), nil
}

// execHasLabel filters traversers whose labels match all specified labels.
func execHasLabel(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser
	var errs []error
	for _, t := range traversers {
		if n, ok := asNode(t); ok {
			if matchLabels(n.Labels, labels) {
				result = append(result, t)
			}
		} else if e, ok := asEdge(t); ok {
			labelSet := make(map[string]bool)
			labelSet[e.Label] = true
			allMatch := true
			for _, l := range labels {
				if !labelSet[l] {
					allMatch = false
					break
				}
			}
			if allMatch {
				result = append(result, t)
			}
		} else {
			errs = append(errs, fmt.Errorf("gremlin: hasLabel() requires element input, got %T", t.Element))
		}
	}
	if len(errs) > 0 {
		return result, errs[0]
	}
	return result, nil
}

// execHasId filters traversers by element ID. Supports both direct values and predicates.
func execHasId(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) > 0 && step.Args[0].Kind == ArgPredicate {
		return filterTraversers(traversers, func(t *Traverser) bool {
			return matchPredicate(t.Element, step.Args[0].Pred)
		}), nil
	}
	ids := make(map[string]bool)
	for _, arg := range step.Args {
		ids[argString(arg)] = true
	}
	return filterTraversers(traversers, func(t *Traverser) bool {
		id := elementID(t.Element)
		return ids[id]
	}), nil
}

// execHasNot filters traversers that do NOT have the specified property key.
func execHasNot(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	key := argString(step.Args[0])
	return filterTraversers(traversers, func(t *Traverser) bool {
		props := elementProps(t.Element)
		_, ok := props[key]
		return !ok
	}), nil
}

// execValues projects element property values. With no arguments, emits all values;
// with keys, emits only the values of matching properties.
func execValues(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	keys := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		props := elementProps(t.Element)
		if len(keys) == 0 {
			for _, v := range props {
				nt := t.clone()
				nt.Element = v
				nt.Path = append(nt.Path, v)
				result = append(result, nt)
			}
		} else {
			for _, key := range keys {
				if v, ok := props[key]; ok {
					nt := t.clone()
					nt.Element = v
					nt.Path = append(nt.Path, v)
					result = append(result, nt)
				}
			}
		}
	}
	return result, nil
}

// isSlice returns true if val is a slice or array.
func isSlice(val any) bool {
	v := reflect.ValueOf(val)
	kind := v.Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

// execValueMap projects element properties as a map. Property values are wrapped in lists
// (Gremlin convention). Includes ~id and ~label metadata.
func execValueMap(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	keys := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		props := elementProps(t.Element)
		m := make(map[string]any)
		if len(keys) == 0 {
			m["~id"] = elementID(t.Element)
			m["~label"] = elementLabel(t.Element)
			for k, v := range props {
				if isSlice(v) {
					m[k] = v
				} else {
					m[k] = []any{v}
				}
			}
		} else {
			for _, key := range keys {
				if key == "id" || key == "~id" {
					m["~id"] = elementID(t.Element)
				} else if key == "label" || key == "~label" {
					m["~label"] = elementLabel(t.Element)
				} else if v, ok := props[key]; ok {
					if isSlice(v) {
						m[key] = v
					} else {
						m[key] = []any{v}
					}
				}
			}
		}
		nt := t.clone()
		nt.Element = m
		nt.Path = append(nt.Path, m)
		result = append(result, nt)
	}
	return result, nil
}

// execProperties projects individual properties as {key, value, ~id, ~label} maps.
func execProperties(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	keys := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		el := t.Element
		var idStr, labelStr string
		switch v := el.(type) {
		case *graphengine.Node:
			idStr = fmt.Sprintf("%d", v.ID)
			if len(v.Labels) > 0 {
				labelStr = v.Labels[0]
			}
		case *graphengine.Edge:
			idStr = fmt.Sprintf("%d", v.ID)
			labelStr = v.Label
		default:
			continue
		}
		props := elementProps(el)
		if len(keys) == 0 {
			for k, v := range props {
				propObj := map[string]any{"key": k, "value": v, "~id": idStr, "~label": labelStr}
				nt := t.clone()
				nt.Element = propObj
				nt.Path = append(nt.Path, propObj)
				result = append(result, nt)
			}
		} else {
			for _, key := range keys {
				if v, ok := props[key]; ok {
					propObj := map[string]any{"key": key, "value": v, "~id": idStr, "~label": labelStr}
					nt := t.clone()
					nt.Element = propObj
					nt.Path = append(nt.Path, propObj)
					result = append(result, nt)
				}
			}
		}
	}
	return result, nil
}

// execPropertyMap projects element properties as a flat map (no list wrapping).
func execPropertyMap(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	keys := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		props := elementProps(t.Element)
		m := make(map[string]any)
		if len(keys) == 0 {
			for k, v := range props {
				m[k] = v
			}
		} else {
			for _, key := range keys {
				if v, ok := props[key]; ok {
					m[key] = v
				}
			}
		}
		nt := t.clone()
		nt.Element = m
		nt.Path = append(nt.Path, m)
		result = append(result, nt)
	}
	return result, nil
}

// execId projects the element ID.
func execId(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		var id any
		switch v := t.Element.(type) {
		case *graphengine.Node:
			id = v.ID
		case *graphengine.Edge:
			id = v.ID
		default:
			id = elementID(t.Element)
		}
		nt := t.clone()
		nt.Element = id
		nt.Path = append(nt.Path, id)
		result = append(result, nt)
	}
	return result, nil
}

// execLabel projects the element label (joined with "::" for multi-label nodes).
func execLabel(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		var lbl string
		switch v := t.Element.(type) {
		case *graphengine.Node:
			lbl = strings.Join(v.Labels, "::")
		case *graphengine.Edge:
			lbl = v.Label
		default:
			lbl = elementLabel(t.Element)
		}
		nt := t.clone()
		nt.Element = lbl
		nt.Path = append(nt.Path, lbl)
		result = append(result, nt)
	}
	return result, nil
}

// execProperty sets or updates element properties. Supports Cardinality enum (single, set, list).
// For set/list cardinality, existing values are accumulated.
func execProperty(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: property requires a GraphWriter")
	}

	args := step.Args
	startIdx := 0
	cardinality := "single"

	if len(args) > 0 && args[0].Kind == ArgEnum && (args[0].Enum.Value == "single" || args[0].Enum.Value == "set" || args[0].Enum.Value == "list") {
		cardinality = args[0].Enum.Value
		startIdx = 1
	}

	if len(args) < startIdx+2 {
		return traversers, nil
	}

	key := argString(args[startIdx])
	val := resolveArgValue(args[startIdx+1])

	for _, t := range traversers {
		switch el := t.Element.(type) {
		case *graphengine.Node:
			switch cardinality {
			case "single":
				if err := ec.Writer.UpdateNode(el.ID, graphengine.Props{key: val}); err != nil {
					return nil, err
				}
			case "set":
				existingProps := el.Props
				if existingVals, ok := existingValsAsList(existingProps[key]); ok {
					if !containsVal(existingVals, val) {
						if err := ec.Writer.UpdateNode(el.ID, graphengine.Props{key: append(existingVals, val)}); err != nil {
							return nil, err
						}
					}
				} else {
					if err := ec.Writer.UpdateNode(el.ID, graphengine.Props{key: []any{existingProps[key], val}}); err != nil {
						return nil, err
					}
				}
			case "list":
				existingProps := el.Props
				if existingVals, ok := existingValsAsList(existingProps[key]); ok {
					if err := ec.Writer.UpdateNode(el.ID, graphengine.Props{key: append(existingVals, val)}); err != nil {
						return nil, err
					}
				} else {
					if err := ec.Writer.UpdateNode(el.ID, graphengine.Props{key: []any{existingProps[key], val}}); err != nil {
						return nil, err
					}
				}
			}
			if refreshed, err := ec.Reader.GetNode(el.ID); err == nil {
				t.Element = refreshed
			}
		case *graphengine.Edge:
			switch cardinality {
			case "single":
				if err := ec.Writer.UpdateEdge(el.ID, graphengine.Props{key: val}); err != nil {
					return nil, err
				}
			case "set":
				existingProps := el.Props
				if existingVals, ok := existingValsAsList(existingProps[key]); ok {
					if !containsVal(existingVals, val) {
						if err := ec.Writer.UpdateEdge(el.ID, graphengine.Props{key: append(existingVals, val)}); err != nil {
							return nil, err
						}
					}
				} else {
					if err := ec.Writer.UpdateEdge(el.ID, graphengine.Props{key: []any{existingProps[key], val}}); err != nil {
						return nil, err
					}
				}
			case "list":
				existingProps := el.Props
				if existingVals, ok := existingValsAsList(existingProps[key]); ok {
					if err := ec.Writer.UpdateEdge(el.ID, graphengine.Props{key: append(existingVals, val)}); err != nil {
						return nil, err
					}
				} else {
					if err := ec.Writer.UpdateEdge(el.ID, graphengine.Props{key: []any{existingProps[key], val}}); err != nil {
						return nil, err
					}
				}
			}
			if refreshed, err := ec.Reader.GetEdge(el.ID); err == nil {
				t.Element = refreshed
			}
		}
	}
	return traversers, nil
}

// existingValsAsList attempts to extract a []any from a property value. Used by
// set/list cardinality property handling.
func existingValsAsList(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	if l, ok := v.([]any); ok {
		return l, true
	}
	return nil, false
}

// execDrop removes all traversed elements from the graph (requires GraphWriter).
func execDrop(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: drop requires a GraphWriter")
	}

	for _, t := range traversers {
		switch el := t.Element.(type) {
		case *graphengine.Node:
			if err := ec.Writer.DeleteNode(el.ID); err != nil {
				return nil, err
			}
		case *graphengine.Edge:
			if err := ec.Writer.DeleteEdge(el.ID); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// execCount returns the number of traversers. Supports global (default) and local scope.
func execCount(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	scope := "global"
	if len(step.Args) > 0 && step.Args[0].Kind == ArgEnum && step.Args[0].Enum.Value == "local" {
		scope = "local"
	}

	if scope == "local" {
		var result []*Traverser
		for _, t := range traversers {
			nt := t.clone()
			switch el := t.Element.(type) {
			case []any:
				nt.Element = int64(len(el))
			case map[string]any:
				nt.Element = int64(len(el))
			case string:
				nt.Element = int64(1)
			default:
				nt.Element = int64(1)
			}
			result = append(result, nt)
		}
		return result, nil
	}

	return []*Traverser{newTraverser(int64(len(traversers)))}, nil
}

// execIdentity is a no-op step that passes traversers through unchanged.
func execIdentity(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

// execKeys projects the property keys of each traverser's element.
func execKeys(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		props := elementProps(t.Element)
		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}
		nt := t.clone()
		nt.Element = keys
		nt.Path = append(nt.Path, keys)
		result = append(result, nt)
	}
	return result, nil
}

// execSample returns a random sample of n traversers.
func execSample(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	n := int(argInt(step.Args[0]))
	if n <= 0 {
		return nil, nil
	}
	if n >= len(traversers) {
		return traversers, nil
	}
	cloned := make([]*Traverser, len(traversers))
	copy(cloned, traversers)
	rand.Shuffle(len(cloned), func(i, j int) {
		cloned[i], cloned[j] = cloned[j], cloned[i]
	})
	return cloned[:n], nil
}

// execSack projects the sack value from each traverser.
func execSack(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		nt := t.clone()
		nt.Element = t.Sack
		nt.Path = append(nt.Path, t.Sack)
		result = append(result, nt)
	}
	return result, nil
}
