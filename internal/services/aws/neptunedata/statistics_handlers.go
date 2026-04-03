package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
)

// GetPropertygraphStatistics returns auto-computed property graph statistics
// including node and edge counts grouped by label.
func (s *NeptuneDataService) GetPropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = req

	s.mu.Lock()
	s.refreshStatistics(reqCtx)

	sigCount := int64(0)
	predCount := int64(0)
	for _, cnt := range s.stats.LabelCounts {
		sigCount += cnt
	}
	for _, cnt := range s.stats.RelCounts {
		predCount += cnt
	}

	result := map[string]interface{}{
		"status": "200",
		"payload": map[string]interface{}{
			"active":       true,
			"autoCompute":  true,
			"date":         time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			"note":         "Automatically computed",
			"statisticsId": "auto-statistics",
			"signatureInfo": map[string]interface{}{
				"signatureCount": sigCount,
				"instanceCount":  s.stats.NodeCount,
				"predicateCount": predCount,
			},
		},
	}
	s.mu.Unlock()
	return result, nil
}

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
		return map[string]interface{}{
			"status": "200",
		}, nil
	case "enableAutoCompute":
		return map[string]interface{}{
			"status": "200",
		}, nil
	case "refresh":
		return map[string]interface{}{
			"status":  "200",
			"payload": map[string]interface{}{"statisticsId": generateQueryID()},
		}, nil
	default:
		return nil, invalidParameter(fmt.Sprintf("unknown mode: %s", params.Mode))
	}
}

func (s *NeptuneDataService) DeletePropertygraphStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	_ = req
	s.mu.Lock()
	s.stats = &GraphStatistics{LabelCounts: make(map[string]int64), RelCounts: make(map[string]int64)}
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

	s.mu.Lock()
	s.refreshStatistics(reqCtx)

	nodeLabels := make([]string, 0, len(s.stats.LabelCounts))
	for label := range s.stats.LabelCounts {
		nodeLabels = append(nodeLabels, label)
	}

	edgeLabels := make([]string, 0, len(s.stats.RelCounts))
	for label := range s.stats.RelCounts {
		edgeLabels = append(edgeLabels, label)
	}

	summary := map[string]interface{}{
		"numNodes":      s.stats.NodeCount,
		"numEdges":      s.stats.EdgeCount,
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
		summary["nodeStructures"] = []interface{}{}
		summary["edgeStructures"] = []interface{}{}
	}

	s.mu.Unlock()

	return map[string]interface{}{
		"payload": map[string]interface{}{
			"version":                       "v1",
			"graphSummary":                  summary,
			"lastStatisticsComputationTime": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		},
	}, nil
}

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
