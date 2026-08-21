// Package paramvalidation provides the shared mechanism for Smithy
// @length and enum membership checks on request parameters. The bounds
// and valid-value sets stay per service (they come from each operation's
// model); only the check mechanism and its empty-value conventions are
// shared. Violations are reported through per-service error factories so
// each service keeps its own AWS error shape and message wording.
package paramvalidation

import "unicode/utf8"

// StringLength checks a string parameter against a Smithy @length range.
// Lengths are counted in Unicode characters, matching how Smithy @length
// traits and the AWS limit descriptions measure strings. Out-of-range
// lengths are reported through errFactory, which receives the field
// name, the actual length and the bounds.
func StringLength(field, value string, min, max int, errFactory func(field string, length, min, max int) error) error {
	n := utf8.RuneCountInString(value)
	if n < min || n > max {
		return errFactory(field, n, min, max)
	}
	return nil
}

// EnumValue checks a single enum parameter against a valid-value set. An
// empty value is valid (the parameter is unset); any other value must be
// a member of the set. Violations are reported through errFactory, which
// receives the field name and the offending value.
func EnumValue(field, value string, validValues map[string]bool, errFactory func(field, value string) error) error {
	if value == "" {
		return nil
	}
	if !validValues[value] {
		return errFactory(field, value)
	}
	return nil
}

// EnumList checks every member of a list parameter against a valid-value
// set. Unlike EnumValue there is no unset convention: every member,
// including an empty one, must be a member of the set. Violations are
// reported through errFactory with the offending value.
func EnumList(field string, values []string, validValues map[string]bool, errFactory func(field, value string) error) error {
	for _, v := range values {
		if !validValues[v] {
			return errFactory(field, v)
		}
	}
	return nil
}
