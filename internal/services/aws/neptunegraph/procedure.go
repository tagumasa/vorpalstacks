package neptunegraph

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/common/request"
	ngstore "vorpalstacks/internal/store/aws/neptunegraph"
	"vorpalstacks/pkg/cypherparser"
	"vorpalstacks/pkg/graphengine"
)

// procedureDispatcher handles execution of vector procedure CALL statements
// within ExecuteQuery. It bridges the parsed Cypher CALL AST to the graphengine
// embedding store operations.
type procedureDispatcher struct {
	db     *graphengine.DB
	dim    int32
	hasVec bool
}

func newProcedureDispatcher(db *graphengine.DB, dim int32) *procedureDispatcher {
	return &procedureDispatcher{db: db, dim: dim, hasVec: dim > 0}
}

// ExecuteCall dispatches a parsed Cypher CALL to the appropriate vector procedure.
func (pd *procedureDispatcher) ExecuteCall(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if !pd.hasVec {
		return nil, fmt.Errorf("ValidationException: vector operations are not available on this graph (no VectorSearchConfiguration)")
	}

	switch call.Name {
	case "neptune.algo.vectors.get":
		return pd.execGet(call, bindings)
	case "neptune.algo.vectors.upsert":
		return pd.execUpsert(call, bindings)
	case "neptune.algo.vectors.remove":
		return pd.execRemove(call, bindings)
	case "neptune.algo.vectors.distance.byNode":
		return pd.execDistanceByNode(call, bindings)
	case "neptune.algo.vectors.distance.byEmbedding":
		return pd.execDistanceByEmbedding(call, bindings)
	case "neptune.algo.vectors.topK.byNode":
		return pd.execTopKByNode(call, bindings)
	case "neptune.algo.vectors.topK.byEmbedding":
		return pd.execTopKByEmbedding(call, bindings)
	default:
		return nil, fmt.Errorf("ValidationException: unknown procedure %s", call.Name)
	}
}

func (pd *procedureDispatcher) execGet(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 1 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.get requires at least 1 argument")
	}

	nodeIDs, err := pd.resolveNodeIDs(call.Args[0], bindings)
	if err != nil {
		return nil, err
	}

	rows := make([]map[string]interface{}, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		emb, err := pd.db.Embeddings.Get(id)
		if err != nil {
			return nil, err
		}
		node, err := pd.db.GetNode(id)
		if err != nil {
			continue
		}
		row := map[string]interface{}{"node": node}
		if emb != nil {
			row["embedding"] = emb
		}
		rows = append(rows, row)
	}

	return &cypherparser.CypherResult{
		Columns: []string{"node", "embedding"},
		Rows:    rows,
	}, nil
}

func (pd *procedureDispatcher) execUpsert(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 2 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.upsert requires 2 arguments (node, embedding)")
	}

	nodeIDs, err := pd.resolveNodeIDs(call.Args[0], bindings)
	if err != nil {
		return nil, err
	}

	emb, err := pd.resolveFloatList(call.Args[1], bindings)
	if err != nil {
		return nil, fmt.Errorf("ValidationException: embedding must be a list of floats: %v", err)
	}

	rows := make([]map[string]interface{}, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		if err := pd.db.Embeddings.Upsert(id, emb, pd.dim); err != nil {
			return nil, fmt.Errorf("ValidationException: %v", err)
		}
		node, err := pd.db.GetNode(id)
		if err != nil {
			node = &graphengine.Node{ID: id}
		}
		rows = append(rows, map[string]interface{}{
			"node":      node,
			"embedding": emb,
			"success":   true,
		})
	}

	return &cypherparser.CypherResult{
		Columns: []string{"node", "embedding", "success"},
		Rows:    rows,
	}, nil
}

func (pd *procedureDispatcher) execRemove(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 1 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.remove requires at least 1 argument")
	}

	nodeIDs, err := pd.resolveNodeIDs(call.Args[0], bindings)
	if err != nil {
		return nil, err
	}

	rows := make([]map[string]interface{}, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		existed, err := pd.db.Embeddings.Remove(id)
		if err != nil {
			return nil, err
		}
		node, err := pd.db.GetNode(id)
		if err != nil {
			node = &graphengine.Node{ID: id}
		}
		rows = append(rows, map[string]interface{}{
			"node":    node,
			"success": existed,
		})
	}

	return &cypherparser.CypherResult{
		Columns: []string{"node", "success"},
		Rows:    rows,
	}, nil
}

