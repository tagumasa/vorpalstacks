package sqs

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
)

// deduplicationCacheMaxSize bounds the in-memory deduplication cache: beyond
// it the earliest-expiring entries are evicted (expired first, then the
// soonest-to-expire live ones). The persisted bucket remains authoritative —
// an evicted live key is re-read from it on the next lookup — so the bound
// costs one bucket read, never correctness.
const deduplicationCacheMaxSize = 500

type deduplicationEntry struct {
	messageID      string
	sequenceNumber string
	expiresAt      time.Time
}

func (s *SQSStore) buildDeduplicationKey(queueURL string, message *Message) string {
	if message.MessageDeduplicationID != "" {
		return queueURL + "#" + message.MessageDeduplicationID
	}
	return queueURL + "#" + calculateMD5(message.Body)
}

func (s *SQSStore) getDeduplicationMessageID(dedupKey string) (messageID, sequenceNumber string, ok bool) {
	s.deduplicationMu.RLock()
	entry, exists := s.deduplicationCache[dedupKey]
	s.deduplicationMu.RUnlock()

	if exists && time.Now().Before(entry.expiresAt) {
		return entry.messageID, entry.sequenceNumber, true
	}

	data, err := s.storage.Bucket(s.dedupBucket).Get([]byte(dedupKey))
	if err != nil {
		// Fail-open by design (a storage blip must not block sends), but
		// never silently: the miss admits a same-key resend as new.
		logs.Warn("SQS: dedup bucket read failed; the window check missed", logs.String("key", dedupKey), logs.Err(err))
	} else if data != nil {
		// The persisted form is "<messageKey>\x01<sequenceNumber>\x01<expiryMs>"
		// (see deduplicationEntryValue): the expiry is the final segment.
		first := bytes.IndexByte(data, '\x01')
		last := bytes.LastIndexByte(data, '\x01')
		if first > 0 && last > first {
			msgKey := string(data[:first])
			seq := string(data[first+1 : last])
			expiryMs, err := strconv.ParseInt(string(data[last+1:]), 10, 64)
			if err == nil {
				if time.Now().UnixMilli() >= expiryMs {
					_ = s.storage.Bucket(s.dedupBucket).Delete([]byte(dedupKey))
					s.deduplicationMu.Lock()
					delete(s.deduplicationCache, dedupKey)
					s.deduplicationMu.Unlock()
					return "", "", false
				}
				s.deduplicationMu.Lock()
				s.deduplicationCache[dedupKey] = &deduplicationEntry{
					messageID:      msgKey,
					sequenceNumber: seq,
					// Arm the in-memory TTL from the persisted expiry, not
					// from now: the interval starts at the first send, so a
					// repeatedly re-read key must still expire on schedule.
					expiresAt: time.UnixMilli(expiryMs),
				}
				// The re-arm insert honours the same bound as the register
				// path: without it a run of re-read keys grows the cache past
				// the cap (the persisted bucket stays authoritative; the
				// eviction trades one read, never a lost window).
				if len(s.deduplicationCache) > deduplicationCacheMaxSize {
					s.cleanupDeduplicationCache()
				}
				s.deduplicationMu.Unlock()
				return msgKey, seq, true
			}
		}
	}

	return "", "", false
}

// deduplicationEntryValue encodes the persisted form of a deduplication
// entry: the recorded message key, its sequence number (the synthetic
// acknowledgment of a resend after the original's deletion echoes it), a
// \x01 separator each, and the window's expiry as unix milliseconds — the
// final segment, so both readers parse the expiry from the tail. Callers
// commit it inside the same transaction as the message persist, so a
// storage failure can never queue a message whose deduplication window is
// silently missing.
func deduplicationEntryValue(messageID, sequenceNumber string) []byte {
	expiry := time.Now().Add(deduplicationWindow).UnixMilli()
	val := append([]byte(messageID), '\x01')
	val = append(val, []byte(sequenceNumber)...)
	val = append(val, '\x01')
	return append(val, []byte(strconv.FormatInt(expiry, 10))...)
}

// registerDeduplicationEntry publishes the in-memory view of a
// deduplication entry whose persisted form the caller has already committed
// (inside the message persist's transaction).
func (s *SQSStore) registerDeduplicationEntry(dedupKey, messageID, sequenceNumber string) {
	s.deduplicationMu.Lock()
	s.deduplicationCache[dedupKey] = &deduplicationEntry{
		messageID:      messageID,
		sequenceNumber: sequenceNumber,
		expiresAt:      time.Now().Add(deduplicationWindow),
	}
	if len(s.deduplicationCache) > deduplicationCacheMaxSize {
		s.cleanupDeduplicationCache()
	}
	s.deduplicationMu.Unlock()
}

func (s *SQSStore) cleanupDeduplicationCache() {
	now := time.Now()
	deleted := 0
	const maxDeletesPerCleanup = 100
	for key, entry := range s.deduplicationCache {
		if now.After(entry.expiresAt) {
			delete(s.deduplicationCache, key)
			deleted++
			if deleted >= maxDeletesPerCleanup {
				break
			}
		}
	}
	// The bound is real: with more live entries than the cap, the
	// earliest-expiring ones go — the persisted bucket answers their next
	// lookup, so eviction trades one read, never a lost window.
	for len(s.deduplicationCache) > deduplicationCacheMaxSize {
		var oldestKey string
		var oldestExpiry time.Time
		first := true
		for key, entry := range s.deduplicationCache {
			if first || entry.expiresAt.Before(oldestExpiry) {
				oldestKey, oldestExpiry, first = key, entry.expiresAt, false
			}
		}
		delete(s.deduplicationCache, oldestKey)
	}
}

func (s *SQSStore) cleanupDeduplicationCacheForQueue(queueURL string) {
	prefix := queueURL + "#"
	for key := range s.deduplicationCache {
		if strings.HasPrefix(key, prefix) {
			delete(s.deduplicationCache, key)
		}
	}
}
