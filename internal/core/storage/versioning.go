package storage

import (
	"encoding/binary"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage/pebbledb"

	"github.com/cockroachdb/pebble/v2"
)

const uint64Size = 8

// PebbleVersionedBucket provides a versioned key-value bucket implementation
// that automatically manages version numbers for each key.
type PebbleVersionedBucket struct {
	db     *pebbledb.DB
	prefix []byte
	mu     sync.Mutex
}

// NewPebbleVersionedBucket creates a new versioned bucket backed by Pebble.
func NewPebbleVersionedBucket(db *pebbledb.DB, prefix []byte) *PebbleVersionedBucket {
	b := &PebbleVersionedBucket{
		db:     db,
		prefix: prefix,
	}
	return b
}

func (b *PebbleVersionedBucket) makeKey(key []byte, version uint64) []byte {
	k := make([]byte, len(b.prefix)+len(key)+1+8)
	copy(k, b.prefix)
	copy(k[len(b.prefix):], key)
	k[len(b.prefix)+len(key)] = '#'
	binary.BigEndian.PutUint64(k[len(b.prefix)+len(key)+1:], version)
	return k
}

func (b *PebbleVersionedBucket) makeMetaKey(key []byte) []byte {
	k := make([]byte, len(b.prefix)+len(MetaPrefix)+len(key))
	copy(k, b.prefix)
	copy(k[len(b.prefix):], []byte(MetaPrefix))
	copy(k[len(b.prefix)+len(MetaPrefix):], key)
	return k
}

func (b *PebbleVersionedBucket) makeCounterKey() []byte {
	k := make([]byte, len(b.prefix)+len(CounterPrefix))
	copy(k, b.prefix)
	copy(k[len(b.prefix):], []byte(CounterPrefix))
	return k
}

func (b *PebbleVersionedBucket) getNextVersion() (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	counterKey := b.makeCounterKey()

	var currentVersion uint64
	data, err := b.db.Get(counterKey)
	if err != nil && err != pebbledb.ErrKeyNotFound {
		return 0, err
	}
	if err == nil && len(data) >= 8 {
		currentVersion = binary.BigEndian.Uint64(data)
	}

	newVersion := currentVersion + 1
	counterData := make([]byte, uint64Size)
	binary.BigEndian.PutUint64(counterData, newVersion)

	batch := b.db.NewBatch()
	if err := batch.Set(counterKey, counterData); err != nil {
		batch.Close()
		return 0, err
	}
	if err := batch.Commit(pebble.Sync); err != nil {
		batch.Close()
		return 0, err
	}
	batch.Close()

	return newVersion, nil
}

