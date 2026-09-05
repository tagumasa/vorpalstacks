package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"vorpalstacks-sdk-tests/config"
)

// s3CreateTestKMSKey creates a throwaway KMS key for S3 encryption tests and
// returns its ARN and key ID together with a cleanup closure that schedules
// the key's deletion with a seven-day pending window.
func (r *TestRunner) s3CreateTestKMSKey(ctx context.Context, description string) (string, string, func(), error) {
	kmsCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("load KMS config: %w", err)
	}
	kmsClient := kms.NewFromConfig(kmsCfg)

	createKeyResp, err := kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
		Description: aws.String(description),
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("CreateKey failed: %w", err)
	}
	keyID := aws.ToString(createKeyResp.KeyMetadata.KeyId)
	keyArn := aws.ToString(createKeyResp.KeyMetadata.Arn)
	cleanup := func() {
		kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(7),
		})
	}
	return keyArn, keyID, cleanup, nil
}

func (r *TestRunner) s3EncryptionTests(ctx context.Context, client *s3.Client, ts string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("s3", "SSES3_PutGetRoundtrip", func() error {
		bucket := s3Bucket(ts, "enc-aes")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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

		getResp, gotBody, err := s3GetRead(ctx, client, bucket, "sse-s3.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected GetObject ServerSideEncryption AES256, got %s", getResp.ServerSideEncryption)
		}

		headResp, err := s3HeadObject(ctx, client, bucket, "sse-s3.txt")
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
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		_, err := client.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
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
		if _, err := s3PutObject(ctx, client, bucket, "default-enc.txt", body); err != nil {
			return err
		}

		getResp, gotBody, err := s3GetRead(ctx, client, bucket, "default-enc.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAes256 {
			return fmt.Errorf("expected GetObject ServerSideEncryption AES256, got %s", getResp.ServerSideEncryption)
		}

		headResp, err := s3HeadObject(ctx, client, bucket, "default-enc.txt")
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
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		encodedKey, encodedMD5 := ssecTestKey()

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

		getResp, gotBody, err := s3GetAndRead(ctx, client, &s3.GetObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("ssec.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
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
		_, keyID, keyCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 SSE-KMS test key")
		if err != nil {
			return err
		}
		defer keyCleanup()

		bucket := s3Bucket(ts, "enc-kms")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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

		getResp, gotBody, err := s3GetRead(ctx, client, bucket, "sse-kms.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected GetObject ServerSideEncryption aws:kms, got %s", getResp.ServerSideEncryption)
		}

		headResp, err := s3HeadObject(ctx, client, bucket, "sse-kms.txt")
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected HeadObject ServerSideEncryption aws:kms, got %s", headResp.ServerSideEncryption)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "SSEKMS_BucketDefaultEncryption", func() error {
		_, keyID, keyCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 SSE-KMS bucket default test key")
		if err != nil {
			return err
		}
		defer keyCleanup()

		bucket := s3Bucket(ts, "enc-kms-default")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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
		if _, err := s3PutObject(ctx, client, bucket, "kms-default.txt", body); err != nil {
			return err
		}

		getResp, gotBody, err := s3GetRead(ctx, client, bucket, "kms-default.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}
		if getResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected GetObject ServerSideEncryption aws:kms, got %s", getResp.ServerSideEncryption)
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_SSES3ToKMS", func() error {
		keyArn, _, keyCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 UpdateObjectEncryption SSE-S3 source test key")
		if err != nil {
			return err
		}
		defer keyCleanup()

		bucket := s3Bucket(ts, "uoe-sse3")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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

		headBefore, err := s3HeadObject(ctx, client, bucket, "uoe.txt")
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

		headAfter, err := s3HeadObject(ctx, client, bucket, "uoe.txt")
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

		getResp, gotBody, err := s3GetRead(ctx, client, bucket, "uoe.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
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
		keyArn, _, keyCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 UpdateObjectEncryption pre-versioning test key")
		if err != nil {
			return err
		}
		defer keyCleanup()

		bucket := s3Bucket(ts, "uoe-prever")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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

		if err := s3EnableVersioning(ctx, client, bucket); err != nil {
			return err
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

		head, err := s3HeadObject(ctx, client, bucket, "prever.txt")
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if head.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected ServerSideEncryption aws:kms, got %s", head.ServerSideEncryption)
		}
		if aws.ToString(head.SSEKMSKeyId) != keyArn {
			return fmt.Errorf("expected SSEKMSKeyId %s, got %s", keyArn, aws.ToString(head.SSEKMSKeyId))
		}

		_, gotBody, err := s3GetRead(ctx, client, bucket, "prever.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_KMSRotation", func() error {
		keyAArn, _, keyACleanup, err := r.s3CreateTestKMSKey(ctx, "S3 UpdateObjectEncryption rotation key A")
		if err != nil {
			return err
		}
		defer keyACleanup()
		keyBArn, _, keyBCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 UpdateObjectEncryption rotation key B")
		if err != nil {
			return err
		}
		defer keyBCleanup()

		bucket := s3Bucket(ts, "uoe-rotate")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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

		headResp, err := s3HeadObject(ctx, client, bucket, "rotate.txt")
		if err != nil {
			return fmt.Errorf("HeadObject failed: %w", err)
		}
		if headResp.ServerSideEncryption != types.ServerSideEncryptionAwsKms {
			return fmt.Errorf("expected ServerSideEncryption aws:kms, got %s", headResp.ServerSideEncryption)
		}
		if aws.ToString(headResp.SSEKMSKeyId) != keyBArn {
			return fmt.Errorf("expected SSEKMSKeyId %s after rotation, got %s", keyBArn, aws.ToString(headResp.SSEKMSKeyId))
		}

		_, gotBody, err := s3GetRead(ctx, client, bucket, "rotate.txt")
		if err != nil {
			return fmt.Errorf("GetObject failed: %w", err)
		}
		if gotBody != body {
			return fmt.Errorf("expected body %q, got %q", body, string(gotBody))
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_UnencryptedSource", func() error {
		keyArn, _, keyCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 UpdateObjectEncryption unencrypted-source test key")
		if err != nil {
			return err
		}
		defer keyCleanup()

		bucket := s3Bucket(ts, "uoe-plain")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		if _, err := s3PutObject(ctx, client, bucket, "plain.txt", "unencrypted"); err != nil {
			return err
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
		if err := expectAWSErrorCode(err, "InvalidRequest"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("s3", "UpdateObjectEncryption_NonexistentKey", func() error {
		keyArn, _, keyCleanup, err := r.s3CreateTestKMSKey(ctx, "S3 UpdateObjectEncryption missing-key test key")
		if err != nil {
			return err
		}
		defer keyCleanup()

		bucket := s3Bucket(ts, "uoe-missing")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
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
		return expectAWSErrorCode(err, "NoSuchKey")
	}))

	results = append(results, r.RunTest("s3", "GetObject_SSECOnPlainObject", func() error {
		bucket := s3Bucket(ts, "ssec-mismatch")
		if err := s3CreateBucket(ctx, client, bucket); err != nil {
			return err
		}
		defer s3CleanupBucket(ctx, client, bucket)

		if _, err := s3PutObject(ctx, client, bucket, "plain.txt", "unencrypted content"); err != nil {
			return err
		}

		encodedKey, encodedMD5 := ssecTestKey()

		_, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket:               aws.String(bucket),
			Key:                  aws.String("plain.txt"),
			SSECustomerAlgorithm: aws.String("AES256"),
			SSECustomerKey:       aws.String(encodedKey),
			SSECustomerKeyMD5:    aws.String(encodedMD5),
		})
		if err == nil {
			return fmt.Errorf("expected error getting non-SSE-C object with SSE-C parameters, got nil")
		}
		if err := expectAWSErrorCode(err, "InvalidRequest"); err != nil {
			return err
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected API error, got: %T: %v", err, err)
		}
		if !strings.Contains(apiErr.ErrorMessage(), "not applicable to this object") {
			return fmt.Errorf("expected 'not applicable to this object' message, got: %s", apiErr.ErrorMessage())
		}
		return nil
	}))

	// Overwriting an object with mismatched SSE settings is rejected with
	// EncryptionTypeMismatch in both directions: a plain PUT against an
	// SSE-C object, and an SSE-C PUT against a plain object.
	results = append(results, r.RunTest("s3", "PutObject_EncryptionTypeMismatch", func() error {
		for _, c := range []struct {
			name         string
			bucketSuffix string
			firstPutSSEC bool // whether the pre-existing object is SSE-C
		}{
			{name: "plain overwrite of an SSE-C object", bucketSuffix: "ssec-overwrite", firstPutSSEC: true},
			{name: "SSE-C overwrite of a plain object", bucketSuffix: "ssec-plain-src", firstPutSSEC: false},
		} {
			bucket := s3Bucket(ts, c.bucketSuffix)
			if err := s3CreateBucket(ctx, client, bucket); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
			defer s3CleanupBucket(ctx, client, bucket)

			encodedKey, encodedMD5 := ssecTestKey()
			ssecPut := func(key, body string) error {
				_, err := client.PutObject(ctx, &s3.PutObjectInput{
					Bucket:               aws.String(bucket),
					Key:                  aws.String(key),
					Body:                 strings.NewReader(body),
					SSECustomerAlgorithm: aws.String("AES256"),
					SSECustomerKey:       aws.String(encodedKey),
					SSECustomerKeyMD5:    aws.String(encodedMD5),
				})
				return err
			}

			var err error
			if c.firstPutSSEC {
				if err = ssecPut("ssec.txt", "customer-encrypted"); err != nil {
					return fmt.Errorf("%s: PutObject SSE-C failed: %w", c.name, err)
				}
				_, err = s3PutObject(ctx, client, bucket, "ssec.txt", "plain overwrite")
			} else {
				if _, err := s3PutObject(ctx, client, bucket, "plain.txt", "plain content"); err != nil {
					return fmt.Errorf("%s: %w", c.name, err)
				}
				err = ssecPut("plain.txt", "ssec overwrite")
			}
			if err == nil {
				return fmt.Errorf("%s: expected EncryptionTypeMismatch for the mismatched overwrite, got nil", c.name)
			}
			if err := expectAWSErrorCode(err, "EncryptionTypeMismatch"); err != nil {
				return fmt.Errorf("%s: %v", c.name, err)
			}
		}
		return nil
	}))

	return results
}
