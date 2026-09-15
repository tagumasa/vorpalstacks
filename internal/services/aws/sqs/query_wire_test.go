package sqs

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// newQueryWireTestService builds the service over a real store on a fresh
// temp directory so query-wire pins run the real handler→Core→store path —
// the generic helpers the query wire bypasses make SDK (JSON-wire) tests
// structurally unable to cover these forms.
func newQueryWireTestService(t *testing.T) (*SQSService, *request.RequestContext, sqsstore.SQSStoreInterface) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	store := sqsstore.NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	// Storage closes after the store: LIFO runs store.Close first, so the
	// move workers Close drains run against live storage.
	t.Cleanup(func() { _ = st.Close() })
	t.Cleanup(func() { store.Close() })

	svc := &SQSService{accountID: "123456789012"}
	svc.stores.Store("us-east-1", store)
	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	return svc, reqCtx, store
}

// TestCreateQueueQueryWireTags pins the Tag.N.Key / Tag.N.Value form the
// tags member carries on the query wire (xmlName "Tag" + xmlFlattened):
// before the stem fallback the queue was silently created with zero tags.
// The JSON member map keeps flowing through the same handler.
func TestCreateQueueQueryWireTags(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":   "query-tags",
			"Tag.1.Key":   "team",
			"Tag.1.Value": "core",
			"Tag.2.Key":   "env",
			"Tag.2.Value": "prod",
		},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	tags, err := store.ListQueueTags(queueURL)
	if err != nil {
		t.Fatalf("ListQueueTags: %v", err)
	}
	if len(tags) != 2 || tags["team"] != "core" || tags["env"] != "prod" {
		t.Fatalf("query-wire tags = %v, want team=core env=prod", tags)
	}

	jsonResp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName": "query-tags-json",
			"Tags":      map[string]interface{}{"stage": "test"},
		},
	})
	if err != nil {
		t.Fatalf("CreateQueue (JSON arm): %v", err)
	}
	jsonMap, _ := jsonResp.(map[string]interface{})
	jsonURL, _ := jsonMap["QueueUrl"].(string)
	jsonTags, err := store.ListQueueTags(jsonURL)
	if err != nil {
		t.Fatalf("ListQueueTags (JSON arm): %v", err)
	}
	if len(jsonTags) != 1 || jsonTags["stage"] != "test" {
		t.Fatalf("JSON-arm tags = %v, want stage=test", jsonTags)
	}
}

// TestUntagQueueQueryWireTagKeys pins the TagKey.N form the TagKeys member
// carries on the query wire (xmlName "TagKey" + xmlFlattened): before the
// parse override the handler parsed zero keys and returned success without
// removing any tag.
func TestUntagQueueQueryWireTagKeys(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "query-untag"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	if err := store.TagQueue(queueURL, map[string]string{"team": "core", "env": "prod", "keep": "yes"}); err != nil {
		t.Fatalf("TagQueue: %v", err)
	}

	if _, err := svc.UntagQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "UntagQueue",
		Parameters: map[string]interface{}{
			"QueueUrl": queueURL,
			"TagKey.1": "team",
			"TagKey.2": "env",
		},
	}); err != nil {
		t.Fatalf("UntagQueue: %v", err)
	}

	tags, err := store.ListQueueTags(queueURL)
	if err != nil {
		t.Fatalf("ListQueueTags: %v", err)
	}
	if len(tags) != 1 || tags["keep"] != "yes" {
		t.Fatalf("tags after query-wire untag = %v, want keep=yes only", tags)
	}
}

