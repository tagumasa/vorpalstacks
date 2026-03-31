package gremlinparser

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"vorpalstacks/pkg/graphengine"
)

func init() {
	RegisterStep("select", execSelect)
	RegisterStep("as", execAs)
	RegisterStep("where", execWhere)
	RegisterStep("filter", execFilter)
	RegisterStep("limit", execLimit)
	RegisterStep("range", execRange)
	RegisterStep("skip", execSkip)
	RegisterStep("dedup", execDedup)
	RegisterStep("order", execOrder)
	RegisterStep("by", execBy)
	RegisterStep("group", execGroup)
	RegisterStep("groupCount", execGroupCount)
	RegisterStep("path", execPath)
	RegisterStep("union", execUnion)
	RegisterStep("is", execIs)
	RegisterStep("not", execNot)
	RegisterStep("and", execAnd)
	RegisterStep("or", execOr)
	RegisterStep("coalesce", execCoalesce)
	RegisterStep("choose", execChoose)
	RegisterStep("optional", execOptional)
	RegisterStep("constant", execConstant)
	RegisterStep("unfold", execUnfold)
	RegisterStep("fold", execFold)
	RegisterStep("map", execMap)
	RegisterStep("flatMap", execFlatMap)
	RegisterStep("repeat", execRepeat)
	RegisterStep("emit", execEmit)
	RegisterStep("until", execUntil)
	RegisterStep("times", execTimes)
	RegisterStep("local", execLocal)
	RegisterStep("simplePath", execSimplePath)
	RegisterStep("tail", execTail)
	RegisterStep("aggregate", execAggregate)
	RegisterStep("cap", execCap)
	RegisterStep("from", execFrom)
	RegisterStep("to", execTo)
	RegisterStep("mergeV", execMergeV)
	RegisterStep("option", execOption)
	RegisterStep("project", execProject)
	RegisterStep("elementMap", execElementMap)
	RegisterStep("mean", execMean)
	RegisterStep("sum", execSum)
	RegisterStep("min", execMin)
	RegisterStep("max", execMax)
	RegisterStep("inject", execInject)
	RegisterStep("coin", execCoin)
}

func execAs(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	label := argString(step.Args[0])
	for _, t := range traversers {
		t.Tags[label] = t.Element
	}
	return traversers, nil
}

func execSelect(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}

	pop := "all"
	remainingArgs := step.Args
	if len(remainingArgs) > 0 && remainingArgs[0].Kind == ArgEnum && remainingArgs[0].Enum != nil && remainingArgs[0].Enum.Category == "Pop" {
		pop = remainingArgs[0].Enum.Value
		remainingArgs = remainingArgs[1:]
	}

	labels := argStrings(remainingArgs)
	byMods := getModulators(step, "by")
	var result []*Traverser

	for _, t := range traversers {
		var vals []any
		for _, label := range labels {
			if val, ok := t.Tags[label]; ok {
				vals = append(vals, val)
				if pop == "first" {
					break
				}
			}
		}

		if len(byMods) > 0 && len(vals) > 0 {
			selected := vals[0]
			if len(byMods[0].Args) > 0 {
				savedElement := t.Element
				t.Element = selected
				projected := computeByValue(ec, t, byMods[0])
				t.Element = savedElement
				nt := t.clone()
				nt.Element = projected
				nt.Path = append(nt.Path, projected)
				result = append(result, nt)
			} else {
				nt := t.clone()
				nt.Element = selected
				nt.Path = append(nt.Path, selected)
				result = append(result, nt)
			}
		} else if len(vals) == 1 {
			nt := t.clone()
			nt.Element = vals[0]
			nt.Path = append(nt.Path, vals[0])
			result = append(result, nt)
		} else if len(vals) > 1 {
			nt := t.clone()
			m := make(map[string]any, len(labels))
			for i, label := range labels {
				if i < len(vals) {
					m[label] = vals[i]
				}
			}
			nt.Element = m
			nt.Path = append(nt.Path, m)
			result = append(result, nt)
		}
	}
	return result, nil
}

