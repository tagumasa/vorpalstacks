package vtl

import (
	"reflect"
	"regexp"
	"strings"
)

var (
	ctxArgsRegex         = regexp.MustCompile(`\$ctx\.args\.([\w.-]+)`)
	ctxSourceRegex       = regexp.MustCompile(`\$ctx\.source\.([\w.-]+)`)
	ctxStashRegex        = regexp.MustCompile(`\$ctx\.stash\.([\w.-]+)`)
	ctxIdentityRegex     = regexp.MustCompile(`\$ctx\.identity\.([\w.-]+)`)
	ctxInfoRegex         = regexp.MustCompile(`\$ctx\.info\.([\w.-]+)`)
	ctxResultNestedRegex = regexp.MustCompile(`\$ctx\.result\.([\w.-]+)`)
	ctxErrorNestedRegex  = regexp.MustCompile(`\$ctx\.error\.([\w.-]+)`)
	ctxRequestRegex      = regexp.MustCompile(`\$ctx\.request\.([\w.-]+)`)
	ctxPrevRegex         = regexp.MustCompile(`\$ctx\.prev\.([\w.-]+)`)
	ctxTriggerRegex      = regexp.MustCompile(`\$ctx\.trigger\.([\w.-]+)`)
	ctxArgsAllRegex      = regexp.MustCompile(`\$ctx\.args\b`)
	ctxSourceAllRegex    = regexp.MustCompile(`\$ctx\.source\b`)
	ctxStashAllRegex     = regexp.MustCompile(`\$ctx\.stash\b`)
	ctxIdentityAllRegex  = regexp.MustCompile(`\$ctx\.identity\b`)
	ctxInfoAllRegex      = regexp.MustCompile(`\$ctx\.info\b`)
	ctxResultAllRegex    = regexp.MustCompile(`\$ctx\.result\b`)
	ctxErrorAllRegex     = regexp.MustCompile(`\$ctx\.error\b`)
	ctxRequestAllRegex   = regexp.MustCompile(`\$ctx\.request\b`)
	ctxPrevAllRegex      = regexp.MustCompile(`\$ctx\.prev\b`)
	ctxTriggerAllRegex   = regexp.MustCompile(`\$ctx\.trigger\b`)
)

func (e *Engine) processAppSyncContext(templateStr string) string {
	if e.AppSyncCtx == nil {
		return templateStr
	}

	result := templateStr
	result = e.processCtxArgsNested(result)
	result = e.processCtxSourceNested(result)
	result = e.processCtxStashNested(result)
	result = e.processCtxIdentityNested(result)
	result = e.processCtxInfoNested(result)
	result = e.processCtxRequestNested(result)
	result = e.processCtxPrevNested(result)
	result = e.processCtxTriggerNested(result)
	result = e.processCtxResultNested(result)
	result = e.processCtxErrorNested(result)
	result = e.processCtxTopLevelRefs(result)
	return result
}

func (e *Engine) processCtxArgsNested(templateStr string) string {
	return ctxArgsRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxArgsRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Args == nil {
			return "null"
		}
		return e.resolveFromMap(e.AppSyncCtx.Args, key)
	})
}

func (e *Engine) processCtxSourceNested(templateStr string) string {
	return ctxSourceRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxSourceRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Source == nil {
			return "null"
		}
		return e.resolveFromInterface(e.AppSyncCtx.Source, key)
	})
}

func (e *Engine) processCtxStashNested(templateStr string) string {
	return ctxStashRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxStashRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Stash == nil {
			return "null"
		}
		return e.resolveFromMap(e.AppSyncCtx.Stash, key)
	})
}

func (e *Engine) processCtxIdentityNested(templateStr string) string {
	return ctxIdentityRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxIdentityRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Identity == nil {
			return "null"
		}
		return e.resolveFromMap(e.AppSyncCtx.Identity, key)
	})
}

