// This is a generated file - do not edit.
//
// Generated from admin_config.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use configEntryDescriptor instead')
const ConfigEntry$json = {
  '1': 'ConfigEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
    {'1': 'type', '3': 3, '4': 1, '5': 9, '10': 'type'},
    {'1': 'source', '3': 4, '4': 1, '5': 9, '10': 'source'},
    {'1': 'description', '3': 5, '4': 1, '5': 9, '10': 'description'},
    {'1': 'read_only', '3': 6, '4': 1, '5': 8, '10': 'readOnly'},
    {'1': 'updated_at', '3': 7, '4': 1, '5': 3, '10': 'updatedAt'},
    {'1': 'env_var', '3': 8, '4': 1, '5': 9, '10': 'envVar'},
    {'1': 'category', '3': 9, '4': 1, '5': 9, '10': 'category'},
  ],
};

/// Descriptor for `ConfigEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List configEntryDescriptor = $convert.base64Decode(
    'CgtDb25maWdFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWUSEg'
    'oEdHlwZRgDIAEoCVIEdHlwZRIWCgZzb3VyY2UYBCABKAlSBnNvdXJjZRIgCgtkZXNjcmlwdGlv'
    'bhgFIAEoCVILZGVzY3JpcHRpb24SGwoJcmVhZF9vbmx5GAYgASgIUghyZWFkT25seRIdCgp1cG'
    'RhdGVkX2F0GAcgASgDUgl1cGRhdGVkQXQSFwoHZW52X3ZhchgIIAEoCVIGZW52VmFyEhoKCGNh'
    'dGVnb3J5GAkgASgJUghjYXRlZ29yeQ==');

@$core.Deprecated('Use getConfigRequestDescriptor instead')
const GetConfigRequest$json = {
  '1': 'GetConfigRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
  ],
};

/// Descriptor for `GetConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConfigRequestDescriptor =
    $convert.base64Decode('ChBHZXRDb25maWdSZXF1ZXN0EhAKA2tleRgBIAEoCVIDa2V5');

@$core.Deprecated('Use getConfigResponseDescriptor instead')
const GetConfigResponse$json = {
  '1': 'GetConfigResponse',
  '2': [
    {
      '1': 'entry',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.admin_config.ConfigEntry',
      '10': 'entry'
    },
  ],
};

/// Descriptor for `GetConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConfigResponseDescriptor = $convert.base64Decode(
    'ChFHZXRDb25maWdSZXNwb25zZRIvCgVlbnRyeRgBIAEoCzIZLmFkbWluX2NvbmZpZy5Db25maW'
    'dFbnRyeVIFZW50cnk=');

@$core.Deprecated('Use listConfigRequestDescriptor instead')
const ListConfigRequest$json = {
  '1': 'ListConfigRequest',
  '2': [
    {'1': 'category', '3': 1, '4': 1, '5': 9, '10': 'category'},
  ],
};

/// Descriptor for `ListConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConfigRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0Q29uZmlnUmVxdWVzdBIaCghjYXRlZ29yeRgBIAEoCVIIY2F0ZWdvcnk=');

@$core.Deprecated('Use listConfigResponseDescriptor instead')
const ListConfigResponse$json = {
  '1': 'ListConfigResponse',
  '2': [
    {
      '1': 'entries',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.admin_config.ConfigEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `ListConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConfigResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0Q29uZmlnUmVzcG9uc2USMwoHZW50cmllcxgBIAMoCzIZLmFkbWluX2NvbmZpZy5Db2'
    '5maWdFbnRyeVIHZW50cmllcw==');

@$core.Deprecated('Use updateConfigRequestDescriptor instead')
const UpdateConfigRequest$json = {
  '1': 'UpdateConfigRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `UpdateConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConfigRequestDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVDb25maWdSZXF1ZXN0EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZQ==');

@$core.Deprecated('Use resetConfigRequestDescriptor instead')
const ResetConfigRequest$json = {
  '1': 'ResetConfigRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
  ],
};

/// Descriptor for `ResetConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetConfigRequestDescriptor = $convert
    .base64Decode('ChJSZXNldENvbmZpZ1JlcXVlc3QSEAoDa2V5GAEgASgJUgNrZXk=');

@$core.Deprecated('Use listServicesRequestDescriptor instead')
const ListServicesRequest$json = {
  '1': 'ListServicesRequest',
};

/// Descriptor for `ListServicesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServicesRequestDescriptor =
    $convert.base64Decode('ChNMaXN0U2VydmljZXNSZXF1ZXN0');

@$core.Deprecated('Use listServicesResponseDescriptor instead')
const ListServicesResponse$json = {
  '1': 'ListServicesResponse',
  '2': [
    {
      '1': 'services',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.admin_config.ServiceInfo',
      '10': 'services'
    },
  ],
};

/// Descriptor for `ListServicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServicesResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0U2VydmljZXNSZXNwb25zZRI1CghzZXJ2aWNlcxgBIAMoCzIZLmFkbWluX2NvbmZpZy'
    '5TZXJ2aWNlSW5mb1IIc2VydmljZXM=');

@$core.Deprecated('Use serviceInfoDescriptor instead')
const ServiceInfo$json = {
  '1': 'ServiceInfo',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'enabled', '3': 2, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'resource_count', '3': 3, '4': 1, '5': 5, '10': 'resourceCount'},
  ],
};

