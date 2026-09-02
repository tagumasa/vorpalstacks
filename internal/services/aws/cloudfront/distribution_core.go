package cloudfront

import (
	"github.com/google/uuid"

	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/resilience"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateDistributionInput carries every field that CreateDistribution needs,
// in a format independent of the wire protocol (HTTP Query/XML vs gRPC-Web).
// Both the HTTP API handler (distribution_operations.go) and the admin
// gRPC handler (admin_handler.go) build this struct from their respective
// request formats and delegate to createDistributionCore, ensuring that
// validation, certificate checking, WAF sync, and persistence follow a
// single code path.
type CreateDistributionInput struct {
	CallerReference string
	Config          *cloudfrontstore.DistributionConfig
	Tags            []types.Tag
	TagsProvided    bool
	// ACMRegion is the region used for ACM certificate lookups. In real AWS,
	// CloudFront only accepts certificates from us-east-1. In vorpalstacks
	// the deployment region is used.
	ACMRegion string
}

// ListDistributionsInput carries the pagination parameters for
// ListDistributions.
type ListDistributionsInput struct {
	Marker   string
	MaxItems int
}

// DeleteDistributionInput carries the parameters for DeleteDistribution.
type DeleteDistributionInput struct {
	Id      string
	IfMatch string
	// ACMRegion is used for best-effort ACM certificate usage cleanup.
	ACMRegion string
}

// UpdateDistributionInput carries every field that UpdateDistribution
// needs, in a format independent of the wire protocol. Both the HTTP API
// handler and future admin gRPC handlers build this struct and delegate
// to updateDistributionCore, ensuring that validation, WAF sync, and ACM
// usage tracking follow a single code path.
type UpdateDistributionInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.DistributionConfig
	// ACMRegion is the region used for ACM certificate lookups.
	ACMRegion string
}

// CreateDistributionResult is the transport-agnostic result of
// createDistributionCore.
type CreateDistributionResult struct {
	Distribution *cloudfrontstore.Distribution
}

// UpdateDistributionResult is the transport-agnostic result of
// updateDistributionCore.
type UpdateDistributionResult struct {
	Distribution *cloudfrontstore.Distribution
}

