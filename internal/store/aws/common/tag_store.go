package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
)

const indexSeparator = "\x00"

// MaxTagsPerResource is the AWS-wide per-resource tag total most services
// document ("at most 50 user-defined tags per resource"); stores whose
// service documents the standard total pass it in their TagBudget.
// A service documenting a different total passes its own constant instead —
// Route 53's API reference, for one, caps a hosted zone or health check at
// ten tags.
const MaxTagsPerResource = types.MaxTagsPerResource

// StandardTagBudget builds the budget of a service documenting the
// AWS-wide fifty, answering an overflow with a wire error of the given
// code — the service's own invalid-input identity, so the error needs no
// per-path mapping on its way to the transport.
func StandardTagBudget(code string) TagBudget {
	return TagBudget{
		MaxKeys:  MaxTagsPerResource,
		Exceeded: NewAWSError(code, "Too many tags.", http.StatusBadRequest),
	}
}

// ErrTagQuotaExceeded reports a tag write whose merged result would carry
// more keys than the calling service's documented per-resource limit; the
// check runs under the store's lock, so the bound holds against concurrent
// writers. A TagBudget whose Exceeded member is set replaces this sentinel
// with the service's own wire identity.
var ErrTagQuotaExceeded = fmt.Errorf("tagstore: merged tag count exceeds the resource's tag limit")

// TagBudget states a service's documented per-resource tag total and the
// wire identity a write past that total answers with.
type TagBudget struct {
	// MaxKeys caps the merged tag count of one resource; zero leaves the
	// merge path unbounded, for services that document no total.
	MaxKeys int
	// Exceeded is returned when a write passes MaxKeys — the service's own
	// invalid-input identity, injected at construction so every write path
	// (tag handler, resource creation, any future caller) answers with it
	// without per-path mapping. Nil falls back to ErrTagQuotaExceeded.
	Exceeded error
}

// TagStore manages resource tags using a Loki-style inverted index.
// The main store holds tag maps keyed by resource key, while the index store
// enables efficient lookup of resources by tag key or key-value pair.
type TagStore struct {
	main  *BaseStore
	index *BaseStore
	// mainName/indexName resolve the same buckets inside a caller's
	// transaction for TagInTxn.
	mainName  string
	indexName string
	// budget is the service's documented tag total and overflow identity,
	// set at construction. Read under mu.
	budget TagBudget
	mu     sync.Mutex
}

// NewTagStore creates a TagStore with region-agnostic bucket names and the
// service's documented tag budget.
func NewTagStore(store storage.BasicStorage, serviceName string, budget TagBudget) *TagStore {
	mainName := serviceName + "-tags"
	indexName := serviceName + "-tag-idx"
	return &TagStore{
		main:      NewBaseStore(store.Bucket(mainName), mainName),
		index:     NewBaseStore(store.Bucket(indexName), indexName),
		mainName:  mainName,
		indexName: indexName,
		budget:    budget,
	}
}

// NewTagStoreWithRegion creates a TagStore with region-scoped bucket names
// and the service's documented tag budget.
func NewTagStoreWithRegion(store storage.BasicStorage, serviceName, region string, budget TagBudget) *TagStore {
	mainName := serviceName + "-tags-" + region
	indexName := serviceName + "-tag-idx-" + region
	return &TagStore{
		main:      NewBaseStore(store.Bucket(mainName), mainName),
		index:     NewBaseStore(store.Bucket(indexName), indexName),
		mainName:  mainName,
		indexName: indexName,
		budget:    budget,
	}
}

// TagInTxn stages a FRESH resource's initial tag set on the caller's
// transaction: the main entry plus the inverted-index segments for every
// key. It exists so a resource's creation and its initial tags commit as
// one unit — a tag failure then leaves no half-created resource for a
// retry to collide with. The store's own mutex is deliberately not taken:
// the calling path holds its resource-level lock and the resource does not
// exist yet, so no concurrent tagger can interleave; stale index cleanup
// is unnecessary for the same reason (a fresh key has no prior segments).
// The merged-count bound TagWithLimit enforces applies to the initial set
// alone here — a fresh resource has nothing to merge with.
func (t *TagStore) TagInTxn(txn storage.Transaction, resourceKey string, tags map[string]string, maxKeys int) error {
	if len(tags) == 0 {
		return nil
	}
	if err := t.ValidateTags(tags); err != nil {
		return err
	}
	if len(tags) > maxKeys {
		return ErrTagQuotaExceeded
	}
	// The main entry is JSON-encoded, the same form BaseStore.Put writes,
	// so List decodes it identically.
	entryBytes, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	if err := txn.Bucket(t.mainName).Put([]byte(resourceKey), entryBytes); err != nil {
		return err
	}
	idxBucket := txn.Bucket(t.indexName)
	for k, v := range tags {
		idxKey := k + "=" + v + indexSeparator + resourceKey
		if err := idxBucket.Put([]byte(idxKey), []byte{0x01}); err != nil {
			return err
		}
	}
	return nil
}

