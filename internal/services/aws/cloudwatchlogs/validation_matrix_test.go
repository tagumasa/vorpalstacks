package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The C2 validation-matrix pin: every row the member-trait census derives
// from the vendored model — required flags, enumerations, lengths,
// patterns and the documented byte limits — is exercised against the Core
// that owns it, with the boundary acceptances alongside the rejections.

// --- PutMetricFilter: metricTransformations is exactly 1..1 ---

func TestPutMetricFilterExactlyOneTransformation(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-metric-group")
	one := []logsstore.MetricTransformation{{
		MetricName: "Count", MetricNamespace: "App", MetricValue: "1",
		MetricNameSet: true, MetricNamespaceSet: true, MetricValueSet: true,
	}}

	if err := svc.putMetricFilterCore("matrix-metric-group", "one", "ERROR", true, one, false, "", nil, "us-east-1"); err != nil {
		t.Fatalf("single transformation: %v", err)
	}
	two := append(one, logsstore.MetricTransformation{MetricName: "Other", MetricNamespace: "App", MetricValue: "1"})
	err := svc.putMetricFilterCore("matrix-metric-group", "two", "ERROR", true, two, false, "", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("two transformations: code=%q err=%v", code, err)
	}
	err = svc.putMetricFilterCore("matrix-metric-group", "zero", "ERROR", true, nil, false, "", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("zero transformations: code=%q err=%v", code, err)
	}
	// The required members whose shapes allow the empty string reject on
	// absence alone: a missing filterPattern, and a transformation missing
	// any of metricName/metricNamespace/metricValue.
	err = svc.putMetricFilterCore("matrix-metric-group", "absent", "", false, one, false, "", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("absent filterPattern: code=%q err=%v", code, err)
	}
	noName := []logsstore.MetricTransformation{one[0]}
	noName[0].MetricNameSet = false
	err = svc.putMetricFilterCore("matrix-metric-group", "no-name", "ERROR", true, noName, false, "", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("absent metricName: code=%q err=%v", code, err)
	}
	noValue := []logsstore.MetricTransformation{one[0]}
	noValue[0].MetricValueSet = false
	err = svc.putMetricFilterCore("matrix-metric-group", "no-value", "ERROR", true, noValue, false, "", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("absent metricValue: code=%q err=%v", code, err)
	}
	// A present empty pattern is the match-everything form and stays
	// legal.
	if err := svc.putMetricFilterCore("matrix-metric-group", "empty-pattern", "", true, one, false, "", nil, "us-east-1"); err != nil {
		t.Fatalf("present empty pattern: %v", err)
	}
}

// --- Tag operations: Tags 1-50, TagKey/TagValue traits, tagKeys min 0 ---

func testTags(n int) map[string]string {
	tags := make(map[string]string, n)
	for i := 0; i < n; i++ {
		tags[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("value-%d", i)
	}
	return tags
}

func TestTagLogGroupValidationRows(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-tag-group")

	if err := svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-tag-group", Tags: testTags(50), Region: "us-east-1"}); err != nil {
		t.Fatalf("50 tags at the ceiling: %v", err)
	}
	err := svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-tag-group", Tags: testTags(51), Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("51 tags: code=%q err=%v", code, err)
	}
	err = svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-tag-group", Tags: map[string]string{}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty tags: code=%q err=%v", code, err)
	}
	err = svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-tag-group", Tags: map[string]string{"bad!key": "v"}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("key outside the charset: code=%q err=%v", code, err)
	}
	err = svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-tag-group", Tags: map[string]string{"k": strings.Repeat("v", 257)}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("value over 256: code=%q err=%v", code, err)
	}
	err = svc.untagLogGroupCore(UntagLogGroupInput{LogGroupName: "matrix-tag-group", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty untag list (TagList min 1): code=%q err=%v", code, err)
	}
}

func tagResourceRequest(arn string, tags map[string]interface{}) *request.ParsedRequest {
	return &request.ParsedRequest{
		Operation:  "TagResource",
		Parameters: map[string]interface{}{"resourceArn": arn, "tags": tags},
	}
}

func tagsParam(n int) map[string]interface{} {
	tags := make(map[string]interface{}, n)
	for i := 0; i < n; i++ {
		tags[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("value-%d", i)
	}
	return tags
}

func TestTagResourceValidationRows(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-tagres-group")
	lg, err := store.GetLogGroup("matrix-tagres-group")
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	cfg := svc.tagHandlerConfig(store)
	ctx := context.Background()

	if _, err := tagutil.HandleTag(ctx, tagResourceRequest(lg.ARN, tagsParam(50)), cfg); err != nil {
		t.Fatalf("50 tags at the ceiling: %v", err)
	}
	_, err = tagutil.HandleTag(ctx, tagResourceRequest(lg.ARN, tagsParam(51)), cfg)
	if code := logsErrorCode(err); code != "TooManyTagsException" {
		t.Fatalf("51 tags: code=%q err=%v", code, err)
	}
	_, err = tagutil.HandleTag(ctx, tagResourceRequest(lg.ARN, map[string]interface{}{}), cfg)
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty tags: code=%q err=%v", code, err)
	}
	_, err = tagutil.HandleTag(ctx, tagResourceRequest(lg.ARN, map[string]interface{}{"no!bang": "v"}), cfg)
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("key outside the charset: code=%q err=%v", code, err)
	}

	// tagKeys carries "Minimum number of 0 items" — an empty list is a
	// valid no-op, not a rejection.
	untagReq := &request.ParsedRequest{
		Operation:  "UntagResource",
		Parameters: map[string]interface{}{"resourceArn": lg.ARN, "tagKeys": []interface{}{}},
	}
	if _, err := tagutil.HandleUntag(ctx, untagReq, cfg); err != nil {
		t.Fatalf("empty tagKeys no-op: %v", err)
	}
}

// --- PutLogGroupDeletionProtection: deletionProtectionEnabled required ---

func TestPutLogGroupDeletionProtectionRequiredMember(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-deletion-group")

	err := svc.putLogGroupDeletionProtectionCore(PutLogGroupDeletionProtectionInput{
		LogGroupIdentifier: "matrix-deletion-group", DeletionProtectionSet: false, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("omitted member: code=%q err=%v", code, err)
	}
	if err := svc.putLogGroupDeletionProtectionCore(PutLogGroupDeletionProtectionInput{
		LogGroupIdentifier: "matrix-deletion-group", DeletionProtectionEnabled: false, DeletionProtectionSet: true, Region: "us-east-1"}); err != nil {
		t.Fatalf("explicit false: %v", err)
	}
}

// --- CreateImportTask: importRoleArn required ---

func TestCreateImportTaskRequiresRoleArn(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-import-group")
	validSource := "arn:aws:s3:::matrix-import-bucket/import.log"

	_, err := svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: validSource, Region: "us-east-1", AccountID: "000000000000"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("missing role: code=%q err=%v", code, err)
	}
	_, err = svc.createImportTaskCore(store, &CreateImportTaskInput{
		ImportSourceArn: validSource, ImportRoleArn: "not-an-arn", Region: "us-east-1", AccountID: "000000000000"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("malformed role: code=%q err=%v", code, err)
	}
}

// --- PutDataProtectionPolicy: policyDocument required, JSON-checked ---

func TestPutDataProtectionPolicyDocumentValidation(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-dpp-group")

	_, err := svc.putDataProtectionPolicyCore("matrix-dpp-group", "", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty document: code=%q err=%v", code, err)
	}
	_, err = svc.putDataProtectionPolicyCore("matrix-dpp-group", "{not json", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("non-JSON document: code=%q err=%v", code, err)
	}

	// The two-block structure: the Audit block (DataIdentifer array plus
	// an Operation whose Audit action contains FindingsDestination) and
	// the Deidentify block (the matching DataIdentifer array plus an
	// Operation whose Deidentify action contains the empty MaskConfig),
	// over the optional metadata fields.
	const valid = `{
		"Name": "DataProtectionPolicy",
		"Description": "mask emails",
		"Version": "2021-06-01",
		"Statement": [
			{"DataIdentifer": ["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],
			 "Operation": {"Audit": {"FindingsDestination": {"CloudWatchLogs": {"LogGroup": "audit-findings"}}}}},
			{"DataIdentifer": ["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],
			 "Operation": {"Deidentify": {"MaskConfig": {}}}}
		]
	}`
	if _, err := svc.putDataProtectionPolicyCore("matrix-dpp-group", valid, "us-east-1"); err != nil {
		t.Fatalf("two-block document: %v", err)
	}

	// Every documented malformation rejects: one block alone, a missing
	// DataIdentifer array, a missing FindingsDestination, mismatched
	// DataIdentifer arrays, a missing and a non-empty MaskConfig.
	for _, tc := range []struct{ name, doc string }{
		{"one block", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}}]}`},
		{"missing audit array", `{"Statement":[{"Operation":{"Audit":{"FindingsDestination":{}}}},{"DataIdentifer":["i"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`},
		{"missing findings destination", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{}}},{"DataIdentifer":["i"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`},
		{"mismatched arrays", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}},{"DataIdentifer":["j"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`},
		{"missing mask config", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}},{"DataIdentifer":["i"],"Operation":{"Deidentify":{}}}]}`},
		{"non-empty mask config", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}},{"DataIdentifer":["i"],"Operation":{"Deidentify":{"MaskConfig":{"Mode":"overwrite"}}}}]}`},
		{"swapped block order", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Deidentify":{"MaskConfig":{}}}},{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}}]}`},
	} {
		_, err := svc.putDataProtectionPolicyCore("matrix-dpp-group", tc.doc, "us-east-1")
		if code := logsErrorCode(err); code != "InvalidParameterException" {
			t.Fatalf("%s: code=%q err=%v", tc.name, code, err)
		}
	}

	// The document carries its own ceiling, tighter than the generic
	// policy-document cap.
	oversize := `{"Name":"` + strings.Repeat("n", logsstore.MaxDataProtectionPolicyDocumentBytes) + `"}`
	if _, err := svc.putDataProtectionPolicyCore("matrix-dpp-group", oversize, "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("oversize document: %v", err)
	}
}

// --- PutQueryDefinition: queryString min length 1 ---

func TestPutQueryDefinitionQueryStringMinLength(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-qd-group")

	_, err := svc.putQueryDefinitionCore("qd", "", "CWLI", "", nil, nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty queryString: code=%q err=%v", code, err)
	}
}

// --- TestMetricFilter: logEventMessages min 1 ---

func TestTestMetricFilterMessagesRequired(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-tmf-group")

	_, err := svc.testMetricFilterCore(testMetricFilterInput{FilterPattern: "ERROR", FilterPatternSet: true})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty messages: code=%q err=%v", code, err)
	}
	if _, err := svc.testMetricFilterCore(testMetricFilterInput{FilterPattern: "ERROR", FilterPatternSet: true, LogEventMessages: []string{"ERROR one"}}); err != nil {
		t.Fatalf("one message: %v", err)
	}
	_, err = svc.testMetricFilterCore(testMetricFilterInput{})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("missing pattern: code=%q err=%v", code, err)
	}
}

// --- PutLogEvents: required members and the documented byte limits ---

