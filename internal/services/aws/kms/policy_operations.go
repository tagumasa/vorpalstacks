package kms

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetKeyPolicy retrieves the key policy for a specified key.
// AWS only supports the policy name "default" for customer managed keys;
// any other PolicyName value returns a ValidationException.
func (s *KMSService) GetKeyPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GetKeyPolicy", nil)
	if err != nil {
		return nil, err
	}

	policy, err := s.getKeyPolicyCore(stores, key, request.GetStringParam(req.Parameters, "PolicyName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy":     policy.PolicyDocument,
		"PolicyName": policy.PolicyName,
	}, nil
}

// PutKeyPolicy sets the key policy for a specified key.
// The policy must be a valid JSON policy document.
// AWS only supports the policy name "default"; any other value returns
// ValidationException. The BypassPolicyLockoutSafetyCheck parameter, when
// false (default), enforces that the policy must not lock out the root
// principal. The previous implementation silently ignored this flag.
func (s *KMSService) PutKeyPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "PutKeyPolicy", nil)
	if err != nil {
		return nil, err
	}

	bypassPolicyLockoutSafetyCheck := false
	if v, ok := req.Parameters["BypassPolicyLockoutSafetyCheck"]; ok {
		if b, ok := v.(bool); ok {
			bypassPolicyLockoutSafetyCheck = b
		}
	}

	if err := s.putKeyPolicyCore(stores, key,
		request.GetStringParam(req.Parameters, "PolicyName"),
		request.GetStringParam(req.Parameters, "Policy"),
		bypassPolicyLockoutSafetyCheck,
		reqCtx.GetAccountID()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// validatePolicyDoesNotLockOutRoot verifies that the supplied policy
// grants the account root principal kms:PutKeyPolicy (or kms:*). AWS
// rejects policies that lock out the root principal unless
// BypassPolicyLockoutSafetyCheck is explicitly set to true.
//
// The accountID parameter comes from the request context, not from the
// policy's Resource ARN. This avoids the wildcard-resource pitfall:
// when Resource is ["*"], extracting the account from the ARN yields
// nothing, and the previous fallback to "000000000000" caused
// legitimate root-only policies to be rejected.
//
// The evaluator handles NotAction, NotPrincipal, Conditions, and Deny
// statements that could implicitly lock out root. The root principal
// ARN is constructed as "arn:vorpalstacks:iam::<accountID>:root".
func validatePolicyDoesNotLockOutRoot(policyDocument string, accountID string) error {
	parsedPolicy, err := policy.ParseDocument(policyDocument)
	if err != nil {
		return ErrMalformedPolicy
	}
	if len(parsedPolicy.Statement) == 0 {
		return ErrMalformedPolicy
	}

	// Construct the root principal from the caller's account ID. The
	// previous implementation extracted the account ID from the policy's
	// Resource ARN, but that fails when Resource is "*" or when the
	// resource belongs to a different account. Using the request
	// context's account ID ensures the lockout check always evaluates
	// against the real account root principal.
	rootPrincipal := fmt.Sprintf("arn:vorpalstacks:iam::%s:root", accountID)

	evaluator := policy.NewPolicyEvaluator()
	ctx := &policy.EvaluationContext{
		Principal: rootPrincipal,
		Action:    "kms:PutKeyPolicy",
		// Resource is set to "*" so that any resource-scoped statement
		// matches; we are testing whether the root principal has the
		// PutKeyPolicy permission, not whether a specific key matches.
		Resource: "*",
	}

	// Test with root principal for kms:PutKeyPolicy.
	decision := evaluator.Evaluate(ctx, []*policy.Document{parsedPolicy})
	if decision.Effect == policy.DecisionEffectAllow {
		// Root can PutKeyPolicy — not locked out.
		return nil
	}

	// Also test with "*" principal wildcard. Some policies use Principal "*"
	// which implicitly includes root.
	ctx.Principal = "*"
	decision = evaluator.Evaluate(ctx, []*policy.Document{parsedPolicy})
	if decision.Effect == policy.DecisionEffectAllow {
		return nil
	}

	return ErrMalformedPolicy
}

// ListKeyPolicies retrieves the names of all key policies attached to a key.
func (s *KMSService) ListKeyPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "ListKeyPolicies", nil)
	if err != nil {
		return nil, err
	}

	result, err := s.listKeyPoliciesCore(stores, key,
		pagination.GetMarker(req.Parameters),
		pagination.GetMaxItems(req.Parameters, 100))
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"PolicyNames": result.Items,
		"Truncated":   result.IsTruncated,
	}
	if result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}
