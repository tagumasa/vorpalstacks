package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"vorpalstacks-sdk-tests/config"
)

type iotTestContext struct {
	client *iot.Client
	ctx    context.Context
	region string
}

func (r *TestRunner) newIoTTestContext() (*iotTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &iotTestContext{
		client: iot.NewFromConfig(cfg),
		ctx:    context.Background(),
		region: r.region,
	}, nil
}

// iotNameSeq provides a monotonic counter so that two uniqueName calls in the
// same nanosecond still produce distinct values (parallel goroutine safe).
var iotNameSeq uint64

// uniqueName returns a collision-free resource name for the given prefix. It
// combines a timestamp with an atomic counter so that parallel test runs,
// multiple regions and re-runs after partial failure never clash.
func uniqueName(prefix string) string {
	n := atomic.AddUint64(&iotNameSeq, 1)
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), n)
}

// arn builds an AWS-style ARN using the test region rather than a hard-coded
// us-east-1, so the IoT data plane stays correct under regional brokers.
func (tc *iotTestContext) arn(service, resourceType, name string) string {
	return fmt.Sprintf("arn:aws:%s:%s:000000000000:%s/%s", service, tc.region, resourceType, name)
}

// iamRoleARN returns a placeholder IAM role ARN for the test region.
func (tc *iotTestContext) iamRoleARN(suffix string) string {
	return fmt.Sprintf("arn:aws:iam::000000000000:role/%s", suffix)
}

// thingGroupExists paginates ListThingGroups and reports whether a group with
// the given name exists. Pagination is mandatory because full regressions run
// many services in parallel and the target group may sit beyond the first page.
func (tc *iotTestContext) thingGroupExists(groupName string) (bool, error) {
	var token *string
	for {
		out, err := tc.client.ListThingGroups(tc.ctx, &iot.ListThingGroupsInput{NextToken: token})
		if err != nil {
			return false, err
		}
		for _, g := range out.ThingGroups {
			if g.GroupName != nil && *g.GroupName == groupName {
				return true, nil
			}
		}
		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		token = out.NextToken
	}
	return false, nil
}

// billingGroupExists paginates ListBillingGroups for groupName.
func (tc *iotTestContext) billingGroupExists(groupName string) (bool, error) {
	var token *string
	for {
		out, err := tc.client.ListBillingGroups(tc.ctx, &iot.ListBillingGroupsInput{NextToken: token})
		if err != nil {
			return false, err
		}
		for _, g := range out.BillingGroups {
			if g.GroupName != nil && *g.GroupName == groupName {
				return true, nil
			}
		}
		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		token = out.NextToken
	}
	return false, nil
}

// thingInGroupExists paginates ListThingsInThingGroup for thingName.
func (tc *iotTestContext) thingInGroupExists(groupName, thingName string) (bool, error) {
	var token *string
	for {
		out, err := tc.client.ListThingsInThingGroup(tc.ctx, &iot.ListThingsInThingGroupInput{
			ThingGroupName: aws.String(groupName),
			NextToken:      token,
		})
		if err != nil {
			return false, err
		}
		for _, t := range out.Things {
			if t == thingName {
				return true, nil
			}
		}
		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		token = out.NextToken
	}
	return false, nil
}

// thingInBillingGroupExists paginates ListThingsInBillingGroup for thingName.
func (tc *iotTestContext) thingInBillingGroupExists(groupName, thingName string) (bool, error) {
	var token *string
	for {
		out, err := tc.client.ListThingsInBillingGroup(tc.ctx, &iot.ListThingsInBillingGroupInput{
			BillingGroupName: aws.String(groupName),
			NextToken:        token,
		})
		if err != nil {
			return false, err
		}
		for _, t := range out.Things {
			if t == thingName {
				return true, nil
			}
		}
		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		token = out.NextToken
	}
	return false, nil
}

// thingExists paginates ListThings for thingName.
func (tc *iotTestContext) thingExists(thingName string) (bool, error) {
	var token *string
	for {
		out, err := tc.client.ListThings(tc.ctx, &iot.ListThingsInput{NextToken: token})
		if err != nil {
			return false, err
		}
		for _, t := range out.Things {
			if t.ThingName != nil && *t.ThingName == thingName {
				return true, nil
			}
		}
		if out.NextToken == nil || *out.NextToken == "" {
			break
		}
		token = out.NextToken
	}
	return false, nil
}

