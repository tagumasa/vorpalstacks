package vtl

import (
	"encoding/json"
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
			b, _ := json.Marshal(info.Variables)
			return string(b)
		case "selectionSetGraphQL":
			return info.SelectionSetGraphQL
		case "selectionSetList":
			if info.SelectionSetList == nil {
				return "[]"
			}
			b, _ := json.Marshal(info.SelectionSetList)
			return string(b)
		case "parentTypeFields":
			if info.ParentTypeFields == nil {
				return "[]"
			}
			b, _ := json.Marshal(info.ParentTypeFields)
			return string(b)
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
		b, _ := json.Marshal(e.AppSyncCtx.Args)
		return string(b)
	})
	result = ctxSourceAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Source == nil {
			return "null"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Source)
		return string(b)
	})
	result = ctxStashAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Stash == nil {
			return "{}"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Stash)
		return string(b)
	})
	result = ctxIdentityAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Identity == nil {
			return "{}"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Identity)
		return string(b)
	})
	result = ctxInfoAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Info == nil {
			return "{}"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Info)
		return string(b)
	})
	result = ctxResultAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Result == nil {
			return "null"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Result)
		return string(b)
	})
	result = ctxErrorAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Error == nil {
			return "null"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Error)
		return string(b)
	})
	result = ctxRequestAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Request == nil {
			return "{}"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Request)
		return string(b)
	})
	result = ctxPrevAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Prev == nil {
			return "{}"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Prev)
		return string(b)
	})
	result = ctxTriggerAllRegex.ReplaceAllStringFunc(result, func(match string) string {
		if e.AppSyncCtx.Trigger == nil {
			return "{}"
		}
		b, _ := json.Marshal(e.AppSyncCtx.Trigger)
		return string(b)
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
	return e.formatValue(current)
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
	return e.formatValue(current)
}
