package rds

import (
	"fmt"
	"sort"
	"sync"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_rds"

	"google.golang.org/protobuf/proto"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const (
	maxEventAge = 7 * 24 * time.Hour
	maxEvents   = 10000
)

// validEventSourceTypes is the closed set of AWS-spec SourceType wire
// strings accepted by RecordEvent. Unknown values are rejected up front
// so that downstream serialisation never has to silently misclassify an
// event. The set mirrors the proto SourceType enum values and the
// sourceTypeToString / stringToSourceType tables in admin_handler.go.
//
// AWS reference:
// https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_DescribeEvents.html
// Valid Values: db-instance | db-parameter-group | db-security-group |
// db-snapshot | db-cluster | db-cluster-snapshot | custom-engine-version |
// db-proxy | blue-green-deployment | db-shard-group | zero-etl
var validEventSourceTypes = map[string]struct{}{
	"db-instance":           {},
	"db-parameter-group":    {},
	"db-security-group":     {},
	"db-snapshot":           {},
	"db-cluster":            {},
	"db-cluster-snapshot":   {},
	"custom-engine-version": {},
	"db-proxy":              {},
	"blue-green-deployment": {},
	"db-shard-group":        {},
	"zero-etl":              {},
}

// ErrInvalidEventSourceType is returned by RecordEvent when the caller
// supplies a SourceType outside validEventSourceTypes.
var ErrInvalidEventSourceType = fmt.Errorf("invalid SourceType")

type BucketConfig struct {
	ClusterBucket           string
	InstanceBucket          string
	SnapshotBucket          string
	InstSnapshotBucket      string
	ClusterParamGroupBucket string
	ParamGroupBucket        string
	SubnetGroupBucket       string
	GlobalClusterBucket     string
	EventSubBucket          string
	EventsBucket            string
	TagsBucket              string
	Namespace               string
}

type RDSStore struct {
	clusters           *common.ProtoStore[DBCluster]
	instances          *common.ProtoStore[DBInstance]
	snapshots          *common.ProtoStore[DBClusterSnapshot]
	instSnapshots      *common.BaseStore
	clusterParamGroups *common.ProtoStore[DBClusterParameterGroup]
	paramGroups        *common.ProtoStore[DBParameterGroup]
	subnetGroups       *common.ProtoStore[DBSubnetGroup]
	globalClusters     *common.ProtoStore[GlobalCluster]
	eventSubs          *common.ProtoStore[EventSubscription]
	events             *common.ProtoStore[Event]
	tags               *common.BaseStore
	mu                 sync.RWMutex
}

func NewRDSStore(store storage.BasicStorage, cfg BucketConfig) *RDSStore {
	ns := cfg.Namespace
	return &RDSStore{
		clusters: common.NewProtoStore(common.ProtoStoreConfig[DBCluster]{
			Store:        common.NewBaseStore(store.Bucket(cfg.ClusterBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBCluster{} },
			ToDomain:     func(m proto.Message) (*DBCluster, error) { return ProtoToCluster(m.(*pb.DBCluster)), nil },
			ToProto:      func(d *DBCluster) (proto.Message, error) { return ClusterToProto(d), nil },
			IDFunc:       func(d *DBCluster) string { return d.DBClusterIdentifier },
			NotFoundErr:  ErrDBClusterNotFound,
			AlreadyExist: ErrDBClusterAlreadyExists,
		}),
		instances: common.NewProtoStore(common.ProtoStoreConfig[DBInstance]{
			Store:        common.NewBaseStore(store.Bucket(cfg.InstanceBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBInstance{} },
			ToDomain:     func(m proto.Message) (*DBInstance, error) { return ProtoToInstance(m.(*pb.DBInstance)), nil },
			ToProto:      func(d *DBInstance) (proto.Message, error) { return InstanceToProto(d), nil },
			IDFunc:       func(d *DBInstance) string { return d.DBInstanceIdentifier },
			NotFoundErr:  ErrDBInstanceNotFound,
			AlreadyExist: ErrDBInstanceAlreadyExists,
		}),
		snapshots: common.NewProtoStore(common.ProtoStoreConfig[DBClusterSnapshot]{
			Store:    common.NewBaseStore(store.Bucket(cfg.SnapshotBucket), ns),
			NewProto: func() proto.Message { return &pb.DBClusterSnapshot{} },
			ToDomain: func(m proto.Message) (*DBClusterSnapshot, error) {
				return ProtoToSnapshot(m.(*pb.DBClusterSnapshot)), nil
			},
			ToProto:      func(d *DBClusterSnapshot) (proto.Message, error) { return SnapshotToProto(d), nil },
			IDFunc:       func(d *DBClusterSnapshot) string { return d.DBClusterSnapshotIdentifier },
			NotFoundErr:  ErrDBClusterSnapshotNotFound,
			AlreadyExist: ErrDBClusterSnapshotAlreadyExists,
		}),
		clusterParamGroups: common.NewProtoStore(common.ProtoStoreConfig[DBClusterParameterGroup]{
			Store:    common.NewBaseStore(store.Bucket(cfg.ClusterParamGroupBucket), ns),
			NewProto: func() proto.Message { return &pb.DBClusterParameterGroup{} },
			ToDomain: func(m proto.Message) (*DBClusterParameterGroup, error) {
				return ProtoToClusterParameterGroup(m.(*pb.DBClusterParameterGroup)), nil
			},
			ToProto:      func(d *DBClusterParameterGroup) (proto.Message, error) { return ClusterParameterGroupToProto(d), nil },
			IDFunc:       func(d *DBClusterParameterGroup) string { return d.DBClusterParameterGroupName },
			NotFoundErr:  ErrDBClusterParameterGroupNotFound,
			AlreadyExist: ErrDBClusterParameterGroupAlreadyExists,
		}),
		paramGroups: common.NewProtoStore(common.ProtoStoreConfig[DBParameterGroup]{
			Store:    common.NewBaseStore(store.Bucket(cfg.ParamGroupBucket), ns),
			NewProto: func() proto.Message { return &pb.DBParameterGroup{} },
			ToDomain: func(m proto.Message) (*DBParameterGroup, error) {
				return ProtoToParameterGroup(m.(*pb.DBParameterGroup)), nil
			},
			ToProto:      func(d *DBParameterGroup) (proto.Message, error) { return ParameterGroupToProto(d), nil },
			IDFunc:       func(d *DBParameterGroup) string { return d.DBParameterGroupName },
			NotFoundErr:  ErrDBParameterGroupNotFound,
			AlreadyExist: ErrDBParameterGroupAlreadyExists,
		}),
		subnetGroups: common.NewProtoStore(common.ProtoStoreConfig[DBSubnetGroup]{
			Store:        common.NewBaseStore(store.Bucket(cfg.SubnetGroupBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBSubnetGroup{} },
			ToDomain:     func(m proto.Message) (*DBSubnetGroup, error) { return ProtoToSubnetGroup(m.(*pb.DBSubnetGroup)), nil },
			ToProto:      func(d *DBSubnetGroup) (proto.Message, error) { return SubnetGroupToProto(d), nil },
			IDFunc:       func(d *DBSubnetGroup) string { return d.DBSubnetGroupName },
			NotFoundErr:  ErrDBSubnetGroupNotFound,
			AlreadyExist: ErrDBSubnetGroupAlreadyExists,
		}),
		globalClusters: common.NewProtoStore(common.ProtoStoreConfig[GlobalCluster]{
			Store:        common.NewBaseStore(store.Bucket(cfg.GlobalClusterBucket), ns),
			NewProto:     func() proto.Message { return &pb.GlobalCluster{} },
			ToDomain:     func(m proto.Message) (*GlobalCluster, error) { return ProtoToGlobalCluster(m.(*pb.GlobalCluster)), nil },
			ToProto:      func(d *GlobalCluster) (proto.Message, error) { return GlobalClusterToProto(d), nil },
			IDFunc:       func(d *GlobalCluster) string { return d.GlobalClusterIdentifier },
			NotFoundErr:  ErrGlobalClusterNotFound,
			AlreadyExist: ErrGlobalClusterAlreadyExists,
		}),
		eventSubs: common.NewProtoStore(common.ProtoStoreConfig[EventSubscription]{
			Store:    common.NewBaseStore(store.Bucket(cfg.EventSubBucket), ns),
			NewProto: func() proto.Message { return &pb.EventSubscription{} },
			ToDomain: func(m proto.Message) (*EventSubscription, error) {
				return ProtoToEventSubscription(m.(*pb.EventSubscription)), nil
			},
			ToProto:      func(d *EventSubscription) (proto.Message, error) { return EventSubscriptionToProto(d), nil },
			IDFunc:       func(d *EventSubscription) string { return d.CustSubscriptionId },
			NotFoundErr:  ErrEventSubscriptionNotFound,
			AlreadyExist: ErrEventSubscriptionAlreadyExists,
		}),
		events: common.NewProtoStore(common.ProtoStoreConfig[Event]{
			Store:        common.NewBaseStore(store.Bucket(cfg.EventsBucket), ns),
			NewProto:     func() proto.Message { return &pb.Event{} },
			ToDomain:     func(m proto.Message) (*Event, error) { return ProtoToEvent(m.(*pb.Event)), nil },
			ToProto:      func(d *Event) (proto.Message, error) { return EventToProto(d), nil },
			IDFunc:       func(d *Event) string { return d.EventID },
			NotFoundErr:  ErrEventNotFound,
			AlreadyExist: ErrEventAlreadyExists,
		}),
		tags:          common.NewBaseStore(store.Bucket(cfg.TagsBucket), ns),
		instSnapshots: common.NewBaseStore(store.Bucket(cfg.InstSnapshotBucket), ns),
	}
}

func (s *RDSStore) CreateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Create(cluster)
}

func (s *RDSStore) GetCluster(id string) (*DBCluster, error) {
	return s.clusters.Get(id)
}

func (s *RDSStore) UpdateCluster(cluster *DBCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Update(cluster)
}

func (s *RDSStore) DeleteCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusters.Delete(id)
}

func (s *RDSStore) ListClusters() ([]*DBCluster, error) {
	return s.clusters.List()
}

func (s *RDSStore) CreateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Create(instance)
}

func (s *RDSStore) GetInstance(id string) (*DBInstance, error) {
	return s.instances.Get(id)
}

func (s *RDSStore) UpdateInstance(instance *DBInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Update(instance)
}

func (s *RDSStore) DeleteInstance(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instances.Delete(id)
}

func (s *RDSStore) ListInstances() ([]*DBInstance, error) {
	return s.instances.List()
}

func (s *RDSStore) CreateSnapshot(snapshot *DBClusterSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots.Create(snapshot)
}

func (s *RDSStore) GetSnapshot(id string) (*DBClusterSnapshot, error) {
	return s.snapshots.Get(id)
}

func (s *RDSStore) DeleteSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots.Delete(id)
}

