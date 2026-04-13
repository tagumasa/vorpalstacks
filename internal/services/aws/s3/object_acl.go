package s3

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// PutObjectAclInput contains the parameters for setting an object's Access Control List.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally specifies a specific version of the object.
// ACL is a canned ACL (e.g., "private", "public-read").
// AccessControlPolicy provides explicit grant configuration.
// GrantFullControl, GrantRead, GrantReadACP, GrantWrite, GrantWriteACP specify grants via headers.
type PutObjectAclInput struct {
	Bucket              string
	Key                 string
	VersionId           string
	ACL                 string
	AccessControlPolicy *s3store.AccessControlPolicy
	GrantFullControl    string
	GrantRead           string
	GrantReadACP        string
	GrantWrite          string
	GrantWriteACP       string
}

// GetObjectAclOutput contains the result of retrieving an object's Access Control List.
// Owner is the owner of the object.
// Grants contains the list of access grants.
type GetObjectAclOutput struct {
	Owner  *s3store.ACLOwner `xml:"Owner"`
	Grants []*s3store.Grant  `xml:"AccessControlList>Grant"`
}

// ToXML serialises the ACL to XML format for S3 API response.
func (o *GetObjectAclOutput) ToXML() string {
	return BuildACLXML(o.Owner, o.Grants)
}

// GetObjectAcl retrieves the Access Control List for an object.
// Returns the owner and list of grants for the specified object version.
func (o *ObjectOperations) GetObjectAcl(ctx context.Context, reqCtx *request.RequestContext, bucket, key, versionId string) (*GetObjectAclOutput, error) {
	if err := o.validateBucketExists(reqCtx, bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(key); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	acp, err := store.objects.GetACLWithVersion(bucket, key, versionId)
	if err != nil {
		return nil, err
	}

	owner := &s3store.ACLOwner{ID: o.svc.accountID, DisplayName: o.svc.accountID}

	if acp == nil {
		return &GetObjectAclOutput{
			Owner: owner,
			Grants: []*s3store.Grant{
				{
					Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: o.svc.accountID, DisplayName: o.svc.accountID},
					Permission: s3store.PermissionFullControl,
				},
			},
		}, nil
	}

	return &GetObjectAclOutput{
		Owner:  acp.Owner,
		Grants: acp.Grants,
	}, nil
}

// PutObjectAcl sets the Access Control List for an object.
// Accepts either a canned ACL string, an AccessControlPolicy, or individual grant headers.
func (o *ObjectOperations) PutObjectAcl(ctx context.Context, reqCtx *request.RequestContext, input *PutObjectAclInput) error {
	store, err := o.svc.store(reqCtx)
	if err != nil {
		return err
	}

	owner := &s3store.ACLOwner{ID: o.svc.accountID, DisplayName: o.svc.accountID}

	var acp *s3store.AccessControlPolicy

	if input.ACL != "" {
		acp, err = CannedACLToPolicy(input.ACL, owner)
		if err != nil {
			return err
		}
	} else if input.AccessControlPolicy != nil {
		acp = input.AccessControlPolicy
	} else {
		grants := ParseGrantHeaders(input.GrantFullControl, input.GrantRead, input.GrantReadACP, input.GrantWrite, input.GrantWriteACP)
		if len(grants) > 0 {
			acp = &s3store.AccessControlPolicy{Owner: owner, Grants: grants}
		} else {
			return fmt.Errorf("missing required ACL specification")
		}
	}

	return store.objects.SetACLWithVersion(input.Bucket, input.Key, input.VersionId, acp)
}