func TestPutLogEventsRequiredMembersAndByteLimits(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-events-group")
	if err := svc.createLogStreamCore(CreateLogStreamInput{LogGroupName: "matrix-events-group", LogStreamName: "s", Region: "us-east-1"}); err != nil {
		t.Fatalf("stream: %v", err)
	}
	now := time.Now().UnixMilli()
	put := func(events ...PutLogEvent) error {
		_, err := svc.putLogEventsCore(PutLogEventsInput{LogGroupName: "matrix-events-group", LogStreamName: "s", Region: "us-east-1", Events: events})
		return err
	}

	err := put(PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now, Message: "no timestamp flag"}, TimestampSet: false})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("absent timestamp: code=%q err=%v", code, err)
	}
	err = put(PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now}, TimestampSet: true})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("empty message: code=%q err=%v", code, err)
	}
	err = put(PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now, Message: strings.Repeat("x", logsstore.MaxLogEventMessageBytes+1)}, TimestampSet: true})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("message over 1 MB: code=%q err=%v", code, err)
	}
	// A present timestamp of zero is a value, not an absence: the age
	// rules reject that event per event, not the request.
	res, err := svc.putLogEventsCore(PutLogEventsInput{LogGroupName: "matrix-events-group", LogStreamName: "s", Region: "us-east-1",
		Events: []PutLogEvent{{LogEntry: logsstore.LogEntry{Message: "epoch zero"}, TimestampSet: true}}})
	if err != nil {
		t.Fatalf("epoch-zero event: %v", err)
	}
	if res.RejectedLogEvents == nil {
		t.Fatalf("epoch-zero event should reject per event, got %+v", res.RejectedLogEvents)
	}

	// Boundaries: the batch ceiling minus the per-event overhead is the
	// largest single message both documented caps admit (a full 1 MB
	// message can never fit the 1 MB batch formula — 1,048,576 message
	// bytes + 26 overhead exceeds 1,048,576 — a tension the two
	// documented sentences genuinely carry, pinned right after).
	if err := put(PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now, Message: strings.Repeat("y", logsstore.MaxPutLogEventsBatchBytes-logsstore.PutLogEventOverheadBytes)}, TimestampSet: true}); err != nil {
		t.Fatalf("message at the batch-admitted maximum: %v", err)
	}
	err = put(PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now, Message: strings.Repeat("y", logsstore.MaxLogEventMessageBytes)}, TimestampSet: true})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("full 1 MB message trips the batch formula: code=%q err=%v", code, err)
	}
	err = put(
		PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now, Message: strings.Repeat("a", 600000)}, TimestampSet: true},
		PutLogEvent{LogEntry: logsstore.LogEntry{Timestamp: now + 1, Message: strings.Repeat("b", 600000)}, TimestampSet: true},
	)
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("batch over 1 MB: code=%q err=%v", code, err)
	}
	_ = store
}

// --- Enumeration members ---

func TestDescribeLogStreamsOrderByEnum(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-order-group")
	if _, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{LogGroupName: "matrix-order-group", OrderBy: "Bogus", Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("garbage orderBy: code=%q err=%v", logsErrorCode(err), err)
	}
	if _, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{LogGroupName: "matrix-order-group", OrderBy: "LastEventTime", Region: "us-east-1"}); err != nil {
		t.Fatalf("LastEventTime: %v", err)
	}
}

func TestDescribeQueriesEnumFilters(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-dq-group")
	if _, _, err := svc.describeQueriesCore(&DescribeQueriesInput{StatusFilter: "Bogus"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("garbage status: %v", err)
	}
	if _, _, err := svc.describeQueriesCore(&DescribeQueriesInput{QueryLanguage: "CWLIX"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("garbage queryLanguage: %v", err)
	}
	if _, _, err := svc.describeQueriesCore(&DescribeQueriesInput{StatusFilter: "Scheduled"}); err != nil {
		t.Fatalf("Scheduled status: %v", err)
	}
}

func TestDescribeTaskStatusEnumFilters(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-taskstatus-group")
	if _, _, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{StatusCode: "BOGUS"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("garbage export statusCode: %v", err)
	}
	if _, _, err := svc.describeExportTasksCore(store, &DescribeExportTasksInput{StatusCode: "RUNNING"}); err != nil {
		t.Fatalf("RUNNING statusCode: %v", err)
	}
	if _, _, err := svc.describeImportTasksCore(store, "", "BOGUS", "", "", 0); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("garbage importStatus: %v", err)
	}
	if _, _, err := svc.describeImportTasksCore(store, "", "COMPLETED", "", "", 0); err != nil {
		t.Fatalf("COMPLETED importStatus: %v", err)
	}
}

func TestGetScheduledQueryHistoryExecutionStatusEnum(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-hist-group")
	if _, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier: "sched", ExecutionStatuses: []string{"Bogus"}}); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("garbage executionStatus: %v", err)
	}
}

// --- DescribeLogGroups: pattern substring, cross-member rules, accounts ---

func TestDescribeLogGroupsPatternAndCrossMemberRules(t *testing.T) {
	svc, store := newReadTestService(t, "DataLogs")
	for _, name := range []string{"aws/DataLogs", "GroupDataLogs", "datalogs", "Groupdata", "Data/log/s"} {
		if err := store.CreateLogGroup(logsstore.NewLogGroup(name, "us-east-1", "000000000000")); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}

	// Case-sensitive substring: DataLogs, aws/DataLogs, GroupDataLogs.
	res, err := svc.describeLogGroupsCore(ListLogGroupsInput{LogGroupNamePattern: "DataLogs", Region: "us-east-1"})
	if err != nil {
		t.Fatalf("pattern listing: %v", err)
	}
	if len(res.LogGroups) != 3 {
		t.Fatalf("substring matches: want 3, got %d", len(res.LogGroups))
	}

	// Mutual exclusion with the prefix member.
	_, err = svc.describeLogGroupsCore(ListLogGroupsInput{LogGroupNamePattern: "Data", LogGroupNamePrefix: "Group", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("pattern+prefix: code=%q err=%v", code, err)
	}
	// logGroupIdentifiers admits includeLinkedAccounts alone.
	for _, tc := range []ListLogGroupsInput{
		{LogGroupIdentifiers: []string{"datalogs"}, LogGroupNamePrefix: "Group"},
		{LogGroupIdentifiers: []string{"datalogs"}, LogGroupNamePattern: "Data"},
		{LogGroupIdentifiers: []string{"datalogs"}, LogGroupClass: "STANDARD"},
	} {
		tc.Region = "us-east-1"
		if _, err := svc.describeLogGroupsCore(tc); logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("identifiers with a fellow filter: %v", err)
		}
	}
	res, err = svc.describeLogGroupsCore(ListLogGroupsInput{
		LogGroupIdentifiers:   []string{"datalogs", "arn:aws:logs:us-east-1:000000000000:log-group:Groupdata"},
		IncludeLinkedAccounts: true, Region: "us-east-1"})
	if err != nil {
		t.Fatalf("identifiers with includeLinkedAccounts: %v", err)
	}
	if len(res.LogGroups) != 2 {
		t.Fatalf("identifier selection (name and ARN forms): want 2, got %d", len(res.LogGroups))
	}

	// Account scoping: accountIdentifiers applies only when
	// includeLinkedAccounts is set; a list without the local account owns
	// nothing on this single-account platform.
	res, err = svc.describeLogGroupsCore(ListLogGroupsInput{
		AccountIdentifiers: []string{"999999999999"}, IncludeLinkedAccounts: true, Region: "us-east-1"})
	if err != nil {
		t.Fatalf("foreign account scope: %v", err)
	}
	if len(res.LogGroups) != 0 {
		t.Fatalf("foreign account scope: want 0, got %d", len(res.LogGroups))
	}
	res, err = svc.describeLogGroupsCore(ListLogGroupsInput{
		AccountIdentifiers: []string{"000000000000"}, IncludeLinkedAccounts: true, Region: "us-east-1"})
	if err != nil {
		t.Fatalf("local account scope: %v", err)
	}
	if len(res.LogGroups) != 6 {
		t.Fatalf("local account scope: want 6, got %d", len(res.LogGroups))
	}
	res, err = svc.describeLogGroupsCore(ListLogGroupsInput{
		AccountIdentifiers: []string{"999999999999"}, Region: "us-east-1"})
	if err != nil {
		t.Fatalf("accountIdentifiers without the include flag: %v", err)
	}
	if len(res.LogGroups) != 6 {
		t.Fatalf("accountIdentifiers without the include flag is not a filter: want 6, got %d", len(res.LogGroups))
	}
}

// --- ListLogGroups pattern alternations ---

func TestListLogGroupsPatternAlternations(t *testing.T) {
	cases := []struct {
		pattern string
		name    string
		want    bool
	}{
		{"DataLogs", "aws/DataLogs", true},
		{"DataLogs", "GroupDataLogs", true},
		{"DataLogs", "datalogs", false},
		{"^/aws", "/aws/lambda", true},
		{"^/aws", "x/aws/lambda", false},
		{"api|web", "my-web-group", true},
		{"api|web", "my-api-group", true},
		{"api|web", "other", false},
		{"^api|^web", "web-front", true},
		{"^api|^web", "back-web", false},
	}
	for _, tc := range cases {
		m, err := parseLogGroupNameMatcher(tc.pattern)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.pattern, err)
		}
		if got := m.matches(tc.name); got != tc.want {
			t.Fatalf("match(%q, %q) = %v, want %v", tc.pattern, tc.name, got, tc.want)
		}
	}

	for _, bad := range []string{"ab", strings.Repeat("x", 25), "a?b", "api|", "|api", strings.Repeat("alt|", 5) + "alt"} {
		if _, err := parseLogGroupNameMatcher(bad); logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("pattern %q should reject: %v", bad, err)
		}
	}
	// Boundaries: five alternations and 24 characters are the ceilings.
	if _, err := parseLogGroupNameMatcher("aaa|bbb|ccc|ddd|eee"); err != nil {
		t.Fatalf("five alternations: %v", err)
	}
	if _, err := parseLogGroupNameMatcher(strings.Repeat("x", 24)); err != nil {
		t.Fatalf("24 characters: %v", err)
	}
}

// --- ListLogGroups: bare pattern matches by substring through the handler ---

func TestListLogGroupsBarePatternSubstrings(t *testing.T) {
	svc, store := newReadTestService(t, "DataLogs")
	for _, name := range []string{"aws/DataLogs", "GroupDataLogs", "datalogs", "Groupdata"} {
		if err := store.CreateLogGroup(logsstore.NewLogGroup(name, "us-east-1", "000000000000")); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}
	reqCtx := request.NewRequestContext(context.Background(), svc.storageManager, "000000000000", "us-east-1")
	resp, err := svc.ListLogGroups(context.Background(), reqCtx, &request.ParsedRequest{
		Parameters: map[string]interface{}{"logGroupNamePattern": "DataLogs"},
	})
	if err != nil {
		t.Fatalf("list with bare pattern: %v", err)
	}
	groups, _ := resp.(map[string]interface{})["logGroups"].([]map[string]interface{})
	names := make(map[string]bool, len(groups))
	for _, g := range groups {
		names[g["logGroupName"].(string)] = true
	}
	// The single bare alternative matches by substring, exactly as the
	// multi-alternative form does: the mid-name and suffixed groups must
	// appear, the case-variant must not.
	for _, want := range []string{"DataLogs", "aws/DataLogs", "GroupDataLogs"} {
		if !names[want] {
			t.Fatalf("bare pattern: %q missing from response (got %v)", want, names)
		}
	}
	if names["datalogs"] || names["Groupdata"] {
		t.Fatalf("bare pattern matched non-substring names: %v", names)
	}
	// The legacy prefix fallback keeps prefix semantics on the same
	// handler: only the prefixed group returns.
	resp, err = svc.ListLogGroups(context.Background(), reqCtx, &request.ParsedRequest{
		Parameters: map[string]interface{}{"logGroupNamePrefix": "Group"},
	})
	if err != nil {
		t.Fatalf("list with legacy prefix: %v", err)
	}
	groups, _ = resp.(map[string]interface{})["logGroups"].([]map[string]interface{})
	for _, g := range groups {
		if name := g["logGroupName"].(string); name != "GroupDataLogs" && name != "Groupdata" {
			t.Fatalf("legacy prefix returned non-prefixed group %q", name)
		}
	}
}

