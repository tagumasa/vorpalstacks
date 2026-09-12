package cognitoidentityprovider

import (
	"context"
	"regexp"
	"testing"
)

// A service-generated secret always satisfies the ClientSecretType shape:
// length inside [24, 64] and characters inside the pattern class [\w+] —
// base64 output could never satisfy it (padding '=' and member '/' fall
// outside the class).
func TestGenerateSecretValueMatchesClientSecretType(t *testing.T) {
	idPattern := regexp.MustCompile(`^[\w+]+$`)
	for i := 0; i < 100; i++ {
		value, err := generateSecretValue()
		if err != nil {
			t.Fatal(err)
		}
		if len(value) < clientSecretMinLength || len(value) > clientSecretMaxLength {
			t.Fatalf("generated length %d outside [%d, %d]", len(value), clientSecretMinLength, clientSecretMaxLength)
		}
		if !idPattern.MatchString(value) {
			t.Fatalf("generated value %q violates the pattern class", value)
		}
	}
}

// A caller-supplied secret must satisfy the ClientSecretType bounds — the
// same validation the client-create path applies.
func TestAddUserPoolClientSecretValidatesCustomSecret(t *testing.T) {
	env := newChallengeTestEnv(t)

	if _, err := env.svc.AddUserPoolClientSecret(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"UserPoolId":   env.pool.ID,
		"ClientId":     challengeTestClientID,
		"ClientSecret": "abc",
	})); err != ErrInvalidParameter {
		t.Fatalf("short custom secret returned %v, want InvalidParameterException", err)
	}
	if _, err := env.svc.AddUserPoolClientSecret(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"UserPoolId":   env.pool.ID,
		"ClientId":     challengeTestClientID,
		"ClientSecret": "contains=padding",
	})); err != ErrInvalidParameter {
		t.Fatalf("out-of-class custom secret returned %v, want InvalidParameterException", err)
	}
}

// The value is revealed only on the Add response of a service-generated
// secret; the List response never carries it, a custom secret is never
// echoed, and the descriptor identifier follows the documented
// <client-id>--<epoch-create-time> format.
func TestAddAndListUserPoolClientSecretRevealRules(t *testing.T) {
	env := newChallengeTestEnv(t)
	idFormat := regexp.MustCompile(`^` + regexp.QuoteMeta(challengeTestClientID) + `--\d+$`)

	resp, err := env.svc.AddUserPoolClientSecret(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"UserPoolId": env.pool.ID,
		"ClientId":   challengeTestClientID,
	}))
	if err != nil {
		t.Fatal(err)
	}
	generated := resp.(map[string]interface{})["ClientSecretDescriptor"].(map[string]interface{})
	if value, ok := generated["ClientSecretValue"].(string); !ok || value == "" {
		t.Fatal("generated secret's Add response must carry the value")
	}
	if id, ok := generated["ClientSecretId"].(string); !ok || !idFormat.MatchString(id) {
		t.Fatalf("descriptor id %v does not follow <client-id>--<epoch-create-time>", generated["ClientSecretId"])
	}

	const customSecret = "custom_secret_value_0123456789ab"
	resp, err = env.svc.AddUserPoolClientSecret(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"UserPoolId":   env.pool.ID,
		"ClientId":     challengeTestClientID,
		"ClientSecret": customSecret,
	}))
	if err != nil {
		t.Fatal(err)
	}
	custom := resp.(map[string]interface{})["ClientSecretDescriptor"].(map[string]interface{})
	if _, ok := custom["ClientSecretValue"]; ok {
		t.Fatal("custom secret's Add response must not echo the value")
	}

	// The stored record holds both values — one generated, one custom.
	client, err := env.store.GetUserPoolClient(env.pool.ID, challengeTestClientID)
	if err != nil {
		t.Fatal(err)
	}
	if len(client.ClientSecrets) != 2 {
		t.Fatalf("stored %d secret descriptors, want 2", len(client.ClientSecrets))
	}
	if !client.ClientSecrets[0].Generated || client.ClientSecrets[0].ClientSecretValue == "" {
		t.Fatal("first descriptor must be the generated secret with its value stored")
	}
	if client.ClientSecrets[1].Generated || client.ClientSecrets[1].ClientSecretValue != customSecret {
		t.Fatal("second descriptor must store the custom value with Generated=false")
	}

	// The list response never reveals a stored value.
	listResp, err := env.svc.ListUserPoolClientSecrets(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"UserPoolId": env.pool.ID,
		"ClientId":   challengeTestClientID,
	}))
	if err != nil {
		t.Fatal(err)
	}
	for _, raw := range listResp.(map[string]interface{})["ClientSecrets"].([]map[string]interface{}) {
		if _, ok := raw["ClientSecretValue"]; ok {
			t.Fatalf("list response revealed a secret value: %#v", raw)
		}
	}
}
