package s3

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"sort"

	"vorpalstacks/internal/services/aws/common/request"
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

const (
	maxBucketNameLength = 63
	minBucketNameLength = 3
)

var bucketNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
var ipAddressRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

var validCORSMethods = map[string]bool{
	"GET":    true,
	"PUT":    true,
	"HEAD":   true,
	"POST":   true,
	"DELETE": true,
}

func validateBucketName(name string) error {
	if len(name) < minBucketNameLength || len(name) > maxBucketNameLength {
		return fmt.Errorf("bucket name must be between 3 and 63 characters")
	}
	if !bucketNameRegex.MatchString(name) {
		return fmt.Errorf("bucket name can only contain lowercase letters, numbers, dots, and hyphens")
	}
	if strings.HasPrefix(name, "xn--") {
		return fmt.Errorf("bucket name cannot start with 'xn--' (punycode)")
	}
	if strings.HasSuffix(name, "-s3alias") {
		return fmt.Errorf("bucket name cannot end with '-s3alias'")
	}
	if strings.HasSuffix(name, "--ol-s3") {
		return fmt.Errorf("bucket name cannot end with '--ol-s3'")
	}
	if strings.HasSuffix(name, ".mrap") {
		return fmt.Errorf("bucket name cannot end with '.mrap'")
	}
	if ipAddressRegex.MatchString(name) {
		return fmt.Errorf("bucket name cannot be formatted as an IP address")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("bucket name cannot contain consecutive periods")
	}
	if strings.Contains(name, ".-") || strings.Contains(name, "-.") {
		return fmt.Errorf("bucket name cannot contain periods adjacent to hyphens")
	}
	return nil
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
		return nil, fmt.Errorf("bucket name is required")
	}

	if err := validateBucketName(input.Bucket); err != nil {
		return nil, err
	}

	region := ctx.GetRegion()
	if input.LocationConstraint != "" && input.LocationConstraint != "us-east-1" {
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
		acp, err := cannedACLToPolicy(input.ACL, &s3store.ACLOwner{ID: o.svc.accountID, DisplayName: o.svc.accountID})
		if err != nil {
			return nil, err
		}
		bucket.ACL = acp
		if err := store.buckets.Put(bucket); err != nil {
			return nil, err
		}
	}

	if input.ObjectLockEnabledForBucket {
		bucket.ObjectLockEnabled = true
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
	count, err := store.objects.CountByBucket(input.Bucket)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("bucket is not empty")
	}

	multipartCount, err := store.objects.CountMultipartUploadsByBucket(input.Bucket)
	if err != nil {
		return err
	}
	if multipartCount > 0 {
		return fmt.Errorf("bucket has ongoing multipart uploads")
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

type HeadBucketOutput struct {
	BucketRegion string
}

func (o *BucketOperations) HeadBucket(ctx *request.RequestContext, input *HeadBucketInput) (*HeadBucketOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, fmt.Errorf("bucket not found")
	}
	return &HeadBucketOutput{
		BucketRegion: bucket.Region,
	}, nil
}

// ListBucketsInput contains the input parameters for the ListBuckets operation.
type ListBucketsInput struct{}

// ListBucketsOutput contains the output result of the ListBuckets operation.
type ListBucketsOutput struct {
	Owner   *Owner
	Buckets []*BucketInfo
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
func (o *BucketOperations) ListBuckets(ctx *request.RequestContext, input *ListBucketsInput) (*ListBucketsOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	buckets, err := store.buckets.List()
	if err != nil {
		return nil, err
	}

	bucketInfos := make([]*BucketInfo, 0)
	for _, b := range buckets {
		bucketInfos = append(bucketInfos, &BucketInfo{
			Name:         b.Name,
			CreationDate: b.CreationDate,
		})
	}

	sort.Slice(bucketInfos, func(i, j int) bool {
		return bucketInfos[i].Name < bucketInfos[j].Name
	})

	return &ListBucketsOutput{
		Owner: &Owner{
			ID:          o.svc.accountID,
			DisplayName: o.svc.accountID,
		},
		Buckets: bucketInfos,
	}, nil
}
