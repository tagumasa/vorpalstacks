package ec2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateVpc creates a VPC with the specified CIDR block.
func (s *EC2Service) CreateVpc(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	cidrBlock := getStringParam(params, "CidrBlock")
	if cidrBlock == "" {
		cidrBlock = "172.31.0.0/16"
	}

	vpcID, err := GenerateVpcID()
	if err != nil {
		return nil, err
	}

	instanceTenancy := getStringParam(params, "InstanceTenancy")
	if instanceTenancy == "" {
		instanceTenancy = "default"
	}

	enableDnsSupport := true
	if v := getStringParam(params, "EnableDnsSupport"); v == "false" {
		enableDnsSupport = false
	}
	enableDnsHostnames := false
	if v := getStringParam(params, "EnableDnsHostnames"); v == "true" {
		enableDnsHostnames = true
	}

	tags := ParseTags(params)

	vpc := &ec2store.VPC{
		VpcId:              vpcID,
		CidrBlock:          cidrBlock,
		State:              "available",
		OwnerId:            s.accountID,
		InstanceTenancy:    instanceTenancy,
		EnableDnsSupport:   enableDnsSupport,
		EnableDnsHostnames: enableDnsHostnames,
		Tags:               tags,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateVPC(vpc); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"Vpc": vpc,
	}, nil
}

// DescribeVpcs describes one or all VPCs.
func (s *EC2Service) DescribeVpcs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	vpcID := getStringParam(params, "VpcId")
	if vpcID != "" {
		vpc, err := store.GetVPC(vpcID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"VpcSet": protocol.XMLElements{ElementName: "item", Items: []interface{}{vpc}},
		}, nil
	}

	vpcs, err := store.ListVPCs()
	if err != nil {
		return nil, err
	}
	items := make([]interface{}, 0, len(vpcs))
	for _, v := range vpcs {
		items = append(items, v)
	}
	return map[string]interface{}{
		"VpcSet": protocol.XMLElements{ElementName: "item", Items: items},
	}, nil
}

// DeleteVpc deletes the specified VPC.
func (s *EC2Service) DeleteVpc(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	vpcID := getStringParam(params, "VpcId")
	if vpcID == "" {
		return nil, fmt.Errorf("ec2: VpcId is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteVPC(vpcID); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}
