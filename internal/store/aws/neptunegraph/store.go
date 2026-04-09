package neptunegraph

import (
	"strings"
	"sync"

	pb "vorpalstacks/internal/pb/storage/storage_neptunegraph"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const (
	graphsBucket             = "neptunegraph_graphs"
	snapshotsBucket          = "neptunegraph_snapshots"
	snapshotsByGraphBucket   = "neptunegraph_snapshots_by_graph"
	endpointsBucket          = "neptunegraph_endpoints"
	tagsBucket               = "neptunegraph_tags"
	queriesBucket            = "neptunegraph_queries"
	importTasksBucket        = "neptunegraph_import_tasks"
	exportTasksBucket        = "neptunegraph_export_tasks"
	exportTasksByGraphBucket = "neptunegraph_export_tasks_by_graph"
)

// NeptuneGraphStoreInterface combines all NeptuneGraph store operations into a single interface.
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

// NeptuneGraphStore provides persistent storage for NeptuneGraph resources using underlying key-value buckets.
type NeptuneGraphStore struct {
	graphs             *common.BaseStore
	snapshots          *common.BaseStore
	snapshotsByGraph   *common.BaseStore
	endpoints          *common.BaseStore
	tags               *common.BaseStore
	queries            *common.BaseStore
	importTasks        *common.BaseStore
	exportTasks        *common.BaseStore
	exportTasksByGraph *common.BaseStore
	mu                 sync.RWMutex
}

var _ NeptuneGraphStoreInterface = (*NeptuneGraphStore)(nil)

// NewNeptuneGraphStore creates a new NeptuneGraphStore backed by the given storage.
func NewNeptuneGraphStore(store storage.BasicStorage) *NeptuneGraphStore {
	return &NeptuneGraphStore{
		graphs:             common.NewBaseStore(store.Bucket(graphsBucket), "neptunegraph"),
		snapshots:          common.NewBaseStore(store.Bucket(snapshotsBucket), "neptunegraph"),
		snapshotsByGraph:   common.NewBaseStore(store.Bucket(snapshotsByGraphBucket), "neptunegraph"),
		endpoints:          common.NewBaseStore(store.Bucket(endpointsBucket), "neptunegraph"),
		tags:               common.NewBaseStore(store.Bucket(tagsBucket), "neptunegraph"),
		queries:            common.NewBaseStore(store.Bucket(queriesBucket), "neptunegraph"),
		importTasks:        common.NewBaseStore(store.Bucket(importTasksBucket), "neptunegraph"),
		exportTasks:        common.NewBaseStore(store.Bucket(exportTasksBucket), "neptunegraph"),
		exportTasksByGraph: common.NewBaseStore(store.Bucket(exportTasksByGraphBucket), "neptunegraph"),
	}
}

// CreateGraph persists a new graph, returning an error if the graph already exists.
func (s *NeptuneGraphStore) CreateGraph(graph *Graph) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := graph.Id
	if s.graphs.Exists(key) {
		return ErrGraphAlreadyExists
	}
	return NewStoreError("create_graph", s.graphs.PutProto(key, graphToProto(graph)))
}

// GetGraph retrieves a graph by its identifier.
func (s *NeptuneGraphStore) GetGraph(id string) (*Graph, error) {
	var p pb.Graph
	if err := s.graphs.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrGraphNotFound
		}
		return nil, NewStoreError("get_graph", err)
	}
	return protoToGraph(&p), nil
}

// UpdateGraph persists updated graph properties.
func (s *NeptuneGraphStore) UpdateGraph(graph *Graph) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.graphs.Exists(graph.Id) {
		return ErrGraphNotFound
	}
	return NewStoreError("update_graph", s.graphs.PutProto(graph.Id, graphToProto(graph)))
}

// DeleteGraph removes a graph and all associated queries, snapshots, and export task indices.
func (s *NeptuneGraphStore) DeleteGraph(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.graphs.Exists(id) {
		return ErrGraphNotFound
	}
	if err := s.graphs.Delete(id); err != nil {
		return NewStoreError("delete_graph", err)
	}
	s.queries.DeleteByPrefix(id + "/")
	s.snapshotsByGraph.DeleteByPrefix(id + "/")
	s.exportTasksByGraph.DeleteByPrefix(id + "/")
	return nil
}

