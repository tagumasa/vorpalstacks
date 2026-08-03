package acm

import (
	"regexp"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

// ---------------------------------------------------------------------------
// Smithy-derived patterns
// ---------------------------------------------------------------------------

// domainNamePattern validates domain names per the Smithy DomainNameString
// constraints: optional wildcard prefix, dot-separated labels where each label
// starts and ends with an alphanumeric character (hyphens allowed in the
// middle).  RE2-compatible rewrite of the Smithy Perl-syntax pattern.
var domainNamePattern = regexp.MustCompile(`^(\*\.)?([A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?\.)+[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?$`)

// idempotencyTokenPattern validates IdempotencyToken per the Smithy
// @pattern(^\w+$) constraint.
var idempotencyTokenPattern = regexp.MustCompile(`^\w+$`)

// tagKeyPattern validates tag keys and values per the Smithy TagKey/TagValue
// pattern constraints.
var tagKeyPattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]*)$`)

// pcaArnPattern validates PCA ARNs per the Smithy PcaArn pattern.
var pcaArnPattern = regexp.MustCompile(`^arn:[\w+=/,.@-]+:acm-pca:[\w+=/,.@-]*:[0-9]+:[\w+=,.@-]+(/[\w+=,.@-]+)*$`)

// ---------------------------------------------------------------------------
// Smithy-derived enum sets
// ---------------------------------------------------------------------------

// validRevocationReasons contains the Smithy RevocationReason enum values.
var validRevocationReasons = map[string]bool{
	"UNSPECIFIED":            true,
	"KEY_COMPROMISE":         true,
	"CA_COMPROMISE":          true,
	"AFFILIATION_CHANGED":    true,
	"SUPERCEDED":             true,
	"SUPERSEDED":             true,
	"CESSATION_OF_OPERATION": true,
	"CERTIFICATE_HOLD":       true,
	"REMOVE_FROM_CRL":        true,
	"PRIVILEGE_WITHDRAWN":    true,
	"A_A_COMPROMISE":         true,
}

// validKeyAlgorithmsMap contains the Smithy KeyAlgorithm enum values.
var validKeyAlgorithmsMap = map[string]bool{
	"RSA_1024":      true,
	"RSA_2048":      true,
	"RSA_3072":      true,
	"RSA_4096":      true,
	"EC_prime256v1": true,
	"EC_secp384r1":  true,
	"EC_secp521r1":  true,
}

// validCertificateStatuses contains the Smithy CertificateStatus enum values.
var validCertificateStatuses = map[string]bool{
	"PENDING_VALIDATION":   true,
	"ISSUED":               true,
	"INACTIVE":             true,
	"EXPIRED":              true,
	"VALIDATION_TIMED_OUT": true,
	"REVOKED":              true,
	"FAILED":               true,
}

// validKeyPairOrigins contains the Smithy CertificateKeyPairOrigin enum values.
var validKeyPairOrigins = map[string]bool{
	"AWS_MANAGED":       true,
	"ACME":              true,
	"CUSTOMER_PROVIDED": true,
}

// validRenewalStatuses contains the Smithy RenewalStatus enum values.
var validRenewalStatuses = map[string]bool{
	"PENDING_AUTO_RENEWAL": true,
	"PENDING_VALIDATION":   true,
	"SUCCESS":              true,
	"FAILED":               true,
}

// validCertTypes contains the Smithy CertificateType enum values.
var validCertTypes = map[string]bool{
	"AMAZON_ISSUED": true,
	"PRIVATE":       true,
	"IMPORTED":      true,
}

// validKeyUsageNames contains the Smithy KeyUsageName enum values.
var validKeyUsageNames = map[string]bool{
	"DIGITAL_SIGNATURE":   true,
	"NON_REPUDIATION":     true,
	"KEY_ENCIPHERMENT":    true,
	"DATA_ENCIPHERMENT":   true,
	"KEY_AGREEMENT":       true,
	"CERTIFICATE_SIGNING": true,
	"CRL_SIGNING":         true,
	"ENCIPHER_ONLY":       true,
	"DECIPHER_ONLY":       true,
	"ANY":                 true,
	"CUSTOM":              true,
}

// validExtendedKeyUsageNames contains the Smithy ExtendedKeyUsageName enum
// values.
var validExtendedKeyUsageNames = map[string]bool{
	"TLS_WEB_SERVER_AUTHENTICATION": true,
	"TLS_WEB_CLIENT_AUTHENTICATION": true,
	"CODE_SIGNING":                  true,
	"EMAIL_PROTECTION":              true,
	"TIME_STAMPING":                 true,
	"OCSP_SIGNING":                  true,
	"IPSEC_END_SYSTEM":              true,
	"IPSEC_TUNNEL":                  true,
	"IPSEC_USER":                    true,
	"ANY":                           true,
	"NONE":                          true,
	"CUSTOM":                        true,
}

// validManagedByValues contains the Smithy ManagedBy enum values.
var validManagedByValues = map[string]bool{
	"CLOUDFRONT": true,
}

// validExportOptionValues contains the Smithy CertificateExport enum values
// used as a filter in ListCertificates.Includes.exportOption.
var validExportOptionValues = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

// ---------------------------------------------------------------------------
// Bool validators
// ---------------------------------------------------------------------------

// isValidRevocationReason returns true when the reason is a valid Smithy
// RevocationReason enum value.
func isValidRevocationReason(reason string) bool {
	return validRevocationReasons[reason]
}

// isValidCertificateArn returns true when the string conforms to the ACM ARN
// structure: arn:partition:acm:region:account-id:certificate/<id>.
// Smithy Arn: @length(20-2048) + pattern.
func isValidCertificateArn(arn string) bool {
	if len(arn) < 20 || len(arn) > 2048 {
		return false
	}
	parts := strings.Split(arn, ":")
	if len(parts) != 6 || parts[0] != "arn" || parts[2] != "acm" {
		return false
	}
	return strings.HasPrefix(parts[5], "certificate/") && len(parts[5]) > len("certificate/")
}

// isValidDomainName returns true when the domain conforms to the ACM
// DomainNameString constraints (253 chars total, 63 chars per label,
// RFC-1123 hostname pattern with optional wildcard prefix).
func isValidDomainName(domain string) bool {
	if domain == "" {
		return false
	}
	if len(domain) > 253 || !domainNamePattern.MatchString(domain) {
		return false
	}
	for _, label := range strings.Split(domain, ".") {
		if label == "*" {
			continue
		}
		if len(label) > 63 {
			return false
		}
	}
	return true
}

// isValidValidationMethod returns true when the method is a valid Smithy
// ValidationMethod enum value (DNS, EMAIL, HTTP).
func isValidValidationMethod(method string) bool {
	return method == "DNS" || method == "EMAIL" || method == "HTTP"
}

// isValidKeyAlgorithm returns true when the algorithm is a valid Smithy
// KeyAlgorithm enum value.
func isValidKeyAlgorithm(algo string) bool {
	return validKeyAlgorithmsMap[algo]
}

// isValidCertificateTransparencyLoggingPreference returns true when the
// preference is a valid Smithy CertificateTransparencyLoggingPreference enum
// value (ENABLED, DISABLED).
func isValidCertificateTransparencyLoggingPreference(pref string) bool {
	return pref == "ENABLED" || pref == "DISABLED"
}

// isValidExportOption returns true when the export option is a valid Smithy
// CertificateExport enum value (ENABLED, DISABLED).
func isValidExportOption(export string) bool {
	return export == "ENABLED" || export == "DISABLED"
}

// isValidManagedBy returns true when the value is a valid Smithy ManagedBy
// enum value.
func isValidManagedBy(mb string) bool {
	return validManagedByValues[mb]
}

// isValidFilterString returns true when the string conforms to the Smithy
// FilterString constraints: @length(1-256). Used by search filter values
// (CommonNameFilter.Value, DnsNameFilter.Value).
func isValidFilterString(s string) bool {
	return len(s) >= 1 && len(s) <= 256
}

// validateNextToken validates an opaque NextToken against the Smithy
// constraints: @length(1-10000). Empty token is allowed (means "first page").
// Non-empty token must be within length range.
func validateNextToken(token string) error {
	if token == "" {
		return nil
	}
	if len(token) > 10000 {
		return awserrors.NewValidationException("NextToken exceeds maximum length of 10000")
	}
	return nil
}
