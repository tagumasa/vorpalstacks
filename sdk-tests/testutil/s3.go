package testutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"vorpalstacks-sdk-tests/config"
)

func s3Bucket(ts string, name string) string {
	return fmt.Sprintf("s3test-%s-%s", name, ts)
}

func s3CleanupBucket(ctx context.Context, client *s3.Client, bucket string) {
	listResp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		return
	}
	if len(listResp.Contents) > 0 {
		var objs []types.ObjectIdentifier
		for _, o := range listResp.Contents {
			objs = append(objs, types.ObjectIdentifier{Key: o.Key})
		}
		client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &types.Delete{Objects: objs},
		})
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

	// ========== BUCKET CORE ==========

	results = append(results, r.RunTest("s3", "CreateBucket", func() error {
		resp, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Location == nil {
			return fmt.Errorf("Location is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListBuckets", func() error {
		resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return err
		}
		if resp.Buckets == nil {
			return fmt.Errorf("buckets list is nil")
		}
		return nil
	}))

	sortBuckets := []string{
		s3Bucket(ts, "sort-z"),
		s3Bucket(ts, "sort-a"),
		s3Bucket(ts, "sort-m"),
	}
	results = append(results, r.RunTest("s3", "ListBuckets_SortedByName", func() error {
		for _, b := range sortBuckets {
			_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(b)})
			if err != nil {
				return fmt.Errorf("failed to create bucket %s: %v", b, err)
			}
		}
		defer func() {
			for _, b := range sortBuckets {
				s3CleanupBucket(ctx, client, b)
			}
		}()
		resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return err
		}
		for i := 1; i < len(resp.Buckets); i++ {
			if *resp.Buckets[i].Name < *resp.Buckets[i-1].Name {
				return fmt.Errorf("buckets not sorted: %s before %s", *resp.Buckets[i-1].Name, *resp.Buckets[i].Name)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadBucket", func() error {
		resp, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.BucketRegion == nil {
			return fmt.Errorf("BucketRegion is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetBucketLocation", func() error {
		resp, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.LocationConstraint != "" && resp.LocationConstraint != "us-east-1" {
			return fmt.Errorf("unexpected LocationConstraint: %s", resp.LocationConstraint)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CreateBucket_DuplicateName", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate bucket name")
		}
		var bae *types.BucketAlreadyOwnedByYou
		var bae2 *types.BucketAlreadyExists
		if !errors.As(err, &bae) && !errors.As(err, &bae2) {
			return fmt.Errorf("expected BucketAlreadyOwnedByYou or BucketAlreadyExists, got: %T: %v", err, err)
		}
		return nil
	}))

	// ========== OBJECT DATA ==========

	results = append(results, r.RunTest("s3", "PutObject", func() error {
		content := "Hello, World!"
		resp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
			Body:   strings.NewReader(content),
		})
		if err != nil {
			return err
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject", func() error {
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err != nil {
			return err
		}
		if resp.Body == nil {
			return fmt.Errorf("body is nil")
		}
		defer resp.Body.Close()
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject", func() error {
		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		if err != nil {
			return err
		}
		if resp.ContentLength == nil {
			return fmt.Errorf("ContentLength is nil")
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Contents == nil {
			return fmt.Errorf("Contents is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "ListObjectsAfterDelete", func() error {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if len(resp.Contents) > 0 {
			var keys []string
			for _, obj := range resp.Contents {
				keys = append(keys, *obj.Key)
			}
			return fmt.Errorf("bucket %s not empty after delete, objects: %v", bucketName, keys)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CopyObject", func() error {
		srcKey := "copy-source.txt"
		dstKey := "copy-dest.txt"
		content := "copy me"

		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(srcKey),
			Body:   strings.NewReader(content),
		})
		if err != nil {
			return fmt.Errorf("put source: %v", err)
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(dstKey),
			CopySource: aws.String(bucketName + "/" + srcKey),
		})
		if err != nil {
			return fmt.Errorf("copy: %v", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(dstKey),
		})
		if err != nil {
			return fmt.Errorf("get copy: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read copy: %v", err)
		}
		if string(body) != content {
			return fmt.Errorf("copy content mismatch: got %q, want %q", string(body), content)
		}
		return nil
	}))

	// ========== OBJECT DATA EDGE CASES ==========

	results = append(results, r.RunTest("s3", "PutObject_GetObject_ContentVerification", func() error {
		content := "Hello, S3 content verification! Japanese test"
		key := "verify-content.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String(key),
			Body:          strings.NewReader(content),
			ContentType:   aws.String("text/plain; charset=utf-8"),
			ContentLength: aws.Int64(int64(len(content))),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body: %v", err)
		}
		if string(body) != content {
			return fmt.Errorf("content mismatch: got %q, want %q", string(body), content)
		}

		if resp.ContentLength == nil || *resp.ContentLength != int64(len(content)) {
			return fmt.Errorf("ContentLength mismatch: got %v, want %d", resp.ContentLength, len(content))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_Overwrite", func() error {
		key := "overwrite-test.txt"
		content1 := "Original content"
		content2 := "Updated content"

		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader(content1),
		})
		if err != nil {
			return fmt.Errorf("put 1: %v", err)
		}

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader(content2),
		})
		if err != nil {
			return fmt.Errorf("put 2: %v", err)
		}

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read: %v", err)
		}
		if string(body) != content2 {
			return fmt.Errorf("after overwrite expected %q, got %q", content2, string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadObject_VerifyMetadata", func() error {
		key := "metadata-test.txt"
		content := "metadata check"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String(key),
			Body:          strings.NewReader(content),
			ContentType:   aws.String("application/json"),
			ContentLength: aws.Int64(int64(len(content))),
			Metadata: map[string]string{
				"custom-key": "custom-value",
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("head: %v", err)
		}
		if resp.ContentType == nil || *resp.ContentType != "application/json" {
			return fmt.Errorf("ContentType mismatch, got %v", resp.ContentType)
		}
		if resp.ContentLength == nil || *resp.ContentLength != int64(len(content)) {
			return fmt.Errorf("ContentLength mismatch, got %v", resp.ContentLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListObjectsV2_MultipleObjects", func() error {
		lobBucket := s3Bucket(ts, "list")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(lobBucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer s3CleanupBucket(ctx, client, lobBucket)

		for i := 0; i < 5; i++ {
			key := fmt.Sprintf("obj%d.txt", i)
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(lobBucket),
				Key:    aws.String(key),
				Body:   strings.NewReader(fmt.Sprintf("content %d", i)),
			})
			if err != nil {
				return fmt.Errorf("put %d: %v", i, err)
			}
		}

		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(lobBucket),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(resp.Contents) != 5 {
			return fmt.Errorf("expected 5 contents, got %d", len(resp.Contents))
		}
		for _, obj := range resp.Contents {
			if obj.Key == nil {
				return fmt.Errorf("key is nil")
			}
			if obj.Size == nil || *obj.Size == 0 {
				return fmt.Errorf("expected non-zero size for %s", *obj.Key)
			}
		}
		return nil
	}))

	// ========== BUCKET TAGGING ==========

	results = append(results, r.RunTest("s3", "PutBucketTagging", func() error {
		_, err := client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
			Bucket: aws.String(bucketName),
			Tagging: &types.Tagging{
				TagSet: []types.Tag{
					{Key: aws.String("Environment"), Value: aws.String("Test")},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketTagging", func() error {
		resp, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.TagSet == nil {
			return fmt.Errorf("TagSet is nil")
		}
		found := false
		for _, t := range resp.TagSet {
			if t.Key != nil && *t.Key == "Environment" && t.Value != nil && *t.Value == "Test" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected tag Environment=Test not found in TagSet")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketTagging", func() error {
		_, err := client.DeleteBucketTagging(ctx, &s3.DeleteBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	// ========== BUCKET ACL ==========

	results = append(results, r.RunTest("s3", "PutBucketAcl", func() error {
		_, err := client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
			Bucket: aws.String(bucketName),
			ACL:    types.BucketCannedACLPrivate,
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketAcl", func() error {
		resp, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Owner == nil {
			return fmt.Errorf("Owner is nil")
		}
		return nil
	}))

	// ========== BUCKET POLICY ==========

	results = append(results, r.RunTest("s3", "PutBucketPolicy", func() error {
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"arn:aws:s3:::` + bucketName + `/*"}]}`
		_, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
			Bucket: aws.String(bucketName),
			Policy: aws.String(policy),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketPolicy", func() error {
		resp, err := client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil {
			return fmt.Errorf("Policy is nil")
		}
		if !strings.Contains(*resp.Policy, "Allow") {
			return fmt.Errorf("policy missing expected content")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketPolicy", func() error {
		_, err := client.DeleteBucketPolicy(ctx, &s3.DeleteBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketPolicy_NotFound", func() error {
		_, err := client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing policy")
		}
		if !strings.Contains(err.Error(), "NoSuchBucketPolicy") {
			return fmt.Errorf("expected NoSuchBucketPolicy, got: %v", err)
		}
		return nil
	}))

	// ========== BUCKET CORS ==========

	results = append(results, r.RunTest("s3", "PutBucketCors", func() error {
		_, err := client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
			Bucket: aws.String(bucketName),
			CORSConfiguration: &types.CORSConfiguration{
				CORSRules: []types.CORSRule{
					{
						AllowedMethods: []string{"GET", "PUT"},
						AllowedOrigins: []string{"https://example.com"},
						MaxAgeSeconds:  aws.Int32(3600),
					},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketCors", func() error {
		resp, err := client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if len(resp.CORSRules) == 0 {
			return fmt.Errorf("expected at least one CORS rule")
		}
		if resp.CORSRules[0].MaxAgeSeconds == nil || *resp.CORSRules[0].MaxAgeSeconds != 3600 {
			return fmt.Errorf("MaxAgeSeconds mismatch, got %v", resp.CORSRules[0].MaxAgeSeconds)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketCors", func() error {
		_, err := client.DeleteBucketCors(ctx, &s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketCors_NotFound", func() error {
		_, err := client.GetBucketCors(ctx, &s3.GetBucketCorsInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing CORS")
		}
		if !strings.Contains(err.Error(), "NoSuchCORSConfiguration") {
			return fmt.Errorf("expected NoSuchCORSConfiguration, got: %v", err)
		}
		return nil
	}))

	// ========== BUCKET ENCRYPTION ==========

	results = append(results, r.RunTest("s3", "PutBucketEncryption", func() error {
		_, err := client.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
			Bucket: aws.String(bucketName),
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{
					{
						ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
							SSEAlgorithm: types.ServerSideEncryptionAes256,
						},
					},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketEncryption", func() error {
		resp, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.ServerSideEncryptionConfiguration == nil {
			return fmt.Errorf("ServerSideEncryptionConfiguration is nil")
		}
		if len(resp.ServerSideEncryptionConfiguration.Rules) == 0 {
			return fmt.Errorf("expected at least one encryption rule")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketEncryption", func() error {
		_, err := client.DeleteBucketEncryption(ctx, &s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketEncryption_NotFound", func() error {
		_, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing encryption")
		}
		if !strings.Contains(err.Error(), "ServerSideEncryptionConfigurationNotFoundError") {
			return fmt.Errorf("expected ServerSideEncryptionConfigurationNotFoundError, got: %v", err)
		}
		return nil
	}))

	// ========== BUCKET VERSIONING ==========

	results = append(results, r.RunTest("s3", "PutBucketVersioning_Enabled", func() error {
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketVersioning_Enabled", func() error {
		resp, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.BucketVersioningStatusEnabled {
			return fmt.Errorf("expected Enabled, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutBucketVersioning_Suspended", func() error {
		_, err := client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusSuspended,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketVersioning_Suspended", func() error {
		resp, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.BucketVersioningStatusSuspended {
			return fmt.Errorf("expected Suspended, got %s", resp.Status)
		}
		return nil
	}))

	// ========== BUCKET LIFECYCLE ==========

	results = append(results, r.RunTest("s3", "PutBucketLifecycleConfiguration", func() error {
		ruleID := "test-expire-rule"
		_, err := client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: []types.LifecycleRule{
					{
						ID:     aws.String(ruleID),
						Status: types.ExpirationStatusEnabled,
						Filter: &types.LifecycleRuleFilter{
							Prefix: aws.String("logs/"),
						},
						Expiration: &types.LifecycleExpiration{
							Days: aws.Int32(30),
						},
					},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketLifecycleConfiguration", func() error {
		resp, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if len(resp.Rules) == 0 {
			return fmt.Errorf("expected at least one lifecycle rule")
		}
		if resp.Rules[0].ID == nil || *resp.Rules[0].ID != "test-expire-rule" {
			return fmt.Errorf("rule ID mismatch, got %v", resp.Rules[0].ID)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketLifecycleConfiguration", func() error {
		_, err := client.DeleteBucketLifecycle(ctx, &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketLifecycleConfiguration_NotFound", func() error {
		_, err := client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing lifecycle")
		}
		if !strings.Contains(err.Error(), "NoSuchLifecycleConfiguration") {
			return fmt.Errorf("expected NoSuchLifecycleConfiguration, got: %v", err)
		}
		return nil
	}))

	// ========== BUCKET WEBSITE ==========

	results = append(results, r.RunTest("s3", "PutBucketWebsite", func() error {
		_, err := client.PutBucketWebsite(ctx, &s3.PutBucketWebsiteInput{
			Bucket: aws.String(bucketName),
			WebsiteConfiguration: &types.WebsiteConfiguration{
				IndexDocument: &types.IndexDocument{
					Suffix: aws.String("index.html"),
				},
				ErrorDocument: &types.ErrorDocument{
					Key: aws.String("error.html"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketWebsite", func() error {
		resp, err := client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.IndexDocument == nil || resp.IndexDocument.Suffix == nil || *resp.IndexDocument.Suffix != "index.html" {
			return fmt.Errorf("IndexDocument Suffix mismatch, got %v", resp.IndexDocument)
		}
		if resp.ErrorDocument == nil || resp.ErrorDocument.Key == nil || *resp.ErrorDocument.Key != "error.html" {
			return fmt.Errorf("ErrorDocument Key mismatch, got %v", resp.ErrorDocument)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketWebsite", func() error {
		_, err := client.DeleteBucketWebsite(ctx, &s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketWebsite_NotFound", func() error {
		_, err := client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing website")
		}
		if !strings.Contains(err.Error(), "NoSuchWebsiteConfiguration") {
			return fmt.Errorf("expected NoSuchWebsiteConfiguration, got: %v", err)
		}
		return nil
	}))

	// ========== OBJECT LOCK CONFIGURATION ==========

	lockBucket := s3Bucket(ts, "lock")
	results = append(results, r.RunTest("s3", "PutObjectLockConfiguration", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket:                     aws.String(lockBucket),
			ObjectLockEnabledForBucket: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}

		_, err = client.PutObjectLockConfiguration(ctx, &s3.PutObjectLockConfigurationInput{
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
		return err
	}))

	results = append(results, r.RunTest("s3", "GetObjectLockConfiguration", func() error {
		resp, err := client.GetObjectLockConfiguration(ctx, &s3.GetObjectLockConfigurationInput{
			Bucket: aws.String(lockBucket),
		})
		if err != nil {
			return err
		}
		if resp.ObjectLockConfiguration == nil {
			return fmt.Errorf("ObjectLockConfiguration is nil")
		}
		if resp.ObjectLockConfiguration.ObjectLockEnabled != types.ObjectLockEnabledEnabled {
			return fmt.Errorf("expected Enabled, got %s", resp.ObjectLockConfiguration.ObjectLockEnabled)
		}
		return nil
	}))

	// ========== BUCKET NOTIFICATION ==========

	results = append(results, r.RunTest("s3", "PutBucketNotificationConfiguration", func() error {
		_, err := client.PutBucketNotificationConfiguration(ctx, &s3.PutBucketNotificationConfigurationInput{
			Bucket: aws.String(bucketName),
			NotificationConfiguration: &types.NotificationConfiguration{
				TopicConfigurations: []types.TopicConfiguration{
					{
						TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
						Events:   []types.Event{types.EventS3ObjectCreatedPut},
					},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketNotificationConfiguration", func() error {
		resp, err := client.GetBucketNotificationConfiguration(ctx, &s3.GetBucketNotificationConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if len(resp.TopicConfigurations) == 0 {
			return fmt.Errorf("expected at least one topic configuration")
		}
		return nil
	}))

	// ========== BUCKET LOGGING ==========

	results = append(results, r.RunTest("s3", "PutBucketLogging", func() error {
		_, err := client.PutBucketLogging(ctx, &s3.PutBucketLoggingInput{
			Bucket: aws.String(bucketName),
			BucketLoggingStatus: &types.BucketLoggingStatus{
				LoggingEnabled: &types.LoggingEnabled{
					TargetBucket: aws.String(bucketName),
					TargetPrefix: aws.String("logs/"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketLogging", func() error {
		resp, err := client.GetBucketLogging(ctx, &s3.GetBucketLoggingInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.LoggingEnabled != nil {
			if resp.LoggingEnabled.TargetPrefix == nil || *resp.LoggingEnabled.TargetPrefix != "logs/" {
				return fmt.Errorf("TargetPrefix mismatch, got %v", resp.LoggingEnabled.TargetPrefix)
			}
		}
		return nil
	}))

	// ========== PUBLIC ACCESS BLOCK ==========

	results = append(results, r.RunTest("s3", "PutPublicAccessBlock", func() error {
		_, err := client.PutPublicAccessBlock(ctx, &s3.PutPublicAccessBlockInput{
			Bucket: aws.String(bucketName),
			PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
				BlockPublicAcls:       aws.Bool(true),
				IgnorePublicAcls:      aws.Bool(true),
				BlockPublicPolicy:     aws.Bool(true),
				RestrictPublicBuckets: aws.Bool(true),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetPublicAccessBlock", func() error {
		resp, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		cfg := resp.PublicAccessBlockConfiguration
		if cfg == nil {
			return fmt.Errorf("PublicAccessBlockConfiguration is nil")
		}
		if cfg.BlockPublicAcls == nil || !*cfg.BlockPublicAcls {
			return fmt.Errorf("BlockPublicAcls should be true")
		}
		if cfg.RestrictPublicBuckets == nil || !*cfg.RestrictPublicBuckets {
			return fmt.Errorf("RestrictPublicBuckets should be true")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeletePublicAccessBlock", func() error {
		_, err := client.DeletePublicAccessBlock(ctx, &s3.DeletePublicAccessBlockInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetPublicAccessBlock_NotFound", func() error {
		_, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing public access block")
		}
		if !strings.Contains(err.Error(), "NoSuchPublicAccessBlockConfiguration") {
			return fmt.Errorf("expected NoSuchPublicAccessBlockConfiguration, got: %v", err)
		}
		return nil
	}))

	// ========== BUCKET OWNERSHIP CONTROLS ==========

	results = append(results, r.RunTest("s3", "PutBucketOwnershipControls", func() error {
		_, err := client.PutBucketOwnershipControls(ctx, &s3.PutBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
			OwnershipControls: &types.OwnershipControls{
				Rules: []types.OwnershipControlsRule{
					{ObjectOwnership: types.ObjectOwnershipBucketOwnerPreferred},
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketOwnershipControls", func() error {
		resp, err := client.GetBucketOwnershipControls(ctx, &s3.GetBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.OwnershipControls == nil || len(resp.OwnershipControls.Rules) == 0 {
			return fmt.Errorf("expected at least one ownership controls rule")
		}
		if resp.OwnershipControls.Rules[0].ObjectOwnership != types.ObjectOwnershipBucketOwnerPreferred {
			return fmt.Errorf("ObjectOwnership mismatch, got %s", resp.OwnershipControls.Rules[0].ObjectOwnership)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketOwnershipControls", func() error {
		_, err := client.DeleteBucketOwnershipControls(ctx, &s3.DeleteBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketOwnershipControls_NotFound", func() error {
		_, err := client.GetBucketOwnershipControls(ctx, &s3.GetBucketOwnershipControlsInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing ownership controls")
		}
		if !strings.Contains(err.Error(), "OwnershipControlsNotFoundError") {
			return fmt.Errorf("expected OwnershipControlsNotFoundError, got: %v", err)
		}
		return nil
	}))

	// ========== BUCKET REQUEST PAYMENT ==========

	results = append(results, r.RunTest("s3", "PutBucketRequestPayment", func() error {
		_, err := client.PutBucketRequestPayment(ctx, &s3.PutBucketRequestPaymentInput{
			Bucket: aws.String(bucketName),
			RequestPaymentConfiguration: &types.RequestPaymentConfiguration{
				Payer: types.PayerRequester,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketRequestPayment", func() error {
		resp, err := client.GetBucketRequestPayment(ctx, &s3.GetBucketRequestPaymentInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Payer != types.PayerRequester {
			return fmt.Errorf("expected Requester, got %s", resp.Payer)
		}
		return nil
	}))

	// ========== BUCKET ACCELERATE CONFIGURATION ==========

	results = append(results, r.RunTest("s3", "PutBucketAccelerateConfiguration", func() error {
		_, err := client.PutBucketAccelerateConfiguration(ctx, &s3.PutBucketAccelerateConfigurationInput{
			Bucket: aws.String(bucketName),
			AccelerateConfiguration: &types.AccelerateConfiguration{
				Status: types.BucketAccelerateStatusSuspended,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetBucketAccelerateConfiguration", func() error {
		resp, err := client.GetBucketAccelerateConfiguration(ctx, &s3.GetBucketAccelerateConfigurationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.BucketAccelerateStatusSuspended {
			return fmt.Errorf("expected Suspended, got %s", resp.Status)
		}
		return nil
	}))

	// ========== OBJECT TAGGING ==========

	_, _ = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("tagged-obj.txt"),
		Body:   strings.NewReader("tag me"),
	})

	results = append(results, r.RunTest("s3", "PutObjectTagging", func() error {
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
		return err
	}))

	results = append(results, r.RunTest("s3", "GetObjectTagging", func() error {
		resp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return err
		}
		if len(resp.TagSet) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(resp.TagSet))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObjectTagging", func() error {
		_, err := client.DeleteObjectTagging(ctx, &s3.DeleteObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetObjectTagging_Empty", func() error {
		resp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return err
		}
		if len(resp.TagSet) != 0 {
			return fmt.Errorf("expected 0 tags after delete, got %d", len(resp.TagSet))
		}
		return nil
	}))

	// ========== OBJECT ACL ==========

	results = append(results, r.RunTest("s3", "GetObjectAcl", func() error {
		resp, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
		})
		if err != nil {
			return err
		}
		if resp.Owner == nil {
			return fmt.Errorf("Owner is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObjectAcl", func() error {
		_, err := client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("tagged-obj.txt"),
			ACL:    types.ObjectCannedACLPrivate,
		})
		return err
	}))

	// ========== OBJECT LEGAL HOLD ==========

	_, _ = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(lockBucket),
		Key:    aws.String("legal-hold-obj.txt"),
		Body:   strings.NewReader("legal hold test"),
	})

	results = append(results, r.RunTest("s3", "PutObjectLegalHold", func() error {
		_, err := client.PutObjectLegalHold(ctx, &s3.PutObjectLegalHoldInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("legal-hold-obj.txt"),
			LegalHold: &types.ObjectLockLegalHold{
				Status: types.ObjectLockLegalHoldStatusOn,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetObjectLegalHold", func() error {
		resp, err := client.GetObjectLegalHold(ctx, &s3.GetObjectLegalHoldInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("legal-hold-obj.txt"),
		})
		if err != nil {
			return err
		}
		if resp.LegalHold == nil {
			return fmt.Errorf("LegalHold is nil")
		}
		if resp.LegalHold.Status != types.ObjectLockLegalHoldStatusOn {
			return fmt.Errorf("expected ON, got %s", resp.LegalHold.Status)
		}
		return nil
	}))

	// ========== OBJECT RETENTION ==========

	_, _ = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(lockBucket),
		Key:    aws.String("retention-obj.txt"),
		Body:   strings.NewReader("retention test"),
	})

	retainUntil := time.Now().Add(24 * time.Hour)
	results = append(results, r.RunTest("s3", "PutObjectRetention", func() error {
		_, err := client.PutObjectRetention(ctx, &s3.PutObjectRetentionInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("retention-obj.txt"),
			Retention: &types.ObjectLockRetention{
				Mode:            types.ObjectLockRetentionModeGovernance,
				RetainUntilDate: &retainUntil,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "GetObjectRetention", func() error {
		resp, err := client.GetObjectRetention(ctx, &s3.GetObjectRetentionInput{
			Bucket: aws.String(lockBucket),
			Key:    aws.String("retention-obj.txt"),
		})
		if err != nil {
			return err
		}
		if resp.Retention == nil {
			return fmt.Errorf("Retention is nil")
		}
		if resp.Retention.Mode != types.ObjectLockRetentionModeGovernance {
			return fmt.Errorf("expected GOVERNANCE, got %s", resp.Retention.Mode)
		}
		return nil
	}))

	// ========== GET OBJECT ATTRIBUTES ==========

	_, _ = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("attrs-obj.txt"),
		Body:   strings.NewReader("object attributes test content"),
	})

	results = append(results, r.RunTest("s3", "GetObjectAttributes", func() error {
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
			return err
		}
		if resp.ObjectSize == nil || *resp.ObjectSize == 0 {
			return fmt.Errorf("ObjectSize is nil or zero")
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		return nil
	}))

	// ========== MULTIPART UPLOAD ==========

	mpuBucket := s3Bucket(ts, "mpu")
	results = append(results, r.RunTest("s3", "CreateMultipartUpload", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(mpuBucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	}))

	var uploadID string
	results = append(results, r.RunTest("s3", "CreateMultipartUpload_Initiate", func() error {
		resp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(mpuBucket),
			Key:    aws.String("multipart-obj.txt"),
		})
		if err != nil {
			return err
		}
		if resp.UploadId == nil {
			return fmt.Errorf("UploadId is nil")
		}
		uploadID = *resp.UploadId
		return nil
	}))

	var part1ETag string
	results = append(results, r.RunTest("s3", "UploadPart", func() error {
		resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			PartNumber: aws.Int32(1),
			UploadId:   aws.String(uploadID),
			Body:       strings.NewReader("part one content"),
		})
		if err != nil {
			return err
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		part1ETag = *resp.ETag
		return nil
	}))

	var part2ETag string
	results = append(results, r.RunTest("s3", "UploadPart_Part2", func() error {
		resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(mpuBucket),
			Key:        aws.String("multipart-obj.txt"),
			PartNumber: aws.Int32(2),
			UploadId:   aws.String(uploadID),
			Body:       strings.NewReader("part two content"),
		})
		if err != nil {
			return err
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		part2ETag = *resp.ETag
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListParts", func() error {
		resp, err := client.ListParts(ctx, &s3.ListPartsInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("multipart-obj.txt"),
			UploadId: aws.String(uploadID),
		})
		if err != nil {
			return err
		}
		if len(resp.Parts) != 2 {
			return fmt.Errorf("expected 2 parts, got %d", len(resp.Parts))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "CompleteMultipartUpload", func() error {
		resp, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(mpuBucket),
			Key:      aws.String("multipart-obj.txt"),
			UploadId: aws.String(uploadID),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{PartNumber: aws.Int32(1), ETag: aws.String(part1ETag)},
					{PartNumber: aws.Int32(2), ETag: aws.String(part2ETag)},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Location == nil {
			return fmt.Errorf("Location is nil")
		}
		if resp.ETag == nil {
			return fmt.Errorf("ETag is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultipartUpload_GetObject", func() error {
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(mpuBucket),
			Key:    aws.String("multipart-obj.txt"),
		})
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		expected := "part one contentpart two content"
		if string(body) != expected {
			return fmt.Errorf("content mismatch: got %q, want %q", string(body), expected)
		}
		return nil
	}))

	// ========== ABORT MULTIPART UPLOAD ==========

	abortBucket := s3Bucket(ts, "abort")
	results = append(results, r.RunTest("s3", "AbortMultipartUpload", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(abortBucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer s3CleanupBucket(ctx, client, abortBucket)

		createResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(abortBucket),
			Key:    aws.String("abort-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("create mpu: %v", err)
		}

		_, err = client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(abortBucket),
			Key:        aws.String("abort-obj.txt"),
			PartNumber: aws.Int32(1),
			UploadId:   createResp.UploadId,
			Body:       strings.NewReader("abort data"),
		})
		if err != nil {
			return fmt.Errorf("upload part: %v", err)
		}

		_, err = client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(abortBucket),
			Key:      aws.String("abort-obj.txt"),
			UploadId: createResp.UploadId,
		})
		return err
	}))

	// ========== LIST MULTIPART UPLOADS ==========

	lmpuBucket := s3Bucket(ts, "listmpu")
	results = append(results, r.RunTest("s3", "ListMultipartUploads", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(lmpuBucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer s3CleanupBucket(ctx, client, lmpuBucket)

		_, err = client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(lmpuBucket),
			Key:    aws.String("listed-mpu-obj.txt"),
		})
		if err != nil {
			return fmt.Errorf("create mpu: %v", err)
		}

		resp, err := client.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
			Bucket: aws.String(lmpuBucket),
		})
		if err != nil {
			return err
		}
		if len(resp.Uploads) == 0 {
			return fmt.Errorf("expected at least one upload")
		}
		return nil
	}))

	// ========== LIST OBJECT VERSIONS ==========

	verBucket := s3Bucket(ts, "versions")
	results = append(results, r.RunTest("s3", "ListObjectVersions", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(verBucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer s3CleanupBucket(ctx, client, verBucket)

		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(verBucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			return fmt.Errorf("enable versioning: %v", err)
		}

		for i := 0; i < 3; i++ {
			_, err = client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(verBucket),
				Key:    aws.String("versioned-obj.txt"),
				Body:   strings.NewReader(fmt.Sprintf("version %d", i)),
			})
			if err != nil {
				return fmt.Errorf("put version %d: %v", i, err)
			}
		}

		resp, err := client.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{
			Bucket: aws.String(verBucket),
		})
		if err != nil {
			return err
		}
		if len(resp.Versions) == 0 {
			return fmt.Errorf("expected at least one version")
		}
		return nil
	}))

	// ========== SELECT OBJECT CONTENT ==========

	_, _ = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("select-data.csv"),
		Body:   strings.NewReader("name,age\nAlice,30\nBob,25\n"),
	})

	results = append(results, r.RunTest("s3", "SelectObjectContent", func() error {
		resp, err := client.SelectObjectContent(ctx, &s3.SelectObjectContentInput{
			Bucket:         aws.String(bucketName),
			Key:            aws.String("select-data.csv"),
			Expression:     aws.String("SELECT s.name FROM s3object s WHERE s.age = '30'"),
			ExpressionType: types.ExpressionTypeSql,
			InputSerialization: &types.InputSerialization{
				CSV: &types.CSVInput{
					FileHeaderInfo: types.FileHeaderInfoUse,
				},
			},
			OutputSerialization: &types.OutputSerialization{
				CSV: &types.CSVOutput{},
			},
		})
		if err != nil {
			return err
		}
		if resp.GetStream() == nil {
			return fmt.Errorf("event stream is nil")
		}
		var results []string
		for event := range resp.GetStream().Events() {
			switch v := event.(type) {
			case *types.SelectObjectContentEventStreamMemberRecords:
				results = append(results, string(v.Value.Payload))
			case *types.SelectObjectContentEventStreamMemberStats:
			case *types.SelectObjectContentEventStreamMemberProgress:
			case *types.SelectObjectContentEventStreamMemberCont:
			case *types.SelectObjectContentEventStreamMemberEnd:
			}
		}
		if len(results) == 0 {
			return fmt.Errorf("expected at least one record result")
		}
		return nil
	}))

	// ========== ERROR / EDGE CASE TESTS ==========

	results = append(results, r.RunTest("s3", "HeadBucket_NonExistent", func() error {
		_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String("nonexistent-bucket-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent bucket")
		}
		var nf *types.NotFound
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFound, got: %T: %v", err, err)
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
		var nsk *types.NoSuchKey
		if !errors.As(err, &nsk) {
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
		var nf *types.NotFound
		if !errors.As(err, &nf) {
			return fmt.Errorf("expected NotFound (404), got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteObject_NonExistentKey", func() error {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("nonexistent-delete-key.txt"),
		})
		if err != nil {
			return fmt.Errorf("delete non-existent key should not error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucket_NotEmpty", func() error {
		neBucket := s3Bucket(ts, "notempty")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(neBucket)})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer s3CleanupBucket(ctx, client, neBucket)

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(neBucket),
			Key:    aws.String("keep-me.txt"),
			Body:   strings.NewReader("data"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(neBucket),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleting non-empty bucket")
		}
		if !strings.Contains(err.Error(), "BucketNotEmpty") {
			return fmt.Errorf("expected BucketNotEmpty, got: %v", err)
		}
		return nil
	}))

	// ========== MULTIBYTE SUPPORT ==========

	results = append(results, r.RunTest("s3", "MultiByte_PutGetRoundtrip", func() error {
		mbCases := []struct {
			key  string
			body string
		}{
			{"テスト/日本語ファイル.txt", "こんにちは世界。これは日本語のテストデータです。"},
			{"文档/简体中文.txt", "你好世界。这是简体中文的测试数据。"},
			{"文件/繁體中文.txt", "你好世界。這是繁體中文的測試資料。"},
		}
		for _, tc := range mbCases {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(tc.key),
				Body:        strings.NewReader(tc.body),
				ContentType: aws.String("text/plain; charset=utf-8"),
			})
			if err != nil {
				return fmt.Errorf("put multibyte object %q: %v", tc.key, err)
			}
			resp, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(tc.key),
			})
			if err != nil {
				return fmt.Errorf("get multibyte object %q: %v", tc.key, err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read body for %q: %v", tc.key, err)
			}
			resp.Body.Close()
			if string(body) != tc.body {
				return fmt.Errorf("body mismatch for %q: expected %q, got %q", tc.key, tc.body, string(body))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_MetadataRoundtrip", func() error {
		mbCases := []struct {
			key  string
			body string
			meta map[string]string
		}{
			{
				key:  "テスト/metadata-日本語.txt",
				body: "メタデータテスト",
				meta: map[string]string{
					"author":  "田中太郎",
					"desc":    "日本語メタデータのテスト",
					"project": "多バイト文字プロジェクト",
				},
			},
			{
				key:  "文档/metadata-中文.txt",
				body: "元数据测试",
				meta: map[string]string{
					"author":  "张三",
					"desc":    "简体中文元数据测试",
					"project": "多字节字符项目",
				},
			},
			{
				key:  "文件/metadata-繁體.txt",
				body: "元資料測試",
				meta: map[string]string{
					"author":  "李四",
					"desc":    "繁體中文元資料測試",
					"project": "多位元組字元專案",
				},
			},
		}
		for _, tc := range mbCases {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(tc.key),
				Body:        strings.NewReader(tc.body),
				ContentType: aws.String("text/plain; charset=utf-8"),
				Metadata:    tc.meta,
			})
			if err != nil {
				return fmt.Errorf("put multibyte metadata object %q: %v", tc.key, err)
			}
			resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(tc.key),
			})
			if err != nil {
				return fmt.Errorf("head multibyte metadata object %q: %v", tc.key, err)
			}
			for k, v := range tc.meta {
				actual, ok := resp.Metadata[k]
				if !ok {
					return fmt.Errorf("missing metadata key %q for object %q", k, tc.key)
				}
				if actual != v {
					return fmt.Errorf("metadata %q mismatch for %q: expected %q, got %q", k, tc.key, v, actual)
				}
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_ListObjectsV2", func() error {
		mbKeys := []string{
			"テスト/リスト日本語.txt",
			"文档/列表简体.txt",
			"文件/列表繁體.txt",
			"mixed/混合日本語中文繁體.txt",
		}
		for _, key := range mbKeys {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   strings.NewReader("list-test-content"),
			})
			if err != nil {
				return fmt.Errorf("put list test object %q: %v", key, err)
			}
		}
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String("テスト/"),
		})
		if err != nil {
			return fmt.Errorf("list objects v2 with Japanese prefix: %v", err)
		}
		found := false
		for _, obj := range resp.Contents {
			if *obj.Key == "テスト/リスト日本語.txt" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("Japanese key not found in ListObjectsV2 results")
		}

		resp2, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String("文档/"),
		})
		if err != nil {
			return fmt.Errorf("list objects v2 with Chinese prefix: %v", err)
		}
		found2 := false
		for _, obj := range resp2.Contents {
			if *obj.Key == "文档/列表简体.txt" {
				found2 = true
				break
			}
		}
		if !found2 {
			return fmt.Errorf("Simplified Chinese key not found in ListObjectsV2 results")
		}

		resp3, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String("文件/"),
		})
		if err != nil {
			return fmt.Errorf("list objects v2 with Taiwanese prefix: %v", err)
		}
		found3 := false
		for _, obj := range resp3.Contents {
			if *obj.Key == "文件/列表繁體.txt" {
				found3 = true
				break
			}
		}
		if !found3 {
			return fmt.Errorf("Traditional Chinese key not found in ListObjectsV2 results")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_CopyObject", func() error {
		srcKey := "テスト/コピー元.txt"
		srcBody := "コピー元のコンテンツです。全角文字含む。"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(srcKey),
			Body:   strings.NewReader(srcBody),
		})
		if err != nil {
			return fmt.Errorf("put source for copy: %v", err)
		}
		dstKey := "テスト/コピー先.txt"
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(dstKey),
			CopySource: aws.String(bucketName + "/" + srcKey),
		})
		if err != nil {
			return fmt.Errorf("copy multibyte key: %v", err)
		}
		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(dstKey),
		})
		if err != nil {
			return fmt.Errorf("get copied multibyte object: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("read copied body: %v", err)
		}
		if string(body) != srcBody {
			return fmt.Errorf("copied body mismatch: expected %q, got %q", srcBody, string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_DeleteObjects", func() error {
		delKeys := []string{
			"テスト/削除対象1.txt",
			"文档/删除目标2.txt",
			"文件/刪除目標3.txt",
		}
		for _, key := range delKeys {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   strings.NewReader("to-delete"),
			})
			if err != nil {
				return fmt.Errorf("put delete target %q: %v", key, err)
			}
		}
		var objs []types.ObjectIdentifier
		for _, key := range delKeys {
			objs = append(objs, types.ObjectIdentifier{Key: aws.String(key)})
		}
		delResp, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{Objects: objs},
		})
		if err != nil {
			return fmt.Errorf("delete multibyte objects: %v", err)
		}
		for _, key := range delKeys {
			found := false
			for _, d := range delResp.Deleted {
				if *d.Key == key {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("key %q not found in Deleted results", key)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_MultipartUpload", func() error {
		mpuKey := "テスト/マルチパート.txt"
		createResp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(mpuKey),
			ContentType: aws.String("application/octet-stream"),
		})
		if err != nil {
			return fmt.Errorf("create multipart upload with multibyte key: %v", err)
		}
		uploadID := createResp.UploadId
		if uploadID == nil {
			return fmt.Errorf("UploadId is nil")
		}

		part1Body := "こんにちは、パート1のデータです。"
		part2Body := "这是第二部分的数据。這是第二部分的數據。"
		uploadResp1, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(mpuKey),
			UploadId:   uploadID,
			PartNumber: aws.Int32(1),
			Body:       strings.NewReader(part1Body),
		})
		if err != nil {
			return fmt.Errorf("upload part 1: %v", err)
		}
		uploadResp2, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(mpuKey),
			UploadId:   uploadID,
			PartNumber: aws.Int32(2),
			Body:       strings.NewReader(part2Body),
		})
		if err != nil {
			return fmt.Errorf("upload part 2: %v", err)
		}

		_, err = client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(bucketName),
			Key:      aws.String(mpuKey),
			UploadId: uploadID,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: []types.CompletedPart{
					{PartNumber: aws.Int32(1), ETag: uploadResp1.ETag},
					{PartNumber: aws.Int32(2), ETag: uploadResp2.ETag},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("complete multipart upload: %v", err)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(mpuKey),
		})
		if err != nil {
			return fmt.Errorf("get completed multipart object: %v", err)
		}
		body, err := io.ReadAll(getResp.Body)
		getResp.Body.Close()
		if err != nil {
			return fmt.Errorf("read multipart body: %v", err)
		}
		expected := part1Body + part2Body
		if string(body) != expected {
			return fmt.Errorf("multipart body mismatch: expected %q, got %q", expected, string(body))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_ObjectTagging", func() error {
		tagKey := "文档/标签测试.txt"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(tagKey),
			Body:   strings.NewReader("tag-test"),
		})
		if err != nil {
			return fmt.Errorf("put tagging test object: %v", err)
		}
		_, err = client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(tagKey),
			Tagging: &types.Tagging{
				TagSet: []types.Tag{
					{Key: aws.String("環境"), Value: aws.String("テスト")},
					{Key: aws.String("说明"), Value: aws.String("简体标签")},
					{Key: aws.String("說明"), Value: aws.String("繁體標籤")},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put multibyte object tagging: %v", err)
		}
		tagResp, err := client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(tagKey),
		})
		if err != nil {
			return fmt.Errorf("get multibyte object tagging: %v", err)
		}
		expectedTags := map[string]string{
			"環境": "テスト",
			"说明": "简体标签",
			"說明": "繁體標籤",
		}
		for k, v := range expectedTags {
			found := false
			for _, t := range tagResp.TagSet {
				if *t.Key == k && *t.Value == v {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("tag %q=%q not found in response", k, v)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "MultiByte_ContentTypeRoundtrip", func() error {
		ctCases := []struct {
			key         string
			contentType string
		}{
			{"テスト/contenttype.txt", "text/plain; charset=shift_jis"},
			{"文档/contenttype.txt", "text/html; charset=gb2312"},
			{"文件/contenttype.txt", "text/html; charset=big5"},
		}
		for _, tc := range ctCases {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(tc.key),
				Body:        strings.NewReader("ct-test"),
				ContentType: aws.String(tc.contentType),
			})
			if err != nil {
				return fmt.Errorf("put ct object %q: %v", tc.key, err)
			}
			resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(tc.key),
			})
			if err != nil {
				return fmt.Errorf("head ct object %q: %v", tc.key, err)
			}
			if resp.ContentType == nil || *resp.ContentType != tc.contentType {
				return fmt.Errorf("ContentType mismatch for %q: expected %q, got %v", tc.key, tc.contentType, resp.ContentType)
			}
		}
		return nil
	}))

	// ========== CLEANUP ==========

	results = append(results, r.RunTest("s3", "DeleteBucket", func() error {
		s3CleanupBucket(ctx, client, bucketName)
		s3CleanupBucket(ctx, client, lockBucket)
		s3CleanupBucket(ctx, client, mpuBucket)
		return nil
	}))

	return results
}
