package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// --- DTOs ---

// DomainNameCreateInput carries the parsed wire members of a CreateDomainName
// request. HasMutualTlsAuthentication and HasEndpointConfiguration preserve
// wire presence so an explicitly provided empty structure still creates the
// stored sub-structure, exactly as the raw parameter walk did.
type DomainNameCreateInput struct {
	DomainName                          string
	CertificateArn                      string
	CertificateName                     string
	CertificateBody                     string
	CertificatePrivateKey               string
	CertificateChain                    string
	RegionalCertificateArn              string
	RegionalCertificateName             string
	SecurityPolicy                      string
	OwnershipVerificationCertificateArn string
	EndpointAccessMode                  string
	Policy                              string
	RoutingMode                         string
	HasMutualTlsAuthentication          bool
	MutualTlsTruststoreUri              string
	MutualTlsTruststoreVersion          string
	HasEndpointConfiguration            bool
	EndpointTypes                       []string
	Tags                                []tagutil.Tag
}

// BasePathMappingInput carries the parsed wire members of a
// CreateBasePathMapping request.
type BasePathMappingInput struct {
	BasePath  string
	RestApiId string
	Stage     string
}

// effectiveCertificateArn returns the ACM certificate ARN that a domain name
// is configured to use, preferring the regional certificate over the edge
// certificate. Returns empty string if no ACM certificate is configured.
func effectiveCertificateArn(d *apigateway.DomainName) string {
	if d.RegionalCertificateArn != "" {
		return d.RegionalCertificateArn
	}
	return d.CertificateArn
}

// resolveDomainNameCore resolves a domain name from either the domainName or
// domainNameId request member. At least one must be provided per the Smithy
// model. When only domainNameId is given, the store is queried to obtain the
// canonical domain name.
func (s *APIGatewayService) resolveDomainNameCore(stores *apiGatewayStores, domainName, domainNameId string) (string, error) {
	if domainName != "" {
		return domainName, nil
	}

	if domainNameId == "" {
		return "", NewBadRequestException("Either domainName or domainNameId must be specified")
	}

	domain, err := stores.domains.GetDomainNameById(domainNameId)
	if err != nil {
		return "", ErrNotFoundException
	}
	return domain.DomainName, nil
}

// createDomainNameCore validates the input, persists the domain name and
// registers its ACM certificate usage. The validation order mirrors the
// original handler: identity, security policy and endpoint access mode
// precede the certificate checks, which precede the routing-mode and policy
// checks.
func (s *APIGatewayService) createDomainNameCore(
	ctx context.Context,
	stores *apiGatewayStores,
	region string,
	in *DomainNameCreateInput,
) (*apigateway.DomainName, error) {
	if in.DomainName == "" {
		return nil, NewBadRequestException("domainName is required")
	}
	if !validateFQDN(in.DomainName) {
		return nil, NewBadRequestException("Invalid domain name: must be a valid FQDN")
	}
	if !validateSecurityPolicy(in.SecurityPolicy) {
		return nil, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
	}
	if !validateEndpointAccessMode(in.EndpointAccessMode) {
		return nil, NewBadRequestException("Invalid endpointAccessMode: must be BASIC or STRICT")
	}

	domain := &apigateway.DomainName{
		DomainName:                          in.DomainName,
		CertificateArn:                      in.CertificateArn,
		CertificateName:                     in.CertificateName,
		RegionalCertificateArn:              in.RegionalCertificateArn,
		RegionalCertificateName:             in.RegionalCertificateName,
		SecurityPolicy:                      in.SecurityPolicy,
		OwnershipVerificationCertificateArn: in.OwnershipVerificationCertificateArn,
		EndpointAccessMode:                  in.EndpointAccessMode,
		Policy:                              in.Policy,
		RoutingMode:                         in.RoutingMode,
	}
	if in.HasMutualTlsAuthentication {
		domain.MutualTlsAuthentication = &apigateway.MutualTlsAuthentication{
			TruststoreUri:     in.MutualTlsTruststoreUri,
			TruststoreVersion: in.MutualTlsTruststoreVersion,
		}
	}
	if in.HasEndpointConfiguration {
		domain.EndpointConfiguration = &apigateway.EndpointConfiguration{Types: in.EndpointTypes}
	}
	domain.Tags = in.Tags

	// Require an ACM certificate ARN — vorpalstacks only supports the ACM
	// certificate ARN model and does not accept the certificateBody /
	// certificatePrivateKey / certificateChain upload path.  Rejecting
	// explicitly instead of silently dropping prevents a fail-OPEN gap where
	// a caller could create a domain with no certificate at all.
	certArn := effectiveCertificateArn(domain)
	if certArn == "" {
		if in.CertificateBody != "" || in.CertificatePrivateKey != "" || in.CertificateChain != "" {
			return nil, NewBadRequestException("certificateBody, certificatePrivateKey, and certificateChain are not supported; use certificateArn or regionalCertificateArn")
		}
		return nil, NewBadRequestException("certificateArn or regionalCertificateArn is required")
	}
	if s.acmInvoker != nil && !s.acmInvoker.CertificateExists(ctx, region, certArn) {
		return nil, NewBadRequestException("The specified certificate ARN does not exist: " + certArn)
	}

	if !validateRoutingMode(domain.RoutingMode) {
		return nil, NewBadRequestException("Invalid routingMode: " + domain.RoutingMode)
	}

	if !validatePolicyJSON(domain.Policy) {
		return nil, NewBadRequestException("Invalid policy: must be valid JSON")
	}

	created, err := stores.domains.CreateDomainName(domain)
	if err != nil {
		return nil, err
	}

	if certArn != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.RegisterCertificateUsage(ctx, region, certArn, created.DomainNameArn); err != nil {
			_ = stores.domains.DeleteDomainName(created.DomainName)
			return nil, fmt.Errorf("failed to register certificate usage: %w", err)
		}
	}

	return created, nil
}

