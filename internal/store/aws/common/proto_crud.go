package common

import (
	"google.golang.org/protobuf/proto"
)

// ProtoStoreConfig configures a generic protobuf CRUD store for a single entity
// type. Each method on ProtoStore[D] delegates to the underlying BaseStore
// using the provided conversion functions and error sentinels.
type ProtoStoreConfig[D any] struct {
	Store        *BaseStore
	NewProto     func() proto.Message
	ToDomain     func(proto.Message) *D
	ToProto      func(*D) proto.Message
	IDFunc       func(*D) string
	NotFoundErr  error
	AlreadyExist error
}

// ProtoStore provides generic CRUD operations backed by a BaseStore and
// protobuf serialisation. Use ProtoStoreConfig to create an instance.
type ProtoStore[D any] struct {
	store        *BaseStore
	newProto     func() proto.Message
	toDomain     func(proto.Message) *D
	toProto      func(*D) proto.Message
	idFunc       func(*D) string
	notFoundErr  error
	alreadyExist error
}

// NewProtoStore creates a ProtoStore from the given configuration.
func NewProtoStore[D any](cfg ProtoStoreConfig[D]) *ProtoStore[D] {
	return &ProtoStore[D]{
		store:        cfg.Store,
		newProto:     cfg.NewProto,
		toDomain:     cfg.ToDomain,
		toProto:      cfg.ToProto,
		idFunc:       cfg.IDFunc,
		notFoundErr:  cfg.NotFoundErr,
		alreadyExist: cfg.AlreadyExist,
	}
}

// Create persists a new entity. Returns AlreadyExist error if the key already
// exists (when configured). The caller must hold any required lock.
func (s *ProtoStore[D]) Create(entity *D) error {
	id := s.idFunc(entity)
	if s.alreadyExist != nil && s.store.Exists(id) {
		return s.alreadyExist
	}
	return s.store.PutProto(id, s.toProto(entity))
}

// Get retrieves an entity by key. Returns NotFoundErr if not found (when
// configured), otherwise returns (nil, nil).
func (s *ProtoStore[D]) Get(id string) (*D, error) {
	pb := s.newProto()
	if err := s.store.GetProto(id, pb); err != nil {
		if IsNotFound(err) {
			if s.notFoundErr != nil {
				return nil, s.notFoundErr
			}
			return nil, nil
		}
		return nil, err
	}
	return s.toDomain(pb), nil
}

// Update replaces the persisted state of an existing entity. Returns
// NotFoundErr if the entity does not exist (when configured). The caller must
// hold any required lock.
func (s *ProtoStore[D]) Update(entity *D) error {
	id := s.idFunc(entity)
	if s.notFoundErr != nil && !s.store.Exists(id) {
		return s.notFoundErr
	}
	return s.store.PutProto(id, s.toProto(entity))
}

// UpdateDirect persists an entity without an existence check. The caller must
// hold any required lock.
func (s *ProtoStore[D]) UpdateDirect(entity *D) error {
	return s.store.PutProto(s.idFunc(entity), s.toProto(entity))
}

// Delete removes an entity by key. The caller must hold any required lock.
func (s *ProtoStore[D]) Delete(id string) error {
	return s.store.Delete(id)
}

// DeleteIfExists removes an entity by key if it exists. Returns NotFoundErr if
// not found (when configured). The caller must hold any required lock.
func (s *ProtoStore[D]) DeleteIfExists(id string) error {
	if s.notFoundErr != nil && !s.store.Exists(id) {
		return s.notFoundErr
	}
	return s.store.Delete(id)
}

// List returns all entities in the store.
func (s *ProtoStore[D]) List() ([]*D, error) {
	var items []*D
	err := s.store.ForEach(func(key string, value []byte) error {
		pb := s.newProto()
		if err := proto.Unmarshal(value, pb); err != nil {
			return err
		}
		items = append(items, s.toDomain(pb))
		return nil
	})
	return items, err
}

