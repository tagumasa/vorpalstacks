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

package filterpattern

// MetricFilter represents a CloudWatch Logs metric filter configuration.
// It defines a filter pattern and optional metric transformations for
// extracting metrics from log events.
type MetricFilter struct {
	FilterName            string
	LogGroupName          string
	FilterPattern         string
	MetricTransformations []MetricTransformation
}

// MetricTransformation defines how to transform a filtered log event
// into a CloudWatch metric.
type MetricTransformation struct {
	MetricName      string
	MetricNamespace string
	MetricValue     string
	DefaultValue    float64
	Unit            string
}

// PatternType represents the type of filter pattern being used.
// Different pattern types require different matching strategies.
type PatternType int

const (
	// PatternTypeUnstructured is a plain text or keyword pattern.
	PatternTypeUnstructured PatternType = iota
	// PatternTypeJSON is a JSON path-based pattern for structured logs.
	PatternTypeJSON
	// PatternTypeDelimited is a pattern for space/comma-delimited logs.
	PatternTypeDelimited
)

// CompiledPattern represents a parsed and compiled filter pattern
// ready for matching against log events.
type CompiledPattern struct {
	Type     PatternType
	Original string
	AST      any
}
