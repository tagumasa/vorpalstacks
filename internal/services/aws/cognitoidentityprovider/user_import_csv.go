package cognitoidentityprovider

import (
	"fmt"
	"strconv"
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
		if strings.TrimSpace(line) == "" {
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
// data: empty names (a stray trailing comma) and names that are neither an
// import-control column, a standard importable attribute, nor an attribute
// declared in the pool schema. Such a file can never produce valid users,
// so the whole job fails rather than silently dropping the column.
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
	for _, name := range header {
		if name == "" {
			return fmt.Errorf("CSV header contains an empty column name")
		}
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

	for _, sa := range pool.SchemaAttributes {
		if !sa.Required {
			continue
		}
		name := schemaAttributeWireName(sa)
		// sub is assigned by the service for every user; it is never
		// supplied through the file.
		if name == "sub" {
			continue
		}
		if cols[name] == "" {
			return out, fmt.Errorf("required attribute %s must have a value", name)
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

	if err := validateImportRowAgainstSchema(pool, out.Attributes); err != nil {
		return out, err
	}

	if err := validateImportMFASetting(pool, cols[importColumnMFA]); err != nil {
		return out, err
	}

	if strings.ToLower(cols[importColumnMFA]) == "true" && out.Attributes[columnPhoneNumber] != "" {
		out.MFAOptions = []*cognitostore.MFAOptionType{{
			DeliveryMedium: "SMS",
			AttributeName:  columnPhoneNumber,
		}}
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

// validateImportRowAgainstSchema checks every populated attribute value
// against the pool's attribute schema (custom definitions override the
// standard defaults). A row whose values do not match the schema fails for
// that user only — the import job itself continues.
func validateImportRowAgainstSchema(pool *cognitostore.UserPool, attrs map[string]string) error {
	schema := make(map[string]cognitostore.SchemaAttributeType)
	for _, sa := range schemaAttributesForDescribe(pool) {
		schema[sa.Name] = sa
	}
	for name, value := range attrs {
		sa, ok := schema[name]
		if !ok {
			continue
		}
		switch sa.AttributeDataType {
		case "Boolean":
			if lower := strings.ToLower(value); lower != "true" && lower != "false" {
				return fmt.Errorf("attribute %s must be true or false", name)
			}
		case "Number":
			num, numErr := strconv.ParseFloat(value, 64)
			if numErr != nil {
				return fmt.Errorf("attribute %s must be a number", name)
			}
			if c := sa.NumberAttributeConstraints; c != nil {
				if c.MinValue != "" {
					if min, minErr := strconv.ParseFloat(c.MinValue, 64); minErr == nil && num < min {
						return fmt.Errorf("attribute %s must be at least %s", name, c.MinValue)
					}
				}
				if c.MaxValue != "" {
					if max, maxErr := strconv.ParseFloat(c.MaxValue, 64); maxErr == nil && num > max {
						return fmt.Errorf("attribute %s must be at most %s", name, c.MaxValue)
					}
				}
			}
		case "String":
			length := utf8.RuneCountInString(value)
			if c := sa.StringAttributeConstraints; c != nil {
				if c.MinLength != "" {
					if min, minErr := strconv.Atoi(c.MinLength); minErr == nil && length < min {
						return fmt.Errorf("attribute %s must be at least %s characters long", name, c.MinLength)
					}
				}
				if c.MaxLength != "" {
					if max, maxErr := strconv.Atoi(c.MaxLength); maxErr == nil && length > max {
						return fmt.Errorf("attribute %s must be at most %s characters long", name, c.MaxLength)
					}
				}
			}
		}
	}
	return nil
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
