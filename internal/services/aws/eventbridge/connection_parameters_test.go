package eventbridge

import (
	"encoding/json"
	"reflect"
	"testing"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// sdkAuthParametersWire is the AuthParameters JSON exactly as the AWS SDK
// serialises it for CreateConnection: every header/query/body family is a
// list of {IsValueSecret, Key, Value} structures.
const sdkAuthParametersWire = `{
	"InvocationHttpParameters": {
		"BodyParameters": [
			{"IsValueSecret": true, "Key": "secretBodyKey", "Value": "secret-body"},
			{"IsValueSecret": false, "Key": "plainBodyKey", "Value": "plain-body"}
		],
		"HeaderParameters": [
			{"IsValueSecret": false, "Key": "X-Plain", "Value": "plain-header"},
			{"IsValueSecret": true, "Key": "X-Secret", "Value": "secret-header"}
		],
		"QueryStringParameters": [
			{"IsValueSecret": true, "Key": "q", "Value": "secret-query"}
		]
	}
}`

// TestParseConnectionHttpParametersSdkListShape pins the wire contract the
// SDK actually sends: parameter families decode from lists of
// {IsValueSecret, Key, Value} structures, keys, values and secret flags all
// surviving the parse.
func TestParseConnectionHttpParametersSdkListShape(t *testing.T) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(sdkAuthParametersWire), &raw); err != nil {
		t.Fatal(err)
	}

	parsed := parseConnectionHttpParameters(raw["InvocationHttpParameters"].(map[string]interface{}))
	if parsed == nil {
		t.Fatal("InvocationHttpParameters must parse from the SDK wire shape")
	}
	hp := parsed
	if hp == nil {
		t.Fatal("connection http parameters must be present")
	}

	wantHeaders := []eventsstore.ConnectionHeaderParameter{
		{Key: "X-Plain", Value: "plain-header", IsValueSecret: false},
		{Key: "X-Secret", Value: "secret-header", IsValueSecret: true},
	}
	if !reflect.DeepEqual(hp.HeaderParameters, wantHeaders) {
		t.Fatalf("header parameters: got %+v want %+v", hp.HeaderParameters, wantHeaders)
	}
	wantQuery := []eventsstore.ConnectionQueryStringParameter{
		{Key: "q", Value: "secret-query", IsValueSecret: true},
	}
	if !reflect.DeepEqual(hp.QueryStringParameters, wantQuery) {
		t.Fatalf("query string parameters: got %+v want %+v", hp.QueryStringParameters, wantQuery)
	}
	wantBody := []eventsstore.ConnectionBodyParameter{
		{Key: "secretBodyKey", Value: "secret-body", IsValueSecret: true},
		{Key: "plainBodyKey", Value: "plain-body", IsValueSecret: false},
	}
	if !reflect.DeepEqual(hp.BodyParameters, wantBody) {
		t.Fatalf("body parameters: got %+v want %+v", hp.BodyParameters, wantBody)
	}
}

// TestParseConnectionHttpParametersOauthTree pins the OAuth branch: the
// client parameters, endpoint and method parse alongside the token-request
// parameter families.
func TestParseConnectionHttpParametersOauthTree(t *testing.T) {
	wire := `{
		"OAuthParameters": {
			"AuthorizationEndpoint": "https://auth.example.com/token",
			"HttpMethod": "POST",
			"ClientParameters": {"ClientID": "id-1", "ClientSecret": "secret-1"},
			"OAuthHttpParameters": {
				"HeaderParameters": [
					{"IsValueSecret": true, "Key": "X-Api-Key", "Value": "token-key"}
				]
			}
		}
	}`
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(wire), &raw); err != nil {
		t.Fatal(err)
	}

	parsed := parseAuthParameters(raw)
	if parsed == nil || parsed.OAuthParameters == nil {
		t.Fatal("OAuth parameters must parse")
	}
	oauth := parsed.OAuthParameters
	if oauth.AuthorizationEndpoint != "https://auth.example.com/token" || oauth.HttpMethod != "POST" {
		t.Fatalf("oauth endpoint/method: %+v", oauth)
	}
	if oauth.ClientParameters == nil || oauth.ClientParameters.ClientID != "id-1" || oauth.ClientParameters.ClientSecret != "secret-1" {
		t.Fatalf("client parameters: %+v", oauth.ClientParameters)
	}
	if oauth.OAuthHttpParameters == nil || len(oauth.OAuthHttpParameters.HeaderParameters) != 1 {
		t.Fatalf("oauth http parameters: %+v", oauth.OAuthHttpParameters)
	}
	if got := oauth.OAuthHttpParameters.HeaderParameters[0]; got.Key != "X-Api-Key" || got.Value != "token-key" || !got.IsValueSecret {
		t.Fatalf("oauth header parameter: %+v", got)
	}
}

