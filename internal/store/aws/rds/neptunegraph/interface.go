package neptunegraph

import (
	"vorpalstacks/internal/store/aws/common"
)

// NeptuneGraphStoreInterface combines all NeptuneGraph store operations into a
// single interface for use by service handlers.
type NeptuneGraphStoreInterface interface {
	GraphOps
	GraphSnapshotOps
	PrivateGraphEndpointOps
	TagOps
	QueryOps
	ImportTaskOps
	ExportTaskOps
}

// GraphOps defines CRUD operations for NeptuneGraph graph resources.
type GraphOps interface {
	CreateGraph(graph *Graph) error
	GetGraph(id string) (*Graph, error)
	UpdateGraph(graph *Graph) error
	DeleteGraph(id string) error
	ListGraphs(opts common.ListOptions) ([]*Graph, string, bool, error)
}

// GraphSnapshotOps defines CRUD operations for NeptuneGraph snapshots.
type GraphSnapshotOps interface {
	CreateSnapshot(snapshot *GraphSnapshot) error
	GetSnapshot(id string) (*GraphSnapshot, error)
	UpdateSnapshot(snapshot *GraphSnapshot) error
	DeleteSnapshot(id string) error
	ListSnapshots(opts common.ListOptions, graphId string) ([]*GraphSnapshot, string, bool, error)
}

// PrivateGraphEndpointOps defines operations for private graph VPC endpoints.
type PrivateGraphEndpointOps interface {
	CreateEndpoint(ep *PrivateGraphEndpoint) error
	GetEndpoint(graphId, vpcId string) (*PrivateGraphEndpoint, error)
	DeleteEndpoint(graphId, vpcId string) error
	ListEndpoints(graphId string) ([]*PrivateGraphEndpoint, error)
}

// TagOps defines operations for resource tagging.
type TagOps interface {
	AddTags(resourceArn string, tags map[string]string) error
	GetTags(resourceArn string) (map[string]string, error)
	RemoveTags(resourceArn string, keys []string) error
}

// QueryOps defines CRUD operations for graph query records.
type QueryOps interface {
	CreateQuery(q *QueryRecord) error
	GetQuery(graphId, id string) (*QueryRecord, error)
	UpdateQuery(q *QueryRecord) error
	DeleteQuery(graphId, id string) error
	ListQueries(graphId string, maxResults int) ([]*QueryRecord, error)
}

// ImportTaskOps defines CRUD and state transition operations for import tasks.
type ImportTaskOps interface {
	CreateImportTask(t *ImportTask) error
	GetImportTask(id string) (*ImportTask, error)
	UpdateImportTask(t *ImportTask) error
	DeleteImportTask(id string) error
	ListImportTasks(opts common.ListOptions) ([]*ImportTask, string, bool, error)
	TryAdvanceImportTask(id string, targetStatus string, apply func(*ImportTask)) error
}

// ExportTaskOps defines CRUD and state transition operations for export tasks.
type ExportTaskOps interface {
	CreateExportTask(t *ExportTask) error
	GetExportTask(id string) (*ExportTask, error)
	UpdateExportTask(t *ExportTask) error
	DeleteExportTask(id string) error
	ListExportTasks(opts common.ListOptions, graphId string) ([]*ExportTask, string, bool, error)
	TryAdvanceExportTask(id string, targetStatus string, apply func(*ExportTask)) error
}

var _ NeptuneGraphStoreInterface = (*NeptuneGraphStore)(nil)