// ListDistributionsResult is the transport-agnostic result of
// listDistributionsCore.
type ListDistributionsResult struct {
	Distributions []*cloudfrontstore.Distribution
	IsTruncated   bool
	NextMarker    string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// validateDistributionConfig enforces the structural requirements of the
// DistributionConfig shape shared by createDistributionCore and
// updateDistributionCore: every origin needs a domain name and a unique
// ID, alias CNAMEs are bounded by the Smithy aliasString length, and the
// enum fields must hold valid values.
func validateDistributionConfig(config *cloudfrontstore.DistributionConfig) error {
	if len(config.Origins.Items) == 0 {
		return invalidArgument("At least one origin is required")
	}
	if config.DefaultCacheBehavior == nil {
		return invalidArgument("DefaultCacheBehavior is required")
	}
	originIDs := make(map[string]bool, len(config.Origins.Items))
	for _, origin := range config.Origins.Items {
		if origin == nil {
			continue
		}
		if origin.DomainName == "" {
			return invalidArgument("Origin.DomainName is required for every origin")
		}
		if origin.ID == "" {
			return invalidArgument("Origin.Id is required for every origin")
		}
		if originIDs[origin.ID] {
			return invalidArgument("Origin.Id must be unique within the distribution: " + origin.ID)
		}
		originIDs[origin.ID] = true
		if origin.CustomOriginConfig != nil {
			opp := origin.CustomOriginConfig.OriginProtocolPolicy
			if opp == "" {
				return invalidArgument("OriginProtocolPolicy is required when CustomOriginConfig is specified")
			}
			if !isValidOriginProtocolPolicy(opp) {
				return invalidArgument("Invalid OriginProtocolPolicy: " + opp)
			}
		}
	}
	if config.Aliases != nil {
		for _, alias := range config.Aliases.Items {
			if len(alias) > cloudfrontstore.MaxAliasItemLength {
				return invalidArgument(fmt.Sprintf("Alias CNAME exceeds the maximum length of %d characters", cloudfrontstore.MaxAliasItemLength))
			}
		}
	}
	if config.PriceClass != "" && !isValidPriceClass(config.PriceClass) {
		return invalidArgument("Invalid PriceClass: " + config.PriceClass)
	}
	if config.HttpVersion != "" && !isValidHttpVersion(config.HttpVersion) {
		return invalidArgument("Invalid HttpVersion: " + config.HttpVersion)
	}
	return nil
}

// ensureCNAMEsAvailable rejects a distribution configuration whose
// alternate domain names are already defined on another non-staging
// distribution, mirroring the CNAMEAlreadyExists error the create,
// update, and copy operations model. A staging configuration is exempt:
// staging distributions never receive viewer traffic, so they may share
// CNAMEs with the distribution they were copied from. selfID excludes
// the distribution being updated from the conflict scan.
func ensureCNAMEsAvailable(stores *cloudfrontStores, config *cloudfrontstore.DistributionConfig, selfID string) error {
	if config == nil || config.Staging || config.Aliases == nil || len(config.Aliases.Items) == 0 {
		return nil
	}
	// Collect every CNAME already in use in a single pass so the check
	// costs one scan regardless of how many aliases the request carries.
	inUse := make(map[string]struct{})
	if _, err := scanDistributions(stores, func(d *cloudfrontstore.Distribution) bool {
		if d.ID == selfID || d.Staging || d.DistributionConfig == nil || d.DistributionConfig.Aliases == nil {
			return false
		}
		for _, other := range d.DistributionConfig.Aliases.Items {
			if other != "" {
				inUse[strings.ToLower(other)] = struct{}{}
			}
		}
		return false
	}); err != nil {
		return awserrors.NewAWSError("InternalError", "Failed to scan distributions for CNAME conflicts: "+err.Error(), 500)
	}
	for _, alias := range config.Aliases.Items {
		if alias == "" {
			continue
		}
		if _, taken := inUse[strings.ToLower(alias)]; taken {
			return awserrors.NewAWSError("CNAMEAlreadyExists",
				"The CNAME specified is already defined for CloudFront: "+alias, 409)
		}
	}
	return nil
}

// validateReferencedPolicies verifies that the cache policies, origin
// request policies, and response headers policies referenced by the cache
// behaviours of a distribution configuration exist. The distribution
// create, update, and copy operations model the NoSuchCachePolicy,
// NoSuchOriginRequestPolicy, and NoSuchResponseHeadersPolicy errors for
// dangling references.
func (s *CloudFrontService) validateReferencedPolicies(stores *cloudfrontStores, config *cloudfrontstore.DistributionConfig) error {
	var behaviours []*cloudfrontstore.CacheBehavior
	if config.DefaultCacheBehavior != nil {
		behaviours = append(behaviours, config.DefaultCacheBehavior)
	}
	if config.CacheBehaviors != nil {
		for _, cb := range config.CacheBehaviors.Items {
			if cb != nil {
				behaviours = append(behaviours, cb)
			}
		}
	}
	for _, cb := range behaviours {
		if cb.CachePolicyId != "" {
			if _, err := stores.cachePolicies.Get(cb.CachePolicyId); err != nil {
				if cloudfrontstore.IsNotFound(err) {
					return awserrors.NewAWSError("NoSuchCachePolicy", "The cache policy does not exist: "+cb.CachePolicyId, 404)
				}
				return err
			}
		}
		if cb.OriginRequestPolicyId != "" {
			if _, err := stores.originRequestPolicies.Get(cb.OriginRequestPolicyId); err != nil {
				if cloudfrontstore.IsNotFound(err) {
					return awserrors.NewAWSError("NoSuchOriginRequestPolicy", "The origin request policy does not exist: "+cb.OriginRequestPolicyId, 404)
				}
				return err
			}
		}
		if cb.ResponseHeadersPolicyId != "" {
			if _, err := stores.responseHeadersPolicies.Get(cb.ResponseHeadersPolicyId); err != nil {
				if cloudfrontstore.IsNotFound(err) {
					return awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "The response headers policy does not exist: "+cb.ResponseHeadersPolicyId, 404)
				}
				return err
			}
		}
	}
	return nil
}

// validateTrustedKeyGroups verifies that every key group a cache behaviour
// trusts exists. The Developer Guide's signing workflow requires a key group
// to be created before it is attached to a distribution, and unknown key
// group identifiers are rejected with InvalidArgument.
func (s *CloudFrontService) validateTrustedKeyGroups(stores *cloudfrontStores, config *cloudfrontstore.DistributionConfig) error {
	var behaviours []*cloudfrontstore.CacheBehavior
	if config.DefaultCacheBehavior != nil {
		behaviours = append(behaviours, config.DefaultCacheBehavior)
	}
	if config.CacheBehaviors != nil {
		for _, cb := range config.CacheBehaviors.Items {
			if cb != nil {
				behaviours = append(behaviours, cb)
			}
		}
	}
	for _, cb := range behaviours {
		if cb.TrustedKeyGroups == nil {
			continue
		}
		for _, kgID := range cb.TrustedKeyGroups.Items {
			if kgID == "" {
				continue
			}
			if _, err := stores.keyGroups.Get(kgID); err != nil {
				if cloudfrontstore.IsNotFound(err) {
					return awserrors.NewAWSError("InvalidArgument", "The specified key group does not exist: "+kgID, 400)
				}
				return err
			}
		}
	}
	return nil
}

