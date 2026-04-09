// Package neptune provides the Pebble-backed storage layer for AWS Neptune
// Management API resources. All data is serialised via protobuf for
// persistence in the underlying Pebble key-value store.
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

// NeptuneStoreInterface aggregates all Neptune resource CRUD operations into a
// single interface for use by the Neptune Management API service handlers.
type NeptuneStoreInterface interface {
	ClusterOps
	InstanceOps
	SnapshotOps
	ClusterParameterGroupOps
	ParameterGroupOps
	SubnetGroupOps
	GlobalClusterOps
	EventSubscriptionOps
	EventOps
	TagOps
	ClusterEndpointOps
}

// NeptuneDataStoreInterface aggregates Neptune Data API persistence operations
// for query states and bulk loader jobs.
type NeptuneDataStoreInterface interface {
	QueryOps
	LoaderJobOps
}

// QueryOps defines CRUD operations for Neptune Data API query states.
type QueryOps interface {
	CreateQuery(q *pb.QueryState) error
	GetQuery(id string) (*pb.QueryState, error)
	UpdateQuery(q *pb.QueryState) error
	DeleteQuery(id string) error
	ListQueries() ([]*pb.QueryState, error)
}

// LoaderJobOps defines CRUD operations for Neptune Data API bulk loader jobs.
type LoaderJobOps interface {
	CreateLoaderJob(job *pb.LoaderJob) error
	GetLoaderJob(id string) (*pb.LoaderJob, error)
	UpdateLoaderJob(job *pb.LoaderJob) error
	DeleteLoaderJob(id string) error
	ListLoaderJobs() ([]*pb.LoaderJob, error)
}

// ClusterOps defines CRUD operations for Neptune DB clusters.
type ClusterOps interface {
	CreateCluster(cluster *DBCluster) error
	GetCluster(id string) (*DBCluster, error)
	UpdateCluster(cluster *DBCluster) error
	DeleteCluster(id string) error
	ListClusters() ([]*DBCluster, error)
}

// InstanceOps defines CRUD operations for Neptune DB instances.
type InstanceOps interface {
	CreateInstance(instance *DBInstance) error
	GetInstance(id string) (*DBInstance, error)
	UpdateInstance(instance *DBInstance) error
	DeleteInstance(id string) error
	ListInstances() ([]*DBInstance, error)
}

// SnapshotOps defines CRUD operations for Neptune cluster snapshots.
type SnapshotOps interface {
	CreateSnapshot(snapshot *DBClusterSnapshot) error
	GetSnapshot(id string) (*DBClusterSnapshot, error)
	DeleteSnapshot(id string) error
	UpdateSnapshot(snapshot *DBClusterSnapshot) error
	ListSnapshots() ([]*DBClusterSnapshot, error)
}

// ClusterParameterGroupOps defines CRUD operations for cluster-level parameter groups.
type ClusterParameterGroupOps interface {
	CreateClusterParameterGroup(pg *DBClusterParameterGroup) error
	GetClusterParameterGroup(name string) (*DBClusterParameterGroup, error)
	UpdateClusterParameterGroup(pg *DBClusterParameterGroup) error
	DeleteClusterParameterGroup(name string) error
	ListClusterParameterGroups() ([]*DBClusterParameterGroup, error)
}

// ParameterGroupOps defines CRUD operations for instance-level parameter groups.
type ParameterGroupOps interface {
	CreateParameterGroup(pg *DBParameterGroup) error
	GetParameterGroup(name string) (*DBParameterGroup, error)
	UpdateParameterGroup(pg *DBParameterGroup) error
	DeleteParameterGroup(name string) error
	ListParameterGroups() ([]*DBParameterGroup, error)
}

// SubnetGroupOps defines CRUD operations for DB subnet groups.
type SubnetGroupOps interface {
	CreateSubnetGroup(sg *DBSubnetGroup) error
	GetSubnetGroup(name string) (*DBSubnetGroup, error)
	UpdateSubnetGroup(sg *DBSubnetGroup) error
	DeleteSubnetGroup(name string) error
	ListSubnetGroups() ([]*DBSubnetGroup, error)
}

// GlobalClusterOps defines CRUD operations for Neptune global clusters.
type GlobalClusterOps interface {
	CreateGlobalCluster(gc *GlobalCluster) error
	GetGlobalCluster(id string) (*GlobalCluster, error)
	UpdateGlobalCluster(gc *GlobalCluster) error
	DeleteGlobalCluster(id string) error
	ListGlobalClusters() ([]*GlobalCluster, error)
}

