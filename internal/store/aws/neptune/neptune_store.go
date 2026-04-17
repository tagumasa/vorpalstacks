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
	if s.clusters.Exists(cluster.DBClusterIdentifier) {
		return ErrDBClusterAlreadyExists
	}
	return s.clusters.PutProto(cluster.DBClusterIdentifier, ClusterToProto(cluster))
}

// GetCluster retrieves a DB cluster by its identifier. Returns
// ErrDBClusterNotFound if the cluster does not exist.
func (s *NeptuneStore) GetCluster(id string) (*DBCluster, error) {
	var pbCluster pb.DBCluster
	if err := s.clusters.GetProto(id, &pbCluster); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBClusterNotFound
		}
		return nil, err
	}
	return ProtoToCluster(&pbCluster), nil
}

// UpdateCluster replaces the persisted state of an existing DB cluster.
func (s *NeptuneStore) UpdateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.clusters.Exists(cluster.DBClusterIdentifier) {
		return ErrDBClusterNotFound
	}
	return s.clusters.PutProto(cluster.DBClusterIdentifier, ClusterToProto(cluster))
}

// DeleteCluster removes a DB cluster from the store.
func (s *NeptuneStore) DeleteCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Delete(id)
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
	if s.instances.Exists(instance.DBInstanceIdentifier) {
		return ErrDBInstanceAlreadyExists
	}
	return s.instances.PutProto(instance.DBInstanceIdentifier, InstanceToProto(instance))
}

// GetInstance retrieves a DB instance by its identifier. Returns
// ErrDBInstanceNotFound if the instance does not exist.
func (s *NeptuneStore) GetInstance(id string) (*DBInstance, error) {
	var pbInstance pb.DBInstance
	if err := s.instances.GetProto(id, &pbInstance); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBInstanceNotFound
		}
		return nil, err
	}
	return ProtoToInstance(&pbInstance), nil
}

// UpdateInstance replaces the persisted state of an existing DB instance.
func (s *NeptuneStore) UpdateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.instances.Exists(instance.DBInstanceIdentifier) {
		return ErrDBInstanceNotFound
	}
	return s.instances.PutProto(instance.DBInstanceIdentifier, InstanceToProto(instance))
}

// DeleteInstance removes a DB instance from the store.
func (s *NeptuneStore) DeleteInstance(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Delete(id)
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
	if s.snapshots.Exists(snapshot.DBClusterSnapshotIdentifier) {
		return ErrDBClusterSnapshotAlreadyExists
	}
	return s.snapshots.PutProto(snapshot.DBClusterSnapshotIdentifier, SnapshotToProto(snapshot))
}

// GetSnapshot retrieves a cluster snapshot by its identifier.
func (s *NeptuneStore) GetSnapshot(id string) (*DBClusterSnapshot, error) {
	var pbSnap pb.DBClusterSnapshot
	if err := s.snapshots.GetProto(id, &pbSnap); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBClusterSnapshotNotFound
		}
		return nil, err
	}
	return ProtoToSnapshot(&pbSnap), nil
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
	if !s.snapshots.Exists(snapshot.DBClusterSnapshotIdentifier) {
		return ErrDBClusterSnapshotNotFound
	}
	return s.snapshots.PutProto(snapshot.DBClusterSnapshotIdentifier, SnapshotToProto(snapshot))
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
	if s.clusterParamGroups.Exists(pg.DBClusterParameterGroupName) {
		return ErrDBClusterParameterGroupAlreadyExists
	}
	return s.clusterParamGroups.PutProto(pg.DBClusterParameterGroupName, ClusterParameterGroupToProto(pg))
}

