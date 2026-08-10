package ec2

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/generators"
)

const securityGroupRulePrefix = "sgr-"

// CreateSecurityGroupInput is the transport-agnostic input for CreateSecurityGroup.
type CreateSecurityGroupInput struct {
	GroupName   string
	Description string
	VpcId       string
	Tags        []ec2Tag
}

// SecurityGroupResult is the transport-agnostic result for SG operations.
type SecurityGroupResult struct {
	SecurityGroup *ec2store.SecurityGroup
	GroupId       string
	Arn           string
}

// SecurityGroupListResult is the transport-agnostic result for DescribeSecurityGroups.
type SecurityGroupListResult struct {
	SecurityGroups []*ec2store.SecurityGroup
	NextToken      string
	IsTruncated    bool
}

// AuthorizeResult is the transport-agnostic result for Authorize operations.
// SecurityGroupRules is the flat list of newly-added rules with IDs.
type AuthorizeResult struct {
	SecurityGroupRules []SecurityGroupRule
}

// SecurityGroupRule is a single rule entry returned by Authorize operations.
type SecurityGroupRule struct {
	RuleId                string
	IsEgress              bool
	SecurityGroupRuleArn  string
	GroupId               string
	GroupOwnerId          string
	IpProtocol            string
	FromPort              int64
	ToPort                int64
	CidrIpv4              string
	CidrIpv6              string
	Description           string
	ReferencedGroupId     string
	ReferencedGroupName   string
	ReferencedGroupUserId string
	PrefixListId          string
}

// RevokeResult is the transport-agnostic result for Revoke operations.
type RevokeResult struct {
	RevokedSecurityGroupRules []SecurityGroupRule
	UnknownIpPermissions      []ec2store.IPRule
}

// generateSecurityGroupRuleID generates a unique security group rule ID.
func generateSecurityGroupRuleID() (string, error) {
	return generators.GenerateIDWithPrefix(securityGroupRulePrefix, 17)
}

// createSecurityGroupCore contains the business logic for CreateSecurityGroup.
func (s *EC2Service) createSecurityGroupCore(store *ec2store.EC2Store, input CreateSecurityGroupInput, region string) (*SecurityGroupResult, error) {
	if input.GroupName == "" {
		return nil, awserrors.NewMissingParameter("GroupName is required")
	}
	if input.Description == "" {
		return nil, awserrors.NewAWSError("MissingParameter",
			"Description is required",
			http.StatusBadRequest)
	}
	if input.VpcId == "" {
		return nil, awserrors.NewAWSError("InvalidParameterValue",
			"VpcId is required for non-default VPC security groups",
			http.StatusBadRequest)
	}

	if _, err := store.GetVPC(input.VpcId); err != nil {
		return nil, translateStoreError(err)
	}

	existingSGs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for _, sg := range existingSGs {
		if sg.GroupName == input.GroupName && sg.VpcId == input.VpcId {
			return nil, awserrors.NewAWSError("InvalidGroup.Duplicate",
				"Security group '"+input.GroupName+"' already exists",
				http.StatusBadRequest)
		}
	}

	groupID, err := GenerateSecurityGroupID()
	if err != nil {
		return nil, err
	}

	arn := fmt.Sprintf("arn:aws:ec2:%s:%s:security-group/%s", region, s.accountID, groupID)

	sg := &ec2store.SecurityGroup{
		GroupId:          groupID,
		GroupName:        input.GroupName,
		Description:      input.Description,
		VpcId:            input.VpcId,
		OwnerId:          s.accountID,
		SecurityGroupArn: arn,
		Tags:             toStoreTags(input.Tags),
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

	return &SecurityGroupResult{
		SecurityGroup: sg,
		GroupId:       groupID,
		Arn:           arn,
	}, nil
}

// describeSecurityGroupsCore contains the business logic for DescribeSecurityGroups.
func (s *EC2Service) describeSecurityGroupsCore(store *ec2store.EC2Store, groupIDs, groupNames []string, filters []ec2Filter, nextToken string, maxResults int) (*SecurityGroupListResult, error) {
	if len(groupIDs) > 0 {
		sgs := make([]*ec2store.SecurityGroup, 0, len(groupIDs))
		for _, id := range groupIDs {
			sg, err := store.GetSecurityGroup(id)
			if err != nil {
				return nil, translateStoreError(err)
			}
			sgs = append(sgs, sg)
		}
		return &SecurityGroupListResult{SecurityGroups: sgs}, nil
	}

	allSGs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}

	filtered := make([]*ec2store.SecurityGroup, 0, len(allSGs))
	for _, sg := range allSGs {
		if len(groupNames) > 0 && !anyMatch(groupNames, sg.GroupName) {
			continue
		}
		if matchesSGFilters(sg, filters) {
			filtered = append(filtered, sg)
		}
	}

	if maxResults <= 0 {
		maxResults = 100
	}
	page := pagination.PaginateSlice(filtered, nextToken, maxResults, func(sg *ec2store.SecurityGroup) string {
		return sg.GroupId
	})
	return &SecurityGroupListResult{
		SecurityGroups: page.Items,
		NextToken:      page.NextMarker,
		IsTruncated:    page.IsTruncated,
	}, nil
}

