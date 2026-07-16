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

// Dialect represents the SQL dialect for parsing.
type Dialect int

const (
	// DialectMySQL uses backticks for identifiers, double quotes for strings.
	DialectMySQL Dialect = iota

	// DialectTimestream uses double quotes for identifiers, single quotes for strings.
	// Also supports PostgreSQL-style cast operator (::).
	DialectTimestream

	// DialectAthena uses double quotes for identifiers, single quotes for strings.
	// Compatible with Timestream dialect.
	DialectAthena

	// DialectPostgreSQL uses double quotes for identifiers, single quotes for strings.
	DialectPostgreSQL

	// DialectPartiQL uses double quotes for identifiers, single quotes for strings.
	// Compatible with DynamoDB PartiQL syntax.
	DialectPartiQL

	// DialectS3Select uses double quotes for identifiers, single quotes for strings.
	// Compatible with S3 Select syntax.
	DialectS3Select
)

// ParserOptions contains configuration for the SQL parser.
type ParserOptions struct {
	// Dialect specifies the SQL dialect to use.
	Dialect Dialect
}

// DefaultParserOptions returns the default parser options (MySQL dialect).
func DefaultParserOptions() ParserOptions {
	return ParserOptions{
		Dialect: DialectMySQL,
	}
}

// currentDialect is the global dialect setting for backward compatibility.
var currentDialect = DialectMySQL

// SetDialect sets the global dialect for the parser.
// This affects the behavior of Parse() and other functions without options.
func SetDialect(d Dialect) {
	currentDialect = d
}

// GetDialect returns the current global dialect.
func GetDialect() Dialect {
	return currentDialect
}

// ParseWithOptions parses the SQL statement with the specified options.
func ParseWithOptions(sql string, opts ParserOptions) (Statement, error) {
	tokenizer := NewStringTokenizerWithOptions(sql, opts)
	if yyParse(tokenizer) != 0 {
		if tokenizer.partialDDL != nil {
			tokenizer.ParseTree = tokenizer.partialDDL
			return tokenizer.ParseTree, nil
		}
		return nil, tokenizer.LastError
	}
	return tokenizer.ParseTree, nil
}

// ParseStrictDDLWithOptions is the same as ParseWithOptions except it errors on
// partially parsed DDL statements.
func ParseStrictDDLWithOptions(sql string, opts ParserOptions) (Statement, error) {
	tokenizer := NewStringTokenizerWithOptions(sql, opts)
	if yyParse(tokenizer) != 0 {
		return nil, tokenizer.LastError
	}
	return tokenizer.ParseTree, nil
}
