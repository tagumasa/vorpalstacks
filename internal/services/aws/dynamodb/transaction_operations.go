// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash/fnv"
	"sort"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TransactGetItems performs multiple GetItem operations in a single transaction with snapshot isolation.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactGetItems.html
func (s *DynamoDBService) TransactGetItems(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transactItems, _ := req.Parameters["TransactItems"].([]interface{})
	return s.transactGetItemsCore(ctx, reqCtx, transactGetItemsInput{
		TransactItems: transactItems,
		Parameters:    req.Parameters,
	})
}

// TransactWriteItems performs multiple write operations in a single transaction with ACID semantics.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html
func (s *DynamoDBService) TransactWriteItems(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transactItems, _ := req.Parameters["TransactItems"].([]interface{})
	return s.transactWriteItemsCore(ctx, reqCtx, transactWriteItemsInput{
		TransactItems: transactItems,
		Parameters:    req.Parameters,
	})
}

// copyAttributes deep-copies an item's attribute map. The update appliers'
// evaluation basis is the item as it was before the expression: the copy
// must not alias any nested map or list whose members an action's write
// replaces or extends, or a later operand's read would observe the
// expression's own writes.
func copyAttributes(attrs map[string]*dbstore.AttributeValue) map[string]*dbstore.AttributeValue {
	if attrs == nil {
		return nil
	}
	cpy := make(map[string]*dbstore.AttributeValue)
	for k, v := range attrs {
		cpy[k] = deepCopyAttributeValue(v)
	}
	return cpy
}

// deepCopyAttributeValue copies one attribute value, recursing into the Map
// and List containers and duplicating every slice and pointed-to scalar.
func deepCopyAttributeValue(v *dbstore.AttributeValue) *dbstore.AttributeValue {
	if v == nil {
		return nil
	}
	cpy := &dbstore.AttributeValue{}
	if v.S != nil {
		s := *v.S
		cpy.S = &s
	}
	if v.N != nil {
		n := *v.N
		cpy.N = &n
	}
	if v.B != nil {
		cpy.B = make([]byte, len(v.B))
		copy(cpy.B, v.B)
	}
	if v.SS != nil {
		cpy.SS = make([]string, len(v.SS))
		copy(cpy.SS, v.SS)
	}
	if v.NS != nil {
		cpy.NS = make([]string, len(v.NS))
		copy(cpy.NS, v.NS)
	}
	if v.BS != nil {
		cpy.BS = make([][]byte, len(v.BS))
		for i, b := range v.BS {
			cpy.BS[i] = make([]byte, len(b))
			copy(cpy.BS[i], b)
		}
	}
	if v.M != nil {
		cpy.M = make(map[string]*dbstore.AttributeValue)
		for k, val := range v.M {
			cpy.M[k] = deepCopyAttributeValue(val)
		}
	}
	if v.L != nil {
		cpy.L = make([]*dbstore.AttributeValue, len(v.L))
		for i, val := range v.L {
			cpy.L[i] = deepCopyAttributeValue(val)
		}
	}
	if v.NULL != nil {
		null := *v.NULL
		cpy.NULL = &null
	}
	if v.BOOL != nil {
		b := *v.BOOL
		cpy.BOOL = &b
	}
	return cpy
}

// buildKeyString renders a table key map as a deterministic identity string
// for TransactWriteItems conflict and idempotency detection. Attribute
// names and values are rendered with the store's key component encoding,
// which is injective over all legal name and value bytes.
func buildKeyString(tableName string, key map[string]*dbstore.AttributeValue) string {
	names := make([]string, 0, len(key))
	for k := range key {
		names = append(names, k)
	}
	sort.Strings(names)

	result := tableName + "#"
	for _, k := range names {
		v := key[k]
		result += dbstore.EncodeKeyValue(&dbstore.AttributeValue{S: &k}) + "=" + dbstore.EncodeKeyValue(v) + ";"
	}
	return result
}

// hashTransactWriteRequest derives a stable digest of the request payload
// (everything except the client token itself) so a token replayed with a
// different payload can be detected within the idempotency window.
func hashTransactWriteRequest(params map[string]interface{}) string {
	payload := make(map[string]interface{}, len(params))
	for k, v := range params {
		if k == "ClientRequestToken" {
			continue
		}
		payload[k] = v
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// idempotencySweepInterval is how often expired client request tokens are
// removed from the per-region idempotency buckets.
const idempotencySweepInterval = time.Minute

// tokenClaimGuard tracks one transaction's idempotency-token claim. The
// deferred release drops the in-progress record whenever the execution
// does not reach its committed path, so no early return — on either
// transaction core — can leave a claimed token locking the client out of
// the idempotency window. Marking the guard committed stops the release:
// after the storage commit the record belongs to the completion path,
// and deleting it would let a retry re-execute the transaction.
type tokenClaimGuard struct {
	claimed   bool
	committed bool
	key       string
	store     dbstore.DynamoDBStoreInterface
}

// claim records a successfully claimed token under its store key.
func (g *tokenClaimGuard) claim(key string) {
	g.key = key
	g.claimed = true
}

// markCommitted moves the guard past the release point — the claim
// survives as the completion path's record.
func (g *tokenClaimGuard) markCommitted() {
	g.committed = true
}

// release drops the in-progress record after a failed execution so the
// client can retry the token; a claim never made and a committed
// transaction are both no-ops.
func (g *tokenClaimGuard) release() {
	if !g.claimed || g.committed {
		return
	}
	if delErr := g.store.Idempotency().Delete(g.key); delErr != nil {
		logs.Error("Failed to release idempotency token claim", logs.Err(delErr))
	}
}

// clientRequestTokenLockShards shards the per-token claim locks: tokens are
// unique per request, so keyed locks would grow without bound, while a
// fixed shard count still serialises every caller of one token (the same
// token always maps to the same shard).
const clientRequestTokenLockShards = 64

// lockClientRequestToken serialises the idempotency claim section for one
// client request token and returns the unlock function. Unrelated tokens
// that share a shard also serialise briefly; the claim section is a single
// store read and write, so the contention is immaterial.
func (s *DynamoDBService) lockClientRequestToken(token string) func() {
	h := fnv.New32a()
	_, _ = h.Write([]byte(token))
	mu := &s.clientRequestTokenLocks[h.Sum32()%clientRequestTokenLockShards]
	mu.Lock()
	return mu.Unlock
}

// sweepStoreIdempotency removes the expired client request tokens of one
// regional store.
func (s *DynamoDBService) sweepStoreIdempotency(store dbstore.DynamoDBStoreInterface) {
	if _, err := store.Idempotency().SweepExpired(time.Now()); err != nil {
		logs.Error("Failed to sweep expired idempotency tokens", logs.Err(err))
	}
}

// ensureIdempotencySweeper starts the background sweeper that removes
// transaction client request tokens once their idempotency window has
// lapsed, across every region store cached by the service.
func (s *DynamoDBService) ensureIdempotencySweeper() {
	s.startIntervalSweeper(&s.idempotencySweepOnce, idempotencySweepInterval, "dynamodb idempotency sweep", s.sweepStoreIdempotency)
}
