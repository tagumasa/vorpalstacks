package neptunegraph

import (
	"strings"
	"sync"

	pb "vorpalstacks/internal/pb/storage/storage_neptunegraph"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
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

// NeptuneGraphStore provides persistent storage for NeptuneGraph resources using underlying key-value buckets.
type NeptuneGraphStore struct {
	graphs             *common.ProtoStore[Graph]
	snapshots          *common.BaseStore
	snapshotsByGraph   *common.BaseStore
	endpoints          *common.BaseStore
	tags               *common.BaseStore
	queries            *common.BaseStore
	importTasks        *common.ProtoStore[ImportTask]
	exportTasks        *common.ProtoStore[ExportTask]
	exportTasksByGraph *common.BaseStore
	mu                 sync.RWMutex
}

// NewNeptuneGraphStore creates a new NeptuneGraphStore backed by the given storage.
func NewNeptuneGraphStore(store storage.BasicStorage) *NeptuneGraphStore {
	return &NeptuneGraphStore{
		graphs: common.NewProtoStore(common.ProtoStoreConfig[Graph]{
			Store:        common.NewBaseStore(store.Bucket(graphsBucket), "neptunegraph"),
			NewProto:     func() proto.Message { return &pb.Graph{} },
			ToDomain:     func(m proto.Message) *Graph { return protoToGraph(m.(*pb.Graph)) },
			ToProto:      func(d *Graph) proto.Message { return graphToProto(d) },
			IDFunc:       func(d *Graph) string { return d.Id },
			NotFoundErr:  ErrGraphNotFound,
			AlreadyExist: ErrGraphAlreadyExists,
		}),
		snapshots:        common.NewBaseStore(store.Bucket(snapshotsBucket), "neptunegraph"),
		snapshotsByGraph: common.NewBaseStore(store.Bucket(snapshotsByGraphBucket), "neptunegraph"),
		endpoints:        common.NewBaseStore(store.Bucket(endpointsBucket), "neptunegraph"),
		tags:             common.NewBaseStore(store.Bucket(tagsBucket), "neptunegraph"),
		queries:          common.NewBaseStore(store.Bucket(queriesBucket), "neptunegraph"),
		importTasks: common.NewProtoStore(common.ProtoStoreConfig[ImportTask]{
			Store:        common.NewBaseStore(store.Bucket(importTasksBucket), "neptunegraph"),
			NewProto:     func() proto.Message { return &pb.ImportTaskRecord{} },
			ToDomain:     func(m proto.Message) *ImportTask { return protoToImportTask(m.(*pb.ImportTaskRecord)) },
			ToProto:      func(d *ImportTask) proto.Message { return importTaskToProto(d) },
			IDFunc:       func(d *ImportTask) string { return d.TaskId },
			NotFoundErr:  ErrImportTaskNotFound,
			AlreadyExist: ErrImportTaskAlreadyExists,
		}),
		exportTasks: common.NewProtoStore(common.ProtoStoreConfig[ExportTask]{
			Store:        common.NewBaseStore(store.Bucket(exportTasksBucket), "neptunegraph"),
			NewProto:     func() proto.Message { return &pb.ExportTaskRecord{} },
			ToDomain:     func(m proto.Message) *ExportTask { return protoToExportTask(m.(*pb.ExportTaskRecord)) },
			ToProto:      func(d *ExportTask) proto.Message { return exportTaskToProto(d) },
			IDFunc:       func(d *ExportTask) string { return d.TaskId },
			NotFoundErr:  ErrExportTaskNotFound,
			AlreadyExist: ErrExportTaskAlreadyExists,
		}),
		exportTasksByGraph: common.NewBaseStore(store.Bucket(exportTasksByGraphBucket), "neptunegraph"),
	}
}

// CreateGraph persists a new graph, returning an error if the graph already exists.
func (s *NeptuneGraphStore) CreateGraph(graph *Graph) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.graphs.Create(graph)
}

