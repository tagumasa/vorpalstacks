package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamInstanceProfileTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "CreateInstanceProfile", func() error {
		resp, err := tc.client.CreateInstanceProfile(tc.ctx, &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return err
		}
		if resp.InstanceProfile == nil {
			return fmt.Errorf("instance profile is nil")
		}
		if aws.ToString(resp.InstanceProfile.InstanceProfileName) != tc.profile {
			return fmt.Errorf("instance profile name mismatch")
		}
		if aws.ToString(resp.InstanceProfile.Arn) == "" {
			return fmt.Errorf("instance profile arn is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetInstanceProfile", func() error {
		resp, err := tc.client.GetInstanceProfile(tc.ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return err
		}
		if resp.InstanceProfile == nil {
			return fmt.Errorf("instance profile is nil")
		}
		if aws.ToString(resp.InstanceProfile.InstanceProfileName) != tc.profile {
			return fmt.Errorf("instance profile name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListInstanceProfiles", func() error {
		var found bool
		var marker *string
		for {
			resp, err := tc.client.ListInstanceProfiles(tc.ctx, &iam.ListInstanceProfilesInput{Marker: marker})
			if err != nil {
				return err
			}
			for _, ip := range resp.InstanceProfiles {
				if aws.ToString(ip.InstanceProfileName) == tc.profile {
					found = true
					break
				}
			}
			if found || !resp.IsTruncated || resp.Marker == nil {
				break
			}
			marker = resp.Marker
		}
		if !found {
			return fmt.Errorf("instance profile %s not found", tc.profile)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AddRoleToInstanceProfile", func() error {
		_, err := tc.client.AddRoleToInstanceProfile(tc.ctx, &iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
			RoleName:            aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetInstanceProfile(tc.ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return fmt.Errorf("GetInstanceProfile after AddRole: %w", err)
		}
		found := false
		for _, r := range resp.InstanceProfile.Roles {
			if aws.ToString(r.RoleName) == tc.role {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("role %s not found in instance profile after AddRoleToInstanceProfile", tc.role)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListInstanceProfilesForRole", func() error {
		resp, err := tc.client.ListInstanceProfilesForRole(tc.ctx, &iam.ListInstanceProfilesForRoleInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		found := false
		for _, ip := range resp.InstanceProfiles {
			if aws.ToString(ip.InstanceProfileName) == tc.profile {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("instance profile %s not found for role %s", tc.profile, tc.role)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "RemoveRoleFromInstanceProfile", func() error {
		_, err := tc.client.RemoveRoleFromInstanceProfile(tc.ctx, &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
			RoleName:            aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetInstanceProfile(tc.ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return fmt.Errorf("GetInstanceProfile after RemoveRole: %w", err)
		}
		for _, r := range resp.InstanceProfile.Roles {
			if aws.ToString(r.RoleName) == tc.role {
				return fmt.Errorf("role %s still in instance profile after RemoveRoleFromInstanceProfile", tc.role)
			}
		}
		return nil
	}))

	// Tags
	results = append(results, r.RunTest("iam", "TagInstanceProfile", func() error {
		_, err := tc.client.TagInstanceProfile(tc.ctx, &iam.TagInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListInstanceProfileTags(tc.ctx, &iam.ListInstanceProfileTagsInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return fmt.Errorf("ListInstanceProfileTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagInstanceProfile")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListInstanceProfileTags", func() error {
		resp, err := tc.client.ListInstanceProfileTags(tc.ctx, &iam.ListInstanceProfileTagsInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagInstanceProfile", func() error {
		_, err := tc.client.UntagInstanceProfile(tc.ctx, &iam.UntagInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
			TagKeys:             []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListInstanceProfileTags(tc.ctx, &iam.ListInstanceProfileTagsInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		if err != nil {
			return err
		}
		if iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment tag should be removed")
		}
		return nil
	}))

	return results
}
