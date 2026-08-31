package ec2

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	ec2store "vorpalstacks/internal/store/aws/ec2"
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
		Description: request.GetParamCaseInsensitive(params, "GroupDescription"),
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
	if err := validateFilterNames(filters, allowedSGFilters); err != nil {
		return nil, err
	}
	// The documented MaxResults exclusivity rule covers lists of IDs
	// ("both a list of IDs and MaxResults"); GroupName is a name list, so it
	// is not included here.
	nextToken, maxResults, err := parsePaginationParams(params, "GroupId")
	if err != nil {
		return nil, err
	}

	result, err := s.describeSecurityGroupsCore(store, groupIDs, groupNames, filters, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.SecurityGroups))
	for _, sg := range result.SecurityGroups {
		items = append(items, securityGroupToXMLMap(sg))
	}
	resp := map[string]interface{}{
		"SecurityGroupInfo": protocol.XMLElements{ElementName: "item", Items: items},
	}
	if result.IsTruncated && result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
}

// DeleteSecurityGroup deletes the specified security group. AWS accepts
// GroupId or (for the default VPC) GroupName; this platform has no default
// VPC, so GroupName is resolved region-wide, consistent with
// resolveSecurityGroupCore.
func (s *EC2Service) DeleteSecurityGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	if err := checkDryRun(params); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sg, err := s.resolveSecurityGroupCore(store,
		request.GetStringParam(params, "GroupId"),
		request.GetStringParam(params, "GroupName"))
	if err != nil {
		return nil, err
	}
	if err := s.deleteSecurityGroupCore(ctx, store, reqCtx.GetRegion(), sg.GroupId); err != nil {
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
// returns the securityGroupRuleSet list. Legacy flat parameters
// (CidrIp/IpProtocol/FromPort/ToPort/SourceSecurityGroup*) are still part of
// the EC2 API surface and are converted into a single permission.
// For Revoke (isAuthorize=false): removes rules by permission or by
// SecurityGroupRuleIds, returns InvalidPermission.NotFound when no rules
// matched.
func (s *EC2Service) authorizeOrRevoke(reqCtx *request.RequestContext, req *request.ParsedRequest, isAuthorize bool, isEgress bool) (interface{}, error) {
	if err := checkDryRun(req.Parameters); err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sg, err := s.resolveSecurityGroupCore(store,
		request.GetStringParam(req.Parameters, "GroupId"),
		request.GetStringParam(req.Parameters, "GroupName"))
	if err != nil {
		return nil, err
	}

	rules, err := parseIPRules(req.Parameters, "IpPermissions")
	if err != nil {
		return nil, err
	}
	ruleIDs := parseSecurityGroupRuleIDs(req.Parameters)

	if isAuthorize {
		// Authorize never accepts SecurityGroupRuleIds; rule properties come
		// from IpPermissions or the legacy flat parameters.
		if len(ruleIDs) > 0 {
			return nil, awserrors.NewAWSError("InvalidParameterCombination",
				"SecurityGroupRuleIds may not be specified for Authorize operations",
				http.StatusBadRequest)
		}
		if len(rules) == 0 {
			rules, err = parseLegacyFlatRule(req.Parameters)
			if err != nil {
				return nil, err
			}
		}
		if len(rules) == 0 {
			return nil, awserrors.NewMissingParameter("IpPermissions or the legacy rule parameters (IpProtocol, CidrIp, ...) are required")
		}
		result, err := s.authorizeSecurityGroupRulesCore(store, sg, rules, isEgress, reqCtx.GetRegion())
		if err != nil {
			return nil, err
		}
		ruleItems := make([]interface{}, 0, len(result.SecurityGroupRules))
		for _, r := range result.SecurityGroupRules {
			ruleItems = append(ruleItems, sgRuleToXMLMap(r))
		}
		return map[string]interface{}{
			"return":               true,
			"SecurityGroupRuleSet": protocol.XMLElements{ElementName: "item", Items: ruleItems},
		}, nil
	}

	// Revoke accepts either rule IDs (SecurityGroupRuleId.N) or rule
	// properties (IpPermissions or the legacy flat parameters), never both.
	if len(ruleIDs) > 0 && (len(rules) > 0 || hasFlatRuleParams(req.Parameters)) {
		return nil, awserrors.NewAWSError("InvalidParameterCombination",
			"SecurityGroupRuleIds and rule properties (IpPermissions or the legacy rule parameters) are mutually exclusive",
			http.StatusBadRequest)
	}
	if len(rules) == 0 && len(ruleIDs) == 0 {
		rules, err = parseLegacyFlatRule(req.Parameters)
		if err != nil {
			return nil, err
		}
	}
	if len(rules) == 0 && len(ruleIDs) == 0 {
		return nil, awserrors.NewMissingParameter("IpPermissions, the legacy rule parameters (IpProtocol, CidrIp, ...) or SecurityGroupRuleIds is required")
	}

	result, err := s.revokeSecurityGroupRulesCore(store, sg, rules, ruleIDs, isEgress)
	if err != nil {
		return nil, err
	}
	revokedItems := make([]interface{}, 0, len(result.RevokedSecurityGroupRules))
	for _, r := range result.RevokedSecurityGroupRules {
		revokedItems = append(revokedItems, revokedRuleToXMLMap(r))
	}
	resp := map[string]interface{}{
		"return":                      true,
		"RevokedSecurityGroupRuleSet": protocol.XMLElements{ElementName: "item", Items: revokedItems},
	}
	if len(result.UnknownIpPermissions) > 0 {
		unknownItems := make([]interface{}, 0, len(result.UnknownIpPermissions))
		for _, u := range result.UnknownIpPermissions {
			unknownItems = append(unknownItems, ipRuleToXMLMap(u))
		}
		resp["UnknownIpPermissionSet"] = protocol.XMLElements{ElementName: "item", Items: unknownItems}
	}
	return resp, nil
}

// matchesSGFilters checks if a security group matches all the given filters.
// Filter names are validated by validateFilterNames before matching.
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
		case "description":
			if !anyMatch(f.Values, sg.Description) {
				return false
			}
		case "vpc-id":
			if !anyMatch(f.Values, sg.VpcId) {
				return false
			}
		case "owner-id":
			if !anyMatch(f.Values, sg.OwnerId) {
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
		default:
			// AWS tag key/value filters use the "tag:<key>" name form; the
			// filter matches when the resource carries that tag key with
			// any of the filter values.
			if strings.HasPrefix(f.Name, "tag:") {
				if !hasTagKeyValues(sg.Tags, strings.TrimPrefix(f.Name, "tag:"), f.Values) {
					return false
				}
				continue
			}
			if !sgMatchesPermissionFilter(sg, f) {
				return false
			}
		}
	}
	return true
}

