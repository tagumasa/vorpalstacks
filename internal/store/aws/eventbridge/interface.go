package eventbridge

import (
	"context"
	"time"
)

// EventsStoreInterface defines operations for managing EventBridge resources.
type EventsStoreInterface interface {
	GetAccountID() string
	GetRegion() string

	CreateEventBus(ctx context.Context, eventBus *EventBus) error
	GetEventBus(ctx context.Context, name string) (*EventBus, error)
	UpdateEventBus(ctx context.Context, eventBus *EventBus) error
	DeleteEventBus(ctx context.Context, name string) error
	ListEventBuses(ctx context.Context, namePrefix string, limit int32, nextToken string) (*EventBusListResult, error)

	CreateRule(ctx context.Context, rule *Rule) error
	GetRule(ctx context.Context, eventBusName, name string) (*Rule, error)
	UpdateRule(ctx context.Context, rule *Rule) error
	DeleteRule(ctx context.Context, eventBusName, name string) error
	ListRules(ctx context.Context, eventBusName, namePrefix string, limit int32, nextToken string) (*RuleListResult, error)

	PutTarget(ctx context.Context, target *Target) error
	GetTarget(ctx context.Context, eventBusName, ruleName, targetID string) (*Target, error)
	DeleteTarget(ctx context.Context, eventBusName, ruleName, targetID string) error
	ListTargetsByRule(ctx context.Context, eventBusName, ruleName string, limit int32, nextToken string) (*TargetListResult, error)
	DeleteTargetsByRule(ctx context.Context, eventBusName, ruleName string, targetIDs []string) error

	CreateArchive(ctx context.Context, archive *Archive) error
	GetArchive(ctx context.Context, name string) (*Archive, error)
	DeleteArchive(ctx context.Context, name string) error
	UpdateArchive(ctx context.Context, archive *Archive) error
	ListArchivesForEventBus(ctx context.Context, eventBusName string) ([]*Archive, error)
	ListArchives(ctx context.Context, namePrefix string, state string, limit int32, nextToken string) (*ArchiveListResult, error)

	StoreArchiveEvent(ctx context.Context, archiveName string, event *ArchivedEvent) error
	GetArchiveEvents(ctx context.Context, archiveName string, startTime, endTime time.Time) ([]*ArchivedEvent, error)
	DeleteArchiveEvents(ctx context.Context, archiveName string) error

	CreateReplay(ctx context.Context, replay *Replay) error
	GetReplay(ctx context.Context, name string) (*Replay, error)
	UpdateReplay(ctx context.Context, replay *Replay) error
	DeleteReplay(ctx context.Context, name string) error
	ListReplays(ctx context.Context, namePrefix string, state ReplayState, limit int32, nextToken string) (*ReplayListResult, error)

	CreateConnection(ctx context.Context, connection *Connection) error
	GetConnection(ctx context.Context, name string) (*Connection, error)
	DeleteConnection(ctx context.Context, name string) error
	UpdateConnection(ctx context.Context, connection *Connection) error
	DeauthorizeConnection(ctx context.Context, name string) error
	ListConnections(ctx context.Context, namePrefix string, state string, limit int32, nextToken string) (*ConnectionListResult, error)

	CreateApiDestination(ctx context.Context, apiDest *ApiDestination) error
	GetApiDestination(ctx context.Context, name string) (*ApiDestination, error)
	DeleteApiDestination(ctx context.Context, name string) error
	UpdateApiDestination(ctx context.Context, apiDest *ApiDestination) error
	ListApiDestinations(ctx context.Context, namePrefix string, limit int32, nextToken string) (*ApiDestinationListResult, error)
}
