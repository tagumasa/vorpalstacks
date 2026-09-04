package testutil

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"vorpalstacks-sdk-tests/config"
)

func s3Bucket(ts string, name string) string {
	return fmt.Sprintf("s3test-%s-%s", name, ts)
}

// waitForObject polls GetObject until the object exists or timeout expires.
// Returns true if the object was found.
func waitForObject(ctx context.Context, client *s3.Client, bucket, key string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err == nil {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

// waitForObjectDeleted polls GetObject until it returns an error (object
// deleted or delete marker created) or timeout expires. Returns true if
// the object is no longer accessible.
func waitForObjectDeleted(ctx context.Context, client *s3.Client, bucket, key string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

// s3GetAndRead performs a GetObject call and drains the response body into
// a string, closing the body before returning so callers cannot leak it.
// The response is returned alongside the body for callers that make further
// assertions on metadata fields (ContentLength, StorageClass, Restore, ...).
func s3GetAndRead(ctx context.Context, client *s3.Client, in *s3.GetObjectInput) (*s3.GetObjectOutput, string, error) {
	resp, err := client.GetObject(ctx, in)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("ReadAll: %w", err)
	}
	return resp, string(body), nil
}

// s3CreateVersionedBucket creates an empty bucket with versioning enabled
// and returns its name together with a cleanup closure that empties and
// deletes it. On a versioning failure the bucket is cleaned up immediately.
func s3CreateVersionedBucket(ctx context.Context, client *s3.Client, name string) (string, func(), error) {
	if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}); err != nil {
		return "", nil, fmt.Errorf("CreateBucket %s: %w", name, err)
	}
	if _, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: aws.String(name),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	}); err != nil {
		s3CleanupBucket(ctx, client, name)
		return "", nil, fmt.Errorf("PutBucketVersioning %s: %w", name, err)
	}
	return name, func() { s3CleanupBucket(ctx, client, name) }, nil
}

func s3CleanupBucket(ctx context.Context, client *s3.Client, bucket string) {
	mpuResp, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{Bucket: aws.String(bucket)})
	if err == nil {
		for _, u := range mpuResp.Uploads {
			client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      u.Key,
				UploadId: u.UploadId,
			})
		}
	}

	// Delete every object version and delete marker (the bucket may have
	// versioning enabled, in which case unversioned deletes would only add
	// markers and leave the bucket non-empty), walking all pages.
	var keyMarker *string
	var versionMarker *string
	for page := 0; page < 1000; page++ {
		listResp, err := client.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{
			Bucket:          aws.String(bucket),
			KeyMarker:       keyMarker,
			VersionIdMarker: versionMarker,
		})
		if err != nil {
			return
		}
		if len(listResp.Versions) == 0 && len(listResp.DeleteMarkers) == 0 {
			break
		}
		var objs []types.ObjectIdentifier
		for _, v := range listResp.Versions {
			objs = append(objs, types.ObjectIdentifier{Key: v.Key, VersionId: v.VersionId})
		}
		for _, m := range listResp.DeleteMarkers {
			objs = append(objs, types.ObjectIdentifier{Key: m.Key, VersionId: m.VersionId})
		}
		if len(objs) > 0 {
			client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &types.Delete{Objects: objs},
			})
		}
		keyMarker = listResp.NextKeyMarker
		versionMarker = listResp.NextVersionIdMarker
		if keyMarker == nil && versionMarker == nil {
			break
		}
	}
	client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
}

func (r *TestRunner) RunS3Tests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "s3",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	ctx := context.Background()
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	bucketName := s3Bucket(ts, "main")
	defer s3CleanupBucket(ctx, client, bucketName)

	results = append(results, r.s3BucketTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3ObjectTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3BucketConfigTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3ObjectConfigTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3MultipartTests(ctx, client, ts)...)
	results = append(results, r.s3MultibyteTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3EncryptionTests(ctx, client, ts)...)
	results = append(results, r.s3AdvancedTests(ctx, client, ts, bucketName)...)

	return results
}
