package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/pebble/v2"
	"vorpalstacks/internal/core/logs"
)

const (
	outboxKeyPrefix = "eb/outbox/"
	// statusIdxPrefix is the secondary index over outbox status. Keys are
	// statusIdxPrefix + <STATUS> + "/" + <zero-padded unix nano> + "/" + eventID,
	// where the timestamp is the delivery time for DELIVERED entries and the
	// creation time for every other status. The fixed-width timestamp makes
	// key order match time order, so the pending scan reads only pending
	// entries and the retention cleanup can identify purge candidates from
	// the key range alone, without decoding any outbox record.
	statusIdxPrefix     = "eb/idx/st/"
	statusIdxTimeFormat = "%020d"
)

// PebbleOutboxStore implements OutboxStore backed by a Pebble key-value
// database, using a per-status secondary index for pending scans and
// time-based cleanup queries.
type PebbleOutboxStore struct {
	db  *pebble.DB
	log logs.Logger

	// statusMu provides per-eventID mutual exclusion so that
	// UpdateStatus's read-modify-write is atomic across concurrent
	// workers. Pebble does not provide native compare-and-set, so we
	// serialise the CAS in-process.
	statusMu sync.Map
}

// NewPebbleOutboxStore creates a PebbleOutboxStore backed by the given
// Pebble database instance.
func NewPebbleOutboxStore(db *pebble.DB) *PebbleOutboxStore {
	return &PebbleOutboxStore{db: db}
}

// tsForIndex returns the timestamp that indexes the entry's current
// status: the delivery time once delivered, the creation time otherwise.
func tsForIndex(entry *OutboxEntry) time.Time {
	if entry.Status == OutboxDelivered && entry.DeliveredAt != nil {
		return *entry.DeliveredAt
	}
	return entry.CreatedAt
}

func statusIdxKey(status OutboxStatus, ts time.Time, eventID string) []byte {
	return []byte(fmt.Sprintf("%s%s/"+statusIdxTimeFormat+"/%s",
		statusIdxPrefix, status.String(), ts.UnixNano(), eventID))
}

func statusIdxKeyFor(entry *OutboxEntry) []byte {
	return statusIdxKey(entry.Status, tsForIndex(entry), entry.EventID)
}

// eventIDFromStatusIdxKey extracts the event ID from the tail of a status
// index key. Event IDs may themselves contain slashes, so everything after
// the timestamp segment belongs to the ID.
func eventIDFromStatusIdxKey(key string) string {
	rest := strings.TrimPrefix(key, statusIdxPrefix)
	parts := strings.SplitN(rest, "/", 3)
	if len(parts) != 3 {
		return ""
	}
	return parts[2]
}

// prefixUpper returns an exclusive upper bound covering every key that
// starts with the given prefix.
func prefixUpper(prefix string) []byte {
	upper := []byte(prefix)
	for i := len(upper) - 1; i >= 0; i-- {
		if upper[i] < 0xFF {
			upper = upper[:i+1]
			upper[i]++
			return upper
		}
	}
	return nil
}

// Write persists an outbox entry together with its status index entry in
// one atomic batch. The batch is fsynced so the outbox's at-least-once
// durability guarantee holds for the record and the index alike.
func (s *PebbleOutboxStore) Write(ctx context.Context, entry *OutboxEntry) error {
	val, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("pebble outbox: marshal: %w", err)
	}

	batch := s.db.NewBatch()
	defer batch.Close()
	if err := batch.Set(outboxKey(entry.EventID), val, nil); err != nil {
		return fmt.Errorf("pebble outbox: write: %w", err)
	}
	if err := batch.Set(statusIdxKeyFor(entry), nil, nil); err != nil {
		return fmt.Errorf("pebble outbox: write index: %w", err)
	}
	if err := batch.Commit(pebble.Sync); err != nil {
		return fmt.Errorf("pebble outbox: write: %w", err)
	}
	return nil
}

// Read retrieves a single outbox entry by event ID.
func (s *PebbleOutboxStore) Read(ctx context.Context, eventID string) (*OutboxEntry, error) {
	key := outboxKey(eventID)
	val, closer, err := s.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("pebble outbox: read: %w", err)
	}
	defer closer.Close()

	var entry OutboxEntry
	if err := json.Unmarshal(val, &entry); err != nil {
		return nil, fmt.Errorf("pebble outbox: unmarshal: %w", err)
	}

	return &entry, nil
}