func (pd *procedureDispatcher) execDistanceByNode(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 2 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.distance.byNode requires at least 2 arguments")
	}

	sourceIDs, err := pd.resolveNodeIDs(call.Args[0], bindings)
	if err != nil {
		return nil, err
	}
	targetIDs, err := pd.resolveNodeIDs(call.Args[1], bindings)
	if err != nil {
		return nil, err
	}

	metric := graphengine.L2Squared
	if len(call.Args) >= 3 {
		opts, err := pd.resolveMap(call.Args[2], bindings)
		if err == nil {
			if m, ok := opts["metric"].(string); ok {
				metric = graphengine.ParseDistanceMetric(m)
			}
		}
	}

	rows := make([]map[string]interface{}, 0)
	for _, srcID := range sourceIDs {
		srcEmb, err := pd.db.Embeddings.Get(srcID)
		if err != nil || srcEmb == nil {
			continue
		}
		srcNode, _ := pd.db.GetNode(srcID)
		for _, tgtID := range targetIDs {
			tgtEmb, err := pd.db.Embeddings.Get(tgtID)
			if err != nil || tgtEmb == nil {
				continue
			}
			tgtNode, _ := pd.db.GetNode(tgtID)
			dist, err := graphengine.ComputeDistance(srcEmb, tgtEmb, metric)
			if err != nil {
				continue
			}
			rows = append(rows, map[string]interface{}{
				"source":   srcNode,
				"target":   tgtNode,
				"distance": dist,
			})
		}
	}

	return &cypherparser.CypherResult{
		Columns: []string{"source", "target", "distance"},
		Rows:    rows,
	}, nil
}

func (pd *procedureDispatcher) execDistanceByEmbedding(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 2 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.distance.byEmbedding requires at least 2 arguments")
	}

	nodeIDs, err := pd.resolveNodeIDs(call.Args[0], bindings)
	if err != nil {
		return nil, err
	}

	opts, err := pd.resolveMap(call.Args[1], bindings)
	if err != nil {
		return nil, fmt.Errorf("ValidationException: second argument must be a map with 'embedding' and optional 'metric'")
	}

	queryEmbRaw, ok := opts["embedding"]
	if !ok {
		return nil, fmt.Errorf("ValidationException: 'embedding' is required in options map")
	}
	queryEmb, err := toFloat64Slice(queryEmbRaw)
	if err != nil {
		return nil, fmt.Errorf("ValidationException: embedding must be a list of floats: %v", err)
	}

	metric := graphengine.L2Squared
	if m, ok := opts["metric"].(string); ok {
		metric = graphengine.ParseDistanceMetric(m)
	}

	rows := make([]map[string]interface{}, 0)
	for _, id := range nodeIDs {
		tgtEmb, err := pd.db.Embeddings.Get(id)
		if err != nil || tgtEmb == nil {
			continue
		}
		tgtNode, _ := pd.db.GetNode(id)
		dist, err := graphengine.ComputeDistance(queryEmb, tgtEmb, metric)
		if err != nil {
			continue
		}
		rows = append(rows, map[string]interface{}{
			"target":   tgtNode,
			"distance": dist,
		})
	}

	return &cypherparser.CypherResult{
		Columns: []string{"target", "distance"},
		Rows:    rows,
	}, nil
}

func (pd *procedureDispatcher) execTopKByNode(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 1 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.topK.byNode requires at least 1 argument")
	}

	nodeIDs, err := pd.resolveNodeIDs(call.Args[0], bindings)
	if err != nil {
		return nil, err
	}

	k := 10
	if len(call.Args) >= 2 {
		opts, err := pd.resolveMap(call.Args[1], bindings)
		if err == nil {
			if v, ok := opts["topK"]; ok {
				k = toInt(v)
				if k <= 0 {
					k = 10
				}
			}
		}
	}

	rows := make([]map[string]interface{}, 0)
	for _, srcID := range nodeIDs {
		srcEmb, err := pd.db.Embeddings.Get(srcID)
		if err != nil || srcEmb == nil {
			continue
		}
		srcNode, _ := pd.db.GetNode(srcID)

		results, err := pd.db.Embeddings.TopK(srcEmb, k, nil)
		if err != nil {
			return nil, err
		}

		for _, r := range results {
			if r.NodeID == srcID {
				continue
			}
			node, _ := pd.db.GetNode(r.NodeID)
			rows = append(rows, map[string]interface{}{
				"source": srcNode,
				"node":   node,
				"score":  r.Score,
			})
		}
	}

	return &cypherparser.CypherResult{
		Columns: []string{"source", "node", "score"},
		Rows:    rows,
	}, nil
}

