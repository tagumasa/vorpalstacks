package dynamodb

import (
	"context"
	"sort"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
)

// ---------------------------------------------------------------------------
// Transaction accounting — the two counter queues an item write feeds and
// the store applies after the carrying transaction commits: contributor
// insights access counters (queued per write, credited per layout under the
// contributor lock) and table metric deltas (item count and size bytes,
// applied under the table's record lock). The commit-time drain is
// DynamoDBStore.Update in store.go; the counter record layouts live in
// contributor_store.go and table_store.go.
// ---------------------------------------------------------------------------

// FlushTableMetrics applies queued per-table metric deltas to the table
// records. Tables are visited in name order and each record is updated
// under its key lock with an immediate read-modify-write, so concurrent
// transactions and the locked TableStore methods cannot lose increments.
// A flush failure is logged and never fails the observed operation — the
// counters are documented approximate values that AWS itself refreshes
// only periodically.
//
// Callers must not hold a table's record lock while queueing metric deltas
// for that table: the flush takes the same lock. No path does — deltas are
// queued exclusively from item-write transactions, which hold no table
// record locks.
func (s *DynamoDBStore) FlushTableMetrics(deltas map[string]TableMetricDelta) {
	if len(deltas) == 0 {
		return
	}
	names := make([]string, 0, len(deltas))
	for name := range deltas {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := s.tables.applyMetricDeltas(name, deltas[name].TableId, deltas[name].ItemCount, deltas[name].SizeBytes); err != nil {
			logs.Warn("failed to apply table metric deltas",
				logs.String("table", name), logs.Err(err))
		}
	}
}

// FlushContributorWrites applies queued contributor events to the access
// counters in one transaction serialised against every other counter
// update, so concurrent reads and writes cannot lose increments. The
// carrying item transaction has already committed: a flush failure is
// logged and never fails the observed operation, and the flush detaches
// from the caller's cancellation (values preserved) — it describes
// already-committed state, so a request ending between the commit and
// this drain must not lose the accounting the commit owes.
func (s *DynamoDBStore) FlushContributorWrites(ctx context.Context, events []ContributorWriteEvent) {
	if len(events) == 0 {
		return
	}
	s.contributorMu.Lock()
	defer s.contributorMu.Unlock()
	err := s.storage.Update(context.WithoutCancel(ctx), func(txn storage.Transaction) error {
		dtxn := &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return dtxn.applyContributorWrites(events)
	})
	if err != nil {
		logs.Warn("failed to record contributor writes", logs.Err(err))
	}
}

// TableMetricDelta accumulates the item-count and size-bytes changes one
// transaction observes for a table. The counters are fields of the whole
// table record, so the deltas are aggregated during the transaction and
// applied once, under the table's record lock, after it commits.
type TableMetricDelta struct {
	// TableId names the generation whose write queued the deltas: the flush
	// applies them only to that generation's record.
	TableId   string
	ItemCount int64
	SizeBytes int64
}

// queueContributorWrite defers the contributor access accounting of an
// item write until the carrying transaction commits. The counters update
// in their own serialised transaction afterwards, which is the only way to
// keep the read-modify-write safe against concurrent updates on a storage
// layer without cross-transaction isolation. Tables without contributor
// insights queue nothing, so their writes open no extra transaction.
func (t *DynamoDBTxn) queueContributorWrite(table *Table, key map[string]*AttributeValue) {
	if table == nil || !table.ContributorInsightsEnabled {
		return
	}
	event := ContributorWriteEvent{TableName: table.Name, Key: make(map[string]*AttributeValue, len(key))}
	for name, value := range key {
		event.Key[name] = value
	}
	t.contributorWrites = append(t.contributorWrites, event)
}

// TakeContributorWrites drains the contributor events queued by the item
// writes of this transaction. The store applies them after the carrying
// transaction commits.
func (t *DynamoDBTxn) TakeContributorWrites() []ContributorWriteEvent {
	events := t.contributorWrites
	t.contributorWrites = nil
	return events
}

// applyContributorWrites credits each queued write event under every
// contributor layout of its table. A write counts as three units of
// ConsumedThroughputUnits. Events on the same counter key are aggregated
// first: the storage layer's read-modify-write reads the committed state,
// so unaggregated repeats inside one transaction would overwrite each
// other.
func (t *DynamoDBTxn) applyContributorWrites(events []ContributorWriteEvent) error {
	tables := make(map[string]*Table)
	keysByTable := make(map[string][]map[string]*AttributeValue)
	var tableOrder []string
	for _, event := range events {
		table := tables[event.TableName]
		if _, resolved := tables[event.TableName]; !resolved {
			fetched, err := t.GetTable(event.TableName)
			if err != nil || fetched == nil || !fetched.ContributorInsightsEnabled {
				fetched = nil
			}
			tables[event.TableName] = fetched
			table = fetched
			if fetched != nil {
				tableOrder = append(tableOrder, event.TableName)
			}
		}
		if table == nil {
			continue
		}
		keysByTable[event.TableName] = append(keysByTable[event.TableName], event.Key)
	}
	at := time.Now()
	for _, name := range tableOrder {
		table := tables[name]
		order, aggregated := contributorAggregateKeys(table, keysByTable[name])
		for _, target := range order {
			if err := RecordAccessTxn(t.txn, t.region(), table.Name, target.layout, target.keyStr, at, aggregated[target], ContributorWriteUnits); err != nil {
				return err
			}
		}
	}
	return nil
}

