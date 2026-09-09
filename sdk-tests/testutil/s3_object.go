package testutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *TestRunner) s3ObjectTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "PutObject", func() error {
		resp, err := s3PutObject(ctx, client, bucketName, "test.txt", "Hello, World!")
		if err != nil {
			return err
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		return nil
	}))

	// Object keys are opaque UTF-8 per the model: dot segments are legal
	// input, distinct keys ("dir/../nested/key.txt" vs its neighbours,
	// "a//b.txt" vs "a/b.txt") must not clobber each other, and a key and
	// its prefix extension are independent objects that coexist.
	results = append(results, r.RunTest("s3", "PutObject_DotDotKeyAccepted", func() error {
		keys := []string{
			"file..backup.txt",
			"a..b",
			"dir/../nested/key.txt",
			"a//b.txt",
			// The flat namespace lets a key and its prefix extension
			// coexist: the first two write the file before its prefix
			// directory exists, the next two the other way round.
			"cofile.txt",
			"cofile.txt/nested.txt",
			"dirfirst/child.txt",
			"dirfirst",
		}
		for _, key := range keys {
			if _, err := s3PutObject(ctx, client, bucketName, key, "body of "+key); err != nil {
				return fmt.Errorf("PutObject(%q) failed: %w", key, err)
			}
		}
		if _, err := s3PutObject(ctx, client, bucketName, "nested/key.txt", "distinct body"); err != nil {
			return err
		}
		for _, key := range keys {
			_, body, err := s3GetRead(ctx, client, bucketName, key)
			if err != nil {
				return fmt.Errorf("GetObject(%q) failed: %w", key, err)
			}
			if body != "body of "+key {
				return fmt.Errorf("GetObject(%q) body = %q, want %q", key, body, "body of "+key)
			}
		}
		_, body, err := s3GetRead(ctx, client, bucketName, "nested/key.txt")
		if err != nil {
			return fmt.Errorf("GetObject(nested/key.txt) failed: %w", err)
		}
		if body != "distinct body" {
			return fmt.Errorf("nested/key.txt was clobbered by a traversal-shaped key: %q", body)
		}
		// The shared bucket must return to its prior emptiness for the
		// later delete-and-list test — deferred, so an assertion failure
		// above cannot leak the traversal-shaped keys into that test.
		defer func() {
			for _, key := range append(keys, "nested/key.txt") {
				_, _ = client.DeleteObject(ctx, &s3.DeleteObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(key),
				})
			}
		}()
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject", func() error {
		resp, gotBody, err := s3GetRead(ctx, client, bucketName, "test.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != "Hello, World!" {
			return fmt.Errorf("expected body %q, got %q", "Hello, World!", gotBody)
		}
		if resp.ContentLength == nil || *resp.ContentLength != 13 {
			return fmt.Errorf("expected ContentLength 13, got %v", resp.ContentLength)
		}
		if resp.ContentType == nil || *resp.ContentType == "" {
			return fmt.Errorf("ContentType is nil or empty")
		}
		if resp.LastModified == nil {
			return fmt.Errorf("LastModified is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject", func() error {
		resp, err := s3HeadObject(ctx, client, bucketName, "test.txt")
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.ContentLength == nil || *resp.ContentLength != 13 {
			return fmt.Errorf("expected ContentLength 13, got %v", resp.ContentLength)
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		if resp.ContentType == nil || *resp.ContentType == "" {
			return fmt.Errorf("ContentType is nil or empty")
		}
		if resp.LastModified == nil {
			return fmt.Errorf("LastModified is nil")
		}
		return nil
	}))

	// Conditional headers evaluate in the RFC 7232 section 6 order the S3
	// API reference points to: the two documented same-family pairs, and
	// the mixed If-None-Match + If-Unmodified-Since pair where the
	// unmodified-since failure (412) takes precedence over the matching
	// none-match (304).
	results = append(results, r.RunTest("s3", "GetObject_ConditionalPrecedence", func() error {
		head, err := s3HeadObject(ctx, client, bucketName, "test.txt")
		if err != nil {
			return fmt.Errorf("HeadObject for the object's ETag failed: %w", err)
		}
		etag := aws.ToString(head.ETag)
		before := aws.ToTime(head.LastModified).Add(-time.Hour)
		after := aws.ToTime(head.LastModified).Add(time.Hour)

		// Documented pair: a true If-Match suppresses a false
		// If-Unmodified-Since — the read succeeds.
		if _, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName), Key: aws.String("test.txt"),
			IfMatch: aws.String(etag), IfUnmodifiedSince: &after,
		}); err != nil {
			return fmt.Errorf("true If-Match with false If-Unmodified-Since must return the object: %w", err)
		}

		// Mixed pair: the false If-Unmodified-Since fails with 412 even
		// though If-None-Match would answer 304.
		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName), Key: aws.String("test.txt"),
			IfNoneMatch: aws.String(etag), IfUnmodifiedSince: &before,
		})
		if errS3 := expectS3Error(err, "PreconditionFailed", 412); errS3 != nil {
			return fmt.Errorf("mixed If-None-Match + false If-Unmodified-Since must fail 412: %w", errS3)
		}

		// Documented pair: a matching If-None-Match answers 304 regardless
		// of a true If-Modified-Since.
		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName), Key: aws.String("test.txt"),
			IfNoneMatch: aws.String(etag), IfModifiedSince: &before,
		})
		if errS3 := expectS3Error(err, "NotModified", 304); errS3 != nil {
			return fmt.Errorf("matching If-None-Match with true If-Modified-Since must be 304: %w", errS3)
		}

		// A differing If-None-Match ignores If-Modified-Since: the
		// unmodified-since state alone cannot produce a 304.
		if _, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName), Key: aws.String("test.txt"),
			IfNoneMatch: aws.String(`"differing"`), IfModifiedSince: &after,
		}); err != nil {
			return fmt.Errorf("differing If-None-Match must ignore If-Modified-Since and return the object: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject_DeleteMarkerVersion", func() error {
		verBucket, cleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "marker-read"))
		if err != nil {
			return fmt.Errorf("versioned bucket: %w", err)
		}
		defer cleanup()

		key := "marker-read.txt"
		if _, err := s3PutObject(ctx, client, verBucket, key, "body"); err != nil {
			return fmt.Errorf("put: %w", err)
		}
		del, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("delete: %w", err)
		}
		if del.DeleteMarker == nil || !*del.DeleteMarker {
			return fmt.Errorf("delete did not report a marker: %+v", del.DeleteMarker)
		}
		markerVersion := aws.ToString(del.VersionId)

		// A versionId addressing the delete marker itself is the
		// documented 405 surface, carrying the marker's Last-Modified and
		// identification headers.
		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(verBucket), Key: aws.String(key), VersionId: del.VersionId,
		})
		if errS3 := expectS3Error(err, "MethodNotAllowed", 405); errS3 != nil {
			return fmt.Errorf("GET of the marker version must be 405 MethodNotAllowed: %w", errS3)
		}
		var httpErr *awshttp.ResponseError
		if !errors.As(err, &httpErr) {
			return fmt.Errorf("GET marker version: no HTTP response on error: %v", err)
		}
		header := httpErr.Response.Header
		if got := header.Get("x-amz-delete-marker"); got != "true" {
			return fmt.Errorf("405 x-amz-delete-marker = %q, want true", got)
		}
		if got := header.Get("x-amz-version-id"); got != markerVersion {
			return fmt.Errorf("405 x-amz-version-id = %q, want %q", got, markerVersion)
		}
		if got := header.Get("Last-Modified"); got == "" {
			return fmt.Errorf("405 must carry the marker's Last-Modified")
		}

		// HEAD shares the contract; its error responses carry no body, so
		// the status and headers are the observables.
		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(verBucket), Key: aws.String(key), VersionId: del.VersionId,
		})
		if got := awsHTTPStatus(err); got != 405 {
			return fmt.Errorf("HEAD of the marker version = HTTP %d (%v), want 405", got, err)
		}
		if !errors.As(err, &httpErr) {
			return fmt.Errorf("HEAD marker version: no HTTP response on error: %v", err)
		}
		header = httpErr.Response.Header
		if got := header.Get("x-amz-delete-marker"); got != "true" {
			return fmt.Errorf("HEAD 405 x-amz-delete-marker = %q, want true", got)
		}
		if got := header.Get("Last-Modified"); got == "" {
			return fmt.Errorf("HEAD 405 must carry the marker's Last-Modified")
		}

		// The plain read of the marker-latest key stays a 404 that still
		// reports the marker identification headers.
		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(verBucket), Key: aws.String(key),
		})
		if errS3 := expectS3Error(err, "NoSuchKey", 404); errS3 != nil {
			return fmt.Errorf("plain GET of the marker-latest key must be 404 NoSuchKey: %w", errS3)
		}
		if !errors.As(err, &httpErr) {
			return fmt.Errorf("plain GET: no HTTP response on error: %v", err)
		}
		header = httpErr.Response.Header
		if got := header.Get("x-amz-delete-marker"); got != "true" {
			return fmt.Errorf("404 x-amz-delete-marker = %q, want true", got)
		}
		if got := header.Get("x-amz-version-id"); got != markerVersion {
			return fmt.Errorf("404 x-amz-version-id = %q, want %q", got, markerVersion)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if resp.Contents == nil {
			return fmt.Errorf("Contents is nil")
		}
		found := false
		for _, obj := range resp.Contents {
			if obj.Key != nil && *obj.Key == "test.txt" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("test.txt not found in listing")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err != nil {
			return fmt.Errorf("DeleteObject failed: %w", err)
		}

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected NoSuchKey error after delete, got nil")
		}
		var noSuchKey *types.NoSuchKey
		if !errors.As(err, &noSuchKey) {
			return fmt.Errorf("expected NoSuchKey, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsAfterDelete", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if len(resp.Contents) != 0 {
			return fmt.Errorf("expected 0 objects, got %d", len(resp.Contents))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CopyObject", func() error {
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:             aws.String(bucketName),
			Key:                aws.String("copy-source.txt"),
			Body:               strings.NewReader("copy me"),
			ContentEncoding:    aws.String("identity-src"),
			ContentDisposition: aws.String("attachment; filename=\"src.txt\""),
			ContentLanguage:    aws.String("en-US"),
			CacheControl:       aws.String("max-age=60"),
		}); err != nil {
			return err
		}

		_, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String("copy-dest.txt"),
			CopySource: aws.String(bucketName + "/copy-source.txt"),
		})
		if err != nil {
			return fmt.Errorf("CopyObject failed: %w", err)
		}

		head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("copy-dest.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject dest failed: %w", err)
		}
		// A plain copy carries the source's content-* metadata over.
		if aws.ToString(head.ContentEncoding) != "identity-src" {
			return fmt.Errorf("expected source ContentEncoding carried by copy, got %q", aws.ToString(head.ContentEncoding))
		}
		if aws.ToString(head.ContentDisposition) != "attachment; filename=\"src.txt\"" {
			return fmt.Errorf("expected source ContentDisposition carried by copy, got %q", aws.ToString(head.ContentDisposition))
		}
		if aws.ToString(head.ContentLanguage) != "en-US" {
			return fmt.Errorf("expected source ContentLanguage carried by copy, got %q", aws.ToString(head.ContentLanguage))
		}
		if aws.ToString(head.CacheControl) != "max-age=60" {
			return fmt.Errorf("expected source CacheControl carried by copy, got %q", aws.ToString(head.CacheControl))
		}

		// REPLACE installs the request's content-* values on the copy.
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:             aws.String(bucketName),
			Key:                aws.String("copy-dest.txt"),
			CopySource:         aws.String(bucketName + "/copy-dest.txt"),
			MetadataDirective:  types.MetadataDirectiveReplace,
			ContentEncoding:    aws.String("gzip"),
			ContentDisposition: aws.String("inline"),
			ContentLanguage:    aws.String("ja-JP"),
			CacheControl:       aws.String("no-cache"),
		})
		if err != nil {
			return fmt.Errorf("CopyObject with REPLACE failed: %w", err)
		}
		replaced, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("copy-dest.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject replaced dest failed: %w", err)
		}
		if aws.ToString(replaced.ContentEncoding) != "gzip" {
			return fmt.Errorf("expected replaced ContentEncoding gzip, got %q", aws.ToString(replaced.ContentEncoding))
		}
		if aws.ToString(replaced.ContentDisposition) != "inline" {
			return fmt.Errorf("expected replaced ContentDisposition inline, got %q", aws.ToString(replaced.ContentDisposition))
		}
		if aws.ToString(replaced.ContentLanguage) != "ja-JP" {
			return fmt.Errorf("expected replaced ContentLanguage ja-JP, got %q", aws.ToString(replaced.ContentLanguage))
		}
		if aws.ToString(replaced.CacheControl) != "no-cache" {
			return fmt.Errorf("expected replaced CacheControl no-cache, got %q", aws.ToString(replaced.CacheControl))
		}

		_, gotBody, err := s3GetRead(ctx, client, bucketName, "copy-dest.txt")
		if err != nil {
			return fmt.Errorf("GetObject dest failed: %w", err)
		}
		if gotBody != "copy me" {
			return fmt.Errorf("expected body %q, got %q", "copy me", gotBody)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_KeyWithHashDoesNotAliasVersion", func() error {
		// "#" is a legal object-key character. A key that embeds another
		// object's version ID must stay its own object: writing it must not
		// touch the versioned object's bytes, and deleting it must not
		// remove the version.
		bucket := s3Bucket(ts, "hash-key")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)
		if err := s3EnableVersioning(ctx, client, bucket); err != nil {
			return err
		}

		putResp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("aliased.txt"),
			Body:   strings.NewReader("version body"),
		})
		if err != nil {
			return fmt.Errorf("PutObject versioned failed: %w", err)
		}
		versionID := aws.ToString(putResp.VersionId)
		if versionID == "" {
			return fmt.Errorf("expected a version ID on the versioned put")
		}

		hashKey := "aliased.txt#" + versionID
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(hashKey),
			Body:   strings.NewReader("hash key body"),
		}); err != nil {
			return fmt.Errorf("PutObject(%q) failed: %w", hashKey, err)
		}

		getVersionResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:    aws.String(bucket),
			Key:       aws.String("aliased.txt"),
			VersionId: aws.String(versionID),
		})
		if err != nil {
			return fmt.Errorf("GetObject(version %s) failed: %w", versionID, err)
		}
		gotVersion, readErr := io.ReadAll(getVersionResp.Body)
		getVersionResp.Body.Close()
		if readErr != nil {
			return fmt.Errorf("GetObject(version %s) read failed: %w", versionID, readErr)
		}
		if string(gotVersion) != "version body" {
			return fmt.Errorf("hash-key put clobbered the versioned object: body = %q", gotVersion)
		}

		_, gotHash, err := s3GetRead(ctx, client, bucket, hashKey)
		if err != nil {
			return fmt.Errorf("GetObject(%q) failed: %w", hashKey, err)
		}
		if gotHash != "hash key body" {
			return fmt.Errorf("hash key body = %q, want %q", gotHash, "hash key body")
		}

		if _, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(hashKey),
		}); err != nil {
			return fmt.Errorf("DeleteObject(%q) failed: %w", hashKey, err)
		}

		afterResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:    aws.String(bucket),
			Key:       aws.String("aliased.txt"),
			VersionId: aws.String(versionID),
		})
		if err != nil {
			return fmt.Errorf("GetObject(version %s) after hash-key delete failed: %w", versionID, err)
		}
		afterDelete, readErr := io.ReadAll(afterResp.Body)
		afterResp.Body.Close()
		if readErr != nil {
			return fmt.Errorf("GetObject(version %s) read after delete failed: %w", versionID, readErr)
		}
		if string(afterDelete) != "version body" {
			return fmt.Errorf("hash-key delete removed the versioned object: body = %q", afterDelete)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CopyObject_VersionedSourceOnSuspendedBucket", func() error {
		// A version-addressed copy source is resolved by the record layout,
		// not the bucket's current versioning status: a suspended bucket
		// still serves its pre-suspension versions to the copier.
		bucket := s3Bucket(ts, "copy-ver")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)
		if err := s3EnableVersioning(ctx, client, bucket); err != nil {
			return err
		}

		putV1, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("versioned-source.txt"),
			Body:   strings.NewReader("version one"),
		})
		if err != nil {
			return fmt.Errorf("PutObject v1 failed: %w", err)
		}
		putV2, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("versioned-source.txt"),
			Body:   strings.NewReader("version two"),
		})
		if err != nil {
			return fmt.Errorf("PutObject v2 failed: %w", err)
		}
		if putV1.VersionId == nil || putV2.VersionId == nil || *putV1.VersionId == *putV2.VersionId {
			return fmt.Errorf("expected two distinct version IDs, got %v and %v", putV1.VersionId, putV2.VersionId)
		}

		if _, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		}); err != nil {
			return fmt.Errorf("PutBucketVersioning (suspend) failed: %w", err)
		}

		if _, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String("versioned-copy.txt"),
			CopySource: aws.String(bucket + "/versioned-source.txt?versionId=" + *putV1.VersionId),
		}); err != nil {
			return fmt.Errorf("CopyObject (versioned source) failed: %w", err)
		}

		_, gotBody, err := s3GetRead(ctx, client, bucket, "versioned-copy.txt")
		if err != nil {
			return fmt.Errorf("GetObject copy dest failed: %w", err)
		}
		if gotBody != "version one" {
			return fmt.Errorf("suspended-bucket versioned copy got %q, want the addressed version's body %q", gotBody, "version one")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CopyObject_MetadataDirective_REPLACEMetadata", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:   aws.String(bucketName),
			Key:      aws.String("metadata-src.txt"),
			Body:     strings.NewReader("original content"),
			Metadata: map[string]string{"custom-key": "original-value"},
		})
		if err != nil {
			return fmt.Errorf("PutObject source failed: %w", err)
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:            aws.String(bucketName),
			Key:               aws.String("metadata-dest.txt"),
			CopySource:        aws.String(bucketName + "/metadata-src.txt"),
			MetadataDirective: types.MetadataDirectiveReplace,
			Metadata:          map[string]string{"replaced-key": "new-value"},
		})
		if err != nil {
			return fmt.Errorf("CopyObject with REPLACE failed: %w", err)
		}

		resp, gotBody, err := s3GetRead(ctx, client, bucketName, "metadata-dest.txt")
		if err != nil {
			return fmt.Errorf("GetObject dest failed: %w", err)
		}
		if gotBody != "original content" {
			return fmt.Errorf("expected body content unchanged with REPLACE, got %q", gotBody)
		}
		if resp.Metadata["replaced-key"] != "new-value" {
			return fmt.Errorf("expected metadata replaced-key=new-value, got %v", resp.Metadata["replaced-key"])
		}
		if _, hasOldKey := resp.Metadata["custom-key"]; hasOldKey {
			return fmt.Errorf("expected old metadata custom-key to be removed with REPLACE, got %v", resp.Metadata["custom-key"])
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_GetObject_ContentVerification", func() error {
		content := "verification content"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String("verify.txt"),
			Body:          strings.NewReader(content),
			ContentType:   aws.String("text/plain"),
			ContentLength: aws.Int64(int64(len(content))),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, gotBody, err := s3GetRead(ctx, client, bucketName, "verify.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != content {
			return fmt.Errorf("expected body %q, got %q", content, gotBody)
		}
		if resp.ContentLength == nil || *resp.ContentLength != int64(len(content)) {
			return fmt.Errorf("expected ContentLength %d, got %v", len(content), resp.ContentLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_Overwrite", func() error {
		if _, err := s3PutObject(ctx, client, bucketName, "overwrite.txt", "first"); err != nil {
			return err
		}

		if _, err := s3PutObject(ctx, client, bucketName, "overwrite.txt", "second"); err != nil {
			return err
		}

		_, gotBody, err := s3GetRead(ctx, client, bucketName, "overwrite.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != "second" {
			return fmt.Errorf("expected body %q, got %q", "second", gotBody)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject_VerifyMetadata", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String("metadata.txt"),
			Body:        strings.NewReader("{}"),
			ContentType: aws.String("application/json"),
			Metadata: map[string]string{
				"custom-key": "custom-value",
			},
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := s3HeadObject(ctx, client, bucketName, "metadata.txt")
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("expected ContentType application/json, got %v", resp.ContentType)
		}
		val, ok := resp.Metadata["custom-key"]
		if !ok || val != "custom-value" {
			return fmt.Errorf("expected Metadata[custom-key]=custom-value, got %q, ok=%v", val, ok)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2_MultipleObjects", func() error {
		listBucket := s3Bucket(ts, "list")
		if err := s3CreateBucket(ctx, client, listBucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, listBucket)

		for i := 0; i < 5; i++ {
			if _, err := s3PutObject(ctx, client, listBucket, fmt.Sprintf("obj-%d.txt", i), fmt.Sprintf("data-%d", i)); err != nil {
				return err
			}
		}

		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(listBucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if len(resp.Contents) != 5 {
			return fmt.Errorf("expected 5 objects, got %d", len(resp.Contents))
		}
		for _, obj := range resp.Contents {
			if obj.Size == nil || *obj.Size == 0 {
				return fmt.Errorf("object %s has zero size", aws.ToString(obj.Key))
			}
			if obj.Owner != nil {
				return fmt.Errorf("object %s: V2 without fetch-owner must omit Owner", aws.ToString(obj.Key))
			}
		}

		// ?fetch-owner=true surfaces the bucket owner on every entry.
		ownerResp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:     aws.String(listBucket),
			FetchOwner: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 fetch-owner failed: %w", err)
		}
		if len(ownerResp.Contents) != 5 {
			return fmt.Errorf("fetch-owner: expected 5 objects, got %d", len(ownerResp.Contents))
		}
		for _, obj := range ownerResp.Contents {
			if obj.Owner == nil || obj.Owner.ID == nil || *obj.Owner.ID == "" {
				return fmt.Errorf("object %s: fetch-owner must populate Owner.ID", aws.ToString(obj.Key))
			}
		}
		return nil
	}))

	// With encodingType=url the server percent-encodes Key values and the
	// client URL-decodes them: keys containing XML-special characters
	// (&, <) and the escape character (%) must survive the round trip.
	results = append(results, r.RunTest("s3", "ListObjects_UrlEncodingTypeRoundTrip", func() error {
		urlEncBucket := s3Bucket(ts, "urlenc")
		if err := s3CreateBucket(ctx, client, urlEncBucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, urlEncBucket)

		keys := []string{"urlenc/a&b.txt", "urlenc/a<b.txt", "urlenc/a%b.txt"}
		for _, k := range keys {
			if _, err := s3PutObject(ctx, client, urlEncBucket, k, "url-encoded key data"); err != nil {
				return err
			}
		}

		decodedKeys := func(raws []*string) error {
			if len(raws) != len(keys) {
				return fmt.Errorf("expected %d objects, got %d", len(keys), len(raws))
			}
			for _, rawPtr := range raws {
				decoded, decErr := url.QueryUnescape(aws.ToString(rawPtr))
				if decErr != nil {
					return fmt.Errorf("key %q is not a valid url encoding: %v", aws.ToString(rawPtr), decErr)
				}
				found := false
				for _, k := range keys {
					if decoded == k {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("decoded key %q (raw %q) does not round-trip to a stored key", decoded, aws.ToString(rawPtr))
				}
			}
			return nil
		}

		v2, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:       aws.String(urlEncBucket),
			Prefix:       aws.String("urlenc/"),
			EncodingType: types.EncodingTypeUrl,
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 with url encoding failed: %w", err)
		}
		var v2Keys []*string
		for _, obj := range v2.Contents {
			v2Keys = append(v2Keys, obj.Key)
		}
		if err := decodedKeys(v2Keys); err != nil {
			return fmt.Errorf("ListObjectsV2: %w", err)
		}

		v1, err := client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket:       aws.String(urlEncBucket),
			Prefix:       aws.String("urlenc/"),
			EncodingType: types.EncodingTypeUrl,
		})
		if err != nil {
			return fmt.Errorf("ListObjects with url encoding failed: %w", err)
		}
		var v1Keys []*string
		for _, obj := range v1.Contents {
			v1Keys = append(v1Keys, obj.Key)
		}
		if err := decodedKeys(v1Keys); err != nil {
			return fmt.Errorf("ListObjects: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2_Pagination", func() error {
		pagBucket := s3Bucket(ts, "pag")
		if err := s3CreateBucket(ctx, client, pagBucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, pagBucket)

		for i := 0; i < 5; i++ {
			if _, err := s3PutObject(ctx, client, pagBucket, fmt.Sprintf("pag/obj-%d.txt", i), fmt.Sprintf("page-data-%d", i)); err != nil {
				return err
			}
		}

		var allKeys []string
		var continuationToken *string
		pageCount := 0
		for {
			pageCount++
			input := &s3.ListObjectsV2Input{
				Bucket:            aws.String(pagBucket),
				Prefix:            aws.String("pag/"),
				MaxKeys:           aws.Int32(2),
				ContinuationToken: continuationToken,
			}
			resp, err := client.ListObjectsV2(ctx, input)
			if err != nil {
				return fmt.Errorf("ListObjectsV2 page %d failed: %w", pageCount, err)
			}
			for _, obj := range resp.Contents {
				if obj.Key != nil {
					allKeys = append(allKeys, *obj.Key)
				}
			}
			if resp.IsTruncated == nil || !*resp.IsTruncated {
				break
			}
			continuationToken = resp.NextContinuationToken
		}
		if len(allKeys) != 5 {
			return fmt.Errorf("expected 5 total keys, got %d", len(allKeys))
		}
		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d", pageCount)
		}
		return nil
	}))

	// An explicit MaxKeys of zero requests an empty page; the S3 list APIs
	// document that the returned key count never exceeds the requested
	// limit, and an empty page is not truncated.
	results = append(results, r.RunTest("s3", "ListObjectsV2_MaxKeysZero", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:  aws.String(bucketName),
			MaxKeys: aws.Int32(0),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 with MaxKeys=0 failed: %w", err)
		}
		if resp.KeyCount == nil || *resp.KeyCount != 0 {
			return fmt.Errorf("expected KeyCount 0, got %v", resp.KeyCount)
		}
		if len(resp.Contents) != 0 {
			return fmt.Errorf("expected no Contents with MaxKeys=0, got %d", len(resp.Contents))
		}
		if resp.IsTruncated == nil || *resp.IsTruncated {
			return fmt.Errorf("expected IsTruncated false with MaxKeys=0, got %v", resp.IsTruncated)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjects_MaxKeysZero", func() error {
		resp, err := client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket:  aws.String(bucketName),
			MaxKeys: aws.Int32(0),
		})
		if err != nil {
			return fmt.Errorf("ListObjects with MaxKeys=0 failed: %w", err)
		}
		if len(resp.Contents) != 0 {
			return fmt.Errorf("expected no Contents with MaxKeys=0, got %d", len(resp.Contents))
		}
		if resp.IsTruncated == nil || *resp.IsTruncated {
			return fmt.Errorf("expected IsTruncated false with MaxKeys=0, got %v", resp.IsTruncated)
		}
		if resp.NextMarker != nil && *resp.NextMarker != "" {
			return fmt.Errorf("expected no NextMarker with MaxKeys=0, got %q", *resp.NextMarker)
		}
		return nil
	}))

	// The ListObjects API reference states that NextMarker "is returned
	// only if you have the delimiter request parameter specified".
	results = append(results, r.RunTest("s3", "ListObjects_NextMarkerOnlyWithDelimiter", func() error {
		v1Bucket := s3Bucket(ts, "v1marker")
		if err := s3CreateBucket(ctx, client, v1Bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, v1Bucket)

		for i := 0; i < 3; i++ {
			if _, err := s3PutObject(ctx, client, v1Bucket, fmt.Sprintf("plain-key-%d.txt", i), fmt.Sprintf("body-%d", i)); err != nil {
				return err
			}
		}

		noDelimiter, err := client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket:  aws.String(v1Bucket),
			MaxKeys: aws.Int32(2),
		})
		if err != nil {
			return fmt.Errorf("ListObjects without delimiter failed: %w", err)
		}
		if noDelimiter.IsTruncated == nil || !*noDelimiter.IsTruncated {
			return fmt.Errorf("expected truncated response, got %v", noDelimiter.IsTruncated)
		}
		if noDelimiter.NextMarker != nil && *noDelimiter.NextMarker != "" {
			return fmt.Errorf("expected no NextMarker without delimiter, got %q", *noDelimiter.NextMarker)
		}

		if _, err := s3PutObject(ctx, client, v1Bucket, "folder/nested.txt", "nested"); err != nil {
			return fmt.Errorf("PutObject nested failed: %w", err)
		}

		withDelimiter, err := client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket:    aws.String(v1Bucket),
			Delimiter: aws.String("/"),
			MaxKeys:   aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("ListObjects with delimiter failed: %w", err)
		}
		if withDelimiter.IsTruncated == nil || !*withDelimiter.IsTruncated {
			return fmt.Errorf("expected truncated response with delimiter, got %v", withDelimiter.IsTruncated)
		}
		if withDelimiter.NextMarker == nil || *withDelimiter.NextMarker == "" {
			return fmt.Errorf("expected NextMarker with delimiter on truncated response")
		}

		// A basic V1 roundtrip also covers the pre-V2 list path end to end:
		// three plain keys plus the nested key, which is only rolled up
		// into a common prefix when a delimiter is requested.
		all, err := client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(v1Bucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjects full page failed: %w", err)
		}
		if len(all.Contents) != 4 {
			return fmt.Errorf("expected 4 keys in full page, got %d", len(all.Contents))
		}
		// V1 returns the owner on every Contents entry unconditionally;
		// the fetch-owner opt-in exists only in V2.
		for _, c := range all.Contents {
			if c.Owner == nil || c.Owner.ID == nil || *c.Owner.ID == "" {
				return fmt.Errorf("expected Owner on V1 Contents entry for key %q", aws.ToString(c.Key))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectVersions", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "versions"))
		if err != nil {
			return err
		}
		defer verCleanup()

		for i := 0; i < 3; i++ {
			if _, err := s3PutObject(ctx, client, verBucket, "versioned-key.txt", fmt.Sprintf("version-%d", i)); err != nil {
				return fmt.Errorf("PutObject version %d failed: %w", i, err)
			}
		}

		resp, err := client.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{
			Bucket: aws.String(verBucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjectVersions failed: %w", err)
		}
		if len(resp.Versions) != 3 {
			return fmt.Errorf("expected 3 versions, got %d", len(resp.Versions))
		}
		// Version and delete-marker entries carry the owner unconditionally.
		for _, v := range resp.Versions {
			if v.Owner == nil || v.Owner.ID == nil || *v.Owner.ID == "" {
				return fmt.Errorf("expected Owner on Version entry for key %q", aws.ToString(v.Key))
			}
		}

		if _, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String("versioned-key.txt"),
		}); err != nil {
			return fmt.Errorf("DeleteObject failed: %w", err)
		}
		withMarker, err := client.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{
			Bucket: aws.String(verBucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjectVersions after delete failed: %w", err)
		}
		if len(withMarker.DeleteMarkers) != 1 {
			return fmt.Errorf("expected 1 delete marker, got %d", len(withMarker.DeleteMarkers))
		}
		if m := withMarker.DeleteMarkers[0]; m.Owner == nil || m.Owner.ID == nil || *m.Owner.ID == "" {
			return fmt.Errorf("expected Owner on DeleteMarker entry")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject_NonExistentKey", func() error {
		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-key.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var noSuchKey *types.NoSuchKey
		if !errors.As(err, &noSuchKey) {
			return fmt.Errorf("expected NoSuchKey, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject_NonExistentKey", func() error {
		_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-key.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var notFound *types.NotFound
		if !errors.As(err, &notFound) {
			return fmt.Errorf("expected NotFound, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject_NonExistentKey", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-delete-key.txt"),
		})
		if err != nil {
			return fmt.Errorf("DeleteObject on non-existent key should not error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject_VersionIdOnNonVersionedBucket", func() error {
		// A bucket that never had versioning only ever holds the null
		// version, so a version-addressed delete references a version that
		// does not exist: NoSuchVersion, not a server-side failure.
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket:    aws.String(bucketName),
			Key:       aws.String("versioned-delete-test.txt"),
			VersionId: aws.String("test-version-id"),
		})
		return expectAWSErrorCode(err, "NoSuchVersion")
	}))

	results = append(results, r.RunTest("s3", "PutObject_SystemMetadata", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:             aws.String(bucketName),
			Key:                aws.String("sysmeta.txt"),
			Body:               strings.NewReader("system metadata test"),
			ContentEncoding:    aws.String("gzip"),
			CacheControl:       aws.String("max-age=3600"),
			ContentDisposition: aws.String("attachment; filename=\"test.txt\""),
			ContentLanguage:    aws.String("en-US"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		headResp, err := s3HeadObject(ctx, client, bucketName, "sysmeta.txt")
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ContentEncoding == nil || *headResp.ContentEncoding != "gzip" {
			return fmt.Errorf("expected ContentEncoding gzip, got %v", headResp.ContentEncoding)
		}
		if headResp.CacheControl == nil || *headResp.CacheControl != "max-age=3600" {
			return fmt.Errorf("expected CacheControl max-age=3600, got %v", headResp.CacheControl)
		}
		if headResp.ContentDisposition == nil || *headResp.ContentDisposition != "attachment; filename=\"test.txt\"" {
			return fmt.Errorf("expected ContentDisposition attachment; filename=\"test.txt\", got %v", headResp.ContentDisposition)
		}
		if headResp.ContentLanguage == nil || *headResp.ContentLanguage != "en-US" {
			return fmt.Errorf("expected ContentLanguage en-US, got %v", headResp.ContentLanguage)
		}

		getResp, _, err := s3GetRead(ctx, client, bucketName, "sysmeta.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if getResp.ContentEncoding == nil || *getResp.ContentEncoding != "gzip" {
			return fmt.Errorf("GetObject: expected ContentEncoding gzip, got %v", getResp.ContentEncoding)
		}
		if getResp.CacheControl == nil || *getResp.CacheControl != "max-age=3600" {
			return fmt.Errorf("GetObject: expected CacheControl max-age=3600, got %v", getResp.CacheControl)
		}
		return nil
	}))

	return results
}
