package kinesis

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

// PutShard stores a shard in the data store.
func (s *KinesisStore) PutShard(shard *Shard) error {
	key := fmt.Sprintf("%s#%s", shard.StreamName, shard.ShardID)
	return s.shardsStore.PutProto(key, ShardToProto(shard))
}

// GetShard retrieves a shard from a stream by its ID.
func (s *KinesisStore) GetShard(streamName, shardID string) (*Shard, error) {
	key := fmt.Sprintf("%s#%s", streamName, shardID)
	var pbShard pb.Shard
	if err := s.shardsStore.GetProto(key, &pbShard); err != nil {
		return nil, ErrShardNotFound
	}
	return ProtoToShard(&pbShard), nil
}

// UpdateShard updates an existing shard in the data store.
func (s *KinesisStore) UpdateShard(shard *Shard) error {
	return s.PutShard(shard)
}

var errStopIteration = errors.New("stop iteration")

// ListShards returns all shards in a stream, optionally filtered by a ShardFilter and paginated.
func (s *KinesisStore) ListShards(streamName string, filter *ShardFilter, exclusiveStartShardID string, limit int) ([]*Shard, error) {
	prefix := streamName + "#"
	var shards []*Shard
	skip := exclusiveStartShardID != ""

	// The trim-horizon-relative filter types need the stream's retention
	// cutoff; an unfiltered scan (nil filter) never consults it.
	var trimHorizon time.Time
	if filter != nil && shardFilterNeedsTrimHorizon(filter.Type) {
		stream, err := s.GetStream(streamName)
		if err != nil {
			return nil, err
		}
		trimHorizon = RetentionCutoff(stream.RetentionPeriodHours)
	}

	err := s.shardsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if skip {
			var pbShard pb.Shard
			if err := proto.Unmarshal(value, &pbShard); err != nil {
				return fmt.Errorf("unmarshal shard for skip check: %w", err)
			}
			if pbShard.ShardId == exclusiveStartShardID {
				skip = false
			}
			return nil
		}

		var pbShard pb.Shard
		if err := proto.Unmarshal(value, &pbShard); err != nil {
			return fmt.Errorf("unmarshal shard: %w", err)
		}
		shard := ProtoToShard(&pbShard)
		if filter == nil || s.shardMatchesFilter(shard, filter, trimHorizon) {
			shards = append(shards, shard)
			if limit > 0 && len(shards) >= limit {
				return errStopIteration
			}
		}
		return nil
	})
	if errors.Is(err, errStopIteration) {
		err = nil
	}
	if err != nil {
		return shards, fmt.Errorf("scan shards prefix: %w", err)
	}
	return shards, nil
}

// shardFilterNeedsTrimHorizon reports whether a filter type's semantics
// reference the stream's trim horizon.
func shardFilterNeedsTrimHorizon(filterType string) bool {
	switch filterType {
	case "AT_TRIM_HORIZON", "FROM_TRIM_HORIZON", "FROM_TIMESTAMP":
		return true
	}
	return false
}

// shardClosedAt reads a closed shard's end timestamp from its ending
// sequence number's arrival-time component. A sequence number that cannot
// be decoded reads as the shard's creation — older than any real close, so
// the horizon-relative comparisons fail closed.
func shardClosedAt(shard *Shard) time.Time {
	if shard.SequenceNumberRange == nil {
		return shard.CreatedAt
	}
	if t, ok := sequenceNumberTime(shard.SequenceNumberRange.EndingSequenceNumber); ok {
		return t
	}
	return shard.CreatedAt
}

// GetChildShards returns the direct child shards of the given shard.
// Child shards are identified by ParentShardID or AdjacentParentShardID
// matching the given shardID. Only open shards are considered because
// the consumer should transition to them after the parent is closed.
func (s *KinesisStore) GetChildShards(streamName, shardID string) ([]*Shard, error) {
	prefix := streamName + "#"
	var children []*Shard

	err := s.shardsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbShard pb.Shard
		if err := proto.Unmarshal(value, &pbShard); err != nil {
			return fmt.Errorf("unmarshal shard for child lookup: %w", err)
		}
		shard := ProtoToShard(&pbShard)
		// Only open shards can be children (closed shards have ended); a
		// shard without a range record reads as open, the same reading the
		// record path's active-shard filter takes.
		if shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber != "" {
			return nil
		}
		if shard.ParentShardID == shardID || shard.AdjacentParentShardID == shardID {
			children = append(children, shard)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan shards for child lookup: %w", err)
	}
	return children, nil
}

