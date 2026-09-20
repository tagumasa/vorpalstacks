package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"vorpalstacks/internal/common/request"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// createFifoTopic creates a FIFO topic with content-based deduplication
// left off, so each pin controls its own deduplication identifiers.
func createFifoTopic(t *testing.T, store snsstore.SNSStoreInterface, name string) *snsstore.Topic {
	t.Helper()
	topic, err := store.CreateTopic(&snsstore.Topic{
		Name: name,
		Attributes: map[string]string{
			snsstore.AttrFifoTopic: "true",
		},
	}, nil)
	if err != nil {
		t.Fatalf("create FIFO topic: %v", err)
	}
	return topic
}

// confirmHTTPSubscription creates one confirmed http subscription against a
// local server that records the arrival order of every notification's
// MessageId. The returned slice pointer is appended under the captured
// mutex; the caller reads it after Close or between synchronous publishes.
func confirmHTTPSubscription(t *testing.T, store snsstore.SNSStoreInterface, topicArn string) (*httptest.Server, *[]string, *sync.Mutex) {
	t.Helper()
	mu := &sync.Mutex{}
	arrived := &[]string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var envelope map[string]interface{}
		_ = json.Unmarshal(body, &envelope)
		mu.Lock()
		*arrived = append(*arrived, envelope["MessageId"].(string))
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	created, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: topicArn,
		Protocol: "http",
		Endpoint: server.URL,
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(created); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}
	return server, arrived, mu
}

// TestFifoPublishDeliveryOrderMatchesSequenceOrder pins the ordering
// design: concurrent publishes to one message group deliver in sequence
// order — the arrival order every endpoint observes is the order the
// SequenceNumbers sort into, not the order the publishes raced in.
func TestFifoPublishDeliveryOrderMatchesSequenceOrder(t *testing.T) {
	svc, store := newFanoutTestService(t)
	topic := createFifoTopic(t, store, "order.fifo")
	_, arrived, mu := confirmHTTPSubscription(t, store, topic.Arn)

	const publishers = 6
	type publishResult struct {
		messageID      string
		sequenceNumber string
	}
	results := make(chan publishResult, publishers)

	var start sync.WaitGroup
	start.Add(1)
	var wg sync.WaitGroup
	for i := 0; i < publishers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start.Wait()
			result, err := svc.publishCore(store, "us-east-1", PublishInput{
				TopicArn:               topic.Arn,
				Message:                fmt.Sprintf("ordered-%d", i),
				MessageGroupId:         "group-a",
				MessageDeduplicationId: fmt.Sprintf("dedup-%d", i),
			})
			if err != nil {
				t.Errorf("publish %d: %v", i, err)
				return
			}
			messageID, _ := result.(map[string]interface{})["MessageId"].(string)
			sequenceNumber, _ := result.(map[string]interface{})["SequenceNumber"].(string)
			results <- publishResult{messageID: messageID, sequenceNumber: sequenceNumber}
		}(i)
	}
	start.Done()
	wg.Wait()
	close(results)

	sequences := make(map[string]int64, publishers)
	for result := range results {
		if result.messageID == "" || result.sequenceNumber == "" {
			t.Fatalf("publish result missing identifiers: %+v", result)
		}
		value, err := strconv.ParseInt(result.sequenceNumber, 10, 64)
		if err != nil {
			t.Fatalf("sequence number %q: %v", result.sequenceNumber, err)
		}
		sequences[result.messageID] = value
	}

	mu.Lock()
	order := make([]string, len(*arrived))
	copy(order, *arrived)
	mu.Unlock()
	if len(order) != publishers {
		t.Fatalf("%d of %d messages arrived — FIFO delivery lost messages", len(order), publishers)
	}
	previous := int64(0)
	for i, messageID := range order {
		sequence, ok := sequences[messageID]
		if !ok {
			t.Fatalf("arrival %d carries unknown MessageId %q", i, messageID)
		}
		if sequence <= previous {
			t.Fatalf("arrival %d (MessageId %s, SequenceNumber %d) did not follow sequence order (previous %d)", i, messageID, sequence, previous)
		}
		previous = sequence
	}
}

// TestFifoDedupHitReturnsOriginalIdentifiers pins the dedup-hit response:
// the duplicate publish is accepted but not delivered and answers with the
// ORIGINAL publish's MessageId and SequenceNumber — the sequence number is
// assigned per message, and a dedup hit is that message.
func TestFifoDedupHitReturnsOriginalIdentifiers(t *testing.T) {
	svc, store := newFanoutTestService(t)
	topic := createFifoTopic(t, store, "dedup-hit.fifo")
	_, arrived, mu := confirmHTTPSubscription(t, store, topic.Arn)

	publish := func() (string, string) {
		t.Helper()
		result, err := svc.publishCore(store, "us-east-1", PublishInput{
			TopicArn:               topic.Arn,
			Message:                "dedup-hit body",
			MessageGroupId:         "group-a",
			MessageDeduplicationId: "same-dedup",
		})
		if err != nil {
			t.Fatalf("publish: %v", err)
		}
		messageID, _ := result.(map[string]interface{})["MessageId"].(string)
		sequenceNumber, _ := result.(map[string]interface{})["SequenceNumber"].(string)
		return messageID, sequenceNumber
	}

	messageID, sequenceNumber := publish()
	retryID, retrySequence := publish()

	if retryID != messageID {
		t.Fatalf("dedup hit returned MessageId %q, want the original %q", retryID, messageID)
	}
	if retrySequence != sequenceNumber {
		t.Fatalf("dedup hit returned SequenceNumber %q, want the original %q", retrySequence, sequenceNumber)
	}
	if sequenceNumber == "" {
		t.Fatal("FIFO publish response carries no SequenceNumber")
	}

	mu.Lock()
	delivered := len(*arrived)
	mu.Unlock()
	if delivered != 1 {
		t.Fatalf("%d deliveries for one deduplicated message, want 1", delivered)
	}
}

