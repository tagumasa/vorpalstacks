// Package sts validators — Smithy trait enforcement for string parameters.
//
// This file centralises all input validation that was previously inlined
// as bare json.Unmarshal calls. Each validator enforces the corresponding
// Smithy shape constraint (length, pattern) before structural parsing.
package sts

import (
	"encoding/json"
	"regexp"
	"unicode/utf8"
)

// sessionPolicyPattern is the Smithy pattern trait shared by
// sessionPolicyDocumentType and unrestrictedSessionPolicyDocumentType.
// Allows TAB, LF, CR, and all characters in the range U+0020–U+00FF
// (printable Latin-1).
var sessionPolicyPattern = regexp.MustCompile(`^[\t\n\r\x20-\xff]+$`)

// Session policy Smithy limits.
const (
	maxSessionPolicyLen    = 2048 // sessionPolicyDocumentType max
	maxSAMLAssertionLen    = 100000
	minSAMLAssertionLen    = 4
	maxProviderIDLen       = 2048 // urlType max
	minProviderIDLen       = 4    // urlType min
	minARNLen              = 20   // arnType min
	maxARNLen              = 2048 // arnType max
	maxContextAssertionLen = 2048 // contextAssertionType max
	minContextAssertionLen = 4    // contextAssertionType min
	maxProvidedContexts    = 5    // ProvidedContextsListType max
	maxPolicyArns          = 10   // AWS docs: "up to 10 managed policy ARNs"
)

// arnTypePattern is the Smithy pattern trait for arnType. Allows TAB,
// LF, CR, space, printable ASCII (0x20-0x7E), and Unicode code points
// U+0085, U+00A0-U+D7FF, U+E000-U+FFFD, U+10000-U+10FFFF.
var arnTypePattern = regexp.MustCompile(`^[\t\n\r\x20-\x7e\x85\xa0-\x{d7ff}\x{e000}-\x{fffd}\x{10000}-\x{10ffff}]+$`)

// validateARN enforces the Smithy arnType constraint: length 20-2048
// (counted in Unicode characters — the arnType pattern admits code points
// beyond ASCII) and the arnType pattern. Returns ErrInvalidPolicyArn on
// failure. Callers validating ProviderArn should translate the error to
// ErrInvalidProviderArn so the field name in the error message is
// accurate.
func validateARN(arn string) error {
	if n := utf8.RuneCountInString(arn); n < minARNLen || n > maxARNLen {
		return ErrInvalidPolicyArn
	}
	if !arnTypePattern.MatchString(arn) {
		return ErrInvalidPolicyArn
	}
	return nil
}

// validateSessionPolicy enforces the Smithy sessionPolicyDocumentType
// constraint: length 1-2048 counted in Unicode characters (the Latin-1
// pattern permits two-byte characters), pattern ^[\t\n\r\x20-\xff]+$, and
// valid JSON structure. Used by AssumeRoleWithSAML,
// AssumeRoleWithWebIdentity, and GetFederationToken.
func validateSessionPolicy(policy string) error {
	if policy == "" {
		return nil // optional parameter
	}
	if utf8.RuneCountInString(policy) > maxSessionPolicyLen {
		return ErrMalformedPolicyDocument
	}
	if !sessionPolicyPattern.MatchString(policy) {
		return ErrMalformedPolicyDocument
	}
	var js interface{}
	if err := json.Unmarshal([]byte(policy), &js); err != nil {
		return ErrMalformedPolicyDocument
	}
	return nil
}

// validateUnrestrictedSessionPolicy enforces the Smithy
// unrestrictedSessionPolicyDocumentType constraint: length 1+, pattern
// ^[\t\n\r\x20-\xff]+$, and valid JSON structure. Used by AssumeRole
// whose Policy member uses the unrestricted variant (no upper length
// bound in Smithy — the packed-policy-size check serves as the
// practical limit).
func validateUnrestrictedSessionPolicy(policy string) error {
	if policy == "" {
		return nil // optional parameter
	}
	if !sessionPolicyPattern.MatchString(policy) {
		return ErrMalformedPolicyDocument
	}
	var js interface{}
	if err := json.Unmarshal([]byte(policy), &js); err != nil {
		return ErrMalformedPolicyDocument
	}
	return nil
}

// validateSAMLAssertion enforces the Smithy SAMLAssertionType constraint:
// length 4-100000 counted in Unicode characters (no pattern).
func validateSAMLAssertion(s string) error {
	if n := utf8.RuneCountInString(s); n < minSAMLAssertionLen || n > maxSAMLAssertionLen {
		return ErrInvalidSAMLAssertion
	}
	return nil
}

// validateRoleArn enforces the Smithy arnType constraint for the RoleArn
// parameter of AssumeRole, AssumeRoleWithSAML, and AssumeRoleWithWebIdentity:
// length 20-2048 (counted in Unicode characters) and the arnType pattern.
// An empty value is rejected separately by the caller (ErrInvalidRoleArn)
// before reaching this function.
func validateRoleArn(arn string) error {
	if n := utf8.RuneCountInString(arn); n < minARNLen || n > maxARNLen {
		return ErrInvalidRoleArnFormat
	}
	if !arnTypePattern.MatchString(arn) {
		return ErrInvalidRoleArnFormat
	}
	return nil
}

// validatePrincipalArn enforces the Smithy arnType constraint for the
// PrincipalArn parameter of AssumeRoleWithSAML: length 20-2048 (counted in
// Unicode characters) and the arnType pattern.
func validatePrincipalArn(arn string) error {
	if n := utf8.RuneCountInString(arn); n < minARNLen || n > maxARNLen {
		return ErrInvalidPrincipalArnFormat
	}
	if !arnTypePattern.MatchString(arn) {
		return ErrInvalidPrincipalArnFormat
	}
	return nil
}

// validateProviderID enforces the Smithy urlType constraint for the
// ProviderId parameter of AssumeRoleWithWebIdentity: length 4-2048 counted
// in Unicode characters (no pattern).
func validateProviderID(s string) error {
	if s == "" {
		return nil // optional parameter
	}
	if n := utf8.RuneCountInString(s); n < minProviderIDLen || n > maxProviderIDLen {
		return ErrInvalidProviderID
	}
	return nil
}
