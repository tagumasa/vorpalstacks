// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/timeutils"
)

// effectiveCertificateArn returns the ACM certificate ARN that a domain name
// is configured to use, preferring the regional certificate over the edge
// certificate. Returns empty string if no ACM certificate is configured.
func effectiveCertificateArn(d *store.DomainName) string {
	if d.RegionalCertificateArn != "" {
		return d.RegionalCertificateArn
	}
	return d.CertificateArn
}

// resolveDomainName resolves a domain name from either the domainName or
// domainNameId request parameter. At least one must be provided per the Smithy
// model. When only domainNameId is given, the store is queried to obtain the
// canonical domain name.
func resolveDomainName(req *request.ParsedRequest, stores *apiGatewayStores) (string, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		domainName = getPathParam(req, "domainName")
	}
	if domainName != "" {
		return domainName, nil
	}

	domainNameId := request.GetStringParam(req.Parameters, "domainNameId")
	if domainNameId == "" {
		return "", NewBadRequestException("Either domainName or domainNameId must be specified")
	}

	domain, err := stores.domains.GetDomainNameById(domainNameId)
	if err != nil {
		return "", ErrNotFoundException
	}
	return domain.DomainName, nil
}

// CreateDomainName creates a new domain name for API Gateway.
func (s *APIGatewayService) CreateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	domainName := request.GetStringParam(req.Parameters, "domainName")
	if domainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}
	if !validateFQDN(domainName) {
		return nil, NewBadRequestException("Invalid domain name: must be a valid FQDN")
	}

	securityPolicy := request.GetStringParam(req.Parameters, "securityPolicy")
	if !validateSecurityPolicy(securityPolicy) {
		return nil, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
	}
	endpointAccessMode := request.GetStringParam(req.Parameters, "endpointAccessMode")
	if !validateEndpointAccessMode(endpointAccessMode) {
		return nil, NewBadRequestException("Invalid endpointAccessMode: must be BASIC or STRICT")
	}

	domain := &store.DomainName{
		DomainName:                          domainName,
		CertificateArn:                      request.GetStringParam(req.Parameters, "certificateArn"),
		CertificateName:                     request.GetStringParam(req.Parameters, "certificateName"),
		RegionalCertificateArn:              request.GetStringParam(req.Parameters, "regionalCertificateArn"),
		RegionalCertificateName:             request.GetStringParam(req.Parameters, "regionalCertificateName"),
		SecurityPolicy:                      request.GetStringParam(req.Parameters, "securityPolicy"),
		OwnershipVerificationCertificateArn: request.GetStringParam(req.Parameters, "ownershipVerificationCertificateArn"),
		EndpointAccessMode:                  request.GetStringParam(req.Parameters, "endpointAccessMode"),
		Policy:                              request.GetStringParam(req.Parameters, "policy"),
		RoutingMode:                         request.GetStringParam(req.Parameters, "routingMode"),
	}

	if mutualTls, ok := req.Parameters["mutualTlsAuthentication"].(map[string]interface{}); ok {
		domain.MutualTlsAuthentication = &store.MutualTlsAuthentication{}
		if v, ok := mutualTls["truststoreUri"].(string); ok {
			domain.MutualTlsAuthentication.TruststoreUri = v
		}
		if v, ok := mutualTls["truststoreVersion"].(string); ok {
			domain.MutualTlsAuthentication.TruststoreVersion = v
		}
	}

	if endpointConfig, ok := req.Parameters["endpointConfiguration"].(map[string]interface{}); ok {
		domain.EndpointConfiguration = &store.EndpointConfiguration{}
		if types, ok := endpointConfig["types"].([]interface{}); ok {
			for _, t := range types {
				if ts, ok := t.(string); ok {
					domain.EndpointConfiguration.Types = append(domain.EndpointConfiguration.Types, ts)
				}
			}
		}
	}

	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		domain.Tags = tagutil.MapInterfaceToTags(tags)
	}

	// Require an ACM certificate ARN — vorpalstacks only supports the ACM
	// certificate ARN model and does not accept the certificateBody /
	// certificatePrivateKey / certificateChain upload path.  Rejecting
	// explicitly instead of silently dropping prevents a fail-OPEN gap where
	// a caller could create a domain with no certificate at all.
	certArn := effectiveCertificateArn(domain)
	if certArn == "" {
		if request.GetStringParam(req.Parameters, "certificateBody") != "" ||
			request.GetStringParam(req.Parameters, "certificatePrivateKey") != "" ||
			request.GetStringParam(req.Parameters, "certificateChain") != "" {
			return nil, NewBadRequestException("certificateBody, certificatePrivateKey, and certificateChain are not supported; use certificateArn or regionalCertificateArn")
		}
		return nil, NewBadRequestException("certificateArn or regionalCertificateArn is required")
	}
	if s.acmInvoker != nil && !s.acmInvoker.CertificateExists(ctx, reqCtx.GetRegion(), certArn) {
		return nil, NewBadRequestException("The specified certificate ARN does not exist: " + certArn)
	}

	if !validateRoutingMode(domain.RoutingMode) {
		return nil, NewBadRequestException("Invalid routingMode: " + domain.RoutingMode)
	}

	if !validatePolicyJSON(domain.Policy) {
		return nil, NewBadRequestException("Invalid policy: must be valid JSON")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.domains.CreateDomainName(domain)
	if err != nil {
		return nil, err
	}

	if certArn != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.RegisterCertificateUsage(ctx, reqCtx.GetRegion(), certArn, created.DomainNameArn); err != nil {
			_ = stores.domains.DeleteDomainName(created.DomainName)
			return nil, fmt.Errorf("failed to register certificate usage: %w", err)
		}
	}

	return s.toDomainNameResponse(created), nil
}

