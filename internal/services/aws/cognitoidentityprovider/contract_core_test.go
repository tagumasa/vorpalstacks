package cognitoidentityprovider

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cognitoidentityprovider"
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

// The whole-configuration validation of a user pool client lives in the
// create/update Cores, so an invalid member must be rejected by the Core no
// matter which transport assembled the client.
func TestCreateUserPoolClientCoreRejectsInvalidConfig(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("client-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	client := cognitostore.NewUserPoolClient(pool.ID, "cfg-client")
	client.AccessTokenValidity = 86401
	if _, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, ""); err == nil {
		t.Fatal("client with out-of-range AccessTokenValidity accepted")
	}

	client = cognitostore.NewUserPoolClient(pool.ID, "cfg-client")
	client.ExplicitAuthFlows = []string{"ALLOW_NOT_A_FLOW"}
	if _, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, ""); err == nil {
		t.Fatal("client with off-enum ExplicitAuthFlows member accepted")
	}

	client = cognitostore.NewUserPoolClient(pool.ID, "cfg-client")
	client.TokenValidityUnits = &cognitostore.TokenValidityUnits{AccessToken: "fortnights"}
	if _, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, ""); err == nil {
		t.Fatal("client with off-enum TokenValidityUnits member accepted")
	}

	client = cognitostore.NewUserPoolClient(pool.ID, "cfg-client")
	if _, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, ""); err != nil {
		t.Fatalf("client with default configuration rejected: %v", err)
	}
}

// The create-only ClientSecret member replaces the generated secret, is
// bound by the ClientSecretType constraints, and cannot be combined with
// GenerateSecret.
func TestCreateUserPoolClientCoreCustomSecret(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("custom-secret-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	const customSecret = "customsecretvalue24charslong_min_ok"

	client := cognitostore.NewUserPoolClient(pool.ID, "secret-client")
	created, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, customSecret)
	if err != nil {
		t.Fatalf("create with custom secret: %v", err)
	}
	if created.ClientSecret != customSecret {
		t.Fatalf("custom secret not stored: got %q", created.ClientSecret)
	}
	stored, err := store.GetUserPoolClient(pool.ID, created.ClientID)
	if err != nil {
		t.Fatalf("load stored client: %v", err)
	}
	if stored.ClientSecret != customSecret {
		t.Fatalf("custom secret not persisted: got %q", stored.ClientSecret)
	}

	client = cognitostore.NewUserPoolClient(pool.ID, "secret-both")
	client.GenerateSecret = true
	if _, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, customSecret); err == nil {
		t.Fatal("custom secret combined with GenerateSecret accepted")
	}

	for _, bad := range []string{"short", "padding-padding-padding-padding-padding-padding-x!!"} {
		client = cognitostore.NewUserPoolClient(pool.ID, "secret-bad")
		if _, err := svc.createUserPoolClientCore(reqCtx.GetRegion(), client, bad); err == nil {
			t.Fatalf("secret violating ClientSecretType accepted: %q", bad)
		}
	}
}

