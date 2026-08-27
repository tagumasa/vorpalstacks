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
