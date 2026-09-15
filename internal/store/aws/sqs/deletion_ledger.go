package sqs

import (
	"strconv"
	"time"

	"vorpalstacks/internal/core/logs"
)

// deletedRecently reports whether a queue with the given URL was deleted
// within the recreate-prohibition window. Stale markers are dropped
// opportunistically so the ledger does not grow without bound.
func (s *SQSStore) deletedRecently(queueURL string) bool {
	bucket := s.storage.Bucket(s.deletionsBucket)
	raw, err := bucket.Get([]byte(queueURL))
	if err != nil {
		// Fail-open by design, never silently: the miss un-arms the
		// recreate-prohibition window for this probe.
		logs.Warn("SQS: deletion-ledger read failed; the recreate window is unenforced for this probe", logs.String("queueURL", queueURL), logs.Err(err))
		return false
	}
	if len(raw) == 0 {
		return false
	}
	deletedAt, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return false
	}
	if time.Since(time.Unix(deletedAt, 0)) < queueDeletionWindow {
		return true
	}
	if err := bucket.Delete([]byte(queueURL)); err != nil {
		// The periodic marker sweep reclaims it; the probe's answer is
		// unaffected. Logged for the same reason as every other fail-open.
		logs.Warn("SQS: stale deletion-marker drop failed; the sweep reclaims it later", logs.String("queueURL", queueURL), logs.Err(err))
	}
	return false
}

// recordQueueDeletion persists the deletion timestamp of a queue so that
// same-name recreation inside the prohibition window can be rejected across
// restarts. The queue deletion itself has already happened when this runs,
// so a failed marker write cannot fail the delete — but it silently voids
// the documented 60-second prohibition for this name, which deserves a loud
// log rather than silence.
func (s *SQSStore) recordQueueDeletion(queueURL string) {
	bucket := s.storage.Bucket(s.deletionsBucket)
	if err := bucket.Put([]byte(queueURL), []byte(strconv.FormatInt(time.Now().Unix(), 10))); err != nil {
		logs.Warn("SQS: failed to persist queue-deletion marker; the recreate-prohibition window is unenforced for this queue", logs.String("queueURL", queueURL), logs.Err(err))
	}
}
