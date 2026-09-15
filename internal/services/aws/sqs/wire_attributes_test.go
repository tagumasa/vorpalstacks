package sqs

import (
	"fmt"
	"testing"

	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// jsonAttributeParams builds a SendMessage parameter map whose
// MessageAttributes arrive as the JSON-protocol map — the form the AWS SDK
// sends.
func jsonAttributeParams() map[string]interface{} {
	return map[string]interface{}{
		"QueueUrl":    "http://localhost:50080/000000000000/queue-a",
		"MessageBody": "body",
		"MessageAttributes": map[string]interface{}{
			"Str": map[string]interface{}{
				"DataType":    "String",
				"StringValue": "value1",
			},
			"Num": map[string]interface{}{
				"DataType":    "Number",
				"StringValue": "42",
			},
			"Bin": map[string]interface{}{
				"DataType":    "Binary",
				"BinaryValue": "dmFsdWU=", // []byte("value")
			},
		},
	}
}

// queryAttributeParams builds the equivalent request with the flattened
// query keys the CLI / query protocol sends.
func queryAttributeParams() map[string]interface{} {
	return map[string]interface{}{
		"QueueUrl":                             "http://localhost:50080/000000000000/queue-a",
		"MessageBody":                          "body",
		"MessageAttribute.1.Name":              "Str",
		"MessageAttribute.1.Value.DataType":    "String",
		"MessageAttribute.1.Value.StringValue": "value1",
		"MessageAttribute.2.Name":              "Num",
		"MessageAttribute.2.Value.DataType":    "Number",
		"MessageAttribute.2.Value.StringValue": "42",
		"MessageAttribute.3.Name":              "Bin",
		"MessageAttribute.3.Value.DataType":    "Binary",
		"MessageAttribute.3.Value.BinaryValue": "dmFsdWU=",
	}
}

// assertAttributeSetsEqual compares two parsed attribute maps member by
// member.
func assertAttributeSetsEqual(t *testing.T, got, want map[string]*sqsstore.MessageAttributeValue) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("parsed %d attributes, want %d", len(got), len(want))
	}
	for name, wantAttr := range want {
		gotAttr, ok := got[name]
		if !ok {
			t.Fatalf("attribute %q missing", name)
		}
		if gotAttr.DataType != wantAttr.DataType {
			t.Errorf("attribute %q DataType = %q, want %q", name, gotAttr.DataType, wantAttr.DataType)
		}
		if (gotAttr.StringValue == nil) != (wantAttr.StringValue == nil) {
			t.Fatalf("attribute %q StringValue presence differs: got %v", name, gotAttr.StringValue)
		}
		if gotAttr.StringValue != nil && *gotAttr.StringValue != *wantAttr.StringValue {
			t.Errorf("attribute %q StringValue = %q, want %q", name, *gotAttr.StringValue, *wantAttr.StringValue)
		}
		if (gotAttr.BinaryValue == nil) != (wantAttr.BinaryValue == nil) {
			t.Fatalf("attribute %q BinaryValue presence differs: got %v", name, gotAttr.BinaryValue)
		}
		if gotAttr.BinaryValue != nil && string(gotAttr.BinaryValue) != string(wantAttr.BinaryValue) {
			t.Errorf("attribute %q BinaryValue = %q, want %q", name, gotAttr.BinaryValue, wantAttr.BinaryValue)
		}
	}
}

