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
	"testing"
)

func TestPreprocessCastOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "SELECT value::double FROM t",
			expected: "SELECT CAST(value AS DOUBLE) FROM t",
		},
		{
			input:    "SELECT name::varchar, age::bigint FROM users",
			expected: "SELECT CAST(name AS VARCHAR), CAST(age AS BIGINT) FROM users",
		},
		{
			input:    "SELECT * FROM t WHERE time::timestamp > '2024-01-01'",
			expected: "SELECT * FROM t WHERE CAST(time AS TIMESTAMP) > '2024-01-01'",
		},
		{
			input:    "SELECT 'Error::404' AS msg",
			expected: "SELECT 'Error::404' AS msg",
		},
		{
			input:    "SELECT value::double WHERE msg = 'a::b::c'",
			expected: "SELECT CAST(value AS DOUBLE) WHERE msg = 'a::b::c'",
		},
		{
			input:    "SELECT 1",
			expected: "SELECT 1",
		},
	}

	pre := NewPreprocessor().AddCastOperator()
	for _, tt := range tests {
		result := pre.Process(tt.input)
		if result != tt.expected {
			t.Errorf("Preprocess(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestTimestreamPreprocessor(t *testing.T) {
	sql := `SELECT * FROM "mydb"."mytable" WHERE time > ago(1h) AND value::double > 10`
	pre := NewTimestreamPreprocessor()
	result := pre.Process(sql)
	expected := `SELECT * FROM "mydb"."mytable" WHERE time > ago(1h) AND CAST(value AS DOUBLE) > 10`
	if result != expected {
		t.Errorf("TimestreamPreprocessor() = %q, want %q", result, expected)
	}
}

func TestParseWithDialect(t *testing.T) {
	mysqlSQL := "SELECT * FROM `users` WHERE name = 'test'"
	opts := ParserOptions{Dialect: DialectMySQL}
	stmt, err := ParseWithOptions(mysqlSQL, opts)
	if err != nil {
		t.Errorf("ParseWithOptions(MySQL) error: %v", err)
	}
	if stmt == nil {
		t.Error("ParseWithOptions(MySQL) returned nil statement")
	}

	timestreamSQL := `SELECT * FROM "mydb"."mytable" WHERE name = 'test'`
	opts = ParserOptions{Dialect: DialectTimestream}
	stmt, err = ParseWithOptions(timestreamSQL, opts)
	if err != nil {
		t.Errorf("ParseWithOptions(Timestream) error: %v", err)
	}
	if stmt == nil {
		t.Error("ParseWithOptions(Timestream) returned nil statement")
	}
}

func TestDoubleQuotedIdentifier(t *testing.T) {
	opts := ParserOptions{Dialect: DialectTimestream}

	tests := []struct {
		sql    string
		hasErr bool
	}{
		{`SELECT * FROM "database"."table"`, false},
		{`SELECT "column" FROM "table"`, false},
		{`SELECT * FROM "table" WHERE "col" = 'value'`, false},
	}

	for _, tt := range tests {
		_, err := ParseWithOptions(tt.sql, opts)
		if (err != nil) != tt.hasErr {
			t.Errorf("ParseWithOptions(%q) error = %v, hasErr %v", tt.sql, err, tt.hasErr)
		}
	}
}

func TestSetGetDialect(t *testing.T) {
	original := GetDialect()
	defer SetDialect(original)

	SetDialect(DialectTimestream)
	if GetDialect() != DialectTimestream {
		t.Errorf("GetDialect() = %v, want %v", GetDialect(), DialectTimestream)
	}

	SetDialect(DialectAthena)
	if GetDialect() != DialectAthena {
		t.Errorf("GetDialect() = %v, want %v", GetDialect(), DialectAthena)
	}
}
