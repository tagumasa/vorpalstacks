package testutil

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"
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
	if err := s3CreateBucket(ctx, client, name); err != nil {
		return "", nil, err
	}
	if err := s3EnableVersioning(ctx, client, name); err != nil {
		s3CleanupBucket(ctx, client, name)
		return "", nil, err
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
			// Governance-mode retentions are bypassed here: cleanup must be
			// able to empty a lock-enabled bucket (test fixtures set
			// short retentions), and the bypass is a no-op for unlocked
			// versions. Compliance-mode versions stay undeletable until
			// expiry — the compliance fixture therefore clears its own
			// retentions before this runs (see complianceFixture.remove).
			client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket:                    aws.String(bucket),
				Delete:                    &types.Delete{Objects: objs},
				BypassGovernanceRetention: aws.Bool(true),
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

// s3ListBucketsAll pages through ListBuckets. The account's bucket list
// is shared with every service running in the same regression pass, so a
// contains-style assertion cannot rely on a single response page.
func s3ListBucketsAll(ctx context.Context, client *s3.Client) ([]types.Bucket, error) {
	var all []types.Bucket
	var token *string
	for page := 0; page < 10000; page++ {
		resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{ContinuationToken: token})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Buckets...)
		if resp.ContinuationToken == nil || *resp.ContinuationToken == "" {
			return all, nil
		}
		token = resp.ContinuationToken
	}
	return nil, fmt.Errorf("ListBuckets pagination did not terminate")
}

// s3ListPartsAll pages through ListParts by part-number marker so part
// assertions see the whole upload, not one response page.
func s3ListPartsAll(ctx context.Context, client *s3.Client, bucket, key, uploadID string) ([]types.Part, error) {
	var all []types.Part
	var marker *string
	for page := 0; page < 10000; page++ {
		resp, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:           aws.String(bucket),
			Key:              aws.String(key),
			UploadId:         aws.String(uploadID),
			PartNumberMarker: marker,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Parts...)
		if !aws.ToBool(resp.IsTruncated) {
			return all, nil
		}
		marker = resp.NextPartNumberMarker
	}
	return nil, fmt.Errorf("ListParts pagination did not terminate")
}

// s3ListVersionsAll pages through ListObjectVersions by key/version marker
// so version assertions see every version of every key, not one page.
func s3ListVersionsAll(ctx context.Context, client *s3.Client, bucket string) ([]types.ObjectVersion, []types.DeleteMarkerEntry, error) {
	var versions []types.ObjectVersion
	var markers []types.DeleteMarkerEntry
	var keyMarker *string
	var versionMarker *string
	for page := 0; page < 10000; page++ {
		resp, err := client.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{
			Bucket:          aws.String(bucket),
			KeyMarker:       keyMarker,
			VersionIdMarker: versionMarker,
		})
		if err != nil {
			return nil, nil, err
		}
		versions = append(versions, resp.Versions...)
		markers = append(markers, resp.DeleteMarkers...)
		if !aws.ToBool(resp.IsTruncated) {
			return versions, markers, nil
		}
		keyMarker = resp.NextKeyMarker
		versionMarker = resp.NextVersionIdMarker
	}
	return nil, nil, fmt.Errorf("ListObjectVersions pagination did not terminate")
}

// s3CreateBucket creates a plain bucket and wraps failures with the bucket
// name so a setup failure names the test's own resource. Tests whose
// subject is CreateBucket itself keep their literal inputs.
func s3CreateBucket(ctx context.Context, client *s3.Client, name string) error {
	if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}); err != nil {
		return fmt.Errorf("CreateBucket %s: %w", name, err)
	}
	return nil
}

// s3PutObject uploads a plain body (Bucket+Key+Body only). Sites carrying
// ContentType, Metadata, StorageClass, Tagging, ACL, or SSE members keep
// their literal PutObjectInput because those members are part of what the
// test exercises.
func s3PutObject(ctx context.Context, client *s3.Client, bucket, key, body string) (*s3.PutObjectOutput, error) {
	resp, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(body),
	})
	if err != nil {
		return nil, fmt.Errorf("PutObject %s/%s: %w", bucket, key, err)
	}
	return resp, nil
}

// s3HeadObject reads the plain HEAD metadata (Bucket+Key only); sites
// addressing a version, a range, or SSE-C keys keep literal inputs.
func s3HeadObject(ctx context.Context, client *s3.Client, bucket, key string) (*s3.HeadObjectOutput, error) {
	resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("HeadObject %s/%s: %w", bucket, key, err)
	}
	return resp, nil
}

// s3GetRead is the plain (Bucket+Key only) form of s3GetAndRead; ranged,
// versioned, and SSE-C reads keep their literal GetObjectInput.
func s3GetRead(ctx context.Context, client *s3.Client, bucket, key string) (*s3.GetObjectOutput, string, error) {
	return s3GetAndRead(ctx, client, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

// s3EnableVersioning turns versioning on for a bucket. Suspension calls
// and the versioning tests themselves keep literal inputs.
func s3EnableVersioning(ctx context.Context, client *s3.Client, bucket string) error {
	if _, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucket),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	}); err != nil {
		return fmt.Errorf("PutBucketVersioning %s: %w", bucket, err)
	}
	return nil
}

