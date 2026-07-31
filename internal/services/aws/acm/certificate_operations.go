package acm

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"
	"vorpalstacks/internal/utils/aws/types"

	vcrypto "vorpalstacks/internal/utils/crypto"

	"golang.org/x/crypto/pbkdf2"
)

func (s *ACMService) acmTagConfig(stores *acmStores, arn string) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:    "CertificateArn",
			TagsParam:        "Tags",
			TagKeysParam:     "Tags",
			TagKeyName:       "Key",
			TagValueName:     "Value",
			RequireTags:      true,
			RequireTagKeys:   true,
			RequireResource:  true,
			UseQueryFallback: true,
		},
		ResourceKey: func(rawKey string) string {
			return arn
		},
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			_, err := stores.certificates.Get(arn)
			if err != nil {
				if acmstorelib.IsNotFound(err) {
					return awserrors.NewResourceNotFoundException("certificate", arn)
				}
				return err
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []types.Tag) error {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return err
			}
			cert.Tags = tagutil.Apply(cert.Tags, tagList)
			return stores.certificates.Update(cert)
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			tagsToRemove := tagutil.ParseTagsWithQueryFallback(params, "Tags")
			keys := make([]string, 0, len(tagsToRemove))
			for _, t := range tagsToRemove {
				keys = append(keys, t.Key)
			}
			return keys
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return err
			}
			tagKeySet := make(map[string]bool, len(tagKeys))
			for _, k := range tagKeys {
				tagKeySet[k] = true
			}
			cert.Tags = tagutil.Remove(cert.Tags, tagKeySet)
			return stores.certificates.Update(cert)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return nil, err
			}
			return cert.Tags, nil
		},
		FormatResponse: func(tagList []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToResponse(tagList),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: func(err error) error {
			switch err.(type) {
			case *tagutil.MissingResourceError:
				return awserrors.NewResourceNotFoundException("certificate", "")
			case *tagutil.MissingTagsError:
				return awserrors.NewValidationException("Tags are required")
			case *tagutil.MissingTagKeysError:
				return awserrors.NewValidationException("Tag keys are required")
			}
			return err
		},
	}
}

// acmGenericTagConfig returns a TagHandlerConfig for the generic TagResource /
// UntagResource / ListTagsForResource operations. These use "ResourceArn" as
// the resource parameter and "TagKeys" (a plain list of key strings) for
// untagging, instead of the certificate-specific parameter names.
func (s *ACMService) acmGenericTagConfig(stores *acmStores, arn string) tagutil.TagHandlerConfig {
	config := s.acmTagConfig(stores, arn)
	config.Param.ResourceParam = "ResourceArn"
	config.Param.TagKeysParam = "TagKeys"
	config.ParseTagKeys = func(params map[string]interface{}) []string {
		return tagutil.ParseTagKeysWithQueryFallback(params, "TagKeys")
	}
	return config
}

