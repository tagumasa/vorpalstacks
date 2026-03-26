// This is a generated file - do not edit.
//
// Generated from admin_auth.proto.

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

import 'admin_auth.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'admin_auth.pb.dart';

@$pb.GrpcServiceName('admin_auth.AdminAuthService')
class AdminAuthServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AdminAuthServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.LoginResponse> login(
    $0.LoginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$login, request, options: options);
  }

  $grpc.ResponseFuture<$0.LoginResponse> refreshToken(
    $0.RefreshTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$refreshToken, request, options: options);
  }

  $grpc.ResponseFuture<$1.Empty> logout(
    $0.LogoutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$logout, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetCurrentUserResponse> getCurrentUser(
    $0.GetCurrentUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCurrentUser, request, options: options);
  }

  // method descriptors

  static final _$login = $grpc.ClientMethod<$0.LoginRequest, $0.LoginResponse>(
      '/admin_auth.AdminAuthService/Login',
      ($0.LoginRequest value) => value.writeToBuffer(),
      $0.LoginResponse.fromBuffer);
  static final _$refreshToken =
      $grpc.ClientMethod<$0.RefreshTokenRequest, $0.LoginResponse>(
          '/admin_auth.AdminAuthService/RefreshToken',
          ($0.RefreshTokenRequest value) => value.writeToBuffer(),
          $0.LoginResponse.fromBuffer);
  static final _$logout = $grpc.ClientMethod<$0.LogoutRequest, $1.Empty>(
      '/admin_auth.AdminAuthService/Logout',
      ($0.LogoutRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$getCurrentUser =
      $grpc.ClientMethod<$0.GetCurrentUserRequest, $0.GetCurrentUserResponse>(
          '/admin_auth.AdminAuthService/GetCurrentUser',
          ($0.GetCurrentUserRequest value) => value.writeToBuffer(),
          $0.GetCurrentUserResponse.fromBuffer);
}

@$pb.GrpcServiceName('admin_auth.AdminAuthService')
abstract class AdminAuthServiceBase extends $grpc.Service {
  $core.String get $name => 'admin_auth.AdminAuthService';

  AdminAuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.LoginRequest, $0.LoginResponse>(
        'Login',
        login_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RefreshTokenRequest, $0.LoginResponse>(
        'RefreshToken',
        refreshToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RefreshTokenRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LogoutRequest, $1.Empty>(
        'Logout',
        logout_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LogoutRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCurrentUserRequest,
            $0.GetCurrentUserResponse>(
        'GetCurrentUser',
        getCurrentUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCurrentUserRequest.fromBuffer(value),
        ($0.GetCurrentUserResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.LoginResponse> login_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LoginRequest> $request) async {
    return login($call, await $request);
  }

  $async.Future<$0.LoginResponse> login(
      $grpc.ServiceCall call, $0.LoginRequest request);

  $async.Future<$0.LoginResponse> refreshToken_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RefreshTokenRequest> $request) async {
    return refreshToken($call, await $request);
  }

  $async.Future<$0.LoginResponse> refreshToken(
      $grpc.ServiceCall call, $0.RefreshTokenRequest request);

  $async.Future<$1.Empty> logout_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LogoutRequest> $request) async {
    return logout($call, await $request);
  }

  $async.Future<$1.Empty> logout(
      $grpc.ServiceCall call, $0.LogoutRequest request);

  $async.Future<$0.GetCurrentUserResponse> getCurrentUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCurrentUserRequest> $request) async {
    return getCurrentUser($call, await $request);
  }

  $async.Future<$0.GetCurrentUserResponse> getCurrentUser(
      $grpc.ServiceCall call, $0.GetCurrentUserRequest request);
}
