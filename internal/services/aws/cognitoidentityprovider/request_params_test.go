package cognitoidentityprovider

import (
	"testing"

	"vorpalstacks/internal/common/request"
)

// The service protocol is awsJson1_1: members are PascalCase keys over typed
// JSON values. A query-string-style body — dotted keys, string-typed numbers
// and booleans — is not a second accepted shape; these pins hold the shared
// readers to the JSON shape alone.

func TestParseNamedAttributeListIgnoresQueryFormKeys(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserAttributes.1.Name":  "email",
		"UserAttributes.1.Value": "user@example.test",
	}}
	if attrs := parseNamedAttributeList(req, "UserAttributes"); len(attrs) != 0 {
		t.Fatalf("query-form attribute keys parsed: %+v", attrs)
	}

	jsonReq := &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserAttributes": []interface{}{
			map[string]interface{}{"Name": "email", "Value": "user@example.test"},
		},
	}}
	attrs := parseNamedAttributeList(jsonReq, "UserAttributes")
	if attrs["email"] != "user@example.test" {
		t.Fatalf("JSON array member not parsed: %+v", attrs)
	}
}

func TestParseClientMetadataIgnoresQueryFormKeys(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"ClientMetadata.key": "value",
	}}
	if md := parseClientMetadata(req); len(md) != 0 {
		t.Fatalf("query-form metadata keys parsed: %+v", md)
	}

	jsonReq := &request.ParsedRequest{Parameters: map[string]interface{}{
		"ClientMetadata": map[string]interface{}{"key": "value"},
	}}
	if md := parseClientMetadata(jsonReq); md["key"] != "value" {
		t.Fatalf("JSON map member not parsed: %+v", md)
	}
}

func TestGetStringSliceParamIgnoresQueryFormKeys(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"CallbackURLs.1": "https://example.test/cb",
		"CallbackURLs.2": "https://other.test/cb",
	}}
	if urls := getStringSliceParam(req, "CallbackURLs"); len(urls) != 0 {
		t.Fatalf("query-form list keys parsed: %+v", urls)
	}

	jsonReq := &request.ParsedRequest{Parameters: map[string]interface{}{
		"CallbackURLs": []interface{}{"https://example.test/cb"},
	}}
	if urls := getStringSliceParam(jsonReq, "CallbackURLs"); len(urls) != 1 || urls[0] != "https://example.test/cb" {
		t.Fatalf("JSON array member not parsed: %+v", urls)
	}
}

func TestGetBoolParamReadsOnlyTypedJSON(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"ForceAliasCreation": true,
		"Permanent":          "true",
	}}
	if !getBoolParam(req, "ForceAliasCreation") {
		t.Fatal("typed JSON boolean not read")
	}
	if getBoolParam(req, "Permanent") {
		t.Fatal("string-typed boolean accepted")
	}
	if _, present := getBoolParamOK(req, "Permanent"); present {
		t.Fatal("string-typed boolean reported present")
	}
	if _, present := getBoolParamOK(req, "Absent"); present {
		t.Fatal("absent member reported present")
	}
}

func TestParseIntParamReadsOnlyTypedJSON(t *testing.T) {
	for name, value := range map[string]interface{}{
		"float64": float64(60),
		"int":     60,
		"int64":   int64(60),
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{"Limit": value}}
		v, ok := parseIntParam(req, "Limit")
		if !ok || v != 60 {
			t.Fatalf("%s: parseIntParam = %d, %v; want 60, true", name, v, ok)
		}
	}

	req := &request.ParsedRequest{Parameters: map[string]interface{}{"Limit": "60"}}
	if v, ok := parseIntParam(req, "Limit"); ok {
		t.Fatalf("string-typed number accepted: %d", v)
	}
}
