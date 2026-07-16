package dispatcher

import (
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/api"
)

func (d *Dispatcher) generateFallbackResponse(operation *api.Operation) map[string]interface{} {
	if operation == nil || operation.OutputShapeID == nil {
		return map[string]interface{}{}
	}
	response, err := d.mockGenerator.GenerateResponse(operation.OutputShapeID)
	if err != nil {
		return map[string]interface{}{"fallback": true, "operation": operation.Name}
	}
	return response
}

func (d *Dispatcher) detectOperationFromParsed(parsedReq *request.ParsedRequest) *api.Operation {
	if parsedReq == nil || parsedReq.Operation == "" {
		return nil
	}

	op, err := d.operationStore.GetByShapeID(parsedReq.Operation)
	if err == nil && op != nil {
		return op
	}

	var found *api.Operation
	_ = d.operationStore.ForEach(func(op *api.Operation) error {
		if op.Name == parsedReq.Operation {
			found = op
			return errors.New("found")
		}
		return nil
	})
	return found
}