// getDomainNameCore resolves the domain identity and returns the stored
// domain name.
func (s *APIGatewayService) getDomainNameCore(stores *apiGatewayStores, domainName, domainNameId string) (*apigateway.DomainName, error) {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return nil, err
	}

	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return nil, ErrNotFoundException
	}
	return domain, nil
}

// deleteDomainNameCore removes a domain name after verifying no base path
// mappings remain, then best-effort unregisters its ACM certificate usage.
func (s *APIGatewayService) deleteDomainNameCore(
	ctx context.Context,
	stores *apiGatewayStores,
	region, domainName, domainNameId string,
) error {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return err
	}

	domain, err := stores.domains.GetDomainName(domainName)
	if err != nil {
		return ErrNotFoundException
	}

	// AWS requires all BasePathMappings to be deleted before a domain name
	// can be removed; otherwise it returns ConflictException. The listing is
	// fail-closed: when the store cannot answer whether mappings exist, the
	// deletion is refused rather than proceeding on an unknown state.
	mappings, err := stores.domains.ListBasePathMappings(domainName, storecommon.ListOptions{MaxItems: 1})
	if err != nil {
		return toApiGatewayError(err)
	}
	if len(mappings.Items) > 0 {
		return NewConflictException("Domain name has active base path mappings; remove them first")
	}

	// Capture cert ARN before deletion for best-effort unregister after.
	certArnForCleanup := effectiveCertificateArn(domain)
	domainArn := domain.DomainNameArn

	if err := stores.domains.DeleteDomainName(domainName); err != nil {
		return ErrNotFoundException
	}

	// Best-effort ACM unregister after successful deletion. A stale InUseBy
	// reference is harmless (only over-blocks certificate deletion), whereas
	// unregistering before deletion risks leaving a live resource unprotected.
	if certArnForCleanup != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.UnregisterCertificateUsage(ctx, region, certArnForCleanup, domainArn); err != nil {
			log.Printf("warning: failed to unregister certificate usage for deleted domain %s: %v", domainName, err)
		}
	}

	return nil
}