// createDistributionCore is the single entry point for distribution creation
// logic shared by the HTTP API and the admin gRPC handler. It performs all
// validation, certificate existence checks, WAF association, persistence,
// and tag application.
func (s *CloudFrontService) createDistributionCore(ctx context.Context, stores *cloudfrontStores, in CreateDistributionInput) (*CreateDistributionResult, error) {
	if in.CallerReference == "" {
		return nil, invalidArgument("CallerReference is required")
	}

	if in.Config == nil {
		return nil, invalidArgument("DistributionConfig is required")
	}
	in.Config.CallerReference = in.CallerReference

	if err := validateDistributionConfig(in.Config); err != nil {
		return nil, err
	}

	if err := s.validateReferencedPolicies(stores, in.Config); err != nil {
		return nil, err
	}

	if err := s.validateDistributionPolicyReference(stores, in.Config); err != nil {
		return nil, err
	}

	if err := s.validateTrustedKeyGroups(stores, in.Config); err != nil {
		return nil, err
	}

	if err := ensureCNAMEsAvailable(stores, in.Config, ""); err != nil {
		return nil, err
	}

	if in.Config.WebACLId != "" && s.wafInvoker != nil && !s.wafInvoker.WebACLExists(ctx, in.Config.WebACLId) {
		return nil, awserrors.NewAWSError("InvalidWebACLId",
			"The specified Web ACL does not exist: "+in.Config.WebACLId, 400)
	}

	certArn := ""
	if in.Config.ViewerCertificate != nil {
		certArn = in.Config.ViewerCertificate.ACMCertificateArn
	}
	if certArn != "" && s.acmInvoker != nil {
		if !s.acmInvoker.CertificateExists(ctx, in.ACMRegion, certArn) {
			return nil, awserrors.NewAWSError("NoSuchCertificate", "The specified certificate ARN does not exist: "+certArn, 404)
		}
	}

	// A CallerReference already used by another distribution must be
	// rejected with the modelled DistributionAlreadyExists error instead
	// of silently returning the existing distribution, which would let a
	// replayed reference mutate an unrelated distribution.
	if _, err := stores.distributions.GetByCallerReference(in.CallerReference); err == nil {
		return nil, awserrors.NewAWSError("DistributionAlreadyExists",
			"The caller reference you attempted to create the distribution with is associated with another distribution.", 409)
	} else if !cloudfrontstore.IsNotFound(err) {
		return nil, err
	}

	distribution, err := stores.distributions.Create(in.CallerReference, in.Config)
	if err != nil {
		return nil, err
	}

	if in.Config.WebACLId != "" && s.wafInvoker != nil {
		if err := s.wafInvoker.AssociateWebACL(in.Config.WebACLId, distribution.ARN); err != nil {
			_ = stores.distributions.Delete(distribution.ID)
			return nil, fmt.Errorf("failed to sync WAF association: %w", err)
		}
	}

	if certArn != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.RegisterCertificateUsage(ctx, in.ACMRegion, certArn, distribution.ARN); err != nil {
			_ = stores.distributions.Delete(distribution.ID)
			return nil, fmt.Errorf("failed to register certificate usage: %w", err)
		}
	}

	if in.TagsProvided && len(in.Tags) > 0 && distribution.ARN != "" {
		if err := stores.tags.Tag(distribution.ARN, in.Tags); err != nil {
			return nil, fmt.Errorf("failed to tag distribution: %w", err)
		}
	}

	go s.transitionDistributionDeployed(stores, distribution.ID)

	return &CreateDistributionResult{Distribution: distribution}, nil
}

// AdminCreateDistributionInput carries the simplified parameters that the
// admin console provides when creating a distribution. Unlike the HTTP API
// which sends a full DistributionConfig XML, the admin console sends only
// the essential fields and relies on sensible defaults for the rest.
type AdminCreateDistributionInput struct {
	CallerReference string
	Comment         string
	Enabled         bool
	OriginID        string
	OriginDomain    string
	ACMRegion       string
}

// createDistributionFromAdmin is the Core entry point for the admin console.
// It builds a full DistributionConfig from the simplified admin parameters
// and delegates to createDistributionCore, ensuring that all validation,
// certificate checking, and persistence follow the single code path.
func (s *CloudFrontService) createDistributionFromAdmin(ctx context.Context, stores *cloudfrontStores, in AdminCreateDistributionInput) (*CreateDistributionResult, error) {
	// A configuration without any origin cannot form a valid distribution.
	if in.OriginDomain == "" {
		return nil, invalidArgument("distributionconfig is required")
	}
	if in.OriginID == "" {
		in.OriginID = "default-origin"
	}
	// The caller reference is CloudFront's idempotency token; mint it from
	// crypto/rand so concurrent creations cannot collide on a clock value.
	if in.CallerReference == "" {
		in.CallerReference = uuid.New().String()
	}
	config := buildDefaultDistributionConfig(in.CallerReference, in.Comment, in.Enabled, in.OriginID, in.OriginDomain)
	return s.createDistributionCore(ctx, stores, CreateDistributionInput{
		CallerReference: in.CallerReference,
		Config:          config,
		ACMRegion:       in.ACMRegion,
	})
}

