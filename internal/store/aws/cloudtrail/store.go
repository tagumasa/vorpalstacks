// Package cloudtrail provides AWS CloudTrail storage functionality for vorpalstacks.
package cloudtrail

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_cloudtrail"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// ErrResourcePolicyNotFound is returned when a resource policy cannot be located.
var ErrResourcePolicyNotFound = fmt.Errorf("resource policy not found")

// CloudTrailStore provides CloudTrail storage operations.
type CloudTrailStore struct {
	*common.BaseStore
	eventsStore       *common.BaseStore
	eventIDIndexStore *common.BaseStore
	*common.TagStore
	arnBuilder          *svcarn.ARNBuilder
	accountID           string
	region              string
	mu                  sync.Mutex
	indexer             *EventIndexManager
	arnIndexStore       *common.BaseStore
	resourcePolicyStore *common.BaseStore
	publicKeyStore      *common.BaseStore
	storage             storage.TransactionalStorageWith2PC
}

func trailBucketName(region string) string {
	return "cloudtrail-trails-" + region
}

func eventBucketName(region string) string {
	return "cloudtrail-events-" + region
}

func arnIndexBucketName(region string) string {
	return "cloudtrail-arn-index-" + region
}

func resourcePolicyBucketName(region string) string {
	return "cloudtrail-resource-policy-" + region
}

func publicKeyBucketName(region string) string {
	return "cloudtrail-public-keys-" + region
}

func eventIDIndexBucketName(region string) string {
	return "cloudtrail-event-id-index-" + region
}

// NewCloudTrailStore creates a new CloudTrail store.
func NewCloudTrailStore(store storage.BasicStorage, accountID, region string) *CloudTrailStore {
	var tstore storage.TransactionalStorageWith2PC
	if ts, ok := store.(storage.TransactionalStorageWith2PC); ok {
		tstore = ts
	}

	return &CloudTrailStore{
		BaseStore:           common.NewBaseStore(store.Bucket(trailBucketName(region)), "cloudtrail-trails"),
		eventsStore:         common.NewBaseStore(store.Bucket(eventBucketName(region)), "cloudtrail-events"),
		eventIDIndexStore:   common.NewBaseStore(store.Bucket(eventIDIndexBucketName(region)), "cloudtrail-event-id-index"),
		TagStore:            common.NewTagStoreWithRegion(store, "cloudtrail", region),
		arnBuilder:          svcarn.NewARNBuilder(accountID, region),
		accountID:           accountID,
		region:              region,
		indexer:             NewEventIndexManager(store, accountID, region),
		arnIndexStore:       common.NewBaseStore(store.Bucket(arnIndexBucketName(region)), "cloudtrail-arn-index"),
		resourcePolicyStore: common.NewBaseStore(store.Bucket(resourcePolicyBucketName(region)), "cloudtrail-resource-policy"),
		publicKeyStore:      common.NewBaseStore(store.Bucket(publicKeyBucketName(region)), "cloudtrail-public-keys"),
		storage:             tstore,
	}
}

