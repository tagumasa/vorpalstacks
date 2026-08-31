package acm

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

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
	cert, err := s.getCertificateCore(reqCtx, request.GetStringParam(req.Parameters, "CertificateArn"))
	if err != nil {
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
	cert, err := s.getCertificateCore(reqCtx, request.GetStringParam(req.Parameters, "CertificateArn"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": certificateToDetailResponse(cert),
	}, nil
}

// ResendValidationEmail resends the domain validation email for a certificate.
func (s *ACMService) ResendValidationEmail(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	err := s.resendValidationEmailCore(reqCtx, ResendValidationEmailInput{
		CertificateArn:   request.GetStringParam(req.Parameters, "CertificateArn"),
		Domain:           request.GetStringParam(req.Parameters, "Domain"),
		ValidationDomain: request.GetStringParam(req.Parameters, "ValidationDomain"),
	})
	if err != nil {
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
	certArn, err := s.importCertificateCore(reqCtx, ImportCertificateInput{
		Certificate:      request.GetStringParam(req.Parameters, "Certificate"),
		PrivateKey:       request.GetStringParam(req.Parameters, "PrivateKey"),
		CertificateChain: request.GetStringParam(req.Parameters, "CertificateChain"),
		ExistingArn:      request.GetStringParam(req.Parameters, "CertificateArn"),
		Tags:             tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certArn,
	}, nil
}

// UpdateCertificateOptions updates the certificate options.
func (s *ACMService) UpdateCertificateOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	err := s.updateCertificateOptionsCore(reqCtx, UpdateCertificateOptionsInput{
		CertificateArn: request.GetStringParam(req.Parameters, "CertificateArn"),
		OptionsRaw:     req.Parameters["Options"],
	})
	if err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// RenewCertificate renews an ACM certificate.
func (s *ACMService) RenewCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.renewCertificateCore(reqCtx, request.GetStringParam(req.Parameters, "CertificateArn")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ExportCertificate exports a private key and certificate.
func (s *ACMService) ExportCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.exportCertificateCore(reqCtx, ExportCertificateInput{
		CertificateArn: request.GetStringParam(req.Parameters, "CertificateArn"),
		Passphrase:     request.GetStringParam(req.Parameters, "Passphrase"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Certificate":      result.Certificate,
		"CertificateChain": result.CertificateChain,
		"PrivateKey":       result.PrivateKey,
	}, nil
}

// RevokeCertificate revokes an ACM certificate.
func (s *ACMService) RevokeCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certArn, err := s.revokeCertificateCore(reqCtx, RevokeCertificateInput{
		CertificateArn:      request.GetStringParam(req.Parameters, "CertificateArn"),
		RevocationReasonRaw: req.Parameters["RevocationReason"],
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CertificateArn": certArn,
	}, nil
}