// GetDomainName retrieves a domain name by its name.
func (s *APIGatewayService) GetDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toDomainNameResponse(domain), nil
}

// DeleteDomainName deletes a domain name.
func (s *APIGatewayService) DeleteDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	// AWS requires all BasePathMappings to be deleted before a domain name
	// can be removed; otherwise it returns ConflictException. The listing is
	// fail-closed: when the store cannot answer whether mappings exist, the
	// deletion is refused rather than proceeding on an unknown state.
	mappings, err := stores.domains.ListBasePathMappings(domainName, common.ListOptions{MaxItems: 1})
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	if len(mappings.Items) > 0 {
		return nil, NewConflictException("Domain name has active base path mappings; remove them first")
	}

	// Capture cert ARN before deletion for best-effort unregister after.
	certArnForCleanup := effectiveCertificateArn(domain)
	domainArn := domain.DomainNameArn

	if err := stores.domains.DeleteDomainName(domainName); err != nil {
		return nil, ErrNotFoundException
	}

	// Best-effort ACM unregister after successful deletion. A stale InUseBy
	// reference is harmless (only over-blocks certificate deletion), whereas
	// unregistering before deletion risks leaving a live resource unprotected.
	if certArnForCleanup != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.UnregisterCertificateUsage(ctx, reqCtx.GetRegion(), certArnForCleanup, domainArn); err != nil {
			log.Printf("warning: failed to unregister certificate usage for deleted domain %s: %v", domainName, err)
		}
	}

	return response.EmptyResponse(), nil
}

