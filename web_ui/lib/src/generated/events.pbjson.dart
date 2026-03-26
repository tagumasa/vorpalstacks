// This is a generated file - do not edit.
//
// Generated from events.proto.

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

@$core.Deprecated('Use apiDestinationHttpMethodDescriptor instead')
const ApiDestinationHttpMethod$json = {
  '1': 'ApiDestinationHttpMethod',
  '2': [
    {'1': 'API_DESTINATION_HTTP_METHOD_OPTIONS', '2': 0},
    {'1': 'API_DESTINATION_HTTP_METHOD_PATCH', '2': 1},
    {'1': 'API_DESTINATION_HTTP_METHOD_POST', '2': 2},
    {'1': 'API_DESTINATION_HTTP_METHOD_HEAD', '2': 3},
    {'1': 'API_DESTINATION_HTTP_METHOD_GET', '2': 4},
    {'1': 'API_DESTINATION_HTTP_METHOD_DELETE', '2': 5},
    {'1': 'API_DESTINATION_HTTP_METHOD_PUT', '2': 6},
  ],
};

/// Descriptor for `ApiDestinationHttpMethod`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List apiDestinationHttpMethodDescriptor = $convert.base64Decode(
    'ChhBcGlEZXN0aW5hdGlvbkh0dHBNZXRob2QSJwojQVBJX0RFU1RJTkFUSU9OX0hUVFBfTUVUSE'
    '9EX09QVElPTlMQABIlCiFBUElfREVTVElOQVRJT05fSFRUUF9NRVRIT0RfUEFUQ0gQARIkCiBB'
    'UElfREVTVElOQVRJT05fSFRUUF9NRVRIT0RfUE9TVBACEiQKIEFQSV9ERVNUSU5BVElPTl9IVF'
    'RQX01FVEhPRF9IRUFEEAMSIwofQVBJX0RFU1RJTkFUSU9OX0hUVFBfTUVUSE9EX0dFVBAEEiYK'
    'IkFQSV9ERVNUSU5BVElPTl9IVFRQX01FVEhPRF9ERUxFVEUQBRIjCh9BUElfREVTVElOQVRJT0'
    '5fSFRUUF9NRVRIT0RfUFVUEAY=');

@$core.Deprecated('Use apiDestinationStateDescriptor instead')
const ApiDestinationState$json = {
  '1': 'ApiDestinationState',
  '2': [
    {'1': 'API_DESTINATION_STATE_ACTIVE', '2': 0},
    {'1': 'API_DESTINATION_STATE_INACTIVE', '2': 1},
  ],
};

/// Descriptor for `ApiDestinationState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List apiDestinationStateDescriptor = $convert.base64Decode(
    'ChNBcGlEZXN0aW5hdGlvblN0YXRlEiAKHEFQSV9ERVNUSU5BVElPTl9TVEFURV9BQ1RJVkUQAB'
    'IiCh5BUElfREVTVElOQVRJT05fU1RBVEVfSU5BQ1RJVkUQAQ==');

@$core.Deprecated('Use archiveStateDescriptor instead')
const ArchiveState$json = {
  '1': 'ArchiveState',
  '2': [
    {'1': 'ARCHIVE_STATE_UPDATING', '2': 0},
    {'1': 'ARCHIVE_STATE_UPDATE_FAILED', '2': 1},
    {'1': 'ARCHIVE_STATE_DISABLED', '2': 2},
    {'1': 'ARCHIVE_STATE_CREATING', '2': 3},
    {'1': 'ARCHIVE_STATE_ENABLED', '2': 4},
    {'1': 'ARCHIVE_STATE_CREATE_FAILED', '2': 5},
  ],
};

/// Descriptor for `ArchiveState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List archiveStateDescriptor = $convert.base64Decode(
    'CgxBcmNoaXZlU3RhdGUSGgoWQVJDSElWRV9TVEFURV9VUERBVElORxAAEh8KG0FSQ0hJVkVfU1'
    'RBVEVfVVBEQVRFX0ZBSUxFRBABEhoKFkFSQ0hJVkVfU1RBVEVfRElTQUJMRUQQAhIaChZBUkNI'
    'SVZFX1NUQVRFX0NSRUFUSU5HEAMSGQoVQVJDSElWRV9TVEFURV9FTkFCTEVEEAQSHwobQVJDSE'
    'lWRV9TVEFURV9DUkVBVEVfRkFJTEVEEAU=');

@$core.Deprecated('Use assignPublicIpDescriptor instead')
const AssignPublicIp$json = {
  '1': 'AssignPublicIp',
  '2': [
    {'1': 'ASSIGN_PUBLIC_IP_DISABLED', '2': 0},
    {'1': 'ASSIGN_PUBLIC_IP_ENABLED', '2': 1},
  ],
};

/// Descriptor for `AssignPublicIp`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List assignPublicIpDescriptor = $convert.base64Decode(
    'Cg5Bc3NpZ25QdWJsaWNJcBIdChlBU1NJR05fUFVCTElDX0lQX0RJU0FCTEVEEAASHAoYQVNTSU'
    'dOX1BVQkxJQ19JUF9FTkFCTEVEEAE=');

@$core.Deprecated('Use connectionAuthorizationTypeDescriptor instead')
const ConnectionAuthorizationType$json = {
  '1': 'ConnectionAuthorizationType',
  '2': [
    {'1': 'CONNECTION_AUTHORIZATION_TYPE_OAUTH_CLIENT_CREDENTIALS', '2': 0},
    {'1': 'CONNECTION_AUTHORIZATION_TYPE_BASIC', '2': 1},
    {'1': 'CONNECTION_AUTHORIZATION_TYPE_API_KEY', '2': 2},
  ],
};

/// Descriptor for `ConnectionAuthorizationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionAuthorizationTypeDescriptor = $convert.base64Decode(
    'ChtDb25uZWN0aW9uQXV0aG9yaXphdGlvblR5cGUSOgo2Q09OTkVDVElPTl9BVVRIT1JJWkFUSU'
    '9OX1RZUEVfT0FVVEhfQ0xJRU5UX0NSRURFTlRJQUxTEAASJwojQ09OTkVDVElPTl9BVVRIT1JJ'
    'WkFUSU9OX1RZUEVfQkFTSUMQARIpCiVDT05ORUNUSU9OX0FVVEhPUklaQVRJT05fVFlQRV9BUE'
    'lfS0VZEAI=');

@$core.Deprecated('Use connectionOAuthHttpMethodDescriptor instead')
const ConnectionOAuthHttpMethod$json = {
  '1': 'ConnectionOAuthHttpMethod',
  '2': [
    {'1': 'CONNECTION_O_AUTH_HTTP_METHOD_POST', '2': 0},
    {'1': 'CONNECTION_O_AUTH_HTTP_METHOD_GET', '2': 1},
    {'1': 'CONNECTION_O_AUTH_HTTP_METHOD_PUT', '2': 2},
  ],
};

/// Descriptor for `ConnectionOAuthHttpMethod`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionOAuthHttpMethodDescriptor = $convert.base64Decode(
    'ChlDb25uZWN0aW9uT0F1dGhIdHRwTWV0aG9kEiYKIkNPTk5FQ1RJT05fT19BVVRIX0hUVFBfTU'
    'VUSE9EX1BPU1QQABIlCiFDT05ORUNUSU9OX09fQVVUSF9IVFRQX01FVEhPRF9HRVQQARIlCiFD'
    'T05ORUNUSU9OX09fQVVUSF9IVFRQX01FVEhPRF9QVVQQAg==');

@$core.Deprecated('Use connectionStateDescriptor instead')
const ConnectionState$json = {
  '1': 'ConnectionState',
  '2': [
    {'1': 'CONNECTION_STATE_DEAUTHORIZED', '2': 0},
    {'1': 'CONNECTION_STATE_UPDATING', '2': 1},
    {'1': 'CONNECTION_STATE_AUTHORIZING', '2': 2},
    {'1': 'CONNECTION_STATE_AUTHORIZED', '2': 3},
    {'1': 'CONNECTION_STATE_DEAUTHORIZING', '2': 4},
    {'1': 'CONNECTION_STATE_DELETING', '2': 5},
    {'1': 'CONNECTION_STATE_CREATING', '2': 6},
  ],
};

/// Descriptor for `ConnectionState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionStateDescriptor = $convert.base64Decode(
    'Cg9Db25uZWN0aW9uU3RhdGUSIQodQ09OTkVDVElPTl9TVEFURV9ERUFVVEhPUklaRUQQABIdCh'
    'lDT05ORUNUSU9OX1NUQVRFX1VQREFUSU5HEAESIAocQ09OTkVDVElPTl9TVEFURV9BVVRIT1JJ'
    'WklORxACEh8KG0NPTk5FQ1RJT05fU1RBVEVfQVVUSE9SSVpFRBADEiIKHkNPTk5FQ1RJT05fU1'
    'RBVEVfREVBVVRIT1JJWklORxAEEh0KGUNPTk5FQ1RJT05fU1RBVEVfREVMRVRJTkcQBRIdChlD'
    'T05ORUNUSU9OX1NUQVRFX0NSRUFUSU5HEAY=');

@$core.Deprecated('Use eventSourceStateDescriptor instead')
const EventSourceState$json = {
  '1': 'EventSourceState',
  '2': [
    {'1': 'EVENT_SOURCE_STATE_PENDING', '2': 0},
    {'1': 'EVENT_SOURCE_STATE_ACTIVE', '2': 1},
    {'1': 'EVENT_SOURCE_STATE_DELETED', '2': 2},
  ],
};

/// Descriptor for `EventSourceState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventSourceStateDescriptor = $convert.base64Decode(
    'ChBFdmVudFNvdXJjZVN0YXRlEh4KGkVWRU5UX1NPVVJDRV9TVEFURV9QRU5ESU5HEAASHQoZRV'
    'ZFTlRfU09VUkNFX1NUQVRFX0FDVElWRRABEh4KGkVWRU5UX1NPVVJDRV9TVEFURV9ERUxFVEVE'
    'EAI=');

@$core.Deprecated('Use launchTypeDescriptor instead')
const LaunchType$json = {
  '1': 'LaunchType',
  '2': [
    {'1': 'LAUNCH_TYPE_FARGATE', '2': 0},
    {'1': 'LAUNCH_TYPE_EXTERNAL', '2': 1},
    {'1': 'LAUNCH_TYPE_EC2', '2': 2},
  ],
};

/// Descriptor for `LaunchType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List launchTypeDescriptor = $convert.base64Decode(
    'CgpMYXVuY2hUeXBlEhcKE0xBVU5DSF9UWVBFX0ZBUkdBVEUQABIYChRMQVVOQ0hfVFlQRV9FWF'
    'RFUk5BTBABEhMKD0xBVU5DSF9UWVBFX0VDMhAC');

@$core.Deprecated('Use placementConstraintTypeDescriptor instead')
const PlacementConstraintType$json = {
  '1': 'PlacementConstraintType',
  '2': [
    {'1': 'PLACEMENT_CONSTRAINT_TYPE_DISTINCT_INSTANCE', '2': 0},
    {'1': 'PLACEMENT_CONSTRAINT_TYPE_MEMBER_OF', '2': 1},
  ],
};

/// Descriptor for `PlacementConstraintType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List placementConstraintTypeDescriptor = $convert.base64Decode(
    'ChdQbGFjZW1lbnRDb25zdHJhaW50VHlwZRIvCitQTEFDRU1FTlRfQ09OU1RSQUlOVF9UWVBFX0'
    'RJU1RJTkNUX0lOU1RBTkNFEAASJwojUExBQ0VNRU5UX0NPTlNUUkFJTlRfVFlQRV9NRU1CRVJf'
    'T0YQAQ==');

@$core.Deprecated('Use placementStrategyTypeDescriptor instead')
const PlacementStrategyType$json = {
  '1': 'PlacementStrategyType',
  '2': [
    {'1': 'PLACEMENT_STRATEGY_TYPE_SPREAD', '2': 0},
    {'1': 'PLACEMENT_STRATEGY_TYPE_RANDOM', '2': 1},
    {'1': 'PLACEMENT_STRATEGY_TYPE_BINPACK', '2': 2},
  ],
};

/// Descriptor for `PlacementStrategyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List placementStrategyTypeDescriptor = $convert.base64Decode(
    'ChVQbGFjZW1lbnRTdHJhdGVneVR5cGUSIgoeUExBQ0VNRU5UX1NUUkFURUdZX1RZUEVfU1BSRU'
    'FEEAASIgoeUExBQ0VNRU5UX1NUUkFURUdZX1RZUEVfUkFORE9NEAESIwofUExBQ0VNRU5UX1NU'
    'UkFURUdZX1RZUEVfQklOUEFDSxAC');

@$core.Deprecated('Use propagateTagsDescriptor instead')
const PropagateTags$json = {
  '1': 'PropagateTags',
  '2': [
    {'1': 'PROPAGATE_TAGS_TASK_DEFINITION', '2': 0},
  ],
};

/// Descriptor for `PropagateTags`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List propagateTagsDescriptor = $convert.base64Decode(
    'Cg1Qcm9wYWdhdGVUYWdzEiIKHlBST1BBR0FURV9UQUdTX1RBU0tfREVGSU5JVElPThAA');

@$core.Deprecated('Use replayStateDescriptor instead')
const ReplayState$json = {
  '1': 'ReplayState',
  '2': [
    {'1': 'REPLAY_STATE_STARTING', '2': 0},
    {'1': 'REPLAY_STATE_RUNNING', '2': 1},
    {'1': 'REPLAY_STATE_CANCELLED', '2': 2},
    {'1': 'REPLAY_STATE_CANCELLING', '2': 3},
    {'1': 'REPLAY_STATE_COMPLETED', '2': 4},
    {'1': 'REPLAY_STATE_FAILED', '2': 5},
  ],
};

/// Descriptor for `ReplayState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List replayStateDescriptor = $convert.base64Decode(
    'CgtSZXBsYXlTdGF0ZRIZChVSRVBMQVlfU1RBVEVfU1RBUlRJTkcQABIYChRSRVBMQVlfU1RBVE'
    'VfUlVOTklORxABEhoKFlJFUExBWV9TVEFURV9DQU5DRUxMRUQQAhIbChdSRVBMQVlfU1RBVEVf'
    'Q0FOQ0VMTElORxADEhoKFlJFUExBWV9TVEFURV9DT01QTEVURUQQBBIXChNSRVBMQVlfU1RBVE'
    'VfRkFJTEVEEAU=');