func (e *Engine) processCtxInfoNested(templateStr string) string {
	return ctxInfoRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxInfoRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Info == nil {
			return ""
		}
		info := e.AppSyncCtx.Info
		switch key {
		case "fieldName":
			return info.FieldName
		case "parentTypeName":
			return info.ParentTypeName
		case "variables":
			if info.Variables == nil {
				return "{}"
			}
			return e.marshalJSON(info.Variables)
		case "selectionSetGraphQL":
			return info.SelectionSetGraphQL
		case "selectionSetList":
			if info.SelectionSetList == nil {
				return "[]"
			}
			return e.marshalJSON(info.SelectionSetList)
		case "parentTypeFields":
			if info.ParentTypeFields == nil {
				return "[]"
			}
			return e.marshalJSON(info.ParentTypeFields)
		case "rootTypeName":
			return info.RootTypeName
		default:
			return ""
		}
	})
}

func (e *Engine) processCtxResultNested(templateStr string) string {
	return ctxResultNestedRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxResultNestedRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Result == nil {
			return "null"
		}
		return e.resolveFromInterface(e.AppSyncCtx.Result, key)
	})
}

func (e *Engine) processCtxErrorNested(templateStr string) string {
	return ctxErrorNestedRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxErrorNestedRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Error == nil {
			return "null"
		}
		return e.resolveFromInterface(e.AppSyncCtx.Error, key)
	})
}

func (e *Engine) processCtxRequestNested(templateStr string) string {
	return ctxRequestRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxRequestRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Request == nil {
			return "null"
		}
		return e.resolveFromMap(e.AppSyncCtx.Request, key)
	})
}

func (e *Engine) processCtxPrevNested(templateStr string) string {
	return ctxPrevRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxPrevRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Prev == nil {
			return "null"
		}
		return e.resolveFromMap(e.AppSyncCtx.Prev, key)
	})
}

func (e *Engine) processCtxTriggerNested(templateStr string) string {
	return ctxTriggerRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := ctxTriggerRegex.FindStringSubmatch(match)[1]
		if e.AppSyncCtx.Trigger == nil {
			return "null"
		}
		return e.resolveFromMap(e.AppSyncCtx.Trigger, key)
	})
}

func (e *Engine) processCtxTopLevelRefs(templateStr string) string {
	result := ctxArgsAllRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		if e.AppSyncCtx.Args == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Args)
	})
	result = ctxSourceAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Source == nil {
			return "null"
		}
		return e.marshalJSON(e.AppSyncCtx.Source)
	})
	result = ctxStashAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Stash == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Stash)
	})
	result = ctxIdentityAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Identity == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Identity)
	})
	result = ctxInfoAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Info == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Info)
	})
	result = ctxResultAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Result == nil {
			return "null"
		}
		return e.marshalJSON(e.AppSyncCtx.Result)
	})
	result = ctxErrorAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Error == nil {
			return "null"
		}
		return e.marshalJSON(e.AppSyncCtx.Error)
	})
	result = ctxRequestAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Request == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Request)
	})
	result = ctxPrevAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Prev == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Prev)
	})
	result = ctxTriggerAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Trigger == nil {
			return "{}"
		}
		return e.marshalJSON(e.AppSyncCtx.Trigger)
	})
	return result
}

func (e *Engine) resolveFromMap(m map[string]interface{}, key string) string {
	parts := strings.Split(key, ".")
	var current interface{} = m
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		case map[string]string:
			current = v[part]
		default:
			return "null"
		}
		if current == nil {
			return "null"
		}
	}
	return e.formatAppSyncValue(current)
}

func (e *Engine) resolveFromInterface(obj interface{}, key string) string {
	parts := strings.Split(key, ".")
	var current interface{} = obj
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		case map[string]string:
			current = v[part]
		default:
			return "null"
		}
		if current == nil {
			return "null"
		}
	}
	return e.formatAppSyncValue(current)
}

func (e *Engine) formatAppSyncValue(val interface{}) string {
	if val == nil {
		return "null"
	}
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.String:
		s := e.marshalJSON(val)
		if len(s) >= 2 {
			return s[1 : len(s)-1]
		}
		return s
	case reflect.Map, reflect.Slice, reflect.Array:
		return e.marshalJSON(val)
	default:
		return e.formatValue(val)
	}
}
