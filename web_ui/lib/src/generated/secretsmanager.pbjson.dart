// This is a generated file - do not edit.
//
// Generated from secretsmanager.proto.

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

@$core.Deprecated('Use filterNameStringTypeDescriptor instead')
const FilterNameStringType$json = {
  '1': 'FilterNameStringType',
  '2': [
    {'1': 'FILTER_NAME_STRING_TYPE_TAG_KEY', '2': 0},
    {'1': 'FILTER_NAME_STRING_TYPE_ALL', '2': 1},
    {'1': 'FILTER_NAME_STRING_TYPE_PRIMARY_REGION', '2': 2},
    {'1': 'FILTER_NAME_STRING_TYPE_OWNING_SERVICE', '2': 3},
    {'1': 'FILTER_NAME_STRING_TYPE_NAME', '2': 4},
    {'1': 'FILTER_NAME_STRING_TYPE_TAG_VALUE', '2': 5},
    {'1': 'FILTER_NAME_STRING_TYPE_DESCRIPTION', '2': 6},
  ],
};

/// Descriptor for `FilterNameStringType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List filterNameStringTypeDescriptor = $convert.base64Decode(
    'ChRGaWx0ZXJOYW1lU3RyaW5nVHlwZRIjCh9GSUxURVJfTkFNRV9TVFJJTkdfVFlQRV9UQUdfS0'
    'VZEAASHwobRklMVEVSX05BTUVfU1RSSU5HX1RZUEVfQUxMEAESKgomRklMVEVSX05BTUVfU1RS'
    'SU5HX1RZUEVfUFJJTUFSWV9SRUdJT04QAhIqCiZGSUxURVJfTkFNRV9TVFJJTkdfVFlQRV9PV0'
    '5JTkdfU0VSVklDRRADEiAKHEZJTFRFUl9OQU1FX1NUUklOR19UWVBFX05BTUUQBBIlCiFGSUxU'
    'RVJfTkFNRV9TVFJJTkdfVFlQRV9UQUdfVkFMVUUQBRInCiNGSUxURVJfTkFNRV9TVFJJTkdfVF'
    'lQRV9ERVNDUklQVElPThAG');

@$core.Deprecated('Use sortByTypeDescriptor instead')
const SortByType$json = {
  '1': 'SortByType',
  '2': [
    {'1': 'SORT_BY_TYPE_LAST_ACCESSED_DATE', '2': 0},
    {'1': 'SORT_BY_TYPE_CREATED_DATE', '2': 1},
    {'1': 'SORT_BY_TYPE_LAST_CHANGED_DATE', '2': 2},
    {'1': 'SORT_BY_TYPE_NAME', '2': 3},
  ],
};

/// Descriptor for `SortByType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sortByTypeDescriptor = $convert.base64Decode(
    'CgpTb3J0QnlUeXBlEiMKH1NPUlRfQllfVFlQRV9MQVNUX0FDQ0VTU0VEX0RBVEUQABIdChlTT1'
    'JUX0JZX1RZUEVfQ1JFQVRFRF9EQVRFEAESIgoeU09SVF9CWV9UWVBFX0xBU1RfQ0hBTkdFRF9E'
    'QVRFEAISFQoRU09SVF9CWV9UWVBFX05BTUUQAw==');

@$core.Deprecated('Use sortOrderTypeDescriptor instead')
const SortOrderType$json = {
  '1': 'SortOrderType',
  '2': [
    {'1': 'SORT_ORDER_TYPE_ASC', '2': 0},
    {'1': 'SORT_ORDER_TYPE_DESC', '2': 1},
  ],
};

/// Descriptor for `SortOrderType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sortOrderTypeDescriptor = $convert.base64Decode(
    'Cg1Tb3J0T3JkZXJUeXBlEhcKE1NPUlRfT1JERVJfVFlQRV9BU0MQABIYChRTT1JUX09SREVSX1'
    'RZUEVfREVTQxAB');

@$core.Deprecated('Use statusTypeDescriptor instead')
const StatusType$json = {
  '1': 'StatusType',
  '2': [
    {'1': 'STATUS_TYPE_INSYNC', '2': 0},
    {'1': 'STATUS_TYPE_FAILED', '2': 1},
    {'1': 'STATUS_TYPE_INPROGRESS', '2': 2},
  ],
};

/// Descriptor for `StatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statusTypeDescriptor = $convert.base64Decode(
    'CgpTdGF0dXNUeXBlEhYKElNUQVRVU19UWVBFX0lOU1lOQxAAEhYKElNUQVRVU19UWVBFX0ZBSU'
    'xFRBABEhoKFlNUQVRVU19UWVBFX0lOUFJPR1JFU1MQAg==');

@$core.Deprecated('Use aPIErrorTypeDescriptor instead')
const APIErrorType$json = {
  '1': 'APIErrorType',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `APIErrorType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aPIErrorTypeDescriptor = $convert.base64Decode(
    'CgxBUElFcnJvclR5cGUSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcnJvcmNvZGUSGwoHbWVzc2'
    'FnZRiFs7twIAEoCVIHbWVzc2FnZRIeCghzZWNyZXRpZBjl1uuiASABKAlSCHNlY3JldGlk');

@$core.Deprecated('Use batchGetSecretValueRequestDescriptor instead')
const BatchGetSecretValueRequest$json = {
  '1': 'BatchGetSecretValueRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.Filter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'secretidlist', '3': 398967021, '4': 3, '5': 9, '10': 'secretidlist'},
  ],
};

/// Descriptor for `BatchGetSecretValueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetSecretValueRequestDescriptor = $convert.base64Decode(
    'ChpCYXRjaEdldFNlY3JldFZhbHVlUmVxdWVzdBIzCgdmaWx0ZXJzGO3N6lkgAygLMhYuc2Vjcm'
    'V0c21hbmFnZXIuRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhy'
    'ZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiYKDHNlY3JldGlkbGlzdB'
    'jtgZ++ASADKAlSDHNlY3JldGlkbGlzdA==');

@$core.Deprecated('Use batchGetSecretValueResponseDescriptor instead')
const BatchGetSecretValueResponse$json = {
  '1': 'BatchGetSecretValueResponse',
  '2': [
    {
      '1': 'errors',
      '3': 166551719,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.APIErrorType',
      '10': 'errors'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'secretvalues',
      '3': 79551512,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.SecretValueEntry',
      '10': 'secretvalues'
    },
  ],
};

/// Descriptor for `BatchGetSecretValueResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetSecretValueResponseDescriptor = $convert.base64Decode(
    'ChtCYXRjaEdldFNlY3JldFZhbHVlUmVzcG9uc2USNwoGZXJyb3JzGKfBtU8gAygLMhwuc2Vjcm'
    'V0c21hbmFnZXIuQVBJRXJyb3JUeXBlUgZlcnJvcnMSHwoJbmV4dHRva2VuGP6EumcgASgJUglu'
    'ZXh0dG9rZW4SRwoMc2VjcmV0dmFsdWVzGJi49yUgAygLMiAuc2VjcmV0c21hbmFnZXIuU2Vjcm'
    'V0VmFsdWVFbnRyeVIMc2VjcmV0dmFsdWVz');

