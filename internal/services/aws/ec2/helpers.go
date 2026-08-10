package ec2

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/generators"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	vpcIDPrefix         = "vpc-"
	subnetIDPrefix      = "subnet-"
	securityGroupPrefix = "sg-"
	ec2IDSuffixLen      = 17
)

// GenerateVpcID creates a new VPC ID in the format vpc-<hex>.
func GenerateVpcID() (string, error) {
	return generators.GenerateIDWithPrefix(vpcIDPrefix, ec2IDSuffixLen)
}

// GenerateSubnetID creates a new subnet ID in the format subnet-<hex>.
func GenerateSubnetID() (string, error) {
	return generators.GenerateIDWithPrefix(subnetIDPrefix, ec2IDSuffixLen)
}

// GenerateSecurityGroupID creates a new security group ID in the format sg-<hex>.
func GenerateSecurityGroupID() (string, error) {
	return generators.GenerateIDWithPrefix(securityGroupPrefix, ec2IDSuffixLen)
}

// calculateAvailableIPs returns the number of usable IP addresses in a subnet
// with the given CIDR block. AWS reserves 5 IPs per subnet.
func calculateAvailableIPs(cidrBlock string) int64 {
	_, ipNet, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		return 0
	}
	ones, bits := ipNet.Mask.Size()
	if bits == 0 || ones >= bits {
		return 0
	}
	total := int64(1) << uint(bits-ones)
	available := total - 5
	if available < 0 {
		available = 0
	}
	return available
}

// parseEC2Tags extracts EC2 tags from request parameters.
// AWS SDK Go v2 sends tags as TagSpecification.N.Tag.N.Key/Value;
// older clients (AWS CLI v1) use Tag.N.Key/Value. Both are handled.
func parseEC2Tags(params map[string]interface{}) []types.Tag {
	if tsTags := parseTagSpecification(params); len(tsTags) > 0 {
		return tsTags
	}
	return tags.ParseTagsWithPrefix(params, "Tag")
}

// checkDryRun returns a DryRunOperation error (HTTP 412) when the DryRun
// query parameter is set to true.
func checkDryRun(params map[string]interface{}) error {
	if v := request.GetStringParam(params, "DryRun"); v == "true" {
		return awserrors.NewAWSError("DryRunOperation",
			"Request would have succeeded, but DryRun flag is set.",
			http.StatusPreconditionFailed)
	}
	return nil
}

// parseTagSpecification parses the TagSpecification.N.Tag.N.Key/Value format
// used by AWS SDK Go v2 for EC2 Query protocol requests.
func parseTagSpecification(params map[string]interface{}) []types.Tag {
	var result []types.Tag
	for i := 1; ; i++ {
		tagPrefix := "TagSpecification." + strconv.Itoa(i) + ".Tag"
		found := false
		for j := 1; ; j++ {
			keyKey := tagPrefix + "." + strconv.Itoa(j) + ".Key"
			valueKey := tagPrefix + "." + strconv.Itoa(j) + ".Value"

			key := request.GetStringParam(params, keyKey)
			if key == "" {
				break
			}

			value := request.GetStringParam(params, valueKey)
			result = append(result, types.Tag{Key: key, Value: value})
			found = true
		}
		if !found {
			break
		}
	}
	return result
}

