package ec2

import (
	"context"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

// CreateSecurityGroup creates a security group in the specified VPC.
// GroupName must be unique within the VPC scope.
func (s *EC2Service) CreateSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
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
		"GroupId": groupID,
	}, nil
}

// DescribeSecurityGroups describes one or more security groups.
// Supports GroupId, GroupName, and Filter.N for filtering.
func (s *EC2Service) DescribeSecurityGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
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
	groupID := request.GetStringParam(params, "GroupId")
	if groupID == "" {
		return nil, awserrors.NewAWSError("MissingParameter", "The request must contain the parameter GroupId", http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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
	return s.modifySecurityGroupRules(reqCtx, req, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissions = mergeIPRules(sg.IpPermissions, rules...)
	})
}

// AuthorizeSecurityGroupEgress adds one or more egress rules to a security group.
func (s *EC2Service) AuthorizeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissionsEgress = mergeIPRules(sg.IpPermissionsEgress, rules...)
	})
}

// RevokeSecurityGroupIngress removes one or more ingress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupIngress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissions = removeIPRules(sg.IpPermissions, rules...)
	})
}

// RevokeSecurityGroupEgress removes one or more egress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.modifySecurityGroupRules(reqCtx, req, func(sg *ec2store.SecurityGroup, rules []ec2store.IPRule) {
		sg.IpPermissionsEgress = removeIPRules(sg.IpPermissionsEgress, rules...)
	})
}

// modifySecurityGroupRules is the common handler for Authorize/Revoke Ingress/Egress.
func (s *EC2Service) modifySecurityGroupRules(reqCtx *request.RequestContext, req *request.ParsedRequest, apply func(*ec2store.SecurityGroup, []ec2store.IPRule)) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := resolveSecurityGroup(store, req.Parameters)
	if err != nil {
		return nil, err
	}

	rules := parseIPRules(req.Parameters, "IpPermissions")
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
