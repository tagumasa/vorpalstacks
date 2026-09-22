package cloudwatchlogs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
)

// The transformer record families: one transformer per log group (the
// ordered processor list) and the per-event transformed messages the
// ingestion seam records and the query plane reads. Both persist in the
// JSON regime; the transformed-message family rides the group's
// lifecycle (the delete-time teardown and the retention purge remove it
// with the group's chunks).

// --- Transformer ---

// transformerRecordMu makes the transformer record's read-modify-write
// and its delete one critical section each, so a Put that read the
// record before a delete cannot write it back afterwards.
var transformerRecordMu sync.Mutex

// PutTransformer stores (or replaces) a log group's transformer,
// preserving the creation stamp across updates and refreshing the
// modification stamp. The write holds the group's write lock — the
// DeleteLogGroup teardown lock — and re-checks the group inside it: the
// service layer's earlier group check runs outside any critical section,
// so without the re-check a teardown interleaving this put would leave
// the transformer permanently outliving its deleted group (the
// metric-filter family's pattern; the record mutex nests inside in the
// same order the teardown's DeleteTransformer takes it in).
func (s *Store) PutTransformer(t *Transformer) error {
	gLock := s.groupLock(t.LogGroupName)
	gLock.Lock()
	defer gLock.Unlock()

	if _, err := s.GetLogGroup(t.LogGroupName); err != nil {
		return err
	}
	transformerRecordMu.Lock()
	defer transformerRecordMu.Unlock()
	now := time.Now().UTC().UnixMilli()
	if existing, err := s.GetTransformer(t.LogGroupName); err == nil {
		t.CreationTime = existing.CreationTime
	} else {
		t.CreationTime = now
	}
	t.LastModifiedTime = now
	return s.Put(s.transformerKey(t.LogGroupName), t)
}

func (s *Store) GetTransformer(logGroupName string) (*Transformer, error) {
	return getJSONRecord[Transformer](s, s.transformerKey(logGroupName), ErrResourceNotFound)
}

func (s *Store) DeleteTransformer(logGroupName string) error {
	transformerRecordMu.Lock()
	defer transformerRecordMu.Unlock()
	if !s.Exists(s.transformerKey(logGroupName)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.transformerKey(logGroupName))
}

// --- Transformed messages ---

// TransformedMessageDigest derives one event's transformed-message key
// component from the identity the chunk records carry. The stream length
// prefixes the digest input: '|' is a legal character in both stream
// names and messages, so a bare separator alone lets two different
// (stream, message) pairs collide on one digest and overwrite each
// other's transformed copies.
func TransformedMessageDigest(timestamp int64, logStream, message string) string {
	sum := sha256.Sum256([]byte(strconv.FormatInt(timestamp, 10) + "|" +
		strconv.Itoa(len(logStream)) + "|" + logStream + "|" + message))
	return hex.EncodeToString(sum[:12])
}

// PutTransformedMessages records the batch's transformed messages
// (digest → transformed message). Entries with an empty message are
// skipped: an event the recipe left without added keys has no
// transformed copy. The write and the group-existence check run under
// the group's write lock — the same lock the whole teardown holds — so
// a batch transformed while its group is being deleted either completes
// before the teardown (and its records die with the group) or observes
// the removed group and writes nothing (a silent skip: the ingestion
// itself already succeeded for the group that no longer exists).
func (s *Store) PutTransformedMessages(logGroupName string, entries map[string]TransformedMessage) error {
	gLock := s.groupLock(logGroupName)
	gLock.Lock()
	defer gLock.Unlock()
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		return nil
	}
	for digest, entry := range entries {
		if entry.Message == "" {
			continue
		}
		if err := s.Put(s.transformedMessageKey(logGroupName, digest), &entry); err != nil {
			return err
		}
	}
	return nil
}

// TransformedMessageMap loads a group's transformed messages keyed by
// digest, the query plane's join form. A decode failure logs and skips
// the record: corruption must stay visible without shrinking the rest
// of the group's transformed set.
func (s *Store) TransformedMessageMap(logGroupName string) (map[string]string, error) {
	result := make(map[string]string)
	prefix := keyPrefixTransformedMessage + escapePath(logGroupName) + ":"
	err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var entry TransformedMessage
		if err := json.Unmarshal(value, &entry); err != nil {
			logs.Warn("Corrupt transformed message record",
				logs.String("key", key), logs.Err(err))
			return nil
		}
		digest := key[len(prefix):]
		result[digest] = entry.Message
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteTransformedMessages removes every transformed-message record of
// a group (the delete-time teardown path).
func (s *Store) DeleteTransformedMessages(logGroupName string) error {
	prefix := keyPrefixTransformedMessage + escapePath(logGroupName) + ":"
	var keys []string
	if err := s.ScanPrefix(prefix, func(key string, _ []byte) error {
		keys = append(keys, key)
		return nil
	}); err != nil {
		return err
	}
	for _, key := range keys {
		if err := s.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// PurgeExpiredTransformedMessages removes a group's transformed-message
// records whose event timestamp fell past the retention cutoff (the
// retention sweep's companion to PurgeExpiredChunks).
func (s *Store) PurgeExpiredTransformedMessages(logGroupName string, cutoff int64) {
	prefix := keyPrefixTransformedMessage + escapePath(logGroupName) + ":"
	var expired []string
	_ = s.ScanPrefix(prefix, func(key string, value []byte) error {
		var entry TransformedMessage
		if err := json.Unmarshal(value, &entry); err != nil {
			return nil
		}
		if entry.Timestamp < cutoff {
			expired = append(expired, key)
		}
		return nil
	})
	for _, key := range expired {
		_ = s.Delete(key)
	}
}

// Transformer represents one log group's transformer: the ordered
// processor list (stored verbatim as the model's Processors vocabulary)
// plus the modification stamps the GetTransformer response carries.
type Transformer struct {
	LogGroupName       string                   `json:"logGroupName"`
	LogGroupIdentifier string                   `json:"logGroupIdentifier"`
	Config             []map[string]interface{} `json:"config"`
	CreationTime       int64                    `json:"creationTime"`
	LastModifiedTime   int64                    `json:"lastModifiedTime"`
}

// TransformedMessage records one event's transformed message — the
// ingestion-time transformation output the query plane reads (originals
// stay the GetLogEvents/FilterLogEvents surface). Keyed by the digest of
// (timestamp, stream, original message).
type TransformedMessage struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}
