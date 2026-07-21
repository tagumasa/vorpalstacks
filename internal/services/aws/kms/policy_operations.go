package kms

import (
	"context"
	"encoding/json"
	"strings"

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

// validatePolicyDoesNotLockOutRoot performs a best-effort check that the
// supplied policy grants the account root principal kms:PutKeyPolicy (or
// kms:*). AWS rejects policies that lock out the root principal unless
// BypassPolicyLockoutSafetyCheck=true. A full IAM policy evaluator is out
// of scope for this implementation; we reject the clearly-broken case of
// a Statement list with no root-permitting statement.
func validatePolicyDoesNotLockOutRoot(policyDocument string) error {
	var doc struct {
		Statement []struct {
			Principal interface{} `json:"Principal"`
			Action    interface{} `json:"Action"`
			Effect    string      `json:"Effect"`
		} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policyDocument), &doc); err != nil {
		return ErrMalformedPolicy
	}
	if len(doc.Statement) == 0 {
		return ErrMalformedPolicy
	}
	for _, stmt := range doc.Statement {
		if stmt.Effect != "Allow" {
			continue
		}
		// Root principal appears as {"AWS": "arn:...:root"} or "*" — both
		// count as granting access to the account root.
		principalStr := ""
		switch p := stmt.Principal.(type) {
		case string:
			principalStr = p
		case map[string]interface{}:
			if v, ok := p["AWS"].(string); ok {
				principalStr = v
			}
		}
		if principalStr != "*" && !strings.HasSuffix(principalStr, ":root") {
			continue
		}
		// Action must cover kms:PutKeyPolicy (or kms:* / "*").
		switch a := stmt.Action.(type) {
		case string:
			if a == "*" || a == "kms:*" || a == "kms:PutKeyPolicy" {
				return nil
			}
		case []interface{}:
			for _, item := range a {
				if s, ok := item.(string); ok {
					if s == "*" || s == "kms:*" || s == "kms:PutKeyPolicy" {
						return nil
					}
				}
			}
		}
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
