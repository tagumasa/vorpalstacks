package cloudwatchlogs

import (
	"errors"
	"sync"
)

// PendingDelivery is one subscription-filter batch whose delivery failed
// and now rides the documented retry window: "Throttled deliverables are
// retried for up to 24 hours. After 24 hours, the failed deliverables are
// dropped" (CloudWatch Logs User Guide, subscription filters). The record
// carries everything a re-drive needs — the destination addressing, the
// matched filter's distribution semantics, and the already-compressed
// payload — plus the window bookkeeping: the first attempt, the drop
// deadline, and the next scheduled attempt under the service's backoff.
type PendingDelivery struct {
	ID           string `json:"id"`
	Region       string `json:"region"`
	DestArn      string `json:"destArn"`
	LogGroup     string `json:"logGroup"`
	LogStream    string `json:"logStream"`
	Distribution string `json:"distribution"`
	Payload      []byte `json:"payload"`
	FirstAttempt int64  `json:"firstAttempt"`
	Deadline     int64  `json:"deadline"`
	NextAttempt  int64  `json:"nextAttempt"`
	Attempts     int    `json:"attempts"`
}

func (s *Store) pendingDeliveryKey(id string) string {
	return keyPrefixPendingDelivery + id
}

// PutPendingDelivery persists or updates one pending delivery. A re-drive
// failure rewrites the same id, so attempts accumulate on one record
// instead of piling up duplicates of the same batch.
func (s *Store) PutPendingDelivery(pd *PendingDelivery) error {
	return s.Put(s.pendingDeliveryKey(pd.ID), pd)
}

// GetPendingDelivery reads one pending delivery by id; absent is the
// not-found sentinel.
func (s *Store) GetPendingDelivery(id string) (*PendingDelivery, error) {
	return getJSONRecord[PendingDelivery](s, s.pendingDeliveryKey(id), ErrResourceNotFound)
}

// DeletePendingDelivery removes a pending delivery after a successful
// re-drive or a deadline drop.
func (s *Store) DeletePendingDelivery(id string) error {
	return s.Delete(s.pendingDeliveryKey(id))
}

// pendingDeliveryMu serialises read-modify-write cycles on pending
// delivery records: the ingestion failure path's attempts merge and the
// retry drain's reschedule decide on the persisted record under this
// mutex, so concurrent writers cannot lose an attempts increment or
// revert each other's schedule.
var pendingDeliveryMu sync.Mutex

// UpsertPendingDelivery merges one pending delivery atomically: fn
// receives the persisted record for the id, or nil when none rides the
// window yet, and returns the record to persist — a nil return writes
// nothing. The read, the merge and the write commit under the pending
// record mutex, so two failure paths of the same batch accumulate
// attempts instead of racing one increment away.
func (s *Store) UpsertPendingDelivery(id string, fn func(existing *PendingDelivery) *PendingDelivery) error {
	pendingDeliveryMu.Lock()
	defer pendingDeliveryMu.Unlock()
	existing, err := s.GetPendingDelivery(id)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return err
	}
	pd := fn(existing)
	if pd == nil {
		return nil
	}
	return s.PutPendingDelivery(pd)
}

// ListPendingDeliveries returns every pending delivery of the store's
// region. The set is bounded by failure volume and each record carries its
// own schedule, so a plain prefix scan is the whole query; the due and
// deadline filtering lives with the service, where the clock lives.
func (s *Store) ListPendingDeliveries() ([]*PendingDelivery, error) {
	return listJSONRecords[PendingDelivery](s, keyPrefixPendingDelivery, "pending delivery record", nil)
}