// RequestCertificate requests a new certificate from ACM.
func (s *ACMService) RequestCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	domainName, err := parseDomainName(params)
	if err != nil {
		return nil, err
	}

	// Smithy IdempotencyToken: @length(1-32) + @pattern(^\w+$).
	// Optional field — validate only if provided.
	if err := validateIdempotencyToken(request.GetStringParam(params, "IdempotencyToken")); err != nil {
		return nil, err
	}

	certId := acmstorelib.GenerateCertificateId()
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	validationMethod, err := parseValidationMethod(params)
	if err != nil {
		return nil, err
	}
	keyAlgorithm, err := parseKeyAlgorithm(params)
	if err != nil {
		return nil, err
	}

	// Parse SubjectAlternativeNames early so they can be embedded in the x509 template.
	// Smithy DomainList requires at least 1 entry when provided (min 1, max 100).
	var sans []string
	if sansRaw := params["SubjectAlternativeNames"]; sansRaw != nil {
		if sansArr, ok := sansRaw.([]interface{}); ok {
			for _, san := range sansArr {
				if s, ok := san.(string); ok {
					sans = append(sans, s)
				}
			}
		}
		if len(sans) == 0 {
			return nil, NewInvalidParameterError("SubjectAlternativeNames must contain at least 1 entry when provided")
		}
		if len(sans) > 100 {
			return nil, NewInvalidParameterError("SubjectAlternativeNames must not exceed 100 entries")
		}
	}

	cert, err := generateAmazonIssuedCert(certificateArn, domainName, sans, keyAlgorithm, validationMethod)
	if err != nil {
		return nil, err
	}

	managedBy, err := parseManagedBy(params)
	if err != nil {
		return nil, err
	}

	pcaArn, err := parseCertificateAuthorityArn(params)
	if err != nil {
		return nil, err
	}

	cert.AccountID = reqCtx.GetAccountID()
	cert.Region = reqCtx.GetRegion()
	cert.CertificateAuthorityArn = pcaArn
	cert.ManagedBy = managedBy

	// Merge user-provided DomainValidationOptions (e.g. ValidationDomain
	// for EMAIL) into the auto-generated ones.
	if err := applyUserDomainValidationOptions(cert.DomainValidationOptions, params); err != nil {
		return nil, err
	}

	tags := tagutil.ParseTagsWithQueryFallback(params, "Tags")
	if len(tags) > 50 {
		return nil, NewTooManyTagsException("Tags must not exceed 50 entries")
	}
	if err := validateACMTags(tags); err != nil {
		return nil, err
	}
	cert.Tags = tags

	if optionsMap, ok := params["Options"].(map[string]interface{}); ok {
		ctlp, err := parseCertificateTransparencyLoggingPreference(optionsMap)
		if err != nil {
			return nil, err
		}
		exportOpt, err := parseExportOption(optionsMap)
		if err != nil {
			return nil, err
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

	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return nil, awserrors.NewConflictException("Certificate already exists")
		}
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certificateArn,
	}, nil
}

// GetCertificate retrieves the certificate and certificate chain for the specified ARN.
func (s *ACMService) GetCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	result := map[string]interface{}{
		"Certificate": cert.Certificate,
	}
	if cert.CertificateChain != "" {
		result["CertificateChain"] = cert.CertificateChain
	}
	return result, nil
}

// ListCertificates retrieves a list of certificates for the account.
func (s *ACMService) ListCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	nextToken := request.GetStringParam(params, "NextToken")
	maxItems := getMaxItems(params)

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	filters := acmstorelib.ListFilters{
		SortBy:    request.GetStringParam(params, "SortBy"),
		SortOrder: request.GetStringParam(params, "SortOrder"),
	}

	// Smithy SortBy enum for ListCertificates only permits CREATED_AT.
	if filters.SortBy != "" && filters.SortBy != "CREATED_AT" {
		return nil, NewInvalidParameterError(fmt.Sprintf("Invalid SortBy: %s. Valid values are CREATED_AT.", filters.SortBy))
	}
	if filters.SortOrder != "" {
		if err := validateSortOrder(filters.SortOrder); err != nil {
			return nil, err
		}
	}

	filters.Statuses = request.GetStringList(params, "CertificateStatuses")
	if err := validateEnumList(filters.Statuses, "CertificateStatuses", validCertificateStatuses); err != nil {
		return nil, err
	}

	filters.Origins = request.GetStringList(params, "CertificateKeyPairOrigins")
	if err := validateEnumList(filters.Origins, "CertificateKeyPairOrigins", validKeyPairOrigins); err != nil {
		return nil, err
	}

	if includes, ok := params["Includes"].(map[string]interface{}); ok {
		filters.KeyTypes = request.GetStringList(includes, "keyTypes")
		filters.KeyUsage = request.GetStringList(includes, "keyUsage")
		filters.ExtendedKeyUsage = request.GetStringList(includes, "extendedKeyUsage")
		filters.ExportOption = request.GetStringParam(includes, "exportOption")
		filters.ManagedBy = request.GetStringParam(includes, "managedBy")
	}
	if err := validateEnumList(filters.KeyTypes, "keyTypes", validKeyAlgorithmsMap); err != nil {
		return nil, err
	}
	if err := validateEnumList(filters.KeyUsage, "keyUsage", validKeyUsageNames); err != nil {
		return nil, err
	}
	if err := validateEnumList(filters.ExtendedKeyUsage, "extendedKeyUsage", validExtendedKeyUsageNames); err != nil {
		return nil, err
	}
	// Validate single-value enum filters that are not lists.
	if filters.ExportOption != "" {
		if err := validateSingleEnum(filters.ExportOption, "exportOption", validExportOptionValues); err != nil {
			return nil, err
		}
	}
	if filters.ManagedBy != "" {
		if err := validateSingleEnum(filters.ManagedBy, "managedBy", validManagedByValues); err != nil {
			return nil, err
		}
	}

	result, err := stores.certificates.ListWithFilters(filters, nextToken, maxItems)
	if err != nil {
		return nil, err
	}

	return listResultToResponse(result), nil
}

