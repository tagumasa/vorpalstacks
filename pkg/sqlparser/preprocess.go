/*
Copyright 2024 Vorpalstacks Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlparser

import (
	"strings"
)

// Preprocessor transforms SQL queries before parsing.
// It supports chained transformations for different SQL dialects.
type Preprocessor struct {
	steps []transformFunc
}

type transformFunc func(sql string) string

// NewPreprocessor creates a new Preprocessor.
func NewPreprocessor() *Preprocessor {
	return &Preprocessor{
		steps: make([]transformFunc, 0),
	}
}

// AddTransform adds a custom transformation function.
func (p *Preprocessor) AddTransform(fn transformFunc) *Preprocessor {
	p.steps = append(p.steps, fn)
	return p
}

// Process applies all transformations to the SQL string.
func (p *Preprocessor) Process(sql string) string {
	for _, step := range p.steps {
		sql = step(sql)
	}
	return sql
}

// NewTimestreamPreprocessor creates a preprocessor for AWS Timestream queries.
// It converts PostgreSQL-style cast operator (::) to standard CAST syntax.
func NewTimestreamPreprocessor() *Preprocessor {
	return NewPreprocessor().AddCastOperator()
}

// NewAthenaPreprocessor creates a preprocessor for AWS Athena queries.
// It converts PostgreSQL-style cast operator (::) to standard CAST syntax.
func NewAthenaPreprocessor() *Preprocessor {
	return NewPreprocessor().AddCastOperator()
}

// AddCastOperator adds transformation for PostgreSQL-style cast operator (::).
// Converts "expr::type" to "CAST(expr AS type)".
func (p *Preprocessor) AddCastOperator() *Preprocessor {
	return p.AddTransform(convertCastOperator)
}

// convertCastOperator converts PostgreSQL-style cast (::) to CAST().
// It properly handles string literals by skipping content inside single quotes.
func convertCastOperator(sql string) string {
	var result strings.Builder
	inString := false
	pending := "" // Characters that might be part of a cast expression
	i := 0

	for i < len(sql) {
		ch := sql[i]

		if ch == '\'' {
			inString = !inString
			if pending != "" {
				result.WriteString(pending)
				pending = ""
			}
			result.WriteByte(ch)
			i++
			continue
		}

		if inString {
			if pending != "" {
				result.WriteString(pending)
				pending = ""
			}
			result.WriteByte(ch)
			i++
			continue
		}

		if i+1 < len(sql) && sql[i:i+2] == "::" {
			// Found ::, convert pending + CAST
			exprStart, typeName, typeEnd := extractCastPartsSimple(sql, i)
			expr := sql[exprStart:i]

			// Write any pending content before the expression
			pendingLen := i - exprStart
			if len(pending) > pendingLen {
				result.WriteString(pending[:len(pending)-pendingLen])
			}
			pending = ""

			result.WriteString("CAST(")
			result.WriteString(expr)
			result.WriteString(" AS ")
			result.WriteString(strings.ToUpper(typeName))
			result.WriteString(")")
			i = typeEnd
			continue
		}

		// Buffer characters that might be part of a cast expression
		if isIdentChar(ch) || ch == '_' || ch == '.' || ch == '"' || ch == '`' || ch == ' ' || ch == '\t' {
			pending += string(ch)
		} else {
			if pending != "" {
				result.WriteString(pending)
				pending = ""
			}
			result.WriteByte(ch)
		}
		i++
	}

	// Write any remaining pending content
	if pending != "" {
		result.WriteString(pending)
	}

	return result.String()
}

// extractCastPartsSimple extracts the identifier before :: and the type after it.
// Returns the start position of expression, type name, and end position of type.
func extractCastPartsSimple(sql string, pos int) (exprStart int, typeName string, typeEnd int) {
	typeStart := pos + 2
	typeEnd = typeStart
	for typeEnd < len(sql) && (isIdentChar(sql[typeEnd]) || sql[typeEnd] == '_') {
		typeEnd++
	}
	typeName = sql[typeStart:typeEnd]

	exprStart = pos
	for exprStart > 0 {
		ch := sql[exprStart-1]
		if isIdentChar(ch) || ch == '_' || ch == '.' || ch == '"' || ch == '`' {
			exprStart--
			continue
		}
		if ch == ')' {
			depth := 1
			exprStart--
			for exprStart > 0 && depth > 0 {
				if sql[exprStart] == ')' {
					depth++
				} else if sql[exprStart] == '(' {
					depth--
				}
				exprStart--
			}
			continue
		}
		break
	}

	return exprStart, typeName, typeEnd
}

func isIdentChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}