// EventSubscriptionOps defines CRUD operations for Neptune event subscriptions.
type EventSubscriptionOps interface {
	CreateEventSubscription(sub *EventSubscription) error
	GetEventSubscription(name string) (*EventSubscription, error)
	UpdateEventSubscription(sub *EventSubscription) error
	DeleteEventSubscription(name string) error
	ListEventSubscriptions() ([]*EventSubscription, error)
}

// TagOps defines operations for managing resource tags.
type TagOps interface {
	AddTags(resourceArn string, tags []Tag) error
	GetTags(resourceArn string) ([]Tag, error)
	RemoveTags(resourceArn string, keys []string) error
}

// ClusterEndpointOps defines CRUD operations for custom cluster endpoints.
type ClusterEndpointOps interface {
	CreateClusterEndpoint(ep *DBClusterEndpoint) error
	GetClusterEndpoint(id string) (*DBClusterEndpoint, error)
	UpdateClusterEndpoint(ep *DBClusterEndpoint) error
	DeleteClusterEndpoint(id string) error
	ListClusterEndpoints(clusterID string) ([]*DBClusterEndpoint, error)
}

// EventOps defines operations for recording and listing Neptune events.
type EventOps interface {
	RecordEvent(evt *Event) error
	ListEvents(opts EventListOptions) (*EventListResult, error)
	PurgeOldEvents() error
}

// NeptuneStore provides the Pebble-backed storage implementation for all Neptune
// Management API resources. Each resource type is stored in its own bucket with
// protobuf serialisation. A single RWMutex guards all write operations.
type NeptuneStore struct {
	clusters           *common.BaseStore
	instances          *common.BaseStore
	snapshots          *common.BaseStore
	clusterParamGroups *common.BaseStore
	paramGroups        *common.BaseStore
	subnetGroups       *common.BaseStore
	globalClusters     *common.BaseStore
	eventSubs          *common.BaseStore
	events             *common.BaseStore
	tags               *common.BaseStore
	clusterEndpoints   *common.BaseStore
	queries            *common.BaseStore
	loaderJobs         *common.BaseStore
	mu                 sync.RWMutex
}

var _ NeptuneStoreInterface = (*NeptuneStore)(nil)
var _ NeptuneDataStoreInterface = (*NeptuneStore)(nil)

// NewNeptuneStore creates a new NeptuneStore backed by the given storage instance.
// Each bucket is allocated a separate BaseStore for the neptune service namespace.
func NewNeptuneStore(store storage.BasicStorage) *NeptuneStore {
	return &NeptuneStore{
		clusters:           common.NewBaseStore(store.Bucket(clusterBucket), "neptune"),
		instances:          common.NewBaseStore(store.Bucket(instanceBucket), "neptune"),
		snapshots:          common.NewBaseStore(store.Bucket(snapshotBucket), "neptune"),
		clusterParamGroups: common.NewBaseStore(store.Bucket(clusterParamGroupBucket), "neptune"),
		paramGroups:        common.NewBaseStore(store.Bucket(paramGroupBucket), "neptune"),
		subnetGroups:       common.NewBaseStore(store.Bucket(subnetGroupBucket), "neptune"),
		globalClusters:     common.NewBaseStore(store.Bucket(globalClusterBucket), "neptune"),
		eventSubs:          common.NewBaseStore(store.Bucket(eventSubBucket), "neptune"),
		events:             common.NewBaseStore(store.Bucket(eventsBucket), "neptune"),
		tags:               common.NewBaseStore(store.Bucket(tagsBucket), "neptune"),
		clusterEndpoints:   common.NewBaseStore(store.Bucket(clusterEndpointBucket), "neptune"),
		queries:            common.NewBaseStore(store.Bucket(queryStateBucket), "neptunedata"),
		loaderJobs:         common.NewBaseStore(store.Bucket(loaderJobBucket), "neptunedata"),
	}
}

// CreateCluster persists a new DB cluster to the store.
func (s *NeptuneStore) CreateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clusters.Bucket().Has([]byte(cluster.DBClusterIdentifier)) {
		return ErrDBClusterAlreadyExists
	}
	data, err := proto.Marshal(ClusterToProto(cluster))
	if err != nil {
		return NewStoreError("create_cluster", err)
	}
	return NewStoreError("create_cluster", s.clusters.Bucket().Put([]byte(cluster.DBClusterIdentifier), data))
}

// GetCluster retrieves a DB cluster by its identifier. Returns
// ErrDBClusterNotFound if the cluster does not exist.
func (s *NeptuneStore) GetCluster(id string) (*DBCluster, error) {
	data, err := s.clusters.Bucket().Get([]byte(id))
	if err != nil {
		return nil, NewStoreError("get_cluster", err)
	}
	if data == nil {
		return nil, ErrDBClusterNotFound
	}
	var pbCluster pb.DBCluster
	if err := proto.Unmarshal(data, &pbCluster); err != nil {
		return nil, NewStoreError("get_cluster", err)
	}
	return ProtoToCluster(&pbCluster), nil
}

