package apigateway

import (
	"context"
	"fmt"
	"slices"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
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
		return "", toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
	}

	if certArn != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.RegisterCertificateUsage(ctx, region, certArn, created.DomainNameArn); err != nil {
			_ = stores.domains.DeleteDomainName(created.DomainName)
			return nil, NewInternalFailureException(fmt.Sprintf("failed to register certificate usage: %v", err))
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
		return nil, toApiGatewayError(err)
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
		return toApiGatewayError(err)
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
		return toApiGatewayError(err)
	}

	// Best-effort ACM unregister after successful deletion. A stale InUseBy
	// reference is harmless (only over-blocks certificate deletion), whereas
	// unregistering before deletion risks leaving a live resource unprotected.
	if certArnForCleanup != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.UnregisterCertificateUsage(ctx, region, certArnForCleanup, domainArn); err != nil {
			logs.Warn("failed to unregister certificate usage for deleted domain", logs.String("domainName", domainName), logs.Err(err))
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
		return nil, toApiGatewayError(err)
	}

	// Capture old state for compensation.
	oldCertArn := effectiveCertificateArn(domain)
	oldCertArnValue := domain.CertificateArn
	oldRegionalCertArnValue := domain.RegionalCertificateArn
	oldCertificateName := domain.CertificateName

	// The certificate and endpoint-type rows exclude add and remove on the
	// same path within one request (their table cells state so verbatim);
	// seenPatchOps carries the first operation seen per path.
	seenPatchOps := make(map[string]string)

	for _, po := range ops {
		handled, err := applyDomainNamePatch(domain, po, seenPatchOps)
		if err != nil {
			return nil, err
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if err := requireDomainCertificate(domain); err != nil {
		return nil, err
	}

	// Pre-validate new cert ARN.
	newCertArn := effectiveCertificateArn(domain)
	if newCertArn != "" && s.acmInvoker != nil && newCertArn != oldCertArn {
		if !s.acmInvoker.CertificateExists(ctx, region, newCertArn) {
			return nil, NewBadRequestException("The specified certificate ARN does not exist: " + newCertArn)
		}
	}

	if err := stores.domains.UpdateDomainName(domain); err != nil {
		return nil, toApiGatewayError(err)
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
					logs.Error("failed to revert domain name after unregister failure", logs.Err(revertErr))
				}
				return nil, NewInternalFailureException(fmt.Sprintf("failed to unregister old certificate usage: %v", err))
			}
		}
		// Step 2: Register new cert.
		if newCertArn != "" {
			if err := s.acmInvoker.RegisterCertificateUsage(ctx, region, newCertArn, domain.DomainNameArn); err != nil {
				// Compensate: re-register old cert (was unregistered in step 1).
				if oldCertArn != "" {
					if revertErr := s.acmInvoker.RegisterCertificateUsage(ctx, region, oldCertArn, domain.DomainNameArn); revertErr != nil {
						logs.Error("failed to re-register old certificate during compensation", logs.Err(revertErr))
					}
				}
				// Revert to old cert values.
				domain.CertificateArn = oldCertArnValue
				domain.RegionalCertificateArn = oldRegionalCertArnValue
				domain.CertificateName = oldCertificateName
				if revertErr := stores.domains.UpdateDomainName(domain); revertErr != nil {
					logs.Error("failed to revert domain name after register failure", logs.Err(revertErr))
				}
				return nil, NewInternalFailureException(fmt.Sprintf("failed to register new certificate usage: %v", err))
			}
		}
	}

	return domain, nil
}

