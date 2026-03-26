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

func TestPartiQLSelect(t *testing.T) {
	opts := ParserOptions{Dialect: DialectPartiQL}

	tests := []struct {
		sql  string
		desc string
	}{
		{`SELECT * FROM "MyTable"`, "simple select all"},
		{`SELECT * FROM "MyTable" WHERE "pk" = 'id1'`, "select with string condition"},
		{`SELECT "id", "name" FROM "users" WHERE "status" = 'active'`, "select specific columns"},
		{`SELECT * FROM "t" WHERE "a" = 1 AND "b" > 5`, "select with AND"},
		{`SELECT * FROM "t" WHERE "a" = 1 OR "b" = 2`, "select with OR"},
		{`SELECT * FROM "t" WHERE "a" BETWEEN 1 AND 10`, "select with BETWEEN"},
		{`SELECT * FROM "t" WHERE "a" IN (1, 2, 3)`, "select with IN"},
		{`SELECT * FROM "t" WHERE "a" IS NULL`, "select with IS NULL"},
		{`SELECT * FROM "t" WHERE "a" IS NOT NULL`, "select with IS NOT NULL"},
		{`SELECT * FROM "t" ORDER BY "a" ASC LIMIT 10`, "select with ORDER BY and LIMIT"},
		{`SELECT COUNT(*) FROM "t"`, "select with aggregate"},
		{`SELECT * FROM "t" WHERE "name" LIKE '%test%'`, "select with LIKE"},
		{`SELECT * FROM "t" WHERE NOT "active" = 1`, "select with NOT"},
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

func TestPartiQLUpdate(t *testing.T) {
	opts := ParserOptions{Dialect: DialectPartiQL}

	tests := []struct {
		sql  string
		desc string
	}{
		{`UPDATE "users" SET "status" = 'inactive'`, "simple update"},
		{`UPDATE "users" SET "status" = 'inactive' WHERE "id" = '123'`, "update with where"},
		{`UPDATE "t" SET "a" = 1, "b" = 2 WHERE "pk" = 'x'`, "update multiple columns"},
		{`UPDATE "t" SET "count" = "count" + 1 WHERE "pk" = 'x'`, "update with expression"},
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

func TestPartiQLDelete(t *testing.T) {
	opts := ParserOptions{Dialect: DialectPartiQL}

	tests := []struct {
		sql  string
		desc string
	}{
		{`DELETE FROM "users" WHERE "id" = '123'`, "delete with where"},
		{`DELETE FROM "t" WHERE "a" = 1 AND "b" = 2`, "delete with AND"},
		{`DELETE FROM "t" WHERE "a" IN (1, 2, 3)`, "delete with IN"},
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

func TestPartiQLDialectQuoting(t *testing.T) {
	opts := ParserOptions{Dialect: DialectPartiQL}

	tests := []struct {
		sql  string
		desc string
	}{
		{`SELECT * FROM "MyTable"`, "double quoted identifier"},
		{`SELECT * FROM "MyTable" WHERE "col" = 'value'`, "double quote id, single quote string"},
		{`SELECT "firstName", "lastName" FROM "Users"`, "camelCase identifiers"},
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

func TestPartiQLInsertValue(t *testing.T) {
	opts := ParserOptions{Dialect: DialectPartiQL}

	tests := []struct {
		sql  string
		desc string
	}{
		{`INSERT INTO "t" VALUE {}`, "empty object"},
		{`INSERT INTO "t" VALUE {'id': '1'}`, "single property string key"},
		{`INSERT INTO "t" VALUE {'id': 1}`, "single property int value"},
		{`INSERT INTO "t" VALUE {'id': '1', 'name': 'test'}`, "multiple properties"},
		{`INSERT INTO "MyTable" VALUE {'pk': 'id1', 'sk': 'sort1', 'data': 'value'}`, "three properties"},
		{`INSERT INTO "t" VALUE {'count': 100, 'active': true}`, "int and bool"},
		{`INSERT INTO "t" VALUE {'id': 'x', 'val': null}`, "with null"},
		{`INSERT INTO "t" VALUE {'a': 1, 'b': 2, 'c': 3, 'd': 4}`, "many properties"},
		{`INSERT INTO "t" VALUE {id: '1'}`, "identifier key"},
		{`INSERT INTO "t" VALUE {id: 1, name: 'test'}`, "multiple identifier keys"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			stmt, err := ParseWithOptions(tt.sql, opts)
			if err != nil {
				t.Errorf("Parse error: %v\nSQL: %s", err, tt.sql)
				return
			}

			ins, ok := stmt.(*Insert)
			if !ok {
				t.Errorf("Expected *Insert, got %T", stmt)
				return
			}

			values, ok := ins.Rows.(Values)
			if !ok {
				t.Errorf("Expected Values, got %T", ins.Rows)
				return
			}

			if len(values) != 1 {
				t.Errorf("Expected 1 row, got %d", len(values))
			}
		})
	}
}

func TestPartiQLInsertValueAST(t *testing.T) {
	opts := ParserOptions{Dialect: DialectPartiQL}

	// First test: simple object
	t.Run("simple_string_value", func(t *testing.T) {
		sql := `INSERT INTO "t" VALUE {'a': '1'}`
		stmt, err := ParseWithOptions(sql, opts)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		_ = stmt
	})

	t.Run("full_test", func(t *testing.T) {
		sql := `INSERT INTO "Music" VALUE {'Artist': 'Acme Band', 'SongTitle': 'PartiQL Rocks'}`
		stmt, err := ParseWithOptions(sql, opts)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		ins, ok := stmt.(*Insert)
		if !ok {
			t.Fatalf("Expected *Insert, got %T", stmt)
		}

		if ins.Table.Name.String() != "Music" {
			t.Errorf("Expected table name 'Music', got %s", ins.Table.Name.String())
		}

		values, ok := ins.Rows.(Values)
		if !ok {
			t.Fatalf("Expected Values, got %T", ins.Rows)
		}

		if len(values) != 1 {
			t.Fatalf("Expected 1 value tuple, got %d", len(values))
		}

		tuple := values[0]
		if len(tuple) != 1 {
			t.Fatalf("Expected 1 element in tuple, got %d", len(tuple))
		}

		obj, ok := tuple[0].(*ObjectLiteral)
		if !ok {
			t.Fatalf("Expected *ObjectLiteral, got %T", tuple[0])
		}

		if len(obj.Properties) != 2 {
			t.Errorf("Expected 2 properties, got %d", len(obj.Properties))
		}

		// Check first property
		prop0 := obj.Properties[0]
		if key, ok := prop0.Key.(*SQLVal); ok {
			if string(key.Val) != "Artist" {
				t.Errorf("Expected key 'Artist', got %s", string(key.Val))
			}
		} else {
			t.Errorf("Expected *SQLVal for key, got %T", prop0.Key)
		}

		// Check second property
		prop1 := obj.Properties[1]
		if key, ok := prop1.Key.(*SQLVal); ok {
			if string(key.Val) != "SongTitle" {
				t.Errorf("Expected key 'SongTitle', got %s", string(key.Val))
			}
		} else {
			t.Errorf("Expected *SQLVal for key, got %T", prop1.Key)
		}
	})
}
