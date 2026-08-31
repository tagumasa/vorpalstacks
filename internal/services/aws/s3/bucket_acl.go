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
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketAclCore(ctx, store, input)
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
	return o.svc.getBucketAclCore(store.buckets, bucket)
}
