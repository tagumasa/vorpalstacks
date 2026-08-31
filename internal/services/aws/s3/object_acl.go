package s3

import (
	"context"

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
func (o *ObjectOperations) GetObjectAcl(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, bucket, key, versionId string) (*GetObjectAclOutput, error) {
	return o.svc.getObjectAclCore(stores, bucket, key, versionId)
}

// PutObjectAcl sets the Access Control List for an object.
// Accepts either a canned ACL string, an AccessControlPolicy, or individual grant headers.
func (o *ObjectOperations) PutObjectAcl(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectAclInput) error {
	return o.svc.putObjectAclCore(ctx, stores, input)
}
