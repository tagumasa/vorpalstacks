package iot

import (
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core persistence layer for the generic-KV CRUD families (CustomMetric,
// Dimension, MitigationAction, FleetMetric, ScheduledAudit, TopicRule
// Destination). State is persisted via the generic-KV store
// (Pebble-backed), so it survives restarts. The CRUD entities implement
// full create/read/update/delete with field echo, NotFound semantics and
// error propagation, matching the AWS IoT Control Plane. The asynchronous
// task families (Detect tasks, on-demand Audit tasks) and the fleet index
// aggregator require dedicated engines that are out of scope; those
// handlers remain structural stubs.

// bulkCreateCore persists a new record under "<category>/<name>",
// capturing the supplied extra fields alongside the name and timestamps.
// Returns the stored record so callers can shape the AWS response. Store
// errors are propagated.
func (s *IoTService) bulkCreateCore(store iotstore.IotStoreInterface, category, name string, extra map[string]interface{}) (map[string]interface{}, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Reject duplicates: AWS returns ResourceAlreadyExistsException (409)
	// when a named resource already exists.
	exists, err := store.GetGenericExists(category+"/"+name, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, iotstore.ErrResourceAlreadyExists
	}
	now := time.Now().Unix()
	rec := map[string]interface{}{
		"name":             name,
		"creationDate":     now,
		"lastModifiedDate": now,
	}
	for k, v := range extra {
		if v != nil {
			rec[k] = v
		}
	}
	if err := store.PutGeneric(category+"/"+name, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// bulkCreateIdempotentCore behaves like bulkCreateCore but honours the
// idempotencyToken trait the create operations carry on their
// clientRequestToken member: replaying the token the existing record was
// created with returns the stored record instead of the duplicate conflict.
func (s *IoTService) bulkCreateIdempotentCore(store iotstore.IotStoreInterface, category, name, clientRequestToken string, extra map[string]interface{}) (map[string]interface{}, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	existing := map[string]interface{}{}
	exists, err := store.GetGenericExists(category+"/"+name, &existing)
	if err != nil {
		return nil, err
	}
	if exists {
		if clientRequestToken != "" && existing["clientRequestToken"] == clientRequestToken {
			return existing, nil
		}
		return nil, iotstore.ErrResourceAlreadyExists
	}
	now := time.Now().Unix()
	rec := map[string]interface{}{
		"name":             name,
		"creationDate":     now,
		"lastModifiedDate": now,
	}
	for k, v := range extra {
		if v != nil {
			rec[k] = v
		}
	}
	if err := store.PutGeneric(category+"/"+name, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// claimClientToken records the clientRequestToken-to-resource-key index the
// dimension and audit-suppression creates require: their model documentation
// states that reusing an existing resource's token for a new resource raises
// an exception. Re-applying a token against the resource key it already
// indexes passes through, preserving the idempotent replay of the same
// create request. The stored index maps the token to the resource record
// key, so the release side needs no second lookup structure.
func (s *IoTService) claimClientToken(store iotstore.IotStoreInterface, kind, token, resourceKey string) error {
	indexKey := "clientToken/" + kind + "/" + token
	index := map[string]interface{}{}
	exists, err := store.GetGenericExists(indexKey, &index)
	if err != nil {
		return err
	}
	if exists {
		if indexed, _ := index["resourceKey"].(string); indexed == resourceKey {
			return nil
		}
		return iotstore.ErrResourceAlreadyExists
	}
	return store.PutGeneric(indexKey, map[string]interface{}{"resourceKey": resourceKey})
}

// releaseClientToken drops the token index when the resource that claimed
// it is deleted, so the token may be reused for a later create.
func (s *IoTService) releaseClientToken(store iotstore.IotStoreInterface, kind, token string) error {
	if token == "" {
		return nil
	}
	return store.DeleteGeneric("clientToken/" + kind + "/" + token)
}

// bulkGetCore loads a single record. exists is false (with nil error) and
// rec is nil when the record is absent.
func (s *IoTService) bulkGetCore(store iotstore.IotStoreInterface, category, name string) (rec map[string]interface{}, exists bool, err error) {
	rec = map[string]interface{}{}
	exists, err = store.GetGenericExists(category+"/"+name, &rec)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return rec, true, nil
}

// bulkUpdateCore merges the supplied fields into an existing record and
// refreshes lastModifiedDate. exists is false (nil error) when absent, so
// callers return ResourceNotFoundException.
func (s *IoTService) bulkUpdateCore(store iotstore.IotStoreInterface, category, name string, merge map[string]interface{}) (map[string]interface{}, bool, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(category+"/"+name, &rec)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	for k, v := range merge {
		if v != nil {
			rec[k] = v
		}
	}
	rec["lastModifiedDate"] = time.Now().Unix()
	if err := store.PutGeneric(category+"/"+name, rec); err != nil {
		return nil, false, err
	}
	return rec, true, nil
}

// bulkDeleteCore removes a record. BaseStore.Delete is idempotent (no
// error on a missing key), matching AWS IoT delete semantics; only
// genuine store errors propagate.
func (s *IoTService) bulkDeleteCore(store iotstore.IotStoreInterface, category, name string) error {
	return store.DeleteGeneric(category + "/" + name)
}

// bulkListCore lists all records under a category prefix.
func (s *IoTService) bulkListCore(store iotstore.IotStoreInterface, category string) ([]map[string]interface{}, error) {
	return store.ListGeneric(category + "/")
}

func bulkName(rec map[string]interface{}) string {
	if name, ok := rec["name"].(string); ok {
		return name
	}
	return ""
}