// s3ReplicationRule builds the single prefix-filtered, delete-marker-
// disabled rule the replication setup tests share; rules with tag filters,
// storage-class overrides, or delete-marker replication keep literal rules.
func s3ReplicationRule(id, destArn, prefix string) types.ReplicationRule {
	return types.ReplicationRule{
		ID:       aws.String(id),
		Status:   types.ReplicationRuleStatusEnabled,
		Priority: aws.Int32(1),
		Filter:   &types.ReplicationRuleFilter{Prefix: aws.String(prefix)},
		Destination: &types.Destination{
			Bucket: aws.String(destArn),
		},
		DeleteMarkerReplication: &types.DeleteMarkerReplication{
			Status: types.DeleteMarkerReplicationStatusDisabled,
		},
	}
}

// s3PutReplication installs a configuration with the shared replication
// test role; the role ARN shape is fixed by the replication engine's test
// role naming. Variadic so multi-rule configurations reuse it too.
func (r *TestRunner) s3PutReplication(ctx context.Context, client *s3.Client, bucket string, rules ...types.ReplicationRule) error {
	_, err := client.PutBucketReplication(ctx, &s3.PutBucketReplicationInput{
		Bucket: aws.String(bucket),
		ReplicationConfiguration: &types.ReplicationConfiguration{
			Role:  aws.String(fmt.Sprintf("arn:aws:iam::%s:role/s3-replication", r.accountID)),
			Rules: rules,
		},
	})
	if err != nil {
		return fmt.Errorf("PutBucketReplication %s: %w", bucket, err)
	}
	return nil
}

// expectS3Error asserts the AWS error code together with the HTTP status
// of the error response; the status is part of the observable contract
// wherever a test pins it.
func expectS3Error(err error, code string, status int) error {
	if err := expectAWSErrorCode(err, code); err != nil {
		return err
	}
	if got := awsHTTPStatus(err); got != status {
		return fmt.Errorf("expected HTTP %d for %s, got %d: %v", status, code, got, err)
	}
	return nil
}

// ssecTestKey derives the deterministic 32-byte SSE-C key material shared
// by the SSE-C tests, in its base64 and MD5 encodings.
func ssecTestKey() (encodedKey, encodedMD5 string) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	encodedKey = base64.StdEncoding.EncodeToString(key)
	keyMD5 := md5.Sum(key)
	encodedMD5 = base64.StdEncoding.EncodeToString(keyMD5[:])
	return encodedKey, encodedMD5
}

// s3PutTwoVersions writes the "first" and "second" versions of a key into
// a versioning-enabled bucket and returns the first version's output so
// tests can address it by VersionId.
func s3PutTwoVersions(ctx context.Context, client *s3.Client, bucket, key string) (*s3.PutObjectOutput, error) {
	putV1, err := s3PutObject(ctx, client, bucket, key, "first")
	if err != nil {
		return nil, fmt.Errorf("PutObject v1 failed: %w", err)
	}
	if _, err := s3PutObject(ctx, client, bucket, key, "second"); err != nil {
		return nil, fmt.Errorf("PutObject v2 failed: %w", err)
	}
	return putV1, nil
}

// s3HasPublicReadGrant reports whether the ACL grants contain the AllUsers
// group READ permission, the signature of a public-read canned ACL.
func s3HasPublicReadGrant(grants []types.Grant) bool {
	for _, grant := range grants {
		if grant.Grantee == nil || grant.Grantee.URI == nil {
			continue
		}
		if *grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" && grant.Permission == types.PermissionRead {
			return true
		}
	}
	return false
}

// s3BucketFixture owns a shared test bucket's lifecycle on the executed
// phase: the wrapper provisions the bucket before the first wrapped
// closure that runs needs it (creating an own bucket in the default
// region is the documented 200 OK path, so a leftover from an aborted run
// is adopted rather than failing the suite), and remove deletes it as the
// service cleanup. Registration itself carries no bucket side effects —
// the go test facade registers the whole suite before any subtest
// executes, and filtered -run selections still find the bucket in place.
type s3BucketFixture struct {
	ctx    context.Context
	client *s3.Client
	name   string
	// provision creates the bucket plus any required setup (versioning
	// for an inventory destination, say); it runs once, before the first
	// wrapped closure.
	provision func() error

	mu      sync.Mutex
	created bool
}

func (f *s3BucketFixture) wrapper(fn func() error) func() error {
	return func() error {
		if err := f.ensure(); err != nil {
			return err
		}
		return fn()
	}
}

func (f *s3BucketFixture) ensure() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.created {
		return nil
	}
	if err := f.provision(); err != nil {
		return err
	}
	f.created = true
	return nil
}

func (f *s3BucketFixture) remove() {
	s3CleanupBucket(f.ctx, f.client, f.name)
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

	// The main bucket lives on the executed phase: the families that
	// receive it register through the fixture wrapper (registration order
	// is preserved), and the registered cleanup removes it after the
	// service's tests finish — in binary mode at builder return, in
	// facade mode after the pending subtests.
	main := &s3BucketFixture{
		ctx:    ctx,
		client: client,
		name:   bucketName,
		provision: func() error {
			return s3CreateBucket(ctx, client, bucketName)
		},
	}
	r.RegisterServiceCleanup("s3", main.remove)
	r.PushClosureWrapper("s3", main.wrapper)
	results = append(results, r.s3BucketTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3ObjectTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3BucketConfigTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3ObjectConfigTests(ctx, client, ts, bucketName)...)
	results = append(results, r.s3MultibyteTests(ctx, client, ts, bucketName)...)
	r.PopClosureWrapper("s3")

	results = append(results, r.s3MultipartTests(ctx, client, ts)...)
	results = append(results, r.s3EncryptionTests(ctx, client, ts)...)

	r.PushClosureWrapper("s3", main.wrapper)
	results = append(results, r.s3AdvancedTests(ctx, client, ts, bucketName)...)
	r.PopClosureWrapper("s3")

	return results
}
