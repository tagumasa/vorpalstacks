package dispatcher

import (
	"net/http"

	"vorpalstacks/internal/services/aws/common/request"
)

func (d *Dispatcher) buildParsedRequest(r *http.Request, dispatchCtx *DispatchContext) *request.ParsedRequest {
	parsedReq, _ := request.ParseAWSRequest(r)
	if parsedReq == nil {
		parsedReq = &request.ParsedRequest{
			Headers:     r.Header,
			QueryParams: r.URL.Query(),
			PathParams:  make(map[string]string),
			Parameters:  make(map[string]interface{}),
		}
	}
	if parsedReq.Parameters == nil {
		parsedReq.Parameters = make(map[string]interface{})
	}

	parsedReq.Operation = dispatchCtx.Operation

	if dispatchCtx.Params != nil {
		for k, v := range dispatchCtx.Params {
			parsedReq.Parameters[k] = v
		}
	}

	return parsedReq
}