// GetClusterParameterGroup retrieves a cluster parameter group by name.
func (s *NeptuneStore) GetClusterParameterGroup(name string) (*DBClusterParameterGroup, error) {
	var pbPG pb.DBClusterParameterGroup
	if err := s.clusterParamGroups.GetProto(name, &pbPG); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBClusterParameterGroupNotFound
		}
		return nil, err
	}
	return ProtoToClusterParameterGroup(&pbPG), nil
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
	if !s.clusterParamGroups.Exists(pg.DBClusterParameterGroupName) {
		return ErrDBClusterParameterGroupNotFound
	}
	return s.clusterParamGroups.PutProto(pg.DBClusterParameterGroupName, ClusterParameterGroupToProto(pg))
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
	if s.paramGroups.Exists(pg.DBParameterGroupName) {
		return ErrDBParameterGroupAlreadyExists
	}
	return s.paramGroups.PutProto(pg.DBParameterGroupName, ParameterGroupToProto(pg))
}

// GetParameterGroup retrieves a DB parameter group by name.
func (s *NeptuneStore) GetParameterGroup(name string) (*DBParameterGroup, error) {
	var pbPG pb.DBParameterGroup
	if err := s.paramGroups.GetProto(name, &pbPG); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBParameterGroupNotFound
		}
		return nil, err
	}
	return ProtoToParameterGroup(&pbPG), nil
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
	if !s.paramGroups.Exists(pg.DBParameterGroupName) {
		return ErrDBParameterGroupNotFound
	}
	return s.paramGroups.PutProto(pg.DBParameterGroupName, ParameterGroupToProto(pg))
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
	if s.subnetGroups.Exists(sg.DBSubnetGroupName) {
		return ErrDBSubnetGroupAlreadyExists
	}
	return s.subnetGroups.PutProto(sg.DBSubnetGroupName, SubnetGroupToProto(sg))
}

// GetSubnetGroup retrieves a DB subnet group by name.
func (s *NeptuneStore) GetSubnetGroup(name string) (*DBSubnetGroup, error) {
	var pbSG pb.DBSubnetGroup
	if err := s.subnetGroups.GetProto(name, &pbSG); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBSubnetGroupNotFound
		}
		return nil, err
	}
	return ProtoToSubnetGroup(&pbSG), nil
}

// UpdateSubnetGroup replaces the persisted state of an existing subnet group.
func (s *NeptuneStore) UpdateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.subnetGroups.Exists(sg.DBSubnetGroupName) {
		return ErrDBSubnetGroupNotFound
	}
	return s.subnetGroups.PutProto(sg.DBSubnetGroupName, SubnetGroupToProto(sg))
}

// DeleteSubnetGroup removes a DB subnet group.
func (s *NeptuneStore) DeleteSubnetGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Delete(name)
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
	if s.globalClusters.Exists(gc.GlobalClusterIdentifier) {
		return ErrGlobalClusterAlreadyExists
	}
	return s.globalClusters.PutProto(gc.GlobalClusterIdentifier, GlobalClusterToProto(gc))
}

// GetGlobalCluster retrieves a global cluster by its identifier.
func (s *NeptuneStore) GetGlobalCluster(id string) (*GlobalCluster, error) {
	var pbGC pb.GlobalCluster
	if err := s.globalClusters.GetProto(id, &pbGC); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrGlobalClusterNotFound
		}
		return nil, err
	}
	return ProtoToGlobalCluster(&pbGC), nil
}

// UpdateGlobalCluster replaces the persisted state of an existing global cluster.
func (s *NeptuneStore) UpdateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.globalClusters.Exists(gc.GlobalClusterIdentifier) {
		return ErrGlobalClusterNotFound
	}
	return s.globalClusters.PutProto(gc.GlobalClusterIdentifier, GlobalClusterToProto(gc))
}

// DeleteGlobalCluster removes a global cluster from the store.
func (s *NeptuneStore) DeleteGlobalCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Delete(id)
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
	if s.eventSubs.Exists(sub.CustSubscriptionId) {
		return ErrEventSubscriptionAlreadyExists
	}
	return s.eventSubs.PutProto(sub.CustSubscriptionId, EventSubscriptionToProto(sub))
}

