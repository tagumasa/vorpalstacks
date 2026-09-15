package sqs

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// The ReceiveMessage emission contract: the attribute name lists are
// opt-in ("you can send a list of attribute names to receive, or you can
// return all of the attributes by specifying All or .* in your request"),
// MD5OfMessageAttributes is the digest of the attributes actually returned,
// and the received Message shape carries no top-level FIFO identifier
// members — those travel as system attributes.

// receiveOne sends one message and receives it once with the given extra
// parameters, returning the emitted message map.
func receiveOne(t *testing.T, svc *SQSService, reqCtx *request.RequestContext, queueURL string, sendExtra, recvExtra map[string]interface{}) map[string]interface{} {
	t.Helper()

	sendParams := map[string]interface{}{
		"QueueUrl":    queueURL,
		"MessageBody": "emission probe",
	}
	for k, v := range sendExtra {
		sendParams[k] = v
	}
	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "SendMessage",
		Parameters: sendParams,
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	recvParams := map[string]interface{}{"QueueUrl": queueURL}
	for k, v := range recvExtra {
		recvParams[k] = v
	}
	recvResp, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "ReceiveMessage",
		Parameters: recvParams,
	})
	if err != nil {
		t.Fatalf("ReceiveMessage: %v", err)
	}
	recvMap, ok := recvResp.(map[string]interface{})
	if !ok {
		t.Fatalf("ReceiveMessage response type %T, want map", recvResp)
	}
	messages, ok := recvMap["Messages"].([]map[string]interface{})
	if !ok || len(messages) != 1 {
		t.Fatalf("ReceiveMessage returned %v, want exactly one message", recvMap)
	}
	return messages[0]
}

// createEmissionQueue creates a standard queue for the emission pins.
func createEmissionQueue(t *testing.T, svc *SQSService, reqCtx *request.RequestContext, name string, extraAttrs map[string]interface{}) string {
	t.Helper()

	params := map[string]interface{}{"QueueName": name}
	i := 1
	for k, v := range extraAttrs {
		params["Attribute."+strconv.Itoa(i)+".Name"] = k
		params["Attribute."+strconv.Itoa(i)+".Value"] = v
		i++
	}
	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: params,
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	return resp.(map[string]interface{})["QueueUrl"].(string)
}

// TestReceiveMessageEmissionDefaultOmitsAttributes pins the opt-in rule:
// a receive that names no attribute list returns the bare message — no
// Attributes, no MessageAttributes and no MD5OfMessageAttributes — even
// though the message carries both kinds.
func TestReceiveMessageEmissionDefaultOmitsAttributes(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "emit-default", nil)

	msg := receiveOne(t, svc, reqCtx, queueURL,
		map[string]interface{}{
			"MessageAttributes": map[string]interface{}{
				"Attr1": map[string]interface{}{"DataType": "String", "StringValue": "v1"},
				"Attr2": map[string]interface{}{"DataType": "String", "StringValue": "v2"},
			},
		}, nil)

	for _, key := range []string{"Attributes", "MessageAttributes", "MD5OfMessageAttributes"} {
		if _, ok := msg[key]; ok {
			t.Errorf("default receive emitted %q, want it omitted (the name lists are opt-in)", key)
		}
	}
	assertMessageVocabulary(t, "default receive", msg)
}

// TestReceiveMessageEmissionAllRequested pins the All wildcard: both
// attribute kinds are returned, with the full-set digest as
// MD5OfMessageAttributes.
func TestReceiveMessageEmissionAllRequested(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "emit-all", nil)

	msg := receiveOne(t, svc, reqCtx, queueURL,
		map[string]interface{}{
			"MessageAttributes": map[string]interface{}{
				"Attr1": map[string]interface{}{"DataType": "String", "StringValue": "v1"},
				"Attr2": map[string]interface{}{"DataType": "String", "StringValue": "v2"},
			},
		},
		map[string]interface{}{
			"MessageAttributeNames": []interface{}{"All"},
			"AttributeNames":        []interface{}{"All"},
		})

	attrs := msg["MessageAttributes"].(map[string]interface{})
	if len(attrs) != 2 {
		t.Fatalf("All-requested receive returned %d attributes, want 2", len(attrs))
	}
	if md5, _ := msg["MD5OfMessageAttributes"].(string); md5 == "" {
		t.Error("All-requested receive omitted MD5OfMessageAttributes")
	}
	sysAttrs := msg["Attributes"].(map[string]string)
	for _, want := range []string{"SentTimestamp", "SenderId", "ApproximateReceiveCount"} {
		if _, ok := sysAttrs[want]; !ok {
			t.Errorf("All-requested receive missing system attribute %q", want)
		}
	}
	assertMessageVocabulary(t, "All-requested receive", msg)
}

