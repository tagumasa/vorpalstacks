package ec2

import (
	"context"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateVpc creates a VPC with the specified CIDR block.
func (s *EC2Service) CreateVpc(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	cidrBlock := request.GetStringParam(params, "CidrBlock")
	if cidrBlock == "" {
		cidrBlock = "172.31.0.0/16"
	}

	vpcID, err := GenerateVpcID()
	if err != nil {
		return nil, err
	}

	instanceTenancy := request.GetStringParam(params, "InstanceTenancy")
	if instanceTenancy == "" {
		instanceTenancy = "default"
	}

	enableDnsSupport := true
	if v := request.GetStringParam(params, "EnableDnsSupport"); v == "false" {
		enableDnsSupport = false
	}
	enableDnsHostnames := false
	if v := request.GetStringParam(params, "EnableDnsHostnames"); v == "true" {
		enableDnsHostnames = true
	}

	vpc := &ec2store.VPC{
		VpcId:              vpcID,
		CidrBlock:          cidrBlock,
		State:              "available",
		OwnerId:            s.accountID,
		InstanceTenancy:    instanceTenancy,
		EnableDnsSupport:   enableDnsSupport,
		EnableDnsHostnames: enableDnsHostnames,
		Tags:               parseEC2Tags(params),
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

// DescribeVpcs describes one or more VPCs. Supports VpcId for single lookup
// and Filter.N for filtering by vpc-id, cidr, state, tag, etc.
func (s *EC2Service) DescribeVpcs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	vpcIDs := request.GetStringList(params, "VpcId")
	if len(vpcIDs) > 0 {
		items := make([]interface{}, 0, len(vpcIDs))
		for _, id := range vpcIDs {
			vpc, err := store.GetVPC(id)
			if err != nil {
				return nil, translateStoreError(err)
			}
			items = append(items, vpc)
		}
		return map[string]interface{}{
			"VpcSet": protocol.XMLElements{ElementName: "item", Items: items},
		}, nil
	}

	vpcs, err := store.ListVPCs()
	if err != nil {
		return nil, err
	}

	filters := parseFilters(params)
	items := make([]interface{}, 0, len(vpcs))
	for _, v := range vpcs {
		if matchesVPCFilters(v, filters) {
			items = append(items, v)
		}
	}
	return map[string]interface{}{
		"VpcSet": protocol.XMLElements{ElementName: "item", Items: items},
	}, nil
}

// DeleteVpc deletes the specified VPC. Returns DependencyViolation if the VPC
// still has dependent resources (subnets or security groups).
func (s *EC2Service) DeleteVpc(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	vpcID := request.GetStringParam(params, "VpcId")
	if vpcID == "" {
		return nil, awserrors.NewMissingParameter("VpcId is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subnets, err := store.ListSubnetsByVPC(vpcID)
	if err != nil {
		return nil, err
	}
	if len(subnets) > 0 {
		return nil, awserrors.NewAWSError("DependencyViolation",
			"The vpc '"+vpcID+"' has dependencies and cannot be deleted",
			http.StatusBadRequest)
	}

	sgs, err := store.ListSecurityGroupsByVPC(vpcID)
	if err != nil {
		return nil, err
	}
	if len(sgs) > 0 {
		return nil, awserrors.NewAWSError("DependencyViolation",
			"The vpc '"+vpcID+"' has dependencies and cannot be deleted",
			http.StatusBadRequest)
	}

	if err := store.DeleteVPC(vpcID); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// matchesVPCFilters checks if a VPC matches all the given filters.
func matchesVPCFilters(vpc *ec2store.VPC, filters []ec2Filter) bool {
	for _, f := range filters {
		switch f.Name {
		case "vpc-id":
			if !anyMatch(f.Values, vpc.VpcId) {
				return false
			}
		case "cidr":
			if !anyMatch(f.Values, vpc.CidrBlock) {
				return false
			}
		case "state":
			if !anyMatch(f.Values, vpc.State) {
				return false
			}
		case "tag-key":
			if !hasTagKey(vpc.Tags, f.Values) {
				return false
			}
		case "tag-value":
			if !hasTagValue(vpc.Tags, f.Values) {
				return false
			}
		case "tag":
			if !hasTagKeyValue(vpc.Tags, f.Values) {
				return false
			}
		}
	}
	return true
}