@$core.Deprecated('Use ruleStateDescriptor instead')
const RuleState$json = {
  '1': 'RuleState',
  '2': [
    {'1': 'RULE_STATE_DISABLED', '2': 0},
    {'1': 'RULE_STATE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `RuleState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List ruleStateDescriptor = $convert.base64Decode(
    'CglSdWxlU3RhdGUSFwoTUlVMRV9TVEFURV9ESVNBQkxFRBAAEhYKElJVTEVfU1RBVEVfRU5BQk'
    'xFRBAB');

@$core.Deprecated('Use activateEventSourceRequestDescriptor instead')
const ActivateEventSourceRequest$json = {
  '1': 'ActivateEventSourceRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `ActivateEventSourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activateEventSourceRequestDescriptor =
    $convert.base64Decode(
        'ChpBY3RpdmF0ZUV2ZW50U291cmNlUmVxdWVzdBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use apiDestinationDescriptor instead')
const ApiDestination$json = {
  '1': 'ApiDestination',
  '2': [
    {
      '1': 'apidestinationarn',
      '3': 90996885,
      '4': 1,
      '5': 9,
      '10': 'apidestinationarn'
    },
    {
      '1': 'apidestinationstate',
      '3': 13153343,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationState',
      '10': 'apidestinationstate'
    },
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'invocationendpoint',
      '3': 411800759,
      '4': 1,
      '5': 9,
      '10': 'invocationendpoint'
    },
    {
      '1': 'invocationratelimitpersecond',
      '3': 295327816,
      '4': 1,
      '5': 5,
      '10': 'invocationratelimitpersecond'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `ApiDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List apiDestinationDescriptor = $convert.base64Decode(
    'Cg5BcGlEZXN0aW5hdGlvbhIvChFhcGlkZXN0aW5hdGlvbmFybhiVgbIrIAEoCVIRYXBpZGVzdG'
    'luYXRpb25hcm4SUAoTYXBpZGVzdGluYXRpb25zdGF0ZRi/6KIGIAEoDjIbLmV2ZW50cy5BcGlE'
    'ZXN0aW5hdGlvblN0YXRlUhNhcGlkZXN0aW5hdGlvbnN0YXRlEicKDWNvbm5lY3Rpb25hcm4YwY'
    '+8WSABKAlSDWNvbm5lY3Rpb25hcm4SJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlv'
    'bnRpbWUSRAoKaHR0cG1ldGhvZBjRjPy9ASABKA4yIC5ldmVudHMuQXBpRGVzdGluYXRpb25IdH'
    'RwTWV0aG9kUgpodHRwbWV0aG9kEjIKEmludm9jYXRpb25lbmRwb2ludBi3qa7EASABKAlSEmlu'
    'dm9jYXRpb25lbmRwb2ludBJGChxpbnZvY2F0aW9ucmF0ZWxpbWl0cGVyc2Vjb25kGMiw6YwBIA'
    'EoBVIcaW52b2NhdGlvbnJhdGVsaW1pdHBlcnNlY29uZBItChBsYXN0bW9kaWZpZWR0aW1lGOCC'
    '/HAgASgJUhBsYXN0bW9kaWZpZWR0aW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use archiveDescriptor instead')
const Archive$json = {
  '1': 'Archive',
  '2': [
    {'1': 'archivename', '3': 88048487, '4': 1, '5': 9, '10': 'archivename'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'eventcount', '3': 128612839, '4': 1, '5': 3, '10': 'eventcount'},
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'retentiondays',
      '3': 267894223,
      '4': 1,
      '5': 5,
      '10': 'retentiondays'
    },
    {'1': 'sizebytes', '3': 486677664, '4': 1, '5': 3, '10': 'sizebytes'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ArchiveState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `Archive`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List archiveDescriptor = $convert.base64Decode(
    'CgdBcmNoaXZlEiMKC2FyY2hpdmVuYW1lGOeG/ikgASgJUgthcmNoaXZlbmFtZRIlCgxjcmVhdG'
    'lvbnRpbWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRIhCgpldmVudGNvdW50GOfzqT0gASgDUgpl'
    'dmVudGNvdW50EioKDmV2ZW50c291cmNlYXJuGMbKipIBIAEoCVIOZXZlbnRzb3VyY2Vhcm4SJw'
    'oNcmV0ZW50aW9uZGF5cxjP+95/IAEoBVINcmV0ZW50aW9uZGF5cxIgCglzaXplYnl0ZXMYoLmI'
    '6AEgASgDUglzaXplYnl0ZXMSLgoFc3RhdGUYl8my7wEgASgOMhQuZXZlbnRzLkFyY2hpdmVTdG'
    'F0ZVIFc3RhdGUSJAoLc3RhdGVyZWFzb24Y89WtswEgASgJUgtzdGF0ZXJlYXNvbg==');

@$core.Deprecated('Use awsVpcConfigurationDescriptor instead')
const AwsVpcConfiguration$json = {
  '1': 'AwsVpcConfiguration',
  '2': [
    {
      '1': 'assignpublicip',
      '3': 461653589,
      '4': 1,
      '5': 14,
      '6': '.events.AssignPublicIp',
      '10': 'assignpublicip'
    },
    {
      '1': 'securitygroups',
      '3': 515282516,
      '4': 3,
      '5': 9,
      '10': 'securitygroups'
    },
    {'1': 'subnets', '3': 414921506, '4': 3, '5': 9, '10': 'subnets'},
  ],
};

/// Descriptor for `AwsVpcConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List awsVpcConfigurationDescriptor = $convert.base64Decode(
    'ChNBd3NWcGNDb25maWd1cmF0aW9uEkIKDmFzc2lnbnB1YmxpY2lwGNWMkdwBIAEoDjIWLmV2ZW'
    '50cy5Bc3NpZ25QdWJsaWNJcFIOYXNzaWducHVibGljaXASKgoOc2VjdXJpdHlncm91cHMY1Kza'
    '9QEgAygJUg5zZWN1cml0eWdyb3VwcxIcCgdzdWJuZXRzGKLm7MUBIAMoCVIHc3VibmV0cw==');

@$core.Deprecated('Use batchArrayPropertiesDescriptor instead')
const BatchArrayProperties$json = {
  '1': 'BatchArrayProperties',
  '2': [
    {'1': 'size', '3': 105352829, '4': 1, '5': 5, '10': 'size'},
  ],
};

/// Descriptor for `BatchArrayProperties`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchArrayPropertiesDescriptor =
    $convert.base64Decode(
        'ChRCYXRjaEFycmF5UHJvcGVydGllcxIVCgRzaXplGP2cnjIgASgFUgRzaXpl');

@$core.Deprecated('Use batchParametersDescriptor instead')
const BatchParameters$json = {
  '1': 'BatchParameters',
  '2': [
    {
      '1': 'arrayproperties',
      '3': 230899444,
      '4': 1,
      '5': 11,
      '6': '.events.BatchArrayProperties',
      '10': 'arrayproperties'
    },
    {
      '1': 'jobdefinition',
      '3': 420132006,
      '4': 1,
      '5': 9,
      '10': 'jobdefinition'
    },
    {'1': 'jobname', '3': 498531160, '4': 1, '5': 9, '10': 'jobname'},
    {
      '1': 'retrystrategy',
      '3': 105516073,
      '4': 1,
      '5': 11,
      '6': '.events.BatchRetryStrategy',
      '10': 'retrystrategy'
    },
  ],
};

/// Descriptor for `BatchParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchParametersDescriptor = $convert.base64Decode(
    'Cg9CYXRjaFBhcmFtZXRlcnMSSQoPYXJyYXlwcm9wZXJ0aWVzGPT9jG4gASgLMhwuZXZlbnRzLk'
    'JhdGNoQXJyYXlQcm9wZXJ0aWVzUg9hcnJheXByb3BlcnRpZXMSKAoNam9iZGVmaW5pdGlvbhim'
    '6arIASABKAlSDWpvYmRlZmluaXRpb24SHAoHam9ibmFtZRjY9tvtASABKAlSB2pvYm5hbWUSQw'
    'oNcmV0cnlzdHJhdGVneRipmKgyIAEoCzIaLmV2ZW50cy5CYXRjaFJldHJ5U3RyYXRlZ3lSDXJl'
    'dHJ5c3RyYXRlZ3k=');

@$core.Deprecated('Use batchRetryStrategyDescriptor instead')
const BatchRetryStrategy$json = {
  '1': 'BatchRetryStrategy',
  '2': [
    {'1': 'attempts', '3': 132533640, '4': 1, '5': 5, '10': 'attempts'},
  ],
};

/// Descriptor for `BatchRetryStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchRetryStrategyDescriptor =
    $convert.base64Decode(
        'ChJCYXRjaFJldHJ5U3RyYXRlZ3kSHQoIYXR0ZW1wdHMYiJuZPyABKAVSCGF0dGVtcHRz');

@$core.Deprecated('Use cancelReplayRequestDescriptor instead')
const CancelReplayRequest$json = {
  '1': 'CancelReplayRequest',
  '2': [
    {'1': 'replayname', '3': 442173850, '4': 1, '5': 9, '10': 'replayname'},
  ],
};

/// Descriptor for `CancelReplayRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelReplayRequestDescriptor = $convert.base64Decode(
    'ChNDYW5jZWxSZXBsYXlSZXF1ZXN0EiIKCnJlcGxheW5hbWUYmpPs0gEgASgJUgpyZXBsYXluYW'
    '1l');

@$core.Deprecated('Use cancelReplayResponseDescriptor instead')
const CancelReplayResponse$json = {
  '1': 'CancelReplayResponse',
  '2': [
    {'1': 'replayarn', '3': 361946526, '4': 1, '5': 9, '10': 'replayarn'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ReplayState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `CancelReplayResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelReplayResponseDescriptor = $convert.base64Decode(
    'ChRDYW5jZWxSZXBsYXlSZXNwb25zZRIgCglyZXBsYXlhcm4YnrvLrAEgASgJUglyZXBsYXlhcm'
    '4SLQoFc3RhdGUYl8my7wEgASgOMhMuZXZlbnRzLlJlcGxheVN0YXRlUgVzdGF0ZRIkCgtzdGF0'
    'ZXJlYXNvbhjz1a2zASABKAlSC3N0YXRlcmVhc29u');

@$core.Deprecated('Use capacityProviderStrategyItemDescriptor instead')
const CapacityProviderStrategyItem$json = {
  '1': 'CapacityProviderStrategyItem',
  '2': [
    {'1': 'base', '3': 500995289, '4': 1, '5': 5, '10': 'base'},
    {
      '1': 'capacityprovider',
      '3': 109086449,
      '4': 1,
      '5': 9,
      '10': 'capacityprovider'
    },
    {'1': 'weight', '3': 278961850, '4': 1, '5': 5, '10': 'weight'},
  ],
};

/// Descriptor for `CapacityProviderStrategyItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderStrategyItemDescriptor =
    $convert.base64Decode(
        'ChxDYXBhY2l0eVByb3ZpZGVyU3RyYXRlZ3lJdGVtEhYKBGJhc2UY2any7gEgASgFUgRiYXNlEi'
        '0KEGNhcGFjaXR5cHJvdmlkZXIY8Y2CNCABKAlSEGNhcGFjaXR5cHJvdmlkZXISGgoGd2VpZ2h0'
        'GLq9goUBIAEoBVIGd2VpZ2h0');

@$core.Deprecated('Use concurrentModificationExceptionDescriptor instead')
const ConcurrentModificationException$json = {
  '1': 'ConcurrentModificationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentModificationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentModificationExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Db25jdXJyZW50TW9kaWZpY2F0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use conditionDescriptor instead')
const Condition$json = {
  '1': 'Condition',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Condition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conditionDescriptor = $convert.base64Decode(
    'CglDb25kaXRpb24SEwoDa2V5GI2S62ggASgJUgNrZXkSFgoEdHlwZRjuoNeKASABKAlSBHR5cG'
    'USGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use connectionDescriptor instead')
const Connection$json = {
  '1': 'Connection',
  '2': [
    {
      '1': 'authorizationtype',
      '3': 481976511,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionAuthorizationType',
      '10': 'authorizationtype'
    },
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastauthorizedtime',
      '3': 250638066,
      '4': 1,
      '5': 9,
      '10': 'lastauthorizedtime'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `Connection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionDescriptor = $convert.base64Decode(
    'CgpDb25uZWN0aW9uElUKEWF1dGhvcml6YXRpb250eXBlGL/B6eUBIAEoDjIjLmV2ZW50cy5Db2'
    '5uZWN0aW9uQXV0aG9yaXphdGlvblR5cGVSEWF1dGhvcml6YXRpb250eXBlEicKDWNvbm5lY3Rp'
    'b25hcm4YwY+8WSABKAlSDWNvbm5lY3Rpb25hcm4SRQoPY29ubmVjdGlvbnN0YXRlGNv65cABIA'
    'EoDjIXLmV2ZW50cy5Db25uZWN0aW9uU3RhdGVSD2Nvbm5lY3Rpb25zdGF0ZRIlCgxjcmVhdGlv'
    'bnRpbWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRIxChJsYXN0YXV0aG9yaXplZHRpbWUY8t3Bdy'
    'ABKAlSEmxhc3RhdXRob3JpemVkdGltZRItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUhBs'
    'YXN0bW9kaWZpZWR0aW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJAoLc3RhdGVyZWFzb24Y89'
    'WtswEgASgJUgtzdGF0ZXJlYXNvbg==');

@$core
    .Deprecated('Use connectionApiKeyAuthResponseParametersDescriptor instead')
const ConnectionApiKeyAuthResponseParameters$json = {
  '1': 'ConnectionApiKeyAuthResponseParameters',
  '2': [
    {'1': 'apikeyname', '3': 43133502, '4': 1, '5': 9, '10': 'apikeyname'},
  ],
};

/// Descriptor for `ConnectionApiKeyAuthResponseParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionApiKeyAuthResponseParametersDescriptor =
    $convert.base64Decode(
        'CiZDb25uZWN0aW9uQXBpS2V5QXV0aFJlc3BvbnNlUGFyYW1ldGVycxIhCgphcGlrZXluYW1lGL'
        '7UyBQgASgJUgphcGlrZXluYW1l');

@$core.Deprecated('Use connectionAuthResponseParametersDescriptor instead')
const ConnectionAuthResponseParameters$json = {
  '1': 'ConnectionAuthResponseParameters',
  '2': [
    {
      '1': 'apikeyauthparameters',
      '3': 110622489,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionApiKeyAuthResponseParameters',
      '10': 'apikeyauthparameters'
    },
    {
      '1': 'basicauthparameters',
      '3': 63965312,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionBasicAuthResponseParameters',
      '10': 'basicauthparameters'
    },
    {
      '1': 'invocationhttpparameters',
      '3': 499128572,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionHttpParameters',
      '10': 'invocationhttpparameters'
    },
    {
      '1': 'oauthparameters',
      '3': 206836569,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionOAuthResponseParameters',
      '10': 'oauthparameters'
    },
  ],
};

/// Descriptor for `ConnectionAuthResponseParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionAuthResponseParametersDescriptor = $convert.base64Decode(
    'CiBDb25uZWN0aW9uQXV0aFJlc3BvbnNlUGFyYW1ldGVycxJlChRhcGlrZXlhdXRocGFyYW1ldG'
    'VycxiZ7t80IAEoCzIuLmV2ZW50cy5Db25uZWN0aW9uQXBpS2V5QXV0aFJlc3BvbnNlUGFyYW1l'
    'dGVyc1IUYXBpa2V5YXV0aHBhcmFtZXRlcnMSYgoTYmFzaWNhdXRocGFyYW1ldGVycxiAkcAeIA'
    'EoCzItLmV2ZW50cy5Db25uZWN0aW9uQmFzaWNBdXRoUmVzcG9uc2VQYXJhbWV0ZXJzUhNiYXNp'
    'Y2F1dGhwYXJhbWV0ZXJzEmAKGGludm9jYXRpb25odHRwcGFyYW1ldGVycxj8sYDuASABKAsyIC'
    '5ldmVudHMuQ29ubmVjdGlvbkh0dHBQYXJhbWV0ZXJzUhhpbnZvY2F0aW9uaHR0cHBhcmFtZXRl'
    'cnMSVgoPb2F1dGhwYXJhbWV0ZXJzGNmm0GIgASgLMikuZXZlbnRzLkNvbm5lY3Rpb25PQXV0aF'
    'Jlc3BvbnNlUGFyYW1ldGVyc1IPb2F1dGhwYXJhbWV0ZXJz');

@$core.Deprecated('Use connectionBasicAuthResponseParametersDescriptor instead')
const ConnectionBasicAuthResponseParameters$json = {
  '1': 'ConnectionBasicAuthResponseParameters',
  '2': [
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ConnectionBasicAuthResponseParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionBasicAuthResponseParametersDescriptor =
    $convert.base64Decode(
        'CiVDb25uZWN0aW9uQmFzaWNBdXRoUmVzcG9uc2VQYXJhbWV0ZXJzEh4KCHVzZXJuYW1lGNqpo+'
        'ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use connectionBodyParameterDescriptor instead')
const ConnectionBodyParameter$json = {
  '1': 'ConnectionBodyParameter',
  '2': [
    {
      '1': 'isvaluesecret',
      '3': 157884907,
      '4': 1,
      '5': 8,
      '10': 'isvaluesecret'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ConnectionBodyParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionBodyParameterDescriptor = $convert.base64Decode(
    'ChdDb25uZWN0aW9uQm9keVBhcmFtZXRlchInCg1pc3ZhbHVlc2VjcmV0GOvDpEsgASgIUg1pc3'
    'ZhbHVlc2VjcmV0EhMKA2tleRiNkutoIAEoCVIDa2V5EhgKBXZhbHVlGOvyn4oBIAEoCVIFdmFs'
    'dWU=');

@$core.Deprecated('Use connectionHeaderParameterDescriptor instead')
const ConnectionHeaderParameter$json = {
  '1': 'ConnectionHeaderParameter',
  '2': [
    {
      '1': 'isvaluesecret',
      '3': 157884907,
      '4': 1,
      '5': 8,
      '10': 'isvaluesecret'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ConnectionHeaderParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionHeaderParameterDescriptor = $convert.base64Decode(
    'ChlDb25uZWN0aW9uSGVhZGVyUGFyYW1ldGVyEicKDWlzdmFsdWVzZWNyZXQY68OkSyABKAhSDW'
    'lzdmFsdWVzZWNyZXQSEwoDa2V5GI2S62ggASgJUgNrZXkSGAoFdmFsdWUY6/KfigEgASgJUgV2'
    'YWx1ZQ==');

@$core.Deprecated('Use connectionHttpParametersDescriptor instead')
const ConnectionHttpParameters$json = {
  '1': 'ConnectionHttpParameters',
  '2': [
    {
      '1': 'bodyparameters',
      '3': 531363370,
      '4': 3,
      '5': 11,
      '6': '.events.ConnectionBodyParameter',
      '10': 'bodyparameters'
    },
    {
      '1': 'headerparameters',
      '3': 148944757,
      '4': 3,
      '5': 11,
      '6': '.events.ConnectionHeaderParameter',
      '10': 'headerparameters'
    },
    {
      '1': 'querystringparameters',
      '3': 258002263,
      '4': 3,
      '5': 11,
      '6': '.events.ConnectionQueryStringParameter',
      '10': 'querystringparameters'
    },
  ],
};

/// Descriptor for `ConnectionHttpParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionHttpParametersDescriptor = $convert.base64Decode(
    'ChhDb25uZWN0aW9uSHR0cFBhcmFtZXRlcnMSSwoOYm9keXBhcmFtZXRlcnMYquyv/QEgAygLMh'
    '8uZXZlbnRzLkNvbm5lY3Rpb25Cb2R5UGFyYW1ldGVyUg5ib2R5cGFyYW1ldGVycxJQChBoZWFk'
    'ZXJwYXJhbWV0ZXJzGPXugkcgAygLMiEuZXZlbnRzLkNvbm5lY3Rpb25IZWFkZXJQYXJhbWV0ZX'
    'JSEGhlYWRlcnBhcmFtZXRlcnMSXwoVcXVlcnlzdHJpbmdwYXJhbWV0ZXJzGNeag3sgAygLMiYu'
    'ZXZlbnRzLkNvbm5lY3Rpb25RdWVyeVN0cmluZ1BhcmFtZXRlclIVcXVlcnlzdHJpbmdwYXJhbW'
    'V0ZXJz');

@$core
    .Deprecated('Use connectionOAuthClientResponseParametersDescriptor instead')
const ConnectionOAuthClientResponseParameters$json = {
  '1': 'ConnectionOAuthClientResponseParameters',
  '2': [
    {'1': 'clientid', '3': 448889284, '4': 1, '5': 9, '10': 'clientid'},
  ],
};

/// Descriptor for `ConnectionOAuthClientResponseParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionOAuthClientResponseParametersDescriptor =
    $convert.base64Decode(
        'CidDb25uZWN0aW9uT0F1dGhDbGllbnRSZXNwb25zZVBhcmFtZXRlcnMSHgoIY2xpZW50aWQYxI'
        'OG1gEgASgJUghjbGllbnRpZA==');

@$core.Deprecated('Use connectionOAuthResponseParametersDescriptor instead')
const ConnectionOAuthResponseParameters$json = {
  '1': 'ConnectionOAuthResponseParameters',
  '2': [
    {
      '1': 'authorizationendpoint',
      '3': 427938596,
      '4': 1,
      '5': 9,
      '10': 'authorizationendpoint'
    },
    {
      '1': 'clientparameters',
      '3': 246864127,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionOAuthClientResponseParameters',
      '10': 'clientparameters'
    },
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionOAuthHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'oauthhttpparameters',
      '3': 10294537,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionHttpParameters',
      '10': 'oauthhttpparameters'
    },
  ],
};

/// Descriptor for `ConnectionOAuthResponseParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionOAuthResponseParametersDescriptor = $convert.base64Decode(
    'CiFDb25uZWN0aW9uT0F1dGhSZXNwb25zZVBhcmFtZXRlcnMSOAoVYXV0aG9yaXphdGlvbmVuZH'
    'BvaW50GKSmh8wBIAEoCVIVYXV0aG9yaXphdGlvbmVuZHBvaW50El4KEGNsaWVudHBhcmFtZXRl'
    'cnMY/7HbdSABKAsyLy5ldmVudHMuQ29ubmVjdGlvbk9BdXRoQ2xpZW50UmVzcG9uc2VQYXJhbW'
    'V0ZXJzUhBjbGllbnRwYXJhbWV0ZXJzEkUKCmh0dHBtZXRob2QY0Yz8vQEgASgOMiEuZXZlbnRz'
    'LkNvbm5lY3Rpb25PQXV0aEh0dHBNZXRob2RSCmh0dHBtZXRob2QSVQoTb2F1dGhodHRwcGFyYW'
    '1ldGVycxiJqvQEIAEoCzIgLmV2ZW50cy5Db25uZWN0aW9uSHR0cFBhcmFtZXRlcnNSE29hdXRo'
    'aHR0cHBhcmFtZXRlcnM=');

@$core.Deprecated('Use connectionQueryStringParameterDescriptor instead')
const ConnectionQueryStringParameter$json = {
  '1': 'ConnectionQueryStringParameter',
  '2': [
    {
      '1': 'isvaluesecret',
      '3': 157884907,
      '4': 1,
      '5': 8,
      '10': 'isvaluesecret'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ConnectionQueryStringParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionQueryStringParameterDescriptor =
    $convert.base64Decode(
        'Ch5Db25uZWN0aW9uUXVlcnlTdHJpbmdQYXJhbWV0ZXISJwoNaXN2YWx1ZXNlY3JldBjrw6RLIA'
        'EoCFINaXN2YWx1ZXNlY3JldBITCgNrZXkYjZLraCABKAlSA2tleRIYCgV2YWx1ZRjr8p+KASAB'
        'KAlSBXZhbHVl');

@$core.Deprecated('Use createApiDestinationRequestDescriptor instead')
const CreateApiDestinationRequest$json = {
  '1': 'CreateApiDestinationRequest',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'invocationendpoint',
      '3': 411800759,
      '4': 1,
      '5': 9,
      '10': 'invocationendpoint'
    },
    {
      '1': 'invocationratelimitpersecond',
      '3': 295327816,
      '4': 1,
      '5': 5,
      '10': 'invocationratelimitpersecond'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateApiDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createApiDestinationRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVBcGlEZXN0aW5hdGlvblJlcXVlc3QSJwoNY29ubmVjdGlvbmFybhjBj7xZIAEoCV'
    'INY29ubmVjdGlvbmFybhIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SRAoK'
    'aHR0cG1ldGhvZBjRjPy9ASABKA4yIC5ldmVudHMuQXBpRGVzdGluYXRpb25IdHRwTWV0aG9kUg'
    'podHRwbWV0aG9kEjIKEmludm9jYXRpb25lbmRwb2ludBi3qa7EASABKAlSEmludm9jYXRpb25l'
    'bmRwb2ludBJGChxpbnZvY2F0aW9ucmF0ZWxpbWl0cGVyc2Vjb25kGMiw6YwBIAEoBVIcaW52b2'
    'NhdGlvbnJhdGVsaW1pdHBlcnNlY29uZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use createApiDestinationResponseDescriptor instead')
const CreateApiDestinationResponse$json = {
  '1': 'CreateApiDestinationResponse',
  '2': [
    {
      '1': 'apidestinationarn',
      '3': 90996885,
      '4': 1,
      '5': 9,
      '10': 'apidestinationarn'
    },
    {
      '1': 'apidestinationstate',
      '3': 13153343,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationState',
      '10': 'apidestinationstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `CreateApiDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createApiDestinationResponseDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVBcGlEZXN0aW5hdGlvblJlc3BvbnNlEi8KEWFwaWRlc3RpbmF0aW9uYXJuGJWBsi'
    'sgASgJUhFhcGlkZXN0aW5hdGlvbmFybhJQChNhcGlkZXN0aW5hdGlvbnN0YXRlGL/oogYgASgO'
    'MhsuZXZlbnRzLkFwaURlc3RpbmF0aW9uU3RhdGVSE2FwaWRlc3RpbmF0aW9uc3RhdGUSJQoMY3'
    'JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbnRpbWUSLQoQbGFzdG1vZGlmaWVkdGltZRjg'
    'gvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZQ==');

@$core.Deprecated('Use createArchiveRequestDescriptor instead')
const CreateArchiveRequest$json = {
  '1': 'CreateArchiveRequest',
  '2': [
    {'1': 'archivename', '3': 88048487, '4': 1, '5': 9, '10': 'archivename'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'retentiondays',
      '3': 267894223,
      '4': 1,
      '5': 5,
      '10': 'retentiondays'
    },
  ],
};

/// Descriptor for `CreateArchiveRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createArchiveRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVBcmNoaXZlUmVxdWVzdBIjCgthcmNoaXZlbmFtZRjnhv4pIAEoCVILYXJjaGl2ZW'
    '5hbWUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEiUKDGV2ZW50cGF0dGVy'
    'bhiA96pvIAEoCVIMZXZlbnRwYXR0ZXJuEioKDmV2ZW50c291cmNlYXJuGMbKipIBIAEoCVIOZX'
    'ZlbnRzb3VyY2Vhcm4SJwoNcmV0ZW50aW9uZGF5cxjP+95/IAEoBVINcmV0ZW50aW9uZGF5cw==');

@$core.Deprecated('Use createArchiveResponseDescriptor instead')
const CreateArchiveResponse$json = {
  '1': 'CreateArchiveResponse',
  '2': [
    {'1': 'archivearn', '3': 56866685, '4': 1, '5': 9, '10': 'archivearn'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ArchiveState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `CreateArchiveResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createArchiveResponseDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVBcmNoaXZlUmVzcG9uc2USIQoKYXJjaGl2ZWFybhj97o4bIAEoCVIKYXJjaGl2ZW'
    'FybhIlCgxjcmVhdGlvbnRpbWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRIuCgVzdGF0ZRiXybLv'
    'ASABKA4yFC5ldmVudHMuQXJjaGl2ZVN0YXRlUgVzdGF0ZRIkCgtzdGF0ZXJlYXNvbhjz1a2zAS'
    'ABKAlSC3N0YXRlcmVhc29u');

@$core.Deprecated(
    'Use createConnectionApiKeyAuthRequestParametersDescriptor instead')
const CreateConnectionApiKeyAuthRequestParameters$json = {
  '1': 'CreateConnectionApiKeyAuthRequestParameters',
  '2': [
    {'1': 'apikeyname', '3': 43133502, '4': 1, '5': 9, '10': 'apikeyname'},
    {'1': 'apikeyvalue', '3': 93786144, '4': 1, '5': 9, '10': 'apikeyvalue'},
  ],
};

/// Descriptor for `CreateConnectionApiKeyAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createConnectionApiKeyAuthRequestParametersDescriptor =
    $convert.base64Decode(
        'CitDcmVhdGVDb25uZWN0aW9uQXBpS2V5QXV0aFJlcXVlc3RQYXJhbWV0ZXJzEiEKCmFwaWtleW'
        '5hbWUYvtTIFCABKAlSCmFwaWtleW5hbWUSIwoLYXBpa2V5dmFsdWUYoKDcLCABKAlSC2FwaWtl'
        'eXZhbHVl');

@$core.Deprecated('Use createConnectionAuthRequestParametersDescriptor instead')
const CreateConnectionAuthRequestParameters$json = {
  '1': 'CreateConnectionAuthRequestParameters',
  '2': [
    {
      '1': 'apikeyauthparameters',
      '3': 110622489,
      '4': 1,
      '5': 11,
      '6': '.events.CreateConnectionApiKeyAuthRequestParameters',
      '10': 'apikeyauthparameters'
    },
    {
      '1': 'basicauthparameters',
      '3': 63965312,
      '4': 1,
      '5': 11,
      '6': '.events.CreateConnectionBasicAuthRequestParameters',
      '10': 'basicauthparameters'
    },
    {
      '1': 'invocationhttpparameters',
      '3': 499128572,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionHttpParameters',
      '10': 'invocationhttpparameters'
    },
    {
      '1': 'oauthparameters',
      '3': 206836569,
      '4': 1,
      '5': 11,
      '6': '.events.CreateConnectionOAuthRequestParameters',
      '10': 'oauthparameters'
    },
  ],
};

/// Descriptor for `CreateConnectionAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionAuthRequestParametersDescriptor = $convert.base64Decode(
    'CiVDcmVhdGVDb25uZWN0aW9uQXV0aFJlcXVlc3RQYXJhbWV0ZXJzEmoKFGFwaWtleWF1dGhwYX'
    'JhbWV0ZXJzGJnu3zQgASgLMjMuZXZlbnRzLkNyZWF0ZUNvbm5lY3Rpb25BcGlLZXlBdXRoUmVx'
    'dWVzdFBhcmFtZXRlcnNSFGFwaWtleWF1dGhwYXJhbWV0ZXJzEmcKE2Jhc2ljYXV0aHBhcmFtZX'
    'RlcnMYgJHAHiABKAsyMi5ldmVudHMuQ3JlYXRlQ29ubmVjdGlvbkJhc2ljQXV0aFJlcXVlc3RQ'
    'YXJhbWV0ZXJzUhNiYXNpY2F1dGhwYXJhbWV0ZXJzEmAKGGludm9jYXRpb25odHRwcGFyYW1ldG'
    'Vycxj8sYDuASABKAsyIC5ldmVudHMuQ29ubmVjdGlvbkh0dHBQYXJhbWV0ZXJzUhhpbnZvY2F0'
    'aW9uaHR0cHBhcmFtZXRlcnMSWwoPb2F1dGhwYXJhbWV0ZXJzGNmm0GIgASgLMi4uZXZlbnRzLk'
    'NyZWF0ZUNvbm5lY3Rpb25PQXV0aFJlcXVlc3RQYXJhbWV0ZXJzUg9vYXV0aHBhcmFtZXRlcnM=');

@$core.Deprecated(
    'Use createConnectionBasicAuthRequestParametersDescriptor instead')
const CreateConnectionBasicAuthRequestParameters$json = {
  '1': 'CreateConnectionBasicAuthRequestParameters',
  '2': [
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `CreateConnectionBasicAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createConnectionBasicAuthRequestParametersDescriptor =
    $convert.base64Decode(
        'CipDcmVhdGVDb25uZWN0aW9uQmFzaWNBdXRoUmVxdWVzdFBhcmFtZXRlcnMSHQoIcGFzc3dvcm'
        'QYuZCMZiABKAlSCHBhc3N3b3JkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated(
    'Use createConnectionOAuthClientRequestParametersDescriptor instead')
const CreateConnectionOAuthClientRequestParameters$json = {
  '1': 'CreateConnectionOAuthClientRequestParameters',
  '2': [
    {'1': 'clientid', '3': 448889284, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
  ],
};

/// Descriptor for `CreateConnectionOAuthClientRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createConnectionOAuthClientRequestParametersDescriptor =
    $convert.base64Decode(
        'CixDcmVhdGVDb25uZWN0aW9uT0F1dGhDbGllbnRSZXF1ZXN0UGFyYW1ldGVycxIeCghjbGllbn'
        'RpZBjEg4bWASABKAlSCGNsaWVudGlkEiYKDGNsaWVudHNlY3JldBj3teLuASABKAlSDGNsaWVu'
        'dHNlY3JldA==');

@$core
    .Deprecated('Use createConnectionOAuthRequestParametersDescriptor instead')
const CreateConnectionOAuthRequestParameters$json = {
  '1': 'CreateConnectionOAuthRequestParameters',
  '2': [
    {
      '1': 'authorizationendpoint',
      '3': 427938596,
      '4': 1,
      '5': 9,
      '10': 'authorizationendpoint'
    },
    {
      '1': 'clientparameters',
      '3': 246864127,
      '4': 1,
      '5': 11,
      '6': '.events.CreateConnectionOAuthClientRequestParameters',
      '10': 'clientparameters'
    },
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionOAuthHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'oauthhttpparameters',
      '3': 10294537,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionHttpParameters',
      '10': 'oauthhttpparameters'
    },
  ],
};

/// Descriptor for `CreateConnectionOAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionOAuthRequestParametersDescriptor = $convert.base64Decode(
    'CiZDcmVhdGVDb25uZWN0aW9uT0F1dGhSZXF1ZXN0UGFyYW1ldGVycxI4ChVhdXRob3JpemF0aW'
    '9uZW5kcG9pbnQYpKaHzAEgASgJUhVhdXRob3JpemF0aW9uZW5kcG9pbnQSYwoQY2xpZW50cGFy'
    'YW1ldGVycxj/sdt1IAEoCzI0LmV2ZW50cy5DcmVhdGVDb25uZWN0aW9uT0F1dGhDbGllbnRSZX'
    'F1ZXN0UGFyYW1ldGVyc1IQY2xpZW50cGFyYW1ldGVycxJFCgpodHRwbWV0aG9kGNGM/L0BIAEo'
    'DjIhLmV2ZW50cy5Db25uZWN0aW9uT0F1dGhIdHRwTWV0aG9kUgpodHRwbWV0aG9kElUKE29hdX'
    'RoaHR0cHBhcmFtZXRlcnMYiar0BCABKAsyIC5ldmVudHMuQ29ubmVjdGlvbkh0dHBQYXJhbWV0'
    'ZXJzUhNvYXV0aGh0dHBwYXJhbWV0ZXJz');

@$core.Deprecated('Use createConnectionRequestDescriptor instead')
const CreateConnectionRequest$json = {
  '1': 'CreateConnectionRequest',
  '2': [
    {
      '1': 'authparameters',
      '3': 258276552,
      '4': 1,
      '5': 11,
      '6': '.events.CreateConnectionAuthRequestParameters',
      '10': 'authparameters'
    },
    {
      '1': 'authorizationtype',
      '3': 481976511,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionAuthorizationType',
      '10': 'authorizationtype'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateConnectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVDb25uZWN0aW9uUmVxdWVzdBJYCg5hdXRocGFyYW1ldGVycxjI+ZN7IAEoCzItLm'
    'V2ZW50cy5DcmVhdGVDb25uZWN0aW9uQXV0aFJlcXVlc3RQYXJhbWV0ZXJzUg5hdXRocGFyYW1l'
    'dGVycxJVChFhdXRob3JpemF0aW9udHlwZRi/wenlASABKA4yIy5ldmVudHMuQ29ubmVjdGlvbk'
    'F1dGhvcml6YXRpb25UeXBlUhFhdXRob3JpemF0aW9udHlwZRIjCgtkZXNjcmlwdGlvbhiK9Pk2'
    'IAEoCVILZGVzY3JpcHRpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use createConnectionResponseDescriptor instead')
const CreateConnectionResponse$json = {
  '1': 'CreateConnectionResponse',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `CreateConnectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionResponseDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVDb25uZWN0aW9uUmVzcG9uc2USJwoNY29ubmVjdGlvbmFybhjBj7xZIAEoCVINY2'
    '9ubmVjdGlvbmFybhJFCg9jb25uZWN0aW9uc3RhdGUY2/rlwAEgASgOMhcuZXZlbnRzLkNvbm5l'
    'Y3Rpb25TdGF0ZVIPY29ubmVjdGlvbnN0YXRlEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3'
    'JlYXRpb250aW1lEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRp'
    'bWU=');

@$core.Deprecated('Use createEventBusRequestDescriptor instead')
const CreateEventBusRequest$json = {
  '1': 'CreateEventBusRequest',
  '2': [
    {
      '1': 'eventsourcename',
      '3': 427669794,
      '4': 1,
      '5': 9,
      '10': 'eventsourcename'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.events.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateEventBusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEventBusRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVFdmVudEJ1c1JlcXVlc3QSLAoPZXZlbnRzb3VyY2VuYW1lGKLy9ssBIAEoCVIPZX'
    'ZlbnRzb3VyY2VuYW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSIwoEdGFncxjBwfa1ASADKAsy'
    'Cy5ldmVudHMuVGFnUgR0YWdz');

@$core.Deprecated('Use createEventBusResponseDescriptor instead')
const CreateEventBusResponse$json = {
  '1': 'CreateEventBusResponse',
  '2': [
    {'1': 'eventbusarn', '3': 515755843, '4': 1, '5': 9, '10': 'eventbusarn'},
  ],
};

/// Descriptor for `CreateEventBusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEventBusResponseDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVFdmVudEJ1c1Jlc3BvbnNlEiQKC2V2ZW50YnVzYXJuGMOe9/UBIAEoCVILZXZlbn'
        'RidXNhcm4=');

@$core.Deprecated('Use createPartnerEventSourceRequestDescriptor instead')
const CreatePartnerEventSourceRequest$json = {
  '1': 'CreatePartnerEventSourceRequest',
  '2': [
    {'1': 'account', '3': 435725053, '4': 1, '5': 9, '10': 'account'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreatePartnerEventSourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPartnerEventSourceRequestDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVQYXJ0bmVyRXZlbnRTb3VyY2VSZXF1ZXN0EhwKB2FjY291bnQY/cXizwEgASgJUg'
        'dhY2NvdW50EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use createPartnerEventSourceResponseDescriptor instead')
const CreatePartnerEventSourceResponse$json = {
  '1': 'CreatePartnerEventSourceResponse',
  '2': [
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
  ],
};

/// Descriptor for `CreatePartnerEventSourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPartnerEventSourceResponseDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVQYXJ0bmVyRXZlbnRTb3VyY2VSZXNwb25zZRIqCg5ldmVudHNvdXJjZWFybhjGyo'
        'qSASABKAlSDmV2ZW50c291cmNlYXJu');

@$core.Deprecated('Use deactivateEventSourceRequestDescriptor instead')
const DeactivateEventSourceRequest$json = {
  '1': 'DeactivateEventSourceRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeactivateEventSourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deactivateEventSourceRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWFjdGl2YXRlRXZlbnRTb3VyY2VSZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deadLetterConfigDescriptor instead')
const DeadLetterConfig$json = {
  '1': 'DeadLetterConfig',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `DeadLetterConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deadLetterConfigDescriptor = $convert
    .base64Decode('ChBEZWFkTGV0dGVyQ29uZmlnEhQKA2Fybhidm+2/ASABKAlSA2Fybg==');

@$core.Deprecated('Use deauthorizeConnectionRequestDescriptor instead')
const DeauthorizeConnectionRequest$json = {
  '1': 'DeauthorizeConnectionRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeauthorizeConnectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deauthorizeConnectionRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWF1dGhvcml6ZUNvbm5lY3Rpb25SZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deauthorizeConnectionResponseDescriptor instead')
const DeauthorizeConnectionResponse$json = {
  '1': 'DeauthorizeConnectionResponse',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastauthorizedtime',
      '3': 250638066,
      '4': 1,
      '5': 9,
      '10': 'lastauthorizedtime'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `DeauthorizeConnectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deauthorizeConnectionResponseDescriptor = $convert.base64Decode(
    'Ch1EZWF1dGhvcml6ZUNvbm5lY3Rpb25SZXNwb25zZRInCg1jb25uZWN0aW9uYXJuGMGPvFkgAS'
    'gJUg1jb25uZWN0aW9uYXJuEkUKD2Nvbm5lY3Rpb25zdGF0ZRjb+uXAASABKA4yFy5ldmVudHMu'
    'Q29ubmVjdGlvblN0YXRlUg9jb25uZWN0aW9uc3RhdGUSJQoMY3JlYXRpb250aW1lGObPqjEgAS'
    'gJUgxjcmVhdGlvbnRpbWUSMQoSbGFzdGF1dGhvcml6ZWR0aW1lGPLdwXcgASgJUhJsYXN0YXV0'
    'aG9yaXplZHRpbWUSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdG'
    'ltZQ==');

@$core.Deprecated('Use deleteApiDestinationRequestDescriptor instead')
const DeleteApiDestinationRequest$json = {
  '1': 'DeleteApiDestinationRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteApiDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteApiDestinationRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVBcGlEZXN0aW5hdGlvblJlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use deleteApiDestinationResponseDescriptor instead')
const DeleteApiDestinationResponse$json = {
  '1': 'DeleteApiDestinationResponse',
};

/// Descriptor for `DeleteApiDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteApiDestinationResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVBcGlEZXN0aW5hdGlvblJlc3BvbnNl');

@$core.Deprecated('Use deleteArchiveRequestDescriptor instead')
const DeleteArchiveRequest$json = {
  '1': 'DeleteArchiveRequest',
  '2': [
    {'1': 'archivename', '3': 88048487, '4': 1, '5': 9, '10': 'archivename'},
  ],
};

/// Descriptor for `DeleteArchiveRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteArchiveRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVBcmNoaXZlUmVxdWVzdBIjCgthcmNoaXZlbmFtZRjnhv4pIAEoCVILYXJjaGl2ZW'
    '5hbWU=');

@$core.Deprecated('Use deleteArchiveResponseDescriptor instead')
const DeleteArchiveResponse$json = {
  '1': 'DeleteArchiveResponse',
};

/// Descriptor for `DeleteArchiveResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteArchiveResponseDescriptor =
    $convert.base64Decode('ChVEZWxldGVBcmNoaXZlUmVzcG9uc2U=');

@$core.Deprecated('Use deleteConnectionRequestDescriptor instead')
const DeleteConnectionRequest$json = {
  '1': 'DeleteConnectionRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteConnectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConnectionRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVDb25uZWN0aW9uUmVxdWVzdBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use deleteConnectionResponseDescriptor instead')
const DeleteConnectionResponse$json = {
  '1': 'DeleteConnectionResponse',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastauthorizedtime',
      '3': 250638066,
      '4': 1,
      '5': 9,
      '10': 'lastauthorizedtime'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `DeleteConnectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConnectionResponseDescriptor = $convert.base64Decode(
    'ChhEZWxldGVDb25uZWN0aW9uUmVzcG9uc2USJwoNY29ubmVjdGlvbmFybhjBj7xZIAEoCVINY2'
    '9ubmVjdGlvbmFybhJFCg9jb25uZWN0aW9uc3RhdGUY2/rlwAEgASgOMhcuZXZlbnRzLkNvbm5l'
    'Y3Rpb25TdGF0ZVIPY29ubmVjdGlvbnN0YXRlEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3'
    'JlYXRpb250aW1lEjEKEmxhc3RhdXRob3JpemVkdGltZRjy3cF3IAEoCVISbGFzdGF1dGhvcml6'
    'ZWR0aW1lEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWU=');

@$core.Deprecated('Use deleteEventBusRequestDescriptor instead')
const DeleteEventBusRequest$json = {
  '1': 'DeleteEventBusRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteEventBusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEventBusRequestDescriptor =
    $convert.base64Decode(
        'ChVEZWxldGVFdmVudEJ1c1JlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use deletePartnerEventSourceRequestDescriptor instead')
const DeletePartnerEventSourceRequest$json = {
  '1': 'DeletePartnerEventSourceRequest',
  '2': [
    {'1': 'account', '3': 435725053, '4': 1, '5': 9, '10': 'account'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeletePartnerEventSourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePartnerEventSourceRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVQYXJ0bmVyRXZlbnRTb3VyY2VSZXF1ZXN0EhwKB2FjY291bnQY/cXizwEgASgJUg'
        'dhY2NvdW50EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteRuleRequestDescriptor instead')
const DeleteRuleRequest$json = {
  '1': 'DeleteRuleRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'force', '3': 526711333, '4': 1, '5': 8, '10': 'force'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteRuleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRuleRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVSdWxlUmVxdWVzdBImCgxldmVudGJ1c25hbWUYvZnr1QEgASgJUgxldmVudGJ1c2'
    '5hbWUSGAoFZm9yY2UYpfST+wEgASgIUgVmb3JjZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use describeApiDestinationRequestDescriptor instead')
const DescribeApiDestinationRequest$json = {
  '1': 'DescribeApiDestinationRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeApiDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeApiDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZUFwaURlc3RpbmF0aW9uUmVxdWVzdBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use describeApiDestinationResponseDescriptor instead')
const DescribeApiDestinationResponse$json = {
  '1': 'DescribeApiDestinationResponse',
  '2': [
    {
      '1': 'apidestinationarn',
      '3': 90996885,
      '4': 1,
      '5': 9,
      '10': 'apidestinationarn'
    },
    {
      '1': 'apidestinationstate',
      '3': 13153343,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationState',
      '10': 'apidestinationstate'
    },
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'invocationendpoint',
      '3': 411800759,
      '4': 1,
      '5': 9,
      '10': 'invocationendpoint'
    },
    {
      '1': 'invocationratelimitpersecond',
      '3': 295327816,
      '4': 1,
      '5': 5,
      '10': 'invocationratelimitpersecond'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeApiDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeApiDestinationResponseDescriptor = $convert.base64Decode(
    'Ch5EZXNjcmliZUFwaURlc3RpbmF0aW9uUmVzcG9uc2USLwoRYXBpZGVzdGluYXRpb25hcm4YlY'
    'GyKyABKAlSEWFwaWRlc3RpbmF0aW9uYXJuElAKE2FwaWRlc3RpbmF0aW9uc3RhdGUYv+iiBiAB'
    'KA4yGy5ldmVudHMuQXBpRGVzdGluYXRpb25TdGF0ZVITYXBpZGVzdGluYXRpb25zdGF0ZRInCg'
    '1jb25uZWN0aW9uYXJuGMGPvFkgASgJUg1jb25uZWN0aW9uYXJuEiUKDGNyZWF0aW9udGltZRjm'
    'z6oxIAEoCVIMY3JlYXRpb250aW1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdG'
    'lvbhJECgpodHRwbWV0aG9kGNGM/L0BIAEoDjIgLmV2ZW50cy5BcGlEZXN0aW5hdGlvbkh0dHBN'
    'ZXRob2RSCmh0dHBtZXRob2QSMgoSaW52b2NhdGlvbmVuZHBvaW50GLeprsQBIAEoCVISaW52b2'
    'NhdGlvbmVuZHBvaW50EkYKHGludm9jYXRpb25yYXRlbGltaXRwZXJzZWNvbmQYyLDpjAEgASgF'
    'UhxpbnZvY2F0aW9ucmF0ZWxpbWl0cGVyc2Vjb25kEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cC'
    'ABKAlSEGxhc3Rtb2RpZmllZHRpbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use describeArchiveRequestDescriptor instead')
const DescribeArchiveRequest$json = {
  '1': 'DescribeArchiveRequest',
  '2': [
    {'1': 'archivename', '3': 88048487, '4': 1, '5': 9, '10': 'archivename'},
  ],
};

/// Descriptor for `DescribeArchiveRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeArchiveRequestDescriptor =
    $convert.base64Decode(
        'ChZEZXNjcmliZUFyY2hpdmVSZXF1ZXN0EiMKC2FyY2hpdmVuYW1lGOeG/ikgASgJUgthcmNoaX'
        'ZlbmFtZQ==');

@$core.Deprecated('Use describeArchiveResponseDescriptor instead')
const DescribeArchiveResponse$json = {
  '1': 'DescribeArchiveResponse',
  '2': [
    {'1': 'archivearn', '3': 56866685, '4': 1, '5': 9, '10': 'archivearn'},
    {'1': 'archivename', '3': 88048487, '4': 1, '5': 9, '10': 'archivename'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'eventcount', '3': 128612839, '4': 1, '5': 3, '10': 'eventcount'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'retentiondays',
      '3': 267894223,
      '4': 1,
      '5': 5,
      '10': 'retentiondays'
    },
    {'1': 'sizebytes', '3': 486677664, '4': 1, '5': 3, '10': 'sizebytes'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ArchiveState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `DescribeArchiveResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeArchiveResponseDescriptor = $convert.base64Decode(
    'ChdEZXNjcmliZUFyY2hpdmVSZXNwb25zZRIhCgphcmNoaXZlYXJuGP3ujhsgASgJUgphcmNoaX'
    'ZlYXJuEiMKC2FyY2hpdmVuYW1lGOeG/ikgASgJUgthcmNoaXZlbmFtZRIlCgxjcmVhdGlvbnRp'
    'bWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3'
    'JpcHRpb24SIQoKZXZlbnRjb3VudBjn86k9IAEoA1IKZXZlbnRjb3VudBIlCgxldmVudHBhdHRl'
    'cm4YgPeqbyABKAlSDGV2ZW50cGF0dGVybhIqCg5ldmVudHNvdXJjZWFybhjGyoqSASABKAlSDm'
    'V2ZW50c291cmNlYXJuEicKDXJldGVudGlvbmRheXMYz/vefyABKAVSDXJldGVudGlvbmRheXMS'
    'IAoJc2l6ZWJ5dGVzGKC5iOgBIAEoA1IJc2l6ZWJ5dGVzEi4KBXN0YXRlGJfJsu8BIAEoDjIULm'
    'V2ZW50cy5BcmNoaXZlU3RhdGVSBXN0YXRlEiQKC3N0YXRlcmVhc29uGPPVrbMBIAEoCVILc3Rh'
    'dGVyZWFzb24=');

@$core.Deprecated('Use describeConnectionRequestDescriptor instead')
const DescribeConnectionRequest$json = {
  '1': 'DescribeConnectionRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeConnectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeConnectionRequestDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZUNvbm5lY3Rpb25SZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use describeConnectionResponseDescriptor instead')
const DescribeConnectionResponse$json = {
  '1': 'DescribeConnectionResponse',
  '2': [
    {
      '1': 'authparameters',
      '3': 258276552,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionAuthResponseParameters',
      '10': 'authparameters'
    },
    {
      '1': 'authorizationtype',
      '3': 481976511,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionAuthorizationType',
      '10': 'authorizationtype'
    },
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'lastauthorizedtime',
      '3': 250638066,
      '4': 1,
      '5': 9,
      '10': 'lastauthorizedtime'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'secretarn', '3': 241012025, '4': 1, '5': 9, '10': 'secretarn'},
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `DescribeConnectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeConnectionResponseDescriptor = $convert.base64Decode(
    'ChpEZXNjcmliZUNvbm5lY3Rpb25SZXNwb25zZRJTCg5hdXRocGFyYW1ldGVycxjI+ZN7IAEoCz'
    'IoLmV2ZW50cy5Db25uZWN0aW9uQXV0aFJlc3BvbnNlUGFyYW1ldGVyc1IOYXV0aHBhcmFtZXRl'
    'cnMSVQoRYXV0aG9yaXphdGlvbnR5cGUYv8Hp5QEgASgOMiMuZXZlbnRzLkNvbm5lY3Rpb25BdX'
    'Rob3JpemF0aW9uVHlwZVIRYXV0aG9yaXphdGlvbnR5cGUSJwoNY29ubmVjdGlvbmFybhjBj7xZ'
    'IAEoCVINY29ubmVjdGlvbmFybhJFCg9jb25uZWN0aW9uc3RhdGUY2/rlwAEgASgOMhcuZXZlbn'
    'RzLkNvbm5lY3Rpb25TdGF0ZVIPY29ubmVjdGlvbnN0YXRlEiUKDGNyZWF0aW9udGltZRjmz6ox'
    'IAEoCVIMY3JlYXRpb250aW1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbh'
    'IxChJsYXN0YXV0aG9yaXplZHRpbWUY8t3BdyABKAlSEmxhc3RhdXRob3JpemVkdGltZRItChBs'
    'YXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUhBsYXN0bW9kaWZpZWR0aW1lEhUKBG5hbWUYh+aBfy'
    'ABKAlSBG5hbWUSHwoJc2VjcmV0YXJuGLma9nIgASgJUglzZWNyZXRhcm4SJAoLc3RhdGVyZWFz'
    'b24Y89WtswEgASgJUgtzdGF0ZXJlYXNvbg==');

@$core.Deprecated('Use describeEventBusRequestDescriptor instead')
const DescribeEventBusRequest$json = {
  '1': 'DescribeEventBusRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeEventBusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEventBusRequestDescriptor =
    $convert.base64Decode(
        'ChdEZXNjcmliZUV2ZW50QnVzUmVxdWVzdBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use describeEventBusResponseDescriptor instead')
const DescribeEventBusResponse$json = {
  '1': 'DescribeEventBusResponse',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
  ],
};

/// Descriptor for `DescribeEventBusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEventBusResponseDescriptor =
    $convert.base64Decode(
        'ChhEZXNjcmliZUV2ZW50QnVzUmVzcG9uc2USFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEhUKBG5hbW'
        'UYh+aBfyABKAlSBG5hbWUSGgoGcG9saWN5GKDv8OABIAEoCVIGcG9saWN5');

@$core.Deprecated('Use describeEventSourceRequestDescriptor instead')
const DescribeEventSourceRequest$json = {
  '1': 'DescribeEventSourceRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeEventSourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEventSourceRequestDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZUV2ZW50U291cmNlUmVxdWVzdBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use describeEventSourceResponseDescriptor instead')
const DescribeEventSourceResponse$json = {
  '1': 'DescribeEventSourceResponse',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdby', '3': 73491847, '4': 1, '5': 9, '10': 'createdby'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'expirationtime',
      '3': 93473378,
      '4': 1,
      '5': 9,
      '10': 'expirationtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.EventSourceState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `DescribeEventSourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEventSourceResponseDescriptor = $convert.base64Decode(
    'ChtEZXNjcmliZUV2ZW50U291cmNlUmVzcG9uc2USFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEh8KCW'
    'NyZWF0ZWRieRiHy4UjIAEoCVIJY3JlYXRlZGJ5EiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIM'
    'Y3JlYXRpb250aW1lEikKDmV4cGlyYXRpb250aW1lGOKUySwgASgJUg5leHBpcmF0aW9udGltZR'
    'IVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjIKBXN0YXRlGJfJsu8BIAEoDjIYLmV2ZW50cy5FdmVu'
    'dFNvdXJjZVN0YXRlUgVzdGF0ZQ==');

@$core.Deprecated('Use describePartnerEventSourceRequestDescriptor instead')
const DescribePartnerEventSourceRequest$json = {
  '1': 'DescribePartnerEventSourceRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribePartnerEventSourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePartnerEventSourceRequestDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZVBhcnRuZXJFdmVudFNvdXJjZVJlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbm'
        'FtZQ==');

@$core.Deprecated('Use describePartnerEventSourceResponseDescriptor instead')
const DescribePartnerEventSourceResponse$json = {
  '1': 'DescribePartnerEventSourceResponse',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribePartnerEventSourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePartnerEventSourceResponseDescriptor =
    $convert.base64Decode(
        'CiJEZXNjcmliZVBhcnRuZXJFdmVudFNvdXJjZVJlc3BvbnNlEhQKA2Fybhidm+2/ASABKAlSA2'
        'FybhIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use describeReplayRequestDescriptor instead')
const DescribeReplayRequest$json = {
  '1': 'DescribeReplayRequest',
  '2': [
    {'1': 'replayname', '3': 442173850, '4': 1, '5': 9, '10': 'replayname'},
  ],
};

/// Descriptor for `DescribeReplayRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeReplayRequestDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZVJlcGxheVJlcXVlc3QSIgoKcmVwbGF5bmFtZRiak+zSASABKAlSCnJlcGxheW'
    '5hbWU=');

@$core.Deprecated('Use describeReplayResponseDescriptor instead')
const DescribeReplayResponse$json = {
  '1': 'DescribeReplayResponse',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.events.ReplayDestination',
      '10': 'destination'
    },
    {'1': 'eventendtime', '3': 30229582, '4': 1, '5': 9, '10': 'eventendtime'},
    {
      '1': 'eventlastreplayedtime',
      '3': 7759073,
      '4': 1,
      '5': 9,
      '10': 'eventlastreplayedtime'
    },
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'eventstarttime',
      '3': 103680757,
      '4': 1,
      '5': 9,
      '10': 'eventstarttime'
    },
    {'1': 'replayarn', '3': 361946526, '4': 1, '5': 9, '10': 'replayarn'},
    {
      '1': 'replayendtime',
      '3': 304261199,
      '4': 1,
      '5': 9,
      '10': 'replayendtime'
    },
    {'1': 'replayname', '3': 442173850, '4': 1, '5': 9, '10': 'replayname'},
    {
      '1': 'replaystarttime',
      '3': 503580492,
      '4': 1,
      '5': 9,
      '10': 'replaystarttime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ReplayState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `DescribeReplayResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeReplayResponseDescriptor = $convert.base64Decode(
    'ChZEZXNjcmliZVJlcGxheVJlc3BvbnNlEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcm'
    'lwdGlvbhI/CgtkZXN0aW5hdGlvbhjgkpDaASABKAsyGS5ldmVudHMuUmVwbGF5RGVzdGluYXRp'
    'b25SC2Rlc3RpbmF0aW9uEiUKDGV2ZW50ZW5kdGltZRjOiLUOIAEoCVIMZXZlbnRlbmR0aW1lEj'
    'cKFWV2ZW50bGFzdHJlcGxheWVkdGltZRjhydkDIAEoCVIVZXZlbnRsYXN0cmVwbGF5ZWR0aW1l'
    'EioKDmV2ZW50c291cmNlYXJuGMbKipIBIAEoCVIOZXZlbnRzb3VyY2Vhcm4SKQoOZXZlbnRzdG'
    'FydHRpbWUY9ZW4MSABKAlSDmV2ZW50c3RhcnR0aW1lEiAKCXJlcGxheWFybhieu8usASABKAlS'
    'CXJlcGxheWFybhIoCg1yZXBsYXllbmR0aW1lGM/QipEBIAEoCVINcmVwbGF5ZW5kdGltZRIiCg'
    'pyZXBsYXluYW1lGJqT7NIBIAEoCVIKcmVwbGF5bmFtZRIsCg9yZXBsYXlzdGFydHRpbWUYzI6Q'
    '8AEgASgJUg9yZXBsYXlzdGFydHRpbWUSLQoFc3RhdGUYl8my7wEgASgOMhMuZXZlbnRzLlJlcG'
    'xheVN0YXRlUgVzdGF0ZRIkCgtzdGF0ZXJlYXNvbhjz1a2zASABKAlSC3N0YXRlcmVhc29u');

@$core.Deprecated('Use describeRuleRequestDescriptor instead')
const DescribeRuleRequest$json = {
  '1': 'DescribeRuleRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeRuleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeRuleRequestDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZVJ1bGVSZXF1ZXN0EiYKDGV2ZW50YnVzbmFtZRi9mevVASABKAlSDGV2ZW50Yn'
    'VzbmFtZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use describeRuleResponseDescriptor instead')
const DescribeRuleResponse$json = {
  '1': 'DescribeRuleResponse',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdby', '3': 73491847, '4': 1, '5': 9, '10': 'createdby'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
    {'1': 'managedby', '3': 455511232, '4': 1, '5': 9, '10': 'managedby'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.RuleState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `DescribeRuleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeRuleResponseDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZVJ1bGVSZXNwb25zZRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SHwoJY3JlYXRlZG'
    'J5GIfLhSMgASgJUgljcmVhdGVkYnkSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0'
    'aW9uEiYKDGV2ZW50YnVzbmFtZRi9mevVASABKAlSDGV2ZW50YnVzbmFtZRIlCgxldmVudHBhdH'
    'Rlcm4YgPeqbyABKAlSDGV2ZW50cGF0dGVybhIgCgltYW5hZ2VkYnkYwJma2QEgASgJUgltYW5h'
    'Z2VkYnkSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIcCgdyb2xlYXJuGIH455kBIAEoCVIHcm9sZW'
    'FybhIyChJzY2hlZHVsZWV4cHJlc3Npb24Y/5Hb1AEgASgJUhJzY2hlZHVsZWV4cHJlc3Npb24S'
    'KwoFc3RhdGUYl8my7wEgASgOMhEuZXZlbnRzLlJ1bGVTdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use disableRuleRequestDescriptor instead')
const DisableRuleRequest$json = {
  '1': 'DisableRuleRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DisableRuleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableRuleRequestDescriptor = $convert.base64Decode(
    'ChJEaXNhYmxlUnVsZVJlcXVlc3QSJgoMZXZlbnRidXNuYW1lGL2Z69UBIAEoCVIMZXZlbnRidX'
    'NuYW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use ecsParametersDescriptor instead')
const EcsParameters$json = {
  '1': 'EcsParameters',
  '2': [
    {
      '1': 'capacityproviderstrategy',
      '3': 273957206,
      '4': 3,
      '5': 11,
      '6': '.events.CapacityProviderStrategyItem',
      '10': 'capacityproviderstrategy'
    },
    {
      '1': 'enableecsmanagedtags',
      '3': 146161174,
      '4': 1,
      '5': 8,
      '10': 'enableecsmanagedtags'
    },
    {
      '1': 'enableexecutecommand',
      '3': 451374779,
      '4': 1,
      '5': 8,
      '10': 'enableexecutecommand'
    },
    {'1': 'group', '3': 91525165, '4': 1, '5': 9, '10': 'group'},
    {
      '1': 'launchtype',
      '3': 184333335,
      '4': 1,
      '5': 14,
      '6': '.events.LaunchType',
      '10': 'launchtype'
    },
    {
      '1': 'networkconfiguration',
      '3': 240088634,
      '4': 1,
      '5': 11,
      '6': '.events.NetworkConfiguration',
      '10': 'networkconfiguration'
    },
    {
      '1': 'placementconstraints',
      '3': 248464365,
      '4': 3,
      '5': 11,
      '6': '.events.PlacementConstraint',
      '10': 'placementconstraints'
    },
    {
      '1': 'placementstrategy',
      '3': 25036678,
      '4': 3,
      '5': 11,
      '6': '.events.PlacementStrategy',
      '10': 'placementstrategy'
    },
    {
      '1': 'platformversion',
      '3': 139924287,
      '4': 1,
      '5': 9,
      '10': 'platformversion'
    },
    {
      '1': 'propagatetags',
      '3': 405557622,
      '4': 1,
      '5': 14,
      '6': '.events.PropagateTags',
      '10': 'propagatetags'
    },
    {'1': 'referenceid', '3': 291739032, '4': 1, '5': 9, '10': 'referenceid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.events.Tag',
      '10': 'tags'
    },
    {'1': 'taskcount', '3': 398407508, '4': 1, '5': 5, '10': 'taskcount'},
    {
      '1': 'taskdefinitionarn',
      '3': 82234477,
      '4': 1,
      '5': 9,
      '10': 'taskdefinitionarn'
    },
  ],
};

/// Descriptor for `EcsParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ecsParametersDescriptor = $convert.base64Decode(
    'Cg1FY3NQYXJhbWV0ZXJzEmQKGGNhcGFjaXR5cHJvdmlkZXJzdHJhdGVneRjWgtGCASADKAsyJC'
    '5ldmVudHMuQ2FwYWNpdHlQcm92aWRlclN0cmF0ZWd5SXRlbVIYY2FwYWNpdHlwcm92aWRlcnN0'
    'cmF0ZWd5EjUKFGVuYWJsZWVjc21hbmFnZWR0YWdzGJb82EUgASgIUhRlbmFibGVlY3NtYW5hZ2'
    'VkdGFncxI2ChRlbmFibGVleGVjdXRlY29tbWFuZBi73Z3XASABKAhSFGVuYWJsZWV4ZWN1dGVj'
    'b21tYW5kEhcKBWdyb3VwGK2g0isgASgJUgVncm91cBI1CgpsYXVuY2h0eXBlGJfo8lcgASgOMh'
    'IuZXZlbnRzLkxhdW5jaFR5cGVSCmxhdW5jaHR5cGUSUwoUbmV0d29ya2NvbmZpZ3VyYXRpb24Y'
    'uuy9ciABKAsyHC5ldmVudHMuTmV0d29ya0NvbmZpZ3VyYXRpb25SFG5ldHdvcmtjb25maWd1cm'
    'F0aW9uElIKFHBsYWNlbWVudGNvbnN0cmFpbnRzGO2HvXYgAygLMhsuZXZlbnRzLlBsYWNlbWVu'
    'dENvbnN0cmFpbnRSFHBsYWNlbWVudGNvbnN0cmFpbnRzEkoKEXBsYWNlbWVudHN0cmF0ZWd5GI'
    'aP+AsgAygLMhkuZXZlbnRzLlBsYWNlbWVudFN0cmF0ZWd5UhFwbGFjZW1lbnRzdHJhdGVneRIr'
    'Cg9wbGF0Zm9ybXZlcnNpb24Yv6bcQiABKAlSD3BsYXRmb3JtdmVyc2lvbhI/Cg1wcm9wYWdhdG'
    'V0YWdzGPaiscEBIAEoDjIVLmV2ZW50cy5Qcm9wYWdhdGVUYWdzUg1wcm9wYWdhdGV0YWdzEiQK'
    'C3JlZmVyZW5jZWlkGJirjosBIAEoCVILcmVmZXJlbmNlaWQSIwoEdGFncxjBwfa1ASADKAsyCy'
    '5ldmVudHMuVGFnUgR0YWdzEiAKCXRhc2tjb3VudBjU7vy9ASABKAVSCXRhc2tjb3VudBIvChF0'
    'YXNrZGVmaW5pdGlvbmFybhjtmJsnIAEoCVIRdGFza2RlZmluaXRpb25hcm4=');

@$core.Deprecated('Use enableRuleRequestDescriptor instead')
const EnableRuleRequest$json = {
  '1': 'EnableRuleRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `EnableRuleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableRuleRequestDescriptor = $convert.base64Decode(
    'ChFFbmFibGVSdWxlUmVxdWVzdBImCgxldmVudGJ1c25hbWUYvZnr1QEgASgJUgxldmVudGJ1c2'
    '5hbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use eventBusDescriptor instead')
const EventBus$json = {
  '1': 'EventBus',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
  ],
};

/// Descriptor for `EventBus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventBusDescriptor = $convert.base64Decode(
    'CghFdmVudEJ1cxIUCgNhcm4YnZvtvwEgASgJUgNhcm4SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZR'
    'IaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2xpY3k=');

@$core.Deprecated('Use eventSourceDescriptor instead')
const EventSource$json = {
  '1': 'EventSource',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdby', '3': 73491847, '4': 1, '5': 9, '10': 'createdby'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'expirationtime',
      '3': 93473378,
      '4': 1,
      '5': 9,
      '10': 'expirationtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.EventSourceState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `EventSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventSourceDescriptor = $convert.base64Decode(
    'CgtFdmVudFNvdXJjZRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SHwoJY3JlYXRlZGJ5GIfLhSMgAS'
    'gJUgljcmVhdGVkYnkSJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbnRpbWUSKQoO'
    'ZXhwaXJhdGlvbnRpbWUY4pTJLCABKAlSDmV4cGlyYXRpb250aW1lEhUKBG5hbWUYh+aBfyABKA'
    'lSBG5hbWUSMgoFc3RhdGUYl8my7wEgASgOMhguZXZlbnRzLkV2ZW50U291cmNlU3RhdGVSBXN0'
    'YXRl');

@$core.Deprecated('Use httpParametersDescriptor instead')
const HttpParameters$json = {
  '1': 'HttpParameters',
  '2': [
    {
      '1': 'headerparameters',
      '3': 148944757,
      '4': 3,
      '5': 11,
      '6': '.events.HttpParameters.HeaderparametersEntry',
      '10': 'headerparameters'
    },
    {
      '1': 'pathparametervalues',
      '3': 489440856,
      '4': 3,
      '5': 9,
      '10': 'pathparametervalues'
    },
    {
      '1': 'querystringparameters',
      '3': 258002263,
      '4': 3,
      '5': 11,
      '6': '.events.HttpParameters.QuerystringparametersEntry',
      '10': 'querystringparameters'
    },
  ],
  '3': [
    HttpParameters_HeaderparametersEntry$json,
    HttpParameters_QuerystringparametersEntry$json
  ],
};

@$core.Deprecated('Use httpParametersDescriptor instead')
const HttpParameters_HeaderparametersEntry$json = {
  '1': 'HeaderparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use httpParametersDescriptor instead')
const HttpParameters_QuerystringparametersEntry$json = {
  '1': 'QuerystringparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `HttpParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List httpParametersDescriptor = $convert.base64Decode(
    'Cg5IdHRwUGFyYW1ldGVycxJbChBoZWFkZXJwYXJhbWV0ZXJzGPXugkcgAygLMiwuZXZlbnRzLk'
    'h0dHBQYXJhbWV0ZXJzLkhlYWRlcnBhcmFtZXRlcnNFbnRyeVIQaGVhZGVycGFyYW1ldGVycxI0'
    'ChNwYXRocGFyYW1ldGVydmFsdWVzGNiMsekBIAMoCVITcGF0aHBhcmFtZXRlcnZhbHVlcxJqCh'
    'VxdWVyeXN0cmluZ3BhcmFtZXRlcnMY15qDeyADKAsyMS5ldmVudHMuSHR0cFBhcmFtZXRlcnMu'
    'UXVlcnlzdHJpbmdwYXJhbWV0ZXJzRW50cnlSFXF1ZXJ5c3RyaW5ncGFyYW1ldGVycxpDChVIZW'
    'FkZXJwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZh'
    'bHVlOgI4ARpIChpRdWVyeXN0cmluZ3BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleR'
    'IUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use illegalStatusExceptionDescriptor instead')
const IllegalStatusException$json = {
  '1': 'IllegalStatusException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IllegalStatusException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List illegalStatusExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbGxlZ2FsU3RhdHVzRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use inputTransformerDescriptor instead')
const InputTransformer$json = {
  '1': 'InputTransformer',
  '2': [
    {
      '1': 'inputpathsmap',
      '3': 154989658,
      '4': 3,
      '5': 11,
      '6': '.events.InputTransformer.InputpathsmapEntry',
      '10': 'inputpathsmap'
    },
    {
      '1': 'inputtemplate',
      '3': 505971626,
      '4': 1,
      '5': 9,
      '10': 'inputtemplate'
    },
  ],
  '3': [InputTransformer_InputpathsmapEntry$json],
};

@$core.Deprecated('Use inputTransformerDescriptor instead')
const InputTransformer_InputpathsmapEntry$json = {
  '1': 'InputpathsmapEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InputTransformer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inputTransformerDescriptor = $convert.base64Decode(
    'ChBJbnB1dFRyYW5zZm9ybWVyElQKDWlucHV0cGF0aHNtYXAY2ujzSSADKAsyKy5ldmVudHMuSW'
    '5wdXRUcmFuc2Zvcm1lci5JbnB1dHBhdGhzbWFwRW50cnlSDWlucHV0cGF0aHNtYXASKAoNaW5w'
    'dXR0ZW1wbGF0ZRiqh6LxASABKAlSDWlucHV0dGVtcGxhdGUaQAoSSW5wdXRwYXRoc21hcEVudH'
    'J5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use internalExceptionDescriptor instead')
const InternalException$json = {
  '1': 'InternalException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalExceptionDescriptor = $convert.base64Decode(
    'ChFJbnRlcm5hbEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidEventPatternExceptionDescriptor instead')
const InvalidEventPatternException$json = {
  '1': 'InvalidEventPatternException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEventPatternException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEventPatternExceptionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkRXZlbnRQYXR0ZXJuRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidStateExceptionDescriptor instead')
const InvalidStateException$json = {
  '1': 'InvalidStateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidStateExceptionDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkU3RhdGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kinesisParametersDescriptor instead')
const KinesisParameters$json = {
  '1': 'KinesisParameters',
  '2': [
    {
      '1': 'partitionkeypath',
      '3': 458722216,
      '4': 1,
      '5': 9,
      '10': 'partitionkeypath'
    },
  ],
};

/// Descriptor for `KinesisParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisParametersDescriptor = $convert.base64Decode(
    'ChFLaW5lc2lzUGFyYW1ldGVycxIuChBwYXJ0aXRpb25rZXlwYXRoGKiX3toBIAEoCVIQcGFydG'
    'l0aW9ua2V5cGF0aA==');

@$core.Deprecated('Use limitExceededExceptionDescriptor instead')
const LimitExceededException$json = {
  '1': 'LimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChZMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use listApiDestinationsRequestDescriptor instead')
const ListApiDestinationsRequest$json = {
  '1': 'ListApiDestinationsRequest',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListApiDestinationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listApiDestinationsRequestDescriptor = $convert.base64Decode(
    'ChpMaXN0QXBpRGVzdGluYXRpb25zUmVxdWVzdBInCg1jb25uZWN0aW9uYXJuGMGPvFkgASgJUg'
    '1jb25uZWN0aW9uYXJuEhgKBWxpbWl0GNWV2cQBIAEoBVIFbGltaXQSIgoKbmFtZXByZWZpeBib'
    '87ysASABKAlSCm5hbWVwcmVmaXgSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listApiDestinationsResponseDescriptor instead')
const ListApiDestinationsResponse$json = {
  '1': 'ListApiDestinationsResponse',
  '2': [
    {
      '1': 'apidestinations',
      '3': 321414539,
      '4': 3,
      '5': 11,
      '6': '.events.ApiDestination',
      '10': 'apidestinations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListApiDestinationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listApiDestinationsResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0QXBpRGVzdGluYXRpb25zUmVzcG9uc2USRAoPYXBpZGVzdGluYXRpb25zGIvLoZkBIA'
        'MoCzIWLmV2ZW50cy5BcGlEZXN0aW5hdGlvblIPYXBpZGVzdGluYXRpb25zEh8KCW5leHR0b2tl'
        'bhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listArchivesRequestDescriptor instead')
const ListArchivesRequest$json = {
  '1': 'ListArchivesRequest',
  '2': [
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ArchiveState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `ListArchivesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listArchivesRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0QXJjaGl2ZXNSZXF1ZXN0EioKDmV2ZW50c291cmNlYXJuGMbKipIBIAEoCVIOZXZlbn'
    'Rzb3VyY2Vhcm4SGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIiCgpuYW1lcHJlZml4GJvzvKwB'
    'IAEoCVIKbmFtZXByZWZpeBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIuCgVzdG'
    'F0ZRiXybLvASABKA4yFC5ldmVudHMuQXJjaGl2ZVN0YXRlUgVzdGF0ZQ==');

@$core.Deprecated('Use listArchivesResponseDescriptor instead')
const ListArchivesResponse$json = {
  '1': 'ListArchivesResponse',
  '2': [
    {
      '1': 'archives',
      '3': 73167779,
      '4': 3,
      '5': 11,
      '6': '.events.Archive',
      '10': 'archives'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListArchivesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listArchivesResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0QXJjaGl2ZXNSZXNwb25zZRIuCghhcmNoaXZlcxij5/EiIAMoCzIPLmV2ZW50cy5Bcm'
    'NoaXZlUghhcmNoaXZlcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listConnectionsRequestDescriptor instead')
const ListConnectionsRequest$json = {
  '1': 'ListConnectionsRequest',
  '2': [
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListConnectionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConnectionsRequestDescriptor = $convert.base64Decode(
    'ChZMaXN0Q29ubmVjdGlvbnNSZXF1ZXN0EkUKD2Nvbm5lY3Rpb25zdGF0ZRjb+uXAASABKA4yFy'
    '5ldmVudHMuQ29ubmVjdGlvblN0YXRlUg9jb25uZWN0aW9uc3RhdGUSGAoFbGltaXQY1ZXZxAEg'
    'ASgFUgVsaW1pdBIiCgpuYW1lcHJlZml4GJvzvKwBIAEoCVIKbmFtZXByZWZpeBIfCgluZXh0dG'
    '9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listConnectionsResponseDescriptor instead')
const ListConnectionsResponse$json = {
  '1': 'ListConnectionsResponse',
  '2': [
    {
      '1': 'connections',
      '3': 98252351,
      '4': 3,
      '5': 11,
      '6': '.events.Connection',
      '10': 'connections'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListConnectionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConnectionsResponseDescriptor = $convert.base64Decode(
    'ChdMaXN0Q29ubmVjdGlvbnNSZXNwb25zZRI3Cgtjb25uZWN0aW9ucxi/7OwuIAMoCzISLmV2ZW'
    '50cy5Db25uZWN0aW9uUgtjb25uZWN0aW9ucxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0'
    'b2tlbg==');

@$core.Deprecated('Use listEventBusesRequestDescriptor instead')
const ListEventBusesRequest$json = {
  '1': 'ListEventBusesRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEventBusesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventBusesRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0RXZlbnRCdXNlc1JlcXVlc3QSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIiCgpuYW'
    '1lcHJlZml4GJvzvKwBIAEoCVIKbmFtZXByZWZpeBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5l'
    'eHR0b2tlbg==');

@$core.Deprecated('Use listEventBusesResponseDescriptor instead')
const ListEventBusesResponse$json = {
  '1': 'ListEventBusesResponse',
  '2': [
    {
      '1': 'eventbuses',
      '3': 125356616,
      '4': 3,
      '5': 11,
      '6': '.events.EventBus',
      '10': 'eventbuses'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEventBusesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventBusesResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0RXZlbnRCdXNlc1Jlc3BvbnNlEjMKCmV2ZW50YnVzZXMYyJTjOyADKAsyEC5ldmVudH'
    'MuRXZlbnRCdXNSCmV2ZW50YnVzZXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listEventSourcesRequestDescriptor instead')
const ListEventSourcesRequest$json = {
  '1': 'ListEventSourcesRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEventSourcesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventSourcesRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0RXZlbnRTb3VyY2VzUmVxdWVzdBIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0EiIKCm'
    '5hbWVwcmVmaXgYm/O8rAEgASgJUgpuYW1lcHJlZml4Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJ'
    'bmV4dHRva2Vu');

@$core.Deprecated('Use listEventSourcesResponseDescriptor instead')
const ListEventSourcesResponse$json = {
  '1': 'ListEventSourcesResponse',
  '2': [
    {
      '1': 'eventsources',
      '3': 368674668,
      '4': 3,
      '5': 11,
      '6': '.events.EventSource',
      '10': 'eventsources'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEventSourcesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventSourcesResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0RXZlbnRTb3VyY2VzUmVzcG9uc2USOwoMZXZlbnRzb3VyY2VzGOyO5q8BIAMoCzITLm'
    'V2ZW50cy5FdmVudFNvdXJjZVIMZXZlbnRzb3VyY2VzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJ'
    'bmV4dHRva2Vu');

@$core.Deprecated('Use listPartnerEventSourceAccountsRequestDescriptor instead')
const ListPartnerEventSourceAccountsRequest$json = {
  '1': 'ListPartnerEventSourceAccountsRequest',
  '2': [
    {
      '1': 'eventsourcename',
      '3': 427669794,
      '4': 1,
      '5': 9,
      '10': 'eventsourcename'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListPartnerEventSourceAccountsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPartnerEventSourceAccountsRequestDescriptor =
    $convert.base64Decode(
        'CiVMaXN0UGFydG5lckV2ZW50U291cmNlQWNjb3VudHNSZXF1ZXN0EiwKD2V2ZW50c291cmNlbm'
        'FtZRii8vbLASABKAlSD2V2ZW50c291cmNlbmFtZRIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0'
        'Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core
    .Deprecated('Use listPartnerEventSourceAccountsResponseDescriptor instead')
const ListPartnerEventSourceAccountsResponse$json = {
  '1': 'ListPartnerEventSourceAccountsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'partnereventsourceaccounts',
      '3': 334415295,
      '4': 3,
      '5': 11,
      '6': '.events.PartnerEventSourceAccount',
      '10': 'partnereventsourceaccounts'
    },
  ],
};

/// Descriptor for `ListPartnerEventSourceAccountsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPartnerEventSourceAccountsResponseDescriptor =
    $convert.base64Decode(
        'CiZMaXN0UGFydG5lckV2ZW50U291cmNlQWNjb3VudHNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/o'
        'S6ZyABKAlSCW5leHR0b2tlbhJlChpwYXJ0bmVyZXZlbnRzb3VyY2VhY2NvdW50cxi/i7ufASAD'
        'KAsyIS5ldmVudHMuUGFydG5lckV2ZW50U291cmNlQWNjb3VudFIacGFydG5lcmV2ZW50c291cm'
        'NlYWNjb3VudHM=');

@$core.Deprecated('Use listPartnerEventSourcesRequestDescriptor instead')
const ListPartnerEventSourcesRequest$json = {
  '1': 'ListPartnerEventSourcesRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListPartnerEventSourcesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPartnerEventSourcesRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0UGFydG5lckV2ZW50U291cmNlc1JlcXVlc3QSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW'
        '1pdBIiCgpuYW1lcHJlZml4GJvzvKwBIAEoCVIKbmFtZXByZWZpeBIfCgluZXh0dG9rZW4Y/oS6'
        'ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listPartnerEventSourcesResponseDescriptor instead')
const ListPartnerEventSourcesResponse$json = {
  '1': 'ListPartnerEventSourcesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'partnereventsources',
      '3': 502828176,
      '4': 3,
      '5': 11,
      '6': '.events.PartnerEventSource',
      '10': 'partnereventsources'
    },
  ],
};

/// Descriptor for `ListPartnerEventSourcesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPartnerEventSourcesResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0UGFydG5lckV2ZW50U291cmNlc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCV'
        'IJbmV4dHRva2VuElAKE3BhcnRuZXJldmVudHNvdXJjZXMYkJni7wEgAygLMhouZXZlbnRzLlBh'
        'cnRuZXJFdmVudFNvdXJjZVITcGFydG5lcmV2ZW50c291cmNlcw==');

@$core.Deprecated('Use listReplaysRequestDescriptor instead')
const ListReplaysRequest$json = {
  '1': 'ListReplaysRequest',
  '2': [
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ReplayState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `ListReplaysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReplaysRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0UmVwbGF5c1JlcXVlc3QSKgoOZXZlbnRzb3VyY2Vhcm4YxsqKkgEgASgJUg5ldmVudH'
    'NvdXJjZWFybhIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0EiIKCm5hbWVwcmVmaXgYm/O8rAEg'
    'ASgJUgpuYW1lcHJlZml4Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEi0KBXN0YX'
    'RlGJfJsu8BIAEoDjITLmV2ZW50cy5SZXBsYXlTdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use listReplaysResponseDescriptor instead')
const ListReplaysResponse$json = {
  '1': 'ListReplaysResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'replays',
      '3': 312313364,
      '4': 3,
      '5': 11,
      '6': '.events.Replay',
      '10': 'replays'
    },
  ],
};

/// Descriptor for `ListReplaysResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReplaysResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0UmVwbGF5c1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEi'
    'wKB3JlcGxheXMYlIz2lAEgAygLMg4uZXZlbnRzLlJlcGxheVIHcmVwbGF5cw==');

@$core.Deprecated('Use listRuleNamesByTargetRequestDescriptor instead')
const ListRuleNamesByTargetRequest$json = {
  '1': 'ListRuleNamesByTargetRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'targetarn', '3': 217664144, '4': 1, '5': 9, '10': 'targetarn'},
  ],
};

/// Descriptor for `ListRuleNamesByTargetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRuleNamesByTargetRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0UnVsZU5hbWVzQnlUYXJnZXRSZXF1ZXN0EiYKDGV2ZW50YnVzbmFtZRi9mevVASABKA'
        'lSDGV2ZW50YnVzbmFtZRIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0Eh8KCW5leHR0b2tlbhj+'
        'hLpnIAEoCVIJbmV4dHRva2VuEh8KCXRhcmdldGFybhiQleVnIAEoCVIJdGFyZ2V0YXJu');

@$core.Deprecated('Use listRuleNamesByTargetResponseDescriptor instead')
const ListRuleNamesByTargetResponse$json = {
  '1': 'ListRuleNamesByTargetResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'rulenames', '3': 267949170, '4': 3, '5': 9, '10': 'rulenames'},
  ],
};

/// Descriptor for `ListRuleNamesByTargetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRuleNamesByTargetResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0UnVsZU5hbWVzQnlUYXJnZXRSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
        '5leHR0b2tlbhIfCglydWxlbmFtZXMY8qjifyADKAlSCXJ1bGVuYW1lcw==');

@$core.Deprecated('Use listRulesRequestDescriptor instead')
const ListRulesRequest$json = {
  '1': 'ListRulesRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListRulesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRulesRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0UnVsZXNSZXF1ZXN0EiYKDGV2ZW50YnVzbmFtZRi9mevVASABKAlSDGV2ZW50YnVzbm'
    'FtZRIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0EiIKCm5hbWVwcmVmaXgYm/O8rAEgASgJUgpu'
    'YW1lcHJlZml4Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listRulesResponseDescriptor instead')
const ListRulesResponse$json = {
  '1': 'ListRulesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.events.Rule',
      '10': 'rules'
    },
  ],
};

/// Descriptor for `ListRulesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRulesResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0UnVsZXNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIlCg'
    'VydWxlcxiB26wUIAMoCzIMLmV2ZW50cy5SdWxlUgVydWxlcw==');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhjtwJmwASABKAlSC3'
        'Jlc291cmNlYXJu');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.events.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USIwoEdGFncxjBwfa1ASADKAsyCy5ldmVudH'
        'MuVGFnUgR0YWdz');

@$core.Deprecated('Use listTargetsByRuleRequestDescriptor instead')
const ListTargetsByRuleRequest$json = {
  '1': 'ListTargetsByRuleRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'rule', '3': 475696372, '4': 1, '5': 9, '10': 'rule'},
  ],
};

/// Descriptor for `ListTargetsByRuleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTargetsByRuleRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0VGFyZ2V0c0J5UnVsZVJlcXVlc3QSJgoMZXZlbnRidXNuYW1lGL2Z69UBIAEoCVIMZX'
    'ZlbnRidXNuYW1lEhgKBWxpbWl0GNWV2cQBIAEoBVIFbGltaXQSHwoJbmV4dHRva2VuGP6Eumcg'
    'ASgJUgluZXh0dG9rZW4SFgoEcnVsZRj0meriASABKAlSBHJ1bGU=');

@$core.Deprecated('Use listTargetsByRuleResponseDescriptor instead')
const ListTargetsByRuleResponse$json = {
  '1': 'ListTargetsByRuleResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.events.Target',
      '10': 'targets'
    },
  ],
};

/// Descriptor for `ListTargetsByRuleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTargetsByRuleResponseDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFyZ2V0c0J5UnVsZVJlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
        'Rva2VuEisKB3RhcmdldHMYgpuCfSADKAsyDi5ldmVudHMuVGFyZ2V0Ugd0YXJnZXRz');

@$core.Deprecated('Use managedRuleExceptionDescriptor instead')
const ManagedRuleException$json = {
  '1': 'ManagedRuleException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ManagedRuleException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedRuleExceptionDescriptor =
    $convert.base64Decode(
        'ChRNYW5hZ2VkUnVsZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use networkConfigurationDescriptor instead')
const NetworkConfiguration$json = {
  '1': 'NetworkConfiguration',
  '2': [
    {
      '1': 'awsvpcconfiguration',
      '3': 223852630,
      '4': 1,
      '5': 11,
      '6': '.events.AwsVpcConfiguration',
      '10': 'awsvpcconfiguration'
    },
  ],
};

/// Descriptor for `NetworkConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkConfigurationDescriptor = $convert.base64Decode(
    'ChROZXR3b3JrQ29uZmlndXJhdGlvbhJQChNhd3N2cGNjb25maWd1cmF0aW9uGNbw3mogASgLMh'
    'suZXZlbnRzLkF3c1ZwY0NvbmZpZ3VyYXRpb25SE2F3c3ZwY2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use operationDisabledExceptionDescriptor instead')
const OperationDisabledException$json = {
  '1': 'OperationDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OperationDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List operationDisabledExceptionDescriptor =
    $convert.base64Decode(
        'ChpPcGVyYXRpb25EaXNhYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use partnerEventSourceDescriptor instead')
const PartnerEventSource$json = {
  '1': 'PartnerEventSource',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `PartnerEventSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partnerEventSourceDescriptor = $convert.base64Decode(
    'ChJQYXJ0bmVyRXZlbnRTb3VyY2USFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEhUKBG5hbWUYh+aBfy'
    'ABKAlSBG5hbWU=');

@$core.Deprecated('Use partnerEventSourceAccountDescriptor instead')
const PartnerEventSourceAccount$json = {
  '1': 'PartnerEventSourceAccount',
  '2': [
    {'1': 'account', '3': 435725053, '4': 1, '5': 9, '10': 'account'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'expirationtime',
      '3': 93473378,
      '4': 1,
      '5': 9,
      '10': 'expirationtime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.EventSourceState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `PartnerEventSourceAccount`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partnerEventSourceAccountDescriptor = $convert.base64Decode(
    'ChlQYXJ0bmVyRXZlbnRTb3VyY2VBY2NvdW50EhwKB2FjY291bnQY/cXizwEgASgJUgdhY2NvdW'
    '50EiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYXRpb250aW1lEikKDmV4cGlyYXRpb250'
    'aW1lGOKUySwgASgJUg5leHBpcmF0aW9udGltZRIyCgVzdGF0ZRiXybLvASABKA4yGC5ldmVudH'
    'MuRXZlbnRTb3VyY2VTdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use placementConstraintDescriptor instead')
const PlacementConstraint$json = {
  '1': 'PlacementConstraint',
  '2': [
    {'1': 'expression', '3': 253079532, '4': 1, '5': 9, '10': 'expression'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.events.PlacementConstraintType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `PlacementConstraint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List placementConstraintDescriptor = $convert.base64Decode(
    'ChNQbGFjZW1lbnRDb25zdHJhaW50EiEKCmV4cHJlc3Npb24Y7N/WeCABKAlSCmV4cHJlc3Npb2'
    '4SNwoEdHlwZRjO4p+JASABKA4yHy5ldmVudHMuUGxhY2VtZW50Q29uc3RyYWludFR5cGVSBHR5'
    'cGU=');

@$core.Deprecated('Use placementStrategyDescriptor instead')
const PlacementStrategy$json = {
  '1': 'PlacementStrategy',
  '2': [
    {'1': 'field', '3': 125985384, '4': 1, '5': 9, '10': 'field'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.events.PlacementStrategyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `PlacementStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List placementStrategyDescriptor = $convert.base64Decode(
    'ChFQbGFjZW1lbnRTdHJhdGVneRIXCgVmaWVsZBjoxIk8IAEoCVIFZmllbGQSNQoEdHlwZRjO4p'
    '+JASABKA4yHS5ldmVudHMuUGxhY2VtZW50U3RyYXRlZ3lUeXBlUgR0eXBl');

@$core.Deprecated('Use policyLengthExceededExceptionDescriptor instead')
const PolicyLengthExceededException$json = {
  '1': 'PolicyLengthExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PolicyLengthExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyLengthExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1Qb2xpY3lMZW5ndGhFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use putEventsRequestDescriptor instead')
const PutEventsRequest$json = {
  '1': 'PutEventsRequest',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.events.PutEventsRequestEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `PutEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventsRequestDescriptor = $convert.base64Decode(
    'ChBQdXRFdmVudHNSZXF1ZXN0EjsKB2VudHJpZXMYlMWy5QEgAygLMh0uZXZlbnRzLlB1dEV2ZW'
    '50c1JlcXVlc3RFbnRyeVIHZW50cmllcw==');

@$core.Deprecated('Use putEventsRequestEntryDescriptor instead')
const PutEventsRequestEntry$json = {
  '1': 'PutEventsRequestEntry',
  '2': [
    {'1': 'detail', '3': 340030389, '4': 1, '5': 9, '10': 'detail'},
    {'1': 'detailtype', '3': 11758589, '4': 1, '5': 9, '10': 'detailtype'},
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'resources', '3': 358918291, '4': 3, '5': 9, '10': 'resources'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {'1': 'time', '3': 535094277, '4': 1, '5': 9, '10': 'time'},
    {'1': 'traceheader', '3': 303628800, '4': 1, '5': 9, '10': 'traceheader'},
  ],
};

/// Descriptor for `PutEventsRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventsRequestEntryDescriptor = $convert.base64Decode(
    'ChVQdXRFdmVudHNSZXF1ZXN0RW50cnkSGgoGZGV0YWlsGLXnkaIBIAEoCVIGZGV0YWlsEiEKCm'
    'RldGFpbHR5cGUY/dfNBSABKAlSCmRldGFpbHR5cGUSJgoMZXZlbnRidXNuYW1lGL2Z69UBIAEo'
    'CVIMZXZlbnRidXNuYW1lEiAKCXJlc291cmNlcxiT0ZKrASADKAlSCXJlc291cmNlcxIZCgZzb3'
    'VyY2UY+ceKDyABKAlSBnNvdXJjZRIWCgR0aW1lGIXIk/8BIAEoCVIEdGltZRIkCgt0cmFjZWhl'
    'YWRlchiAhOSQASABKAlSC3RyYWNlaGVhZGVy');

@$core.Deprecated('Use putEventsResponseDescriptor instead')
const PutEventsResponse$json = {
  '1': 'PutEventsResponse',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.events.PutEventsResultEntry',
      '10': 'entries'
    },
    {
      '1': 'failedentrycount',
      '3': 458519576,
      '4': 1,
      '5': 5,
      '10': 'failedentrycount'
    },
  ],
};

/// Descriptor for `PutEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventsResponseDescriptor = $convert.base64Decode(
    'ChFQdXRFdmVudHNSZXNwb25zZRI6CgdlbnRyaWVzGJTFsuUBIAMoCzIcLmV2ZW50cy5QdXRFdm'
    'VudHNSZXN1bHRFbnRyeVIHZW50cmllcxIuChBmYWlsZWRlbnRyeWNvdW50GJjo0doBIAEoBVIQ'
    'ZmFpbGVkZW50cnljb3VudA==');

@$core.Deprecated('Use putEventsResultEntryDescriptor instead')
const PutEventsResultEntry$json = {
  '1': 'PutEventsResultEntry',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
  ],
};

/// Descriptor for `PutEventsResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventsResultEntryDescriptor = $convert.base64Decode(
    'ChRQdXRFdmVudHNSZXN1bHRFbnRyeRIfCgllcnJvcmNvZGUYmdbDECABKAlSCWVycm9yY29kZR'
    'ImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USHAoHZXZlbnRpZBjTlt2z'
    'ASABKAlSB2V2ZW50aWQ=');

@$core.Deprecated('Use putPartnerEventsRequestDescriptor instead')
const PutPartnerEventsRequest$json = {
  '1': 'PutPartnerEventsRequest',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.events.PutPartnerEventsRequestEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `PutPartnerEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putPartnerEventsRequestDescriptor =
    $convert.base64Decode(
        'ChdQdXRQYXJ0bmVyRXZlbnRzUmVxdWVzdBJCCgdlbnRyaWVzGJTFsuUBIAMoCzIkLmV2ZW50cy'
        '5QdXRQYXJ0bmVyRXZlbnRzUmVxdWVzdEVudHJ5UgdlbnRyaWVz');

@$core.Deprecated('Use putPartnerEventsRequestEntryDescriptor instead')
const PutPartnerEventsRequestEntry$json = {
  '1': 'PutPartnerEventsRequestEntry',
  '2': [
    {'1': 'detail', '3': 340030389, '4': 1, '5': 9, '10': 'detail'},
    {'1': 'detailtype', '3': 11758589, '4': 1, '5': 9, '10': 'detailtype'},
    {'1': 'resources', '3': 358918291, '4': 3, '5': 9, '10': 'resources'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {'1': 'time', '3': 535094277, '4': 1, '5': 9, '10': 'time'},
  ],
};

/// Descriptor for `PutPartnerEventsRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putPartnerEventsRequestEntryDescriptor = $convert.base64Decode(
    'ChxQdXRQYXJ0bmVyRXZlbnRzUmVxdWVzdEVudHJ5EhoKBmRldGFpbBi155GiASABKAlSBmRldG'
    'FpbBIhCgpkZXRhaWx0eXBlGP3XzQUgASgJUgpkZXRhaWx0eXBlEiAKCXJlc291cmNlcxiT0ZKr'
    'ASADKAlSCXJlc291cmNlcxIZCgZzb3VyY2UY+ceKDyABKAlSBnNvdXJjZRIWCgR0aW1lGIXIk/'
    '8BIAEoCVIEdGltZQ==');

@$core.Deprecated('Use putPartnerEventsResponseDescriptor instead')
const PutPartnerEventsResponse$json = {
  '1': 'PutPartnerEventsResponse',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.events.PutPartnerEventsResultEntry',
      '10': 'entries'
    },
    {
      '1': 'failedentrycount',
      '3': 458519576,
      '4': 1,
      '5': 5,
      '10': 'failedentrycount'
    },
  ],
};

/// Descriptor for `PutPartnerEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putPartnerEventsResponseDescriptor = $convert.base64Decode(
    'ChhQdXRQYXJ0bmVyRXZlbnRzUmVzcG9uc2USQQoHZW50cmllcxiUxbLlASADKAsyIy5ldmVudH'
    'MuUHV0UGFydG5lckV2ZW50c1Jlc3VsdEVudHJ5UgdlbnRyaWVzEi4KEGZhaWxlZGVudHJ5Y291'
    'bnQYmOjR2gEgASgFUhBmYWlsZWRlbnRyeWNvdW50');

@$core.Deprecated('Use putPartnerEventsResultEntryDescriptor instead')
const PutPartnerEventsResultEntry$json = {
  '1': 'PutPartnerEventsResultEntry',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
  ],
};

/// Descriptor for `PutPartnerEventsResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putPartnerEventsResultEntryDescriptor =
    $convert.base64Decode(
        'ChtQdXRQYXJ0bmVyRXZlbnRzUmVzdWx0RW50cnkSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcn'
        'JvcmNvZGUSJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdlEhwKB2V2ZW50'
        'aWQY05bdswEgASgJUgdldmVudGlk');

@$core.Deprecated('Use putPermissionRequestDescriptor instead')
const PutPermissionRequest$json = {
  '1': 'PutPermissionRequest',
  '2': [
    {'1': 'action', '3': 175614240, '4': 1, '5': 9, '10': 'action'},
    {
      '1': 'condition',
      '3': 212017247,
      '4': 1,
      '5': 11,
      '6': '.events.Condition',
      '10': 'condition'
    },
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'principal', '3': 361640138, '4': 1, '5': 9, '10': 'principal'},
    {'1': 'statementid', '3': 169602348, '4': 1, '5': 9, '10': 'statementid'},
  ],
};

/// Descriptor for `PutPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putPermissionRequestDescriptor = $convert.base64Decode(
    'ChRQdXRQZXJtaXNzaW9uUmVxdWVzdBIZCgZhY3Rpb24YoNLeUyABKAlSBmFjdGlvbhIyCgljb2'
    '5kaXRpb24Y38CMZSABKAsyES5ldmVudHMuQ29uZGl0aW9uUgljb25kaXRpb24SJgoMZXZlbnRi'
    'dXNuYW1lGL2Z69UBIAEoCVIMZXZlbnRidXNuYW1lEhoKBnBvbGljeRig7/DgASABKAlSBnBvbG'
    'ljeRIgCglwcmluY2lwYWwYyuG4rAEgASgJUglwcmluY2lwYWwSIwoLc3RhdGVtZW50aWQYrNrv'
    'UCABKAlSC3N0YXRlbWVudGlk');

@$core.Deprecated('Use putRuleRequestDescriptor instead')
const PutRuleRequest$json = {
  '1': 'PutRuleRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.RuleState',
      '10': 'state'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.events.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `PutRuleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRuleRequestDescriptor = $convert.base64Decode(
    'Cg5QdXRSdWxlUmVxdWVzdBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SJg'
    'oMZXZlbnRidXNuYW1lGL2Z69UBIAEoCVIMZXZlbnRidXNuYW1lEiUKDGV2ZW50cGF0dGVybhiA'
    '96pvIAEoCVIMZXZlbnRwYXR0ZXJuEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSHAoHcm9sZWFybh'
    'iB+OeZASABKAlSB3JvbGVhcm4SMgoSc2NoZWR1bGVleHByZXNzaW9uGP+R29QBIAEoCVISc2No'
    'ZWR1bGVleHByZXNzaW9uEisKBXN0YXRlGJfJsu8BIAEoDjIRLmV2ZW50cy5SdWxlU3RhdGVSBX'
    'N0YXRlEiMKBHRhZ3MYwcH2tQEgAygLMgsuZXZlbnRzLlRhZ1IEdGFncw==');

@$core.Deprecated('Use putRuleResponseDescriptor instead')
const PutRuleResponse$json = {
  '1': 'PutRuleResponse',
  '2': [
    {'1': 'rulearn', '3': 217507655, '4': 1, '5': 9, '10': 'rulearn'},
  ],
};

/// Descriptor for `PutRuleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRuleResponseDescriptor = $convert.base64Decode(
    'Cg9QdXRSdWxlUmVzcG9uc2USGwoHcnVsZWFybhjHzttnIAEoCVIHcnVsZWFybg==');

@$core.Deprecated('Use putTargetsRequestDescriptor instead')
const PutTargetsRequest$json = {
  '1': 'PutTargetsRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'rule', '3': 475696372, '4': 1, '5': 9, '10': 'rule'},
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.events.Target',
      '10': 'targets'
    },
  ],
};

/// Descriptor for `PutTargetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putTargetsRequestDescriptor = $convert.base64Decode(
    'ChFQdXRUYXJnZXRzUmVxdWVzdBImCgxldmVudGJ1c25hbWUYvZnr1QEgASgJUgxldmVudGJ1c2'
    '5hbWUSFgoEcnVsZRj0meriASABKAlSBHJ1bGUSKwoHdGFyZ2V0cxiCm4J9IAMoCzIOLmV2ZW50'
    'cy5UYXJnZXRSB3RhcmdldHM=');

@$core.Deprecated('Use putTargetsResponseDescriptor instead')
const PutTargetsResponse$json = {
  '1': 'PutTargetsResponse',
  '2': [
    {
      '1': 'failedentries',
      '3': 86416685,
      '4': 3,
      '5': 11,
      '6': '.events.PutTargetsResultEntry',
      '10': 'failedentries'
    },
    {
      '1': 'failedentrycount',
      '3': 458519576,
      '4': 1,
      '5': 5,
      '10': 'failedentrycount'
    },
  ],
};

/// Descriptor for `PutTargetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putTargetsResponseDescriptor = $convert.base64Decode(
    'ChJQdXRUYXJnZXRzUmVzcG9uc2USRgoNZmFpbGVkZW50cmllcxitupopIAMoCzIdLmV2ZW50cy'
    '5QdXRUYXJnZXRzUmVzdWx0RW50cnlSDWZhaWxlZGVudHJpZXMSLgoQZmFpbGVkZW50cnljb3Vu'
    'dBiY6NHaASABKAVSEGZhaWxlZGVudHJ5Y291bnQ=');

@$core.Deprecated('Use putTargetsResultEntryDescriptor instead')
const PutTargetsResultEntry$json = {
  '1': 'PutTargetsResultEntry',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'targetid', '3': 46993334, '4': 1, '5': 9, '10': 'targetid'},
  ],
};

/// Descriptor for `PutTargetsResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putTargetsResultEntryDescriptor = $convert.base64Decode(
    'ChVQdXRUYXJnZXRzUmVzdWx0RW50cnkSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcnJvcmNvZG'
    'USJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdlEh0KCHRhcmdldGlkGLaf'
    'tBYgASgJUgh0YXJnZXRpZA==');

@$core.Deprecated('Use redshiftDataParametersDescriptor instead')
const RedshiftDataParameters$json = {
  '1': 'RedshiftDataParameters',
  '2': [
    {'1': 'database', '3': 278147289, '4': 1, '5': 9, '10': 'database'},
    {'1': 'dbuser', '3': 392658689, '4': 1, '5': 9, '10': 'dbuser'},
    {
      '1': 'secretmanagerarn',
      '3': 530394088,
      '4': 1,
      '5': 9,
      '10': 'secretmanagerarn'
    },
    {'1': 'sql', '3': 13608736, '4': 1, '5': 9, '10': 'sql'},
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
    {'1': 'withevent', '3': 5518724, '4': 1, '5': 8, '10': 'withevent'},
  ],
};

/// Descriptor for `RedshiftDataParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redshiftDataParametersDescriptor = $convert.base64Decode(
    'ChZSZWRzaGlmdERhdGFQYXJhbWV0ZXJzEh4KCGRhdGFiYXNlGNnh0IQBIAEoCVIIZGF0YWJhc2'
    'USGgoGZGJ1c2VyGIH+nbsBIAEoCVIGZGJ1c2VyEi4KEHNlY3JldG1hbmFnZXJhcm4Y6Nf0/AEg'
    'ASgJUhBzZWNyZXRtYW5hZ2VyYXJuEhMKA3NxbBigzr4GIAEoCVIDc3FsEicKDXN0YXRlbWVudG'
    '5hbWUY9t3+CiABKAlSDXN0YXRlbWVudG5hbWUSHwoJd2l0aGV2ZW50GITr0AIgASgIUgl3aXRo'
    'ZXZlbnQ=');

