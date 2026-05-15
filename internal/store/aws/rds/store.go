package rds

import (
	"fmt"
	"sort"
	"sync"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_rds"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	maxEventAge = 7 * 24 * time.Hour
	maxEvents   = 10000
)

type BucketConfig struct {
	ClusterBucket           string
	InstanceBucket          string
	SnapshotBucket          string
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
			ToDomain:     func(m proto.Message) *DBCluster { return ProtoToCluster(m.(*pb.DBCluster)) },
			ToProto:      func(d *DBCluster) proto.Message { return ClusterToProto(d) },
			IDFunc:       func(d *DBCluster) string { return d.DBClusterIdentifier },
			NotFoundErr:  ErrDBClusterNotFound,
			AlreadyExist: ErrDBClusterAlreadyExists,
		}),
		instances: common.NewProtoStore(common.ProtoStoreConfig[DBInstance]{
			Store:        common.NewBaseStore(store.Bucket(cfg.InstanceBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBInstance{} },
			ToDomain:     func(m proto.Message) *DBInstance { return ProtoToInstance(m.(*pb.DBInstance)) },
			ToProto:      func(d *DBInstance) proto.Message { return InstanceToProto(d) },
			IDFunc:       func(d *DBInstance) string { return d.DBInstanceIdentifier },
			NotFoundErr:  ErrDBInstanceNotFound,
			AlreadyExist: ErrDBInstanceAlreadyExists,
		}),
		snapshots: common.NewProtoStore(common.ProtoStoreConfig[DBClusterSnapshot]{
			Store:        common.NewBaseStore(store.Bucket(cfg.SnapshotBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBClusterSnapshot{} },
			ToDomain:     func(m proto.Message) *DBClusterSnapshot { return ProtoToSnapshot(m.(*pb.DBClusterSnapshot)) },
			ToProto:      func(d *DBClusterSnapshot) proto.Message { return SnapshotToProto(d) },
			IDFunc:       func(d *DBClusterSnapshot) string { return d.DBClusterSnapshotIdentifier },
			NotFoundErr:  ErrDBClusterSnapshotNotFound,
			AlreadyExist: ErrDBClusterSnapshotAlreadyExists,
		}),
		clusterParamGroups: common.NewProtoStore(common.ProtoStoreConfig[DBClusterParameterGroup]{
			Store:    common.NewBaseStore(store.Bucket(cfg.ClusterParamGroupBucket), ns),
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
			Store:        common.NewBaseStore(store.Bucket(cfg.ParamGroupBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBParameterGroup{} },
			ToDomain:     func(m proto.Message) *DBParameterGroup { return ProtoToParameterGroup(m.(*pb.DBParameterGroup)) },
			ToProto:      func(d *DBParameterGroup) proto.Message { return ParameterGroupToProto(d) },
			IDFunc:       func(d *DBParameterGroup) string { return d.DBParameterGroupName },
			NotFoundErr:  ErrDBParameterGroupNotFound,
			AlreadyExist: ErrDBParameterGroupAlreadyExists,
		}),
		subnetGroups: common.NewProtoStore(common.ProtoStoreConfig[DBSubnetGroup]{
			Store:        common.NewBaseStore(store.Bucket(cfg.SubnetGroupBucket), ns),
			NewProto:     func() proto.Message { return &pb.DBSubnetGroup{} },
			ToDomain:     func(m proto.Message) *DBSubnetGroup { return ProtoToSubnetGroup(m.(*pb.DBSubnetGroup)) },
			ToProto:      func(d *DBSubnetGroup) proto.Message { return SubnetGroupToProto(d) },
			IDFunc:       func(d *DBSubnetGroup) string { return d.DBSubnetGroupName },
			NotFoundErr:  ErrDBSubnetGroupNotFound,
			AlreadyExist: ErrDBSubnetGroupAlreadyExists,
		}),
		globalClusters: common.NewProtoStore(common.ProtoStoreConfig[GlobalCluster]{
			Store:        common.NewBaseStore(store.Bucket(cfg.GlobalClusterBucket), ns),
			NewProto:     func() proto.Message { return &pb.GlobalCluster{} },
			ToDomain:     func(m proto.Message) *GlobalCluster { return ProtoToGlobalCluster(m.(*pb.GlobalCluster)) },
			ToProto:      func(d *GlobalCluster) proto.Message { return GlobalClusterToProto(d) },
			IDFunc:       func(d *GlobalCluster) string { return d.GlobalClusterIdentifier },
			NotFoundErr:  ErrGlobalClusterNotFound,
			AlreadyExist: ErrGlobalClusterAlreadyExists,
		}),
		eventSubs: common.NewProtoStore(common.ProtoStoreConfig[EventSubscription]{
			Store:        common.NewBaseStore(store.Bucket(cfg.EventSubBucket), ns),
			NewProto:     func() proto.Message { return &pb.EventSubscription{} },
			ToDomain:     func(m proto.Message) *EventSubscription { return ProtoToEventSubscription(m.(*pb.EventSubscription)) },
			ToProto:      func(d *EventSubscription) proto.Message { return EventSubscriptionToProto(d) },
			IDFunc:       func(d *EventSubscription) string { return d.CustSubscriptionId },
			NotFoundErr:  ErrEventSubscriptionNotFound,
			AlreadyExist: ErrEventSubscriptionAlreadyExists,
		}),
		events: common.NewProtoStore(common.ProtoStoreConfig[Event]{
			Store:        common.NewBaseStore(store.Bucket(cfg.EventsBucket), ns),
			NewProto:     func() proto.Message { return &pb.Event{} },
			ToDomain:     func(m proto.Message) *Event { return ProtoToEvent(m.(*pb.Event)) },
			ToProto:      func(d *Event) proto.Message { return EventToProto(d) },
			IDFunc:       func(d *Event) string { return d.EventID },
			NotFoundErr:  ErrEventNotFound,
			AlreadyExist: ErrEventAlreadyExists,
		}),
		tags: common.NewBaseStore(store.Bucket(cfg.TagsBucket), ns),
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if evt.EventID == "" {
		evt.EventID = fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	return s.events.Create(evt)
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
		return true
	})
	if err != nil {
		return nil, err
	}

	started := opts.Marker == ""
	var remaining []*Event
	for _, evt := range allEvents {
		if !started {
			if evt.EventID == opts.Marker {
				started = true
			}
			continue
		}
		remaining = append(remaining, evt)
		if len(remaining) > opts.MaxRecords {
			break
		}
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
