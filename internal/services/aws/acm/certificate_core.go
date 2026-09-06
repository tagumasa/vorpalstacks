package acm

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	types "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"

	vcrypto "vorpalstacks/internal/utils/crypto"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// RequestCertificateInput carries every field that RequestCertificate needs,
// in a format independent of the wire protocol (HTTP Query/JSON vs gRPC-Web).
// Both the HTTP API handler (certificate_operations.go) and the admin gRPC
// handler (admin_handler.go) build this struct from their respective request
// formats and delegate to requestCertificateCore, ensuring that validation,
// certificate generation, and persistence follow a single code path.
type RequestCertificateInput struct {
	DomainName       string
	SANs             []string
	SANsProvided     bool // true when the SANs key was present (empty-list detection)
	IdempotencyToken string
	KeyAlgorithm     string // default "RSA_2048" when empty
	ValidationMethod string // default "DNS" when empty
	ManagedBy        string
	PCAArn           string
	Tags             []types.Tag
	TagsProvided     bool // true when the Tags key was present
	DVOs             []DomainValidationOverride
	DVOsProvided     bool                     // true when the DVOs key was present
	Options          *CertificateOptionsInput // nil = use default
	AccountID        string
	Region           string
}

// DomainValidationOverride is the transport-agnostic representation of a
// user-provided DomainValidationOption entry. Only the DomainName and
// ValidationDomain fields are user-overridable; the rest are auto-generated.
type DomainValidationOverride struct {
	DomainName       string
	ValidationDomain string
}

