package s3

import (
	"fmt"

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
		grants := ParseGrantHeaders(input.GrantFullControl, input.GrantRead, input.GrantReadACP, input.GrantWrite, input.GrantWriteACP)
		if len(grants) > 0 {
			acp = &s3store.AccessControlPolicy{Owner: owner, Grants: grants}
		} else {
			return fmt.Errorf("missing required ACL specification")
		}
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	publicAccessBlock, _ := store.buckets.GetPublicAccessBlock(input.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicAcls {
		if input.ACL == "public-read" || input.ACL == "public-read-write" {
			return fmt.Errorf("bucket has BlockPublicAcls enabled")
		}
		if acpContainsPublicAccess(acp) {
			return fmt.Errorf("bucket has BlockPublicAcls enabled")
		}
	}

	return store.buckets.SetACL(input.Bucket, acp)
}

func acpContainsPublicAccess(acp *s3store.AccessControlPolicy) bool {
	if acp == nil {
		return false
	}
	for _, grant := range acp.Grants {
		if grant.Grantee != nil && grant.Grantee.URI == s3store.AllUsersGroup {
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

func cannedACLToPolicy(cannedACL string, owner *s3store.ACLOwner) (*s3store.AccessControlPolicy, error) {
	if owner == nil {
		owner = &s3store.ACLOwner{ID: "owner", DisplayName: "owner"}
	}

	grants := []*s3store.Grant{
		{
			Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: owner.ID, DisplayName: owner.DisplayName},
			Permission: s3store.PermissionFullControl,
		},
	}

	switch cannedACL {
	case "private":
	case "public-read":
		grants = append(grants, &s3store.Grant{
			Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AllUsersGroup},
			Permission: s3store.PermissionRead,
		})
	case "public-read-write":
		grants = append(grants,
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AllUsersGroup}, Permission: s3store.PermissionRead},
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AllUsersGroup}, Permission: s3store.PermissionWrite},
		)
	case "authenticated-read":
		grants = append(grants, &s3store.Grant{
			Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.AuthenticatedUsersGroup},
			Permission: s3store.PermissionRead,
		})
	case "log-delivery-write":
		grants = append(grants,
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.LogDeliveryGroup}, Permission: s3store.PermissionReadACP},
			&s3store.Grant{Grantee: &s3store.Grantee{Type: s3store.GranteeTypeGroup, URI: s3store.LogDeliveryGroup}, Permission: s3store.PermissionWrite},
		)
	case "bucket-owner-read", "bucket-owner-full-control", "aws-exec-read":
	default:
		return nil, fmt.Errorf("invalid canned ACL: %s", cannedACL)
	}

	return &s3store.AccessControlPolicy{Owner: owner, Grants: grants}, nil
}
