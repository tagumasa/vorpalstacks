package wafv2

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
)

// AdminHandler provides WAFv2 service administration functionality.
// It implements the WAFV2ServiceHandler interface for gRPC-Web communication.
type AdminHandler struct{}

var _ wafv2connect.WAFV2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new WAFv2 AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListWebACLs lists web access control lists in WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	return connect.NewResponse(&pb.ListWebACLsResponse{
		Webacls: []*pb.WebACLSummary{},
	}), nil
}

// AssociateWebACL associates a web ACL with a regional resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateWebACL(ctx context.Context, req *connect.Request[pb.AssociateWebACLRequest]) (*connect.Response[pb.AssociateWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.AssociateWebACL is not implemented"))
}

// CheckCapacity checks the capacity for a rule group or web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CheckCapacity(ctx context.Context, req *connect.Request[pb.CheckCapacityRequest]) (*connect.Response[pb.CheckCapacityResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.CheckCapacity is not implemented"))
}

// CreateAPIKey creates an API key for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAPIKey(ctx context.Context, req *connect.Request[pb.CreateAPIKeyRequest]) (*connect.Response[pb.CreateAPIKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.CreateAPIKey is not implemented"))
}

// CreateIPSet creates an IP set for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateIPSet(ctx context.Context, req *connect.Request[pb.CreateIPSetRequest]) (*connect.Response[pb.CreateIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.CreateIPSet is not implemented"))
}

// CreateRegexPatternSet creates a regex pattern set for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRegexPatternSet(ctx context.Context, req *connect.Request[pb.CreateRegexPatternSetRequest]) (*connect.Response[pb.CreateRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.CreateRegexPatternSet is not implemented"))
}

// CreateRuleGroup creates a rule group for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRuleGroup(ctx context.Context, req *connect.Request[pb.CreateRuleGroupRequest]) (*connect.Response[pb.CreateRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.CreateRuleGroup is not implemented"))
}

// CreateWebACL creates a web access control list for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateWebACL(ctx context.Context, req *connect.Request[pb.CreateWebACLRequest]) (*connect.Response[pb.CreateWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.CreateWebACL is not implemented"))
}

// DeleteAPIKey deletes an API key for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAPIKey(ctx context.Context, req *connect.Request[pb.DeleteAPIKeyRequest]) (*connect.Response[pb.DeleteAPIKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteAPIKey is not implemented"))
}

// DeleteFirewallManagerRuleGroups deletes Firewall Manager rule groups from a web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteFirewallManagerRuleGroups(ctx context.Context, req *connect.Request[pb.DeleteFirewallManagerRuleGroupsRequest]) (*connect.Response[pb.DeleteFirewallManagerRuleGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteFirewallManagerRuleGroups is not implemented"))
}

// DeleteIPSet deletes an IP set from WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteIPSet(ctx context.Context, req *connect.Request[pb.DeleteIPSetRequest]) (*connect.Response[pb.DeleteIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteIPSet is not implemented"))
}

// DeleteLoggingConfiguration deletes a logging configuration from WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteLoggingConfiguration(ctx context.Context, req *connect.Request[pb.DeleteLoggingConfigurationRequest]) (*connect.Response[pb.DeleteLoggingConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteLoggingConfiguration is not implemented"))
}

// DeletePermissionPolicy deletes a permission policy from a rule group or web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePermissionPolicy(ctx context.Context, req *connect.Request[pb.DeletePermissionPolicyRequest]) (*connect.Response[pb.DeletePermissionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeletePermissionPolicy is not implemented"))
}

// DeleteRegexPatternSet deletes a regex pattern set from WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRegexPatternSet(ctx context.Context, req *connect.Request[pb.DeleteRegexPatternSetRequest]) (*connect.Response[pb.DeleteRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteRegexPatternSet is not implemented"))
}

// DeleteRuleGroup deletes a rule group from WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRuleGroup(ctx context.Context, req *connect.Request[pb.DeleteRuleGroupRequest]) (*connect.Response[pb.DeleteRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteRuleGroup is not implemented"))
}

// DeleteWebACL deletes a web access control list from WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteWebACL(ctx context.Context, req *connect.Request[pb.DeleteWebACLRequest]) (*connect.Response[pb.DeleteWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DeleteWebACL is not implemented"))
}

// DescribeAllManagedProducts describes all managed products for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeAllManagedProducts(ctx context.Context, req *connect.Request[pb.DescribeAllManagedProductsRequest]) (*connect.Response[pb.DescribeAllManagedProductsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DescribeAllManagedProducts is not implemented"))
}

// DescribeManagedProductsByVendor describes managed products for a specific vendor.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeManagedProductsByVendor(ctx context.Context, req *connect.Request[pb.DescribeManagedProductsByVendorRequest]) (*connect.Response[pb.DescribeManagedProductsByVendorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DescribeManagedProductsByVendor is not implemented"))
}

