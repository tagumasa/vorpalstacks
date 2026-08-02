package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

	ts := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)

	roleName := fmt.Sprintf("TestRole-%s", ts)
	trustPolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": "arn:aws:iam::%s:root"},
			"Action": "sts:AssumeRole"
		}]
	}`, r.AccountID())

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
	})
	if err != nil {
		results = append(results, SetupFailResult("sts", "failed to create test role: %v", err))
		return nil, results
	}

	samlRole := fmt.Sprintf("SAMLRole-%s", ts)
	samlTrustPolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"Federated": "arn:aws:iam::%s:saml-provider/TestProvider"},
			"Action": "sts:AssumeRoleWithSAML"
		}]
	}`, r.AccountID())

	_, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(samlRole),
		AssumeRolePolicyDocument: aws.String(samlTrustPolicy),
	})
	if err != nil {
		results = append(results, SetupFailResult("sts", "failed to create SAML test role: %v", err))
		return nil, results
	}

	webIdRole := fmt.Sprintf("WebIdRole-%s", ts)
	webIdTrustPolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"Federated": "arn:aws:iam::%s:oidc-provider/example.com"},
			"Action": "sts:AssumeRoleWithWebIdentity"
		}]
	}`, r.AccountID())

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
