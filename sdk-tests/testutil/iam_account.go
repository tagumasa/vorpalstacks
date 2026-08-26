package testutil

import (
	"fmt"
	"time"

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
		if resp.SummaryMap["UsersQuota"] != 5000 {
			return fmt.Errorf("UsersQuota: got %d, want 5000", resp.SummaryMap["UsersQuota"])
		}
		if resp.SummaryMap["InstanceProfilesQuota"] != 1000 {
			return fmt.Errorf("InstanceProfilesQuota: got %d, want 1000", resp.SummaryMap["InstanceProfilesQuota"])
		}
		if resp.SummaryMap["GroupsQuota"] != 300 {
			return fmt.Errorf("GroupsQuota: got %d, want 300", resp.SummaryMap["GroupsQuota"])
		}
		if resp.SummaryMap["ServerCertificatesQuota"] != 20 {
			return fmt.Errorf("ServerCertificatesQuota: got %d, want 20", resp.SummaryMap["ServerCertificatesQuota"])
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

	results = append(results, r.RunTest("iam", "GetAccountAuthorizationDetails_DetailFields", func() error {
		user := fmt.Sprintf("AuthzDetail-User-%s", tc.ts)
		role := fmt.Sprintf("AuthzDetail-Role-%s", tc.ts)
		profile := fmt.Sprintf("AuthzDetail-Prof-%s", tc.ts)
		policy := fmt.Sprintf("AuthzDetail-Policy-%s", tc.ts)
		doc := iamAllowPolicy("s3:ListBucket")

		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		defer tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)})

		policyArn, cleanupPolicy, err := tc.createPolicy(policy, doc)
		if err != nil {
			return err
		}
		defer cleanupPolicy()
		if _, err := tc.client.CreatePolicyVersion(tc.ctx, &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(policyArn),
			PolicyDocument: aws.String(doc),
		}); err != nil {
			return err
		}

		if _, err := tc.client.PutUserPermissionsBoundary(tc.ctx, &iam.PutUserPermissionsBoundaryInput{
			UserName:            aws.String(user),
			PermissionsBoundary: aws.String(policyArn),
		}); err != nil {
			return err
		}
		if _, err := tc.client.TagUser(tc.ctx, &iam.TagUserInput{
			UserName: aws.String(user),
			Tags:     []types.Tag{{Key: aws.String("Purpose"), Value: aws.String("authz-detail")}},
		}); err != nil {
			return err
		}

		cleanupRole, err := tc.createRole(role)
		if err != nil {
			return err
		}
		defer cleanupRole()
		if _, err := tc.client.TagRole(tc.ctx, &iam.TagRoleInput{
			RoleName: aws.String(role),
			Tags:     []types.Tag{{Key: aws.String("Purpose"), Value: aws.String("authz-detail")}},
		}); err != nil {
			return err
		}
		if _, err := tc.client.CreateInstanceProfile(tc.ctx, &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
		}); err != nil {
			return err
		}
		defer tc.client.DeleteInstanceProfile(tc.ctx, &iam.DeleteInstanceProfileInput{InstanceProfileName: aws.String(profile)})
		if _, err := tc.client.AddRoleToInstanceProfile(tc.ctx, &iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
			RoleName:            aws.String(role),
		}); err != nil {
			return err
		}
		defer tc.client.RemoveRoleFromInstanceProfile(tc.ctx, &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
			RoleName:            aws.String(role),
		})

		// Walk every page: during full regression the account carries
		// thousands of entities from other suites' resources, so a single
		// MaxItems-bounded page truncates the tail and hides freshly
		// created entities.
		var userDetailList []types.UserDetail
		var roleDetailList []types.RoleDetail
		var policyDetails []types.ManagedPolicyDetail
		var marker *string
		for {
			resp, err := tc.client.GetAccountAuthorizationDetails(tc.ctx, &iam.GetAccountAuthorizationDetailsInput{
				MaxItems: aws.Int32(1000),
				Marker:   marker,
			})
			if err != nil {
				return err
			}
			userDetailList = append(userDetailList, resp.UserDetailList...)
			roleDetailList = append(roleDetailList, resp.RoleDetailList...)
			policyDetails = append(policyDetails, resp.Policies...)
			if !resp.IsTruncated || resp.Marker == nil {
				break
			}
			marker = resp.Marker
		}

		var userDetail *types.UserDetail
		for i := range userDetailList {
			if aws.ToString(userDetailList[i].UserName) == user {
				userDetail = &userDetailList[i]
				break
			}
		}
		if userDetail == nil {
			return fmt.Errorf("user %s not found in UserDetailList", user)
		}
		if userDetail.PermissionsBoundary == nil || aws.ToString(userDetail.PermissionsBoundary.PermissionsBoundaryArn) != policyArn {
			return fmt.Errorf("user detail missing PermissionsBoundary %s", policyArn)
		}
		if !iamTagPresent(userDetail.Tags, "Purpose", "authz-detail") {
			return fmt.Errorf("user detail missing Tags")
		}

		var roleDetail *types.RoleDetail
		for i := range roleDetailList {
			if aws.ToString(roleDetailList[i].RoleName) == role {
				roleDetail = &roleDetailList[i]
				break
			}
		}
		if roleDetail == nil {
			return fmt.Errorf("role %s not found in RoleDetailList", role)
		}
		if len(roleDetail.InstanceProfileList) != 1 || aws.ToString(roleDetail.InstanceProfileList[0].InstanceProfileName) != profile {
			return fmt.Errorf("role detail missing instance profile %s", profile)
		}
		if !iamTagPresent(roleDetail.Tags, "Purpose", "authz-detail") {
			return fmt.Errorf("role detail missing Tags")
		}

		var policyDetail *types.ManagedPolicyDetail
		for i := range policyDetails {
			if aws.ToString(policyDetails[i].PolicyName) == policy {
				policyDetail = &policyDetails[i]
				break
			}
		}
		if policyDetail == nil {
			return fmt.Errorf("policy %s not found in Policies", policy)
		}
		if policyDetail.CreateDate == nil {
			return fmt.Errorf("policy detail missing CreateDate")
		}
		if policyDetail.UpdateDate == nil {
			return fmt.Errorf("policy detail missing UpdateDate")
		}
		if !policyDetail.IsAttachable {
			return fmt.Errorf("policy detail missing IsAttachable=true")
		}
		if len(policyDetail.PolicyVersionList) != 2 {
			return fmt.Errorf("policy detail version count: got %d, want 2", len(policyDetail.PolicyVersionList))
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

	results = append(results, r.RunTest("iam", "UpdateAccountPasswordPolicy_ReplaceSemantics", func() error {
		// Configure a deliberately strict custom policy first.
		_, err := tc.client.UpdateAccountPasswordPolicy(tc.ctx, &iam.UpdateAccountPasswordPolicyInput{
			MinimumPasswordLength:      aws.Int32(31),
			RequireSymbols:             true,
			RequireNumbers:             true,
			RequireUppercaseCharacters: true,
			RequireLowercaseCharacters: true,
			AllowUsersToChangePassword: true,
			MaxPasswordAge:             aws.Int32(90),
			PasswordReusePrevention:    aws.Int32(10),
			HardExpiry:                 aws.Bool(true),
		})
		if err != nil {
			return err
		}

		// A follow-up update naming only one parameter must replace the
		// whole policy: every unmentioned parameter reverts to its
		// documented default value instead of being merged.
		if _, err := tc.client.UpdateAccountPasswordPolicy(tc.ctx, &iam.UpdateAccountPasswordPolicyInput{
			MinimumPasswordLength: aws.Int32(10),
		}); err != nil {
			return err
		}

		resp, err := tc.client.GetAccountPasswordPolicy(tc.ctx, &iam.GetAccountPasswordPolicyInput{})
		if err != nil {
			return err
		}
		if resp.PasswordPolicy == nil {
			return fmt.Errorf("password policy is nil")
		}
		pp := resp.PasswordPolicy
		if aws.ToInt32(pp.MinimumPasswordLength) != 10 {
			return fmt.Errorf("minimum password length: got %d, want 10", aws.ToInt32(pp.MinimumPasswordLength))
		}
		if pp.RequireSymbols {
			return fmt.Errorf("RequireSymbols must revert to the default value false")
		}
		if pp.RequireNumbers {
			return fmt.Errorf("RequireNumbers must revert to the default value false")
		}
		if pp.RequireUppercaseCharacters {
			return fmt.Errorf("RequireUppercaseCharacters must revert to the default value false")
		}
		if pp.RequireLowercaseCharacters {
			return fmt.Errorf("RequireLowercaseCharacters must revert to the default value false")
		}
		if pp.AllowUsersToChangePassword {
			return fmt.Errorf("AllowUsersToChangePassword must revert to the default value false")
		}
		if aws.ToInt32(pp.MaxPasswordAge) != 0 {
			return fmt.Errorf("MaxPasswordAge must revert to the default value 0, got %d", aws.ToInt32(pp.MaxPasswordAge))
		}
		if aws.ToInt32(pp.PasswordReusePrevention) != 0 {
			return fmt.Errorf("PasswordReusePrevention must revert to the default value 0, got %d", aws.ToInt32(pp.PasswordReusePrevention))
		}
		if aws.ToBool(pp.HardExpiry) {
			return fmt.Errorf("HardExpiry must revert to the default value false")
		}

		// An update with no parameters at all resets everything to the
		// per-parameter defaults.
		if _, err := tc.client.UpdateAccountPasswordPolicy(tc.ctx, &iam.UpdateAccountPasswordPolicyInput{}); err != nil {
			return err
		}
		resp, err = tc.client.GetAccountPasswordPolicy(tc.ctx, &iam.GetAccountPasswordPolicyInput{})
		if err != nil {
			return err
		}
		if aws.ToInt32(resp.PasswordPolicy.MinimumPasswordLength) != 6 {
			return fmt.Errorf("bare update must reset the minimum length to 6, got %d", aws.ToInt32(resp.PasswordPolicy.MinimumPasswordLength))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "OutboundWebIdentityFederation_Lifecycle", func() error {
		// Start from a disabled feature regardless of prior state.
		_, _ = tc.client.DisableOutboundWebIdentityFederation(tc.ctx, &iam.DisableOutboundWebIdentityFederationInput{})

		enable, err := tc.client.EnableOutboundWebIdentityFederation(tc.ctx, &iam.EnableOutboundWebIdentityFederationInput{})
		if err != nil {
			return err
		}
		if aws.ToString(enable.IssuerIdentifier) == "" {
			return fmt.Errorf("enable returned an empty issuer identifier")
		}

		info, err := tc.client.GetOutboundWebIdentityFederationInfo(tc.ctx, &iam.GetOutboundWebIdentityFederationInfoInput{})
		if err != nil {
			return err
		}
		if aws.ToString(info.IssuerIdentifier) != aws.ToString(enable.IssuerIdentifier) {
			return fmt.Errorf("issuer identifier mismatch between enable and get")
		}

		// Enabling twice fails with FeatureEnabled.
		if _, err := tc.client.EnableOutboundWebIdentityFederation(tc.ctx, &iam.EnableOutboundWebIdentityFederationInput{}); err == nil {
			return fmt.Errorf("second enable must fail with FeatureEnabled")
		} else if !containsErrorCode(err, "FeatureEnabled") {
			return fmt.Errorf("second enable: got %v, want FeatureEnabled", err)
		}

		if _, err := tc.client.DisableOutboundWebIdentityFederation(tc.ctx, &iam.DisableOutboundWebIdentityFederationInput{}); err != nil {
			return err
		}

		// Disabling twice and reading a disabled feature both fail with
		// FeatureDisabled.
		if _, err := tc.client.DisableOutboundWebIdentityFederation(tc.ctx, &iam.DisableOutboundWebIdentityFederationInput{}); err == nil {
			return fmt.Errorf("second disable must fail with FeatureDisabled")
		} else if !containsErrorCode(err, "FeatureDisabled") {
			return fmt.Errorf("second disable: got %v, want FeatureDisabled", err)
		}
		if _, err := tc.client.GetOutboundWebIdentityFederationInfo(tc.ctx, &iam.GetOutboundWebIdentityFederationInfoInput{}); err == nil {
			return fmt.Errorf("get while disabled must fail with FeatureDisabled")
		} else if !containsErrorCode(err, "FeatureDisabled") {
			return fmt.Errorf("get while disabled: got %v, want FeatureDisabled", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GenerateServiceLastAccessedDetails_EmptyArn", func() error {
		_, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
			Arn: aws.String(""),
		})
		if err == nil {
			return fmt.Errorf("an empty Arn must be rejected")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServiceLastAccessedDetails", func() error {
		user := fmt.Sprintf("SLA-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()
		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)

		gen, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
			Arn: aws.String(userArn),
		})
		if err != nil {
			return err
		}
		if gen.JobId == nil || *gen.JobId == "" {
			return fmt.Errorf("generate returned an empty job id")
		}

		// The report generation is asynchronous; poll until completion.
		var status string
		for i := 0; i < 20; i++ {
			resp, err := tc.client.GetServiceLastAccessedDetails(tc.ctx, &iam.GetServiceLastAccessedDetailsInput{
				JobId: gen.JobId,
			})
			if err != nil {
				return err
			}
			if resp.JobStatus == types.JobStatusTypeCompleted {
				status = string(resp.JobStatus)
				break
			}
			if resp.JobStatus == types.JobStatusTypeFailed {
				return fmt.Errorf("service last accessed job failed")
			}
			time.Sleep(500 * time.Millisecond)
		}
		if status == "" {
			return fmt.Errorf("service last accessed job did not complete in time")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServiceLastAccessedDetailsWithEntities", func() error {
		role := fmt.Sprintf("SLAEnt-%s", tc.ts)
		cleanupRole, err := tc.createRole(role)
		if err != nil {
			return err
		}
		defer cleanupRole()
		roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, role)

		gen, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
			Arn: aws.String(roleArn),
		})
		if err != nil {
			return err
		}

		// The entity-level report becomes available once the job completes.
		for i := 0; i < 20; i++ {
			resp, err := tc.client.GetServiceLastAccessedDetails(tc.ctx, &iam.GetServiceLastAccessedDetailsInput{
				JobId: gen.JobId,
			})
			if err != nil {
				return err
			}
			if resp.JobStatus == types.JobStatusTypeCompleted {
				break
			}
			if resp.JobStatus == types.JobStatusTypeFailed {
				return fmt.Errorf("service last accessed job failed")
			}
			time.Sleep(500 * time.Millisecond)
		}

		_, err = tc.client.GetServiceLastAccessedDetailsWithEntities(tc.ctx, &iam.GetServiceLastAccessedDetailsWithEntitiesInput{
			JobId:            gen.JobId,
			ServiceNamespace: aws.String("s3"),
		})
		if err != nil {
			return fmt.Errorf("GetServiceLastAccessedDetailsWithEntities: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SetSecurityTokenServicePreferences", func() error {
		if _, err := tc.client.SetSecurityTokenServicePreferences(tc.ctx, &iam.SetSecurityTokenServicePreferencesInput{
			GlobalEndpointTokenVersion: types.GlobalEndpointTokenVersionV2Token,
		}); err != nil {
			return err
		}
		_, err := tc.client.SetSecurityTokenServicePreferences(tc.ctx, &iam.SetSecurityTokenServicePreferencesInput{
			GlobalEndpointTokenVersion: "bogus",
		})
		if err == nil {
			return fmt.Errorf("an invalid token version must be rejected")
		}
		if !isInvalidInputError(err) {
			return fmt.Errorf("invalid token version: got %v, want InvalidInput", err)
		}
		return nil
	}))

	return results
}
