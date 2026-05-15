package neptune

import (
	"fmt"
	"sync"
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_neptune"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	rds "vorpalstacks/internal/store/aws/rds"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	clusterEndpointBucket = "neptune_cluster_endpoints"
	queryStateBucket      = "neptunedata_queries"
	loaderJobBucket       = "neptunedata_loader_jobs"
)

type NeptuneStore struct {
	*rds.RDSStore
	clusterEndpoints *common.ProtoStore[DBClusterEndpoint]
	queries          *common.RawProtoStore[*pb.QueryState]
	loaderJobs       *common.RawProtoStore[*pb.LoaderJob]
	neptuneMu        sync.RWMutex
}

func NewNeptuneStore(store storage.BasicStorage) *NeptuneStore {
	rdsStore := rds.NewRDSStore(store, rds.BucketConfig{
		ClusterBucket:           "neptune_clusters",
		InstanceBucket:          "neptune_instances",
		SnapshotBucket:          "neptune_snapshots",
		ClusterParamGroupBucket: "neptune_cluster_param_groups",
		ParamGroupBucket:        "neptune_param_groups",
		SubnetGroupBucket:       "neptune_subnet_groups",
		GlobalClusterBucket:     "neptune_global_clusters",
		EventSubBucket:          "neptune_event_subscriptions",
		EventsBucket:            "neptune_events",
		TagsBucket:              "neptune_tags",
		Namespace:               "neptune",
	})
	return &NeptuneStore{
		RDSStore: rdsStore,
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
	}
}

func (s *NeptuneStore) CreateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.clusterEndpoints.Create(ep)
}

func (s *NeptuneStore) GetClusterEndpoint(id string) (*DBClusterEndpoint, error) {
	return s.clusterEndpoints.Get(id)
}

func (s *NeptuneStore) UpdateClusterEndpoint(ep *DBClusterEndpoint) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.clusterEndpoints.Update(ep)
}

func (s *NeptuneStore) DeleteClusterEndpoint(id string) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.clusterEndpoints.DeleteIfExists(id)
}

func (s *NeptuneStore) ListClusterEndpoints(clusterID string) ([]*DBClusterEndpoint, error) {
	return s.clusterEndpoints.ListFiltered(func(ep *DBClusterEndpoint) bool {
		return clusterID == "" || ep.DBClusterIdentifier == clusterID
	})
}

func (s *NeptuneStore) CreateQuery(q *pb.QueryState) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.queries.Create(q)
}

func (s *NeptuneStore) GetQuery(id string) (*pb.QueryState, error) {
	return s.queries.Get(id)
}

func (s *NeptuneStore) UpdateQuery(q *pb.QueryState) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.queries.Update(q)
}

func (s *NeptuneStore) DeleteQuery(id string) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.queries.Delete(id)
}

func (s *NeptuneStore) ListQueries() ([]*pb.QueryState, error) {
	return s.queries.List()
}

func (s *NeptuneStore) CreateLoaderJob(job *pb.LoaderJob) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.loaderJobs.Create(job)
}

func (s *NeptuneStore) GetLoaderJob(id string) (*pb.LoaderJob, error) {
	return s.loaderJobs.Get(id)
}

func (s *NeptuneStore) UpdateLoaderJob(job *pb.LoaderJob) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.loaderJobs.Update(job)
}

func (s *NeptuneStore) DeleteLoaderJob(id string) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	return s.loaderJobs.Delete(id)
}

func (s *NeptuneStore) ListLoaderJobs() ([]*pb.LoaderJob, error) {
	return s.loaderJobs.List()
}

func (s *NeptuneStore) RecordEvent(evt *rds.Event) error {
	s.neptuneMu.Lock()
	defer s.neptuneMu.Unlock()
	if evt.EventID == "" {
		evt.EventID = fmt.Sprintf("evt-%d", time.Now().UnixNano())
	}
	return s.RDSStore.RecordEvent(evt)
}

func (s *NeptuneStore) ListEvents(opts rds.EventListOptions) (*rds.EventListResult, error) {
	return s.RDSStore.ListEvents(opts)
}

func (s *NeptuneStore) PurgeOldEvents() error {
	return s.RDSStore.PurgeOldEvents()
}

func (s *NeptuneStore) AddTags(resourceArn string, tags []types.Tag) error {
	return s.RDSStore.AddTags(resourceArn, tags)
}

func (s *NeptuneStore) GetTags(resourceArn string) ([]types.Tag, error) {
	return s.RDSStore.GetTags(resourceArn)
}

func (s *NeptuneStore) RemoveTags(resourceArn string, keys []string) error {
	return s.RDSStore.RemoveTags(resourceArn, keys)
}

var _ NeptuneStoreInterface = (*NeptuneStore)(nil)
