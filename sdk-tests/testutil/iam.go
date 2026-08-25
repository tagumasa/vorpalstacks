package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"vorpalstacks-sdk-tests/config"
)

type iamTestContext struct {
	client    *iam.Client
	ctx       context.Context
	ts        string
	user      string
	group     string
	role      string
	policy    string
	profile   string
	region    string
	accountID string

	userInlinePolicy    string
	roleInlinePolicy    string
	groupInlinePolicy   string
	accessKeyId         string
	policyArn           string
	accountAlias        string
	samlProviderArn     string
	samlPrivateKeyId    string
	virtualMFASerial    string
	svcLinkedRoleName   string
	deletionTaskId      string
	oidcProviderArn     string
	sshPublicKeyId      string
	serverCertArn       string
	serviceCredId       string
	serviceCredPassword string
}

func newIAMTestContext(endpoint, region, accountID string) (*iamTestContext, error) {
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
		region:            region,
		accountID:         accountID,
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

// ec2AssumeRolePolicy is the standard EC2 trust policy used whenever a
// test needs a throwaway role.
const ec2AssumeRolePolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`

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
  <!-- padding for Smithy SAMLMetadataDocumentType min length 1000: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX -->
  </md:IDPSSODescriptor>
</md:EntityDescriptor>`

func (r *TestRunner) RunIAMTests() []TestResult {
	var results []TestResult

	tc, err := newIAMTestContext(r.endpoint, r.region, r.accountID)
	if err != nil {
		return []TestResult{{Service: "iam", TestName: "Setup", Status: "FAIL", Error: err.Error()}}
	}

	results = append(results, r.iamUserTests(tc)...)
	results = append(results, r.iamSigningCertificateTests(tc)...)
	results = append(results, r.iamSSHPublicKeyTests(tc)...)
	results = append(results, r.iamGroupTests(tc)...)
	results = append(results, r.iamDeleteOperationTests(tc)...)
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

// isInvalidInputError reports whether err carries the InvalidInput error
// code, matching either the SDK error type or the wire error string.
func isInvalidInputError(err error) bool {
	var invalidInput *types.InvalidInputException
	if errors.As(err, &invalidInput) {
		return true
	}
	return err != nil && strings.Contains(err.Error(), "InvalidInput")
}

// containsErrorCode reports whether err carries the given AWS error code.
// The SDK surfaces modelled exceptions by shape name (e.g.
// LimitExceededException) while the wire code lacks the suffix, so both
// spellings are accepted.
func containsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}
	norm := func(c string) string { return strings.TrimSuffix(c, "Exception") }
	var awsErr interface{ ErrorCode() string }
	if errors.As(err, &awsErr) && norm(awsErr.ErrorCode()) == norm(code) {
		return true
	}
	return strings.Contains(err.Error(), "api error "+code)
}

// iamAssertNoSuchEntity verifies that err carries NoSuchEntity, as
// required after a successful delete of the named entity.
func iamAssertNoSuchEntity(err error, what string) error {
	if err == nil || !containsErrorCode(err, "NoSuchEntity") {
		return fmt.Errorf("%s: got %v, want NoSuchEntity", what, err)
	}
	return nil
}

// iamAllowPolicy renders a single-statement allow-all-resources policy
// document for the given action.
func iamAllowPolicy(action string) string {
	return fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":%q,"Resource":"*"}]}`, action)
}

// iamPaginate walks every page of an IAM list operation, concatenating
// the page items. Every list-type test must traverse all pages: during a
// full regression other services create resources concurrently, so a
// single page is never guaranteed to hold the entity under test.
func iamPaginate[T any](fetch func(marker *string) ([]T, *string, error)) ([]T, error) {
	return paginate(fetch)
}

// createUser creates a throwaway user and returns its cleanup.
func (tc *iamTestContext) createUser(name string) (func(), error) {
	if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(name)}); err != nil {
		return nil, err
	}
	return func() {
		tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(name)})
	}, nil
}

// createRole creates a throwaway role carrying the standard EC2 trust
// policy and returns its cleanup.
func (tc *iamTestContext) createRole(name string) (func(), error) {
	if _, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		AssumeRolePolicyDocument: aws.String(ec2AssumeRolePolicy),
	}); err != nil {
		return nil, err
	}
	return func() {
		tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
	}, nil
}

// createPolicy creates a customer-managed policy from doc and returns
// its ARN plus cleanup.
func (tc *iamTestContext) createPolicy(name, doc string) (string, func(), error) {
	resp, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(doc),
	})
	if err != nil {
		return "", nil, err
	}
	arn := aws.ToString(resp.Policy.Arn)
	return arn, func() {
		tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(arn)})
	}, nil
}
