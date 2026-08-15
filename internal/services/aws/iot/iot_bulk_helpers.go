package iot

import (
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Handlers for Security, Detect, Audit, Fleet Indexing, Logging,
// Encryption and Topic Rule Destinations. State is persisted via the generic-KV
// store (Pebble-backed), so it survives restarts. The CRUD entities
// (CustomMetric, Dimension, MitigationAction, FleetMetric, ScheduledAudit,
// TopicRuleDestination) implement full create/read/update/delete with field
// echo, NotFound semantics and error propagation, matching the AWS IoT
// Control Plane. The asynchronous task families (Detect tasks, on-demand
// Audit tasks) and the fleet index aggregator require dedicated engines that
// are out of scope; those handlers remain structural stubs.

// ---- bulk generic-KV helpers ------------------------------------------------

// bulkCreate persists a new record under "<category>/<name>", capturing the
// supplied extra fields alongside the name and timestamps. Returns the stored
// record so callers can shape the AWS response. Store errors are propagated.
func (s *IoTService) bulkCreate(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string, extra map[string]interface{}) (map[string]interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Reject duplicates: AWS returns ResourceAlreadyExistsException (409)
	// when a named resource already exists. Only check when the caller
	// supplied a name (auto-generated UUIDs are inherently unique).
	if userProvided := request.GetParamCaseInsensitive(req.Parameters, nameKey); userProvided != "" {
		exists, err := store.GetGenericExists(category+"/"+name, &map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, iotstore.ErrResourceAlreadyExists
		}
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

// bulkGet loads a single record. exists is false (with nil error) and rec
// is nil when the record is absent.
func (s *IoTService) bulkGet(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string) (rec map[string]interface{}, name string, exists bool, err error) {
	name = request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, sErr := s.store(reqCtx)
	if sErr != nil {
		return nil, name, false, sErr
	}
	rec = map[string]interface{}{}
	exists, err = store.GetGenericExists(category+"/"+name, &rec)
	if err != nil {
		return nil, name, false, err
	}
	if !exists {
		return nil, name, false, nil
	}
	return rec, name, true, nil
}

// bulkUpdate merges the supplied fields into an existing record and refreshes
// lastModifiedDate. exists is false (nil error) when absent, so callers return
// ResourceNotFoundException.
func (s *IoTService) bulkUpdate(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string, merge map[string]interface{}) (map[string]interface{}, bool, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, false, err
	}
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

// bulkDelete removes a record. BaseStore.Delete is idempotent (no error on a
// missing key), matching AWS IoT delete semantics; only genuine store errors
// propagate.
func (s *IoTService) bulkDelete(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string) error {
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	return store.DeleteGeneric(category + "/" + name)
}

// bulkList lists all records under a category prefix.
func (s *IoTService) bulkList(reqCtx *request.RequestContext, category string) ([]map[string]interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return store.ListGeneric(category + "/")
}

func bulkName(rec map[string]interface{}) string {
	if name, ok := rec["name"].(string); ok {
		return name
	}
	return ""
}
