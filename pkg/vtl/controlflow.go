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
	"regexp"
	"strings"
)

// Regular expressions for matching VTL control flow directives.
// These patterns are used to identify and parse various control structures.
var (
	ifStartRe      = regexp.MustCompile(`#if\s*\(`)
	elseifRe       = regexp.MustCompile(`#elseif\s*\(`)
	elseRe         = regexp.MustCompile(`#else`)
	endifRe        = regexp.MustCompile(`#end\s*`)
	foreachStartRe = regexp.MustCompile(`#foreach\s*\(\s*\$(\w+)\s+in\s*`)
	endforeachRe   = regexp.MustCompile(`#end\s*`)
	setRe          = regexp.MustCompile(`#set\s*\(\s*\$(\w+(?:\.\w+)*)\s*=\s*`)
	commentRe      = regexp.MustCompile(`#\*.*?\*#`)
)

// function is inlined into findMatchingEndImpl (loop body yields wrong results).
//
//go:noinline — required: Go 1.25 inliner produces incorrect code when this
func findBalancedCloseParen(s string, startIdx int) (closeParenIdx int, content string) {
	depth := 1
	inString := false
	stringChar := byte(0)
	for i := startIdx; i < len(s); i++ {
		ch := s[i]
		if inString {
			if ch == stringChar && (i == 0 || s[i-1] != '\\') {
				inString = false
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			inString = true
			stringChar = ch
		} else if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
			if depth == 0 {
				return i, strings.TrimSpace(s[startIdx:i])
			}
		}
	}
	return -1, ""
}

// processControlFlow is the main entry point for processing all control flow
// directives in a template. It handles comment removal first, then iteratively
// processes blocks until no more directives remain.
func (e *Engine) processControlFlow(template string) string {
	template = commentRe.ReplaceAllString(template, "")

	for {
		newTemplate, processed := e.processNextBlock(template)
		if !processed {
			break
		}
		template = newTemplate
	}

	return template
}

// processNextBlock attempts to process the next control flow block in the
// template. It finds whichever directive (#foreach, #if, #set) appears
// first in the template and processes it. Returns the modified template
// and a boolean indicating whether a block was processed.
func (e *Engine) processNextBlock(template string) (string, bool) {
	foreachMatch := foreachStartRe.FindStringIndex(template)
	ifMatch := ifStartRe.FindStringIndex(template)
	setMatch := setRe.FindStringIndex(template)

	bestPos := len(template) + 1
	bestType := ""

	if foreachMatch != nil && foreachMatch[0] < bestPos {
		bestPos = foreachMatch[0]
		bestType = "foreach"
	}
	if ifMatch != nil && ifMatch[0] < bestPos {
		bestPos = ifMatch[0]
		bestType = "if"
	}
	if setMatch != nil && setMatch[0] < bestPos {
		bestPos = setMatch[0]
		bestType = "set"
	}

	switch bestType {
	case "foreach":
		result, found := e.processForeachBlock(template)
		return result, found
	case "if":
		result, found := e.processIfBlock(template)
		return result, found
	case "set":
		result, found := e.processSetBlock(template)
		return result, found
	}
	return template, false
}