// updateDistributionCore is the single entry point for distribution
// updates, shared by the HTTP API and future admin gRPC handlers. It
// applies the same validation as creation, persists the new
// configuration first, and then synchronises WAF and ACM certificate
// usage; every external step has a compensating action that restores the
// previous stored configuration and WAF association on failure.
func (s *CloudFrontService) updateDistributionCore(ctx context.Context, stores *cloudfrontStores, in UpdateDistributionInput) (*UpdateDistributionResult, error) {
	if in.Id == "" {
		return nil, invalidArgument("Id is required")
	}
	if in.Config == nil {
		return nil, invalidArgument("DistributionConfig is required")
	}

	distribution, err := stores.distributions.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	if err := verifyIfMatch(in.IfMatch, distribution.ETag); err != nil {
		return nil, err
	}

	// Capture old state for compensation.
	oldConfig := distribution.DistributionConfig
	oldEnabled := distribution.Enabled
	oldWebACLId := ""
	oldCertArn := ""
	if oldConfig != nil {
		oldWebACLId = oldConfig.WebACLId
		if oldConfig.ViewerCertificate != nil {
			oldCertArn = oldConfig.ViewerCertificate.ACMCertificateArn
		}
	}

	if in.Config.CallerReference == "" {
		in.Config.CallerReference = distribution.CallerReference
	}

	if err := validateDistributionConfig(in.Config); err != nil {
		return nil, err
	}

	if err := s.validateReferencedPolicies(stores, in.Config); err != nil {
		return nil, err
	}

	if err := s.validateDistributionPolicyReference(stores, in.Config); err != nil {
		return nil, err
	}

	if err := s.validateTrustedKeyGroups(stores, in.Config); err != nil {
		return nil, err
	}

	if err := ensureCNAMEsAvailable(stores, in.Config, in.Id); err != nil {
		return nil, err
	}

	if in.Config.WebACLId != "" && s.wafInvoker != nil && !s.wafInvoker.WebACLExists(ctx, in.Config.WebACLId) {
		return nil, awserrors.NewAWSError("InvalidWebACLId",
			"The specified Web ACL does not exist: "+in.Config.WebACLId, 400)
	}

	// Pre-validate new cert ARN before saving.
	newCertArn := ""
	if in.Config.ViewerCertificate != nil {
		newCertArn = in.Config.ViewerCertificate.ACMCertificateArn
	}
	if newCertArn != "" && s.acmInvoker != nil && newCertArn != oldCertArn {
		if !s.acmInvoker.CertificateExists(ctx, in.ACMRegion, newCertArn) {
			return nil, awserrors.NewAWSError("NoSuchCertificate", "The specified certificate ARN does not exist: "+newCertArn, 404)
		}
	}

	// The store is the primary transaction: persist the new
	// configuration before touching the external WAF and ACM state so
	// that any later failure can be compensated from a consistent base.
	distribution.DistributionConfig = in.Config
	distribution.Enabled = in.Config.Enabled

	if err := stores.distributions.UpdateWithLastModified(in.Id, distribution); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}

	revert := func() {
		s.revertDistributionUpdate(stores, in.Id, distribution, oldConfig, oldEnabled, oldWebACLId, in.Config.WebACLId)
	}

	if oldWebACLId != in.Config.WebACLId && s.wafInvoker != nil {
		distArn := distribution.ARN
		if oldWebACLId != "" {
			if err := s.wafInvoker.DisassociateWebACL(oldWebACLId, distArn); err != nil {
				revert()
				return nil, fmt.Errorf("failed to remove old WAF association: %w", err)
			}
		}
		if in.Config.WebACLId != "" {
			if err := s.wafInvoker.AssociateWebACL(in.Config.WebACLId, distArn); err != nil {
				// The disassociation above already succeeded, so the
				// revert inside re-establishes the old association.
				revert()
				return nil, fmt.Errorf("failed to sync WAF association: %w", err)
			}
		}
	}

	// ACM cert operations with compensating transactions on failure; a
	// failed certificate step also reverts the configuration and the WAF
	// association switched above.
	if s.acmInvoker != nil && oldCertArn != newCertArn {
		distArn := distribution.ARN
		// Step 1: Unregister old cert.
		if oldCertArn != "" {
			if err := s.acmInvoker.UnregisterCertificateUsage(ctx, in.ACMRegion, oldCertArn, distArn); err != nil {
				revert()
				return nil, fmt.Errorf("failed to unregister old certificate usage: %w", err)
			}
		}
		// Step 2: Register new cert.
		if newCertArn != "" {
			if err := s.acmInvoker.RegisterCertificateUsage(ctx, in.ACMRegion, newCertArn, distArn); err != nil {
				// Compensate: re-register old cert (was unregistered in step 1).
				if oldCertArn != "" {
					_ = s.acmInvoker.RegisterCertificateUsage(ctx, in.ACMRegion, oldCertArn, distArn)
				}
				revert()
				return nil, fmt.Errorf("failed to register new certificate usage: %w", err)
			}
		}
	}

	// A configuration change invalidates the distribution's edge cache.
	if s.distributionServer != nil {
		s.distributionServer.PurgeDistribution(in.Id)
		s.distributionServer.PurgeCertificates()
	}

	return &UpdateDistributionResult{Distribution: distribution}, nil
}

