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
	"fmt"
)

// Transform processes a VTL template string and returns the transformed result.
// This is the main entry point for template processing. The template is processed
// through multiple phases: control flow directives, input variables, utility
// functions, context variables, and stage variables.
//
// Returns an empty string if the input template is empty, or an error if
// template processing fails.
func (e *Engine) Transform(templateStr string) (string, error) {
	if templateStr == "" {
		return "", nil
	}

	e.transformErr = nil
	result := e.processAllPhases(templateStr)
	if e.transformErr != nil {
		return result, e.transformErr
	}
	return result, nil
}

// marshalJSON is a DRY helper that serialises a value to JSON.  When
// serialisation fails the first error is captured in e.transformErr so that
// Transform() can return it to the caller.  Using this helper everywhere
// eliminates the pervasive "jsonBytes, _ := json.Marshal(...)" silent-error
// anti-pattern across all VTL phases.
func (e *Engine) marshalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		if e.transformErr == nil {
			e.transformErr = fmt.Errorf("vtl: json marshal failed: %w", err)
		}
		return ""
	}
	return string(b)
}

// processAllPhases runs the full 8-phase VTL pipeline on the given template
// string. This is the single source of truth for the phase ordering and is
// used by both Transform and nested block processing (#foreach, #if) to
// ensure that variable substitution works consistently at every nesting
// level.
func (e *Engine) processAllPhases(result string) string {
	result = e.processControlFlow(result)
	result = e.processAppSyncUtil(result)
	result = e.processInput(result)
	result = e.processUtil(result)
	result = e.processAppSyncContext(result)
	result = e.processUtilToJsonFinal(result)
	result = e.processContext(result)
	result = e.processStageVariables(result)
	return result
}

// TransformRequest processes a request template using the provided body and
// content type. This is typically used for API Gateway request templates to
// transform incoming request data before it reaches the backend.
//
// The method selects the appropriate template based on the content type,
// builds the input context from the body, and then processes the template.
// Returns the transformed body as a byte slice.
func (e *Engine) TransformRequest(requestTemplates map[string]string, contentType string, body []byte) ([]byte, error) {
	if requestTemplates == nil {
		return body, nil
	}

	templateStr, ok := requestTemplates[contentType]
	if !ok {
		templateStr, ok = requestTemplates["*/*"]
		if !ok {
			return body, nil
		}
	}

	e.buildInputContext(body)

	result, err := e.Transform(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to transform request: %w", err)
	}

	return []byte(result), nil
}

// TransformResponse processes a response template using the provided body and
// content type. This is typically used for API Gateway response templates to
// transform backend responses before they are returned to the client.
//
// The method selects the appropriate template based on the content type,
// builds the input context from the body, and then processes the template.
// Returns the transformed body as a byte slice.
func (e *Engine) TransformResponse(responseTemplates map[string]string, contentType string, body []byte) ([]byte, error) {
	if responseTemplates == nil {
		return body, nil
	}

	templateStr, ok := responseTemplates[contentType]
	if !ok {
		templateStr, ok = responseTemplates["*/*"]
		if !ok {
			return body, nil
		}
	}

	e.buildInputContext(body)

	result, err := e.Transform(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to transform response: %w", err)
	}

	return []byte(result), nil
}
