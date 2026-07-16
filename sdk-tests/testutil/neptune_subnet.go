package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneSubnetGroupTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "SetupVPC", func() error {
		vpcResp, err := tc.ec2User.CreateVpc(tc.ctx, &ec2.CreateVpcInput{
			CidrBlock: aws.String("10.0.0.0/16"),
		})
		if err != nil {
			return fmt.Errorf("failed to create VPC: %v", err)
		}
		if vpcResp.Vpc == nil || vpcResp.Vpc.VpcId == nil {
			return fmt.Errorf("expected Vpc with VpcId in CreateVpc response")
		}
		tc.vpcID = vpcResp.Vpc.VpcId

		for i, cidr := range []string{"10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"} {
			az := fmt.Sprintf("%s%c", tc.region, 'a'+byte(i))
			subResp, err := tc.ec2User.CreateSubnet(tc.ctx, &ec2.CreateSubnetInput{
				VpcId:            tc.vpcID,
				CidrBlock:        aws.String(cidr),
				AvailabilityZone: aws.String(az),
			})
			if err != nil {
				return fmt.Errorf("failed to create subnet %s: %v", cidr, err)
			}
			if subResp.Subnet == nil || subResp.Subnet.SubnetId == nil {
				return fmt.Errorf("expected Subnet with SubnetId for %s", cidr)
			}
			tc.allSubnetIds = append(tc.allSubnetIds, *subResp.Subnet.SubnetId)
		}
		tc.subnetIds = tc.allSubnetIds[:2]
		return nil
	}))

	results = append(results, r.RunTest("neptune", "CreateDBSubnetGroup", func() error {
		resp, err := tc.client.CreateDBSubnetGroup(tc.ctx, &neptune.CreateDBSubnetGroupInput{
			DBSubnetGroupName:        aws.String(tc.subnetGroupName),
			DBSubnetGroupDescription: aws.String("Test subnet group"),
			SubnetIds:                tc.subnetIds,
		})
		if err != nil {
			return err
		}
		if resp.DBSubnetGroup == nil {
			return fmt.Errorf("expected DBSubnetGroup in response")
		}
		sg := resp.DBSubnetGroup
		if sg.DBSubnetGroupName == nil || *sg.DBSubnetGroupName != tc.subnetGroupName {
			return fmt.Errorf("expected DBSubnetGroupName=%s, got %v", tc.subnetGroupName, sg.DBSubnetGroupName)
		}
		if sg.DBSubnetGroupDescription == nil || *sg.DBSubnetGroupDescription != "Test subnet group" {
			return fmt.Errorf("expected description 'Test subnet group', got %v", sg.DBSubnetGroupDescription)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBSubnetGroups", func() error {
		resp, err := tc.client.DescribeDBSubnetGroups(tc.ctx, &neptune.DescribeDBSubnetGroupsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, sg := range resp.DBSubnetGroups {
			if sg.DBSubnetGroupName != nil && *sg.DBSubnetGroupName == tc.subnetGroupName {
				found = true
				if sg.VpcId == nil || *sg.VpcId == "" {
					return fmt.Errorf("expected non-empty VpcId on subnet group")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created subnet group not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "ModifyDBSubnetGroup", func() error {
		resp, err := tc.client.ModifyDBSubnetGroup(tc.ctx, &neptune.ModifyDBSubnetGroupInput{
			DBSubnetGroupName:        aws.String(tc.subnetGroupName),
			DBSubnetGroupDescription: aws.String("Modified test subnet group"),
			SubnetIds:                tc.allSubnetIds,
		})
		if err != nil {
			return err
		}
		if resp.DBSubnetGroup == nil {
			return fmt.Errorf("expected DBSubnetGroup in ModifyDBSubnetGroup response")
		}
		if resp.DBSubnetGroup.DBSubnetGroupDescription == nil || *resp.DBSubnetGroup.DBSubnetGroupDescription != "Modified test subnet group" {
			return fmt.Errorf("expected description 'Modified test subnet group' in ModifyDBSubnetGroup response, got %v", resp.DBSubnetGroup.DBSubnetGroupDescription)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeDBSubnetGroups_FilterByName", func() error {
		resp, err := tc.client.DescribeDBSubnetGroups(tc.ctx, &neptune.DescribeDBSubnetGroupsInput{
			DBSubnetGroupName: aws.String(tc.subnetGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.DBSubnetGroups) != 1 {
			return fmt.Errorf("expected 1 subnet group, got %d", len(resp.DBSubnetGroups))
		}
		sg := resp.DBSubnetGroups[0]
		if sg.DBSubnetGroupDescription == nil || *sg.DBSubnetGroupDescription != "Modified test subnet group" {
			return fmt.Errorf("expected description 'Modified test subnet group', got %v", sg.DBSubnetGroupDescription)
		}
		return nil
	}))

	return results
}