// The update Core rejects an invalid assembled configuration and leaves the
// stored client untouched.
func TestUpdateUserPoolClientCoreRejectsInvalidConfig(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("client-update-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	client := cognitostore.NewUserPoolClient(pool.ID, "cfg-client")
	if err := store.CreateUserPoolClient(client); err != nil {
		t.Fatalf("create client: %v", err)
	}

	client.PreventUserExistenceErrors = "OVERRIDE"
	if err := svc.updateUserPoolClientCore(reqCtx.GetRegion(), client); err == nil {
		t.Fatal("update with off-enum PreventUserExistenceErrors accepted")
	}
	stored, err := store.GetUserPoolClient(pool.ID, client.ClientID)
	if err != nil {
		t.Fatalf("get client: %v", err)
	}
	if stored.PreventUserExistenceErrors == "OVERRIDE" {
		t.Fatal("invalid value persisted")
	}
}

// The boolean client members are presence-aware on the shared write path:
// an omitted member keeps the stored value, EnableTokenRevocation is
// active on new clients unless explicitly deactivated, GenerateSecret
// exists only on the create request (an update wire carrying it must not
// strip the stored secret pair), and the secret is returned by every
// response whose client carries one.
func TestUserPoolClientBooleanPresenceAndSecretRules(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("client-bool-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	clientResponse := func(resp interface{}) map[string]interface{} {
		t.Helper()
		m, ok := resp.(map[string]interface{})["UserPoolClient"].(map[string]interface{})
		if !ok {
			t.Fatalf("response carries no UserPoolClient map: %T", resp)
		}
		return m
	}

	createResp, err := svc.CreateUserPoolClient(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolId": pool.ID,
		"ClientName": "bool-client",
	}})
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	created := clientResponse(createResp)
	if created["EnableTokenRevocation"] != true {
		t.Fatalf("EnableTokenRevocation default: got %v, want true", created["EnableTokenRevocation"])
	}
	if created["AllowedOAuthFlowsUserPoolClient"] != false {
		t.Fatalf("AllowedOAuthFlowsUserPoolClient default: got %v, want false", created["AllowedOAuthFlowsUserPoolClient"])
	}
	if _, has := created["ClientSecret"]; has {
		t.Fatal("secret returned for a client created without GenerateSecret")
	}
	clientID, _ := created["ClientId"].(string)

	// An update that omits the boolean members keeps every stored value.
	updateReq := func(params map[string]interface{}) *request.ParsedRequest {
		params["UserPoolId"] = pool.ID
		params["ClientId"] = clientID
		return &request.ParsedRequest{Parameters: params}
	}
	renamed, err := svc.UpdateUserPoolClient(context.Background(), reqCtx, updateReq(map[string]interface{}{
		"ClientName": "bool-client-renamed",
	}))
	if err != nil {
		t.Fatalf("rename update: %v", err)
	}
	if renamed := clientResponse(renamed); renamed["EnableTokenRevocation"] != true {
		t.Fatalf("omitted EnableTokenRevocation flipped: got %v, want true", renamed["EnableTokenRevocation"])
	}

	// An explicit value applies, and a later omission keeps the applied
	// value — the create-time default never resurfaces on update.
	disabled, err := svc.UpdateUserPoolClient(context.Background(), reqCtx, updateReq(map[string]interface{}{
		"ClientName":                      "bool-client-disabled",
		"EnableTokenRevocation":           false,
		"AllowedOAuthFlowsUserPoolClient": true,
	}))
	if err != nil {
		t.Fatalf("explicit-boolean update: %v", err)
	}
	disabledClient := clientResponse(disabled)
	if disabledClient["EnableTokenRevocation"] != false {
		t.Fatalf("explicit EnableTokenRevocation=false not applied: got %v", disabledClient["EnableTokenRevocation"])
	}
	if disabledClient["AllowedOAuthFlowsUserPoolClient"] != true {
		t.Fatalf("explicit AllowedOAuthFlowsUserPoolClient=true not applied: got %v", disabledClient["AllowedOAuthFlowsUserPoolClient"])
	}
	kept, err := svc.UpdateUserPoolClient(context.Background(), reqCtx, updateReq(map[string]interface{}{
		"ClientName": "bool-client-kept",
	}))
	if err != nil {
		t.Fatalf("post-explicit rename update: %v", err)
	}
	if keptClient := clientResponse(kept); keptClient["EnableTokenRevocation"] != false || keptClient["AllowedOAuthFlowsUserPoolClient"] != true {
		t.Fatalf("explicit booleans not kept after omission: revocation=%v flows=%v",
			keptClient["EnableTokenRevocation"], keptClient["AllowedOAuthFlowsUserPoolClient"])
	}

	// GenerateSecret is a create-only member: the stored flag and secret
	// survive an update wire that carries the member, and the secret is
	// visible on describe exactly as on create and update.
	secretCreate, err := svc.CreateUserPoolClient(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolId":     pool.ID,
		"ClientName":     "secret-client",
		"GenerateSecret": true,
	}})
	if err != nil {
		t.Fatalf("create secret client: %v", err)
	}
	secretClient := clientResponse(secretCreate)
	secret, _ := secretClient["ClientSecret"].(string)
	if secret == "" {
		t.Fatal("GenerateSecret=true produced no client secret")
	}
	secretID, _ := secretClient["ClientId"].(string)

	describeResp, err := svc.DescribeUserPoolClient(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolId": pool.ID,
		"ClientId":   secretID,
	}})
	if err != nil {
		t.Fatalf("describe secret client: %v", err)
	}
	if described := clientResponse(describeResp); described["ClientSecret"] != secret {
		t.Fatalf("describe returned a different secret: got %v, want the stored one", described["ClientSecret"])
	}

	secretUpdate, err := svc.UpdateUserPoolClient(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolId":     pool.ID,
		"ClientId":       secretID,
		"ClientName":     "secret-client-renamed",
		"GenerateSecret": false,
	}})
	if err != nil {
		t.Fatalf("update secret client: %v", err)
	}
	if updated := clientResponse(secretUpdate); updated["ClientSecret"] != secret {
		t.Fatalf("update altered the client secret: got %v, want the stored one", updated["ClientSecret"])
	}
	storedSecret, err := store.GetUserPoolClient(pool.ID, secretID)
	if err != nil {
		t.Fatalf("get secret client: %v", err)
	}
	if !storedSecret.GenerateSecret || storedSecret.ClientSecret != secret {
		t.Fatalf("update stripped the stored secret pair: GenerateSecret=%v secret-matches=%v",
			storedSecret.GenerateSecret, storedSecret.ClientSecret == secret)
	}
}

