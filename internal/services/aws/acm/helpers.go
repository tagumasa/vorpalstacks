package acm

import (
	"fmt"
	"regexp"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// domainNamePattern validates domain names per the Smithy DomainNameString
// constraints: optional wildcard prefix, dot-separated labels where each label
// starts and ends with an alphanumeric character (hyphens allowed in the middle).
// RE2-compatible rewrite of the Smithy Perl-syntax pattern.
var domainNamePattern = regexp.MustCompile(`^(\*\.)?([A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?\.)+[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?$`)

// idempotencyTokenPattern validates IdempotencyToken per the Smithy
// @pattern(^\w+$) constraint.
var idempotencyTokenPattern = regexp.MustCompile(`^\w+$`)

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

// validExtendedKeyUsageNames contains the Smithy ExtendedKeyUsageName enum values.
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

// validateEnumList validates that every value in the list is a member of the
// given enum set. Returns nil if the list is empty.
func validateEnumList(values []string, paramName string, validSet map[string]bool) error {
	for _, v := range values {
		if !validSet[v] {
			return NewInvalidParameterError(fmt.Sprintf("Invalid %s: %s", paramName, v))
		}
	}
	return nil
}

// validateSingleEnum validates a single-value enum parameter. Returns nil if
// the value is empty (parameter not provided).
func validateSingleEnum(value, paramName string, validSet map[string]bool) error {
	if value == "" {
		return nil
	}
	if !validSet[value] {
		return NewInvalidParameterError(fmt.Sprintf("Invalid %s: %s", paramName, value))
	}
	return nil
}

func isValidRevocationReason(reason string) bool {
	return validRevocationReasons[reason]
}

func getMaxItems(params map[string]interface{}) int {
	return pagination.GetMaxItems(params, pagination.DefaultMaxItems)
}

func parseCertificateArn(params map[string]interface{}, paramName string) (string, error) {
	arn := request.GetStringParam(params, paramName)
	if arn == "" {
		return "", awserrors.NewValidationException(paramName + " is required")
	}
	if err := validateCertificateArn(arn); err != nil {
		return "", err
	}
	return arn, nil
}

// validateCertificateArn validates that the given string conforms to the
// ACM ARN structure: arn:partition:acm:region:account-id:certificate/<id>.
// Shared by the HTTP API (via parseCertificateArn) and the gRPC admin handler.
func validateCertificateArn(arn string) error {
	// ARN structure: arn:partition:service:region:account-id:resource
	parts := strings.Split(arn, ":")
	if len(parts) != 6 || parts[0] != "arn" || parts[2] != "acm" {
		return NewInvalidArnError(arn)
	}
	if !strings.HasPrefix(parts[5], "certificate/") || len(parts[5]) <= len("certificate/") {
		return NewInvalidArnError(arn)
	}
	return nil
}

// validateDomainName normalises and validates a domain name string against the
// ACM DomainNameString constraints (253 chars total, 63 chars per label,
// RFC-1123 hostname pattern). Wildcard labels ("*") are allowed in the
// leftmost position. Both the HTTP API (parseDomainName) and the gRPC admin
// handler call this function to ensure identical validation logic.
func validateDomainName(domain string) (string, error) {
	if domain == "" {
		return "", awserrors.NewValidationException("DomainName is required")
	}
	domain = strings.ToLower(domain)
	if len(domain) > 253 || !domainNamePattern.MatchString(domain) {
		return "", awserrors.NewValidationException("Invalid domain name: " + domain)
	}
	// Validate per-label 63 char limit per Smithy pattern.
	for _, label := range strings.Split(domain, ".") {
		if label == "*" {
			continue
		}
		if len(label) > 63 {
			return "", awserrors.NewValidationException("Invalid domain name: label exceeds 63 characters: " + label)
		}
	}
	return domain, nil
}

// validateDomainValidationFields validates the DomainName and ValidationDomain
// fields of a DomainValidationOption against the Smithy constraints. Shared by
// the HTTP API (applyUserDomainValidationOptions) and the admin handler to
// ensure consistent validation across both paths.
func validateDomainValidationFields(domainName, validationDomain string) error {
	if domainName == "" {
		return NewInvalidDomainValidationOptionsException("DomainName is required in DomainValidationOptions")
	}
	if validationDomain == "" {
		return NewInvalidDomainValidationOptionsException("ValidationDomain is required in DomainValidationOptions")
	}
	if !domainNamePattern.MatchString(domainName) {
		return NewInvalidDomainValidationOptionsException(fmt.Sprintf("Invalid DomainName in DomainValidationOptions: %s", domainName))
	}
	if !domainNamePattern.MatchString(validationDomain) {
		return NewInvalidDomainValidationOptionsException(fmt.Sprintf("Invalid ValidationDomain in DomainValidationOptions: %s", validationDomain))
	}
	return nil
}

func parseDomainName(params map[string]interface{}) (string, error) {
	return validateDomainName(request.GetStringParam(params, "DomainName"))
}

func parseValidationMethod(params map[string]interface{}) (string, error) {
	method := request.GetStringParam(params, "ValidationMethod")
	if method == "" {
		return "DNS", nil
	}
	switch method {
	case "DNS", "EMAIL", "HTTP":
		return method, nil
	default:
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid ValidationMethod: %s. Valid values are DNS, EMAIL, HTTP.", method))
	}
}

func parseKeyAlgorithm(params map[string]interface{}) (string, error) {
	algo := request.GetStringParam(params, "KeyAlgorithm")
	if algo == "" {
		return "RSA_2048", nil
	}
	switch algo {
	case "RSA_1024", "RSA_2048", "RSA_3072", "RSA_4096",
		"EC_prime256v1", "EC_secp384r1", "EC_secp521r1":
		return algo, nil
	default:
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid KeyAlgorithm: %s. Valid values are RSA_1024, RSA_2048, RSA_3072, RSA_4096, EC_prime256v1, EC_secp384r1, EC_secp521r1.", algo))
	}
}

func parseCertificateTransparencyLoggingPreference(params map[string]interface{}) (string, error) {
	pref := request.GetStringParam(params, "CertificateTransparencyLoggingPreference")
	if pref == "" {
		return "ENABLED", nil
	}
	switch pref {
	case "ENABLED", "DISABLED":
		return pref, nil
	default:
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid CertificateTransparencyLoggingPreference: %s. Valid values are ENABLED, DISABLED.", pref))
	}
}

func parseExportOption(params map[string]interface{}) (string, error) {
	export := request.GetStringParam(params, "Export")
	if export == "" {
		return "DISABLED", nil
	}
	switch export {
	case "ENABLED", "DISABLED":
		return export, nil
	default:
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid Export: %s. Valid values are ENABLED, DISABLED.", export))
	}
}

var validManagedByValues = map[string]bool{
	"CLOUDFRONT": true,
}

// validExportOptionValues contains the Smithy CertificateExport enum values
// used as a filter in ListCertificates.Includes.exportOption.
var validExportOptionValues = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

// validateIdempotencyToken validates an IdempotencyToken string against the
// Smithy constraints (@length 1-32, @pattern ^\w+$). Shared by the HTTP API
// (RequestCertificate) and the admin handler. No-op when token is empty.
func validateIdempotencyToken(token string) error {
	if token == "" {
		return nil
	}
	if len(token) > 32 {
		return awserrors.NewValidationException("IdempotencyToken must not exceed 32 characters")
	}
	if !idempotencyTokenPattern.MatchString(token) {
		return awserrors.NewValidationException("IdempotencyToken must contain only alphanumeric characters and underscores")
	}
	return nil
}

// validateACMTags validates tag keys and values against the Smithy
// constraints: TagKey @length(1-128) + pattern, TagValue @length(max 256)
// + pattern. Returns InvalidTagException on violation.
var tagKeyPattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]*)$`)

func validateACMTags(tags []types.Tag) error {
	for _, t := range tags {
		if len(t.Key) < 1 || len(t.Key) > 128 {
			return NewInvalidTagException(fmt.Sprintf("Tag key length must be 1-128: got %d", len(t.Key)))
		}
		if len(t.Value) > 256 {
			return NewInvalidTagException(fmt.Sprintf("Tag value length must not exceed 256: got %d", len(t.Value)))
		}
		if !tagKeyPattern.MatchString(t.Key) {
			return NewInvalidTagException(fmt.Sprintf("Tag key contains invalid characters: %s", t.Key))
		}
		if !tagKeyPattern.MatchString(t.Value) {
			return NewInvalidTagException(fmt.Sprintf("Tag value contains invalid characters: %s", t.Value))
		}
	}
	return nil
}

func parseManagedBy(params map[string]interface{}) (string, error) {
	mb := request.GetStringParam(params, "ManagedBy")
	if mb == "" {
		return "", nil
	}
	if !validManagedByValues[mb] {
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid ManagedBy value: %s", mb))
	}
	return mb, nil
}

// pcaArnPattern validates PCA ARNs per the Smithy PcaArn pattern.
var pcaArnPattern = regexp.MustCompile(`^arn:[\w+=/,.@-]+:acm-pca:[\w+=/,.@-]*:[0-9]+:[\w+=,.@-]+(/[\w+=,.@-]+)*$`)

// validateCertificateAuthorityArn validates a PCA ARN string against the
// Smithy PcaArn constraints. Shared by the HTTP API (via
// parseCertificateAuthorityArn) and the admin handler.
func validateCertificateAuthorityArn(arn string) (string, error) {
	if arn == "" {
		return "", nil
	}
	// Smithy PcaArn: @length(20-2048) + pattern.
	if len(arn) < 20 || len(arn) > 2048 {
		return "", awserrors.NewInvalidParameterException("CertificateAuthorityArn length must be 20-2048")
	}
	if !pcaArnPattern.MatchString(arn) {
		return "", awserrors.NewInvalidParameterException("Invalid CertificateAuthorityArn format")
	}
	return arn, nil
}

func parseCertificateAuthorityArn(params map[string]interface{}) (string, error) {
	return validateCertificateAuthorityArn(request.GetStringParam(params, "CertificateAuthorityArn"))
}

func certificateOptionsToResponse(opts *acmstorelib.CertificateOptions) map[string]interface{} {
	if opts == nil {
		return map[string]interface{}{
			"CertificateTransparencyLoggingPreference": "ENABLED",
			"Export": "DISABLED",
		}
	}
	return map[string]interface{}{
		"CertificateTransparencyLoggingPreference": opts.CertificateTransparencyLoggingPreference,
		"Export": opts.Export,
	}
}

func domainValidationToResponse(dv *acmstorelib.DomainValidation) map[string]interface{} {
	if dv == nil {
		return nil
	}
	result := map[string]interface{}{
		"DomainName":       dv.DomainName,
		"ValidationMethod": dv.ValidationMethod,
		"ValidationStatus": dv.ValidationStatus,
		"ValidationDomain": dv.ValidationDomain,
	}
	if len(dv.ValidationEmails) > 0 {
		result["ValidationEmails"] = dv.ValidationEmails
	}
	if dv.ResourceRecord != nil {
		result["ResourceRecord"] = map[string]interface{}{
			"Name":  dv.ResourceRecord.Name,
			"Type":  dv.ResourceRecord.Type,
			"Value": dv.ResourceRecord.Value,
		}
	}
	if dv.HttpRedirect != nil {
		result["HttpRedirect"] = map[string]interface{}{
			"StatusCode": dv.HttpRedirect.StatusCode,
			"Location":   dv.HttpRedirect.Location,
		}
	}
	return result
}

func certificateToDetailResponse(cert *acmstorelib.Certificate) map[string]interface{} {
	result := map[string]interface{}{
		"CertificateArn":     cert.CertificateArn,
		"DomainName":         cert.DomainName,
		"Subject":            cert.Subject,
		"Issuer":             cert.Issuer,
		"Serial":             cert.Serial,
		"Status":             cert.Status,
		"Type":               cert.Type,
		"KeyAlgorithm":       cert.KeyAlgorithm,
		"SignatureAlgorithm": cert.SignatureAlgorithm,
		"RenewalEligibility": cert.RenewalEligibility,
		"CreatedAt":          timeutils.FormatEpochSeconds(cert.CreatedAt),
		"Options":            certificateOptionsToResponse(cert.Options),
	}

	if !cert.NotBefore.IsZero() {
		result["NotBefore"] = timeutils.FormatEpochSeconds(cert.NotBefore)
	}
	if !cert.NotAfter.IsZero() {
		result["NotAfter"] = timeutils.FormatEpochSeconds(cert.NotAfter)
	}

	if len(cert.InUseBy) > 0 {
		result["InUseBy"] = cert.InUseBy
	}

	if len(cert.SubjectAlternativeNames) > 0 {
		result["SubjectAlternativeNames"] = cert.SubjectAlternativeNames
	}

	if len(cert.DomainValidationOptions) > 0 {
		dvos := make([]interface{}, len(cert.DomainValidationOptions))
		for i, dv := range cert.DomainValidationOptions {
			dvos[i] = domainValidationToResponse(dv)
		}
		result["DomainValidationOptions"] = dvos
	}

	if !cert.IssuedAt.IsZero() {
		result["IssuedAt"] = timeutils.FormatEpochSeconds(cert.IssuedAt)
	}

	if !cert.ImportedAt.IsZero() {
		result["ImportedAt"] = timeutils.FormatEpochSeconds(cert.ImportedAt)
	}

	if !cert.RevokedAt.IsZero() {
		result["RevokedAt"] = timeutils.FormatEpochSeconds(cert.RevokedAt)
	}

	if cert.RevocationReason != "" {
		result["RevocationReason"] = cert.RevocationReason
	}

	if cert.FailureReason != "" {
		result["FailureReason"] = cert.FailureReason
	}

	if cert.CertificateAuthorityArn != "" {
		result["CertificateAuthorityArn"] = cert.CertificateAuthorityArn
	}

	if len(cert.KeyUsages) > 0 {
		kus := make([]interface{}, len(cert.KeyUsages))
		for i, ku := range cert.KeyUsages {
			kus[i] = map[string]interface{}{"Name": ku.Name}
		}
		result["KeyUsages"] = kus
	}

	if len(cert.ExtendedKeyUsages) > 0 {
		ekus := make([]interface{}, len(cert.ExtendedKeyUsages))
		for i, eku := range cert.ExtendedKeyUsages {
			ekus[i] = map[string]interface{}{"Name": eku.Name, "OID": eku.OID}
		}
		result["ExtendedKeyUsages"] = ekus
	}

	if cert.RenewalSummary != nil {
		var renewalDvos []interface{}
		if len(cert.RenewalSummary.DomainValidationOptions) > 0 {
			renewalDvos = make([]interface{}, len(cert.RenewalSummary.DomainValidationOptions))
			for i, dv := range cert.RenewalSummary.DomainValidationOptions {
				renewalDvos[i] = domainValidationToResponse(dv)
			}
		}
		result["RenewalSummary"] = map[string]interface{}{
			"RenewalStatus":           cert.RenewalSummary.RenewalStatus,
			"RenewalStatusReason":     cert.RenewalSummary.RenewalStatusReason,
			"DomainValidationOptions": renewalDvos,
			"UpdatedAt":               timeutils.FormatEpochSeconds(cert.RenewalSummary.UpdatedAt),
		}
	}

	if cert.ManagedBy != "" {
		result["ManagedBy"] = cert.ManagedBy
	}
	if cert.CertificateKeyPairOrigin != "" {
		result["CertificateKeyPairOrigin"] = cert.CertificateKeyPairOrigin
	}

	return result
}

func certificateSummaryToResponse(summary *acmstorelib.CertificateSummary) map[string]interface{} {
	result := map[string]interface{}{
		"CertificateArn":     summary.CertificateArn,
		"DomainName":         summary.DomainName,
		"Status":             summary.Status,
		"Type":               summary.Type,
		"RenewalEligibility": summary.RenewalEligibility,
		"KeyAlgorithm":       summary.KeyAlgorithm,
		"InUse":              summary.InUse,
		"Exported":           summary.Exported,
	}

	if len(summary.SubjectAlternativeNameSummaries) > 0 {
		result["SubjectAlternativeNameSummaries"] = summary.SubjectAlternativeNameSummaries
	}
	if summary.HasAdditionalSubjectAlternativeNames {
		result["HasAdditionalSubjectAlternativeNames"] = summary.HasAdditionalSubjectAlternativeNames
	}
	if len(summary.KeyUsages) > 0 {
		result["KeyUsages"] = summary.KeyUsages
	}
	if len(summary.ExtendedKeyUsages) > 0 {
		result["ExtendedKeyUsages"] = summary.ExtendedKeyUsages
	}
	if summary.NotBefore != 0 {
		result["NotBefore"] = summary.NotBefore
	}
	if summary.NotAfter != 0 {
		result["NotAfter"] = summary.NotAfter
	}
	if summary.CreatedAt != 0 {
		result["CreatedAt"] = summary.CreatedAt
	}
	if summary.IssuedAt != 0 {
		result["IssuedAt"] = summary.IssuedAt
	}
	if summary.ImportedAt != 0 {
		result["ImportedAt"] = summary.ImportedAt
	}
	if summary.ExportOption != "" {
		result["ExportOption"] = summary.ExportOption
	}
	if summary.ManagedBy != "" {
		result["ManagedBy"] = summary.ManagedBy
	}
	if summary.CertificateKeyPairOrigin != "" {
		result["CertificateKeyPairOrigin"] = summary.CertificateKeyPairOrigin
	}
	if summary.RevokedAt != 0 {
		result["RevokedAt"] = summary.RevokedAt
	}

	return result
}

func listResultToResponse(result *acmstorelib.CertificateListResult) map[string]interface{} {
	certs := make([]interface{}, len(result.Certificates))
	for i, cert := range result.Certificates {
		certs[i] = certificateSummaryToResponse(cert)
	}
	resp := map[string]interface{}{
		"CertificateSummaryList": certs,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp
}
