package sqs

import "testing"

func strPtr(s string) *string { return &s }

func TestValidateMessageAttributes(t *testing.T) {
	valid := map[string]*MessageAttributeValue{
		"attr_1":     {DataType: "String", StringValue: strPtr("value")},
		"attr-2":     {DataType: "Number", StringValue: strPtr("42")},
		"a.b":        {DataType: "Binary", BinaryValue: []byte{1, 2, 3}},
		"Number.int": {DataType: "Number.int", StringValue: strPtr("7")},
	}
	if err := ValidateMessageAttributes(valid); err != nil {
		t.Fatalf("ValidateMessageAttributes(valid) = %v, want nil", err)
	}

	if err := ValidateMessageAttributes(nil); err != nil {
		t.Fatalf("ValidateMessageAttributes(nil) = %v, want nil", err)
	}

	tooMany := make(map[string]*MessageAttributeValue, MaxMessageAttributes+1)
	for i := 0; i < MaxMessageAttributes+1; i++ {
		tooMany[string(rune('a'+i%26))+string(rune('0'+i/26))] = &MessageAttributeValue{DataType: "String"}
	}
	if err := ValidateMessageAttributes(tooMany); err == nil {
		t.Fatal("ValidateMessageAttributes(11 attrs) = nil, want error")
	}

	invalidNames := []string{
		"",
		"aws.reserved",
		"Amazon.Reserved",
		".leading",
		"trailing.",
		"double..period",
		"invalid char!",
	}
	for _, name := range invalidNames {
		attrs := map[string]*MessageAttributeValue{name: {DataType: "String"}}
		if err := ValidateMessageAttributes(attrs); err == nil {
			t.Errorf("ValidateMessageAttributes(name=%q) = nil, want error", name)
		}
	}

	oversized := string(make([]byte, MaxMessageAttributeNameLength+1))
	for i := range oversized {
		oversized = oversized[:i] + "a" + oversized[i+1:]
	}
	attrs := map[string]*MessageAttributeValue{oversized: {DataType: "String"}}
	if err := ValidateMessageAttributes(attrs); err == nil {
		t.Error("ValidateMessageAttributes(257-char name) = nil, want error")
	}

	invalidTypes := []string{
		"",
		"NotAType",
		string(make([]byte, 0)),
	}
	for _, dataType := range invalidTypes {
		if dataType == "" {
			// Empty DataType must be rejected outright.
			single := map[string]*MessageAttributeValue{"ok": {DataType: ""}}
			if err := ValidateMessageAttributes(single); err == nil {
				t.Error("ValidateMessageAttributes(empty DataType) = nil, want error")
			}
			continue
		}
		single := map[string]*MessageAttributeValue{"ok": {DataType: dataType}}
		if err := ValidateMessageAttributes(single); err == nil {
			t.Errorf("ValidateMessageAttributes(dataType=%q) = nil, want error", dataType)
		}
	}

	// A nil attribute value must be rejected.
	nilAttr := map[string]*MessageAttributeValue{"ok": nil}
	if err := ValidateMessageAttributes(nilAttr); err == nil {
		t.Error("ValidateMessageAttributes(nil attr) = nil, want error")
	}
}

func TestMessageSize(t *testing.T) {
	attrs := map[string]*MessageAttributeValue{
		"name": {DataType: "String", StringValue: strPtr("value")},
	}
	// body(5) + name(4) + dataType(6) + value(5) = 20
	if got := MessageSize("body!", attrs); got != 20 {
		t.Errorf("MessageSize = %d, want 20", got)
	}
	if got := MessageSize("abc", nil); got != 3 {
		t.Errorf("MessageSize(nil attrs) = %d, want 3", got)
	}
}