// UpdateStatus atomically transitions an entry's status from 'from' to
// 'to', moving its status index entry in the same batch. Atomicity is
// guaranteed by a per-eventID in-process lock, because Pebble does not
// provide native compare-and-set semantics.
//
// When the destination status is terminal (Delivered or Failed) the
// per-eventID mutex is dropped from the statusMu map after the lock is
// released so the map cannot grow without bound over the lifetime of the
// process. The map entry is removed only after the lock is released so
// that any concurrent goroutine already blocked on the same mutex keeps
// the old pointer and wakes up correctly; new callers will allocate a
// fresh mutex, which is safe because terminal events must not transition
// again.
func (s *PebbleOutboxStore) UpdateStatus(ctx context.Context, eventID string, from, to OutboxStatus) (bool, error) {
	transitioned, err := s.transitionUnderLock(eventID, from, to)
	if err != nil {
		return false, err
	}
	if transitioned && (to == OutboxDelivered || to == OutboxFailed) {
		s.statusMu.Delete(eventID)
	}
	return transitioned, nil
}

// transitionUnderLock performs the compare-and-set under the per-eventID
// mutex and reports whether the transition actually occurred.
func (s *PebbleOutboxStore) transitionUnderLock(eventID string, from, to OutboxStatus) (bool, error) {
	unlock := s.lockEvent(eventID)
	defer unlock()

	key := outboxKey(eventID)

	val, closer, err := s.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("pebble outbox: get for CAS: %w", err)
	}

	var entry OutboxEntry
	if err := json.Unmarshal(val, &entry); err != nil {
		closer.Close()
		return false, fmt.Errorf("pebble outbox: unmarshal for CAS: %w", err)
	}
	closer.Close()

	if entry.Status != from {
		return false, nil
	}

	oldIdx := statusIdxKeyFor(&entry)
	entry.Status = to
	if to == OutboxDelivered {
		now := time.Now().UTC()
		entry.DeliveredAt = &now
	}

	newVal, err := json.Marshal(&entry)
	if err != nil {
		return false, fmt.Errorf("pebble outbox: marshal for CAS: %w", err)
	}

	batch := s.db.NewBatch()
	defer batch.Close()
	if err := batch.Set(key, newVal, nil); err != nil {
		return false, fmt.Errorf("pebble outbox: write for CAS: %w", err)
	}
	if err := batch.Delete(oldIdx, nil); err != nil {
		return false, fmt.Errorf("pebble outbox: index for CAS: %w", err)
	}
	if err := batch.Set(statusIdxKeyFor(&entry), nil, nil); err != nil {
		return false, fmt.Errorf("pebble outbox: index for CAS: %w", err)
	}
	if err := batch.Commit(pebble.Sync); err != nil {
		return false, fmt.Errorf("pebble outbox: write for CAS: %w", err)
	}

	return true, nil
}

