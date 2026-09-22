package cloudwatchlogs

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
)

// The field index policy record family: one policy per log group (the
// group-level surface PutIndexPolicy manages), persisted in the JSON
// regime and torn down with the group. The account-level FIELD_INDEX_POLICY
// family is the existing AccountPolicy record set; the effective policy
// resolution reads both.

// indexPolicyRecordMu serialises the two-record policy+trail writes:
// the replacement's INACTIVE-trail merge reads the current policy and
// trail and writes both records under this mutex, so concurrent policy
// changes on one group cannot compute from the same base and drop each
// other's trail entries.
var indexPolicyRecordMu sync.Mutex

// ReplaceIndexPolicy runs the policy+trail replacement as one locked
// section: build sees the current policy (nil when absent) and the
// current trail, and the records it returns persist together. An error
// from build aborts without writing either record.
func (s *Store) ReplaceIndexPolicy(groupName string, build func(current *IndexPolicy, trail *IndexInactiveTrail) (*IndexPolicy, *IndexInactiveTrail, error)) error {
	indexPolicyRecordMu.Lock()
	defer indexPolicyRecordMu.Unlock()
	current, err := s.GetIndexPolicy(groupName)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		// Only a missing record means "no policy": a failed read is not
		// absence, and building from a phantom nil would overwrite the
		// unreadable record's trail with a hole (and the rollback path
		// would then delete it) — the error propagates instead.
		return err
	}
	if err != nil {
		current = nil
	}
	trail := s.GetIndexInactiveTrail(groupName)
	policy, newTrail, err := build(current, trail)
	if err != nil {
		return err
	}
	policy.LastUpdateTime = time.Now().UTC().UnixMilli()
	if s.ts != nil {
		// The pair commits as one storage transaction where the backend
		// supports it: a mid-pair failure leaves neither record written,
		// so the policy and its INACTIVE trail cannot diverge.
		policyData, err := json.Marshal(policy)
		if err != nil {
			return err
		}
		trailData, err := json.Marshal(newTrail)
		if err != nil {
			return err
		}
		ctx := context.Background()
		return s.ts.Update(ctx, func(txn storage.Transaction) error {
			if err := txn.Bucket(s.bucketName).Put([]byte(s.indexPolicyKey(groupName)), policyData); err != nil {
				return err
			}
			trailKey := []byte(s.indexInactiveTrailKey(groupName))
			if len(newTrail.Fields) == 0 {
				// An empty trail is the record's absence (the writer's
				// own convention).
				return txn.Bucket(s.bucketName).Delete(trailKey)
			}
			return txn.Bucket(s.bucketName).Put(trailKey, trailData)
		})
	}
	// Fallback (non-transactional backend): sequential writes with
	// compensation — a trail-write failure restores the pre-replacement
	// policy record (or its absence), so the pair still cannot diverge.
	if err := s.putIndexPolicyStamped(policy); err != nil {
		return err
	}
	if err := s.PutIndexInactiveTrail(newTrail); err != nil {
		if current != nil {
			if rbErr := s.putIndexPolicyStamped(current); rbErr != nil {
				logs.Error("Failed to restore the index policy after a trail-write failure",
					logs.String("logGroupName", groupName), logs.Err(rbErr))
			}
		} else if rbErr := s.deleteIndexPolicyLocked(groupName); rbErr != nil && !errors.Is(rbErr, ErrResourceNotFound) {
			logs.Error("Failed to remove the index policy after a trail-write failure",
				logs.String("logGroupName", groupName), logs.Err(rbErr))
		}
		return err
	}
	return nil
}

// RemoveIndexPolicy runs the delete-path trail update and the policy
// removal as one locked section (the delete's trail merge reads the
// policy it is removing).
func (s *Store) RemoveIndexPolicy(groupName string, build func(current *IndexPolicy, trail *IndexInactiveTrail) (*IndexInactiveTrail, error)) error {
	indexPolicyRecordMu.Lock()
	defer indexPolicyRecordMu.Unlock()
	current, err := s.GetIndexPolicy(groupName)
	if err != nil {
		return err
	}
	trail := s.GetIndexInactiveTrail(groupName)
	newTrail, err := build(current, trail)
	if err != nil {
		return err
	}
	if err := s.deleteIndexPolicyLocked(groupName); err != nil {
		return err
	}
	return s.PutIndexInactiveTrail(newTrail)
}

// putIndexPolicyStamped writes a policy record whose modification stamp
// the caller already set — the replacement's compensation path, which
// must restore the pre-replacement record byte-identically rather than
// re-stamp it.
func (s *Store) putIndexPolicyStamped(p *IndexPolicy) error {
	return s.Put(s.indexPolicyKey(p.LogGroupName), p)
}