// --- PutResourcePolicy: revisionId issuance ---

func TestPutResourcePolicyQuotaAndRevisionRules(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-rp-group")
	doc := `{"Version":"2012-10-17","Statement":[]}`
	lg, err := store.GetLogGroup("matrix-rp-group")
	if err != nil {
		t.Fatal(err)
	}
	groupArn := lg.ARN

	// The account-wide form carries no revision identity: "The revision ID
	// of the created or updated resource policy. Only returned for
	// resource-scoped policies."
	first, err := svc.putResourcePolicyCore("p1", doc, "", "", "us-east-1")
	if err != nil {
		t.Fatalf("account-wide put: %v", err)
	}
	if first.RevisionId != "" {
		t.Fatalf("account-wide put issued a revisionId %q", first.RevisionId)
	}
	if _, err := svc.putResourcePolicyCore("p1", doc, "", "", "us-east-1"); err != nil {
		t.Fatalf("account-wide re-put needs no revision: %v", err)
	}

	// The resource-scoped form: "Required when resourceArn is provided to
	// prevent concurrent modifications. Use null when creating a resource
	// policy for the first time" — the first-time null creates and issues
	// a revision; an update must quote it; a stale quote rejects.
	rp1, err := svc.putResourcePolicyCore("rp1", doc, groupArn, "", "us-east-1")
	if err != nil {
		t.Fatalf("first-time resource-scoped put: %v", err)
	}
	if rp1.RevisionId == "" {
		t.Fatal("resource-scoped put issued no revisionId")
	}
	_, err = svc.putResourcePolicyCore("rp1", doc, groupArn, "", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("update without a revision quote: code=%q err=%v", code, err)
	}
	_, err = svc.putResourcePolicyCore("rp1", doc, groupArn, "rev-stale", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("stale revision quote: code=%q err=%v", code, err)
	}
	rp2, err := svc.putResourcePolicyCore("rp1", doc, groupArn, rp1.RevisionId, "us-east-1")
	if err != nil {
		t.Fatalf("current revision quote: %v", err)
	}
	if rp2.RevisionId == rp1.RevisionId {
		t.Fatalf("successful re-put carried the old revisionId %q", rp2.RevisionId)
	}

	// "One per LogGroup resourceARN": another name on the same ARN
	// rejects, while the policy's own name keeps its slot.
	_, err = svc.putResourcePolicyCore("rp-other", doc, groupArn, "", "us-east-1")
	if code := logsErrorCode(err); code != "LimitExceededException" {
		t.Fatalf("second policy on one resourceArn: code=%q err=%v", code, err)
	}
	if _, err := svc.putResourcePolicyCore("rp1", doc, groupArn, rp2.RevisionId, "us-east-1"); err != nil {
		t.Fatalf("replacing keeps its slot: %v", err)
	}

	// "An account can have a maximum of 10 policies without resourceARN":
	// p1 plus nine more fill the ceiling and the eleventh rejects.
	for i := 0; i < logsstore.MaxAccountScopeResourcePolicies-1; i++ {
		if _, err := svc.putResourcePolicyCore(fmt.Sprintf("acct-%d", i), doc, "", "", "us-east-1"); err != nil {
			t.Fatalf("account-wide policy %d: %v", i, err)
		}
	}
	_, err = svc.putResourcePolicyCore("acct-overflow", doc, "", "", "us-east-1")
	if code := logsErrorCode(err); code != "LimitExceededException" {
		t.Fatalf("eleventh account-wide policy: code=%q err=%v", code, err)
	}
}

// --- Resource-policy delete and listing contract ---

// DeleteResourcePolicy carries the resourceArn member and the
// expectedRevisionId rule its API reference documents ("Required when
// deleting a resource-scoped policy to prevent concurrent modifications");
// DescribeResourcePolicies' policyScope documents "Valid values are ACCOUNT
// or RESOURCE. When not specified, defaults to ACCOUNT."
func TestResourcePolicyDeleteAndScopeContract(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-rpd-group")
	doc := `{"Version":"2012-10-17","Statement":[]}`
	lg, err := store.GetLogGroup("matrix-rpd-group")
	if err != nil {
		t.Fatal(err)
	}
	groupArn := lg.ARN

	if _, err := svc.putResourcePolicyCore("wide", doc, "", "", "us-east-1"); err != nil {
		t.Fatalf("account-wide put: %v", err)
	}
	rp, err := svc.putResourcePolicyCore("scoped", doc, groupArn, "", "us-east-1")
	if err != nil {
		t.Fatalf("resource-scoped put: %v", err)
	}

	// The resource-scoped delete quotes its revision or rejects.
	err = svc.deleteResourcePolicyCore("scoped", "", "", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("resource-scoped delete without a revision quote: code=%q err=%v", code, err)
	}
	err = svc.deleteResourcePolicyCore("scoped", "", "rev-stale", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("resource-scoped delete with a stale quote: code=%q err=%v", code, err)
	}
	if err := svc.deleteResourcePolicyCore("scoped", "", rp.RevisionId, "us-east-1"); err != nil {
		t.Fatalf("resource-scoped delete with the current quote: %v", err)
	}
	// The account-wide form carries no revision identity and needs none.
	if err := svc.deleteResourcePolicyCore("wide", "", "", "us-east-1"); err != nil {
		t.Fatalf("account-wide delete without a quote: %v", err)
	}

	// The resourceArn form: "Currently only supports LogGroup ARN."
	for name, arn := range map[string]string{
		"foreign service":    "arn:aws:s3:::bucket",
		"non-group resource": "arn:aws:logs:us-east-1:000000000000:destination:d1",
		"stream namespace":   "arn:aws:logs:us-east-1:000000000000:log-group:matrix-rpd-group:log-stream:s1",
		"not an ARN":         "matrix-rpd-group",
	} {
		_, err := svc.putResourcePolicyCore("bad-arn", doc, arn, "", "us-east-1")
		if code := logsErrorCode(err); code != "InvalidParameterException" {
			t.Fatalf("%s resourceArn: code=%q err=%v", name, code, err)
		}
	}

	// The listing scope: absent defaults to ACCOUNT, the enum rejects a
	// foreign value, and RESOURCE selects the resource-scoped form.
	_, err = svc.putResourcePolicyCore("wide2", doc, "", "", "us-east-1")
	if err != nil {
		t.Fatalf("account-wide re-put: %v", err)
	}
	_, err = svc.putResourcePolicyCore("scoped2", doc, groupArn, "", "us-east-1")
	if err != nil {
		t.Fatalf("resource-scoped re-put: %v", err)
	}
	def, _, err := svc.describeResourcePoliciesCore("", "", "", "us-east-1", 0)
	if err != nil {
		t.Fatalf("default-scope describe: %v", err)
	}
	for _, p := range def {
		if p.PolicyScope != "ACCOUNT" {
			t.Fatalf("default scope listed a %s policy named %q", p.PolicyScope, p.PolicyName)
		}
	}
	res, _, err := svc.describeResourcePoliciesCore("", "RESOURCE", "", "us-east-1", 0)
	if err != nil {
		t.Fatalf("RESOURCE-scope describe: %v", err)
	}
	if len(res) != 1 || res[0].PolicyName != "scoped2" {
		t.Fatalf("RESOURCE scope returned %d policies, want the scoped one alone", len(res))
	}
	_, _, err = svc.describeResourcePoliciesCore("", "BOTH", "", "us-east-1", 0)
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("foreign policyScope: code=%q err=%v", code, err)
	}
}

// The admission serialises: concurrent Puts quoting the same revision
// admit exactly one, and the account-wide census ceiling holds under a
// concurrent burst. (The pre-fix check-then-put race was structural —
// the store's write carried no lock and the census decided on an
// unlocked read — so the pin guards the regression rather than
// demonstrating the race.)
func TestResourcePolicyAdmissionSerialises(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-rps-group")
	doc := `{"Version":"2012-10-17","Statement":[]}`
	lg, err := store.GetLogGroup("matrix-rps-group")
	if err != nil {
		t.Fatal(err)
	}
	groupArn := lg.ARN

	rp, err := svc.putResourcePolicyCore("raced", doc, groupArn, "", "us-east-1")
	if err != nil {
		t.Fatalf("first-time resource-scoped put: %v", err)
	}
	var wg sync.WaitGroup
	succeeded := make(chan *logsstore.ResourcePolicy, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got, err := svc.putResourcePolicyCore("raced", doc, groupArn, rp.RevisionId, "us-east-1"); err == nil {
				succeeded <- got
			}
		}()
	}
	wg.Wait()
	close(succeeded)
	winners := 0
	var winnerRevision string
	for got := range succeeded {
		winners++
		winnerRevision = got.RevisionId
	}
	if winners != 1 {
		t.Fatalf("%d concurrent puts quoting one revision succeeded, want exactly one", winners)
	}
	stored, err := store.GetResourcePolicy("raced")
	if err != nil {
		t.Fatal(err)
	}
	if stored.RevisionId != winnerRevision {
		t.Fatalf("stored revision %q is not the admitted winner's %q", stored.RevisionId, winnerRevision)
	}

	// The census under a burst: with one slot free, a concurrent volley
	// for distinct names admits exactly one more policy and the ceiling
	// stands at the documented count afterwards.
	for i := 0; i < logsstore.MaxAccountScopeResourcePolicies-1; i++ {
		if _, err := svc.putResourcePolicyCore(fmt.Sprintf("fill-%d", i), doc, "", "", "us-east-1"); err != nil {
			t.Fatalf("fill policy %d: %v", i, err)
		}
	}
	var censusWg sync.WaitGroup
	censusSucceeded := make(chan string, 6)
	for i := 0; i < 6; i++ {
		censusWg.Add(1)
		go func(i int) {
			defer censusWg.Done()
			if _, err := svc.putResourcePolicyCore(fmt.Sprintf("burst-%d", i), doc, "", "", "us-east-1"); err == nil {
				censusSucceeded <- fmt.Sprintf("burst-%d", i)
			}
		}(i)
	}
	censusWg.Wait()
	close(censusSucceeded)
	if len(censusSucceeded) != 1 {
		t.Fatalf("%d burst puts admitted into one free census slot, want exactly one", len(censusSucceeded))
	}
	all, err := store.ListResourcePolicies("")
	if err != nil {
		t.Fatal(err)
	}
	accountScoped := 0
	for _, p := range all {
		if p.ResourceArn == "" {
			accountScoped++
		}
	}
	if accountScoped != logsstore.MaxAccountScopeResourcePolicies {
		t.Fatalf("account-wide census counted %d policies, ceiling %d", accountScoped, logsstore.MaxAccountScopeResourcePolicies)
	}
}

