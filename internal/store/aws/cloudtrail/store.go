// Package cloudtrail provides AWS CloudTrail storage functionality for vorpalstacks.
package cloudtrail

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/pagination"
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

func importBucketName(region string) string {
	return "cloudtrail-imports-" + region
}

func eventIDIndexBucketName(region string) string {
	return "cloudtrail-event-id-index-" + region
}

// NewCloudTrailStore creates a new CloudTrail store.
//
// Persistence regimes (recorded decision): trails, events, ARN indexes,
// resource policies and public keys persist through storage_cloudtrail.proto
// (typed schemas with generated converters); the CloudTrail Lake families —
// event data stores, queries, channels, event configurations, imports —
// persist as JSON through BaseStore, so a model-shaped member rides a struct
// change without a proto regeneration. The split stands deliberately:
// unifying either direction would rewrite every family's converters and
// reset the data directory for no behavioural gain. New record families
// persist as JSON unless the event write path's compactness demands binary.
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

	// The trail quota holds at admission: "Trails per Region — 5 ... This
	// quota cannot be increased" (Quotas in AWS CloudTrail). Counting under
	// the creation mutex makes the check-and-write one step.
	trailCount := 0
	if err := s.BaseStore.ForEach(func(_ string, _ []byte) error {
		trailCount++
		return nil
	}); err != nil {
		return nil, err
	}
	if trailCount >= MaxTrailsPerRegion {
		return nil, ErrTrailQuotaExceeded
	}

	now := time.Now().UTC()
	trail.TrailARN = s.BuildTrailARN(trail.Name)
	trail.HomeRegion = s.region
	trail.CreatedAt = now
	trail.LastUpdated = now
	trail.IsLogging = false

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

	if trail.LogFileValidationEnabled {
		if _, err := s.GenerateAndStorePublicKey(trail.Name); err != nil {
			return nil, fmt.Errorf("failed to generate public key for trail: %w", err)
		}
	}

	return trail, nil
}

// GetTrail retrieves a CloudTrail trail by name. A missing record yields
// ErrTrailNotFound; every other failure (I/O, a corrupt record) propagates
// so callers can tell an absent trail from a broken read.
func (s *CloudTrailStore) GetTrail(trailName string) (*Trail, error) {
	var p pb.Trail
	if err := s.BaseStore.GetProto(trailName, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrTrailNotFound
		}
		return nil, err
	}
	return ProtoToTrail(&p), nil
}

// GetTrailByARN retrieves a CloudTrail trail by ARN.
func (s *CloudTrailStore) GetTrailByARN(trailARN string) (*Trail, error) {
	if s.arnIndexStore != nil && s.arnIndexStore.Exists(trailARN) {
		var trailName string
		if err := s.arnIndexStore.Get(trailARN, &trailName); err == nil {
			return s.GetTrail(trailName)
		}
	}

	trails, err := common.ListMatchingProto[*pb.Trail](s.BaseStore, "", func() *pb.Trail { return &pb.Trail{} }, func(t *pb.Trail) bool {
		return t.TrailArn == trailARN
	})
	if err != nil {
		return nil, err
	}
	if len(trails) > 0 {
		return ProtoToTrail(trails[0]), nil
	}
	return nil, ErrTrailNotFound
}

// ResolveTrail resolves a trail by name or ARN.
func (s *CloudTrailStore) ResolveTrail(nameOrARN string) (*Trail, error) {
	_, _, _, _, resource := svcarn.SplitARN(nameOrARN)
	if strings.HasPrefix(resource, "trail/") {
		return s.GetTrailByARN(nameOrARN)
	}
	return s.GetTrail(nameOrARN)
}

// MutateTrail loads a trail by name, applies the mutation under the store
// mutex, and persists the result — the load-apply-persist surface trail
// writers use (API updates, delivery bookkeeping) so concurrent writers
// cannot lose each other's changes. A missing record answers
// ErrTrailNotFound; any other load failure propagates so callers can tell
// an absent trail from a broken read. An error from apply aborts the
// mutation verbatim; ErrUnchanged leaves the record unwritten.
func (s *CloudTrailStore) MutateTrail(trailName string, apply func(*Trail) error) (*Trail, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var p pb.Trail
	if err := s.BaseStore.GetProto(trailName, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrTrailNotFound
		}
		return nil, err
	}
	trail := ProtoToTrail(&p)
	if err := apply(trail); err != nil {
		if errors.Is(err, ErrUnchanged) {
			return trail, nil
		}
		return nil, err
	}
	trail.LastUpdated = time.Now().UTC()
	if err := s.PutProto(trail.Name, TrailToProto(trail)); err != nil {
		return nil, err
	}
	return trail, nil
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