// ListGraphs returns a paginated list of graphs.
func (s *NeptuneGraphStore) ListGraphs(opts common.ListOptions) ([]*Graph, string, bool, error) {
	result, err := common.ListProto[*pb.Graph](s.graphs, opts, func() *pb.Graph { return &pb.Graph{} }, nil)
	if err != nil {
		return nil, "", false, NewStoreError("list_graphs", err)
	}
	graphs := make([]*Graph, len(result.Items))
	for i, item := range result.Items {
		graphs[i] = protoToGraph(item)
	}
	return graphs, result.NextMarker, result.IsTruncated, nil
}

// CreateSnapshot persists a new graph snapshot and maintains a secondary index by source graph.
func (s *NeptuneGraphStore) CreateSnapshot(snapshot *GraphSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.snapshots.Exists(snapshot.Id) {
		return ErrSnapshotAlreadyExists
	}
	if err := s.snapshots.PutProto(snapshot.Id, snapshotToProto(snapshot)); err != nil {
		return NewStoreError("create_snapshot", err)
	}
	idxKey := snapshot.SourceGraphId + "/" + snapshot.Id
	return NewStoreError("create_snapshot_idx", s.snapshotsByGraph.Put(idxKey, ""))
}

// GetSnapshot retrieves a snapshot by its identifier.
func (s *NeptuneGraphStore) GetSnapshot(id string) (*GraphSnapshot, error) {
	var p pb.GraphSnapshot
	if err := s.snapshots.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrSnapshotNotFound
		}
		return nil, NewStoreError("get_snapshot", err)
	}
	return protoToSnapshot(&p), nil
}

// UpdateSnapshot persists updated snapshot properties.
func (s *NeptuneGraphStore) UpdateSnapshot(snapshot *GraphSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.snapshots.Exists(snapshot.Id) {
		return ErrSnapshotNotFound
	}
	return NewStoreError("update_snapshot", s.snapshots.PutProto(snapshot.Id, snapshotToProto(snapshot)))
}

// DeleteSnapshot removes a snapshot and its secondary graph index entry.
func (s *NeptuneGraphStore) DeleteSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var p pb.GraphSnapshot
	if err := s.snapshots.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return ErrSnapshotNotFound
		}
		return NewStoreError("delete_snapshot", err)
	}
	idxKey := p.GetSourceGraphId() + "/" + id
	s.snapshotsByGraph.Delete(idxKey)
	return NewStoreError("delete_snapshot", s.snapshots.Delete(id))
}

// ListSnapshots returns a paginated list of snapshots, optionally filtered by source graph identifier.
func (s *NeptuneGraphStore) ListSnapshots(opts common.ListOptions, graphId string) ([]*GraphSnapshot, string, bool, error) {
	if graphId != "" {
		var ids []string
		err := s.snapshotsByGraph.ScanPrefix(graphId+"/", func(key string, value []byte) error {
			parts := strings.SplitN(key, "/", 2)
			if len(parts) == 2 {
				ids = append(ids, parts[1])
			}
			return nil
		})
		if err != nil {
			return nil, "", false, NewStoreError("list_snapshots_idx", err)
		}

		startIdx := 0
		if opts.Marker != "" {
			for i, id := range ids {
				if id == opts.Marker {
					startIdx = i + 1
					break
				}
			}
		}

		if startIdx >= len(ids) {
			return nil, "", false, nil
		}

		ids = ids[startIdx:]
		if opts.MaxItems > 0 && len(ids) > opts.MaxItems {
			ids = ids[:opts.MaxItems]
		}

		snapshots := make([]*GraphSnapshot, 0, len(ids))
		for _, id := range ids {
			snap, err := s.GetSnapshot(id)
			if err != nil {
				continue
			}
			snapshots = append(snapshots, snap)
		}

		truncated := false
		nextToken := ""
		if opts.MaxItems > 0 && len(snapshots) == opts.MaxItems {
			nextToken = snapshots[len(snapshots)-1].Id
			truncated = true
		}
		return snapshots, nextToken, truncated, nil
	}
	result, err := common.ListProto[*pb.GraphSnapshot](s.snapshots, opts, func() *pb.GraphSnapshot { return &pb.GraphSnapshot{} }, nil)
	if err != nil {
		return nil, "", false, NewStoreError("list_snapshots", err)
	}
	snapshots := make([]*GraphSnapshot, len(result.Items))
	for i, item := range result.Items {
		snapshots[i] = protoToSnapshot(item)
	}
	return snapshots, result.NextMarker, result.IsTruncated, nil
}