// contributorAggregateKeys folds one table's keys into per-layout counter
// targets, counting repeats and remembering first-seen order: the storage
// layer's read-modify-write reads the committed state, so unaggregated
// repeats inside one transaction would overwrite each other.
func contributorAggregateKeys(table *Table, keys []map[string]*AttributeValue) ([]contributorTarget, map[contributorTarget]int64) {
	aggregated := make(map[contributorTarget]int64)
	var order []contributorTarget
	for _, key := range keys {
		for _, layout := range ContributorLayouts(table) {
			target := contributorTarget{layout: layout, keyStr: ContributorKeyString(table, key, layout)}
			if _, seen := aggregated[target]; !seen {
				order = append(order, target)
			}
			aggregated[target]++
		}
	}
	return order, aggregated
}

// contributorTarget names one contributor counter: the layout a write is
// credited under and the layout's rendered key.
type contributorTarget struct {
	layout string
	keyStr string
}

// contributorRecordReads credits one read event for every key under each
// contributor layout of the table, inside the caller's transaction. Reads
// on the same counter key (items sharing a partition key) are aggregated
// first: the storage layer's read-modify-write reads the committed state,
// so unaggregated repeats inside one transaction would overwrite each
// other.
func (t *DynamoDBTxn) contributorRecordReads(tableName string, keys []map[string]*AttributeValue) error {
	if len(keys) == 0 {
		return nil
	}
	table, err := t.GetTable(tableName)
	if err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	order, aggregated := contributorAggregateKeys(table, keys)
	at := time.Now()
	for _, target := range order {
		if err := RecordAccessTxn(t.txn, t.region(), table.Name, target.layout, target.keyStr, at, aggregated[target], ContributorReadUnits); err != nil {
			return err
		}
	}
	return nil
}

// contributorRecordQueryEvent credits the single event a Query contributes,
// on the partition-key series only: a Query is one read event regardless of
// how many items it returns, and a result set spans many sort keys.
func (t *DynamoDBTxn) contributorRecordQueryEvent(tableName string, key map[string]*AttributeValue) error {
	table, err := t.GetTable(tableName)
	if err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	keyStr := ContributorKeyString(table, key, ContributorLayoutPartitionKey)
	if keyStr == "" {
		return nil
	}
	return RecordAccessTxn(t.txn, t.region(), table.Name, ContributorLayoutPartitionKey, keyStr, time.Now(), 1, ContributorReadUnits)
}

// RecordContributorReads credits one read event per key under every
// contributor layout of the table. The update runs under the contributor
// lock, so concurrent counter updates cannot lose increments.
func (s *DynamoDBStore) RecordContributorReads(ctx context.Context, tableName string, keys []map[string]*AttributeValue) error {
	if len(keys) == 0 {
		return nil
	}
	if table, err := s.Tables().Get(tableName); err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	s.contributorMu.Lock()
	defer s.contributorMu.Unlock()
	return s.storage.Update(ctx, func(txn storage.Transaction) error {
		dtxn := &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return dtxn.contributorRecordReads(tableName, keys)
	})
}

// RecordContributorQuery credits the single read event a Query contributes
// to the partition-key series. The update runs under the contributor lock,
// so concurrent counter updates cannot lose increments.
func (s *DynamoDBStore) RecordContributorQuery(ctx context.Context, tableName string, key map[string]*AttributeValue) error {
	if table, err := s.Tables().Get(tableName); err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	s.contributorMu.Lock()
	defer s.contributorMu.Unlock()
	return s.storage.Update(ctx, func(txn storage.Transaction) error {
		dtxn := &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return dtxn.contributorRecordQueryEvent(tableName, key)
	})
}

// queueTableMetric records the generation the queued deltas belong to: the
// flush applies them only to that generation's record, so a same-name
// successor created between the write's commit and the flush never inherits
// the old generation's counters.
func (t *DynamoDBTxn) queueTableMetric(tableName string, mutate func(*TableMetricDelta)) error {
	if t.tableMetricDeltas == nil {
		t.tableMetricDeltas = make(map[string]TableMetricDelta)
	}
	d := t.tableMetricDeltas[tableName]
	if d.TableId == "" {
		table, err := t.GetTable(tableName)
		if err != nil {
			return err
		}
		d.TableId = table.TableId
	}
	mutate(&d)
	t.tableMetricDeltas[tableName] = d
	return nil
}

