package acm

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	acmstorelib "vorpalstacks/internal/store/aws/acm"

	vcrypto "vorpalstacks/internal/utils/crypto"
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
		TagFunc: func(ctx context.Context, resourceKey string, tagList []tagutil.Tag) error {
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
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			cert, err := stores.certificates.Get(resourceKey)
			if err != nil {
				return nil, err
			}
			return cert.Tags, nil
		},
		FormatResponse: func(tagList []tagutil.Tag, _ string) (interface{}, error) {
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

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := RequestCertificateInput{
		DomainName:       request.GetStringParam(params, "DomainName"),
		IdempotencyToken: request.GetStringParam(params, "IdempotencyToken"),
		KeyAlgorithm:     request.GetStringParam(params, "KeyAlgorithm"),
		ValidationMethod: request.GetStringParam(params, "ValidationMethod"),
		ManagedBy:        request.GetStringParam(params, "ManagedBy"),
		PCAArn:           request.GetStringParam(params, "CertificateAuthorityArn"),
		AccountID:        reqCtx.GetAccountID(),
		Region:           reqCtx.GetRegion(),
	}

	// SubjectAlternativeNames: detect presence for empty-list validation.
	if sansRaw, ok := params["SubjectAlternativeNames"]; ok && sansRaw != nil {
		input.SANsProvided = true
		if sansArr, ok := sansRaw.([]interface{}); ok {
			for _, san := range sansArr {
				if s, ok := san.(string); ok {
					input.SANs = append(input.SANs, s)
				}
			}
		}
	}

	// Tags: detect presence for empty-list validation.
	if _, ok := params["Tags"]; ok {
		input.TagsProvided = true
	}
	input.Tags = tagutil.ParseTagsWithQueryFallback(params, "Tags")

	// DomainValidationOptions: detect presence + parse overrides.
	if rawDVOs, ok := params["DomainValidationOptions"]; ok && rawDVOs != nil {
		input.DVOsProvided = true
		if arr, ok := rawDVOs.([]interface{}); ok {
			for _, item := range arr {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				dn, _ := m["DomainName"].(string)
				vd, _ := m["ValidationDomain"].(string)
				input.DVOs = append(input.DVOs, DomainValidationOverride{
					DomainName:       dn,
					ValidationDomain: vd,
				})
			}
		}
	}

	// Options.
	if optionsMap, ok := params["Options"].(map[string]interface{}); ok {
		input.Options = &CertificateOptionsInput{
			CTLPreference: request.GetStringParam(optionsMap, "CertificateTransparencyLoggingPreference"),
			Export:        request.GetStringParam(optionsMap, "Export"),
		}
	}

	certArn, err := s.requestCertificateCore(ctx, stores, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certArn,
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
	if err := validateNextToken(nextToken); err != nil {
		return nil, err
	}

	// Smithy MaxItems: @range(1-1000). When explicitly provided, validate
	// the range; when absent, default to pagination.DefaultMaxItems.
	maxItems := pagination.DefaultMaxItems
	if _, ok := params["MaxItems"]; ok {
		mi := request.GetIntParam(params, "MaxItems")
		if mi < 1 || mi > 1000 {
			return nil, awserrors.NewValidationException("MaxItems must be between 1 and 1000")
		}
		maxItems = mi
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := ListCertificatesInput{
		NextToken: nextToken,
		MaxItems:  maxItems,
		SortBy:    request.GetStringParam(params, "SortBy"),
		SortOrder: request.GetStringParam(params, "SortOrder"),
	}

	// Smithy SortBy enum for ListCertificates only permits CREATED_AT.
	if input.SortBy != "" && input.SortBy != "CREATED_AT" {
		return nil, NewInvalidParameterError(fmt.Sprintf("Invalid SortBy: %s. Valid values are CREATED_AT.", input.SortBy))
	}
	if input.SortOrder != "" {
		if err := validateSortOrder(input.SortOrder); err != nil {
			return nil, err
		}
	}

	input.Statuses = request.GetStringList(params, "CertificateStatuses")
	if err := validateEnumList(input.Statuses, "CertificateStatuses", validCertificateStatuses); err != nil {
		return nil, err
	}

	input.Origins = request.GetStringList(params, "CertificateKeyPairOrigins")
	// Smithy CertificateKeyPairOrigins: @length(1-3) when provided.
	if _, ok := params["CertificateKeyPairOrigins"]; ok {
		if len(input.Origins) == 0 {
			return nil, NewInvalidParameterError("CertificateKeyPairOrigins must contain at least 1 entry when provided")
		}
		if len(input.Origins) > 3 {
			return nil, NewInvalidParameterError("CertificateKeyPairOrigins must not exceed 3 entries")
		}
	}
	if err := validateEnumList(input.Origins, "CertificateKeyPairOrigins", validKeyPairOrigins); err != nil {
		return nil, err
	}

	if includes, ok := params["Includes"].(map[string]interface{}); ok {
		input.KeyTypes = request.GetStringList(includes, "keyTypes")
		input.KeyUsage = request.GetStringList(includes, "keyUsage")
		input.ExtendedKeyUsage = request.GetStringList(includes, "extendedKeyUsage")
		input.ExportOption = request.GetStringParam(includes, "exportOption")
		input.ManagedBy = request.GetStringParam(includes, "managedBy")
	}
	if err := validateEnumList(input.KeyTypes, "keyTypes", validKeyAlgorithmsMap); err != nil {
		return nil, err
	}
	if err := validateEnumList(input.KeyUsage, "keyUsage", validKeyUsageNames); err != nil {
		return nil, err
	}
	if err := validateEnumList(input.ExtendedKeyUsage, "extendedKeyUsage", validExtendedKeyUsageNames); err != nil {
		return nil, err
	}
	// Validate single-value enum filters that are not lists.
	if input.ExportOption != "" {
		if err := validateSingleEnum(input.ExportOption, "exportOption", validExportOptionValues); err != nil {
			return nil, err
		}
	}
	if input.ManagedBy != "" {
		if err := validateSingleEnum(input.ManagedBy, "managedBy", validManagedByValues); err != nil {
			return nil, err
		}
	}

	result, err := s.listCertificatesCore(stores, input)
	if err != nil {
		return nil, err
	}

	return listResultToResponse(result), nil
}

// DeleteCertificate deletes the specified certificate.
func (s *ACMService) DeleteCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	arn := request.GetStringParam(req.Parameters, "CertificateArn")
	if err := s.deleteCertificateCore(stores, arn); err != nil {
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
	certificate, err := decodeBase64PEM(certificate)
	if err != nil {
		return nil, err
	}
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
		privateKey, err = decodeBase64PEM(privateKey)
		if err != nil {
			return nil, err
		}
		if len(privateKey) > 5120 {
			return nil, awserrors.NewValidationException("PrivateKey exceeds maximum length of 5120 bytes")
		}
	}
	certificateChain := request.GetStringParam(params, "CertificateChain")
	if certificateChain != "" {
		certificateChain, err = decodeBase64PEM(certificateChain)
		if err != nil {
			return nil, err
		}
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
		Serial:                   extractSerialFromPEM(certificate),
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

	// Regenerate the X.509 certificate material (new key pair, new serial,
	// new PEM). The ARN is preserved per AWS specification: "When ACM renews
	// a certificate, the certificate's Amazon Resource Name (ARN) remains
	// the same."
	if err := renewCertificateMaterial(cert); err != nil {
		return nil, NewInternalServerException(fmt.Sprintf("Failed to regenerate certificate material: %v", err))
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
