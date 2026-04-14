// Package acm provides ACM (AWS Certificate Manager) service operations for vorpalstacks.
package acm

import (
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
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

func filteredListToResponse(summaries []*acmstorelib.CertificateSummary, marker string, maxItems int) map[string]interface{} {
	if maxItems <= 0 {
		maxItems = 100
	}
	started := marker == ""
	var filtered []*acmstorelib.CertificateSummary
	for _, s := range summaries {
		if !started {
			if s.CertificateArn > marker {
				started = true
			}
			if !started {
				continue
			}
		}
		filtered = append(filtered, s)
	}
	certs := make([]interface{}, 0, len(filtered))
	hasMore := false
	var lastArn string
	for i, s := range filtered {
		if i >= maxItems {
			hasMore = true
			break
		}
		certs = append(certs, certificateSummaryToResponse(s))
		lastArn = s.CertificateArn
	}
	nextToken := ""
	if hasMore && lastArn != "" {
		nextToken = lastArn
	}
	return map[string]interface{}{
		"NextToken":              nextToken,
		"CertificateSummaryList": certs,
	}
}
