// This is a generated file - do not edit.
//
// Generated from cognito_identity.proto.

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

import 'cognito_identity.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'cognito_identity.pb.dart';

/// CognitoIdentityService provides cognito_identity API operations.
@$pb.GrpcServiceName('cognito_identity.CognitoIdentityService')
class CognitoIdentityServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CognitoIdentityServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Creates a new identity pool. The identity pool is a store of user identity information that is specific to your Amazon Web Services account. The keys for SupportedLoginProviders are as follows: Fac...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.IdentityPool> createIdentityPool(
    $0.CreateIdentityPoolInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createIdentityPool, request, options: options);
  }

  /// Deletes identities from an identity pool. You can specify a list of 1-60 identities that you want to delete. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteIdentitiesResponse> deleteIdentities(
    $0.DeleteIdentitiesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIdentities, request, options: options);
  }

  /// Deletes an identity pool. Once a pool is deleted, users will not be able to authenticate with the pool. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteIdentityPool(
    $0.DeleteIdentityPoolInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIdentityPool, request, options: options);
  }

  /// Returns metadata related to the given identity, including when the identity was created and any associated linked logins. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.IdentityDescription> describeIdentity(
    $0.DescribeIdentityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeIdentity, request, options: options);
  }

  /// Gets details about a particular identity pool, including the pool name, ID description, creation date, and current number of users. You must use Amazon Web Services developer credentials to call th...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.IdentityPool> describeIdentityPool(
    $0.DescribeIdentityPoolInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeIdentityPool, request, options: options);
  }

  /// Returns credentials for the provided identity ID. Any provided logins will be validated against supported login providers. If the token is for cognito-identity.amazonaws.com, it will be passed thro...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCredentialsForIdentityResponse>
      getCredentialsForIdentity(
    $0.GetCredentialsForIdentityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCredentialsForIdentity, request,
        options: options);
  }

  /// Generates (or retrieves) IdentityID. Supplying multiple logins will create an implicit linked account. This is a public API. You do not need any credentials to call this API.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetIdResponse> getId(
    $0.GetIdInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getId, request, options: options);
  }

  /// Gets the roles for an identity pool. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetIdentityPoolRolesResponse> getIdentityPoolRoles(
    $0.GetIdentityPoolRolesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIdentityPoolRoles, request, options: options);
  }

  /// Gets an OpenID token, using a known Cognito ID. This known Cognito ID is returned by GetId. You can optionally add additional logins for the identity. Supplying multiple logins creates an implicit ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetOpenIdTokenResponse> getOpenIdToken(
    $0.GetOpenIdTokenInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOpenIdToken, request, options: options);
  }

  /// Registers (or retrieves) a Cognito IdentityId and an OpenID Connect token for a user authenticated by your backend authentication process. Supplying multiple logins will create an implicit linked a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetOpenIdTokenForDeveloperIdentityResponse>
      getOpenIdTokenForDeveloperIdentity(
    $0.GetOpenIdTokenForDeveloperIdentityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOpenIdTokenForDeveloperIdentity, request,
        options: options);
  }

  /// Use GetPrincipalTagAttributeMap to list all mappings between PrincipalTags and user attributes.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPrincipalTagAttributeMapResponse>
      getPrincipalTagAttributeMap(
    $0.GetPrincipalTagAttributeMapInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPrincipalTagAttributeMap, request,
        options: options);
  }

  /// Lists the identities in an identity pool. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListIdentitiesResponse> listIdentities(
    $0.ListIdentitiesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listIdentities, request, options: options);
  }

  /// Lists all of the Cognito identity pools registered for your account. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListIdentityPoolsResponse> listIdentityPools(
    $0.ListIdentityPoolsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listIdentityPools, request, options: options);
  }

  /// Lists the tags that are assigned to an Amazon Cognito identity pool. A tag is a label that you can apply to identity pools to categorize and manage them in different ways, such as by purpose, owner...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Retrieves the IdentityID associated with a DeveloperUserIdentifier or the list of DeveloperUserIdentifier values associated with an IdentityId for an existing identity. Either IdentityID or Develop...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.LookupDeveloperIdentityResponse>
      lookupDeveloperIdentity(
    $0.LookupDeveloperIdentityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$lookupDeveloperIdentity, request,
        options: options);
  }

  /// Merges two users having different IdentityIds, existing in the same identity pool, and identified by the same developer provider. You can use this action to request that discrete users be merged an...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.MergeDeveloperIdentitiesResponse>
      mergeDeveloperIdentities(
    $0.MergeDeveloperIdentitiesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$mergeDeveloperIdentities, request,
        options: options);
  }

  /// Sets the roles for an identity pool. These roles are used when making calls to GetCredentialsForIdentity action. You must use Amazon Web Services developer credentials to call this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> setIdentityPoolRoles(
    $0.SetIdentityPoolRolesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setIdentityPoolRoles, request, options: options);
  }

  /// You can use this operation to use default (username and clientID) attribute or custom attribute mappings.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetPrincipalTagAttributeMapResponse>
      setPrincipalTagAttributeMap(
    $0.SetPrincipalTagAttributeMapInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setPrincipalTagAttributeMap, request,
        options: options);
  }

  /// Assigns a set of tags to the specified Amazon Cognito identity pool. A tag is a label that you can use to categorize and manage identity pools in different ways, such as by purpose, owner, environm...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Unlinks a DeveloperUserIdentifier from an existing identity. Unlinked developer users will be considered new identities next time they are seen. If, for a given Cognito identity, you remove all fed...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> unlinkDeveloperIdentity(
    $0.UnlinkDeveloperIdentityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$unlinkDeveloperIdentity, request,
        options: options);
  }

  /// Unlinks a federated identity from an existing account. Unlinked logins will be considered new identities next time they are seen. Removing the last linked login will make this identity inaccessible...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> unlinkIdentity(
    $0.UnlinkIdentityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$unlinkIdentity, request, options: options);
  }

  /// Removes the specified tags from the specified Amazon Cognito identity pool. You can use this action up to 5 times per second, per account
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates the configuration of an identity pool. If you don't provide a value for a parameter, Amazon Cognito sets it to its default value. You must use Amazon Web Services developer credentials to c...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.IdentityPool> updateIdentityPool(
    $0.IdentityPool request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateIdentityPool, request, options: options);
  }

  // method descriptors

  static final _$createIdentityPool =
      $grpc.ClientMethod<$0.CreateIdentityPoolInput, $0.IdentityPool>(
          '/cognito_identity.CognitoIdentityService/CreateIdentityPool',
          ($0.CreateIdentityPoolInput value) => value.writeToBuffer(),
          $0.IdentityPool.fromBuffer);
  static final _$deleteIdentities =
      $grpc.ClientMethod<$0.DeleteIdentitiesInput, $0.DeleteIdentitiesResponse>(
          '/cognito_identity.CognitoIdentityService/DeleteIdentities',
          ($0.DeleteIdentitiesInput value) => value.writeToBuffer(),
          $0.DeleteIdentitiesResponse.fromBuffer);
  static final _$deleteIdentityPool =
      $grpc.ClientMethod<$0.DeleteIdentityPoolInput, $1.Empty>(
          '/cognito_identity.CognitoIdentityService/DeleteIdentityPool',
          ($0.DeleteIdentityPoolInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeIdentity =
      $grpc.ClientMethod<$0.DescribeIdentityInput, $0.IdentityDescription>(
          '/cognito_identity.CognitoIdentityService/DescribeIdentity',
          ($0.DescribeIdentityInput value) => value.writeToBuffer(),
          $0.IdentityDescription.fromBuffer);
  static final _$describeIdentityPool =
      $grpc.ClientMethod<$0.DescribeIdentityPoolInput, $0.IdentityPool>(
          '/cognito_identity.CognitoIdentityService/DescribeIdentityPool',
          ($0.DescribeIdentityPoolInput value) => value.writeToBuffer(),
          $0.IdentityPool.fromBuffer);
  static final _$getCredentialsForIdentity = $grpc.ClientMethod<
          $0.GetCredentialsForIdentityInput,
          $0.GetCredentialsForIdentityResponse>(
      '/cognito_identity.CognitoIdentityService/GetCredentialsForIdentity',
      ($0.GetCredentialsForIdentityInput value) => value.writeToBuffer(),
      $0.GetCredentialsForIdentityResponse.fromBuffer);
  static final _$getId = $grpc.ClientMethod<$0.GetIdInput, $0.GetIdResponse>(
      '/cognito_identity.CognitoIdentityService/GetId',
      ($0.GetIdInput value) => value.writeToBuffer(),
      $0.GetIdResponse.fromBuffer);
  static final _$getIdentityPoolRoles = $grpc.ClientMethod<
          $0.GetIdentityPoolRolesInput, $0.GetIdentityPoolRolesResponse>(
      '/cognito_identity.CognitoIdentityService/GetIdentityPoolRoles',
      ($0.GetIdentityPoolRolesInput value) => value.writeToBuffer(),
      $0.GetIdentityPoolRolesResponse.fromBuffer);
  static final _$getOpenIdToken =
      $grpc.ClientMethod<$0.GetOpenIdTokenInput, $0.GetOpenIdTokenResponse>(
          '/cognito_identity.CognitoIdentityService/GetOpenIdToken',
          ($0.GetOpenIdTokenInput value) => value.writeToBuffer(),
          $0.GetOpenIdTokenResponse.fromBuffer);
  static final _$getOpenIdTokenForDeveloperIdentity = $grpc.ClientMethod<
          $0.GetOpenIdTokenForDeveloperIdentityInput,
          $0.GetOpenIdTokenForDeveloperIdentityResponse>(
      '/cognito_identity.CognitoIdentityService/GetOpenIdTokenForDeveloperIdentity',
      ($0.GetOpenIdTokenForDeveloperIdentityInput value) =>
          value.writeToBuffer(),
      $0.GetOpenIdTokenForDeveloperIdentityResponse.fromBuffer);
  static final _$getPrincipalTagAttributeMap = $grpc.ClientMethod<
          $0.GetPrincipalTagAttributeMapInput,
          $0.GetPrincipalTagAttributeMapResponse>(
      '/cognito_identity.CognitoIdentityService/GetPrincipalTagAttributeMap',
      ($0.GetPrincipalTagAttributeMapInput value) => value.writeToBuffer(),
      $0.GetPrincipalTagAttributeMapResponse.fromBuffer);
  static final _$listIdentities =
      $grpc.ClientMethod<$0.ListIdentitiesInput, $0.ListIdentitiesResponse>(
          '/cognito_identity.CognitoIdentityService/ListIdentities',
          ($0.ListIdentitiesInput value) => value.writeToBuffer(),
          $0.ListIdentitiesResponse.fromBuffer);
  static final _$listIdentityPools = $grpc.ClientMethod<
          $0.ListIdentityPoolsInput, $0.ListIdentityPoolsResponse>(
      '/cognito_identity.CognitoIdentityService/ListIdentityPools',
      ($0.ListIdentityPoolsInput value) => value.writeToBuffer(),
      $0.ListIdentityPoolsResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceInput, $0.ListTagsForResourceResponse>(
      '/cognito_identity.CognitoIdentityService/ListTagsForResource',
      ($0.ListTagsForResourceInput value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$lookupDeveloperIdentity = $grpc.ClientMethod<
          $0.LookupDeveloperIdentityInput, $0.LookupDeveloperIdentityResponse>(
      '/cognito_identity.CognitoIdentityService/LookupDeveloperIdentity',
      ($0.LookupDeveloperIdentityInput value) => value.writeToBuffer(),
      $0.LookupDeveloperIdentityResponse.fromBuffer);
  static final _$mergeDeveloperIdentities = $grpc.ClientMethod<
          $0.MergeDeveloperIdentitiesInput,
          $0.MergeDeveloperIdentitiesResponse>(
      '/cognito_identity.CognitoIdentityService/MergeDeveloperIdentities',
      ($0.MergeDeveloperIdentitiesInput value) => value.writeToBuffer(),
      $0.MergeDeveloperIdentitiesResponse.fromBuffer);
  static final _$setIdentityPoolRoles =
      $grpc.ClientMethod<$0.SetIdentityPoolRolesInput, $1.Empty>(
          '/cognito_identity.CognitoIdentityService/SetIdentityPoolRoles',
          ($0.SetIdentityPoolRolesInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setPrincipalTagAttributeMap = $grpc.ClientMethod<
          $0.SetPrincipalTagAttributeMapInput,
          $0.SetPrincipalTagAttributeMapResponse>(
      '/cognito_identity.CognitoIdentityService/SetPrincipalTagAttributeMap',
      ($0.SetPrincipalTagAttributeMapInput value) => value.writeToBuffer(),
      $0.SetPrincipalTagAttributeMapResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $0.TagResourceResponse>(
          '/cognito_identity.CognitoIdentityService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$unlinkDeveloperIdentity =
      $grpc.ClientMethod<$0.UnlinkDeveloperIdentityInput, $1.Empty>(
          '/cognito_identity.CognitoIdentityService/UnlinkDeveloperIdentity',
          ($0.UnlinkDeveloperIdentityInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$unlinkIdentity =
      $grpc.ClientMethod<$0.UnlinkIdentityInput, $1.Empty>(
          '/cognito_identity.CognitoIdentityService/UnlinkIdentity',
          ($0.UnlinkIdentityInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $0.UntagResourceResponse>(
          '/cognito_identity.CognitoIdentityService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateIdentityPool =
      $grpc.ClientMethod<$0.IdentityPool, $0.IdentityPool>(
          '/cognito_identity.CognitoIdentityService/UpdateIdentityPool',
          ($0.IdentityPool value) => value.writeToBuffer(),
          $0.IdentityPool.fromBuffer);
}

@$pb.GrpcServiceName('cognito_identity.CognitoIdentityService')
abstract class CognitoIdentityServiceBase extends $grpc.Service {
  $core.String get $name => 'cognito_identity.CognitoIdentityService';

  CognitoIdentityServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateIdentityPoolInput, $0.IdentityPool>(
        'CreateIdentityPool',
        createIdentityPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateIdentityPoolInput.fromBuffer(value),
        ($0.IdentityPool value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteIdentitiesInput,
            $0.DeleteIdentitiesResponse>(
        'DeleteIdentities',
        deleteIdentities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteIdentitiesInput.fromBuffer(value),
        ($0.DeleteIdentitiesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteIdentityPoolInput, $1.Empty>(
        'DeleteIdentityPool',
        deleteIdentityPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteIdentityPoolInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeIdentityInput, $0.IdentityDescription>(
            'DescribeIdentity',
            describeIdentity_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeIdentityInput.fromBuffer(value),
            ($0.IdentityDescription value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeIdentityPoolInput, $0.IdentityPool>(
            'DescribeIdentityPool',
            describeIdentityPool_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeIdentityPoolInput.fromBuffer(value),
            ($0.IdentityPool value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCredentialsForIdentityInput,
            $0.GetCredentialsForIdentityResponse>(
        'GetCredentialsForIdentity',
        getCredentialsForIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCredentialsForIdentityInput.fromBuffer(value),
        ($0.GetCredentialsForIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIdInput, $0.GetIdResponse>(
        'GetId',
        getId_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetIdInput.fromBuffer(value),
        ($0.GetIdResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIdentityPoolRolesInput,
            $0.GetIdentityPoolRolesResponse>(
        'GetIdentityPoolRoles',
        getIdentityPoolRoles_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetIdentityPoolRolesInput.fromBuffer(value),
        ($0.GetIdentityPoolRolesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetOpenIdTokenInput, $0.GetOpenIdTokenResponse>(
            'GetOpenIdToken',
            getOpenIdToken_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetOpenIdTokenInput.fromBuffer(value),
            ($0.GetOpenIdTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOpenIdTokenForDeveloperIdentityInput,
            $0.GetOpenIdTokenForDeveloperIdentityResponse>(
        'GetOpenIdTokenForDeveloperIdentity',
        getOpenIdTokenForDeveloperIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOpenIdTokenForDeveloperIdentityInput.fromBuffer(value),
        ($0.GetOpenIdTokenForDeveloperIdentityResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPrincipalTagAttributeMapInput,
            $0.GetPrincipalTagAttributeMapResponse>(
        'GetPrincipalTagAttributeMap',
        getPrincipalTagAttributeMap_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPrincipalTagAttributeMapInput.fromBuffer(value),
        ($0.GetPrincipalTagAttributeMapResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListIdentitiesInput, $0.ListIdentitiesResponse>(
            'ListIdentities',
            listIdentities_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListIdentitiesInput.fromBuffer(value),
            ($0.ListIdentitiesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListIdentityPoolsInput,
            $0.ListIdentityPoolsResponse>(
        'ListIdentityPools',
        listIdentityPools_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListIdentityPoolsInput.fromBuffer(value),
        ($0.ListIdentityPoolsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceInput,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceInput.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LookupDeveloperIdentityInput,
            $0.LookupDeveloperIdentityResponse>(
        'LookupDeveloperIdentity',
        lookupDeveloperIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.LookupDeveloperIdentityInput.fromBuffer(value),
        ($0.LookupDeveloperIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.MergeDeveloperIdentitiesInput,
            $0.MergeDeveloperIdentitiesResponse>(
        'MergeDeveloperIdentities',
        mergeDeveloperIdentities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.MergeDeveloperIdentitiesInput.fromBuffer(value),
        ($0.MergeDeveloperIdentitiesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetIdentityPoolRolesInput, $1.Empty>(
        'SetIdentityPoolRoles',
        setIdentityPoolRoles_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetIdentityPoolRolesInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetPrincipalTagAttributeMapInput,
            $0.SetPrincipalTagAttributeMapResponse>(
        'SetPrincipalTagAttributeMap',
        setPrincipalTagAttributeMap_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetPrincipalTagAttributeMapInput.fromBuffer(value),
        ($0.SetPrincipalTagAttributeMapResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $0.TagResourceResponse>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UnlinkDeveloperIdentityInput, $1.Empty>(
        'UnlinkDeveloperIdentity',
        unlinkDeveloperIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UnlinkDeveloperIdentityInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UnlinkIdentityInput, $1.Empty>(
        'UnlinkIdentity',
        unlinkIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UnlinkIdentityInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceInput, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceInput.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.IdentityPool, $0.IdentityPool>(
        'UpdateIdentityPool',
        updateIdentityPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.IdentityPool.fromBuffer(value),
        ($0.IdentityPool value) => value.writeToBuffer()));
  }

  $async.Future<$0.IdentityPool> createIdentityPool_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateIdentityPoolInput> $request) async {
    return createIdentityPool($call, await $request);
  }

  $async.Future<$0.IdentityPool> createIdentityPool(
      $grpc.ServiceCall call, $0.CreateIdentityPoolInput request);

  $async.Future<$0.DeleteIdentitiesResponse> deleteIdentities_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteIdentitiesInput> $request) async {
    return deleteIdentities($call, await $request);
  }

  $async.Future<$0.DeleteIdentitiesResponse> deleteIdentities(
      $grpc.ServiceCall call, $0.DeleteIdentitiesInput request);

  $async.Future<$1.Empty> deleteIdentityPool_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteIdentityPoolInput> $request) async {
    return deleteIdentityPool($call, await $request);
  }

  $async.Future<$1.Empty> deleteIdentityPool(
      $grpc.ServiceCall call, $0.DeleteIdentityPoolInput request);

  $async.Future<$0.IdentityDescription> describeIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeIdentityInput> $request) async {
    return describeIdentity($call, await $request);
  }

  $async.Future<$0.IdentityDescription> describeIdentity(
      $grpc.ServiceCall call, $0.DescribeIdentityInput request);

  $async.Future<$0.IdentityPool> describeIdentityPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeIdentityPoolInput> $request) async {
    return describeIdentityPool($call, await $request);
  }

  $async.Future<$0.IdentityPool> describeIdentityPool(
      $grpc.ServiceCall call, $0.DescribeIdentityPoolInput request);

  $async.Future<$0.GetCredentialsForIdentityResponse>
      getCredentialsForIdentity_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetCredentialsForIdentityInput> $request) async {
    return getCredentialsForIdentity($call, await $request);
  }

  $async.Future<$0.GetCredentialsForIdentityResponse> getCredentialsForIdentity(
      $grpc.ServiceCall call, $0.GetCredentialsForIdentityInput request);

  $async.Future<$0.GetIdResponse> getId_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.GetIdInput> $request) async {
    return getId($call, await $request);
  }

  $async.Future<$0.GetIdResponse> getId(
      $grpc.ServiceCall call, $0.GetIdInput request);

  $async.Future<$0.GetIdentityPoolRolesResponse> getIdentityPoolRoles_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetIdentityPoolRolesInput> $request) async {
    return getIdentityPoolRoles($call, await $request);
  }

  $async.Future<$0.GetIdentityPoolRolesResponse> getIdentityPoolRoles(
      $grpc.ServiceCall call, $0.GetIdentityPoolRolesInput request);

  $async.Future<$0.GetOpenIdTokenResponse> getOpenIdToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetOpenIdTokenInput> $request) async {
    return getOpenIdToken($call, await $request);
  }

  $async.Future<$0.GetOpenIdTokenResponse> getOpenIdToken(
      $grpc.ServiceCall call, $0.GetOpenIdTokenInput request);

  $async.Future<$0.GetOpenIdTokenForDeveloperIdentityResponse>
      getOpenIdTokenForDeveloperIdentity_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetOpenIdTokenForDeveloperIdentityInput>
              $request) async {
    return getOpenIdTokenForDeveloperIdentity($call, await $request);
  }

  $async.Future<$0.GetOpenIdTokenForDeveloperIdentityResponse>
      getOpenIdTokenForDeveloperIdentity($grpc.ServiceCall call,
          $0.GetOpenIdTokenForDeveloperIdentityInput request);

  $async.Future<$0.GetPrincipalTagAttributeMapResponse>
      getPrincipalTagAttributeMap_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetPrincipalTagAttributeMapInput> $request) async {
    return getPrincipalTagAttributeMap($call, await $request);
  }

  $async.Future<$0.GetPrincipalTagAttributeMapResponse>
      getPrincipalTagAttributeMap(
          $grpc.ServiceCall call, $0.GetPrincipalTagAttributeMapInput request);

  $async.Future<$0.ListIdentitiesResponse> listIdentities_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListIdentitiesInput> $request) async {
    return listIdentities($call, await $request);
  }

  $async.Future<$0.ListIdentitiesResponse> listIdentities(
      $grpc.ServiceCall call, $0.ListIdentitiesInput request);

  $async.Future<$0.ListIdentityPoolsResponse> listIdentityPools_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListIdentityPoolsInput> $request) async {
    return listIdentityPools($call, await $request);
  }

  $async.Future<$0.ListIdentityPoolsResponse> listIdentityPools(
      $grpc.ServiceCall call, $0.ListIdentityPoolsInput request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceInput> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceInput request);

  $async.Future<$0.LookupDeveloperIdentityResponse> lookupDeveloperIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.LookupDeveloperIdentityInput> $request) async {
    return lookupDeveloperIdentity($call, await $request);
  }

  $async.Future<$0.LookupDeveloperIdentityResponse> lookupDeveloperIdentity(
      $grpc.ServiceCall call, $0.LookupDeveloperIdentityInput request);

  $async.Future<$0.MergeDeveloperIdentitiesResponse>
      mergeDeveloperIdentities_Pre($grpc.ServiceCall $call,
          $async.Future<$0.MergeDeveloperIdentitiesInput> $request) async {
    return mergeDeveloperIdentities($call, await $request);
  }

  $async.Future<$0.MergeDeveloperIdentitiesResponse> mergeDeveloperIdentities(
      $grpc.ServiceCall call, $0.MergeDeveloperIdentitiesInput request);

  $async.Future<$1.Empty> setIdentityPoolRoles_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetIdentityPoolRolesInput> $request) async {
    return setIdentityPoolRoles($call, await $request);
  }

  $async.Future<$1.Empty> setIdentityPoolRoles(
      $grpc.ServiceCall call, $0.SetIdentityPoolRolesInput request);

  $async.Future<$0.SetPrincipalTagAttributeMapResponse>
      setPrincipalTagAttributeMap_Pre($grpc.ServiceCall $call,
          $async.Future<$0.SetPrincipalTagAttributeMapInput> $request) async {
    return setPrincipalTagAttributeMap($call, await $request);
  }

  $async.Future<$0.SetPrincipalTagAttributeMapResponse>
      setPrincipalTagAttributeMap(
          $grpc.ServiceCall call, $0.SetPrincipalTagAttributeMapInput request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$1.Empty> unlinkDeveloperIdentity_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UnlinkDeveloperIdentityInput> $request) async {
    return unlinkDeveloperIdentity($call, await $request);
  }

  $async.Future<$1.Empty> unlinkDeveloperIdentity(
      $grpc.ServiceCall call, $0.UnlinkDeveloperIdentityInput request);

  $async.Future<$1.Empty> unlinkIdentity_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UnlinkIdentityInput> $request) async {
    return unlinkIdentity($call, await $request);
  }

  $async.Future<$1.Empty> unlinkIdentity(
      $grpc.ServiceCall call, $0.UnlinkIdentityInput request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);

  $async.Future<$0.IdentityPool> updateIdentityPool_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.IdentityPool> $request) async {
    return updateIdentityPool($call, await $request);
  }

  $async.Future<$0.IdentityPool> updateIdentityPool(
      $grpc.ServiceCall call, $0.IdentityPool request);
}