// UpdateIdentityProvider keeps absent members, replaces present ones, cannot
// retarget the provider type (the update request model has no ProviderType
// member), and an explicitly empty identifier list clears the stored one.
func TestUpdateIdentityProviderCorePresenceSemantics(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("idp-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	ip := &cognitostore.IdentityProvider{
		UserPoolID:       pool.ID,
		ProviderName:     "oidc-main",
		ProviderType:     "OIDC",
		ProviderDetails:  map[string]string{"client_id": "a", "client_secret": "b"},
		AttributeMapping: map[string]string{"email": "email"},
		IdpIdentifiers:   []string{"alias-1"},
	}
	if err := store.CreateIdentityProvider(ip); err != nil {
		t.Fatalf("create provider: %v", err)
	}

	updated, err := svc.updateIdentityProviderCore(reqCtx.GetRegion(), UpdateIdentityProviderInput{
		UserPoolID:      pool.ID,
		ProviderName:    "oidc-main",
		ProviderDetails: map[string]string{"client_id": "c"},
	})
	if err != nil {
		t.Fatalf("update provider: %v", err)
	}
	if updated.ProviderDetails["client_id"] != "c" || len(updated.ProviderDetails) != 1 {
		t.Fatalf("ProviderDetails not replaced: %v", updated.ProviderDetails)
	}
	if updated.AttributeMapping["email"] != "email" {
		t.Fatalf("absent AttributeMapping not kept: %v", updated.AttributeMapping)
	}
	if len(updated.IdpIdentifiers) != 1 || updated.IdpIdentifiers[0] != "alias-1" {
		t.Fatalf("absent IdpIdentifiers not kept: %v", updated.IdpIdentifiers)
	}
	if updated.ProviderType != "OIDC" {
		t.Fatalf("update retargeted the provider type: %q", updated.ProviderType)
	}

	updated, err = svc.updateIdentityProviderCore(reqCtx.GetRegion(), UpdateIdentityProviderInput{
		UserPoolID:             pool.ID,
		ProviderName:           "oidc-main",
		IdpIdentifiersProvided: true,
	})
	if err != nil {
		t.Fatalf("clear identifiers: %v", err)
	}
	if len(updated.IdpIdentifiers) != 0 {
		t.Fatalf("explicit empty identifier list did not clear: %v", updated.IdpIdentifiers)
	}

	if _, err := svc.updateIdentityProviderCore(reqCtx.GetRegion(), UpdateIdentityProviderInput{
		UserPoolID:   pool.ID,
		ProviderName: "ghost-provider",
	}); err != ErrResourceNotFound {
		t.Fatalf("update of unknown provider returned %v, want ErrResourceNotFound", err)
	}
}