// DeleteCertificate deletes the specified certificate.
func (s *ACMService) DeleteCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if len(cert.InUseBy) > 0 {
		return nil, NewResourceInUseError("certificate", arn)
	}

	if err := stores.certificates.Delete(arn); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeCertificate retrieves detailed information about a certificate.
func (s *ACMService) DescribeCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": certificateToDetailResponse(cert),
	}, nil
}

// ResendValidationEmail resends the domain validation email for a certificate.
func (s *ACMService) ResendValidationEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.Type == "IMPORTED" {
		return nil, NewInvalidStateException("Certificate is not in PENDING_VALIDATION state")
	}
	// AWS only allows ResendValidationEmail for certificates in
	// PENDING_VALIDATION status. ISSUED certificates have already completed
	// validation and resend is a no-op that AWS rejects.
	if cert.Status != "PENDING_VALIDATION" {
		return nil, NewInvalidStateException("Certificate is not in PENDING_VALIDATION state")
	}

	// Domain and ValidationDomain are required per Smithy model.
	domain := request.GetStringParam(params, "Domain")
	validationDomain := request.GetStringParam(params, "ValidationDomain")
	if domain == "" {
		return nil, awserrors.NewValidationException("Domain is required")
	}
	if validationDomain == "" {
		return nil, awserrors.NewValidationException("ValidationDomain is required")
	}

	// In AWS, ResendValidationEmail resends the domain validation email.
	// For edge deployment, email sending is not available, but we honour the
	// parameters by updating ValidationDomain.
	changed := false
	for _, dvo := range cert.DomainValidationOptions {
		if dvo.DomainName != domain {
			continue
		}
		dvo.ValidationDomain = validationDomain
		changed = true
	}

	if !changed {
		return nil, NewInvalidDomainValidationOptionsException(fmt.Sprintf("Domain %s not found in certificate validation options", domain))
	}

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AddTagsToCertificate adds one or more tags to a certificate.
func (s *ACMService) AddTagsToCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.acmTagConfig(stores, arn))
}

// RemoveTagsFromCertificate removes one or more tags from a certificate.
func (s *ACMService) RemoveTagsFromCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.acmTagConfig(stores, arn))
}

// ListTagsForCertificate lists the tags associated with a certificate.
func (s *ACMService) ListTagsForCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.acmTagConfig(stores, arn))
}

// TagResource adds one or more tags to an ACM resource (generic API).
func (s *ACMService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := parseCertificateArn(req.Parameters, "ResourceArn")
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.acmGenericTagConfig(stores, arn))
}

// UntagResource removes one or more tags from an ACM resource (generic API).
func (s *ACMService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := parseCertificateArn(req.Parameters, "ResourceArn")
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.acmGenericTagConfig(stores, arn))
}

// ListTagsForResource lists the tags associated with an ACM resource (generic API).
func (s *ACMService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn, err := parseCertificateArn(req.Parameters, "ResourceArn")
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.acmGenericTagConfig(stores, arn))
}