// PutEvent stores a CloudTrail event. Besides the event history entry, the
// record is copied into every event data store whose advanced event
// selectors match it, making each store an independent data boundary (Lake
// queries read only their own store's copies).
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

	// The EDS list is read before the transaction opens so the fan-out
	// writes join the event's own single transaction.
	edsList, err := s.listEventDataStoresRaw()
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
				if err := s.indexer.AddIndexInTxn(txn, event); err != nil {
					return err
				}
			}
			return s.ingestEventIntoEDSs(txn, edsList, event, key, eventData)
		})
	}

	if err := s.eventsStore.PutProto(key, EventToProto(event)); err != nil {
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
	return s.ingestEventIntoEDSs(nil, edsList, event, key, eventData)
}

// LookupEvents looks up CloudTrail events by query. Results are ordered
// most recent first on every path (the LookupEvents contract). Pagination
// runs through an opaque IndexCursor encoded into the returned nextToken —
// every path, including the filterless scan, issues this one token format —
// and a nextToken that does not decode yields ErrInvalidNextToken rather
// than silently restarting the walk. Callers should treat nextToken as
// opaque and pass it back unchanged in subsequent calls.
func (s *CloudTrailStore) LookupEvents(query EventQuery) ([]*Event, string, error) {
	if query.MaxResults <= 0 {
		query.MaxResults = DefaultLookupEventsResults
	}

	// Decode the incoming nextToken into an IndexCursor. A non-empty token
	// that does not decode was never issued by this store.
	cursor, cursorErr := decodeIndexCursor(query.NextToken)
	if cursorErr != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrInvalidNextToken, cursorErr)
	}

	var eventIDs []string
	var nextCursor IndexCursor
	var err error

	switch {
	case query.EventID != "":
		event, getErr := s.GetEventByID(query.EventID)
		if getErr != nil {
			if errors.Is(getErr, ErrEventNotFound) {
				// An EventId that matches no recorded event is a valid
				// empty result, not a failure.
				return nil, "", nil
			}
			return nil, "", getErr
		}
		return []*Event{event}, "", nil
	case len(query.EventNames) > 0 && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByEventName(query.EventNames, query.MaxResults, cursor)
	case query.Username != "" && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByUsername(query.Username, query.MaxResults, cursor)
	case query.EventSource != "" && s.indexer != nil:
		eventIDs, nextCursor, err = s.indexer.QueryByEventSource(query.EventSource, query.MaxResults, cursor)
	case (query.StartTime != nil || query.EndTime != nil) && s.indexer != nil:
		// An unbounded side is resolved against the recorded event span,
		// so a single-bound lookup walks every hour from its bound to the
		// newest (or from the oldest to its bound) recorded event.
		start, end, spanErr := s.resolveTimeBounds(query.StartTime, query.EndTime)
		if spanErr != nil {
			return nil, "", spanErr
		}
		if start == nil || end == nil {
			// No recorded events at all.
			return nil, "", nil
		}
		eventIDs, nextCursor, err = s.indexer.QueryByTime(start, end, query.MaxResults, cursor)
	default:
		return s.lookupEventsScan(query, cursor)
	}

	if err != nil {
		return nil, "", err
	}

	var events []*Event
	for _, id := range eventIDs {
		if int32(len(events)) >= query.MaxResults {
			break
		}
		event, getErr := s.GetEventByID(id)
		if getErr != nil {
			// An index entry whose event is gone (purged between the
			// index walk and the fetch) is skipped; a genuine read
			// failure surfaces.
			if errors.Is(getErr, ErrEventNotFound) {
				continue
			}
			return nil, "", getErr
		}
		if s.eventMatchesQuery(event, query) {
			events = append(events, event)
		}
	}

	return events, encodeIndexCursor(nextCursor), nil
}

