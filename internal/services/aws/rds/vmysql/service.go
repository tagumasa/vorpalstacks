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

	// done is closed by the goroutine that runs srv.Start() when the
	// server has stopped (either because Close() called srv.Close() or
	// because Start returned an error on its own). Callers waiting on
	// done can be certain that srv will not produce further events.
	done chan struct{}
}

// Service manages per-instance MySQL engine lifecycles. Each DB instance
// created via the RDS API gets its own isolated engine and TCP listener
// with a dynamically allocated port from the shared range [50200, 50400].
type Service struct {
	mu             sync.Mutex
	storageManager *storage.RegionStorageManager
	portAlloc      *portalloc.Allocator
	instances      map[string]*instanceEntry

	// regionForSnapshots records the region passed to the most recent
	// Open() call. Snapshot / Restore calls use it to address the same
	// regional storage as the source instance. Per-region isolation is
	// already enforced by StorageManager; this field just remembers
	// which region the snapshots should live in.
	regionForSnapshots string
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

// Open creates and starts an isolated MySQL engine for the given DB instance.
// Returns the dynamically allocated port number.
func (s *Service) Open(region, instanceID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storageManager == nil {
		return 0, fmt.Errorf("vmysql: StorageManager not set")
	}

	// Remember the region for snapshot/restore operations. Only the
	// latest region is recorded; multi-region snapshot-from-A-restore-
	// to-B is out of scope for the local engine.
	s.regionForSnapshots = region

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

	// srv.Start() blocks on Listener.Accept(); it only returns when the
	// listener is closed. We need to (a) confirm the server is actually
	// accepting connections before Open() returns, and (b) detect if
	// Start() ever returns so Close() can wait for full shutdown.
	done := make(chan struct{})
	go func() {
		defer close(done)
		err := srv.Start()
		if err != nil {
			logs.Warn("vmysql: server stopped with error",
				logs.Err(err), logs.String("instance", instanceID))
		}
	}()

	entry := &instanceEntry{
		store:    store,
		provider: provider,
		engine:   engine,
		srv:      srv,
		port:     port,
		done:     done,
	}
	s.instances[instanceID] = entry

	// Probe the listener with a real TCP connection so we know the
	// dolt server's Accept loop is servicing connections before
	// declaring the instance "available". Without this probe, Open()
	// can return success while the goroutine has not yet entered
	// Accept — a client connecting in that window would block until
	// Accept runs.
	probeAddr := fmt.Sprintf("127.0.0.1:%d", port)
	probeDeadline := time.Now().Add(3 * time.Second)
	for {
		c, derr := net.DialTimeout("tcp", probeAddr, 500*time.Millisecond)
		if derr == nil {
			c.Close()
			break
		}
		// If srv.Start() already returned (server failed), bail out
		// immediately rather than spinning until the deadline.
		select {
		case <-done:
			s.mu.Lock()
			delete(s.instances, instanceID)
			s.mu.Unlock()
			store.Close()
			s.portAlloc.Release("vmysql", instanceID)
			return 0, fmt.Errorf("vmysql: server failed to start for instance %s", instanceID)
		default:
		}
		if time.Now().After(probeDeadline) {
			s.mu.Lock()
			delete(s.instances, instanceID)
			s.mu.Unlock()
			// Trigger shutdown of the goroutine. srv.Close() will
			// cause srv.Start() to return; we then wait on done.
			srv.Close()
			<-done
			store.Close()
			s.portAlloc.Release("vmysql", instanceID)
			return 0, fmt.Errorf("vmysql: server did not become ready for instance %s within 3s", instanceID)
		}
		time.Sleep(50 * time.Millisecond)
	}

	logs.Info("vmysql: instance engine started", logs.String("instance", instanceID), logs.Int("port", port))
	return port, nil
}

// Close stops the MySQL engine for the given instance and releases its port.
// It blocks until the dolt server's accept loop has actually exited, so that
// callers can be certain no in-flight connection handlers will touch the
// store after Close returns. Without the wait, store.Close() could mark the
// store closed while a handler was still mid-query.
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
		// dolt's Server.Close calls Listener.Close which in turn calls
		// eg.Wait() on the connection-handler errgroup, so by the time
		// srv.Close returns no handler goroutine is still running.
		if err := entry.srv.Close(); err != nil {
			logs.Warn("vmysql: srv.Close returned error",
				logs.Err(err), logs.String("instance", instanceID))
		}
	}
	// Wait for the goroutine that ran srv.Start() to observe the
	// listener closure and exit. This closes the race between
	// srv.Close() returning and store.Close() running.
	if entry.done != nil {
		<-entry.done
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

// Shutdown stops all running instance engines. It waits for every
// engine's accept loop to exit before returning so that no handler
// goroutine touches a closed store after Shutdown returns.
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
			if err := entry.srv.Close(); err != nil {
				logs.Warn("vmysql: Shutdown srv.Close returned error",
					logs.Err(err), logs.String("instance", id))
			}
		}
		if entry.done != nil {
			<-entry.done
		}
		entry.store.Close()
		if err := s.portAlloc.Release("vmysql", id); err != nil {
			logs.Warn("vmysql: Shutdown port release failed",
				logs.Err(err), logs.String("instance", id))
		}
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

