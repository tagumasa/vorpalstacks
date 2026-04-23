// Cypher full query execution pipeline: WITH, UNWIND, write clauses, and multi-hop matching.
//
// Based on goraphdb/cypher_exec.go match binding and multi-pattern execution.
// Handles the complete query lifecycle including WITH projection, UNWIND expansion,
// SET/DELETE/REMOVE write operations, comma-separated MATCH patterns, and
// variable-length path binding.

package cypherparser

import (
	"context"
	"fmt"

	"vorpalstacks/pkg/graphengine"
)

func executeFullQuery(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	var bindings []map[string]any

	if q.Unwind != nil && q.With == nil {
		val, err := evalExpr(&EvalContext{}, &q.Unwind.Expr)
		if err != nil {
			return nil, fmt.Errorf("cypher: UNWIND expression eval failed: %w", err)
		}
		list, ok := val.([]any)
		if !ok {
			return nil, fmt.Errorf("cypher: UNWIND requires a list, got %T", val)
		}
		for _, item := range list {
			bindings = append(bindings, map[string]any{q.Unwind.Var: item})
		}
	} else {
		bindings = []map[string]any{{}}
	}

	hasMatch := len(q.Match.Pattern.Nodes) > 0
	if hasMatch {
		var err error
		bindings, err = matchBindings(ctx, reader, q, bindings)
		if err != nil {
			return nil, err
		}
	}

	for i := 1; i < len(q.Matches); i++ {
		subQ := &CypherQuery{Match: q.Matches[i]}
		var err error
		bindings, err = matchBindings(ctx, reader, subQ, bindings)
		if err != nil {
			return nil, err
		}
		if len(bindings) == 0 {
			break
		}
	}

	if q.Where != nil && q.With == nil {
		filtered := make([]map[string]any, 0, len(bindings))
		for _, b := range bindings {
			ok, err := evalBool(&EvalContext{Bindings: b}, q.Where)
			if err != nil {
				return nil, err
			}
			if ok {
				filtered = append(filtered, b)
			}
		}
		bindings = filtered
	}

	if q.With != nil {
		var err error
		bindings, err = applyWith(q.With, q.With.Where, bindings)
		if err != nil {
			return nil, err
		}
	}

	if q.Unwind != nil && q.With != nil {
		var expanded []map[string]any
		for _, binding := range bindings {
			unwindVal := evalUnwindExpr(&q.Unwind.Expr, binding)
			if unwindList, ok := unwindVal.([]any); ok {
				for _, item := range unwindList {
					newBinding := make(map[string]any, len(binding)+1)
					for k, v := range binding {
						newBinding[k] = v
					}
					newBinding[q.Unwind.Var] = item
					expanded = append(expanded, newBinding)
				}
			}
		}
		bindings = expanded
	}

	if len(q.Set) > 0 || len(q.Delete) > 0 || len(q.Remove) > 0 {
		return nil, fmt.Errorf("cypher: write clauses (SET/DELETE/REMOVE) require GraphWriter; use ExecuteQueryWrite")
	}

	return projectBindings(reader, q, bindings)
}

func evalUnwindExpr(e *Expression, binding map[string]any) any {
	if e.Kind == ExprVarRef {
		if val, ok := binding[e.Variable]; ok {
			return val
		}
		return nil
	}
	if e.Kind == ExprListLiteral {
		result := make([]any, len(e.ListItems))
		for i := range e.ListItems {
			result[i] = evalUnwindExpr(&e.ListItems[i], binding)
		}
		return result
	}
	val, _ := evalExpr(&EvalContext{Bindings: binding}, e)
	return val
}

func matchBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seedBindings []map[string]any) ([]map[string]any, error) {
	pat := q.Match.Pattern

	switch {
	case len(pat.Nodes) == 1 && len(pat.Rels) == 0:
		return matchNodeBindings(ctx, reader, q, seedBindings)
	case len(pat.Nodes) == 2 && len(pat.Rels) == 1:
		rel := pat.Rels[0]
		if rel.VarLength {
			return matchVarLengthBindings(ctx, reader, q, seedBindings)
		}
		return matchSingleHopBindings(ctx, reader, q, seedBindings)
	default:
		return matchMultiHopBindings(ctx, reader, q, seedBindings)
	}
}

func matchWhere(q *CypherQuery, bindings map[string]any) bool {
	if q.Where == nil {
		return true
	}
	if q.With != nil {
		return true
	}
	ok, err := evalBool(&EvalContext{Bindings: bindings}, q.Where)
	if err != nil {
		return false
	}
	return ok
}

func matchNodeBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seeds []map[string]any) ([]map[string]any, error) {
	np := q.Match.Pattern.Nodes[0]
	varName := np.Variable
	if varName == "" {
		varName = "_n"
	}

	hasExprProps := hasExpressionProps(np.Props)

	var candidates []*graphengine.Node
	var err error

	if hasExprProps {
		if len(np.Labels) > 0 {
			ids, err := reader.FindByLabel(np.Labels[0])
			if err != nil {
				return nil, err
			}
			for _, id := range ids {
				n, err := reader.GetNode(id)
				if err != nil {
					continue
				}
				candidates = append(candidates, n)
			}
		} else {
			err = reader.ForEachNode(func(n *graphengine.Node) error {
				candidates = append(candidates, n)
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		candidates, err = findCandidates(reader, np)
		if err != nil {
			return nil, err
		}
	}

	var bindings []map[string]any
	for _, n := range candidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		for _, seed := range seeds {
			if boundVal, ok := seed[varName]; ok {
				if bn, ok := boundVal.(*graphengine.Node); ok && bn.ID != n.ID {
					continue
				}
			}
			if hasExprProps {
				if !matchPropsWithBinding(n.Props, np.Props, seed) {
					continue
				}
			} else {
				if !matchProps(n.Props, np.Props) {
					continue
				}
			}

			binding := make(map[string]any, len(seed)+1)
			for k, v := range seed {
				binding[k] = v
			}
			binding[varName] = n

			if !matchWhere(q, binding) {
				continue
			}

			bindings = append(bindings, binding)
		}
	}

	return bindings, nil
}

func matchSingleHopBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seeds []map[string]any) ([]map[string]any, error) {
	pat := q.Match.Pattern
	aPat := pat.Nodes[0]
	rel := pat.Rels[0]
	bPat := pat.Nodes[1]

	aVar := aPat.Variable
	if aVar == "" {
		aVar = "_a"
	}
	rVar := rel.Variable
	bVar := bPat.Variable
	if bVar == "" {
		bVar = "_b"
	}

	aHasExprProps := hasExpressionProps(aPat.Props)
	hasExprProps := aHasExprProps || hasExpressionProps(bPat.Props) || hasExpressionProps(rel.Props)

	var bindings []map[string]any

	if boundVal, ok := seeds[0][aVar]; ok && len(seeds) == 1 {
		sourceNode, ok := boundVal.(*graphengine.Node)
		if !ok {
			return nil, nil
		}

		edges, err := reader.GetEdges(sourceNode.ID, rel.Dir, rel.Label)
		if err != nil {
			return nil, err
		}

		for _, e := range edges {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			targetID := e.To
			if rel.Dir == graphengine.Incoming {
				targetID = e.From
			} else if rel.Dir == graphengine.Both {
				if e.To == sourceNode.ID {
					targetID = e.From
				}
			}

			bNode, err := reader.GetNode(targetID)
			if err != nil {
				continue
			}

			if len(bPat.Labels) > 0 && !matchLabels(bNode.Labels, bPat.Labels) {
				continue
			}
			if hasExprProps {
				if !matchPropsWithBinding(bNode.Props, bPat.Props, seeds[0]) {
					continue
				}
			} else {
				if !matchProps(bNode.Props, bPat.Props) {
					continue
				}
			}

			if len(rel.Props) > 0 {
				if hasExprProps {
					if !matchPropsWithBinding(e.Props, rel.Props, seeds[0]) {
						continue
					}
				} else {
					if !matchEdgeProps(e, rel.Props) {
						continue
					}
				}
			}

			binding := make(map[string]any, len(seeds[0])+3)
			for k, v := range seeds[0] {
				binding[k] = v
			}
			binding[bVar] = bNode
			if rVar != "" {
				binding[rVar] = e
			}

			if !matchWhere(q, binding) {
				continue
			}

			bindings = append(bindings, binding)
		}

		return bindings, nil
	}

	scanPat := aPat
	if aHasExprProps {
		scanPat = NodePattern{Variable: aPat.Variable, Labels: aPat.Labels}
	}
	aCandidates, err := findCandidates(reader, scanPat)
	if err != nil {
		return nil, err
	}

	for _, a := range aCandidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		edges, err := reader.GetEdges(a.ID, rel.Dir, rel.Label)
		if err != nil {
			return nil, err
		}

		for _, e := range edges {
			targetID := e.To
			if rel.Dir == graphengine.Incoming {
				targetID = e.From
			} else if rel.Dir == graphengine.Both {
				if e.To == a.ID {
					targetID = e.From
				}
			}

			bNode, err := reader.GetNode(targetID)
			if err != nil {
				continue
			}

			if len(bPat.Labels) > 0 && !matchLabels(bNode.Labels, bPat.Labels) {
				continue
			}

			for _, seed := range seeds {
				if boundVal, ok := seed[aVar]; ok {
					if bn, ok := boundVal.(*graphengine.Node); ok && bn.ID != a.ID {
						continue
					}
				}
				if aHasExprProps && !matchPropsWithBinding(a.Props, aPat.Props, seed) {
					continue
				}
				if hasExprProps {
					if !matchPropsWithBinding(bNode.Props, bPat.Props, seed) {
						continue
					}
				} else {
					if !matchProps(bNode.Props, bPat.Props) {
						continue
					}
				}

				if len(rel.Props) > 0 {
					if hasExprProps {
						if !matchPropsWithBinding(e.Props, rel.Props, seed) {
							continue
						}
					} else {
						if !matchEdgeProps(e, rel.Props) {
							continue
						}
					}
				}

				binding := make(map[string]any, len(seed)+3)
				for k, v := range seed {
					binding[k] = v
				}
				binding[aVar] = a
				binding[bVar] = bNode
				if rVar != "" {
					binding[rVar] = e
				}

				if !matchWhere(q, binding) {
					continue
				}

				bindings = append(bindings, binding)
			}
		}
	}

	return bindings, nil
}

func executeMultiMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	bindings := []map[string]any{{}}

	for _, mc := range q.Matches {
		subQ := &CypherQuery{Match: mc}
		var err error
		bindings, err = matchBindings(ctx, reader, subQ, bindings)
		if err != nil {
			return nil, err
		}
		if len(bindings) == 0 {
			break
		}
	}

	if q.Where != nil {
		filtered := make([]map[string]any, 0, len(bindings))
		for _, b := range bindings {
			ok, err := evalBool(&EvalContext{Bindings: b}, q.Where)
			if err != nil {
				return nil, err
			}
			if ok {
				filtered = append(filtered, b)
			}
		}
		bindings = filtered
	}

	return projectBindings(reader, q, bindings)
}

func execMultiHopMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	pat := q.Match.Pattern
	firstPat := pat.Nodes[0]

	candidates, err := findCandidates(reader, firstPat)
	if err != nil {
		return nil, err
	}

	var allBindings []map[string]any
	for _, startNode := range candidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		chainResults := chainMultiHop(ctx, reader, pat, startNode, 0, nil)
		for _, cb := range chainResults {
			if q.Where != nil {
				ok, err := evalBool(&EvalContext{Bindings: cb}, q.Where)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}
			allBindings = append(allBindings, cb)
		}
	}

	return projectBindings(reader, q, allBindings)
}

func execCommaMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	pat := q.Match.Pattern
	candidateSets := make([][]*graphengine.Node, len(pat.Nodes))
	varNames := make([]string, len(pat.Nodes))
	for i, np := range pat.Nodes {
		candidates, err := findCandidates(reader, np)
		if err != nil {
			return nil, err
		}
		var matched []*graphengine.Node
		for _, n := range candidates {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if !matchProps(n.Props, np.Props) {
				continue
			}
			matched = append(matched, n)
		}
		candidateSets[i] = matched
		varNames[i] = np.Variable
		if varNames[i] == "" {
			varNames[i] = fmt.Sprintf("_n%d", i)
		}
	}

	var allBindings []map[string]any
	crossProductBindings(candidateSets, varNames, 0, map[string]any{}, func(binding map[string]any) {
		if matchWhere(q, binding) {
			cp := make(map[string]any, len(binding))
			for k, v := range binding {
				cp[k] = v
			}
			allBindings = append(allBindings, cp)
		}
	})

	return projectBindings(reader, q, allBindings)
}

func crossProductBindings(sets [][]*graphengine.Node, varNames []string, idx int, acc map[string]any, emit func(map[string]any)) {
	if idx >= len(sets) {
		emit(acc)
		return
	}
	for _, n := range sets[idx] {
		newAcc := make(map[string]any, len(acc)+1)
		for k, v := range acc {
			newAcc[k] = v
		}
		newAcc[varNames[idx]] = n
		crossProductBindings(sets, varNames, idx+1, newAcc, emit)
	}
}

func matchMultiHopBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seeds []map[string]any) ([]map[string]any, error) {
	pat := q.Match.Pattern
	firstPat := pat.Nodes[0]
	aVar := firstPat.Variable
	if aVar == "" {
		aVar = "_a"
	}

	var candidates []*graphengine.Node

	if boundVal, ok := seeds[0][aVar]; ok && len(seeds) == 1 {
		sourceNode, ok := boundVal.(*graphengine.Node)
		if !ok {
			return nil, nil
		}
		candidates = []*graphengine.Node{sourceNode}
	} else {
		var err error
		candidates, err = findCandidates(reader, firstPat)
		if err != nil {
			return nil, err
		}
	}

	var bindings []map[string]any
	for _, startNode := range candidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		chainResults := chainMultiHop(ctx, reader, pat, startNode, 0, nil)
		for _, cb := range chainResults {
			for _, seed := range seeds {
				if boundVal, ok := seed[aVar]; ok {
					if bn, ok := boundVal.(*graphengine.Node); ok && bn.ID != startNode.ID {
						continue
					}
				}
				merged := make(map[string]any, len(seed)+len(cb))
				for k, v := range seed {
					merged[k] = v
				}
				for k, v := range cb {
					merged[k] = v
				}
				if matchWhere(q, merged) {
					bindings = append(bindings, merged)
				}
			}
		}
	}
	return bindings, nil
}

