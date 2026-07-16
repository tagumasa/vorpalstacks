package testutil

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func (r *TestRunner) s3BucketTests(ctx context.Context, client *s3.Client, ts string, bucketName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "CreateBucket", func() error {
		resp, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		if resp.Location == nil {
			return fmt.Errorf("Location is nil")
		}
		if *resp.Location != "/"+bucketName {
			return fmt.Errorf("expected Location /%s, got %s", bucketName, *resp.Location)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListBuckets", func() error {
		resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return fmt.Errorf("ListBuckets failed: %w", err)
		}
		if resp.Buckets == nil {
			return fmt.Errorf("Buckets is nil")
		}
		found := false
		for _, b := range resp.Buckets {
			if b.Name != nil && *b.Name == bucketName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("bucket %q not found in list", bucketName)
		}
		for _, b := range resp.Buckets {
			if b.Name != nil && *b.Name == bucketName {
				if b.CreationDate == nil {
					return fmt.Errorf("bucket %q has nil CreationDate", bucketName)
				}
				break
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "ListBuckets_SortedByName", func() error {
		sortNames := []string{"sort-z", "sort-a", "sort-m"}
		var sortBuckets []string
		for _, n := range sortNames {
			bn := s3Bucket(ts, n)
			sortBuckets = append(sortBuckets, bn)
			_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(bn),
			})
			if err != nil {
				return fmt.Errorf("CreateBucket %s failed: %w", bn, err)
			}
		}
		defer func() {
			for _, bn := range sortBuckets {
				s3CleanupBucket(ctx, client, bn)
			}
		}()

		resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return fmt.Errorf("ListBuckets failed: %w", err)
		}

		var found []string
		for _, b := range resp.Buckets {
			if b.Name == nil {
				continue
			}
			for _, sb := range sortBuckets {
				if *b.Name == sb {
					found = append(found, sb)
				}
			}
		}
		if len(found) != 3 {
			return fmt.Errorf("expected 3 sort buckets, found %d", len(found))
		}

		sorted := make([]string, len(found))
		copy(sorted, found)
		sort.Strings(sorted)
		for i := range found {
			if found[i] != sorted[i] {
				return fmt.Errorf("buckets not in lexicographic order: got %v, want %v", found, sorted)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadBucket", func() error {
		resp, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("HeadBucket failed: %w", err)
		}
		if resp.BucketRegion == nil {
			return fmt.Errorf("BucketRegion is nil")
		}
		if *resp.BucketRegion == "" {
			return fmt.Errorf("BucketRegion is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetBucketLocation", func() error {
		resp, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("GetBucketLocation failed: %w", err)
		}
		loc := string(resp.LocationConstraint)
		if loc == "" || loc == "us-east-1" {
			return nil
		}
		return fmt.Errorf("unexpected location: %s", loc)
	}))

	results = append(results, r.RunTest("s3", "CreateBucket_DuplicateName", func() error {
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate bucket name")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		code := apiErr.ErrorCode()
		if code != "BucketAlreadyOwnedByYou" && code != "BucketAlreadyExists" {
			return fmt.Errorf("expected BucketAlreadyOwnedByYou or BucketAlreadyExists, got %s: %v", code, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "HeadBucket_NonExistent", func() error {
		_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String("nonexistent-bucket-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent bucket")
		}
		var notFound *types.NotFound
		if !errors.As(err, &notFound) {
			return fmt.Errorf("expected NotFound, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "DeleteBucket_NotEmpty", func() error {
		notEmptyBucket := s3Bucket(ts, "notempty")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(notEmptyBucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, notEmptyBucket)

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(notEmptyBucket),
			Key:    aws.String("blocker-object.txt"),
			Body:   nil,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(notEmptyBucket),
		})
		if err == nil {
			return fmt.Errorf("expected error when deleting non-empty bucket")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "BucketNotEmpty" {
			return fmt.Errorf("expected BucketNotEmpty, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	return results
}