// processForeachBlock handles the #foreach directive which iterates over
// collections or arrays. It resolves the iterable expression, iterates
// through each item, and processes the body template for each iteration.
// Loop metadata such as index, first, last is made available via $foreach.
func (e *Engine) processForeachBlock(template string) (string, bool) {
	matches := foreachStartRe.FindStringSubmatchIndex(template)
	if matches == nil {
		return template, false
	}

	varName := template[matches[2]:matches[3]]
	closeIdx, iterExpr := findBalancedCloseParen(template, matches[1])
	if closeIdx == -1 {
		return template, false
	}

	blockStart := matches[0]
	tagEnd := closeIdx + 1

	endIdx := findMatchingEnd(template, tagEnd, "foreach")
	if endIdx == -1 {
		return template, false
	}

	body := template[tagEnd:endIdx]
	endTag := endforeachRe.FindString(template[endIdx:])
	blockEnd := endIdx + len(endTag)

	iterable := e.resolveIterable(iterExpr)
	if iterable == nil {
		iterable = []interface{}{}
	}

	items, ok := iterable.([]interface{})
	if !ok {
		if strSlice, ok := iterable.([]string); ok {
			items = make([]interface{}, len(strSlice))
			for i, s := range strSlice {
				items[i] = s
			}
		} else {
			items = []interface{}{}
		}
	}

	var result strings.Builder
	for i, item := range items {
		oldVal, hadVar := e.context.Context[varName]
		e.context.Context[varName] = item

		loopVars := map[string]interface{}{
			"index":     i + 1,
			"index0":    i,
			"count":     i + 1,
			"first":     i == 0,
			"last":      i == len(items)-1,
			"hasNext":   i < len(items)-1,
			"hasPrev":   i > 0,
			"length":    len(items),
			"revindex":  len(items) - i,
			"revindex0": len(items) - i - 1,
		}
		oldLoop, hadLoop := e.context.Context["foreach"]
		e.context.Context["foreach"] = loopVars

		rendered := e.processControlFlow(body)
		rendered = e.processAppSyncUtil(rendered)
		rendered = e.processAppSyncContext(rendered)
		rendered = e.processSimpleVariables(rendered)
		result.WriteString(rendered)

		if hadLoop {
			e.context.Context["foreach"] = oldLoop
		} else {
			delete(e.context.Context, "foreach")
		}

		if hadVar {
			e.context.Context[varName] = oldVal
		} else {
			delete(e.context.Context, varName)
		}
	}

	return template[:blockStart] + result.String() + template[blockEnd:], true
}

// processIfBlock handles the #if, #elseif, and #else directives.
// It evaluates conditions and processes the body of the first branch
// whose condition evaluates to true.
func (e *Engine) processIfBlock(template string) (string, bool) {
	matches := ifStartRe.FindStringIndex(template)
	if matches == nil {
		return template, false
	}

	closeIdx, condition := findBalancedCloseParen(template, matches[1])
	if closeIdx == -1 {
		return template, false
	}

	condition = strings.TrimSpace(condition)
	blockStart := matches[0]
	tagEnd := closeIdx + 1

	endIdx := findMatchingEnd(template, tagEnd, "if")
	if endIdx == -1 {
		return template, false
	}

	body := template[tagEnd:endIdx]
	endTag := endifRe.FindString(template[endIdx:])
	blockEnd := endIdx + len(endTag)

	branches := e.parseIfBranches(body, condition)

	for _, branch := range branches {
		if e.evaluateCondition(branch.condition) {
			rendered := e.processControlFlow(branch.body)
			rendered = e.processAppSyncUtil(rendered)
			rendered = e.processAppSyncContext(rendered)
			return template[:blockStart] + rendered + template[blockEnd:], true
		}
	}

	return template[:blockStart] + template[blockEnd:], true
}

// ifBranch represents a single branch in an #if/#elseif/#else structure.
type ifBranch struct {
	condition string
	body      string
}

// parseIfBranches parses the body of an #if block and extracts all branches
// including the initial #if, any #elseif branches, and an optional #else branch.
func (e *Engine) parseIfBranches(body, firstCondition string) []ifBranch {
	var branches []ifBranch
	remaining := body

	branches = append(branches, ifBranch{condition: firstCondition, body: ""})

	for {
		elseifPos, elseifCond := findElseif(remaining)
		elsePos := findElse(remaining)

		if elseifPos >= 0 && (elsePos < 0 || elseifPos < elsePos) {
			branches[len(branches)-1].body = remaining[:elseifPos]
			elseifTag := elseifRe.FindString(remaining[elseifPos:])
			closeIdx, _ := findBalancedCloseParen(remaining, elseifPos+len(elseifTag))
			if closeIdx != -1 {
				remaining = remaining[closeIdx+1:]
			} else {
				remaining = remaining[elseifPos+len(elseifTag):]
			}
			branches = append(branches, ifBranch{condition: elseifCond, body: ""})
		} else if elsePos >= 0 {
			branches[len(branches)-1].body = remaining[:elsePos]
			elseTag := elseRe.FindString(remaining[elsePos:])
			remaining = remaining[elsePos+len(elseTag):]
			branches = append(branches, ifBranch{condition: "true", body: remaining})
			break
		} else {
			branches[len(branches)-1].body = remaining
			break
		}
	}

	return branches
}