// resolveTimeBounds completes a partially-bounded time query from the
// recorded event span: a missing start becomes the oldest recorded event's
// time and a missing end the newest's. Both bounds are returned nil when no
// events are recorded at all.
func (s *CloudTrailStore) resolveTimeBounds(startTime, endTime *time.Time) (*time.Time, *time.Time, error) {
	if startTime != nil && endTime != nil {
		return startTime, endTime, nil
	}
	oldest, newest, err := s.eventSpanBounds()
	if err != nil {
		return nil, nil, err
	}
	if startTime == nil {
		startTime = oldest
	}
	if endTime == nil {
		endTime = newest
	}
	return startTime, endTime, nil
}

// eventSpanBounds reads the oldest and newest recorded event times from the
// endpoint keys of the events bucket ("<unixnano>#<eventID>", whose
// lexicographic order is chronological). Nil, nil means the bucket is
// empty.
func (s *CloudTrailStore) eventSpanBounds() (*time.Time, *time.Time, error) {
	parse := func(key string) (*time.Time, error) {
		nanos, _, ok := strings.Cut(key, "#")
		if !ok {
			return nil, fmt.Errorf("cloudtrail store: malformed event key %q", key)
		}
		ts, err := strconv.ParseInt(nanos, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cloudtrail store: malformed event key %q: %w", key, err)
		}
		t := time.Unix(0, ts).UTC()
		return &t, nil
	}

	iter := s.eventsStore.Bucket().ScanPrefix(nil)
	var oldest *time.Time
	if iter.Next() {
		t, err := parse(string(iter.Key()))
		if err != nil {
			iter.Close()
			return nil, nil, err
		}
		oldest = t
	}
	if err := iter.Error(); err != nil {
		iter.Close()
		return nil, nil, err
	}
	iter.Close()
	if oldest == nil {
		return nil, nil, nil
	}

	revIter := s.eventsStore.Bucket().ScanPrefixReverse(nil, nil)
	var newest *time.Time
	if revIter.Next() {
		t, err := parse(string(revIter.Key()))
		if err != nil {
			revIter.Close()
			return nil, nil, err
		}
		newest = t
	}
	if err := revIter.Error(); err != nil {
		revIter.Close()
		return nil, nil, err
	}
	revIter.Close()

	return oldest, newest, nil
}

func (s *CloudTrailStore) eventIDIndexBucket() storage.Bucket {
	if s.storage != nil {
		if basic, ok := s.storage.(storage.BasicStorage); ok {
			return basic.Bucket(eventIDIndexBucketName(s.region))
		}
	}
	return nil
}