// revertDistributionUpdate restores a distribution to its previous
// configuration after a failed external side effect, including the WAF
// association when it had already been switched. The revert is best
// effort: callers report the original error, not a compensation failure.
func (s *CloudFrontService) revertDistributionUpdate(stores *cloudfrontStores, id string, distribution *cloudfrontstore.Distribution, oldConfig *cloudfrontstore.DistributionConfig, oldEnabled bool, oldWebACLId, newWebACLId string) {
	if s.wafInvoker != nil && oldWebACLId != newWebACLId {
		if newWebACLId != "" {
			_ = s.wafInvoker.DisassociateWebACL(newWebACLId, distribution.ARN)
		}
		if oldWebACLId != "" {
			_ = s.wafInvoker.AssociateWebACL(oldWebACLId, distribution.ARN)
		}
	}
	distribution.DistributionConfig = oldConfig
	distribution.Enabled = oldEnabled
	_ = stores.distributions.UpdateWithLastModified(id, distribution)
}

// buildDefaultDistributionConfig constructs a store-type DistributionConfig
// with sensible defaults for fields not provided by the admin console.
func buildDefaultDistributionConfig(callerRef, comment string, enabled bool, originID, originDomain string) *cloudfrontstore.DistributionConfig {
	return &cloudfrontstore.DistributionConfig{
		CallerReference: callerRef,
		Comment:         comment,
		Enabled:         enabled,
		Origins: cloudfrontstore.Origins{
			Quantity: 1,
			Items: []*cloudfrontstore.Origin{
				{
					ID:         originID,
					DomainName: originDomain,
					CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{
						HTTPPort:             80,
						HTTPSPort:            443,
						OriginProtocolPolicy: "https-only",
					},
					ConnectionAttempts: 3,
					ConnectionTimeout:  10,
				},
			},
		},
		DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
			TargetOriginId:       originID,
			ViewerProtocolPolicy: "allow-all",
			AllowedMethods: &cloudfrontstore.AllowedMethods{
				Quantity: 2,
				Items:    []string{"GET", "HEAD"},
			},
			ForwardedValues: &cloudfrontstore.ForwardedValues{
				QueryString: false,
				Cookies:     &cloudfrontstore.CookiePreferences{Forward: "none"},
			},
			MinTTL: 0,
		},
		PriceClass:    "PriceClass_All",
		HttpVersion:   "http2and3",
		IsIPV6Enabled: true,
		ViewerCertificate: &cloudfrontstore.ViewerCertificate{
			CloudFrontDefaultCertificate: true,
			CertificateSource:            "cloudfront",
		},
		Restrictions: &cloudfrontstore.Restrictions{
			GeoRestriction: cloudfrontstore.GeoRestriction{
				RestrictionType: "none",
			},
		},
	}
}

// listDistributionsCore is the single entry point for listing distributions.
func (s *CloudFrontService) listDistributionsCore(stores *cloudfrontStores, in ListDistributionsInput) (*ListDistributionsResult, error) {
	in.MaxItems = resolveListMaxItems(in.MaxItems)
	result, err := stores.distributions.List(in.Marker, in.MaxItems)
	if err != nil {
		return nil, err
	}
	return &ListDistributionsResult{
		Distributions: result.Distributions,
		IsTruncated:   result.IsTruncated,
		NextMarker:    result.NextMarker,
	}, nil
}

// deleteDistributionCore is the single entry point for distribution deletion.
// It performs all validation: ETag check, DistributionNotDisabled check,
// WAF cleanup, tag cleanup, and invalidation cleanup. This ensures the admin
// handler and HTTP API follow identical deletion semantics.
func (s *CloudFrontService) deleteDistributionCore(ctx context.Context, stores *cloudfrontStores, in DeleteDistributionInput) error {
	if in.Id == "" {
		return invalidArgument("Id is required")
	}

	distribution, err := stores.distributions.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return err
	}

	if err := verifyIfMatch(in.IfMatch, distribution.ETag); err != nil {
		return err
	}

	if distribution.Enabled {
		return awserrors.NewAWSError("DistributionNotDisabled",
			"Distribution must be disabled before deletion", 409)
	}

	certArnForCleanup := ""
	if distribution.DistributionConfig != nil && distribution.DistributionConfig.ViewerCertificate != nil {
		certArnForCleanup = distribution.DistributionConfig.ViewerCertificate.ACMCertificateArn
	}

	if distribution.DistributionConfig != nil && s.wafInvoker != nil {
		webACLId := distribution.DistributionConfig.WebACLId
		if webACLId != "" {
			if err := s.wafInvoker.DisassociateWebACL(webACLId, distribution.ARN); err != nil {
				return fmt.Errorf("failed to remove WAF association: %w", err)
			}
		}
	}

	if distribution.ARN != "" {
		_ = stores.tags.TagStore.Delete(distribution.ARN)
	}

	_ = stores.invalidations.DeleteByDistribution(in.Id)

	if err := stores.distributions.Delete(in.Id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return err
	}

	// Deleting the distribution drops its cached entries.
	if s.distributionServer != nil {
		s.distributionServer.PurgeDistribution(in.Id)
		s.distributionServer.PurgeCertificates()
	}

	if certArnForCleanup != "" && s.acmInvoker != nil {
		if err := s.acmInvoker.UnregisterCertificateUsage(ctx, in.ACMRegion, certArnForCleanup, distribution.ARN); err != nil {
			slog.Error("failed to unregister certificate usage for deleted distribution", "distributionId", in.Id, "error", err)
		}
	}

	return nil
}

