package ec2

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
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
		return nil, translateStoreError(err)
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

	page, err := paginateEC2(filtered, nextToken, maxResults, func(sg *ec2store.SecurityGroup) string {
		return sg.GroupId
	})
	if err != nil {
		return nil, err
	}
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

// resolveSecurityGroupCore is the single security-group resolution path for
// the Delete and Authorize/Revoke operations: GroupId is preferred when
// present, otherwise the group is looked up by GroupName region-wide (this
// platform has no default VPC, so GroupName is never scoped to one).
func (s *EC2Service) resolveSecurityGroupCore(store *ec2store.EC2Store, groupID, groupName string) (*ec2store.SecurityGroup, error) {
	if groupID != "" {
		sg, err := store.GetSecurityGroup(groupID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return sg, nil
	}

	if groupName == "" {
		return nil, awserrors.NewMissingParameter("GroupId or GroupName is required")
	}

	sgs, err := store.ListSecurityGroups()
	if err != nil {
		return nil, translateStoreError(err)
	}
	for _, sg := range sgs {
		if sg.GroupName == groupName {
			return sg, nil
		}
	}
	return nil, awserrors.NewAWSError("InvalidGroup.NotFound", "The security group does not exist", http.StatusNotFound)
}

// authorizeSecurityGroupRulesCore adds rules to a security group and returns
// the newly created SecurityGroupRules with generated rule IDs. An identical
// existing source entry (same protocol/port plus identical CIDR, group pair,
// or prefix list) yields InvalidPermission.Duplicate, per the AWS IpRange
// documentation on duplicate rule errors.
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

	var existing []ec2store.IPRule
	if isEgress {
		existing = sg.IpPermissionsEgress
	} else {
		existing = sg.IpPermissions
	}
	if dup := findDuplicateEntry(existing, rules); dup != "" {
		return nil, awserrors.NewAWSError("InvalidPermission.Duplicate",
			fmt.Sprintf("The permission '%s' has already been granted", dup),
			http.StatusBadRequest)
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

// findDuplicateEntry returns a description of the first requested source
// entry that already exists verbatim in the security group, or "" when all
// entries are new. A different description on an otherwise identical entry
// also counts as a duplicate (descriptions are updated via
// ModifySecurityGroupRules, not by re-authorising).
func findDuplicateEntry(existing, requested []ec2store.IPRule) string {
	for _, er := range existing {
		for _, nr := range requested {
			if !rulesMatch(er, nr) {
				continue
			}
			for _, x := range nr.IpRanges {
				for _, e := range er.IpRanges {
					if e.CidrIp == x.CidrIp {
						return fmt.Sprintf("rule with ipProtocol=%s ports=%d-%d cidr=%s already exists", nr.IpProtocol, nr.FromPort, nr.ToPort, x.CidrIp)
					}
				}
			}
			for _, x := range nr.Ipv6Ranges {
				for _, e := range er.Ipv6Ranges {
					if e.CidrIp == x.CidrIp {
						return fmt.Sprintf("rule with ipProtocol=%s ports=%d-%d cidrIpv6=%s already exists", nr.IpProtocol, nr.FromPort, nr.ToPort, x.CidrIp)
					}
				}
			}
			for _, x := range nr.UserIdGroupPairs {
				for _, e := range er.UserIdGroupPairs {
					if e.GroupId == x.GroupId && e.UserId == x.UserId {
						return fmt.Sprintf("rule with ipProtocol=%s ports=%d-%d groupId=%s already exists", nr.IpProtocol, nr.FromPort, nr.ToPort, x.GroupId)
					}
				}
			}
			for _, x := range nr.PrefixListIds {
				for _, e := range er.PrefixListIds {
					if e.PrefixListId == x.PrefixListId {
						return fmt.Sprintf("rule with ipProtocol=%s ports=%d-%d prefixListId=%s already exists", nr.IpProtocol, nr.FromPort, nr.ToPort, x.PrefixListId)
					}
				}
			}
		}
	}
	return ""
}

// revokeSecurityGroupRulesCore removes rules from a security group, either by
// permission match or by SecurityGroupRuleIds. It returns the rules that were
// actually revoked (identified by their stored rule IDs) and the requested
// permissions that matched nothing (unknownIpPermissionSet). Returns
// InvalidPermission.NotFound when nothing was removed (AWS behaviour for
// non-default VPCs).
func (s *EC2Service) revokeSecurityGroupRulesCore(store *ec2store.EC2Store, sg *ec2store.SecurityGroup, rules []ec2store.IPRule, ruleIDs []string, isEgress bool) (*RevokeResult, error) {
	if len(rules) == 0 && len(ruleIDs) == 0 {
		return nil, awserrors.NewMissingParameter("IpPermissions, the legacy rule parameters (IpProtocol, CidrIp, ...) or SecurityGroupRuleIds is required")
	}
	var existingRules []ec2store.IPRule
	if isEgress {
		existingRules = sg.IpPermissionsEgress
	} else {
		existingRules = sg.IpPermissions
	}

	// Requested permissions whose protocol/port and source match no existing
	// entry are reported in unknownIpPermissionSet.
	var unknown []ec2store.IPRule
	for _, nr := range rules {
		if !permissionExists(existingRules, nr) {
			unknown = append(unknown, nr)
		}
	}

	var revoked []SecurityGroupRule
	updated := make([]ec2store.IPRule, 0, len(existingRules))
	removedCount := 0

	for _, er := range existingRules {
		keep := er
		if len(ruleIDs) > 0 {
			keep, revoked = removeByRuleIDs(er, ruleIDs, revoked, isEgress, sg.GroupId)
		} else {
			keep, revoked = removeMatchingEntries(er, rules, revoked, isEgress, sg.GroupId)
		}
		removedCount += ruleEntryCount(er) - ruleEntryCount(keep)
		if ruleEntryCount(keep) > 0 {
			updated = append(updated, keep)
		}
	}

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

	if err := store.UpdateSecurityGroup(sg); err != nil {
		return nil, translateStoreError(err)
	}

	return &RevokeResult{
		RevokedSecurityGroupRules: revoked,
		UnknownIpPermissions:      unknown,
	}, nil
}

// permissionExists reports whether the requested permission's protocol/port
// matches an existing rule with at least one identical source entry.
func permissionExists(existing []ec2store.IPRule, nr ec2store.IPRule) bool {
	for _, er := range existing {
		if !rulesMatch(er, nr) {
			continue
		}
		for _, r := range er.IpRanges {
			for _, x := range nr.IpRanges {
				if r.CidrIp == x.CidrIp {
					return true
				}
			}
		}
		for _, r := range er.Ipv6Ranges {
			for _, x := range nr.Ipv6Ranges {
				if r.CidrIp == x.CidrIp {
					return true
				}
			}
		}
		for _, g := range er.UserIdGroupPairs {
			for _, x := range nr.UserIdGroupPairs {
				if g.GroupId == x.GroupId && g.UserId == x.UserId {
					return true
				}
			}
		}
		for _, p := range er.PrefixListIds {
			for _, x := range nr.PrefixListIds {
				if p.PrefixListId == x.PrefixListId {
					return true
				}
			}
		}
	}
	return false
}

// ruleEntryCount counts the source entries of a rule (ranges, pairs, prefix
// lists in both address families).
func ruleEntryCount(r ec2store.IPRule) int {
	return len(r.IpRanges) + len(r.Ipv6Ranges) + len(r.UserIdGroupPairs) + len(r.PrefixListIds)
}

// removeEntriesByID strips entries whose stored rule ID matches any of the
// requested IDs and appends a revoked-rule record for each removal. It is
// generic over the source-family entry type (IPv4 ranges, IPv6 ranges,
// group pairs, prefix list IDs) so that every family shares one traversal;
// per-family copies previously diverged on which slice they iterated and
// reported IPv4 ranges as IPv6.
func removeEntriesByID[T any](src []T, ruleIDs []string, revoked *[]SecurityGroupRule, ruleIDOf func(T) string, toRevoked func(T) SecurityGroupRule) []T {
	keep := make([]T, 0, len(src))
	for _, e := range src {
		if containsString(ruleIDs, ruleIDOf(e)) {
			*revoked = append(*revoked, toRevoked(e))
			continue
		}
		keep = append(keep, e)
	}
	return keep
}

// removeByRuleIDs strips entries whose stored RuleId matches any of the given
// IDs and appends them to the revoked list.
func removeByRuleIDs(er ec2store.IPRule, ruleIDs []string, revoked []SecurityGroupRule, isEgress bool, groupID string) (ec2store.IPRule, []SecurityGroupRule) {
	keep := er
	keep.IpRanges = removeEntriesByID(er.IpRanges, ruleIDs, &revoked,
		func(e ec2store.IPRange) string { return e.RuleId },
		func(e ec2store.IPRange) SecurityGroupRule {
			return buildRevokedEntry(er, e.RuleId, isEgress, groupID, e.CidrIp, "", e.Description, "", "")
		})
	keep.Ipv6Ranges = removeEntriesByID(er.Ipv6Ranges, ruleIDs, &revoked,
		func(e ec2store.IPRange) string { return e.RuleId },
		func(e ec2store.IPRange) SecurityGroupRule {
			return buildRevokedEntry(er, e.RuleId, isEgress, groupID, "", e.CidrIp, e.Description, "", "")
		})
	keep.UserIdGroupPairs = removeEntriesByID(er.UserIdGroupPairs, ruleIDs, &revoked,
		func(e ec2store.GroupPair) string { return e.RuleId },
		func(e ec2store.GroupPair) SecurityGroupRule {
			return buildRevokedEntry(er, e.RuleId, isEgress, groupID, "", "", e.Description, e.GroupId, "")
		})
	keep.PrefixListIds = removeEntriesByID(er.PrefixListIds, ruleIDs, &revoked,
		func(e ec2store.PrefixListId) string { return e.RuleId },
		func(e ec2store.PrefixListId) SecurityGroupRule {
			return buildRevokedEntry(er, e.RuleId, isEgress, groupID, "", "", e.Description, "", e.PrefixListId)
		})
	return keep, revoked
}

// removeMatchingEntries strips entries that match the requested permissions
// (protocol/port plus identical source) and appends them to the revoked list.
func removeMatchingEntries(er ec2store.IPRule, toRemove []ec2store.IPRule, revoked []SecurityGroupRule, isEgress bool, groupID string) (ec2store.IPRule, []SecurityGroupRule) {
	keep := er
	for _, nr := range toRemove {
		if !rulesMatch(er, nr) {
			continue
		}
		keep.IpRanges, revoked = removeIPRangesMatching(er.IpRanges, nr.IpRanges, revoked, isEgress, groupID, er)
		keep.Ipv6Ranges, revoked = removeIPv6RangesMatching(er.Ipv6Ranges, nr.Ipv6Ranges, revoked, isEgress, groupID, er)
		keep.UserIdGroupPairs, revoked = removeGroupPairsMatching(er.UserIdGroupPairs, nr.UserIdGroupPairs, revoked, isEgress, groupID, er)
		keep.PrefixListIds, revoked = removePrefixListIDsMatching(er.PrefixListIds, nr.PrefixListIds, revoked, isEgress, groupID, er)
	}
	return keep, revoked
}

// buildRevokedEntry constructs a revoked-rule record for a removed entry,
// preserving the stored rule ID.
func buildRevokedEntry(rule ec2store.IPRule, ruleID string, isEgress bool, groupID string, cidrIpv4, cidrIpv6, description, refGroupID, prefixListID string) SecurityGroupRule {
	return SecurityGroupRule{
		RuleId:            ruleID,
		IsEgress:          isEgress,
		GroupId:           groupID,
		IpProtocol:        rule.IpProtocol,
		FromPort:          rule.FromPort,
		ToPort:            rule.ToPort,
		CidrIpv4:          cidrIpv4,
		CidrIpv6:          cidrIpv6,
		Description:       description,
		ReferencedGroupId: refGroupID,
		PrefixListId:      prefixListID,
	}
}

// removeIPRangesMatching removes IPv4 ranges identical to any requested range.
func removeIPRangesMatching(ranges []ec2store.IPRange, toRemove []ec2store.IPRange, revoked []SecurityGroupRule, isEgress bool, groupID string, rule ec2store.IPRule) ([]ec2store.IPRange, []SecurityGroupRule) {
	result := make([]ec2store.IPRange, 0, len(ranges))
	for _, er := range ranges {
		found := false
		for _, nr := range toRemove {
			if er.CidrIp == nr.CidrIp {
				found = true
				break
			}
		}
		if found {
			revoked = append(revoked, buildRevokedEntry(rule, er.RuleId, isEgress, groupID, er.CidrIp, "", er.Description, "", ""))
			continue
		}
		result = append(result, er)
	}
	return result, revoked
}

// removeIPv6RangesMatching removes IPv6 ranges identical to any requested range.
func removeIPv6RangesMatching(ranges []ec2store.IPRange, toRemove []ec2store.IPRange, revoked []SecurityGroupRule, isEgress bool, groupID string, rule ec2store.IPRule) ([]ec2store.IPRange, []SecurityGroupRule) {
	result := make([]ec2store.IPRange, 0, len(ranges))
	for _, er := range ranges {
		found := false
		for _, nr := range toRemove {
			if er.CidrIp == nr.CidrIp {
				found = true
				break
			}
		}
		if found {
			revoked = append(revoked, buildRevokedEntry(rule, er.RuleId, isEgress, groupID, "", er.CidrIp, er.Description, "", ""))
			continue
		}
		result = append(result, er)
	}
	return result, revoked
}

// removeGroupPairsMatching removes group pairs identical to any requested pair.
func removeGroupPairsMatching(pairs []ec2store.GroupPair, toRemove []ec2store.GroupPair, revoked []SecurityGroupRule, isEgress bool, groupID string, rule ec2store.IPRule) ([]ec2store.GroupPair, []SecurityGroupRule) {
	result := make([]ec2store.GroupPair, 0, len(pairs))
	for _, er := range pairs {
		found := false
		for _, nr := range toRemove {
			if er.GroupId == nr.GroupId && er.UserId == nr.UserId {
				found = true
				break
			}
		}
		if found {
			revoked = append(revoked, buildRevokedEntry(rule, er.RuleId, isEgress, groupID, "", "", er.Description, er.GroupId, ""))
			continue
		}
		result = append(result, er)
	}
	return result, revoked
}

// removePrefixListIDsMatching removes prefix list entries identical to any
// requested entry.
func removePrefixListIDsMatching(pls []ec2store.PrefixListId, toRemove []ec2store.PrefixListId, revoked []SecurityGroupRule, isEgress bool, groupID string, rule ec2store.IPRule) ([]ec2store.PrefixListId, []SecurityGroupRule) {
	result := make([]ec2store.PrefixListId, 0, len(pls))
	for _, er := range pls {
		found := false
		for _, nr := range toRemove {
			if er.PrefixListId == nr.PrefixListId {
				found = true
				break
			}
		}
		if found {
			revoked = append(revoked, buildRevokedEntry(rule, er.RuleId, isEgress, groupID, "", "", er.Description, "", er.PrefixListId))
			continue
		}
		result = append(result, er)
	}
	return result, revoked
}

// containsString reports whether the slice contains the target.
func containsString(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
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
