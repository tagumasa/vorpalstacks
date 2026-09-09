package testutil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func (r *TestRunner) s3ObjectConfigTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	// The lock bucket and the seed objects (tagging/attributes objects in
	// the main bucket, legal-hold/retention objects in the lock bucket)
	// live on the executed phase like the main bucket itself: registration
	// must not create, seed or delete anything, because the go test facade
	// registers the whole suite before any subtest runs. Provisioning runs
	// once, when the first wrapped closure executes; the main-bucket
	// wrapper stacked outside this one has ensured the main bucket exists
	// by then, so the seed puts land in an existing bucket.
	lockBucket := s3Bucket(ts, "lock")
	lockFixture := &s3BucketFixture{
		ctx:    ctx,
		client: client,
		name:   lockBucket,
		provision: func() error {
			if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket:                     aws.String(lockBucket),
				ObjectLockEnabledForBucket: aws.Bool(true),
			}); err != nil {
				return fmt.Errorf("CreateBucket (lock) failed: %w", err)
			}
			seeds := []struct{ bucket, key, body string }{
				{bucketName, "tagged-obj.txt", "tag me"},
				{bucketName, "attrs-obj.txt", "attributes content"},
				{lockBucket, "legal-hold-obj.txt", "legal hold content"},
				{lockBucket, "retention-obj.txt", "retention content"},
			}
			for _, seed := range seeds {
				if _, err := s3PutObject(ctx, client, seed.bucket, seed.key, seed.body); err != nil {
					return fmt.Errorf("seed PutObject %s failed: %w", seed.key, err)
				}
			}
			return nil
		},
	}
	// The legal hold must be lifted before the bucket sweep: unlike
	// governance retentions it is not bypassable, only releasable, and an
	// ON hold would keep the version (and the bucket) undeletable.
	r.RegisterServiceCleanup("s3", func() {
		_, _ = client.PutObjectLegalHold(ctx, &s3.PutObjectLegalHoldInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("legal-hold-obj.txt"),
			LegalHold: &types.ObjectLockLegalHold{
				Status: types.ObjectLockLegalHoldStatusOff,
			},
		})
	})
	r.RegisterServiceCleanup("s3", lockFixture.remove)
	r.PushClosureWrapper("s3", lockFixture.wrapper)
	defer r.PopClosureWrapper("s3")

	results = append(results, r.RunTest("s3", "PutObjectTagging_GetVerify", func() error {
		_, err := client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
			Tagging: &types.Tagging{
				TagSet: []types.Tag{
					{Key: aws.String("env"), Value: aws.String("prod")},
					{Key: aws.String("team"), Value: aws.String("backend")},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutObjectTagging failed: %w", err)
		}
		getResp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging failed: %w", err)
		}
		if len(getResp.TagSet) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(getResp.TagSet))
		}
		hasEnv := false
		hasTeam := false
		for _, tag := range getResp.TagSet {
			if tag.Key != nil && *tag.Key == "env" && tag.Value != nil && *tag.Value == "prod" {
				hasEnv = true
			}
			if tag.Key != nil && *tag.Key == "team" && tag.Value != nil && *tag.Value == "backend" {
				hasTeam = true
			}
		}
		if !hasEnv {
			return fmt.Errorf("tag env=prod not found in TagSet")
		}
		if !hasTeam {
			return fmt.Errorf("tag team=backend not found in TagSet")
		}
		return nil
	}))

	// The x-amz-tagging header is URL query-parameter encoded per the API
	// contract; the server decodes it, so a compliant client's tags read
	// back decoded.
	results = append(results, r.RunTest("s3", "PutObject_TaggingHeaderDecoded", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-header-obj.txt"),
			Body:   strings.NewReader("tagged via header"),
			// The SDK passes the Tagging string through as the
			// x-amz-tagging header, so the encoded form goes on the wire.
			Tagging: aws.String("a%20b=c%21d"),
		})
		if err != nil {
			return fmt.Errorf("PutObject with Tagging failed: %w", err)
		}
		getResp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-header-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging failed: %w", err)
		}
		if len(getResp.TagSet) != 1 {
			return fmt.Errorf("expected 1 tag, got %d", len(getResp.TagSet))
		}
		tag := getResp.TagSet[0]
		if aws.ToString(tag.Key) != "a b" || aws.ToString(tag.Value) != "c!d" {
			return fmt.Errorf("tag = %q=%q, want a b=c!d (decoded)", aws.ToString(tag.Key), aws.ToString(tag.Value))
		}
		return nil
	}))

	// The tag-set rules apply to the upload header exactly as they apply to
	// PutObjectTagging, and validation happens before the write: a rejected
	// header must leave no stored object behind.
	results = append(results, r.RunTest("s3", "PutObject_InvalidTaggingRejectedBeforeWrite", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:  aws.String(bucketName),
			Key:     aws.String("invalid-tagging-obj.txt"),
			Body:    strings.NewReader("must never be stored"),
			Tagging: aws.String("aws:reserved=denied"),
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag on upload, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "InvalidArgument" {
			return fmt.Errorf("expected InvalidArgument, got %s: %v", apiErr.ErrorCode(), err)
		}
		if code := awsHTTPStatus(err); code != http.StatusBadRequest {
			return fmt.Errorf("expected HTTP 400 for invalid tagging, got %d: %v", code, err)
		}
		_, getErr := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("invalid-tagging-obj.txt"),
		})
		if getErr == nil {
			return fmt.Errorf("object was stored despite the rejected tagging header")
		}
		if code := awsHTTPStatus(getErr); code != http.StatusNotFound {
			return fmt.Errorf("expected HTTP 404 after rejected put, got %d: %v", code, getErr)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObjectTagging_VerifyEmpty", func() error {
		_, err := client.DeleteObjectTagging(ctx, &s3.DeleteObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("DeleteObjectTagging failed: %w", err)
		}
		getResp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectTagging failed: %w", err)
		}
		if len(getResp.TagSet) != 0 {
			return fmt.Errorf("expected 0 tags after delete, got %d", len(getResp.TagSet))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObjectAcl_GetVerify", func() error {
		_, err := client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
			ACL:    types.ObjectCannedACLPrivate,
		})
		if err != nil {
			return fmt.Errorf("PutObjectAcl failed: %w", err)
		}
		getResp, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectAcl failed: %w", err)
		}
		if getResp.Owner == nil || getResp.Owner.ID == nil || *getResp.Owner.ID == "" {
			return fmt.Errorf("Owner is missing: %+v", getResp.Owner)
		}
		// The private canned ACL is the owner's FULL_CONTROL alone: any
		// other grant set means the canned ACL was not applied.
		if len(getResp.Grants) != 1 {
			return fmt.Errorf("expected exactly 1 grant for the private ACL, got %+v", getResp.Grants)
		}
		g := getResp.Grants[0]
		if g.Permission != types.PermissionFullControl || g.Grantee == nil || g.Grantee.ID == nil || *g.Grantee.ID != *getResp.Owner.ID {
			return fmt.Errorf("expected owner FULL_CONTROL grant, got %+v", g)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObjectLockConfiguration_GetVerify", func() error {
		_, err := client.PutObjectLockConfiguration(ctx, &s3.PutObjectLockConfigurationInput{
			Bucket: aws.String(lockBucket),
			ObjectLockConfiguration: &types.ObjectLockConfiguration{
				ObjectLockEnabled: types.ObjectLockEnabledEnabled,
				Rule: &types.ObjectLockRule{
					DefaultRetention: &types.DefaultRetention{
						Mode: types.ObjectLockRetentionModeGovernance,
						Days: aws.Int32(10),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutObjectLockConfiguration failed: %w", err)
		}
		getResp, err := client.GetObjectLockConfiguration(ctx, &s3.GetObjectLockConfigurationInput{
			Bucket: aws.String(lockBucket),
		})
		if err != nil {
			return fmt.Errorf("GetObjectLockConfiguration failed: %w", err)
		}
		if getResp.ObjectLockConfiguration == nil {
			return fmt.Errorf("ObjectLockConfiguration is nil")
		}
		if getResp.ObjectLockConfiguration.ObjectLockEnabled != types.ObjectLockEnabledEnabled {
			return fmt.Errorf("expected ObjectLockEnabled Enabled, got %s", getResp.ObjectLockConfiguration.ObjectLockEnabled)
		}
		if getResp.ObjectLockConfiguration.Rule == nil {
			return fmt.Errorf("Rule is nil")
		}
		dr := getResp.ObjectLockConfiguration.Rule.DefaultRetention
		if dr == nil {
			return fmt.Errorf("DefaultRetention is nil")
		}
		if dr.Mode != types.ObjectLockRetentionModeGovernance {
			return fmt.Errorf("expected DefaultRetention.Mode Governance, got %s", dr.Mode)
		}
		if dr.Days == nil || *dr.Days != 10 {
			return fmt.Errorf("expected DefaultRetention.Days 10, got %v", dr.Days)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObjectLegalHold_GetVerify", func() error {
		_, err := client.PutObjectLegalHold(ctx, &s3.PutObjectLegalHoldInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("legal-hold-obj.txt"),
			LegalHold: &types.ObjectLockLegalHold{
				Status: types.ObjectLockLegalHoldStatusOn,
			},
		})
		if err != nil {
			return fmt.Errorf("PutObjectLegalHold failed: %w", err)
		}
		getResp, err := client.GetObjectLegalHold(ctx, &s3.GetObjectLegalHoldInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("legal-hold-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectLegalHold failed: %w", err)
		}
		if getResp.LegalHold == nil {
			return fmt.Errorf("LegalHold is nil")
		}
		if getResp.LegalHold.Status != types.ObjectLockLegalHoldStatusOn {
			return fmt.Errorf("expected LegalHold Status ON, got %s", getResp.LegalHold.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObjectRetention_GetVerify", func() error {
		retainUntil := time.Now().Add(24 * time.Hour)
		_, err := client.PutObjectRetention(ctx, &s3.PutObjectRetentionInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("retention-obj.txt"),
			Retention: &types.ObjectLockRetention{
				Mode:            types.ObjectLockRetentionModeGovernance,
				RetainUntilDate: aws.Time(retainUntil),
			},
		})
		if err != nil {
			return fmt.Errorf("PutObjectRetention failed: %w", err)
		}
		getResp, err := client.GetObjectRetention(ctx, &s3.GetObjectRetentionInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("retention-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObjectRetention failed: %w", err)
		}
		if getResp.Retention == nil {
			return fmt.Errorf("Retention is nil")
		}
		if getResp.Retention.Mode != types.ObjectLockRetentionModeGovernance {
			return fmt.Errorf("expected Mode Governance, got %s", getResp.Retention.Mode)
		}
		if getResp.Retention.RetainUntilDate == nil {
			return fmt.Errorf("RetainUntilDate is nil")
		}
		if getResp.Retention.RetainUntilDate.Before(time.Now()) {
			return fmt.Errorf("RetainUntilDate is in the past: %v", getResp.Retention.RetainUntilDate)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ObjectLock_LegalHold_BlocksDelete", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("legal-hold-obj.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error deleting object with legal hold ON, got nil")
		}
		if err := expectS3Error(err, "AccessDenied", http.StatusForbidden); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ObjectLock_Retention_BlocksDelete", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("retention-obj.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error deleting object with active GOVERNANCE retention, got nil")
		}
		if err := expectS3Error(err, "AccessDenied", http.StatusForbidden); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ObjectLock_Governance_BypassDelete", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket:                    aws.String(lockBucket),
			Key:                       aws.String("retention-obj.txt"),
			BypassGovernanceRetention: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("expected bypass delete to succeed, got: %w", err)
		}
		return nil
	}))

	// The compliance-mode fixture follows the same executed-phase pattern:
	// the bucket, its locked object and the COMPLIANCE retention itself are
	// established when the first wrapped closure runs.
	complianceBucket := s3Bucket(ts, "compliance-lock")
	complianceFixture := &s3BucketFixture{
		ctx:    ctx,
		client: client,
		name:   complianceBucket,
		provision: func() error {
			if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket:                     aws.String(complianceBucket),
				ObjectLockEnabledForBucket: aws.Bool(true),
			}); err != nil {
				return fmt.Errorf("CreateBucket (compliance) failed: %w", err)
			}
			if _, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(complianceBucket),
				Key:    aws.String("compliance-obj.txt"),
				Body:   strings.NewReader("compliance content"),
			}); err != nil {
				return fmt.Errorf("seed PutObject compliance-obj.txt failed: %w", err)
			}
			// The retention window is the shortest future date the
			// scenario permits: COMPLIANCE mode can be neither shortened
			// nor bypassed, so the bucket outlives this run by exactly
			// this window — two minutes, not a day.
			if _, err := client.PutObjectRetention(ctx, &s3.PutObjectRetentionInput{
				Bucket: aws.String(complianceBucket),
				Key:    aws.String("compliance-obj.txt"),
				Retention: &types.ObjectLockRetention{
					Mode:            types.ObjectLockRetentionModeCompliance,
					RetainUntilDate: aws.Time(time.Now().Add(2 * time.Minute)),
				},
			}); err != nil {
				return fmt.Errorf("PutObjectRetention (compliance) failed: %w", err)
			}
			return nil
		},
	}
	r.RegisterServiceCleanup("s3", complianceFixture.remove)
	r.PushClosureWrapper("s3", complianceFixture.wrapper)
	defer r.PopClosureWrapper("s3")

	results = append(results, r.RunTest("s3", "ObjectLock_Compliance_NoBypass", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket:                    aws.String(complianceBucket),
			Key:                       aws.String("compliance-obj.txt"),
			BypassGovernanceRetention: aws.Bool(true),
		})
		if err == nil {
			return fmt.Errorf("expected error deleting COMPLIANCE retention object even with bypass, got nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObjectAttributes_VerifyAll", func() error {
		resp, err := client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("attrs-obj.txt"),
			ObjectAttributes: []types.ObjectAttributes{
				types.ObjectAttributesObjectSize,
				types.ObjectAttributesEtag,
				types.ObjectAttributesStorageClass,
			},
		})
		if err != nil {
			return fmt.Errorf("GetObjectAttributes failed: %w", err)
		}
		if resp.ObjectSize == nil || *resp.ObjectSize == 0 {
			return fmt.Errorf("expected ObjectSize > 0, got %v", resp.ObjectSize)
		}
		if *resp.ObjectSize != int64(len("attributes content")) {
			return fmt.Errorf("expected ObjectSize %d, got %d", len("attributes content"), *resp.ObjectSize)
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		if resp.StorageClass != types.StorageClassStandard {
			return fmt.Errorf("expected StorageClass STANDARD, got %s", resp.StorageClass)
		}

		// Unrequested fields are not returned: a checksum-only request must
		// carry no ETag, size or storage class.
		onlyChecksum, err := client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("attrs-obj.txt"),
			ObjectAttributes: []types.ObjectAttributes{
				types.ObjectAttributesChecksum,
			},
		})
		if err != nil {
			return fmt.Errorf("GetObjectAttributes (checksum only) failed: %w", err)
		}
		if onlyChecksum.ETag != nil || onlyChecksum.ObjectSize != nil || onlyChecksum.StorageClass != "" {
			return fmt.Errorf("unrequested attributes returned: ETag=%v ObjectSize=%v StorageClass=%s", onlyChecksum.ETag, onlyChecksum.ObjectSize, onlyChecksum.StorageClass)
		}

		// A value outside the modelled attribute set is InvalidArgument.
		_, err = client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("attrs-obj.txt"),
			ObjectAttributes: []types.ObjectAttributes{
				types.ObjectAttributes("Bogus"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for unmodelled attribute value, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "InvalidArgument" {
			return fmt.Errorf("expected InvalidArgument for unmodelled attribute, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	// Every storage class this platform persists must survive the
	// put→head round trip, and the hardware-bound classes of the AWS enum
	// are rejected at acceptance.
	results = append(results, r.RunTest("s3", "PutObject_StorageClassRoundTrip", func() error {
		accepted := []types.StorageClass{
			types.StorageClassStandard,
			types.StorageClassReducedRedundancy,
			types.StorageClassStandardIa,
			types.StorageClassOnezoneIa,
			types.StorageClassIntelligentTiering,
			types.StorageClassGlacier,
			types.StorageClassGlacierIr,
			types.StorageClassDeepArchive,
		}
		for _, class := range accepted {
			key := "sc-roundtrip/" + string(class) + ".txt"
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:       aws.String(bucketName),
				Key:          aws.String(key),
				Body:         strings.NewReader("storage class round trip"),
				StorageClass: class,
			})
			if err != nil {
				return fmt.Errorf("PutObject(%s) failed: %w", class, err)
			}
			head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			})
			if err != nil {
				return fmt.Errorf("HeadObject(%s) failed: %w", class, err)
			}
			if head.StorageClass != class {
				return fmt.Errorf("StorageClass round trip: put %s, head %s", class, head.StorageClass)
			}
			if head.ContentLength == nil || *head.ContentLength != int64(len("storage class round trip")) {
				return fmt.Errorf("ContentLength round trip for %s: got %v", class, head.ContentLength)
			}
		}

		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String("sc-roundtrip/outposts.txt"),
			Body:         strings.NewReader("hardware-bound class"),
			StorageClass: types.StorageClassOutposts,
		})
		if err == nil {
			return fmt.Errorf("PutObject with OUTPOSTS must be rejected on this platform, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error for OUTPOSTS, got %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "InvalidArgument" {
			return fmt.Errorf("expected InvalidArgument for OUTPOSTS, got %s: %v", apiErr.ErrorCode(), err)
		}
		if code := awsHTTPStatus(err); code != http.StatusBadRequest {
			return fmt.Errorf("expected HTTP 400 for OUTPOSTS, got %d: %v", code, err)
		}
		return nil
	}))

	return results
}
