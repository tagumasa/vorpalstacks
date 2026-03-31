package gremlinparser

import (
	"fmt"
	"math/rand"
	"os"
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
		id := toNodeID(arg)
		n, err := ec.Reader.GetNode(id)
		if err != nil {
			continue
		}
		traversers = append(traversers, newTraverser(n))
	}
	return traversers, nil
}

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
		id := toEdgeID(arg)
		e, err := ec.Reader.GetEdge(id)
		if err != nil {
			continue
		}
		traversers = append(traversers, newTraverser(e))
	}
	return traversers, nil
}

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

	var edgeID graphengine.EdgeID
	var err error
	edgeID, err = ec.Writer.AddEdge(fromNodeID, toNodeIDVal, label, nil)
	if err != nil {
		return nil, err
	}

	e, err := ec.Reader.GetEdge(edgeID)
	if err != nil {
		return nil, err
	}

	return []*Traverser{newTraverser(e)}, nil
}

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
			otherID = e.From
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

func execHas(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return filterTraversers(traversers, func(t *Traverser) bool {
		return matchHasStep(ec, t, step.Args)
	}), nil
}

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

func isSlice(val any) bool {
	v := reflect.ValueOf(val)
	kind := v.Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func asSlice(val any) []any {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	}
	return nil
}

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

func existingValsAsList(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	if l, ok := v.([]any); ok {
		return l, true
	}
	return nil, false
}

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

func execIdentity(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

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
			return ok && val == resolveArgValue(args[1])
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
			return ok && val == resolveArgValue(args[1])
		}
	}

	return true
}

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

func matchPredicate(val any, pred *Predicate) bool {
	if pred == nil {
		return true
	}

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
	}

	return false
}

func filterTraversers(traversers []*Traverser, fn func(*Traverser) bool) []*Traverser {
	result := make([]*Traverser, 0, len(traversers))
	for _, t := range traversers {
		if fn(t) {
			result = append(result, t)
		}
	}
	return result
}

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

func toNodeID(arg Argument) graphengine.NodeID {
	switch arg.Kind {
	case ArgInt:
		return graphengine.NodeID(arg.Int)
	case ArgFloat:
		return graphengine.NodeID(int64(arg.Float))
	case ArgString:
		s := strings.TrimPrefix(arg.Str, "n")
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gremlin: toNodeID parse error: %q: %v\n", arg.Str, err)
			return 0
		}
		return graphengine.NodeID(id)
	default:
		fmt.Fprintf(os.Stderr, "gremlin: toNodeID unsupported arg kind: %v\n", arg.Kind)
		return 0
	}
}

func toEdgeID(arg Argument) graphengine.EdgeID {
	switch arg.Kind {
	case ArgInt:
		return graphengine.EdgeID(arg.Int)
	case ArgFloat:
		return graphengine.EdgeID(int64(arg.Float))
	case ArgString:
		s := strings.TrimPrefix(arg.Str, "n")
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gremlin: toEdgeID parse error: %q: %v\n", arg.Str, err)
			return 0
		}
		return graphengine.EdgeID(id)
	default:
		fmt.Fprintf(os.Stderr, "gremlin: toEdgeID unsupported arg kind: %v\n", arg.Kind)
		return 0
	}
}

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

func containsStr(s, substr string) bool {
	return strings.Contains(s, substr)
}

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

func execMatch(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	return nil, fmt.Errorf("gremlin: match() step is not yet supported")
}

func execMath(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	return nil, fmt.Errorf("gremlin: math() step is not yet supported")
}

func singleLabelFilter(labels []string) string {
	if len(labels) == 1 {
		return labels[0]
	}
	return ""
}

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

func containsVal(vals []any, val any) bool {
	for _, v := range vals {
		if elementOrValueEqual(v, val) {
			return true
		}
	}
	return false
}