@$core.Deprecated('Use cancelRotateSecretRequestDescriptor instead')
const CancelRotateSecretRequest$json = {
  '1': 'CancelRotateSecretRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `CancelRotateSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelRotateSecretRequestDescriptor =
    $convert.base64Decode(
        'ChlDYW5jZWxSb3RhdGVTZWNyZXRSZXF1ZXN0Eh4KCHNlY3JldGlkGOXW66IBIAEoCVIIc2Vjcm'
        'V0aWQ=');

@$core.Deprecated('Use cancelRotateSecretResponseDescriptor instead')
const CancelRotateSecretResponse$json = {
  '1': 'CancelRotateSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `CancelRotateSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelRotateSecretResponseDescriptor =
    $convert.base64Decode(
        'ChpDYW5jZWxSb3RhdGVTZWNyZXRSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm4SFQoEbm'
        'FtZRiH5oF/IAEoCVIEbmFtZRIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use createSecretRequestDescriptor instead')
const CreateSecretRequest$json = {
  '1': 'CreateSecretRequest',
  '2': [
    {
      '1': 'addreplicaregions',
      '3': 461171870,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ReplicaRegionType',
      '10': 'addreplicaregions'
    },
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'forceoverwritereplicasecret',
      '3': 247407324,
      '4': 1,
      '5': 8,
      '10': 'forceoverwritereplicasecret'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'secretbinary', '3': 94375681, '4': 1, '5': 12, '10': 'secretbinary'},
    {'1': 'secretstring', '3': 190782253, '4': 1, '5': 9, '10': 'secretstring'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.Tag',
      '10': 'tags'
    },
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `CreateSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSecretRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVTZWNyZXRSZXF1ZXN0ElMKEWFkZHJlcGxpY2FyZWdpb25zGJ7Z89sBIAMoCzIhLn'
    'NlY3JldHNtYW5hZ2VyLlJlcGxpY2FSZWdpb25UeXBlUhFhZGRyZXBsaWNhcmVnaW9ucxIyChJj'
    'bGllbnRyZXF1ZXN0dG9rZW4Y8e+i2QEgASgJUhJjbGllbnRyZXF1ZXN0dG9rZW4SIwoLZGVzY3'
    'JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEkMKG2ZvcmNlb3ZlcndyaXRlcmVwbGljYXNl'
    'Y3JldBjcxfx1IAEoCFIbZm9yY2VvdmVyd3JpdGVyZXBsaWNhc2VjcmV0Eh0KCGttc2tleWlkGI'
    '3JlxYgASgJUghrbXNrZXlpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiUKDHNlY3JldGJpbmFy'
    'eRiBnoAtIAEoDFIMc2VjcmV0YmluYXJ5EiUKDHNlY3JldHN0cmluZxittvxaIAEoCVIMc2Vjcm'
    'V0c3RyaW5nEisKBHRhZ3MYwcH2tQEgAygLMhMuc2VjcmV0c21hbmFnZXIuVGFnUgR0YWdzEhYK'
    'BHR5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use createSecretResponseDescriptor instead')
const CreateSecretResponse$json = {
  '1': 'CreateSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ReplicationStatusType',
      '10': 'replicationstatus'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `CreateSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSecretResponseDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVTZWNyZXRSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm4SFQoEbmFtZRiH5o'
    'F/IAEoCVIEbmFtZRJXChFyZXBsaWNhdGlvbnN0YXR1cxiMqqX8ASADKAsyJS5zZWNyZXRzbWFu'
    'YWdlci5SZXBsaWNhdGlvblN0YXR1c1R5cGVSEXJlcGxpY2F0aW9uc3RhdHVzEiAKCXZlcnNpb2'
    '5pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use decryptionFailureDescriptor instead')
const DecryptionFailure$json = {
  '1': 'DecryptionFailure',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DecryptionFailure`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List decryptionFailureDescriptor = $convert.base64Decode(
    'ChFEZWNyeXB0aW9uRmFpbHVyZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use deleteResourcePolicyRequestDescriptor instead')
const DeleteResourcePolicyRequest$json = {
  '1': 'DeleteResourcePolicyRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `DeleteResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVSZXNvdXJjZVBvbGljeVJlcXVlc3QSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZW'
        'NyZXRpZA==');

@$core.Deprecated('Use deleteResourcePolicyResponseDescriptor instead')
const DeleteResourcePolicyResponse$json = {
  '1': 'DeleteResourcePolicyResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyResponseDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVSZXNvdXJjZVBvbGljeVJlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIVCg'
        'RuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use deleteSecretRequestDescriptor instead')
const DeleteSecretRequest$json = {
  '1': 'DeleteSecretRequest',
  '2': [
    {
      '1': 'forcedeletewithoutrecovery',
      '3': 179882315,
      '4': 1,
      '5': 8,
      '10': 'forcedeletewithoutrecovery'
    },
    {
      '1': 'recoverywindowindays',
      '3': 141399243,
      '4': 1,
      '5': 3,
      '10': 'recoverywindowindays'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `DeleteSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSecretRequestDescriptor = $convert.base64Decode(
    'ChNEZWxldGVTZWNyZXRSZXF1ZXN0EkEKGmZvcmNlZGVsZXRld2l0aG91dHJlY292ZXJ5GMuS41'
    'UgASgIUhpmb3JjZWRlbGV0ZXdpdGhvdXRyZWNvdmVyeRI1ChRyZWNvdmVyeXdpbmRvd2luZGF5'
    'cxjLqbZDIAEoA1IUcmVjb3Zlcnl3aW5kb3dpbmRheXMSHgoIc2VjcmV0aWQY5dbrogEgASgJUg'
    'hzZWNyZXRpZA==');

@$core.Deprecated('Use deleteSecretResponseDescriptor instead')
const DeleteSecretResponse$json = {
  '1': 'DeleteSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'deletiondate', '3': 347845564, '4': 1, '5': 9, '10': 'deletiondate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSecretResponseDescriptor = $convert.base64Decode(
    'ChREZWxldGVTZWNyZXRSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm4SJgoMZGVsZXRpb2'
    '5kYXRlGLzn7qUBIAEoCVIMZGVsZXRpb25kYXRlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use describeSecretRequestDescriptor instead')
const DescribeSecretRequest$json = {
  '1': 'DescribeSecretRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `DescribeSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeSecretRequestDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZVNlY3JldFJlcXVlc3QSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZXRpZA'
    '==');

@$core.Deprecated('Use describeSecretResponseDescriptor instead')
const DescribeSecretResponse$json = {
  '1': 'DescribeSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'deleteddate', '3': 516314255, '4': 1, '5': 9, '10': 'deleteddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'externalsecretrotationmetadata',
      '3': 52900542,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ExternalSecretRotationMetadataItem',
      '10': 'externalsecretrotationmetadata'
    },
    {
      '1': 'externalsecretrotationrolearn',
      '3': 470712576,
      '4': 1,
      '5': 9,
      '10': 'externalsecretrotationrolearn'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastaccesseddate',
      '3': 194418963,
      '4': 1,
      '5': 9,
      '10': 'lastaccesseddate'
    },
    {
      '1': 'lastchangeddate',
      '3': 314015460,
      '4': 1,
      '5': 9,
      '10': 'lastchangeddate'
    },
    {
      '1': 'lastrotateddate',
      '3': 501475691,
      '4': 1,
      '5': 9,
      '10': 'lastrotateddate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextrotationdate',
      '3': 192035355,
      '4': 1,
      '5': 9,
      '10': 'nextrotationdate'
    },
    {
      '1': 'owningservice',
      '3': 462487817,
      '4': 1,
      '5': 9,
      '10': 'owningservice'
    },
    {
      '1': 'primaryregion',
      '3': 480901186,
      '4': 1,
      '5': 9,
      '10': 'primaryregion'
    },
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ReplicationStatusType',
      '10': 'replicationstatus'
    },
    {
      '1': 'rotationenabled',
      '3': 209507301,
      '4': 1,
      '5': 8,
      '10': 'rotationenabled'
    },
    {
      '1': 'rotationlambdaarn',
      '3': 335026080,
      '4': 1,
      '5': 9,
      '10': 'rotationlambdaarn'
    },
    {
      '1': 'rotationrules',
      '3': 259458135,
      '4': 1,
      '5': 11,
      '6': '.secretsmanager.RotationRulesType',
      '10': 'rotationrules'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.Tag',
      '10': 'tags'
    },
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {
      '1': 'versionidstostages',
      '3': 90698314,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.DescribeSecretResponse.VersionidstostagesEntry',
      '10': 'versionidstostages'
    },
  ],
  '3': [DescribeSecretResponse_VersionidstostagesEntry$json],
};

@$core.Deprecated('Use describeSecretResponseDescriptor instead')
const DescribeSecretResponse_VersionidstostagesEntry$json = {
  '1': 'VersionidstostagesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DescribeSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeSecretResponseDescriptor = $convert.base64Decode(
    'ChZEZXNjcmliZVNlY3JldFJlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIkCgtjcmVhdG'
    'VkZGF0ZRiwsOfGASABKAlSC2NyZWF0ZWRkYXRlEiQKC2RlbGV0ZWRkYXRlGI+pmfYBIAEoCVIL'
    'ZGVsZXRlZGRhdGUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEn0KHmV4dG'
    'VybmFsc2VjcmV0cm90YXRpb25tZXRhZGF0YRi+5ZwZIAMoCzIyLnNlY3JldHNtYW5hZ2VyLkV4'
    'dGVybmFsU2VjcmV0Um90YXRpb25NZXRhZGF0YUl0ZW1SHmV4dGVybmFsc2VjcmV0cm90YXRpb2'
    '5tZXRhZGF0YRJICh1leHRlcm5hbHNlY3JldHJvdGF0aW9ucm9sZWFybhiAgrrgASABKAlSHWV4'
    'dGVybmFsc2VjcmV0cm90YXRpb25yb2xlYXJuEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNrZX'
    'lpZBItChBsYXN0YWNjZXNzZWRkYXRlGJOy2lwgASgJUhBsYXN0YWNjZXNzZWRkYXRlEiwKD2xh'
    'c3RjaGFuZ2VkZGF0ZRjk/d2VASABKAlSD2xhc3RjaGFuZ2VkZGF0ZRIsCg9sYXN0cm90YXRlZG'
    'RhdGUY69KP7wEgASgJUg9sYXN0cm90YXRlZGRhdGUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIt'
    'ChBuZXh0cm90YXRpb25kYXRlGJv0yFsgASgJUhBuZXh0cm90YXRpb25kYXRlEigKDW93bmluZ3'
    'NlcnZpY2UYiYLE3AEgASgJUg1vd25pbmdzZXJ2aWNlEigKDXByaW1hcnlyZWdpb24YwvCn5QEg'
    'ASgJUg1wcmltYXJ5cmVnaW9uElcKEXJlcGxpY2F0aW9uc3RhdHVzGIyqpfwBIAMoCzIlLnNlY3'
    'JldHNtYW5hZ2VyLlJlcGxpY2F0aW9uU3RhdHVzVHlwZVIRcmVwbGljYXRpb25zdGF0dXMSKwoP'
    'cm90YXRpb25lbmFibGVkGOWn82MgASgIUg9yb3RhdGlvbmVuYWJsZWQSMAoRcm90YXRpb25sYW'
    '1iZGFhcm4YoK/gnwEgASgJUhFyb3RhdGlvbmxhbWJkYWFybhJKCg1yb3RhdGlvbnJ1bGVzGNeI'
    '3HsgASgLMiEuc2VjcmV0c21hbmFnZXIuUm90YXRpb25SdWxlc1R5cGVSDXJvdGF0aW9ucnVsZX'
    'MSKwoEdGFncxjBwfa1ASADKAsyEy5zZWNyZXRzbWFuYWdlci5UYWdSBHRhZ3MSFgoEdHlwZRju'
    'oNeKASABKAlSBHR5cGUScQoSdmVyc2lvbmlkc3Rvc3RhZ2VzGMrknysgAygLMj4uc2VjcmV0c2'
    '1hbmFnZXIuRGVzY3JpYmVTZWNyZXRSZXNwb25zZS5WZXJzaW9uaWRzdG9zdGFnZXNFbnRyeVIS'
    'dmVyc2lvbmlkc3Rvc3RhZ2VzGkUKF1ZlcnNpb25pZHN0b3N0YWdlc0VudHJ5EhAKA2tleRgBIA'
    'EoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use encryptionFailureDescriptor instead')
const EncryptionFailure$json = {
  '1': 'EncryptionFailure',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EncryptionFailure`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionFailureDescriptor = $convert.base64Decode(
    'ChFFbmNyeXB0aW9uRmFpbHVyZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use externalSecretRotationMetadataItemDescriptor instead')
const ExternalSecretRotationMetadataItem$json = {
  '1': 'ExternalSecretRotationMetadataItem',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ExternalSecretRotationMetadataItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List externalSecretRotationMetadataItemDescriptor =
    $convert.base64Decode(
        'CiJFeHRlcm5hbFNlY3JldFJvdGF0aW9uTWV0YWRhdGFJdGVtEhMKA2tleRiNkutoIAEoCVIDa2'
        'V5EhgKBXZhbHVlGOvyn4oBIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use filterDescriptor instead')
const Filter$json = {
  '1': 'Filter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.secretsmanager.FilterNameStringType',
      '10': 'key'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `Filter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterDescriptor = $convert.base64Decode(
    'CgZGaWx0ZXISOQoDa2V5GI2S62ggASgOMiQuc2VjcmV0c21hbmFnZXIuRmlsdGVyTmFtZVN0cm'
    'luZ1R5cGVSA2tleRIZCgZ2YWx1ZXMY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use getRandomPasswordRequestDescriptor instead')
const GetRandomPasswordRequest$json = {
  '1': 'GetRandomPasswordRequest',
  '2': [
    {
      '1': 'excludecharacters',
      '3': 335851390,
      '4': 1,
      '5': 9,
      '10': 'excludecharacters'
    },
    {
      '1': 'excludelowercase',
      '3': 225858843,
      '4': 1,
      '5': 8,
      '10': 'excludelowercase'
    },
    {
      '1': 'excludenumbers',
      '3': 216382246,
      '4': 1,
      '5': 8,
      '10': 'excludenumbers'
    },
    {
      '1': 'excludepunctuation',
      '3': 78530732,
      '4': 1,
      '5': 8,
      '10': 'excludepunctuation'
    },
    {
      '1': 'excludeuppercase',
      '3': 345594144,
      '4': 1,
      '5': 8,
      '10': 'excludeuppercase'
    },
    {'1': 'includespace', '3': 216842628, '4': 1, '5': 8, '10': 'includespace'},
    {
      '1': 'passwordlength',
      '3': 216611365,
      '4': 1,
      '5': 3,
      '10': 'passwordlength'
    },
    {
      '1': 'requireeachincludedtype',
      '3': 176215282,
      '4': 1,
      '5': 8,
      '10': 'requireeachincludedtype'
    },
  ],
};

/// Descriptor for `GetRandomPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRandomPasswordRequestDescriptor = $convert.base64Decode(
    'ChhHZXRSYW5kb21QYXNzd29yZFJlcXVlc3QSMAoRZXhjbHVkZWNoYXJhY3RlcnMY/t6SoAEgAS'
    'gJUhFleGNsdWRlY2hhcmFjdGVycxItChBleGNsdWRlbG93ZXJjYXNlGJuq2WsgASgIUhBleGNs'
    'dWRlbG93ZXJjYXNlEikKDmV4Y2x1ZGVudW1iZXJzGKb2lmcgASgIUg5leGNsdWRlbnVtYmVycx'
    'IxChJleGNsdWRlcHVuY3R1YXRpb24YrJG5JSABKAhSEmV4Y2x1ZGVwdW5jdHVhdGlvbhIuChBl'
    'eGNsdWRldXBwZXJjYXNlGKCy5aQBIAEoCFIQZXhjbHVkZXVwcGVyY2FzZRIlCgxpbmNsdWRlc3'
    'BhY2UYhIOzZyABKAhSDGluY2x1ZGVzcGFjZRIpCg5wYXNzd29yZGxlbmd0aBil9KRnIAEoA1IO'
    'cGFzc3dvcmRsZW5ndGgSOwoXcmVxdWlyZWVhY2hpbmNsdWRlZHR5cGUY8qmDVCABKAhSF3JlcX'
    'VpcmVlYWNoaW5jbHVkZWR0eXBl');

@$core.Deprecated('Use getRandomPasswordResponseDescriptor instead')
const GetRandomPasswordResponse$json = {
  '1': 'GetRandomPasswordResponse',
  '2': [
    {
      '1': 'randompassword',
      '3': 477267856,
      '4': 1,
      '5': 9,
      '10': 'randompassword'
    },
  ],
};

/// Descriptor for `GetRandomPasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRandomPasswordResponseDescriptor =
    $convert.base64Decode(
        'ChlHZXRSYW5kb21QYXNzd29yZFJlc3BvbnNlEioKDnJhbmRvbXBhc3N3b3JkGJCPyuMBIAEoCV'
        'IOcmFuZG9tcGFzc3dvcmQ=');

@$core.Deprecated('Use getResourcePolicyRequestDescriptor instead')
const GetResourcePolicyRequest$json = {
  '1': 'GetResourcePolicyRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `GetResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZX'
        'RpZA==');

@$core.Deprecated('Use getResourcePolicyResponseDescriptor instead')
const GetResourcePolicyResponse$json = {
  '1': 'GetResourcePolicyResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
  ],
};