@$core.Deprecated('Use removePermissionRequestDescriptor instead')
const RemovePermissionRequest$json = {
  '1': 'RemovePermissionRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {
      '1': 'removeallpermissions',
      '3': 8309441,
      '4': 1,
      '5': 8,
      '10': 'removeallpermissions'
    },
    {'1': 'statementid', '3': 169602348, '4': 1, '5': 9, '10': 'statementid'},
  ],
};

/// Descriptor for `RemovePermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removePermissionRequestDescriptor = $convert.base64Decode(
    'ChdSZW1vdmVQZXJtaXNzaW9uUmVxdWVzdBImCgxldmVudGJ1c25hbWUYvZnr1QEgASgJUgxldm'
    'VudGJ1c25hbWUSNQoUcmVtb3ZlYWxscGVybWlzc2lvbnMYwZX7AyABKAhSFHJlbW92ZWFsbHBl'
    'cm1pc3Npb25zEiMKC3N0YXRlbWVudGlkGKza71AgASgJUgtzdGF0ZW1lbnRpZA==');

@$core.Deprecated('Use removeTargetsRequestDescriptor instead')
const RemoveTargetsRequest$json = {
  '1': 'RemoveTargetsRequest',
  '2': [
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'force', '3': 526711333, '4': 1, '5': 8, '10': 'force'},
    {'1': 'ids', '3': 56356874, '4': 3, '5': 9, '10': 'ids'},
    {'1': 'rule', '3': 475696372, '4': 1, '5': 9, '10': 'rule'},
  ],
};

