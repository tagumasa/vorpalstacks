package ec2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateSubnet creates a subnet in the specified VPC.
func (s *EC2Service) CreateSubnet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	vpcID := getStringParam(params, "VpcId")
	if vpcID == "" {
		return nil, fmt.Errorf("ec2: VpcId is required")
	}
	cidrBlock := getStringParam(params, "CidrBlock")
	if cidrBlock == "" {
		return nil, fmt.Errorf("ec2: CidrBlock is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetVPC(vpcID); err != nil {
		return nil, translateStoreError(err)
	}

	subnetID, err := GenerateSubnetID()
	if err != nil {
		return nil, err
	}

	az := getStringParam(params, "AvailabilityZone")
	if az == "" {
		az = reqCtx.GetRegion() + "a"
	}

	mapPublicIpOnLaunch := false
	if v := getStringParam(params, "MapPublicIpOnLaunch"); v == "true" {
		mapPublicIpOnLaunch = true
	}

	tags := ParseTags(params)

	subnet := &ec2store.Subnet{
		SubnetId:            subnetID,
		VpcId:               vpcID,
		CidrBlock:           cidrBlock,
		AvailabilityZone:    az,
		State:               "available",
		OwnerId:             s.accountID,
		MapPublicIpOnLaunch: mapPublicIpOnLaunch,
		Tags:                tags,
	}

	if err := store.CreateSubnet(subnet); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"Subnet": subnet,
	}, nil
}

// DescribeSubnets describes one or all subnets.
func (s *EC2Service) DescribeSubnets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subnetID := getStringParam(params, "SubnetId")
	if subnetID != "" {
		subnet, err := store.GetSubnet(subnetID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"SubnetSet": protocol.XMLElements{ElementName: "item", Items: []interface{}{subnet}},
		}, nil
	}

	filterVpcID := getStringParam(params, "Filter.1.Value")

	subnets, err := store.ListSubnets()
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(subnets))
	for _, sn := range subnets {
		if filterVpcID != "" && sn.VpcId != filterVpcID {
			continue
		}
		items = append(items, sn)
	}
	return map[string]interface{}{
		"SubnetSet": protocol.XMLElements{ElementName: "item", Items: items},
	}, nil
}

// DeleteSubnet deletes the specified subnet.
func (s *EC2Service) DeleteSubnet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	subnetID := getStringParam(params, "SubnetId")
	if subnetID == "" {
		return nil, fmt.Errorf("ec2: SubnetId is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteSubnet(subnetID); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}