// ImportCertificate imports a certificate into ACM. When CertificateArn is
// provided, the existing certificate is updated (re-import); otherwise a new
// IMPORTED certificate is created.
func (s *ACMService) ImportCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	certificate := request.GetStringParam(params, "Certificate")
	if certificate == "" {
		return nil, awserrors.NewValidationException("Certificate is required")
	}
	certificate = decodeBase64PEM(certificate)
	if len(certificate) > 32768 {
		return nil, awserrors.NewValidationException("Certificate exceeds maximum length of 32768 bytes")
	}

	privateKey := request.GetStringParam(params, "PrivateKey")

	existingArn := request.GetStringParam(params, "CertificateArn")

	// PrivateKey is required for initial import (no existingArn).
	// For re-import (existingArn set), PrivateKey is optional — the
	// existing key is retained if not provided (Smithy: required=false).
	if existingArn == "" && privateKey == "" {
		return nil, awserrors.NewValidationException("PrivateKey is required for initial certificate import")
	}
	if privateKey != "" {
		privateKey = decodeBase64PEM(privateKey)
		if len(privateKey) > 5120 {
			return nil, awserrors.NewValidationException("PrivateKey exceeds maximum length of 5120 bytes")
		}
	}
	certificateChain := request.GetStringParam(params, "CertificateChain")
	if certificateChain != "" {
		certificateChain = decodeBase64PEM(certificateChain)
		if len(certificateChain) > 2097152 {
			return nil, awserrors.NewValidationException("CertificateChain exceeds maximum length of 2097152 bytes")
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// --- Re-import path: update an existing certificate ---
	if existingArn != "" {
		if err := validateCertificateArn(existingArn); err != nil {
			return nil, err
		}
		existing, err := stores.certificates.Get(existingArn)
		if err != nil {
			if acmstorelib.IsNotFound(err) {
				return nil, awserrors.NewResourceNotFoundException("certificate", existingArn)
			}
			return nil, err
		}

		// Only IMPORTED certificates can be re-imported.
		if existing.Type != "IMPORTED" {
			return nil, awserrors.NewInvalidParameterException("Only imported certificates can be re-imported")
		}

		// If Tags are provided on re-import, validate and replace.
		// If not provided, existing tags are retained.
		reImportTags := tagutil.ParseTagsWithQueryFallback(params, "Tags")
		if len(reImportTags) > 0 {
			if len(reImportTags) > 50 {
				return nil, NewTooManyTagsException("Tags must not exceed 50 entries")
			}
			if err := validateACMTags(reImportTags); err != nil {
				return nil, err
			}
			existing.Tags = reImportTags
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
			if len(parsedCert.DNSNames) > 0 {
				existing.DomainName = parsedCert.DNSNames[0]
			} else if parsedCert.Subject.CommonName != "" {
				existing.DomainName = parsedCert.Subject.CommonName
			}
		}

		if err := stores.certificates.Update(existing); err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"CertificateArn": existingArn,
		}, nil
	}

	// --- Initial import path: create a new certificate ---
	tags := tagutil.ParseTagsWithQueryFallback(params, "Tags")
	if err := validateACMTags(tags); err != nil {
		return nil, err
	}

	certId := acmstorelib.GenerateCertificateId()
	certificateArn := stores.arnBuilder.BuildCertificateARN(certId)

	now := time.Now().UTC()
	domainName, err := extractDomainFromCert(certificate)
	if err != nil {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Invalid certificate: %v", err))
	}
	cert := &acmstorelib.Certificate{
		CertificateArn:           certificateArn,
		DomainName:               domainName,
		Serial:                   acmstorelib.GenerateCertificateSerial(),
		Status:                   "ISSUED",
		Type:                     "IMPORTED",
		KeyAlgorithm:             determineKeyAlgorithm(certificate),
		SignatureAlgorithm:       "SHA256WITHRSA",
		RenewalEligibility:       "INELIGIBLE",
		CreatedAt:                now,
		ImportedAt:               now,
		NotBefore:                now,
		NotAfter:                 now.AddDate(1, 0, 0),
		Certificate:              certificate,
		CertificateChain:         certificateChain,
		PrivateKey:               privateKey,
		Tags:                     tags,
		AccountID:                reqCtx.GetAccountID(),
		Region:                   reqCtx.GetRegion(),
		CertificateKeyPairOrigin: "CUSTOMER_PROVIDED",
	}

	if parsedCert, _ := vcrypto.ParseCertificatePEM([]byte(certificate)); parsedCert != nil {
		cert.NotBefore = parsedCert.NotBefore
		cert.NotAfter = parsedCert.NotAfter
		cert.KeyAlgorithm = determineKeyAlgorithmFromParsed(parsedCert)
		cert.SignatureAlgorithm = determineSignatureAlgorithmFromParsed(parsedCert)
		cert.Subject = parsedCert.Subject.String()
		cert.Issuer = parsedCert.Issuer.String()
	}

	if err := stores.certificates.Create(cert); err != nil {
		if acmstorelib.IsAlreadyExists(err) {
			return nil, awserrors.NewConflictException("Certificate already exists")
		}
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certificateArn,
	}, nil
}