/// Descriptor for `RemoveTargetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTargetsRequestDescriptor = $convert.base64Decode(
    'ChRSZW1vdmVUYXJnZXRzUmVxdWVzdBImCgxldmVudGJ1c25hbWUYvZnr1QEgASgJUgxldmVudG'
    'J1c25hbWUSGAoFZm9yY2UYpfST+wEgASgIUgVmb3JjZRITCgNpZHMYiuDvGiADKAlSA2lkcxIW'
    'CgRydWxlGPSZ6uIBIAEoCVIEcnVsZQ==');

@$core.Deprecated('Use removeTargetsResponseDescriptor instead')
const RemoveTargetsResponse$json = {
  '1': 'RemoveTargetsResponse',
  '2': [
    {
      '1': 'failedentries',
      '3': 86416685,
      '4': 3,
      '5': 11,
      '6': '.events.RemoveTargetsResultEntry',
      '10': 'failedentries'
    },
    {
      '1': 'failedentrycount',
      '3': 458519576,
      '4': 1,
      '5': 5,
      '10': 'failedentrycount'
    },
  ],
};

/// Descriptor for `RemoveTargetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTargetsResponseDescriptor = $convert.base64Decode(
    'ChVSZW1vdmVUYXJnZXRzUmVzcG9uc2USSQoNZmFpbGVkZW50cmllcxitupopIAMoCzIgLmV2ZW'
    '50cy5SZW1vdmVUYXJnZXRzUmVzdWx0RW50cnlSDWZhaWxlZGVudHJpZXMSLgoQZmFpbGVkZW50'
    'cnljb3VudBiY6NHaASABKAVSEGZhaWxlZGVudHJ5Y291bnQ=');