// UpdateCluster replaces the persisted state of an existing DB cluster.
func (s *NeptuneStore) UpdateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.clusters.Bucket().Has([]byte(cluster.DBClusterIdentifier)) {
		return ErrDBClusterNotFound
	}
	data, err := proto.Marshal(ClusterToProto(cluster))
	if err != nil {
		return NewStoreError("update_cluster", err)
	}
	return NewStoreError("update_cluster", s.clusters.Bucket().Put([]byte(cluster.DBClusterIdentifier), data))
}

// DeleteCluster removes a DB cluster from the store.
func (s *NeptuneStore) DeleteCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_cluster", s.clusters.Delete(id))
}

// ListClusters returns all DB clusters in the store.
func (s *NeptuneStore) ListClusters() ([]*DBCluster, error) {
	var clusters []*DBCluster
	err := s.clusters.ForEach(func(key string, value []byte) error {
		var pbCluster pb.DBCluster
		if err := proto.Unmarshal(value, &pbCluster); err != nil {
			return err
		}
		clusters = append(clusters, ProtoToCluster(&pbCluster))
		return nil
	})
	return clusters, err
}

// CreateInstance persists a new DB instance to the store.
func (s *NeptuneStore) CreateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.instances.Bucket().Has([]byte(instance.DBInstanceIdentifier)) {
		return ErrDBInstanceAlreadyExists
	}
	data, err := proto.Marshal(InstanceToProto(instance))
	if err != nil {
		return NewStoreError("create_instance", err)
	}
	return NewStoreError("create_instance", s.instances.Bucket().Put([]byte(instance.DBInstanceIdentifier), data))
}

// GetInstance retrieves a DB instance by its identifier. Returns
// ErrDBInstanceNotFound if the instance does not exist.
func (s *NeptuneStore) GetInstance(id string) (*DBInstance, error) {
	data, err := s.instances.Bucket().Get([]byte(id))
	if err != nil || data == nil {
		return nil, ErrDBInstanceNotFound
	}
	var pbInstance pb.DBInstance
	if err := proto.Unmarshal(data, &pbInstance); err != nil {
		return nil, NewStoreError("get_instance", err)
	}
	return ProtoToInstance(&pbInstance), nil
}

// UpdateInstance replaces the persisted state of an existing DB instance.
func (s *NeptuneStore) UpdateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.instances.Bucket().Has([]byte(instance.DBInstanceIdentifier)) {
		return ErrDBInstanceNotFound
	}
	data, err := proto.Marshal(InstanceToProto(instance))
	if err != nil {
		return NewStoreError("update_instance", err)
	}
	return NewStoreError("update_instance", s.instances.Bucket().Put([]byte(instance.DBInstanceIdentifier), data))
}

// DeleteInstance removes a DB instance from the store.
func (s *NeptuneStore) DeleteInstance(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_instance", s.instances.Delete(id))
}

// ListInstances returns all DB instances in the store.
func (s *NeptuneStore) ListInstances() ([]*DBInstance, error) {
	var instances []*DBInstance
	err := s.instances.ForEach(func(key string, value []byte) error {
		var pbInst pb.DBInstance
		if err := proto.Unmarshal(value, &pbInst); err != nil {
			return err
		}
		instances = append(instances, ProtoToInstance(&pbInst))
		return nil
	})
	return instances, err
}

// CreateSnapshot persists a new cluster snapshot to the store.
func (s *NeptuneStore) CreateSnapshot(snapshot *DBClusterSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.snapshots.Bucket().Has([]byte(snapshot.DBClusterSnapshotIdentifier)) {
		return ErrDBClusterSnapshotAlreadyExists
	}
	data, err := proto.Marshal(SnapshotToProto(snapshot))
	if err != nil {
		return NewStoreError("create_snapshot", err)
	}
	return NewStoreError("create_snapshot", s.snapshots.Bucket().Put([]byte(snapshot.DBClusterSnapshotIdentifier), data))
}

// GetSnapshot retrieves a cluster snapshot by its identifier.
func (s *NeptuneStore) GetSnapshot(id string) (*DBClusterSnapshot, error) {
	data, err := s.snapshots.Bucket().Get([]byte(id))
	if err != nil || data == nil {
		return nil, ErrDBClusterSnapshotNotFound
	}
	var pbSnap pb.DBClusterSnapshot
	if err := proto.Unmarshal(data, &pbSnap); err != nil {
		return nil, NewStoreError("get_snapshot", err)
	}
	return ProtoToSnapshot(&pbSnap), nil
}

