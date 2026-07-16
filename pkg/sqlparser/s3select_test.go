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

import "testing"

func TestS3SelectParse(t *testing.T) {
	opts := ParserOptions{Dialect: DialectS3Select}

	tests := []struct {
		sql  string
		desc string
	}{
		{`SELECT * FROM s3object`, "select all"},
		{`SELECT * FROM s3object s`, "select with alias"},
		{`SELECT _1, _2 FROM s3object`, "select csv columns"},
		{`SELECT s._1, s._2 FROM s3object s`, "select with alias prefix"},
		{`SELECT * FROM s3object WHERE _1 = 'value'`, "select with where"},
		{`SELECT * FROM s3object s WHERE s._1 = 'value'`, "select with alias and where"},
		{`SELECT _1, _2 FROM s3object WHERE _3 > 100`, "select with numeric condition"},
		{`SELECT COUNT(*) FROM s3object`, "select with aggregate"},
		{`SELECT * FROM s3object WHERE _1 LIKE '%test%'`, "select with LIKE"},
		{`SELECT * FROM s3object WHERE _1 IS NULL`, "select with IS NULL"},
		{`SELECT * FROM s3object WHERE _1 IN ('a', 'b', 'c')`, "select with IN"},
		{`SELECT * FROM s3object s WHERE s._1 = 'a' AND s._2 > 5`, "select with AND"},
		{`SELECT * FROM s3object s WHERE s._1 = 'a' OR s._2 = 'b'`, "select with OR"},
		{`SELECT _1 AS col1, _2 AS col2 FROM s3object`, "select with aliases"},
		{`SELECT * FROM s3object ORDER BY _1 LIMIT 10`, "select with ORDER BY and LIMIT"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := ParseWithOptions(tt.sql, opts)
			if err != nil {
				t.Errorf("Parse error: %v\nSQL: %s", err, tt.sql)
			}
		})
	}
}

func TestS3SelectDialectQuoting(t *testing.T) {
	opts := ParserOptions{Dialect: DialectS3Select}

	tests := []struct {
		sql  string
		desc string
	}{
		{`SELECT "col1" FROM s3object`, "double quoted identifier"},
		{`SELECT * FROM s3object WHERE "col" = 'value'`, "double quote id, single quote string"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := ParseWithOptions(tt.sql, opts)
			if err != nil {
				t.Errorf("Parse error: %v\nSQL: %s", err, tt.sql)
			}
		})
	}
}