// SnapshotData copies every row, schema, and index from the source
// instance's rdbengine store into a snapshot store whose key prefix is
// 'snap_<snapshotID>'. Called by CreateDBSnapshot (manual snapshots)
// and DeleteDBInstance (final snapshots) so that
// RestoreDBInstanceFromDBSnapshot can recover user tables.
//
// The source is read through a Pebble point-in-time snapshot so that
// all tables observe a consistent state regardless of concurrent writes.
//
// Idempotent: any pre-existing snapshot data for this ID is deleted
// before the copy begins, so re-runs after a partial failure succeed.
//
// The copy is non-transactional across tables: each table is copied in
// its own Pebble batch. A crash mid-snapshot leaves a partial snapshot
// that will fail cleanly at restore time rather than silently producing
// a database with half the tables.
func (s *Service) SnapshotData(instanceID, snapshotID string) (err error) {
	if snapshotID == "" {
		return fmt.Errorf("vmysql: SnapshotData: snapshotID is empty")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return fmt.Errorf("vmysql: SnapshotData: %w", err)
	}
	src := s.GetInstanceStore(instanceID)
	if src == nil {
		return fmt.Errorf("vmysql: SnapshotData: instance %s is not open", instanceID)
	}

	// Create a point-in-time snapshot of the source store so that all
	// table scans observe a consistent view. Without this, writes that
	// commit between scanning table A and table B would produce a
	// cross-table inconsistent snapshot. The snapshot reader is closed
	// before return to release the underlying Pebble snapshot handle.
	snapSrc, err := src.NewSnapshotReader()
	if err != nil {
		return fmt.Errorf("vmysql: SnapshotData: create snapshot reader: %w", err)
	}
	defer snapSrc.CloseSnapshot()

	// Delete any pre-existing snapshot data for this ID so the snapshot
	// is idempotent. Without this, a re-run after a partial failure
	// would collide on ErrAlreadyExists at the first duplicate row.
	_ = s.DeleteSnapshotData(snapshotID)

	// If the copy fails partway through, clean up the partial snapshot
	// data so the snap_<snapshotID> bucket does not leak.
	defer func() {
		if err != nil {
			if cleanupErr := s.DeleteSnapshotData(snapshotID); cleanupErr != nil {
				logs.Warn("vmysql: SnapshotData cleanup failed after error",
					logs.String("snapshot", snapshotID), logs.Err(cleanupErr))
			}
		}
	}()

	snapStore, err := s.openSnapshotStore(snapshotID)
	if err != nil {
		return err
	}
	defer snapStore.Close()

	ctx := context.Background()
	dbs, err := snapSrc.ListDatabases(ctx)

	// An instance with no databases yet is a valid snapshot target —
	// the snapshot is simply empty. Don't treat ListDatabases nil as an
	// error in that case.
	if err != nil {
		return fmt.Errorf("vmysql: SnapshotData: list databases: %w", err)
	}

	for _, dbName := range dbs {
		if err := snapStore.CreateDatabase(ctx, dbName); err != nil && err != rdbengine.ErrAlreadyExists {
			return fmt.Errorf("vmysql: SnapshotData: create database %s: %w", dbName, err)
		}
		tables, err := snapSrc.ListTables(ctx, dbName)
		if err != nil {
			return fmt.Errorf("vmysql: SnapshotData: list tables in %s: %w", dbName, err)
		}
		for _, tbl := range tables {
			schema, err := snapSrc.GetTableSchema(ctx, dbName, tbl)
			if err != nil {
				return fmt.Errorf("vmysql: SnapshotData: schema for %s.%s: %w", dbName, tbl, err)
			}
			if err := snapStore.CreateTable(ctx, dbName, schema); err != nil && err != rdbengine.ErrAlreadyExists {
				return fmt.Errorf("vmysql: SnapshotData: create table %s.%s: %w", dbName, tbl, err)
			}
			iter, err := snapSrc.ScanRows(ctx, dbName, tbl, rdbengine.ScanOptions{})
			if err != nil {
				return fmt.Errorf("vmysql: SnapshotData: scan %s.%s: %w", dbName, tbl, err)
			}
			for iter.Next() {
				row := iter.Row()
				engineRow := make(rdbengine.Row, len(row))
				for k, v := range row {
					engineRow[k] = v
				}
				pk, err := rdbengine.EncodePK(schema, row)
				if err != nil {
					iter.Close()
					return fmt.Errorf("vmysql: SnapshotData: encode pk for %s.%s: %w", dbName, tbl, err)
				}
				if err := snapStore.InsertRow(ctx, dbName, tbl, pk, engineRow); err != nil {
					iter.Close()
					return fmt.Errorf("vmysql: SnapshotData: insert row into %s.%s: %w", dbName, tbl, err)
				}
			}
			if iterErr := iter.Error(); iterErr != nil {
				iter.Close()
				return fmt.Errorf("vmysql: SnapshotData: iterate %s.%s: %w", dbName, tbl, iterErr)
			}
			iter.Close()

			indexes, err := snapSrc.ListIndexes(ctx, dbName, tbl)
			if err != nil {
				return fmt.Errorf("vmysql: SnapshotData: list indexes for %s.%s: %w", dbName, tbl, err)
			}
			for _, idx := range indexes {
				if err := snapStore.CreateIndex(ctx, dbName, tbl, idx.Name, idx.Columns, idx.Unique); err != nil && err != rdbengine.ErrAlreadyExists {
					return fmt.Errorf("vmysql: SnapshotData: create index %s on %s.%s: %w", idx.Name, dbName, tbl, err)
				}
				// When the index already exists, verify the definition
				// matches so a stale index with different columns or
				// uniqueness does not silently diverge from the source.
				if err == rdbengine.ErrAlreadyExists {
					existingIndexes, lookupErr := snapStore.ListIndexes(ctx, dbName, tbl)
					if lookupErr != nil {
						return fmt.Errorf("vmysql: SnapshotData: verify indexes on %s.%s: %w", dbName, tbl, lookupErr)
					}
					for _, ei := range existingIndexes {
						if ei.Name == idx.Name {
							if ei.Unique != idx.Unique || !equalStringSlices(ei.Columns, idx.Columns) {
								return fmt.Errorf("vmysql: SnapshotData: index %s on %s.%s already exists with a different definition", idx.Name, dbName, tbl)
							}
							break
						}
					}
				}
			}
		}
	}
	return nil
}

