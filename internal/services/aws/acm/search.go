package acm

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
	"vorpalstacks/internal/utils/timeutils"
)

// SearchCertificates retrieves a list of certificates matching search criteria.
func (s *ACMService) SearchCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters

	nextToken := request.GetStringParam(params, "NextToken")
	maxResults := 100
	if _, ok := params["MaxResults"]; ok {
		mr := request.GetIntParam(params, "MaxResults")
		if mr < 1 || mr > 500 {
			return nil, awserrors.NewValidationException(fmt.Sprintf("MaxResults must be between 1 and 500, got %d", mr))
		}
		maxResults = mr
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
	filtered, err := applyCertificateFilters(allCerts, params)
	if err != nil {
		return nil, err
	}

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
// Returns an error if the filter structure violates Smithy constraints
// (e.g. list length outside 1-15, missing REQUIRED fields on filter members).
func applyCertificateFilters(certs []*acmstorelib.Certificate, params map[string]interface{}) ([]*acmstorelib.Certificate, error) {
	filterStatement := getNestedValue(params, "FilterStatement")
	if filterStatement == nil {
		return certs, nil
	}
	if err := validateFilterStatement(filterStatement); err != nil {
		return nil, err
	}
	var result []*acmstorelib.Certificate
	for _, c := range certs {
		if evaluateFilterStatement(c, filterStatement) {
			result = append(result, c)
		}
	}
	return result, nil
}

// validateFilterStatement recursively walks a CertificateFilterStatement
// union and validates Smithy constraints:
//   - CertificateFilterStatementList (And/Or): length 1-15
//   - CommonNameFilter: Value + ComparisonOperator both REQUIRED
//   - DnsNameFilter: Value + ComparisonOperator both REQUIRED
func validateFilterStatement(statement interface{}) error {
	if statement == nil {
		return nil
	}
	stmtMap, ok := statement.(map[string]interface{})
	if !ok {
		return awserrors.NewValidationException("FilterStatement must be a structure")
	}
	if andList, ok := stmtMap["And"]; ok {
		if err := validateFilterList(andList, "And"); err != nil {
			return err
		}
	}
	if orList, ok := stmtMap["Or"]; ok {
		if err := validateFilterList(orList, "Or"); err != nil {
			return err
		}
	}
	if notStmt, ok := stmtMap["Not"]; ok {
		if err := validateFilterStatement(notStmt); err != nil {
			return err
		}
	}
	if filter, ok := stmtMap["Filter"]; ok {
		if err := validateCertificateFilter(filter); err != nil {
			return err
		}
	}
	return nil
}

// validateFilterList validates that a CertificateFilterStatementList
// has between 1 and 15 items, then recursively validates each statement.
func validateFilterList(list interface{}, label string) error {
	arr, ok := list.([]interface{})
	if !ok {
		return awserrors.NewValidationException(fmt.Sprintf("%s must be a list of filter statements", label))
	}
	if len(arr) < 1 || len(arr) > 15 {
		return awserrors.NewValidationException(fmt.Sprintf("%s list must contain between 1 and 15 items, got %d", label, len(arr)))
	}
	for _, item := range arr {
		if err := validateFilterStatement(item); err != nil {
			return err
		}
	}
	return nil
}

// validateCertificateFilter validates a CertificateFilter union, which
// contains nested filter unions (CertificateArn, X509AttributeFilter,
// AcmCertificateMetadataFilter).
func validateCertificateFilter(filter interface{}) error {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return awserrors.NewValidationException("Filter must be a structure")
	}
	if x509, ok := filterMap["X509AttributeFilter"]; ok {
		if err := validateX509Filter(x509); err != nil {
			return err
		}
	}
	return nil
}

