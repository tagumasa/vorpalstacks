// Package bucketname holds the AWS general-purpose S3 bucket naming
// constraints that AWS restates across services whose APIs reference S3
// buckets: S3 itself, CloudTrail trail destinations and Timestream
// scheduled-query error reporting. Validate applies the full naming rule
// set so every referencing service enforces the same contract: 3 to 63
// characters of lowercase letters, digits, hyphens and periods, beginning
// and ending with a letter or digit, no two adjacent periods, no ".-"
// or "-." pairs, no IP-formatted names, and none of the reserved
// prefixes (xn--, sthree-, amzn-s3-demo-) or suffixes (-s3alias,
// --ol-s3, .mrap, --x-s3, --table-s3).
package bucketname

import (
	"regexp"
	"strings"
)

const (
	// MinLength is the minimum S3 bucket name length.
	MinLength = 3

	// MaxLength is the maximum S3 bucket name length.
	MaxLength = 63
)

// DNSPattern matches a general-purpose S3 bucket name: lowercase letters,
// digits, hyphens and periods, starting and ending with a letter or digit.
var DNSPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)

// ipAddressPattern matches IPv4 address forms, which AWS forbids as
// bucket names (for example 192.168.5.4).
var ipAddressPattern = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)

// reservedPrefixes are bucket name prefixes AWS reserves for its own
// use; creating a bucket with one of them is rejected.
var reservedPrefixes = []string{"xn--", "sthree-", "amzn-s3-demo-"}

// reservedSuffixes are bucket name suffixes reserved for access point
// aliases, Object Lambda aliases, Multi-Region Access Points, directory
// buckets and table buckets; creating a bucket with one of them is
// rejected.
var reservedSuffixes = []string{"-s3alias", "--ol-s3", ".mrap", "--x-s3", "--table-s3"}

// Validate reports whether name is a valid general-purpose S3 bucket
// name per the AWS bucket naming rules. It is the single entry point for
// services that reference bucket names; the exported pattern and length
// constants remain available for error messages.
func Validate(name string) bool {
	if len(name) < MinLength || len(name) > MaxLength {
		return false
	}
	if !DNSPattern.MatchString(name) {
		return false
	}
	if ipAddressPattern.MatchString(name) {
		return false
	}
	if strings.Contains(name, "..") || strings.Contains(name, ".-") || strings.Contains(name, "-.") {
		return false
	}
	for _, prefix := range reservedPrefixes {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}
	for _, suffix := range reservedSuffixes {
		if strings.HasSuffix(name, suffix) {
			return false
		}
	}
	return true
}
