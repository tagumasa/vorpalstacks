package ec2

import (
	"fmt"
	"strconv"

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

// parseEC2Tags extracts EC2 tags from request parameters.
// EC2 uses Tag.N.Key / Tag.N.Value format in the Query protocol.
func parseEC2Tags(params map[string]interface{}) []types.Tag {
	return tags.ParseTagsWithPrefix(params, "Tag")
}

// parseInt64Port parses a port string to int64, returning -1 for empty/invalid.
// EC2 uses -1 to mean "all ports" (with protocol -1 meaning all protocols).
func parseInt64Port(s string) int64 {
	if s == "" {
		return -1
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return n
}

// parseIPRules extracts IP permission rules from request parameters.
// AWS SDK sends nested format:
//   - IpPermissions.N.IpRanges.M.CidrIp
//   - IpPermissions.N.IpRanges.M.Description
//   - IpPermissions.N.Groups.M.GroupId
//   - IpPermissions.N.Groups.M.UserId
//   - IpPermissions.N.Ipv6Ranges.M.CidrIpv6
func parseIPRules(params map[string]interface{}, prefix string) []ec2store.IPRule {
	var rules []ec2store.IPRule

	for i := 1; ; i++ {
		ipProtocol := request.GetStringParam(params, fmt.Sprintf("%s.%d.IpProtocol", prefix, i))
		if ipProtocol == "" {
			break
		}

		rule := ec2store.IPRule{
			IpProtocol: ipProtocol,
			FromPort:   parseInt64Port(request.GetStringParam(params, fmt.Sprintf("%s.%d.FromPort", prefix, i))),
			ToPort:     parseInt64Port(request.GetStringParam(params, fmt.Sprintf("%s.%d.ToPort", prefix, i))),
		}

		ipPrefix := fmt.Sprintf("%s.%d.IpRanges", prefix, i)
		for j := 1; ; j++ {
			cidr := request.GetStringParam(params, fmt.Sprintf("%s.%d.CidrIp", ipPrefix, j))
			if cidr == "" {
				break
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
			ipRange := ec2store.IPRange{CidrIp: cidr}
			if desc := request.GetStringParam(params, fmt.Sprintf("%s.%d.Description", ipv6Prefix, j)); desc != "" {
				ipRange.Description = desc
			}
			rule.Ipv6Ranges = append(rule.Ipv6Ranges, ipRange)
		}

		groupPrefix := fmt.Sprintf("%s.%d.Groups", prefix, i)
		for j := 1; ; j++ {
			groupId := request.GetStringParam(params, fmt.Sprintf("%s.%d.GroupId", groupPrefix, j))
			if groupId == "" {
				groupName := request.GetStringParam(params, fmt.Sprintf("%s.%d.GroupName", groupPrefix, j))
				if groupName == "" {
					break
				}
			}
			pair := ec2store.GroupPair{
				GroupId:   groupId,
				GroupName: request.GetStringParam(params, fmt.Sprintf("%s.%d.GroupName", groupPrefix, j)),
				UserId:    request.GetStringParam(params, fmt.Sprintf("%s.%d.UserId", groupPrefix, j)),
				VpcId:     request.GetStringParam(params, fmt.Sprintf("%s.%d.VpcId", groupPrefix, j)),
			}
			rule.UserIdGroupPairs = append(rule.UserIdGroupPairs, pair)
		}

		rules = append(rules, rule)
	}

	return rules
}

// ipRuleEquals checks if two IPRules match on protocol, ports, and all
// ranges/groups. This is used for duplicate detection and rule revocation.
func ipRuleEquals(a, b ec2store.IPRule) bool {
	if a.IpProtocol != b.IpProtocol || a.FromPort != b.FromPort || a.ToPort != b.ToPort {
		return false
	}
	if len(a.IpRanges) != len(b.IpRanges) || len(a.UserIdGroupPairs) != len(b.UserIdGroupPairs) || len(a.Ipv6Ranges) != len(b.Ipv6Ranges) {
		return false
	}
	for _, ar := range a.IpRanges {
		found := false
		for _, br := range b.IpRanges {
			if ar.CidrIp == br.CidrIp {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for _, ag := range a.UserIdGroupPairs {
		found := false
		for _, bg := range b.UserIdGroupPairs {
			if ag.GroupId == bg.GroupId && ag.UserId == bg.UserId {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for _, ar := range a.Ipv6Ranges {
		found := false
		for _, br := range b.Ipv6Ranges {
			if ar.CidrIp == br.CidrIp {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// mergeIPRules appends new rules to an existing rule set, avoiding full duplicates.
func mergeIPRules(existing []ec2store.IPRule, newRules ...ec2store.IPRule) []ec2store.IPRule {
	for _, nr := range newRules {
		duplicate := false
		for _, er := range existing {
			if ipRuleEquals(er, nr) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			existing = append(existing, nr)
		}
	}
	return existing
}

// removeIPRules removes rules matching the given criteria from an existing rule set.
func removeIPRules(existing []ec2store.IPRule, toRemove ...ec2store.IPRule) []ec2store.IPRule {
	result := make([]ec2store.IPRule, 0, len(existing))
	for _, er := range existing {
		removed := false
		for _, nr := range toRemove {
			if ipRuleEquals(er, nr) {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, er)
		}
	}
	return result
}
