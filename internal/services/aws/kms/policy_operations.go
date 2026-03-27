package kms

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
)

// GetKeyPolicy retrieves the key policy for a specified key.
// You can specify the policy name; if not provided, the default policy is returned.
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

	policyDocument := request.GetStringParam(req.Parameters, "Policy")
	if policyDocument == "" {
		return nil, ErrMalformedPolicy
	}

	var js interface{}
	if err := json.Unmarshal([]byte(policyDocument), &js); err != nil {
		return nil, ErrMalformedPolicy
	}

	if err := stores.keyPolicies.Put(key.KeyID, policyName, policyDocument); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
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

	return map[string]interface{}{
		"PolicyNames": policies,
		"Truncated":   false,
	}, nil
}
