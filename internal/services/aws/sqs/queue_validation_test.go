package sqs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/kmsutil"
	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// The write-path validations below run before any store access, so a nil
// store is safe for the negative cases.

// TestQueueAttributeNameValidationOnWritePaths pins the InvalidAttributeName
// contract the model lists on CreateQueue and SetQueueAttributes: an
// attribute name outside the store's shared 22-name vocabulary is rejected
// on both write paths, on both planes' shared Core (the rejection was
// previously reachable only on the GET path).
func TestQueueAttributeNameValidationOnWritePaths(t *testing.T) {
	s := &SQSService{}
	bogus := map[string]string{"Bogus": "x"}

	if _, err := s.createQueueCore(context.Background(), nil, CreateQueueInput{
		QueueName: "queue-a",
		Region:    "us-east-1",
		Attrs:     bogus,
	}); err != ErrInvalidAttributeName {
		t.Errorf("createQueueCore: got %v, want ErrInvalidAttributeName", err)
	}

	if err := s.setQueueAttributesCore(context.Background(), nil, SetQueueAttributesInput{
		QueueURL: "http://localhost:50080/000000000000/queue-a",
		Attrs:    bogus,
	}); err != ErrInvalidAttributeName {
		t.Errorf("setQueueAttributesCore: got %v, want ErrInvalidAttributeName", err)
	}
}

// TestKmsMasterKeyIdFormatValidatedOnCorePaths pins the format rule that
// previously lived only in the store's SetQueueAttributes switch: a
// malformed KmsMasterKeyId is rejected by the Core on both write paths, so
// the admin gRPC plane validates it identically to the HTTP plane.
func TestKmsMasterKeyIdFormatValidatedOnCorePaths(t *testing.T) {
	s := &SQSService{}
	badKey := map[string]string{"KmsMasterKeyId": "not a valid key!"}

	if _, err := s.createQueueCore(context.Background(), nil, CreateQueueInput{
		QueueName: "queue-a",
		Region:    "us-east-1",
		Attrs:     badKey,
	}); err != ErrInvalidParameterValue {
		t.Errorf("createQueueCore: got %v, want ErrInvalidParameterValue", err)
	}

	if err := s.setQueueAttributesCore(context.Background(), nil, SetQueueAttributesInput{
		QueueURL: "http://localhost:50080/000000000000/queue-a",
		Attrs:    badKey,
	}); err != ErrInvalidParameterValue {
		t.Errorf("setQueueAttributesCore: got %v, want ErrInvalidParameterValue", err)
	}
}

// TestFifoQueueNamingCrossRuleShared pins the bidirectional queue-type
// naming rule through the shared exported store function the Core now uses:
// FifoQueue=true requires the ".fifo" suffix and vice-versa.
func TestFifoQueueNamingCrossRuleShared(t *testing.T) {
	s := &SQSService{}

	if _, err := s.createQueueCore(context.Background(), nil, CreateQueueInput{
		QueueName: "queue-a",
		Region:    "us-east-1",
		Attrs:     map[string]string{"FifoQueue": "true"},
	}); err != ErrInvalidParameterValue {
		t.Errorf("FifoQueue=true without .fifo suffix: got %v, want ErrInvalidParameterValue", err)
	}

	if _, err := s.createQueueCore(context.Background(), nil, CreateQueueInput{
		QueueName: "queue-a.fifo",
		Region:    "us-east-1",
		Attrs:     map[string]string{},
	}); err != ErrInvalidParameterValue {
		t.Errorf(".fifo suffix without FifoQueue=true: got %v, want ErrInvalidParameterValue", err)
	}
}

