package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfEdgeTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListKeyGroups_VerifyFields", func() error {
		resp, err := client.ListKeyGroups(ctx, &cloudfront.ListKeyGroupsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.KeyGroupList == nil {
			return fmt.Errorf("key group list is nil")
		}
		if resp.KeyGroupList.MaxItems == nil {
			return fmt.Errorf("maxItems is nil")
		}
		if resp.KeyGroupList.Quantity == nil {
			return fmt.Errorf("quantity is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetInvalidation_NonExistentDist", func() error {
		_, err := client.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
			DistributionId: aws.String("nonexistent-dist-id"),
			Id:             aws.String("any-inv-id"),
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateInvalidation_NonExistentDist", func() error {
		_, err := client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
			DistributionId: aws.String("nonexistent-dist-id"),
			InvalidationBatch: &types.InvalidationBatch{
				CallerReference: aws.String("ref"),
				Paths: &types.Paths{
					Quantity: aws.Int32(1),
					Items:    []string{"/"},
				},
			},
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreatePublicKey_MissingCallerReference_Rejected", func() error {
		_, err := client.CreatePublicKey(ctx, &cloudfront.CreatePublicKeyInput{
			PublicKeyConfig: &types.PublicKeyConfig{
				Name:       aws.String(tc.uniquePrefix("pk-neg")),
				EncodedKey: aws.String(cloudfrontEncodedKeyB64),
				// An empty caller reference passes client-side
				// validation but must be rejected by the server, which
				// models the member as required.
				CallerReference: aws.String(""),
			},
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "PublicKey_RoundTrip", func() error {
		name := tc.uniquePrefix("pk-rt")
		created, err := client.CreatePublicKey(ctx, &cloudfront.CreatePublicKeyInput{
			PublicKeyConfig: &types.PublicKeyConfig{
				Name:            aws.String(name),
				EncodedKey:      aws.String(cloudfrontEncodedKeyB64),
				CallerReference: aws.String(tc.uniquePrefix("pkref")),
			},
		})
		if err != nil {
			return err
		}
		pkID := aws.ToString(created.PublicKey.Id)
		etag := aws.ToString(created.ETag)

		got, err := client.GetPublicKey(ctx, &cloudfront.GetPublicKeyInput{Id: aws.String(pkID)})
		if err != nil {
			return err
		}
		if aws.ToString(got.PublicKey.PublicKeyConfig.Name) != name {
			return fmt.Errorf("public key name mismatch: got %q", aws.ToString(got.PublicKey.PublicKeyConfig.Name))
		}

		cfgResp, err := client.GetPublicKeyConfig(ctx, &cloudfront.GetPublicKeyConfigInput{Id: aws.String(pkID)})
		if err != nil {
			return err
		}
		cfgResp.PublicKeyConfig.Comment = aws.String("round-trip comment")

		updated, err := client.UpdatePublicKey(ctx, &cloudfront.UpdatePublicKeyInput{
			Id:              aws.String(pkID),
			IfMatch:         aws.String(etag),
			PublicKeyConfig: cfgResp.PublicKeyConfig,
		})
		if err != nil {
			return err
		}
		if aws.ToString(updated.ETag) == etag {
			return fmt.Errorf("ETag should change after update")
		}

		listResp, err := client.ListPublicKeys(ctx, &cloudfront.ListPublicKeysInput{MaxItems: aws.Int32(100)})
		if err != nil {
			return err
		}
		found := false
		for _, item := range listResp.PublicKeyList.Items {
			if aws.ToString(item.Id) == pkID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("public key %q not listed", pkID)
		}

		if _, err := client.DeletePublicKey(ctx, &cloudfront.DeletePublicKeyInput{
			Id:      aws.String(pkID),
			IfMatch: aws.String(aws.ToString(updated.ETag)),
		}); err != nil {
			return err
		}
		_, err = client.GetPublicKey(ctx, &cloudfront.GetPublicKeyInput{Id: aws.String(pkID)})
		return AssertErrorContains(err, "NoSuchPublicKey")
	}))

	return results
}