// UpdateCertificateOptions updates the certificate options.
func (s *ACMService) UpdateCertificateOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	// Smithy declares InvalidStateException for UpdateCertificateOptions —
	// only ISSUED certs can have their options updated.
	if cert.Status != "ISSUED" {
		return nil, NewInvalidStateException("Certificate is not in a valid state for updating options.")
	}

	optionsMap, ok := params["Options"].(map[string]interface{})
	if !ok {
		return nil, awserrors.NewValidationException("Options are required")
	}

	ctlp, err := parseCertificateTransparencyLoggingPreference(optionsMap)
	if err != nil {
		return nil, err
	}
	exportOpt, err := parseExportOption(optionsMap)
	if err != nil {
		return nil, err
	}
	cert.Options = &acmstorelib.CertificateOptions{
		CertificateTransparencyLoggingPreference: ctlp,
		Export:                                   exportOpt,
	}

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RenewCertificate renews an ACM certificate.
func (s *ACMService) RenewCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.Type != "AMAZON_ISSUED" {
		return nil, awserrors.NewValidationException("Certificate type is not supported. Only Amazon-issued certificates can be renewed.")
	}

	// Smithy declares RequestInProgressException for RenewCertificate —
	// a PENDING_VALIDATION cert is already being processed.
	if cert.Status == "PENDING_VALIDATION" {
		return nil, NewRequestInProgressException("Certificate request is in progress.")
	}

	// Only ISSUED or EXPIRED certs can be renewed. All other states
	// (FAILED, REVOKED, etc.) are invalid for renewal.
	if cert.Status != "ISSUED" && cert.Status != "EXPIRED" {
		return nil, awserrors.NewValidationException("Certificate is not in a valid state for renewal.")
	}

	if cert.RenewalEligibility == "INELIGIBLE" {
		return nil, awserrors.NewValidationException("Certificate is not eligible for renewal.")
	}

	now := time.Now().UTC()
	cert.NotBefore = now
	cert.NotAfter = now.AddDate(1, 0, 0)
	// ACM renewal transitions through PENDING_VALIDATION → SUCCESS; the
	// platform re-signs synchronously so we persist the final state directly.
	cert.IssuedAt = now
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

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ExportCertificate exports a private key and certificate.
func (s *ACMService) ExportCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	arn, err := parseCertificateArn(params, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.PrivateKey == "" {
		return nil, awserrors.NewValidationException("Certificate does not have an exportable private key. Only imported certificates with private keys can be exported.")
	}

	passphrase := request.GetStringParam(params, "Passphrase")
	if passphrase == "" {
		return nil, awserrors.NewValidationException("Passphrase is required")
	}
	// Smithy defines Passphrase as a Blob (base64-encoded). Reject if the
	// input is not valid base64 rather than silently falling back to raw
	// bytes, which would be a fail-OPEN acceptance of malformed input.
	passphraseBytes, err := base64.StdEncoding.DecodeString(passphrase)
	if err != nil {
		return nil, awserrors.NewValidationException("Passphrase must be valid base64-encoded data")
	}
	if passLen := len(passphraseBytes); passLen < 4 || passLen > 128 {
		return nil, awserrors.NewValidationException("Passphrase must be between 4 and 128 bytes")
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

	return map[string]interface{}{
		"Certificate":      cert.Certificate,
		"CertificateChain": cert.CertificateChain,
		"PrivateKey":       encryptedKey,
	}, nil
}

// RevokeCertificate revokes an ACM certificate.
func (s *ACMService) RevokeCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn, err := parseCertificateArn(req.Parameters, "CertificateArn")
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := stores.certificates.Get(arn)
	if err != nil {
		if acmstorelib.IsNotFound(err) {
			return nil, awserrors.NewResourceNotFoundException("certificate", arn)
		}
		return nil, err
	}

	if cert.Status == "REVOKED" {
		return nil, awserrors.NewResourceNotFoundException("certificate", arn)
	}

	// Smithy declares ResourceInUseException for RevokeCertificate — reject
	// when the cert is actively in use by CloudFront, API Gateway, etc.
	if len(cert.InUseBy) > 0 {
		return nil, NewResourceInUseError("certificate", arn)
	}

	if cert.Status != "ISSUED" {
		return nil, awserrors.NewValidationException("Certificate is not in a valid state for revocation.")
	}

	// RevocationReason is required per Smithy model.
	reasonRaw, ok := req.Parameters["RevocationReason"]
	if !ok {
		return nil, awserrors.NewValidationException("RevocationReason is required")
	}
	reason, ok := reasonRaw.(string)
	if !ok || !isValidRevocationReason(reason) {
		return nil, NewInvalidParameterError(fmt.Sprintf("Invalid RevocationReason: %v", reasonRaw))
	}
	cert.RevocationReason = reason
	cert.Status = "REVOKED"
	cert.RevokedAt = time.Now().UTC()

	if err := stores.certificates.Update(cert); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": arn,
	}, nil
}

