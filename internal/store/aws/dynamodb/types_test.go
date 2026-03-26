package dynamodb

import (
	"testing"
)

func TestAttributeValue_IsString(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: true,
		},
		{
			name: "number value",
			av:   &AttributeValue{N: strPtr("123")},
			want: false,
		},
		{
			name: "binary value",
			av:   &AttributeValue{B: []byte("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsNumber(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "number value",
			av:   &AttributeValue{N: strPtr("123")},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsNumber(); got != tt.want {
				t.Errorf("IsNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsBinary(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "binary value",
			av:   &AttributeValue{B: []byte("test")},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsBinary(); got != tt.want {
				t.Errorf("IsBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsBool(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "bool true value",
			av:   &AttributeValue{BOOL: boolPtr(true)},
			want: true,
		},
		{
			name: "bool false value",
			av:   &AttributeValue{BOOL: boolPtr(false)},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsBool(); got != tt.want {
				t.Errorf("IsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsNull(t *testing.T) {
	trueVal := true
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "null value",
			av:   &AttributeValue{NULL: &trueVal},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsNull(); got != tt.want {
				t.Errorf("IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsMap(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "map value",
			av:   &AttributeValue{M: map[string]*AttributeValue{}},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsMap(); got != tt.want {
				t.Errorf("IsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsList(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "list value",
			av:   &AttributeValue{L: []*AttributeValue{}},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsList(); got != tt.want {
				t.Errorf("IsList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsStringSet(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "string set value",
			av:   &AttributeValue{SS: []string{"a", "b"}},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsStringSet(); got != tt.want {
				t.Errorf("IsStringSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsNumberSet(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "number set value",
			av:   &AttributeValue{NS: []string{"1", "2"}},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsNumberSet(); got != tt.want {
				t.Errorf("IsNumberSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValue_IsBinarySet(t *testing.T) {
	tests := []struct {
		name string
		av   *AttributeValue
		want bool
	}{
		{
			name: "binary set value",
			av:   &AttributeValue{BS: [][]byte{[]byte("a"), []byte("b")}},
			want: true,
		},
		{
			name: "string value",
			av:   &AttributeValue{S: strPtr("test")},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.av.IsBinarySet(); got != tt.want {
				t.Errorf("IsBinarySet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	got := StringValue("test")
	if got.S == nil || *got.S != "test" {
		t.Errorf("StringValue() = %v, want S pointer to 'test'", got)
	}
}

func TestNumberValue(t *testing.T) {
	got := NumberValue("123")
	if got.N == nil || *got.N != "123" {
		t.Errorf("NumberValue() = %v, want N pointer to '123'", got)
	}
}

func TestBinaryValue(t *testing.T) {
	data := []byte("test")
	got := BinaryValue(data)
	if got.B == nil || string(got.B) != "test" {
		t.Errorf("BinaryValue() = %v, want B to be 'test'", got)
	}
}

func TestBoolValue(t *testing.T) {
	got := BoolValue(true)
	if got.BOOL == nil || *got.BOOL != true {
		t.Errorf("BoolValue() = %v, want BOOL to be true", got)
	}
}

func TestNullValue(t *testing.T) {
	got := NullValue()
	if got.NULL == nil || *got.NULL != true {
		t.Errorf("NullValue() = %v, want NULL to be true", got)
	}
}

func TestMapValue(t *testing.T) {
	m := map[string]*AttributeValue{
		"key": StringValue("value"),
	}
	got := MapValue(m)
	if got.M == nil {
		t.Errorf("MapValue() = %v, want M to be non-nil", got)
	}
	if got.M["key"] == nil || *got.M["key"].S != "value" {
		t.Errorf("MapValue() key value = %v, want 'value'", got.M["key"])
	}
}

func TestListValue(t *testing.T) {
	l := []*AttributeValue{StringValue("a"), StringValue("b")}
	got := ListValue(l)
	if got.L == nil || len(got.L) != 2 {
		t.Errorf("ListValue() = %v, want L to have 2 elements", got)
	}
}

func TestStringSet(t *testing.T) {
	ss := []string{"a", "b", "c"}
	got := StringSet(ss)
	if got.SS == nil || len(got.SS) != 3 {
		t.Errorf("StringSet() = %v, want SS to have 3 elements", got)
	}
}

func TestNumberSet(t *testing.T) {
	ns := []string{"1", "2", "3"}
	got := NumberSet(ns)
	if got.NS == nil || len(got.NS) != 3 {
		t.Errorf("NumberSet() = %v, want NS to have 3 elements", got)
	}
}

func TestBinarySet(t *testing.T) {
	bs := [][]byte{[]byte("a"), []byte("b")}
	got := BinarySet(bs)
	if got.BS == nil || len(got.BS) != 2 {
		t.Errorf("BinarySet() = %v, want BS to have 2 elements", got)
	}
}

func TestItem_GetKeyAttributeValue(t *testing.T) {
	item := &Item{
		Key: map[string]*AttributeValue{
			"pk": StringValue("partition-key"),
			"sk": StringValue("sort-key"),
		},
	}

	tests := []struct {
		name      string
		attrName  string
		wantFound bool
	}{
		{
			name:      "existing key",
			attrName:  "pk",
			wantFound: true,
		},
		{
			name:      "existing sort key",
			attrName:  "sk",
			wantFound: true,
		},
		{
			name:      "non-existing key",
			attrName:  "unknown",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := item.GetKeyAttributeValue(tt.attrName)
			if tt.wantFound && got == nil {
				t.Errorf("GetKeyAttributeValue() = nil, want non-nil for %s", tt.attrName)
			}
			if !tt.wantFound && got != nil {
				t.Errorf("GetKeyAttributeValue() = %v, want nil for %s", got, tt.attrName)
			}
		})
	}
}

func TestItem_GetAttribute(t *testing.T) {
	item := &Item{
		Attributes: map[string]*AttributeValue{
			"name": StringValue("test-name"),
			"age":  NumberValue("30"),
		},
	}

	tests := []struct {
		name      string
		attrName  string
		wantFound bool
	}{
		{
			name:      "existing attribute",
			attrName:  "name",
			wantFound: true,
		},
		{
			name:      "non-existing attribute",
			attrName:  "unknown",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := item.GetAttribute(tt.attrName)
			if tt.wantFound && got == nil {
				t.Errorf("GetAttribute() = nil, want non-nil for %s", tt.attrName)
			}
			if !tt.wantFound && got != nil {
				t.Errorf("GetAttribute() = %v, want nil for %s", got, tt.attrName)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