/// Descriptor for `GetResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyResponseDescriptor = $convert.base64Decode(
    'ChlHZXRSZXNvdXJjZVBvbGljeVJlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIVCgRuYW'
    '1lGIfmgX8gASgJUgRuYW1lEikKDnJlc291cmNlcG9saWN5GLCUwQcgASgJUg5yZXNvdXJjZXBv'
    'bGljeQ==');

@$core.Deprecated('Use getSecretValueRequestDescriptor instead')
const GetSecretValueRequest$json = {
  '1': 'GetSecretValueRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {'1': 'versionstage', '3': 229692340, '4': 1, '5': 9, '10': 'versionstage'},
  ],
};

/// Descriptor for `GetSecretValueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSecretValueRequestDescriptor = $convert.base64Decode(
    'ChVHZXRTZWNyZXRWYWx1ZVJlcXVlc3QSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZXRpZB'
    'IgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQSJQoMdmVyc2lvbnN0YWdlGLSnw20g'
    'ASgJUgx2ZXJzaW9uc3RhZ2U=');

@$core.Deprecated('Use getSecretValueResponseDescriptor instead')
const GetSecretValueResponse$json = {
  '1': 'GetSecretValueResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'secretbinary', '3': 94375681, '4': 1, '5': 12, '10': 'secretbinary'},
    {'1': 'secretstring', '3': 190782253, '4': 1, '5': 9, '10': 'secretstring'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {
      '1': 'versionstages',
      '3': 224220993,
      '4': 3,
      '5': 9,
      '10': 'versionstages'
    },
  ],
};