@$core.Deprecated('Use removeTargetsResultEntryDescriptor instead')
const RemoveTargetsResultEntry$json = {
  '1': 'RemoveTargetsResultEntry',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'targetid', '3': 46993334, '4': 1, '5': 9, '10': 'targetid'},
  ],
};

/// Descriptor for `RemoveTargetsResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTargetsResultEntryDescriptor = $convert.base64Decode(
    'ChhSZW1vdmVUYXJnZXRzUmVzdWx0RW50cnkSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcnJvcm'
    'NvZGUSJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdlEh0KCHRhcmdldGlk'
    'GLaftBYgASgJUgh0YXJnZXRpZA==');

@$core.Deprecated('Use replayDescriptor instead')
const Replay$json = {
  '1': 'Replay',
  '2': [
    {'1': 'eventendtime', '3': 30229582, '4': 1, '5': 9, '10': 'eventendtime'},
    {
      '1': 'eventlastreplayedtime',
      '3': 7759073,
      '4': 1,
      '5': 9,
      '10': 'eventlastreplayedtime'
    },
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'eventstarttime',
      '3': 103680757,
      '4': 1,
      '5': 9,
      '10': 'eventstarttime'
    },
    {
      '1': 'replayendtime',
      '3': 304261199,
      '4': 1,
      '5': 9,
      '10': 'replayendtime'
    },
    {'1': 'replayname', '3': 442173850, '4': 1, '5': 9, '10': 'replayname'},
    {
      '1': 'replaystarttime',
      '3': 503580492,
      '4': 1,
      '5': 9,
      '10': 'replaystarttime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ReplayState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `Replay`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replayDescriptor = $convert.base64Decode(
    'CgZSZXBsYXkSJQoMZXZlbnRlbmR0aW1lGM6ItQ4gASgJUgxldmVudGVuZHRpbWUSNwoVZXZlbn'
    'RsYXN0cmVwbGF5ZWR0aW1lGOHJ2QMgASgJUhVldmVudGxhc3RyZXBsYXllZHRpbWUSKgoOZXZl'
    'bnRzb3VyY2Vhcm4YxsqKkgEgASgJUg5ldmVudHNvdXJjZWFybhIpCg5ldmVudHN0YXJ0dGltZR'
    'j1lbgxIAEoCVIOZXZlbnRzdGFydHRpbWUSKAoNcmVwbGF5ZW5kdGltZRjP0IqRASABKAlSDXJl'
    'cGxheWVuZHRpbWUSIgoKcmVwbGF5bmFtZRiak+zSASABKAlSCnJlcGxheW5hbWUSLAoPcmVwbG'
    'F5c3RhcnR0aW1lGMyOkPABIAEoCVIPcmVwbGF5c3RhcnR0aW1lEi0KBXN0YXRlGJfJsu8BIAEo'
    'DjITLmV2ZW50cy5SZXBsYXlTdGF0ZVIFc3RhdGUSJAoLc3RhdGVyZWFzb24Y89WtswEgASgJUg'
    'tzdGF0ZXJlYXNvbg==');

@$core.Deprecated('Use replayDestinationDescriptor instead')
const ReplayDestination$json = {
  '1': 'ReplayDestination',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'filterarns', '3': 135157548, '4': 3, '5': 9, '10': 'filterarns'},
  ],
};

