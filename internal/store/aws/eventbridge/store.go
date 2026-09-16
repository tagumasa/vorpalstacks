// Package eventbridge provides EventBridge storage functionality for vorpalstacks.
package eventbridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// EventsStore provides EventBridge storage functionality.
type EventsStore struct {
	*common.BaseStore
	rulesStore           *common.BaseStore
	targetsStore         *common.BaseStore
	archivesStore        *common.BaseStore
	archivedEventsStore  *common.BaseStore
	replaysStore         *common.BaseStore
	connectionsStore     *common.BaseStore
	apiDestinationsStore *common.BaseStore
	// archiveBusIndexStore holds the per-bus archive index: one row per
	// (source bus, archive) pair, key "<bus>:<archive>", value the archive
	// name. The delivery path resolves a bus's archives per delivered event,
	// so the index answers with a prefix scan of the bus's own rows instead
	// of walking and decoding every archive record in the region per event.
	archiveBusIndexStore *common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	createMu   sync.Mutex
}

// The recordWriteMu family serialises record read-modify-write and delete
// cycles per family across EventsStore instances: the scheduler, replay
// and delivery workers and the API handlers operate separate store
// instances over the same Pebble keyspace, so instance-scoped locking
// cannot protect a record cycle. Every Mutate* callback and every Delete*
// runs under its family's mutex, which also closes the delete-vs-mutate
// resurrection window (a cycle that read before the delete can no longer
// write the record back after it).
var (
	eventBusRecordWriteMu       sync.Mutex
	targetRecordWriteMu         sync.Mutex
	archiveRecordWriteMu        sync.Mutex
	connectionRecordWriteMu     sync.Mutex
	apiDestinationRecordWriteMu sync.Mutex
	replayRecordWriteMu         sync.Mutex
)

// ruleRecordWriteMu serialises rule-record read-modify-write cycles
// across EventsStore instances: the scheduler worker and the API
// handlers operate separate store instances over the same Pebble
// keyspace, so TouchRuleLastFired's read-modify-write could otherwise
// lose a concurrent UpdateRule write (or vice versa).
var ruleRecordWriteMu sync.Mutex

// Region returns the region this store instance serves. Cross-service
// calls keyed by region (e.g. the connection-credential secret in the
// region's Secrets Manager) use it to address the same region as the
// records they belong to.
func (s *EventsStore) Region() string { return s.region }

// NewEventsStore creates a new EventBridge events store.
func NewEventsStore(store storage.BasicStorage, accountID, region string) *EventsStore {
	return &EventsStore{
		BaseStore:            common.NewBaseStore(store.Bucket("events-eventbuses-"+region), "events-eventbuses"),
		rulesStore:           common.NewBaseStore(store.Bucket("events-rules-"+region), "events-rules"),
		targetsStore:         common.NewBaseStore(store.Bucket("events-targets-"+region), "events-targets"),
		archivesStore:        common.NewBaseStore(store.Bucket("events-archives-"+region), "events-archives"),
		archivedEventsStore:  common.NewBaseStore(store.Bucket("events-archived-events-"+region), "events-archived-events"),
		archiveBusIndexStore: common.NewBaseStore(store.Bucket("events-archive-bus-index-"+region), "events-archive-bus-index"),
		replaysStore:         common.NewBaseStore(store.Bucket("events-replays-"+region), "events-replays"),
		connectionsStore:     common.NewBaseStore(store.Bucket("events-connections-"+region), "events-connections"),
		apiDestinationsStore: common.NewBaseStore(store.Bucket("events-apidestinations-"+region), "events-apidestinations"),
		TagStore:             common.NewTagStoreWithRegion(store, "events", region),
		arnBuilder:           svcarn.NewARNBuilder(accountID, region),
		accountID:            accountID,
		region:               region,
	}
}

// GetAccountID returns the AWS account ID associated with this store.
func (s *EventsStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the AWS region associated with this store.
func (s *EventsStore) GetRegion() string {
	return s.region
}

func (s *EventsStore) buildEventBusARN(name string) string {
	return s.arnBuilder.Events().EventBus(name)
}

func (s *EventsStore) buildRuleARN(eventBusName, ruleName string) string {
	return s.arnBuilder.Events().RuleOnBus(eventBusName, ruleName)
}

func (s *EventsStore) buildArchiveARN(name string) string {
	return s.arnBuilder.Events().Archive(name)
}

func (s *EventsStore) buildConnectionARN(name string) string {
	// The ARN carries a unique id after the name per the Smithy ConnectionArn
	// pattern (connection/<name>/<id>).
	return s.arnBuilder.Events().Connection(name, uuid.NewString())
}

func (s *EventsStore) buildApiDestinationARN(name string) string {
	// The ARN carries a unique id after the name per the Smithy
	// ApiDestinationArn pattern (api-destination/<name>/<id>).
	return s.arnBuilder.Events().ApiDestination(name, uuid.NewString())
}

// EventBus operations

// CreateEventBus creates a new event bus.
//
// Parameters:
//   - ctx: The context
//   - eventBus: The event bus to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateEventBus(ctx context.Context, eventBus *EventBus) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if eventBus.Name == "" {
		return ErrEmptyResourceName
	}

	arn := s.buildEventBusARN(eventBus.Name)
	if s.Exists(arn) {
		return ErrEventBusAlreadyExists
	}

	now := time.Now().UTC()
	eventBus.ARN = arn
	eventBus.Region = s.region
	eventBus.AccountID = s.accountID
	eventBus.CreatedAt = now
	eventBus.LastModifiedAt = now

	return s.Put(arn, eventBus)
}

