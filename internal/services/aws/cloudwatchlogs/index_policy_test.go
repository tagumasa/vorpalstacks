package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// newIndexPolicyTestEnv builds a service with one Standard-class group
// and one stream.
func newIndexPolicyTestEnv(t *testing.T, group string) (*LogsService, *logsstore.Store) {
	t.Helper()
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", group)); err != nil {
		t.Fatalf("create log stream: %v", err)
	}
	return svc, store
}

// putIndexEvents writes one batch of JSON-or-plain events at the given
// millisecond timestamps.
func putIndexEvents(t *testing.T, svc *LogsService, group string, at int64, messages ...string) {
	t.Helper()
	events := make([]PutLogEvent, len(messages))
	for i, message := range messages {
		events[i] = PutLogEvent{
			LogEntry:     logsstore.LogEntry{Timestamp: at + int64(i), Message: message},
			TimestampSet: true,
		}
	}
	if _, err := svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: group, LogStreamName: "s1", Events: events, Region: "us-east-1",
	}); err != nil {
		t.Fatalf("put log events: %v", err)
	}
}

// fieldIndexRows runs DescribeFieldIndexes over one group and indexes
// the rows by field name and category (a policy may index a field the
// DEFAULT category already lists — @logStream is a documented example
// of both).
func fieldIndexRows(t *testing.T, svc *LogsService, group string, categories ...string) map[string]map[string]interface{} {
	t.Helper()
	entries, next, err := svc.describeFieldIndexesCore([]string{group}, categories, "us-east-1", "")
	if err != nil {
		t.Fatalf("describe field indexes: %v", err)
	}
	if next != "" {
		t.Fatalf("unexpected pagination: %v", next)
	}
	rows := make(map[string]map[string]interface{}, len(entries))
	for _, entry := range entries {
		rows[entry.fieldIndexName+"\x00"+entry.indexCategory] = formatFieldIndex(entry)
	}
	return rows
}

// quoteJoin renders a JSON string array body (["a","b"]).
func quoteJoin(items []string) string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = `"` + item + `"`
	}
	return strings.Join(quoted, ",")
}

// TestIndexPolicyDocumentRows pins the policy document's validation
// matrix: the documented quotas (at least one field, twenty fields,
// hundred-character names), the FieldsV2 type enum, the mutual
// exclusivity rule and the @-doubling rule for custom fields.
func TestIndexPolicyDocumentRows(t *testing.T) {
	over20 := make([]string, 21)
	for i := range over20 {
		over20[i] = fmt.Sprintf("F%02d", i)
	}
	longName := strings.Repeat("F", logsstore.IndexPolicyFieldNameMax+1)
	rows := []struct {
		name string
		doc  string
		err  string // "" accepts
	}{
		{"fields array accepts", `{"Fields":["RequestId","TransactionId"]}`, ""},
		{"fieldsV2 both types accept", `{"FieldsV2":{"RequestId":{"type":"FIELD_INDEX"},"APIName":{"type":"FACET"}}}`, ""},
		{"both members accept", `{"Fields":["TransactionId"],"FieldsV2":{"RequestId":{"type":"FIELD_INDEX"}}}`, ""},
		{"empty document rejects", `{}`, "at least one field index"},
		{"twenty-first field rejects", `{"Fields": [` + quoteJoin(over20) + `]}`, "20"},
		{"hundred-first character rejects", `{"Fields":["` + longName + `"]}`, "100"},
		{"overlapping member rejects", `{"Fields":["RequestId"],"FieldsV2":{"RequestId":{"type":"FACET"}}}`, "mutually exclusive"},
		{"unsupported type rejects", `{"FieldsV2":{"A":{"type":"PARTIAL"}}}`, "FIELD_INDEX and FACET"},
		{"single-at custom field rejects", `{"Fields":["@userId"]}`, "extra @"},
		{"generated field accepts", `{"Fields":["@logStream"]}`, ""},
		{"doubled-at custom field accepts", `{"Fields":["@@userId"]}`, ""},
		{"malformed document rejects", `{Fields`, "JSON"},
	}
	for _, row := range rows {
		_, err := parseIndexPolicyDocument(row.doc)
		if row.err == "" {
			if err != nil {
				t.Errorf("%s: %v", row.name, err)
			}
			continue
		}
		if err == nil || !strings.Contains(err.Error(), row.err) {
			t.Errorf("%s: got %v, want %q", row.name, err, row.err)
		}
	}
}