func (s *RDSStore) UpdateSnapshot(snapshot *DBClusterSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots.Update(snapshot)
}

func (s *RDSStore) ListSnapshots() ([]*DBClusterSnapshot, error) {
	return s.snapshots.List()
}

func (s *RDSStore) CreateInstanceSnapshot(snap *DBInstanceSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.instSnapshots.Exists(snap.DBSnapshotIdentifier) {
		return ErrDBSnapshotAlreadyExists
	}
	return s.instSnapshots.Put(snap.DBSnapshotIdentifier, snap)
}

func (s *RDSStore) GetInstanceSnapshot(id string) (*DBInstanceSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var snap DBInstanceSnapshot
	if err := s.instSnapshots.Get(id, &snap); err != nil {
		return nil, ErrDBSnapshotNotFound
	}
	return &snap, nil
}

func (s *RDSStore) ListInstanceSnapshots() ([]*DBInstanceSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return common.ListAll[DBInstanceSnapshot](s.instSnapshots)
}

func (s *RDSStore) DeleteInstanceSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instSnapshots.Delete(id)
}

func (s *RDSStore) UpdateInstanceSnapshot(snap *DBInstanceSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instSnapshots.Put(snap.DBSnapshotIdentifier, snap)
}

func (s *RDSStore) CreateClusterParameterGroup(pg *DBClusterParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterParamGroups.Create(pg)
}

