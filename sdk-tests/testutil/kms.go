package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunKMSTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "kms",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := kms.NewFromConfig(cfg)
	ctx := context.Background()

	keyDescription := fmt.Sprintf("Test Key %d", time.Now().UnixNano())
	keyAlias := fmt.Sprintf("alias/test-key-%d", time.Now().UnixNano())

	var keyID string
	results = append(results, r.RunTest("kms", "CreateKey", func() error {
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String(keyDescription),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil || resp.KeyMetadata.KeyId == nil {
			return fmt.Errorf("key metadata is nil")
		}
		keyID = *resp.KeyMetadata.KeyId
		return nil
	}))

	results = append(results, r.RunTest("kms", "DescribeKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil {
			return fmt.Errorf("key metadata is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "EnableKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.EnableKey(ctx, &kms.EnableKeyInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DisableKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.DisableKey(ctx, &kms.DisableKeyInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "EnableKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.EnableKey(ctx, &kms.EnableKeyInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "UpdateKeyDescription", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		newDescription := fmt.Sprintf("Updated Key %d", time.Now().UnixNano())
		resp, err := client.UpdateKeyDescription(ctx, &kms.UpdateKeyDescriptionInput{
			KeyId:       aws.String(keyID),
			Description: aws.String(newDescription),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateAlias", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(keyAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases", func() error {
		resp, err := client.ListAliases(ctx, &kms.ListAliasesInput{})
		if err != nil {
			return err
		}
		if resp.Aliases == nil {
			return fmt.Errorf("aliases list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DescribeKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil {
			return fmt.Errorf("key metadata is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		plaintext := []byte("Hello, KMS!")
		resp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil {
			return fmt.Errorf("ciphertext blob is nil")
		}
		return nil
	}))

	var ciphertextBlob []byte
	results = append(results, r.RunTest("kms", "Encrypt (for Decrypt)", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		plaintext := []byte("Hello, KMS!")
		resp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return err
		}
		ciphertextBlob = resp.CiphertextBlob
		return nil
	}))

	results = append(results, r.RunTest("kms", "Decrypt", func() error {
		if ciphertextBlob == nil {
			return fmt.Errorf("ciphertext not available")
		}
		_, err := client.Decrypt(ctx, &kms.DecryptInput{
			CiphertextBlob: ciphertextBlob,
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
			KeyId:   aws.String(keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil || len(resp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext blob is nil or empty")
		}
		if resp.Plaintext == nil || len(resp.Plaintext) != 32 {
			return fmt.Errorf("expected 32-byte plaintext, got %d", len(resp.Plaintext))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKeyWithoutPlaintext", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.GenerateDataKeyWithoutPlaintext(ctx, &kms.GenerateDataKeyWithoutPlaintextInput{
			KeyId:   aws.String(keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil || len(resp.CiphertextBlob) == 0 {
			return fmt.Errorf("ciphertext blob is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateRandom", func() error {
		resp, err := client.GenerateRandom(ctx, &kms.GenerateRandomInput{
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

	results = append(results, r.RunTest("kms", "ReEncrypt", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		if ciphertextBlob == nil {
			return fmt.Errorf("ciphertext not available")
		}
		resp, err := client.ReEncrypt(ctx, &kms.ReEncryptInput{
			CiphertextBlob:   ciphertextBlob,
			DestinationKeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.CiphertextBlob == nil {
			return fmt.Errorf("ciphertext blob is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "EnableKeyRotation", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.EnableKeyRotation(ctx, &kms.EnableKeyRotationInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.GetKeyRotationStatus(ctx, &kms.GetKeyRotationStatusInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DisableKeyRotation", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.DisableKeyRotation(ctx, &kms.DisableKeyRotationInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateGrant", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		granteePrincipal := fmt.Sprintf("arn:aws:iam::000000000000:user/TestUser")
		resp, err := client.CreateGrant(ctx, &kms.CreateGrantInput{
			KeyId:            aws.String(keyID),
			GranteePrincipal: aws.String(granteePrincipal),
			Operations:       []types.GrantOperation{types.GrantOperationEncrypt},
		})
		if err != nil {
			return err
		}
		if resp.GrantToken == nil || *resp.GrantToken == "" {
			return fmt.Errorf("grant token is nil or empty")
		}
		if resp.GrantId == nil || *resp.GrantId == "" {
			return fmt.Errorf("grant ID is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListGrants", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.ListGrants(ctx, &kms.ListGrantsInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.Grants == nil {
			return fmt.Errorf("grants list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListRetirableGrants", func() error {
		resp, err := client.ListRetirableGrants(ctx, &kms.ListRetirableGrantsInput{
			RetiringPrincipal: aws.String(fmt.Sprintf("arn:aws:iam::000000000000:user/TestUser")),
		})
		if err != nil {
			return err
		}
		if resp.Grants == nil {
			return fmt.Errorf("grants list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "PutKeyPolicy", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": "kms:*",
				"Resource": "*"
			}]
		}`
		_, err := client.PutKeyPolicy(ctx, &kms.PutKeyPolicyInput{
			KeyId:      aws.String(keyID),
			PolicyName: aws.String("default"),
			Policy:     aws.String(policy),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyPolicy", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.GetKeyPolicy(ctx, &kms.GetKeyPolicyInput{
			KeyId:      aws.String(keyID),
			PolicyName: aws.String("default"),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListKeyPolicies", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.ListKeyPolicies(ctx, &kms.ListKeyPoliciesInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.PolicyNames == nil {
			return fmt.Errorf("policy names list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "TagResource", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.TagResource(ctx, &kms.TagResourceInput{
			KeyId: aws.String(keyID),
			Tags: []types.Tag{
				{TagKey: aws.String("Environment"), TagValue: aws.String("test")},
				{TagKey: aws.String("Project"), TagValue: aws.String("sdk-tests")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListResourceTags", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.ListResourceTags(ctx, &kms.ListResourceTagsInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "UntagResource", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.UntagResource(ctx, &kms.UntagResourceInput{
			KeyId:   aws.String(keyID),
			TagKeys: []string{"Environment"},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "UpdateAlias", func() error {
		newAlias := fmt.Sprintf("alias/test-key-updated-%d", time.Now().UnixNano())
		_, err := client.UpdateAlias(ctx, &kms.UpdateAliasInput{
			AliasName:   aws.String(keyAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		listResp, err := client.ListAliases(ctx, &kms.ListAliasesInput{})
		if err != nil {
			return fmt.Errorf("list aliases after update: %v", err)
		}
		found := false
		for _, a := range listResp.Aliases {
			if aws.ToString(a.AliasName) == keyAlias {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("alias %q not found after update", keyAlias)
		}
		keyAlias = newAlias
		return nil
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(7),
		})
		if err != nil {
			return err
		}
		if resp.DeletionDate == nil {
			return fmt.Errorf("deletion date is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CancelKeyDeletion", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.CancelKeyDeletion(ctx, &kms.CancelKeyDeletionInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("key ID in response is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DeleteAlias", func() error {
		resp, err := client.DeleteAlias(ctx, &kms.DeleteAliasInput{
			AliasName: aws.String(keyAlias),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("kms", "DescribeKey_NonExistentKey", func() error {
		_, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{
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

	results = append(results, r.RunTest("kms", "Encrypt_DecryptRoundtrip", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		plaintext := []byte("roundtrip-test-data-12345")
		encResp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt: %v", err)
		}
		decResp, err := client.Decrypt(ctx, &kms.DecryptInput{
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

	results = append(results, r.RunTest("kms", "GenerateDataKey_ContentVerify", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
			KeyId:   aws.String(keyID),
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

	results = append(results, r.RunTest("kms", "CreateAlias_Duplicate", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		dupAlias := fmt.Sprintf("alias/dup-test-%d", time.Now().UnixNano())
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(dupAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteAlias(ctx, &kms.DeleteAliasInput{AliasName: aws.String(dupAlias)})
		_, err = client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(dupAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate alias")
		}
		var aee *types.AlreadyExistsException
		if !errors.As(err, &aee) {
			return fmt.Errorf("expected AlreadyExistsException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_DisabledKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}
		defer client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})

		_, err = client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyID),
			Plaintext: []byte("should fail"),
		})
		if err == nil {
			return fmt.Errorf("expected error when encrypting with disabled key")
		}
		var de *types.DisabledException
		if !errors.As(err, &de) {
			return fmt.Errorf("expected DisabledException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_InvalidWindow", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(3),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid pending window (3 days, min is 7)")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases_ContainsCreated", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		testAlias := fmt.Sprintf("alias/list-test-%d", time.Now().UnixNano())
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(testAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer client.DeleteAlias(ctx, &kms.DeleteAliasInput{AliasName: aws.String(testAlias)})

		resp, err := client.ListAliases(ctx, &kms.ListAliasesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, a := range resp.Aliases {
			if a.AliasName != nil && *a.AliasName == testAlias {
				found = true
				if a.TargetKeyId == nil || *a.TargetKeyId != keyID {
					return fmt.Errorf("alias target key mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created alias %q not found in ListAliases", testAlias)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyPolicy_ContentVerify", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":"kms:*","Resource":"*"}]}`
		_, err := client.PutKeyPolicy(ctx, &kms.PutKeyPolicyInput{
			KeyId:      aws.String(keyID),
			PolicyName: aws.String("default"),
			Policy:     aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}
		resp, err := client.GetKeyPolicy(ctx, &kms.GetKeyPolicyInput{
			KeyId:      aws.String(keyID),
			PolicyName: aws.String("default"),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		if *resp.Policy != policy {
			return fmt.Errorf("policy content mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ReEncrypt_WithDifferentKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		newKeyResp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("ReEncrypt target key"),
		})
		if err != nil {
			return fmt.Errorf("create key: %v", err)
		}
		defer client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               newKeyResp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})

		plaintext := []byte("re-encrypt-test")
		encResp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyID),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt: %v", err)
		}

		reResp, err := client.ReEncrypt(ctx, &kms.ReEncryptInput{
			CiphertextBlob:   encResp.CiphertextBlob,
			DestinationKeyId: newKeyResp.KeyMetadata.KeyId,
		})
		if err != nil {
			return fmt.Errorf("re-encrypt: %v", err)
		}

		decResp, err := client.Decrypt(ctx, &kms.DecryptInput{
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

	return results
}