/// Descriptor for `ReplayDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replayDestinationDescriptor = $convert.base64Decode(
    'ChFSZXBsYXlEZXN0aW5hdGlvbhIUCgNhcm4YnZvtvwEgASgJUgNhcm4SIQoKZmlsdGVyYXJucx'
    'isrrlAIAMoCVIKZmlsdGVyYXJucw==');

@$core.Deprecated('Use resourceAlreadyExistsExceptionDescriptor instead')
const ResourceAlreadyExistsException$json = {
  '1': 'ResourceAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'Ch5SZXNvdXJjZUFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use retryPolicyDescriptor instead')
const RetryPolicy$json = {
  '1': 'RetryPolicy',
  '2': [
    {
      '1': 'maximumeventageinseconds',
      '3': 393041563,
      '4': 1,
      '5': 5,
      '10': 'maximumeventageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
  ],
};

/// Descriptor for `RetryPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List retryPolicyDescriptor = $convert.base64Decode(
    'CgtSZXRyeVBvbGljeRI+ChhtYXhpbXVtZXZlbnRhZ2VpbnNlY29uZHMYm621uwEgASgFUhhtYX'
    'hpbXVtZXZlbnRhZ2VpbnNlY29uZHMSNQoUbWF4aW11bXJldHJ5YXR0ZW1wdHMYwKi5NSABKAVS'
    'FG1heGltdW1yZXRyeWF0dGVtcHRz');

@$core.Deprecated('Use ruleDescriptor instead')
const Rule$json = {
  '1': 'Rule',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'eventbusname', '3': 448449725, '4': 1, '5': 9, '10': 'eventbusname'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
    {'1': 'managedby', '3': 455511232, '4': 1, '5': 9, '10': 'managedby'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.RuleState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `Rule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ruleDescriptor = $convert.base64Decode(
    'CgRSdWxlEhQKA2Fybhidm+2/ASABKAlSA2FybhIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZG'
    'VzY3JpcHRpb24SJgoMZXZlbnRidXNuYW1lGL2Z69UBIAEoCVIMZXZlbnRidXNuYW1lEiUKDGV2'
    'ZW50cGF0dGVybhiA96pvIAEoCVIMZXZlbnRwYXR0ZXJuEiAKCW1hbmFnZWRieRjAmZrZASABKA'
    'lSCW1hbmFnZWRieRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhwKB3JvbGVhcm4YgfjnmQEgASgJ'
    'Ugdyb2xlYXJuEjIKEnNjaGVkdWxlZXhwcmVzc2lvbhj/kdvUASABKAlSEnNjaGVkdWxlZXhwcm'
    'Vzc2lvbhIrCgVzdGF0ZRiXybLvASABKA4yES5ldmVudHMuUnVsZVN0YXRlUgVzdGF0ZQ==');

@$core.Deprecated('Use runCommandParametersDescriptor instead')
const RunCommandParameters$json = {
  '1': 'RunCommandParameters',
  '2': [
    {
      '1': 'runcommandtargets',
      '3': 422375352,
      '4': 3,
      '5': 11,
      '6': '.events.RunCommandTarget',
      '10': 'runcommandtargets'
    },
  ],
};

/// Descriptor for `RunCommandParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runCommandParametersDescriptor = $convert.base64Decode(
    'ChRSdW5Db21tYW5kUGFyYW1ldGVycxJKChFydW5jb21tYW5kdGFyZ2V0cxi437PJASADKAsyGC'
    '5ldmVudHMuUnVuQ29tbWFuZFRhcmdldFIRcnVuY29tbWFuZHRhcmdldHM=');

@$core.Deprecated('Use runCommandTargetDescriptor instead')
const RunCommandTarget$json = {
  '1': 'RunCommandTarget',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `RunCommandTarget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runCommandTargetDescriptor = $convert.base64Decode(
    'ChBSdW5Db21tYW5kVGFyZ2V0EhMKA2tleRiNkutoIAEoCVIDa2V5EhkKBnZhbHVlcxjcxLRqIA'
    'MoCVIGdmFsdWVz');

@$core.Deprecated('Use sageMakerPipelineParameterDescriptor instead')
const SageMakerPipelineParameter$json = {
  '1': 'SageMakerPipelineParameter',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `SageMakerPipelineParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sageMakerPipelineParameterDescriptor =
    $convert.base64Decode(
        'ChpTYWdlTWFrZXJQaXBlbGluZVBhcmFtZXRlchIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhgKBX'
        'ZhbHVlGOvyn4oBIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use sageMakerPipelineParametersDescriptor instead')
const SageMakerPipelineParameters$json = {
  '1': 'SageMakerPipelineParameters',
  '2': [
    {
      '1': 'pipelineparameterlist',
      '3': 198270119,
      '4': 3,
      '5': 11,
      '6': '.events.SageMakerPipelineParameter',
      '10': 'pipelineparameterlist'
    },
  ],
};

/// Descriptor for `SageMakerPipelineParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sageMakerPipelineParametersDescriptor =
    $convert.base64Decode(
        'ChtTYWdlTWFrZXJQaXBlbGluZVBhcmFtZXRlcnMSWwoVcGlwZWxpbmVwYXJhbWV0ZXJsaXN0GK'
        'e5xV4gAygLMiIuZXZlbnRzLlNhZ2VNYWtlclBpcGVsaW5lUGFyYW1ldGVyUhVwaXBlbGluZXBh'
        'cmFtZXRlcmxpc3Q=');

@$core.Deprecated('Use sqsParametersDescriptor instead')
const SqsParameters$json = {
  '1': 'SqsParameters',
  '2': [
    {
      '1': 'messagegroupid',
      '3': 419537435,
      '4': 1,
      '5': 9,
      '10': 'messagegroupid'
    },
  ],
};

/// Descriptor for `SqsParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sqsParametersDescriptor = $convert.base64Decode(
    'Cg1TcXNQYXJhbWV0ZXJzEioKDm1lc3NhZ2Vncm91cGlkGJvEhsgBIAEoCVIObWVzc2FnZWdyb3'
    'VwaWQ=');

@$core.Deprecated('Use startReplayRequestDescriptor instead')
const StartReplayRequest$json = {
  '1': 'StartReplayRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.events.ReplayDestination',
      '10': 'destination'
    },
    {'1': 'eventendtime', '3': 30229582, '4': 1, '5': 9, '10': 'eventendtime'},
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'eventstarttime',
      '3': 103680757,
      '4': 1,
      '5': 9,
      '10': 'eventstarttime'
    },
    {'1': 'replayname', '3': 442173850, '4': 1, '5': 9, '10': 'replayname'},
  ],
};

