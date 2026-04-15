// Filter, branch, and projection step executors.
//
// Registers the following steps:
//   - select, as: label-based side-effect storage and retrieval
//   - where, filter: traversal-based filtering
//   - limit, range, skip: result set slicing
//   - dedup: deduplication by element, label, or by() modulator
//   - order, by: sorting (by() provides sort keys)
//   - group, groupCount: aggregation into maps
//   - path: path history projection
//   - union, coalesce, choose, optional: branching
//   - is, not, and, or: boolean filtering
//   - repeat, emit, until, times: looping/gremlin traversal
//   - local: scoped sub-traversal
//   - simplePath: cycle detection in paths
//   - tail: last-n elements
//   - constant, unfold, fold: element transformation
//   - map, flatMap: projection and flattening
//   - aggregate, cap: side-effect storage and retrieval
//   - from, to: edge source/target resolution for addE
//   - mergeV: upsert vertex (match or create)
//   - option, project, elementMap: advanced projection
//   - mean, sum, min, max: numeric aggregation
//   - inject, coin: element injection and random sampling

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

// execAs labels the current element in each traverser's tags for later retrieval via select().
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

// execSelect retrieves previously labelled elements from traverser tags.
// Supports Pop enum (all, first) and by() modulators for projection.
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

// execWhere filters traversers using a nested traversal, a label-predicate pair,
// or by() modulators that extract values for truthiness checking.
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

// execFilter filters traversers using a nested traversal (retains those producing results).
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

// execLimit retains only the first n traversers.
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

// execRange retains traversers from index lo (inclusive) to hi (exclusive).
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

// execSkip discards the first n traversers.
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

// execDedup removes duplicate traversers based on element identity, label, or by() modulator.
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

// dedupKey produces a deduplication key string for a traverser element.
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

// execOrder sorts traversers in-place by one or more by() modulators. Supports asc/desc via enum args.
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

// execBy is a no-op when used as a standalone step. As a modulator, it provides
// sort keys to order() and group keys to group()/groupCount()/dedup().
func execBy(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	_ = ec
	_ = traversers
	_ = step
	return traversers, nil
}

// execGroup groups traversers into a map keyed by the first by() modulator,
// with values collected from the second by() modulator (or raw elements).
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

// execGroupCount groups traversers and counts occurrences per key.
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

// execPath projects the full traversal path history as the element.
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

// execUnion merges the results of all nested traversal arguments into a single stream.
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

// execCoalesce returns the results of the first nested traversal that produces output for each traverser.
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

// execChoose branches each traverser based on a condition traversal (args[0]).
// If the condition produces results, args[1] traversal is executed; otherwise args[2].
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

// execOptional applies a nested traversal to each traverser, keeping the original
// traverser if the nested traversal yields no results.
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

// execIs filters traversers by a predicate or exact value match.
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

// execNot retains traversers for which the nested traversal produces no results.
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

// execAnd retains traversers for which ALL nested traversal arguments produce results.
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

// execOr retains traversers for which ANY nested traversal argument produces results.
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

// execRepeat iterates a nested traversal with support for emit(), until(), and times()
// modulators. The nested traversal is executed repeatedly on each traverser until a
// stop condition is met (until modulator, max iterations from times, or exhaustion).
// When emit is active, traversers are yielded at each iteration.
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

	var untilTrav *Traversal
	var untilHasTrav bool
	if hasUntil && len(untilMod.Args) > 0 && untilMod.Args[0].Kind == ArgNestedTraversal {
		untilTrav = untilMod.Args[0].Trav
		untilHasTrav = true
	}

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

