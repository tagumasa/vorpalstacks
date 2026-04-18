package neptune

import (
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_neptune"
	"vorpalstacks/internal/utils/aws/types"
)

// NeptuneStoreInterface aggregates all Neptune Management API resource CRUD
// operations into a single interface for use by service handlers.
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

// ParameterGroupOps defines CRUD operations for DB parameter groups.
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
	AddTags(resourceArn string, tags []types.Tag) error
	GetTags(resourceArn string) ([]types.Tag, error)
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

// EventOps defines operations for recording and listing Neptune database events.
type EventOps interface {
	RecordEvent(evt *Event) error
	ListEvents(opts EventListOptions) (*EventListResult, error)
	PurgeOldEvents() error
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

var (
	_ NeptuneStoreInterface     = (*NeptuneStore)(nil)
	_ NeptuneDataStoreInterface = (*NeptuneStore)(nil)
)
