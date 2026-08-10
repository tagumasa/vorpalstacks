package ec2

import (
	"context"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateSubnet creates a subnet in the specified VPC.
func (s *EC2Service) CreateSubnet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	mapPublicIpOnLaunch := false
	if v := request.GetStringParam(params, "MapPublicIpOnLaunch"); v == "true" {
		mapPublicIpOnLaunch = true
	}

	result, err := s.createSubnetCore(store, CreateSubnetInput{
		VpcId:               request.GetStringParam(params, "VpcId"),
		CidrBlock:           request.GetStringParam(params, "CidrBlock"),
		AvailabilityZone:    request.GetStringParam(params, "AvailabilityZone"),
		MapPublicIpOnLaunch: mapPublicIpOnLaunch,
		Tags:                parseTagsToCore(parseEC2Tags(params)),
		Region:              reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"Subnet": result.Subnet}, nil
}

// DescribeSubnets describes one or more subnets.
func (s *EC2Service) DescribeSubnets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subnetIDs := request.GetStringList(params, "SubnetId")
	filters, err := parseFilters(params)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetStringParam(params, "NextToken")
	maxResults := 0
	if mr := request.GetStringParam(params, "MaxResults"); mr != "" {
		if v, e := parseInt64Param(mr); e == nil {
			maxResults = int(v)
		}
	}

	result, err := s.describeSubnetsCore(store, subnetIDs, filters, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Subnets))
	for _, sn := range result.Subnets {
		items = append(items, sn)
	}
	resp := map[string]interface{}{
		"SubnetSet": protocol.XMLElements{ElementName: "item", Items: items},
	}
	if result.IsTruncated && result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
}

// DeleteSubnet deletes the specified subnet.
func (s *EC2Service) DeleteSubnet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSubnetCore(ctx, store, reqCtx.GetRegion(), request.GetStringParam(params, "SubnetId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{"return": true}, nil
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
