package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/pkg/gremlinparser"
)

// ExecuteGremlinQuery parses and executes a Gremlin traversal query against the
// graph engine, returning results in the AWS response shape.
func (s *NeptuneDataService) ExecuteGremlinQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	body := req.Body
	var params struct {
		Gremlin string `json:"gremlin"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	qid := generateQueryID()
	s.trackQuery(qid, params.Gremlin, "gremlin")

	reader := reqCtx.GraphReader()
	writer := reqCtx.GraphWriter()
	parsed, err := gremlinparser.Parse(params.Gremlin)
	if err != nil {
		s.resolveQuery(qid, nil, err)
		return nil, malformedQuery(err.Error())
	}

	result, execErr := gremlinparser.ExecuteQuery(ctx, reader, writer, parsed, nil)
	s.resolveQuery(qid, result, execErr)
	if execErr != nil {
		return nil, malformedQuery(execErr.Error())
	}

	return map[string]interface{}{
		"requestId": qid,
		"status": map[string]interface{}{
			"code":       float64(200),
			"message":    "",
			"attributes": map[string]interface{}{},
		},
		"result": result,
		"meta": map[string]interface{}{
			"requestId": qid,
		},
	}, nil
}

// ExecuteGremlinExplainQuery returns an explain plan for a Gremlin query without
// executing it.
func (s *NeptuneDataService) ExecuteGremlinExplainQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Gremlin string `json:"gremlin"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	return map[string]interface{}{
		"code": "gremlin-traversal-explanation",
		"plan": explainGremlinQuery(params.Gremlin),
	}, nil
}

// ExecuteGremlinProfileQuery returns a profiled explain plan for a Gremlin query.
func (s *NeptuneDataService) ExecuteGremlinProfileQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Gremlin string `json:"gremlin"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	return map[string]interface{}{
		"code":    "gremlin-traversal-profile",
		"profile": profileGremlinQuery(params.Gremlin),
	}, nil
}

func (s *NeptuneDataService) GetGremlinQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	qs, ok := s.queryState[queryId]
	if !ok {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}

	evalStats := map[string]interface{}{
		"cancelled": qs.Status == "cancelled",
		"elapsed":   elapsedMs(qs.StartedAt, qs.EndedAt),
		"waited":    0,
	}

	return map[string]interface{}{
		"queryId":        qs.ID,
		"queryString":    qs.Query,
		"queryEvalStats": evalStats,
	}, nil
}

func (s *NeptuneDataService) ListGremlinQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	includeWaiting := request.GetBoolParam(req.Parameters, "includeWaiting")
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpired()

	var queries []interface{}
	var acceptedCount, runningCount int32

	for _, qs := range s.queryState {
		if qs.QueryType != "gremlin" {
			continue
		}
		if qs.Status == "waiting" && !includeWaiting {
			continue
		}
		entry := map[string]interface{}{
			"queryId":     qs.ID,
			"queryString": qs.Query,
		}
		switch qs.Status {
		case "running":
			runningCount++
		default:
			acceptedCount++
		}
		queries = append(queries, entry)
	}

	return map[string]interface{}{
		"queries":            queries,
		"acceptedQueryCount": acceptedCount,
		"runningQueryCount":  runningCount,
	}, nil
}

func (s *NeptuneDataService) CancelGremlinQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	qs, ok := s.queryState[queryId]
	if !ok {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}
	qs.Status = "cancelled"
	qs.EndedAt = time.Now()

	return map[string]interface{}{
		"status": "200 OK",
	}, nil
}