// GetGraph retrieves a graph by its identifier.
func (s *NeptuneGraphStore) GetGraph(id string) (*Graph, error) {
	return s.graphs.Get(id)
}

// UpdateGraph persists updated graph properties.
func (s *NeptuneGraphStore) UpdateGraph(graph *Graph) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.graphs.Update(graph)
}

// DeleteGraph removes a graph and all associated queries, snapshots, and export task indices.
func (s *NeptuneGraphStore) DeleteGraph(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.graphs.Exists(id) {
		return ErrGraphNotFound
	}
	if err := s.graphs.Store().Delete(id); err != nil {
		return err
	}
	if err := s.queries.DeleteByPrefix(id + "/"); err != nil {
		logs.Warn("Failed to delete queries for graph", logs.String("graphId", id), logs.Err(err))
	}
	if err := s.snapshotsByGraph.DeleteByPrefix(id + "/"); err != nil {
		logs.Warn("Failed to delete snapshots for graph", logs.String("graphId", id), logs.Err(err))
	}
	if err := s.exportTasksByGraph.DeleteByPrefix(id + "/"); err != nil {
		logs.Warn("Failed to delete export tasks for graph", logs.String("graphId", id), logs.Err(err))
	}
	return nil
}

// ListGraphs returns a paginated list of graphs.
func (s *NeptuneGraphStore) ListGraphs(opts common.ListOptions) ([]*Graph, string, bool, error) {
	result, err := common.ListProto[*pb.Graph](s.graphs.Store(), opts, func() *pb.Graph { return &pb.Graph{} }, nil)
	if err != nil {
		return nil, "", false, err
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
		return err
	}
	idxKey := snapshot.SourceGraphId + "/" + snapshot.Id
	return s.snapshotsByGraph.Put(idxKey, "")
}

// GetSnapshot retrieves a snapshot by its identifier.
func (s *NeptuneGraphStore) GetSnapshot(id string) (*GraphSnapshot, error) {
	var p pb.GraphSnapshot
	if err := s.snapshots.GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrSnapshotNotFound
		}
		return nil, err
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
	return s.snapshots.PutProto(snapshot.Id, snapshotToProto(snapshot))
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
		return err
	}
	idxKey := p.GetSourceGraphId() + "/" + id
	_ = s.snapshotsByGraph.Delete(idxKey)
	return s.snapshots.Delete(id)
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
			return nil, "", false, err
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
		return nil, "", false, err
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
	return s.endpoints.PutProto(key, endpointToProto(ep))
}

// GetEndpoint retrieves a private graph endpoint by graph and VPC identifiers.
func (s *NeptuneGraphStore) GetEndpoint(graphId, vpcId string) (*PrivateGraphEndpoint, error) {
	key := graphId + ":" + vpcId
	var p pb.PrivateGraphEndpoint
	if err := s.endpoints.GetProto(key, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEndpointNotFound
		}
		return nil, err
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
	return s.endpoints.Delete(key)
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
		return nil, err
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
	return s.tags.PutProto(resourceArn, tagsToProto(existing))
}

// GetTags retrieves all tags associated with a resource ARN.
func (s *NeptuneGraphStore) GetTags(resourceArn string) (map[string]string, error) {
	var tp pb.TagMap
	if err := s.tags.GetProto(resourceArn, &tp); err != nil {
		if common.IsNotFound(err) {
			return map[string]string{}, nil
		}
		return nil, err
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
	return s.tags.PutProto(resourceArn, tagsToProto(existing))
}

// CreateQuery persists a new query record.
func (s *NeptuneGraphStore) CreateQuery(q *QueryRecord) error {
	key := q.GraphId + "/" + q.Id
	return s.queries.PutProto(key, queryToProto(q))
}

// GetQuery retrieves a query record by graph and query identifiers.
func (s *NeptuneGraphStore) GetQuery(graphId, id string) (*QueryRecord, error) {
	var p pb.QueryRecord
	if err := s.queries.GetProto(graphId+"/"+id, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrQueryNotFound
		}
		return nil, err
	}
	return protoToQuery(&p), nil
}