// The secondary paths carry their member traits: ListLogGroups'
// accountIdentifiers rides the same AccountIds 0-20 bound DescribeLogGroups
// enforces; the deletion-protection and bearer-token logGroupIdentifier
// members carry the shared shape's per-element traits; and the KmsKeyId
// length ceiling applies on every kmsKeyId surface, not the lookup-table
// and destination paths alone.
func TestSecondaryPathTraitGaps(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-spt-group")

	// ListLogGroups: "You can specify as many as 20 account IDs in the
	// array" — the shared AccountIds shape's bound, on the list path too.
	ids := make([]string, logsstore.MaxListAccountIdentifiers+1)
	for i := range ids {
		ids[i] = "000000000000"
	}
	if _, err := svc.listLogGroupsCore(ListLogGroupsInput{AccountIdentifiers: ids, Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("listLogGroups with %d account identifiers: code=%q err=%v", len(ids), logsErrorCode(err), err)
	}

	// The identifier members: over-alphabet and over-length input is a
	// parameter error, not a not-found.
	badIdentifier := "not!an@identifier!"
	if err := svc.putLogGroupDeletionProtectionCore(PutLogGroupDeletionProtectionInput{
		LogGroupIdentifier: badIdentifier, DeletionProtectionEnabled: true, DeletionProtectionSet: true, Region: "us-east-1",
	}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("deletion protection with a foreign identifier: code=%q err=%v", logsErrorCode(err), err)
	}
	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier: badIdentifier, BearerTokenAuthenticationEnabled: true, BearerTokenAuthenticationSet: true, Region: "us-east-1",
	}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("bearer token with a foreign identifier: code=%q err=%v", logsErrorCode(err), err)
	}

	// The KmsKeyId length ceiling on the group surfaces: an alias ARN
	// past the bound rejects on CreateLogGroup and AssociateKmsKey.
	longAlias := fmt.Sprintf("arn:aws:kms:us-east-1:000000000000:alias/%s", strings.Repeat("a", logsstore.MaxKmsKeyIdLength))
	if _, err := svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "matrix-spt-kms", KmsKeyId: longAlias, Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("createLogGroup with an over-bound kmsKeyId: code=%q err=%v", logsErrorCode(err), err)
	}
	if err := svc.associateKmsKeyCore(AssociateKmsKeyInput{LogGroupName: "matrix-spt-group", KmsKeyId: longAlias, Region: "us-east-1"}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("associateKmsKey with an over-bound kmsKeyId: code=%q err=%v", logsErrorCode(err), err)
	}
}

// The LogGroup shape's configuration members ride the group record:
// dataProtectionStatus follows the protection-policy lifecycle (empty
// before the first policy, ACTIVATED while one is stored, DELETED after
// its deletion) and bearerTokenAuthenticationEnabled follows the switch.
func TestGroupRecordConfigMembers(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-grc-group")
	const validDPP = `{"Statement":[` +
		`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Audit":{"FindingsDestination":{"CloudWatchLogs":{"LogGroup":"audit"}}}}},` +
		`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`

	untouched, err := store.GetLogGroup("matrix-grc-group")
	if err != nil {
		t.Fatal(err)
	}
	if untouched.DataProtectionStatus != "" || untouched.BearerTokenAuthenticationEnabled {
		t.Fatalf("untouched group carries status %q bearer=%v, want empty/false", untouched.DataProtectionStatus, untouched.BearerTokenAuthenticationEnabled)
	}

	if _, err := svc.putDataProtectionPolicyCore("matrix-grc-group", validDPP, "us-east-1"); err != nil {
		t.Fatalf("put data protection policy: %v", err)
	}
	activated, err := store.GetLogGroup("matrix-grc-group")
	if err != nil {
		t.Fatal(err)
	}
	if activated.DataProtectionStatus != "ACTIVATED" {
		t.Fatalf("after put: status %q, want ACTIVATED", activated.DataProtectionStatus)
	}

	if err := svc.deleteDataProtectionPolicyCore("matrix-grc-group", "us-east-1"); err != nil {
		t.Fatalf("delete data protection policy: %v", err)
	}
	deleted, err := store.GetLogGroup("matrix-grc-group")
	if err != nil {
		t.Fatal(err)
	}
	if deleted.DataProtectionStatus != "DELETED" {
		t.Fatalf("after delete: status %q, want DELETED", deleted.DataProtectionStatus)
	}
	// The delete of a policy the group never had is a not-found, so it
	// cannot mint the DELETED marker on an untouched group.
	if err := svc.deleteDataProtectionPolicyCore("matrix-grc-group", "us-east-1"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("second delete: code=%q err=%v", logsErrorCode(err), err)
	}

	if err := svc.putBearerTokenAuthenticationCore(PutBearerTokenAuthenticationInput{
		LogGroupIdentifier: "matrix-grc-group", BearerTokenAuthenticationEnabled: true,
		BearerTokenAuthenticationSet: true, Region: "us-east-1",
	}); err != nil {
		t.Fatalf("enable bearer token authentication: %v", err)
	}
	enabled, err := store.GetLogGroup("matrix-grc-group")
	if err != nil {
		t.Fatal(err)
	}
	if !enabled.BearerTokenAuthenticationEnabled {
		t.Fatal("after enable: the record still reports disabled")
	}

	// The listing core serves the same records, so the DescribeLogGroups
	// serialiser's source fields are these.
	listed, err := svc.describeLogGroupsCore(ListLogGroupsInput{LogGroupNamePrefix: "matrix-grc-group", Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(listed.LogGroups) != 1 {
		t.Fatalf("describe returned %d groups, want 1", len(listed.LogGroups))
	}
	if listed.LogGroups[0].DataProtectionStatus != "DELETED" || !listed.LogGroups[0].BearerTokenAuthenticationEnabled {
		t.Fatalf("described group carries status %q bearer=%v", listed.LogGroups[0].DataProtectionStatus, listed.LogGroups[0].BearerTokenAuthenticationEnabled)
	}
}

// --- AssociateKmsKey: query-result ARN form and member exclusivity ---

func TestAssociateKmsKeyQueryResultForm(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-kms-group")
	queryResultArn := "arn:aws:logs:us-east-1:000000000000:query-result:*"
	keyArn := "arn:aws:kms:us-east-1:000000000000:key/1234abcd-12ab-34cd-56ef-1234567890ab"

	err := svc.associateKmsKeyCore(AssociateKmsKeyInput{
		LogGroupName: "matrix-kms-group", ResourceIdentifier: queryResultArn, KmsKeyId: keyArn, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("both target members: code=%q err=%v", code, err)
	}
	if err := svc.associateKmsKeyCore(AssociateKmsKeyInput{
		ResourceIdentifier: queryResultArn, KmsKeyId: keyArn, Region: "us-east-1"}); err != nil {
		t.Fatalf("query-result association: %v", err)
	}
	rec, err := store.GetQueryResultKmsKey()
	if err != nil || rec.KmsKeyId != keyArn {
		t.Fatalf("recorded association: rec=%+v err=%v", rec, err)
	}
	if err := svc.disassociateKmsKeyCore(DisassociateKmsKeyInput{
		ResourceIdentifier: queryResultArn, Region: "us-east-1"}); err != nil {
		t.Fatalf("query-result disassociation: %v", err)
	}
	if _, err := store.GetQueryResultKmsKey(); err == nil {
		t.Fatalf("record should be gone after disassociation")
	}
	// A resource that is neither form reports as not found.
	err = svc.associateKmsKeyCore(AssociateKmsKeyInput{
		ResourceIdentifier: "arn:aws:logs:us-east-1:000000000000:query-result:other", KmsKeyId: keyArn, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "ResourceNotFoundException" {
		t.Fatalf("foreign resource form: code=%q err=%v", code, err)
	}
}

// --- PutDestination: targetArn and roleArn are required members ---

func TestPutDestinationRequiredMembers(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-dest-group")

	_, err := svc.putDestinationCore("matrix-dest", "arn:aws:iam::000000000000:role/LogsRole", "", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("missing targetArn: code=%q err=%v", code, err)
	}
	_, err = svc.putDestinationCore("matrix-dest", "", "arn:aws:kinesis:us-east-1:000000000000:stream/Target", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("missing roleArn: code=%q err=%v", code, err)
	}
}

// --- Limit overflow identity per operation family: the delivery and
// scheduled-query families declare ValidationException; a Limit above
// the bound rejects with that code, not InvalidParameterException. ---

func TestLimitOverflowRejectsWithDeclaredFamilyIdentity(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-limit-group")

	_, _, err := svc.describeDeliverySourcesCore("us-east-1", "", logsstore.MaxDescribeLimit+1)
	if code := logsErrorCode(err); code != "ValidationException" {
		t.Fatalf("DescribeDeliverySources limit overflow: code=%q err=%v", code, err)
	}
	_, _, err = svc.listScheduledQueriesCore(store, "", "", logsstore.MaxListMaxResults+1, "")
	if code := logsErrorCode(err); code != "ValidationException" {
		t.Fatalf("ListScheduledQueries maxResults overflow: code=%q err=%v", code, err)
	}
	_, err = svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier: "arn:aws:logs:us-east-1:000000000000:scheduled-query:limit",
		StartTime:  1, EndTime: 2,
		MaxResults: logsstore.MaxListMaxResults + 1,
	})
	if code := logsErrorCode(err); code != "ValidationException" {
		t.Fatalf("GetScheduledQueryHistory maxResults overflow: code=%q err=%v", code, err)
	}

	// The legacy families keep their declared InvalidParameterException.
	_, _, err = svc.describeSubscriptionFiltersCore(store, &DescribeSubscriptionFiltersInput{
		LogGroupName: "matrix-limit-group", Limit: logsstore.MaxDescribeLimit + 1})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("DescribeSubscriptionFilters limit overflow: code=%q err=%v", code, err)
	}
}

// --- Filter pattern syntax and the per-group regex quota ---

// TestPutFilterPatternSyntaxAndRegexQuota pins the pattern grammar gate
// (the documented dialect and shape rejections surface as
// InvalidParameterException instead of storing match-nothing filters)
// and the per-group census: "There is a maximum of 5 filter patterns
// containing regex for each log group when creating metric filters or
// subscription filters" — counted across both families, with a replaced
// filter keeping its slot and non-regex patterns not counted.
func TestPutFilterPatternSyntaxAndRegexQuota(t *testing.T) {
	svc, store := newReadTestService(t, "pattern-quota-group")

	mk := func(name, pattern string) error {
		return svc.putMetricFilterCore("pattern-quota-group", name, pattern, true,
			[]logsstore.MetricTransformation{{
				MetricName: "Count", MetricNamespace: "App", MetricValue: "1",
				MetricNameSet: true, MetricNamespaceSet: true, MetricValueSet: true,
			}},
			false, "", nil, "us-east-1")
	}

	for _, tc := range []struct{ name, pattern string }{
		{"unsupported symbol", `%something!%`},
		{"parentheses subpattern", `%(alpha|beta)%`},
		{"multi-byte regex", `%マルチバイト%`},
		{"third regex in JSON pattern", `{ $.a = %r1% && $.b = %r2% && $.c = %r3% }`},
		{"unterminated delimited form", `[w1=%time%`},
		{"unbalanced quotes", `{ $.a = "unbalanced }`},
	} {
		if err := mk(tc.name, tc.pattern); logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("%s: code=%q err=%v", tc.name, logsErrorCode(err), err)
		}
	}

	for i := 1; i <= 4; i++ {
		if err := mk(fmt.Sprintf("rx-%d", i), `%[a-z]+%`); err != nil {
			t.Fatalf("regex-bearing metric filter %d: %v", i, err)
		}
	}
	if err := svc.putSubscriptionFilterCore(context.Background(), store, &PutSubscriptionFilterInput{
		LogGroupName:     "pattern-quota-group",
		FilterName:       "rx-sub",
		FilterPattern:    `%ERROR%`,
		FilterPatternSet: true,
		DestinationArn:   "arn:aws:kinesis:us-east-1:000000000000:stream/logs-target",
		Distribution:     "ByLogStream",
	}); err != nil {
		t.Fatalf("regex-bearing subscription filter (fifth slot): %v", err)
	}
	if err := mk("rx-sixth", `%[0-9]+%`); logsErrorCode(err) != "LimitExceededException" {
		t.Fatalf("sixth regex-bearing pattern: code=%q err=%v", logsErrorCode(err), err)
	}
	if err := mk("plain", "ERROR"); err != nil {
		t.Fatalf("non-regex pattern is not counted against the quota: %v", err)
	}
	if err := mk("rx-1", `%[a-z]+%`); err != nil {
		t.Fatalf("replacing a regex-bearing filter keeps its slot: %v", err)
	}
}

// --- Scheduled queries: window requiredness, not-found, member traits,
// enum filters ---

func TestScheduledQueryWindowRequirednessAndNotFound(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-sq-window-group")

	// Both window members are required: an omitting request rejects, an
	// explicit zero window is the documented minimum and reaches the
	// identifier resolution, where an unknown query rejects with the
	// operation's declared ResourceNotFoundException.
	if _, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier: "no-such-scheduled-query"}); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("omitted window: %v", err)
	}
	if _, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier: "no-such-scheduled-query", StartTimeSet: true, EndTimeSet: true}); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("unknown identifier: %v", err)
	}
	if _, err := svc.getScheduledQueryHistoryCore(store, &GetScheduledQueryHistoryInput{
		Identifier: "no-such-scheduled-query", StartTime: -1, EndTime: 2, StartTimeSet: true, EndTimeSet: true}); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("negative window: %v", err)
	}

	// The list filters are enumerations; values outside the vocabulary
	// reject instead of silently matching nothing.
	if _, _, err := svc.listScheduledQueriesCore(store, "enabled", "", 0, ""); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("lowercase state filter: %v", err)
	}
	if _, _, err := svc.listScheduledQueriesCore(store, "", "BOGUS", 0, ""); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("garbage scheduleType filter: %v", err)
	}
}

