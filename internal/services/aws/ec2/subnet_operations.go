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
	vpcID := request.GetStringParam(params, "VpcId")
	if vpcID == "" {
		return nil, fmt.Errorf("ec2: VpcId is required")
	}
	cidrBlock := request.GetStringParam(params, "CidrBlock")
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

	az := request.GetStringParam(params, "AvailabilityZone")
	if az == "" {
		az = reqCtx.GetRegion() + "a"
	}

	mapPublicIpOnLaunch := false
	if v := request.GetStringParam(params, "MapPublicIpOnLaunch"); v == "true" {
		mapPublicIpOnLaunch = true
	}

	subnet := &ec2store.Subnet{
		SubnetId:            subnetID,
		VpcId:               vpcID,
		CidrBlock:           cidrBlock,
		AvailabilityZone:    az,
		State:               "available",
		OwnerId:             s.accountID,
		MapPublicIpOnLaunch: mapPublicIpOnLaunch,
		Tags:                parseEC2Tags(params),
	}

	if err := store.CreateSubnet(subnet); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"Subnet": subnet,
	}, nil
}

// DescribeSubnets describes one or more subnets. Supports SubnetId for single
// lookup and Filter.N for filtering by vpc-id, subnet-id, cidr-block, state,
// availability-zone, tag, etc.
func (s *EC2Service) DescribeSubnets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subnetID := request.GetStringParam(params, "SubnetId")
	if subnetID != "" {
		subnet, err := store.GetSubnet(subnetID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"SubnetSet": protocol.XMLElements{ElementName: "item", Items: []interface{}{subnet}},
		}, nil
	}

	subnets, err := store.ListSubnets()
	if err != nil {
		return nil, err
	}

	filters := parseFilters(params)
	items := make([]interface{}, 0, len(subnets))
	for _, sn := range subnets {
		if matchesSubnetFilters(sn, filters) {
			items = append(items, sn)
		}
	}
	return map[string]interface{}{
		"SubnetSet": protocol.XMLElements{ElementName: "item", Items: items},
	}, nil
}

// DeleteSubnet deletes the specified subnet.
func (s *EC2Service) DeleteSubnet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	subnetID := request.GetStringParam(params, "SubnetId")
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

// matchesSubnetFilters checks if a subnet matches all the given filters.
func matchesSubnetFilters(sn *ec2store.Subnet, filters []ec2Filter) bool {
	for _, f := range filters {
		switch f.Name {
		case "vpc-id":
			if !anyMatch(f.Values, sn.VpcId) {
				return false
			}
		case "subnet-id":
			if !anyMatch(f.Values, sn.SubnetId) {
				return false
			}
		case "cidr-block":
			if !anyMatch(f.Values, sn.CidrBlock) {
				return false
			}
		case "state":
			if !anyMatch(f.Values, sn.State) {
				return false
			}
		case "availability-zone":
			if !anyMatch(f.Values, sn.AvailabilityZone) {
				return false
			}
		case "tag-key":
			if !hasTagKey(sn.Tags, f.Values) {
				return false
			}
		case "tag-value":
			if !hasTagValue(sn.Tags, f.Values) {
				return false
			}
		case "tag":
			if !hasTagKeyValue(sn.Tags, f.Values) {
				return false
			}
		}
	}
	return true
}
