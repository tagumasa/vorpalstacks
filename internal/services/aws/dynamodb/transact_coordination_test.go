// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestTransactWriteSurvivesCompletionRecordFailure(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnIdemTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	requestParams := func(token string) map[string]interface{} {
		params := map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{
				"TableName": "TxnIdemTable",
				"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("committed")},
			}}},
		}
		if token != "" {
			params["ClientRequestToken"] = token
		}
		return params
	}
	restoreSeam := func() {
		recordIdempotencyCompletion = func(store dbstore.DynamoDBStoreInterface, token, requestHash string, expiresAt time.Time, readUnits map[string]float64) error {
			return store.Idempotency().Record(token, requestHash, dbstore.IdempotencyStateCompleted, expiresAt, readUnits)
		}
	}

	// Fault injection: the completion record fails for the first token.
	recordIdempotencyCompletion = func(dbstore.DynamoDBStoreInterface, string, string, time.Time, map[string]float64) error {
		return errors.New("injected completion-record failure")
	}
	defer restoreSeam()

	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: requestParams("tok-fail")}); err != nil {
		t.Fatalf("completion-record failure: the committed transaction must still succeed, got %v", err)
	}
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnIdemTable", "Key": map[string]interface{}{"id": sVal("a")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if got := resp.(map[string]interface{})["Item"].(map[string]interface{})["v"].(map[string]interface{})["S"]; got != "committed" {
		t.Fatalf("completion-record failure: expected the written v=committed, got %v", got)
	}

	// The token stays claimed in-progress: an in-window retry answers the
	// documented TransactionInProgress — bounded, and never a re-execution.
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: requestParams("tok-fail")}); !errors.Is(err, ErrTransactionInProgress) {
		t.Fatalf("in-window retry after failed completion record: expected TransactionInProgress, got %v", err)
	}

	// With the seam restored, the completion record works and a replay of
	// the same token and payload answers from the record.
	restoreSeam()
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: requestParams("tok-ok")}); err != nil {
		t.Fatalf("first token-ok call: unexpected error %v", err)
	}
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: requestParams("tok-ok")}); err != nil {
		t.Fatalf("token-ok replay: expected the recorded outcome, got %v", err)
	}
}

