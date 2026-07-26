package ec2

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
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
	vpcID := request.GetStringParam(params, "VpcId")
	if vpcID == "" {
		return nil, awserrors.NewMissingParameter("VpcId is required")
	}
	cidrBlock := request.GetStringParam(params, "CidrBlock")
	if cidrBlock == "" {
		return nil, awserrors.NewMissingParameter("CidrBlock is required")
	}
	if err := validateCIDRBlock(cidrBlock); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	vpc, err := store.GetVPC(vpcID)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := validateSubnetInVPC(cidrBlock, vpc.CidrBlock); err != nil {
		return nil, err
	}

	// CIDR overlap check: reject if another subnet in the same VPC has an
	// overlapping CIDR block.
	existingSubnets, err := store.ListSubnetsByVPC(vpcID)
	if err != nil {
		return nil, translateStoreError(err)
	}
	if err := validateSubnetCIDROverlap(cidrBlock, existingSubnets); err != nil {
		return nil, err
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
		SubnetId:                subnetID,
		VpcId:                   vpcID,
		CidrBlock:               cidrBlock,
		AvailabilityZone:        az,
		AvailableIpAddressCount: calculateAvailableIPs(cidrBlock),
		State:                   "available",
		OwnerId:                 s.accountID,
		SubnetArn:               fmt.Sprintf("arn:aws:ec2:%s:%s:subnet/%s", reqCtx.GetRegion(), s.accountID, subnetID),
		MapPublicIpOnLaunch:     mapPublicIpOnLaunch,
		SubnetType:              "ipv4",
		Tags:                    parseEC2Tags(params),
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
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subnetIDs := request.GetStringList(params, "SubnetId")
	if len(subnetIDs) > 0 {
		items := make([]interface{}, 0, len(subnetIDs))
		for _, id := range subnetIDs {
			subnet, err := store.GetSubnet(id)
			if err != nil {
				return nil, translateStoreError(err)
			}
			items = append(items, subnet)
		}
		return map[string]interface{}{
			"SubnetSet": protocol.XMLElements{ElementName: "item", Items: items},
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
// AWS checks for instances, ENIs, NAT gateways, and VPC endpoints before
// allowing deletion. On this edge platform, those resources do not exist.
// Cross-service consumers (Lambda VpcConfig, Neptune DB subnet groups) are
// checked via registered SubnetUsageCheckers on the event bus.
func (s *EC2Service) DeleteSubnet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	subnetID := request.GetStringParam(params, "SubnetId")
	if subnetID == "" {
		return nil, awserrors.NewMissingParameter("SubnetId is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Existence check first (AWS spec ordering): return InvalidSubnetID.NotFound
	// before any dependency evaluation.
	if _, err := store.GetSubnet(subnetID); err != nil {
		return nil, translateStoreError(err)
	}

	// Cross-service dependency check: verify no Lambda function VpcConfig or
	// Neptune DB subnet group references this subnet.
	if s.bus != nil {
		for _, checker := range s.bus.SubnetUsageCheckers() {
			if checker.IsSubnetInUse(ctx, reqCtx.GetRegion(), subnetID) {
				return nil, awserrors.NewAWSError(
					"DependencyViolation",
					fmt.Sprintf("The subnet '%s' has dependencies and cannot be deleted", subnetID),
					http.StatusBadRequest,
				)
			}
		}
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
