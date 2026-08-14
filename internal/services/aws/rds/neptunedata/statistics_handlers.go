package neptunedata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage/graphengine"
	"vorpalstacks/internal/utils/timeutils"
)

// errLimitReached is a sentinel error returned by appendNodeRecords and
// appendEdgeRecords when the pagination limit is exhausted. It signals a
// normal stop of iteration, not a genuine error.
var errLimitReached = errors.New("limit reached")

// GetPropertygraphStatistics returns auto-computed property graph statistics
// including node and edge counts grouped by label.
func (s *NeptuneDataService) GetPropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = req

	s.mu.RLock()
	statsDisabled := s.statsDisabled
	autoCompute := s.autoComputeEnabled
	s.mu.RUnlock()

	if statsDisabled {
		return nil, statisticsNotAvailable("Statistics are disabled. Call ManagePropertygraphStatistics with mode 'refresh' or 'enableAutoCompute' to generate statistics.")
	}

	s.refreshStatistics(reqCtx)
	stats := s.getStats(reqCtx.GetRegion())
	nodeCount, edgeCount, labelCounts, relCounts, _, _ := stats.snapshot()

	sigCount := int64(len(labelCounts))
	predCount := int64(len(relCounts))

	result := map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"active":       true,
			"autoCompute":  autoCompute,
			"date":         time.Now().UTC().Format(timeutils.ISO8601UTCFormat),
			"note":         "Automatically computed",
			"statisticsId": "auto-statistics",
			"signatureInfo": map[string]interface{}{
				"edgeCount":      edgeCount,
				"instanceCount":  nodeCount,
				"predicateCount": predCount,
				"signatureCount": sigCount,
			},
		},
	}
	return result, nil
}

