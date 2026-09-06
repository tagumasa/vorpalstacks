package acm

import (
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
	acmstorelib "vorpalstacks/internal/store/aws/acm"
	"vorpalstacks/internal/utils/timeutils"
)

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
// SearchCertificates declares ValidationException (not
// InvalidParameterException), so violations surface as ValidationException.
func validateSearchSortBy(sortBy string) error {
	if sortBy == "" || sortBy == "CREATED_AT" {
		return nil
	}
	if !validSearchSortByValues[sortBy] {
		return awserrors.NewValidationException(fmt.Sprintf("Invalid SortBy: %s", sortBy))
	}
	return nil
}

// validateSortOrder validates SortOrder parameter (ASCENDING/DESCENDING).
// Shared by ListCertificates and SearchCertificates; both declare
// ValidationException, not InvalidParameterException.
func validateSortOrder(sortOrder string) error {
	switch sortOrder {
	case "ASCENDING", "DESCENDING":
		return nil
	default:
		return awserrors.NewValidationException(fmt.Sprintf("Invalid SortOrder: %s. Valid values are ASCENDING, DESCENDING.", sortOrder))
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
	if cert.WasExported {
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
		"SerialNumber": formatSerialNumberHex(parsed.SerialNumber),
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
// UnknownExtKeyUsage). ExtKeyUsageAny is anyExtendedKeyUsage (OID
// 2.5.29.37.0), the ACM enum's ANY value.
var extKeyUsageMap = map[x509.ExtKeyUsage]string{
	x509.ExtKeyUsageAny:             "ANY",
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

// keyUsageList converts ACM KeyUsageName strings to the store's KeyUsage
// entries.
func keyUsageList(names []string) []*acmstorelib.KeyUsage {
	if len(names) == 0 {
		return nil
	}
	kus := make([]*acmstorelib.KeyUsage, len(names))
	for i, name := range names {
		kus[i] = &acmstorelib.KeyUsage{Name: name}
	}
	return kus
}

// x509UsageFields extracts the Key Usage and Extended Key Usage store fields
// from a parsed certificate, so imports persist the same usage data the
// search response derives from the PEM. Extended usages outside the ACM enum
// are recorded by OID with no enum name — the filter input is enum-validated,
// so they correctly match no filter.
func x509UsageFields(parsed *x509.Certificate) ([]*acmstorelib.KeyUsage, []*acmstorelib.ExtendedKeyUsage) {
	kus := keyUsageList(keyUsageFlagsToStrings(parsed.KeyUsage))
	var ekus []*acmstorelib.ExtendedKeyUsage
	for _, eku := range parsed.ExtKeyUsage {
		if name, ok := extKeyUsageMap[eku]; ok {
			ekus = append(ekus, &acmstorelib.ExtendedKeyUsage{Name: name})
		}
	}
	for _, oid := range parsed.UnknownExtKeyUsage {
		ekus = append(ekus, &acmstorelib.ExtendedKeyUsage{OID: oid.String()})
	}
	return kus, ekus
}

// pkixName is an alias for crypto/x509/pkix.Name to simplify the
// buildDistinguishedName signature.
type pkixName = pkix.Name

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

func parseIntToken(token string) (int, error) {
	return strconv.Atoi(token)
}

func formatIntToken(n int) string {
	return strconv.Itoa(n)
}