// GetEventByID retrieves a CloudTrail event by ID. A missing event yields
// ErrEventNotFound; every other failure (I/O, a corrupt record) propagates
// so callers can tell an absent event from a broken read.
func (s *CloudTrailStore) GetEventByID(eventID string) (*Event, error) {
	var fullKey string
	if bucket := s.eventIDIndexBucket(); bucket != nil {
		fullKeyBytes, err := bucket.Get([]byte(eventID))
		if err != nil {
			return nil, err
		}
		if fullKeyBytes == nil {
			return nil, ErrEventNotFound
		}
		fullKey = string(fullKeyBytes)
	} else if err := s.eventIDIndexStore.Get(eventID, &fullKey); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	var p pb.Event
	if err := s.eventsStore.GetProto(fullKey, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return ProtoToEvent(&p), nil
}

// lookupEventsScan serves the filterless default path: a newest-first walk
// over the whole events bucket ("<unixnano>#<eventID>" keys order
// chronologically), filtered by the query, paginated through the shared
// IndexCursor token — the cursor's Key is the last served storage key and
// iteration resumes strictly below it.
func (s *CloudTrailStore) lookupEventsScan(query EventQuery, cursor IndexCursor) ([]*Event, string, error) {
	before := []byte(cursor.Key)
	if len(before) == 0 {
		before = nil
	}
	iter := s.eventsStore.Bucket().ScanPrefixReverse(nil, before)
	defer iter.Close()

	var events []*Event
	var lastKey string
	for iter.Next() {
		if int32(len(events)) >= query.MaxResults {
			// Budget reached while another key remains below: issue a
			// continuation token so the walk can resume.
			return events, encodeIndexCursor(IndexCursor{Key: lastKey}), nil
		}
		key := string(iter.Key())
		var p pb.Event
		if err := proto.Unmarshal(iter.Value(), &p); err != nil {
			return nil, "", err
		}
		lastKey = key
		if protoMatchesQuery(&p, query) {
			events = append(events, ProtoToEvent(&p))
		}
	}
	if err := iter.Error(); err != nil {
		return nil, "", err
	}

	return events, "", nil
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
func (s *CloudTrailStore) RecordServiceEvent(eventName, eventSource string, userIdentity *UserIdentity, sourceIP, accessKeyID, userAgent string, readOnly bool, errorCode, errorMessage string, requestParams, responseElements map[string]interface{}, resources []Resource) error {
	event := NewEvent(eventName, eventSource, userIdentity, readOnly)
	event.AwsRegion = s.region
	event.RequestParameters = requestParams
	event.ResponseElements = responseElements
	event.SourceIPAddress = sourceIP
	event.AccessKeyId = accessKeyID
	event.UserAgent = userAgent
	event.ErrorCode = errorCode
	event.ErrorMessage = errorMessage
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
		MaxResults: DefaultLookupEventsResults,
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

// GenerateAndStorePublicKey creates a new RSA key pair and stores the public
// key alongside its private half: the trail's hourly digest chain signs with
// this key across restarts, so the private material must survive the process.
func (s *CloudTrailStore) GenerateAndStorePublicKey(trailName string) (*PublicKey, error) {
	pk, priv, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	pk.TrailName = trailName
	pk.PrivateKeyDER = x509.MarshalPKCS1PrivateKey(priv)
	if err := s.StorePublicKey(pk); err != nil {
		return nil, err
	}
	return pk, nil
}

// LoadTrailSigningKey returns the trail's newest validation key pair. A
// stored key without private material (issued before the digest chain
// persisted it) is replaced by a freshly generated pair, so the caller
// always receives a usable signer.
func (s *CloudTrailStore) LoadTrailSigningKey(trailName string) (*PublicKey, *rsa.PrivateKey, error) {
	keys, err := s.ListPublicKeys(nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var newest *PublicKey
	for _, pk := range keys {
		if pk.TrailName != trailName || len(pk.PrivateKeyDER) == 0 {
			continue
		}
		if newest == nil || pk.ValidityStartTime.After(newest.ValidityStartTime) {
			newest = pk
		}
	}
	if newest != nil {
		priv, err := x509.ParsePKCS1PrivateKey(newest.PrivateKeyDER)
		if err == nil {
			return newest, priv, nil
		}
	}
	pub, err := s.GenerateAndStorePublicKey(trailName)
	if err != nil {
		return nil, nil, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(pub.PrivateKeyDER)
	if err != nil {
		return nil, nil, err
	}
	return pub, priv, nil
}

// CreateAndStoreSigningKey generates a fresh RSA signing key for a query
// result sign file, persists the public half as a region key (no trail
// association, so trail deletion never collects it), and returns both
// halves: the private key signs the sign file now, the stored public key
// is what ListPublicKeys later serves to validators matching its
// fingerprint.
func (s *CloudTrailStore) CreateAndStoreSigningKey() (*PublicKey, *rsa.PrivateKey, error) {
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	if err := s.StorePublicKey(pub); err != nil {
		return nil, nil, err
	}
	return pub, priv, nil
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

	// The event data store quota counts every lifecycle stage — the
	// PENDING_DELETION records above included: "Event data stores — 10 ...
	// This includes event data stores in any lifecycle stage" (Quotas in
	// AWS CloudTrail). The list was loaded under this mutex, so the count
	// and the write are one step.
	if len(existing) >= MaxEventDataStoresPerRegion {
		return nil, ErrEventDataStoreQuotaExceeded
	}

	eds.CreatedTimestamp = time.Now().UTC()
	eds.UpdatedTimestamp = eds.CreatedTimestamp

	if err := s.eventDataStoreStore.Put(eds.EventDataStoreID, eds); err != nil {
		return nil, fmt.Errorf("failed to store event data store: %w", err)
	}

	return eds, nil
}

// GetEventDataStore retrieves an event data store by ID or ARN. A missing
// record yields ErrEventDataStoreNotFound; every other failure (I/O, a
// corrupt record) propagates so guards can tell an absent destination from
// a broken one.
func (s *CloudTrailStore) GetEventDataStore(idOrARN string) (*EventDataStore, error) {
	id := ExtractEventDataStoreID(idOrARN)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventDataStoreNotFound
		}
		return nil, err
	}
	return &eds, nil
}

// ListEventDataStores returns all event data stores.
func (s *CloudTrailStore) ListEventDataStores(opts common.ListOptions) (*common.ListResult[EventDataStore], error) {
	return common.List[EventDataStore](s.eventDataStoreStore, opts, nil)
}

// ListEventDataStoresAll drains every ListEventDataStores page. The
// retention sweep and the hard-delete sweeper must see every store, not the
// first page.
func (s *CloudTrailStore) ListEventDataStoresAll() ([]*EventDataStore, error) {
	var all []*EventDataStore
	marker := ""
	for {
		result, err := s.ListEventDataStores(common.ListOptions{
			MaxItems: MaxListEventDataStoresResults,
			Marker:   marker,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, result.Items...)
		if result.NextMarker == "" {
			return all, nil
		}
		marker = result.NextMarker
	}
}

// ErrUnchanged aborts a mutation without writing: an apply function returns
// it to leave the loaded record exactly as it was. The executor transitions
// use it to refuse overwriting a status another caller has already settled.
var ErrUnchanged = errors.New("record unchanged")

// MutateEventDataStore loads the event data store addressed by ID or ARN,
// applies apply under the store mutex, and persists the result. The whole
// read-modify-write is atomic: concurrent updates, ingestion toggles,
// federation changes, and lifecycle transitions cannot lose writes or
// interleave with each other. An error from apply aborts the mutation and
// is returned verbatim. The store mutex is held while apply runs, so apply
// must not call mutating store methods; read-only ones are safe.
func (s *CloudTrailStore) MutateEventDataStore(idOrARN string, apply func(*EventDataStore) error) (*EventDataStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := ExtractEventDataStoreID(idOrARN)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventDataStoreNotFound
		}
		return nil, err
	}
	if err := apply(&eds); err != nil {
		if errors.Is(err, ErrUnchanged) {
			return &eds, nil
		}
		return nil, err
	}
	eds.UpdatedTimestamp = time.Now().UTC()
	if err := s.eventDataStoreStore.Put(eds.EventDataStoreID, &eds); err != nil {
		return nil, err
	}
	return &eds, nil
}

// RestoreEventDataStore restores a PENDING_DELETION event data store to
// ENABLED status. The status check and the write run under the store mutex
// as one step.
func (s *CloudTrailStore) RestoreEventDataStore(id string) (*EventDataStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id = ExtractEventDataStoreID(id)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventDataStoreNotFound
		}
		return nil, err
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

// DeleteEventDataStoreIf deletes an event data store's record when guard
// raises no objection. The guard runs under the store mutex with the loaded
// record, so its verdict cannot race with a concurrent update or restore;
// ErrUnchanged from the guard refuses the delete without an error. The
// store's event bucket is NOT touched — the caller owns dropping it after
// the record is gone.
func (s *CloudTrailStore) DeleteEventDataStoreIf(idOrARN string, guard func(*EventDataStore) error) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := ExtractEventDataStoreID(idOrARN)
	var eds EventDataStore
	if err := s.eventDataStoreStore.Get(id, &eds); err != nil {
		if common.IsNotFound(err) {
			return false, ErrEventDataStoreNotFound
		}
		return false, err
	}
	if err := guard(&eds); err != nil {
		if errors.Is(err, ErrUnchanged) {
			return false, nil
		}
		return false, err
	}
	if err := s.eventDataStoreStore.Delete(id); err != nil {
		return false, err
	}
	return true, nil
}

// listEventDataStoresRaw returns all stored event data stores. A record
// that fails to unmarshal aborts the listing — corrupt state must surface,
// not silently shrink the result set.
func (s *CloudTrailStore) listEventDataStoresRaw() ([]*EventDataStore, error) {
	var result []*EventDataStore
	err := s.eventDataStoreStore.ForEach(func(_ string, value []byte) error {
		var eds EventDataStore
		if err := json.Unmarshal(value, &eds); err != nil {
			return err
		}
		result = append(result, &eds)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list event data stores: %w", err)
	}
	return result, nil
}

// ExtractEventDataStoreID extracts the UUID from an event data store ID or
// ARN. It is the single definition of EDS identity normalisation: every
// caller that must turn an ID-or-ARN wire value into the storage key routes
// through here.
func ExtractEventDataStoreID(idOrARN string) string {
	if idx := strings.LastIndex(idOrARN, "/"); idx >= 0 {
		return idOrARN[idx+1:]
	}
	return idOrARN
}

// --- Query operations ---

// SaveQuery persists a new query record. It is the creation write; every
// later status transition goes through MutateQuery so a terminal status
// written by CancelQuery can never be blindly overwritten.
func (s *CloudTrailStore) SaveQuery(qr *QueryRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queryStore.Put(qr.QueryID, qr)
}

// AdmitQuery persists a new query record unless the store already holds
// maxRunning non-terminal queries, in which case it refuses the admission
// with ErrMaxConcurrentQueries. The count and the write run under the
// store mutex as one step, so two simultaneous admissions cannot both
// slip past the bound.
func (s *CloudTrailStore) AdmitQuery(qr *QueryRecord, maxRunning int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	running := 0
	err := s.queryStore.ForEach(func(_ string, value []byte) error {
		var existing QueryRecord
		if err := json.Unmarshal(value, &existing); err != nil {
			return err
		}
		if !QueryTerminalStatus(existing.QueryStatus) {
			running++
		}
		return nil
	})
	if err != nil {
		return err
	}
	if running >= maxRunning {
		return ErrMaxConcurrentQueries
	}
	return s.queryStore.Put(qr.QueryID, qr)
}

// GetQuery retrieves a query record by ID.
func (s *CloudTrailStore) GetQuery(queryID string) (*QueryRecord, error) {
	var qr QueryRecord
	if err := s.queryStore.Get(queryID, &qr); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrQueryNotFound
		}
		return nil, err
	}
	return &qr, nil
}