// GetEventBus retrieves an event bus by name. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself so the
// Cores can tell a permanent absence from a transient failure.
//
// Parameters:
//   - ctx: The context
//   - name: The event bus name
//
// Returns:
//   - *EventBus: The event bus if found
//   - error: ErrEventBusNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetEventBus(ctx context.Context, name string) (*EventBus, error) {
	arn := s.buildEventBusARN(name)
	var eventBus EventBus
	if err := s.BaseStore.Get(arn, &eventBus); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrEventBusNotFound
		}
		return nil, err
	}
	return &eventBus, nil
}

// MutateEventBus applies fn to the event bus record inside the write mutex
// so the whole read-modify-write cycle is atomic with respect to every
// other event bus record writer (description updates, permission-policy
// merges and deletes run on separate store instances over the same Pebble
// keyspace, hence the package-scope lock). fn mutates the record in place;
// fn returning an error aborts the write. Callers stamp LastModifiedAt
// themselves.
//
// Parameters:
//   - ctx: The context
//   - name: The event bus name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrEventBusNotFound if the event bus does not exist, the
//     storage fault otherwise
func (s *EventsStore) MutateEventBus(ctx context.Context, name string, fn func(*EventBus) error) error {
	arn := s.buildEventBusARN(name)
	eventBusRecordWriteMu.Lock()
	defer eventBusRecordWriteMu.Unlock()

	var eventBus EventBus
	if err := s.BaseStore.Get(arn, &eventBus); err != nil {
		if common.IsNotFound(err) {
			return ErrEventBusNotFound
		}
		return err
	}
	if err := fn(&eventBus); err != nil {
		return err
	}
	return s.Put(arn, &eventBus)
}

// DeleteEventBus deletes an event bus by name.
//
// Parameters:
//   - ctx: The context
//   - name: The event bus name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteEventBus(ctx context.Context, name string) error {
	arn := s.buildEventBusARN(name)
	// The record lock closes the resurrection window: without it a
	// MutateEventBus cycle that read before the delete could write the
	// record back after it.
	eventBusRecordWriteMu.Lock()
	defer eventBusRecordWriteMu.Unlock()
	if !s.Exists(arn) {
		return ErrEventBusNotFound
	}
	_ = s.TagStore.Delete(arn)
	return s.BaseStore.Delete(arn)
}

// ListEventBuses lists event buses with optional filtering and pagination.
//
// Parameters:
//   - ctx: The context
//   - namePrefix: Optional name prefix filter
//   - limit: Maximum number of results
//   - nextToken: Pagination token
//
// Returns:
//   - *EventBusListResult: The list result with event buses and next token
//   - error: An error if listing fails
//
// eventBusKeyPrefix is the key head every event bus record carries: bus
// records are stored under their full ARN, which this prefix matches. The
// events-eventbuses bucket also holds residue rows from the retired
// raw-write cross-service ingress plane (keyed "events:<bus>:<id>"); the
// prefix keeps that residue invisible to listings without deleting data.
const eventBusKeyPrefix = "arn:aws:events:"

func (s *EventsStore) ListEventBuses(ctx context.Context, namePrefix string, limit int32, nextToken string) (*EventBusListResult, error) {
	opts := common.ListOptions{
		Prefix:   eventBusKeyPrefix,
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[EventBus](s.BaseStore, opts, func(e *EventBus) bool {
		if namePrefix == "" {
			return true
		}
		return strings.HasPrefix(e.Name, namePrefix)
	})
	if err != nil {
		return nil, err
	}

	return &EventBusListResult{
		EventBuses: result.Items,
		NextToken:  result.NextMarker,
	}, nil
}

// Rule operations

func (s *EventsStore) buildRuleKey(eventBusName, ruleName string) string {
	return fmt.Sprintf("%s:%s", eventBusName, ruleName)
}

// CreateRule creates a new rule on an event bus without a count cap —
// the test and internal seeding path.
//
// Parameters:
//   - ctx: The context
//   - rule: The rule to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateRule(ctx context.Context, rule *Rule) error {
	return s.CreateRuleCapped(ctx, rule, 0)
}

// CreateRuleCapped creates a new rule on an event bus under the per-bus
// rule-count quota: it mirrors CreateRule's contract and adds
// ErrRuleCapReached when creating a NEW rule would push the bus past
// maxPerBus records. The existence check, the count and the write run
// inside the one create lock, so concurrent creates cannot both pass
// the gate; an existing rule name never hits the cap — the quota bounds
// the creation of new rules, and updating a rule at the quota stays
// allowed. maxPerBus <= 0 disables the cap. The count lists at most
// maxPerBus+1 records: the gate needs the threshold comparison alone,
// not the exact total.
//
// Parameters:
//   - ctx: The context
//   - rule: The rule to create
//   - maxPerBus: The per-bus rule-count ceiling
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateRuleCapped(ctx context.Context, rule *Rule, maxPerBus int) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if rule.Name == "" {
		return ErrEmptyResourceName
	}

	key := s.buildRuleKey(rule.EventBusName, rule.Name)
	if s.rulesStore.Exists(key) {
		return ErrRuleAlreadyExists
	}
	if maxPerBus > 0 {
		counted, err := common.List[Rule](s.rulesStore, common.ListOptions{
			Prefix:   rule.EventBusName + ":",
			MaxItems: maxPerBus + 1,
		}, func(*Rule) bool { return true })
		if err != nil {
			return err
		}
		if len(counted.Items) >= maxPerBus {
			return ErrRuleCapReached
		}
	}

	now := time.Now().UTC()
	rule.ARN = s.buildRuleARN(rule.EventBusName, rule.Name)
	rule.Region = s.region
	rule.AccountID = s.accountID
	rule.CreatedAt = now
	rule.LastModifiedAt = now
	if rule.State == "" {
		rule.State = RuleStateEnabled
	}

	return s.rulesStore.Put(key, rule)
}