// findElseif searches for the next #elseif directive at the current nesting
// level. Returns the position and condition string of the #elseif, or -1 and
// empty string if not found.
func findElseif(body string) (int, string) {
	depth := 0
	idx := 0
	for idx < len(body) {
		ifMatch := ifStartRe.FindStringIndex(body[idx:])
		elseifMatch := elseifRe.FindStringIndex(body[idx:])
		elseMatch := elseRe.FindStringIndex(body[idx:])
		endifMatch := endifRe.FindStringIndex(body[idx:])

		var nextPos int = len(body) + 1
		var nextType string = ""

		if endifMatch != nil && endifMatch[0] < nextPos {
			nextPos = endifMatch[0]
			nextType = "endif"
		}
		if ifMatch != nil && ifMatch[0] < nextPos {
			nextPos = ifMatch[0]
			nextType = "if"
		}
		if elseifMatch != nil && elseifMatch[0] < nextPos {
			nextPos = elseifMatch[0]
			nextType = "elseif"
		}
		if elseMatch != nil && elseMatch[0] < nextPos {
			nextPos = elseMatch[0]
			nextType = "else"
		}

		if nextPos > len(body) {
			break
		}

		switch nextType {
		case "if":
			depth++
			closeIdx, _ := findBalancedCloseParen(body, idx+ifMatch[1])
			if closeIdx != -1 {
				idx += closeIdx + 1
			} else {
				idx += ifMatch[1]
			}
		case "endif":
			if depth == 0 {
				return -1, ""
			}
			depth--
			idx += endifMatch[1]
		case "elseif":
			contentStart := idx + elseifMatch[1]
			closeIdx, cond := findBalancedCloseParen(body, contentStart)
			if closeIdx == -1 {
				idx += elseifMatch[1]
				continue
			}
			if depth == 0 {
				return idx + nextPos, cond
			}
			idx += closeIdx + 1
		case "else":
			if depth == 0 {
				return idx + nextPos, ""
			}
			idx += elseMatch[1]
		default:
			idx++
		}
	}
	return -1, ""
}

// findElse searches for the next #else directive at the current nesting level.
// Returns the position of the #else or -1 if not found.
func findElse(body string) int {
	depth := 0
	idx := 0
	for idx < len(body) {
		ifMatch := ifStartRe.FindStringIndex(body[idx:])
		elseifMatch := elseifRe.FindStringIndex(body[idx:])
		elseMatch := elseRe.FindStringIndex(body[idx:])
		endifMatch := endifRe.FindStringIndex(body[idx:])

		var nextPos int = len(body) + 1
		var nextType string = ""

		if endifMatch != nil && endifMatch[0] < nextPos {
			nextPos = endifMatch[0]
			nextType = "endif"
		}
		if ifMatch != nil && ifMatch[0] < nextPos {
			nextPos = ifMatch[0]
			nextType = "if"
		}
		if elseifMatch != nil && elseifMatch[0] < nextPos {
			nextPos = elseifMatch[0]
			nextType = "elseif"
		}
		if elseMatch != nil && elseMatch[0] < nextPos {
			nextPos = elseMatch[0]
			nextType = "else"
		}

		if nextPos > len(body) {
			break
		}

		switch nextType {
		case "if":
			depth++
			closeIdx, _ := findBalancedCloseParen(body, idx+ifMatch[1])
			if closeIdx != -1 {
				idx += closeIdx + 1
			} else {
				idx += ifMatch[1]
			}
		case "endif":
			if depth == 0 {
				return -1
			}
			depth--
			idx += endifMatch[1]
		case "elseif":
			closeIdx, _ := findBalancedCloseParen(body, idx+elseifMatch[1])
			if closeIdx != -1 {
				idx += closeIdx + 1
			} else {
				idx += elseifMatch[1]
			}
		case "else":
			if depth == 0 {
				return idx + nextPos
			}
			idx += elseMatch[1]
		default:
			idx++
		}
	}
	return -1
}

