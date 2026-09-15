package sqs

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

// cleanupExpiredMessages periodically scans all queues and deletes messages
// that have exceeded their MessageRetentionPeriod. Expired messages are never
// delivered by ReceiveMessage (filtered by isMessageExpired), so this cleanup
// can run without holding msgMutex.
func (s *SQSStore) cleanupExpiredMessages() {
	defer s.wg.Done()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.runCleanupTick()
		}
	}
}

// runCleanupTick runs one full cleanup pass. Each sweep recovers its own
// panics: a transient panic in one sweep is logged and the pass continues —
// one panic must not kill retention cleanup, receipt sweeping and dedup
// sweeping for the process lifetime.
func (s *SQSStore) runCleanupTick() {
	s.runSweep("retention", s.doMessageRetentionCleanup)
	s.runSweep("receipts", s.doReceiptHandleCleanup)
	s.runSweep("deduplication", s.doDeduplicationCleanup)
	s.runSweep("deletion markers", s.doDeletionMarkerCleanup)
}

// runSweep executes one cleanup pass, logging and swallowing a panic so the
// caller's schedule keeps running.
func (s *SQSStore) runSweep(name string, sweep func()) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("SQS: panic in "+name+" cleanup sweep", logs.Any("panic", r))
		}
	}()
	sweep()
}

// doMessageRetentionCleanup scans ALL queues (with pagination) and deletes
// expired messages from each. Expired messages are never delivered by
// ReceiveMessage (filtered by isMessageExpired), so this cleanup can run
// without holding msgMutex. Errors are logged but do not stop the cleanup.
func (s *SQSStore) doMessageRetentionCleanup() {
	now := time.Now().UTC()
	const pageSize = 100
	var marker string

	for {
		if s.ctx.Err() != nil {
			return
		}

		opts := common.ListOptions{MaxItems: pageSize}
		if marker != "" {
			opts.Marker = marker
		}

		result, err := s.ListQueues(opts, "")
		if err != nil {
			logs.Warn("SQS retention cleanup: ListQueues failed", logs.Err(err))
			return
		}

		for _, queue := range result.Items {
			if s.ctx.Err() != nil {
				return
			}
			retentionCutoff := now.Add(-time.Duration(queue.MessageRetentionPeriod) * time.Second)
			prefix := messagePrefix(queue.URL)

			var expiredKeys []string
			if err := common.ForEachAllProto[*pb.Message](s.messagesStore, prefix, func() *pb.Message { return &pb.Message{} }, nil, func(msgPb *pb.Message) error {
				if s.isMessageExpired(msgPb, retentionCutoff) {
					expiredKeys = append(expiredKeys, messageKey(queue.URL, msgPb.Id))
				}
				return nil
			}); err != nil {
				logs.Warn("SQS retention cleanup: message scan failed", logs.String("queue", queue.URL), logs.Err(err))
				continue
			}

			for _, key := range expiredKeys {
				if err := s.messagesStore.Delete(key); err != nil {
					logs.Warn("SQS retention cleanup: delete failed", logs.String("key", key), logs.Err(err))
				}
			}
		}

		if !result.IsTruncated || result.NextMarker == "" {
			break
		}
		marker = result.NextMarker
	}
}

// doReceiptHandleCleanup deletes receipt-handle entries older than
// receiptHandleRetention. Receipt handles embed their issue instant
// ("<uuid>#<unix-nano>", see generateReceiptHandle), so the age is read from
// the key alone; entries with an unparseable key are left untouched. Every
// receive adds an entry, and old handles that are never reused have no other
// reclamation path on queues that are neither purged nor deleted, so this
// sweep is what bounds the receipts bucket.
func (s *SQSStore) doReceiptHandleCleanup() {
	if s.ctx.Err() != nil {
		return
	}
	cutoff := time.Now().UTC().Add(-receiptHandleRetention)
	receiptsBucket := s.storage.Bucket(s.receiptsBucket)
	var stale [][]byte
	if err := receiptsBucket.ForEach(func(k, _ []byte) error {
		sep := strings.LastIndexByte(string(k), '#')
		if sep < 0 || sep == len(k)-1 {
			return nil
		}
		issuedNano, err := strconv.ParseInt(string(k[sep+1:]), 10, 64)
		if err != nil {
			return nil
		}
		if time.Unix(0, issuedNano).Before(cutoff) {
			keyCopy := make([]byte, len(k))
			copy(keyCopy, k)
			stale = append(stale, keyCopy)
		}
		return nil
	}); err != nil {
		logs.Warn("SQS receipt cleanup: scan failed; stale entries stay until the next sweep", logs.Err(err))
	}
	for _, k := range stale {
		if err := receiptsBucket.Delete(k); err != nil {
			logs.Warn("SQS receipt cleanup: delete failed", logs.String("handle", string(k)), logs.Err(err))
		}
	}
}

