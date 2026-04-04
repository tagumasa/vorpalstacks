package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
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
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
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

		_, err = client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyID),
			Plaintext: []byte("should fail"),
		})
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
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

	// === EXPANDED TESTS ===

	results = append(results, r.RunTest("kms", "ListKeys_Basic", func() error {
		resp, err := client.ListKeys(ctx, &kms.ListKeysInput{})
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
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.ListKeys(ctx, &kms.ListKeysInput{})
		if err != nil {
			return err
		}
		found := false
		for _, k := range resp.Keys {
			if aws.ToString(k.KeyId) == keyID {
				found = true
				if k.KeyArn == nil || *k.KeyArn == "" {
					return fmt.Errorf("key ARN is nil or empty")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created key %q not found in ListKeys", keyID)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_ContentVerify", func() error {
		desc := fmt.Sprintf("ContentVerify %d", time.Now().UnixNano())
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String(desc),
		})
		if err != nil {
			return err
		}
		m := resp.KeyMetadata
		if m.KeyId == nil || *m.KeyId == "" {
			return fmt.Errorf("KeyId is nil or empty")
		}
		if m.Arn == nil || *m.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		if m.KeyState != types.KeyStateEnabled {
			return fmt.Errorf("expected KeyState=Enabled, got %s", m.KeyState)
		}
		if !m.Enabled {
			return fmt.Errorf("expected Enabled=true")
		}
		if m.KeyUsage != types.KeyUsageTypeEncryptDecrypt {
			return fmt.Errorf("expected KeyUsage=ENCRYPT_DECRYPT, got %s", m.KeyUsage)
		}
		if m.KeySpec != types.KeySpecSymmetricDefault {
			return fmt.Errorf("expected KeySpec=SYMMETRIC_DEFAULT, got %s", m.KeySpec)
		}
		if m.Origin != types.OriginTypeAwsKms {
			return fmt.Errorf("expected Origin=AWS_KMS, got %s", m.Origin)
		}
		if m.KeyManager != types.KeyManagerTypeCustomer {
			return fmt.Errorf("expected KeyManager=CUSTOMER, got %s", m.KeyManager)
		}
		if m.Description == nil || *m.Description != desc {
			return fmt.Errorf("Description mismatch")
		}
		if m.EncryptionAlgorithms == nil || len(m.EncryptionAlgorithms) == 0 {
			return fmt.Errorf("EncryptionAlgorithms is empty")
		}
		defer client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               m.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})
		return nil
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_ReturnsKeyID", func() error {
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
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		if *resp.KeyId != keyID {
			return fmt.Errorf("expected KeyId=%q, got %q", keyID, *resp.KeyId)
		}
		if resp.DeletionDate == nil {
			return fmt.Errorf("DeletionDate is nil")
		}
		if resp.KeyState != types.KeyStatePendingDeletion {
			return fmt.Errorf("expected KeyState=PendingDeletion, got %s", resp.KeyState)
		}
		if resp.PendingWindowInDays == nil || *resp.PendingWindowInDays != 7 {
			return fmt.Errorf("expected PendingWindowInDays=7")
		}
		_, err = client.CancelKeyDeletion(ctx, &kms.CancelKeyDeletionInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("cancel deletion: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CancelKeyDeletion_RestoresEnabledState", func() error {
		desc := fmt.Sprintf("CancelRestore %d", time.Now().UnixNano())
		createResp, err := client.CreateKey(ctx, &kms.CreateKeyInput{Description: aws.String(desc)})
		if err != nil {
			return err
		}
		newKeyID := createResp.KeyMetadata.KeyId
		defer client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               newKeyID,
			PendingWindowInDays: aws.Int32(7),
		})

		_, err = client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: newKeyID})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               newKeyID,
			PendingWindowInDays: aws.Int32(7),
		})
		if err != nil {
			return fmt.Errorf("schedule deletion: %v", err)
		}

		_, err = client.CancelKeyDeletion(ctx, &kms.CancelKeyDeletionInput{KeyId: newKeyID})
		if err != nil {
			return fmt.Errorf("cancel deletion: %v", err)
		}

		descResp, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: newKeyID})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.KeyMetadata.Enabled {
			return fmt.Errorf("expected key to be Disabled after cancel (was disabled before scheduling deletion)")
		}
		if descResp.KeyMetadata.KeyState != types.KeyStateDisabled {
			return fmt.Errorf("expected KeyState=Disabled, got %s", descResp.KeyMetadata.KeyState)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "TagResource_ByAlias", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tagAlias := fmt.Sprintf("alias/tag-test-%d", time.Now().UnixNano())
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(tagAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer client.DeleteAlias(ctx, &kms.DeleteAliasInput{AliasName: aws.String(tagAlias)})

		_, err = client.TagResource(ctx, &kms.TagResourceInput{
			KeyId: aws.String(tagAlias),
			Tags: []types.Tag{
				{TagKey: aws.String("AliasTag"), TagValue: aws.String("test-value")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag by alias: %v", err)
		}

		tagResp, err := client.ListResourceTags(ctx, &kms.ListResourceTagsInput{KeyId: aws.String(tagAlias)})
		if err != nil {
			return fmt.Errorf("list tags by alias: %v", err)
		}
		found := false
		for _, t := range tagResp.Tags {
			if aws.ToString(t.TagKey) == "AliasTag" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tag set via alias not found")
		}
		client.UntagResource(ctx, &kms.UntagResourceInput{
			KeyId:   aws.String(keyID),
			TagKeys: []string{"AliasTag"},
		})
		return nil
	}))

	var rsaKeyID string
	results = append(results, r.RunTest("kms", "CreateKey_RSA", func() error {
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("RSA Sign/Verify Key"),
			KeyUsage:    types.KeyUsageTypeSignVerify,
			KeySpec:     types.KeySpecRsa2048,
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil {
			return fmt.Errorf("key metadata is nil")
		}
		rsaKeyID = *resp.KeyMetadata.KeyId
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetPublicKey_RSA", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		resp, err := client.GetPublicKey(ctx, &kms.GetPublicKeyInput{
			KeyId: aws.String(rsaKeyID),
		})
		if err != nil {
			return err
		}
		if resp.PublicKey == nil || len(resp.PublicKey) == 0 {
			return fmt.Errorf("public key is nil or empty")
		}
		if resp.KeySpec != types.KeySpecRsa2048 {
			return fmt.Errorf("expected KeySpec=RSA_2048, got %s", resp.KeySpec)
		}
		if resp.KeyUsage != types.KeyUsageTypeSignVerify {
			return fmt.Errorf("expected KeyUsage=SIGN_VERIFY, got %s", resp.KeyUsage)
		}
		if resp.SigningAlgorithms == nil || len(resp.SigningAlgorithms) == 0 {
			return fmt.Errorf("SigningAlgorithms is nil or empty")
		}
		return nil
	}))

	var signature []byte
	results = append(results, r.RunTest("kms", "Sign_RSA", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		message := []byte("message to sign")
		resp, err := client.Sign(ctx, &kms.SignInput{
			KeyId:            aws.String(rsaKeyID),
			Message:          message,
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		if err != nil {
			return err
		}
		if resp.Signature == nil || len(resp.Signature) == 0 {
			return fmt.Errorf("signature is nil or empty")
		}
		signature = resp.Signature
		return nil
	}))

	results = append(results, r.RunTest("kms", "Verify_RSA", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		if signature == nil {
			return fmt.Errorf("signature not available")
		}
		message := []byte("message to sign")
		resp, err := client.Verify(ctx, &kms.VerifyInput{
			KeyId:            aws.String(rsaKeyID),
			Message:          message,
			Signature:        signature,
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		if err != nil {
			return err
		}
		if !resp.SignatureValid {
			return fmt.Errorf("expected signature to be valid")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Verify_RSA_InvalidSignature", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		message := []byte("different message")
		badSig := make([]byte, 256)
		for i := range badSig {
			badSig[i] = 0xFF
		}
		resp, err := client.Verify(ctx, &kms.VerifyInput{
			KeyId:            aws.String(rsaKeyID),
			Message:          message,
			Signature:        badSig,
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		if err != nil {
			return err
		}
		if resp.SignatureValid {
			return fmt.Errorf("expected invalid signature")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Sign_DisabledKey", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		_, err := client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: aws.String(rsaKeyID)})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = client.Sign(ctx, &kms.SignInput{
			KeyId:            aws.String(rsaKeyID),
			Message:          []byte("test"),
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		})
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(rsaKeyID)})
		if err == nil {
			return fmt.Errorf("expected error when signing with disabled key")
		}
		var de *types.DisabledException
		if !errors.As(err, &de) {
			return fmt.Errorf("expected DisabledException, got: %T: %v", err, err)
		}
		return nil
	}))

	var hmacKeyID string
	results = append(results, r.RunTest("kms", "CreateKey_HMAC", func() error {
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("HMAC Key"),
			KeyUsage:    types.KeyUsageTypeGenerateVerifyMac,
			KeySpec:     types.KeySpecHmac256,
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil {
			return fmt.Errorf("key metadata is nil")
		}
		hmacKeyID = *resp.KeyMetadata.KeyId
		return nil
	}))

	var macValue []byte
	results = append(results, r.RunTest("kms", "GenerateMac", func() error {
		if hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		message := []byte("message to mac")
		resp, err := client.GenerateMac(ctx, &kms.GenerateMacInput{
			KeyId:        aws.String(hmacKeyID),
			Message:      message,
			MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
		})
		if err != nil {
			return err
		}
		if resp.Mac == nil || len(resp.Mac) == 0 {
			return fmt.Errorf("MAC is nil or empty")
		}
		macValue = resp.Mac
		return nil
	}))

	results = append(results, r.RunTest("kms", "VerifyMac", func() error {
		if hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		if macValue == nil {
			return fmt.Errorf("MAC value not available")
		}
		message := []byte("message to mac")
		resp, err := client.VerifyMac(ctx, &kms.VerifyMacInput{
			KeyId:        aws.String(hmacKeyID),
			Message:      message,
			Mac:          macValue,
			MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
		})
		if err != nil {
			return err
		}
		if !resp.MacValid {
			return fmt.Errorf("expected MAC to be valid")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "VerifyMac_InvalidMac", func() error {
		if hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		message := []byte("message to mac")
		badMac := make([]byte, 32)
		for i := range badMac {
			badMac[i] = 0xFF
		}
		resp, err := client.VerifyMac(ctx, &kms.VerifyMacInput{
			KeyId:        aws.String(hmacKeyID),
			Message:      message,
			Mac:          badMac,
			MacAlgorithm: types.MacAlgorithmSpecHmacSha256,
		})
		if err != nil {
			return err
		}
		if resp.MacValid {
			return fmt.Errorf("expected invalid MAC")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKeyPair", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
		resp, err := client.GenerateDataKeyPair(ctx, &kms.GenerateDataKeyPairInput{
			KeyId:       aws.String(keyID),
			KeyPairSpec: types.DataKeyPairSpecRsa2048,
		})
		if err != nil {
			return err
		}
		if resp.PrivateKeyCiphertextBlob == nil || len(resp.PrivateKeyCiphertextBlob) == 0 {
			return fmt.Errorf("private key ciphertext is nil or empty")
		}
		if resp.PrivateKeyPlaintext == nil || len(resp.PrivateKeyPlaintext) == 0 {
			return fmt.Errorf("private key plaintext is nil or empty")
		}
		if resp.PublicKey == nil || len(resp.PublicKey) == 0 {
			return fmt.Errorf("public key is nil or empty")
		}
		if resp.KeyPairSpec != types.DataKeyPairSpecRsa2048 {
			return fmt.Errorf("expected KeyPairSpec=RSA_2048, got %s", resp.KeyPairSpec)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKeyPairWithoutPlaintext", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
		resp, err := client.GenerateDataKeyPairWithoutPlaintext(ctx, &kms.GenerateDataKeyPairWithoutPlaintextInput{
			KeyId:       aws.String(keyID),
			KeyPairSpec: types.DataKeyPairSpecRsa2048,
		})
		if err != nil {
			return err
		}
		if resp.PrivateKeyCiphertextBlob == nil || len(resp.PrivateKeyCiphertextBlob) == 0 {
			return fmt.Errorf("private key ciphertext is nil or empty")
		}
		if resp.PublicKey == nil || len(resp.PublicKey) == 0 {
			return fmt.Errorf("public key is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListKeyRotations", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.EnableKeyRotation(ctx, &kms.EnableKeyRotationInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("enable rotation: %v", err)
		}
		resp, err := client.ListKeyRotations(ctx, &kms.ListKeyRotationsInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.Rotations == nil {
			return fmt.Errorf("rotations is nil")
		}
		if len(resp.Rotations) != 1 {
			return fmt.Errorf("expected 1 rotation, got %d", len(resp.Rotations))
		}
		if resp.Rotations[0].RotationType != types.RotationTypeAutomatic {
			return fmt.Errorf("expected RotationType=AUTOMATIC, got %s", resp.Rotations[0].RotationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus_ContentVerify", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := client.GetKeyRotationStatus(ctx, &kms.GetKeyRotationStatusInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if !resp.KeyRotationEnabled {
			return fmt.Errorf("expected KeyRotationEnabled=true")
		}
		if resp.RotationPeriodInDays == nil || *resp.RotationPeriodInDays != 365 {
			return fmt.Errorf("expected RotationPeriodInDays=365, got %d", aws.ToInt32(resp.RotationPeriodInDays))
		}
		return nil
	}))

	var grantID string
	results = append(results, r.RunTest("kms", "RetireGrant", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		granteePrincipal := fmt.Sprintf("arn:aws:iam::000000000000:user/RetireUser")
		createResp, err := client.CreateGrant(ctx, &kms.CreateGrantInput{
			KeyId:            aws.String(keyID),
			GranteePrincipal: aws.String(granteePrincipal),
			Operations:       []types.GrantOperation{types.GrantOperationDecrypt},
		})
		if err != nil {
			return fmt.Errorf("create grant: %v", err)
		}
		grantID = *createResp.GrantId
		_ = grantID

		_, err = client.RetireGrant(ctx, &kms.RetireGrantInput{
			GrantId: createResp.GrantId,
			KeyId:   aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("retire grant: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "RevokeGrant", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		granteePrincipal := fmt.Sprintf("arn:aws:iam::000000000000:user/RevokeUser")
		createResp, err := client.CreateGrant(ctx, &kms.CreateGrantInput{
			KeyId:            aws.String(keyID),
			GranteePrincipal: aws.String(granteePrincipal),
			Operations:       []types.GrantOperation{types.GrantOperationEncrypt},
		})
		if err != nil {
			return fmt.Errorf("create grant: %v", err)
		}

		_, err = client.RevokeGrant(ctx, &kms.RevokeGrantInput{
			KeyId:   aws.String(keyID),
			GrantId: createResp.GrantId,
		})
		if err != nil {
			return fmt.Errorf("revoke grant: %v", err)
		}

		listResp, err := client.ListGrants(ctx, &kms.ListGrantsInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("list grants: %v", err)
		}
		for _, g := range listResp.Grants {
			if aws.ToString(g.GrantId) == *createResp.GrantId {
				return fmt.Errorf("revoked grant still in list")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListResourceTags_ContentVerify", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		testTags := []types.Tag{
			{TagKey: aws.String("Env"), TagValue: aws.String("prod")},
			{TagKey: aws.String("Team"), TagValue: aws.String("backend")},
		}
		_, err := client.TagResource(ctx, &kms.TagResourceInput{
			KeyId: aws.String(keyID),
			Tags:  testTags,
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}
		resp, err := client.ListResourceTags(ctx, &kms.ListResourceTagsInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		found := 0
		for _, t := range resp.Tags {
			for _, expected := range testTags {
				if aws.ToString(t.TagKey) == *expected.TagKey && aws.ToString(t.TagValue) == *expected.TagValue {
					found++
					break
				}
			}
		}
		if found < 2 {
			return fmt.Errorf("expected to find 2 tags, found %d", found)
		}
		client.UntagResource(ctx, &kms.UntagResourceInput{
			KeyId:   aws.String(keyID),
			TagKeys: []string{"Env", "Team"},
		})
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_ByKeyARN", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
		descResp, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		keyARN := *descResp.KeyMetadata.Arn
		plaintext := []byte("encrypt-by-arn-test")
		encResp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(keyARN),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt by ARN: %v", err)
		}
		decResp, err := client.Decrypt(ctx, &kms.DecryptInput{
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
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
		testAlias := fmt.Sprintf("alias/encrypt-test-%d", time.Now().UnixNano())
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(testAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer client.DeleteAlias(ctx, &kms.DeleteAliasInput{AliasName: aws.String(testAlias)})

		plaintext := []byte("encrypt-by-alias-test")
		encResp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(testAlias),
			Plaintext: plaintext,
		})
		if err != nil {
			return fmt.Errorf("encrypt by alias: %v", err)
		}
		decResp, err := client.Decrypt(ctx, &kms.DecryptInput{
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

	results = append(results, r.RunTest("kms", "CreateKey_InvalidKeyUsageKeySpec", func() error {
		_, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("Invalid combo"),
			KeyUsage:    types.KeyUsageTypeEncryptDecrypt,
			KeySpec:     types.KeySpecHmac256,
		})
		if err == nil {
			return fmt.Errorf("expected error for ENCRYPT_DECRYPT + HMAC_256")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey_NumberOfBytes", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
		resp, err := client.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
			KeyId:         aws.String(keyID),
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

	results = append(results, r.RunTest("kms", "GenerateRandom_VariousSizes", func() error {
		for _, size := range []int32{1, 16, 128, 1024} {
			resp, err := client.GenerateRandom(ctx, &kms.GenerateRandomInput{
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

	results = append(results, r.RunTest("kms", "CreateKey_MultiRegion", func() error {
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("Multi-Region Key"),
			MultiRegion: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata.MultiRegion == nil || !*resp.KeyMetadata.MultiRegion {
			return fmt.Errorf("expected MultiRegion=true")
		}
		defer client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               resp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_WithTags", func() error {
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("Key with tags"),
			Tags: []types.Tag{
				{TagKey: aws.String("Purpose"), TagValue: aws.String("testing")},
			},
		})
		if err != nil {
			return err
		}
		tagResp, err := client.ListResourceTags(ctx, &kms.ListResourceTagsInput{
			KeyId: resp.KeyMetadata.KeyId,
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		found := false
		for _, t := range tagResp.Tags {
			if aws.ToString(t.TagKey) == "Purpose" && aws.ToString(t.TagValue) == "testing" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tag not found after create")
		}
		defer client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               resp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})
		return nil
	}))

	results = append(results, r.RunTest("kms", "DescribeKey_ByAlias", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		testAlias := fmt.Sprintf("alias/desc-test-%d", time.Now().UnixNano())
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(testAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer client.DeleteAlias(ctx, &kms.DeleteAliasInput{AliasName: aws.String(testAlias)})

		resp, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: aws.String(testAlias)})
		if err != nil {
			return fmt.Errorf("describe by alias: %v", err)
		}
		if resp.KeyMetadata == nil || aws.ToString(resp.KeyMetadata.KeyId) != keyID {
			return fmt.Errorf("key ID mismatch when describing by alias")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "UpdateKeyDescription_VerifyChange", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		newDesc := fmt.Sprintf("Verified %d", time.Now().UnixNano())
		_, err := client.UpdateKeyDescription(ctx, &kms.UpdateKeyDescriptionInput{
			KeyId:       aws.String(keyID),
			Description: aws.String(newDesc),
		})
		if err != nil {
			return err
		}
		resp, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: aws.String(keyID)})
		if err != nil {
			return err
		}
		if aws.ToString(resp.KeyMetadata.Description) != newDesc {
			return fmt.Errorf("description not updated")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Sign_InvalidAlgorithm", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		_, err := client.Sign(ctx, &kms.SignInput{
			KeyId:            aws.String(rsaKeyID),
			Message:          []byte("test"),
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: "INVALID_ALGORITHM",
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid algorithm")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_WrongKeyUsage", func() error {
		if hmacKeyID == "" {
			return fmt.Errorf("HMAC key ID not available")
		}
		_, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(hmacKeyID),
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

	results = append(results, r.RunTest("kms", "DisableKey_NonExistent", func() error {
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
		_, err := client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: aws.String(fakeKeyID)})
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
		_, err := client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(fakeKeyID)})
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
		_, err := client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
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
		_, err := client.GetPublicKey(ctx, &kms.GetPublicKeyInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListGrants_NonExistent", func() error {
		fakeKeyID := "arn:aws:kms:us-east-1:000000000000:key/ffffffff-ffff-ffff-ffff-ffffffffffff"
		_, err := client.ListGrants(ctx, &kms.ListGrantsInput{KeyId: aws.String(fakeKeyID)})
		if err == nil {
			return fmt.Errorf("expected error for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GenerateDataKey_DisabledKey", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = client.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
			KeyId:   aws.String(keyID),
			KeySpec: types.DataKeySpecAes256,
		})
		client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: aws.String(keyID)})
		if err == nil {
			return fmt.Errorf("expected error when generating data key with disabled key")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DeleteAlias_NonExistent", func() error {
		_, err := client.DeleteAlias(ctx, &kms.DeleteAliasInput{
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
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String("alias/aws/test"),
			TargetKeyId: aws.String(keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for alias/aws/ prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateAlias_WithoutPrefix", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String("no-prefix-alias"),
			TargetKeyId: aws.String(keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for alias without alias/ prefix")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "PutKeyPolicy_InvalidJSON", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.PutKeyPolicy(ctx, &kms.PutKeyPolicyInput{
			KeyId:      aws.String(keyID),
			PolicyName: aws.String("default"),
			Policy:     aws.String("not valid json {{{"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid JSON policy")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_SignVerifyKey", func() error {
		if rsaKeyID == "" {
			return fmt.Errorf("RSA key ID not available")
		}
		_, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:     aws.String(rsaKeyID),
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

	results = append(results, r.RunTest("kms", "ReEncrypt_InvalidCiphertext", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.ReEncrypt(ctx, &kms.ReEncryptInput{
			CiphertextBlob:   []byte("not valid ciphertext"),
			DestinationKeyId: aws.String(keyID),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid ciphertext")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListAliases_FilterByKeyID", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		filterAlias := fmt.Sprintf("alias/filter-test-%d", time.Now().UnixNano())
		_, err := client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(filterAlias),
			TargetKeyId: aws.String(keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		defer client.DeleteAlias(ctx, &kms.DeleteAliasInput{AliasName: aws.String(filterAlias)})

		resp, err := client.ListAliases(ctx, &kms.ListAliasesInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		for _, a := range resp.Aliases {
			if a.AliasName != nil && strings.HasPrefix(*a.AliasName, "alias/aws/") {
				continue
			}
			if a.TargetKeyId != nil && *a.TargetKeyId != keyID {
				return fmt.Errorf("alias %q has wrong target key %q", *a.AliasName, *a.TargetKeyId)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus_DisabledRotation", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := client.DisableKeyRotation(ctx, &kms.DisableKeyRotationInput{KeyId: aws.String(keyID)})
		if err != nil {
			return fmt.Errorf("disable rotation: %v", err)
		}
		resp, err := client.GetKeyRotationStatus(ctx, &kms.GetKeyRotationStatusInput{
			KeyId: aws.String(keyID),
		})
		if err != nil {
			return err
		}
		if resp.KeyRotationEnabled {
			return fmt.Errorf("expected KeyRotationEnabled=false")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_ExternalOrigin", func() error {
		resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
			Description: aws.String("External origin key"),
			Origin:      types.OriginTypeExternal,
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata.Origin != types.OriginTypeExternal {
			return fmt.Errorf("expected Origin=EXTERNAL, got %s", resp.KeyMetadata.Origin)
		}
		defer client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               resp.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7),
		})
		return nil
	}))

	results = append(results, r.RunTest("kms", "Encrypt_EncryptionContext", func() error {
		if keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		plaintext := []byte("context-test-data")
		encResp, err := client.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               aws.String(keyID),
			Plaintext:           plaintext,
			EncryptionContext:   map[string]string{"project": "test", "stage": "dev"},
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecSymmetricDefault,
		})
		if err != nil {
			return fmt.Errorf("encrypt: %v", err)
		}
		decResp, err := client.Decrypt(ctx, &kms.DecryptInput{
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

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("kms", "ListKeys_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgKeyIDs []string
		for i := 0; i < 5; i++ {
			resp, err := client.CreateKey(ctx, &kms.CreateKeyInput{
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
			resp, err := client.ListKeys(ctx, &kms.ListKeysInput{
				Marker: marker,
				Limit:  aws.Int32(2),
			})
			if err != nil {
				for _, kid := range pgKeyIDs {
					client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
						KeyId:               aws.String(kid),
						PendingWindowInDays: aws.Int32(7),
					})
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
			client.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
				KeyId:               aws.String(kid),
				PendingWindowInDays: aws.Int32(7),
			})
		}
		if len(allKeys) < 5 {
			return fmt.Errorf("expected at least 5 keys across pages, got %d", len(allKeys))
		}
		return nil
	}))

	return results
}
