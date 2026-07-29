package neptunegraph

import (
	"context"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
	"vorpalstacks/internal/utils/timeutils"
)

// populateDetailedStats computes property-level statistics for DETAILED mode
// GetGraphSummary responses. It iterates all nodes and edges once, collecting:
//   - numNodeProperties / numEdgeProperties: distinct property name counts
//   - totalNodePropertyValues / totalEdgePropertyValues: total key-value pairs
//   - nodeProperties / edgeProperties: per-label property name → count maps
//   - nodeStructures / edgeStructures: per-label structural summaries
func populateDetailedStats(db *graphengine.DB, summary *ngstore.GraphDataSummary) {
	nodePropNames := make(map[string]bool)
	edgePropNames := make(map[string]bool)
	var totalNodePropVals, totalEdgePropVals int64

	// Per-label property accumulators.
	nodeLabelProps := make(map[string]map[string]int64)
	edgeLabelProps := make(map[string]map[string]int64)
	nodeLabelCounts := make(map[string]int64)
	edgeLabelCounts := make(map[string]int64)

	// Track distinct outgoing edge labels per node label.
	nodeLabelOutLabels := make(map[string]map[string]bool)

	_ = db.ForEachNode(func(node *graphengine.Node) error {
		for _, label := range node.Labels {
			nodeLabelCounts[label]++
		}
		for k := range node.Props {
			if k == "~id" {
				continue
			}
			nodePropNames[k] = true
			totalNodePropVals++
			for _, label := range node.Labels {
				if nodeLabelProps[label] == nil {
					nodeLabelProps[label] = make(map[string]int64)
				}
				nodeLabelProps[label][k]++
			}
		}
		return nil
	})

	_ = db.ForEachEdge(func(edge *graphengine.Edge) error {
		edgeLabelCounts[edge.Label]++
		if edge.Props != nil {
			for k := range edge.Props {
				edgePropNames[k] = true
				totalEdgePropVals++
				if edgeLabelProps[edge.Label] == nil {
					edgeLabelProps[edge.Label] = make(map[string]int64)
				}
				edgeLabelProps[edge.Label][k]++
			}
		}

		// Track outgoing edge labels per source node label.
		fromNode, err := db.GetNode(edge.From)
		if err == nil && fromNode != nil {
			for _, label := range fromNode.Labels {
				if nodeLabelOutLabels[label] == nil {
					nodeLabelOutLabels[label] = make(map[string]bool)
				}
				nodeLabelOutLabels[label][edge.Label] = true
			}
		}
		return nil
	})

	summary.NumNodeProperties = int64Ptr(int64(len(nodePropNames)))
	summary.NumEdgeProperties = int64Ptr(int64(len(edgePropNames)))
	summary.TotalNodePropertyValues = int64Ptr(totalNodePropVals)
	summary.TotalEdgePropertyValues = int64Ptr(totalEdgePropVals)

	// nodeProperties: per-label sorted property name → count map
	if len(nodeLabelProps) > 0 {
		labels := make([]string, 0, len(nodeLabelProps))
		for label := range nodeLabelProps {
			labels = append(labels, label)
		}
		sort.Strings(labels)
		for _, label := range labels {
			props := nodeLabelProps[label]
			m := make(map[string]int64)
			propNames := make([]string, 0, len(props))
			for p := range props {
				propNames = append(propNames, p)
			}
			sort.Strings(propNames)
			for _, p := range propNames {
				m[p] = props[p]
			}
			summary.NodeProperties = append(summary.NodeProperties, m)
		}
	}

	// edgeProperties: per-label sorted property name → count map
	if len(edgeLabelProps) > 0 {
		labels := make([]string, 0, len(edgeLabelProps))
		for label := range edgeLabelProps {
			labels = append(labels, label)
		}
		sort.Strings(labels)
		for _, label := range labels {
			props := edgeLabelProps[label]
			m := make(map[string]int64)
			propNames := make([]string, 0, len(props))
			for p := range props {
				propNames = append(propNames, p)
			}
			sort.Strings(propNames)
			for _, p := range propNames {
				m[p] = props[p]
			}
			summary.EdgeProperties = append(summary.EdgeProperties, m)
		}
	}

	// nodeStructures: per-label structure with count, properties, outgoing labels
	if len(nodeLabelCounts) > 0 {
		labels := make([]string, 0, len(nodeLabelCounts))
		for label := range nodeLabelCounts {
			labels = append(labels, label)
		}
		sort.Strings(labels)
		for _, label := range labels {
			ns := ngstore.NodeStructure{
				Count: int64Ptr(nodeLabelCounts[label]),
			}
			if props, ok := nodeLabelProps[label]; ok {
				propNames := make([]string, 0, len(props))
				for p := range props {
					propNames = append(propNames, p)
				}
				sort.Strings(propNames)
				ns.NodeProperties = propNames
			}
			if outLabels, ok := nodeLabelOutLabels[label]; ok && len(outLabels) > 0 {
				sorted := make([]string, 0, len(outLabels))
				for ol := range outLabels {
					sorted = append(sorted, ol)
				}
				sort.Strings(sorted)
				ns.DistinctOutgoingEdgeLabels = sorted
			}
			summary.NodeStructures = append(summary.NodeStructures, ns)
		}
	}

	// edgeStructures: per-label structure with count, properties
	if len(edgeLabelCounts) > 0 {
		labels := make([]string, 0, len(edgeLabelCounts))
		for label := range edgeLabelCounts {
			labels = append(labels, label)
		}
		sort.Strings(labels)
		for _, label := range labels {
			es := ngstore.EdgeStructure{
				Count: int64Ptr(edgeLabelCounts[label]),
			}
			if props, ok := edgeLabelProps[label]; ok {
				propNames := make([]string, 0, len(props))
				for p := range props {
					propNames = append(propNames, p)
				}
				sort.Strings(propNames)
				es.EdgeProperties = propNames
			}
			summary.EdgeStructures = append(summary.EdgeStructures, es)
		}
	}
}