// TestRedrivePolicyClearOnCoreApply pins the services-layer clear mirror:
// applyQueueAttributes treats the empty RedrivePolicy value as a clear (nil
// typed policy, raw attribute removed), and the CreateQueue idempotence
// comparison matches the clear form against a queue without a policy.
func TestRedrivePolicyClearOnCoreApply(t *testing.T) {
	dlqARN := "arn:aws:sqs:us-east-1:123456789012:dlq"

	queue := sqsstore.NewQueue("clear-src", "us-east-1", "123456789012")
	rdp := fmt.Sprintf(`{"deadLetterTargetArn":%q,"maxReceiveCount":3}`, dlqARN)
	if err := applyQueueAttributes(map[string]string{"RedrivePolicy": rdp}, queue); err != nil {
		t.Fatalf("applyQueueAttributes set: %v", err)
	}
	if queue.RedrivePolicy == nil || queue.RedrivePolicy.DeadLetterTargetARN != dlqARN {
		t.Fatalf("applyQueueAttributes set: RedrivePolicy = %+v, want target %s", queue.RedrivePolicy, dlqARN)
	}

	if err := applyQueueAttributes(map[string]string{"RedrivePolicy": ""}, queue); err != nil {
		t.Fatalf("applyQueueAttributes clear: %v", err)
	}
	if queue.RedrivePolicy != nil {
		t.Errorf("applyQueueAttributes clear: RedrivePolicy = %+v, want nil", queue.RedrivePolicy)
	}
	if _, stillThere := queue.Attributes["RedrivePolicy"]; stillThere {
		t.Errorf("applyQueueAttributes clear: raw attribute map still carries %q", queue.Attributes["RedrivePolicy"])
	}

	if !requestAttrsMatchExisting(map[string]string{"RedrivePolicy": ""}, queue) {
		t.Error("requestAttrsMatchExisting: the clear form must match a queue without a redrive policy")
	}
	if requestAttrsMatchExisting(map[string]string{"RedrivePolicy": rdp}, queue) {
		t.Error("requestAttrsMatchExisting: an armed policy must not match a cleared queue")
	}
}

// fakeKMSChecker stands in for the platform's KMS key checker, returning a
// fixed error per key id so the Kms* error mapping is exercised without a
// KMS service.
type fakeKMSChecker map[string]error

func (f fakeKMSChecker) CheckKey(ctx context.Context, region, keyID string) error {
	if err, ok := f[keyID]; ok {
		return err
	}
	return nil
}

// TestMapKMSError pins the KMS error family: each kmsutil sentinel a key
// check can return maps to its documented Kms* wire error, and an unknown
// failure degrades to InvalidParameterValue.
func TestMapKMSError(t *testing.T) {
	cases := map[string]struct {
		err  error
		want string
	}{
		"not found":       {kmsutil.ErrKeyNotFound, "KmsNotFound"},
		"disabled":        {kmsutil.ErrKeyDisabled, "KmsDisabled"},
		"invalid state":   {kmsutil.ErrKeyInvalidState, "KmsInvalidState"},
		"invalid usage":   {kmsutil.ErrKeyInvalidUsage, "KmsInvalidKeyUsage"},
		"unknown failure": {errors.New("kaboom"), "InvalidParameterValue"},
	}
	for name, tc := range cases {
		got := mapKMSError(tc.err)
		apiErr, ok := got.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("%s: mapped to %T, want *AWSError", name, got)
		}
		if apiErr.GetCode() != tc.want {
			t.Fatalf("%s: code %q, want %q", name, apiErr.GetCode(), tc.want)
		}
	}
}

