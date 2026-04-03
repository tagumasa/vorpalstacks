package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/pkg/cypherparser"
)

// ExecuteOpenCypherQuery parses and executes an OpenCypher query against the
// graph engine, dispatching to the appropriate executor based on query type.
func (s *NeptuneDataService) ExecuteOpenCypherQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	body := req.Body
	var params struct {
		Query      string          `json:"query"`
		Parameters json.RawMessage `json:"parameters"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Query == "" {
		return nil, missingParameter("query")
	}

	var cypherParams map[string]any
	if len(params.Parameters) > 0 && string(params.Parameters) != "" {
		trimmed := strings.TrimSpace(string(params.Parameters))
		if len(trimmed) > 0 && trimmed[0] == '"' {
			var paramStr string
			if err := json.Unmarshal(params.Parameters, &paramStr); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
			if err := json.Unmarshal([]byte(paramStr), &cypherParams); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
		} else {
			if err := json.Unmarshal(params.Parameters, &cypherParams); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
		}
	}

	qid := generateQueryID()
	s.trackQuery(qid, params.Query, "opencypher")

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		s.resolveQuery(qid, nil, err)
		return nil, malformedQuery(err.Error())
	}

	reader := reqCtx.GraphReader()
	writer := reqCtx.GraphWriter()

	var result *cypherparser.CypherResult
	switch {
	case parsed.Write != nil:
		result, err = cypherparser.ExecuteWrite(ctx, reader, writer, parsed.Write, cypherParams)
	case parsed.Merge != nil:
		result, err = cypherparser.ExecuteMerge(ctx, reader, writer, parsed.Merge, cypherParams)
	case parsed.Read != nil && (len(parsed.Read.Set) > 0 || len(parsed.Read.Delete) > 0 || len(parsed.Read.Remove) > 0 || parsed.Read.Create != nil):
		result, err = cypherparser.ExecuteQueryWrite(ctx, reader, writer, parsed.Read, cypherParams)
	default:
		result, err = cypherparser.Execute(ctx, reader, parsed.Read, cypherParams)
	}
	s.resolveQuery(qid, result, err)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	return map[string]interface{}{
		"results": result,
	}, nil
}

// ExecuteOpenCypherExplainQuery returns an explain plan for an OpenCypher query.
// Only read queries are supported; write queries return BadRequestException.
func (s *NeptuneDataService) ExecuteOpenCypherExplainQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	body := req.Body
	var params struct {
		Query      string          `json:"query"`
		Explain    string          `json:"explain"`
		Parameters json.RawMessage `json:"parameters"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Query == "" {
		return nil, missingParameter("query")
	}

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	if parsed.Read == nil {
		return nil, badRequest("EXPLAIN is only supported for read queries")
	}

	plan := cypherparser.BuildExplainPlan(parsed.Read, reqCtx.GraphReader())
	return map[string]interface{}{
		"explain": plan,
	}, nil
}

func (s *NeptuneDataService) GetOpenCypherQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

func (s *NeptuneDataService) ListOpenCypherQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	includeWaiting := request.GetBoolParam(req.Parameters, "includeWaiting")
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpired()

	var queries []interface{}
	var acceptedCount, runningCount int32

	for _, qs := range s.queryState {
		if qs.QueryType != "opencypher" {
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

func (s *NeptuneDataService) CancelOpenCypherQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	silent := request.GetBoolParam(req.Parameters, "silent")
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

	if silent {
		return map[string]interface{}{}, nil
	}

	return map[string]interface{}{
		"payload": true,
		"status":  "200 OK",
	}, nil
}
