package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"vorpalstacks-sdk-tests/config"
)

type route53TestContext struct {
	client    *route53.Client
	ec2Client *ec2.Client
	ctx       context.Context
	region    string
	uniq      int64
}

func (r *TestRunner) RunRoute53Tests() []TestResult {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{
			Service:  "route53",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}

	tc := &route53TestContext{
		client:    route53.NewFromConfig(cfg),
		ec2Client: ec2.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		uniq:      time.Now().UnixNano(),
	}

	var results []TestResult
	results = append(results, r.runRoute53ZoneTests(tc)...)
	results = append(results, r.runRoute53RecordTests(tc)...)
	results = append(results, r.runRoute53HealthCheckTests(tc)...)
	results = append(results, r.runRoute53VPCTests(tc)...)
	results = append(results, r.runRoute53TagTests(tc)...)
	results = append(results, r.runRoute53CidrCollectionTests(tc)...)
	results = append(results, r.runRoute53EdgeTests(tc)...)
	return results
}

func (tc *route53TestContext) domain(suffix string) string {
	return fmt.Sprintf("%s-%d.com.", suffix, tc.uniq)
}

func (tc *route53TestContext) callerRef(suffix string) string {
	return fmt.Sprintf("%s-%d", suffix, tc.uniq)
}

// uniqName returns a per-run unique resource name.
func (tc *route53TestContext) uniqName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, tc.uniq)
}

// createZone creates a public hosted zone with the given name and caller
// reference.
func (tc *route53TestContext) createZone(name, ref string) (*route53.CreateHostedZoneOutput, error) {
	return tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
		Name:            aws.String(name),
		CallerReference: aws.String(ref),
	})
}

// createPrivateZone creates a private hosted zone attached to vpcID; an
// empty comment omits the Comment member.
func (tc *route53TestContext) createPrivateZone(name, ref, vpcID, comment string) (*route53.CreateHostedZoneOutput, error) {
	config := &types.HostedZoneConfig{PrivateZone: true}
	if comment != "" {
		config.Comment = aws.String(comment)
	}
	return tc.client.CreateHostedZone(tc.ctx, &route53.CreateHostedZoneInput{
		Name:             aws.String(name),
		CallerReference:  aws.String(ref),
		HostedZoneConfig: config,
		VPC: &types.VPC{
			VPCId:     aws.String(vpcID),
			VPCRegion: types.VPCRegion(tc.region),
		},
	})
}

func (tc *route53TestContext) getZone(id string) (*route53.GetHostedZoneOutput, error) {
	return tc.client.GetHostedZone(tc.ctx, &route53.GetHostedZoneInput{Id: aws.String(id)})
}

func (tc *route53TestContext) deleteZone(id string) {
	tc.client.DeleteHostedZone(tc.ctx, &route53.DeleteHostedZoneInput{Id: aws.String(id)})
}

// changeRecords submits a change batch of record-set changes against the
// hosted zone.
func (tc *route53TestContext) changeRecords(zoneID string, changes ...types.Change) (*route53.ChangeResourceRecordSetsOutput, error) {
	return tc.client.ChangeResourceRecordSets(tc.ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch:  &types.ChangeBatch{Changes: changes},
	})
}

// rrChange builds a single change applying action to a record set with the
// given values; no values yields the empty ResourceRecords list.
func rrChange(action types.ChangeAction, name string, rrType types.RRType, ttl int64, values ...string) types.Change {
	records := make([]types.ResourceRecord, 0, len(values))
	for _, v := range values {
		records = append(records, types.ResourceRecord{Value: aws.String(v)})
	}
	return types.Change{
		Action: action,
		ResourceRecordSet: &types.ResourceRecordSet{
			Name:            aws.String(name),
			Type:            rrType,
			TTL:             aws.Int64(ttl),
			ResourceRecords: records,
		},
	}
}

func (tc *route53TestContext) listRecords(zoneID string) (*route53.ListResourceRecordSetsOutput, error) {
	return tc.client.ListResourceRecordSets(tc.ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	})
}

func (tc *route53TestContext) createHealthCheck(ref string, config *types.HealthCheckConfig) (*route53.CreateHealthCheckOutput, error) {
	return tc.client.CreateHealthCheck(tc.ctx, &route53.CreateHealthCheckInput{
		CallerReference:   aws.String(ref),
		HealthCheckConfig: config,
	})
}

// tcpHealthCheck builds a TCP health check config targeting fqdn:port.
func tcpHealthCheck(fqdn string, port int32) *types.HealthCheckConfig {
	return &types.HealthCheckConfig{
		Type:                     types.HealthCheckTypeTcp,
		FullyQualifiedDomainName: aws.String(fqdn),
		Port:                     aws.Int32(port),
	}
}

func (tc *route53TestContext) getHealthCheck(id string) (*route53.GetHealthCheckOutput, error) {
	return tc.client.GetHealthCheck(tc.ctx, &route53.GetHealthCheckInput{HealthCheckId: aws.String(id)})
}

func (tc *route53TestContext) deleteHealthCheck(id string) {
	tc.client.DeleteHealthCheck(tc.ctx, &route53.DeleteHealthCheckInput{HealthCheckId: aws.String(id)})
}

func (tc *route53TestContext) createCidrCollection(name, ref string) (*route53.CreateCidrCollectionOutput, error) {
	return tc.client.CreateCidrCollection(tc.ctx, &route53.CreateCidrCollectionInput{
		Name:            aws.String(name),
		CallerReference: aws.String(ref),
	})
}

func (tc *route53TestContext) deleteCidrCollection(id string) {
	tc.client.DeleteCidrCollection(tc.ctx, &route53.DeleteCidrCollectionInput{Id: aws.String(id)})
}

// expectRoute53Error asserts that err carries the given error code substring
// and HTTP status.
func expectRoute53Error(err error, code string, status int) error {
	if err == nil {
		return fmt.Errorf("expected %s error, got nil", code)
	}
	if cerr := AssertErrorContains(err, code); cerr != nil {
		return cerr
	}
	if got := awsHTTPStatus(err); got != status {
		return fmt.Errorf("expected HTTP %d, got %d", status, got)
	}
	return nil
}

// createTestVPC creates a VPC in EC2 and returns its ID. The VPC should be
// deleted by the caller via defer.
func (tc *route53TestContext) createTestVPC(cidr string) (string, error) {
	resp, err := tc.ec2Client.CreateVpc(tc.ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidr),
	})
	if err != nil {
		return "", fmt.Errorf("create VPC: %w", err)
	}
	if resp.Vpc == nil || resp.Vpc.VpcId == nil {
		return "", fmt.Errorf("create VPC: nil response")
	}
	return *resp.Vpc.VpcId, nil
}

func (tc *route53TestContext) deleteTestVPC(vpcID string) {
	tc.ec2Client.DeleteVpc(tc.ctx, &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcID),
	})
}