// RestoreData copies every row, schema, and index from a snapshot
// store (key prefix 'snap_<snapshotID>') into the destination
// instance's rdbengine store (key prefix 'mysql_<instanceID>').
// Used by RestoreDBInstanceFromDBSnapshot to recover user data.
//
// The destination must already exist and have its engine open. Rows
// from the snapshot are inserted into the destination; existing rows
// with the same primary key are left untouched (idempotent restore).
func (s *Service) RestoreData(snapshotID, instanceID string) error {
	if snapshotID == "" {
		return fmt.Errorf("vmysql: RestoreData: snapshotID is empty")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return fmt.Errorf("vmysql: RestoreData: %w", err)
	}
	dst := s.GetInstanceStore(instanceID)
	if dst == nil {
		return fmt.Errorf("vmysql: RestoreData: instance %s is not open", instanceID)
	}

	snapStore, err := s.openSnapshotStore(snapshotID)
	if err != nil {
		return err
	}
	defer snapStore.Close()

	ctx := context.Background()
	dbs, err := snapStore.ListDatabases(ctx)
	if err != nil {
		return fmt.Errorf("vmysql: RestoreData: list snapshot databases: %w", err)
	}
	for _, dbName := range dbs {
		if err := dst.CreateDatabase(ctx, dbName); err != nil && err != rdbengine.ErrAlreadyExists {
			return fmt.Errorf("vmysql: RestoreData: create database %s: %w", dbName, err)
		}
		tables, err := snapStore.ListTables(ctx, dbName)
		if err != nil {
			return fmt.Errorf("vmysql: RestoreData: list tables in %s: %w", dbName, err)
		}
		for _, tbl := range tables {
			schema, err := snapStore.GetTableSchema(ctx, dbName, tbl)
			if err != nil {
				return fmt.Errorf("vmysql: RestoreData: schema for %s.%s: %w", dbName, tbl, err)
			}
			if err := dst.CreateTable(ctx, dbName, schema); err != nil && err != rdbengine.ErrAlreadyExists {
				return fmt.Errorf("vmysql: RestoreData: create table %s.%s: %w", dbName, tbl, err)
			}
			iter, err := snapStore.ScanRows(ctx, dbName, tbl, rdbengine.ScanOptions{})
			if err != nil {
				return fmt.Errorf("vmysql: RestoreData: scan %s.%s: %w", dbName, tbl, err)
			}
			for iter.Next() {
				row := iter.Row()
				engineRow := make(rdbengine.Row, len(row))
				for k, v := range row {
					engineRow[k] = v
				}
				pk, err := rdbengine.EncodePK(schema, row)
				if err != nil {
					iter.Close()
					return fmt.Errorf("vmysql: RestoreData: encode pk for %s.%s: %w", dbName, tbl, err)
				}
				// InsertRow returns ErrAlreadyExists for duplicate PKs;
				// we treat that as success so a re-restore over an
				// existing instance is idempotent.
				if err := dst.InsertRow(ctx, dbName, tbl, pk, engineRow); err != nil && err != rdbengine.ErrAlreadyExists {
					iter.Close()
					return fmt.Errorf("vmysql: RestoreData: insert row into %s.%s: %w", dbName, tbl, err)
				}
			}
			if iterErr := iter.Error(); iterErr != nil {
				iter.Close()
				return fmt.Errorf("vmysql: RestoreData: iterate %s.%s: %w", dbName, tbl, iterErr)
			}
			iter.Close()
		}
	}
	return nil
}