// GetRule retrieves a rule by event bus name and rule name. A missing
// record is the not-found sentinel; any other storage fault surfaces as
// itself so the Cores can tell a permanent absence from a transient
// failure.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - name: The rule name
//
// Returns:
//   - *Rule: The rule if found
//   - error: ErrRuleNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetRule(ctx context.Context, eventBusName, name string) (*Rule, error) {
	key := s.buildRuleKey(eventBusName, name)
	var rule Rule
	if err := s.rulesStore.Get(key, &rule); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}
	return &rule, nil
}

// MutateRule applies fn to the rule record inside the write mutex so the
// whole read-modify-write cycle is atomic with respect to every other rule
// record writer (the delivery worker and the service hold separate store
// instances over the same Pebble keyspace, hence the package-scope lock).
// fn mutates the record in place; callers decide whether to stamp
// LastModifiedAt — user-initiated updates do, internal markers must not.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - name: The rule name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrRuleNotFound if the rule does not exist, the storage
//     fault otherwise
func (s *EventsStore) MutateRule(ctx context.Context, eventBusName, name string, fn func(*Rule) error) error {
	key := s.buildRuleKey(eventBusName, name)
	ruleRecordWriteMu.Lock()
	defer ruleRecordWriteMu.Unlock()

	var rule Rule
	if err := s.rulesStore.Get(key, &rule); err != nil {
		if common.IsNotFound(err) {
			return ErrRuleNotFound
		}
		return err
	}
	if err := fn(&rule); err != nil {
		return err
	}
	return s.rulesStore.Put(key, &rule)
}

// TouchRuleLastFired records the schedule boundary a rule just fired
// under. The write carries the record as read and never stamps
// LastModifiedAt: firing is not a rule modification. The boundary only
// advances, so a late in-flight fire cannot regress a newer one.
func (s *EventsStore) TouchRuleLastFired(ctx context.Context, eventBusName, name string, firedAt time.Time) error {
	return s.MutateRule(ctx, eventBusName, name, func(rule *Rule) error {
		if !firedAt.After(rule.LastFiredAt) {
			return nil
		}
		rule.LastFiredAt = firedAt
		return nil
	})
}

// DeleteRule deletes a rule by event bus name and rule name.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - name: The rule name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteRule(ctx context.Context, eventBusName, name string) error {
	key := s.buildRuleKey(eventBusName, name)
	// The record lock closes the resurrection window: without it a
	// TouchRuleLastFired or MutateRule cycle that read before the delete
	// could write the record back after it.
	ruleRecordWriteMu.Lock()
	defer ruleRecordWriteMu.Unlock()
	if !s.rulesStore.Exists(key) {
		return ErrRuleNotFound
	}
	// Tags follow the record under the same lock, like every other
	// resource family's delete: a caller-side clean would leak tag rows to
	// any future direct store caller.
	_ = s.TagStore.Delete(s.buildRuleARN(eventBusName, name))
	return s.rulesStore.Delete(key)
}

