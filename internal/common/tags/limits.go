package tags

import (
	"sort"
	"strings"
	"unicode/utf8"
)

// TagLimits describes the tag constraints of one service API. The zero
// value enforces no limits at all, so services built their profiles by
// setting only the fields their AWS documentation specifies.
type TagLimits struct {
	// MaxCount is the maximum number of tags on one resource; 0 disables
	// the count check.
	MaxCount int

	// MinKeyLength is the minimum tag key length; 0 disables the lower
	// bound.
	MinKeyLength int

	// MaxKeyLength is the maximum tag key length; 0 disables the upper
	// bound.
	MaxKeyLength int

	// MaxValueLength is the maximum tag value length; 0 disables the
	// upper bound.
	MaxValueLength int

	// ReservedPrefix is the key prefix reserved for AWS use (conventionally
	// "aws:"); empty disables the reservation check.
	ReservedPrefix string

	// ReservedCaseSensitive controls whether the reserved prefix comparison
	// is case sensitive. AWS treats the aws: reservation case-insensitively
	// for the services that document it, so the default comparison lower
	// cases both sides.
	ReservedCaseSensitive bool
}

// StandardLimits returns the AWS-wide tag limits: at most MaxTagsPerResource
// tags, keys of 1 to MaxTagKeyLength characters, values of at most
// MaxTagValueLength characters and the aws: key prefix reserved. Services
// whose documented limits differ must build their own profile instead of
// weakening this one.
func StandardLimits() TagLimits {
	return TagLimits{
		MaxCount:       MaxTagsPerResource,
		MinKeyLength:   1,
		MaxKeyLength:   MaxTagKeyLength,
		MaxValueLength: MaxTagValueLength,
		ReservedPrefix: "aws:",
	}
}

// Violation classifies why a tag list failed a TagLimits check.
type Violation int

const (
	// OK means the tag list satisfies the limits.
	OK Violation = iota
	// TooManyTags means the list exceeds MaxCount.
	TooManyTags
	// TagKeyTooShort means a key is shorter than MinKeyLength.
	TagKeyTooShort
	// TagKeyTooLong means a key is longer than MaxKeyLength.
	TagKeyTooLong
	// TagValueTooLong means a value is longer than MaxValueLength.
	TagValueTooLong
	// ReservedTagKey means a key starts with ReservedPrefix.
	ReservedTagKey
)

// CheckTags validates a typed tag list against the limits and reports the
// first violation in canonical order: count, then per-tag reserved prefix,
// key length and value length. Key and value lengths are counted in
// Unicode characters, matching the AWS tag limit definitions. Reporting
// the key-side violations before the value-side ones keeps the error
// code deterministic when a tag violates several rules at once.
//
// Alongside the violation it returns the offending tag key (list order for
// the slice form), or the empty string when there is no offender to name.
// Consumers must format error messages from this finding instead of
// re-deriving the offender with a second, drifting scan.
func CheckTags(tagList []Tag, limits TagLimits) (Violation, string) {
	if limits.MaxCount > 0 && len(tagList) > limits.MaxCount {
		return TooManyTags, ""
	}
	for _, t := range tagList {
		if violatesReservedPrefix(t.Key, limits) {
			return ReservedTagKey, t.Key
		}
		if limits.MinKeyLength > 0 && utf8.RuneCountInString(t.Key) < limits.MinKeyLength {
			return TagKeyTooShort, t.Key
		}
		if limits.MaxKeyLength > 0 && utf8.RuneCountInString(t.Key) > limits.MaxKeyLength {
			return TagKeyTooLong, t.Key
		}
		if limits.MaxValueLength > 0 && utf8.RuneCountInString(t.Value) > limits.MaxValueLength {
			return TagValueTooLong, t.Key
		}
	}
	return OK, ""
}

// CheckStringTags validates a string-map tag form against the limits,
// using the same canonical order as CheckTags. Map iteration order is
// normalised by walking the keys in sorted order so the reported
// violation and offending key are deterministic.
//
// Alongside the violation it returns the offending tag key, or the empty
// string when there is no offender to name; consumers must format error
// messages from this finding instead of re-deriving the offender.
func CheckStringTags(tags map[string]string, limits TagLimits) (Violation, string) {
	if limits.MaxCount > 0 && len(tags) > limits.MaxCount {
		return TooManyTags, ""
	}
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if violatesReservedPrefix(k, limits) {
			return ReservedTagKey, k
		}
		if limits.MinKeyLength > 0 && utf8.RuneCountInString(k) < limits.MinKeyLength {
			return TagKeyTooShort, k
		}
		if limits.MaxKeyLength > 0 && utf8.RuneCountInString(k) > limits.MaxKeyLength {
			return TagKeyTooLong, k
		}
		if limits.MaxValueLength > 0 && utf8.RuneCountInString(tags[k]) > limits.MaxValueLength {
			return TagValueTooLong, k
		}
	}
	return OK, ""
}

func violatesReservedPrefix(key string, limits TagLimits) bool {
	if limits.ReservedPrefix == "" {
		return false
	}
	if limits.ReservedCaseSensitive {
		return strings.HasPrefix(key, limits.ReservedPrefix)
	}
	return strings.HasPrefix(strings.ToLower(key), strings.ToLower(limits.ReservedPrefix))
}
