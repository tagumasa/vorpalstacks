package ec2

import (
	"context"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/types"
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

	return map[string]interface{}{"Vpc": result.Vpc}, nil
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
	nextToken := request.GetStringParam(params, "NextToken")
	maxResults := 0
	if mr := request.GetStringParam(params, "MaxResults"); mr != "" {
		if v, e := parseInt64Param(mr); e == nil {
			maxResults = int(v)
		}
	}

	result, err := s.describeVpcsCore(store, vpcIDs, filters, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Vpcs))
	for _, v := range result.Vpcs {
		items = append(items, v)
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