func (s *KinesisStore) listShardsInTxn(txn storage.Transaction, streamName string) ([]*Shard, error) {
	bucket := txn.Bucket(s.shardsBucketName())
	prefix := []byte(streamName + "#")
	var shards []*Shard

	iter := bucket.ScanPrefix(prefix)
	for iter.Next() {
		var pbShard pb.Shard
		if err := proto.Unmarshal(iter.Value(), &pbShard); err != nil {
			return nil, fmt.Errorf("unmarshal shard: %w", err)
		}
		shard := ProtoToShard(&pbShard)
		// A shard without a range record reads as open, the same reading
		// the record path's active-shard filter takes.
		if shard.SequenceNumberRange == nil || shard.SequenceNumberRange.EndingSequenceNumber == "" {
			shards = append(shards, shard)
		}
	}
	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("scan shards iterator: %w", err)
	}
	return shards, nil
}

// shardMatchesFilter evaluates one shard against a ShardFilter per the
// Type member's documented semantics. trimHorizon is the stream's
// retention cutoff (now minus the retention period); a closed shard's end
// timestamp is the arrival time encoded in its ending sequence number. A
// filter type outside the enum matches nothing — the Core rejects those
// before the scan, so the default is a defensive no-match.
func (s *KinesisStore) shardMatchesFilter(shard *Shard, filter *ShardFilter, trimHorizon time.Time) bool {
	open := shard.SequenceNumberRange == nil || shard.SequenceNumberRange.EndingSequenceNumber == ""
	switch filter.Type {
	case "AFTER_SHARD_ID":
		// All shards — open and closed — starting with the shard whose ID
		// immediately follows the provided ShardId.
		return shard.ShardID > filter.ShardID
	case "AT_LATEST":
		// Only the currently open shards.
		return open
	case "AT_TRIM_HORIZON":
		// The shards that were open at the trim horizon: created at or
		// before it, and either still open or closed at or after it.
		return !shard.CreatedAt.After(trimHorizon) && (open || !shardClosedAt(shard).Before(trimHorizon))
	case "FROM_TRIM_HORIZON":
		// All shards within the retention period (trim to tip): open
		// shards, and closed shards closed at or after the trim horizon.
		return open || !shardClosedAt(shard).Before(trimHorizon)
	case "AT_TIMESTAMP":
		if filter.Timestamp == nil {
			return false
		}
		// Shards whose start timestamp is at or before the given timestamp
		// and whose end timestamp is at or after it, or that are still open.
		return !shard.CreatedAt.After(*filter.Timestamp) && (open || !shardClosedAt(shard).Before(*filter.Timestamp))
	case "FROM_TIMESTAMP":
		if filter.Timestamp == nil {
			return false
		}
		// All closed shards whose end timestamp is at or after the given
		// timestamp, and all open shards; a timestamp before the trim
		// horizon is corrected to the trim horizon.
		from := *filter.Timestamp
		if from.Before(trimHorizon) {
			from = trimHorizon
		}
		return open || !shardClosedAt(shard).Before(from)
	}
	return false
}

func (s *KinesisStore) getShardInTxn(txn storage.Transaction, streamName, shardID string) (*Shard, error) {
	bucket := txn.Bucket(s.shardsBucketName())
	key := []byte(fmt.Sprintf("%s#%s", streamName, shardID))
	data, err := bucket.Get(key)
	if err != nil {
		return nil, fmt.Errorf("get shard from bucket: %w", err)
	}
	if data == nil {
		return nil, ErrShardNotFound
	}
	var pbShard pb.Shard
	if err := proto.Unmarshal(data, &pbShard); err != nil {
		return nil, fmt.Errorf("unmarshal shard data: %w", err)
	}
	return ProtoToShard(&pbShard), nil
}

