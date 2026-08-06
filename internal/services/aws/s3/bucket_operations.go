package s3

import (
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
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
	LocationConstraint         string
	ObjectLockEnabledForBucket bool
}

// CreateBucketOutput contains the output result of the CreateBucket operation.
type CreateBucketOutput struct {
	Location string
}

// CreateBucket creates a new bucket.
func (o *BucketOperations) CreateBucket(ctx *request.RequestContext, input *CreateBucketInput) (*CreateBucketOutput, error) {
	if input.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket name is required")
	}

	if err := validateBucketName(input.Bucket); err != nil {
		return nil, err
	}

	region := ctx.GetRegion()
	if input.LocationConstraint != "" && input.LocationConstraint != defaults.DefaultRegion {
		region = input.LocationConstraint
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	bucket, err := store.buckets.Create(input.Bucket, region)
	if err != nil {
		return nil, err
	}

	if input.ACL != "" {
		acp, err := CannedACLToPolicy(input.ACL, &s3store.ACLOwner{ID: o.svc.accountID, DisplayName: o.svc.accountID})
		if err != nil {
			return nil, err
		}
		bucket.ACL = acp
	}

	if input.ObjectLockEnabledForBucket {
		bucket.ObjectLockEnabled = true
	}

	if input.ACL != "" || input.ObjectLockEnabledForBucket {
		if err := store.buckets.Put(bucket); err != nil {
			return nil, err
		}
	}

	return &CreateBucketOutput{
		Location: "/" + input.Bucket,
	}, nil
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

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if bucket.ObjectLockEnabled {
		logs.Warn("s3: deleting bucket with Object Lock enabled", logs.String("bucket", input.Bucket))
	}

	count, err := store.objects.CountByBucket(input.Bucket)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrBucketNotEmpty
	}

	multipartCount, err := store.objects.CountMultipartUploadsByBucket(input.Bucket)
	if err != nil {
		return err
	}
	if multipartCount > 0 {
		return ErrBucketNotEmpty
	}

	o.svc.encryptionManager.DeleteBucketKey(input.Bucket)

	return store.buckets.Delete(input.Bucket)
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
	return store.buckets.Get(input.Bucket)
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
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}
	return &HeadBucketOutput{
		BucketRegion: bucket.Region,
	}, nil
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
	buckets, err := store.buckets.List()
	if err != nil {
		return nil, err
	}

	bucketInfos := make([]*BucketInfo, 0, len(buckets))
	for _, b := range buckets {
		if input.Prefix != "" && !strings.HasPrefix(b.Name, input.Prefix) {
			continue
		}
		if input.BucketRegion != "" && b.Region != input.BucketRegion {
			continue
		}
		bucketInfos = append(bucketInfos, &BucketInfo{
			Name:         b.Name,
			CreationDate: b.CreationDate,
		})
	}

	sort.Slice(bucketInfos, func(i, j int) bool {
		return bucketInfos[i].Name < bucketInfos[j].Name
	})

	// Apply ContinuationToken (offset by name).
	startIdx := 0
	if input.ContinuationToken != "" {
		for i, b := range bucketInfos {
			if b.Name > input.ContinuationToken {
				startIdx = i
				break
			}
			// If we reach the end, no more results.
			if i == len(bucketInfos)-1 {
				startIdx = len(bucketInfos)
			}
		}
	}

	// Default page size is 10000 per AWS spec when no pagination params given.
	// Smithy range for MaxBuckets is 1-10000.
	maxBuckets := input.MaxBuckets
	if maxBuckets <= 0 || maxBuckets > 10000 {
		maxBuckets = 10000
	}

	endIdx := startIdx + maxBuckets
	var nextToken string
	isTruncated := false
	if endIdx < len(bucketInfos) {
		nextToken = bucketInfos[endIdx-1].Name
		isTruncated = true
	} else {
		endIdx = len(bucketInfos)
	}

	paged := bucketInfos[startIdx:endIdx]

	return &ListBucketsOutput{
		Owner: &Owner{
			ID:          o.svc.accountID,
			DisplayName: o.svc.accountID,
		},
		Buckets:           paged,
		Prefix:            input.Prefix,
		ContinuationToken: nextToken,
		IsTruncated:       isTruncated,
	}, nil
}