func (pd *procedureDispatcher) execTopKByEmbedding(call *cypherparser.CypherCall, bindings map[string]interface{}) (*cypherparser.CypherResult, error) {
	if len(call.Args) < 1 {
		return nil, fmt.Errorf("ValidationException: neptune.algo.vectors.topK.byEmbedding requires at least 1 argument")
	}

	opts, err := pd.resolveMap(call.Args[0], bindings)
	if err != nil {
		return nil, fmt.Errorf("ValidationException: argument must be a map with 'embedding' and optional 'topK', 'vertexFilter'")
	}

	queryEmbRaw, ok := opts["embedding"]
	if !ok {
		return nil, fmt.Errorf("ValidationException: 'embedding' is required")
	}
	queryEmb, err := toFloat64Slice(queryEmbRaw)
	if err != nil {
		return nil, fmt.Errorf("ValidationException: embedding must be a list of floats: %v", err)
	}

	k := 10
	if v, ok := opts["topK"]; ok {
		k = toInt(v)
		if k <= 0 {
			k = 10
		}
	}

	var filter graphengine.VertexFilter
	if vfStr, ok := opts["vertexFilter"].(string); ok && vfStr != "" {
		filter, err = graphengine.ParseVertexFilter(vfStr)
		if err != nil {
			return nil, fmt.Errorf("ValidationException: invalid vertexFilter: %v", err)
		}
	}

	var filterFn func(graphengine.NodeID) bool
	if filter != nil {
		filterFn = func(id graphengine.NodeID) bool {
			node, err := pd.db.GetNode(id)
			if err != nil {
				return false
			}
			return filter.Evaluate(node)
		}
	}

	results, err := pd.db.Embeddings.TopK(queryEmb, k, filterFn)
	if err != nil {
		return nil, err
	}

	rows := make([]map[string]interface{}, 0, len(results))
	for _, r := range results {
		node, _ := pd.db.GetNode(r.NodeID)
		rows = append(rows, map[string]interface{}{
			"node":  node,
			"score": r.Score,
		})
	}

	return &cypherparser.CypherResult{
		Columns: []string{"node", "score"},
		Rows:    rows,
	}, nil
}

// resolveNodeIDs extracts node IDs from a CALL argument. Supports:
//   - A bound variable referencing a graphengine.Node or []graphengine.Node
//   - A string literal (treated as node ~id)
//   - A list of string literals
func (pd *procedureDispatcher) resolveNodeIDs(arg cypherparser.Expression, bindings map[string]interface{}) ([]graphengine.NodeID, error) {
	if arg.Kind == cypherparser.ExprVarRef {
		val, ok := bindings[arg.Variable]
		if !ok {
			return nil, fmt.Errorf("ValidationException: unbound variable %s", arg.Variable)
		}
		return pd.nodeIDsFromValue(val)
	}

	if arg.Kind == cypherparser.ExprLiteral {
		switch v := arg.LitValue.(type) {
		case string:
			return pd.lookupNodeByID(v)
		}
	}

	if arg.Kind == cypherparser.ExprListLiteral {
		var ids []graphengine.NodeID
		for _, item := range arg.ListItems {
			subIDs, err := pd.resolveNodeIDs(item, bindings)
			if err != nil {
				return nil, err
			}
			ids = append(ids, subIDs...)
		}
		return ids, nil
	}

	return nil, fmt.Errorf("ValidationException: cannot resolve node from argument")
}

func (pd *procedureDispatcher) nodeIDsFromValue(val interface{}) ([]graphengine.NodeID, error) {
	switch v := val.(type) {
	case *graphengine.Node:
		return []graphengine.NodeID{v.ID}, nil
	case graphengine.Node:
		return []graphengine.NodeID{v.ID}, nil
	case []*graphengine.Node:
		ids := make([]graphengine.NodeID, 0, len(v))
		for _, n := range v {
			if n != nil {
				ids = append(ids, n.ID)
			}
		}
		return ids, nil
	case string:
		return pd.lookupNodeByID(v)
	case []string:
		var ids []graphengine.NodeID
		for _, s := range v {
			sub, err := pd.lookupNodeByID(s)
			if err != nil {
				return nil, err
			}
			ids = append(ids, sub...)
		}
		return ids, nil
	case []interface{}:
		var ids []graphengine.NodeID
		for _, item := range v {
			sub, err := pd.nodeIDsFromValue(item)
			if err != nil {
				return nil, err
			}
			ids = append(ids, sub...)
		}
		return ids, nil
	default:
		return nil, fmt.Errorf("ValidationException: unsupported node reference type %T", val)
	}
}

