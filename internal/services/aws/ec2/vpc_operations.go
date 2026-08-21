package ec2

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateVpc creates a VPC with the specified CIDR block.
func (s *EC2Service) CreateVpc(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}

	enableDnsSupport := true
	if v := request.GetStringParam(params, "EnableDnsSupport"); v == "false" {
		enableDnsSupport = false
	}
	enableDnsHostnames := false
	if v := request.GetStringParam(params, "EnableDnsHostnames"); v == "true" {
		enableDnsHostnames = true
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createVpcCore(store, CreateVpcInput{
		CidrBlock:          request.GetStringParam(params, "CidrBlock"),
		InstanceTenancy:    request.GetStringParam(params, "InstanceTenancy"),
		EnableDnsSupport:   &enableDnsSupport,
		EnableDnsHostnames: &enableDnsHostnames,
		Tags:               parseTagsToCore(parseEC2Tags(params)),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"Vpc": vpcToXMLMap(result.Vpc)}, nil
}

// DescribeVpcs describes one or more VPCs.
func (s *EC2Service) DescribeVpcs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	vpcIDs := request.GetStringList(params, "VpcId")
	filters, err := parseFilters(params)
	if err != nil {
		return nil, err
	}
	if err := validateFilterNames(filters, allowedVPCFilters); err != nil {
		return nil, err
	}
	nextToken, maxResults, err := parsePaginationParams(params, "VpcId")
	if err != nil {
		return nil, err
	}

	result, err := s.describeVpcsCore(store, vpcIDs, filters, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Vpcs))
	for _, v := range result.Vpcs {
		items = append(items, vpcToXMLMap(v))
	}
	resp := map[string]interface{}{
		"VpcSet": protocol.XMLElements{ElementName: "item", Items: items},
	}
	if result.IsTruncated && result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
}

// DeleteVpc deletes the specified VPC.
func (s *EC2Service) DeleteVpc(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteVpcCore(store, request.GetStringParam(params, "VpcId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{"return": true}, nil
}

// matchesVPCFilters checks if a VPC matches all the given filters.
// Filter names are validated by validateFilterNames before matching.
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
		case "cidr-block-association.cidr-block", "cidr-block-association.association-id", "cidr-block-association.state":
			if !vpcMatchesCidrAssociationFilter(vpc, f) {
				return false
			}
		case "ipv6-cidr-block-association.ipv6-cidr-block",
			"ipv6-cidr-block-association.ipv6-pool",
			"ipv6-cidr-block-association.association-id",
			"ipv6-cidr-block-association.state":
			// This implementation associates no IPv6 CIDR blocks with VPCs,
			// so no VPC matches these documented filters.
			return false
		case "dhcp-options-id":
			if !anyMatch(f.Values, vpc.DhcpOptionsId) {
				return false
			}
		case "owner-id":
			if !anyMatch(f.Values, vpc.OwnerId) {
				return false
			}
		case "is-default":
			if !anyMatchBool(f.Values, vpc.IsDefault) {
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
		default:
			// AWS tag key/value filters use the "tag:<key>" name form; the
			// filter matches when the resource carries that tag key with
			// any of the filter values.
			if strings.HasPrefix(f.Name, "tag:") {
				if !hasTagKeyValues(vpc.Tags, strings.TrimPrefix(f.Name, "tag:"), f.Values) {
					return false
				}
			}
		}
	}
	return true
}

// vpcMatchesCidrAssociationFilter matches a cidr-block-association.* filter
// against any of the VPC's CIDR block associations.
func vpcMatchesCidrAssociationFilter(vpc *ec2store.VPC, f ec2Filter) bool {
	for _, assoc := range vpc.CidrBlockAssociationSet {
		switch f.Name {
		case "cidr-block-association.cidr-block":
			if anyMatch(f.Values, assoc.CidrBlock) {
				return true
			}
		case "cidr-block-association.association-id":
			if anyMatch(f.Values, assoc.AssociationId) {
				return true
			}
		case "cidr-block-association.state":
			if anyMatch(f.Values, assoc.CidrBlockState.State) {
				return true
			}
		}
	}
	return false
}

// parseTagsToCore converts store tags to core ec2Tag slice.
func parseTagsToCore(tags []types.Tag) []ec2Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]ec2Tag, len(tags))
	for i, t := range tags {
		result[i] = ec2Tag{Key: t.Key, Value: t.Value}
	}
	return result
}
