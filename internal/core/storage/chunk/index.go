package chunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cockroachdb/pebble/v2"
)

var (
	// ErrIndexNotInitialized is returned when the index has not been initialised.
	ErrIndexNotInitialized = errors.New("index not initialized")
	// ErrInvalidIndexKey is returned when an invalid key is used for index operations.
	ErrInvalidIndexKey = errors.New("invalid index key")
	// ErrChunkNotFound is returned when a requested chunk does not exist in the index.
	ErrChunkNotFound = errors.New("chunk not found")
)

// PebbleIndex provides a Pebble-based index for storing and retrieving chunk metadata.
// It offers fast lookups by chunk ID and time range queries.
type PebbleIndex struct {
	db   *pebble.DB
	path string
	mu   sync.RWMutex
}

// IndexOptions contains configuration options for creating a PebbleIndex.
type IndexOptions struct {
	Path string
}

// NewPebbleIndex creates a new PebbleIndex with the specified options.
// If options are nil or the path is empty, a default path in the data directory is used.
func NewPebbleIndex(opts *IndexOptions) (*PebbleIndex, error) {
	if opts == nil || opts.Path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		opts = &IndexOptions{
			Path: filepath.Join(home, ".vorpalstacks", "data", "chunks", "index.db"),
		}
	}

	db, err := pebble.Open(opts.Path, &pebble.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to open pebble index: %w", err)
	}

	return &PebbleIndex{
		db:   db,
		path: opts.Path,
	}, nil
}

// Meta represents metadata about a stored chunk including its
// time range, entry count, and file location.
type Meta struct {
	ChunkID    string            `json:"chunk_id"`
	MinTs      int64             `json:"min_ts"`
	MaxTs      int64             `json:"max_ts"`
	EntryCount int               `json:"entry_count"`
	ChunkPath  string            `json:"chunk_path"`
	Tags       map[string]string `json:"tags,omitempty"`
}

// Add stores a chunk metadata entry in the index.
func (idx *PebbleIndex) Add(meta Meta) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	batch := idx.db.NewBatch()
	defer batch.Close()

	chunkKey := idx.makeKey(meta)
	value, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal chunk meta: %w", err)
	}

	if err := batch.Set(chunkKey, value, pebble.NoSync); err != nil {
		return err
	}

	tsKey := []byte(fmt.Sprintf("ts:%016X:%s", meta.MinTs, meta.ChunkID))
	if err := batch.Set(tsKey, value, pebble.NoSync); err != nil {
		return err
	}

	for k, v := range meta.Tags {
		tagKey := []byte(fmt.Sprintf("tag:%s:%s:%s", k, v, meta.ChunkID))
		if err := batch.Set(tagKey, value, pebble.NoSync); err != nil {
			return err
		}
	}

	return batch.Commit(pebble.NoSync)
}

// Update updates an existing chunk metadata entry, removing old index keys first.
func (idx *PebbleIndex) Update(meta Meta) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	existing, err := idx.getMetaLocked(meta.ChunkID)
	if err != nil && !errors.Is(err, ErrChunkNotFound) {
		return err
	}

	batch := idx.db.NewBatch()
	defer batch.Close()

	if existing != nil {
		if err := idx.deleteMetaLocked(batch, existing); err != nil {
			return err
		}
	}

	chunkKey := idx.makeKey(meta)
	value, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal chunk meta: %w", err)
	}

	if err := batch.Set(chunkKey, value, pebble.NoSync); err != nil {
		return err
	}

	tsKey := []byte(fmt.Sprintf("ts:%016X:%s", meta.MinTs, meta.ChunkID))
	if err := batch.Set(tsKey, value, pebble.NoSync); err != nil {
		return err
	}

	for k, v := range meta.Tags {
		tagKey := []byte(fmt.Sprintf("tag:%s:%s:%s", k, v, meta.ChunkID))
		if err := batch.Set(tagKey, value, pebble.NoSync); err != nil {
			return err
		}
	}

	return batch.Commit(pebble.NoSync)
}

