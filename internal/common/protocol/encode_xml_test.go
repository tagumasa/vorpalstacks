package protocol

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeXMLResponse_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	err := EncodeXMLResponse(w, "Test", map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	assert.Contains(t, w.Body.String(), "<Result></Result>")
}

func TestEncodeXMLResponse_Simple(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"Message": "Hello",
	}
	err := EncodeXMLResponse(w, "SendMessage", response)
	require.NoError(t, err)
	assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "<SendMessageResponse>")
	assert.Contains(t, w.Body.String(), "<SendMessageResult>")
	assert.Contains(t, w.Body.String(), "<Message>Hello</Message>")
}

func TestEncodeXMLResponse_NestedMap(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"Attributes": map[string]interface{}{
			"Key": "Value",
		},
	}
	err := EncodeXMLResponse(w, "GetAttributes", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<Attributes>")
	assert.Contains(t, w.Body.String(), "<Key>Value</Key>")
}

func TestEncodeXMLResponse_StringMap(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"Tags": map[string]string{
			"Name":  "test",
			"Value": "example",
		},
	}
	err := EncodeXMLResponse(w, "ListTags", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<Tags>")
	assert.Contains(t, w.Body.String(), "<entry>")
	assert.Contains(t, w.Body.String(), "<key>Name</key>")
	assert.Contains(t, w.Body.String(), "<value>test</value>")
}

func TestEncodeXMLResponse_Slice(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"Items": []interface{}{
			map[string]interface{}{"Id": "1"},
			map[string]interface{}{"Id": "2"},
		},
	}
	err := EncodeXMLResponse(w, "ListItems", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<Items>")
	assert.Contains(t, w.Body.String(), "<member>")
	assert.Contains(t, w.Body.String(), "<Id>1</Id>")
	assert.Contains(t, w.Body.String(), "<Id>2</Id>")
}

func TestEncodeXMLResponse_StringSlice(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"QueueUrls": []string{"url1", "url2"},
	}
	err := EncodeXMLResponse(w, "ListQueues", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<QueueUrls>")
	assert.Contains(t, w.Body.String(), "<member>url1</member>")
	assert.Contains(t, w.Body.String(), "<member>url2</member>")
}

func TestEncodeXMLResponse_MapSlice(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"Messages": []map[string]interface{}{
			{"Body": "msg1"},
			{"Body": "msg2"},
		},
	}
	err := EncodeXMLResponse(w, "ReceiveMessage", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<Messages>")
	assert.Contains(t, w.Body.String(), "<Body>msg1</Body>")
	assert.Contains(t, w.Body.String(), "<Body>msg2</Body>")
}

func TestEncodeXMLResponse_VariousTypes(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"StringVal": "text",
		"IntVal":    42,
		"Int64Val":  int64(100),
		"FloatVal":  3.14,
		"BoolVal":   true,
	}
	err := EncodeXMLResponse(w, "Test", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<StringVal>text</StringVal>")
	assert.Contains(t, w.Body.String(), "<IntVal>42</IntVal>")
	assert.Contains(t, w.Body.String(), "<BoolVal>true</BoolVal>")
}

func TestEncodeQueryXMLResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"QueueUrl": "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue",
	}
	err := EncodeQueryXMLResponse(w, "CreateQueue", response)
	require.NoError(t, err)
	assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "<CreateQueueResponse>")
	assert.Contains(t, w.Body.String(), "<CreateQueueResult>")
	assert.Contains(t, w.Body.String(), "<QueueUrl>")
	assert.Contains(t, w.Body.String(), "<ResponseMetadata>")
	assert.Contains(t, w.Body.String(), "<RequestId>")
}

func TestEncodeQueryXMLResponse_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	err := EncodeQueryXMLResponse(w, "DeleteQueue", map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<DeleteQueueResponse>")
	assert.Contains(t, w.Body.String(), "<DeleteQueueResult>")
	assert.Contains(t, w.Body.String(), "</DeleteQueueResult>")
}

func TestAnyToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"int64", int64(100), "100"},
		{"float64", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anyToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"zero", 0, "0"},
		{"positive", 42, "42"},
		{"negative", -10, "-10"},
		{"large", 123456789, "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetProtocolHeaders_All(t *testing.T) {
	w := httptest.NewRecorder()
	SetProtocolHeaders(w, "application/json", "SQS.SendMessage", "rpc-v2", true)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "SQS.SendMessage", w.Header().Get("X-Amz-Target"))
	assert.Equal(t, "rpc-v2", w.Header().Get("smithy-protocol"))
	assert.Equal(t, "true", w.Header().Get("X-Amzn-Query-Mode"))
}

func TestSetProtocolHeaders_Partial(t *testing.T) {
	w := httptest.NewRecorder()
	SetProtocolHeaders(w, "application/xml", "", "", false)

	assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
	assert.Empty(t, w.Header().Get("X-Amz-Target"))
	assert.Empty(t, w.Header().Get("smithy-protocol"))
	assert.Empty(t, w.Header().Get("X-Amzn-Query-Mode"))
}

func TestSetProtocolHeaders_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	SetProtocolHeaders(w, "", "", "", false)

	assert.Empty(t, w.Header().Get("Content-Type"))
	assert.Empty(t, w.Header().Get("X-Amz-Target"))
}

func TestEncodeXMLResponse_SliceWithSimpleValues(t *testing.T) {
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"Numbers": []interface{}{1, 2, 3},
	}
	err := EncodeXMLResponse(w, "Test", response)
	require.NoError(t, err)
	assert.Contains(t, w.Body.String(), "<Numbers>")
	assert.Contains(t, w.Body.String(), "<member>1</member>")
	assert.Contains(t, w.Body.String(), "<member>2</member>")
	assert.Contains(t, w.Body.String(), "<member>3</member>")
}

func TestNumberToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"int", 42, "42"},
		{"int64", int64(100), "100"},
		{"float64", 3.14, "3.14"},
		{"other", "string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := numberToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
