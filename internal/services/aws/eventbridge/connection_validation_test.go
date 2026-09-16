package eventbridge

import (
	"context"
	"strings"
	"testing"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// seedValidationStore builds a store with one authorised BASIC connection
// the API-destination cores can reference, created through the service core
// (wired to a fake Secrets Manager invoker) so the full create lifecycle —
// including the credential secret and the authorization stamp — runs.
func seedValidationStore(t *testing.T) (*eventsstore.EventsStore, *EventsService) {
	t.Helper()
	svc, store, _ := newConnectionTestService(t)
	if _, err := svc.createConnectionCore(context.Background(), store, CreateConnectionInput{
		Name:              "ref-conn",
		AuthorizationType: "BASIC",
		AuthParameters: &eventsstore.AuthParameters{
			BasicAuthParameters: &eventsstore.BasicAuthParameters{Username: "u", Password: "p"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	return store, svc
}

// TestCreateApiDestinationValidation pins the CreateApiDestination
// validation contract: HttpMethod is a required member with no default, the
// invocation endpoint must be an https URL within the HttpsEndpoint
// pattern, the rate limit's @range floor is 1 (explicit 0 is rejected; no
// maximum is documented), and ConnectionArn must match the model pattern
// and name an existing connection.
func TestCreateApiDestinationValidation(t *testing.T) {
	ctx := context.Background()
	store, svc := seedValidationStore(t)
	connARN := connectionARN(t, store, "ref-conn")

	cases := []struct {
		name  string
		input CreateApiDestinationInput
		check func(t *testing.T, err error)
	}{
		{
			name: "valid",
			input: CreateApiDestinationInput{
				Name: "valid-dest", ConnectionArn: connARN, HttpMethod: "POST",
				InvocationEndpoint: "https://example.com/webhook",
			},
			check: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("valid destination rejected: %v", err)
				}
			},
		},
		{
			name: "httpMethod omitted is rejected (no default)",
			input: CreateApiDestinationInput{
				Name: "no-method", ConnectionArn: connARN,
				InvocationEndpoint: "https://example.com/webhook",
			},
			check: wantValidation(t, "HttpMethod is required"),
		},
		{
			name: "http endpoint rejected",
			input: CreateApiDestinationInput{
				Name: "http-dest", ConnectionArn: connARN, HttpMethod: "POST",
				InvocationEndpoint: "http://example.com/webhook",
			},
			check: wantValidation(t, "InvocationEndpoint"),
		},
		{
			name: "rate limit zero rejected",
			input: CreateApiDestinationInput{
				Name: "rate-zero", ConnectionArn: connARN, HttpMethod: "POST",
				InvocationEndpoint:     "https://example.com/webhook",
				InvocationRateLimitSet: true,
				InvocationRateLimit:    0,
			},
			check: wantValidation(t, "InvocationRateLimitPerSecond"),
		},
		{
			name: "rate limit negative rejected",
			input: CreateApiDestinationInput{
				Name: "rate-neg", ConnectionArn: connARN, HttpMethod: "POST",
				InvocationEndpoint:     "https://example.com/webhook",
				InvocationRateLimitSet: true,
				InvocationRateLimit:    -3,
			},
			check: wantValidation(t, "InvocationRateLimitPerSecond"),
		},
		{
			name: "malformed connection ARN rejected",
			input: CreateApiDestinationInput{
				Name: "bad-arn", ConnectionArn: "arn:aws:events:us-east-1:000000000000:connection/no-suffix",
				HttpMethod: "POST", InvocationEndpoint: "https://example.com/webhook",
			},
			check: wantValidation(t, "does not match the required pattern"),
		},
		{
			name: "unknown connection is ResourceNotFound",
			input: CreateApiDestinationInput{
				Name: "ghost-conn", ConnectionArn: replaceConnectionName(connARN, "ghost-conn"),
				HttpMethod: "POST", InvocationEndpoint: "https://example.com/webhook",
			},
			check: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("unknown connection must be rejected")
				}
				awsErr, ok := err.(*awserrors.AWSError)
				if !ok {
					t.Fatalf("want AWSError, got %T: %v", err, err)
				}
				if awsErr.Code != "ResourceNotFoundException" {
					t.Fatalf("want ResourceNotFoundException, got %s: %s", awsErr.Code, awsErr.Message)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.createApiDestinationCore(ctx, store, tc.input)
			tc.check(t, err)
		})
	}
}

// wantValidation returns a check asserting a ValidationException whose
// message mentions want.
func wantValidation(t *testing.T, want string) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expected ValidationException mentioning %q, got nil", want)
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("want AWSError, got %T: %v", err, err)
		}
		if awsErr.Code != "ValidationException" {
			t.Fatalf("want ValidationException, got %s: %s", awsErr.Code, awsErr.Message)
		}
	}
}

