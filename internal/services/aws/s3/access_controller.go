// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/iam/policy"
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
func (ac *AccessController) CheckAccess(
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

	if err := ac.evaluateACL(check, bucket, stores); err == nil {
		return nil
	}

	if bucket.Policy != "" {
		if err := ac.evaluateBucketPolicy(check, bucket); err == nil {
			return nil
		}
	}

	return ErrAccessDenied
}

// CheckObjectAccess evaluates whether an operation on an object should be allowed.
// It checks bucket ownership, object ACL, and bucket policy.
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

	if err := ac.evaluateObjectACL(ctx, check, bucket, stores); err == nil {
		return nil
	}

	if bucket.Policy != "" {
		if err := ac.evaluateBucketPolicy(check, bucket); err == nil {
			return nil
		}
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

func (ac *AccessController) evaluateBucketPolicy(check *AccessCheck, bucket *s3store.Bucket) error {
	if bucket.PublicAccessBlock != nil && bucket.PublicAccessBlock.RestrictPublicBuckets {
		if check.PrincipalType == request.PrincipalTypeAnonymous {
			return ErrAccessDenied
		}
	}

	policyDoc, err := policy.ParseDocument(bucket.Policy)
	if err != nil {
		return ErrAccessDenied
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

	decision := ac.policyEvaluator.Evaluate(evalCtx, []*policy.Document{policyDoc})
	if decision.Effect == policy.DecisionEffectAllow {
		return nil
	}

	return ErrAccessDenied
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
			action == "s3:ListBucketVersions" || action == "s3:ListMultipartUploadParts"
	case s3store.PermissionWrite:
		if isObject {
			return action == "s3:PutObject" || action == "s3:DeleteObject"
		}
		return action == "s3:PutObject" || action == "s3:DeleteObject" ||
			action == "s3:AbortMultipartUpload" || action == "s3:CreateMultipartUpload" ||
			action == "s3:UploadPart" || action == "s3:CompleteMultipartUpload"
	case s3store.PermissionReadACP:
		if isObject {
			return action == "s3:GetObjectAcl"
		}
		return action == "s3:GetBucketAcl"
	case s3store.PermissionWriteACP:
		if isObject {
			return action == "s3:PutObjectAcl"
		}
		return action == "s3:PutBucketAcl"
	}
	return false
}

func (ac *AccessController) aclContainsPublicAccess(acl *s3store.AccessControlPolicy) bool {
	if acl == nil {
		return false
	}
	for _, grant := range acl.Grants {
		if grant.Grantee != nil && grant.Grantee.URI == s3store.AllUsersGroup {
			return true
		}
	}
	return false
}

func (ac *AccessController) extractAccountFromPrincipal(principal string) string {
	if principal == "" || principal == "*" {
		return ""
	}

	if strings.HasPrefix(principal, "arn:aws:iam::") {
		parts := strings.Split(principal, ":")
		if len(parts) >= 5 {
			return parts[4]
		}
	}
	return ""
}

func actionFromMethod(method string, isBucket bool, queryKey string) string {
	if isBucket {
		switch method {
		case "GET":
			switch {
			case queryKey == "acl":
				return "s3:GetBucketAcl"
			case queryKey == "policy":
				return "s3:GetBucketPolicy"
			case queryKey == "versioning":
				return "s3:GetBucketVersioning"
			case queryKey == "encryption":
				return "s3:GetEncryptionConfiguration"
			case queryKey == "cors":
				return "s3:GetBucketCORS"
			case queryKey == "tagging":
				return "s3:GetBucketTagging"
			case queryKey == "location":
				return "s3:GetBucketLocation"
			case queryKey == "versions":
				return "s3:ListBucketVersions"
			case queryKey == "uploads":
				return "s3:ListBucketMultipartUploads"
			default:
				return "s3:ListBucket"
			}
		case "PUT":
			switch {
			case queryKey == "acl":
				return "s3:PutBucketAcl"
			case queryKey == "policy":
				return "s3:PutBucketPolicy"
			case queryKey == "versioning":
				return "s3:PutBucketVersioning"
			case queryKey == "encryption":
				return "s3:PutEncryptionConfiguration"
			case queryKey == "cors":
				return "s3:PutBucketCORS"
			case queryKey == "tagging":
				return "s3:PutBucketTagging"
			default:
				return "s3:CreateBucket"
			}
		case "DELETE":
			switch {
			case queryKey == "policy":
				return "s3:DeleteBucketPolicy"
			case queryKey == "encryption":
				return "s3:DeleteEncryptionConfiguration"
			case queryKey == "cors":
				return "s3:DeleteBucketCORS"
			case queryKey == "tagging":
				return "s3:DeleteBucketTagging"
			default:
				return "s3:DeleteBucket"
			}
		case "HEAD":
			return "s3:ListBucket"
		}
	} else {
		switch method {
		case "GET":
			switch {
			case queryKey == "acl":
				return "s3:GetObjectAcl"
			case queryKey == "uploadId":
				return "s3:ListMultipartUploadParts"
			case queryKey == "uploads":
				return "s3:ListBucketMultipartUploads"
			default:
				return "s3:GetObject"
			}
		case "PUT":
			switch {
			case queryKey == "acl":
				return "s3:PutObjectAcl"
			case queryKey == "uploadId":
				return "s3:UploadPart"
			case queryKey == "uploads":
				return "s3:CreateMultipartUpload"
			default:
				return "s3:PutObject"
			}
		case "DELETE":
			return "s3:DeleteObject"
		case "HEAD":
			return "s3:GetObject"
		case "POST":
			switch {
			case queryKey == "uploadId":
				return "s3:CompleteMultipartUpload"
			case queryKey == "delete":
				return "s3:DeleteObject"
			case queryKey == "uploads":
				return "s3:CreateMultipartUpload"
			default:
				return "s3:PutObject"
			}
		}
	}
	return ""
}

func buildResource(accountID, region, bucket, key string) string {
	if key != "" {
		return "arn:aws:s3:::" + bucket + "/" + key
	}
	return "arn:aws:s3:::" + bucket
}
