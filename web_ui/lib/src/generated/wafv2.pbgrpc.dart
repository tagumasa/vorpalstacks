// This is a generated file - do not edit.
//
// Generated from wafv2.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'wafv2.pb.dart' as $0;

export 'wafv2.pb.dart';

/// WAFV2Service provides wafv2 API operations.
@$pb.GrpcServiceName('wafv2.WAFV2Service')
class WAFV2ServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  WAFV2ServiceClient(super.channel, {super.options, super.interceptors});

  /// Associates a web ACL with a resource, to protect the resource. Use this for all resource types except for Amazon CloudFront distributions. For Amazon CloudFront, call UpdateDistribution for the dis...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AssociateWebACLResponse> associateWebACL(
    $0.AssociateWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateWebACL, request, options: options);
  }

  /// Returns the web ACL capacity unit (WCU) requirements for a specified scope and set of rules. You can use this to check the capacity requirements for the rules you want to use in a RuleGroup or WebA...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CheckCapacityResponse> checkCapacity(
    $0.CheckCapacityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkCapacity, request, options: options);
  }

  /// Creates an API key that contains a set of token domains. API keys are required for the integration of the CAPTCHA API in your JavaScript client applications. The API lets you customize the placemen...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateAPIKeyResponse> createAPIKey(
    $0.CreateAPIKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAPIKey, request, options: options);
  }

  /// Creates an IPSet, which you use to identify web requests that originate from specific IP addresses or ranges of IP addresses. For example, if you're receiving a lot of requests from a ranges of IP ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateIPSetResponse> createIPSet(
    $0.CreateIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createIPSet, request, options: options);
  }

  /// Creates a RegexPatternSet, which you reference in a RegexPatternSetReferenceStatement, to have WAF inspect a web request component for the specified patterns.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRegexPatternSetResponse> createRegexPatternSet(
    $0.CreateRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRegexPatternSet, request, options: options);
  }

  /// Creates a RuleGroup per the specifications provided. A rule group defines a collection of rules to inspect and control web requests that you can use in a WebACL. When you create a rule group, you d...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRuleGroupResponse> createRuleGroup(
    $0.CreateRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRuleGroup, request, options: options);
  }

  /// Creates a WebACL per the specifications provided. A web ACL defines a collection of rules to use to inspect and control web requests. Each rule has a statement that defines what to look for in web ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateWebACLResponse> createWebACL(
    $0.CreateWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createWebACL, request, options: options);
  }

  /// Deletes the specified API key. After you delete a key, it can take up to 24 hours for WAF to disallow use of the key in all regions.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteAPIKeyResponse> deleteAPIKey(
    $0.DeleteAPIKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAPIKey, request, options: options);
  }

  /// Deletes all rule groups that are managed by Firewall Manager from the specified WebACL. You can only use this if ManagedByFirewallManager and RetrofittedByFirewallManager are both false in the web ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteFirewallManagerRuleGroupsResponse>
      deleteFirewallManagerRuleGroups(
    $0.DeleteFirewallManagerRuleGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFirewallManagerRuleGroups, request,
        options: options);
  }

  /// Deletes the specified IPSet.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteIPSetResponse> deleteIPSet(
    $0.DeleteIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIPSet, request, options: options);
  }

  /// Deletes the LoggingConfiguration from the specified web ACL.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteLoggingConfigurationResponse>
      deleteLoggingConfiguration(
    $0.DeleteLoggingConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteLoggingConfiguration, request,
        options: options);
  }

  /// Permanently deletes an IAM policy from the specified rule group. You must be the owner of the rule group to perform this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeletePermissionPolicyResponse>
      deletePermissionPolicy(
    $0.DeletePermissionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePermissionPolicy, request,
        options: options);
  }

  /// Deletes the specified RegexPatternSet.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRegexPatternSetResponse> deleteRegexPatternSet(
    $0.DeleteRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRegexPatternSet, request, options: options);
  }

  /// Deletes the specified RuleGroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRuleGroupResponse> deleteRuleGroup(
    $0.DeleteRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRuleGroup, request, options: options);
  }

  /// Deletes the specified WebACL. You can only use this if ManagedByFirewallManager is false in the web ACL. Before deleting any web ACL, first disassociate it from all resources. To retrieve a list of...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteWebACLResponse> deleteWebACL(
    $0.DeleteWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteWebACL, request, options: options);
  }

  /// Provides high-level information for the Amazon Web Services Managed Rules rule groups and Amazon Web Services Marketplace managed rule groups.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAllManagedProductsResponse>
      describeAllManagedProducts(
    $0.DescribeAllManagedProductsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAllManagedProducts, request,
        options: options);
  }

  /// Provides high-level information for the managed rule groups owned by a specific vendor.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeManagedProductsByVendorResponse>
      describeManagedProductsByVendor(
    $0.DescribeManagedProductsByVendorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeManagedProductsByVendor, request,
        options: options);
  }

  /// Provides high-level information for a managed rule group, including descriptions of the rules.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeManagedRuleGroupResponse>
      describeManagedRuleGroup(
    $0.DescribeManagedRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeManagedRuleGroup, request,
        options: options);
  }

  /// Disassociates the specified resource from its web ACL association, if it has one. Use this for all resource types except for Amazon CloudFront distributions. For Amazon CloudFront, call UpdateDistr...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DisassociateWebACLResponse> disassociateWebACL(
    $0.DisassociateWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateWebACL, request, options: options);
  }

  /// Generates a presigned download URL for the specified release of the mobile SDK. The mobile SDK is not generally available. Customers who have access to the mobile SDK can use it to establish and ma...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateMobileSdkReleaseUrlResponse>
      generateMobileSdkReleaseUrl(
    $0.GenerateMobileSdkReleaseUrlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateMobileSdkReleaseUrl, request,
        options: options);
  }

  /// Returns your API key in decrypted form. Use this to check the token domains that you have defined for the key. API keys are required for the integration of the CAPTCHA API in your JavaScript client...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDecryptedAPIKeyResponse> getDecryptedAPIKey(
    $0.GetDecryptedAPIKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDecryptedAPIKey, request, options: options);
  }

  /// Retrieves the specified IPSet.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetIPSetResponse> getIPSet(
    $0.GetIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIPSet, request, options: options);
  }

  /// Returns the LoggingConfiguration for the specified web ACL.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLoggingConfigurationResponse>
      getLoggingConfiguration(
    $0.GetLoggingConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLoggingConfiguration, request,
        options: options);
  }

  /// Retrieves the specified managed rule set. This is intended for use only by vendors of managed rule sets. Vendors are Amazon Web Services and Amazon Web Services Marketplace sellers. Vendors, you ca...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetManagedRuleSetResponse> getManagedRuleSet(
    $0.GetManagedRuleSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getManagedRuleSet, request, options: options);
  }

  /// Retrieves information for the specified mobile SDK release, including release notes and tags. The mobile SDK is not generally available. Customers who have access to the mobile SDK can use it to es...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetMobileSdkReleaseResponse> getMobileSdkRelease(
    $0.GetMobileSdkReleaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMobileSdkRelease, request, options: options);
  }

  /// Returns the IAM policy that is attached to the specified rule group. You must be the owner of the rule group to perform this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPermissionPolicyResponse> getPermissionPolicy(
    $0.GetPermissionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPermissionPolicy, request, options: options);
  }

  /// Retrieves the IP addresses that are currently blocked by a rate-based rule instance. This is only available for rate-based rules that aggregate solely on the IP address or on the forwarded IP addre...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRateBasedStatementManagedKeysResponse>
      getRateBasedStatementManagedKeys(
    $0.GetRateBasedStatementManagedKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRateBasedStatementManagedKeys, request,
        options: options);
  }

  /// Retrieves the specified RegexPatternSet.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRegexPatternSetResponse> getRegexPatternSet(
    $0.GetRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRegexPatternSet, request, options: options);
  }

  /// Retrieves the specified RuleGroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRuleGroupResponse> getRuleGroup(
    $0.GetRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRuleGroup, request, options: options);
  }

  /// Gets detailed information about a specified number of requests--a sample--that WAF randomly selects from among the first 5,000 requests that your Amazon Web Services resource received during a time...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSampledRequestsResponse> getSampledRequests(
    $0.GetSampledRequestsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSampledRequests, request, options: options);
  }

  /// Retrieves aggregated statistics about the top URI paths accessed by bot traffic for a specified web ACL and time window. You can use this operation to analyze which paths on your web application re...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetTopPathStatisticsByTrafficResponse>
      getTopPathStatisticsByTraffic(
    $0.GetTopPathStatisticsByTrafficRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTopPathStatisticsByTraffic, request,
        options: options);
  }

  /// Retrieves the specified WebACL.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetWebACLResponse> getWebACL(
    $0.GetWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getWebACL, request, options: options);
  }

  /// Retrieves the WebACL for the specified resource. This call uses GetWebACL, to verify that your account has permission to access the retrieved web ACL. If you get an error that indicates that your a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetWebACLForResourceResponse> getWebACLForResource(
    $0.GetWebACLForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getWebACLForResource, request, options: options);
  }

  /// Retrieves a list of the API keys that you've defined for the specified scope. API keys are required for the integration of the CAPTCHA API in your JavaScript client applications. The API lets you c...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAPIKeysResponse> listAPIKeys(
    $0.ListAPIKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAPIKeys, request, options: options);
  }

  /// Retrieves an array of managed rule groups that are available for you to use. This list includes all Amazon Web Services Managed Rules rule groups and all of the Amazon Web Services Marketplace mana...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAvailableManagedRuleGroupsResponse>
      listAvailableManagedRuleGroups(
    $0.ListAvailableManagedRuleGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAvailableManagedRuleGroups, request,
        options: options);
  }

  /// Returns a list of the available versions for the specified managed rule group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAvailableManagedRuleGroupVersionsResponse>
      listAvailableManagedRuleGroupVersions(
    $0.ListAvailableManagedRuleGroupVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAvailableManagedRuleGroupVersions, request,
        options: options);
  }

  /// Retrieves an array of IPSetSummary objects for the IP sets that you manage.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListIPSetsResponse> listIPSets(
    $0.ListIPSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listIPSets, request, options: options);
  }

  /// Retrieves an array of your LoggingConfiguration objects.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListLoggingConfigurationsResponse>
      listLoggingConfigurations(
    $0.ListLoggingConfigurationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listLoggingConfigurations, request,
        options: options);
  }

  /// Retrieves the managed rule sets that you own. This is intended for use only by vendors of managed rule sets. Vendors are Amazon Web Services and Amazon Web Services Marketplace sellers. Vendors, yo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListManagedRuleSetsResponse> listManagedRuleSets(
    $0.ListManagedRuleSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listManagedRuleSets, request, options: options);
  }

  /// Retrieves a list of the available releases for the mobile SDK and the specified device platform. The mobile SDK is not generally available. Customers who have access to the mobile SDK can use it to...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListMobileSdkReleasesResponse> listMobileSdkReleases(
    $0.ListMobileSdkReleasesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMobileSdkReleases, request, options: options);
  }

  /// Retrieves an array of RegexPatternSetSummary objects for the regex pattern sets that you manage.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRegexPatternSetsResponse> listRegexPatternSets(
    $0.ListRegexPatternSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRegexPatternSets, request, options: options);
  }

  /// Retrieves an array of the Amazon Resource Names (ARNs) for the resources that are associated with the specified web ACL. For Amazon CloudFront, don't use this call. Instead, use the CloudFront call...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListResourcesForWebACLResponse>
      listResourcesForWebACL(
    $0.ListResourcesForWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourcesForWebACL, request,
        options: options);
  }

  /// Retrieves an array of RuleGroupSummary objects for the rule groups that you manage.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRuleGroupsResponse> listRuleGroups(
    $0.ListRuleGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRuleGroups, request, options: options);
  }

  /// Retrieves the TagInfoForResource for the specified resource. Tags are key:value pairs that you can use to categorize and manage your resources, for purposes like billing. For example, you might set...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Retrieves an array of WebACLSummary objects for the web ACLs that you manage.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListWebACLsResponse> listWebACLs(
    $0.ListWebACLsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listWebACLs, request, options: options);
  }

  /// Enables the specified LoggingConfiguration, to start logging from a web ACL, according to the configuration provided. If you configure data protection for the web ACL, the protection applies to the...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutLoggingConfigurationResponse>
      putLoggingConfiguration(
    $0.PutLoggingConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putLoggingConfiguration, request,
        options: options);
  }

  /// Defines the versions of your managed rule set that you are offering to the customers. Customers see your offerings as managed rule groups with versioning. This is intended for use only by vendors o...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutManagedRuleSetVersionsResponse>
      putManagedRuleSetVersions(
    $0.PutManagedRuleSetVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putManagedRuleSetVersions, request,
        options: options);
  }

  /// Use this to share a rule group with other accounts. This action attaches an IAM policy to the specified resource. You must be the owner of the rule group to perform this operation. This action is s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutPermissionPolicyResponse> putPermissionPolicy(
    $0.PutPermissionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putPermissionPolicy, request, options: options);
  }

  /// Associates tags with the specified Amazon Web Services resource. Tags are key:value pairs that you can use to categorize and manage your resources, for purposes like billing. For example, you might...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Disassociates tags from an Amazon Web Services resource. Tags are key:value pairs that you can associate with Amazon Web Services resources. For example, the tag key might be "customer" and the tag...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates the specified IPSet. This operation completely replaces the mutable specifications that you already have for the IP set with the ones that you provide to this call. To modify an IP set, do ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateIPSetResponse> updateIPSet(
    $0.UpdateIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateIPSet, request, options: options);
  }

  /// Updates the expiration information for your managed rule set. Use this to initiate the expiration of a managed rule group version. After you initiate expiration for a version, WAF excludes it from ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateManagedRuleSetVersionExpiryDateResponse>
      updateManagedRuleSetVersionExpiryDate(
    $0.UpdateManagedRuleSetVersionExpiryDateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateManagedRuleSetVersionExpiryDate, request,
        options: options);
  }

  /// Updates the specified RegexPatternSet. This operation completely replaces the mutable specifications that you already have for the regex pattern set with the ones that you provide to this call. To ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRegexPatternSetResponse> updateRegexPatternSet(
    $0.UpdateRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRegexPatternSet, request, options: options);
  }

  /// Updates the specified RuleGroup. This operation completely replaces the mutable specifications that you already have for the rule group with the ones that you provide to this call. To modify a rule...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRuleGroupResponse> updateRuleGroup(
    $0.UpdateRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRuleGroup, request, options: options);
  }

  /// Updates the specified WebACL. While updating a web ACL, WAF provides continuous coverage to the resources that you have associated with the web ACL. This operation completely replaces the mutable s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateWebACLResponse> updateWebACL(
    $0.UpdateWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateWebACL, request, options: options);
  }

  // method descriptors

  static final _$associateWebACL =
      $grpc.ClientMethod<$0.AssociateWebACLRequest, $0.AssociateWebACLResponse>(
          '/wafv2.WAFV2Service/AssociateWebACL',
          ($0.AssociateWebACLRequest value) => value.writeToBuffer(),
          $0.AssociateWebACLResponse.fromBuffer);
  static final _$checkCapacity =
      $grpc.ClientMethod<$0.CheckCapacityRequest, $0.CheckCapacityResponse>(
          '/wafv2.WAFV2Service/CheckCapacity',
          ($0.CheckCapacityRequest value) => value.writeToBuffer(),
          $0.CheckCapacityResponse.fromBuffer);
  static final _$createAPIKey =
      $grpc.ClientMethod<$0.CreateAPIKeyRequest, $0.CreateAPIKeyResponse>(
          '/wafv2.WAFV2Service/CreateAPIKey',
          ($0.CreateAPIKeyRequest value) => value.writeToBuffer(),
          $0.CreateAPIKeyResponse.fromBuffer);
  static final _$createIPSet =
      $grpc.ClientMethod<$0.CreateIPSetRequest, $0.CreateIPSetResponse>(
          '/wafv2.WAFV2Service/CreateIPSet',
          ($0.CreateIPSetRequest value) => value.writeToBuffer(),
          $0.CreateIPSetResponse.fromBuffer);
  static final _$createRegexPatternSet = $grpc.ClientMethod<
          $0.CreateRegexPatternSetRequest, $0.CreateRegexPatternSetResponse>(
      '/wafv2.WAFV2Service/CreateRegexPatternSet',
      ($0.CreateRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.CreateRegexPatternSetResponse.fromBuffer);
  static final _$createRuleGroup =
      $grpc.ClientMethod<$0.CreateRuleGroupRequest, $0.CreateRuleGroupResponse>(
          '/wafv2.WAFV2Service/CreateRuleGroup',
          ($0.CreateRuleGroupRequest value) => value.writeToBuffer(),
          $0.CreateRuleGroupResponse.fromBuffer);
  static final _$createWebACL =
      $grpc.ClientMethod<$0.CreateWebACLRequest, $0.CreateWebACLResponse>(
          '/wafv2.WAFV2Service/CreateWebACL',
          ($0.CreateWebACLRequest value) => value.writeToBuffer(),
          $0.CreateWebACLResponse.fromBuffer);
  static final _$deleteAPIKey =
      $grpc.ClientMethod<$0.DeleteAPIKeyRequest, $0.DeleteAPIKeyResponse>(
          '/wafv2.WAFV2Service/DeleteAPIKey',
          ($0.DeleteAPIKeyRequest value) => value.writeToBuffer(),
          $0.DeleteAPIKeyResponse.fromBuffer);
  static final _$deleteFirewallManagerRuleGroups = $grpc.ClientMethod<
          $0.DeleteFirewallManagerRuleGroupsRequest,
          $0.DeleteFirewallManagerRuleGroupsResponse>(
      '/wafv2.WAFV2Service/DeleteFirewallManagerRuleGroups',
      ($0.DeleteFirewallManagerRuleGroupsRequest value) =>
          value.writeToBuffer(),
      $0.DeleteFirewallManagerRuleGroupsResponse.fromBuffer);
  static final _$deleteIPSet =
      $grpc.ClientMethod<$0.DeleteIPSetRequest, $0.DeleteIPSetResponse>(
          '/wafv2.WAFV2Service/DeleteIPSet',
          ($0.DeleteIPSetRequest value) => value.writeToBuffer(),
          $0.DeleteIPSetResponse.fromBuffer);
  static final _$deleteLoggingConfiguration = $grpc.ClientMethod<
          $0.DeleteLoggingConfigurationRequest,
          $0.DeleteLoggingConfigurationResponse>(
      '/wafv2.WAFV2Service/DeleteLoggingConfiguration',
      ($0.DeleteLoggingConfigurationRequest value) => value.writeToBuffer(),
      $0.DeleteLoggingConfigurationResponse.fromBuffer);
  static final _$deletePermissionPolicy = $grpc.ClientMethod<
          $0.DeletePermissionPolicyRequest, $0.DeletePermissionPolicyResponse>(
      '/wafv2.WAFV2Service/DeletePermissionPolicy',
      ($0.DeletePermissionPolicyRequest value) => value.writeToBuffer(),
      $0.DeletePermissionPolicyResponse.fromBuffer);
  static final _$deleteRegexPatternSet = $grpc.ClientMethod<
          $0.DeleteRegexPatternSetRequest, $0.DeleteRegexPatternSetResponse>(
      '/wafv2.WAFV2Service/DeleteRegexPatternSet',
      ($0.DeleteRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.DeleteRegexPatternSetResponse.fromBuffer);
  static final _$deleteRuleGroup =
      $grpc.ClientMethod<$0.DeleteRuleGroupRequest, $0.DeleteRuleGroupResponse>(
          '/wafv2.WAFV2Service/DeleteRuleGroup',
          ($0.DeleteRuleGroupRequest value) => value.writeToBuffer(),
          $0.DeleteRuleGroupResponse.fromBuffer);
  static final _$deleteWebACL =
      $grpc.ClientMethod<$0.DeleteWebACLRequest, $0.DeleteWebACLResponse>(
          '/wafv2.WAFV2Service/DeleteWebACL',
          ($0.DeleteWebACLRequest value) => value.writeToBuffer(),
          $0.DeleteWebACLResponse.fromBuffer);
  static final _$describeAllManagedProducts = $grpc.ClientMethod<
          $0.DescribeAllManagedProductsRequest,
          $0.DescribeAllManagedProductsResponse>(
      '/wafv2.WAFV2Service/DescribeAllManagedProducts',
      ($0.DescribeAllManagedProductsRequest value) => value.writeToBuffer(),
      $0.DescribeAllManagedProductsResponse.fromBuffer);
  static final _$describeManagedProductsByVendor = $grpc.ClientMethod<
          $0.DescribeManagedProductsByVendorRequest,
          $0.DescribeManagedProductsByVendorResponse>(
      '/wafv2.WAFV2Service/DescribeManagedProductsByVendor',
      ($0.DescribeManagedProductsByVendorRequest value) =>
          value.writeToBuffer(),
      $0.DescribeManagedProductsByVendorResponse.fromBuffer);
  static final _$describeManagedRuleGroup = $grpc.ClientMethod<
          $0.DescribeManagedRuleGroupRequest,
          $0.DescribeManagedRuleGroupResponse>(
      '/wafv2.WAFV2Service/DescribeManagedRuleGroup',
      ($0.DescribeManagedRuleGroupRequest value) => value.writeToBuffer(),
      $0.DescribeManagedRuleGroupResponse.fromBuffer);
  static final _$disassociateWebACL = $grpc.ClientMethod<
          $0.DisassociateWebACLRequest, $0.DisassociateWebACLResponse>(
      '/wafv2.WAFV2Service/DisassociateWebACL',
      ($0.DisassociateWebACLRequest value) => value.writeToBuffer(),
      $0.DisassociateWebACLResponse.fromBuffer);
  static final _$generateMobileSdkReleaseUrl = $grpc.ClientMethod<
          $0.GenerateMobileSdkReleaseUrlRequest,
          $0.GenerateMobileSdkReleaseUrlResponse>(
      '/wafv2.WAFV2Service/GenerateMobileSdkReleaseUrl',
      ($0.GenerateMobileSdkReleaseUrlRequest value) => value.writeToBuffer(),
      $0.GenerateMobileSdkReleaseUrlResponse.fromBuffer);
  static final _$getDecryptedAPIKey = $grpc.ClientMethod<
          $0.GetDecryptedAPIKeyRequest, $0.GetDecryptedAPIKeyResponse>(
      '/wafv2.WAFV2Service/GetDecryptedAPIKey',
      ($0.GetDecryptedAPIKeyRequest value) => value.writeToBuffer(),
      $0.GetDecryptedAPIKeyResponse.fromBuffer);
  static final _$getIPSet =
      $grpc.ClientMethod<$0.GetIPSetRequest, $0.GetIPSetResponse>(
          '/wafv2.WAFV2Service/GetIPSet',
          ($0.GetIPSetRequest value) => value.writeToBuffer(),
          $0.GetIPSetResponse.fromBuffer);
  static final _$getLoggingConfiguration = $grpc.ClientMethod<
          $0.GetLoggingConfigurationRequest,
          $0.GetLoggingConfigurationResponse>(
      '/wafv2.WAFV2Service/GetLoggingConfiguration',
      ($0.GetLoggingConfigurationRequest value) => value.writeToBuffer(),
      $0.GetLoggingConfigurationResponse.fromBuffer);
  static final _$getManagedRuleSet = $grpc.ClientMethod<
          $0.GetManagedRuleSetRequest, $0.GetManagedRuleSetResponse>(
      '/wafv2.WAFV2Service/GetManagedRuleSet',
      ($0.GetManagedRuleSetRequest value) => value.writeToBuffer(),
      $0.GetManagedRuleSetResponse.fromBuffer);
  static final _$getMobileSdkRelease = $grpc.ClientMethod<
          $0.GetMobileSdkReleaseRequest, $0.GetMobileSdkReleaseResponse>(
      '/wafv2.WAFV2Service/GetMobileSdkRelease',
      ($0.GetMobileSdkReleaseRequest value) => value.writeToBuffer(),
      $0.GetMobileSdkReleaseResponse.fromBuffer);
  static final _$getPermissionPolicy = $grpc.ClientMethod<
          $0.GetPermissionPolicyRequest, $0.GetPermissionPolicyResponse>(
      '/wafv2.WAFV2Service/GetPermissionPolicy',
      ($0.GetPermissionPolicyRequest value) => value.writeToBuffer(),
      $0.GetPermissionPolicyResponse.fromBuffer);
  static final _$getRateBasedStatementManagedKeys = $grpc.ClientMethod<
          $0.GetRateBasedStatementManagedKeysRequest,
          $0.GetRateBasedStatementManagedKeysResponse>(
      '/wafv2.WAFV2Service/GetRateBasedStatementManagedKeys',
      ($0.GetRateBasedStatementManagedKeysRequest value) =>
          value.writeToBuffer(),
      $0.GetRateBasedStatementManagedKeysResponse.fromBuffer);
  static final _$getRegexPatternSet = $grpc.ClientMethod<
          $0.GetRegexPatternSetRequest, $0.GetRegexPatternSetResponse>(
      '/wafv2.WAFV2Service/GetRegexPatternSet',
      ($0.GetRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.GetRegexPatternSetResponse.fromBuffer);
  static final _$getRuleGroup =
      $grpc.ClientMethod<$0.GetRuleGroupRequest, $0.GetRuleGroupResponse>(
          '/wafv2.WAFV2Service/GetRuleGroup',
          ($0.GetRuleGroupRequest value) => value.writeToBuffer(),
          $0.GetRuleGroupResponse.fromBuffer);
  static final _$getSampledRequests = $grpc.ClientMethod<
          $0.GetSampledRequestsRequest, $0.GetSampledRequestsResponse>(
      '/wafv2.WAFV2Service/GetSampledRequests',
      ($0.GetSampledRequestsRequest value) => value.writeToBuffer(),
      $0.GetSampledRequestsResponse.fromBuffer);
  static final _$getTopPathStatisticsByTraffic = $grpc.ClientMethod<
          $0.GetTopPathStatisticsByTrafficRequest,
          $0.GetTopPathStatisticsByTrafficResponse>(
      '/wafv2.WAFV2Service/GetTopPathStatisticsByTraffic',
      ($0.GetTopPathStatisticsByTrafficRequest value) => value.writeToBuffer(),
      $0.GetTopPathStatisticsByTrafficResponse.fromBuffer);
  static final _$getWebACL =
      $grpc.ClientMethod<$0.GetWebACLRequest, $0.GetWebACLResponse>(
          '/wafv2.WAFV2Service/GetWebACL',
          ($0.GetWebACLRequest value) => value.writeToBuffer(),
          $0.GetWebACLResponse.fromBuffer);
  static final _$getWebACLForResource = $grpc.ClientMethod<
          $0.GetWebACLForResourceRequest, $0.GetWebACLForResourceResponse>(
      '/wafv2.WAFV2Service/GetWebACLForResource',
      ($0.GetWebACLForResourceRequest value) => value.writeToBuffer(),
      $0.GetWebACLForResourceResponse.fromBuffer);
  static final _$listAPIKeys =
      $grpc.ClientMethod<$0.ListAPIKeysRequest, $0.ListAPIKeysResponse>(
          '/wafv2.WAFV2Service/ListAPIKeys',
          ($0.ListAPIKeysRequest value) => value.writeToBuffer(),
          $0.ListAPIKeysResponse.fromBuffer);
  static final _$listAvailableManagedRuleGroups = $grpc.ClientMethod<
          $0.ListAvailableManagedRuleGroupsRequest,
          $0.ListAvailableManagedRuleGroupsResponse>(
      '/wafv2.WAFV2Service/ListAvailableManagedRuleGroups',
      ($0.ListAvailableManagedRuleGroupsRequest value) => value.writeToBuffer(),
      $0.ListAvailableManagedRuleGroupsResponse.fromBuffer);
  static final _$listAvailableManagedRuleGroupVersions = $grpc.ClientMethod<
          $0.ListAvailableManagedRuleGroupVersionsRequest,
          $0.ListAvailableManagedRuleGroupVersionsResponse>(
      '/wafv2.WAFV2Service/ListAvailableManagedRuleGroupVersions',
      ($0.ListAvailableManagedRuleGroupVersionsRequest value) =>
          value.writeToBuffer(),
      $0.ListAvailableManagedRuleGroupVersionsResponse.fromBuffer);
  static final _$listIPSets =
      $grpc.ClientMethod<$0.ListIPSetsRequest, $0.ListIPSetsResponse>(
          '/wafv2.WAFV2Service/ListIPSets',
          ($0.ListIPSetsRequest value) => value.writeToBuffer(),
          $0.ListIPSetsResponse.fromBuffer);
  static final _$listLoggingConfigurations = $grpc.ClientMethod<
          $0.ListLoggingConfigurationsRequest,
          $0.ListLoggingConfigurationsResponse>(
      '/wafv2.WAFV2Service/ListLoggingConfigurations',
      ($0.ListLoggingConfigurationsRequest value) => value.writeToBuffer(),
      $0.ListLoggingConfigurationsResponse.fromBuffer);
  static final _$listManagedRuleSets = $grpc.ClientMethod<
          $0.ListManagedRuleSetsRequest, $0.ListManagedRuleSetsResponse>(
      '/wafv2.WAFV2Service/ListManagedRuleSets',
      ($0.ListManagedRuleSetsRequest value) => value.writeToBuffer(),
      $0.ListManagedRuleSetsResponse.fromBuffer);
  static final _$listMobileSdkReleases = $grpc.ClientMethod<
          $0.ListMobileSdkReleasesRequest, $0.ListMobileSdkReleasesResponse>(
      '/wafv2.WAFV2Service/ListMobileSdkReleases',
      ($0.ListMobileSdkReleasesRequest value) => value.writeToBuffer(),
      $0.ListMobileSdkReleasesResponse.fromBuffer);
  static final _$listRegexPatternSets = $grpc.ClientMethod<
          $0.ListRegexPatternSetsRequest, $0.ListRegexPatternSetsResponse>(
      '/wafv2.WAFV2Service/ListRegexPatternSets',
      ($0.ListRegexPatternSetsRequest value) => value.writeToBuffer(),
      $0.ListRegexPatternSetsResponse.fromBuffer);
  static final _$listResourcesForWebACL = $grpc.ClientMethod<
          $0.ListResourcesForWebACLRequest, $0.ListResourcesForWebACLResponse>(
      '/wafv2.WAFV2Service/ListResourcesForWebACL',
      ($0.ListResourcesForWebACLRequest value) => value.writeToBuffer(),
      $0.ListResourcesForWebACLResponse.fromBuffer);
  static final _$listRuleGroups =
      $grpc.ClientMethod<$0.ListRuleGroupsRequest, $0.ListRuleGroupsResponse>(
          '/wafv2.WAFV2Service/ListRuleGroups',
          ($0.ListRuleGroupsRequest value) => value.writeToBuffer(),
          $0.ListRuleGroupsResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/wafv2.WAFV2Service/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listWebACLs =
      $grpc.ClientMethod<$0.ListWebACLsRequest, $0.ListWebACLsResponse>(
          '/wafv2.WAFV2Service/ListWebACLs',
          ($0.ListWebACLsRequest value) => value.writeToBuffer(),
          $0.ListWebACLsResponse.fromBuffer);
  static final _$putLoggingConfiguration = $grpc.ClientMethod<
          $0.PutLoggingConfigurationRequest,
          $0.PutLoggingConfigurationResponse>(
      '/wafv2.WAFV2Service/PutLoggingConfiguration',
      ($0.PutLoggingConfigurationRequest value) => value.writeToBuffer(),
      $0.PutLoggingConfigurationResponse.fromBuffer);
  static final _$putManagedRuleSetVersions = $grpc.ClientMethod<
          $0.PutManagedRuleSetVersionsRequest,
          $0.PutManagedRuleSetVersionsResponse>(
      '/wafv2.WAFV2Service/PutManagedRuleSetVersions',
      ($0.PutManagedRuleSetVersionsRequest value) => value.writeToBuffer(),
      $0.PutManagedRuleSetVersionsResponse.fromBuffer);
  static final _$putPermissionPolicy = $grpc.ClientMethod<
          $0.PutPermissionPolicyRequest, $0.PutPermissionPolicyResponse>(
      '/wafv2.WAFV2Service/PutPermissionPolicy',
      ($0.PutPermissionPolicyRequest value) => value.writeToBuffer(),
      $0.PutPermissionPolicyResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/wafv2.WAFV2Service/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/wafv2.WAFV2Service/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateIPSet =
      $grpc.ClientMethod<$0.UpdateIPSetRequest, $0.UpdateIPSetResponse>(
          '/wafv2.WAFV2Service/UpdateIPSet',
          ($0.UpdateIPSetRequest value) => value.writeToBuffer(),
          $0.UpdateIPSetResponse.fromBuffer);
  static final _$updateManagedRuleSetVersionExpiryDate = $grpc.ClientMethod<
          $0.UpdateManagedRuleSetVersionExpiryDateRequest,
          $0.UpdateManagedRuleSetVersionExpiryDateResponse>(
      '/wafv2.WAFV2Service/UpdateManagedRuleSetVersionExpiryDate',
      ($0.UpdateManagedRuleSetVersionExpiryDateRequest value) =>
          value.writeToBuffer(),
      $0.UpdateManagedRuleSetVersionExpiryDateResponse.fromBuffer);
  static final _$updateRegexPatternSet = $grpc.ClientMethod<
          $0.UpdateRegexPatternSetRequest, $0.UpdateRegexPatternSetResponse>(
      '/wafv2.WAFV2Service/UpdateRegexPatternSet',
      ($0.UpdateRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.UpdateRegexPatternSetResponse.fromBuffer);
  static final _$updateRuleGroup =
      $grpc.ClientMethod<$0.UpdateRuleGroupRequest, $0.UpdateRuleGroupResponse>(
          '/wafv2.WAFV2Service/UpdateRuleGroup',
          ($0.UpdateRuleGroupRequest value) => value.writeToBuffer(),
          $0.UpdateRuleGroupResponse.fromBuffer);
  static final _$updateWebACL =
      $grpc.ClientMethod<$0.UpdateWebACLRequest, $0.UpdateWebACLResponse>(
          '/wafv2.WAFV2Service/UpdateWebACL',
          ($0.UpdateWebACLRequest value) => value.writeToBuffer(),
          $0.UpdateWebACLResponse.fromBuffer);
}

@$pb.GrpcServiceName('wafv2.WAFV2Service')
abstract class WAFV2ServiceBase extends $grpc.Service {
  $core.String get $name => 'wafv2.WAFV2Service';

  WAFV2ServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AssociateWebACLRequest,
            $0.AssociateWebACLResponse>(
        'AssociateWebACL',
        associateWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateWebACLRequest.fromBuffer(value),
        ($0.AssociateWebACLResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CheckCapacityRequest, $0.CheckCapacityResponse>(
            'CheckCapacity',
            checkCapacity_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CheckCapacityRequest.fromBuffer(value),
            ($0.CheckCapacityResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateAPIKeyRequest, $0.CreateAPIKeyResponse>(
            'CreateAPIKey',
            createAPIKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateAPIKeyRequest.fromBuffer(value),
            ($0.CreateAPIKeyResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateIPSetRequest, $0.CreateIPSetResponse>(
            'CreateIPSet',
            createIPSet_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateIPSetRequest.fromBuffer(value),
            ($0.CreateIPSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRegexPatternSetRequest,
            $0.CreateRegexPatternSetResponse>(
        'CreateRegexPatternSet',
        createRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRegexPatternSetRequest.fromBuffer(value),
        ($0.CreateRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRuleGroupRequest,
            $0.CreateRuleGroupResponse>(
        'CreateRuleGroup',
        createRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRuleGroupRequest.fromBuffer(value),
        ($0.CreateRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateWebACLRequest, $0.CreateWebACLResponse>(
            'CreateWebACL',
            createWebACL_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateWebACLRequest.fromBuffer(value),
            ($0.CreateWebACLResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteAPIKeyRequest, $0.DeleteAPIKeyResponse>(
            'DeleteAPIKey',
            deleteAPIKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteAPIKeyRequest.fromBuffer(value),
            ($0.DeleteAPIKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFirewallManagerRuleGroupsRequest,
            $0.DeleteFirewallManagerRuleGroupsResponse>(
        'DeleteFirewallManagerRuleGroups',
        deleteFirewallManagerRuleGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFirewallManagerRuleGroupsRequest.fromBuffer(value),
        ($0.DeleteFirewallManagerRuleGroupsResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteIPSetRequest, $0.DeleteIPSetResponse>(
            'DeleteIPSet',
            deleteIPSet_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteIPSetRequest.fromBuffer(value),
            ($0.DeleteIPSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteLoggingConfigurationRequest,
            $0.DeleteLoggingConfigurationResponse>(
        'DeleteLoggingConfiguration',
        deleteLoggingConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteLoggingConfigurationRequest.fromBuffer(value),
        ($0.DeleteLoggingConfigurationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePermissionPolicyRequest,
            $0.DeletePermissionPolicyResponse>(
        'DeletePermissionPolicy',
        deletePermissionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePermissionPolicyRequest.fromBuffer(value),
        ($0.DeletePermissionPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRegexPatternSetRequest,
            $0.DeleteRegexPatternSetResponse>(
        'DeleteRegexPatternSet',
        deleteRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRegexPatternSetRequest.fromBuffer(value),
        ($0.DeleteRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRuleGroupRequest,
            $0.DeleteRuleGroupResponse>(
        'DeleteRuleGroup',
        deleteRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRuleGroupRequest.fromBuffer(value),
        ($0.DeleteRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteWebACLRequest, $0.DeleteWebACLResponse>(
            'DeleteWebACL',
            deleteWebACL_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteWebACLRequest.fromBuffer(value),
            ($0.DeleteWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAllManagedProductsRequest,
            $0.DescribeAllManagedProductsResponse>(
        'DescribeAllManagedProducts',
        describeAllManagedProducts_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAllManagedProductsRequest.fromBuffer(value),
        ($0.DescribeAllManagedProductsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeManagedProductsByVendorRequest,
            $0.DescribeManagedProductsByVendorResponse>(
        'DescribeManagedProductsByVendor',
        describeManagedProductsByVendor_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeManagedProductsByVendorRequest.fromBuffer(value),
        ($0.DescribeManagedProductsByVendorResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeManagedRuleGroupRequest,
            $0.DescribeManagedRuleGroupResponse>(
        'DescribeManagedRuleGroup',
        describeManagedRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeManagedRuleGroupRequest.fromBuffer(value),
        ($0.DescribeManagedRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisassociateWebACLRequest,
            $0.DisassociateWebACLResponse>(
        'DisassociateWebACL',
        disassociateWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateWebACLRequest.fromBuffer(value),
        ($0.DisassociateWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateMobileSdkReleaseUrlRequest,
            $0.GenerateMobileSdkReleaseUrlResponse>(
        'GenerateMobileSdkReleaseUrl',
        generateMobileSdkReleaseUrl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateMobileSdkReleaseUrlRequest.fromBuffer(value),
        ($0.GenerateMobileSdkReleaseUrlResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDecryptedAPIKeyRequest,
            $0.GetDecryptedAPIKeyResponse>(
        'GetDecryptedAPIKey',
        getDecryptedAPIKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDecryptedAPIKeyRequest.fromBuffer(value),
        ($0.GetDecryptedAPIKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIPSetRequest, $0.GetIPSetResponse>(
        'GetIPSet',
        getIPSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetIPSetRequest.fromBuffer(value),
        ($0.GetIPSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetLoggingConfigurationRequest,
            $0.GetLoggingConfigurationResponse>(
        'GetLoggingConfiguration',
        getLoggingConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetLoggingConfigurationRequest.fromBuffer(value),
        ($0.GetLoggingConfigurationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetManagedRuleSetRequest,
            $0.GetManagedRuleSetResponse>(
        'GetManagedRuleSet',
        getManagedRuleSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetManagedRuleSetRequest.fromBuffer(value),
        ($0.GetManagedRuleSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMobileSdkReleaseRequest,
            $0.GetMobileSdkReleaseResponse>(
        'GetMobileSdkRelease',
        getMobileSdkRelease_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMobileSdkReleaseRequest.fromBuffer(value),
        ($0.GetMobileSdkReleaseResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPermissionPolicyRequest,
            $0.GetPermissionPolicyResponse>(
        'GetPermissionPolicy',
        getPermissionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPermissionPolicyRequest.fromBuffer(value),
        ($0.GetPermissionPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRateBasedStatementManagedKeysRequest,
            $0.GetRateBasedStatementManagedKeysResponse>(
        'GetRateBasedStatementManagedKeys',
        getRateBasedStatementManagedKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRateBasedStatementManagedKeysRequest.fromBuffer(value),
        ($0.GetRateBasedStatementManagedKeysResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRegexPatternSetRequest,
            $0.GetRegexPatternSetResponse>(
        'GetRegexPatternSet',
        getRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRegexPatternSetRequest.fromBuffer(value),
        ($0.GetRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetRuleGroupRequest, $0.GetRuleGroupResponse>(
            'GetRuleGroup',
            getRuleGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetRuleGroupRequest.fromBuffer(value),
            ($0.GetRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSampledRequestsRequest,
            $0.GetSampledRequestsResponse>(
        'GetSampledRequests',
        getSampledRequests_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSampledRequestsRequest.fromBuffer(value),
        ($0.GetSampledRequestsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTopPathStatisticsByTrafficRequest,
            $0.GetTopPathStatisticsByTrafficResponse>(
        'GetTopPathStatisticsByTraffic',
        getTopPathStatisticsByTraffic_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTopPathStatisticsByTrafficRequest.fromBuffer(value),
        ($0.GetTopPathStatisticsByTrafficResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetWebACLRequest, $0.GetWebACLResponse>(
        'GetWebACL',
        getWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetWebACLRequest.fromBuffer(value),
        ($0.GetWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetWebACLForResourceRequest,
            $0.GetWebACLForResourceResponse>(
        'GetWebACLForResource',
        getWebACLForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetWebACLForResourceRequest.fromBuffer(value),
        ($0.GetWebACLForResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListAPIKeysRequest, $0.ListAPIKeysResponse>(
            'ListAPIKeys',
            listAPIKeys_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListAPIKeysRequest.fromBuffer(value),
            ($0.ListAPIKeysResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAvailableManagedRuleGroupsRequest,
            $0.ListAvailableManagedRuleGroupsResponse>(
        'ListAvailableManagedRuleGroups',
        listAvailableManagedRuleGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAvailableManagedRuleGroupsRequest.fromBuffer(value),
        ($0.ListAvailableManagedRuleGroupsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListAvailableManagedRuleGroupVersionsRequest,
            $0.ListAvailableManagedRuleGroupVersionsResponse>(
        'ListAvailableManagedRuleGroupVersions',
        listAvailableManagedRuleGroupVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAvailableManagedRuleGroupVersionsRequest.fromBuffer(value),
        ($0.ListAvailableManagedRuleGroupVersionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListIPSetsRequest, $0.ListIPSetsResponse>(
        'ListIPSets',
        listIPSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListIPSetsRequest.fromBuffer(value),
        ($0.ListIPSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListLoggingConfigurationsRequest,
            $0.ListLoggingConfigurationsResponse>(
        'ListLoggingConfigurations',
        listLoggingConfigurations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListLoggingConfigurationsRequest.fromBuffer(value),
        ($0.ListLoggingConfigurationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListManagedRuleSetsRequest,
            $0.ListManagedRuleSetsResponse>(
        'ListManagedRuleSets',
        listManagedRuleSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListManagedRuleSetsRequest.fromBuffer(value),
        ($0.ListManagedRuleSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMobileSdkReleasesRequest,
            $0.ListMobileSdkReleasesResponse>(
        'ListMobileSdkReleases',
        listMobileSdkReleases_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMobileSdkReleasesRequest.fromBuffer(value),
        ($0.ListMobileSdkReleasesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRegexPatternSetsRequest,
            $0.ListRegexPatternSetsResponse>(
        'ListRegexPatternSets',
        listRegexPatternSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRegexPatternSetsRequest.fromBuffer(value),
        ($0.ListRegexPatternSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourcesForWebACLRequest,
            $0.ListResourcesForWebACLResponse>(
        'ListResourcesForWebACL',
        listResourcesForWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourcesForWebACLRequest.fromBuffer(value),
        ($0.ListResourcesForWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRuleGroupsRequest,
            $0.ListRuleGroupsResponse>(
        'ListRuleGroups',
        listRuleGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRuleGroupsRequest.fromBuffer(value),
        ($0.ListRuleGroupsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListWebACLsRequest, $0.ListWebACLsResponse>(
            'ListWebACLs',
            listWebACLs_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListWebACLsRequest.fromBuffer(value),
            ($0.ListWebACLsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutLoggingConfigurationRequest,
            $0.PutLoggingConfigurationResponse>(
        'PutLoggingConfiguration',
        putLoggingConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutLoggingConfigurationRequest.fromBuffer(value),
        ($0.PutLoggingConfigurationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutManagedRuleSetVersionsRequest,
            $0.PutManagedRuleSetVersionsResponse>(
        'PutManagedRuleSetVersions',
        putManagedRuleSetVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutManagedRuleSetVersionsRequest.fromBuffer(value),
        ($0.PutManagedRuleSetVersionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutPermissionPolicyRequest,
            $0.PutPermissionPolicyResponse>(
        'PutPermissionPolicy',
        putPermissionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutPermissionPolicyRequest.fromBuffer(value),
        ($0.PutPermissionPolicyResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
            'TagResource',
            tagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagResourceRequest.fromBuffer(value),
            ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceRequest.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateIPSetRequest, $0.UpdateIPSetResponse>(
            'UpdateIPSet',
            updateIPSet_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateIPSetRequest.fromBuffer(value),
            ($0.UpdateIPSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateManagedRuleSetVersionExpiryDateRequest,
            $0.UpdateManagedRuleSetVersionExpiryDateResponse>(
        'UpdateManagedRuleSetVersionExpiryDate',
        updateManagedRuleSetVersionExpiryDate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateManagedRuleSetVersionExpiryDateRequest.fromBuffer(value),
        ($0.UpdateManagedRuleSetVersionExpiryDateResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRegexPatternSetRequest,
            $0.UpdateRegexPatternSetResponse>(
        'UpdateRegexPatternSet',
        updateRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRegexPatternSetRequest.fromBuffer(value),
        ($0.UpdateRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRuleGroupRequest,
            $0.UpdateRuleGroupResponse>(
        'UpdateRuleGroup',
        updateRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRuleGroupRequest.fromBuffer(value),
        ($0.UpdateRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateWebACLRequest, $0.UpdateWebACLResponse>(
            'UpdateWebACL',
            updateWebACL_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateWebACLRequest.fromBuffer(value),
            ($0.UpdateWebACLResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AssociateWebACLResponse> associateWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AssociateWebACLRequest> $request) async {
    return associateWebACL($call, await $request);
  }

  $async.Future<$0.AssociateWebACLResponse> associateWebACL(
      $grpc.ServiceCall call, $0.AssociateWebACLRequest request);

  $async.Future<$0.CheckCapacityResponse> checkCapacity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CheckCapacityRequest> $request) async {
    return checkCapacity($call, await $request);
  }

  $async.Future<$0.CheckCapacityResponse> checkCapacity(
      $grpc.ServiceCall call, $0.CheckCapacityRequest request);

  $async.Future<$0.CreateAPIKeyResponse> createAPIKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateAPIKeyRequest> $request) async {
    return createAPIKey($call, await $request);
  }

  $async.Future<$0.CreateAPIKeyResponse> createAPIKey(
      $grpc.ServiceCall call, $0.CreateAPIKeyRequest request);

  $async.Future<$0.CreateIPSetResponse> createIPSet_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateIPSetRequest> $request) async {
    return createIPSet($call, await $request);
  }

  $async.Future<$0.CreateIPSetResponse> createIPSet(
      $grpc.ServiceCall call, $0.CreateIPSetRequest request);

  $async.Future<$0.CreateRegexPatternSetResponse> createRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRegexPatternSetRequest> $request) async {
    return createRegexPatternSet($call, await $request);
  }

  $async.Future<$0.CreateRegexPatternSetResponse> createRegexPatternSet(
      $grpc.ServiceCall call, $0.CreateRegexPatternSetRequest request);

  $async.Future<$0.CreateRuleGroupResponse> createRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRuleGroupRequest> $request) async {
    return createRuleGroup($call, await $request);
  }

  $async.Future<$0.CreateRuleGroupResponse> createRuleGroup(
      $grpc.ServiceCall call, $0.CreateRuleGroupRequest request);

  $async.Future<$0.CreateWebACLResponse> createWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateWebACLRequest> $request) async {
    return createWebACL($call, await $request);
  }

  $async.Future<$0.CreateWebACLResponse> createWebACL(
      $grpc.ServiceCall call, $0.CreateWebACLRequest request);

  $async.Future<$0.DeleteAPIKeyResponse> deleteAPIKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteAPIKeyRequest> $request) async {
    return deleteAPIKey($call, await $request);
  }

  $async.Future<$0.DeleteAPIKeyResponse> deleteAPIKey(
      $grpc.ServiceCall call, $0.DeleteAPIKeyRequest request);

  $async.Future<$0.DeleteFirewallManagerRuleGroupsResponse>
      deleteFirewallManagerRuleGroups_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeleteFirewallManagerRuleGroupsRequest>
              $request) async {
    return deleteFirewallManagerRuleGroups($call, await $request);
  }

  $async.Future<$0.DeleteFirewallManagerRuleGroupsResponse>
      deleteFirewallManagerRuleGroups($grpc.ServiceCall call,
          $0.DeleteFirewallManagerRuleGroupsRequest request);

  $async.Future<$0.DeleteIPSetResponse> deleteIPSet_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteIPSetRequest> $request) async {
    return deleteIPSet($call, await $request);
  }

  $async.Future<$0.DeleteIPSetResponse> deleteIPSet(
      $grpc.ServiceCall call, $0.DeleteIPSetRequest request);

  $async.Future<$0.DeleteLoggingConfigurationResponse>
      deleteLoggingConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteLoggingConfigurationRequest> $request) async {
    return deleteLoggingConfiguration($call, await $request);
  }

  $async.Future<$0.DeleteLoggingConfigurationResponse>
      deleteLoggingConfiguration(
          $grpc.ServiceCall call, $0.DeleteLoggingConfigurationRequest request);

  $async.Future<$0.DeletePermissionPolicyResponse> deletePermissionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeletePermissionPolicyRequest> $request) async {
    return deletePermissionPolicy($call, await $request);
  }

  $async.Future<$0.DeletePermissionPolicyResponse> deletePermissionPolicy(
      $grpc.ServiceCall call, $0.DeletePermissionPolicyRequest request);

  $async.Future<$0.DeleteRegexPatternSetResponse> deleteRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRegexPatternSetRequest> $request) async {
    return deleteRegexPatternSet($call, await $request);
  }

  $async.Future<$0.DeleteRegexPatternSetResponse> deleteRegexPatternSet(
      $grpc.ServiceCall call, $0.DeleteRegexPatternSetRequest request);

  $async.Future<$0.DeleteRuleGroupResponse> deleteRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRuleGroupRequest> $request) async {
    return deleteRuleGroup($call, await $request);
  }

  $async.Future<$0.DeleteRuleGroupResponse> deleteRuleGroup(
      $grpc.ServiceCall call, $0.DeleteRuleGroupRequest request);

  $async.Future<$0.DeleteWebACLResponse> deleteWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteWebACLRequest> $request) async {
    return deleteWebACL($call, await $request);
  }

  $async.Future<$0.DeleteWebACLResponse> deleteWebACL(
      $grpc.ServiceCall call, $0.DeleteWebACLRequest request);

  $async.Future<$0.DescribeAllManagedProductsResponse>
      describeAllManagedProducts_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeAllManagedProductsRequest> $request) async {
    return describeAllManagedProducts($call, await $request);
  }

  $async.Future<$0.DescribeAllManagedProductsResponse>
      describeAllManagedProducts(
          $grpc.ServiceCall call, $0.DescribeAllManagedProductsRequest request);

  $async.Future<$0.DescribeManagedProductsByVendorResponse>
      describeManagedProductsByVendor_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeManagedProductsByVendorRequest>
              $request) async {
    return describeManagedProductsByVendor($call, await $request);
  }

  $async.Future<$0.DescribeManagedProductsByVendorResponse>
      describeManagedProductsByVendor($grpc.ServiceCall call,
          $0.DescribeManagedProductsByVendorRequest request);

  $async.Future<$0.DescribeManagedRuleGroupResponse>
      describeManagedRuleGroup_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeManagedRuleGroupRequest> $request) async {
    return describeManagedRuleGroup($call, await $request);
  }

  $async.Future<$0.DescribeManagedRuleGroupResponse> describeManagedRuleGroup(
      $grpc.ServiceCall call, $0.DescribeManagedRuleGroupRequest request);

  $async.Future<$0.DisassociateWebACLResponse> disassociateWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DisassociateWebACLRequest> $request) async {
    return disassociateWebACL($call, await $request);
  }

  $async.Future<$0.DisassociateWebACLResponse> disassociateWebACL(
      $grpc.ServiceCall call, $0.DisassociateWebACLRequest request);

  $async.Future<$0.GenerateMobileSdkReleaseUrlResponse>
      generateMobileSdkReleaseUrl_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GenerateMobileSdkReleaseUrlRequest> $request) async {
    return generateMobileSdkReleaseUrl($call, await $request);
  }

  $async.Future<$0.GenerateMobileSdkReleaseUrlResponse>
      generateMobileSdkReleaseUrl($grpc.ServiceCall call,
          $0.GenerateMobileSdkReleaseUrlRequest request);

  $async.Future<$0.GetDecryptedAPIKeyResponse> getDecryptedAPIKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDecryptedAPIKeyRequest> $request) async {
    return getDecryptedAPIKey($call, await $request);
  }

  $async.Future<$0.GetDecryptedAPIKeyResponse> getDecryptedAPIKey(
      $grpc.ServiceCall call, $0.GetDecryptedAPIKeyRequest request);

  $async.Future<$0.GetIPSetResponse> getIPSet_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetIPSetRequest> $request) async {
    return getIPSet($call, await $request);
  }

  $async.Future<$0.GetIPSetResponse> getIPSet(
      $grpc.ServiceCall call, $0.GetIPSetRequest request);

  $async.Future<$0.GetLoggingConfigurationResponse> getLoggingConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLoggingConfigurationRequest> $request) async {
    return getLoggingConfiguration($call, await $request);
  }

  $async.Future<$0.GetLoggingConfigurationResponse> getLoggingConfiguration(
      $grpc.ServiceCall call, $0.GetLoggingConfigurationRequest request);

  $async.Future<$0.GetManagedRuleSetResponse> getManagedRuleSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetManagedRuleSetRequest> $request) async {
    return getManagedRuleSet($call, await $request);
  }

  $async.Future<$0.GetManagedRuleSetResponse> getManagedRuleSet(
      $grpc.ServiceCall call, $0.GetManagedRuleSetRequest request);

  $async.Future<$0.GetMobileSdkReleaseResponse> getMobileSdkRelease_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMobileSdkReleaseRequest> $request) async {
    return getMobileSdkRelease($call, await $request);
  }

  $async.Future<$0.GetMobileSdkReleaseResponse> getMobileSdkRelease(
      $grpc.ServiceCall call, $0.GetMobileSdkReleaseRequest request);

  $async.Future<$0.GetPermissionPolicyResponse> getPermissionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPermissionPolicyRequest> $request) async {
    return getPermissionPolicy($call, await $request);
  }

  $async.Future<$0.GetPermissionPolicyResponse> getPermissionPolicy(
      $grpc.ServiceCall call, $0.GetPermissionPolicyRequest request);

  $async.Future<$0.GetRateBasedStatementManagedKeysResponse>
      getRateBasedStatementManagedKeys_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetRateBasedStatementManagedKeysRequest>
              $request) async {
    return getRateBasedStatementManagedKeys($call, await $request);
  }

  $async.Future<$0.GetRateBasedStatementManagedKeysResponse>
      getRateBasedStatementManagedKeys($grpc.ServiceCall call,
          $0.GetRateBasedStatementManagedKeysRequest request);

  $async.Future<$0.GetRegexPatternSetResponse> getRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRegexPatternSetRequest> $request) async {
    return getRegexPatternSet($call, await $request);
  }

  $async.Future<$0.GetRegexPatternSetResponse> getRegexPatternSet(
      $grpc.ServiceCall call, $0.GetRegexPatternSetRequest request);

  $async.Future<$0.GetRuleGroupResponse> getRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRuleGroupRequest> $request) async {
    return getRuleGroup($call, await $request);
  }

  $async.Future<$0.GetRuleGroupResponse> getRuleGroup(
      $grpc.ServiceCall call, $0.GetRuleGroupRequest request);

  $async.Future<$0.GetSampledRequestsResponse> getSampledRequests_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSampledRequestsRequest> $request) async {
    return getSampledRequests($call, await $request);
  }

  $async.Future<$0.GetSampledRequestsResponse> getSampledRequests(
      $grpc.ServiceCall call, $0.GetSampledRequestsRequest request);

  $async.Future<$0.GetTopPathStatisticsByTrafficResponse>
      getTopPathStatisticsByTraffic_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetTopPathStatisticsByTrafficRequest>
              $request) async {
    return getTopPathStatisticsByTraffic($call, await $request);
  }

  $async.Future<$0.GetTopPathStatisticsByTrafficResponse>
      getTopPathStatisticsByTraffic($grpc.ServiceCall call,
          $0.GetTopPathStatisticsByTrafficRequest request);

  $async.Future<$0.GetWebACLResponse> getWebACL_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetWebACLRequest> $request) async {
    return getWebACL($call, await $request);
  }

  $async.Future<$0.GetWebACLResponse> getWebACL(
      $grpc.ServiceCall call, $0.GetWebACLRequest request);

  $async.Future<$0.GetWebACLForResourceResponse> getWebACLForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetWebACLForResourceRequest> $request) async {
    return getWebACLForResource($call, await $request);
  }

  $async.Future<$0.GetWebACLForResourceResponse> getWebACLForResource(
      $grpc.ServiceCall call, $0.GetWebACLForResourceRequest request);

  $async.Future<$0.ListAPIKeysResponse> listAPIKeys_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListAPIKeysRequest> $request) async {
    return listAPIKeys($call, await $request);
  }

  $async.Future<$0.ListAPIKeysResponse> listAPIKeys(
      $grpc.ServiceCall call, $0.ListAPIKeysRequest request);

  $async.Future<$0.ListAvailableManagedRuleGroupsResponse>
      listAvailableManagedRuleGroups_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListAvailableManagedRuleGroupsRequest>
              $request) async {
    return listAvailableManagedRuleGroups($call, await $request);
  }

  $async.Future<$0.ListAvailableManagedRuleGroupsResponse>
      listAvailableManagedRuleGroups($grpc.ServiceCall call,
          $0.ListAvailableManagedRuleGroupsRequest request);

  $async.Future<$0.ListAvailableManagedRuleGroupVersionsResponse>
      listAvailableManagedRuleGroupVersions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListAvailableManagedRuleGroupVersionsRequest>
              $request) async {
    return listAvailableManagedRuleGroupVersions($call, await $request);
  }

  $async.Future<$0.ListAvailableManagedRuleGroupVersionsResponse>
      listAvailableManagedRuleGroupVersions($grpc.ServiceCall call,
          $0.ListAvailableManagedRuleGroupVersionsRequest request);

  $async.Future<$0.ListIPSetsResponse> listIPSets_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListIPSetsRequest> $request) async {
    return listIPSets($call, await $request);
  }

  $async.Future<$0.ListIPSetsResponse> listIPSets(
      $grpc.ServiceCall call, $0.ListIPSetsRequest request);

  $async.Future<$0.ListLoggingConfigurationsResponse>
      listLoggingConfigurations_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListLoggingConfigurationsRequest> $request) async {
    return listLoggingConfigurations($call, await $request);
  }

  $async.Future<$0.ListLoggingConfigurationsResponse> listLoggingConfigurations(
      $grpc.ServiceCall call, $0.ListLoggingConfigurationsRequest request);

  $async.Future<$0.ListManagedRuleSetsResponse> listManagedRuleSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListManagedRuleSetsRequest> $request) async {
    return listManagedRuleSets($call, await $request);
  }

  $async.Future<$0.ListManagedRuleSetsResponse> listManagedRuleSets(
      $grpc.ServiceCall call, $0.ListManagedRuleSetsRequest request);

  $async.Future<$0.ListMobileSdkReleasesResponse> listMobileSdkReleases_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListMobileSdkReleasesRequest> $request) async {
    return listMobileSdkReleases($call, await $request);
  }

  $async.Future<$0.ListMobileSdkReleasesResponse> listMobileSdkReleases(
      $grpc.ServiceCall call, $0.ListMobileSdkReleasesRequest request);

  $async.Future<$0.ListRegexPatternSetsResponse> listRegexPatternSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRegexPatternSetsRequest> $request) async {
    return listRegexPatternSets($call, await $request);
  }

  $async.Future<$0.ListRegexPatternSetsResponse> listRegexPatternSets(
      $grpc.ServiceCall call, $0.ListRegexPatternSetsRequest request);

  $async.Future<$0.ListResourcesForWebACLResponse> listResourcesForWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListResourcesForWebACLRequest> $request) async {
    return listResourcesForWebACL($call, await $request);
  }

  $async.Future<$0.ListResourcesForWebACLResponse> listResourcesForWebACL(
      $grpc.ServiceCall call, $0.ListResourcesForWebACLRequest request);

  $async.Future<$0.ListRuleGroupsResponse> listRuleGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRuleGroupsRequest> $request) async {
    return listRuleGroups($call, await $request);
  }

  $async.Future<$0.ListRuleGroupsResponse> listRuleGroups(
      $grpc.ServiceCall call, $0.ListRuleGroupsRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListWebACLsResponse> listWebACLs_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListWebACLsRequest> $request) async {
    return listWebACLs($call, await $request);
  }

  $async.Future<$0.ListWebACLsResponse> listWebACLs(
      $grpc.ServiceCall call, $0.ListWebACLsRequest request);

  $async.Future<$0.PutLoggingConfigurationResponse> putLoggingConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutLoggingConfigurationRequest> $request) async {
    return putLoggingConfiguration($call, await $request);
  }

  $async.Future<$0.PutLoggingConfigurationResponse> putLoggingConfiguration(
      $grpc.ServiceCall call, $0.PutLoggingConfigurationRequest request);

  $async.Future<$0.PutManagedRuleSetVersionsResponse>
      putManagedRuleSetVersions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.PutManagedRuleSetVersionsRequest> $request) async {
    return putManagedRuleSetVersions($call, await $request);
  }

  $async.Future<$0.PutManagedRuleSetVersionsResponse> putManagedRuleSetVersions(
      $grpc.ServiceCall call, $0.PutManagedRuleSetVersionsRequest request);

  $async.Future<$0.PutPermissionPolicyResponse> putPermissionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutPermissionPolicyRequest> $request) async {
    return putPermissionPolicy($call, await $request);
  }

  $async.Future<$0.PutPermissionPolicyResponse> putPermissionPolicy(
      $grpc.ServiceCall call, $0.PutPermissionPolicyRequest request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateIPSetResponse> updateIPSet_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateIPSetRequest> $request) async {
    return updateIPSet($call, await $request);
  }

  $async.Future<$0.UpdateIPSetResponse> updateIPSet(
      $grpc.ServiceCall call, $0.UpdateIPSetRequest request);

  $async.Future<$0.UpdateManagedRuleSetVersionExpiryDateResponse>
      updateManagedRuleSetVersionExpiryDate_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateManagedRuleSetVersionExpiryDateRequest>
              $request) async {
    return updateManagedRuleSetVersionExpiryDate($call, await $request);
  }

  $async.Future<$0.UpdateManagedRuleSetVersionExpiryDateResponse>
      updateManagedRuleSetVersionExpiryDate($grpc.ServiceCall call,
          $0.UpdateManagedRuleSetVersionExpiryDateRequest request);

  $async.Future<$0.UpdateRegexPatternSetResponse> updateRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRegexPatternSetRequest> $request) async {
    return updateRegexPatternSet($call, await $request);
  }

  $async.Future<$0.UpdateRegexPatternSetResponse> updateRegexPatternSet(
      $grpc.ServiceCall call, $0.UpdateRegexPatternSetRequest request);

  $async.Future<$0.UpdateRuleGroupResponse> updateRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRuleGroupRequest> $request) async {
    return updateRuleGroup($call, await $request);
  }

  $async.Future<$0.UpdateRuleGroupResponse> updateRuleGroup(
      $grpc.ServiceCall call, $0.UpdateRuleGroupRequest request);

  $async.Future<$0.UpdateWebACLResponse> updateWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateWebACLRequest> $request) async {
    return updateWebACL($call, await $request);
  }

  $async.Future<$0.UpdateWebACLResponse> updateWebACL(
      $grpc.ServiceCall call, $0.UpdateWebACLRequest request);
}
