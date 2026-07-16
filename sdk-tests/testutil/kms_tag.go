package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSTagTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "TagResource", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.TagResource(tc.ctx, &kms.TagResourceInput{
			KeyId: aws.String(tc.keyID),
			Tags: []types.Tag{
				{TagKey: aws.String("Environment"), TagValue: aws.String("test")},
				{TagKey: aws.String("Project"), TagValue: aws.String("sdk-tests")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "ListResourceTags", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.ListResourceTags(tc.ctx, &kms.ListResourceTagsInput{
			KeyId: aws.String(tc.keyID),
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
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.UntagResource(tc.ctx, &kms.UntagResourceInput{
			KeyId:   aws.String(tc.keyID),
			TagKeys: []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "TagResource_ByAlias", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		tagAlias := fmt.Sprintf("alias/tag-test-%d", len(tc.keyID))
		_, err := tc.client.CreateAlias(tc.ctx, &kms.CreateAliasInput{
			AliasName:   aws.String(tagAlias),
			TargetKeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return fmt.Errorf("create alias: %v", err)
		}
		tc.addCleanupAlias(tagAlias)

		_, err = tc.client.TagResource(tc.ctx, &kms.TagResourceInput{
			KeyId: aws.String(tagAlias),
			Tags: []types.Tag{
				{TagKey: aws.String("AliasTag"), TagValue: aws.String("test-value")},
			},
		})
		if err != nil {
			return fmt.Errorf("tag by alias: %v", err)
		}

		tagResp, err := tc.client.ListResourceTags(tc.ctx, &kms.ListResourceTagsInput{KeyId: aws.String(tagAlias)})
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
		tc.client.UntagResource(tc.ctx, &kms.UntagResourceInput{
			KeyId:   aws.String(tc.keyID),
			TagKeys: []string{"AliasTag"},
		})
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListResourceTags_ContentVerify", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		testTags := []types.Tag{
			{TagKey: aws.String("Env"), TagValue: aws.String("prod")},
			{TagKey: aws.String("Team"), TagValue: aws.String("backend")},
		}
		_, err := tc.client.TagResource(tc.ctx, &kms.TagResourceInput{
			KeyId: aws.String(tc.keyID),
			Tags:  testTags,
		})
		if err != nil {
			return fmt.Errorf("tag resource: %v", err)
		}
		resp, err := tc.client.ListResourceTags(tc.ctx, &kms.ListResourceTagsInput{KeyId: aws.String(tc.keyID)})
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
		tc.client.UntagResource(tc.ctx, &kms.UntagResourceInput{
			KeyId:   aws.String(tc.keyID),
			TagKeys: []string{"Env", "Team"},
		})
		return nil
	}))

	return results
}
