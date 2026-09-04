package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

// cfKeyGroupTests exercises the key-group item operations end to end. The
// family shares its public-key infrastructure, so every test that stores a
// key group first creates a real public key and removes both at the end.
func cfKeyGroupTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_RoundTrip", func() error {
		pkID, pkETag, err := tc.createPublicKey("kg-rt-pk")
		if err != nil {
			return err
		}
		pkDeleted := false
		defer func() {
			if !pkDeleted {
				_, _ = client.DeletePublicKey(ctx, &cloudfront.DeletePublicKeyInput{Id: aws.String(pkID), IfMatch: aws.String(pkETag)})
			}
		}()

		name := tc.uniquePrefix("kg-rt")
		created, err := client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{
				Name:    aws.String(name),
				Items:   []string{pkID},
				Comment: aws.String("round trip"),
			},
		})
		if err != nil {
			return err
		}
		kgID := aws.ToString(created.KeyGroup.Id)
		etag := aws.ToString(created.ETag)
		if kgID == "" || etag == "" {
			return fmt.Errorf("create returned empty id %q or etag %q", kgID, etag)
		}
		if aws.ToString(created.KeyGroup.KeyGroupConfig.Name) != name {
			return fmt.Errorf("created config name mismatch: got %q", aws.ToString(created.KeyGroup.KeyGroupConfig.Name))
		}

		got, err := client.GetKeyGroup(ctx, &cloudfront.GetKeyGroupInput{Id: aws.String(kgID)})
		if err != nil {
			return err
		}
		if aws.ToString(got.KeyGroup.KeyGroupConfig.Name) != name {
			return fmt.Errorf("get config name mismatch: got %q", aws.ToString(got.KeyGroup.KeyGroupConfig.Name))
		}
		if len(got.KeyGroup.KeyGroupConfig.Items) != 1 || got.KeyGroup.KeyGroupConfig.Items[0] != pkID {
			return fmt.Errorf("get items mismatch: got %v", got.KeyGroup.KeyGroupConfig.Items)
		}

		cfgResp, err := client.GetKeyGroupConfig(ctx, &cloudfront.GetKeyGroupConfigInput{Id: aws.String(kgID)})
		if err != nil {
			return err
		}
		if aws.ToString(cfgResp.KeyGroupConfig.Comment) != "round trip" {
			return fmt.Errorf("get-config comment mismatch: got %q", aws.ToString(cfgResp.KeyGroupConfig.Comment))
		}
		if aws.ToString(cfgResp.ETag) != etag {
			return fmt.Errorf("get-config etag mismatch: got %q want %q", aws.ToString(cfgResp.ETag), etag)
		}

		cfgResp.KeyGroupConfig.Comment = aws.String("updated comment")
		updated, err := client.UpdateKeyGroup(ctx, &cloudfront.UpdateKeyGroupInput{
			Id:             aws.String(kgID),
			IfMatch:        aws.String(etag),
			KeyGroupConfig: cfgResp.KeyGroupConfig,
		})
		if err != nil {
			return err
		}
		newETag := aws.ToString(updated.ETag)
		if newETag == etag {
			return fmt.Errorf("ETag should change after update")
		}

		groups, err := paginate(func(next *string) ([]types.KeyGroupSummary, *string, error) {
			resp, lerr := client.ListKeyGroups(ctx, &cloudfront.ListKeyGroupsInput{Marker: next})
			if lerr != nil {
				return nil, nil, lerr
			}
			if resp.KeyGroupList == nil {
				return nil, nil, nil
			}
			return resp.KeyGroupList.Items, resp.KeyGroupList.NextMarker, nil
		})
		if err != nil {
			return err
		}
		found := false
		for _, item := range groups {
			if aws.ToString(item.KeyGroup.Id) == kgID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("key group %q not listed", kgID)
		}

		if _, err := client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{
			Id:      aws.String(kgID),
			IfMatch: aws.String(newETag),
		}); err != nil {
			return err
		}
		_, err = client.GetKeyGroup(ctx, &cloudfront.GetKeyGroupInput{Id: aws.String(kgID)})
		if aerr := AssertErrorContains(err, "NoSuchResource"); aerr != nil {
			return aerr
		}
		if _, err := client.DeletePublicKey(ctx, &cloudfront.DeletePublicKeyInput{
			Id:      aws.String(pkID),
			IfMatch: aws.String(pkETag),
		}); err != nil {
			return err
		}
		pkDeleted = true
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_DuplicateName_Rejected", func() error {
		pkID, pkETag, err := tc.createPublicKey("kg-dup-pk")
		if err != nil {
			return err
		}
		defer func() {
			_, _ = client.DeletePublicKey(ctx, &cloudfront.DeletePublicKeyInput{Id: aws.String(pkID), IfMatch: aws.String(pkETag)})
		}()

		name := tc.uniquePrefix("kg-dup")
		created, err := client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{Name: aws.String(name), Items: []string{pkID}},
		})
		if err != nil {
			return err
		}
		defer func() {
			_, _ = client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{
				Id:      created.KeyGroup.Id,
				IfMatch: created.ETag,
			})
		}()

		_, err = client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{Name: aws.String(name), Items: []string{pkID}},
		})
		return AssertErrorContains(err, "KeyGroupAlreadyExists")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_UnknownID_Rejected", func() error {
		randomID := tc.uniquePrefix("no-such-key-group")
		_, err := client.GetKeyGroup(ctx, &cloudfront.GetKeyGroupInput{Id: aws.String(randomID)})
		if aerr := AssertErrorContains(err, "NoSuchResource"); aerr != nil {
			return aerr
		}
		_, err = client.GetKeyGroupConfig(ctx, &cloudfront.GetKeyGroupConfigInput{Id: aws.String(randomID)})
		if aerr := AssertErrorContains(err, "NoSuchResource"); aerr != nil {
			return aerr
		}
		_, err = client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{Id: aws.String(randomID), IfMatch: aws.String("stub-etag")})
		return AssertErrorContains(err, "NoSuchResource")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_Preconditions_Rejected", func() error {
		pkID, pkETag, err := tc.createPublicKey("kg-pre-pk")
		if err != nil {
			return err
		}
		defer func() {
			_, _ = client.DeletePublicKey(ctx, &cloudfront.DeletePublicKeyInput{Id: aws.String(pkID), IfMatch: aws.String(pkETag)})
		}()

		created, err := client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{Name: aws.String(tc.uniquePrefix("kg-pre")), Items: []string{pkID}},
		})
		if err != nil {
			return err
		}
		kgID := aws.ToString(created.KeyGroup.Id)
		etag := aws.ToString(created.ETag)

		_, err = client.UpdateKeyGroup(ctx, &cloudfront.UpdateKeyGroupInput{
			Id:             aws.String(kgID),
			KeyGroupConfig: created.KeyGroup.KeyGroupConfig,
		})
		if aerr := AssertErrorContains(err, "InvalidIfMatchVersion"); aerr != nil {
			return aerr
		}

		_, err = client.UpdateKeyGroup(ctx, &cloudfront.UpdateKeyGroupInput{
			Id:             aws.String(kgID),
			IfMatch:        aws.String("bogus-etag"),
			KeyGroupConfig: created.KeyGroup.KeyGroupConfig,
		})
		if aerr := AssertErrorContains(err, "PreconditionFailed"); aerr != nil {
			return aerr
		}

		_, err = client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{Id: aws.String(kgID)})
		if aerr := AssertErrorContains(err, "InvalidIfMatchVersion"); aerr != nil {
			return aerr
		}

		_, err = client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{Id: aws.String(kgID), IfMatch: aws.String(etag)})
		return err
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_TooManyPublicKeys_Rejected", func() error {
		items := make([]string, 6)
		for i := range items {
			items[i] = tc.uniquePrefix(fmt.Sprintf("kg-over-pk-%d", i))
		}
		_, err := client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{Name: aws.String(tc.uniquePrefix("kg-over")), Items: items},
		})
		return AssertErrorContains(err, "TooManyPublicKeysInKeyGroup")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_UnknownPublicKey_Rejected", func() error {
		_, err := client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{
				Name:  aws.String(tc.uniquePrefix("kg-unknown-pk")),
				Items: []string{tc.uniquePrefix("no-such-public-key")},
			},
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateDistribution_UnknownTrustedKeyGroup_Rejected", func() error {
		cfg := tc.baseDistributionConfig(tc.uniquePrefix("kg-unknown-kg"), "unknown trusted key group", "example.org")
		cfg.DefaultCacheBehavior.TrustedKeyGroups = &types.TrustedKeyGroups{
			Enabled:  aws.Bool(true),
			Quantity: aws.Int32(1),
			Items:    []string{tc.uniquePrefix("no-such-key-group")},
		}
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributionsByKeyGroup_AttachedDistribution", func() error {
		pkID, pkETag, err := tc.createPublicKey("kg-att-pk")
		if err != nil {
			return err
		}
		defer func() {
			_, _ = client.DeletePublicKey(ctx, &cloudfront.DeletePublicKeyInput{Id: aws.String(pkID), IfMatch: aws.String(pkETag)})
		}()

		created, err := client.CreateKeyGroup(ctx, &cloudfront.CreateKeyGroupInput{
			KeyGroupConfig: &types.KeyGroupConfig{Name: aws.String(tc.uniquePrefix("kg-att")), Items: []string{pkID}},
		})
		if err != nil {
			return err
		}
		kgID := aws.ToString(created.KeyGroup.Id)
		kgETag := aws.ToString(created.ETag)
		keyGroupDeleted := false
		defer func() {
			if !keyGroupDeleted {
				_, _ = client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{Id: aws.String(kgID), IfMatch: aws.String(kgETag)})
			}
		}()

		cfg := tc.baseDistributionConfig(tc.uniquePrefix("kg-att-dist"), "key group attachment", "example.org")
		cfg.DefaultCacheBehavior.TrustedKeyGroups = &types.TrustedKeyGroups{
			Enabled:  aws.Bool(true),
			Quantity: aws.Int32(1),
			Items:    []string{kgID},
		}
		distResp, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{DistributionConfig: cfg})
		if err != nil {
			return err
		}
		distID := aws.ToString(distResp.Distribution.Id)
		distETag := aws.ToString(distResp.ETag)
		defer func() {
			_ = tc.disableAndDeleteDistribution(distID, distETag)
		}()

		ids, err := paginate(func(next *string) ([]string, *string, error) {
			resp, lerr := client.ListDistributionsByKeyGroup(ctx, &cloudfront.ListDistributionsByKeyGroupInput{
				KeyGroupId: aws.String(kgID),
				Marker:     next,
			})
			if lerr != nil {
				return nil, nil, lerr
			}
			if resp.DistributionIdList == nil {
				return nil, nil, nil
			}
			return resp.DistributionIdList.Items, resp.DistributionIdList.NextMarker, nil
		})
		if err != nil {
			return err
		}
		found := false
		for _, id := range ids {
			if id == distID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("distribution %q not listed under key group %q", distID, kgID)
		}

		_, err = client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{Id: aws.String(kgID), IfMatch: aws.String(kgETag)})
		if aerr := AssertErrorContains(err, "ResourceInUse"); aerr != nil {
			return aerr
		}

		if err := tc.disableAndDeleteDistribution(distID, distETag); err != nil {
			return err
		}
		if _, err := client.DeleteKeyGroup(ctx, &cloudfront.DeleteKeyGroupInput{Id: aws.String(kgID), IfMatch: aws.String(kgETag)}); err != nil {
			return err
		}
		keyGroupDeleted = true
		_, err = client.GetKeyGroup(ctx, &cloudfront.GetKeyGroupInput{Id: aws.String(kgID)})
		return AssertErrorContains(err, "NoSuchResource")
	}))

	return results
}
