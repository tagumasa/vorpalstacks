package ec2

import (
	"context"
	"strings"

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
	if err := validateFilterNames(filters, allowedSubnetFilters); err != nil {
		return nil, err
	}
	nextToken, maxResults, err := parsePaginationParams(params, "SubnetId")
	if err != nil {
		return nil, err
	}

	result, err := s.describeSubnetsCore(store, subnetIDs, filters, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Subnets))
	for _, sn := range result.Subnets {
		items = append(items, subnetToXMLMap(sn))
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
// Filter names are validated by validateFilterNames before matching.
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
		case "cidr-block", "cidr", "cidrBlock":
			if !anyMatch(f.Values, sn.CidrBlock) {
				return false
			}
		case "state":
			if !anyMatch(f.Values, sn.State) {
				return false
			}
		case "availability-zone", "availabilityZone":
			if !anyMatch(f.Values, sn.AvailabilityZone) {
				return false
			}
		case "availability-zone-id", "availabilityZoneId":
			if !anyMatch(f.Values, sn.AvailabilityZoneId) {
				return false
			}
		case "available-ip-address-count":
			if !anyMatchInt64(f.Values, sn.AvailableIpAddressCount) {
				return false
			}
		case "default-for-az", "defaultForAz":
			if !anyMatchBool(f.Values, sn.DefaultForAz) {
				return false
			}
		case "map-public-ip-on-launch":
			if !anyMatchBool(f.Values, sn.MapPublicIpOnLaunch) {
				return false
			}
		case "enable-dns64":
			if !anyMatchBool(f.Values, sn.EnableDns64) {
				return false
			}
		case "map-customer-owned-ip-on-launch":
			// Subnets in this implementation never map customer-owned
			// addresses; the attribute reports its default value.
			if !anyMatchBool(f.Values, false) {
				return false
			}
		case "enable-lni-at-device-index", "customer-owned-ipv4-pool", "outpost-arn",
			"ipv6-cidr-block-association.ipv6-cidr-block",
			"ipv6-cidr-block-association.association-id",
			"ipv6-cidr-block-association.state":
			// Attributes not carried by this implementation (local network
			// interfaces, customer-owned pools, Outposts, associated IPv6
			// CIDR blocks); such subnets match none of these filter values.
			return false
		case "private-dns-name-options-on-launch.hostname-type":
			// IPv4-only subnets assign instance hostnames from the IPv4
			// address by default.
			if !anyMatch(f.Values, "ip-name") {
				return false
			}
		case "private-dns-name-options-on-launch.enable-resource-name-dns-a-record",
			"private-dns-name-options-on-launch.enable-resource-name-dns-aaaa-record":
			if !anyMatchBool(f.Values, false) {
				return false
			}
		case "owner-id":
			if !anyMatch(f.Values, sn.OwnerId) {
				return false
			}
		case "subnet-arn":
			if !anyMatch(f.Values, sn.SubnetArn) {
				return false
			}
		case "ipv6-native":
			if !anyMatchBool(f.Values, sn.Ipv6Native) {
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
		default:
			// AWS tag key/value filters use the "tag:<key>" name form; the
			// filter matches when the resource carries that tag key with
			// any of the filter values.
			if strings.HasPrefix(f.Name, "tag:") {
				if !hasTagKeyValues(sn.Tags, strings.TrimPrefix(f.Name, "tag:"), f.Values) {
					return false
				}
			}
		}
	}
	return true
}
