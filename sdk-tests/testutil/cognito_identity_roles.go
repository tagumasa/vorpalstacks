package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
)

func (r *TestRunner) cognitoIdentityRolesTests(ctx context.Context, client *cognitoidentity.Client, poolID string) []TestResult {
	var results []TestResult
	acct := r.accountID

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles", func() error {
		_, err := client.SetIdentityPoolRoles(ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated":   fmt.Sprintf("arn:aws:iam::%s:role/auth-role", acct),
				"unauthenticated": fmt.Sprintf("arn:aws:iam::%s:role/unauth-role", acct),
			},
		})
		if err != nil {
			return err
		}
		rolesResp, err := client.GetIdentityPoolRoles(ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return fmt.Errorf("GetIdentityPoolRoles after set: %v", err)
		}
		if rolesResp.Roles["authenticated"] != fmt.Sprintf("arn:aws:iam::%s:role/auth-role", acct) {
			return fmt.Errorf("authenticated role not saved")
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
			return fmt.Errorf("roles is nil")
		}
		if resp.Roles["authenticated"] != fmt.Sprintf("arn:aws:iam::%s:role/auth-role", acct) {
			return fmt.Errorf("unexpected authenticated role: %v", resp.Roles["authenticated"])
		}
		if resp.Roles["unauthenticated"] != fmt.Sprintf("arn:aws:iam::%s:role/unauth-role", acct) {
			return fmt.Errorf("unexpected unauthenticated role: %v", resp.Roles["unauthenticated"])
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles_WithMappings", func() error {
		_, err := client.SetIdentityPoolRoles(ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated": fmt.Sprintf("arn:aws:iam::%s:role/auth-role", acct),
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

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles_RuleMappings", func() error {
		name := fmt.Sprintf("test-idpool-rules-%d", poolID[len(poolID)-5:])
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
				"authenticated": fmt.Sprintf("arn:aws:iam::%s:role/auth", acct),
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
								RoleARN:   aws.String(fmt.Sprintf("arn:aws:iam::%s:role/admin", acct)),
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
		client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: aws.String(pid),
		})
		return nil
	}))

	return results
}