// doDeduplicationCleanup deletes persisted deduplication entries whose
// interval has expired. Entries carry their expiry in the value's final
// segment ("<messageKey>\x01<sequenceNumber>\x01<unix-millis>", see
// deduplicationEntryValue), so the sweep reads it from the tail alone;
// entries without a parseable expiry are left to the expired-re-read path.
// FIFO queues accumulate one entry per distinct dedup key, and only an
// expired re-read or DeleteQueue removed them, so this sweep is what bounds
// the bucket.
func (s *SQSStore) doDeduplicationCleanup() {
	if s.ctx.Err() != nil {
		return
	}
	nowMs := time.Now().UnixMilli()
	bucket := s.storage.Bucket(s.dedupBucket)
	type staleEntry struct {
		key   []byte
		value []byte
	}
	var stale []staleEntry
	if err := bucket.ForEach(func(k, v []byte) error {
		last := bytes.LastIndexByte(v, '\x01')
		if last <= 0 {
			return nil
		}
		expiryMs, err := strconv.ParseInt(string(v[last+1:]), 10, 64)
		if err != nil {
			return nil
		}
		if nowMs >= expiryMs {
			entry := staleEntry{
				key:   append([]byte(nil), k...),
				value: append([]byte(nil), v...),
			}
			stale = append(stale, entry)
		}
		return nil
	}); err != nil {
		logs.Warn("SQS dedup cleanup: scan failed; expired entries stay until the next sweep", logs.Err(err))
	}
	for _, e := range stale {
		// Re-check under the per-key lock, the same mutual exclusion the
		// sends hold: a send that re-armed this exact key after the sweep
		// read it replaced the value, and deleting unconditionally would
		// erase the fresh window and admit a same-key resend as new.
		s.dedupKeyLocker.Lock(string(e.key))
		current, err := bucket.Get(e.key)
		if err != nil || !bytes.Equal(current, e.value) {
			s.dedupKeyLocker.Unlock(string(e.key))
			continue
		}
		if err := bucket.Delete(e.key); err != nil {
			logs.Warn("SQS dedup cleanup: delete failed", logs.String("key", string(e.key)), logs.Err(err))
		}
		s.dedupKeyLocker.Unlock(string(e.key))
	}
}

// doDeletionMarkerCleanup deletes deletion-ledger entries older than the
// recreate-prohibition window. Markers whose queue name is never recreated
// have no other reclamation path (deletedRecently drops a marker only when
// its own queue URL is probed), so this sweep is what bounds the bucket.
// Entries with an unparseable timestamp are removed: they can no longer
// answer the window question either way.
func (s *SQSStore) doDeletionMarkerCleanup() {
	if s.ctx.Err() != nil {
		return
	}
	cutoff := time.Now().Add(-queueDeletionWindow)
	bucket := s.storage.Bucket(s.deletionsBucket)
	var stale [][]byte
	if err := bucket.ForEach(func(k, v []byte) error {
		deletedAt, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil || time.Unix(deletedAt, 0).Before(cutoff) {
			keyCopy := make([]byte, len(k))
			copy(keyCopy, k)
			stale = append(stale, keyCopy)
		}
		return nil
	}); err != nil {
		logs.Warn("SQS deletion-marker cleanup: scan failed; stale markers stay until the next sweep", logs.Err(err))
	}
	for _, k := range stale {
		if err := bucket.Delete(k); err != nil {
			logs.Warn("SQS deletion-marker cleanup: delete failed", logs.String("queueURL", string(k)), logs.Err(err))
		}
	}
}