// deleteSecurityGroupCore contains the business logic for DeleteSecurityGroup.
func (s *EC2Service) deleteSecurityGroupCore(ctx context.Context, store *ec2store.EC2Store, region, groupID string) error {
	if groupID == "" {
		return awserrors.NewAWSError("MissingParameter", "The request must contain the parameter GroupId", http.StatusBadRequest)
	}

	if _, err := store.GetSecurityGroup(groupID); err != nil {
		return translateStoreError(err)
	}

	allSGs, err := store.ListSecurityGroups()
	if err != nil {
		return translateStoreError(err)
	}
	for _, sg := range allSGs {
		if sg.GroupId == groupID {
			continue
		}
		if sgReferencesGroup(sg.IpPermissions, groupID) || sgReferencesGroup(sg.IpPermissionsEgress, groupID) {
			return awserrors.NewAWSError("DependencyViolation",
				"The security group '"+groupID+"' is being referenced by security group '"+sg.GroupId+"'",
				http.StatusBadRequest)
		}
	}

	if s.bus != nil {
		for _, checker := range s.bus.SecurityGroupUsageCheckers() {
			if checker.IsSecurityGroupInUse(ctx, region, groupID) {
				return awserrors.NewAWSError(
					"DependencyViolation",
					"The security group '"+groupID+"' is being used by another resource",
					http.StatusBadRequest,
				)
			}
		}
	}

	return translateStoreError(store.DeleteSecurityGroup(groupID))
}

// authorizeSecurityGroupRulesCore adds rules to a security group and returns
// the newly created SecurityGroupRules with generated rule IDs.
func (s *EC2Service) authorizeSecurityGroupRulesCore(store *ec2store.EC2Store, sg *ec2store.SecurityGroup, rules []ec2store.IPRule, isEgress bool, region string) (*AuthorizeResult, error) {
	if err := validateGroupReferences(store, rules); err != nil {
		return nil, err
	}
	for _, r := range rules {
		if err := validateIPRule(r); err != nil {
			return nil, err
		}
		if err := validateIPRuleRanges(r); err != nil {
			return nil, err
		}
	}

	var addedRules []SecurityGroupRule
	for _, rule := range rules {
		for i := range rule.IpRanges {
			if rule.IpRanges[i].RuleId == "" {
				id, err := generateSecurityGroupRuleID()
				if err != nil {
					return nil, err
				}
				rule.IpRanges[i].RuleId = id
			}
			addedRules = append(addedRules, buildSGRule(sg, rule, rule.IpRanges[i].RuleId, isEgress, region, rule.IpRanges[i].CidrIp, "", rule.IpRanges[i].Description, ec2store.GroupPair{}, ""))
		}
		for i := range rule.Ipv6Ranges {
			if rule.Ipv6Ranges[i].RuleId == "" {
				id, err := generateSecurityGroupRuleID()
				if err != nil {
					return nil, err
				}
				rule.Ipv6Ranges[i].RuleId = id
			}
			addedRules = append(addedRules, buildSGRule(sg, rule, rule.Ipv6Ranges[i].RuleId, isEgress, region, "", rule.Ipv6Ranges[i].CidrIp, rule.Ipv6Ranges[i].Description, ec2store.GroupPair{}, ""))
		}
		for i := range rule.UserIdGroupPairs {
			if rule.UserIdGroupPairs[i].RuleId == "" {
				id, err := generateSecurityGroupRuleID()
				if err != nil {
					return nil, err
				}
				rule.UserIdGroupPairs[i].RuleId = id
			}
			addedRules = append(addedRules, buildSGRule(sg, rule, rule.UserIdGroupPairs[i].RuleId, isEgress, region, "", "", rule.UserIdGroupPairs[i].Description, rule.UserIdGroupPairs[i], ""))
		}
		for i := range rule.PrefixListIds {
			if rule.PrefixListIds[i].RuleId == "" {
				id, err := generateSecurityGroupRuleID()
				if err != nil {
					return nil, err
				}
				rule.PrefixListIds[i].RuleId = id
			}
			addedRules = append(addedRules, buildSGRule(sg, rule, rule.PrefixListIds[i].RuleId, isEgress, region, "", "", rule.PrefixListIds[i].Description, ec2store.GroupPair{}, rule.PrefixListIds[i].PrefixListId))
		}
	}

	if isEgress {
		sg.IpPermissionsEgress = mergeIPRules(sg.IpPermissionsEgress, rules...)
	} else {
		sg.IpPermissions = mergeIPRules(sg.IpPermissions, rules...)
	}

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return &AuthorizeResult{SecurityGroupRules: addedRules}, nil
}

