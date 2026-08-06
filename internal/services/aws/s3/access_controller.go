package s3

import (
	"context"
	"time"

	arnutil "vorpalstacks/internal/utils/aws/arn"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// AccessController handles S3 access control evaluation.
type AccessController struct {
	policyEvaluator *policy.PolicyEvaluator
	accountID       string
}

// NewAccessController creates a new access controller.
func NewAccessController(accountID string) *AccessController {
	return &AccessController{
		policyEvaluator: policy.NewPolicyEvaluator(),
		accountID:       accountID,
	}
}

// AccessCheck contains parameters for access control checks.
type AccessCheck struct {
	Principal       string
	PrincipalID     string
	PrincipalType   request.PrincipalType
	Action          string
	Resource        string
	Bucket          string
	Key             string
	SourceIP        string
	Referer         string
	SecureTransport bool
}

// CheckAccess evaluates whether an operation should be allowed.
// Evaluation order follows AWS S3 semantics:
// 1. Owner → Allow (root bypass)
// 2. Bucket Policy explicit Deny → Deny (overrides ACL)
// 3. Bucket Policy explicit Allow → Allow
// 4. ACL Allow → Allow
// 5. Default Deny
func (ac *AccessController) CheckAccess(
	ctx context.Context,
	stores *s3Stores,
	check *AccessCheck,
) error {
	// Service-level operations (ListAllMyBuckets) and CreateBucket do not
	// target an existing bucket, so skip the bucket lookup and policy/ACL
	// evaluation.  IAM policy is evaluated upstream by dispatcher.Authorize.
	if check.Action == "s3:ListAllMyBuckets" || check.Action == "s3:CreateBucket" {
		return nil
	}

	bucket, err := stores.buckets.Get(check.Bucket)
	if err != nil {
		return ErrNoSuchBucket
	}

	if ac.isOwner(check, bucket) {
		return nil
	}

	// Evaluate bucket policy first so explicit Deny overrides ACL Allow.
	if bucket.Policy != "" {
		decision := ac.evaluateBucketPolicyDecision(check, bucket)
		if decision.Effect == policy.DecisionEffectDeny {
			return ErrAccessDenied
		}
		if decision.Effect == policy.DecisionEffectAllow {
			return nil
		}
		// DefaultDeny: fall through to ACL
	}

	if err := ac.evaluateACL(check, bucket, stores); err == nil {
		return nil
	}

	return ErrAccessDenied
}

// CheckObjectAccess evaluates whether an operation on an object should be allowed.
// Evaluation order follows AWS S3 semantics:
// 1. Owner → Allow (root bypass)
// 2. Bucket Policy explicit Deny → Deny (overrides ACL)
// 3. Bucket Policy explicit Allow → Allow
// 4. Object ACL Allow → Allow
// 5. Default Deny
func (ac *AccessController) CheckObjectAccess(
	ctx context.Context,
	stores *s3Stores,
	check *AccessCheck,
) error {
	bucket, err := stores.buckets.Get(check.Bucket)
	if err != nil {
		return ErrNoSuchBucket
	}

	if ac.isOwner(check, bucket) {
		return nil
	}

	// Evaluate bucket policy first so explicit Deny overrides ACL Allow.
	if bucket.Policy != "" {
		decision := ac.evaluateBucketPolicyDecision(check, bucket)
		if decision.Effect == policy.DecisionEffectDeny {
			return ErrAccessDenied
		}
		if decision.Effect == policy.DecisionEffectAllow {
			return nil
		}
		// DefaultDeny: fall through to ACL
	}

	if err := ac.evaluateObjectACL(ctx, check, bucket, stores); err == nil {
		return nil
	}

	return ErrAccessDenied
}

func (ac *AccessController) isOwner(check *AccessCheck, bucket *s3store.Bucket) bool {
	if bucket.ACL != nil && bucket.ACL.Owner != nil {
		if check.PrincipalID == bucket.ACL.Owner.ID {
			return true
		}
	}
	if check.PrincipalID == ac.accountID {
		return true
	}
	return false
}

func (ac *AccessController) evaluateACL(check *AccessCheck, bucket *s3store.Bucket, stores *s3Stores) error {
	acl := bucket.ACL
	if acl == nil {
		if bucket.PublicAccessBlock != nil && bucket.PublicAccessBlock.IgnorePublicAcls {
			return ErrAccessDenied
		}
		return ErrAccessDenied
	}

	if bucket.PublicAccessBlock != nil && bucket.PublicAccessBlock.IgnorePublicAcls {
		if ac.aclContainsPublicAccess(acl) {
			return ErrAccessDenied
		}
	}

	for _, grant := range acl.Grants {
		if ac.grantMatchesPrincipal(grant, check) {
			if ac.permissionMatchesAction(grant.Permission, check.Action, false) {
				return nil
			}
		}
	}

	return ErrAccessDenied
}

func (ac *AccessController) evaluateObjectACL(
	ctx context.Context,
	check *AccessCheck,
	bucket *s3store.Bucket,
	stores *s3Stores,
) error {
	if check.Key == "" {
		return ac.evaluateACL(check, bucket, stores)
	}

	obj, err := stores.objects.Head(ctx, check.Bucket, check.Key)
	if err != nil {
		return ac.evaluateACL(check, bucket, stores)
	}

	if obj.ACL != nil {
		if bucket.PublicAccessBlock != nil && bucket.PublicAccessBlock.IgnorePublicAcls {
			if ac.aclContainsPublicAccess(obj.ACL) {
				return ErrAccessDenied
			}
		}

		for _, grant := range obj.ACL.Grants {
			if ac.grantMatchesPrincipal(grant, check) {
				if ac.permissionMatchesAction(grant.Permission, check.Action, true) {
					return nil
				}
			}
		}
	}

	return ac.evaluateACL(check, bucket, stores)
}

// evaluateBucketPolicyDecision evaluates the bucket policy and returns the
// full decision (Allow, explicit Deny, or DefaultDeny).
func (ac *AccessController) evaluateBucketPolicyDecision(check *AccessCheck, bucket *s3store.Bucket) *policy.Decision {
	if bucket.PublicAccessBlock != nil && bucket.PublicAccessBlock.RestrictPublicBuckets {
		if check.PrincipalType == request.PrincipalTypeAnonymous {
			return &policy.Decision{Effect: policy.DecisionEffectDeny, Reason: "anonymous access restricted by PublicAccessBlock"}
		}
	}

	policyDoc, err := policy.ParseDocument(bucket.Policy)
	if err != nil {
		return &policy.Decision{Effect: policy.DecisionEffectDeny, Reason: "invalid bucket policy"}
	}

	evalCtx := &policy.EvaluationContext{
		Principal:        check.Principal,
		PrincipalAccount: ac.extractAccountFromPrincipal(check.Principal),
		Action:           check.Action,
		Resource:         check.Resource,
		SourceIP:         check.SourceIP,
		RequestTime:      time.Now(),
		Referer:          check.Referer,
		SecureTransport:  check.SecureTransport,
	}

	return ac.policyEvaluator.Evaluate(evalCtx, []*policy.Document{policyDoc})
}

func (ac *AccessController) grantMatchesPrincipal(grant *s3store.Grant, check *AccessCheck) bool {
	if grant.Grantee == nil {
		return false
	}

	switch grant.Grantee.Type {
	case s3store.GranteeTypeCanonicalUser:
		return grant.Grantee.ID == check.PrincipalID || grant.Grantee.ID == ac.accountID
	case s3store.GranteeTypeGroup:
		switch grant.Grantee.URI {
		case s3store.AllUsersGroup:
			return true
		case s3store.AuthenticatedUsersGroup:
			return check.PrincipalType != request.PrincipalTypeAnonymous
		case s3store.LogDeliveryGroup:
			return check.PrincipalType != request.PrincipalTypeAnonymous
		}
	}
	return false
}

func (ac *AccessController) permissionMatchesAction(perm s3store.Permission, action string, isObject bool) bool {
	switch perm {
	case s3store.PermissionFullControl:
		return true
	case s3store.PermissionRead:
		if isObject {
			return action == "s3:GetObject"
		}
		return action == "s3:GetObject" || action == "s3:ListBucket" ||
			action == "s3:ListBucketVersions" || action == "s3:ListMultipartUploadParts" ||
			action == "s3:ListBucketMultipartUploads"
	case s3store.PermissionWrite:
		if isObject {
			return action == "s3:PutObject" || action == "s3:DeleteObject" ||
				action == "s3:PutObjectTagging" || action == "s3:DeleteObjectTagging"
		}
		return action == "s3:PutObject" || action == "s3:DeleteObject" ||
			action == "s3:AbortMultipartUpload" || action == "s3:CreateMultipartUpload" ||
			action == "s3:UploadPart" || action == "s3:CompleteMultipartUpload"
	case s3store.PermissionReadACP:
		if isObject {
			return action == "s3:GetObjectAcl"
		}
		return action == "s3:GetBucketAcl" || action == "s3:GetBucketPolicyStatus"
	case s3store.PermissionWriteACP:
		if isObject {
			return action == "s3:PutObjectAcl"
		}
		return action == "s3:PutBucketAcl" || action == "s3:PutBucketPolicy"
	}
	return false
}

func (ac *AccessController) aclContainsPublicAccess(acl *s3store.AccessControlPolicy) bool {
	if acl == nil {
		return false
	}
	for _, grant := range acl.Grants {
		if grant.Grantee == nil {
			continue
		}
		if grant.Grantee.URI == s3store.AllUsersGroup ||
			grant.Grantee.URI == s3store.AuthenticatedUsersGroup {
			return true
		}
	}
	return false
}

func (ac *AccessController) extractAccountFromPrincipal(principal string) string {
	if principal == "" || principal == "*" {
		return ""
	}

	if _, service, _, accountID, _ := arnutil.SplitARN(principal); service == "iam" {
		return accountID
	}
	return ""
}

func buildResource(accountID, region, bucket, key string) string {
	if key != "" {
		return arnutil.NewARNBuilder("", "").S3().Object(bucket, key)
	}
	return arnutil.NewARNBuilder("", "").S3().Bucket(bucket)
}
