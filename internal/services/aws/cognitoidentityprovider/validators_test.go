package cognitoidentityprovider

import (
	"strings"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// TestValidateCustomAttributeNameUnicodeLengths pins that
// CustomAttributeNameType follows the Smithy @length(1, 20) trait counted in
// Unicode characters; the shape's pattern uses Unicode categories, so
// multibyte names are valid input and must not be rejected on byte length.
func TestValidateCustomAttributeNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateCustomAttributeName(strings.Repeat(cjk, 20)); err != nil {
		t.Errorf("20-character CJK custom attribute name rejected: %v", err)
	}
	if err := validateCustomAttributeName(strings.Repeat(cjk, 21)); err == nil {
		t.Error("21-character CJK custom attribute name accepted")
	}
}

// TestValidateUserPoolConfigAcceptsStandardSchemaAttributes pins that the
// user pool schema name domain includes the standard attribute names: AWS
// documents modifying standard attribute properties through the
// CreateUserPool Schema parameter, and phone_number_verified — the one
// standard name longer than the custom attribute limit — is valid input.
func TestValidateUserPoolConfigAcceptsStandardSchemaAttributes(t *testing.T) {
	pool := &cognitostore.UserPool{
		Name: "standard-schema",
		SchemaAttributes: []cognitostore.SchemaAttributeType{
			{Name: "email", AttributeDataType: "String", Required: true},
			{Name: "phone_number_verified", AttributeDataType: "Boolean", Mutable: true},
		},
	}
	if err := validateUserPoolConfig(pool); err != nil {
		t.Fatalf("standard schema attributes rejected: %v", err)
	}

	// The custom attribute name limit still applies to non-standard names.
	custom := &cognitostore.UserPool{
		Name: "standard-schema",
		SchemaAttributes: []cognitostore.SchemaAttributeType{
			{Name: strings.Repeat("c", 21), AttributeDataType: "String"},
		},
	}
	if err := validateUserPoolConfig(custom); err == nil {
		t.Fatal("21-character non-standard attribute name accepted")
	}
}

// TestValidateUsernamePatternUnicodeLengths pins that UsernameType follows
// the Smithy @length(1, 128) trait counted in Unicode characters; the
// pattern uses Unicode categories (\p{L} and friends).
func TestValidateUsernamePatternUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if !validateUsernamePattern(strings.Repeat(cjk, 128)) {
		t.Error("128-character CJK username rejected")
	}
	if validateUsernamePattern(strings.Repeat(cjk, 129)) {
		t.Error("129-character CJK username accepted")
	}
	if validateUsernamePattern("") {
		t.Error("empty username accepted")
	}
}

// TestValidateRegionNameUnicodeLengths pins that RegionNameType follows the
// Smithy @length(5, 32) trait counted in Unicode characters (no pattern).
func TestValidateRegionNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateRegionName(strings.Repeat(cjk, 32)); err != nil {
		t.Errorf("32-character CJK region name rejected: %v", err)
	}
	if err := validateRegionName(strings.Repeat(cjk, 33)); err == nil {
		t.Error("33-character CJK region name accepted")
	}
}

// TestValidateClientNamePattern pins the Smithy ClientNameType constraints:
// length 1-128 and pattern ^[\w\s+=,.@-]+$ (shared with UserPoolNameType).
func TestValidateClientNamePattern(t *testing.T) {
	if validateClientNamePattern("") {
		t.Error("empty client name accepted")
	}
	if !validateClientNamePattern("web-app_1@example") {
		t.Error("valid client name rejected")
	}
	if !validateClientNamePattern(strings.Repeat("a", 128)) {
		t.Error("128-character client name rejected")
	}
	if validateClientNamePattern(strings.Repeat("a", 129)) {
		t.Error("129-character client name accepted")
	}
	if validateClientNamePattern("bad:name") {
		t.Error("colon in client name accepted")
	}
	if validateClientNamePattern("bad/name") {
		t.Error("slash in client name accepted")
	}
}

