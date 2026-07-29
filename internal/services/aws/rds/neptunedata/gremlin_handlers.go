package neptunedata

import (
	"context"
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/pkg/gremlinparser"
)

// ExecuteGremlinQuery parses and executes a Gremlin traversal query against the
// graph engine, returning results in the AWS response shape. The Accept header
// (Smithy httpHeader "accept") controls the desired serializer format.
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

	// M11: Read the Accept header (Smithy httpHeader trait on serializer).
	// Currently only GraphSON v3 (application/json) is supported for HTTP.
	// Other formats are acknowledged but not yet implemented.
	acceptHeader := req.Headers.Get("Accept")
	if acceptHeader != "" && acceptHeader != "application/vnd.gremlin-v3.0+json" && acceptHeader != "application/json" {
		logs.Debug("ExecuteGremlinQuery: unsupported Accept header, falling back to GraphSON v3",
			logs.String("accept", acceptHeader))
	}

	qid := generateQueryID()

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, internalFailure(err.Error())
	}

	s.trackQuery(store, qid, params.Gremlin, "gremlin")

	reader := reqCtx.GraphReader()
	if reader == nil {
		return nil, internalFailure("graph reader not available")
	}
	writer := reqCtx.GraphWriter()
	if writer == nil {
		return nil, internalFailure("graph writer not available")
	}
	parsed, err := gremlinparser.Parse(params.Gremlin)
	if err != nil {
		s.resolveQuery(store, qid, nil, err)
		return nil, malformedQuery(err.Error())
	}

	result, execErr := gremlinparser.ExecuteQuery(ctx, reader, writer, parsed, nil)
	s.resolveQuery(store, qid, result, execErr)
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
// executing it. The output field (Smithy httpPayload) contains the formatted
// explain plan text.
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

	plan, err := explainGremlinQuery(params.Gremlin)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	// M12: Smithy output field is "output" with httpPayload trait.
	// Return the formatted explain plan as a text block matching
	// AWS Neptune's explain response format.
	outputText := formatExplainOutput(plan)
	return map[string]interface{}{
		"output": outputText,
	}, nil
}

// ExecuteGremlinProfileQuery returns a profiled explain plan for a Gremlin query.
// Supports the profile.* parameters from the Smithy jsonName traits:
// profile.results, profile.chop, profile.serializer, profile.indexOps.
func (s *NeptuneDataService) ExecuteGremlinProfileQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	_ = reqCtx
	body := req.Body
	var params struct {
		Gremlin           string `json:"gremlin"`
		ProfileResults    *bool  `json:"profile.results"`
		ProfileChop       *int   `json:"profile.chop"`
		ProfileSerializer string `json:"profile.serializer"`
		ProfileIndexOps   *bool  `json:"profile.indexOps"`
	}
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, badRequest(fmt.Sprintf("invalid request body: %v", err))
	}
	if params.Gremlin == "" {
		return nil, missingParameter("gremlin")
	}

	profOpts := profileOptions{
		results:    params.ProfileResults,
		chop:       params.ProfileChop,
		serializer: params.ProfileSerializer,
		indexOps:   params.ProfileIndexOps,
	}
	plan, err := profileGremlinQueryEx(params.Gremlin, profOpts)
	if err != nil {
		return nil, malformedQuery(err.Error())
	}

	// M12/M13: Smithy output field is "output" with httpPayload trait.
	outputText := formatProfileOutput(plan, profOpts)
	return map[string]interface{}{
		"output": outputText,
	}, nil
}

// GetGremlinQueryStatus returns the current status and evaluation statistics of
// a previously submitted Gremlin query.
func (s *NeptuneDataService) GetGremlinQueryStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.getQueryStatus(reqCtx, req)
}

// ListGremlinQueries returns all submitted Gremlin queries, optionally
// including those in a waiting state.
func (s *NeptuneDataService) ListGremlinQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.listQueries(reqCtx, req, "gremlin")
}

// CancelGremlinQuery cancels a running Gremlin query and marks its status as
// cancelled.
func (s *NeptuneDataService) CancelGremlinQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_ = ctx
	return s.cancelQuery(reqCtx, req, false, false)
}
