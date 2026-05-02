package testutil

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"vorpalstacks-sdk-tests/config"
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
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected GetObject ServerSideEncryption AES256, got %s", getResp.ServerSideEncryption)
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
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected GetObject ServerSideEncryption AES256, got %s", getResp.ServerSideEncryption)
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
		if getResp.SSECustomerAlgorithm == nil || *getResp.SSECustomerAlgorithm != "AES256" {
			return fmt.Errorf("expected GetObject SSECustomerAlgorithm AES256, got %v", getResp.SSECustomerAlgorithm)
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

	results = append(results, r.RunTest("s3", "SSEKMS_PutGetRoundtrip", func() error {
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 SSE-KMS test key"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey failed: %w", err)
		}
		keyID := *createKeyResp.KeyMetadata.KeyId
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "enc-kms")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		body := "kms-encrypted-data"
		putResp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("sse-kms.txt"),
			Body:                 strings.NewReader(body),
			ServerSideEncryption: types.ServerSideEncryptionAwsKms,
			SSEKMSKeyId:          aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("PutObject SSE-KMS failed: %w", err)
		}
		if putResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected ServerSideEncryption aws:kms, got %s", putResp.ServerSideEncryption)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("sse-kms.txt"),
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
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected GetObject ServerSideEncryption aws:kms, got %s", getResp.ServerSideEncryption)
		}

		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("sse-kms.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected HeadObject ServerSideEncryption aws:kms, got %s", headResp.ServerSideEncryption)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "SSEKMS_BucketDefaultEncryption", func() error {
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 SSE-KMS bucket default test key"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey failed: %w", err)
		}
		keyID := *createKeyResp.KeyMetadata.KeyId
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "enc-kms-default")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
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
							SSEAlgorithm:   types.ServerSideEncryptionAwsKms,
							KMSMasterKeyID: aws.String(keyID),
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketEncryption SSE-KMS failed: %w", err)
		}

		body := "kms-default-encrypted"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("kms-default.txt"),
			Body:   strings.NewReader(body),
		})
		if err != nil {
			return fmt.Errorf("PutObject with bucket default KMS encryption failed: %w", err)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("kms-default.txt"),
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
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected GetObject ServerSideEncryption aws:kms, got %s", getResp.ServerSideEncryption)
		}
		return nil
	}))

	return results
}
