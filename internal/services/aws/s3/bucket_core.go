package s3

import (
	"context"
	"errors"
	"sort"
	"strings"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// Bucket lifecycle and listing cores: create, delete and list buckets and
// their objects for both the admin gRPC-Web console and the AWS API plane.
// Object read/write/delete cores live in object_core.go, the streaming
// put/get/copy engines in object_stream_core.go.

// AdminListBucketsInput carries the fields needed for ListBuckets.
type AdminListBucketsInput struct {
	Prefix            string
	BucketRegion      string
	MaxBuckets        int
	ContinuationToken string
}

// AdminCreateBucketInput carries the fields needed for CreateBucket.
type AdminCreateBucketInput struct {
	Bucket                     string
	Region                     string
	ACL                        string
	GrantFullControl           string
	GrantRead                  string
	GrantReadACP               string
	GrantWrite                 string
	GrantWriteACP              string
	ObjectOwnership            string
	ObjectLockEnabledForBucket bool
}

// AdminDeleteBucketInput carries the fields needed for DeleteBucket.
type AdminDeleteBucketInput struct {
	Bucket string
}

// AdminListObjectsInput carries the fields needed for ListObjectsV2.
type AdminListObjectsInput struct {
	Bucket    string
	Prefix    string
	Delimiter string
	Marker    string
	MaxKeys   int
}

// AdminListBucketsResult holds the transport-agnostic result of ListBuckets.
// Returns raw store buckets so each transport layer can format as needed.
type AdminListBucketsResult struct {
	Buckets           []*s3store.Bucket
	ContinuationToken string
	IsTruncated       bool
}

// AdminCreateBucketResult holds the transport-agnostic result of CreateBucket.
type AdminCreateBucketResult struct {
	Location string
}

// AdminDeleteBucketResult holds the transport-agnostic result of DeleteBucket.
type AdminDeleteBucketResult struct{}

// AdminListObjectsResult holds the transport-agnostic result of ListObjectsV2.
type AdminListObjectsResult struct {
	Objects        []*s3store.Object
	CommonPrefixes []string
	IsTruncated    bool
	NextMarker     string
}

// listBucketsCore returns all buckets in the regional store, optionally
// filtered by prefix and bucket-region, sorted by name, and paginated.
func (s *S3Service) listBucketsCore(bucketStore s3store.BucketStoreInterface, in AdminListBucketsInput) (*AdminListBucketsResult, error) {
	buckets, err := bucketStore.List()
	if err != nil {
		return nil, err
	}

	filtered := make([]*s3store.Bucket, 0, len(buckets))
	for _, b := range buckets {
		if in.Prefix != "" && !strings.HasPrefix(b.Name, in.Prefix) {
			continue
		}
		if in.BucketRegion != "" && b.Region != in.BucketRegion {
			continue
		}
		filtered = append(filtered, b)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name < filtered[j].Name
	})

	startIdx := 0
	if in.ContinuationToken != "" {
		for i, b := range filtered {
			if b.Name > in.ContinuationToken {
				startIdx = i
				break
			}
			if i == len(filtered)-1 {
				startIdx = len(filtered)
			}
		}
	}

	maxBuckets := in.MaxBuckets
	if maxBuckets <= 0 || maxBuckets > s3MaxBuckets {
		maxBuckets = s3MaxBuckets
	}

	endIdx := startIdx + maxBuckets
	var nextToken string
	isTruncated := false
	if endIdx < len(filtered) {
		nextToken = filtered[endIdx-1].Name
		isTruncated = true
	} else {
		endIdx = len(filtered)
	}

	return &AdminListBucketsResult{
		Buckets:           filtered[startIdx:endIdx],
		ContinuationToken: nextToken,
		IsTruncated:       isTruncated,
	}, nil
}

