package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// TTL Core — single validation + persistence path
//
// These methods encapsulate TTL lifecycle logic. Both the HTTP API handlers
// (ttl_operations.go) and any future admin handler delegate to these
// methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// describeTimeToLiveCore returns the TTL specification for the named table.
// A nil specification means TTL has never been configured.
func (s *DynamoDBService) describeTimeToLiveCore(store dbstore.DynamoDBStoreInterface, tableName string) (*dbstore.TimeToLiveSpecification, error) {
	return store.Tables().GetTimeToLive(tableName)
}

// UpdateTimeToLiveInput is the service-layer DTO for updating TTL settings.
type UpdateTimeToLiveInput struct {
	TableName     string
	Enabled       bool
	AttributeName string
}

// updateTimeToLiveCore persists the requested TTL specification and starts
// the background status transition (ENABLING -> ENABLED / DISABLING ->
// DISABLED). It returns the persisted specification in its initial
// transition state.
func (s *DynamoDBService) updateTimeToLiveCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, in UpdateTimeToLiveInput) (*dbstore.TimeToLiveSpecification, error) {
	ttl := &dbstore.TimeToLiveSpecification{
		Enabled:       in.Enabled,
		AttributeName: in.AttributeName,
	}
	if in.Enabled {
		ttl.Status = dbstore.TTLStatusEnabling
	} else {
		ttl.Status = dbstore.TTLStatusDisabling
	}

	if err := store.Tables().SetTimeToLive(in.TableName, ttl); err != nil {
		return nil, err
	}

	// Background transition to final state.
	s.bgWg.Add(1)
	go func() {
		defer func() { resilience.RecoverPanic("dynamodb TTL state transition") }()
		defer s.bgWg.Done()
		select {
		case <-time.After(1 * time.Second):
		case <-s.bgCtx.Done():
			return
		}
		if in.Enabled {
			ttl.Status = dbstore.TTLStatusEnabled
		} else {
			ttl.Status = dbstore.TTLStatusDisabled
		}
		if err := store.Tables().SetTimeToLive(in.TableName, ttl); err != nil {
			logs.Error("Failed to transition TTL to final state",
				logs.Err(err),
				logs.String("tableName", in.TableName),
			)
		}
	}()

	return ttl, nil
}