// ListRules lists rules for an event bus with optional filtering.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - namePrefix: Optional name prefix filter
//   - limit: Maximum number of results
//   - nextToken: Pagination token
//
// Returns:
//   - *RuleListResult: The list result with rules and next token
//   - error: An error if listing fails
//
// ListRules lists rules, optionally scoped to one event bus. An empty
// eventBusName sweeps every bus — the scheduler's tick path, which must
// reach scheduled rules on all of them; the API plane always resolves a
// concrete bus name (defaulting to "default") before calling.
func (s *EventsStore) ListRules(ctx context.Context, eventBusName, namePrefix string, limit int32, nextToken string) (*RuleListResult, error) {
	prefix := eventBusName + ":"
	if eventBusName == "" {
		prefix = ""
	}
	opts := common.ListOptions{
		Prefix:   prefix,
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Rule](s.rulesStore, opts, func(r *Rule) bool {
		if namePrefix == "" {
			return true
		}
		return strings.HasPrefix(r.Name, namePrefix)
	})
	if err != nil {
		return nil, err
	}

	return &RuleListResult{
		Rules:     result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// Target operations

func (s *EventsStore) buildTargetKey(eventBusName, ruleName, targetID string) string {
	return fmt.Sprintf("%s:%s:%s", eventBusName, ruleName, targetID)
}

// PutTarget adds or updates a target for a rule. The CreatedAt
// read-preserve-write runs inside the record lock so a concurrent put of
// the same target ID can never observe or restore a zero creation time
// (putters race across separate store instances over the same Pebble
// keyspace, hence the package-scope lock).
//
// Parameters:
//   - ctx: The context
//   - target: The target to put
//
// Returns:
//   - error: An error if the operation fails
func (s *EventsStore) PutTarget(ctx context.Context, target *Target) error {
	key := s.buildTargetKey(target.EventBusName, target.RuleName, target.ID)
	targetRecordWriteMu.Lock()
	defer targetRecordWriteMu.Unlock()

	var existing Target
	if err := s.targetsStore.Get(key, &existing); err == nil && !existing.CreatedAt.IsZero() {
		target.CreatedAt = existing.CreatedAt
	} else {
		target.CreatedAt = time.Now().UTC()
	}
	return s.targetsStore.Put(key, target)
}

// GetTarget retrieves a target by event bus name, rule name, and target
// ID. A missing record is the not-found sentinel; any other storage fault
// surfaces as itself so the Cores can tell a permanent absence from a
// transient failure.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - ruleName: The rule name
//   - targetID: The target ID
//
// Returns:
//   - *Target: The target if found
//   - error: ErrTargetNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetTarget(ctx context.Context, eventBusName, ruleName, targetID string) (*Target, error) {
	key := s.buildTargetKey(eventBusName, ruleName, targetID)
	var target Target
	if err := s.targetsStore.Get(key, &target); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return &target, nil
}

// DeleteTarget deletes a target by event bus name, rule name, and target ID.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - ruleName: The rule name
//   - targetID: The target ID to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteTarget(ctx context.Context, eventBusName, ruleName, targetID string) error {
	key := s.buildTargetKey(eventBusName, ruleName, targetID)
	// The record lock closes the resurrection window against a concurrent
	// PutTarget of the same key.
	targetRecordWriteMu.Lock()
	defer targetRecordWriteMu.Unlock()
	if !s.targetsStore.Exists(key) {
		return ErrTargetNotFound
	}
	return s.targetsStore.Delete(key)
}

// ListTargetsByRule lists targets for a rule.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - ruleName: The rule name
//   - limit: Maximum number of results
//   - nextToken: Pagination token
//
// Returns:
//   - *TargetListResult: The list result with targets and next token
//   - error: An error if listing fails
func (s *EventsStore) ListTargetsByRule(ctx context.Context, eventBusName, ruleName string, limit int32, nextToken string) (*TargetListResult, error) {
	prefix := fmt.Sprintf("%s:%s:", eventBusName, ruleName)
	opts := common.ListOptions{
		Prefix:   prefix,
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Target](s.targetsStore, opts, nil)
	if err != nil {
		return nil, err
	}

	return &TargetListResult{
		Targets:   result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// Archive operations

// CreateArchive creates a new event archive.
//
// Parameters:
//   - ctx: The context
//   - archive: The archive to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateArchive(ctx context.Context, archive *Archive) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if archive.Name == "" {
		return ErrEmptyResourceName
	}

	if s.archivesStore.Exists(archive.Name) {
		return ErrArchiveAlreadyExists
	}

	now := time.Now().UTC()
	archive.ARN = s.buildArchiveARN(archive.Name)
	archive.Region = s.region
	archive.AccountID = s.accountID
	archive.CreatedAt = now
	if archive.State == "" {
		archive.State = ArchiveStateEnabled
	}

	// The per-bus index row lands with the record so the delivery path can
	// resolve the bus's archives from the moment the archive exists.
	if err := s.archiveBusIndexStore.Put(archive.EventBusName+":"+archive.Name, archive.Name); err != nil {
		return err
	}
	return s.archivesStore.Put(archive.Name, archive)
}

// GetArchive retrieves an archive by name. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself so the
// Cores can tell a permanent absence from a transient failure.
//
// Parameters:
//   - ctx: The context
//   - name: The archive name
//
// Returns:
//   - *Archive: The archive if found
//   - error: ErrArchiveNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetArchive(ctx context.Context, name string) (*Archive, error) {
	var archive Archive
	if err := s.archivesStore.Get(name, &archive); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrArchiveNotFound
		}
		return nil, err
	}
	return &archive, nil
}

// DeleteArchive deletes an archive by name: its stored event rows, the
// record, its tags and its bus-index row.
//
// Parameters:
//   - ctx: The context
//   - name: The archive name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteArchive(ctx context.Context, name string) error {
	// The record lock closes the resurrection window: without it a
	// MutateArchive cycle that read before the delete could write the
	// record back after it. The event-row sweep runs inside the same
	// locked section, so an ingress write (StoreArchiveEvent holds the same
	// lock and re-checks existence) cannot slip between the sweep and the
	// record delete — no orphaned event row survives the archive to
	// surface in a same-named archive created later.
	archiveRecordWriteMu.Lock()
	defer archiveRecordWriteMu.Unlock()
	if !s.archivesStore.Exists(name) {
		return ErrArchiveNotFound
	}
	// The index row's bus comes from the record itself, read under the
	// same lock as the delete so a concurrent create cannot interleave.
	current, err := s.GetArchive(ctx, name)
	if err != nil {
		return err
	}
	// Events go before the record: a fault between the two leaves a
	// retryable record whose events are partially gone — never orphaned
	// events under a deleted record (the retention worker no longer
	// iterates a deleted archive).
	if err := s.deleteArchiveEventsByName(name); err != nil {
		return err
	}
	_ = s.TagStore.Delete(s.buildArchiveARN(name))
	if err := s.archivesStore.Delete(name); err != nil {
		return err
	}
	// The index row goes after the record delete: a fault between the two
	// leaves a stale row the reader skips (its Get misses), never a live
	// record hidden from delivery.
	return s.archiveBusIndexStore.Delete(current.EventBusName + ":" + name)
}

