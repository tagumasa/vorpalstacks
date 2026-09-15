package sqs

import (
	"sort"
	"testing"

	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// The output-shape vocabulary of the five DTO-typed operations, transcribed
// from the service model (sqs-2012-11-05, jq census at fix time). Every
// response member key these operations emit must appear here with exactly
// this casing — ListDeadLetterSourceQueues carries the model's own
// asymmetry, queueUrls, while ListQueues emits QueueUrls. The model defines
// no other members for these shapes, so an emitted key outside its entry's
// vocabulary is a fix-introduced defect.
var modelOutputVocabulary = map[string][]string{
	"GetQueueAttributes":         {"Attributes"},
	"ListDeadLetterSourceQueues": {"queueUrls", "NextToken"},
	"StartMessageMoveTask":       {"TaskHandle"},
	"CancelMessageMoveTask":      {"ApproximateNumberOfMessagesMoved"},
	"ListMessageMoveTasks":       {"Results"},
	"ListMessageMoveTasksEntry": {
		"TaskHandle", "Status", "SourceArn", "DestinationArn",
		"MaxNumberOfMessagesPerSecond", "ApproximateNumberOfMessagesMoved",
		"ApproximateNumberOfMessagesToMove", "FailureReason", "StartedTimestamp",
	},
	// The received Message shape (jq census of com.amazonaws.sqs#Message):
	// exactly seven members. FIFO identifiers are NOT members — they travel
	// as entries of the Attributes map when requested.
	"ReceiveMessage Message": {
		"MessageId", "ReceiptHandle", "MD5OfBody", "Body",
		"MD5OfMessageAttributes", "MessageAttributes", "Attributes",
	},
}

// alwaysPresentMoveTaskEntryKeys are emitted for every task entry; the other
// entry members are conditional on being set.
var alwaysPresentMoveTaskEntryKeys = []string{
	"Status", "SourceArn", "ApproximateNumberOfMessagesMoved", "StartedTimestamp",
}

func responseKeySet(t *testing.T, resp map[string]interface{}) []string {
	t.Helper()
	keys := make([]string, 0, len(resp))
	for k := range resp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func assertKeySet(t *testing.T, context string, got []string, want ...string) {
	t.Helper()
	sorted := append([]string(nil), want...)
	sort.Strings(sorted)
	if len(got) != len(sorted) {
		t.Fatalf("%s: emitted keys %v, want exactly %v", context, got, sorted)
	}
	for i := range got {
		if got[i] != sorted[i] {
			t.Fatalf("%s: emitted keys %v, want exactly %v", context, got, sorted)
		}
	}
}

// TestResponseVocabularyConversionPins walks the DTO→wire conversions with
// representative DTOs and asserts the exact emitted key sets against the
// model vocabulary above, including the conditional members of a
// ListMessageMoveTasks entry (full DTO emits every member, minimal DTO only
// the always-present four).
func TestResponseVocabularyConversionPins(t *testing.T) {
	fullEntry := MessageMoveTaskDescription{
		TaskHandle:                        "task-1",
		Status:                            "RUNNING",
		SourceArn:                         "arn:aws:sqs:us-east-1:123456789012:src",
		DestinationArn:                    "arn:aws:sqs:us-east-1:123456789012:dst",
		MaxNumberOfMessagesPerSecond:      50,
		ApproximateNumberOfMessagesMoved:  10,
		ApproximateNumberOfMessagesToMove: 100,
		FailureReason:                     "boom",
		StartedTimestamp:                  1700000000000,
	}
	cases := []struct {
		name string
		resp map[string]interface{}
		want []string
	}{
		{"GetQueueAttributes", getQueueAttributesResponse(&GetQueueAttributesResult{
			Attributes: map[string]string{"QueueArn": "arn"},
		}), modelOutputVocabulary["GetQueueAttributes"]},
		{"ListDeadLetterSourceQueues unpaginated", listDeadLetterSourceQueuesResponse(&ListDeadLetterSourceQueuesResult{
			QueueURLs: []string{"http://q"},
		}), []string{"queueUrls"}},
		{"ListDeadLetterSourceQueues paginated", listDeadLetterSourceQueuesResponse(&ListDeadLetterSourceQueuesResult{
			QueueURLs: []string{"http://q"}, NextToken: "marker",
		}), modelOutputVocabulary["ListDeadLetterSourceQueues"]},
		{"StartMessageMoveTask", startMessageMoveTaskResponse(&StartMessageMoveTaskResult{
			TaskHandle: "task-1",
		}), modelOutputVocabulary["StartMessageMoveTask"]},
		{"CancelMessageMoveTask", cancelMessageMoveTaskResponse(&CancelMessageMoveTaskResult{
			ApproximateNumberOfMessagesMoved: 7,
		}), modelOutputVocabulary["CancelMessageMoveTask"]},
	}
	for _, tc := range cases {
		assertKeySet(t, tc.name, responseKeySet(t, tc.resp), tc.want...)
	}

	full := listMessageMoveTasksResponse(&ListMessageMoveTasksResult{Results: []MessageMoveTaskDescription{fullEntry}})
	assertKeySet(t, "ListMessageMoveTasks", responseKeySet(t, full), modelOutputVocabulary["ListMessageMoveTasks"]...)
	fullEntries := full["Results"].([]map[string]interface{})
	assertKeySet(t, "ListMessageMoveTasks full entry", responseKeySet(t, fullEntries[0]), modelOutputVocabulary["ListMessageMoveTasksEntry"]...)

	minimal := listMessageMoveTasksResponse(&ListMessageMoveTasksResult{Results: []MessageMoveTaskDescription{{
		Status:           "RUNNING",
		SourceArn:        "arn:aws:sqs:us-east-1:123456789012:src",
		StartedTimestamp: 1700000000000,
	}}})
	minimalEntries := minimal["Results"].([]map[string]interface{})
	assertKeySet(t, "ListMessageMoveTasks minimal entry", responseKeySet(t, minimalEntries[0]), alwaysPresentMoveTaskEntryKeys...)
}

// newMoveTaskVocabularyEnv builds the query-wire service scaffold with one
// redrive pair seeded: a source queue whose RedrivePolicy targets the DLQ,
// and one message held in flight in the DLQ so a started move task has work
// outstanding, the listed task reports a non-zero
// ApproximateNumberOfMessagesToMove, and the cancel pin below lands on a
// RUNNING task deterministically instead of racing the worker's completion.
func newMoveTaskVocabularyEnv(t *testing.T) (*SQSService, sqsstore.SQSStoreInterface, string, string, string) {
	t.Helper()
	svc, _, store := newQueryWireTestService(t)
	dlq, err := store.CreateQueue(sqsstore.NewQueue("vocab-dlq", "us-east-1", "123456789012"))
	if err != nil {
		t.Fatalf("create DLQ: %v", err)
	}
	source, err := store.CreateQueue(sqsstore.NewQueue("vocab-source", "us-east-1", "123456789012"))
	if err != nil {
		t.Fatalf("create source: %v", err)
	}
	rdp := `{"deadLetterTargetArn":"` + dlq.ARN + `","maxReceiveCount":3}`
	if err := store.SetQueueAttributes(source.URL, map[string]string{"RedrivePolicy": rdp}); err != nil {
		t.Fatalf("set redrive policy: %v", err)
	}
	if _, err := store.SendMessage(dlq.URL, sqsstore.NewMessage("redriven payload")); err != nil {
		t.Fatalf("seed DLQ message: %v", err)
	}
	// Hold the seeded message in flight under a long visibility timeout:
	// the move task's only message is then in flight, so the worker polls
	// and the task stays RUNNING.
	vt := int32(60)
	if _, err := store.ReceiveMessage(dlq.URL, 1, &vt, 0, ""); err != nil {
		t.Fatalf("hold DLQ message in flight: %v", err)
	}
	return svc, store, source.URL, dlq.URL, dlq.ARN
}

// TestResponseVocabularyRealPath exercises the five DTO-typed cores and
// their conversions through the real store, asserting the emitted response
// keys and casing against the model vocabulary: every key of the
// GetQueueAttributes attribute map is a known attribute name, the move-task
// entries carry only modelled members, and the queueUrls casing asymmetry
// survives the DTO round-trip.
func TestResponseVocabularyRealPath(t *testing.T) {
	s, store, sourceURL, dlqURL, dlqARN := newMoveTaskVocabularyEnv(t)

	// "All" keeps this a full-vocabulary diff: an omitted AttributeNames is
	// the documented empty-result request and would leave the attribute-name
	// loop below vacuous.
	attrs, err := s.getQueueAttributesCore(store, GetQueueAttributesInput{QueueURL: sourceURL, AttributeNames: []string{"All"}})
	if err != nil {
		t.Fatalf("getQueueAttributesCore: %v", err)
	}
	attrsResp := getQueueAttributesResponse(attrs)
	assertKeySet(t, "GetQueueAttributes real path", responseKeySet(t, attrsResp), modelOutputVocabulary["GetQueueAttributes"]...)
	for name := range attrsResp["Attributes"].(map[string]string) {
		if !sqsstore.IsValidAttributeName(name) {
			t.Errorf("GetQueueAttributes emitted unknown attribute name %q", name)
		}
	}

	sources, err := s.listDeadLetterSourceQueuesCore(store, ListDeadLetterSourceQueuesInput{QueueURL: dlqURL})
	if err != nil {
		t.Fatalf("listDeadLetterSourceQueuesCore: %v", err)
	}
	sourcesResp := listDeadLetterSourceQueuesResponse(sources)
	assertKeySet(t, "ListDeadLetterSourceQueues real path", responseKeySet(t, sourcesResp), "queueUrls")
	if urls := sourcesResp["queueUrls"].([]string); len(urls) != 1 || urls[0] != sourceURL {
		t.Errorf("ListDeadLetterSourceQueues result = %v, want [%s]", urls, sourceURL)
	}

	started, err := s.startMessageMoveTaskCore(store, StartMessageMoveTaskInput{SourceARN: dlqARN})
	if err != nil {
		t.Fatalf("startMessageMoveTaskCore: %v", err)
	}
	assertKeySet(t, "StartMessageMoveTask real path",
		responseKeySet(t, startMessageMoveTaskResponse(started)), modelOutputVocabulary["StartMessageMoveTask"]...)

	listed, err := s.listMessageMoveTasksCore(store, ListMessageMoveTasksInput{SourceARN: dlqARN})
	if err != nil {
		t.Fatalf("listMessageMoveTasksCore: %v", err)
	}
	listedResp := listMessageMoveTasksResponse(listed)
	assertKeySet(t, "ListMessageMoveTasks real path", responseKeySet(t, listedResp), modelOutputVocabulary["ListMessageMoveTasks"]...)
	entries := listedResp["Results"].([]map[string]interface{})
	if len(entries) == 0 {
		t.Fatal("ListMessageMoveTasks real path: no task entries")
	}
	vocab := modelOutputVocabulary["ListMessageMoveTasksEntry"]
	allowed := make(map[string]bool, len(vocab))
	for _, k := range vocab {
		allowed[k] = true
	}
	for k := range entries[0] {
		if !allowed[k] {
			t.Errorf("ListMessageMoveTasks entry emitted unmodelled key %q", k)
		}
	}
	for _, k := range alwaysPresentMoveTaskEntryKeys {
		if _, ok := entries[0][k]; !ok {
			t.Errorf("ListMessageMoveTasks entry missing always-present key %q", k)
		}
	}

	// The env holds the task's only message in flight, so the task stays
	// RUNNING and the cancel lands deterministically — no retry loop, no
	// skippable pin.
	cancelled, cerr := s.cancelMessageMoveTaskCore(store, CancelMessageMoveTaskInput{TaskHandle: started.TaskHandle})
	if cerr != nil {
		t.Fatalf("cancelMessageMoveTaskCore: %v", cerr)
	}
	assertKeySet(t, "CancelMessageMoveTask real path",
		responseKeySet(t, cancelMessageMoveTaskResponse(cancelled)), modelOutputVocabulary["CancelMessageMoveTask"]...)
}
