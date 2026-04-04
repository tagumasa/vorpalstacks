package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

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

// GetGremlinQueryStatus returns the current status and evaluation statistics of
// a previously submitted Gremlin query.
func (s *NeptuneDataService) GetGremlinQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	qr, err := s.store.GetQuery(queryId)
	if err != nil || qr == nil {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}

	var elapsed int64
	if qr.EndTime != nil && qr.StartTime != nil {
		elapsed = qr.EndTime.AsTime().Sub(qr.StartTime.AsTime()).Milliseconds()
	}

	evalStats := map[string]interface{}{
		"cancelled": qr.Status == "cancelled",
		"elapsed":   elapsed,
		"waited":    0,
	}

	return map[string]interface{}{
		"queryId":        qr.GetQueryId(),
		"queryString":    qr.GetQueryType(),
		"queryEvalStats": evalStats,
	}, nil
}

// ListGremlinQueries returns all submitted Gremlin queries, optionally
// including those in a waiting state.
func (s *NeptuneDataService) ListGremlinQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	includeWaiting := request.GetBoolParam(req.Parameters, "includeWaiting")

	queries, err := s.store.ListQueries()
	if err != nil {
		return nil, err
	}

	var result []interface{}
	var acceptedCount, runningCount int32

	for _, qr := range queries {
		if qr.GetQueryType() != "gremlin" {
			continue
		}
		if qr.GetStatus() == "waiting" && !includeWaiting {
			continue
		}
		entry := map[string]interface{}{
			"queryId":     qr.GetQueryId(),
			"queryString": qr.GetQueryType(),
		}
		if qr.GetStatus() == "running" {
			runningCount++
		} else {
			acceptedCount++
		}
		result = append(result, entry)
	}

	return map[string]interface{}{
		"queries":            result,
		"acceptedQueryCount": acceptedCount,
		"runningQueryCount":  runningCount,
	}, nil
}

// CancelGremlinQuery cancels a running Gremlin query and marks its status as
// cancelled.
func (s *NeptuneDataService) CancelGremlinQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	if queryId == "" {
		return nil, missingParameter("queryId")
	}

	qr, err := s.store.GetQuery(queryId)
	if err != nil || qr == nil {
		return nil, badRequest(fmt.Sprintf("query not found: %s", queryId))
	}
	qr.Status = "cancelled"
	qr.EndTime = timestamppb.Now()
	_ = s.store.UpdateQuery(qr)

	return map[string]interface{}{
		"status": "200 OK",
	}, nil
}
