package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamSAMLTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	samlProviderName := fmt.Sprintf("TestSAML-%s", tc.ts)

	results = append(results, r.RunTest("iam", "CreateSAMLProvider", func() error {
		resp, err := tc.client.CreateSAMLProvider(tc.ctx, &iam.CreateSAMLProviderInput{
			Name:                    aws.String(samlProviderName),
			SAMLMetadataDocument:    aws.String(samlMetadata),
			AssertionEncryptionMode: types.AssertionEncryptionModeTypeAllowed,
			AddPrivateKey:           aws.String(testSAMLPrivateKey),
			Tags: []types.Tag{
				{Key: aws.String("Source"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.SAMLProviderArn == nil {
			return fmt.Errorf("saml provider arn is nil")
		}
		tc.samlProviderArn = *resp.SAMLProviderArn
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetSAMLProvider", func() error {
		resp, err := tc.client.GetSAMLProvider(tc.ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		if resp.SAMLMetadataDocument == nil || *resp.SAMLMetadataDocument == "" {
			return fmt.Errorf("saml metadata document is empty")
		}
		// AssertionEncryptionMode must round-trip
		if resp.AssertionEncryptionMode != types.AssertionEncryptionModeTypeAllowed {
			return fmt.Errorf("AssertionEncryptionMode mismatch: got %s, want Allowed", resp.AssertionEncryptionMode)
		}
		// PrivateKeyList must contain at least 1 entry
		if len(resp.PrivateKeyList) == 0 {
			return fmt.Errorf("PrivateKeyList is empty")
		}
		// KeyId must be a privateKeyIdType (len [22,64], uppercase
		// alphanumeric) — not a synthetic placeholder.
		firstKeyId := aws.ToString(resp.PrivateKeyList[0].KeyId)
		if len(firstKeyId) < 22 || len(firstKeyId) > 64 {
			return fmt.Errorf("PrivateKey KeyId length %d outside [22,64]: %q", len(firstKeyId), firstKeyId)
		}
		for _, c := range firstKeyId {
			if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				return fmt.Errorf("PrivateKey KeyId contains non-alphanumeric char: %q", firstKeyId)
			}
		}
		tc.samlPrivateKeyId = firstKeyId
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSAMLProviders", func() error {
		resp, err := tc.client.ListSAMLProviders(tc.ctx, &iam.ListSAMLProvidersInput{})
		if err != nil {
			return err
		}
		if resp.SAMLProviderList == nil {
			return fmt.Errorf("saml provider list is nil")
		}
		found := false
		for _, p := range resp.SAMLProviderList {
			if aws.ToString(p.Arn) == tc.samlProviderArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("saml provider %s not found in list", tc.samlProviderArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateSAMLProvider", func() error {
		resp, err := tc.client.UpdateSAMLProvider(tc.ctx, &iam.UpdateSAMLProviderInput{
			SAMLProviderArn:      aws.String(tc.samlProviderArn),
			SAMLMetadataDocument: aws.String(samlMetadata),
		})
		if err != nil {
			return err
		}
		if resp.SAMLProviderArn == nil {
			return fmt.Errorf("saml provider arn is nil")
		}
		return nil
	}))

	// UpdateSAMLProvider.AddPrivateKey appends a new key, then
	// RemovePrivateKey removes a specific key by KeyId.
	results = append(results, r.RunTest("iam", "UpdateSAMLProvider_AddRemovePrivateKey", func() error {
		// Add a second key via Update
		_, err := tc.client.UpdateSAMLProvider(tc.ctx, &iam.UpdateSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
			AddPrivateKey:   aws.String(testSAMLPrivateKeySecond),
		})
		if err != nil {
			return fmt.Errorf("UpdateSAMLProvider AddPrivateKey failed: %w", err)
		}

		// Verify 2 keys now
		getResp, err := tc.client.GetSAMLProvider(tc.ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return fmt.Errorf("GetSAMLProvider after AddPrivateKey: %w", err)
		}
		if len(getResp.PrivateKeyList) != 2 {
			return fmt.Errorf("expected 2 private keys after AddPrivateKey, got %d", len(getResp.PrivateKeyList))
		}

		// Remove the original key by its KeyId
		_, err = tc.client.UpdateSAMLProvider(tc.ctx, &iam.UpdateSAMLProviderInput{
			SAMLProviderArn:  aws.String(tc.samlProviderArn),
			RemovePrivateKey: aws.String(tc.samlPrivateKeyId),
		})
		if err != nil {
			return fmt.Errorf("UpdateSAMLProvider RemovePrivateKey failed: %w", err)
		}

		// Verify 1 key remains (the second one), proving selective removal
		getResp2, err := tc.client.GetSAMLProvider(tc.ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return fmt.Errorf("GetSAMLProvider after RemovePrivateKey: %w", err)
		}
		if len(getResp2.PrivateKeyList) != 1 {
			return fmt.Errorf("expected 1 private key after RemovePrivateKey, got %d", len(getResp2.PrivateKeyList))
		}
		// The remaining key must NOT be the one we removed
		remainingKeyId := aws.ToString(getResp2.PrivateKeyList[0].KeyId)
		if remainingKeyId == tc.samlPrivateKeyId {
			return fmt.Errorf("wrong key was removed: removed KeyId %s still present", tc.samlPrivateKeyId)
		}
		return nil
	}))

	// RemovePrivateKey with a non-existent KeyId must return an error
	results = append(results, r.RunTest("iam", "Error_UpdateSAMLProvider_RemovePrivateKeyNotFound", func() error {
		_, err := tc.client.UpdateSAMLProvider(tc.ctx, &iam.UpdateSAMLProviderInput{
			SAMLProviderArn:  aws.String(tc.samlProviderArn),
			RemovePrivateKey: aws.String("INVALIDKEYIDDOESNOTEXIST1234567"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent RemovePrivateKey KeyId")
		}
		return nil
	}))

	// SAML tags
	results = append(results, r.RunTest("iam", "TagSAMLProvider", func() error {
		_, err := tc.client.TagSAMLProvider(tc.ctx, &iam.TagSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListSAMLProviderTags(tc.ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return fmt.Errorf("ListSAMLProviderTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagSAMLProvider")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSAMLProviderTags", func() error {
		resp, err := tc.client.ListSAMLProviderTags(tc.ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		if !iamTagPresent(resp.Tags, "Source", "test") {
			return fmt.Errorf("Source=test tag not found (from create)")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagSAMLProvider", func() error {
		_, err := tc.client.UntagSAMLProvider(tc.ctx, &iam.UntagSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
			TagKeys:         []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListSAMLProviderTags(tc.ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		if iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment tag should be removed")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteSAMLProvider", func() error {
		_, err := tc.client.DeleteSAMLProvider(tc.ctx, &iam.DeleteSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetSAMLProvider(tc.ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err == nil {
			return fmt.Errorf("GetSAMLProvider should fail after DeleteSAMLProvider")
		}
		return nil
	}))

	return results
}
