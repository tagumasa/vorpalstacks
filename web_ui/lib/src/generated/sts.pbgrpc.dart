// This is a generated file - do not edit.
//
// Generated from sts.proto.

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

import 'sts.pb.dart' as $0;

export 'sts.pb.dart';

/// STSService provides sts API operations.
@$pb.GrpcServiceName('sts.STSService')
class STSServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  STSServiceClient(super.channel, {super.options, super.interceptors});

  /// Returns a set of temporary security credentials that you can use to access Amazon Web Services resources. These temporary credentials consist of an access key ID, a secret access key, and a securit...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.AssumeRoleResponse> assumeRole(
    $0.AssumeRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$assumeRole, request, options: options);
  }

  /// Returns a set of temporary security credentials for users who have been authenticated via a SAML authentication response. This operation provides a mechanism for tying an enterprise identity store ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.AssumeRoleWithSAMLResponse> assumeRoleWithSAML(
    $0.AssumeRoleWithSAMLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$assumeRoleWithSAML, request, options: options);
  }

  /// Returns a set of temporary security credentials for users who have been authenticated in a mobile or web application with a web identity provider. Example providers include the OAuth 2.0 providers ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.AssumeRoleWithWebIdentityResponse>
      assumeRoleWithWebIdentity(
    $0.AssumeRoleWithWebIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$assumeRoleWithWebIdentity, request,
        options: options);
  }

  /// Returns a set of short term credentials you can use to perform privileged tasks on a member account in your organization. You must use credentials from an Organizations management account or a dele...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.AssumeRootResponse> assumeRoot(
    $0.AssumeRootRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$assumeRoot, request, options: options);
  }

  /// Decodes additional information about the authorization status of a request from an encoded message returned in response to an Amazon Web Services request. For example, if a user is not authorized t...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.DecodeAuthorizationMessageResponse>
      decodeAuthorizationMessage(
    $0.DecodeAuthorizationMessageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$decodeAuthorizationMessage, request,
        options: options);
  }

  /// Returns the account identifier for the specified access key ID. Access keys consist of two parts: an access key ID (for example, AKIAIOSFODNN7EXAMPLE) and a secret access key (for example, wJalrXUt...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetAccessKeyInfoResponse> getAccessKeyInfo(
    $0.GetAccessKeyInfoRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccessKeyInfo, request, options: options);
  }

  /// Returns details about the IAM user or role whose credentials are used to call the operation. No permissions are required to perform this operation. If an administrator attaches a policy to your ide...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetCallerIdentityResponse> getCallerIdentity(
    $0.GetCallerIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCallerIdentity, request, options: options);
  }

  /// Exchanges a trade-in token for temporary Amazon Web Services credentials with the permissions associated with the assumed principal. This operation allows you to obtain credentials for a specific p...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetDelegatedAccessTokenResponse>
      getDelegatedAccessToken(
    $0.GetDelegatedAccessTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDelegatedAccessToken, request,
        options: options);
  }

  /// Returns a set of temporary security credentials (consisting of an access key ID, a secret access key, and a security token) for a user. A typical use is in a proxy application that gets temporary s...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetFederationTokenResponse> getFederationToken(
    $0.GetFederationTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFederationToken, request, options: options);
  }

  /// Returns a set of temporary credentials for an Amazon Web Services account or IAM user. The credentials consist of an access key ID, a secret access key, and a security token. Typically, you use Get...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetSessionTokenResponse> getSessionToken(
    $0.GetSessionTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSessionToken, request, options: options);
  }

  /// Returns a signed JSON Web Token (JWT) that represents the calling Amazon Web Services identity. The returned JWT can be used to authenticate with external services that support OIDC discovery. The ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetWebIdentityTokenResponse> getWebIdentityToken(
    $0.GetWebIdentityTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getWebIdentityToken, request, options: options);
  }

  // method descriptors

  static final _$assumeRole =
      $grpc.ClientMethod<$0.AssumeRoleRequest, $0.AssumeRoleResponse>(
          '/sts.STSService/AssumeRole',
          ($0.AssumeRoleRequest value) => value.writeToBuffer(),
          $0.AssumeRoleResponse.fromBuffer);
  static final _$assumeRoleWithSAML = $grpc.ClientMethod<
          $0.AssumeRoleWithSAMLRequest, $0.AssumeRoleWithSAMLResponse>(
      '/sts.STSService/AssumeRoleWithSAML',
      ($0.AssumeRoleWithSAMLRequest value) => value.writeToBuffer(),
      $0.AssumeRoleWithSAMLResponse.fromBuffer);
  static final _$assumeRoleWithWebIdentity = $grpc.ClientMethod<
          $0.AssumeRoleWithWebIdentityRequest,
          $0.AssumeRoleWithWebIdentityResponse>(
      '/sts.STSService/AssumeRoleWithWebIdentity',
      ($0.AssumeRoleWithWebIdentityRequest value) => value.writeToBuffer(),
      $0.AssumeRoleWithWebIdentityResponse.fromBuffer);
  static final _$assumeRoot =
      $grpc.ClientMethod<$0.AssumeRootRequest, $0.AssumeRootResponse>(
          '/sts.STSService/AssumeRoot',
          ($0.AssumeRootRequest value) => value.writeToBuffer(),
          $0.AssumeRootResponse.fromBuffer);
  static final _$decodeAuthorizationMessage = $grpc.ClientMethod<
          $0.DecodeAuthorizationMessageRequest,
          $0.DecodeAuthorizationMessageResponse>(
      '/sts.STSService/DecodeAuthorizationMessage',
      ($0.DecodeAuthorizationMessageRequest value) => value.writeToBuffer(),
      $0.DecodeAuthorizationMessageResponse.fromBuffer);
  static final _$getAccessKeyInfo = $grpc.ClientMethod<
          $0.GetAccessKeyInfoRequest, $0.GetAccessKeyInfoResponse>(
      '/sts.STSService/GetAccessKeyInfo',
      ($0.GetAccessKeyInfoRequest value) => value.writeToBuffer(),
      $0.GetAccessKeyInfoResponse.fromBuffer);
  static final _$getCallerIdentity = $grpc.ClientMethod<
          $0.GetCallerIdentityRequest, $0.GetCallerIdentityResponse>(
      '/sts.STSService/GetCallerIdentity',
      ($0.GetCallerIdentityRequest value) => value.writeToBuffer(),
      $0.GetCallerIdentityResponse.fromBuffer);
  static final _$getDelegatedAccessToken = $grpc.ClientMethod<
          $0.GetDelegatedAccessTokenRequest,
          $0.GetDelegatedAccessTokenResponse>(
      '/sts.STSService/GetDelegatedAccessToken',
      ($0.GetDelegatedAccessTokenRequest value) => value.writeToBuffer(),
      $0.GetDelegatedAccessTokenResponse.fromBuffer);
  static final _$getFederationToken = $grpc.ClientMethod<
          $0.GetFederationTokenRequest, $0.GetFederationTokenResponse>(
      '/sts.STSService/GetFederationToken',
      ($0.GetFederationTokenRequest value) => value.writeToBuffer(),
      $0.GetFederationTokenResponse.fromBuffer);
  static final _$getSessionToken =
      $grpc.ClientMethod<$0.GetSessionTokenRequest, $0.GetSessionTokenResponse>(
          '/sts.STSService/GetSessionToken',
          ($0.GetSessionTokenRequest value) => value.writeToBuffer(),
          $0.GetSessionTokenResponse.fromBuffer);
  static final _$getWebIdentityToken = $grpc.ClientMethod<
          $0.GetWebIdentityTokenRequest, $0.GetWebIdentityTokenResponse>(
      '/sts.STSService/GetWebIdentityToken',
      ($0.GetWebIdentityTokenRequest value) => value.writeToBuffer(),
      $0.GetWebIdentityTokenResponse.fromBuffer);
}

