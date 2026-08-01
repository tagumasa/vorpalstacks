package wafv2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

// ManagedRuleGroupSummary represents a high-level summary of an AWS-managed rule group.
type ManagedRuleGroupSummary struct {
	Name                *string `json:"Name,omitempty"`
	VendorName          *string `json:"VendorName,omitempty"`
	Description         *string `json:"Description,omitempty"`
	VersioningSupported *bool   `json:"VersioningSupported,omitempty"`
}

// awsManagedRuleGroups contains the list of available AWS-managed rule groups.
var awsManagedRuleGroups = []ManagedRuleGroupSummary{
	{Name: proto.String("AWSManagedRulesCommonRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Core rule set containing rules generally applicable to web applications."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesAmazonIpReputationList"), VendorName: proto.String("AWS"), Description: proto.String("Rules based on Amazon threat intelligence."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesKnownBadInputsRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block request patterns that are known to be invalid."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesSQLiRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block SQL injection attacks."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesLinuxRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block exploitation of Linux specific vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesUnixRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block exploitation of Unix specific vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesWindowsRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block exploitation of Windows specific vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesPHPRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block exploitation of PHP specific vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesWordPressRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block exploitation of WordPress specific vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesBotControlRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to detect and block bot traffic."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesATPRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Account Takeover Protection rules."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesACFPRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Account Creation Fraud Prevention rules."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesAnonymousIpList"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block requests from anonymous IP addresses."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesJavaRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Rules to block exploitation of Java specific vulnerabilities."), VersioningSupported: proto.Bool(true)},
	{Name: proto.String("AWSManagedRulesJTRRuleSet"), VendorName: proto.String("AWS"), Description: proto.String("Jack the Ripper rules for detecting credential stuffing patterns."), VersioningSupported: proto.Bool(true)},
}

// ListAvailableManagedRuleGroups returns a paginated list of all available AWS-managed rule groups.
func (s *WAFv2Service) ListAvailableManagedRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	_ = scope

	limit := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	startIdx := 0
	if nextMarker != "" {
		for i, rg := range awsManagedRuleGroups {
			if rg.Name != nil && *rg.Name == nextMarker {
				startIdx = i
				break
			}
		}
	}

	endIdx := startIdx + limit
	var respNextMarker *string
	if endIdx < len(awsManagedRuleGroups) {
		respNextMarker = awsManagedRuleGroups[endIdx].Name
	} else {
		endIdx = len(awsManagedRuleGroups)
	}

	resp := map[string]interface{}{
		"ManagedRuleGroups": awsManagedRuleGroups[startIdx:endIdx],
	}
	var nextMarkerVal string
	if respNextMarker != nil {
		nextMarkerVal = *respNextMarker
	}
	pagination.SetNextToken(resp, "NextMarker", nextMarkerVal)
	return resp, nil
}

// DescribeManagedRuleGroup provides details about the specified AWS-managed rule group, including capacity and available labels.
func (s *WAFv2Service) DescribeManagedRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")
	versionName := request.GetStringParam(req.Parameters, "VersionName")

	if name == "" || vendorName == "" {
		return nil, invalidParamError("Name and VendorName are required")
	}

	for _, rg := range awsManagedRuleGroups {
		if rg.Name != nil && rg.VendorName != nil && *rg.Name == name && *rg.VendorName == vendorName {
			capacity := int64(100)
			if name == "AWSManagedRulesBotControlRuleSet" {
				capacity = 50
			} else if name == "AWSManagedRulesATPRuleSet" {
				capacity = 50
			}

			labelNamespace := "awswaf:" + vendorName + ":" + name + ":"

			var versionNameResp *string
			if versionName != "" {
				versionNameResp = &versionName
			}

			return map[string]interface{}{
				"Capacity":       capacity,
				"LabelNamespace": labelNamespace,
				"VersionName":    versionNameResp,
				"SnsTopicArn":    fmt.Sprintf("arn:aws:sns:us-east-1:123456789012:aws-managed-waf-%s", strings.ToLower(name)),
				"Rules": []map[string]interface{}{
					{"Name": name + "_Rule1", "Action": map[string]interface{}{"Block": map[string]interface{}{}}},
					{"Name": name + "_Rule2", "Action": map[string]interface{}{"Count": map[string]interface{}{}}},
				},
				"AvailableLabels": []map[string]interface{}{
					{"Name": labelNamespace + "rule1"},
					{"Name": labelNamespace + "rule2"},
				},
				"ConsumedLabels": []map[string]interface{}{
					{"Name": labelNamespace + "rule1"},
				},
			}, nil
		}
	}

	return nil, notFoundError("ManagedRuleGroup")
}

// ListAvailableManagedRuleGroupVersions returns the available versions for the specified AWS-managed rule group.
func (s *WAFv2Service) ListAvailableManagedRuleGroupVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	vendorName := request.GetStringParam(req.Parameters, "VendorName")

	if name == "" || vendorName == "" {
		return nil, invalidParamError("Name and VendorName are required")
	}

	found := false
	for _, rg := range awsManagedRuleGroups {
		if rg.Name != nil && rg.VendorName != nil && *rg.Name == name && *rg.VendorName == vendorName {
			found = true
			break
		}
	}

	if !found {
		return nil, notFoundError("ManagedRuleGroup")
	}

	now := time.Now()
	versions := make([]map[string]interface{}, 0, 3)
	for i := 0; i < 3; i++ {
		t := now.AddDate(0, -i, 0)
		versions = append(versions, map[string]interface{}{
			"Name":        fmt.Sprintf("Version_%04d_%02d_01", t.Year(), t.Month()),
			"LastUpdated": t.Unix(),
		})
	}

	return map[string]interface{}{
		"Versions": versions,
	}, nil
}
