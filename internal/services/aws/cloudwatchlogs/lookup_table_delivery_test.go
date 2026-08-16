package cloudwatchlogs

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// fakeKMSInvoker is a test stand-in for the event bus KMS invoker: data keys
// are 32 random bytes "wrapped" by XOR with a static test key.
type fakeKMSInvoker struct {
	keys    map[string]bool
	wrapKey []byte
}

func newFakeKMSInvoker(keyIDs ...string) *fakeKMSInvoker {
	wrapKey := make([]byte, 32)
	if _, err := rand.Read(wrapKey); err != nil {
		panic(err)
	}
	keys := make(map[string]bool, len(keyIDs))
	for _, id := range keyIDs {
		keys[id] = true
	}
	return &fakeKMSInvoker{keys: keys, wrapKey: wrapKey}
}

func (f *fakeKMSInvoker) GenerateDataKey(_ context.Context, keyID, _ string, _ map[string]string, _ string) (*eventbus.KMSDataKeyResult, error) {
	if !f.keys[keyID] {
		return nil, fmt.Errorf("key %s not found", keyID)
	}
	plaintext := make([]byte, 32)
	if _, err := rand.Read(plaintext); err != nil {
		return nil, err
	}
	wrapped := make([]byte, 32)
	for i := range plaintext {
		wrapped[i] = plaintext[i] ^ f.wrapKey[i]
	}
	return &eventbus.KMSDataKeyResult{Plaintext: plaintext, CiphertextBlob: wrapped}, nil
}

func (f *fakeKMSInvoker) Decrypt(_ context.Context, keyID string, ciphertext []byte, _ map[string]string, _ string) ([]byte, error) {
	if !f.keys[keyID] {
		return nil, fmt.Errorf("key %s not found", keyID)
	}
	plaintext := make([]byte, len(ciphertext))
	for i := range ciphertext {
		plaintext[i] = ciphertext[i] ^ f.wrapKey[i]
	}
	return plaintext, nil
}

func (f *fakeKMSInvoker) KeyExists(_ context.Context, keyID string) bool {
	return f.keys[keyID]
}

func newDeliveryTestService(keyIDs ...string) *LogsService {
	return &LogsService{accountID: "000000000000", kms: newFakeKMSInvoker(keyIDs...)}
}