func TestCreateScheduledQueryMemberTraits(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-sq-traits-group")
	base := func() *CreateScheduledQueryInput {
		return &CreateScheduledQueryInput{
			Name:               "sq-traits",
			QueryString:        "fields @message",
			QueryLanguage:      "CWLI",
			ScheduleExpression: "rate(1 hour)",
			ExecutionRoleArn:   "arn:aws:iam::000000000000:role/deliver",
		}
	}

	in := base()
	in.Description = strings.Repeat("d", logsstore.MaxScheduledQueryDescriptionLength+1)
	if _, err := svc.createScheduledQueryCore(store, in); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("oversized description: %v", err)
	}

	in = base()
	in.LogGroupIdentifiers = []string{strings.Repeat("a", logsstore.MaxLogGroupIdentifierLength+1)}
	if _, err := svc.createScheduledQueryCore(store, in); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("oversized identifier element: %v", err)
	}
	in = base()
	in.LogGroupIdentifiers = []string{"bad identifier!"}
	if _, err := svc.createScheduledQueryCore(store, in); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("identifier element alphabet: %v", err)
	}

	in = base()
	in.Tags = make(map[string]string)
	for i := 0; i < logsstore.MaxScheduledQueryTags+1; i++ {
		in.Tags[fmt.Sprintf("k%d", i)] = "v"
	}
	if _, err := svc.createScheduledQueryCore(store, in); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("oversized tag set: %v", err)
	}
	in = base()
	in.Tags = map[string]string{"bad key!": "v"}
	if _, err := svc.createScheduledQueryCore(store, in); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("tag key alphabet: %v", err)
	}

	in = base()
	in.ScheduleStartTime = -1
	if _, err := svc.createScheduledQueryCore(store, in); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("negative scheduleStartTime: %v", err)
	}
}

// DescribeLogGroups' list members carry their documented bounds: "You
// can specify as many as 50 log groups in the array" and "You can
// specify as many as 20 account IDs in the array" — an oversized list
// rejects instead of being silently processed.
func TestDescribeLogGroupsListBounds(t *testing.T) {
	svc, _ := newReadTestService(t, "dlg-bounds-group")

	groups := make([]string, 51)
	for i := range groups {
		groups[i] = fmt.Sprintf("dlg-bounds-group-%d", i)
	}
	if _, err := svc.describeLogGroupsCore(ListLogGroupsInput{
		LogGroupIdentifiers: groups, Region: "us-east-1",
	}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("fifty-one identifiers = %v, want InvalidParameterException", err)
	}

	accounts := make([]string, 21)
	for i := range accounts {
		accounts[i] = fmt.Sprintf("11111111111%d", i)
	}
	if _, err := svc.describeLogGroupsCore(ListLogGroupsInput{
		AccountIdentifiers: accounts, IncludeLinkedAccounts: true, Region: "us-east-1",
	}); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("twenty-one account IDs = %v, want InvalidParameterException", err)
	}
}

// --- Tag trio: destinations are taggable resources, and the 50-tag
// ceiling binds the cumulative total ---

func untagResourceRequest(arn string, keys ...string) *request.ParsedRequest {
	return &request.ParsedRequest{
		Operation:  "UntagResource",
		Parameters: map[string]interface{}{"resourceArn": arn, "tagKeys": toIface(keys)},
	}
}

func listTagsResourceRequest(arn string) *request.ParsedRequest {
	return &request.ParsedRequest{
		Operation:  "ListTagsForResource",
		Parameters: map[string]interface{}{"resourceArn": arn},
	}
}

func toIface(keys []string) []interface{} {
	out := make([]interface{}, len(keys))
	for i, k := range keys {
		out[i] = k
	}
	return out
}

func TestTagTrioDestinationsAndCumulativeCeiling(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-dest-group")
	ctx := context.Background()
	cfg := svc.tagHandlerConfig(store)

	// "Currently, the only CloudWatch Logs resources that can be tagged
	// are log groups and destinations" — the trio resolves a destination
	// ARN, and the tags PutDestination carried join the same tag store.
	dest, err := svc.putDestinationCore("matrix-dest",
		"arn:aws:iam::000000000000:role/deliverer",
		"arn:aws:kinesis:us-east-1:000000000000:stream/matrix",
		map[string]string{"origin": "put"}, "us-east-1")
	if err != nil {
		t.Fatalf("put destination: %v", err)
	}
	if _, err := tagutil.HandleTag(ctx, tagResourceRequest(dest.ARN, map[string]interface{}{"trio": "v"}), cfg); err != nil {
		t.Fatalf("tag destination: %v", err)
	}
	listed, err := tagutil.HandleList(ctx, listTagsResourceRequest(dest.ARN), cfg)
	if err != nil {
		t.Fatalf("list destination tags: %v", err)
	}
	got := listed.(map[string]interface{})["tags"].(map[string]string)
	if got["origin"] != "put" || got["trio"] != "v" {
		t.Fatalf("destination tags = %v, want the PutDestination pair plus the trio write", got)
	}
	if _, err := tagutil.HandleUntag(ctx, untagResourceRequest(dest.ARN, "origin"), cfg); err != nil {
		t.Fatalf("untag destination: %v", err)
	}
	listed, err = tagutil.HandleList(ctx, listTagsResourceRequest(dest.ARN), cfg)
	if err != nil {
		t.Fatal(err)
	}
	got = listed.(map[string]interface{})["tags"].(map[string]string)
	if _, still := got["origin"]; still {
		t.Fatalf("untag left the key behind: %v", got)
	}

	missing := strings.Replace(dest.ARN, "matrix-dest", "matrix-missing", 1)
	_, err = tagutil.HandleTag(ctx, tagResourceRequest(missing, map[string]interface{}{"k": "v"}), cfg)
	if code := logsErrorCode(err); code != "ResourceNotFoundException" {
		t.Fatalf("missing destination: code=%q err=%v", code, err)
	}

	// "You can associate as many as 50 tags with a CloudWatch Logs
	// resource": a group holding 30 tags rejects the write that would
	// carry it past 50 — TooManyTagsException on the tag trio, the legacy
	// operation's parameter identity on TagLogGroup.
	lg, err := store.GetLogGroup("matrix-dest-group")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-dest-group", Tags: testTags(30), Region: "us-east-1"}); err != nil {
		t.Fatalf("seed 30 tags: %v", err)
	}
	disjoint := make(map[string]interface{}, 21)
	for i := 0; i < 21; i++ {
		disjoint[fmt.Sprintf("trio-%d", i)] = "v"
	}
	_, err = tagutil.HandleTag(ctx, tagResourceRequest(lg.ARN, disjoint), cfg)
	if code := logsErrorCode(err); code != "TooManyTagsException" {
		t.Fatalf("cumulative 51 via trio: code=%q err=%v", code, err)
	}
	legacy := func(n int) map[string]string {
		tags := make(map[string]string, n)
		for i := 0; i < n; i++ {
			tags[fmt.Sprintf("legacy-%d", i)] = "v"
		}
		return tags
	}
	err = svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-dest-group", Tags: legacy(21), Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("cumulative 51 via legacy op: code=%q err=%v", code, err)
	}
	// The merge up to exactly 50 still succeeds.
	if err := svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-dest-group", Tags: legacy(20), Region: "us-east-1"}); err != nil {
		t.Fatalf("fill to the 50 ceiling: %v", err)
	}
}

// The reserved aws: prefix rejects on every tag write surface.
func TestTagReservedPrefixRejects(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-res-group")
	lg, err := store.GetLogGroup("matrix-res-group")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	cfg := svc.tagHandlerConfig(store)

	_, err = tagutil.HandleTag(ctx, tagResourceRequest(lg.ARN, map[string]interface{}{"aws:reserved": "v"}), cfg)
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("trio aws: prefix: code=%q err=%v", code, err)
	}
	err = svc.tagLogGroupCore(TagLogGroupInput{LogGroupName: "matrix-res-group", Tags: map[string]string{"aws:reserved": "v"}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("legacy aws: prefix: code=%q err=%v", code, err)
	}
	_, err = svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "matrix-res-new", Tags: map[string]string{"aws:reserved": "v"}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("createLogGroup aws: prefix: code=%q err=%v", code, err)
	}
}

// CreateLogGroup's reserved prefix and KMS class checks: "Log group names
// can't start with the string aws/" rejects the plain prefix while the
// platform's leading-slash /aws/ form stays legal; a KMS key that is
// missing, disabled or asymmetric rejects with InvalidParameterException on
// both the create and associate paths ("If you attempt to associate a KMS
// key with the log group but the KMS key does not exist or the KMS key is
// disabled, you receive an InvalidParameterException error"; "CloudWatch
// Logs supports only symmetric KMS keys"), while a usable symmetric key
// passes. The per-region group counter counts the groups the store holds.
func TestCreateLogGroupReservedPrefixKmsClassAndQuotaCount(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-create-group")

	_, err := svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "aws/lambda/hidden", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("aws/ prefix: code=%q err=%v", code, err)
	}
	if _, err := svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "/aws/lambda/allowed", Region: "us-east-1"}); err != nil {
		t.Fatalf("leading-slash /aws/ form rejected: %v", err)
	}

	const keyARN = "arn:aws:kms:us-east-1:000000000000:key/11111111-1111-1111-1111-111111111111"
	kmsFake := newFakeKMSInvoker(keyARN)
	svc.kms = kmsFake

	// A key the platform does not know at all.
	_, err = svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "matrix-kms-create", KmsKeyId: keyARN + "0", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("create with unknown key: code=%q err=%v", code, err)
	}

	// A key that exists but fails the enabled-and-symmetric check.
	kmsFake.markUnusable(keyARN)
	_, err = svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "matrix-kms-create", KmsKeyId: keyARN, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("create with unusable key: code=%q err=%v", code, err)
	}
	err = svc.associateKmsKeyCore(AssociateKmsKeyInput{LogGroupName: "matrix-create-group", KmsKeyId: keyARN, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("associate unusable key: code=%q err=%v", code, err)
	}

	// A usable symmetric key passes both paths.
	delete(kmsFake.unusable, keyARN)
	if _, err := svc.createLogGroupCore(CreateLogGroupInput{LogGroupName: "matrix-kms-create", KmsKeyId: keyARN, Region: "us-east-1"}); err != nil {
		t.Fatalf("create with usable key: %v", err)
	}
	if err := svc.associateKmsKeyCore(AssociateKmsKeyInput{LogGroupName: "matrix-create-group", KmsKeyId: keyARN, Region: "us-east-1"}); err != nil {
		t.Fatalf("associate usable key: %v", err)
	}

	count, err := store.CountLogGroups()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("CountLogGroups = %d, want the 3 groups created (fixture, /aws/ form, KMS-carrying)", count)
	}
}