// ListFiltered returns all entities that pass the filter function.
func (s *ProtoStore[D]) ListFiltered(filter func(*D) bool) ([]*D, error) {
	var items []*D
	err := s.store.ForEach(func(key string, value []byte) error {
		pb := s.newProto()
		if err := proto.Unmarshal(value, pb); err != nil {
			return err
		}
		d := s.toDomain(pb)
		if filter == nil || filter(d) {
			items = append(items, d)
		}
		return nil
	})
	return items, err
}

// ListProto returns all raw protobuf messages in the store, without domain
// conversion.
func (s *ProtoStore[D]) ListProto() ([]proto.Message, error) {
	var items []proto.Message
	err := s.store.ForEach(func(key string, value []byte) error {
		pb := s.newProto()
		if err := proto.Unmarshal(value, pb); err != nil {
			return err
		}
		items = append(items, pb)
		return nil
	})
	return items, err
}

// Exists checks if an entity exists by key.
func (s *ProtoStore[D]) Exists(id string) bool {
	return s.store.Exists(id)
}

// Store returns the underlying BaseStore for advanced operations (e.g.
// ScanPrefix, DeleteByPrefix).
func (s *ProtoStore[D]) Store() *BaseStore {
	return s.store
}

// RawProtoStore provides CRUD operations for entities that are stored as
// protobuf messages directly, without a separate domain type.
type RawProtoStore[P proto.Message] struct {
	store       *BaseStore
	newProto    func() P
	idFunc      func(P) string
	notFoundErr error
}

// RawProtoStoreConfig configures a RawProtoStore.
type RawProtoStoreConfig[P proto.Message] struct {
	Store       *BaseStore
	NewProto    func() P
	IDFunc      func(P) string
	NotFoundErr error
}

// NewRawProtoStore creates a RawProtoStore from the given configuration.
func NewRawProtoStore[P proto.Message](cfg RawProtoStoreConfig[P]) *RawProtoStore[P] {
	return &RawProtoStore[P]{
		store:       cfg.Store,
		newProto:    cfg.NewProto,
		idFunc:      cfg.IDFunc,
		notFoundErr: cfg.NotFoundErr,
	}
}

// Create persists a new protobuf entity. The caller must hold any required lock.
func (s *RawProtoStore[P]) Create(entity P) error {
	return s.store.PutProto(s.idFunc(entity), entity)
}

// Get retrieves a protobuf entity by key. Returns NotFoundErr if configured
// and not found, otherwise (nil, nil).
func (s *RawProtoStore[P]) Get(id string) (P, error) {
	var zero P
	pb := s.newProto()
	if err := s.store.GetProto(id, pb); err != nil {
		if IsNotFound(err) {
			if s.notFoundErr != nil {
				return zero, s.notFoundErr
			}
			return zero, nil
		}
		return zero, err
	}
	return pb, nil
}

// Update persists an updated protobuf entity. The caller must hold any required lock.
func (s *RawProtoStore[P]) Update(entity P) error {
	return s.store.PutProto(s.idFunc(entity), entity)
}

// Delete removes a protobuf entity by key. The caller must hold any required lock.
func (s *RawProtoStore[P]) Delete(id string) error {
	return s.store.Delete(id)
}

// DeleteIfExists removes a protobuf entity if it exists. The caller must hold
// any required lock.
func (s *RawProtoStore[P]) DeleteIfExists(id string) error {
	if s.notFoundErr != nil && !s.store.Exists(id) {
		return s.notFoundErr
	}
	return s.store.Delete(id)
}

// List returns all protobuf entities in the store.
func (s *RawProtoStore[P]) List() ([]P, error) {
	var items []P
	err := s.store.ForEach(func(key string, value []byte) error {
		pb := s.newProto()
		if err := proto.Unmarshal(value, pb); err != nil {
			return err
		}
		items = append(items, pb)
		return nil
	})
	return items, err
}

// Exists checks if an entity exists by key.
func (s *RawProtoStore[P]) Exists(id string) bool {
	return s.store.Exists(id)
}

// Store returns the underlying BaseStore for advanced operations.
func (s *RawProtoStore[P]) Store() *BaseStore {
	return s.store
}