// CertificateOptionsInput is the transport-agnostic representation of the
// CertificateOptions structure. Both fields use the Smithy enum string
// values ("ENABLED" / "DISABLED").
type CertificateOptionsInput struct {
	CTLPreference string // CertificateTransparencyLoggingPreference
	Export        string // Export option
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// requestCertificateCore is the single entry point for certificate creation
// logic shared by the HTTP API and the admin gRPC handler. It performs all
// Smithy-conformant validation, generates the X.509 material, and persists
// the certificate to the store.
func (s *ACMService) requestCertificateCore(stores *acmStores, in RequestCertificateInput) (string, error) {
	// 1. DomainName (Smithy DomainNameString: REQUIRED).
	domainName, err := validateDomainName(in.DomainName)
	if err != nil {
		return "", err
	}

	// 2. IdempotencyToken (Smithy: OPTIONAL, @length 1-32, @pattern ^\w+$).
	if err := validateIdempotencyToken(in.IdempotencyToken); err != nil {
		return "", err
	}

	// 3. SubjectAlternativeNames (Smithy DomainList: @length 1-100 when provided).
	sans := in.SANs
	if in.SANsProvided {
		if len(sans) == 0 {
			return "", awserrors.NewValidationException("SubjectAlternativeNames must contain at least 1 entry when provided")
		}
		if len(sans) > 100 {
			return "", awserrors.NewValidationException("SubjectAlternativeNames must not exceed 100 entries")
		}
	}
	for _, san := range sans {
		if _, err := validateDomainName(san); err != nil {
			return "", err
		}
	}

	// 4. KeyAlgorithm (Smithy KeyAlgorithm enum, default RSA_2048).
	// validKeyAlgorithmsMap includes RSA_1024 because it is a valid enum
	// value for ImportCertificate (determining algorithm from PEM) and for
	// ListCertificates keyTypes filtering of imported certificates.
	// However, the platform cannot issue new RSA_1024 certificates because
	// vcrypto.GenerateRSAKey enforces a minimum of 2048 bits. Reject RSA_1024
	// at the issuance validation stage with an accurate error message rather
	// than letting it fail at the key-generation stage with a misleading
	// "Unsupported KeyAlgorithm" error.
	keyAlgorithm := in.KeyAlgorithm
	if keyAlgorithm == "" {
		keyAlgorithm = "RSA_2048"
	}
	if !isValidKeyAlgorithm(keyAlgorithm) {
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid KeyAlgorithm: %s. Valid values are RSA_1024, RSA_2048, RSA_3072, RSA_4096, EC_prime256v1, EC_secp384r1, EC_secp521r1.", keyAlgorithm))
	}
	if keyAlgorithm == "RSA_1024" {
		return "", NewInvalidParameterError("KeyAlgorithm RSA_1024 is not available for new certificates. Minimum supported key size is RSA_2048.")
	}

	// 5. ValidationMethod (Smithy ValidationMethod enum, default DNS).
	validationMethod := in.ValidationMethod
	if validationMethod == "" {
		validationMethod = "DNS"
	}
	if !isValidValidationMethod(validationMethod) {
		return "", NewInvalidParameterError(fmt.Sprintf("Invalid ValidationMethod: %s. Valid values are DNS, EMAIL, HTTP.", validationMethod))
	}

	// 6. Generate certificate ID and ARN.
	certId := acmstorelib.GenerateCertificateId()
	certArn := stores.arnBuilder.BuildCertificateARN(certId)

	// 7. Generate the Amazon-issued certificate (key pair + X.509 + PEM).
	cert, err := generateAmazonIssuedCert(certArn, domainName, sans, keyAlgorithm, validationMethod)
	if err != nil {
		return "", err
	}

	// 8. Account / region metadata.
	cert.AccountID = in.AccountID
	cert.Region = in.Region

	// 9. ManagedBy (Smithy CertificateManagedBy enum: CLOUDFRONT only).
	if in.ManagedBy != "" {
		if !isValidManagedBy(in.ManagedBy) {
			return "", NewInvalidParameterError(fmt.Sprintf("Invalid ManagedBy value: %s", in.ManagedBy))
		}
		cert.ManagedBy = in.ManagedBy
	}

	// 10. CertificateAuthorityArn (Smithy PcaArn: OPTIONAL).
	pcaArn, err := validateCertificateAuthorityArn(in.PCAArn)
	if err != nil {
		return "", err
	}
	cert.CertificateAuthorityArn = pcaArn

	// 11. Tags (Smithy TagList: @length 1-50 when provided).
	if in.TagsProvided {
		if len(in.Tags) == 0 {
			return "", awserrors.NewValidationException("Tags must contain at least 1 entry when provided")
		}
		if len(in.Tags) > maxTagsPerCertificate {
			return "", NewTooManyTagsException(fmt.Sprintf("Tags must not exceed %d entries", maxTagsPerCertificate))
		}
	}
	if err := validateACMTags(in.Tags); err != nil {
		return "", err
	}
	cert.Tags = in.Tags

	// 12. DomainValidationOptions overrides (Smithy: @length 1-100 when provided).
	if in.DVOsProvided {
		if len(in.DVOs) == 0 {
			return "", awserrors.NewValidationException("DomainValidationOptions must contain at least 1 entry when provided")
		}
		if len(in.DVOs) > 100 {
			return "", awserrors.NewValidationException("DomainValidationOptions must not exceed 100 entries")
		}
	}
	userMap := make(map[string]string)
	for _, dvo := range in.DVOs {
		if err := validateDomainValidationFields(dvo.DomainName, dvo.ValidationDomain); err != nil {
			return "", err
		}
		userMap[strings.ToLower(dvo.DomainName)] = strings.ToLower(dvo.ValidationDomain)
	}
	for _, dv := range cert.DomainValidationOptions {
		if vd, ok := userMap[strings.ToLower(dv.DomainName)]; ok {
			dv.ValidationDomain = vd
		}
	}

	// 13. Options (Smithy CertificateOptions, default CTL=ENABLED, Export=DISABLED).
	if in.Options != nil {
		ctlp := in.Options.CTLPreference
		if ctlp == "" {
			ctlp = "ENABLED"
		}
		if !isValidCertificateTransparencyLoggingPreference(ctlp) {
			return "", NewInvalidParameterError(fmt.Sprintf("Invalid CertificateTransparencyLoggingPreference: %s. Valid values are ENABLED, DISABLED.", ctlp))
		}
		exportOpt := in.Options.Export
		if exportOpt == "" {
			exportOpt = "DISABLED"
		}
		if !isValidExportOption(exportOpt) {
			return "", NewInvalidParameterError(fmt.Sprintf("Invalid Export: %s. Valid values are ENABLED, DISABLED.", exportOpt))
		}
		cert.Options = &acmstorelib.CertificateOptions{
			CertificateTransparencyLoggingPreference: ctlp,
			Export:                                   exportOpt,
		}
	} else {
		cert.Options = &acmstorelib.CertificateOptions{
			CertificateTransparencyLoggingPreference: "ENABLED",
			Export:                                   "DISABLED",
		}
	}

	// 14. Persist to store.
	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return "", awserrors.NewConflictException("Certificate already exists")
		}
		return "", err
	}

	return certArn, nil
}

