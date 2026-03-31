// Cypher write executor — handles CREATE, MERGE, SET, DELETE, REMOVE.
//
// Based on goraphdb/cypher_write.go executeCreate, executeCreatePattern,
// executeMerge, extended with:
//   - GraphReader/GraphWriter separation (goraphdb used *DB directly)
//   - Edge property support in CREATE (goraphdb did not support this)
//   - DETACH DELETE support
//   - REMOVE label/property support
//   - SET/DELETE/REMOVE on read-query results (MATCH ... SET/DELETE/REMOVE)
//
// Write execution model:
//   1. For CypherWrite (pure CREATE): create patterns, bind variables, project RETURN
//   2. For CypherMerge: match-or-create node, apply ON CREATE/MATCH SET, project RETURN
//   3. For CypherQuery with Set/Delete/Remove: run read phase first, apply mutations per row

package cypherparser

import (
	"context"
	"fmt"

	"vorpalstacks/pkg/graphengine"
)

// WriteStats tracks what was mutated by a write query.
type WriteStats struct {
	NodesCreated  int `json:"nodes_created"`
	NodesDeleted  int `json:"nodes_deleted"`
	EdgesCreated  int `json:"edges_created"`
	EdgesDeleted  int `json:"edges_deleted"`
	LabelsSet     int `json:"labels_set"`
	PropsSet      int `json:"properties_set"`
	LabelsRemoved int `json:"labels_removed"`
	PropsRemoved  int `json:"properties_removed"`
}

// ExecuteWrite runs a parsed CREATE query against the graph store.
func ExecuteWrite(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, w *CypherWrite, params map[string]any) (*CypherResult, error) {
	if err := resolveWriteParams(w, params); err != nil {
		return nil, err
	}
	return executeCreate(ctx, reader, writer, w)
}

// ExecuteMerge runs a parsed MERGE query against the graph store.
func ExecuteMerge(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, m *CypherMerge, params map[string]any) (*CypherResult, error) {
	if err := resolveMergeParams(m, params); err != nil {
		return nil, err
	}
	return executeMerge(ctx, reader, writer, m)
}

// ExecuteQueryWrite runs a CypherQuery that contains write clauses (SET, DELETE, REMOVE)
// after the read phase. Raw variable bindings are collected during the read phase,
// mutations are applied, then results are projected from the updated bindings.
func ExecuteQueryWrite(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, q *CypherQuery, params map[string]any) (*CypherResult, error) {
	if err := ResolveParams(q, params); err != nil {
		return nil, err
	}

	if len(q.Set) > 0 || len(q.Delete) > 0 || len(q.Remove) > 0 || q.Create != nil {
		bindings, err := collectReadBindings(ctx, reader, q)
		if err != nil {
			return nil, err
		}

		var stats WriteStats

		if q.Create != nil {
			if err := resolveWriteParams(q.Create, params); err != nil {
				return nil, err
			}
			for _, binding := range bindings {
				for _, cp := range q.Create.Creates {
					if err := createPattern(ctx, reader, writer, cp, binding, &stats); err != nil {
						return nil, err
					}
				}
			}
		}

		for _, binding := range bindings {
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

		return projectBindings(reader, q, bindings)
	}

	return executeCypher(ctx, reader, q)
}

// collectReadBindings runs the read phase and returns raw variable-binding maps
// (variable name → entity) instead of projected column-name maps.
func collectReadBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) ([]map[string]any, error) {
	if q.OptionalMatch != nil {
		return collectOptionalBindings(ctx, reader, q)
	}

	pat := q.Match.Pattern
	varName := ""
	if len(pat.Nodes) > 0 && pat.Nodes[0].Variable != "" {
		varName = pat.Nodes[0].Variable
	}

	switch {
	case len(pat.Nodes) == 1 && len(pat.Rels) == 0:
		return collectSingleNodeBindings(ctx, reader, q, pat.Nodes[0], varName)
	case len(pat.Nodes) == 2 && len(pat.Rels) == 1:
		return collectSingleHopBindings(ctx, reader, q, pat)
	case len(pat.Rels) == 0:
		return collectCommaBindings(ctx, reader, q, pat)
	default:
		return collectMultiHopBindings(ctx, reader, q, pat)
	}
}

func collectCommaBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, pat Pattern) ([]map[string]any, error) {
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

	var bindings []map[string]any
	crossProduct(candidateSets, varNames, 0, map[string]any{}, func(binding map[string]any) {
		if q.Where != nil {
			ok, err := evalBool(&EvalContext{Bindings: binding}, q.Where)
			if err != nil || !ok {
				return
			}
		}
		cp := make(map[string]any, len(binding))
		for k, v := range binding {
			cp[k] = v
		}
		bindings = append(bindings, cp)
	})
	return bindings, nil
}

func crossProduct(sets [][]*graphengine.Node, varNames []string, idx int, acc map[string]any, emit func(map[string]any)) {
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
		crossProduct(sets, varNames, idx+1, newAcc, emit)
	}
}

func collectSingleNodeBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, np NodePattern, varName string) ([]map[string]any, error) {
	if varName == "" {
		varName = "_n"
	}

	candidates, err := findCandidates(reader, np)
	if err != nil {
		return nil, err
	}

	var bindings []map[string]any
	for _, n := range candidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !matchProps(n.Props, np.Props) {
			continue
		}
		binding := map[string]any{varName: n}
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

func collectSingleHopBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, pat Pattern) ([]map[string]any, error) {
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

	var bindings []map[string]any
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
			if !matchProps(bNode.Props, bPat.Props) {
				continue
			}

			if len(rel.Props) > 0 && !matchEdgeProps(e, rel.Props) {
				continue
			}

			binding := map[string]any{aVar: a, bVar: bNode}
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

func collectOptionalBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery) ([]map[string]any, error) {
	mainBindings, err := collectMainBindings(ctx, reader, q)
	if err != nil {
		return nil, err
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

	return allBindings, nil
}

// ---------------------------------------------------------------------------
// CREATE executor
// ---------------------------------------------------------------------------

func executeCreate(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, w *CypherWrite) (*CypherResult, error) {
	bindings := make(map[string]any)
	var stats WriteStats

	for _, cp := range w.Creates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := createPattern(ctx, reader, writer, cp, bindings, &stats); err != nil {
			return nil, err
		}
	}

	result := &CypherResult{}
	if w.Return != nil {
		for _, item := range w.Return.Items {
			result.Columns = append(result.Columns, returnItemName(item))
		}
		row := make(map[string]any, len(w.Return.Items))
		for _, item := range w.Return.Items {
			colName := returnItemName(item)
			val, err := evalExpr(&EvalContext{Bindings: bindings}, &item.Expr)
			if err != nil {
				return nil, err
			}
			row[colName] = val
		}
		result.Rows = append(result.Rows, row)
	}

	if len(w.OrderBy) > 0 {
		sortRows(result.Rows, w.OrderBy)
	}

	if w.Skip != nil && *w.Skip > 0 {
		if *w.Skip < len(result.Rows) {
			result.Rows = result.Rows[*w.Skip:]
		} else {
			result.Rows = nil
		}
	}

	if w.Limit != nil && *w.Limit > 0 && len(result.Rows) > *w.Limit {
		result.Rows = result.Rows[:*w.Limit]
	}

	return result, nil
}