// Group create validation (required members, name pattern, precedence range)
// lives in the single create Core shared by both planes.
func TestCreateGroupValidatedCoreRejectsInvalidMembers(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("group-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if _, err := svc.createGroupValidatedCore(context.Background(), reqCtx.GetRegion(), CreateGroupInput{
		UserPoolID: pool.ID,
		GroupName:  "bad group",
	}); err == nil {
		t.Fatal("group name violating the GroupNameType pattern accepted")
	}

	negative := -1
	if _, err := svc.createGroupValidatedCore(context.Background(), reqCtx.GetRegion(), CreateGroupInput{
		UserPoolID: pool.ID,
		GroupName:  "team-a",
		Precedence: &negative,
	}); err == nil {
		t.Fatal("negative precedence accepted")
	}

	group, err := svc.createGroupValidatedCore(context.Background(), reqCtx.GetRegion(), CreateGroupInput{
		UserPoolID:  pool.ID,
		GroupName:   "team-a",
		Description: "first line of defence",
	})
	if err != nil {
		t.Fatalf("valid group rejected: %v", err)
	}
	if group.Description != "first line of defence" {
		t.Fatalf("description not applied: %q", group.Description)
	}
	if _, err := store.GetGroup(pool.ID, "team-a"); err != nil {
		t.Fatalf("group not persisted: %v", err)
	}
}