/// Descriptor for `ServiceInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceInfoDescriptor = $convert.base64Decode(
    'CgtTZXJ2aWNlSW5mbxISCgRuYW1lGAEgASgJUgRuYW1lEhgKB2VuYWJsZWQYAiABKAhSB2VuYW'
    'JsZWQSJQoOcmVzb3VyY2VfY291bnQYAyABKAVSDXJlc291cmNlQ291bnQ=');

@$core.Deprecated('Use getServiceStatusRequestDescriptor instead')
const GetServiceStatusRequest$json = {
  '1': 'GetServiceStatusRequest',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetServiceStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServiceStatusRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRTZXJ2aWNlU3RhdHVzUmVxdWVzdBISCgRuYW1lGAEgASgJUgRuYW1l');

@$core.Deprecated('Use serviceStatusDescriptor instead')
const ServiceStatus$json = {
  '1': 'ServiceStatus',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'enabled', '3': 2, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'mode', '3': 3, '4': 1, '5': 9, '10': 'mode'},
    {'1': 'resource_count', '3': 4, '4': 1, '5': 5, '10': 'resourceCount'},
    {'1': 'last_request_at', '3': 5, '4': 1, '5': 3, '10': 'lastRequestAt'},
  ],
};

/// Descriptor for `ServiceStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceStatusDescriptor = $convert.base64Decode(
    'Cg1TZXJ2aWNlU3RhdHVzEhIKBG5hbWUYASABKAlSBG5hbWUSGAoHZW5hYmxlZBgCIAEoCFIHZW'
    '5hYmxlZBISCgRtb2RlGAMgASgJUgRtb2RlEiUKDnJlc291cmNlX2NvdW50GAQgASgFUg1yZXNv'
    'dXJjZUNvdW50EiYKD2xhc3RfcmVxdWVzdF9hdBgFIAEoA1INbGFzdFJlcXVlc3RBdA==');

@$core.Deprecated('Use getResourcePortRequestDescriptor instead')
const GetResourcePortRequest$json = {
  '1': 'GetResourcePortRequest',
  '2': [
    {'1': 'service_port_key', '3': 1, '4': 1, '5': 9, '10': 'servicePortKey'},
    {'1': 'resource_id', '3': 2, '4': 1, '5': 9, '10': 'resourceId'},
  ],
};

/// Descriptor for `GetResourcePortRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePortRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRSZXNvdXJjZVBvcnRSZXF1ZXN0EigKEHNlcnZpY2VfcG9ydF9rZXkYASABKAlSDnNlcn'
        'ZpY2VQb3J0S2V5Eh8KC3Jlc291cmNlX2lkGAIgASgJUgpyZXNvdXJjZUlk');

@$core.Deprecated('Use getResourcePortResponseDescriptor instead')
const GetResourcePortResponse$json = {
  '1': 'GetResourcePortResponse',
  '2': [
    {'1': 'port', '3': 1, '4': 1, '5': 5, '10': 'port'},
    {'1': 'source', '3': 2, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `GetResourcePortResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePortResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRSZXNvdXJjZVBvcnRSZXNwb25zZRISCgRwb3J0GAEgASgFUgRwb3J0EhYKBnNvdXJjZR'
        'gCIAEoCVIGc291cmNl');

@$core.Deprecated('Use setResourcePortRequestDescriptor instead')
const SetResourcePortRequest$json = {
  '1': 'SetResourcePortRequest',
  '2': [
    {'1': 'service_port_key', '3': 1, '4': 1, '5': 9, '10': 'servicePortKey'},
    {'1': 'resource_id', '3': 2, '4': 1, '5': 9, '10': 'resourceId'},
    {'1': 'port', '3': 3, '4': 1, '5': 5, '10': 'port'},
  ],
};

/// Descriptor for `SetResourcePortRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setResourcePortRequestDescriptor = $convert.base64Decode(
    'ChZTZXRSZXNvdXJjZVBvcnRSZXF1ZXN0EigKEHNlcnZpY2VfcG9ydF9rZXkYASABKAlSDnNlcn'
    'ZpY2VQb3J0S2V5Eh8KC3Jlc291cmNlX2lkGAIgASgJUgpyZXNvdXJjZUlkEhIKBHBvcnQYAyAB'
    'KAVSBHBvcnQ=');
