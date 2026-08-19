package apps

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// stubDynamoDBStoreProvider satisfies dynamoDBStoreProvider without wiring
// the whole DynamoDBService.
type stubDynamoDBStoreProvider struct {
	store dynamodbstore.DynamoDBStoreInterface
}

func (p stubDynamoDBStoreProvider) GetStoreForRegion(string) (dynamodbstore.DynamoDBStoreInterface, error) {
	return p.store, nil
}

func strAttr(v string) *dynamodbstore.AttributeValue {
	return &dynamodbstore.AttributeValue{S: &v}
}

// newInvokerTestStore opens a store over a temporary directory with one
// pk+sk table. Insights are enabled only when the caller asks for them.
func newInvokerTestStore(t *testing.T, enableInsights bool) dynamodbstore.DynamoDBStoreInterface {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	store := dynamodbstore.NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(
		"Tbl",
		[]*dynamodbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dynamodbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dynamodbstore.KeyTypeRange},
		},
		[]*dynamodbstore.AttributeDefinition{
			{AttributeName: "pk", AttributeType: dynamodbstore.ScalarAttributeTypeS},
			{AttributeName: "sk", AttributeType: dynamodbstore.ScalarAttributeTypeS},
		},
		dynamodbstore.BillingModePayPerRequest, nil, nil, nil, nil, nil, false,
	); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if enableInsights {
		tbl, err := store.Tables().Get("Tbl")
		if err != nil {
			t.Fatalf("get table: %v", err)
		}
		tbl.ContributorInsightsEnabled = true
		if err := store.Tables().Put(tbl); err != nil {
			t.Fatalf("enable insights: %v", err)
		}
	}
	return store
}

// TestDynamoDBInvokerWritesFeedContributorInsights pins that item writes
// issued through the cross-service invoker are counted by contributor
// insights exactly like direct data-plane writes: every PutItem, UpdateItem
// and DeleteItem is one write event on the touched key.
func TestDynamoDBInvokerWritesFeedContributorInsights(t *testing.T) {
	store := newInvokerTestStore(t, true)
	adapter := &dynamoDBInvokerAdapter{provider: stubDynamoDBStoreProvider{store: store}}
	ctx := t.Context()

	if _, err := adapter.PutItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "A", "sk": "1"},
		map[string]interface{}{"val": "x"}); err != nil {
		t.Fatalf("put A: %v", err)
	}
	if _, err := adapter.PutItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "B", "sk": "1"},
		map[string]interface{}{"val": "y"}); err != nil {
		t.Fatalf("put B: %v", err)
	}
	if err := adapter.UpdateItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "B", "sk": "1"},
		map[string]interface{}{"val": "z"}); err != nil {
		t.Fatalf("update B: %v", err)
	}
	if err := adapter.DeleteItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "A", "sk": "1"}); err != nil {
		t.Fatalf("delete A: %v", err)
	}

	tbl, err := store.Tables().Get("Tbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	keyA := map[string]*dynamodbstore.AttributeValue{"pk": strAttr("A"), "sk": strAttr("1")}
	keyB := map[string]*dynamodbstore.AttributeValue{"pk": strAttr("B"), "sk": strAttr("1")}
	pkOf := func(key map[string]*dynamodbstore.AttributeValue) string {
		return dynamodbstore.ContributorKeyString(tbl, key, dynamodbstore.ContributorLayoutPartitionKey)
	}
	fullOf := func(key map[string]*dynamodbstore.AttributeValue) string {
		return dynamodbstore.ContributorKeyString(tbl, key, dynamodbstore.ContributorLayoutFullKey)
	}

	now := time.Now()
	window := func(layout string) map[string]dynamodbstore.ContributorKeyStat {
		stats, err := store.Contributors().TopKeys("Tbl", layout, now.Add(-time.Minute), now.Add(time.Minute), 10)
		if err != nil {
			t.Fatalf("top keys %s: %v", layout, err)
		}
		byKey := make(map[string]dynamodbstore.ContributorKeyStat, len(stats))
		for _, s := range stats {
			byKey[s.Key] = s
		}
		return byKey
	}

	pkc := window(dynamodbstore.ContributorLayoutPartitionKey)
	if len(pkc) != 2 {
		t.Fatalf("expected both partitions counted, got %v", pkc)
	}
	for key, want := range map[string]int64{pkOf(keyA): 2, pkOf(keyB): 2} {
		got := pkc[key]
		if got.Count != want {
			t.Fatalf("partition %s: expected %d events, got %d", key, want, got.Count)
		}
		if got.Units != float64(want)*dynamodbstore.ContributorWriteUnits {
			t.Fatalf("partition %s: expected %.0f units, got %f", key, float64(want)*dynamodbstore.ContributorWriteUnits, got.Units)
		}
	}

	full := window(dynamodbstore.ContributorLayoutFullKey)
	if len(full) != 2 {
		t.Fatalf("expected both full keys counted, got %v", full)
	}
	for key, want := range map[string]int64{fullOf(keyA): 2, fullOf(keyB): 2} {
		if got := full[key]; got.Count != want {
			t.Fatalf("full key %s: expected %d events, got %d", key, want, got.Count)
		}
	}
}

