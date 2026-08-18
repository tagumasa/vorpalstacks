package testutil

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

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

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_SSES3ToKMS", func() error {
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 UpdateObjectEncryption SSE-S3 source test key"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey failed: %w", err)
		}
		keyArn := *createKeyResp.KeyMetadata.Arn
		keyID := *createKeyResp.KeyMetadata.KeyId
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "uoe-sse3")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		body := "update-encryption-source"
		putResp, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("uoe.txt"),
			Body:                 strings.NewReader(body),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		})
		if err != nil {
			return fmt.Errorf("PutObject SSE-S3 failed: %w", err)
		}
		originalETag := aws.ToString(putResp.ETag)

		headBefore, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("uoe.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject (before) failed: %w", err)
		}

		_, err = client.UpdateObjectEncryption(ctx, &s3.UpdateObjectEncryptionInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("uoe.txt"),
			ObjectEncryption: &types.ObjectEncryptionMemberSSEKMS{
				Value: types.SSEKMSEncryption{
					KMSKeyArn:        aws.String(keyArn),
					BucketKeyEnabled: aws.Bool(true),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateObjectEncryption failed: %w", err)
		}

		headAfter, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("uoe.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject (after) failed: %w", err)
		}
		if headAfter.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected ServerSideEncryption aws:kms after update, got %s", headAfter.ServerSideEncryption)
		}
		if aws.ToString(headAfter.SSEKMSKeyId) != keyArn {
			return fmt.Errorf("expected SSEKMSKeyId %s, got %s", keyArn, aws.ToString(headAfter.SSEKMSKeyId))
		}
		if aws.ToString(headAfter.ETag) != originalETag {
			return fmt.Errorf("ETag not preserved: before %s after %s", originalETag, aws.ToString(headAfter.ETag))
		}
		if headAfter.ContentLength == nil || *headAfter.ContentLength != int64(len(body)) {
			return fmt.Errorf("expected ContentLength %d, got %v", len(body), headAfter.ContentLength)
		}
		if headBefore.LastModified == nil || headAfter.LastModified == nil || !headAfter.LastModified.Equal(*headBefore.LastModified) {
			return fmt.Errorf("LastModified not preserved: before %v after %v", headBefore.LastModified, headAfter.LastModified)
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("uoe.txt"),
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

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_PreVersioningObject", func() error {
		// An SSE object written before versioning was enabled has only the
		// null-version record; encryption updates must resolve it through the
		// same fallback path as metadata reads instead of reporting NoSuchKey.
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 UpdateObjectEncryption pre-versioning test key"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey failed: %w", err)
		}
		keyArn := *createKeyResp.KeyMetadata.Arn
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               createKeyResp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "uoe-prever")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		body := "pre-versioning encryption source"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("prever.txt"),
			Body:                 strings.NewReader(body),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		})
		if err != nil {
			return fmt.Errorf("PutObject SSE-S3 failed: %w", err)
		}

		_, err = client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
			Bucket: aws.String(bucket),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})
		if err != nil {
			return fmt.Errorf("PutBucketVersioning failed: %w", err)
		}

		_, err = client.UpdateObjectEncryption(ctx, &s3.UpdateObjectEncryptionInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("prever.txt"),
			ObjectEncryption: &types.ObjectEncryptionMemberSSEKMS{
				Value: types.SSEKMSEncryption{
					KMSKeyArn: aws.String(keyArn),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateObjectEncryption on pre-versioning object failed: %w", err)
		}

		head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("prever.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if head.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected ServerSideEncryption aws:kms, got %s", head.ServerSideEncryption)
		}
		if aws.ToString(head.SSEKMSKeyId) != keyArn {
			return fmt.Errorf("expected SSEKMSKeyId %s, got %s", keyArn, aws.ToString(head.SSEKMSKeyId))
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("prever.txt"),
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
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_KMSRotation", func() error {
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		keyA, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 UpdateObjectEncryption rotation key A"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey (A) failed: %w", err)
		}
		keyAArn := *keyA.KeyMetadata.Arn
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               keyA.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})
		keyB, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 UpdateObjectEncryption rotation key B"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey (B) failed: %w", err)
		}
		keyBArn := *keyB.KeyMetadata.Arn
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               keyB.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "uoe-rotate")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		body := "kms-rotation-data"
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("rotate.txt"),
			Body:                 strings.NewReader(body),
			ServerSideEncryption: types.ServerSideEncryptionAwsKms,
			SSEKMSKeyId:          aws.String(keyAArn),
		})
		if err != nil {
			return fmt.Errorf("PutObject SSE-KMS (key A) failed: %w", err)
		}

		_, err = client.UpdateObjectEncryption(ctx, &s3.UpdateObjectEncryptionInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("rotate.txt"),
			ObjectEncryption: &types.ObjectEncryptionMemberSSEKMS{
				Value: types.SSEKMSEncryption{
					KMSKeyArn: aws.String(keyBArn),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateObjectEncryption (rotation) failed: %w", err)
		}

		headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("rotate.txt"),
		})
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected ServerSideEncryption aws:kms, got %s", headResp.ServerSideEncryption)
		}
		if aws.ToString(headResp.SSEKMSKeyId) != keyBArn {
			return fmt.Errorf("expected SSEKMSKeyId %s after rotation, got %s", keyBArn, aws.ToString(headResp.SSEKMSKeyId))
		}

		getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("rotate.txt"),
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
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_UnencryptedSource", func() error {
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 UpdateObjectEncryption unencrypted-source test key"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey failed: %w", err)
		}
		keyArn := *createKeyResp.KeyMetadata.Arn
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               createKeyResp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "uoe-plain")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("plain.txt"),
			Body:   strings.NewReader("unencrypted"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		_, err = client.UpdateObjectEncryption(ctx, &s3.UpdateObjectEncryptionInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("plain.txt"),
			ObjectEncryption: &types.ObjectEncryptionMemberSSEKMS{
				Value: types.SSEKMSEncryption{
					KMSKeyArn: aws.String(keyArn),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error updating unencrypted source object, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "InvalidRequest" {
			return fmt.Errorf("expected InvalidRequest, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_NonexistentKey", func() error {
		kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load KMS config: %w", err)
		}
		kmsClient := kms.NewFromConfig(kmsCfg)

		createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("S3 UpdateObjectEncryption missing-key test key"),
		})
		if err != nil {
			return fmt.Errorf("CreateKey failed: %w", err)
		}
		keyArn := *createKeyResp.KeyMetadata.Arn
		defer kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               createKeyResp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})

		bucket := s3Bucket(ts, "uoe-missing")
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		_, err = client.UpdateObjectEncryption(ctx, &s3.UpdateObjectEncryptionInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("does-not-exist.txt"),
			ObjectEncryption: &types.ObjectEncryptionMemberSSEKMS{
				Value: types.SSEKMSEncryption{
					KMSKeyArn: aws.String(keyArn),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for nonexistent key, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "NoSuchKey" {
			return fmt.Errorf("expected NoSuchKey, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "GetObject_SSECOnPlainObject", func() error {
		bucket := s3Bucket(ts, "ssec-mismatch")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("plain.txt"),
			Body:   strings.NewReader("unencrypted content"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		customerKey := make([]byte, 32)
		for i := range customerKey {
			customerKey[i] = byte(i)
		}
		encodedKey := base64.StdEncoding.EncodeToString(customerKey)
		keyMD5 := md5.Sum(customerKey)
		encodedMD5 := base64.StdEncoding.EncodeToString(keyMD5[:])

		_, err = client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("plain.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err == nil {
			return fmt.Errorf("expected error getting non-SSE-C object with SSE-C parameters, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "InvalidRequest" {
			return fmt.Errorf("expected InvalidRequest, got %s: %v", apiErr.ErrorCode(), err)
		}
		if !strings.Contains(apiErr.ErrorMessage(), "not applicable to this object") {
			return fmt.Errorf("expected 'not applicable to this object' message, got: %s", apiErr.ErrorMessage())
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_EncryptionTypeMismatch_OverwriteWithoutKey", func() error {
		bucket := s3Bucket(ts, "ssec-overwrite")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		customerKey := make([]byte, 32)
		for i := range customerKey {
			customerKey[i] = byte(i)
		}
		encodedKey := base64.StdEncoding.EncodeToString(customerKey)
		keyMD5 := md5.Sum(customerKey)
		encodedMD5 := base64.StdEncoding.EncodeToString(keyMD5[:])

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec.txt"),
			Body:                 strings.NewReader("customer-encrypted"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("PutObject SSE-C failed: %w", err)
		}

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("ssec.txt"),
			Body:   strings.NewReader("plain overwrite"),
		})
		if err == nil {
			return fmt.Errorf("expected EncryptionTypeMismatch overwriting SSE-C object without customer key, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "EncryptionTypeMismatch" {
			return fmt.Errorf("expected EncryptionTypeMismatch, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "PutObject_EncryptionTypeMismatch_SSECOverwritePlain", func() error {
		bucket := s3Bucket(ts, "ssec-plain-src")
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("CreateBucket failed: %w", err)
		}
		defer s3CleanupBucket(ctx, client, bucket)

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("plain.txt"),
			Body:   strings.NewReader("plain content"),
		})
		if err != nil {
			return fmt.Errorf("PutObject failed: %w", err)
		}

		customerKey := make([]byte, 32)
		for i := range customerKey {
			customerKey[i] = byte(i)
		}
		encodedKey := base64.StdEncoding.EncodeToString(customerKey)
		keyMD5 := md5.Sum(customerKey)
		encodedMD5 := base64.StdEncoding.EncodeToString(keyMD5[:])

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("plain.txt"),
			Body:                 strings.NewReader("ssec overwrite"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err == nil {
			return fmt.Errorf("expected EncryptionTypeMismatch overwriting plain object with SSE-C parameters, got nil")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if apiErr.ErrorCode() != "EncryptionTypeMismatch" {
			return fmt.Errorf("expected EncryptionTypeMismatch, got %s: %v", apiErr.ErrorCode(), err)
		}
		return nil
	}))

	return results
}
