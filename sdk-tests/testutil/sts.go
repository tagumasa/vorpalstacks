package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"vorpalstacks-sdk-tests/config"
)

type stsTestContext struct {
	client    *sts.Client
	iamClient *iam.Client
	ctx       context.Context
	ts        string
	roleName  string
	samlRole  string
	webIdRole string
	region    string
	accountID string
}

func newSTSTestContext(r *TestRunner) (*stsTestContext, []TestResult) {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		results = append(results, TestResult{
			Service:  "sts",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
		return nil, results
	}

	client := sts.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	ctx := context.Background()

	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	roleName := fmt.Sprintf("TestRole-%s", ts)
	trustPolicy := stsTrustPolicy("AWS", fmt.Sprintf("arn:aws:iam::%s:root", r.AccountID()), "sts:AssumeRole")

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
	})
	if err != nil {
		results = append(results, SetupFailResult("sts", "failed to create test role: %v", err))
		return nil, results
	}

	samlRole := fmt.Sprintf("SAMLRole-%s", ts)
	samlTrustPolicy := stsTrustPolicy("Federated",
		fmt.Sprintf("arn:aws:iam::%s:saml-provider/TestProvider", r.AccountID()), "sts:AssumeRoleWithSAML")

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(samlRole),
		AssumeRolePolicyDocument: aws.String(samlTrustPolicy),
	})
	if err != nil {
		results = append(results, SetupFailResult("sts", "failed to create SAML test role: %v", err))
		return nil, results
	}

	webIdRole := fmt.Sprintf("WebIdRole-%s", ts)
	webIdTrustPolicy := stsTrustPolicy("Federated",
		fmt.Sprintf("arn:aws:iam::%s:oidc-provider/example.com", r.AccountID()), "sts:AssumeRoleWithWebIdentity")

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(webIdRole),
		AssumeRolePolicyDocument: aws.String(webIdTrustPolicy),
	})
	if err != nil {
		results = append(results, SetupFailResult("sts", "failed to create WebIdentity test role: %v", err))
		return nil, results
	}

	tc := &stsTestContext{
		client:    client,
		iamClient: iamClient,
		ctx:       ctx,
		ts:        ts,
		roleName:  roleName,
		samlRole:  samlRole,
		webIdRole: webIdRole,
		region:    r.region,
		accountID: r.accountID,
	}

	return tc, results
}

func (tc *stsTestContext) cleanup() {
	_, _ = tc.iamClient.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(tc.roleName)})
	_, _ = tc.iamClient.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(tc.samlRole)})
	_, _ = tc.iamClient.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(tc.webIdRole)})
}

func (tc *stsTestContext) roleARN() string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, tc.roleName)
}

func (tc *stsTestContext) samlRoleARN() string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, tc.samlRole)
}

func (tc *stsTestContext) webIdRoleARN() string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, tc.webIdRole)
}

// stsTrustPolicy assembles an AssumeRole trust policy document for the
// given principal kind ("AWS" or "Federated"), principal ARN, and
// allowed action.
func stsTrustPolicy(principalKind, principal, action string) string {
	return fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"%s": "%s"},
			"Action": "%s"
		}]
	}`, principalKind, principal, action)
}

// stsAssertCredentials verifies the four mandatory members of a
// credentials block returned by an STS session operation.
func stsAssertCredentials(creds *types.Credentials) error {
	if creds == nil {
		return fmt.Errorf("credentials is nil")
	}
	if creds.AccessKeyId == nil || *creds.AccessKeyId == "" {
		return fmt.Errorf("access key ID is nil or empty")
	}
	if creds.SecretAccessKey == nil || *creds.SecretAccessKey == "" {
		return fmt.Errorf("secret access key is nil or empty")
	}
	if creds.SessionToken == nil || *creds.SessionToken == "" {
		return fmt.Errorf("session token is nil or empty")
	}
	if creds.Expiration.IsZero() {
		return fmt.Errorf("expiration is zero")
	}
	return nil
}

// stsAssertAssumedRoleUser verifies the assumed-role user block
// returned by the AssumeRole family of operations.
func stsAssertAssumedRoleUser(user *types.AssumedRoleUser) error {
	if user == nil {
		return fmt.Errorf("assumed role user is nil")
	}
	if user.AssumedRoleId == nil || *user.AssumedRoleId == "" {
		return fmt.Errorf("assumed role ID is nil or empty")
	}
	if user.Arn == nil || *user.Arn == "" {
		return fmt.Errorf("assumed role user ARN is nil or empty")
	}
	return nil
}

// stsAssertPackedPolicySize verifies the packed-policy-size member
// returned when an inline session policy is supplied.
func stsAssertPackedPolicySize(size *int32) error {
	if size == nil || *size == 0 {
		return fmt.Errorf("PackedPolicySize should be > 0, got: %v", size)
	}
	return nil
}

func (r *TestRunner) RunSTSTests() []TestResult {
	tc, results := newSTSTestContext(r)
	if tc == nil {
		return results
	}
	defer tc.cleanup()

	results = append(results, r.runSTSIdentityTests(tc)...)
	results = append(results, r.runSTSAssumeTests(tc)...)
	results = append(results, r.runSTSSAMLTests(tc)...)
	results = append(results, r.runSTSWebIdentityTests(tc)...)
	results = append(results, r.runSTSFederationTests(tc)...)
	results = append(results, r.runSTSDecodeTests(tc)...)

	return results
}
