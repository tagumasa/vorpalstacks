package gremlinparser

import (
	"fmt"
	"math/rand"

	"vorpalstacks/pkg/graphengine"
)

func execTail(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		if len(traversers) == 0 {
			return nil, nil
		}
		return []*Traverser{traversers[len(traversers)-1]}, nil
	}
	n := int(argInt(step.Args[0]))
	if n <= 0 || len(traversers) == 0 {
		return nil, nil
	}
	if n >= len(traversers) {
		return traversers, nil
	}
	return traversers[len(traversers)-n:], nil
}

func execConstant(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	val := resolveArgValue(step.Args[0])
	var result []*Traverser
	for _, t := range traversers {
		nt := t.clone()
		nt.Element = val
		nt.Path = append(nt.Path, val)
		result = append(result, nt)
	}
	return result, nil
}

func execUnfold(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		switch el := t.Element.(type) {
		case []any:
			for _, item := range el {
				nt := t.clone()
				nt.Element = item
				nt.Path = append(nt.Path, item)
				result = append(result, nt)
			}
		case map[string]any:
			for _, v := range el {
				nt := t.clone()
				nt.Element = v
				nt.Path = append(nt.Path, v)
				result = append(result, nt)
			}
		default:
			result = append(result, t)
		}
	}
	return result, nil
}

func execFold(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(traversers) == 0 {
		return []*Traverser{newTraverser([]any{})}, nil
	}
	list := make([]any, len(traversers))
	for i, t := range traversers {
		list[i] = t.Element
	}
	return []*Traverser{newTraverser(list)}, nil
}

func execMap(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 || step.Args[0].Kind != ArgNestedTraversal {
		return traversers, nil
	}
	var result []*Traverser
	for _, t := range traversers {
		nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
		if err != nil || len(nested) == 0 {
			continue
		}
		vals := make([]any, len(nested))
		for i, n := range nested {
			vals[i] = n.Element
		}
		nt := t.clone()
		nt.Element = vals
		result = append(result, nt)
	}
	return result, nil
}

func execFlatMap(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 || step.Args[0].Kind != ArgNestedTraversal {
		return traversers, nil
	}
	var result []*Traverser
	for _, t := range traversers {
		nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
		if err != nil || len(nested) == 0 {
			continue
		}
		result = append(result, nested...)
	}
	return result, nil
}

func execAggregate(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	sideEffectKey := argString(step.Args[0])

	if ec.SideEffects[sideEffectKey] == nil {
		ec.SideEffects[sideEffectKey] = []any{}
	}

	collection, ok := ec.SideEffects[sideEffectKey].([]any)
	if !ok {
		collection = []any{}
	}
	for _, t := range traversers {
		collection = append(collection, t.Element)
	}
	ec.SideEffects[sideEffectKey] = collection
	return traversers, nil
}

func execCap(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	sideEffectKey := argString(step.Args[0])

	val, ok := ec.SideEffects[sideEffectKey]
	if !ok {
		return []*Traverser{newTraverser([]any{})}, nil
	}
	return []*Traverser{newTraverser(val)}, nil
}

func execFrom(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	if len(traversers) == 0 {
		return traversers, nil
	}
	if step.Args[0].Kind == ArgString {
		label := argString(step.Args[0])
		result := make([]*Traverser, len(traversers))
		for i, t := range traversers {
			val, ok := t.Tags[label]
			if !ok {
				result[i] = t
				continue
			}
			if n, ok := val.(*graphengine.Node); ok {
				nt := t.clone()
				nt.Tags["__from"] = n.ID
				result[i] = nt
			} else {
				result[i] = t
			}
		}
		return result, nil
	}
	if step.Args[0].Kind == ArgNestedTraversal {
		result := make([]*Traverser, len(traversers))
		for i, t := range traversers {
			nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
			if err != nil || len(nested) == 0 {
				result[i] = t
				continue
			}
			if n, ok := asNode(nested[0]); ok {
				nt := t.clone()
				nt.Tags["__from"] = n.ID
				result[i] = nt
			} else {
				result[i] = t
			}
		}
		return result, nil
	}
	return traversers, nil
}

func execTo(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	if len(traversers) == 0 {
		return traversers, nil
	}
	if step.Args[0].Kind == ArgString {
		label := argString(step.Args[0])
		result := make([]*Traverser, len(traversers))
		for i, t := range traversers {
			val, ok := t.Tags[label]
			if !ok {
				result[i] = t
				continue
			}
			if n, ok := val.(*graphengine.Node); ok {
				nt := t.clone()
				nt.Tags["__to"] = n.ID
				result[i] = nt
			} else {
				result[i] = t
			}
		}
		return result, nil
	}
	if step.Args[0].Kind == ArgNestedTraversal {
		result := make([]*Traverser, len(traversers))
		for i, t := range traversers {
			nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
			if err != nil || len(nested) == 0 {
				result[i] = t
				continue
			}
			if n, ok := asNode(nested[0]); ok {
				nt := t.clone()
				nt.Tags["__to"] = n.ID
				result[i] = nt
			} else {
				result[i] = t
			}
		}
		return result, nil
	}
	return traversers, nil
}

