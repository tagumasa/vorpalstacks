package gremlinparser

import (
	"fmt"
	"sort"
	"strings"

	"vorpalstacks/internal/core/storage/graphengine"
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
	if len(step.Args) == 0 && len(step.Modulators) == 0 {
		return traversers, nil
	}

	if len(step.Args) == 0 && len(step.Modulators) > 0 {
		for _, mod := range step.Modulators {
			if mod.Name == "by" {
				traversers = filterTraversers(traversers, func(t *Traverser) bool {
					val := computeByValue(ec, t, mod)
					return val != nil && !isZeroValue(val)
				})
			}
		}
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
