package ec2

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateSecurityGroup creates a security group in the specified VPC.
// GroupName must be unique within the VPC scope.
func (s *EC2Service) CreateSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	groupName := request.GetStringParam(params, "GroupName")
	if groupName == "" {
		return nil, awserrors.NewMissingParameter("GroupName is required")
	}
	description := request.GetStringParam(params, "Description")
	if description == "" {
		description = groupName
	}
	vpcID := request.GetStringParam(params, "VpcId")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if vpcID != "" {
		if _, err := store.GetVPC(vpcID); err != nil {
			return nil, translateStoreError(err)
		}
	}

	existingSGs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for _, sg := range existingSGs {
		if sg.GroupName == groupName && sg.VpcId == vpcID {
			return nil, awserrors.NewAWSError("InvalidGroup.Duplicate",
				"Security group '"+groupName+"' already exists",
				http.StatusBadRequest)
		}
	}

	groupID, err := GenerateSecurityGroupID()
	if err != nil {
		return nil, err
	}

	sg := &ec2store.SecurityGroup{
		GroupId:     groupID,
		GroupName:   groupName,
		Description: description,
		VpcId:       vpcID,
		OwnerId:     s.accountID,
		Tags:        parseEC2Tags(params),
		IpPermissionsEgress: []ec2store.IPRule{
			{
				IpProtocol: "-1",
				FromPort:   -1,
				ToPort:     -1,
				IpRanges:   []ec2store.IPRange{{CidrIp: "0.0.0.0/0"}},
			},
		},
	}

	if err := store.CreateSecurityGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"GroupId":          groupID,
		"SecurityGroupArn": fmt.Sprintf("arn:aws:ec2:%s:%s:security-group/%s", reqCtx.GetRegion(), s.accountID, groupID),
		"TagSet":           protocol.XMLElements{ElementName: "item", Items: sgTagsToInterface(sg.Tags)},
	}, nil
}

// DescribeSecurityGroups describes one or more security groups.
// Supports GroupId, GroupName, and Filter.N for filtering.
func (s *EC2Service) DescribeSecurityGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	groupIDs := request.GetStringList(params, "GroupId")
	groupNames := request.GetStringList(params, "GroupName")

	if len(groupIDs) > 0 {
		items := make([]interface{}, 0, len(groupIDs))
		for _, id := range groupIDs {
			sg, err := store.GetSecurityGroup(id)
			if err != nil {
				return nil, translateStoreError(err)
			}
			items = append(items, sg)
		}
		return map[string]interface{}{
			"SecurityGroupInfo": protocol.XMLElements{ElementName: "item", Items: items},
		}, nil
	}

	sgs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, err
	}

	filters := parseFilters(params)
	items := make([]interface{}, 0, len(sgs))
	for _, sg := range sgs {
		if len(groupNames) > 0 && !anyMatch(groupNames, sg.GroupName) {
			continue
		}
		if matchesSGFilters(sg, filters) {
			items = append(items, sg)
		}
	}
	return map[string]interface{}{
		"SecurityGroupInfo": protocol.XMLElements{ElementName: "item", Items: items},
	}, nil
}

// DeleteSecurityGroup deletes the specified security group.
func (s *EC2Service) DeleteSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	groupID := request.GetStringParam(params, "GroupId")
	if groupID == "" {
		return nil, awserrors.NewAWSError("MissingParameter", "The request must contain the parameter GroupId", http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Existence check first (AWS spec ordering): return InvalidGroup.NotFound
	// before any dependency evaluation.
	if _, err := store.GetSecurityGroup(groupID); err != nil {
		return nil, translateStoreError(err)
	}

	// DependencyViolation: reject if another SG references this one via
	// UserIdGroupPairs in its ingress or egress rules.
	allSGs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}
	for _, sg := range allSGs {
		if sg.GroupId == groupID {
			continue
		}
		if sgReferencesGroup(sg.IpPermissions, groupID) || sgReferencesGroup(sg.IpPermissionsEgress, groupID) {
			return nil, awserrors.NewAWSError("DependencyViolation",
				"The security group '"+groupID+"' is being referenced by security group '"+sg.GroupId+"'",
				http.StatusBadRequest)
		}
	}

	// Cross-service dependency check: verify no Lambda function VpcConfig
	// references this security group.
	if s.bus != nil {
		for _, checker := range s.bus.SecurityGroupUsageCheckers() {
			if checker.IsSecurityGroupInUse(ctx, reqCtx.GetRegion(), groupID) {
				return nil, awserrors.NewAWSError(
					"DependencyViolation",
					"The security group '"+groupID+"' is being used by another resource",
					http.StatusBadRequest,
				)
			}
		}
	}

	if err := store.DeleteSecurityGroup(groupID); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// AuthorizeSecurityGroupIngress adds one or more ingress rules to a security group.
func (s *EC2Service) AuthorizeSecurityGroupIngress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, true, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissions = mergeIPRules(sg.IpPermissions, rules...)
	})
}

