package cloudtrail

import (
	"crypto/rsa"
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
	MutateTrail(trailName string, apply func(*Trail) error) (*Trail, error)
	DeleteTrail(trailName string) error
	ResolveTrail(nameOrARN string) (*Trail, error)
	ListTrails(opts common.ListOptions) (*common.ListResult[Trail], error)
	StartLogging(trailName string) error
	StopLogging(trailName string) error
	PutEventSelector(trailName string, eventSelectors []EventSelector) error
	PutAdvancedEventSelectors(trailName string, selectors []AdvancedEventSelector) error
	PutInsightSelectors(trailName string, insightSelectors []InsightSelector) error
	PutEvent(event *Event) error
	PurgeEventHistoryBefore(cutoff time.Time, batchSize int) (*PurgeResult, error)
	LookupEvents(query EventQuery) ([]*Event, string, error)
	GetEventByID(eventID string) (*Event, error)
	RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP, accessKeyID, userAgent string, readOnly bool, errorCode, errorMessage string, requestParams, responseElements map[string]interface{}, resources []Resource) error
	GetResourcePolicy(resourceARN string) (*ResourcePolicy, error)
	PutResourcePolicy(resourceARN string, policy string) error
	DeleteResourcePolicy(resourceARN string) error
	Tag(trailName string, tags map[string]string) error
	Untag(trailName string, tagKeys []string) error
	ListAsSlice(trailName string) ([]types.Tag, error)
	ListPublicKeys(startTime, endTime *time.Time) ([]*PublicKey, error)
	GenerateAndStorePublicKey(trailName string) (*PublicKey, error)
	LoadTrailSigningKey(trailName string) (*PublicKey, *rsa.PrivateKey, error)
	CreateAndStoreSigningKey() (*PublicKey, *rsa.PrivateKey, error)
	DeletePublicKeysByTrail(trailName string) error
	CreateEventDataStore(eds *EventDataStore) (*EventDataStore, error)
	GetEventDataStore(idOrARN string) (*EventDataStore, error)
	ListEventDataStores(opts common.ListOptions) (*common.ListResult[EventDataStore], error)
	MutateEventDataStore(idOrARN string, apply func(*EventDataStore) error) (*EventDataStore, error)
	RestoreEventDataStore(id string) (*EventDataStore, error)
	DeleteEventDataStoreIf(idOrARN string, guard func(*EventDataStore) error) (bool, error)
	ListEventDataStoresAll() ([]*EventDataStore, error)
	LookupEDSEvents(edsID string, query EDSQuery) ([]*Event, string, error)
	PutEventIntoEDS(edsID string, event *Event) error
	PutEventIntoEDSs(edsIDs []string, event *Event) error
	PurgeEDSEventsBefore(edsID string, cutoff time.Time, batchSize int) (int, error)
	DropEDSEvents(edsID string) error
	SaveQuery(qr *QueryRecord) error
	AdmitQuery(qr *QueryRecord, maxRunning int) error
	GetQuery(queryID string) (*QueryRecord, error)
	MutateQuery(queryID string, apply func(*QueryRecord) error) (*QueryRecord, error)
	ListQueriesByEDS(edsID string) ([]*QueryRecord, error)
	PurgeQueriesBefore(cutoff time.Time) (int, error)
	CreateChannel(ch *Channel) (*Channel, error)
	GetChannel(arn string) (*Channel, error)
	MutateChannel(arn string, apply func(*Channel) error) (*Channel, error)
	DeleteChannelIf(arn string, guard func(*Channel) error) error
	ListChannels(opts common.ListOptions) (*common.ListResult[Channel], error)
	GetEventConfiguration(trailName, edsID string) (map[string]interface{}, error)
	PutEventConfiguration(trailName, edsID string, config map[string]interface{}) error
	DeleteEventConfiguration(trailName, edsID string) error
	CreateImport(imp *Import) (*Import, error)
	CreateImportIfNoOngoing(imp *Import) (*Import, error)
	GetImport(importID string) (*Import, error)
	MutateImport(importID string, apply func(*Import) error) (*Import, error)
	ListImports(opts common.ListOptions, destination, statusFilter string) (*common.ListResult[Import], error)
	ListImportFailures(importID string, opts common.ListOptions) (*common.ListResult[ImportFailure], error)
}

var _ CloudTrailStoreInterface = (*CloudTrailStore)(nil)
