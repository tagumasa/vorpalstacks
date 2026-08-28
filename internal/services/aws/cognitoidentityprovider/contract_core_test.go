package cognitoidentityprovider

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// The typed AWS SDK validates the required members of
// UpdateUserPoolReplica, CreateTerms and CreateManagedLoginBranding
// client-side, so the server-side rejections below are unreachable through
// the SDK client; these unit tests pin the Core behaviour directly.

func newContractTestService(t *testing.T) (*CognitoService, *request.RequestContext, cognitostore.CognitoStoreInterface) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	store := cognitostore.NewCognitoStore(st, "123456789012", "us-east-1")

	svc := NewCognitoService("123456789012", "us-east-1")
	svc.stores.Store("us-east-1", store)
	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	return svc, reqCtx, store
}

func TestUpdateUserPoolReplicaRequiresStatus(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("replica-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	replica := &cognitostore.UserPoolReplica{
		UserPoolID:   pool.ID,
		RegionName:   "us-west-2",
		Status:       "ACTIVE",
		Role:         "SECONDARY",
		CreationDate: time.Now().UTC(),
	}
	if err := store.SaveUserPoolReplica(replica); err != nil {
		t.Fatalf("save replica: %v", err)
	}

	if _, err := svc.updateUserPoolReplicaCore(reqCtx, UpdateUserPoolReplicaInput{
		UserPoolID: pool.ID,
		RegionName: "us-west-2",
	}); err == nil {
		t.Fatal("update without Status succeeded")
	}
}

func TestCreateTermsRequiresAllModelMembers(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("terms-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	cases := []struct {
		name string
		in   CreateTermsInput
	}{
		{"missing ClientId", CreateTermsInput{UserPoolID: pool.ID, TermsName: "terms-of-use", TermsSource: "LINK", Enforcement: "NONE"}},
		{"missing TermsSource", CreateTermsInput{UserPoolID: pool.ID, ClientID: "client", TermsName: "terms-of-use", Enforcement: "NONE"}},
		{"missing Enforcement", CreateTermsInput{UserPoolID: pool.ID, ClientID: "client", TermsName: "terms-of-use", TermsSource: "LINK"}},
		{"off-enum TermsSource", CreateTermsInput{UserPoolID: pool.ID, ClientID: "client", TermsName: "terms-of-use", TermsSource: "DOCUMENT", Enforcement: "NONE"}},
		{"off-enum Enforcement", CreateTermsInput{UserPoolID: pool.ID, ClientID: "client", TermsName: "terms-of-use", TermsSource: "LINK", Enforcement: "REQUIRED"}},
	}
	for _, tc := range cases {
		if _, err := svc.createTermsCore(reqCtx, tc.in); err == nil {
			t.Fatalf("%s: create succeeded", tc.name)
		}
	}
}

func TestCreateManagedLoginBrandingRequiresClientId(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("branding-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if _, err := svc.createManagedLoginBrandingCore(reqCtx, CreateManagedLoginBrandingInput{
		UserPoolID: pool.ID,
	}); err == nil {
		t.Fatal("create without ClientId succeeded")
	}
}
