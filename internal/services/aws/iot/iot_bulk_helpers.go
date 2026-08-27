package iot

import (
	"vorpalstacks/internal/common/request"
)

// Transport adapters for the generic-KV Core functions in
// iot_bulk_core.go. Handlers for Security, Detect, Audit, Fleet Indexing,
// Logging, Encryption and Topic Rule Destinations. State is persisted via
// the generic-KV store (Pebble-backed), so it survives restarts. The CRUD
// entities (CustomMetric, Dimension, MitigationAction, FleetMetric,
// ScheduledAudit, TopicRuleDestination) implement full
// create/read/update/delete with field echo, NotFound semantics and error
// propagation, matching the AWS IoT Control Plane. The asynchronous task
// families (Detect tasks, on-demand Audit tasks) and the fleet index
// aggregator require dedicated engines that are out of scope; those
// handlers remain structural stubs.

// ---- bulk generic-KV transport adapters --------------------------------------

// bulkCreate persists a new record under "<category>/<name>", capturing the
// supplied extra fields alongside the name and timestamps. Returns the stored
// record so callers can shape the AWS response. Store errors are propagated.
func (s *IoTService) bulkCreate(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string, extra map[string]interface{}) (map[string]interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetParamCaseInsensitive(req.Parameters, nameKey)
	return s.bulkCreateCore(store, category, name, extra)
}

// bulkGet loads a single record. exists is false (with nil error) and rec
// is nil when the record is absent.
func (s *IoTService) bulkGet(reqCtx *request.RequestContext, category string, req *request.ParsedRequest, nameKey string) (rec map[string]interface{}, name string, exists bool, err error) {
	name = request.GetParamCaseInsensitive(req.Parameters, nameKey)
	store, sErr := s.store(reqCtx)
	if sErr != nil {
		return nil, name, false, sErr
	}
	rec, exists, err = s.bulkGetCore(store, category, name)
	return rec, name, exists, err
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
	return s.bulkUpdateCore(store, category, name, merge)
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
	return s.bulkDeleteCore(store, category, name)
}

// bulkList lists all records under a category prefix.
func (s *IoTService) bulkList(reqCtx *request.RequestContext, category string) ([]map[string]interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.bulkListCore(store, category)
}