// DeleteSnapshot removes a cluster snapshot from the store.
func (s *NeptuneStore) DeleteSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_snapshot", s.snapshots.Delete(id))
}

// UpdateSnapshot persists modifications to an existing cluster snapshot.
func (s *NeptuneStore) UpdateSnapshot(snapshot *DBClusterSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.snapshots.Bucket().Has([]byte(snapshot.DBClusterSnapshotIdentifier)) {
		return ErrDBClusterSnapshotNotFound
	}
	data, err := proto.Marshal(SnapshotToProto(snapshot))
	if err != nil {
		return NewStoreError("update_snapshot", err)
	}
	return NewStoreError("update_snapshot", s.snapshots.Bucket().Put([]byte(snapshot.DBClusterSnapshotIdentifier), data))
}

// ListSnapshots returns all cluster snapshots in the store.
func (s *NeptuneStore) ListSnapshots() ([]*DBClusterSnapshot, error) {
	var snapshots []*DBClusterSnapshot
	err := s.snapshots.ForEach(func(key string, value []byte) error {
		var pbSnap pb.DBClusterSnapshot
		if err := proto.Unmarshal(value, &pbSnap); err != nil {
			return err
		}
		snapshots = append(snapshots, ProtoToSnapshot(&pbSnap))
		return nil
	})
	return snapshots, err
}

// CreateClusterParameterGroup persists a new cluster parameter group.
func (s *NeptuneStore) CreateClusterParameterGroup(pg *DBClusterParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clusterParamGroups.Bucket().Has([]byte(pg.DBClusterParameterGroupName)) {
		return ErrDBClusterParameterGroupAlreadyExists
	}
	data, err := proto.Marshal(ClusterParameterGroupToProto(pg))
	if err != nil {
		return NewStoreError("create_cluster_param_group", err)
	}
	return NewStoreError("create_cluster_param_group", s.clusterParamGroups.Bucket().Put([]byte(pg.DBClusterParameterGroupName), data))
}

// GetClusterParameterGroup retrieves a cluster parameter group by name.
func (s *NeptuneStore) GetClusterParameterGroup(name string) (*DBClusterParameterGroup, error) {
	data, err := s.clusterParamGroups.Bucket().Get([]byte(name))
	if err != nil || data == nil {
		return nil, ErrDBClusterParameterGroupNotFound
	}
	var pbPG pb.DBClusterParameterGroup
	if err := proto.Unmarshal(data, &pbPG); err != nil {
		return nil, NewStoreError("get_cluster_param_group", err)
	}
	return ProtoToClusterParameterGroup(&pbPG), nil
}

// DeleteClusterParameterGroup removes a cluster parameter group.
func (s *NeptuneStore) DeleteClusterParameterGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_cluster_param_group", s.clusterParamGroups.Delete(name))
}

// UpdateClusterParameterGroup replaces the persisted state of a cluster parameter group.
func (s *NeptuneStore) UpdateClusterParameterGroup(pg *DBClusterParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.clusterParamGroups.Bucket().Has([]byte(pg.DBClusterParameterGroupName)) {
		return ErrDBClusterParameterGroupNotFound
	}
	data, err := proto.Marshal(ClusterParameterGroupToProto(pg))
	if err != nil {
		return NewStoreError("update_cluster_param_group", err)
	}
	return NewStoreError("update_cluster_param_group", s.clusterParamGroups.Bucket().Put([]byte(pg.DBClusterParameterGroupName), data))
}

// ListClusterParameterGroups returns all cluster parameter groups.
func (s *NeptuneStore) ListClusterParameterGroups() ([]*DBClusterParameterGroup, error) {
	var groups []*DBClusterParameterGroup
	err := s.clusterParamGroups.ForEach(func(key string, value []byte) error {
		var pbPG pb.DBClusterParameterGroup
		if err := proto.Unmarshal(value, &pbPG); err != nil {
			return err
		}
		groups = append(groups, ProtoToClusterParameterGroup(&pbPG))
		return nil
	})
	return groups, err
}

// CreateParameterGroup persists a new DB parameter group.
func (s *NeptuneStore) CreateParameterGroup(pg *DBParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.paramGroups.Bucket().Has([]byte(pg.DBParameterGroupName)) {
		return ErrDBParameterGroupAlreadyExists
	}
	data, err := proto.Marshal(ParameterGroupToProto(pg))
	if err != nil {
		return NewStoreError("create_param_group", err)
	}
	return NewStoreError("create_param_group", s.paramGroups.Bucket().Put([]byte(pg.DBParameterGroupName), data))
}

