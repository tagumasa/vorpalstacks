package ec2

import (
	"fmt"

	ec2store "vorpalstacks/internal/store/aws/ec2"
	"vorpalstacks/internal/utils/aws/generators"
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

// ParseTags extracts EC2 tags from request parameters.
// Supports both flat Tag.N.Key/Tag.N.Value and TagSpecifications formats.
func ParseTags(params map[string]interface{}) []ec2store.Tag {
	tags := make([]ec2store.Tag, 0)
	for i := 1; ; i++ {
		key := getStringParam(params, fmt.Sprintf("Tag.%d.Key", i))
		if key == "" {
			break
		}
		value := getStringParam(params, fmt.Sprintf("Tag.%d.Value", i))
		tags = append(tags, ec2store.Tag{Key: key, Value: value})
	}
	if len(tags) > 0 {
		return tags
	}
	if tagList, ok := params["Tag"].([]interface{}); ok {
		for _, t := range tagList {
			if m, ok := t.(map[string]interface{}); ok {
				k, _ := m["Key"].(string)
				v, _ := m["Value"].(string)
				if k != "" {
					tags = append(tags, ec2store.Tag{Key: k, Value: v})
				}
			}
		}
	}
	return tags
}

// getStringParam retrieves a string parameter, trying lowercase key first then title case.
func getStringParam(params map[string]interface{}, key string) string {
	for _, k := range []string{key, lowerFirst(key)} {
		if v, ok := params[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// lowerFirst converts the first character of a string to lowercase.
func lowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string([]byte{s[0] + 32}) + s[1:]
}

// mergeIPRules appends new rules to an existing rule set, avoiding duplicates
// based on protocol, from/to port, and CIDR.
func mergeIPRules(existing []ec2store.IPRule, newRules ...ec2store.IPRule) []ec2store.IPRule {
	for _, nr := range newRules {
		duplicate := false
		for _, er := range existing {
			if er.IpProtocol == nr.IpProtocol && er.FromPort == nr.FromPort && er.ToPort == nr.ToPort {
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
			if er.IpProtocol == nr.IpProtocol && er.FromPort == nr.FromPort && er.ToPort == nr.ToPort {
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