@$pb.GrpcServiceName('sts.STSService')
abstract class STSServiceBase extends $grpc.Service {
  $core.String get $name => 'sts.STSService';

  STSServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AssumeRoleRequest, $0.AssumeRoleResponse>(
        'AssumeRole',
        assumeRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.AssumeRoleRequest.fromBuffer(value),
        ($0.AssumeRoleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssumeRoleWithSAMLRequest,
            $0.AssumeRoleWithSAMLResponse>(
        'AssumeRoleWithSAML',
        assumeRoleWithSAML_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssumeRoleWithSAMLRequest.fromBuffer(value),
        ($0.AssumeRoleWithSAMLResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssumeRoleWithWebIdentityRequest,
            $0.AssumeRoleWithWebIdentityResponse>(
        'AssumeRoleWithWebIdentity',
        assumeRoleWithWebIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssumeRoleWithWebIdentityRequest.fromBuffer(value),
        ($0.AssumeRoleWithWebIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssumeRootRequest, $0.AssumeRootResponse>(
        'AssumeRoot',
        assumeRoot_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.AssumeRootRequest.fromBuffer(value),
        ($0.AssumeRootResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DecodeAuthorizationMessageRequest,
            $0.DecodeAuthorizationMessageResponse>(
        'DecodeAuthorizationMessage',
        decodeAuthorizationMessage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DecodeAuthorizationMessageRequest.fromBuffer(value),
        ($0.DecodeAuthorizationMessageResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccessKeyInfoRequest,
            $0.GetAccessKeyInfoResponse>(
        'GetAccessKeyInfo',
        getAccessKeyInfo_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAccessKeyInfoRequest.fromBuffer(value),
        ($0.GetAccessKeyInfoResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCallerIdentityRequest,
            $0.GetCallerIdentityResponse>(
        'GetCallerIdentity',
        getCallerIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCallerIdentityRequest.fromBuffer(value),
        ($0.GetCallerIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDelegatedAccessTokenRequest,
            $0.GetDelegatedAccessTokenResponse>(
        'GetDelegatedAccessToken',
        getDelegatedAccessToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDelegatedAccessTokenRequest.fromBuffer(value),
        ($0.GetDelegatedAccessTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFederationTokenRequest,
            $0.GetFederationTokenResponse>(
        'GetFederationToken',
        getFederationToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFederationTokenRequest.fromBuffer(value),
        ($0.GetFederationTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSessionTokenRequest,
            $0.GetSessionTokenResponse>(
        'GetSessionToken',
        getSessionToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSessionTokenRequest.fromBuffer(value),
        ($0.GetSessionTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetWebIdentityTokenRequest,
            $0.GetWebIdentityTokenResponse>(
        'GetWebIdentityToken',
        getWebIdentityToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetWebIdentityTokenRequest.fromBuffer(value),
        ($0.GetWebIdentityTokenResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AssumeRoleResponse> assumeRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AssumeRoleRequest> $request) async {
    return assumeRole($call, await $request);
  }

  $async.Future<$0.AssumeRoleResponse> assumeRole(
      $grpc.ServiceCall call, $0.AssumeRoleRequest request);

  $async.Future<$0.AssumeRoleWithSAMLResponse> assumeRoleWithSAML_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AssumeRoleWithSAMLRequest> $request) async {
    return assumeRoleWithSAML($call, await $request);
  }

  $async.Future<$0.AssumeRoleWithSAMLResponse> assumeRoleWithSAML(
      $grpc.ServiceCall call, $0.AssumeRoleWithSAMLRequest request);

  $async.Future<$0.AssumeRoleWithWebIdentityResponse>
      assumeRoleWithWebIdentity_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AssumeRoleWithWebIdentityRequest> $request) async {
    return assumeRoleWithWebIdentity($call, await $request);
  }

  $async.Future<$0.AssumeRoleWithWebIdentityResponse> assumeRoleWithWebIdentity(
      $grpc.ServiceCall call, $0.AssumeRoleWithWebIdentityRequest request);

  $async.Future<$0.AssumeRootResponse> assumeRoot_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AssumeRootRequest> $request) async {
    return assumeRoot($call, await $request);
  }

  $async.Future<$0.AssumeRootResponse> assumeRoot(
      $grpc.ServiceCall call, $0.AssumeRootRequest request);

  $async.Future<$0.DecodeAuthorizationMessageResponse>
      decodeAuthorizationMessage_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DecodeAuthorizationMessageRequest> $request) async {
    return decodeAuthorizationMessage($call, await $request);
  }

  $async.Future<$0.DecodeAuthorizationMessageResponse>
      decodeAuthorizationMessage(
          $grpc.ServiceCall call, $0.DecodeAuthorizationMessageRequest request);

  $async.Future<$0.GetAccessKeyInfoResponse> getAccessKeyInfo_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAccessKeyInfoRequest> $request) async {
    return getAccessKeyInfo($call, await $request);
  }

  $async.Future<$0.GetAccessKeyInfoResponse> getAccessKeyInfo(
      $grpc.ServiceCall call, $0.GetAccessKeyInfoRequest request);

  $async.Future<$0.GetCallerIdentityResponse> getCallerIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCallerIdentityRequest> $request) async {
    return getCallerIdentity($call, await $request);
  }

  $async.Future<$0.GetCallerIdentityResponse> getCallerIdentity(
      $grpc.ServiceCall call, $0.GetCallerIdentityRequest request);

  $async.Future<$0.GetDelegatedAccessTokenResponse> getDelegatedAccessToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDelegatedAccessTokenRequest> $request) async {
    return getDelegatedAccessToken($call, await $request);
  }

  $async.Future<$0.GetDelegatedAccessTokenResponse> getDelegatedAccessToken(
      $grpc.ServiceCall call, $0.GetDelegatedAccessTokenRequest request);

  $async.Future<$0.GetFederationTokenResponse> getFederationToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetFederationTokenRequest> $request) async {
    return getFederationToken($call, await $request);
  }

  $async.Future<$0.GetFederationTokenResponse> getFederationToken(
      $grpc.ServiceCall call, $0.GetFederationTokenRequest request);

  $async.Future<$0.GetSessionTokenResponse> getSessionToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSessionTokenRequest> $request) async {
    return getSessionToken($call, await $request);
  }

  $async.Future<$0.GetSessionTokenResponse> getSessionToken(
      $grpc.ServiceCall call, $0.GetSessionTokenRequest request);

  $async.Future<$0.GetWebIdentityTokenResponse> getWebIdentityToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetWebIdentityTokenRequest> $request) async {
    return getWebIdentityToken($call, await $request);
  }

  $async.Future<$0.GetWebIdentityTokenResponse> getWebIdentityToken(
      $grpc.ServiceCall call, $0.GetWebIdentityTokenRequest request);
}
