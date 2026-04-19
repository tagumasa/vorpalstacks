package cloudtrail

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
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
	ListTrails() ([]*Trail, error)
	StartLogging(trailName string) error
	StopLogging(trailName string) error
	PutEventSelector(trailName string, eventSelectors []EventSelector) error
	GetEventSelector(trailName string) ([]EventSelector, error)
	PutInsightSelectors(trailName string, insightSelectors []InsightSelector) error
	GetInsightSelectors(trailName string) ([]InsightSelector, error)
	PutEvent(event *Event) error
	LookupEvents(query EventQuery) ([]*Event, string, error)
	GetEventByID(eventID string) (*Event, error)
	RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP string, requestParams, responseElements map[string]interface{}) error
	GetResourcePolicy(resourceARN string) (*ResourcePolicy, error)
	PutResourcePolicy(resourceARN string, policy string) error
	DeleteResourcePolicy(resourceARN string) error
	Tag(trailName string, tags map[string]string) error
	Untag(trailName string, tagKeys []string) error
	ListAsSlice(trailName string) ([]types.Tag, error)
	ListPublicKeys(startTime, endTime *time.Time) ([]*PublicKey, error)
	GenerateAndStorePublicKey(trailName string) (*PublicKey, error)
}

var _ CloudTrailStoreInterface = (*CloudTrailStore)(nil)