// CreateEndpoint persists a new private graph VPC endpoint.
func (s *NeptuneGraphStore) CreateEndpoint(ep *PrivateGraphEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := ep.GraphId + ":" + ep.VpcId
	if s.endpoints.Exists(key) {
		return ErrEndpointAlreadyExists
	}
	return NewStoreError("create_endpoint", s.endpoints.PutProto(key, endpointToProto(ep)))
}

// GetEndpoint retrieves a private graph endpoint by graph and VPC identifiers.
func (s *NeptuneGraphStore) GetEndpoint(graphId, vpcId string) (*PrivateGraphEndpoint, error) {
	key := graphId + ":" + vpcId
	var p pb.PrivateGraphEndpoint
	if err := s.endpoints.GetProto(key, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEndpointNotFound
		}
		return nil, NewStoreError("get_endpoint", err)
	}
	return protoToEndpoint(&p), nil
}

// DeleteEndpoint removes a private graph endpoint identified by graph and VPC identifiers.
func (s *NeptuneGraphStore) DeleteEndpoint(graphId, vpcId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := graphId + ":" + vpcId
	if !s.endpoints.Exists(key) {
		return ErrEndpointNotFound
	}
	return NewStoreError("delete_endpoint", s.endpoints.Delete(key))
}

// ListEndpoints returns all private graph endpoints for a given graph identifier.
func (s *NeptuneGraphStore) ListEndpoints(graphId string) ([]*PrivateGraphEndpoint, error) {
	prefix := graphId + ":"
	var endpoints []*PrivateGraphEndpoint
	err := s.endpoints.ScanPrefix(prefix, func(key string, value []byte) error {
		var p pb.PrivateGraphEndpoint
		if err := proto.Unmarshal(value, &p); err != nil {
			return err
		}
		endpoints = append(endpoints, protoToEndpoint(&p))
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_endpoints", err)
	}
	return endpoints, nil
}

// AddTags merges new tags into the existing tag set for a resource.
func (s *NeptuneGraphStore) AddTags(resourceArn string, tags map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing := make(map[string]string)
	if s.tags.Exists(resourceArn) {
		var tp pb.TagMap
		if err := s.tags.GetProto(resourceArn, &tp); err == nil {
			existing = tp.GetTags()
		}
	}
	for k, v := range tags {
		existing[k] = v
	}
	return NewStoreError("add_tags", s.tags.PutProto(resourceArn, tagsToProto(existing)))
}

// GetTags retrieves all tags associated with a resource ARN.
func (s *NeptuneGraphStore) GetTags(resourceArn string) (map[string]string, error) {
	var tp pb.TagMap
	if err := s.tags.GetProto(resourceArn, &tp); err != nil {
		if common.IsNotFound(err) {
			return map[string]string{}, nil
		}
		return nil, NewStoreError("get_tags", err)
	}
	return protoToTags(&tp), nil
}

// RemoveTags deletes specified tag keys from a resource's tag set.
func (s *NeptuneGraphStore) RemoveTags(resourceArn string, keys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing := make(map[string]string)
	if s.tags.Exists(resourceArn) {
		var tp pb.TagMap
		if err := s.tags.GetProto(resourceArn, &tp); err == nil {
			existing = tp.GetTags()
		}
	}
	for _, k := range keys {
		delete(existing, k)
	}
	return NewStoreError("remove_tags", s.tags.PutProto(resourceArn, tagsToProto(existing)))
}

// CreateQuery persists a new query record.
func (s *NeptuneGraphStore) CreateQuery(q *QueryRecord) error {
	key := q.GraphId + "/" + q.Id
	return NewStoreError("create_query", s.queries.PutProto(key, queryToProto(q)))
}

// GetQuery retrieves a query record by graph and query identifiers.
func (s *NeptuneGraphStore) GetQuery(graphId, id string) (*QueryRecord, error) {
	var p pb.QueryRecord
	if err := s.queries.GetProto(graphId+"/"+id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrQueryNotFound
		}
		return nil, NewStoreError("get_query", err)
	}
	return protoToQuery(&p), nil
}