// TestReceiveMessageEmissionFilteredSubsetMD5 pins the MD5 basis: a
// filtered receive returns only the requested attribute and
// MD5OfMessageAttributes is the digest of THAT subset — not the send-time
// full-set digest.
func TestReceiveMessageEmissionFilteredSubsetMD5(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "emit-filter", nil)

	sendResp, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":    queueURL,
			"MessageBody": "emission probe",
			"MessageAttributes": map[string]interface{}{
				"Attr1": map[string]interface{}{"DataType": "String", "StringValue": "v1"},
				"Attr2": map[string]interface{}{"DataType": "String", "StringValue": "v2"},
			},
		},
	})
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	fullMD5 := sendResp.(map[string]interface{})["MD5OfMessageAttributes"].(string)

	recvResp, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ReceiveMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":              queueURL,
			"MessageAttributeNames": []interface{}{"Attr1"},
		},
	})
	if err != nil {
		t.Fatalf("ReceiveMessage: %v", err)
	}
	msg := recvResp.(map[string]interface{})["Messages"].([]map[string]interface{})[0]

	attrs := msg["MessageAttributes"].(map[string]interface{})
	if len(attrs) != 1 || attrs["Attr1"] == nil {
		t.Fatalf("filtered receive returned %v, want exactly Attr1", attrs)
	}
	subsetMD5 := msg["MD5OfMessageAttributes"].(string)
	v1 := "v1"
	want := sqsstore.CalculateMessageAttributesMD5(map[string]*sqsstore.MessageAttributeValue{
		"Attr1": {DataType: "String", StringValue: &v1},
	})
	if subsetMD5 != want {
		t.Errorf("filtered MD5OfMessageAttributes = %q, want digest of the returned subset %q", subsetMD5, want)
	}
	if subsetMD5 == fullMD5 {
		t.Error("filtered MD5OfMessageAttributes equals the full-set digest; the basis must be the returned subset")
	}
	assertMessageVocabulary(t, "filtered receive", msg)
}

// TestReceiveMessageEmissionFilterEmptyOmitsMD5 pins the empty-filter
// emission: a name list that matches nothing yields no MessageAttributes
// and no MD5OfMessageAttributes.
func TestReceiveMessageEmissionFilterEmptyOmitsMD5(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "emit-empty-filter", nil)

	msg := receiveOne(t, svc, reqCtx, queueURL,
		map[string]interface{}{
			"MessageAttributes": map[string]interface{}{
				"Attr1": map[string]interface{}{"DataType": "String", "StringValue": "v1"},
			},
		},
		map[string]interface{}{"MessageAttributeNames": []interface{}{"NoSuchAttr"}},
	)

	for _, key := range []string{"MessageAttributes", "MD5OfMessageAttributes"} {
		if _, ok := msg[key]; ok {
			t.Errorf("empty-filter receive emitted %q, want it omitted", key)
		}
	}
	assertMessageVocabulary(t, "empty-filter receive", msg)
}

// TestReceiveMessageEmissionSystemAttributesOptIn pins the system-attribute
// list on both wire forms: omitted returns no Attributes, a specific name
// returns exactly that attribute, and the query-wire stem form
// (AttributeName.N) drives the same opt-in rule.
func TestReceiveMessageEmissionSystemAttributesOptIn(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "emit-sys", nil)

	msg := receiveOne(t, svc, reqCtx, queueURL, nil, nil)
	if _, ok := msg["Attributes"]; ok {
		t.Error("default receive emitted Attributes, want it omitted")
	}

	msg = receiveOne(t, svc, reqCtx, queueURL, nil,
		map[string]interface{}{"AttributeNames": []interface{}{"SentTimestamp"}})
	sysAttrs := msg["Attributes"].(map[string]string)
	if len(sysAttrs) != 1 || sysAttrs["SentTimestamp"] == "" {
		t.Fatalf("SentTimestamp-filtered receive returned %v, want exactly SentTimestamp", sysAttrs)
	}

	msg = receiveOne(t, svc, reqCtx, queueURL, nil,
		map[string]interface{}{"AttributeName.1": "SentTimestamp"})
	sysAttrs = msg["Attributes"].(map[string]string)
	if len(sysAttrs) != 1 || sysAttrs["SentTimestamp"] == "" {
		t.Fatalf("query-stem filtered receive returned %v, want exactly SentTimestamp", sysAttrs)
	}
}

