// Cypher match strategy executors: single-hop and variable-length path matching.
//
// Based on goraphdb/cypher_exec.go execSingleHopMatch, execVarLengthMatch.
// Strategy 2 handles two-node one-relationship patterns with index-aware
// candidate finding. Strategy 3 handles variable-length path expansion
// with configurable min/max hop bounds.

package cypherparser

import (
	"context"
	"errors"

	"vorpalstacks/pkg/graphengine"
)

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
