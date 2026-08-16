package cloudfront

import (
	awserrors "vorpalstacks/internal/common/errors"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// CreateResponseHeadersPolicyInput carries the parsed response headers
// policy configuration.
type CreateResponseHeadersPolicyInput struct {
	Config *cloudfrontstore.ResponseHeadersPolicyConfig
}

// UpdateResponseHeadersPolicyInput carries the parameters for updating a
// response headers policy.
type UpdateResponseHeadersPolicyInput struct {
	Id      string
	IfMatch string
	Config  *cloudfrontstore.ResponseHeadersPolicyConfig
}

// ResponseHeadersPoliciesResult is the transport-agnostic list result for
// response headers policies, already filtered by the requested policy type.
type ResponseHeadersPoliciesResult struct {
	Policies          []*cloudfrontstore.ResponseHeadersPolicy
	EffectiveMaxItems int
	NextMarker        string
}

func (s *CloudFrontService) createResponseHeadersPolicyCore(stores *cloudfrontStores, in CreateResponseHeadersPolicyInput) (*cloudfrontstore.ResponseHeadersPolicy, error) {
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	existing, _ := stores.responseHeadersPolicies.GetByName(in.Config.Name)
	if existing != nil {
		return nil, awserrors.NewAWSError("ResponseHeadersPolicyAlreadyExists", "Response headers policy with this name already exists", 409)
	}
	return stores.responseHeadersPolicies.Create(in.Config)
}

func (s *CloudFrontService) getResponseHeadersPolicyCore(stores *cloudfrontStores, id string) (*cloudfrontstore.ResponseHeadersPolicy, error) {
	if err := requireID(id); err != nil {
		return nil, err
	}
	policy, err := stores.responseHeadersPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}
	return policy, nil
}

func (s *CloudFrontService) updateResponseHeadersPolicyCore(stores *cloudfrontStores, in UpdateResponseHeadersPolicyInput) (*cloudfrontstore.ResponseHeadersPolicy, error) {
	if err := requireID(in.Id); err != nil {
		return nil, err
	}
	if in.Config == nil || in.Config.Name == "" {
		return nil, invalidArgument("Name is required")
	}
	existing, err := stores.responseHeadersPolicies.Get(in.Id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}
	if err := verifyIfMatch(in.IfMatch, existing.ETag); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(existing.Name, in.Config.Name,
		func(name string) bool { dup, _ := stores.responseHeadersPolicies.GetByName(name); return dup != nil },
		awserrors.NewAWSError("ResponseHeadersPolicyAlreadyExists", "Response headers policy with this name already exists", 409)); err != nil {
		return nil, err
	}
	policy, err := stores.responseHeadersPolicies.Update(in.Id, in.Config)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return nil, err
	}
	return policy, nil
}

func (s *CloudFrontService) deleteResponseHeadersPolicyCore(stores *cloudfrontStores, id, ifMatch string) error {
	if err := requireID(id); err != nil {
		return err
	}
	existing, err := stores.responseHeadersPolicies.Get(id)
	if err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return err
	}
	if err := verifyIfMatch(ifMatch, existing.ETag); err != nil {
		return err
	}
	attached, err := isResponseHeadersPolicyAttached(stores, id)
	if err != nil {
		return awserrors.NewAWSError("InternalError", "Failed to check response headers policy usage: "+err.Error(), 500)
	}
	if attached {
		return awserrors.NewAWSError("ResponseHeadersPolicyInUse",
			"Cannot delete this response headers policy because it is attached to one or more distributions", 409)
	}
	if err := stores.responseHeadersPolicies.Delete(id); err != nil {
		if cloudfrontstore.IsNotFound(err) {
			return awserrors.NewAWSError("NoSuchResponseHeadersPolicy", "Response headers policy not found", 404)
		}
		return err
	}
	return nil
}

func (s *CloudFrontService) listResponseHeadersPoliciesCore(stores *cloudfrontStores, in ListPoliciesInput) (*ResponseHeadersPoliciesResult, error) {
	if in.TypeFilter != "" && !isValidPolicyListType(in.TypeFilter) {
		return nil, invalidArgument("Invalid Type: " + in.TypeFilter)
	}
	maxItems := resolveListMaxItems(in.MaxItems)
	result, err := stores.responseHeadersPolicies.List(in.Marker, maxItems)
	if err != nil {
		return nil, err
	}
	policies := make([]*cloudfrontstore.ResponseHeadersPolicy, 0, len(result.ResponseHeadersPolicies))
	for _, p := range result.ResponseHeadersPolicies {
		if policyMatchesListType(p.IsManaged, in.TypeFilter) {
			policies = append(policies, p)
		}
	}
	return &ResponseHeadersPoliciesResult{
		Policies:          policies,
		EffectiveMaxItems: maxItems,
		NextMarker:        result.NextMarker,
	}, nil
}