// GetParameterGroup retrieves a DB parameter group by name.
func (s *NeptuneStore) GetParameterGroup(name string) (*DBParameterGroup, error) {
	data, err := s.paramGroups.Bucket().Get([]byte(name))
	if err != nil || data == nil {
		return nil, ErrDBParameterGroupNotFound
	}
	var pbPG pb.DBParameterGroup
	if err := proto.Unmarshal(data, &pbPG); err != nil {
		return nil, NewStoreError("get_param_group", err)
	}
	return ProtoToParameterGroup(&pbPG), nil
}

// DeleteParameterGroup removes a DB parameter group.
func (s *NeptuneStore) DeleteParameterGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_param_group", s.paramGroups.Delete(name))
}

// UpdateParameterGroup replaces the persisted state of a DB parameter group.
func (s *NeptuneStore) UpdateParameterGroup(pg *DBParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.paramGroups.Bucket().Has([]byte(pg.DBParameterGroupName)) {
		return ErrDBParameterGroupNotFound
	}
	data, err := proto.Marshal(ParameterGroupToProto(pg))
	if err != nil {
		return NewStoreError("update_param_group", err)
	}
	return NewStoreError("update_param_group", s.paramGroups.Bucket().Put([]byte(pg.DBParameterGroupName), data))
}

// ListParameterGroups returns all DB parameter groups.
func (s *NeptuneStore) ListParameterGroups() ([]*DBParameterGroup, error) {
	var groups []*DBParameterGroup
	err := s.paramGroups.ForEach(func(key string, value []byte) error {
		var pbPG pb.DBParameterGroup
		if err := proto.Unmarshal(value, &pbPG); err != nil {
			return err
		}
		groups = append(groups, ProtoToParameterGroup(&pbPG))
		return nil
	})
	return groups, err
}

// CreateSubnetGroup persists a new DB subnet group.
func (s *NeptuneStore) CreateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subnetGroups.Bucket().Has([]byte(sg.DBSubnetGroupName)) {
		return ErrDBSubnetGroupAlreadyExists
	}
	data, err := proto.Marshal(SubnetGroupToProto(sg))
	if err != nil {
		return NewStoreError("create_subnet_group", err)
	}
	return NewStoreError("create_subnet_group", s.subnetGroups.Bucket().Put([]byte(sg.DBSubnetGroupName), data))
}

// GetSubnetGroup retrieves a DB subnet group by name.
func (s *NeptuneStore) GetSubnetGroup(name string) (*DBSubnetGroup, error) {
	data, err := s.subnetGroups.Bucket().Get([]byte(name))
	if err != nil || data == nil {
		return nil, ErrDBSubnetGroupNotFound
	}
	var pbSG pb.DBSubnetGroup
	if err := proto.Unmarshal(data, &pbSG); err != nil {
		return nil, NewStoreError("get_subnet_group", err)
	}
	return ProtoToSubnetGroup(&pbSG), nil
}

// UpdateSubnetGroup replaces the persisted state of an existing subnet group.
func (s *NeptuneStore) UpdateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.subnetGroups.Bucket().Has([]byte(sg.DBSubnetGroupName)) {
		return ErrDBSubnetGroupNotFound
	}
	data, err := proto.Marshal(SubnetGroupToProto(sg))
	if err != nil {
		return NewStoreError("update_subnet_group", err)
	}
	return NewStoreError("update_subnet_group", s.subnetGroups.Bucket().Put([]byte(sg.DBSubnetGroupName), data))
}

// DeleteSubnetGroup removes a DB subnet group.
func (s *NeptuneStore) DeleteSubnetGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_subnet_group", s.subnetGroups.Delete(name))
}

// ListSubnetGroups returns all DB subnet groups.
func (s *NeptuneStore) ListSubnetGroups() ([]*DBSubnetGroup, error) {
	var groups []*DBSubnetGroup
	err := s.subnetGroups.ForEach(func(key string, value []byte) error {
		var pbSG pb.DBSubnetGroup
		if err := proto.Unmarshal(value, &pbSG); err != nil {
			return err
		}
		groups = append(groups, ProtoToSubnetGroup(&pbSG))
		return nil
	})
	return groups, err
}

// CreateGlobalCluster persists a new global cluster.
func (s *NeptuneStore) CreateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.globalClusters.Bucket().Has([]byte(gc.GlobalClusterIdentifier)) {
		return ErrGlobalClusterAlreadyExists
	}
	data, err := proto.Marshal(GlobalClusterToProto(gc))
	if err != nil {
		return NewStoreError("create_global_cluster", err)
	}
	return NewStoreError("create_global_cluster", s.globalClusters.Bucket().Put([]byte(gc.GlobalClusterIdentifier), data))
}

