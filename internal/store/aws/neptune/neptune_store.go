package neptune

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_neptune"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	clusterBucket           = "neptune_clusters"
	instanceBucket          = "neptune_instances"
	snapshotBucket          = "neptune_snapshots"
	clusterParamGroupBucket = "neptune_cluster_param_groups"
	paramGroupBucket        = "neptune_param_groups"
	subnetGroupBucket       = "neptune_subnet_groups"
	globalClusterBucket     = "neptune_global_clusters"
	eventSubBucket          = "neptune_event_subscriptions"
	eventsBucket            = "neptune_events"
	tagsBucket              = "neptune_tags"
	clusterEndpointBucket   = "neptune_cluster_endpoints"
	queryStateBucket        = "neptunedata_queries"
	loaderJobBucket         = "neptunedata_loader_jobs"
)

// NeptuneStore provides the Pebble-backed storage implementation for all Neptune
// Management API resources. Each resource type is stored in its own bucket with
// protobuf serialisation. A single RWMutex guards all write operations.
type NeptuneStore struct {
	clusters           *common.ProtoStore[DBCluster]
	instances          *common.ProtoStore[DBInstance]
	snapshots          *common.ProtoStore[DBClusterSnapshot]
	clusterParamGroups *common.ProtoStore[DBClusterParameterGroup]
	paramGroups        *common.ProtoStore[DBParameterGroup]
	subnetGroups       *common.ProtoStore[DBSubnetGroup]
	globalClusters     *common.ProtoStore[GlobalCluster]
	eventSubs          *common.ProtoStore[EventSubscription]
	clusterEndpoints   *common.ProtoStore[DBClusterEndpoint]
	queries            *common.RawProtoStore[*pb.QueryState]
	loaderJobs         *common.RawProtoStore[*pb.LoaderJob]
	events             *common.BaseStore
	tags               *common.BaseStore
	mu                 sync.RWMutex
}