// deleteCertificateCore is the single entry point for certificate deletion
// shared by the HTTP API and the admin gRPC handler.
func (s *ACMService) deleteCertificateCore(stores *acmStores, arn string) error {
	if arn == "" {
		return awserrors.NewValidationException("CertificateArn is required")
	}
	if err := validateCertificateArn(arn); err != nil {
		return err
	}

	cert, err := s.fetchCertificate(stores, arn)
	if err != nil {
		return err
	}

	if len(cert.InUseBy) > 0 {
		return NewResourceInUseError("certificate", arn)
	}

	return stores.certificates.Delete(arn)
}

// ---------------------------------------------------------------------------
// Transport-agnostic DTOs for ListCertificates
// ---------------------------------------------------------------------------

// ListCertificatesInput carries every field that ListCertificates needs,
// in a format independent of the wire protocol. Both the HTTP API handler
// and the admin gRPC handler build this struct and delegate to
// listCertificatesCore, ensuring that pagination validation and filter
// application follow a single code path.
type ListCertificatesInput struct {
	NextToken        string
	MaxItems         int
	MaxItemsSet      bool
	Statuses         []string
	KeyTypes         []string
	KeyUsage         []string
	ExtendedKeyUsage []string
	ExportOption     string
	ManagedBy        string
	Origins          []string
	OriginsProvided  bool
	SortBy           string
	SortOrder        string
}

// CertSummaryOut is the transport-agnostic certificate summary returned by
// listCertificatesCore. It mirrors the store-layer CertificateSummary but
// does not carry Pebble JSON tags, keeping the store layer an internal
// implementation detail invisible to transport handlers.
type CertSummaryOut struct {
	CertificateArn                       string
	DomainName                           string
	SubjectAlternativeNameSummaries      []string
	HasAdditionalSubjectAlternativeNames bool
	Status                               string
	Type                                 string
	RenewalEligibility                   string
	KeyAlgorithm                         string
	ManagedBy                            string
	CertificateKeyPairOrigin             string
	KeyUsages                            []string
	ExtendedKeyUsages                    []string
	InUse                                bool
	NotBefore                            float64
	NotAfter                             float64
	CreatedAt                            float64
	IssuedAt                             float64
	ImportedAt                           float64
	Exported                             bool
	ExportOption                         string
	RevokedAt                            float64
}

