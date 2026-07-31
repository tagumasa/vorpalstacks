package acm

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

// domainNamePattern validates domain names per the Smithy DomainNameString
// constraints: optional wildcard prefix, dot-separated labels where each label
// starts and ends with an alphanumeric character (hyphens allowed in the middle).
// RE2-compatible rewrite of the Smithy Perl-syntax pattern.
var domainNamePattern = regexp.MustCompile(`^(\*\.)?([A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?\.)+[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?$`)

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

func isValidRevocationReason(reason string) bool {
	return validRevocationReasons[reason]
}

func formatEpochSeconds(t time.Time) float64 {
	return float64(t.Unix()) + float64(t.Nanosecond())/1e9
}

func getMaxItems(params map[string]interface{}) int {
	return pagination.GetMaxItems(params, pagination.DefaultMaxItems)
}

func parseCertificateArn(params map[string]interface{}, paramName string) (string, error) {
	arn := request.GetStringParam(params, paramName)
	if arn == "" {
		return "", awserrors.NewValidationException(paramName + " is required")
	}
	if !strings.HasPrefix(arn, "arn:") || !strings.Contains(arn, ":certificate/") {
		return "", NewInvalidArnError(arn)
	}
	return arn, nil
}

func parseDomainName(params map[string]interface{}) (string, error) {
	domain := request.GetStringParam(params, "DomainName")
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

func parseCertificateAuthorityArn(params map[string]interface{}) string {
	return request.GetStringParam(params, "CertificateAuthorityArn")
}

func parseStringList(params map[string]interface{}, key string) []string {
	var result []string
	if raw, ok := params[key]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					result = append(result, s)
				}
			}
		}
	}
	return result
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
		"CreatedAt":          formatEpochSeconds(cert.CreatedAt),
		"Options":            certificateOptionsToResponse(cert.Options),
	}

	if !cert.NotBefore.IsZero() {
		result["NotBefore"] = formatEpochSeconds(cert.NotBefore)
	}
	if !cert.NotAfter.IsZero() {
		result["NotAfter"] = formatEpochSeconds(cert.NotAfter)
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
		result["IssuedAt"] = formatEpochSeconds(cert.IssuedAt)
	}

	if !cert.ImportedAt.IsZero() {
		result["ImportedAt"] = formatEpochSeconds(cert.ImportedAt)
	}

	if !cert.RevokedAt.IsZero() {
		result["RevokedAt"] = formatEpochSeconds(cert.RevokedAt)
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
			"UpdatedAt":               formatEpochSeconds(cert.RenewalSummary.UpdatedAt),
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
