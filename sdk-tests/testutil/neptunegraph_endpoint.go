package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) runNeptunegraphEndpointTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult

	var endpointVpcID string
	var endpointSubnetIDs []string

	results = append(results, r.RunTest("neptunegraph", "SetupVpcForEndpoint", func() error {
		ec2Cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("failed to load EC2 config: %v", err)
		}
		ec2Client := ec2.NewFromConfig(ec2Cfg)

		vpcResp, err := ec2Client.CreateVpc(tc.ctx, &ec2.CreateVpcInput{
			CidrBlock: aws.String("10.99.0.0/16"),
		})
		if err != nil {
			return fmt.Errorf("failed to create VPC: %v", err)
		}
		endpointVpcID = *vpcResp.Vpc.VpcId

		for _, cidr := range []string{"10.99.1.0/24", "10.99.2.0/24"} {
			subResp, err := ec2Client.CreateSubnet(tc.ctx, &ec2.CreateSubnetInput{
				VpcId:            aws.String(endpointVpcID),
				CidrBlock:        aws.String(cidr),
				AvailabilityZone: aws.String(r.region + "a"),
			})
			if err != nil {
				return fmt.Errorf("failed to create subnet %s: %v", cidr, err)
			}
			endpointSubnetIDs = append(endpointSubnetIDs, *subResp.Subnet.SubnetId)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "CreatePrivateGraphEndpoint", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if endpointVpcID == "" || len(endpointSubnetIDs) == 0 {
			return fmt.Errorf("VPC/subnet setup failed")
		}
		resp, err := tc.client.CreatePrivateGraphEndpoint(tc.ctx, &neptunegraph.CreatePrivateGraphEndpointInput{
			GraphIdentifier: aws.String(tc.graphID),
			VpcId:           aws.String(endpointVpcID),
			SubnetIds:       endpointSubnetIDs,
		})
		if err != nil {
			return err
		}
		if resp.VpcId == nil || *resp.VpcId != endpointVpcID {
			return fmt.Errorf("expected vpcId=%s, got %v", endpointVpcID, resp.VpcId)
		}
		if resp.Status == "" {
			return fmt.Errorf("expected non-nil endpoint status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetPrivateGraphEndpoint", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if endpointVpcID == "" {
			return fmt.Errorf("VPC setup failed")
		}
		resp, err := tc.client.GetPrivateGraphEndpoint(tc.ctx, &neptunegraph.GetPrivateGraphEndpointInput{
			GraphIdentifier: aws.String(tc.graphID),
			VpcId:           aws.String(endpointVpcID),
		})
		if err != nil {
			return err
		}
		if resp.VpcId == nil || *resp.VpcId != endpointVpcID {
			return fmt.Errorf("expected vpcId=%s, got %v", endpointVpcID, resp.VpcId)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListPrivateGraphEndpoints", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.ListPrivateGraphEndpoints(tc.ctx, &neptunegraph.ListPrivateGraphEndpointsInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
		if err != nil {
			return err
		}
		if resp.PrivateGraphEndpoints == nil {
			return fmt.Errorf("expected non-nil PrivateGraphEndpoints list")
		}
		if len(resp.PrivateGraphEndpoints) == 0 {
			return fmt.Errorf("expected at least one private endpoint")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "DeletePrivateGraphEndpoint", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if endpointVpcID == "" {
			return fmt.Errorf("VPC setup failed")
		}
		resp, err := tc.client.DeletePrivateGraphEndpoint(tc.ctx, &neptunegraph.DeletePrivateGraphEndpointInput{
			GraphIdentifier: aws.String(tc.graphID),
			VpcId:           aws.String(endpointVpcID),
		})
		if err != nil {
			return err
		}
		if resp.VpcId == nil || *resp.VpcId != endpointVpcID {
			return fmt.Errorf("expected vpcId=%s, got %v", endpointVpcID, resp.VpcId)
		}
		return nil
	}))

	return results
}
