package dynamodb

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Import / Export Core — single validation + persistence path for table
// import and export operations.
//
// These methods encapsulate import/export lifecycle logic, including S3
// invocations for data transfer. Both the HTTP API handlers
// (import_export_operations.go) and any future admin handler delegate to
// these methods.
// ---------------------------------------------------------------------------
//
// The operation families themselves live in export_core.go and
// import_core.go; this file carries the machinery both families share:
// client-token idempotency lookup, the S3 failure-code mapping, and the
// point-in-time item snapshot reader (which backup restore also consumes).

// clientTokenWindow is the documented idempotency window of an export or
// import ClientToken: a token is valid for eight hours after the first
// request that used it completed.
const clientTokenWindow = 8 * time.Hour

// clientTokenHash derives the idempotency payload hash of an export or
// import request: the canonical JSON of the request parameters with the
// ClientToken itself removed (map marshalling is key-sorted, so the hash is
// stable across identical retries).
func clientTokenHash(params map[string]interface{}) string {
	filtered := make(map[string]interface{}, len(params))
	for k, v := range params {
		if k != "ClientToken" {
			filtered[k] = v
		}
	}
	encoded, err := json.Marshal(filtered)
	if err != nil {
		encoded = []byte(fmt.Sprintf("%v", filtered))
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

// findByClientToken walks one job family's records, page by page, for the
// record a ClientToken created, or nil when none carries it: an idempotent
// retry must find its record however many records the table holds. A failed
// page read yields nil — the retry then creates a replacement record.
func findByClientToken[T any](listPage func(marker string) ([]T, string, error), tokenOf func(T) string, clientToken string) T {
	var zero T
	marker := ""
	for {
		records, next, err := listPage(marker)
		if err != nil {
			return zero
		}
		for _, rec := range records {
			if tokenOf(rec) == clientToken {
				return rec
			}
		}
		if next == "" {
			return zero
		}
		marker = next
	}
}

// s3FailureCode maps an S3-plane error to the failure-code family the
// polling client sees: a not-found store class names the missing resource,
// and everything else keeps the access-class code. The code follows the
// error value, not the call site, so retry logic can tell a vanished
// resource from a permissions problem.
func s3FailureCode(err error) string {
	if commonstore.IsNotFound(err) {
		return "S3NoSuchBucket"
	}
	return "S3AccessDenied"
}

// snapshotItemsAsOf returns the table's items at the given time: the
// current state with every journaled mutation newer than the point undone.
// The undo replays newest first, so each before-image overwrites the state
// the newer mutations left behind.
func snapshotItemsAsOf(store dbstore.DynamoDBStoreInterface, tableName string, at time.Time) ([]*dbstore.Item, error) {
	items := make(map[string]*dbstore.Item)
	order := make([]string, 0)
	if err := store.Items().Scan(tableName, func(item *dbstore.Item) error {
		key := itemKeyString(item.Key)
		if _, seen := items[key]; !seen {
			order = append(order, key)
		}
		items[key] = item
		return nil
	}); err != nil {
		return nil, err
	}

	if err := store.Journal().ReverseReplay(tableName, at, func(change *dbstore.JournalChange) error {
		key := itemKeyString(change.Key)
		if change.BeforeImage == nil {
			delete(items, key)
			return nil
		}
		if _, seen := items[key]; !seen {
			order = append(order, key)
		}
		items[key] = &dbstore.Item{
			TableName:  tableName,
			Key:        change.Key,
			Attributes: change.BeforeImage,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	snapshot := make([]*dbstore.Item, 0, len(order))
	for _, key := range order {
		if item, present := items[key]; present {
			snapshot = append(snapshot, item)
		}
	}
	return snapshot, nil
}

// itemKeyString renders a primary key map as a canonical string. Go's JSON
// encoder writes map keys in sorted order, so equal keys always render
// equally.
func itemKeyString(key map[string]*dbstore.AttributeValue) string {
	data, err := json.Marshal(key)
	if err != nil {
		return ""
	}
	return string(data)
}