// ListCertificatesResult is the transport-agnostic result of listing
// certificates.
type ListCertificatesResult struct {
	Certificates []CertSummaryOut
	NextToken    string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// listCertificatesCore is the single entry point for certificate listing
// shared by the HTTP API and the admin gRPC handler. All parameter
// validation happens here so both planes enforce identical rules; the
// HTTP API passes filters parsed from request parameters, the admin
// handler passes proto-extracted values.
func (s *ACMService) listCertificatesCore(stores *acmStores, in ListCertificatesInput) (*ListCertificatesResult, error) {
	if err := validateNextToken(in.NextToken); err != nil {
		return nil, err
	}

	// Smithy MaxItems: @range(1-1000). When explicitly provided, validate
	// the range; when absent, default to pagination.DefaultMaxItems.
	maxItems := pagination.DefaultMaxItems
	if in.MaxItemsSet {
		if in.MaxItems < 1 || in.MaxItems > listCertificatesMaxItems {
			return nil, awserrors.NewValidationException(
				fmt.Sprintf("MaxItems must be between 1 and %d", listCertificatesMaxItems))
		}
		maxItems = in.MaxItems
	}

	// Smithy SortBy enum for ListCertificates only permits CREATED_AT.
	// ListCertificates declares InvalidArgsException and ValidationException
	// only, so constraint violations surface as ValidationException.
	if in.SortBy != "" && in.SortBy != "CREATED_AT" {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Invalid SortBy: %s. Valid values are CREATED_AT.", in.SortBy))
	}
	if in.SortOrder != "" {
		if err := validateSortOrder(in.SortOrder); err != nil {
			return nil, err
		}
	}
	// Both sort members document the pair as mandatory: "If you specify
	// SortBy, you must also specify SortOrder." and the reverse.
	if (in.SortBy != "") != (in.SortOrder != "") {
		return nil, awserrors.NewValidationException("If you specify SortBy, you must also specify SortOrder; if you specify SortOrder, you must also specify SortBy.")
	}

	if err := validateEnumList(in.Statuses, "CertificateStatuses", validCertificateStatuses); err != nil {
		return nil, err
	}

	// Smithy CertificateKeyPairOrigins: @length(1-3) when provided.
	if in.OriginsProvided {
		if len(in.Origins) == 0 {
			return nil, awserrors.NewValidationException("CertificateKeyPairOrigins must contain at least 1 entry when provided")
		}
		if len(in.Origins) > 3 {
			return nil, awserrors.NewValidationException("CertificateKeyPairOrigins must not exceed 3 entries")
		}
	}
	if err := validateEnumList(in.Origins, "CertificateKeyPairOrigins", validKeyPairOrigins); err != nil {
		return nil, err
	}

	if err := validateEnumList(in.KeyTypes, "keyTypes", validKeyAlgorithmsMap); err != nil {
		return nil, err
	}
	if err := validateEnumList(in.KeyUsage, "keyUsage", validKeyUsageNames); err != nil {
		return nil, err
	}
	if err := validateEnumList(in.ExtendedKeyUsage, "extendedKeyUsage", validExtendedKeyUsageNames); err != nil {
		return nil, err
	}
	// Validate single-value enum filters that are not lists.
	if in.ExportOption != "" {
		if err := validateSingleEnum(in.ExportOption, "exportOption", validExportOptionValues); err != nil {
			return nil, err
		}
	}
	if in.ManagedBy != "" {
		if err := validateSingleEnum(in.ManagedBy, "managedBy", validManagedByValues); err != nil {
			return nil, err
		}
	}

	filters := acmstorelib.ListFilters{
		Statuses:         in.Statuses,
		KeyTypes:         in.KeyTypes,
		KeyUsage:         in.KeyUsage,
		ExtendedKeyUsage: in.ExtendedKeyUsage,
		ExportOption:     in.ExportOption,
		ManagedBy:        in.ManagedBy,
		Origins:          in.Origins,
		SortBy:           in.SortBy,
		SortOrder:        in.SortOrder,
	}

	storeResult, err := stores.certificates.ListWithFilters(filters, in.NextToken, maxItems)
	if err != nil {
		return nil, err
	}

	certs := make([]CertSummaryOut, len(storeResult.Certificates))
	for i, sc := range storeResult.Certificates {
		certs[i] = certSummaryStoreToService(sc)
	}

	return &ListCertificatesResult{
		Certificates: certs,
		NextToken:    storeResult.NextToken,
	}, nil
}

// certSummaryStoreToService converts a store-layer CertificateSummary to
// the transport-agnostic CertSummaryOut.
func certSummaryStoreToService(sc *acmstorelib.CertificateSummary) CertSummaryOut {
	return CertSummaryOut{
		CertificateArn:                       sc.CertificateArn,
		DomainName:                           sc.DomainName,
		SubjectAlternativeNameSummaries:      sc.SubjectAlternativeNameSummaries,
		HasAdditionalSubjectAlternativeNames: sc.HasAdditionalSubjectAlternativeNames,
		Status:                               sc.Status,
		Type:                                 sc.Type,
		RenewalEligibility:                   sc.RenewalEligibility,
		KeyAlgorithm:                         sc.KeyAlgorithm,
		ManagedBy:                            sc.ManagedBy,
		CertificateKeyPairOrigin:             sc.CertificateKeyPairOrigin,
		KeyUsages:                            sc.KeyUsages,
		ExtendedKeyUsages:                    sc.ExtendedKeyUsages,
		InUse:                                sc.InUse,
		NotBefore:                            sc.NotBefore,
		NotAfter:                             sc.NotAfter,
		CreatedAt:                            sc.CreatedAt,
		IssuedAt:                             sc.IssuedAt,
		ImportedAt:                           sc.ImportedAt,
		Exported:                             sc.Exported,
		ExportOption:                         sc.ExportOption,
		RevokedAt:                            sc.RevokedAt,
	}
}

// ---------------------------------------------------------------------------
// Certificate read / lifecycle Core functions
// ---------------------------------------------------------------------------

// Smithy length and range limits for the import/export paths (decoded byte
// counts) plus the per-certificate tag cap. These are the single definition
// sites for the raw numbers.
const (
	maxCertificateBodyBytes  = 32768
	maxPrivateKeyBytes       = 5120
	maxCertificateChainBytes = 2097152
	minPassphraseBytes       = 4
	maxPassphraseBytes       = 128
	maxTagsPerCertificate    = 50
)

// requireCertificateArn is the validation half of parseCertificateArn: it
// rejects an absent or malformed certificate ARN. It lives in the core file so
// that handlers carry wire extraction only; parseCertificateArn (helpers.go)
// delegates here for the tag handlers that still parse from parameters.
func requireCertificateArn(arn, paramName string) error {
	if arn == "" {
		return awserrors.NewValidationException(paramName + " is required")
	}
	return validateCertificateArn(arn)
}

// fetchCertificate loads a certificate by ARN with the shared not-found
// mapping: a missing certificate is reported as ResourceNotFoundException
// and real storage errors propagate unchanged.
func (s *ACMService) fetchCertificate(stores *acmStores, arn string) (*acmstorelib.Certificate, error) {
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}
	return cert, nil
}