func (s *KinesisStore) putShardInTxn(txn storage.Transaction, shard *Shard) error {
	bucket := txn.Bucket(s.shardsBucketName())
	key := []byte(fmt.Sprintf("%s#%s", shard.StreamName, shard.ShardID))
	data, err := proto.Marshal(ShardToProto(shard))
	if err != nil {
		return fmt.Errorf("marshal shard: %w", err)
	}
	return bucket.Put(key, data)
}

// nextShardIDNumberInTxn allocates the next shard ID number for a stream:
// per-stream uniqueness is both the storage-key contract and the AWS one.
// The number is one past the highest numeric suffix among the stream's
// existing shards — open and closed — so a fresh child can never collide
// with, and overwrite, a live or closed record the way a store-lifetime
// counter could against a stream whose initial shards carry index-based
// IDs.
func (s *KinesisStore) nextShardIDNumberInTxn(txn storage.Transaction, streamName string) (int64, error) {
	bucket := txn.Bucket(s.shardsBucketName())
	prefix := []byte(streamName + "#")

	var maxID int64
	iter := bucket.ScanPrefix(prefix)
	for iter.Next() {
		shardID := string(iter.Key()[len(prefix):])
		numStr := strings.TrimPrefix(shardID, "shardId-")
		if num, err := strconv.ParseInt(numStr, 10, 64); err == nil && num > maxID {
			maxID = num
		}
	}
	if err := iter.Error(); err != nil {
		return 0, fmt.Errorf("scan shards for ID allocation: %w", err)
	}
	return maxID + 1, nil
}

func formatShardID(number int64) string {
	return fmt.Sprintf("shardId-%012d", number)
}

// SplitShard splits a shard into two new shards.
func (s *KinesisStore) SplitShard(streamName, shardID, newStartingHashKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.splitShardInternal(streamName, shardID, newStartingHashKey)
}

func (s *KinesisStore) splitShardInternal(streamName, shardID, newStartingHashKey string) error {
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return s.splitShardInTxn(txn, streamName, shardID, newStartingHashKey)
	})
}
func (s *KinesisStore) splitShardInTxn(txn storage.Transaction, streamName, shardID, newStartingHashKey string) error {
	shard, err := s.getShardInTxn(txn, streamName, shardID)
	if err != nil {
		return fmt.Errorf("get existing shard: %w", err)
	}

	// A shard without a range record reads as open — the same reading the
	// record path's active-shard filter takes.
	if shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber != "" {
		return ErrShardClosed
	}

	nextNumber, err := s.nextShardIDNumberInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("allocate child shard IDs: %w", err)
	}

	closed, child1, child2, err := s.splitShardRecord(shard, newStartingHashKey, formatShardID(nextNumber), formatShardID(nextNumber+1))
	if err != nil {
		return err
	}

	if err := s.putShardInTxn(txn, closed); err != nil {
		return fmt.Errorf("put closed shard: %w", err)
	}
	if err := s.putShardInTxn(txn, child1); err != nil {
		return fmt.Errorf("put new shard 1: %w", err)
	}
	if err := s.putShardInTxn(txn, child2); err != nil {
		return fmt.Errorf("put new shard 2: %w", err)
	}

	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("get stream for split: %w", err)
	}
	stream.ShardCount++
	return s.updateStreamInTxn(txn, stream)
}

// closeShard seals a shard's sequence-number range at its latest written
// position — the starting sequence number when nothing was ever written.
// A shard record without a range submessage gets one here: sealing writes
// the range it will now carry.
func closeShard(shard *Shard) {
	if shard.SequenceNumberRange == nil {
		shard.SequenceNumberRange = &SequenceNumberRange{}
	}
	shard.SequenceNumberRange.EndingSequenceNumber = shard.LatestSequenceNumber
	if shard.LatestSequenceNumber == "" {
		shard.SequenceNumberRange.EndingSequenceNumber = shard.SequenceNumberRange.StartingSequenceNumber
	}
}

