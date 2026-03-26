package athena

import "errors"

var (
	// ErrWorkGroupNotFound is returned when the specified Athena work group
	// does not exist.
	ErrWorkGroupNotFound = errors.New("workgroup not found")

	// ErrWorkGroupAlreadyExists is returned when attempting to create a work group
	// that already exists.
	ErrWorkGroupAlreadyExists = errors.New("workgroup already exists")

	// ErrQueryExecutionNotFound is returned when the specified query execution
	// does not exist.
	ErrQueryExecutionNotFound = errors.New("query execution not found")

	// ErrNamedQueryNotFound is returned when the specified named query
	// does not exist.
	ErrNamedQueryNotFound = errors.New("named query not found")

	// ErrNamedQueryAlreadyExists is returned when attempting to create a named query
	// that already exists.
	ErrNamedQueryAlreadyExists = errors.New("named query already exists")

	// ErrPreparedStatementNotFound is returned when the specified prepared statement
	// does not exist.
	ErrPreparedStatementNotFound = errors.New("prepared statement not found")

	// ErrPreparedStatementAlreadyExists is returned when attempting to create a
	// prepared statement that already exists.
	ErrPreparedStatementAlreadyExists = errors.New("prepared statement already exists")

	// ErrDataCatalogNotFound is returned when the specified data catalog
	// does not exist.
	ErrDataCatalogNotFound = errors.New("data catalog not found")

	// ErrDataCatalogAlreadyExists is returned when attempting to create a data
	// catalog that already exists.
	ErrDataCatalogAlreadyExists = errors.New("data catalog already exists")

	// ErrDatabaseNotFound is returned when the specified database does not exist.
	ErrDatabaseNotFound = errors.New("database not found")

	// ErrTableNotFound is returned when the specified table does not exist.
	ErrTableNotFound = errors.New("table not found")

	// ErrInvalidRequest is returned when the request is not valid.
	ErrInvalidRequest = errors.New("invalid request")
)