// getCertificateCore fetches a certificate by ARN, mapping a missing
// certificate to ResourceNotFoundException. Shared by the GetCertificate and
// DescribeCertificate handlers.
func (s *ACMService) getCertificateCore(stores *acmStores, arn string) (*acmstorelib.Certificate, error) {
	if err := requireCertificateArn(arn, "CertificateArn"); err != nil {
		return nil, err
	}
	return s.fetchCertificate(stores, arn)
}

// ResendValidationEmailInput carries the wire-extracted fields of
// ResendValidationEmail.
type ResendValidationEmailInput struct {
	CertificateArn   string
	Domain           string
	ValidationDomain string
}

// resendValidationEmailCore is the single validation + persistence path for
// ResendValidationEmail.
func (s *ACMService) resendValidationEmailCore(stores *acmStores, in ResendValidationEmailInput) error {
	if err := requireCertificateArn(in.CertificateArn, "CertificateArn"); err != nil {
		return err
	}
	cert, err := s.fetchCertificate(stores, in.CertificateArn)
	if err != nil {
		return err
	}

	if cert.Type == "IMPORTED" {
		return NewInvalidStateException("Certificate is not in PENDING_VALIDATION state")
	}
	// AWS only allows ResendValidationEmail for certificates in
	// PENDING_VALIDATION status. ISSUED certificates have already completed
	// validation and resend is a no-op that AWS rejects.
	if cert.Status != "PENDING_VALIDATION" {
		return NewInvalidStateException("Certificate is not in PENDING_VALIDATION state")
	}

	// Domain and ValidationDomain are required per Smithy model.
	if in.Domain == "" {
		return awserrors.NewValidationException("Domain is required")
	}
	if in.ValidationDomain == "" {
		return awserrors.NewValidationException("ValidationDomain is required")
	}

	// In AWS, ResendValidationEmail resends the domain validation email.
	// For edge deployment, email sending is not available, but we honour the
	// parameters by updating ValidationDomain.
	changed := false
	for _, dvo := range cert.DomainValidationOptions {
		if dvo.DomainName != in.Domain {
			continue
		}
		dvo.ValidationDomain = in.ValidationDomain
		changed = true
	}

	if !changed {
		return NewInvalidDomainValidationOptionsException(fmt.Sprintf("Domain %s not found in certificate validation options", in.Domain))
	}

	return stores.certificates.Update(cert)
}

// ImportCertificateInput carries the raw (still base64-encoded) wire values of
// ImportCertificate. Decoding and length validation happen in the Core so the
// handler performs wire extraction only.
type ImportCertificateInput struct {
	Certificate      string
	PrivateKey       string
	CertificateChain string
	ExistingArn      string
	Tags             []types.Tag
}