// GetGlobalCluster retrieves a global cluster by its identifier.
func (s *NeptuneStore) GetGlobalCluster(id string) (*GlobalCluster, error) {
	data, err := s.globalClusters.Bucket().Get([]byte(id))
	if err != nil || data == nil {
		return nil, ErrGlobalClusterNotFound
	}
	var pbGC pb.GlobalCluster
	if err := proto.Unmarshal(data, &pbGC); err != nil {
		return nil, NewStoreError("get_global_cluster", err)
	}
	return ProtoToGlobalCluster(&pbGC), nil
}

// UpdateGlobalCluster replaces the persisted state of an existing global cluster.
func (s *NeptuneStore) UpdateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.globalClusters.Bucket().Has([]byte(gc.GlobalClusterIdentifier)) {
		return ErrGlobalClusterNotFound
	}
	data, err := proto.Marshal(GlobalClusterToProto(gc))
	if err != nil {
		return NewStoreError("update_global_cluster", err)
	}
	return NewStoreError("update_global_cluster", s.globalClusters.Bucket().Put([]byte(gc.GlobalClusterIdentifier), data))
}

// DeleteGlobalCluster removes a global cluster from the store.
func (s *NeptuneStore) DeleteGlobalCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_global_cluster", s.globalClusters.Delete(id))
}

// ListGlobalClusters returns all global clusters in the store.
func (s *NeptuneStore) ListGlobalClusters() ([]*GlobalCluster, error) {
	var clusters []*GlobalCluster
	err := s.globalClusters.ForEach(func(key string, value []byte) error {
		var pbGC pb.GlobalCluster
		if err := proto.Unmarshal(value, &pbGC); err != nil {
			return err
		}
		clusters = append(clusters, ProtoToGlobalCluster(&pbGC))
		return nil
	})
	return clusters, err
}

// CreateEventSubscription persists a new event subscription.
func (s *NeptuneStore) CreateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.eventSubs.Bucket().Has([]byte(sub.CustSubscriptionId)) {
		return ErrEventSubscriptionAlreadyExists
	}
	data, err := proto.Marshal(EventSubscriptionToProto(sub))
	if err != nil {
		return NewStoreError("create_event_subscription", err)
	}
	return NewStoreError("create_event_subscription", s.eventSubs.Bucket().Put([]byte(sub.CustSubscriptionId), data))
}

// GetEventSubscription retrieves an event subscription by its subscription ID.
func (s *NeptuneStore) GetEventSubscription(name string) (*EventSubscription, error) {
	data, err := s.eventSubs.Bucket().Get([]byte(name))
	if err != nil || data == nil {
		return nil, ErrEventSubscriptionNotFound
	}
	var pbSub pb.EventSubscription
	if err := proto.Unmarshal(data, &pbSub); err != nil {
		return nil, NewStoreError("get_event_subscription", err)
	}
	return ProtoToEventSubscription(&pbSub), nil
}

// UpdateEventSubscription replaces the persisted state of an event subscription.
func (s *NeptuneStore) UpdateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.eventSubs.Bucket().Has([]byte(sub.CustSubscriptionId)) {
		return ErrEventSubscriptionNotFound
	}
	data, err := proto.Marshal(EventSubscriptionToProto(sub))
	if err != nil {
		return NewStoreError("update_event_subscription", err)
	}
	return NewStoreError("update_event_subscription", s.eventSubs.Bucket().Put([]byte(sub.CustSubscriptionId), data))
}

// DeleteEventSubscription removes an event subscription.
func (s *NeptuneStore) DeleteEventSubscription(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_event_subscription", s.eventSubs.Delete(name))
}

// ListEventSubscriptions returns all event subscriptions.
func (s *NeptuneStore) ListEventSubscriptions() ([]*EventSubscription, error) {
	var subs []*EventSubscription
	err := s.eventSubs.ForEach(func(key string, value []byte) error {
		var pbSub pb.EventSubscription
		if err := proto.Unmarshal(value, &pbSub); err != nil {
			return err
		}
		subs = append(subs, ProtoToEventSubscription(&pbSub))
		return nil
	})
	return subs, err
}

// AddTags merges the given tags into the existing tags for a resource. Tags
// with the same key are overwritten.
func (s *NeptuneStore) AddTags(resourceArn string, tags []Tag) error {
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
	merged := make([]Tag, 0, len(tagMap))
	for k, v := range tagMap {
		merged = append(merged, Tag{Key: k, Value: v})
	}
	pbTags := &pb.TagList{Tags: TagsToProto(merged)}
	data, err := proto.Marshal(pbTags)
	if err != nil {
		return NewStoreError("add_tags", err)
	}
	return NewStoreError("add_tags", s.tags.Bucket().Put([]byte(resourceArn), data))
}

