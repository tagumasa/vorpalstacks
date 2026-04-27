package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *TestRunner) s3ObjectConfigTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	lockBucket := s3Bucket(ts, "lock")
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket:                     aws.String(lockBucket),
		ObjectLockEnabledForBucket: aws.Bool(true),
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "s3",
			TestName: "SetupLockBucket",
			Status:   "FAIL",
			Error:    fmt.Sprintf("CreateBucket (lock) failed: %v", err),
		})
	}
	defer s3CleanupBucket(ctx, client, lockBucket)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("tagged-obj.txt"),
		Body:   strings.NewReader("tag me"),
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "s3",
			TestName: "SetupTaggedObject",
			Status:   "FAIL",
			Error:    fmt.Sprintf("PutObject tagged-obj.txt failed: %v", err),
		})
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(lockBucket),
		Key:    aws.String("legal-hold-obj.txt"),
		Body:   strings.NewReader("legal hold content"),
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "s3",
			TestName: "SetupLegalHoldObject",
			Status:   "FAIL",
			Error:    fmt.Sprintf("PutObject legal-hold-obj.txt failed: %v", err),
		})
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(lockBucket),
		Key:    aws.String("retention-obj.txt"),
		Body:   strings.NewReader("retention content"),
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "s3",
			TestName: "SetupRetentionObject",
			Status:   "FAIL",
			Error:    fmt.Sprintf("PutObject retention-obj.txt failed: %v", err),
		})
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("attrs-obj.txt"),
		Body:   strings.NewReader("attributes content"),
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "s3",
			TestName: "SetupAttrsObject",
			Status:   "FAIL",
			Error:    fmt.Sprintf("PutObject attrs-obj.txt failed: %v", err),
		})
	}

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
		if getResp.Owner == nil {
			return fmt.Errorf("Owner is nil")
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
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		if resp.StorageClass != types.StorageClassStandard {
			return fmt.Errorf("expected StorageClass STANDARD, got %s", resp.StorageClass)
		}
		return nil
	}))

	return results
}
