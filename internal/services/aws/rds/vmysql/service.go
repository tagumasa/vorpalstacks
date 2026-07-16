package vmysql

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/vitess/go/mysql"
	"vorpalstacks/internal/core/logs"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/rdbengine"
	"vorpalstacks/internal/server/portalloc"
	rdssvc "vorpalstacks/internal/services/aws/rds"
)

var (
	_ rdssvc.Engine    = (*Service)(nil)
	_ rdssvc.GetPorter = (*Service)(nil)
)

const bucketName = "rds_data"

type instanceEntry struct {
	store    *rdbengine.Store
	provider *pebbleProvider
	engine   *sqle.Engine
	srv      *server.Server
	port     int
}

// Service manages per-instance MySQL engine lifecycles. Each DB instance
// created via the RDS API gets its own isolated engine and TCP listener
// with a dynamically allocated port from the shared range [50200, 50400].
type Service struct {
	mu             sync.Mutex
	storageManager *storage.RegionStorageManager
	portAlloc      *portalloc.Allocator
	instances      map[string]*instanceEntry
}

func NewService(allocator *portalloc.Allocator) *Service {
	return &Service{
		portAlloc: allocator,
		instances: make(map[string]*instanceEntry),
	}
}

func (s *Service) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

func (s *Service) SetRegion(region string) {}

// Open creates and starts an isolated MySQL engine for the given DB instance.
// Returns the dynamically allocated port number.
func (s *Service) Open(region, instanceID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storageManager == nil {
		return 0, fmt.Errorf("vmysql: StorageManager not set")
	}

	if _, exists := s.instances[instanceID]; exists {
		return 0, fmt.Errorf("vmysql: instance %s already open", instanceID)
	}

	rs, err := s.storageManager.GetStorage(region)
	if err != nil {
		return 0, fmt.Errorf("vmysql: failed to get regional storage: %w", err)
	}
	bkt := rs.Bucket(bucketName)
	bb, ok := bkt.(storage.BatchBucket)
	if !ok {
		return 0, fmt.Errorf("vmysql: bucket %q does not support batch operations", bucketName)
	}

	store, err := rdbengine.New(bb, rdbengine.Options{Engine: "mysql_" + instanceID})
	if err != nil {
		return 0, fmt.Errorf("vmysql: failed to create rdbengine store: %w", err)
	}

	provider := newPebbleProvider(store)
	engine := sqle.NewDefault(provider)

	port, err := s.portAlloc.Get("vmysql", instanceID)
	if err != nil {
		store.Close()
		return 0, fmt.Errorf("vmysql: failed to allocate port: %w", err)
	}

	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		store.Close()
		s.portAlloc.Release("vmysql", instanceID)
		return 0, fmt.Errorf("vmysql: failed to listen on %s: %w", addr, err)
	}

	ctxFactory := func(ctx context.Context, opts ...sql.ContextOption) *sql.Context {
		return sql.NewContext(ctx, opts...)
	}
	sessionBuilder := func(ctx context.Context, conn *mysql.Conn, addr string) (sql.Session, error) {
		return newPebbleSession(sql.NewBaseSession(), provider, store), nil
	}

	cfg := server.Config{
		Protocol: "tcp",
		Address:  addr,
		Listener: ln,
	}

	srv, err := server.NewServer(cfg, engine, ctxFactory, sessionBuilder, &noopListener{})
	if err != nil {
		ln.Close()
		store.Close()
		s.portAlloc.Release("vmysql", instanceID)
		return 0, fmt.Errorf("vmysql: failed to create server: %w", err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			logs.Warn("vmysql: server stopped", logs.Err(err), logs.String("instance", instanceID))
		}
	}()

	s.instances[instanceID] = &instanceEntry{
		store:    store,
		provider: provider,
		engine:   engine,
		srv:      srv,
		port:     port,
	}

	logs.Info("vmysql: instance engine started", logs.String("instance", instanceID), logs.Int("port", port))
	return port, nil
}

// Close stops the MySQL engine for the given instance and releases its port.
func (s *Service) Close(instanceID string) error {
	s.mu.Lock()
	entry, ok := s.instances[instanceID]
	if ok {
		delete(s.instances, instanceID)
	}
	s.mu.Unlock()

	if !ok {
		return nil
	}

	if entry.srv != nil {
		entry.srv.Close()
	}
	entry.store.Close()

	if err := s.portAlloc.Release("vmysql", instanceID); err != nil {
		logs.Warn("vmysql: failed to release port", logs.Err(err), logs.String("instance", instanceID))
	}

	logs.Info("vmysql: instance engine stopped", logs.String("instance", instanceID))
	return nil
}

// GetPort returns the allocated port for an instance.
func (s *Service) GetPort(instanceID string) (int, error) {
	s.mu.Lock()
	_, exists := s.instances[instanceID]
	s.mu.Unlock()
	if !exists {
		return 0, fmt.Errorf("vmysql: instance %s not found", instanceID)
	}
	return s.portAlloc.Get("vmysql", instanceID)
}

func (s *Service) EngineType() string { return "mysql" }

// Shutdown stops all running instance engines.
func (s *Service) Shutdown() {
	s.mu.Lock()
	instances := make(map[string]*instanceEntry, len(s.instances))
	for k, v := range s.instances {
		instances[k] = v
	}
	s.instances = make(map[string]*instanceEntry)
	s.mu.Unlock()

	for id, entry := range instances {
		if entry.srv != nil {
			entry.srv.Close()
		}
		entry.store.Close()
		s.portAlloc.Release("vmysql", id)
	}
}

// GetInstanceStore returns the rdbengine store for a running instance.
func (s *Service) GetInstanceStore(instanceID string) *rdbengine.Store {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.instances[instanceID]; ok {
		return entry.store
	}
	return nil
}

// GetEngine returns the sqle.Engine for a running instance, or nil if not found.
func (s *Service) GetEngine(instanceID string) *sqle.Engine {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.instances[instanceID]; ok {
		return entry.engine
	}
	return nil
}

func (s *Service) NewContext(instanceID string, database string) *sql.Context {
	s.mu.Lock()
	entry, ok := s.instances[instanceID]
	s.mu.Unlock()
	if !ok {
		return nil
	}
	return newEngineContext(entry.engine, entry.provider, entry.store, database)
}

type noopListener struct{}

func (n *noopListener) ClientConnected()                             {}
func (n *noopListener) ClientDisconnected()                          {}
func (n *noopListener) QueryStarted()                                {}
func (n *noopListener) QueryCompleted(success bool, d time.Duration) {}
