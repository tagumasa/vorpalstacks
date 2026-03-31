// Cypher read executor — runs a parsed CypherQuery against a graphengine.GraphReader.
//
// Based on goraphdb/cypher_exec.go executeCypher/executeCypherNormal strategy dispatch,
// execNodeMatch, execSingleHopMatch, execVarLengthMatch, findCandidates, projectResults,
// projectPatternResults, extractWhereEquality, and goraphdb/cypher_optional.go
// execWithOptionalMatch/collectMainBindings/attemptOptionalMatch/nullifyOptionalBindings/
// projectBindings.
//
// Architecture:
//   - Execute delegates to the appropriate strategy based on pattern shape
//   - Strategy 1: simple node match (1 node, 0 rels)
//   - Strategy 2: single-hop pattern (2 nodes, 1 rel)
//   - Strategy 3: variable-length path (2 nodes, 1 var-length rel)
//   - OPTIONAL MATCH: left-outer-join semantics via collectMainBindings + attemptOptionalMatch
//   - All strategies use index-aware candidate finding (label index → property index → full scan)
//   - LIMIT push-down: early exit in all scan loops when no ORDER BY
//   - Aggregation: two-phase (collect rows → compute aggregation → project)
//   - Result projection: ORDER BY, LIMIT, DISTINCT, SKIP

package cypherparser

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"vorpalstacks/pkg/graphengine"
)

var (
	errLimitReached = errors.New("cypher: limit reached")
	errContextDone  = errors.New("cypher: context done")
)

const defaultMaxDepth = 50

// Execute runs a parsed CypherQuery against a graphengine.GraphReader.
// It handles EXPLAIN, PROFILE, WITH/UNWIND pipeline projection, and
// dispatches to the appropriate strategy.
func Execute(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, params map[string]any) (*CypherResult, error) {
	if err := ResolveParams(q, params); err != nil {
		return nil, err
	}

	return executeCypher(ctx, reader, q)
}

func executeCypher(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	if q.Explain == ExplainOnly {
		plan := BuildExplainPlan(q, reader)
		return &CypherResult{
			Columns: []string{"plan"},
			Rows:    []map[string]any{{"plan": plan.String()}},
		}, nil
	}

	if q.Explain == ExplainProfile {
		normalQ := *q
		normalQ.Explain = ExplainNone
		result, err := executeCypherNormal(ctx, reader, &normalQ)
		if err != nil {
			return nil, err
		}
		plan := BuildProfilePlan(q, reader, result)
		return &CypherResult{
			Columns: []string{"plan", "results"},
			Rows:    []map[string]any{{"plan": plan.String(), "results": result}},
		}, nil
	}

	return executeCypherNormal(ctx, reader, q)
}

func executeCypherNormal(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	if q.OptionalMatch != nil {
		return execWithOptionalMatch(ctx, reader, q)
	}

	hasWrite := len(q.Set) > 0 || len(q.Delete) > 0 || len(q.Remove) > 0

	if q.Unwind != nil || q.With != nil || hasWrite {
		return executeFullQuery(ctx, reader, q)
	}

	if len(q.Matches) > 1 {
		return executeMultiMatch(ctx, reader, q)
	}

	pat := q.Match.Pattern

	switch {
	case len(pat.Nodes) == 1 && len(pat.Rels) == 0:
		return execNodeMatch(ctx, reader, q)

	case len(pat.Nodes) == 2 && len(pat.Rels) == 1:
		rel := pat.Rels[0]
		if rel.VarLength {
			return execVarLengthMatch(ctx, reader, q)
		}
		return execSingleHopMatch(ctx, reader, q)

	case len(pat.Rels) == 0:
		return execCommaMatch(ctx, reader, q)

	default:
		return execMultiHopMatch(ctx, reader, q)
	}
}

// ---------------------------------------------------------------------------
// Strategy 1: Simple node match — MATCH (n) / MATCH (n {name: "Alice"})
// ---------------------------------------------------------------------------

func execNodeMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	nodePat := q.Match.Pattern.Nodes[0]
	varName := nodePat.Variable
	if varName == "" {
		varName = "_n"
	}

	hasOrderBy := len(q.OrderBy) > 0
	limitVal := derefInt(q.Limit)
	hasAgg := hasAggregationInReturn(q)

	if len(nodePat.Labels) > 0 {
		candidates, err := reader.FindByLabel(nodePat.Labels[0])
		if err != nil {
			return nil, err
		}
		var nodes []*graphengine.Node
		for _, id := range candidates {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			n, err := reader.GetNode(id)
			if err != nil {
				continue
			}
			if !matchLabels(n.Labels, nodePat.Labels) {
				continue
			}
			if !matchProps(n.Props, nodePat.Props) {
				continue
			}
			if q.Where != nil {
				ectx := &EvalContext{Bindings: map[string]any{varName: n}}
				ok, err := evalBool(ectx, q.Where)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}
			nodes = append(nodes, n)
			if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(nodes) >= limitVal+derefInt(q.Skip) {
				break
			}
		}
		return projectResults(reader, q, nodes, nil, varName, "")
	}

	if len(nodePat.Props) > 0 {
		for key, val := range nodePat.Props {
			candidates, err := reader.FindByProperty(key, val)
			if err != nil {
				continue
			}
			var nodes []*graphengine.Node
			for _, id := range candidates {
				n, err := reader.GetNode(id)
				if err != nil {
					continue
				}
				if !matchProps(n.Props, nodePat.Props) {
					continue
				}
				if q.Where != nil {
					ectx := &EvalContext{Bindings: map[string]any{varName: n}}
					ok, err := evalBool(ectx, q.Where)
					if err != nil {
						return nil, err
					}
					if !ok {
						continue
					}
				}
				nodes = append(nodes, n)
			}
			return projectResults(reader, q, nodes, nil, varName, "")
		}
	}

	if q.Where != nil && len(nodePat.Props) == 0 {
		if prop, val, ok := extractWhereEquality(q.Where, varName); ok {
			candidates, err := reader.FindByProperty(prop, val)
			if err == nil {
				var nodes []*graphengine.Node
				for _, id := range candidates {
					n, err := reader.GetNode(id)
					if err != nil {
						continue
					}
					if q.Where != nil {
						ectx := &EvalContext{Bindings: map[string]any{varName: n}}
						ok, err := evalBool(ectx, q.Where)
						if err != nil {
							return nil, err
						}
						if !ok {
							continue
						}
					}
					nodes = append(nodes, n)
				}
				return projectResults(reader, q, nodes, nil, varName, "")
			}
		}
	}

	var nodes []*graphengine.Node
	err := reader.ForEachNode(func(n *graphengine.Node) error {
		if err := ctx.Err(); err != nil {
			return errContextDone
		}
		if !matchProps(n.Props, nodePat.Props) {
			return nil
		}
		if q.Where != nil {
			ectx := &EvalContext{Bindings: map[string]any{varName: n}}
			ok, err := evalBool(ectx, q.Where)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
		}
		nodes = append(nodes, n)
		if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(nodes) >= limitVal+derefInt(q.Skip) {
			return errLimitReached
		}
		return nil
	})
	if errors.Is(err, errContextDone) {
		return nil, ctx.Err()
	}
	if err != nil && !errors.Is(err, errLimitReached) {
		return nil, err
	}

	return projectResults(reader, q, nodes, nil, varName, "")
}

// ---------------------------------------------------------------------------
// Strategy 2: Single-hop pattern — MATCH (a)-[r:LABEL]->(b)
// ---------------------------------------------------------------------------

func execSingleHopMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
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

	if len(aPat.Props) == 0 && len(aPat.Labels) == 0 && rel.Label != "" && rel.Dir == graphengine.Outgoing {
		return execSingleHopByEdgeType(ctx, reader, q, aPat, rel, bPat, aVar, rVar, bVar)
	}

	limitVal := derefInt(q.Limit)
	hasOrderBy := len(q.OrderBy) > 0
	hasAgg := hasAggregationInReturn(q)

	aCandidates, err := findCandidates(reader, aPat)
	if err != nil {
		return nil, err
	}

	var rows []resultRow

	for _, a := range aCandidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(rows) >= limitVal+derefInt(q.Skip) {
			break
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
			if !matchProps(bNode.Props, bPat.Props) {
				continue
			}

			bindings := map[string]any{aVar: a, bVar: bNode}
			if rVar != "" {
				bindings[rVar] = e
			}

			if q.Where != nil {
				ectx := &EvalContext{Bindings: bindings}
				ok, err := evalBool(ectx, q.Where)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}

			if len(rel.Props) > 0 {
				if !matchEdgeProps(e, rel.Props) {
					continue
				}
			}

			rows = append(rows, resultRow{a: a, r: e, b: bNode})

			if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(rows) >= limitVal+derefInt(q.Skip) {
				break
			}
		}
	}

	return projectPatternResults(reader, q, rows, aVar, rVar, bVar)
}