// MutateQuery loads the query record, applies apply under the store mutex,
// and persists the result — the compare-and-set the executor and CancelQuery
// coordinate through: both check the persisted status inside apply, so the
// status one caller wrote is never silently overwritten by the other. An
// error from apply aborts the mutation and is returned verbatim.
func (s *CloudTrailStore) MutateQuery(queryID string, apply func(*QueryRecord) error) (*QueryRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var qr QueryRecord
	if err := s.queryStore.Get(queryID, &qr); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrQueryNotFound
		}
		return nil, err
	}
	if err := apply(&qr); err != nil {
		if errors.Is(err, ErrUnchanged) {
			return &qr, nil
		}
		return nil, err
	}
	if err := s.queryStore.Put(qr.QueryID, &qr); err != nil {
		return nil, err
	}
	return &qr, nil
}

// ListQueriesByEDS lists queries for an event data store. A record that
// fails to unmarshal aborts the listing — corrupt state must surface, not
// silently shrink the result set.
func (s *CloudTrailStore) ListQueriesByEDS(edsID string) ([]*QueryRecord, error) {
	var result []*QueryRecord
	err := s.queryStore.ForEach(func(_ string, value []byte) error {
		var qr QueryRecord
		if err := json.Unmarshal(value, &qr); err != nil {
			return err
		}
		if ExtractEventDataStoreID(qr.EventDataStore) == edsID {
			result = append(result, &qr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PurgeQueriesBefore deletes every query record whose StartTime precedes
// the cutoff, collecting the doomed keys first and deleting after the
// walk — the store's purge pattern. "Returns a list of queries and query
// statuses for the past seven days" (ListQueries) bounds the records'
// lifetime as well as the listing — the retention sweep keeps the bucket
// from growing without bound.
func (s *CloudTrailStore) PurgeQueriesBefore(cutoff time.Time) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	purged := 0
	var doomed []string
	err := s.queryStore.ForEach(func(key string, value []byte) error {
		var qr QueryRecord
		if err := json.Unmarshal(value, &qr); err != nil {
			return err
		}
		if qr.StartTime.Before(cutoff) {
			doomed = append(doomed, key)
		}
		return nil
	})
	if err != nil {
		return purged, err
	}
	for _, key := range doomed {
		if err := s.queryStore.Delete(key); err != nil {
			return purged, err
		}
		purged++
	}
	return purged, nil
}

// --- Channel operations ---

// CreateChannel persists a new channel. Channel names are unique within
// the account and region: a name already carried by another channel is
// rejected with ErrChannelAlreadyExists. Sources are unique too — "A
// maximum of one channel is allowed per source" (CreateChannel Source) —
// so a source already carried by another channel is rejected with
// ErrChannelSourceInUse.
func (s *CloudTrailStore) CreateChannel(ch *Channel) (*Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taken, err := s.channelTaken("", func(c *Channel) bool { return c.Name == ch.Name })
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrChannelAlreadyExists
	}
	sourceTaken, err := s.channelTaken("", func(c *Channel) bool { return c.Source == ch.Source })
	if err != nil {
		return nil, err
	}
	if sourceTaken {
		return nil, ErrChannelSourceInUse
	}

	// The channel quota holds at admission: "Channels — 25 ... This quota
	// cannot be increased" (Quotas in AWS CloudTrail); the count runs under
	// the creation mutex so the check and the write are one step.
	channelCount := 0
	if err := s.channelStore.ForEach(func(_ string, _ []byte) error {
		channelCount++
		return nil
	}); err != nil {
		return nil, err
	}
	if channelCount >= MaxChannelsPerRegion {
		return nil, ErrChannelQuotaExceeded
	}

	ch.CreatedAt = time.Now().UTC()
	ch.UpdatedAt = ch.CreatedAt
	return ch, s.channelStore.Put(ch.ChannelARN, ch)
}

// channelTaken reports whether a channel other than excludeARN matches
// match — the one scan the name- and source-uniqueness guards share,
// running under the caller's store mutex.
func (s *CloudTrailStore) channelTaken(excludeARN string, match func(*Channel) bool) (bool, error) {
	taken := false
	err := s.channelStore.ForEach(func(_ string, value []byte) error {
		var ch Channel
		if err := json.Unmarshal(value, &ch); err != nil {
			return nil
		}
		if ch.ChannelARN != excludeARN && match(&ch) {
			taken = true
		}
		return nil
	})
	return taken, err
}

// GetChannel retrieves a channel by ARN.
func (s *CloudTrailStore) GetChannel(arn string) (*Channel, error) {
	var ch Channel
	if err := s.channelStore.Get(arn, &ch); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}
	return &ch, nil
}

// MutateChannel loads the channel, applies apply under the store mutex,
// enforces channel-name uniqueness (the applied record may carry a new
// name), and persists the result. The whole read-modify-write is atomic.
// An error from apply aborts the mutation and is returned verbatim.
func (s *CloudTrailStore) MutateChannel(arn string, apply func(*Channel) error) (*Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ch Channel
	if err := s.channelStore.Get(arn, &ch); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}
	if err := apply(&ch); err != nil {
		if errors.Is(err, ErrUnchanged) {
			return &ch, nil
		}
		return nil, err
	}
	taken, err := s.channelTaken(ch.ChannelARN, func(c *Channel) bool { return c.Name == ch.Name })
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrChannelAlreadyExists
	}
	ch.UpdatedAt = time.Now().UTC()
	if err := s.channelStore.Put(ch.ChannelARN, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// DeleteChannelIf deletes a channel when guard raises no objection. The
// guard runs under the store mutex with the loaded record, so its verdict
// cannot race with a concurrent channel update; an error from guard aborts
// the delete and is returned verbatim.
func (s *CloudTrailStore) DeleteChannelIf(arn string, guard func(*Channel) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ch Channel
	if err := s.channelStore.Get(arn, &ch); err != nil {
		if common.IsNotFound(err) {
			return ErrChannelNotFound
		}
		return err
	}
	if err := guard(&ch); err != nil {
		return err
	}
	return s.channelStore.Delete(arn)
}

// ListChannels lists channels with pagination support.
func (s *CloudTrailStore) ListChannels(opts common.ListOptions) (*common.ListResult[Channel], error) {
	return common.List[Channel](s.channelStore, opts, nil)
}

// --- Event Configuration ---

// GetEventConfiguration retrieves event configuration for a trail or EDS.
// A missing record yields ErrEventConfigurationNotFound; every other
// failure propagates (the Lake families' not-found/I-O distinction).
func (s *CloudTrailStore) GetEventConfiguration(trailName, edsID string) (map[string]interface{}, error) {
	key := eventConfigKey(trailName, edsID)
	var config map[string]interface{}
	if err := s.eventConfigStore.Get(key, &config); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventConfigurationNotFound
		}
		return nil, err
	}
	return config, nil
}