// validateX509Filter validates an X509AttributeFilter union. The Subject
// member contains a CommonNameFilter (Value + ComparisonOperator, both
// REQUIRED) and the SubjectAlternativeName member contains a DnsNameFilter
// (Value + ComparisonOperator, both REQUIRED).
func validateX509Filter(filter interface{}) error {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return awserrors.NewValidationException("X509AttributeFilter must be a structure")
	}
	if subjectRaw, ok := filterMap["Subject"]; ok {
		subjectMap, ok := subjectRaw.(map[string]interface{})
		if !ok {
			return awserrors.NewValidationException("Subject filter must be a structure")
		}
		if cnRaw, ok := subjectMap["CommonName"]; ok {
			cnMap, ok := cnRaw.(map[string]interface{})
			if !ok {
				return awserrors.NewValidationException("CommonName filter must be a structure")
			}
			if err := validateRequiredComparisonFilter(cnMap, "CommonName"); err != nil {
				return err
			}
		}
	}
	if sanRaw, ok := filterMap["SubjectAlternativeName"]; ok {
		sanMap, ok := sanRaw.(map[string]interface{})
		if !ok {
			return awserrors.NewValidationException("SubjectAlternativeName filter must be a structure")
		}
		if dnRaw, ok := sanMap["DnsName"]; ok {
			dnMap, ok := dnRaw.(map[string]interface{})
			if !ok {
				return awserrors.NewValidationException("DnsName filter must be a structure")
			}
			if err := validateRequiredComparisonFilter(dnMap, "DnsName"); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateRequiredComparisonFilter validates that a filter structure (e.g.
// CommonNameFilter, DnsNameFilter) contains both Value (REQUIRED) and
// ComparisonOperator (REQUIRED) with valid values.
func validateRequiredComparisonFilter(m map[string]interface{}, label string) error {
	value, _ := m["Value"].(string)
	if value == "" {
		return awserrors.NewValidationException(fmt.Sprintf("%s filter Value is required", label))
	}
	op, ok := m["ComparisonOperator"].(string)
	if !ok || op == "" {
		return awserrors.NewValidationException(fmt.Sprintf("%s filter ComparisonOperator is required", label))
	}
	if op != "EQUALS" && op != "CONTAINS" {
		return awserrors.NewValidationException(fmt.Sprintf("%s filter ComparisonOperator must be EQUALS or CONTAINS, got %s", label, op))
	}
	return nil
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
	if sanFilter, ok := filterMap["SubjectAlternativeName"]; ok {
		return evaluateSubjectAlternativeNameFilter(cert, sanFilter)
	}
	if keyAlgo, ok := filterMap["KeyAlgorithm"]; ok {
		algoStr, _ := keyAlgo.(string)
		return algoStr == cert.KeyAlgorithm
	}
	if keyUsage, ok := filterMap["KeyUsage"]; ok {
		kuStr, _ := keyUsage.(string)
		for _, ku := range cert.KeyUsages {
			if ku.Name == kuStr {
				return true
			}
		}
		return false
	}
	if extKeyUsage, ok := filterMap["ExtendedKeyUsage"]; ok {
		ekuStr, _ := extKeyUsage.(string)
		for _, eku := range cert.ExtendedKeyUsages {
			if eku.Name == ekuStr || eku.OID == ekuStr {
				return true
			}
		}
		return false
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

// extractCommonName extracts the CN value from an RFC 2253 DN string
// produced by pkix.Name.String() (e.g. "CN=example.com,O=Org" → "example.com").
func extractCommonName(subject string) string {
	for _, part := range strings.Split(subject, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "CN=") {
			return part[3:]
		}
	}
	return ""
}

// evaluateSubjectFilter evaluates a SubjectFilter union against a cert.
func evaluateSubjectFilter(cert *acmstorelib.Certificate, filter interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	commonNameFilter, ok := filterMap["CommonName"]
	if !ok {
		return false
	}
	// CommonNameFilter is a structure with Value (required) and
	// ComparisonOperator (required, Smithy enum: EQUALS | CONTAINS).
	cnFilterMap, ok := commonNameFilter.(map[string]interface{})
	if !ok {
		return false
	}
	value, _ := cnFilterMap["Value"].(string)
	if value == "" {
		return false
	}
	op, _ := cnFilterMap["ComparisonOperator"].(string)
	cn := strings.ToLower(extractCommonName(cert.Subject))
	searchVal := strings.ToLower(value)
	switch op {
	case "EQUALS":
		return cn == searchVal
	default:
		// CONTAINS is the AWS default when the field is absent.
		return strings.Contains(cn, searchVal)
	}
}

// evaluateSubjectAlternativeNameFilter evaluates a SubjectAlternativeNameFilter
// union against a cert. The only supported member is DnsName, which is a
// DnsNameFilter structure containing Value (FilterString, REQUIRED) and
// ComparisonOperator (REQUIRED: EQUALS | CONTAINS).
func evaluateSubjectAlternativeNameFilter(cert *acmstorelib.Certificate, filter interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	if dnsNameRaw, ok := filterMap["DnsName"]; ok {
		// DnsName targets a DnsNameFilter structure, not a plain string.
		dnFilterMap, ok := dnsNameRaw.(map[string]interface{})
		if !ok {
			return false
		}
		value, _ := dnFilterMap["Value"].(string)
		if value == "" {
			return false
		}
		op, _ := dnFilterMap["ComparisonOperator"].(string)
		searchVal := strings.ToLower(value)
		for _, san := range cert.SubjectAlternativeNames {
			sanLower := strings.ToLower(san)
			if op == "EQUALS" {
				if sanLower == searchVal {
					return true
				}
			} else {
				// CONTAINS is the default.
				if strings.Contains(sanLower, searchVal) {
					return true
				}
			}
		}
		return false
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
// cmpTime compares two timestamps for sorting. Returns -1, 0, or 1.
func cmpTime(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}

// sortKeyFuncs maps each SortBy value to a comparator returning
// negative/zero/positive for less/equal/greater.
var sortKeyFuncs = map[string]func(a, b *acmstorelib.Certificate) int{
	"CREATED_AT":      func(a, b *acmstorelib.Certificate) int { return cmpTime(a.CreatedAt, b.CreatedAt) },
	"NOT_BEFORE":      func(a, b *acmstorelib.Certificate) int { return cmpTime(a.NotBefore, b.NotBefore) },
	"NOT_AFTER":       func(a, b *acmstorelib.Certificate) int { return cmpTime(a.NotAfter, b.NotAfter) },
	"ISSUED_AT":       func(a, b *acmstorelib.Certificate) int { return cmpTime(a.IssuedAt, b.IssuedAt) },
	"IMPORTED_AT":     func(a, b *acmstorelib.Certificate) int { return cmpTime(a.ImportedAt, b.ImportedAt) },
	"REVOKED_AT":      func(a, b *acmstorelib.Certificate) int { return cmpTime(a.RevokedAt, b.RevokedAt) },
	"STATUS":          func(a, b *acmstorelib.Certificate) int { return strings.Compare(a.Status, b.Status) },
	"TYPE":            func(a, b *acmstorelib.Certificate) int { return strings.Compare(a.Type, b.Type) },
	"KEY_ALGORITHM":   func(a, b *acmstorelib.Certificate) int { return strings.Compare(a.KeyAlgorithm, b.KeyAlgorithm) },
	"CERTIFICATE_ARN": func(a, b *acmstorelib.Certificate) int { return strings.Compare(a.CertificateArn, b.CertificateArn) },
	"COMMON_NAME":     func(a, b *acmstorelib.Certificate) int { return strings.Compare(a.Subject, b.Subject) },
	"RENEWAL_ELIGIBILITY": func(a, b *acmstorelib.Certificate) int {
		return strings.Compare(a.RenewalEligibility, b.RenewalEligibility)
	},
	"MANAGED_BY": func(a, b *acmstorelib.Certificate) int { return strings.Compare(a.ManagedBy, b.ManagedBy) },
	"CERTIFICATE_KEY_PAIR_ORIGIN": func(a, b *acmstorelib.Certificate) int {
		return strings.Compare(a.CertificateKeyPairOrigin, b.CertificateKeyPairOrigin)
	},
	"RENEWAL_STATUS": func(a, b *acmstorelib.Certificate) int {
		return strings.Compare(renewalStatus(a), renewalStatus(b))
	},
	"EXPORTED": func(a, b *acmstorelib.Certificate) int {
		return isExported(a) - isExported(b)
	},
	"IN_USE": func(a, b *acmstorelib.Certificate) int {
		return isInUse(a) - isInUse(b)
	},
	"EXPORT_OPTION": func(a, b *acmstorelib.Certificate) int {
		return strings.Compare(exportOption(a), exportOption(b))
	},
	"VALIDATION_METHOD": func(a, b *acmstorelib.Certificate) int {
		return strings.Compare(validationMethod(a), validationMethod(b))
	},
}

func sortCertificateSearchResults(certs []*acmstorelib.Certificate, sortBy, sortOrder string) {
	keyFunc, ok := sortKeyFuncs[sortBy]
	if !ok {
		keyFunc = sortKeyFuncs["CREATED_AT"]
	}
	ascending := sortOrder != "DESCENDING"
	sort.Slice(certs, func(i, j int) bool {
		c := keyFunc(certs[i], certs[j])
		if ascending {
			return c < 0
		}
		return c > 0
	})
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
// Per the Smithy model, CertificateMetadata is a union whose sole member
// is AcmCertificateMetadata — the metadata fields must be nested inside it.
// Only fields declared in the AcmCertificateMetadata shape are emitted;
// X.509 attributes (KeyAlgorithm, NotBefore, NotAfter) belong in
// X509Attributes, which is populated separately by extractX509Attributes.
func buildCertificateSearchResult(cert *acmstorelib.Certificate) map[string]interface{} {
	metadata := map[string]interface{}{
		"Status":                   cert.Status,
		"Type":                     cert.Type,
		"RenewalEligibility":       cert.RenewalEligibility,
		"CreatedAt":                timeutils.FormatEpochSeconds(cert.CreatedAt),
		"Exported":                 cert.WasExported,
		"InUse":                    len(cert.InUseBy) > 0,
		"CertificateKeyPairOrigin": cert.CertificateKeyPairOrigin,
		"ManagedBy":                cert.ManagedBy,
		"ExportOption":             "DISABLED",
	}

	if cert.Options != nil && cert.Options.Export != "" {
		metadata["ExportOption"] = cert.Options.Export
	}

	if !cert.IssuedAt.IsZero() {
		metadata["IssuedAt"] = timeutils.FormatEpochSeconds(cert.IssuedAt)
	}
	if !cert.ImportedAt.IsZero() {
		metadata["ImportedAt"] = timeutils.FormatEpochSeconds(cert.ImportedAt)
	}
	if !cert.RevokedAt.IsZero() {
		metadata["RevokedAt"] = timeutils.FormatEpochSeconds(cert.RevokedAt)
	}
	if cert.RenewalSummary != nil {
		metadata["RenewalStatus"] = cert.RenewalSummary.RenewalStatus
	}
	if len(cert.DomainValidationOptions) > 0 {
		metadata["ValidationMethod"] = cert.DomainValidationOptions[0].ValidationMethod
	}

	result := map[string]interface{}{
		"CertificateArn": cert.CertificateArn,
		"CertificateMetadata": map[string]interface{}{
			"AcmCertificateMetadata": metadata,
		},
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
		attrs["NotBefore"] = timeutils.FormatEpochSeconds(parsed.NotBefore)
	}
	if !parsed.NotAfter.IsZero() {
		attrs["NotAfter"] = timeutils.FormatEpochSeconds(parsed.NotAfter)
	}

	if len(parsed.DNSNames) > 0 {
		sans := make([]interface{}, len(parsed.DNSNames))
		for i, d := range parsed.DNSNames {
			sans[i] = map[string]interface{}{"DnsName": d}
		}
		attrs["SubjectAlternativeNames"] = sans
	}

	// Extract KeyUsage X.509 v3 extension flags as ACM KeyUsageName strings.
	keyUsages := keyUsageFlagsToStrings(parsed.KeyUsage)
	if len(keyUsages) > 0 {
		attrs["KeyUsages"] = keyUsages
	}

	// Extract ExtendedKeyUsage OIDs as ACM ExtendedKeyUsageName strings.
	extKeyUsages := extKeyUsageEnumsToStrings(parsed.ExtKeyUsage, parsed.UnknownExtKeyUsage)
	if len(extKeyUsages) > 0 {
		attrs["ExtendedKeyUsages"] = extKeyUsages
	}

	return attrs, nil
}

// keyUsageFlagsToStrings converts Go x509.KeyUsage bitmask flags to ACM
// KeyUsageName enum strings (per the Smithy KeyUsageName shape).
func keyUsageFlagsToStrings(ku x509.KeyUsage) []string {
	type flagPair struct {
		flag x509.KeyUsage
		name string
	}
	pairs := []flagPair{
		{x509.KeyUsageDigitalSignature, "DIGITAL_SIGNATURE"},
		{x509.KeyUsageContentCommitment, "NON_REPUDIATION"},
		{x509.KeyUsageKeyEncipherment, "KEY_ENCIPHERMENT"},
		{x509.KeyUsageDataEncipherment, "DATA_ENCIPHERMENT"},
		{x509.KeyUsageKeyAgreement, "KEY_AGREEMENT"},
		{x509.KeyUsageCertSign, "CERTIFICATE_SIGNING"},
		{x509.KeyUsageCRLSign, "CRL_SIGNING"},
		{x509.KeyUsageEncipherOnly, "ENCIPHER_ONLY"},
		{x509.KeyUsageDecipherOnly, "DECIPHER_ONLY"},
	}
	var result []string
	for _, p := range pairs {
		if ku&p.flag != 0 {
			result = append(result, p.name)
		}
	}
	return result
}

// extKeyUsageEnumsToStrings converts Go x509.ExtKeyUsage enum values and
// any unknown OIDs to ACM ExtendedKeyUsageName strings.
func extKeyUsageEnumsToStrings(ekus []x509.ExtKeyUsage, unknownOIDs []asn1.ObjectIdentifier) []string {
	var result []string
	for _, eku := range ekus {
		name, ok := extKeyUsageMap[eku]
		if ok {
			result = append(result, name)
		}
	}
	for _, oid := range unknownOIDs {
		result = append(result, oid.String())
	}
	return result
}

// extKeyUsageMap maps Go x509.ExtKeyUsage values to ACM ExtendedKeyUsageName
// enum strings (per the Smithy ExtendedKeyUsageName shape). Only values that
// have a direct 1-to-1 correspondence are included; SGC and other Go EKU
// constants without an ACM equivalent are omitted (the if-ok check in
// extKeyUsageEnumsToStrings skips them, and unknown OIDs surface via
// UnknownExtKeyUsage).
var extKeyUsageMap = map[x509.ExtKeyUsage]string{
	x509.ExtKeyUsageServerAuth:      "TLS_WEB_SERVER_AUTHENTICATION",
	x509.ExtKeyUsageClientAuth:      "TLS_WEB_CLIENT_AUTHENTICATION",
	x509.ExtKeyUsageCodeSigning:     "CODE_SIGNING",
	x509.ExtKeyUsageEmailProtection: "EMAIL_PROTECTION",
	x509.ExtKeyUsageIPSECEndSystem:  "IPSEC_END_SYSTEM",
	x509.ExtKeyUsageIPSECTunnel:     "IPSEC_TUNNEL",
	x509.ExtKeyUsageIPSECUser:       "IPSEC_USER",
	x509.ExtKeyUsageTimeStamping:    "TIME_STAMPING",
	x509.ExtKeyUsageOCSPSigning:     "OCSP_SIGNING",
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