func execSingleHopByEdgeType(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, aPat NodePattern, rel RelPattern, bPat NodePattern, aVar, rVar, bVar string) (*CypherResult, error) {
	var rows []resultRow
	limitVal := derefInt(q.Limit)
	hasOrderBy := len(q.OrderBy) > 0
	hasAgg := hasAggregationInReturn(q)

	err := reader.ScanEdgesByType(rel.Label, func(edge *graphengine.Edge, src, dst *graphengine.Node) error {
		if err := ctx.Err(); err != nil {
			return errContextDone
		}

		if rel.Dir == graphengine.Incoming {
			src, dst = dst, src
		}

		if len(aPat.Labels) > 0 && !matchLabels(src.Labels, aPat.Labels) {
			return nil
		}
		if !matchProps(src.Props, aPat.Props) {
			return nil
		}
		if len(bPat.Labels) > 0 && !matchLabels(dst.Labels, bPat.Labels) {
			return nil
		}
		if !matchProps(dst.Props, bPat.Props) {
			return nil
		}

		bindings := map[string]any{aVar: src, bVar: dst}
		if rVar != "" {
			bindings[rVar] = edge
		}

		if q.Where != nil {
			ectx := &EvalContext{Bindings: bindings}
			ok, err := evalBool(ectx, q.Where)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
		}

		if len(rel.Props) > 0 {
			if !matchEdgeProps(edge, rel.Props) {
				return nil
			}
		}

		rows = append(rows, resultRow{a: src, r: edge, b: dst})

		if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(rows) >= limitVal+derefInt(q.Skip) {
			return errLimitReached
		}
		return nil
	})
	if errors.Is(err, errContextDone) {
		return nil, ctx.Err()
	}
	if err != nil && !errors.Is(err, errLimitReached) {
		return nil, err
	}

	return projectPatternResults(reader, q, rows, aVar, rVar, bVar)
}

// ---------------------------------------------------------------------------
// Strategy 3: Variable-length path — MATCH (a)-[:LABEL*min..max]->(b)
// ---------------------------------------------------------------------------

func execVarLengthMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
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

	limitVal := derefInt(q.Limit)
	hasOrderBy := len(q.OrderBy) > 0
	hasAgg := hasAggregationInReturn(q)

	aCandidates, err := findCandidates(reader, aPat)
	if err != nil {
		return nil, err
	}

	maxDepth := rel.MaxHops
	if maxDepth < 0 {
		maxDepth = defaultMaxDepth
	}
	minDepth := rel.MinHops

	var rows []resultRow

	rVar := rel.Variable

	for _, a := range aCandidates {
		if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(rows) >= limitVal+derefInt(q.Skip) {
			break
		}
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

			bindings := map[string]any{aVar: a, bVar: bNode}
			if rVar != "" {
				bindings[rVar] = tr.Path
			}

			if q.Where != nil {
				ectx := &EvalContext{Bindings: bindings}
				ok, err := evalBool(ectx, q.Where)
				if err != nil {
					bfsErr = err
					return false
				}
				if !ok {
					return true
				}
			}

			rows = append(rows, resultRow{a: a, b: bNode, rPath: tr.Path})

			if limitVal > 0 && !hasOrderBy && !hasAgg && !q.Return.Distinct && len(rows) >= limitVal+derefInt(q.Skip) {
				return false
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

	return projectPatternResults(reader, q, rows, aVar, rVar, bVar)
}

// ---------------------------------------------------------------------------
// OPTIONAL MATCH — left-outer-join semantics.
// ---------------------------------------------------------------------------

func execWithOptionalMatch(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) (*CypherResult, error) {
	mainQ := &CypherQuery{
		Matches: q.Matches,
		Match:   q.Match,
		Where:   q.Where,
	}

	mainBindings, err := collectMainBindings(ctx, reader, mainQ)
	if err != nil {
		return nil, err
	}

	if len(q.Matches) > 1 {
		for _, mc := range q.Matches[1:] {
			subQ := &CypherQuery{Match: mc}
			var chainErr error
			mainBindings, chainErr = matchBindings(ctx, reader, subQ, mainBindings)
			if chainErr != nil {
				return nil, chainErr
			}
			if len(mainBindings) == 0 {
				break
			}
		}
	}

	optPat := q.OptionalMatch.Pattern

	var allBindings []map[string]any

	for _, binding := range mainBindings {
		optRows, err := attemptOptionalMatch(ctx, reader, binding, optPat, q.OptionalWhere)
		if err != nil {
			return nil, err
		}

		if len(optRows) == 0 {
			nulled := nullifyOptionalBindings(binding, optPat)
			allBindings = append(allBindings, nulled)
		} else {
			allBindings = append(allBindings, optRows...)
		}
	}

	return projectBindings(reader, q, allBindings)
}

func collectMainBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) ([]map[string]any, error) {
	pat := q.Match.Pattern
	var bindings []map[string]any

	switch {
	case len(pat.Nodes) == 1 && len(pat.Rels) == 0:
		np := pat.Nodes[0]
		varName := np.Variable
		if varName == "" {
			varName = "_n"
		}

		candidates, err := findCandidates(reader, np)
		if err != nil {
			return nil, err
		}

		for _, n := range candidates {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if !matchProps(n.Props, np.Props) {
				continue
			}
			if q.Where != nil {
				ectx := &EvalContext{Bindings: map[string]any{varName: n}}
				ok, err := evalBool(ectx, q.Where)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}
			bindings = append(bindings, map[string]any{varName: n})
		}

	case len(pat.Nodes) == 2 && len(pat.Rels) == 1:
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

		aCandidates, err := findCandidates(reader, aPat)
		if err != nil {
			return nil, err
		}

		for _, a := range aCandidates {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			edges, err := reader.GetEdges(a.ID, rel.Dir, rel.Label)
			if err != nil {
				continue
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
				if !matchProps(bNode.Props, bPat.Props) {
					continue
				}

				row := map[string]any{aVar: a, bVar: bNode}
				if rVar != "" {
					row[rVar] = e
				}

				if q.Where != nil {
					ectx := &EvalContext{Bindings: row}
					ok, err := evalBool(ectx, q.Where)
					if err != nil {
						return nil, err
					}
					if !ok {
						continue
					}
				}

				if len(rel.Props) > 0 && !matchEdgeProps(e, rel.Props) {
					continue
				}

				bindings = append(bindings, row)
			}
		}
	}

	return bindings, nil
}

func attemptOptionalMatch(ctx context.Context, reader graphengine.GraphReader, binding map[string]any, optPat Pattern, optWhere *Expression) ([]map[string]any, error) {
	if len(optPat.Nodes) == 2 && len(optPat.Rels) == 1 {
		aPat := optPat.Nodes[0]
		rel := optPat.Rels[0]
		bPat := optPat.Nodes[1]

		aVar := aPat.Variable
		if aVar == "" {
			aVar = "_a"
		}
		rVar := rel.Variable
		bVar := bPat.Variable
		if bVar == "" {
			bVar = "_b"
		}

		boundVal, ok := binding[aVar]
		if !ok {
			return nil, nil
		}
		sourceNode, ok := boundVal.(*graphengine.Node)
		if !ok {
			return nil, nil
		}

		edges, err := reader.GetEdges(sourceNode.ID, rel.Dir, rel.Label)
		if err != nil {
			return nil, nil
		}

		var rows []map[string]any
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
			if !matchProps(bNode.Props, bPat.Props) {
				continue
			}

			row := make(map[string]any, len(binding)+2)
			for k, v := range binding {
				row[k] = v
			}
			row[bVar] = bNode
			if rVar != "" {
				row[rVar] = e
			}

			if optWhere != nil {
				ectx := &EvalContext{Bindings: row}
				ok, err := evalBool(ectx, optWhere)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}

			rows = append(rows, row)
		}
		return rows, nil
	}

	if len(optPat.Nodes) == 1 && len(optPat.Rels) == 0 {
		np := optPat.Nodes[0]
		varName := np.Variable
		if varName == "" {
			varName = "_opt"
		}

		candidates, err := findCandidates(reader, np)
		if err != nil {
			return nil, nil
		}

		var rows []map[string]any
		for _, n := range candidates {
			row := make(map[string]any, len(binding)+1)
			for k, v := range binding {
				row[k] = v
			}
			row[varName] = n

			if optWhere != nil {
				ectx := &EvalContext{Bindings: row}
				ok, err := evalBool(ectx, optWhere)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}

			rows = append(rows, row)
		}
		return rows, nil
	}

	return nil, nil
}