// readIndexPolicyRecord is the policy-record read behind GetIndexPolicy,
// indirected so tests can exercise the replacement's read-failure path —
// a failed read must propagate, never read as policy-absence.
var readIndexPolicyRecord = func(s *Store, key string) (*IndexPolicy, error) {
	return getJSONRecord[IndexPolicy](s, key, ErrResourceNotFound)
}

func (s *Store) GetIndexPolicy(logGroupName string) (*IndexPolicy, error) {
	return readIndexPolicyRecord(s, s.indexPolicyKey(logGroupName))
}

// ReadIndexPolicySnapshot returns the group's policy record (nil when
// absent) and INACTIVE trail as one mutex-held pair: every trail write
// happens inside this mutex's locked sections, so a replacement cannot
// commit between the two reads and leave a listing that carries a
// replaced field under neither its active nor its inactive category.
func (s *Store) ReadIndexPolicySnapshot(logGroupName string) (*IndexPolicy, *IndexInactiveTrail) {
	indexPolicyRecordMu.Lock()
	defer indexPolicyRecordMu.Unlock()
	policy, err := s.GetIndexPolicy(logGroupName)
	if err != nil {
		policy = nil
	}
	return policy, s.GetIndexInactiveTrail(logGroupName)
}

// deleteIndexPolicyLocked removes the policy record; the caller holds
// indexPolicyRecordMu (RemoveIndexPolicy's locked section reuses it, so
// the exported delete must not re-enter the mutex).
func (s *Store) deleteIndexPolicyLocked(groupName string) error {
	if !s.Exists(s.indexPolicyKey(groupName)) {
		return ErrResourceNotFound
	}
	return s.Delete(s.indexPolicyKey(groupName))
}

func (s *Store) DeleteIndexPolicy(logGroupName string) error {
	// Under the family mutex so ReplaceIndexPolicy's two-record write
	// cannot re-create the policy after the delete.
	indexPolicyRecordMu.Lock()
	defer indexPolicyRecordMu.Unlock()
	return s.deleteIndexPolicyLocked(logGroupName)
}

// --- INACTIVE trail ---

// PutIndexInactiveTrail replaces a group's INACTIVE field trail.
func (s *Store) PutIndexInactiveTrail(t *IndexInactiveTrail) error {
	if len(t.Fields) == 0 {
		// An empty trail is the absence of the record: dropping the last
		// inactive field removes the trail entirely.
		return s.Delete(s.indexInactiveTrailKey(t.LogGroupName))
	}
	return s.Put(s.indexInactiveTrailKey(t.LogGroupName), t)
}

// GetIndexInactiveTrail loads a group's INACTIVE field trail; a group
// without one reports the empty trail.
func (s *Store) GetIndexInactiveTrail(logGroupName string) *IndexInactiveTrail {
	var t IndexInactiveTrail
	if err := s.Get(s.indexInactiveTrailKey(logGroupName), &t); err != nil {
		return &IndexInactiveTrail{LogGroupName: logGroupName}
	}
	return &t
}

// DeleteIndexInactiveTrail removes the trail (the group teardown path).
func (s *Store) DeleteIndexInactiveTrail(logGroupName string) error {
	// Under the family mutex for the same reason as DeleteIndexPolicy.
	indexPolicyRecordMu.Lock()
	defer indexPolicyRecordMu.Unlock()
	return s.Delete(s.indexInactiveTrailKey(logGroupName))
}

// IndexFieldSpec pairs one indexed field name with its index type (the
// FieldsV2 vocabulary; Fields-array entries are FIELD_INDEX).
type IndexFieldSpec struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// IndexPolicy is one log group's field index policy: the policy document
// verbatim and the ARN echo the responses carry. The INACTIVE field
// trail is the separate IndexInactiveTrail family (it outlives this
// record's deletion).
type IndexPolicy struct {
	LogGroupName       string `json:"logGroupName"`
	LogGroupIdentifier string `json:"logGroupIdentifier"`
	PolicyDocument     string `json:"policyDocument"`
	LastUpdateTime     int64  `json:"lastUpdateTime"`
}

// IndexInactiveTrail records the fields a group previously indexed but
// no longer does — the INACTIVE index category's substrate. It survives
// the policy record's deletion (the trail outlives the policy that
// produced it) and dies with the group.
type IndexInactiveTrail struct {
	LogGroupName string           `json:"logGroupName"`
	Fields       []IndexFieldSpec `json:"fields,omitempty"`
}