// UpdateDomainName updates an existing domain name.
func (s *APIGatewayService) UpdateDomainName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return nil, ErrNotFoundException
	}

	// Capture old state for compensation.
	oldCertArn := effectiveCertificateArn(domain)
	oldCertArnValue := domain.CertificateArn
	oldRegionalCertArnValue := domain.RegionalCertificateArn
	oldCertificateName := domain.CertificateName

	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	for _, po := range ops {
		switch {
		case po.Path == "/certificateArn":
			if po.Value == "" {
				return nil, NewBadRequestException("certificateArn cannot be cleared; provide a new certificate ARN")
			}
			domain.CertificateArn = po.Value
		case po.Path == "/regionalCertificateArn":
			if po.Value == "" {
				return nil, NewBadRequestException("regionalCertificateArn cannot be cleared; provide a new certificate ARN")
			}
			domain.RegionalCertificateArn = po.Value
		case po.Path == "/certificateName":
			domain.CertificateName = po.Value
		case po.Path == "/securityPolicy":
			if !validateSecurityPolicy(po.Value) {
				return nil, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
			}
			domain.SecurityPolicy = po.Value
		case po.Path == "/ownershipVerificationCertificateArn":
			domain.OwnershipVerificationCertificateArn = po.Value
		case po.Path == "/routingMode":
			if !validateRoutingMode(po.Value) {
				return nil, NewBadRequestException("Invalid routingMode: " + po.Value)
			}
			domain.RoutingMode = po.Value
		case po.Path == "/policy":
			if !validatePolicyJSON(po.Value) {
				return nil, NewBadRequestException("Invalid policy: must be valid JSON")
			}
			domain.Policy = po.Value
		case po.Path == "/managementPolicy":
			if !validatePolicyJSON(po.Value) {
				return nil, NewBadRequestException("Invalid managementPolicy: must be valid JSON")
			}
			domain.ManagementPolicy = po.Value
		case strings.HasPrefix(po.Path, "/endpointConfiguration/types"):
			if domain.EndpointConfiguration == nil {
				domain.EndpointConfiguration = &store.EndpointConfiguration{}
			}
			typeName := strings.TrimPrefix(po.Path, "/endpointConfiguration/types/")
			if typeName == "" || typeName == "/endpointConfiguration/types" {
				typeName = po.Value
			}
			if po.Op == "remove" {
				domain.EndpointConfiguration.Types = removeString(domain.EndpointConfiguration.Types, typeName)
			} else {
				if !sliceContains(domain.EndpointConfiguration.Types, typeName) {
					domain.EndpointConfiguration.Types = append(domain.EndpointConfiguration.Types, typeName)
				}
			}
		case po.Path == "/mutualTlsAuthentication/truststoreUri":
			if domain.MutualTlsAuthentication == nil {
				domain.MutualTlsAuthentication = &store.MutualTlsAuthentication{}
			}
			if po.Op == "remove" {
				domain.MutualTlsAuthentication.TruststoreUri = ""
			} else {
				domain.MutualTlsAuthentication.TruststoreUri = po.Value
			}
		case po.Path == "/mutualTlsAuthentication/truststoreVersion":
			if domain.MutualTlsAuthentication == nil {
				domain.MutualTlsAuthentication = &store.MutualTlsAuthentication{}
			}
			if po.Op == "remove" {
				domain.MutualTlsAuthentication.TruststoreVersion = ""
			} else {
				domain.MutualTlsAuthentication.TruststoreVersion = po.Value
			}
		}
	}

	// Pre-validate new cert ARN.
	newCertArn := effectiveCertificateArn(domain)
	if newCertArn != "" && s.acmInvoker != nil && newCertArn != oldCertArn {
		if !s.acmInvoker.CertificateExists(ctx, reqCtx.GetRegion(), newCertArn) {
			return nil, NewBadRequestException("The specified certificate ARN does not exist: " + newCertArn)
		}
	}

	if err := stores.domains.UpdateDomainName(domain); err != nil {
		return nil, err
	}

	// ACM cert operations with compensating transaction on failure.
	if s.acmInvoker != nil && oldCertArn != newCertArn {
		// Step 1: Unregister old cert.
		if oldCertArn != "" {
			if err := s.acmInvoker.UnregisterCertificateUsage(ctx, reqCtx.GetRegion(), oldCertArn, domain.DomainNameArn); err != nil {
				// Compensate: revert to old cert values.
				domain.CertificateArn = oldCertArnValue
				domain.RegionalCertificateArn = oldRegionalCertArnValue
				domain.CertificateName = oldCertificateName
				if revertErr := stores.domains.UpdateDomainName(domain); revertErr != nil {
					log.Printf("error: failed to revert domain name after unregister failure: %v", revertErr)
				}
				return nil, fmt.Errorf("failed to unregister old certificate usage: %w", err)
			}
		}
		// Step 2: Register new cert.
		if newCertArn != "" {
			if err := s.acmInvoker.RegisterCertificateUsage(ctx, reqCtx.GetRegion(), newCertArn, domain.DomainNameArn); err != nil {
				// Compensate: re-register old cert (was unregistered in step 1).
				if oldCertArn != "" {
					if revertErr := s.acmInvoker.RegisterCertificateUsage(ctx, reqCtx.GetRegion(), oldCertArn, domain.DomainNameArn); revertErr != nil {
						log.Printf("error: failed to re-register old certificate during compensation: %v", revertErr)
					}
				}
				// Revert to old cert values.
				domain.CertificateArn = oldCertArnValue
				domain.RegionalCertificateArn = oldRegionalCertArnValue
				domain.CertificateName = oldCertificateName
				if revertErr := stores.domains.UpdateDomainName(domain); revertErr != nil {
					log.Printf("error: failed to revert domain name after register failure: %v", revertErr)
				}
				return nil, fmt.Errorf("failed to register new certificate usage: %w", err)
			}
		}
	}

	return s.toDomainNameResponse(domain), nil
}

