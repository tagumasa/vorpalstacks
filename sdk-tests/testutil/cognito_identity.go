package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCognitoIdentityTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cognito-identity",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cognitoidentity.NewFromConfig(cfg)
	ctx := context.Background()

	poolName := fmt.Sprintf("test-idpool-%d", time.Now().UnixNano())

	var poolID string
	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool", func() error {
		resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(poolName),
			AllowUnauthenticatedIdentities: true,
		})
		if err != nil {
			return err
		}
		if resp.IdentityPoolId == nil {
			return fmt.Errorf("IdentityPoolId is nil")
		}
		poolID = *resp.IdentityPoolId
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentityPool", func() error {
		resp, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolName != poolName {
			return fmt.Errorf("expected pool name %s, got %s", poolName, *resp.IdentityPoolName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool_WithOptions", func() error {
		name := fmt.Sprintf("test-idpool-opts-%d", time.Now().UnixNano())
		resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(name),
			AllowUnauthenticatedIdentities: false,
			AllowClassicFlow:               aws.Bool(true),
			DeveloperProviderName:          aws.String("my-dev-provider"),
			OpenIdConnectProviderARNs:      []string{"arn:aws:iam::123456789012:oidc-provider/example.com"},
			SamlProviderARNs:               []string{"arn:aws:iam::123456789012:saml-provider/example.com"},
			CognitoIdentityProviders: []types.CognitoIdentityProvider{
				{
					ProviderName:         aws.String("cognito-idp.us-east-1.amazonaws.com/us-east-1_xxxxx"),
					ClientId:             aws.String("abc123"),
					ServerSideTokenCheck: aws.Bool(true),
				},
			},
			SupportedLoginProviders: map[string]string{
				"graph.facebook.com": "1234567890",
			},
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolName != name {
			return fmt.Errorf("expected name %s, got %s", name, *resp.IdentityPoolName)
		}
		if resp.AllowUnauthenticatedIdentities {
			return fmt.Errorf("expected AllowUnauthenticatedIdentities false")
		}
		if resp.AllowClassicFlow == nil || !*resp.AllowClassicFlow {
			return fmt.Errorf("expected AllowClassicFlow true")
		}
		if resp.DeveloperProviderName == nil || *resp.DeveloperProviderName != "my-dev-provider" {
			return fmt.Errorf("expected DeveloperProviderName my-dev-provider")
		}
		if len(resp.OpenIdConnectProviderARNs) != 1 {
			return fmt.Errorf("expected 1 OpenIdConnectProviderARN")
		}
		if len(resp.SamlProviderARNs) != 1 {
			return fmt.Errorf("expected 1 SamlProviderARN")
		}
		if len(resp.CognitoIdentityProviders) != 1 {
			return fmt.Errorf("expected 1 CognitoIdentityProvider")
		}
		if len(resp.SupportedLoginProviders) != 1 {
			return fmt.Errorf("expected 1 SupportedLoginProvider")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListIdentityPools", func() error {
		resp, err := client.ListIdentityPools(ctx, &cognitoidentity.ListIdentityPoolsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if len(resp.IdentityPools) < 1 {
			return fmt.Errorf("expected at least 1 pool, got %d", len(resp.IdentityPools))
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UpdateIdentityPool", func() error {
		newName := poolName + "-updated"
		_, err := client.UpdateIdentityPool(ctx, &cognitoidentity.UpdateIdentityPoolInput{
			IdentityPoolId:                 aws.String(poolID),
			IdentityPoolName:               aws.String(newName),
			AllowUnauthenticatedIdentities: false,
		})
		if err != nil {
			return err
		}
		resp, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolName != newName {
			return fmt.Errorf("expected updated name %s, got %s", newName, *resp.IdentityPoolName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetId", func() error {
		resp, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		return nil
	}))

	var identityID string
	results = append(results, r.RunTest("cognito-identity", "GetId_WithLogins", func() error {
		resp, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"graph.facebook.com": "test-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		identityID = *resp.IdentityId
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentity", func() error {
		resp, err := client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != identityID {
			return fmt.Errorf("expected identity ID %s, got %s", identityID, *resp.IdentityId)
		}
		found := false
		for _, l := range resp.Logins {
			if l == "graph.facebook.com" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected graph.facebook.com in Logins")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity", func() error {
		resp, err := client.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: aws.String(identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != identityID {
			return fmt.Errorf("expected identity ID %s, got %s", identityID, *resp.IdentityId)
		}
		if resp.Credentials == nil {
			return fmt.Errorf("Credentials is nil")
		}
		if resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("AccessKeyId is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles", func() error {
		_, err := client.SetIdentityPoolRoles(ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated":   "arn:aws:iam::123456789012:role/auth-role",
				"unauthenticated": "arn:aws:iam::123456789012:role/unauth-role",
			},
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetIdentityPoolRoles", func() error {
		resp, err := client.GetIdentityPoolRoles(ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.Roles == nil {
			return fmt.Errorf("Roles is nil")
		}
		if resp.Roles["authenticated"] != "arn:aws:iam::123456789012:role/auth-role" {
			return fmt.Errorf("unexpected authenticated role: %v", resp.Roles["authenticated"])
		}
		if resp.Roles["unauthenticated"] != "arn:aws:iam::123456789012:role/unauth-role" {
			return fmt.Errorf("unexpected unauthenticated role: %v", resp.Roles["unauthenticated"])
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles_WithMappings", func() error {
		_, err := client.SetIdentityPoolRoles(ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated": "arn:aws:iam::123456789012:role/auth-role",
			},
			RoleMappings: map[string]types.RoleMapping{
				"graph.facebook.com": {
					Type:                    types.RoleMappingTypeToken,
					AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeAuthenticatedRole,
				},
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.GetIdentityPoolRoles(ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.RoleMappings == nil {
			return fmt.Errorf("RoleMappings is nil")
		}
		if m, ok := resp.RoleMappings["graph.facebook.com"]; !ok {
			return fmt.Errorf("expected graph.facebook.com in RoleMappings")
		} else if m.Type != types.RoleMappingTypeToken {
			return fmt.Errorf("expected Token type, got %s", m.Type)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "TagResource", func() error {
		_, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		_, err = client.TagResource(ctx, &cognitoidentity.TagResourceInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:cognito-identity:%s:000000000000:identitypool/%s", r.region, poolID)),
			Tags: map[string]string{
				"Environment": "test",
				"Team":        "sdk-tests",
			},
		})
		if err != nil {
			return err
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
			return fmt.Errorf("Tags is nil")
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
			return fmt.Errorf("Team tag should have been removed")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListIdentities", func() error {
		resp, err := client.ListIdentities(ctx, &cognitoidentity.ListIdentitiesInput{
			IdentityPoolId: aws.String(poolID),
			MaxResults:     aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolId != poolID {
			return fmt.Errorf("expected pool ID %s, got %s", poolID, *resp.IdentityPoolId)
		}
		if len(resp.Identities) < 1 {
			return fmt.Errorf("expected at least 1 identity, got %d", len(resp.Identities))
		}
		return nil
	}))

	var secondIdentityID string
	results = append(results, r.RunTest("cognito-identity", "GetId_SecondIdentity", func() error {
		resp, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"accounts.google.com": "google-token",
			},
		})
		if err != nil {
			return err
		}
		secondIdentityID = *resp.IdentityId
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentities", func() error {
		resp, err := client.DeleteIdentities(ctx, &cognitoidentity.DeleteIdentitiesInput{
			IdentityIdsToDelete: []string{secondIdentityID},
		})
		if err != nil {
			return err
		}
		if len(resp.UnprocessedIdentityIds) > 0 {
			return fmt.Errorf("unexpected unprocessed identities: %v", resp.UnprocessedIdentityIds)
		}
		_, err = client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(secondIdentityID),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken", func() error {
		resp, err := client.GetOpenIdToken(ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != identityID {
			return fmt.Errorf("expected identity ID %s, got %s", identityID, *resp.IdentityId)
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("Token is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken_WithLogins", func() error {
		resp, err := client.GetOpenIdToken(ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(identityID),
			Logins: map[string]string{
				"graph.facebook.com": "new-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("Token is nil or empty")
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
		if *resp.IdentityPoolId != poolID {
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

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdTokenForDeveloperIdentity", func() error {
		resp, err := client.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"my-dev-provider": "dev-user-1",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("Token is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "LookupDeveloperIdentity", func() error {
		resp, err := client.LookupDeveloperIdentity(ctx, &cognitoidentity.LookupDeveloperIdentityInput{
			IdentityPoolId:          aws.String(poolID),
			DeveloperUserIdentifier: aws.String("dev-user-1"),
		})
		if err != nil {
			return err
		}
		if len(resp.DeveloperUserIdentifierList) < 1 {
			return fmt.Errorf("expected at least 1 developer user identifier")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdTokenForDeveloperIdentity_Reuse", func() error {
		resp, err := client.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"my-dev-provider": "dev-user-1",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "MergeDeveloperIdentities", func() error {
		_, err := client.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"my-dev-provider": "dev-user-2",
			},
		})
		if err != nil {
			return err
		}

		resp, err := client.MergeDeveloperIdentities(ctx, &cognitoidentity.MergeDeveloperIdentitiesInput{
			SourceUserIdentifier:      aws.String("dev-user-1"),
			DestinationUserIdentifier: aws.String("dev-user-2"),
			DeveloperProviderName:     aws.String("my-dev-provider"),
			IdentityPoolId:            aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty after merge")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UnlinkDeveloperIdentity", func() error {
		_, err := client.UnlinkDeveloperIdentity(ctx, &cognitoidentity.UnlinkDeveloperIdentityInput{
			IdentityPoolId:          aws.String(poolID),
			IdentityId:              aws.String(identityID),
			DeveloperProviderName:   aws.String("my-dev-provider"),
			DeveloperUserIdentifier: aws.String("dev-user-1"),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UnlinkIdentity", func() error {
		_, err := client.UnlinkIdentity(ctx, &cognitoidentity.UnlinkIdentityInput{
			IdentityId: aws.String(identityID),
			Logins: map[string]string{
				"graph.facebook.com": "token",
			},
			LoginsToRemove: []string{"graph.facebook.com"},
		})
		if err != nil {
			return err
		}
		resp, err := client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(identityID),
		})
		if err != nil {
			return err
		}
		for _, l := range resp.Logins {
			if l == "graph.facebook.com" {
				return fmt.Errorf("graph.facebook.com should have been unlinked")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentityPool", func() error {
		_, err := client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		_, err = client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted pool")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentityPool_NonExistent", func() error {
		_, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent pool")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentityPool_NonExistent", func() error {
		_, err := client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent pool")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentity_NonExistent", func() error {
		_, err := client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetId_NonExistentPool", func() error {
		_, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent pool")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool_WithTags", func() error {
		name := fmt.Sprintf("test-idpool-tags-%d", time.Now().UnixNano())
		resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(name),
			AllowUnauthenticatedIdentities: true,
			IdentityPoolTags: map[string]string{
				"Env":  "production",
				"Cost": "high",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityPoolTags == nil {
			return fmt.Errorf("IdentityPoolTags is nil in CreateIdentityPool response")
		}
		if resp.IdentityPoolTags["Env"] != "production" {
			return fmt.Errorf("expected Env=production")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles_RuleMappings", func() error {
		name := fmt.Sprintf("test-idpool-rules-%d", time.Now().UnixNano())
		createResp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(name),
			AllowUnauthenticatedIdentities: true,
		})
		if err != nil {
			return err
		}
		pid := *createResp.IdentityPoolId

		_, err = client.SetIdentityPoolRoles(ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(pid),
			Roles: map[string]string{
				"authenticated": "arn:aws:iam::123456789012:role/auth",
			},
			RoleMappings: map[string]types.RoleMapping{
				"graph.facebook.com": {
					Type:                    types.RoleMappingTypeRules,
					AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeDeny,
					RulesConfiguration: &types.RulesConfigurationType{
						Rules: []types.MappingRule{
							{
								Claim:     aws.String("isAdmin"),
								MatchType: types.MappingRuleMatchTypeEquals,
								Value:     aws.String("true"),
								RoleARN:   aws.String("arn:aws:iam::123456789012:role/admin"),
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		resp, err := client.GetIdentityPoolRoles(ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(pid),
		})
		if err != nil {
			return err
		}
		if resp.RoleMappings == nil {
			return fmt.Errorf("RoleMappings is nil")
		}
		m := resp.RoleMappings["graph.facebook.com"]
		if m.Type != types.RoleMappingTypeRules {
			return fmt.Errorf("expected Rules type")
		}
		if m.RulesConfiguration == nil || len(m.RulesConfiguration.Rules) != 1 {
			return fmt.Errorf("expected 1 rule")
		}
		rule := m.RulesConfiguration.Rules[0]
		if *rule.Claim != "isAdmin" {
			return fmt.Errorf("expected claim isAdmin")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity_NonExistent", func() error {
		_, err := client.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken_NonExistent", func() error {
		_, err := client.GetOpenIdToken(ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentities_NonExistent", func() error {
		resp, err := client.DeleteIdentities(ctx, &cognitoidentity.DeleteIdentitiesInput{
			IdentityIdsToDelete: []string{"00000000-0000-0000-0000-000000000000"},
		})
		if err != nil {
			return err
		}
		if len(resp.UnprocessedIdentityIds) != 1 {
			return fmt.Errorf("expected 1 unprocessed identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetIdentityPoolRoles_NonExistent", func() error {
		_, err := client.GetIdentityPoolRoles(ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent pool")
		}
		return nil
	}))

	return results
}