func (s *RDSStore) GetClusterParameterGroup(name string) (*DBClusterParameterGroup, error) {
	return s.clusterParamGroups.Get(name)
}

func (s *RDSStore) DeleteClusterParameterGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterParamGroups.Delete(name)
}

func (s *RDSStore) UpdateClusterParameterGroup(pg *DBClusterParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clusterParamGroups.Update(pg)
}

func (s *RDSStore) ListClusterParameterGroups() ([]*DBClusterParameterGroup, error) {
	return s.clusterParamGroups.List()
}

func (s *RDSStore) CreateParameterGroup(pg *DBParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paramGroups.Create(pg)
}

func (s *RDSStore) GetParameterGroup(name string) (*DBParameterGroup, error) {
	return s.paramGroups.Get(name)
}

func (s *RDSStore) DeleteParameterGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paramGroups.Delete(name)
}

func (s *RDSStore) UpdateParameterGroup(pg *DBParameterGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paramGroups.Update(pg)
}

func (s *RDSStore) ListParameterGroups() ([]*DBParameterGroup, error) {
	return s.paramGroups.List()
}

func (s *RDSStore) CreateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Create(sg)
}

func (s *RDSStore) GetSubnetGroup(name string) (*DBSubnetGroup, error) {
	return s.subnetGroups.Get(name)
}

