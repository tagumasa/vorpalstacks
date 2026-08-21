package acm

import (
	"context"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
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
func (s *ACMService) requestCertificateCore(ctx context.Context, stores *acmStores, in RequestCertificateInput) (string, error) {
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
		if len(in.Tags) > 50 {
			return "", NewTooManyTagsException("Tags must not exceed 50 entries")
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

	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return awserrors.NewResourceNotFoundException("certificate", arn)
		}
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
	Statuses         []string
	KeyTypes         []string
	KeyUsage         []string
	ExtendedKeyUsage []string
	ExportOption     string
	ManagedBy        string
	Origins          []string
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
// shared by the HTTP API and the admin gRPC handler. The HTTP API passes
// filters parsed from request parameters; the admin handler passes empty
// filters to list all certificates.
func (s *ACMService) listCertificatesCore(stores *acmStores, in ListCertificatesInput) (*ListCertificatesResult, error) {
	// Smithy MaxItems: @range(1, 1000). Enforce upper bound; default is
	// applied by the caller when maxItems <= 0.
	if in.MaxItems > 1000 {
		return nil, awserrors.NewValidationException("MaxItems must not exceed 1000")
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

	storeResult, err := stores.certificates.ListWithFilters(filters, in.NextToken, in.MaxItems)
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
