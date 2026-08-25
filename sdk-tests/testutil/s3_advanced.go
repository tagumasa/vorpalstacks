package testutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *TestRunner) s3AdvancedTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "RangeGet_VerifyPartialContent", func() error {
		content := "0123456789ABCDEFGHIJ"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("range-test.txt"),
			Body:   strings.NewReader(content),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("range-test.txt"),
			Range:  aws.String("bytes=0-4"),
		})
		if err != nil {
			return fmt.Errorf("GetObject with Range failed: %w", err)
		}
		if gotBody != "01234" {
			return fmt.Errorf("expected body %q, got %q", "01234", gotBody)
		}
		if resp.ContentRange == nil {
			return fmt.Errorf("ContentRange is nil")
		}
		if !strings.Contains(*resp.ContentRange, "bytes 0-4/20") {
			return fmt.Errorf("expected ContentRange to contain %q, got %q", "bytes 0-4/20", *resp.ContentRange)
		}
		return nil
	}))

	// The response-* query parameters override the corresponding response
	// headers on successful GET responses.
	results = append(results, r.RunTest("s3", "GetObject_ResponseHeaderOverrides", func() error {
		key := "response-override.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        strings.NewReader("override me"),
			ContentType: aws.String("text/plain"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:                     aws.String(bucketName),
			Key:                        aws.String(key),
			ResponseContentType:        aws.String("application/x-test-override"),
			ResponseCacheControl:       aws.String("max-age=60"),
			ResponseContentDisposition: aws.String("attachment; filename=overridden.txt"),
			ResponseContentLanguage:    aws.String("en-GB"),
		})
		if err != nil {
			return fmt.Errorf("GetObject with response overrides failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.ContentType == nil || *resp.ContentType != "application/x-test-override" {
			return fmt.Errorf("expected overridden Content-Type, got %v", resp.ContentType)
		}
		if resp.CacheControl == nil || *resp.CacheControl != "max-age=60" {
			return fmt.Errorf("expected overridden Cache-Control, got %v", resp.CacheControl)
		}
		if resp.ContentDisposition == nil || *resp.ContentDisposition != "attachment; filename=overridden.txt" {
			return fmt.Errorf("expected overridden Content-Disposition, got %v", resp.ContentDisposition)
		}
		if resp.ContentLanguage == nil || *resp.ContentLanguage != "en-GB" {
			return fmt.Errorf("expected overridden Content-Language, got %v", resp.ContentLanguage)
		}
		return nil
	}))

	// Every storage class accepted on PUT must round-trip its class value
	// through HEAD metadata together with the object content length.
	results = append(results, r.RunTest("s3", "StorageClass_HeadRoundTrip", func() error {
		for _, sc := range []struct {
			name  string
			class types.StorageClass
			key   string
			body  string
		}{
			{"STANDARD_IA", types.StorageClassStandardIa, "sc-ia.txt", "standard-ia"},
			{"GLACIER", types.StorageClassGlacier, "sc-glacier.txt", "glacier"},
			{"ONEZONE_IA", types.StorageClassOnezoneIa, "sc-1ia.txt", "onezone-ia"},
			{"INTELLIGENT_TIERING", types.StorageClassIntelligentTiering, "sc-it.txt", "intelligent-tiering"},
			{"REDUCED_REDUNDANCY", types.StorageClassReducedRedundancy, "sc-rr.txt", "reduced-redundancy"},
		} {
			if _, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:       aws.String(bucketName),
				Key:          aws.String(sc.key),
				Body:         strings.NewReader(sc.body),
				StorageClass: sc.class,
			}); err != nil {
				return fmt.Errorf("PutObject %s failed: %w", sc.name, err)
			}

			resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(sc.key),
			})
			if err != nil {
				return fmt.Errorf("HeadObject %s failed: %w", sc.name, err)
			}
			if resp.StorageClass != sc.class {
				return fmt.Errorf("%s: expected StorageClass %s, got %s", sc.name, sc.class, resp.StorageClass)
			}
			if resp.ContentLength == nil || *resp.ContentLength != int64(len(sc.body)) {
				return fmt.Errorf("%s: expected ContentLength %d, got %v", sc.name, len(sc.body), resp.ContentLength)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObjects_MultiDelete", func() error {
		multiDelBucket := s3Bucket(ts, "multidel")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(multiDelBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, multiDelBucket)

		keys := []string{"del-a.txt", "del-b.txt", "del-c.txt"}
		for _, k := range keys {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(multiDelBucket),
				Key:    aws.String(k),
				Body:   strings.NewReader("delete-me"),
			})
			if err != nil {
				return fmt.Errorf("PutObject %s failed: %w", k, err)
			}
		}

		var objs []types.ObjectIdentifier
		for _, k := range keys {
			objs = append(objs, types.ObjectIdentifier{Key: aws.String(k)})
		}

		delResp, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(multiDelBucket),
			Delete: &types.Delete{Objects: objs},
		})
		if err != nil {
			return fmt.Errorf("DeleteObjects failed: %w", err)
		}

		if len(delResp.Deleted) != 3 {
			return fmt.Errorf("expected 3 deleted, got %d", len(delResp.Deleted))
		}
		for _, k := range keys {
			found := false
			for _, d := range delResp.Deleted {
				if d.Key != nil && *d.Key == k {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("key %s not found in Deleted response", k)
			}
		}

		listResp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(multiDelBucket),
		})
		if err != nil {
			return fmt.Errorf("ListObjectsV2 failed: %w", err)
		}
		if len(listResp.Contents) != 0 {
			return fmt.Errorf("expected 0 objects after multi-delete, got %d", len(listResp.Contents))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PresignedURL_GetObject", func() error {
		presigner := s3.NewPresignClient(client)
		presignedReq, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("range-test.txt"),
		})
		if err != nil {
			return fmt.Errorf("PresignGetObject failed: %w", err)
		}

		httpResp, err := http.Get(presignedReq.URL)
		if err != nil {
			return fmt.Errorf("http.Get failed: %w", err)
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected status 200, got %d", httpResp.StatusCode)
		}

		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != "0123456789ABCDEFGHIJ" {
			return fmt.Errorf("expected body %q, got %q", "0123456789ABCDEFGHIJ", string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject_GlacierInvalidObjectState", func() error {
		key := "glacier-gate.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String(key),
			Body:         strings.NewReader("archived"),
			StorageClass: types.StorageClassGlacier,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err == nil {
			return fmt.Errorf("expected InvalidObjectState for archived object, got nil")
		}
		if err := expectAWSErrorCode(err, "InvalidObjectState"); err != nil {
			return err
		}
		// The S3 API reference documents InvalidObjectState with HTTP 403.
		if code := awsHTTPStatus(err); code != http.StatusForbidden {
			return fmt.Errorf("expected HTTP 403 for InvalidObjectState, got %d", code)
		}

		// HEAD remains available for archived objects.
		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("HeadObject on archived object failed: %w", err)
		}
		if headResp.StorageClass != types.StorageClassGlacier {
			return fmt.Errorf("expected StorageClass GLACIER on HEAD, got %s", headResp.StorageClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CopyObject_GlacierSourceNotInActiveTier", func() error {
		srcKey := "glacier-copy-src.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String(srcKey),
			Body:         strings.NewReader("archived source"),
			StorageClass: types.StorageClassGlacier,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String("glacier-copy-dst.txt"),
			CopySource: aws.String(bucketName + "/" + srcKey),
		})
		if err == nil {
			return fmt.Errorf("expected ObjectNotInActiveTierError for unrestored source, got nil")
		}
		if err := expectAWSErrorCode(err, "ObjectNotInActiveTierError"); err != nil {
			return err
		}
		// The S3 API reference documents ObjectNotInActiveTierError with
		// HTTP 403.
		if code := awsHTTPStatus(err); code != http.StatusForbidden {
			return fmt.Errorf("expected HTTP 403 for ObjectNotInActiveTierError, got %d", code)
		}

		// After restoring the source the copy succeeds.
		_, err = client.RestoreObject(ctx, &s3.RestoreObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(srcKey),
			RestoreRequest: &types.RestoreRequest{
				Days: aws.Int32(1),
			},
		})
		if err != nil {
			return fmt.Errorf("RestoreObject failed: %w", err)
		}
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String("glacier-copy-dst.txt"),
			CopySource: aws.String(bucketName + "/" + srcKey),
		})
		if err != nil {
			return fmt.Errorf("CopyObject after restore failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "RestoreObject_StandardClass", func() error {
		key := "restore-standard.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader("standard object"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.RestoreObject(ctx, &s3.RestoreObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			RestoreRequest: &types.RestoreRequest{
				Days: aws.Int32(1),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error restoring STANDARD object, got nil")
		}
		// Restoring an object that is not archived is rejected with
		// InvalidObjectState ("Restore is not allowed for the object's
		// current storage class"), reported with HTTP 403.
		if err := expectAWSErrorCode(err, "InvalidObjectState"); err != nil {
			return err
		}
		if code := awsHTTPStatus(err); code != http.StatusForbidden {
			return fmt.Errorf("expected HTTP 403 for InvalidObjectState, got %d", code)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "RestoreObject_SpecificVersion", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "restore-ver"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "versioned-restore.txt"
		putArchived, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(verBucket),
			Key:          aws.String(key),
			Body:         strings.NewReader("archived version"),
			StorageClass: types.StorageClassGlacier,
		})
		if err != nil {
			return fmt.Errorf("PutObject archived version failed: %w", err)
		}
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("standard version"),
		})
		if err != nil {
			return fmt.Errorf("PutObject standard version failed: %w", err)
		}

		headBefore, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putArchived.VersionId,
		})
		if err != nil {
			return fmt.Errorf("HeadObject on archived version failed: %w", err)
		}
		if headBefore.StorageClass != types.StorageClassGlacier {
			return fmt.Errorf("expected archived version GLACIER before restore, got %s", headBefore.StorageClass)
		}

		_, err = client.RestoreObject(ctx, &s3.RestoreObjectInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putArchived.VersionId,
			RestoreRequest: &types.RestoreRequest{
				Days: aws.Int32(1),
			},
		})
		if err != nil {
			return fmt.Errorf("RestoreObject with VersionId failed: %w", err)
		}

		headAfter, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putArchived.VersionId,
		})
		if err != nil {
			return fmt.Errorf("HeadObject on restored version failed: %w", err)
		}
		// The restore creates a temporary copy: the storage class stays
		// GLACIER and the restore status is reported via x-amz-restore.
		if headAfter.StorageClass != types.StorageClassGlacier {
			return fmt.Errorf("expected restored version to keep GLACIER, got %s", headAfter.StorageClass)
		}
		if headAfter.Restore == nil ||
			!strings.Contains(*headAfter.Restore, `ongoing-request="false"`) ||
			!strings.Contains(*headAfter.Restore, "expiry-date=") {
			return fmt.Errorf("expected x-amz-restore with expiry-date on restored version, got %v", headAfter.Restore)
		}

		latest, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("HeadObject on latest version failed: %w", err)
		}
		if latest.StorageClass != types.StorageClassStandard || latest.Restore != nil {
			return fmt.Errorf("expected untouched latest version (STANDARD, no restore header), got class=%s restore=%v", latest.StorageClass, latest.Restore)
		}
		return nil
	}))

	// A restored archive object serves its content on GET while the storage
	// class stays GLACIER and the response reports the restore status.
	results = append(results, r.RunTest("s3", "GetObject_RestoredArchiveServesContent", func() error {
		key := "restored-archive.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String(key),
			Body:         strings.NewReader("frozen data"),
			StorageClass: types.StorageClassGlacier,
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err == nil {
			return fmt.Errorf("expected InvalidObjectState for unrestored archive object, got nil")
		}
		if err := expectAWSErrorCode(err, "InvalidObjectState"); err != nil {
			return fmt.Errorf("expected InvalidObjectState for unrestored archive object: %v", err)
		}

		if _, err := client.RestoreObject(ctx, &s3.RestoreObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			RestoreRequest: &types.RestoreRequest{
				Days: aws.Int32(1),
			},
		}); err != nil {
			return fmt.Errorf("RestoreObject failed: %w", err)
		}

		get, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("GetObject after restore failed: %w", err)
		}
		if gotBody != "frozen data" {
			return fmt.Errorf("expected restored content, got %q", gotBody)
		}
		if get.Restore == nil || !strings.Contains(*get.Restore, `ongoing-request="false"`) {
			return fmt.Errorf("expected x-amz-restore on GET, got %v", get.Restore)
		}
		if get.StorageClass != types.StorageClassGlacier {
			return fmt.Errorf("expected storage class to stay GLACIER while restored, got %s", get.StorageClass)
		}
		return nil
	}))

	// Re-issuing a restore on an already-restored object extends the
	// temporary copy's expiry relative to the current time.
	results = append(results, r.RunTest("s3", "RestoreObject_ReRestoreExtendsExpiry", func() error {
		key := "re-restore.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String(key),
			Body:         strings.NewReader("cold"),
			StorageClass: types.StorageClassGlacier,
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		restore := func(days int32) error {
			_, err := client.RestoreObject(ctx, &s3.RestoreObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				RestoreRequest: &types.RestoreRequest{
					Days: aws.Int32(days),
				},
			})
			return err
		}
		if err := restore(1); err != nil {
			return fmt.Errorf("initial RestoreObject failed: %w", err)
		}
		first, err := restoreExpiryDate(ctx, client, bucketName, key)
		if err != nil {
			return err
		}
		if err := restore(3); err != nil {
			return fmt.Errorf("second RestoreObject failed: %w", err)
		}
		second, err := restoreExpiryDate(ctx, client, bucketName, key)
		if err != nil {
			return err
		}
		if second.Sub(first) < 24*time.Hour {
			return fmt.Errorf("expected extended expiry-date, first=%s second=%s", first, second)
		}
		return nil
	}))

	// Object tagging operations honour the versionId parameter: they read
	// and write the tag set of the addressed version, not the current one.
	results = append(results, r.RunTest("s3", "GetObjectTagging_VersionedObject", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "tagging-ver"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "versioned-tagging.txt"
		putV1, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("first"),
		})
		if err != nil {
			return fmt.Errorf("PutObject v1 failed: %w", err)
		}
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("second"),
		}); err != nil {
			return fmt.Errorf("PutObject v2 failed: %w", err)
		}

		if _, err := client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putV1.VersionId,
			Tagging: &types.Tagging{TagSet: []types.Tag{
				{Key: aws.String("stage"), Value: aws.String("original")},
			}},
		}); err != nil {
			return fmt.Errorf("PutObjectTagging with VersionId failed: %w", err)
		}
		if _, err := client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Tagging: &types.Tagging{TagSet: []types.Tag{
				{Key: aws.String("stage"), Value: aws.String("current")},
			}},
		}); err != nil {
			return fmt.Errorf("PutObjectTagging on latest failed: %w", err)
		}

		v1Tags, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putV1.VersionId,
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging with VersionId failed: %w", err)
		}
		if got := tagSetValue(v1Tags.TagSet, "stage"); got != "original" {
			return fmt.Errorf("expected original tag on v1, got %q (set %v)", got, v1Tags.TagSet)
		}

		latestTags, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging on latest failed: %w", err)
		}
		if got := tagSetValue(latestTags.TagSet, "stage"); got != "current" {
			return fmt.Errorf("expected current tag on latest, got %q (set %v)", got, latestTags.TagSet)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObjectTagging_SpecificVersionOnlyUpdatesThatVersion", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "tagging-put-ver"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "versioned-tagging-put.txt"
		putV1, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("first"),
		})
		if err != nil {
			return fmt.Errorf("PutObject v1 failed: %w", err)
		}
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("second"),
		}); err != nil {
			return fmt.Errorf("PutObject v2 failed: %w", err)
		}

		if _, err := client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putV1.VersionId,
			Tagging: &types.Tagging{TagSet: []types.Tag{
				{Key: aws.String("stage"), Value: aws.String("original")},
			}},
		}); err != nil {
			return fmt.Errorf("PutObjectTagging with VersionId failed: %w", err)
		}

		latestTags, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging on latest failed: %w", err)
		}
		if len(latestTags.TagSet) != 0 {
			return fmt.Errorf("expected latest version untagged, got %v", latestTags.TagSet)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObjectTagging_SpecificVersion", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "tagging-del-ver"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "versioned-tagging-del.txt"
		putV1, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("first"),
		})
		if err != nil {
			return fmt.Errorf("PutObject v1 failed: %w", err)
		}
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("second"),
		}); err != nil {
			return fmt.Errorf("PutObject v2 failed: %w", err)
		}

		if _, err := client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putV1.VersionId,
			Tagging: &types.Tagging{TagSet: []types.Tag{
				{Key: aws.String("stage"), Value: aws.String("original")},
			}},
		}); err != nil {
			return fmt.Errorf("PutObjectTagging on v1 failed: %w", err)
		}
		if _, err := client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Tagging: &types.Tagging{TagSet: []types.Tag{
				{Key: aws.String("stage"), Value: aws.String("current")},
			}},
		}); err != nil {
			return fmt.Errorf("PutObjectTagging on latest failed: %w", err)
		}

		if _, err := client.DeleteObjectTagging(ctx, &s3.DeleteObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putV1.VersionId,
		}); err != nil {
			return fmt.Errorf("DeleteObjectTagging with VersionId failed: %w", err)
		}

		v1Tags, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: putV1.VersionId,
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging with VersionId failed: %w", err)
		}
		if len(v1Tags.TagSet) != 0 {
			return fmt.Errorf("expected v1 untagged after delete, got %v", v1Tags.TagSet)
		}

		latestTags, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging on latest failed: %w", err)
		}
		if got := tagSetValue(latestTags.TagSet, "stage"); got != "current" {
			return fmt.Errorf("expected latest tag preserved, got %q (set %v)", got, latestTags.TagSet)
		}
		return nil
	}))

	// Addressing a version that does not exist is NoSuchVersion, not
	// NoSuchKey: the key exists, only the requested version is missing.
	results = append(results, r.RunTest("s3", "GetObject_NonExistentVersion", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "get-nover"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "nonexistent-version.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("body"),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: aws.String("nonexistent-version-id-0000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected NoSuchVersion for nonexistent version, got nil")
		}
		if err := expectAWSErrorCode(err, "NoSuchVersion"); err != nil {
			return err
		}
		if code := awsHTTPStatus(err); code != http.StatusNotFound {
			return fmt.Errorf("expected HTTP 404 for NoSuchVersion, got %d", code)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject_NonExistentVersion", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "head-nover"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "nonexistent-version.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("body"),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: aws.String("nonexistent-version-id-0000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected an error for nonexistent version, got nil")
		}
		// HEAD error responses carry no body, so the SDK surfaces the
		// modelled NotFound error for any 404; the HTTP status is the
		// observable contract here.
		if err := expectAWSErrorCode(err, "NotFound"); err != nil {
			return err
		}
		if code := awsHTTPStatus(err); code != http.StatusNotFound {
			return fmt.Errorf("expected HTTP 404 for nonexistent version, got %d", code)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObjectTagging_NonExistentVersion", func() error {
		verBucket, verCleanup, err := s3CreateVersionedBucket(ctx, client, s3Bucket(ts, "tagging-nover"))
		if err != nil {
			return err
		}
		defer verCleanup()

		key := "nonexistent-version.txt"
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(verBucket),
			Key:    aws.String(key),
			Body:   strings.NewReader("body"),
		}); err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket:    aws.String(verBucket),
			Key:       aws.String(key),
			VersionId: aws.String("nonexistent-version-id-0000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected NoSuchVersion for nonexistent version, got nil")
		}
		if err := expectAWSErrorCode(err, "NoSuchVersion"); err != nil {
			return err
		}
		if code := awsHTTPStatus(err); code != http.StatusNotFound {
			return fmt.Errorf("expected HTTP 404 for NoSuchVersion, got %d", code)
		}
		return nil
	}))

	// HeadObject evaluates the same conditional headers as GetObject:
	// failed If-Match / If-Unmodified-Since give 412 PreconditionFailed,
	// matching If-None-Match / unmodified If-Modified-Since give 304.
	results = append(results, r.RunTest("s3", "HeadObject_ConditionalRequests", func() error {
		key := "head-conditional.txt"
		put, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader("head conditional body"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:  aws.String(bucketName),
			Key:     aws.String(key),
			IfMatch: put.ETag,
		})
		if err != nil {
			return fmt.Errorf("HeadObject with matching If-Match failed: %w", err)
		}

		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:  aws.String(bucketName),
			Key:     aws.String(key),
			IfMatch: aws.String(`"00000000000000000000000000000000"`),
		})
		if err == nil {
			return fmt.Errorf("expected PreconditionFailed for mismatched If-Match, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusPreconditionFailed {
			return fmt.Errorf("expected HTTP 412 for mismatched If-Match, got %d: %v", code, err)
		}

		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			IfNoneMatch: put.ETag,
		})
		if err == nil {
			return fmt.Errorf("expected NotModified for matching If-None-Match, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusNotModified {
			return fmt.Errorf("expected HTTP 304 for matching If-None-Match, got %d: %v", code, err)
		}

		future := time.Now().Add(time.Hour)
		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:          aws.String(bucketName),
			Key:             aws.String(key),
			IfModifiedSince: &future,
		})
		if err == nil {
			return fmt.Errorf("expected NotModified for future If-Modified-Since, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusNotModified {
			return fmt.Errorf("expected HTTP 304 for future If-Modified-Since, got %d: %v", code, err)
		}

		past := time.Now().Add(-24 * time.Hour)
		_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:            aws.String(bucketName),
			Key:               aws.String(key),
			IfUnmodifiedSince: &past,
		})
		if err == nil {
			return fmt.Errorf("expected PreconditionFailed for past If-Unmodified-Since, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusPreconditionFailed {
			return fmt.Errorf("expected HTTP 412 for past If-Unmodified-Since, got %d: %v", code, err)
		}
		return nil
	}))

	// The x-amz-copy-source-if-* headers are preconditions on the copy
	// source object; a failed precondition fails the copy with 412
	// PreconditionFailed.
	results = append(results, r.RunTest("s3", "CopyObject_CopySourceIfConditions", func() error {
		srcKey := "copy-if-src.txt"
		head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(srcKey),
		})
		if err == nil && head != nil {
			// Source exists from an earlier run; refresh it.
			_, _ = client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(srcKey),
			})
		}
		put, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(srcKey),
			Body:   strings.NewReader("copy source body"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		src := bucketName + "/" + srcKey

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:            aws.String(bucketName),
			Key:               aws.String("copy-if-dst.txt"),
			CopySource:        aws.String(src),
			CopySourceIfMatch: aws.String(`"00000000000000000000000000000000"`),
		})
		if err == nil {
			return fmt.Errorf("expected PreconditionFailed for mismatched copy-source If-Match, got nil")
		}
		if err := expectAWSErrorCode(err, "PreconditionFailed"); err != nil {
			return fmt.Errorf("expected PreconditionFailed: %v", err)
		}
		if code := awsHTTPStatus(err); code != http.StatusPreconditionFailed {
			return fmt.Errorf("expected HTTP 412 for mismatched copy-source If-Match, got %d", code)
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:            aws.String(bucketName),
			Key:               aws.String("copy-if-dst.txt"),
			CopySource:        aws.String(src),
			CopySourceIfMatch: put.ETag,
		})
		if err != nil {
			return fmt.Errorf("CopyObject with matching If-Match failed: %w", err)
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:                aws.String(bucketName),
			Key:                   aws.String("copy-if-dst.txt"),
			CopySource:            aws.String(src),
			CopySourceIfNoneMatch: put.ETag,
		})
		if err == nil {
			return fmt.Errorf("expected PreconditionFailed for matching copy-source If-None-Match, got nil")
		}
		if code := awsHTTPStatus(err); code != http.StatusPreconditionFailed {
			return fmt.Errorf("expected HTTP 412 for matching copy-source If-None-Match, got %d", code)
		}
		return nil
	}))

	// UploadPartCopy copies a byte range of an existing object into a part
	// of an in-progress multipart upload.
	results = append(results, r.RunTest("s3", "UploadPartCopy_BasicRoundtrip", func() error {
		upcBucket := s3Bucket(ts, "upcopy")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(upcBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, upcBucket)

		srcKey := "upcopy-src.txt"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(upcBucket),
			Key:    aws.String(srcKey),
			Body:   strings.NewReader("0123456789"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		createResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(upcBucket),
			Key:    aws.String("upcopy-dst.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}

		copyResp, err := client.UploadPartCopy(ctx, &s3.UploadPartCopyInput{
			Bucket:          aws.String(upcBucket),
			Key:             aws.String("upcopy-dst.txt"),
			UploadId:        createResp.UploadId,
			PartNumber:      aws.Int32(1),
			CopySource:      aws.String(upcBucket + "/" + srcKey),
			CopySourceRange: aws.String("bytes=2-6"),
		})
		if err != nil {
			return fmt.Errorf("UploadPartCopy failed: %w", err)
		}
		if copyResp.CopyPartResult == nil || copyResp.CopyPartResult.ETag == nil {
			return fmt.Errorf("expected CopyPartResult with ETag, got %v", copyResp.CopyPartResult)
		}

		_, err = client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(upcBucket),
			Key:      aws.String("upcopy-dst.txt"),
			UploadId: createResp.UploadId,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{ETag: copyResp.CopyPartResult.ETag, PartNumber: aws.Int32(1)},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CompleteMultipartUpload failed: %w", err)
		}

		_, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket: aws.String(upcBucket),
			Key:    aws.String("upcopy-dst.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != "23456" {
			return fmt.Errorf("expected copied range %q, got %q", "23456", gotBody)
		}
		return nil
	}))

	// Copying a part from an unrestored archived source is rejected with
	// the same error the single-request copy reports.
	results = append(results, r.RunTest("s3", "UploadPartCopy_UnrestoredArchiveSourceRejected", func() error {
		upcBucket := s3Bucket(ts, "upcopy-arc")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(upcBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, upcBucket)

		srcKey := "glacier-src.txt"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(upcBucket),
			Key:          aws.String(srcKey),
			Body:         strings.NewReader("archived part source"),
			StorageClass: types.StorageClassGlacier,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		createResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(upcBucket),
			Key:    aws.String("glacier-dst.txt"),
		})
		if err != nil {
			return fmt.Errorf("CreateMultipartUpload failed: %w", err)
		}

		_, err = client.UploadPartCopy(ctx, &s3.UploadPartCopyInput{
			Bucket:     aws.String(upcBucket),
			Key:        aws.String("glacier-dst.txt"),
			UploadId:   createResp.UploadId,
			PartNumber: aws.Int32(1),
			CopySource: aws.String(upcBucket + "/" + srcKey),
		})
		if err == nil {
			return fmt.Errorf("expected ObjectNotInActiveTierError for unrestored source, got nil")
		}
		if err := expectAWSErrorCode(err, "ObjectNotInActiveTierError"); err != nil {
			return err
		}
		// The S3 API reference documents ObjectNotInActiveTierError with
		// HTTP 403.
		if code := awsHTTPStatus(err); code != http.StatusForbidden {
			return fmt.Errorf("expected HTTP 403 for ObjectNotInActiveTierError, got %d", code)
		}

		// After restoring the source the part copy succeeds.
		_, err = client.RestoreObject(ctx, &s3.RestoreObjectInput{
			Bucket: aws.String(upcBucket),
			Key:    aws.String(srcKey),
			RestoreRequest: &types.RestoreRequest{
				Days: aws.Int32(1),
			},
		})
		if err != nil {
			return fmt.Errorf("RestoreObject failed: %w", err)
		}
		_, err = client.UploadPartCopy(ctx, &s3.UploadPartCopyInput{
			Bucket:     aws.String(upcBucket),
			Key:        aws.String("glacier-dst.txt"),
			UploadId:   createResp.UploadId,
			PartNumber: aws.Int32(1),
			CopySource: aws.String(upcBucket + "/" + srcKey),
		})
		if err != nil {
			return fmt.Errorf("UploadPartCopy after restore failed: %w", err)
		}
		return nil
	}))

	// A HEAD with a Range header reports the partial metadata (206,
	// Content-Range, ranged Content-Length) without a body.
	results = append(results, r.RunTest("s3", "HeadObject_Range", func() error {
		key := "head-range.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader("0123456789"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Range:  aws.String("bytes=2-5"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject with Range failed: %w", err)
		}
		if resp.ContentLength == nil || *resp.ContentLength != 4 {
			return fmt.Errorf("expected ContentLength 4, got %v", resp.ContentLength)
		}
		if resp.ContentRange == nil || *resp.ContentRange != "bytes 2-5/10" {
			return fmt.Errorf("expected ContentRange bytes 2-5/10, got %v", resp.ContentRange)
		}
		return nil
	}))

	return results
}

// tagSetValue returns the value of the given tag key, or the empty string
// when the key is absent from the set.
func tagSetValue(tagSet []types.Tag, key string) string {
	for _, tag := range tagSet {
		if aws.ToString(tag.Key) == key {
			return aws.ToString(tag.Value)
		}
	}
	return ""
}

// restoreExpiryDate reads the expiry-date from the x-amz-restore header of
// a HEAD response, failing when the header or the date is absent.
func restoreExpiryDate(ctx context.Context, client *s3.Client, bucket, key string) (time.Time, error) {
	head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("HeadObject failed: %w", err)
	}
	if head.Restore == nil {
		return time.Time{}, errors.New("missing x-amz-restore header")
	}
	const marker = `expiry-date="`
	idx := strings.Index(*head.Restore, marker)
	if idx < 0 {
		return time.Time{}, fmt.Errorf("no expiry-date in %q", *head.Restore)
	}
	rest := (*head.Restore)[idx+len(marker):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return time.Time{}, fmt.Errorf("unterminated expiry-date in %q", *head.Restore)
	}
	return http.ParseTime(rest[:end])
}
