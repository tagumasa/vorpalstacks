package s3

import (
	"encoding/json"
	"strings"

	"vorpalstacks/internal/common/request"
)

// PutBucketPolicyInput contains the request parameters for the PutBucketPolicy operation.
type PutBucketPolicyInput struct {
	Bucket string
	Policy string
}

// PutBucketPolicy applies a policy to a bucket.
// The bucket owner must have PutBucketPolicy permission.
// Returns an error if BlockPublicPolicy is enabled and the policy grants public access.
func (o *BucketOperations) PutBucketPolicy(ctx *request.RequestContext, input *PutBucketPolicyInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	if err := validatePolicyDocument(input.Policy); err != nil {
		return err
	}

	publicAccessBlock, _ := store.buckets.GetPublicAccessBlock(input.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicPolicy {
		if policyContainsPublicAccess(input.Policy) {
			return NewInvalidArgumentError("bucket has BlockPublicPolicy enabled")
		}
	}

	return store.buckets.SetPolicy(input.Bucket, input.Policy)
}

// policyContainsPublicAccess checks whether a bucket policy grants public
// access. A statement is considered public if:
//   - Effect is "Allow" AND
//   - Principal is "*" (or contains wildcard account) AND
//   - Either no Condition is set, or the Condition does not restrict access
//     (e.g. IpAddress with 0.0.0.0/0 or ::/0)
func policyContainsPublicAccess(policy string) bool {
	var p struct {
		Statement []struct {
			Effect    string                 `json:"Effect"`
			Principal interface{}            `json:"Principal"`
			Condition map[string]interface{} `json:"Condition"`
		} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policy), &p); err != nil {
		return false
	}
	for _, stmt := range p.Statement {
		if !strings.EqualFold(stmt.Effect, "Allow") {
			continue
		}
		if !isPublicPrincipal(stmt.Principal) {
			continue
		}
		// Principal is "*" and Effect is Allow — check if Condition restricts.
		if stmt.Condition == nil || conditionIsEffectivelyPublic(stmt.Condition) {
			return true
		}
	}
	return false
}

// conditionIsEffectivelyPublic returns true when a Condition block either
// does not restrict access or uses a permissive IP range like 0.0.0.0/0.
func conditionIsEffectivelyPublic(condition map[string]interface{}) bool {
	for condOp, condVal := range condition {
		vals, ok := condVal.(map[string]interface{})
		if !ok {
			continue
		}
		opLower := strings.ToLower(condOp)
		switch opLower {
		case "ipaddress":
			// IpAddress with 0.0.0.0/0 or ::/0 is effectively public.
			for _, v := range vals {
				if ipMatchesAll(v) {
					return true
				}
			}
		case "stringlike":
			// StringLike with "*" wildcard on any key is effectively public.
			for _, v := range vals {
				if v == "*" {
					return true
				}
			}
		}
	}
	return false
}

// ipMatchesAll checks whether a CIDR value matches all IP addresses.
func ipMatchesAll(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return v == "0.0.0.0/0" || v == "::/0"
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && (s == "0.0.0.0/0" || s == "::/0") {
				return true
			}
		}
	}
	return false
}

func isPublicPrincipal(principal interface{}) bool {
	switch v := principal.(type) {
	case string:
		return v == "*"
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && s == "*" {
				return true
			}
		}
	case map[string]interface{}:
		for _, val := range v {
			if s, ok := val.(string); ok && s == "*" {
				return true
			}
			if arr, ok := val.([]interface{}); ok {
				for _, item := range arr {
					if s, ok := item.(string); ok && s == "*" {
						return true
					}
				}
			}
		}
	}
	return false
}

// GetBucketPolicyInput contains the request parameters for the GetBucketPolicy operation.
type GetBucketPolicyInput struct {
	Bucket string
}

// GetBucketPolicyOutput contains the result of the GetBucketPolicy operation.
type GetBucketPolicyOutput struct {
	Policy string `xml:"Policy"`
}

// GetBucketPolicy retrieves the policy attached to a bucket.
func (o *BucketOperations) GetBucketPolicy(ctx *request.RequestContext, input *GetBucketPolicyInput) (*GetBucketPolicyOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.Policy == "" {
		return nil, ErrNoSuchBucketPolicy
	}

	return &GetBucketPolicyOutput{
		Policy: bucket.Policy,
	}, nil
}

// DeleteBucketPolicyInput contains the request parameters for the DeleteBucketPolicy operation.
type DeleteBucketPolicyInput struct {
	Bucket string
}

// DeleteBucketPolicy removes the policy from a bucket.
func (o *BucketOperations) DeleteBucketPolicy(ctx *request.RequestContext, input *DeleteBucketPolicyInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetPolicy(input.Bucket, "")
}