// TestDynamoDBInvokerReadsFeedContributorInsights pins that reads issued
// through the cross-service invoker are counted by contributor insights
// exactly like direct data-plane reads: a found GetItem is one read event,
// a Scan counts every returned item, and a Query (paginated or not) is a
// single event on the partition-key series. A GetItem for a missing key is
// not counted, mirroring the data-plane GetItem behaviour.
func TestDynamoDBInvokerReadsFeedContributorInsights(t *testing.T) {
	store := newInvokerTestStore(t, true)
	adapter := &dynamoDBInvokerAdapter{provider: stubDynamoDBStoreProvider{store: store}}
	ctx := t.Context()

	put := func(pk, sk string) {
		t.Helper()
		if _, err := adapter.PutItem(ctx, "us-east-1", "Tbl",
			map[string]interface{}{"pk": pk, "sk": sk},
			map[string]interface{}{"val": "x"}); err != nil {
			t.Fatalf("put %s/%s: %v", pk, sk, err)
		}
	}
	put("A", "1")
	put("A", "2")
	put("B", "1")

	if item, err := adapter.GetItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "A", "sk": "1"}); err != nil || item == nil {
		t.Fatalf("get A/1: %v %v", item, err)
	}
	if _, err := adapter.GetItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "C", "sk": "9"}); err == nil {
		t.Fatalf("get of a missing key must fail through the invoker")
	}
	if items, err := adapter.Scan(ctx, "us-east-1", "Tbl", 0); err != nil || len(items) != 3 {
		t.Fatalf("scan: %d items, err %v", len(items), err)
	}
	if items, err := adapter.Query(ctx, "us-east-1", "Tbl", "A", 0); err != nil || len(items) != 2 {
		t.Fatalf("query A: %d items, err %v", len(items), err)
	}
	pageItems, marker, err := adapter.ScanWithPagination(ctx, "us-east-1", "Tbl", 2, "")
	if err != nil || len(pageItems) != 2 || marker == "" {
		t.Fatalf("scan page: %d items, marker %q, err %v", len(pageItems), marker, err)
	}
	queryItems, queryMarker, err := adapter.QueryWithPagination(ctx, "us-east-1", "Tbl", "A", 1, "")
	if err != nil || len(queryItems) != 1 || queryMarker == "" {
		t.Fatalf("query page: %d items, marker %q, err %v", len(queryItems), queryMarker, err)
	}

	tbl, err := store.Tables().Get("Tbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	stat := func(layout, key string) dynamodbstore.ContributorKeyStat {
		now := time.Now()
		stats, err := store.Contributors().TopKeys("Tbl", layout, now.Add(-time.Minute), now.Add(time.Minute), 10)
		if err != nil {
			t.Fatalf("top keys %s: %v", layout, err)
		}
		for _, s := range stats {
			if s.Key == key {
				return s
			}
		}
		return dynamodbstore.ContributorKeyStat{}
	}
	keyOf := func(pk, sk string) string {
		return dynamodbstore.ContributorKeyString(tbl,
			map[string]*dynamodbstore.AttributeValue{"pk": strAttr(pk), "sk": strAttr(sk)},
			dynamodbstore.ContributorLayoutFullKey)
	}
	pkOf := func(pk string) string {
		return dynamodbstore.ContributorKeyString(tbl,
			map[string]*dynamodbstore.AttributeValue{"pk": strAttr(pk)},
			dynamodbstore.ContributorLayoutPartitionKey)
	}

	// writes: A x2, B x1; GetItem A/1; Scan x3; Query A; query page A;
	// scan page A/1+A/2.
	wantPKC := map[string]int64{pkOf("A"): 9, pkOf("B"): 2}
	for key, want := range wantPKC {
		if got := stat(dynamodbstore.ContributorLayoutPartitionKey, key); got.Count != want {
			t.Fatalf("partition %s: expected %d events, got %d", key, want, got.Count)
		}
	}
	wantFull := map[string]int64{keyOf("A", "1"): 4, keyOf("A", "2"): 3, keyOf("B", "1"): 2}
	for key, want := range wantFull {
		if got := stat(dynamodbstore.ContributorLayoutFullKey, key); got.Count != want {
			t.Fatalf("full key %s: expected %d events, got %d", key, want, got.Count)
		}
	}
	if got := stat(dynamodbstore.ContributorLayoutFullKey, keyOf("C", "9")); got.Count != 0 {
		t.Fatalf("a GetItem for a missing key must not be counted, got %d", got.Count)
	}
}