// TestRedactAuthParametersResponseShape pins the DescribeConnection
// AuthParameters response structure: the credential members (Password,
// ClientSecret, ApiKeyValue) are absent from the response shape entirely,
// while Username, ClientID and ApiKeyName are returned as stored and the
// parameter lists keep every entry's Key, Value and IsValueSecret.
func TestRedactAuthParametersResponseShape(t *testing.T) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(sdkAuthParametersWire), &raw); err != nil {
		t.Fatal(err)
	}
	hp := parseConnectionHttpParameters(raw["InvocationHttpParameters"].(map[string]interface{}))
	full := &eventsstore.AuthParameters{
		BasicAuthParameters: &eventsstore.BasicAuthParameters{Username: "user-1", Password: "pass-1"},
		OAuthParameters: &eventsstore.OAuthParameters{
			AuthorizationEndpoint: "https://auth.example.com/token",
			HttpMethod:            "POST",
			ClientParameters:      &eventsstore.OAuthClientParameters{ClientID: "id-1", ClientSecret: "secret-1"},
		},
		ApiKeyAuthParameters:     &eventsstore.ApiKeyAuthParameters{ApiKeyName: "X-Api-Key", ApiKeyValue: "key-value"},
		InvocationHttpParameters: hp,
	}

	out := redactAuthParameters(full)

	basic := out["BasicAuthParameters"].(map[string]interface{})
	if basic["Username"] != "user-1" {
		t.Fatalf("Username must be returned as stored, got %v", basic["Username"])
	}
	if _, ok := basic["Password"]; ok {
		t.Fatal("Password is absent from the response shape and must not be rendered")
	}

	oauth := out["OAuthParameters"].(map[string]interface{})
	client := oauth["ClientParameters"].(map[string]interface{})
	if client["ClientID"] != "id-1" {
		t.Fatalf("ClientID must be returned as stored, got %v", client["ClientID"])
	}
	if _, ok := client["ClientSecret"]; ok {
		t.Fatal("ClientSecret is absent from the response shape and must not be rendered")
	}

	apiKey := out["ApiKeyAuthParameters"].(map[string]interface{})
	if apiKey["ApiKeyName"] != "X-Api-Key" {
		t.Fatalf("ApiKeyName must be returned as stored, got %v", apiKey["ApiKeyName"])
	}
	if _, ok := apiKey["ApiKeyValue"]; ok {
		t.Fatal("ApiKeyValue is absent from the response shape and must not be rendered")
	}

	invocation := out["InvocationHttpParameters"].(map[string]interface{})
	headers := invocation["HeaderParameters"].([]map[string]interface{})
	if len(headers) != 2 {
		t.Fatalf("both header entries must be returned, got %d", len(headers))
	}
	if headers[1]["Key"] != "X-Secret" || headers[1]["Value"] != "secret-header" || headers[1]["IsValueSecret"] != true {
		t.Fatalf("header entry must keep Key, Value and IsValueSecret: %+v", headers[1])
	}
	bodies := invocation["BodyParameters"].([]map[string]interface{})
	if len(bodies) != 2 || bodies[0]["Key"] != "secretBodyKey" {
		t.Fatalf("body entries must be returned with their keys: %+v", bodies)
	}
}
