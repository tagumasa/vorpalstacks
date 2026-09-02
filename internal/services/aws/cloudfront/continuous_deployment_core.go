package cloudfront

import (
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// Core functions for the continuous deployment policy operations. The
// HTTP API handlers in continuous_deployment_operations.go are thin
// adapters that parse the wire request, delegate all validation and
// persistence here, and serialise the result.

// CreateContinuousDeploymentPolicyInput carries the parameters for
// CreateContinuousDeploymentPolicy.
type CreateContinuousDeploymentPolicyInput struct {
	Config *cloudfrontstore.ContinuousDeploymentPolicyConfig
}

// UpdateContinuousDeploymentPolicyInput carries the parameters for
// UpdateContinuousDeploymentPolicy.
type UpdateContinuousDeploymentPolicyInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.ContinuousDeploymentPolicyConfig
}

// ListContinuousDeploymentPoliciesInput carries the pagination parameters
// for ListContinuousDeploymentPolicies.
type ListContinuousDeploymentPoliciesInput struct {
	Marker   string
	MaxItems int
}

// createContinuousDeploymentPolicyCore is the single entry point for
// creating a continuous deployment policy: it validates the
// configuration, enforces the account quota and staging-distribution
// exclusivity, and persists the policy.
func (s *CloudFrontService) createContinuousDeploymentPolicyCore(stores *cloudfrontStores, in CreateContinuousDeploymentPolicyInput) (*cloudfrontstore.ContinuousDeploymentPolicy, error) {
	if in.Config == nil {
		return nil, invalidArgument("ContinuousDeploymentPolicyConfig is required")
	}
	if err := s.validateContinuousDeploymentPolicyConfig(stores, in.Config, ""); err != nil {
		return nil, err
	}
	if stores.deploymentPolicies.Count() >= cloudfrontstore.MaxContinuousDeploymentPolicies {
		return nil, awserrors.NewAWSError("TooManyContinuousDeploymentPolicies",
			fmt.Sprintf("You have reached the maximum number of continuous deployment policies: %d", cloudfrontstore.MaxContinuousDeploymentPolicies), 400)
	}
	return stores.deploymentPolicies.Create(in.Config)
}

// getContinuousDeploymentPolicyCore is the single entry point for
// fetching a policy, mapping a missing policy to the modelled
// NoSuchContinuousDeploymentPolicy error.
func (s *CloudFrontService) getContinuousDeploymentPolicyCore(stores *cloudfrontStores, id string) (*cloudfrontstore.ContinuousDeploymentPolicy, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	policy, err := stores.deploymentPolicies.Get(id)
	if err != nil {
		return nil, awserrors.NewAWSError("NoSuchContinuousDeploymentPolicy",
			fmt.Sprintf("The specified continuous deployment policy does not exist: %s", id), 404)
	}
	return policy, nil
}

// updateContinuousDeploymentPolicyCore is the single entry point for
// updating a policy. The If-Match version check and the
// staging-distribution exclusivity re-check run before persistence.
func (s *CloudFrontService) updateContinuousDeploymentPolicyCore(stores *cloudfrontStores, in UpdateContinuousDeploymentPolicyInput) (*cloudfrontstore.ContinuousDeploymentPolicy, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if in.Config == nil {
		return nil, invalidArgument("ContinuousDeploymentPolicyConfig is required")
	}
	existing, err := stores.deploymentPolicies.Get(in.Id)
	if err != nil {
		return nil, awserrors.NewAWSError("NoSuchContinuousDeploymentPolicy",
			fmt.Sprintf("The specified continuous deployment policy does not exist: %s", in.Id), 404)
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if err := s.validateContinuousDeploymentPolicyConfig(stores, in.Config, in.Id); err != nil {
		return nil, err
	}
	return stores.deploymentPolicies.Update(in.Id, in.Config)
}

// deleteContinuousDeploymentPolicyCore is the single entry point for
// deleting a policy. A policy still attached to a distribution cannot be
// deleted (ContinuousDeploymentPolicyInUse).
func (s *CloudFrontService) deleteContinuousDeploymentPolicyCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	policy, err := stores.deploymentPolicies.Get(id)
	if err != nil {
		return awserrors.NewAWSError("NoSuchContinuousDeploymentPolicy",
			fmt.Sprintf("The specified continuous deployment policy does not exist: %s", id), 404)
	}
	if err := verifyIfMatch(ifMatch, policy.ETag); err != nil {
		return err
	}
	referenced, err := scanDistributions(stores, func(dist *cloudfrontstore.Distribution) bool {
		return dist.DistributionConfig != nil && dist.DistributionConfig.ContinuousDeploymentPolicyId == id
	})
	if err != nil {
		return err
	}
	if referenced {
		return awserrors.NewAWSError("ContinuousDeploymentPolicyInUse",
			"The continuous deployment policy is attached to one or more distributions", 409)
	}
	return stores.deploymentPolicies.Delete(id)
}

// listContinuousDeploymentPoliciesCore is the single entry point for
// listing continuous deployment policies.
func (s *CloudFrontService) listContinuousDeploymentPoliciesCore(stores *cloudfrontStores, in ListContinuousDeploymentPoliciesInput) (*cloudfrontstore.ContinuousDeploymentPolicyListResult, error) {
	return stores.deploymentPolicies.List(in.Marker, resolveListMaxItems(in.MaxItems))
}