// TestSendMessageAttributeParitySingleVsBatch pins the single-parser rule:
// the same logical attribute set parses to the same store representation and
// the same MD5OfMessageAttributes basis on every path — single send and
// batch, JSON and flattened query. The five former parser copies drifted
// exactly here (empty StringValue kept on one path and dropped on another,
// BinaryValue dropped on the query arms).
func TestSendMessageAttributeParitySingleVsBatch(t *testing.T) {
	single, err := parseRequestMessageAttributes(jsonAttributeParams())
	if err != nil {
		t.Fatalf("single-send JSON parse failed: %v", err)
	}

	query, err := parseRequestMessageAttributes(queryAttributeParams())
	if err != nil {
		t.Fatalf("single-send query parse failed: %v", err)
	}

	batchJSON, err := parseBatchSendEntries(map[string]interface{}{
		"QueueUrl": "http://localhost:50080/000000000000/queue-a",
		"Entries": []interface{}{
			map[string]interface{}{
				"Id":                "entry-1",
				"MessageBody":       "body",
				"MessageAttributes": jsonAttributeParams()["MessageAttributes"],
			},
		},
	})
	if err != nil {
		t.Fatalf("batch JSON parse failed: %v", err)
	}

	queryParams := queryAttributeParams()
	batchQueryParams := map[string]interface{}{
		"QueueUrl":                                                            queryParams["QueueUrl"],
		"SendMessageBatchRequestEntry.1.Id":                                   "entry-1",
		"SendMessageBatchRequestEntry.1.MessageBody":                          "body",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Name":              "Str",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Value.DataType":    "String",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Value.StringValue": "value1",
		"SendMessageBatchRequestEntry.1.MessageAttribute.2.Name":              "Num",
		"SendMessageBatchRequestEntry.1.MessageAttribute.2.Value.DataType":    "Number",
		"SendMessageBatchRequestEntry.1.MessageAttribute.2.Value.StringValue": "42",
		"SendMessageBatchRequestEntry.1.MessageAttribute.3.Name":              "Bin",
		"SendMessageBatchRequestEntry.1.MessageAttribute.3.Value.DataType":    "Binary",
		"SendMessageBatchRequestEntry.1.MessageAttribute.3.Value.BinaryValue": "dmFsdWU=",
	}
	batchQuery, err := parseBatchSendEntries(batchQueryParams)
	if err != nil {
		t.Fatalf("batch query parse failed: %v", err)
	}

	// The query arms must decode BinaryValue, not accept-and-drop it.
	assertAttributeSetsEqual(t, query, single)
	assertAttributeSetsEqual(t, batchJSON[0].message.MessageAttributes, single)
	assertAttributeSetsEqual(t, batchQuery[0].message.MessageAttributes, single)

	wantMD5 := sqsstore.CalculateMessageAttributesMD5(single)
	for name, attrs := range map[string]map[string]*sqsstore.MessageAttributeValue{
		"single JSON":  single,
		"single query": query,
		"batch JSON":   batchJSON[0].message.MessageAttributes,
		"batch query":  batchQuery[0].message.MessageAttributes,
	} {
		if got := sqsstore.CalculateMessageAttributesMD5(attrs); got != wantMD5 {
			t.Errorf("%s: MD5 basis %s, want %s", name, got, wantMD5)
		}
	}
}

// TestEmptyAttributeValueRejected pins the empty-value rule from the
// MessageAttributeValue documentation trait ("Name, type, value and the
// message body must not be empty or null"): an attribute carrying no
// non-empty StringValue and no non-empty BinaryValue is rejected on every
// wire path, never silently stored as a value-less attribute.
func TestEmptyAttributeValueRejected(t *testing.T) {
	cases := []struct {
		name string
		run  func() error
	}{
		{"single JSON empty StringValue", func() error {
			_, err := parseRequestMessageAttributes(map[string]interface{}{
				"MessageAttributes": map[string]interface{}{
					"Str": map[string]interface{}{"DataType": "String", "StringValue": ""},
				},
			})
			return err
		}},
		{"single query value-less attribute", func() error {
			_, err := parseRequestMessageAttributes(map[string]interface{}{
				"MessageAttribute.1.Name":           "Str",
				"MessageAttribute.1.Value.DataType": "String",
			})
			return err
		}},
		{"batch JSON empty StringValue", func() error {
			_, err := parseBatchSendEntries(map[string]interface{}{
				"Entries": []interface{}{
					map[string]interface{}{
						"Id":          "entry-1",
						"MessageBody": "body",
						"MessageAttributes": map[string]interface{}{
							"Str": map[string]interface{}{"DataType": "String", "StringValue": ""},
						},
					},
				},
			})
			return err
		}},
		{"system attribute without value", func() error {
			_, err := parseSystemAttributes(map[string]interface{}{
				"MessageSystemAttribute.1.Name":           "AWSTraceHeader",
				"MessageSystemAttribute.1.Value.DataType": "String",
			}, "MessageSystemAttribute.")
			return err
		}},
	}
	for _, tc := range cases {
		if err := tc.run(); err != ErrInvalidParameterValue {
			t.Errorf("%s: got %v, want ErrInvalidParameterValue", tc.name, err)
		}
	}
}

// TestQueryBatchMalformedAttributeRejectedNotTruncated pins the
// malformed-attribute policy on the query batch arm: an attribute with a
// Name but no value members is rejected for the whole request — the old
// parser `break`ed at the first missing DataType and silently truncated
// every remaining attribute of the entry.
func TestQueryBatchMalformedAttributeRejectedNotTruncated(t *testing.T) {
	params := map[string]interface{}{
		"QueueUrl":                                                            "http://localhost:50080/000000000000/queue-a",
		"SendMessageBatchRequestEntry.1.Id":                                   "entry-1",
		"SendMessageBatchRequestEntry.1.MessageBody":                          "body",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Name":              "Attr1",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Value.DataType":    "String",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Value.StringValue": "value1",
		"SendMessageBatchRequestEntry.1.MessageAttribute.2.Name":              "Attr2",
		"SendMessageBatchRequestEntry.1.MessageAttribute.2.Value.DataType":    "String",
	}
	_, err := parseBatchSendEntries(params)
	if err != ErrInvalidParameterValue {
		t.Fatalf("value-less second attribute must be rejected, got %v", err)
	}
}

