package cloudfront

import (
	"fmt"
	"log/slog"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/resilience"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// Core functions for the cache-invalidation operations. The HTTP API
// handlers in invalidation_operations.go are thin adapters that receive
// the global store bundle, parse the wire request, and delegate all
// validation and persistence here.

// CreateInvalidationInput carries the parameters for CreateInvalidation in
// a format independent of the wire protocol.
type CreateInvalidationInput struct {
	Id              string
	CallerReference string
	Paths           []string
	// DeclaredQuantity is the Paths.Quantity value from the request; the
	// batch consistency check compares it against the parsed path items.
	DeclaredQuantity int
}

// GetInvalidationInput carries the identifiers for GetInvalidation.
type GetInvalidationInput struct {
	Id             string
	InvalidationId string
}

// ListInvalidationsInput carries the pagination parameters for
// ListInvalidations.
type ListInvalidationsInput struct {
	Id       string
	Marker   string
	MaxItems int
}

// createInvalidationCore is the single entry point for creating a cache
// invalidation: it verifies the distribution exists, validates the
// invalidation batch, persists it, and starts the asynchronous status
// transition to Completed.
func (s *CloudFrontService) createInvalidationCore(stores *cloudfrontStores, in CreateInvalidationInput) (*cloudfrontstore.Invalidation, error) {
	if in.Id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Required parameter Id is missing.", 400)
	}

	if _, err := stores.distributions.Get(in.Id); err != nil {
		return nil, awserrors.NewAWSError("NoSuchDistribution", fmt.Sprintf("The specified distribution does not exist: %s", in.Id), 404)
	}

	if in.CallerReference == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "CallerReference is required.", 400)
	}

	if in.DeclaredQuantity != len(in.Paths) {
		return nil, awserrors.NewAWSError("InconsistentQuantities",
			fmt.Sprintf("The Quantity value (%d) does not match the number of path items (%d)", in.DeclaredQuantity, len(in.Paths)), 400)
	}

	if len(in.Paths) == 0 {
		return nil, awserrors.NewAWSError("InvalidArgument", "At least one path is required.", 400)
	}
	if len(in.Paths) > cloudfrontstore.MaxInvalidationPathsPerRequest {
		return nil, awserrors.NewAWSError("BatchTooLarge", "Invalidation batch specified is too large.", 413)
	}

	inv, err := stores.invalidations.Create(in.Id, in.CallerReference, in.Paths)
	if err != nil {
		return nil, awserrors.NewInternalFailureException(fmt.Sprintf("Failed to create invalidation: %v", err))
	}

	// The invalidation takes effect on the live edge cache immediately.
	if s.distributionServer != nil {
		s.distributionServer.InvalidatePaths(in.Id, in.Paths)
	}

	go s.transitionInvalidation(stores, inv)

	return inv, nil
}

// transitionInvalidation asynchronously transitions an invalidation from
// InProgress to Completed, simulating the real CloudFront edge propagation.
func (s *CloudFrontService) transitionInvalidation(stores *cloudfrontStores, inv *cloudfrontstore.Invalidation) {
	defer func() {
		if r := resilience.RecoverPanic("cloudfront invalidation status transition"); r != nil {
			slog.Error("panic during invalidation status transition",
				"invalidationId", inv.ID, "distributionId", inv.DistributionID, "panic", r)
		}
	}()

	time.Sleep(2 * time.Second)

	inv.Status = "Completed"
	if err := stores.invalidations.Update(inv); err != nil {
		slog.Error("failed to persist invalidation status transition",
			"invalidationId", inv.ID, "distributionId", inv.DistributionID, "error", err)
	}
}

// getInvalidationCore is the single entry point for fetching an
// invalidation, mapping a missing distribution or invalidation to the
// modelled NoSuchDistribution and NoSuchInvalidation errors.
func (s *CloudFrontService) getInvalidationCore(stores *cloudfrontStores, in GetInvalidationInput) (*cloudfrontstore.Invalidation, error) {
	if in.Id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Required parameter Id is missing.", 400)
	}

	if in.InvalidationId == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Required parameter invalidationId is missing.", 400)
	}

	if _, err := stores.distributions.Get(in.Id); err != nil {
		return nil, awserrors.NewAWSError("NoSuchDistribution", fmt.Sprintf("The specified distribution does not exist: %s", in.Id), 404)
	}

	inv, err := stores.invalidations.Get(in.Id, in.InvalidationId)
	if err != nil {
		return nil, awserrors.NewAWSError("NoSuchInvalidation", fmt.Sprintf("The specified invalidation does not exist: %s", in.InvalidationId), 404)
	}

	return inv, nil
}

// listInvalidationsCore is the single entry point for listing the
// invalidations of a distribution.
func (s *CloudFrontService) listInvalidationsCore(stores *cloudfrontStores, in ListInvalidationsInput) (*cloudfrontstore.InvalidationListResult, error) {
	if in.Id == "" {
		return nil, awserrors.NewAWSError("InvalidArgument", "Required parameter Id is missing.", 400)
	}

	if _, err := stores.distributions.Get(in.Id); err != nil {
		return nil, awserrors.NewAWSError("NoSuchDistribution", fmt.Sprintf("The specified distribution does not exist: %s", in.Id), 404)
	}

	result, err := stores.invalidations.List(in.Id, in.Marker, in.MaxItems)
	if err != nil {
		return nil, awserrors.NewInternalFailureException(fmt.Sprintf("Failed to list invalidations: %v", err))
	}
	return result, nil
}
