package dynamodb

import (
	"sync"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// Lock modes of the item registry: a transaction write lock is held by an
// in-flight TransactWriteItems for its whole two-phase commit; a
// single-item write lock is held by one PutItem, UpdateItem or DeleteItem
// Core for its read-modify-write transaction.
const (
	itemLockModeTxn  = "txn"
	itemLockModeItem = "item"
)

// itemLockRegistry coordinates the in-flight write planes around one item.
// The storage engine's read-write transactions read committed state and
// commit as last-writer-wins batches, so without cross-transaction
// coordination two transactions addressing one item both commit and the
// earlier write is silently lost. The registry closes that window
// in-process; acquisition is all-or-nothing under one mutex with no
// waiting, so no deadlock is possible and a conflict answers immediately.
//
// The mode matrix follows the documented conflict outcomes (Transaction
// conflict handling, developer guide):
//
//   - an item-level request of TransactWriteItems, or the item's read
//     counterpart in TransactGetItems, rejected over any held lock fails
//     the whole request with TransactionCanceledException carrying a
//     TransactionConflict cancellation reason at the rejected item's slot;
//   - a PutItem, UpdateItem or DeleteItem rejected over an in-flight
//     TransactWriteItems fails with TransactionConflictException, while
//     two concurrent single-item writes on one item remain legal
//     last-writer-wins and never conflict;
//   - a TransactGetItems also answers a conflict when the item is part of
//     an ongoing single-item write.
//
// Reads outside the transaction APIs (GetItem, Query, Scan, BatchGetItem)
// never conflict: the guide documents that a concurrent GetItem and an
// ongoing transaction both succeed, and nothing is rejected on account of
// an ongoing read.
//
// A key therefore carries either one transaction holder or a COUNT of
// single-item holders: several single-item writers may share a key, and
// the key must stay locked until the LAST of them finishes — an early
// releaser would open the lost-update window the registry exists to close
// for a transaction acquiring in between.
var itemLockRegistry = struct {
	sync.Mutex
	held map[string]*itemLockHolders
}{held: make(map[string]*itemLockHolders)}

// itemLockHolders counts the writers one key carries: at most one
// transaction, or any number of concurrent single-item writers.
type itemLockHolders struct {
	txn   bool
	items int
}

// heldAny reports whether the holder record keeps the key locked.
func (h *itemLockHolders) heldAny() bool {
	return h != nil && (h.txn || h.items > 0)
}

// itemLockKey namespaces an item for the registry: locks are scoped per
// region and table, so a same-named table in another region never
// conflicts.
func itemLockKey(region, tableName string, key map[string]*dbstore.AttributeValue) string {
	return region + "|" + tableName + "|" + buildKeyString(tableName, key)
}

// firstConflictingKey scans the deduplicated keys in request order and
// returns the first whose holder matches the mode's conflict rule — the
// deterministic conflict-selection core both lock planes share. The
// reported key, and the cancellation slot it maps to, is always the
// first conflicting key in request order, never a map-iteration pick.
// Caller holds the registry mutex.
func firstConflictingKey(keys []string, conflicts func(*itemLockHolders) bool) (string, bool) {
	seen := make(map[string]bool, len(keys))
	for _, k := range keys {
		if seen[k] {
			continue
		}
		seen[k] = true
		if conflicts(itemLockRegistry.held[k]) {
			return k, true
		}
	}
	return "", false
}

// tryLockItems acquires the write locks for every key in one
// all-or-nothing step, reporting the first conflicting key. A transaction
// acquisition conflicts with any holder; a single-item acquisition
// conflicts only with a transaction holder, because single-item writes
// are last-writer-wins among themselves — the co-holders are counted so
// each writer releases only its own share.
func tryLockItems(mode string, keys []string) (string, bool) {
	itemLockRegistry.Lock()
	defer itemLockRegistry.Unlock()

	conflictKey, conflict := firstConflictingKey(keys, func(h *itemLockHolders) bool {
		return h.heldAny() && (mode != itemLockModeItem || h.txn)
	})
	if conflict {
		return conflictKey, false
	}
	unique := make(map[string]bool, len(keys))
	for _, k := range keys {
		unique[k] = true
	}
	for k := range unique {
		holder := itemLockRegistry.held[k]
		if holder == nil {
			holder = &itemLockHolders{}
			itemLockRegistry.held[k] = holder
		}
		if mode == itemLockModeTxn {
			holder.txn = true
		} else {
			holder.items++
		}
	}
	return "", true
}

// firstLockedItem reports the first key carrying any held write lock: the
// read plane's conflict check. It marks nothing — an ongoing
// TransactGetItems never rejects another request.
func firstLockedItem(keys []string) (string, bool) {
	itemLockRegistry.Lock()
	defer itemLockRegistry.Unlock()

	return firstConflictingKey(keys, (*itemLockHolders).heldAny)
}

// unlockItems releases the acquiring mode's hold on previously acquired
// keys. A single-item release drops one share of a possibly shared key;
// the key leaves the registry only when its last holder is gone.
func unlockItems(mode string, keys []string) {
	itemLockRegistry.Lock()
	defer itemLockRegistry.Unlock()

	unique := make(map[string]bool, len(keys))
	for _, k := range keys {
		unique[k] = true
	}
	for k := range unique {
		holder := itemLockRegistry.held[k]
		if holder == nil {
			continue
		}
		if mode == itemLockModeTxn {
			holder.txn = false
		} else if holder.items > 0 {
			holder.items--
		}
		if !holder.heldAny() {
			delete(itemLockRegistry.held, k)
		}
	}
}