/// Descriptor for `StartReplayRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startReplayRequestDescriptor = $convert.base64Decode(
    'ChJTdGFydFJlcGxheVJlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uEj8KC2Rlc3RpbmF0aW9uGOCSkNoBIAEoCzIZLmV2ZW50cy5SZXBsYXlEZXN0aW5hdGlvblIL'
    'ZGVzdGluYXRpb24SJQoMZXZlbnRlbmR0aW1lGM6ItQ4gASgJUgxldmVudGVuZHRpbWUSKgoOZX'
    'ZlbnRzb3VyY2Vhcm4YxsqKkgEgASgJUg5ldmVudHNvdXJjZWFybhIpCg5ldmVudHN0YXJ0dGlt'
    'ZRj1lbgxIAEoCVIOZXZlbnRzdGFydHRpbWUSIgoKcmVwbGF5bmFtZRiak+zSASABKAlSCnJlcG'
    'xheW5hbWU=');

@$core.Deprecated('Use startReplayResponseDescriptor instead')
const StartReplayResponse$json = {
  '1': 'StartReplayResponse',
  '2': [
    {'1': 'replayarn', '3': 361946526, '4': 1, '5': 9, '10': 'replayarn'},
    {
      '1': 'replaystarttime',
      '3': 503580492,
      '4': 1,
      '5': 9,
      '10': 'replaystarttime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ReplayState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `StartReplayResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startReplayResponseDescriptor = $convert.base64Decode(
    'ChNTdGFydFJlcGxheVJlc3BvbnNlEiAKCXJlcGxheWFybhieu8usASABKAlSCXJlcGxheWFybh'
    'IsCg9yZXBsYXlzdGFydHRpbWUYzI6Q8AEgASgJUg9yZXBsYXlzdGFydHRpbWUSLQoFc3RhdGUY'
    'l8my7wEgASgOMhMuZXZlbnRzLlJlcGxheVN0YXRlUgVzdGF0ZRIkCgtzdGF0ZXJlYXNvbhjz1a'
    '2zASABKAlSC3N0YXRlcmVhc29u');

@$core.Deprecated('Use tagDescriptor instead')
const Tag$json = {
  '1': 'Tag',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Tag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagDescriptor = $convert.base64Decode(
    'CgNUYWcSEwoDa2V5GI2S62ggASgJUgNrZXkSGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.events.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZXNvdXJjZW'
    'FybhIjCgR0YWdzGMHB9rUBIAMoCzILLmV2ZW50cy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

@$core.Deprecated('Use targetDescriptor instead')
const Target$json = {
  '1': 'Target',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'batchparameters',
      '3': 242480292,
      '4': 1,
      '5': 11,
      '6': '.events.BatchParameters',
      '10': 'batchparameters'
    },
    {
      '1': 'deadletterconfig',
      '3': 79786642,
      '4': 1,
      '5': 11,
      '6': '.events.DeadLetterConfig',
      '10': 'deadletterconfig'
    },
    {
      '1': 'ecsparameters',
      '3': 501521183,
      '4': 1,
      '5': 11,
      '6': '.events.EcsParameters',
      '10': 'ecsparameters'
    },
    {
      '1': 'httpparameters',
      '3': 367051514,
      '4': 1,
      '5': 11,
      '6': '.events.HttpParameters',
      '10': 'httpparameters'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'input', '3': 529785116, '4': 1, '5': 9, '10': 'input'},
    {'1': 'inputpath', '3': 223809421, '4': 1, '5': 9, '10': 'inputpath'},
    {
      '1': 'inputtransformer',
      '3': 31280095,
      '4': 1,
      '5': 11,
      '6': '.events.InputTransformer',
      '10': 'inputtransformer'
    },
    {
      '1': 'kinesisparameters',
      '3': 70111902,
      '4': 1,
      '5': 11,
      '6': '.events.KinesisParameters',
      '10': 'kinesisparameters'
    },
    {
      '1': 'redshiftdataparameters',
      '3': 254855317,
      '4': 1,
      '5': 11,
      '6': '.events.RedshiftDataParameters',
      '10': 'redshiftdataparameters'
    },
    {
      '1': 'retrypolicy',
      '3': 266827188,
      '4': 1,
      '5': 11,
      '6': '.events.RetryPolicy',
      '10': 'retrypolicy'
    },
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'runcommandparameters',
      '3': 351941188,
      '4': 1,
      '5': 11,
      '6': '.events.RunCommandParameters',
      '10': 'runcommandparameters'
    },
    {
      '1': 'sagemakerpipelineparameters',
      '3': 513800606,
      '4': 1,
      '5': 11,
      '6': '.events.SageMakerPipelineParameters',
      '10': 'sagemakerpipelineparameters'
    },
    {
      '1': 'sqsparameters',
      '3': 91345999,
      '4': 1,
      '5': 11,
      '6': '.events.SqsParameters',
      '10': 'sqsparameters'
    },
  ],
};

/// Descriptor for `Target`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetDescriptor = $convert.base64Decode(
    'CgZUYXJnZXQSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEkQKD2JhdGNocGFyYW1ldGVycxik6c9zIA'
    'EoCzIXLmV2ZW50cy5CYXRjaFBhcmFtZXRlcnNSD2JhdGNocGFyYW1ldGVycxJHChBkZWFkbGV0'
    'dGVyY29uZmlnGJLlhSYgASgLMhguZXZlbnRzLkRlYWRMZXR0ZXJDb25maWdSEGRlYWRsZXR0ZX'
    'Jjb25maWcSPwoNZWNzcGFyYW1ldGVycxiftpLvASABKAsyFS5ldmVudHMuRWNzUGFyYW1ldGVy'
    'c1INZWNzcGFyYW1ldGVycxJCCg5odHRwcGFyYW1ldGVycxj6hYOvASABKAsyFi5ldmVudHMuSH'
    'R0cFBhcmFtZXRlcnNSDmh0dHBwYXJhbWV0ZXJzEhIKAmlkGIHyorcBIAEoCVICaWQSGAoFaW5w'
    'dXQYnMLP/AEgASgJUgVpbnB1dBIfCglpbnB1dHBhdGgYjZ/caiABKAlSCWlucHV0cGF0aBJHCh'
    'BpbnB1dHRyYW5zZm9ybWVyGN+X9Q4gASgLMhguZXZlbnRzLklucHV0VHJhbnNmb3JtZXJSEGlu'
    'cHV0dHJhbnNmb3JtZXISSgoRa2luZXNpc3BhcmFtZXRlcnMYnqW3ISABKAsyGS5ldmVudHMuS2'
    'luZXNpc1BhcmFtZXRlcnNSEWtpbmVzaXNwYXJhbWV0ZXJzElkKFnJlZHNoaWZ0ZGF0YXBhcmFt'
    'ZXRlcnMYlZHDeSABKAsyHi5ldmVudHMuUmVkc2hpZnREYXRhUGFyYW1ldGVyc1IWcmVkc2hpZn'
    'RkYXRhcGFyYW1ldGVycxI4CgtyZXRyeXBvbGljeRi0651/IAEoCzITLmV2ZW50cy5SZXRyeVBv'
    'bGljeVILcmV0cnlwb2xpY3kSHAoHcm9sZWFybhiB+OeZASABKAlSB3JvbGVhcm4SVAoUcnVuY2'
    '9tbWFuZHBhcmFtZXRlcnMYxOTopwEgASgLMhwuZXZlbnRzLlJ1bkNvbW1hbmRQYXJhbWV0ZXJz'
    'UhRydW5jb21tYW5kcGFyYW1ldGVycxJpChtzYWdlbWFrZXJwaXBlbGluZXBhcmFtZXRlcnMYnv'
    'P/9AEgASgLMiMuZXZlbnRzLlNhZ2VNYWtlclBpcGVsaW5lUGFyYW1ldGVyc1Ibc2FnZW1ha2Vy'
    'cGlwZWxpbmVwYXJhbWV0ZXJzEj4KDXNxc3BhcmFtZXRlcnMYz6jHKyABKAsyFS5ldmVudHMuU3'
    'FzUGFyYW1ldGVyc1INc3FzcGFyYW1ldGVycw==');