// DescribeManagedRuleGroup describes a managed rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeManagedRuleGroup(ctx context.Context, req *connect.Request[pb.DescribeManagedRuleGroupRequest]) (*connect.Response[pb.DescribeManagedRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DescribeManagedRuleGroup is not implemented"))
}

// DisassociateWebACL disassociates a web ACL from a regional resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisassociateWebACL(ctx context.Context, req *connect.Request[pb.DisassociateWebACLRequest]) (*connect.Response[pb.DisassociateWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.DisassociateWebACL is not implemented"))
}

// GenerateMobileSdkReleaseUrl generates a URL for downloading a mobile SDK release.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateMobileSdkReleaseUrl(ctx context.Context, req *connect.Request[pb.GenerateMobileSdkReleaseUrlRequest]) (*connect.Response[pb.GenerateMobileSdkReleaseUrlResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GenerateMobileSdkReleaseUrl is not implemented"))
}

// GetDecryptedAPIKey returns the decrypted API key for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDecryptedAPIKey(ctx context.Context, req *connect.Request[pb.GetDecryptedAPIKeyRequest]) (*connect.Response[pb.GetDecryptedAPIKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetDecryptedAPIKey is not implemented"))
}

// GetIPSet returns details about an IP set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetIPSet(ctx context.Context, req *connect.Request[pb.GetIPSetRequest]) (*connect.Response[pb.GetIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetIPSet is not implemented"))
}

// GetLoggingConfiguration returns details about a logging configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetLoggingConfiguration(ctx context.Context, req *connect.Request[pb.GetLoggingConfigurationRequest]) (*connect.Response[pb.GetLoggingConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetLoggingConfiguration is not implemented"))
}

// GetManagedRuleSet returns details about a managed rule set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetManagedRuleSet(ctx context.Context, req *connect.Request[pb.GetManagedRuleSetRequest]) (*connect.Response[pb.GetManagedRuleSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetManagedRuleSet is not implemented"))
}

// GetMobileSdkRelease returns details about a mobile SDK release.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMobileSdkRelease(ctx context.Context, req *connect.Request[pb.GetMobileSdkReleaseRequest]) (*connect.Response[pb.GetMobileSdkReleaseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetMobileSdkRelease is not implemented"))
}

// GetPermissionPolicy returns the permission policy for a rule group or web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPermissionPolicy(ctx context.Context, req *connect.Request[pb.GetPermissionPolicyRequest]) (*connect.Response[pb.GetPermissionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetPermissionPolicy is not implemented"))
}

// GetRateBasedStatementManagedKeys returns the managed keys for a rate-based rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRateBasedStatementManagedKeys(ctx context.Context, req *connect.Request[pb.GetRateBasedStatementManagedKeysRequest]) (*connect.Response[pb.GetRateBasedStatementManagedKeysResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetRateBasedStatementManagedKeys is not implemented"))
}

// GetRegexPatternSet returns details about a regex pattern set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRegexPatternSet(ctx context.Context, req *connect.Request[pb.GetRegexPatternSetRequest]) (*connect.Response[pb.GetRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetRegexPatternSet is not implemented"))
}

// GetRuleGroup returns details about a rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRuleGroup(ctx context.Context, req *connect.Request[pb.GetRuleGroupRequest]) (*connect.Response[pb.GetRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetRuleGroup is not implemented"))
}

// GetSampledRequests returns a sample of web requests for a rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSampledRequests(ctx context.Context, req *connect.Request[pb.GetSampledRequestsRequest]) (*connect.Response[pb.GetSampledRequestsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetSampledRequests is not implemented"))
}

// GetTopPathStatisticsByTraffic returns top path statistics by traffic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTopPathStatisticsByTraffic(ctx context.Context, req *connect.Request[pb.GetTopPathStatisticsByTrafficRequest]) (*connect.Response[pb.GetTopPathStatisticsByTrafficResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetTopPathStatisticsByTraffic is not implemented"))
}

// GetWebACL returns details about a web access control list.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetWebACL(ctx context.Context, req *connect.Request[pb.GetWebACLRequest]) (*connect.Response[pb.GetWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetWebACL is not implemented"))
}

// GetWebACLForResource returns the web ACL associated with a resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetWebACLForResource(ctx context.Context, req *connect.Request[pb.GetWebACLForResourceRequest]) (*connect.Response[pb.GetWebACLForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.GetWebACLForResource is not implemented"))
}

// ListAPIKeys lists API keys for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAPIKeys(ctx context.Context, req *connect.Request[pb.ListAPIKeysRequest]) (*connect.Response[pb.ListAPIKeysResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListAPIKeys is not implemented"))
}

// ListAvailableManagedRuleGroups lists available managed rule groups.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAvailableManagedRuleGroups(ctx context.Context, req *connect.Request[pb.ListAvailableManagedRuleGroupsRequest]) (*connect.Response[pb.ListAvailableManagedRuleGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListAvailableManagedRuleGroups is not implemented"))
}

