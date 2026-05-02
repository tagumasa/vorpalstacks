package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSKeyTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	keyDescription := fmt.Sprintf("Test Key %d", time.Now().UnixNano())
	keyAlias := fmt.Sprintf("alias/test-key-%d", time.Now().UnixNano())
	tc.keyAlias = keyAlias

	results = append(results, r.RunTest("kms", "CreateKey", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String(keyDescription),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil || resp.KeyMetadata.KeyId == nil {
			return fmt.Errorf("key metadata is nil")
		}
		tc.keyID = *resp.KeyMetadata.KeyId
		if resp.KeyMetadata.Arn != nil {
			tc.keyARN = *resp.KeyMetadata.Arn
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DescribeKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata == nil {
			return fmt.Errorf("key metadata is nil")
		}
		if resp.KeyMetadata.KeyId == nil || *resp.KeyMetadata.KeyId != tc.keyID {
			return fmt.Errorf("key ID mismatch: got %q", aws.ToString(resp.KeyMetadata.KeyId))
		}
		if resp.KeyMetadata.KeyState != types.KeyStateEnabled {
			return fmt.Errorf("expected KeyState=Enabled, got %s", resp.KeyMetadata.KeyState)
		}
		if !resp.KeyMetadata.Enabled {
			return fmt.Errorf("expected Enabled=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "EnableKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.EnableKey(tc.ctx, &kms.EnableKeyInput{
			KeyId: aws.String(tc.keyID),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "DisableKey", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{
			KeyId: aws.String(tc.keyID),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "UpdateKeyDescription", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		newDescription := fmt.Sprintf("Updated Key %d", time.Now().UnixNano())
		_, err := tc.client.UpdateKeyDescription(tc.ctx, &kms.UpdateKeyDescriptionInput{
			KeyId:       aws.String(tc.keyID),
			Description: aws.String(newDescription),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "CreateKey_ContentVerify", func() error {
		desc := fmt.Sprintf("ContentVerify %d", time.Now().UnixNano())
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
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
			return fmt.Errorf("arn is nil or empty")
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
			return fmt.Errorf("description mismatch")
		}
		if len(m.EncryptionAlgorithms) == 0 {
			return fmt.Errorf("EncryptionAlgorithms is empty")
		}
		tc.addCleanupKey(*m.KeyId)
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_MultiRegion", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("Multi-Region Key"),
			MultiRegion: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata.MultiRegion == nil || !*resp.KeyMetadata.MultiRegion {
			return fmt.Errorf("expected MultiRegion=true")
		}
		tc.addCleanupKey(*resp.KeyMetadata.KeyId)
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_WithTags", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("Key with tags"),
			Tags: []types.Tag{
				{TagKey: aws.String("Purpose"), TagValue: aws.String("testing")},
			},
		})
		if err != nil {
			return err
		}
		tagResp, err := tc.client.ListResourceTags(tc.ctx, &kms.ListResourceTagsInput{
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
		tc.addCleanupKey(*resp.KeyMetadata.KeyId)
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_ExternalOrigin", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("External origin key"),
			Origin:      types.OriginTypeExternal,
		})
		if err != nil {
			return err
		}
		if resp.KeyMetadata.Origin != types.OriginTypeExternal {
			return fmt.Errorf("expected Origin=EXTERNAL, got %s", resp.KeyMetadata.Origin)
		}
		if resp.KeyMetadata.KeyState != types.KeyStatePendingImport {
			return fmt.Errorf("expected KeyState=PendingImport for EXTERNAL origin, got %s", resp.KeyMetadata.KeyState)
		}
		if resp.KeyMetadata.Enabled {
			return fmt.Errorf("expected Enabled=false for EXTERNAL origin")
		}
		tc.addCleanupKey(*resp.KeyMetadata.KeyId)
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_RSA", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
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
		if resp.KeyMetadata.KeyUsage != types.KeyUsageTypeSignVerify {
			return fmt.Errorf("expected KeyUsage=SIGN_VERIFY, got %s", resp.KeyMetadata.KeyUsage)
		}
		if resp.KeyMetadata.KeySpec != types.KeySpecRsa2048 {
			return fmt.Errorf("expected KeySpec=RSA_2048, got %s", resp.KeyMetadata.KeySpec)
		}
		tc.rsaKeyID = *resp.KeyMetadata.KeyId
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_HMAC", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
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
		if resp.KeyMetadata.KeyUsage != types.KeyUsageTypeGenerateVerifyMac {
			return fmt.Errorf("expected KeyUsage=GENERATE_VERIFY_MAC, got %s", resp.KeyMetadata.KeyUsage)
		}
		if resp.KeyMetadata.KeySpec != types.KeySpecHmac256 {
			return fmt.Errorf("expected KeySpec=HMAC_256, got %s", resp.KeyMetadata.KeySpec)
		}
		tc.hmacKeyID = *resp.KeyMetadata.KeyId
		return nil
	}))

	results = append(results, r.RunTest("kms", "CreateKey_InvalidKeyUsageKeySpec", func() error {
		_, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("Invalid combo"),
			KeyUsage:    types.KeyUsageTypeEncryptDecrypt,
			KeySpec:     types.KeySpecHmac256,
		})
		if err == nil {
			return fmt.Errorf("expected error for ENCRYPT_DECRYPT + HMAC_256")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ScheduleKeyDeletion_ReturnsKeyID", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(tc.keyID),
			PendingWindowInDays: aws.Int32(7),
		})
		if err != nil {
			return err
		}
		if resp.KeyId == nil || *resp.KeyId == "" {
			return fmt.Errorf("KeyId in response is nil or empty")
		}
		if *resp.KeyId != tc.keyID {
			return fmt.Errorf("expected KeyId=%q, got %q", tc.keyID, *resp.KeyId)
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
		_, err = tc.client.CancelKeyDeletion(tc.ctx, &kms.CancelKeyDeletionInput{KeyId: aws.String(tc.keyID)})
		if err != nil {
			return fmt.Errorf("cancel deletion: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "CancelKeyDeletion_RestoresEnabledState", func() error {
		desc := fmt.Sprintf("CancelRestore %d", time.Now().UnixNano())
		createResp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{Description: aws.String(desc)})
		if err != nil {
			return err
		}
		newKeyID := createResp.KeyMetadata.KeyId
		tc.addCleanupKey(*newKeyID)

		_, err = tc.client.DisableKey(tc.ctx, &kms.DisableKeyInput{KeyId: newKeyID})
		if err != nil {
			return fmt.Errorf("disable: %v", err)
		}

		_, err = tc.client.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               newKeyID,
			PendingWindowInDays: aws.Int32(7),
		})
		if err != nil {
			return fmt.Errorf("schedule deletion: %v", err)
		}

		_, err = tc.client.CancelKeyDeletion(tc.ctx, &kms.CancelKeyDeletionInput{KeyId: newKeyID})
		if err != nil {
			return fmt.Errorf("cancel deletion: %v", err)
		}

		descResp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: newKeyID})
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

	results = append(results, r.RunTest("kms", "DescribeKey_ByAlias", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		testAlias := fmt.Sprintf("alias/desc-test-%d", time.Now().UnixNano())
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(testAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		tc.addCleanupAlias(testAlias)

		resp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: aws.String(testAlias)})
		if err != nil {
			return fmt.Errorf("describe by alias: %v", err)
		}
		if resp.KeyMetadata == nil || aws.ToString(resp.KeyMetadata.KeyId) != tc.keyID {
			return fmt.Errorf("key ID mismatch when describing by alias")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "UpdateKeyDescription_VerifyChange", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		newDesc := fmt.Sprintf("Verified %d", time.Now().UnixNano())
		_, err := tc.client.UpdateKeyDescription(tc.ctx, &kms.UpdateKeyDescriptionInput{
			KeyId:       aws.String(tc.keyID),
			Description: aws.String(newDesc),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: aws.String(tc.keyID)})
		if err != nil {
			return err
		}
		if aws.ToString(resp.KeyMetadata.Description) != newDesc {
			return fmt.Errorf("description not updated")
		}
		return nil
	}))

	return results
}