// TestIndexPolicyCRUDAndPrecedence pins the policy surface's lifecycle
// and the documented precedence: the group-level policy is returned
// while it exists, the applying account-level policy after it is
// deleted, and nothing when neither exists.
func TestIndexPolicyCRUDAndPrecedence(t *testing.T) {
	const group = "idx-group"
	svc, store := newIndexPolicyTestEnv(t, group)

	policy, err := svc.putIndexPolicyCore(group, `{"Fields":["RequestId"],"FieldsV2":{"APIName":{"type":"FACET"}}}`, "us-east-1")
	if err != nil {
		t.Fatalf("put index policy: %v", err)
	}
	if policy.LogGroupIdentifier == "" || policy.LogGroupIdentifier == group {
		t.Fatalf("policy ARN echo: %q", policy.LogGroupIdentifier)
	}

	views, err := svc.describeIndexPoliciesCore([]string{group}, "us-east-1")
	if err != nil || len(views) != 1 {
		t.Fatalf("describe: %v %d", err, len(views))
	}
	formatted := formatIndexPolicy(views[0].groupARN, views[0].eff.policyDocument, views[0].eff.lastUpdateTime, views[0].eff.source, views[0].eff.policyName)
	if formatted["source"] != "LOG_GROUP" {
		t.Fatalf("group policy source: %v", formatted["source"])
	}
	if _, has := formatted["policyName"]; has {
		t.Fatal("group-level policy carries a policyName; log group-level policies don't have names")
	}
	if formatted["logGroupIdentifier"] != policy.LogGroupIdentifier {
		t.Fatalf("describe ARN: %v", formatted["logGroupIdentifier"])
	}

	// An account-level policy applies to the group but never overrides
	// the group-level one.
	if _, err := svc.putAccountPolicyCore("idx-account", `{"Fields":["SessionId"]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = idx-", "us-east-1"); err != nil {
		t.Fatalf("put account policy: %v", err)
	}
	views, err = svc.describeIndexPoliciesCore([]string{group}, "us-east-1")
	if err != nil || len(views) != 1 || views[0].eff.source != "LOG_GROUP" {
		t.Fatalf("group policy does not override the account policy: %v", views)
	}

	if err := svc.deleteIndexPolicyCore(group, "us-east-1"); err != nil {
		t.Fatalf("delete index policy: %v", err)
	}
	views, err = svc.describeIndexPoliciesCore([]string{group}, "us-east-1")
	if err != nil || len(views) != 1 || views[0].eff.source != "ACCOUNT" {
		t.Fatalf("account policy does not apply after deletion: %v", views)
	}
	formatted = formatIndexPolicy(views[0].groupARN, views[0].eff.policyDocument, views[0].eff.lastUpdateTime, views[0].eff.source, views[0].eff.policyName)
	if formatted["policyName"] != "idx-account" {
		t.Fatalf("account policy name: %v", formatted["policyName"])
	}

	if err := svc.deleteIndexPolicyCore(group, "us-east-1"); err == nil {
		t.Fatal("second delete succeeded; the policy is gone")
	}
	if _, err := svc.describeIndexPoliciesCore([]string{group, "idx-group"}, "us-east-1"); err == nil {
		t.Fatal("two identifiers accepted; the member's length trait is 1..1")
	}

	// A group outside the prefix carries no policy at all.
	if err := store.CreateLogGroup(logsstore.NewLogGroup("other-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create other group: %v", err)
	}
	views, err = svc.describeIndexPoliciesCore([]string{"other-group"}, "us-east-1")
	if err != nil || len(views) != 0 {
		t.Fatalf("unindexed group describes policies: %v", views)
	}

	// The group teardown removes the policy record and the trail: a
	// same-named recreation starts unindexed.
	putIndexEvents(t, svc, group, time.Now().UnixMilli(), `{"RequestId":"r1"}`)
	if _, err := svc.putIndexPolicyCore(group, `{"Fields":["RequestId"]}`, "us-east-1"); err != nil {
		t.Fatalf("re-put index policy: %v", err)
	}
	if err := store.DeleteLogGroup(group); err != nil {
		t.Fatalf("delete log group: %v", err)
	}
	if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("recreate log group: %v", err)
	}
	if _, err := store.GetIndexPolicy(group); err == nil {
		t.Fatal("index policy survived the group deletion")
	}
	if trail := store.GetIndexInactiveTrail(group); len(trail.Fields) != 0 {
		t.Fatalf("inactive trail survived the group deletion: %v", trail.Fields)
	}
}

// TestDescribeFieldIndexesDerivation pins the listing's derivation: the
// fixed DEFAULT category, the effective policy's CUSTOM fields with the
// bounds the ingestion scans accumulated (events ingested before the
// policy existed never carry — firstEventTime is "the earliest log event
// that matches this field index, after the index policy that contains it
// was created"), the dotted and array-segment JSON paths, the identity
// fields' whole-history DEFAULT bounds (stream-derived) and the INACTIVE
// trail's frozen bounds.
func TestDescribeFieldIndexesDerivation(t *testing.T) {
	const group = "idx-scan-group"
	svc, _ := newIndexPolicyTestEnv(t, group)
	base := time.Now().UnixMilli() - 60000

	putIndexEvents(t, svc, group, base, `{"TransactionId":"t-old","level":"INFO"}`)

	policyStamp := time.Now().UnixMilli()
	time.Sleep(5 * time.Millisecond)
	if _, err := svc.putIndexPolicyCore(group, `{"Fields":["TransactionId","userIdentity.accessKeyId","items.0.id"],"FieldsV2":{"@logStream":{"type":"FIELD_INDEX"}}}`, "us-east-1"); err != nil {
		t.Fatalf("put index policy: %v", err)
	}
	putIndexEvents(t, svc, group, policyStamp+1000,
		`{"TransactionId":"t-new","userIdentity":{"accessKeyId":"AKIA"},"items":[{"id":"i1"}]}`,
		`{"level":"WARN"}`,
		"plain text line",
	)

	rows := fieldIndexRows(t, svc, group)
	defaultSeen := 0
	for _, name := range defaultFieldIndexes {
		row, ok := rows[name+"\x00DEFAULT"]
		if !ok {
			t.Fatalf("DEFAULT field %s missing from the default listing", name)
		}
		if row["indexCategory"] != "DEFAULT" || row["type"] != "FIELD_INDEX" {
			t.Fatalf("DEFAULT row %s: %v", name, row)
		}
		if _, has := row["logGroupIdentifier"]; has {
			t.Fatalf("DEFAULT row %s carries a logGroupIdentifier; default fields appear in no single-group policy", name)
		}
		defaultSeen++
	}
	if defaultSeen != len(defaultFieldIndexes) {
		t.Fatalf("default listing: %d rows", defaultSeen)
	}

	custom := rows["TransactionId\x00CUSTOM"]
	if custom == nil || custom["indexCategory"] != "CUSTOM" {
		t.Fatalf("custom TransactionId row: %v", custom)
	}
	if first := custom["firstEventTime"].(int64); first < policyStamp {
		t.Fatalf("firstEventTime %d precedes the policy stamp %d; events before the policy was created do not count", first, policyStamp)
	}
	if custom["logGroupIdentifier"] == nil {
		t.Fatal("single-group CUSTOM row omits logGroupIdentifier")
	}
	if rows["userIdentity.accessKeyId\x00CUSTOM"] == nil || rows["items.0.id\x00CUSTOM"] == nil {
		t.Fatal("dotted and array-segment JSON paths missing from the CUSTOM listing")
	}
	// The @logStream policy field is CUSTOM on top of its DEFAULT entry.
	if rows["@logStream\x00CUSTOM"]["indexCategory"] != "CUSTOM" {
		t.Fatalf("@logStream policy row category: %v", rows["@logStream\x00CUSTOM"])
	}
	// The two @logStream rows answer different questions: the DEFAULT
	// row derives from the stream records and spans the whole event
	// history (the identity field matches every event by construction),
	// while the CUSTOM row counts only the scans the policy's creation
	// gated — the pre-policy event at base carries no CUSTOM bound.
	if row := rows["@logStream\x00DEFAULT"]; row["firstEventTime"].(int64) != base || row["lastEventTime"].(int64) != policyStamp+1002 {
		t.Fatalf("@logStream DEFAULT bounds span the history: first=%v last=%v want %d..%d",
			row["firstEventTime"], row["lastEventTime"], base, policyStamp+1002)
	}
	if row := rows["@logStream\x00CUSTOM"]; row["firstEventTime"].(int64) != policyStamp+1000 || row["lastEventTime"].(int64) != policyStamp+1002 {
		t.Fatalf("@logStream CUSTOM bounds are policy-gated: first=%v last=%v want %d..%d",
			row["firstEventTime"], row["lastEventTime"], policyStamp+1000, policyStamp+1002)
	}

	customOnly := fieldIndexRows(t, svc, group, "CUSTOM")
	if _, has := customOnly["@aws.region\x00DEFAULT"]; has {
		t.Fatal("category filter ignored: DEFAULT row present in a CUSTOM-only listing")
	}
	autoOnly := fieldIndexRows(t, svc, group, "AUTO")
	if len(autoOnly) != 0 {
		t.Fatalf("AUTO category carries fields: %v", autoOnly)
	}
	if _, _, err := svc.describeFieldIndexesCore([]string{group}, []string{"NOT_A_CATEGORY"}, "us-east-1", ""); err == nil {
		t.Fatal("invalid index category accepted")
	}
	if _, _, err := svc.describeFieldIndexesCore([]string{group, "missing-group-entirely"}, []string{"AUTO"}, "us-east-1", ""); err == nil {
		t.Fatal("missing group accepted")
	}

	// Dropping a field from the policy moves it to INACTIVE; its bounds
	// freeze at the range the scans accumulated while the policy watched
	// the field ("the earliest log event that matches this field index,
	// after the index policy that contains it was created" — the
	// pre-policy event at base was never scanned for this index).
	if _, err := svc.putIndexPolicyCore(group, `{"Fields":["userIdentity.accessKeyId"]}`, "us-east-1"); err != nil {
		t.Fatalf("replace index policy: %v", err)
	}
	inactive := fieldIndexRows(t, svc, group, "INACTIVE")
	row, ok := inactive["TransactionId\x00INACTIVE"]
	if !ok {
		t.Fatalf("dropped field missing from INACTIVE: %v", inactive)
	}
	if first := row["firstEventTime"].(int64); first != policyStamp+1000 {
		t.Fatalf("INACTIVE firstEventTime %d, want the first post-policy match at %d", first, policyStamp+1000)
	}
	if row["logGroupIdentifier"] != nil {
		t.Fatal("INACTIVE row carries logGroupIdentifier; it appears in no single-group policy")
	}
	// Events ingested while no policy watches the field carry nothing
	// into its frozen range.
	putIndexEvents(t, svc, group, policyStamp+90000, `{"TransactionId":"t-idle"}`)
	// Re-adding the field resumes from the frozen bounds: the index the
	// original policy created still owns the accumulated range.
	if _, err := svc.putIndexPolicyCore(group, `{"Fields":["userIdentity.accessKeyId","TransactionId"]}`, "us-east-1"); err != nil {
		t.Fatalf("re-add index policy: %v", err)
	}
	resumed := fieldIndexRows(t, svc, group, "CUSTOM")
	if row := resumed["TransactionId\x00CUSTOM"]; row["firstEventTime"].(int64) != policyStamp+1000 || row["lastEventTime"].(int64) != policyStamp+1000 {
		t.Fatalf("re-added field bounds: first=%v last=%v want the frozen %d..%d",
			row["firstEventTime"], row["lastEventTime"], policyStamp+1000, policyStamp+1000)
	}
	putIndexEvents(t, svc, group, policyStamp+95000, `{"TransactionId":"t-resume"}`)
	extended := fieldIndexRows(t, svc, group, "CUSTOM")
	if row := extended["TransactionId\x00CUSTOM"]; row["firstEventTime"].(int64) != policyStamp+1000 || row["lastEventTime"].(int64) != policyStamp+95000 {
		t.Fatalf("resumed field bounds: first=%v last=%v want %d..%d",
			row["firstEventTime"], row["lastEventTime"], policyStamp+1000, policyStamp+95000)
	}
	// The account policy taking over a deleted group policy subtracts
	// its own fields from the trail.
	if err := svc.deleteIndexPolicyCore(group, "us-east-1"); err != nil {
		t.Fatalf("delete index policy: %v", err)
	}
	trail := fieldIndexRows(t, svc, group, "INACTIVE")
	if _, has := trail["userIdentity.accessKeyId\x00INACTIVE"]; !has {
		t.Fatalf("deleted policy's field missing from INACTIVE: %v", trail)
	}
}

// TestIngestedFieldIndexScanUsesTransformedForm pins the scan's message
// form: a group with an effective transformer is scanned over the
// transformed copy — a field the transformer exposes gains bounds and a
// field it deletes loses them, because the query plane never serves the
// raw message for a transformed group.
func TestIngestedFieldIndexScanUsesTransformedForm(t *testing.T) {
	const group = "idx-transform-group"
	svc, _ := newIndexPolicyTestEnv(t, group)

	if _, err := svc.putTransformerCore(group, []map[string]interface{}{
		{"parseJSON": map[string]interface{}{}},
		{"addKeys": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"key": "correlationId", "value": "c-1"},
		}}},
		{"deleteKeys": map[string]interface{}{"withKeys": []interface{}{"level"}}},
	}, "us-east-1"); err != nil {
		t.Fatalf("put transformer: %v", err)
	}
	if _, err := svc.putIndexPolicyCore(group, `{"Fields":["correlationId","level"]}`, "us-east-1"); err != nil {
		t.Fatalf("put index policy: %v", err)
	}
	at := time.Now().UnixMilli()
	putIndexEvents(t, svc, group, at, `{"level":"INFO"}`)

	rows := fieldIndexRows(t, svc, group, "CUSTOM")
	if row := rows["correlationId\x00CUSTOM"]; row == nil || row["firstEventTime"].(int64) != at {
		t.Fatalf("transformer-injected field bounds: %v", row)
	}
	if row := rows["level\x00CUSTOM"]; row == nil || row["firstEventTime"] != nil {
		t.Fatalf("transformer-deleted field must carry no bounds: %v", row)
	}
}

// TestFieldIndexAccountPolicySelection pins the account-level
// FIELD_INDEX_POLICY rules at PutAccountPolicy: the shared document
// syntax, the prefix overlap rule, the account-wide/prefix mutual
// exclusion, the data-source criteria's platform rejection and the
// prefix-policy quota.
func TestFieldIndexAccountPolicySelection(t *testing.T) {
	const group = "idx-account-group"
	svc, _ := newIndexPolicyTestEnv(t, group)

	if _, err := svc.putAccountPolicyCore("fip-1", `{"Fields":["A"]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = app-", "us-east-1"); err != nil {
		t.Fatalf("prefix policy: %v", err)
	}
	if _, err := svc.putAccountPolicyCore("fip-2", `{"Fields":["B"]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = app-", "us-east-1"); err == nil {
		t.Fatal("same prefix accepted twice")
	}
	if _, err := svc.putAccountPolicyCore("fip-3", `{"Fields":["B"]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = app-prod-", "us-east-1"); err == nil {
		t.Fatal("contained prefix accepted")
	}
	if _, err := svc.putAccountPolicyCore("fip-wide", `{"Fields":["B"]}`, "FIELD_INDEX_POLICY", "", "", "us-east-1"); err == nil {
		t.Fatal("account-wide policy accepted alongside prefix policies")
	}
	if _, err := svc.putAccountPolicyCore("fip-ds", `{"Fields":["B"]}`, "FIELD_INDEX_POLICY", "", "DataSourceName = amazon_vpc AND DataSourceType = flow", "us-east-1"); err == nil {
		t.Fatal("data-source criteria accepted; the platform carries no vended data sources")
	}
	if _, err := svc.putAccountPolicyCore("fip-bad", `{"Nope":[]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = zzz-", "us-east-1"); err == nil {
		t.Fatal("field-less policy document accepted")
	}
	// Same-name replacement is exempt from its own overlap row.
	if _, err := svc.putAccountPolicyCore("fip-1", `{"Fields":["A2"]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = app-", "us-east-1"); err != nil {
		t.Fatalf("self-replacement: %v", err)
	}

	// The account-wide slot opens once the prefix policies are gone,
	// and the prefix quota rejects the twenty-first prefix policy.
	if err := svc.deleteAccountPolicyCore("fip-1", "FIELD_INDEX_POLICY", "us-east-1"); err != nil {
		t.Fatalf("delete account policy: %v", err)
	}
	if _, err := svc.putAccountPolicyCore("fip-wide", `{"Fields":["B"]}`, "FIELD_INDEX_POLICY", "", "", "us-east-1"); err != nil {
		t.Fatalf("account-wide policy after cleanup: %v", err)
	}
	if err := svc.deleteAccountPolicyCore("fip-wide", "FIELD_INDEX_POLICY", "us-east-1"); err != nil {
		t.Fatalf("delete account-wide policy: %v", err)
	}
	// "You can have one account-level field index policy that applies to
	// all log groups in the account" (PutAccountPolicy): the second
	// account-wide policy under another name rejects — its effective
	// resolution would be listing-order-dependent — and a same-name put
	// stays a replacement.
	if _, err := svc.putAccountPolicyCore("fip-wide", `{"Fields":["B"]}`, "FIELD_INDEX_POLICY", "", "", "us-east-1"); err != nil {
		t.Fatalf("first account-wide policy: %v", err)
	}
	if _, err := svc.putAccountPolicyCore("fip-wide2", `{"Fields":["C"]}`, "FIELD_INDEX_POLICY", "", "", "us-east-1"); err == nil {
		t.Fatal("second account-wide policy under another name accepted")
	}
	if _, err := svc.putAccountPolicyCore("fip-wide", `{"Fields":["C"]}`, "FIELD_INDEX_POLICY", "", "", "us-east-1"); err != nil {
		t.Fatalf("account-wide self-replacement: %v", err)
	}
	if err := svc.deleteAccountPolicyCore("fip-wide", "FIELD_INDEX_POLICY", "us-east-1"); err != nil {
		t.Fatalf("delete account-wide policy again: %v", err)
	}
	for i := 0; i < logsstore.FieldIndexPolicyPrefixQuota; i++ {
		if _, err := svc.putAccountPolicyCore(fmt.Sprintf("fip-q%02d", i), `{"Fields":["A"]}`, "FIELD_INDEX_POLICY", "", fmt.Sprintf("LogGroupNamePrefix = q%02d-", i), "us-east-1"); err != nil {
			t.Fatalf("quota policy %d: %v", i, err)
		}
	}
	_, err := svc.putAccountPolicyCore("fip-over", `{"Fields":["A"]}`, "FIELD_INDEX_POLICY", "", "LogGroupNamePrefix = overflow-", "us-east-1")
	if err == nil || !strings.Contains(fmt.Sprint(err), "LimitExceededException") {
		t.Fatalf("twenty-first prefix policy: %v", err)
	}
}

// TestIndexPolicySDKFacingVocabulary pins the index family's response
// shapes through the handlers: PutIndexPolicy's IndexPolicy member set
// (no policyName at the log group level), DescribeIndexPolicies' list,
// DeleteIndexPolicy's Unit output and DescribeFieldIndexes' FieldIndex
// rows.
func TestIndexPolicySDKFacingVocabulary(t *testing.T) {
	svc, _ := newReadTestService(t, "vocab-index-group")
	reqCtx := request.NewRequestContext(context.Background(), nil, "000000000000", "us-east-1")
	ctx := context.Background()
	req := func(params map[string]interface{}) *request.ParsedRequest {
		return &request.ParsedRequest{Operation: "Vocab", Parameters: params}
	}

	putResp, err := svc.PutIndexPolicy(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifier": "vocab-index-group",
		"policyDocument":     `{"Fields":["RequestId"],"FieldsV2":{"APIName":{"type":"FACET"}}}`,
	}))
	if err != nil {
		t.Fatalf("put index policy: %v", err)
	}
	policy := putResp.(map[string]interface{})["indexPolicy"].(map[string]interface{})
	for k := range policy {
		switch k {
		case "logGroupIdentifier", "policyDocument", "lastUpdateTime", "source":
		default:
			t.Fatalf("IndexPolicy carries the unmodelled member %q", k)
		}
	}
	if policy["source"] != "LOG_GROUP" {
		t.Fatalf("put source: %v", policy["source"])
	}

	descResp, err := svc.DescribeIndexPolicies(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifiers": []interface{}{"vocab-index-group"},
	}))
	if err != nil {
		t.Fatalf("describe index policies: %v", err)
	}
	policies := descResp.(map[string]interface{})["indexPolicies"].([]map[string]interface{})
	if len(policies) != 1 {
		t.Fatalf("describe index policies: %d rows", len(policies))
	}
	for k := range policies[0] {
		switch k {
		case "logGroupIdentifier", "policyDocument", "lastUpdateTime", "source", "policyName":
		default:
			t.Fatalf("IndexPolicy carries the unmodelled member %q", k)
		}
	}

	fieldsResp, err := svc.DescribeFieldIndexes(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifiers": []interface{}{"vocab-index-group"},
	}))
	if err != nil {
		t.Fatalf("describe field indexes: %v", err)
	}
	fieldIndexes := fieldsResp.(map[string]interface{})["fieldIndexes"].([]map[string]interface{})
	if len(fieldIndexes) == 0 {
		t.Fatal("field index listing is empty")
	}
	for k := range fieldIndexes[0] {
		switch k {
		case "fieldIndexName", "logGroupIdentifier", "lastScanTime", "firstEventTime", "lastEventTime", "type", "indexCategory":
		default:
			t.Fatalf("FieldIndex carries the unmodelled member %q", k)
		}
	}

	delResp, err := svc.DeleteIndexPolicy(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifier": "vocab-index-group",
	}))
	if err != nil {
		t.Fatalf("delete index policy: %v", err)
	}
	if len(delResp.(map[string]interface{})) != 0 {
		t.Fatalf("DeleteIndexPolicy response is not empty: %v", delResp)
	}
}

