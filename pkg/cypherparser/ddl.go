// Cypher DDL executor — handles CREATE INDEX, CREATE CONSTRAINT, SHOW INDEXES,
// SHOW CONSTRAINTS, DROP INDEX, DROP CONSTRAINT.
//
// Original implementation — goraphdb does not have DDL support.
//
// Execution model:
//   - CREATE INDEX: registers a property index for a label+property pair
//   - CREATE CONSTRAINT: registers a unique constraint for a label+property pair
//   - SHOW INDEXES: lists all registered property indexes
//   - SHOW CONSTRAINTS: lists all registered unique constraints
//   - DROP INDEX: removes a registered property index
//   - DROP CONSTRAINT: removes a registered unique constraint

package cypherparser

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/storage/graphengine"
)

// ExecuteDDL runs a parsed DDL statement against the graph store.
func ExecuteDDL(_ context.Context, ddl graphengine.GraphDDL, stmt *DDLStatement) (*CypherResult, error) {
	switch stmt.Kind {
	case "INDEX":
		if err := ddl.CreateIndex(stmt.Label, stmt.Property); err != nil {
			return nil, fmt.Errorf("cypher DDL: %w", err)
		}
		return &CypherResult{
			Columns: []string{"label", "property"},
			Rows: []map[string]any{
				{"label": stmt.Label, "property": stmt.Property},
			},
		}, nil

	case "CONSTRAINT":
		if err := ddl.CreateUniqueConstraint(stmt.Label, stmt.Property); err != nil {
			return nil, fmt.Errorf("cypher DDL: %w", err)
		}
		return &CypherResult{
			Columns: []string{"label", "property", "type"},
			Rows: []map[string]any{
				{"label": stmt.Label, "property": stmt.Property, "type": "UNIQUE"},
			},
		}, nil

	case "SHOW_INDEXES":
		indexes, err := ddl.ShowIndexes()
		if err != nil {
			return nil, fmt.Errorf("cypher DDL: %w", err)
		}
		rows := make([]map[string]any, 0, len(indexes))
		for _, idx := range indexes {
			rows = append(rows, map[string]any{"label": idx.Label, "property": idx.Property})
		}
		return &CypherResult{
			Columns: []string{"label", "property"},
			Rows:    rows,
		}, nil

	case "SHOW_CONSTRAINTS":
		constraints, err := ddl.ShowConstraints()
		if err != nil {
			return nil, fmt.Errorf("cypher DDL: %w", err)
		}
		rows := make([]map[string]any, 0, len(constraints))
		for _, c := range constraints {
			rows = append(rows, map[string]any{"label": c.Label, "property": c.Property, "type": "UNIQUE"})
		}
		return &CypherResult{
			Columns: []string{"label", "property", "type"},
			Rows:    rows,
		}, nil

	case "DROP_INDEX":
		if err := ddl.DropIndex(stmt.Label, stmt.Property); err != nil {
			return nil, fmt.Errorf("cypher DDL: %w", err)
		}
		return &CypherResult{
			Columns: []string{},
			Rows:    []map[string]any{},
		}, nil

	case "DROP_CONSTRAINT":
		if err := ddl.DropUniqueConstraint(stmt.Label, stmt.Property); err != nil {
			return nil, fmt.Errorf("cypher DDL: %w", err)
		}
		return &CypherResult{
			Columns: []string{},
			Rows:    []map[string]any{},
		}, nil

	default:
		return nil, fmt.Errorf("cypher DDL: unsupported DDL kind: %s", stmt.Kind)
	}
}
