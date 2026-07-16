package vmysql

import (
	"context"
	"fmt"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"

	"vorpalstacks/internal/core/storage/rdbengine"
)

type pebbleTransaction struct {
	readOnly bool
}

var _ sql.Transaction = (*pebbleTransaction)(nil)

func (t *pebbleTransaction) String() string   { return "pebble-vmysql transaction" }
func (t *pebbleTransaction) IsReadOnly() bool { return t.readOnly }

type pebbleSession struct {
	*sql.BaseSession
	provider *pebbleProvider
	inTx     bool
	txnBatch *rdbengine.TxnBatch
	store    *rdbengine.Store
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

func (s *pebbleSession) StartTransaction(ctx *sql.Context, tCharacteristic sql.TransactionCharacteristic) (sql.Transaction, error) {
	s.inTx = true
	if s.store != nil {
		s.txnBatch = s.store.NewTxnBatch()
	}
	return &pebbleTransaction{readOnly: tCharacteristic == sql.ReadOnly}, nil
}

func (s *pebbleSession) CommitTransaction(ctx *sql.Context, tx sql.Transaction) error {
	if s.txnBatch != nil {
		err := s.txnBatch.Commit()
		s.txnBatch = nil
		s.inTx = false
		return err
	}
	s.inTx = false
	return nil
}

func (s *pebbleSession) Rollback(ctx *sql.Context, transaction sql.Transaction) error {
	if s.txnBatch != nil {
		s.txnBatch.Rollback()
		s.txnBatch = nil
	}
	s.inTx = false
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

func (s *pebbleSession) isInTx() bool {
	return s.inTx
}

func (s *pebbleSession) txnInsertRow(db, table string, pk []byte, row rdbengine.Row) error {
	return s.store.TxnInsertRow(s.txnBatch, db, table, pk, row)
}

func (s *pebbleSession) txnUpdateRow(db, table string, pk []byte, row rdbengine.Row) error {
	return s.store.TxnUpdateRow(s.txnBatch, db, table, pk, row)
}

func (s *pebbleSession) txnDeleteRow(db, table string, pk []byte) error {
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