// The family's identifier members ride the shared LogGroupIdentifier
// element traits (1..2048 of "[\w#+=/:,.@-]"): a foreign character
// rejects as InvalidParameterException — and the alphabet excludes '*',
// so PutIndexPolicy's "Don't include an * at the end" rejects in the
// same check instead of the suffix being trimmed away.
func TestIndexPolicyIdentifierTraits(t *testing.T) {
	svc, _ := newIndexPolicyTestEnv(t, "ip-trait-group")

	rows := []struct {
		name string
		call func() error
	}{
		{"put with space", func() error {
			_, err := svc.putIndexPolicyCore("ip trait group", `{"Fields":["RequestId"]}`, "us-east-1")
			return err
		}},
		{"put with asterisk suffix", func() error {
			_, err := svc.putIndexPolicyCore("arn:aws:logs:us-east-1:000000000000:log-group:ip-trait-group:*", `{"Fields":["RequestId"]}`, "us-east-1")
			return err
		}},
		{"delete with space", func() error {
			return svc.deleteIndexPolicyCore("ip trait group", "us-east-1")
		}},
		{"describe with space", func() error {
			_, err := svc.describeIndexPoliciesCore([]string{"ip trait group"}, "us-east-1")
			return err
		}},
		{"field indexes with space", func() error {
			_, _, err := svc.describeFieldIndexesCore([]string{"ip trait group"}, []string{"AUTO"}, "us-east-1", "")
			return err
		}},
	}
	for _, row := range rows {
		if err := row.call(); logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("%s: err = %v, want InvalidParameterException", row.name, err)
		}
	}

	// The ARN form passes the traits and addresses the group.
	if _, err := svc.putIndexPolicyCore("arn:aws:logs:us-east-1:000000000000:log-group:ip-trait-group", `{"Fields":["RequestId"]}`, "us-east-1"); err != nil {
		t.Fatalf("ARN-form identifier rejected: %v", err)
	}
}