// UpdateQuery persists updated query record fields.
func (s *NeptuneGraphStore) UpdateQuery(q *QueryRecord) error {
	key := q.GraphId + "/" + q.Id
	return NewStoreError("update_query", s.queries.PutProto(key, queryToProto(q)))
}

// DeleteQuery removes a query record by graph and query identifiers.
func (s *NeptuneGraphStore) DeleteQuery(graphId, id string) error {
	key := graphId + "/" + id
	if !s.queries.Exists(key) {
		return ErrQueryNotFound
	}
	return NewStoreError("delete_query", s.queries.Delete(key))
}

// ListQueries returns query records for a graph, optionally limited by maxResults.
func (s *NeptuneGraphStore) ListQueries(graphId string, maxResults int) ([]*QueryRecord, error) {
	prefix := graphId + "/"
	var queries []*QueryRecord
	err := s.queries.ScanPrefix(prefix, func(key string, value []byte) error {
		var p pb.QueryRecord
		if err := proto.Unmarshal(value, &p); err != nil {
			return err
		}
		queries = append(queries, protoToQuery(&p))
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_queries", err)
	}
	if maxResults > 0 && len(queries) > maxResults {
		queries = queries[:maxResults]
	}
	return queries, nil
}

// CreateImportTask persists a new import task, returning an error if the task already exists.
func (s *NeptuneGraphStore) CreateImportTask(t *ImportTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.importTasks.Exists(t.TaskId) {
		return ErrImportTaskAlreadyExists
	}
	return NewStoreError("create_import_task", s.importTasks.PutProto(t.TaskId, importTaskToProto(t)))
}

// GetImportTask retrieves an import task by its identifier.
func (s *NeptuneGraphStore) GetImportTask(id string) (*ImportTask, error) {
	var p pb.ImportTaskRecord
	if err := s.importTasks.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrImportTaskNotFound
		}
		return nil, NewStoreError("get_import_task", err)
	}
	return protoToImportTask(&p), nil
}

// UpdateImportTask persists updated import task properties.
func (s *NeptuneGraphStore) UpdateImportTask(t *ImportTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("update_import_task", s.importTasks.PutProto(t.TaskId, importTaskToProto(t)))
}

// DeleteImportTask removes an import task by its identifier.
func (s *NeptuneGraphStore) DeleteImportTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.importTasks.Exists(id) {
		return ErrImportTaskNotFound
	}
	return NewStoreError("delete_import_task", s.importTasks.Delete(id))
}

// ListImportTasks returns a paginated list of import tasks.
func (s *NeptuneGraphStore) ListImportTasks(opts common.ListOptions) ([]*ImportTask, string, bool, error) {
	result, err := common.ListProto[*pb.ImportTaskRecord](s.importTasks, opts, func() *pb.ImportTaskRecord { return &pb.ImportTaskRecord{} }, nil)
	if err != nil {
		return nil, "", false, NewStoreError("list_import_tasks", err)
	}
	tasks := make([]*ImportTask, len(result.Items))
	for i, item := range result.Items {
		tasks[i] = protoToImportTask(item)
	}
	return tasks, result.NextMarker, result.IsTruncated, nil
}

// CreateExportTask persists a new export task and maintains a secondary index by graph.
func (s *NeptuneGraphStore) CreateExportTask(t *ExportTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.exportTasks.Exists(t.TaskId) {
		return ErrExportTaskAlreadyExists
	}
	if err := s.exportTasks.PutProto(t.TaskId, exportTaskToProto(t)); err != nil {
		return NewStoreError("create_export_task", err)
	}
	if t.GraphId != "" {
		idxKey := t.GraphId + "/" + t.TaskId
		return NewStoreError("create_export_task_idx", s.exportTasksByGraph.Put(idxKey, ""))
	}
	return nil
}