func nullifyOptionalBindings(binding map[string]any, optPat Pattern) map[string]any {
	row := make(map[string]any, len(binding)+3)
	for k, v := range binding {
		row[k] = v
	}

	for _, np := range optPat.Nodes {
		if np.Variable != "" {
			if _, exists := row[np.Variable]; !exists {
				row[np.Variable] = nil
			}
		}
	}
	for _, rp := range optPat.Rels {
		if rp.Variable != "" {
			if _, exists := row[rp.Variable]; !exists {
				row[rp.Variable] = nil
			}
		}
	}

	return row
}

// ---------------------------------------------------------------------------
// Candidate finding — index-aware
// ---------------------------------------------------------------------------

func findCandidates(reader graphengine.GraphReader, np NodePattern) ([]*graphengine.Node, error) {
	if len(np.Labels) > 0 {
		candidates, err := reader.FindByLabel(np.Labels[0])
		if err != nil {
			return nil, err
		}
		if len(np.Labels) > 1 || len(np.Props) > 0 {
			var filtered []*graphengine.Node
			for _, id := range candidates {
				n, err := reader.GetNode(id)
				if err != nil {
					continue
				}
				if matchLabels(n.Labels, np.Labels) && matchProps(n.Props, np.Props) {
					filtered = append(filtered, n)
				}
			}
			return filtered, nil
		}
		var nodes []*graphengine.Node
		for _, id := range candidates {
			n, err := reader.GetNode(id)
			if err != nil {
				continue
			}
			nodes = append(nodes, n)
		}
		return nodes, nil
	}

	if len(np.Props) == 0 {
		var nodes []*graphengine.Node
		err := reader.ForEachNode(func(n *graphengine.Node) error {
			nodes = append(nodes, n)
			return nil
		})
		return nodes, err
	}

	for key, val := range np.Props {
		candidates, err := reader.FindByProperty(key, val)
		if err != nil {
			continue
		}
		var filtered []*graphengine.Node
		for _, id := range candidates {
			n, err := reader.GetNode(id)
			if err != nil {
				continue
			}
			if matchProps(n.Props, np.Props) {
				filtered = append(filtered, n)
			}
		}
		return filtered, nil
	}

	var nodes []*graphengine.Node
	err := reader.ForEachNode(func(n *graphengine.Node) error {
		if matchProps(n.Props, np.Props) {
			nodes = append(nodes, n)
		}
		return nil
	})
	return nodes, err
}

// ---------------------------------------------------------------------------
// Result projection
// ---------------------------------------------------------------------------

type resultRow struct {
	a     *graphengine.Node
	r     *graphengine.Edge
	b     *graphengine.Node
	rPath []*graphengine.Edge
}

func bindRelationshipVar(bindings map[string]any, rVar string, rr resultRow) {
	if rr.r != nil {
		bindings[rVar] = rr.r
	} else if rr.rPath != nil {
		bindings[rVar] = rr.rPath
	}
}