/// Descriptor for `GetSecretValueResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSecretValueResponseDescriptor = $convert.base64Decode(
    'ChZHZXRTZWNyZXRWYWx1ZVJlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIkCgtjcmVhdG'
    'VkZGF0ZRiwsOfGASABKAlSC2NyZWF0ZWRkYXRlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJQoM'
    'c2VjcmV0YmluYXJ5GIGegC0gASgMUgxzZWNyZXRiaW5hcnkSJQoMc2VjcmV0c3RyaW5nGK22/F'
    'ogASgJUgxzZWNyZXRzdHJpbmcSIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlkEicK'
    'DXZlcnNpb25zdGFnZXMYwa71aiADKAlSDXZlcnNpb25zdGFnZXM=');

@$core.Deprecated('Use internalServiceErrorDescriptor instead')
const InternalServiceError$json = {
  '1': 'InternalServiceError',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServiceError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServiceErrorDescriptor =
    $convert.base64Decode(
        'ChRJbnRlcm5hbFNlcnZpY2VFcnJvchIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidNextTokenExceptionDescriptor instead')
const InvalidNextTokenException$json = {
  '1': 'InvalidNextTokenException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidNextTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNextTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkTmV4dFRva2VuRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidParameterExceptionDescriptor instead')
const InvalidParameterException$json = {
  '1': 'InvalidParameterException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidRequestExceptionDescriptor instead')
const InvalidRequestException$json = {
  '1': 'InvalidRequestException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRequestExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkUmVxdWVzdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use limitExceededExceptionDescriptor instead')
const LimitExceededException$json = {
  '1': 'LimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChZMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use listSecretVersionIdsRequestDescriptor instead')
const ListSecretVersionIdsRequest$json = {
  '1': 'ListSecretVersionIdsRequest',
  '2': [
    {
      '1': 'includedeprecated',
      '3': 299058751,
      '4': 1,
      '5': 8,
      '10': 'includedeprecated'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `ListSecretVersionIdsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSecretVersionIdsRequestDescriptor = $convert.base64Decode(
    'ChtMaXN0U2VjcmV0VmVyc2lvbklkc1JlcXVlc3QSMAoRaW5jbHVkZWRlcHJlY2F0ZWQYv4zNjg'
    'EgASgIUhFpbmNsdWRlZGVwcmVjYXRlZBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVz'
    'dWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIeCghzZWNyZXRpZBjl1uuiAS'
    'ABKAlSCHNlY3JldGlk');

@$core.Deprecated('Use listSecretVersionIdsResponseDescriptor instead')
const ListSecretVersionIdsResponse$json = {
  '1': 'ListSecretVersionIdsResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'versions',
      '3': 252099085,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.SecretVersionsListEntry',
      '10': 'versions'
    },
  ],
};

/// Descriptor for `ListSecretVersionIdsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSecretVersionIdsResponseDescriptor = $convert.base64Decode(
    'ChxMaXN0U2VjcmV0VmVyc2lvbklkc1Jlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIVCg'
    'RuYW1lGIfmgX8gASgJUgRuYW1lEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEkYK'
    'CHZlcnNpb25zGI30mnggAygLMicuc2VjcmV0c21hbmFnZXIuU2VjcmV0VmVyc2lvbnNMaXN0RW'
    '50cnlSCHZlcnNpb25z');

@$core.Deprecated('Use listSecretsRequestDescriptor instead')
const ListSecretsRequest$json = {
  '1': 'ListSecretsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.Filter',
      '10': 'filters'
    },
    {
      '1': 'includeplanneddeletion',
      '3': 64231622,
      '4': 1,
      '5': 8,
      '10': 'includeplanneddeletion'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'sortby',
      '3': 186052369,
      '4': 1,
      '5': 14,
      '6': '.secretsmanager.SortByType',
      '10': 'sortby'
    },
    {
      '1': 'sortorder',
      '3': 274231684,
      '4': 1,
      '5': 14,
      '6': '.secretsmanager.SortOrderType',
      '10': 'sortorder'
    },
  ],
};

/// Descriptor for `ListSecretsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSecretsRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0U2VjcmV0c1JlcXVlc3QSMwoHZmlsdGVycxjtzepZIAMoCzIWLnNlY3JldHNtYW5hZ2'
    'VyLkZpbHRlclIHZmlsdGVycxI5ChZpbmNsdWRlcGxhbm5lZGRlbGV0aW9uGMax0B4gASgIUhZp'
    'bmNsdWRlcGxhbm5lZGRlbGV0aW9uEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bH'
    'RzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEjUKBnNvcnRieRiR3ttYIAEoDjIa'
    'LnNlY3JldHNtYW5hZ2VyLlNvcnRCeVR5cGVSBnNvcnRieRI/Cglzb3J0b3JkZXIYhOPhggEgAS'
    'gOMh0uc2VjcmV0c21hbmFnZXIuU29ydE9yZGVyVHlwZVIJc29ydG9yZGVy');

@$core.Deprecated('Use listSecretsResponseDescriptor instead')
const ListSecretsResponse$json = {
  '1': 'ListSecretsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'secretlist',
      '3': 280320894,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.SecretListEntry',
      '10': 'secretlist'
    },
  ],
};

