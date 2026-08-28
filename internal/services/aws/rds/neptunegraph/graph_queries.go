package neptunegraph

import (
	"context"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/request"
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

	summary.NumNodeProperties = proto.Int64(int64(len(nodePropNames)))
	summary.NumEdgeProperties = proto.Int64(int64(len(edgePropNames)))
	summary.TotalNodePropertyValues = proto.Int64(totalNodePropVals)
	summary.TotalEdgePropertyValues = proto.Int64(totalEdgePropVals)

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
				Count: proto.Int64(nodeLabelCounts[label]),
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
				Count: proto.Int64(edgeLabelCounts[label]),
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

	in := &GetQueryInput{
		QueryID:         request.GetStringParam(req.Parameters, "queryId"),
		GraphIdentifier: resolveGraphIdentifier(req.Parameters),
	}

	query, err := s.getQueryCore(store, in)
	if err != nil {
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

	in := &ListQueriesInput{
		GraphIdentifier: resolveGraphIdentifier(req.Parameters),
		MaxResults:      clampMaxResults(request.GetIntParam(req.Parameters, "maxResults")),
		State:           strings.ToUpper(request.GetStringParam(req.Parameters, "state")),
	}

	queries, err := s.listQueriesCore(store, in)
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
func (s *NeptuneGraphService) CancelQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CancelQueryInput{
		QueryID:         request.GetStringParam(req.Parameters, "queryId"),
		GraphIdentifier: resolveGraphIdentifier(req.Parameters),
	}

	query, err := s.cancelQueryCore(store, in)
	if err != nil {
		return nil, err
	}
	return queryToResponse(query), nil
}

// GetGraphSummary returns structural statistics about the specified graph's data.
func (s *NeptuneGraphService) GetGraphSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &GetGraphSummaryInput{
		GraphIdentifier: resolveGraphIdentifier(req.Parameters),
		Mode:            strings.ToUpper(request.GetStringParam(req.Parameters, "mode")),
	}

	res, err := s.getGraphSummaryCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"graphSummary":                  graphDataSummaryToResponse(res.Summary),
		"lastStatisticsComputationTime": res.StatsTime.Format(timeutils.ISO8601UTCFormat),
		"version":                       "1.0",
	}, nil
}