// validateContinuousDeploymentPolicyConfig applies the documented
// configuration rules: staging DNS names must belong to existing staging
// distributions, each staging distribution backs at most one policy, the
// traffic type selects exactly one traffic configuration, the weight is
// bounded by the 15 percent quota, session stickiness TTLs live in
// 300-3600 seconds with idle not exceeding maximum, and single-header
// header names must carry the aws-cf-cd- prefix.
func (s *CloudFrontService) validateContinuousDeploymentPolicyConfig(stores *cloudfrontStores, config *cloudfrontstore.ContinuousDeploymentPolicyConfig, selfID string) error {
	names := config.StagingDistributionDnsNames
	if names == nil || len(names.Items) == 0 {
		return invalidArgument("StagingDistributionDnsNames must contain at least one staging distribution domain name")
	}
	if names.Quantity != len(names.Items) {
		return awserrors.NewAWSError("InconsistentQuantities",
			fmt.Sprintf("The Quantity value (%d) does not match the number of DNS name items (%d)", names.Quantity, len(names.Items)), 400)
	}
	for _, name := range names.Items {
		staging, err := stores.distributions.GetByDomainName(name)
		if err != nil || staging == nil {
			return invalidArgument(fmt.Sprintf("The specified staging distribution does not exist: %s", name))
		}
		if !staging.Staging {
			return invalidArgument(fmt.Sprintf("The specified distribution is not a staging distribution: %s", name))
		}
		inUse, err := s.stagingNameUsedByPolicy(stores, name, selfID)
		if err != nil {
			return err
		}
		if inUse {
			return awserrors.NewAWSError("StagingDistributionInUse",
				fmt.Sprintf("The staging distribution is already used by another continuous deployment policy: %s", name), 409)
		}
	}

	if config.TrafficConfig == nil {
		return nil
	}
	traffic := config.TrafficConfig
	switch traffic.Type {
	case "SingleWeight":
		if traffic.SingleHeaderConfig != nil {
			return invalidArgument("SingleHeaderConfig must not be present when the traffic type is SingleWeight")
		}
		weight := traffic.SingleWeightConfig
		if weight == nil {
			return invalidArgument("SingleWeightConfig is required when the traffic type is SingleWeight")
		}
		if weight.Weight < 0 || weight.Weight > cloudfrontstore.MaxContinuousDeploymentWeight {
			return invalidArgument(fmt.Sprintf("Weight must be between 0 and %v", cloudfrontstore.MaxContinuousDeploymentWeight))
		}
		if stickiness := weight.SessionStickinessConfig; stickiness != nil {
			if stickiness.IdleTTL < cloudfrontstore.MinSessionStickinessTTLSeconds || stickiness.IdleTTL > cloudfrontstore.MaxSessionStickinessTTLSeconds {
				return invalidArgument(fmt.Sprintf("IdleTTL must be between %d and %d seconds",
					cloudfrontstore.MinSessionStickinessTTLSeconds, cloudfrontstore.MaxSessionStickinessTTLSeconds))
			}
			if stickiness.MaximumTTL < cloudfrontstore.MinSessionStickinessTTLSeconds || stickiness.MaximumTTL > cloudfrontstore.MaxSessionStickinessTTLSeconds {
				return invalidArgument(fmt.Sprintf("MaximumTTL must be between %d and %d seconds",
					cloudfrontstore.MinSessionStickinessTTLSeconds, cloudfrontstore.MaxSessionStickinessTTLSeconds))
			}
			if stickiness.IdleTTL > stickiness.MaximumTTL {
				return invalidArgument("IdleTTL must be less than or equal to MaximumTTL")
			}
		}
	case "SingleHeader":
		if traffic.SingleWeightConfig != nil {
			return invalidArgument("SingleWeightConfig must not be present when the traffic type is SingleHeader")
		}
		header := traffic.SingleHeaderConfig
		if header == nil {
			return invalidArgument("SingleHeaderConfig is required when the traffic type is SingleHeader")
		}
		if header.Header == "" || header.Value == "" {
			return invalidArgument("SingleHeaderConfig requires a header name and value")
		}
		if !strings.HasPrefix(header.Header, "aws-cf-cd-") {
			return invalidArgument("The single-header configuration header name must contain the prefix aws-cf-cd-")
		}
	default:
		return invalidArgument("TrafficConfig Type must be SingleWeight or SingleHeader")
	}
	return nil
}

// stagingNameUsedByPolicy reports whether a staging distribution domain
// name is already referenced by another continuous deployment policy.
func (s *CloudFrontService) stagingNameUsedByPolicy(stores *cloudfrontStores, name, selfID string) (bool, error) {
	marker := ""
	for {
		result, err := stores.deploymentPolicies.List(marker, cloudfrontstore.DefaultListMaxItems)
		if err != nil {
			return false, err
		}
		for _, policy := range result.Policies {
			if policy.ID == selfID || policy.ContinuousDeploymentPolicyConfig == nil {
				continue
			}
			names := policy.ContinuousDeploymentPolicyConfig.StagingDistributionDnsNames
			if names == nil {
				continue
			}
			for _, item := range names.Items {
				if strings.EqualFold(item, name) {
					return true, nil
				}
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			return false, nil
		}
		marker = result.NextMarker
	}
}

// validateDistributionPolicyReference checks a distribution
// configuration's ContinuousDeploymentPolicyId: the policy must exist and
// a staging distribution cannot itself carry a policy.
func (s *CloudFrontService) validateDistributionPolicyReference(stores *cloudfrontStores, config *cloudfrontstore.DistributionConfig) error {
	if config.ContinuousDeploymentPolicyId == "" {
		return nil
	}
	if config.Staging {
		return invalidArgument("A staging distribution cannot have a continuous deployment policy")
	}
	if _, err := stores.deploymentPolicies.Get(config.ContinuousDeploymentPolicyId); err != nil {
		return awserrors.NewAWSError("NoSuchContinuousDeploymentPolicy",
			fmt.Sprintf("The specified continuous deployment policy does not exist: %s", config.ContinuousDeploymentPolicyId), 404)
	}
	return nil
}
