package vmysql

import (
	"context"
	"fmt"
	"sync"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"

	"vorpalstacks/internal/core/storage/rdbengine"
)

// readOnlyTxn is returned by StartTransaction when the caller requests
// a read-only transaction. It carries no state other than the readOnly
// flag; the enforcement happens in the txn* helpers on pebbleSession,
// which consult currentTxReadOnly before forwarding writes to the
// underlying TxnBatch.
type pebbleTransaction struct {
	readOnly bool
}

var _ sql.Transaction = (*pebbleTransaction)(nil)

func (t *pebbleTransaction) String() string   { return "pebble-vmysql transaction" }
func (t *pebbleTransaction) IsReadOnly() bool { return t.readOnly }

// pebbleSession implements sql.Session / sql.TransactionSession on top
// of an rdbengine.Store. Transaction state (inTx, txnBatch, readOnly)
// is guarded by txMu so that concurrent queries on the same session —
// which go-mysql-server may issue when running parallel branches of a
// plan — cannot race on the per-session TxnBatch.
type pebbleSession struct {
	*sql.BaseSession
	provider *pebbleProvider
	store    *rdbengine.Store

	txMu     sync.Mutex
	inTx     bool
	txnBatch *rdbengine.TxnBatch
	readOnly bool
}

var _ sql.Session = (*pebbleSession)(nil)
var _ sql.TransactionSession = (*pebbleSession)(nil)

func newPebbleSession(base *sql.BaseSession, provider *pebbleProvider, store *rdbengine.Store) *pebbleSession {
	return &pebbleSession{
		BaseSession: base,
		provider:    provider,
		store:       store,
	}
}

func sessionFromCtx(ctx *sql.Context) *pebbleSession {
	s, ok := ctx.Session.(*pebbleSession)
	if !ok {
		return nil
	}
	return s
}

// StartTransaction opens a new transactional TxnBatch on the underlying
// store. Returns an error if the session has no store (which would
// otherwise cause CommitTransaction to silently no-op and lose every
// write the caller believed was committed).
func (s *pebbleSession) StartTransaction(ctx *sql.Context, tCharacteristic sql.TransactionCharacteristic) (sql.Transaction, error) {
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.store == nil {
		return nil, fmt.Errorf("vmysql: cannot start transaction: session store is nil")
	}
	// SQL standard: only one active transaction per session. If a
	// second Start arrives without a Commit/Rollback, surface it as an
	// error rather than silently orphaning the previous batch.
	if s.inTx {
		return nil, fmt.Errorf("vmysql: transaction already in progress on this session")
	}
	s.inTx = true
	s.txnBatch = s.store.NewTxnBatch()
	s.readOnly = tCharacteristic == sql.ReadOnly
	return &pebbleTransaction{readOnly: s.readOnly}, nil
}

func (s *pebbleSession) CommitTransaction(ctx *sql.Context, tx sql.Transaction) error {
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.txnBatch == nil {
		// CommitTransaction without a matching StartTransaction is a
		// client protocol error; surface it rather than silently
		// returning nil which would mask the lost writes.
		return fmt.Errorf("vmysql: CommitTransaction called with no active transaction")
	}
	err := s.txnBatch.Commit()
	s.txnBatch = nil
	s.inTx = false
	s.readOnly = false
	return err
}

// Rollback aborts the active transaction. Unlike CommitTransaction,
// Rollback without an active transaction is a silent no-op rather than
// an error. This matches MySQL semantics where ROLLBACK without a
// transaction produces a warning, not an error, and prevents breakage
// if go-mysql-server calls Rollback defensively.
func (s *pebbleSession) Rollback(ctx *sql.Context, transaction sql.Transaction) error {
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.txnBatch != nil {
		s.txnBatch.Rollback()
		s.txnBatch = nil
	}
	s.inTx = false
	s.readOnly = false
	return nil
}

func (s *pebbleSession) CreateSavepoint(ctx *sql.Context, transaction sql.Transaction, name string) error {
	return fmt.Errorf("savepoints are not supported")
}

func (s *pebbleSession) RollbackToSavepoint(ctx *sql.Context, transaction sql.Transaction, name string) error {
	return fmt.Errorf("savepoints are not supported")
}

func (s *pebbleSession) ReleaseSavepoint(ctx *sql.Context, transaction sql.Transaction, name string) error {
	return fmt.Errorf("savepoints are not supported")
}

// isInTx reports whether the session is currently inside a transaction.
// The method acquires txMu and is goroutine-safe. The return value is
// a point-in-time snapshot: a concurrent CommitTransaction or Rollback
// may change the answer before the caller acts on it. Callers that
// need atomicity should hold txMu for the full check-then-act sequence
// rather than calling isInTx in isolation.
func (s *pebbleSession) isInTx() bool {
	s.txMu.Lock()
	defer s.txMu.Unlock()
	return s.inTx
}

// enforceWritable is called by every transactional write helper. It
// rejects writes issued inside a READ ONLY transaction, matching the
// SQL standard and MySQL behaviour for START TRANSACTION READ ONLY.
func (s *pebbleSession) enforceWritable() error {
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.inTx && s.readOnly {
		return fmt.Errorf("vmysql: cannot execute write in a READ ONLY transaction")
	}
	return nil
}

func (s *pebbleSession) txnInsertRow(db, table string, pk []byte, row rdbengine.Row) error {
	if s.store == nil {
		return fmt.Errorf("vmysql: txn insert: session store is nil")
	}
	if err := s.enforceWritable(); err != nil {
		return err
	}
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.txnBatch == nil {
		return fmt.Errorf("vmysql: txn insert: no active transaction")
	}
	return s.store.TxnInsertRow(s.txnBatch, db, table, pk, row)
}

func (s *pebbleSession) txnUpdateRow(db, table string, pk []byte, row rdbengine.Row) error {
	if s.store == nil {
		return fmt.Errorf("vmysql: txn update: session store is nil")
	}
	if err := s.enforceWritable(); err != nil {
		return err
	}
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.txnBatch == nil {
		return fmt.Errorf("vmysql: txn update: no active transaction")
	}
	return s.store.TxnUpdateRow(s.txnBatch, db, table, pk, row)
}

func (s *pebbleSession) txnDeleteRow(db, table string, pk []byte) error {
	if s.store == nil {
		return fmt.Errorf("vmysql: txn delete: session store is nil")
	}
	if err := s.enforceWritable(); err != nil {
		return err
	}
	s.txMu.Lock()
	defer s.txMu.Unlock()
	if s.txnBatch == nil {
		return fmt.Errorf("vmysql: txn delete: no active transaction")
	}
	return s.store.TxnDeleteRow(s.txnBatch, db, table, pk)
}

func newEngineContext(engine *sqle.Engine, provider *pebbleProvider, store *rdbengine.Store, database string) *sql.Context {
	sess := newPebbleSession(sql.NewBaseSession(), provider, store)
	ctx := sql.NewContext(context.Background(), sql.WithSession(sess))
	if database != "" {
		ctx.SetCurrentDatabase(database)
	}
	return ctx
}