// lockEvent acquires a per-eventID mutex and returns an unlock function.
// UpdateStatus drops the mutex from the map once the event reaches a
// terminal state so the map remains bounded.
func (s *PebbleOutboxStore) lockEvent(eventID string) func() {
	v, _ := s.statusMu.LoadOrStore(eventID, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// UpdateEntry unconditionally overwrites the stored outbox entry and moves
// its status index entry to match. The previous record is read under the
// per-eventID lock so the old index key is removed in the same batch; the
// lock entry is dropped when the new status is terminal.
func (s *PebbleOutboxStore) UpdateEntry(ctx context.Context, entry *OutboxEntry) error {
	unlock := s.lockEvent(entry.EventID)
	defer unlock()

	val, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("pebble outbox: marshal update: %w", err)
	}

	key := outboxKey(entry.EventID)
	batch := s.db.NewBatch()
	defer batch.Close()

	if oldVal, closer, getErr := s.db.Get(key); getErr == nil {
		var old OutboxEntry
		if jsonErr := json.Unmarshal(oldVal, &old); jsonErr == nil {
			_ = batch.Delete(statusIdxKeyFor(&old), nil)
		}
		closer.Close()
	} else if getErr != pebble.ErrNotFound {
		return fmt.Errorf("pebble outbox: get for update: %w", getErr)
	}

	if err := batch.Set(key, val, nil); err != nil {
		return fmt.Errorf("pebble outbox: write update: %w", err)
	}
	if err := batch.Set(statusIdxKeyFor(entry), nil, nil); err != nil {
		return fmt.Errorf("pebble outbox: index update: %w", err)
	}
	if err := batch.Commit(pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: write update: %w", err)
	}

	if entry.Status == OutboxDelivered || entry.Status == OutboxFailed {
		s.statusMu.Delete(entry.EventID)
	}
	return nil
}

// ListPending returns up to 'limit' outbox entries with status OutboxPending.
func (s *PebbleOutboxStore) ListPending(ctx context.Context, limit int) ([]*OutboxEntry, error) {
	entries, _, err := s.ListPendingFrom(ctx, limit, "")
	return entries, err
}

// ListPendingFrom returns up to 'limit' pending entries plus the opaque
// cursor of the last returned entry. Passing the cursor back resumes
// strictly after it; an empty cursor starts from the head of the pending
// set, so pages tile the backlog without overlap or gaps. Only pending
// index keys and pending records are read: delivered and failed entries
// inside their retention windows cost the scan nothing.
func (s *PebbleOutboxStore) ListPendingFrom(ctx context.Context, limit int, afterCursor string) ([]*OutboxEntry, string, error) {
	prefix := statusIdxPrefix + OutboxPending.String() + "/"

	var lower []byte
	if afterCursor != "" {
		// The strict successor of the cursor key: excludes the cursor
		// entry itself while staying inside the pending index range.
		lower = append([]byte(afterCursor), 0x00)
	} else {
		lower = []byte(prefix)
	}

	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: lower,
		UpperBound: prefixUpper(prefix),
	})
	if err != nil {
		return nil, "", fmt.Errorf("pebble outbox: create iterator: %w", err)
	}

	var entries []*OutboxEntry
	var staleIdxKeys []string
	cursor := ""
	for iter.First(); iter.Valid(); iter.Next() {
		if len(entries) >= limit {
			break
		}

		key := string(iter.Key())
		eventID := eventIDFromStatusIdxKey(key)
		if eventID == "" {
			continue
		}

		val, closer, getErr := s.db.Get(outboxKey(eventID))
		if getErr == pebble.ErrNotFound {
			// Index entry outlived its record; drop it so it stops
			// occupying the pending scan.
			staleIdxKeys = append(staleIdxKeys, key)
			continue
		}
		if getErr != nil {
			iter.Close()
			return entries, cursor, fmt.Errorf("pebble outbox: iterate pending: %w", getErr)
		}

		var entry OutboxEntry
		jsonErr := json.Unmarshal(val, &entry)
		closer.Close()
		if jsonErr != nil || entry.Status != OutboxPending {
			staleIdxKeys = append(staleIdxKeys, key)
			continue
		}

		entries = append(entries, &entry)
		cursor = key
	}

	iterErr := iter.Error()
	iter.Close()
	if iterErr != nil {
		return entries, cursor, fmt.Errorf("pebble outbox: iterate pending: %w", iterErr)
	}

	for _, key := range staleIdxKeys {
		if err := s.db.Delete([]byte(key), pebble.NoSync); err != nil {
			s.logWarn("failed to remove stale status index entry", "key", key, "error", err)
		}
	}

	return entries, cursor, nil
}

// ResetStaleProcessing finds entries in OutboxProcessing status and
// resets them to OutboxPending. This is called during EventBus Start to
// recover entries left behind by a crash or panic during processing. The
// processing status index supplies the candidate set, so no other outbox
// record is decoded.
func (s *PebbleOutboxStore) ResetStaleProcessing(ctx context.Context) (int, error) {
	prefix := statusIdxPrefix + OutboxProcessing.String() + "/"

	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: []byte(prefix),
		UpperBound: prefixUpper(prefix),
	})
	if err != nil {
		return 0, fmt.Errorf("pebble outbox: create iterator (reset stale): %w", err)
	}

	type recovery struct {
		entry  OutboxEntry
		idxKey string
	}
	var found []recovery
	var staleIdxKeys []string
	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())
		eventID := eventIDFromStatusIdxKey(key)
		if eventID == "" {
			continue
		}

		val, closer, getErr := s.db.Get(outboxKey(eventID))
		if getErr == pebble.ErrNotFound {
			staleIdxKeys = append(staleIdxKeys, key)
			continue
		}
		if getErr != nil {
			iter.Close()
			return len(found), fmt.Errorf("pebble outbox: iterate stale reset: %w", getErr)
		}

		var entry OutboxEntry
		jsonErr := json.Unmarshal(val, &entry)
		closer.Close()
		if jsonErr != nil || entry.Status != OutboxProcessing {
			staleIdxKeys = append(staleIdxKeys, key)
			continue
		}
		found = append(found, recovery{entry: entry, idxKey: key})
	}

	iterErr := iter.Error()
	iter.Close()
	if iterErr != nil {
		return 0, fmt.Errorf("pebble outbox: iterate stale reset: %w", iterErr)
	}

	count := 0
	for _, r := range found {
		r.entry.Status = OutboxPending
		newVal, err := json.Marshal(&r.entry)
		if err != nil {
			s.logWarn("failed to marshal stale entry", "event_id", r.entry.EventID, "error", err)
			continue
		}

		batch := s.db.NewBatch()
		if err := batch.Set(outboxKey(r.entry.EventID), newVal, nil); err != nil {
			batch.Close()
			s.logWarn("failed to persist stale entry reset", "event_id", r.entry.EventID, "error", err)
			continue
		}
		_ = batch.Delete([]byte(r.idxKey), nil)
		_ = batch.Set(statusIdxKeyFor(&r.entry), nil, nil)
		if err := batch.Commit(pebble.Sync); err != nil {
			batch.Close()
			s.logWarn("failed to persist stale entry reset", "event_id", r.entry.EventID, "error", err)
			continue
		}
		batch.Close()
		count++
	}

	for _, key := range staleIdxKeys {
		if err := s.db.Delete([]byte(key), pebble.NoSync); err != nil {
			s.logWarn("failed to remove stale status index entry", "key", key, "error", err)
		}
	}

	return count, nil
}