func (pd *procedureDispatcher) lookupNodeByID(idStr string) ([]graphengine.NodeID, error) {
	id, err := parseNodeIDString(idStr)
	if err != nil {
		return nil, fmt.Errorf("ValidationException: invalid node id '%s': %v", idStr, err)
	}
	if exists, _ := pd.db.NodeExists(id); !exists {
		return nil, fmt.Errorf("ResourceNotFoundException: node with id %s not found", idStr)
	}
	return []graphengine.NodeID{id}, nil
}

func parseNodeIDString(s string) (graphengine.NodeID, error) {
	var id uint64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			id = id*10 + uint64(c-'0')
		} else {
			return 0, fmt.Errorf("not a valid numeric node id")
		}
	}
	if id == 0 {
		return 0, fmt.Errorf("node id must be > 0")
	}
	return graphengine.NodeID(id), nil
}

func (pd *procedureDispatcher) resolveFloatList(arg cypherparser.Expression, bindings map[string]interface{}) ([]float64, error) {
	if arg.Kind == cypherparser.ExprVarRef {
		val, ok := bindings[arg.Variable]
		if !ok {
			return nil, fmt.Errorf("unbound variable %s", arg.Variable)
		}
		return toFloat64Slice(val)
	}

	if arg.Kind == cypherparser.ExprListLiteral {
		emb := make([]float64, 0, len(arg.ListItems))
		for _, item := range arg.ListItems {
			if item.Kind == cypherparser.ExprLiteral {
				f := toFloat64Value(item.LitValue)
				emb = append(emb, f)
			} else if item.Kind == cypherparser.ExprVarRef {
				val, ok := bindings[item.Variable]
				if !ok {
					return nil, fmt.Errorf("unbound variable %s", item.Variable)
				}
				slice, err := toFloat64Slice(val)
				if err != nil {
					return nil, err
				}
				emb = append(emb, slice...)
			} else {
				return nil, fmt.Errorf("embedding list item must be a literal or variable")
			}
		}
		return emb, nil
	}

	return nil, fmt.Errorf("expected a list of floats")
}

func (pd *procedureDispatcher) resolveMap(arg cypherparser.Expression, bindings map[string]interface{}) (map[string]interface{}, error) {
	if arg.Kind == cypherparser.ExprMapLiteral {
		m := make(map[string]interface{})
		for _, pair := range arg.MapPairs {
			key := pair.Key
			val := evalExprValue(pair.Value, bindings)
			m[key] = val
		}
		return m, nil
	}

	if arg.Kind == cypherparser.ExprVarRef {
		val, ok := bindings[arg.Variable]
		if ok {
			if m, ok := val.(map[string]interface{}); ok {
				return m, nil
			}
		}
	}

	return nil, fmt.Errorf("expected a map expression")
}

func evalExprValue(expr cypherparser.Expression, bindings map[string]interface{}) interface{} {
	switch expr.Kind {
	case cypherparser.ExprLiteral:
		return expr.LitValue
	case cypherparser.ExprVarRef:
		if val, ok := bindings[expr.Variable]; ok {
			return val
		}
		return nil
	case cypherparser.ExprListLiteral:
		items := make([]interface{}, 0, len(expr.ListItems))
		for _, item := range expr.ListItems {
			items = append(items, evalExprValue(item, bindings))
		}
		return items
	case cypherparser.ExprMapLiteral:
		m := make(map[string]interface{})
		for _, pair := range expr.MapPairs {
			m[pair.Key] = evalExprValue(pair.Value, bindings)
		}
		return m
	default:
		return nil
	}
}

func toFloat64Slice(val interface{}) ([]float64, error) {
	switch v := val.(type) {
	case []float64:
		return v, nil
	case []interface{}:
		result := make([]float64, 0, len(v))
		for _, item := range v {
			result = append(result, toFloat64Value(item))
		}
		return result, nil
	default:
		return nil, fmt.Errorf("expected a list, got %T", val)
	}
}

func toFloat64Value(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		return 0
	}
}

func toInt(val interface{}) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

