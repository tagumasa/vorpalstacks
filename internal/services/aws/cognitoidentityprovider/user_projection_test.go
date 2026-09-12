package cognitoidentityprovider

import (
	"context"
	"errors"
	"strings"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// attrMap converts a wire attribute list back into a name/value map.
func attrMap(list []map[string]string) map[string]string {
	m := make(map[string]string, len(list))
	for _, entry := range list {
		m[entry["Name"]] = entry["Value"]
	}
	return m
}

// The single attribute projection materialises the required sub attribute
// from the user's identifier for records whose stored attribute map carries
// no copy, and preserves an already-materialised copy.
func TestUserAttributeListMaterialisesSub(t *testing.T) {
	legacy := cognitostore.NewUser("us-east-1_POOL", "legacy-user")
	legacy.Attributes = map[string]string{"email": "legacy@example.com"}
	if got := attrMap(userAttributeList(legacy))["sub"]; got != legacy.ID {
		t.Fatalf("sub = %q for a record without a stored sub, want the user ID %q", got, legacy.ID)
	}

	current := cognitostore.NewUser("us-east-1_POOL", "current-user")
	current.Attributes = map[string]string{"sub": current.ID, "email": "current@example.com"}
	if got := attrMap(userAttributeList(current))["sub"]; got != current.ID {
		t.Fatalf("materialised sub not preserved: %q", got)
	}
}

// The GetUser and AdminGetUser projections carry exactly the members the
// model defines — the UserType status members stay off GetUser, both use
// UserAttributes (never UserType's Attributes), and the MFA preference
// members follow the activated-MFA names the model documents.
func TestUserResponseShapesMatchModel(t *testing.T) {
	user := cognitostore.NewUser("us-east-1_POOL", "mfa-user")
	user.UserStatus = "CONFIRMED"
	user.Attributes = map[string]string{"sub": user.ID}
	user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{Enabled: true, Verified: true, PreferredMfa: true}
	user.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true}

	getUser := formatGetUserResponse(user)
	for _, key := range []string{"Username", "UserAttributes", "PreferredMfaSetting", "UserMFASettingList"} {
		if _, ok := getUser[key]; !ok {
			t.Fatalf("GetUser response missing model member %q: %v", key, getUser)
		}
	}
	for _, key := range []string{"Attributes", "Enabled", "UserStatus", "UserCreateDate", "UserLastModifiedDate"} {
		if _, ok := getUser[key]; ok {
			t.Fatalf("GetUser response carries non-model member %q", key)
		}
	}
	if getUser["PreferredMfaSetting"] != "SOFTWARE_TOKEN_MFA" {
		t.Fatalf("PreferredMfaSetting = %v, want SOFTWARE_TOKEN_MFA", getUser["PreferredMfaSetting"])
	}
	settings, ok := getUser["UserMFASettingList"].([]string)
	if !ok || len(settings) != 2 || settings[0] != "SOFTWARE_TOKEN_MFA" || settings[1] != "SMS_MFA" {
		t.Fatalf("UserMFASettingList = %v, want [SOFTWARE_TOKEN_MFA SMS_MFA]", getUser["UserMFASettingList"])
	}
	if got := attrMap(getUser["UserAttributes"].([]map[string]string))["sub"]; got != user.ID {
		t.Fatalf("GetUser UserAttributes sub = %q, want the user ID", got)
	}

	adminGetUser := formatAdminGetUserResponse(user)
	want := map[string]bool{
		"Username": true, "Enabled": true, "UserStatus": true,
		"UserCreateDate": true, "UserLastModifiedDate": true,
		"UserAttributes": true, "PreferredMfaSetting": true, "UserMFASettingList": true,
	}
	if len(adminGetUser) != len(want) {
		t.Fatalf("AdminGetUser response members %v, want exactly %v", adminGetUser, want)
	}
	for key := range want {
		if _, ok := adminGetUser[key]; !ok {
			t.Fatalf("AdminGetUser response missing model member %q: %v", key, adminGetUser)
		}
	}

	// A user with no activated MFA gets neither MFA member.
	plain := cognitostore.NewUser("us-east-1_POOL", "plain-user")
	plain.Attributes = map[string]string{}
	plainResp := formatGetUserResponse(plain)
	if _, ok := plainResp["PreferredMfaSetting"]; ok {
		t.Fatal("PreferredMfaSetting set for a user without activated MFA")
	}
	if _, ok := plainResp["UserMFASettingList"]; ok {
		t.Fatal("UserMFASettingList set for a user without activated MFA")
	}

	// Activated MFA without a preferred factor: the model defines
	// UserMFASettingList as the activated factors, independent of any
	// preference, so the list must not hide behind a PreferredMfaSetting.
	noPreference := cognitostore.NewUser("us-east-1_POOL", "sms-only-user")
	noPreference.UserStatus = "CONFIRMED"
	noPreference.Attributes = map[string]string{"sub": noPreference.ID}
	noPreference.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true}
	projections := map[string]map[string]interface{}{
		"GetUser":         formatGetUserResponse(noPreference),
		"AdminGetUser":    formatAdminGetUserResponse(noPreference),
		"UserAuthFactors": computeUserAuthFactors(noPreference),
	}
	for name, resp := range projections {
		settings, ok := resp["UserMFASettingList"].([]string)
		if !ok || len(settings) != 1 || settings[0] != "SMS_MFA" {
			t.Fatalf("%s UserMFASettingList = %v, want [SMS_MFA]", name, resp["UserMFASettingList"])
		}
		if _, ok := resp["PreferredMfaSetting"]; ok {
			t.Fatalf("%s PreferredMfaSetting set without a preferred factor", name)
		}
	}
}

