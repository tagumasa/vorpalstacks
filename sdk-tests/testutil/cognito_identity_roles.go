package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
)

func (r *TestRunner) cognitoIdentityRolesTests(tc *cognitoIdentityContext) []TestResult {
	var results []TestResult
	acct := r.accountID

	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles", func() error {
		_, err := tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
			Roles: map[string]string{
				"authenticated":   fmt.Sprintf("arn:aws:iam::%s:role/auth-role", acct),
				"unauthenticated": fmt.Sprintf("arn:aws:iam::%s:role/unauth-role", acct),
			},
		})
		if err != nil {
			return err
		}
		rolesResp, err := tc.client.GetIdentityPoolRoles(tc.ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
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
		resp, err := tc.client.GetIdentityPoolRoles(tc.ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
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
		_, err := tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
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
		resp, err := tc.client.GetIdentityPoolRoles(tc.ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
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
		pid, cleanupPid, err := tc.createIdPool(fmt.Sprintf("test-idpool-rules-%s", tc.poolID[len(tc.poolID)-5:]))
		if err != nil {
			return err
		}
		defer cleanupPid()

		_, err = tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
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

		resp, err := tc.client.GetIdentityPoolRoles(tc.ctx, &cognitoidentity.GetIdentityPoolRolesInput{
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

	// The RoleMappings key shape is IdentityProviderName: a pattern-less
	// Smithy @length(1, 128) trait counted in Unicode characters, so a
	// 100-character CJK key (300 bytes) must be accepted.
	results = append(results, r.RunTest("cognito-identity", "SetIdentityPoolRoles_MultibyteMappingKeyAccepted", func() error {
		cjkKey := strings.Repeat("\u65e5", 100)
		baseRoles := map[string]string{
			"authenticated":   fmt.Sprintf("arn:aws:iam::%s:role/auth-role", acct),
			"unauthenticated": fmt.Sprintf("arn:aws:iam::%s:role/unauth-role", acct),
		}
		_, err := tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
			Roles:          baseRoles,
			RoleMappings: map[string]types.RoleMapping{
				cjkKey: {
					Type:                    types.RoleMappingTypeToken,
					AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeAuthenticatedRole,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("set with multibyte mapping key: %v", err)
		}
		resp, err := tc.client.GetIdentityPoolRoles(tc.ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
		})
		if err != nil {
			return err
		}
		if _, ok := resp.RoleMappings[cjkKey]; !ok {
			return fmt.Errorf("multibyte role-mapping key not persisted")
		}
		// Restore the plain two-role state for later readers.
		_, err = tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(tc.poolID),
			Roles:          baseRoles,
		})
		return err
	}))

	return results
}
