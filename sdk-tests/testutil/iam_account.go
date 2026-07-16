package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamAccountTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "GetAccountSummary", func() error {
		resp, err := tc.client.GetAccountSummary(tc.ctx, &iam.GetAccountSummaryInput{})
		if err != nil {
			return err
		}
		if resp.SummaryMap == nil {
			return fmt.Errorf("summary map is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetAccountAuthorizationDetails", func() error {
		resp, err := tc.client.GetAccountAuthorizationDetails(tc.ctx, &iam.GetAccountAuthorizationDetailsInput{
			Filter: []types.EntityType{types.EntityTypeUser},
		})
		if err != nil {
			return err
		}
		if resp.UserDetailList == nil {
			return fmt.Errorf("user detail list is nil")
		}
		return nil
	}))

	// Account aliases
	results = append(results, r.RunTest("iam", "CreateAccountAlias", func() error {
		_, err := tc.client.CreateAccountAlias(tc.ctx, &iam.CreateAccountAliasInput{
			AccountAlias: aws.String(tc.accountAlias),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAccountAliases(tc.ctx, &iam.ListAccountAliasesInput{})
		if err != nil {
			return fmt.Errorf("ListAccountAliases after create: %w", err)
		}
		found := false
		for _, a := range resp.AccountAliases {
			if a == tc.accountAlias {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("alias %s not found after CreateAccountAlias", tc.accountAlias)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListAccountAliases", func() error {
		resp, err := tc.client.ListAccountAliases(tc.ctx, &iam.ListAccountAliasesInput{})
		if err != nil {
			return err
		}
		found := false
		for _, alias := range resp.AccountAliases {
			if alias == tc.accountAlias {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("account alias %s not found", tc.accountAlias)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteAccountAlias", func() error {
		_, err := tc.client.DeleteAccountAlias(tc.ctx, &iam.DeleteAccountAliasInput{
			AccountAlias: aws.String(tc.accountAlias),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAccountAliases(tc.ctx, &iam.ListAccountAliasesInput{})
		if err != nil {
			return fmt.Errorf("ListAccountAliases after delete: %w", err)
		}
		for _, a := range resp.AccountAliases {
			if a == tc.accountAlias {
				return fmt.Errorf("alias %s still exists after DeleteAccountAlias", tc.accountAlias)
			}
		}
		return nil
	}))

	// Password policy
	results = append(results, r.RunTest("iam", "UpdateAccountPasswordPolicy", func() error {
		_, err := tc.client.UpdateAccountPasswordPolicy(tc.ctx, &iam.UpdateAccountPasswordPolicyInput{
			MinimumPasswordLength:      aws.Int32(12),
			RequireUppercaseCharacters: true,
			RequireLowercaseCharacters: true,
			RequireNumbers:             true,
			RequireSymbols:             true,
			AllowUsersToChangePassword: true,
			MaxPasswordAge:             aws.Int32(90),
			PasswordReusePrevention:    aws.Int32(5),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccountPasswordPolicy(tc.ctx, &iam.GetAccountPasswordPolicyInput{})
		if err != nil {
			return fmt.Errorf("GetAccountPasswordPolicy after update: %w", err)
		}
		if aws.ToInt32(resp.PasswordPolicy.MinimumPasswordLength) != 12 {
			return fmt.Errorf("minimum password length: got %d, want 12", aws.ToInt32(resp.PasswordPolicy.MinimumPasswordLength))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetAccountPasswordPolicy", func() error {
		resp, err := tc.client.GetAccountPasswordPolicy(tc.ctx, &iam.GetAccountPasswordPolicyInput{})
		if err != nil {
			return err
		}
		if resp.PasswordPolicy == nil {
			return fmt.Errorf("password policy is nil")
		}
		pp := resp.PasswordPolicy
		if aws.ToInt32(pp.MinimumPasswordLength) != 12 {
			return fmt.Errorf("minimum password length mismatch: got %d, want 12", aws.ToInt32(pp.MinimumPasswordLength))
		}
		if !pp.RequireUppercaseCharacters {
			return fmt.Errorf("require uppercase should be true")
		}
		if !pp.RequireLowercaseCharacters {
			return fmt.Errorf("require lowercase should be true")
		}
		if !pp.RequireNumbers {
			return fmt.Errorf("require numbers should be true")
		}
		if !pp.RequireSymbols {
			return fmt.Errorf("require symbols should be true")
		}
		if !pp.AllowUsersToChangePassword {
			return fmt.Errorf("allow users to change password should be true")
		}
		if aws.ToInt32(pp.MaxPasswordAge) != 90 {
			return fmt.Errorf("max password age mismatch: got %d, want 90", aws.ToInt32(pp.MaxPasswordAge))
		}
		if aws.ToInt32(pp.PasswordReusePrevention) != 5 {
			return fmt.Errorf("password reuse prevention mismatch: got %d, want 5", aws.ToInt32(pp.PasswordReusePrevention))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteAccountPasswordPolicy", func() error {
		_, err := tc.client.DeleteAccountPasswordPolicy(tc.ctx, &iam.DeleteAccountPasswordPolicyInput{})
		if err != nil {
			return err
		}
		_, err = tc.client.GetAccountPasswordPolicy(tc.ctx, &iam.GetAccountPasswordPolicyInput{})
		if err == nil {
			return fmt.Errorf("GetAccountPasswordPolicy should fail after DeleteAccountPasswordPolicy")
		}
		return nil
	}))

	return results
}