// The sub filter matches both a record with a materialised sub and one
// whose stored attribute map carries no copy.
func TestMatchUserFilterBySub(t *testing.T) {
	legacy := cognitostore.NewUser("us-east-1_POOL", "legacy-user")
	legacy.Attributes = map[string]string{"email": "legacy@example.com"}
	if !matchUserFilter(legacy, `sub = "`+legacy.ID+`"`) {
		t.Fatal("sub filter did not match a record without a stored sub")
	}

	current := cognitostore.NewUser("us-east-1_POOL", "current-user")
	current.Attributes = map[string]string{"sub": current.ID}
	if !matchUserFilter(current, `sub = "`+current.ID+`"`) {
		t.Fatal("sub filter did not match a materialised sub")
	}
	if matchUserFilter(current, `sub = "some-other-sub"`) {
		t.Fatal("sub filter matched a different subject")
	}
}

// The shared Core validation enforces the pool schema on every attribute
// write: unknown names are rejected, immutable attributes cannot be
// updated, values must satisfy their data type and constraints, and
// creation requires the pool's required attributes.
func TestValidateUserAttributesAgainstSchemaRules(t *testing.T) {
	pool := &cognitostore.UserPool{
		AliasAttributes:    []string{},
		UsernameAttributes: []string{},
		SchemaAttributes: []cognitostore.SchemaAttributeType{
			{Name: "rank", AttributeDataType: "String", Mutable: true,
				StringAttributeConstraints: &cognitostore.StringAttributeConstraints{MinLength: "2", MaxLength: "10"}},
			{Name: "tier", AttributeDataType: "Number", Mutable: false, Required: true,
				NumberAttributeConstraints: &cognitostore.NumberAttributeConstraints{MinValue: "1", MaxValue: "5"}},
		},
	}

	cases := []struct {
		name     string
		attrs    map[string]string
		isCreate bool
		wantErr  string
	}{
		{"unknown custom attribute", map[string]string{"custom:undefined": "x"}, true, "does not exist in the pool schema"},
		{"immutable custom attribute on update", map[string]string{"custom:tier": "3"}, false, "is immutable"},
		{"sub on update", map[string]string{"sub": "abc"}, false, "is immutable"},
		// Creation paths strip a client-supplied sub before calling the
		// validator; in isolation a non-empty sub is just another String
		// attribute as far as the constraints go.
		{"non-empty sub passes constraints in isolation", map[string]string{"sub": "abc", "custom:tier": "3"}, true, ""},
		{"boolean value required", map[string]string{"email_verified": "yes"}, true, "must be true or false"},
		{"number value required", map[string]string{"custom:tier": "gold"}, true, "must be a number"},
		{"number range", map[string]string{"custom:tier": "9"}, true, "must be at most 5"},
		{"string length minimum", map[string]string{"custom:rank": "g"}, true, "at least 2 characters"},
		{"birthdate length", map[string]string{"birthdate": "1990-01-011"}, true, "at most 10 characters"},
		{"required attribute missing on create", map[string]string{"email": "a@b.c"}, true, "required attribute custom:tier"},
		{"required attribute present", map[string]string{"custom:tier": "3"}, true, ""},
		{"update skips required enforcement", map[string]string{"email": "new@example.com"}, false, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateUserAttributesAgainstSchema(pool, tc.attrs, tc.isCreate)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("want error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

// The user-creation and attribute-update cores route every attribute write
// through the schema validation, and creation materialises the required
// sub attribute into the stored record.
func TestUserWriteCoresValidateAgainstPoolSchema(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(func() *cognitostore.UserPool {
		p := cognitostore.NewUserPool("schema-pool", "us-east-1")
		p.SchemaAttributes = []cognitostore.SchemaAttributeType{
			{Name: "rank", AttributeDataType: "String", Mutable: true},
		}
		return p
	}())
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID, ClientID: "schema-client", ClientName: "schema",
	}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	ctx := context.Background()

	if _, err := svc.adminCreateUserCore(ctx, reqCtx, AdminCreateUserInput{
		UserPoolID: pool.ID, Username: "schema-admin", MessageAction: "SUPPRESS",
		UserAttributes: map[string]string{"custom:undefined": "x"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("AdminCreateUser with an undefined attribute returned %v, want InvalidParameter", err)
	}

	if _, err := svc.adminCreateUserCore(ctx, reqCtx, AdminCreateUserInput{
		UserPoolID: pool.ID, Username: "schema-admin", MessageAction: "SUPPRESS",
		UserAttributes: map[string]string{"custom:rank": "gold", "email": "schema@example.com"},
	}); err != nil {
		t.Fatalf("AdminCreateUser with schema-conforming attributes failed: %v", err)
	}
	created, err := store.GetUser(pool.ID, "schema-admin")
	if err != nil {
		t.Fatal(err)
	}
	if created.Attributes["sub"] != created.ID {
		t.Fatalf("AdminCreateUser stored sub = %q, want the user ID %q", created.Attributes["sub"], created.ID)
	}

	if _, err := svc.adminUpdateUserAttributesCore(ctx, reqCtx, AdminUpdateUserAttributesInput{
		UserPoolID: pool.ID, Username: "schema-admin",
		UserAttributes: map[string]string{"sub": "forged"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("AdminUpdateUserAttributes writing sub returned %v, want InvalidParameter", err)
	}
	if _, err := svc.adminUpdateUserAttributesCore(ctx, reqCtx, AdminUpdateUserAttributesInput{
		UserPoolID: pool.ID, Username: "schema-admin",
		UserAttributes: map[string]string{"custom:undefined": "x"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("AdminUpdateUserAttributes with an undefined attribute returned %v, want InvalidParameter", err)
	}

	if _, err := svc.signUpCore(ctx, reqCtx, SignUpInput{
		ClientID: "schema-client", Username: "schema-signup", Password: "SchemaPass123!",
		UserAttributes: map[string]string{"custom:undefined": "x"},
	}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("SignUp with an undefined attribute returned %v, want InvalidParameter", err)
	}

	if _, err := svc.signUpCore(ctx, reqCtx, SignUpInput{
		ClientID: "schema-client", Username: "schema-signup", Password: "SchemaPass123!",
		UserAttributes: map[string]string{"custom:rank": "gold", "email": "signup@example.com"},
	}); err != nil {
		t.Fatalf("SignUp with schema-conforming attributes failed: %v", err)
	}
	signedUp, err := store.GetUser(pool.ID, "schema-signup")
	if err != nil {
		t.Fatal(err)
	}
	if signedUp.Attributes["sub"] != signedUp.ID {
		t.Fatalf("SignUp stored sub = %q, want the user ID %q", signedUp.Attributes["sub"], signedUp.ID)
	}
}

// The access-token attribute update path applies the same schema
// validation as the admin path.
func TestUpdateUserAttributesByAccessTokenValidates(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("token-schema-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	// Token issuance reads the client record and fails closed when it is
	// missing, so the fixture registers the app client it issues against.
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "token-schema-client",
		ClientName: "token-schema-client",
	}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	user := cognitostore.NewUser(pool.ID, "token-schema-user")
	user.UserStatus = "CONFIRMED"
	user.Attributes = map[string]string{"sub": user.ID, "email": "token@example.com"}
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}
	accessToken, _, _, _, err := svc.CreateTokens(reqCtx, pool.ID, user.ID, "token-schema-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatalf("create tokens: %v", err)
	}

	err = svc.updateUserAttributesByAccessTokenCore(context.Background(), reqCtx, accessToken, map[string]string{"sub": "forged"})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("UpdateUserAttributes writing sub returned %v, want InvalidParameter", err)
	}
	err = svc.updateUserAttributesByAccessTokenCore(context.Background(), reqCtx, accessToken, map[string]string{"birthdate": "1990"})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("UpdateUserAttributes with a malformed birthdate returned %v, want InvalidParameter", err)
	}
	if err := svc.updateUserAttributesByAccessTokenCore(context.Background(), reqCtx, accessToken, map[string]string{"name": "Token User"}); err != nil {
		t.Fatalf("UpdateUserAttributes with conforming attributes failed: %v", err)
	}
}

// The attribute-verification round exists for the contact attributes
// alone: both code issuance and verification refuse any other attribute
// name before touching the caller's tokens.
func TestVerificationRoundRestrictedToContactAttributes(t *testing.T) {
	svc, reqCtx, _ := newContractTestService(t)

	_, err := svc.getUserAttributeVerificationCodeCore(context.Background(), reqCtx, GetUserAttributeVerificationCodeInput{
		AccessToken: "any-token", AttributeName: "name",
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("code issuance for a non-contact attribute returned %v, want InvalidParameter", err)
	}

	_, err = svc.verifyUserAttributeCore(reqCtx, VerifyUserAttributeInput{
		AccessToken: "any-token", AttributeName: "name", Code: "123456",
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("verification for a non-contact attribute returned %v, want InvalidParameter", err)
	}
}