func (idx *PebbleIndex) getMetaLocked(chunkID string) (*Meta, error) {
	key := []byte("chunk:" + chunkID)
	value, closer, err := idx.db.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, ErrChunkNotFound
		}
		return nil, err
	}
	defer closer.Close()

	var meta Meta
	if err := json.Unmarshal(value, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chunk meta: %w", err)
	}

	return &meta, nil
}

func (idx *PebbleIndex) deleteMetaLocked(batch *pebble.Batch, meta *Meta) error {
	chunkKey := idx.makeKey(*meta)
	if err := batch.Delete(chunkKey, pebble.NoSync); err != nil {
		return err
	}

	tsKey := []byte(fmt.Sprintf("ts:%016X:%s", meta.MinTs, meta.ChunkID))
	if err := batch.Delete(tsKey, pebble.NoSync); err != nil {
		return err
	}

	for k, v := range meta.Tags {
		tagKey := []byte(fmt.Sprintf("tag:%s:%s:%s", k, v, meta.ChunkID))
		if err := batch.Delete(tagKey, pebble.NoSync); err != nil {
			return err
		}
	}

	return nil
}

// AddBatch stores multiple chunk metadata entries in the index in a single transaction.
func (idx *PebbleIndex) AddBatch(metas []Meta) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	batch := idx.db.NewBatch()
	defer batch.Close()

	for _, meta := range metas {
		chunkKey := idx.makeKey(meta)
		value, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("failed to marshal chunk meta: %w", err)
		}
		if err := batch.Set(chunkKey, value, pebble.NoSync); err != nil {
			return err
		}

		tsKey := []byte(fmt.Sprintf("ts:%016X:%s", meta.MinTs, meta.ChunkID))
		if err := batch.Set(tsKey, value, pebble.NoSync); err != nil {
			return err
		}

		for k, v := range meta.Tags {
			tagKey := []byte(fmt.Sprintf("tag:%s:%s:%s", k, v, meta.ChunkID))
			if err := batch.Set(tagKey, value, pebble.NoSync); err != nil {
				return err
			}
		}
	}

	return batch.Commit(pebble.NoSync)
}

// Get retrieves chunk metadata by its unique chunk ID.
func (idx *PebbleIndex) Get(chunkID string) (*Meta, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	key := []byte("chunk:" + chunkID)
	value, closer, err := idx.db.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, ErrChunkNotFound
		}
		return nil, err
	}
	defer closer.Close()

	var meta Meta
	if err := json.Unmarshal(value, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chunk meta: %w", err)
	}

	return &meta, nil
}