// GetAccountID returns the AWS account ID.
func (s *CloudTrailStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the AWS region.
func (s *CloudTrailStore) GetRegion() string {
	return s.region
}

// BuildTrailARN builds the ARN for a CloudTrail trail.
func (s *CloudTrailStore) BuildTrailARN(trailName string) string {
	return s.arnBuilder.CloudTrail().Trail(trailName)
}

// CreateTrail creates a new CloudTrail trail.
func (s *CloudTrailStore) CreateTrail(trail *Trail) (*Trail, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if trail.Name == "" {
		return nil, ErrInvalidTrailName
	}

	if s.Exists(trail.Name) {
		return nil, ErrTrailAlreadyExists
	}

	now := time.Now().UTC()
	trail.TrailARN = s.BuildTrailARN(trail.Name)
	trail.HomeRegion = s.region
	trail.CreatedAt = now
	trail.LastUpdated = now
	trail.IsLogging = false

	if trail.Tags == nil {
		trail.Tags = make(map[string]string)
	}

	if s.storage != nil {
		err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			trailData, err := proto.Marshal(TrailToProto(trail))
			if err != nil {
				return err
			}

			trailsBucket := txn.Bucket(trailBucketName(s.region))
			if err := trailsBucket.Put([]byte(trail.Name), trailData); err != nil {
				return err
			}

			if s.arnIndexStore != nil {
				arnIndexBucket := txn.Bucket(arnIndexBucketName(s.region))
				if err := arnIndexBucket.Put([]byte(trail.TrailARN), []byte(trail.Name)); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		if len(trail.Tags) > 0 {
			if err := s.TagStore.TagResource(trail.Name, trail.Tags); err != nil {
				return nil, err
			}
		}

		if trail.LogFileValidationEnabled {
			if _, err := s.GenerateAndStorePublicKey(trail.Name); err != nil {
				return nil, fmt.Errorf("failed to generate public key for trail: %w", err)
			}
		}

		return trail, nil
	}

	if err := s.PutProto(trail.Name, TrailToProto(trail)); err != nil {
		return nil, err
	}

	if s.arnIndexStore != nil {
		if err := s.arnIndexStore.Put(trail.TrailARN, trail.Name); err != nil {
			return nil, err
		}
	}

	if len(trail.Tags) > 0 {
		if err := s.TagStore.TagResource(trail.Name, trail.Tags); err != nil {
			return nil, err
		}
	}

	if trail.LogFileValidationEnabled {
		if _, err := s.GenerateAndStorePublicKey(trail.Name); err != nil {
			return nil, fmt.Errorf("failed to generate public key for trail: %w", err)
		}
	}

	return trail, nil
}

// GetTrail retrieves a CloudTrail trail by name.
func (s *CloudTrailStore) GetTrail(trailName string) (*Trail, error) {
	var p pb.Trail
	if err := s.BaseStore.GetProto(trailName, &p); err != nil {
		return nil, ErrTrailNotFound
	}
	return ProtoToTrail(&p), nil
}

// GetTrailByARN retrieves a CloudTrail trail by ARN.
func (s *CloudTrailStore) GetTrailByARN(trailARN string) (*Trail, error) {
	normalizedARN := s.normalizeARN(trailARN)

	if s.arnIndexStore != nil && s.arnIndexStore.Exists(normalizedARN) {
		var trailName string
		if err := s.arnIndexStore.Get(normalizedARN, &trailName); err == nil {
			return s.GetTrail(trailName)
		}
	}

	trails, err := s.ListTrails()
	if err != nil {
		return nil, err
	}
	for _, trail := range trails {
		if s.normalizeARN(trail.TrailARN) == normalizedARN {
			return trail, nil
		}
	}
	return nil, ErrTrailNotFound
}

func (s *CloudTrailStore) normalizeARN(arn string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) >= 5 && parts[4] == "" {
		parts[4] = s.accountID
		return strings.Join(parts, ":")
	}
	return arn
}

func (s *CloudTrailStore) resolveTrail(nameOrARN string) (*Trail, error) {
	if strings.Contains(nameOrARN, ":trail/") {
		return s.GetTrailByARN(nameOrARN)
	}
	return s.GetTrail(nameOrARN)
}

// UpdateTrail updates an existing CloudTrail trail.
func (s *CloudTrailStore) UpdateTrail(trail *Trail) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateTrailInternal(trail)
}

func (s *CloudTrailStore) updateTrailInternal(trail *Trail) error {
	if !s.Exists(trail.Name) {
		return ErrTrailNotFound
	}
	trail.LastUpdated = time.Now().UTC()
	return s.PutProto(trail.Name, TrailToProto(trail))
}

