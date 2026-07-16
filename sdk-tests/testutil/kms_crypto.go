package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSCryptoTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "Encrypt", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		plaintext := []byte("Hello, KMS!")
		resp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil || len(resp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext blob is nil or empty")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt (for Decrypt)", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		plaintext := []byte("Hello, KMS!")
		resp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil {
			return fmt.Errorf("ciphertext blob is nil")
		}
		tc.ciphertextBlob = resp.CiphertextBlob
		return nil
	}))

	results = append(results, r.RunTest("kms", "Decrypt", func() error {
		if tc.ciphertextBlob == nil {
			return fmt.Errorf("ciphertext not available")
		}
		resp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob: tc.ciphertextBlob,
		})
		if err != nil {
			return err
		}
		if resp.Plaintext == nil {
			return fmt.Errorf("plaintext is nil")
		}
		if string(resp.Plaintext) != "Hello, KMS!" {
			return fmt.Errorf("plaintext mismatch: got %q", string(resp.Plaintext))
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_DecryptRoundtrip", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		plaintext := []byte("roundtrip-test-data-12345")
		encResp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt: %v", err)
		}
		decResp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob: encResp.CiphertextBlob,
		})
		if err != nil {
			return fmt.Errorf("decrypt: %v", err)
		}
		if string(decResp.Plaintext) != string(plaintext) {
			return fmt.Errorf("plaintext mismatch: got %q, want %q", string(decResp.Plaintext), string(plaintext))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_ByKeyARN", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		descResp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: aws.String(tc.keyID)})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		keyARN := *descResp.KeyMetadata.Arn
		plaintext := []byte("encrypt-by-arn-test")
		encResp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyARN),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt by ARN: %v", err)
		}
		decResp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob: encResp.CiphertextBlob,
		})
		if err != nil {
			return fmt.Errorf("decrypt: %v", err)
		}
		if string(decResp.Plaintext) != string(plaintext) {
			return fmt.Errorf("plaintext mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_ByAlias", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		testAlias := fmt.Sprintf("alias/encrypt-test-%d", tc.keyID[len(tc.keyID)-8:])
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(testAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		tc.addCleanupAlias(testAlias)

		plaintext := []byte("encrypt-by-alias-test")
		encResp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(testAlias),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt by alias: %v", err)
		}
		decResp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob: encResp.CiphertextBlob,
		})
		if err != nil {
			return fmt.Errorf("decrypt: %v", err)
		}
		if string(decResp.Plaintext) != string(plaintext) {
			return fmt.Errorf("plaintext mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_EncryptionContext", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		plaintext := []byte("context-test-data")
		encResp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:               aws.String(tc.keyID),
			Plaintext:           plaintext,
			EncryptionContext:   map[string]string{"project": "test", "stage": "dev"},
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecSymmetricDefault,
		})
		if err != nil {
			return fmt.Errorf("encrypt: %v", err)
		}
		decResp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob:    encResp.CiphertextBlob,
			EncryptionContext: map[string]string{"project": "test", "stage": "dev"},
		})
		if err != nil {
			return fmt.Errorf("decrypt with context: %v", err)
		}
		if string(decResp.Plaintext) != string(plaintext) {
			return fmt.Errorf("plaintext mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ReEncrypt", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		if tc.ciphertextBlob == nil {
			return fmt.Errorf("ciphertext not available")
		}
		resp, err := tc.client.ReEncrypt(tc.ctx, &kms.ReEncryptInput{
			CiphertextBlob:   tc.ciphertextBlob,
			DestinationKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil {
			return fmt.Errorf("ciphertext blob is nil")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ReEncrypt_WithDifferentKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		newKeyResp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("ReEncrypt target key"),
		})
		if err != nil {
			return fmt.Errorf("create key: %v", err)
		}
		tc.addCleanupKey(*newKeyResp.KeyMetadata.KeyId)

		plaintext := []byte("re-encrypt-test")
		encResp, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt: %v", err)
		}

		reResp, err := tc.client.ReEncrypt(tc.ctx, &kms.ReEncryptInput{
			CiphertextBlob:   encResp.CiphertextBlob,
			DestinationKeyId: newKeyResp.KeyMetadata.KeyId,
		})
		if err != nil {
			return fmt.Errorf("re-encrypt: %v", err)
		}

		decResp, err := tc.client.Decrypt(tc.ctx, &kms.DecryptInput{
			CiphertextBlob: reResp.CiphertextBlob,
			KeyId:          newKeyResp.KeyMetadata.KeyId,
		})
		if err != nil {
			return fmt.Errorf("decrypt re-encrypted: %v", err)
		}
		if string(decResp.Plaintext) != string(plaintext) {
			return fmt.Errorf("plaintext mismatch after re-encrypt")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.GenerateDataKey(tc.ctx, &kms.GenerateDataKeyInput{
			KeyId:   aws.String(tc.keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		if err != nil {
			return err
		}
		if len(resp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext blob is nil or empty")
		}
		if resp.Plaintext == nil || len(resp.Plaintext) != 32 {
			return fmt.Errorf("expected 32-byte plaintext, got %d", len(resp.Plaintext))
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKeyWithoutPlaintext", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.GenerateDataKeyWithoutPlaintext(tc.ctx, &kms.GenerateDataKeyWithoutPlaintextInput{
			KeyId:   aws.String(tc.keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		if err != nil {
			return err
		}
		if len(resp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext blob is nil or empty")
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey_NumberOfBytes", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		resp, err := tc.client.GenerateDataKey(tc.ctx, &kms.GenerateDataKeyInput{
			KeyId:         aws.String(tc.keyID),
			NumberOfBytes: aws.Int32(64),
		})
		if err != nil {
			return err
		}
		if len(resp.Plaintext) != 64 {
			return fmt.Errorf("expected 64-byte plaintext, got %d", len(resp.Plaintext))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey_ContentVerify", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		resp, err := tc.client.GenerateDataKey(tc.ctx, &kms.GenerateDataKeyInput{
			KeyId:   aws.String(tc.keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		if err != nil {
			return err
		}
		if len(resp.Plaintext) != 32 {
			return fmt.Errorf("expected 32-byte plaintext, got %d", len(resp.Plaintext))
		}
		if len(resp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext blob is empty")
		}
		if len(resp.Plaintext) == len(resp.CiphertextBlob) {
			return fmt.Errorf("plaintext and ciphertext should have different lengths")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateRandom", func() error {
		resp, err := tc.client.GenerateRandom(tc.ctx, &kms.GenerateRandomInput{
			NumberOfBytes: aws.Int32(32),
		})
		if err != nil {
			return err
		}
		if resp.Plaintext == nil || len(resp.Plaintext) != 32 {
			return fmt.Errorf("expected 32-byte plaintext, got %d", len(resp.Plaintext))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateRandom_VariousSizes", func() error {
		for _, size := range []int32{1, 16, 128, 1024} {
			resp, err := tc.client.GenerateRandom(tc.ctx, &kms.GenerateRandomInput{
				NumberOfBytes: aws.Int32(size),
			})
			if err != nil {
				return fmt.Errorf("size %d: %v", size, err)
			}
			if len(resp.Plaintext) != int(size) {
				return fmt.Errorf("expected %d bytes, got %d", size, len(resp.Plaintext))
			}
		}
		return nil
	}))

	return results
}
