package dynamodb

import (
	"context"

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
// The wire members are extracted tolerantly — an absent specification leaves
// SpecPresent false, an Enabled of the wrong JSON type sets Malformed, an
// absent or non-string AttributeName leaves it empty — and updateTimeToLiveCore,
// not the handler, rejects each of those.
type UpdateTimeToLiveInput struct {
	TableName     string
	Enabled       bool
	AttributeName string
	SpecPresent   bool // TimeToLiveSpecification was present on the wire
	Malformed     bool // a specification member carried the wrong JSON type
}

// parseTimeToLiveSpec extracts the wire's TimeToLiveSpecification into the
// Core input. It never rejects: absent members take their zero values so the
// rejection stays the Core's decision, keeping the HTTP handler a transport
// adapter on the TTL path.
func parseTimeToLiveSpec(v interface{}) UpdateTimeToLiveInput {
	var in UpdateTimeToLiveInput
	ttlSpec, ok := v.(map[string]interface{})
	if !ok {
		return in
	}
	in.SpecPresent = true
	if b, ok := ttlSpec["Enabled"].(bool); ok {
		in.Enabled = b
	} else if _, present := ttlSpec["Enabled"]; present {
		in.Malformed = true
	}
	in.AttributeName, _ = ttlSpec["AttributeName"].(string)
	return in
}

// updateTimeToLiveCore validates the requested TTL specification and
// persists it. The enable or disable is effective when the call returns —
// the store applies expiry atomically, so the status is written in its
// final state (ENABLED / DISABLED) rather than through a transition
// window the single-node store never actually has.
func (s *DynamoDBService) updateTimeToLiveCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, in UpdateTimeToLiveInput) (*dbstore.TimeToLiveSpecification, error) {
	if !in.SpecPresent || in.Malformed || in.AttributeName == "" {
		return nil, ErrInvalidParameter
	}
	// The TTL attribute name carries the model's 1-255 length constraint.
	if !validateTimeToLiveAttributeName(in.AttributeName) {
		return nil, ErrInvalidParameter
	}
	// Enabling TTL on a table that already has TTL enabled is rejected;
	// renaming the TTL attribute requires disabling TTL first. Disabling
	// is always allowed and is the documented path for changing the
	// attribute.
	existingTTL, _ := store.Tables().GetTimeToLive(in.TableName)
	if in.Enabled && existingTTL != nil && existingTTL.Enabled {
		return nil, ErrInvalidParameter
	}

	ttl := &dbstore.TimeToLiveSpecification{
		Enabled:       in.Enabled,
		AttributeName: in.AttributeName,
	}
	if in.Enabled {
		ttl.Status = dbstore.TTLStatusEnabled
	} else {
		ttl.Status = dbstore.TTLStatusDisabled
	}

	if err := store.Tables().SetTimeToLive(in.TableName, ttl); err != nil {
		return nil, err
	}

	return ttl, nil
}
