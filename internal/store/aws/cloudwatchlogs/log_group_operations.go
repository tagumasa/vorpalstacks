package cloudwatchlogs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
	"vorpalstacks/internal/store/aws/common"
)

// CreateLogGroup creates a new CloudWatch Logs log group. The existence
// check and the write run under the group's lock, the same lock the
// deletes hold: a create interleaving with DeleteLogGroup cannot
// resurrect a family the delete is tearing down, and two concurrent
// creators of one name admit exactly one record.
func (s *Store) CreateLogGroup(lg *LogGroup) error {
	gLock := s.groupLock(lg.Name)
	gLock.Lock()
	defer gLock.Unlock()

	key := s.logGroupKey(lg.Name)
	if s.Exists(key) {
		return ErrLogGroupAlreadyExists
	}
	lg.ARN = s.arnBuilder.CloudWatch().LogGroup(lg.Name)
	return s.PutProto(key, LogGroupToProto(lg))
}

// GetLogGroup retrieves a CloudWatch Logs log group by name. A missing
// record reports the not-found sentinel; any other storage failure
// (corrupt record, backend error) propagates as a storage error so the
// caller does not mistake an outage for an absent resource.
func (s *Store) GetLogGroup(name string) (*LogGroup, error) {
	key := s.logGroupKey(name)
	var p pb.LogGroup
	if err := s.GetProto(key, &p); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrLogGroupNotFound
		}
		return nil, err
	}
	return ProtoToLogGroup(&p), nil
}

// MutateLogGroup reads a log group record, applies fn, and persists the
// result under the group's lock — the same lock that serialises
// ingestion's metadata commit, so a configuration change cannot lose a
// concurrent StoredBytes or MetricFilterCount update and vice versa. fn
// must not retain the record beyond the call.
func (s *Store) MutateLogGroup(name string, fn func(*LogGroup) error) error {
	gLock := s.groupLock(name)
	gLock.Lock()
	defer gLock.Unlock()
	return s.mutateLogGroupLocked(name, fn)
}

func (s *Store) mutateLogGroupLocked(name string, fn func(*LogGroup) error) error {
	lg, err := s.GetLogGroup(name)
	if err != nil {
		return err
	}
	if err := fn(lg); err != nil {
		return err
	}
	return s.putLogGroupLocked(lg)
}

// PutLogGroup writes a log group record under the group's lock: a
// full-record overwrite serialises with ingestion's StoredBytes commits
// and the mutate seam, so it cannot clobber a concurrent metadata update
// with a stale snapshot.
func (s *Store) PutLogGroup(lg *LogGroup) error {
	gLock := s.groupLock(lg.Name)
	gLock.Lock()
	defer gLock.Unlock()
	return s.putLogGroupLocked(lg)
}

// putLogGroupLocked is the lock-free write body for callers already
// holding the group's lock (the mutate seam and the ingestion commit
// paths).
func (s *Store) putLogGroupLocked(lg *LogGroup) error {
	key := s.logGroupKey(lg.Name)
	return s.PutProto(key, LogGroupToProto(lg))
}

// DeleteLogGroup deletes a CloudWatch Logs log group together with every
// record and chunk file it owns — streams, chunks (with a backstop for
// orphaned index records), metric filters, subscription filters, the
// data protection policy and tags are the complete per-group family
// list; a family added to the key layout joins this list or it leaks on
// every group delete. The teardown continues past individual
// failures and tolerates files already removed by an earlier partial
// delete, so a retry always makes progress. The group record itself is
// removed only when every step succeeded: deleting it while sub-resources
// survive would orphan them permanently, because a retry stops at the
// missing group before it can reach them. The whole teardown holds the
// group's lock, so an in-flight PutLogEvents either commits before the
// delete (and is then removed with the rest) or observes the group gone
// once it acquires the lock — a write can no longer resurrect the
// deleted records.
func (s *Store) DeleteLogGroup(name string) error {
	gLock := s.groupLock(name)
	gLock.Lock()
	defer gLock.Unlock()
	return s.deleteLogGroupLocked(name, false)
}

