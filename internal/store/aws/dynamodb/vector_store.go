package dynamodb

import (
	"fmt"
	"sort"
	"strconv"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/utils/vecdist"
)

// Vector index entries record one vector per (index, item) pair in a dedicated
// bucket. The key is table\x00index\x00primaryKey, where primaryKey is the
// same encoded item key the GSI/LSI entry values carry, so a search resolves
// matching items without a second lookup index; the value is the proto
// VectorIndexEntry holding the vector itself. All maintenance runs inside the
// caller's transaction, mirroring IndexStore, and is reached through
// DynamoDBTxn.PutIndexEntries/DeleteIndexEntries plus the vector-specific
// methods below.

// extractVector reads the index's vector attribute from an item as a float64
// slice of float32-quantised elements — SearchVector and indexed vector
// elements are 32-bit IEEE-754 values (SearchVectors API reference), so every
// element is rounded to its float32 value and the persisted entry and the
// distance arithmetic only ever see values an exact float32 can represent.
// It returns nil when the attribute is absent, is not a list of numbers, or
// its length does not match the index dimension count: such items are simply
// not indexed — the write itself succeeds.
func extractVector(vi *VectorIndex, item *Item) []float64 {
	var av *AttributeValue
	if item.Attributes != nil {
		av = item.Attributes[vi.VectorAttributeName]
	}
	if av == nil && item.Key != nil {
		av = item.Key[vi.VectorAttributeName]
	}
	if av == nil || av.L == nil || int64(len(av.L)) != vi.Dimensions {
		return nil
	}
	vec := make([]float64, len(av.L))
	for i, el := range av.L {
		if el == nil || el.N == nil {
			return nil
		}
		// bitSize 32 rounds to the float32 value and returns it as an exact
		// float64, so the proto's double field holds float32-valued doubles.
		f, err := strconv.ParseFloat(*el.N, 32)
		if err != nil {
			return nil
		}
		vec[i] = f
	}
	return vec
}

func vectorEntryKey(tableName, indexName, primaryKey string) string {
	return tableName + KeySep + indexName + KeySep + primaryKey
}

// putVectorEntries writes (or clears, when the item no longer carries a valid
// vector) the entry of every vector index on the table for the item.
func putVectorEntries(txn storage.Transaction, region string, table *Table, item *Item) error {
	bucket := txn.Bucket(vectorIndexBucketName(region))
	primaryKey := EncodeItemKey(table.Name, item.Key, table)
	if primaryKey == "" {
		return nil
	}
	for _, vi := range table.VectorIndexes {
		key := vectorEntryKey(table.Name, vi.IndexName, primaryKey)
		vec := extractVector(vi, item)
		if vec == nil {
			// The attribute is missing, malformed, or has the wrong
			// dimension count: whatever entry existed for this item must
			// not survive under the new value.
			if err := bucket.Delete([]byte(key)); err != nil {
				return err
			}
			continue
		}
		value, err := proto.Marshal(&pb.VectorIndexEntry{Vector: vec})
		if err != nil {
			return fmt.Errorf("marshal vector entry for %s: %w", vi.IndexName, err)
		}
		if err := bucket.Put([]byte(key), value); err != nil {
			return err
		}
	}
	return nil
}

// deleteVectorEntries removes the item's entry from every vector index on
// the table.
func deleteVectorEntries(txn storage.Transaction, region string, table *Table, item *Item) error {
	bucket := txn.Bucket(vectorIndexBucketName(region))
	primaryKey := EncodeItemKey(table.Name, item.Key, table)
	if primaryKey == "" {
		return nil
	}
	for _, vi := range table.VectorIndexes {
		key := vectorEntryKey(table.Name, vi.IndexName, primaryKey)
		if err := bucket.Delete([]byte(key)); err != nil {
			return err
		}
	}
	return nil
}

// putVectorEntriesForIndex writes only the named vector index's entry for an
// item, leaving every other vector index untouched — the backfill twin of
// IndexStore.PutIndexEntriesForIndex. A name absent from the table's schema
// writes nothing: the index may have been deleted between the schema update
// and this backfill write.
func putVectorEntriesForIndex(txn storage.Transaction, region string, table *Table, indexName string, item *Item) error {
	for _, vi := range table.VectorIndexes {
		if vi.IndexName != indexName {
			continue
		}
		bucket := txn.Bucket(vectorIndexBucketName(region))
		primaryKey := EncodeItemKey(table.Name, item.Key, table)
		if primaryKey == "" {
			return nil
		}
		key := vectorEntryKey(table.Name, vi.IndexName, primaryKey)
		vec := extractVector(vi, item)
		if vec == nil {
			if err := bucket.Delete([]byte(key)); err != nil {
				return err
			}
			return nil
		}
		value, err := proto.Marshal(&pb.VectorIndexEntry{Vector: vec})
		if err != nil {
			return fmt.Errorf("marshal vector entry for %s: %w", vi.IndexName, err)
		}
		return bucket.Put([]byte(key), value)
	}
	return nil
}

