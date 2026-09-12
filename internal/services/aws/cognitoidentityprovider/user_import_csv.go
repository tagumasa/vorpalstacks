package cognitoidentityprovider

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// This file implements the CSV dialect and row semantics of the Cognito
// user import file. The dialect follows the Amazon Cognito developer guide
// ("Importing users into user pools from a CSV file", Formatting the CSV
// file): values are unquoted, a backslash escapes a comma inside a value,
// the first row is the header, and the column order does not matter. This
// differs from the RFC 4180 dialects other services import (DynamoDB
// ImportTable, Timestream BatchLoad), so the parser is service-local like
// those, on top of the shared limits defined in the store package.

// Columns with import-control meaning; they never become user attributes.
const (
	importColumnUsername  = "cognito:username"
	importColumnMFA       = "cognito:mfa_enabled"
	importColumnHash      = "password_hash"
	columnEmailVerified   = "email_verified"
	columnPhoneVerified   = "phone_number_verified"
	columnEmail           = "email"
	columnPhoneNumber     = "phone_number"
	columnBirthdate       = "birthdate"
	importStatusReset     = "RESET_REQUIRED"
	importStatusConfirmed = "CONFIRMED"
)

// userImportRow is the validated projection of one CSV data row.
type userImportRow struct {
	Username         string
	Attributes       map[string]string
	UserStatus       string
	PasswordHash     string
	PasswordHashAlgo string
	MFAOptions       []*cognitostore.MFAOptionType
	// SmsMfaEnabled mirrors the cognito:mfa_enabled column onto the modern
	// SMS factor so the imported user is actually challenged at sign-in;
	// MFAOptions above stays for the admin projections.
	SmsMfaEnabled bool
}

// importCSVRow is one data row of the import file together with its
// 1-based line number in the file. Amazon Cognito's per-user outcome logs
// refer to line numbers of the whole file — the header row is line 1, so
// the first user is "Line Number 2" — and blank lines keep their position
// in that numbering.
type importCSVRow struct {
	Fields     []string
	LineNumber int
}

// splitImportCSVLine splits one CSV line on unescaped commas and restores
// escaped characters: `\,` yields a literal comma in the value. Any other
// backslash pair yields the escaped character itself.
func splitImportCSVLine(line string) []string {
	var fields []string
	var b strings.Builder
	escaped := false
	for _, r := range line {
		switch {
		case escaped:
			b.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case r == ',':
			fields = append(fields, b.String())
			b.Reset()
		default:
			b.WriteRune(r)
		}
	}
	if escaped {
		b.WriteRune('\\')
	}
	fields = append(fields, b.String())
	return fields
}

// parseImportCSV splits the uploaded file into its header row and data
// rows. It strips a UTF-8 BOM, accepts CRLF and LF line endings, skips
// blank lines, and enforces the documented row-count ceiling at the job
// level. Per-row validation (attribute rules) happens in applyImportRow so
// one bad row cannot fail the whole job; the header itself is validated
// against the pool schema by validateImportCSVHeader.
func parseImportCSV(data []byte) (header []string, rows []importCSVRow, err error) {
	text := strings.TrimPrefix(string(data), "\ufeff")
	lines := strings.Split(text, "\n")
	for i, raw := range lines {
		lineNumber := i + 1
		line := strings.TrimSuffix(raw, "\r")
		// A whitespace-only line is blank whatever its length: it carries
		// no data row, so the skip runs before the row-length ceiling and
		// a long blank separator cannot fail the whole job. Blank lines
		// keep their position in the per-user line numbering.
		if i > 0 && strings.TrimSpace(line) == "" {
			continue
		}
		// The documented row-length ceiling counts characters of the
		// whole line, not fields: the file is unparsable line by line
		// beyond it.
		if utf8.RuneCountInString(line) > cognitostore.MaxImportCSVRowLengthChars {
			return nil, nil, fmt.Errorf("CSV line %d has %d characters, exceeding the maximum of %d", lineNumber, utf8.RuneCountInString(line), cognitostore.MaxImportCSVRowLengthChars)
		}
		if i == 0 {
			for _, name := range splitImportCSVLine(line) {
				header = append(header, strings.TrimSpace(name))
			}
			continue
		}
		rows = append(rows, importCSVRow{Fields: splitImportCSVLine(line), LineNumber: lineNumber})
	}
	if len(header) == 0 {
		return nil, nil, fmt.Errorf("CSV file has no header row")
	}
	if !hasColumn(header, importColumnUsername) {
		return nil, nil, fmt.Errorf("CSV header is missing the required %s column", importColumnUsername)
	}
	if len(rows) > cognitostore.MaxImportCSVRows {
		return nil, nil, fmt.Errorf("CSV file has %d rows, exceeding the maximum of %d", len(rows), cognitostore.MaxImportCSVRows)
	}
	return header, rows, nil
}

