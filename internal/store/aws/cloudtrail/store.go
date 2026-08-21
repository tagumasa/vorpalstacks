// Package cloudtrail provides AWS CloudTrail storage functionality for vorpalstacks.
package cloudtrail

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

// ErrImportNotFound is returned when an import cannot be located.
var ErrImportNotFound = fmt.Errorf("import not found")

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
	eventDataStoreStore *common.BaseStore
	queryStore          *common.BaseStore
	channelStore        *common.BaseStore
	eventConfigStore    *common.BaseStore
	delegatedAdminStore *common.BaseStore
	importStore         *common.BaseStore
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

func eventDataStoreBucketName(region string) string {
	return "cloudtrail-event-data-stores-" + region
}

func queryBucketName(region string) string {
	return "cloudtrail-queries-" + region
}

func channelBucketName(region string) string {
	return "cloudtrail-channels-" + region
}

func eventConfigBucketName(region string) string {
	return "cloudtrail-event-config-" + region
}

func delegatedAdminBucketName(region string) string {
	return "cloudtrail-delegated-admins-" + region
}

func importBucketName(region string) string {
	return "cloudtrail-imports-" + region
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
		eventDataStoreStore: common.NewBaseStore(store.Bucket(eventDataStoreBucketName(region)), "cloudtrail-event-data-stores"),
		queryStore:          common.NewBaseStore(store.Bucket(queryBucketName(region)), "cloudtrail-queries"),
		channelStore:        common.NewBaseStore(store.Bucket(channelBucketName(region)), "cloudtrail-channels"),
		eventConfigStore:    common.NewBaseStore(store.Bucket(eventConfigBucketName(region)), "cloudtrail-event-config"),
		delegatedAdminStore: common.NewBaseStore(store.Bucket(delegatedAdminBucketName(region)), "cloudtrail-delegated-admins"),
		importStore:         common.NewBaseStore(store.Bucket(importBucketName(region)), "cloudtrail-imports"),
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

	trailData, err := proto.Marshal(TrailToProto(trail))
	if err != nil {
		return nil, err
	}

	if s.storage != nil {
		if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			if err := txn.Bucket(trailBucketName(s.region)).Put([]byte(trail.Name), trailData); err != nil {
				return err
			}
			if s.arnIndexStore != nil {
				if err := txn.Bucket(arnIndexBucketName(s.region)).Put([]byte(trail.TrailARN), []byte(trail.Name)); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		if err := s.BaseStore.PutProto(trail.Name, TrailToProto(trail)); err != nil {
			return nil, err
		}
		if s.arnIndexStore != nil {
			if err := s.arnIndexStore.Put(trail.TrailARN, trail.Name); err != nil {
				return nil, err
			}
		}
	}

	if len(trail.Tags) > 0 {
		if err := s.TagStore.Tag(trail.Name, trail.Tags); err != nil {
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

	trails, err := common.ListMatchingProto[*pb.Trail](s.BaseStore, "", func() *pb.Trail { return &pb.Trail{} }, func(t *pb.Trail) bool {
		return s.normalizeARN(t.TrailArn) == normalizedARN
	})
	if err != nil {
		return nil, err
	}
	if len(trails) > 0 {
		return ProtoToTrail(trails[0]), nil
	}
	return nil, ErrTrailNotFound
}

// normalizeARN fills in the account-id slot of trail ARNs that were
// recorded without one; resource parts containing colons (versioned or
// qualified resources) are preserved by the shared splitter.
func (s *CloudTrailStore) normalizeARN(arn string) string {
	partition, service, region, accountID, resource := svcarn.SplitARN(arn)
	if service == "" || accountID != "" {
		return arn
	}
	return "arn:" + partition + ":" + service + ":" + region + ":" + s.accountID + ":" + resource
}

// ResolveTrail resolves a trail by name or ARN.
func (s *CloudTrailStore) ResolveTrail(nameOrARN string) (*Trail, error) {
	_, _, _, _, resource := svcarn.SplitARN(nameOrARN)
	if strings.HasPrefix(resource, "trail/") {
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
		if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			if err := txn.Bucket(trailBucketName(s.region)).Delete([]byte(trailName)); err != nil {
				return err
			}
			if s.arnIndexStore != nil {
				if err := txn.Bucket(arnIndexBucketName(s.region)).Delete([]byte(trail.TrailARN)); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		// Clean up tags using TagStore which handles both the main bucket
		// (resourceKey\x00tagKey entries) and the inverted index bucket.
		// This cannot be done inside the transaction above because the
		// tag entry keys use a \x00 separator and a blind Delete(trailName)
		// would not match any entries.
		return s.TagStore.Delete(trailName)
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

// ListTrails returns CloudTrail trails with pagination support.
func (s *CloudTrailStore) ListTrails(opts common.ListOptions) (*common.ListResult[Trail], error) {
	result, err := common.ListProto[*pb.Trail](s.BaseStore, opts, func() *pb.Trail { return &pb.Trail{} }, nil)
	if err != nil {
		return nil, err
	}
	var trails []*Trail
	for _, t := range result.Items {
		trails = append(trails, ProtoToTrail(t))
	}
	return &common.ListResult[Trail]{
		Items:       trails,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

// StartLogging starts logging for a CloudTrail trail.
func (s *CloudTrailStore) StartLogging(trailName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trail, err := s.ResolveTrail(trailName)
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

	trail, err := s.ResolveTrail(trailName)
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
	trail.AdvancedEventSelectors = nil
	trail.HasCustomEventSelectors = true

	return s.updateTrailInternal(trail)
}

// PutAdvancedEventSelectors sets advanced event selectors for a CloudTrail trail.
// Providing advanced selectors clears basic event selectors per AWS spec.
func (s *CloudTrailStore) PutAdvancedEventSelectors(trailName string, selectors []AdvancedEventSelector) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trail, err := s.GetTrail(trailName)
	if err != nil {
		return err
	}

	trail.AdvancedEventSelectors = selectors
	trail.EventSelectors = nil
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
	eventData, err := proto.Marshal(EventToProto(event))
	if err != nil {
		return err
	}

	if s.storage != nil {
		return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			if err := txn.Bucket(eventBucketName(s.region)).Put([]byte(key), eventData); err != nil {
				return err
			}
			if err := txn.Bucket(eventIDIndexBucketName(s.region)).Put([]byte(event.EventID), []byte(key)); err != nil {
				return err
			}
			if s.indexer != nil {
				return s.indexer.AddIndexInTxn(txn, event)
			}
			return nil
		})
	}

	if err := s.eventsStore.PutProto(key, EventToProto(event)); err != nil {
		return err
	}
	if err := s.eventIDIndexStore.Put(event.EventID, key); err != nil {
		return err
	}
	if s.indexer != nil {
		return s.indexer.AddIndex(event)
	}
	return nil
}

// LookupEvents looks up CloudTrail events by query. The indexer paths
// (EventName, Username, EventSource, Time) paginate via an opaque
// IndexCursor encoded into the returned nextToken. The default scan path
// uses a plain marker string. Callers should treat nextToken as opaque
// and pass it back unchanged in subsequent calls.
func (s *CloudTrailStore) LookupEvents(query EventQuery) ([]*Event, string, error) {
	if query.MaxResults <= 0 {
		query.MaxResults = 50
	}

	// Decode the incoming nextToken into an IndexCursor for indexer paths.
	// Non-indexed tokens (empty or scan-path markers) yield an empty cursor,
	// causing the query to start from the beginning.
	cursor, cursorErr := decodeIndexCursor(query.NextToken)
	if cursorErr != nil {
		cursor = IndexCursor{}
	}

	var eventIDs []string
	var nextCursor IndexCursor
	var err error

	switch {
	case query.EventID != "":
		event, getErr := s.GetEventByID(query.EventID)
		if getErr != nil {
			return nil, "", nil
		}
		return []*Event{event}, "", nil
	case len(query.EventNames) > 0 && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByEventName(query.EventNames, query.MaxResults, cursor)
	case query.Username != "" && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByUsername(query.Username, query.MaxResults, cursor)
	case query.EventSource != "" && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByEventSource(query.EventSource, query.MaxResults, cursor)
	case (query.StartTime != nil || query.EndTime != nil) && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByTime(query.StartTime, query.EndTime, query.MaxResults, cursor)
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

	return events, encodeIndexCursor(nextCursor), nil
}

func (s *CloudTrailStore) eventIDIndexBucket() storage.Bucket {
	if s.storage != nil {
		return s.storage.(storage.BasicStorage).Bucket(eventIDIndexBucketName(s.region))
	}
	return nil
}

// GetEventByID retrieves a CloudTrail event by ID.
func (s *CloudTrailStore) GetEventByID(eventID string) (*Event, error) {
	var fullKey string
	if bucket := s.eventIDIndexBucket(); bucket != nil {
		fullKeyBytes, err := bucket.Get([]byte(eventID))
		if err != nil || fullKeyBytes == nil {
			return nil, ErrEventNotFound
		}
		fullKey = string(fullKeyBytes)
	} else if err := s.eventIDIndexStore.Get(eventID, &fullKey); err != nil {
		return nil, ErrEventNotFound
	}

	var p pb.Event
	if err := s.eventsStore.GetProto(fullKey, &p); err != nil {
		return nil, ErrEventNotFound
	}
	return ProtoToEvent(&p), nil
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

	if query.EventSource != "" && event.GetEventSource() != query.EventSource {
		return false
	}

	if query.AccessKeyID != "" && event.GetAccessKeyId() != query.AccessKeyID {
		return false
	}

	if query.EventID != "" && event.GetEventId() != query.EventID {
		return false
	}

	if query.ReadOnly == "true" && event.GetReadOnly() != "true" {
		return false
	}
	if query.ReadOnly == "false" && event.GetReadOnly() != "false" {
		return false
	}

	if query.EventCategory != "" && event.GetEventCategory() != query.EventCategory {
		return false
	}

	return true
}

func (s *CloudTrailStore) eventMatchesQuery(event *Event, query EventQuery) bool {
	return protoMatchesQuery(EventToProto(event), query)
}

// RecordServiceEvent records a service event to CloudTrail.
func (s *CloudTrailStore) RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP, accessKeyID string, requestParams, responseElements map[string]interface{}, resources []Resource) error {
	event := NewEvent(eventName, eventSource, userIdentity)
	event.RequestParameters = requestParams
	event.ResponseElements = responseElements
	event.SourceIPAddress = sourceIP
	event.AccessKeyId = accessKeyID
	event.UserAgent = "vorpalstacks-internal"
	for _, r := range resources {
		event.Resources = append(event.Resources, Resource{ResourceType: r.ResourceType, ResourceName: r.ResourceName})
	}
	event.generateCloudTrailEvent()
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
	EventSource   string
	AccessKeyID   string
	EventID       string
	ReadOnly      string
	EventCategory string
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

// DeletePublicKeysByTrail removes all public keys associated with the
// given trail name.  This is called during trail deletion to prevent
// orphaned key material from lingering after the trail is gone.
func (s *CloudTrailStore) DeletePublicKeysByTrail(trailName string) error {
	var toDelete []string
	err := s.publicKeyStore.ForEach(func(key string, value []byte) error {
		var p pb.PublicKey
		if err := proto.Unmarshal(value, &p); err != nil {
			return err
		}
		if p.TrailName == trailName {
			toDelete = append(toDelete, key)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, key := range toDelete {
		if err := s.publicKeyStore.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// --- Event Data Store operations ---

// CreateEventDataStore persists a new event data store.
func (s *CloudTrailStore) CreateEventDataStore(eds *EventDataStore) (*EventDataStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.listEventDataStoresRaw()
	if err != nil {
		return nil, err
	}
	for _, e := range existing {
		if e.Name == eds.Name && e.Status != "PENDING_DELETION" {
			return nil, ErrEventDataStoreAlreadyExists
		}
	}

	eds.CreatedTimestamp = time.Now().UTC()
	eds.UpdatedTimestamp = eds.CreatedTimestamp

	if err := s.eventDataStoreStore.Put(eds.EventDataStoreID, eds); err != nil {
		return nil, fmt.Errorf("failed to store event data store: %w", err)
	}

	return eds, nil
}

// GetEventDataStore retrieves an event data store by ID or ARN.
func (s *CloudTrailStore) GetEventDataStore(idOrARN string) (*EventDataStore, error) {
	id := extractEventDataStoreID(idOrARN)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		return nil, ErrEventDataStoreNotFound
	}
	return &eds, nil
}

// ListEventDataStores returns all event data stores.
func (s *CloudTrailStore) ListEventDataStores(opts common.ListOptions) (*common.ListResult[EventDataStore], error) {
	return common.List[EventDataStore](s.eventDataStoreStore, opts, nil)
}

// UpdateEventDataStore updates an existing event data store.
func (s *CloudTrailStore) UpdateEventDataStore(eds *EventDataStore) error {
	eds.UpdatedTimestamp = time.Now().UTC()
	return s.eventDataStoreStore.Put(eds.EventDataStoreID, eds)
}

// DeleteEventDataStore soft-deletes an event data store by setting its status
// to PENDING_DELETION. AWS spec: after 7 days it is permanently deleted.
func (s *CloudTrailStore) DeleteEventDataStore(id string) error {
	id = extractEventDataStoreID(id)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		return ErrEventDataStoreNotFound
	}
	eds.Status = "PENDING_DELETION"
	now := time.Now().UTC()
	eds.DeletedTimestamp = &now
	return s.eventDataStoreStore.Put(eds.EventDataStoreID, &eds)
}

// RestoreEventDataStore restores a PENDING_DELETION event data store to
// ENABLED status.
func (s *CloudTrailStore) RestoreEventDataStore(id string) (*EventDataStore, error) {
	id = extractEventDataStoreID(id)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		return nil, ErrEventDataStoreNotFound
	}
	if eds.Status != "PENDING_DELETION" {
		return nil, ErrEventDataStoreNotPendingDeletion
	}
	eds.Status = "ENABLED"
	eds.DeletedTimestamp = nil
	eds.UpdatedTimestamp = time.Now().UTC()
	if err := s.eventDataStoreStore.Put(eds.EventDataStoreID, &eds); err != nil {
		return nil, err
	}
	return &eds, nil
}

// listEventDataStoresRaw returns all stored event data stores.
func (s *CloudTrailStore) listEventDataStoresRaw() ([]*EventDataStore, error) {
	var result []*EventDataStore
	err := s.eventDataStoreStore.ForEach(func(_ string, value []byte) error {
		var eds EventDataStore
		if err := json.Unmarshal(value, &eds); err != nil {
			return nil
		}
		result = append(result, &eds)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list event data stores: %w", err)
	}
	return result, nil
}

// extractEventDataStoreID extracts the UUID from an ID or ARN.
func extractEventDataStoreID(idOrARN string) string {
	if idx := strings.LastIndex(idOrARN, "/"); idx >= 0 {
		return idOrARN[idx+1:]
	}
	return idOrARN
}

// --- Query operations ---

// SaveQuery persists a query record.
func (s *CloudTrailStore) SaveQuery(qr *QueryRecord) error {
	return s.queryStore.Put(qr.QueryID, qr)
}

// GetQuery retrieves a query record by ID.
func (s *CloudTrailStore) GetQuery(queryID string) (*QueryRecord, error) {
	var qr QueryRecord
	if err := s.queryStore.Get(queryID, &qr); err != nil {
		return nil, ErrQueryNotFound
	}
	return &qr, nil
}

// ListQueriesByEDS lists queries for an event data store.
func (s *CloudTrailStore) ListQueriesByEDS(edsID string) ([]*QueryRecord, error) {
	var result []*QueryRecord
	err := s.queryStore.ForEach(func(_ string, value []byte) error {
		var qr QueryRecord
		if err := json.Unmarshal(value, &qr); err != nil {
			return nil
		}
		if extractEventDataStoreID(qr.EventDataStore) == edsID {
			result = append(result, &qr)
		}
		return nil
	})
	return result, err
}

// --- Channel operations ---

// CreateChannel persists a new channel.
func (s *CloudTrailStore) CreateChannel(ch *Channel) (*Channel, error) {
	ch.CreatedAt = time.Now().UTC()
	ch.UpdatedAt = ch.CreatedAt
	return ch, s.channelStore.Put(ch.ChannelARN, ch)
}

// GetChannel retrieves a channel by ARN.
func (s *CloudTrailStore) GetChannel(arn string) (*Channel, error) {
	var ch Channel
	if err := s.channelStore.Get(arn, &ch); err != nil {
		return nil, ErrChannelNotFound
	}
	return &ch, nil
}

// UpdateChannel updates an existing channel.
func (s *CloudTrailStore) UpdateChannel(ch *Channel) error {
	ch.UpdatedAt = time.Now().UTC()
	return s.channelStore.Put(ch.ChannelARN, ch)
}

// DeleteChannel deletes a channel by ARN.
func (s *CloudTrailStore) DeleteChannel(arn string) error {
	return s.channelStore.Delete(arn)
}

// ListChannels lists channels with pagination support.
func (s *CloudTrailStore) ListChannels(opts common.ListOptions) (*common.ListResult[Channel], error) {
	return common.List[Channel](s.channelStore, opts, nil)
}

// --- Event Configuration ---

// GetEventConfiguration retrieves event configuration for a trail or EDS.
func (s *CloudTrailStore) GetEventConfiguration(trailName, edsID string) (map[string]interface{}, error) {
	key := eventConfigKey(trailName, edsID)
	var config map[string]interface{}
	if err := s.eventConfigStore.Get(key, &config); err != nil {
		return nil, fmt.Errorf("event configuration not found")
	}
	return config, nil
}

// PutEventConfiguration stores event configuration.
func (s *CloudTrailStore) PutEventConfiguration(trailName, edsID string, config map[string]interface{}) error {
	key := eventConfigKey(trailName, edsID)
	return s.eventConfigStore.Put(key, config)
}

func eventConfigKey(trailName, edsID string) string {
	if trailName != "" {
		return "trail:" + trailName
	}
	return "eds:" + edsID
}

// --- Delegated Admin ---

// RegisterDelegatedAdmin registers a delegated admin account.
func (s *CloudTrailStore) RegisterDelegatedAdmin(accountID string) error {
	return s.delegatedAdminStore.Put(accountID, true)
}

// DeregisterDelegatedAdmin removes a delegated admin registration.
func (s *CloudTrailStore) DeregisterDelegatedAdmin(accountID string) error {
	return s.delegatedAdminStore.Delete(accountID)
}

// IsDelegatedAdmin checks if an account is a registered delegated admin.
func (s *CloudTrailStore) IsDelegatedAdmin(accountID string) bool {
	return s.delegatedAdminStore.Exists(accountID)
}

// --- Import operations ---

// CreateImport persists a new import record.
func (s *CloudTrailStore) CreateImport(imp *Import) (*Import, error) {
	return imp, s.importStore.Put(imp.ImportID, imp)
}

// GetImport retrieves an import by ID.
func (s *CloudTrailStore) GetImport(importID string) (*Import, error) {
	var imp Import
	if err := s.importStore.Get(importID, &imp); err != nil {
		return nil, ErrImportNotFound
	}
	return &imp, nil
}

// UpdateImport updates an existing import record.
func (s *CloudTrailStore) UpdateImport(imp *Import) error {
	imp.UpdatedTimestamp = time.Now().UTC()
	return s.importStore.Put(imp.ImportID, imp)
}

// ListImports lists imports with optional destination and status filters.
func (s *CloudTrailStore) ListImports(opts common.ListOptions, destination, statusFilter string) (*common.ListResult[Import], error) {
	return common.List[Import](s.importStore, opts, func(imp *Import) bool {
		if destination != "" {
			found := false
			for _, d := range imp.Destinations {
				if d == destination {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if statusFilter != "" && imp.ImportStatus != statusFilter {
			return false
		}
		return true
	})
}

// ListImportFailures lists failures for a specific import.
func (s *CloudTrailStore) ListImportFailures(importID string, opts common.ListOptions) (*common.ListResult[ImportFailure], error) {
	imp, err := s.GetImport(importID)
	if err != nil {
		return nil, err
	}

	total := len(imp.Failures)
	offset := 0
	if opts.Marker != "" {
		if n, err := strconv.Atoi(opts.Marker); err == nil && n >= 0 {
			offset = n
		}
	}
	maxItems := opts.MaxItems
	if maxItems <= 0 {
		maxItems = 50
	}

	end := offset + maxItems
	if end > total {
		end = total
	}

	var items []*ImportFailure
	for i := offset; i < end; i++ {
		items = append(items, &imp.Failures[i])
	}

	result := &common.ListResult[ImportFailure]{
		Items: items,
	}
	if end < total {
		result.NextMarker = strconv.Itoa(end)
		result.IsTruncated = true
	}
	return result, nil
}