func projectResults(reader graphengine.GraphReader, q *CypherQuery, nodes []*graphengine.Node, edges []*graphengine.Edge, nodeVar, edgeVar string) (*CypherResult, error) {
	if q.Limit != nil && *q.Limit == 0 {
		return &CypherResult{Columns: buildColumns(q.Return.Items)}, nil
	}

	if hasAggregationInReturn(q) {
		return projectAggregatedResults(reader, q, nodes, nodeVar)
	}

	result := &CypherResult{}

	for _, item := range q.Return.Items {
		result.Columns = append(result.Columns, returnItemName(item))
	}

	if len(q.OrderBy) > 0 && q.Limit != nil && *q.Limit > 0 && !q.Return.Distinct {
		heapSize := *q.Limit + derefInt(q.Skip)
		if heapSize <= 0 {
			heapSize = *q.Limit
		}
		h := newTopKHeap(q.OrderBy, heapSize)

		for _, n := range nodes {
			bindings := map[string]any{nodeVar: n}
			sortKey := evalSortKey(q.OrderBy, bindings)
			h.offer(topKItem{sortKey: sortKey, source: n})
		}

		sorted := h.sorted()
		for _, item := range sorted {
			n := item.source.(*graphengine.Node)
			bindings := map[string]any{nodeVar: n}
			row := make(map[string]any)
			for _, ri := range q.Return.Items {
				colName := returnItemName(ri)
				val, err := evalExpr(&EvalContext{Bindings: bindings}, &ri.Expr)
				if err != nil {
					return nil, err
				}
				row[colName] = val
			}
			result.Rows = append(result.Rows, row)
		}
		return result, nil
	}

	for _, n := range nodes {
		bindings := map[string]any{nodeVar: n}
		row := make(map[string]any)
		for _, item := range q.Return.Items {
			colName := returnItemName(item)
			val, err := evalExpr(&EvalContext{Bindings: bindings}, &item.Expr)
			if err != nil {
				return nil, err
			}
			row[colName] = val
		}
		result.Rows = append(result.Rows, row)
	}

	if q.Return.Distinct {
		result.Rows = distinctRows(result.Rows, result.Columns)
	}

	if len(q.OrderBy) > 0 {
		sortRows(result.Rows, q.OrderBy)
	}

	if q.Skip != nil && *q.Skip > 0 {
		if *q.Skip < len(result.Rows) {
			result.Rows = result.Rows[*q.Skip:]
		} else {
			result.Rows = []map[string]any{}
		}
	}

	if q.Limit != nil && *q.Limit > 0 && len(result.Rows) > *q.Limit {
		result.Rows = result.Rows[:*q.Limit]
	}

	return result, nil
}

func projectPatternResults(reader graphengine.GraphReader, q *CypherQuery, rows []resultRow, aVar, rVar, bVar string) (*CypherResult, error) {
	if q.Limit != nil && *q.Limit == 0 {
		return &CypherResult{Columns: buildColumns(q.Return.Items)}, nil
	}

	if hasAggregationInReturn(q) {
		return projectAggregatedPatternResults(reader, q, rows, aVar, rVar, bVar)
	}

	result := &CypherResult{}

	for _, item := range q.Return.Items {
		result.Columns = append(result.Columns, returnItemName(item))
	}

	if len(q.OrderBy) > 0 && q.Limit != nil && *q.Limit > 0 && !q.Return.Distinct {
		heapSize := *q.Limit + derefInt(q.Skip)
		if heapSize <= 0 {
			heapSize = *q.Limit
		}
		h := newTopKHeap(q.OrderBy, heapSize)

		for _, rr := range rows {
			bindings := map[string]any{aVar: rr.a, bVar: rr.b}
			if rVar != "" {
				bindRelationshipVar(bindings, rVar, rr)
			}
			sortKey := evalSortKey(q.OrderBy, bindings)
			h.offer(topKItem{sortKey: sortKey, source: rr})
		}

		sorted := h.sorted()
		for _, item := range sorted {
			rr := item.source.(resultRow)
			bindings := map[string]any{aVar: rr.a, bVar: rr.b}
			if rVar != "" {
				bindRelationshipVar(bindings, rVar, rr)
			}
			row := make(map[string]any)
			for _, ri := range q.Return.Items {
				colName := returnItemName(ri)
				val, err := evalExpr(&EvalContext{Bindings: bindings}, &ri.Expr)
				if err != nil {
					return nil, err
				}
				row[colName] = val
			}
			result.Rows = append(result.Rows, row)
		}
		return result, nil
	}

	for _, rr := range rows {
		bindings := map[string]any{aVar: rr.a, bVar: rr.b}
		if rVar != "" {
			bindRelationshipVar(bindings, rVar, rr)
		}

		row := make(map[string]any)
		for _, item := range q.Return.Items {
			colName := returnItemName(item)
			val, err := evalExpr(&EvalContext{Bindings: bindings}, &item.Expr)
			if err != nil {
				return nil, err
			}
			row[colName] = val
		}
		result.Rows = append(result.Rows, row)
	}

	if q.Return.Distinct {
		result.Rows = distinctRows(result.Rows, result.Columns)
	}

	if len(q.OrderBy) > 0 {
		sortRows(result.Rows, q.OrderBy)
	}

	if q.Skip != nil && *q.Skip > 0 {
		if *q.Skip < len(result.Rows) {
			result.Rows = result.Rows[*q.Skip:]
		} else {
			result.Rows = []map[string]any{}
		}
	}

	if q.Limit != nil && *q.Limit > 0 && len(result.Rows) > *q.Limit {
		result.Rows = result.Rows[:*q.Limit]
	}

	return result, nil
}

