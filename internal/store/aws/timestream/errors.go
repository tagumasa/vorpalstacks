// Package timestream provides Timestream storage functionality for vorpalstacks.
package timestream

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrDatabaseNotFound is returned when a database is not found.
	ErrDatabaseNotFound = errors.New("database not found")
	// ErrDatabaseAlreadyExists is returned when a database already exists.
	ErrDatabaseAlreadyExists = errors.New("database already exists")
	// ErrDatabaseNotEmpty is returned when a database is not empty.
	ErrDatabaseNotEmpty = errors.New("database is not empty")
	// ErrTableNotFound is returned when a table is not found.
	ErrTableNotFound = errors.New("table not found")
	// ErrTableAlreadyExists is returned when a table already exists.
	ErrTableAlreadyExists = errors.New("table already exists")
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = common.ErrInvalidParameter
	// ErrResourceNotFound is returned when a resource is not found.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrScheduledQueryNotFound is returned when a scheduled query is not found.
	ErrScheduledQueryNotFound = errors.New("scheduled query not found")
	// ErrScheduledQueryAlreadyExists is returned when a scheduled query already exists.
	ErrScheduledQueryAlreadyExists = errors.New("scheduled query already exists")
	// ErrScheduledQueryRunNotFound is returned when a scheduled query run is not found.
	ErrScheduledQueryRunNotFound = errors.New("scheduled query run not found")
	// ErrBatchLoadTaskNotFound is returned when a batch load task is not found.
	ErrBatchLoadTaskNotFound = errors.New("batch load task not found")
	// ErrBatchLoadTaskAlreadyExists is returned when a batch load task already exists.
	ErrBatchLoadTaskAlreadyExists = errors.New("batch load task already exists")
)
