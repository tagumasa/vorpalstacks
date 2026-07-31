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
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

// SearchCertificates retrieves a list of certificates matching search criteria.
func (s *ACMService) SearchCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters

	nextToken := request.GetStringParam(params, "NextToken")
	maxResults := 500
	if mrRaw, ok := params["MaxResults"]; ok {
		mr := request.GetIntParam(params, "MaxResults")
		if mr < 1 {
			return nil, awserrors.NewValidationException(fmt.Sprintf("MaxResults must be between 1 and 500, got %v", mrRaw))
		}
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
	if err := validateSearchSortBy(sortBy); err != nil {
		return nil, err
	}
	if sortOrder != "" {
		if err := validateSortOrder(sortOrder); err != nil {
			return nil, err
		}
	}
	sortCertificateSearchResults(filtered, sortBy, sortOrder)

	// Paginate using offset-based tokens.
	offset := 0
	if nextToken != "" {
		n, parseErr := parseIntToken(nextToken)
		if parseErr != nil {
			return nil, NewInvalidParameterError(fmt.Sprintf("Invalid NextToken: %s", nextToken))
		}
		if n < 0 {
			return nil, NewInvalidParameterError("Invalid NextToken: negative offset")
		}
		offset = n
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
// parameters. Supports the full CertificateFilterStatement union type
// including And/Or/Not compound operators and nested filter unions.
func applyCertificateFilters(certs []*acmstorelib.Certificate, params map[string]interface{}) []*acmstorelib.Certificate {
	filterStatement := getNestedValue(params, "FilterStatement")
	if filterStatement == nil {
		return certs
	}
	var result []*acmstorelib.Certificate
	for _, c := range certs {
		if evaluateFilterStatement(c, filterStatement) {
			result = append(result, c)
		}
	}
	return result
}

// getNestedValue extracts a nested value from a map by walking the key path.
func getNestedValue(m map[string]interface{}, keys ...string) interface{} {
	var current interface{} = m
	for _, key := range keys {
		cm, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = cm[key]
		if !ok {
			return nil
		}
	}
	return current
}

// evaluateFilterStatement recursively evaluates a CertificateFilterStatement
// against a certificate. The statement is a union type with And/Or/Not/Filter.
func evaluateFilterStatement(cert *acmstorelib.Certificate, statement interface{}) bool {
	if statement == nil {
		return false
	}
	stmtMap, ok := statement.(map[string]interface{})
	if !ok {
		return false
	}
	if andList, ok := stmtMap["And"]; ok {
		return evaluateAndFilter(cert, andList)
	}
	if orList, ok := stmtMap["Or"]; ok {
		return evaluateOrFilter(cert, orList)
	}
	if notStmt, ok := stmtMap["Not"]; ok {
		return !evaluateFilterStatement(cert, notStmt)
	}
	if filter, ok := stmtMap["Filter"]; ok {
		return evaluateCertificateFilter(cert, filter)
	}
	return false
}

// evaluateAndFilter evaluates an And compound filter (all must match).
func evaluateAndFilter(cert *acmstorelib.Certificate, list interface{}) bool {
	arr, ok := list.([]interface{})
	if !ok || len(arr) == 0 {
		return false
	}
	for _, item := range arr {
		if !evaluateFilterStatement(cert, item) {
			return false
		}
	}
	return true
}

// evaluateOrFilter evaluates an Or compound filter (any must match).
func evaluateOrFilter(cert *acmstorelib.Certificate, list interface{}) bool {
	arr, ok := list.([]interface{})
	if !ok || len(arr) == 0 {
		return false
	}
	for _, item := range arr {
		if evaluateFilterStatement(cert, item) {
			return true
		}
	}
	return false
}

// evaluateCertificateFilter evaluates a CertificateFilter union against a cert.
func evaluateCertificateFilter(cert *acmstorelib.Certificate, filter interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	if arnFilter, ok := filterMap["CertificateArn"]; ok {
		arnStr, _ := arnFilter.(string)
		return arnStr == cert.CertificateArn
	}
	if x509Filter, ok := filterMap["X509AttributeFilter"]; ok {
		return evaluateX509Filter(cert, x509Filter)
	}
	if metaFilter, ok := filterMap["AcmCertificateMetadataFilter"]; ok {
		return evaluateMetadataFilter(cert, metaFilter)
	}
	return false
}

// evaluateX509Filter evaluates an X509AttributeFilter union against a cert.
func evaluateX509Filter(cert *acmstorelib.Certificate, filter interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	if subjectFilter, ok := filterMap["Subject"]; ok {
		return evaluateSubjectFilter(cert, subjectFilter)
	}
	if keyAlgo, ok := filterMap["KeyAlgorithm"]; ok {
		algoStr, _ := keyAlgo.(string)
		return algoStr == cert.KeyAlgorithm
	}
	if serialNum, ok := filterMap["SerialNumber"]; ok {
		serialStr, _ := serialNum.(string)
		return serialStr == cert.Serial
	}
	if notAfter, ok := filterMap["NotAfter"]; ok {
		return evaluateTimestampRange(cert.NotAfter, notAfter)
	}
	if notBefore, ok := filterMap["NotBefore"]; ok {
		return evaluateTimestampRange(cert.NotBefore, notBefore)
	}
	return false
}

// evaluateSubjectFilter evaluates a SubjectFilter union against a cert.
func evaluateSubjectFilter(cert *acmstorelib.Certificate, filter interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	if commonName, ok := filterMap["CommonName"]; ok {
		cnStr, _ := commonName.(string)
		return strings.Contains(strings.ToLower(cert.Subject), strings.ToLower(cnStr))
	}
	return false
}

// evaluateMetadataFilter evaluates an AcmCertificateMetadataFilter union against a cert.
func evaluateMetadataFilter(cert *acmstorelib.Certificate, filter interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	if status, ok := filterMap["Status"]; ok {
		statusStr, _ := status.(string)
		return statusStr == cert.Status
	}
	if renewalStatus, ok := filterMap["RenewalStatus"]; ok {
		rsStr, _ := renewalStatus.(string)
		if cert.RenewalSummary == nil {
			return rsStr == ""
		}
		return rsStr == cert.RenewalSummary.RenewalStatus
	}
	if certType, ok := filterMap["Type"]; ok {
		typeStr, _ := certType.(string)
		return typeStr == cert.Type
	}
	if inUse, ok := filterMap["InUse"]; ok {
		inUseVal, _ := inUse.(bool)
		return inUseVal == (len(cert.InUseBy) > 0)
	}
	if exported, ok := filterMap["Exported"]; ok {
		exportedVal, _ := exported.(bool)
		isExported := cert.Options != nil && cert.Options.Export == "ENABLED"
		return exportedVal == isExported
	}
	if exportOption, ok := filterMap["ExportOption"]; ok {
		eoStr, _ := exportOption.(string)
		if cert.Options == nil {
			return eoStr == "DISABLED"
		}
		return eoStr == cert.Options.Export
	}
	if managedBy, ok := filterMap["ManagedBy"]; ok {
		mbStr, _ := managedBy.(string)
		return mbStr == cert.ManagedBy
	}
	if validationMethod, ok := filterMap["ValidationMethod"]; ok {
		vmStr, _ := validationMethod.(string)
		for _, dvo := range cert.DomainValidationOptions {
			if dvo.ValidationMethod == vmStr {
				return true
			}
		}
		return false
	}
	if certKeyPairOrigin, ok := filterMap["CertificateKeyPairOrigin"]; ok {
		ckpoStr, _ := certKeyPairOrigin.(string)
		return ckpoStr == cert.CertificateKeyPairOrigin
	}
	return false
}

// evaluateTimestampRange checks if a time falls within a TimestampRange.
func evaluateTimestampRange(t time.Time, rangeVal interface{}) bool {
	rangeMap, ok := rangeVal.(map[string]interface{})
	if !ok {
		return false
	}
	start, startOk := rangeMap["Start"].(float64)
	end, endOk := rangeMap["End"].(float64)
	if !startOk && !endOk {
		return false
	}
	if startOk {
		startTime := time.Unix(int64(start), 0)
		if t.Before(startTime) {
			return false
		}
	}
	if endOk {
		endTime := time.Unix(int64(end), 0)
		if t.After(endTime) {
			return false
		}
	}
	return true
}

// validSearchSortByValues contains all supported SortBy values for
// SearchCertificates.  ACME scope-out values (ACME_ENDPOINT_ARN,
// ACME_ACCOUNT_ID) are excluded as the ACME API is not implemented.
var validSearchSortByValues = map[string]bool{
	"CREATED_AT":                  true,
	"NOT_AFTER":                   true,
	"STATUS":                      true,
	"RENEWAL_STATUS":              true,
	"EXPORTED":                    true,
	"IN_USE":                      true,
	"NOT_BEFORE":                  true,
	"KEY_ALGORITHM":               true,
	"TYPE":                        true,
	"CERTIFICATE_ARN":             true,
	"COMMON_NAME":                 true,
	"REVOKED_AT":                  true,
	"RENEWAL_ELIGIBILITY":         true,
	"ISSUED_AT":                   true,
	"MANAGED_BY":                  true,
	"EXPORT_OPTION":               true,
	"VALIDATION_METHOD":           true,
	"IMPORTED_AT":                 true,
	"CERTIFICATE_KEY_PAIR_ORIGIN": true,
}

// validateSearchSortBy returns an error if sortBy is not a supported value.
func validateSearchSortBy(sortBy string) error {
	if sortBy == "" || sortBy == "CREATED_AT" {
		return nil
	}
	if !validSearchSortByValues[sortBy] {
		return NewInvalidParameterError(fmt.Sprintf("Invalid SortBy: %s", sortBy))
	}
	return nil
}

// validateSortOrder validates SortOrder parameter (ASCENDING/DESCENDING).
func validateSortOrder(sortOrder string) error {
	switch sortOrder {
	case "ASCENDING", "DESCENDING":
		return nil
	default:
		return NewInvalidParameterError(fmt.Sprintf("Invalid SortOrder: %s. Valid values are ASCENDING, DESCENDING.", sortOrder))
	}
}

// sortCertificateSearchResults sorts the slice by the given field and order.
func sortCertificateSearchResults(certs []*acmstorelib.Certificate, sortBy, sortOrder string) {
	ascending := sortOrder != "DESCENDING"
	less := func(i, j int) bool {
		a, b := certs[i], certs[j]
		switch sortBy {
		case "NOT_AFTER":
			if ascending {
				return a.NotAfter.Before(b.NotAfter)
			}
			return a.NotAfter.After(b.NotAfter)
		case "STATUS":
			if ascending {
				return a.Status < b.Status
			}
			return a.Status > b.Status
		case "RENEWAL_STATUS":
			aVal, bVal := renewalStatus(a), renewalStatus(b)
			if ascending {
				return aVal < bVal
			}
			return aVal > bVal
		case "EXPORTED":
			aVal, bVal := isExported(a), isExported(b)
			if ascending {
				return aVal < bVal
			}
			return aVal > bVal
		case "IN_USE":
			aVal, bVal := isInUse(a), isInUse(b)
			if ascending {
				return aVal < bVal
			}
			return aVal > bVal
		case "NOT_BEFORE":
			if ascending {
				return a.NotBefore.Before(b.NotBefore)
			}
			return a.NotBefore.After(b.NotBefore)
		case "KEY_ALGORITHM":
			if ascending {
				return a.KeyAlgorithm < b.KeyAlgorithm
			}
			return a.KeyAlgorithm > b.KeyAlgorithm
		case "TYPE":
			if ascending {
				return a.Type < b.Type
			}
			return a.Type > b.Type
		case "CERTIFICATE_ARN":
			if ascending {
				return a.CertificateArn < b.CertificateArn
			}
			return a.CertificateArn > b.CertificateArn
		case "COMMON_NAME":
			if ascending {
				return a.Subject < b.Subject
			}
			return a.Subject > b.Subject
		case "REVOKED_AT":
			if ascending {
				return a.RevokedAt.Before(b.RevokedAt)
			}
			return a.RevokedAt.After(b.RevokedAt)
		case "RENEWAL_ELIGIBILITY":
			if ascending {
				return a.RenewalEligibility < b.RenewalEligibility
			}
			return a.RenewalEligibility > b.RenewalEligibility
		case "ISSUED_AT":
			if ascending {
				return a.IssuedAt.Before(b.IssuedAt)
			}
			return a.IssuedAt.After(b.IssuedAt)
		case "MANAGED_BY":
			if ascending {
				return a.ManagedBy < b.ManagedBy
			}
			return a.ManagedBy > b.ManagedBy
		case "EXPORT_OPTION":
			aVal, bVal := exportOption(a), exportOption(b)
			if ascending {
				return aVal < bVal
			}
			return aVal > bVal
		case "VALIDATION_METHOD":
			aVal, bVal := validationMethod(a), validationMethod(b)
			if ascending {
				return aVal < bVal
			}
			return aVal > bVal
		case "IMPORTED_AT":
			if ascending {
				return a.ImportedAt.Before(b.ImportedAt)
			}
			return a.ImportedAt.After(b.ImportedAt)
		case "CERTIFICATE_KEY_PAIR_ORIGIN":
			if ascending {
				return a.CertificateKeyPairOrigin < b.CertificateKeyPairOrigin
			}
			return a.CertificateKeyPairOrigin > b.CertificateKeyPairOrigin
		default:
			if ascending {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			return a.CreatedAt.After(b.CreatedAt)
		}
	}
	sort.Slice(certs, less)
}

// renewalStatus returns the renewal status string, or empty if not set.
func renewalStatus(cert *acmstorelib.Certificate) string {
	if cert.RenewalSummary != nil {
		return cert.RenewalSummary.RenewalStatus
	}
	return ""
}

// isExported returns 1 if the certificate is exported, 0 otherwise.
func isExported(cert *acmstorelib.Certificate) int {
	if cert.Options != nil && cert.Options.Export == "ENABLED" {
		return 1
	}
	return 0
}

// isInUse returns 1 if the certificate is in use, 0 otherwise.
func isInUse(cert *acmstorelib.Certificate) int {
	if len(cert.InUseBy) > 0 {
		return 1
	}
	return 0
}

// exportOption returns the export option string, or empty if not set.
func exportOption(cert *acmstorelib.Certificate) string {
	if cert.Options != nil {
		return cert.Options.Export
	}
	return ""
}

// validationMethod returns the first validation method from domain
// validation options, or empty if not set.
func validationMethod(cert *acmstorelib.Certificate) string {
	if len(cert.DomainValidationOptions) > 0 {
		return cert.DomainValidationOptions[0].ValidationMethod
	}
	return ""
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
	x509Attrs, _ := extractX509Attributes(cert.Certificate)
	if x509Attrs != nil {
		result["X509Attributes"] = x509Attrs
	}

	return result
}

// extractX509Attributes parses the PEM certificate and extracts X.509
// attributes for the search result.
func extractX509Attributes(certPEM string) (map[string]interface{}, error) {
	if certPEM == "" {
		return nil, nil
	}
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM certificate")
	}
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X.509 certificate: %w", err)
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

	return attrs, nil
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

func parseIntToken(token string) (int, error) {
	return strconv.Atoi(token)
}

func formatIntToken(n int) string {
	return strconv.Itoa(n)
}