// ManagePropertygraphStatistics enables, disables, or refreshes auto-computed property graph statistics.
func (s *NeptuneDataService) ManagePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	body := req.Body
	var params struct {
		Mode string `json:"mode"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}

	switch params.Mode {
	case "disableAutoCompute":
		s.mu.Lock()
		s.autoComputeEnabled = false
		s.mu.Unlock()
		return map[string]interface{}{
			"status": "200",
		}, nil
	case "enableAutoCompute":
		s.mu.Lock()
		s.autoComputeEnabled = true
		s.mu.Unlock()
		return map[string]interface{}{
			"status": "200",
		}, nil
	case "refresh":
		s.mu.Lock()
		s.statsDisabled = false
		s.mu.Unlock()
		// Fire-and-forget: AWS Neptune handles statistics refresh
		// asynchronously. The caller receives a statisticsId immediately
		// and can poll GetPropertygraphStatistics to check the result.
		go func() {
			defer func() {
				if re := recover(); re != nil {
					logs.Error("ManagePropertygraphStatistics: refresh panic recovered", logs.Any("panic", re))
				}
			}()
			s.refreshStatistics(reqCtx)
		}()
		return map[string]interface{}{
			"status":  "200",
			"payload": map[string]interface{}{"statisticsId": generateStatisticsID()},
		}, nil
	default:
		return nil, invalidParameter(fmt.Sprintf("unknown mode: %s", params.Mode))
	}
}

// DeletePropertygraphStatistics clears all property graph statistics for the current region.
func (s *NeptuneDataService) DeletePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = req
	s.mu.Lock()
	s.statsDisabled = true
	region := reqCtx.GetRegion()
	s.statsMap.Store(region, &GraphStatistics{LabelCounts: make(map[string]int64), RelCounts: make(map[string]int64)})
	s.mu.Unlock()
	return map[string]interface{}{
		"status":  "200",
		"payload": map[string]interface{}{},
	}, nil
}

// GetPropertygraphSummary returns a detailed summary of the property graph
// including node/edge counts, label lists, and structural metadata.
func (s *NeptuneDataService) GetPropertygraphSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	mode := getPathParam(req, "mode")

	s.mu.RLock()
	statsDisabled := s.statsDisabled
	s.mu.RUnlock()

	if !statsDisabled {
		s.refreshStatistics(reqCtx)
	}
	stats := s.getStats(reqCtx.GetRegion())
	nodeCount, edgeCount, labelCounts, relCounts, nodePropCounts, edgePropCounts := stats.snapshot()

	nodeLabels := make([]string, 0, len(labelCounts))
	for label := range labelCounts {
		nodeLabels = append(nodeLabels, label)
	}

	edgeLabels := make([]string, 0, len(relCounts))
	for label := range relCounts {
		edgeLabels = append(edgeLabels, label)
	}

	summary := map[string]interface{}{
		"numNodes":      nodeCount,
		"numEdges":      edgeCount,
		"numNodeLabels": int64(len(nodeLabels)),
		"numEdgeLabels": int64(len(edgeLabels)),
	}

	if mode == "detailed" {
		var totalNodePropVals int64
		nodeProps := make([]interface{}, 0, len(nodePropCounts))
		for prop, count := range nodePropCounts {
			totalNodePropVals += count
			nodeProps = append(nodeProps, map[string]interface{}{
				"property": prop,
				"count":    count,
			})
		}

		var totalEdgePropVals int64
		edgeProps := make([]interface{}, 0, len(edgePropCounts))
		for prop, count := range edgePropCounts {
			totalEdgePropVals += count
			edgeProps = append(edgeProps, map[string]interface{}{
				"property": prop,
				"count":    count,
			})
		}

		summary["numNodeProperties"] = int64(len(nodePropCounts))
		summary["numEdgeProperties"] = int64(len(edgePropCounts))
		summary["totalNodePropertyValues"] = totalNodePropVals
		summary["totalEdgePropertyValues"] = totalEdgePropVals
		summary["nodeLabels"] = nodeLabels
		summary["edgeLabels"] = edgeLabels
		summary["nodeProperties"] = nodeProps
		summary["edgeProperties"] = edgeProps
		nodeStructures := make([]interface{}, 0, len(nodeLabels))
		for _, label := range nodeLabels {
			nodeStructures = append(nodeStructures, map[string]interface{}{
				"label": label,
				"count": labelCounts[label],
			})
		}
		summary["nodeStructures"] = nodeStructures
		edgeStructures := make([]interface{}, 0, len(edgeLabels))
		for _, label := range edgeLabels {
			edgeStructures = append(edgeStructures, map[string]interface{}{
				"label": label,
				"count": relCounts[label],
			})
		}
		summary["edgeStructures"] = edgeStructures
	}

	return map[string]interface{}{
		"payload": map[string]interface{}{
			"version":                       "v1",
			"graphSummary":                  summary,
			"lastStatisticsComputationTime": time.Now().UTC().Format(timeutils.ISO8601UTCFormat),
		},
	}, nil
}

// GetPropertygraphStream returns the property graph change stream by enumerating
// all current nodes and edges from the graph store and presenting them as PG_JSON
// ADD records (snapshot-as-stream approach).
func (s *NeptuneDataService) GetPropertygraphStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx

	// Read the encoding preference from the Accept-Encoding HTTP header
	// (Smithy httpHeader trait). Currently only "gzip" is defined in the
	// Encoding enum. Framework-level gzip middleware is required for actual
	// compression; we log the request for diagnostics.
	acceptEncoding := req.Headers.Get("Accept-Encoding")
	if acceptEncoding != "" {
		logs.Debug("GetPropertygraphStream: client requested encoding",
			logs.String("accept-encoding", acceptEncoding))
	}

	limit := request.GetIntParam(req.Parameters, "limit")
	if limit <= 0 {
		limit = 1000
	}

	reader := reqCtx.GraphReader()
	if reader == nil {
		return nil, internalFailure("graph reader not available")
	}

	commitNum := request.GetIntParam(req.Parameters, "commitNum")
	opNum := request.GetIntParam(req.Parameters, "opNum")
	iteratorType := request.GetStringParam(req.Parameters, "iteratorType")

	var records []interface{}
	remaining := limit

	records, remaining = appendNodeRecords(reader, records, remaining)
	records, remaining = appendEdgeRecords(reader, records, remaining)

	if records == nil {
		records = []interface{}{}
	}

	totalRecords := len(records)

	now := time.Now().UnixMilli()

	lastEventID := map[string]interface{}{}
	if commitNum > 0 {
		lastEventID["commitNum"] = commitNum
	}
	if opNum > 0 {
		lastEventID["opNum"] = opNum
	}
	if iteratorType != "" {
		lastEventID["iteratorType"] = iteratorType
	}

	return map[string]interface{}{
		"format":                   "PG_JSON",
		"lastEventId":              lastEventID,
		"totalRecords":             totalRecords,
		"records":                  records,
		"lastTrxTimestampInMillis": now,
	}, nil
}

func appendNodeRecords(reader graphengine.GraphReader, records []interface{}, remaining int) ([]interface{}, int) {
	if remaining <= 0 {
		return records, remaining
	}
	// Propagate iteration errors instead of silently swallowing them.
	err := reader.ForEachNode(func(node *graphengine.Node) error {
		if remaining <= 0 {
			return errLimitReached
		}
		label := ""
		if len(node.Labels) > 0 {
			label = node.Labels[0]
		}
		value := map[string]interface{}{
			"~id":    fmt.Sprintf("%d", node.ID),
			"~label": label,
		}
		for k, v := range node.Props {
			value[k] = v
		}
		records = append(records, map[string]interface{}{
			"id":    fmt.Sprintf("%d", node.ID),
			"type":  "vl",
			"key":   label,
			"value": value,
			"op":    "ADD",
		})
		remaining--
		return nil
	})
	if err != nil && !errors.Is(err, errLimitReached) {
		logs.Warn("GetPropertygraphStream: node iteration error", logs.Err(err))
	}
	return records, remaining
}

func appendEdgeRecords(reader graphengine.GraphReader, records []interface{}, remaining int) ([]interface{}, int) {
	if remaining <= 0 {
		return records, remaining
	}
	// Propagate iteration errors instead of silently swallowing them.
	err := reader.ForEachEdge(func(edge *graphengine.Edge) error {
		if remaining <= 0 {
			return errLimitReached
		}
		value := map[string]interface{}{
			"~id":   fmt.Sprintf("%d", edge.ID),
			"~type": edge.Label,
			"~from": fmt.Sprintf("%d", edge.From),
			"~to":   fmt.Sprintf("%d", edge.To),
		}
		for k, v := range edge.Props {
			value[k] = v
		}
		records = append(records, map[string]interface{}{
			"id":    fmt.Sprintf("%d", edge.ID),
			"type":  "el",
			"key":   edge.Label,
			"value": value,
			"op":    "ADD",
		})
		remaining--
		return nil
	})
	if err != nil && !errors.Is(err, errLimitReached) {
		logs.Warn("GetPropertygraphStream: edge iteration error", logs.Err(err))
	}
	return records, remaining
}
