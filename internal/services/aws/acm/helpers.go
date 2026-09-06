package acm

import (
	"fmt"
	"strings"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/paramvalidation"
	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
	"vorpalstacks/internal/utils/timeutils"
)

// validateEnumList validates that every value in the list is a member of the
// given enum set. Returns nil if the list is empty. The callers serve
// operations that declare ValidationException (not InvalidParameterException),
// so violations surface as ValidationException.
func validateEnumList(values []string, paramName string, validSet map[string]bool) error {
	return paramvalidation.EnumList(paramName, values, validSet,
		func(field, value string) error {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid %s: %s", field, value))
		})
}

// validateSingleEnum validates a single-value enum parameter. Returns nil
// if the value is empty (parameter not provided).
func validateSingleEnum(value, paramName string, validSet map[string]bool) error {
	return paramvalidation.EnumValue(paramName, value, validSet,
		func(field, value string) error {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid %s: %s", field, value))
		})
}

func parseCertificateArn(params map[string]interface{}, paramName string) (string, error) {
	arn := request.GetStringParam(params, paramName)
	if err := requireCertificateArn(arn, paramName); err != nil {
		return "", err
	}
	return arn, nil
}

// validateCertificateArn validates that the given string conforms to the
// ACM ARN structure: arn:partition:acm:region:account-id:certificate/<id>.
// Shared by the HTTP API (via parseCertificateArn) and the gRPC admin handler.
func validateCertificateArn(arn string) error {
	if !isValidCertificateArn(arn) {
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
	normalised := strings.ToLower(domain)
	if !isValidDomainName(normalised) {
		for _, label := range strings.Split(normalised, ".") {
			if label == "*" {
				continue
			}
			if len(label) > 63 {
				return "", awserrors.NewValidationException("Invalid domain name: label exceeds 63 characters: " + label)
			}
		}
		return "", awserrors.NewValidationException("Invalid domain name: " + normalised)
	}
	return normalised, nil
}

// validateDomainValidationFields validates the DomainName and ValidationDomain
// fields of a DomainValidationOption against the Smithy constraints (Smithy
// DomainNameString: @length 1-253 + pattern). Shared by the HTTP API core
// function and the admin handler to ensure consistent validation.
func validateDomainValidationFields(domainName, validationDomain string) error {
	if domainName == "" {
		return NewInvalidDomainValidationOptionsException("DomainName is required in DomainValidationOptions")
	}
	if validationDomain == "" {
		return NewInvalidDomainValidationOptionsException("ValidationDomain is required in DomainValidationOptions")
	}
	if !isValidDomainName(domainName) {
		return NewInvalidDomainValidationOptionsException(fmt.Sprintf("Invalid DomainName in DomainValidationOptions: %s", domainName))
	}
	if !isValidDomainName(validationDomain) {
		return NewInvalidDomainValidationOptionsException(fmt.Sprintf("Invalid ValidationDomain in DomainValidationOptions: %s", validationDomain))
	}
	return nil
}

func parseCertificateTransparencyLoggingPreference(params map[string]interface{}) (string, error) {
	pref := request.GetStringParam(params, "CertificateTransparencyLoggingPreference")
	if pref == "" {
		return "ENABLED", nil
	}
	if !isValidCertificateTransparencyLoggingPreference(pref) {
		return "", awserrors.NewValidationException(fmt.Sprintf("Invalid CertificateTransparencyLoggingPreference: %s. Valid values are ENABLED, DISABLED.", pref))
	}
	return pref, nil
}

func parseExportOption(params map[string]interface{}) (string, error) {
	export := request.GetStringParam(params, "Export")
	if export == "" {
		return "DISABLED", nil
	}
	if !isValidExportOption(export) {
		return "", awserrors.NewValidationException(fmt.Sprintf("Invalid Export: %s. Valid values are ENABLED, DISABLED.", export))
	}
	return export, nil
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
func validateACMTags(tags []types.Tag) error {
	for _, t := range tags {
		if keyLen := utf8.RuneCountInString(t.Key); keyLen < 1 || keyLen > types.MaxTagKeyLength {
			return NewInvalidTagException(fmt.Sprintf("Tag key length must be 1-%d: got %d", types.MaxTagKeyLength, keyLen))
		}
		if valueLen := utf8.RuneCountInString(t.Value); valueLen > types.MaxTagValueLength {
			return NewInvalidTagException(fmt.Sprintf("Tag value length must not exceed %d: got %d", types.MaxTagValueLength, valueLen))
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

// validateCertificateAuthorityArn validates a PCA ARN string against the
// Smithy PcaArn constraints. Shared by the HTTP API (via the core function)
// and the admin handler.
func validateCertificateAuthorityArn(arn string) (string, error) {
	if arn == "" {
		return "", nil
	}
	if len(arn) < 20 || len(arn) > 2048 {
		return "", awserrors.NewInvalidParameterException("CertificateAuthorityArn length must be 20-2048")
	}
	if !pcaArnPattern.MatchString(arn) {
		return "", awserrors.NewInvalidParameterException("Invalid CertificateAuthorityArn format")
	}
	return arn, nil
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

func certificateSummaryToResponse(summary *CertSummaryOut) map[string]interface{} {
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

func listResultToResponse(result *ListCertificatesResult) map[string]interface{} {
	certs := make([]interface{}, len(result.Certificates))
	for i := range result.Certificates {
		certs[i] = certificateSummaryToResponse(&result.Certificates[i])
	}
	resp := map[string]interface{}{
		"CertificateSummaryList": certs,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp
}
