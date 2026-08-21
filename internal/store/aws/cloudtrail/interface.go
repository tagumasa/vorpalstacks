package cloudtrail

import (
	"time"

	types "vorpalstacks/internal/common/tags"

	common "vorpalstacks/internal/store/aws/common"
)

// CloudTrailStoreInterface defines operations for managing CloudTrail trails and events.
type CloudTrailStoreInterface interface {
	GetAccountID() string
	GetRegion() string
	BuildTrailARN(trailName string) string
	CreateTrail(trail *Trail) (*Trail, error)
	GetTrail(trailName string) (*Trail, error)
	GetTrailByARN(trailARN string) (*Trail, error)
	UpdateTrail(trail *Trail) error
	DeleteTrail(trailName string) error
	ResolveTrail(nameOrARN string) (*Trail, error)
	ListTrails(opts common.ListOptions) (*common.ListResult[Trail], error)
	StartLogging(trailName string) error
	StopLogging(trailName string) error
	PutEventSelector(trailName string, eventSelectors []EventSelector) error
	PutAdvancedEventSelectors(trailName string, selectors []AdvancedEventSelector) error
	GetEventSelector(trailName string) ([]EventSelector, error)
	PutInsightSelectors(trailName string, insightSelectors []InsightSelector) error
	GetInsightSelectors(trailName string) ([]InsightSelector, error)
	PutEvent(event *Event) error
	LookupEvents(query EventQuery) ([]*Event, string, error)
	GetEventByID(eventID string) (*Event, error)
	RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP, accessKeyID string, requestParams, responseElements map[string]interface{}, resources []Resource) error
	GetResourcePolicy(resourceARN string) (*ResourcePolicy, error)
	PutResourcePolicy(resourceARN string, policy string) error
	DeleteResourcePolicy(resourceARN string) error
	Tag(trailName string, tags map[string]string) error
	Untag(trailName string, tagKeys []string) error
	ListAsSlice(trailName string) ([]types.Tag, error)
	ListPublicKeys(startTime, endTime *time.Time) ([]*PublicKey, error)
	GenerateAndStorePublicKey(trailName string) (*PublicKey, error)
	DeletePublicKeysByTrail(trailName string) error
	CreateEventDataStore(eds *EventDataStore) (*EventDataStore, error)
	GetEventDataStore(idOrARN string) (*EventDataStore, error)
	ListEventDataStores(opts common.ListOptions) (*common.ListResult[EventDataStore], error)
	UpdateEventDataStore(eds *EventDataStore) error
	DeleteEventDataStore(id string) error
	RestoreEventDataStore(id string) (*EventDataStore, error)
	SaveQuery(qr *QueryRecord) error
	GetQuery(queryID string) (*QueryRecord, error)
	ListQueriesByEDS(edsID string) ([]*QueryRecord, error)
	CreateChannel(ch *Channel) (*Channel, error)
	GetChannel(arn string) (*Channel, error)
	UpdateChannel(ch *Channel) error
	DeleteChannel(arn string) error
	ListChannels(opts common.ListOptions) (*common.ListResult[Channel], error)
	GetEventConfiguration(trailName, edsID string) (map[string]interface{}, error)
	PutEventConfiguration(trailName, edsID string, config map[string]interface{}) error
	RegisterDelegatedAdmin(accountID string) error
	DeregisterDelegatedAdmin(accountID string) error
	CreateImport(imp *Import) (*Import, error)
	GetImport(importID string) (*Import, error)
	UpdateImport(imp *Import) error
	ListImports(opts common.ListOptions, destination, statusFilter string) (*common.ListResult[Import], error)
	ListImportFailures(importID string, opts common.ListOptions) (*common.ListResult[ImportFailure], error)
}

var _ CloudTrailStoreInterface = (*CloudTrailStore)(nil)