// processSetBlock handles the #set directive which assigns values to variables.
// The value can be a string literal, boolean, null, or a variable reference.
func (e *Engine) processSetBlock(template string) (string, bool) {
	matches := setRe.FindStringSubmatchIndex(template)
	if matches == nil {
		return template, false
	}

	varName := template[matches[2]:matches[3]]
	closeIdx, valueExpr := findBalancedCloseParen(template, matches[1])
	if closeIdx == -1 {
		return template, false
	}

	value := e.resolveValue(valueExpr)

	if strings.Contains(varName, ".") {
		parts := strings.SplitN(varName, ".", 2)
		if parts[0] == "ctx" && e.AppSyncCtx != nil {
			e.setAppSyncCtxPath(parts[1], value)
		} else if parts[0] == "context" {
			e.setContextPath(parts[1], value)
		} else {
			e.setNestedVariable(varName, value)
		}
	} else {
		if strings.Contains(valueExpr, ".add(") {
			e.processMethodCall(varName, valueExpr)
		} else {
			e.context.Context[varName] = value
		}
	}

	return template[:matches[0]] + template[closeIdx+1:], true
}

func (e *Engine) setAppSyncCtxPath(path string, value interface{}) {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) < 2 {
		return
	}
	topLevel := parts[0]
	field := parts[1]

	switch topLevel {
	case "stash":
		if e.AppSyncCtx.Stash == nil {
			e.AppSyncCtx.Stash = map[string]interface{}{}
		}
		e.AppSyncCtx.Stash[field] = value
	case "args":
		if e.AppSyncCtx.Args == nil {
			e.AppSyncCtx.Args = map[string]interface{}{}
		}
		e.AppSyncCtx.Args[field] = value
	}
}

func (e *Engine) setContextPath(path string, value interface{}) {
	e.context.Context[path] = value
}

func (e *Engine) setNestedVariable(varName string, value interface{}) {
	parts := strings.SplitN(varName, ".", 2)
	if len(parts) < 2 {
		e.context.Context[varName] = value
		return
	}
	base, ok := e.context.Context[parts[0]]
	if !ok {
		base = map[string]interface{}{}
	}
	if m, ok := base.(map[string]interface{}); ok {
		m[parts[1]] = value
		e.context.Context[parts[0]] = m
	}
}

var methodCallRe = regexp.MustCompile(`\$([\w]+)\.add\(([^)]+)\)`)

func (e *Engine) processMethodCall(varName string, valueExpr string) {
	matches := methodCallRe.FindStringSubmatch(valueExpr)
	if matches == nil {
		e.context.Context[varName] = e.resolveValue(valueExpr)
		return
	}
	targetVar := matches[1]
	argExpr := strings.TrimSpace(matches[2])

	current, ok := e.context.Context[targetVar]
	if !ok {
		return
	}

	arr, ok := current.([]interface{})
	if !ok {
		return
	}

	arg := e.resolveValue(argExpr)
	e.context.Context[targetVar] = append(arr, arg)
	e.context.Context[varName] = true
}