// DeleteTrail deletes a CloudTrail trail by name.
func (s *CloudTrailStore) DeleteTrail(trailName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(trailName) {
		return ErrTrailNotFound
	}

	trail, err := s.GetTrail(trailName)
	if err != nil {
		return err
	}

	if s.storage != nil {
		return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			trailsBucket := txn.Bucket(trailBucketName(s.region))
			if err := trailsBucket.Delete([]byte(trailName)); err != nil {
				return err
			}

			if s.arnIndexStore != nil {
				arnIndexBucket := txn.Bucket(arnIndexBucketName(s.region))
				if err := arnIndexBucket.Delete([]byte(trail.TrailARN)); err != nil {
					return err
				}
			}

			tagBucketName := "cloudtrail-tags-" + s.region
			tagBucket := txn.Bucket(tagBucketName)
			if err := tagBucket.Delete([]byte(trailName)); err != nil {
				return err
			}

			return nil
		})
	}

	if s.arnIndexStore != nil {
		if err := s.arnIndexStore.Delete(trail.TrailARN); err != nil {
			return err
		}
	}

	if err := s.TagStore.Delete(trailName); err != nil {
		return err
	}

	return s.BaseStore.Delete(trailName)
}

// ListTrails returns all CloudTrail trails.
func (s *CloudTrailStore) ListTrails() ([]*Trail, error) {
	opts := common.ListOptions{MaxItems: 10000}
	result, err := common.ListProto[*pb.Trail](s.BaseStore, opts, func() *pb.Trail { return &pb.Trail{} }, nil)
	if err != nil {
		return nil, err
	}
	var trails []*Trail
	for _, t := range result.Items {
		trails = append(trails, ProtoToTrail(t))
	}
	return trails, nil
}

// StartLogging starts logging for a CloudTrail trail.
func (s *CloudTrailStore) StartLogging(trailName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trail, err := s.resolveTrail(trailName)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	trail.IsLogging = true
	trail.StartedLoggingAt = &now
	trail.StoppedLoggingAt = nil

	return s.updateTrailInternal(trail)
}

// StopLogging stops logging for a CloudTrail trail.
func (s *CloudTrailStore) StopLogging(trailName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trail, err := s.resolveTrail(trailName)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	trail.IsLogging = false
	trail.StoppedLoggingAt = &now

	return s.updateTrailInternal(trail)
}

// PutEventSelector sets event selectors for a CloudTrail trail.
func (s *CloudTrailStore) PutEventSelector(trailName string, eventSelectors []EventSelector) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trail, err := s.GetTrail(trailName)
	if err != nil {
		return err
	}

	trail.EventSelectors = eventSelectors
	trail.HasCustomEventSelectors = true

	return s.updateTrailInternal(trail)
}

// GetEventSelector retrieves event selectors for a CloudTrail trail.
func (s *CloudTrailStore) GetEventSelector(trailName string) ([]EventSelector, error) {
	trail, err := s.GetTrail(trailName)
	if err != nil {
		return nil, err
	}
	return trail.EventSelectors, nil
}

// PutInsightSelectors sets insight selectors for a CloudTrail trail.
func (s *CloudTrailStore) PutInsightSelectors(trailName string, insightSelectors []InsightSelector) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trail, err := s.GetTrail(trailName)
	if err != nil {
		return err
	}

	trail.InsightSelectors = insightSelectors
	trail.HasInsightSelectors = true

	return s.updateTrailInternal(trail)
}

// GetInsightSelectors retrieves insight selectors for a CloudTrail trail.
func (s *CloudTrailStore) GetInsightSelectors(trailName string) ([]InsightSelector, error) {
	trail, err := s.GetTrail(trailName)
	if err != nil {
		return nil, err
	}
	return trail.InsightSelectors, nil
}