// ListAvailableManagedRuleGroupVersions lists versions of a managed rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAvailableManagedRuleGroupVersions(ctx context.Context, req *connect.Request[pb.ListAvailableManagedRuleGroupVersionsRequest]) (*connect.Response[pb.ListAvailableManagedRuleGroupVersionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListAvailableManagedRuleGroupVersions is not implemented"))
}

// ListIPSets lists IP sets in WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListIPSets(ctx context.Context, req *connect.Request[pb.ListIPSetsRequest]) (*connect.Response[pb.ListIPSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListIPSets is not implemented"))
}

// ListLoggingConfigurations lists logging configurations in WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListLoggingConfigurations(ctx context.Context, req *connect.Request[pb.ListLoggingConfigurationsRequest]) (*connect.Response[pb.ListLoggingConfigurationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListLoggingConfigurations is not implemented"))
}

// ListManagedRuleSets lists managed rule sets in WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListManagedRuleSets(ctx context.Context, req *connect.Request[pb.ListManagedRuleSetsRequest]) (*connect.Response[pb.ListManagedRuleSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListManagedRuleSets is not implemented"))
}

// ListMobileSdkReleases lists mobile SDK releases for WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListMobileSdkReleases(ctx context.Context, req *connect.Request[pb.ListMobileSdkReleasesRequest]) (*connect.Response[pb.ListMobileSdkReleasesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListMobileSdkReleases is not implemented"))
}

// ListRegexPatternSets lists regex pattern sets in WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRegexPatternSets(ctx context.Context, req *connect.Request[pb.ListRegexPatternSetsRequest]) (*connect.Response[pb.ListRegexPatternSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListRegexPatternSets is not implemented"))
}

// ListResourcesForWebACL lists resources associated with a web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListResourcesForWebACL(ctx context.Context, req *connect.Request[pb.ListResourcesForWebACLRequest]) (*connect.Response[pb.ListResourcesForWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListResourcesForWebACL is not implemented"))
}

// ListRuleGroups lists rule groups in WAFv2.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRuleGroups(ctx context.Context, req *connect.Request[pb.ListRuleGroupsRequest]) (*connect.Response[pb.ListRuleGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListRuleGroups is not implemented"))
}

// ListTagsForResource lists tags for a WAFv2 resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.ListTagsForResource is not implemented"))
}

// PutLoggingConfiguration creates or updates a logging configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutLoggingConfiguration(ctx context.Context, req *connect.Request[pb.PutLoggingConfigurationRequest]) (*connect.Response[pb.PutLoggingConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.PutLoggingConfiguration is not implemented"))
}

// PutManagedRuleSetVersions updates the versions of a managed rule set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutManagedRuleSetVersions(ctx context.Context, req *connect.Request[pb.PutManagedRuleSetVersionsRequest]) (*connect.Response[pb.PutManagedRuleSetVersionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.PutManagedRuleSetVersions is not implemented"))
}

// PutPermissionPolicy creates or updates a permission policy for a rule group or web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutPermissionPolicy(ctx context.Context, req *connect.Request[pb.PutPermissionPolicyRequest]) (*connect.Response[pb.PutPermissionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.PutPermissionPolicy is not implemented"))
}

// TagResource adds tags to a WAFv2 resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pb.TagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.TagResource is not implemented"))
}

// UntagResource removes tags from a WAFv2 resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pb.UntagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.UntagResource is not implemented"))
}

// UpdateIPSet updates an IP set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateIPSet(ctx context.Context, req *connect.Request[pb.UpdateIPSetRequest]) (*connect.Response[pb.UpdateIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.UpdateIPSet is not implemented"))
}

// UpdateManagedRuleSetVersionExpiryDate updates the expiry date for a managed rule set version.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateManagedRuleSetVersionExpiryDate(ctx context.Context, req *connect.Request[pb.UpdateManagedRuleSetVersionExpiryDateRequest]) (*connect.Response[pb.UpdateManagedRuleSetVersionExpiryDateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.UpdateManagedRuleSetVersionExpiryDate is not implemented"))
}

// UpdateRegexPatternSet updates a regex pattern set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRegexPatternSet(ctx context.Context, req *connect.Request[pb.UpdateRegexPatternSetRequest]) (*connect.Response[pb.UpdateRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.UpdateRegexPatternSet is not implemented"))
}

// UpdateRuleGroup updates a rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRuleGroup(ctx context.Context, req *connect.Request[pb.UpdateRuleGroupRequest]) (*connect.Response[pb.UpdateRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.UpdateRuleGroup is not implemented"))
}

// UpdateWebACL updates a web access control list.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateWebACL(ctx context.Context, req *connect.Request[pb.UpdateWebACLRequest]) (*connect.Response[pb.UpdateWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wafv2.WAFV2Service.UpdateWebACL is not implemented"))
}