// findMatchingEnd finds the matching #end tag for a block starting at the
// given index. It properly handles nested control structures by tracking
// the nesting depth. Returns the position of the matching #end or -1 if
// no matching end is found.
// function (identical logic inlined at call sites returns -1; direct call works).
//
//go:noinline — required: Go 1.25 inliner produces incorrect code for this
func findMatchingEndImpl(template string, startIdx int) int {
	depth := 1
	idx := startIdx
	for idx < len(template) && depth > 0 {
		ifMatch := ifStartRe.FindStringIndex(template[idx:])
		foreachMatch := foreachStartRe.FindStringIndex(template[idx:])
		endMatch := endifRe.FindStringIndex(template[idx:])

		var nextIf, nextForeach, nextEnd int
		if ifMatch != nil {
			nextIf = ifMatch[0]
		} else {
			nextIf = len(template[idx:]) + 1
		}
		if foreachMatch != nil {
			nextForeach = foreachMatch[0]
		} else {
			nextForeach = len(template[idx:]) + 1
		}
		if endMatch != nil {
			nextEnd = endMatch[0]
		} else {
			return -1
		}

		nextStart := nextIf
		if nextForeach < nextStart {
			nextStart = nextForeach
		}

		if nextEnd < nextStart {
			depth--
			if depth == 0 {
				return idx + nextEnd
			}
			idx += endMatch[1]
		} else {
			depth++
			if nextIf < nextForeach {
				closeIdx, _ := findBalancedCloseParen(template, idx+ifMatch[1])
				if closeIdx != -1 {
					idx = closeIdx + 1
				} else {
					idx += ifMatch[1]
				}
			} else {
				closeIdx, _ := findBalancedCloseParen(template, idx+foreachMatch[1])
				if closeIdx != -1 {
					idx = closeIdx + 1
				} else {
					idx += foreachMatch[1]
				}
			}
		}
	}

	return -1
}

func findMatchingEnd(template string, startIdx int, blockType string) int {
	return findMatchingEndImpl(template, startIdx)
}

// resolveInputExpr attempts to resolve a $input.* expression string to its
// concrete value by running the input processing functions on it.
func (e *Engine) resolveInputExpr(expr string) (interface{}, bool) {
	if !strings.HasPrefix(expr, "$input") {
		return nil, false
	}

	resolved := e.processInputPath(expr)
	if resolved != expr {
		return resolved, true
	}

	resolved = e.processInputJson(expr)
	if resolved != expr {
		var parsed interface{}
		if err := json.Unmarshal([]byte(resolved), &parsed); err == nil {
			return parsed, true
		}
		return resolved, true
	}

	resolved = e.processInputBody(expr)
	if resolved != expr {
		return resolved, true
	}

	return nil, false
}

// resolveIterable resolves an expression to an iterable value (array or slice).
// It supports variable references that start with $.
func (e *Engine) resolveIterable(expr string) interface{} {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "$") {
		path := expr[1:]
		val := e.getVariableByPath(path)
		if val != nil {
			return val
		}
		if e.AppSyncCtx != nil {
			if appSyncVal := e.resolveAppSyncPath(path); appSyncVal != nil {
				return appSyncVal
			}
		}
		if inputVal, ok := e.resolveInputExpr(expr); ok {
			return inputVal
		}
	}

	return nil
}

