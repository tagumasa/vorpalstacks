package cloudwatchlogs

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
)

// --- Query Definition ---

func (s *Store) PutQueryDefinitionEntry(qd *QueryDefinition) error {
	qd.LastModified = time.Now().UTC().UnixMilli()
	return s.Put(s.queryDefinitionKey(qd.QueryDefinitionId), qd)
}

func (s *Store) DeleteQueryDefinitionEntry(id string) error {
	key := s.queryDefinitionKey(id)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListQueryDefinitions(namePrefix string) ([]*QueryDefinition, error) {
	var defs []*QueryDefinition
	if err := s.ScanPrefix("query-definition:", func(key string, value []byte) error {
		var qd QueryDefinition
		if err := json.Unmarshal(value, &qd); err != nil {
			return nil
		}
		if namePrefix == "" || strings.HasPrefix(qd.Name, namePrefix) {
			defs = append(defs, &qd)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return defs, nil
}

// --- Scheduled Query ---

func (s *Store) PutScheduledQuery(sq *ScheduledQuery) error {
	sq.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.scheduledQueryKey(sq.Id), sq)
}

func (s *Store) GetScheduledQuery(id string) (*ScheduledQuery, error) {
	var sq ScheduledQuery
	if err := s.Get(s.scheduledQueryKey(id), &sq); err != nil {
		return nil, ErrResourceNotFound
	}
	return &sq, nil
}

// scheduledQueryRecordWriteMu serialises read-modify-write cycles on
// scheduled query records. The delivery worker and the API handlers
// write these records concurrently: a full-record put from either side
// must not interleave with the other's read-modify-write, and a write
// racing a delete must not resurrect the deleted record.
var scheduledQueryRecordWriteMu sync.Mutex

// MutateScheduledQuery loads the stored scheduled query, applies fn, and
// persists the result as one atomic read-modify-write. The put advances
// LastUpdatedTime, so this is the path for user-driven mutations.
func (s *Store) MutateScheduledQuery(id string, fn func(*ScheduledQuery) error) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(id)
	if err != nil {
		return err
	}
	if err := fn(sq); err != nil {
		return err
	}
	return s.PutScheduledQuery(sq)
}

// TouchScheduledQueryDelivery records the outcome of a scheduled query
// execution: the consumed schedule boundary (the internal
// deduplication marker, only ever advancing), the execution clock
// (surfaced as lastTriggeredTime, the timestamp the query was last
// executed), and the execution status. It never advances
// LastUpdatedTime — an execution is not an update.
func (s *Store) TouchScheduledQueryDelivery(id string, boundary, executedAt int64, status string) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(id)
	if err != nil {
		return err
	}
	if boundary > sq.LastExecutedBoundary {
		sq.LastExecutedBoundary = boundary
	}
	sq.LastTriggeredTime = executedAt
	sq.LastExecutionStatus = status
	return s.Put(s.scheduledQueryKey(sq.Id), sq)
}

func (s *Store) DeleteScheduledQuery(id string) error {
	key := s.scheduledQueryKey(id)
	// The lock closes the resurrection window: a delivery touch or a
	// mutation that read the record before the delete must not write it
	// back afterwards.
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListScheduledQueries(state string) ([]*ScheduledQuery, error) {
	var queries []*ScheduledQuery
	if err := s.ScanPrefix("scheduled-query:", func(key string, value []byte) error {
		var sq ScheduledQuery
		if err := json.Unmarshal(value, &sq); err != nil {
			return nil
		}
		if state == "" || sq.State == state {
			queries = append(queries, &sq)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return queries, nil
}

// --- Scheduled Query Execution History ---

func (s *Store) PutScheduledQueryExecution(exec *ScheduledQueryExecution) error {
	return s.Put(s.scheduledQueryExecutionKey(exec.ScheduledQueryId, exec.TriggerTime), exec)
}

func (s *Store) ListScheduledQueryExecutions(sqId string, startTime, endTime int64) ([]*ScheduledQueryExecution, error) {
	var execs []*ScheduledQueryExecution
	prefix := "sq-execution:" + sqId + ":"
	if err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var exec ScheduledQueryExecution
		if err := json.Unmarshal(value, &exec); err != nil {
			return nil
		}
		if startTime > 0 && exec.TriggerTime < startTime {
			return nil
		}
		if endTime > 0 && exec.TriggerTime > endTime {
			return nil
		}
		execs = append(execs, &exec)
		return nil
	}); err != nil {
		return nil, err
	}
	return execs, nil
}