func (s *RDSStore) UpdateSubnetGroup(sg *DBSubnetGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Update(sg)
}

func (s *RDSStore) DeleteSubnetGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subnetGroups.Delete(name)
}

func (s *RDSStore) ListSubnetGroups() ([]*DBSubnetGroup, error) {
	return s.subnetGroups.List()
}

func (s *RDSStore) CreateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Create(gc)
}

func (s *RDSStore) GetGlobalCluster(id string) (*GlobalCluster, error) {
	return s.globalClusters.Get(id)
}

func (s *RDSStore) UpdateGlobalCluster(gc *GlobalCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Update(gc)
}

func (s *RDSStore) DeleteGlobalCluster(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.globalClusters.Delete(id)
}

func (s *RDSStore) ListGlobalClusters() ([]*GlobalCluster, error) {
	return s.globalClusters.List()
}

func (s *RDSStore) CreateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Create(sub)
}

func (s *RDSStore) GetEventSubscription(name string) (*EventSubscription, error) {
	return s.eventSubs.Get(name)
}

func (s *RDSStore) UpdateEventSubscription(sub *EventSubscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Update(sub)
}

func (s *RDSStore) DeleteEventSubscription(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventSubs.Delete(name)
}

func (s *RDSStore) ListEventSubscriptions() ([]*EventSubscription, error) {
	return s.eventSubs.List()
}

// ListClusterEndpoints returns cluster endpoints filtered by cluster ID.
// The base implementation returns an empty list; NeptuneStore overrides this
// with its own ProtoStore-backed implementation.
func (s *RDSStore) ListClusterEndpoints(clusterID string) ([]*DBClusterEndpoint, error) {
	return nil, nil
}

func (s *RDSStore) AddTags(resourceArn string, tags []types.Tag) error {
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
	sort.Slice(merged, func(i, j int) bool { return merged[i].Key < merged[j].Key })
	return s.tags.PutProto(resourceArn, &pb.TagList{Tags: TagsToProto(merged)})
}

func (s *RDSStore) GetTags(resourceArn string) ([]types.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getTagsUnlocked(resourceArn)
}