// TestDeauthorizeClearsAuthParameters pins the DeauthorizeConnection
// contract: the stored credentials are removed ("Removes all authorization
// parameters from the connection") while LastAuthorizedAt survives as the
// historical record of the connection's last authorization.
func TestDeauthorizeClearsAuthParameters(t *testing.T) {
	ctx := context.Background()
	store, svc := seedValidationStore(t)

	if _, err := svc.deauthorizeConnectionCore(ctx, store, "ref-conn"); err != nil {
		t.Fatal(err)
	}

	connection, err := store.GetConnection(ctx, "ref-conn")
	if err != nil {
		t.Fatal(err)
	}
	if connection.AuthParameters != nil {
		t.Fatalf("authorization parameters must be removed, got %+v", connection.AuthParameters)
	}
	if connection.State != eventsstore.ConnectionStateDeauthorized {
		t.Fatalf("state = %s", connection.State)
	}
	if connection.LastAuthorizedAt.IsZero() {
		t.Fatal("LastAuthorizedAt records the last authorization and must survive deauthorization")
	}
}

// TestValidateAuthParametersOAuthMethodEnum pins the OAuth token-request
// method contract: ConnectionOAuthHttpMethod admits GET, POST and PUT only.
func TestValidateAuthParametersOAuthMethodEnum(t *testing.T) {
	for _, method := range []string{"GET", "POST", "PUT"} {
		if err := validateAuthParameters("OAUTH_CLIENT_CREDENTIALS", oauthParams(method)); err != nil {
			t.Fatalf("method %s must be accepted: %v", method, err)
		}
	}
	for _, method := range []string{"DELETE", "PATCH", "HEAD", "OPTIONS"} {
		if err := validateAuthParameters("OAUTH_CLIENT_CREDENTIALS", oauthParams(method)); err == nil {
			t.Fatalf("method %s must be rejected", method)
		}
	}
}

func oauthParams(method string) *eventsstore.AuthParameters {
	return &eventsstore.AuthParameters{
		OAuthParameters: &eventsstore.OAuthParameters{
			AuthorizationEndpoint: "https://auth.example.com/token",
			HttpMethod:            method,
			ClientParameters:      &eventsstore.OAuthClientParameters{ClientID: "id", ClientSecret: "secret"},
		},
	}
}