// sgMatchesPermissionFilter handles the ip-permission.* (inbound),
// ingress.ip-permission.* (alias for inbound) and egress.ip-permission.*
// (outbound) filter families. A permission filter matches when any rule in
// the corresponding direction satisfies the value.
func sgMatchesPermissionFilter(sg *ec2store.SecurityGroup, f ec2Filter) bool {
	var rules []ec2store.IPRule
	switch {
	case strings.HasPrefix(f.Name, "ip-permission."), strings.HasPrefix(f.Name, "ingress.ip-permission."):
		rules = sg.IpPermissions
	case strings.HasPrefix(f.Name, "egress.ip-permission."):
		rules = sg.IpPermissionsEgress
	default:
		return false
	}
	sub := strings.TrimPrefix(strings.TrimPrefix(f.Name, "ingress."), "egress.")
	sub = strings.TrimPrefix(sub, "ip-permission.")
	for _, r := range rules {
		switch sub {
		case "protocol":
			if anyMatch(f.Values, r.IpProtocol) {
				return true
			}
		case "from-port":
			if anyMatchInt64Port(f.Values, r.FromPort) {
				return true
			}
		case "to-port":
			if anyMatchInt64Port(f.Values, r.ToPort) {
				return true
			}
		case "cidr":
			for _, ip := range r.IpRanges {
				if anyMatch(f.Values, ip.CidrIp) {
					return true
				}
			}
		case "ipv6-cidr":
			for _, ip := range r.Ipv6Ranges {
				if anyMatch(f.Values, ip.CidrIp) {
					return true
				}
			}
		case "group-id":
			for _, g := range r.UserIdGroupPairs {
				if anyMatch(f.Values, g.GroupId) {
					return true
				}
			}
		case "group-name":
			for _, g := range r.UserIdGroupPairs {
				if anyMatch(f.Values, g.GroupName) {
					return true
				}
			}
		case "user-id":
			for _, g := range r.UserIdGroupPairs {
				if anyMatch(f.Values, g.UserId) {
					return true
				}
			}
		case "prefix-list-id":
			for _, pl := range r.PrefixListIds {
				if anyMatch(f.Values, pl.PrefixListId) {
					return true
				}
			}
		}
	}
	return false
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

// revokedRuleToXMLMap converts a revoked SecurityGroupRule to the wire shape
// of RevokedSecurityGroupRule (securityGroupRuleId + referencedGroupId,
// without referencedGroupInfo or the rule ARN).
func revokedRuleToXMLMap(r SecurityGroupRule) map[string]interface{} {
	m := map[string]interface{}{
		"SecurityGroupRuleId": r.RuleId,
		"IsEgress":            r.IsEgress,
		"GroupId":             r.GroupId,
		"IpProtocol":          r.IpProtocol,
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
	if r.PrefixListId != "" {
		m["PrefixListId"] = r.PrefixListId
	}
	if r.ReferencedGroupId != "" {
		m["ReferencedGroupId"] = r.ReferencedGroupId
	}
	if r.Description != "" {
		m["Description"] = r.Description
	}
	return m
}

// parseSecurityGroupRuleIDs extracts SecurityGroupRuleId.N from request
// parameters (the EC2 Query wire name of the Revoke SecurityGroupRuleIds
// list member is SecurityGroupRuleId).
func parseSecurityGroupRuleIDs(params map[string]interface{}) []string {
	var ids []string
	for i := 1; ; i++ {
		v := request.GetStringParam(params, "SecurityGroupRuleId."+strconv.Itoa(i))
		if v == "" {
			break
		}
		ids = append(ids, v)
	}
	return ids
}

// flatRuleParamNames are the legacy single-rule query parameters accepted by
// the Authorize/Revoke operations alongside IpPermissions.
var flatRuleParamNames = []string{"IpProtocol", "CidrIp", "FromPort", "ToPort", "SourceSecurityGroupName", "SourceSecurityGroupOwnerId"}

// hasFlatRuleParams reports whether any legacy flat rule parameter is set.
func hasFlatRuleParams(params map[string]interface{}) bool {
	for _, name := range flatRuleParamNames {
		if request.GetStringParam(params, name) != "" {
			return true
		}
	}
	return false
}

// parseLegacyFlatRule converts the deprecated flat Authorize/Revoke request
// parameters (CidrIp/IpProtocol/FromPort/ToPort/SourceSecurityGroupName/
// SourceSecurityGroupOwnerId) into a single IP permission rule.
func parseLegacyFlatRule(params map[string]interface{}) ([]ec2store.IPRule, error) {
	ipProtocol := request.GetStringParam(params, "IpProtocol")
	cidr := request.GetStringParam(params, "CidrIp")
	fromPortRaw := request.GetStringParam(params, "FromPort")
	toPortRaw := request.GetStringParam(params, "ToPort")
	sourceGroup := request.GetStringParam(params, "SourceSecurityGroupName")
	sourceOwner := request.GetStringParam(params, "SourceSecurityGroupOwnerId")

	if !hasFlatRuleParams(params) {
		return nil, nil
	}
	if ipProtocol == "" {
		return nil, awserrors.NewMissingParameter("IpPermissions or IpProtocol is required")
	}

	fromPort := int64(-1)
	if fromPortRaw != "" {
		v, err := parseInt64Param(fromPortRaw, "FromPort")
		if err != nil {
			return nil, err
		}
		fromPort = v
	}
	toPort := int64(-1)
	if toPortRaw != "" {
		v, err := parseInt64Param(toPortRaw, "ToPort")
		if err != nil {
			return nil, err
		}
		toPort = v
	}

	rule := ec2store.IPRule{
		IpProtocol: ipProtocol,
		FromPort:   fromPort,
		ToPort:     toPort,
	}
	if cidr != "" {
		if err := validateIPv4CIDRInRule(cidr); err != nil {
			return nil, err
		}
		rule.IpRanges = []ec2store.IPRange{{CidrIp: cidr}}
	}
	if sourceGroup != "" {
		rule.UserIdGroupPairs = []ec2store.GroupPair{{
			GroupName: sourceGroup,
			UserId:    sourceOwner,
		}}
	}
	if len(rule.IpRanges) == 0 && len(rule.UserIdGroupPairs) == 0 {
		return nil, awserrors.NewAWSError("InvalidPermission.Malformed",
			"A permission must specify at least one of CidrIp or SourceSecurityGroupName",
			http.StatusBadRequest)
	}
	return []ec2store.IPRule{rule}, nil
}
