package neptunegraph

import (
	"errors"
	"fmt"
)

func NewStoreError(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("neptunegraph store: %s: %w", op, err)
}

var (
	ErrGraphNotFound           = errors.New("neptunegraph: graph not found")
	ErrGraphAlreadyExists      = errors.New("neptunegraph: graph already exists")
	ErrSnapshotNotFound        = errors.New("neptunegraph: graph snapshot not found")
	ErrSnapshotAlreadyExists   = errors.New("neptunegraph: graph snapshot already exists")
	ErrEndpointNotFound        = errors.New("neptunegraph: private graph endpoint not found")
	ErrEndpointAlreadyExists   = errors.New("neptunegraph: private graph endpoint already exists")
	ErrImportTaskNotFound      = errors.New("neptunegraph: import task not found")
	ErrImportTaskAlreadyExists = errors.New("neptunegraph: import task already exists")
	ErrExportTaskNotFound      = errors.New("neptunegraph: export task not found")
	ErrExportTaskAlreadyExists = errors.New("neptunegraph: export task already exists")
	ErrQueryNotFound           = errors.New("neptunegraph: query not found")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrGraphNotFound) ||
		errors.Is(err, ErrSnapshotNotFound) ||
		errors.Is(err, ErrEndpointNotFound) ||
		errors.Is(err, ErrImportTaskNotFound) ||
		errors.Is(err, ErrExportTaskNotFound) ||
		errors.Is(err, ErrQueryNotFound)
}

func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrGraphAlreadyExists) ||
		errors.Is(err, ErrSnapshotAlreadyExists) ||
		errors.Is(err, ErrEndpointAlreadyExists) ||
		errors.Is(err, ErrImportTaskAlreadyExists) ||
		errors.Is(err, ErrExportTaskAlreadyExists)
}