// splitShardRecord computes the closed parent and the two open children of
// splitting one shard.
//
// When the caller specifies NewStartingHashKey (NSK), per the Smithy
// documentation NSK is the starting hash key of child 2 (upper):
//
//	child1 = [start, NSK-1],  child2 = [NSK, end]
//
// When auto-midpoint is used, midHash is the ending key of child 1:
//
//	child1 = [start, mid],    child2 = [mid+1, end]
func (s *KinesisStore) splitShardRecord(shard *Shard, newStartingHashKey, childID1, childID2 string) (*Shard, *Shard, *Shard, error) {
	startHash, ok := new(big.Int).SetString(shard.HashKeyRange.StartingHashKey, 10)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid starting hash key: %s", shard.HashKeyRange.StartingHashKey)
	}
	endHash, ok := new(big.Int).SetString(shard.HashKeyRange.EndingHashKey, 10)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid ending hash key: %s", shard.HashKeyRange.EndingHashKey)
	}
	var midHash *big.Int
	if newStartingHashKey != "" {
		var ok bool
		midHash, ok = new(big.Int).SetString(newStartingHashKey, 10)
		if !ok {
			return nil, nil, nil, fmt.Errorf("invalid new starting hash key: %s", newStartingHashKey)
		}
	} else {
		midHash = new(big.Int).Div(new(big.Int).Add(startHash, endHash), big.NewInt(2))
	}

	// Compute child shard hash key boundaries.
	var child1End, child2Start *big.Int
	if newStartingHashKey != "" {
		child2Start = midHash
		child1End = new(big.Int).Sub(midHash, big.NewInt(1))
	} else {
		child1End = midHash
		child2Start = new(big.Int).Add(midHash, big.NewInt(1))
	}

	closed := *shard
	closeShard(&closed)

	newShard1 := NewShard(childID1, shard.StreamName, shard.HashKeyRange.StartingHashKey, child1End.String())
	newShard1.ParentShardID = shard.ShardID
	newShard1.SequenceNumberRange = &SequenceNumberRange{
		StartingSequenceNumber: s.generateSequenceNumber(childID1),
	}

	newShard2 := NewShard(childID2, shard.StreamName, child2Start.String(), shard.HashKeyRange.EndingHashKey)
	newShard2.ParentShardID = shard.ShardID
	newShard2.SequenceNumberRange = &SequenceNumberRange{
		StartingSequenceNumber: s.generateSequenceNumber(childID2),
	}

	return &closed, newShard1, newShard2, nil
}

// MergeShards merges two adjacent shards into one.
func (s *KinesisStore) MergeShards(streamName, shardID1, shardID2 string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mergeShardsInternal(streamName, shardID1, shardID2)
}

func (s *KinesisStore) mergeShardsInternal(streamName, shardID1, shardID2 string) error {
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return s.mergeShardsInTxn(txn, streamName, shardID1, shardID2)
	})
}

func (s *KinesisStore) mergeShardsInTxn(txn storage.Transaction, streamName, shardID1, shardID2 string) error {
	shard1, err := s.getShardInTxn(txn, streamName, shardID1)
	if err != nil {
		return fmt.Errorf("get shard 1: %w", err)
	}
	shard2, err := s.getShardInTxn(txn, streamName, shardID2)
	if err != nil {
		return fmt.Errorf("get shard 2: %w", err)
	}

	if !s.areAdjacentShards(shard1, shard2) {
		return ErrShardsNotAdjacent
	}

	nextNumber, err := s.nextShardIDNumberInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("allocate merged shard ID: %w", err)
	}

	closed1, closed2, child, err := s.mergedShardRecord(shard1, shard2, formatShardID(nextNumber))
	if err != nil {
		return err
	}

	if err := s.putShardInTxn(txn, closed1); err != nil {
		return fmt.Errorf("put closed shard 1: %w", err)
	}
	if err := s.putShardInTxn(txn, closed2); err != nil {
		return fmt.Errorf("put closed shard 2: %w", err)
	}
	if err := s.putShardInTxn(txn, child); err != nil {
		return fmt.Errorf("put merged shard: %w", err)
	}

	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("get stream for merge: %w", err)
	}
	stream.ShardCount--
	return s.updateStreamInTxn(txn, stream)
}