// List returns all tags for the given resource as a key-value map.
// Returns an empty map if the resource has no tags.
func (t *TagStore) List(resourceKey string) (map[string]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var tags map[string]string
	if err := t.main.Get(resourceKey, &tags); err != nil {
		if IsNotFound(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	return tags, nil
}

// ListAsSlice returns all tags for the given resource as a slice of Tag structs.
func (t *TagStore) ListAsSlice(resourceKey string) ([]types.Tag, error) {
	tags, err := t.List(resourceKey)
	if err != nil {
		return nil, err
	}
	result := make([]types.Tag, 0, len(tags))
	for k, v := range tags {
		result = append(result, types.Tag{Key: k, Value: v})
	}
	return result, nil
}

// Tag merges the given tags into the existing tags for a resource.
// Existing tag keys are overwritten; new tag keys are added.
// For overwritten keys, stale index entries are removed before inserting new ones.
//
// ValidateTags performs the deterministic validation applied by Tag
// without writing anything, so callers can reject invalid input before
// mutating any persisted state.
func (t *TagStore) ValidateTags(newTags map[string]string) error {
	for k, v := range newTags {
		if strings.ContainsRune(k, '\x00') || strings.ContainsRune(v, '\x00') {
			return fmt.Errorf("tagstore: tag key or value contains forbidden NUL byte (key=%q)", k)
		}
	}
	return nil
}

// Tag keys and values MUST NOT contain the NUL byte (\x00). The byte is
// used as the index-segment separator and would corrupt resource-key
// extraction in FindByTag. AWS itself rejects NUL in tag values, so
// this guard primarily defends internal consistency.
func (t *TagStore) Tag(resourceKey string, newTags map[string]string) error {
	if len(newTags) == 0 {
		return nil
	}
	if err := t.ValidateTags(newTags); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	exceeded, err := t.enforceLimitLocked(resourceKey, newTags, t.budget.MaxKeys)
	if err != nil {
		return err
	}
	if exceeded {
		return t.quotaExceeded()
	}
	return t.mergeLocked(resourceKey, newTags)
}

// enforceLimitLocked reports whether merging newTags onto the resource's
// live set would carry more keys than maxKeys: existing keys the write
// overwrites cost nothing. maxKeys zero or below disables the check. The
// caller must hold t.mu.
func (t *TagStore) enforceLimitLocked(resourceKey string, newTags map[string]string, maxKeys int) (bool, error) {
	if maxKeys <= 0 {
		return false, nil
	}
	existing, err := t.listUnlocked(resourceKey)
	if err != nil {
		return false, err
	}
	merged := len(existing)
	for k := range newTags {
		if _, ok := existing[k]; !ok {
			merged++
		}
	}
	return merged > maxKeys, nil
}

// quotaExceeded is the error a write past the documented total answers
// with: the service-injected identity from the construction budget, or the
// generic sentinel when the budget carries none.
func (t *TagStore) quotaExceeded() error {
	if t.budget.Exceeded != nil {
		return t.budget.Exceeded
	}
	return ErrTagQuotaExceeded
}

// mergeLocked merges newTags into the resource's live set and rewrites the
// main entry plus the touched index segments; the caller must hold t.mu.
func (t *TagStore) mergeLocked(resourceKey string, newTags map[string]string) error {
	existing, err := t.listUnlocked(resourceKey)
	if err != nil {
		return err
	}
	for k, v := range newTags {
		existing[k] = v
	}
	if err := t.main.Put(resourceKey, existing); err != nil {
		return err
	}

	suffix := indexSeparator + resourceKey
	for k, v := range newTags {
		if err := t.deleteIndexEntriesUnlocked(k, suffix); err != nil {
			return err
		}
		idxKey := k + "=" + v + indexSeparator + resourceKey
		if err := t.index.Put(idxKey, []byte{0x01}); err != nil {
			return err
		}
	}
	return nil
}

// Replace swaps the full tag set of a resource: keys absent from the new set
// are removed, present keys are set to the new values, and the main entry and
// the inverted index are rewritten under one lock with a single main write.
// Callers whose input is the desired end state use this instead of an
// Untag/Tag pair, which needs a rollback path when the second write fails.
// The input set is itself bounded by the store's documented cap: an
// end-state write past the cap is a quota violation like any merge.
func (t *TagStore) Replace(resourceKey string, tags map[string]string) error {
	if err := t.ValidateTags(tags); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.budget.MaxKeys > 0 && len(tags) > t.budget.MaxKeys {
		return t.quotaExceeded()
	}

	existing, err := t.listUnlocked(resourceKey)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		if err := t.main.Delete(resourceKey); err != nil {
			return err
		}
	} else if err := t.main.Put(resourceKey, tags); err != nil {
		return err
	}

	suffix := indexSeparator + resourceKey
	// Drop the index entries of the old set, then recreate the entries of the
	// new set, so removed keys and changed values both leave no stale entries.
	for k := range existing {
		if err := t.deleteIndexEntriesUnlocked(k, suffix); err != nil {
			return err
		}
	}
	for k, v := range tags {
		idxKey := k + "=" + v + indexSeparator + resourceKey
		if err := t.index.Put(idxKey, []byte{0x01}); err != nil {
			return err
		}
	}
	return nil
}

// TagFromSlice converts a slice of Tag structs to a map and delegates to Tag.
func (t *TagStore) TagFromSlice(resourceKey string, tags []types.Tag) error {
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return t.Tag(resourceKey, tagMap)
}

// Untag removes the specified tag keys from a resource.
// If all tags are removed, the resource entry is deleted from the main store.
// The corresponding inverted index entries are also removed.
func (t *TagStore) Untag(resourceKey string, tagKeys []string) error {
	if len(tagKeys) == 0 {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	tags, err := t.listUnlocked(resourceKey)
	if err != nil {
		return err
	}

	removeSet := make(map[string]bool, len(tagKeys))
	for _, k := range tagKeys {
		delete(tags, k)
		removeSet[k] = true
	}

	if len(tags) == 0 {
		if err := t.main.Delete(resourceKey); err != nil {
			return err
		}
	} else {
		if err := t.main.Put(resourceKey, tags); err != nil {
			return err
		}
	}

	suffix := indexSeparator + resourceKey
	for k := range removeSet {
		if err := t.deleteIndexEntriesUnlocked(k, suffix); err != nil {
			return err
		}
	}
	return nil
}

// deleteIndexEntriesUnlocked removes the index entries of one tag key for the
// resource identified by suffix; the caller must hold the store lock.
func (t *TagStore) deleteIndexEntriesUnlocked(tagKey, suffix string) error {
	prefix := tagKey + "="
	return t.index.ScanPrefix(prefix, func(idxKey string, _ []byte) error {
		if strings.HasSuffix(idxKey, suffix) {
			return t.index.Delete(idxKey)
		}
		return nil
	})
}

// Delete removes all tags and index entries for a resource.
func (t *TagStore) Delete(resourceKey string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.main.Delete(resourceKey); err != nil {
		return err
	}

	suffix := indexSeparator + resourceKey
	return t.index.ScanPrefix("", func(idxKey string, _ []byte) error {
		if strings.HasSuffix(idxKey, suffix) {
			return t.index.Delete(idxKey)
		}
		return nil
	})
}

// FindByTag returns all resource keys that have a tag with the given key (any value).
func (t *TagStore) FindByTag(tagKey string) ([]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	prefix := tagKey + "="
	result := make([]string, 0)
	err := t.index.ScanPrefix(prefix, func(idxKey string, _ []byte) error {
		idx := bytes.Index([]byte(idxKey), []byte(indexSeparator))
		if idx >= 0 {
			result = append(result, idxKey[idx+1:])
		}
		return nil
	})
	return result, err
}

// FindByTagValue returns all resource keys that have a tag matching both the given key and value.
func (t *TagStore) FindByTagValue(tagKey, tagValue string) ([]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	prefix := tagKey + "=" + tagValue + indexSeparator
	result := make([]string, 0)
	err := t.index.ScanPrefix(prefix, func(idxKey string, _ []byte) error {
		result = append(result, idxKey[len(prefix):])
		return nil
	})
	return result, err
}

// RebuildIndex deletes and recreates the inverted index from the main tag store.
// Use this after bulk tag migrations or data corruption recovery.
func (t *TagStore) RebuildIndex() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.index.DeleteByPrefix(""); err != nil {
		return err
	}

	return t.main.ForEach(func(resourceKey string, value []byte) error {
		var tags map[string]string
		if err := t.main.Get(resourceKey, &tags); err != nil {
			return err
		}
		for k, v := range tags {
			idxKey := k + "=" + v + indexSeparator + resourceKey
			if err := t.index.Put(idxKey, []byte{0x01}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *TagStore) listUnlocked(resourceKey string) (map[string]string, error) {
	var tags map[string]string
	if err := t.main.Get(resourceKey, &tags); err != nil {
		if IsNotFound(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	return tags, nil
}