// generateAmazonIssuedCert generates a self-signed AMAZON_ISSUED certificate
// with the given parameters. Shared by the HTTP API (RequestCertificate) and
// the admin handler. The caller is responsible for setting AccountID, Region,
// Tags, ManagedBy, CertificateAuthorityArn, and Options afterwards.
func generateAmazonIssuedCert(certArn, domainName string, sans []string, keyAlgorithm, validationMethod string) (*acmstorelib.Certificate, error) {
	// Build the complete DNSNames list for the certificate template.
	dnsNames := make([]string, 0, 1+len(sans))
	dnsNames = append(dnsNames, domainName)
	dnsNames = append(dnsNames, sans...)

	now := time.Now().UTC()
	key, err := generateKeyForKeyAlgorithm(keyAlgorithm)
	if err != nil {
		return nil, NewInvalidParameterError(fmt.Sprintf("Unsupported KeyAlgorithm: %s", keyAlgorithm))
	}

	serialBigInt, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return nil, NewInternalServerException("Failed to generate serial")
	}
	notAfter := now.AddDate(1, 0, 0)
	template := &x509.Certificate{
		SerialNumber: serialBigInt,
		Subject:      pkix.Name{CommonName: domainName},
		NotBefore:    now,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:     dnsNames,
	}

	certDER, err := vcrypto.CreateCertificate(template, template, key.Public(), key)
	if err != nil {
		return nil, NewInternalServerException("Failed to create certificate")
	}

	certPEM := vcrypto.EncodeCertificatePEM(certDER)

	// Build domain validation options. ACM creates certs in PENDING_VALIDATION
	// state; the platform validates synchronously (self-signed) so domain
	// validation options transition to SUCCESS immediately.
	domainValidationOptions := buildDomainValidationOptions(domainName, validationMethod, sans)
	for i := range domainValidationOptions {
		domainValidationOptions[i].ValidationStatus = "SUCCESS"
	}

	return &acmstorelib.Certificate{
		CertificateArn:           certArn,
		DomainName:               domainName,
		Serial:                   acmstorelib.GenerateCertificateSerial(),
		Status:                   "ISSUED",
		Type:                     "AMAZON_ISSUED",
		KeyAlgorithm:             keyAlgorithm,
		SignatureAlgorithm:       signatureAlgorithmForKeyAlgorithm(keyAlgorithm),
		RenewalEligibility:       "ELIGIBLE",
		CreatedAt:                now,
		Subject:                  domainName,
		Issuer:                   domainName,
		Certificate:              certPEM,
		NotBefore:                now,
		NotAfter:                 notAfter,
		IssuedAt:                 now,
		SubjectAlternativeNames:  sans,
		DomainValidationOptions:  domainValidationOptions,
		CertificateKeyPairOrigin: "AWS_MANAGED",
	}, nil
}

