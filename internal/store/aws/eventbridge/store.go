// Package events provides EventBridge storage functionality for vorpalstacks.
package eventbridge

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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
	*common.TagStore
	arnBuilder        *svcarn.ARNBuilder
	accountID         string
	region            string
	createMu          sync.Mutex
	archiveCountersMu sync.Mutex
}

// NewEventsStore creates a new EventBridge events store.
func NewEventsStore(store storage.BasicStorage, accountID, region string) *EventsStore {
	return &EventsStore{
		BaseStore:            common.NewBaseStore(store.Bucket("events-eventbuses-"+region), "events-eventbuses"),
		rulesStore:           common.NewBaseStore(store.Bucket("events-rules-"+region), "events-rules"),
		targetsStore:         common.NewBaseStore(store.Bucket("events-targets-"+region), "events-targets"),
		archivesStore:        common.NewBaseStore(store.Bucket("events-archives-"+region), "events-archives"),
		archivedEventsStore:  common.NewBaseStore(store.Bucket("events-archived-events-"+region), "events-archived-events"),
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
	return s.arnBuilder.Events().Connection(name)
}

func (s *EventsStore) buildApiDestinationARN(name string) string {
	return s.arnBuilder.Events().ApiDestination(name)
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
		return ErrInvalidARN
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

// GetEventBus retrieves an event bus by name.
//
// Parameters:
//   - ctx: The context
//   - name: The event bus name
//
// Returns:
//   - *EventBus: The event bus if found
//   - error: An error if not found
func (s *EventsStore) GetEventBus(ctx context.Context, name string) (*EventBus, error) {
	arn := s.buildEventBusARN(name)
	var eventBus EventBus
	if err := s.BaseStore.Get(arn, &eventBus); err != nil {
		return nil, ErrEventBusNotFound
	}
	return &eventBus, nil
}

// UpdateEventBus updates an existing event bus.
//
// Parameters:
//   - ctx: The context
//   - eventBus: The event bus to update
//
// Returns:
//   - error: An error if update fails
func (s *EventsStore) UpdateEventBus(ctx context.Context, eventBus *EventBus) error {
	if !s.Exists(eventBus.ARN) {
		return ErrEventBusNotFound
	}
	eventBus.LastModifiedAt = time.Now().UTC()
	return s.Put(eventBus.ARN, eventBus)
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
	if !s.Exists(arn) {
		return ErrEventBusNotFound
	}
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
func (s *EventsStore) ListEventBuses(ctx context.Context, namePrefix string, limit int32, nextToken string) (*EventBusListResult, error) {
	opts := common.ListOptions{
		Prefix:   "",
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[EventBus](s.BaseStore, opts, func(e *EventBus) bool {
		if namePrefix == "" {
			return true
		}
		return len(e.Name) >= len(namePrefix) && e.Name[:len(namePrefix)] == namePrefix
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

// CreateRule creates a new rule on an event bus.
//
// Parameters:
//   - ctx: The context
//   - rule: The rule to create
//
// Returns:
//   - error: An error if creation fails
func (s *EventsStore) CreateRule(ctx context.Context, rule *Rule) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if rule.Name == "" {
		return ErrInvalidARN
	}

	key := s.buildRuleKey(rule.EventBusName, rule.Name)
	if s.rulesStore.Exists(key) {
		return ErrRuleAlreadyExists
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

// GetRule retrieves a rule by event bus name and rule name.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - name: The rule name
//
// Returns:
//   - *Rule: The rule if found
//   - error: An error if not found
func (s *EventsStore) GetRule(ctx context.Context, eventBusName, name string) (*Rule, error) {
	key := s.buildRuleKey(eventBusName, name)
	var rule Rule
	if err := s.rulesStore.Get(key, &rule); err != nil {
		return nil, ErrRuleNotFound
	}
	return &rule, nil
}

// UpdateRule updates an existing rule.
//
// Parameters:
//   - ctx: The context
//   - rule: The rule to update
//
// Returns:
//   - error: An error if update fails
func (s *EventsStore) UpdateRule(ctx context.Context, rule *Rule) error {
	key := s.buildRuleKey(rule.EventBusName, rule.Name)
	if !s.rulesStore.Exists(key) {
		return ErrRuleNotFound
	}
	rule.LastModifiedAt = time.Now().UTC()
	return s.rulesStore.Put(key, rule)
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
	if !s.rulesStore.Exists(key) {
		return ErrRuleNotFound
	}
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
func (s *EventsStore) ListRules(ctx context.Context, eventBusName, namePrefix string, limit int32, nextToken string) (*RuleListResult, error) {
	prefix := eventBusName + ":"
	opts := common.ListOptions{
		Prefix:   prefix,
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Rule](s.rulesStore, opts, func(r *Rule) bool {
		if namePrefix == "" {
			return true
		}
		return len(r.Name) >= len(namePrefix) && r.Name[:len(namePrefix)] == namePrefix
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

// PutTarget adds or updates a target for a rule.
//
// Parameters:
//   - ctx: The context
//   - target: The target to put
//
// Returns:
//   - error: An error if the operation fails
func (s *EventsStore) PutTarget(ctx context.Context, target *Target) error {
	key := s.buildTargetKey(target.EventBusName, target.RuleName, target.ID)
	var existing Target
	if err := s.targetsStore.Get(key, &existing); err == nil && !existing.CreatedAt.IsZero() {
		target.CreatedAt = existing.CreatedAt
	} else {
		target.CreatedAt = time.Now().UTC()
	}
	return s.targetsStore.Put(key, target)
}

// GetTarget retrieves a target by event bus name, rule name, and target ID.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - ruleName: The rule name
//   - targetID: The target ID
//
// Returns:
//   - *Target: The target if found
//   - error: An error if not found
func (s *EventsStore) GetTarget(ctx context.Context, eventBusName, ruleName, targetID string) (*Target, error) {
	key := s.buildTargetKey(eventBusName, ruleName, targetID)
	var target Target
	if err := s.targetsStore.Get(key, &target); err != nil {
		return nil, ErrTargetNotFound
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

// DeleteTargetsByRule deletes multiple targets for a rule.
//
// Parameters:
//   - ctx: The context
//   - eventBusName: The event bus name
//   - ruleName: The rule name
//   - targetIDs: The list of target IDs to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteTargetsByRule(ctx context.Context, eventBusName, ruleName string, targetIDs []string) error {
	for _, targetID := range targetIDs {
		key := s.buildTargetKey(eventBusName, ruleName, targetID)
		if err := s.targetsStore.Delete(key); err != nil {
			return fmt.Errorf("failed to delete target %s: %w", targetID, err)
		}
	}
	return nil
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
		return ErrInvalidARN
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

	return s.archivesStore.Put(archive.Name, archive)
}

// GetArchive retrieves an archive by name.
//
// Parameters:
//   - ctx: The context
//   - name: The archive name
//
// Returns:
//   - *Archive: The archive if found
//   - error: An error if not found
func (s *EventsStore) GetArchive(ctx context.Context, name string) (*Archive, error) {
	var archive Archive
	if err := s.archivesStore.Get(name, &archive); err != nil {
		return nil, ErrArchiveNotFound
	}
	return &archive, nil
}

// DeleteArchive deletes an archive by name.
//
// Parameters:
//   - ctx: The context
//   - name: The archive name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteArchive(ctx context.Context, name string) error {
	if !s.archivesStore.Exists(name) {
		return ErrArchiveNotFound
	}
	return s.archivesStore.Delete(name)
}

// IncrementArchiveCounters atomically increments EventCount and SizeBytes for an archive.
func (s *EventsStore) IncrementArchiveCounters(ctx context.Context, archiveName string, eventSize int64) error {
	s.archiveCountersMu.Lock()
	defer s.archiveCountersMu.Unlock()

	var archive Archive
	if err := s.archivesStore.Get(archiveName, &archive); err != nil {
		return err
	}
	archive.EventCount++
	archive.SizeBytes += eventSize
	return s.archivesStore.Put(archiveName, &archive)
}

// UpdateArchive updates an existing archive.
//
// Parameters:
//   - ctx: The context
//   - archive: The archive to update
//
// Returns:
//   - error: An error if update fails
func (s *EventsStore) UpdateArchive(ctx context.Context, archive *Archive) error {
	if !s.archivesStore.Exists(archive.Name) {
		return ErrArchiveNotFound
	}
	return s.archivesStore.Put(archive.Name, archive)
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
	opts := common.ListOptions{MaxItems: 1000}
	result, err := common.List[Archive](s.archivesStore, opts, func(a *Archive) bool {
		return a.EventBusName == eventBusName
	})
	if err != nil {
		return nil, err
	}
	return result.Items, nil
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
		return ErrInvalidARN
	}

	if s.connectionsStore.Exists(connection.Name) {
		return ErrConnectionAlreadyExists
	}

	now := time.Now().UTC()
	connection.ARN = s.buildConnectionARN(connection.Name)
	connection.Region = s.region
	connection.AccountID = s.accountID
	connection.CreatedAt = now
	if connection.State == "" {
		connection.State = ConnectionStateAuthorized
	}

	return s.connectionsStore.Put(connection.Name, connection)
}

// GetConnection retrieves a connection by name.
//
// Parameters:
//   - ctx: The context
//   - name: The connection name
//
// Returns:
//   - *Connection: The connection if found
//   - error: An error if not found
func (s *EventsStore) GetConnection(ctx context.Context, name string) (*Connection, error) {
	var connection Connection
	if err := s.connectionsStore.Get(name, &connection); err != nil {
		return nil, ErrConnectionNotFound
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
	if !s.connectionsStore.Exists(name) {
		return ErrConnectionNotFound
	}
	return s.connectionsStore.Delete(name)
}

// UpdateConnection updates an existing EventBridge connection.
func (s *EventsStore) UpdateConnection(ctx context.Context, connection *Connection) error {
	if !s.connectionsStore.Exists(connection.Name) {
		return ErrConnectionNotFound
	}
	return s.connectionsStore.Put(connection.Name, connection)
}

// DeauthorizeConnection deauthorises an EventBridge connection, setting its state to deauthorised.
func (s *EventsStore) DeauthorizeConnection(ctx context.Context, name string) error {
	var connection Connection
	if err := s.connectionsStore.Get(name, &connection); err != nil {
		return ErrConnectionNotFound
	}
	connection.State = ConnectionStateDeauthorized
	connection.StateReason = "User initiated deauthorization"
	connection.LastAuthorizedAt = time.Time{}
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
		return ErrInvalidARN
	}

	if s.apiDestinationsStore.Exists(apiDest.Name) {
		return ErrApiDestinationAlreadyExists
	}

	now := time.Now().UTC()
	apiDest.ARN = s.buildApiDestinationARN(apiDest.Name)
	apiDest.Region = s.region
	apiDest.AccountID = s.accountID
	apiDest.CreatedAt = now
	if apiDest.State == "" {
		apiDest.State = ApiDestinationStateActive
	}

	return s.apiDestinationsStore.Put(apiDest.Name, apiDest)
}

// GetApiDestination retrieves an API destination by name.
//
// Parameters:
//   - ctx: The context
//   - name: The API destination name
//
// Returns:
//   - *ApiDestination: The API destination if found
//   - error: An error if not found
func (s *EventsStore) GetApiDestination(ctx context.Context, name string) (*ApiDestination, error) {
	var apiDest ApiDestination
	if err := s.apiDestinationsStore.Get(name, &apiDest); err != nil {
		return nil, ErrApiDestinationNotFound
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
	if !s.apiDestinationsStore.Exists(name) {
		return ErrApiDestinationNotFound
	}
	return s.apiDestinationsStore.Delete(name)
}

// UpdateApiDestination updates an existing EventBridge API destination.
func (s *EventsStore) UpdateApiDestination(ctx context.Context, apiDest *ApiDestination) error {
	if !s.apiDestinationsStore.Exists(apiDest.Name) {
		return ErrApiDestinationNotFound
	}
	return s.apiDestinationsStore.Put(apiDest.Name, apiDest)
}

func (s *EventsStore) buildReplayARN(name string) string {
	return s.arnBuilder.Events().Replay(name)
}

// StoreArchiveEvent stores an event in an archive.
//
// Parameters:
//   - ctx: The context
//   - archiveName: The archive name
//   - event: The archived event to store
//
// Returns:
//   - error: An error if storage fails
func (s *EventsStore) StoreArchiveEvent(ctx context.Context, archiveName string, event *ArchivedEvent) error {
	key := fmt.Sprintf("%s:%d:%s", archiveName, event.Timestamp.UnixNano(), event.ID)
	return s.archivedEventsStore.Put(key, event)
}

// GetArchiveEvents retrieves archived events for an archive within a time range.
//
// Parameters:
//   - ctx: The context
//   - archiveName: The archive name
//   - startTime: The start time for retrieval
//   - endTime: The end time for retrieval
//
// Returns:
//   - []*ArchivedEvent: The list of archived events
//   - error: An error if retrieval fails
func (s *EventsStore) GetArchiveEvents(ctx context.Context, archiveName string, startTime, endTime time.Time) ([]*ArchivedEvent, error) {
	prefix := archiveName + ":"
	opts := common.ListOptions{
		Prefix:   prefix,
		MaxItems: 10000,
	}

	result, err := common.List[ArchivedEvent](s.archivedEventsStore, opts, func(e *ArchivedEvent) bool {
		return (startTime.IsZero() || e.Timestamp.After(startTime) || e.Timestamp.Equal(startTime)) &&
			(endTime.IsZero() || e.Timestamp.Before(endTime) || e.Timestamp.Equal(endTime))
	})
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// DeleteArchiveEvents deletes all archived events for an archive.
//
// Parameters:
//   - ctx: The context
//   - archiveName: The archive name
//
// Returns:
//   - error: An error if deletion fails
func (s *EventsStore) DeleteArchiveEvents(ctx context.Context, archiveName string) error {
	prefix := archiveName + ":"
	opts := common.ListOptions{Prefix: prefix, MaxItems: 100000}
	result, err := common.List[ArchivedEvent](s.archivedEventsStore, opts, nil)
	if err != nil {
		return err
	}
	for _, e := range result.Items {
		key := prefix + fmt.Sprintf("%d:", e.Timestamp.UnixNano()) + e.ID
		if err := s.archivedEventsStore.Delete(key); err != nil {
			return err
		}
	}
	return nil
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
	if replay.Name == "" {
		return ErrInvalidARN
	}

	if s.replaysStore.Exists(replay.Name) {
		return ErrReplayAlreadyExists
	}

	replay.ARN = s.buildReplayARN(replay.Name)
	replay.Region = s.region
	replay.AccountID = s.accountID
	if replay.State == "" {
		replay.State = ReplayStateStarting
	}

	return s.replaysStore.Put(replay.Name, replay)
}

// GetReplay retrieves a replay by name.
//
// Parameters:
//   - ctx: The context
//   - name: The replay name
//
// Returns:
//   - *Replay: The replay if found
//   - error: An error if not found
func (s *EventsStore) GetReplay(ctx context.Context, name string) (*Replay, error) {
	var replay Replay
	if err := s.replaysStore.Get(name, &replay); err != nil {
		return nil, ErrReplayNotFound
	}
	return &replay, nil
}

// UpdateReplay updates an existing replay.
//
// Parameters:
//   - ctx: The context
//   - replay: The replay to update
//
// Returns:
//   - error: An error if update fails
func (s *EventsStore) UpdateReplay(ctx context.Context, replay *Replay) error {
	if !s.replaysStore.Exists(replay.Name) {
		return ErrReplayNotFound
	}
	return s.replaysStore.Put(replay.Name, replay)
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
func (s *EventsStore) ListReplays(ctx context.Context, namePrefix string, state ReplayState, limit int32, nextToken string) (*ReplayListResult, error) {
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
func (s *EventsStore) ListArchives(ctx context.Context, namePrefix string, state string, limit int32, nextToken string) (*ArchiveListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Archive](s.archivesStore, opts, func(a *Archive) bool {
		if namePrefix != "" && !strings.HasPrefix(a.Name, namePrefix) {
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
func (s *EventsStore) ListApiDestinations(ctx context.Context, namePrefix string, limit int32, nextToken string) (*ApiDestinationListResult, error) {
	opts := common.ListOptions{
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[ApiDestination](s.apiDestinationsStore, opts, func(d *ApiDestination) bool {
		if namePrefix != "" && !strings.HasPrefix(d.Name, namePrefix) {
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
