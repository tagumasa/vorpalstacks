package cognitoidentity

import (
	"context"
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// capturingIssuer records the role and the session tags each credential
// session was issued for.
type capturingIssuer struct {
	roleArn string
	tags    map[string]string
}

func (c *capturingIssuer) IssueSession(roleArn, roleSessionName string, durationSeconds int, tags map[string]string) (*CredentialResult, error) {
	c.roleArn = roleArn
	c.tags = tags
	return &CredentialResult{
		AccessKeyID:     "AKIA-TEST",
		SecretAccessKey: "secret",
		SessionToken:    "session",
		Expiration:      time.Now().Add(time.Hour),
	}, nil
}

// stubIDTokenClaims stands in for the user-pool service behind the
// CognitoIDTokenClaimResolver contract: it hands each presented token the
// claim set the test wants the token to carry, and can refuse everything.
type stubIDTokenClaims struct {
	byToken map[string]map[string]string
	err     error
}

func (s *stubIDTokenClaims) IDTokenClaimsForPool(_ context.Context, _, _, token string) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.byToken[token], nil
}

const (
	testAuthRole   = "arn:aws:iam::000000000000:role/auth-role"
	testUnauthRole = "arn:aws:iam::000000000000:role/unauth-role"
	testMappedRole = "arn:aws:iam::000000000000:role/mapped-role"
	testOtherRole  = "arn:aws:iam::000000000000:role/other-role"

	// testUserPoolProvider is the platform user-pool issuer form a Logins key
	// carries when the login is a user-pool ID token.
	testUserPoolProvider = "cognito-idp.us-east-1.amazonaws.com/us-east-1_testpool"
)

// seedRolePool creates a pool with classic roles, an engineering-claim login
// identity and the given role mappings for the login provider. Claims arrive
// through the stub resolver the caller installs.
func seedRolePool(t *testing.T, store cognitoidentitystore.CognitoIdentityStoreInterface, mappings map[string]cognitoidentitystore.RoleMapping) string {
	t.Helper()
	pool := cognitoidentitystore.NewIdentityPool("role-pool", false, "us-east-1")
	created, err := store.CreateIdentityPool(pool)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := store.SetIdentityPoolRoles(created.ID, testAuthRole, testUnauthRole, mappings); err != nil {
		t.Fatalf("set pool roles: %v", err)
	}
	identity, err := store.GetOrCreateIdentityByLogins(created.ID, map[string]string{testUserPoolProvider: "user-token"})
	if err != nil {
		t.Fatalf("create login identity: %v", err)
	}
	return identity.ID
}

func rulesMapping(resolution string, rules ...cognitoidentitystore.MappingRule) map[string]cognitoidentitystore.RoleMapping {
	return map[string]cognitoidentitystore.RoleMapping{
		testUserPoolProvider: {
			Type:                    "Rules",
			AmbiguousRoleResolution: resolution,
			RulesConfiguration:      &cognitoidentitystore.RulesConfiguration{Rules: rules},
		},
	}
}

// A matching rule selects its role for the credential session.
func TestGetCredentialsForIdentityRulesMappingSelectsMappedRole(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"department": "engineering"},
	}})
	identityID := seedRolePool(t, svcStore(t, svc), rulesMapping("AuthenticatedRole", cognitoidentitystore.MappingRule{
		Claim: "department", MatchType: "Equals", Value: "engineering", RoleARN: testMappedRole,
	}))

	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatalf("GetCredentialsForIdentity through a matching rule: %v", err)
	}
	if issuer.roleArn != testMappedRole {
		t.Fatalf("session issued for %s, want the rule-mapped role %s", issuer.roleArn, testMappedRole)
	}
}

// A Token mapping reads the cognito:preferred_role claim of the presented
// user-pool ID token.
func TestGetCredentialsForIdentityTokenMappingUsesPreferredRole(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"cognito:preferred_role": testMappedRole},
	}})
	identityID := seedRolePool(t, svcStore(t, svc), map[string]cognitoidentitystore.RoleMapping{
		testUserPoolProvider: {Type: "Token", AmbiguousRoleResolution: "Deny"},
	})

	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatalf("GetCredentialsForIdentity through a preferred role: %v", err)
	}
	if issuer.roleArn != testMappedRole {
		t.Fatalf("session issued for %s, want the preferred role %s", issuer.roleArn, testMappedRole)
	}
}