// GetEventSubscription retrieves an event subscription by its subscription ID.
func (s *NeptuneStore) GetEventSubscription(name string) (*EventSubscription, error) {
	var pbSub pb.EventSubscription
	if err := s.eventSubs.GetProto(name, &pbSub); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventSubscriptionNotFound
		}
		return nil, err
	}
	return ProtoToEventSubscription(&pbSub), nil
}

// UpdateEventSubscription replaces the persisted state of an event subscription.
func (s *NeptuneStore) UpdateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.eventSubs.Exists(sub.CustSubscriptionId) {
		return ErrEventSubscriptionNotFound
	}
	return s.eventSubs.PutProto(sub.CustSubscriptionId, EventSubscriptionToProto(sub))
}

// DeleteEventSubscription removes an event subscription.
func (s *NeptuneStore) DeleteEventSubscription(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Delete(name)
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
	return s.tags.PutProto(resourceArn, &pb.TagList{Tags: TagsToProto(merged)})
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
	var pbTags pb.TagList
	if err := s.tags.GetProto(resourceArn, &pbTags); err != nil {
		if common.IsNotFound(err) {
			return []Tag{}, nil
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
	filtered := make([]Tag, 0, len(existing))
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
	if s.clusterEndpoints.Exists(ep.DBClusterEndpointIdentifier) {
		return ErrDBClusterEndpointAlreadyExists
	}
	return s.clusterEndpoints.PutProto(ep.DBClusterEndpointIdentifier, ClusterEndpointToProto(ep))
}

// GetClusterEndpoint retrieves a cluster endpoint by its identifier.
func (s *NeptuneStore) GetClusterEndpoint(id string) (*DBClusterEndpoint, error) {
	var pbEP pb.DBClusterEndpoint
	if err := s.clusterEndpoints.GetProto(id, &pbEP); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrDBClusterEndpointNotFound
		}
		return nil, err
	}
	return ProtoToClusterEndpoint(&pbEP), nil
}

// UpdateClusterEndpoint replaces the persisted state of a cluster endpoint.
func (s *NeptuneStore) UpdateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.clusterEndpoints.Exists(ep.DBClusterEndpointIdentifier) {
		return ErrDBClusterEndpointNotFound
	}
	return s.clusterEndpoints.PutProto(ep.DBClusterEndpointIdentifier, ClusterEndpointToProto(ep))
}

// DeleteClusterEndpoint removes a cluster endpoint from the store.
func (s *NeptuneStore) DeleteClusterEndpoint(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.clusterEndpoints.Exists(id) {
		return ErrDBClusterEndpointNotFound
	}
	return s.clusterEndpoints.Delete(id)
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
	return s.queries.PutProto(q.GetQueryId(), q)
}

// GetQuery retrieves a Neptune Data API query state by ID.
func (s *NeptuneStore) GetQuery(id string) (*pb.QueryState, error) {
	var q pb.QueryState
	if err := s.queries.GetProto(id, &q); err != nil {
		if common.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

// UpdateQuery replaces the persisted state of a Neptune Data API query.
func (s *NeptuneStore) UpdateQuery(q *pb.QueryState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queries.PutProto(q.GetQueryId(), q)
}

// DeleteQuery removes a Neptune Data API query state.
func (s *NeptuneStore) DeleteQuery(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queries.Delete(id)
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
	return s.loaderJobs.PutProto(job.GetLoadId(), job)
}

// GetLoaderJob retrieves a Neptune Data API bulk loader job by ID.
func (s *NeptuneStore) GetLoaderJob(id string) (*pb.LoaderJob, error) {
	var job pb.LoaderJob
	if err := s.loaderJobs.GetProto(id, &job); err != nil {
		if common.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// UpdateLoaderJob replaces the persisted state of a Neptune Data API bulk loader job.
func (s *NeptuneStore) UpdateLoaderJob(job *pb.LoaderJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loaderJobs.PutProto(job.GetLoadId(), job)
}

// DeleteLoaderJob removes a Neptune Data API bulk loader job.
func (s *NeptuneStore) DeleteLoaderJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loaderJobs.Delete(id)
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
