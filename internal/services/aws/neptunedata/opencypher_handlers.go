package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

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

// GetOpenCypherQueryStatus returns the current status and evaluation statistics
// of a previously submitted OpenCypher query.
func (s *NeptuneDataService) GetOpenCypherQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		"cancelled": qr.GetStatus() == "cancelled",
		"elapsed":   elapsed,
		"waited":    0,
	}

	return map[string]interface{}{
		"queryId":        qr.GetQueryId(),
		"queryString":    qr.GetQueryType(),
		"queryEvalStats": evalStats,
	}, nil
}

// ListOpenCypherQueries returns all submitted OpenCypher queries, optionally
// including those in a waiting state.
func (s *NeptuneDataService) ListOpenCypherQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		if qr.GetQueryType() != "opencypher" {
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

// CancelOpenCypherQuery cancels a running OpenCypher query and marks its status
// as cancelled.
func (s *NeptuneDataService) CancelOpenCypherQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	queryId := getPathParam(req, "queryId")
	silent := request.GetBoolParam(req.Parameters, "silent")
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

	if silent {
		return map[string]interface{}{}, nil
	}

	return map[string]interface{}{
		"payload": true,
		"status":  "200 OK",
	}, nil
}
