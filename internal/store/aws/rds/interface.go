package rds

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

type StoreInterface interface {
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
}

type ClusterOps interface {
	CreateCluster(cluster *DBCluster) error
	GetCluster(id string) (*DBCluster, error)
	UpdateCluster(cluster *DBCluster) error
	DeleteCluster(id string) error
	ListClusters() ([]*DBCluster, error)
}

type InstanceOps interface {
	CreateInstance(instance *DBInstance) error
	GetInstance(id string) (*DBInstance, error)
	UpdateInstance(instance *DBInstance) error
	DeleteInstance(id string) error
	ListInstances() ([]*DBInstance, error)
}

type SnapshotOps interface {
	CreateSnapshot(snapshot *DBClusterSnapshot) error
	GetSnapshot(id string) (*DBClusterSnapshot, error)
	DeleteSnapshot(id string) error
	UpdateSnapshot(snapshot *DBClusterSnapshot) error
	ListSnapshots() ([]*DBClusterSnapshot, error)
}

type ClusterParameterGroupOps interface {
	CreateClusterParameterGroup(pg *DBClusterParameterGroup) error
	GetClusterParameterGroup(name string) (*DBClusterParameterGroup, error)
	UpdateClusterParameterGroup(pg *DBClusterParameterGroup) error
	DeleteClusterParameterGroup(name string) error
	ListClusterParameterGroups() ([]*DBClusterParameterGroup, error)
}

type ParameterGroupOps interface {
	CreateParameterGroup(pg *DBParameterGroup) error
	GetParameterGroup(name string) (*DBParameterGroup, error)
	UpdateParameterGroup(pg *DBParameterGroup) error
	DeleteParameterGroup(name string) error
	ListParameterGroups() ([]*DBParameterGroup, error)
}

type SubnetGroupOps interface {
	CreateSubnetGroup(sg *DBSubnetGroup) error
	GetSubnetGroup(name string) (*DBSubnetGroup, error)
	UpdateSubnetGroup(sg *DBSubnetGroup) error
	DeleteSubnetGroup(name string) error
	ListSubnetGroups() ([]*DBSubnetGroup, error)
}

type GlobalClusterOps interface {
	CreateGlobalCluster(gc *GlobalCluster) error
	GetGlobalCluster(id string) (*GlobalCluster, error)
	UpdateGlobalCluster(gc *GlobalCluster) error
	DeleteGlobalCluster(id string) error
	ListGlobalClusters() ([]*GlobalCluster, error)
}

type EventSubscriptionOps interface {
	CreateEventSubscription(sub *EventSubscription) error
	GetEventSubscription(name string) (*EventSubscription, error)
	UpdateEventSubscription(sub *EventSubscription) error
	DeleteEventSubscription(name string) error
	ListEventSubscriptions() ([]*EventSubscription, error)
}

type TagOps interface {
	AddTags(resourceArn string, tags []types.Tag) error
	GetTags(resourceArn string) ([]types.Tag, error)
	RemoveTags(resourceArn string, keys []string) error
}

type EventOps interface {
	RecordEvent(evt *Event) error
	ListEvents(opts EventListOptions) (*EventListResult, error)
	PurgeOldEvents() error
}

type Event struct {
	EventID          string    `json:"event_id"`
	Date             time.Time `json:"date"`
	EventCategories  []string  `json:"event_categories"`
	Message          string    `json:"message"`
	SourceArn        string    `json:"source_arn"`
	SourceIdentifier string    `json:"source_identifier"`
	SourceType       string    `json:"source_type"`
}

type EventListOptions struct {
	SourceType       string
	SourceIdentifier string
	StartTime        time.Time
	EndTime          time.Time
	Marker           string
	MaxRecords       int
}

type EventListResult struct {
	Events      []*Event
	Marker      string
	IsTruncated bool
}
