package eventbridge

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/store/aws/eventbridge"
)

// fakeSecretsInvoker is an in-memory invokers.SecretsManagerInvoker for
// service tests: it hands out ARNs of the real Secrets Manager form and
// records deletions so credential-lifecycle assertions can observe them.
type fakeSecretsInvoker struct {
	mu        sync.Mutex
	nextID    int
	secrets   map[string]string // ARN -> secret string
	deleted   []string
	getRegion string
	deleteErr error
}

func newFakeSecretsInvoker() *fakeSecretsInvoker {
	return &fakeSecretsInvoker{secrets: map[string]string{}}
}

func (f *fakeSecretsInvoker) CreateServiceSecret(ctx context.Context, region, name, secretString, description string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextID++
	arn := fmt.Sprintf("arn:aws:secretsmanager:%s:000000000000:secret:%s-AbCdEf", region, name)
	f.secrets[arn] = secretString
	return arn, nil
}

func (f *fakeSecretsInvoker) GetServiceSecretString(ctx context.Context, region, secretId string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.getRegion = region
	value, ok := f.secrets[secretId]
	if !ok {
		return "", awserrors.NewAWSError("ResourceNotFoundException", "secret not found", http.StatusNotFound)
	}
	return value, nil
}

func (f *fakeSecretsInvoker) UpdateServiceSecretString(ctx context.Context, region, secretId, secretString string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.secrets[secretId]; !ok {
		return awserrors.NewAWSError("ResourceNotFoundException", "secret not found", http.StatusNotFound)
	}
	f.secrets[secretId] = secretString
	return nil
}

func (f *fakeSecretsInvoker) DeleteServiceSecret(ctx context.Context, region, secretId string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.secrets, secretId)
	f.deleted = append(f.deleted, secretId)
	return nil
}

// secretValue returns the stored payload of a secret ARN.
func (f *fakeSecretsInvoker) secretValue(arn string) (string, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	value, ok := f.secrets[arn]
	return value, ok
}

// deletedARNs returns every deleted secret ARN.
func (f *fakeSecretsInvoker) deletedARNs() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.deleted...)
}

// newConnectionTestService builds a service wired to a bus carrying a fake
// Secrets Manager invoker, a fresh events store, and the fake itself.
func newConnectionTestService(t *testing.T) (*EventsService, *eventbridge.EventsStore, *fakeSecretsInvoker) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventbridge.NewEventsStore(st, "000000000000", "us-east-1")
	svc := NewEventsService(nil, "000000000000")
	fake := newFakeSecretsInvoker()
	bus := eventbus.NewEventBus()
	bus.SetSecretsManagerInvoker(fake)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(svc.Close)
	return svc, store, fake
}

