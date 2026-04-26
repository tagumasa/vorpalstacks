// Cypher OPTIONAL MATCH execution, candidate finding, and result projection.
//
// Based on goraphdb/cypher_exec.go findCandidates, projectResults,
// projectPatternResults, and goraphdb/cypher_optional.go
// execWithOptionalMatch/collectMainBindings/attemptOptionalMatch/nullifyOptionalBindings/
// projectBindings. Extended with aggregation projection and DISTINCT support.

package cypherparser

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/core/storage/graphengine"
)

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
		startIdx := derefInt(q.Skip)
		if startIdx > len(sorted) {
			startIdx = len(sorted)
		}
		limitEnd := len(sorted)
		if q.Limit != nil && *q.Limit > 0 && startIdx+*q.Limit < limitEnd {
			limitEnd = startIdx + *q.Limit
		}
		for _, item := range sorted[startIdx:limitEnd] {
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
		startIdx := derefInt(q.Skip)
		if startIdx > len(sorted) {
			startIdx = len(sorted)
		}
		limitEnd := len(sorted)
		if q.Limit != nil && *q.Limit > 0 && startIdx+*q.Limit < limitEnd {
			limitEnd = startIdx + *q.Limit
		}
		for _, item := range sorted[startIdx:limitEnd] {
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