func buildDomainValidationOptions(domainName, validationMethod string, sans []string) []*acmstorelib.DomainValidation {
	// Build validation entries for the primary domain and all SANs.
	allDomains := make([]string, 0, 1+len(sans))
	allDomains = append(allDomains, domainName)
	allDomains = append(allDomains, sans...)

	options := make([]*acmstorelib.DomainValidation, 0, len(allDomains))
	for _, d := range allDomains {
		dv := &acmstorelib.DomainValidation{
			DomainName:       d,
			ValidationDomain: d,
			ValidationMethod: validationMethod,
			ValidationStatus: "PENDING_VALIDATION",
		}

		if validationMethod == "DNS" {
			dv.ResourceRecord = &acmstorelib.ResourceRecord{
				Name:  acmstorelib.GenerateDomainValidationRecordName(d),
				Type:  "CNAME",
				Value: acmstorelib.GenerateDomainValidationRecordValue(),
			}
		}

		options = append(options, dv)
	}

	return options
}

// applyUserDomainValidationOptions merges user-provided DomainValidationOptions
// (e.g. custom ValidationDomain for EMAIL validation) into the auto-generated
// options, matching by DomainName. Both DomainName and ValidationDomain are
// REQUIRED per the Smithy DomainValidationOption shape — entries missing
// either field are rejected.
func applyUserDomainValidationOptions(options []*acmstorelib.DomainValidation, params map[string]interface{}) error {
	raw, ok := params["DomainValidationOptions"]
	if !ok {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return awserrors.NewValidationException("DomainValidationOptions must be a list")
	}
	// Smithy DomainValidationOptionList: @length(min=1, max=100).
	if len(arr) == 0 {
		return awserrors.NewValidationException("DomainValidationOptions must contain at least 1 entry when provided")
	}
	if len(arr) > 100 {
		return awserrors.NewValidationException("DomainValidationOptions must not exceed 100 entries")
	}
	// Build a lookup of user-provided ValidationDomain by DomainName.
	userMap := make(map[string]string)
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		dn, _ := m["DomainName"].(string)
		vd, _ := m["ValidationDomain"].(string)
		if err := validateDomainValidationFields(dn, vd); err != nil {
			return err
		}
		userMap[strings.ToLower(dn)] = strings.ToLower(vd)
	}
	// Apply user-provided ValidationDomain to matching entries.
	for _, dv := range options {
		if vd, ok := userMap[strings.ToLower(dv.DomainName)]; ok {
			dv.ValidationDomain = vd
		}
	}
	return nil
}

// generateKeyForKeyAlgorithm generates a private key matching the ACM
// KeyAlgorithm enum value. Returns a crypto.Signer that works with both
// RSA and ECDSA key types for x509 certificate creation.
func generateKeyForKeyAlgorithm(keyAlgorithm string) (crypto.Signer, error) {
	switch keyAlgorithm {
	case "", "RSA_2048":
		return vcrypto.GenerateRSAKey(2048)
	case "RSA_3072":
		return vcrypto.GenerateRSAKey(3072)
	case "RSA_4096":
		return vcrypto.GenerateRSAKey(4096)
	case "EC_prime256v1":
		return vcrypto.GenerateECDSAKey(elliptic.P256())
	case "EC_secp384r1":
		return vcrypto.GenerateECDSAKey(elliptic.P384())
	case "EC_secp521r1":
		return vcrypto.GenerateECDSAKey(elliptic.P521())
	default:
		return nil, fmt.Errorf("unsupported key algorithm: %s", keyAlgorithm)
	}
}

// signatureAlgorithmForKeyAlgorithm returns the ACM SignatureAlgorithm
// string that corresponds to the given KeyAlgorithm.
func signatureAlgorithmForKeyAlgorithm(keyAlgorithm string) string {
	switch keyAlgorithm {
	case "EC_prime256v1":
		return "SHA256WITHECDSA"
	case "EC_secp384r1":
		return "SHA384WITHECDSA"
	case "EC_secp521r1":
		return "SHA512WITHECDSA"
	default:
		return "SHA256WITHRSA"
	}
}

