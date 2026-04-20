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

// CreateSecurityGroup creates a security group in the specified VPC.
func (s *EC2Service) CreateSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	groupName := getStringParam(params, "GroupName")
	if groupName == "" {
		return nil, fmt.Errorf("ec2: GroupName is required")
	}
	description := getStringParam(params, "Description")
	if description == "" {
		description = groupName
	}
	vpcID := getStringParam(params, "VpcId")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if vpcID != "" {
		if _, err := store.GetVPC(vpcID); err != nil {
			return nil, translateStoreError(err)
		}
	}

	groupID, err := GenerateSecurityGroupID()
	if err != nil {
		return nil, err
	}

	tags := ParseTags(params)

	sg := &ec2store.SecurityGroup{
		GroupId:     groupID,
		GroupName:   groupName,
		Description: description,
		VpcId:       vpcID,
		OwnerId:     s.accountID,
		Tags:        tags,
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

// DescribeSecurityGroups describes one or all security groups.
func (s *EC2Service) DescribeSecurityGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	groupID := getStringParam(params, "GroupId")
	groupName := getStringParam(params, "GroupName")

	if groupID != "" {
		sg, err := store.GetSecurityGroup(groupID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"SecurityGroupInfo": protocol.XMLElements{ElementName: "item", Items: []interface{}{sg}},
		}, nil
	}

	sgs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(sgs))
	for _, sg := range sgs {
		if groupName != "" && sg.GroupName != groupName {
			continue
		}
		items = append(items, sg)
	}
	return map[string]interface{}{
		"SecurityGroupInfo": protocol.XMLElements{ElementName: "item", Items: items},
	}, nil
}

// DeleteSecurityGroup deletes the specified security group.
func (s *EC2Service) DeleteSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	groupID := getStringParam(params, "GroupId")
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := resolveSecurityGroup(store, req.Parameters)
	if err != nil {
		return nil, err
	}

	newRules := parseIPRules(req.Parameters, "IpPermissions")
	sg.IpPermissions = mergeIPRules(sg.IpPermissions, newRules...)

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// AuthorizeSecurityGroupEgress adds one or more egress rules to a security group.
func (s *EC2Service) AuthorizeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := resolveSecurityGroup(store, req.Parameters)
	if err != nil {
		return nil, err
	}

	newRules := parseIPRules(req.Parameters, "IpPermissions")
	sg.IpPermissionsEgress = mergeIPRules(sg.IpPermissionsEgress, newRules...)

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// RevokeSecurityGroupIngress removes one or more ingress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupIngress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := resolveSecurityGroup(store, req.Parameters)
	if err != nil {
		return nil, err
	}

	removeRules := parseIPRules(req.Parameters, "IpPermissions")
	sg.IpPermissions = removeIPRules(sg.IpPermissions, removeRules...)

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// RevokeSecurityGroupEgress removes one or more egress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := resolveSecurityGroup(store, req.Parameters)
	if err != nil {
		return nil, err
	}

	removeRules := parseIPRules(req.Parameters, "IpPermissions")
	sg.IpPermissionsEgress = removeIPRules(sg.IpPermissionsEgress, removeRules...)

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"return": true,
	}, nil
}

// resolveSecurityGroup finds the security group by GroupId or GroupName.
func resolveSecurityGroup(store *ec2store.EC2Store, params map[string]interface{}) (*ec2store.SecurityGroup, error) {
	groupID := getStringParam(params, "GroupId")
	if groupID != "" {
		sg, err := store.GetSecurityGroup(groupID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return sg, nil
	}

	groupName := getStringParam(params, "GroupName")
	if groupName == "" {
		return nil, fmt.Errorf("ec2: GroupId or GroupName is required")
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

// parseIPRules extracts IP permission rules from request parameters.
func parseIPRules(params map[string]interface{}, prefix string) []ec2store.IPRule {
	var rules []ec2store.IPRule

	for i := 1; ; i++ {
		ipProtocol := getStringParam(params, fmt.Sprintf("%s.%d.IpProtocol", prefix, i))
		if ipProtocol == "" {
			break
		}

		rule := ec2store.IPRule{
			IpProtocol: ipProtocol,
		}

		fromPort := getStringParam(params, fmt.Sprintf("%s.%d.FromPort", prefix, i))
		if fromPort != "" {
			rule.FromPort = parseInt64(fromPort)
		} else {
			rule.FromPort = -1
		}

		toPort := getStringParam(params, fmt.Sprintf("%s.%d.ToPort", prefix, i))
		if toPort != "" {
			rule.ToPort = parseInt64(toPort)
		} else {
			rule.ToPort = -1
		}

		groupID := getStringParam(params, fmt.Sprintf("%s.%d.GroupId", prefix, i))
		if groupID != "" {
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, ec2store.GroupPair{GroupId: groupID})
		}

		cidrIP := getStringParam(params, fmt.Sprintf("%s.%d.CidrIp", prefix, i))
		if cidrIP != "" {
			rule.IpRanges = append(rule.IpRanges, ec2store.IPRange{CidrIp: cidrIP})
		}

		rules = append(rules, rule)
	}

	return rules
}

// parseInt64 parses a decimal string to int64, ignoring non-digit characters.
func parseInt64(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		}
	}
	return n
}