// createBucketCore validates the bucket name, creates the bucket in the
// regional store, and optionally applies a canned ACL and/or Object Lock.
func (s *S3Service) createBucketCore(bucketStore s3store.BucketStoreInterface, in AdminCreateBucketInput) (*AdminCreateBucketResult, error) {
	if in.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket name is required")
	}

	if err := validateBucketName(in.Bucket); err != nil {
		return nil, err
	}

	// Bucket names form a single global namespace and this platform has a
	// single account, so a duplicate name always collides with a bucket the
	// requester owns. S3 returns BucketAlreadyOwnedByYou for that case in
	// every Region except North Virginia, where legacy clients rely on a
	// 200 OK response that resets the existing bucket's ACLs.
	if bucketStore.Exists(in.Bucket) {
		if in.Region == defaults.DefaultRegion {
			existing, getErr := bucketStore.Get(in.Bucket)
			if getErr != nil {
				return nil, getErr
			}
			existing.ACL = nil
			if putErr := bucketStore.Put(existing); putErr != nil {
				return nil, putErr
			}
			return &AdminCreateBucketResult{Location: "/" + in.Bucket}, nil
		}
		return nil, ErrBucketAlreadyOwnedByYou
	}
	if s.s3Store != nil {
		if found, _ := s.s3Store.FindBucket(in.Bucket); found != nil {
			return nil, ErrBucketAlreadyOwnedByYou
		}
	}

	bucket, err := bucketStore.Create(in.Bucket, in.Region)
	if err != nil {
		if errors.Is(err, s3store.ErrBucketAlreadyExists) {
			return nil, ErrBucketAlreadyOwnedByYou
		}
		return nil, err
	}

	if in.ACL != "" || in.GrantFullControl != "" || in.GrantRead != "" || in.GrantReadACP != "" || in.GrantWrite != "" || in.GrantWriteACP != "" {
		acp, err := buildACLFromInput(in.ACL, nil, in.GrantFullControl, in.GrantRead, in.GrantReadACP, in.GrantWrite, in.GrantWriteACP,
			&s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID})
		if err != nil {
			return nil, err
		}
		bucket.ACL = acp
	}

	if in.ObjectOwnership != "" {
		if err := validateOwnershipControls([]OwnershipControlsRuleInput{{ObjectOwnership: in.ObjectOwnership}}); err != nil {
			return nil, err
		}
		// With ACLs disabled the bucket "accepts only PUT requests that do
		// not specify an ACL or PUT requests with bucket owner full control
		// ACLs"; other ACLs fail with 400 AccessControlListNotSupported.
		if in.ObjectOwnership == "BucketOwnerEnforced" {
			if in.ACL != "" && !aclAllowedWithAclsDisabled(in.ACL) {
				return nil, ErrAccessControlListNotSupported
			}
			if in.GrantFullControl != "" || in.GrantRead != "" || in.GrantReadACP != "" || in.GrantWrite != "" || in.GrantWriteACP != "" {
				return nil, ErrAccessControlListNotSupported
			}
		}
		bucket.OwnershipControls = &s3store.OwnershipControls{
			Rules: []s3store.OwnershipControlsRule{{ObjectOwnership: in.ObjectOwnership}},
		}
	}

	if in.ObjectLockEnabledForBucket {
		bucket.ObjectLockEnabled = true
	}

	if bucket.ACL != nil || bucket.OwnershipControls != nil || in.ObjectLockEnabledForBucket {
		if err := bucketStore.Put(bucket); err != nil {
			return nil, err
		}
	}

	return &AdminCreateBucketResult{Location: "/" + in.Bucket}, nil
}

// deleteBucketCore verifies the bucket is empty (no objects, no incomplete
// multipart uploads), purges the encryption key, and deletes the bucket.
func (s *S3Service) deleteBucketCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminDeleteBucketInput) (*AdminDeleteBucketResult, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}
	if bucket.ObjectLockEnabled {
		logs.Warn("s3: deleting bucket with Object Lock enabled", logs.String("bucket", in.Bucket))
	}

	count, err := objectStore.CountByBucket(in.Bucket)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrBucketNotEmpty
	}

	multipartCount, err := objectStore.CountMultipartUploadsByBucket(in.Bucket)
	if err != nil {
		return nil, err
	}
	if multipartCount > 0 {
		return nil, ErrBucketNotEmpty
	}

	if s.encryptionManager != nil {
		s.encryptionManager.DeleteBucketKey(in.Bucket)
	}

	if err := bucketStore.Delete(in.Bucket); err != nil {
		return nil, err
	}

	// The gates above proved no object records remain, so the blob tree
	// holds at most empty directories; removing it reclaims the skeleton
	// the object writes created. The record delete is already committed,
	// and a failure here leaves only orphan directories — the lesser
	// residue — so it is logged rather than failing the request.
	if s.blobStore != nil {
		if err := s.blobStore.DeleteBucket(ctx, in.Bucket); err != nil {
			logs.Warn("s3: blob directory removal failed after bucket delete",
				logs.String("bucket", in.Bucket), logs.Err(err))
		}
	}
	return &AdminDeleteBucketResult{}, nil
}

// listObjectsCore lists objects in a bucket with pagination.  MaxKeys is
// expected to be resolved by the caller (an absent limit becomes the
// 1000-page default at the transport edge); here it is only clamped.
func (s *S3Service) listObjectsCore(objectStore s3store.ObjectStoreInterface, in AdminListObjectsInput) (*AdminListObjectsResult, error) {
	maxKeys := in.MaxKeys
	if maxKeys < 0 {
		maxKeys = 0
	}
	if maxKeys > s3MaxKeys {
		maxKeys = s3MaxKeys
	}

	result, err := objectStore.List(in.Bucket, in.Prefix, in.Delimiter, in.Marker, maxKeys)
	if err != nil {
		return nil, err
	}

	return &AdminListObjectsResult{
		Objects:        result.Objects,
		CommonPrefixes: result.CommonPrefixes,
		IsTruncated:    result.IsTruncated,
		NextMarker:     result.NextMarker,
	}, nil
}