// Group update keeps absent members, applies present ones (an explicit empty
// description clears), validates the precedence range in the Core, and
// reports a dedicated not-found for an unknown group.
func TestUpdateGroupCoreValidatesAndMerges(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("group-update-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	precedence := 5
	if _, err := svc.createGroupValidatedCore(context.Background(), reqCtx.GetRegion(), CreateGroupInput{
		UserPoolID:  pool.ID,
		GroupName:   "team-b",
		Description: "keep me",
		Precedence:  &precedence,
	}); err != nil {
		t.Fatalf("create group: %v", err)
	}

	empty := ""
	updated, err := svc.updateGroupCore(context.Background(), reqCtx.GetRegion(), UpdateGroupInput{
		UserPoolID:  pool.ID,
		GroupName:   "team-b",
		Description: &empty,
	})
	if err != nil {
		t.Fatalf("clear description: %v", err)
	}
	if updated.Description != "" {
		t.Fatalf("explicit empty description did not clear: %q", updated.Description)
	}
	if updated.Precedence == nil || *updated.Precedence != 5 {
		t.Fatalf("absent precedence not kept: %v", updated.Precedence)
	}

	negative := -3
	if _, err := svc.updateGroupCore(context.Background(), reqCtx.GetRegion(), UpdateGroupInput{
		UserPoolID: pool.ID,
		GroupName:  "team-b",
		Precedence: &negative,
	}); err == nil {
		t.Fatal("negative precedence accepted on update")
	}
	stored, err := store.GetGroup(pool.ID, "team-b")
	if err != nil {
		t.Fatalf("get group: %v", err)
	}
	if stored.Precedence == nil || *stored.Precedence != 5 {
		t.Fatalf("stored precedence changed by rejected update: %v", stored.Precedence)
	}

	if _, err := svc.updateGroupCore(context.Background(), reqCtx.GetRegion(), UpdateGroupInput{
		UserPoolID: pool.ID,
		GroupName:  "ghost-group",
	}); err != ErrGroupNotFound {
		t.Fatalf("update of unknown group returned %v, want ErrGroupNotFound", err)
	}
}

// The tag Cores resolve the ResourceArn to an existing user pool: a
// well-formed pool ARN naming no pool fails with ResourceNotFoundException on
// all three operations, and a structurally wrong ARN (other service, other
// resource type, or pattern violation) is a malformed parameter.
func TestTagCoresResolveResourceToExistingUserPool(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("tag-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	ghost := "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_nosuchpoo"
	if err := svc.tagResourceCore(reqCtx.GetRegion(), ghost, map[string]string{"Environment": "test"}); err != ErrResourceNotFound {
		t.Fatalf("tag of nonexistent pool returned %v, want ErrResourceNotFound", err)
	}
	if err := svc.untagResourceCore(reqCtx.GetRegion(), ghost, []string{"Environment"}); err != ErrResourceNotFound {
		t.Fatalf("untag of nonexistent pool returned %v, want ErrResourceNotFound", err)
	}
	if _, err := svc.listTagsForResourceCore(reqCtx.GetRegion(), ghost); err != ErrResourceNotFound {
		t.Fatalf("list tags of nonexistent pool returned %v, want ErrResourceNotFound", err)
	}

	for _, bad := range []string{
		"arn:aws:cognito-identity:us-east-1:123456789012:identitypool/us-east-1_x",
		"arn:aws:cognito-idp:us-east-1:123456789012:group/us-east-1_x/team",
		"not-an-arn-but-longer-than-twenty-chars",
	} {
		if err := svc.tagResourceCore(reqCtx.GetRegion(), bad, map[string]string{"Environment": "test"}); err != ErrInvalidParameter {
			t.Fatalf("tag with malformed ARN %q returned %v, want ErrInvalidParameter", bad, err)
		}
	}

	if err := svc.tagResourceCore(reqCtx.GetRegion(), pool.Arn, map[string]string{"Environment": "test"}); err != nil {
		t.Fatalf("tag of real pool: %v", err)
	}
	tags, err := svc.listTagsForResourceCore(reqCtx.GetRegion(), pool.Arn)
	if err != nil {
		t.Fatalf("list tags of real pool: %v", err)
	}
	if tags["Environment"] != "test" {
		t.Fatalf("expected Environment=test, got %v", tags)
	}
	if err := svc.untagResourceCore(reqCtx.GetRegion(), pool.Arn, []string{"Environment"}); err != nil {
		t.Fatalf("untag of real pool: %v", err)
	}
	tags, err = svc.listTagsForResourceCore(reqCtx.GetRegion(), pool.Arn)
	if err != nil {
		t.Fatalf("list tags after untag: %v", err)
	}
	if _, ok := tags["Environment"]; ok {
		t.Fatalf("Environment tag survived untag: %v", tags)
	}
}

// The domain Cores enforce the binding rules on the wire: duplicate
// bindings surface as InvalidParameterException (the domain family models
// no conflict shape), and delete/update against another pool's domain
// surface as ResourceNotFoundException.
func TestUserPoolDomainCoreBindingRules(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	poolA, err := store.CreateUserPool(cognitostore.NewUserPool("domain-core-a", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool A: %v", err)
	}
	poolB, err := store.CreateUserPool(cognitostore.NewUserPool("domain-core-b", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool B: %v", err)
	}

	if _, err := svc.createUserPoolDomainCore(CreateUserPoolDomainInput{
		Region: reqCtx.GetRegion(), Domain: "core-alpha", UserPoolID: poolA.ID,
	}); err != nil {
		t.Fatalf("initial domain bind: %v", err)
	}

	// Claiming pool A's domain for pool B is a parameter violation.
	if _, err := svc.createUserPoolDomainCore(CreateUserPoolDomainInput{
		Region: reqCtx.GetRegion(), Domain: "core-alpha", UserPoolID: poolB.ID,
	}); err != ErrInvalidParameter {
		t.Fatalf("cross-pool domain claim returned %v, want ErrInvalidParameter", err)
	}

	// A second domain for pool A is a parameter violation too.
	if _, err := svc.createUserPoolDomainCore(CreateUserPoolDomainInput{
		Region: reqCtx.GetRegion(), Domain: "core-beta", UserPoolID: poolA.ID,
	}); err != ErrInvalidParameter {
		t.Fatalf("second domain for one pool returned %v, want ErrInvalidParameter", err)
	}

	// Updating another pool's domain is reported as not found.
	if _, err := svc.updateUserPoolDomainCore(UpdateUserPoolDomainInput{
		Region: reqCtx.GetRegion(), Domain: "core-alpha", UserPoolID: poolB.ID,
	}); err != ErrResourceNotFound {
		t.Fatalf("foreign update returned %v, want ErrResourceNotFound", err)
	}

	// The owner updates its own domain.
	if _, err := svc.updateUserPoolDomainCore(UpdateUserPoolDomainInput{
		Region: reqCtx.GetRegion(), Domain: "core-alpha", UserPoolID: poolA.ID,
	}); err != nil {
		t.Fatalf("owning update: %v", err)
	}

	// Deleting another pool's domain fails and leaves the binding intact.
	if err := svc.deleteUserPoolDomainCore(reqCtx.GetRegion(), poolB.ID, "core-alpha"); err != ErrResourceNotFound {
		t.Fatalf("foreign delete returned %v, want ErrResourceNotFound", err)
	}
	if _, err := store.GetUserPoolDomain("core-alpha"); err != nil {
		t.Fatalf("foreign delete removed the binding: %v", err)
	}

	// Deleting a domain that does not exist, from a pool that does not
	// exist: both are reported as not found.
	if err := svc.deleteUserPoolDomainCore(reqCtx.GetRegion(), poolA.ID, "no-such-domain"); err != ErrResourceNotFound {
		t.Fatalf("missing domain delete returned %v, want ErrResourceNotFound", err)
	}
	if err := svc.deleteUserPoolDomainCore(reqCtx.GetRegion(), "us-east-1_ghostpool", "core-alpha"); err != ErrResourceNotFound {
		t.Fatalf("ghost pool delete returned %v, want ErrResourceNotFound", err)
	}

	// The owner deletes its own domain.
	if err := svc.deleteUserPoolDomainCore(reqCtx.GetRegion(), poolA.ID, "core-alpha"); err != nil {
		t.Fatalf("owning delete: %v", err)
	}
}

// Provisioned limits answer the documented free quota until a value is
// provisioned, persist the provisioned value in the regional store, and
// reject definitions outside the model's enum and adjustable-category set
// and values below the free quota.
func TestProvisionedLimitCoreRules(t *testing.T) {
	svc, reqCtx, _ := newContractTestService(t)

	// The typed AWS SDK validates the required LimitDefinition and
	// RequestedLimitValue client-side, so the server-side rejections below
	// are unreachable through the SDK client; the Core is pinned directly.
	unset, err := svc.getProvisionedLimitCore(ProvisionedLimitInput{
		Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserAuthentication",
	})
	if err != nil {
		t.Fatalf("get unset limit: %v", err)
	}
	if unset.ProvisionedLimitValue != 120 || unset.FreeLimitValue != 120 {
		t.Fatalf("unset UserAuthentication answers %d/%d, want 120/120", unset.ProvisionedLimitValue, unset.FreeLimitValue)
	}

	updated, err := svc.updateProvisionedLimitCore(ProvisionedLimitInput{
		Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserAuthentication",
		RequestedValue: 300, RequestedValueSet: true,
	})
	if err != nil {
		t.Fatalf("update limit: %v", err)
	}
	if updated.ProvisionedLimitValue != 300 || updated.FreeLimitValue != 120 {
		t.Fatalf("update answers %d/%d, want 300/120", updated.ProvisionedLimitValue, updated.FreeLimitValue)
	}

	got, err := svc.getProvisionedLimitCore(ProvisionedLimitInput{
		Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserAuthentication",
	})
	if err != nil {
		t.Fatalf("get provisioned limit: %v", err)
	}
	if got.ProvisionedLimitValue != 300 || got.FreeLimitValue != 120 {
		t.Fatalf("provisioned value not persisted: got %d/%d", got.ProvisionedLimitValue, got.FreeLimitValue)
	}

	// Reducing back to the free quota is the documented lower bound.
	reduced, err := svc.updateProvisionedLimitCore(ProvisionedLimitInput{
		Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserAuthentication",
		RequestedValue: 120, RequestedValueSet: true,
	})
	if err != nil {
		t.Fatalf("reduce to free quota: %v", err)
	}
	if reduced.ProvisionedLimitValue != 120 {
		t.Fatalf("reduction answers %d, want 120", reduced.ProvisionedLimitValue)
	}

	// Definition-level violations reject both operations identically.
	definitionCases := []struct {
		name string
		in   ProvisionedLimitInput
		want error
	}{
		{
			"off-enum limit class",
			ProvisionedLimitInput{Region: reqCtx.GetRegion(), LimitClass: "ACCOUNT_LEVEL", Category: "UserAuthentication"},
			ErrInvalidParameter,
		},
		{
			"missing category",
			ProvisionedLimitInput{Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY"},
			ErrInvalidParameter,
		},
		{
			"unknown category",
			ProvisionedLimitInput{Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "NopeLimit"},
			ErrResourceNotFound,
		},
		{
			"non-adjustable category",
			ProvisionedLimitInput{Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserList"},
			ErrResourceNotFound,
		},
	}
	for _, tc := range definitionCases {
		if _, err := svc.updateProvisionedLimitCore(tc.in); err != tc.want {
			t.Fatalf("%s: update returned %v, want %v", tc.name, err, tc.want)
		}
		if _, err := svc.getProvisionedLimitCore(tc.in); err != tc.want {
			t.Fatalf("%s: get returned %v, want %v", tc.name, err, tc.want)
		}
	}

	// Update-specific violations: a missing or below-free requested value.
	updateCases := []struct {
		name string
		in   ProvisionedLimitInput
	}{
		{
			"missing RequestedLimitValue",
			ProvisionedLimitInput{Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserAuthentication"},
		},
		{
			"value below free quota",
			ProvisionedLimitInput{Region: reqCtx.GetRegion(), LimitClass: "API_CATEGORY", Category: "UserAuthentication", RequestedValue: 10, RequestedValueSet: true},
		},
	}
	for _, tc := range updateCases {
		if _, err := svc.updateProvisionedLimitCore(tc.in); err != ErrInvalidParameter {
			t.Fatalf("%s: update returned %v, want %v", tc.name, err, ErrInvalidParameter)
		}
	}
}

// The create operation's modelled error surface defines no duplicate-name
// rejection for app clients: two clients bearing the same name in one pool
// coexist, each with its own generated identifier.
func TestCreateUserPoolClientAllowsDuplicateNames(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("dup-client-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	create := func(name string) string {
		t.Helper()
		resp, err := svc.CreateUserPoolClient(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"UserPoolId": pool.ID,
			"ClientName": name,
		}})
		if err != nil {
			t.Fatalf("create client %q: %v", name, err)
		}
		client, ok := resp.(map[string]interface{})["UserPoolClient"].(map[string]interface{})
		if !ok {
			t.Fatalf("response carries no UserPoolClient map: %T", resp)
		}
		id, _ := client["ClientId"].(string)
		if id == "" {
			t.Fatal("create response carries no ClientId")
		}
		return id
	}
	first := create("shared-name")
	second := create("shared-name")
	if first == second {
		t.Fatalf("duplicate-name clients share one identifier: %q", first)
	}
}

// The pool Status member carries the model's StatusType enum value Enabled;
// the admin converter maps that value onto the proto enum so the console
// renders every pool as enabled rather than disabled.
func TestUserPoolStatusIsEnabledEnum(t *testing.T) {
	_, _, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("status-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if pool.Status != cognitostore.UserPoolStatusEnabled {
		t.Fatalf("stored pool Status: got %q, want %q", pool.Status, cognitostore.UserPoolStatusEnabled)
	}
	if got := formatUserPool(pool)["Status"]; got != cognitostore.UserPoolStatusEnabled {
		t.Fatalf("described pool Status: got %v, want %q", got, cognitostore.UserPoolStatusEnabled)
	}
	if statusToProto(pool.Status) != pb.StatusType_STATUS_TYPE_ENABLED {
		t.Fatal("admin converter does not map the stored status to the enabled proto enum")
	}
}

// DeleteManagedLoginBranding and DeleteTerms report a missing record as
// ResourceNotFoundException — a silent success for an absent record is not
// on either operation's error surface.
func TestDeleteManagedLoginBrandingUnknownStyleIsNotFound(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("branding-del-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if _, err := svc.deleteManagedLoginBrandingCore(reqCtx, DeleteManagedLoginBrandingInput{
		UserPoolID:             pool.ID,
		ManagedLoginBrandingID: "6f9619ff-8b86-d011-b42d-00c04fc964ff",
	}); err != ErrResourceNotFound {
		t.Fatalf("delete unknown branding: got %v, want ErrResourceNotFound", err)
	}
}

func TestDeleteTermsUnknownDocumentIsNotFound(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("terms-del-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if _, err := svc.deleteTermsCore(reqCtx, DeleteTermsInput{
		UserPoolID: pool.ID,
		TermsID:    "6f9619ff-8b86-d011-b42d-00c04fc964ff",
	}); err != ErrResourceNotFound {
		t.Fatalf("delete unknown terms: got %v, want ErrResourceNotFound", err)
	}
}

// GetSigningCertificate answers with a self-signed X.509 certificate derived
// from the pool's signing key — not the bare public key — and the derivation
// is deterministic: the same pool always yields the identical certificate,
// valid for the documented ten years.
func TestGetSigningCertificateIssuesStableX509Certificate(t *testing.T) {
	svc, _, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("signcert-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	first, err := svc.getSigningCertificateCore("us-east-1", pool.ID)
	if err != nil {
		t.Fatalf("get signing certificate: %v", err)
	}
	second, err := svc.getSigningCertificateCore("us-east-1", pool.ID)
	if err != nil {
		t.Fatalf("get signing certificate again: %v", err)
	}
	if first != second {
		t.Fatal("certificate is not stable across calls")
	}
	block, _ := pem.Decode([]byte(first))
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("response is not a PEM certificate: %.40s", first)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
	}
	if cert.Subject.CommonName != pool.ID {
		t.Fatalf("certificate CN: got %q, want %q", cert.Subject.CommonName, pool.ID)
	}
	if !cert.NotAfter.Equal(cert.NotBefore.AddDate(signingCertificateValidityYears, 0, 0)) {
		t.Fatalf("validity span: NotBefore %v NotAfter %v", cert.NotBefore, cert.NotAfter)
	}
	// The certificate is a leaf-style signing credential without CA
	// capability, so the self-signature is verified directly over the
	// to-be-signed bytes rather than through CA-chain verification.
	if err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature); err != nil {
		t.Fatalf("certificate is not signed by its own key: %v", err)
	}
}

// Explicitly supplied token-validity values out of the model ranges are
// rejected — an explicit 0 or a value above the maximum reaches the Core's
// validation instead of being silently dropped.
func TestUserPoolClientTokenValidityRanges(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("validity-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	created, err := svc.CreateUserPoolClient(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolId": pool.ID,
		"ClientName": "validity-client",
	}})
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	clientID, _ := created.(map[string]interface{})["UserPoolClient"].(map[string]interface{})["ClientId"].(string)
	updateReq := func(validity string, value int) *request.ParsedRequest {
		return &request.ParsedRequest{Parameters: map[string]interface{}{
			"UserPoolId": pool.ID,
			"ClientId":   clientID,
			"ClientName": "validity-client",
			validity:     value,
		}}
	}

	for _, tc := range []struct {
		member string
		value  int
	}{
		{"AccessTokenValidity", 0},
		{"AccessTokenValidity", 86401},
		{"IdTokenValidity", 0},
		{"IdTokenValidity", 86401},
		{"RefreshTokenValidity", -1},
	} {
		if _, err := svc.UpdateUserPoolClient(context.Background(), reqCtx, updateReq(tc.member, tc.value)); err != ErrInvalidParameter {
			t.Fatalf("%s=%d returned %v, want InvalidParameterException", tc.member, tc.value, err)
		}
	}

	if _, err := svc.UpdateUserPoolClient(context.Background(), reqCtx, updateReq("AccessTokenValidity", 3600)); err != nil {
		t.Fatalf("in-range AccessTokenValidity rejected: %v", err)
	}
}
