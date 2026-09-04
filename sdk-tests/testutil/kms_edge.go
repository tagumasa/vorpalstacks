package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSEdgeTests(tc *kmsTestContext) []TestResult {
	var results []TestResult
	reg := tc.region
	acct := tc.accountID

	results = append(results, r.RunTest("kms", "ListKeys_Basic", func() error {
		resp, err := tc.client.ListKeys(tc.ctx, &kms.ListKeysInput{})
		if err != nil {
			return err
		}
		if resp.Keys == nil {
			return fmt.Errorf("keys list is nil")
		}
		if len(resp.Keys) == 0 {
			return fmt.Errorf("expected at least one key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListKeys_ContainsCreatedKey", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		var marker *string
		for page := 0; page < 100; page++ {
			resp, err := tc.client.ListKeys(tc.ctx, &kms.ListKeysInput{
				Marker: marker,
				Limit:  aws.Int32(1000),
			})
			if err != nil {
				return err
			}
			for _, k := range resp.Keys {
				if aws.ToString(k.KeyId) == tc.keyID {
					if k.KeyArn == nil || *k.KeyArn == "" {
						return fmt.Errorf("key ARN is nil or empty")
					}
					return nil
				}
			}
			if resp.NextMarker == nil || !resp.Truncated {
				break
			}
			marker = resp.NextMarker
		}
		return fmt.Errorf("created key %q not found in ListKeys", tc.keyID)
	}))

	results = append(results, r.RunTest("kms", "DescribeKey_NonExistentKey", func() error {
		_, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{
			KeyId: aws.String(fmt.Sprintf("arn:aws:kms:%s:%s:key/00000000-0000-0000-0000-000000000000", reg, acct)),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return kmsAssertErrorAs[types.NotFoundException](err)
	}))

	results = append(results, r.RunTest("kms", "Encrypt_DisabledKey", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: aws.String(tc.keyID)})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.keyID),
			Plaintext: []byte("should fail"),
		})
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		if err == nil {
			return fmt.Errorf("expected error when encrypting with disabled key")
		}
		return kmsAssertErrorAs[types.DisabledException](err)
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey_DisabledKey", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: aws.String(tc.keyID)})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = tc.client.GenerateDataKey(tc.ctx, &kms.GenerateDataKeyInput{
			KeyId:   aws.String(tc.keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.keyID)})
		if err == nil {
			return fmt.Errorf("expected error when generating data key with disabled key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Sign_DisabledKey", func() error {
		if err := tc.requireRSAKeyID(); err != nil {
			return err
		}
		_, err := tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: aws.String(tc.rsaKeyID)})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = tc.client.Sign(tc.ctx, &kms.SignInput{
			KeyId:            aws.String(tc.rsaKeyID),
			Message:          []byte("test"),
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(tc.rsaKeyID)})
		if err == nil {
			return fmt.Errorf("expected error when signing with disabled key")
		}
		return kmsAssertErrorAs[types.DisabledException](err)
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_InvalidWindow", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(tc.keyID),
			PendingWindowInDays: aws.Int32(3),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid pending window (3 days, min is 7)")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_SignVerifyKey", func() error {
		if err := tc.requireRSAKeyID(); err != nil {
			return err
		}
		_, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.rsaKeyID),
			Plaintext: []byte("should fail"),
		})
		if err == nil {
			return fmt.Errorf("expected error for encrypt with SIGN_VERIFY key")
		}
		return kmsAssertErrorAs[types.InvalidKeyUsageException](err)
	}))

	results = append(results, r.RunTest("kms", "Encrypt_WrongKeyUsage", func() error {
		if err := tc.requireHMACKeyID(); err != nil {
			return err
		}
		_, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.hmacKeyID),
			Plaintext: []byte("test"),
		})
		if err == nil {
			return fmt.Errorf("expected error for encrypt with HMAC key")
		}
		return kmsAssertErrorAs[types.InvalidKeyUsageException](err)
	}))

	results = append(results, r.RunTest("kms", "ReEncrypt_InvalidCiphertext", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.ReEncrypt(tc.ctx, &kms.ReEncryptInput{
			CiphertextBlob:   []byte("not valid ciphertext"),
			DestinationKeyId: aws.String(tc.keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid ciphertext")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Sign_InvalidAlgorithm", func() error {
		if err := tc.requireRSAKeyID(); err != nil {
			return err
		}
		_, err := tc.client.Sign(tc.ctx, &kms.SignInput{
			KeyId:            aws.String(tc.rsaKeyID),
			Message:          []byte("test"),
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: "INVALID_ALGORITHM",
		})
		// A SigningAlgorithm outside the Smithy enum is a shape violation:
		// the aws-json-1.1 protocol rejects it with SerializationException.
		return AssertErrorContains(err, "SerializationException")
	}))

	results = append(results, r.RunTest("kms", "DisableKey_NonExistent", func() error {
		fakeKeyID := tc.nonexistentKeyARN()
		_, err := tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return kmsAssertErrorAs[types.NotFoundException](err)
	}))

	results = append(results, r.RunTest("kms", "EnableKey_NonExistent", func() error {
		fakeKeyID := tc.nonexistentKeyARN()
		_, err := tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return kmsAssertErrorAs[types.NotFoundException](err)
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_NonExistent", func() error {
		fakeKeyID := tc.nonexistentKeyARN()
		_, err := tc.client.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(fakeKeyID),
			PendingWindowInDays: aws.Int32(7),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetPublicKey_NonExistent", func() error {
		fakeKeyID := tc.nonexistentKeyARN()
		_, err := tc.client.GetPublicKey(tc.ctx, &kms.GetPublicKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListGrants_NonExistent", func() error {
		fakeKeyID := tc.nonexistentKeyARN()
		_, err := tc.client.ListGrants(tc.ctx, &kms.ListGrantsInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DeleteAlias_NonExistent", func() error {
		_, err := tc.client.DeleteAlias(tc.ctx, &kms.DeleteAliasInput{
			AliasName: aws.String("alias/nonexistent-test-alias"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent alias")
		}
		return kmsAssertErrorAs[types.NotFoundException](err)
	}))

	results = append(results, r.RunTest("kms", "CreateAlias_AliasAWSReserved", func() error {
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String("alias/aws/test"),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for alias/aws/ prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateAlias_WithoutPrefix", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String("no-prefix-alias"),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for alias without alias/ prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "PutKeyPolicy_InvalidJSON", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.PutKeyPolicy(tc.ctx, &kms.PutKeyPolicyInput{
			KeyId:      aws.String(tc.keyID),
			PolicyName: aws.String("default"),
			Policy:     aws.String("not valid json {{{"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid JSON policy")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListKeys_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgKeyIDs []string
		for i := 0; i < 5; i++ {
			resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
				Description: aws.String(fmt.Sprintf("pag-key-%s-%d", pgTs, i)),
			})
			if err != nil {
				return fmt.Errorf("create key: %v", err)
			}
			pgKeyIDs = append(pgKeyIDs, aws.ToString(resp.KeyMetadata.KeyId))
		}

		var allKeys []string
		var marker *string
		for {
			resp, err := tc.client.ListKeys(tc.ctx, &kms.ListKeysInput{
				Marker: marker,
				Limit:  aws.Int32(2),
			})
			if err != nil {
				for _, kid := range pgKeyIDs {
					tc.scheduleDeletion(kid)
				}
				return fmt.Errorf("list keys page: %v", err)
			}
			for _, k := range resp.Keys {
				allKeys = append(allKeys, aws.ToString(k.KeyId))
			}
			if resp.Truncated && resp.NextMarker != nil {
				marker = resp.NextMarker
			} else {
				break
			}
		}

		for _, kid := range pgKeyIDs {
			tc.scheduleDeletion(kid)
		}
		if len(allKeys) < 5 {
			return fmt.Errorf("expected at least 5 keys across pages, got %d", len(allKeys))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases_Pagination", func() error {
		var aliasNames []string
		for i := 0; i < 5; i++ {
			aliasName := fmt.Sprintf("alias/pag-alias-%d-%d", time.Now().UnixNano(), i)
			_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
				AliasName:   aws.String(aliasName),
				TargetKeyId: aws.String(tc.keyID),
			})
			if err != nil {
				for _, a := range aliasNames {
					tc.client.DeleteAlias(tc.ctx, &kms.DeleteAliasInput{AliasName: aws.String(a)})
				}
				return fmt.Errorf("create alias: %v", err)
			}
			aliasNames = append(aliasNames, aliasName)
		}

		var allAliases []string
		var marker *string
		for {
			resp, err := tc.client.ListAliases(tc.ctx, &kms.ListAliasesInput{
				Marker: marker,
				Limit:  aws.Int32(2),
			})
			if err != nil {
				for _, a := range aliasNames {
					tc.client.DeleteAlias(tc.ctx, &kms.DeleteAliasInput{AliasName: aws.String(a)})
				}
				return fmt.Errorf("list aliases page: %v", err)
			}
			for _, a := range resp.Aliases {
				if a.AliasName != nil {
					allAliases = append(allAliases, *a.AliasName)
				}
			}
			if resp.Truncated && resp.NextMarker != nil {
				marker = resp.NextMarker
			} else {
				break
			}
		}

		for _, a := range aliasNames {
			tc.client.DeleteAlias(tc.ctx, &kms.DeleteAliasInput{AliasName: aws.String(a)})
		}
		foundCount := 0
		for _, created := range aliasNames {
			for _, listed := range allAliases {
				if listed == created {
					foundCount++
					break
				}
			}
		}
		if foundCount < 5 {
			return fmt.Errorf("expected 5 created aliases in paginated list, found %d", foundCount)
		}
		return nil
	}))

	// CustomKeyStoreId must be rejected because the platform does
	// not implement Custom Key Stores. Without the check the parameter
	// is silently dropped, causing the caller to believe a key was
	// created in a custom key store.
	results = append(results, r.RunTest("kms", "CreateKey_CustomKeyStoreIdRejected", func() error {
		_, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			CustomKeyStoreId: aws.String("cks-1234567890abcdef0"),
		})
		if err == nil {
			return fmt.Errorf("expected error for CustomKeyStoreId without custom key store support")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	// DisableKeyRotation must reject asymmetric keys with
	// UnsupportedOperationException, matching the EnableKeyRotation
	// behaviour.
	results = append(results, r.RunTest("kms", "DisableKeyRotation_AsymmetricKey", func() error {
		if err := tc.requireRSAKeyID(); err != nil {
			return err
		}
		_, err := tc.client.DisableKeyRotation(tc.ctx, &kms.DisableKeyRotationInput{
			KeyId: aws.String(tc.rsaKeyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for DisableKeyRotation on asymmetric key")
		}
		return AssertErrorContains(err, "UnsupportedOperationException")
	}))

	// GenerateDataKeyPair with an unsupported KeyPairSpec (SM2) must
	// return ValidationException, not KMSInternalException. The HSM
	// backend does not implement SM2 key pair generation.
	results = append(results, r.RunTest("kms", "GenerateDataKeyPair_UnsupportedSpec", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.GenerateDataKeyPair(tc.ctx, &kms.GenerateDataKeyPairInput{
			KeyId:       aws.String(tc.keyID),
			KeyPairSpec: types.DataKeyPairSpecSm2,
		})
		if err == nil {
			return fmt.Errorf("expected error for unsupported KeyPairSpec SM2")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	return results
}
