package s3

import (
	"context"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// PutBucketAclInput is the input for PutBucketAcl.
type PutBucketAclInput struct {
	Bucket              string
	ACL                 string
	AccessControlPolicy *s3store.AccessControlPolicy
	GrantFullControl    string
	GrantRead           string
	GrantReadACP        string
	GrantWrite          string
	GrantWriteACP       string
}

// GetBucketAclOutput is the output of GetBucketAcl containing owner and grants.
type GetBucketAclOutput struct {
	Owner  *s3store.ACLOwner `xml:"Owner"`
	Grants []*s3store.Grant  `xml:"AccessControlList>Grant"`
}

// ToXML converts the ACL output to XML format.
func (o *GetBucketAclOutput) ToXML() string {
	return BuildACLXML(o.Owner, o.Grants)
}

// PutBucketAcl sets the access control list (ACL) for an S3 bucket.
func (o *BucketOperations) PutBucketAcl(ctx *request.RequestContext, input *PutBucketAclInput) error {
	owner := &s3store.ACLOwner{ID: o.svc.accountID, DisplayName: o.svc.accountID}

	var acp *s3store.AccessControlPolicy
	var err error

	if input.ACL != "" {
		acp, err = CannedACLToPolicy(input.ACL, owner)
		if err != nil {
			return err
		}
	} else if input.AccessControlPolicy != nil {
		acp = input.AccessControlPolicy
	} else {
		grants, err := ParseGrantHeaders(input.GrantFullControl, input.GrantRead, input.GrantReadACP, input.GrantWrite, input.GrantWriteACP)
		if err != nil {
			return NewInvalidArgumentError(err.Error())
		}
		if len(grants) > 0 {
			acp = &s3store.AccessControlPolicy{Owner: owner, Grants: grants}
		} else {
			return NewInvalidArgumentError("missing required ACL specification")
		}
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	publicAccessBlock, _ := store.buckets.GetPublicAccessBlock(input.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicAcls {
		if isPublicCannedACL(input.ACL) {
			return NewInvalidArgumentError("bucket has BlockPublicAcls enabled")
		}
		if acpContainsPublicAccess(acp) {
			return NewInvalidArgumentError("bucket has BlockPublicAcls enabled")
		}
	}

	// With Object Ownership set to BucketOwnerEnforced, "requests to set or
	// update ACLs fail" with AccessControlListNotSupported.
	if aclsDisabled, _ := o.svc.bucketACLsDisabled(ctx, store, input.Bucket); aclsDisabled {
		return ErrAccessControlListNotSupported
	}

	return store.buckets.SetACL(input.Bucket, acp)
}

// bucketACLsDisabled reports whether the bucket's Object Ownership setting
// is BucketOwnerEnforced, under which ACLs are disabled and set/update ACL
// requests fail.
func (s *S3Service) bucketACLsDisabled(ctx context.Context, store *s3Stores, bucket string) (bool, error) {
	b, err := store.buckets.Get(bucket)
	if err != nil {
		return false, err
	}
	return b.OwnershipControls != nil &&
		len(b.OwnershipControls.Rules) == 1 &&
		b.OwnershipControls.Rules[0].ObjectOwnership == "BucketOwnerEnforced", nil
}

// resolveUploadACL validates the ACL request headers of an upload (PutObject,
// CreateMultipartUpload, CopyObject) and builds the object ACL they express.
// With ACLs disabled (BucketOwnerEnforced), object ownership "accepts only
// PUT requests that do not specify an ACL or PUT requests with bucket owner
// full control ACLs"; anything else fails with 400
// AccessControlListNotSupported. A nil result means no ACL was requested
// (or the request is a no-op under disabled ACLs).
func (s *S3Service) resolveUploadACL(ctx context.Context, store *s3Stores, bucket string, headers aclHeaders) (*s3store.AccessControlPolicy, error) {
	specified := headers.ACL != "" || headers.GrantFullControl != "" || headers.GrantRead != "" ||
		headers.GrantReadACP != "" || headers.GrantWrite != "" || headers.GrantWriteACP != ""
	if !specified {
		return nil, nil
	}

	disabled, err := s.bucketACLsDisabled(ctx, store, bucket)
	if err != nil {
		return nil, err
	}
	if disabled {
		if headers.ACL != "" && headers.ACL != "private" && headers.ACL != "bucket-owner-full-control" {
			return nil, ErrAccessControlListNotSupported
		}
		if headers.GrantFullControl != "" || headers.GrantRead != "" || headers.GrantReadACP != "" ||
			headers.GrantWrite != "" || headers.GrantWriteACP != "" {
			return nil, ErrAccessControlListNotSupported
		}
		return nil, nil
	}

	owner := &s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID}
	if headers.ACL != "" {
		return CannedACLToPolicy(headers.ACL, owner)
	}
	grants, err := ParseGrantHeaders(headers.GrantFullControl, headers.GrantRead, headers.GrantReadACP, headers.GrantWrite, headers.GrantWriteACP)
	if err != nil {
		return nil, NewInvalidArgumentError(err.Error())
	}
	return &s3store.AccessControlPolicy{Owner: owner, Grants: grants}, nil
}

func acpContainsPublicAccess(acp *s3store.AccessControlPolicy) bool {
	if acp == nil {
		return false
	}
	for _, grant := range acp.Grants {
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

// GetBucketAcl retrieves the access control list (ACL) for an S3 bucket.
func (o *BucketOperations) GetBucketAcl(ctx *request.RequestContext, bucket string) (*GetBucketAclOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	b, err := store.buckets.Get(bucket)
	if err != nil {
		return nil, err
	}

	owner := &s3store.ACLOwner{ID: o.svc.accountID, DisplayName: o.svc.accountID}

	if b.ACL == nil {
		return &GetBucketAclOutput{
			Owner: owner,
			Grants: []*s3store.Grant{
				{
					Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: o.svc.accountID, DisplayName: o.svc.accountID},
					Permission: s3store.PermissionFullControl,
				},
			},
		}, nil
	}

	return &GetBucketAclOutput{
		Owner:  b.ACL.Owner,
		Grants: b.ACL.Grants,
	}, nil
}