// transitionDistributionDeployed asynchronously transitions a distribution
// from InProgress to Deployed, simulating real CloudFront deployment.
func (s *CloudFrontService) transitionDistributionDeployed(stores *cloudfrontStores, distID string) {
	defer func() {
		if r := resilience.RecoverPanic("cloudfront distribution deploy transition"); r != nil {
			slog.Error("panic during distribution status transition",
				"distributionId", distID, "panic", r)
		}
	}()

	time.Sleep(3 * time.Second)

	dist, err := stores.distributions.Get(distID)
	if err != nil {
		slog.Error("failed to get distribution for status transition",
			"distributionId", distID, "error", err)
		return
	}
	dist.Status = "Deployed"
	if err := stores.distributions.UpdateWithLastModified(distID, dist); err != nil {
		slog.Error("failed to persist distribution status transition",
			"distributionId", distID, "error", err)
	}
}

// AssociateDistributionWebACLInput carries the parameters for
// AssociateDistributionWebACL.
type AssociateDistributionWebACLInput struct {
	Id        string
	WebACLArn string
	IfMatch   string
}

// associateDistributionWebACLCore is the single entry point for
// associating a Web ACL with a distribution. The If-Match version check
// is optional for this operation.
func (s *CloudFrontService) associateDistributionWebACLCore(ctx context.Context, stores *cloudfrontStores, in AssociateDistributionWebACLInput) (*cloudfrontstore.Distribution, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if in.WebACLArn == "" {
		return nil, invalidArgument("WebACLArn is required")
	}
	distribution, err := stores.distributions.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}
	if in.IfMatch != "" {
		if err := verifyIfMatch(in.IfMatch, distribution.ETag); err != nil {
			return nil, err
		}
	}
	// A distribution with an attached continuous deployment policy
	// cannot change its Web ACL association; the policy must be deleted
	// first. Re-associating the web ACL that is already in place stays
	// allowed — the continuous deployment quotas page forbids only an
	// association that is the first one for that ACL on the
	// distribution.
	if distribution.DistributionConfig != nil &&
		distribution.DistributionConfig.ContinuousDeploymentPolicyId != "" &&
		distribution.DistributionConfig.WebACLId != in.WebACLArn {
		return nil, awserrors.NewAWSError("InvalidArgument",
			"You cannot associate a web ACL with a distribution that has a continuous deployment policy. Delete the continuous deployment policy first.", 400)
	}
	// This operation models EntityNotFound rather than InvalidWebACLId
	// (which belongs to the distribution create/update family) for a Web
	// ACL that does not exist.
	if s.wafInvoker != nil && !s.wafInvoker.WebACLExists(ctx, in.WebACLArn) {
		return nil, awserrors.NewAWSError("EntityNotFound",
			"The specified Web ACL does not exist: "+in.WebACLArn, 404)
	}
	if distribution.DistributionConfig != nil && distribution.DistributionConfig.WebACLId != in.WebACLArn {
		old := distribution.DistributionConfig.WebACLId
		// Persist the new association first; a failed WAF sync then
		// restores the previous association instead of leaving the
		// external state ahead of the stored one.
		distribution.DistributionConfig.WebACLId = in.WebACLArn
		if err := stores.distributions.UpdateWithLastModified(in.Id, distribution); err != nil {
			return nil, err
		}
		if s.wafInvoker != nil {
			if old != "" {
				if err := s.wafInvoker.DisassociateWebACL(old, distribution.ARN); err != nil {
					distribution.DistributionConfig.WebACLId = old
					_ = stores.distributions.UpdateWithLastModified(in.Id, distribution)
					return nil, fmt.Errorf("failed to remove old WAF association: %w", err)
				}
			}
			if err := s.wafInvoker.AssociateWebACL(in.WebACLArn, distribution.ARN); err != nil {
				// The disassociation above already succeeded, so
				// re-establish the old association.
				if old != "" {
					_ = s.wafInvoker.AssociateWebACL(old, distribution.ARN)
				}
				distribution.DistributionConfig.WebACLId = old
				_ = stores.distributions.UpdateWithLastModified(in.Id, distribution)
				return nil, fmt.Errorf("failed to sync WAF association: %w", err)
			}
		}
	}
	return distribution, nil
}

