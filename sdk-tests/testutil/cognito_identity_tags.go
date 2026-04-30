package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

func (r *TestRunner) cognitoIdentityTagsTests(ctx context.Context, client *cognitoidentity.Client, poolID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "TagResource", func() error {
		arn := fmt.Sprintf("arn:aws:cognito-identity:%s:000000000000:identitypool/%s", r.region, poolID)
		_, err := client.TagResource(ctx, &cognitoidentity.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags: map[string]string{
				"Environment": "test",
				"Team":        "sdk-tests",
			},
		})
		if err != nil {
			return err
		}
		listResp, err := client.ListTagsForResource(ctx, &cognitoidentity.ListTagsForResourceInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("ListTagsForResource after TagResource: %v", err)
		}
		if listResp.Tags == nil || listResp.Tags["Environment"] != "test" {
			return fmt.Errorf("tag Environment not found after TagResource")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &cognitoidentity.ListTagsForResourceInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:cognito-identity:%s:000000000000:identitypool/%s", r.region, poolID)),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["Environment"] != "test" {
			return fmt.Errorf("expected Environment=test, got %s", resp.Tags["Environment"])
		}
		if resp.Tags["Team"] != "sdk-tests" {
			return fmt.Errorf("expected Team=sdk-tests, got %s", resp.Tags["Team"])
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UntagResource", func() error {
		arn := fmt.Sprintf("arn:aws:cognito-identity:%s:000000000000:identitypool/%s", r.region, poolID)
		_, err := client.UntagResource(ctx, &cognitoidentity.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"Team"},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &cognitoidentity.ListTagsForResourceInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return err
		}
		if _, ok := resp.Tags["Team"]; ok {
			return fmt.Errorf("team tag should have been removed")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "SetPrincipalTagAttributeMap", func() error {
		_, err := client.SetPrincipalTagAttributeMap(ctx, &cognitoidentity.SetPrincipalTagAttributeMapInput{
			IdentityPoolId:       aws.String(poolID),
			IdentityProviderName: aws.String("graph.facebook.com"),
			PrincipalTags: map[string]string{
				"email":    "email",
				"username": "sub",
			},
			UseDefaults: aws.Bool(false),
		})
		if err != nil {
			return err
		}
		resp, err := client.GetPrincipalTagAttributeMap(ctx, &cognitoidentity.GetPrincipalTagAttributeMapInput{
			IdentityPoolId:       aws.String(poolID),
			IdentityProviderName: aws.String("graph.facebook.com"),
		})
		if err != nil {
			return fmt.Errorf("GetPrincipalTagAttributeMap after set: %v", err)
		}
		if resp.PrincipalTags["email"] != "email" {
			return fmt.Errorf("expected email->email mapping")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetPrincipalTagAttributeMap", func() error {
		resp, err := client.GetPrincipalTagAttributeMap(ctx, &cognitoidentity.GetPrincipalTagAttributeMapInput{
			IdentityPoolId:       aws.String(poolID),
			IdentityProviderName: aws.String("graph.facebook.com"),
		})
		if err != nil {
			return err
		}
		if resp.IdentityPoolId == nil || *resp.IdentityPoolId != poolID {
			return fmt.Errorf("expected pool ID %s", poolID)
		}
		if resp.PrincipalTags["email"] != "email" {
			return fmt.Errorf("expected email->email mapping")
		}
		if resp.PrincipalTags["username"] != "sub" {
			return fmt.Errorf("expected username->sub mapping")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetPrincipalTagAttributeMap_Defaults", func() error {
		resp, err := client.GetPrincipalTagAttributeMap(ctx, &cognitoidentity.GetPrincipalTagAttributeMapInput{
			IdentityPoolId:       aws.String(poolID),
			IdentityProviderName: aws.String("accounts.google.com"),
		})
		if err != nil {
			return err
		}
		if resp.UseDefaults == nil || !*resp.UseDefaults {
			return fmt.Errorf("expected UseDefaults=true for non-existent mapping")
		}
		return nil
	}))

	return results
}