// TestStandardTopicAcceptsMessageGroupId pins the documented fair-queues
// behaviour: MessageGroupId is optional on standard topics — accepted and
// forwarded to SQS subscriptions — while MessageDeduplicationId applies
// only to FIFO topics.
func TestStandardTopicAcceptsMessageGroupId(t *testing.T) {
	svc, store := newFanoutTestService(t)
	topic, err := store.CreateTopic(&snsstore.Topic{Name: "fair-queues-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	result, err := svc.publishCore(store, "us-east-1", PublishInput{
		TopicArn:       topic.Arn,
		Message:        "fair queues message",
		MessageGroupId: "fair-group",
	})
	if err != nil {
		t.Fatalf("publish with MessageGroupId on a standard topic: %v", err)
	}
	row := result.(map[string]interface{})
	if id, _ := row["MessageId"].(string); id == "" {
		t.Fatal("standard publish response carries no MessageId")
	}
	if _, has := row["SequenceNumber"]; has {
		t.Fatal("standard publish response carries SequenceNumber — the member applies only to FIFO topics")
	}

	if _, err := svc.publishCore(store, "us-east-1", PublishInput{
		TopicArn:               topic.Arn,
		Message:                "m",
		MessageDeduplicationId: "dedup",
	}); err == nil {
		t.Fatal("MessageDeduplicationId was accepted on a standard topic — it applies only to FIFO topics")
	}
}

// TestPublishBatchFifoDedupHitCarriesSequence pins the batch twin: two
// entries sharing a deduplication identifier admit exactly one copy — the
// first entry delivers with its own identifiers, the second answers with
// the first's MessageId and SequenceNumber.
func TestPublishBatchFifoDedupHitCarriesSequence(t *testing.T) {
	svc, store := newFanoutTestService(t)
	topic := createFifoTopic(t, store, "batch-dedup.fifo")

	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	result, err := svc.publishBatchCore(store, reqCtx, PublishBatchInput{
		TopicArn: topic.Arn,
		Entries: []map[string]interface{}{
			{"Id": "a", "Message": "first", "MessageGroupId": "group-a", "MessageDeduplicationId": "shared"},
			{"Id": "b", "Message": "second", "MessageGroupId": "group-a", "MessageDeduplicationId": "shared"},
		},
	})
	if err != nil {
		t.Fatalf("publish batch: %v", err)
	}
	rows := result.(map[string]interface{})["Successful"].([]map[string]interface{})
	if len(rows) != 2 {
		t.Fatalf("batch Successful rows = %d, want 2 (both entries accepted)", len(rows))
	}
	first, second := rows[0], rows[1]
	if first["Id"] != "a" || second["Id"] != "b" {
		t.Fatalf("batch rows out of order: %v", rows)
	}
	if first["SequenceNumber"] == "" || first["MessageId"] == "" {
		t.Fatalf("first row missing identifiers: %v", first)
	}
	if second["MessageId"] != first["MessageId"] {
		t.Fatalf("dedup-hit row MessageId %v, want the original %v", second["MessageId"], first["MessageId"])
	}
	if second["SequenceNumber"] != first["SequenceNumber"] {
		t.Fatalf("dedup-hit row SequenceNumber %v, want the original %v", second["SequenceNumber"], first["SequenceNumber"])
	}
}

// dedupFailsOnSecondCheckStore passes the first deduplication check (the
// batch's validation pass) and fails every later one — the harness that
// drives a per-entry Pass-2 failure.
type dedupFailsOnSecondCheckStore struct {
	snsstore.SNSStoreInterface
	calls int
}

func (w *dedupFailsOnSecondCheckStore) CheckFifoDeduplication(topicArn, group, id string, perGroup bool) (string, string, bool, error) {
	w.calls++
	if w.calls > 1 {
		return "", "", false, fmt.Errorf("injected dedup fault")
	}
	return w.SNSStoreInterface.CheckFifoDeduplication(topicArn, group, id, perGroup)
}

// TestPublishBatchFailureRowCodeIsInternalError pins the failed-entry
// vocabulary: a Pass-2 internal failure reports the model's InternalError
// code (InternalErrorException's awsQueryError trait) — the same word the
// dispatcher's internal mapping uses — never a foreign spelling.
func TestPublishBatchFailureRowCodeIsInternalError(t *testing.T) {
	svc, store := newFanoutTestService(t)
	topic := createFifoTopic(t, store, "batch-fault.fifo")
	wrapped := &dedupFailsOnSecondCheckStore{SNSStoreInterface: store}

	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	result, err := svc.publishBatchCore(wrapped, reqCtx, PublishBatchInput{
		TopicArn: topic.Arn,
		Entries: []map[string]interface{}{
			{"Id": "a", "Message": "one", "MessageGroupId": "g", "MessageDeduplicationId": "d1"},
		},
	})
	if err != nil {
		t.Fatalf("publish batch: %v", err)
	}
	failed := result.(map[string]interface{})["Failed"].([]map[string]interface{})
	if len(failed) != 1 {
		t.Fatalf("Failed rows = %d, want 1 (the Pass-2 fault)", len(failed))
	}
	if failed[0]["Code"] != "InternalError" {
		t.Fatalf("Failed row Code = %v, want InternalError", failed[0]["Code"])
	}
}