func createPattern(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, cp CreatePattern, bindings map[string]any, stats *WriteStats) error {
	if len(cp.Nodes) == 0 {
		return fmt.Errorf("cypher exec: CREATE pattern has no nodes")
	}

	nodeIDs := make([]graphengine.NodeID, len(cp.Nodes))
	for i, np := range cp.Nodes {
		if err := ctx.Err(); err != nil {
			return err
		}

		if np.Variable != "" {
			if existing, ok := bindings[np.Variable]; ok {
				if n, ok := existing.(*graphengine.Node); ok {
					nodeIDs[i] = n.ID
					continue
				}
			}
		}

		props := make(graphengine.Props, len(np.Props))
		for k, v := range np.Props {
			props[k] = v
		}

		id, err := writer.AddNode(np.Labels, props)
		if err != nil {
			return fmt.Errorf("cypher exec: CREATE node failed: %w", err)
		}
		nodeIDs[i] = id
		stats.NodesCreated++
		if len(np.Labels) > 0 {
			stats.LabelsSet += len(np.Labels)
		}
		stats.PropsSet += len(np.Props)

		if np.Variable != "" {
			n, err := reader.GetNode(id)
			if err != nil {
				return fmt.Errorf("cypher exec: CREATE node lookup failed: %w", err)
			}
			bindings[np.Variable] = n
		}
	}

	for i, rp := range cp.Rels {
		if err := ctx.Err(); err != nil {
			return err
		}

		fromID := nodeIDs[i]
		toID := nodeIDs[i+1]

		if rp.Dir == graphengine.Incoming {
			fromID, toID = toID, fromID
		}

		if rp.Label == "" {
			return fmt.Errorf("cypher exec: CREATE relationship requires a label (type)")
		}

		edgeProps := make(graphengine.Props, len(rp.Props))
		for k, v := range rp.Props {
			edgeProps[k] = v
		}

		edgeID, err := writer.AddEdge(fromID, toID, rp.Label, edgeProps)
		if err != nil {
			return fmt.Errorf("cypher exec: CREATE edge failed: %w", err)
		}
		stats.EdgesCreated++
		stats.PropsSet += len(rp.Props)

		if rp.Variable != "" {
			edges, err := reader.GetEdges(fromID, graphengine.Outgoing, rp.Label)
			if err == nil {
				for _, e := range edges {
					if e.ID == edgeID {
						bindings[rp.Variable] = e
						break
					}
				}
			}
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// MERGE executor
// ---------------------------------------------------------------------------

func executeMerge(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, m *CypherMerge) (*CypherResult, error) {
	mp := m.Pattern
	bindings := make(map[string]any)

	var matched *graphengine.Node

	if len(mp.Labels) > 0 && len(mp.Props) > 0 {
		for propKey, propVal := range mp.Props {
			if reader.HasUniqueConstraint(mp.Labels[0], propKey) {
				foundID, err := reader.FindByUniqueConstraint(mp.Labels[0], propKey, propVal)
				if err != nil {
					return nil, fmt.Errorf("cypher exec: MERGE lookup failed: %w", err)
				}
				if foundID != 0 {
					found, err := reader.GetNode(foundID)
					if err != nil {
						return nil, fmt.Errorf("cypher exec: MERGE get matched node failed: %w", err)
					}
					if matchLabels(found.Labels, mp.Labels) && matchProps(found.Props, mp.Props) {
						matched = found
					}
				}
				break
			}
		}
	}

	if matched == nil && len(mp.Labels) > 0 {
		candidates, err := reader.FindByLabel(mp.Labels[0])
		if err != nil {
			return nil, fmt.Errorf("cypher exec: MERGE label scan failed: %w", err)
		}
		for _, id := range candidates {
			n, err := reader.GetNode(id)
			if err != nil {
				continue
			}
			if matchLabels(n.Labels, mp.Labels) && matchProps(n.Props, mp.Props) {
				matched = n
				break
			}
		}
	}

	wasCreated := false
	if matched == nil {
		props := make(graphengine.Props, len(mp.Props))
		for k, v := range mp.Props {
			props[k] = v
		}

		id, err := writer.AddNode(mp.Labels, props)
		if err != nil {
			return nil, fmt.Errorf("cypher exec: MERGE create failed: %w", err)
		}

		matched, err = reader.GetNode(id)
		if err != nil {
			return nil, fmt.Errorf("cypher exec: MERGE get created node failed: %w", err)
		}
		wasCreated = true
	}

	if mp.Variable != "" {
		bindings[mp.Variable] = matched
	}

	var setItems []SetItem
	if wasCreated {
		setItems = m.OnCreateSet
	} else {
		setItems = m.OnMatchSet
	}
	if len(setItems) > 0 {
		stats := &WriteStats{}
		if err := applySet(reader, writer, setItems, bindings, stats); err != nil {
			return nil, err
		}
	}

	result := &CypherResult{}
	if m.Return != nil {
		for _, item := range m.Return.Items {
			result.Columns = append(result.Columns, returnItemName(item))
		}
		row := make(map[string]any, len(m.Return.Items))
		for _, item := range m.Return.Items {
			colName := returnItemName(item)
			val, err := evalExpr(&EvalContext{Bindings: bindings}, &item.Expr)
			if err != nil {
				return nil, err
			}
			row[colName] = val
		}
		result.Rows = append(result.Rows, row)
	}

	if len(m.OrderBy) > 0 {
		sortRows(result.Rows, m.OrderBy)
	}

	if m.Skip != nil && *m.Skip > 0 {
		if *m.Skip < len(result.Rows) {
			result.Rows = result.Rows[*m.Skip:]
		} else {
			result.Rows = nil
		}
	}

	if m.Limit != nil && *m.Limit > 0 && len(result.Rows) > *m.Limit {
		result.Rows = result.Rows[:*m.Limit]
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Post-read mutations (SET, DELETE, REMOVE on MATCH results)
// ---------------------------------------------------------------------------

func applySet(reader graphengine.GraphReader, writer graphengine.GraphWriter, items []SetItem, row map[string]any, stats *WriteStats) error {
	for _, si := range items {
		entity, ok := row[si.Variable]
		if !ok {
			continue
		}

		if si.SetLabel {
			switch e := entity.(type) {
			case *graphengine.Node:
				if err := writer.AddLabel(e.ID, si.Label); err != nil {
					return fmt.Errorf("cypher exec: SET label failed: %w", err)
				}
				refreshed, err := reader.GetNode(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: SET label refresh failed: %w", err)
				}
				row[si.Variable] = refreshed
				stats.LabelsSet++
			}
			continue
		}

		val, err := evalExpr(&EvalContext{Bindings: row}, &si.Value)
		if err != nil {
			return fmt.Errorf("cypher exec: SET %s eval failed: %w", si.Variable, err)
		}

		if si.Merge {
			switch e := entity.(type) {
			case *graphengine.Node:
				merged, ok := val.(map[string]any)
				if !ok {
					return fmt.Errorf("cypher exec: SET += requires a map, got %T", val)
				}
				updateProps := make(graphengine.Props, len(merged))
				for k, v := range merged {
					updateProps[k] = v
				}
				if err := writer.UpdateNode(e.ID, updateProps); err != nil {
					return fmt.Errorf("cypher exec: SET += node properties failed: %w", err)
				}
				refreshed, err := reader.GetNode(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: SET += node refresh failed: %w", err)
				}
				row[si.Variable] = refreshed
				stats.PropsSet += len(merged)
			case *graphengine.Edge:
				merged, ok := val.(map[string]any)
				if !ok {
					return fmt.Errorf("cypher exec: SET += requires a map, got %T", val)
				}
				updateProps := make(graphengine.Props, len(merged))
				for k, v := range merged {
					updateProps[k] = v
				}
				if err := writer.UpdateEdge(e.ID, updateProps); err != nil {
					return fmt.Errorf("cypher exec: SET += edge properties failed: %w", err)
				}
				refreshed, err := reader.GetEdge(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: SET += edge refresh failed: %w", err)
				}
				row[si.Variable] = refreshed
				stats.PropsSet += len(merged)
			}
			continue
		}

		if si.ReplaceAll {
			switch e := entity.(type) {
			case *graphengine.Node:
				newProps, ok := val.(map[string]any)
				if !ok {
					return fmt.Errorf("cypher exec: SET = requires a map, got %T", val)
				}
				if e.Props != nil {
					for k := range e.Props {
						if _, exists := newProps[k]; !exists {
							if err := writer.RemoveProperty(e.ID, k); err != nil {
								return fmt.Errorf("cypher exec: SET = remove property %s failed: %w", k, err)
							}
						}
					}
				}
				updateProps := make(graphengine.Props, len(newProps))
				for k, v := range newProps {
					updateProps[k] = v
				}
				if err := writer.UpdateNode(e.ID, updateProps); err != nil {
					return fmt.Errorf("cypher exec: SET = node properties failed: %w", err)
				}
				refreshed, err := reader.GetNode(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: SET = node refresh failed: %w", err)
				}
				row[si.Variable] = refreshed
				stats.PropsSet += len(newProps)
			case *graphengine.Edge:
				newProps, ok := val.(map[string]any)
				if !ok {
					return fmt.Errorf("cypher exec: SET = requires a map, got %T", val)
				}
				if e.Props != nil {
					for k := range e.Props {
						if _, exists := newProps[k]; !exists {
							if err := writer.RemoveEdgeProperty(e.ID, k); err != nil {
								return fmt.Errorf("cypher exec: SET = remove edge property %s failed: %w", k, err)
							}
						}
					}
				}
				updateProps := make(graphengine.Props, len(newProps))
				for k, v := range newProps {
					updateProps[k] = v
				}
				if err := writer.UpdateEdge(e.ID, updateProps); err != nil {
					return fmt.Errorf("cypher exec: SET = edge properties failed: %w", err)
				}
				refreshed, err := reader.GetEdge(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: SET = edge refresh failed: %w", err)
				}
				row[si.Variable] = refreshed
				stats.PropsSet += len(newProps)
			}
			continue
		}

		switch e := entity.(type) {
		case *graphengine.Node:
			if err := writer.UpdateNode(e.ID, graphengine.Props{si.Property: val}); err != nil {
				return fmt.Errorf("cypher exec: SET node property failed: %w", err)
			}
			refreshed, err := reader.GetNode(e.ID)
			if err != nil {
				return fmt.Errorf("cypher exec: SET node refresh failed: %w", err)
			}
			row[si.Variable] = refreshed
			stats.PropsSet++
		case *graphengine.Edge:
			if err := writer.UpdateEdge(e.ID, graphengine.Props{si.Property: val}); err != nil {
				return fmt.Errorf("cypher exec: SET edge property failed: %w", err)
			}
			refreshed, err := reader.GetEdge(e.ID)
			if err != nil {
				return fmt.Errorf("cypher exec: SET edge refresh failed: %w", err)
			}
			row[si.Variable] = refreshed
			stats.PropsSet++
		}
	}
	return nil
}

func applyDelete(reader graphengine.GraphReader, writer graphengine.GraphWriter, vars []string, row map[string]any, stats *WriteStats, detach bool) error {
	for _, v := range vars {
		entity, ok := row[v]
		if !ok {
			continue
		}

		switch e := entity.(type) {
		case *graphengine.Node:
			if err := writer.DeleteNode(e.ID); err != nil {
				return fmt.Errorf("cypher exec: DELETE node %d failed: %w", e.ID, err)
			}
			stats.NodesDeleted++
		case *graphengine.Edge:
			if err := writer.DeleteEdge(e.ID); err != nil {
				return fmt.Errorf("cypher exec: DELETE edge %d failed: %w", e.ID, err)
			}
			stats.EdgesDeleted++
		}
	}
	return nil
}

func applyRemove(reader graphengine.GraphReader, writer graphengine.GraphWriter, items []RemoveItem, row map[string]any, stats *WriteStats) error {
	for _, ri := range items {
		entity, ok := row[ri.Variable]
		if !ok {
			continue
		}

		switch e := entity.(type) {
		case *graphengine.Node:
			switch ri.Kind {
			case RemoveLabel:
				if err := writer.RemoveLabel(e.ID, ri.Label); err != nil {
					return fmt.Errorf("cypher exec: REMOVE label %s from node %d failed: %w", ri.Label, e.ID, err)
				}
				refreshed, err := reader.GetNode(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: REMOVE label refresh failed: %w", err)
				}
				row[ri.Variable] = refreshed
				stats.LabelsRemoved++
			case RemoveProperty:
				if err := writer.RemoveProperty(e.ID, ri.Property); err != nil {
					return fmt.Errorf("cypher exec: REMOVE property %s from node %d failed: %w", ri.Property, e.ID, err)
				}
				refreshed, err := reader.GetNode(e.ID)
				if err != nil {
					return fmt.Errorf("cypher exec: REMOVE property refresh failed: %w", err)
				}
				row[ri.Variable] = refreshed
				stats.PropsRemoved++
			}
		case *graphengine.Edge:
			if ri.Kind == RemoveProperty {
				if err := writer.RemoveEdgeProperty(e.ID, ri.Property); err != nil {
					return fmt.Errorf("cypher exec: REMOVE property %s from edge %d failed: %w", ri.Property, e.ID, err)
				}
				stats.PropsRemoved++
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Parameter resolution for write queries
// ---------------------------------------------------------------------------

func resolveWriteParams(w *CypherWrite, params map[string]any) error {
	for _, cp := range w.Creates {
		for i := range cp.Nodes {
			if err := resolveNodeProps(&cp.Nodes[i], params); err != nil {
				return err
			}
		}
		for i := range cp.Rels {
			if err := resolveRelPropsWrite(&cp.Rels[i], params); err != nil {
				return err
			}
		}
	}
	if w.Return != nil {
		for i := range w.Return.Items {
			if err := resolveExprParams(&w.Return.Items[i].Expr, params); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveMergeParams(m *CypherMerge, params map[string]any) error {
	if err := resolveNodePropsDirect(m.Pattern.Props, params); err != nil {
		return err
	}
	for i := range m.OnCreateSet {
		if err := resolveExprParams(&m.OnCreateSet[i].Value, params); err != nil {
			return err
		}
	}
	for i := range m.OnMatchSet {
		if err := resolveExprParams(&m.OnMatchSet[i].Value, params); err != nil {
			return err
		}
	}
	if m.Return != nil {
		for i := range m.Return.Items {
			if err := resolveExprParams(&m.Return.Items[i].Expr, params); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveNodeProps(np *NodePattern, params map[string]any) error {
	if np.Props == nil {
		return nil
	}
	for k, v := range np.Props {
		if p, ok := v.(paramRef); ok {
			val, ok := params[string(p)]
			if !ok {
				return fmt.Errorf("cypher exec: parameter $%s not provided", p)
			}
			np.Props[k] = val
		} else if expr, ok := v.(Expression); ok {
			if err := resolveExprParams(&expr, params); err != nil {
				return err
			}
			val, err := evalExpr(&EvalContext{Params: params}, &expr)
			if err != nil {
				return fmt.Errorf("cypher exec: failed to evaluate property expression: %w", err)
			}
			np.Props[k] = val
		}
	}
	return nil
}

func resolveRelPropsWrite(rp *RelPattern, params map[string]any) error {
	if rp.Props == nil {
		return nil
	}
	for k, v := range rp.Props {
		if p, ok := v.(paramRef); ok {
			val, ok := params[string(p)]
			if !ok {
				return fmt.Errorf("cypher exec: parameter $%s not provided", p)
			}
			rp.Props[k] = val
		} else if expr, ok := v.(Expression); ok {
			if err := resolveExprParams(&expr, params); err != nil {
				return err
			}
			val, err := evalExpr(&EvalContext{Params: params}, &expr)
			if err != nil {
				return fmt.Errorf("cypher exec: failed to evaluate edge property expression: %w", err)
			}
			rp.Props[k] = val
		}
	}
	return nil
}

func resolveNodePropsDirect(props map[string]any, params map[string]any) error {
	if props == nil {
		return nil
	}
	for k, v := range props {
		if p, ok := v.(paramRef); ok {
			val, ok := params[string(p)]
			if !ok {
				return fmt.Errorf("cypher exec: parameter $%s not provided", p)
			}
			props[k] = val
		} else if expr, ok := v.(Expression); ok {
			if err := resolveExprParams(&expr, params); err != nil {
				return err
			}
			val, err := evalExpr(&EvalContext{Params: params}, &expr)
			if err != nil {
				return fmt.Errorf("cypher exec: failed to evaluate property expression: %w", err)
			}
			props[k] = val
		}
	}
	return nil
}

func collectMultiHopBindings(ctx context.Context, reader graphengine.GraphReader, q *CypherQuery, pat Pattern) ([]map[string]any, error) {
	firstPat := pat.Nodes[0]
	aVar := firstPat.Variable
	if aVar == "" {
		aVar = "_a"
	}

	candidates, err := findCandidates(reader, firstPat)
	if err != nil {
		return nil, err
	}

	var bindings []map[string]any
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
			bindings = append(bindings, cb)
		}
	}
	return bindings, nil
}
