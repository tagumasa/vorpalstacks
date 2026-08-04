package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runSTSAssumeTests(tc *stsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sts", "AssumeRole_ContentVerify", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("VerifySession"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("access key ID is nil or empty")
		}
		if resp.Credentials.SecretAccessKey == nil || *resp.Credentials.SecretAccessKey == "" {
			return fmt.Errorf("secret access key is nil or empty")
		}
		if resp.Credentials.SessionToken == nil || *resp.Credentials.SessionToken == "" {
			return fmt.Errorf("session token is nil or empty")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		if resp.AssumedRoleUser == nil {
			return fmt.Errorf("assumed role user is nil")
		}
		if resp.AssumedRoleUser.AssumedRoleId == nil || *resp.AssumedRoleUser.AssumedRoleId == "" {
			return fmt.Errorf("assumed role ID is nil or empty")
		}
		if resp.AssumedRoleUser.Arn == nil || *resp.AssumedRoleUser.Arn == "" {
			return fmt.Errorf("assumed role user ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithSourceIdentity", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("SourceIdSession"),
			SourceIdentity:  aws.String("AdminUser"),
		})
		if err != nil {
			return err
		}
		if resp.SourceIdentity == nil || *resp.SourceIdentity != "AdminUser" {
			return fmt.Errorf("SourceIdentity not returned correctly, got: %v", resp.SourceIdentity)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithPolicy", func() error {
		inlinePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("PolicySession"),
			Policy:          aws.String(inlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PackedPolicySize == nil || *resp.PackedPolicySize == 0 {
			return fmt.Errorf("PackedPolicySize should be > 0, got: %v", resp.PackedPolicySize)
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithDurationSeconds", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("DurationSession"),
			DurationSeconds: aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_NonExistentRole", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/NonExistentRole", tc.accountID)),
			RoleSessionName: aws.String("TestSession"),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_InvalidDuration", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("DurationSession"),
			DurationSeconds: aws.Int32(100),
		})
		if err := AssertErrorContains(err, "InvalidDuration"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_EmptySessionName", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String(""),
		})
		if err := AssertErrorContains(err, "ValidationError"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_InvalidArn", func() error {
		// "not-a-valid-arn" (15 chars) fails the Smithy arnType min
		// length constraint (20).  AWS returns ValidationError for
		// constraint violations, not InvalidRoleArn.
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String("not-a-valid-arn"),
			RoleSessionName: aws.String("TestSession"),
		})
		if err := AssertErrorContains(err, "ValidationError"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_ShortRoleArn", func() error {
		// 14 characters — well below the arnType minimum (20).
		shortArn := "arn:aws:iam::1"
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(shortArn),
			RoleSessionName: aws.String("TestSession"),
		})
		if err := AssertErrorContains(err, "ValidationError"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoot_Basic", func() error {
		resp, err := tc.client.AssumeRoot(tc.ctx, &sts.AssumeRootInput{
			TargetPrincipal: aws.String(tc.accountID),
			TaskPolicyArn: &types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/IAMAuditRootUserCredentials"),
			},
			DurationSeconds: aws.Int32(900),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("access key ID is nil or empty")
		}
		if resp.Credentials.SecretAccessKey == nil || *resp.Credentials.SecretAccessKey == "" {
			return fmt.Errorf("secret access key is nil or empty")
		}
		if resp.Credentials.SessionToken == nil || *resp.Credentials.SessionToken == "" {
			return fmt.Errorf("session token is nil or empty")
		}
		if resp.Credentials.Expiration.IsZero() {
			return fmt.Errorf("expiration is zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoot_MissingTaskPolicyArn", func() error {
		_, err := tc.client.AssumeRoot(tc.ctx, &sts.AssumeRootInput{
			TargetPrincipal: aws.String(tc.accountID),
			DurationSeconds: aws.Int32(900),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing TaskPolicyArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoot_DurationExceedsMax", func() error {
		_, err := tc.client.AssumeRoot(tc.ctx, &sts.AssumeRootInput{
			TargetPrincipal: aws.String(tc.accountID),
			TaskPolicyArn: &types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/IAMAuditRootUserCredentials"),
			},
			DurationSeconds: aws.Int32(43200),
		})
		if err == nil {
			return fmt.Errorf("expected error for DurationSeconds exceeding 900")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoot_MissingTargetPrincipal", func() error {
		_, err := tc.client.AssumeRoot(tc.ctx, &sts.AssumeRootInput{
			TaskPolicyArn: &types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/IAMAuditRootUserCredentials"),
			},
			DurationSeconds: aws.Int32(900),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing TargetPrincipal")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRoot_InvalidTaskPolicyArn", func() error {
		_, err := tc.client.AssumeRoot(tc.ctx, &sts.AssumeRootInput{
			TargetPrincipal: aws.String(tc.accountID),
			TaskPolicyArn: &types.PolicyDescriptorType{
				Arn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
			},
			DurationSeconds: aws.Int32(900),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid TaskPolicyArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithExternalId", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("ExtIdSession"),
			ExternalId:      aws.String("my-external-id-123"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithMalformedPolicy", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("BadPolicySession"),
			Policy:          aws.String("not-valid-json{"),
		})
		if err == nil {
			return fmt.Errorf("expected error for malformed policy")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "GetFederationToken_InvalidName", func() error {
		_, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name: aws.String("x"),
		})
		if err == nil {
			return fmt.Errorf("expected error for too-short Name")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_WithTags", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("TagSession"),
			Tags: []types.Tag{
				{Key: aws.String("Department"), Value: aws.String("Engineering")},
				{Key: aws.String("Project"), Value: aws.String("Alpha")},
			},
			TransitiveTagKeys: []string{"Department"},
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sts", "AssumeRole_DuplicateTags", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("DupTagSession"),
			Tags: []types.Tag{
				{Key: aws.String("Dept"), Value: aws.String("Eng")},
				{Key: aws.String("Dept"), Value: aws.String("Sales")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate tag keys")
		}
		return nil
	}))

	// sessionPolicyDocumentType max length 2048.
	results = append(results, r.RunTest("sts", "GetFederationToken_PolicyTooLarge", func() error {
		oversized := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"` +
			strings.Repeat("a", 2100) + `"}]}`
		_, err := tc.client.GetFederationToken(tc.ctx, &sts.GetFederationTokenInput{
			Name:   aws.String("BigPolicyFed"),
			Policy: aws.String(oversized),
		})
		if err == nil {
			return fmt.Errorf("expected error for oversized policy (>2048)")
		}
		return nil
	}))

	// PolicyDescriptorType arnType validation.
	results = append(results, r.RunTest("sts", "AssumeRole_InvalidPolicyArnFormat", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("BadArnSession"),
			PolicyArns: []types.PolicyDescriptorType{
				{Arn: aws.String("not-an-arn")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid PolicyArn format")
		}
		return nil
	}))

	// Chained AssumeRole: use temporary credentials from a first
	// AssumeRole to call GetCallerIdentity.  The caller ARN must be
	// the assumed-role ARN (arn:aws:sts::account:assumed-role/...),
	// not the bare role ARN.  This verifies the resolveCallerIdentity
	// fix.
	results = append(results, r.RunTest("sts", "AssumeRole_ChainedCallerIdentity", func() error {
		resp, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("ChainSession1"),
		})
		if err != nil {
			return err
		}
		if resp.Credentials == nil {
			return fmt.Errorf("first AssumeRole returned nil credentials")
		}

		// Build a second STS client authenticated with the temporary
		// credentials from the first AssumeRole.
		chainedCfg, _ := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		chainedCfg.Credentials = credentials.NewStaticCredentialsProvider(
			*resp.Credentials.AccessKeyId,
			*resp.Credentials.SecretAccessKey,
			*resp.Credentials.SessionToken,
		)
		chainedClient := sts.NewFromConfig(chainedCfg)

		idResp, err := chainedClient.GetCallerIdentity(tc.ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}
		if idResp.Arn == nil || *idResp.Arn == "" {
			return fmt.Errorf("chained GetCallerIdentity returned empty ARN")
		}
		// The ARN must be an assumed-role ARN, not a role ARN.
		if !strings.Contains(*idResp.Arn, ":assumed-role/") {
			return fmt.Errorf("expected assumed-role ARN, got: %s", *idResp.Arn)
		}
		if !strings.Contains(*idResp.Arn, "ChainSession1") {
			return fmt.Errorf("expected session name in ARN, got: %s", *idResp.Arn)
		}
		return nil
	}))

	// PolicyArns max 10 limit (AWS docs: "up to 10 managed policy
	// ARNs"). The Smithy policyDescriptorListType has no length trait,
	// so the SDK does not validate client-side.
	results = append(results, r.RunTest("sts", "AssumeRole_TooManyPolicyArns", func() error {
		policyArns := make([]types.PolicyDescriptorType, 11)
		for i := range policyArns {
			policyArns[i] = types.PolicyDescriptorType{
				Arn: aws.String(fmt.Sprintf("arn:aws:iam::123456789012:policy/policy-%d", i)),
			}
		}
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("TooManyArns"),
			PolicyArns:      policyArns,
		})
		if err == nil {
			return fmt.Errorf("expected error for >10 PolicyArns")
		}
		return AssertErrorContains(err, "ValidationError")
	}))

	// TransitiveTagKeys must have a corresponding session tag.
	results = append(results, r.RunTest("sts", "AssumeRole_TransitiveKeyNotInTags", func() error {
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("BadTransitive"),
			Tags: []types.Tag{
				{Key: aws.String("Project"), Value: aws.String("Alpha")},
			},
			TransitiveTagKeys: []string{"Project", "NonExistentKey"},
		})
		if err == nil {
			return fmt.Errorf("expected error for transitive key not in tags")
		}
		return AssertErrorContains(err, "ValidationError")
	}))

	// Dedup check: duplicate PolicyArns are rejected, mirroring the
	// ErrDuplicateSessionTagKey guard for session tags.
	results = append(results, r.RunTest("sts", "AssumeRole_DuplicatePolicyArn", func() error {
		dupArn := "arn:aws:iam::123456789012:policy/dup-policy"
		_, err := tc.client.AssumeRole(tc.ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(tc.roleARN()),
			RoleSessionName: aws.String("DupArnSession"),
			PolicyArns: []types.PolicyDescriptorType{
				{Arn: aws.String(dupArn)},
				{Arn: aws.String(dupArn)},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate PolicyArn")
		}
		return AssertErrorContains(err, "ValidationError")
	}))

	return results
}
