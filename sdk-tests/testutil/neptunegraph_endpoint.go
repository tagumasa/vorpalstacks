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

	endpointVpcID, endpointSubnetIDs, err := r.setupEndpointVPC(tc)
	if err != nil {
		return append(results, TestResult{
			Service:  "neptunegraph",
			TestName: "SetupVpcForEndpoint",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}

	results = append(results, r.RunTest("neptunegraph", "CreatePrivateGraphEndpoint", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
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
		if err := tc.requireGraph(); err != nil {
			return err
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
		if err := tc.requireGraph(); err != nil {
			return err
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
		if err := tc.requireGraph(); err != nil {
			return err
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

// setupEndpointVPC provisions a throwaway VPC with two subnets used by the
// private-endpoint tests. It runs at section start and its failure surfaces
// as the SetupVpcForEndpoint FAIL row.
func (r *TestRunner) setupEndpointVPC(tc *neptunegraphContext) (string, []string, error) {
	ec2Cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to load EC2 config: %v", err)
	}
	ec2Client := ec2.NewFromConfig(ec2Cfg)

	vpcResp, err := ec2Client.CreateVpc(tc.ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.99.0.0/16"),
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to create VPC: %v", err)
	}
	if vpcResp.Vpc == nil || vpcResp.Vpc.VpcId == nil {
		return "", nil, fmt.Errorf("expected Vpc with VpcId in CreateVpc response")
	}
	vpcID := *vpcResp.Vpc.VpcId

	var subnetIDs []string
	for _, cidr := range []string{"10.99.1.0/24", "10.99.2.0/24"} {
		subResp, err := ec2Client.CreateSubnet(tc.ctx, &ec2.CreateSubnetInput{
			VpcId:            aws.String(vpcID),
			CidrBlock:        aws.String(cidr),
			AvailabilityZone: aws.String(r.region + "a"),
		})
		if err != nil {
			return "", nil, fmt.Errorf("failed to create subnet %s: %v", cidr, err)
		}
		if subResp.Subnet == nil || subResp.Subnet.SubnetId == nil {
			return "", nil, fmt.Errorf("expected Subnet with SubnetId for %s", cidr)
		}
		subnetIDs = append(subnetIDs, *subResp.Subnet.SubnetId)
	}
	return vpcID, subnetIDs, nil
}
