// Package acm provides ACM (AWS Certificate Manager) service operations for vorpalstacks.
package acm

import (
	"strings"
	"time"

	acmpb "vorpalstacks/internal/pb/aws/acm"
	"vorpalstacks/internal/services/aws/common/pagination"
	"vorpalstacks/internal/services/aws/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

func formatEpochSeconds(t time.Time) float64 {
	return float64(t.Unix()) + float64(t.Nanosecond())/1e9
}

func getMaxItems(params map[string]interface{}) int {
	return pagination.GetMaxItems(params, pagination.DefaultMaxItems)
}

func parseCertificateArn(params map[string]interface{}, paramName string) (string, error) {
	arn := request.GetStringParam(params, paramName)
	if arn == "" {
		return "", NewValidationException(paramName + " is required")
	}
	return arn, nil
}

func parseDomainName(params map[string]interface{}) (string, error) {
	domain := request.GetStringParam(params, "DomainName")
	if domain == "" {
		return "", NewValidationException("DomainName is required")
	}
	return strings.ToLower(domain), nil
}

func parseValidationMethod(params map[string]interface{}) string {
	method := request.GetStringParam(params, "ValidationMethod")
	if method == "" {
		return "DNS"
	}
	return method
}

func parseKeyAlgorithm(params map[string]interface{}) string {
	algo := request.GetStringParam(params, "KeyAlgorithm")
	if algo == "" {
		return "RSA_2048"
	}
	return algo
}

func parseCertificateTransparencyLoggingPreference(params map[string]interface{}) string {
	pref := request.GetStringParam(params, "CertificateTransparencyLoggingPreference")
	if pref == "" {
		return "ENABLED"
	}
	return pref
}

func certificateStatusToString(status int32) string {
	switch status {
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_PENDING_VALIDATION):
		return "PENDING_VALIDATION"
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_ISSUED):
		return "ISSUED"
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_INACTIVE):
		return "INACTIVE"
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_EXPIRED):
		return "EXPIRED"
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_VALIDATION_TIMED_OUT):
		return "VALIDATION_TIMED_OUT"
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_REVOKED):
		return "REVOKED"
	case int32(acmpb.CertificateStatus_CERTIFICATE_STATUS_FAILED):
		return "FAILED"
	default:
		return "PENDING_VALIDATION"
	}
}

func keyAlgorithmToString(algo int32) string {
	switch algo {
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_RSA_1024):
		return "RSA_1024"
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_RSA_2048):
		return "RSA_2048"
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_RSA_3072):
		return "RSA_3072"
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_RSA_4096):
		return "RSA_4096"
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_EC_PRIME256V1):
		return "EC_prime256v1"
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_EC_SECP384R1):
		return "EC_secp384r1"
	case int32(acmpb.KeyAlgorithm_KEY_ALGORITHM_EC_SECP521R1):
		return "EC_secp521r1"
	default:
		return "RSA_2048"
	}
}

func certificateTypeToString(certType int32) string {
	switch certType {
	case int32(acmpb.CertificateType_CERTIFICATE_TYPE_IMPORTED):
		return "IMPORTED"
	case int32(acmpb.CertificateType_CERTIFICATE_TYPE_AMAZON_ISSUED):
		return "AMAZON_ISSUED"
	case int32(acmpb.CertificateType_CERTIFICATE_TYPE_PRIVATE):
		return "PRIVATE"
	default:
		return "AMAZON_ISSUED"
	}
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
	return result
}

func certificateToDetailResponse(cert *acmstorelib.Certificate) map[string]interface{} {
	now := time.Now()
	notBefore := cert.NotBefore
	notAfter := cert.NotAfter
	if notBefore.IsZero() {
		notBefore = now
	}
	if notAfter.IsZero() {
		notAfter = now.AddDate(1, 0, 0)
	}

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
		"NotBefore":          formatEpochSeconds(notBefore),
		"NotAfter":           formatEpochSeconds(notAfter),
		"CreatedAt":          formatEpochSeconds(cert.CreatedAt),
		"Options":            certificateOptionsToResponse(cert.Options),
		"InUseBy":            cert.InUseBy,
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

	return result
}

func listCertificatesToResponse(result *acmstorelib.CertificateListResult) map[string]interface{} {
	certs := make([]interface{}, len(result.Certificates))
	for i, cert := range result.Certificates {
		certs[i] = certificateSummaryToResponse(cert)
	}
	return map[string]interface{}{
		"NextToken":              result.NextToken,
		"CertificateSummaryList": certs,
	}
}
