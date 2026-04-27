package testutil

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *TestRunner) s3EncryptionTests(ctx context.Context, client *s3.Client, ts string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "SSES3_PutGetRoundtrip", func() error {
		bucket := s3Bucket(ts, "enc-aes")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		body := "encrypted-with-sse-s3"
		putResp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("sse-s3.txt"),
			Body:                 strings.NewReader(body),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if putResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected ServerSideEncryption AES256, got %s", putResp.ServerSideEncryption)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("sse-s3.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer getResp.Body.Close()
		gotBody, err := io.ReadAll(getResp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(gotBody) != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}

		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("sse-s3.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected HeadObject ServerSideEncryption AES256, got %s", headResp.ServerSideEncryption)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "SSES3_BucketDefaultEncryption", func() error {
		bucket := s3Bucket(ts, "enc-default")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		_, err = client.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
			Bucket: aws.String(bucket),
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
		if err != nil {
			return fmt.Errorf("PutBucketEncryption failed: %w", err)
		}

		body := "default-encrypted"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("default-enc.txt"),
			Body:   strings.NewReader(body),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("default-enc.txt"),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer getResp.Body.Close()
		gotBody, err := io.ReadAll(getResp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(gotBody) != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}

		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("default-enc.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected HeadObject ServerSideEncryption AES256, got %s", headResp.ServerSideEncryption)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "SSEC_PutGetRoundtrip", func() error {
		bucket := s3Bucket(ts, "enc-cust")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		key := make([]byte, 32)
		for i := range key {
			key[i] = byte(i)
		}
		encodedKey := base64.StdEncoding.EncodeToString(key)
		keyMD5 := md5.Sum(key)
		encodedMD5 := base64.StdEncoding.EncodeToString(keyMD5[:])

		body := "customer-encrypted-data"
		putResp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec.txt"),
			Body:                 strings.NewReader(body),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}
		if putResp.SSECustomerAlgorithm == nil || *putResp.SSECustomerAlgorithm != "AES256" {
			return fmt.Errorf("expected SSECustomerAlgorithm AES256, got %v", putResp.SSECustomerAlgorithm)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		defer getResp.Body.Close()
		gotBody, err := io.ReadAll(getResp.Body)
		if err != nil {
			return fmt.Errorf("ReadAll failed: %w", err)
		}
		if string(gotBody) != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}

		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ContentLength == nil || *headResp.ContentLength != int64(len(body)) {
			return fmt.Errorf("expected ContentLength %d, got %v", len(body), headResp.ContentLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "SSEC_GetWithoutKeyFails", func() error {
		bucket := s3Bucket(ts, "enc-cust")
		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("ssec.txt"),
		})
		if err == nil {
			return fmt.Errorf("expected error when getting SSE-C object without customer key, got nil")
		}
		return nil
	}))

	return results
}