// disassociateDistributionWebACLCore is the single entry point for
// removing a distribution's Web ACL association.
func (s *CloudFrontService) disassociateDistributionWebACLCore(ctx context.Context, stores *cloudfrontStores, id, ifMatch string) (*cloudfrontstore.Distribution, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	distribution, err := stores.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}
	if ifMatch != "" {
		if err := verifyIfMatch(ifMatch, distribution.ETag); err != nil {
			return nil, err
		}
	}
	// A distribution with an attached continuous deployment policy
	// cannot drop its Web ACL association either; the policy must be
	// deleted first.
	if distribution.DistributionConfig != nil && distribution.DistributionConfig.ContinuousDeploymentPolicyId != "" {
		return nil, awserrors.NewAWSError("InvalidArgument",
			"You cannot disassociate a web ACL from a distribution that has a continuous deployment policy. Delete the continuous deployment policy first.", 400)
	}
	if distribution.DistributionConfig != nil && distribution.DistributionConfig.WebACLId != "" {
		old := distribution.DistributionConfig.WebACLId
		// Persist the cleared association first; a failed WAF call then
		// restores the previous association instead of losing it.
		distribution.DistributionConfig.WebACLId = ""
		if err := stores.distributions.UpdateWithLastModified(id, distribution); err != nil {
			return nil, err
		}
		if s.wafInvoker != nil {
			if err := s.wafInvoker.DisassociateWebACL(old, distribution.ARN); err != nil {
				distribution.DistributionConfig.WebACLId = old
				_ = stores.distributions.UpdateWithLastModified(id, distribution)
				return nil, fmt.Errorf("failed to remove WAF association: %w", err)
			}
		}
	}
	return distribution, nil
}

// CopyDistributionInput carries the parameters for CopyDistribution.
type CopyDistributionInput struct {
	PrimaryId       string
	CallerReference string
	// Enabled is the explicit enabled state for the staging copy. When
	// EnabledProvided is false the AWS default of true applies.
	Enabled         bool
	EnabledProvided bool
	// StagingProvided and Staging carry the Staging header; the only
	// valid value is true.
	StagingProvided bool
	Staging         bool
	IfMatch         string
	ACMRegion       string
}

// copyDistributionCore is the single entry point for creating a staging
// copy of an existing distribution. The primary's configuration is
// deep-copied so the two distributions never share state.
func (s *CloudFrontService) copyDistributionCore(ctx context.Context, stores *cloudfrontStores, in CopyDistributionInput) (*cloudfrontstore.Distribution, error) {
	if err := requireID(in.PrimaryId); err != nil {
		return nil, err
	}
	if in.CallerReference == "" {
		return nil, invalidArgument("CallerReference is required")
	}
	if in.StagingProvided && !in.Staging {
		return nil, invalidArgument("Staging must be true")
	}
	primary, err := stores.distributions.Get(in.PrimaryId)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}
	if in.IfMatch != "" {
		if err := verifyIfMatch(in.IfMatch, primary.ETag); err != nil {
			return nil, err
		}
	}
	if primary.DistributionConfig == nil {
		return nil, invalidArgument("The primary distribution has no configuration to copy")
	}

	config, err := cloneDistributionConfig(primary.DistributionConfig)
	if err != nil {
		return nil, err
	}
	config.CallerReference = in.CallerReference
	if in.EnabledProvided {
		config.Enabled = in.Enabled
	} else {
		config.Enabled = true
	}
	config.Staging = true

	result, err := s.createDistributionCore(ctx, stores, CreateDistributionInput{
		CallerReference: in.CallerReference,
		Config:          config,
		ACMRegion:       in.ACMRegion,
	})
	if err != nil {
		return nil, err
	}
	// The store derives the top-level Staging flag from the copied
	// configuration, so no post-create adjustment is needed here.
	return result.Distribution, nil
}

// cloneDistributionConfig deep-copies a distribution configuration so a
// staging copy never shares pointers with its primary distribution.
func cloneDistributionConfig(config *cloudfrontstore.DistributionConfig) (*cloudfrontstore.DistributionConfig, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise distribution config for copy: %w", err)
	}
	var clone cloudfrontstore.DistributionConfig
	if err := json.Unmarshal(data, &clone); err != nil {
		return nil, fmt.Errorf("failed to deserialise distribution config for copy: %w", err)
	}
	return &clone, nil
}

// getDistributionCore is the single entry point for fetching a
// distribution by ID, mapping a missing record to the modelled
// NoSuchDistribution error.
func (s *CloudFrontService) getDistributionCore(stores *cloudfrontStores, id string) (*cloudfrontstore.Distribution, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	distribution, err := stores.distributions.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchDistribution", "Distribution not found", 404)
		}
		return nil, err
	}
	return distribution, nil
}

// DistributionDetail holds the store-derived response parts that every
// distribution-returning operation serialises: the count of in-progress
// invalidation batches and the ActiveTrustedSigners /
// ActiveTrustedKeyGroups shapes, whose key-pair IDs resolve through the
// key group store.
type DistributionDetail struct {
	InProgressInvalidations int
	ActiveSigners           map[string]interface{}
	ActiveKeyGroups         map[string]interface{}
}