// No rule matching plus a Deny ambiguous resolution refuses the session.
func TestGetCredentialsForIdentityDenyResolutionRefuses(t *testing.T) {
	svc, _ := newMergeTestService(t)
	svc.SetCredentialIssuer(&capturingIssuer{})
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"department": "engineering"},
	}})
	identityID := seedRolePool(t, svcStore(t, svc), rulesMapping("Deny", cognitoidentitystore.MappingRule{
		Claim: "department", MatchType: "Equals", Value: "sales", RoleARN: testMappedRole,
	}))

	_, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	})
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("unmatched rule with Deny resolution returned %v, want NotAuthorizedException", err)
	}
}

// A CustomRoleArn is ignored when the mapping grants exactly one role: the
// model documents the parameter as selecting "when multiple roles were
// received in the token", so a single grant has nothing to select among.
func TestGetCredentialsForIdentityCustomRoleArnIgnoredForSingleGrantedRole(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"department": "engineering"},
	}})
	identityID := seedRolePool(t, svcStore(t, svc), rulesMapping("AuthenticatedRole", cognitoidentitystore.MappingRule{
		Claim: "department", MatchType: "Equals", Value: "engineering", RoleARN: testMappedRole,
	}))

	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID:    identityID,
		CustomRoleARN: testOtherRole,
		Logins:        map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatalf("CustomRoleArn alongside a single granted role: %v", err)
	}
	if issuer.roleArn != testMappedRole {
		t.Fatalf("session issued for %s, want the single granted role %s", issuer.roleArn, testMappedRole)
	}
}

// A CustomRoleArn outside the roles a multi-role grant offers is refused.
func TestGetCredentialsForIdentityCustomRoleArnRequiresAuthorisation(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"cognito:roles": testMappedRole + "," + testOtherRole},
	}})
	store := svcStore(t, svc)

	// A Token mapping whose token carries several cognito:roles is the
	// multi-role grant CustomRoleArn selects among: a role from the list is
	// assumed, a role outside it is refused.
	multiID := seedRolePool(t, store, map[string]cognitoidentitystore.RoleMapping{
		testUserPoolProvider: {Type: "Token", AmbiguousRoleResolution: "AuthenticatedRole"},
	})
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID:    multiID,
		CustomRoleARN: testOtherRole,
		Logins:        map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatalf("CustomRoleArn selecting within the multi-role grant: %v", err)
	}
	if issuer.roleArn != testOtherRole {
		t.Fatalf("session issued for %s, want the selected role %s", issuer.roleArn, testOtherRole)
	}
	_, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID:    multiID,
		CustomRoleARN: "arn:aws:iam::000000000000:role/foreign-role",
		Logins:        map[string]string{testUserPoolProvider: "user-token"},
	})
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("CustomRoleArn outside the multi-role grant returned %v, want NotAuthorizedException", err)
	}
}

// Without a governing mapping the pool grants only its authenticated role —
// a single granted role, so a CustomRoleArn is ignored rather than enforced;
// the unauthenticated grant, with no token roles at all, ignores it the same
// way.
func TestGetCredentialsForIdentityCustomRoleArnIgnoredWithoutSelection(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"department": "engineering"},
	}})
	store := svcStore(t, svc)

	bareID := seedRolePool(t, store, nil)
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID:    bareID,
		CustomRoleARN: testOtherRole,
		Logins:        map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatalf("CustomRoleArn without a mapping: %v", err)
	}
	if issuer.roleArn != testAuthRole {
		t.Fatalf("session issued for %s, want the authenticated role %s", issuer.roleArn, testAuthRole)
	}

	unauth := cognitoidentitystore.NewIdentity(poolOf(t, store, bareID))
	if err := store.CreateIdentity(unauth); err != nil {
		t.Fatalf("create unauthenticated identity: %v", err)
	}
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID:    unauth.ID,
		CustomRoleARN: testUnauthRole,
	}); err != nil {
		t.Fatalf("CustomRoleArn for an unauthenticated identity: %v", err)
	}
	if issuer.roleArn != testUnauthRole {
		t.Fatalf("session issued for %s, want the unauthenticated role %s", issuer.roleArn, testUnauthRole)
	}
}

// Rules are evaluated in order and the first match specifies the role: a
// later rule that also matches neither replaces nor joins the grant.
func TestGetCredentialsForIdentityRulesFirstMatchWins(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"department": "engineering"},
	}})
	identityID := seedRolePool(t, svcStore(t, svc), rulesMapping("AuthenticatedRole",
		cognitoidentitystore.MappingRule{Claim: "department", MatchType: "Equals", Value: "engineering", RoleARN: testMappedRole},
		cognitoidentitystore.MappingRule{Claim: "department", MatchType: "Contains", Value: "gine", RoleARN: testOtherRole},
	))
	// department=engineering matches both rules; the first rule's role is
	// the grant.
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatalf("GetCredentialsForIdentity through matching rules: %v", err)
	}
	if issuer.roleArn != testMappedRole {
		t.Fatalf("session issued for %s, want the first matching rule's role %s", issuer.roleArn, testMappedRole)
	}
}