// deleteArchiveEventsByName removes every stored event row of the archive.
// The caller holds the archive family lock (DeleteArchive) or accepts the
// unlocked sweep (DeleteExpiredArchiveEvents filters by time instead).
func (s *EventsStore) deleteArchiveEventsByName(name string) error {
	prefix := name + ":"
	return common.ForEachAll[ArchivedEvent](s.archivedEventsStore, prefix, nil, func(e *ArchivedEvent) error {
		key := prefix + fmt.Sprintf("%d:", e.Timestamp.UnixNano()) + e.ID
		return s.archivedEventsStore.Delete(key)
	})
}

// MutateArchive applies fn to the archive record inside the write mutex so
// the whole read-modify-write cycle is atomic with respect to every other
// archive record writer. Counter increments from the event delivery path
// and user-initiated configuration merges run on separate store instances
// over the same Pebble keyspace, hence the package-scope lock: without it a
// merge that read the record before a concurrent increment silently
// regresses EventCount/SizeBytes. fn mutates the record in place; fn
// returning an error aborts the write.
//
// Parameters:
//   - ctx: The context
//   - name: The archive name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrArchiveNotFound if the archive does not exist, the storage
//     fault otherwise
func (s *EventsStore) MutateArchive(ctx context.Context, name string, fn func(*Archive) error) error {
	archiveRecordWriteMu.Lock()
	defer archiveRecordWriteMu.Unlock()

	var archive Archive
	if err := s.archivesStore.Get(name, &archive); err != nil {
		if common.IsNotFound(err) {
			return ErrArchiveNotFound
		}
		return err
	}
	if err := fn(&archive); err != nil {
		return err
	}
	return s.archivesStore.Put(name, &archive)
}

// IncrementArchiveCounters increments EventCount and SizeBytes for an
// archive through the atomic record mutation so the increment can never be
// lost to (or regress) a concurrent configuration update.
func (s *EventsStore) IncrementArchiveCounters(ctx context.Context, archiveName string, eventSize int64) error {
	return s.MutateArchive(ctx, archiveName, func(archive *Archive) error {
		archive.EventCount++
		archive.SizeBytes += eventSize
		return nil
	})
}

// ArchiveListResult represents the result of listing archives.
type ArchiveListResult struct {
	Archives  []*Archive
	NextToken string
}

// ListArchivesForEventBus lists archives for a specific event bus.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//
// Returns:
//   - []*Archive: The list of archives
//   - error: An error if listing fails
func (s *EventsStore) ListArchivesForEventBus(ctx context.Context, eventBusName string) ([]*Archive, error) {
	// The per-bus index answers with a prefix scan of the bus's own rows;
	// bus names cannot contain ":" (the EventBusName Smithy pattern excludes
	// it), so "<bus>:" cannot bleed into another bus's rows. A stale index
	// row — whose archive was deleted between the index write and this read
	// — is skipped; a real storage fault surfaces as itself.
	var archives []*Archive
	err := s.archiveBusIndexStore.ScanPrefix(eventBusName+":", func(key string, value []byte) error {
		var name string
		if err := json.Unmarshal(value, &name); err != nil {
			return err
		}
		archive, err := s.GetArchive(ctx, name)
		if err != nil {
			if errors.Is(err, ErrArchiveNotFound) {
				return nil
			}
			return err
		}
		archives = append(archives, archive)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return archives, nil
}

// Connection operations

// CreateConnection creates a new EventBridge connection.
//
// Parameters:
//   - ctx: The context
//   - connection: The connection to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateConnection(ctx context.Context, connection *Connection) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if connection.Name == "" {
		return ErrEmptyResourceName
	}

	if s.connectionsStore.Exists(connection.Name) {
		return ErrConnectionAlreadyExists
	}

	now := time.Now().UTC()
	connection.ARN = s.buildConnectionARN(connection.Name)
	connection.Region = s.region
	connection.AccountID = s.accountID
	connection.CreatedAt = now
	connection.LastModifiedAt = now
	if connection.State == "" {
		connection.State = ConnectionStateAuthorized
	}

	return s.connectionsStore.Put(connection.Name, connection)
}

// GetConnection retrieves a connection by name. A missing record is the
// not-found sentinel; any other storage fault surfaces as itself so the
// Cores can tell a permanent absence from a transient failure.
//
// Parameters:
//   - ctx: The context
//   - name: The connection name
//
// Returns:
//   - *Connection: The connection if found
//   - error: ErrConnectionNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetConnection(ctx context.Context, name string) (*Connection, error) {
	var connection Connection
	if err := s.connectionsStore.Get(name, &connection); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrConnectionNotFound
		}
		return nil, err
	}
	return &connection, nil
}

// DeleteConnection deletes a connection by name.
//
// Parameters:
//   - ctx: The context
//   - name: The connection name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteConnection(ctx context.Context, name string) error {
	// The record lock closes the resurrection window: without it a
	// MutateConnection cycle that read before the delete could write the
	// record back after it.
	connectionRecordWriteMu.Lock()
	defer connectionRecordWriteMu.Unlock()
	if !s.connectionsStore.Exists(name) {
		return ErrConnectionNotFound
	}
	// Tags are keyed by the stored ARN; the ARN carries a unique id minted
	// at creation, so rebuilding it here would target a different key.
	if connection, err := s.GetConnection(ctx, name); err == nil {
		_ = s.TagStore.Delete(connection.ARN)
	}
	return s.connectionsStore.Delete(name)
}

