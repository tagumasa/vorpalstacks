package neptunegraph

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrGraphNotFound is returned when the requested NeptuneGraph graph does not exist.
	ErrGraphNotFound = errors.New("neptunegraph: graph not found")
	// ErrGraphAlreadyExists is returned when attempting to create a graph that already exists.
	ErrGraphAlreadyExists = errors.New("neptunegraph: graph already exists")
	// ErrSnapshotNotFound is returned when the requested graph snapshot does not exist.
	ErrSnapshotNotFound = errors.New("neptunegraph: graph snapshot not found")
	// ErrSnapshotAlreadyExists is returned when attempting to create a snapshot that already exists.
	ErrSnapshotAlreadyExists = errors.New("neptunegraph: graph snapshot already exists")
	// ErrEndpointNotFound is returned when the requested private graph endpoint does not exist.
	ErrEndpointNotFound = errors.New("neptunegraph: private graph endpoint not found")
	// ErrEndpointAlreadyExists is returned when attempting to create an endpoint that already exists.
	ErrEndpointAlreadyExists = errors.New("neptunegraph: private graph endpoint already exists")
	// ErrImportTaskNotFound is returned when the requested import task does not exist.
	ErrImportTaskNotFound = errors.New("neptunegraph: import task not found")
	// ErrImportTaskAlreadyExists is returned when attempting to create an import task that already exists.
	ErrImportTaskAlreadyExists = errors.New("neptunegraph: import task already exists")
	// ErrExportTaskNotFound is returned when the requested export task does not exist.
	ErrExportTaskNotFound = errors.New("neptunegraph: export task not found")
	// ErrExportTaskAlreadyExists is returned when attempting to create an export task that already exists.
	ErrExportTaskAlreadyExists = errors.New("neptunegraph: export task already exists")
	// ErrQueryNotFound is returned when the requested query does not exist.
	ErrQueryNotFound = errors.New("neptunegraph: query not found")
)

// IsNotFound reports whether err represents a NeptuneGraph not-found condition.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrGraphNotFound) ||
		errors.Is(err, ErrSnapshotNotFound) ||
		errors.Is(err, ErrEndpointNotFound) ||
		errors.Is(err, ErrImportTaskNotFound) ||
		errors.Is(err, ErrExportTaskNotFound) ||
		errors.Is(err, ErrQueryNotFound) ||
		common.IsNotFound(err)
}

// IsAlreadyExists reports whether err represents a NeptuneGraph already-exists condition.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrGraphAlreadyExists) ||
		errors.Is(err, ErrSnapshotAlreadyExists) ||
		errors.Is(err, ErrEndpointAlreadyExists) ||
		errors.Is(err, ErrImportTaskAlreadyExists) ||
		errors.Is(err, ErrExportTaskAlreadyExists) ||
		common.IsAlreadyExists(err)
}