func execWhere(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}

	if step.Args[0].Kind == ArgNestedTraversal {
		return filterTraversers(traversers, func(t *Traverser) bool {
			nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
			if err != nil || len(nested) == 0 {
				return false
			}
			return true
		}), nil
	}

	if step.Args[0].Kind == ArgString && len(step.Args) > 1 && step.Args[1].Kind == ArgPredicate {
		label := argString(step.Args[0])
		pred := step.Args[1].Pred
		return filterTraversers(traversers, func(t *Traverser) bool {
			val, ok := t.Tags[label]
			if !ok {
				return false
			}
			return matchPredicate(val, pred)
		}), nil
	}

	return traversers, nil
}

func execFilter(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}

	if step.Args[0].Kind == ArgNestedTraversal {
		return filterTraversers(traversers, func(t *Traverser) bool {
			nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
			return err == nil && len(nested) > 0
		}), nil
	}

	return traversers, nil
}

func execLimit(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
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
	return traversers[:n], nil
}

func execRange(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) < 2 {
		return traversers, nil
	}
	lo := int(argInt(step.Args[0]))
	hi := int(argInt(step.Args[1]))
	if lo < 0 {
		lo = 0
	}
	if hi > len(traversers) {
		hi = len(traversers)
	}
	if lo >= hi {
		return nil, nil
	}
	return traversers[lo:hi], nil
}

func execSkip(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}
	n := int(argInt(step.Args[0]))
	if n <= 0 {
		return traversers, nil
	}
	if n >= len(traversers) {
		return nil, nil
	}
	return traversers[n:], nil
}

func execDedup(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	seen := make(map[string]bool)
	var result []*Traverser

	byMods := getModulators(step, "by")
	labels := argStrings(step.Args)
	for _, t := range traversers {
		var key string
		if len(byMods) > 0 {
			key = computeByValueAsString(ec, t, byMods[0])
		} else if len(labels) > 0 {
			keys := make([]string, 0, len(labels))
			for _, label := range labels {
				if val, ok := t.Tags[label]; ok {
					keys = append(keys, dedupKey(val))
				}
			}
			key = strings.Join(keys, "|")
		} else {
			key = dedupKey(t.Element)
		}
		if !seen[key] {
			seen[key] = true
			result = append(result, t)
		}
	}
	return result, nil
}

