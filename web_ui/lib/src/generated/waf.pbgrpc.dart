// This is a generated file - do not edit.
//
// Generated from waf.proto.

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

import 'waf.pb.dart' as $0;

export 'waf.pb.dart';

/// WAFService provides waf API operations.
@$pb.GrpcServiceName('waf.WAFService')
class WAFServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  WAFServiceClient(super.channel, {super.options, super.interceptors});

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateByteMatchSetResponse> createByteMatchSet(
    $0.CreateByteMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createByteMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateGeoMatchSetResponse> createGeoMatchSet(
    $0.CreateGeoMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createGeoMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateIPSetResponse> createIPSet(
    $0.CreateIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createIPSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRateBasedRuleResponse> createRateBasedRule(
    $0.CreateRateBasedRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRateBasedRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRegexMatchSetResponse> createRegexMatchSet(
    $0.CreateRegexMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRegexMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRegexPatternSetResponse> createRegexPatternSet(
    $0.CreateRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRegexPatternSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRuleResponse> createRule(
    $0.CreateRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateRuleGroupResponse> createRuleGroup(
    $0.CreateRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRuleGroup, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateSizeConstraintSetResponse>
      createSizeConstraintSet(
    $0.CreateSizeConstraintSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSizeConstraintSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateSqlInjectionMatchSetResponse>
      createSqlInjectionMatchSet(
    $0.CreateSqlInjectionMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSqlInjectionMatchSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateWebACLResponse> createWebACL(
    $0.CreateWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createWebACL, request, options: options);
  }

  /// Creates an AWS CloudFormation WAFV2 template for the specified web ACL in the specified Amazon S3 bucket. Then, in CloudFormation, you create a stack from the template, to create the web ACL and it...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateWebACLMigrationStackResponse>
      createWebACLMigrationStack(
    $0.CreateWebACLMigrationStackRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createWebACLMigrationStack, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateXssMatchSetResponse> createXssMatchSet(
    $0.CreateXssMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createXssMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteByteMatchSetResponse> deleteByteMatchSet(
    $0.DeleteByteMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteByteMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteGeoMatchSetResponse> deleteGeoMatchSet(
    $0.DeleteGeoMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteGeoMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteIPSetResponse> deleteIPSet(
    $0.DeleteIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIPSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
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

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
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

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRateBasedRuleResponse> deleteRateBasedRule(
    $0.DeleteRateBasedRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRateBasedRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRegexMatchSetResponse> deleteRegexMatchSet(
    $0.DeleteRegexMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRegexMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRegexPatternSetResponse> deleteRegexPatternSet(
    $0.DeleteRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRegexPatternSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRuleResponse> deleteRule(
    $0.DeleteRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteRuleGroupResponse> deleteRuleGroup(
    $0.DeleteRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRuleGroup, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteSizeConstraintSetResponse>
      deleteSizeConstraintSet(
    $0.DeleteSizeConstraintSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSizeConstraintSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteSqlInjectionMatchSetResponse>
      deleteSqlInjectionMatchSet(
    $0.DeleteSqlInjectionMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSqlInjectionMatchSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteWebACLResponse> deleteWebACL(
    $0.DeleteWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteWebACL, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteXssMatchSetResponse> deleteXssMatchSet(
    $0.DeleteXssMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteXssMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetByteMatchSetResponse> getByteMatchSet(
    $0.GetByteMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getByteMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetChangeTokenResponse> getChangeToken(
    $0.GetChangeTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getChangeToken, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetChangeTokenStatusResponse> getChangeTokenStatus(
    $0.GetChangeTokenStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getChangeTokenStatus, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetGeoMatchSetResponse> getGeoMatchSet(
    $0.GetGeoMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGeoMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetIPSetResponse> getIPSet(
    $0.GetIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIPSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
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

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPermissionPolicyResponse> getPermissionPolicy(
    $0.GetPermissionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPermissionPolicy, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRateBasedRuleResponse> getRateBasedRule(
    $0.GetRateBasedRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRateBasedRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRateBasedRuleManagedKeysResponse>
      getRateBasedRuleManagedKeys(
    $0.GetRateBasedRuleManagedKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRateBasedRuleManagedKeys, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRegexMatchSetResponse> getRegexMatchSet(
    $0.GetRegexMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRegexMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRegexPatternSetResponse> getRegexPatternSet(
    $0.GetRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRegexPatternSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRuleResponse> getRule(
    $0.GetRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRuleGroupResponse> getRuleGroup(
    $0.GetRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRuleGroup, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSampledRequestsResponse> getSampledRequests(
    $0.GetSampledRequestsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSampledRequests, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSizeConstraintSetResponse> getSizeConstraintSet(
    $0.GetSizeConstraintSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSizeConstraintSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSqlInjectionMatchSetResponse>
      getSqlInjectionMatchSet(
    $0.GetSqlInjectionMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSqlInjectionMatchSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetWebACLResponse> getWebACL(
    $0.GetWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getWebACL, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetXssMatchSetResponse> getXssMatchSet(
    $0.GetXssMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getXssMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListActivatedRulesInRuleGroupResponse>
      listActivatedRulesInRuleGroup(
    $0.ListActivatedRulesInRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listActivatedRulesInRuleGroup, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListByteMatchSetsResponse> listByteMatchSets(
    $0.ListByteMatchSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listByteMatchSets, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListGeoMatchSetsResponse> listGeoMatchSets(
    $0.ListGeoMatchSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGeoMatchSets, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListIPSetsResponse> listIPSets(
    $0.ListIPSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listIPSets, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
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

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRateBasedRulesResponse> listRateBasedRules(
    $0.ListRateBasedRulesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRateBasedRules, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRegexMatchSetsResponse> listRegexMatchSets(
    $0.ListRegexMatchSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRegexMatchSets, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRegexPatternSetsResponse> listRegexPatternSets(
    $0.ListRegexPatternSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRegexPatternSets, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRuleGroupsResponse> listRuleGroups(
    $0.ListRuleGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRuleGroups, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRulesResponse> listRules(
    $0.ListRulesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRules, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSizeConstraintSetsResponse>
      listSizeConstraintSets(
    $0.ListSizeConstraintSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSizeConstraintSets, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSqlInjectionMatchSetsResponse>
      listSqlInjectionMatchSets(
    $0.ListSqlInjectionMatchSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSqlInjectionMatchSets, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSubscribedRuleGroupsResponse>
      listSubscribedRuleGroups(
    $0.ListSubscribedRuleGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSubscribedRuleGroups, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListWebACLsResponse> listWebACLs(
    $0.ListWebACLsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listWebACLs, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListXssMatchSetsResponse> listXssMatchSets(
    $0.ListXssMatchSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listXssMatchSets, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
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

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutPermissionPolicyResponse> putPermissionPolicy(
    $0.PutPermissionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putPermissionPolicy, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateByteMatchSetResponse> updateByteMatchSet(
    $0.UpdateByteMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateByteMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateGeoMatchSetResponse> updateGeoMatchSet(
    $0.UpdateGeoMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateGeoMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateIPSetResponse> updateIPSet(
    $0.UpdateIPSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateIPSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRateBasedRuleResponse> updateRateBasedRule(
    $0.UpdateRateBasedRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRateBasedRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRegexMatchSetResponse> updateRegexMatchSet(
    $0.UpdateRegexMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRegexMatchSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRegexPatternSetResponse> updateRegexPatternSet(
    $0.UpdateRegexPatternSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRegexPatternSet, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRuleResponse> updateRule(
    $0.UpdateRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRule, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateRuleGroupResponse> updateRuleGroup(
    $0.UpdateRuleGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRuleGroup, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateSizeConstraintSetResponse>
      updateSizeConstraintSet(
    $0.UpdateSizeConstraintSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSizeConstraintSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateSqlInjectionMatchSetResponse>
      updateSqlInjectionMatchSet(
    $0.UpdateSqlInjectionMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSqlInjectionMatchSet, request,
        options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateWebACLResponse> updateWebACL(
    $0.UpdateWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateWebACL, request, options: options);
  }

  /// This is AWS WAF Classic documentation. For more information, see AWS WAF Classic in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API and see the AWS WAF Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateXssMatchSetResponse> updateXssMatchSet(
    $0.UpdateXssMatchSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateXssMatchSet, request, options: options);
  }

  // method descriptors

  static final _$createByteMatchSet = $grpc.ClientMethod<
          $0.CreateByteMatchSetRequest, $0.CreateByteMatchSetResponse>(
      '/waf.WAFService/CreateByteMatchSet',
      ($0.CreateByteMatchSetRequest value) => value.writeToBuffer(),
      $0.CreateByteMatchSetResponse.fromBuffer);
  static final _$createGeoMatchSet = $grpc.ClientMethod<
          $0.CreateGeoMatchSetRequest, $0.CreateGeoMatchSetResponse>(
      '/waf.WAFService/CreateGeoMatchSet',
      ($0.CreateGeoMatchSetRequest value) => value.writeToBuffer(),
      $0.CreateGeoMatchSetResponse.fromBuffer);
  static final _$createIPSet =
      $grpc.ClientMethod<$0.CreateIPSetRequest, $0.CreateIPSetResponse>(
          '/waf.WAFService/CreateIPSet',
          ($0.CreateIPSetRequest value) => value.writeToBuffer(),
          $0.CreateIPSetResponse.fromBuffer);
  static final _$createRateBasedRule = $grpc.ClientMethod<
          $0.CreateRateBasedRuleRequest, $0.CreateRateBasedRuleResponse>(
      '/waf.WAFService/CreateRateBasedRule',
      ($0.CreateRateBasedRuleRequest value) => value.writeToBuffer(),
      $0.CreateRateBasedRuleResponse.fromBuffer);
  static final _$createRegexMatchSet = $grpc.ClientMethod<
          $0.CreateRegexMatchSetRequest, $0.CreateRegexMatchSetResponse>(
      '/waf.WAFService/CreateRegexMatchSet',
      ($0.CreateRegexMatchSetRequest value) => value.writeToBuffer(),
      $0.CreateRegexMatchSetResponse.fromBuffer);
  static final _$createRegexPatternSet = $grpc.ClientMethod<
          $0.CreateRegexPatternSetRequest, $0.CreateRegexPatternSetResponse>(
      '/waf.WAFService/CreateRegexPatternSet',
      ($0.CreateRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.CreateRegexPatternSetResponse.fromBuffer);
  static final _$createRule =
      $grpc.ClientMethod<$0.CreateRuleRequest, $0.CreateRuleResponse>(
          '/waf.WAFService/CreateRule',
          ($0.CreateRuleRequest value) => value.writeToBuffer(),
          $0.CreateRuleResponse.fromBuffer);
  static final _$createRuleGroup =
      $grpc.ClientMethod<$0.CreateRuleGroupRequest, $0.CreateRuleGroupResponse>(
          '/waf.WAFService/CreateRuleGroup',
          ($0.CreateRuleGroupRequest value) => value.writeToBuffer(),
          $0.CreateRuleGroupResponse.fromBuffer);
  static final _$createSizeConstraintSet = $grpc.ClientMethod<
          $0.CreateSizeConstraintSetRequest,
          $0.CreateSizeConstraintSetResponse>(
      '/waf.WAFService/CreateSizeConstraintSet',
      ($0.CreateSizeConstraintSetRequest value) => value.writeToBuffer(),
      $0.CreateSizeConstraintSetResponse.fromBuffer);
  static final _$createSqlInjectionMatchSet = $grpc.ClientMethod<
          $0.CreateSqlInjectionMatchSetRequest,
          $0.CreateSqlInjectionMatchSetResponse>(
      '/waf.WAFService/CreateSqlInjectionMatchSet',
      ($0.CreateSqlInjectionMatchSetRequest value) => value.writeToBuffer(),
      $0.CreateSqlInjectionMatchSetResponse.fromBuffer);
  static final _$createWebACL =
      $grpc.ClientMethod<$0.CreateWebACLRequest, $0.CreateWebACLResponse>(
          '/waf.WAFService/CreateWebACL',
          ($0.CreateWebACLRequest value) => value.writeToBuffer(),
          $0.CreateWebACLResponse.fromBuffer);
  static final _$createWebACLMigrationStack = $grpc.ClientMethod<
          $0.CreateWebACLMigrationStackRequest,
          $0.CreateWebACLMigrationStackResponse>(
      '/waf.WAFService/CreateWebACLMigrationStack',
      ($0.CreateWebACLMigrationStackRequest value) => value.writeToBuffer(),
      $0.CreateWebACLMigrationStackResponse.fromBuffer);
  static final _$createXssMatchSet = $grpc.ClientMethod<
          $0.CreateXssMatchSetRequest, $0.CreateXssMatchSetResponse>(
      '/waf.WAFService/CreateXssMatchSet',
      ($0.CreateXssMatchSetRequest value) => value.writeToBuffer(),
      $0.CreateXssMatchSetResponse.fromBuffer);
  static final _$deleteByteMatchSet = $grpc.ClientMethod<
          $0.DeleteByteMatchSetRequest, $0.DeleteByteMatchSetResponse>(
      '/waf.WAFService/DeleteByteMatchSet',
      ($0.DeleteByteMatchSetRequest value) => value.writeToBuffer(),
      $0.DeleteByteMatchSetResponse.fromBuffer);
  static final _$deleteGeoMatchSet = $grpc.ClientMethod<
          $0.DeleteGeoMatchSetRequest, $0.DeleteGeoMatchSetResponse>(
      '/waf.WAFService/DeleteGeoMatchSet',
      ($0.DeleteGeoMatchSetRequest value) => value.writeToBuffer(),
      $0.DeleteGeoMatchSetResponse.fromBuffer);
  static final _$deleteIPSet =
      $grpc.ClientMethod<$0.DeleteIPSetRequest, $0.DeleteIPSetResponse>(
          '/waf.WAFService/DeleteIPSet',
          ($0.DeleteIPSetRequest value) => value.writeToBuffer(),
          $0.DeleteIPSetResponse.fromBuffer);
  static final _$deleteLoggingConfiguration = $grpc.ClientMethod<
          $0.DeleteLoggingConfigurationRequest,
          $0.DeleteLoggingConfigurationResponse>(
      '/waf.WAFService/DeleteLoggingConfiguration',
      ($0.DeleteLoggingConfigurationRequest value) => value.writeToBuffer(),
      $0.DeleteLoggingConfigurationResponse.fromBuffer);
  static final _$deletePermissionPolicy = $grpc.ClientMethod<
          $0.DeletePermissionPolicyRequest, $0.DeletePermissionPolicyResponse>(
      '/waf.WAFService/DeletePermissionPolicy',
      ($0.DeletePermissionPolicyRequest value) => value.writeToBuffer(),
      $0.DeletePermissionPolicyResponse.fromBuffer);
  static final _$deleteRateBasedRule = $grpc.ClientMethod<
          $0.DeleteRateBasedRuleRequest, $0.DeleteRateBasedRuleResponse>(
      '/waf.WAFService/DeleteRateBasedRule',
      ($0.DeleteRateBasedRuleRequest value) => value.writeToBuffer(),
      $0.DeleteRateBasedRuleResponse.fromBuffer);
  static final _$deleteRegexMatchSet = $grpc.ClientMethod<
          $0.DeleteRegexMatchSetRequest, $0.DeleteRegexMatchSetResponse>(
      '/waf.WAFService/DeleteRegexMatchSet',
      ($0.DeleteRegexMatchSetRequest value) => value.writeToBuffer(),
      $0.DeleteRegexMatchSetResponse.fromBuffer);
  static final _$deleteRegexPatternSet = $grpc.ClientMethod<
          $0.DeleteRegexPatternSetRequest, $0.DeleteRegexPatternSetResponse>(
      '/waf.WAFService/DeleteRegexPatternSet',
      ($0.DeleteRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.DeleteRegexPatternSetResponse.fromBuffer);
  static final _$deleteRule =
      $grpc.ClientMethod<$0.DeleteRuleRequest, $0.DeleteRuleResponse>(
          '/waf.WAFService/DeleteRule',
          ($0.DeleteRuleRequest value) => value.writeToBuffer(),
          $0.DeleteRuleResponse.fromBuffer);
  static final _$deleteRuleGroup =
      $grpc.ClientMethod<$0.DeleteRuleGroupRequest, $0.DeleteRuleGroupResponse>(
          '/waf.WAFService/DeleteRuleGroup',
          ($0.DeleteRuleGroupRequest value) => value.writeToBuffer(),
          $0.DeleteRuleGroupResponse.fromBuffer);
  static final _$deleteSizeConstraintSet = $grpc.ClientMethod<
          $0.DeleteSizeConstraintSetRequest,
          $0.DeleteSizeConstraintSetResponse>(
      '/waf.WAFService/DeleteSizeConstraintSet',
      ($0.DeleteSizeConstraintSetRequest value) => value.writeToBuffer(),
      $0.DeleteSizeConstraintSetResponse.fromBuffer);
  static final _$deleteSqlInjectionMatchSet = $grpc.ClientMethod<
          $0.DeleteSqlInjectionMatchSetRequest,
          $0.DeleteSqlInjectionMatchSetResponse>(
      '/waf.WAFService/DeleteSqlInjectionMatchSet',
      ($0.DeleteSqlInjectionMatchSetRequest value) => value.writeToBuffer(),
      $0.DeleteSqlInjectionMatchSetResponse.fromBuffer);
  static final _$deleteWebACL =
      $grpc.ClientMethod<$0.DeleteWebACLRequest, $0.DeleteWebACLResponse>(
          '/waf.WAFService/DeleteWebACL',
          ($0.DeleteWebACLRequest value) => value.writeToBuffer(),
          $0.DeleteWebACLResponse.fromBuffer);
  static final _$deleteXssMatchSet = $grpc.ClientMethod<
          $0.DeleteXssMatchSetRequest, $0.DeleteXssMatchSetResponse>(
      '/waf.WAFService/DeleteXssMatchSet',
      ($0.DeleteXssMatchSetRequest value) => value.writeToBuffer(),
      $0.DeleteXssMatchSetResponse.fromBuffer);
  static final _$getByteMatchSet =
      $grpc.ClientMethod<$0.GetByteMatchSetRequest, $0.GetByteMatchSetResponse>(
          '/waf.WAFService/GetByteMatchSet',
          ($0.GetByteMatchSetRequest value) => value.writeToBuffer(),
          $0.GetByteMatchSetResponse.fromBuffer);
  static final _$getChangeToken =
      $grpc.ClientMethod<$0.GetChangeTokenRequest, $0.GetChangeTokenResponse>(
          '/waf.WAFService/GetChangeToken',
          ($0.GetChangeTokenRequest value) => value.writeToBuffer(),
          $0.GetChangeTokenResponse.fromBuffer);
  static final _$getChangeTokenStatus = $grpc.ClientMethod<
          $0.GetChangeTokenStatusRequest, $0.GetChangeTokenStatusResponse>(
      '/waf.WAFService/GetChangeTokenStatus',
      ($0.GetChangeTokenStatusRequest value) => value.writeToBuffer(),
      $0.GetChangeTokenStatusResponse.fromBuffer);
  static final _$getGeoMatchSet =
      $grpc.ClientMethod<$0.GetGeoMatchSetRequest, $0.GetGeoMatchSetResponse>(
          '/waf.WAFService/GetGeoMatchSet',
          ($0.GetGeoMatchSetRequest value) => value.writeToBuffer(),
          $0.GetGeoMatchSetResponse.fromBuffer);
  static final _$getIPSet =
      $grpc.ClientMethod<$0.GetIPSetRequest, $0.GetIPSetResponse>(
          '/waf.WAFService/GetIPSet',
          ($0.GetIPSetRequest value) => value.writeToBuffer(),
          $0.GetIPSetResponse.fromBuffer);
  static final _$getLoggingConfiguration = $grpc.ClientMethod<
          $0.GetLoggingConfigurationRequest,
          $0.GetLoggingConfigurationResponse>(
      '/waf.WAFService/GetLoggingConfiguration',
      ($0.GetLoggingConfigurationRequest value) => value.writeToBuffer(),
      $0.GetLoggingConfigurationResponse.fromBuffer);
  static final _$getPermissionPolicy = $grpc.ClientMethod<
          $0.GetPermissionPolicyRequest, $0.GetPermissionPolicyResponse>(
      '/waf.WAFService/GetPermissionPolicy',
      ($0.GetPermissionPolicyRequest value) => value.writeToBuffer(),
      $0.GetPermissionPolicyResponse.fromBuffer);
  static final _$getRateBasedRule = $grpc.ClientMethod<
          $0.GetRateBasedRuleRequest, $0.GetRateBasedRuleResponse>(
      '/waf.WAFService/GetRateBasedRule',
      ($0.GetRateBasedRuleRequest value) => value.writeToBuffer(),
      $0.GetRateBasedRuleResponse.fromBuffer);
  static final _$getRateBasedRuleManagedKeys = $grpc.ClientMethod<
          $0.GetRateBasedRuleManagedKeysRequest,
          $0.GetRateBasedRuleManagedKeysResponse>(
      '/waf.WAFService/GetRateBasedRuleManagedKeys',
      ($0.GetRateBasedRuleManagedKeysRequest value) => value.writeToBuffer(),
      $0.GetRateBasedRuleManagedKeysResponse.fromBuffer);
  static final _$getRegexMatchSet = $grpc.ClientMethod<
          $0.GetRegexMatchSetRequest, $0.GetRegexMatchSetResponse>(
      '/waf.WAFService/GetRegexMatchSet',
      ($0.GetRegexMatchSetRequest value) => value.writeToBuffer(),
      $0.GetRegexMatchSetResponse.fromBuffer);
  static final _$getRegexPatternSet = $grpc.ClientMethod<
          $0.GetRegexPatternSetRequest, $0.GetRegexPatternSetResponse>(
      '/waf.WAFService/GetRegexPatternSet',
      ($0.GetRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.GetRegexPatternSetResponse.fromBuffer);
  static final _$getRule =
      $grpc.ClientMethod<$0.GetRuleRequest, $0.GetRuleResponse>(
          '/waf.WAFService/GetRule',
          ($0.GetRuleRequest value) => value.writeToBuffer(),
          $0.GetRuleResponse.fromBuffer);
  static final _$getRuleGroup =
      $grpc.ClientMethod<$0.GetRuleGroupRequest, $0.GetRuleGroupResponse>(
          '/waf.WAFService/GetRuleGroup',
          ($0.GetRuleGroupRequest value) => value.writeToBuffer(),
          $0.GetRuleGroupResponse.fromBuffer);
  static final _$getSampledRequests = $grpc.ClientMethod<
          $0.GetSampledRequestsRequest, $0.GetSampledRequestsResponse>(
      '/waf.WAFService/GetSampledRequests',
      ($0.GetSampledRequestsRequest value) => value.writeToBuffer(),
      $0.GetSampledRequestsResponse.fromBuffer);
  static final _$getSizeConstraintSet = $grpc.ClientMethod<
          $0.GetSizeConstraintSetRequest, $0.GetSizeConstraintSetResponse>(
      '/waf.WAFService/GetSizeConstraintSet',
      ($0.GetSizeConstraintSetRequest value) => value.writeToBuffer(),
      $0.GetSizeConstraintSetResponse.fromBuffer);
  static final _$getSqlInjectionMatchSet = $grpc.ClientMethod<
          $0.GetSqlInjectionMatchSetRequest,
          $0.GetSqlInjectionMatchSetResponse>(
      '/waf.WAFService/GetSqlInjectionMatchSet',
      ($0.GetSqlInjectionMatchSetRequest value) => value.writeToBuffer(),
      $0.GetSqlInjectionMatchSetResponse.fromBuffer);
  static final _$getWebACL =
      $grpc.ClientMethod<$0.GetWebACLRequest, $0.GetWebACLResponse>(
          '/waf.WAFService/GetWebACL',
          ($0.GetWebACLRequest value) => value.writeToBuffer(),
          $0.GetWebACLResponse.fromBuffer);
  static final _$getXssMatchSet =
      $grpc.ClientMethod<$0.GetXssMatchSetRequest, $0.GetXssMatchSetResponse>(
          '/waf.WAFService/GetXssMatchSet',
          ($0.GetXssMatchSetRequest value) => value.writeToBuffer(),
          $0.GetXssMatchSetResponse.fromBuffer);
  static final _$listActivatedRulesInRuleGroup = $grpc.ClientMethod<
          $0.ListActivatedRulesInRuleGroupRequest,
          $0.ListActivatedRulesInRuleGroupResponse>(
      '/waf.WAFService/ListActivatedRulesInRuleGroup',
      ($0.ListActivatedRulesInRuleGroupRequest value) => value.writeToBuffer(),
      $0.ListActivatedRulesInRuleGroupResponse.fromBuffer);
  static final _$listByteMatchSets = $grpc.ClientMethod<
          $0.ListByteMatchSetsRequest, $0.ListByteMatchSetsResponse>(
      '/waf.WAFService/ListByteMatchSets',
      ($0.ListByteMatchSetsRequest value) => value.writeToBuffer(),
      $0.ListByteMatchSetsResponse.fromBuffer);
  static final _$listGeoMatchSets = $grpc.ClientMethod<
          $0.ListGeoMatchSetsRequest, $0.ListGeoMatchSetsResponse>(
      '/waf.WAFService/ListGeoMatchSets',
      ($0.ListGeoMatchSetsRequest value) => value.writeToBuffer(),
      $0.ListGeoMatchSetsResponse.fromBuffer);
  static final _$listIPSets =
      $grpc.ClientMethod<$0.ListIPSetsRequest, $0.ListIPSetsResponse>(
          '/waf.WAFService/ListIPSets',
          ($0.ListIPSetsRequest value) => value.writeToBuffer(),
          $0.ListIPSetsResponse.fromBuffer);
  static final _$listLoggingConfigurations = $grpc.ClientMethod<
          $0.ListLoggingConfigurationsRequest,
          $0.ListLoggingConfigurationsResponse>(
      '/waf.WAFService/ListLoggingConfigurations',
      ($0.ListLoggingConfigurationsRequest value) => value.writeToBuffer(),
      $0.ListLoggingConfigurationsResponse.fromBuffer);
  static final _$listRateBasedRules = $grpc.ClientMethod<
          $0.ListRateBasedRulesRequest, $0.ListRateBasedRulesResponse>(
      '/waf.WAFService/ListRateBasedRules',
      ($0.ListRateBasedRulesRequest value) => value.writeToBuffer(),
      $0.ListRateBasedRulesResponse.fromBuffer);
  static final _$listRegexMatchSets = $grpc.ClientMethod<
          $0.ListRegexMatchSetsRequest, $0.ListRegexMatchSetsResponse>(
      '/waf.WAFService/ListRegexMatchSets',
      ($0.ListRegexMatchSetsRequest value) => value.writeToBuffer(),
      $0.ListRegexMatchSetsResponse.fromBuffer);
  static final _$listRegexPatternSets = $grpc.ClientMethod<
          $0.ListRegexPatternSetsRequest, $0.ListRegexPatternSetsResponse>(
      '/waf.WAFService/ListRegexPatternSets',
      ($0.ListRegexPatternSetsRequest value) => value.writeToBuffer(),
      $0.ListRegexPatternSetsResponse.fromBuffer);
  static final _$listRuleGroups =
      $grpc.ClientMethod<$0.ListRuleGroupsRequest, $0.ListRuleGroupsResponse>(
          '/waf.WAFService/ListRuleGroups',
          ($0.ListRuleGroupsRequest value) => value.writeToBuffer(),
          $0.ListRuleGroupsResponse.fromBuffer);
  static final _$listRules =
      $grpc.ClientMethod<$0.ListRulesRequest, $0.ListRulesResponse>(
          '/waf.WAFService/ListRules',
          ($0.ListRulesRequest value) => value.writeToBuffer(),
          $0.ListRulesResponse.fromBuffer);
  static final _$listSizeConstraintSets = $grpc.ClientMethod<
          $0.ListSizeConstraintSetsRequest, $0.ListSizeConstraintSetsResponse>(
      '/waf.WAFService/ListSizeConstraintSets',
      ($0.ListSizeConstraintSetsRequest value) => value.writeToBuffer(),
      $0.ListSizeConstraintSetsResponse.fromBuffer);
  static final _$listSqlInjectionMatchSets = $grpc.ClientMethod<
          $0.ListSqlInjectionMatchSetsRequest,
          $0.ListSqlInjectionMatchSetsResponse>(
      '/waf.WAFService/ListSqlInjectionMatchSets',
      ($0.ListSqlInjectionMatchSetsRequest value) => value.writeToBuffer(),
      $0.ListSqlInjectionMatchSetsResponse.fromBuffer);
  static final _$listSubscribedRuleGroups = $grpc.ClientMethod<
          $0.ListSubscribedRuleGroupsRequest,
          $0.ListSubscribedRuleGroupsResponse>(
      '/waf.WAFService/ListSubscribedRuleGroups',
      ($0.ListSubscribedRuleGroupsRequest value) => value.writeToBuffer(),
      $0.ListSubscribedRuleGroupsResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/waf.WAFService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listWebACLs =
      $grpc.ClientMethod<$0.ListWebACLsRequest, $0.ListWebACLsResponse>(
          '/waf.WAFService/ListWebACLs',
          ($0.ListWebACLsRequest value) => value.writeToBuffer(),
          $0.ListWebACLsResponse.fromBuffer);
  static final _$listXssMatchSets = $grpc.ClientMethod<
          $0.ListXssMatchSetsRequest, $0.ListXssMatchSetsResponse>(
      '/waf.WAFService/ListXssMatchSets',
      ($0.ListXssMatchSetsRequest value) => value.writeToBuffer(),
      $0.ListXssMatchSetsResponse.fromBuffer);
  static final _$putLoggingConfiguration = $grpc.ClientMethod<
          $0.PutLoggingConfigurationRequest,
          $0.PutLoggingConfigurationResponse>(
      '/waf.WAFService/PutLoggingConfiguration',
      ($0.PutLoggingConfigurationRequest value) => value.writeToBuffer(),
      $0.PutLoggingConfigurationResponse.fromBuffer);
  static final _$putPermissionPolicy = $grpc.ClientMethod<
          $0.PutPermissionPolicyRequest, $0.PutPermissionPolicyResponse>(
      '/waf.WAFService/PutPermissionPolicy',
      ($0.PutPermissionPolicyRequest value) => value.writeToBuffer(),
      $0.PutPermissionPolicyResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/waf.WAFService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/waf.WAFService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateByteMatchSet = $grpc.ClientMethod<
          $0.UpdateByteMatchSetRequest, $0.UpdateByteMatchSetResponse>(
      '/waf.WAFService/UpdateByteMatchSet',
      ($0.UpdateByteMatchSetRequest value) => value.writeToBuffer(),
      $0.UpdateByteMatchSetResponse.fromBuffer);
  static final _$updateGeoMatchSet = $grpc.ClientMethod<
          $0.UpdateGeoMatchSetRequest, $0.UpdateGeoMatchSetResponse>(
      '/waf.WAFService/UpdateGeoMatchSet',
      ($0.UpdateGeoMatchSetRequest value) => value.writeToBuffer(),
      $0.UpdateGeoMatchSetResponse.fromBuffer);
  static final _$updateIPSet =
      $grpc.ClientMethod<$0.UpdateIPSetRequest, $0.UpdateIPSetResponse>(
          '/waf.WAFService/UpdateIPSet',
          ($0.UpdateIPSetRequest value) => value.writeToBuffer(),
          $0.UpdateIPSetResponse.fromBuffer);
  static final _$updateRateBasedRule = $grpc.ClientMethod<
          $0.UpdateRateBasedRuleRequest, $0.UpdateRateBasedRuleResponse>(
      '/waf.WAFService/UpdateRateBasedRule',
      ($0.UpdateRateBasedRuleRequest value) => value.writeToBuffer(),
      $0.UpdateRateBasedRuleResponse.fromBuffer);
  static final _$updateRegexMatchSet = $grpc.ClientMethod<
          $0.UpdateRegexMatchSetRequest, $0.UpdateRegexMatchSetResponse>(
      '/waf.WAFService/UpdateRegexMatchSet',
      ($0.UpdateRegexMatchSetRequest value) => value.writeToBuffer(),
      $0.UpdateRegexMatchSetResponse.fromBuffer);
  static final _$updateRegexPatternSet = $grpc.ClientMethod<
          $0.UpdateRegexPatternSetRequest, $0.UpdateRegexPatternSetResponse>(
      '/waf.WAFService/UpdateRegexPatternSet',
      ($0.UpdateRegexPatternSetRequest value) => value.writeToBuffer(),
      $0.UpdateRegexPatternSetResponse.fromBuffer);
  static final _$updateRule =
      $grpc.ClientMethod<$0.UpdateRuleRequest, $0.UpdateRuleResponse>(
          '/waf.WAFService/UpdateRule',
          ($0.UpdateRuleRequest value) => value.writeToBuffer(),
          $0.UpdateRuleResponse.fromBuffer);
  static final _$updateRuleGroup =
      $grpc.ClientMethod<$0.UpdateRuleGroupRequest, $0.UpdateRuleGroupResponse>(
          '/waf.WAFService/UpdateRuleGroup',
          ($0.UpdateRuleGroupRequest value) => value.writeToBuffer(),
          $0.UpdateRuleGroupResponse.fromBuffer);
  static final _$updateSizeConstraintSet = $grpc.ClientMethod<
          $0.UpdateSizeConstraintSetRequest,
          $0.UpdateSizeConstraintSetResponse>(
      '/waf.WAFService/UpdateSizeConstraintSet',
      ($0.UpdateSizeConstraintSetRequest value) => value.writeToBuffer(),
      $0.UpdateSizeConstraintSetResponse.fromBuffer);
  static final _$updateSqlInjectionMatchSet = $grpc.ClientMethod<
          $0.UpdateSqlInjectionMatchSetRequest,
          $0.UpdateSqlInjectionMatchSetResponse>(
      '/waf.WAFService/UpdateSqlInjectionMatchSet',
      ($0.UpdateSqlInjectionMatchSetRequest value) => value.writeToBuffer(),
      $0.UpdateSqlInjectionMatchSetResponse.fromBuffer);
  static final _$updateWebACL =
      $grpc.ClientMethod<$0.UpdateWebACLRequest, $0.UpdateWebACLResponse>(
          '/waf.WAFService/UpdateWebACL',
          ($0.UpdateWebACLRequest value) => value.writeToBuffer(),
          $0.UpdateWebACLResponse.fromBuffer);
  static final _$updateXssMatchSet = $grpc.ClientMethod<
          $0.UpdateXssMatchSetRequest, $0.UpdateXssMatchSetResponse>(
      '/waf.WAFService/UpdateXssMatchSet',
      ($0.UpdateXssMatchSetRequest value) => value.writeToBuffer(),
      $0.UpdateXssMatchSetResponse.fromBuffer);
}

@$pb.GrpcServiceName('waf.WAFService')
abstract class WAFServiceBase extends $grpc.Service {
  $core.String get $name => 'waf.WAFService';

  WAFServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateByteMatchSetRequest,
            $0.CreateByteMatchSetResponse>(
        'CreateByteMatchSet',
        createByteMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateByteMatchSetRequest.fromBuffer(value),
        ($0.CreateByteMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateGeoMatchSetRequest,
            $0.CreateGeoMatchSetResponse>(
        'CreateGeoMatchSet',
        createGeoMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateGeoMatchSetRequest.fromBuffer(value),
        ($0.CreateGeoMatchSetResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateIPSetRequest, $0.CreateIPSetResponse>(
            'CreateIPSet',
            createIPSet_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateIPSetRequest.fromBuffer(value),
            ($0.CreateIPSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRateBasedRuleRequest,
            $0.CreateRateBasedRuleResponse>(
        'CreateRateBasedRule',
        createRateBasedRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRateBasedRuleRequest.fromBuffer(value),
        ($0.CreateRateBasedRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRegexMatchSetRequest,
            $0.CreateRegexMatchSetResponse>(
        'CreateRegexMatchSet',
        createRegexMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRegexMatchSetRequest.fromBuffer(value),
        ($0.CreateRegexMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRegexPatternSetRequest,
            $0.CreateRegexPatternSetResponse>(
        'CreateRegexPatternSet',
        createRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRegexPatternSetRequest.fromBuffer(value),
        ($0.CreateRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRuleRequest, $0.CreateRuleResponse>(
        'CreateRule',
        createRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateRuleRequest.fromBuffer(value),
        ($0.CreateRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRuleGroupRequest,
            $0.CreateRuleGroupResponse>(
        'CreateRuleGroup',
        createRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRuleGroupRequest.fromBuffer(value),
        ($0.CreateRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateSizeConstraintSetRequest,
            $0.CreateSizeConstraintSetResponse>(
        'CreateSizeConstraintSet',
        createSizeConstraintSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateSizeConstraintSetRequest.fromBuffer(value),
        ($0.CreateSizeConstraintSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateSqlInjectionMatchSetRequest,
            $0.CreateSqlInjectionMatchSetResponse>(
        'CreateSqlInjectionMatchSet',
        createSqlInjectionMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateSqlInjectionMatchSetRequest.fromBuffer(value),
        ($0.CreateSqlInjectionMatchSetResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateWebACLRequest, $0.CreateWebACLResponse>(
            'CreateWebACL',
            createWebACL_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateWebACLRequest.fromBuffer(value),
            ($0.CreateWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateWebACLMigrationStackRequest,
            $0.CreateWebACLMigrationStackResponse>(
        'CreateWebACLMigrationStack',
        createWebACLMigrationStack_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateWebACLMigrationStackRequest.fromBuffer(value),
        ($0.CreateWebACLMigrationStackResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateXssMatchSetRequest,
            $0.CreateXssMatchSetResponse>(
        'CreateXssMatchSet',
        createXssMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateXssMatchSetRequest.fromBuffer(value),
        ($0.CreateXssMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteByteMatchSetRequest,
            $0.DeleteByteMatchSetResponse>(
        'DeleteByteMatchSet',
        deleteByteMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteByteMatchSetRequest.fromBuffer(value),
        ($0.DeleteByteMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteGeoMatchSetRequest,
            $0.DeleteGeoMatchSetResponse>(
        'DeleteGeoMatchSet',
        deleteGeoMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteGeoMatchSetRequest.fromBuffer(value),
        ($0.DeleteGeoMatchSetResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.DeleteRateBasedRuleRequest,
            $0.DeleteRateBasedRuleResponse>(
        'DeleteRateBasedRule',
        deleteRateBasedRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRateBasedRuleRequest.fromBuffer(value),
        ($0.DeleteRateBasedRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRegexMatchSetRequest,
            $0.DeleteRegexMatchSetResponse>(
        'DeleteRegexMatchSet',
        deleteRegexMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRegexMatchSetRequest.fromBuffer(value),
        ($0.DeleteRegexMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRegexPatternSetRequest,
            $0.DeleteRegexPatternSetResponse>(
        'DeleteRegexPatternSet',
        deleteRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRegexPatternSetRequest.fromBuffer(value),
        ($0.DeleteRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRuleRequest, $0.DeleteRuleResponse>(
        'DeleteRule',
        deleteRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteRuleRequest.fromBuffer(value),
        ($0.DeleteRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRuleGroupRequest,
            $0.DeleteRuleGroupResponse>(
        'DeleteRuleGroup',
        deleteRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRuleGroupRequest.fromBuffer(value),
        ($0.DeleteRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSizeConstraintSetRequest,
            $0.DeleteSizeConstraintSetResponse>(
        'DeleteSizeConstraintSet',
        deleteSizeConstraintSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSizeConstraintSetRequest.fromBuffer(value),
        ($0.DeleteSizeConstraintSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSqlInjectionMatchSetRequest,
            $0.DeleteSqlInjectionMatchSetResponse>(
        'DeleteSqlInjectionMatchSet',
        deleteSqlInjectionMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSqlInjectionMatchSetRequest.fromBuffer(value),
        ($0.DeleteSqlInjectionMatchSetResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteWebACLRequest, $0.DeleteWebACLResponse>(
            'DeleteWebACL',
            deleteWebACL_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteWebACLRequest.fromBuffer(value),
            ($0.DeleteWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteXssMatchSetRequest,
            $0.DeleteXssMatchSetResponse>(
        'DeleteXssMatchSet',
        deleteXssMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteXssMatchSetRequest.fromBuffer(value),
        ($0.DeleteXssMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetByteMatchSetRequest,
            $0.GetByteMatchSetResponse>(
        'GetByteMatchSet',
        getByteMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetByteMatchSetRequest.fromBuffer(value),
        ($0.GetByteMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChangeTokenRequest,
            $0.GetChangeTokenResponse>(
        'GetChangeToken',
        getChangeToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetChangeTokenRequest.fromBuffer(value),
        ($0.GetChangeTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChangeTokenStatusRequest,
            $0.GetChangeTokenStatusResponse>(
        'GetChangeTokenStatus',
        getChangeTokenStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetChangeTokenStatusRequest.fromBuffer(value),
        ($0.GetChangeTokenStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetGeoMatchSetRequest,
            $0.GetGeoMatchSetResponse>(
        'GetGeoMatchSet',
        getGeoMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetGeoMatchSetRequest.fromBuffer(value),
        ($0.GetGeoMatchSetResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.GetPermissionPolicyRequest,
            $0.GetPermissionPolicyResponse>(
        'GetPermissionPolicy',
        getPermissionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPermissionPolicyRequest.fromBuffer(value),
        ($0.GetPermissionPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRateBasedRuleRequest,
            $0.GetRateBasedRuleResponse>(
        'GetRateBasedRule',
        getRateBasedRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRateBasedRuleRequest.fromBuffer(value),
        ($0.GetRateBasedRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRateBasedRuleManagedKeysRequest,
            $0.GetRateBasedRuleManagedKeysResponse>(
        'GetRateBasedRuleManagedKeys',
        getRateBasedRuleManagedKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRateBasedRuleManagedKeysRequest.fromBuffer(value),
        ($0.GetRateBasedRuleManagedKeysResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRegexMatchSetRequest,
            $0.GetRegexMatchSetResponse>(
        'GetRegexMatchSet',
        getRegexMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRegexMatchSetRequest.fromBuffer(value),
        ($0.GetRegexMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRegexPatternSetRequest,
            $0.GetRegexPatternSetResponse>(
        'GetRegexPatternSet',
        getRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRegexPatternSetRequest.fromBuffer(value),
        ($0.GetRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRuleRequest, $0.GetRuleResponse>(
        'GetRule',
        getRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetRuleRequest.fromBuffer(value),
        ($0.GetRuleResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.GetSizeConstraintSetRequest,
            $0.GetSizeConstraintSetResponse>(
        'GetSizeConstraintSet',
        getSizeConstraintSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSizeConstraintSetRequest.fromBuffer(value),
        ($0.GetSizeConstraintSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSqlInjectionMatchSetRequest,
            $0.GetSqlInjectionMatchSetResponse>(
        'GetSqlInjectionMatchSet',
        getSqlInjectionMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSqlInjectionMatchSetRequest.fromBuffer(value),
        ($0.GetSqlInjectionMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetWebACLRequest, $0.GetWebACLResponse>(
        'GetWebACL',
        getWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetWebACLRequest.fromBuffer(value),
        ($0.GetWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetXssMatchSetRequest,
            $0.GetXssMatchSetResponse>(
        'GetXssMatchSet',
        getXssMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetXssMatchSetRequest.fromBuffer(value),
        ($0.GetXssMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListActivatedRulesInRuleGroupRequest,
            $0.ListActivatedRulesInRuleGroupResponse>(
        'ListActivatedRulesInRuleGroup',
        listActivatedRulesInRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListActivatedRulesInRuleGroupRequest.fromBuffer(value),
        ($0.ListActivatedRulesInRuleGroupResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListByteMatchSetsRequest,
            $0.ListByteMatchSetsResponse>(
        'ListByteMatchSets',
        listByteMatchSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListByteMatchSetsRequest.fromBuffer(value),
        ($0.ListByteMatchSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGeoMatchSetsRequest,
            $0.ListGeoMatchSetsResponse>(
        'ListGeoMatchSets',
        listGeoMatchSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListGeoMatchSetsRequest.fromBuffer(value),
        ($0.ListGeoMatchSetsResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.ListRateBasedRulesRequest,
            $0.ListRateBasedRulesResponse>(
        'ListRateBasedRules',
        listRateBasedRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRateBasedRulesRequest.fromBuffer(value),
        ($0.ListRateBasedRulesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRegexMatchSetsRequest,
            $0.ListRegexMatchSetsResponse>(
        'ListRegexMatchSets',
        listRegexMatchSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRegexMatchSetsRequest.fromBuffer(value),
        ($0.ListRegexMatchSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRegexPatternSetsRequest,
            $0.ListRegexPatternSetsResponse>(
        'ListRegexPatternSets',
        listRegexPatternSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRegexPatternSetsRequest.fromBuffer(value),
        ($0.ListRegexPatternSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRuleGroupsRequest,
            $0.ListRuleGroupsResponse>(
        'ListRuleGroups',
        listRuleGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRuleGroupsRequest.fromBuffer(value),
        ($0.ListRuleGroupsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRulesRequest, $0.ListRulesResponse>(
        'ListRules',
        listRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListRulesRequest.fromBuffer(value),
        ($0.ListRulesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSizeConstraintSetsRequest,
            $0.ListSizeConstraintSetsResponse>(
        'ListSizeConstraintSets',
        listSizeConstraintSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSizeConstraintSetsRequest.fromBuffer(value),
        ($0.ListSizeConstraintSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSqlInjectionMatchSetsRequest,
            $0.ListSqlInjectionMatchSetsResponse>(
        'ListSqlInjectionMatchSets',
        listSqlInjectionMatchSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSqlInjectionMatchSetsRequest.fromBuffer(value),
        ($0.ListSqlInjectionMatchSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSubscribedRuleGroupsRequest,
            $0.ListSubscribedRuleGroupsResponse>(
        'ListSubscribedRuleGroups',
        listSubscribedRuleGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSubscribedRuleGroupsRequest.fromBuffer(value),
        ($0.ListSubscribedRuleGroupsResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.ListXssMatchSetsRequest,
            $0.ListXssMatchSetsResponse>(
        'ListXssMatchSets',
        listXssMatchSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListXssMatchSetsRequest.fromBuffer(value),
        ($0.ListXssMatchSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutLoggingConfigurationRequest,
            $0.PutLoggingConfigurationResponse>(
        'PutLoggingConfiguration',
        putLoggingConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutLoggingConfigurationRequest.fromBuffer(value),
        ($0.PutLoggingConfigurationResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.UpdateByteMatchSetRequest,
            $0.UpdateByteMatchSetResponse>(
        'UpdateByteMatchSet',
        updateByteMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateByteMatchSetRequest.fromBuffer(value),
        ($0.UpdateByteMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateGeoMatchSetRequest,
            $0.UpdateGeoMatchSetResponse>(
        'UpdateGeoMatchSet',
        updateGeoMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateGeoMatchSetRequest.fromBuffer(value),
        ($0.UpdateGeoMatchSetResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateIPSetRequest, $0.UpdateIPSetResponse>(
            'UpdateIPSet',
            updateIPSet_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateIPSetRequest.fromBuffer(value),
            ($0.UpdateIPSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRateBasedRuleRequest,
            $0.UpdateRateBasedRuleResponse>(
        'UpdateRateBasedRule',
        updateRateBasedRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRateBasedRuleRequest.fromBuffer(value),
        ($0.UpdateRateBasedRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRegexMatchSetRequest,
            $0.UpdateRegexMatchSetResponse>(
        'UpdateRegexMatchSet',
        updateRegexMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRegexMatchSetRequest.fromBuffer(value),
        ($0.UpdateRegexMatchSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRegexPatternSetRequest,
            $0.UpdateRegexPatternSetResponse>(
        'UpdateRegexPatternSet',
        updateRegexPatternSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRegexPatternSetRequest.fromBuffer(value),
        ($0.UpdateRegexPatternSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRuleRequest, $0.UpdateRuleResponse>(
        'UpdateRule',
        updateRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateRuleRequest.fromBuffer(value),
        ($0.UpdateRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRuleGroupRequest,
            $0.UpdateRuleGroupResponse>(
        'UpdateRuleGroup',
        updateRuleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRuleGroupRequest.fromBuffer(value),
        ($0.UpdateRuleGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateSizeConstraintSetRequest,
            $0.UpdateSizeConstraintSetResponse>(
        'UpdateSizeConstraintSet',
        updateSizeConstraintSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateSizeConstraintSetRequest.fromBuffer(value),
        ($0.UpdateSizeConstraintSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateSqlInjectionMatchSetRequest,
            $0.UpdateSqlInjectionMatchSetResponse>(
        'UpdateSqlInjectionMatchSet',
        updateSqlInjectionMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateSqlInjectionMatchSetRequest.fromBuffer(value),
        ($0.UpdateSqlInjectionMatchSetResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateWebACLRequest, $0.UpdateWebACLResponse>(
            'UpdateWebACL',
            updateWebACL_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateWebACLRequest.fromBuffer(value),
            ($0.UpdateWebACLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateXssMatchSetRequest,
            $0.UpdateXssMatchSetResponse>(
        'UpdateXssMatchSet',
        updateXssMatchSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateXssMatchSetRequest.fromBuffer(value),
        ($0.UpdateXssMatchSetResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.CreateByteMatchSetResponse> createByteMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateByteMatchSetRequest> $request) async {
    return createByteMatchSet($call, await $request);
  }

  $async.Future<$0.CreateByteMatchSetResponse> createByteMatchSet(
      $grpc.ServiceCall call, $0.CreateByteMatchSetRequest request);

  $async.Future<$0.CreateGeoMatchSetResponse> createGeoMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateGeoMatchSetRequest> $request) async {
    return createGeoMatchSet($call, await $request);
  }

  $async.Future<$0.CreateGeoMatchSetResponse> createGeoMatchSet(
      $grpc.ServiceCall call, $0.CreateGeoMatchSetRequest request);

  $async.Future<$0.CreateIPSetResponse> createIPSet_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateIPSetRequest> $request) async {
    return createIPSet($call, await $request);
  }

  $async.Future<$0.CreateIPSetResponse> createIPSet(
      $grpc.ServiceCall call, $0.CreateIPSetRequest request);

  $async.Future<$0.CreateRateBasedRuleResponse> createRateBasedRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRateBasedRuleRequest> $request) async {
    return createRateBasedRule($call, await $request);
  }

  $async.Future<$0.CreateRateBasedRuleResponse> createRateBasedRule(
      $grpc.ServiceCall call, $0.CreateRateBasedRuleRequest request);

  $async.Future<$0.CreateRegexMatchSetResponse> createRegexMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRegexMatchSetRequest> $request) async {
    return createRegexMatchSet($call, await $request);
  }

  $async.Future<$0.CreateRegexMatchSetResponse> createRegexMatchSet(
      $grpc.ServiceCall call, $0.CreateRegexMatchSetRequest request);

  $async.Future<$0.CreateRegexPatternSetResponse> createRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRegexPatternSetRequest> $request) async {
    return createRegexPatternSet($call, await $request);
  }

  $async.Future<$0.CreateRegexPatternSetResponse> createRegexPatternSet(
      $grpc.ServiceCall call, $0.CreateRegexPatternSetRequest request);

  $async.Future<$0.CreateRuleResponse> createRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateRuleRequest> $request) async {
    return createRule($call, await $request);
  }

  $async.Future<$0.CreateRuleResponse> createRule(
      $grpc.ServiceCall call, $0.CreateRuleRequest request);

  $async.Future<$0.CreateRuleGroupResponse> createRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRuleGroupRequest> $request) async {
    return createRuleGroup($call, await $request);
  }

  $async.Future<$0.CreateRuleGroupResponse> createRuleGroup(
      $grpc.ServiceCall call, $0.CreateRuleGroupRequest request);

  $async.Future<$0.CreateSizeConstraintSetResponse> createSizeConstraintSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateSizeConstraintSetRequest> $request) async {
    return createSizeConstraintSet($call, await $request);
  }

  $async.Future<$0.CreateSizeConstraintSetResponse> createSizeConstraintSet(
      $grpc.ServiceCall call, $0.CreateSizeConstraintSetRequest request);

  $async.Future<$0.CreateSqlInjectionMatchSetResponse>
      createSqlInjectionMatchSet_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateSqlInjectionMatchSetRequest> $request) async {
    return createSqlInjectionMatchSet($call, await $request);
  }

  $async.Future<$0.CreateSqlInjectionMatchSetResponse>
      createSqlInjectionMatchSet(
          $grpc.ServiceCall call, $0.CreateSqlInjectionMatchSetRequest request);

  $async.Future<$0.CreateWebACLResponse> createWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateWebACLRequest> $request) async {
    return createWebACL($call, await $request);
  }

  $async.Future<$0.CreateWebACLResponse> createWebACL(
      $grpc.ServiceCall call, $0.CreateWebACLRequest request);

  $async.Future<$0.CreateWebACLMigrationStackResponse>
      createWebACLMigrationStack_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateWebACLMigrationStackRequest> $request) async {
    return createWebACLMigrationStack($call, await $request);
  }

  $async.Future<$0.CreateWebACLMigrationStackResponse>
      createWebACLMigrationStack(
          $grpc.ServiceCall call, $0.CreateWebACLMigrationStackRequest request);

  $async.Future<$0.CreateXssMatchSetResponse> createXssMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateXssMatchSetRequest> $request) async {
    return createXssMatchSet($call, await $request);
  }

  $async.Future<$0.CreateXssMatchSetResponse> createXssMatchSet(
      $grpc.ServiceCall call, $0.CreateXssMatchSetRequest request);

  $async.Future<$0.DeleteByteMatchSetResponse> deleteByteMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteByteMatchSetRequest> $request) async {
    return deleteByteMatchSet($call, await $request);
  }

  $async.Future<$0.DeleteByteMatchSetResponse> deleteByteMatchSet(
      $grpc.ServiceCall call, $0.DeleteByteMatchSetRequest request);

  $async.Future<$0.DeleteGeoMatchSetResponse> deleteGeoMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteGeoMatchSetRequest> $request) async {
    return deleteGeoMatchSet($call, await $request);
  }

  $async.Future<$0.DeleteGeoMatchSetResponse> deleteGeoMatchSet(
      $grpc.ServiceCall call, $0.DeleteGeoMatchSetRequest request);

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

  $async.Future<$0.DeleteRateBasedRuleResponse> deleteRateBasedRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRateBasedRuleRequest> $request) async {
    return deleteRateBasedRule($call, await $request);
  }

  $async.Future<$0.DeleteRateBasedRuleResponse> deleteRateBasedRule(
      $grpc.ServiceCall call, $0.DeleteRateBasedRuleRequest request);

  $async.Future<$0.DeleteRegexMatchSetResponse> deleteRegexMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRegexMatchSetRequest> $request) async {
    return deleteRegexMatchSet($call, await $request);
  }

  $async.Future<$0.DeleteRegexMatchSetResponse> deleteRegexMatchSet(
      $grpc.ServiceCall call, $0.DeleteRegexMatchSetRequest request);

  $async.Future<$0.DeleteRegexPatternSetResponse> deleteRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRegexPatternSetRequest> $request) async {
    return deleteRegexPatternSet($call, await $request);
  }

  $async.Future<$0.DeleteRegexPatternSetResponse> deleteRegexPatternSet(
      $grpc.ServiceCall call, $0.DeleteRegexPatternSetRequest request);

  $async.Future<$0.DeleteRuleResponse> deleteRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRuleRequest> $request) async {
    return deleteRule($call, await $request);
  }

  $async.Future<$0.DeleteRuleResponse> deleteRule(
      $grpc.ServiceCall call, $0.DeleteRuleRequest request);

  $async.Future<$0.DeleteRuleGroupResponse> deleteRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRuleGroupRequest> $request) async {
    return deleteRuleGroup($call, await $request);
  }

  $async.Future<$0.DeleteRuleGroupResponse> deleteRuleGroup(
      $grpc.ServiceCall call, $0.DeleteRuleGroupRequest request);

  $async.Future<$0.DeleteSizeConstraintSetResponse> deleteSizeConstraintSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteSizeConstraintSetRequest> $request) async {
    return deleteSizeConstraintSet($call, await $request);
  }

  $async.Future<$0.DeleteSizeConstraintSetResponse> deleteSizeConstraintSet(
      $grpc.ServiceCall call, $0.DeleteSizeConstraintSetRequest request);

  $async.Future<$0.DeleteSqlInjectionMatchSetResponse>
      deleteSqlInjectionMatchSet_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteSqlInjectionMatchSetRequest> $request) async {
    return deleteSqlInjectionMatchSet($call, await $request);
  }

  $async.Future<$0.DeleteSqlInjectionMatchSetResponse>
      deleteSqlInjectionMatchSet(
          $grpc.ServiceCall call, $0.DeleteSqlInjectionMatchSetRequest request);

  $async.Future<$0.DeleteWebACLResponse> deleteWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteWebACLRequest> $request) async {
    return deleteWebACL($call, await $request);
  }

  $async.Future<$0.DeleteWebACLResponse> deleteWebACL(
      $grpc.ServiceCall call, $0.DeleteWebACLRequest request);

  $async.Future<$0.DeleteXssMatchSetResponse> deleteXssMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteXssMatchSetRequest> $request) async {
    return deleteXssMatchSet($call, await $request);
  }

  $async.Future<$0.DeleteXssMatchSetResponse> deleteXssMatchSet(
      $grpc.ServiceCall call, $0.DeleteXssMatchSetRequest request);

  $async.Future<$0.GetByteMatchSetResponse> getByteMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetByteMatchSetRequest> $request) async {
    return getByteMatchSet($call, await $request);
  }

  $async.Future<$0.GetByteMatchSetResponse> getByteMatchSet(
      $grpc.ServiceCall call, $0.GetByteMatchSetRequest request);

  $async.Future<$0.GetChangeTokenResponse> getChangeToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetChangeTokenRequest> $request) async {
    return getChangeToken($call, await $request);
  }

  $async.Future<$0.GetChangeTokenResponse> getChangeToken(
      $grpc.ServiceCall call, $0.GetChangeTokenRequest request);

  $async.Future<$0.GetChangeTokenStatusResponse> getChangeTokenStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetChangeTokenStatusRequest> $request) async {
    return getChangeTokenStatus($call, await $request);
  }

  $async.Future<$0.GetChangeTokenStatusResponse> getChangeTokenStatus(
      $grpc.ServiceCall call, $0.GetChangeTokenStatusRequest request);

  $async.Future<$0.GetGeoMatchSetResponse> getGeoMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetGeoMatchSetRequest> $request) async {
    return getGeoMatchSet($call, await $request);
  }

  $async.Future<$0.GetGeoMatchSetResponse> getGeoMatchSet(
      $grpc.ServiceCall call, $0.GetGeoMatchSetRequest request);

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

  $async.Future<$0.GetPermissionPolicyResponse> getPermissionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPermissionPolicyRequest> $request) async {
    return getPermissionPolicy($call, await $request);
  }

  $async.Future<$0.GetPermissionPolicyResponse> getPermissionPolicy(
      $grpc.ServiceCall call, $0.GetPermissionPolicyRequest request);

  $async.Future<$0.GetRateBasedRuleResponse> getRateBasedRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRateBasedRuleRequest> $request) async {
    return getRateBasedRule($call, await $request);
  }

  $async.Future<$0.GetRateBasedRuleResponse> getRateBasedRule(
      $grpc.ServiceCall call, $0.GetRateBasedRuleRequest request);

  $async.Future<$0.GetRateBasedRuleManagedKeysResponse>
      getRateBasedRuleManagedKeys_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetRateBasedRuleManagedKeysRequest> $request) async {
    return getRateBasedRuleManagedKeys($call, await $request);
  }

  $async.Future<$0.GetRateBasedRuleManagedKeysResponse>
      getRateBasedRuleManagedKeys($grpc.ServiceCall call,
          $0.GetRateBasedRuleManagedKeysRequest request);

  $async.Future<$0.GetRegexMatchSetResponse> getRegexMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRegexMatchSetRequest> $request) async {
    return getRegexMatchSet($call, await $request);
  }

  $async.Future<$0.GetRegexMatchSetResponse> getRegexMatchSet(
      $grpc.ServiceCall call, $0.GetRegexMatchSetRequest request);

  $async.Future<$0.GetRegexPatternSetResponse> getRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRegexPatternSetRequest> $request) async {
    return getRegexPatternSet($call, await $request);
  }

  $async.Future<$0.GetRegexPatternSetResponse> getRegexPatternSet(
      $grpc.ServiceCall call, $0.GetRegexPatternSetRequest request);

  $async.Future<$0.GetRuleResponse> getRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRuleRequest> $request) async {
    return getRule($call, await $request);
  }

  $async.Future<$0.GetRuleResponse> getRule(
      $grpc.ServiceCall call, $0.GetRuleRequest request);

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

  $async.Future<$0.GetSizeConstraintSetResponse> getSizeConstraintSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSizeConstraintSetRequest> $request) async {
    return getSizeConstraintSet($call, await $request);
  }

  $async.Future<$0.GetSizeConstraintSetResponse> getSizeConstraintSet(
      $grpc.ServiceCall call, $0.GetSizeConstraintSetRequest request);

  $async.Future<$0.GetSqlInjectionMatchSetResponse> getSqlInjectionMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSqlInjectionMatchSetRequest> $request) async {
    return getSqlInjectionMatchSet($call, await $request);
  }

  $async.Future<$0.GetSqlInjectionMatchSetResponse> getSqlInjectionMatchSet(
      $grpc.ServiceCall call, $0.GetSqlInjectionMatchSetRequest request);

  $async.Future<$0.GetWebACLResponse> getWebACL_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetWebACLRequest> $request) async {
    return getWebACL($call, await $request);
  }

  $async.Future<$0.GetWebACLResponse> getWebACL(
      $grpc.ServiceCall call, $0.GetWebACLRequest request);

  $async.Future<$0.GetXssMatchSetResponse> getXssMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetXssMatchSetRequest> $request) async {
    return getXssMatchSet($call, await $request);
  }

  $async.Future<$0.GetXssMatchSetResponse> getXssMatchSet(
      $grpc.ServiceCall call, $0.GetXssMatchSetRequest request);

  $async.Future<$0.ListActivatedRulesInRuleGroupResponse>
      listActivatedRulesInRuleGroup_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListActivatedRulesInRuleGroupRequest>
              $request) async {
    return listActivatedRulesInRuleGroup($call, await $request);
  }

  $async.Future<$0.ListActivatedRulesInRuleGroupResponse>
      listActivatedRulesInRuleGroup($grpc.ServiceCall call,
          $0.ListActivatedRulesInRuleGroupRequest request);

  $async.Future<$0.ListByteMatchSetsResponse> listByteMatchSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListByteMatchSetsRequest> $request) async {
    return listByteMatchSets($call, await $request);
  }

  $async.Future<$0.ListByteMatchSetsResponse> listByteMatchSets(
      $grpc.ServiceCall call, $0.ListByteMatchSetsRequest request);

  $async.Future<$0.ListGeoMatchSetsResponse> listGeoMatchSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListGeoMatchSetsRequest> $request) async {
    return listGeoMatchSets($call, await $request);
  }

  $async.Future<$0.ListGeoMatchSetsResponse> listGeoMatchSets(
      $grpc.ServiceCall call, $0.ListGeoMatchSetsRequest request);

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

  $async.Future<$0.ListRateBasedRulesResponse> listRateBasedRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRateBasedRulesRequest> $request) async {
    return listRateBasedRules($call, await $request);
  }

  $async.Future<$0.ListRateBasedRulesResponse> listRateBasedRules(
      $grpc.ServiceCall call, $0.ListRateBasedRulesRequest request);

  $async.Future<$0.ListRegexMatchSetsResponse> listRegexMatchSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRegexMatchSetsRequest> $request) async {
    return listRegexMatchSets($call, await $request);
  }

  $async.Future<$0.ListRegexMatchSetsResponse> listRegexMatchSets(
      $grpc.ServiceCall call, $0.ListRegexMatchSetsRequest request);

  $async.Future<$0.ListRegexPatternSetsResponse> listRegexPatternSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRegexPatternSetsRequest> $request) async {
    return listRegexPatternSets($call, await $request);
  }

  $async.Future<$0.ListRegexPatternSetsResponse> listRegexPatternSets(
      $grpc.ServiceCall call, $0.ListRegexPatternSetsRequest request);

  $async.Future<$0.ListRuleGroupsResponse> listRuleGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRuleGroupsRequest> $request) async {
    return listRuleGroups($call, await $request);
  }

  $async.Future<$0.ListRuleGroupsResponse> listRuleGroups(
      $grpc.ServiceCall call, $0.ListRuleGroupsRequest request);

  $async.Future<$0.ListRulesResponse> listRules_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListRulesRequest> $request) async {
    return listRules($call, await $request);
  }

  $async.Future<$0.ListRulesResponse> listRules(
      $grpc.ServiceCall call, $0.ListRulesRequest request);

  $async.Future<$0.ListSizeConstraintSetsResponse> listSizeConstraintSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSizeConstraintSetsRequest> $request) async {
    return listSizeConstraintSets($call, await $request);
  }

  $async.Future<$0.ListSizeConstraintSetsResponse> listSizeConstraintSets(
      $grpc.ServiceCall call, $0.ListSizeConstraintSetsRequest request);

  $async.Future<$0.ListSqlInjectionMatchSetsResponse>
      listSqlInjectionMatchSets_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListSqlInjectionMatchSetsRequest> $request) async {
    return listSqlInjectionMatchSets($call, await $request);
  }

  $async.Future<$0.ListSqlInjectionMatchSetsResponse> listSqlInjectionMatchSets(
      $grpc.ServiceCall call, $0.ListSqlInjectionMatchSetsRequest request);

  $async.Future<$0.ListSubscribedRuleGroupsResponse>
      listSubscribedRuleGroups_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListSubscribedRuleGroupsRequest> $request) async {
    return listSubscribedRuleGroups($call, await $request);
  }

  $async.Future<$0.ListSubscribedRuleGroupsResponse> listSubscribedRuleGroups(
      $grpc.ServiceCall call, $0.ListSubscribedRuleGroupsRequest request);

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

  $async.Future<$0.ListXssMatchSetsResponse> listXssMatchSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListXssMatchSetsRequest> $request) async {
    return listXssMatchSets($call, await $request);
  }

  $async.Future<$0.ListXssMatchSetsResponse> listXssMatchSets(
      $grpc.ServiceCall call, $0.ListXssMatchSetsRequest request);

  $async.Future<$0.PutLoggingConfigurationResponse> putLoggingConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutLoggingConfigurationRequest> $request) async {
    return putLoggingConfiguration($call, await $request);
  }

  $async.Future<$0.PutLoggingConfigurationResponse> putLoggingConfiguration(
      $grpc.ServiceCall call, $0.PutLoggingConfigurationRequest request);

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

  $async.Future<$0.UpdateByteMatchSetResponse> updateByteMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateByteMatchSetRequest> $request) async {
    return updateByteMatchSet($call, await $request);
  }

  $async.Future<$0.UpdateByteMatchSetResponse> updateByteMatchSet(
      $grpc.ServiceCall call, $0.UpdateByteMatchSetRequest request);

  $async.Future<$0.UpdateGeoMatchSetResponse> updateGeoMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateGeoMatchSetRequest> $request) async {
    return updateGeoMatchSet($call, await $request);
  }

  $async.Future<$0.UpdateGeoMatchSetResponse> updateGeoMatchSet(
      $grpc.ServiceCall call, $0.UpdateGeoMatchSetRequest request);

  $async.Future<$0.UpdateIPSetResponse> updateIPSet_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateIPSetRequest> $request) async {
    return updateIPSet($call, await $request);
  }

  $async.Future<$0.UpdateIPSetResponse> updateIPSet(
      $grpc.ServiceCall call, $0.UpdateIPSetRequest request);

  $async.Future<$0.UpdateRateBasedRuleResponse> updateRateBasedRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRateBasedRuleRequest> $request) async {
    return updateRateBasedRule($call, await $request);
  }

  $async.Future<$0.UpdateRateBasedRuleResponse> updateRateBasedRule(
      $grpc.ServiceCall call, $0.UpdateRateBasedRuleRequest request);

  $async.Future<$0.UpdateRegexMatchSetResponse> updateRegexMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRegexMatchSetRequest> $request) async {
    return updateRegexMatchSet($call, await $request);
  }

  $async.Future<$0.UpdateRegexMatchSetResponse> updateRegexMatchSet(
      $grpc.ServiceCall call, $0.UpdateRegexMatchSetRequest request);

  $async.Future<$0.UpdateRegexPatternSetResponse> updateRegexPatternSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRegexPatternSetRequest> $request) async {
    return updateRegexPatternSet($call, await $request);
  }

  $async.Future<$0.UpdateRegexPatternSetResponse> updateRegexPatternSet(
      $grpc.ServiceCall call, $0.UpdateRegexPatternSetRequest request);

  $async.Future<$0.UpdateRuleResponse> updateRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateRuleRequest> $request) async {
    return updateRule($call, await $request);
  }

  $async.Future<$0.UpdateRuleResponse> updateRule(
      $grpc.ServiceCall call, $0.UpdateRuleRequest request);

  $async.Future<$0.UpdateRuleGroupResponse> updateRuleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRuleGroupRequest> $request) async {
    return updateRuleGroup($call, await $request);
  }

  $async.Future<$0.UpdateRuleGroupResponse> updateRuleGroup(
      $grpc.ServiceCall call, $0.UpdateRuleGroupRequest request);

  $async.Future<$0.UpdateSizeConstraintSetResponse> updateSizeConstraintSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateSizeConstraintSetRequest> $request) async {
    return updateSizeConstraintSet($call, await $request);
  }

  $async.Future<$0.UpdateSizeConstraintSetResponse> updateSizeConstraintSet(
      $grpc.ServiceCall call, $0.UpdateSizeConstraintSetRequest request);

  $async.Future<$0.UpdateSqlInjectionMatchSetResponse>
      updateSqlInjectionMatchSet_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateSqlInjectionMatchSetRequest> $request) async {
    return updateSqlInjectionMatchSet($call, await $request);
  }

  $async.Future<$0.UpdateSqlInjectionMatchSetResponse>
      updateSqlInjectionMatchSet(
          $grpc.ServiceCall call, $0.UpdateSqlInjectionMatchSetRequest request);

  $async.Future<$0.UpdateWebACLResponse> updateWebACL_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateWebACLRequest> $request) async {
    return updateWebACL($call, await $request);
  }

  $async.Future<$0.UpdateWebACLResponse> updateWebACL(
      $grpc.ServiceCall call, $0.UpdateWebACLRequest request);

  $async.Future<$0.UpdateXssMatchSetResponse> updateXssMatchSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateXssMatchSetRequest> $request) async {
    return updateXssMatchSet($call, await $request);
  }

  $async.Future<$0.UpdateXssMatchSetResponse> updateXssMatchSet(
      $grpc.ServiceCall call, $0.UpdateXssMatchSetRequest request);
}