// AuthorizeSecurityGroupEgress adds one or more egress rules to a security group.
func (s *EC2Service) AuthorizeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, true, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissionsEgress = mergeIPRules(sg.IpPermissionsEgress, rules...)
	})
}

// RevokeSecurityGroupIngress removes one or more ingress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupIngress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, false, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissions = removeIPRules(sg.IpPermissions, rules...)
	})
}

// RevokeSecurityGroupEgress removes one or more egress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, false, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissionsEgress = removeIPRules(sg.IpPermissionsEgress, rules...)
	})
}

// modifySecurityGroupRules is the common handler for Authorize/Revoke Ingress/Egress.
// validateRules controls whether new rules are validated: Authorize passes true
// so malformed rules are rejected; Revoke passes false so legacy rules that were
// created with looser validation can still be removed by their owner.
func (s *EC2Service) modifySecurityGroupRules(reqCtx *request.RequestContext, req *request.ParsedRequest, validateRules bool, apply func(*ec2store.SecurityGroup, []ec2store.IPRule)) (interface{}, error) {
	if err := checkDryRun(req.Parameters); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := resolveSecurityGroup(store, req.Parameters)
	if err != nil {
		return nil, err
	}

	rules := parseIPRules(req.Parameters, "IpPermissions")
	if validateRules {
		for _, r := range rules {
			if err := validateIPRule(r); err != nil {
				return nil, err
			}
		}
	}
	apply(sg, rules)

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// resolveSecurityGroup finds the security group by GroupId or GroupName.
func resolveSecurityGroup(store *ec2store.EC2Store, params map[string]interface{}) (*ec2store.SecurityGroup, error) {
	groupID := request.GetStringParam(params, "GroupId")
	if groupID != "" {
		sg, err := store.GetSecurityGroup(groupID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return sg, nil
	}

	groupName := request.GetStringParam(params, "GroupName")
	if groupName == "" {
		return nil, awserrors.NewMissingParameter("GroupId or GroupName is required")
	}

	sgs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for _, sg := range sgs {
		if sg.GroupName == groupName {
			return sg, nil
		}
	}
	return nil, awserrors.NewAWSError("InvalidGroup.NotFound", "The security group does not exist", http.StatusNotFound)
}

// matchesSGFilters checks if a security group matches all the given filters.
func matchesSGFilters(sg *ec2store.SecurityGroup, filters []ec2Filter) bool {
	for _, f := range filters {
		switch f.Name {
		case "group-id":
			if !anyMatch(f.Values, sg.GroupId) {
				return false
			}
		case "group-name":
			if !anyMatch(f.Values, sg.GroupName) {
				return false
			}
		case "vpc-id":
			if !anyMatch(f.Values, sg.VpcId) {
				return false
			}
		case "tag-key":
			if !hasTagKey(sg.Tags, f.Values) {
				return false
			}
		case "tag-value":
			if !hasTagValue(sg.Tags, f.Values) {
				return false
			}
		case "tag":
			if !hasTagKeyValue(sg.Tags, f.Values) {
				return false
			}
		}
	}
	return true
}

// sgReferencesGroup returns true if any rule in the given list contains a
// UserIdGroupPair referencing the specified groupId.
func sgReferencesGroup(rules []ec2store.IPRule, groupId string) bool {
	for _, rule := range rules {
		for _, pair := range rule.UserIdGroupPairs {
			if pair.GroupId == groupId {
				return true
			}
		}
	}
	return false
}

// sgTagsToInterface converts store tags to interface slice for XML output.
func sgTagsToInterface(tags []types.Tag) []interface{} {
	items := make([]interface{}, 0, len(tags))
	for _, t := range tags {
		items = append(items, t)
	}
	return items
}
