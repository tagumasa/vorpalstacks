package neptune

import (
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	rds "vorpalstacks/internal/store/aws/rds"
)

type NeptuneStoreInterface interface {
	rds.StoreInterface
	NeptuneDataStoreInterface
	ClusterEndpointOps
}

type NeptuneDataStoreInterface interface {
	QueryOps
	LoaderJobOps
}

type QueryOps interface {
	CreateQuery(q *pb.QueryState) error
	GetQuery(id string) (*pb.QueryState, error)
	UpdateQuery(q *pb.QueryState) error
	DeleteQuery(id string) error
	ListQueries() ([]*pb.QueryState, error)
}

type LoaderJobOps interface {
	CreateLoaderJob(job *pb.LoaderJob) error
	GetLoaderJob(id string) (*pb.LoaderJob, error)
	UpdateLoaderJob(job *pb.LoaderJob) error
	DeleteLoaderJob(id string) error
	ListLoaderJobs() ([]*pb.LoaderJob, error)
}

type ClusterEndpointOps interface {
	CreateClusterEndpoint(ep *DBClusterEndpoint) error
	GetClusterEndpoint(id string) (*DBClusterEndpoint, error)
	UpdateClusterEndpoint(ep *DBClusterEndpoint) error
	DeleteClusterEndpoint(id string) error
	ListClusterEndpoints(clusterID string) ([]*DBClusterEndpoint, error)
}

type Event = rds.Event
type EventListOptions = rds.EventListOptions
type EventListResult = rds.EventListResult

type TagOps = rds.TagOps
type EventOps = rds.EventOps
type ClusterOps = rds.ClusterOps
type InstanceOps = rds.InstanceOps
type SnapshotOps = rds.SnapshotOps

var _ NeptuneStoreInterface = (*NeptuneStore)(nil)