func dedupKey(val any) string {
	switch v := val.(type) {
	case *graphengine.Node:
		return fmt.Sprintf("node:%d", v.ID)
	case *graphengine.Edge:
		return fmt.Sprintf("edge:%d", v.ID)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func execOrder(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	byMods := getModulators(step, "by")
	if len(byMods) == 0 {
		return traversers, nil
	}

	sort.SliceStable(traversers, func(i, j int) bool {
		for _, bm := range byMods {
			asc := true
			if len(bm.Args) > 1 && bm.Args[1].Kind == ArgEnum {
				switch bm.Args[1].Enum.Value {
				case "desc":
					asc = false
				case "asc":
					asc = true
				}
			}
			vi := computeByValue(ec, traversers[i], bm)
			vj := computeByValue(ec, traversers[j], bm)
			if fi, ok := toFloat64(vi); ok {
				if fj, ok := toFloat64(vj); ok {
					if fi != fj {
						if asc {
							return fi < fj
						}
						return fi > fj
					}
					continue
				}
			}
			ki := computeByValueAsString(ec, traversers[i], bm)
			kj := computeByValueAsString(ec, traversers[j], bm)
			if ki != kj {
				if asc {
					return ki < kj
				}
				return ki > kj
			}
		}
		return false
	})

	return traversers, nil
}

func execBy(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	_ = ec
	_ = traversers
	_ = step
	return traversers, nil
}

func execGroup(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	byMods := getModulators(step, "by")
	groups := make(map[string][]any)
	var keys []string

	keyMod := Step{}
	if len(byMods) >= 1 {
		keyMod = byMods[0]
	}
	valMod := Step{}
	if len(byMods) >= 2 {
		valMod = byMods[1]
	}

	for _, t := range traversers {
		var key string
		if keyMod.Name != "" {
			keyElement := computeByValue(ec, t, keyMod)
			key = fmt.Sprintf("%v", keyElement)
		} else {
			key = dedupKey(t.Element)
		}
		if _, exists := groups[key]; !exists {
			keys = append(keys, key)
		}

		if valMod.Name != "" {
			val := computeByValue(ec, t, valMod)
			groups[key] = append(groups[key], val)
		} else {
			groups[key] = append(groups[key], t.Element)
		}
	}

	result := make(map[string]any)
	for _, key := range keys {
		result[key] = groups[key]
	}

	return []*Traverser{newTraverser(result)}, nil
}

func execGroupCount(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	byMods := getModulators(step, "by")
	counts := make(map[string]int64)
	var keys []string

	keyMod := Step{}
	if len(byMods) >= 1 {
		keyMod = byMods[0]
	}

	for _, t := range traversers {
		var key string
		if keyMod.Name != "" {
			key = computeByValueAsString(ec, t, keyMod)
		} else {
			key = dedupKey(t.Element)
		}
		if _, exists := counts[key]; !exists {
			keys = append(keys, key)
		}
		counts[key]++
	}

	result := make(map[string]any)
	for _, key := range keys {
		result[key] = counts[key]
	}

	return []*Traverser{newTraverser(result)}, nil
}

func execPath(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		pathCopy := make([]any, len(t.Path))
		copy(pathCopy, t.Path)
		nt := t.clone()
		nt.Element = pathCopy
		result = append(result, nt)
	}
	return result, nil
}

func execUnion(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser

	for _, t := range traversers {
		for _, arg := range step.Args {
			if arg.Kind == ArgNestedTraversal {
				nested, err := executeNestedTraversal(ec, t, arg.Trav)
				if err != nil {
					continue
				}
				result = append(result, nested...)
			}
		}
	}
	return result, nil
}

func execCoalesce(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser

	for _, t := range traversers {
		found := false
		for _, arg := range step.Args {
			if arg.Kind == ArgNestedTraversal {
				nested, err := executeNestedTraversal(ec, t, arg.Trav)
				if err != nil || len(nested) == 0 {
					continue
				}
				result = append(result, nested...)
				found = true
				break
			}
		}
		if !found {
			continue
		}
	}
	return result, nil
}

func execChoose(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser

	for _, t := range traversers {
		branchIdx := 0
		if len(step.Args) > 0 && step.Args[0].Kind == ArgNestedTraversal {
			nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
			if err != nil || len(nested) == 0 {
				branchIdx = 2
			} else {
				branchIdx = 1
			}
		} else if len(step.Args) > 0 && step.Args[0].Kind == ArgPredicate {
			if step.Args[0].Pred == nil || !matchPredicate(t.Element, step.Args[0].Pred) {
				branchIdx = 2
			} else {
				branchIdx = 1
			}
		}

		if branchIdx < len(step.Args) && step.Args[branchIdx].Kind == ArgNestedTraversal {
			nested, err := executeNestedTraversal(ec, t, step.Args[branchIdx].Trav)
			if err == nil {
				result = append(result, nested...)
			}
		} else {
			result = append(result, t)
		}
	}
	return result, nil
}

func execOptional(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser

	for _, t := range traversers {
		if len(step.Args) > 0 && step.Args[0].Kind == ArgNestedTraversal {
			nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
			if err == nil && len(nested) > 0 {
				result = append(result, nested...)
			} else {
				result = append(result, t)
			}
		} else {
			result = append(result, t)
		}
	}
	return result, nil
}

func execIs(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		return traversers, nil
	}

	if step.Args[0].Kind == ArgPredicate {
		return filterTraversers(traversers, func(t *Traverser) bool {
			return matchPredicate(t.Element, step.Args[0].Pred)
		}), nil
	}

	val := resolveArgValue(step.Args[0])
	return filterTraversers(traversers, func(t *Traverser) bool {
		return t.Element == val
	}), nil
}

func execNot(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 || step.Args[0].Kind != ArgNestedTraversal {
		return traversers, nil
	}

	var errs []error
	filtered := filterTraversers(traversers, func(t *Traverser) bool {
		nested, err := executeNestedTraversal(ec, t, step.Args[0].Trav)
		if err != nil {
			errs = append(errs, err)
			return true
		}
		return len(nested) == 0
	})
	if len(errs) > 0 {
		return filtered, errs[0]
	}
	return filtered, nil
}