// TestBatchEmptyMessageBodyRejected pins the @required MessageBody member of
// SendMessageBatchRequestEntry on both wire formats: an empty body is
// rejected exactly as the single-send path rejects it (the SDK cannot reach
// this path — its client-side validation demands a body — but the CLI and
// raw query forms can).
func TestBatchEmptyMessageBodyRejected(t *testing.T) {
	jsonErr := func() error {
		_, err := parseBatchSendEntries(map[string]interface{}{
			"Entries": []interface{}{
				map[string]interface{}{"Id": "entry-1", "MessageBody": ""},
			},
		})
		return err
	}()
	if jsonErr != ErrMissingParameter {
		t.Errorf("JSON empty body: got %v, want ErrMissingParameter", jsonErr)
	}

	queryErr := func() error {
		_, err := parseBatchSendEntries(map[string]interface{}{
			"SendMessageBatchRequestEntry.1.Id": "entry-1",
		})
		return err
	}()
	if queryErr != ErrMissingParameter {
		t.Errorf("query empty body: got %v, want ErrMissingParameter", queryErr)
	}
}

// TestValidateBatchEntryIDsOrder pins the shared batch-entry rule order:
// the entry count is checked first, then per-Id format, then distinctness.
func TestValidateBatchEntryIDsOrder(t *testing.T) {
	eleven := make([]string, 11)
	for i := range eleven {
		eleven[i] = fmt.Sprintf("entry-%d", i)
	}
	if err := validateBatchEntryIDs(eleven); err != ErrTooManyEntriesInBatch {
		t.Errorf("eleven entries: got %v, want ErrTooManyEntriesInBatch", err)
	}
	eleven[2] = "bad id!"
	if err := validateBatchEntryIDs(eleven); err != ErrTooManyEntriesInBatch {
		t.Errorf("count must outrank a malformed Id: got %v, want ErrTooManyEntriesInBatch", err)
	}
	if err := validateBatchEntryIDs([]string{"entry-1", "entry-1"}); err != ErrBatchEntryIdsNotDistinct {
		t.Errorf("duplicate Ids: got %v, want ErrBatchEntryIdsNotDistinct", err)
	}
	if err := validateBatchEntryIDs([]string{"bad id!"}); err != ErrInvalidBatchEntryId {
		t.Errorf("malformed Id: got %v, want ErrInvalidBatchEntryId", err)
	}
}

// TestSystemAttributesQueryArm pins the consolidated system-attribute
// parsing: the query arm resolves AWSTraceHeader with its value and rejects
// any other system attribute name.
func TestSystemAttributesQueryArm(t *testing.T) {
	attrs, err := parseSystemAttributes(map[string]interface{}{
		"MessageSystemAttribute.1.Name":              "AWSTraceHeader",
		"MessageSystemAttribute.1.Value.DataType":    "String",
		"MessageSystemAttribute.1.Value.StringValue": "trace",
	}, "MessageSystemAttribute.")
	if err != nil {
		t.Fatalf("valid AWSTraceHeader must parse: %v", err)
	}
	if th := attrs["AWSTraceHeader"]; th == nil || th.StringValue == nil || *th.StringValue != "trace" {
		t.Fatalf("AWSTraceHeader not parsed: %+v", attrs)
	}

	_, err = parseSystemAttributes(map[string]interface{}{
		"MessageSystemAttribute.1.Name":              "Other",
		"MessageSystemAttribute.1.Value.DataType":    "String",
		"MessageSystemAttribute.1.Value.StringValue": "value",
	}, "MessageSystemAttribute.")
	if err != ErrInvalidParameterValue {
		t.Errorf("unknown system attribute: got %v, want ErrInvalidParameterValue", err)
	}
}

// TestJSONNonObjectAttributeValueRejected pins the malformed-shape policy
// on the JSON attribute arm: a message-attribute value that is not an
// object violates the MessageAttributeValue shape and rejects the request —
// never a silent drop of the attribute from the map.
func TestJSONNonObjectAttributeValueRejected(t *testing.T) {
	scalar := map[string]interface{}{
		"MessageAttributes": map[string]interface{}{
			"foo": "bar",
		},
	}
	if _, err := parseRequestMessageAttributes(scalar); err != ErrSerializationException {
		t.Fatalf("scalar attribute value: got %v, want ErrSerializationException", err)
	}

	nested := map[string]interface{}{
		"Entries": []interface{}{
			map[string]interface{}{
				"Id":                "a",
				"MessageBody":       "body",
				"MessageAttributes": map[string]interface{}{"foo": 42},
			},
		},
	}
	if _, err := parseBatchEntriesJSON(nested["Entries"].([]interface{})); err != ErrSerializationException {
		t.Fatalf("numeric attribute value in a batch entry: got %v, want ErrSerializationException", err)
	}
}