@$core.Deprecated('Use testEventPatternRequestDescriptor instead')
const TestEventPatternRequest$json = {
  '1': 'TestEventPatternRequest',
  '2': [
    {'1': 'event', '3': 271274816, '4': 1, '5': 9, '10': 'event'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
  ],
};

/// Descriptor for `TestEventPatternRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testEventPatternRequestDescriptor =
    $convert.base64Decode(
        'ChdUZXN0RXZlbnRQYXR0ZXJuUmVxdWVzdBIYCgVldmVudBjApq2BASABKAlSBWV2ZW50EiUKDG'
        'V2ZW50cGF0dGVybhiA96pvIAEoCVIMZXZlbnRwYXR0ZXJu');

@$core.Deprecated('Use testEventPatternResponseDescriptor instead')
const TestEventPatternResponse$json = {
  '1': 'TestEventPatternResponse',
  '2': [
    {'1': 'result', '3': 273346629, '4': 1, '5': 8, '10': 'result'},
  ],
};

/// Descriptor for `TestEventPatternResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testEventPatternResponseDescriptor =
    $convert.base64Decode(
        'ChhUZXN0RXZlbnRQYXR0ZXJuUmVzcG9uc2USGgoGcmVzdWx0GMXgq4IBIAEoCFIGcmVzdWx0');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhjtwJmwASABKAlSC3Jlc291cm'
    'NlYXJuEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagResourceResponseDescriptor instead')
const UntagResourceResponse$json = {
  '1': 'UntagResourceResponse',
};

/// Descriptor for `UntagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceResponseDescriptor =
    $convert.base64Decode('ChVVbnRhZ1Jlc291cmNlUmVzcG9uc2U=');

@$core.Deprecated('Use updateApiDestinationRequestDescriptor instead')
const UpdateApiDestinationRequest$json = {
  '1': 'UpdateApiDestinationRequest',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'invocationendpoint',
      '3': 411800759,
      '4': 1,
      '5': 9,
      '10': 'invocationendpoint'
    },
    {
      '1': 'invocationratelimitpersecond',
      '3': 295327816,
      '4': 1,
      '5': 5,
      '10': 'invocationratelimitpersecond'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateApiDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateApiDestinationRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVBcGlEZXN0aW5hdGlvblJlcXVlc3QSJwoNY29ubmVjdGlvbmFybhjBj7xZIAEoCV'
    'INY29ubmVjdGlvbmFybhIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SRAoK'
    'aHR0cG1ldGhvZBjRjPy9ASABKA4yIC5ldmVudHMuQXBpRGVzdGluYXRpb25IdHRwTWV0aG9kUg'
    'podHRwbWV0aG9kEjIKEmludm9jYXRpb25lbmRwb2ludBi3qa7EASABKAlSEmludm9jYXRpb25l'
    'bmRwb2ludBJGChxpbnZvY2F0aW9ucmF0ZWxpbWl0cGVyc2Vjb25kGMiw6YwBIAEoBVIcaW52b2'
    'NhdGlvbnJhdGVsaW1pdHBlcnNlY29uZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use updateApiDestinationResponseDescriptor instead')
const UpdateApiDestinationResponse$json = {
  '1': 'UpdateApiDestinationResponse',
  '2': [
    {
      '1': 'apidestinationarn',
      '3': 90996885,
      '4': 1,
      '5': 9,
      '10': 'apidestinationarn'
    },
    {
      '1': 'apidestinationstate',
      '3': 13153343,
      '4': 1,
      '5': 14,
      '6': '.events.ApiDestinationState',
      '10': 'apidestinationstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `UpdateApiDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateApiDestinationResponseDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVBcGlEZXN0aW5hdGlvblJlc3BvbnNlEi8KEWFwaWRlc3RpbmF0aW9uYXJuGJWBsi'
    'sgASgJUhFhcGlkZXN0aW5hdGlvbmFybhJQChNhcGlkZXN0aW5hdGlvbnN0YXRlGL/oogYgASgO'
    'MhsuZXZlbnRzLkFwaURlc3RpbmF0aW9uU3RhdGVSE2FwaWRlc3RpbmF0aW9uc3RhdGUSJQoMY3'
    'JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbnRpbWUSLQoQbGFzdG1vZGlmaWVkdGltZRjg'
    'gvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZQ==');

@$core.Deprecated('Use updateArchiveRequestDescriptor instead')
const UpdateArchiveRequest$json = {
  '1': 'UpdateArchiveRequest',
  '2': [
    {'1': 'archivename', '3': 88048487, '4': 1, '5': 9, '10': 'archivename'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'eventpattern', '3': 233487232, '4': 1, '5': 9, '10': 'eventpattern'},
    {
      '1': 'retentiondays',
      '3': 267894223,
      '4': 1,
      '5': 5,
      '10': 'retentiondays'
    },
  ],
};

/// Descriptor for `UpdateArchiveRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateArchiveRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVBcmNoaXZlUmVxdWVzdBIjCgthcmNoaXZlbmFtZRjnhv4pIAEoCVILYXJjaGl2ZW'
    '5hbWUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEiUKDGV2ZW50cGF0dGVy'
    'bhiA96pvIAEoCVIMZXZlbnRwYXR0ZXJuEicKDXJldGVudGlvbmRheXMYz/vefyABKAVSDXJldG'
    'VudGlvbmRheXM=');

@$core.Deprecated('Use updateArchiveResponseDescriptor instead')
const UpdateArchiveResponse$json = {
  '1': 'UpdateArchiveResponse',
  '2': [
    {'1': 'archivearn', '3': 56866685, '4': 1, '5': 9, '10': 'archivearn'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.events.ArchiveState',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
  ],
};

/// Descriptor for `UpdateArchiveResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateArchiveResponseDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVBcmNoaXZlUmVzcG9uc2USIQoKYXJjaGl2ZWFybhj97o4bIAEoCVIKYXJjaGl2ZW'
    'FybhIlCgxjcmVhdGlvbnRpbWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRIuCgVzdGF0ZRiXybLv'
    'ASABKA4yFC5ldmVudHMuQXJjaGl2ZVN0YXRlUgVzdGF0ZRIkCgtzdGF0ZXJlYXNvbhjz1a2zAS'
    'ABKAlSC3N0YXRlcmVhc29u');

@$core.Deprecated(
    'Use updateConnectionApiKeyAuthRequestParametersDescriptor instead')
const UpdateConnectionApiKeyAuthRequestParameters$json = {
  '1': 'UpdateConnectionApiKeyAuthRequestParameters',
  '2': [
    {'1': 'apikeyname', '3': 43133502, '4': 1, '5': 9, '10': 'apikeyname'},
    {'1': 'apikeyvalue', '3': 93786144, '4': 1, '5': 9, '10': 'apikeyvalue'},
  ],
};

/// Descriptor for `UpdateConnectionApiKeyAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateConnectionApiKeyAuthRequestParametersDescriptor =
    $convert.base64Decode(
        'CitVcGRhdGVDb25uZWN0aW9uQXBpS2V5QXV0aFJlcXVlc3RQYXJhbWV0ZXJzEiEKCmFwaWtleW'
        '5hbWUYvtTIFCABKAlSCmFwaWtleW5hbWUSIwoLYXBpa2V5dmFsdWUYoKDcLCABKAlSC2FwaWtl'
        'eXZhbHVl');

@$core.Deprecated('Use updateConnectionAuthRequestParametersDescriptor instead')
const UpdateConnectionAuthRequestParameters$json = {
  '1': 'UpdateConnectionAuthRequestParameters',
  '2': [
    {
      '1': 'apikeyauthparameters',
      '3': 110622489,
      '4': 1,
      '5': 11,
      '6': '.events.UpdateConnectionApiKeyAuthRequestParameters',
      '10': 'apikeyauthparameters'
    },
    {
      '1': 'basicauthparameters',
      '3': 63965312,
      '4': 1,
      '5': 11,
      '6': '.events.UpdateConnectionBasicAuthRequestParameters',
      '10': 'basicauthparameters'
    },
    {
      '1': 'invocationhttpparameters',
      '3': 499128572,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionHttpParameters',
      '10': 'invocationhttpparameters'
    },
    {
      '1': 'oauthparameters',
      '3': 206836569,
      '4': 1,
      '5': 11,
      '6': '.events.UpdateConnectionOAuthRequestParameters',
      '10': 'oauthparameters'
    },
  ],
};

/// Descriptor for `UpdateConnectionAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionAuthRequestParametersDescriptor = $convert.base64Decode(
    'CiVVcGRhdGVDb25uZWN0aW9uQXV0aFJlcXVlc3RQYXJhbWV0ZXJzEmoKFGFwaWtleWF1dGhwYX'
    'JhbWV0ZXJzGJnu3zQgASgLMjMuZXZlbnRzLlVwZGF0ZUNvbm5lY3Rpb25BcGlLZXlBdXRoUmVx'
    'dWVzdFBhcmFtZXRlcnNSFGFwaWtleWF1dGhwYXJhbWV0ZXJzEmcKE2Jhc2ljYXV0aHBhcmFtZX'
    'RlcnMYgJHAHiABKAsyMi5ldmVudHMuVXBkYXRlQ29ubmVjdGlvbkJhc2ljQXV0aFJlcXVlc3RQ'
    'YXJhbWV0ZXJzUhNiYXNpY2F1dGhwYXJhbWV0ZXJzEmAKGGludm9jYXRpb25odHRwcGFyYW1ldG'
    'Vycxj8sYDuASABKAsyIC5ldmVudHMuQ29ubmVjdGlvbkh0dHBQYXJhbWV0ZXJzUhhpbnZvY2F0'
    'aW9uaHR0cHBhcmFtZXRlcnMSWwoPb2F1dGhwYXJhbWV0ZXJzGNmm0GIgASgLMi4uZXZlbnRzLl'
    'VwZGF0ZUNvbm5lY3Rpb25PQXV0aFJlcXVlc3RQYXJhbWV0ZXJzUg9vYXV0aHBhcmFtZXRlcnM=');

@$core.Deprecated(
    'Use updateConnectionBasicAuthRequestParametersDescriptor instead')
const UpdateConnectionBasicAuthRequestParameters$json = {
  '1': 'UpdateConnectionBasicAuthRequestParameters',
  '2': [
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateConnectionBasicAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateConnectionBasicAuthRequestParametersDescriptor =
    $convert.base64Decode(
        'CipVcGRhdGVDb25uZWN0aW9uQmFzaWNBdXRoUmVxdWVzdFBhcmFtZXRlcnMSHQoIcGFzc3dvcm'
        'QYuZCMZiABKAlSCHBhc3N3b3JkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated(
    'Use updateConnectionOAuthClientRequestParametersDescriptor instead')
const UpdateConnectionOAuthClientRequestParameters$json = {
  '1': 'UpdateConnectionOAuthClientRequestParameters',
  '2': [
    {'1': 'clientid', '3': 448889284, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
  ],
};

/// Descriptor for `UpdateConnectionOAuthClientRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateConnectionOAuthClientRequestParametersDescriptor =
    $convert.base64Decode(
        'CixVcGRhdGVDb25uZWN0aW9uT0F1dGhDbGllbnRSZXF1ZXN0UGFyYW1ldGVycxIeCghjbGllbn'
        'RpZBjEg4bWASABKAlSCGNsaWVudGlkEiYKDGNsaWVudHNlY3JldBj3teLuASABKAlSDGNsaWVu'
        'dHNlY3JldA==');

@$core
    .Deprecated('Use updateConnectionOAuthRequestParametersDescriptor instead')
const UpdateConnectionOAuthRequestParameters$json = {
  '1': 'UpdateConnectionOAuthRequestParameters',
  '2': [
    {
      '1': 'authorizationendpoint',
      '3': 427938596,
      '4': 1,
      '5': 9,
      '10': 'authorizationendpoint'
    },
    {
      '1': 'clientparameters',
      '3': 246864127,
      '4': 1,
      '5': 11,
      '6': '.events.UpdateConnectionOAuthClientRequestParameters',
      '10': 'clientparameters'
    },
    {
      '1': 'httpmethod',
      '3': 398394961,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionOAuthHttpMethod',
      '10': 'httpmethod'
    },
    {
      '1': 'oauthhttpparameters',
      '3': 10294537,
      '4': 1,
      '5': 11,
      '6': '.events.ConnectionHttpParameters',
      '10': 'oauthhttpparameters'
    },
  ],
};

/// Descriptor for `UpdateConnectionOAuthRequestParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionOAuthRequestParametersDescriptor = $convert.base64Decode(
    'CiZVcGRhdGVDb25uZWN0aW9uT0F1dGhSZXF1ZXN0UGFyYW1ldGVycxI4ChVhdXRob3JpemF0aW'
    '9uZW5kcG9pbnQYpKaHzAEgASgJUhVhdXRob3JpemF0aW9uZW5kcG9pbnQSYwoQY2xpZW50cGFy'
    'YW1ldGVycxj/sdt1IAEoCzI0LmV2ZW50cy5VcGRhdGVDb25uZWN0aW9uT0F1dGhDbGllbnRSZX'
    'F1ZXN0UGFyYW1ldGVyc1IQY2xpZW50cGFyYW1ldGVycxJFCgpodHRwbWV0aG9kGNGM/L0BIAEo'
    'DjIhLmV2ZW50cy5Db25uZWN0aW9uT0F1dGhIdHRwTWV0aG9kUgpodHRwbWV0aG9kElUKE29hdX'
    'RoaHR0cHBhcmFtZXRlcnMYiar0BCABKAsyIC5ldmVudHMuQ29ubmVjdGlvbkh0dHBQYXJhbWV0'
    'ZXJzUhNvYXV0aGh0dHBwYXJhbWV0ZXJz');

@$core.Deprecated('Use updateConnectionRequestDescriptor instead')
const UpdateConnectionRequest$json = {
  '1': 'UpdateConnectionRequest',
  '2': [
    {
      '1': 'authparameters',
      '3': 258276552,
      '4': 1,
      '5': 11,
      '6': '.events.UpdateConnectionAuthRequestParameters',
      '10': 'authparameters'
    },
    {
      '1': 'authorizationtype',
      '3': 481976511,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionAuthorizationType',
      '10': 'authorizationtype'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateConnectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionRequestDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVDb25uZWN0aW9uUmVxdWVzdBJYCg5hdXRocGFyYW1ldGVycxjI+ZN7IAEoCzItLm'
    'V2ZW50cy5VcGRhdGVDb25uZWN0aW9uQXV0aFJlcXVlc3RQYXJhbWV0ZXJzUg5hdXRocGFyYW1l'
    'dGVycxJVChFhdXRob3JpemF0aW9udHlwZRi/wenlASABKA4yIy5ldmVudHMuQ29ubmVjdGlvbk'
    'F1dGhvcml6YXRpb25UeXBlUhFhdXRob3JpemF0aW9udHlwZRIjCgtkZXNjcmlwdGlvbhiK9Pk2'
    'IAEoCVILZGVzY3JpcHRpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use updateConnectionResponseDescriptor instead')
const UpdateConnectionResponse$json = {
  '1': 'UpdateConnectionResponse',
  '2': [
    {
      '1': 'connectionarn',
      '3': 187631553,
      '4': 1,
      '5': 9,
      '10': 'connectionarn'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.events.ConnectionState',
      '10': 'connectionstate'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastauthorizedtime',
      '3': 250638066,
      '4': 1,
      '5': 9,
      '10': 'lastauthorizedtime'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `UpdateConnectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionResponseDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVDb25uZWN0aW9uUmVzcG9uc2USJwoNY29ubmVjdGlvbmFybhjBj7xZIAEoCVINY2'
    '9ubmVjdGlvbmFybhJFCg9jb25uZWN0aW9uc3RhdGUY2/rlwAEgASgOMhcuZXZlbnRzLkNvbm5l'
    'Y3Rpb25TdGF0ZVIPY29ubmVjdGlvbnN0YXRlEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3'
    'JlYXRpb250aW1lEjEKEmxhc3RhdXRob3JpemVkdGltZRjy3cF3IAEoCVISbGFzdGF1dGhvcml6'
    'ZWR0aW1lEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWU=');