// UpdateQuery persists updated query record fields.
func (s *NeptuneGraphStore) UpdateQuery(q *QueryRecord) error {
	key := q.GraphId + "/" + q.Id
	return s.queries.PutProto(key, queryToProto(q))
}

// DeleteQuery removes a query record by graph and query identifiers.
func (s *NeptuneGraphStore) DeleteQuery(graphId, id string) error {
	key := graphId + "/" + id
	if !s.queries.Exists(key) {
		return ErrQueryNotFound
	}
	return s.queries.Delete(key)
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
		return nil, err
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
	return s.importTasks.Create(t)
}

// GetImportTask retrieves an import task by its identifier.
func (s *NeptuneGraphStore) GetImportTask(id string) (*ImportTask, error) {
	return s.importTasks.Get(id)
}

// UpdateImportTask persists updated import task properties.
func (s *NeptuneGraphStore) UpdateImportTask(t *ImportTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.importTasks.UpdateDirect(t)
}

// DeleteImportTask removes an import task by its identifier.
func (s *NeptuneGraphStore) DeleteImportTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.importTasks.DeleteIfExists(id)
}

// ListImportTasks returns a paginated list of import tasks.
func (s *NeptuneGraphStore) ListImportTasks(opts common.ListOptions) ([]*ImportTask, string, bool, error) {
	result, err := common.ListProto[*pb.ImportTaskRecord](s.importTasks.Store(), opts, func() *pb.ImportTaskRecord { return &pb.ImportTaskRecord{} }, nil)
	if err != nil {
		return nil, "", false, err
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
	if s.exportTasks.Store().Exists(t.TaskId) {
		return ErrExportTaskAlreadyExists
	}
	if err := s.exportTasks.Store().PutProto(t.TaskId, exportTaskToProto(t)); err != nil {
		return err
	}
	if t.GraphId != "" {
		idxKey := t.GraphId + "/" + t.TaskId
		return s.exportTasksByGraph.Put(idxKey, "")
	}
	return nil
}

// GetExportTask retrieves an export task by its identifier.
func (s *NeptuneGraphStore) GetExportTask(id string) (*ExportTask, error) {
	return s.exportTasks.Get(id)
}

// UpdateExportTask persists updated export task properties.
func (s *NeptuneGraphStore) UpdateExportTask(t *ExportTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.exportTasks.UpdateDirect(t)
}

// DeleteExportTask removes an export task by its identifier.
func (s *NeptuneGraphStore) DeleteExportTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.exportTasks.DeleteIfExists(id)
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
			return nil, "", false, err
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
	result, err := common.ListProto[*pb.ExportTaskRecord](s.exportTasks.Store(), opts, func() *pb.ExportTaskRecord { return &pb.ExportTaskRecord{} }, nil)
	if err != nil {
		return nil, "", false, err
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
	if err := s.importTasks.Store().GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return ErrImportTaskNotFound
		}
		return err
	}
	if p.GetStatus() != expectedStatus {
		return nil
	}
	task := protoToImportTask(&p)
	if apply != nil {
		apply(task)
	}
	return s.importTasks.Store().PutProto(id, importTaskToProto(task))
}

// TryAdvanceExportTask atomically advances an export task status if it matches the expected current status.
func (s *NeptuneGraphStore) TryAdvanceExportTask(id string, expectedStatus string, apply func(*ExportTask)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var p pb.ExportTaskRecord
	if err := s.exportTasks.Store().GetProto(id, &p); err != nil {
		if common.IsNotFound(err) {
			return ErrExportTaskNotFound
		}
		return err
	}
	if p.GetStatus() != expectedStatus {
		return nil
	}
	task := protoToExportTask(&p)
	if apply != nil {
		apply(task)
	}
	return s.exportTasks.Store().PutProto(id, exportTaskToProto(task))
}
