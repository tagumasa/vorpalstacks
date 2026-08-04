package rdsdata

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
)

// commitTransactionCore claims the transaction entry from the map, waits
// for any in-flight or background statements, then executes COMMIT.
//
// On COMMIT failure the AWS spec requires the transaction to be rolled
// back ("If the COMMIT fails, the transaction is rolled back"). The core
// executes ROLLBACK for cleanup so the engine state and the map state
// remain consistent: the transaction is gone from both the map and the
// engine by the time the error is returned.
//
// The entry is removed from the map BEFORE executing COMMIT so that a
// concurrent caller cannot double-commit. This is safe because COMMIT
// failure triggers an immediate ROLLBACK — the transaction is terminal
// regardless of the outcome.
func (s *RDSDataService) commitTransactionCore(txID string) error {
	entry, err := s.claimTransaction(txID)
	if err != nil {
		return err
	}
	defer entry.execMu.Unlock()

	commitCtx := entry.sqlCtx
	if commitCtx == nil {
		commitCtx = newSQLContext(entry.database)
	}

	if _, err := executeSQL(entry.engine, commitCtx, "COMMIT", false, "", ""); err != nil {
		// AWS spec: "If the COMMIT fails, the transaction is rolled
		// back." Execute ROLLBACK to clean up the engine-side
		// transaction so it does not leak as an orphan.
		_, _ = executeSQL(entry.engine, commitCtx, "ROLLBACK", false, "", "")
		return mapSQLError(err)
	}

	logs.Info("rdsdata: CommitTransaction", logs.String("txId", txID))
	return nil
}

// rollbackTransactionCore claims the transaction entry from the map, waits
// for any in-flight or background statements, then executes ROLLBACK.
//
// The entry is removed from the map BEFORE executing ROLLBACK so that a
// concurrent caller cannot double-rollback. If ROLLBACK fails at the engine
// level the entry stays removed — the engine transaction may leak, but
// re-inserting a potentially corrupt sqlCtx would be worse.
func (s *RDSDataService) rollbackTransactionCore(txID string) error {
	entry, err := s.claimTransaction(txID)
	if err != nil {
		return err
	}
	defer entry.execMu.Unlock()

	rollbackCtx := entry.sqlCtx
	if rollbackCtx == nil {
		rollbackCtx = newSQLContext(entry.database)
	}

	if _, err := executeSQL(entry.engine, rollbackCtx, "ROLLBACK", false, "", ""); err != nil {
		return mapSQLError(err)
	}

	logs.Info("rdsdata: RollbackTransaction", logs.String("txId", txID))
	return nil
}

// claimTransaction atomically removes a transaction from the map and
// returns its entry. It performs expiry checks, bounded waits for
// background statements (ContinueAfterTimeout), and serialisation against
// in-flight ExecuteStatement calls. If any pre-condition fails the entry
// is NOT removed from the map, keeping the transaction retryable.
//
// On success the entry has been removed from the map and execMu has been
// acquired (released by the caller via the returned entry's defer).
// Callers must call entry.execMu.Unlock() when done, or use a defer
// immediately after claimTransaction returns.
func (s *RDSDataService) claimTransaction(txID string) (*staleEntry, error) {
	s.mu.Lock()
	entry, ok := s.transactions[txID]
	s.mu.Unlock()

	if !ok {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", txID))
	}

	// Pre-check expiry: if the transaction has already exceeded its idle
	// or max-life deadline, purgeExpired should have already rolled it
	// back. But the purge ticker has a gap, so an expired-but-not-yet-
	// purged transaction may still be in the map.
	if entry.isExpired(time.Now()) {
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", txID))
	}

	// Wait for outstanding ContinueAfterTimeout background statements
	// before committing/rolling back. The wait is bounded by the Data
	// API's 45-second call timeout. If the background statement has not
	// finished, the transaction stays in the map for retry.
	if !waitForBg(&entry.bgWg, defaultStatementTimeout) {
		return nil, statementTimeout(fmt.Sprintf(
			"transaction %s timed out waiting for a background statement; retry", txID), "")
	}

	// Acquire execMu to wait for any in-flight (non-background)
	// ExecuteStatement or BatchExecuteStatement running on this
	// transaction's sql.Context.
	entry.execMu.Lock()

	// Re-check map presence under s.mu to prevent a double-commit when
	// a concurrent caller already claimed the entry between our initial
	// lookup and acquiring execMu.
	s.mu.Lock()
	if _, stillThere := s.transactions[txID]; !stillThere {
		s.mu.Unlock()
		entry.execMu.Unlock()
		return nil, transactionNotFound(fmt.Sprintf("transaction %s not found or expired", txID))
	}
	delete(s.transactions, txID)
	s.mu.Unlock()

	return entry, nil
}