func resolveGraphIdentifier(params map[string]interface{}) string {
	if id := request.GetStringParam(params, "graphIdentifier"); id != "" {
		return id
	}
	return request.GetStringParam(params, "graphidentifier")
}

// ExecuteQuery runs a Cypher query against the specified graph's query engine.
func (s *NeptuneGraphService) ExecuteQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := resolveGraphIdentifier(req.Parameters)
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	s.enginesMu.Lock()
	entry, ok := s.activeEngines[graphID]
	if !ok || entry.stopped {
		s.enginesMu.Unlock()
		return nil, newValidationException("UNSUPPORTED_OPERATION", "graph is not available")
	}
	entry.wg.Add(1)
	s.enginesMu.Unlock()
	defer entry.wg.Done()

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	return executeCypherQuery(ctx, s, reqCtx, req, graphID, entry, store)
}

// GetQuery retrieves the details and results of a previously executed query.
func (s *NeptuneGraphService) GetQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryID := request.GetStringParam(req.Parameters, "queryId")
	if queryID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "queryId")
	}

	graphID := resolveGraphIdentifier(req.Parameters)
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	query, err := store.GetQuery(graphID, queryID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("query", queryID)
		}
		return nil, err
	}

	return queryToResponse(query), nil
}

// ListQueries returns query records for a graph, optionally filtered by state.
func (s *NeptuneGraphService) ListQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	graphID := resolveGraphIdentifier(req.Parameters)
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	maxResults := request.GetIntParam(req.Parameters, "maxResults")
	if maxResults < 1 || maxResults > 100 {
		maxResults = 100
	}

	stateFilter := request.GetStringParam(req.Parameters, "state")
	if stateFilter != "" && stateFilter != "ALL" {
		allQueries, err := store.ListQueries(graphID, 0)
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, maxResults)
		for _, q := range allQueries {
			if q.State != stateFilter {
				continue
			}
			items = append(items, queryToResponse(q))
			if len(items) >= maxResults {
				break
			}
		}
		return map[string]interface{}{
			"queries": items,
		}, nil
	}

	queries, err := store.ListQueries(graphID, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(queries))
	for _, q := range queries {
		items = append(items, queryToResponse(q))
	}

	return map[string]interface{}{
		"queries": items,
	}, nil
}

