package neptunedata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage/graphengine"
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
	if len(params.Parameters) > 0 {
		raw := bytes.TrimSpace(params.Parameters)
		if len(raw) > 0 && raw[0] == '"' {
			var paramStr string
			if err := json.Unmarshal(raw, &paramStr); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
			if err := json.Unmarshal([]byte(paramStr), &cypherParams); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
		} else {
			if err := json.Unmarshal(raw, &cypherParams); err != nil {
				return nil, invalidParameter(fmt.Sprintf("invalid parameters: %v", err))
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	qid := generateQueryID()
	s.trackQuery(store, qid, params.Query, "opencypher")

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		s.resolveQuery(store, qid, nil, err)
		return nil, malformedQuery(err.Error())
	}

	reader := reqCtx.GraphReader()
	if reader == nil {
		return nil, internalFailure("graph reader not available")
	}
	writer := reqCtx.GraphWriter()
	if writer == nil {
		return nil, internalFailure("graph writer not available")
	}

	var result *cypherparser.CypherResult
	switch {
	case parsed.DDL != nil:
		ddl, ok := reader.(graphengine.GraphDDL)
		if !ok {
			s.resolveQuery(store, qid, nil, fmt.Errorf("DDL not available"))
			return nil, internalFailure("DDL interface not available")
		}
		result, err = cypherparser.ExecuteDDL(ctx, ddl, parsed.DDL)
	case parsed.Write != nil:
		result, err = cypherparser.ExecuteWrite(ctx, reader, writer, parsed.Write, cypherParams)
	case parsed.Merge != nil:
		result, err = cypherparser.ExecuteMerge(ctx, reader, writer, parsed.Merge, cypherParams)
	case parsed.Read != nil && (len(parsed.Read.Set) > 0 || len(parsed.Read.Delete) > 0 || len(parsed.Read.Remove) > 0 || parsed.Read.Create != nil):
		result, err = cypherparser.ExecuteQueryWrite(ctx, reader, writer, parsed.Read, cypherParams)
	default:
		if parsed.Read == nil {
			s.resolveQuery(store, qid, nil, fmt.Errorf("unsupported query type"))
			return nil, badRequest("unsupported query type")
		}
		result, err = cypherparser.Execute(ctx, reader, parsed.Read, cypherParams)
	}
	s.resolveQuery(store, qid, result, err)
	if err != nil {
		return nil, failureByQuery(err.Error())
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
	if !validateExplainMode(params.Explain) {
		return nil, invalidParameter(fmt.Sprintf("invalid explain mode: %s (valid values: static, details, dynamic)", params.Explain))
	}

	parsed, err := cypherparser.Parse(params.Query)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	if parsed.Read == nil {
		return nil, badRequest("EXPLAIN is only supported for read queries")
	}

	reader := reqCtx.GraphReader()
	if reader == nil {
		return nil, internalFailure("graph reader not available")
	}
	plan := cypherparser.BuildExplainPlan(parsed.Read, reader)
	return map[string]interface{}{
		"explain": plan,
	}, nil
}

// GetOpenCypherQueryStatus returns the current status and evaluation statistics
// of a previously submitted OpenCypher query.
func (s *NeptuneDataService) GetOpenCypherQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.getQueryStatus(reqCtx, req)
}

// ListOpenCypherQueries returns all submitted OpenCypher queries, optionally
// including those in a waiting state.
func (s *NeptuneDataService) ListOpenCypherQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.listQueries(reqCtx, req, "opencypher")
}

// CancelOpenCypherQuery cancels a running OpenCypher query and marks its status
// as cancelled.
func (s *NeptuneDataService) CancelOpenCypherQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	silent := request.GetBoolParam(req.Parameters, "silent")
	return s.cancelQuery(reqCtx, req, silent, true)
}