// PutEvent stores a CloudTrail event.
func (s *CloudTrailStore) PutEvent(event *Event) error {
	if event.EventID == "" {
		event.EventID = uuid.New().String()
	}
	if event.EventTime.IsZero() {
		event.EventTime = time.Now().UTC()
	}

	key := fmt.Sprintf("%d#%s", event.EventTime.UnixNano(), event.EventID)

	if s.storage != nil {
		return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			eventData, err := proto.Marshal(EventToProto(event))
			if err != nil {
				return err
			}

			eventsBucket := txn.Bucket(eventBucketName(s.region))
			if err := eventsBucket.Put([]byte(key), eventData); err != nil {
				return err
			}

			eventIDIndexBucket := txn.Bucket(eventIDIndexBucketName(s.region))
			if err := eventIDIndexBucket.Put([]byte(event.EventID), []byte(key)); err != nil {
				return err
			}

			if s.indexer != nil {
				if err := s.indexer.AddIndexInTxn(txn, event); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := s.eventsStore.Put(key, event); err != nil {
		return err
	}

	if err := s.eventIDIndexStore.Put(event.EventID, key); err != nil {
		return err
	}

	if s.indexer != nil {
		if err := s.indexer.AddIndex(event); err != nil {
			return err
		}
	}

	return nil
}

// LookupEvents looks up CloudTrail events by query.
func (s *CloudTrailStore) LookupEvents(query EventQuery) ([]*Event, string, error) {
	if query.MaxResults <= 0 {
		query.MaxResults = 50
	}
	var eventIDs []string
	var err error

	switch {
	case len(query.EventNames) > 0 && s.indexer != nil:
		eventIDs, err = s.indexer.QueryByEventName(query.EventNames, query.MaxResults)
	case query.Username != "" && s.indexer != nil:
		eventIDs, err = s.indexer.QueryByUsername(query.Username, query.MaxResults)
	case (query.StartTime != nil || query.EndTime != nil) && s.indexer != nil:
		eventIDs, err = s.indexer.QueryByTime(query.StartTime, query.EndTime, query.MaxResults)
	default:
		return s.lookupEventsScan(query)
	}

	if err != nil {
		return nil, "", err
	}

	var events []*Event
	for _, id := range eventIDs {
		if int32(len(events)) >= query.MaxResults {
			break
		}
		event, err := s.GetEventByID(id)
		if err != nil {
			continue
		}
		if s.eventMatchesQuery(event, query) {
			events = append(events, event)
		}
	}

	return events, "", nil
}

// GetEventByID retrieves a CloudTrail event by ID.
func (s *CloudTrailStore) GetEventByID(eventID string) (*Event, error) {
	var fullKey string
	if err := s.eventIDIndexStore.Get(eventID, &fullKey); err != nil {
		return nil, ErrEventNotFound
	}

	var event Event
	if err := s.eventsStore.Get(fullKey, &event); err != nil {
		return nil, ErrEventNotFound
	}
	return &event, nil
}

func (s *CloudTrailStore) lookupEventsScan(query EventQuery) ([]*Event, string, error) {
	opts := common.ListOptions{
		Marker:   query.NextToken,
		MaxItems: int(query.MaxResults),
	}
	result, err := common.ListProto[*pb.Event](s.eventsStore, opts, func() *pb.Event { return &pb.Event{} }, func(e *pb.Event) bool {
		return protoMatchesQuery(e, query)
	})
	if err != nil {
		return nil, "", err
	}
	events := make([]*Event, len(result.Items))
	for i, p := range result.Items {
		events[i] = ProtoToEvent(p)
	}
	return events, result.NextMarker, nil
}

func protoMatchesQuery(event *pb.Event, query EventQuery) bool {
	eventTime := time.UnixMilli(event.GetEventTime())
	if query.StartTime != nil && eventTime.Before(*query.StartTime) {
		return false
	}
	if query.EndTime != nil && eventTime.After(*query.EndTime) {
		return false
	}

	if len(query.EventNames) > 0 {
		found := false
		for _, name := range query.EventNames {
			if event.GetEventName() == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if query.Username != "" {
		if event.GetUserIdentity() == nil || event.GetUserIdentity().GetUserName() != query.Username {
			return false
		}
	}

	if len(query.ResourceNames) > 0 {
		if len(event.GetResources()) == 0 {
			return false
		}
		found := false
		for _, rn := range query.ResourceNames {
			for _, res := range event.GetResources() {
				if res.GetResourceName() == rn {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	if query.ResourceType != "" {
		if len(event.GetResources()) == 0 {
			return false
		}
		found := false
		for _, res := range event.GetResources() {
			if res.GetResourceType() == query.ResourceType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (s *CloudTrailStore) eventMatchesQuery(event *Event, query EventQuery) bool {
	if query.StartTime != nil && event.EventTime.Before(*query.StartTime) {
		return false
	}
	if query.EndTime != nil && event.EventTime.After(*query.EndTime) {
		return false
	}

	if len(query.EventNames) > 0 {
		found := false
		for _, name := range query.EventNames {
			if event.EventName == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if query.Username != "" {
		if event.UserIdentity == nil || event.UserIdentity.UserName != query.Username {
			return false
		}
	}

	if len(query.ResourceNames) > 0 {
		if event.Resources == nil {
			return false
		}
		found := false
		for _, rn := range query.ResourceNames {
			for _, res := range event.Resources {
				if res.ResourceName == rn {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	if query.ResourceType != "" {
		if event.Resources == nil {
			return false
		}
		found := false
		for _, res := range event.Resources {
			if res.ResourceType == query.ResourceType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// RecordServiceEvent records a service event to CloudTrail.
func (s *CloudTrailStore) RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP string, requestParams, responseElements map[string]interface{}) error {
	event := NewEvent(eventName, eventSource, userIdentity)
	event.RequestParameters = requestParams
	event.ResponseElements = responseElements
	event.SourceIPAddress = sourceIP
	event.UserAgent = "vorpalstacks-internal"
	return s.PutEvent(event)
}

// EventQuery represents a query for looking up CloudTrail events.
type EventQuery struct {
	StartTime     *time.Time
	EndTime       *time.Time
	EventNames    []string
	Username      string
	ResourceNames []string
	ResourceType  string
	MaxResults    int32
	NextToken     string
}

// NewEventQuery creates a new CloudTrail event query with default values.
func NewEventQuery() EventQuery {
	return EventQuery{
		MaxResults: 50,
	}
}

// GetResourcePolicy retrieves a resource policy for CloudTrail
func (s *CloudTrailStore) GetResourcePolicy(resourceARN string) (*ResourcePolicy, error) {
	var p pb.ResourcePolicy
	if err := s.resourcePolicyStore.GetProto(resourceARN, &p); err != nil {
		return nil, ErrResourcePolicyNotFound
	}
	return ProtoToResourcePolicy(&p), nil
}

// PutResourcePolicy stores a resource policy for CloudTrail.
func (s *CloudTrailStore) PutResourcePolicy(resourceARN string, policy string) error {
	rp := &ResourcePolicy{
		ResourceARN: resourceARN,
		Policy:      policy,
	}
	return s.resourcePolicyStore.PutProto(resourceARN, ResourcePolicyToProto(rp))
}

// DeleteResourcePolicy deletes a resource policy for CloudTrail.
func (s *CloudTrailStore) DeleteResourcePolicy(resourceARN string) error {
	return s.resourcePolicyStore.Delete(resourceARN)
}

// StorePublicKey persists a public key for log file validation.
func (s *CloudTrailStore) StorePublicKey(pk *PublicKey) error {
	return s.publicKeyStore.PutProto(pk.PublicKeyID, PublicKeyToProto(pk))
}

// ListPublicKeys returns all stored public keys, optionally filtered by time range.
func (s *CloudTrailStore) ListPublicKeys(startTime, endTime *time.Time) ([]*PublicKey, error) {
	var keys []*PublicKey
	err := s.publicKeyStore.ForEach(func(key string, value []byte) error {
		var p pb.PublicKey
		if err := proto.Unmarshal(value, &p); err != nil {
			return err
		}
		pk := ProtoToPublicKey(&p)
		if startTime != nil && pk.ValidityEndTime.Before(*startTime) {
			return nil
		}
		if endTime != nil && pk.ValidityStartTime.After(*endTime) {
			return nil
		}
		keys = append(keys, pk)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// GenerateAndStorePublicKey creates a new RSA key pair and stores the public key.
func (s *CloudTrailStore) GenerateAndStorePublicKey(trailName string) (*PublicKey, error) {
	pk, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	pk.TrailName = trailName
	if err := s.StorePublicKey(pk); err != nil {
		return nil, err
	}
	return pk, nil
}
