package cognitoidentityprovider

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func TestSplitImportCSVLine(t *testing.T) {
	cases := []struct {
		name string
		line string
		want []string
	}{
		{"plain", "a,b,c", []string{"a", "b", "c"}},
		{"escaped comma", `123 Any Street\, Suite 5,johndoe`, []string{"123 Any Street, Suite 5", "johndoe"}},
		{"escaped backslash", `a,b\\c,d`, []string{"a", `b\c`, "d"}},
		{"trailing lone backslash", `a,b\`, []string{"a", `b\`}},
		{"single field", "only", []string{"only"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitImportCSVLine(tc.line)
			if len(got) != len(tc.want) {
				t.Fatalf("got %d fields %v, want %d %v", len(got), got, len(tc.want), tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("field %d = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestParseImportCSV(t *testing.T) {
	t.Run("bom crlf and blank lines", func(t *testing.T) {
		data := "\ufeffcognito:username,email,email_verified\r\njohndoe,j@example.com,TRUE\r\n\r\n"
		header, rows, err := parseImportCSV([]byte(data))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(header) != 3 || header[0] != "cognito:username" {
			t.Fatalf("header = %v", header)
		}
		if len(rows) != 1 || rows[0].Fields[0] != "johndoe" {
			t.Fatalf("rows = %v", rows)
		}
		// The header is file line 1, so the first user row is line 2.
		if rows[0].LineNumber != 2 {
			t.Fatalf("line number = %d, want 2", rows[0].LineNumber)
		}
	})
	t.Run("missing username column fails the job", func(t *testing.T) {
		if _, _, err := parseImportCSV([]byte("email,email_verified\na@b.c,true\n")); err == nil || !strings.Contains(err.Error(), "cognito:username") {
			t.Fatalf("want missing-username error, got %v", err)
		}
	})
	t.Run("empty file", func(t *testing.T) {
		if _, _, err := parseImportCSV([]byte("")); err == nil {
			t.Fatal("want error for empty file")
		}
	})
}

func importTestPool() *cognitostore.UserPool {
	return &cognitostore.UserPool{
		AutoVerifiedAttributes: []string{"email"},
		SchemaAttributes: []cognitostore.SchemaAttributeType{
			{Name: "email", Required: true},
			{Name: "rank"},
		},
	}
}

func TestApplyImportRowHappyPath(t *testing.T) {
	header := []string{"cognito:username", "email", "email_verified", "birthdate", "updated_at", "custom:rank", "cognito:mfa_enabled", "phone_number", "password_hash"}
	row := []string{"johndoe", "johndoe@example.com", "TRUE", "02/01/1985", "1471453471", "gold", "", "+14325551212", ""}

	got, err := applyImportRow(importTestPool(), "", header, row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != "johndoe" {
		t.Errorf("username = %q", got.Username)
	}
	if got.UserStatus != importStatusReset {
		t.Errorf("status without password hash = %q, want RESET_REQUIRED", got.UserStatus)
	}
	if got.Attributes["birthdate"] != "1985-02-01" {
		t.Errorf("birthdate = %q, want converted 1985-02-01", got.Attributes["birthdate"])
	}
	if got.Attributes["updated_at"] != "1471453471" {
		t.Errorf("updated_at = %q, want passthrough", got.Attributes["updated_at"])
	}
	if got.Attributes["custom:rank"] != "gold" {
		t.Errorf("custom:rank = %q", got.Attributes["custom:rank"])
	}
	if got.Attributes["email_verified"] != "true" {
		t.Errorf("email_verified = %q, want lowercased true", got.Attributes["email_verified"])
	}
	if len(got.MFAOptions) != 0 {
		t.Errorf("blank mfa_enabled must not set MFAOptions: %v", got.MFAOptions)
	}
}

func TestApplyImportRowMFAEnabled(t *testing.T) {
	header := []string{"cognito:username", "email", "email_verified", "phone_number", "cognito:mfa_enabled"}
	row := []string{"janedoe", "j@example.com", "true", "+14325559999", "true"}
	got, err := applyImportRow(importTestPool(), "", header, row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.MFAOptions) != 1 || got.MFAOptions[0].DeliveryMedium != "SMS" || got.MFAOptions[0].AttributeName != "phone_number" {
		t.Fatalf("MFAOptions = %v", got.MFAOptions)
	}
}

func TestApplyImportRowRejections(t *testing.T) {
	header := []string{"cognito:username", "email", "email_verified"}
	cases := []struct {
		name string
		row  []string
		want string
	}{
		{"missing username", []string{"", "a@b.c", "true"}, "cognito:username column is required"},
		{"username with space", []string{"john doe", "a@b.c", "true"}, "must not contain spaces or tabs"},
		{"no auto-verified true", []string{"johndoe", "a@b.c", "false"}, "auto-verified"},
		{"verified without value", []string{"johndoe", "", "true"}, "must have a value"},
		{"invalid verified flag", []string{"johndoe", "a@b.c", "yes"}, "true or false"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := applyImportRow(importTestPool(), "", header, tc.row)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestApplyImportRowRequiredAttributes(t *testing.T) {
	pool := importTestPool()
	pool.SchemaAttributes = append(pool.SchemaAttributes, cognitostore.SchemaAttributeType{Name: "loyalty", Required: true})
	header := []string{"cognito:username", "email", "email_verified", "custom:loyalty"}
	if _, err := applyImportRow(pool, "", header, []string{"johndoe", "a@b.c", "true", ""}); err == nil || !strings.Contains(err.Error(), "required attribute custom:loyalty") {
		t.Fatalf("want required-attribute error, got %v", err)
	}
	if _, err := applyImportRow(pool, "", header, []string{"johndoe", "a@b.c", "true", "silver"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyImportRowBirthdateFormat(t *testing.T) {
	header := []string{"cognito:username", "email", "email_verified", "birthdate"}
	if _, err := applyImportRow(importTestPool(), "", header, []string{"johndoe", "a@b.c", "true", "1985-02-01"}); err == nil || !strings.Contains(err.Error(), "mm/dd/yyyy") {
		t.Fatalf("want birthdate format error, got %v", err)
	}
}

func TestApplyImportRowFieldCountMismatch(t *testing.T) {
	header := []string{"cognito:username", "email", "email_verified"}
	cases := []struct {
		name string
		row  []string
	}{
		{"short row", []string{"johndoe", "a@b.c"}},
		{"long row", []string{"johndoe", "a@b.c", "true", "extra"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := applyImportRow(importTestPool(), "", header, tc.row)
			if err == nil || !strings.Contains(err.Error(), "header has") {
				t.Fatalf("want field-count error, got %v", err)
			}
		})
	}
}

func TestValidateImportCSVHeader(t *testing.T) {
	t.Run("empty column name", func(t *testing.T) {
		header := []string{"cognito:username", ""}
		if err := validateImportCSVHeader(importTestPool(), header); err == nil || !strings.Contains(err.Error(), "empty column name") {
			t.Fatalf("want empty-column error, got %v", err)
		}
	})
	t.Run("unknown column", func(t *testing.T) {
		header := []string{"cognito:username", "not_an_attribute"}
		if err := validateImportCSVHeader(importTestPool(), header); err == nil || !strings.Contains(err.Error(), "not_an_attribute") {
			t.Fatalf("want unknown-column error, got %v", err)
		}
	})
	t.Run("standard and custom columns accepted", func(t *testing.T) {
		header := []string{"cognito:username", "email", "email_verified", "custom:rank", "password_hash", "cognito:mfa_enabled"}
		if err := validateImportCSVHeader(importTestPool(), header); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestApplyImportRowSchemaValidation(t *testing.T) {
	pool := importTestPool()
	pool.SchemaAttributes = append(pool.SchemaAttributes,
		cognitostore.SchemaAttributeType{
			Name:              "rank",
			AttributeDataType: "Number",
			NumberAttributeConstraints: &cognitostore.NumberAttributeConstraints{
				MinValue: "1", MaxValue: "10",
			},
		},
		cognitostore.SchemaAttributeType{
			Name:              "level",
			AttributeDataType: "String",
			StringAttributeConstraints: &cognitostore.StringAttributeConstraints{
				MinLength: "2", MaxLength: "4",
			},
		},
	)
	header := []string{"cognito:username", "email", "email_verified", "custom:rank", "custom:level", "updated_at"}
	cases := []struct {
		name string
		row  []string
		want string
	}{
		{"number not numeric", []string{"u", "a@b.c", "true", "high", "ab", "1"}, "must be a number"},
		{"number below min", []string{"u", "a@b.c", "true", "0", "ab", "1"}, "at least 1"},
		{"number above max", []string{"u", "a@b.c", "true", "11", "ab", "1"}, "at most 10"},
		{"string below min length", []string{"u", "a@b.c", "true", "5", "a", "1"}, "at least 2 characters"},
		{"string above max length", []string{"u", "a@b.c", "true", "5", "abcde", "1"}, "at most 4 characters"},
		{"updated_at not numeric", []string{"u", "a@b.c", "true", "5", "ab", "yesterday"}, "must be a number"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := applyImportRow(pool, "", header, tc.row)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
	// Boundaries pass, including the standard Boolean attribute.
	if _, err := applyImportRow(pool, "", header, []string{"u", "a@b.c", "true", "10", "abcd", "1471453471"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyImportRowBooleanSchemaValidation(t *testing.T) {
	header := []string{"cognito:username", "email", "email_verified"}
	// The standard email_verified column is Boolean-typed; a non-boolean
	// value is rejected by the schema check as well as the auto-verified
	// flag check.
	if _, err := applyImportRow(importTestPool(), "", header, []string{"u", "a@b.c", "maybe"}); err == nil {
		t.Fatal("want boolean validation error, got nil")
	}
}

func TestValidateImportMFASetting(t *testing.T) {
	cases := []struct {
		name      string
		poolMFA   string
		value     string
		wantError bool
	}{
		{"required pool true", "ON", "true", false},
		{"required pool blank", "ON", "", false},
		{"required pool false", "ON", "false", true},
		{"off pool false", "OFF", "false", false},
		{"off pool blank", "OFF", "", false},
		{"off pool true", "OFF", "true", true},
		{"optional pool anything", "OPTIONAL", "true", false},
		{"invalid literal", "OFF", "yes", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pool := importTestPool()
			pool.MfaConfiguration = tc.poolMFA
			err := validateImportMFASetting(pool, tc.value)
			if tc.wantError && err == nil {
				t.Fatal("want error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateImportJobName(t *testing.T) {
	cases := []struct {
		name  string
		value string
		ok    bool
	}{
		{"plain", "my-job", true},
		{"allowed punctuation", "job 2+ =,.-@_", true},
		{"empty", "", false},
		{"over 128 chars", strings.Repeat("a", 129), false},
		{"exactly 128 chars", strings.Repeat("a", 128), true},
		{"disallowed character", "job:name", false},
		{"disallowed slash", "job/name", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := validateImportJobName(tc.value); got != tc.ok {
				t.Fatalf("validateImportJobName(%q) = %v, want %v", tc.value, got, tc.ok)
			}
		})
	}
}

func TestApplyImportRowPasswordHash(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("s3cret!"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate bcrypt: %v", err)
	}
	header := []string{"cognito:username", "email", "email_verified", "password_hash"}

	got, err := applyImportRow(importTestPool(), "BCRYPT", header, []string{"johndoe", "a@b.c", "true", string(hash)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserStatus != importStatusConfirmed {
		t.Errorf("status with password hash = %q, want CONFIRMED", got.UserStatus)
	}
	if got.PasswordHash != string(hash) || got.PasswordHashAlgo != "BCRYPT" {
		t.Errorf("hash not carried: algo=%q", got.PasswordHashAlgo)
	}

	if _, err := applyImportRow(importTestPool(), "", header, []string{"johndoe", "a@b.c", "true", string(hash)}); err == nil || !strings.Contains(err.Error(), "no hashing algorithm") {
		t.Fatalf("want missing-algorithm error, got %v", err)
	}
	if _, err := applyImportRow(importTestPool(), "BCRYPT", header, []string{"johndoe", "a@b.c", "true", "not-a-bcrypt-hash"}); err == nil || !strings.Contains(err.Error(), "malformed BCRYPT") {
		t.Fatalf("want malformed error, got %v", err)
	}
}

func TestValidateImportedHashParamsBounds(t *testing.T) {
	salt := make([]byte, 16)
	for i := range salt {
		salt[i] = byte(i)
	}
	scryptHash, err := scrypt.Key([]byte("pw"), salt, 16384, 8, 1, 32)
	if err != nil {
		t.Fatalf("scrypt: %v", err)
	}
	overScrypt, err := scrypt.Key([]byte("pw"), salt, 131072, 8, 1, 32)
	if err != nil {
		t.Fatalf("scrypt over: %v", err)
	}
	cases := []struct {
		name string
		algo string
		hash string
		ok   bool
	}{
		{"scrypt in bounds", "SCRYPT", fmtScrypt(16384, 8, 1, salt, scryptHash), true},
		{"scrypt N over", "SCRYPT", fmtScrypt(131072, 8, 1, salt, overScrypt), false},
		{"argon2 in bounds", "ARGON2ID", fmtArgon2("pw", 19456, 2, 1, salt), true},
		{"argon2 memory over", "ARGON2ID", fmtArgon2("pw", 32768, 2, 1, salt), false},
		{"argon2 time over", "ARGON2ID", fmtArgon2("pw", 19456, 3, 1, salt), false},
		{"argon2 parallelism over", "ARGON2ID", fmtArgon2("pw", 19456, 2, 2, salt), false},
		{"pbkdf2 in bounds", "PBKDF2_SHA256", fmtPbkdf2("pw", 1000, salt), true},
		{"pbkdf2 iterations over", "PBKDF2_SHA256", fmtPbkdf2("pw", 600001, salt), false},
		{"bcrypt cost over", "BCRYPT", "$2b$13$CtA.Rcu/szzn9U00wpUjOuN3vrgJRZycv4aOzcP3GzqzO8UDPEFq", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateImportedHashParams(tc.algo, tc.hash)
			if tc.ok && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.ok && err == nil {
				t.Fatal("want bounds error, got nil")
			}
		})
	}
}

func TestVerifyImportedPasswordHash(t *testing.T) {
	password := "s3cret!"
	salt := make([]byte, 16)
	for i := range salt {
		salt[i] = byte(i + 1)
	}

	bhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	shash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		t.Fatalf("scrypt: %v", err)
	}

	cases := []struct {
		algo string
		hash string
	}{
		{"BCRYPT", string(bhash)},
		{"SCRYPT", fmtScrypt(16384, 8, 1, salt, shash)},
		{"ARGON2ID", fmtArgon2(password, 19456, 2, 1, salt)},
		{"PBKDF2_SHA256", fmtPbkdf2(password, 1000, salt)},
	}
	for _, tc := range cases {
		t.Run(tc.algo+" positive", func(t *testing.T) {
			if !verifyImportedPasswordHash(tc.algo, tc.hash, password) {
				t.Fatal("want verification to succeed")
			}
		})
		t.Run(tc.algo+" negative", func(t *testing.T) {
			if verifyImportedPasswordHash(tc.algo, tc.hash, "wrong") {
				t.Fatal("want verification to fail for wrong password")
			}
		})
	}
}

func fmtScrypt(n, r, p int, salt, hash []byte) string {
	return strconv.Itoa(n) + "$" + strconv.Itoa(r) + "$" + strconv.Itoa(p) +
		"$" + hex.EncodeToString(salt) + "$" + hex.EncodeToString(hash)
}

func fmtArgon2(password string, mem, tCost, par int, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, uint32(tCost), uint32(mem), uint8(par), 32)
	return "$argon2id$v=19$m=" + strconv.Itoa(mem) + ",t=" + strconv.Itoa(tCost) + ",p=" + strconv.Itoa(par) +
		"$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
}

func fmtPbkdf2(password string, iterations int, salt []byte) string {
	hash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return "$pbkdf2-sha256$" + strconv.Itoa(iterations) +
		"$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
}

func TestParseImportCSVRowLengthLimit(t *testing.T) {
	header := "cognito:username,email,email_verified"
	atLimit := header + "\n" + "johndoe," + strings.Repeat("a", cognitostore.MaxImportCSVRowLengthChars-len("johndoe,"))
	if _, _, err := parseImportCSV([]byte(atLimit)); err != nil {
		t.Fatalf("row at the limit must parse: %v", err)
	}
	overLimit := header + "\n" + "johndoe," + strings.Repeat("a", cognitostore.MaxImportCSVRowLengthChars-len("johndoe,")+1)
	_, _, err := parseImportCSV([]byte(overLimit))
	if err == nil {
		t.Fatal("row over the character limit must be rejected")
	}
}

func TestValidateImportedHashParamsRejectsNonPositiveParameters(t *testing.T) {
	cases := []struct{ algo, hash string }{
		{"BCRYPT", "$2b$00$" + strings.Repeat("a", 53)},
		{"SCRYPT", fmtScrypt(0, 8, 1, []byte("salt"), nil)},
		{"SCRYPT", fmtScrypt(1024, 0, 1, []byte("salt"), nil)},
		{"ARGON2ID", "$argon2id$v=19$m=1024,t=0,p=1$c2FsdA$aGFzaA"},
		{"ARGON2ID", "$argon2id$v=19$m=1024,t=2,p=0$c2FsdA$aGFzaA"},
		{"ARGON2ID", "$argon2id$v=19$m=4,t=2,p=1$c2FsdA$aGFzaA"},
		{"PBKDF2_SHA256", "$pbkdf2-sha256$0$c2FsdA$aGFzaA"},
	}
	for _, tc := range cases {
		if err := validateImportedHashParams(tc.algo, tc.hash); err == nil {
			t.Errorf("%s hash %q with non-positive parameters must be rejected", tc.algo, tc.hash)
		}
	}
}

// A stored hash with zero-valued parameters must fail verification
// cleanly; the argon2 library panics on rounds or parallelism below one.
func TestVerifyImportedPasswordHashZeroParamsNoPanic(t *testing.T) {
	if verifyImportedPasswordHash("ARGON2ID", "$argon2id$v=19$m=1024,t=0,p=1$c2FsdA$aGFzaA", "password") {
		t.Fatal("verification must fail for a zero-parameter hash")
	}
	if verifyImportedPasswordHash("ARGON2ID", "$argon2id$v=19$m=1024,t=2,p=0$c2FsdA$aGFzaA", "password") {
		t.Fatal("verification must fail for a zero-parameter hash")
	}
}