func extractDomainFromCert(cert string) (string, error) {
	parsed, err := vcrypto.ParseCertificatePEM([]byte(cert))
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	if len(parsed.DNSNames) > 0 {
		return parsed.DNSNames[0], nil
	}

	if parsed.Subject.CommonName != "" {
		return parsed.Subject.CommonName, nil
	}

	return "", fmt.Errorf("no domain name found in certificate")
}

func determineKeyAlgorithm(cert string) string {
	parsed, err := vcrypto.ParseCertificatePEM([]byte(cert))
	if err != nil {
		return "RSA_2048"
	}
	return determineKeyAlgorithmFromParsed(parsed)
}

func determineKeyAlgorithmFromParsed(cert *x509.Certificate) string {
	switch cert.PublicKeyAlgorithm {
	case x509.ECDSA:
		if key, ok := cert.PublicKey.(*ecdsa.PublicKey); ok {
			switch key.Curve {
			case elliptic.P384():
				return "EC_secp384r1"
			case elliptic.P521():
				return "EC_secp521r1"
			}
		}
		return "EC_prime256v1"
	case x509.RSA:
		bits := 0
		if key, ok := cert.PublicKey.(*rsa.PublicKey); ok {
			bits = key.N.BitLen()
		}
		switch bits {
		case 1024:
			return "RSA_1024"
		case 2048:
			return "RSA_2048"
		case 3072:
			return "RSA_3072"
		case 4096:
			return "RSA_4096"
		default:
			return "RSA_2048"
		}
	default:
		return "RSA_2048"
	}
}

func determineSignatureAlgorithmFromParsed(cert *x509.Certificate) string {
	switch cert.SignatureAlgorithm {
	case x509.SHA256WithRSA:
		return "SHA256WITHRSA"
	case x509.SHA384WithRSA:
		return "SHA384WITHRSA"
	case x509.SHA512WithRSA:
		return "SHA512WITHRSA"
	case x509.ECDSAWithSHA256:
		return "SHA256WITHECDSA"
	case x509.ECDSAWithSHA384:
		return "SHA384WITHECDSA"
	case x509.ECDSAWithSHA512:
		return "SHA512WITHECDSA"
	default:
		return "SHA256WITHRSA"
	}
}

func encryptPrivateKey(privateKeyPEM, passphrase string) (string, error) {
	keyBytes := []byte(privateKeyPEM)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		if decoded, err := base64.StdEncoding.DecodeString(privateKeyPEM); err == nil {
			block, _ = pem.Decode(decoded)
		}
	}
	if block == nil {
		if restored := pemFixNewlines(privateKeyPEM); restored != privateKeyPEM {
			block, _ = pem.Decode([]byte(restored))
		}
	}
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(passphrase), salt, 100000, 32, sha256.New)
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nil, nonce, block.Bytes, nil)
	result := append(salt, nonce...)
	result = append(result, encrypted...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func decodeBase64PEM(s string) string {
	if strings.Contains(s, "-----BEGIN") {
		return s
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}
	return strings.TrimSpace(string(decoded))
}

func pemFixNewlines(pemStr string) string {
	pemStr = strings.TrimSpace(pemStr)
	if strings.Contains(pemStr, "\n") {
		return pemStr
	}
	begin := strings.Index(pemStr, "-----BEGIN")
	if begin < 0 {
		return pemStr
	}
	end := strings.Index(pemStr, "-----END")
	if end < 0 {
		return pemStr
	}
	end += len("-----END")
	typeEnd := strings.Index(pemStr[begin:], "-----")
	if typeEnd < 0 {
		return pemStr
	}
	typeEnd += begin
	var b strings.Builder
	b.WriteString(pemStr[begin : typeEnd+5])
	b.WriteByte('\n')
	b64 := pemStr[typeEnd+5 : end]
	for i := 0; i < len(b64); i += 64 {
		if i+64 > len(b64) {
			b.WriteString(b64[i:])
		} else {
			b.WriteString(b64[i : i+64])
			b.WriteByte('\n')
		}
	}
	b.WriteString(pemStr[end:])
	return b.String()
}