func (s *RDSStore) getTagsUnlocked(resourceArn string) ([]types.Tag, error) {
	var pbTags pb.TagList
	if err := s.tags.GetProto(resourceArn, &pbTags); err != nil {
		if common.IsNotFound(err) {
			return []types.Tag{}, nil
		}
		return nil, err
	}
	return ProtoToTags(pbTags.Tags), nil
}

func (s *RDSStore) RemoveTags(resourceArn string, keys []string) error {
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

func (s *RDSStore) RecordEvent(evt *Event) error {
	// Reject unknown SourceType at the entry point so we never have to
	// fabricate a category during serialisation (stringToSourceType in
	// admin_handler previously mapped unknowns to BLUE_GREEN_DEPLOYMENT,
	// the proto3 zero value, silently corrupting audit data). An empty
	// SourceType is permitted because RecordEvent is also called for
	// ad-hoc platform events that do not map to a single source.
	if evt.SourceType != "" {
		if _, ok := validEventSourceTypes[evt.SourceType]; !ok {
			return fmt.Errorf("%w: %q", ErrInvalidEventSourceType, evt.SourceType)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// Allocate a local ID if the caller did not provide one. Mutating the
	// caller-supplied Event is a surprising side-effect (the caller may
	// reuse the struct across goroutines or test cases); copy first.
	id := evt.EventID
	if id == "" {
		id = fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	stored := *evt
	stored.EventID = id
	return s.events.Create(&stored)
}

func (s *RDSStore) ListEvents(opts EventListOptions) (*EventListResult, error) {
	if opts.MaxRecords <= 0 {
		opts.MaxRecords = 100
	}

	cutoff := time.Now().Add(-maxEventAge)

	allEvents, err := s.events.ListFiltered(func(evt *Event) bool {
		if evt.Date.Before(cutoff) {
			return false
		}
		if opts.SourceType != "" && evt.SourceType != opts.SourceType {
			return false
		}
		if opts.SourceIdentifier != "" && evt.SourceIdentifier != opts.SourceIdentifier {
			return false
		}
		if !opts.StartTime.IsZero() && evt.Date.Before(opts.StartTime) {
			return false
		}
		if !opts.EndTime.IsZero() && evt.Date.After(opts.EndTime) {
			return false
		}
		if len(opts.EventCategories) > 0 && !intersectsCategories(evt.EventCategories, opts.EventCategories) {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	started := opts.Marker == ""
	markerFound := started // when Marker is empty we start at the head
	var remaining []*Event
	for _, evt := range allEvents {
		if !started {
			if evt.EventID == opts.Marker {
				started = true
				markerFound = true
			}
			continue
		}
		remaining = append(remaining, evt)
		if len(remaining) > opts.MaxRecords {
			break
		}
	}
	if !markerFound {
		// AWS RDS returns InvalidParameterValue when Marker does not match
		// any event. Silent empty results hide pagination bugs in callers
		// (e.g., a marker from a stale or purged event set).
		return nil, ErrInvalidEventMarker
	}

	result := &EventListResult{Events: remaining}
	if len(remaining) > opts.MaxRecords {
		result.Events = remaining[:opts.MaxRecords]
		result.IsTruncated = true
		result.Marker = result.Events[len(result.Events)-1].EventID
	}
	return result, nil
}

func (s *RDSStore) PurgeOldEvents() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-maxEventAge)

	allEvents, err := s.events.List()
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

var _ StoreInterface = (*RDSStore)(nil)

// intersectsCategories reports whether the event's own category list shares
// at least one element with the requested filter set. AWS RDS treats the
// EventCategories filter as a union (OR) within a single source type.
func intersectsCategories(eventCats, filterCats []string) bool {
	if len(filterCats) == 0 {
		return true
	}
	want := make(map[string]struct{}, len(filterCats))
	for _, c := range filterCats {
		want[c] = struct{}{}
	}
	for _, c := range eventCats {
		if _, ok := want[c]; ok {
			return true
		}
	}
	return false
}
