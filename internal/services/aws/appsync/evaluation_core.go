package appsync

// Package appsync: evaluation_core.go holds the evaluation operations'
// execution path: request-contract validation, the JS (APPSYNC_JS) and VTL
// evaluation engines, and the AWS response/error semantics. The HTTP
// handlers only parse the wire body and delegate here.

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/google/uuid"

	"vorpalstacks/internal/utils/timeutils"
	"vorpalstacks/pkg/vtl"
)

// EvaluateCodeInput is the transport-agnostic input for EvaluateCode.
type EvaluateCodeInput struct {
	Code        string
	Context     json.RawMessage
	Function    string
	RuntimeName string
}

// EvaluateMappingTemplateInput is the transport-agnostic input for
// EvaluateMappingTemplate.
type EvaluateMappingTemplateInput struct {
	Context  string
	Template string
}

// evaluateCodeCore executes an AppSync function resolver code snippet and
// returns the result.
func evaluateCodeCore(in *EvaluateCodeInput) (map[string]interface{}, error) {
	// Required-member checks interleaved with the content validations
	// (code first, then context, then the runtime name) preserve the
	// original failure precedence.
	if in.Code == "" {
		return nil, NewBadRequestException("code is required")
	}
	if err := validateCode(in.Code); err != nil {
		return nil, err
	}
	if len(in.Context) == 0 {
		return nil, NewBadRequestException("context is required")
	}
	if err := validateContext(string(in.Context)); err != nil {
		return nil, err
	}
	if in.RuntimeName == "" {
		return nil, NewBadRequestException("runtime is required")
	}

	ctxMap := make(map[string]interface{})
	if len(in.Context) > 0 {
		if err := json.Unmarshal(in.Context, &ctxMap); err != nil {
			return nil, NewBadRequestException("invalid context JSON")
		}
	}

	vm := goja.New()
	// Enforce a maximum execution time to prevent infinite loops in
	// user-supplied JavaScript. The interrupt is delivered asynchronously
	// and causes the VM to return an error on the next instruction boundary.
	timer := time.AfterFunc(5*time.Second, func() {
		vm.Interrupt(fmt.Errorf("evaluation timed out after 5 seconds"))
	})
	defer timer.Stop()
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

	_, err := vm.RunString(in.Code)
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

	fnName := in.Function
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
			// AppSync documents util.autoId() as a 128-bit random UUID;
			// the VTL implementation already returns one, so the JS
			// runtime must match it rather than minting a clock value.
			return vm.ToValue(uuid.New().String())
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
				return vm.ToValue(time.Now().UTC().Format(timeutils.ISO8601UTCFormat))
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

// evaluateMappingTemplateCore evaluates a VTL mapping template.
func evaluateMappingTemplateCore(in *EvaluateMappingTemplateInput) (map[string]interface{}, error) {
	// Required-member checks interleaved with the content validations
	// (context first, then template) preserve the original failure
	// precedence.
	if in.Context == "" {
		return nil, NewBadRequestException("context is required")
	}
	if err := validateContext(in.Context); err != nil {
		return nil, err
	}
	if in.Template == "" {
		return nil, NewBadRequestException("template is required")
	}
	if err := validateTemplate(in.Template); err != nil {
		return nil, err
	}

	ctxMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(in.Context), &ctxMap); err != nil {
		return nil, NewBadRequestException("invalid context JSON")
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

	evalResult, err := engine.Transform(in.Template)

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
		firstErr := engine.AppSyncCtx.Errors[0]
		eObj := map[string]interface{}{"message": firstErr.Message}
		if firstErr.ErrorType != "" {
			eObj["errorType"] = firstErr.ErrorType
		}
		if firstErr.Data != nil {
			eObj["data"] = firstErr.Data
		}
		errorObj = eObj
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
