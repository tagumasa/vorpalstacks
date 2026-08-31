package s3

import (
	"strings"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// BucketOperations provides bucket-level S3 operations.
type BucketOperations struct {
	svc *S3Service
}

// NewBucketOperations creates a new BucketOperations instance.
func NewBucketOperations(svc *S3Service) *BucketOperations {
	return &BucketOperations{svc: svc}
}

// CreateBucketInput contains the input parameters for the CreateBucket operation.
type CreateBucketInput struct {
	Bucket                     string
	ACL                        string
	GrantFullControl           string
	GrantRead                  string
	GrantReadACP               string
	GrantWrite                 string
	GrantWriteACP              string
	ObjectOwnership            string
	LocationConstraint         string
	ObjectLockEnabledForBucket bool
}

// CreateBucketOutput contains the output result of the CreateBucket operation.
type CreateBucketOutput struct {
	Location string
}

// CreateBucket creates a new bucket.
func (o *BucketOperations) CreateBucket(ctx *request.RequestContext, input *CreateBucketInput) (*CreateBucketOutput, error) {
	region := ctx.GetRegion()
	if input.LocationConstraint != "" && input.LocationConstraint != defaults.DefaultRegion {
		region = input.LocationConstraint
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	result, err := o.svc.createBucketCore(store.buckets, AdminCreateBucketInput{
		Bucket:                     input.Bucket,
		Region:                     region,
		ACL:                        input.ACL,
		GrantFullControl:           input.GrantFullControl,
		GrantRead:                  input.GrantRead,
		GrantReadACP:               input.GrantReadACP,
		GrantWrite:                 input.GrantWrite,
		GrantWriteACP:              input.GrantWriteACP,
		ObjectOwnership:            input.ObjectOwnership,
		ObjectLockEnabledForBucket: input.ObjectLockEnabledForBucket,
	})
	if err != nil {
		return nil, err
	}

	return &CreateBucketOutput{Location: result.Location}, nil
}

// DeleteBucketInput contains the input parameters for the DeleteBucket operation.
type DeleteBucketInput struct {
	Bucket string
}

// DeleteBucket deletes a bucket.
func (o *BucketOperations) DeleteBucket(ctx *request.RequestContext, input *DeleteBucketInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	_, err = o.svc.deleteBucketCore(store.buckets, store.objects, AdminDeleteBucketInput{
		Bucket: input.Bucket,
	})
	return err
}

// GetBucketInput contains the input parameters for the GetBucket operation.
type GetBucketInput struct {
	Bucket string
}

// GetBucket retrieves information about a bucket.
func (o *BucketOperations) GetBucket(ctx *request.RequestContext, input *GetBucketInput) (*s3store.Bucket, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketCore(store.buckets, input)
}

// HeadBucketInput contains the input parameters for the HeadBucket operation.
type HeadBucketInput struct {
	Bucket string
}

// HeadBucketOutput contains the output result of the HeadBucket operation.
type HeadBucketOutput struct {
	BucketRegion string
}

// HeadBucket retrieves the metadata headers for an S3 bucket without listing its contents.
func (o *BucketOperations) HeadBucket(ctx *request.RequestContext, input *HeadBucketInput) (*HeadBucketOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.headBucketCore(store.buckets, input)
}

// ListBucketsInput contains the input parameters for the ListBuckets operation.
type ListBucketsInput struct {
	MaxBuckets        int
	ContinuationToken string
	Prefix            string
	BucketRegion      string
}

// ListBucketsOutput contains the output result of the ListBuckets operation.
type ListBucketsOutput struct {
	Owner             *Owner
	Buckets           []*BucketInfo
	Prefix            string
	ContinuationToken string
	IsTruncated       bool
}

// ToXML converts the ListBucketsOutput to XML format.
func (o *ListBucketsOutput) ToXML() string {
	var result strings.Builder
	result.WriteString(`<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	if o.Owner != nil {
		result.WriteString(`<Owner><ID>`)
		result.WriteString(xmlEscape(o.Owner.ID))
		result.WriteString(`</ID><DisplayName>`)
		result.WriteString(xmlEscape(o.Owner.DisplayName))
		result.WriteString(`</DisplayName></Owner>`)
	}
	if o.Prefix != "" {
		result.WriteString(`<Prefix>`)
		result.WriteString(xmlEscape(o.Prefix))
		result.WriteString(`</Prefix>`)
	}
	if o.ContinuationToken != "" {
		result.WriteString(`<ContinuationToken>`)
		result.WriteString(xmlEscape(o.ContinuationToken))
		result.WriteString(`</ContinuationToken>`)
	}
	if o.IsTruncated {
		result.WriteString(`<IsTruncated>true</IsTruncated>`)
	} else {
		result.WriteString(`<IsTruncated>false</IsTruncated>`)
	}
	result.WriteString(`<Buckets>`)
	if len(o.Buckets) > 0 {
		for _, b := range o.Buckets {
			result.WriteString(`<Bucket><Name>`)
			result.WriteString(xmlEscape(b.Name))
			result.WriteString(`</Name><CreationDate>`)
			result.WriteString(b.CreationDate.Format(time.RFC3339))
			result.WriteString(`</CreationDate></Bucket>`)
		}
	}
	result.WriteString(`</Buckets>`)
	result.WriteString(`</ListAllMyBucketsResult>`)
	return result.String()
}

func xmlEscape(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			result.WriteString("&amp;")
		case '<':
			result.WriteString("&lt;")
		case '>':
			result.WriteString("&gt;")
		case '"':
			result.WriteString("&quot;")
		case '\'':
			result.WriteString("&apos;")
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Owner represents the owner of a bucket.
type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

// BucketInfo contains metadata about a bucket.
type BucketInfo struct {
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
}

// ListBuckets lists all buckets in the account.
// Supports pagination via MaxBuckets/ContinuationToken (AWS S3 ListBuckets v2).
func (o *BucketOperations) ListBuckets(ctx *request.RequestContext, input *ListBucketsInput) (*ListBucketsOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	result, err := o.svc.listBucketsCore(store.buckets, AdminListBucketsInput{
		Prefix:            input.Prefix,
		BucketRegion:      input.BucketRegion,
		MaxBuckets:        input.MaxBuckets,
		ContinuationToken: input.ContinuationToken,
	})
	if err != nil {
		return nil, err
	}

	bucketInfos := make([]*BucketInfo, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		bucketInfos = append(bucketInfos, &BucketInfo{
			Name:         b.Name,
			CreationDate: b.CreationDate,
		})
	}

	return &ListBucketsOutput{
		Owner: &Owner{
			ID:          o.svc.accountID,
			DisplayName: o.svc.accountID,
		},
		Buckets:           bucketInfos,
		Prefix:            input.Prefix,
		ContinuationToken: result.ContinuationToken,
		IsTruncated:       result.IsTruncated,
	}, nil
}
