// This is a generated file - do not edit.
//
// Generated from admin_config.proto.

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

import 'admin_config.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'admin_config.pb.dart';

@$pb.GrpcServiceName('admin_config.AdminConfigService')
class AdminConfigServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AdminConfigServiceClient(super.channel, {super.options, super.interceptors});

  /// Configuration management (full CRUD)
  $grpc.ResponseFuture<$0.GetConfigResponse> getConfig(
    $0.GetConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConfig, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListConfigResponse> listConfig(
    $0.ListConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConfig, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConfigEntry> updateConfig(
    $0.UpdateConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateConfig, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConfigEntry> resetConfig(
    $0.ResetConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resetConfig, request, options: options);
  }

  /// Service information (read-only)
  $grpc.ResponseFuture<$0.ListServicesResponse> listServices(
    $0.ListServicesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listServices, request, options: options);
  }

  $grpc.ResponseFuture<$0.ServiceStatus> getServiceStatus(
    $0.GetServiceStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getServiceStatus, request, options: options);
  }

  /// Port mappings for resources
  $grpc.ResponseFuture<$0.GetResourcePortResponse> getResourcePort(
    $0.GetResourcePortRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePort, request, options: options);
  }

  $grpc.ResponseFuture<$1.Empty> setResourcePort(
    $0.SetResourcePortRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setResourcePort, request, options: options);
  }

  // method descriptors

  static final _$getConfig =
      $grpc.ClientMethod<$0.GetConfigRequest, $0.GetConfigResponse>(
          '/admin_config.AdminConfigService/GetConfig',
          ($0.GetConfigRequest value) => value.writeToBuffer(),
          $0.GetConfigResponse.fromBuffer);
  static final _$listConfig =
      $grpc.ClientMethod<$0.ListConfigRequest, $0.ListConfigResponse>(
          '/admin_config.AdminConfigService/ListConfig',
          ($0.ListConfigRequest value) => value.writeToBuffer(),
          $0.ListConfigResponse.fromBuffer);
  static final _$updateConfig =
      $grpc.ClientMethod<$0.UpdateConfigRequest, $0.ConfigEntry>(
          '/admin_config.AdminConfigService/UpdateConfig',
          ($0.UpdateConfigRequest value) => value.writeToBuffer(),
          $0.ConfigEntry.fromBuffer);
  static final _$resetConfig =
      $grpc.ClientMethod<$0.ResetConfigRequest, $0.ConfigEntry>(
          '/admin_config.AdminConfigService/ResetConfig',
          ($0.ResetConfigRequest value) => value.writeToBuffer(),
          $0.ConfigEntry.fromBuffer);
  static final _$listServices =
      $grpc.ClientMethod<$0.ListServicesRequest, $0.ListServicesResponse>(
          '/admin_config.AdminConfigService/ListServices',
          ($0.ListServicesRequest value) => value.writeToBuffer(),
          $0.ListServicesResponse.fromBuffer);
  static final _$getServiceStatus =
      $grpc.ClientMethod<$0.GetServiceStatusRequest, $0.ServiceStatus>(
          '/admin_config.AdminConfigService/GetServiceStatus',
          ($0.GetServiceStatusRequest value) => value.writeToBuffer(),
          $0.ServiceStatus.fromBuffer);
  static final _$getResourcePort =
      $grpc.ClientMethod<$0.GetResourcePortRequest, $0.GetResourcePortResponse>(
          '/admin_config.AdminConfigService/GetResourcePort',
          ($0.GetResourcePortRequest value) => value.writeToBuffer(),
          $0.GetResourcePortResponse.fromBuffer);
  static final _$setResourcePort =
      $grpc.ClientMethod<$0.SetResourcePortRequest, $1.Empty>(
          '/admin_config.AdminConfigService/SetResourcePort',
          ($0.SetResourcePortRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('admin_config.AdminConfigService')
abstract class AdminConfigServiceBase extends $grpc.Service {
  $core.String get $name => 'admin_config.AdminConfigService';

  AdminConfigServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetConfigRequest, $0.GetConfigResponse>(
        'GetConfig',
        getConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetConfigRequest.fromBuffer(value),
        ($0.GetConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConfigRequest, $0.ListConfigResponse>(
        'ListConfig',
        listConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListConfigRequest.fromBuffer(value),
        ($0.ListConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateConfigRequest, $0.ConfigEntry>(
        'UpdateConfig',
        updateConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateConfigRequest.fromBuffer(value),
        ($0.ConfigEntry value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResetConfigRequest, $0.ConfigEntry>(
        'ResetConfig',
        resetConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResetConfigRequest.fromBuffer(value),
        ($0.ConfigEntry value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListServicesRequest, $0.ListServicesResponse>(
            'ListServices',
            listServices_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListServicesRequest.fromBuffer(value),
            ($0.ListServicesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetServiceStatusRequest, $0.ServiceStatus>(
            'GetServiceStatus',
            getServiceStatus_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetServiceStatusRequest.fromBuffer(value),
            ($0.ServiceStatus value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePortRequest,
            $0.GetResourcePortResponse>(
        'GetResourcePort',
        getResourcePort_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePortRequest.fromBuffer(value),
        ($0.GetResourcePortResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetResourcePortRequest, $1.Empty>(
        'SetResourcePort',
        setResourcePort_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetResourcePortRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$0.GetConfigResponse> getConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetConfigRequest> $request) async {
    return getConfig($call, await $request);
  }

  $async.Future<$0.GetConfigResponse> getConfig(
      $grpc.ServiceCall call, $0.GetConfigRequest request);

  $async.Future<$0.ListConfigResponse> listConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListConfigRequest> $request) async {
    return listConfig($call, await $request);
  }

  $async.Future<$0.ListConfigResponse> listConfig(
      $grpc.ServiceCall call, $0.ListConfigRequest request);

  $async.Future<$0.ConfigEntry> updateConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateConfigRequest> $request) async {
    return updateConfig($call, await $request);
  }

  $async.Future<$0.ConfigEntry> updateConfig(
      $grpc.ServiceCall call, $0.UpdateConfigRequest request);

  $async.Future<$0.ConfigEntry> resetConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ResetConfigRequest> $request) async {
    return resetConfig($call, await $request);
  }

  $async.Future<$0.ConfigEntry> resetConfig(
      $grpc.ServiceCall call, $0.ResetConfigRequest request);

  $async.Future<$0.ListServicesResponse> listServices_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListServicesRequest> $request) async {
    return listServices($call, await $request);
  }

  $async.Future<$0.ListServicesResponse> listServices(
      $grpc.ServiceCall call, $0.ListServicesRequest request);

  $async.Future<$0.ServiceStatus> getServiceStatus_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetServiceStatusRequest> $request) async {
    return getServiceStatus($call, await $request);
  }

  $async.Future<$0.ServiceStatus> getServiceStatus(
      $grpc.ServiceCall call, $0.GetServiceStatusRequest request);

  $async.Future<$0.GetResourcePortResponse> getResourcePort_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePortRequest> $request) async {
    return getResourcePort($call, await $request);
  }

  $async.Future<$0.GetResourcePortResponse> getResourcePort(
      $grpc.ServiceCall call, $0.GetResourcePortRequest request);

  $async.Future<$1.Empty> setResourcePort_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetResourcePortRequest> $request) async {
    return setResourcePort($call, await $request);
  }

  $async.Future<$1.Empty> setResourcePort(
      $grpc.ServiceCall call, $0.SetResourcePortRequest request);
}