// importCertificateCore imports a certificate into ACM. When ExistingArn is
// provided, the existing certificate is updated (re-import); otherwise a new
// IMPORTED certificate is created.
func (s *ACMService) importCertificateCore(stores *acmStores, in ImportCertificateInput) (string, error) {
	certificate := in.Certificate
	if certificate == "" {
		return "", awserrors.NewValidationException("Certificate is required")
	}
	certificate, err := decodeBase64PEM(certificate)
	if err != nil {
		return "", err
	}
	if len(certificate) > maxCertificateBodyBytes {
		return "", awserrors.NewValidationException(fmt.Sprintf("Certificate exceeds maximum length of %d bytes", maxCertificateBodyBytes))
	}

	privateKey := in.PrivateKey

	existingArn := in.ExistingArn

	// PrivateKey is required for initial import (no existingArn).
	// For re-import (existingArn set), PrivateKey is optional — the
	// existing key is retained if not provided.
	if existingArn == "" && privateKey == "" {
		return "", awserrors.NewValidationException("PrivateKey is required for initial certificate import")
	}
	if privateKey != "" {
		privateKey, err = decodeBase64PEM(privateKey)
		if err != nil {
			return "", err
		}
		if len(privateKey) > maxPrivateKeyBytes {
			return "", awserrors.NewValidationException(fmt.Sprintf("PrivateKey exceeds maximum length of %d bytes", maxPrivateKeyBytes))
		}
	}
	certificateChain := in.CertificateChain
	if certificateChain != "" {
		certificateChain, err = decodeBase64PEM(certificateChain)
		if err != nil {
			return "", err
		}
		if len(certificateChain) > maxCertificateChainBytes {
			return "", awserrors.NewValidationException(fmt.Sprintf("CertificateChain exceeds maximum length of %d bytes", maxCertificateChainBytes))
		}
	}

	// --- Re-import path: update an existing certificate ---
	if existingArn != "" {
		if err := validateCertificateArn(existingArn); err != nil {
			return "", err
		}
		existing, err := s.fetchCertificate(stores, existingArn)
		if err != nil {
			return "", err
		}

		// Only IMPORTED certificates can be re-imported.
		if existing.Type != "IMPORTED" {
			return "", awserrors.NewInvalidParameterException("Only imported certificates can be re-imported")
		}

		// If Tags are provided on re-import, validate and replace.
		// If not provided, existing tags are retained.
		if len(in.Tags) > 0 {
			if len(in.Tags) > maxTagsPerCertificate {
				return "", NewTooManyTagsException(fmt.Sprintf("Tags must not exceed %d entries", maxTagsPerCertificate))
			}
			if err := validateACMTags(in.Tags); err != nil {
				return "", err
			}
			existing.Tags = in.Tags
		}

		// Replace certificate material. Fields not provided on re-import
		// (e.g. PrivateKey) retain their existing values.
		existing.Certificate = certificate
		if certificateChain != "" {
			existing.CertificateChain = certificateChain
		}
		if privateKey != "" {
			existing.PrivateKey = privateKey
		}
		existing.ImportedAt = time.Now().UTC()

		// Re-extract X.509 attributes from the updated certificate PEM.
		if parsedCert, _ := vcrypto.ParseCertificatePEM([]byte(certificate)); parsedCert != nil {
			existing.NotBefore = parsedCert.NotBefore
			existing.NotAfter = parsedCert.NotAfter
			existing.KeyAlgorithm = determineKeyAlgorithmFromParsed(parsedCert)
			existing.SignatureAlgorithm = determineSignatureAlgorithmFromParsed(parsedCert)
			existing.Subject = parsedCert.Subject.String()
			existing.Issuer = parsedCert.Issuer.String()
			if domain, err := domainFromParsedCert(parsedCert); err == nil {
				existing.DomainName = domain
			}
		}

		if err := stores.certificates.Update(existing); err != nil {
			return "", err
		}

		return existingArn, nil
	}

	// --- Initial import path: create a new certificate ---
	tags := in.Tags
	if err := validateACMTags(tags); err != nil {
		return "", err
	}

	certId := acmstorelib.GenerateCertificateId()
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	// Parse the PEM once and derive every X.509-derived field from the
	// single result; the divergent silent fallbacks of the former
	// per-field helpers (empty serial, default key algorithm) were
	// dead-in-practice branches around this same parse.
	parsedCert, err := vcrypto.ParseCertificatePEM([]byte(certificate))
	if err != nil {
		return "", awserrors.NewValidationException(fmt.Sprintf("Invalid certificate: failed to parse certificate: %v", err))
	}
	domainName, err := domainFromParsedCert(parsedCert)
	if err != nil {
		return "", awserrors.NewValidationException(fmt.Sprintf("Invalid certificate: %v", err))
	}
	now := time.Now().UTC()
	certKeyUsages, certExtKeyUsages := x509UsageFields(parsedCert)
	cert := &acmstorelib.Certificate{
		CertificateArn:           certificateArn,
		DomainName:               domainName,
		Serial:                   formatSerialNumberHex(parsedCert.SerialNumber),
		Status:                   "ISSUED",
		Type:                     "IMPORTED",
		KeyAlgorithm:             determineKeyAlgorithmFromParsed(parsedCert),
		SignatureAlgorithm:       determineSignatureAlgorithmFromParsed(parsedCert),
		RenewalEligibility:       "INELIGIBLE",
		CreatedAt:                now,
		ImportedAt:               now,
		NotBefore:                parsedCert.NotBefore,
		NotAfter:                 parsedCert.NotAfter,
		KeyUsages:                certKeyUsages,
		ExtendedKeyUsages:        certExtKeyUsages,
		Certificate:              certificate,
		CertificateChain:         certificateChain,
		PrivateKey:               privateKey,
		Tags:                     tags,
		Subject:                  parsedCert.Subject.String(),
		Issuer:                   parsedCert.Issuer.String(),
		AccountID:                stores.accountID,
		Region:                   stores.region,
		CertificateKeyPairOrigin: "CUSTOMER_PROVIDED",
	}

	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return "", awserrors.NewConflictException("Certificate already exists")
		}
		return "", err
	}

	return certificateArn, nil
}