// deleteVectorEntriesForIndex removes every entry of the named vector index
// so a deleted index's data does not outlive the index, mirroring
// IndexStore.DeleteIndexEntriesForIndex with the same batched deletes.
func deleteVectorEntriesForIndex(txn storage.Transaction, region, tableName, indexName string) error {
	return deletePrefixBatched(txn.Bucket(vectorIndexBucketName(region)), tableName+KeySep+indexName+KeySep)
}

// VectorSearchHit pairs the resolved item with its similarity score. Whether
// a lower or higher score means "more similar" depends on the index's
// distance function; hits are always ordered most-similar-first.
type VectorSearchHit struct {
	Item  *Item
	Score float64
}

// vectorTopK performs a brute-force similarity search over one vector index:
// every entry's vector is scored with the shared distance algorithms, entries
// whose dimensions no longer match the query are skipped, candidates are
// ordered best-first, and the filter (the search-condition expression) prunes
// items as they are loaded, so at most k passing items are read. COSINE and
// EUCLIDEAN scores are distances (k smallest, ascending); DOT_PRODUCT scores
// are similarities (k highest, descending) — the SearchVectors API reference
// semantics.
func vectorTopK(txn storage.Transaction, region string, table *Table, indexName string, query []float64, k int, filter func(*Item) bool) ([]VectorSearchHit, error) {
	if k < 1 {
		return nil, fmt.Errorf("vector top-k requires k >= 1, got %d", k)
	}

	var vi *VectorIndex
	for _, cand := range table.VectorIndexes {
		if cand.IndexName == indexName {
			vi = cand
			break
		}
	}
	if vi == nil {
		return nil, fmt.Errorf("vector index %s not found on table %s", indexName, table.Name)
	}

	var metric vecdist.DistanceMetric
	switch vi.DistanceFunction {
	case "COSINE":
		metric = vecdist.CosineDistance
	case "EUCLIDEAN":
		metric = vecdist.L2
	case "DOT_PRODUCT":
		metric = vecdist.DotProduct
	default:
		return nil, fmt.Errorf("unknown distance function %q for index %s", vi.DistanceFunction, indexName)
	}

	bucket := txn.Bucket(vectorIndexBucketName(region))
	prefix := table.Name + KeySep + indexName + KeySep

	type candidate struct {
		primaryKey string
		score      float64
	}
	candidates := make([]candidate, 0, k)

	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()
	for iter.Next() {
		var entry pb.VectorIndexEntry
		if err := proto.Unmarshal(iter.Value(), &entry); err != nil {
			return nil, fmt.Errorf("unmarshal vector entry under %s: %w", iter.Key(), err)
		}
		if len(entry.Vector) != len(query) {
			continue
		}
		score, err := vecdist.ComputeDistance(query, entry.Vector, metric)
		if err != nil {
			continue
		}
		candidates = append(candidates, candidate{
			primaryKey: string(iter.Key())[len(prefix):],
			score:      score,
		})
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}

	descending := vi.DistanceFunction == "DOT_PRODUCT"
	sort.Slice(candidates, func(i, j int) bool {
		if descending {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].score < candidates[j].score
	})

	itemBucket := txn.Bucket(itemBucketName(region))
	hits := make([]VectorSearchHit, 0, k)
	for _, c := range candidates {
		if len(hits) >= k {
			break
		}
		data, err := itemBucket.Get([]byte(c.primaryKey))
		if err != nil {
			return nil, fmt.Errorf("get item %s for vector search: %w", c.primaryKey, err)
		}
		if data == nil {
			// The item was deleted after its entry was written; the
			// entry is stale and scores nothing.
			continue
		}
		var pbItem pb.Item
		if err := proto.Unmarshal(data, &pbItem); err != nil {
			return nil, fmt.Errorf("unmarshal item %s for vector search: %w", c.primaryKey, err)
		}
		item := itemFromProto(&pbItem)
		if filter != nil && !filter(item) {
			continue
		}
		hits = append(hits, VectorSearchHit{Item: item, Score: c.score})
	}
	return hits, nil
}
