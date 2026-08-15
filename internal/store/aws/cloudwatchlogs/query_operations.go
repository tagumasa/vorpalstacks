package cloudwatchlogs

import (
	"encoding/json"
	"strings"
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

func (s *Store) DeleteScheduledQuery(id string) error {
	key := s.scheduledQueryKey(id)
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