// TestReceiveMessageEmissionFifoIdsThroughAttributes pins the FIFO
// identifier channel: a FIFO receive exposes MessageGroupId,
// MessageDeduplicationId and SequenceNumber as system attributes when
// requested, and never as top-level Message members (the model's Message
// shape has no such members).
func TestReceiveMessageEmissionFifoIdsThroughAttributes(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "emit-fifo.fifo", map[string]interface{}{
		"FifoQueue":                 "true",
		"ContentBasedDeduplication": "true",
	})

	msg := receiveOne(t, svc, reqCtx, queueURL,
		map[string]interface{}{
			"MessageGroupId":         "emit-group",
			"MessageDeduplicationId": "emit-dedup",
		},
		map[string]interface{}{"AttributeNames": []interface{}{"All"}})

	for _, key := range []string{"MessageGroupId", "MessageDeduplicationId"} {
		if _, ok := msg[key]; ok {
			t.Errorf("received message emitted top-level %q; FIFO identifiers travel as system attributes", key)
		}
	}
	sysAttrs := msg["Attributes"].(map[string]string)
	if sysAttrs["MessageGroupId"] != "emit-group" {
		t.Errorf("Attributes[MessageGroupId] = %q, want emit-group", sysAttrs["MessageGroupId"])
	}
	if sysAttrs["MessageDeduplicationId"] != "emit-dedup" {
		t.Errorf("Attributes[MessageDeduplicationId] = %q, want emit-dedup", sysAttrs["MessageDeduplicationId"])
	}
	if sysAttrs["SequenceNumber"] == "" {
		t.Error("FIFO receive Attributes missing SequenceNumber under All")
	}
	assertMessageVocabulary(t, "FIFO receive", msg)
}

// assertMessageVocabulary checks a received message's emitted keys against
// the model's Message shape vocabulary (see response_vocabulary_test.go).
func assertMessageVocabulary(t *testing.T, context string, msg map[string]interface{}) {
	t.Helper()
	allowed := make(map[string]bool, len(modelOutputVocabulary["ReceiveMessage Message"]))
	for _, k := range modelOutputVocabulary["ReceiveMessage Message"] {
		allowed[k] = true
	}
	for k := range msg {
		if !allowed[k] {
			t.Errorf("%s emitted unmodelled Message key %q", context, k)
		}
	}
}

// TestReceiveMessageSystemAttributeNameVocabulary pins the filter-enum
// contract of the receive name lists: both system channels reject an
// out-of-vocabulary token (including the MessageAttributeNames-only ".*"
// wildcard), while documented names and "All" stay accepted.
func TestReceiveMessageSystemAttributeNameVocabulary(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "vocab-probe", nil)

	cases := []struct {
		name      string
		recvExtra map[string]interface{}
	}{
		{"deprecated channel unknown name", map[string]interface{}{
			"AttributeNames": []interface{}{"NotAnAttribute"},
		}},
		{"deprecated channel wildcard form", map[string]interface{}{
			"AttributeNames": []interface{}{".*"},
		}},
		{"system channel unknown name", map[string]interface{}{
			"MessageSystemAttributeNames": []interface{}{"NotAnAttribute"},
		}},
		{"system channel wildcard form", map[string]interface{}{
			"MessageSystemAttributeNames": []interface{}{".*"},
		}},
	}
	for _, tc := range cases {
		_, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "ReceiveMessage",
			Parameters: mergeRecvParams(queueURL, tc.recvExtra),
		})
		if err == nil || !strings.Contains(err.Error(), "InvalidAttributeName") {
			t.Fatalf("%s: expected InvalidAttributeName, got %v", tc.name, err)
		}
	}

	// Documented names and the wildcard forms of the free-string member stay
	// accepted (the queue is empty; the request must simply not error).
	for name, extra := range map[string]map[string]interface{}{
		"deprecated channel documented name": {"AttributeNames": []interface{}{"SentTimestamp"}},
		"deprecated channel All":             {"AttributeNames": []interface{}{"All"}},
		"system channel All":                 {"MessageSystemAttributeNames": []interface{}{"All"}},
		"message attributes wildcard":        {"MessageAttributeNames": []interface{}{".*"}},
	} {
		if _, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "ReceiveMessage",
			Parameters: mergeRecvParams(queueURL, extra),
		}); err != nil {
			t.Fatalf("%s: accepted, got %v", name, err)
		}
	}
}

// mergeRecvParams builds a receive parameter map over the queue URL.
func mergeRecvParams(queueURL string, extra map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{"QueueUrl": queueURL}
	for k, v := range extra {
		params[k] = v
	}
	return params
}

// TestReceiveMessageSqsManagedSseEnabledEmitted pins the documented
// receivable system attribute SqsManagedSseEnabled: a queue configured with
// the attribute reports its value per message when the deprecated
// AttributeNames channel asks for it.
func TestReceiveMessageSqsManagedSseEnabledEmitted(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)
	queueURL := createEmissionQueue(t, svc, reqCtx, "sse-probe", map[string]interface{}{
		"SqsManagedSseEnabled": "true",
	})

	msg := receiveOne(t, svc, reqCtx, queueURL, nil, map[string]interface{}{
		"AttributeNames": []interface{}{"SqsManagedSseEnabled"},
	})
	attrs, ok := msg["Attributes"].(map[string]string)
	if !ok || attrs["SqsManagedSseEnabled"] != "true" {
		t.Fatalf("expected Attributes[SqsManagedSseEnabled]=true, got %v", msg["Attributes"])
	}
}
