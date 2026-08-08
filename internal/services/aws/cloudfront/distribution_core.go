package cloudfront

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/resilience"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
	"vorpalstacks/internal/utils/aws/types"
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

// CreateDistributionResult is the transport-agnostic result of
// createDistributionCore.
type CreateDistributionResult struct {
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

	for _, origin := range in.Config.Origins.Items {
		if origin.CustomOriginConfig != nil {
			opp := origin.CustomOriginConfig.OriginProtocolPolicy
			if opp == "" {
				return nil, invalidArgument("OriginProtocolPolicy is required when CustomOriginConfig is specified")
			}
			if !isValidOriginProtocolPolicy(opp) {
				return nil, invalidArgument("Invalid OriginProtocolPolicy: " + opp)
			}
		}
	}
	if in.Config.PriceClass != "" && !isValidPriceClass(in.Config.PriceClass) {
		return nil, invalidArgument("Invalid PriceClass: " + in.Config.PriceClass)
	}
	if in.Config.HttpVersion != "" && !isValidHttpVersion(in.Config.HttpVersion) {
		return nil, invalidArgument("Invalid HttpVersion: " + in.Config.HttpVersion)
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
	config := buildDefaultDistributionConfig(in.CallerReference, in.Comment, in.Enabled, in.OriginID, in.OriginDomain)
	return s.createDistributionCore(ctx, stores, CreateDistributionInput{
		CallerReference: in.CallerReference,
		Config:          config,
		ACMRegion:       in.ACMRegion,
	})
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
	if in.MaxItems <= 0 {
		in.MaxItems = 100
	}
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

	if in.IfMatch == "" {
		return awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}
	if in.IfMatch != "*" && distribution.ETag != in.IfMatch {
		return awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
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