func newDeliveryTestStore(t *testing.T) *logsstore.Store {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store, err := logsstore.NewStore(st, st.Bucket("logs-us-east-1"), "000000000000", "us-east-1", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestValidateDestinationConfiguration(t *testing.T) {
	validLookup := map[string]interface{}{
		"lookupTableConfiguration": map[string]interface{}{
			"tableName": "users",
			"roleArn":   "arn:aws:iam::000000000000:role/deliver",
		},
	}
	validS3 := map[string]interface{}{
		"s3Configuration": map[string]interface{}{
			"destinationIdentifier": "s3://results-bucket/prefix",
			"roleArn":               "arn:aws:iam::000000000000:role/deliver",
		},
	}
	cases := []struct {
		name string
		dc   map[string]interface{}
		ok   bool
	}{
		{"valid lookup table destination", validLookup, true},
		{"valid s3 destination", validS3, true},
		{"both destinations", map[string]interface{}{
			"lookupTableConfiguration": validLookup["lookupTableConfiguration"],
			"s3Configuration":          validS3["s3Configuration"],
		}, true},
		{"empty configuration", map[string]interface{}{}, true},
		{"lookup missing table name", map[string]interface{}{
			"lookupTableConfiguration": map[string]interface{}{
				"roleArn": "arn:aws:iam::000000000000:role/deliver",
			},
		}, false},
		{"lookup invalid table name", map[string]interface{}{
			"lookupTableConfiguration": map[string]interface{}{
				"tableName": "has-dash",
				"roleArn":   "arn:aws:iam::000000000000:role/deliver",
			},
		}, false},
		{"lookup missing role arn", map[string]interface{}{
			"lookupTableConfiguration": map[string]interface{}{
				"tableName": "users",
			},
		}, false},
		{"lookup non iam role arn", map[string]interface{}{
			"lookupTableConfiguration": map[string]interface{}{
				"tableName": "users",
				"roleArn":   "arn:aws:s3:::bucket",
			},
		}, false},
		{"s3 missing uri", map[string]interface{}{
			"s3Configuration": map[string]interface{}{
				"roleArn": "arn:aws:iam::000000000000:role/deliver",
			},
		}, false},
		{"s3 invalid uri scheme", map[string]interface{}{
			"s3Configuration": map[string]interface{}{
				"destinationIdentifier": "https://bucket/prefix",
				"roleArn":               "arn:aws:iam::000000000000:role/deliver",
			},
		}, false},
		{"s3 invalid owner account", map[string]interface{}{
			"s3Configuration": map[string]interface{}{
				"destinationIdentifier": "s3://results-bucket",
				"roleArn":               "arn:aws:iam::000000000000:role/deliver",
				"ownerAccountId":        "123",
			},
		}, false},
	}
	for _, tc := range cases {
		err := validateDestinationConfiguration(tc.dc)
		if tc.ok && err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("%s: expected error, got none", tc.name)
		}
	}
}

func TestValidateLookupTableKmsKey(t *testing.T) {
	svc := newDeliveryTestService("alias/valid", strings.Repeat("k", 256))
	if err := svc.validateLookupTableKmsKey("alias/valid"); err != nil {
		t.Fatalf("valid key rejected: %v", err)
	}
	if err := svc.validateLookupTableKmsKey(""); err != nil {
		t.Fatalf("empty key rejected: %v", err)
	}
	if err := svc.validateLookupTableKmsKey("alias/unknown"); err == nil {
		t.Fatal("unknown key accepted")
	}
	if err := svc.validateLookupTableKmsKey(strings.Repeat("k", 257)); err == nil {
		t.Fatal("over-length key accepted")
	}
}

func TestLookupTableKMSEncryptionRoundTrip(t *testing.T) {
	svc := newDeliveryTestService("key-1")
	body := "id,name\n1,Alice\n2,Bob"
	lt := &logsstore.LookupTable{Name: "users_enc", KmsKeyId: "key-1"}
	if err := svc.applyLookupTableBody(lt, body, "us-east-1"); err != nil {
		t.Fatal(err)
	}
	if lt.TableBody != "" {
		t.Fatal("encrypted table still stores plaintext body")
	}
	if len(lt.EncryptedBody) == 0 || len(lt.EncryptedDataKey) == 0 || len(lt.ContentNonce) == 0 {
		t.Fatal("encrypted table is missing envelope fields")
	}
	if strings.Contains(string(lt.EncryptedBody), "Alice") {
		t.Fatal("ciphertext contains plaintext")
	}
	if lt.RecordsCount != 2 || len(lt.TableFields) != 2 || lt.SizeBytes != int64(len(body)) {
		t.Fatalf("metadata not derived: %+v", lt)
	}
	plain, err := svc.lookupTablePlainBody(lt, "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if plain != body {
		t.Fatalf("round trip mismatch: %q", plain)
	}

	plainTable := &logsstore.LookupTable{Name: "users"}
	if err := svc.applyLookupTableBody(plainTable, body, "us-east-1"); err != nil {
		t.Fatal(err)
	}
	if plainTable.TableBody != body || len(plainTable.EncryptedBody) != 0 {
		t.Fatal("table without KMS key must store plaintext")
	}
}

func TestDeliverToLookupTableCreatesAndRefreshes(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)
	cfg := map[string]interface{}{
		"tableName":   "refreshed",
		"roleArn":     "arn:aws:iam::000000000000:role/deliver",
		"description": "auto-refreshed",
		"tags":        map[string]interface{}{"team": "ops"},
	}
	rows := []queryResultRow{
		{columns: []string{"id", "name"}, fields: map[string]string{"id": "1", "name": "Alice"}},
		{columns: []string{"id", "name"}, fields: map[string]string{"id": "2", "name": "Bob"}},
	}

	dest := svc.deliverToLookupTable("us-east-1", store, cfg, rows)
	if dest.Status != destinationStatusComplete {
		t.Fatalf("create delivery failed: %s", dest.ErrorMessage)
	}
	lt, err := store.GetLookupTable("refreshed")
	if err != nil {
		t.Fatal(err)
	}
	if lt.Description != "auto-refreshed" || lt.Tags["team"] != "ops" {
		t.Fatalf("creation metadata missing: %+v", lt)
	}
	if lt.RecordsCount != 2 || lt.TableBody == "" {
		t.Fatalf("unexpected body: %+v", lt)
	}

	updated := []queryResultRow{
		{columns: []string{"id", "name"}, fields: map[string]string{"id": "3", "name": "Carol"}},
	}
	dest = svc.deliverToLookupTable("us-east-1", store, cfg, updated)
	if dest.Status != destinationStatusComplete {
		t.Fatalf("refresh delivery failed: %s", dest.ErrorMessage)
	}
	lt, err = store.GetLookupTable("refreshed")
	if err != nil {
		t.Fatal(err)
	}
	if lt.RecordsCount != 1 || !strings.Contains(lt.TableBody, "Carol") || strings.Contains(lt.TableBody, "Alice") {
		t.Fatalf("refresh did not replace content: %q", lt.TableBody)
	}
	if lt.Description != "auto-refreshed" || lt.Tags["team"] != "ops" {
		t.Fatal("refresh overwrote creation-only metadata")
	}

	empty := svc.deliverToLookupTable("us-east-1", store, cfg, nil)
	if empty.Status != destinationStatusClientError {
		t.Fatalf("empty results must be a client error, got %s", empty.Status)
	}
}