// PutEventConfiguration stores event configuration.
func (s *CloudTrailStore) PutEventConfiguration(trailName, edsID string, config map[string]interface{}) error {
	key := eventConfigKey(trailName, edsID)
	return s.eventConfigStore.Put(key, config)
}

// DeleteEventConfiguration removes the event configuration for a trail or
// event data store — the trail-deletion cleanup path, so a deleted trail
// leaves no orphan configuration record behind.
func (s *CloudTrailStore) DeleteEventConfiguration(trailName, edsID string) error {
	key := eventConfigKey(trailName, edsID)
	return s.eventConfigStore.Delete(key)
}

func eventConfigKey(trailName, edsID string) string {
	if trailName != "" {
		return "trail:" + trailName
	}
	return "eds:" + edsID
}

// --- Import operations ---

// CreateImport persists a new import record.
func (s *CloudTrailStore) CreateImport(imp *Import) (*Import, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return imp, s.importStore.Put(imp.ImportID, imp)
}

// ErrImportOngoing is returned by CreateImportIfNoOngoing when any import
// record is INITIALIZING or IN_PROGRESS.
var ErrImportOngoing = errors.New("an import is already in progress")

// Resource-quota sentinels: the creation methods enforce the fetched
// "Quotas in AWS CloudTrail" counts under the store mutex; the service
// layer maps each to the model's declared error shape.
var (
	ErrTrailQuotaExceeded          = errors.New("the trail quota for the region is exceeded")
	ErrEventDataStoreQuotaExceeded = errors.New("the event data store quota for the region is exceeded")
	ErrChannelQuotaExceeded        = errors.New("the channel quota for the region is exceeded")
)

