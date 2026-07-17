package acm

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

// SearchCertificates retrieves a list of certificates matching search criteria.
func (s *ACMService) SearchCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters

	nextToken := request.GetStringParam(params, "NextToken")
	maxResults := 500
	if mr := request.GetIntParam(params, "MaxResults"); mr > 0 {
		maxResults = mr
		if maxResults > 500 {
			maxResults = 500
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allCerts, err := stores.certificates.ListAll()
	if err != nil {
		return nil, err
	}

	// Apply filters from the FilterStatement.
	filtered := applyCertificateFilters(allCerts, params)

	// Sort results.
	sortBy := request.GetStringParam(params, "SortBy")
	sortOrder := request.GetStringParam(params, "SortOrder")
	sortCertificateSearchResults(filtered, sortBy, sortOrder)

	// Paginate using offset-based tokens.
	offset := 0
	if nextToken != "" {
		if n, parseErr := parseIntToken(nextToken); parseErr == nil && n >= 0 {
			offset = n
		}
	}
	if offset > len(filtered) {
		offset = len(filtered)
	}

	end := offset + maxResults
	if end > len(filtered) {
		end = len(filtered)
	}

	page := filtered[offset:end]

	results := make([]interface{}, 0, len(page))
	for _, c := range page {
		results = append(results, buildCertificateSearchResult(c))
	}

	resp := map[string]interface{}{
		"Results": results,
	}
	if end < len(filtered) {
		resp["NextToken"] = formatIntToken(end)
	}

	return resp, nil
}

// applyCertificateFilters filters certificates based on FilterStatement
// parameters. Supports CertificateArn exact match, AcmCertificateMetadataFilter
// (Status, Type, KeyAlgorithm, ValidationMethod, Exported, InUse), and
// X509AttributeFilter (KeyAlgorithm, Subject.CommonName, SerialNumber).
func applyCertificateFilters(certs []*acmstorelib.Certificate, params map[string]interface{}) []*acmstorelib.Certificate {
	// Check for direct CertificateArn filter.
	if arnFilter := getFilterValue(params, "CertificateArn"); arnFilter != "" {
		var result []*acmstorelib.Certificate
		for _, c := range certs {
			if c.CertificateArn == arnFilter {
				result = append(result, c)
			}
		}
		return result
	}

	// Extract metadata filter values.
	statuses := getFilterList(params, "AcmCertificateMetadataFilter.Status", "Status")
	types := getFilterList(params, "AcmCertificateMetadataFilter.Type", "Type")
	keyAlgos := getFilterList(params, "X509AttributeFilter.KeyAlgorithm", "KeyAlgorithm")
	validationMethods := getFilterList(params, "AcmCertificateMetadataFilter.ValidationMethod", "ValidationMethod")
	commonName := getFilterValue(params, "Subject.CommonName")
	if commonName == "" {
		commonName = getFilterValue(params, "X509AttributeFilter.Subject.CommonName")
	}

	var filtered []*acmstorelib.Certificate
	for _, c := range certs {
		if len(statuses) > 0 && !containsString(statuses, c.Status) {
			continue
		}
		if len(types) > 0 && !containsString(types, c.Type) {
			continue
		}
		if len(validationMethods) > 0 {
			vmMatched := false
			for _, dvo := range c.DomainValidationOptions {
				if containsString(validationMethods, dvo.ValidationMethod) {
					vmMatched = true
					break
				}
			}
			if !vmMatched {
				continue
			}
		}
		if len(keyAlgos) > 0 && !containsString(keyAlgos, c.KeyAlgorithm) {
			continue
		}
		if commonName != "" && !strings.Contains(strings.ToLower(c.Subject), strings.ToLower(commonName)) {
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}

// getFilterValue extracts a scalar filter value from the FilterStatement
// parameter hierarchy. It tries multiple key patterns.
func getFilterValue(params map[string]interface{}, keys ...string) string {
	prefixes := []string{
		"FilterStatement.Filter.",
		"FilterStatement.",
		"Filter.",
		"",
	}
	for _, prefix := range prefixes {
		for _, key := range keys {
			fullKey := prefix + key
			if v := request.GetStringParam(params, fullKey); v != "" {
				return v
			}
		}
	}
	return ""
}

// getFilterList extracts a list filter value from the FilterStatement
// parameter hierarchy. Tries multiple key patterns.
func getFilterList(params map[string]interface{}, keys ...string) []string {
	prefixes := []string{
		"FilterStatement.Filter.",
		"FilterStatement.",
		"Filter.",
		"",
	}
	for _, prefix := range prefixes {
		for _, key := range keys {
			fullKey := prefix + key
			// Try member.N format.
			if vals := request.GetStringList(params, fullKey); len(vals) > 0 {
				return vals
			}
			// Try direct list.
			if raw, ok := params[fullKey]; ok {
				if arr, ok := raw.([]interface{}); ok {
					var result []string
					for _, v := range arr {
						if s, ok := v.(string); ok {
							result = append(result, s)
						}
					}
					if len(result) > 0 {
						return result
					}
				}
			}
		}
	}
	return nil
}

// sortCertificateSearchResults sorts the slice by the given field and order.
func sortCertificateSearchResults(certs []*acmstorelib.Certificate, sortBy, sortOrder string) {
	ascending := sortOrder != "DESCENDING"
	less := func(i, j int) bool {
		a, b := certs[i], certs[j]
		switch sortBy {
		case "CREATED_AT":
			if ascending {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			return a.CreatedAt.After(b.CreatedAt)
		case "CERTIFICATE_ARN":
			if ascending {
				return a.CertificateArn < b.CertificateArn
			}
			return a.CertificateArn > b.CertificateArn
		case "STATUS":
			if ascending {
				return a.Status < b.Status
			}
			return a.Status > b.Status
		case "TYPE":
			if ascending {
				return a.Type < b.Type
			}
			return a.Type > b.Type
		case "KEY_ALGORITHM":
			if ascending {
				return a.KeyAlgorithm < b.KeyAlgorithm
			}
			return a.KeyAlgorithm > b.KeyAlgorithm
		case "NOT_BEFORE":
			if ascending {
				return a.NotBefore.Before(b.NotBefore)
			}
			return a.NotBefore.After(b.NotBefore)
		case "NOT_AFTER":
			if ascending {
				return a.NotAfter.Before(b.NotAfter)
			}
			return a.NotAfter.After(b.NotAfter)
		case "IMPORTED_AT":
			ai := a.ImportedAt
			bi := b.ImportedAt
			if ascending {
				return ai.Before(bi)
			}
			return ai.After(bi)
		default:
			if ascending {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			return a.CreatedAt.After(b.CreatedAt)
		}
	}
	sort.Slice(certs, less)
}

// buildCertificateSearchResult constructs the response object with
// CertificateArn, CertificateMetadata, and X509Attributes.
func buildCertificateSearchResult(cert *acmstorelib.Certificate) map[string]interface{} {
	result := map[string]interface{}{
		"CertificateArn": cert.CertificateArn,
		"CertificateMetadata": map[string]interface{}{
			"Status":             cert.Status,
			"Type":               cert.Type,
			"KeyAlgorithm":       cert.KeyAlgorithm,
			"RenewalEligibility": cert.RenewalEligibility,
			"NotBefore":          formatEpochSeconds(cert.NotBefore),
			"NotAfter":           formatEpochSeconds(cert.NotAfter),
			"CreatedAt":          formatEpochSeconds(cert.CreatedAt),
			"Exported":           false,
			"InUse":              len(cert.InUseBy) > 0,
			"ExportOption":       "DISABLED",
		},
	}

	if cert.Options != nil && cert.Options.Export != "" {
		result["CertificateMetadata"].(map[string]interface{})["ExportOption"] = cert.Options.Export
	}

	if !cert.IssuedAt.IsZero() {
		result["CertificateMetadata"].(map[string]interface{})["IssuedAt"] = formatEpochSeconds(cert.IssuedAt)
	}
	if !cert.ImportedAt.IsZero() {
		result["CertificateMetadata"].(map[string]interface{})["ImportedAt"] = formatEpochSeconds(cert.ImportedAt)
	}

	// Parse X.509 attributes from the certificate PEM.
	x509Attrs := extractX509Attributes(cert.Certificate)
	if x509Attrs != nil {
		result["X509Attributes"] = x509Attrs
	}

	return result
}

// extractX509Attributes parses the PEM certificate and extracts X.509
// attributes for the search result.
func extractX509Attributes(certPEM string) map[string]interface{} {
	if certPEM == "" {
		return nil
	}
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}

	attrs := map[string]interface{}{
		"Subject":      buildDistinguishedName(parsed.Subject),
		"Issuer":       buildDistinguishedName(parsed.Issuer),
		"SerialNumber": parsed.SerialNumber.String(),
		"KeyAlgorithm": keyAlgorithmFromX509(parsed.PublicKey),
	}

	if !parsed.NotBefore.IsZero() {
		attrs["NotBefore"] = formatEpochSeconds(parsed.NotBefore)
	}
	if !parsed.NotAfter.IsZero() {
		attrs["NotAfter"] = formatEpochSeconds(parsed.NotAfter)
	}

	if len(parsed.DNSNames) > 0 {
		sans := make([]interface{}, len(parsed.DNSNames))
		for i, d := range parsed.DNSNames {
			sans[i] = map[string]interface{}{"DnsName": d}
		}
		attrs["SubjectAlternativeNames"] = sans
	}

	return attrs
}

// buildDistinguishedName converts a pkix.Name to the ACM DistinguishedName
// response structure.
func buildDistinguishedName(name pkixName) map[string]interface{} {
	dn := map[string]interface{}{}
	if name.CommonName != "" {
		dn["CommonName"] = name.CommonName
	}
	if len(name.Country) > 0 {
		dn["Country"] = name.Country[0]
	}
	if len(name.Locality) > 0 {
		dn["Locality"] = name.Locality[0]
	}
	if len(name.Organization) > 0 {
		dn["Organization"] = name.Organization[0]
	}
	if len(name.OrganizationalUnit) > 0 {
		dn["OrganizationalUnit"] = name.OrganizationalUnit[0]
	}
	if len(name.Province) > 0 {
		dn["State"] = name.Province[0]
	}
	if name.SerialNumber != "" {
		dn["SerialNumber"] = name.SerialNumber
	}
	return dn
}

// keyAlgorithmFromX509 derives the ACM KeyAlgorithm enum string from the
// actual public key by inspecting the key type and bit length / curve.
func keyAlgorithmFromX509(pub crypto.PublicKey) string {
	switch k := pub.(type) {
	case *rsa.PublicKey:
		return fmt.Sprintf("RSA_%d", k.N.BitLen())
	case *ecdsa.PublicKey:
		return ecCurveToAlgorithm(k.Curve)
	default:
		return "UNKNOWN"
	}
}

// ecCurveToAlgorithm maps a Go elliptic.Curve to the ACM KeyAlgorithm enum
// string used for EC certificates.
func ecCurveToAlgorithm(curve elliptic.Curve) string {
	switch curve {
	case elliptic.P256():
		return "EC_prime256v1"
	case elliptic.P384():
		return "EC_secp384r1"
	case elliptic.P521():
		return "EC_secp521r1"
	default:
		return "EC_unknown"
	}
}

// pkixName is an alias for crypto/x509/pkix.Name to simplify the
// buildDistinguishedName signature.
type pkixName = pkix.Name

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func parseIntToken(token string) (int, error) {
	return strconv.Atoi(token)
}

func formatIntToken(n int) string {
	return strconv.Itoa(n)
}
