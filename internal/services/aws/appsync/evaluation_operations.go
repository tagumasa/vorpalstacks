package appsync

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/google/uuid"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/pkg/vtl"
)

// StartDataSourceIntrospection initiates an introspection of a data source.
// In the local environment external RDS is unavailable, so the introspection
// is immediately marked as FAILED with a descriptive status detail.
// POST /v1/datasources/introspections
func (s *AppSyncService) StartDataSourceIntrospection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	introspectionId := uuid.New().String()
	return map[string]interface{}{
		"introspectionId":           introspectionId,
		"introspectionStatus":       "FAILED",
		"introspectionStatusDetail": "Data source introspection is not supported in the local environment",
	}, nil
}

// GetDataSourceIntrospection retrieves the result of a data source introspection.
// GET /v1/datasources/introspections/{introspectionId}
func (s *AppSyncService) GetDataSourceIntrospection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	introspectionId := request.GetStringParam(req.Parameters, "introspectionId")
	if introspectionId == "" {
		return nil, NewBadRequestException("introspectionId is required")
	}
	return map[string]interface{}{
		"introspectionId":           introspectionId,
		"introspectionStatus":       "FAILED",
		"introspectionStatusDetail": "Data source introspection is not supported in the local environment",
	}, nil
}

// ListTypesByAssociation returns types from the source API of a merged API association.
// GET /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}/types
func (s *AppSyncService) ListTypesByAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	associationId := request.GetStringParam(req.Parameters, "associationId")
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return mapStoreError(err)
	}

	format := request.GetStringParam(req.Parameters, "format")
	if format == "" {
		return nil, NewBadRequestException("format is required")
	}

	opts := parsePaginationOptions(req)
	types, nextToken, err := store.ListTypes(assoc.SourceApiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(types))
	for _, t := range types {
		items = append(items, typeToMap(t))
	}

	response := map[string]interface{}{
		"types": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

func (s *AppSyncService) EvaluateCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var body struct {
		Code     string          `json:"code"`
		Context  json.RawMessage `json:"context"`
		Function string          `json:"function"`
		Runtime  struct {
			Name           string `json:"runtimeName"`
			RuntimeVersion string `json:"runtimeVersion"`
		} `json:"runtime"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return map[string]interface{}{
			"evaluationResult": "",
			"error":            map[string]string{"message": "invalid request body"},
			"logs":             []interface{}{},
			"outErrors":        "",
			"stash":            "{}",
		}, nil
	}

	if body.Code == "" {
		return map[string]interface{}{
			"evaluationResult": "",
			"error":            map[string]string{"message": "code is required"},
			"logs":             []interface{}{},
			"outErrors":        "",
			"stash":            "{}",
		}, nil
	}

	ctxMap := make(map[string]interface{})
	if len(body.Context) > 0 {
		if err := json.Unmarshal(body.Context, &ctxMap); err != nil {
			return map[string]interface{}{
				"evaluationResult": "",
				"error":            map[string]string{"message": "invalid context JSON"},
				"logs":             []interface{}{},
				"outErrors":        "",
				"stash":            "{}",
			}, nil
		}
	}

	vm := goja.New()
	logs := []interface{}{}
	hasError := false
	errorResult := ""
	var stashVal map[string]interface{}

	if v, ok := ctxMap["stash"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			stashVal = m
		}
	}
	if stashVal == nil {
		stashVal = make(map[string]interface{})
	}

	jsStash := vm.NewObject()
	for k, v := range stashVal {
		jsStash.Set(k, v)
	}

	ctxObj := vm.NewObject()
	for k, v := range ctxMap {
		if k == "stash" {
			ctxObj.Set(k, jsStash)
			continue
		}
		ctxObj.Set(k, v)
	}

	quietMode := false
	utilObj := buildUtilObject(vm, &logs, &hasError, &errorResult, &quietMode)
	vm.Set("util", utilObj)
	vm.Set("context", ctxObj)
	vm.Set("console", map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			parts := make([]string, len(call.Arguments))
			for i, a := range call.Arguments {
				parts[i] = a.String()
			}
			logs = append(logs, strings.Join(parts, " "))
			return goja.Undefined()
		},
		"error": func(call goja.FunctionCall) goja.Value {
			parts := make([]string, len(call.Arguments))
			for i, a := range call.Arguments {
				parts[i] = a.String()
			}
			logs = append(logs, strings.Join(parts, " "))
			return goja.Undefined()
		},
	})

	_, err := vm.RunString(body.Code)
	if err != nil {
		errMsg := err.Error()
		if jsErr, ok := err.(*goja.Exception); ok {
			errMsg = jsErr.Error()
		}
		outErr := errMsg
		if quietMode {
			outErr = ""
		}
		return map[string]interface{}{
			"evaluationResult": "",
			"error":            map[string]string{"message": errMsg},
			"logs":             logs,
			"outErrors":        outErr,
			"stash":            extractStashFromVM(vm, ctxObj),
		}, nil
	}

	var evalResult interface{}

	fnName := body.Function
	if fnName == "" {
		fnName = "handler"
	}

	if handlerFn := vm.Get(fnName); handlerFn != nil && !goja.IsNull(handlerFn) && !goja.IsUndefined(handlerFn) {
		if fn, ok := goja.AssertFunction(handlerFn); ok {
			ret, err := fn(goja.Undefined(), ctxObj)
			if err != nil {
				errMsg := err.Error()
				if jsErr, ok := err.(*goja.Exception); ok {
					errMsg = jsErr.Error()
				}
				outErr := errMsg
				if quietMode {
					outErr = ""
				}
				return map[string]interface{}{
					"evaluationResult": "",
					"error":            map[string]string{"message": errMsg},
					"logs":             logs,
					"outErrors":        outErr,
					"stash":            extractStashFromVM(vm, ctxObj),
				}, nil
			}
			if ret != nil && !goja.IsUndefined(ret) && !goja.IsNull(ret) {
				evalResult = ret.Export()
			}
		}
	}

	if hasError {
		outErr := errorResult
		if quietMode {
			outErr = ""
		}
		return map[string]interface{}{
			"evaluationResult": "",
			"error":            map[string]string{"message": errorResult},
			"logs":             logs,
			"outErrors":        outErr,
			"stash":            extractStashFromVM(vm, ctxObj),
		}, nil
	}

	evalResultStr := ""
	if evalResult != nil {
		b, err := json.Marshal(evalResult)
		if err == nil {
			evalResultStr = string(b)
		} else {
			evalResultStr = fmt.Sprintf("%v", evalResult)
		}
	}

	return map[string]interface{}{
		"evaluationResult": evalResultStr,
		"error":            nil,
		"logs":             logs,
		"outErrors":        "",
		"stash":            extractStashFromVM(vm, ctxObj),
	}, nil
}

func buildUtilObject(vm *goja.Runtime, logs *[]interface{}, hasError *bool, errorResult *string, quietMode *bool) map[string]interface{} {
	return map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			for _, a := range call.Arguments {
				*logs = append(*logs, a.String())
			}
			return goja.Undefined()
		},
		"error": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				*errorResult = call.Arguments[0].String()
			} else {
				*errorResult = "Unknown error"
			}
			*hasError = true
			panic(vm.NewGoError(fmt.Errorf("%s", *errorResult)))
		},
		"appendError": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				*logs = append(*logs, call.Arguments[0].String())
			}
			return goja.Undefined()
		},
		"autoId": func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(fmt.Sprintf("%08x", time.Now().UnixNano()))
		},
		"autoUlid": func(call goja.FunctionCall) goja.Value {
			return vm.ToValue(generateULID())
		},
		"isNull": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(true)
			}
			return vm.ToValue(goja.IsNull(call.Arguments[0]) || goja.IsUndefined(call.Arguments[0]))
		},
		"isUndefined": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(true)
			}
			return vm.ToValue(goja.IsUndefined(call.Arguments[0]))
		},
		"isString": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(false)
			}
			return vm.ToValue(call.Arguments[0].ExportType().Kind() == reflect.String)
		},
		"isNumber": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(false)
			}
			return vm.ToValue(goja.IsNumber(call.Arguments[0]))
		},
		"isBoolean": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(false)
			}
			exported := call.Arguments[0].Export()
			_, ok := exported.(bool)
			return vm.ToValue(ok)
		},
		"isList": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(false)
			}
			return vm.ToValue(call.Arguments[0].ExportType().Kind() == reflect.Slice)
		},
		"isMap": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(false)
			}
			return vm.ToValue(call.Arguments[0].ExportType().Kind() == reflect.Map)
		},
		"isObject": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return vm.ToValue(false)
			}
			k := call.Arguments[0].ExportType().Kind()
			return vm.ToValue(k == reflect.Map || k == reflect.Slice)
		},
		"qr": func(call goja.FunctionCall) goja.Value {
			*quietMode = true
			return goja.Undefined()
		},
		"runtime": map[string]interface{}{
			"name":           "APPSYNC_JS",
			"runtimeVersion": "1.0.0",
		},
		"time": map[string]interface{}{
			"nowISO8601": func(call goja.FunctionCall) goja.Value {
				return vm.ToValue(time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
			},
			"nowEpochSeconds": func(call goja.FunctionCall) goja.Value {
				return vm.ToValue(time.Now().Unix())
			},
			"nowEpochMillis": func(call goja.FunctionCall) goja.Value {
				return vm.ToValue(time.Now().UnixMilli())
			},
		},
		"defaultIfNull": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return goja.Undefined()
			}
			if goja.IsNull(call.Arguments[0]) || goja.IsUndefined(call.Arguments[0]) {
				if len(call.Arguments) > 1 {
					return call.Arguments[1]
				}
				return goja.Null()
			}
			return call.Arguments[0]
		},
		"unauthorized": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				*errorResult = "Unauthorized: " + call.Arguments[0].String()
			} else {
				*errorResult = "Unauthorized"
			}
			*hasError = true
			panic(vm.NewGoError(fmt.Errorf("%s", *errorResult)))
		},
	}
}

func extractStashFromVM(vm *goja.Runtime, ctxObj *goja.Object) string {
	stashVal := ctxObj.Get("stash")
	if stashVal == nil || goja.IsUndefined(stashVal) || goja.IsNull(stashVal) {
		return "{}"
	}
	exported := stashVal.Export()
	if exported == nil {
		return "{}"
	}
	b, err := json.Marshal(exported)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// EvaluateMappingTemplate evaluates a VTL mapping template.
// POST /v1/dataplane-evaluatetemplate
func (s *AppSyncService) EvaluateMappingTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var body struct {
		Context  string `json:"context"`
		Template string `json:"template"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return map[string]interface{}{
			"evaluationResult": "",
			"error":            map[string]string{"message": "invalid request body"},
			"logs":             []interface{}{},
			"outErrors":        "",
			"stash":            "{}",
		}, nil
	}

	ctxMap := make(map[string]interface{})
	if body.Context != "" {
		if err := json.Unmarshal([]byte(body.Context), &ctxMap); err != nil {
			return map[string]interface{}{
				"evaluationResult": "",
				"error":            map[string]string{"message": "invalid context JSON"},
				"logs":             []interface{}{},
				"outErrors":        "",
				"stash":            "{}",
			}, nil
		}
	}

	stash := make(map[string]interface{})
	if v, ok := ctxMap["stash"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			stash = m
		}
	}

	var identity map[string]interface{}
	if v, ok := ctxMap["identity"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			identity = m
		}
	}

	var reqMap map[string]interface{}
	if v, ok := ctxMap["request"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			reqMap = m
		}
	}

	var info *vtl.AppSyncFieldInfo
	if v, ok := ctxMap["info"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			info = &vtl.AppSyncFieldInfo{}
			if fn, ok := m["fieldName"].(string); ok {
				info.FieldName = fn
			}
			if pt, ok := m["parentTypeName"].(string); ok {
				info.ParentTypeName = pt
			}
			if ss, ok := m["selectionSetGraphQL"].(string); ok {
				info.SelectionSetGraphQL = ss
			}
			if sl, ok := m["selectionSetList"].([]interface{}); ok {
				for _, s := range sl {
					if str, ok := s.(string); ok {
						info.SelectionSetList = append(info.SelectionSetList, str)
					}
				}
			}
		}
	}

	var resultVal interface{}
	if v, ok := ctxMap["result"]; ok {
		resultVal = v
	}

	engine := vtl.NewEngine()
	engine.AppSyncCtx = &vtl.AppSyncContext{
		Args:     extractMap(ctxMap, "args"),
		Source:   ctxMap["source"],
		Stash:    stash,
		Identity: identity,
		Info:     info,
		Result:   resultVal,
		Request:  reqMap,
		Errors:   []vtl.AppSyncError{},
	}

	evalResult, err := engine.Transform(body.Template)

	if err != nil {
		return map[string]interface{}{
			"evaluationResult": "",
			"error":            map[string]string{"message": err.Error()},
			"logs":             []interface{}{},
			"outErrors":        "",
			"stash":            stashJSON(engine.AppSyncCtx.Stash),
		}, nil
	}

	var errorObj interface{}
	outErrors := ""
	if engine.AppSyncCtx != nil && len(engine.AppSyncCtx.Errors) > 0 {
		errArr := make([]map[string]interface{}, 0, len(engine.AppSyncCtx.Errors))
		for _, ae := range engine.AppSyncCtx.Errors {
			eObj := map[string]interface{}{"message": ae.Message}
			if ae.ErrorType != "" {
				eObj["errorType"] = ae.ErrorType
			}
			if ae.Data != nil {
				eObj["data"] = ae.Data
			}
			errArr = append(errArr, eObj)
		}
		errorObj = errArr
		errBytes, _ := json.Marshal(engine.AppSyncCtx.Errors)
		outErrors = string(errBytes)
	}

	return map[string]interface{}{
		"evaluationResult": evalResult,
		"error":            errorObj,
		"logs":             []interface{}{},
		"outErrors":        outErrors,
		"stash":            stashJSON(engine.AppSyncCtx.Stash),
	}, nil
}

func extractMap(m map[string]interface{}, key string) map[string]interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}
	if mapped, ok := v.(map[string]interface{}); ok {
		return mapped
	}
	return nil
}

func stashJSON(stash map[string]interface{}) string {
	if stash == nil {
		return "{}"
	}
	b, err := json.Marshal(stash)
	if err != nil {
		return "{}"
	}
	return string(b)
}

const crockfordBase32 = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func generateULID() string {
	ms := time.Now().UnixMilli()
	b := make([]byte, 16)
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	_, _ = rand.Read(b[6:])
	var sb strings.Builder
	bits := uint(0)
	n := uint(0)
	for _, v := range b {
		bits = bits<<8 | uint(v)
		n += 8
		for n >= 5 {
			n -= 5
			sb.WriteByte(crockfordBase32[bits>>n&0x1f])
		}
	}
	if n > 0 {
		sb.WriteByte(crockfordBase32[bits<<(5-n)&0x1f])
	}
	return sb.String()
}
