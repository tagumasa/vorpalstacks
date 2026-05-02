package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"vorpalstacks-sdk-tests/config"
)

type iamTestContext struct {
	client  *iam.Client
	ctx     context.Context
	ts      string
	user    string
	group   string
	role    string
	policy  string
	profile string

	userInlinePolicy    string
	roleInlinePolicy    string
	groupInlinePolicy   string
	accessKeyId         string
	policyArn           string
	accountAlias        string
	samlProviderArn     string
	virtualMFASerial    string
	svcLinkedRoleName   string
	oidcProviderArn     string
	sshPublicKeyId      string
	serverCertArn       string
	serviceCredId       string
	serviceCredPassword string
}

func newIAMTestContext(endpoint, region string) (*iamTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: endpoint,
		Region:   region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	return &iamTestContext{
		client:            iam.NewFromConfig(cfg),
		ctx:               context.Background(),
		ts:                ts,
		user:              fmt.Sprintf("TestUser-%s", ts),
		group:             fmt.Sprintf("TestGroup-%s", ts),
		role:              fmt.Sprintf("TestRole-%s", ts),
		policy:            fmt.Sprintf("TestPolicy-%s", ts),
		profile:           fmt.Sprintf("TestProfile-%s", ts),
		userInlinePolicy:  fmt.Sprintf("UserPolicy-%s", ts),
		roleInlinePolicy:  fmt.Sprintf("RolePolicy-%s", ts),
		groupInlinePolicy: fmt.Sprintf("GroupPolicy-%s", ts),
		accountAlias:      fmt.Sprintf("test-alias-%s", ts),
	}, nil
}

const assumeRolePolicy = `{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Principal": {"Service": "lambda.amazonaws.com"},
		"Action": "sts:AssumeRole"
	}]
}`

const s3FullAccessPolicy = `{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Action": "s3:*",
		"Resource": "*"
	}]
}`

const logsFullAccessPolicy = `{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Action": "logs:*",
		"Resource": "*"
	}]
}`

const samlMetadata = `<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://idp.example.com">
  <md:IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.example.com/sso"/>
  </md:IDPSSODescriptor>
</md:EntityDescriptor>`

func (r *TestRunner) RunIAMTests() []TestResult {
	var results []TestResult

	tc, err := newIAMTestContext(r.endpoint, r.region)
	if err != nil {
		return []TestResult{{Service: "iam", TestName: "Setup", Status: "FAIL", Error: err.Error()}}
	}

	results = append(results, r.iamUserTests(tc)...)
	results = append(results, r.iamGroupTests(tc)...)
	results = append(results, r.iamRoleTests(tc)...)
	results = append(results, r.iamPolicyTests(tc)...)
	results = append(results, r.iamPermissionsBoundaryTests(tc)...)
	results = append(results, r.iamInstanceProfileTests(tc)...)
	results = append(results, r.iamAccountTests(tc)...)
	results = append(results, r.iamSAMLTests(tc)...)
	results = append(results, r.iamMFATests(tc)...)
	results = append(results, r.iamAdvancedTests(tc)...)
	results = append(results, r.iamEdgeTests(tc)...)
	results = append(results, r.iamCleanup(tc)...)
	results = append(results, r.iamPaginationTests(tc)...)

	return results
}

func iamFindUserByName(users []types.User, name string) (*types.User, bool) {
	for i := range users {
		if aws.ToString(users[i].UserName) == name {
			return &users[i], true
		}
	}
	return nil, false
}

func iamFindGroupByName(groups []types.Group, name string) (*types.Group, bool) {
	for i := range groups {
		if aws.ToString(groups[i].GroupName) == name {
			return &groups[i], true
		}
	}
	return nil, false
}

func iamFindRoleByPrefix(roles []types.Role, prefix string) []types.Role {
	var matched []types.Role
	for _, r := range roles {
		if strings.HasPrefix(aws.ToString(r.RoleName), prefix) {
			matched = append(matched, r)
		}
	}
	return matched
}

func iamFindAttachedPolicy(policies []types.AttachedPolicy, arn string) bool {
	for _, p := range policies {
		if aws.ToString(p.PolicyArn) == arn {
			return true
		}
	}
	return false
}

func iamTagPresent(tags []types.Tag, key, value string) bool {
	for _, t := range tags {
		if aws.ToString(t.Key) == key && aws.ToString(t.Value) == value {
			return true
		}
	}
	return false
}