// resolveValue resolves an expression to its value. It handles string literals
// (both double and single quoted), variable references, and boolean/null literals.
func (e *Engine) resolveValue(expr string) interface{} {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`) {
		inner := expr[1 : len(expr)-1]
		return e.interpolateStringVars(inner)
	}
	if strings.HasPrefix(expr, `'`) && strings.HasSuffix(expr, `'`) {
		return expr[1 : len(expr)-1]
	}

	if strings.HasPrefix(expr, "$") {
		path := expr[1:]
		val := e.getVariableByPath(path)
		if val != nil {
			return val
		}
		if e.AppSyncCtx != nil {
			if appSyncVal := e.resolveAppSyncPath(path); appSyncVal != nil {
				return appSyncVal
			}
		}
		if inputVal, ok := e.resolveInputExpr(expr); ok {
			return inputVal
		}
		if strings.HasPrefix(path, "util.") {
			resolved := e.processAppSyncUtil(expr)
			if resolved != expr {
				return resolved
			}
		}
		return ""
	}

	if expr == "true" {
		return true
	}
	if expr == "false" {
		return false
	}
	if expr == "null" {
		return nil
	}
	if expr == "[]" {
		return []interface{}{}
	}
	if expr == "{}" {
		return map[string]interface{}{}
	}

	return expr
}

// evaluateCondition evaluates a VTL condition expression and returns the
// boolean result. It supports logical operators (&&, ||, !), comparison
// operators (==, !=, >, <, >=, <=), and treats values as booleans.
func (e *Engine) evaluateCondition(cond string) bool {
	cond = strings.TrimSpace(cond)

	if cond == "true" {
		return true
	}
	if cond == "false" {
		return false
	}

	if strings.Contains(cond, " && ") {
		parts := strings.Split(cond, " && ")
		for _, p := range parts {
			if !e.evaluateCondition(p) {
				return false
			}
		}
		return true
	}

	if strings.Contains(cond, " || ") {
		parts := strings.Split(cond, " || ")
		for _, p := range parts {
			if e.evaluateCondition(p) {
				return true
			}
		}
		return false
	}

	if strings.HasPrefix(cond, "!") {
		return !e.evaluateCondition(cond[1:])
	}

	for _, op := range []string{"!=", "==", ">=", "<=", ">", "<"} {
		if strings.Contains(cond, op) {
			parts := strings.Split(cond, op)
			if len(parts) == 2 {
				left := e.resolveValue(strings.TrimSpace(parts[0]))
				right := e.resolveValue(strings.TrimSpace(parts[1]))
				return e.compareValues(left, right, op)
			}
		}
	}

	val := e.resolveValue(cond)
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case int, int64, float64:
		return v != 0
	default:
		return true
	}
}

// compareValues compares two values using the specified operator.
// Values are compared as strings for simplicity.
func (e *Engine) compareValues(left, right interface{}, op string) bool {
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)

	switch op {
	case "==":
		return leftStr == rightStr
	case "!=":
		return leftStr != rightStr
	case ">":
		return leftStr > rightStr
	case "<":
		return leftStr < rightStr
	case ">=":
		return leftStr >= rightStr
	case "<=":
		return leftStr <= rightStr
	}
	return false
}

func (e *Engine) resolveAppSyncPath(path string) interface{} {
	if e.AppSyncCtx == nil {
		return nil
	}

	parts := strings.SplitN(path, ".", 2)
	if len(parts) < 2 || parts[0] != "ctx" {
		return nil
	}

	subPath := parts[1]
	subParts := strings.SplitN(subPath, ".", 2)
	topLevel := subParts[0]
	fieldPath := ""
	if len(subParts) > 1 {
		fieldPath = subParts[1]
	}

	var base interface{}
	switch topLevel {
	case "args":
		base = e.AppSyncCtx.Args
	case "source":
		base = e.AppSyncCtx.Source
	case "stash":
		base = e.AppSyncCtx.Stash
	case "identity":
		base = e.AppSyncCtx.Identity
	case "result":
		base = e.AppSyncCtx.Result
	case "request":
		base = e.AppSyncCtx.Request
	case "error":
		base = e.AppSyncCtx.Error
	case "prev":
		base = e.AppSyncCtx.Prev
	case "trigger":
		base = e.AppSyncCtx.Trigger
	case "info":
		if e.AppSyncCtx.Info == nil {
			return nil
		}
		info := e.AppSyncCtx.Info
		if fieldPath == "" {
			return info
		}
		switch fieldPath {
		case "fieldName":
			return info.FieldName
		case "parentTypeName":
			return info.ParentTypeName
		case "selectionSetGraphQL":
			return info.SelectionSetGraphQL
		case "selectionSetList":
			return info.SelectionSetList
		}
		return nil
	default:
		return nil
	}

	if base == nil {
		return nil
	}
	if fieldPath == "" {
		return base
	}

	return navigatePath(base, fieldPath)
}

var stringVarRe = regexp.MustCompile(`\$([\w]+(?:\.[\w]+)*)`)

func (e *Engine) interpolateStringVars(s string) string {
	return stringVarRe.ReplaceAllStringFunc(s, func(match string) string {
		path := stringVarRe.FindStringSubmatch(match)[1]
		if strings.HasPrefix(path, "context.") || strings.HasPrefix(path, "stageVariables.") {
			return match
		}
		if val := e.getVariableByPath(path); val != nil {
			return e.formatValue(val)
		}
		if e.AppSyncCtx != nil {
			if val := e.resolveAppSyncPath(path); val != nil {
				return e.formatValue(val)
			}
		}
		return match
	})
}