// MutateConnection applies fn to the connection record inside the write
// mutex so the whole read-modify-write cycle is atomic with respect to
// every other connection record writer (configuration updates and
// deauthorization run on separate store instances over the same Pebble
// keyspace, hence the package-scope lock). fn mutates the record in place;
// fn returning an error aborts the write.
//
// Parameters:
//   - ctx: The context
//   - name: The connection name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrConnectionNotFound if the connection does not exist, the
//     storage fault otherwise
func (s *EventsStore) MutateConnection(ctx context.Context, name string, fn func(*Connection) error) error {
	connectionRecordWriteMu.Lock()
	defer connectionRecordWriteMu.Unlock()

	var connection Connection
	if err := s.connectionsStore.Get(name, &connection); err != nil {
		if common.IsNotFound(err) {
			return ErrConnectionNotFound
		}
		return err
	}
	if err := fn(&connection); err != nil {
		return err
	}
	return s.connectionsStore.Put(name, &connection)
}

// ApiDestination operations

// CreateApiDestination creates a new API destination.
//
// Parameters:
//   - ctx: The context
//   - apiDest: The API destination to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateApiDestination(ctx context.Context, apiDest *ApiDestination) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if apiDest.Name == "" {
		return ErrEmptyResourceName
	}

	if s.apiDestinationsStore.Exists(apiDest.Name) {
		return ErrApiDestinationAlreadyExists
	}

	now := time.Now().UTC()
	apiDest.ARN = s.buildApiDestinationARN(apiDest.Name)
	apiDest.Region = s.region
	apiDest.AccountID = s.accountID
	apiDest.CreatedAt = now
	apiDest.LastModifiedAt = now
	if apiDest.State == "" {
		apiDest.State = ApiDestinationStateActive
	}

	return s.apiDestinationsStore.Put(apiDest.Name, apiDest)
}

// GetApiDestination retrieves an API destination by name. A missing record
// is the not-found sentinel; any other storage fault surfaces as itself so
// the Cores can tell a permanent absence from a transient failure.
//
// Parameters:
//   - ctx: The context
//   - name: The API destination name
//
// Returns:
//   - *ApiDestination: The API destination if found
//   - error: ErrApiDestinationNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetApiDestination(ctx context.Context, name string) (*ApiDestination, error) {
	var apiDest ApiDestination
	if err := s.apiDestinationsStore.Get(name, &apiDest); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrApiDestinationNotFound
		}
		return nil, err
	}
	return &apiDest, nil
}

// DeleteApiDestination deletes an API destination by name.
//
// Parameters:
//   - ctx: The context
//   - name: The API destination name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteApiDestination(ctx context.Context, name string) error {
	// The record lock closes the resurrection window against a concurrent
	// MutateApiDestination cycle.
	apiDestinationRecordWriteMu.Lock()
	defer apiDestinationRecordWriteMu.Unlock()
	if !s.apiDestinationsStore.Exists(name) {
		return ErrApiDestinationNotFound
	}
	// Tags are keyed by the stored ARN; the ARN carries a unique id minted
	// at creation, so rebuilding it here would target a different key.
	if apiDest, err := s.GetApiDestination(ctx, name); err == nil {
		_ = s.TagStore.Delete(apiDest.ARN)
	}
	return s.apiDestinationsStore.Delete(name)
}

// MutateApiDestination applies fn to the API destination record inside the
// write mutex so the whole read-modify-write cycle is atomic with respect
// to every other API destination record writer (updaters run on separate
// store instances over the same Pebble keyspace, hence the package-scope
// lock). fn mutates the record in place; fn returning an error aborts the
// write.
//
// Parameters:
//   - ctx: The context
//   - name: The API destination name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrApiDestinationNotFound if the API destination does not
//     exist, the storage fault otherwise
func (s *EventsStore) MutateApiDestination(ctx context.Context, name string, fn func(*ApiDestination) error) error {
	apiDestinationRecordWriteMu.Lock()
	defer apiDestinationRecordWriteMu.Unlock()

	var apiDest ApiDestination
	if err := s.apiDestinationsStore.Get(name, &apiDest); err != nil {
		if common.IsNotFound(err) {
			return ErrApiDestinationNotFound
		}
		return err
	}
	if err := fn(&apiDest); err != nil {
		return err
	}
	return s.apiDestinationsStore.Put(name, &apiDest)
}

func (s *EventsStore) buildReplayARN(name string) string {
	return s.arnBuilder.Events().Replay(name)
}

// StoreArchiveEvent stores an event in an archive. The write runs under the
// archive family's record lock and re-checks the archive's existence, so an
// ingress racing the archive's delete cannot land an event row under a dead
// archive name — a row a same-named archive created later would surface as
// its own.
//
// Parameters:
//   - ctx: The context
//   - archiveName: The archive name
//   - event: The archived event to store
//
// Returns:
//   - error: ErrArchiveNotFound when the archive record is gone, the storage
//     fault otherwise
func (s *EventsStore) StoreArchiveEvent(ctx context.Context, archiveName string, event *ArchivedEvent) error {
	archiveRecordWriteMu.Lock()
	defer archiveRecordWriteMu.Unlock()
	if !s.archivesStore.Exists(archiveName) {
		return ErrArchiveNotFound
	}
	key := fmt.Sprintf("%s:%d:%s", archiveName, event.Timestamp.UnixNano(), event.ID)
	return s.archivedEventsStore.Put(key, event)
}