// parseIPRules extracts IP permission rules from request parameters.
// Returns an error if any port value or CIDR is malformed.
func parseIPRules(params map[string]interface{}, prefix string) ([]ec2store.IPRule, error) {
	var rules []ec2store.IPRule

	for i := 1; ; i++ {
		ipProtocol := request.GetStringParam(params, fmt.Sprintf("%s.%d.IpProtocol", prefix, i))
		if ipProtocol == "" {
			break
		}

		fromPort, err := parseInt64Param(request.GetStringParam(params, fmt.Sprintf("%s.%d.FromPort", prefix, i)))
		if err != nil {
			return nil, err
		}
		toPort, err := parseInt64Param(request.GetStringParam(params, fmt.Sprintf("%s.%d.ToPort", prefix, i)))
		if err != nil {
			return nil, err
		}

		rule := ec2store.IPRule{
			IpProtocol: ipProtocol,
			FromPort:   fromPort,
			ToPort:     toPort,
		}

		ipPrefix := fmt.Sprintf("%s.%d.IpRanges", prefix, i)
		for j := 1; ; j++ {
			cidr := request.GetStringParam(params, fmt.Sprintf("%s.%d.CidrIp", ipPrefix, j))
			if cidr == "" {
				break
			}
			if err := validateIPv4CIDRInRule(cidr); err != nil {
				return nil, err
			}
			ipRange := ec2store.IPRange{CidrIp: cidr}
			if desc := request.GetStringParam(params, fmt.Sprintf("%s.%d.Description", ipPrefix, j)); desc != "" {
				ipRange.Description = desc
			}
			rule.IpRanges = append(rule.IpRanges, ipRange)
		}

		ipv6Prefix := fmt.Sprintf("%s.%d.Ipv6Ranges", prefix, i)
		for j := 1; ; j++ {
			cidr := request.GetStringParam(params, fmt.Sprintf("%s.%d.CidrIpv6", ipv6Prefix, j))
			if cidr == "" {
				break
			}
			if err := validateIPv6CIDRInRule(cidr); err != nil {
				return nil, err
			}
			ipRange := ec2store.IPRange{CidrIp: cidr}
			if desc := request.GetStringParam(params, fmt.Sprintf("%s.%d.Description", ipv6Prefix, j)); desc != "" {
				ipRange.Description = desc
			}
			rule.Ipv6Ranges = append(rule.Ipv6Ranges, ipRange)
		}

		groupPrefix := fmt.Sprintf("%s.%d.Groups", prefix, i)
		for j := 1; ; j++ {
			groupId := request.GetStringParam(params, fmt.Sprintf("%s.%d.GroupId", groupPrefix, j))
			groupName := request.GetStringParam(params, fmt.Sprintf("%s.%d.GroupName", groupPrefix, j))
			if groupId == "" && groupName == "" {
				break
			}
			pair := ec2store.GroupPair{
				GroupId:                groupId,
				GroupName:              groupName,
				UserId:                 request.GetStringParam(params, fmt.Sprintf("%s.%d.UserId", groupPrefix, j)),
				VpcId:                  request.GetStringParam(params, fmt.Sprintf("%s.%d.VpcId", groupPrefix, j)),
				Description:            request.GetStringParam(params, fmt.Sprintf("%s.%d.Description", groupPrefix, j)),
				VpcPeeringConnectionId: request.GetStringParam(params, fmt.Sprintf("%s.%d.VpcPeeringConnectionId", groupPrefix, j)),
			}
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, pair)
		}

		plPrefix := fmt.Sprintf("%s.%d.PrefixListIds", prefix, i)
		for j := 1; ; j++ {
			plId := request.GetStringParam(params, fmt.Sprintf("%s.%d.PrefixListId", plPrefix, j))
			if plId == "" {
				break
			}
			pl := ec2store.PrefixListId{PrefixListId: plId}
			if desc := request.GetStringParam(params, fmt.Sprintf("%s.%d.Description", plPrefix, j)); desc != "" {
				pl.Description = desc
			}
			rule.PrefixListIds = append(rule.PrefixListIds, pl)
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// rulesMatch checks if two rules share the same protocol and port range.
func rulesMatch(a, b ec2store.IPRule) bool {
	return a.IpProtocol == b.IpProtocol && a.FromPort == b.FromPort && a.ToPort == b.ToPort
}

// mergeIPRules merges new rules into the existing set.
func mergeIPRules(existing []ec2store.IPRule, newRules ...ec2store.IPRule) []ec2store.IPRule {
	for _, nr := range newRules {
		merged := false
		for i := range existing {
			if rulesMatch(existing[i], nr) {
				existing[i].IpRanges = mergeIPRanges(existing[i].IpRanges, nr.IpRanges)
				existing[i].Ipv6Ranges = mergeIPRanges(existing[i].Ipv6Ranges, nr.Ipv6Ranges)
				existing[i].UserIdGroupPairs = mergeGroupPairs(existing[i].UserIdGroupPairs, nr.UserIdGroupPairs)
				merged = true
				break
			}
		}
		if !merged {
			existing = append(existing, nr)
		}
	}
	return existing
}

func mergeIPRanges(existing []ec2store.IPRange, newRanges []ec2store.IPRange) []ec2store.IPRange {
	for _, nr := range newRanges {
		found := false
		for i := range existing {
			if existing[i].CidrIp == nr.CidrIp {
				if nr.Description != "" {
					existing[i].Description = nr.Description
				}
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, nr)
		}
	}
	return existing
}

func mergeGroupPairs(existing []ec2store.GroupPair, newPairs []ec2store.GroupPair) []ec2store.GroupPair {
	for _, np := range newPairs {
		found := false
		for _, ep := range existing {
			if ep.GroupId == np.GroupId && ep.UserId == np.UserId {
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, np)
		}
	}
	return existing
}

// removeIPRules removes specific ranges/pairs from matching rules.
// Returns the filtered result and the count of rules actually removed.
func removeIPRules(existing []ec2store.IPRule, toRemove ...ec2store.IPRule) ([]ec2store.IPRule, int) {
	result := make([]ec2store.IPRule, 0, len(existing))
	removedCount := 0
	for _, er := range existing {
		beforeRanges := len(er.IpRanges) + len(er.Ipv6Ranges) + len(er.UserIdGroupPairs) + len(er.PrefixListIds)
		for _, nr := range toRemove {
			if rulesMatch(er, nr) {
				er.IpRanges = removeIPRanges(er.IpRanges, nr.IpRanges)
				er.Ipv6Ranges = removeIPRanges(er.Ipv6Ranges, nr.Ipv6Ranges)
				er.UserIdGroupPairs = removeGroupPairs(er.UserIdGroupPairs, nr.UserIdGroupPairs)
				er.PrefixListIds = removePrefixListIds(er.PrefixListIds, nr.PrefixListIds)
			}
		}
		afterRanges := len(er.IpRanges) + len(er.Ipv6Ranges) + len(er.UserIdGroupPairs) + len(er.PrefixListIds)
		if afterRanges < beforeRanges {
			removedCount += beforeRanges - afterRanges
		}
		if afterRanges > 0 {
			result = append(result, er)
		}
	}
	return result, removedCount
}

func removeIPRanges(existing []ec2store.IPRange, toRemove []ec2store.IPRange) []ec2store.IPRange {
	result := make([]ec2store.IPRange, 0, len(existing))
	for _, er := range existing {
		found := false
		for _, nr := range toRemove {
			if er.CidrIp == nr.CidrIp {
				found = true
				break
			}
		}
		if !found {
			result = append(result, er)
		}
	}
	return result
}

func removeGroupPairs(existing []ec2store.GroupPair, toRemove []ec2store.GroupPair) []ec2store.GroupPair {
	result := make([]ec2store.GroupPair, 0, len(existing))
	for _, er := range existing {
		found := false
		for _, nr := range toRemove {
			if er.GroupId == nr.GroupId && er.UserId == nr.UserId {
				found = true
				break
			}
		}
		if !found {
			result = append(result, er)
		}
	}
	return result
}

func removePrefixListIds(existing []ec2store.PrefixListId, toRemove []ec2store.PrefixListId) []ec2store.PrefixListId {
	result := make([]ec2store.PrefixListId, 0, len(existing))
	for _, er := range existing {
		found := false
		for _, nr := range toRemove {
			if er.PrefixListId == nr.PrefixListId {
				found = true
				break
			}
		}
		if !found {
			result = append(result, er)
		}
	}
	return result
}
