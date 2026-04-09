// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import (
	"encoding/json"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"

	"github.com/google/uuid"
)

// EventSourceStore manages Lambda event source mappings.
type EventSourceStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	accountId  string
	region     string
	mu         sync.Mutex

	indexByKey          map[string]string
	indexByFunction     map[string][]string
	indexPathByKey      map[string]string
	indexPathByFunction map[string]string
	indexOnce           sync.Once
}

// NewEventSourceStore creates a new EventSourceStore.
func NewEventSourceStore(store storage.BasicStorage, accountId, region string) *EventSourceStore {
	bucket := store.Bucket("lambda-event-sources-" + region)
	return &EventSourceStore{
		BaseStore:  common.NewBaseStore(bucket, "lambda-event-sources"),
		arnBuilder: NewARNBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

func (s *EventSourceStore) ensureIndex() {
	s.indexOnce.Do(func() {
		s.indexByKey = make(map[string]string)
		s.indexByFunction = make(map[string][]string)
		s.indexPathByKey = make(map[string]string)
		s.indexPathByFunction = make(map[string]string)

		if err := s.ForEach(func(key string, value []byte) error {
			var mapping EventSourceMapping
			if err := json.Unmarshal(value, &mapping); err != nil {
				return nil
			}
			s.indexMapping(key, &mapping)
			return nil
		}); err != nil {
			logs.Error("failed to build event source index", logs.String("error", err.Error()))
		}
	})
}

func (s *EventSourceStore) indexMapping(key string, mapping *EventSourceMapping) {
	indexPath := mapping.EventSourceArn + "|" + mapping.FunctionArn
	s.indexByKey[key] = indexPath
	s.indexPathByKey[indexPath] = key

	funcList := s.indexByFunction[mapping.FunctionArn]
	s.indexByFunction[mapping.FunctionArn] = append(funcList, key)
	s.indexPathByFunction[mapping.FunctionArn+"|"+key] = key
}

func (s *EventSourceStore) unindexMapping(key string, mapping *EventSourceMapping) {
	indexPath := mapping.EventSourceArn + "|" + mapping.FunctionArn
	delete(s.indexByKey, key)
	delete(s.indexPathByKey, indexPath)

	funcList := s.indexByFunction[mapping.FunctionArn]
	for i, k := range funcList {
		if k == key {
			s.indexByFunction[mapping.FunctionArn] = append(funcList[:i], funcList[i+1:]...)
			break
		}
	}
	delete(s.indexPathByFunction, mapping.FunctionArn+"|"+key)
}

// Create creates a new event source mapping.
func (s *EventSourceStore) Create(mapping *EventSourceMapping) (*EventSourceMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureIndex()

	indexPath := mapping.EventSourceArn + "|" + mapping.FunctionArn
	if _, exists := s.indexPathByKey[indexPath]; exists {
		return nil, ErrEventSourceAlreadyExists
	}

	if mapping.UUID == "" {
		mapping.UUID = uuid.New().String()
	}

	mapping.LastModified = time.Now().UTC()
	mapping.State = "Enabled"
	mapping.StateTransitionReason = "User action"

	if err := s.Put(mapping.UUID, mapping); err != nil {
		return nil, err
	}

	s.indexMapping(mapping.UUID, mapping)

	return mapping, nil
}

// FindByEventSourceAndFunction finds an event source mapping by event source and function ARN.
func (s *EventSourceStore) FindByEventSourceAndFunction(eventSourceArn, functionArn string) (*EventSourceMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureIndex()

	indexPath := eventSourceArn + "|" + functionArn
	key, exists := s.indexPathByKey[indexPath]
	if !exists {
		return nil, ErrEventSourceNotFound
	}

	return s.Get(key)
}

// Get retrieves an event source mapping by UUID.
func (s *EventSourceStore) Get(uuid string) (*EventSourceMapping, error) {
	var mapping EventSourceMapping
	if err := s.BaseStore.Get(uuid, &mapping); err != nil {
		return nil, ErrEventSourceNotFound
	}
	return &mapping, nil
}

// Update updates an event source mapping.
func (s *EventSourceStore) Update(mapping *EventSourceMapping) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInternal(mapping)
}

func (s *EventSourceStore) updateInternal(mapping *EventSourceMapping) error {
	if !s.Exists(mapping.UUID) {
		return ErrEventSourceNotFound
	}

	oldMapping, err := s.Get(mapping.UUID)
	if err != nil {
		return err
	}

	mapping.LastModified = time.Now().UTC()
	if err := s.Put(mapping.UUID, mapping); err != nil {
		return err
	}

	s.ensureIndex()
	s.unindexMapping(mapping.UUID, oldMapping)
	s.indexMapping(mapping.UUID, mapping)

	return nil
}

// Delete removes an event source mapping.
func (s *EventSourceStore) Delete(uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	mapping, err := s.Get(uuid)
	if err != nil {
		return ErrEventSourceNotFound
	}

	if err := s.BaseStore.Delete(uuid); err != nil {
		return err
	}

	s.ensureIndex()
	s.unindexMapping(uuid, mapping)

	return nil
}

// List returns event source mappings with pagination.
func (s *EventSourceStore) List(opts common.ListOptions) (*common.ListResult[EventSourceMapping], error) {
	return common.List[EventSourceMapping](s.BaseStore, opts, nil)
}

// ListByFunction returns event source mappings for a specific function.
func (s *EventSourceStore) ListByFunction(functionArn string) ([]*EventSourceMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureIndex()

	keys, exists := s.indexByFunction[functionArn]
	if !exists {
		return []*EventSourceMapping{}, nil
	}

	mappings := make([]*EventSourceMapping, 0, len(keys))
	for _, key := range keys {
		mapping, err := s.Get(key)
		if err != nil {
			continue
		}
		mappings = append(mappings, mapping)
	}
	return mappings, nil
}

// SetState sets the state of an event source mapping.
func (s *EventSourceStore) SetState(uuid, state, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var mapping EventSourceMapping
	if err := s.BaseStore.Get(uuid, &mapping); err != nil {
		return ErrEventSourceNotFound
	}
	mapping.State = state
	mapping.StateTransitionReason = reason
	return s.updateInternal(&mapping)
}

// Activate enables an event source mapping.
func (s *EventSourceStore) Activate(uuid string) error {
	return s.SetState(uuid, "Enabled", "User action")
}

// Deactivate disables an event source mapping.
func (s *EventSourceStore) Deactivate(uuid string) error {
	return s.SetState(uuid, "Disabled", "User action")
}
