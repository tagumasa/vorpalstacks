// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"hash/fnv"
	"sort"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TransactGetItems performs multiple GetItem operations in a single transaction with snapshot isolation.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactGetItems.html
func (s *DynamoDBService) TransactGetItems(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transactItems, ok := req.Parameters["TransactItems"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	return s.transactGetItemsCore(ctx, reqCtx, transactGetItemsInput{
		TransactItems: transactItems,
		Parameters:    req.Parameters,
	})
}

// TransactWriteItems performs multiple write operations in a single transaction with ACID semantics.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html
func (s *DynamoDBService) TransactWriteItems(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transactItems, ok := req.Parameters["TransactItems"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	return s.transactWriteItemsCore(ctx, reqCtx, transactWriteItemsInput{
		TransactItems: transactItems,
		Parameters:    req.Parameters,
	})
}

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

func buildKeyString(tableName string, key map[string]*dbstore.AttributeValue) string {
	names := make([]string, 0, len(key))
	for k := range key {
		names = append(names, k)
	}
	sort.Strings(names)

	result := tableName + "#"
	for _, k := range names {
		v := key[k]
		result += k + "="
		if v.S != nil {
			result += "S:" + escapeKeyPart(*v.S)
		} else if v.N != nil {
			result += "N:" + *v.N
		} else if v.B != nil {
			result += "B:" + base64.StdEncoding.EncodeToString(v.B)
		}
		result += ";"
	}
	return result
}

func escapeKeyPart(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '\\', ';', '=', '#':
			result += "\\" + string(c)
		default:
			result += string(c)
		}
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

// ensureIdempotencySweeper starts the background sweeper that removes
// transaction client request tokens once their idempotency window has
// lapsed, across every region store cached by the service.
func (s *DynamoDBService) ensureIdempotencySweeper() {
	s.idempotencySweepOnce.Do(func() {
		s.bgWg.Add(1)
		go func() {
			defer func() { resilience.RecoverPanic("dynamodb idempotency sweep") }()
			defer s.bgWg.Done()
			ticker := time.NewTicker(idempotencySweepInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					s.stores.Range(func(_, v any) bool {
						store, ok := v.(dbstore.DynamoDBStoreInterface)
						if !ok {
							return true
						}
						if _, err := store.Idempotency().SweepExpired(time.Now()); err != nil {
							logs.Error("Failed to sweep expired idempotency tokens", logs.Err(err))
						}
						return true
					})
				case <-s.bgCtx.Done():
					return
				}
			}
		}()
	})
}