func TestTransactItemConflictDetection(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TxnConflictTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	keyA := map[string]interface{}{"id": sVal("a")}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnConflictTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	lockKeyA := itemLockKey("us-east-1", "TxnConflictTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("a")})
	txnPutA := func() error {
		_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{
				"TableName": "TxnConflictTable",
				"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("y")},
			}}},
		}})
		return err
	}
	tgiGetA := func() error {
		_, err := svc.TransactGetItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Get": map[string]interface{}{
				"TableName": "TxnConflictTable", "Key": keyA,
			}}},
		}})
		return err
	}
	canceledWithConflictReason := func(err error, label string) {
		t.Helper()
		var canceled *TransactionCanceledError
		if err == nil || !errors.As(err, &canceled) {
			t.Fatalf("%s: expected TransactionCanceledException, got %v", label, err)
		}
		if len(canceled.CancellationReasons) != 1 || canceled.CancellationReasons[0].Code != "TransactionConflict" {
			t.Fatalf("%s: expected a TransactionConflict cancellation reason, got %+v", label, canceled.CancellationReasons)
		}
		if got := canceled.CancellationReasons[0].Message; got != transactionConflictReasonMessage {
			t.Fatalf("%s: expected the documented conflict reason message, got %+v", label, canceled.CancellationReasons[0])
		}
	}

	// An in-flight TransactWriteItems on the item: the incoming transaction
	// planes and every single-item write reject; other items and plain
	// reads proceed.
	if _, free := tryLockItems(itemLockModeTxn, []string{lockKeyA}); !free {
		t.Fatal("txn lock acquisition: expected the uncontended key to be free")
	}
	canceledWithConflictReason(txnPutA(), "transact write over in-flight transaction")
	canceledWithConflictReason(tgiGetA(), "transact get over in-flight transaction")
	for label, call := range map[string]func() error{
		"put over in-flight transaction": func() error {
			_, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "TxnConflictTable", "Item": map[string]interface{}{"id": sVal("a"), "v": sVal("z")},
			}})
			return err
		},
		"update over in-flight transaction": func() error {
			_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "TxnConflictTable", "Key": keyA,
				"UpdateExpression":          "SET v = :v",
				"ExpressionAttributeValues": map[string]interface{}{":v": sVal("z")},
			}})
			return err
		},
		"delete over in-flight transaction": func() error {
			_, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "TxnConflictTable", "Key": keyA,
			}})
			return err
		},
	} {
		if err := call(); !errors.Is(err, ErrTransactionConflict) {
			t.Fatalf("%s: expected TransactionConflictException, got %v", label, err)
		}
	}
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnConflictTable", "Item": map[string]interface{}{"id": sVal("b"), "v": sVal("ok")},
	}}); err != nil {
		t.Fatalf("other item over in-flight transaction: unexpected error %v", err)
	}
	if _, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnConflictTable", "Key": keyA,
	}}); err != nil {
		t.Fatalf("plain read over in-flight transaction: unexpected error %v", err)
	}
	unlockItems(itemLockModeTxn, []string{lockKeyA})

	// An in-flight single-item write: the transaction planes still reject
	// (the serialisation runs both directions), while a concurrent
	// single-item write stays legal last-writer-wins.
	if _, free := tryLockItems(itemLockModeItem, []string{lockKeyA}); !free {
		t.Fatal("item lock acquisition: expected the released key to be free")
	}
	canceledWithConflictReason(txnPutA(), "transact write over in-flight single-item write")
	canceledWithConflictReason(tgiGetA(), "transact get over in-flight single-item write")
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TxnConflictTable", "Item": map[string]interface{}{"id": sVal("a"), "v": sVal("wins")},
	}}); err != nil {
		t.Fatalf("single-item write over in-flight single-item write: unexpected error %v", err)
	}
	unlockItems(itemLockModeItem, []string{lockKeyA})

	// Two concurrent single-item writers share one key: the first to
	// finish releases only its own share, and a transaction acquiring in
	// that window still meets the conflict the registry exists to raise —
	// the early full release is the lost-update window.
	if _, free := tryLockItems(itemLockModeItem, []string{lockKeyA}); !free {
		t.Fatal("first co-writer acquisition: expected the key to be free")
	}
	if _, free := tryLockItems(itemLockModeItem, []string{lockKeyA}); !free {
		t.Fatal("second co-writer acquisition: expected item-mode co-tenancy")
	}
	unlockItems(itemLockModeItem, []string{lockKeyA})
	canceledWithConflictReason(txnPutA(), "transact write over the surviving co-writer")
	unlockItems(itemLockModeItem, []string{lockKeyA})
	if _, locked := firstLockedItem([]string{lockKeyA}); locked {
		t.Fatal("fully released key: expected no lock to remain")
	}

	// Uncontended, the transaction commits.
	if err := txnPutA(); err != nil {
		t.Fatalf("uncontended transact write: unexpected error %v", err)
	}
}