// TestReceiveMessageQueryWireFilters pins the AttributeName.N and
// MessageAttributeName.N filter forms (the members' xmlName stems): the
// requested filter must arrive at the emission filter. Restricting to a
// proper subset proves arrival — an unread filter leaves the list empty and
// the default-All emission returns everything stored.
func TestReceiveMessageQueryWireFilters(t *testing.T) {
	svc, reqCtx, store := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "query-filters"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	attr1 := "value-one"
	attr2 := "value-two"
	if _, err := store.SendMessage(queueURL, &sqsstore.Message{
		Body: "filtered-body",
		MessageAttributes: map[string]*sqsstore.MessageAttributeValue{
			"attr1": {DataType: "String", StringValue: &attr1},
			"attr2": {DataType: "String", StringValue: &attr2},
		},
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	receive := func(params map[string]interface{}) map[string]interface{} {
		t.Helper()
		params["QueueUrl"] = queueURL
		params["MaxNumberOfMessages"] = "1"
		// VisibilityTimeout 0 returns the message to the queue at once so
		// the second receive observes the same message.
		params["VisibilityTimeout"] = "0"
		r, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "ReceiveMessage",
			Parameters: params,
		})
		if err != nil {
			t.Fatalf("ReceiveMessage: %v", err)
		}
		rMap, _ := r.(map[string]interface{})
		messages, ok := rMap["Messages"].([]map[string]interface{})
		if !ok || len(messages) != 1 {
			t.Fatalf("receive returned %v messages, want 1", r)
		}
		return messages[0]
	}

	// System-attribute stem: restrict to SentTimestamp — SenderId and
	// ApproximateReceiveCount, also present on the message, are dropped.
	sysMsg := receive(map[string]interface{}{"AttributeName.1": "SentTimestamp"})
	sysAttrs, ok := sysMsg["Attributes"].(map[string]string)
	if !ok || len(sysAttrs) != 1 {
		t.Fatalf("Attributes = %v, want exactly the SentTimestamp entry", sysMsg["Attributes"])
	}
	if _, ok := sysAttrs["SentTimestamp"]; !ok {
		t.Fatalf("Attributes = %v, want SentTimestamp present", sysAttrs)
	}

	// Message-attribute stem: restrict to attr1 — attr2 is dropped.
	msgMsg := receive(map[string]interface{}{"MessageAttributeName.1": "attr1"})
	msgAttrs, ok := msgMsg["MessageAttributes"].(map[string]interface{})
	if !ok || len(msgAttrs) != 1 {
		t.Fatalf("MessageAttributes = %v, want exactly the attr1 entry", msgMsg["MessageAttributes"])
	}
	if _, ok := msgAttrs["attr1"]; !ok {
		t.Fatalf("MessageAttributes = %v, want attr1 present", msgAttrs)
	}
}

