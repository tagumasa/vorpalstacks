package testutil

import (
	"fmt"
	"strings"

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

		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

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
		cleanupProfile, err := tc.createInstanceProfile(profile)
		if err != nil {
			return err
		}
		defer cleanupProfile()
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

	// A paged walk must lose no policy: the response Marker names the
	// first item of the next page, so resuming at — not after — the
	// marker keeps every entity. Eight fresh policies force several
	// MaxItems=3 pages even on a freshly initialised account, and the
	// union of the small pages must equal the single large page.
	results = append(results, r.RunTest("iam", "GetAccountAuthorizationDetails_PaginationComplete", func() error {
		var cleanups []func()
		defer func() {
			for i := len(cleanups) - 1; i >= 0; i-- {
				cleanups[i]()
			}
		}()
		doc := iamAllowPolicy("s3:GetObject")
		for i := 0; i < 8; i++ {
			name := fmt.Sprintf("AuthzPage-Policy-%s-%d", tc.ts, i)
			_, cleanup, err := tc.createPolicy(name, doc)
			if err != nil {
				return err
			}
			cleanups = append(cleanups, cleanup)
		}

		walk := func(maxItems int32) (map[string]bool, int, error) {
			seen := map[string]bool{}
			pages := 0
			var marker *string
			for {
				input := &iam.GetAccountAuthorizationDetailsInput{
					Filter:   []types.EntityType{types.EntityTypeLocalManagedPolicy},
					MaxItems: aws.Int32(maxItems),
				}
				if marker != nil {
					input.Marker = marker
				}
				resp, err := tc.client.GetAccountAuthorizationDetails(tc.ctx, input)
				if err != nil {
					return nil, 0, err
				}
				for _, p := range resp.Policies {
					seen[aws.ToString(p.Arn)] = true
				}
				pages++
				if !resp.IsTruncated || resp.Marker == nil {
					break
				}
				marker = resp.Marker
			}
			return seen, pages, nil
		}

		paged, pages, err := walk(3)
		if err != nil {
			return err
		}
		if pages < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d", pages)
		}
		direct, _, err := walk(1000)
		if err != nil {
			return err
		}
		if len(paged) != len(direct) {
			missing := []string{}
			for arn := range direct {
				if !paged[arn] {
					missing = append(missing, arn)
				}
			}
			return fmt.Errorf("paged walk saw %d policies over %d pages, single page saw %d; missing: %v",
				len(paged), pages, len(direct), missing)
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

	results = append(results, r.RunTest("iam", "GenerateServiceLastAccessedDetails_ArnLengthRejected", func() error {
		for _, arn := range []string{"a", fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, strings.Repeat("x", 2100))} {
			_, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
				Arn: aws.String(arn),
			})
			if err == nil {
				return fmt.Errorf("Arn of length %d must be rejected", len(arn))
			}
			if !containsErrorCode(err, "InvalidInput") {
				return fmt.Errorf("Arn length %d: got %v, want InvalidInput", len(arn), err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GenerateServiceLastAccessedDetails_InvalidGranularityRejected", func() error {
		_, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
			Arn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:user/granularity-invalid", tc.accountID)),
			Granularity: types.AccessAdvisorUsageGranularityType("P30D"),
		})
		if err == nil {
			return fmt.Errorf("an ISO-duration granularity must be rejected")
		}
		if !containsErrorCode(err, "InvalidInput") {
			return fmt.Errorf("invalid granularity: got %v, want InvalidInput", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServiceLastAccessedDetails", func() error {
		// A job id of any length other than the fixed 36 characters is
		// malformed input, not a missing report.
		if _, err := tc.client.GetServiceLastAccessedDetails(tc.ctx, &iam.GetServiceLastAccessedDetailsInput{
			JobId: aws.String("short"),
		}); err == nil {
			return fmt.Errorf("a malformed job id must be rejected")
		} else if !containsErrorCode(err, "InvalidInput") {
			return fmt.Errorf("malformed job id: got %v, want InvalidInput", err)
		}

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
		completed, err := tc.awaitServiceLastAccessedDetails(gen.JobId)
		if err != nil {
			return err
		}

		// A generate call without Granularity defaults to SERVICE_LEVEL, and
		// the completed response reports the granularity through JobType.
		if completed.JobType != types.AccessAdvisorUsageGranularityTypeServiceLevel {
			return fmt.Errorf("defaulted report JobType: got %q, want SERVICE_LEVEL", completed.JobType)
		}
		if completed.ServicesLastAccessed == nil {
			return fmt.Errorf("completed report must carry the ServicesLastAccessed member")
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
		if _, err := tc.awaitServiceLastAccessedDetails(gen.JobId); err != nil {
			return err
		}

		resp, err := tc.client.GetServiceLastAccessedDetailsWithEntities(tc.ctx, &iam.GetServiceLastAccessedDetailsWithEntitiesInput{
			JobId:            gen.JobId,
			ServiceNamespace: aws.String("s3"),
		})
		if err != nil {
			return fmt.Errorf("GetServiceLastAccessedDetailsWithEntities: %w", err)
		}
		if resp.JobStatus != types.JobStatusTypeCompleted {
			return fmt.Errorf("entity-level report JobStatus: got %q, want COMPLETED", resp.JobStatus)
		}
		if resp.EntityDetailsList == nil {
			return fmt.Errorf("entity-level report must carry the EntityDetailsList member")
		}
		// Each entity detail carries the documented EntityInfo shape.
		for _, detail := range resp.EntityDetailsList {
			if detail.EntityInfo == nil {
				return fmt.Errorf("entity detail must carry EntityInfo")
			}
			if detail.EntityInfo.Arn == nil || *detail.EntityInfo.Arn == "" {
				return fmt.Errorf("entity EntityInfo must carry the entity ARN")
			}
			if detail.EntityInfo.Type != types.PolicyOwnerEntityTypeUser && detail.EntityInfo.Type != types.PolicyOwnerEntityTypeRole && detail.EntityInfo.Type != types.PolicyOwnerEntityTypeGroup {
				return fmt.Errorf("entity EntityInfo Type: got %q, want USER/ROLE/GROUP", detail.EntityInfo.Type)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ServiceLastAccessed_ActionLevelGranularity", func() error {
		user := fmt.Sprintf("SLAAct-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()
		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)

		gen, err := tc.client.GenerateServiceLastAccessedDetails(tc.ctx, &iam.GenerateServiceLastAccessedDetailsInput{
			Arn:         aws.String(userArn),
			Granularity: types.AccessAdvisorUsageGranularityTypeActionLevel,
		})
		if err != nil {
			return err
		}
		if gen.JobId == nil || *gen.JobId == "" {
			return fmt.Errorf("generate returned an empty job id")
		}

		completed, err := tc.awaitServiceLastAccessedDetails(gen.JobId)
		if err != nil {
			return err
		}
		if completed.JobType != types.AccessAdvisorUsageGranularityTypeActionLevel {
			return fmt.Errorf("action-level report JobType: got %q, want ACTION_LEVEL", completed.JobType)
		}

		// The entity-level view scoped to a namespace stays callable on the
		// action-level job.
		if _, err := tc.client.GetServiceLastAccessedDetailsWithEntities(tc.ctx, &iam.GetServiceLastAccessedDetailsWithEntitiesInput{
			JobId:            gen.JobId,
			ServiceNamespace: aws.String("iam"),
		}); err != nil {
			return fmt.Errorf("entity view on action-level job: %w", err)
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