func TestExecuteTransactionClientTokenIdempotency(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "PartiQLIdemTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	execute := func(token, statement string) error {
		params := map[string]interface{}{
			"TransactStatements": []interface{}{map[string]interface{}{"Statement": statement}},
		}
		if token != "" {
			params["ClientRequestToken"] = token
		}
		_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}

	insert := `INSERT INTO "PartiQLIdemTable" VALUE {'id': 'a', 'v': '1'}`
	if err := execute("etok", insert); err != nil {
		t.Fatalf("first token call: unexpected error %v", err)
	}

	// The same token and payload replays: an executed INSERT onto the
	// existing item would cancel the transaction, so success proves the
	// write did not re-run.
	if err := execute("etok", insert); err != nil {
		t.Fatalf("token replay: expected the recorded outcome, got %v", err)
	}
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "PartiQLIdemTable", "Key": map[string]interface{}{"id": sVal("a")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	if _, present := resp.(map[string]interface{})["Item"]; !present {
		t.Fatal("token replay: expected the inserted item to remain")
	}

	// The same token with a different payload is rejected.
	if err := execute("etok", `INSERT INTO "PartiQLIdemTable" VALUE {'id': 'b', 'v': '2'}`); !errors.Is(err, ErrIdempotentParameterMismatch) {
		t.Fatalf("token payload mismatch: expected IdempotentParameterMismatch, got %v", err)
	}

	// An in-progress record answers the documented in-progress error; the
	// record's hash must match the payload, or the mismatch fires first.
	params2 := map[string]interface{}{
		"TransactStatements": []interface{}{map[string]interface{}{"Statement": `INSERT INTO "PartiQLIdemTable" VALUE {'id': 'c', 'v': '3'}`}},
		"ClientRequestToken": "etok2",
	}
	if err := store.Idempotency().Record("ExecuteTransaction|etok2", hashTransactWriteRequest(params2),
		dbstore.IdempotencyStateInProgress, time.Now().Add(5*time.Minute), nil); err != nil {
		t.Fatalf("seed in-progress record: %v", err)
	}
	if err := execute("etok2", `INSERT INTO "PartiQLIdemTable" VALUE {'id': 'c', 'v': '3'}`); !errors.Is(err, ErrTransactionInProgress) {
		t.Fatalf("in-progress token: expected TransactionInProgress, got %v", err)
	}
	if _, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "PartiQLIdemTable", "Key": map[string]interface{}{"id": sVal("c")},
	}}); err != nil {
		t.Fatalf("missing-target get: %v", err)
	}

	// The token keyspace is namespaced per API: the same token string on
	// TransactWriteItems addresses an independent record and executes.
	if _, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{
			"TableName": "PartiQLIdemTable",
			"Item":      map[string]interface{}{"id": sVal("d"), "v": sVal("x")},
		}}},
		"ClientRequestToken": "etok",
	}}); err != nil {
		t.Fatalf("cross-API token reuse: unexpected error %v", err)
	}

	// A cancelled-before-execution transaction returns its claimed token:
	// an item-registry conflict rejects the write before it starts, and a
	// retry of the same token after the conflicting holder releases must
	// execute rather than meet a lingering in-progress record.
	conflictVal := "e"
	conflictKey := itemLockKey("us-east-1", "PartiQLIdemTable", map[string]*dbstore.AttributeValue{
		"id": {S: &conflictVal},
	})
	if _, free := tryLockItems(itemLockModeTxn, []string{conflictKey}); !free {
		t.Fatal("seed conflicting item lock")
	}
	conflictCall := func() error {
		_, err := svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactItems": []interface{}{map[string]interface{}{"Put": map[string]interface{}{
				"TableName": "PartiQLIdemTable",
				"Item":      map[string]interface{}{"id": sVal("e"), "v": sVal("y")},
			}}},
			"ClientRequestToken": "etok3",
		}})
		return err
	}
	if err := conflictCall(); err == nil {
		t.Fatal("conflicting holder: expected the transaction to cancel")
	}
	unlockItems(itemLockModeTxn, []string{conflictKey})
	if err := conflictCall(); err != nil {
		t.Fatalf("token retry after a cancelled transaction: %v", err)
	}
}