// The event read/write plane carries the modelled input traits: name
// patterns reject as InvalidParameterException rather than surfacing as a
// wrong ResourceNotFound identity, the stream array carries its
// InputLogStreamNames ceiling ("Maximum number of 100 items"), the window
// members carry the Timestamp minimum of 0, the stream-prefix members
// target LogStreamName's pattern, and a DescribeLogStreams listing against
// a missing group fails with the operation's declared
// ResourceNotFoundException instead of serving an empty page.
func TestEventPathTraitsAndMissingGroupListing(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-events-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("matrix-events-stream", "matrix-events-group")); err != nil {
		t.Fatal(err)
	}
	oneEvent := []PutLogEvent{{LogEntry: logsstore.LogEntry{Timestamp: 1, Message: "m"}, TimestampSet: true}}

	// Malformed names are parameter errors on both event planes, not a
	// ResourceNotFound from the store walk.
	_, err := svc.putLogEventsCore(PutLogEventsInput{LogGroupName: "matrix-events-group", LogStreamName: "bad:stream", Events: oneEvent, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("putLogEvents stream pattern: code=%q err=%v", code, err)
	}
	_, err = svc.getLogEventsCore(GetLogEventsInput{LogGroupName: "matrix-events-group", LogStreamName: "bad:stream", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("getLogEvents stream pattern: code=%q err=%v", code, err)
	}
	_, err = svc.filterLogEventsCore(FilterLogEventsInput{LogGroupName: "matrix:group", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("filterLogEvents group pattern: code=%q err=%v", code, err)
	}

	// The stream array's InputLogStreamNames ceiling.
	names := make([]string, logsstore.MaxInputLogStreamNames+1)
	for i := range names {
		names[i] = "stream"
	}
	_, err = svc.filterLogEventsCore(FilterLogEventsInput{LogGroupName: "matrix-events-group", LogStreamNames: names, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("logStreamNames over ceiling: code=%q err=%v", code, err)
	}

	// The window members' Timestamp minimum of 0.
	_, err = svc.getLogEventsCore(GetLogEventsInput{LogGroupName: "matrix-events-group", LogStreamName: "matrix-events-stream", StartTime: -1, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("getLogEvents negative startTime: code=%q err=%v", code, err)
	}
	_, err = svc.filterLogEventsCore(FilterLogEventsInput{LogGroupName: "matrix-events-group", EndTime: -1, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("filterLogEvents negative endTime: code=%q err=%v", code, err)
	}

	// The stream-prefix members target LogStreamName's pattern.
	_, err = svc.describeLogStreamsCore(DescribeLogStreamsInput{LogGroupName: "matrix-events-group", LogStreamNamePrefix: "has:colon", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("describeLogStreams prefix pattern: code=%q err=%v", code, err)
	}
	_, err = svc.filterLogEventsCore(FilterLogEventsInput{LogGroupName: "matrix-events-group", LogStreamNamePref: "has:colon", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("filterLogEvents prefix pattern: code=%q err=%v", code, err)
	}

	// A listing against a missing group fails with the operation's
	// declared identity instead of an empty page.
	_, err = svc.describeLogStreamsCore(DescribeLogStreamsInput{LogGroupName: "matrix-missing-group", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "ResourceNotFoundException" {
		t.Fatalf("describeLogStreams missing group: code=%q err=%v", code, err)
	}
}

// ListLogGroups carries its documented filter members: logGroupClass
// filters and reaches the summary entries; logGroupTags combines with AND
// across filters, any-of within a key's values, a values-less filter on
// key existence, case-insensitive wildcard and case-sensitive negation
// matching; fieldIndexNames selects groups whose field-index set (the
// DEFAULT category plus the effective policy's fields) contains every
// requested name; dataSources selects through the delivery-source
// association with OR within a dimension and AND across dimensions; and
// the members' bounds and element traits reject.
func TestListLogGroupsFilterMembersAndBounds(t *testing.T) {
	svc, store := newReadTestService(t, "lfg-anchor")
	if err := store.DeleteLogGroup("lfg-anchor"); err != nil {
		t.Fatal(err)
	}
	mustGroup := func(name, class string, tags map[string]string) {
		lg := logsstore.NewLogGroup(name, "us-east-1", "000000000000")
		lg.LogGroupClass = class
		if err := store.CreateLogGroup(lg); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
		if len(tags) > 0 {
			if err := store.Tags().Tag(lg.ARN, tags); err != nil {
				t.Fatalf("tag %s: %v", name, err)
			}
		}
	}
	mustGroup("lfg-prod-core", "STANDARD", map[string]string{"env": "prod", "team": "core"})
	mustGroup("lfg-dev", "STANDARD", map[string]string{"env": "dev"})
	mustGroup("lfg-prod-ia", "INFREQUENT_ACCESS", map[string]string{"env": "prod"})

	prodCore, err := store.GetLogGroup("lfg-prod-core")
	if err != nil {
		t.Fatal(err)
	}
	// A custom field index on one group; the other groups ride the
	// DEFAULT category alone.
	if err := store.ReplaceIndexPolicy("lfg-prod-core", func(current *logsstore.IndexPolicy, trail *logsstore.IndexInactiveTrail) (*logsstore.IndexPolicy, *logsstore.IndexInactiveTrail, error) {
		return &logsstore.IndexPolicy{
			LogGroupName:   "lfg-prod-core",
			PolicyDocument: `{"Fields":["TransactionId"]}`,
		}, trail, nil
	}); err != nil {
		t.Fatal(err)
	}
	// The data-source association: one delivery source addresses the
	// tagged STANDARD group alone.
	if err := store.PutDeliverySource(&logsstore.DeliverySource{
		Name:        "lfg-source",
		ResourceArn: prodCore.ARN,
		Service:     "logs",
		LogType:     "AWS::Logs::LogGroup",
	}); err != nil {
		t.Fatal(err)
	}

	names := func(in ListLogGroupsInput) []string {
		in.Region = "us-east-1"
		res, err := svc.listLogGroupsCore(in)
		if err != nil {
			t.Fatalf("listing: %v", err)
		}
		got := make([]string, 0, len(res.LogGroups))
		for _, lg := range res.LogGroups {
			got = append(got, lg.Name)
		}
		return got
	}

	if got := names(ListLogGroupsInput{LogGroupClass: "INFREQUENT_ACCESS"}); len(got) != 1 || got[0] != "lfg-prod-ia" {
		t.Fatalf("class filter = %v", got)
	}
	if got := names(ListLogGroupsInput{LogGroupClass: "STANDARD"}); len(got) != 2 {
		t.Fatalf("class filter STANDARD = %v", got)
	}

	if got := names(ListLogGroupsInput{LogGroupTags: []TagFilterInput{{Key: "env", Values: []string{"prod"}}}}); len(got) != 2 {
		t.Fatalf("tag any-of = %v", got)
	}
	if got := names(ListLogGroupsInput{LogGroupTags: []TagFilterInput{
		{Key: "env", Values: []string{"prod"}}, {Key: "team"},
	}}); len(got) != 1 || got[0] != "lfg-prod-core" {
		t.Fatalf("tag AND with values-less key = %v", got)
	}
	// Wildcard matching is case-insensitive; negation is case-sensitive.
	if got := names(ListLogGroupsInput{LogGroupTags: []TagFilterInput{{Key: "env", Values: []string{"PR*"}}}}); len(got) != 2 {
		t.Fatalf("tag wildcard = %v", got)
	}
	if got := names(ListLogGroupsInput{LogGroupTags: []TagFilterInput{{Key: "env", Values: []string{"!prod"}}}}); len(got) != 1 || got[0] != "lfg-dev" {
		t.Fatalf("tag negation = %v", got)
	}

	if got := names(ListLogGroupsInput{FieldIndexNames: []string{"traceId"}}); len(got) != 3 {
		t.Fatalf("default field index = %v", got)
	}
	if got := names(ListLogGroupsInput{FieldIndexNames: []string{"TransactionId"}}); len(got) != 1 || got[0] != "lfg-prod-core" {
		t.Fatalf("custom field index = %v", got)
	}
	if got := names(ListLogGroupsInput{FieldIndexNames: []string{"TransactionId", "traceId"}}); len(got) != 1 || got[0] != "lfg-prod-core" {
		t.Fatalf("all-specified field indexes = %v", got)
	}
	if got := names(ListLogGroupsInput{FieldIndexNames: []string{"NoSuchField"}}); len(got) != 0 {
		t.Fatalf("unknown field index = %v", got)
	}
	// The member's element alphabet carries no @, so the @-prefixed
	// generated-field indexes are not addressable through it.
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{FieldIndexNames: []string{"@logStream"}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("@-named field index: code=%q err=%v", code, err)
	}

	if got := names(ListLogGroupsInput{DataSources: []DataSourceFilterInput{{Name: "lfg-source"}}}); len(got) != 1 || got[0] != "lfg-prod-core" {
		t.Fatalf("data source by name = %v", got)
	}
	if got := names(ListLogGroupsInput{DataSources: []DataSourceFilterInput{{Name: "lfg-source", Type: "AWS::Logs::LogGroup"}}}); len(got) != 1 {
		t.Fatalf("data source name and type = %v", got)
	}
	if got := names(ListLogGroupsInput{DataSources: []DataSourceFilterInput{{Name: "lfg-source", Type: "Other"}}}); len(got) != 0 {
		t.Fatalf("data source AND across dimensions = %v", got)
	}
	if got := names(ListLogGroupsInput{DataSources: []DataSourceFilterInput{{Name: "missing-source"}}}); len(got) != 0 {
		t.Fatalf("unknown data source = %v", got)
	}

	// The members' bounds and element traits.
	six := func(n int) []TagFilterInput {
		filters := make([]TagFilterInput, n)
		for i := range filters {
			filters[i] = TagFilterInput{Key: fmt.Sprintf("k%d", i)}
		}
		return filters
	}
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{LogGroupTags: six(logsstore.MaxListLogGroupTagFilters + 1), Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("tag filter bound: code=%q err=%v", code, err)
	}
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{LogGroupTags: []TagFilterInput{{Key: "env", Values: []string{"!bad*value*"}}}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("tag value trait: code=%q err=%v", code, err)
	}
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{DataSources: make([]DataSourceFilterInput, logsstore.MaxListDataSourceFilters+1), Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("data source bound: code=%q err=%v", code, err)
	}
	fieldNames := make([]string, logsstore.MaxListFieldIndexNames+1)
	for i := range fieldNames {
		fieldNames[i] = "f"
	}
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{FieldIndexNames: fieldNames, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("field index bound: code=%q err=%v", code, err)
	}
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{AccountIdentifiers: []string{"1234"}, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("account identifier trait: code=%q err=%v", code, err)
	}
	_, err = svc.listLogGroupsCore(ListLogGroupsInput{LogGroupClass: "BOGUS", Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("class enum: code=%q err=%v", code, err)
	}
}

// DescribeMetricFilters' metric members filter and require each other:
// "Filters results to include only those with the specified metric name.
// If you include this parameter in your request, you must also include
// the metricNamespace parameter" (and vice versa) — one without the other
// rejects, and the pair selects the filters whose transformations publish
// that metric, with the page token scoped to the metric members.
func TestDescribeMetricFiltersMetricMembers(t *testing.T) {
	svc, store := newReadTestService(t, "matrix-desc-metric")
	put := func(filter, metricName, namespace string) {
		if err := store.PutMetricFilter(&logsstore.MetricFilter{
			Name: filter, LogGroupName: "matrix-desc-metric", FilterPattern: "ERROR",
			MetricTransformations: []logsstore.MetricTransformation{{
				MetricName: metricName, MetricNamespace: namespace, MetricValue: "1",
			}},
		}); err != nil {
			t.Fatal(err)
		}
	}
	put("errors", "ErrorCount", "App")
	put("errors-too", "ErrorCount", "App")
	put("latency", "Latency", "App")
	put("other-ns", "ErrorCount", "Other")

	_, _, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "matrix-desc-metric", MetricNameSet: true,
		MetricName: "ErrorCount", Limit: 50, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("metricName without namespace: code=%q err=%v", code, err)
	}
	_, _, err = svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "matrix-desc-metric", MetricNamespaceSet: true,
		MetricNamespace: "App", Limit: 50, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("metricNamespace without name: code=%q err=%v", code, err)
	}

	res, token, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "matrix-desc-metric", MetricNameSet: true, MetricNamespaceSet: true,
		MetricName: "ErrorCount", MetricNamespace: "App", Limit: 50, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 || res[0].Name != "errors" || res[1].Name != "errors-too" || token != "" {
		t.Fatalf("metric filter = %v, token %q, want the two ErrorCount/App filters", filterNames(res), token)
	}

	// The page token is scoped to the metric members: a token minted under
	// one metric pair rejects against another.
	res1, pairToken, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "matrix-desc-metric", MetricNameSet: true, MetricNamespaceSet: true,
		MetricName: "ErrorCount", MetricNamespace: "App", Limit: 1, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res1) != 1 || pairToken == "" {
		t.Fatalf("limited metric page = %v, token %q", filterNames(res1), pairToken)
	}
	_, _, err = svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "matrix-desc-metric", MetricNameSet: true, MetricNamespaceSet: true,
		MetricName: "Latency", MetricNamespace: "App", NextToken: pairToken, Limit: 50, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("token replay across metric pairs: code=%q err=%v", code, err)
	}
}

func filterNames(filters []*logsstore.MetricFilter) []string {
	names := make([]string, 0, len(filters))
	for _, f := range filters {
		names = append(names, f.Name)
	}
	return names
}

// The filter Describe operations declare ResourceNotFoundException: a
// listing against a group that does not exist fails instead of serving
// an empty page (the sibling stream listing carries the same check).
func TestDescribeFilterFamiliesRejectMissingGroup(t *testing.T) {
	svc, store := newReadTestService(t, "describe-filters-present")

	_, _, err := svc.describeSubscriptionFiltersCore(store, &DescribeSubscriptionFiltersInput{
		LogGroupName: "no-such-filter-group", Limit: 10})
	if code := logsErrorCode(err); code != "ResourceNotFoundException" {
		t.Fatalf("DescribeSubscriptionFilters missing group: code=%q err=%v", code, err)
	}
	_, _, err = svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "no-such-filter-group", Limit: 10, Region: "us-east-1"})
	if code := logsErrorCode(err); code != "ResourceNotFoundException" {
		t.Fatalf("DescribeMetricFilters missing group: code=%q err=%v", code, err)
	}

	// A present group with no filters still lists: an empty filter set
	// is a valid page, distinct from a missing group.
	filters, _, err := svc.describeSubscriptionFiltersCore(store, &DescribeSubscriptionFiltersInput{
		LogGroupName: "describe-filters-present", Limit: 10})
	if err != nil || len(filters) != 0 {
		t.Fatalf("present group listing: filters=%v err=%v", filters, err)
	}
	metricFilters, _, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "describe-filters-present", Limit: 10, Region: "us-east-1"})
	if err != nil || len(metricFilters) != 0 {
		t.Fatalf("present group metric listing: filters=%v err=%v", metricFilters, err)
	}
}

// The transformation's dimensions and unit members carry their documented
// contract: "One metric filter can include as many as three dimensions"
// with emitSystemFieldDimensions counting toward the total, the conflict
// rule "If you assign dimensions to a metric created by a metric filter,
// you can't assign a default value for that metric", the unit enum ("The
// unit to assign to the metric. If you omit this, the unit is set as
// None"), the system-field value set ("Valid values are @aws.account and
// @aws.region."), and the per-key 255-character ceiling.
func TestPutMetricFilterDimensionsAndUnit(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-dims-group")
	base := func(muts ...func(*logsstore.MetricTransformation)) []logsstore.MetricTransformation {
		t := logsstore.MetricTransformation{
			MetricName: "Count", MetricNamespace: "App", MetricValue: "1",
			MetricNameSet: true, MetricNamespaceSet: true, MetricValueSet: true,
		}
		for _, mut := range muts {
			mut(&t)
		}
		return []logsstore.MetricTransformation{t}
	}
	put := func(name string, ts []logsstore.MetricTransformation, emit ...string) error {
		return svc.putMetricFilterCore("matrix-dims-group", name, "ERROR", true, ts, false, "", emit, "us-east-1")
	}

	if err := put("plain", base()); err != nil {
		t.Fatalf("undimensioned filter: %v", err)
	}
	if err := put("with-dims", base(func(t *logsstore.MetricTransformation) {
		t.Dimensions = map[string]string{"Request": "$request"}
		t.Unit = "Count"
	})); err != nil {
		t.Fatalf("dimensioned filter with unit: %v", err)
	}
	// Three transformation dimensions alone, and fewer with system
	// fields, stay within the ceiling.
	if err := put("three-dims", base(func(t *logsstore.MetricTransformation) {
		t.Dimensions = map[string]string{"a": "1", "b": "2", "c": "3"}
	})); err != nil {
		t.Fatalf("three dimensions: %v", err)
	}
	if err := put("system-fields", base(), "@aws.account", "@aws.region"); err != nil {
		t.Fatalf("system fields alone: %v", err)
	}

	err := put("over-ceiling", base(func(t *logsstore.MetricTransformation) {
		t.Dimensions = map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"}
	}))
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("four dimensions: code=%q err=%v", code, err)
	}
	err = put("system-over-ceiling", base(func(t *logsstore.MetricTransformation) {
		t.Dimensions = map[string]string{"a": "1", "b": "2"}
	}), "@aws.account", "@aws.region")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("dimensions plus system fields over ceiling: code=%q err=%v", code, err)
	}
	err = put("default-conflict", base(func(t *logsstore.MetricTransformation) {
		t.Dimensions = map[string]string{"a": "1"}
		t.DefaultValue, t.DefaultValueSet = 0, true
	}))
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("dimensions with default value: code=%q err=%v", code, err)
	}
	err = put("bad-unit", base(func(t *logsstore.MetricTransformation) { t.Unit = "Furlongs" }))
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("unknown unit: code=%q err=%v", code, err)
	}
	err = put("long-key", base(func(t *logsstore.MetricTransformation) {
		t.Dimensions = map[string]string{strings.Repeat("k", 256): "v"}
	}))
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("256-character dimension key: code=%q err=%v", code, err)
	}
	err = put("bad-system-field", base(), "@aws.availabilityZone")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("unknown system field: code=%q err=%v", code, err)
	}

	// The stored transformation round-trips its dimensions and unit
	// through the listing the DescribeMetricFilters core serves.
	res, _, err := svc.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName: "matrix-dims-group", Limit: 50, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res {
		if f.Name != "with-dims" {
			continue
		}
		tt := f.MetricTransformations[0]
		if len(tt.Dimensions) != 1 || tt.Dimensions["Request"] != "$request" || tt.Unit != "Count" {
			t.Fatalf("round-tripped transformation = %+v", tt)
		}
	}
}