func execAnd(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return filterTraversers(traversers, func(t *Traverser) bool {
		for _, arg := range step.Args {
			if arg.Kind != ArgNestedTraversal {
				continue
			}
			nested, err := executeNestedTraversal(ec, t, arg.Trav)
			if err != nil || len(nested) == 0 {
				return false
			}
		}
		return true
	}), nil
}

func execOr(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return filterTraversers(traversers, func(t *Traverser) bool {
		for _, arg := range step.Args {
			if arg.Kind != ArgNestedTraversal {
				continue
			}
			nested, err := executeNestedTraversal(ec, t, arg.Trav)
			if err == nil && len(nested) > 0 {
				return true
			}
		}
		return false
	}), nil
}

func execRepeat(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 || step.Args[0].Kind != ArgNestedTraversal {
		return traversers, nil
	}

	nested := step.Args[0].Trav
	var result []*Traverser

	_, hasEmit := getFirstModulator(step, "emit")
	untilMod, hasUntil := getFirstModulator(step, "until")
	timesMod, hasTimes := getFirstModulator(step, "times")

	stepLabel := ""
	for _, m := range step.Modulators {
		if m.Name == "as" && len(m.Args) > 0 {
			stepLabel = argString(m.Args[0])
			break
		}
	}

	maxLoops := 50
	if hasTimes && len(timesMod.Args) > 0 {
		maxLoops = int(timesMod.Args[0].Int)
		if maxLoops <= 0 {
			maxLoops = 50
		}
	}

	untilTrav, untilHasTrav := untilMod.Args[0].Trav, hasUntil && len(untilMod.Args) > 0 && untilMod.Args[0].Kind == ArgNestedTraversal

	const maxTotalTraversers = 100000

	for _, t := range traversers {
		current := []*Traverser{t}

		for i := 0; i < maxLoops; i++ {
			var next []*Traverser
			for _, ct := range current {
				if stepLabel != "" {
					ct.Loops[stepLabel] = i + 1
				}
				executed, err := executeNestedTraversal(ec, ct, nested)
				if err != nil {
					continue
				}
				next = append(next, executed...)
			}

			if hasEmit {
				for _, nt := range next {
					nt2 := nt.clone()
					result = append(result, nt2)
				}
				if len(result) > maxTotalTraversers {
					return nil, fmt.Errorf("gremlin: repeat() exceeded maximum traverser limit of %d", maxTotalTraversers)
				}
			}

			if len(next) == 0 {
				break
			}

			shouldStop := false
			if untilHasTrav {
				for _, nt := range next {
					check, err := executeNestedTraversal(ec, nt, untilTrav)
					if err == nil && len(check) > 0 {
						shouldStop = true
						break
					}
				}
			} else if hasUntil && len(untilMod.Args) > 0 && untilMod.Args[0].Kind == ArgPredicate {
				for _, nt := range next {
					if matchPredicate(nt.Element, untilMod.Args[0].Pred) {
						shouldStop = true
						break
					}
				}
			}

			if shouldStop {
				if !hasEmit {
					result = append(result, next...)
				}
				break
			}

			current = next
		}

		if !hasEmit && len(current) > 0 {
			result = append(result, current...)
		}
		if len(result) > maxTotalTraversers {
			return nil, fmt.Errorf("gremlin: repeat() exceeded maximum traverser limit of %d", maxTotalTraversers)
		}
	}

	return result, nil
}

func execEmit(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

func execUntil(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

func execTimes(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

func execLocal(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
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

func execSimplePath(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return filterTraversers(traversers, func(t *Traverser) bool {
		seen := make(map[string]bool)
		for _, p := range t.Path {
			key := dedupKey(p)
			if seen[key] {
				return false
			}
			seen[key] = true
		}
		return true
	}), nil
}

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
		ec.Reader.ForEachNode(func(n *graphengine.Node) error {
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
		for _, step := range trav.Steps {
			exec, ok := stepRegistry[step.Name]
			if !ok {
				return nil, fmt.Errorf("gremlin: unknown step %q", step.Name)
			}
			var err error
			inputs := make([]*Traverser, len(current))
			for i, t := range current {
				inputs[i] = t
			}
			if len(inputs) == 0 {
				inputs = []*Traverser{parent.clone()}
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
