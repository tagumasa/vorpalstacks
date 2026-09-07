package dynamodb

import (
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	commonstore "vorpalstacks/internal/store/aws/common"
)

func idempotencyBucketName(region string) string {
	return "dynamodb_idempotency-" + region
}

// IdempotencyStore records TransactWriteItems client request tokens for the
// documented idempotency window so a retry with the same token and payload
// is recognised without re-executing the transaction. A token replayed with
// a different payload is a client error the service layer reports.
type IdempotencyStore struct {
	*commonstore.BaseStore
}

// NewIdempotencyStore creates an IdempotencyStore for the given region.
func NewIdempotencyStore(store storage.BasicStorage, region string) *IdempotencyStore {
	return &IdempotencyStore{
		BaseStore: commonstore.NewBaseStore(store.Bucket(idempotencyBucketName(region)), "dynamodb_idempotency"),
	}
}

// Idempotency record states: a token is claimed as in-progress before the
// transaction executes and promoted to completed once it has committed.
const (
	IdempotencyStateInProgress = "in_progress"
	IdempotencyStateCompleted  = "completed"
)

// idempotencyRecord is the in-memory form of one token: the hash of the
// request the token was first used with, the claim state, when the record
// lapses, and — on a completed TransactWriteItems record — the per-table
// read capacity units a same-token replay reports.
type idempotencyRecord struct {
	RequestHash string
	State       string
	ExpiresAt   int64
	ReadUnits   map[string]float64
}

// idempotencyRecordFromBytes decodes a persisted token record; undecodable
// bytes are treated as no record.
func idempotencyRecordFromBytes(data []byte) (*idempotencyRecord, error) {
	var pbRecord pb.IdempotencyRecord
	if err := proto.Unmarshal(data, &pbRecord); err != nil {
		return nil, err
	}
	return &idempotencyRecord{
		RequestHash: pbRecord.RequestHash,
		State:       pbRecord.State,
		ExpiresAt:   pbRecord.ExpiresAt,
		ReadUnits:   pbRecord.ReplayReadUnits,
	}, nil
}

// Lookup returns the request hash, state, and recorded read units for the
// token when a live record exists. A missing record is reported as absent;
// expired records are treated as absent and removed. Any other storage
// failure — and undecodable record bytes — is returned as an error: treating
// corruption as absence would silently re-execute a request the caller
// believes is deduplicated.
func (s *IdempotencyStore) Lookup(token string) (string, string, map[string]float64, bool, error) {
	data, err := s.BaseStore.GetRaw(token)
	if err != nil {
		if commonstore.IsNotFound(err) {
			return "", "", nil, false, nil
		}
		return "", "", nil, false, err
	}
	record, err := idempotencyRecordFromBytes(data)
	if err != nil {
		return "", "", nil, false, err
	}
	if record.ExpiresAt <= time.Now().Unix() {
		_ = s.BaseStore.Delete(token)
		return "", "", nil, false, nil
	}
	return record.RequestHash, record.State, record.ReadUnits, true, nil
}

// Record stores the token with its request hash, claim state, expiry, and —
// for completed TransactWriteItems records — the per-table read units a
// replay reports (nil leaves the record without replay units).
func (s *IdempotencyStore) Record(token, requestHash, state string, expiresAt time.Time, readUnits map[string]float64) error {
	data, err := proto.Marshal(&pb.IdempotencyRecord{
		RequestHash:     requestHash,
		State:           state,
		ExpiresAt:       expiresAt.Unix(),
		ReplayReadUnits: readUnits,
	})
	if err != nil {
		return err
	}
	return s.BaseStore.PutRaw(token, data)
}

// SweepExpired removes every record whose idempotency window has lapsed and
// returns how many were removed. Tokens are single-use records with a short
// window, so without a sweep the bucket grows without bound: lazy expiry in
// Lookup only fires when the same token is requested again.
func (s *IdempotencyStore) SweepExpired(now time.Time) (int, error) {
	var expired []string
	if eachErr := s.BaseStore.ForEach(func(token string, data []byte) error {
		record, err := idempotencyRecordFromBytes(data)
		if err != nil || record.ExpiresAt > now.Unix() {
			return nil
		}
		expired = append(expired, token)
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