/// Descriptor for `ListSecretsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSecretsResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0U2VjcmV0c1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEk'
    'MKCnNlY3JldGxpc3QY/rbVhQEgAygLMh8uc2VjcmV0c21hbmFnZXIuU2VjcmV0TGlzdEVudHJ5'
    'UgpzZWNyZXRsaXN0');

@$core.Deprecated('Use malformedPolicyDocumentExceptionDescriptor instead')
const MalformedPolicyDocumentException$json = {
  '1': 'MalformedPolicyDocumentException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MalformedPolicyDocumentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List malformedPolicyDocumentExceptionDescriptor =
    $convert.base64Decode(
        'CiBNYWxmb3JtZWRQb2xpY3lEb2N1bWVudEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use preconditionNotMetExceptionDescriptor instead')
const PreconditionNotMetException$json = {
  '1': 'PreconditionNotMetException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PreconditionNotMetException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preconditionNotMetExceptionDescriptor =
    $convert.base64Decode(
        'ChtQcmVjb25kaXRpb25Ob3RNZXRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use publicPolicyExceptionDescriptor instead')
const PublicPolicyException$json = {
  '1': 'PublicPolicyException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PublicPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicPolicyExceptionDescriptor = $convert.base64Decode(
    'ChVQdWJsaWNQb2xpY3lFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use putResourcePolicyRequestDescriptor instead')
const PutResourcePolicyRequest$json = {
  '1': 'PutResourcePolicyRequest',
  '2': [
    {
      '1': 'blockpublicpolicy',
      '3': 505379326,
      '4': 1,
      '5': 8,
      '10': 'blockpublicpolicy'
    },
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `PutResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyRequestDescriptor = $convert.base64Decode(
    'ChhQdXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSMAoRYmxvY2twdWJsaWNwb2xpY3kY/vP98AEgAS'
    'gIUhFibG9ja3B1YmxpY3BvbGljeRIpCg5yZXNvdXJjZXBvbGljeRiwlMEHIAEoCVIOcmVzb3Vy'
    'Y2Vwb2xpY3kSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZXRpZA==');

@$core.Deprecated('Use putResourcePolicyResponseDescriptor instead')
const PutResourcePolicyResponse$json = {
  '1': 'PutResourcePolicyResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `PutResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyResponseDescriptor =
    $convert.base64Decode(
        'ChlQdXRSZXNvdXJjZVBvbGljeVJlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIVCgRuYW'
        '1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use putSecretValueRequestDescriptor instead')
const PutSecretValueRequest$json = {
  '1': 'PutSecretValueRequest',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {
      '1': 'rotationtoken',
      '3': 292175477,
      '4': 1,
      '5': 9,
      '10': 'rotationtoken'
    },
    {'1': 'secretbinary', '3': 94375681, '4': 1, '5': 12, '10': 'secretbinary'},
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
    {'1': 'secretstring', '3': 190782253, '4': 1, '5': 9, '10': 'secretstring'},
    {
      '1': 'versionstages',
      '3': 224220993,
      '4': 3,
      '5': 9,
      '10': 'versionstages'
    },
  ],
};

/// Descriptor for `PutSecretValueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putSecretValueRequestDescriptor = $convert.base64Decode(
    'ChVQdXRTZWNyZXRWYWx1ZVJlcXVlc3QSMgoSY2xpZW50cmVxdWVzdHRva2VuGPHvotkBIAEoCV'
    'ISY2xpZW50cmVxdWVzdHRva2VuEigKDXJvdGF0aW9udG9rZW4Y9fyoiwEgASgJUg1yb3RhdGlv'
    'bnRva2VuEiUKDHNlY3JldGJpbmFyeRiBnoAtIAEoDFIMc2VjcmV0YmluYXJ5Eh4KCHNlY3JldG'
    'lkGOXW66IBIAEoCVIIc2VjcmV0aWQSJQoMc2VjcmV0c3RyaW5nGK22/FogASgJUgxzZWNyZXRz'
    'dHJpbmcSJwoNdmVyc2lvbnN0YWdlcxjBrvVqIAMoCVINdmVyc2lvbnN0YWdlcw==');

@$core.Deprecated('Use putSecretValueResponseDescriptor instead')
const PutSecretValueResponse$json = {
  '1': 'PutSecretValueResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {
      '1': 'versionstages',
      '3': 224220993,
      '4': 3,
      '5': 9,
      '10': 'versionstages'
    },
  ],
};

/// Descriptor for `PutSecretValueResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putSecretValueResponseDescriptor = $convert.base64Decode(
    'ChZQdXRTZWNyZXRWYWx1ZVJlc3BvbnNlEhQKA2FybhidnK+9ASABKAlSA2FybhIVCgRuYW1lGI'
    'fmgX8gASgJUgRuYW1lEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZBInCg12ZXJz'
    'aW9uc3RhZ2VzGMGu9WogAygJUg12ZXJzaW9uc3RhZ2Vz');

