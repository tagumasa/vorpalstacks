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

	"vorpalstacks/pkg/graphengine"
)

var (
	errLimitReached = errors.New("cypher: limit reached")
	errContextDone  = errors.New("cypher: context done")
)

const defaultMaxDepth = 50

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