// With several linked providers each carrying a mapping, the governing
// mapping is chosen deterministically: providers are evaluated in sorted
// order, so the alphabetically-first linked provider's grant wins on every
// request, never the vagaries of map iteration.
func TestGetCredentialsForIdentityMultipleProvidersDeterministic(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	alphaProvider := "cognito-idp.us-east-1.amazonaws.com/us-east-1_alpha"
	betaProvider := "cognito-idp.us-east-1.amazonaws.com/us-east-1_beta"
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"token-a": {"cognito:preferred_role": testMappedRole},
		"token-b": {},
	}})
	store := svcStore(t, svc)

	created, err := store.CreateIdentityPool(cognitoidentitystore.NewIdentityPool("multi-provider", false, "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	// The sorted-first provider (alpha) carries the granting Token mapping;
	// beta's Rules mapping matches nothing and resolves Deny. Whichever
	// mapping the iteration happens to visit first decides, so only sorted
	// evaluation makes the outcome stable.
	mappings := map[string]cognitoidentitystore.RoleMapping{
		alphaProvider: {Type: "Token", AmbiguousRoleResolution: "Deny"},
		betaProvider: {
			Type: "Rules", AmbiguousRoleResolution: "Deny",
			RulesConfiguration: &cognitoidentitystore.RulesConfiguration{Rules: []cognitoidentitystore.MappingRule{
				{Claim: "audience", MatchType: "Equals", Value: "no-match", RoleARN: testOtherRole},
			}},
		},
	}
	if err := store.SetIdentityPoolRoles(created.ID, testAuthRole, testUnauthRole, mappings); err != nil {
		t.Fatalf("set pool roles: %v", err)
	}
	identity, err := store.GetOrCreateIdentityByLogins(created.ID, map[string]string{
		alphaProvider: "token-a",
		betaProvider:  "token-b",
	})
	if err != nil {
		t.Fatalf("create login identity: %v", err)
	}

	for i := 0; i < 25; i++ {
		if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
			IdentityID: identity.ID,
			Logins: map[string]string{
				betaProvider:  "token-b",
				alphaProvider: "token-a",
			},
		}); err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if issuer.roleArn != testMappedRole {
			t.Fatalf("iteration %d resolved %s, want the sorted-first provider's role %s", i, issuer.roleArn, testMappedRole)
		}
	}
}

// A presented platform user-pool token that fails validation — a bad
// signature in production — fails the credential request closed instead of
// resolving roles from an unverified claim set.
func TestGetCredentialsForIdentityInvalidPlatformTokenRefused(t *testing.T) {
	svc, _ := newMergeTestService(t)
	svc.SetCredentialIssuer(&capturingIssuer{})
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{err: errors.New("bad signature")})
	identityID := seedRolePool(t, svcStore(t, svc), rulesMapping("AuthenticatedRole", cognitoidentitystore.MappingRule{
		Claim: "department", MatchType: "Equals", Value: "engineering", RoleARN: testMappedRole,
	}))

	_, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "forged-token"},
	})
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("forged platform token returned %v, want NotAuthorizedException", err)
	}
}

// The issuer-form parser accepts the documented host with and without the
// scheme and rejects foreign provider names.
func TestPlatformUserPoolProvider(t *testing.T) {
	for name, tc := range map[string]struct {
		provider string
		region   string
		poolID   string
		ok       bool
	}{
		"bare":        {provider: "cognito-idp.us-east-1.amazonaws.com/us-east-1_POOL", region: "us-east-1", poolID: "us-east-1_POOL", ok: true},
		"schemed":     {provider: "https://cognito-idp.us-west-2.amazonaws.com/us-west-2_X", region: "us-west-2", poolID: "us-west-2_X", ok: true},
		"facebook":    {provider: "graph.facebook.com"},
		"google":      {provider: "accounts.google.com"},
		"amazon":      {provider: "www.amazon.com"},
		"oidc-issuer": {provider: "https://accounts.example.com"},
		"no-pool":     {provider: "cognito-idp.us-east-1.amazonaws.com/"},
		"bare-host":   {provider: "cognito-idp.us-east-1.amazonaws.com"},
	} {
		region, poolID, ok := platformUserPoolProvider(tc.provider)
		if ok != tc.ok || (ok && (region != tc.region || poolID != tc.poolID)) {
			t.Fatalf("%s: parsed (%q, %q, %v), want (%q, %q, %v)", name, region, poolID, ok, tc.region, tc.poolID, tc.ok)
		}
	}
}