// TestValidateConnectionHttpParameters pins the ConnectionHttpParameters
// trait contract: each family list is capped at its modelled entry count,
// header and query-string keys and values are bounded at 512 characters
// under their model character patterns (an RFC 7230 token header key,
// printable-ASCII header values, control-free query-string members with
// CR/LF still admissible in values), the length basis is Unicode code
// points, and the body family carries the list cap alone — its Key and
// Value shapes hold no modelled traits.
func TestValidateConnectionHttpParameters(t *testing.T) {
	entryAt := func(n int) *eventsstore.ConnectionHttpParameters {
		p := &eventsstore.ConnectionHttpParameters{}
		for i := 0; i < n; i++ {
			p.HeaderParameters = append(p.HeaderParameters, eventsstore.ConnectionHeaderParameter{Key: "X-Key", Value: "v"})
		}
		return p
	}

	cases := []struct {
		name    string
		params  *eventsstore.ConnectionHttpParameters
		wantErr bool
	}{
		{name: "valid families", params: &eventsstore.ConnectionHttpParameters{
			HeaderParameters:      []eventsstore.ConnectionHeaderParameter{{Key: "X-Api_Key", Value: " token  value "}},
			QueryStringParameters: []eventsstore.ConnectionQueryStringParameter{{Key: "名前", Value: "line1\nline2"}},
			BodyParameters:        []eventsstore.ConnectionBodyParameter{{Key: "any key ✓", Value: "any value ✓"}},
		}},
		{name: "header list at the cap", params: entryAt(eventsstore.ConnectionHttpParametersMaxEntries)},
		{name: "header list over the cap", params: entryAt(eventsstore.ConnectionHttpParametersMaxEntries + 1), wantErr: true},
		{name: "query list over the cap", params: &eventsstore.ConnectionHttpParameters{QueryStringParameters: make([]eventsstore.ConnectionQueryStringParameter, eventsstore.ConnectionHttpParametersMaxEntries+1)}, wantErr: true},
		{name: "body list over the cap", params: &eventsstore.ConnectionHttpParameters{BodyParameters: make([]eventsstore.ConnectionBodyParameter, eventsstore.ConnectionHttpParametersMaxEntries+1)}, wantErr: true},
		{name: "header key with a space", params: &eventsstore.ConnectionHttpParameters{HeaderParameters: []eventsstore.ConnectionHeaderParameter{{Key: "X Api-Key", Value: "v"}}}, wantErr: true},
		{name: "header value with a newline", params: &eventsstore.ConnectionHttpParameters{HeaderParameters: []eventsstore.ConnectionHeaderParameter{{Key: "X-Key", Value: "a\nb"}}}, wantErr: true},
		{name: "header value at the cap", params: &eventsstore.ConnectionHttpParameters{HeaderParameters: []eventsstore.ConnectionHeaderParameter{{Key: "X-Key", Value: strings.Repeat("v", eventsstore.ConnectionParameterMaxLength)}}}},
		{name: "header value over the cap", params: &eventsstore.ConnectionHttpParameters{HeaderParameters: []eventsstore.ConnectionHeaderParameter{{Key: "X-Key", Value: strings.Repeat("v", eventsstore.ConnectionParameterMaxLength+1)}}}, wantErr: true},
		{name: "header key over the cap", params: &eventsstore.ConnectionHttpParameters{HeaderParameters: []eventsstore.ConnectionHeaderParameter{{Key: strings.Repeat("k", eventsstore.ConnectionParameterMaxLength+1), Value: "v"}}}, wantErr: true},
		{name: "query key with a control character", params: &eventsstore.ConnectionHttpParameters{QueryStringParameters: []eventsstore.ConnectionQueryStringParameter{{Key: "a\x01b", Value: "v"}}}, wantErr: true},
		{name: "query value with a forbidden control character", params: &eventsstore.ConnectionHttpParameters{QueryStringParameters: []eventsstore.ConnectionQueryStringParameter{{Key: "k", Value: "a\x01b"}}}, wantErr: true},
		// The @length basis is Unicode code points: a 512-code-point
		// multibyte key is at the cap even though its byte length is
		// three times that.
		{name: "multibyte query key at the cap counts code points", params: &eventsstore.ConnectionHttpParameters{QueryStringParameters: []eventsstore.ConnectionQueryStringParameter{{Key: strings.Repeat("キ", eventsstore.ConnectionParameterMaxLength), Value: "v"}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConnectionHttpParameters(tc.params)
			if tc.wantErr && err == nil {
				t.Fatal("expected a ValidationException, got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected acceptance, got %v", err)
			}
		})
	}

	// The bounds ride the auth-parameters validation for both nesting
	// sites: the invocation parameters and the OAuth token-request
	// parameters share the ConnectionHttpParameters shape.
	apiKey := &eventsstore.AuthParameters{
		ApiKeyAuthParameters:     &eventsstore.ApiKeyAuthParameters{ApiKeyName: "k", ApiKeyValue: "v"},
		InvocationHttpParameters: entryAt(eventsstore.ConnectionHttpParametersMaxEntries + 1),
	}
	if err := validateAuthParameters("API_KEY", apiKey); err == nil {
		t.Fatal("invocation HTTP parameters over the cap must fail validation")
	}
	oauth := oauthParams("POST")
	oauth.OAuthParameters.OAuthHttpParameters = entryAt(eventsstore.ConnectionHttpParametersMaxEntries + 1)
	if err := validateAuthParameters("OAUTH_CLIENT_CREDENTIALS", oauth); err == nil {
		t.Fatal("OAuth HTTP parameters over the cap must fail validation")
	}
}

// TestValidateKmsKeyIdentifierForms pins the KMS identifier contract: the
// key ARN, key ID, key alias and key alias ARN are all accepted; non-KMS
// ARNs and malformed aliases are rejected.
func TestValidateKmsKeyIdentifierForms(t *testing.T) {
	accepted := []string{
		"",
		"arn:aws:kms:us-east-1:000000000000:key/1234abcd-12ab-34cd-56ef-1234567890ab",
		"arn:aws-us-gov:kms:us-gov-west-1:000000000000:key/mrk-1234abcd12ab34cd56ef1234567890ab",
		"arn:aws:kms:us-east-1:000000000000:alias/my-alias",
		"alias/my-alias",
		"1234abcd-12ab-34cd-56ef-1234567890ab",
		"mrk-1234abcd12ab34cd56ef1234567890ab",
	}
	for _, id := range accepted {
		if !validateKmsKeyIdentifier(id) {
			t.Errorf("identifier %q must be accepted", id)
		}
	}
	rejected := []string{
		"arn:aws:s3:us-east-1:000000000000:bucket/x",
		"arn:aws:kms:us-east-1:000000000000:related/thing",
		"alias/",
		"short",
		"has spaces",
	}
	for _, id := range rejected {
		if validateKmsKeyIdentifier(id) {
			t.Errorf("identifier %q must be rejected", id)
		}
	}

	// The shape's @length(0,2048) ceiling bounds every form: an alias at
	// the bound passes, one character over fails. The four identifier
	// forms are ASCII by their own patterns, so characters and bytes
	// coincide for every reachable form.
	atBound := "alias/" + strings.Repeat("a", eventsstore.KmsKeyIdentifierMaxLength-6)
	if utf8.RuneCountInString(atBound) != eventsstore.KmsKeyIdentifierMaxLength || !validateKmsKeyIdentifier(atBound) {
		t.Errorf("an identifier of exactly %d characters must be accepted", eventsstore.KmsKeyIdentifierMaxLength)
	}
	overBound := "alias/" + strings.Repeat("a", eventsstore.KmsKeyIdentifierMaxLength-5)
	if validateKmsKeyIdentifier(overBound) {
		t.Errorf("an identifier over %d characters must be rejected", eventsstore.KmsKeyIdentifierMaxLength)
	}
}

// connectionARN returns the stored ARN of a seeded connection.
func connectionARN(t *testing.T, store *eventsstore.EventsStore, name string) string {
	t.Helper()
	connection, err := store.GetConnection(context.Background(), name)
	if err != nil {
		t.Fatal(err)
	}
	return connection.ARN
}

// replaceConnectionName swaps the name segment of a connection ARN,
// producing an ARN that names a different (possibly non-existent)
// connection.
func replaceConnectionName(arn, name string) string {
	return strings.Replace(arn, "connection/ref-conn/", "connection/"+name+"/", 1)
}