// ListArchiveEvents retrieves one page of archived events for an archive
// within a time range (inclusive bounds).
//
// Parameters:
//   - ctx: The context
//   - archiveName: The archive name
//   - startTime: The start time for retrieval
//   - endTime: The end time for retrieval
//   - limit: The maximum number of events to return
//   - nextToken: The continuation token from the previous page
//
// Returns:
//   - *ArchivedEventListResult: One page of archived events plus the
//     continuation token
//   - error: An error if retrieval fails
//
// The scan is page-bounded so a replay walks its window in pages instead of
// materialising the archive's entire filtered set before delivering the
// first event.
func (s *EventsStore) ListArchiveEvents(ctx context.Context, archiveName string, startTime, endTime time.Time, limit int32, nextToken string) (*ArchivedEventListResult, error) {
	opts := common.ListOptions{
		Prefix:   archiveName + ":",
		Marker:   nextToken,
		MaxItems: int(limit),
	}
	result, err := common.List[ArchivedEvent](s.archivedEventsStore, opts, func(e *ArchivedEvent) bool {
		return (startTime.IsZero() || !e.Timestamp.Before(startTime)) &&
			(endTime.IsZero() || !e.Timestamp.After(endTime))
	})
	if err != nil {
		return nil, err
	}
	return &ArchivedEventListResult{
		Events:    result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// DeleteExpiredArchiveEvents deletes archived events older than the cutoff
// time and decrements the archive counters accordingly.
func (s *EventsStore) DeleteExpiredArchiveEvents(ctx context.Context, archiveName string, cutoff time.Time) error {
	prefix := archiveName + ":"
	return common.ForEachAll[ArchivedEvent](s.archivedEventsStore, prefix, func(e *ArchivedEvent) bool {
		return e.Timestamp.Before(cutoff)
	}, func(e *ArchivedEvent) error {
		key := prefix + fmt.Sprintf("%d:", e.Timestamp.UnixNano()) + e.ID
		if err := s.archivedEventsStore.Delete(key); err != nil {
			return err
		}
		// The counter decrement runs through the atomic record mutation;
		// a concurrently deleted archive is a legitimate outcome of the
		// sweep (the retention worker iterates live archives), so its
		// not-found error is tolerated rather than aborting the pass.
		if err := s.MutateArchive(ctx, archiveName, func(archive *Archive) error {
			if archive.EventCount > 0 {
				archive.EventCount--
			}
			eventSize := int64(0)
			if eventBytes, err := json.Marshal(e.Event); err == nil {
				eventSize = int64(len(eventBytes))
			}
			if archive.SizeBytes >= eventSize {
				archive.SizeBytes -= eventSize
			} else {
				archive.SizeBytes = 0
			}
			return nil
		}); err != nil && err != ErrArchiveNotFound {
			return err
		}
		return nil
	})
}

// CreateReplay creates a new event replay.
//
// Parameters:
//   - ctx: The context
//   - replay: The replay to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateReplay(ctx context.Context, replay *Replay) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	return s.createReplayLocked(ctx, replay)
}

// CreateReplayCapped counts the store's active (non-terminal) replays and
// creates the record inside one locked section, so the concurrent-replay cap
// holds under concurrent creators: the count-then-create pair is atomic with
// every other capped create and with the state transitions that free cap
// slots (the replay family's record lock serialises them all). Lock order is
// createMu → replayRecordWriteMu; no other path nests these.
//
// Parameters:
//   - ctx: The context
//   - replay: The replay to create
//   - maxActive: The maximum of non-terminal replays the store may hold
//
// Returns:
//   - error: ErrReplayCapReached when maxActive is already reached, the
//     CreateReplay errors otherwise
func (s *EventsStore) CreateReplayCapped(ctx context.Context, replay *Replay, maxActive int) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	replayRecordWriteMu.Lock()
	defer replayRecordWriteMu.Unlock()

	active := 0
	nextToken := ""
	for {
		result, err := s.ListReplays(ctx, "", "", "", ListLimitMaximum, nextToken)
		if err != nil {
			return err
		}
		for _, r := range result.Replays {
			switch r.State {
			case ReplayStateStarting, ReplayStateRunning, ReplayStateCancelling:
				active++
			}
		}
		if result.NextToken == "" {
			break
		}
		nextToken = result.NextToken
	}
	if active >= maxActive {
		return ErrReplayCapReached
	}
	return s.createReplayLocked(ctx, replay)
}

// createReplayLocked is the create body shared by CreateReplay and
// CreateReplayCapped; the caller holds the create lock.
func (s *EventsStore) createReplayLocked(ctx context.Context, replay *Replay) error {
	if replay.Name == "" {
		return ErrEmptyResourceName
	}

	if s.replaysStore.Exists(replay.Name) {
		return ErrReplayAlreadyExists
	}

	replay.ARN = s.buildReplayARN(replay.Name)
	replay.Region = s.region
	replay.AccountID = s.accountID
	replay.CreatedAt = time.Now().UTC()
	if replay.State == "" {
		replay.State = ReplayStateStarting
	}

	return s.replaysStore.Put(replay.Name, replay)
}