func TestExecuteTransactionItemLockCoordination(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "ExecTxnLockTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "ExecTxnLockTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	execTxn := func(stmt, token string) error {
		m := map[string]interface{}{"Statement": stmt}
		params := map[string]interface{}{"TransactStatements": []interface{}{m}}
		if token != "" {
			params["ClientRequestToken"] = token
		}
		_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	key := func(id string) string {
		return itemLockKey("us-east-1", "ExecTxnLockTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue(id)})
	}
	expectConflict := func(err error, label string) {
		t.Helper()
		// The operation's declared error set carries
		// TransactionCanceledException alone, so a lock conflict answers
		// through the cancellation envelope with the documented
		// TransactionConflict reason at the conflicting statement's slot.
		var canceled *TransactionCanceledError
		if !errors.As(err, &canceled) {
			t.Fatalf("%s: expected TransactionCanceledException, got %v", label, err)
		}
		if len(canceled.CancellationReasons) != 1 ||
			canceled.CancellationReasons[0].Code != "TransactionConflict" ||
			canceled.CancellationReasons[0].Message != transactionConflictReasonMessage {
			t.Fatalf("%s: expected a TransactionConflict cancellation reason, got %+v", label, canceled.CancellationReasons)
		}
	}

	// A transaction-mode holder on the item: the ExecuteTransaction UPDATE
	// on the same item rejects.
	lockA := key("a")
	if _, ok := tryLockItems(itemLockModeTxn, []string{lockA}); !ok {
		t.Fatal("txn-mode lock acquisition: unexpectedly conflicted")
	}
	expectConflict(execTxn(`UPDATE "ExecTxnLockTable" SET v = 'y' WHERE id = 'a'`, ""), "txn-mode holder UPDATE")

	// The rejection released the claimed idempotency token: the same
	// ClientRequestToken retries once the holder is gone instead of
	// answering TransactionInProgress for the window.
	expectConflict(execTxn(`UPDATE "ExecTxnLockTable" SET v = 'y' WHERE id = 'a'`, "coord-token"), "tokened rejection")
	unlockItems(itemLockModeTxn, []string{lockA})
	if err := execTxn(`UPDATE "ExecTxnLockTable" SET v = 'y' WHERE id = 'a'`, "coord-token"); err != nil {
		t.Fatalf("token retry after released holder: %v", err)
	}

	// An item-mode holder (a single-item write in flight): the transaction
	// acquisition conflicts with any holder.
	if _, ok := tryLockItems(itemLockModeItem, []string{lockA}); !ok {
		t.Fatal("item-mode lock acquisition: unexpectedly conflicted")
	}
	expectConflict(execTxn(`DELETE FROM "ExecTxnLockTable" WHERE id = 'a'`, ""), "item-mode holder DELETE")
	unlockItems(itemLockModeItem, []string{lockA})

	// An INSERT locks its target key the same way.
	lockB := key("b")
	if _, ok := tryLockItems(itemLockModeTxn, []string{lockB}); !ok {
		t.Fatal("insert-target lock acquisition: unexpectedly conflicted")
	}
	expectConflict(execTxn(`INSERT INTO "ExecTxnLockTable" VALUE {'id': 'b', 'v': 'z'}`, ""), "txn-mode holder INSERT")
	unlockItems(itemLockModeTxn, []string{lockB})

	// The locks leave with the response: after a committed ExecuteTransaction
	// the same item's key acquires freely.
	if err := execTxn(`INSERT INTO "ExecTxnLockTable" VALUE {'id': 'b', 'v': 'z'}`, ""); err != nil {
		t.Fatalf("post-holder INSERT: %v", err)
	}
	if _, ok := tryLockItems(itemLockModeTxn, []string{lockB}); !ok {
		t.Fatal("lock released: the committed ExecuteTransaction left its lock behind")
	}
	unlockItems(itemLockModeTxn, []string{lockB})
}

// TestExecuteTransactionReleasesTokenOnDerivationFailures pins the claim
// guard: a request-level failure inside the lock-key derivation releases
// the claimed token, so the same token retried answers the same
// request-level error instead of the ten-minute TransactionInProgress
// lockout. The legs cover both derivation branches (INSERT and
// UPDATE/DELETE) under both failure shapes (missing table, non-ACTIVE
// table).
func TestExecuteTransactionReleasesTokenOnDerivationFailures(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "ExecTxnTokenTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	creating, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "ExecTxnCreatingTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	if err != nil {
		t.Fatalf("create CREATING table: %v", err)
	}
	creating.Status = dbstore.TableStatusCreating
	if err := store.Tables().Put(creating); err != nil {
		t.Fatalf("persist CREATING status: %v", err)
	}

	exec := func(stmt, token string) error {
		_, err := svc.ExecuteTransaction(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TransactStatements": []interface{}{map[string]interface{}{"Statement": stmt}},
			"ClientRequestToken": token,
		}})
		return err
	}

	legs := []struct{ name, token, stmt string }{
		{"INSERT missing table", "deriv-miss-i", `INSERT INTO "ExecTxnNoTable" VALUE {'id': 'a'}`},
		{"UPDATE missing table", "deriv-miss-u", `UPDATE "ExecTxnNoTable" SET v = 'x' WHERE id = 'a'`},
		{"DELETE missing table", "deriv-miss-d", `DELETE FROM "ExecTxnNoTable" WHERE id = 'a'`},
		{"INSERT non-ACTIVE table", "deriv-creating-i", `INSERT INTO "ExecTxnCreatingTable" VALUE {'id': 'a'}`},
		{"UPDATE non-ACTIVE table", "deriv-creating-u", `UPDATE "ExecTxnCreatingTable" SET v = 'x' WHERE id = 'a'`},
		{"DELETE non-ACTIVE table", "deriv-creating-d", `DELETE FROM "ExecTxnCreatingTable" WHERE id = 'a'`},
	}
	for _, leg := range legs {
		first := exec(leg.stmt, leg.token)
		if first == nil {
			t.Fatalf("%s: expected the request-level failure, got success", leg.name)
		}
		retry := exec(leg.stmt, leg.token)
		if errors.Is(retry, ErrTransactionInProgress) {
			t.Fatalf("%s: retry answered TransactionInProgress — the claimed token leaked", leg.name)
		}
		if retry.Error() != first.Error() {
			t.Fatalf("%s: retry answered %v, want the first call's %v", leg.name, retry, first)
		}
	}
}