// QueryByTimeRange retrieves all chunk metadata entries within the specified timestamp range.
// A chunk is included if its time range overlaps with the query range.
func (idx *PebbleIndex) QueryByTimeRange(minTs, maxTs int64) ([]Meta, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var results []Meta
	seen := make(map[string]bool)

	lowerBound := []byte("ts:")
	upperBound := []byte(fmt.Sprintf("ts:%016X:", maxTs+1))

	iter, err := idx.db.NewIter(&pebble.IterOptions{
		LowerBound: lowerBound,
		UpperBound: upperBound,
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		var meta Meta
		if err := json.Unmarshal(iter.Value(), &meta); err != nil {
			continue
		}
		if meta.MinTs <= maxTs && meta.MaxTs >= minTs {
			if !seen[meta.ChunkID] {
				seen[meta.ChunkID] = true
				results = append(results, meta)
			}
		}
	}

	return results, nil
}

// QueryByTag retrieves all chunk metadata entries with the specified tag key-value pair.
func (idx *PebbleIndex) QueryByTag(key, value string) ([]Meta, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var results []Meta

	prefix := []byte(fmt.Sprintf("tag:%s:%s:", key, value))
	iter, err := idx.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.SeekGE(prefix); iter.Valid(); iter.Next() {
		if !iterHasPrefix(iter.Key(), prefix) {
			break
		}
		var meta Meta
		if err := json.Unmarshal(iter.Value(), &meta); err != nil {
			continue
		}
		results = append(results, meta)
	}

	return results, nil
}

// Delete removes a chunk metadata entry by its chunk ID.
func (idx *PebbleIndex) Delete(chunkID string) error {
	chunkKey := []byte("chunk:" + chunkID)
	value, closer, err := idx.db.Get(chunkKey)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return ErrChunkNotFound
		}
		return err
	}
	defer closer.Close()

	var meta Meta
	if err := json.Unmarshal(value, &meta); err != nil {
		return fmt.Errorf("failed to unmarshal chunk meta: %w", err)
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	batch := idx.db.NewBatch()
	defer batch.Close()

	if err := batch.Delete(chunkKey, pebble.NoSync); err != nil {
		return err
	}

	tsKey := []byte(fmt.Sprintf("ts:%016X:%s", meta.MinTs, chunkID))
	if err := batch.Delete(tsKey, pebble.NoSync); err != nil {
		return err
	}

	for k, v := range meta.Tags {
		tagKey := []byte(fmt.Sprintf("tag:%s:%s:%s", k, v, chunkID))
		if err := batch.Delete(tagKey, pebble.NoSync); err != nil {
			return err
		}
	}

	if err := batch.Commit(pebble.NoSync); err != nil {
		return err
	}

	return nil
}

// DeleteBefore removes all chunk metadata entries with timestamps before the specified value.
// Returns the number of entries deleted.
func (idx *PebbleIndex) DeleteBefore(ts int64) (int, error) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	var count int
	var toDelete []Meta

	iter, err := idx.db.NewIter(&pebble.IterOptions{
		LowerBound: []byte("chunk:"),
		UpperBound: []byte("chunk~"),
	})
	if err != nil {
		return 0, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		var meta Meta
		if err := json.Unmarshal(iter.Value(), &meta); err != nil {
			continue
		}
		if meta.MaxTs < ts {
			toDelete = append(toDelete, meta)
		}
	}

	if len(toDelete) == 0 {
		return 0, nil
	}

	for _, meta := range toDelete {
		if err := idx.deleteMeta(&meta); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (idx *PebbleIndex) deleteMeta(meta *Meta) error {
	batch := idx.db.NewBatch()
	defer batch.Close()

	chunkKey := idx.makeKey(*meta)
	if err := batch.Delete(chunkKey, pebble.NoSync); err != nil {
		return err
	}

	tsKey := []byte(fmt.Sprintf("ts:%016X:%s", meta.MinTs, meta.ChunkID))
	if err := batch.Delete(tsKey, pebble.NoSync); err != nil {
		return err
	}

	for k, v := range meta.Tags {
		tagKey := []byte(fmt.Sprintf("tag:%s:%s:%s", k, v, meta.ChunkID))
		if err := batch.Delete(tagKey, pebble.NoSync); err != nil {
			return err
		}
	}

	return batch.Commit(pebble.NoSync)
}

// Count returns the total number of entries in the index.
func (idx *PebbleIndex) Count() (int, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var count int
	iter, err := idx.db.NewIter(&pebble.IterOptions{
		LowerBound: []byte("chunk:"),
		UpperBound: []byte("chunk~"),
	})
	if err != nil {
		return 0, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		count++
	}

	return count, nil
}

// Close closes the Pebble index and releases all resources.
func (idx *PebbleIndex) Close() error {
	return idx.db.Close()
}

func (idx *PebbleIndex) makeKey(meta Meta) []byte {
	return []byte("chunk:" + meta.ChunkID)
}

func iterHasPrefix(key, prefix []byte) bool {
	if len(key) < len(prefix) {
		return false
	}
	for i := range prefix {
		if key[i] != prefix[i] {
			return false
		}
	}
	return true
}
