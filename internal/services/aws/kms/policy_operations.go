package kms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	if policyName == "" {
		policyName = "default"
	}
	// AWS: only "default" is a valid policy name for KMS keys.
	if policyName != "default" {
		return nil, ErrValidation
	}

	policy, err := stores.keyPolicies.Get(key.KeyID, policyName)
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

	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	if policyName == "" {
		policyName = "default"
	}
	if policyName != "default" {
		return nil, ErrValidation
	}

	policyDocument := request.GetStringParam(req.Parameters, "Policy")
	if policyDocument == "" {
		return nil, ErrMalformedPolicy
	}

	var js interface{}
	if err := json.Unmarshal([]byte(policyDocument), &js); err != nil {
		return nil, ErrMalformedPolicy
	}

	bypassPolicyLockoutSafetyCheck := false
	if v, ok := req.Parameters["BypassPolicyLockoutSafetyCheck"]; ok {
		if b, ok := v.(bool); ok {
			bypassPolicyLockoutSafetyCheck = b
		}
	}

	if !bypassPolicyLockoutSafetyCheck {
		if err := validatePolicyDoesNotLockOutRoot(policyDocument); err != nil {
			return nil, err
		}
	}

	if err := stores.keyPolicies.Put(key.KeyID, policyName, policyDocument); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// validatePolicyDoesNotLockOutRoot verifies that the supplied policy
// grants the account root principal kms:PutKeyPolicy (or kms:*). AWS
// rejects policies that lock out the root principal unless
// BypassPolicyLockoutSafetyCheck=true.
//
// The check uses the full IAM PolicyEvaluator (not manual JSON parsing)
// to correctly handle conditions, NotPrincipal, NotAction, and
// Deny statements that could implicitly lock out root. The root
// principal ARN is constructed from the policy's resource account.
func validatePolicyDoesNotLockOutRoot(policyDocument string) error {
	parsedPolicy, err := policy.ParseDocument(policyDocument)
	if err != nil {
		return ErrMalformedPolicy
	}
	if len(parsedPolicy.Statement) == 0 {
		return ErrMalformedPolicy
	}

	// Extract the account ID from the first resource ARN in the policy.
	// The root principal ARN is "arn:vorpalstacks:iam::<account>:root".
	accountID := "000000000000"
	for _, stmt := range parsedPolicy.Statement {
		for _, res := range stmt.Resource {
			if id := extractAccountFromARN(res); id != "" {
				accountID = id
				break
			}
		}
		if accountID != "000000000000" {
			break
		}
	}
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

// extractAccountFromARN extracts the account ID from an ARN string.
func extractAccountFromARN(arnStr string) string {
	parts := strings.Split(arnStr, ":")
	if len(parts) >= 5 {
		return parts[4]
	}
	return ""
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

	policies, err := stores.keyPolicies.List(key.KeyID)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	result := pagination.PaginateSlice(policies, marker, maxItems, func(p string) string {
		return p
	})

	response := map[string]interface{}{
		"PolicyNames": result.Items,
		"Truncated":   result.IsTruncated,
	}
	if result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}
