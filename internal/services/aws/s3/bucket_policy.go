package s3

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
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

	publicAccessBlock, _ := store.buckets.GetPublicAccessBlock(input.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicPolicy {
		if policyContainsPublicAccess(input.Policy) {
			return fmt.Errorf("bucket has BlockPublicPolicy enabled")
		}
	}

	return store.buckets.SetPolicy(input.Bucket, input.Policy)
}

func policyContainsPublicAccess(policy string) bool {
	return strings.Contains(policy, `"Principal": "*"`) ||
		strings.Contains(policy, `"Principal":"*"`) ||
		strings.Contains(policy, `"Principal": {"AWS": "*"`) ||
		strings.Contains(policy, `"Principal":{"AWS":"*"`) ||
		strings.Contains(strings.ToLower(policy), `"principal": "*"`) ||
		strings.Contains(strings.ToLower(policy), `"principal":{"aws":"*"`)
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
