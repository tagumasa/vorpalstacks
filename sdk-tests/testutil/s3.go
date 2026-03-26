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

	bucketName := fmt.Sprintf("test-bucket-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("s3", "CreateBucket", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		return err
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

	ts := time.Now().UnixNano()
	sortBuckets := []string{
		fmt.Sprintf("z-bucket-sort-%d", ts),
		fmt.Sprintf("a-bucket-sort-%d", ts),
		fmt.Sprintf("m-bucket-sort-%d", ts),
	}
	results = append(results, r.RunTest("s3", "ListBuckets_SortedByName", func() error {
		for _, b := range sortBuckets {
			_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(b),
			})
			if err != nil {
				return fmt.Errorf("failed to create bucket %s: %v", b, err)
			}
		}
		defer func() {
			for _, b := range sortBuckets {
				client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(b)})
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
		_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "PutObject", func() error {
		content := "Hello, World!"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
			Body:   strings.NewReader(content),
		})
		return err
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
		_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("test.txt"),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "ListObjects", func() error {
		_, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		return err
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
		_, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	results = append(results, r.RunTest("s3", "DeleteBucketTagging", func() error {
		_, err := client.DeleteBucketTagging(ctx, &s3.DeleteBucketTaggingInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS (before DeleteBucket) ===

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
		lobBucket := fmt.Sprintf("lob-bucket-%d", time.Now().UnixNano())
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(lobBucket),
		})
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer func() {
			client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(lobBucket)})
		}()

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

	results = append(results, r.RunTest("s3", "DeleteBucket", func() error {
		var objectsToDelete []types.ObjectIdentifier
		listResp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
		for _, obj := range listResp.Contents {
			objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{Key: obj.Key})
		}
		if len(objectsToDelete) > 0 {
			_, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucketName),
				Delete: &types.Delete{
					Objects: objectsToDelete,
				},
			})
			if err != nil {
				return err
			}
		}
		_, err = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		return err
	}))

	return results
}