// GetVersion retrieves a specific version of a key.
func (b *PebbleVersionedBucket) GetVersion(key []byte, version uint64) (*VersionedValue, error) {
	fullKey := b.makeKey(key, version)
	data, err := b.db.Get(fullKey)
	if err != nil {
		if err == pebbledb.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return b.decodeVersionedValue(data), nil
}

// GetLatest retrieves the latest version of a key.
func (b *PebbleVersionedBucket) GetLatest(key []byte) (*VersionedValue, error) {
	metaKey := b.makeMetaKey(key)
	metaData, err := b.db.Get(metaKey)
	if err != nil {
		if err == pebbledb.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if len(metaData) < 8 {
		return nil, nil
	}
	latestVersion := binary.BigEndian.Uint64(metaData[:8])

	return b.GetVersion(key, latestVersion)
}

// PutWithVersion stores a value with automatic version incrementing.
// Returns the new versioned value.
func (b *PebbleVersionedBucket) PutWithVersion(key, value []byte) (*VersionedValue, error) {
	version, err := b.getNextVersion()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()

	vv := &VersionedValue{
		Value:     value,
		Version:   version,
		Timestamp: now,
		Deleted:   false,
	}

	encoded := b.encodeVersionedValue(vv)
	fullKey := b.makeKey(key, version)

	batch := b.db.NewBatch()
	defer batch.Close()

	if err := batch.Set(fullKey, encoded); err != nil {
		return nil, err
	}

	metaKey := b.makeMetaKey(key)
	metaData := make([]byte, uint64Size)
	binary.BigEndian.PutUint64(metaData, version)
	if err := batch.Set(metaKey, metaData); err != nil {
		return nil, err
	}

	if err := batch.Commit(pebble.Sync); err != nil {
		return nil, err
	}

	return vv, nil
}

// DeleteWithVersion marks a key as deleted with a new version.
// This creates a tombstone that can be filtered in list operations.
func (b *PebbleVersionedBucket) DeleteWithVersion(key []byte) (*VersionedValue, error) {
	version, err := b.getNextVersion()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()

	vv := &VersionedValue{
		Value:     nil,
		Version:   version,
		Timestamp: now,
		Deleted:   true,
	}

	encoded := b.encodeVersionedValue(vv)
	fullKey := b.makeKey(key, version)

	batch := b.db.NewBatch()
	defer batch.Close()

	if err := batch.Set(fullKey, encoded); err != nil {
		return nil, err
	}

	metaKey := b.makeMetaKey(key)
	metaData := make([]byte, uint64Size)
	binary.BigEndian.PutUint64(metaData, version)
	if err := batch.Set(metaKey, metaData); err != nil {
		return nil, err
	}

	if err := batch.Commit(pebble.Sync); err != nil {
		return nil, err
	}

	return vv, nil
}

// ListVersions returns all versions of a key, optionally filtered and sorted.
func (b *PebbleVersionedBucket) ListVersions(key []byte, opts VersionListOptions) ([]*VersionedValue, error) {
	prefix := make([]byte, len(b.prefix)+len(key)+1)
	copy(prefix, b.prefix)
	copy(prefix[len(b.prefix):], key)
	prefix[len(b.prefix)+len(key)] = '#'

	upperBound := make([]byte, len(prefix)+1)
	copy(upperBound, prefix)
	upperBound[len(prefix)] = 0xFF

	iterOpts := &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound,
	}

	iter, err := b.db.NewIter(iterOpts)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var versions []*VersionedValue
	count := 0

	moveFn := iter.Next
	if opts.Reverse {
		moveFn = iter.Prev
		if !iter.SeekLT(upperBound) {
			iter.First()
		}
	} else {
		if !iter.SeekGE(prefix) {
			return versions, nil
		}
	}

	for ; iter.Valid(); moveFn() {
		if opts.Limit > 0 && count >= opts.Limit {
			break
		}

		data, err := iter.ValueAndErr()
		if err != nil {
			continue
		}
		vv := b.decodeVersionedValue(data)

		if opts.StartVersion > 0 && vv.Version < opts.StartVersion {
			if opts.Reverse {
				continue
			}
			break
		}

		if !vv.Deleted || opts.IncludeDeleted {
			versions = append(versions, vv)
			count++
		}
	}

	return versions, iter.Error()
}

// PurgeVersions removes old versions of a key, keeping only the specified number.
func (b *PebbleVersionedBucket) PurgeVersions(key []byte, keepVersions int) error {
	versions, err := b.ListVersions(key, VersionListOptions{IncludeDeleted: true, Reverse: true})
	if err != nil {
		return err
	}

	if len(versions) <= keepVersions {
		return nil
	}

	toDelete := versions[keepVersions:]
	batch := b.db.NewBatch()
	defer batch.Close()

	for _, vv := range toDelete {
		fullKey := b.makeKey(key, vv.Version)
		if err := batch.Delete(fullKey); err != nil {
			return err
		}
	}

	return batch.Commit(&pebble.WriteOptions{Sync: true})
}

func (b *PebbleVersionedBucket) encodeVersionedValue(vv *VersionedValue) []byte {
	size := 1 + 8 + 8 + 4 + len(vv.Value)
	buf := make([]byte, size)
	offset := 0

	if vv.Deleted {
		buf[offset] = 1
	} else {
		buf[offset] = 0
	}
	offset++

	binary.BigEndian.PutUint64(buf[offset:], vv.Version)
	offset += 8

	binary.BigEndian.PutUint64(buf[offset:], uint64(vv.Timestamp))
	offset += 8

	binary.BigEndian.PutUint32(buf[offset:], uint32(len(vv.Value)))
	offset += 4

	copy(buf[offset:], vv.Value)

	return buf
}

func (b *PebbleVersionedBucket) decodeVersionedValue(data []byte) *VersionedValue {
	if len(data) < 21 {
		return &VersionedValue{}
	}

	vv := &VersionedValue{}
	offset := 0

	vv.Deleted = data[offset] == 1
	offset++

	vv.Version = binary.BigEndian.Uint64(data[offset:])
	offset += 8

	if offset+8 > len(data) {
		return vv
	}
	vv.Timestamp = int64(binary.BigEndian.Uint64(data[offset:]))
	offset += 8

	valueLen := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	if int(valueLen) > 0 && offset+int(valueLen) <= len(data) {
		vv.Value = make([]byte, valueLen)
		copy(vv.Value, data[offset:offset+int(valueLen)])
	}

	return vv
}

// VersionConflictError is returned when a version conflict is detected
// during optimistic locking operations.
type VersionConflictError struct {
	Key []byte
	Err error
}

// Error returns the error message for the version conflict.
func (e *VersionConflictError) Error() string {
	if e.Err != nil {
		return "version conflict: " + e.Err.Error()
	}
	return "version conflict: failed to acquire unique version after retries"
}

// Unwrap returns the underlying error, if any.
func (e *VersionConflictError) Unwrap() error { return e.Err }