// revokeSecurityGroupRulesCore removes rules from a security group.
// Returns InvalidPermission.NotFound when no rules were actually removed
// (AWS behaviour for non-default VPCs).
func (s *EC2Service) revokeSecurityGroupRulesCore(store *ec2store.EC2Store, sg *ec2store.SecurityGroup, rules []ec2store.IPRule, isEgress bool) (*RevokeResult, error) {
	var existingRules []ec2store.IPRule
	if isEgress {
		existingRules = sg.IpPermissionsEgress
	} else {
		existingRules = sg.IpPermissions
	}

	updated, removedCount := removeIPRules(existingRules, rules...)
	if removedCount == 0 {
		return nil, awserrors.NewAWSError("InvalidPermission.NotFound",
			"The specified rules were not found in the security group",
			http.StatusBadRequest)
	}

	if isEgress {
		sg.IpPermissionsEgress = updated
	} else {
		sg.IpPermissions = updated
	}

	var revoked []SecurityGroupRule
	for _, rule := range rules {
		for _, ip := range rule.IpRanges {
			revoked = append(revoked, SecurityGroupRule{
				IpProtocol:  rule.IpProtocol,
				FromPort:    rule.FromPort,
				ToPort:      rule.ToPort,
				IsEgress:    isEgress,
				CidrIpv4:    ip.CidrIp,
				Description: ip.Description,
			})
		}
		for _, ip := range rule.Ipv6Ranges {
			revoked = append(revoked, SecurityGroupRule{
				IpProtocol:  rule.IpProtocol,
				FromPort:    rule.FromPort,
				ToPort:      rule.ToPort,
				IsEgress:    isEgress,
				CidrIpv6:    ip.CidrIp,
				Description: ip.Description,
			})
		}
		for _, pair := range rule.UserIdGroupPairs {
			revoked = append(revoked, SecurityGroupRule{
				IpProtocol:            rule.IpProtocol,
				FromPort:              rule.FromPort,
				ToPort:                rule.ToPort,
				IsEgress:              isEgress,
				ReferencedGroupId:     pair.GroupId,
				ReferencedGroupName:   pair.GroupName,
				ReferencedGroupUserId: pair.UserId,
				Description:           pair.Description,
			})
		}
		for _, pl := range rule.PrefixListIds {
			revoked = append(revoked, SecurityGroupRule{
				IpProtocol:   rule.IpProtocol,
				FromPort:     rule.FromPort,
				ToPort:       rule.ToPort,
				IsEgress:     isEgress,
				PrefixListId: pl.PrefixListId,
				Description:  pl.Description,
			})
		}
	}

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return &RevokeResult{
		RevokedSecurityGroupRules: revoked,
	}, nil
}

// buildSGRule constructs a SecurityGroupRule from its components.
func buildSGRule(sg *ec2store.SecurityGroup, rule ec2store.IPRule, ruleId string, isEgress bool, region string, cidrIpv4, cidrIpv6, description string, pair ec2store.GroupPair, prefixListId string) SecurityGroupRule {
	arn := fmt.Sprintf("arn:aws:ec2:%s:%s:security-group-rule/%s", region, sg.OwnerId, ruleId)
	r := SecurityGroupRule{
		RuleId:               ruleId,
		IsEgress:             isEgress,
		SecurityGroupRuleArn: arn,
		GroupId:              sg.GroupId,
		GroupOwnerId:         sg.OwnerId,
		IpProtocol:           rule.IpProtocol,
		FromPort:             rule.FromPort,
		ToPort:               rule.ToPort,
		CidrIpv4:             cidrIpv4,
		CidrIpv6:             cidrIpv6,
		Description:          description,
		PrefixListId:         prefixListId,
	}
	if pair.GroupId != "" {
		r.ReferencedGroupId = pair.GroupId
		r.ReferencedGroupName = pair.GroupName
		r.ReferencedGroupUserId = pair.UserId
	}
	return r
}

// resolveSecurityGroupCore finds a security group by GroupId or GroupName.
func resolveSecurityGroupCore(store *ec2store.EC2Store, params map[string]interface{}) (*ec2store.SecurityGroup, error) {
	return resolveSecurityGroup(store, params)
}

// sgRuleSummary returns a human-readable summary of a rule for error messages.
func sgRuleSummary(rule ec2store.IPRule) string {
	parts := []string{rule.IpProtocol}
	if rule.FromPort >= 0 {
		parts = append(parts, fmt.Sprintf("%d", rule.FromPort))
	}
	if rule.ToPort >= 0 {
		parts = append(parts, fmt.Sprintf("%d", rule.ToPort))
	}
	return strings.Join(parts, ":")
}
