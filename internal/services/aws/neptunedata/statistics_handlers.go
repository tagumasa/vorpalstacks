package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage/graphengine"
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
	nodeCount, edgeCount, labelCounts, relCounts := stats.snapshot()

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

// GetPropertygraphStream returns the property graph change stream by enumerating
// all current nodes and edges from the graph store and presenting them as PG_JSON
// ADD records (snapshot-as-stream approach).
func (s *NeptuneDataService) GetPropertygraphStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx

	limit := request.GetIntParam(req.Parameters, "limit")
	if limit <= 0 {
		limit = 1000
	}

	reader, ok := reqCtx.GraphReader().(graphengine.GraphReader)
	if !ok {
		return nil, internalFailure("graph reader not available")
	}

	commitNum := request.GetIntParam(req.Parameters, "commitNum")
	opNum := request.GetIntParam(req.Parameters, "opNum")
	iteratorType := request.GetStringParam(req.Parameters, "iteratorType")

	var records []interface{}
	remaining := limit

	records, remaining = appendNodeRecords(reader, records, remaining)
	records, _ = appendEdgeRecords(reader, records, remaining)

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
	_ = reader.ForEachNode(func(node *graphengine.Node) error {
		if remaining <= 0 {
			return fmt.Errorf("limit reached")
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
	return records, remaining
}

func appendEdgeRecords(reader graphengine.GraphReader, records []interface{}, remaining int) ([]interface{}, int) {
	if remaining <= 0 {
		return records, remaining
	}
	_ = reader.ForEachEdge(func(edge *graphengine.Edge) error {
		if remaining <= 0 {
			return fmt.Errorf("limit reached")
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
	return records, remaining
}
