package apigateway

import (
	"strconv"
)

// PatchOperation represents a JSON Patch operation for updating API Gateway resources.
type PatchOperation struct {
	Op    string
	Path  string
	Value string
}

func parsePatchOperations(params map[string]interface{}) []PatchOperation {
	var ops []PatchOperation
	patchOps, ok := params["patchOperations"].([]interface{})
	if !ok {
		return ops
	}

	for _, op := range patchOps {
		if opMap, ok := op.(map[string]interface{}); ok {
			po := PatchOperation{}
			if o, ok := opMap["op"].(string); ok {
				po.Op = o
			}
			if p, ok := opMap["path"].(string); ok {
				po.Path = p
			}
			if v, ok := opMap["value"].(string); ok {
				po.Value = v
			}
			ops = append(ops, po)
		}
	}
	return ops
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func parseInt32(s string) int32 {
	v, _ := strconv.ParseInt(s, 10, 32)
	return int32(v)
}

func parseFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
