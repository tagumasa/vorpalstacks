package testutil

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoPoolCoreTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "DescribeUserPool", func() error {
		resp, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.UserPool == nil {
			return fmt.Errorf("UserPool is nil")
		}
		if resp.UserPool.Id == nil || *resp.UserPool.Id != tc.userPoolID {
			return fmt.Errorf("UserPool.Id mismatch: got %v, want %s", resp.UserPool.Id, tc.userPoolID)
		}
		if resp.UserPool.Name == nil || *resp.UserPool.Name == "" {
			return fmt.Errorf("UserPool.Name is nil or empty")
		}
		if resp.UserPool.Arn == nil || *resp.UserPool.Arn == "" {
			return fmt.Errorf("UserPool.Arn is nil or empty")
		}
		if resp.UserPool.Status != types.StatusTypeEnabled {
			return fmt.Errorf("UserPool.Status: got %q, want Enabled", resp.UserPool.Status)
		}
		if resp.UserPool.Policies == nil || resp.UserPool.Policies.PasswordPolicy == nil {
			return fmt.Errorf("PasswordPolicy is nil")
		}
		if resp.UserPool.Policies.PasswordPolicy.MinimumLength == nil || *resp.UserPool.Policies.PasswordPolicy.MinimumLength != 8 {
			return fmt.Errorf("MinimumLength mismatch: got %v, want 8", resp.UserPool.Policies.PasswordPolicy.MinimumLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetSigningCertificate", func() error {
		resp, err := tc.client.GetSigningCertificate(tc.ctx, &cognitoidentityprovider.GetSigningCertificateInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.Certificate == nil {
			return fmt.Errorf("Certificate is nil")
		}
		block, _ := pem.Decode([]byte(*resp.Certificate))
		if block == nil || block.Type != "CERTIFICATE" {
			return fmt.Errorf("Certificate is not a PEM-encoded X.509 certificate: %.40s", *resp.Certificate)
		}
		if _, err := x509.ParseCertificate(block.Bytes); err != nil {
			return fmt.Errorf("certificate does not parse: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUserPools", func() error {
		pools, err := paginate(func(next *string) ([]types.UserPoolDescriptionType, *string, error) {
			resp, err := tc.client.ListUserPools(tc.ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(10),
				NextToken:  next,
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.UserPools, resp.NextToken, nil
		})
		if err != nil {
			return err
		}
		found := false
		for _, pool := range pools {
			if pool.Id != nil && *pool.Id == tc.userPoolID {
				found = true
				if pool.Name == nil || *pool.Name == "" {
					return fmt.Errorf("pool Name is nil or empty in listing")
				}
			}
		}
		if !found {
			return fmt.Errorf("created pool %s not found in ListUserPools", tc.userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPool", func() error {
		_, err := tc.client.UpdateUserPool(tc.ctx, &cognitoidentityprovider.UpdateUserPoolInput{
			UserPoolId: aws.String(tc.userPoolID),
			Policies: &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength:    aws.Int32(10),
					RequireUppercase: true,
					RequireLowercase: true,
					RequireNumbers:   true,
					RequireSymbols:   true,
				},
			},
		})
		if err != nil {
			return err
		}
		descResp, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return fmt.Errorf("DescribeUserPool after update: %v", err)
		}
		pp := descResp.UserPool.Policies.PasswordPolicy
		if pp.MinimumLength == nil || *pp.MinimumLength != 10 {
			return fmt.Errorf("MinimumLength not updated: got %v, want 10", pp.MinimumLength)
		}
		if !pp.RequireSymbols {
			return fmt.Errorf("RequireSymbols not updated to true")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UserPoolConfig_SignInPolicyEumsSmsAdditionalFlows", func() error {
		poolID, cleanup, err := tc.createUserPool(tc.unique("config-pool"), func(input *cognitoidentityprovider.CreateUserPoolInput) {
			input.Policies = &types.UserPoolPolicyType{
				SignInPolicy: &types.SignInPolicyType{
					AllowedFirstAuthFactors: []types.AuthFactorType{
						types.AuthFactorTypePassword, types.AuthFactorTypeEmailOtp,
					},
				},
			}
			input.SmsConfiguration = &types.SmsConfigurationType{
				EumsSms: &types.EumsSmsConfigurationType{
					CallerArn:           aws.String("arn:aws:iam::123456789012:role/eums-sender"),
					ExternalId:          aws.String("eums-external-id"),
					OriginationIdentity: aws.String("+12065550100"),
				},
			}
			input.UserPoolAddOns = &types.UserPoolAddOnsType{
				AdvancedSecurityMode: types.AdvancedSecurityModeTypeAudit,
				AdvancedSecurityAdditionalFlows: &types.AdvancedSecurityAdditionalFlowsType{
					CustomAuthMode: types.AdvancedSecurityEnabledModeTypeEnforced,
				},
			}
		})
		if err != nil {
			return err
		}
		defer cleanup()

		desc, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if desc.UserPool.Policies == nil || desc.UserPool.Policies.SignInPolicy == nil {
			return fmt.Errorf("describe lost the sign-in policy")
		}
		factors := desc.UserPool.Policies.SignInPolicy.AllowedFirstAuthFactors
		if len(factors) != 2 || factors[0] != types.AuthFactorTypePassword || factors[1] != types.AuthFactorTypeEmailOtp {
			return fmt.Errorf("AllowedFirstAuthFactors round-trip mismatch: %v", factors)
		}
		sms := desc.UserPool.SmsConfiguration
		if sms == nil || sms.EumsSms == nil || sms.EumsSms.CallerArn == nil ||
			*sms.EumsSms.CallerArn != "arn:aws:iam::123456789012:role/eums-sender" ||
			sms.EumsSms.ExternalId == nil || *sms.EumsSms.ExternalId != "eums-external-id" ||
			sms.EumsSms.OriginationIdentity == nil || *sms.EumsSms.OriginationIdentity != "+12065550100" {
			return fmt.Errorf("EumsSms round-trip mismatch: %+v", sms)
		}
		addOns := desc.UserPool.UserPoolAddOns
		if addOns == nil || addOns.AdvancedSecurityAdditionalFlows == nil ||
			addOns.AdvancedSecurityAdditionalFlows.CustomAuthMode != types.AdvancedSecurityEnabledModeTypeEnforced {
			return fmt.Errorf("AdvancedSecurityAdditionalFlows round-trip mismatch: %+v", addOns)
		}

		// The update request rebuilds the configuration: omitting these
		// members reverts them to their defaults (no sign-in policy, no SMS
		// configuration, no add-ons).
		if _, err := tc.client.UpdateUserPool(tc.ctx, &cognitoidentityprovider.UpdateUserPoolInput{
			UserPoolId: aws.String(poolID),
		}); err != nil {
			return err
		}
		desc, err = tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if desc.UserPool.Policies != nil && desc.UserPool.Policies.SignInPolicy != nil {
			return fmt.Errorf("an update omitting Policies kept the sign-in policy")
		}
		if desc.UserPool.SmsConfiguration != nil {
			return fmt.Errorf("an update omitting SmsConfiguration kept the EumsSms configuration")
		}
		if desc.UserPool.UserPoolAddOns != nil {
			return fmt.Errorf("an update omitting UserPoolAddOns kept the additional flows")
		}

		var invalidParam *types.InvalidParameterException
		if _, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("signin-software-token")),
			Policies: &types.UserPoolPolicyType{
				SignInPolicy: &types.SignInPolicyType{
					AllowedFirstAuthFactors: []types.AuthFactorType{
						types.AuthFactorTypePassword, types.AuthFactorTypeSoftwareToken,
					},
				},
			},
		}); !errors.As(err, &invalidParam) {
			return fmt.Errorf("SOFTWARE_TOKEN as a first auth factor must be rejected as InvalidParameter, got %v", err)
		}
		if _, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("sms-both-channels")),
			SmsConfiguration: &types.SmsConfigurationType{
				SnsCallerArn: aws.String("arn:aws:iam::123456789012:role/sns-sender"),
				EumsSms: &types.EumsSmsConfigurationType{
					CallerArn: aws.String("arn:aws:iam::123456789012:role/eums-sender"),
				},
			},
		}); !errors.As(err, &invalidParam) {
			return fmt.Errorf("SnsCallerArn combined with EumsSms must be rejected as InvalidParameter, got %v", err)
		}
		// CustomAuthMode values outside the enabled-mode enum are pinned at
		// the unit level: the SDK validates the closed enum client-side, so
		// an off-enum value never reaches the server through this client.
		return nil
	}))

	results = append(results, r.RunTest("cognito", "SetUserPoolMfaConfig", func() error {
		_, err := tc.client.SetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:       aws.String(tc.userPoolID),
			MfaConfiguration: types.UserPoolMfaTypeOn,
		})
		if err != nil {
			return fmt.Errorf("SetUserPoolMfaConfig failed: %v", err)
		}
		mfaResp, err := tc.client.GetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return fmt.Errorf("GetUserPoolMfaConfig failed: %v", err)
		}
		if mfaResp.MfaConfiguration != types.UserPoolMfaTypeOn {
			return fmt.Errorf("expected MfaConfiguration ON, got %v", mfaResp.MfaConfiguration)
		}
		_, err = tc.client.SetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:       aws.String(tc.userPoolID),
			MfaConfiguration: types.UserPoolMfaTypeOff,
		})
		if err != nil {
			return fmt.Errorf("SetUserPoolMfaConfig (OFF) failed: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetUserPoolMfaConfig", func() error {
		resp, err := tc.client.GetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.MfaConfiguration != types.UserPoolMfaTypeOff {
			return fmt.Errorf("expected MfaConfiguration OFF after reset, got %v", resp.MfaConfiguration)
		}
		return nil
	}))

	domainName := tc.unique("test-domain")
	results = append(results, r.RunTest("cognito", "CreateUserPoolDomain", func() error {
		resp, err := tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.CloudFrontDomain == nil || *resp.CloudFrontDomain == "" {
			return fmt.Errorf("CloudFrontDomain is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeUserPoolDomain", func() error {
		resp, err := tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainDescription == nil {
			return fmt.Errorf("DomainDescription is nil")
		}
		if resp.DomainDescription.UserPoolId == nil || *resp.DomainDescription.UserPoolId != tc.userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.DomainDescription.UserPoolId, tc.userPoolID)
		}
		if resp.DomainDescription.Domain == nil || *resp.DomainDescription.Domain != domainName {
			return fmt.Errorf("domain mismatch: got %v, want %s", resp.DomainDescription.Domain, domainName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPoolDomain", func() error {
		// A pool holds one domain at a time, so the update exercises the
		// domain the pool already owns rather than creating a second one.
		resp, err := tc.client.UpdateUserPoolDomain(tc.ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
			Domain:              aws.String(domainName),
			UserPoolId:          aws.String(tc.userPoolID),
			ManagedLoginVersion: aws.Int32(2),
			CustomDomainConfig: &types.CustomDomainConfigType{
				CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/ud-test"),
			},
			Routing: &types.RoutingType{
				Failover: &types.FailoverType{
					SecondaryRegion:             aws.String("us-west-2"),
					PrimaryRoute53HealthCheckId: aws.String("hc-ud-test"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateUserPoolDomain failed: %v", err)
		}
		if resp.CloudFrontDomain == nil || *resp.CloudFrontDomain == "" {
			return fmt.Errorf("CloudFrontDomain is nil or empty")
		}
		if resp.ManagedLoginVersion == nil || *resp.ManagedLoginVersion != 2 {
			return fmt.Errorf("ManagedLoginVersion mismatch: got %v, want 2", resp.ManagedLoginVersion)
		}
		if resp.Routing == nil || resp.Routing.Failover == nil || resp.Routing.Failover.SecondaryRegion == nil || *resp.Routing.Failover.SecondaryRegion != "us-west-2" {
			return fmt.Errorf("Routing not reflected in update response: %v", resp.Routing)
		}
		// A second update touching only ManagedLoginVersion must preserve
		// the members set by the first update.
		_, err = tc.client.UpdateUserPoolDomain(tc.ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
			Domain:              aws.String(domainName),
			UserPoolId:          aws.String(tc.userPoolID),
			ManagedLoginVersion: aws.Int32(3),
		})
		if err != nil {
			return fmt.Errorf("UpdateUserPoolDomain (ManagedLoginVersion only) failed: %v", err)
		}
		desc, err := tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		d := desc.DomainDescription
		if d == nil {
			return fmt.Errorf("DomainDescription is nil after update")
		}
		if d.ManagedLoginVersion == nil || *d.ManagedLoginVersion != 3 {
			return fmt.Errorf("stored ManagedLoginVersion mismatch: got %v, want 3", d.ManagedLoginVersion)
		}
		if d.CustomDomainConfig == nil || d.CustomDomainConfig.CertificateArn == nil || *d.CustomDomainConfig.CertificateArn == "" {
			return fmt.Errorf("first-update CustomDomainConfig not preserved by second update")
		}
		if d.Routing == nil || d.Routing.Failover == nil || d.Routing.Failover.SecondaryRegion == nil || *d.Routing.Failover.SecondaryRegion != "us-west-2" {
			return fmt.Errorf("first-update Routing not preserved by second update")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPoolDomain_SecurityPolicyEnumRejected", func() error {
		_, err := tc.client.UpdateUserPoolDomain(tc.ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(tc.userPoolID),
			CustomDomainConfig: &types.CustomDomainConfigType{
				// CertificateArn is SDK-required; the off-enum SecurityPolicy is
				// the member under test (the SDK does not validate enum values).
				CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/sp-test"),
				SecurityPolicy: types.SecurityPolicyType("not-a-policy"),
			},
		})
		if err := expectAWSErrorCode(err, "InvalidParameterException"); err != nil {
			return err
		}
		// A rejected update must leave the binding untouched.
		desc, derr := tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if derr != nil {
			return fmt.Errorf("describe after rejected update: %v", derr)
		}
		if desc.DomainDescription == nil || desc.DomainDescription.UserPoolId == nil || *desc.DomainDescription.UserPoolId != tc.userPoolID {
			return fmt.Errorf("rejected update disturbed the domain binding: %v", desc.DomainDescription)
		}
		return nil
	}))

	// A domain string binds to at most one pool, and a pool owns at most one
	// domain: both duplicate paths are InvalidParameterException, the error
	// every domain operation's model error list documents for duplicate
	// bindings.
	results = append(results, r.RunTest("cognito", "CreateUserPoolDomain_DuplicateRejected", func() error {
		secondPoolID, _, cleanupPool, err := tc.createUserPoolWithArn(tc.unique("test-pool-domain-dup"))
		if err != nil {
			return err
		}
		defer cleanupPool()
		fresh := tc.unique("dup-fresh-domain")
		if _, err := tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(fresh),
			UserPoolId: aws.String(secondPoolID),
		}); err != nil {
			return fmt.Errorf("create fresh domain: %v", err)
		}
		defer tc.client.DeleteUserPoolDomain(tc.ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(fresh),
			UserPoolId: aws.String(secondPoolID),
		})
		_, err = tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(tc.unique("dup-second-domain")),
			UserPoolId: aws.String(secondPoolID),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("second domain for one pool: %v", err)
		}
		_, err = tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(secondPoolID),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("cross-pool domain claim: %v", err)
		}
		return nil
	}))

	// Delete and update carry the owning pool's ID: operating on another
	// pool's domain is ResourceNotFoundException and must not disturb the
	// owner's binding.
	results = append(results, r.RunTest("cognito", "UserPoolDomain_ForeignPoolRejected", func() error {
		secondPoolID, _, cleanupPool, err := tc.createUserPoolWithArn(tc.unique("test-pool-domain-foreign"))
		if err != nil {
			return err
		}
		defer cleanupPool()
		_, err = tc.client.DeleteUserPoolDomain(tc.ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(secondPoolID),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("delete of foreign domain: %v", err)
		}
		_, err = tc.client.UpdateUserPoolDomain(tc.ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
			Domain:              aws.String(domainName),
			UserPoolId:          aws.String(secondPoolID),
			ManagedLoginVersion: aws.Int32(9),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("update of foreign domain: %v", err)
		}
		desc, err := tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err != nil {
			return fmt.Errorf("describe after foreign calls: %v", err)
		}
		if desc.DomainDescription == nil || desc.DomainDescription.UserPoolId == nil || *desc.DomainDescription.UserPoolId != tc.userPoolID {
			return fmt.Errorf("foreign calls disturbed the owner binding: %v", desc.DomainDescription)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteUserPoolDomain", func() error {
		_, err := tc.client.DeleteUserPoolDomain(tc.ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted domain")
		}
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("deleted domain describe: %v", err)
		}
		return nil
	}))

	// Deleting a pool releases everything the pool owned: the domain string
	// of a deleted pool binds again to a fresh pool, where an orphaned
	// binding would reject the claim as a duplicate.
	results = append(results, r.RunTest("cognito", "DeleteUserPool_CascadeReleasesDomain", func() error {
		firstPoolID, _, cleanupFirst, err := tc.createUserPoolWithArn(tc.unique("test-pool-cascade-a"))
		if err != nil {
			return err
		}
		cascadeDomain := tc.unique("cascade-domain")
		if _, err := tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(cascadeDomain),
			UserPoolId: aws.String(firstPoolID),
		}); err != nil {
			cleanupFirst()
			return fmt.Errorf("bind domain: %v", err)
		}
		if _, err := tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: aws.String(firstPoolID),
		}); err != nil {
			cleanupFirst()
			return fmt.Errorf("delete first pool: %v", err)
		}
		secondPoolID, _, cleanupSecond, err := tc.createUserPoolWithArn(tc.unique("test-pool-cascade-b"))
		if err != nil {
			return err
		}
		defer cleanupSecond()
		if _, err := tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(cascadeDomain),
			UserPoolId: aws.String(secondPoolID),
		}); err != nil {
			return fmt.Errorf("rebind domain after pool deletion: %v", err)
		}
		desc, err := tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(cascadeDomain),
		})
		if err != nil {
			return fmt.Errorf("describe rebound domain: %v", err)
		}
		if desc.DomainDescription == nil || desc.DomainDescription.UserPoolId == nil || *desc.DomainDescription.UserPoolId != secondPoolID {
			return fmt.Errorf("rebound domain description: %+v", desc.DomainDescription)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetCSVHeader", func() error {
		resp, err := tc.client.GetCSVHeader(tc.ctx, &cognitoidentityprovider.GetCSVHeaderInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if len(resp.CSVHeader) == 0 {
			return fmt.Errorf("expected non-empty CSV header")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeRiskConfiguration", func() error {
		resp, err := tc.client.DescribeRiskConfiguration(tc.ctx, &cognitoidentityprovider.DescribeRiskConfigurationInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.RiskConfiguration == nil {
			return fmt.Errorf("RiskConfiguration is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "TagResource", func() error {
		_, poolArn, cleanupPool, err := tc.createUserPoolWithArn(tc.unique("test-pool-tags"))
		if err != nil {
			return err
		}
		defer cleanupPool()
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: aws.String(poolArn),
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: aws.String(poolArn),
		})
		if err != nil {
			return err
		}
		if listResp.Tags == nil {
			return fmt.Errorf("tags is nil after tagging")
		}
		if listResp.Tags["Environment"] != "test" {
			return fmt.Errorf("tag Environment not found after TagResource")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "TagResource_ReservedPrefixRejected", func() error {
		_, poolArn, cleanupPool, err := tc.createUserPoolWithArn(tc.unique("test-pool-tagres"))
		if err != nil {
			return err
		}
		defer cleanupPool()
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: aws.String(poolArn),
			Tags: map[string]string{
				"aws:reserved": "v",
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag key")
		}
		return nil
	}))

	// Tag operations against a user pool that does not exist fail with
	// ResourceNotFoundException, as the service model specifies.
	results = append(results, r.RunTest("cognito", "TagResource_NonExistentPool", func() error {
		arn := fmt.Sprintf("arn:aws:cognito-idp:%s:%s:userpool/%s_no-such-pool", r.region, r.accountID, r.region)
		_, err := tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags:        map[string]string{"Environment": "test"},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("TagResource: %v", err)
		}
		_, err = tc.client.UntagResource(tc.ctx, &cognitoidentityprovider.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"Environment"},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("UntagResource: %v", err)
		}
		_, err = tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: aws.String(arn),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("ListTagsForResource: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListTagsForResource", func() error {
		_, poolArn, cleanupPool, err := tc.createUserPoolWithArn(tc.unique("test-pool-listtags"))
		if err != nil {
			return err
		}
		defer cleanupPool()
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: aws.String(poolArn),
			Tags: map[string]string{
				"Test": "value",
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: aws.String(poolArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["Test"] != "value" {
			return fmt.Errorf("expected tag Test=value, got %v", resp.Tags["Test"])
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UntagResource", func() error {
		_, poolArn, cleanupPool, err := tc.createUserPoolWithArn(tc.unique("test-pool-untag"))
		if err != nil {
			return err
		}
		defer cleanupPool()
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: aws.String(poolArn),
			Tags: map[string]string{
				"Test": "value",
			},
		})
		if err != nil {
			return err
		}
		_, err = tc.client.UntagResource(tc.ctx, &cognitoidentityprovider.UntagResourceInput{
			ResourceArn: aws.String(poolArn),
			TagKeys:     []string{"Test"},
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: aws.String(poolArn),
		})
		if err != nil {
			return err
		}
		if _, exists := listResp.Tags["Test"]; exists {
			return fmt.Errorf("tag Test should have been removed after UntagResource")
		}
		return nil
	}))

	return results
}