// The account-policy plane enforces its documented requiredness, enum and
// selectionCriteria forms: policyType is required and enum-bound on
// DescribeAccountPolicies ("Required: Yes"; an invalid value is a
// parameter error, not an empty listing) and DeleteAccountPolicy (whose
// missing-member identities are per-member); "Specifying
// selectionCriteria is valid only when you specify
// SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY or TRANSFORMER_POLICY";
// and "If policyType is SUBSCRIPTION_FILTER_POLICY, the only supported
// selectionCriteria filter is LogGroupName NOT IN []". accountIdentifiers
// carries its one-ID documented form, and a source account other than the
// local one owns no policies on this single-account platform.
func TestAccountPolicyRequirednessAndSelectionForms(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-ap-group")
	doc := `{"Fields":["SessionId"]}`

	if _, _, err := svc.describeAccountPoliciesCore("", "", nil, "", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("describe without policyType: %v", err)
	}
	if _, _, err := svc.describeAccountPoliciesCore("NOT_A_POLICY_TYPE", "", nil, "", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("describe with a foreign policyType: %v", err)
	}
	if err := svc.deleteAccountPolicyCore("", "FIELD_INDEX_POLICY", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("delete without policyName: %v", err)
	}
	if err := svc.deleteAccountPolicyCore("p", "", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("delete without policyType: %v", err)
	}
	if err := svc.deleteAccountPolicyCore("p", "NOT_A_POLICY_TYPE", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("delete with a foreign policyType: %v", err)
	}

	// selectionCriteria rides only the three selection-supporting types.
	_, err := svc.putAccountPolicyCore("dp-sc", doc, "DATA_PROTECTION_POLICY", "", "LogGroupNamePrefix = x-", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("data-protection policy with selectionCriteria: code=%q err=%v", code, err)
	}
	// The SUBSCRIPTION_FILTER_POLICY form: the empty selection, the
	// exclusion list and the empty list are legal; any other filter form
	// rejects.
	if _, err := svc.putAccountPolicyCore("sfp-empty", doc, "SUBSCRIPTION_FILTER_POLICY", "", "", "us-east-1"); err != nil {
		t.Fatalf("subscription policy without selectionCriteria: %v", err)
	}
	if _, err := svc.putAccountPolicyCore("sfp-list", doc, "SUBSCRIPTION_FILTER_POLICY", "", `LogGroupName NOT IN ["alpha","beta"]`, "us-east-1"); err != nil {
		t.Fatalf("subscription policy with an exclusion list: %v", err)
	}
	if _, err := svc.putAccountPolicyCore("sfp-none", doc, "SUBSCRIPTION_FILTER_POLICY", "", "LogGroupName NOT IN []", "us-east-1"); err != nil {
		t.Fatalf("subscription policy with an empty exclusion list: %v", err)
	}
	_, err = svc.putAccountPolicyCore("sfp-prefix", doc, "SUBSCRIPTION_FILTER_POLICY", "", "LogGroupNamePrefix = x-", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("subscription policy with a prefix filter: code=%q err=%v", code, err)
	}
	_, err = svc.putAccountPolicyCore("sfp-badname", doc, "SUBSCRIPTION_FILTER_POLICY", "", `LogGroupName NOT IN ["bad:name"]`, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("subscription policy with a malformed group name: code=%q err=%v", code, err)
	}

	// accountIdentifiers: one ID of twelve digits, and a foreign source
	// account owns nothing here.
	if _, _, err := svc.describeAccountPoliciesCore("FIELD_INDEX_POLICY", "", []string{"1234"}, "", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("short account identifier: %v", err)
	}
	if _, _, err := svc.describeAccountPoliciesCore("FIELD_INDEX_POLICY", "", []string{"000000000000", "999999999999"}, "", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("two account identifiers: %v", err)
	}
	foreign, _, err := svc.describeAccountPoliciesCore("SUBSCRIPTION_FILTER_POLICY", "", []string{"999999999999"}, "", "us-east-1")
	if err != nil {
		t.Fatalf("foreign account describe: %v", err)
	}
	if len(foreign) != 0 {
		t.Fatalf("foreign account owns no policies here, got %d", len(foreign))
	}
	local, _, err := svc.describeAccountPoliciesCore("SUBSCRIPTION_FILTER_POLICY", "", []string{"000000000000"}, "", "us-east-1")
	if err != nil {
		t.Fatalf("local account describe: %v", err)
	}
	if len(local) == 0 {
		t.Fatal("local account describe returned nothing")
	}
}