// GetTags returns the tags associated with a resource ARN.
func (s *NeptuneStore) GetTags(resourceArn string) ([]Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getTagsUnlocked(resourceArn)
}

// getTagsUnlocked retrieves tags without acquiring the lock. Caller must hold
// at least a read lock.
func (s *NeptuneStore) getTagsUnlocked(resourceArn string) ([]Tag, error) {
	data, err := s.tags.Bucket().Get([]byte(resourceArn))
	if err != nil {
		return nil, NewStoreError("get_tags", err)
	}
	if data == nil {
		return []Tag{}, nil
	}
	var pbTags pb.TagList
	if err := proto.Unmarshal(data, &pbTags); err != nil {
		return nil, NewStoreError("get_tags_unmarshal", err)
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
	filtered := make([]Tag, 0, len(existing))
	for _, t := range existing {
		if !removeSet[t.Key] {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == 0 {
		return NewStoreError("remove_tags", s.tags.Delete(resourceArn))
	}
	pbTags := &pb.TagList{Tags: TagsToProto(filtered)}
	data, err := proto.Marshal(pbTags)
	if err != nil {
		return NewStoreError("remove_tags", err)
	}
	return NewStoreError("remove_tags", s.tags.Bucket().Put([]byte(resourceArn), data))
}

// CreateClusterEndpoint persists a new custom cluster endpoint.
func (s *NeptuneStore) CreateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clusterEndpoints.Bucket().Has([]byte(ep.DBClusterEndpointIdentifier)) {
		return ErrDBClusterEndpointAlreadyExists
	}
	data, err := proto.Marshal(ClusterEndpointToProto(ep))
	if err != nil {
		return NewStoreError("create_cluster_endpoint", err)
	}
	return NewStoreError("create_cluster_endpoint", s.clusterEndpoints.Bucket().Put([]byte(ep.DBClusterEndpointIdentifier), data))
}

// GetClusterEndpoint retrieves a cluster endpoint by its identifier.
func (s *NeptuneStore) GetClusterEndpoint(id string) (*DBClusterEndpoint, error) {
	data, err := s.clusterEndpoints.Bucket().Get([]byte(id))
	if err != nil || data == nil {
		return nil, ErrDBClusterEndpointNotFound
	}
	var pbEP pb.DBClusterEndpoint
	if err := proto.Unmarshal(data, &pbEP); err != nil {
		return nil, NewStoreError("get_cluster_endpoint", err)
	}
	return ProtoToClusterEndpoint(&pbEP), nil
}

// UpdateClusterEndpoint replaces the persisted state of a cluster endpoint.
func (s *NeptuneStore) UpdateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.clusterEndpoints.Bucket().Has([]byte(ep.DBClusterEndpointIdentifier)) {
		return ErrDBClusterEndpointNotFound
	}
	data, err := proto.Marshal(ClusterEndpointToProto(ep))
	if err != nil {
		return NewStoreError("update_cluster_endpoint", err)
	}
	return NewStoreError("update_cluster_endpoint", s.clusterEndpoints.Bucket().Put([]byte(ep.DBClusterEndpointIdentifier), data))
}

// DeleteClusterEndpoint removes a cluster endpoint from the store.
func (s *NeptuneStore) DeleteClusterEndpoint(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.clusterEndpoints.Bucket().Get([]byte(id))
	if err != nil {
		return NewStoreError("delete_cluster_endpoint", err)
	}
	if data == nil {
		return ErrDBClusterEndpointNotFound
	}
	return NewStoreError("delete_cluster_endpoint", s.clusterEndpoints.Delete(id))
}

// ListClusterEndpoints returns all cluster endpoints, optionally filtered by
// the parent cluster ID. An empty clusterID returns all endpoints.
func (s *NeptuneStore) ListClusterEndpoints(clusterID string) ([]*DBClusterEndpoint, error) {
	var endpoints []*DBClusterEndpoint
	err := s.clusterEndpoints.ForEach(func(key string, value []byte) error {
		var pbEP pb.DBClusterEndpoint
		if err := proto.Unmarshal(value, &pbEP); err != nil {
			return err
		}
		if clusterID == "" || pbEP.GetDbClusterIdentifier() == clusterID {
			endpoints = append(endpoints, ProtoToClusterEndpoint(&pbEP))
		}
		return nil
	})
	return endpoints, err
}

// CreateQuery persists a new Neptune Data API query state.
func (s *NeptuneStore) CreateQuery(q *pb.QueryState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := proto.Marshal(q)
	if err != nil {
		return NewStoreError("create_query", err)
	}
	return NewStoreError("create_query", s.queries.Bucket().Put([]byte(q.GetQueryId()), data))
}

// GetQuery retrieves a Neptune Data API query state by ID.
func (s *NeptuneStore) GetQuery(id string) (*pb.QueryState, error) {
	data, err := s.queries.Bucket().Get([]byte(id))
	if err != nil || data == nil {
		return nil, nil
	}
	var q pb.QueryState
	if err := proto.Unmarshal(data, &q); err != nil {
		return nil, NewStoreError("get_query", err)
	}
	return &q, nil
}

// UpdateQuery replaces the persisted state of a Neptune Data API query.
func (s *NeptuneStore) UpdateQuery(q *pb.QueryState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := proto.Marshal(q)
	if err != nil {
		return NewStoreError("update_query", err)
	}
	return NewStoreError("update_query", s.queries.Bucket().Put([]byte(q.GetQueryId()), data))
}

// DeleteQuery removes a Neptune Data API query state.
func (s *NeptuneStore) DeleteQuery(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_query", s.queries.Delete(id))
}

// ListQueries returns all Neptune Data API query states.
func (s *NeptuneStore) ListQueries() ([]*pb.QueryState, error) {
	var queries []*pb.QueryState
	err := s.queries.ForEach(func(key string, value []byte) error {
		var q pb.QueryState
		if err := proto.Unmarshal(value, &q); err != nil {
			return err
		}
		queries = append(queries, &q)
		return nil
	})
	return queries, err
}

// CreateLoaderJob persists a new Neptune Data API bulk loader job.
func (s *NeptuneStore) CreateLoaderJob(job *pb.LoaderJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := proto.Marshal(job)
	if err != nil {
		return NewStoreError("create_loader_job", err)
	}
	return NewStoreError("create_loader_job", s.loaderJobs.Bucket().Put([]byte(job.GetLoadId()), data))
}

// GetLoaderJob retrieves a Neptune Data API bulk loader job by ID.
func (s *NeptuneStore) GetLoaderJob(id string) (*pb.LoaderJob, error) {
	data, err := s.loaderJobs.Bucket().Get([]byte(id))
	if err != nil || data == nil {
		return nil, nil
	}
	var job pb.LoaderJob
	if err := proto.Unmarshal(data, &job); err != nil {
		return nil, NewStoreError("get_loader_job", err)
	}
	return &job, nil
}

// UpdateLoaderJob replaces the persisted state of a Neptune Data API bulk loader job.
func (s *NeptuneStore) UpdateLoaderJob(job *pb.LoaderJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := proto.Marshal(job)
	if err != nil {
		return NewStoreError("update_loader_job", err)
	}
	return NewStoreError("update_loader_job", s.loaderJobs.Bucket().Put([]byte(job.GetLoadId()), data))
}

// DeleteLoaderJob removes a Neptune Data API bulk loader job.
func (s *NeptuneStore) DeleteLoaderJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return NewStoreError("delete_loader_job", s.loaderJobs.Delete(id))
}