// executeCypherQuery handles the full ExecuteQuery flow: parse, dispatch,
// track query state, and return results.
func executeCypherQuery(ctx context.Context, s *NeptuneGraphService, reqCtx *request.RequestContext, req *request.ParsedRequest, graphID string, entry *engineEntry, store *ngstore.NeptuneGraphStore) (interface{}, error) {
	var params struct {
		Query      string `json:"query"`
		Language   string `json:"language"`
		Parameters string `json:"parameters"`
	}
	if err := json.Unmarshal(req.Body, &params); err != nil {
		return nil, newValidationException("BAD_REQUEST", "invalid request body")
	}
	if params.Query == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "query is required")
	}

	lang := strings.ToUpper(params.Language)
	if lang == "" {
		lang = "CYPHER"
	}
	if lang != "CYPHER" && lang != "OPENCYPHER" && lang != "OPEN_CYPHER" && lang != "GREMLIN" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("unsupported language: %s", lang))
	}

	queryID := generateID("q-")
	now := time.Now().UTC()

	queryRecord := &ngstore.QueryRecord{
		Id:          queryID,
		QueryString: params.Query,
		Language:    lang,
		State:       "RUNNING",
		GraphId:     graphID,
		StartedAt:   now,
	}
	if err := store.CreateQuery(queryRecord); err != nil {
		logs.Warn("failed to create query record", logs.Err(err))
	}

	finaliseQuery := func(state string) {
		elapsed := int32(time.Since(now).Milliseconds())
		queryRecord.Elapsed = elapsed
		queryRecord.State = state
		if err := store.UpdateQuery(queryRecord); err != nil {
			logs.Warn("failed to update query state", logs.Err(err))
		}
	}

	if lang == "GREMLIN" {
		finaliseQuery("FAILED")
		return nil, newValidationException("UNSUPPORTED_OPERATION", "Gremlin queries are not supported")
	}

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		finaliseQuery("FAILED")
		return nil, newValidationException("MALFORMED_QUERY", err.Error())
	}

	var execErr error
	var result interface{}

	switch {
	case parsed.Call != nil:
		result, execErr = executeCallQuery(parsed.Call, entry, store, graphID)
	case parsed.Write != nil:
		var r *cypherparser.CypherResult
		r, execErr = cypherparser.ExecuteWrite(ctx, entry.db, entry.db, parsed.Write, nil)
		if r != nil {
			result = map[string]interface{}{"results": r}
		}
	case parsed.Merge != nil:
		var r *cypherparser.CypherResult
		r, execErr = cypherparser.ExecuteMerge(ctx, entry.db, entry.db, parsed.Merge, nil)
		if r != nil {
			result = map[string]interface{}{"results": r}
		}
	case parsed.Read != nil:
		if len(parsed.Read.Set) > 0 || len(parsed.Read.Delete) > 0 || len(parsed.Read.Remove) > 0 || parsed.Read.Create != nil {
			var r *cypherparser.CypherResult
			r, execErr = cypherparser.ExecuteQueryWrite(ctx, entry.db, entry.db, parsed.Read, nil)
			if r != nil {
				result = map[string]interface{}{"results": r}
			}
		} else {
			var r *cypherparser.CypherResult
			r, execErr = cypherparser.Execute(ctx, entry.db, parsed.Read, nil)
			if r != nil {
				result = map[string]interface{}{"results": r}
			}
		}
	default:
		finaliseQuery("FAILED")
		return nil, newValidationException("MALFORMED_QUERY", "unsupported query type")
	}

	if execErr != nil {
		finaliseQuery("FAILED")
		return nil, newValidationException("MALFORMED_QUERY", execErr.Error())
	}

	finaliseQuery("COMPLETE")
	return result, nil
}

func executeCallQuery(call *cypherparser.CypherCall, entry *engineEntry, store *ngstore.NeptuneGraphStore, graphID string) (interface{}, error) {
	graph, err := store.GetGraph(graphID)
	if err != nil {
		return nil, newResourceNotFoundException("graph", graphID)
	}

	var dim int32
	if graph.VectorSearchConfiguration != nil {
		dim = graph.VectorSearchConfiguration.Dimension
	}

	pd := newProcedureDispatcher(entry.db, dim)
	result, err := pd.ExecuteCall(call, nil)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "ValidationException") {
			return nil, newValidationException("ILLEGAL_ARGUMENT", strings.TrimPrefix(errStr, "ValidationException: "))
		}
		if strings.Contains(errStr, "ResourceNotFoundException") {
			return nil, newResourceNotFoundException("node", "")
		}
		return nil, newHTTPError(500, "InternalServerException", errStr)
	}

	return map[string]interface{}{"results": result}, nil
}

// validateEmbedding checks that a float slice does not contain INF or NaN.
func validateEmbedding(emb []float64) error {
	for _, v := range emb {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return fmt.Errorf("embedding contains invalid float value (INF or NaN)")
		}
	}
	return nil
}