// TestConnectionCredentialsLiveInTheSecret pins the credential-storage
// design: the record carries the non-credential halves plus the secret ARN
// (SecretArn is populated — "The ARN of the secret created from the
// authorization parameters specified for the connection"), and the secret
// holds the full tree including every credential.
func TestConnectionCredentialsLiveInTheSecret(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)

	connection, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "cred-conn",
		AuthorizationType: "BASIC",
		AuthParameters: &eventbridge.AuthParameters{
			BasicAuthParameters: &eventbridge.BasicAuthParameters{Username: "user-9", Password: "pass-9"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if connection.SecretArn == "" {
		t.Fatal("SecretArn must be populated on create")
	}
	stored, err := store.GetConnection(ctx, "cred-conn")
	if err != nil {
		t.Fatal(err)
	}
	if stored.AuthParameters.BasicAuthParameters.Password != "" {
		t.Fatal("the record must not keep the Basic password in plaintext")
	}
	if stored.AuthParameters.BasicAuthParameters.Username != "user-9" {
		t.Fatalf("Username is a non-credential member and stays on the record, got %q", stored.AuthParameters.BasicAuthParameters.Username)
	}
	value, ok := fake.secretValue(connection.SecretArn)
	if !ok {
		t.Fatal("credential secret must exist")
	}
	if !strings.Contains(value, "pass-9") {
		t.Fatalf("the secret carries the full credentials, got %s", value)
	}
}

// TestConnectionCredentialLifecycle pins the secret's follow-the-record
// lifecycle: an update rewrites the secret, a deauthorization removes it,
// and deleting the connection removes it too.
func TestConnectionCredentialLifecycle(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)

	connection, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "life-conn",
		AuthorizationType: "API_KEY",
		AuthParameters: &eventbridge.AuthParameters{
			ApiKeyAuthParameters: &eventbridge.ApiKeyAuthParameters{ApiKeyName: "X-Key", ApiKeyValue: "v1"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := svc.updateConnectionCore(ctx, store, UpdateConnectionInput{
		Name:              "life-conn",
		AuthParametersSet: true,
		AuthParameters: &eventbridge.AuthParameters{
			ApiKeyAuthParameters: &eventbridge.ApiKeyAuthParameters{ApiKeyName: "X-Key", ApiKeyValue: "v2"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.SecretArn != connection.SecretArn {
		t.Fatalf("an update rewrites the existing secret, ARN changed %s -> %s", connection.SecretArn, updated.SecretArn)
	}
	if value, _ := fake.secretValue(connection.SecretArn); !strings.Contains(value, "v2") {
		t.Fatalf("the secret must carry the updated credentials, got %s", value)
	}
	stored, err := store.GetConnection(ctx, "life-conn")
	if err != nil {
		t.Fatal(err)
	}
	if stored.AuthParameters.ApiKeyAuthParameters.ApiKeyValue != "" {
		t.Fatal("the updated record must not keep the API key value in plaintext")
	}

	if _, err := svc.deauthorizeConnectionCore(ctx, store, "life-conn"); err != nil {
		t.Fatal(err)
	}
	if len(fake.deletedARNs()) != 1 || fake.deletedARNs()[0] != connection.SecretArn {
		t.Fatalf("deauthorization must remove the credential secret, deleted %v", fake.deletedARNs())
	}
	if _, ok := fake.secretValue(connection.SecretArn); ok {
		t.Fatal("the credential secret must be gone after deauthorization")
	}

	// A deauthorized connection can be re-authorized by an update: the new
	// credentials get a fresh secret.
	reauthed, err := svc.updateConnectionCore(ctx, store, UpdateConnectionInput{
		Name:              "life-conn",
		AuthParametersSet: true,
		AuthParameters: &eventbridge.AuthParameters{
			ApiKeyAuthParameters: &eventbridge.ApiKeyAuthParameters{ApiKeyName: "X-Key", ApiKeyValue: "v3"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if reauthed.SecretArn == "" || reauthed.SecretArn == connection.SecretArn {
		t.Fatalf("re-authorization mints a fresh secret, got %q", reauthed.SecretArn)
	}
}

// TestUpdateConnectionAuthorizationTypeSemantics pins the authorization-type
// update contract: the re-authorization trigger is a credential update, so a
// type-only change is rejected (its credentials cannot be established from
// the previous type's secret), a type change with its parameters re-marks
// AUTHORIZED, and a type echo never re-authorizes a deauthorized connection.
func TestUpdateConnectionAuthorizationTypeSemantics(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)

	_, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "type-conn",
		AuthorizationType: "OAUTH_CLIENT_CREDENTIALS",
		AuthParameters:    oauthParamsForUpdate("POST"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// A type-only change is rejected and the record keeps the old type.
	_, err = svc.updateConnectionCore(ctx, store, UpdateConnectionInput{
		Name:                 "type-conn",
		AuthorizationTypeSet: true,
		AuthorizationType:    "BASIC",
	})
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
		t.Fatalf("a type-only change must be rejected as ValidationException, got %v", err)
	}
	stored, err := store.GetConnection(ctx, "type-conn")
	if err != nil {
		t.Fatal(err)
	}
	if stored.AuthorizationType != "OAUTH_CLIENT_CREDENTIALS" || stored.State != eventbridge.ConnectionStateAuthorized {
		t.Fatalf("the rejected update must leave the record untouched, got type %s state %s",
			stored.AuthorizationType, stored.State)
	}

	// The type change with its credentials re-marks AUTHORIZED and rewrites
	// the secret with the new type's parameters.
	updated, err := svc.updateConnectionCore(ctx, store, UpdateConnectionInput{
		Name:                 "type-conn",
		AuthorizationTypeSet: true,
		AuthorizationType:    "BASIC",
		AuthParametersSet:    true,
		AuthParameters: &eventbridge.AuthParameters{
			BasicAuthParameters: &eventbridge.BasicAuthParameters{Username: "u", Password: "p"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.AuthorizationType != "BASIC" || updated.State != eventbridge.ConnectionStateAuthorized {
		t.Fatalf("type change with credentials: type %s state %s", updated.AuthorizationType, updated.State)
	}
	if updated.AuthParameters.BasicAuthParameters == nil {
		t.Fatal("the record must carry the new type's redacted parameters")
	}
	if value, _ := fake.secretValue(updated.SecretArn); !strings.Contains(value, "basicAuthParameters") {
		t.Fatalf("the secret must carry the new type's credentials, got %s", value)
	}

	// An AuthorizationType echo changes no credentials, so a deauthorized
	// connection stays deauthorized.
	if _, err := svc.deauthorizeConnectionCore(ctx, store, "type-conn"); err != nil {
		t.Fatal(err)
	}
	echoed, err := svc.updateConnectionCore(ctx, store, UpdateConnectionInput{
		Name:                 "type-conn",
		AuthorizationTypeSet: true,
		AuthorizationType:    "BASIC",
		DescriptionSet:       true,
		Description:          "echo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if echoed.State != eventbridge.ConnectionStateDeauthorized {
		t.Fatalf("a type echo must not re-authorize a deauthorized connection, state %s", echoed.State)
	}
	if echoed.Description != "echo" {
		t.Fatalf("the sibling member still applies, description %q", echoed.Description)
	}
}

func oauthParamsForUpdate(method string) *eventbridge.AuthParameters {
	return &eventbridge.AuthParameters{
		OAuthParameters: &eventbridge.OAuthParameters{
			AuthorizationEndpoint: "https://auth.example.com/token",
			HttpMethod:            method,
			ClientParameters:      &eventbridge.OAuthClientParameters{ClientID: "id", ClientSecret: "secret"},
		},
	}
}

// TestConnectionCreateValidatesBeforeSecretCreation pins the create path's
// side-effect ordering: an input that fails validation must fail while the
// request is still side-effect free — no credential secret is created, so a
// rejected request cannot leak one on every retry.
func TestConnectionCreateValidatesBeforeSecretCreation(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)

	base := CreateConnectionInput{
		AuthorizationType: "BASIC",
		AuthParameters: &eventbridge.AuthParameters{
			BasicAuthParameters: &eventbridge.BasicAuthParameters{Username: "u", Password: "p"},
		},
	}

	long := base
	long.Name = "leak-desc"
	long.DescriptionSet = true
	long.Description = strings.Repeat("d", 513)
	_, err := svc.createConnectionCore(ctx, store, long)
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
		t.Fatalf("over-long Description must be a ValidationException, got %v", err)
	}

	badKey := base
	badKey.Name = "leak-kms"
	badKey.KmsKeyIdentifierSet = true
	badKey.KmsKeyIdentifier = "not-a-key-identifier!"
	_, err = svc.createConnectionCore(ctx, store, badKey)
	if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
		t.Fatalf("invalid KmsKeyIdentifier must be a ValidationException, got %v", err)
	}

	if len(fake.secrets) != 0 || len(fake.deletedARNs()) != 0 {
		t.Fatalf("rejected inputs must not create credential secrets, created %v, deleted %v",
			fake.secrets, fake.deletedARNs())
	}
	for _, name := range []string{"leak-desc", "leak-kms"} {
		if _, err := store.GetConnection(ctx, name); err != eventbridge.ErrConnectionNotFound {
			t.Fatalf("no record may exist for rejected input %s, got err %v", name, err)
		}
	}
}

// TestConnectionDeauthorizeAtomicity pins the mutation-scoped
// deauthorization: a failing secret deletion aborts the whole record
// change (the connection keeps its authorisation and its claim on the
// secret), and a deauthorization deletes the secret the record names at
// mutation time — after a re-authorization minted a fresh secret, the
// next deauthorization removes that one, never an older ARN.
func TestConnectionDeauthorizeAtomicity(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)

	connection, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "atomic-conn",
		AuthorizationType: "API_KEY",
		AuthParameters: &eventbridge.AuthParameters{
			ApiKeyAuthParameters: &eventbridge.ApiKeyAuthParameters{ApiKeyName: "k", ApiKeyValue: "v1"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// A failing deletion is an atomic abort: error returned, record and
	// secret untouched.
	fake.deleteErr = errors.New("secrets manager unavailable")
	if _, err := svc.deauthorizeConnectionCore(ctx, store, "atomic-conn"); err == nil {
		t.Fatal("a failing secret deletion must fail the deauthorization")
	}
	stored, err := store.GetConnection(ctx, "atomic-conn")
	if err != nil {
		t.Fatal(err)
	}
	if stored.State != eventbridge.ConnectionStateAuthorized || stored.SecretArn != connection.SecretArn {
		t.Fatalf("an aborted deauthorization must leave the record intact, got state %s secret %q",
			stored.State, stored.SecretArn)
	}
	if _, ok := fake.secretValue(connection.SecretArn); !ok {
		t.Fatal("the secret must survive an aborted deauthorization")
	}

	// Deauthorize, re-authorize (fresh secret), deauthorize again: the
	// second deauthorization must remove the secret the record names then.
	fake.deleteErr = nil
	if _, err := svc.deauthorizeConnectionCore(ctx, store, "atomic-conn"); err != nil {
		t.Fatal(err)
	}
	reauthed, err := svc.updateConnectionCore(ctx, store, UpdateConnectionInput{
		Name:              "atomic-conn",
		AuthParametersSet: true,
		AuthParameters: &eventbridge.AuthParameters{
			ApiKeyAuthParameters: &eventbridge.ApiKeyAuthParameters{ApiKeyName: "k", ApiKeyValue: "v2"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if reauthed.SecretArn == "" || reauthed.SecretArn == connection.SecretArn {
		t.Fatalf("re-authorization mints a fresh secret, got %q", reauthed.SecretArn)
	}
	deauthorized, err := svc.deauthorizeConnectionCore(ctx, store, "atomic-conn")
	if err != nil {
		t.Fatal(err)
	}
	if deauthorized.State != eventbridge.ConnectionStateDeauthorized || deauthorized.SecretArn != "" {
		t.Fatalf("deauthorized record = state %s secret %q", deauthorized.State, deauthorized.SecretArn)
	}
	deleted := fake.deletedARNs()
	if len(deleted) != 2 || deleted[0] != connection.SecretArn || deleted[1] != reauthed.SecretArn {
		t.Fatalf("deauthorizations must delete the record's current secret each time, deleted %v", deleted)
	}
	if _, ok := fake.secretValue(reauthed.SecretArn); ok {
		t.Fatal("the re-authorization's secret must be removed by the second deauthorization")
	}
}

// TestConnectionDeleteRemovesCredentialSecret pins the delete path's
// compensating secret removal.
func TestConnectionDeleteRemovesCredentialSecret(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)

	connection, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "del-conn",
		AuthorizationType: "BASIC",
		AuthParameters: &eventbridge.AuthParameters{
			BasicAuthParameters: &eventbridge.BasicAuthParameters{Username: "u", Password: "p"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.deleteConnectionCore(ctx, store, "del-conn"); err != nil {
		t.Fatal(err)
	}
	if len(fake.deletedARNs()) != 1 || fake.deletedARNs()[0] != connection.SecretArn {
		t.Fatalf("deleting the connection must remove its credential secret, deleted %v", fake.deletedARNs())
	}
}

// TestConnectionCreateRequiresSecretsManager pins the dependency error: a
// service without the Secrets Manager invoker cannot create credentials
// (AWS likewise requires Secrets Manager permission for connections).
func TestConnectionCreateRequiresSecretsManager(t *testing.T) {
	ctx := context.Background()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventbridge.NewEventsStore(st, "000000000000", "us-east-1")
	svc := NewEventsService(nil, "000000000000")
	t.Cleanup(svc.Close)

	_, err = svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "no-bus",
		AuthorizationType: "BASIC",
		AuthParameters: &eventbridge.AuthParameters{
			BasicAuthParameters: &eventbridge.BasicAuthParameters{Username: "u", Password: "p"},
		},
	})
	if err == nil {
		t.Fatal("connection creation must fail without the Secrets Manager service")
	}
	if _, err := store.GetConnection(ctx, "no-bus"); err != eventbridge.ErrConnectionNotFound {
		t.Fatalf("no record may be written when the secret cannot be stored, got err %v", err)
	}
}