// TestItemLockConflictSelectionIsDeterministic pins the ordered conflict
// scan: when several requested keys are held, the reported conflicting
// key — and the cancellation slot it maps to — is the first one in
// request order, on both the acquire path and the read-only check.
func TestItemLockConflictSelectionIsDeterministic(t *testing.T) {
	kFirst, kSecond := "det|T|first", "det|T|second"
	keys := []string{kFirst, kSecond}

	if _, free := tryLockItems(itemLockModeTxn, keys); !free {
		t.Fatal("holder acquisition: expected both keys free")
	}
	if k, free := tryLockItems(itemLockModeTxn, keys); free {
		t.Fatal("both-held probe: expected a conflict")
	} else if k != kFirst {
		t.Fatalf("both-held probe: expected the first key in request order, got %q", k)
	}
	unlockItems(itemLockModeTxn, []string{kFirst})
	if k, free := tryLockItems(itemLockModeTxn, keys); free {
		t.Fatal("second-held probe: expected a conflict")
	} else if k != kSecond {
		t.Fatalf("second-held probe: expected the remaining held key, got %q", k)
	}

	// The read plane's check selects through the same ordered scan.
	if k, locked := firstLockedItem(keys); !locked || k != kSecond {
		t.Fatalf("read-plane probe: expected the remaining held key, got (%q, %v)", k, locked)
	}
	unlockItems(itemLockModeTxn, []string{kSecond})
}

func TestStandaloneWritePlanesCoordinateWithTheItemRegistry(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "LockPlaneTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LockPlaneTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "v": sVal("x")},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	exec := func(stmt string) error {
		_, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"Statement": stmt}})
		return err
	}
	batchPut := func(id, v string) map[string]interface{} {
		t.Helper()
		resp, err := svc.BatchWriteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"RequestItems": map[string]interface{}{
				"LockPlaneTable": []interface{}{
					map[string]interface{}{"PutRequest": map[string]interface{}{"Item": map[string]interface{}{"id": sVal(id), "v": sVal(v)}}},
				},
			},
		}})
		if err != nil {
			t.Fatalf("batch put %q: %v", id, err)
		}
		return resp.(map[string]interface{})
	}

	// An in-flight transaction holder on the item: every standalone write
	// verb addressing it rejects with TransactionConflictException.
	lockA := itemLockKey("us-east-1", "LockPlaneTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("a")})
	if _, ok := tryLockItems(itemLockModeTxn, []string{lockA}); !ok {
		t.Fatal("txn-mode lock acquisition: unexpectedly conflicted")
	}
	for _, stmt := range []string{
		`UPDATE "LockPlaneTable" SET v = 'y' WHERE id = 'a'`,
		`DELETE FROM "LockPlaneTable" WHERE id = 'a'`,
		`INSERT INTO "LockPlaneTable" VALUE {'id': 'a', 'v': 'z'}`,
	} {
		if err := exec(stmt); !errors.Is(err, ErrTransactionConflict) {
			t.Fatalf("locked %q: expected TransactionConflictException, got %v", stmt, err)
		}
	}

	// A batch op over the holder reports its key through UnprocessedItems
	// while the batch itself succeeds, and an unlocked key in the same
	// plane still commits.
	lockedResp := batchPut("a", "w")
	if unprocessed := lockedResp["UnprocessedItems"].(map[string]interface{}); len(unprocessed) != 1 {
		t.Fatalf("locked batch: expected the key in UnprocessedItems, got %#v", unprocessed)
	}
	freeResp := batchPut("b", "free")
	if unprocessed := freeResp["UnprocessedItems"].(map[string]interface{}); len(unprocessed) != 0 {
		t.Fatalf("unlocked batch: expected no UnprocessedItems, got %#v", unprocessed)
	}

	// After the holder releases, the same statements commit and the locks
	// leave with the responses.
	unlockItems(itemLockModeTxn, []string{lockA})
	if err := exec(`UPDATE "LockPlaneTable" SET v = 'y' WHERE id = 'a'`); err != nil {
		t.Fatalf("post-holder UPDATE: %v", err)
	}
	item, err := store.Items().Get("LockPlaneTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("a")})
	if err != nil {
		t.Fatalf("post-holder read: %v", err)
	}
	if got := *item.Attributes["v"].S; got != "y" {
		t.Fatalf("post-holder UPDATE: expected v=y, got %q", got)
	}
	batchPut("a", "w")
	if item, err = store.Items().Get("LockPlaneTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("a")}); err != nil {
		t.Fatalf("post-holder batch read: %v", err)
	}
	if got := *item.Attributes["v"].S; got != "w" {
		t.Fatalf("post-holder batch put: expected v=w, got %q", got)
	}
	if _, ok := tryLockItems(itemLockModeTxn, []string{lockA}); !ok {
		t.Fatal("lock released: a committed standalone write left its lock behind")
	}
	unlockItems(itemLockModeTxn, []string{lockA})
}
