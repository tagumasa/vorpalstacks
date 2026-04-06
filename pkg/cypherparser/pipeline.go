// Cypher pipeline executor — handles WITH and UNWIND clauses.
//
// Original implementation — goraphdb does not have WITH/UNWIND/pipeline support.
//
// Execution model:
//   - WITH: evaluates expressions from current bindings, projects them as new
//     variable bindings, optionally applies WHERE filtering, DISTINCT, ORDER BY,
//     SKIP, and LIMIT. The projected bindings then feed into subsequent clauses.
//   - UNWIND: expands a list expression into individual rows, binding each element
//     to the specified variable.
//   - Pipeline: chains multiple segments, each producing bindings for the next.
//     Segment ordering: WITH → UNWIND → MATCH → OPTIONAL MATCH → SET/DELETE/REMOVE → RETURN.
//
// WITH clause acts like an intermediate RETURN — it projects a subset of variables
// and optionally filters/reorders. After WITH, only the projected variables are
// visible to subsequent clauses.
//
// UNWIND clause expands a list (literal, parameter, or expression) into rows.
// Each row has one binding for the UNWIND variable. Combined with MATCH, this
// enables patterns like:
//
//	UNWIND $names AS name MATCH (n:Person {name: name}) RETURN n

package cypherparser

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"vorpalstacks/pkg/graphengine"
)

