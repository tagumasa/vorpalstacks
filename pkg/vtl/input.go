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
	"encoding/json"
	"regexp"
	"strings"
)

// Regular expressions for matching $input variable references in templates.
// These patterns identify various ways to access request input data.
var (
	inputPathRegex      = regexp.MustCompile(`\$input\.path\(['"]([^'"]+)['"]\)`)
	inputJsonRegex      = regexp.MustCompile(`\$input\.json\(['"]([^'"]*)['"]\)`)
	inputBodyRegex      = regexp.MustCompile(`\$input\.body`)
	inputParamsRegex    = regexp.MustCompile(`\$input\.params\(['"]([^'"]+)['"]\)`)
	inputParamsAllRegex = regexp.MustCompile(`\$input\.params\(\)`)
)

// processInput is the main entry point for processing all $input variables
// in a template. It processes body, path, JSON, and parameter references.
func (e *Engine) processInput(templateStr string) string {
	result := templateStr
	result = e.processInputBody(result)
	result = e.processInputPath(result)
	result = e.processInputJson(result)
	result = e.processInputParams(result)
	result = e.processInputParamsAll(result)
	return result
}

// processInputBody handles the $input.body reference which returns the
// raw request body as a string.
func (e *Engine) processInputBody(templateStr string) string {
	return inputBodyRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return e.context.Body
	})
}

// processInputPath handles the $input.path() function which extracts a value
// from the JSON body using a JSON path expression and returns it as a string.
func (e *Engine) processInputPath(templateStr string) string {
	return inputPathRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		path := inputPathRegex.FindStringSubmatch(match)[1]
		if e.context.JSONBody == nil {
			return ""
		}
		return ExtractJSONPath(e.context.JSONBody, path)
	})
}

// processInputJson handles the $input.json() function which parses the JSON
// body and returns either the entire body or a specific path as a JSON string.
func (e *Engine) processInputJson(templateStr string) string {
	return inputJsonRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		path := inputJsonRegex.FindStringSubmatch(match)[1]
		if e.context.JSONBody == nil {
			return ""
		}
		if path == "" || path == "$" {
			jsonBytes, _ := json.Marshal(e.context.JSONBody)
			return string(jsonBytes)
		}
		value := ExtractJSONPathRaw(e.context.JSONBody, path)
		if value == nil {
			return ""
		}
		jsonBytes, _ := json.Marshal(value)
		return string(jsonBytes)
	})
}

// processInputParams handles the $input.params() function which retrieves
// a specific request parameter (query string, header, or path parameter).
func (e *Engine) processInputParams(templateStr string) string {
	return inputParamsRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		paramName := inputParamsRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.Params[paramName]; ok {
			return val
		}
		return ""
	})
}

// processInputParamsAll handles the $input.params() function with no arguments,
// which returns all parameters as a JSON object.
func (e *Engine) processInputParamsAll(templateStr string) string {
	return inputParamsAllRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		if len(e.context.Params) == 0 {
			return "{}"
		}
		jsonBytes, _ := json.Marshal(e.context.Params)
		return string(jsonBytes)
	})
}

// buildInputContext parses the request body and stores both the raw body
// string and the parsed JSON object in the context for use by other
// $input functions.
func (e *Engine) buildInputContext(body []byte) {
	if len(body) > 0 {
		e.context.Body = string(body)
		var bodyObj interface{}
		if err := json.Unmarshal(body, &bodyObj); err == nil {
			e.context.JSONBody = bodyObj
		}
	}
}

// GetParamValue retrieves a parameter value by name from a map of parameters.
// The search is case-insensitive.
func GetParamValue(params map[string]string, name string) string {
	for key, val := range params {
		if strings.EqualFold(key, name) {
			return val
		}
	}
	return ""
}

// GetAllParams converts a map of string parameters to a map of interface{}
// values, which is the format expected by the VTL engine.
func GetAllParams(params map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range params {
		result[k] = v
	}
	return result
}