// TestDynamoDBInvokerPutItemMergesKeyAndAttributes pins the response of the
// invoker PutItem: the stored item is the key attributes merged with the
// payload attributes.
func TestDynamoDBInvokerPutItemMergesKeyAndAttributes(t *testing.T) {
	store := newInvokerTestStore(t, false)
	adapter := &dynamoDBInvokerAdapter{provider: stubDynamoDBStoreProvider{store: store}}

	item, err := adapter.PutItem(t.Context(), "us-east-1", "Tbl",
		map[string]interface{}{"pk": "A", "sk": "1"},
		map[string]interface{}{"val": "x"})
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	for k, want := range map[string]interface{}{"pk": "A", "sk": "1", "val": "x"} {
		if item[k] != want {
			t.Fatalf("result[%s]: expected %v, got %v", k, want, item[k])
		}
	}

	stored, err := store.Items().Get("Tbl", map[string]*dynamodbstore.AttributeValue{"pk": strAttr("A"), "sk": strAttr("1")})
	if err != nil {
		t.Fatalf("get stored item: %v", err)
	}
	if stored == nil || stored.Attributes["val"] == nil || *stored.Attributes["val"].S != "x" {
		t.Fatalf("stored item missing merged attributes: %+v", stored)
	}
}

// TestDynamoDBInvokerPutItemJournalsUnderPITR pins that invoker writes to a
// point-in-time-recovery table append journal records, so restore keeps
// working for items written through cross-service integrations.
func TestDynamoDBInvokerPutItemJournalsUnderPITR(t *testing.T) {
	store := newInvokerTestStore(t, false)
	tbl, err := store.Tables().Get("Tbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	tbl.PointInTimeRecovery = &dynamodbstore.PointInTimeRecoveryDescription{Status: dynamodbstore.PITRStatusEnabled}
	if err := store.Tables().Put(tbl); err != nil {
		t.Fatalf("enable pitr: %v", err)
	}

	adapter := &dynamoDBInvokerAdapter{provider: stubDynamoDBStoreProvider{store: store}}
	ctx := t.Context()
	if _, err := adapter.PutItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "A", "sk": "1"},
		map[string]interface{}{"val": "x"}); err != nil {
		t.Fatalf("put: %v", err)
	}
	if err := adapter.DeleteItem(ctx, "us-east-1", "Tbl",
		map[string]interface{}{"pk": "A", "sk": "1"}); err != nil {
		t.Fatalf("delete: %v", err)
	}

	records := 0
	err = store.Journal().ReverseReplay("Tbl", time.Now().Add(-time.Minute), func(*dynamodbstore.JournalChange) error {
		records++
		return nil
	})
	if err != nil {
		t.Fatalf("replay journal: %v", err)
	}
	if records != 2 {
		t.Fatalf("expected the put and the delete to be journaled, got %d records", records)
	}
}
