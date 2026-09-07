package dynamodb

import (
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func guardTestTable() *dbstore.Table {
	return &dbstore.Table{
		Name: "t",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
	}
}

func TestUpdateClauseTargetNames(t *testing.T) {
	_, clauses, _ := parseUpdateStatement(`UPDATE t SET a = 1, b = 'x' REMOVE extra ADD n 5 DELETE tags 'v' WHERE pk = 'k1'`)
	got := updateClauseTargetNames(clauses)
	want := []string{"a", "b", "extra", "n", "tags"}
	if len(got) != len(want) {
		t.Fatalf("collected targets %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("target %d = %q, want %q (all: %v)", i, got[i], want[i], got)
		}
	}
}

func TestPartiQLUpdateRejectsKeyAttributeTargets(t *testing.T) {
	table := guardTestTable()

	cases := []struct {
		name      string
		statement string
	}{
		{"SET partition key", `UPDATE t SET pk = 'other' WHERE pk = 'k1' AND sk = 's1'`},
		{"REMOVE sort key", `UPDATE t SET a = 1 REMOVE sk WHERE pk = 'k1' AND sk = 's1'`},
		{"ADD partition key", `UPDATE t SET a = 1 ADD pk 5 WHERE pk = 'k1' AND sk = 's1'`},
	}
	for _, tc := range cases {
		_, clauses, _ := parseUpdateStatement(tc.statement)
		err := validateNotKeyAttributes(table, updateClauseTargetNames(clauses))
		if err == nil {
			t.Fatalf("%s: expected rejection", tc.name)
		}
		apiErr, ok := err.(*APIError)
		if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
			t.Fatalf("%s: expected ValidationException, got %v", tc.name, err)
		}
	}

	// A statement whose clauses touch only non-key attributes passes the guard.
	_, clauses, _ := parseUpdateStatement(`UPDATE t SET a = 1 REMOVE extra WHERE pk = 'k1' AND sk = 's1'`)
	if err := validateNotKeyAttributes(table, updateClauseTargetNames(clauses)); err != nil {
		t.Fatalf("non-key statement rejected: %v", err)
	}
}

func matchSAttr(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{S: &v} }

func TestMatchSingleTargetItemEnforcesSingleItemContract(t *testing.T) {
	table := guardTestTable()
	mkItem := func(sk string) *dbstore.Item {
		return &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{
			"pk": matchSAttr("a"),
			"sk": matchSAttr(sk),
		}}
	}
	items := []*dbstore.Item{mkItem("s1"), mkItem("s2"), mkItem("s3")}

	// The scan adapter emulates a partition scan: it feeds every item and
	// stops as soon as the matcher's callback signals it has enough.
	run := func(where string, wantOldAll bool) (item *dbstore.Item, preFilter []*dbstore.Item, err error, invoked int) {
		_, whereExpr := parseSelectStatement("SELECT * FROM t WHERE " + where)
		scan := func(pkValue string, cb func(*dbstore.Item) error) error {
			for _, it := range items {
				invoked++
				if cbErr := cb(it); cbErr != nil {
					return cbErr
				}
			}
			return nil
		}
		item, preFilter, _, err = matchSingleTargetItem(scan, table, whereExpr, nil, "UPDATE", wantOldAll)
		return item, preFilter, err, invoked
	}

	// A nil WHERE clause is rejected.
	if _, _, err, _ := run("", false); err == nil {
		t.Fatal("empty WHERE text: expected rejection")
	}

	// A WHERE without a partition-key equality is rejected before any scan.
	if _, _, err, invoked := run("v = 'x'", false); err == nil || invoked != 0 {
		t.Fatalf("non-key WHERE: expected rejection without scanning, got err=%v invoked=%d", err, invoked)
	}

	// Exactly one match: the item is returned and every item was scanned.
	item, _, err, invoked := run("pk = 'a' AND sk = 's2'", false)
	if err != nil || item == nil {
		t.Fatalf("single match: expected item, got err=%v item=%v", err, item)
	}
	if got := *item.Attributes["sk"].S; got != "s2" {
		t.Fatalf("single match: expected sk=s2, got %q", got)
	}
	if invoked != 3 {
		t.Fatalf("single match: expected 3 scanned items, got %d", invoked)
	}

	// More than one match: ValidationException, and the scan stopped after
	// the second match.
	_, _, err, invoked = run("pk = 'a'", false)
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazonaws.dynamodb.v20120810#ValidationException" {
		t.Fatalf("multi match: expected ValidationException, got %v", err)
	}
	if invoked != 2 {
		t.Fatalf("multi match: expected the scan to stop after 2 items, got %d", invoked)
	}

	// Zero matches with ALL_OLD reporting: no item, every scanned item kept.
	item, preFilter, err, _ := run("pk = 'a' AND sk = 'zz'", true)
	if err != nil || item != nil {
		t.Fatalf("zero match: expected no error and no item, got err=%v item=%v", err, item)
	}
	if len(preFilter) != 3 {
		t.Fatalf("zero match: expected 3 pre-filter items for ALL_OLD, got %d", len(preFilter))
	}
}