// UpdateCertificateOptionsInput carries the raw wire values of
// UpdateCertificateOptions; the Options structure is type-checked in the Core.
type UpdateCertificateOptionsInput struct {
	CertificateArn string
	OptionsRaw     interface{}
}

// updateCertificateOptionsCore is the single validation + persistence path for
// UpdateCertificateOptions.
func (s *ACMService) updateCertificateOptionsCore(stores *acmStores, in UpdateCertificateOptionsInput) error {
	if err := requireCertificateArn(in.CertificateArn, "CertificateArn"); err != nil {
		return err
	}
	cert, err := s.fetchCertificate(stores, in.CertificateArn)
	if err != nil {
		return err
	}

	// Smithy declares InvalidStateException for UpdateCertificateOptions —
	// only ISSUED certs can have their options updated.
	if cert.Status != "ISSUED" {
		return NewInvalidStateException("Certificate is not in a valid state for updating options.")
	}

	optionsMap, ok := in.OptionsRaw.(map[string]interface{})
	if !ok {
		return awserrors.NewValidationException("Options are required")
	}

	ctlp, err := parseCertificateTransparencyLoggingPreference(optionsMap)
	if err != nil {
		return err
	}
	exportOpt, err := parseExportOption(optionsMap)
	if err != nil {
		return err
	}
	cert.Options = &acmstorelib.CertificateOptions{
		CertificateTransparencyLoggingPreference: ctlp,
		Export:                                   exportOpt,
	}

	return stores.certificates.Update(cert)
}

// renewCertificateCore is the single validation + persistence path for
// RenewCertificate.
func (s *ACMService) renewCertificateCore(stores *acmStores, arn string) error {
	if err := requireCertificateArn(arn, "CertificateArn"); err != nil {
		return err
	}
	cert, err := s.fetchCertificate(stores, arn)
	if err != nil {
		return err
	}

	if cert.Type != "AMAZON_ISSUED" {
		return awserrors.NewValidationException("Certificate type is not supported. Only Amazon-issued certificates can be renewed.")
	}

	// Smithy declares RequestInProgressException for RenewCertificate —
	// a PENDING_VALIDATION cert is already being processed.
	if cert.Status == "PENDING_VALIDATION" {
		return NewRequestInProgressException("Certificate request is in progress.")
	}

	// Only ISSUED or EXPIRED certs can be renewed. All other states
	// (FAILED, REVOKED, etc.) are invalid for renewal.
	if cert.Status != "ISSUED" && cert.Status != "EXPIRED" {
		return awserrors.NewValidationException("Certificate is not in a valid state for renewal.")
	}

	if cert.RenewalEligibility == "INELIGIBLE" {
		return awserrors.NewValidationException("Certificate is not eligible for renewal.")
	}

	// Regenerate the X.509 certificate material (new key pair, new serial,
	// new PEM). The ARN is preserved per AWS specification: "When ACM renews
	// a certificate, the certificate's Amazon Resource Name (ARN) remains
	// the same."
	if err := renewCertificateMaterial(cert); err != nil {
		return NewInternalServerException(fmt.Sprintf("Failed to regenerate certificate material: %v", err))
	}

	now := time.Now().UTC()
	// Populate RenewalSummary.DomainValidationOptions (Smithy REQUIRED)
	// by deep-copying the existing DVOs with updated validation status.
	renewalDVOs := make([]*acmstorelib.DomainValidation, len(cert.DomainValidationOptions))
	for i, dvo := range cert.DomainValidationOptions {
		renewalDVOs[i] = &acmstorelib.DomainValidation{
			DomainName:       dvo.DomainName,
			ValidationDomain: dvo.ValidationDomain,
			ValidationMethod: dvo.ValidationMethod,
			ValidationStatus: "SUCCESS",
			ValidationEmails: dvo.ValidationEmails,
			ResourceRecord:   dvo.ResourceRecord,
			HttpRedirect:     dvo.HttpRedirect,
		}
	}
	cert.RenewalSummary = &acmstorelib.RenewalSummary{
		RenewalStatus:           "SUCCESS",
		UpdatedAt:               now,
		DomainValidationOptions: renewalDVOs,
	}
	for i := range cert.DomainValidationOptions {
		cert.DomainValidationOptions[i].ValidationStatus = "SUCCESS"
	}

	return stores.certificates.Update(cert)
}

// ExportCertificateInput carries the wire-extracted fields of
// ExportCertificate; Passphrase is still base64-encoded wire text.
type ExportCertificateInput struct {
	CertificateArn string
	Passphrase     string
}

