/*
 * Copyright 2026 Vorpalstacks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vtl

import (
	"fmt"
	"regexp"
	"strings"
)

// Regex patterns for matching context variables in templates.
// These patterns are used to identify and extract various types of
// context variables during template processing.
var (
	contextVariablesRegex        = regexp.MustCompile(`\$context\.([\w.-]+)`)
	contextAuthorizerRegex       = regexp.MustCompile(`\$context\.authorizer\.([\w.-]+)`)
	contextIdentityRegex         = regexp.MustCompile(`\$context\.identity\.([\w.-]+)`)
	contextRequestOverrideRegex  = regexp.MustCompile(`\$context\.requestOverride\.([\w.-]+)`)
	contextResponseOverrideRegex = regexp.MustCompile(`\$context\.responseOverride\.([\w.-]+)`)
	simpleVariableRegex          = regexp.MustCompile(`\$([\w]+(?:\.[\w]+)*)`)
)

// Reserved prefixes that should not be processed as simple variables.
// These are handled by their specific processing functions.
var reservedPrefixes = []string{"authorizer.", "identity.", "requestOverride.", "responseOverride.", "stageVariables."}

// processContext is the main entry point for processing all context-related
// variables in a template. It processes context variables in a specific order:
// authorizer, identity, request override, response override, generic context variables,
// and finally simple variables.
func (e *Engine) processContext(templateStr string) string {
	result := templateStr
	result = e.processContextAuthorizer(result)
	result = e.processContextIdentity(result)
	result = e.processContextRequestOverride(result)
	result = e.processContextResponseOverride(result)
	result = e.processContextVariables(result)
	result = e.processSimpleVariables(result)
	return result
}

// processSimpleVariables handles simple variable interpolation using the $var
// or $obj.prop syntax. It resolves variables from the context map, user-defined
// variables, and special loop metadata.
func (e *Engine) processSimpleVariables(templateStr string) string {
	return simpleVariableRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		path := simpleVariableRegex.FindStringSubmatch(match)[1]
		if strings.HasPrefix(path, "context.") {
			return match
		}
		for _, prefix := range reservedPrefixes {
			if strings.HasPrefix(path, prefix) {
				return match
			}
		}
		if e.isInsideDirective(templateStr, match) {
			return match
		}
		val := e.getVariableByPath(path)
		if val == nil {
			return match
		}
		return e.formatValue(val)
	})
}

// isInsideDirective determines whether a variable reference is located inside
// a VTL directive (such as #if, #foreach, or #set). This is used to prevent
// variable substitution within directive expressions.
func (e *Engine) isInsideDirective(template, match string) bool {
	idx := strings.Index(template, match)
	if idx == -1 {
		return false
	}
	before := template[:idx]
	inDirective := false
	directiveDepth := 0
	for i := 0; i < len(before); i++ {
		if i+1 < len(before) && before[i] == '#' && (before[i+1] == 'i' || before[i+1] == 'f' || before[i+1] == 'e' || before[i+1] == 's') {
			if i+3 < len(before) && before[i:i+4] == "#if " {
				inDirective = true
				directiveDepth++
			} else if i+5 < len(before) && before[i:i+6] == "#foreach" {
				inDirective = true
				directiveDepth++
			} else if i+4 < len(before) && before[i:i+5] == "#else" {
			} else if i+3 < len(before) && before[i:i+4] == "#set" {
				inDirective = true
				directiveDepth++
			}
		}
		if before[i] == '(' && inDirective {
			directiveDepth++
		}
		if before[i] == ')' && inDirective {
			directiveDepth--
			if directiveDepth == 0 {
				inDirective = false
			}
		}
	}
	return inDirective
}

// getVariableByPath resolves a dot-separated path to its corresponding value.
// It searches in the context map first, then in user-defined variables,
// and supports special handling for loop metadata.
func (e *Engine) getVariableByPath(path string) interface{} {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil
	}

	var current interface{}
	if val, ok := e.context.Context[parts[0]]; ok {
		current = val
	} else if val, ok := e.variables[parts[0]]; ok {
		current = val
	} else {
		return nil
	}

	for i := 1; i < len(parts); i++ {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, ok := v[parts[i]]; ok {
				current = val
			} else {
				return nil
			}
		case map[string]string:
			if val, ok := v[parts[i]]; ok {
				current = val
			} else {
				return nil
			}
		case ContextLoopMetadata:
			if parts[i] == "index" {
				return v.Index
			} else if parts[i] == "count" {
				return v.Count
			} else if parts[i] == "first" {
				return v.First
			} else if parts[i] == "last" {
				return v.Last
			} else if parts[i] == "hasNext" {
				return v.HasNext
			} else if parts[i] == "hasPrev" {
				return v.HasPrev
			}
			return nil
		default:
			return nil
		}
	}

	return current
}

// formatValue converts a Go value to its string representation suitable
// for template output. It handles various primitive types including strings,
// integers, floats, and booleans.
func (e *Engine) formatValue(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// processContextVariables handles generic $context.* variables that are
// not covered by more specific handlers (authorizer, identity, etc.).
func (e *Engine) processContextVariables(templateStr string) string {
	return contextVariablesRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		varName := contextVariablesRegex.FindStringSubmatch(match)[1]
		for _, prefix := range reservedPrefixes {
			if strings.HasPrefix(varName, prefix) {
				return ""
			}
		}
		if val, ok := e.context.Context[varName]; ok {
			return e.formatValue(val)
		}
		return ""
	})
}

// processContextAuthorizer handles $context.authorizer.* variables which
// provide access to authoriser context data from API Gateway Lambda authorisers.
func (e *Engine) processContextAuthorizer(templateStr string) string {
	return contextAuthorizerRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := contextAuthorizerRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.Authorizer[key]; ok {
			switch v := val.(type) {
			case string:
				return v
			default:
				return ""
			}
		}
		return ""
	})
}

// processContextIdentity handles $context.identity.* variables which provide
// access to caller identity information such as user ARN, account ID, etc.
func (e *Engine) processContextIdentity(templateStr string) string {
	return contextIdentityRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := contextIdentityRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.Identity[key]; ok {
			switch v := val.(type) {
			case string:
				return v
			default:
				return ""
			}
		}
		return ""
	})
}

// processContextRequestOverride handles $context.requestOverride.* variables
// which allow access to request parameters that can override defaults.
func (e *Engine) processContextRequestOverride(templateStr string) string {
	return contextRequestOverrideRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := contextRequestOverrideRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.RequestOverride[key]; ok {
			return val
		}
		return ""
	})
}

// processContextResponseOverride handles $context.responseOverride.* variables
// which allow modification of response headers and status codes.
func (e *Engine) processContextResponseOverride(templateStr string) string {
	return contextResponseOverrideRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		key := contextResponseOverrideRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.ResponseOverride[key]; ok {
			return val
		}
		return ""
	})
}
