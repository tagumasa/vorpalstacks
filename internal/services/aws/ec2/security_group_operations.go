package ec2

import (
	"context"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateSecurityGroup creates a security group in the specified VPC.
func (s *EC2Service) CreateSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createSecurityGroupCore(store, CreateSecurityGroupInput{
		GroupName:   request.GetStringParam(params, "GroupName"),
		Description: request.GetStringParam(params, "Description"),
		VpcId:       request.GetStringParam(params, "VpcId"),
		Tags:        parseTagsToCore(parseEC2Tags(params)),
	}, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"GroupId":          result.GroupId,
		"SecurityGroupArn": result.Arn,
		"TagSet":           protocol.XMLElements{ElementName: "item", Items: sgTagsToInterface(result.SecurityGroup.Tags)},
	}, nil
}

// DescribeSecurityGroups describes one or more security groups.
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

	result, err := s.describeSecurityGroupsCore(store, groupIDs, groupNames, filters, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.SecurityGroups))
	for _, sg := range result.SecurityGroups {
		items = append(items, sg)
	}
	resp := map[string]interface{}{
		"SecurityGroupInfo": protocol.XMLElements{ElementName: "item", Items: items},
	}
	if result.IsTruncated && result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
}

// DeleteSecurityGroup deletes the specified security group.
func (s *EC2Service) DeleteSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSecurityGroupCore(ctx, store, reqCtx.GetRegion(), request.GetStringParam(params, "GroupId")); err != nil {
		return nil, err
	}
	return map[string]interface{}{"return": true}, nil
}

// AuthorizeSecurityGroupIngress adds one or more ingress rules to a security group.
func (s *EC2Service) AuthorizeSecurityGroupIngress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.authorizeOrRevoke(reqCtx, req, true, false)
}

// AuthorizeSecurityGroupEgress adds one or more egress rules to a security group.
func (s *EC2Service) AuthorizeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.authorizeOrRevoke(reqCtx, req, true, true)
}

// RevokeSecurityGroupIngress removes one or more ingress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupIngress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.authorizeOrRevoke(reqCtx, req, false, false)
}

// RevokeSecurityGroupEgress removes one or more egress rules from a security group.
func (s *EC2Service) RevokeSecurityGroupEgress(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.authorizeOrRevoke(reqCtx, req, false, true)
}

// authorizeOrRevoke handles both Authorize and Revoke for ingress and egress.
// For Authorize (isAuthorize=true): validates rules, generates rule IDs,
// returns SecurityGroupRules list.
// For Revoke (isAuthorize=false): removes rules, returns InvalidPermission.NotFound
// when no rules matched.
func (s *EC2Service) authorizeOrRevoke(reqCtx *request.RequestContext, req *request.ParsedRequest, isAuthorize bool, isEgress bool) (interface{}, error) {
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

	rules, err := parseIPRules(req.Parameters, "IpPermissions")
	if err != nil {
		return nil, err
	}

	if isAuthorize {
		result, err := s.authorizeSecurityGroupRulesCore(store, sg, rules, isEgress, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
		ruleItems := make([]interface{}, 0, len(result.SecurityGroupRules))
		for _, r := range result.SecurityGroupRules {
			ruleItems = append(ruleItems, sgRuleToXMLMap(r))
		}
		return map[string]interface{}{
			"SecurityGroupRules": protocol.XMLElements{ElementName: "item", Items: ruleItems},
		}, nil
	}

	result, err := s.revokeSecurityGroupRulesCore(store, sg, rules, isEgress)
	if err != nil {
		return nil, err
	}
	revokedItems := make([]interface{}, 0, len(result.RevokedSecurityGroupRules))
	for _, r := range result.RevokedSecurityGroupRules {
		revokedItems = append(revokedItems, sgRuleToXMLMap(r))
	}
	return map[string]interface{}{
		"RevokedSecurityGroupRules": protocol.XMLElements{ElementName: "item", Items: revokedItems},
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

// sgReferencesGroup returns true if any rule references the specified groupId.
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

// sgRuleToXMLMap converts a SecurityGroupRule to a map for XML output.
func sgRuleToXMLMap(r SecurityGroupRule) map[string]interface{} {
	m := map[string]interface{}{
		"SecurityGroupRuleId":  r.RuleId,
		"IsEgress":             r.IsEgress,
		"SecurityGroupRuleArn": r.SecurityGroupRuleArn,
		"GroupId":              r.GroupId,
		"GroupOwnerId":         r.GroupOwnerId,
		"IpProtocol":           r.IpProtocol,
	}
	if r.IpProtocol != "-1" {
		m["FromPort"] = r.FromPort
		m["ToPort"] = r.ToPort
	}
	if r.CidrIpv4 != "" {
		m["CidrIpv4"] = r.CidrIpv4
	}
	if r.CidrIpv6 != "" {
		m["CidrIpv6"] = r.CidrIpv6
	}
	if r.Description != "" {
		m["Description"] = r.Description
	}
	if r.ReferencedGroupId != "" {
		m["ReferencedGroupInfo"] = map[string]interface{}{
			"GroupId": r.ReferencedGroupId,
			"UserId":  r.ReferencedGroupUserId,
		}
		if r.ReferencedGroupName != "" {
			m["ReferencedGroupInfo"].(map[string]interface{})["GroupName"] = r.ReferencedGroupName
		}
	}
	if r.PrefixListId != "" {
		m["PrefixListId"] = r.PrefixListId
	}
	return m
}
