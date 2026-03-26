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
	"regexp"
)

// Regular expressions for matching $stageVariables references in templates.
// Supports both dot notation and bracket notation for variable access.
var (
	stageVariablesDotRegex     = regexp.MustCompile(`\$stageVariables\.(\w+)`)
	stageVariablesBracketRegex = regexp.MustCompile(`\$\{?stageVariables\['([^']+)'\]\}?`)
)

// processStageVariables is the main entry point for processing all
// $stageVariables references in a template. It handles both dot and
// bracket notation for accessing stage variables.
func (e *Engine) processStageVariables(templateStr string) string {
	result := templateStr
	result = e.processStageVariablesDot(result)
	result = e.processStageVariablesBracket(result)
	return result
}

// processStageVariablesDot handles the $stageVariables.varName syntax
// for accessing stage variables.
func (e *Engine) processStageVariablesDot(templateStr string) string {
	return stageVariablesDotRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		varName := stageVariablesDotRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.StageVars[varName]; ok {
			return val
		}
		return ""
	})
}

// processStageVariablesBracket handles the $stageVariables['name'] and
// ${stageVariables['name']} syntaxes for accessing stage variables.
func (e *Engine) processStageVariablesBracket(templateStr string) string {
	return stageVariablesBracketRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		varName := stageVariablesBracketRegex.FindStringSubmatch(match)[1]
		if val, ok := e.context.StageVars[varName]; ok {
			return val
		}
		return ""
	})
}
