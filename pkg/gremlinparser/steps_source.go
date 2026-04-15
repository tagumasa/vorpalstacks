// Source step executors: graph element access, mutation, and property operations.
//
// Registers the following steps:
//   - V, E: source steps for iterating all vertices/edges or looking up by ID
//   - addV, addE: vertex/edge creation (requires GraphWriter)
//   - out, in, both: vertex-to-vertex traversal via edges
//   - outE, inE, bothE: vertex-to-edge traversal
//   - outV, inV, otherV: edge-to-vertex traversal
//   - has, hasLabel, hasId, hasNot: element filtering
//   - values, valueMap, properties, propertyMap: property projection
//   - id, label: element metadata access
//   - property: property mutation (supports single/set/list cardinality)
//   - drop: element deletion (requires GraphWriter)
//   - count, identity, keys, sample, sack, match, math: miscellaneous steps

package gremlinparser

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/pkg/graphengine"
)

func init() {
	RegisterStep("V", execV)
	RegisterStep("E", execE)
	RegisterStep("addV", execAddV)
	RegisterStep("addE", execAddE)
	RegisterStep("out", execOut)
	RegisterStep("in", execIn)
	RegisterStep("both", execBoth)
	RegisterStep("outE", execOutE)
	RegisterStep("inE", execInE)
	RegisterStep("bothE", execBothE)
	RegisterStep("outV", execOutV)
	RegisterStep("inV", execInV)
	RegisterStep("otherV", execOtherV)
	RegisterStep("has", execHas)
	RegisterStep("hasLabel", execHasLabel)
	RegisterStep("hasId", execHasId)
	RegisterStep("hasNot", execHasNot)
	RegisterStep("values", execValues)
	RegisterStep("valueMap", execValueMap)
	RegisterStep("properties", execProperties)
	RegisterStep("propertyMap", execPropertyMap)
	RegisterStep("id", execId)
	RegisterStep("label", execLabel)
	RegisterStep("property", execProperty)
	RegisterStep("drop", execDrop)
	RegisterStep("count", execCount)
	RegisterStep("identity", execIdentity)
	RegisterStep("keys", execKeys)
	RegisterStep("sample", execSample)
	RegisterStep("sack", execSack)
	RegisterStep("match", execMatch)
	RegisterStep("math", execMath)
}

// execV is the vertex source step. With no arguments, iterates all vertices;
// with arguments, looks up vertices by ID.
func execV(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		var traversers []*Traverser
		err := ec.Reader.ForEachNode(func(n *graphengine.Node) error {
			traversers = append(traversers, newTraverser(n))
			return nil
		})
		return traversers, err
	}

	var traversers []*Traverser
	for _, arg := range step.Args {
		id, err := toNodeID(arg)
		if err != nil {
			continue
		}
		n, err := ec.Reader.GetNode(id)
		if err != nil {
			continue
		}
		traversers = append(traversers, newTraverser(n))
	}
	return traversers, nil
}

// execE is the edge source step. With no arguments, iterates all edges;
// with arguments, looks up edges by ID.
func execE(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		var traversers []*Traverser
		err := ec.Reader.ForEachEdge(func(e *graphengine.Edge) error {
			traversers = append(traversers, newTraverser(e))
			return nil
		})
		return traversers, err
	}

	var traversers []*Traverser
	for _, arg := range step.Args {
		id, err := toEdgeID(arg)
		if err != nil {
			continue
		}
		e, err := ec.Reader.GetEdge(id)
		if err != nil {
			continue
		}
		traversers = append(traversers, newTraverser(e))
	}
	return traversers, nil
}

// execAddV creates a new vertex with the given label(s).
func execAddV(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: addV requires a GraphWriter")
	}

	var labels []string
	if len(step.Args) > 0 {
		labels = []string{argString(step.Args[0])}
	}

	id, err := ec.Writer.AddNode(labels, nil)
	if err != nil {
		return nil, err
	}
	n, err := ec.Reader.GetNode(id)
	if err != nil {
		return nil, err
	}
	return []*Traverser{newTraverser(n)}, nil
}

// execAddE creates a new edge with the given label. The source vertex comes from
// the incoming traverser (or a from() modulator), and the target from a to() modulator.
func execAddE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: addE requires a GraphWriter")
	}

	if len(step.Args) == 0 {
		return nil, fmt.Errorf("gremlin: addE requires a label argument")
	}
	label := argString(step.Args[0])

	if len(traversers) == 0 {
		return nil, fmt.Errorf("gremlin: addE requires a source traverser")
	}

	tags := traversers[0].Tags
	fromID, hasFrom := tags["__from"]
	toID, hasTo := tags["__to"]

	if !hasFrom {
		sourceNode, ok := asNode(traversers[0])
		if ok {
			for _, t := range traversers {
				t.Tags["__from"] = sourceNode.ID
			}
			fromID = sourceNode.ID
			hasFrom = true
		}
	}

	if !hasTo {
		if len(traversers) > 1 {
			targetNode, ok := asNode(traversers[1])
			if ok {
				for _, t := range traversers {
					t.Tags["__to"] = targetNode.ID
				}
				toID = targetNode.ID
				hasTo = true
			}
		} else {
			return nil, fmt.Errorf("gremlin: addE requires a to() step when there is only one traverser")
		}
	}

	if !hasFrom || !hasTo {
		return nil, fmt.Errorf("gremlin: addE requires both from() and to() steps")
	}

	fromNodeID, ok := fromID.(graphengine.NodeID)
	if !ok {
		return nil, fmt.Errorf("gremlin: addE from ID is not a valid NodeID")
	}
	toNodeIDVal, ok := toID.(graphengine.NodeID)
	if !ok {
		return nil, fmt.Errorf("gremlin: addE to ID is not a valid NodeID")
	}

	var result []*Traverser
	for _, t := range traversers {
		edgeID, err := ec.Writer.AddEdge(fromNodeID, toNodeIDVal, label, nil)
		if err != nil {
			return nil, err
		}

		e, err := ec.Reader.GetEdge(edgeID)
		if err != nil {
			return nil, err
		}

		nt := t.clone()
		nt.Element = e
		nt.Path = append(nt.Path, e)
		result = append(result, nt)
	}

	return result, nil
}