// NewNeptuneStore creates a new NeptuneStore backed by the given storage instance.
// Each bucket is allocated a separate BaseStore for the neptune service namespace.
func NewNeptuneStore(store storage.BasicStorage) *NeptuneStore {
	return &NeptuneStore{
		clusters: common.NewProtoStore(common.ProtoStoreConfig[DBCluster]{
			Store:        common.NewBaseStore(store.Bucket(clusterBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.DBCluster{} },
			ToDomain:     func(m proto.Message) *DBCluster { return ProtoToCluster(m.(*pb.DBCluster)) },
			ToProto:      func(d *DBCluster) proto.Message { return ClusterToProto(d) },
			IDFunc:       func(d *DBCluster) string { return d.DBClusterIdentifier },
			NotFoundErr:  ErrDBClusterNotFound,
			AlreadyExist: ErrDBClusterAlreadyExists,
		}),
		instances: common.NewProtoStore(common.ProtoStoreConfig[DBInstance]{
			Store:        common.NewBaseStore(store.Bucket(instanceBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.DBInstance{} },
			ToDomain:     func(m proto.Message) *DBInstance { return ProtoToInstance(m.(*pb.DBInstance)) },
			ToProto:      func(d *DBInstance) proto.Message { return InstanceToProto(d) },
			IDFunc:       func(d *DBInstance) string { return d.DBInstanceIdentifier },
			NotFoundErr:  ErrDBInstanceNotFound,
			AlreadyExist: ErrDBInstanceAlreadyExists,
		}),
		snapshots: common.NewProtoStore(common.ProtoStoreConfig[DBClusterSnapshot]{
			Store:        common.NewBaseStore(store.Bucket(snapshotBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.DBClusterSnapshot{} },
			ToDomain:     func(m proto.Message) *DBClusterSnapshot { return ProtoToSnapshot(m.(*pb.DBClusterSnapshot)) },
			ToProto:      func(d *DBClusterSnapshot) proto.Message { return SnapshotToProto(d) },
			IDFunc:       func(d *DBClusterSnapshot) string { return d.DBClusterSnapshotIdentifier },
			NotFoundErr:  ErrDBClusterSnapshotNotFound,
			AlreadyExist: ErrDBClusterSnapshotAlreadyExists,
		}),
		clusterParamGroups: common.NewProtoStore(common.ProtoStoreConfig[DBClusterParameterGroup]{
			Store:    common.NewBaseStore(store.Bucket(clusterParamGroupBucket), "neptune"),
			NewProto: func() proto.Message { return &pb.DBClusterParameterGroup{} },
			ToDomain: func(m proto.Message) *DBClusterParameterGroup {
				return ProtoToClusterParameterGroup(m.(*pb.DBClusterParameterGroup))
			},
			ToProto:      func(d *DBClusterParameterGroup) proto.Message { return ClusterParameterGroupToProto(d) },
			IDFunc:       func(d *DBClusterParameterGroup) string { return d.DBClusterParameterGroupName },
			NotFoundErr:  ErrDBClusterParameterGroupNotFound,
			AlreadyExist: ErrDBClusterParameterGroupAlreadyExists,
		}),
		paramGroups: common.NewProtoStore(common.ProtoStoreConfig[DBParameterGroup]{
			Store:        common.NewBaseStore(store.Bucket(paramGroupBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.DBParameterGroup{} },
			ToDomain:     func(m proto.Message) *DBParameterGroup { return ProtoToParameterGroup(m.(*pb.DBParameterGroup)) },
			ToProto:      func(d *DBParameterGroup) proto.Message { return ParameterGroupToProto(d) },
			IDFunc:       func(d *DBParameterGroup) string { return d.DBParameterGroupName },
			NotFoundErr:  ErrDBParameterGroupNotFound,
			AlreadyExist: ErrDBParameterGroupAlreadyExists,
		}),
		subnetGroups: common.NewProtoStore(common.ProtoStoreConfig[DBSubnetGroup]{
			Store:        common.NewBaseStore(store.Bucket(subnetGroupBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.DBSubnetGroup{} },
			ToDomain:     func(m proto.Message) *DBSubnetGroup { return ProtoToSubnetGroup(m.(*pb.DBSubnetGroup)) },
			ToProto:      func(d *DBSubnetGroup) proto.Message { return SubnetGroupToProto(d) },
			IDFunc:       func(d *DBSubnetGroup) string { return d.DBSubnetGroupName },
			NotFoundErr:  ErrDBSubnetGroupNotFound,
			AlreadyExist: ErrDBSubnetGroupAlreadyExists,
		}),
		globalClusters: common.NewProtoStore(common.ProtoStoreConfig[GlobalCluster]{
			Store:        common.NewBaseStore(store.Bucket(globalClusterBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.GlobalCluster{} },
			ToDomain:     func(m proto.Message) *GlobalCluster { return ProtoToGlobalCluster(m.(*pb.GlobalCluster)) },
			ToProto:      func(d *GlobalCluster) proto.Message { return GlobalClusterToProto(d) },
			IDFunc:       func(d *GlobalCluster) string { return d.GlobalClusterIdentifier },
			NotFoundErr:  ErrGlobalClusterNotFound,
			AlreadyExist: ErrGlobalClusterAlreadyExists,
		}),
		eventSubs: common.NewProtoStore(common.ProtoStoreConfig[EventSubscription]{
			Store:        common.NewBaseStore(store.Bucket(eventSubBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.EventSubscription{} },
			ToDomain:     func(m proto.Message) *EventSubscription { return ProtoToEventSubscription(m.(*pb.EventSubscription)) },
			ToProto:      func(d *EventSubscription) proto.Message { return EventSubscriptionToProto(d) },
			IDFunc:       func(d *EventSubscription) string { return d.CustSubscriptionId },
			NotFoundErr:  ErrEventSubscriptionNotFound,
			AlreadyExist: ErrEventSubscriptionAlreadyExists,
		}),
		clusterEndpoints: common.NewProtoStore(common.ProtoStoreConfig[DBClusterEndpoint]{
			Store:        common.NewBaseStore(store.Bucket(clusterEndpointBucket), "neptune"),
			NewProto:     func() proto.Message { return &pb.DBClusterEndpoint{} },
			ToDomain:     func(m proto.Message) *DBClusterEndpoint { return ProtoToClusterEndpoint(m.(*pb.DBClusterEndpoint)) },
			ToProto:      func(d *DBClusterEndpoint) proto.Message { return ClusterEndpointToProto(d) },
			IDFunc:       func(d *DBClusterEndpoint) string { return d.DBClusterEndpointIdentifier },
			NotFoundErr:  ErrDBClusterEndpointNotFound,
			AlreadyExist: ErrDBClusterEndpointAlreadyExists,
		}),
		queries: common.NewRawProtoStore(common.RawProtoStoreConfig[*pb.QueryState]{
			Store:    common.NewBaseStore(store.Bucket(queryStateBucket), "neptunedata"),
			NewProto: func() *pb.QueryState { return &pb.QueryState{} },
			IDFunc:   func(p *pb.QueryState) string { return p.GetQueryId() },
		}),
		loaderJobs: common.NewRawProtoStore(common.RawProtoStoreConfig[*pb.LoaderJob]{
			Store:    common.NewBaseStore(store.Bucket(loaderJobBucket), "neptunedata"),
			NewProto: func() *pb.LoaderJob { return &pb.LoaderJob{} },
			IDFunc:   func(p *pb.LoaderJob) string { return p.GetLoadId() },
		}),
		events: common.NewBaseStore(store.Bucket(eventsBucket), "neptune"),
		tags:   common.NewBaseStore(store.Bucket(tagsBucket), "neptune"),
	}
}

// CreateCluster persists a new DB cluster to the store.
func (s *NeptuneStore) CreateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Create(cluster)
}

// GetCluster retrieves a DB cluster by its identifier. Returns
// ErrDBClusterNotFound if the cluster does not exist.
func (s *NeptuneStore) GetCluster(id string) (*DBCluster, error) {
	return s.clusters.Get(id)
}

// UpdateCluster replaces the persisted state of an existing DB cluster.
func (s *NeptuneStore) UpdateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Update(cluster)
}

// DeleteCluster removes a DB cluster from the store.
func (s *NeptuneStore) DeleteCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Delete(id)
}

// ListClusters returns all DB clusters in the store.
func (s *NeptuneStore) ListClusters() ([]*DBCluster, error) {
	return s.clusters.List()
}

// CreateInstance persists a new DB instance to the store.
func (s *NeptuneStore) CreateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Create(instance)
}

// GetInstance retrieves a DB instance by its identifier. Returns
// ErrDBInstanceNotFound if the instance does not exist.
func (s *NeptuneStore) GetInstance(id string) (*DBInstance, error) {
	return s.instances.Get(id)
}

// UpdateInstance replaces the persisted state of an existing DB instance.
func (s *NeptuneStore) UpdateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Update(instance)
}

// DeleteInstance removes a DB instance from the store.
func (s *NeptuneStore) DeleteInstance(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Delete(id)
}

// ListInstances returns all DB instances in the store.
func (s *NeptuneStore) ListInstances() ([]*DBInstance, error) {
	return s.instances.List()
}

// CreateSnapshot persists a new cluster snapshot to the store.
func (s *NeptuneStore) CreateSnapshot(snapshot *DBClusterSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots.Create(snapshot)
}

// GetSnapshot retrieves a cluster snapshot by its identifier.
func (s *NeptuneStore) GetSnapshot(id string) (*DBClusterSnapshot, error) {
	return s.snapshots.Get(id)
}

// DeleteSnapshot removes a cluster snapshot from the store.
func (s *NeptuneStore) DeleteSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots.Delete(id)
}

// UpdateSnapshot persists modifications to an existing cluster snapshot.
func (s *NeptuneStore) UpdateSnapshot(snapshot *DBClusterSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots.Update(snapshot)
}

// ListSnapshots returns all cluster snapshots in the store.
func (s *NeptuneStore) ListSnapshots() ([]*DBClusterSnapshot, error) {
	return s.snapshots.List()
}

// CreateClusterParameterGroup persists a new cluster parameter group.
func (s *NeptuneStore) CreateClusterParameterGroup(pg *DBClusterParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterParamGroups.Create(pg)
}

// GetClusterParameterGroup retrieves a cluster parameter group by name.
func (s *NeptuneStore) GetClusterParameterGroup(name string) (*DBClusterParameterGroup, error) {
	return s.clusterParamGroups.Get(name)
}

// DeleteClusterParameterGroup removes a cluster parameter group.
func (s *NeptuneStore) DeleteClusterParameterGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterParamGroups.Delete(name)
}