func TestDeliverToLookupTableEncrypted(t *testing.T) {
	svc := newDeliveryTestService("key-1")
	store := newDeliveryTestStore(t)
	cfg := map[string]interface{}{
		"tableName": "secure",
		"roleArn":   "arn:aws:iam::000000000000:role/deliver",
		"kmsKeyId":  "key-1",
	}
	rows := []queryResultRow{
		{columns: []string{"id"}, fields: map[string]string{"id": "1"}},
	}
	dest := svc.deliverToLookupTable("us-east-1", store, cfg, rows)
	if dest.Status != destinationStatusComplete {
		t.Fatalf("delivery failed: %s", dest.ErrorMessage)
	}
	lt, err := store.GetLookupTable("secure")
	if err != nil {
		t.Fatal(err)
	}
	if lt.TableBody != "" || len(lt.EncryptedBody) == 0 {
		t.Fatal("destination table was not encrypted")
	}
	plain, err := svc.lookupTablePlainBody(lt, "us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(plain, "id") {
		t.Fatalf("decrypted body lost content: %q", plain)
	}
}

func TestDeliverToS3WithoutBus(t *testing.T) {
	svc := newDeliveryTestService()
	cfg := map[string]interface{}{
		"destinationIdentifier": "s3://results-bucket/prefix",
		"roleArn":               "arn:aws:iam::000000000000:role/deliver",
	}
	dest := svc.deliverToS3("us-east-1", cfg, "sq-1", nil)
	if dest.Status != destinationStatusFailed {
		t.Fatalf("expected FAILED without S3 invoker, got %s", dest.Status)
	}
	if bucket, prefix := splitS3DestinationURI("s3://results-bucket/a/b"); bucket != "results-bucket" || prefix != "a/b" {
		t.Fatalf("split failed: %q %q", bucket, prefix)
	}
	if bucket, prefix := splitS3DestinationURI("s3://results-bucket"); bucket != "results-bucket" || prefix != "" {
		t.Fatalf("split without prefix failed: %q %q", bucket, prefix)
	}
}

// TestScheduledQueryTriggerDelivers exercises the full trigger path: events
// in the log group flow through the query into the configured lookup table
// destination and the execution records the delivery outcome.
func TestScheduledQueryTriggerDelivers(t *testing.T) {
	svc := newDeliveryTestService()
	store := newDeliveryTestStore(t)
	if err := store.CreateLogGroup(&logsstore.LogGroup{Name: "app"}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateLogStream(&logsstore.LogStream{LogGroupName: "app", Name: "s1"}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := store.PutLogEvents("app", "s1", []logsstore.LogEntry{
		{Timestamp: now, Message: `{"id":"u1","name":"Alice"}`, IngestionTime: now},
	}); err != nil {
		t.Fatal(err)
	}

	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-deliver",
		Name:                "deliver",
		QueryString:         "fields id, name",
		LogGroupIdentifiers: []string{"app"},
		State:               "ENABLED",
		DestinationConfiguration: map[string]interface{}{
			"lookupTableConfiguration": map[string]interface{}{
				"tableName": "triggered",
				"roleArn":   "arn:aws:iam::000000000000:role/deliver",
			},
		},
		StartTimeOffset: 60 * 60 * 1000,
		CreationTime:    now,
	}
	svc.triggerScheduledQuery("us-east-1", store, sq)

	execs, err := store.ListScheduledQueryExecutions("sq-deliver", 0, time.Now().Add(time.Hour).UnixMilli())
	if err != nil {
		t.Fatal(err)
	}
	if len(execs) != 1 {
		t.Fatalf("expected one execution record, got %d", len(execs))
	}
	exec := execs[0]
	if exec.Status != "SUCCESS" {
		t.Fatalf("execution status %s: %s", exec.Status, exec.ErrorMessage)
	}
	if len(exec.Destinations) != 1 || exec.Destinations[0].Status != destinationStatusComplete {
		t.Fatalf("unexpected destinations: %+v", exec.Destinations)
	}
	if exec.Destinations[0].DestinationType != "LOOKUP_TABLE" || exec.Destinations[0].ProcessedIdentifier != "triggered" {
		t.Fatalf("unexpected destination record: %+v", exec.Destinations[0])
	}

	lt, err := store.GetLookupTable("triggered")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(lt.TableBody, "Alice") {
		t.Fatalf("delivered body missing query results: %q", lt.TableBody)
	}
}