func projectBindings(reader graphengine.GraphReader, q *CypherQuery, bindings []map[string]any) (*CypherResult, error) {
	if q.Limit != nil && *q.Limit == 0 {
		return &CypherResult{Columns: buildColumns(q.Return.Items)}, nil
	}

	if hasAggregationInReturn(q) {
		return projectAggregatedBindings(reader, q, bindings)
	}

	result := &CypherResult{}

	for _, item := range q.Return.Items {
		result.Columns = append(result.Columns, returnItemName(item))
	}

	for _, binding := range bindings {
		row := make(map[string]any, len(q.Return.Items))
		for _, item := range q.Return.Items {
			colName := returnItemName(item)
			val, err := evalExpr(&EvalContext{Bindings: binding}, &item.Expr)
			if err != nil {
				return nil, err
			}
			row[colName] = val
		}
		result.Rows = append(result.Rows, row)
	}

	if q.Return.Distinct {
		result.Rows = distinctRows(result.Rows, result.Columns)
	}

	if len(q.OrderBy) > 0 {
		sortRows(result.Rows, q.OrderBy)
	}

	if q.Skip != nil && *q.Skip > 0 {
		if *q.Skip < len(result.Rows) {
			result.Rows = result.Rows[*q.Skip:]
		} else {
			result.Rows = []map[string]any{}
		}
	}

	if q.Limit != nil && *q.Limit > 0 && len(result.Rows) > *q.Limit {
		result.Rows = result.Rows[:*q.Limit]
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Aggregation projection
// ---------------------------------------------------------------------------

func projectAggregatedResults(reader graphengine.GraphReader, q *CypherQuery, nodes []*graphengine.Node, nodeVar string) (*CypherResult, error) {
	result := &CypherResult{}
	for _, item := range q.Return.Items {
		result.Columns = append(result.Columns, returnItemName(item))
	}

	bindingRows := make([]map[string]any, len(nodes))
	for i, n := range nodes {
		bindingRows[i] = map[string]any{nodeVar: n}
	}

	aggResult := computeAggregations(q.Return.Items, bindingRows)
	nonAggRow := make(map[string]any)
	for _, item := range q.Return.Items {
		if item.Expr.Kind == ExprAggregation {
			continue
		}
		colName := returnItemName(item)
		var val any
		if len(bindingRows) > 0 {
			bindings := bindingRows[0]
			var err error
			val, err = evalExpr(&EvalContext{Bindings: bindings}, &item.Expr)
			if err != nil {
				return nil, err
			}
		}
		nonAggRow[colName] = val
	}
	for colName, val := range aggResult {
		nonAggRow[colName] = val
	}

	result.Rows = append(result.Rows, nonAggRow)
	return result, nil
}

func projectAggregatedPatternResults(reader graphengine.GraphReader, q *CypherQuery, rows []resultRow, aVar, rVar, bVar string) (*CypherResult, error) {
	result := &CypherResult{}
	for _, item := range q.Return.Items {
		result.Columns = append(result.Columns, returnItemName(item))
	}

	bindingRows := make([]map[string]any, len(rows))
	for i, rr := range rows {
		bindingRows[i] = map[string]any{aVar: rr.a, bVar: rr.b}
		if rVar != "" {
			bindRelationshipVar(bindingRows[i], rVar, rr)
		}
	}

	aggResult := computeAggregations(q.Return.Items, bindingRows)
	nonAggRow := make(map[string]any)
	for _, item := range q.Return.Items {
		if item.Expr.Kind == ExprAggregation {
			continue
		}
		colName := returnItemName(item)
		var val any
		if len(bindingRows) > 0 {
			bindings := bindingRows[0]
			var err error
			val, err = evalExpr(&EvalContext{Bindings: bindings}, &item.Expr)
			if err != nil {
				return nil, err
			}
		}
		nonAggRow[colName] = val
	}
	for colName, val := range aggResult {
		nonAggRow[colName] = val
	}

	result.Rows = append(result.Rows, nonAggRow)
	return result, nil
}

func projectAggregatedBindings(reader graphengine.GraphReader, q *CypherQuery, bindings []map[string]any) (*CypherResult, error) {
	result := &CypherResult{}
	for _, item := range q.Return.Items {
		result.Columns = append(result.Columns, returnItemName(item))
	}

	aggResult := computeAggregations(q.Return.Items, bindings)
	nonAggRow := make(map[string]any)
	for _, item := range q.Return.Items {
		if item.Expr.Kind == ExprAggregation {
			continue
		}
		colName := returnItemName(item)
		var val any
		if len(bindings) > 0 {
			binding := bindings[0]
			var err error
			val, err = evalExpr(&EvalContext{Bindings: binding}, &item.Expr)
			if err != nil {
				return nil, err
			}
		}
		nonAggRow[colName] = val
	}
	for colName, val := range aggResult {
		nonAggRow[colName] = val
	}

	result.Rows = append(result.Rows, nonAggRow)
	return result, nil
}

func computeAggregations(items []ReturnItem, rows []map[string]any) map[string]any {
	result := make(map[string]any)
	for _, item := range items {
		if item.Expr.Kind != ExprAggregation {
			continue
		}
		colName := returnItemName(item)
		val, _ := ComputeAggregation(item.Expr.AggFunc, item.Expr.AggArg, item.Expr.AggDistinct, rows)
		result[colName] = val
	}
	return result
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func extractWhereEquality(where *Expression, nodeVar string) (string, any, bool) {
	if where == nil {
		return "", nil, false
	}

	if where.Kind == ExprComparison && where.Op == OpEq {
		if where.Left != nil && where.Left.Kind == ExprPropAccess && where.Left.Object == nodeVar &&
			where.Right != nil && where.Right.Kind == ExprLiteral {
			return where.Left.Property, where.Right.LitValue, true
		}
		if where.Right != nil && where.Right.Kind == ExprPropAccess && where.Right.Object == nodeVar &&
			where.Left != nil && where.Left.Kind == ExprLiteral {
			return where.Right.Property, where.Left.LitValue, true
		}
	}

	if where.Kind == ExprAnd {
		for i := range where.Operands {
			if prop, val, ok := extractWhereEquality(&where.Operands[i], nodeVar); ok {
				return prop, val, true
			}
		}
	}

	return "", nil, false
}

func matchEdgeProps(edge *graphengine.Edge, constraints map[string]any) bool {
	if constraints == nil {
		return true
	}
	if edge.Props == nil {
		return false
	}
	for key, expected := range constraints {
		actual, ok := edge.Props[key]
		if !ok {
			return false
		}
		if compareValues(actual, expected) != 0 {
			return false
		}
	}
	return true
}

func hasAggregationInReturn(q *CypherQuery) bool {
	for _, item := range q.Return.Items {
		if HasAggregation(&item.Expr) {
			return true
		}
	}
	return false
}

func distinctRows(rows []map[string]any, columns []string) []map[string]any {
	seen := make(map[string]bool, len(rows))
	var result []map[string]any
	for _, row := range rows {
		key := rowKey(row, columns)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, row)
	}
	return result
}

func rowKey(row map[string]any, columns []string) string {
	var b strings.Builder
	for _, col := range columns {
		val := row[col]
		switch v := val.(type) {
		case *graphengine.Node:
			fmt.Fprintf(&b, "node:%d|", v.ID)
		case *graphengine.Edge:
			fmt.Fprintf(&b, "edge:%d|", v.ID)
		default:
			fmt.Fprintf(&b, "%T:%v|", val, val)
		}
	}
	return b.String()
}

func buildColumns(items []ReturnItem) []string {
	cols := make([]string, len(items))
	for i, item := range items {
		cols[i] = returnItemName(item)
	}
	return cols
}

// ---------------------------------------------------------------------------
// Full query execution (WITH, UNWIND, write clauses)
// ---------------------------------------------------------------------------

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
		bindings, err = applyWith(q.With, nil, bindings)
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
	aVar := firstPat.Variable
	if aVar == "" {
		aVar = "_a"
	}

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
