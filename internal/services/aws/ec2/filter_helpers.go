package ec2

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

// ec2Filter represents a parsed EC2 Filter.N.Name + Filter.N.Value.M pair.
type ec2Filter struct {
	Name   string
	Values []string
}

// parseFilters extracts EC2 Filter.N.Name/Filter.N.Value.M from request params.
// Returns an error when a filter has a Name but no Values (AWS InvalidFilter).
func parseFilters(params map[string]interface{}) ([]ec2Filter, error) {
	var filters []ec2Filter
	for i := 1; ; i++ {
		name := request.GetStringParam(params, fmt.Sprintf("Filter.%d.Name", i))
		if name == "" {
			break
		}
		var values []string
		for j := 1; ; j++ {
			v := request.GetStringParam(params, fmt.Sprintf("Filter.%d.Value.%d", i, j))
			if v == "" {
				break
			}
			values = append(values, v)
		}
		if len(values) == 0 {
			return nil, awserrors.NewAWSError("InvalidFilter",
				fmt.Sprintf("Filter %d ('%s') requires at least one value.", i, name),
				http.StatusBadRequest)
		}
		filters = append(filters, ec2Filter{Name: name, Values: values})
	}
	return filters, nil
}

// anyMatch returns true if target matches any of the values (case-insensitive).
func anyMatch(values []string, target string) bool {
	for _, v := range values {
		if strings.EqualFold(v, target) {
			return true
		}
	}
	return false
}

// anyMatchBool returns true if the target boolean matches any value
// (AWS boolean filters use the strings "true" and "false").
func anyMatchBool(values []string, target bool) bool {
	return anyMatch(values, strconv.FormatBool(target))
}

// anyMatchInt64 returns true if the target integer equals any numeric value.
func anyMatchInt64(values []string, target int64) bool {
	for _, v := range values {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n == target {
			return true
		}
	}
	return false
}

// anyMatchInt64Port matches port-number filter values; -1 (all ports/ICMP
// types) is also accepted as a comparable value.
func anyMatchInt64Port(values []string, target int64) bool {
	return anyMatchInt64(values, target)
}

// validateFilterNames rejects filter names that are not recognised for the
// given resource kind. AWS returns InvalidParameterValue for unknown filter
// names rather than ignoring them.
func validateFilterNames(filters []ec2Filter, allowed func(name string) bool) error {
	for _, f := range filters {
		if !allowed(f.Name) {
			return awserrors.NewAWSError("InvalidParameterValue",
				fmt.Sprintf("Value (%s) for parameter Filter is invalid. The filter '%s' is not recognised", f.Name, f.Name),
				http.StatusBadRequest)
		}
	}
	return nil
}

// allowedVPCFilters lists the AWS-documented DescribeVpcs filter names.
// The ipv6-cidr-block-association.* family is accepted even though this
// implementation stores no VPC IPv6 associations; those filters match
// nothing, mirroring an account without associated IPv6 blocks.
func allowedVPCFilters(name string) bool {
	switch name {
	case "cidr", "cidr-block-association.cidr-block", "cidr-block-association.association-id",
		"cidr-block-association.state", "dhcp-options-id", "owner-id", "state", "vpc-id",
		"is-default", "tag-key", "tag-value",
		"ipv6-cidr-block-association.ipv6-cidr-block", "ipv6-cidr-block-association.ipv6-pool",
		"ipv6-cidr-block-association.association-id", "ipv6-cidr-block-association.state":
		return true
	}
	return strings.HasPrefix(name, "tag:")
}

// allowedSubnetFilters lists the AWS-documented DescribeSubnets filter names
// supported by this implementation, including the documented alternative
// spellings (availabilityZone, availabilityZoneId, cidr, cidrBlock,
// defaultForAz). Filters for attributes this implementation does not carry
// (Outposts, customer-owned IPv4 pools, LNI device index, IPv6 CIDR
// associations, private DNS name options) are accepted and match nothing or
// the attribute's default value, so that documented filters are not rejected.
func allowedSubnetFilters(name string) bool {
	switch name {
	case "availability-zone", "availabilityZone", "availability-zone-id", "availabilityZoneId",
		"available-ip-address-count", "cidr-block", "cidr", "cidrBlock", "default-for-az",
		"defaultForAz", "enable-dns64", "enable-lni-at-device-index", "ipv6-native",
		"map-customer-owned-ip-on-launch", "map-public-ip-on-launch", "owner-id", "state",
		"subnet-id", "subnet-arn", "vpc-id",
		"customer-owned-ipv4-pool", "outpost-arn",
		"ipv6-cidr-block-association.ipv6-cidr-block",
		"ipv6-cidr-block-association.association-id",
		"ipv6-cidr-block-association.state",
		"private-dns-name-options-on-launch.hostname-type",
		"private-dns-name-options-on-launch.enable-resource-name-dns-a-record",
		"private-dns-name-options-on-launch.enable-resource-name-dns-aaaa-record",
		"tag-key", "tag-value":
		return true
	}
	return strings.HasPrefix(name, "tag:")
}

// allowedSGFilters lists the AWS-documented DescribeSecurityGroups filter
// names supported by this implementation: the plain ip-permission.* family
// (inbound rules, the form used by the AWS documentation examples), the
// ingress.ip-permission.* alias, the egress.ip-permission.* family, and the
// scalar filters.
func allowedSGFilters(name string) bool {
	switch name {
	case "description", "group-id", "group-name", "owner-id", "vpc-id",
		"tag-key", "tag-value":
		return true
	}
	if strings.HasPrefix(name, "tag:") {
		return true
	}
	for _, prefix := range []string{"ip-permission.", "ingress.ip-permission.", "egress.ip-permission."} {
		if strings.HasPrefix(name, prefix) {
			switch strings.TrimPrefix(name, prefix) {
			case "cidr", "from-port", "group-id", "group-name", "ipv6-cidr",
				"prefix-list-id", "protocol", "to-port", "user-id":
				return true
			}
		}
	}
	return false
}

// hasTagKey returns true if any tag has a key matching any of the given values.
func hasTagKey(tags []types.Tag, values []string) bool {
	for _, t := range tags {
		if anyMatch(values, t.Key) {
			return true
		}
	}
	return false
}

// hasTagValue returns true if any tag has a value matching any of the given values.
func hasTagValue(tags []types.Tag, values []string) bool {
	for _, t := range tags {
		if anyMatch(values, t.Value) {
			return true
		}
	}
	return false
}

// hasTagKeyValues returns true if any tag has the given key and a value
// matching any of the values (used by the "tag:<key>" filter name form).
func hasTagKeyValues(tags []types.Tag, key string, values []string) bool {
	for _, t := range tags {
		if t.Key != key {
			continue
		}
		for _, v := range values {
			if t.Value == v {
				return true
			}
		}
	}
	return false
}

// hasTagKeyValue returns true if any tag matches any "key:value" entry in values.
func hasTagKeyValue(tags []types.Tag, values []string) bool {
	for _, v := range values {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) != 2 {
			continue
		}
		for _, t := range tags {
			if t.Key == parts[0] && t.Value == parts[1] {
				return true
			}
		}
	}
	return false
}