// UpdateItemCount queues an item-count change for a table. The counter is a
// field of the whole table record and the storage layer is read-committed:
// reading and rewriting the record inside the caller's transaction would
// overwrite same-transaction and concurrent increments. The queued delta is
// applied once, under the table's record lock, after the carrying
// transaction commits. Table existence is guaranteed by the item write that
// precedes every counter call on each path.
func (t *DynamoDBTxn) UpdateItemCount(tableName string, delta int64) error {
	return t.queueTableMetric(tableName, func(d *TableMetricDelta) { d.ItemCount += delta })
}

// UpdateTableSize queues a table-size change for a table, aggregated and
// applied like UpdateItemCount after the carrying transaction commits.
func (t *DynamoDBTxn) UpdateTableSize(tableName string, delta int64) error {
	return t.queueTableMetric(tableName, func(d *TableMetricDelta) { d.SizeBytes += delta })
}

// TakeTableMetricDeltas drains the metric deltas queued by this
// transaction. The store applies them after the carrying transaction
// commits.
func (t *DynamoDBTxn) TakeTableMetricDeltas() map[string]TableMetricDelta {
	deltas := t.tableMetricDeltas
	t.tableMetricDeltas = nil
	return deltas
}

// MergeTableMetricDeltas folds every delta of src into dst and returns dst.
// Two-phase transactions run one wrapper per operation; merging their
// queued deltas keeps the post-commit flush one application per table. The
// fold carries the generation id: every wrapper of one transaction resolves
// the id inside the same transaction, so the merged flush still lands under
// the same-generation guard applyMetricDeltas enforces.
func MergeTableMetricDeltas(dst, src map[string]TableMetricDelta) map[string]TableMetricDelta {
	if dst == nil {
		dst = make(map[string]TableMetricDelta, len(src))
	}
	for name, d := range src {
		acc := dst[name]
		acc.ItemCount += d.ItemCount
		acc.SizeBytes += d.SizeBytes
		if acc.TableId == "" {
			acc.TableId = d.TableId
		}
		dst[name] = acc
	}
	return dst
}

// restartOrphanedJobCode marks a job the process interrupted before it
// reached a terminal state: the carrying goroutine died with the process,
// so the record cannot complete on its own.
const restartOrphanedJobCode = "InternalFailure"

// SweepOrphanedJobs moves export and import records stranded in IN_PROGRESS
// by a process restart to FAILED. The jobs run as goroutines inside one
// process — a restart kills every carrier — so a record still IN_PROGRESS
// when a store opens can never complete: leaving it would answer Describe
// forever with a job that died. The sweep runs once per store construction
// (the region store is cached per process), before any request can observe
// the stale state.
func (s *DynamoDBStore) SweepOrphanedJobs() {
	if err := forEachRecordPage(
		func(marker string) ([]*ExportDescription, string, error) { return s.exports.List("", marker, 100) },
		func(export *ExportDescription) error {
			if export.ExportStatus != ExportStatusInProgress {
				return nil
			}
			export.ExportStatus = ExportStatusFailed
			export.FailureCode = restartOrphanedJobCode
			export.FailureMessage = "the export was interrupted by a process restart"
			export.EndTime = time.Now().UTC()
			if err := s.exports.Put(export); err != nil {
				logs.Warn("DynamoDB: failed to fail a restart-orphaned export",
					logs.String("exportArn", export.ExportArn), logs.Err(err))
			}
			return nil
		},
	); err != nil {
		logs.Warn("DynamoDB: failed to scan exports for restart-orphaned jobs", logs.Err(err))
	}
	if err := forEachRecordPage(
		func(marker string) ([]*ImportTableDescription, string, error) { return s.imports.List("", marker, 100) },
		func(imp *ImportTableDescription) error {
			if imp.ImportStatus != ImportStatusInProgress {
				return nil
			}
			imp.ImportStatus = ImportStatusFailed
			imp.FailureCode = restartOrphanedJobCode
			imp.FailureMessage = "the import was interrupted by a process restart"
			imp.EndTime = time.Now().UTC()
			if err := s.imports.Put(imp); err != nil {
				logs.Warn("DynamoDB: failed to fail a restart-orphaned import",
					logs.String("importArn", imp.ImportArn), logs.Err(err))
			}
			return nil
		},
	); err != nil {
		logs.Warn("DynamoDB: failed to scan imports for restart-orphaned jobs", logs.Err(err))
	}
}

// forEachRecordPage walks one family's List pagination to exhaustion,
// handing each record of every page to fn. The page size stays the store's
// List default; the walk ends at the empty continuation marker.
func forEachRecordPage[T any](list func(marker string) ([]T, string, error), fn func(T) error) error {
	for marker := ""; ; {
		records, next, err := list(marker)
		if err != nil {
			return err
		}
		for _, record := range records {
			if err := fn(record); err != nil {
				return err
			}
		}
		if next == "" {
			return nil
		}
		marker = next
	}
}