func chainMultiHop(ctx context.Context, reader graphengine.GraphReader, pat Pattern, current *graphengine.Node, edgeIdx int, acc map[string]any) []map[string]any {
	if acc == nil {
		aVar := pat.Nodes[0].Variable
		if aVar == "" {
			aVar = "_a"
		}
		acc = map[string]any{aVar: current}
	}

	if edgeIdx >= len(pat.Rels) {
		return []map[string]any{acc}
	}

	rel := pat.Rels[edgeIdx]
	nextPat := pat.Nodes[edgeIdx+1]
	nextVar := nextPat.Variable
	if nextVar == "" {
		nextVar = fmt.Sprintf("_n%d", edgeIdx+1)
	}
	rVar := rel.Variable

	edges, err := reader.GetEdges(current.ID, rel.Dir, rel.Label)
	if err != nil {
		return nil
	}

	var results []map[string]any
	for _, e := range edges {
		if err := ctx.Err(); err != nil {
			return nil
		}

		targetID := e.To
		if rel.Dir == graphengine.Incoming {
			targetID = e.From
		} else if rel.Dir == graphengine.Both {
			if e.To == current.ID {
				targetID = e.From
			}
		}

		nextNode, err := reader.GetNode(targetID)
		if err != nil {
			continue
		}

		if len(nextPat.Labels) > 0 && !matchLabels(nextNode.Labels, nextPat.Labels) {
			continue
		}
		if !matchProps(nextNode.Props, nextPat.Props) {
			continue
		}
		if len(rel.Props) > 0 && !matchEdgeProps(e, rel.Props) {
			continue
		}

		newAcc := make(map[string]any, len(acc)+2)
		for k, v := range acc {
			newAcc[k] = v
		}
		newAcc[nextVar] = nextNode
		if rVar != "" {
			newAcc[rVar] = e
		}

		results = append(results, chainMultiHop(ctx, reader, pat, nextNode, edgeIdx+1, newAcc)...)
	}
	return results
}

func hasExpressionProps(props map[string]any) bool {
	if props == nil {
		return false
	}
	for _, v := range props {
		if _, ok := v.(Expression); ok {
			return true
		}
	}
	return false
}

func matchVarLengthBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seeds []map[string]any) ([]map[string]any, error) {
	pat := q.Match.Pattern
	aPat := pat.Nodes[0]
	rel := pat.Rels[0]
	bPat := pat.Nodes[1]

	aVar := aPat.Variable
	if aVar == "" {
		aVar = "_a"
	}
	bVar := bPat.Variable
	if bVar == "" {
		bVar = "_b"
	}
	rVar := rel.Variable

	var aCandidates []*graphengine.Node

	if boundVal, ok := seeds[0][aVar]; ok && len(seeds) == 1 {
		sourceNode, ok := boundVal.(*graphengine.Node)
		if !ok {
			return nil, nil
		}
		aCandidates = []*graphengine.Node{sourceNode}
	} else {
		var err error
		aCandidates, err = findCandidates(reader, aPat)
		if err != nil {
			return nil, err
		}
	}

	maxDepth := rel.MaxHops
	if maxDepth < 0 {
		maxDepth = defaultMaxDepth
	}
	minDepth := rel.MinHops

	var bindings []map[string]any

	for _, a := range aCandidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		var edgeFilter graphengine.EdgeFilter
		if rel.Label != "" {
			label := rel.Label
			edgeFilter = func(e *graphengine.Edge) bool {
				return e.Label == label
			}
		}

		var bfsErr error
		err := reader.BFS(a.ID, maxDepth, rel.Dir, edgeFilter, func(tr *graphengine.TraversalResult) bool {
			if tr.Depth < minDepth {
				return true
			}

			bNode := tr.Node

			if len(bPat.Labels) > 0 && !matchLabels(bNode.Labels, bPat.Labels) {
				return true
			}
			if !matchProps(bNode.Props, bPat.Props) {
				return true
			}

			for _, seed := range seeds {
				binding := make(map[string]any, len(seed)+3)
				for k, v := range seed {
					binding[k] = v
				}
				binding[aVar] = a
				binding[bVar] = bNode
				if rVar != "" {
					binding[rVar] = tr.Path
				}

				if q.Where != nil && q.With == nil {
					ok, err := evalBool(&EvalContext{Bindings: binding}, q.Where)
					if err != nil {
						bfsErr = err
						return false
					}
					if !ok {
						continue
					}
				}

				bindings = append(bindings, binding)
			}
			return true
		})
		if bfsErr != nil {
			return nil, bfsErr
		}
		if err != nil {
			return nil, err
		}
	}

	return bindings, nil
}