// updateDomainNameCore applies JSON Patch operations to a domain name with
// ACM certificate compensation on failure.
func (s *APIGatewayService) updateDomainNameCore(
	ctx context.Context,
	stores *apiGatewayStores,
	region, domainName, domainNameId string,
	ops []PatchOperation,
) (*apigateway.DomainName, error) {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
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
				domain.EndpointConfiguration = &apigateway.EndpointConfiguration{}
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
				domain.MutualTlsAuthentication = &apigateway.MutualTlsAuthentication{}
			}
			if po.Op == "remove" {
				domain.MutualTlsAuthentication.TruststoreUri = ""
			} else {
				domain.MutualTlsAuthentication.TruststoreUri = po.Value
			}
		case po.Path == "/mutualTlsAuthentication/truststoreVersion":
			if domain.MutualTlsAuthentication == nil {
				domain.MutualTlsAuthentication = &apigateway.MutualTlsAuthentication{}
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
		if !s.acmInvoker.CertificateExists(ctx, region, newCertArn) {
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
			if err := s.acmInvoker.UnregisterCertificateUsage(ctx, region, oldCertArn, domain.DomainNameArn); err != nil {
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
			if err := s.acmInvoker.RegisterCertificateUsage(ctx, region, newCertArn, domain.DomainNameArn); err != nil {
				// Compensate: re-register old cert (was unregistered in step 1).
				if oldCertArn != "" {
					if revertErr := s.acmInvoker.RegisterCertificateUsage(ctx, region, oldCertArn, domain.DomainNameArn); revertErr != nil {
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

	return domain, nil
}

// listDomainNamesCore returns a page of domain names.
func (s *APIGatewayService) listDomainNamesCore(stores *apiGatewayStores, marker string, maxItems int) (*storecommon.ListResult[apigateway.DomainName], error) {
	return stores.domains.ListDomainNames(storecommon.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
}

// createBasePathMappingCore validates and persists a base path mapping under
// the resolved domain name.
func (s *APIGatewayService) createBasePathMappingCore(
	stores *apiGatewayStores,
	domainName, domainNameId string,
	in *BasePathMappingInput,
) (*apigateway.BasePathMapping, error) {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return nil, err
	}

	if in.RestApiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	mapping := &apigateway.BasePathMapping{
		BasePath:  in.BasePath,
		RestApiId: in.RestApiId,
		Stage:     in.Stage,
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

	return created, nil
}

// getBasePathMappingCore retrieves a base path mapping by domain and base
// path.
func (s *APIGatewayService) getBasePathMappingCore(stores *apiGatewayStores, domainName, domainNameId, basePath string) (*apigateway.BasePathMapping, error) {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return nil, err
	}

	if basePath == "" {
		return nil, NewBadRequestException("basePath is required")
	}

	mapping, err := stores.domains.GetBasePathMapping(domainName, basePath)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return mapping, nil
}

// deleteBasePathMappingCore removes a base path mapping.
func (s *APIGatewayService) deleteBasePathMappingCore(stores *apiGatewayStores, domainName, domainNameId, basePath string) error {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return err
	}

	if basePath == "" {
		return NewBadRequestException("basePath is required")
	}

	if err := stores.domains.DeleteBasePathMapping(domainName, basePath); err != nil {
		return ErrNotFoundException
	}

	return nil
}

// updateBasePathMappingCore applies JSON Patch operations to a base path
// mapping under the domain-and-base-path key lock, including the rename
// delete-and-recreate path with its compensating restore.
func (s *APIGatewayService) updateBasePathMappingCore(
	stores *apiGatewayStores,
	domainName, domainNameId, basePath string,
	ops []PatchOperation,
) (*apigateway.BasePathMapping, error) {
	domainName, err := s.resolveDomainNameCore(stores, domainName, domainNameId)
	if err != nil {
		return nil, err
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

	return mapping, nil
}

// listBasePathMappingsCore returns a page of base path mappings for a domain
// name. The caller resolves the domain identity first so the pagination
// validation keeps its original position after it.
func (s *APIGatewayService) listBasePathMappingsCore(
	stores *apiGatewayStores,
	domainName, marker string,
	maxItems int,
) (*storecommon.ListResult[apigateway.BasePathMapping], error) {
	return stores.domains.ListBasePathMappings(domainName, storecommon.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
}