// UpdateClusterParameterGroup replaces the persisted state of a cluster parameter group.
func (s *NeptuneStore) UpdateClusterParameterGroup(pg *DBClusterParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterParamGroups.Update(pg)
}

// ListClusterParameterGroups returns all cluster parameter groups.
func (s *NeptuneStore) ListClusterParameterGroups() ([]*DBClusterParameterGroup, error) {
	return s.clusterParamGroups.List()
}

// CreateParameterGroup persists a new DB parameter group.
func (s *NeptuneStore) CreateParameterGroup(pg *DBParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paramGroups.Create(pg)
}

// GetParameterGroup retrieves a DB parameter group by name.
func (s *NeptuneStore) GetParameterGroup(name string) (*DBParameterGroup, error) {
	return s.paramGroups.Get(name)
}

// DeleteParameterGroup removes a DB parameter group.
func (s *NeptuneStore) DeleteParameterGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paramGroups.Delete(name)
}

// UpdateParameterGroup replaces the persisted state of a DB parameter group.
func (s *NeptuneStore) UpdateParameterGroup(pg *DBParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paramGroups.Update(pg)
}

// ListParameterGroups returns all DB parameter groups.
func (s *NeptuneStore) ListParameterGroups() ([]*DBParameterGroup, error) {
	return s.paramGroups.List()
}

