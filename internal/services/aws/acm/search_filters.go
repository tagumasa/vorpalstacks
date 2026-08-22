package acm

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
)

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
	if certArn, ok := filterMap["CertificateArn"]; ok {
		arnStr, ok := certArn.(string)
		if !ok {
			return awserrors.NewValidationException("CertificateArn must be a string")
		}
		if !isValidCertificateArn(arnStr) {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid CertificateArn: %s", arnStr))
		}
	}
	if x509, ok := filterMap["X509AttributeFilter"]; ok {
		if err := validateX509Filter(x509); err != nil {
			return err
		}
	}
	if meta, ok := filterMap["AcmCertificateMetadataFilter"]; ok {
		if err := validateAcmCertificateMetadataFilter(meta); err != nil {
			return err
		}
	}
	return nil
}

// validateAcmCertificateMetadataFilter validates each enum-typed field in the
// AcmCertificateMetadataFilter against the corresponding Smithy enum.
func validateAcmCertificateMetadataFilter(filter interface{}) error {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return awserrors.NewValidationException("AcmCertificateMetadataFilter must be a structure")
	}
	if status, ok := filterMap["Status"]; ok {
		statusStr, ok := status.(string)
		if !ok || !validCertificateStatuses[statusStr] {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid Status filter value: %v", status))
		}
	}
	if renewalStatus, ok := filterMap["RenewalStatus"]; ok {
		rsStr, ok := renewalStatus.(string)
		if !ok {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid RenewalStatus filter value: %v", renewalStatus))
		}
		if !validRenewalStatuses[rsStr] {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid RenewalStatus filter value: %s", rsStr))
		}
	}
	if certType, ok := filterMap["Type"]; ok {
		typeStr, ok := certType.(string)
		if !ok {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid Type filter value: %v", certType))
		}
		if !validCertTypes[typeStr] {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid Type filter value: %s", typeStr))
		}
	}
	if _, ok := filterMap["InUse"]; ok {
		if _, ok := filterMap["InUse"].(bool); !ok {
			return awserrors.NewValidationException("InUse filter value must be a boolean")
		}
	}
	if _, ok := filterMap["Exported"]; ok {
		if _, ok := filterMap["Exported"].(bool); !ok {
			return awserrors.NewValidationException("Exported filter value must be a boolean")
		}
	}
	if exportOption, ok := filterMap["ExportOption"]; ok {
		eoStr, ok := exportOption.(string)
		if !ok || !validExportOptionValues[eoStr] {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid ExportOption filter value: %v", exportOption))
		}
	}
	if managedBy, ok := filterMap["ManagedBy"]; ok {
		mbStr, ok := managedBy.(string)
		if !ok || !validManagedByValues[mbStr] {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid ManagedBy filter value: %v", managedBy))
		}
	}
	if validationMethod, ok := filterMap["ValidationMethod"]; ok {
		vmStr, ok := validationMethod.(string)
		if !ok || !isValidValidationMethod(vmStr) {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid ValidationMethod filter value: %v", validationMethod))
		}
	}
	if ckpo, ok := filterMap["CertificateKeyPairOrigin"]; ok {
		ckpoStr, ok := ckpo.(string)
		if !ok || !validKeyPairOrigins[ckpoStr] {
			return awserrors.NewValidationException(fmt.Sprintf("Invalid CertificateKeyPairOrigin filter value: %v", ckpo))
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
	value, ok := m["Value"].(string)
	if !ok || value == "" {
		return awserrors.NewValidationException(fmt.Sprintf("%s filter Value is required", label))
	}
	// Smithy FilterString: @length(min:1, max:256), counted in Unicode
	// characters.
	if utf8.RuneCountInString(value) > maxFilterStringLength {
		return awserrors.NewValidationException(fmt.Sprintf("%s filter Value must not exceed %d characters", label, maxFilterStringLength))
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
		arnStr, ok := arnFilter.(string)
		if !ok {
			return false
		}
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
		algoStr, ok := keyAlgo.(string)
		if !ok {
			return false
		}
		return algoStr == cert.KeyAlgorithm
	}
	if keyUsage, ok := filterMap["KeyUsage"]; ok {
		kuStr, ok := keyUsage.(string)
		if !ok {
			return false
		}
		for _, ku := range cert.KeyUsages {
			if ku.Name == kuStr {
				return true
			}
		}
		return false
	}
	if extKeyUsage, ok := filterMap["ExtendedKeyUsage"]; ok {
		ekuStr, ok := extKeyUsage.(string)
		if !ok {
			return false
		}
		for _, eku := range cert.ExtendedKeyUsages {
			if eku.Name == ekuStr || eku.OID == ekuStr {
				return true
			}
		}
		return false
	}
	if serialNum, ok := filterMap["SerialNumber"]; ok {
		serialStr, ok := serialNum.(string)
		if !ok {
			return false
		}
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
	value, ok := cnFilterMap["Value"].(string)
	if !ok {
		return false
	}
	// Smithy FilterString: @length(1-256).
	if value == "" || !isValidFilterString(value) {
		return false
	}
	op, ok := cnFilterMap["ComparisonOperator"].(string)
	if !ok {
		return false
	}
	cn := strings.ToLower(extractCommonName(cert.Subject))
	searchVal := strings.ToLower(value)
	switch op {
	case "EQUALS":
		return cn == searchVal
	case "CONTAINS":
		return strings.Contains(cn, searchVal)
	default:
		return false
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
		value, ok := dnFilterMap["Value"].(string)
		if !ok {
			return false
		}
		// Smithy FilterString: @length(1-256).
		if value == "" || !isValidFilterString(value) {
			return false
		}
		op, ok := dnFilterMap["ComparisonOperator"].(string)
		if !ok {
			return false
		}
		searchVal := strings.ToLower(value)
		for _, san := range cert.SubjectAlternativeNames {
			sanLower := strings.ToLower(san)
			switch op {
			case "EQUALS":
				if sanLower == searchVal {
					return true
				}
			case "CONTAINS":
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
		statusStr, ok := status.(string)
		if !ok {
			return false
		}
		return statusStr == cert.Status
	}
	if renewalStatus, ok := filterMap["RenewalStatus"]; ok {
		rsStr, ok := renewalStatus.(string)
		if !ok {
			return false
		}
		if cert.RenewalSummary == nil {
			return false
		}
		return rsStr == cert.RenewalSummary.RenewalStatus
	}
	if certType, ok := filterMap["Type"]; ok {
		typeStr, ok := certType.(string)
		if !ok {
			return false
		}
		return typeStr == cert.Type
	}
	if inUse, ok := filterMap["InUse"]; ok {
		inUseVal, ok := inUse.(bool)
		if !ok {
			return false
		}
		return inUseVal == (len(cert.InUseBy) > 0)
	}
	if exported, ok := filterMap["Exported"]; ok {
		exportedVal, ok := exported.(bool)
		if !ok {
			return false
		}
		return exportedVal == cert.WasExported
	}
	if exportOption, ok := filterMap["ExportOption"]; ok {
		eoStr, ok := exportOption.(string)
		if !ok {
			return false
		}
		if cert.Options == nil {
			return eoStr == "DISABLED"
		}
		return eoStr == cert.Options.Export
	}
	if managedBy, ok := filterMap["ManagedBy"]; ok {
		mbStr, ok := managedBy.(string)
		if !ok {
			return false
		}
		return mbStr == cert.ManagedBy
	}
	if validationMethod, ok := filterMap["ValidationMethod"]; ok {
		vmStr, ok := validationMethod.(string)
		if !ok {
			return false
		}
		for _, dvo := range cert.DomainValidationOptions {
			if dvo.ValidationMethod == vmStr {
				return true
			}
		}
		return false
	}
	if certKeyPairOrigin, ok := filterMap["CertificateKeyPairOrigin"]; ok {
		ckpoStr, ok := certKeyPairOrigin.(string)
		if !ok {
			return false
		}
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