@$core.Deprecated('Use removeRegionsFromReplicationRequestDescriptor instead')
const RemoveRegionsFromReplicationRequest$json = {
  '1': 'RemoveRegionsFromReplicationRequest',
  '2': [
    {
      '1': 'removereplicaregions',
      '3': 367753789,
      '4': 3,
      '5': 9,
      '10': 'removereplicaregions'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `RemoveRegionsFromReplicationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeRegionsFromReplicationRequestDescriptor =
    $convert.base64Decode(
        'CiNSZW1vdmVSZWdpb25zRnJvbVJlcGxpY2F0aW9uUmVxdWVzdBI2ChRyZW1vdmVyZXBsaWNhcm'
        'VnaW9ucxi99K2vASADKAlSFHJlbW92ZXJlcGxpY2FyZWdpb25zEh4KCHNlY3JldGlkGOXW66IB'
        'IAEoCVIIc2VjcmV0aWQ=');

@$core.Deprecated('Use removeRegionsFromReplicationResponseDescriptor instead')
const RemoveRegionsFromReplicationResponse$json = {
  '1': 'RemoveRegionsFromReplicationResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ReplicationStatusType',
      '10': 'replicationstatus'
    },
  ],
};

/// Descriptor for `RemoveRegionsFromReplicationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeRegionsFromReplicationResponseDescriptor =
    $convert.base64Decode(
        'CiRSZW1vdmVSZWdpb25zRnJvbVJlcGxpY2F0aW9uUmVzcG9uc2USFAoDYXJuGJ2cr70BIAEoCV'
        'IDYXJuElcKEXJlcGxpY2F0aW9uc3RhdHVzGIyqpfwBIAMoCzIlLnNlY3JldHNtYW5hZ2VyLlJl'
        'cGxpY2F0aW9uU3RhdHVzVHlwZVIRcmVwbGljYXRpb25zdGF0dXM=');

@$core.Deprecated('Use replicaRegionTypeDescriptor instead')
const ReplicaRegionType$json = {
  '1': 'ReplicaRegionType',
  '2': [
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
  ],
};

/// Descriptor for `ReplicaRegionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaRegionTypeDescriptor = $convert.base64Decode(
    'ChFSZXBsaWNhUmVnaW9uVHlwZRIdCghrbXNrZXlpZBiNyZcWIAEoCVIIa21za2V5aWQSGQoGcm'
    'VnaW9uGJ7xuUkgASgJUgZyZWdpb24=');

@$core.Deprecated('Use replicateSecretToRegionsRequestDescriptor instead')
const ReplicateSecretToRegionsRequest$json = {
  '1': 'ReplicateSecretToRegionsRequest',
  '2': [
    {
      '1': 'addreplicaregions',
      '3': 461171870,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ReplicaRegionType',
      '10': 'addreplicaregions'
    },
    {
      '1': 'forceoverwritereplicasecret',
      '3': 247407324,
      '4': 1,
      '5': 8,
      '10': 'forceoverwritereplicasecret'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `ReplicateSecretToRegionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicateSecretToRegionsRequestDescriptor = $convert.base64Decode(
    'Ch9SZXBsaWNhdGVTZWNyZXRUb1JlZ2lvbnNSZXF1ZXN0ElMKEWFkZHJlcGxpY2FyZWdpb25zGJ'
    '7Z89sBIAMoCzIhLnNlY3JldHNtYW5hZ2VyLlJlcGxpY2FSZWdpb25UeXBlUhFhZGRyZXBsaWNh'
    'cmVnaW9ucxJDChtmb3JjZW92ZXJ3cml0ZXJlcGxpY2FzZWNyZXQY3MX8dSABKAhSG2ZvcmNlb3'
    'ZlcndyaXRlcmVwbGljYXNlY3JldBIeCghzZWNyZXRpZBjl1uuiASABKAlSCHNlY3JldGlk');

@$core.Deprecated('Use replicateSecretToRegionsResponseDescriptor instead')
const ReplicateSecretToRegionsResponse$json = {
  '1': 'ReplicateSecretToRegionsResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ReplicationStatusType',
      '10': 'replicationstatus'
    },
  ],
};

/// Descriptor for `ReplicateSecretToRegionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicateSecretToRegionsResponseDescriptor =
    $convert.base64Decode(
        'CiBSZXBsaWNhdGVTZWNyZXRUb1JlZ2lvbnNSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm'
        '4SVwoRcmVwbGljYXRpb25zdGF0dXMYjKql/AEgAygLMiUuc2VjcmV0c21hbmFnZXIuUmVwbGlj'
        'YXRpb25TdGF0dXNUeXBlUhFyZXBsaWNhdGlvbnN0YXR1cw==');

@$core.Deprecated('Use replicationStatusTypeDescriptor instead')
const ReplicationStatusType$json = {
  '1': 'ReplicationStatusType',
  '2': [
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastaccesseddate',
      '3': 194418963,
      '4': 1,
      '5': 9,
      '10': 'lastaccesseddate'
    },
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.secretsmanager.StatusType',
      '10': 'status'
    },
    {
      '1': 'statusmessage',
      '3': 72590095,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
  ],
};

/// Descriptor for `ReplicationStatusType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationStatusTypeDescriptor = $convert.base64Decode(
    'ChVSZXBsaWNhdGlvblN0YXR1c1R5cGUSHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEi'
    '0KEGxhc3RhY2Nlc3NlZGRhdGUYk7LaXCABKAlSEGxhc3RhY2Nlc3NlZGRhdGUSGQoGcmVnaW9u'
    'GJ7xuUkgASgJUgZyZWdpb24SNQoGc3RhdHVzGJDk+wIgASgOMhouc2VjcmV0c21hbmFnZXIuU3'
    'RhdHVzVHlwZVIGc3RhdHVzEicKDXN0YXR1c21lc3NhZ2UYj8bOIiABKAlSDXN0YXR1c21lc3Nh'
    'Z2U=');

@$core.Deprecated('Use resourceExistsExceptionDescriptor instead')
const ResourceExistsException$json = {
  '1': 'ResourceExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChdSZXNvdXJjZUV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use restoreSecretRequestDescriptor instead')
const RestoreSecretRequest$json = {
  '1': 'RestoreSecretRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `RestoreSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreSecretRequestDescriptor = $convert.base64Decode(
    'ChRSZXN0b3JlU2VjcmV0UmVxdWVzdBIeCghzZWNyZXRpZBjl1uuiASABKAlSCHNlY3JldGlk');

@$core.Deprecated('Use restoreSecretResponseDescriptor instead')
const RestoreSecretResponse$json = {
  '1': 'RestoreSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `RestoreSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreSecretResponseDescriptor = $convert.base64Decode(
    'ChVSZXN0b3JlU2VjcmV0UmVzcG9uc2USFAoDYXJuGJ2cr70BIAEoCVIDYXJuEhUKBG5hbWUYh+'
    'aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use rotateSecretRequestDescriptor instead')
const RotateSecretRequest$json = {
  '1': 'RotateSecretRequest',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {
      '1': 'externalsecretrotationmetadata',
      '3': 52900542,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ExternalSecretRotationMetadataItem',
      '10': 'externalsecretrotationmetadata'
    },
    {
      '1': 'externalsecretrotationrolearn',
      '3': 470712576,
      '4': 1,
      '5': 9,
      '10': 'externalsecretrotationrolearn'
    },
    {
      '1': 'rotateimmediately',
      '3': 265384053,
      '4': 1,
      '5': 8,
      '10': 'rotateimmediately'
    },
    {
      '1': 'rotationlambdaarn',
      '3': 335026080,
      '4': 1,
      '5': 9,
      '10': 'rotationlambdaarn'
    },
    {
      '1': 'rotationrules',
      '3': 259458135,
      '4': 1,
      '5': 11,
      '6': '.secretsmanager.RotationRulesType',
      '10': 'rotationrules'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `RotateSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rotateSecretRequestDescriptor = $convert.base64Decode(
    'ChNSb3RhdGVTZWNyZXRSZXF1ZXN0EjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhJ9Ch5leHRlcm5hbHNlY3JldHJvdGF0aW9ubWV0YWRhdGEYvuWc'
    'GSADKAsyMi5zZWNyZXRzbWFuYWdlci5FeHRlcm5hbFNlY3JldFJvdGF0aW9uTWV0YWRhdGFJdG'
    'VtUh5leHRlcm5hbHNlY3JldHJvdGF0aW9ubWV0YWRhdGESSAodZXh0ZXJuYWxzZWNyZXRyb3Rh'
    'dGlvbnJvbGVhcm4YgIK64AEgASgJUh1leHRlcm5hbHNlY3JldHJvdGF0aW9ucm9sZWFybhIvCh'
    'Fyb3RhdGVpbW1lZGlhdGVseRj14MV+IAEoCFIRcm90YXRlaW1tZWRpYXRlbHkSMAoRcm90YXRp'
    'b25sYW1iZGFhcm4YoK/gnwEgASgJUhFyb3RhdGlvbmxhbWJkYWFybhJKCg1yb3RhdGlvbnJ1bG'
    'VzGNeI3HsgASgLMiEuc2VjcmV0c21hbmFnZXIuUm90YXRpb25SdWxlc1R5cGVSDXJvdGF0aW9u'
    'cnVsZXMSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZXRpZA==');

@$core.Deprecated('Use rotateSecretResponseDescriptor instead')
const RotateSecretResponse$json = {
  '1': 'RotateSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `RotateSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rotateSecretResponseDescriptor = $convert.base64Decode(
    'ChRSb3RhdGVTZWNyZXRSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm4SFQoEbmFtZRiH5o'
    'F/IAEoCVIEbmFtZRIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use rotationRulesTypeDescriptor instead')
const RotationRulesType$json = {
  '1': 'RotationRulesType',
  '2': [
    {
      '1': 'automaticallyafterdays',
      '3': 350893940,
      '4': 1,
      '5': 3,
      '10': 'automaticallyafterdays'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 9, '10': 'duration'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
  ],
};

/// Descriptor for `RotationRulesType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rotationRulesTypeDescriptor = $convert.base64Decode(
    'ChFSb3RhdGlvblJ1bGVzVHlwZRI6ChZhdXRvbWF0aWNhbGx5YWZ0ZXJkYXlzGPTuqKcBIAEoA1'
    'IWYXV0b21hdGljYWxseWFmdGVyZGF5cxIeCghkdXJhdGlvbhiukp2mASABKAlSCGR1cmF0aW9u'
    'EjIKEnNjaGVkdWxlZXhwcmVzc2lvbhj/kdvUASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbg==');

@$core.Deprecated('Use secretListEntryDescriptor instead')
const SecretListEntry$json = {
  '1': 'SecretListEntry',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'deleteddate', '3': 516314255, '4': 1, '5': 9, '10': 'deleteddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'externalsecretrotationmetadata',
      '3': 52900542,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ExternalSecretRotationMetadataItem',
      '10': 'externalsecretrotationmetadata'
    },
    {
      '1': 'externalsecretrotationrolearn',
      '3': 470712576,
      '4': 1,
      '5': 9,
      '10': 'externalsecretrotationrolearn'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastaccesseddate',
      '3': 194418963,
      '4': 1,
      '5': 9,
      '10': 'lastaccesseddate'
    },
    {
      '1': 'lastchangeddate',
      '3': 314015460,
      '4': 1,
      '5': 9,
      '10': 'lastchangeddate'
    },
    {
      '1': 'lastrotateddate',
      '3': 501475691,
      '4': 1,
      '5': 9,
      '10': 'lastrotateddate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextrotationdate',
      '3': 192035355,
      '4': 1,
      '5': 9,
      '10': 'nextrotationdate'
    },
    {
      '1': 'owningservice',
      '3': 462487817,
      '4': 1,
      '5': 9,
      '10': 'owningservice'
    },
    {
      '1': 'primaryregion',
      '3': 480901186,
      '4': 1,
      '5': 9,
      '10': 'primaryregion'
    },
    {
      '1': 'rotationenabled',
      '3': 209507301,
      '4': 1,
      '5': 8,
      '10': 'rotationenabled'
    },
    {
      '1': 'rotationlambdaarn',
      '3': 335026080,
      '4': 1,
      '5': 9,
      '10': 'rotationlambdaarn'
    },
    {
      '1': 'rotationrules',
      '3': 259458135,
      '4': 1,
      '5': 11,
      '6': '.secretsmanager.RotationRulesType',
      '10': 'rotationrules'
    },
    {
      '1': 'secretversionstostages',
      '3': 356331823,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.SecretListEntry.SecretversionstostagesEntry',
      '10': 'secretversionstostages'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.Tag',
      '10': 'tags'
    },
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
  '3': [SecretListEntry_SecretversionstostagesEntry$json],
};

@$core.Deprecated('Use secretListEntryDescriptor instead')
const SecretListEntry_SecretversionstostagesEntry$json = {
  '1': 'SecretversionstostagesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SecretListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List secretListEntryDescriptor = $convert.base64Decode(
    'Cg9TZWNyZXRMaXN0RW50cnkSFAoDYXJuGJ2cr70BIAEoCVIDYXJuEiQKC2NyZWF0ZWRkYXRlGL'
    'Cw58YBIAEoCVILY3JlYXRlZGRhdGUSJAoLZGVsZXRlZGRhdGUYj6mZ9gEgASgJUgtkZWxldGVk'
    'ZGF0ZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SfQoeZXh0ZXJuYWxzZW'
    'NyZXRyb3RhdGlvbm1ldGFkYXRhGL7lnBkgAygLMjIuc2VjcmV0c21hbmFnZXIuRXh0ZXJuYWxT'
    'ZWNyZXRSb3RhdGlvbk1ldGFkYXRhSXRlbVIeZXh0ZXJuYWxzZWNyZXRyb3RhdGlvbm1ldGFkYX'
    'RhEkgKHWV4dGVybmFsc2VjcmV0cm90YXRpb25yb2xlYXJuGICCuuABIAEoCVIdZXh0ZXJuYWxz'
    'ZWNyZXRyb3RhdGlvbnJvbGVhcm4SHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEi0KEG'
    'xhc3RhY2Nlc3NlZGRhdGUYk7LaXCABKAlSEGxhc3RhY2Nlc3NlZGRhdGUSLAoPbGFzdGNoYW5n'
    'ZWRkYXRlGOT93ZUBIAEoCVIPbGFzdGNoYW5nZWRkYXRlEiwKD2xhc3Ryb3RhdGVkZGF0ZRjr0o'
    '/vASABKAlSD2xhc3Ryb3RhdGVkZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEi0KEG5leHRy'
    'b3RhdGlvbmRhdGUYm/TIWyABKAlSEG5leHRyb3RhdGlvbmRhdGUSKAoNb3duaW5nc2VydmljZR'
    'iJgsTcASABKAlSDW93bmluZ3NlcnZpY2USKAoNcHJpbWFyeXJlZ2lvbhjC8KflASABKAlSDXBy'
    'aW1hcnlyZWdpb24SKwoPcm90YXRpb25lbmFibGVkGOWn82MgASgIUg9yb3RhdGlvbmVuYWJsZW'
    'QSMAoRcm90YXRpb25sYW1iZGFhcm4YoK/gnwEgASgJUhFyb3RhdGlvbmxhbWJkYWFybhJKCg1y'
    'b3RhdGlvbnJ1bGVzGNeI3HsgASgLMiEuc2VjcmV0c21hbmFnZXIuUm90YXRpb25SdWxlc1R5cG'
    'VSDXJvdGF0aW9ucnVsZXMSdwoWc2VjcmV0dmVyc2lvbnN0b3N0YWdlcxiv4vSpASADKAsyOy5z'
    'ZWNyZXRzbWFuYWdlci5TZWNyZXRMaXN0RW50cnkuU2VjcmV0dmVyc2lvbnN0b3N0YWdlc0VudH'
    'J5UhZzZWNyZXR2ZXJzaW9uc3Rvc3RhZ2VzEisKBHRhZ3MYwcH2tQEgAygLMhMuc2VjcmV0c21h'
    'bmFnZXIuVGFnUgR0YWdzEhYKBHR5cGUY7qDXigEgASgJUgR0eXBlGkkKG1NlY3JldHZlcnNpb2'
    '5zdG9zdGFnZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6'
    'AjgB');

@$core.Deprecated('Use secretValueEntryDescriptor instead')
const SecretValueEntry$json = {
  '1': 'SecretValueEntry',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'secretbinary', '3': 94375681, '4': 1, '5': 12, '10': 'secretbinary'},
    {'1': 'secretstring', '3': 190782253, '4': 1, '5': 9, '10': 'secretstring'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {
      '1': 'versionstages',
      '3': 224220993,
      '4': 3,
      '5': 9,
      '10': 'versionstages'
    },
  ],
};

/// Descriptor for `SecretValueEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List secretValueEntryDescriptor = $convert.base64Decode(
    'ChBTZWNyZXRWYWx1ZUVudHJ5EhQKA2FybhidnK+9ASABKAlSA2FybhIkCgtjcmVhdGVkZGF0ZR'
    'iwsOfGASABKAlSC2NyZWF0ZWRkYXRlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJQoMc2VjcmV0'
    'YmluYXJ5GIGegC0gASgMUgxzZWNyZXRiaW5hcnkSJQoMc2VjcmV0c3RyaW5nGK22/FogASgJUg'
    'xzZWNyZXRzdHJpbmcSIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlkEicKDXZlcnNp'
    'b25zdGFnZXMYwa71aiADKAlSDXZlcnNpb25zdGFnZXM=');

@$core.Deprecated('Use secretVersionsListEntryDescriptor instead')
const SecretVersionsListEntry$json = {
  '1': 'SecretVersionsListEntry',
  '2': [
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'kmskeyids', '3': 478641518, '4': 3, '5': 9, '10': 'kmskeyids'},
    {
      '1': 'lastaccesseddate',
      '3': 194418963,
      '4': 1,
      '5': 9,
      '10': 'lastaccesseddate'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {
      '1': 'versionstages',
      '3': 224220993,
      '4': 3,
      '5': 9,
      '10': 'versionstages'
    },
  ],
};

/// Descriptor for `SecretVersionsListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List secretVersionsListEntryDescriptor = $convert.base64Decode(
    'ChdTZWNyZXRWZXJzaW9uc0xpc3RFbnRyeRIkCgtjcmVhdGVkZGF0ZRiwsOfGASABKAlSC2NyZW'
    'F0ZWRkYXRlEiAKCWttc2tleWlkcxju+p3kASADKAlSCWttc2tleWlkcxItChBsYXN0YWNjZXNz'
    'ZWRkYXRlGJOy2lwgASgJUhBsYXN0YWNjZXNzZWRkYXRlEiAKCXZlcnNpb25pZBib4ZmhASABKA'
    'lSCXZlcnNpb25pZBInCg12ZXJzaW9uc3RhZ2VzGMGu9WogAygJUg12ZXJzaW9uc3RhZ2Vz');

@$core.Deprecated('Use stopReplicationToReplicaRequestDescriptor instead')
const StopReplicationToReplicaRequest$json = {
  '1': 'StopReplicationToReplicaRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `StopReplicationToReplicaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopReplicationToReplicaRequestDescriptor =
    $convert.base64Decode(
        'Ch9TdG9wUmVwbGljYXRpb25Ub1JlcGxpY2FSZXF1ZXN0Eh4KCHNlY3JldGlkGOXW66IBIAEoCV'
        'IIc2VjcmV0aWQ=');

@$core.Deprecated('Use stopReplicationToReplicaResponseDescriptor instead')
const StopReplicationToReplicaResponse$json = {
  '1': 'StopReplicationToReplicaResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `StopReplicationToReplicaResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopReplicationToReplicaResponseDescriptor =
    $convert.base64Decode(
        'CiBTdG9wUmVwbGljYXRpb25Ub1JlcGxpY2FSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm'
        '4=');

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
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZXRpZBIrCg'
    'R0YWdzGMHB9rUBIAMoCzITLnNlY3JldHNtYW5hZ2VyLlRhZ1IEdGFncw==');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIeCghzZWNyZXRpZBjl1uuiASABKAlSCHNlY3JldGlkEh'
    'wKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use updateSecretRequestDescriptor instead')
const UpdateSecretRequest$json = {
  '1': 'UpdateSecretRequest',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'secretbinary', '3': 94375681, '4': 1, '5': 12, '10': 'secretbinary'},
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
    {'1': 'secretstring', '3': 190782253, '4': 1, '5': 9, '10': 'secretstring'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `UpdateSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSecretRequestDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVTZWNyZXRSZXF1ZXN0EjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24S'
    'HQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEiUKDHNlY3JldGJpbmFyeRiBnoAtIAEoDF'
    'IMc2VjcmV0YmluYXJ5Eh4KCHNlY3JldGlkGOXW66IBIAEoCVIIc2VjcmV0aWQSJQoMc2VjcmV0'
    'c3RyaW5nGK22/FogASgJUgxzZWNyZXRzdHJpbmcSFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use updateSecretResponseDescriptor instead')
const UpdateSecretResponse$json = {
  '1': 'UpdateSecretResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `UpdateSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSecretResponseDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVTZWNyZXRSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm4SFQoEbmFtZRiH5o'
    'F/IAEoCVIEbmFtZRIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use updateSecretVersionStageRequestDescriptor instead')
const UpdateSecretVersionStageRequest$json = {
  '1': 'UpdateSecretVersionStageRequest',
  '2': [
    {
      '1': 'movetoversionid',
      '3': 509017411,
      '4': 1,
      '5': 9,
      '10': 'movetoversionid'
    },
    {
      '1': 'removefromversionid',
      '3': 194704147,
      '4': 1,
      '5': 9,
      '10': 'removefromversionid'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
    {'1': 'versionstage', '3': 229692340, '4': 1, '5': 9, '10': 'versionstage'},
  ],
};

/// Descriptor for `UpdateSecretVersionStageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSecretVersionStageRequestDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVTZWNyZXRWZXJzaW9uU3RhZ2VSZXF1ZXN0EiwKD21vdmV0b3ZlcnNpb25pZBjD+t'
        'vyASABKAlSD21vdmV0b3ZlcnNpb25pZBIzChNyZW1vdmVmcm9tdmVyc2lvbmlkGJPm61wgASgJ'
        'UhNyZW1vdmVmcm9tdmVyc2lvbmlkEh4KCHNlY3JldGlkGOXW66IBIAEoCVIIc2VjcmV0aWQSJQ'
        'oMdmVyc2lvbnN0YWdlGLSnw20gASgJUgx2ZXJzaW9uc3RhZ2U=');

@$core.Deprecated('Use updateSecretVersionStageResponseDescriptor instead')
const UpdateSecretVersionStageResponse$json = {
  '1': 'UpdateSecretVersionStageResponse',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateSecretVersionStageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSecretVersionStageResponseDescriptor =
    $convert.base64Decode(
        'CiBVcGRhdGVTZWNyZXRWZXJzaW9uU3RhZ2VSZXNwb25zZRIUCgNhcm4YnZyvvQEgASgJUgNhcm'
        '4SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use validateResourcePolicyRequestDescriptor instead')
const ValidateResourcePolicyRequest$json = {
  '1': 'ValidateResourcePolicyRequest',
  '2': [
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
    {'1': 'secretid', '3': 341502821, '4': 1, '5': 9, '10': 'secretid'},
  ],
};

/// Descriptor for `ValidateResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'Ch1WYWxpZGF0ZVJlc291cmNlUG9saWN5UmVxdWVzdBIpCg5yZXNvdXJjZXBvbGljeRiwlMEHIA'
        'EoCVIOcmVzb3VyY2Vwb2xpY3kSHgoIc2VjcmV0aWQY5dbrogEgASgJUghzZWNyZXRpZA==');

@$core.Deprecated('Use validateResourcePolicyResponseDescriptor instead')
const ValidateResourcePolicyResponse$json = {
  '1': 'ValidateResourcePolicyResponse',
  '2': [
    {
      '1': 'policyvalidationpassed',
      '3': 298018737,
      '4': 1,
      '5': 8,
      '10': 'policyvalidationpassed'
    },
    {
      '1': 'validationerrors',
      '3': 381330622,
      '4': 3,
      '5': 11,
      '6': '.secretsmanager.ValidationErrorsEntry',
      '10': 'validationerrors'
    },
  ],
};

/// Descriptor for `ValidateResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateResourcePolicyResponseDescriptor =
    $convert.base64Decode(
        'Ch5WYWxpZGF0ZVJlc291cmNlUG9saWN5UmVzcG9uc2USOgoWcG9saWN5dmFsaWRhdGlvbnBhc3'
        'NlZBixz42OASABKAhSFnBvbGljeXZhbGlkYXRpb25wYXNzZWQSVQoQdmFsaWRhdGlvbmVycm9y'
        'cxi+yeq1ASADKAsyJS5zZWNyZXRzbWFuYWdlci5WYWxpZGF0aW9uRXJyb3JzRW50cnlSEHZhbG'
        'lkYXRpb25lcnJvcnM=');

@$core.Deprecated('Use validationErrorsEntryDescriptor instead')
const ValidationErrorsEntry$json = {
  '1': 'ValidationErrorsEntry',
  '2': [
    {'1': 'checkname', '3': 68143813, '4': 1, '5': 9, '10': 'checkname'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
  ],
};

/// Descriptor for `ValidationErrorsEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationErrorsEntryDescriptor = $convert.base64Decode(
    'ChVWYWxpZGF0aW9uRXJyb3JzRW50cnkSHwoJY2hlY2tuYW1lGMWVvyAgASgJUgljaGVja25hbW'
    'USJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdl');
