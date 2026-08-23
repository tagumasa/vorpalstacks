package sqs

import (
	"fmt"
	"testing"
)

func TestParseBatchEntriesQueryRejectsEleventhEntry(t *testing.T) {
	// Query-protocol batch entries are contiguous and the entry-count check
	// must reject an eleventh entry rather than silently ignoring it.
	params := map[string]interface{}{}
	for i := 1; i <= 11; i++ {
		params[fmt.Sprintf("SendMessageBatchRequestEntry.%d.Id", i)] = fmt.Sprintf("entry-%d", i)
		params[fmt.Sprintf("SendMessageBatchRequestEntry.%d.MessageBody", i)] = fmt.Sprintf("body-%d", i)
	}
	_, err := parseBatchEntriesQuery(params)
	if err == nil {
		t.Fatalf("11 query-format entries must fail with TooManyEntriesInBatchRequest")
	}
	if err != ErrTooManyEntriesInBatch {
		t.Errorf("got %v, want ErrTooManyEntriesInBatch", err)
	}
}

func TestParseBatchEntriesQueryTenEntries(t *testing.T) {
	params := map[string]interface{}{}
	for i := 1; i <= 10; i++ {
		params[fmt.Sprintf("SendMessageBatchRequestEntry.%d.Id", i)] = fmt.Sprintf("entry-%d", i)
		params[fmt.Sprintf("SendMessageBatchRequestEntry.%d.MessageBody", i)] = fmt.Sprintf("body-%d", i)
	}
	entries, err := parseBatchEntriesQuery(params)
	if err != nil {
		t.Fatalf("10 query-format entries must parse: %v", err)
	}
	if len(entries) != 10 {
		t.Errorf("parsed %d entries, want 10", len(entries))
	}
}

func TestRejectAttributeListValuesJSON(t *testing.T) {
	params := map[string]interface{}{
		"MessageAttributes": map[string]interface{}{
			"Attr1": map[string]interface{}{
				"DataType":         "String",
				"StringListValues": []interface{}{"value1"},
			},
		},
	}
	if err := rejectAttributeListValues(params); err != ErrUnsupportedOperation {
		t.Errorf("JSON StringListValues: got %v, want ErrUnsupportedOperation", err)
	}

	params = map[string]interface{}{
		"MessageAttributes": map[string]interface{}{
			"Attr1": map[string]interface{}{
				"DataType":         "Binary",
				"BinaryListValues": []interface{}{"dmFsdWU="},
			},
		},
	}
	if err := rejectAttributeListValues(params); err != ErrUnsupportedOperation {
		t.Errorf("JSON BinaryListValues: got %v, want ErrUnsupportedOperation", err)
	}

	params = map[string]interface{}{
		"MessageSystemAttributes": map[string]interface{}{
			"AWSTraceHeader": map[string]interface{}{
				"DataType":         "String",
				"StringValue":      "trace",
				"StringListValues": []interface{}{"value1"},
			},
		},
	}
	if err := rejectAttributeListValues(params); err != ErrUnsupportedOperation {
		t.Errorf("JSON system StringListValues: got %v, want ErrUnsupportedOperation", err)
	}

	params = map[string]interface{}{
		"Entries": []interface{}{
			map[string]interface{}{
				"Id":          "entry-1",
				"MessageBody": "body",
				"MessageAttributes": map[string]interface{}{
					"Attr1": map[string]interface{}{
						"DataType":         "String",
						"StringListValues": []interface{}{"value1"},
					},
				},
			},
		},
	}
	if err := rejectAttributeListValues(params); err != ErrUnsupportedOperation {
		t.Errorf("JSON batch entry StringListValues: got %v, want ErrUnsupportedOperation", err)
	}
}

func TestRejectAttributeListValuesQuery(t *testing.T) {
	params := map[string]interface{}{
		"MessageAttribute.1.Name":                    "Attr1",
		"MessageAttribute.1.Value.DataType":          "String",
		"MessageAttribute.1.Value.StringListValue.1": "value1",
	}
	if err := rejectAttributeListValues(params); err != ErrUnsupportedOperation {
		t.Errorf("query StringListValue: got %v, want ErrUnsupportedOperation", err)
	}

	params = map[string]interface{}{
		"SendMessageBatchRequestEntry.1.Id":                                         "entry-1",
		"SendMessageBatchRequestEntry.1.MessageAttribute.1.Value.BinaryListValue.1": "dmFsdWU=",
	}
	if err := rejectAttributeListValues(params); err != ErrUnsupportedOperation {
		t.Errorf("query batch BinaryListValue: got %v, want ErrUnsupportedOperation", err)
	}
}

func TestRejectAttributeListValuesCleanRequest(t *testing.T) {
	params := map[string]interface{}{
		"MessageBody": "body",
		"MessageAttributes": map[string]interface{}{
			"Attr1": map[string]interface{}{
				"DataType":    "String",
				"StringValue": "value1",
			},
		},
	}
	if err := rejectAttributeListValues(params); err != nil {
		t.Errorf("ordinary attributes must pass, got %v", err)
	}
}