// validateImportCSVHeader rejects header columns that cannot carry user
// data: empty names (a stray trailing comma), duplicated names (the row
// projection keys columns by name, so a duplicate would silently collapse
// two values into one), and names that are neither an import-control
// column, a standard importable attribute, nor an attribute declared in the
// pool schema. Such a file can never produce valid users, so the whole job
// fails rather than silently dropping the column.
func validateImportCSVHeader(pool *cognitostore.UserPool, header []string) error {
	allowed := make(map[string]bool, len(csvHeaderBase)+len(pool.SchemaAttributes)+1)
	for _, name := range csvHeaderBase {
		allowed[name] = true
	}
	allowed[importColumnHash] = true
	for _, sa := range pool.SchemaAttributes {
		if sa.Name == "" || standardSchemaAttributeNames[sa.Name] {
			continue
		}
		allowed[schemaAttributeWireName(sa)] = true
	}
	seen := make(map[string]bool, len(header))
	for _, name := range header {
		if name == "" {
			return fmt.Errorf("CSV header contains an empty column name")
		}
		if seen[name] {
			return fmt.Errorf("CSV header contains the column %q more than once", name)
		}
		seen[name] = true
		if !allowed[name] {
			return fmt.Errorf("CSV header column %q is not a recognised user attribute for this user pool", name)
		}
	}
	return nil
}

func hasColumn(header []string, name string) bool {
	for _, h := range header {
		if h == name {
			return true
		}
	}
	return false
}