// CreateSubnetGroup persists a new DB subnet group.
func (s *NeptuneStore) CreateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Create(sg)
}

// GetSubnetGroup retrieves a DB subnet group by name.
func (s *NeptuneStore) GetSubnetGroup(name string) (*DBSubnetGroup, error) {
	return s.subnetGroups.Get(name)
}

// UpdateSubnetGroup replaces the persisted state of an existing subnet group.
func (s *NeptuneStore) UpdateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Update(sg)
}

// DeleteSubnetGroup removes a DB subnet group.
func (s *NeptuneStore) DeleteSubnetGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Delete(name)
}

// ListSubnetGroups returns all DB subnet groups.
func (s *NeptuneStore) ListSubnetGroups() ([]*DBSubnetGroup, error) {
	return s.subnetGroups.List()
}

// CreateGlobalCluster persists a new global cluster.
func (s *NeptuneStore) CreateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Create(gc)
}

// GetGlobalCluster retrieves a global cluster by its identifier.
func (s *NeptuneStore) GetGlobalCluster(id string) (*GlobalCluster, error) {
	return s.globalClusters.Get(id)
}

// UpdateGlobalCluster replaces the persisted state of an existing global cluster.
func (s *NeptuneStore) UpdateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Update(gc)
}

// DeleteGlobalCluster removes a global cluster from the store.
func (s *NeptuneStore) DeleteGlobalCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Delete(id)
}

// ListGlobalClusters returns all global clusters in the store.
func (s *NeptuneStore) ListGlobalClusters() ([]*GlobalCluster, error) {
	return s.globalClusters.List()
}

// CreateEventSubscription persists a new event subscription.
func (s *NeptuneStore) CreateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Create(sub)
}

// GetEventSubscription retrieves an event subscription by its subscription ID.
func (s *NeptuneStore) GetEventSubscription(name string) (*EventSubscription, error) {
	return s.eventSubs.Get(name)
}

// UpdateEventSubscription replaces the persisted state of an event subscription.
func (s *NeptuneStore) UpdateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Update(sub)
}

// DeleteEventSubscription removes an event subscription.
func (s *NeptuneStore) DeleteEventSubscription(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Delete(name)
}