// GetExportTask retrieves an export task by its identifier.
func (s *NeptuneGraphStore) GetExportTask(id string) (*ExportTask, error) {
	var p pb.ExportTaskRecord
	if err := s.exportTasks.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrExportTaskNotFound
		}
		return nil, NewStoreError("get_export_task", err)
	}
	return protoToExportTask(&p), nil
}

// UpdateExportTask persists updated export task properties.
func (s *NeptuneGraphStore) UpdateExportTask(t *ExportTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("update_export_task", s.exportTasks.PutProto(t.TaskId, exportTaskToProto(t)))
}

// DeleteExportTask removes an export task by its identifier.
func (s *NeptuneGraphStore) DeleteExportTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.exportTasks.Exists(id) {
		return ErrExportTaskNotFound
	}
	return NewStoreError("delete_export_task", s.exportTasks.Delete(id))
}

// ListExportTasks returns a paginated list of export tasks, optionally filtered by graph identifier.
func (s *NeptuneGraphStore) ListExportTasks(opts common.ListOptions, graphId string) ([]*ExportTask, string, bool, error) {
	if graphId != "" {
		var ids []string
		err := s.exportTasksByGraph.ScanPrefix(graphId+"/", func(key string, value []byte) error {
			parts := strings.SplitN(key, "/", 2)
			if len(parts) == 2 {
				ids = append(ids, parts[1])
			}
			return nil
		})
		if err != nil {
			return nil, "", false, NewStoreError("list_export_tasks_idx", err)
		}

		startIdx := 0
		if opts.Marker != "" {
			for i, id := range ids {
				if id == opts.Marker {
					startIdx = i + 1
					break
				}
			}
		}

		if startIdx >= len(ids) {
			return nil, "", false, nil
		}

		ids = ids[startIdx:]
		if opts.MaxItems > 0 && len(ids) > opts.MaxItems {
			ids = ids[:opts.MaxItems]
		}

		tasks := make([]*ExportTask, 0, len(ids))
		for _, id := range ids {
			task, err := s.GetExportTask(id)
			if err != nil {
				continue
			}
			tasks = append(tasks, task)
		}

		truncated := false
		nextToken := ""
		if opts.MaxItems > 0 && len(tasks) == opts.MaxItems {
			nextToken = tasks[len(tasks)-1].TaskId
			truncated = true
		}
		return tasks, nextToken, truncated, nil
	}
	result, err := common.ListProto[*pb.ExportTaskRecord](s.exportTasks, opts, func() *pb.ExportTaskRecord { return &pb.ExportTaskRecord{} }, nil)
	if err != nil {
		return nil, "", false, NewStoreError("list_export_tasks", err)
	}
	tasks := make([]*ExportTask, len(result.Items))
	for i, item := range result.Items {
		tasks[i] = protoToExportTask(item)
	}
	return tasks, result.NextMarker, result.IsTruncated, nil
}

// TryAdvanceImportTask atomically advances an import task status if it matches the expected current status.
func (s *NeptuneGraphStore) TryAdvanceImportTask(id string, expectedStatus string, apply func(*ImportTask)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var p pb.ImportTaskRecord
	if err := s.importTasks.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return ErrImportTaskNotFound
		}
		return NewStoreError("advance_import_task", err)
	}
	if p.GetStatus() != expectedStatus {
		return nil
	}
	task := protoToImportTask(&p)
	if apply != nil {
		apply(task)
	}
	return NewStoreError("advance_import_task", s.importTasks.PutProto(id, importTaskToProto(task)))
}

// TryAdvanceExportTask atomically advances an export task status if it matches the expected current status.
func (s *NeptuneGraphStore) TryAdvanceExportTask(id string, expectedStatus string, apply func(*ExportTask)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var p pb.ExportTaskRecord
	if err := s.exportTasks.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return ErrExportTaskNotFound
		}
		return NewStoreError("advance_export_task", err)
	}
	if p.GetStatus() != expectedStatus {
		return nil
	}
	task := protoToExportTask(&p)
	if apply != nil {
		apply(task)
	}
	return NewStoreError("advance_export_task", s.exportTasks.PutProto(id, exportTaskToProto(task)))
}