// The account-level data-protection policy carries the same document
// contract as its group-level twin: the two-block Audit/Deidentify
// structure and its own tighter character ceiling, both read from the
// PutAccountPolicy policyDocument member documentation.
func TestAccountDataProtectionPolicyStructure(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-ap-dpp-group")

	const valid = `{
		"Name": "DataProtectionPolicy",
		"Version": "2021-06-01",
		"Statement": [
			{"DataIdentifer": ["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],
			 "Operation": {"Audit": {"FindingsDestination": {"CloudWatchLogs": {"LogGroup": "audit-findings"}}}}},
			{"DataIdentifer": ["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],
			 "Operation": {"Deidentify": {"MaskConfig": {}}}}
		]
	}`
	if _, err := svc.putAccountPolicyCore("ap-dpp", valid, "DATA_PROTECTION_POLICY", "", "", "us-east-1"); err != nil {
		t.Fatalf("two-block account document: %v", err)
	}

	for _, tc := range []struct{ name, doc string }{
		{"one block", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}}]}`},
		{"non-empty mask config", `{"Statement":[{"DataIdentifer":["i"],"Operation":{"Audit":{"FindingsDestination":{}}}},{"DataIdentifer":["i"],"Operation":{"Deidentify":{"MaskConfig":{"Mode":"overwrite"}}}}]}`},
	} {
		_, err := svc.putAccountPolicyCore("ap-dpp", tc.doc, "DATA_PROTECTION_POLICY", "", "", "us-east-1")
		if code := logsErrorCode(err); code != "InvalidParameterException" {
			t.Fatalf("%s: code=%q err=%v", tc.name, code, err)
		}
	}

	// A document past the data-protection ceiling but inside the generic
	// policy-document cap rejects on the tighter bound alone.
	oversize := `{"Name":"` + strings.Repeat("n", logsstore.MaxDataProtectionPolicyDocumentBytes) + `"}`
	if runeCount := len([]rune(oversize)); runeCount <= int(logsstore.MaxDataProtectionPolicyDocumentBytes) || runeCount > int(logsstore.MaxPolicyDocumentLength) {
		t.Fatalf("fixture must sit in the window past the data-protection ceiling and inside the generic cap, got %d characters", runeCount)
	}
	_, err := svc.putAccountPolicyCore("ap-dpp", oversize, "DATA_PROTECTION_POLICY", "", "", "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("oversize account document: code=%q err=%v", code, err)
	}
}

// The stored-never-applied members now apply: fieldSelectionCriteria
// ("A filter expression that specifies which log events should be
// processed by this metric filter based on system fields such as source
// account and source region. Uses selection criteria syntax with
// operators like =, !=, AND, OR, IN, NOT IN") gates both filter families
// at evaluation time and validates its grammar; the subscription
// emitSystemFields member ("A list of system fields to include in the log
// events sent to the subscription destination. Valid values are
// @aws.account, @aws.region, and @source.log.") rides the delivered
// payload's events and validates its vocabulary. (The metric family's
// emitSystemFieldDimensions application and vocabulary are pinned by the
// dimensions-and-unit test.)
func TestFieldSelectionCriteriaAndEmitSystemFields(t *testing.T) {
	// The criteria evaluator: every documented operator form, the AND/OR
	// composition (an AND binding tighter than an enclosing OR), and the
	// absent-criteria pass-through.
	cases := []struct {
		criteria string
		expected bool
	}{
		{``, true},
		{`@aws.region = "us-east-1"`, true},
		{`@aws.region = "eu-west-1"`, false},
		{`@aws.region != "eu-west-1"`, true},
		{`@aws.account IN ["123456789012", "000000000000"]`, true},
		{`@aws.account IN ["999999999999"]`, false},
		{`@aws.region NOT IN ["cn-north-1"]`, true},
		{`@aws.region NOT IN ["us-east-1"]`, false},
		{`@aws.region = "us-east-1" AND @aws.account = "000000000000"`, true},
		{`@aws.region = "us-east-1" AND @aws.account = "111111111111"`, false},
		{`@aws.region = "eu-west-1" OR @aws.account = "000000000000"`, true},
		{`@aws.region = "eu-west-1" OR @aws.account = "000000000000" AND @aws.region = "us-east-1"`, true},
		{`@aws.region = "us-east-1" AND @aws.account = "000000000000" OR @aws.region = "eu-west-1"`, true},
		{`@nosuch.field = "x"`, false},
		// Operator fragments inside quoted values are value text, never
		// operators: each leaf below is well-formed against the
		// us-east-1/000000000000 system fields.
		{`@aws.region != "x IN y"`, true},
		{`@aws.region IN ["us-east-1", "a IN b"]`, true},
		{`@aws.region = "US IN X"`, false},
		{`@aws.account = "0 != 000000000000"`, false},
	}
	for _, tc := range cases {
		if got := evalFieldSelectionCriteria(tc.criteria, "us-east-1", "000000000000"); got != tc.expected {
			t.Errorf("evalFieldSelectionCriteria(%q) = %v, expected %v", tc.criteria, got, tc.expected)
		}
	}

	// The grammar validates at Put on both families: a malformed leaf and
	// a foreign system field reject; the documented forms store.
	svc, _ := newReadTestService(t, "matrix-fsc-group")
	base := []logsstore.MetricTransformation{{
		MetricName: "Count", MetricNamespace: "App", MetricValue: "1",
		MetricNameSet: true, MetricNamespaceSet: true, MetricValueSet: true,
	}}
	if err := svc.putMetricFilterCore("matrix-fsc-group", "ok", "ERROR", true, base, false,
		`@aws.region = "us-east-1" OR @aws.account IN ["000000000000"]`, nil, "us-east-1"); err != nil {
		t.Fatalf("documented criteria: %v", err)
	}
	if err := svc.putMetricFilterCore("matrix-fsc-group", "quoted-op", "ERROR", true, base, false,
		`@aws.region != "x IN y"`, nil, "us-east-1"); err != nil {
		t.Fatalf("operator fragment inside a quoted value is value text: %v", err)
	}
	err := svc.putMetricFilterCore("matrix-fsc-group", "bad-leaf", "ERROR", true, base, false,
		"@aws.region us-east-1", nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("malformed leaf: code=%q err=%v", code, err)
	}
	err = svc.putMetricFilterCore("matrix-fsc-group", "bad-field", "ERROR", true, base, false,
		`@aws.availabilityZone = "x"`, nil, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("foreign system field: code=%q err=%v", code, err)
	}

	// The metric gate: a filter whose criteria the batch's system fields
	// fail never emits, a passing one does.
	metrics := &valueRecordingMetricInvoker{}
	svc.SetCloudWatchMetricInvoker(metrics)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	gated := []logsstore.MetricTransformation{base[0]}
	gated[0].MetricName = "GatedCount"
	passing := []logsstore.MetricTransformation{base[0]}
	passing[0].MetricName = "PassingCount"
	if err := store.PutMetricFilter(&logsstore.MetricFilter{
		Name: "gated", LogGroupName: "matrix-fsc-group", FilterPattern: "ERROR",
		FieldSelectionCriteria: `@aws.region = "eu-west-1"`,
		MetricTransformations:  gated,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutMetricFilter(&logsstore.MetricFilter{
		Name: "passing", LogGroupName: "matrix-fsc-group", FilterPattern: "ERROR",
		FieldSelectionCriteria: `@aws.region = "us-east-1"`,
		MetricTransformations:  passing,
	}); err != nil {
		t.Fatal(err)
	}
	svc.evaluateMetricFilters(store, "us-east-1", "matrix-fsc-group", "s",
		[]logsstore.LogEntry{{Timestamp: 1, Message: "ERROR boom", IngestionTime: 1}}, nil)
	emitted := map[string]bool{}
	for _, c := range metrics.calls {
		emitted[c.metric] = true
	}
	if emitted["GatedCount"] {
		t.Error("a criteria-failing filter must not emit")
	}
	if !emitted["PassingCount"] {
		t.Error("a criteria-passing filter must emit")
	}

	// The subscription vocabulary: the documented three values accept and
	// ride the payload's events; a foreign value rejects at Put.
	if err := svc.putSubscriptionFilterCore(context.Background(), store, &PutSubscriptionFilterInput{
		LogGroupName: "matrix-fsc-group", FilterName: "sysfields", FilterPattern: "ERROR", FilterPatternSet: true,
		DestinationArn:   "arn:aws:kinesis:us-east-1:000000000000:stream/fsc-target",
		EmitSystemFields: []string{"@aws.account", "@aws.region", "@source.log"},
		Region:           "us-east-1",
	}); err != nil {
		t.Fatalf("documented system fields: %v", err)
	}
	err = svc.putSubscriptionFilterCore(context.Background(), store, &PutSubscriptionFilterInput{
		LogGroupName: "matrix-fsc-group", FilterName: "sysfields-bad", FilterPattern: "ERROR", FilterPatternSet: true,
		DestinationArn:   "arn:aws:kinesis:us-east-1:000000000000:stream/fsc-target",
		EmitSystemFields: []string{"@logStream"},
		Region:           "us-east-1",
	})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("foreign system field on subscription: code=%q err=%v", code, err)
	}
	filters, err := store.ListSubscriptionFilters("matrix-fsc-group", "")
	if err != nil || len(filters) == 0 {
		t.Fatalf("list subscription filters: %v (%d)", err, len(filters))
	}
	var sys *logsstore.SubscriptionFilter
	for _, f := range filters {
		if f.FilterName == "sysfields" {
			sys = f
		}
	}
	if sys == nil {
		t.Fatal("sysfields filter missing")
	}
	payload := svc.buildSubscriptionPayload(sys, "us-east-1", "matrix-fsc-group", "s",
		[]logsstore.LogEntry{{Timestamp: 1, Message: "ERROR boom", IngestionTime: 1}})
	events := payload["logEvents"].([]map[string]interface{})
	if got := events[0]["@aws.account"]; got != "000000000000" {
		t.Fatalf("@aws.account = %v, want the owner account", got)
	}
	if got := events[0]["@aws.region"]; got != "us-east-1" {
		t.Fatalf("@aws.region = %v, want the group's region", got)
	}
	if got := events[0]["@source.log"]; got != "matrix-fsc-group" {
		t.Fatalf("@source.log = %v, want the log group name", got)
	}
}

// --- FilterLogEvents: logStreamNames @length(min 1) ---

func TestFilterLogEventsEmptyStreamNamesRejects(t *testing.T) {
	svc, _ := newReadTestService(t, "matrix-fle-group")

	// An explicitly empty array violates the InputLogStreamNames length
	// trait's minimum; the non-nil empty slice is the present-member mark.
	_, err := svc.filterLogEventsCore(FilterLogEventsInput{
		LogGroupName: "matrix-fle-group", LogStreamNames: []string{}, Region: "us-east-1",
	})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("present empty array: code=%q err=%v", code, err)
	}

	// An absent array still means every stream in the group.
	if _, err := svc.filterLogEventsCore(FilterLogEventsInput{
		LogGroupName: "matrix-fle-group", Region: "us-east-1",
	}); err != nil {
		t.Fatalf("absent array: %v", err)
	}
}
