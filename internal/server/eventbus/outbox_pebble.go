package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/pebble/v2"
	"vorpalstacks/internal/core/logs"
)

const (
	outboxKeyPrefix  = "eb/outbox/"
	createdIdxPrefix = "eb/idx/created/"
)

type PebbleOutboxStore struct {
	db  *pebble.DB
	log logs.Logger
}

func NewPebbleOutboxStore(db *pebble.DB) *PebbleOutboxStore {
	return &PebbleOutboxStore{db: db}
}

func (s *PebbleOutboxStore) Write(ctx context.Context, entry *OutboxEntry) error {
	key := outboxKey(entry.EventID)

	val, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("pebble outbox: marshal: %w", err)
	}

	if err := s.db.Set(key, val, pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: write: %w", err)
	}

	idxKey := createdIdxKey(entry.CreatedAt, entry.EventID)
	if err := s.db.Set(idxKey, nil, pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: write index: %w", err)
	}

	return nil
}

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

func (s *PebbleOutboxStore) UpdateStatus(ctx context.Context, eventID string, from, to OutboxStatus) (bool, error) {
	key := outboxKey(eventID)

	val, closer, err := s.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("pebble outbox: get for CAS: %w", err)
	}
	defer closer.Close()

	var entry OutboxEntry
	if err := json.Unmarshal(val, &entry); err != nil {
		return false, fmt.Errorf("pebble outbox: unmarshal for CAS: %w", err)
	}

	if entry.Status != from {
		return false, nil
	}

	entry.Status = to
	if to == OutboxDelivered {
		now := time.Now().UTC()
		entry.DeliveredAt = &now
	}

	newVal, err := json.Marshal(&entry)
	if err != nil {
		return false, fmt.Errorf("pebble outbox: marshal for CAS: %w", err)
	}

	if err := s.db.Set(key, newVal, pebble.NoSync); err != nil {
		return false, fmt.Errorf("pebble outbox: write for CAS: %w", err)
	}

	return true, nil
}

func (s *PebbleOutboxStore) UpdateEntry(ctx context.Context, entry *OutboxEntry) error {
	key := outboxKey(entry.EventID)

	val, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("pebble outbox: marshal update: %w", err)
	}

	if err := s.db.Set(key, val, pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: write update: %w", err)
	}

	return nil
}

func (s *PebbleOutboxStore) ListPending(ctx context.Context, limit int) ([]*OutboxEntry, error) {
	var entries []*OutboxEntry

	iterOpts := &pebble.IterOptions{
		LowerBound: []byte(outboxKeyPrefix),
		UpperBound: []byte(outboxKey("zzzzzzzzz")),
	}

	iter, err := s.db.NewIter(iterOpts)
	if err != nil {
		return nil, fmt.Errorf("pebble outbox: create iterator: %w", err)
	}
	defer iter.Close()

	count := 0
	for iter.First(); iter.Valid(); iter.Next() {
		if count >= limit {
			break
		}

		var entry OutboxEntry
		if err := json.Unmarshal(iter.Value(), &entry); err != nil {
			s.logWarn("skipping malformed outbox entry", "key", string(iter.Key()), "error", err)
			continue
		}

		if entry.Status == OutboxPending {
			entries = append(entries, &entry)
			count++
		}
	}

	if err := iter.Error(); err != nil {
		return entries, fmt.Errorf("pebble outbox: iterate pending: %w", err)
	}

	return entries, nil
}

func (s *PebbleOutboxStore) Delete(ctx context.Context, eventID string) error {
	key := outboxKey(eventID)

	val, closer, err := s.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil
		}
		return fmt.Errorf("pebble outbox: get for delete: %w", err)
	}
	defer closer.Close()

	var entry OutboxEntry
	if err := json.Unmarshal(val, &entry); err != nil {
		return fmt.Errorf("pebble outbox: unmarshal for delete: %w", err)
	}

	if err := s.db.Delete(key, pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: delete entry: %w", err)
	}

	idxKey := createdIdxKey(entry.CreatedAt, eventID)
	if err := s.db.Delete(idxKey, pebble.NoSync); err != nil {
		return fmt.Errorf("pebble outbox: delete index: %w", err)
	}

	return nil
}

func (s *PebbleOutboxStore) Cleanup(ctx context.Context, deliveredBefore time.Time, failedBefore time.Time) (int, error) {
	lower := []byte(createdIdxPrefix)
	upper := []byte(createdIdxPrefix + "z")

	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: lower,
		UpperBound: upper,
	})
	if err != nil {
		return 0, fmt.Errorf("pebble outbox: cleanup iterator: %w", err)
	}
	defer iter.Close()

	type deleteEntry struct {
		eventID string
		idxKey  string
	}
	var toDelete []deleteEntry

	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())

		parts := strings.SplitN(strings.TrimPrefix(key, createdIdxPrefix), "/", 2)
		if len(parts) != 2 {
			continue
		}
		eventID := parts[1]

		entryKey := outboxKey(eventID)
		entryVal, closer, err := s.db.Get(entryKey)
		if err != nil {
			if err == pebble.ErrNotFound {
				_ = s.db.Delete([]byte(key), pebble.NoSync)
				continue
			}
			continue
		}

		var entry OutboxEntry
		if err := json.Unmarshal(entryVal, &entry); err != nil {
			closer.Close()
			continue
		}
		closer.Close()

		switch entry.Status {
		case OutboxDelivered:
			if entry.DeliveredAt != nil && entry.DeliveredAt.Before(deliveredBefore) {
				toDelete = append(toDelete, deleteEntry{eventID: eventID, idxKey: key})
			}
		case OutboxFailed:
			createdCutoff := failedBefore
			if entry.CreatedAt.Before(createdCutoff) {
				toDelete = append(toDelete, deleteEntry{eventID: eventID, idxKey: key})
			}
		}
	}

	if err := iter.Error(); err != nil {
		return 0, fmt.Errorf("pebble outbox: cleanup iterate: %w", err)
	}

	count := 0
	batch := s.db.NewBatch()
	for _, item := range toDelete {
		if err := batch.Delete(outboxKey(item.eventID), nil); err != nil {
			continue
		}
		_ = batch.Delete([]byte(item.idxKey), nil)
		count++
	}
	if count > 0 {
		if err := batch.Commit(pebble.NoSync); err != nil {
			return 0, fmt.Errorf("pebble outbox: cleanup commit: %w", err)
		}
	}

	return count, nil
}

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

func createdIdxKey(createdAt time.Time, eventID string) []byte {
	ts := fmt.Sprintf("%020d", createdAt.UnixNano())
	return []byte(createdIdxPrefix + ts + "/" + eventID)
}

var _ OutboxStore = (*PebbleOutboxStore)(nil)