// TestValidatePasswordPolicyRanges pins the Smithy ranges for the password
// policy members: MinimumLength {6, 99}, TemporaryPasswordValidityDays
// {0, 365}, PasswordHistorySize {0, 24}. A zero MinimumLength means the
// member was never supplied and the default applies instead.
func TestValidatePasswordPolicyRanges(t *testing.T) {
	cases := []struct {
		name   string
		policy cognitostore.PasswordPolicy
		valid  bool
	}{
		{"minimum length lower boundary", cognitostore.PasswordPolicy{MinimumLength: MinPasswordPolicyMinimumLength}, true},
		{"minimum length upper boundary", cognitostore.PasswordPolicy{MinimumLength: MaxPasswordPolicyMinimumLength}, true},
		{"minimum length below range", cognitostore.PasswordPolicy{MinimumLength: MinPasswordPolicyMinimumLength - 1}, false},
		{"minimum length above range", cognitostore.PasswordPolicy{MinimumLength: MaxPasswordPolicyMinimumLength + 1}, false},
		{"temporary days upper boundary", cognitostore.PasswordPolicy{TemporaryPasswordValidityDays: MaxTemporaryPasswordValidityDays}, true},
		{"temporary days above range", cognitostore.PasswordPolicy{TemporaryPasswordValidityDays: MaxTemporaryPasswordValidityDays + 1}, false},
		{"history size upper boundary", cognitostore.PasswordPolicy{PasswordHistorySize: MaxPasswordHistorySize}, true},
		{"history size above range", cognitostore.PasswordPolicy{PasswordHistorySize: MaxPasswordHistorySize + 1}, false},
	}
	for _, tc := range cases {
		err := validatePasswordPolicyRanges(&tc.policy)
		if tc.valid && err != nil {
			t.Errorf("%s rejected: %v", tc.name, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("%s accepted", tc.name)
		}
	}
}

// TestValidateUserPoolConfig pins the whole-pool validation shared by the
// create and update paths: enum members, the alias/username exclusivity
// rule, ArnType members, and the admin-create unused-account range.
func TestValidateUserPoolConfig(t *testing.T) {
	valid := func(mutate func(*cognitostore.UserPool)) *cognitostore.UserPool {
		pool := cognitostore.NewUserPool("valid-name", "us-east-1")
		if mutate != nil {
			mutate(pool)
		}
		return pool
	}

	if err := validateUserPoolConfig(valid(nil)); err != nil {
		t.Errorf("valid pool rejected: %v", err)
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) { p.Name = "bad:name" })); err == nil {
		t.Error("invalid pool name accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) {
		p.AliasAttributes = []string{"email"}
		p.UsernameAttributes = []string{"email"}
	})); err == nil {
		t.Error("alias and username attributes together accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) { p.MfaConfiguration = "SOMETIMES" })); err == nil {
		t.Error("invalid MfaConfiguration accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) { p.DeletionProtection = "MAYBE" })); err == nil {
		t.Error("invalid DeletionProtection accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) {
		p.LambdaConfig = &cognitostore.LambdaConfig{PreSignUp: "not-an-arn"}
	})); err == nil {
		t.Error("malformed LambdaConfig trigger ARN accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) {
		p.EmailConfiguration = &cognitostore.EmailConfiguration{EmailSendingAccount: "CARRIER_PIGEON"}
	})); err == nil {
		t.Error("invalid EmailSendingAccount accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) {
		p.AdminCreateUserConfig = &cognitostore.AdminCreateUserConfig{UnusedAccountValidityDays: MaxUnusedAccountValidityDays + 1}
	})); err == nil {
		t.Error("out-of-range UnusedAccountValidityDays accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) {
		p.SchemaAttributes = []cognitostore.SchemaAttributeType{{Name: "", AttributeDataType: "String"}}
	})); err == nil {
		t.Error("empty schema attribute name accepted")
	}
	if err := validateUserPoolConfig(valid(func(p *cognitostore.UserPool) {
		p.SchemaAttributes = []cognitostore.SchemaAttributeType{{Name: "custom", AttributeDataType: "VIBES"}}
	})); err == nil {
		t.Error("invalid schema attribute data type accepted")
	}
}

// TestValidateArnType pins the Smithy ArnType constraints: length 20-2048
// and the generic ARN pattern.
func TestValidateArnType(t *testing.T) {
	if validateArnType("arn:aws:lambda:us-east-1:123456789012:function:myfn") != true {
		t.Error("valid ARN rejected")
	}
	if validateArnType("not-an-arn") {
		t.Error("plain string accepted as ARN")
	}
	if validateArnType("arn:aws:lambda:us-east-1:123456789012:function:" + strings.Repeat("f", 2048)) {
		t.Error("over-length ARN accepted")
	}
}

// TestValidatePasswordSymbolClassification pins that the symbol requirement
// is satisfied only by the special characters AWS documents (plus
// non-leading, non-trailing spaces), while non-basic-Latin letters satisfy
// no requirement — the required classes are basic Latin letters, numbers,
// and the documented specials.
func TestValidatePasswordSymbolClassification(t *testing.T) {
	policy := &cognitostore.PasswordPolicy{
		MinimumLength:    8,
		RequireUppercase: true, RequireLowercase: true,
		RequireNumbers: true, RequireSymbols: true,
	}
	if err := validatePassword("Passw0rd!", policy); err != nil {
		t.Errorf("password with documented special rejected: %v", err)
	}
	if err := validatePassword("Passw0rd ", policy); err == nil {
		t.Error("trailing space accepted as symbol")
	}
	if err := validatePassword(" Password0a", policy); err == nil {
		t.Error("leading space accepted as symbol")
	}
	if err := validatePassword("Passw0r d", policy); err != nil {
		t.Errorf("interior space not accepted as symbol: %v", err)
	}
	if err := validatePassword("P\u00e9ssw0rd", policy); err == nil {
		t.Error("non-basic-Latin letter accepted as symbol")
	}
}

// TestGenerateTemporaryPasswordPolicyCompliance pins that a generated
// temporary password satisfies the pool policy for every combination of
// character-class requirements and for elevated minimum lengths.
func TestGenerateTemporaryPasswordPolicyCompliance(t *testing.T) {
	policies := []*cognitostore.PasswordPolicy{
		nil,
		{MinimumLength: 8},
		{MinimumLength: 12, RequireUppercase: true, RequireLowercase: true, RequireNumbers: true, RequireSymbols: true},
		{MinimumLength: 20, RequireSymbols: true},
	}
	for i, p := range policies {
		pw, err := generateTemporaryPassword(p)
		if err != nil {
			t.Fatalf("policy %d: %v", i, err)
		}
		if err := validatePassword(pw, p); err != nil {
			t.Errorf("policy %d: generated password %q rejected: %v", i, pw, err)
		}
	}
}