// ExportCertificateResult is the transport-agnostic export payload. The
// private key is encrypted with the caller-supplied passphrase.
type ExportCertificateResult struct {
	Certificate      string
	CertificateChain string
	PrivateKey       string
}

// exportCertificateCore is the single validation + persistence path for
// ExportCertificate.
func (s *ACMService) exportCertificateCore(stores *acmStores, in ExportCertificateInput) (*ExportCertificateResult, error) {
	if err := requireCertificateArn(in.CertificateArn, "CertificateArn"); err != nil {
		return nil, err
	}
	cert, err := s.fetchCertificate(stores, in.CertificateArn)
	if err != nil {
		return nil, err
	}

	// Amazon-issued certificates are never exportable even though the
	// platform retains their key pairs for TLS termination; only imported
	// (and private-CA-issued) material can leave the store.
	if cert.Type == "AMAZON_ISSUED" {
		return nil, awserrors.NewValidationException("Certificate does not have an exportable private key. Only imported certificates with private keys can be exported.")
	}

	if cert.PrivateKey == "" {
		return nil, awserrors.NewValidationException("Certificate does not have an exportable private key. Only imported certificates with private keys can be exported.")
	}

	if in.Passphrase == "" {
		return nil, awserrors.NewValidationException("Passphrase is required")
	}
	// Smithy defines Passphrase as a Blob (base64-encoded). Reject if the
	// input is not valid base64 rather than silently falling back to raw
	// bytes, which would be a fail-OPEN acceptance of malformed input.
	passphraseBytes, err := base64.StdEncoding.DecodeString(in.Passphrase)
	if err != nil {
		return nil, awserrors.NewValidationException("Passphrase must be valid base64-encoded data")
	}
	if passLen := len(passphraseBytes); passLen < minPassphraseBytes || passLen > maxPassphraseBytes {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Passphrase must be between %d and %d bytes", minPassphraseBytes, maxPassphraseBytes))
	}

	encryptedKey, err := encryptPrivateKey(cert.PrivateKey, string(passphraseBytes))
	if err != nil {
		return nil, awserrors.NewValidationException("Failed to encrypt private key")
	}

	// Record that this certificate has been exported so that ListCertificates
	// and SearchCertificates report the correct Exported state.
	if !cert.WasExported {
		cert.WasExported = true
		if err := stores.certificates.Update(cert); err != nil {
			return nil, NewInternalServerException("Failed to update certificate export state")
		}
	}

	return &ExportCertificateResult{
		Certificate:      cert.Certificate,
		CertificateChain: cert.CertificateChain,
		PrivateKey:       encryptedKey,
	}, nil
}

// RevokeCertificateInput carries the wire-extracted fields of
// RevokeCertificate. RevocationReasonRaw is the raw wire value; presence,
// type, and enum validation happen in the Core.
type RevokeCertificateInput struct {
	CertificateArn      string
	RevocationReasonRaw interface{}
}

// revokeCertificateCore is the single validation + persistence path for
// RevokeCertificate. It returns the ARN of the revoked certificate.
func (s *ACMService) revokeCertificateCore(stores *acmStores, in RevokeCertificateInput) (string, error) {
	arn := in.CertificateArn
	if err := requireCertificateArn(arn, "CertificateArn"); err != nil {
		return "", err
	}
	cert, err := s.fetchCertificate(stores, arn)
	if err != nil {
		return "", err
	}

	if cert.Status == "REVOKED" {
		return "", awserrors.NewResourceNotFoundException("certificate", arn)
	}

	// Smithy declares ResourceInUseException for RevokeCertificate — reject
	// when the cert is actively in use by CloudFront, API Gateway, etc.
	if len(cert.InUseBy) > 0 {
		return "", NewResourceInUseError("certificate", arn)
	}

	if cert.Status != "ISSUED" {
		return "", awserrors.NewValidationException("Certificate is not in a valid state for revocation.")
	}

	// RevocationReason is required per Smithy model.
	if in.RevocationReasonRaw == nil {
		return "", awserrors.NewValidationException("RevocationReason is required")
	}
	reason, ok := in.RevocationReasonRaw.(string)
	if !ok || !isValidRevocationReason(reason) {
		return "", awserrors.NewValidationException(fmt.Sprintf("Invalid RevocationReason: %v", in.RevocationReasonRaw))
	}
	cert.RevocationReason = reason
	cert.Status = "REVOKED"
	cert.RevokedAt = time.Now().UTC()

	if err := stores.certificates.Update(cert); err != nil {
		return "", err
	}

	return arn, nil
}
