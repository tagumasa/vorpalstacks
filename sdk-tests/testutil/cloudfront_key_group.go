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

	createPublicKey := func(tag string) (id, etag string, err error) {
		resp, cerr := client.CreatePublicKey(ctx, &cloudfront.CreatePublicKeyInput{
			PublicKeyConfig: &types.PublicKeyConfig{
				Name:            aws.String(tc.uniquePrefix(tag)),
				EncodedKey:      aws.String(cloudfrontEncodedKeyB64),
				CallerReference: aws.String(tc.uniqueCallerRef(tag + "-ref")),
			},
		})
		if cerr != nil {
			return "", "", cerr
		}
		return aws.ToString(resp.PublicKey.Id), aws.ToString(resp.ETag), nil
	}

	results = append(results, tc.runner.RunTest("cloudfront", "KeyGroup_RoundTrip", func() error {
		pkID, pkETag, err := createPublicKey("kg-rt-pk")
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

		found := false
		marker := ""
		for {
			input := &cloudfront.ListKeyGroupsInput{}
			if marker != "" {
				input.Marker = aws.String(marker)
			}
			listResp, err := client.ListKeyGroups(ctx, input)
			if err != nil {
				return err
			}
			for _, item := range listResp.KeyGroupList.Items {
				if aws.ToString(item.KeyGroup.Id) == kgID {
					found = true
				}
			}
			if found || listResp.KeyGroupList.NextMarker == nil || aws.ToString(listResp.KeyGroupList.NextMarker) == "" {
				break
			}
			marker = aws.ToString(listResp.KeyGroupList.NextMarker)
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
		pkID, pkETag, err := createPublicKey("kg-dup-pk")
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
		pkID, pkETag, err := createPublicKey("kg-pre-pk")
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
		originID := tc.uniquePrefix("origin")
		_, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{
			DistributionConfig: &types.DistributionConfig{
				CallerReference: aws.String(tc.uniqueCallerRef("kg-unknown-kg")),
				Enabled:         aws.Bool(true),
				Comment:         aws.String("unknown trusted key group"),
				Origins: &types.Origins{
					Quantity: aws.Int32(1),
					Items: []types.Origin{
						{
							Id:         aws.String(originID),
							DomainName: aws.String("example.org"),
							CustomOriginConfig: &types.CustomOriginConfig{
								HTTPPort:             aws.Int32(80),
								HTTPSPort:            aws.Int32(443),
								OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
							},
						},
					},
				},
				DefaultCacheBehavior: &types.DefaultCacheBehavior{
					TargetOriginId:       aws.String(originID),
					ViewerProtocolPolicy: types.ViewerProtocolPolicyAllowAll,
					TrustedKeyGroups: &types.TrustedKeyGroups{
						Enabled:  aws.Bool(true),
						Quantity: aws.Int32(1),
						Items:    []string{tc.uniquePrefix("no-such-key-group")},
					},
					ForwardedValues: &types.ForwardedValues{
						QueryString: aws.Bool(false),
						Cookies:     &types.CookiePreference{Forward: types.ItemSelectionNone},
					},
				},
				ViewerCertificate: &types.ViewerCertificate{CloudFrontDefaultCertificate: aws.Bool(true)},
				Restrictions: &types.Restrictions{
					GeoRestriction: &types.GeoRestriction{RestrictionType: types.GeoRestrictionTypeNone, Quantity: aws.Int32(0)},
				},
			},
		})
		return AssertErrorContains(err, "InvalidArgument")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "ListDistributionsByKeyGroup_AttachedDistribution", func() error {
		pkID, pkETag, err := createPublicKey("kg-att-pk")
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

		originID := tc.uniquePrefix("origin")
		distResp, err := client.CreateDistribution(ctx, &cloudfront.CreateDistributionInput{
			DistributionConfig: &types.DistributionConfig{
				CallerReference: aws.String(tc.uniqueCallerRef("kg-att-dist")),
				Enabled:         aws.Bool(true),
				Comment:         aws.String("key group attachment"),
				Origins: &types.Origins{
					Quantity: aws.Int32(1),
					Items: []types.Origin{
						{
							Id:         aws.String(originID),
							DomainName: aws.String("example.org"),
							CustomOriginConfig: &types.CustomOriginConfig{
								HTTPPort:             aws.Int32(80),
								HTTPSPort:            aws.Int32(443),
								OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
							},
						},
					},
				},
				DefaultCacheBehavior: &types.DefaultCacheBehavior{
					TargetOriginId:       aws.String(originID),
					ViewerProtocolPolicy: types.ViewerProtocolPolicyAllowAll,
					TrustedKeyGroups: &types.TrustedKeyGroups{
						Enabled:  aws.Bool(true),
						Quantity: aws.Int32(1),
						Items:    []string{kgID},
					},
					ForwardedValues: &types.ForwardedValues{
						QueryString: aws.Bool(false),
						Cookies:     &types.CookiePreference{Forward: types.ItemSelectionNone},
					},
				},
				ViewerCertificate: &types.ViewerCertificate{CloudFrontDefaultCertificate: aws.Bool(true)},
				Restrictions: &types.Restrictions{
					GeoRestriction: &types.GeoRestriction{RestrictionType: types.GeoRestrictionTypeNone, Quantity: aws.Int32(0)},
				},
			},
		})
		if err != nil {
			return err
		}
		distID := aws.ToString(distResp.Distribution.Id)
		distETag := aws.ToString(distResp.ETag)
		defer func() {
			_ = tc.disableAndDeleteDistribution(distID, distETag)
		}()

		found := false
		marker := ""
		for {
			input := &cloudfront.ListDistributionsByKeyGroupInput{KeyGroupId: aws.String(kgID)}
			if marker != "" {
				input.Marker = aws.String(marker)
			}
			listResp, err := client.ListDistributionsByKeyGroup(ctx, input)
			if err != nil {
				return err
			}
			for _, id := range listResp.DistributionIdList.Items {
				if id == distID {
					found = true
				}
			}
			if found || listResp.DistributionIdList.NextMarker == nil || aws.ToString(listResp.DistributionIdList.NextMarker) == "" {
				break
			}
			marker = aws.ToString(listResp.DistributionIdList.NextMarker)
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
