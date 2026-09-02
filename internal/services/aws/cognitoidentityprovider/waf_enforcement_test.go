package cognitoidentityprovider

import (
	"encoding/json"
	"testing"
)

// TestRedactPIIForWAF pins the PII the user-pool WAF documentation keeps
// away from AWS WAF: usernames, passwords, phone numbers and email
// addresses are not available to web ACL rules, so the members are
// dropped before inspection while unrelated content survives.
func TestRedactPIIForWAF(t *testing.T) {
	untouched := redactPIIForWAF([]byte(`{"ClientId":"abc","ClientMetadata":{"trace":"keep"}}`))
	if string(untouched) != `{"ClientId":"abc","ClientMetadata":{"trace":"keep"}}` {
		t.Fatalf("body without PII was rewritten: %s", untouched)
	}

	redacted := redactPIIForWAF([]byte(`{"Username":"eve","Password":"s3cret","UserAttributes":[
		{"Name":"email","Value":"eve@example.test"},
		{"Name":"custom:team","Value":"ops"},
		{"Name":"phone_number","Value":"+819000000000"}],
		"AuthParameters":{"USERNAME":"eve","PASSWORD":"p","SECRET_HASH":"h","NEW_PASSWORD":"n"}}`))
	var parsed map[string]interface{}
	if err := json.Unmarshal(redacted, &parsed); err != nil {
		t.Fatalf("redacted body is not JSON: %v", err)
	}
	for _, key := range []string{"Username", "Password"} {
		if _, ok := parsed[key]; ok {
			t.Errorf("%s survived the redaction", key)
		}
	}
	attributes, _ := parsed["UserAttributes"].([]interface{})
	if len(attributes) != 1 {
		t.Fatalf("non-PII attributes must survive, got %v", parsed["UserAttributes"])
	}
	if attribute, _ := attributes[0].(map[string]interface{}); attribute["Name"] != "custom:team" {
		t.Errorf("unexpected surviving attribute %v", attributes[0])
	}
	authParams, _ := parsed["AuthParameters"].(map[string]interface{})
	if _, ok := authParams["USERNAME"]; !ok {
		t.Error("AuthParameters USERNAME must survive (it is not in the documented PII set)")
	}
	for _, key := range []string{"PASSWORD", "SECRET_HASH", "NEW_PASSWORD"} {
		if _, ok := authParams[key]; ok {
			t.Errorf("AuthParameters %s survived the redaction", key)
		}
	}

	if got := string(redactPIIForWAF([]byte("not json"))); got != "not json" {
		t.Fatalf("unparseable body must pass through, got %s", got)
	}
	if got := redactPIIForWAF(nil); got != nil {
		t.Fatalf("empty body must pass through, got %v", got)
	}
}