// certificateExists paginates ListCertificates (marker-based) for the given id.
func (tc *iotTestContext) certificateExists(certID string) (bool, error) {
	var marker *string
	for {
		out, err := tc.client.ListCertificates(tc.ctx, &iot.ListCertificatesInput{Marker: marker})
		if err != nil {
			return false, err
		}
		for _, c := range out.Certificates {
			if c.CertificateId != nil && *c.CertificateId == certID {
				return true, nil
			}
		}
		if out.NextMarker == nil || *out.NextMarker == "" {
			break
		}
		marker = out.NextMarker
	}
	return false, nil
}

// awsHTTPStatus extracts the HTTP status code from an AWS SDK v2 error, or 0
// when the error carries no HTTP status information.
func awsHTTPStatus(err error) int {
	var httpErr interface{ HTTPStatusCode() int }
	if errors.As(err, &httpErr) {
		return httpErr.HTTPStatusCode()
	}
	return 0
}

// errorMatches returns true when err carries one of the substrings (case
// insensitive) or the given HTTP status code.
func errorMatches(err error, status int, codeSubstrings ...string) bool {
	if err == nil {
		return false
	}
	if status > 0 && awsHTTPStatus(err) == status {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, s := range codeSubstrings {
		if strings.Contains(msg, strings.ToLower(s)) {
			return true
		}
	}
	return false
}

// expectNotFound asserts that err represents a resource-not-found condition.
// It returns a descriptive failure when err is nil or does not look like a
// NotFound response.
func expectNotFound(err error) error {
	if err == nil {
		return fmt.Errorf("expected a NotFound error, got nil")
	}
	if errorMatches(err, 404, "ResourceNotFound", "NotFound", "not found") {
		return nil
	}
	return fmt.Errorf("expected a NotFound error, got: %w", err)
}

// expectConflict asserts that err represents a resource-already-exists /
// conflict condition.
func expectConflict(err error) error {
	if err == nil {
		return fmt.Errorf("expected a Conflict error, got nil")
	}
	if errorMatches(err, 409, "ResourceAlreadyExists", "AlreadyExists", "Conflict") {
		return nil
	}
	return fmt.Errorf("expected a Conflict error, got: %w", err)
}

// expectValidationError asserts that err represents a validation / bad
// request condition.
func expectValidationError(err error) error {
	if err == nil {
		return fmt.Errorf("expected a Validation error, got nil")
	}
	if errorMatches(err, 400, "Validation", "InvalidRequest", "InvalidParameter") {
		return nil
	}
	return fmt.Errorf("expected a Validation error, got: %w", err)
}

func (r *TestRunner) RunIoTTests() []TestResult {
	tc, err := r.newIoTTestContext()
	if err != nil {
		return []TestResult{{
			Service:  "iot",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	var results []TestResult
	results = append(results, r.runIoTThingTests(tc)...)
	results = append(results, r.runIoTThingTypeTests(tc)...)
	results = append(results, r.runIoTThingGroupTests(tc)...)
	results = append(results, r.runIoTBillingGroupTests(tc)...)
	results = append(results, r.runIoTCertTests(tc)...)
	results = append(results, r.runIoTPolicyTests(tc)...)
	results = append(results, r.runIoTJobTests(tc)...)
	results = append(results, r.runIoTAuthorizerTests(tc)...)
	results = append(results, r.runIoTSecurityProfileTests(tc)...)
	results = append(results, r.runIoTDomainConfigTests(tc)...)
	results = append(results, r.runIoTProvisioningTemplateTests(tc)...)
	results = append(results, r.runIoTCACertificateTests(tc)...)
	results = append(results, r.runIoTReadonlyListTests(tc)...)
	results = append(results, r.runIoTCRUDExtTests(tc)...)
	results = append(results, r.runIoTReadonlyExtTests(tc)...)
	results = append(results, r.runIoTMoreCRUDTests(tc)...)
	results = append(results, r.runIoTAuditStreamProvTests(tc)...)
	results = append(results, r.runIoTFinalBatchTests(tc)...)
	results = append(results, r.runIoTRemainingSmokeTests(tc)...)
	results = append(results, r.runIoTUpdateAndTailTests(tc)...)
	results = append(results, r.runIoTLastBatchTests(tc)...)
	results = append(results, r.runIoTAuditSuppressionTests(tc)...)
	results = append(results, r.runIoTFinalRemainingTests(tc)...)
	results = append(results, r.runIoTShadowTests(tc)...)
	results = append(results, r.runIoTStreamRegistrationTests(tc)...)
	results = append(results, r.runIoTPackageCommandProviderTests(tc)...)
	results = append(results, r.runIoTCoverageGapTests(tc)...)
	results = append(results, r.runIoTIntegrationRuleActionTests(tc)...)
	results = append(results, r.runIoTMQTTTests(tc)...)
	return results
}
