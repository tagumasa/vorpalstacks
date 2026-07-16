// Package rdbengine provides an embedded relational database engine backed by
// Pebble key-value storage, with support for row-level CRUD, secondary indexes,
// catalog management, and type-aware column encoding for MySQL and PostgreSQL
// compatibility.
package rdbengine

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNotFound is returned when a requested row, table, or database does not exist.
	ErrNotFound = errors.New("rdbengine: not found")
	// ErrAlreadyExists is returned when attempting to create a resource that already exists.
	ErrAlreadyExists = errors.New("rdbengine: already exists")
	// ErrConstraintViolation is returned when a unique constraint is violated.
	ErrConstraintViolation = errors.New("rdbengine: unique constraint violation")
)

// ColumnType represents the SQL column type for schema definitions.
type ColumnType int

const (
	ColumnTypeUnknown ColumnType = iota
	ColumnTypeBool
	ColumnTypeInt8
	ColumnTypeInt16
	ColumnTypeInt32
	ColumnTypeInt64
	ColumnTypeFloat32
	ColumnTypeFloat64
	ColumnTypeDecimal
	ColumnTypeString
	ColumnTypeBytes
	ColumnTypeDate
	ColumnTypeTimestamp
	ColumnTypeTimestampTZ
	ColumnTypeJSON
	ColumnTypeUUID
	ColumnTypeSerial
	ColumnTypeBigSerial
	ColumnTypeText
	ColumnTypeBlob
)

func (t ColumnType) String() string {
	switch t {
	case ColumnTypeBool:
		return "BOOL"
	case ColumnTypeInt8:
		return "TINYINT"
	case ColumnTypeInt16:
		return "SMALLINT"
	case ColumnTypeInt32:
		return "INT"
	case ColumnTypeInt64:
		return "BIGINT"
	case ColumnTypeFloat32:
		return "FLOAT"
	case ColumnTypeFloat64:
		return "DOUBLE"
	case ColumnTypeDecimal:
		return "DECIMAL"
	case ColumnTypeString:
		return "VARCHAR"
	case ColumnTypeBytes:
		return "BLOB"
	case ColumnTypeDate:
		return "DATE"
	case ColumnTypeTimestamp:
		return "TIMESTAMP"
	case ColumnTypeTimestampTZ:
		return "TIMESTAMPTZ"
	case ColumnTypeJSON:
		return "JSON"
	case ColumnTypeUUID:
		return "UUID"
	case ColumnTypeSerial:
		return "SERIAL"
	case ColumnTypeBigSerial:
		return "BIGSERIAL"
	case ColumnTypeText:
		return "TEXT"
	case ColumnTypeBlob:
		return "BLOB"
	default:
		return "UNKNOWN"
	}
}

// ColumnDef describes a single column within a table schema.
type ColumnDef struct {
	Name         string     `json:"name"`
	Type         ColumnType `json:"type"`
	Nullable     bool       `json:"nullable"`
	PrimaryKey   bool       `json:"primary_key"`
	AutoIncr     bool       `json:"auto_increment"`
	DefaultValue *string    `json:"default_value,omitempty"`
}

// TableSchema defines the structure of a database table.
type TableSchema struct {
	Name    string      `json:"name"`
	Columns []ColumnDef `json:"columns"`
}

// PrimaryKeyColumns returns the columns that form the primary key.
func (s *TableSchema) PrimaryKeyColumns() []ColumnDef {
	var pk []ColumnDef
	for _, c := range s.Columns {
		if c.PrimaryKey {
			pk = append(pk, c)
		}
	}
	return pk
}

// ColumnByName returns the column definition for the given name, or nil.
func (s *TableSchema) ColumnByName(name string) *ColumnDef {
	for i := range s.Columns {
		if s.Columns[i].Name == name {
			return &s.Columns[i]
		}
	}
	return nil
}

// IndexDef describes a secondary index on a table.
type IndexDef struct {
	Name    string   `json:"name"`
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// DatabaseMeta holds metadata about a database within the catalog.
type DatabaseMeta struct {
	Name      string    `json:"name"`
	Engine    string    `json:"engine"`
	CreatedAt time.Time `json:"created_at"`
}

// ColumnValue represents a typed value that can be stored in a row.
// The Go type must correspond to the ColumnType of the column.
type ColumnValue struct {
	Type  ColumnType
	Value interface{}
}

// Row is a map of column name to typed value.
type Row map[string]ColumnValue

// RowIterator provides forward iteration over rows returned by a scan.
type RowIterator interface {
	Next() bool
	Row() Row
	Error() error
	Close()
}

// ScanOptions controls the behaviour of a row scan operation.
type ScanOptions struct {
	Limit   int
	Offset  int
	StartPK []byte
	EndPK   []byte
	Reverse bool
}

// IndexScanOptions controls the behaviour of an index scan operation.
type IndexScanOptions struct {
	Limit   int
	Offset  int
	Start   []byte
	End     []byte
	Reverse bool
}

// Options configures the rdbengine initialisation.
type Options struct {
	Engine string
}

// DefaultOptions returns sensible defaults for rdbengine initialisation.
func DefaultOptions() Options {
	return Options{Engine: "mysql"}
}

func fmtErr(op string, err error) error {
	return fmt.Errorf("rdbengine: %s: %w", op, err)
}
