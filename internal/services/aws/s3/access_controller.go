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
//
// The S3 HTTP plane is mounted ahead of the gRPC dispatcher, so identity
// policies are not evaluated here: bucket policies and ACLs are the only
// authorisation inputs on this path.
// evaluatePolicyThenACL runs the evaluation ladder the bucket and object
// planes share: owner → Allow, bucket-policy explicit Deny → Deny, bucket
// policy explicit Allow → Allow, then the plane-specific ACL evaluator.
func (ac *AccessController) evaluatePolicyThenACL(
	stores *s3Stores,
	check *AccessCheck,
	aclEvaluator func(bucket *s3store.Bucket) error,
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

	if err := aclEvaluator(bucket); err == nil {
		return nil
	}

	return ErrAccessDenied
}

func (ac *AccessController) CheckAccess(
	ctx context.Context,
	stores *s3Stores,
	check *AccessCheck,
) error {
	// Service-level operations (ListAllMyBuckets) and CreateBucket do not
	// target an existing bucket, so skip the bucket lookup and policy/ACL
	// evaluation.
	if check.Action == "s3:ListAllMyBuckets" || check.Action == "s3:CreateBucket" {
		return nil
	}

	return ac.evaluatePolicyThenACL(stores, check, func(bucket *s3store.Bucket) error {
		return ac.evaluateACL(check, bucket, stores)
	})
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
	return ac.evaluatePolicyThenACL(stores, check, func(bucket *s3store.Bucket) error {
		return ac.evaluateObjectACL(ctx, check, bucket, stores)
	})
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
	// A bucket with no ACL denies every ACL-based check: both a pristine
	// bucket and one whose ACL was never written behave the same whether
	// or not public-access blocks are set.
	if acl == nil {
		return ErrAccessDenied
	}

	if bucket.PublicAccessBlock != nil && bucket.PublicAccessBlock.IgnorePublicAcls {
		if acpContainsPublicAccess(acl) {
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
			if acpContainsPublicAccess(obj.ACL) {
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

// permissionMatchesAction applies the ACL-to-action mapping (AWS ACL
// overview, "Mapping of ACL permissions and access policy permissions").
// The table is encoded action-first so each action names the permissions
// whose row carries it: READ/READ_ACP/WRITE_ACP rows include the Version
// twins of their base actions (the classifier emits the versioned action
// for version-addressed requests, and the ACL grant covers both in AWS).
// WRITE on an object maps to no action ("Not applicable" in the table);
// FULL_CONTROL is the union of the four rows for the resource, never a
// blanket match for arbitrary actions.
func (ac *AccessController) permissionMatchesAction(perm s3store.Permission, action string, isObject bool) bool {
	grants := func(perms ...s3store.Permission) bool {
		for _, p := range perms {
			if perm == p {
				return true
			}
		}
		return false
	}

	if isObject {
		switch action {
		case "s3:GetObject", "s3:GetObjectVersion":
			return grants(s3store.PermissionRead, s3store.PermissionFullControl)
		case "s3:GetObjectAcl", "s3:GetObjectVersionAcl":
			return grants(s3store.PermissionReadACP, s3store.PermissionFullControl)
		case "s3:PutObjectAcl", "s3:PutObjectVersionAcl":
			return grants(s3store.PermissionWriteACP, s3store.PermissionFullControl)
		}
		return false
	}

	switch action {
	case "s3:ListBucket", "s3:ListBucketVersions", "s3:ListBucketMultipartUploads":
		return grants(s3store.PermissionRead, s3store.PermissionFullControl)
	case "s3:PutObject":
		return grants(s3store.PermissionWrite, s3store.PermissionFullControl)
	case "s3:GetBucketAcl":
		return grants(s3store.PermissionReadACP, s3store.PermissionFullControl)
	case "s3:PutBucketAcl":
		return grants(s3store.PermissionWriteACP, s3store.PermissionFullControl)
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
