package lambda

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// LambdaStore provides access to all Lambda stores.
type LambdaStore struct {
	functions    *FunctionStore
	layers       *LayerStore
	eventSources *EventSourceStore
}

// NewLambdaStore creates a new LambdaStore with the given storage.
func NewLambdaStore(store storage.BasicStorage, accountId, region string) *LambdaStore {
	return &LambdaStore{
		functions:    NewFunctionStore(store, accountId, region),
		layers:       NewLayerStore(store, accountId, region),
		eventSources: NewEventSourceStore(store, accountId, region),
	}
}

// NewLambdaStoreFromStores creates a new LambdaStore with the given stores.
func NewLambdaStoreFromStores(functions *FunctionStore, layers *LayerStore, eventSources *EventSourceStore) *LambdaStore {
	return &LambdaStore{
		functions:    functions,
		layers:       layers,
		eventSources: eventSources,
	}
}

// Functions returns the function store.
func (s *LambdaStore) Functions() FunctionStoreInterface {
	return s.functions
}

// Layers returns the layer store.
func (s *LambdaStore) Layers() LayerStoreInterface {
	return s.layers
}

// EventSources returns the event source store.
func (s *LambdaStore) EventSources() EventSourceStoreInterface {
	return s.eventSources
}

// FunctionsRaw returns the raw function store.
func (s *LambdaStore) FunctionsRaw() *FunctionStore {
	return s.functions
}

// LayersRaw returns the raw layer store.
func (s *LambdaStore) LayersRaw() *LayerStore {
	return s.layers
}

// EventSourcesRaw returns the raw event source store.
func (s *LambdaStore) EventSourcesRaw() *EventSourceStore {
	return s.eventSources
}

var _ LambdaStoreInterface = (*LambdaStore)(nil)

var (
	_ FunctionStoreInterface    = (*FunctionStore)(nil)
	_ EventSourceStoreInterface = (*EventSourceStore)(nil)
	_ LayerStoreInterface       = (*LayerStore)(nil)
)

// WrapFunctionStore wraps a FunctionStore in its interface.
func WrapFunctionStore(s *FunctionStore) FunctionStoreInterface {
	return s
}

// WrapEventSourceStore wraps an EventSourceStore in its interface.
func WrapEventSourceStore(s *EventSourceStore) EventSourceStoreInterface {
	return s
}

// WrapLayerStore wraps a LayerStore in its interface.
func WrapLayerStore(s *LayerStore) LayerStoreInterface {
	return s
}

// LambdaStoreFactory creates Lambda stores.
type LambdaStoreFactory struct {
	storage   storage.BasicStorage
	accountID string
}

// NewLambdaStoreFactory creates a new Lambda store factory.
func NewLambdaStoreFactory(store storage.BasicStorage, accountID string) *LambdaStoreFactory {
	return &LambdaStoreFactory{
		storage:   store,
		accountID: accountID,
	}
}

// CreateStore creates a new Lambda store for the given region.
func (f *LambdaStoreFactory) CreateStore(region string) *LambdaStore {
	return NewLambdaStore(f.storage, f.accountID, region)
}

// CreateFunctionStore creates a new function store for the given region.
func (f *LambdaStoreFactory) CreateFunctionStore(region string) *FunctionStore {
	return NewFunctionStore(f.storage, f.accountID, region)
}

// CreateLayerStore creates a new layer store for the given region.
func (f *LambdaStoreFactory) CreateLayerStore(region string) *LayerStore {
	return NewLayerStore(f.storage, f.accountID, region)
}

// CreateEventSourceStore creates a new event source store for the given region.
func (f *LambdaStoreFactory) CreateEventSourceStore(region string) *EventSourceStore {
	return NewEventSourceStore(f.storage, f.accountID, region)
}

// ListOptions defines pagination options for listing resources.
type ListOptions = common.ListOptions

// ListResult defines the result of a list operation with pagination information.
type ListResult[T any] = common.ListResult[T]