// TestQueryWireGappedStemListsReject pins the malformed-list rule: a stem
// list with a missing interior slot (AttributeName.1 then AttributeName.3)
// rejects the request instead of silently truncating the tail — the same
// rule the batch-entry parsers apply to entry-shaped keys without their Id.
func TestQueryWireGappedStemListsReject(t *testing.T) {
	svc, reqCtx, _ := newQueryWireTestService(t)

	resp, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation:  "CreateQueue",
		Parameters: map[string]interface{}{"QueueName": "gapped-stems"},
	})
	if err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	respMap, _ := resp.(map[string]interface{})
	queueURL, _ := respMap["QueueUrl"].(string)

	// GetQueueAttributes: gap in AttributeName.N.
	if _, err := svc.GetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "GetQueueAttributes",
		Parameters: map[string]interface{}{
			"QueueUrl":        queueURL,
			"AttributeName.1": "All",
			"AttributeName.3": "VisibilityTimeout",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("gapped AttributeName list: expected InvalidParameterValue, got %v", err)
	}

	// ReceiveMessage: gap in MessageAttributeName.N.
	if _, err := svc.ReceiveMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "ReceiveMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":               queueURL,
			"MessageAttributeName.1": "attr1",
			"MessageAttributeName.3": "attr2",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("gapped MessageAttributeName list: expected InvalidParameterValue, got %v", err)
	}

	// SendMessage: an attribute slot with value members but no Name.
	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":                             queueURL,
			"MessageBody":                          "body",
			"MessageAttribute.1.Name":              "attr1",
			"MessageAttribute.1.Value.DataType":    "String",
			"MessageAttribute.1.Value.StringValue": "value-one",
			"MessageAttribute.2.Value.DataType":    "String",
			"MessageAttribute.2.Value.StringValue": "value-two",
			"MessageAttribute.3.Name":              "attr3",
			"MessageAttribute.3.Value.DataType":    "String",
			"MessageAttribute.3.Value.StringValue": "value-three",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("attribute slot without a Name: expected InvalidParameterValue, got %v", err)
	}

	// SendMessage: a gapped attribute index — entry 1, no entry 2 at all,
	// entry 3 — must reject rather than parse entry 1 and silently drop
	// entry 3.
	if _, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":                             queueURL,
			"MessageBody":                          "body",
			"MessageAttribute.1.Name":              "attr1",
			"MessageAttribute.1.Value.DataType":    "String",
			"MessageAttribute.1.Value.StringValue": "value-one",
			"MessageAttribute.3.Name":              "attr3",
			"MessageAttribute.3.Value.DataType":    "String",
			"MessageAttribute.3.Value.StringValue": "value-three",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("gapped attribute index: expected InvalidParameterValue, got %v", err)
	}

	// CreateQueue: a gapped Attribute.N entry list (1 then 3) — the
	// map-shaped stems follow the flattened lists' contiguity rule.
	if _, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":         "gapped-attrs",
			"Attribute.1.Name":  "VisibilityTimeout",
			"Attribute.1.Value": "30",
			"Attribute.3.Name":  "MaximumMessageSize",
			"Attribute.3.Value": "1024",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("gapped Attribute entry list: expected InvalidParameterValue, got %v", err)
	}

	// CreateQueue: a Tag entry slot whose Value arrived without its Key.
	if _, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":   "keyless-tag",
			"Tag.1.Key":   "team",
			"Tag.1.Value": "core",
			"Tag.2.Value": "orphan",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("tag slot without a Key: expected InvalidParameterValue, got %v", err)
	}

	// CreateQueue: a gapped Tag.N entry list (1 then 3).
	if _, err := svc.CreateQueue(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "CreateQueue",
		Parameters: map[string]interface{}{
			"QueueName":   "gapped-tags",
			"Tag.1.Key":   "team",
			"Tag.1.Value": "core",
			"Tag.3.Key":   "env",
			"Tag.3.Value": "prod",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("gapped Tag entry list: expected InvalidParameterValue, got %v", err)
	}

	// SetQueueAttributes: an attribute slot with a Value but no Name.
	if _, err := svc.SetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SetQueueAttributes",
		Parameters: map[string]interface{}{
			"QueueUrl":          queueURL,
			"Attribute.1.Value": "30",
		},
	}); err == nil || !strings.Contains(err.Error(), "InvalidParameterValue") {
		t.Fatalf("attribute slot without a Name on SetQueueAttributes: expected InvalidParameterValue, got %v", err)
	}

	// Contiguous lists still parse: three well-formed attributes all land.
	sendResp, err := svc.SendMessage(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SendMessage",
		Parameters: map[string]interface{}{
			"QueueUrl":                             queueURL,
			"MessageBody":                          "body",
			"MessageAttribute.1.Name":              "attr1",
			"MessageAttribute.1.Value.DataType":    "String",
			"MessageAttribute.1.Value.StringValue": "value-one",
			"MessageAttribute.2.Name":              "attr2",
			"MessageAttribute.2.Value.DataType":    "String",
			"MessageAttribute.2.Value.StringValue": "value-two",
			"MessageAttribute.3.Name":              "attr3",
			"MessageAttribute.3.Value.DataType":    "String",
			"MessageAttribute.3.Value.StringValue": "value-three",
		},
	})
	if err != nil {
		t.Fatalf("contiguous attribute list rejected: %v", err)
	}
	if _, ok := sendResp.(map[string]interface{})["MD5OfMessageAttributes"]; !ok {
		t.Fatalf("contiguous attribute list did not round-trip attributes: %v", sendResp)
	}

	// Contiguous map entries still parse: SetQueueAttributes applies both
	// entries and the change reads back.
	if _, err := svc.SetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "SetQueueAttributes",
		Parameters: map[string]interface{}{
			"QueueUrl":          queueURL,
			"Attribute.1.Name":  "VisibilityTimeout",
			"Attribute.1.Value": "30",
			"Attribute.2.Name":  "MaximumMessageSize",
			"Attribute.2.Value": "2048",
		},
	}); err != nil {
		t.Fatalf("contiguous Attribute entry list rejected: %v", err)
	}
	gqaResp, err := svc.GetQueueAttributes(context.Background(), reqCtx, &request.ParsedRequest{
		Operation: "GetQueueAttributes",
		Parameters: map[string]interface{}{
			"QueueUrl":        queueURL,
			"AttributeName.1": "VisibilityTimeout",
		},
	})
	if err != nil {
		t.Fatalf("read back applied attributes: %v", err)
	}
	applied := gqaResp.(map[string]interface{})["Attributes"].(map[string]string)
	if applied["VisibilityTimeout"] != "30" {
		t.Fatalf("contiguous Attribute entry list did not apply: %v", applied)
	}
}