func execMergeV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: mergeV requires a GraphWriter")
	}

	var props map[string]any
	if len(step.Args) > 0 && step.Args[0].Kind == ArgMap {
		props = make(map[string]any)
		for _, entry := range step.Args[0].Map {
			props[resolveArgValue(entry.Key).(string)] = resolveArgValue(entry.Value)
		}
	}

	var labels []string
	if label, ok := props["~label"]; ok {
		labels = []string{fmt.Sprintf("%v", label)}
		delete(props, "~label")
	}

	searchLabel := ""
	if len(labels) > 0 {
		searchLabel = labels[0]
	}

	if searchLabel != "" && len(props) > 0 {
		var existingIDs []graphengine.NodeID
		_ = ec.Reader.ForEachNode(func(n *graphengine.Node) error {
			if !matchLabels(n.Labels, []string{searchLabel}) {
				return nil
			}
			allMatch := true
			for k, v := range props {
				if nv, ok := n.Props[k]; !ok || !valuesEqual(nv, v) {
					allMatch = false
					break
				}
			}
			if allMatch {
				existingIDs = append(existingIDs, n.ID)
			}
			return nil
		})
		if len(existingIDs) > 0 {
			n, err := ec.Reader.GetNode(existingIDs[0])
			if err == nil {
				return []*Traverser{newTraverser(n)}, nil
			}
		}
	}

	n, err := ec.Writer.AddNode(labels, graphengine.Props(props))
	if err != nil {
		return nil, err
	}
	return []*Traverser{newTraverser(n)}, nil
}

func execOption(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) > 0 && step.Args[0].Kind == ArgNestedTraversal {
		if len(step.Args) > 1 && step.Args[1].Kind == ArgNestedTraversal {
			ec.SideEffects["__options"] = append(getSideEffectOptions(ec), optionEntry{
				condition: step.Args[0].Trav,
				traversal: step.Args[1].Trav,
			})
		}
	}
	return traversers, nil
}

type optionEntry struct {
	condition *Traversal
	traversal *Traversal
}

func getSideEffectOptions(ec *ExecContext) []optionEntry {
	opts, _ := ec.SideEffects["__options"].([]optionEntry)
	return opts
}

func execProject(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}

	keys := argStrings(step.Args)
	byMods := getModulators(step, "by")
	var result []*Traverser

	for _, t := range traversers {
		m := make(map[string]any)
		for i, key := range keys {
			if val, ok := t.Tags[key]; ok {
				if len(byMods) > i {
					savedElement := t.Element
					t.Element = val
					m[key] = computeByValue(ec, t, byMods[i])
					t.Element = savedElement
				} else {
					m[key] = val
				}
			} else {
				props := elementProps(t.Element)
				if v, ok := props[key]; ok {
					m[key] = v
				} else {
					m[key] = nil
				}
			}
		}
		nt := t.clone()
		nt.Element = m
		result = append(result, nt)
	}
	return result, nil
}

func execElementMap(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	keys := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		m := make(map[string]any)
		switch el := t.Element.(type) {
		case *graphengine.Node:
			if len(keys) == 0 || sliceContains(keys, "label") || sliceContains(keys, "~label") {
				if len(el.Labels) > 0 {
					m["~label"] = el.Labels[0]
				}
			}
			if len(keys) == 0 || sliceContains(keys, "id") || sliceContains(keys, "~id") {
				m["~id"] = fmt.Sprintf("%d", el.ID)
			}
			for k, v := range el.Props {
				if len(keys) == 0 || sliceContains(keys, k) {
					m[k] = v
				}
			}
		case *graphengine.Edge:
			if len(keys) == 0 || sliceContains(keys, "label") || sliceContains(keys, "~label") {
				m["~label"] = el.Label
			}
			if len(keys) == 0 || sliceContains(keys, "id") || sliceContains(keys, "~id") {
				m["~id"] = fmt.Sprintf("%d", el.ID)
			}
			if len(keys) == 0 || sliceContains(keys, "from") || sliceContains(keys, "~from") {
				m["~from"] = fmt.Sprintf("%d", el.From)
			}
			if len(keys) == 0 || sliceContains(keys, "to") || sliceContains(keys, "~to") {
				m["~to"] = fmt.Sprintf("%d", el.To)
			}
			for k, v := range el.Props {
				if len(keys) == 0 || sliceContains(keys, k) {
					m[k] = v
				}
			}
		}
		nt := t.clone()
		nt.Element = m
		result = append(result, nt)
	}
	return result, nil
}

