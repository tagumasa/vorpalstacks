package integration

import (
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/pkg/vtl"
)

// newIntegrationVTLEngine builds a VTL engine carrying the request context
// integration templates expect: the body (raw and, when parseable, as JSON),
// the merged path/query/header parameters, the stage/api/requestId context
// map, and the stage variables. Callers with extra context members (the mock
// executor adds path and method) set them on the returned engine.
func newIntegrationVTLEngine(req *IntegrationRequest, body []byte) *vtl.Engine {
	engine := vtl.NewEngine()

	engine.GetContext().Body = string(body)
	if len(body) > 0 {
		var bodyObj interface{}
		if err := json.Unmarshal(body, &bodyObj); err == nil {
			engine.GetContext().JSONBody = bodyObj
		}
	}

	params := make(map[string]string)
	for k, v := range req.PathParams {
		params[k] = v
	}
	for k, v := range req.QueryParams {
		params[k] = v
	}
	for k, v := range req.Headers {
		params[k] = v
	}
	engine.GetContext().Params = params

	engine.GetContext().Context = map[string]interface{}{
		"stage":     req.StageName,
		"apiId":     req.RestApiId,
		"requestId": fmt.Sprintf("%x", time.Now().UnixNano()),
	}

	engine.GetContext().StageVars = req.StageVariables

	return engine
}