// TestGetQueueAttributesSelection pins the attribute-selection block of
// GetQueueAttributes on the real path: an explicit "All" returns the full
// set, a proper subset returns exactly the requested names, and an unknown
// name rejects with InvalidAttributeName (the selection block was
// previously unexercised — every test omitted AttributeNames).
func TestGetQueueAttributesSelection(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "gqa-select", map[string]interface{}{
		"VisibilityTimeout": "45",
	})

	call := func(names []interface{}) (map[string]string, error) {
		resp, err := svc.GetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
			Operation: "GetQueueAttributes",
			Parameters: map[string]interface{}{
				"QueueUrl":       queueURL,
				"AttributeNames": names,
			},
		})
		if err != nil {
			return nil, err
		}
		attrs, _ := resp.(map[string]interface{})["Attributes"].(map[string]string)
		return attrs, nil
	}

	all, err := call([]interface{}{"All"})
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if all["QueueArn"] == "" || all["VisibilityTimeout"] != "45" || all["CreatedTimestamp"] == "" {
		t.Fatalf("All returned an incomplete set: %v", all)
	}

	subset, err := call([]interface{}{"VisibilityTimeout"})
	if err != nil {
		t.Fatalf("subset: %v", err)
	}
	if len(subset) != 1 || subset["VisibilityTimeout"] != "45" {
		t.Fatalf("subset returned %v, want exactly VisibilityTimeout=45", subset)
	}

	// An omitted AttributeNames list is the documented empty-result request,
	// not a request for everything; "All" is the documented wildcard.
	omitted, err := call(nil)
	if err != nil {
		t.Fatalf("omitted: %v", err)
	}
	if len(omitted) != 0 {
		t.Fatalf("omitted AttributeNames returned %d attributes, want the documented empty result", len(omitted))
	}

	if _, err := call([]interface{}{"Bogus"}); err == nil || !strings.Contains(err.Error(), "InvalidAttributeName") {
		t.Fatalf("unknown name: expected InvalidAttributeName, got %v", err)
	}
}

// TestListDeadLetterSourceQueuesNextTokenGate pins the documented pagination
// contract: "You must set MaxResults to receive a value for NextToken in the
// response" — a truncated default page (1000 sources and one more) surfaces no
// token, an explicit MaxResults does, and following that token yields the
// remaining source.
func TestListDeadLetterSourceQueuesNextTokenGate(t *testing.T) {
	svc, _, store := newQueryWireTestService(t)

	dlq, err := store.CreateQueue(sqsstore.NewQueue("gate-token-dlq", "us-east-1", "123456789012"))
	if err != nil {
		t.Fatalf("create DLQ: %v", err)
	}
	for i := 0; i <= int(sqsstore.MaxListResults); i++ {
		q := sqsstore.NewQueue(fmt.Sprintf("gate-token-src-%04d", i), "us-east-1", "123456789012")
		policy := fmt.Sprintf(`{"deadLetterTargetArn":%q,"maxReceiveCount":3}`, dlq.ARN)
		q.Attributes["RedrivePolicy"] = policy
		q.RedrivePolicy = &sqsstore.RedrivePolicy{DeadLetterTargetARN: dlq.ARN, MaxReceiveCount: 3}
		if _, err := store.CreateQueue(q); err != nil {
			t.Fatalf("create source %04d: %v", i, err)
		}
	}

	omitted, err := svc.listDeadLetterSourceQueuesCore(store, ListDeadLetterSourceQueuesInput{QueueURL: dlq.URL})
	if err != nil {
		t.Fatalf("omitted MaxResults: %v", err)
	}
	if len(omitted.QueueURLs) != int(sqsstore.MaxListResults) {
		t.Fatalf("default page returned %d sources, want the full %d", len(omitted.QueueURLs), sqsstore.MaxListResults)
	}
	if omitted.NextToken != "" {
		t.Fatalf("omitted MaxResults surfaced NextToken %q; the token requires an explicit MaxResults", omitted.NextToken)
	}

	explicit, err := svc.listDeadLetterSourceQueuesCore(store, ListDeadLetterSourceQueuesInput{
		QueueURL:      dlq.URL,
		MaxResults:    int32(sqsstore.MaxListResults),
		MaxResultsSet: true,
	})
	if err != nil {
		t.Fatalf("explicit MaxResults: %v", err)
	}
	if explicit.NextToken == "" {
		t.Fatal("explicit MaxResults did not surface NextToken on a truncated result")
	}

	page2, err := svc.listDeadLetterSourceQueuesCore(store, ListDeadLetterSourceQueuesInput{
		QueueURL:      dlq.URL,
		MaxResults:    int32(sqsstore.MaxListResults),
		MaxResultsSet: true,
		NextToken:     explicit.NextToken,
	})
	if err != nil {
		t.Fatalf("second page: %v", err)
	}
	if len(page2.QueueURLs) != 1 {
		t.Fatalf("second page returned %d sources, want the 1 remaining", len(page2.QueueURLs))
	}
	if page2.NextToken != "" {
		t.Fatalf("exhausted pagination surfaced NextToken %q", page2.NextToken)
	}
}