// applyImportRow validates one data row against the pool configuration and
// projects it into the user write. Rules from the developer guide: the
// username is required and must not contain spaces or tabs; at least one of
// the pool's auto-verified attributes must be true (when the pool has
// exactly one, that one must be true); email/phone_number must be present
// when their verified flag is true; every attribute the pool marks required
// must have a value; birthdate is mm/dd/yyyy in the file and stored as
// YYYY-MM-DD; leading/trailing whitespace is trimmed; a password_hash value
// requires the job's hashing algorithm and yields a CONFIRMED user while a
// missing hash yields RESET_REQUIRED.
func applyImportRow(pool *cognitostore.UserPool, hashAlgo string, header []string, row []string) (userImportRow, error) {
	out := userImportRow{Attributes: map[string]string{}}
	if len(row) != len(header) {
		return out, fmt.Errorf("the row has %d fields but the header has %d columns", len(row), len(header))
	}

	cols := map[string]string{}
	for i, name := range header {
		cols[name] = strings.TrimSpace(row[i])
	}

	out.Username = cols[importColumnUsername]
	if out.Username == "" {
		return out, fmt.Errorf("the %s column is required", importColumnUsername)
	}
	if strings.ContainsAny(out.Username, " \t") {
		return out, fmt.Errorf("the %s value must not contain spaces or tabs", importColumnUsername)
	}

	autoVerified := pool.AutoVerifiedAttributes
	if len(autoVerified) == 0 {
		autoVerified = []string{columnEmail, columnPhoneNumber}
	}
	oneTrue := false
	for _, attr := range autoVerified {
		verifiedCol := attr + "_verified"
		flag := strings.ToLower(cols[verifiedCol])
		if flag != "" && flag != "true" && flag != "false" {
			return out, fmt.Errorf("%s must be true or false", verifiedCol)
		}
		if flag == "true" {
			oneTrue = true
		}
	}
	if !oneTrue {
		return out, fmt.Errorf("the user record does not set any of the auto-verified attributes to true")
	}
	for _, pair := range [][2]string{{columnEmailVerified, columnEmail}, {columnPhoneVerified, columnPhoneNumber}} {
		if strings.ToLower(cols[pair[0]]) == "true" && cols[pair[1]] == "" {
			return out, fmt.Errorf("%s must have a value when %s is true", pair[1], pair[0])
		}
	}

	for name, value := range cols {
		switch name {
		case importColumnUsername, importColumnMFA:
			continue
		case importColumnHash:
			continue
		}
		if value == "" {
			continue
		}
		if name == columnBirthdate {
			converted, convErr := convertImportBirthdate(value)
			if convErr != nil {
				return out, convErr
			}
			out.Attributes[name] = converted
			continue
		}
		if name == columnEmailVerified || name == columnPhoneVerified {
			out.Attributes[name] = strings.ToLower(value)
			continue
		}
		out.Attributes[name] = value
	}

	// The shared Core validation covers the schema side of the row: value
	// constraints, and the required attributes every import record must
	// populate. A row that fails does so for that user only — the import
	// job itself continues.
	if err := validateUserAttributesAgainstSchema(pool, out.Attributes, true); err != nil {
		return out, err
	}

	if err := validateImportMFASetting(pool, cols[importColumnMFA]); err != nil {
		return out, err
	}

	if strings.ToLower(cols[importColumnMFA]) == "true" && out.Attributes[columnPhoneNumber] != "" {
		// cognito:mfa_enabled turns on SMS MFA for the imported user: the
		// entry is recorded in both models — the legacy MFAOptions the
		// admin projections read, and the modern SmsMfa flag the sign-in
		// second-factor decision challenges with. A row with the flag but
		// no phone number imports all the same: the developer guide
		// records that imported users may hold the MFA-enabled state
		// without a valid factor — such users cannot complete sign-in
		// until they configure an email attribute, phone number, or TOTP
		// that is a valid factor in their pool.
		out.MFAOptions = []*cognitostore.MFAOptionType{{
			DeliveryMedium: "SMS",
			AttributeName:  columnPhoneNumber,
		}}
		out.SmsMfaEnabled = true
	}

	out.UserStatus = importStatusReset
	if hash := cols[importColumnHash]; hash != "" {
		if hashAlgo == "" {
			return out, fmt.Errorf("a password_hash value was supplied but the import job specifies no hashing algorithm")
		}
		if err := validateImportedHashParams(hashAlgo, hash); err != nil {
			return out, err
		}
		out.PasswordHash = hash
		out.PasswordHashAlgo = hashAlgo
		out.UserStatus = importStatusConfirmed
	}
	return out, nil
}

// convertImportBirthdate converts the file's mm/dd/yyyy representation to
// the stored YYYY-MM-DD attribute format.
func convertImportBirthdate(value string) (string, error) {
	t, err := time.Parse("01/02/2006", value)
	if err != nil {
		return "", fmt.Errorf("birthdate %q must be in mm/dd/yyyy format", value)
	}
	return t.Format("2006-01-02"), nil
}

// validateImportMFASetting enforces the developer-guide rule that the
// cognito:mfa_enabled column must correspond to the pool's MFA
// configuration: a required-MFA pool accepts only true or blank, an
// MFA-off pool accepts only false or blank. A blank value simply adopts
// the pool-required state.
func validateImportMFASetting(pool *cognitostore.UserPool, value string) error {
	if value == "" {
		return nil
	}
	lower := strings.ToLower(value)
	if lower != "true" && lower != "false" {
		return fmt.Errorf("%s must be true or false", importColumnMFA)
	}
	switch pool.MfaConfiguration {
	case "ON":
		if lower != "true" {
			return fmt.Errorf("%s must be true or blank because the user pool requires MFA", importColumnMFA)
		}
	case "OFF":
		if lower != "false" {
			return fmt.Errorf("%s must be false or blank because the user pool has MFA disabled", importColumnMFA)
		}
	}
	return nil
}
