// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestListContributorInsightsListsOnlyStateCarryingTables pins the list's
// membership contract: the unfiltered walk keeps the tables whose record an
// UpdateContributorInsights call has stamped (ENABLED and DISABLED alike —
// both are insights state) and omits tables insights were never configured
// on, paging by the last consumed table name when a result page fills before
// the record pages do; the scoped call answers the named table's own
// summary, DISABLED included, the way DescribeContributorInsights does.
func TestListContributorInsightsListsOnlyStateCarryingTables(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	create := func(name string) {
		t.Helper()
		if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":            name,
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		}}); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	// Sorted record order: CiA, CiB, CiC. CiA and CiB carry insights state
	// (enabled and disabled); CiC was never configured.
	create("CiA")
	create("CiB")
	create("CiC")
	update := func(table, action string) {
		t.Helper()
		if _, err := svc.UpdateContributorInsights(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": table, "ContributorInsightsAction": action,
		}}); err != nil {
			t.Fatalf("update insights on %s: %v", table, err)
		}
	}
	update("CiA", "ENABLE")
	update("CiB", "DISABLE")

	list := func(params map[string]interface{}) (interface{}, string) {
		t.Helper()
		resp, err := svc.ListContributorInsights(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		if err != nil {
			t.Fatalf("list contributor insights: %v", err)
		}
		m := resp.(map[string]interface{})
		token, _ := m["NextToken"].(string)
		return m["ContributorInsightsSummaries"], token
	}
	names := func(raw interface{}) []string {
		t.Helper()
		entries := raw.([]map[string]interface{})
		got := make([]string, 0, len(entries))
		for _, e := range entries {
			got = append(got, e["TableName"].(string))
		}
		return got
	}

	// The unfiltered list answers exactly the state-carrying tables.
	summaries, token := list(nil)
	if got := names(summaries); len(got) != 2 || got[0] != "CiA" || got[1] != "CiB" {
		t.Fatalf("unfiltered summaries = %v, want [CiA CiB] (CiC carries no state)", got)
	}
	if token != "" {
		t.Fatalf("unfiltered list carried a token: %q", token)
	}
	for _, e := range summaries.([]map[string]interface{}) {
		want := "DISABLED"
		if e["TableName"] == "CiA" {
			want = "ENABLED"
		}
		if e["ContributorInsightsStatus"] != want {
			t.Fatalf("summary %v: status = %v, want %s", e["TableName"], e["ContributorInsightsStatus"], want)
		}
	}

	// A page size below the state-carrying count pages by the last consumed
	// record: the first page holds CiA alone, the marker resumes after it
	// (CiC sits between the state-carrying records of successive pages
	// without ever appearing), and the final page ends the traversal.
	summaries, token = list(map[string]interface{}{"MaxResults": 1.0})
	if got := names(summaries); len(got) != 1 || got[0] != "CiA" {
		t.Fatalf("first page = %v, want [CiA]", got)
	}
	if token != "CiA" {
		t.Fatalf("first page token = %q, want CiA", token)
	}
	summaries, token = list(map[string]interface{}{"MaxResults": 1.0, "NextToken": token})
	if got := names(summaries); len(got) != 1 || got[0] != "CiB" {
		t.Fatalf("second page = %v, want [CiB]", got)
	}
	if token == "" {
		t.Fatal("second page ended the traversal early: CiC was not yet consumed")
	}
	summaries, token = list(map[string]interface{}{"MaxResults": 1.0, "NextToken": token})
	if got := names(summaries); len(got) != 0 {
		t.Fatalf("final page = %v, want empty", got)
	}
	if token != "" {
		t.Fatalf("final page token = %q, want none", token)
	}

	// The scoped call answers the named table's own summary — DISABLED
	// included for a table insights were never configured on, matching
	// DescribeContributorInsights.
	summaries, _ = list(map[string]interface{}{"TableName": "CiC"})
	if got := names(summaries); len(got) != 1 || got[0] != "CiC" {
		t.Fatalf("scoped never-configured summary = %v, want [CiC]", got)
	}
	if entries := summaries.([]map[string]interface{}); entries[0]["ContributorInsightsStatus"] != "DISABLED" {
		t.Fatalf("scoped never-configured status = %v, want DISABLED", entries[0]["ContributorInsightsStatus"])
	}
}