// applyDomainNamePatch applies one patch operation to a domain name per the
// official UpdateDomainName table. seenOps carries the first operation seen
// per path so the certificate and endpoint-type rows can enforce their
// documented same-request add/remove exclusion. Returns (handled, err):
// handled=false signals an unrecognised path for the unknown-patch-path
// error.
func applyDomainNamePatch(domain *apigateway.DomainName, po PatchOperation, seenOps map[string]string) (bool, error) {
	// rejectCertExclusion enforces the cells' "This operation cannot be
	// included with the add/remove operation in the same request".
	rejectCertExclusion := func() error {
		first, ok := seenOps[po.Path]
		if !ok {
			seenOps[po.Path] = po.Op
			return nil
		}
		if (first == "add" && po.Op == "remove") || (first == "remove" && po.Op == "add") {
			return NewBadRequestException(fmt.Sprintf(
				"the %s and %s operations on '%s' cannot be included in the same request", first, po.Op, po.Path))
		}
		return nil
	}

	switch po.Path {
	case "/certificateArn", "/regionalCertificateArn":
		// The certificate ARN rows document add, replace and remove
		// conditionally (edge/regional transitions); remove clears the
		// addressed certificate, and the requireDomainCertificate invariant
		// keeps one certificate on the domain.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return true, err
		}
		if err := rejectCertExclusion(); err != nil {
			return true, err
		}
		if po.Op == "remove" {
			if po.Path == "/certificateArn" {
				domain.CertificateArn = ""
			} else {
				domain.RegionalCertificateArn = ""
			}
			return true, nil
		}
		if po.Value == "" {
			if po.Path == "/certificateArn" {
				return true, NewBadRequestException("certificateArn cannot be cleared; provide a new certificate ARN")
			}
			return true, NewBadRequestException("regionalCertificateArn cannot be cleared; provide a new certificate ARN")
		}
		if po.Path == "/certificateArn" {
			domain.CertificateArn = po.Value
		} else {
			domain.RegionalCertificateArn = po.Value
		}
	case "/certificateName", "/regionalCertificateName":
		// The certificate name rows document add, replace and remove with
		// the same same-request exclusion.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return true, err
		}
		if err := rejectCertExclusion(); err != nil {
			return true, err
		}
		if po.Op == "remove" {
			if po.Path == "/certificateName" {
				domain.CertificateName = ""
			} else {
				domain.RegionalCertificateName = ""
			}
			return true, nil
		}
		if po.Path == "/certificateName" {
			domain.CertificateName = po.Value
		} else {
			domain.RegionalCertificateName = po.Value
		}
	case "/securityPolicy":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateSecurityPolicy(po.Value) {
			return true, NewBadRequestException("Invalid securityPolicy: must be TLS_1_0, TLS_1_2, or start with SecurityPolicy_")
		}
		domain.SecurityPolicy = po.Value
	case "/ownershipVerificationCertificateArn":
		// The row documents add, replace and remove.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return true, err
		}
		domain.OwnershipVerificationCertificateArn = po.Value
	case "/routingMode":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateRoutingMode(po.Value) {
			return true, NewBadRequestException("Invalid routingMode: " + po.Value)
		}
		domain.RoutingMode = po.Value
	case "/policy":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validatePolicyJSON(po.Value) {
			return true, NewBadRequestException("Invalid policy: must be valid JSON")
		}
		domain.Policy = po.Value
	case "/managementPolicy":
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validatePolicyJSON(po.Value) {
			return true, NewBadRequestException("Invalid managementPolicy: must be valid JSON")
		}
		domain.ManagementPolicy = po.Value
	case "/endpointConfiguration/types":
		// The row documents add and remove for "updates between
		// edge-optimized and regional endpoints"; replace is not supported,
		// and add excludes remove within the same request. The developer
		// guide's migration flow gives the operations list semantics: the
		// new type is added to the existing types list (its output example
		// shows "types": ["EDGE", "REGIONAL"] with both coexisting until the
		// DNS cutover), and the obsolete type is removed from the list in a
		// later request.
		if err := requirePatchOp(po, opAdd|opRemove); err != nil {
			return true, err
		}
		if err := rejectCertExclusion(); err != nil {
			return true, err
		}
		if po.Value != "EDGE" && po.Value != "REGIONAL" {
			return true, NewBadRequestException("Invalid endpoint type: must be EDGE or REGIONAL")
		}
		if domain.EndpointConfiguration == nil {
			domain.EndpointConfiguration = &apigateway.EndpointConfiguration{}
		}
		types := domain.EndpointConfiguration.Types
		if po.Op == "remove" {
			types = slices.DeleteFunc(types, func(t string) bool { return t == po.Value })
		} else if !slices.Contains(types, po.Value) {
			types = append(types, po.Value)
		}
		domain.EndpointConfiguration.Types = types
	case "/endpointConfiguration/ipAddressType":
		// The row documents replace only, "Only dualstack and ipv4 are
		// supported"; the model notes dualstack is the only value a
		// PRIVATE endpoint admits.
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateIpAddressType(po.Value) {
			return true, NewBadRequestException("Invalid ipAddressType: must be ipv4 or dualstack")
		}
		if domain.EndpointConfiguration != nil &&
			slices.Contains(domain.EndpointConfiguration.Types, "PRIVATE") && po.Value != "dualstack" {
			return true, NewBadRequestException("ipAddressType for a PRIVATE endpoint must be dualstack")
		}
		if domain.EndpointConfiguration == nil {
			domain.EndpointConfiguration = &apigateway.EndpointConfiguration{}
		}
		domain.EndpointConfiguration.IpAddressType = po.Value
	case "/endpointAccessMode":
		// The row documents replace only.
		if err := requirePatchOp(po, opReplace); err != nil {
			return true, err
		}
		if !validateEndpointAccessMode(po.Value) {
			return true, NewBadRequestException("Invalid endpointAccessMode: must be BASIC or STRICT")
		}
		domain.EndpointAccessMode = po.Value
	case "/mutualTlsAuthentication/truststoreUri":
		// The row documents add, replace and remove.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return true, err
		}
		if domain.MutualTlsAuthentication == nil {
			domain.MutualTlsAuthentication = &apigateway.MutualTlsAuthentication{}
		}
		if po.Op == "remove" {
			domain.MutualTlsAuthentication.TruststoreUri = ""
		} else {
			domain.MutualTlsAuthentication.TruststoreUri = po.Value
		}
	case "/mutualTlsAuthentication/truststoreVersion":
		// The row documents add, replace and remove.
		if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
			return true, err
		}
		if domain.MutualTlsAuthentication == nil {
			domain.MutualTlsAuthentication = &apigateway.MutualTlsAuthentication{}
		}
		if po.Op == "remove" {
			domain.MutualTlsAuthentication.TruststoreVersion = ""
		} else {
			domain.MutualTlsAuthentication.TruststoreVersion = po.Value
		}
	default:
		return false, nil
	}
	return true, nil
}

