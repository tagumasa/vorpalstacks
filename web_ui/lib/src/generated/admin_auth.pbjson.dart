// This is a generated file - do not edit.
//
// Generated from admin_auth.proto.

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

@$core.Deprecated('Use loginRequestDescriptor instead')
const LoginRequest$json = {
  '1': 'LoginRequest',
  '2': [
    {'1': 'username', '3': 1, '4': 1, '5': 9, '10': 'username'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `LoginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginRequestDescriptor = $convert.base64Decode(
    'CgxMb2dpblJlcXVlc3QSGgoIdXNlcm5hbWUYASABKAlSCHVzZXJuYW1lEhoKCHBhc3N3b3JkGA'
    'IgASgJUghwYXNzd29yZA==');

@$core.Deprecated('Use loginResponseDescriptor instead')
const LoginResponse$json = {
  '1': 'LoginResponse',
  '2': [
    {'1': 'access_token', '3': 1, '4': 1, '5': 9, '10': 'accessToken'},
    {'1': 'refresh_token', '3': 2, '4': 1, '5': 9, '10': 'refreshToken'},
    {'1': 'id_token', '3': 3, '4': 1, '5': 9, '10': 'idToken'},
    {'1': 'expires_in', '3': 4, '4': 1, '5': 5, '10': 'expiresIn'},
    {'1': 'token_type', '3': 5, '4': 1, '5': 9, '10': 'tokenType'},
    {
      '1': 'refresh_expires_in',
      '3': 6,
      '4': 1,
      '5': 5,
      '10': 'refreshExpiresIn'
    },
  ],
};

/// Descriptor for `LoginResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginResponseDescriptor = $convert.base64Decode(
    'Cg1Mb2dpblJlc3BvbnNlEiEKDGFjY2Vzc190b2tlbhgBIAEoCVILYWNjZXNzVG9rZW4SIwoNcm'
    'VmcmVzaF90b2tlbhgCIAEoCVIMcmVmcmVzaFRva2VuEhkKCGlkX3Rva2VuGAMgASgJUgdpZFRv'
    'a2VuEh0KCmV4cGlyZXNfaW4YBCABKAVSCWV4cGlyZXNJbhIdCgp0b2tlbl90eXBlGAUgASgJUg'
    'l0b2tlblR5cGUSLAoScmVmcmVzaF9leHBpcmVzX2luGAYgASgFUhByZWZyZXNoRXhwaXJlc0lu');

@$core.Deprecated('Use refreshTokenRequestDescriptor instead')
const RefreshTokenRequest$json = {
  '1': 'RefreshTokenRequest',
  '2': [
    {'1': 'refresh_token', '3': 1, '4': 1, '5': 9, '10': 'refreshToken'},
  ],
};

/// Descriptor for `RefreshTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refreshTokenRequestDescriptor = $convert.base64Decode(
    'ChNSZWZyZXNoVG9rZW5SZXF1ZXN0EiMKDXJlZnJlc2hfdG9rZW4YASABKAlSDHJlZnJlc2hUb2'
    'tlbg==');

@$core.Deprecated('Use logoutRequestDescriptor instead')
const LogoutRequest$json = {
  '1': 'LogoutRequest',
  '2': [
    {'1': 'access_token', '3': 1, '4': 1, '5': 9, '10': 'accessToken'},
  ],
};

/// Descriptor for `LogoutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logoutRequestDescriptor = $convert.base64Decode(
    'Cg1Mb2dvdXRSZXF1ZXN0EiEKDGFjY2Vzc190b2tlbhgBIAEoCVILYWNjZXNzVG9rZW4=');

@$core.Deprecated('Use getCurrentUserRequestDescriptor instead')
const GetCurrentUserRequest$json = {
  '1': 'GetCurrentUserRequest',
  '2': [
    {'1': 'access_token', '3': 1, '4': 1, '5': 9, '10': 'accessToken'},
  ],
};

/// Descriptor for `GetCurrentUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCurrentUserRequestDescriptor = $convert.base64Decode(
    'ChVHZXRDdXJyZW50VXNlclJlcXVlc3QSIQoMYWNjZXNzX3Rva2VuGAEgASgJUgthY2Nlc3NUb2'
    'tlbg==');

@$core.Deprecated('Use getCurrentUserResponseDescriptor instead')
const GetCurrentUserResponse$json = {
  '1': 'GetCurrentUserResponse',
  '2': [
    {'1': 'username', '3': 1, '4': 1, '5': 9, '10': 'username'},
    {'1': 'arn', '3': 2, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'user_id', '3': 3, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'groups', '3': 4, '4': 3, '5': 9, '10': 'groups'},
    {
      '1': 'attached_policies',
      '3': 5,
      '4': 3,
      '5': 9,
      '10': 'attachedPolicies'
    },
    {'1': 'create_date', '3': 6, '4': 1, '5': 9, '10': 'createDate'},
  ],
};

/// Descriptor for `GetCurrentUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCurrentUserResponseDescriptor = $convert.base64Decode(
    'ChZHZXRDdXJyZW50VXNlclJlc3BvbnNlEhoKCHVzZXJuYW1lGAEgASgJUgh1c2VybmFtZRIQCg'
    'Nhcm4YAiABKAlSA2FybhIXCgd1c2VyX2lkGAMgASgJUgZ1c2VySWQSFgoGZ3JvdXBzGAQgAygJ'
    'UgZncm91cHMSKwoRYXR0YWNoZWRfcG9saWNpZXMYBSADKAlSEGF0dGFjaGVkUG9saWNpZXMSHw'
    'oLY3JlYXRlX2RhdGUYBiABKAlSCmNyZWF0ZURhdGU=');
