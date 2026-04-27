package testutil

import (
	"context"
	"fmt"
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