// CreateImportIfNoOngoing admits a new import only when no other import is
// INITIALIZING or IN_PROGRESS: the ongoing scan and the creation write run
// under the store mutex as one step, so two concurrent StartImport calls
// cannot both pass the one-ongoing-import contract (AdmitQuery is the
// precedent for count-and-write admission).
func (s *CloudTrailStore) CreateImportIfNoOngoing(imp *Import) (*Import, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.importStore.ForEach(func(_ string, value []byte) error {
		var existing Import
		if err := json.Unmarshal(value, &existing); err != nil {
			return err
		}
		if existing.ImportStatus == "INITIALIZING" || existing.ImportStatus == "IN_PROGRESS" {
			return ErrImportOngoing
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return imp, s.importStore.Put(imp.ImportID, imp)
}

// GetImport retrieves an import by ID.
func (s *CloudTrailStore) GetImport(importID string) (*Import, error) {
	var imp Import
	if err := s.importStore.Get(importID, &imp); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrImportNotFound
		}
		return nil, err
	}
	return &imp, nil
}

// MutateImport loads the import record, applies apply under the store mutex,
// and persists the result. The status transitions of the import executor and
// StopImport coordinate through it: both check the persisted status inside
// apply, so a STOPPED verdict written mid-execution is never overwritten.
// An error from apply aborts the mutation and is returned verbatim.
func (s *CloudTrailStore) MutateImport(importID string, apply func(*Import) error) (*Import, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var imp Import
	if err := s.importStore.Get(importID, &imp); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrImportNotFound
		}
		return nil, err
	}
	if err := apply(&imp); err != nil {
		if errors.Is(err, ErrUnchanged) {
			return &imp, nil
		}
		return nil, err
	}
	imp.UpdatedTimestamp = time.Now().UTC()
	if err := s.importStore.Put(imp.ImportID, &imp); err != nil {
		return nil, err
	}
	return &imp, nil
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

	maxItems := opts.MaxItems
	if maxItems <= 0 {
		maxItems = DefaultListImportsResults
	}

	paged := pagination.PaginateSliceByPosition(imp.Failures, opts.Marker, maxItems)
	items := make([]*ImportFailure, 0, len(paged.Items))
	for i := range paged.Items {
		items = append(items, &paged.Items[i])
	}

	return &common.ListResult[ImportFailure]{
		Items:       items,
		NextMarker:  paged.NextMarker,
		IsTruncated: paged.IsTruncated,
	}, nil
}