// distributionDetailCore collects the store-derived parts of a
// distribution response so that every plane computes them through one
// code path.
func (s *CloudFrontService) distributionDetailCore(stores *cloudfrontStores, d *cloudfrontstore.Distribution) DistributionDetail {
	inProgressCount, _ := stores.invalidations.CountInProgress(d.ID)
	return DistributionDetail{
		InProgressInvalidations: inProgressCount,
		ActiveSigners:           computeActiveTrustedSigners(d),
		ActiveKeyGroups:         computeActiveTrustedKeyGroups(d, stores),
	}
}

// computeActiveTrustedKeyGroups inspects all cache behaviours for
// TrustedKeyGroups with Enabled=true and produces the
// ActiveTrustedKeyGroups output shape. Each key group's PublicKey IDs are
// resolved from the KeyGroup store to populate KeyPairIds.
func computeActiveTrustedKeyGroups(d *cloudfrontstore.Distribution, stores *cloudfrontStores) map[string]interface{} {
	kgIDs := collectTrustedKeyGroupIDs(d)
	if len(kgIDs) == 0 {
		return map[string]interface{}{"Enabled": false, "Quantity": 0}
	}
	items := make([]interface{}, 0, len(kgIDs))
	for kgID := range kgIDs {
		keyPairIds := map[string]interface{}{"Quantity": 0}
		if stores != nil {
			if kg, err := stores.keyGroups.Get(kgID); err == nil && len(kg.KeyGroupConfig.Items) > 0 {
				kpItems := make([]interface{}, len(kg.KeyGroupConfig.Items))
				for i, kp := range kg.KeyGroupConfig.Items {
					kpItems[i] = kp
				}
				keyPairIds = map[string]interface{}{
					"Quantity": len(kg.KeyGroupConfig.Items),
					"Items":    protocol.XMLElements{ElementName: "KeyPairId", Items: kpItems},
				}
			}
		}
		items = append(items, map[string]interface{}{
			"KeyGroupId": kgID,
			"KeyPairIds": keyPairIds,
		})
	}
	return map[string]interface{}{
		"Enabled":  true,
		"Quantity": len(kgIDs),
		"Items":    protocol.XMLElements{ElementName: "KeyGroup", Items: items},
	}
}

// ListDistributionsByWebACLIdInput carries the parameters for
// ListDistributionsByWebACLId.
type ListDistributionsByWebACLIdInput struct {
	WebACLId string
}

// listDistributionsByWebACLIdCore returns every distribution whose
// configuration references the given Web ACL.
func (s *CloudFrontService) listDistributionsByWebACLIdCore(ctx context.Context, stores *cloudfrontStores, in ListDistributionsByWebACLIdInput) ([]*cloudfrontstore.Distribution, error) {
	// This operation models InvalidWebACLId, returned when the specified
	// Web ACL does not exist; the listing must not silently succeed with
	// an empty result for an unknown ACL.
	if s.wafInvoker != nil && !s.wafInvoker.WebACLExists(ctx, in.WebACLId) {
		return nil, awserrors.NewAWSError("InvalidWebACLId", "The specified Web ACL does not exist: "+in.WebACLId, 400)
	}

	// The association match must run over every stored distribution: the
	// shared iterator replaces a zero MaxItems with the platform default
	// page size, so a single unbounded-looking List call would silently
	// stop after the first page and miss associations whose records sit
	// beyond it. Page with NextMarker until the store is exhausted.
	var matched []*cloudfrontstore.Distribution
	scanMarker := ""
	for {
		page, err := stores.distributions.List(scanMarker, 0)
		if err != nil {
			return nil, err
		}
		for _, d := range page.Distributions {
			if d.DistributionConfig != nil && d.DistributionConfig.WebACLId == in.WebACLId {
				matched = append(matched, d)
			}
		}
		if !page.IsTruncated || page.NextMarker == "" {
			break
		}
		scanMarker = page.NextMarker
	}
	return matched, nil
}

// collectDistributions traverses every distribution page by page and
// returns the IDs of all distributions for which fn returns true.
func collectDistributions(stores *cloudfrontStores, fn func(*cloudfrontstore.Distribution) bool) ([]string, error) {
	marker := ""
	var ids []string
	for {
		result, err := stores.distributions.List(marker, cloudfrontstore.DefaultListMaxItems)
		if err != nil {
			return nil, err
		}
		for _, dist := range result.Distributions {
			if fn(dist) {
				ids = append(ids, dist.ID)
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			return ids, nil
		}
		marker = result.NextMarker
	}
}

// requireReferenceIdCore enforces the required reference member shared by
// the list-distributions-by-reference operations.
func requireReferenceIdCore(id, param string) error {
	if id == "" {
		return awserrors.NewAWSError("InvalidArgument", param+" is required", 400)
	}
	return nil
}