// svcStore returns the us-east-1 store the merge-test service was seeded
// with, as the interface the Cores consume.
func svcStore(t *testing.T, svc *CognitoIdentityService) cognitoidentitystore.CognitoIdentityStoreInterface {
	t.Helper()
	v, ok := svc.stores.Load("us-east-1")
	if !ok {
		t.Fatal("service carries no us-east-1 store")
	}
	return v.(cognitoidentitystore.CognitoIdentityStoreInterface)
}

// poolOf resolves the pool an identity belongs to.
func poolOf(t *testing.T, store cognitoidentitystore.CognitoIdentityStoreInterface, identityID string) string {
	t.Helper()
	identity, err := store.GetIdentityByID(identityID)
	if err != nil {
		t.Fatalf("resolve identity %s: %v", identityID, err)
	}
	return identity.IdentityPoolID
}

// The principal-tag attribute maps attach session tags to the issued
// credentials: custom mappings take their values from the mapped claims,
// UseDefaults applies the aud and sub defaults, a map for a non-issuer
// provider contributes nothing, and no map issues an untagged session.
func TestGetCredentialsForIdentityPrincipalTagsBecomeSessionTags(t *testing.T) {
	svc, _ := newMergeTestService(t)
	issuer := &capturingIssuer{}
	svc.SetCredentialIssuer(issuer)
	svc.SetCognitoIDTokenClaimResolver(&stubIDTokenClaims{byToken: map[string]map[string]string{
		"user-token": {"department": "engineering", "aud": "7l2kkspaff2m0jkkn1o7ouvi25", "sub": "57e7b692-4f66-480d-98b8-45a6729b4c88"},
	}})
	store := svcStore(t, svc)

	identityID := seedRolePool(t, store, rulesMapping("AuthenticatedRole", cognitoidentitystore.MappingRule{
		Claim: "department", MatchType: "Equals", Value: "engineering", RoleARN: testMappedRole,
	}))
	poolID := poolOf(t, store, identityID)

	// No map stored: the session carries no tags.
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatal(err)
	}
	if issuer.tags != nil {
		t.Fatalf("session tagged without any attribute map: %v", issuer.tags)
	}

	// A custom mapping takes its tag value from the named claim; the
	// defaults are off, so aud/sub stay absent.
	if _, err := svc.setPrincipalTagAttributeMapCore(&request.RequestContext{Region: "us-east-1"}, SetPrincipalTagAttributeMapInput{
		IdentityPoolID:       poolID,
		IdentityProviderName: testUserPoolProvider,
		PrincipalTags:        map[string]string{"department": "department", "unresolved": "absent-claim"},
	}); err != nil {
		t.Fatalf("set attribute map: %v", err)
	}
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"department": "engineering"}
	if len(issuer.tags) != len(want) || issuer.tags["department"] != want["department"] {
		t.Fatalf("tags = %v, want %v", issuer.tags, want)
	}

	// UseDefaults adds the aud and sub defaults alongside the custom map.
	if _, err := svc.setPrincipalTagAttributeMapCore(&request.RequestContext{Region: "us-east-1"}, SetPrincipalTagAttributeMapInput{
		IdentityPoolID:       poolID,
		IdentityProviderName: testUserPoolProvider,
		PrincipalTags:        map[string]string{},
		UseDefaults:          true,
	}); err != nil {
		t.Fatalf("set attribute map with defaults: %v", err)
	}
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatal(err)
	}
	if issuer.tags["aud"] != "7l2kkspaff2m0jkkn1o7ouvi25" || issuer.tags["sub"] != "57e7b692-4f66-480d-98b8-45a6729b4c88" {
		t.Fatalf("default tags = %v, want the aud and sub claim values", issuer.tags)
	}

	// A map for a non-issuer provider has no verifiable claims to read and
	// contributes nothing.
	if _, err := svc.setPrincipalTagAttributeMapCore(&request.RequestContext{Region: "us-east-1"}, SetPrincipalTagAttributeMapInput{
		IdentityPoolID:       poolID,
		IdentityProviderName: "graph.facebook.com",
		PrincipalTags:        map[string]string{"department": "department"},
	}); err != nil {
		t.Fatalf("set external-provider map: %v", err)
	}
	if _, err := svc.getCredentialsForIdentityCore(&request.RequestContext{Region: "us-east-1"}, GetCredentialsForIdentityInput{
		IdentityID: identityID,
		Logins:     map[string]string{testUserPoolProvider: "user-token"},
	}); err != nil {
		t.Fatal(err)
	}
	if len(issuer.tags) != 2 {
		t.Fatalf("external-provider map contributed tags: %v", issuer.tags)
	}
}
