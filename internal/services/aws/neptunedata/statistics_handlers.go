package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
)

// GetPropertygraphStatistics returns auto-computed property graph statistics
// including node and edge counts grouped by label.
func (s *NeptuneDataService) GetPropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = req

	s.mu.RLock()
	statsDisabled := s.statsDisabled
	autoCompute := s.autoComputeEnabled
	s.mu.RUnlock()

	if !statsDisabled {
		s.refreshStatistics(reqCtx)
	}
	stats := s.getStats(reqCtx.GetRegion())
	nodeCount, _, labelCounts, relCounts := stats.snapshot()

	sigCount := int64(len(labelCounts))
	predCount := int64(len(relCounts))

	result := map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"active":       !statsDisabled,
			"autoCompute":  autoCompute,
			"date":         time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			"note":         "Automatically computed",
			"statisticsId": "auto-statistics",
			"signatureInfo": map[string]interface{}{
				"signatureCount": sigCount,
				"instanceCount":  nodeCount,
				"predicateCount": predCount,
			},
		},
	}
	return result, nil
}

// ManagePropertygraphStatistics enables, disables, or refreshes auto-computed property graph statistics.
func (s *NeptuneDataService) ManagePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
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
		s.refreshStatistics(reqCtx)
		return map[string]interface{}{
			"status":  "200",
			"payload": map[string]interface{}{"statisticsId": generateQueryID()},
		}, nil
	default:
		return nil, invalidParameter(fmt.Sprintf("unknown mode: %s", params.Mode))
	}
}

// DeletePropertygraphStatistics clears all property graph statistics for the current region.
func (s *NeptuneDataService) DeletePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
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
	nodeCount, edgeCount, labelCounts, relCounts := stats.snapshot()

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
		summary["numNodeProperties"] = 0
		summary["numEdgeProperties"] = 0
		summary["totalNodePropertyValues"] = 0
		summary["totalEdgePropertyValues"] = 0
		summary["nodeLabels"] = nodeLabels
		summary["edgeLabels"] = edgeLabels
		summary["nodeProperties"] = []interface{}{}
		summary["edgeProperties"] = []interface{}{}
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
			"lastStatisticsComputationTime": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		},
	}, nil
}

// GetPropertygraphStream returns an empty property graph change stream (stub implementation).
func (s *NeptuneDataService) GetPropertygraphStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = request.GetIntParam(req.Parameters, "commitNum")
	_ = request.GetStringParam(req.Parameters, "iteratorType")
	_ = request.GetIntParam(req.Parameters, "limit")
	_ = request.GetIntParam(req.Parameters, "opNum")
	return map[string]interface{}{
		"format":                   "PG_JSON",
		"lastEventId":              map[string]string{},
		"totalRecords":             0,
		"records":                  []interface{}{},
		"lastTrxTimestampInMillis": int64(0),
	}, nil
}