// GetDomainNames lists all domain names.
func (s *APIGatewayService) GetDomainNames(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.domains.ListDomainNames(common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, d := range result.Items {
		items = append(items, s.toDomainNameResponse(d))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toDomainNameResponse(d *store.DomainName) map[string]interface{} {
	response := map[string]interface{}{
		"domainName": d.DomainName,
	}

	if d.DomainNameId != "" {
		response["domainNameId"] = d.DomainNameId
	}
	if d.DomainNameArn != "" {
		response["domainNameArn"] = d.DomainNameArn
	}
	if d.CertificateArn != "" {
		response["certificateArn"] = d.CertificateArn
	}
	if d.CertificateName != "" {
		response["certificateName"] = d.CertificateName
	}
	if !d.CertificateUploadDate.IsZero() {
		response["certificateUploadDate"] = timeutils.FormatEpochSeconds(d.CertificateUploadDate)
	}
	if d.DistributionDomainName != "" {
		response["distributionDomainName"] = d.DistributionDomainName
	}
	if d.DistributionHostedZoneId != "" {
		response["distributionHostedZoneId"] = d.DistributionHostedZoneId
	}
	if d.RegionalDomainName != "" {
		response["regionalDomainName"] = d.RegionalDomainName
	}
	if d.RegionalHostedZoneId != "" {
		response["regionalHostedZoneId"] = d.RegionalHostedZoneId
	}
	if d.RegionalCertificateArn != "" {
		response["regionalCertificateArn"] = d.RegionalCertificateArn
	}
	if d.RegionalCertificateName != "" {
		response["regionalCertificateName"] = d.RegionalCertificateName
	}
	if d.DomainNameStatus != "" {
		response["domainNameStatus"] = d.DomainNameStatus
	}
	if d.DomainNameStatusMessage != "" {
		response["domainNameStatusMessage"] = d.DomainNameStatusMessage
	}
	if d.SecurityPolicy != "" {
		response["securityPolicy"] = d.SecurityPolicy
	}
	if d.EndpointAccessMode != "" {
		response["endpointAccessMode"] = d.EndpointAccessMode
	}
	if d.MutualTlsAuthentication != nil {
		mtls := map[string]interface{}{}
		if d.MutualTlsAuthentication.TruststoreUri != "" {
			mtls["truststoreUri"] = d.MutualTlsAuthentication.TruststoreUri
		}
		if d.MutualTlsAuthentication.TruststoreVersion != "" {
			mtls["truststoreVersion"] = d.MutualTlsAuthentication.TruststoreVersion
		}
		if len(d.MutualTlsAuthentication.TruststoreWarnings) > 0 {
			mtls["truststoreWarnings"] = d.MutualTlsAuthentication.TruststoreWarnings
		}
		response["mutualTlsAuthentication"] = mtls
	}
	if d.OwnershipVerificationCertificateArn != "" {
		response["ownershipVerificationCertificateArn"] = d.OwnershipVerificationCertificateArn
	}
	if d.Policy != "" {
		response["policy"] = d.Policy
	}
	if d.RoutingMode != "" {
		response["routingMode"] = d.RoutingMode
	}
	if d.ManagementPolicy != "" {
		response["managementPolicy"] = d.ManagementPolicy
	}
	if d.EndpointConfiguration != nil {
		response["endpointConfiguration"] = map[string]interface{}{
			"types": d.EndpointConfiguration.Types,
		}
	}
	if len(d.Tags) > 0 {
		response["tags"] = tagutil.ToMap(d.Tags)
	}

	return response
}

// CreateBasePathMapping creates a new base path mapping.
func (s *APIGatewayService) CreateBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	restApiId := request.GetStringParam(req.Parameters, "restApiId")
	if restApiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	mapping := &store.BasePathMapping{
		BasePath:  request.GetStringParam(req.Parameters, "basePath"),
		RestApiId: restApiId,
		Stage:     request.GetStringParam(req.Parameters, "stage"),
	}

	if mapping.BasePath == "" {
		mapping.BasePath = "(none)"
	}
	if !validateBasePath(mapping.BasePath) {
		return nil, NewBadRequestException("basePath must contain only alphanumeric characters, hyphens, underscores, periods, and forward slashes")
	}

	created, err := stores.domains.CreateBasePathMapping(domainName, mapping)
	if err != nil {
		return nil, err
	}

	return s.toBasePathMappingResponse(created), nil
}

// GetBasePathMapping retrieves a base path mapping.
func (s *APIGatewayService) GetBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	mapping, err := stores.domains.GetBasePathMapping(domainName, basePath)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toBasePathMappingResponse(mapping), nil
}

// DeleteBasePathMapping deletes a base path mapping.
func (s *APIGatewayService) DeleteBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	if err := stores.domains.DeleteBasePathMapping(domainName, basePath); err != nil {
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// UpdateBasePathMapping updates an existing base path mapping.
func (s *APIGatewayService) UpdateBasePathMapping(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	basePath := request.GetStringParam(req.Parameters, "basePath")
	if basePath == "" {
		basePath = getPathParam(req, "basePath")
	}
	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	stores.keyLocker.Lock(domainName + ":" + basePath)
	defer stores.keyLocker.Unlock(domainName + ":" + basePath)

	mapping, err := stores.domains.GetBasePathMapping(domainName, basePath)
	if err != nil {
		return nil, ErrNotFoundException
	}

	renamed := false
	oldBasePath := ""
	oldRestApiId := mapping.RestApiId
	oldStage := mapping.Stage
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	for _, po := range ops {
		switch po.Path {
		case "/restApiId":
			mapping.RestApiId = po.Value
		case "/stage":
			mapping.Stage = po.Value
		case "/basePath":
			if !validateBasePath(po.Value) {
				return nil, NewBadRequestException("basePath must contain only alphanumeric characters, hyphens, underscores, periods, and forward slashes")
			}
			oldBasePath = basePath
			basePath = po.Value
			mapping.BasePath = po.Value
			renamed = true
		}
	}

	if renamed {
		// Pre-check: reject if the target basePath already exists, avoiding
		// a destructive delete-then-fail-then-restore cycle.
		if _, err := stores.domains.GetBasePathMapping(domainName, basePath); err == nil {
			return nil, NewConflictException(fmt.Sprintf("basePath '%s' already exists for this domain", basePath))
		}
		if err := stores.domains.DeleteBasePathMapping(domainName, oldBasePath); err != nil {
			return nil, err
		}
		if _, err := stores.domains.CreateBasePathMapping(domainName, mapping); err != nil {
			// Compensating restore: re-create the original mapping so the
			// rename failure does not result in data loss. Restore all
			// fields that may have been modified by earlier patch ops.
			mapping.BasePath = oldBasePath
			mapping.RestApiId = oldRestApiId
			mapping.Stage = oldStage
			_, _ = stores.domains.CreateBasePathMapping(domainName, mapping)
			return nil, err
		}
	} else {
		if err := stores.domains.UpdateBasePathMapping(domainName, basePath, mapping); err != nil {
			return nil, err
		}
	}

	return s.toBasePathMappingResponse(mapping), nil
}

// GetBasePathMappings lists all base path mappings for a domain name.
func (s *APIGatewayService) GetBasePathMappings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	domainName, err := resolveDomainName(req, stores)
	if err != nil {
		return nil, err
	}

	maxItems, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	marker := request.GetStringParam(req.Parameters, "position")

	result, err := stores.domains.ListBasePathMappings(domainName, common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, m := range result.Items {
		items = append(items, s.toBasePathMappingResponse(m))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toBasePathMappingResponse(m *store.BasePathMapping) map[string]interface{} {
	return map[string]interface{}{
		"basePath":  m.BasePath,
		"restApiId": m.RestApiId,
		"stage":     m.Stage,
	}
}