// mergedShardRecord computes the closed pair and the single open child of
// merging two adjacent shards: the child spans both ranges, in key order
// regardless of the argument order.
func (s *KinesisStore) mergedShardRecord(shard1, shard2 *Shard, childID string) (*Shard, *Shard, *Shard, error) {
	start1, ok1 := new(big.Int).SetString(shard1.HashKeyRange.StartingHashKey, 10)
	start2, ok2 := new(big.Int).SetString(shard2.HashKeyRange.StartingHashKey, 10)
	if !ok1 || !ok2 {
		return nil, nil, nil, fmt.Errorf("invalid hash key range")
	}

	lower, upper := shard1, shard2
	if start1.Cmp(start2) > 0 {
		lower, upper = shard2, shard1
	}

	closed1 := *shard1
	closeShard(&closed1)
	closed2 := *shard2
	closeShard(&closed2)

	newShard := NewShard(childID, shard1.StreamName, lower.HashKeyRange.StartingHashKey, upper.HashKeyRange.EndingHashKey)
	newShard.ParentShardID = shard1.ShardID
	newShard.AdjacentParentShardID = shard2.ShardID
	newShard.SequenceNumberRange = &SequenceNumberRange{
		StartingSequenceNumber: s.generateSequenceNumber(childID),
	}
	return &closed1, &closed2, newShard, nil
}

func (s *KinesisStore) areAdjacentShards(shard1, shard2 *Shard) bool {
	end1, ok1 := new(big.Int).SetString(shard1.HashKeyRange.EndingHashKey, 10)
	start2, ok2 := new(big.Int).SetString(shard2.HashKeyRange.StartingHashKey, 10)
	if ok1 && ok2 {
		diff := new(big.Int).Sub(start2, end1)
		if diff.Cmp(big.NewInt(1)) == 0 {
			return true
		}
	}

	end2, ok2 := new(big.Int).SetString(shard2.HashKeyRange.EndingHashKey, 10)
	start1, ok1 := new(big.Int).SetString(shard1.HashKeyRange.StartingHashKey, 10)
	if ok1 && ok2 {
		diff := new(big.Int).Sub(start1, end2)
		if diff.Cmp(big.NewInt(1)) == 0 {
			return true
		}
	}

	return false
}

// hashKeyLess compares two decimal hash keys numerically. Hash key strings
// span 1 to 20 digits, so lexicographic order diverges from numeric order
// exactly on the initial shard layout of a four-shard stream. An
// unparseable key falls back to string order so the comparison stays
// total; malformed keys surface as errors in the adjacency check.
func hashKeyLess(a, b string) bool {
	x, okX := new(big.Int).SetString(a, 10)
	y, okY := new(big.Int).SetString(b, 10)
	if !okX || !okY {
		return a < b
	}
	return x.Cmp(y) < 0
}

// hashRangeWidth returns the inclusive span of a shard's hash key range;
// an unparseable boundary reads as zero width, deferring to any well-formed
// competitor.
func hashRangeWidth(shard *Shard) *big.Int {
	if shard.HashKeyRange == nil {
		return big.NewInt(0)
	}
	lo, okLo := new(big.Int).SetString(shard.HashKeyRange.StartingHashKey, 10)
	hi, okHi := new(big.Int).SetString(shard.HashKeyRange.EndingHashKey, 10)
	if !okLo || !okHi {
		return big.NewInt(0)
	}
	return new(big.Int).Add(new(big.Int).Sub(hi, lo), big.NewInt(1))
}

// shardIDNumber extracts the numeric suffix of a shard ID.
func shardIDNumber(shard *Shard) int64 {
	num, err := strconv.ParseInt(strings.TrimPrefix(shard.ShardID, "shardId-"), 10, 64)
	if err != nil {
		return 0
	}
	return num
}

// UpdateShardCount updates the shard count for a stream.
func (s *KinesisStore) UpdateShardCount(streamName string, targetCount int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return s.updateShardCountInTxn(txn, streamName, targetCount)
	})
}