// execEmit is a no-op when used standalone. As a modulator on repeat(), it causes
// traversers to be emitted at each iteration rather than only at the end.
func execEmit(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

// execUntil is a no-op when used standalone. As a modulator on repeat(), it defines
// the stop condition for loop iteration.
func execUntil(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

// execTimes is a no-op when used standalone. As a modulator on repeat(), it limits
// the number of loop iterations.
func execTimes(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return traversers, nil
}

// execLocal executes a nested traversal scoped to each individual traverser,
// allowing per-element sub-traversals without polluting the outer traversal stream.
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

// execSimplePath filters out traversers whose path history contains duplicate elements,
// effectively removing cycles from traversal paths.
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

// execTail returns the last n traversers from the stream, or the final element if no count is given.
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

// execConstant replaces each traverser's element with the specified constant value,
// useful for injecting fixed values mid-traversal.
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

// execUnfold iterates over traverser elements that are lists or maps, emitting one traverser
// per item. Non-iterable elements pass through unchanged.
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

// execFold aggregates all traverser elements into a single list traverser.
// The inverse of unfold(); an empty stream yields an empty list.
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

// execMap applies a nested traversal to each traverser and collects the results as a list element,
// producing one output traverser per input traverser.
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

// execFlatMap applies a nested traversal to each traverser and flattens all results into the stream.
// Unlike execMap, each nested result becomes a separate traverser in the output.
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

// execAggregate collects all traverser elements into a named side-effect collection,
// which can later be retrieved with cap(). The original traversers pass through unchanged.
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

// execCap retrieves a named side-effect collection (populated by aggregate() or group())
// and emits it as a single traverser. Returns an empty list if the key does not exist.
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

// execFrom resolves the source vertex for addE() from a labelled tag or nested traversal.
// Stores the resolved node ID in the "__from" tag for the subsequent addE step to consume.
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

// execTo resolves the target vertex for addE() from a labelled tag or nested traversal.
// Stores the resolved node ID in the "__to" tag for the subsequent addE step to consume.
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

// execMergeV upserts a vertex: returns an existing match (by label + properties) or creates
// a new one with the given properties. Requires a GraphWriter in the execution context.
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

// execOption registers a conditional branch pair (condition traversal + result traversal)
// as a side-effect for use with choose() or coalesce().
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

// optionEntry pairs a condition traversal with the traversal to execute when the condition matches.
type optionEntry struct {
	condition *Traversal
	traversal *Traversal
}

// getSideEffectOptions retrieves all registered option entries from the execution context's
// "__options" side-effect key.
func getSideEffectOptions(ec *ExecContext) []optionEntry {
	opts, _ := ec.SideEffects["__options"].([]optionEntry)
	return opts
}

// execProject projects traverser elements into a map using the specified keys, looking up
// values from tags or element properties. Optional by() modulators control value extraction.
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

// execElementMap converts each traverser's element (Node or Edge) into a map containing
// its label, id, and properties. Optional key arguments restrict which fields are included.
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

// sliceContains reports whether item appears in the string slice s.
func sliceContains(s []string, item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

// execMean computes the arithmetic mean of all numeric traverser elements,
// returning a single traverser with the result. Non-numeric elements are skipped.
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

// execSum computes the sum of all numeric traverser elements, returning a single
// traverser with the result. Non-numeric elements are skipped.
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

// execMin finds the minimum element across all traversers using numeric comparison,
// returning a single traverser. Returns nil for an empty input.
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

// execMax finds the maximum element across all traversers using numeric comparison,
// returning a single traverser. Returns nil for an empty input.
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

// execInject inserts the specified values at the start of the traverser stream,
// before the existing traversers. Useful for seeding a traversal with initial values.
func execInject(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, arg := range step.Args {
		result = append(result, newTraverser(resolveArgValue(arg)))
	}
	result = append(result, traversers...)
	return result, nil
}

// execCoin randomly filters traversers, retaining each with the given probability (default 0.5).
// Produces a probabilistic sample of the traverser stream.
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

// executeNestedTraversal runs a sub-traversal in the context of a parent traverser.
// Handles two cases: (1) source steps (V, E, addV, etc.) start fresh, and
// (2) mid-traversal steps use the parent traverser as their starting element.
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

// getModulators returns all modulators on a step with the given name (e.g. "by").
func getModulators(step Step, name string) []Step {
	var result []Step
	for _, m := range step.Modulators {
		if m.Name == name {
			result = append(result, m)
		}
	}
	return result
}

// getFirstModulator returns the first modulator on a step with the given name.
func getFirstModulator(step Step, name string) (Step, bool) {
	for _, m := range step.Modulators {
		if m.Name == name {
			return m, true
		}
	}
	return Step{}, false
}

// computeByValue resolves the value that a by() modulator extracts from a traverser.
// Supports nested traversals, property key lookups (including ~id and ~label), and
// identity (no args returns the element as-is).
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

// computeByValueAsString resolves a by() modulator value and formats it as a string.
// Used by dedup() and groupCount() for map key computation.
func computeByValueAsString(ec *ExecContext, t *Traverser, byMod Step) string {
	v := computeByValue(ec, t, byMod)
	return fmt.Sprintf("%v", v)
}

// valuesEqual compares two values for equality, handling numeric coercion (both values
// convertible to float64 are compared numerically) and string representation fallback.
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
