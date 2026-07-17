package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"vorpalstacks-sdk-tests/config"
)

type route53TestContext struct {
	client    *route53.Client
	ec2Client *ec2.Client
	ctx       context.Context
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
		uniq:      time.Now().UnixNano(),
	}

	var results []TestResult
	results = append(results, r.runRoute53ZoneTests(tc)...)
	results = append(results, r.runRoute53RecordTests(tc)...)
	results = append(results, r.runRoute53HealthCheckTests(tc)...)
	results = append(results, r.runRoute53VPCTests(tc)...)
	results = append(results, r.runRoute53TagTests(tc)...)
	results = append(results, r.runRoute53EdgeTests(tc)...)
	return results
}

func (tc *route53TestContext) domain(suffix string) string {
	return fmt.Sprintf("%s-%d.com.", suffix, tc.uniq)
}

func (tc *route53TestContext) callerRef(suffix string) string {
	return fmt.Sprintf("%s-%d", suffix, tc.uniq)
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
