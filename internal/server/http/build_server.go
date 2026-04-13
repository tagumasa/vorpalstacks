// Package http provides HTTP server functionality for vorpalstacks.
package http

import (
	"context"
	"net/http"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/http/chain"
	"vorpalstacks/internal/server/http/classifier"
	"vorpalstacks/internal/server/listener"
	"vorpalstacks/internal/store/api"
	iamstore "vorpalstacks/internal/store/aws/iam"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// Config holds configuration for the HTTP server.
type Config struct {
	Port            string
	DataPath        string
	AccountID       string
	Region          string
	SignatureConfig SignatureConfig
	UseChainGateway bool
	TLSConfig       TLSConfig
}

// SignatureConfig holds AWS signature verification configuration.
type SignatureConfig struct {
	Enabled         bool
	Service         string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// TLSConfig holds TLS configuration.
type TLSConfig struct {
	Enabled  bool
	Port     string
	CertPath string
	KeyPath  string
	Hostname string
}

// ShutdownHook is a function called during graceful shutdown.
type ShutdownHook func(ctx context.Context)

// Server represents the HTTP server that handles AWS API requests.
type Server struct {
	config            *Config
	storageManager    *storage.RegionStorageManager
	eventBus          *eventbus.EventBus
	httpServer        *http.Server
	httpServerMu      sync.Mutex
	dispatcher        *dispatcher.Dispatcher
	s3Handler         http.Handler
	apiGatewayRuntime http.Handler
	jwksHandler       http.Handler
	handlerMu         sync.RWMutex
	classifier        *classifier.Classifier
	serviceStore      *api.ServiceStore
	operationStore    *api.OperationStore
	chainGateway      *chain.Gateway
	shutdownCtx       context.Context
	shutdownCancel    context.CancelFunc
	closeOnce         sync.Once
	startOnce         sync.Once
	tlsManager        *TLSManager
	tlsServer         *http.Server
	tlsServerMu       sync.Mutex
	iamStore          iamstore.IAMStoreInterface
	s3Store           s3store.S3StoreInterface
	blobStore         storage.BlobStore
	shutdownHooks     []ShutdownHook
	shutdownHooksMu   sync.Mutex
	listenerManager   *listener.Manager
}

// Classifier returns the request classifier.
//
// Returns:
//   - *classifier.Classifier: The request classifier
func (s *Server) Classifier() *classifier.Classifier {
	return s.classifier
}

// ServiceStore returns the service store.
//
// Returns:
//   - *api.ServiceStore: The service store
func (s *Server) ServiceStore() *api.ServiceStore {
	return s.serviceStore
}

// OperationStore returns the operation store.
//
// Returns:
//   - *api.OperationStore: The operation store
func (s *Server) OperationStore() *api.OperationStore {
	return s.operationStore
}

// ChainGateway returns the chain gateway.
//
// Returns:
//   - *chain.Gateway: The chain gateway
func (s *Server) ChainGateway() *chain.Gateway {
	return s.chainGateway
}

// RegisterHandler registers a handler for an operation.
//
// Parameters:
//   - operationName: The name of the operation
//   - handler: The handler function to register
func (s *Server) RegisterHandler(operationName string, handler dispatcher.Handler) {
	s.dispatcher.RegisterHandler(operationName, handler)
}

// Dispatcher returns the dispatcher.
//
// Returns:
//   - *dispatcher.Dispatcher: The dispatcher
func (s *Server) Dispatcher() *dispatcher.Dispatcher {
	return s.dispatcher
}

// Storage returns the global storage.
//
// Returns:
//   - storage.BasicStorage: The global storage
func (s *Server) Storage() storage.BasicStorage {
	globalStore, _ := s.storageManager.GetGlobalStorage()
	return globalStore
}

// StorageManager returns the storage manager.
//
// Returns:
//   - *storage.RegionStorageManager: The storage manager
func (s *Server) StorageManager() *storage.RegionStorageManager {
	return s.storageManager
}

// IAMStore returns the IAM store interface.
func (s *Server) IAMStore() iamstore.IAMStoreInterface {
	return s.iamStore
}

// S3Store returns the S3 store interface.
func (s *Server) S3Store() s3store.S3StoreInterface {
	return s.s3Store
}

// BlobStore returns the blob store interface.
func (s *Server) BlobStore() storage.BlobStore {
	return s.blobStore
}

// EventBus returns the event bus.
func (s *Server) EventBus() *eventbus.EventBus {
	return s.eventBus
}

// SetEventBus sets the event bus.
func (s *Server) SetEventBus(bus *eventbus.EventBus) {
	s.eventBus = bus
}

// RegisterS3Handler registers an S3 handler.
//
// Parameters:
//   - handler: The S3 handler to register
func (s *Server) RegisterS3Handler(handler http.Handler) {
	s.handlerMu.Lock()
	s.s3Handler = handler
	s.handlerMu.Unlock()
}

// S3Handler returns the S3 handler.
//
// Returns:
//   - http.Handler: The S3 handler
func (s *Server) S3Handler() http.Handler {
	s.handlerMu.RLock()
	defer s.handlerMu.RUnlock()
	return s.s3Handler
}

// RegisterAPIGatewayRuntimeHandler registers an API Gateway runtime handler.
//
// Parameters:
//   - handler: The API Gateway runtime handler to register
func (s *Server) RegisterAPIGatewayRuntimeHandler(handler http.Handler) {
	s.handlerMu.Lock()
	s.apiGatewayRuntime = handler
	s.handlerMu.Unlock()
}

// APIGatewayRuntimeHandler returns the API Gateway runtime handler.
//
// Returns:
//   - http.Handler: The API Gateway runtime handler
func (s *Server) APIGatewayRuntimeHandler() http.Handler {
	s.handlerMu.RLock()
	defer s.handlerMu.RUnlock()
	return s.apiGatewayRuntime
}

// RegisterJWKSHandler registers a JWKS handler.
//
// Parameters:
//   - handler: The JWKS handler to register
func (s *Server) RegisterJWKSHandler(handler http.Handler) {
	s.handlerMu.Lock()
	s.jwksHandler = handler
	s.handlerMu.Unlock()
}

// JWKSHandler returns the JWKS handler.
//
// Returns:
//   - http.Handler: The JWKS handler
func (s *Server) JWKSHandler() http.Handler {
	s.handlerMu.RLock()
	defer s.handlerMu.RUnlock()
	return s.jwksHandler
}

// RegisterShutdownHook adds a function to be called during graceful shutdown.
// Hooks are called in reverse registration order (LIFO).
func (s *Server) RegisterShutdownHook(hook ShutdownHook) {
	s.shutdownHooksMu.Lock()
	defer s.shutdownHooksMu.Unlock()
	s.shutdownHooks = append(s.shutdownHooks, hook)
}

// SetListenerManager sets the secondary listener manager.
func (s *Server) SetListenerManager(m *listener.Manager) {
	s.listenerManager = m
}

// ListenerManager returns the secondary listener manager.
func (s *Server) ListenerManager() *listener.Manager {
	return s.listenerManager
}