// Delete removes an outbox entry and its status index entry.
func (s *PebbleOutboxStore) Delete(ctx context.Context, eventID string) error {
	key := outboxKey(eventID)

	val, closer, err := s.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil
		}
		return fmt.Errorf("pebble outbox: get for delete: %w", err)
	}

	var entry OutboxEntry
	if err := json.Unmarshal(val, &entry); err != nil {
		closer.Close()
		return fmt.Errorf("pebble outbox: unmarshal for delete: %w", err)
	}
	closer.Close()

	batch := s.db.NewBatch()
	defer batch.Close()
	if err := batch.Delete(key, nil); err != nil {
		return fmt.Errorf("pebble outbox: delete entry: %w", err)
	}
	if err := batch.Delete(statusIdxKeyFor(&entry), nil); err != nil {
		return fmt.Errorf("pebble outbox: delete index: %w", err)
	}
	if err := batch.Commit(pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: delete: %w", err)
	}
	return nil
}

// Cleanup purges delivered entries older than deliveredBefore and failed
// entries older than failedBefore. Pending and processing entries are
// never purged: they hold undelivered events and must survive until they
// deliver or transition to Failed, which retention eventually reaps. The
// status index keys embed the relevant timestamp, so purge candidates are
// identified from the key range alone and no outbox record is decoded.
func (s *PebbleOutboxStore) Cleanup(ctx context.Context, deliveredBefore time.Time, failedBefore time.Time) (int, error) {
	batch := s.db.NewBatch()
	defer batch.Close()

	count := 0
	for _, target := range []struct {
		status OutboxStatus
		before time.Time
	}{
		{OutboxDelivered, deliveredBefore},
		{OutboxFailed, failedBefore},
	} {
		prefix := statusIdxPrefix + target.status.String() + "/"
		upper := []byte(prefix + fmt.Sprintf(statusIdxTimeFormat, target.before.UnixNano()))

		iter, err := s.db.NewIter(&pebble.IterOptions{
			LowerBound: []byte(prefix),
			UpperBound: upper,
		})
		if err != nil {
			return 0, fmt.Errorf("pebble outbox: cleanup iterator: %w", err)
		}
		for iter.First(); iter.Valid(); iter.Next() {
			eventID := eventIDFromStatusIdxKey(string(iter.Key()))
			if eventID == "" {
				continue
			}
			if err := batch.Delete(outboxKey(eventID), nil); err != nil {
				continue
			}
			// Copy the key: iter.Key is only valid until the next advance.
			_ = batch.Delete(append([]byte(nil), iter.Key()...), nil)
			count++
		}
		if err := iter.Error(); err != nil {
			iter.Close()
			return 0, fmt.Errorf("pebble outbox: cleanup iterate: %w", err)
		}
		iter.Close()
	}

	if count > 0 {
		if err := batch.Commit(pebble.NoSync); err != nil {
			return 0, fmt.Errorf("pebble outbox: cleanup commit: %w", err)
		}
	}

	return count, nil
}

// Close is a no-op for PebbleOutboxStore as the underlying database is
// managed externally.
func (s *PebbleOutboxStore) Close() error {
	return nil
}

func (s *PebbleOutboxStore) logWarn(msg string, keyvals ...interface{}) {
	if s.log != nil {
		fields := make([]logs.Field, 0, len(keyvals)/2)
		for i := 0; i+1 < len(keyvals); i += 2 {
			fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
		}
		s.log.Warn(msg, fields...)
	}
}

func outboxKey(eventID string) []byte {
	return []byte(outboxKeyPrefix + eventID)
}

var _ OutboxStore = (*PebbleOutboxStore)(nil)
