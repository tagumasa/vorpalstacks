package testutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

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

		resp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("range-test.txt"),
			Range:  aws.String("bytes=0-4"),
		})
		if err != nil {
			return fmt.Errorf("GetObject with Range failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(body) != "01234" {
			return fmt.Errorf("expected body %q, got %q", "01234", string(body))
		}
		if resp.ContentRange == nil {
			return fmt.Errorf("ContentRange is nil")
		}
		if !strings.Contains(*resp.ContentRange, "bytes 0-4/20") {
			return fmt.Errorf("expected ContentRange to contain %q, got %q", "bytes 0-4/20", *resp.ContentRange)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "StorageClass_STANDARD_IA", func() error {
		content := "standard-ia"
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String("sc-ia.txt"),
			Body:         strings.NewReader(content),
			StorageClass: types.StorageClassStandardIa,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sc-ia.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.StorageClass != types.StorageClassStandardIa {
			return fmt.Errorf("expected StorageClass %s, got %s", types.StorageClassStandardIa, resp.StorageClass)
		}
		if resp.ContentLength == nil || *resp.ContentLength != int64(len(content)) {
			return fmt.Errorf("expected ContentLength %d, got %v", len(content), resp.ContentLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "StorageClass_GLACIER", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String("sc-glacier.txt"),
			Body:         strings.NewReader("glacier"),
			StorageClass: types.StorageClassGlacier,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sc-glacier.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.StorageClass != types.StorageClassGlacier {
			return fmt.Errorf("expected StorageClass %s, got %s", types.StorageClassGlacier, resp.StorageClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "StorageClass_ONEZONE_IA", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String("sc-1ia.txt"),
			Body:         strings.NewReader("onezone-ia"),
			StorageClass: types.StorageClassOnezoneIa,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sc-1ia.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.StorageClass != types.StorageClassOnezoneIa {
			return fmt.Errorf("expected StorageClass %s, got %s", types.StorageClassOnezoneIa, resp.StorageClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "StorageClass_INTELLIGENT_TIERING", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String("sc-it.txt"),
			Body:         strings.NewReader("intelligent-tiering"),
			StorageClass: types.StorageClassIntelligentTiering,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sc-it.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.StorageClass != types.StorageClassIntelligentTiering {
			return fmt.Errorf("expected StorageClass %s, got %s", types.StorageClassIntelligentTiering, resp.StorageClass)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "StorageClass_REDUCED_REDUNDANCY", func() error {
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String("sc-rr.txt"),
			Body:         strings.NewReader("reduced-redundancy"),
			StorageClass: types.StorageClassReducedRedundancy,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		resp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("sc-rr.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if resp.StorageClass != types.StorageClassReducedRedundancy {
			return fmt.Errorf("expected StorageClass %s, got %s", types.StorageClassReducedRedundancy, resp.StorageClass)
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

	return results
}
