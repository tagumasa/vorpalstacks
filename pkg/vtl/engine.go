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

	result := templateStr
	result = e.processControlFlow(result)
	result = e.processAppSyncUtil(result)
	result = e.processInput(result)
	result = e.processUtil(result)
	result = e.processAppSyncContext(result)
	result = e.processContext(result)
	result = e.processStageVariables(result)

	return result, nil
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