// ListLoaderJobs returns all Neptune Data API bulk loader jobs.
func (s *NeptuneStore) ListLoaderJobs() ([]*pb.LoaderJob, error) {
	var jobs []*pb.LoaderJob
	err := s.loaderJobs.ForEach(func(key string, value []byte) error {
		var job pb.LoaderJob
		if err := proto.Unmarshal(value, &job); err != nil {
			return err
		}
		jobs = append(jobs, &job)
		return nil
	})
	return jobs, err
}

// Event represents a Neptune database event with metadata.
type Event struct {
	EventID          string    `json:"event_id"`
	Date             time.Time `json:"date"`
	EventCategories  []string  `json:"event_categories"`
	Message          string    `json:"message"`
	SourceArn        string    `json:"source_arn"`
	SourceIdentifier string    `json:"source_identifier"`
	SourceType       string    `json:"source_type"`
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
	data, err := json.Marshal(evt)
	if err != nil {
		return NewStoreError("record_event", err)
	}
	return NewStoreError("record_event", s.events.Bucket().Put([]byte(evt.EventID), data))
}

// EventListOptions controls filtering and pagination for event listing.
type EventListOptions struct {
	SourceType       string
	SourceIdentifier string
	StartTime        time.Time
	EndTime          time.Time
	Marker           string
	MaxRecords       int
}

// EventListResult holds a page of events with pagination metadata.
type EventListResult struct {
	Events      []*Event
	Marker      string
	IsTruncated bool
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