// DeleteLogGroupIfUnprotected deletes a log group unless deletion
// protection is enabled, deciding inside the delete's critical section:
// a protection enable that commits after a caller-side read is still
// honoured, and no window exists between the check and the teardown.
// The plain DeleteLogGroup deliberately bypasses the guard — it serves
// the create-path rollback of a group the caller just created.
func (s *Store) DeleteLogGroupIfUnprotected(name string) error {
	gLock := s.groupLock(name)
	gLock.Lock()
	defer gLock.Unlock()
	return s.deleteLogGroupLocked(name, true)
}

func (s *Store) deleteLogGroupLocked(name string, honourProtection bool) error {
	lg, err := s.GetLogGroup(name)
	if err != nil {
		return err
	}
	if honourProtection && lg.DeletionProtectionEnabled {
		return ErrLogGroupDeletionProtected
	}

	var errs []error

	if err := s.deleteAllLogStreams(name); err != nil {
		errs = append(errs, fmt.Errorf("delete log streams: %w", err))
	}

	// Backstop for chunk index records whose stream record is already
	// gone: DeleteLogStream covers only chunks listed under a live
	// stream. The file may already have been removed by an earlier
	// partial delete, which is not an error here. A listing failure is
	// reported: without the list the backstop cannot know what to remove,
	// and the group record must survive so a retry can reach the records.
	chunks, listErr := s.ListChunksForLogGroup(name)
	if listErr != nil {
		errs = append(errs, fmt.Errorf("list chunks: %w", listErr))
	}
	for _, chunk := range chunks {
		cp, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := os.Remove(cp); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
		if err := s.Delete(s.chunkIndexKey(name, chunk.LogStream, chunk.ChunkID)); err != nil {
			errs = append(errs, err)
		}
	}
	metricFilterMarker := ""
	for {
		metricFilters, nextToken, listErr := s.ListMetricFilters(name, "", metricFilterMarker, ListingPageSize)
		if listErr != nil {
			// Without the list the filters cannot be deleted; report the
			// failure so the group record survives and a retry can reach them.
			errs = append(errs, fmt.Errorf("list metric filters: %w", listErr))
			break
		}
		for _, mf := range metricFilters {
			if err := s.Delete(s.metricFilterKey(name, mf.Name)); err != nil {
				errs = append(errs, fmt.Errorf("delete metric filter %s: %w", mf.Name, err))
			}
		}
		if nextToken == "" {
			break
		}
		metricFilterMarker = nextToken
	}
	subFilters, listErr := s.ListSubscriptionFilters(name, "")
	if listErr != nil {
		errs = append(errs, fmt.Errorf("list subscription filters: %w", listErr))
	}
	for _, sf := range subFilters {
		if err := s.Delete(s.subscriptionFilterKey(name, sf.FilterName)); err != nil {
			errs = append(errs, fmt.Errorf("delete subscription filter %s: %w", sf.FilterName, err))
		}
	}
	// The group's data protection policy is keyed on the resolved group
	// name (the Put path resolves name-or-ARN before storing); a group
	// without a policy deletes nothing here, and a stale policy can never
	// survive into a same-named recreation.
	if err := s.Delete(s.dataProtectionPolicyKey(name)); err != nil {
		errs = append(errs, fmt.Errorf("delete data protection policy: %w", err))
	}
	// The bearer token authentication switch rides the group record
	// itself, so it dies with the record below — a same-named recreation
	// starts with the switch disabled.
	// The group's transformer and its transformed-message records die
	// with the group: a same-named recreation starts without them.
	if s.Exists(s.transformerKey(name)) {
		if err := s.DeleteTransformer(name); err != nil {
			errs = append(errs, fmt.Errorf("delete transformer: %w", err))
		}
	}
	if err := s.DeleteTransformedMessages(name); err != nil {
		errs = append(errs, fmt.Errorf("delete transformed messages: %w", err))
	}
	// The group's field index policy, its INACTIVE field trail and its
	// accumulated field-index bounds die with the group: a same-named
	// recreation starts unindexed.
	if s.Exists(s.indexPolicyKey(name)) {
		if err := s.DeleteIndexPolicy(name); err != nil {
			errs = append(errs, fmt.Errorf("delete index policy: %w", err))
		}
	}
	if err := s.DeleteIndexInactiveTrail(name); err != nil {
		errs = append(errs, fmt.Errorf("delete index inactive trail: %w", err))
	}
	if err := s.DeleteFieldIndexBounds(name); err != nil {
		errs = append(errs, fmt.Errorf("delete field index bounds: %w", err))
	}
	arn := s.arnBuilder.CloudWatch().LogGroup(name)
	if err := s.tagStore.Delete(arn); err != nil {
		errs = append(errs, fmt.Errorf("delete tags for log group %s: %w", name, err))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	if err := s.Delete(s.logGroupKey(name)); err != nil {
		return err
	}
	return nil
}

// deleteAllLogStreams removes every stream of a log group. The caller
// must hold the group's lock (DeleteLogGroup does).
func (s *Store) deleteAllLogStreams(logGroupName string) error {
	marker := ""
	for {
		streams, nextMarker, err := s.ListLogStreams(logGroupName, "", marker, ListingPageSize)
		if err != nil {
			return err
		}
		for _, stream := range streams {
			if err := s.deleteLogStreamLocked(logGroupName, stream.Name); err != nil {
				return err
			}
		}
		if nextMarker == "" {
			return nil
		}
		marker = nextMarker
	}
}

// CountLogGroups counts the log groups the region holds — the CreateLogGroup
// quota's input ("You can create up to 1,000,000 log groups per Region per
// account"). The count-then-create check is advisory against concurrent
// creations on one region; the uniqueness check inside CreateLogGroup
// remains the authority for collisions.
func (s *Store) CountLogGroups() (int, error) {
	count := 0
	if err := s.ScanPrefix(keyPrefixLogGroup, func(_ string, _ []byte) error {
		count++
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}

// ListLogGroups lists CloudWatch Logs log groups with optional prefix and pagination.
func (s *Store) ListLogGroups(prefix, marker string, maxItems int) ([]*LogGroup, string, error) {
	if maxItems <= 0 {
		maxItems = DefaultListLogGroupsLimit
	}

	opts := common.ListOptions{
		Prefix:   keyPrefixLogGroup,
		Marker:   marker,
		MaxItems: maxItems,
	}

	result, err := common.ListProto[*pb.LogGroup](s.BaseStore, opts, func() *pb.LogGroup { return &pb.LogGroup{} }, func(lg *pb.LogGroup) bool {
		if prefix != "" && lg.Name != prefix && !strings.HasPrefix(lg.Name, prefix) {
			return false
		}
		return true
	})
	if err != nil {
		return nil, "", err
	}

	groups := make([]*LogGroup, len(result.Items))
	for i, p := range result.Items {
		groups[i] = ProtoToLogGroup(p)
	}
	return groups, result.NextMarker, nil
}

// streamBounds is one stream's surviving timestamp range across a purge
// pass — kept chunks and expired-but-removal-failed chunks alike, every
// chunk whose index record still serves reads when the pass ends.
type streamBounds struct{ first, last int64 }

// PurgeExpiredChunks removes expired log chunks from a log group based on retention policy.
func (s *Store) PurgeExpiredChunks(logGroupName string, cutoffTime int64) (int64, error) {
	// The whole purge holds the group's lock: chunk files are created and
	// indexed under the same lock, and the StoredBytes decrement must not
	// lose a concurrent ingestion's increment.
	gLock := s.groupLock(logGroupName)
	gLock.Lock()
	defer gLock.Unlock()

	if _, err := s.GetLogGroup(logGroupName); err != nil {
		return 0, err
	}

	chunks, err := s.ListChunksForLogGroup(logGroupName)
	if err != nil {
		return 0, err
	}
	var totalBytesRemoved int64

	// The purge advances the affected streams' boundary members to the
	// data that survives it: firstEventTimestamp and lastEventTimestamp
	// describe stored events, and the purged era's values would otherwise
	// outlive the events they describe.
	remaining := map[string]*streamBounds{}
	purged := map[string]bool{}

	for _, chunk := range chunks {
		if chunk.MaxTs >= cutoffTime {
			mergeRemainingBounds(remaining, chunk)
			continue
		}

		chunkPath, err := s.safeChunkPath(chunk.ChunkPath)
		if err != nil {
			// The chunk stays on the read surface (its index record
			// survives): it is a survivor for the boundary rewrite, not
			// a purged entry.
			mergeRemainingBounds(remaining, chunk)
			continue
		}

		if err := os.Remove(chunkPath); err != nil && !os.IsNotExist(err) {
			mergeRemainingBounds(remaining, chunk)
			continue
		}

		indexKey := s.chunkIndexKey(logGroupName, chunk.LogStream, chunk.ChunkID)
		if err := s.Delete(indexKey); err != nil {
			// The index record survives, so the next purge re-lists this
			// chunk and retries the file removal and the decrement —
			// counting its bytes now would double-decrement on the retry.
			// The record surviving also keeps the chunk readable, so its
			// bounds belong to the stream's surviving set.
			logs.Error("Failed to delete chunk index", logs.String("key", indexKey), logs.Err(err))
			mergeRemainingBounds(remaining, chunk)
			continue
		}
		purged[chunk.LogStream] = true

		// The decrement uses the chunk's recorded ingested message bytes —
		// the same basis PutLogEvents incremented on — and lands
		// immediately after the index record it belongs to: the deleted
		// record is what stops a later purge decrementing the chunk twice,
		// and batching decrements to the end of the pass would lose every
		// one of them if the group write failed after the records were
		// already gone.
		if err := s.mutateLogGroupLocked(logGroupName, func(lg *LogGroup) error {
			lg.StoredBytes -= chunk.ByteSize
			if lg.StoredBytes < 0 {
				lg.StoredBytes = 0
			}
			return nil
		}); err != nil {
			return totalBytesRemoved, err
		}
		totalBytesRemoved += chunk.ByteSize
	}

	for stream := range purged {
		ls, err := s.GetLogStream(logGroupName, stream)
		if err != nil {
			continue
		}
		if b := remaining[stream]; b != nil {
			ls.FirstEventTs, ls.LastEventTs = b.first, b.last
		} else {
			ls.FirstEventTs, ls.LastEventTs = 0, 0
		}
		if err := s.PutProto(s.logStreamKey(logGroupName, stream), LogStreamToProto(ls)); err != nil {
			return totalBytesRemoved, err
		}
	}

	return totalBytesRemoved, nil
}

// mergeRemainingBounds folds one surviving chunk's timestamp range into
// the stream's surviving bounds — a kept chunk, or an expired chunk whose
// removal failed and whose index record therefore still serves reads.
// The boundary rewrite must describe exactly the chunks the read surface
// still holds.
func mergeRemainingBounds(remaining map[string]*streamBounds, chunk *ChunkMeta) {
	if b := remaining[chunk.LogStream]; b == nil {
		remaining[chunk.LogStream] = &streamBounds{first: chunk.MinTs, last: chunk.MaxTs}
	} else {
		if chunk.MinTs < b.first {
			b.first = chunk.MinTs
		}
		if chunk.MaxTs > b.last {
			b.last = chunk.MaxTs
		}
	}
}

// PurgeAllExpiredChunks purges expired chunks from all log groups based on their retention policies.
func (s *Store) PurgeAllExpiredChunks() error {
	now := time.Now().UnixMilli()

	marker := ""
	for {
		groups, nextMarker, err := s.ListLogGroups("", marker, ListingPageSize)
		if err != nil {
			return err
		}

		for _, lg := range groups {
			if lg.RetentionInDays <= 0 {
				continue
			}

			cutoffTime := now - int64(lg.RetentionInDays)*24*60*60*1000
			if _, err := s.PurgeExpiredChunks(lg.Name, cutoffTime); err != nil {
				continue
			}
			// The transformed-message family expires with the chunks it
			// mirrors.
			s.PurgeExpiredTransformedMessages(lg.Name, cutoffTime)
		}

		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}

	return nil
}
