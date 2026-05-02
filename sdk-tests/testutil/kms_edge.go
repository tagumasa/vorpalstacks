package testutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSEdgeTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

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
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.ListKeys(tc.ctx, &kms.ListKeysInput{})
		if err != nil {
			return err
		}
		found := false
		for _, k := range resp.Keys {
			if aws.ToString(k.KeyId) == tc.keyID {
				found = true
				if k.KeyArn == nil || *k.KeyArn == "" {
					return fmt.Errorf("key ARN is nil or empty")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created key %q not found in ListKeys", tc.keyID)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DescribeKey_NonExistentKey", func() error {
		_, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{
			KeyId: aws.String("arn:aws:kms:us-east-1:000000000000:key/00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var nfe *types.NotFoundException
		if !errors.As(err, &nfe) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_DisabledKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		var de *types.DisabledException
		if !errors.As(err, &de) {
			return fmt.Errorf("expected DisabledException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey_DisabledKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
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
		var de *types.DisabledException
		if !errors.As(err, &de) {
			return fmt.Errorf("expected DisabledException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_InvalidWindow", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		_, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.rsaKeyID),
			Plaintext: []byte("should fail"),
		})
		if err == nil {
			return fmt.Errorf("expected error for encrypt with SIGN_VERIFY key")
		}
		var ikue *types.InvalidKeyUsageException
		if !errors.As(err, &ikue) {
			return fmt.Errorf("expected InvalidKeyUsageException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_WrongKeyUsage", func() error {
		if tc.hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		_, err := tc.client.Encrypt(tc.ctx, &kms.EncryptInput{
			KeyId:     aws.String(tc.hmacKeyID),
			Plaintext: []byte("test"),
		})
		if err == nil {
			return fmt.Errorf("expected error for encrypt with HMAC key")
		}
		var ikue *types.InvalidKeyUsageException
		if !errors.As(err, &ikue) {
			return fmt.Errorf("expected InvalidKeyUsageException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ReEncrypt_InvalidCiphertext", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		if tc.rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		_, err := tc.client.Sign(tc.ctx, &kms.SignInput{
			KeyId:            aws.String(tc.rsaKeyID),
			Message:          []byte("test"),
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: "INVALID_ALGORITHM",
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid algorithm")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DisableKey_NonExistent", func() error {
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
		_, err := tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var nfe *types.NotFoundException
		if !errors.As(err, &nfe) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "EnableKey_NonExistent", func() error {
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
		_, err := tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		var nfe *types.NotFoundException
		if !errors.As(err, &nfe) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_NonExistent", func() error {
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
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
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
		_, err := tc.client.GetPublicKey(tc.ctx, &kms.GetPublicKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListGrants_NonExistent", func() error {
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
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
		var nfe *types.NotFoundException
		if !errors.As(err, &nfe) {
			return fmt.Errorf("expected NotFoundException, got: %T: %v", err, err)
		}
		return nil
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
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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

	return results
}