// TestWritePathsRejectReadOnlyAttributeNames pins the write-vocabulary
// split: the write paths (CreateQueue, SetQueueAttributes) reject the
// read-only names — "All" and the computed/report-only attributes — with
// InvalidAttributeName (the error SetQueueAttributes documents for unknown
// attribute names); the read path keeps accepting the full vocabulary.
func TestWritePathsRejectReadOnlyAttributeNames(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	created, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":        "write-vocab",
			"Attribute.1.Name": "CreatedTimestamp", "Attribute.1.Value": "1700000000",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "InvalidAttributeName") {
		t.Fatalf("CreateQueue with a read-only name: expected InvalidAttributeName, got %v", err)
	}
	if created != nil {
		t.Fatal("CreateQueue with a read-only name returned a result")
	}

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "write-vocab-2"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	if _, err := svc.SetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SetQueueAttributes",
		Parameters: map[string]interface{}{
			"QueueUrl":         queueURL,
			"Attribute.1.Name": "All", "Attribute.1.Value": "1",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidAttributeName") {
		t.Fatalf("SetQueueAttributes with All: expected InvalidAttributeName, got %v", err)
	}

	// The read side keeps the full vocabulary: "All" is the documented
	// wildcard there.
	if _, err := svc.GetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "GetQueueAttributes",
		Parameters: map[string]interface{}{"QueueUrl": queueURL, "AttributeName.1": "All"},
	}); err != nil {
		t.Fatalf("GetQueueAttributes with All on the read path: %v", err)
	}
}

// TestSetQueueAttributesExistencePrecedesKMSCheck pins the error precedence:
// a request naming a nonexistent queue answers QueueDoesNotExist even when
// the KMS key it carries is bad — the KMS check is a platform extension on
// the request's attribute payload, not the request's target.
func TestSetQueueAttributesExistencePrecedesKMSCheck(t *testing.T) {
	svc, _, store := newQueryWireTestService(t)
	svc.kmsChecker = fakeKMSChecker{"alias/bad-key": kmsutil.ErrKeyNotFound}

	err := svc.setQueueAttributesCore(context.Background(), store, SetQueueAttributesInput{
		QueueURL: "http://localhost:50080/123456789012/no-such-queue",
		Region:   "us-east-1",
		Attrs:    map[string]string{"KmsMasterKeyId": "alias/bad-key"},
	})
	if err == nil || !strings.Contains(err.Error(), "QueueDoesNotExist") {
		t.Fatalf("nonexistent queue with a bad key: expected QueueDoesNotExist, got %v", err)
	}
}

// TestBatchFailedEntryCarriesMappedMessage pins the Failed-entry text: the
// Message is the mapped single-operation error's wording, not the store
// sentinel's internal phrasing.
func TestBatchFailedEntryCarriesMappedMessage(t *testing.T) {
	entry := batchFailedEntry("entry-1", sqsstore.ErrMessageTooLarge)
	message, _ := entry["Message"].(string)
	if !strings.Contains(message, "Message must be shorter than") {
		t.Fatalf("Failed-entry Message = %q; want the mapped AWS wording quoting the byte bound", message)
	}
	if message == sqsstore.ErrMessageTooLarge.Error() {
		t.Fatalf("Failed-entry Message is the store sentinel text %q", message)
	}
}
