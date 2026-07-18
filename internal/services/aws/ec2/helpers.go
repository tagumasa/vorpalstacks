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

// validateCIDRBlock validates that the given string is a valid IPv4 CIDR block.
func validateCIDRBlock(cidr string) error {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter cidrBlock is invalid. This is not a valid CIDR block.", cidr),
			http.StatusBadRequest)
	}
	if ip.To4() == nil {
		return awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter cidrBlock is invalid. This is not a valid IPv4 CIDR block.", cidr),
			http.StatusBadRequest)
	}
	return nil
}

// validateSubnetInVPC checks that the subnet CIDR falls within the VPC CIDR range.
func validateSubnetInVPC(subnetCIDR, vpcCIDR string) error {
	_, subnetNet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return validateCIDRBlock(subnetCIDR)
	}
	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if err != nil {
		return validateCIDRBlock(vpcCIDR)
	}

	// Subnet prefix must be longer than or equal to VPC prefix.
	onesSubnet, _ := subnetNet.Mask.Size()
	onesVPC, _ := vpcNet.Mask.Size()
	if onesSubnet < onesVPC {
		return awserrors.NewAWSError("InvalidSubnet.Range",
			fmt.Sprintf("The CIDR '%s' does not fall within the VPC CIDR range '%s'.", subnetCIDR, vpcCIDR),
			http.StatusBadRequest)
	}

	// Subnet network address must be within the VPC network.
	if !vpcNet.Contains(subnetNet.IP) {
		return awserrors.NewAWSError("InvalidSubnet.Range",
			fmt.Sprintf("The CIDR '%s' does not fall within the VPC CIDR range '%s'.", subnetCIDR, vpcCIDR),
			http.StatusBadRequest)
	}

	return nil
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

// rulesMatch checks if two rules share the same protocol and port range.
func rulesMatch(a, b ec2store.IPRule) bool {
	return a.IpProtocol == b.IpProtocol && a.FromPort == b.FromPort && a.ToPort == b.ToPort
}

// mergeIPRules merges new rules into the existing set. Rules with matching
// protocol/port have their IP ranges consolidated — duplicate CIDRs update
// the Description rather than being silently dropped.
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

// mergeIPRanges merges IP ranges. If a CIDR already exists, the Description
// is updated when the incoming Description is non-empty (AWS idempotent behaviour).
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

// mergeGroupPairs merges user ID group pairs, avoiding duplicates by GroupId+UserId.
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

// removeIPRules removes specific ranges/pairs from matching rules. If all
// ranges are removed from a rule, the rule itself is removed.
func removeIPRules(existing []ec2store.IPRule, toRemove ...ec2store.IPRule) []ec2store.IPRule {
	result := make([]ec2store.IPRule, 0, len(existing))
	for _, er := range existing {
		for _, nr := range toRemove {
			if rulesMatch(er, nr) {
				er.IpRanges = removeIPRanges(er.IpRanges, nr.IpRanges)
				er.Ipv6Ranges = removeIPRanges(er.Ipv6Ranges, nr.Ipv6Ranges)
				er.UserIdGroupPairs = removeGroupPairs(er.UserIdGroupPairs, nr.UserIdGroupPairs)
			}
		}
		// Keep the rule only if it still has sources.
		if len(er.IpRanges) > 0 || len(er.Ipv6Ranges) > 0 || len(er.UserIdGroupPairs) > 0 {
			result = append(result, er)
		}
	}
	return result
}

// removeIPRanges removes IP ranges matching the given CIDRs.
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

// removeGroupPairs removes user ID group pairs matching by GroupId+UserId.
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

// validateIPRule validates the protocol/port combination for a security group rule.
// EC2 rules: tcp/udp accept ports 0-65535 or -1 (all); icmp/icmpv6 accept type/code
// 0-255 or -1; protocol "-1" or numeric protocols must have ports -1.
func validateIPRule(rule ec2store.IPRule) error {
	proto := rule.IpProtocol
	from := rule.FromPort
	to := rule.ToPort

	switch proto {
	case "-1":
		if from != -1 || to != -1 {
			return awserrors.NewAWSError("InvalidParameterValue",
				"Ports are not allowed for protocol '-1'", http.StatusBadRequest)
		}
	case "tcp", "udp":
		if !validPortRange(from, to) {
			return awserrors.NewAWSError("InvalidParameterValue",
				fmt.Sprintf("Invalid port range %d-%d for protocol %s", from, to, proto),
				http.StatusBadRequest)
		}
	case "icmp", "icmpv6":
		if !validICMPRange(from, to) {
			return awserrors.NewAWSError("InvalidParameterValue",
				fmt.Sprintf("Invalid ICMP type/code %d-%d for protocol %s", from, to, proto),
				http.StatusBadRequest)
		}
	default:
		// Numeric protocols (e.g. "50" ESP, "51" AH, "89" OSPF) must not specify ports.
		if from != -1 || to != -1 {
			return awserrors.NewAWSError("InvalidParameterValue",
				fmt.Sprintf("Ports are not allowed for protocol '%s'", proto),
				http.StatusBadRequest)
		}
	}
	return nil
}

// validPortRange checks that from/to form a valid TCP/UDP port range (-1 means all).
func validPortRange(from, to int64) bool {
	if from == -1 && to == -1 {
		return true
	}
	if from < 0 || from > 65535 || to < 0 || to > 65535 {
		return false
	}
	return from <= to
}

// validICMPRange checks that from/to form a valid ICMP type/code pair (-1 means all).
func validICMPRange(from, to int64) bool {
	if from == -1 && to == -1 {
		return true
	}
	if from < -1 || from > 255 || to < -1 || to > 255 {
		return false
	}
	return true
}