// DeleteSnapshotData removes the snapshot store's data for the given
// snapshotID. Called when a DBInstanceSnapshot is deleted so the
// snapshot's row data does not leak indefinitely.
func (s *Service) DeleteSnapshotData(snapshotID string) error {
	if snapshotID == "" {
		return fmt.Errorf("vmysql: DeleteSnapshotData: snapshotID is empty")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return fmt.Errorf("vmysql: DeleteSnapshotData: %w", err)
	}
	if s.storageManager == nil {
		return fmt.Errorf("vmysql: DeleteSnapshotData: StorageManager not set")
	}
	snapStore, err := s.openSnapshotStore(snapshotID)
	if err != nil {
		return err
	}
	// DropDatabase deletes all tables, indexes, rows atomically per db.
	// Iterate all databases and drop each.
	ctx := context.Background()
	dbs, err := snapStore.ListDatabases(ctx)
	if err != nil {
		snapStore.Close()
		return err
	}
	for _, db := range dbs {
		if err := snapStore.DropDatabase(ctx, db); err != nil {
			snapStore.Close()
			return err
		}
	}
	snapStore.Close()
	return nil
}

// openSnapshotStore creates a rdbengine.Store whose key prefix is
// 'snap_<snapshotID>'. Snapshots share the source instance's bucket so
// data never crosses regional storage boundaries.
func (s *Service) openSnapshotStore(snapshotID string) (*rdbengine.Store, error) {
	if s.storageManager == nil {
		return nil, fmt.Errorf("vmysql: StorageManager not set")
	}
	// Snapshots are stored under the region of the Service (set via
	// SetStorageManager). Use the same path the Open call uses.
	rs, err := s.storageManager.GetStorage(s.regionForSnapshots)
	if err != nil {
		return nil, fmt.Errorf("vmysql: snapshot: get regional storage: %w", err)
	}
	bkt := rs.Bucket(bucketName)
	bb, ok := bkt.(storage.BatchBucket)
	if !ok {
		return nil, fmt.Errorf("vmysql: bucket %q does not support batch operations", bucketName)
	}
	return rdbengine.New(bb, rdbengine.Options{Engine: "snap_" + snapshotID})
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

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