func sliceContains(s []string, item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func execMean(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(traversers) == 0 {
		return []*Traverser{newTraverser(float64(0))}, nil
	}
	var sum float64
	count := 0
	for _, t := range traversers {
		if f, ok := toFloat64(t.Element); ok {
			sum += f
			count++
		}
	}
	if count == 0 {
		return []*Traverser{newTraverser(float64(0))}, nil
	}
	return []*Traverser{newTraverser(sum / float64(count))}, nil
}

func execSum(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(traversers) == 0 {
		return []*Traverser{newTraverser(float64(0))}, nil
	}
	var sum float64
	for _, t := range traversers {
		if f, ok := toFloat64(t.Element); ok {
			sum += f
		}
	}
	return []*Traverser{newTraverser(sum)}, nil
}

func execMin(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(traversers) == 0 {
		return nil, nil
	}
	min := traversers[0].Element
	for _, t := range traversers[1:] {
		if compareNumeric(t.Element, min) < 0 {
			min = t.Element
		}
	}
	return []*Traverser{newTraverser(min)}, nil
}

func execMax(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(traversers) == 0 {
		return nil, nil
	}
	max := traversers[0].Element
	for _, t := range traversers[1:] {
		if compareNumeric(t.Element, max) > 0 {
			max = t.Element
		}
	}
	return []*Traverser{newTraverser(max)}, nil
}

func execInject(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, arg := range step.Args {
		result = append(result, newTraverser(resolveArgValue(arg)))
	}
	result = append(result, traversers...)
	return result, nil
}

func execCoin(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	probability := 0.5
	if step.Args[0].Kind == ArgFloat {
		probability = step.Args[0].Float
	} else if step.Args[0].Kind == ArgInt {
		probability = float64(step.Args[0].Int)
	}

	var result []*Traverser
	for _, t := range traversers {
		if rand.Float64() < probability {
			result = append(result, t)
		}
	}
	return result, nil
}

func executeNestedTraversal(ec *ExecContext, parent *Traverser, trav *Traversal) ([]*Traverser, error) {
	if len(trav.Steps) == 0 {
		return nil, nil
	}

	var current []*Traverser

	sourceStepNames := map[string]bool{"V": true, "E": true, "addV": true, "addE": true, "mergeV": true, "inject": true}
	if sourceStepNames[trav.Steps[0].Name] {
		exec, ok := stepRegistry[trav.Steps[0].Name]
		if !ok {
			return nil, fmt.Errorf("gremlin: unknown step %q", trav.Steps[0].Name)
		}
		result, err := exec(ec, nil, trav.Steps[0])
		if err != nil {
			return nil, err
		}
		current, err = applyModulators(ec, result, trav.Steps[0])
		if err != nil {
			return nil, err
		}
		for i := 1; i < len(trav.Steps); i++ {
			step := trav.Steps[i]
			exec, ok := stepRegistry[step.Name]
			if !ok {
				return nil, fmt.Errorf("gremlin: unknown step %q", step.Name)
			}
			var err error
			current, err = exec(ec, current, step)
			if err != nil {
				return nil, err
			}
			current, err = applyModulators(ec, current, step)
			if err != nil {
				return nil, err
			}
		}
	} else {
		for i, step := range trav.Steps {
			exec, ok := stepRegistry[step.Name]
			if !ok {
				return nil, fmt.Errorf("gremlin: unknown step %q", step.Name)
			}
			var err error
			var inputs []*Traverser
			if i == 0 && len(current) == 0 {
				inputs = []*Traverser{parent.clone()}
			} else {
				inputs = current
			}
			if len(inputs) == 0 {
				continue
			}
			current, err = exec(ec, inputs, step)
			if err != nil {
				return nil, err
			}
			current, err = applyModulators(ec, current, step)
			if err != nil {
				return nil, err
			}
			if len(current) > 0 {
				parent = current[len(current)-1]
			}
		}
	}

	return current, nil
}

func getModulators(step Step, name string) []Step {
	var result []Step
	for _, m := range step.Modulators {
		if m.Name == name {
			result = append(result, m)
		}
	}
	return result
}

func getFirstModulator(step Step, name string) (Step, bool) {
	for _, m := range step.Modulators {
		if m.Name == name {
			return m, true
		}
	}
	return Step{}, false
}

func computeByValue(ec *ExecContext, t *Traverser, byMod Step) any {
	if len(byMod.Args) == 0 {
		return t.Element
	}
	if byMod.Args[0].Kind == ArgNestedTraversal {
		nested, err := executeNestedTraversal(ec, t, byMod.Args[0].Trav)
		if err != nil || len(nested) == 0 {
			return nil
		}
		return nested[0].Element
	}
	if byMod.Args[0].Kind == ArgString {
		key := byMod.Args[0].Str
		if key == "~id" || key == "id" {
			return elementID(t.Element)
		}
		if key == "~label" || key == "label" {
			return elementLabel(t.Element)
		}
		props := elementProps(t.Element)
		if v, ok := props[key]; ok {
			return v
		}
		return nil
	}
	return t.Element
}

func computeByValueAsString(ec *ExecContext, t *Traverser, byMod Step) string {
	v := computeByValue(ec, t, byMod)
	return fmt.Sprintf("%v", v)
}

func valuesEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af == bf
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func isZeroValue(v any) bool {
	switch val := v.(type) {
	case nil:
		return true
	case bool:
		return !val
	case int:
		return val == 0
	case int64:
		return val == 0
	case float64:
		return val == 0
	case string:
		return val == ""
	}
	return false
}