// GetReplay retrieves a replay by name. A missing record is the not-found
// sentinel; any other storage fault surfaces as itself so the Cores can
// tell a permanent absence from a transient failure.
//
// Parameters:
//   - ctx: The context
//   - name: The replay name
//
// Returns:
//   - *Replay: The replay if found
//   - error: ErrReplayNotFound if not found, the storage fault otherwise
func (s *EventsStore) GetReplay(ctx context.Context, name string) (*Replay, error) {
	var replay Replay
	if err := s.replaysStore.Get(name, &replay); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrReplayNotFound
		}
		return nil, err
	}
	return &replay, nil
}

// MutateReplay applies fn to the replay record inside the write mutex so
// the whole read-modify-write cycle is atomic with respect to every other
// replay record writer: the replay worker's state transitions, the cancel
// path and the delete all run on separate store instances over the same
// Pebble keyspace, hence the package-scope lock — a terminal-write race
// between cancel and complete is decided by lock order, not by which
// goroutine re-read first. fn mutates the record in place; fn returning an
// error aborts the write.
//
// Parameters:
//   - ctx: The context
//   - name: The replay name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrReplayNotFound if the replay does not exist, the storage
//     fault otherwise
func (s *EventsStore) MutateReplay(ctx context.Context, name string, fn func(*Replay) error) error {
	replayRecordWriteMu.Lock()
	defer replayRecordWriteMu.Unlock()

	var replay Replay
	if err := s.replaysStore.Get(name, &replay); err != nil {
		if common.IsNotFound(err) {
			return ErrReplayNotFound
		}
		return err
	}
	if err := fn(&replay); err != nil {
		return err
	}
	return s.replaysStore.Put(name, &replay)
}

// DeleteReplay deletes a replay by name.
//
// Parameters:
//   - ctx: The context
//   - name: The replay name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteReplay(ctx context.Context, name string) error {
	// The record lock closes the resurrection window against a concurrent
	// MutateReplay cycle.
	replayRecordWriteMu.Lock()
	defer replayRecordWriteMu.Unlock()
	if !s.replaysStore.Exists(name) {
		return ErrReplayNotFound
	}
	return s.replaysStore.Delete(name)
}

// ConnectionListResult represents the result of listing connections.
type ConnectionListResult struct {
	Connections []*Connection
	NextToken   string
}

// ApiDestinationListResult represents the result of listing API destinations.
type ApiDestinationListResult struct {
	ApiDestinations []*ApiDestination
	NextToken       string
}

// ReplayListResult represents the result of listing replays.
type ReplayListResult struct {
	Replays   []*Replay
	NextToken string
}

// ArchivedEventListResult represents one page of archived events.
type ArchivedEventListResult struct {
	Events    []*ArchivedEvent
	NextToken string
}

// ListReplays lists replays with optional filtering.
//
// Parameters:
//   - ctx: The context
//   - namePrefix: Optional name prefix filter
//   - state: Optional state filter
//   - limit: Maximum number of results
//   - nextToken: Pagination token
//
// Returns:
//   - *ReplayListResult: The list result with replays and next token
//   - error: An error if listing fails
func (s *EventsStore) ListReplays(ctx context.Context, namePrefix string, state ReplayState, eventSourceArn string, limit int32, nextToken string) (*ReplayListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Replay](s.replaysStore, opts, func(r *Replay) bool {
		if namePrefix != "" && !strings.HasPrefix(r.Name, namePrefix) {
			return false
		}
		if state != "" && r.State != state {
			return false
		}
		if eventSourceArn != "" && r.EventSourceARN != eventSourceArn {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &ReplayListResult{
		Replays:   result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// ListArchives lists archives with optional filtering.
func (s *EventsStore) ListArchives(ctx context.Context, namePrefix, eventSourceArn, state string, limit int32, nextToken string) (*ArchiveListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Archive](s.archivesStore, opts, func(a *Archive) bool {
		if namePrefix != "" && !strings.HasPrefix(a.Name, namePrefix) {
			return false
		}
		if eventSourceArn != "" && a.EventSourceARN != eventSourceArn {
			return false
		}
		if state != "" && string(a.State) != state {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &ArchiveListResult{
		Archives:  result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// ListConnections lists connections with optional filtering.
func (s *EventsStore) ListConnections(ctx context.Context, namePrefix string, state string, limit int32, nextToken string) (*ConnectionListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Connection](s.connectionsStore, opts, func(c *Connection) bool {
		if namePrefix != "" && !strings.HasPrefix(c.Name, namePrefix) {
			return false
		}
		if state != "" && string(c.State) != state {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &ConnectionListResult{
		Connections: result.Items,
		NextToken:   result.NextMarker,
	}, nil
}

// ListApiDestinations lists API destinations with optional filtering.
func (s *EventsStore) ListApiDestinations(ctx context.Context, namePrefix, connectionArn string, limit int32, nextToken string) (*ApiDestinationListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[ApiDestination](s.apiDestinationsStore, opts, func(d *ApiDestination) bool {
		if namePrefix != "" && !strings.HasPrefix(d.Name, namePrefix) {
			return false
		}
		if connectionArn != "" && d.ConnectionARN != connectionArn {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &ApiDestinationListResult{
		ApiDestinations: result.Items,
		NextToken:       result.NextMarker,
	}, nil
}