// ListEventSubscriptions returns all event subscriptions.
func (s *NeptuneStore) ListEventSubscriptions() ([]*EventSubscription, error) {
	return s.eventSubs.List()
}

// AddTags merges the given tags into the existing tags for a resource. Tags
// with the same key are overwritten.
func (s *NeptuneStore) AddTags(resourceArn string, tags []types.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.getTagsUnlocked(resourceArn)
	if err != nil {
		return err
	}
	tagMap := make(map[string]string)
	for _, t := range existing {
		tagMap[t.Key] = t.Value
	}
	for _, t := range tags {
		tagMap[t.Key] = t.Value
	}
	merged := make([]types.Tag, 0, len(tagMap))
	for k, v := range tagMap {
		merged = append(merged, types.Tag{Key: k, Value: v})
	}
	return s.tags.PutProto(resourceArn, &pb.TagList{Tags: TagsToProto(merged)})
}

// GetTags returns the tags associated with a resource ARN.
func (s *NeptuneStore) GetTags(resourceArn string) ([]types.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getTagsUnlocked(resourceArn)
}

// getTagsUnlocked retrieves tags without acquiring the lock. Caller must hold
// at least a read lock.
func (s *NeptuneStore) getTagsUnlocked(resourceArn string) ([]types.Tag, error) {
	var pbTags pb.TagList
	if err := s.tags.GetProto(resourceArn, &pbTags); err != nil {
		if common.IsNotFound(err) {
			return []types.Tag{}, nil
		}
		return nil, err
	}
	return ProtoToTags(pbTags.Tags), nil
}

// RemoveTags removes the tags with the specified keys from a resource.
func (s *NeptuneStore) RemoveTags(resourceArn string, keys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.getTagsUnlocked(resourceArn)
	if err != nil {
		return err
	}
	removeSet := make(map[string]bool, len(keys))
	for _, k := range keys {
		removeSet[k] = true
	}
	filtered := make([]types.Tag, 0, len(existing))
	for _, t := range existing {
		if !removeSet[t.Key] {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == 0 {
		return s.tags.Delete(resourceArn)
	}
	return s.tags.PutProto(resourceArn, &pb.TagList{Tags: TagsToProto(filtered)})
}

// CreateClusterEndpoint persists a new custom cluster endpoint.
func (s *NeptuneStore) CreateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterEndpoints.Create(ep)
}

// GetClusterEndpoint retrieves a cluster endpoint by its identifier.
func (s *NeptuneStore) GetClusterEndpoint(id string) (*DBClusterEndpoint, error) {
	return s.clusterEndpoints.Get(id)
}

// UpdateClusterEndpoint replaces the persisted state of a cluster endpoint.
func (s *NeptuneStore) UpdateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterEndpoints.Update(ep)
}

// DeleteClusterEndpoint removes a cluster endpoint from the store.
func (s *NeptuneStore) DeleteClusterEndpoint(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterEndpoints.DeleteIfExists(id)
}

// ListClusterEndpoints returns all cluster endpoints, optionally filtered by
// the parent cluster ID. An empty clusterID returns all endpoints.
func (s *NeptuneStore) ListClusterEndpoints(clusterID string) ([]*DBClusterEndpoint, error) {
	return s.clusterEndpoints.ListFiltered(func(ep *DBClusterEndpoint) bool {
		return clusterID == "" || ep.DBClusterIdentifier == clusterID
	})
}

// CreateQuery persists a new Neptune Data API query state.
func (s *NeptuneStore) CreateQuery(q *pb.QueryState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queries.Create(q)
}

// GetQuery retrieves a Neptune Data API query state by ID.
func (s *NeptuneStore) GetQuery(id string) (*pb.QueryState, error) {
	return s.queries.Get(id)
}

// UpdateQuery replaces the persisted state of a Neptune Data API query.
func (s *NeptuneStore) UpdateQuery(q *pb.QueryState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queries.Update(q)
}

// DeleteQuery removes a Neptune Data API query state.
func (s *NeptuneStore) DeleteQuery(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queries.Delete(id)
}