// updateShardCountInTxn reshapes the stream to the target shard count.
// Transactional reads observe the state the transaction began with — the
// storage layer's documented semantics keep pending batch writes invisible
// to Get and iterators — so the reshaping is computed in memory over the
// working set and one write pass persists every outcome; a refresh read
// between merges would see the pre-scaling set forever and the loop would
// stall without converging.
func (s *KinesisStore) updateShardCountInTxn(txn storage.Transaction, streamName string, targetCount int32) error {
	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("get stream: %w", err)
	}

	if stream.StreamModeDetails != nil && stream.StreamModeDetails.StreamMode == StreamModeOnDemand {
		return ErrOperationNotSupportedOnOnDemandStream
	}

	open, err := s.listShardsInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("list current shards: %w", err)
	}

	currentCount := int32(len(open))
	if currentCount == targetCount {
		return nil
	}
	if currentCount == 0 {
		return ErrInvalidShardCount
	}

	ratio := float64(targetCount) / float64(currentCount)
	if ratio > 2.0 || ratio < 0.5 {
		return ErrInvalidShardCount
	}

	// The working set of open shards, the running maximum shard ID for
	// collision-free child allocation, and the record forms to persist.
	work := make([]*Shard, len(open))
	copy(work, open)
	maxNumber := int64(0)
	for _, shard := range open {
		if n := shardIDNumber(shard); n > maxNumber {
			maxNumber = n
		}
	}
	var writes []*Shard

	record := func(shards ...*Shard) {
		writes = append(writes, shards...)
	}

	if targetCount > currentCount {
		for int32(len(work)) < targetCount {
			if len(work) == 0 {
				break
			}
			// Split the widest hash range. The working set keeps encounter
			// order — originals hold their positions and children append —
			// so among equal widths an unsplit original is always chosen
			// before a child: no shard is split again while a same-width
			// shard of an earlier generation waits, which would close and
			// reopen a phantom generation inside one operation.
			parent := work[0]
			for _, shard := range work[1:] {
				if hashRangeWidth(shard).Cmp(hashRangeWidth(parent)) > 0 {
					parent = shard
				}
			}

			closed, child1, child2, err := s.splitShardRecord(parent, "", formatShardID(maxNumber+1), formatShardID(maxNumber+2))
			if err != nil {
				return err
			}
			maxNumber += 2
			record(closed, child1, child2)

			for i, shard := range work {
				if shard == parent {
					work = append(work[:i], work[i+1:]...)
					break
				}
			}
			work = append(work, child1, child2)
		}
	} else {
		sort.Slice(work, func(i, j int) bool {
			return hashKeyLess(work[i].HashKeyRange.StartingHashKey, work[j].HashKeyRange.StartingHashKey)
		})
		for int32(len(work)) > targetCount && len(work) >= 2 {
			merged := false
			for j := 0; j < len(work)-1; j++ {
				if !s.areAdjacentShards(work[j], work[j+1]) {
					continue
				}
				closed1, closed2, child, err := s.mergedShardRecord(work[j], work[j+1], formatShardID(maxNumber+1))
				if err != nil {
					return err
				}
				maxNumber++
				record(closed1, closed2, child)

				// The child starts where work[j] started, so replacing
				// the pair in place keeps the numeric order.
				spliced := make([]*Shard, 0, len(work)-1)
				spliced = append(spliced, work[:j]...)
				spliced = append(spliced, child)
				spliced = append(spliced, work[j+2:]...)
				work = spliced
				merged = true
				break
			}
			if !merged {
				break
			}
		}
	}

	// The reshape must land exactly on the requested target. A well-formed
	// stream's geometry always allows it; falling short anyway is a server
	// defect whose honest answer is the failure — never an intermediate
	// count persisted and reported as success.
	if int32(len(work)) != targetCount {
		return ErrInvalidShardCount
	}

	for _, shard := range writes {
		if err := s.putShardInTxn(txn, shard); err != nil {
			return fmt.Errorf("put reshaped shard %s: %w", shard.ShardID, err)
		}
	}

	stream.ShardCount = int32(len(work))
	return s.updateStreamInTxn(txn, stream)
}