// execOut traverses outgoing edges to adjacent vertices, filtered by optional edge labels.
func execOut(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			continue
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Outgoing, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			target, err := ec.Reader.GetNode(e.To)
			if err != nil {
				continue
			}
			nt := t.clone()
			nt.Element = target
			nt.Path = append(nt.Path, target)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execIn traverses incoming edges to adjacent vertices, filtered by optional edge labels.
func execIn(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			continue
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Incoming, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			target, err := ec.Reader.GetNode(e.From)
			if err != nil {
				continue
			}
			nt := t.clone()
			nt.Element = target
			nt.Path = append(nt.Path, target)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execBoth traverses both incoming and outgoing edges to adjacent vertices.
func execBoth(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			continue
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Both, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			var targetID graphengine.NodeID
			if e.To == n.ID {
				targetID = e.From
			} else {
				targetID = e.To
			}
			target, err := ec.Reader.GetNode(targetID)
			if err != nil {
				continue
			}
			nt := t.clone()
			nt.Element = target
			nt.Path = append(nt.Path, target)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execOutE traverses outgoing edges from each vertex, filtered by optional edge labels.
func execOutE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: outE() requires vertex input, got %T", t.Element)
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Outgoing, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			nt := t.clone()
			nt.Element = e
			nt.Path = append(nt.Path, e)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execInE traverses incoming edges to each vertex, filtered by optional edge labels.
func execInE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: inE() requires vertex input, got %T", t.Element)
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Incoming, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			nt := t.clone()
			nt.Element = e
			nt.Path = append(nt.Path, e)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execBothE traverses both incoming and outgoing edges from each vertex.
func execBothE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: bothE() requires vertex input, got %T", t.Element)
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Both, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			nt := t.clone()
			nt.Element = e
			nt.Path = append(nt.Path, e)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execOutV returns the source vertex of each edge.
func execOutV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		e, ok := asEdge(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: outV() requires edge input, got %T", t.Element)
		}
		n, err := ec.Reader.GetNode(e.From)
		if err != nil {
			continue
		}
		nt := t.clone()
		nt.Element = n
		nt.Path = append(nt.Path, n)
		result = append(result, nt)
	}
	return result, nil
}

// execInV returns the target vertex of each edge.
func execInV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		e, ok := asEdge(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: inV() requires edge input, got %T", t.Element)
		}
		n, err := ec.Reader.GetNode(e.To)
		if err != nil {
			continue
		}
		nt := t.clone()
		nt.Element = n
		nt.Path = append(nt.Path, n)
		result = append(result, nt)
	}
	return result, nil
}

// execOtherV returns the vertex at the opposite end of each edge from the traverser's path.
func execOtherV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		e, ok := asEdge(t)
		if !ok {
			continue
		}
		var otherID graphengine.NodeID
		if len(t.Path) > 1 {
			if prevNode, ok := t.Path[len(t.Path)-2].(*graphengine.Node); ok {
				if e.From == prevNode.ID {
					otherID = e.To
				} else {
					otherID = e.From
				}
			}
		}
		if otherID == 0 {
			continue
		}
		n, err := ec.Reader.GetNode(otherID)
		if err != nil {
			continue
		}
		nt := t.clone()
		nt.Element = n
		nt.Path = append(nt.Path, n)
		result = append(result, nt)
	}
	return result, nil
}

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

// elementOrValueEqual compares two values for equality, handling node and edge identity comparison.
func elementOrValueEqual(a, b any) bool {
	if a == b {
		return true
	}
	an, aok := a.(*graphengine.Node)
	bn, bok := b.(*graphengine.Node)
	if aok && bok {
		return an.ID == bn.ID
	}
	ae, aok := a.(*graphengine.Edge)
	be, bok := b.(*graphengine.Edge)
	if aok && bok {
		return ae.ID == be.ID
	}
	return false
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

// execMatch is a stub for the match() step (not yet implemented).
func execMatch(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	return nil, fmt.Errorf("gremlin: match() step is not yet supported")
}

// execMath is a stub for the math() step (not yet implemented).
func execMath(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	return nil, fmt.Errorf("gremlin: math() step is not yet supported")
}

// singleLabelFilter returns the single label if exactly one is provided, otherwise ""
// (the graph engine returns all edges when label filter is empty).
func singleLabelFilter(labels []string) string {
	if len(labels) == 1 {
		return labels[0]
	}
	return ""
}

// filterEdgesByLabels filters edges by label when more than one label is specified.
// For a single label, the graph engine already handles filtering efficiently.
func filterEdgesByLabels(edges []*graphengine.Edge, labels []string) []*graphengine.Edge {
	if len(labels) <= 1 {
		return edges
	}
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}
	filtered := make([]*graphengine.Edge, 0, len(edges))
	for _, e := range edges {
		if labelSet[e.Label] {
			filtered = append(filtered, e)
		}
	}
	return filtered
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