// ListQueries returns all Neptune Data API query states.
func (s *NeptuneStore) ListQueries() ([]*pb.QueryState, error) {
	return s.queries.List()
}

// CreateLoaderJob persists a new Neptune Data API bulk loader job.
func (s *NeptuneStore) CreateLoaderJob(job *pb.LoaderJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loaderJobs.Create(job)
}

// GetLoaderJob retrieves a Neptune Data API bulk loader job by ID.
func (s *NeptuneStore) GetLoaderJob(id string) (*pb.LoaderJob, error) {
	return s.loaderJobs.Get(id)
}

// UpdateLoaderJob replaces the persisted state of a Neptune Data API bulk loader job.
func (s *NeptuneStore) UpdateLoaderJob(job *pb.LoaderJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loaderJobs.Update(job)
}

// DeleteLoaderJob removes a Neptune Data API bulk loader job.
func (s *NeptuneStore) DeleteLoaderJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loaderJobs.Delete(id)
}

// ListLoaderJobs returns all Neptune Data API bulk loader jobs.
func (s *NeptuneStore) ListLoaderJobs() ([]*pb.LoaderJob, error) {
	return s.loaderJobs.List()
}

const (
	maxEventAge = 7 * 24 * time.Hour
	maxEvents   = 10000
)

// RecordEvent persists a Neptune database event to the store.
func (s *NeptuneStore) RecordEvent(evt *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if evt.EventID == "" {
		evt.EventID = fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	return s.events.Put(evt.EventID, evt)
}

// ListEvents returns events matching the given filters with pagination.
func (s *NeptuneStore) ListEvents(opts EventListOptions) (*EventListResult, error) {
	if opts.MaxRecords <= 0 {
		opts.MaxRecords = 100
	}

	cutoff := time.Now().Add(-maxEventAge)

	var allEvents []*Event
	err := s.events.ForEach(func(key string, value []byte) error {
		var evt Event
		if err := json.Unmarshal(value, &evt); err != nil {
			return err
		}
		if evt.Date.Before(cutoff) {
			return nil
		}
		if opts.SourceType != "" && evt.SourceType != opts.SourceType {
			return nil
		}
		if opts.SourceIdentifier != "" && evt.SourceIdentifier != opts.SourceIdentifier {
			return nil
		}
		if !opts.StartTime.IsZero() && evt.Date.Before(opts.StartTime) {
			return nil
		}
		if !opts.EndTime.IsZero() && evt.Date.After(opts.EndTime) {
			return nil
		}
		allEvents = append(allEvents, &evt)
		return nil
	})
	if err != nil {
		return nil, err
	}

	started := opts.Marker == ""
	var page []*Event
	for _, evt := range allEvents {
		if !started {
			if evt.EventID == opts.Marker {
				started = true
			}
			continue
		}
		page = append(page, evt)
		if len(page) >= opts.MaxRecords {
			break
		}
	}

	result := &EventListResult{Events: page}
	if len(page) < len(allEvents) {
		result.IsTruncated = true
		result.Marker = page[len(page)-1].EventID
	}
	return result, nil
}

// PurgeOldEvents removes events older than maxEventAge and trims the total
// count to maxEvents by deleting the oldest entries first.
func (s *NeptuneStore) PurgeOldEvents() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-maxEventAge)

	var allEvents []*Event
	err := s.events.ForEach(func(key string, value []byte) error {
		var evt Event
		if err := json.Unmarshal(value, &evt); err != nil {
			return err
		}
		allEvents = append(allEvents, &evt)
		return nil
	})
	if err != nil {
		return err
	}

	keep := make([]*Event, 0, len(allEvents))
	for _, evt := range allEvents {
		if evt.Date.Before(cutoff) {
			_ = s.events.Delete(evt.EventID)
			continue
		}
		keep = append(keep, evt)
	}

	if len(keep) > maxEvents {
		sort.Slice(keep, func(i, j int) bool { return keep[i].Date.Before(keep[j].Date) })
		for i := 0; i < len(keep)-maxEvents; i++ {
			_ = s.events.Delete(keep[i].EventID)
		}
	}
	return nil
}