// requireDomainCertificate returns an error when a patch request leaves the
// domain without any certificate: the create contract requires certificateArn
// or regionalCertificateArn, and the documented certificate removes serve
// edge/regional transitions only.
func requireDomainCertificate(domain *apigateway.DomainName) error {
	if effectiveCertificateArn(domain) == "" {
		return NewBadRequestException("a domain name must keep either certificateArn or regionalCertificateArn")
	}
	return nil
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
		return nil, toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
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
		return toApiGatewayError(err)
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
		return nil, toApiGatewayError(err)
	}

	renamed := false
	oldBasePath := ""
	oldRestApiId := mapping.RestApiId
	oldStage := mapping.Stage
	for _, po := range ops {
		handled := false
		switch po.Path {
		case "/restApiId", "/restapiId":
			handled = true
			// The patch table prints both rows lowercase ("/basepath",
			// "/restapiId") while the official CLI example uses the member
			// casing ("path='/basePath'"); each row accepts both spellings.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			mapping.RestApiId = po.Value
		case "/stage":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			mapping.Stage = po.Value
		case "/basePath", "/basepath":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateBasePath(po.Value) {
				return nil, NewBadRequestException("basePath must contain only alphanumeric characters, hyphens, underscores, periods, and forward slashes")
			}
			oldBasePath = basePath
			basePath = po.Value
			mapping.BasePath = po.Value
			renamed = true
		}
		if !handled {
			return nil, unknownPatchPathError(po)
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