// CancelQuery cancels a running query by transitioning it to CANCELLING state.
// Per Smithy QueryState model, only RUNNING/WAITING/CANCELLING are valid states.
// Terminal queries are deleted (not stored with a terminal state), so a query
// that still exists is either RUNNING or CANCELLING.
func (s *NeptuneGraphService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	queryID := request.GetStringParam(req.Parameters, "queryId")
	if queryID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "queryId")
	}

	graphID := resolveGraphIdentifier(req.Parameters)
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	query, err := store.GetQuery(graphID, queryID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("query", queryID)
		}
		return nil, err
	}

	// If already CANCELLING, return idempotently.
	if query.State == "CANCELLING" || query.State == "WAITING" {
		return queryToResponse(query), nil
	}

	query.State = "CANCELLING"
	if err := store.UpdateQuery(query); err != nil {
		logs.Warn("failed to cancel query", logs.String("queryId", queryID), logs.Err(err))
		return nil, err
	}

	return queryToResponse(query), nil
}

// GetGraphSummary returns structural statistics about the specified graph's data.
func (s *NeptuneGraphService) GetGraphSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	graphID := resolveGraphIdentifier(req.Parameters)
	if graphID == "" {
		return nil, newValidationException("ILLEGAL_ARGUMENT", "graphIdentifier header required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	graph, err := s.resolveGraphIdentifier(store, graphID)
	if err != nil {
		if ngstore.IsNotFound(err) {
			return nil, newResourceNotFoundException("graph", graphID)
		}
		return nil, err
	}
	graphID = graph.Id

	s.enginesMu.Lock()
	entry, ok := s.activeEngines[graphID]
	if !ok {
		s.enginesMu.Unlock()
		return nil, newResourceNotFoundException("graph", graphID)
	}
	if entry.stopped {
		s.enginesMu.Unlock()
		return nil, newResourceNotFoundException("graph", graphID)
	}
	entry.wg.Add(1)
	s.enginesMu.Unlock()
	defer entry.wg.Done()

	entry.mu.RLock()
	stats := entry.db.Stats()
	entry.mu.RUnlock()
	now := time.Now().UTC()

	summary := &ngstore.GraphDataSummary{
		NumNodes:      int64Ptr(stats.NodeCount),
		NumEdges:      int64Ptr(stats.EdgeCount),
		NumNodeLabels: int64Ptr(int64(len(stats.LabelCounts))),
		NumEdgeLabels: int64Ptr(int64(len(stats.RelCounts))),
	}

	if len(stats.LabelCounts) > 0 {
		labels := make([]string, 0, len(stats.LabelCounts))
		for label := range stats.LabelCounts {
			labels = append(labels, label)
		}
		summary.NodeLabels = labels
	}

	if len(stats.RelCounts) > 0 {
		labels := make([]string, 0, len(stats.RelCounts))
		for label := range stats.RelCounts {
			labels = append(labels, label)
		}
		summary.EdgeLabels = labels
	}

	// DETAILED mode: compute property statistics by iterating all nodes
	// and edges. This is O(n) but GetGraphSummary is not a hot path.
	mode := request.GetStringParam(req.Parameters, "mode")
	if strings.ToUpper(mode) == "DETAILED" {
		populateDetailedStats(entry.db, summary)
	}

	return map[string]interface{}{
		"graphSummary":                  graphDataSummaryToResponse(summary),
		"lastStatisticsComputationTime": now.Format(timeutils.ISO8601UTCFormat),
		"version":                       "1.0",
	}, nil
}
