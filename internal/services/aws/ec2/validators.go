package ec2

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	ec2store "vorpalstacks/internal/store/aws/ec2"
)

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

// canonicalizeCIDR normalises a CIDR block to its canonical network form.
// AWS modifies the specified CIDR block to its canonical form; for example,
// if you specify 100.68.0.18/18, AWS modifies it to 100.68.0.0/18.
func canonicalizeCIDR(cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter cidrBlock is invalid. This is not a valid CIDR block.", cidr),
			http.StatusBadRequest)
	}
	return ipNet.String(), nil
}

// validateIPv4CIDRInRule validates a CIDR within a security group IP range.
func validateIPv4CIDRInRule(cidr string) error {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter ipRange is invalid. This is not a valid CIDR block.", cidr),
			http.StatusBadRequest)
	}
	if ip.To4() == nil {
		return awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter ipRange is invalid. This is not a valid IPv4 CIDR block.", cidr),
			http.StatusBadRequest)
	}
	return nil
}

// validateIPv6CIDRInRule validates an IPv6 CIDR within a security group IP range.
func validateIPv6CIDRInRule(cidr string) error {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter ipv6Range is invalid. This is not a valid CIDR block.", cidr),
			http.StatusBadRequest)
	}
	if ip.To4() != nil {
		return awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Value (%s) for parameter ipv6Range is invalid. This is not a valid IPv6 CIDR block.", cidr),
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

	onesSubnet, _ := subnetNet.Mask.Size()
	onesVPC, _ := vpcNet.Mask.Size()
	if onesSubnet < onesVPC {
		return awserrors.NewAWSError("InvalidSubnet.Range",
			fmt.Sprintf("The CIDR '%s' does not fall within the VPC CIDR range '%s'.", subnetCIDR, vpcCIDR),
			http.StatusBadRequest)
	}

	if !vpcNet.Contains(subnetNet.IP) {
		return awserrors.NewAWSError("InvalidSubnet.Range",
			fmt.Sprintf("The CIDR '%s' does not fall within the VPC CIDR range '%s'.", subnetCIDR, vpcCIDR),
			http.StatusBadRequest)
	}

	return nil
}

// validateSubnetCIDROverlap checks if the new CIDR overlaps with any existing
// subnet in the same VPC.
func validateSubnetCIDROverlap(newCIDR string, existingSubnets []*ec2store.Subnet) error {
	_, newNet, err := net.ParseCIDR(newCIDR)
	if err != nil {
		return validateCIDRBlock(newCIDR)
	}
	for _, sub := range existingSubnets {
		_, existingNet, err := net.ParseCIDR(sub.CidrBlock)
		if err != nil {
			continue
		}
		if newNet.Contains(existingNet.IP) || existingNet.Contains(newNet.IP) {
			return awserrors.NewAWSError("InvalidSubnet.Conflict",
				"The CIDR '"+newCIDR+"' conflicts with another subnet",
				http.StatusBadRequest)
		}
	}
	return nil
}

// parseInt64Param parses a string parameter into an int64. An empty string
// returns -1 (used as a sentinel for "not specified" in port and MaxResults
// contexts). Invalid input yields an InvalidParameterValue error naming the
// offending parameter.
func parseInt64Param(s string, paramName string) (int64, error) {
	if s == "" {
		return -1, nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, awserrors.NewAWSError("InvalidParameterValue",
			fmt.Sprintf("Invalid value '%s' for parameter %s", s, paramName),
			http.StatusBadRequest)
	}
	return n, nil
}

// validateInstanceTenancy checks the InstanceTenancy value against the
// Smithy Tenancy enum (default, dedicated, host).
func validateInstanceTenancy(tenancy string) error {
	switch tenancy {
	case "", "default", "dedicated", "host":
		return nil
	}
	return awserrors.NewAWSError("InvalidParameterValue",
		fmt.Sprintf("Value (%s) for parameter instanceTenancy is invalid. Valid values: default, dedicated, host", tenancy),
		http.StatusBadRequest)
}

// validateGroupReferences checks that all UserIdGroupPair GroupId references
// exist in the store before applying rules.
func validateGroupReferences(store *ec2store.EC2Store, rules []ec2store.IPRule) error {
	for _, rule := range rules {
		for _, pair := range rule.UserIdGroupPairs {
			if pair.GroupId == "" {
				continue
			}
			if _, err := store.GetSecurityGroup(pair.GroupId); err != nil {
				return awserrors.NewAWSError("InvalidGroup.NotFound",
					fmt.Sprintf("The security group '%s' does not exist", pair.GroupId),
					http.StatusNotFound)
			}
		}
	}
	return nil
}

// validateIPRule validates the protocol/port combination for a security group rule.
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
		if from != -1 || to != -1 {
			return awserrors.NewAWSError("InvalidParameterValue",
				fmt.Sprintf("Ports are not allowed for protocol '%s'", proto),
				http.StatusBadRequest)
		}
	}
	return nil
}

// validateIPRuleRanges validates CIDR format for all IP ranges in a rule.
func validateIPRuleRanges(rule ec2store.IPRule) error {
	for _, r := range rule.IpRanges {
		if r.CidrIp != "" {
			if err := validateIPv4CIDRInRule(r.CidrIp); err != nil {
				return err
			}
		}
	}
	for _, r := range rule.Ipv6Ranges {
		if r.CidrIp != "" {
			if err := validateIPv6CIDRInRule(r.CidrIp); err != nil {
				return err
			}
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
