package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSMultiRegionTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	var multiKeyID string

	results = append(results, r.RunTest("kms", "ReplicateKey", func() error {
		resp, err := tc.client.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("Multi-Region primary for ReplicateKey"),
			MultiRegion: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create multi-region key: %v", err)
		}
		multiKeyID = *resp.KeyMetadata.KeyId
		tc.addCleanupKey(multiKeyID)

		if resp.KeyMetadata.MultiRegion == nil || !*resp.KeyMetadata.MultiRegion {
			return fmt.Errorf("expected MultiRegion=true")
		}

		replResp, err := tc.client.ReplicateKey(tc.ctx, &kms.ReplicateKeyInput{
			KeyId:         aws.String(multiKeyID),
			ReplicaRegion: aws.String("us-west-2"),
		})
		if err != nil {
			return fmt.Errorf("replicate key: %v", err)
		}
		if replResp.ReplicaKeyMetadata == nil {
			return fmt.Errorf("replica key metadata is nil")
		}
		if replResp.ReplicaKeyMetadata.KeyId == nil || *replResp.ReplicaKeyMetadata.KeyId == "" {
			return fmt.Errorf("replica key ID is nil or empty")
		}
		if replResp.ReplicaKeyMetadata.MultiRegion == nil || !*replResp.ReplicaKeyMetadata.MultiRegion {
			return fmt.Errorf("expected replica MultiRegion=true")
		}
		tc.addCleanupKey(*replResp.ReplicaKeyMetadata.KeyId)
		return nil
	}))

	results = append(results, r.RunTest("kms", "UpdatePrimaryRegion", func() error {
		if multiKeyID == "" {
			return fmt.Errorf("multi-region key ID not available")
		}
		_, err := tc.client.UpdatePrimaryRegion(tc.ctx, &kms.UpdatePrimaryRegionInput{
			KeyId:         aws.String(multiKeyID),
			PrimaryRegion: aws.String("us-east-1"),
		})
		if err != nil {
			return fmt.Errorf("update primary region: %v", err)
		}

		descResp, err := tc.client.DescribeKey(tc.ctx, &kms.DescribeKeyInput{KeyId: aws.String(multiKeyID)})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if descResp.KeyMetadata == nil {
			return fmt.Errorf("key metadata is nil")
		}
		if descResp.KeyMetadata.KeyState != types.KeyStateEnabled {
			return fmt.Errorf("expected KeyState=Enabled, got %s", descResp.KeyMetadata.KeyState)
		}
		return nil
	}))

	return results
}
