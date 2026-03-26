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

// Context represents the execution context for VTL template processing.
// It contains all the data that can be accessed within templates including
// request body, parameters, context variables, stage variables, and
// request/response overrides.
type Context struct {
	Body             string
	JSONBody         interface{}
	Params           map[string]string
	Context          map[string]interface{}
	StageVars        map[string]string
	Authorizer       map[string]interface{}
	Identity         map[string]interface{}
	RequestOverride  map[string]string
	ResponseOverride map[string]string
	RequestTime      string
	HTTPMethod       string
	ResourcePath     string
	Path             string
}

// ContextLoopMetadata contains metadata about the current iteration in a
// #foreach loop. This information is available via the $foreach variable
// within loop bodies.
type ContextLoopMetadata struct {
	Index   int
	Count   int
	First   bool
	Last    bool
	HasNext bool
	HasPrev bool
}

// Engine is the core VTL template processing engine. It maintains the
// execution context and provides methods for transforming templates.
// The engine processes templates through multiple phases: control flow,
// input parsing, utility functions, context variables, and stage variables.
type Engine struct {
	context   *Context
	variables map[string]interface{}
}

// NewEngine creates a new VTL engine instance with default context
// and empty variable storage.
func NewEngine() *Engine {
	return &Engine{
		context: &Context{
			Params:           make(map[string]string),
			Context:          make(map[string]interface{}),
			StageVars:        make(map[string]string),
			Authorizer:       make(map[string]interface{}),
			Identity:         make(map[string]interface{}),
			RequestOverride:  make(map[string]string),
			ResponseOverride: make(map[string]string),
		},
		variables: make(map[string]interface{}),
	}
}

// SetContext sets the execution context for the engine. The context
// contains all the data that will be available during template processing.
func (e *Engine) SetContext(ctx *Context) {
	if ctx != nil {
		e.context = ctx
	}
}

// GetContext returns the current execution context.
func (e *Engine) GetContext() *Context {
	return e.context
}
