package dynamodb

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

func idempotencyBucketName(region string) string {
	return "dynamodb_idempotency-" + region
}

// IdempotencyStore records TransactWriteItems client request tokens for the
// documented idempotency window so a retry with the same token and payload
// is recognised without re-executing the transaction. A token replayed with
// a different payload is a client error the service layer reports.
type IdempotencyStore struct {
	*common.BaseStore
}

// NewIdempotencyStore creates an IdempotencyStore for the given region.
func NewIdempotencyStore(store storage.BasicStorage, region string) *IdempotencyStore {
	return &IdempotencyStore{
		BaseStore: common.NewBaseStore(store.Bucket(idempotencyBucketName(region)), "dynamodb_idempotency"),
	}
}

// Idempotency record states: a token is claimed as in-progress before the
// transaction executes and promoted to completed once it has committed.
const (
	IdempotencyStateInProgress = "in_progress"
	IdempotencyStateCompleted  = "completed"
)

// idempotencyRecord is the persisted form of one token: the hash of the
// request the token was first used with, the claim state, and when the
// record lapses.
type idempotencyRecord struct {
	RequestHash string `json:"request_hash"`
	State       string `json:"state,omitempty"`
	ExpiresAt   int64  `json:"expires_at"`
}

// Lookup returns the request hash and state recorded for the token when a
// live record exists. Expired records are treated as absent and removed.
func (s *IdempotencyStore) Lookup(token string) (string, string, bool, error) {
	data, err := s.BaseStore.GetRaw(token)
	if err != nil {
		return "", "", false, nil
	}
	var record idempotencyRecord
	if jsonErr := json.Unmarshal(data, &record); jsonErr != nil {
		return "", "", false, nil
	}
	if record.ExpiresAt <= time.Now().Unix() {
		_ = s.BaseStore.Delete(token)
		return "", "", false, nil
	}
	return record.RequestHash, record.State, true, nil
}

// Record stores the token with its request hash, claim state, and expiry.
func (s *IdempotencyStore) Record(token, requestHash, state string, expiresAt time.Time) error {
	record := idempotencyRecord{
		RequestHash: requestHash,
		State:       state,
		ExpiresAt:   expiresAt.Unix(),
	}
	return s.BaseStore.Put(token, record)
}

// SweepExpired removes every record whose idempotency window has lapsed and
// returns how many were removed. Tokens are single-use records with a short
// window, so without a sweep the bucket grows without bound: lazy expiry in
// Lookup only fires when the same token is requested again.
func (s *IdempotencyStore) SweepExpired(now time.Time) (int, error) {
	var expired []string
	if eachErr := s.BaseStore.ForEach(func(token string, data []byte) error {
		var record idempotencyRecord
		if jsonErr := json.Unmarshal(data, &record); jsonErr != nil {
			return nil
		}
		if record.ExpiresAt <= now.Unix() {
			expired = append(expired, token)
		}
		return nil
	}); eachErr != nil {
		return 0, eachErr
	}
	for _, token := range expired {
		if delErr := s.BaseStore.Delete(token); delErr != nil {
			return len(expired), delErr
		}
	}
	return len(expired), nil
}

// Delete removes a token record.
func (s *IdempotencyStore) Delete(token string) error {
	return s.BaseStore.Delete(token)
}
