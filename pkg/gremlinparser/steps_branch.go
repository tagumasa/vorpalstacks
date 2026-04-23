package gremlinparser

import (
	"fmt"
)

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