// ExecutePipeline runs a multi-segment CypherPipeline.
// Each segment's WITH clause produces bindings that feed into the next segment.
func ExecutePipeline(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, pipe *CypherPipeline, params map[string]any) (*CypherResult, error) {
	if len(pipe.Segments) == 0 {
		return &CypherResult{Columns: []string{}, Rows: []map[string]any{}}, nil
	}

	var prevBindings []map[string]any

	for i, seg := range pipe.Segments {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if seg.With != nil && prevBindings != nil {
			var err error
			prevBindings, err = applyWith(seg.With, seg.With.Where, prevBindings)
			if err != nil {
				return nil, err
			}
		}

		q := seg.Query
		if err := ResolveParams(q, params); err != nil {
			return nil, err
		}

		hasWrite := len(q.Set) > 0 || len(q.Delete) > 0 || len(q.Remove) > 0

		if i < len(pipe.Segments)-1 {
			result, err := executeSegmentRead(ctx, reader, q, prevBindings)
			if err != nil {
				return nil, err
			}
			prevBindings = result
		} else {
			if hasWrite {
				result, err := executeSegmentWrite(ctx, reader, writer, q, prevBindings)
				if err != nil {
					return nil, err
				}
				return result, nil
			}
			result, err := executeSegmentFinal(ctx, reader, q, prevBindings)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}

	return &CypherResult{Columns: []string{}, Rows: []map[string]any{}}, nil
}

func executeSegmentRead(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, prevBindings []map[string]any) ([]map[string]any, error) {
	inputBindings := prevBindings
	if inputBindings == nil {
		inputBindings = []map[string]any{{}}
	}

	var result []map[string]any

	for _, binding := range inputBindings {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if q.Unwind != nil {
			unwound, err := applyUnwind(binding, q.Unwind)
			if err != nil {
				return nil, err
			}
			for _, ub := range unwound {
				segmentResult, err := runSegmentClauses(ctx, reader, q, ub)
				if err != nil {
					return nil, err
				}
				result = append(result, segmentResult...)
			}
		} else {
			segmentResult, err := runSegmentClauses(ctx, reader, q, binding)
			if err != nil {
				return nil, err
			}
			result = append(result, segmentResult...)
		}
	}

	if q.With != nil {
		return applyWith(q.With, q.With.Where, result)
	}

	return result, nil
}

func executeSegmentFinal(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, prevBindings []map[string]any) (*CypherResult, error) {
	inputBindings := prevBindings
	if inputBindings == nil {
		inputBindings = []map[string]any{{}}
	}

	var allBindings []map[string]any

	for _, binding := range inputBindings {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if q.Unwind != nil {
			unwound, err := applyUnwind(binding, q.Unwind)
			if err != nil {
				return nil, err
			}
			for _, ub := range unwound {
				segmentResult, err := runSegmentClauses(ctx, reader, q, ub)
				if err != nil {
					return nil, err
				}
				allBindings = append(allBindings, segmentResult...)
			}
		} else {
			segmentResult, err := runSegmentClauses(ctx, reader, q, binding)
			if err != nil {
				return nil, err
			}
			allBindings = append(allBindings, segmentResult...)
		}
	}

	if q.With != nil {
		var err error
		allBindings, err = applyWith(q.With, q.With.Where, allBindings)
		if err != nil {
			return nil, err
		}
	}

	return projectBindings(reader, q, allBindings)
}

func executeSegmentWrite(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, q *CypherQuery, prevBindings []map[string]any) (*CypherResult, error) {
	inputBindings := prevBindings
	if inputBindings == nil {
		inputBindings = []map[string]any{{}}
	}

	var allBindings []map[string]any

	for _, binding := range inputBindings {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if q.Unwind != nil {
			unwound, err := applyUnwind(binding, q.Unwind)
			if err != nil {
				return nil, err
			}
			for _, ub := range unwound {
				segmentResult, err := runSegmentClauses(ctx, reader, q, ub)
				if err != nil {
					return nil, err
				}
				allBindings = append(allBindings, segmentResult...)
			}
		} else {
			segmentResult, err := runSegmentClauses(ctx, reader, q, binding)
			if err != nil {
				return nil, err
			}
			allBindings = append(allBindings, segmentResult...)
		}
	}

	if q.With != nil {
		var err error
		allBindings, err = applyWith(q.With, q.With.Where, allBindings)
		if err != nil {
			return nil, err
		}
	}

	var stats WriteStats
	for _, binding := range allBindings {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if len(q.Set) > 0 {
			if err := applySet(reader, writer, q.Set, binding, &stats); err != nil {
				return nil, err
			}
		}
		if len(q.Delete) > 0 {
			if err := applyDelete(reader, writer, q.Delete, binding, &stats, q.DeleteDetach); err != nil {
				return nil, err
			}
		}
		if len(q.Remove) > 0 {
			if err := applyRemove(reader, writer, q.Remove, binding, &stats); err != nil {
				return nil, err
			}
		}
	}

	return projectBindings(reader, q, allBindings)
}

func runSegmentClauses(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seedBinding map[string]any) ([]map[string]any, error) {
	pat := q.Match.Pattern

	if len(pat.Nodes) == 0 {
		return []map[string]any{seedBinding}, nil
	}

	switch {
	case len(pat.Nodes) == 1 && len(pat.Rels) == 0:
		return runSegmentNodeMatch(ctx, reader, q, seedBinding)
	case len(pat.Nodes) == 2 && len(pat.Rels) == 1:
		rel := pat.Rels[0]
		if rel.VarLength {
			return runSegmentVarLengthMatch(ctx, reader, q, seedBinding)
		}
		return runSegmentSingleHopMatch(ctx, reader, q, seedBinding)
	default:
		return nil, fmt.Errorf("cypher: unsupported pattern with %d nodes and %d relationships",
			len(pat.Nodes), len(pat.Rels))
	}
}

func runSegmentNodeMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seedBinding map[string]any) ([]map[string]any, error) {
	nodePat := q.Match.Pattern.Nodes[0]
	varName := nodePat.Variable
	if varName == "" {
		varName = "_n"
	}

	candidates, err := findCandidates(reader, nodePat)
	if err != nil {
		return nil, err
	}

	var bindings []map[string]any
	for _, n := range candidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !matchProps(n.Props, nodePat.Props) {
			continue
		}

		binding := make(map[string]any, len(seedBinding)+1)
		for k, v := range seedBinding {
			binding[k] = v
		}
		binding[varName] = n

		if q.Where != nil {
			ok, err := evalBool(&EvalContext{Bindings: binding}, q.Where)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}

		bindings = append(bindings, binding)
	}

	return bindings, nil
}

func runSegmentSingleHopMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seedBinding map[string]any) ([]map[string]any, error) {
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

	aHasExprProps := hasExpressionProps(aPat.Props)
	bHasExprProps := hasExpressionProps(bPat.Props)
	rHasExprProps := hasExpressionProps(rel.Props)
	_ = aHasExprProps

	if boundVal, ok := seedBinding[aVar]; ok {
		sourceNode, ok := boundVal.(*graphengine.Node)
		if !ok {
			return nil, nil
		}

		edges, err := reader.GetEdges(sourceNode.ID, rel.Dir, rel.Label)
		if err != nil {
			return nil, err
		}

		var bindings []map[string]any
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
			if bHasExprProps {
				if !matchPropsWithBinding(bNode.Props, bPat.Props, seedBinding) {
					continue
				}
			} else {
				if !matchProps(bNode.Props, bPat.Props) {
					continue
				}
			}

			if rHasExprProps {
				if !matchPropsWithBinding(e.Props, rel.Props, seedBinding) {
					continue
				}
			} else if len(rel.Props) > 0 && !matchEdgeProps(e, rel.Props) {
				continue
			}

			binding := make(map[string]any, len(seedBinding)+3)
			for k, v := range seedBinding {
				binding[k] = v
			}
			binding[bVar] = bNode
			if rVar != "" {
				binding[rVar] = e
			}

			if q.Where != nil {
				ok, err := evalBool(&EvalContext{Bindings: binding}, q.Where)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
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

	var bindings []map[string]any
	for _, a := range aCandidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if aHasExprProps && !matchPropsWithBinding(a.Props, aPat.Props, seedBinding) {
			continue
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
			if bHasExprProps {
				if !matchPropsWithBinding(bNode.Props, bPat.Props, seedBinding) {
					continue
				}
			} else {
				if !matchProps(bNode.Props, bPat.Props) {
					continue
				}
			}

			if rHasExprProps {
				if !matchPropsWithBinding(e.Props, rel.Props, seedBinding) {
					continue
				}
			} else if len(rel.Props) > 0 && !matchEdgeProps(e, rel.Props) {
				continue
			}

			binding := make(map[string]any, len(seedBinding)+3)
			for k, v := range seedBinding {
				binding[k] = v
			}
			binding[aVar] = a
			binding[bVar] = bNode
			if rVar != "" {
				binding[rVar] = e
			}

			if q.Where != nil {
				ok, err := evalBool(&EvalContext{Bindings: binding}, q.Where)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}

			bindings = append(bindings, binding)
		}
	}

	return bindings, nil
}

func runSegmentVarLengthMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, seedBinding map[string]any) ([]map[string]any, error) {
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

	if boundVal, ok := seedBinding[aVar]; ok {
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

			binding := make(map[string]any, len(seedBinding)+3)
			for k, v := range seedBinding {
				binding[k] = v
			}
			binding[aVar] = a
			binding[bVar] = bNode
			if rVar != "" {
				binding[rVar] = tr.Path
			}

			if q.Where != nil {
				ok, err := evalBool(&EvalContext{Bindings: binding}, q.Where)
				if err != nil {
					bfsErr = err
					return false
				}
				if !ok {
					return true
				}
			}

			bindings = append(bindings, binding)
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

// ---------------------------------------------------------------------------
// WITH projection
// ---------------------------------------------------------------------------

// applyWith projects bindings through a WITH clause.
// Each input binding is evaluated against the WITH items to produce a new binding
// with only the projected variables.
// If any WITH item contains an aggregation, all non-aggregated items take their
// value from the first row (standard Cypher GROUP BY semantics).
func applyWith(wc *WithClause, where *Expression, bindings []map[string]any) ([]map[string]any, error) {
	if wc == nil {
		return bindings, nil
	}

	if len(bindings) == 0 {
		return []map[string]any{}, nil
	}

	hasAgg := false
	for _, item := range wc.Items {
		if HasAggregation(&item.Expr) {
			hasAgg = true
			break
		}
	}

	if hasAgg {
		return applyWithAggregation(wc, where, bindings)
	}

	var result []map[string]any

	for _, binding := range bindings {
		newBinding := make(map[string]any, len(wc.Items))
		for _, item := range wc.Items {
			val, err := evalExpr(&EvalContext{Bindings: binding}, &item.Expr)
			if err != nil {
				return nil, err
			}

			name := item.Alias
			if name == "" {
				if item.Expr.Kind == ExprVarRef {
					name = item.Expr.Variable
				} else {
					name = exprName(item.Expr)
				}
			}

			newBinding[name] = val
		}

		if where != nil {
			ok, err := evalBool(&EvalContext{Bindings: newBinding}, where)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}

		result = append(result, newBinding)
	}

	if wc.Distinct {
		result = distinctBindings(result)
	}

	if len(wc.OrderBy) > 0 {
		sortRows(result, wc.OrderBy)
	}

	if wc.Skip != nil && *wc.Skip > 0 {
		if *wc.Skip < len(result) {
			result = result[*wc.Skip:]
		} else {
			result = nil
		}
	}

	if wc.Limit != nil && *wc.Limit > 0 && len(result) > *wc.Limit {
		result = result[:*wc.Limit]
	}

	return result, nil
}

func applyWithAggregation(wc *WithClause, where *Expression, bindings []map[string]any) ([]map[string]any, error) {
	aggResult := make(map[string]any)

	for _, item := range wc.Items {
		name := item.Alias
		if name == "" {
			if item.Expr.Kind == ExprVarRef {
				name = item.Expr.Variable
			} else {
				name = exprName(item.Expr)
			}
		}

		if item.Expr.Kind == ExprAggregation {
			val, err := ComputeAggregation(item.Expr.AggFunc, item.Expr.AggArg, item.Expr.AggDistinct, bindings)
			if err != nil {
				return nil, err
			}
			aggResult[name] = val
		} else if HasAggregation(&item.Expr) {
			ctx := &EvalContext{Bindings: bindings[0], AggValues: make(map[*Expression]any)}
			for _, aggItem := range wc.Items {
				if aggItem.Expr.Kind == ExprAggregation {
					val, err := ComputeAggregation(aggItem.Expr.AggFunc, aggItem.Expr.AggArg, aggItem.Expr.AggDistinct, bindings)
					if err != nil {
						return nil, err
					}
					ctx.AggValues[&aggItem.Expr] = val
				}
			}
			val, err := evalExpr(ctx, &item.Expr)
			if err != nil {
				return nil, err
			}
			aggResult[name] = val
		} else {
			val, err := evalExpr(&EvalContext{Bindings: bindings[0]}, &item.Expr)
			if err != nil {
				return nil, err
			}
			aggResult[name] = val
		}
	}

	if wc.Distinct {
		return []map[string]any{aggResult}, nil
	}

	if where != nil {
		ok, err := evalBool(&EvalContext{Bindings: aggResult}, where)
		if err != nil {
			return nil, err
		}
		if !ok {
			return []map[string]any{}, nil
		}
	}

	result := []map[string]any{aggResult}

	if len(wc.OrderBy) > 0 {
		sortRows(result, wc.OrderBy)
	}

	if wc.Skip != nil && *wc.Skip > 0 {
		if *wc.Skip < len(result) {
			result = result[*wc.Skip:]
		} else {
			result = []map[string]any{}
		}
	}

	if wc.Limit != nil && *wc.Limit > 0 && len(result) > *wc.Limit {
		result = result[:*wc.Limit]
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// UNWIND expansion
// ---------------------------------------------------------------------------

// applyUnwind expands a list expression from the current binding into
// individual bindings, each with the UNWIND variable set to one element.
func applyUnwind(binding map[string]any, uc *UnwindClause) ([]map[string]any, error) {
	val, err := evalExpr(&EvalContext{Bindings: binding}, &uc.Expr)
	if err != nil {
		return nil, fmt.Errorf("cypher: UNWIND expression eval failed: %w", err)
	}

	list, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("cypher: UNWIND requires a list, got %T", val)
	}

	result := make([]map[string]any, len(list))
	for i, item := range list {
		newBinding := make(map[string]any, len(binding)+1)
		for k, v := range binding {
			newBinding[k] = v
		}
		newBinding[uc.Var] = item
		result[i] = newBinding
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func distinctBindings(bindings []map[string]any) []map[string]any {
	seen := make(map[string]bool, len(bindings))
	var result []map[string]any
	for _, b := range bindings {
		key := bindingKey(b)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, b)
	}
	return result
}

func bindingKey(binding map[string]any) string {
	keys := make([]string, 0, len(binding))
	for k := range binding {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		v := binding[k]
		fmt.Fprintf(&b, "%s=%T:%v|", k, v, v)
	}
	return b.String()
}
