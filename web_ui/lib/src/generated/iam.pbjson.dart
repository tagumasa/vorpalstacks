// This is a generated file - do not edit.
//
// Generated from iam.proto.

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

@$core.Deprecated('Use accessAdvisorUsageGranularityTypeDescriptor instead')
const AccessAdvisorUsageGranularityType$json = {
  '1': 'AccessAdvisorUsageGranularityType',
  '2': [
    {'1': 'ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_SERVICE_LEVEL', '2': 0},
    {'1': 'ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_ACTION_LEVEL', '2': 1},
  ],
};

/// Descriptor for `AccessAdvisorUsageGranularityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accessAdvisorUsageGranularityTypeDescriptor =
    $convert.base64Decode(
        'CiFBY2Nlc3NBZHZpc29yVXNhZ2VHcmFudWxhcml0eVR5cGUSNwozQUNDRVNTX0FEVklTT1JfVV'
        'NBR0VfR1JBTlVMQVJJVFlfVFlQRV9TRVJWSUNFX0xFVkVMEAASNgoyQUNDRVNTX0FEVklTT1Jf'
        'VVNBR0VfR1JBTlVMQVJJVFlfVFlQRV9BQ1RJT05fTEVWRUwQAQ==');

@$core.Deprecated('Use contextKeyTypeEnumDescriptor instead')
const ContextKeyTypeEnum$json = {
  '1': 'ContextKeyTypeEnum',
  '2': [
    {'1': 'CONTEXT_KEY_TYPE_ENUM_NUMERIC_LIST', '2': 0},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_BOOLEAN_LIST', '2': 1},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_BINARY', '2': 2},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_DATE_LIST', '2': 3},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_NUMERIC', '2': 4},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_STRING_LIST', '2': 5},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_STRING', '2': 6},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_IP', '2': 7},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_BOOLEAN', '2': 8},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_IP_LIST', '2': 9},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_DATE', '2': 10},
    {'1': 'CONTEXT_KEY_TYPE_ENUM_BINARY_LIST', '2': 11},
  ],
};

/// Descriptor for `ContextKeyTypeEnum`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contextKeyTypeEnumDescriptor = $convert.base64Decode(
    'ChJDb250ZXh0S2V5VHlwZUVudW0SJgoiQ09OVEVYVF9LRVlfVFlQRV9FTlVNX05VTUVSSUNfTE'
    'lTVBAAEiYKIkNPTlRFWFRfS0VZX1RZUEVfRU5VTV9CT09MRUFOX0xJU1QQARIgChxDT05URVhU'
    'X0tFWV9UWVBFX0VOVU1fQklOQVJZEAISIwofQ09OVEVYVF9LRVlfVFlQRV9FTlVNX0RBVEVfTE'
    'lTVBADEiEKHUNPTlRFWFRfS0VZX1RZUEVfRU5VTV9OVU1FUklDEAQSJQohQ09OVEVYVF9LRVlf'
    'VFlQRV9FTlVNX1NUUklOR19MSVNUEAUSIAocQ09OVEVYVF9LRVlfVFlQRV9FTlVNX1NUUklORx'
    'AGEhwKGENPTlRFWFRfS0VZX1RZUEVfRU5VTV9JUBAHEiEKHUNPTlRFWFRfS0VZX1RZUEVfRU5V'
    'TV9CT09MRUFOEAgSIQodQ09OVEVYVF9LRVlfVFlQRV9FTlVNX0lQX0xJU1QQCRIeChpDT05URV'
    'hUX0tFWV9UWVBFX0VOVU1fREFURRAKEiUKIUNPTlRFWFRfS0VZX1RZUEVfRU5VTV9CSU5BUllf'
    'TElTVBAL');

@$core.Deprecated('Use deletionTaskStatusTypeDescriptor instead')
const DeletionTaskStatusType$json = {
  '1': 'DeletionTaskStatusType',
  '2': [
    {'1': 'DELETION_TASK_STATUS_TYPE_IN_PROGRESS', '2': 0},
    {'1': 'DELETION_TASK_STATUS_TYPE_SUCCEEDED', '2': 1},
    {'1': 'DELETION_TASK_STATUS_TYPE_FAILED', '2': 2},
    {'1': 'DELETION_TASK_STATUS_TYPE_NOT_STARTED', '2': 3},
  ],
};

/// Descriptor for `DeletionTaskStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deletionTaskStatusTypeDescriptor = $convert.base64Decode(
    'ChZEZWxldGlvblRhc2tTdGF0dXNUeXBlEikKJURFTEVUSU9OX1RBU0tfU1RBVFVTX1RZUEVfSU'
    '5fUFJPR1JFU1MQABInCiNERUxFVElPTl9UQVNLX1NUQVRVU19UWVBFX1NVQ0NFRURFRBABEiQK'
    'IERFTEVUSU9OX1RBU0tfU1RBVFVTX1RZUEVfRkFJTEVEEAISKQolREVMRVRJT05fVEFTS19TVE'
    'FUVVNfVFlQRV9OT1RfU1RBUlRFRBAD');

@$core.Deprecated('Use entityTypeDescriptor instead')
const EntityType$json = {
  '1': 'EntityType',
  '2': [
    {'1': 'ENTITY_TYPE_GROUP', '2': 0},
    {'1': 'ENTITY_TYPE_LOCALMANAGEDPOLICY', '2': 1},
    {'1': 'ENTITY_TYPE_USER', '2': 2},
    {'1': 'ENTITY_TYPE_AWSMANAGEDPOLICY', '2': 3},
    {'1': 'ENTITY_TYPE_ROLE', '2': 4},
  ],
};

/// Descriptor for `EntityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List entityTypeDescriptor = $convert.base64Decode(
    'CgpFbnRpdHlUeXBlEhUKEUVOVElUWV9UWVBFX0dST1VQEAASIgoeRU5USVRZX1RZUEVfTE9DQU'
    'xNQU5BR0VEUE9MSUNZEAESFAoQRU5USVRZX1RZUEVfVVNFUhACEiAKHEVOVElUWV9UWVBFX0FX'
    'U01BTkFHRURQT0xJQ1kQAxIUChBFTlRJVFlfVFlQRV9ST0xFEAQ=');

@$core.Deprecated('Use featureTypeDescriptor instead')
const FeatureType$json = {
  '1': 'FeatureType',
  '2': [
    {'1': 'FEATURE_TYPE_ROOT_SESSIONS', '2': 0},
    {'1': 'FEATURE_TYPE_ROOT_CREDENTIALS_MANAGEMENT', '2': 1},
  ],
};

/// Descriptor for `FeatureType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List featureTypeDescriptor = $convert.base64Decode(
    'CgtGZWF0dXJlVHlwZRIeChpGRUFUVVJFX1RZUEVfUk9PVF9TRVNTSU9OUxAAEiwKKEZFQVRVUk'
    'VfVFlQRV9ST09UX0NSRURFTlRJQUxTX01BTkFHRU1FTlQQAQ==');

@$core.Deprecated('Use permissionsBoundaryAttachmentTypeDescriptor instead')
const PermissionsBoundaryAttachmentType$json = {
  '1': 'PermissionsBoundaryAttachmentType',
  '2': [
    {'1': 'PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY', '2': 0},
  ],
};

/// Descriptor for `PermissionsBoundaryAttachmentType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List permissionsBoundaryAttachmentTypeDescriptor =
    $convert.base64Decode(
        'CiFQZXJtaXNzaW9uc0JvdW5kYXJ5QXR0YWNobWVudFR5cGUSLworUEVSTUlTU0lPTlNfQk9VTk'
        'RBUllfQVRUQUNITUVOVF9UWVBFX1BPTElDWRAA');

@$core.Deprecated('Use policyEvaluationDecisionTypeDescriptor instead')
const PolicyEvaluationDecisionType$json = {
  '1': 'PolicyEvaluationDecisionType',
  '2': [
    {'1': 'POLICY_EVALUATION_DECISION_TYPE_IMPLICIT_DENY', '2': 0},
    {'1': 'POLICY_EVALUATION_DECISION_TYPE_EXPLICIT_DENY', '2': 1},
    {'1': 'POLICY_EVALUATION_DECISION_TYPE_ALLOWED', '2': 2},
  ],
};

/// Descriptor for `PolicyEvaluationDecisionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyEvaluationDecisionTypeDescriptor = $convert.base64Decode(
    'ChxQb2xpY3lFdmFsdWF0aW9uRGVjaXNpb25UeXBlEjEKLVBPTElDWV9FVkFMVUFUSU9OX0RFQ0'
    'lTSU9OX1RZUEVfSU1QTElDSVRfREVOWRAAEjEKLVBPTElDWV9FVkFMVUFUSU9OX0RFQ0lTSU9O'
    'X1RZUEVfRVhQTElDSVRfREVOWRABEisKJ1BPTElDWV9FVkFMVUFUSU9OX0RFQ0lTSU9OX1RZUE'
    'VfQUxMT1dFRBAC');

@$core.Deprecated('Use policyParameterTypeEnumDescriptor instead')
const PolicyParameterTypeEnum$json = {
  '1': 'PolicyParameterTypeEnum',
  '2': [
    {'1': 'POLICY_PARAMETER_TYPE_ENUM_STRING_LIST', '2': 0},
    {'1': 'POLICY_PARAMETER_TYPE_ENUM_STRING', '2': 1},
  ],
};

/// Descriptor for `PolicyParameterTypeEnum`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyParameterTypeEnumDescriptor = $convert.base64Decode(
    'ChdQb2xpY3lQYXJhbWV0ZXJUeXBlRW51bRIqCiZQT0xJQ1lfUEFSQU1FVEVSX1RZUEVfRU5VTV'
    '9TVFJJTkdfTElTVBAAEiUKIVBPTElDWV9QQVJBTUVURVJfVFlQRV9FTlVNX1NUUklORxAB');

@$core.Deprecated('Use policySourceTypeDescriptor instead')
const PolicySourceType$json = {
  '1': 'PolicySourceType',
  '2': [
    {'1': 'POLICY_SOURCE_TYPE_GROUP', '2': 0},
    {'1': 'POLICY_SOURCE_TYPE_USER_MANAGED', '2': 1},
    {'1': 'POLICY_SOURCE_TYPE_AWS_MANAGED', '2': 2},
    {'1': 'POLICY_SOURCE_TYPE_NONE', '2': 3},
    {'1': 'POLICY_SOURCE_TYPE_RESOURCE', '2': 4},
    {'1': 'POLICY_SOURCE_TYPE_ROLE', '2': 5},
    {'1': 'POLICY_SOURCE_TYPE_USER', '2': 6},
  ],
};

/// Descriptor for `PolicySourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policySourceTypeDescriptor = $convert.base64Decode(
    'ChBQb2xpY3lTb3VyY2VUeXBlEhwKGFBPTElDWV9TT1VSQ0VfVFlQRV9HUk9VUBAAEiMKH1BPTE'
    'lDWV9TT1VSQ0VfVFlQRV9VU0VSX01BTkFHRUQQARIiCh5QT0xJQ1lfU09VUkNFX1RZUEVfQVdT'
    'X01BTkFHRUQQAhIbChdQT0xJQ1lfU09VUkNFX1RZUEVfTk9ORRADEh8KG1BPTElDWV9TT1VSQ0'
    'VfVFlQRV9SRVNPVVJDRRAEEhsKF1BPTElDWV9TT1VSQ0VfVFlQRV9ST0xFEAUSGwoXUE9MSUNZ'
    'X1NPVVJDRV9UWVBFX1VTRVIQBg==');

@$core.Deprecated('Use policyUsageTypeDescriptor instead')
const PolicyUsageType$json = {
  '1': 'PolicyUsageType',
  '2': [
    {'1': 'POLICY_USAGE_TYPE_PERMISSIONSPOLICY', '2': 0},
    {'1': 'POLICY_USAGE_TYPE_PERMISSIONSBOUNDARY', '2': 1},
  ],
};

/// Descriptor for `PolicyUsageType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyUsageTypeDescriptor = $convert.base64Decode(
    'Cg9Qb2xpY3lVc2FnZVR5cGUSJwojUE9MSUNZX1VTQUdFX1RZUEVfUEVSTUlTU0lPTlNQT0xJQ1'
    'kQABIpCiVQT0xJQ1lfVVNBR0VfVFlQRV9QRVJNSVNTSU9OU0JPVU5EQVJZEAE=');

@$core.Deprecated('Use reportFormatTypeDescriptor instead')
const ReportFormatType$json = {
  '1': 'ReportFormatType',
  '2': [
    {'1': 'REPORT_FORMAT_TYPE_TEXT_CSV', '2': 0},
  ],
};

/// Descriptor for `ReportFormatType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reportFormatTypeDescriptor = $convert.base64Decode(
    'ChBSZXBvcnRGb3JtYXRUeXBlEh8KG1JFUE9SVF9GT1JNQVRfVFlQRV9URVhUX0NTVhAA');

@$core.Deprecated('Use reportStateTypeDescriptor instead')
const ReportStateType$json = {
  '1': 'ReportStateType',
  '2': [
    {'1': 'REPORT_STATE_TYPE_STARTED', '2': 0},
    {'1': 'REPORT_STATE_TYPE_INPROGRESS', '2': 1},
    {'1': 'REPORT_STATE_TYPE_COMPLETE', '2': 2},
  ],
};

/// Descriptor for `ReportStateType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reportStateTypeDescriptor = $convert.base64Decode(
    'Cg9SZXBvcnRTdGF0ZVR5cGUSHQoZUkVQT1JUX1NUQVRFX1RZUEVfU1RBUlRFRBAAEiAKHFJFUE'
    '9SVF9TVEFURV9UWVBFX0lOUFJPR1JFU1MQARIeChpSRVBPUlRfU1RBVEVfVFlQRV9DT01QTEVU'
    'RRAC');

@$core.Deprecated('Use assertionEncryptionModeTypeDescriptor instead')
const assertionEncryptionModeType$json = {
  '1': 'assertionEncryptionModeType',
  '2': [
    {'1': 'ASSERTION_ENCRYPTION_MODE_TYPE_REQUIRED', '2': 0},
    {'1': 'ASSERTION_ENCRYPTION_MODE_TYPE_ALLOWED', '2': 1},
  ],
};

/// Descriptor for `assertionEncryptionModeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List assertionEncryptionModeTypeDescriptor =
    $convert.base64Decode(
        'Chthc3NlcnRpb25FbmNyeXB0aW9uTW9kZVR5cGUSKwonQVNTRVJUSU9OX0VOQ1JZUFRJT05fTU'
        '9ERV9UWVBFX1JFUVVJUkVEEAASKgomQVNTRVJUSU9OX0VOQ1JZUFRJT05fTU9ERV9UWVBFX0FM'
        'TE9XRUQQAQ==');

@$core.Deprecated('Use assignmentStatusTypeDescriptor instead')
const assignmentStatusType$json = {
  '1': 'assignmentStatusType',
  '2': [
    {'1': 'ASSIGNMENT_STATUS_TYPE_ANY', '2': 0},
    {'1': 'ASSIGNMENT_STATUS_TYPE_UNASSIGNED', '2': 1},
    {'1': 'ASSIGNMENT_STATUS_TYPE_ASSIGNED', '2': 2},
  ],
};

/// Descriptor for `assignmentStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List assignmentStatusTypeDescriptor = $convert.base64Decode(
    'ChRhc3NpZ25tZW50U3RhdHVzVHlwZRIeChpBU1NJR05NRU5UX1NUQVRVU19UWVBFX0FOWRAAEi'
    'UKIUFTU0lHTk1FTlRfU1RBVFVTX1RZUEVfVU5BU1NJR05FRBABEiMKH0FTU0lHTk1FTlRfU1RB'
    'VFVTX1RZUEVfQVNTSUdORUQQAg==');

@$core.Deprecated('Use encodingTypeDescriptor instead')
const encodingType$json = {
  '1': 'encodingType',
  '2': [
    {'1': 'ENCODING_TYPE_PEM', '2': 0},
    {'1': 'ENCODING_TYPE_SSH', '2': 1},
  ],
};

/// Descriptor for `encodingType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encodingTypeDescriptor = $convert.base64Decode(
    'CgxlbmNvZGluZ1R5cGUSFQoRRU5DT0RJTkdfVFlQRV9QRU0QABIVChFFTkNPRElOR19UWVBFX1'
    'NTSBAB');

@$core.Deprecated('Use globalEndpointTokenVersionDescriptor instead')
const globalEndpointTokenVersion$json = {
  '1': 'globalEndpointTokenVersion',
  '2': [
    {'1': 'GLOBAL_ENDPOINT_TOKEN_VERSION_V2TOKEN', '2': 0},
    {'1': 'GLOBAL_ENDPOINT_TOKEN_VERSION_V1TOKEN', '2': 1},
  ],
};

/// Descriptor for `globalEndpointTokenVersion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List globalEndpointTokenVersionDescriptor =
    $convert.base64Decode(
        'ChpnbG9iYWxFbmRwb2ludFRva2VuVmVyc2lvbhIpCiVHTE9CQUxfRU5EUE9JTlRfVE9LRU5fVk'
        'VSU0lPTl9WMlRPS0VOEAASKQolR0xPQkFMX0VORFBPSU5UX1RPS0VOX1ZFUlNJT05fVjFUT0tF'
        'ThAB');

@$core.Deprecated('Use jobStatusTypeDescriptor instead')
const jobStatusType$json = {
  '1': 'jobStatusType',
  '2': [
    {'1': 'JOB_STATUS_TYPE_IN_PROGRESS', '2': 0},
    {'1': 'JOB_STATUS_TYPE_COMPLETED', '2': 1},
    {'1': 'JOB_STATUS_TYPE_FAILED', '2': 2},
  ],
};

/// Descriptor for `jobStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List jobStatusTypeDescriptor = $convert.base64Decode(
    'Cg1qb2JTdGF0dXNUeXBlEh8KG0pPQl9TVEFUVVNfVFlQRV9JTl9QUk9HUkVTUxAAEh0KGUpPQl'
    '9TVEFUVVNfVFlQRV9DT01QTEVURUQQARIaChZKT0JfU1RBVFVTX1RZUEVfRkFJTEVEEAI=');

@$core.Deprecated('Use permissionCheckResultTypeDescriptor instead')
const permissionCheckResultType$json = {
  '1': 'permissionCheckResultType',
  '2': [
    {'1': 'PERMISSION_CHECK_RESULT_TYPE_UNSURE', '2': 0},
    {'1': 'PERMISSION_CHECK_RESULT_TYPE_ALLOWED', '2': 1},
    {'1': 'PERMISSION_CHECK_RESULT_TYPE_DENIED', '2': 2},
  ],
};

/// Descriptor for `permissionCheckResultType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List permissionCheckResultTypeDescriptor = $convert.base64Decode(
    'ChlwZXJtaXNzaW9uQ2hlY2tSZXN1bHRUeXBlEicKI1BFUk1JU1NJT05fQ0hFQ0tfUkVTVUxUX1'
    'RZUEVfVU5TVVJFEAASKAokUEVSTUlTU0lPTl9DSEVDS19SRVNVTFRfVFlQRV9BTExPV0VEEAES'
    'JwojUEVSTUlTU0lPTl9DSEVDS19SRVNVTFRfVFlQRV9ERU5JRUQQAg==');

@$core.Deprecated('Use permissionCheckStatusTypeDescriptor instead')
const permissionCheckStatusType$json = {
  '1': 'permissionCheckStatusType',
  '2': [
    {'1': 'PERMISSION_CHECK_STATUS_TYPE_IN_PROGRESS', '2': 0},
    {'1': 'PERMISSION_CHECK_STATUS_TYPE_COMPLETE', '2': 1},
    {'1': 'PERMISSION_CHECK_STATUS_TYPE_FAILED', '2': 2},
  ],
};

/// Descriptor for `permissionCheckStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List permissionCheckStatusTypeDescriptor = $convert.base64Decode(
    'ChlwZXJtaXNzaW9uQ2hlY2tTdGF0dXNUeXBlEiwKKFBFUk1JU1NJT05fQ0hFQ0tfU1RBVFVTX1'
    'RZUEVfSU5fUFJPR1JFU1MQABIpCiVQRVJNSVNTSU9OX0NIRUNLX1NUQVRVU19UWVBFX0NPTVBM'
    'RVRFEAESJwojUEVSTUlTU0lPTl9DSEVDS19TVEFUVVNfVFlQRV9GQUlMRUQQAg==');

@$core.Deprecated('Use policyOwnerEntityTypeDescriptor instead')
const policyOwnerEntityType$json = {
  '1': 'policyOwnerEntityType',
  '2': [
    {'1': 'POLICY_OWNER_ENTITY_TYPE_GROUP', '2': 0},
    {'1': 'POLICY_OWNER_ENTITY_TYPE_ROLE', '2': 1},
    {'1': 'POLICY_OWNER_ENTITY_TYPE_USER', '2': 2},
  ],
};

/// Descriptor for `policyOwnerEntityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyOwnerEntityTypeDescriptor = $convert.base64Decode(
    'ChVwb2xpY3lPd25lckVudGl0eVR5cGUSIgoeUE9MSUNZX09XTkVSX0VOVElUWV9UWVBFX0dST1'
    'VQEAASIQodUE9MSUNZX09XTkVSX0VOVElUWV9UWVBFX1JPTEUQARIhCh1QT0xJQ1lfT1dORVJf'
    'RU5USVRZX1RZUEVfVVNFUhAC');

@$core.Deprecated('Use policyScopeTypeDescriptor instead')
const policyScopeType$json = {
  '1': 'policyScopeType',
  '2': [
    {'1': 'POLICY_SCOPE_TYPE_LOCAL', '2': 0},
    {'1': 'POLICY_SCOPE_TYPE_ALL', '2': 1},
    {'1': 'POLICY_SCOPE_TYPE_AWS', '2': 2},
  ],
};

/// Descriptor for `policyScopeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyScopeTypeDescriptor = $convert.base64Decode(
    'Cg9wb2xpY3lTY29wZVR5cGUSGwoXUE9MSUNZX1NDT1BFX1RZUEVfTE9DQUwQABIZChVQT0xJQ1'
    'lfU0NPUEVfVFlQRV9BTEwQARIZChVQT0xJQ1lfU0NPUEVfVFlQRV9BV1MQAg==');

@$core.Deprecated('Use policyTypeDescriptor instead')
const policyType$json = {
  '1': 'policyType',
  '2': [
    {'1': 'POLICY_TYPE_INLINE', '2': 0},
    {'1': 'POLICY_TYPE_MANAGED', '2': 1},
  ],
};

/// Descriptor for `policyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyTypeDescriptor = $convert.base64Decode(
    'Cgpwb2xpY3lUeXBlEhYKElBPTElDWV9UWVBFX0lOTElORRAAEhcKE1BPTElDWV9UWVBFX01BTk'
    'FHRUQQAQ==');

@$core.Deprecated('Use sortKeyTypeDescriptor instead')
const sortKeyType$json = {
  '1': 'sortKeyType',
  '2': [
    {'1': 'SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_DESCENDING', '2': 0},
    {'1': 'SORT_KEY_TYPE_SERVICE_NAMESPACE_DESCENDING', '2': 1},
    {'1': 'SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_ASCENDING', '2': 2},
    {'1': 'SORT_KEY_TYPE_SERVICE_NAMESPACE_ASCENDING', '2': 3},
  ],
};

/// Descriptor for `sortKeyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sortKeyTypeDescriptor = $convert.base64Decode(
    'Cgtzb3J0S2V5VHlwZRI0CjBTT1JUX0tFWV9UWVBFX0xBU1RfQVVUSEVOVElDQVRFRF9USU1FX0'
    'RFU0NFTkRJTkcQABIuCipTT1JUX0tFWV9UWVBFX1NFUlZJQ0VfTkFNRVNQQUNFX0RFU0NFTkRJ'
    'TkcQARIzCi9TT1JUX0tFWV9UWVBFX0xBU1RfQVVUSEVOVElDQVRFRF9USU1FX0FTQ0VORElORx'
    'ACEi0KKVNPUlRfS0VZX1RZUEVfU0VSVklDRV9OQU1FU1BBQ0VfQVNDRU5ESU5HEAM=');

@$core.Deprecated('Use stateTypeDescriptor instead')
const stateType$json = {
  '1': 'stateType',
  '2': [
    {'1': 'STATE_TYPE_ACCEPTED', '2': 0},
    {'1': 'STATE_TYPE_PENDING_APPROVAL', '2': 1},
    {'1': 'STATE_TYPE_FINALIZED', '2': 2},
    {'1': 'STATE_TYPE_UNASSIGNED', '2': 3},
    {'1': 'STATE_TYPE_REJECTED', '2': 4},
    {'1': 'STATE_TYPE_EXPIRED', '2': 5},
    {'1': 'STATE_TYPE_ASSIGNED', '2': 6},
  ],
};

/// Descriptor for `stateType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateTypeDescriptor = $convert.base64Decode(
    'CglzdGF0ZVR5cGUSFwoTU1RBVEVfVFlQRV9BQ0NFUFRFRBAAEh8KG1NUQVRFX1RZUEVfUEVORE'
    'lOR19BUFBST1ZBTBABEhgKFFNUQVRFX1RZUEVfRklOQUxJWkVEEAISGQoVU1RBVEVfVFlQRV9V'
    'TkFTU0lHTkVEEAMSFwoTU1RBVEVfVFlQRV9SRUpFQ1RFRBAEEhYKElNUQVRFX1RZUEVfRVhQSV'
    'JFRBAFEhcKE1NUQVRFX1RZUEVfQVNTSUdORUQQBg==');

@$core.Deprecated('Use statusTypeDescriptor instead')
const statusType$json = {
  '1': 'statusType',
  '2': [
    {'1': 'STATUS_TYPE_ACTIVE', '2': 0},
    {'1': 'STATUS_TYPE_EXPIRED', '2': 1},
    {'1': 'STATUS_TYPE_INACTIVE', '2': 2},
  ],
};

/// Descriptor for `statusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statusTypeDescriptor = $convert.base64Decode(
    'CgpzdGF0dXNUeXBlEhYKElNUQVRVU19UWVBFX0FDVElWRRAAEhcKE1NUQVRVU19UWVBFX0VYUE'
    'lSRUQQARIYChRTVEFUVVNfVFlQRV9JTkFDVElWRRAC');

@$core.Deprecated('Use summaryKeyTypeDescriptor instead')
const summaryKeyType$json = {
  '1': 'summaryKeyType',
  '2': [
    {'1': 'SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERGROUPQUOTA', '2': 0},
    {'1': 'SUMMARY_KEY_TYPE_ROLES', '2': 1},
    {'1': 'SUMMARY_KEY_TYPE_VERSIONSPERPOLICYQUOTA', '2': 2},
    {'1': 'SUMMARY_KEY_TYPE_MFADEVICES', '2': 3},
    {'1': 'SUMMARY_KEY_TYPE_ACCESSKEYSPERUSERQUOTA', '2': 4},
    {'1': 'SUMMARY_KEY_TYPE_POLICIES', '2': 5},
    {'1': 'SUMMARY_KEY_TYPE_INSTANCEPROFILES', '2': 6},
    {'1': 'SUMMARY_KEY_TYPE_ROLEPOLICYSIZEQUOTA', '2': 7},
    {'1': 'SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERROLEQUOTA', '2': 8},
    {'1': 'SUMMARY_KEY_TYPE_ASSUMEROLEPOLICYSIZEQUOTA', '2': 9},
    {'1': 'SUMMARY_KEY_TYPE_SIGNINGCERTIFICATESPERUSERQUOTA', '2': 10},
    {'1': 'SUMMARY_KEY_TYPE_SERVERCERTIFICATESQUOTA', '2': 11},
    {'1': 'SUMMARY_KEY_TYPE_INSTANCEPROFILESQUOTA', '2': 12},
    {'1': 'SUMMARY_KEY_TYPE_ROLESQUOTA', '2': 13},
    {'1': 'SUMMARY_KEY_TYPE_GROUPPOLICYSIZEQUOTA', '2': 14},
    {'1': 'SUMMARY_KEY_TYPE_POLICYVERSIONSINUSE', '2': 15},
    {'1': 'SUMMARY_KEY_TYPE_POLICIESQUOTA', '2': 16},
    {'1': 'SUMMARY_KEY_TYPE_USERSQUOTA', '2': 17},
    {'1': 'SUMMARY_KEY_TYPE_USERPOLICYSIZEQUOTA', '2': 18},
    {'1': 'SUMMARY_KEY_TYPE_GLOBALENDPOINTTOKENVERSION', '2': 19},
    {'1': 'SUMMARY_KEY_TYPE_POLICYSIZEQUOTA', '2': 20},
    {'1': 'SUMMARY_KEY_TYPE_SERVERCERTIFICATES', '2': 21},
    {'1': 'SUMMARY_KEY_TYPE_ACCOUNTPASSWORDPRESENT', '2': 22},
    {'1': 'SUMMARY_KEY_TYPE_ACCOUNTSIGNINGCERTIFICATESPRESENT', '2': 23},
    {'1': 'SUMMARY_KEY_TYPE_USERS', '2': 24},
    {'1': 'SUMMARY_KEY_TYPE_POLICYVERSIONSINUSEQUOTA', '2': 25},
    {'1': 'SUMMARY_KEY_TYPE_MFADEVICESINUSE', '2': 26},
    {'1': 'SUMMARY_KEY_TYPE_GROUPS', '2': 27},
    {'1': 'SUMMARY_KEY_TYPE_ACCOUNTACCESSKEYSPRESENT', '2': 28},
    {'1': 'SUMMARY_KEY_TYPE_GROUPSQUOTA', '2': 29},
    {'1': 'SUMMARY_KEY_TYPE_PROVIDERS', '2': 30},
    {'1': 'SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERUSERQUOTA', '2': 31},
    {'1': 'SUMMARY_KEY_TYPE_GROUPSPERUSERQUOTA', '2': 32},
    {'1': 'SUMMARY_KEY_TYPE_ACCOUNTMFAENABLED', '2': 33},
  ],
};

/// Descriptor for `summaryKeyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List summaryKeyTypeDescriptor = $convert.base64Decode(
    'Cg5zdW1tYXJ5S2V5VHlwZRIyCi5TVU1NQVJZX0tFWV9UWVBFX0FUVEFDSEVEUE9MSUNJRVNQRV'
    'JHUk9VUFFVT1RBEAASGgoWU1VNTUFSWV9LRVlfVFlQRV9ST0xFUxABEisKJ1NVTU1BUllfS0VZ'
    'X1RZUEVfVkVSU0lPTlNQRVJQT0xJQ1lRVU9UQRACEh8KG1NVTU1BUllfS0VZX1RZUEVfTUZBRE'
    'VWSUNFUxADEisKJ1NVTU1BUllfS0VZX1RZUEVfQUNDRVNTS0VZU1BFUlVTRVJRVU9UQRAEEh0K'
    'GVNVTU1BUllfS0VZX1RZUEVfUE9MSUNJRVMQBRIlCiFTVU1NQVJZX0tFWV9UWVBFX0lOU1RBTk'
    'NFUFJPRklMRVMQBhIoCiRTVU1NQVJZX0tFWV9UWVBFX1JPTEVQT0xJQ1lTSVpFUVVPVEEQBxIx'
    'Ci1TVU1NQVJZX0tFWV9UWVBFX0FUVEFDSEVEUE9MSUNJRVNQRVJST0xFUVVPVEEQCBIuCipTVU'
    '1NQVJZX0tFWV9UWVBFX0FTU1VNRVJPTEVQT0xJQ1lTSVpFUVVPVEEQCRI0CjBTVU1NQVJZX0tF'
    'WV9UWVBFX1NJR05JTkdDRVJUSUZJQ0FURVNQRVJVU0VSUVVPVEEQChIsCihTVU1NQVJZX0tFWV'
    '9UWVBFX1NFUlZFUkNFUlRJRklDQVRFU1FVT1RBEAsSKgomU1VNTUFSWV9LRVlfVFlQRV9JTlNU'
    'QU5DRVBST0ZJTEVTUVVPVEEQDBIfChtTVU1NQVJZX0tFWV9UWVBFX1JPTEVTUVVPVEEQDRIpCi'
    'VTVU1NQVJZX0tFWV9UWVBFX0dST1VQUE9MSUNZU0laRVFVT1RBEA4SKAokU1VNTUFSWV9LRVlf'
    'VFlQRV9QT0xJQ1lWRVJTSU9OU0lOVVNFEA8SIgoeU1VNTUFSWV9LRVlfVFlQRV9QT0xJQ0lFU1'
    'FVT1RBEBASHwobU1VNTUFSWV9LRVlfVFlQRV9VU0VSU1FVT1RBEBESKAokU1VNTUFSWV9LRVlf'
    'VFlQRV9VU0VSUE9MSUNZU0laRVFVT1RBEBISLworU1VNTUFSWV9LRVlfVFlQRV9HTE9CQUxFTk'
    'RQT0lOVFRPS0VOVkVSU0lPThATEiQKIFNVTU1BUllfS0VZX1RZUEVfUE9MSUNZU0laRVFVT1RB'
    'EBQSJwojU1VNTUFSWV9LRVlfVFlQRV9TRVJWRVJDRVJUSUZJQ0FURVMQFRIrCidTVU1NQVJZX0'
    'tFWV9UWVBFX0FDQ09VTlRQQVNTV09SRFBSRVNFTlQQFhI2CjJTVU1NQVJZX0tFWV9UWVBFX0FD'
    'Q09VTlRTSUdOSU5HQ0VSVElGSUNBVEVTUFJFU0VOVBAXEhoKFlNVTU1BUllfS0VZX1RZUEVfVV'
    'NFUlMQGBItCilTVU1NQVJZX0tFWV9UWVBFX1BPTElDWVZFUlNJT05TSU5VU0VRVU9UQRAZEiQK'
    'IFNVTU1BUllfS0VZX1RZUEVfTUZBREVWSUNFU0lOVVNFEBoSGwoXU1VNTUFSWV9LRVlfVFlQRV'
    '9HUk9VUFMQGxItCilTVU1NQVJZX0tFWV9UWVBFX0FDQ09VTlRBQ0NFU1NLRVlTUFJFU0VOVBAc'
    'EiAKHFNVTU1BUllfS0VZX1RZUEVfR1JPVVBTUVVPVEEQHRIeChpTVU1NQVJZX0tFWV9UWVBFX1'
    'BST1ZJREVSUxAeEjEKLVNVTU1BUllfS0VZX1RZUEVfQVRUQUNIRURQT0xJQ0lFU1BFUlVTRVJR'
    'VU9UQRAfEicKI1NVTU1BUllfS0VZX1RZUEVfR1JPVVBTUEVSVVNFUlFVT1RBECASJgoiU1VNTU'
    'FSWV9LRVlfVFlQRV9BQ0NPVU5UTUZBRU5BQkxFRBAh');

@$core.Deprecated('Use summaryStateTypeDescriptor instead')
const summaryStateType$json = {
  '1': 'summaryStateType',
  '2': [
    {'1': 'SUMMARY_STATE_TYPE_AVAILABLE', '2': 0},
    {'1': 'SUMMARY_STATE_TYPE_NOT_AVAILABLE', '2': 1},
    {'1': 'SUMMARY_STATE_TYPE_NOT_SUPPORTED', '2': 2},
    {'1': 'SUMMARY_STATE_TYPE_FAILED', '2': 3},
  ],
};

/// Descriptor for `summaryStateType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List summaryStateTypeDescriptor = $convert.base64Decode(
    'ChBzdW1tYXJ5U3RhdGVUeXBlEiAKHFNVTU1BUllfU1RBVEVfVFlQRV9BVkFJTEFCTEUQABIkCi'
    'BTVU1NQVJZX1NUQVRFX1RZUEVfTk9UX0FWQUlMQUJMRRABEiQKIFNVTU1BUllfU1RBVEVfVFlQ'
    'RV9OT1RfU1VQUE9SVEVEEAISHQoZU1VNTUFSWV9TVEFURV9UWVBFX0ZBSUxFRBAD');

@$core.Deprecated('Use acceptDelegationRequestRequestDescriptor instead')
const AcceptDelegationRequestRequest$json = {
  '1': 'AcceptDelegationRequestRequest',
  '2': [
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
  ],
};

/// Descriptor for `AcceptDelegationRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List acceptDelegationRequestRequestDescriptor =
    $convert.base64Decode(
        'Ch5BY2NlcHREZWxlZ2F0aW9uUmVxdWVzdFJlcXVlc3QSNAoTZGVsZWdhdGlvbnJlcXVlc3RpZB'
        'iet4CFASABKAlSE2RlbGVnYXRpb25yZXF1ZXN0aWQ=');

@$core.Deprecated('Use accessDetailDescriptor instead')
const AccessDetail$json = {
  '1': 'AccessDetail',
  '2': [
    {'1': 'entitypath', '3': 93829450, '4': 1, '5': 9, '10': 'entitypath'},
    {
      '1': 'lastauthenticatedtime',
      '3': 258126198,
      '4': 1,
      '5': 9,
      '10': 'lastauthenticatedtime'
    },
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
    {
      '1': 'servicenamespace',
      '3': 432654764,
      '4': 1,
      '5': 9,
      '10': 'servicenamespace'
    },
    {
      '1': 'totalauthenticatedentities',
      '3': 289355788,
      '4': 1,
      '5': 5,
      '10': 'totalauthenticatedentities'
    },
  ],
};

/// Descriptor for `AccessDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDetailDescriptor = $convert.base64Decode(
    'CgxBY2Nlc3NEZXRhaWwSIQoKZW50aXR5cGF0aBjK8t4sIAEoCVIKZW50aXR5cGF0aBI3ChVsYX'
    'N0YXV0aGVudGljYXRlZHRpbWUY9uKKeyABKAlSFWxhc3RhdXRoZW50aWNhdGVkdGltZRIZCgZy'
    'ZWdpb24YnvG5SSABKAlSBnJlZ2lvbhIjCgtzZXJ2aWNlbmFtZRjgqttBIAEoCVILc2VydmljZW'
    '5hbWUSLgoQc2VydmljZW5hbWVzcGFjZRisk6fOASABKAlSEHNlcnZpY2VuYW1lc3BhY2USQgoa'
    'dG90YWxhdXRoZW50aWNhdGVkZW50aXRpZXMYjPD8iQEgASgFUhp0b3RhbGF1dGhlbnRpY2F0ZW'
    'RlbnRpdGllcw==');

@$core.Deprecated('Use accessKeyDescriptor instead')
const AccessKey$json = {
  '1': 'AccessKey',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'secretaccesskey',
      '3': 172288487,
      '4': 1,
      '5': 9,
      '10': 'secretaccesskey'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AccessKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessKeyDescriptor = $convert.base64Decode(
    'CglBY2Nlc3NLZXkSJAoLYWNjZXNza2V5aWQYoLe32AEgASgJUgthY2Nlc3NrZXlpZBIhCgpjcm'
    'VhdGVkYXRlGJK5/BEgASgJUgpjcmVhdGVkYXRlEisKD3NlY3JldGFjY2Vzc2tleRjn05NSIAEo'
    'CVIPc2VjcmV0YWNjZXNza2V5EioKBnN0YXR1cxiQ5PsCIAEoDjIPLmlhbS5zdGF0dXNUeXBlUg'
    'ZzdGF0dXMSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use accessKeyLastUsedDescriptor instead')
const AccessKeyLastUsed$json = {
  '1': 'AccessKeyLastUsed',
  '2': [
    {'1': 'lastuseddate', '3': 70377689, '4': 1, '5': 9, '10': 'lastuseddate'},
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
  ],
};

/// Descriptor for `AccessKeyLastUsed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessKeyLastUsedDescriptor = $convert.base64Decode(
    'ChFBY2Nlc3NLZXlMYXN0VXNlZBIlCgxsYXN0dXNlZGRhdGUY2cHHISABKAlSDGxhc3R1c2VkZG'
    'F0ZRIZCgZyZWdpb24YnvG5SSABKAlSBnJlZ2lvbhIjCgtzZXJ2aWNlbmFtZRjgqttBIAEoCVIL'
    'c2VydmljZW5hbWU=');

@$core.Deprecated('Use accessKeyMetadataDescriptor instead')
const AccessKeyMetadata$json = {
  '1': 'AccessKeyMetadata',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AccessKeyMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessKeyMetadataDescriptor = $convert.base64Decode(
    'ChFBY2Nlc3NLZXlNZXRhZGF0YRIkCgthY2Nlc3NrZXlpZBigt7fYASABKAlSC2FjY2Vzc2tleW'
    'lkEiEKCmNyZWF0ZWRhdGUYkrn8ESABKAlSCmNyZWF0ZWRhdGUSKgoGc3RhdHVzGJDk+wIgASgO'
    'Mg8uaWFtLnN0YXR1c1R5cGVSBnN0YXR1cxIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW'
    '1l');

@$core.Deprecated(
    'Use accountNotManagementOrDelegatedAdministratorExceptionDescriptor instead')
const AccountNotManagementOrDelegatedAdministratorException$json = {
  '1': 'AccountNotManagementOrDelegatedAdministratorException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccountNotManagementOrDelegatedAdministratorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    accountNotManagementOrDelegatedAdministratorExceptionDescriptor =
    $convert.base64Decode(
        'CjVBY2NvdW50Tm90TWFuYWdlbWVudE9yRGVsZWdhdGVkQWRtaW5pc3RyYXRvckV4Y2VwdGlvbh'
        'IbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use addClientIDToOpenIDConnectProviderRequestDescriptor instead')
const AddClientIDToOpenIDConnectProviderRequest$json = {
  '1': 'AddClientIDToOpenIDConnectProviderRequest',
  '2': [
    {'1': 'clientid', '3': 448889284, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
  ],
};

/// Descriptor for `AddClientIDToOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    addClientIDToOpenIDConnectProviderRequestDescriptor = $convert.base64Decode(
        'CilBZGRDbGllbnRJRFRvT3BlbklEQ29ubmVjdFByb3ZpZGVyUmVxdWVzdBIeCghjbGllbnRpZB'
        'jEg4bWASABKAlSCGNsaWVudGlkEj4KGG9wZW5pZGNvbm5lY3Rwcm92aWRlcmFybhjj84/pASAB'
        'KAlSGG9wZW5pZGNvbm5lY3Rwcm92aWRlcmFybg==');

@$core.Deprecated('Use addRoleToInstanceProfileRequestDescriptor instead')
const AddRoleToInstanceProfileRequest$json = {
  '1': 'AddRoleToInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `AddRoleToInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addRoleToInstanceProfileRequestDescriptor =
    $convert.base64Decode(
        'Ch9BZGRSb2xlVG9JbnN0YW5jZVByb2ZpbGVSZXF1ZXN0EjQKE2luc3RhbmNlcHJvZmlsZW5hbW'
        'UY7c682gEgASgJUhNpbnN0YW5jZXByb2ZpbGVuYW1lEh4KCHJvbGVuYW1lGLPzvMIBIAEoCVII'
        'cm9sZW5hbWU=');

@$core.Deprecated('Use addUserToGroupRequestDescriptor instead')
const AddUserToGroupRequest$json = {
  '1': 'AddUserToGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AddUserToGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addUserToGroupRequestDescriptor = $convert.base64Decode(
    'ChVBZGRVc2VyVG9Hcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW'
    '1lEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use associateDelegationRequestRequestDescriptor instead')
const AssociateDelegationRequestRequest$json = {
  '1': 'AssociateDelegationRequestRequest',
  '2': [
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
  ],
};

/// Descriptor for `AssociateDelegationRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateDelegationRequestRequestDescriptor =
    $convert.base64Decode(
        'CiFBc3NvY2lhdGVEZWxlZ2F0aW9uUmVxdWVzdFJlcXVlc3QSNAoTZGVsZWdhdGlvbnJlcXVlc3'
        'RpZBiet4CFASABKAlSE2RlbGVnYXRpb25yZXF1ZXN0aWQ=');

@$core.Deprecated('Use attachGroupPolicyRequestDescriptor instead')
const AttachGroupPolicyRequest$json = {
  '1': 'AttachGroupPolicyRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
  ],
};

/// Descriptor for `AttachGroupPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachGroupPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhBdHRhY2hHcm91cFBvbGljeVJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdX'
        'BuYW1lEiAKCXBvbGljeWFybhjD0rrtASABKAlSCXBvbGljeWFybg==');

@$core.Deprecated('Use attachRolePolicyRequestDescriptor instead')
const AttachRolePolicyRequest$json = {
  '1': 'AttachRolePolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `AttachRolePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachRolePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdBdHRhY2hSb2xlUG9saWN5UmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xpY3'
        'lhcm4SHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use attachUserPolicyRequestDescriptor instead')
const AttachUserPolicyRequest$json = {
  '1': 'AttachUserPolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AttachUserPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachUserPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdBdHRhY2hVc2VyUG9saWN5UmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xpY3'
        'lhcm4SHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use attachedPermissionsBoundaryDescriptor instead')
const AttachedPermissionsBoundary$json = {
  '1': 'AttachedPermissionsBoundary',
  '2': [
    {
      '1': 'permissionsboundaryarn',
      '3': 488828173,
      '4': 1,
      '5': 9,
      '10': 'permissionsboundaryarn'
    },
    {
      '1': 'permissionsboundarytype',
      '3': 523357118,
      '4': 1,
      '5': 14,
      '6': '.iam.PermissionsBoundaryAttachmentType',
      '10': 'permissionsboundarytype'
    },
  ],
};

/// Descriptor for `AttachedPermissionsBoundary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachedPermissionsBoundaryDescriptor = $convert.base64Decode(
    'ChtBdHRhY2hlZFBlcm1pc3Npb25zQm91bmRhcnkSOgoWcGVybWlzc2lvbnNib3VuZGFyeWFybh'
    'iN2ovpASABKAlSFnBlcm1pc3Npb25zYm91bmRhcnlhcm4SZAoXcGVybWlzc2lvbnNib3VuZGFy'
    'eXR5cGUYvpfH+QEgASgOMiYuaWFtLlBlcm1pc3Npb25zQm91bmRhcnlBdHRhY2htZW50VHlwZV'
    'IXcGVybWlzc2lvbnNib3VuZGFyeXR5cGU=');

@$core.Deprecated('Use attachedPolicyDescriptor instead')
const AttachedPolicy$json = {
  '1': 'AttachedPolicy',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `AttachedPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachedPolicyDescriptor = $convert.base64Decode(
    'Cg5BdHRhY2hlZFBvbGljeRIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xpY3lhcm4SIQoKcG'
    '9saWN5bmFtZRi99Yd/IAEoCVIKcG9saWN5bmFtZQ==');

@$core.Deprecated('Use callerIsNotManagementAccountExceptionDescriptor instead')
const CallerIsNotManagementAccountException$json = {
  '1': 'CallerIsNotManagementAccountException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CallerIsNotManagementAccountException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callerIsNotManagementAccountExceptionDescriptor =
    $convert.base64Decode(
        'CiVDYWxsZXJJc05vdE1hbmFnZW1lbnRBY2NvdW50RXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cC'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use changePasswordRequestDescriptor instead')
const ChangePasswordRequest$json = {
  '1': 'ChangePasswordRequest',
  '2': [
    {'1': 'newpassword', '3': 425247981, '4': 1, '5': 9, '10': 'newpassword'},
    {'1': 'oldpassword', '3': 410224204, '4': 1, '5': 9, '10': 'oldpassword'},
  ],
};

/// Descriptor for `ChangePasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changePasswordRequestDescriptor = $convert.base64Decode(
    'ChVDaGFuZ2VQYXNzd29yZFJlcXVlc3QSJAoLbmV3cGFzc3dvcmQY7YnjygEgASgJUgtuZXdwYX'
    'Nzd29yZBIkCgtvbGRwYXNzd29yZBjMjM7DASABKAlSC29sZHBhc3N3b3Jk');

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

@$core.Deprecated('Use contextEntryDescriptor instead')
const ContextEntry$json = {
  '1': 'ContextEntry',
  '2': [
    {
      '1': 'contextkeyname',
      '3': 259692145,
      '4': 1,
      '5': 9,
      '10': 'contextkeyname'
    },
    {
      '1': 'contextkeytype',
      '3': 4954324,
      '4': 1,
      '5': 14,
      '6': '.iam.ContextKeyTypeEnum',
      '10': 'contextkeytype'
    },
    {
      '1': 'contextkeyvalues',
      '3': 410115170,
      '4': 3,
      '5': 9,
      '10': 'contextkeyvalues'
    },
  ],
};

/// Descriptor for `ContextEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextEntryDescriptor = $convert.base64Decode(
    'CgxDb250ZXh0RW50cnkSKQoOY29udGV4dGtleW5hbWUY8azqeyABKAlSDmNvbnRleHRrZXluYW'
    '1lEkIKDmNvbnRleHRrZXl0eXBlGNSxrgIgASgOMhcuaWFtLkNvbnRleHRLZXlUeXBlRW51bVIO'
    'Y29udGV4dGtleXR5cGUSLgoQY29udGV4dGtleXZhbHVlcxjiuMfDASADKAlSEGNvbnRleHRrZX'
    'l2YWx1ZXM=');

@$core.Deprecated('Use createAccessKeyRequestDescriptor instead')
const CreateAccessKeyRequest$json = {
  '1': 'CreateAccessKeyRequest',
  '2': [
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `CreateAccessKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAccessKeyRequestDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVBY2Nlc3NLZXlSZXF1ZXN0Eh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbW'
        'U=');

@$core.Deprecated('Use createAccessKeyResponseDescriptor instead')
const CreateAccessKeyResponse$json = {
  '1': 'CreateAccessKeyResponse',
  '2': [
    {
      '1': 'accesskey',
      '3': 268627635,
      '4': 1,
      '5': 11,
      '6': '.iam.AccessKey',
      '10': 'accesskey'
    },
  ],
};

/// Descriptor for `CreateAccessKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAccessKeyResponseDescriptor =
    $convert.base64Decode(
        'ChdDcmVhdGVBY2Nlc3NLZXlSZXNwb25zZRIwCglhY2Nlc3NrZXkYs92LgAEgASgLMg4uaWFtLk'
        'FjY2Vzc0tleVIJYWNjZXNza2V5');

@$core.Deprecated('Use createAccountAliasRequestDescriptor instead')
const CreateAccountAliasRequest$json = {
  '1': 'CreateAccountAliasRequest',
  '2': [
    {'1': 'accountalias', '3': 446269289, '4': 1, '5': 9, '10': 'accountalias'},
  ],
};

/// Descriptor for `CreateAccountAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAccountAliasRequestDescriptor =
    $convert.base64Decode(
        'ChlDcmVhdGVBY2NvdW50QWxpYXNSZXF1ZXN0EiYKDGFjY291bnRhbGlhcxjpjubUASABKAlSDG'
        'FjY291bnRhbGlhcw==');

@$core.Deprecated('Use createDelegationRequestRequestDescriptor instead')
const CreateDelegationRequestRequest$json = {
  '1': 'CreateDelegationRequestRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'notificationchannel',
      '3': 377061570,
      '4': 1,
      '5': 9,
      '10': 'notificationchannel'
    },
    {
      '1': 'onlysendbyowner',
      '3': 218835688,
      '4': 1,
      '5': 8,
      '10': 'onlysendbyowner'
    },
    {
      '1': 'owneraccountid',
      '3': 369721751,
      '4': 1,
      '5': 9,
      '10': 'owneraccountid'
    },
    {
      '1': 'permissions',
      '3': 117505412,
      '4': 1,
      '5': 11,
      '6': '.iam.DelegationPermission',
      '10': 'permissions'
    },
    {'1': 'redirecturl', '3': 50921923, '4': 1, '5': 9, '10': 'redirecturl'},
    {
      '1': 'requestmessage',
      '3': 58996154,
      '4': 1,
      '5': 9,
      '10': 'requestmessage'
    },
    {
      '1': 'requestorworkflowid',
      '3': 505587788,
      '4': 1,
      '5': 9,
      '10': 'requestorworkflowid'
    },
    {
      '1': 'sessionduration',
      '3': 413635428,
      '4': 1,
      '5': 5,
      '10': 'sessionduration'
    },
  ],
};

/// Descriptor for `CreateDelegationRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDelegationRequestRequestDescriptor = $convert.base64Decode(
    'Ch5DcmVhdGVEZWxlZ2F0aW9uUmVxdWVzdFJlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKA'
    'lSC2Rlc2NyaXB0aW9uEjQKE25vdGlmaWNhdGlvbmNoYW5uZWwYwoHmswEgASgJUhNub3RpZmlj'
    'YXRpb25jaGFubmVsEisKD29ubHlzZW5kYnlvd25lchjo1axoIAEoCFIPb25seXNlbmRieW93bm'
    'VyEioKDm93bmVyYWNjb3VudGlkGJeDprABIAEoCVIOb3duZXJhY2NvdW50aWQSPgoLcGVybWlz'
    'c2lvbnMYhPuDOCABKAsyGS5pYW0uRGVsZWdhdGlvblBlcm1pc3Npb25SC3Blcm1pc3Npb25zEi'
    'MKC3JlZGlyZWN0dXJsGMODpBggASgJUgtyZWRpcmVjdHVybBIpCg5yZXF1ZXN0bWVzc2FnZRi6'
    '65AcIAEoCVIOcmVxdWVzdG1lc3NhZ2USNAoTcmVxdWVzdG9yd29ya2Zsb3dpZBjM0IrxASABKA'
    'lSE3JlcXVlc3RvcndvcmtmbG93aWQSLAoPc2Vzc2lvbmR1cmF0aW9uGOSmnsUBIAEoBVIPc2Vz'
    'c2lvbmR1cmF0aW9u');

@$core.Deprecated('Use createDelegationRequestResponseDescriptor instead')
const CreateDelegationRequestResponse$json = {
  '1': 'CreateDelegationRequestResponse',
  '2': [
    {
      '1': 'consoledeeplink',
      '3': 60202521,
      '4': 1,
      '5': 9,
      '10': 'consoledeeplink'
    },
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
  ],
};

/// Descriptor for `CreateDelegationRequestResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDelegationRequestResponseDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVEZWxlZ2F0aW9uUmVxdWVzdFJlc3BvbnNlEisKD2NvbnNvbGVkZWVwbGluaxiZvN'
        'ocIAEoCVIPY29uc29sZWRlZXBsaW5rEjQKE2RlbGVnYXRpb25yZXF1ZXN0aWQYnreAhQEgASgJ'
        'UhNkZWxlZ2F0aW9ucmVxdWVzdGlk');

@$core.Deprecated('Use createGroupRequestDescriptor instead')
const CreateGroupRequest$json = {
  '1': 'CreateGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `CreateGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGroupRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVHcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEh'
    'UKBHBhdGgY18ibWyABKAlSBHBhdGg=');

@$core.Deprecated('Use createGroupResponseDescriptor instead')
const CreateGroupResponse$json = {
  '1': 'CreateGroupResponse',
  '2': [
    {
      '1': 'group',
      '3': 91525165,
      '4': 1,
      '5': 11,
      '6': '.iam.Group',
      '10': 'group'
    },
  ],
};

/// Descriptor for `CreateGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGroupResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVHcm91cFJlc3BvbnNlEiMKBWdyb3VwGK2g0isgASgLMgouaWFtLkdyb3VwUgVncm'
    '91cA==');

@$core.Deprecated('Use createInstanceProfileRequestDescriptor instead')
const CreateInstanceProfileRequest$json = {
  '1': 'CreateInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createInstanceProfileRequestDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVJbnN0YW5jZVByb2ZpbGVSZXF1ZXN0EjQKE2luc3RhbmNlcHJvZmlsZW5hbWUY7c'
        '682gEgASgJUhNpbnN0YW5jZXByb2ZpbGVuYW1lEhUKBHBhdGgY18ibWyABKAlSBHBhdGgSIAoE'
        'dGFncxjBwfa1ASADKAsyCC5pYW0uVGFnUgR0YWdz');

@$core.Deprecated('Use createInstanceProfileResponseDescriptor instead')
const CreateInstanceProfileResponse$json = {
  '1': 'CreateInstanceProfileResponse',
  '2': [
    {
      '1': 'instanceprofile',
      '3': 57726800,
      '4': 1,
      '5': 11,
      '6': '.iam.InstanceProfile',
      '10': 'instanceprofile'
    },
  ],
};

/// Descriptor for `CreateInstanceProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createInstanceProfileResponseDescriptor =
    $convert.base64Decode(
        'Ch1DcmVhdGVJbnN0YW5jZVByb2ZpbGVSZXNwb25zZRJBCg9pbnN0YW5jZXByb2ZpbGUY0K7DGy'
        'ABKAsyFC5pYW0uSW5zdGFuY2VQcm9maWxlUg9pbnN0YW5jZXByb2ZpbGU=');

@$core.Deprecated('Use createLoginProfileRequestDescriptor instead')
const CreateLoginProfileRequest$json = {
  '1': 'CreateLoginProfileRequest',
  '2': [
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {
      '1': 'passwordresetrequired',
      '3': 330943683,
      '4': 1,
      '5': 8,
      '10': 'passwordresetrequired'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `CreateLoginProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createLoginProfileRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVMb2dpblByb2ZpbGVSZXF1ZXN0Eh0KCHBhc3N3b3JkGLmQjGYgASgJUghwYXNzd2'
    '9yZBI4ChVwYXNzd29yZHJlc2V0cmVxdWlyZWQYw5nnnQEgASgIUhVwYXNzd29yZHJlc2V0cmVx'
    'dWlyZWQSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use createLoginProfileResponseDescriptor instead')
const CreateLoginProfileResponse$json = {
  '1': 'CreateLoginProfileResponse',
  '2': [
    {
      '1': 'loginprofile',
      '3': 465493040,
      '4': 1,
      '5': 11,
      '6': '.iam.LoginProfile',
      '10': 'loginprofile'
    },
  ],
};

/// Descriptor for `CreateLoginProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createLoginProfileResponseDescriptor =
    $convert.base64Decode(
        'ChpDcmVhdGVMb2dpblByb2ZpbGVSZXNwb25zZRI5Cgxsb2dpbnByb2ZpbGUYsLj73QEgASgLMh'
        'EuaWFtLkxvZ2luUHJvZmlsZVIMbG9naW5wcm9maWxl');

@$core.Deprecated('Use createOpenIDConnectProviderRequestDescriptor instead')
const CreateOpenIDConnectProviderRequest$json = {
  '1': 'CreateOpenIDConnectProviderRequest',
  '2': [
    {'1': 'clientidlist', '3': 428113132, '4': 3, '5': 9, '10': 'clientidlist'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {
      '1': 'thumbprintlist',
      '3': 88528351,
      '4': 3,
      '5': 9,
      '10': 'thumbprintlist'
    },
    {'1': 'url', '3': 354018239, '4': 1, '5': 9, '10': 'url'},
  ],
};

/// Descriptor for `CreateOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOpenIDConnectProviderRequestDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVPcGVuSURDb25uZWN0UHJvdmlkZXJSZXF1ZXN0EiYKDGNsaWVudGlkbGlzdBjs+Z'
        'HMASADKAlSDGNsaWVudGlkbGlzdBIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3MS'
        'KQoOdGh1bWJwcmludGxpc3QY36ubKiADKAlSDnRodW1icHJpbnRsaXN0EhQKA3VybBi/x+eoAS'
        'ABKAlSA3VybA==');

@$core.Deprecated('Use createOpenIDConnectProviderResponseDescriptor instead')
const CreateOpenIDConnectProviderResponse$json = {
  '1': 'CreateOpenIDConnectProviderResponse',
  '2': [
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateOpenIDConnectProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOpenIDConnectProviderResponseDescriptor =
    $convert.base64Decode(
        'CiNDcmVhdGVPcGVuSURDb25uZWN0UHJvdmlkZXJSZXNwb25zZRI+ChhvcGVuaWRjb25uZWN0cH'
        'JvdmlkZXJhcm4Y4/OP6QEgASgJUhhvcGVuaWRjb25uZWN0cHJvdmlkZXJhcm4SIAoEdGFncxjB'
        'wfa1ASADKAsyCC5pYW0uVGFnUgR0YWdz');

@$core.Deprecated('Use createPolicyRequestDescriptor instead')
const CreatePolicyRequest$json = {
  '1': 'CreatePolicyRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreatePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPolicyRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVQb2xpY3lSZXF1ZXN0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdG'
    'lvbhIVCgRwYXRoGNfIm1sgASgJUgRwYXRoEikKDnBvbGljeWRvY3VtZW50GMuuwXEgASgJUg5w'
    'b2xpY3lkb2N1bWVudBIhCgpwb2xpY3luYW1lGL31h38gASgJUgpwb2xpY3luYW1lEiAKBHRhZ3'
    'MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createPolicyResponseDescriptor instead')
const CreatePolicyResponse$json = {
  '1': 'CreatePolicyResponse',
  '2': [
    {
      '1': 'policy',
      '3': 471611296,
      '4': 1,
      '5': 11,
      '6': '.iam.Policy',
      '10': 'policy'
    },
  ],
};

/// Descriptor for `CreatePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPolicyResponseDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVQb2xpY3lSZXNwb25zZRInCgZwb2xpY3kYoO/w4AEgASgLMgsuaWFtLlBvbGljeV'
    'IGcG9saWN5');

@$core.Deprecated('Use createPolicyVersionRequestDescriptor instead')
const CreatePolicyVersionRequest$json = {
  '1': 'CreatePolicyVersionRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'setasdefault', '3': 428021657, '4': 1, '5': 8, '10': 'setasdefault'},
  ],
};

/// Descriptor for `CreatePolicyVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPolicyVersionRequestDescriptor =
    $convert.base64Decode(
        'ChpDcmVhdGVQb2xpY3lWZXJzaW9uUmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2'
        'xpY3lhcm4SKQoOcG9saWN5ZG9jdW1lbnQYy67BcSABKAlSDnBvbGljeWRvY3VtZW50EiYKDHNl'
        'dGFzZGVmYXVsdBiZr4zMASABKAhSDHNldGFzZGVmYXVsdA==');

@$core.Deprecated('Use createPolicyVersionResponseDescriptor instead')
const CreatePolicyVersionResponse$json = {
  '1': 'CreatePolicyVersionResponse',
  '2': [
    {
      '1': 'policyversion',
      '3': 266935938,
      '4': 1,
      '5': 11,
      '6': '.iam.PolicyVersion',
      '10': 'policyversion'
    },
  ],
};

/// Descriptor for `CreatePolicyVersionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPolicyVersionResponseDescriptor =
    $convert.base64Decode(
        'ChtDcmVhdGVQb2xpY3lWZXJzaW9uUmVzcG9uc2USOwoNcG9saWN5dmVyc2lvbhiCvaR/IAEoCz'
        'ISLmlhbS5Qb2xpY3lWZXJzaW9uUg1wb2xpY3l2ZXJzaW9u');

@$core.Deprecated('Use createRoleRequestDescriptor instead')
const CreateRoleRequest$json = {
  '1': 'CreateRoleRequest',
  '2': [
    {
      '1': 'assumerolepolicydocument',
      '3': 500384765,
      '4': 1,
      '5': 9,
      '10': 'assumerolepolicydocument'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'maxsessionduration',
      '3': 391414912,
      '4': 1,
      '5': 5,
      '10': 'maxsessionduration'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 9,
      '10': 'permissionsboundary'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRoleRequestDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVSb2xlUmVxdWVzdBI+Chhhc3N1bWVyb2xlcG9saWN5ZG9jdW1lbnQY/YfN7gEgAS'
    'gJUhhhc3N1bWVyb2xlcG9saWN5ZG9jdW1lbnQSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rl'
    'c2NyaXB0aW9uEjIKEm1heHNlc3Npb25kdXJhdGlvbhiAidK6ASABKAVSEm1heHNlc3Npb25kdX'
    'JhdGlvbhIVCgRwYXRoGNfIm1sgASgJUgRwYXRoEjMKE3Blcm1pc3Npb25zYm91bmRhcnkYtsOg'
    'CCABKAlSE3Blcm1pc3Npb25zYm91bmRhcnkSHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbm'
    'FtZRIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use createRoleResponseDescriptor instead')
const CreateRoleResponse$json = {
  '1': 'CreateRoleResponse',
  '2': [
    {
      '1': 'role',
      '3': 271285818,
      '4': 1,
      '5': 11,
      '6': '.iam.Role',
      '10': 'role'
    },
  ],
};

/// Descriptor for `CreateRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRoleResponseDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVSb2xlUmVzcG9uc2USIQoEcm9sZRi6/K2BASABKAsyCS5pYW0uUm9sZVIEcm9sZQ'
    '==');

@$core.Deprecated('Use createSAMLProviderRequestDescriptor instead')
const CreateSAMLProviderRequest$json = {
  '1': 'CreateSAMLProviderRequest',
  '2': [
    {
      '1': 'addprivatekey',
      '3': 205767607,
      '4': 1,
      '5': 9,
      '10': 'addprivatekey'
    },
    {
      '1': 'assertionencryptionmode',
      '3': 474560298,
      '4': 1,
      '5': 14,
      '6': '.iam.assertionEncryptionModeType',
      '10': 'assertionencryptionmode'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'samlmetadatadocument',
      '3': 282432645,
      '4': 1,
      '5': 9,
      '10': 'samlmetadatadocument'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateSAMLProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSAMLProviderRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVTQU1MUHJvdmlkZXJSZXF1ZXN0EicKDWFkZHByaXZhdGVrZXkYt4ePYiABKAlSDW'
    'FkZHByaXZhdGVrZXkSXgoXYXNzZXJ0aW9uZW5jcnlwdGlvbm1vZGUYqu6k4gEgASgOMiAuaWFt'
    'LmFzc2VydGlvbkVuY3J5cHRpb25Nb2RlVHlwZVIXYXNzZXJ0aW9uZW5jcnlwdGlvbm1vZGUSFQ'
    'oEbmFtZRiH5oF/IAEoCVIEbmFtZRI2ChRzYW1sbWV0YWRhdGFkb2N1bWVudBiFqdaGASABKAlS'
    'FHNhbWxtZXRhZGF0YWRvY3VtZW50EiAKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncw'
    '==');

@$core.Deprecated('Use createSAMLProviderResponseDescriptor instead')
const CreateSAMLProviderResponse$json = {
  '1': 'CreateSAMLProviderResponse',
  '2': [
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateSAMLProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSAMLProviderResponseDescriptor =
    $convert.base64Decode(
        'ChpDcmVhdGVTQU1MUHJvdmlkZXJSZXNwb25zZRIsCg9zYW1scHJvdmlkZXJhcm4Ypc2ojAEgAS'
        'gJUg9zYW1scHJvdmlkZXJhcm4SIAoEdGFncxjBwfa1ASADKAsyCC5pYW0uVGFnUgR0YWdz');

@$core.Deprecated('Use createServiceLinkedRoleRequestDescriptor instead')
const CreateServiceLinkedRoleRequest$json = {
  '1': 'CreateServiceLinkedRoleRequest',
  '2': [
    {
      '1': 'awsservicename',
      '3': 378200975,
      '4': 1,
      '5': 9,
      '10': 'awsservicename'
    },
    {'1': 'customsuffix', '3': 269317868, '4': 1, '5': 9, '10': 'customsuffix'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `CreateServiceLinkedRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createServiceLinkedRoleRequestDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVTZXJ2aWNlTGlua2VkUm9sZVJlcXVlc3QSKgoOYXdzc2VydmljZW5hbWUYj8ertA'
        'EgASgJUg5hd3NzZXJ2aWNlbmFtZRImCgxjdXN0b21zdWZmaXgY7O21gAEgASgJUgxjdXN0b21z'
        'dWZmaXgSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9u');

@$core.Deprecated('Use createServiceLinkedRoleResponseDescriptor instead')
const CreateServiceLinkedRoleResponse$json = {
  '1': 'CreateServiceLinkedRoleResponse',
  '2': [
    {
      '1': 'role',
      '3': 271285818,
      '4': 1,
      '5': 11,
      '6': '.iam.Role',
      '10': 'role'
    },
  ],
};

/// Descriptor for `CreateServiceLinkedRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createServiceLinkedRoleResponseDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVTZXJ2aWNlTGlua2VkUm9sZVJlc3BvbnNlEiEKBHJvbGUYuvytgQEgASgLMgkuaW'
        'FtLlJvbGVSBHJvbGU=');

@$core
    .Deprecated('Use createServiceSpecificCredentialRequestDescriptor instead')
const CreateServiceSpecificCredentialRequest$json = {
  '1': 'CreateServiceSpecificCredentialRequest',
  '2': [
    {
      '1': 'credentialagedays',
      '3': 119555957,
      '4': 1,
      '5': 5,
      '10': 'credentialagedays'
    },
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `CreateServiceSpecificCredentialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createServiceSpecificCredentialRequestDescriptor =
    $convert.base64Decode(
        'CiZDcmVhdGVTZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFsUmVxdWVzdBIvChFjcmVkZW50aWFsYW'
        'dlZGF5cxj1joE5IAEoBVIRY3JlZGVudGlhbGFnZWRheXMSIwoLc2VydmljZW5hbWUY4KrbQSAB'
        'KAlSC3NlcnZpY2VuYW1lEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core
    .Deprecated('Use createServiceSpecificCredentialResponseDescriptor instead')
const CreateServiceSpecificCredentialResponse$json = {
  '1': 'CreateServiceSpecificCredentialResponse',
  '2': [
    {
      '1': 'servicespecificcredential',
      '3': 186794126,
      '4': 1,
      '5': 11,
      '6': '.iam.ServiceSpecificCredential',
      '10': 'servicespecificcredential'
    },
  ],
};

/// Descriptor for `CreateServiceSpecificCredentialResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createServiceSpecificCredentialResponseDescriptor =
    $convert.base64Decode(
        'CidDcmVhdGVTZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFsUmVzcG9uc2USXwoZc2VydmljZXNwZW'
        'NpZmljY3JlZGVudGlhbBiOgYlZIAEoCzIeLmlhbS5TZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFs'
        'UhlzZXJ2aWNlc3BlY2lmaWNjcmVkZW50aWFs');

@$core.Deprecated('Use createUserRequestDescriptor instead')
const CreateUserRequest$json = {
  '1': 'CreateUserRequest',
  '2': [
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 9,
      '10': 'permissionsboundary'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `CreateUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserRequestDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVVc2VyUmVxdWVzdBIVCgRwYXRoGNfIm1sgASgJUgRwYXRoEjMKE3Blcm1pc3Npb2'
    '5zYm91bmRhcnkYtsOgCCABKAlSE3Blcm1pc3Npb25zYm91bmRhcnkSIAoEdGFncxjBwfa1ASAD'
    'KAsyCC5pYW0uVGFnUgR0YWdzEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use createUserResponseDescriptor instead')
const CreateUserResponse$json = {
  '1': 'CreateUserResponse',
  '2': [
    {
      '1': 'user',
      '3': 10894867,
      '4': 1,
      '5': 11,
      '6': '.iam.User',
      '10': 'user'
    },
  ],
};

/// Descriptor for `CreateUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserResponseDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVVc2VyUmVzcG9uc2USIAoEdXNlchiT/JgFIAEoCzIJLmlhbS5Vc2VyUgR1c2Vy');

@$core.Deprecated('Use createVirtualMFADeviceRequestDescriptor instead')
const CreateVirtualMFADeviceRequest$json = {
  '1': 'CreateVirtualMFADeviceRequest',
  '2': [
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {
      '1': 'virtualmfadevicename',
      '3': 531043720,
      '4': 1,
      '5': 9,
      '10': 'virtualmfadevicename'
    },
  ],
};

/// Descriptor for `CreateVirtualMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVirtualMFADeviceRequestDescriptor =
    $convert.base64Decode(
        'Ch1DcmVhdGVWaXJ0dWFsTUZBRGV2aWNlUmVxdWVzdBIVCgRwYXRoGNfIm1sgASgJUgRwYXRoEi'
        'AKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncxI2ChR2aXJ0dWFsbWZhZGV2aWNlbmFt'
        'ZRiIq5z9ASABKAlSFHZpcnR1YWxtZmFkZXZpY2VuYW1l');

@$core.Deprecated('Use createVirtualMFADeviceResponseDescriptor instead')
const CreateVirtualMFADeviceResponse$json = {
  '1': 'CreateVirtualMFADeviceResponse',
  '2': [
    {
      '1': 'virtualmfadevice',
      '3': 105575533,
      '4': 1,
      '5': 11,
      '6': '.iam.VirtualMFADevice',
      '10': 'virtualmfadevice'
    },
  ],
};

/// Descriptor for `CreateVirtualMFADeviceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVirtualMFADeviceResponseDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVWaXJ0dWFsTUZBRGV2aWNlUmVzcG9uc2USRAoQdmlydHVhbG1mYWRldmljZRjt6K'
        'syIAEoCzIVLmlhbS5WaXJ0dWFsTUZBRGV2aWNlUhB2aXJ0dWFsbWZhZGV2aWNl');

@$core.Deprecated('Use credentialReportExpiredExceptionDescriptor instead')
const CredentialReportExpiredException$json = {
  '1': 'CredentialReportExpiredException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CredentialReportExpiredException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List credentialReportExpiredExceptionDescriptor =
    $convert.base64Decode(
        'CiBDcmVkZW50aWFsUmVwb3J0RXhwaXJlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use credentialReportNotPresentExceptionDescriptor instead')
const CredentialReportNotPresentException$json = {
  '1': 'CredentialReportNotPresentException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CredentialReportNotPresentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List credentialReportNotPresentExceptionDescriptor =
    $convert.base64Decode(
        'CiNDcmVkZW50aWFsUmVwb3J0Tm90UHJlc2VudEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use credentialReportNotReadyExceptionDescriptor instead')
const CredentialReportNotReadyException$json = {
  '1': 'CredentialReportNotReadyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CredentialReportNotReadyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List credentialReportNotReadyExceptionDescriptor =
    $convert.base64Decode(
        'CiFDcmVkZW50aWFsUmVwb3J0Tm90UmVhZHlFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use deactivateMFADeviceRequestDescriptor instead')
const DeactivateMFADeviceRequest$json = {
  '1': 'DeactivateMFADeviceRequest',
  '2': [
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeactivateMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deactivateMFADeviceRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWFjdGl2YXRlTUZBRGV2aWNlUmVxdWVzdBImCgxzZXJpYWxudW1iZXIY5bq5xwEgASgJUg'
        'xzZXJpYWxudW1iZXISHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use delegationPermissionDescriptor instead')
const DelegationPermission$json = {
  '1': 'DelegationPermission',
  '2': [
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyParameter',
      '10': 'parameters'
    },
    {
      '1': 'policytemplatearn',
      '3': 456924085,
      '4': 1,
      '5': 9,
      '10': 'policytemplatearn'
    },
  ],
};

/// Descriptor for `DelegationPermission`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationPermissionDescriptor = $convert.base64Decode(
    'ChREZWxlZ2F0aW9uUGVybWlzc2lvbhI4CgpwYXJhbWV0ZXJzGPqn/usBIAMoCzIULmlhbS5Qb2'
    'xpY3lQYXJhbWV0ZXJSCnBhcmFtZXRlcnMSMAoRcG9saWN5dGVtcGxhdGVhcm4Ytbfw2QEgASgJ'
    'UhFwb2xpY3l0ZW1wbGF0ZWFybg==');

@$core.Deprecated('Use delegationRequestDescriptor instead')
const DelegationRequest$json = {
  '1': 'DelegationRequest',
  '2': [
    {'1': 'approverid', '3': 358748230, '4': 1, '5': 9, '10': 'approverid'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'expirationtime',
      '3': 93473378,
      '4': 1,
      '5': 9,
      '10': 'expirationtime'
    },
    {'1': 'notes', '3': 267393899, '4': 1, '5': 9, '10': 'notes'},
    {
      '1': 'onlysendbyowner',
      '3': 218835688,
      '4': 1,
      '5': 8,
      '10': 'onlysendbyowner'
    },
    {
      '1': 'owneraccountid',
      '3': 369721751,
      '4': 1,
      '5': 9,
      '10': 'owneraccountid'
    },
    {'1': 'ownerid', '3': 375630298, '4': 1, '5': 9, '10': 'ownerid'},
    {
      '1': 'permissionpolicy',
      '3': 482652389,
      '4': 1,
      '5': 9,
      '10': 'permissionpolicy'
    },
    {
      '1': 'permissions',
      '3': 117505412,
      '4': 1,
      '5': 11,
      '6': '.iam.DelegationPermission',
      '10': 'permissions'
    },
    {'1': 'redirecturl', '3': 50921923, '4': 1, '5': 9, '10': 'redirecturl'},
    {
      '1': 'rejectionreason',
      '3': 114778373,
      '4': 1,
      '5': 9,
      '10': 'rejectionreason'
    },
    {
      '1': 'requestmessage',
      '3': 58996154,
      '4': 1,
      '5': 9,
      '10': 'requestmessage'
    },
    {'1': 'requestorid', '3': 231558211, '4': 1, '5': 9, '10': 'requestorid'},
    {
      '1': 'requestorname',
      '3': 16690957,
      '4': 1,
      '5': 9,
      '10': 'requestorname'
    },
    {
      '1': 'rolepermissionrestrictionarns',
      '3': 47253157,
      '4': 3,
      '5': 9,
      '10': 'rolepermissionrestrictionarns'
    },
    {
      '1': 'sessionduration',
      '3': 413635428,
      '4': 1,
      '5': 5,
      '10': 'sessionduration'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.iam.stateType',
      '10': 'state'
    },
    {'1': 'updatedtime', '3': 156274570, '4': 1, '5': 9, '10': 'updatedtime'},
  ],
};

/// Descriptor for `DelegationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationRequestDescriptor = $convert.base64Decode(
    'ChFEZWxlZ2F0aW9uUmVxdWVzdBIiCgphcHByb3ZlcmlkGMagiKsBIAEoCVIKYXBwcm92ZXJpZB'
    'IhCgpjcmVhdGVkYXRlGJK5/BEgASgJUgpjcmVhdGVkYXRlEjQKE2RlbGVnYXRpb25yZXF1ZXN0'
    'aWQYnreAhQEgASgJUhNkZWxlZ2F0aW9ucmVxdWVzdGlkEiMKC2Rlc2NyaXB0aW9uGIr0+TYgAS'
    'gJUgtkZXNjcmlwdGlvbhIpCg5leHBpcmF0aW9udGltZRjilMksIAEoCVIOZXhwaXJhdGlvbnRp'
    'bWUSFwoFbm90ZXMY67bAfyABKAlSBW5vdGVzEisKD29ubHlzZW5kYnlvd25lchjo1axoIAEoCF'
    'IPb25seXNlbmRieW93bmVyEioKDm93bmVyYWNjb3VudGlkGJeDprABIAEoCVIOb3duZXJhY2Nv'
    'dW50aWQSHAoHb3duZXJpZBja046zASABKAlSB293bmVyaWQSLgoQcGVybWlzc2lvbnBvbGljeR'
    'jl4ZLmASABKAlSEHBlcm1pc3Npb25wb2xpY3kSPgoLcGVybWlzc2lvbnMYhPuDOCABKAsyGS5p'
    'YW0uRGVsZWdhdGlvblBlcm1pc3Npb25SC3Blcm1pc3Npb25zEiMKC3JlZGlyZWN0dXJsGMODpB'
    'ggASgJUgtyZWRpcmVjdHVybBIrCg9yZWplY3Rpb25yZWFzb24YhcLdNiABKAlSD3JlamVjdGlv'
    'bnJlYXNvbhIpCg5yZXF1ZXN0bWVzc2FnZRi665AcIAEoCVIOcmVxdWVzdG1lc3NhZ2USIwoLcm'
    'VxdWVzdG9yaWQYw5i1biABKAlSC3JlcXVlc3RvcmlkEicKDXJlcXVlc3Rvcm5hbWUYjd76ByAB'
    'KAlSDXJlcXVlc3Rvcm5hbWUSRwodcm9sZXBlcm1pc3Npb25yZXN0cmljdGlvbmFybnMYpY3EFi'
    'ADKAlSHXJvbGVwZXJtaXNzaW9ucmVzdHJpY3Rpb25hcm5zEiwKD3Nlc3Npb25kdXJhdGlvbhjk'
    'pp7FASABKAVSD3Nlc3Npb25kdXJhdGlvbhIoCgVzdGF0ZRiXybLvASABKA4yDi5pYW0uc3RhdG'
    'VUeXBlUgVzdGF0ZRIjCgt1cGRhdGVkdGltZRiKn8JKIAEoCVILdXBkYXRlZHRpbWU=');

@$core.Deprecated('Use deleteAccessKeyRequestDescriptor instead')
const DeleteAccessKeyRequest$json = {
  '1': 'DeleteAccessKeyRequest',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteAccessKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAccessKeyRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVBY2Nlc3NLZXlSZXF1ZXN0EiQKC2FjY2Vzc2tleWlkGKC3t9gBIAEoCVILYWNjZX'
        'Nza2V5aWQSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use deleteAccountAliasRequestDescriptor instead')
const DeleteAccountAliasRequest$json = {
  '1': 'DeleteAccountAliasRequest',
  '2': [
    {'1': 'accountalias', '3': 446269289, '4': 1, '5': 9, '10': 'accountalias'},
  ],
};

/// Descriptor for `DeleteAccountAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAccountAliasRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVBY2NvdW50QWxpYXNSZXF1ZXN0EiYKDGFjY291bnRhbGlhcxjpjubUASABKAlSDG'
        'FjY291bnRhbGlhcw==');

@$core.Deprecated('Use deleteConflictExceptionDescriptor instead')
const DeleteConflictException$json = {
  '1': 'DeleteConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DeleteConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use deleteGroupPolicyRequestDescriptor instead')
const DeleteGroupPolicyRequest$json = {
  '1': 'DeleteGroupPolicyRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `DeleteGroupPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteGroupPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVHcm91cFBvbGljeVJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdX'
        'BuYW1lEiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWU=');

@$core.Deprecated('Use deleteGroupRequestDescriptor instead')
const DeleteGroupRequest$json = {
  '1': 'DeleteGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
  ],
};

/// Descriptor for `DeleteGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteGroupRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVHcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1l');

@$core.Deprecated('Use deleteInstanceProfileRequestDescriptor instead')
const DeleteInstanceProfileRequest$json = {
  '1': 'DeleteInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
  ],
};

/// Descriptor for `DeleteInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteInstanceProfileRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVJbnN0YW5jZVByb2ZpbGVSZXF1ZXN0EjQKE2luc3RhbmNlcHJvZmlsZW5hbWUY7c'
        '682gEgASgJUhNpbnN0YW5jZXByb2ZpbGVuYW1l');

@$core.Deprecated('Use deleteLoginProfileRequestDescriptor instead')
const DeleteLoginProfileRequest$json = {
  '1': 'DeleteLoginProfileRequest',
  '2': [
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteLoginProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteLoginProfileRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVMb2dpblByb2ZpbGVSZXF1ZXN0Eh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm'
        '5hbWU=');

@$core.Deprecated('Use deleteOpenIDConnectProviderRequestDescriptor instead')
const DeleteOpenIDConnectProviderRequest$json = {
  '1': 'DeleteOpenIDConnectProviderRequest',
  '2': [
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
  ],
};

/// Descriptor for `DeleteOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOpenIDConnectProviderRequestDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVPcGVuSURDb25uZWN0UHJvdmlkZXJSZXF1ZXN0Ej4KGG9wZW5pZGNvbm5lY3Rwcm'
        '92aWRlcmFybhjj84/pASABKAlSGG9wZW5pZGNvbm5lY3Rwcm92aWRlcmFybg==');

@$core.Deprecated('Use deletePolicyRequestDescriptor instead')
const DeletePolicyRequest$json = {
  '1': 'DeletePolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
  ],
};

/// Descriptor for `DeletePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePolicyRequestDescriptor = $convert.base64Decode(
    'ChNEZWxldGVQb2xpY3lSZXF1ZXN0EiAKCXBvbGljeWFybhjD0rrtASABKAlSCXBvbGljeWFybg'
    '==');

@$core.Deprecated('Use deletePolicyVersionRequestDescriptor instead')
const DeletePolicyVersionRequest$json = {
  '1': 'DeletePolicyVersionRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeletePolicyVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePolicyVersionRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVQb2xpY3lWZXJzaW9uUmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2'
        'xpY3lhcm4SIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use deleteRolePermissionsBoundaryRequestDescriptor instead')
const DeleteRolePermissionsBoundaryRequest$json = {
  '1': 'DeleteRolePermissionsBoundaryRequest',
  '2': [
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `DeleteRolePermissionsBoundaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRolePermissionsBoundaryRequestDescriptor =
    $convert.base64Decode(
        'CiREZWxldGVSb2xlUGVybWlzc2lvbnNCb3VuZGFyeVJlcXVlc3QSHgoIcm9sZW5hbWUYs/O8wg'
        'EgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use deleteRolePolicyRequestDescriptor instead')
const DeleteRolePolicyRequest$json = {
  '1': 'DeleteRolePolicyRequest',
  '2': [
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `DeleteRolePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRolePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVSb2xlUG9saWN5UmVxdWVzdBIhCgpwb2xpY3luYW1lGL31h38gASgJUgpwb2xpY3'
        'luYW1lEh4KCHJvbGVuYW1lGLPzvMIBIAEoCVIIcm9sZW5hbWU=');

@$core.Deprecated('Use deleteRoleRequestDescriptor instead')
const DeleteRoleRequest$json = {
  '1': 'DeleteRoleRequest',
  '2': [
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `DeleteRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRoleRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVSb2xlUmVxdWVzdBIeCghyb2xlbmFtZRiz87zCASABKAlSCHJvbGVuYW1l');

@$core.Deprecated('Use deleteSAMLProviderRequestDescriptor instead')
const DeleteSAMLProviderRequest$json = {
  '1': 'DeleteSAMLProviderRequest',
  '2': [
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
  ],
};

/// Descriptor for `DeleteSAMLProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSAMLProviderRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVTQU1MUHJvdmlkZXJSZXF1ZXN0EiwKD3NhbWxwcm92aWRlcmFybhilzaiMASABKA'
        'lSD3NhbWxwcm92aWRlcmFybg==');

@$core.Deprecated('Use deleteSSHPublicKeyRequestDescriptor instead')
const DeleteSSHPublicKeyRequest$json = {
  '1': 'DeleteSSHPublicKeyRequest',
  '2': [
    {
      '1': 'sshpublickeyid',
      '3': 154499415,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeyid'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteSSHPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSSHPublicKeyRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVTU0hQdWJsaWNLZXlSZXF1ZXN0EikKDnNzaHB1YmxpY2tleWlkGNfy1UkgASgJUg'
        '5zc2hwdWJsaWNrZXlpZBIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use deleteServerCertificateRequestDescriptor instead')
const DeleteServerCertificateRequest$json = {
  '1': 'DeleteServerCertificateRequest',
  '2': [
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
  ],
};

/// Descriptor for `DeleteServerCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteServerCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVTZXJ2ZXJDZXJ0aWZpY2F0ZVJlcXVlc3QSNwoVc2VydmVyY2VydGlmaWNhdGVuYW'
        '1lGP+Amx4gASgJUhVzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbWU=');

@$core.Deprecated('Use deleteServiceLinkedRoleRequestDescriptor instead')
const DeleteServiceLinkedRoleRequest$json = {
  '1': 'DeleteServiceLinkedRoleRequest',
  '2': [
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `DeleteServiceLinkedRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteServiceLinkedRoleRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVTZXJ2aWNlTGlua2VkUm9sZVJlcXVlc3QSHgoIcm9sZW5hbWUYs/O8wgEgASgJUg'
        'hyb2xlbmFtZQ==');

@$core.Deprecated('Use deleteServiceLinkedRoleResponseDescriptor instead')
const DeleteServiceLinkedRoleResponse$json = {
  '1': 'DeleteServiceLinkedRoleResponse',
  '2': [
    {
      '1': 'deletiontaskid',
      '3': 379712140,
      '4': 1,
      '5': 9,
      '10': 'deletiontaskid'
    },
  ],
};

/// Descriptor for `DeleteServiceLinkedRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteServiceLinkedRoleResponseDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVTZXJ2aWNlTGlua2VkUm9sZVJlc3BvbnNlEioKDmRlbGV0aW9udGFza2lkGIzlh7'
        'UBIAEoCVIOZGVsZXRpb250YXNraWQ=');

@$core
    .Deprecated('Use deleteServiceSpecificCredentialRequestDescriptor instead')
const DeleteServiceSpecificCredentialRequest$json = {
  '1': 'DeleteServiceSpecificCredentialRequest',
  '2': [
    {
      '1': 'servicespecificcredentialid',
      '3': 426936633,
      '4': 1,
      '5': 9,
      '10': 'servicespecificcredentialid'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteServiceSpecificCredentialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteServiceSpecificCredentialRequestDescriptor =
    $convert.base64Decode(
        'CiZEZWxldGVTZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFsUmVxdWVzdBJEChtzZXJ2aWNlc3BlY2'
        'lmaWNjcmVkZW50aWFsaWQYuZLKywEgASgJUhtzZXJ2aWNlc3BlY2lmaWNjcmVkZW50aWFsaWQS'
        'HgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use deleteSigningCertificateRequestDescriptor instead')
const DeleteSigningCertificateRequest$json = {
  '1': 'DeleteSigningCertificateRequest',
  '2': [
    {
      '1': 'certificateid',
      '3': 149606126,
      '4': 1,
      '5': 9,
      '10': 'certificateid'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteSigningCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSigningCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVTaWduaW5nQ2VydGlmaWNhdGVSZXF1ZXN0EicKDWNlcnRpZmljYXRlaWQY7p2rRy'
        'ABKAlSDWNlcnRpZmljYXRlaWQSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use deleteUserPermissionsBoundaryRequestDescriptor instead')
const DeleteUserPermissionsBoundaryRequest$json = {
  '1': 'DeleteUserPermissionsBoundaryRequest',
  '2': [
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteUserPermissionsBoundaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPermissionsBoundaryRequestDescriptor =
    $convert.base64Decode(
        'CiREZWxldGVVc2VyUGVybWlzc2lvbnNCb3VuZGFyeVJlcXVlc3QSHgoIdXNlcm5hbWUY+sHU4Q'
        'EgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use deleteUserPolicyRequestDescriptor instead')
const DeleteUserPolicyRequest$json = {
  '1': 'DeleteUserPolicyRequest',
  '2': [
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteUserPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVVc2VyUG9saWN5UmVxdWVzdBIhCgpwb2xpY3luYW1lGL31h38gASgJUgpwb2xpY3'
        'luYW1lEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use deleteUserRequestDescriptor instead')
const DeleteUserRequest$json = {
  '1': 'DeleteUserRequest',
  '2': [
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DeleteUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVVc2VyUmVxdWVzdBIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use deleteVirtualMFADeviceRequestDescriptor instead')
const DeleteVirtualMFADeviceRequest$json = {
  '1': 'DeleteVirtualMFADeviceRequest',
  '2': [
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
  ],
};

/// Descriptor for `DeleteVirtualMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVirtualMFADeviceRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVWaXJ0dWFsTUZBRGV2aWNlUmVxdWVzdBImCgxzZXJpYWxudW1iZXIY5bq5xwEgAS'
        'gJUgxzZXJpYWxudW1iZXI=');

@$core.Deprecated('Use deletionTaskFailureReasonTypeDescriptor instead')
const DeletionTaskFailureReasonType$json = {
  '1': 'DeletionTaskFailureReasonType',
  '2': [
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {
      '1': 'roleusagelist',
      '3': 363512789,
      '4': 3,
      '5': 11,
      '6': '.iam.RoleUsageType',
      '10': 'roleusagelist'
    },
  ],
};

/// Descriptor for `DeletionTaskFailureReasonType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletionTaskFailureReasonTypeDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGlvblRhc2tGYWlsdXJlUmVhc29uVHlwZRIZCgZyZWFzb24YuoLFCSABKAlSBnJlYX'
        'NvbhI8Cg1yb2xldXNhZ2VsaXN0GNWHq60BIAMoCzISLmlhbS5Sb2xlVXNhZ2VUeXBlUg1yb2xl'
        'dXNhZ2VsaXN0');

@$core.Deprecated('Use detachGroupPolicyRequestDescriptor instead')
const DetachGroupPolicyRequest$json = {
  '1': 'DetachGroupPolicyRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
  ],
};

/// Descriptor for `DetachGroupPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detachGroupPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhEZXRhY2hHcm91cFBvbGljeVJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdX'
        'BuYW1lEiAKCXBvbGljeWFybhjD0rrtASABKAlSCXBvbGljeWFybg==');

@$core.Deprecated('Use detachRolePolicyRequestDescriptor instead')
const DetachRolePolicyRequest$json = {
  '1': 'DetachRolePolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `DetachRolePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detachRolePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdEZXRhY2hSb2xlUG9saWN5UmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xpY3'
        'lhcm4SHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use detachUserPolicyRequestDescriptor instead')
const DetachUserPolicyRequest$json = {
  '1': 'DetachUserPolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `DetachUserPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detachUserPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdEZXRhY2hVc2VyUG9saWN5UmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xpY3'
        'lhcm4SHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated(
    'Use disableOrganizationsRootCredentialsManagementRequestDescriptor instead')
const DisableOrganizationsRootCredentialsManagementRequest$json = {
  '1': 'DisableOrganizationsRootCredentialsManagementRequest',
};

/// Descriptor for `DisableOrganizationsRootCredentialsManagementRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    disableOrganizationsRootCredentialsManagementRequestDescriptor =
    $convert.base64Decode(
        'CjREaXNhYmxlT3JnYW5pemF0aW9uc1Jvb3RDcmVkZW50aWFsc01hbmFnZW1lbnRSZXF1ZXN0');

@$core.Deprecated(
    'Use disableOrganizationsRootCredentialsManagementResponseDescriptor instead')
const DisableOrganizationsRootCredentialsManagementResponse$json = {
  '1': 'DisableOrganizationsRootCredentialsManagementResponse',
  '2': [
    {
      '1': 'enabledfeatures',
      '3': 211387242,
      '4': 3,
      '5': 14,
      '6': '.iam.FeatureType',
      '10': 'enabledfeatures'
    },
    {
      '1': 'organizationid',
      '3': 311006120,
      '4': 1,
      '5': 9,
      '10': 'organizationid'
    },
  ],
};

/// Descriptor for `DisableOrganizationsRootCredentialsManagementResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    disableOrganizationsRootCredentialsManagementResponseDescriptor =
    $convert.base64Decode(
        'CjVEaXNhYmxlT3JnYW5pemF0aW9uc1Jvb3RDcmVkZW50aWFsc01hbmFnZW1lbnRSZXNwb25zZR'
        'I9Cg9lbmFibGVkZmVhdHVyZXMY6obmZCADKA4yEC5pYW0uRmVhdHVyZVR5cGVSD2VuYWJsZWRm'
        'ZWF0dXJlcxIqCg5vcmdhbml6YXRpb25pZBiop6aUASABKAlSDm9yZ2FuaXphdGlvbmlk');

@$core
    .Deprecated('Use disableOrganizationsRootSessionsRequestDescriptor instead')
const DisableOrganizationsRootSessionsRequest$json = {
  '1': 'DisableOrganizationsRootSessionsRequest',
};

/// Descriptor for `DisableOrganizationsRootSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableOrganizationsRootSessionsRequestDescriptor =
    $convert.base64Decode(
        'CidEaXNhYmxlT3JnYW5pemF0aW9uc1Jvb3RTZXNzaW9uc1JlcXVlc3Q=');

@$core.Deprecated(
    'Use disableOrganizationsRootSessionsResponseDescriptor instead')
const DisableOrganizationsRootSessionsResponse$json = {
  '1': 'DisableOrganizationsRootSessionsResponse',
  '2': [
    {
      '1': 'enabledfeatures',
      '3': 211387242,
      '4': 3,
      '5': 14,
      '6': '.iam.FeatureType',
      '10': 'enabledfeatures'
    },
    {
      '1': 'organizationid',
      '3': 311006120,
      '4': 1,
      '5': 9,
      '10': 'organizationid'
    },
  ],
};

/// Descriptor for `DisableOrganizationsRootSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableOrganizationsRootSessionsResponseDescriptor =
    $convert.base64Decode(
        'CihEaXNhYmxlT3JnYW5pemF0aW9uc1Jvb3RTZXNzaW9uc1Jlc3BvbnNlEj0KD2VuYWJsZWRmZW'
        'F0dXJlcxjqhuZkIAMoDjIQLmlhbS5GZWF0dXJlVHlwZVIPZW5hYmxlZGZlYXR1cmVzEioKDm9y'
        'Z2FuaXphdGlvbmlkGKinppQBIAEoCVIOb3JnYW5pemF0aW9uaWQ=');

@$core.Deprecated('Use duplicateCertificateExceptionDescriptor instead')
const DuplicateCertificateException$json = {
  '1': 'DuplicateCertificateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DuplicateCertificateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateCertificateExceptionDescriptor =
    $convert.base64Decode(
        'Ch1EdXBsaWNhdGVDZXJ0aWZpY2F0ZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use duplicateSSHPublicKeyExceptionDescriptor instead')
const DuplicateSSHPublicKeyException$json = {
  '1': 'DuplicateSSHPublicKeyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DuplicateSSHPublicKeyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateSSHPublicKeyExceptionDescriptor =
    $convert.base64Decode(
        'Ch5EdXBsaWNhdGVTU0hQdWJsaWNLZXlFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use enableMFADeviceRequestDescriptor instead')
const EnableMFADeviceRequest$json = {
  '1': 'EnableMFADeviceRequest',
  '2': [
    {
      '1': 'authenticationcode1',
      '3': 444488540,
      '4': 1,
      '5': 9,
      '10': 'authenticationcode1'
    },
    {
      '1': 'authenticationcode2',
      '3': 461266159,
      '4': 1,
      '5': 9,
      '10': 'authenticationcode2'
    },
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `EnableMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableMFADeviceRequestDescriptor = $convert.base64Decode(
    'ChZFbmFibGVNRkFEZXZpY2VSZXF1ZXN0EjQKE2F1dGhlbnRpY2F0aW9uY29kZTEY3Lb50wEgAS'
    'gJUhNhdXRoZW50aWNhdGlvbmNvZGUxEjQKE2F1dGhlbnRpY2F0aW9uY29kZTIY77n52wEgASgJ'
    'UhNhdXRoZW50aWNhdGlvbmNvZGUyEiYKDHNlcmlhbG51bWJlchjlurnHASABKAlSDHNlcmlhbG'
    '51bWJlchIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated(
    'Use enableOrganizationsRootCredentialsManagementRequestDescriptor instead')
const EnableOrganizationsRootCredentialsManagementRequest$json = {
  '1': 'EnableOrganizationsRootCredentialsManagementRequest',
};

/// Descriptor for `EnableOrganizationsRootCredentialsManagementRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    enableOrganizationsRootCredentialsManagementRequestDescriptor =
    $convert.base64Decode(
        'CjNFbmFibGVPcmdhbml6YXRpb25zUm9vdENyZWRlbnRpYWxzTWFuYWdlbWVudFJlcXVlc3Q=');

@$core.Deprecated(
    'Use enableOrganizationsRootCredentialsManagementResponseDescriptor instead')
const EnableOrganizationsRootCredentialsManagementResponse$json = {
  '1': 'EnableOrganizationsRootCredentialsManagementResponse',
  '2': [
    {
      '1': 'enabledfeatures',
      '3': 211387242,
      '4': 3,
      '5': 14,
      '6': '.iam.FeatureType',
      '10': 'enabledfeatures'
    },
    {
      '1': 'organizationid',
      '3': 311006120,
      '4': 1,
      '5': 9,
      '10': 'organizationid'
    },
  ],
};

/// Descriptor for `EnableOrganizationsRootCredentialsManagementResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    enableOrganizationsRootCredentialsManagementResponseDescriptor =
    $convert.base64Decode(
        'CjRFbmFibGVPcmdhbml6YXRpb25zUm9vdENyZWRlbnRpYWxzTWFuYWdlbWVudFJlc3BvbnNlEj'
        '0KD2VuYWJsZWRmZWF0dXJlcxjqhuZkIAMoDjIQLmlhbS5GZWF0dXJlVHlwZVIPZW5hYmxlZGZl'
        'YXR1cmVzEioKDm9yZ2FuaXphdGlvbmlkGKinppQBIAEoCVIOb3JnYW5pemF0aW9uaWQ=');

@$core
    .Deprecated('Use enableOrganizationsRootSessionsRequestDescriptor instead')
const EnableOrganizationsRootSessionsRequest$json = {
  '1': 'EnableOrganizationsRootSessionsRequest',
};

/// Descriptor for `EnableOrganizationsRootSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableOrganizationsRootSessionsRequestDescriptor =
    $convert.base64Decode(
        'CiZFbmFibGVPcmdhbml6YXRpb25zUm9vdFNlc3Npb25zUmVxdWVzdA==');

@$core
    .Deprecated('Use enableOrganizationsRootSessionsResponseDescriptor instead')
const EnableOrganizationsRootSessionsResponse$json = {
  '1': 'EnableOrganizationsRootSessionsResponse',
  '2': [
    {
      '1': 'enabledfeatures',
      '3': 211387242,
      '4': 3,
      '5': 14,
      '6': '.iam.FeatureType',
      '10': 'enabledfeatures'
    },
    {
      '1': 'organizationid',
      '3': 311006120,
      '4': 1,
      '5': 9,
      '10': 'organizationid'
    },
  ],
};

/// Descriptor for `EnableOrganizationsRootSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableOrganizationsRootSessionsResponseDescriptor =
    $convert.base64Decode(
        'CidFbmFibGVPcmdhbml6YXRpb25zUm9vdFNlc3Npb25zUmVzcG9uc2USPQoPZW5hYmxlZGZlYX'
        'R1cmVzGOqG5mQgAygOMhAuaWFtLkZlYXR1cmVUeXBlUg9lbmFibGVkZmVhdHVyZXMSKgoOb3Jn'
        'YW5pemF0aW9uaWQYqKemlAEgASgJUg5vcmdhbml6YXRpb25pZA==');

@$core.Deprecated(
    'Use enableOutboundWebIdentityFederationResponseDescriptor instead')
const EnableOutboundWebIdentityFederationResponse$json = {
  '1': 'EnableOutboundWebIdentityFederationResponse',
  '2': [
    {
      '1': 'issueridentifier',
      '3': 64259590,
      '4': 1,
      '5': 9,
      '10': 'issueridentifier'
    },
  ],
};

/// Descriptor for `EnableOutboundWebIdentityFederationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    enableOutboundWebIdentityFederationResponseDescriptor =
    $convert.base64Decode(
        'CitFbmFibGVPdXRib3VuZFdlYklkZW50aXR5RmVkZXJhdGlvblJlc3BvbnNlEi0KEGlzc3Vlcm'
        'lkZW50aWZpZXIYhozSHiABKAlSEGlzc3VlcmlkZW50aWZpZXI=');

@$core.Deprecated('Use entityAlreadyExistsExceptionDescriptor instead')
const EntityAlreadyExistsException$json = {
  '1': 'EntityAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EntityAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChxFbnRpdHlBbHJlYWR5RXhpc3RzRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use entityDetailsDescriptor instead')
const EntityDetails$json = {
  '1': 'EntityDetails',
  '2': [
    {
      '1': 'entityinfo',
      '3': 401992803,
      '4': 1,
      '5': 11,
      '6': '.iam.EntityInfo',
      '10': 'entityinfo'
    },
    {
      '1': 'lastauthenticated',
      '3': 389149429,
      '4': 1,
      '5': 9,
      '10': 'lastauthenticated'
    },
  ],
};

/// Descriptor for `EntityDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityDetailsDescriptor = $convert.base64Decode(
    'Cg1FbnRpdHlEZXRhaWxzEjMKCmVudGl0eWluZm8Y49jXvwEgASgLMg8uaWFtLkVudGl0eUluZm'
    '9SCmVudGl0eWluZm8SMAoRbGFzdGF1dGhlbnRpY2F0ZWQY9eXHuQEgASgJUhFsYXN0YXV0aGVu'
    'dGljYXRlZA==');

@$core.Deprecated('Use entityInfoDescriptor instead')
const EntityInfo$json = {
  '1': 'EntityInfo',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.iam.policyOwnerEntityType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `EntityInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityInfoDescriptor = $convert.base64Decode(
    'CgpFbnRpdHlJbmZvEhQKA2Fybhidm+2/ASABKAlSA2FybhISCgJpZBiB8qK3ASABKAlSAmlkEh'
    'UKBG5hbWUYh+aBfyABKAlSBG5hbWUSFQoEcGF0aBjXyJtbIAEoCVIEcGF0aBIyCgR0eXBlGO6g'
    '14oBIAEoDjIaLmlhbS5wb2xpY3lPd25lckVudGl0eVR5cGVSBHR5cGU=');

@$core
    .Deprecated('Use entityTemporarilyUnmodifiableExceptionDescriptor instead')
const EntityTemporarilyUnmodifiableException$json = {
  '1': 'EntityTemporarilyUnmodifiableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EntityTemporarilyUnmodifiableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityTemporarilyUnmodifiableExceptionDescriptor =
    $convert.base64Decode(
        'CiZFbnRpdHlUZW1wb3JhcmlseVVubW9kaWZpYWJsZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyC'
        'cgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use errorDetailsDescriptor instead')
const ErrorDetails$json = {
  '1': 'ErrorDetails',
  '2': [
    {'1': 'code', '3': 425572629, '4': 1, '5': 9, '10': 'code'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ErrorDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorDetailsDescriptor = $convert.base64Decode(
    'CgxFcnJvckRldGFpbHMSFgoEY29kZRiV8vbKASABKAlSBGNvZGUSGwoHbWVzc2FnZRiFs7twIA'
    'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use evaluationResultDescriptor instead')
const EvaluationResult$json = {
  '1': 'EvaluationResult',
  '2': [
    {
      '1': 'evalactionname',
      '3': 342009271,
      '4': 1,
      '5': 9,
      '10': 'evalactionname'
    },
    {
      '1': 'evaldecision',
      '3': 346306806,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyEvaluationDecisionType',
      '10': 'evaldecision'
    },
    {
      '1': 'evaldecisiondetails',
      '3': 70732822,
      '4': 3,
      '5': 11,
      '6': '.iam.EvaluationResult.EvaldecisiondetailsEntry',
      '10': 'evaldecisiondetails'
    },
    {
      '1': 'evalresourcename',
      '3': 9391249,
      '4': 1,
      '5': 9,
      '10': 'evalresourcename'
    },
    {
      '1': 'matchedstatements',
      '3': 389771374,
      '4': 3,
      '5': 11,
      '6': '.iam.Statement',
      '10': 'matchedstatements'
    },
    {
      '1': 'missingcontextvalues',
      '3': 526483253,
      '4': 3,
      '5': 9,
      '10': 'missingcontextvalues'
    },
    {
      '1': 'organizationsdecisiondetail',
      '3': 530188257,
      '4': 1,
      '5': 11,
      '6': '.iam.OrganizationsDecisionDetail',
      '10': 'organizationsdecisiondetail'
    },
    {
      '1': 'permissionsboundarydecisiondetail',
      '3': 480692611,
      '4': 1,
      '5': 11,
      '6': '.iam.PermissionsBoundaryDecisionDetail',
      '10': 'permissionsboundarydecisiondetail'
    },
    {
      '1': 'resourcespecificresults',
      '3': 439535362,
      '4': 3,
      '5': 11,
      '6': '.iam.ResourceSpecificResult',
      '10': 'resourcespecificresults'
    },
  ],
  '3': [EvaluationResult_EvaldecisiondetailsEntry$json],
};

@$core.Deprecated('Use evaluationResultDescriptor instead')
const EvaluationResult_EvaldecisiondetailsEntry$json = {
  '1': 'EvaldecisiondetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyEvaluationDecisionType',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `EvaluationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List evaluationResultDescriptor = $convert.base64Decode(
    'ChBFdmFsdWF0aW9uUmVzdWx0EioKDmV2YWxhY3Rpb25uYW1lGLfLiqMBIAEoCVIOZXZhbGFjdG'
    'lvbm5hbWUSSQoMZXZhbGRlY2lzaW9uGPbxkKUBIAEoDjIhLmlhbS5Qb2xpY3lFdmFsdWF0aW9u'
    'RGVjaXNpb25UeXBlUgxldmFsZGVjaXNpb24SYwoTZXZhbGRlY2lzaW9uZGV0YWlscxiWmN0hIA'
    'MoCzIuLmlhbS5FdmFsdWF0aW9uUmVzdWx0LkV2YWxkZWNpc2lvbmRldGFpbHNFbnRyeVITZXZh'
    'bGRlY2lzaW9uZGV0YWlscxItChBldmFscmVzb3VyY2VuYW1lGJGZvQQgASgJUhBldmFscmVzb3'
    'VyY2VuYW1lEkAKEW1hdGNoZWRzdGF0ZW1lbnRzGO7g7bkBIAMoCzIOLmlhbS5TdGF0ZW1lbnRS'
    'EW1hdGNoZWRzdGF0ZW1lbnRzEjYKFG1pc3Npbmdjb250ZXh0dmFsdWVzGLX+hfsBIAMoCVIUbW'
    'lzc2luZ2NvbnRleHR2YWx1ZXMSZgobb3JnYW5pemF0aW9uc2RlY2lzaW9uZGV0YWlsGOGP6PwB'
    'IAEoCzIgLmlhbS5Pcmdhbml6YXRpb25zRGVjaXNpb25EZXRhaWxSG29yZ2FuaXphdGlvbnNkZW'
    'Npc2lvbmRldGFpbBJ4CiFwZXJtaXNzaW9uc2JvdW5kYXJ5ZGVjaXNpb25kZXRhaWwYg5Ob5QEg'
    'ASgLMiYuaWFtLlBlcm1pc3Npb25zQm91bmRhcnlEZWNpc2lvbkRldGFpbFIhcGVybWlzc2lvbn'
    'Nib3VuZGFyeWRlY2lzaW9uZGV0YWlsElkKF3Jlc291cmNlc3BlY2lmaWNyZXN1bHRzGIKOy9EB'
    'IAMoCzIbLmlhbS5SZXNvdXJjZVNwZWNpZmljUmVzdWx0UhdyZXNvdXJjZXNwZWNpZmljcmVzdW'
    'x0cxppChhFdmFsZGVjaXNpb25kZXRhaWxzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSNwoFdmFs'
    'dWUYAiABKA4yIS5pYW0uUG9saWN5RXZhbHVhdGlvbkRlY2lzaW9uVHlwZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use featureDisabledExceptionDescriptor instead')
const FeatureDisabledException$json = {
  '1': 'FeatureDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FeatureDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List featureDisabledExceptionDescriptor =
    $convert.base64Decode(
        'ChhGZWF0dXJlRGlzYWJsZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use featureEnabledExceptionDescriptor instead')
const FeatureEnabledException$json = {
  '1': 'FeatureEnabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FeatureEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List featureEnabledExceptionDescriptor =
    $convert.base64Decode(
        'ChdGZWF0dXJlRW5hYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use generateCredentialReportResponseDescriptor instead')
const GenerateCredentialReportResponse$json = {
  '1': 'GenerateCredentialReportResponse',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.iam.ReportStateType',
      '10': 'state'
    },
  ],
};

/// Descriptor for `GenerateCredentialReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateCredentialReportResponseDescriptor =
    $convert.base64Decode(
        'CiBHZW5lcmF0ZUNyZWRlbnRpYWxSZXBvcnRSZXNwb25zZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IA'
        'EoCVILZGVzY3JpcHRpb24SLgoFc3RhdGUYl8my7wEgASgOMhQuaWFtLlJlcG9ydFN0YXRlVHlw'
        'ZVIFc3RhdGU=');

@$core.Deprecated(
    'Use generateOrganizationsAccessReportRequestDescriptor instead')
const GenerateOrganizationsAccessReportRequest$json = {
  '1': 'GenerateOrganizationsAccessReportRequest',
  '2': [
    {'1': 'entitypath', '3': 93829450, '4': 1, '5': 9, '10': 'entitypath'},
    {
      '1': 'organizationspolicyid',
      '3': 16852361,
      '4': 1,
      '5': 9,
      '10': 'organizationspolicyid'
    },
  ],
};

/// Descriptor for `GenerateOrganizationsAccessReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateOrganizationsAccessReportRequestDescriptor =
    $convert.base64Decode(
        'CihHZW5lcmF0ZU9yZ2FuaXphdGlvbnNBY2Nlc3NSZXBvcnRSZXF1ZXN0EiEKCmVudGl0eXBhdG'
        'gYyvLeLCABKAlSCmVudGl0eXBhdGgSNwoVb3JnYW5pemF0aW9uc3BvbGljeWlkGInLhAggASgJ'
        'UhVvcmdhbml6YXRpb25zcG9saWN5aWQ=');

@$core.Deprecated(
    'Use generateOrganizationsAccessReportResponseDescriptor instead')
const GenerateOrganizationsAccessReportResponse$json = {
  '1': 'GenerateOrganizationsAccessReportResponse',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `GenerateOrganizationsAccessReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    generateOrganizationsAccessReportResponseDescriptor = $convert.base64Decode(
        'CilHZW5lcmF0ZU9yZ2FuaXphdGlvbnNBY2Nlc3NSZXBvcnRSZXNwb25zZRIXCgVqb2JpZBjS1N'
        '0zIAEoCVIFam9iaWQ=');

@$core.Deprecated(
    'Use generateServiceLastAccessedDetailsRequestDescriptor instead')
const GenerateServiceLastAccessedDetailsRequest$json = {
  '1': 'GenerateServiceLastAccessedDetailsRequest',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'granularity',
      '3': 498840280,
      '4': 1,
      '5': 14,
      '6': '.iam.AccessAdvisorUsageGranularityType',
      '10': 'granularity'
    },
  ],
};

/// Descriptor for `GenerateServiceLastAccessedDetailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    generateServiceLastAccessedDetailsRequestDescriptor = $convert.base64Decode(
        'CilHZW5lcmF0ZVNlcnZpY2VMYXN0QWNjZXNzZWREZXRhaWxzUmVxdWVzdBIUCgNhcm4YnZvtvw'
        'EgASgJUgNhcm4STAoLZ3JhbnVsYXJpdHkY2OXu7QEgASgOMiYuaWFtLkFjY2Vzc0Fkdmlzb3JV'
        'c2FnZUdyYW51bGFyaXR5VHlwZVILZ3JhbnVsYXJpdHk=');

@$core.Deprecated(
    'Use generateServiceLastAccessedDetailsResponseDescriptor instead')
const GenerateServiceLastAccessedDetailsResponse$json = {
  '1': 'GenerateServiceLastAccessedDetailsResponse',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `GenerateServiceLastAccessedDetailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    generateServiceLastAccessedDetailsResponseDescriptor =
    $convert.base64Decode(
        'CipHZW5lcmF0ZVNlcnZpY2VMYXN0QWNjZXNzZWREZXRhaWxzUmVzcG9uc2USFwoFam9iaWQY0t'
        'TdMyABKAlSBWpvYmlk');

@$core.Deprecated('Use getAccessKeyLastUsedRequestDescriptor instead')
const GetAccessKeyLastUsedRequest$json = {
  '1': 'GetAccessKeyLastUsedRequest',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
  ],
};

/// Descriptor for `GetAccessKeyLastUsedRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccessKeyLastUsedRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRBY2Nlc3NLZXlMYXN0VXNlZFJlcXVlc3QSJAoLYWNjZXNza2V5aWQYoLe32AEgASgJUg'
        'thY2Nlc3NrZXlpZA==');

@$core.Deprecated('Use getAccessKeyLastUsedResponseDescriptor instead')
const GetAccessKeyLastUsedResponse$json = {
  '1': 'GetAccessKeyLastUsedResponse',
  '2': [
    {
      '1': 'accesskeylastused',
      '3': 104233392,
      '4': 1,
      '5': 11,
      '6': '.iam.AccessKeyLastUsed',
      '10': 'accesskeylastused'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetAccessKeyLastUsedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccessKeyLastUsedResponseDescriptor =
    $convert.base64Decode(
        'ChxHZXRBY2Nlc3NLZXlMYXN0VXNlZFJlc3BvbnNlEkcKEWFjY2Vzc2tleWxhc3R1c2VkGLDz2T'
        'EgASgLMhYuaWFtLkFjY2Vzc0tleUxhc3RVc2VkUhFhY2Nlc3NrZXlsYXN0dXNlZBIeCgh1c2Vy'
        'bmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use getAccountAuthorizationDetailsRequestDescriptor instead')
const GetAccountAuthorizationDetailsRequest$json = {
  '1': 'GetAccountAuthorizationDetailsRequest',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 3,
      '5': 14,
      '6': '.iam.EntityType',
      '10': 'filter'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `GetAccountAuthorizationDetailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountAuthorizationDetailsRequestDescriptor =
    $convert.base64Decode(
        'CiVHZXRBY2NvdW50QXV0aG9yaXphdGlvbkRldGFpbHNSZXF1ZXN0EisKBmZpbHRlchiYgaelAS'
        'ADKA4yDy5pYW0uRW50aXR5VHlwZVIGZmlsdGVyEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2Vy'
        'Eh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core
    .Deprecated('Use getAccountAuthorizationDetailsResponseDescriptor instead')
const GetAccountAuthorizationDetailsResponse$json = {
  '1': 'GetAccountAuthorizationDetailsResponse',
  '2': [
    {
      '1': 'groupdetaillist',
      '3': 351520026,
      '4': 3,
      '5': 11,
      '6': '.iam.GroupDetail',
      '10': 'groupdetaillist'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.iam.ManagedPolicyDetail',
      '10': 'policies'
    },
    {
      '1': 'roledetaillist',
      '3': 282084497,
      '4': 3,
      '5': 11,
      '6': '.iam.RoleDetail',
      '10': 'roledetaillist'
    },
    {
      '1': 'userdetaillist',
      '3': 495221388,
      '4': 3,
      '5': 11,
      '6': '.iam.UserDetail',
      '10': 'userdetaillist'
    },
  ],
};

/// Descriptor for `GetAccountAuthorizationDetailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountAuthorizationDetailsResponseDescriptor = $convert.base64Decode(
    'CiZHZXRBY2NvdW50QXV0aG9yaXphdGlvbkRldGFpbHNSZXNwb25zZRI+Cg9ncm91cGRldGFpbG'
    'xpc3QYmorPpwEgAygLMhAuaWFtLkdyb3VwRGV0YWlsUg9ncm91cGRldGFpbGxpc3QSIwoLaXN0'
    'cnVuY2F0ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2'
    'VyEjcKCHBvbGljaWVzGJisihMgAygLMhguaWFtLk1hbmFnZWRQb2xpY3lEZXRhaWxSCHBvbGlj'
    'aWVzEjsKDnJvbGVkZXRhaWxsaXN0GJGJwYYBIAMoCzIPLmlhbS5Sb2xlRGV0YWlsUg5yb2xlZG'
    'V0YWlsbGlzdBI7Cg51c2VyZGV0YWlsbGlzdBiM9ZHsASADKAsyDy5pYW0uVXNlckRldGFpbFIO'
    'dXNlcmRldGFpbGxpc3Q=');

@$core.Deprecated('Use getAccountPasswordPolicyResponseDescriptor instead')
const GetAccountPasswordPolicyResponse$json = {
  '1': 'GetAccountPasswordPolicyResponse',
  '2': [
    {
      '1': 'passwordpolicy',
      '3': 235700227,
      '4': 1,
      '5': 11,
      '6': '.iam.PasswordPolicy',
      '10': 'passwordpolicy'
    },
  ],
};

/// Descriptor for `GetAccountPasswordPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountPasswordPolicyResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRBY2NvdW50UGFzc3dvcmRQb2xpY3lSZXNwb25zZRI+Cg5wYXNzd29yZHBvbGljeRiDgL'
        'JwIAEoCzITLmlhbS5QYXNzd29yZFBvbGljeVIOcGFzc3dvcmRwb2xpY3k=');

@$core.Deprecated('Use getAccountSummaryResponseDescriptor instead')
const GetAccountSummaryResponse$json = {
  '1': 'GetAccountSummaryResponse',
  '2': [
    {
      '1': 'summarymap',
      '3': 121019076,
      '4': 3,
      '5': 11,
      '6': '.iam.GetAccountSummaryResponse.SummarymapEntry',
      '10': 'summarymap'
    },
  ],
  '3': [GetAccountSummaryResponse_SummarymapEntry$json],
};

@$core.Deprecated('Use getAccountSummaryResponseDescriptor instead')
const GetAccountSummaryResponse_SummarymapEntry$json = {
  '1': 'SummarymapEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetAccountSummaryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountSummaryResponseDescriptor = $convert.base64Decode(
    'ChlHZXRBY2NvdW50U3VtbWFyeVJlc3BvbnNlElEKCnN1bW1hcnltYXAYxLXaOSADKAsyLi5pYW'
    '0uR2V0QWNjb3VudFN1bW1hcnlSZXNwb25zZS5TdW1tYXJ5bWFwRW50cnlSCnN1bW1hcnltYXAa'
    'PQoPU3VtbWFyeW1hcEVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgFUgV2YW'
    'x1ZToCOAE=');

@$core.Deprecated('Use getContextKeysForCustomPolicyRequestDescriptor instead')
const GetContextKeysForCustomPolicyRequest$json = {
  '1': 'GetContextKeysForCustomPolicyRequest',
  '2': [
    {
      '1': 'policyinputlist',
      '3': 320766346,
      '4': 3,
      '5': 9,
      '10': 'policyinputlist'
    },
  ],
};

/// Descriptor for `GetContextKeysForCustomPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContextKeysForCustomPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRDb250ZXh0S2V5c0ZvckN1c3RvbVBvbGljeVJlcXVlc3QSLAoPcG9saWN5aW5wdXRsaX'
        'N0GIqD+pgBIAMoCVIPcG9saWN5aW5wdXRsaXN0');

@$core.Deprecated('Use getContextKeysForPolicyResponseDescriptor instead')
const GetContextKeysForPolicyResponse$json = {
  '1': 'GetContextKeysForPolicyResponse',
  '2': [
    {
      '1': 'contextkeynames',
      '3': 16393914,
      '4': 3,
      '5': 9,
      '10': 'contextkeynames'
    },
  ],
};

/// Descriptor for `GetContextKeysForPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContextKeysForPolicyResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXRDb250ZXh0S2V5c0ZvclBvbGljeVJlc3BvbnNlEisKD2NvbnRleHRrZXluYW1lcxi6ze'
        'gHIAMoCVIPY29udGV4dGtleW5hbWVz');

@$core
    .Deprecated('Use getContextKeysForPrincipalPolicyRequestDescriptor instead')
const GetContextKeysForPrincipalPolicyRequest$json = {
  '1': 'GetContextKeysForPrincipalPolicyRequest',
  '2': [
    {
      '1': 'policyinputlist',
      '3': 320766346,
      '4': 3,
      '5': 9,
      '10': 'policyinputlist'
    },
    {
      '1': 'policysourcearn',
      '3': 84285350,
      '4': 1,
      '5': 9,
      '10': 'policysourcearn'
    },
  ],
};

/// Descriptor for `GetContextKeysForPrincipalPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContextKeysForPrincipalPolicyRequestDescriptor =
    $convert.base64Decode(
        'CidHZXRDb250ZXh0S2V5c0ZvclByaW5jaXBhbFBvbGljeVJlcXVlc3QSLAoPcG9saWN5aW5wdX'
        'RsaXN0GIqD+pgBIAMoCVIPcG9saWN5aW5wdXRsaXN0EisKD3BvbGljeXNvdXJjZWFybhimr5go'
        'IAEoCVIPcG9saWN5c291cmNlYXJu');

@$core.Deprecated('Use getCredentialReportResponseDescriptor instead')
const GetCredentialReportResponse$json = {
  '1': 'GetCredentialReportResponse',
  '2': [
    {'1': 'content', '3': 23568227, '4': 1, '5': 12, '10': 'content'},
    {
      '1': 'generatedtime',
      '3': 88097636,
      '4': 1,
      '5': 9,
      '10': 'generatedtime'
    },
    {
      '1': 'reportformat',
      '3': 67774421,
      '4': 1,
      '5': 14,
      '6': '.iam.ReportFormatType',
      '10': 'reportformat'
    },
  ],
};

/// Descriptor for `GetCredentialReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCredentialReportResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRDcmVkZW50aWFsUmVwb3J0UmVzcG9uc2USGwoHY29udGVudBjjvp4LIAEoDFIHY29udG'
        'VudBInCg1nZW5lcmF0ZWR0aW1lGOSGgSogASgJUg1nZW5lcmF0ZWR0aW1lEjwKDHJlcG9ydGZv'
        'cm1hdBjVz6ggIAEoDjIVLmlhbS5SZXBvcnRGb3JtYXRUeXBlUgxyZXBvcnRmb3JtYXQ=');

@$core.Deprecated('Use getDelegationRequestRequestDescriptor instead')
const GetDelegationRequestRequest$json = {
  '1': 'GetDelegationRequestRequest',
  '2': [
    {
      '1': 'delegationpermissioncheck',
      '3': 10515097,
      '4': 1,
      '5': 8,
      '10': 'delegationpermissioncheck'
    },
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
  ],
};

/// Descriptor for `GetDelegationRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDelegationRequestRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXREZWxlZ2F0aW9uUmVxdWVzdFJlcXVlc3QSPwoZZGVsZWdhdGlvbnBlcm1pc3Npb25jaG'
        'VjaxiZ5YEFIAEoCFIZZGVsZWdhdGlvbnBlcm1pc3Npb25jaGVjaxI0ChNkZWxlZ2F0aW9ucmVx'
        'dWVzdGlkGJ63gIUBIAEoCVITZGVsZWdhdGlvbnJlcXVlc3RpZA==');

@$core.Deprecated('Use getDelegationRequestResponseDescriptor instead')
const GetDelegationRequestResponse$json = {
  '1': 'GetDelegationRequestResponse',
  '2': [
    {
      '1': 'delegationrequest',
      '3': 331878529,
      '4': 1,
      '5': 11,
      '6': '.iam.DelegationRequest',
      '10': 'delegationrequest'
    },
    {
      '1': 'permissioncheckresult',
      '3': 120682140,
      '4': 1,
      '5': 14,
      '6': '.iam.permissionCheckResultType',
      '10': 'permissioncheckresult'
    },
    {
      '1': 'permissioncheckstatus',
      '3': 164847565,
      '4': 1,
      '5': 14,
      '6': '.iam.permissionCheckStatusType',
      '10': 'permissioncheckstatus'
    },
  ],
};

/// Descriptor for `GetDelegationRequestResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDelegationRequestResponseDescriptor = $convert.base64Decode(
    'ChxHZXREZWxlZ2F0aW9uUmVxdWVzdFJlc3BvbnNlEkgKEWRlbGVnYXRpb25yZXF1ZXN0GIGhoJ'
    '4BIAEoCzIWLmlhbS5EZWxlZ2F0aW9uUmVxdWVzdFIRZGVsZWdhdGlvbnJlcXVlc3QSVwoVcGVy'
    'bWlzc2lvbmNoZWNrcmVzdWx0GJztxTkgASgOMh4uaWFtLnBlcm1pc3Npb25DaGVja1Jlc3VsdF'
    'R5cGVSFXBlcm1pc3Npb25jaGVja3Jlc3VsdBJXChVwZXJtaXNzaW9uY2hlY2tzdGF0dXMYzb/N'
    'TiABKA4yHi5pYW0ucGVybWlzc2lvbkNoZWNrU3RhdHVzVHlwZVIVcGVybWlzc2lvbmNoZWNrc3'
    'RhdHVz');

@$core.Deprecated('Use getGroupPolicyRequestDescriptor instead')
const GetGroupPolicyRequest$json = {
  '1': 'GetGroupPolicyRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `GetGroupPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGroupPolicyRequestDescriptor = $convert.base64Decode(
    'ChVHZXRHcm91cFBvbGljeVJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW'
    '1lEiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWU=');

@$core.Deprecated('Use getGroupPolicyResponseDescriptor instead')
const GetGroupPolicyResponse$json = {
  '1': 'GetGroupPolicyResponse',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `GetGroupPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGroupPolicyResponseDescriptor = $convert.base64Decode(
    'ChZHZXRHcm91cFBvbGljeVJlc3BvbnNlEiAKCWdyb3VwbmFtZRjIyqCqASABKAlSCWdyb3Vwbm'
    'FtZRIpCg5wb2xpY3lkb2N1bWVudBjLrsFxIAEoCVIOcG9saWN5ZG9jdW1lbnQSIQoKcG9saWN5'
    'bmFtZRi99Yd/IAEoCVIKcG9saWN5bmFtZQ==');

@$core.Deprecated('Use getGroupRequestDescriptor instead')
const GetGroupRequest$json = {
  '1': 'GetGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `GetGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGroupRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRHcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEhkKBm'
    '1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use getGroupResponseDescriptor instead')
const GetGroupResponse$json = {
  '1': 'GetGroupResponse',
  '2': [
    {
      '1': 'group',
      '3': 91525165,
      '4': 1,
      '5': 11,
      '6': '.iam.Group',
      '10': 'group'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'users',
      '3': 112472756,
      '4': 3,
      '5': 11,
      '6': '.iam.User',
      '10': 'users'
    },
  ],
};

/// Descriptor for `GetGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGroupResponseDescriptor = $convert.base64Decode(
    'ChBHZXRHcm91cFJlc3BvbnNlEiMKBWdyb3VwGK2g0isgASgLMgouaWFtLkdyb3VwUgVncm91cB'
    'IjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJ'
    'UgZtYXJrZXISIgoFdXNlcnMYtOXQNSADKAsyCS5pYW0uVXNlclIFdXNlcnM=');

@$core.Deprecated('Use getHumanReadableSummaryRequestDescriptor instead')
const GetHumanReadableSummaryRequest$json = {
  '1': 'GetHumanReadableSummaryRequest',
  '2': [
    {'1': 'entityarn', '3': 422623110, '4': 1, '5': 9, '10': 'entityarn'},
    {'1': 'locale', '3': 493771704, '4': 1, '5': 9, '10': 'locale'},
  ],
};

/// Descriptor for `GetHumanReadableSummaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHumanReadableSummaryRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXRIdW1hblJlYWRhYmxlU3VtbWFyeVJlcXVlc3QSIAoJZW50aXR5YXJuGIbvwskBIAEoCV'
        'IJZW50aXR5YXJuEhoKBmxvY2FsZRi4t7nrASABKAlSBmxvY2FsZQ==');

@$core.Deprecated('Use getHumanReadableSummaryResponseDescriptor instead')
const GetHumanReadableSummaryResponse$json = {
  '1': 'GetHumanReadableSummaryResponse',
  '2': [
    {'1': 'locale', '3': 493771704, '4': 1, '5': 9, '10': 'locale'},
    {
      '1': 'summarycontent',
      '3': 116389717,
      '4': 1,
      '5': 9,
      '10': 'summarycontent'
    },
    {
      '1': 'summarystate',
      '3': 124704997,
      '4': 1,
      '5': 14,
      '6': '.iam.summaryStateType',
      '10': 'summarystate'
    },
  ],
};

/// Descriptor for `GetHumanReadableSummaryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHumanReadableSummaryResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXRIdW1hblJlYWRhYmxlU3VtbWFyeVJlc3BvbnNlEhoKBmxvY2FsZRi4t7nrASABKAlSBm'
        'xvY2FsZRIpCg5zdW1tYXJ5Y29udGVudBjV7r83IAEoCVIOc3VtbWFyeWNvbnRlbnQSPAoMc3Vt'
        'bWFyeXN0YXRlGOWxuzsgASgOMhUuaWFtLnN1bW1hcnlTdGF0ZVR5cGVSDHN1bW1hcnlzdGF0ZQ'
        '==');

@$core.Deprecated('Use getInstanceProfileRequestDescriptor instead')
const GetInstanceProfileRequest$json = {
  '1': 'GetInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
  ],
};

/// Descriptor for `GetInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInstanceProfileRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRJbnN0YW5jZVByb2ZpbGVSZXF1ZXN0EjQKE2luc3RhbmNlcHJvZmlsZW5hbWUY7c682g'
        'EgASgJUhNpbnN0YW5jZXByb2ZpbGVuYW1l');

@$core.Deprecated('Use getInstanceProfileResponseDescriptor instead')
const GetInstanceProfileResponse$json = {
  '1': 'GetInstanceProfileResponse',
  '2': [
    {
      '1': 'instanceprofile',
      '3': 57726800,
      '4': 1,
      '5': 11,
      '6': '.iam.InstanceProfile',
      '10': 'instanceprofile'
    },
  ],
};

/// Descriptor for `GetInstanceProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInstanceProfileResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRJbnN0YW5jZVByb2ZpbGVSZXNwb25zZRJBCg9pbnN0YW5jZXByb2ZpbGUY0K7DGyABKA'
        'syFC5pYW0uSW5zdGFuY2VQcm9maWxlUg9pbnN0YW5jZXByb2ZpbGU=');

@$core.Deprecated('Use getLoginProfileRequestDescriptor instead')
const GetLoginProfileRequest$json = {
  '1': 'GetLoginProfileRequest',
  '2': [
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetLoginProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLoginProfileRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRMb2dpblByb2ZpbGVSZXF1ZXN0Eh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbW'
        'U=');

@$core.Deprecated('Use getLoginProfileResponseDescriptor instead')
const GetLoginProfileResponse$json = {
  '1': 'GetLoginProfileResponse',
  '2': [
    {
      '1': 'loginprofile',
      '3': 465493040,
      '4': 1,
      '5': 11,
      '6': '.iam.LoginProfile',
      '10': 'loginprofile'
    },
  ],
};

/// Descriptor for `GetLoginProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLoginProfileResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRMb2dpblByb2ZpbGVSZXNwb25zZRI5Cgxsb2dpbnByb2ZpbGUYsLj73QEgASgLMhEuaW'
        'FtLkxvZ2luUHJvZmlsZVIMbG9naW5wcm9maWxl');

@$core.Deprecated('Use getMFADeviceRequestDescriptor instead')
const GetMFADeviceRequest$json = {
  '1': 'GetMFADeviceRequest',
  '2': [
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMFADeviceRequestDescriptor = $convert.base64Decode(
    'ChNHZXRNRkFEZXZpY2VSZXF1ZXN0EiYKDHNlcmlhbG51bWJlchjlurnHASABKAlSDHNlcmlhbG'
    '51bWJlchIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use getMFADeviceResponseDescriptor instead')
const GetMFADeviceResponse$json = {
  '1': 'GetMFADeviceResponse',
  '2': [
    {
      '1': 'certifications',
      '3': 399647051,
      '4': 3,
      '5': 11,
      '6': '.iam.GetMFADeviceResponse.CertificationsEntry',
      '10': 'certifications'
    },
    {'1': 'enabledate', '3': 99725723, '4': 1, '5': 9, '10': 'enabledate'},
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [GetMFADeviceResponse_CertificationsEntry$json],
};

@$core.Deprecated('Use getMFADeviceResponseDescriptor instead')
const GetMFADeviceResponse_CertificationsEntry$json = {
  '1': 'CertificationsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetMFADeviceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMFADeviceResponseDescriptor = $convert.base64Decode(
    'ChRHZXRNRkFEZXZpY2VSZXNwb25zZRJZCg5jZXJ0aWZpY2F0aW9ucxjLwsi+ASADKAsyLS5pYW'
    '0uR2V0TUZBRGV2aWNlUmVzcG9uc2UuQ2VydGlmaWNhdGlvbnNFbnRyeVIOY2VydGlmaWNhdGlv'
    'bnMSIQoKZW5hYmxlZGF0ZRib48YvIAEoCVIKZW5hYmxlZGF0ZRImCgxzZXJpYWxudW1iZXIY5b'
    'q5xwEgASgJUgxzZXJpYWxudW1iZXISHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZRpB'
    'ChNDZXJ0aWZpY2F0aW9uc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZToCOAE=');

@$core.Deprecated('Use getOpenIDConnectProviderRequestDescriptor instead')
const GetOpenIDConnectProviderRequest$json = {
  '1': 'GetOpenIDConnectProviderRequest',
  '2': [
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
  ],
};

/// Descriptor for `GetOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpenIDConnectProviderRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRPcGVuSURDb25uZWN0UHJvdmlkZXJSZXF1ZXN0Ej4KGG9wZW5pZGNvbm5lY3Rwcm92aW'
        'RlcmFybhjj84/pASABKAlSGG9wZW5pZGNvbm5lY3Rwcm92aWRlcmFybg==');

@$core.Deprecated('Use getOpenIDConnectProviderResponseDescriptor instead')
const GetOpenIDConnectProviderResponse$json = {
  '1': 'GetOpenIDConnectProviderResponse',
  '2': [
    {'1': 'clientidlist', '3': 428113132, '4': 3, '5': 9, '10': 'clientidlist'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {
      '1': 'thumbprintlist',
      '3': 88528351,
      '4': 3,
      '5': 9,
      '10': 'thumbprintlist'
    },
    {'1': 'url', '3': 354018239, '4': 1, '5': 9, '10': 'url'},
  ],
};

/// Descriptor for `GetOpenIDConnectProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpenIDConnectProviderResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRPcGVuSURDb25uZWN0UHJvdmlkZXJSZXNwb25zZRImCgxjbGllbnRpZGxpc3QY7PmRzA'
        'EgAygJUgxjbGllbnRpZGxpc3QSIQoKY3JlYXRlZGF0ZRiSufwRIAEoCVIKY3JlYXRlZGF0ZRIg'
        'CgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3MSKQoOdGh1bWJwcmludGxpc3QY36ubKi'
        'ADKAlSDnRodW1icHJpbnRsaXN0EhQKA3VybBi/x+eoASABKAlSA3VybA==');

@$core.Deprecated('Use getOrganizationsAccessReportRequestDescriptor instead')
const GetOrganizationsAccessReportRequest$json = {
  '1': 'GetOrganizationsAccessReportRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'sortkey',
      '3': 407935321,
      '4': 1,
      '5': 14,
      '6': '.iam.sortKeyType',
      '10': 'sortkey'
    },
  ],
};

/// Descriptor for `GetOrganizationsAccessReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOrganizationsAccessReportRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXRPcmdhbml6YXRpb25zQWNjZXNzUmVwb3J0UmVxdWVzdBIXCgVqb2JpZBjS1N0zIAEoCV'
        'IFam9iaWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgF'
        'UghtYXhpdGVtcxIuCgdzb3J0a2V5GNmywsIBIAEoDjIQLmlhbS5zb3J0S2V5VHlwZVIHc29ydG'
        'tleQ==');

@$core.Deprecated('Use getOrganizationsAccessReportResponseDescriptor instead')
const GetOrganizationsAccessReportResponse$json = {
  '1': 'GetOrganizationsAccessReportResponse',
  '2': [
    {
      '1': 'accessdetails',
      '3': 80267372,
      '4': 3,
      '5': 11,
      '6': '.iam.AccessDetail',
      '10': 'accessdetails'
    },
    {
      '1': 'errordetails',
      '3': 369268426,
      '4': 1,
      '5': 11,
      '6': '.iam.ErrorDetails',
      '10': 'errordetails'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'jobcompletiondate',
      '3': 517991255,
      '4': 1,
      '5': 9,
      '10': 'jobcompletiondate'
    },
    {
      '1': 'jobcreationdate',
      '3': 201415034,
      '4': 1,
      '5': 9,
      '10': 'jobcreationdate'
    },
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.iam.jobStatusType',
      '10': 'jobstatus'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'numberofservicesaccessible',
      '3': 27175708,
      '4': 1,
      '5': 5,
      '10': 'numberofservicesaccessible'
    },
    {
      '1': 'numberofservicesnotaccessed',
      '3': 127942410,
      '4': 1,
      '5': 5,
      '10': 'numberofservicesnotaccessed'
    },
  ],
};

/// Descriptor for `GetOrganizationsAccessReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOrganizationsAccessReportResponseDescriptor = $convert.base64Decode(
    'CiRHZXRPcmdhbml6YXRpb25zQWNjZXNzUmVwb3J0UmVzcG9uc2USOgoNYWNjZXNzZGV0YWlscx'
    'jskKMmIAMoCzIRLmlhbS5BY2Nlc3NEZXRhaWxSDWFjY2Vzc2RldGFpbHMSOQoMZXJyb3JkZXRh'
    'aWxzGMqtirABIAEoCzIRLmlhbS5FcnJvckRldGFpbHNSDGVycm9yZGV0YWlscxIjCgtpc3RydW'
    '5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSMAoRam9iY29tcGxldGlvbmRhdGUY19b/9gEg'
    'ASgJUhFqb2Jjb21wbGV0aW9uZGF0ZRIrCg9qb2JjcmVhdGlvbmRhdGUY+rKFYCABKAlSD2pvYm'
    'NyZWF0aW9uZGF0ZRIzCglqb2JzdGF0dXMYx5z7MyABKA4yEi5pYW0uam9iU3RhdHVzVHlwZVIJ'
    'am9ic3RhdHVzEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEkEKGm51bWJlcm9mc2VydmljZX'
    'NhY2Nlc3NpYmxlGJzW+gwgASgFUhpudW1iZXJvZnNlcnZpY2VzYWNjZXNzaWJsZRJDChtudW1i'
    'ZXJvZnNlcnZpY2Vzbm90YWNjZXNzZWQYiv6APSABKAVSG251bWJlcm9mc2VydmljZXNub3RhY2'
    'Nlc3NlZA==');

@$core.Deprecated(
    'Use getOutboundWebIdentityFederationInfoResponseDescriptor instead')
const GetOutboundWebIdentityFederationInfoResponse$json = {
  '1': 'GetOutboundWebIdentityFederationInfoResponse',
  '2': [
    {
      '1': 'issueridentifier',
      '3': 64259590,
      '4': 1,
      '5': 9,
      '10': 'issueridentifier'
    },
    {
      '1': 'jwtvendingenabled',
      '3': 148790521,
      '4': 1,
      '5': 8,
      '10': 'jwtvendingenabled'
    },
  ],
};

/// Descriptor for `GetOutboundWebIdentityFederationInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getOutboundWebIdentityFederationInfoResponseDescriptor =
    $convert.base64Decode(
        'CixHZXRPdXRib3VuZFdlYklkZW50aXR5RmVkZXJhdGlvbkluZm9SZXNwb25zZRItChBpc3N1ZX'
        'JpZGVudGlmaWVyGIaM0h4gASgJUhBpc3N1ZXJpZGVudGlmaWVyEi8KEWp3dHZlbmRpbmdlbmFi'
        'bGVkGPm5+UYgASgIUhFqd3R2ZW5kaW5nZW5hYmxlZA==');

@$core.Deprecated('Use getPolicyRequestDescriptor instead')
const GetPolicyRequest$json = {
  '1': 'GetPolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
  ],
};

/// Descriptor for `GetPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPolicyRequestDescriptor = $convert.base64Decode(
    'ChBHZXRQb2xpY3lSZXF1ZXN0EiAKCXBvbGljeWFybhjD0rrtASABKAlSCXBvbGljeWFybg==');

@$core.Deprecated('Use getPolicyResponseDescriptor instead')
const GetPolicyResponse$json = {
  '1': 'GetPolicyResponse',
  '2': [
    {
      '1': 'policy',
      '3': 471611296,
      '4': 1,
      '5': 11,
      '6': '.iam.Policy',
      '10': 'policy'
    },
  ],
};

/// Descriptor for `GetPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPolicyResponseDescriptor = $convert.base64Decode(
    'ChFHZXRQb2xpY3lSZXNwb25zZRInCgZwb2xpY3kYoO/w4AEgASgLMgsuaWFtLlBvbGljeVIGcG'
    '9saWN5');

@$core.Deprecated('Use getPolicyVersionRequestDescriptor instead')
const GetPolicyVersionRequest$json = {
  '1': 'GetPolicyVersionRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetPolicyVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPolicyVersionRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRQb2xpY3lWZXJzaW9uUmVxdWVzdBIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xpY3'
        'lhcm4SIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use getPolicyVersionResponseDescriptor instead')
const GetPolicyVersionResponse$json = {
  '1': 'GetPolicyVersionResponse',
  '2': [
    {
      '1': 'policyversion',
      '3': 266935938,
      '4': 1,
      '5': 11,
      '6': '.iam.PolicyVersion',
      '10': 'policyversion'
    },
  ],
};

/// Descriptor for `GetPolicyVersionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPolicyVersionResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXRQb2xpY3lWZXJzaW9uUmVzcG9uc2USOwoNcG9saWN5dmVyc2lvbhiCvaR/IAEoCzISLm'
        'lhbS5Qb2xpY3lWZXJzaW9uUg1wb2xpY3l2ZXJzaW9u');

@$core.Deprecated('Use getRolePolicyRequestDescriptor instead')
const GetRolePolicyRequest$json = {
  '1': 'GetRolePolicyRequest',
  '2': [
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `GetRolePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRolePolicyRequestDescriptor = $convert.base64Decode(
    'ChRHZXRSb2xlUG9saWN5UmVxdWVzdBIhCgpwb2xpY3luYW1lGL31h38gASgJUgpwb2xpY3luYW'
    '1lEh4KCHJvbGVuYW1lGLPzvMIBIAEoCVIIcm9sZW5hbWU=');

@$core.Deprecated('Use getRolePolicyResponseDescriptor instead')
const GetRolePolicyResponse$json = {
  '1': 'GetRolePolicyResponse',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `GetRolePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRolePolicyResponseDescriptor = $convert.base64Decode(
    'ChVHZXRSb2xlUG9saWN5UmVzcG9uc2USKQoOcG9saWN5ZG9jdW1lbnQYy67BcSABKAlSDnBvbG'
    'ljeWRvY3VtZW50EiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWUSHgoIcm9sZW5h'
    'bWUYs/O8wgEgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use getRoleRequestDescriptor instead')
const GetRoleRequest$json = {
  '1': 'GetRoleRequest',
  '2': [
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `GetRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRoleRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRSb2xlUmVxdWVzdBIeCghyb2xlbmFtZRiz87zCASABKAlSCHJvbGVuYW1l');

@$core.Deprecated('Use getRoleResponseDescriptor instead')
const GetRoleResponse$json = {
  '1': 'GetRoleResponse',
  '2': [
    {
      '1': 'role',
      '3': 271285818,
      '4': 1,
      '5': 11,
      '6': '.iam.Role',
      '10': 'role'
    },
  ],
};

/// Descriptor for `GetRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRoleResponseDescriptor = $convert.base64Decode(
    'Cg9HZXRSb2xlUmVzcG9uc2USIQoEcm9sZRi6/K2BASABKAsyCS5pYW0uUm9sZVIEcm9sZQ==');

@$core.Deprecated('Use getSAMLProviderRequestDescriptor instead')
const GetSAMLProviderRequest$json = {
  '1': 'GetSAMLProviderRequest',
  '2': [
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
  ],
};

/// Descriptor for `GetSAMLProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSAMLProviderRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRTQU1MUHJvdmlkZXJSZXF1ZXN0EiwKD3NhbWxwcm92aWRlcmFybhilzaiMASABKAlSD3'
        'NhbWxwcm92aWRlcmFybg==');

@$core.Deprecated('Use getSAMLProviderResponseDescriptor instead')
const GetSAMLProviderResponse$json = {
  '1': 'GetSAMLProviderResponse',
  '2': [
    {
      '1': 'assertionencryptionmode',
      '3': 474560298,
      '4': 1,
      '5': 14,
      '6': '.iam.assertionEncryptionModeType',
      '10': 'assertionencryptionmode'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'privatekeylist',
      '3': 145689700,
      '4': 3,
      '5': 11,
      '6': '.iam.SAMLPrivateKey',
      '10': 'privatekeylist'
    },
    {
      '1': 'samlmetadatadocument',
      '3': 282432645,
      '4': 1,
      '5': 9,
      '10': 'samlmetadatadocument'
    },
    {
      '1': 'samlprovideruuid',
      '3': 241210667,
      '4': 1,
      '5': 9,
      '10': 'samlprovideruuid'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {'1': 'validuntil', '3': 366644404, '4': 1, '5': 9, '10': 'validuntil'},
  ],
};

/// Descriptor for `GetSAMLProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSAMLProviderResponseDescriptor = $convert.base64Decode(
    'ChdHZXRTQU1MUHJvdmlkZXJSZXNwb25zZRJeChdhc3NlcnRpb25lbmNyeXB0aW9ubW9kZRiq7q'
    'TiASABKA4yIC5pYW0uYXNzZXJ0aW9uRW5jcnlwdGlvbk1vZGVUeXBlUhdhc3NlcnRpb25lbmNy'
    'eXB0aW9ubW9kZRIhCgpjcmVhdGVkYXRlGJK5/BEgASgJUgpjcmVhdGVkYXRlEj4KDnByaXZhdG'
    'VrZXlsaXN0GOSYvEUgAygLMhMuaWFtLlNBTUxQcml2YXRlS2V5Ug5wcml2YXRla2V5bGlzdBI2'
    'ChRzYW1sbWV0YWRhdGFkb2N1bWVudBiFqdaGASABKAlSFHNhbWxtZXRhZGF0YWRvY3VtZW50Ei'
    '0KEHNhbWxwcm92aWRlcnV1aWQYq6qCcyABKAlSEHNhbWxwcm92aWRlcnV1aWQSIAoEdGFncxjB'
    'wfa1ASADKAsyCC5pYW0uVGFnUgR0YWdzEiIKCnZhbGlkdW50aWwYtJnqrgEgASgJUgp2YWxpZH'
    'VudGls');

@$core.Deprecated('Use getSSHPublicKeyRequestDescriptor instead')
const GetSSHPublicKeyRequest$json = {
  '1': 'GetSSHPublicKeyRequest',
  '2': [
    {
      '1': 'encoding',
      '3': 237643025,
      '4': 1,
      '5': 14,
      '6': '.iam.encodingType',
      '10': 'encoding'
    },
    {
      '1': 'sshpublickeyid',
      '3': 154499415,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeyid'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetSSHPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSSHPublicKeyRequestDescriptor = $convert.base64Decode(
    'ChZHZXRTU0hQdWJsaWNLZXlSZXF1ZXN0EjAKCGVuY29kaW5nGJHKqHEgASgOMhEuaWFtLmVuY2'
    '9kaW5nVHlwZVIIZW5jb2RpbmcSKQoOc3NocHVibGlja2V5aWQY1/LVSSABKAlSDnNzaHB1Ymxp'
    'Y2tleWlkEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use getSSHPublicKeyResponseDescriptor instead')
const GetSSHPublicKeyResponse$json = {
  '1': 'GetSSHPublicKeyResponse',
  '2': [
    {
      '1': 'sshpublickey',
      '3': 520385596,
      '4': 1,
      '5': 11,
      '6': '.iam.SSHPublicKey',
      '10': 'sshpublickey'
    },
  ],
};

/// Descriptor for `GetSSHPublicKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSSHPublicKeyResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRTU0hQdWJsaWNLZXlSZXNwb25zZRI5Cgxzc2hwdWJsaWNrZXkYvOiR+AEgASgLMhEuaW'
        'FtLlNTSFB1YmxpY0tleVIMc3NocHVibGlja2V5');

@$core.Deprecated('Use getServerCertificateRequestDescriptor instead')
const GetServerCertificateRequest$json = {
  '1': 'GetServerCertificateRequest',
  '2': [
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
  ],
};

/// Descriptor for `GetServerCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServerCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRTZXJ2ZXJDZXJ0aWZpY2F0ZVJlcXVlc3QSNwoVc2VydmVyY2VydGlmaWNhdGVuYW1lGP'
        '+Amx4gASgJUhVzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbWU=');

@$core.Deprecated('Use getServerCertificateResponseDescriptor instead')
const GetServerCertificateResponse$json = {
  '1': 'GetServerCertificateResponse',
  '2': [
    {
      '1': 'servercertificate',
      '3': 68775774,
      '4': 1,
      '5': 11,
      '6': '.iam.ServerCertificate',
      '10': 'servercertificate'
    },
  ],
};

/// Descriptor for `GetServerCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServerCertificateResponseDescriptor =
    $convert.base64Decode(
        'ChxHZXRTZXJ2ZXJDZXJ0aWZpY2F0ZVJlc3BvbnNlEkcKEXNlcnZlcmNlcnRpZmljYXRlGN7e5S'
        'AgASgLMhYuaWFtLlNlcnZlckNlcnRpZmljYXRlUhFzZXJ2ZXJjZXJ0aWZpY2F0ZQ==');

@$core.Deprecated('Use getServiceLastAccessedDetailsRequestDescriptor instead')
const GetServiceLastAccessedDetailsRequest$json = {
  '1': 'GetServiceLastAccessedDetailsRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `GetServiceLastAccessedDetailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServiceLastAccessedDetailsRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRTZXJ2aWNlTGFzdEFjY2Vzc2VkRGV0YWlsc1JlcXVlc3QSFwoFam9iaWQY0tTdMyABKA'
        'lSBWpvYmlkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEo'
        'BVIIbWF4aXRlbXM=');

@$core.Deprecated('Use getServiceLastAccessedDetailsResponseDescriptor instead')
const GetServiceLastAccessedDetailsResponse$json = {
  '1': 'GetServiceLastAccessedDetailsResponse',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.iam.ErrorDetails',
      '10': 'error'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'jobcompletiondate',
      '3': 517991255,
      '4': 1,
      '5': 9,
      '10': 'jobcompletiondate'
    },
    {
      '1': 'jobcreationdate',
      '3': 201415034,
      '4': 1,
      '5': 9,
      '10': 'jobcreationdate'
    },
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.iam.jobStatusType',
      '10': 'jobstatus'
    },
    {
      '1': 'jobtype',
      '3': 279889397,
      '4': 1,
      '5': 14,
      '6': '.iam.AccessAdvisorUsageGranularityType',
      '10': 'jobtype'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'serviceslastaccessed',
      '3': 264776823,
      '4': 3,
      '5': 11,
      '6': '.iam.ServiceLastAccessed',
      '10': 'serviceslastaccessed'
    },
  ],
};

/// Descriptor for `GetServiceLastAccessedDetailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServiceLastAccessedDetailsResponseDescriptor = $convert.base64Decode(
    'CiVHZXRTZXJ2aWNlTGFzdEFjY2Vzc2VkRGV0YWlsc1Jlc3BvbnNlEisKBWVycm9yGPK5tpwBIA'
    'EoCzIRLmlhbS5FcnJvckRldGFpbHNSBWVycm9yEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtp'
    'c3RydW5jYXRlZBIwChFqb2Jjb21wbGV0aW9uZGF0ZRjX1v/2ASABKAlSEWpvYmNvbXBsZXRpb2'
    '5kYXRlEisKD2pvYmNyZWF0aW9uZGF0ZRj6soVgIAEoCVIPam9iY3JlYXRpb25kYXRlEjMKCWpv'
    'YnN0YXR1cxjHnPszIAEoDjISLmlhbS5qb2JTdGF0dXNUeXBlUglqb2JzdGF0dXMSRAoHam9idH'
    'lwZRj1i7uFASABKA4yJi5pYW0uQWNjZXNzQWR2aXNvclVzYWdlR3JhbnVsYXJpdHlUeXBlUgdq'
    'b2J0eXBlEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEk8KFHNlcnZpY2VzbGFzdGFjY2Vzc2'
    'VkGPfYoH4gAygLMhguaWFtLlNlcnZpY2VMYXN0QWNjZXNzZWRSFHNlcnZpY2VzbGFzdGFjY2Vz'
    'c2Vk');

@$core.Deprecated(
    'Use getServiceLastAccessedDetailsWithEntitiesRequestDescriptor instead')
const GetServiceLastAccessedDetailsWithEntitiesRequest$json = {
  '1': 'GetServiceLastAccessedDetailsWithEntitiesRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'servicenamespace',
      '3': 432654764,
      '4': 1,
      '5': 9,
      '10': 'servicenamespace'
    },
  ],
};

/// Descriptor for `GetServiceLastAccessedDetailsWithEntitiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getServiceLastAccessedDetailsWithEntitiesRequestDescriptor =
    $convert.base64Decode(
        'CjBHZXRTZXJ2aWNlTGFzdEFjY2Vzc2VkRGV0YWlsc1dpdGhFbnRpdGllc1JlcXVlc3QSFwoFam'
        '9iaWQY0tTdMyABKAlSBWpvYmlkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0'
        'ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSLgoQc2VydmljZW5hbWVzcGFjZRisk6fOASABKAlSEH'
        'NlcnZpY2VuYW1lc3BhY2U=');

@$core.Deprecated(
    'Use getServiceLastAccessedDetailsWithEntitiesResponseDescriptor instead')
const GetServiceLastAccessedDetailsWithEntitiesResponse$json = {
  '1': 'GetServiceLastAccessedDetailsWithEntitiesResponse',
  '2': [
    {
      '1': 'entitydetailslist',
      '3': 371739953,
      '4': 3,
      '5': 11,
      '6': '.iam.EntityDetails',
      '10': 'entitydetailslist'
    },
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.iam.ErrorDetails',
      '10': 'error'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'jobcompletiondate',
      '3': 517991255,
      '4': 1,
      '5': 9,
      '10': 'jobcompletiondate'
    },
    {
      '1': 'jobcreationdate',
      '3': 201415034,
      '4': 1,
      '5': 9,
      '10': 'jobcreationdate'
    },
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.iam.jobStatusType',
      '10': 'jobstatus'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `GetServiceLastAccessedDetailsWithEntitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getServiceLastAccessedDetailsWithEntitiesResponseDescriptor =
    $convert.base64Decode(
        'CjFHZXRTZXJ2aWNlTGFzdEFjY2Vzc2VkRGV0YWlsc1dpdGhFbnRpdGllc1Jlc3BvbnNlEkQKEW'
        'VudGl0eWRldGFpbHNsaXN0GLGaobEBIAMoCzISLmlhbS5FbnRpdHlEZXRhaWxzUhFlbnRpdHlk'
        'ZXRhaWxzbGlzdBIrCgVlcnJvchjyubacASABKAsyES5pYW0uRXJyb3JEZXRhaWxzUgVlcnJvch'
        'IjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSMAoRam9iY29tcGxldGlvbmRh'
        'dGUY19b/9gEgASgJUhFqb2Jjb21wbGV0aW9uZGF0ZRIrCg9qb2JjcmVhdGlvbmRhdGUY+rKFYC'
        'ABKAlSD2pvYmNyZWF0aW9uZGF0ZRIzCglqb2JzdGF0dXMYx5z7MyABKA4yEi5pYW0uam9iU3Rh'
        'dHVzVHlwZVIJam9ic3RhdHVzEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2Vy');

@$core.Deprecated(
    'Use getServiceLinkedRoleDeletionStatusRequestDescriptor instead')
const GetServiceLinkedRoleDeletionStatusRequest$json = {
  '1': 'GetServiceLinkedRoleDeletionStatusRequest',
  '2': [
    {
      '1': 'deletiontaskid',
      '3': 379712140,
      '4': 1,
      '5': 9,
      '10': 'deletiontaskid'
    },
  ],
};

/// Descriptor for `GetServiceLinkedRoleDeletionStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getServiceLinkedRoleDeletionStatusRequestDescriptor = $convert.base64Decode(
        'CilHZXRTZXJ2aWNlTGlua2VkUm9sZURlbGV0aW9uU3RhdHVzUmVxdWVzdBIqCg5kZWxldGlvbn'
        'Rhc2tpZBiM5Ye1ASABKAlSDmRlbGV0aW9udGFza2lk');

@$core.Deprecated(
    'Use getServiceLinkedRoleDeletionStatusResponseDescriptor instead')
const GetServiceLinkedRoleDeletionStatusResponse$json = {
  '1': 'GetServiceLinkedRoleDeletionStatusResponse',
  '2': [
    {
      '1': 'reason',
      '3': 20005178,
      '4': 1,
      '5': 11,
      '6': '.iam.DeletionTaskFailureReasonType',
      '10': 'reason'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.DeletionTaskStatusType',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetServiceLinkedRoleDeletionStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getServiceLinkedRoleDeletionStatusResponseDescriptor =
    $convert.base64Decode(
        'CipHZXRTZXJ2aWNlTGlua2VkUm9sZURlbGV0aW9uU3RhdHVzUmVzcG9uc2USPQoGcmVhc29uGL'
        'qCxQkgASgLMiIuaWFtLkRlbGV0aW9uVGFza0ZhaWx1cmVSZWFzb25UeXBlUgZyZWFzb24SNgoG'
        'c3RhdHVzGJDk+wIgASgOMhsuaWFtLkRlbGV0aW9uVGFza1N0YXR1c1R5cGVSBnN0YXR1cw==');

@$core.Deprecated('Use getUserPolicyRequestDescriptor instead')
const GetUserPolicyRequest$json = {
  '1': 'GetUserPolicyRequest',
  '2': [
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetUserPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserPolicyRequestDescriptor = $convert.base64Decode(
    'ChRHZXRVc2VyUG9saWN5UmVxdWVzdBIhCgpwb2xpY3luYW1lGL31h38gASgJUgpwb2xpY3luYW'
    '1lEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use getUserPolicyResponseDescriptor instead')
const GetUserPolicyResponse$json = {
  '1': 'GetUserPolicyResponse',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetUserPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserPolicyResponseDescriptor = $convert.base64Decode(
    'ChVHZXRVc2VyUG9saWN5UmVzcG9uc2USKQoOcG9saWN5ZG9jdW1lbnQYy67BcSABKAlSDnBvbG'
    'ljeWRvY3VtZW50EiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWUSHgoIdXNlcm5h'
    'bWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use getUserRequestDescriptor instead')
const GetUserRequest$json = {
  '1': 'GetUserRequest',
  '2': [
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRVc2VyUmVxdWVzdBIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use getUserResponseDescriptor instead')
const GetUserResponse$json = {
  '1': 'GetUserResponse',
  '2': [
    {
      '1': 'user',
      '3': 10894867,
      '4': 1,
      '5': 11,
      '6': '.iam.User',
      '10': 'user'
    },
  ],
};

/// Descriptor for `GetUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserResponseDescriptor = $convert.base64Decode(
    'Cg9HZXRVc2VyUmVzcG9uc2USIAoEdXNlchiT/JgFIAEoCzIJLmlhbS5Vc2VyUgR1c2Vy');

@$core.Deprecated('Use groupDescriptor instead')
const Group$json = {
  '1': 'Group',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {'1': 'groupid', '3': 73973250, '4': 1, '5': 9, '10': 'groupid'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `Group`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List groupDescriptor = $convert.base64Decode(
    'CgVHcm91cBIUCgNhcm4YnZvtvwEgASgJUgNhcm4SIQoKY3JlYXRlZGF0ZRiSufwRIAEoCVIKY3'
    'JlYXRlZGF0ZRIbCgdncm91cGlkGIL8oiMgASgJUgdncm91cGlkEiAKCWdyb3VwbmFtZRjIyqCq'
    'ASABKAlSCWdyb3VwbmFtZRIVCgRwYXRoGNfIm1sgASgJUgRwYXRo');

@$core.Deprecated('Use groupDetailDescriptor instead')
const GroupDetail$json = {
  '1': 'GroupDetail',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'attachedmanagedpolicies',
      '3': 155564765,
      '4': 3,
      '5': 11,
      '6': '.iam.AttachedPolicy',
      '10': 'attachedmanagedpolicies'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {'1': 'groupid', '3': 73973250, '4': 1, '5': 9, '10': 'groupid'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {
      '1': 'grouppolicylist',
      '3': 442500627,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyDetail',
      '10': 'grouppolicylist'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `GroupDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List groupDetailDescriptor = $convert.base64Decode(
    'CgtHcm91cERldGFpbBIUCgNhcm4YnZvtvwEgASgJUgNhcm4SUAoXYXR0YWNoZWRtYW5hZ2VkcG'
    '9saWNpZXMY3fWWSiADKAsyEy5pYW0uQXR0YWNoZWRQb2xpY3lSF2F0dGFjaGVkbWFuYWdlZHBv'
    'bGljaWVzEiEKCmNyZWF0ZWRhdGUYkrn8ESABKAlSCmNyZWF0ZWRhdGUSGwoHZ3JvdXBpZBiC/K'
    'IjIAEoCVIHZ3JvdXBpZBIgCglncm91cG5hbWUYyMqgqgEgASgJUglncm91cG5hbWUSPwoPZ3Jv'
    'dXBwb2xpY3lsaXN0GJOMgNMBIAMoCzIRLmlhbS5Qb2xpY3lEZXRhaWxSD2dyb3VwcG9saWN5bG'
    'lzdBIVCgRwYXRoGNfIm1sgASgJUgRwYXRo');

@$core.Deprecated('Use instanceProfileDescriptor instead')
const InstanceProfile$json = {
  '1': 'InstanceProfile',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'instanceprofileid',
      '3': 97287331,
      '4': 1,
      '5': 9,
      '10': 'instanceprofileid'
    },
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'roles',
      '3': 511168127,
      '4': 3,
      '5': 11,
      '6': '.iam.Role',
      '10': 'roles'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `InstanceProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceProfileDescriptor = $convert.base64Decode(
    'Cg9JbnN0YW5jZVByb2ZpbGUSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEiEKCmNyZWF0ZWRhdGUYkr'
    'n8ESABKAlSCmNyZWF0ZWRhdGUSLwoRaW5zdGFuY2Vwcm9maWxlaWQYo/mxLiABKAlSEWluc3Rh'
    'bmNlcHJvZmlsZWlkEjQKE2luc3RhbmNlcHJvZmlsZW5hbWUY7c682gEgASgJUhNpbnN0YW5jZX'
    'Byb2ZpbGVuYW1lEhUKBHBhdGgY18ibWyABKAlSBHBhdGgSIwoFcm9sZXMY/5zf8wEgAygLMgku'
    'aWFtLlJvbGVSBXJvbGVzEiAKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use invalidAuthenticationCodeExceptionDescriptor instead')
const InvalidAuthenticationCodeException$json = {
  '1': 'InvalidAuthenticationCodeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAuthenticationCodeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAuthenticationCodeExceptionDescriptor =
    $convert.base64Decode(
        'CiJJbnZhbGlkQXV0aGVudGljYXRpb25Db2RlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidCertificateExceptionDescriptor instead')
const InvalidCertificateException$json = {
  '1': 'InvalidCertificateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidCertificateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidCertificateExceptionDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkQ2VydGlmaWNhdGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use invalidInputExceptionDescriptor instead')
const InvalidInputException$json = {
  '1': 'InvalidInputException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInputException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInputExceptionDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkSW5wdXRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidPublicKeyExceptionDescriptor instead')
const InvalidPublicKeyException$json = {
  '1': 'InvalidPublicKeyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidPublicKeyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPublicKeyExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUHVibGljS2V5RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidUserTypeExceptionDescriptor instead')
const InvalidUserTypeException$json = {
  '1': 'InvalidUserTypeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidUserTypeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidUserTypeExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkVXNlclR5cGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use keyPairMismatchExceptionDescriptor instead')
const KeyPairMismatchException$json = {
  '1': 'KeyPairMismatchException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeyPairMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyPairMismatchExceptionDescriptor =
    $convert.base64Decode(
        'ChhLZXlQYWlyTWlzbWF0Y2hFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

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

@$core.Deprecated('Use listAccessKeysRequestDescriptor instead')
const ListAccessKeysRequest$json = {
  '1': 'ListAccessKeysRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListAccessKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAccessKeysRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0QWNjZXNzS2V5c1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbW'
    'F4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJu'
    'YW1l');

@$core.Deprecated('Use listAccessKeysResponseDescriptor instead')
const ListAccessKeysResponse$json = {
  '1': 'ListAccessKeysResponse',
  '2': [
    {
      '1': 'accesskeymetadata',
      '3': 74401512,
      '4': 3,
      '5': 11,
      '6': '.iam.AccessKeyMetadata',
      '10': 'accesskeymetadata'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListAccessKeysResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAccessKeysResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0QWNjZXNzS2V5c1Jlc3BvbnNlEkcKEWFjY2Vzc2tleW1ldGFkYXRhGOiNvSMgAygLMh'
    'YuaWFtLkFjY2Vzc0tleU1ldGFkYXRhUhFhY2Nlc3NrZXltZXRhZGF0YRIjCgtpc3RydW5jYXRl'
    'ZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXI=');

@$core.Deprecated('Use listAccountAliasesRequestDescriptor instead')
const ListAccountAliasesRequest$json = {
  '1': 'ListAccountAliasesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListAccountAliasesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAccountAliasesRequestDescriptor =
    $convert.base64Decode(
        'ChlMaXN0QWNjb3VudEFsaWFzZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh'
        '4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listAccountAliasesResponseDescriptor instead')
const ListAccountAliasesResponse$json = {
  '1': 'ListAccountAliasesResponse',
  '2': [
    {
      '1': 'accountaliases',
      '3': 476963149,
      '4': 3,
      '5': 9,
      '10': 'accountaliases'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListAccountAliasesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAccountAliasesResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0QWNjb3VudEFsaWFzZXNSZXNwb25zZRIqCg5hY2NvdW50YWxpYXNlcxjNwrfjASADKA'
        'lSDmFjY291bnRhbGlhc2VzEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIZ'
        'CgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listAttachedGroupPoliciesRequestDescriptor instead')
const ListAttachedGroupPoliciesRequest$json = {
  '1': 'ListAttachedGroupPoliciesRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
  ],
};

/// Descriptor for `ListAttachedGroupPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAttachedGroupPoliciesRequestDescriptor =
    $convert.base64Decode(
        'CiBMaXN0QXR0YWNoZWRHcm91cFBvbGljaWVzUmVxdWVzdBIgCglncm91cG5hbWUYyMqgqgEgAS'
        'gJUglncm91cG5hbWUSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba'
        '8QEgASgFUghtYXhpdGVtcxIiCgpwYXRocHJlZml4GMu35r0BIAEoCVIKcGF0aHByZWZpeA==');

@$core.Deprecated('Use listAttachedGroupPoliciesResponseDescriptor instead')
const ListAttachedGroupPoliciesResponse$json = {
  '1': 'ListAttachedGroupPoliciesResponse',
  '2': [
    {
      '1': 'attachedpolicies',
      '3': 60262760,
      '4': 3,
      '5': 11,
      '6': '.iam.AttachedPolicy',
      '10': 'attachedpolicies'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListAttachedGroupPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAttachedGroupPoliciesResponseDescriptor =
    $convert.base64Decode(
        'CiFMaXN0QXR0YWNoZWRHcm91cFBvbGljaWVzUmVzcG9uc2USQgoQYXR0YWNoZWRwb2xpY2llcx'
        'jokt4cIAMoCzITLmlhbS5BdHRhY2hlZFBvbGljeVIQYXR0YWNoZWRwb2xpY2llcxIjCgtpc3Ry'
        'dW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZX'
        'I=');

@$core.Deprecated('Use listAttachedRolePoliciesRequestDescriptor instead')
const ListAttachedRolePoliciesRequest$json = {
  '1': 'ListAttachedRolePoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `ListAttachedRolePoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAttachedRolePoliciesRequestDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0QXR0YWNoZWRSb2xlUG9saWNpZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbW'
        'Fya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIgoKcGF0aHByZWZpeBjLt+a9'
        'ASABKAlSCnBhdGhwcmVmaXgSHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use listAttachedRolePoliciesResponseDescriptor instead')
const ListAttachedRolePoliciesResponse$json = {
  '1': 'ListAttachedRolePoliciesResponse',
  '2': [
    {
      '1': 'attachedpolicies',
      '3': 60262760,
      '4': 3,
      '5': 11,
      '6': '.iam.AttachedPolicy',
      '10': 'attachedpolicies'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListAttachedRolePoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAttachedRolePoliciesResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0QXR0YWNoZWRSb2xlUG9saWNpZXNSZXNwb25zZRJCChBhdHRhY2hlZHBvbGljaWVzGO'
        'iS3hwgAygLMhMuaWFtLkF0dGFjaGVkUG9saWN5UhBhdHRhY2hlZHBvbGljaWVzEiMKC2lzdHJ1'
        'bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlcg'
        '==');

@$core.Deprecated('Use listAttachedUserPoliciesRequestDescriptor instead')
const ListAttachedUserPoliciesRequest$json = {
  '1': 'ListAttachedUserPoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListAttachedUserPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAttachedUserPoliciesRequestDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0QXR0YWNoZWRVc2VyUG9saWNpZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbW'
        'Fya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIgoKcGF0aHByZWZpeBjLt+a9'
        'ASABKAlSCnBhdGhwcmVmaXgSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use listAttachedUserPoliciesResponseDescriptor instead')
const ListAttachedUserPoliciesResponse$json = {
  '1': 'ListAttachedUserPoliciesResponse',
  '2': [
    {
      '1': 'attachedpolicies',
      '3': 60262760,
      '4': 3,
      '5': 11,
      '6': '.iam.AttachedPolicy',
      '10': 'attachedpolicies'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListAttachedUserPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAttachedUserPoliciesResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0QXR0YWNoZWRVc2VyUG9saWNpZXNSZXNwb25zZRJCChBhdHRhY2hlZHBvbGljaWVzGO'
        'iS3hwgAygLMhMuaWFtLkF0dGFjaGVkUG9saWN5UhBhdHRhY2hlZHBvbGljaWVzEiMKC2lzdHJ1'
        'bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlcg'
        '==');

@$core.Deprecated('Use listDelegationRequestsRequestDescriptor instead')
const ListDelegationRequestsRequest$json = {
  '1': 'ListDelegationRequestsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'ownerid', '3': 375630298, '4': 1, '5': 9, '10': 'ownerid'},
  ],
};

/// Descriptor for `ListDelegationRequestsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDelegationRequestsRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0RGVsZWdhdGlvblJlcXVlc3RzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcm'
        'tlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEhwKB293bmVyaWQY2tOOswEgASgJ'
        'Ugdvd25lcmlk');

@$core.Deprecated('Use listDelegationRequestsResponseDescriptor instead')
const ListDelegationRequestsResponse$json = {
  '1': 'ListDelegationRequestsResponse',
  '2': [
    {
      '1': 'delegationrequests',
      '3': 385003146,
      '4': 3,
      '5': 11,
      '6': '.iam.DelegationRequest',
      '10': 'delegationrequests'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'istruncated', '3': 459873914, '4': 1, '5': 8, '10': 'istruncated'},
  ],
};

/// Descriptor for `ListDelegationRequestsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDelegationRequestsResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0RGVsZWdhdGlvblJlcXVlc3RzUmVzcG9uc2USSgoSZGVsZWdhdGlvbnJlcXVlc3RzGI'
        'rdyrcBIAMoCzIWLmlhbS5EZWxlZ2F0aW9uUmVxdWVzdFISZGVsZWdhdGlvbnJlcXVlc3RzEhkK'
        'Bm1hcmtlchi43c0qIAEoCVIGbWFya2VyEiQKC2lzdHJ1bmNhdGVkGPq8pNsBIAEoCFILaXN0cn'
        'VuY2F0ZWQ=');

@$core.Deprecated('Use listEntitiesForPolicyRequestDescriptor instead')
const ListEntitiesForPolicyRequest$json = {
  '1': 'ListEntitiesForPolicyRequest',
  '2': [
    {
      '1': 'entityfilter',
      '3': 325795861,
      '4': 1,
      '5': 14,
      '6': '.iam.EntityType',
      '10': 'entityfilter'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {
      '1': 'policyusagefilter',
      '3': 102535305,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyUsageType',
      '10': 'policyusagefilter'
    },
  ],
};

/// Descriptor for `ListEntitiesForPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEntitiesForPolicyRequestDescriptor = $convert.base64Decode(
    'ChxMaXN0RW50aXRpZXNGb3JQb2xpY3lSZXF1ZXN0EjcKDGVudGl0eWZpbHRlchiVgK2bASABKA'
    '4yDy5pYW0uRW50aXR5VHlwZVIMZW50aXR5ZmlsdGVyEhkKBm1hcmtlchi43c0qIAEoCVIGbWFy'
    'a2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIgoKcGF0aHByZWZpeBjLt+a9AS'
    'ABKAlSCnBhdGhwcmVmaXgSIAoJcG9saWN5YXJuGMPSuu0BIAEoCVIJcG9saWN5YXJuEkUKEXBv'
    'bGljeXVzYWdlZmlsdGVyGImh8jAgASgOMhQuaWFtLlBvbGljeVVzYWdlVHlwZVIRcG9saWN5dX'
    'NhZ2VmaWx0ZXI=');

@$core.Deprecated('Use listEntitiesForPolicyResponseDescriptor instead')
const ListEntitiesForPolicyResponse$json = {
  '1': 'ListEntitiesForPolicyResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'policygroups',
      '3': 187209232,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyGroup',
      '10': 'policygroups'
    },
    {
      '1': 'policyroles',
      '3': 475716745,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyRole',
      '10': 'policyroles'
    },
    {
      '1': 'policyusers',
      '3': 263543182,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyUser',
      '10': 'policyusers'
    },
  ],
};

/// Descriptor for `ListEntitiesForPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEntitiesForPolicyResponseDescriptor = $convert.base64Decode(
    'Ch1MaXN0RW50aXRpZXNGb3JQb2xpY3lSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCF'
    'ILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISNwoMcG9saWN5Z3JvdXBz'
    'GJCsolkgAygLMhAuaWFtLlBvbGljeUdyb3VwUgxwb2xpY3lncm91cHMSNQoLcG9saWN5cm9sZX'
    'MYibnr4gEgAygLMg8uaWFtLlBvbGljeVJvbGVSC3BvbGljeXJvbGVzEjQKC3BvbGljeXVzZXJz'
    'GI6z1X0gAygLMg8uaWFtLlBvbGljeVVzZXJSC3BvbGljeXVzZXJz');

@$core.Deprecated('Use listGroupPoliciesRequestDescriptor instead')
const ListGroupPoliciesRequest$json = {
  '1': 'ListGroupPoliciesRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListGroupPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupPoliciesRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0R3JvdXBQb2xpY2llc1JlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdX'
    'BuYW1lEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVII'
    'bWF4aXRlbXM=');

@$core.Deprecated('Use listGroupPoliciesResponseDescriptor instead')
const ListGroupPoliciesResponse$json = {
  '1': 'ListGroupPoliciesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'policynames', '3': 264098782, '4': 3, '5': 9, '10': 'policynames'},
  ],
};

/// Descriptor for `ListGroupPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupPoliciesResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0R3JvdXBQb2xpY2llc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3'
    'RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIjCgtwb2xpY3luYW1lcxjep/d9'
    'IAMoCVILcG9saWN5bmFtZXM=');

@$core.Deprecated('Use listGroupsForUserRequestDescriptor instead')
const ListGroupsForUserRequest$json = {
  '1': 'ListGroupsForUserRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListGroupsForUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupsForUserRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0R3JvdXBzRm9yVXNlclJlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHg'
    'oIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIeCgh1c2VybmFtZRj6wdThASABKAlSCHVz'
    'ZXJuYW1l');

@$core.Deprecated('Use listGroupsForUserResponseDescriptor instead')
const ListGroupsForUserResponse$json = {
  '1': 'ListGroupsForUserResponse',
  '2': [
    {
      '1': 'groups',
      '3': 360662414,
      '4': 3,
      '5': 11,
      '6': '.iam.Group',
      '10': 'groups'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListGroupsForUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupsForUserResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0R3JvdXBzRm9yVXNlclJlc3BvbnNlEiYKBmdyb3VwcxiOi/2rASADKAsyCi5pYW0uR3'
    'JvdXBSBmdyb3VwcxIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFy'
    'a2VyGLjdzSogASgJUgZtYXJrZXI=');

@$core.Deprecated('Use listGroupsRequestDescriptor instead')
const ListGroupsRequest$json = {
  '1': 'ListGroupsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
  ],
};

/// Descriptor for `ListGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupsRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0R3JvdXBzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdG'
    'VtcxiU1trxASABKAVSCG1heGl0ZW1zEiIKCnBhdGhwcmVmaXgYy7fmvQEgASgJUgpwYXRocHJl'
    'Zml4');

@$core.Deprecated('Use listGroupsResponseDescriptor instead')
const ListGroupsResponse$json = {
  '1': 'ListGroupsResponse',
  '2': [
    {
      '1': 'groups',
      '3': 360662414,
      '4': 3,
      '5': 11,
      '6': '.iam.Group',
      '10': 'groups'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListGroupsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0R3JvdXBzUmVzcG9uc2USJgoGZ3JvdXBzGI6L/asBIAMoCzIKLmlhbS5Hcm91cFIGZ3'
    'JvdXBzEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3N'
    'KiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listInstanceProfileTagsRequestDescriptor instead')
const ListInstanceProfileTagsRequest$json = {
  '1': 'ListInstanceProfileTagsRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListInstanceProfileTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInstanceProfileTagsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0SW5zdGFuY2VQcm9maWxlVGFnc1JlcXVlc3QSNAoTaW5zdGFuY2Vwcm9maWxlbmFtZR'
        'jtzrzaASABKAlSE2luc3RhbmNlcHJvZmlsZW5hbWUSGQoGbWFya2VyGLjdzSogASgJUgZtYXJr'
        'ZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listInstanceProfileTagsResponseDescriptor instead')
const ListInstanceProfileTagsResponse$json = {
  '1': 'ListInstanceProfileTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListInstanceProfileTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInstanceProfileTagsResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0SW5zdGFuY2VQcm9maWxlVGFnc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgAS'
        'gIUgtpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIgCgR0YWdzGMHB9rUB'
        'IAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listInstanceProfilesForRoleRequestDescriptor instead')
const ListInstanceProfilesForRoleRequest$json = {
  '1': 'ListInstanceProfilesForRoleRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `ListInstanceProfilesForRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInstanceProfilesForRoleRequestDescriptor =
    $convert.base64Decode(
        'CiJMaXN0SW5zdGFuY2VQcm9maWxlc0ZvclJvbGVSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCV'
        'IGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSHgoIcm9sZW5hbWUYs/O8'
        'wgEgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use listInstanceProfilesForRoleResponseDescriptor instead')
const ListInstanceProfilesForRoleResponse$json = {
  '1': 'ListInstanceProfilesForRoleResponse',
  '2': [
    {
      '1': 'instanceprofiles',
      '3': 111334261,
      '4': 3,
      '5': 11,
      '6': '.iam.InstanceProfile',
      '10': 'instanceprofiles'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListInstanceProfilesForRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInstanceProfilesForRoleResponseDescriptor =
    $convert.base64Decode(
        'CiNMaXN0SW5zdGFuY2VQcm9maWxlc0ZvclJvbGVSZXNwb25zZRJDChBpbnN0YW5jZXByb2ZpbG'
        'VzGPWmizUgAygLMhQuaWFtLkluc3RhbmNlUHJvZmlsZVIQaW5zdGFuY2Vwcm9maWxlcxIjCgtp'
        'c3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYX'
        'JrZXI=');

@$core.Deprecated('Use listInstanceProfilesRequestDescriptor instead')
const ListInstanceProfilesRequest$json = {
  '1': 'ListInstanceProfilesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
  ],
};

/// Descriptor for `ListInstanceProfilesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInstanceProfilesRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0SW5zdGFuY2VQcm9maWxlc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZX'
        'ISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpwYXRocHJlZml4GMu35r0BIAEo'
        'CVIKcGF0aHByZWZpeA==');

@$core.Deprecated('Use listInstanceProfilesResponseDescriptor instead')
const ListInstanceProfilesResponse$json = {
  '1': 'ListInstanceProfilesResponse',
  '2': [
    {
      '1': 'instanceprofiles',
      '3': 111334261,
      '4': 3,
      '5': 11,
      '6': '.iam.InstanceProfile',
      '10': 'instanceprofiles'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListInstanceProfilesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInstanceProfilesResponseDescriptor = $convert.base64Decode(
    'ChxMaXN0SW5zdGFuY2VQcm9maWxlc1Jlc3BvbnNlEkMKEGluc3RhbmNlcHJvZmlsZXMY9aaLNS'
    'ADKAsyFC5pYW0uSW5zdGFuY2VQcm9maWxlUhBpbnN0YW5jZXByb2ZpbGVzEiMKC2lzdHJ1bmNh'
    'dGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listMFADeviceTagsRequestDescriptor instead')
const ListMFADeviceTagsRequest$json = {
  '1': 'ListMFADeviceTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
  ],
};

/// Descriptor for `ListMFADeviceTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMFADeviceTagsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0TUZBRGV2aWNlVGFnc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHg'
    'oIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxImCgxzZXJpYWxudW1iZXIY5bq5xwEgASgJ'
    'UgxzZXJpYWxudW1iZXI=');

@$core.Deprecated('Use listMFADeviceTagsResponseDescriptor instead')
const ListMFADeviceTagsResponse$json = {
  '1': 'ListMFADeviceTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListMFADeviceTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMFADeviceTagsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0TUZBRGV2aWNlVGFnc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3'
    'RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIgCgR0YWdzGMHB9rUBIAMoCzII'
    'LmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listMFADevicesRequestDescriptor instead')
const ListMFADevicesRequest$json = {
  '1': 'ListMFADevicesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListMFADevicesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMFADevicesRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0TUZBRGV2aWNlc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbW'
    'F4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJu'
    'YW1l');

@$core.Deprecated('Use listMFADevicesResponseDescriptor instead')
const ListMFADevicesResponse$json = {
  '1': 'ListMFADevicesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'mfadevices',
      '3': 17178165,
      '4': 3,
      '5': 11,
      '6': '.iam.MFADevice',
      '10': 'mfadevices'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListMFADevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMFADevicesResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0TUZBRGV2aWNlc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW'
    '5jYXRlZBIxCgptZmFkZXZpY2VzGLW8mAggAygLMg4uaWFtLk1GQURldmljZVIKbWZhZGV2aWNl'
    'cxIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listOpenIDConnectProviderTagsRequestDescriptor instead')
const ListOpenIDConnectProviderTagsRequest$json = {
  '1': 'ListOpenIDConnectProviderTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
  ],
};

/// Descriptor for `ListOpenIDConnectProviderTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpenIDConnectProviderTagsRequestDescriptor =
    $convert.base64Decode(
        'CiRMaXN0T3BlbklEQ29ubmVjdFByb3ZpZGVyVGFnc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogAS'
        'gJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxI+ChhvcGVuaWRjb25u'
        'ZWN0cHJvdmlkZXJhcm4Y4/OP6QEgASgJUhhvcGVuaWRjb25uZWN0cHJvdmlkZXJhcm4=');

@$core.Deprecated('Use listOpenIDConnectProviderTagsResponseDescriptor instead')
const ListOpenIDConnectProviderTagsResponse$json = {
  '1': 'ListOpenIDConnectProviderTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListOpenIDConnectProviderTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpenIDConnectProviderTagsResponseDescriptor =
    $convert.base64Decode(
        'CiVMaXN0T3BlbklEQ29ubmVjdFByb3ZpZGVyVGFnc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGN'
        'qfuHMgASgIUgtpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIgCgR0YWdz'
        'GMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listOpenIDConnectProvidersRequestDescriptor instead')
const ListOpenIDConnectProvidersRequest$json = {
  '1': 'ListOpenIDConnectProvidersRequest',
};

/// Descriptor for `ListOpenIDConnectProvidersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpenIDConnectProvidersRequestDescriptor =
    $convert.base64Decode('CiFMaXN0T3BlbklEQ29ubmVjdFByb3ZpZGVyc1JlcXVlc3Q=');

@$core.Deprecated('Use listOpenIDConnectProvidersResponseDescriptor instead')
const ListOpenIDConnectProvidersResponse$json = {
  '1': 'ListOpenIDConnectProvidersResponse',
  '2': [
    {
      '1': 'openidconnectproviderlist',
      '3': 57978792,
      '4': 3,
      '5': 11,
      '6': '.iam.OpenIDConnectProviderListEntry',
      '10': 'openidconnectproviderlist'
    },
  ],
};

/// Descriptor for `ListOpenIDConnectProvidersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpenIDConnectProvidersResponseDescriptor =
    $convert.base64Decode(
        'CiJMaXN0T3BlbklEQ29ubmVjdFByb3ZpZGVyc1Jlc3BvbnNlEmQKGW9wZW5pZGNvbm5lY3Rwcm'
        '92aWRlcmxpc3QYqN/SGyADKAsyIy5pYW0uT3BlbklEQ29ubmVjdFByb3ZpZGVyTGlzdEVudHJ5'
        'UhlvcGVuaWRjb25uZWN0cHJvdmlkZXJsaXN0');

@$core.Deprecated('Use listOrganizationsFeaturesRequestDescriptor instead')
const ListOrganizationsFeaturesRequest$json = {
  '1': 'ListOrganizationsFeaturesRequest',
};

/// Descriptor for `ListOrganizationsFeaturesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOrganizationsFeaturesRequestDescriptor =
    $convert.base64Decode('CiBMaXN0T3JnYW5pemF0aW9uc0ZlYXR1cmVzUmVxdWVzdA==');

@$core.Deprecated('Use listOrganizationsFeaturesResponseDescriptor instead')
const ListOrganizationsFeaturesResponse$json = {
  '1': 'ListOrganizationsFeaturesResponse',
  '2': [
    {
      '1': 'enabledfeatures',
      '3': 211387242,
      '4': 3,
      '5': 14,
      '6': '.iam.FeatureType',
      '10': 'enabledfeatures'
    },
    {
      '1': 'organizationid',
      '3': 311006120,
      '4': 1,
      '5': 9,
      '10': 'organizationid'
    },
  ],
};

/// Descriptor for `ListOrganizationsFeaturesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOrganizationsFeaturesResponseDescriptor =
    $convert.base64Decode(
        'CiFMaXN0T3JnYW5pemF0aW9uc0ZlYXR1cmVzUmVzcG9uc2USPQoPZW5hYmxlZGZlYXR1cmVzGO'
        'qG5mQgAygOMhAuaWFtLkZlYXR1cmVUeXBlUg9lbmFibGVkZmVhdHVyZXMSKgoOb3JnYW5pemF0'
        'aW9uaWQYqKemlAEgASgJUg5vcmdhbml6YXRpb25pZA==');

@$core
    .Deprecated('Use listPoliciesGrantingServiceAccessEntryDescriptor instead')
const ListPoliciesGrantingServiceAccessEntry$json = {
  '1': 'ListPoliciesGrantingServiceAccessEntry',
  '2': [
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyGrantingServiceAccess',
      '10': 'policies'
    },
    {
      '1': 'servicenamespace',
      '3': 432654764,
      '4': 1,
      '5': 9,
      '10': 'servicenamespace'
    },
  ],
};

/// Descriptor for `ListPoliciesGrantingServiceAccessEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPoliciesGrantingServiceAccessEntryDescriptor =
    $convert.base64Decode(
        'CiZMaXN0UG9saWNpZXNHcmFudGluZ1NlcnZpY2VBY2Nlc3NFbnRyeRI/Cghwb2xpY2llcxiYrI'
        'oTIAMoCzIgLmlhbS5Qb2xpY3lHcmFudGluZ1NlcnZpY2VBY2Nlc3NSCHBvbGljaWVzEi4KEHNl'
        'cnZpY2VuYW1lc3BhY2UYrJOnzgEgASgJUhBzZXJ2aWNlbmFtZXNwYWNl');

@$core.Deprecated(
    'Use listPoliciesGrantingServiceAccessRequestDescriptor instead')
const ListPoliciesGrantingServiceAccessRequest$json = {
  '1': 'ListPoliciesGrantingServiceAccessRequest',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'servicenamespaces',
      '3': 279494409,
      '4': 3,
      '5': 9,
      '10': 'servicenamespaces'
    },
  ],
};

/// Descriptor for `ListPoliciesGrantingServiceAccessRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPoliciesGrantingServiceAccessRequestDescriptor =
    $convert.base64Decode(
        'CihMaXN0UG9saWNpZXNHcmFudGluZ1NlcnZpY2VBY2Nlc3NSZXF1ZXN0EhQKA2Fybhidm+2/AS'
        'ABKAlSA2FybhIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIwChFzZXJ2aWNlbmFtZXNwYWNl'
        'cxiJ/qKFASADKAlSEXNlcnZpY2VuYW1lc3BhY2Vz');

@$core.Deprecated(
    'Use listPoliciesGrantingServiceAccessResponseDescriptor instead')
const ListPoliciesGrantingServiceAccessResponse$json = {
  '1': 'ListPoliciesGrantingServiceAccessResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'policiesgrantingserviceaccess',
      '3': 90939387,
      '4': 3,
      '5': 11,
      '6': '.iam.ListPoliciesGrantingServiceAccessEntry',
      '10': 'policiesgrantingserviceaccess'
    },
  ],
};

/// Descriptor for `ListPoliciesGrantingServiceAccessResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listPoliciesGrantingServiceAccessResponseDescriptor = $convert.base64Decode(
        'CilMaXN0UG9saWNpZXNHcmFudGluZ1NlcnZpY2VBY2Nlc3NSZXNwb25zZRIjCgtpc3RydW5jYX'
        'RlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISdAod'
        'cG9saWNpZXNncmFudGluZ3NlcnZpY2VhY2Nlc3MY+7+uKyADKAsyKy5pYW0uTGlzdFBvbGljaW'
        'VzR3JhbnRpbmdTZXJ2aWNlQWNjZXNzRW50cnlSHXBvbGljaWVzZ3JhbnRpbmdzZXJ2aWNlYWNj'
        'ZXNz');

@$core.Deprecated('Use listPoliciesRequestDescriptor instead')
const ListPoliciesRequest$json = {
  '1': 'ListPoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'onlyattached', '3': 264348868, '4': 1, '5': 8, '10': 'onlyattached'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
    {
      '1': 'policyusagefilter',
      '3': 102535305,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyUsageType',
      '10': 'policyusagefilter'
    },
    {
      '1': 'scope',
      '3': 65430924,
      '4': 1,
      '5': 14,
      '6': '.iam.policyScopeType',
      '10': 'scope'
    },
  ],
};

/// Descriptor for `ListPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPoliciesRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0UG9saWNpZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heG'
    'l0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSJQoMb25seWF0dGFjaGVkGMTJhn4gASgIUgxvbmx5'
    'YXR0YWNoZWQSIgoKcGF0aHByZWZpeBjLt+a9ASABKAlSCnBhdGhwcmVmaXgSRQoRcG9saWN5dX'
    'NhZ2VmaWx0ZXIYiaHyMCABKA4yFC5pYW0uUG9saWN5VXNhZ2VUeXBlUhFwb2xpY3l1c2FnZWZp'
    'bHRlchItCgVzY29wZRiMy5kfIAEoDjIULmlhbS5wb2xpY3lTY29wZVR5cGVSBXNjb3Bl');

@$core.Deprecated('Use listPoliciesResponseDescriptor instead')
const ListPoliciesResponse$json = {
  '1': 'ListPoliciesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.iam.Policy',
      '10': 'policies'
    },
  ],
};

/// Descriptor for `ListPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPoliciesResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0UG9saWNpZXNSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2'
    'F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISKgoIcG9saWNpZXMYmKyKEyADKAsyCy5p'
    'YW0uUG9saWN5Ughwb2xpY2llcw==');

@$core.Deprecated('Use listPolicyTagsRequestDescriptor instead')
const ListPolicyTagsRequest$json = {
  '1': 'ListPolicyTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
  ],
};

/// Descriptor for `ListPolicyTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPolicyTagsRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0UG9saWN5VGFnc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbW'
    'F4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIgCglwb2xpY3lhcm4Yw9K67QEgASgJUglwb2xp'
    'Y3lhcm4=');

@$core.Deprecated('Use listPolicyTagsResponseDescriptor instead')
const ListPolicyTagsResponse$json = {
  '1': 'ListPolicyTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListPolicyTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPolicyTagsResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0UG9saWN5VGFnc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW'
    '5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIgCgR0YWdzGMHB9rUBIAMoCzIILmlh'
    'bS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listPolicyVersionsRequestDescriptor instead')
const ListPolicyVersionsRequest$json = {
  '1': 'ListPolicyVersionsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
  ],
};

/// Descriptor for `ListPolicyVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPolicyVersionsRequestDescriptor = $convert.base64Decode(
    'ChlMaXN0UG9saWN5VmVyc2lvbnNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh'
    '4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIAoJcG9saWN5YXJuGMPSuu0BIAEoCVIJ'
    'cG9saWN5YXJu');

@$core.Deprecated('Use listPolicyVersionsResponseDescriptor instead')
const ListPolicyVersionsResponse$json = {
  '1': 'ListPolicyVersionsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'versions',
      '3': 252099085,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyVersion',
      '10': 'versions'
    },
  ],
};

/// Descriptor for `ListPolicyVersionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPolicyVersionsResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0UG9saWN5VmVyc2lvbnNSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaX'
        'N0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISMQoIdmVyc2lvbnMYjfSaeCAD'
        'KAsyEi5pYW0uUG9saWN5VmVyc2lvblIIdmVyc2lvbnM=');

@$core.Deprecated('Use listRolePoliciesRequestDescriptor instead')
const ListRolePoliciesRequest$json = {
  '1': 'ListRolePoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `ListRolePoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRolePoliciesRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0Um9sZVBvbGljaWVzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCg'
    'htYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEh4KCHJvbGVuYW1lGLPzvMIBIAEoCVIIcm9s'
    'ZW5hbWU=');

@$core.Deprecated('Use listRolePoliciesResponseDescriptor instead')
const ListRolePoliciesResponse$json = {
  '1': 'ListRolePoliciesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'policynames', '3': 264098782, '4': 3, '5': 9, '10': 'policynames'},
  ],
};

/// Descriptor for `ListRolePoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRolePoliciesResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0Um9sZVBvbGljaWVzUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2lzdH'
    'J1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEiMKC3BvbGljeW5hbWVzGN6n930g'
    'AygJUgtwb2xpY3luYW1lcw==');

@$core.Deprecated('Use listRoleTagsRequestDescriptor instead')
const ListRoleTagsRequest$json = {
  '1': 'ListRoleTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `ListRoleTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRoleTagsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0Um9sZVRhZ3NSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heG'
    'l0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbmFt'
    'ZQ==');

@$core.Deprecated('Use listRoleTagsResponseDescriptor instead')
const ListRoleTagsResponse$json = {
  '1': 'ListRoleTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListRoleTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRoleTagsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0Um9sZVRhZ3NSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2'
    'F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISIAoEdGFncxjBwfa1ASADKAsyCC5pYW0u'
    'VGFnUgR0YWdz');

@$core.Deprecated('Use listRolesRequestDescriptor instead')
const ListRolesRequest$json = {
  '1': 'ListRolesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
  ],
};

/// Descriptor for `ListRolesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRolesRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0Um9sZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW'
    '1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIgoKcGF0aHByZWZpeBjLt+a9ASABKAlSCnBhdGhwcmVm'
    'aXg=');

@$core.Deprecated('Use listRolesResponseDescriptor instead')
const ListRolesResponse$json = {
  '1': 'ListRolesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'roles',
      '3': 511168127,
      '4': 3,
      '5': 11,
      '6': '.iam.Role',
      '10': 'roles'
    },
  ],
};

/// Descriptor for `ListRolesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRolesResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0Um9sZXNSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZW'
    'QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISIwoFcm9sZXMY/5zf8wEgAygLMgkuaWFtLlJv'
    'bGVSBXJvbGVz');

@$core.Deprecated('Use listSAMLProviderTagsRequestDescriptor instead')
const ListSAMLProviderTagsRequest$json = {
  '1': 'ListSAMLProviderTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
  ],
};

/// Descriptor for `ListSAMLProviderTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSAMLProviderTagsRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0U0FNTFByb3ZpZGVyVGFnc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZX'
        'ISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIsCg9zYW1scHJvdmlkZXJhcm4Ypc2o'
        'jAEgASgJUg9zYW1scHJvdmlkZXJhcm4=');

@$core.Deprecated('Use listSAMLProviderTagsResponseDescriptor instead')
const ListSAMLProviderTagsResponse$json = {
  '1': 'ListSAMLProviderTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListSAMLProviderTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSAMLProviderTagsResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0U0FNTFByb3ZpZGVyVGFnc1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUg'
        'tpc3RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIgCgR0YWdzGMHB9rUBIAMo'
        'CzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listSAMLProvidersRequestDescriptor instead')
const ListSAMLProvidersRequest$json = {
  '1': 'ListSAMLProvidersRequest',
};

/// Descriptor for `ListSAMLProvidersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSAMLProvidersRequestDescriptor =
    $convert.base64Decode('ChhMaXN0U0FNTFByb3ZpZGVyc1JlcXVlc3Q=');

@$core.Deprecated('Use listSAMLProvidersResponseDescriptor instead')
const ListSAMLProvidersResponse$json = {
  '1': 'ListSAMLProvidersResponse',
  '2': [
    {
      '1': 'samlproviderlist',
      '3': 308003146,
      '4': 3,
      '5': 11,
      '6': '.iam.SAMLProviderListEntry',
      '10': 'samlproviderlist'
    },
  ],
};

/// Descriptor for `ListSAMLProvidersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSAMLProvidersResponseDescriptor =
    $convert.base64Decode(
        'ChlMaXN0U0FNTFByb3ZpZGVyc1Jlc3BvbnNlEkoKEHNhbWxwcm92aWRlcmxpc3QYyoLvkgEgAy'
        'gLMhouaWFtLlNBTUxQcm92aWRlckxpc3RFbnRyeVIQc2FtbHByb3ZpZGVybGlzdA==');

@$core.Deprecated('Use listSSHPublicKeysRequestDescriptor instead')
const ListSSHPublicKeysRequest$json = {
  '1': 'ListSSHPublicKeysRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListSSHPublicKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSSHPublicKeysRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0U1NIUHVibGljS2V5c1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHg'
    'oIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIeCgh1c2VybmFtZRj6wdThASABKAlSCHVz'
    'ZXJuYW1l');

@$core.Deprecated('Use listSSHPublicKeysResponseDescriptor instead')
const ListSSHPublicKeysResponse$json = {
  '1': 'ListSSHPublicKeysResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'sshpublickeys',
      '3': 469962073,
      '4': 3,
      '5': 11,
      '6': '.iam.SSHPublicKeyMetadata',
      '10': 'sshpublickeys'
    },
  ],
};

/// Descriptor for `ListSSHPublicKeysResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSSHPublicKeysResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0U1NIUHVibGljS2V5c1Jlc3BvbnNlEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3'
    'RydW5jYXRlZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchJDCg1zc2hwdWJsaWNrZXlzGNma'
    'jOABIAMoCzIZLmlhbS5TU0hQdWJsaWNLZXlNZXRhZGF0YVINc3NocHVibGlja2V5cw==');

@$core.Deprecated('Use listServerCertificateTagsRequestDescriptor instead')
const ListServerCertificateTagsRequest$json = {
  '1': 'ListServerCertificateTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
  ],
};

/// Descriptor for `ListServerCertificateTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServerCertificateTagsRequestDescriptor =
    $convert.base64Decode(
        'CiBMaXN0U2VydmVyQ2VydGlmaWNhdGVUYWdzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm'
        '1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEjcKFXNlcnZlcmNlcnRpZmlj'
        'YXRlbmFtZRj/gJseIAEoCVIVc2VydmVyY2VydGlmaWNhdGVuYW1l');

@$core.Deprecated('Use listServerCertificateTagsResponseDescriptor instead')
const ListServerCertificateTagsResponse$json = {
  '1': 'ListServerCertificateTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListServerCertificateTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServerCertificateTagsResponseDescriptor =
    $convert.base64Decode(
        'CiFMaXN0U2VydmVyQ2VydGlmaWNhdGVUYWdzUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cy'
        'ABKAhSC2lzdHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEiAKBHRhZ3MYwcH2'
        'tQEgAygLMgguaWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use listServerCertificatesRequestDescriptor instead')
const ListServerCertificatesRequest$json = {
  '1': 'ListServerCertificatesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
  ],
};

/// Descriptor for `ListServerCertificatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServerCertificatesRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0U2VydmVyQ2VydGlmaWNhdGVzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcm'
        'tlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEiIKCnBhdGhwcmVmaXgYy7fmvQEg'
        'ASgJUgpwYXRocHJlZml4');

@$core.Deprecated('Use listServerCertificatesResponseDescriptor instead')
const ListServerCertificatesResponse$json = {
  '1': 'ListServerCertificatesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'servercertificatemetadatalist',
      '3': 230344753,
      '4': 3,
      '5': 11,
      '6': '.iam.ServerCertificateMetadata',
      '10': 'servercertificatemetadatalist'
    },
  ],
};

/// Descriptor for `ListServerCertificatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServerCertificatesResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0U2VydmVyQ2VydGlmaWNhdGVzUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cyABKA'
        'hSC2lzdHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEmcKHXNlcnZlcmNlcnRp'
        'ZmljYXRlbWV0YWRhdGFsaXN0GLGQ620gAygLMh4uaWFtLlNlcnZlckNlcnRpZmljYXRlTWV0YW'
        'RhdGFSHXNlcnZlcmNlcnRpZmljYXRlbWV0YWRhdGFsaXN0');

@$core.Deprecated('Use listServiceSpecificCredentialsRequestDescriptor instead')
const ListServiceSpecificCredentialsRequest$json = {
  '1': 'ListServiceSpecificCredentialsRequest',
  '2': [
    {'1': 'allusers', '3': 97090269, '4': 1, '5': 8, '10': 'allusers'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListServiceSpecificCredentialsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServiceSpecificCredentialsRequestDescriptor =
    $convert.base64Decode(
        'CiVMaXN0U2VydmljZVNwZWNpZmljQ3JlZGVudGlhbHNSZXF1ZXN0Eh0KCGFsbHVzZXJzGN31pS'
        '4gASgIUghhbGx1c2VycxIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU'
        '1trxASABKAVSCG1heGl0ZW1zEiMKC3NlcnZpY2VuYW1lGOCq20EgASgJUgtzZXJ2aWNlbmFtZR'
        'IeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core
    .Deprecated('Use listServiceSpecificCredentialsResponseDescriptor instead')
const ListServiceSpecificCredentialsResponse$json = {
  '1': 'ListServiceSpecificCredentialsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'servicespecificcredentials',
      '3': 82595803,
      '4': 3,
      '5': 11,
      '6': '.iam.ServiceSpecificCredentialMetadata',
      '10': 'servicespecificcredentials'
    },
  ],
};

/// Descriptor for `ListServiceSpecificCredentialsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listServiceSpecificCredentialsResponseDescriptor =
    $convert.base64Decode(
        'CiZMaXN0U2VydmljZVNwZWNpZmljQ3JlZGVudGlhbHNSZXNwb25zZRIjCgtpc3RydW5jYXRlZB'
        'jan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISaQoac2Vy'
        'dmljZXNwZWNpZmljY3JlZGVudGlhbHMY25+xJyADKAsyJi5pYW0uU2VydmljZVNwZWNpZmljQ3'
        'JlZGVudGlhbE1ldGFkYXRhUhpzZXJ2aWNlc3BlY2lmaWNjcmVkZW50aWFscw==');

@$core.Deprecated('Use listSigningCertificatesRequestDescriptor instead')
const ListSigningCertificatesRequest$json = {
  '1': 'ListSigningCertificatesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListSigningCertificatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSigningCertificatesRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0U2lnbmluZ0NlcnRpZmljYXRlc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYX'
        'JrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIeCgh1c2VybmFtZRj6wdThASAB'
        'KAlSCHVzZXJuYW1l');

@$core.Deprecated('Use listSigningCertificatesResponseDescriptor instead')
const ListSigningCertificatesResponse$json = {
  '1': 'ListSigningCertificatesResponse',
  '2': [
    {
      '1': 'certificates',
      '3': 411978970,
      '4': 3,
      '5': 11,
      '6': '.iam.SigningCertificate',
      '10': 'certificates'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListSigningCertificatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSigningCertificatesResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0U2lnbmluZ0NlcnRpZmljYXRlc1Jlc3BvbnNlEj8KDGNlcnRpZmljYXRlcxjambnEAS'
        'ADKAsyFy5pYW0uU2lnbmluZ0NlcnRpZmljYXRlUgxjZXJ0aWZpY2F0ZXMSIwoLaXN0cnVuY2F0'
        'ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2Vy');

@$core.Deprecated('Use listUserPoliciesRequestDescriptor instead')
const ListUserPoliciesRequest$json = {
  '1': 'ListUserPoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListUserPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoliciesRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0VXNlclBvbGljaWVzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCg'
    'htYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNl'
    'cm5hbWU=');

@$core.Deprecated('Use listUserPoliciesResponseDescriptor instead')
const ListUserPoliciesResponse$json = {
  '1': 'ListUserPoliciesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'policynames', '3': 264098782, '4': 3, '5': 9, '10': 'policynames'},
  ],
};

/// Descriptor for `ListUserPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoliciesResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0VXNlclBvbGljaWVzUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2lzdH'
    'J1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEiMKC3BvbGljeW5hbWVzGN6n930g'
    'AygJUgtwb2xpY3luYW1lcw==');

@$core.Deprecated('Use listUserTagsRequestDescriptor instead')
const ListUserTagsRequest$json = {
  '1': 'ListUserTagsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ListUserTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserTagsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0VXNlclRhZ3NSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heG'
    'l0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFt'
    'ZQ==');

@$core.Deprecated('Use listUserTagsResponseDescriptor instead')
const ListUserTagsResponse$json = {
  '1': 'ListUserTagsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListUserTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserTagsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0VXNlclRhZ3NSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2'
    'F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISIAoEdGFncxjBwfa1ASADKAsyCC5pYW0u'
    'VGFnUgR0YWdz');

@$core.Deprecated('Use listUsersRequestDescriptor instead')
const ListUsersRequest$json = {
  '1': 'ListUsersRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'pathprefix', '3': 398040011, '4': 1, '5': 9, '10': 'pathprefix'},
  ],
};

/// Descriptor for `ListUsersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUsersRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0VXNlcnNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW'
    '1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIgoKcGF0aHByZWZpeBjLt+a9ASABKAlSCnBhdGhwcmVm'
    'aXg=');

@$core.Deprecated('Use listUsersResponseDescriptor instead')
const ListUsersResponse$json = {
  '1': 'ListUsersResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'users',
      '3': 112472756,
      '4': 3,
      '5': 11,
      '6': '.iam.User',
      '10': 'users'
    },
  ],
};

/// Descriptor for `ListUsersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUsersResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0VXNlcnNSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZW'
    'QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISIgoFdXNlcnMYtOXQNSADKAsyCS5pYW0uVXNl'
    'clIFdXNlcnM=');

@$core.Deprecated('Use listVirtualMFADevicesRequestDescriptor instead')
const ListVirtualMFADevicesRequest$json = {
  '1': 'ListVirtualMFADevicesRequest',
  '2': [
    {
      '1': 'assignmentstatus',
      '3': 277782963,
      '4': 1,
      '5': 14,
      '6': '.iam.assignmentStatusType',
      '10': 'assignmentstatus'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListVirtualMFADevicesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVirtualMFADevicesRequestDescriptor = $convert.base64Decode(
    'ChxMaXN0VmlydHVhbE1GQURldmljZXNSZXF1ZXN0EkkKEGFzc2lnbm1lbnRzdGF0dXMYs8O6hA'
    'EgASgOMhkuaWFtLmFzc2lnbm1lbnRTdGF0dXNUeXBlUhBhc3NpZ25tZW50c3RhdHVzEhkKBm1h'
    'cmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listVirtualMFADevicesResponseDescriptor instead')
const ListVirtualMFADevicesResponse$json = {
  '1': 'ListVirtualMFADevicesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'virtualmfadevices',
      '3': 117329102,
      '4': 3,
      '5': 11,
      '6': '.iam.VirtualMFADevice',
      '10': 'virtualmfadevices'
    },
  ],
};

/// Descriptor for `ListVirtualMFADevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVirtualMFADevicesResponseDescriptor = $convert.base64Decode(
    'Ch1MaXN0VmlydHVhbE1GQURldmljZXNSZXNwb25zZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCF'
    'ILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISRgoRdmlydHVhbG1mYWRl'
    'dmljZXMYzpn5NyADKAsyFS5pYW0uVmlydHVhbE1GQURldmljZVIRdmlydHVhbG1mYWRldmljZX'
    'M=');

@$core.Deprecated('Use loginProfileDescriptor instead')
const LoginProfile$json = {
  '1': 'LoginProfile',
  '2': [
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'passwordresetrequired',
      '3': 330943683,
      '4': 1,
      '5': 8,
      '10': 'passwordresetrequired'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `LoginProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginProfileDescriptor = $convert.base64Decode(
    'CgxMb2dpblByb2ZpbGUSIQoKY3JlYXRlZGF0ZRiSufwRIAEoCVIKY3JlYXRlZGF0ZRI4ChVwYX'
    'Nzd29yZHJlc2V0cmVxdWlyZWQYw5nnnQEgASgIUhVwYXNzd29yZHJlc2V0cmVxdWlyZWQSHgoI'
    'dXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use mFADeviceDescriptor instead')
const MFADevice$json = {
  '1': 'MFADevice',
  '2': [
    {'1': 'enabledate', '3': 99725723, '4': 1, '5': 9, '10': 'enabledate'},
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `MFADevice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mFADeviceDescriptor = $convert.base64Decode(
    'CglNRkFEZXZpY2USIQoKZW5hYmxlZGF0ZRib48YvIAEoCVIKZW5hYmxlZGF0ZRImCgxzZXJpYW'
    'xudW1iZXIY5bq5xwEgASgJUgxzZXJpYWxudW1iZXISHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1'
    'c2VybmFtZQ==');

@$core.Deprecated('Use malformedCertificateExceptionDescriptor instead')
const MalformedCertificateException$json = {
  '1': 'MalformedCertificateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MalformedCertificateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List malformedCertificateExceptionDescriptor =
    $convert.base64Decode(
        'Ch1NYWxmb3JtZWRDZXJ0aWZpY2F0ZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use malformedPolicyDocumentExceptionDescriptor instead')
const MalformedPolicyDocumentException$json = {
  '1': 'MalformedPolicyDocumentException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MalformedPolicyDocumentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List malformedPolicyDocumentExceptionDescriptor =
    $convert.base64Decode(
        'CiBNYWxmb3JtZWRQb2xpY3lEb2N1bWVudEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use managedPolicyDetailDescriptor instead')
const ManagedPolicyDetail$json = {
  '1': 'ManagedPolicyDetail',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'attachmentcount',
      '3': 468165416,
      '4': 1,
      '5': 5,
      '10': 'attachmentcount'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'defaultversionid',
      '3': 296003120,
      '4': 1,
      '5': 9,
      '10': 'defaultversionid'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'isattachable', '3': 2814573, '4': 1, '5': 8, '10': 'isattachable'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundaryusagecount',
      '3': 271917656,
      '4': 1,
      '5': 5,
      '10': 'permissionsboundaryusagecount'
    },
    {'1': 'policyid', '3': 299520499, '4': 1, '5': 9, '10': 'policyid'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policyversionlist',
      '3': 173611150,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyVersion',
      '10': 'policyversionlist'
    },
    {'1': 'updatedate', '3': 402617681, '4': 1, '5': 9, '10': 'updatedate'},
  ],
};

/// Descriptor for `ManagedPolicyDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedPolicyDetailDescriptor = $convert.base64Decode(
    'ChNNYW5hZ2VkUG9saWN5RGV0YWlsEhQKA2Fybhidm+2/ASABKAlSA2FybhIsCg9hdHRhY2htZW'
    '50Y291bnQYqMae3wEgASgFUg9hdHRhY2htZW50Y291bnQSIQoKY3JlYXRlZGF0ZRiSufwRIAEo'
    'CVIKY3JlYXRlZGF0ZRIuChBkZWZhdWx0dmVyc2lvbmlkGLDMko0BIAEoCVIQZGVmYXVsdHZlcn'
    'Npb25pZBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SJQoMaXNhdHRhY2hh'
    'YmxlGO3kqwEgASgIUgxpc2F0dGFjaGFibGUSFQoEcGF0aBjXyJtbIAEoCVIEcGF0aBJICh1wZX'
    'JtaXNzaW9uc2JvdW5kYXJ5dXNhZ2Vjb3VudBjYxNSBASABKAVSHXBlcm1pc3Npb25zYm91bmRh'
    'cnl1c2FnZWNvdW50Eh4KCHBvbGljeWlkGPOj6Y4BIAEoCVIIcG9saWN5aWQSIQoKcG9saWN5bm'
    'FtZRi99Yd/IAEoCVIKcG9saWN5bmFtZRJDChFwb2xpY3l2ZXJzaW9ubGlzdBiOseRSIAMoCzIS'
    'LmlhbS5Qb2xpY3lWZXJzaW9uUhFwb2xpY3l2ZXJzaW9ubGlzdBIiCgp1cGRhdGVkYXRlGNHq/b'
    '8BIAEoCVIKdXBkYXRlZGF0ZQ==');

@$core.Deprecated('Use noSuchEntityExceptionDescriptor instead')
const NoSuchEntityException$json = {
  '1': 'NoSuchEntityException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchEntityException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchEntityExceptionDescriptor = $convert.base64Decode(
    'ChVOb1N1Y2hFbnRpdHlFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use openIDConnectProviderListEntryDescriptor instead')
const OpenIDConnectProviderListEntry$json = {
  '1': 'OpenIDConnectProviderListEntry',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `OpenIDConnectProviderListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openIDConnectProviderListEntryDescriptor =
    $convert.base64Decode(
        'Ch5PcGVuSURDb25uZWN0UHJvdmlkZXJMaXN0RW50cnkSFAoDYXJuGJ2b7b8BIAEoCVIDYXJu');

@$core.Deprecated('Use openIdIdpCommunicationErrorExceptionDescriptor instead')
const OpenIdIdpCommunicationErrorException$json = {
  '1': 'OpenIdIdpCommunicationErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpenIdIdpCommunicationErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openIdIdpCommunicationErrorExceptionDescriptor =
    $convert.base64Decode(
        'CiRPcGVuSWRJZHBDb21tdW5pY2F0aW9uRXJyb3JFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use organizationNotFoundExceptionDescriptor instead')
const OrganizationNotFoundException$json = {
  '1': 'OrganizationNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OrganizationNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List organizationNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ch1Pcmdhbml6YXRpb25Ob3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated(
    'Use organizationNotInAllFeaturesModeExceptionDescriptor instead')
const OrganizationNotInAllFeaturesModeException$json = {
  '1': 'OrganizationNotInAllFeaturesModeException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OrganizationNotInAllFeaturesModeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    organizationNotInAllFeaturesModeExceptionDescriptor = $convert.base64Decode(
        'CilPcmdhbml6YXRpb25Ob3RJbkFsbEZlYXR1cmVzTW9kZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use organizationsDecisionDetailDescriptor instead')
const OrganizationsDecisionDetail$json = {
  '1': 'OrganizationsDecisionDetail',
  '2': [
    {
      '1': 'allowedbyorganizations',
      '3': 387308871,
      '4': 1,
      '5': 8,
      '10': 'allowedbyorganizations'
    },
  ],
};

/// Descriptor for `OrganizationsDecisionDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List organizationsDecisionDetailDescriptor =
    $convert.base64Decode(
        'ChtPcmdhbml6YXRpb25zRGVjaXNpb25EZXRhaWwSOgoWYWxsb3dlZGJ5b3JnYW5pemF0aW9ucx'
        'jHute4ASABKAhSFmFsbG93ZWRieW9yZ2FuaXphdGlvbnM=');

@$core.Deprecated('Use passwordPolicyDescriptor instead')
const PasswordPolicy$json = {
  '1': 'PasswordPolicy',
  '2': [
    {
      '1': 'allowuserstochangepassword',
      '3': 286109075,
      '4': 1,
      '5': 8,
      '10': 'allowuserstochangepassword'
    },
    {
      '1': 'expirepasswords',
      '3': 487573877,
      '4': 1,
      '5': 8,
      '10': 'expirepasswords'
    },
    {'1': 'hardexpiry', '3': 422015840, '4': 1, '5': 8, '10': 'hardexpiry'},
    {
      '1': 'maxpasswordage',
      '3': 58745976,
      '4': 1,
      '5': 5,
      '10': 'maxpasswordage'
    },
    {
      '1': 'minimumpasswordlength',
      '3': 202098235,
      '4': 1,
      '5': 5,
      '10': 'minimumpasswordlength'
    },
    {
      '1': 'passwordreuseprevention',
      '3': 469673813,
      '4': 1,
      '5': 5,
      '10': 'passwordreuseprevention'
    },
    {
      '1': 'requirelowercasecharacters',
      '3': 83014890,
      '4': 1,
      '5': 8,
      '10': 'requirelowercasecharacters'
    },
    {
      '1': 'requirenumbers',
      '3': 385325811,
      '4': 1,
      '5': 8,
      '10': 'requirenumbers'
    },
    {
      '1': 'requiresymbols',
      '3': 98968906,
      '4': 1,
      '5': 8,
      '10': 'requiresymbols'
    },
    {
      '1': 'requireuppercasecharacters',
      '3': 187130985,
      '4': 1,
      '5': 8,
      '10': 'requireuppercasecharacters'
    },
  ],
};

/// Descriptor for `PasswordPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List passwordPolicyDescriptor = $convert.base64Decode(
    'Cg5QYXNzd29yZFBvbGljeRJCChphbGxvd3VzZXJzdG9jaGFuZ2VwYXNzd29yZBiT27aIASABKA'
    'hSGmFsbG93dXNlcnN0b2NoYW5nZXBhc3N3b3JkEiwKD2V4cGlyZXBhc3N3b3Jkcxj1kr/oASAB'
    'KAhSD2V4cGlyZXBhc3N3b3JkcxIiCgpoYXJkZXhwaXJ5GODmnckBIAEoCFIKaGFyZGV4cGlyeR'
    'IpCg5tYXhwYXNzd29yZGFnZRj4yIEcIAEoBVIObWF4cGFzc3dvcmRhZ2USNwoVbWluaW11bXBh'
    'c3N3b3JkbGVuZ3RoGLuMr2AgASgFUhVtaW5pbXVtcGFzc3dvcmRsZW5ndGgSPAoXcGFzc3dvcm'
    'RyZXVzZXByZXZlbnRpb24Y1c763wEgASgFUhdwYXNzd29yZHJldXNlcHJldmVudGlvbhJBChpy'
    'ZXF1aXJlbG93ZXJjYXNlY2hhcmFjdGVycxjq6conIAEoCFIacmVxdWlyZWxvd2VyY2FzZWNoYX'
    'JhY3RlcnMSKgoOcmVxdWlyZW51bWJlcnMY87XetwEgASgIUg5yZXF1aXJlbnVtYmVycxIpCg5y'
    'ZXF1aXJlc3ltYm9scxjKypgvIAEoCFIOcmVxdWlyZXN5bWJvbHMSQQoacmVxdWlyZXVwcGVyY2'
    'FzZWNoYXJhY3RlcnMY6cidWSABKAhSGnJlcXVpcmV1cHBlcmNhc2VjaGFyYWN0ZXJz');

@$core.Deprecated('Use passwordPolicyViolationExceptionDescriptor instead')
const PasswordPolicyViolationException$json = {
  '1': 'PasswordPolicyViolationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PasswordPolicyViolationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List passwordPolicyViolationExceptionDescriptor =
    $convert.base64Decode(
        'CiBQYXNzd29yZFBvbGljeVZpb2xhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use permissionsBoundaryDecisionDetailDescriptor instead')
const PermissionsBoundaryDecisionDetail$json = {
  '1': 'PermissionsBoundaryDecisionDetail',
  '2': [
    {
      '1': 'allowedbypermissionsboundary',
      '3': 78647437,
      '4': 1,
      '5': 8,
      '10': 'allowedbypermissionsboundary'
    },
  ],
};

/// Descriptor for `PermissionsBoundaryDecisionDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List permissionsBoundaryDecisionDetailDescriptor =
    $convert.base64Decode(
        'CiFQZXJtaXNzaW9uc0JvdW5kYXJ5RGVjaXNpb25EZXRhaWwSRQocYWxsb3dlZGJ5cGVybWlzc2'
        'lvbnNib3VuZGFyeRiNocAlIAEoCFIcYWxsb3dlZGJ5cGVybWlzc2lvbnNib3VuZGFyeQ==');

@$core.Deprecated('Use policyDescriptor instead')
const Policy$json = {
  '1': 'Policy',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'attachmentcount',
      '3': 468165416,
      '4': 1,
      '5': 5,
      '10': 'attachmentcount'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'defaultversionid',
      '3': 296003120,
      '4': 1,
      '5': 9,
      '10': 'defaultversionid'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'isattachable', '3': 2814573, '4': 1, '5': 8, '10': 'isattachable'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundaryusagecount',
      '3': 271917656,
      '4': 1,
      '5': 5,
      '10': 'permissionsboundaryusagecount'
    },
    {'1': 'policyid', '3': 299520499, '4': 1, '5': 9, '10': 'policyid'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {'1': 'updatedate', '3': 402617681, '4': 1, '5': 9, '10': 'updatedate'},
  ],
};

/// Descriptor for `Policy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyDescriptor = $convert.base64Decode(
    'CgZQb2xpY3kSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEiwKD2F0dGFjaG1lbnRjb3VudBioxp7fAS'
    'ABKAVSD2F0dGFjaG1lbnRjb3VudBIhCgpjcmVhdGVkYXRlGJK5/BEgASgJUgpjcmVhdGVkYXRl'
    'Ei4KEGRlZmF1bHR2ZXJzaW9uaWQYsMySjQEgASgJUhBkZWZhdWx0dmVyc2lvbmlkEiMKC2Rlc2'
    'NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIlCgxpc2F0dGFjaGFibGUY7eSrASABKAhS'
    'DGlzYXR0YWNoYWJsZRIVCgRwYXRoGNfIm1sgASgJUgRwYXRoEkgKHXBlcm1pc3Npb25zYm91bm'
    'Rhcnl1c2FnZWNvdW50GNjE1IEBIAEoBVIdcGVybWlzc2lvbnNib3VuZGFyeXVzYWdlY291bnQS'
    'HgoIcG9saWN5aWQY86PpjgEgASgJUghwb2xpY3lpZBIhCgpwb2xpY3luYW1lGL31h38gASgJUg'
    'pwb2xpY3luYW1lEiAKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncxIiCgp1cGRhdGVk'
    'YXRlGNHq/b8BIAEoCVIKdXBkYXRlZGF0ZQ==');

@$core.Deprecated('Use policyDetailDescriptor instead')
const PolicyDetail$json = {
  '1': 'PolicyDetail',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `PolicyDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyDetailDescriptor = $convert.base64Decode(
    'CgxQb2xpY3lEZXRhaWwSKQoOcG9saWN5ZG9jdW1lbnQYy67BcSABKAlSDnBvbGljeWRvY3VtZW'
    '50EiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWU=');

@$core.Deprecated('Use policyEvaluationExceptionDescriptor instead')
const PolicyEvaluationException$json = {
  '1': 'PolicyEvaluationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PolicyEvaluationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyEvaluationExceptionDescriptor =
    $convert.base64Decode(
        'ChlQb2xpY3lFdmFsdWF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use policyGrantingServiceAccessDescriptor instead')
const PolicyGrantingServiceAccess$json = {
  '1': 'PolicyGrantingServiceAccess',
  '2': [
    {'1': 'entityname', '3': 79565666, '4': 1, '5': 9, '10': 'entityname'},
    {
      '1': 'entitytype',
      '3': 327275163,
      '4': 1,
      '5': 14,
      '6': '.iam.policyOwnerEntityType',
      '10': 'entitytype'
    },
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policytype',
      '3': 470619144,
      '4': 1,
      '5': 14,
      '6': '.iam.policyType',
      '10': 'policytype'
    },
  ],
};

/// Descriptor for `PolicyGrantingServiceAccess`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyGrantingServiceAccessDescriptor = $convert.base64Decode(
    'ChtQb2xpY3lHcmFudGluZ1NlcnZpY2VBY2Nlc3MSIQoKZW50aXR5bmFtZRjipvglIAEoCVIKZW'
    '50aXR5bmFtZRI+CgplbnRpdHl0eXBlGJulh5wBIAEoDjIaLmlhbS5wb2xpY3lPd25lckVudGl0'
    'eVR5cGVSCmVudGl0eXR5cGUSIAoJcG9saWN5YXJuGMPSuu0BIAEoCVIJcG9saWN5YXJuEiEKCn'
    'BvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWUSMwoKcG9saWN5dHlwZRiIqLTgASABKA4y'
    'Dy5pYW0ucG9saWN5VHlwZVIKcG9saWN5dHlwZQ==');

@$core.Deprecated('Use policyGroupDescriptor instead')
const PolicyGroup$json = {
  '1': 'PolicyGroup',
  '2': [
    {'1': 'groupid', '3': 73973250, '4': 1, '5': 9, '10': 'groupid'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
  ],
};

/// Descriptor for `PolicyGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyGroupDescriptor = $convert.base64Decode(
    'CgtQb2xpY3lHcm91cBIbCgdncm91cGlkGIL8oiMgASgJUgdncm91cGlkEiAKCWdyb3VwbmFtZR'
    'jIyqCqASABKAlSCWdyb3VwbmFtZQ==');

@$core.Deprecated('Use policyNotAttachableExceptionDescriptor instead')
const PolicyNotAttachableException$json = {
  '1': 'PolicyNotAttachableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PolicyNotAttachableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyNotAttachableExceptionDescriptor =
    $convert.base64Decode(
        'ChxQb2xpY3lOb3RBdHRhY2hhYmxlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use policyParameterDescriptor instead')
const PolicyParameter$json = {
  '1': 'PolicyParameter',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyParameterTypeEnum',
      '10': 'type'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `PolicyParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyParameterDescriptor = $convert.base64Decode(
    'Cg9Qb2xpY3lQYXJhbWV0ZXISFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRI0CgR0eXBlGO6g14oBIA'
    'EoDjIcLmlhbS5Qb2xpY3lQYXJhbWV0ZXJUeXBlRW51bVIEdHlwZRIZCgZ2YWx1ZXMY3MS0aiAD'
    'KAlSBnZhbHVlcw==');

@$core.Deprecated('Use policyRoleDescriptor instead')
const PolicyRole$json = {
  '1': 'PolicyRole',
  '2': [
    {'1': 'roleid', '3': 486019037, '4': 1, '5': 9, '10': 'roleid'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `PolicyRole`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyRoleDescriptor = $convert.base64Decode(
    'CgpQb2xpY3lSb2xlEhoKBnJvbGVpZBjdn+DnASABKAlSBnJvbGVpZBIeCghyb2xlbmFtZRiz87'
    'zCASABKAlSCHJvbGVuYW1l');

@$core.Deprecated('Use policyUserDescriptor instead')
const PolicyUser$json = {
  '1': 'PolicyUser',
  '2': [
    {'1': 'userid', '3': 10274112, '4': 1, '5': 9, '10': 'userid'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `PolicyUser`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyUserDescriptor = $convert.base64Decode(
    'CgpQb2xpY3lVc2VyEhkKBnVzZXJpZBjAivMEIAEoCVIGdXNlcmlkEh4KCHVzZXJuYW1lGPrB1O'
    'EBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use policyVersionDescriptor instead')
const PolicyVersion$json = {
  '1': 'PolicyVersion',
  '2': [
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {'1': 'document', '3': 407108341, '4': 1, '5': 9, '10': 'document'},
    {
      '1': 'isdefaultversion',
      '3': 465655635,
      '4': 1,
      '5': 8,
      '10': 'isdefaultversion'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PolicyVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyVersionDescriptor = $convert.base64Decode(
    'Cg1Qb2xpY3lWZXJzaW9uEiEKCmNyZWF0ZWRhdGUYkrn8ESABKAlSCmNyZWF0ZWRhdGUSHgoIZG'
    '9jdW1lbnQY9fWPwgEgASgJUghkb2N1bWVudBIuChBpc2RlZmF1bHR2ZXJzaW9uGNOuhd4BIAEo'
    'CFIQaXNkZWZhdWx0dmVyc2lvbhIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use positionDescriptor instead')
const Position$json = {
  '1': 'Position',
  '2': [
    {'1': 'column', '3': 46213704, '4': 1, '5': 5, '10': 'column'},
    {'1': 'line', '3': 403240264, '4': 1, '5': 5, '10': 'line'},
  ],
};

/// Descriptor for `Position`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List positionDescriptor = $convert.base64Decode(
    'CghQb3NpdGlvbhIZCgZjb2x1bW4YyNSEFiABKAVSBmNvbHVtbhIWCgRsaW5lGMjqo8ABIAEoBV'
    'IEbGluZQ==');

@$core.Deprecated('Use putGroupPolicyRequestDescriptor instead')
const PutGroupPolicyRequest$json = {
  '1': 'PutGroupPolicyRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `PutGroupPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putGroupPolicyRequestDescriptor = $convert.base64Decode(
    'ChVQdXRHcm91cFBvbGljeVJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW'
    '1lEikKDnBvbGljeWRvY3VtZW50GMuuwXEgASgJUg5wb2xpY3lkb2N1bWVudBIhCgpwb2xpY3lu'
    'YW1lGL31h38gASgJUgpwb2xpY3luYW1l');

@$core.Deprecated('Use putRolePermissionsBoundaryRequestDescriptor instead')
const PutRolePermissionsBoundaryRequest$json = {
  '1': 'PutRolePermissionsBoundaryRequest',
  '2': [
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 9,
      '10': 'permissionsboundary'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `PutRolePermissionsBoundaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRolePermissionsBoundaryRequestDescriptor =
    $convert.base64Decode(
        'CiFQdXRSb2xlUGVybWlzc2lvbnNCb3VuZGFyeVJlcXVlc3QSMwoTcGVybWlzc2lvbnNib3VuZG'
        'FyeRi2w6AIIAEoCVITcGVybWlzc2lvbnNib3VuZGFyeRIeCghyb2xlbmFtZRiz87zCASABKAlS'
        'CHJvbGVuYW1l');

@$core.Deprecated('Use putRolePolicyRequestDescriptor instead')
const PutRolePolicyRequest$json = {
  '1': 'PutRolePolicyRequest',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `PutRolePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRolePolicyRequestDescriptor = $convert.base64Decode(
    'ChRQdXRSb2xlUG9saWN5UmVxdWVzdBIpCg5wb2xpY3lkb2N1bWVudBjLrsFxIAEoCVIOcG9saW'
    'N5ZG9jdW1lbnQSIQoKcG9saWN5bmFtZRi99Yd/IAEoCVIKcG9saWN5bmFtZRIeCghyb2xlbmFt'
    'ZRiz87zCASABKAlSCHJvbGVuYW1l');

@$core.Deprecated('Use putUserPermissionsBoundaryRequestDescriptor instead')
const PutUserPermissionsBoundaryRequest$json = {
  '1': 'PutUserPermissionsBoundaryRequest',
  '2': [
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 9,
      '10': 'permissionsboundary'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `PutUserPermissionsBoundaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putUserPermissionsBoundaryRequestDescriptor =
    $convert.base64Decode(
        'CiFQdXRVc2VyUGVybWlzc2lvbnNCb3VuZGFyeVJlcXVlc3QSMwoTcGVybWlzc2lvbnNib3VuZG'
        'FyeRi2w6AIIAEoCVITcGVybWlzc2lvbnNib3VuZGFyeRIeCgh1c2VybmFtZRj6wdThASABKAlS'
        'CHVzZXJuYW1l');

@$core.Deprecated('Use putUserPolicyRequestDescriptor instead')
const PutUserPolicyRequest$json = {
  '1': 'PutUserPolicyRequest',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `PutUserPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putUserPolicyRequestDescriptor = $convert.base64Decode(
    'ChRQdXRVc2VyUG9saWN5UmVxdWVzdBIpCg5wb2xpY3lkb2N1bWVudBjLrsFxIAEoCVIOcG9saW'
    'N5ZG9jdW1lbnQSIQoKcG9saWN5bmFtZRi99Yd/IAEoCVIKcG9saWN5bmFtZRIeCgh1c2VybmFt'
    'ZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use rejectDelegationRequestRequestDescriptor instead')
const RejectDelegationRequestRequest$json = {
  '1': 'RejectDelegationRequestRequest',
  '2': [
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
    {'1': 'notes', '3': 267393899, '4': 1, '5': 9, '10': 'notes'},
  ],
};

/// Descriptor for `RejectDelegationRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectDelegationRequestRequestDescriptor =
    $convert.base64Decode(
        'Ch5SZWplY3REZWxlZ2F0aW9uUmVxdWVzdFJlcXVlc3QSNAoTZGVsZWdhdGlvbnJlcXVlc3RpZB'
        'iet4CFASABKAlSE2RlbGVnYXRpb25yZXF1ZXN0aWQSFwoFbm90ZXMY67bAfyABKAlSBW5vdGVz');

@$core.Deprecated(
    'Use removeClientIDFromOpenIDConnectProviderRequestDescriptor instead')
const RemoveClientIDFromOpenIDConnectProviderRequest$json = {
  '1': 'RemoveClientIDFromOpenIDConnectProviderRequest',
  '2': [
    {'1': 'clientid', '3': 448889284, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
  ],
};

/// Descriptor for `RemoveClientIDFromOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    removeClientIDFromOpenIDConnectProviderRequestDescriptor =
    $convert.base64Decode(
        'Ci5SZW1vdmVDbGllbnRJREZyb21PcGVuSURDb25uZWN0UHJvdmlkZXJSZXF1ZXN0Eh4KCGNsaW'
        'VudGlkGMSDhtYBIAEoCVIIY2xpZW50aWQSPgoYb3BlbmlkY29ubmVjdHByb3ZpZGVyYXJuGOPz'
        'j+kBIAEoCVIYb3BlbmlkY29ubmVjdHByb3ZpZGVyYXJu');

@$core.Deprecated('Use removeRoleFromInstanceProfileRequestDescriptor instead')
const RemoveRoleFromInstanceProfileRequest$json = {
  '1': 'RemoveRoleFromInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `RemoveRoleFromInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeRoleFromInstanceProfileRequestDescriptor =
    $convert.base64Decode(
        'CiRSZW1vdmVSb2xlRnJvbUluc3RhbmNlUHJvZmlsZVJlcXVlc3QSNAoTaW5zdGFuY2Vwcm9maW'
        'xlbmFtZRjtzrzaASABKAlSE2luc3RhbmNlcHJvZmlsZW5hbWUSHgoIcm9sZW5hbWUYs/O8wgEg'
        'ASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use removeUserFromGroupRequestDescriptor instead')
const RemoveUserFromGroupRequest$json = {
  '1': 'RemoveUserFromGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `RemoveUserFromGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeUserFromGroupRequestDescriptor =
    $convert.base64Decode(
        'ChpSZW1vdmVVc2VyRnJvbUdyb3VwUmVxdWVzdBIgCglncm91cG5hbWUYyMqgqgEgASgJUglncm'
        '91cG5hbWUSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core
    .Deprecated('Use reportGenerationLimitExceededExceptionDescriptor instead')
const ReportGenerationLimitExceededException$json = {
  '1': 'ReportGenerationLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ReportGenerationLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reportGenerationLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiZSZXBvcnRHZW5lcmF0aW9uTGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyC'
        'cgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use resetServiceSpecificCredentialRequestDescriptor instead')
const ResetServiceSpecificCredentialRequest$json = {
  '1': 'ResetServiceSpecificCredentialRequest',
  '2': [
    {
      '1': 'servicespecificcredentialid',
      '3': 426936633,
      '4': 1,
      '5': 9,
      '10': 'servicespecificcredentialid'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ResetServiceSpecificCredentialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetServiceSpecificCredentialRequestDescriptor =
    $convert.base64Decode(
        'CiVSZXNldFNlcnZpY2VTcGVjaWZpY0NyZWRlbnRpYWxSZXF1ZXN0EkQKG3NlcnZpY2VzcGVjaW'
        'ZpY2NyZWRlbnRpYWxpZBi5ksrLASABKAlSG3NlcnZpY2VzcGVjaWZpY2NyZWRlbnRpYWxpZBIe'
        'Cgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core
    .Deprecated('Use resetServiceSpecificCredentialResponseDescriptor instead')
const ResetServiceSpecificCredentialResponse$json = {
  '1': 'ResetServiceSpecificCredentialResponse',
  '2': [
    {
      '1': 'servicespecificcredential',
      '3': 186794126,
      '4': 1,
      '5': 11,
      '6': '.iam.ServiceSpecificCredential',
      '10': 'servicespecificcredential'
    },
  ],
};

/// Descriptor for `ResetServiceSpecificCredentialResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetServiceSpecificCredentialResponseDescriptor =
    $convert.base64Decode(
        'CiZSZXNldFNlcnZpY2VTcGVjaWZpY0NyZWRlbnRpYWxSZXNwb25zZRJfChlzZXJ2aWNlc3BlY2'
        'lmaWNjcmVkZW50aWFsGI6BiVkgASgLMh4uaWFtLlNlcnZpY2VTcGVjaWZpY0NyZWRlbnRpYWxS'
        'GXNlcnZpY2VzcGVjaWZpY2NyZWRlbnRpYWw=');

@$core.Deprecated('Use resourceSpecificResultDescriptor instead')
const ResourceSpecificResult$json = {
  '1': 'ResourceSpecificResult',
  '2': [
    {
      '1': 'evaldecisiondetails',
      '3': 70732822,
      '4': 3,
      '5': 11,
      '6': '.iam.ResourceSpecificResult.EvaldecisiondetailsEntry',
      '10': 'evaldecisiondetails'
    },
    {
      '1': 'evalresourcedecision',
      '3': 33317058,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyEvaluationDecisionType',
      '10': 'evalresourcedecision'
    },
    {
      '1': 'evalresourcename',
      '3': 9391249,
      '4': 1,
      '5': 9,
      '10': 'evalresourcename'
    },
    {
      '1': 'matchedstatements',
      '3': 389771374,
      '4': 3,
      '5': 11,
      '6': '.iam.Statement',
      '10': 'matchedstatements'
    },
    {
      '1': 'missingcontextvalues',
      '3': 526483253,
      '4': 3,
      '5': 9,
      '10': 'missingcontextvalues'
    },
    {
      '1': 'permissionsboundarydecisiondetail',
      '3': 480692611,
      '4': 1,
      '5': 11,
      '6': '.iam.PermissionsBoundaryDecisionDetail',
      '10': 'permissionsboundarydecisiondetail'
    },
  ],
  '3': [ResourceSpecificResult_EvaldecisiondetailsEntry$json],
};

@$core.Deprecated('Use resourceSpecificResultDescriptor instead')
const ResourceSpecificResult_EvaldecisiondetailsEntry$json = {
  '1': 'EvaldecisiondetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicyEvaluationDecisionType',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ResourceSpecificResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceSpecificResultDescriptor = $convert.base64Decode(
    'ChZSZXNvdXJjZVNwZWNpZmljUmVzdWx0EmkKE2V2YWxkZWNpc2lvbmRldGFpbHMYlpjdISADKA'
    'syNC5pYW0uUmVzb3VyY2VTcGVjaWZpY1Jlc3VsdC5FdmFsZGVjaXNpb25kZXRhaWxzRW50cnlS'
    'E2V2YWxkZWNpc2lvbmRldGFpbHMSWAoUZXZhbHJlc291cmNlZGVjaXNpb24YwsHxDyABKA4yIS'
    '5pYW0uUG9saWN5RXZhbHVhdGlvbkRlY2lzaW9uVHlwZVIUZXZhbHJlc291cmNlZGVjaXNpb24S'
    'LQoQZXZhbHJlc291cmNlbmFtZRiRmb0EIAEoCVIQZXZhbHJlc291cmNlbmFtZRJAChFtYXRjaG'
    'Vkc3RhdGVtZW50cxju4O25ASADKAsyDi5pYW0uU3RhdGVtZW50UhFtYXRjaGVkc3RhdGVtZW50'
    'cxI2ChRtaXNzaW5nY29udGV4dHZhbHVlcxi1/oX7ASADKAlSFG1pc3Npbmdjb250ZXh0dmFsdW'
    'VzEngKIXBlcm1pc3Npb25zYm91bmRhcnlkZWNpc2lvbmRldGFpbBiDk5vlASABKAsyJi5pYW0u'
    'UGVybWlzc2lvbnNCb3VuZGFyeURlY2lzaW9uRGV0YWlsUiFwZXJtaXNzaW9uc2JvdW5kYXJ5ZG'
    'VjaXNpb25kZXRhaWwaaQoYRXZhbGRlY2lzaW9uZGV0YWlsc0VudHJ5EhAKA2tleRgBIAEoCVID'
    'a2V5EjcKBXZhbHVlGAIgASgOMiEuaWFtLlBvbGljeUV2YWx1YXRpb25EZWNpc2lvblR5cGVSBX'
    'ZhbHVlOgI4AQ==');

@$core.Deprecated('Use resyncMFADeviceRequestDescriptor instead')
const ResyncMFADeviceRequest$json = {
  '1': 'ResyncMFADeviceRequest',
  '2': [
    {
      '1': 'authenticationcode1',
      '3': 444488540,
      '4': 1,
      '5': 9,
      '10': 'authenticationcode1'
    },
    {
      '1': 'authenticationcode2',
      '3': 461266159,
      '4': 1,
      '5': 9,
      '10': 'authenticationcode2'
    },
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ResyncMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resyncMFADeviceRequestDescriptor = $convert.base64Decode(
    'ChZSZXN5bmNNRkFEZXZpY2VSZXF1ZXN0EjQKE2F1dGhlbnRpY2F0aW9uY29kZTEY3Lb50wEgAS'
    'gJUhNhdXRoZW50aWNhdGlvbmNvZGUxEjQKE2F1dGhlbnRpY2F0aW9uY29kZTIY77n52wEgASgJ'
    'UhNhdXRoZW50aWNhdGlvbmNvZGUyEiYKDHNlcmlhbG51bWJlchjlurnHASABKAlSDHNlcmlhbG'
    '51bWJlchIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use roleDescriptor instead')
const Role$json = {
  '1': 'Role',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'assumerolepolicydocument',
      '3': 500384765,
      '4': 1,
      '5': 9,
      '10': 'assumerolepolicydocument'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'maxsessionduration',
      '3': 391414912,
      '4': 1,
      '5': 5,
      '10': 'maxsessionduration'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 11,
      '6': '.iam.AttachedPermissionsBoundary',
      '10': 'permissionsboundary'
    },
    {'1': 'roleid', '3': 486019037, '4': 1, '5': 9, '10': 'roleid'},
    {
      '1': 'rolelastused',
      '3': 298072393,
      '4': 1,
      '5': 11,
      '6': '.iam.RoleLastUsed',
      '10': 'rolelastused'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `Role`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List roleDescriptor = $convert.base64Decode(
    'CgRSb2xlEhQKA2Fybhidm+2/ASABKAlSA2FybhI+Chhhc3N1bWVyb2xlcG9saWN5ZG9jdW1lbn'
    'QY/YfN7gEgASgJUhhhc3N1bWVyb2xlcG9saWN5ZG9jdW1lbnQSIQoKY3JlYXRlZGF0ZRiSufwR'
    'IAEoCVIKY3JlYXRlZGF0ZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SMg'
    'oSbWF4c2Vzc2lvbmR1cmF0aW9uGICJ0roBIAEoBVISbWF4c2Vzc2lvbmR1cmF0aW9uEhUKBHBh'
    'dGgY18ibWyABKAlSBHBhdGgSVQoTcGVybWlzc2lvbnNib3VuZGFyeRi2w6AIIAEoCzIgLmlhbS'
    '5BdHRhY2hlZFBlcm1pc3Npb25zQm91bmRhcnlSE3Blcm1pc3Npb25zYm91bmRhcnkSGgoGcm9s'
    'ZWlkGN2f4OcBIAEoCVIGcm9sZWlkEjkKDHJvbGVsYXN0dXNlZBjJ8pCOASABKAsyES5pYW0uUm'
    '9sZUxhc3RVc2VkUgxyb2xlbGFzdHVzZWQSHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbmFt'
    'ZRIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use roleDetailDescriptor instead')
const RoleDetail$json = {
  '1': 'RoleDetail',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'assumerolepolicydocument',
      '3': 500384765,
      '4': 1,
      '5': 9,
      '10': 'assumerolepolicydocument'
    },
    {
      '1': 'attachedmanagedpolicies',
      '3': 155564765,
      '4': 3,
      '5': 11,
      '6': '.iam.AttachedPolicy',
      '10': 'attachedmanagedpolicies'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'instanceprofilelist',
      '3': 261628088,
      '4': 3,
      '5': 11,
      '6': '.iam.InstanceProfile',
      '10': 'instanceprofilelist'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 11,
      '6': '.iam.AttachedPermissionsBoundary',
      '10': 'permissionsboundary'
    },
    {'1': 'roleid', '3': 486019037, '4': 1, '5': 9, '10': 'roleid'},
    {
      '1': 'rolelastused',
      '3': 298072393,
      '4': 1,
      '5': 11,
      '6': '.iam.RoleLastUsed',
      '10': 'rolelastused'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
    {
      '1': 'rolepolicylist',
      '3': 536765692,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyDetail',
      '10': 'rolepolicylist'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `RoleDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List roleDetailDescriptor = $convert.base64Decode(
    'CgpSb2xlRGV0YWlsEhQKA2Fybhidm+2/ASABKAlSA2FybhI+Chhhc3N1bWVyb2xlcG9saWN5ZG'
    '9jdW1lbnQY/YfN7gEgASgJUhhhc3N1bWVyb2xlcG9saWN5ZG9jdW1lbnQSUAoXYXR0YWNoZWRt'
    'YW5hZ2VkcG9saWNpZXMY3fWWSiADKAsyEy5pYW0uQXR0YWNoZWRQb2xpY3lSF2F0dGFjaGVkbW'
    'FuYWdlZHBvbGljaWVzEiEKCmNyZWF0ZWRhdGUYkrn8ESABKAlSCmNyZWF0ZWRhdGUSSQoTaW5z'
    'dGFuY2Vwcm9maWxlbGlzdBi4weB8IAMoCzIULmlhbS5JbnN0YW5jZVByb2ZpbGVSE2luc3Rhbm'
    'NlcHJvZmlsZWxpc3QSFQoEcGF0aBjXyJtbIAEoCVIEcGF0aBJVChNwZXJtaXNzaW9uc2JvdW5k'
    'YXJ5GLbDoAggASgLMiAuaWFtLkF0dGFjaGVkUGVybWlzc2lvbnNCb3VuZGFyeVITcGVybWlzc2'
    'lvbnNib3VuZGFyeRIaCgZyb2xlaWQY3Z/g5wEgASgJUgZyb2xlaWQSOQoMcm9sZWxhc3R1c2Vk'
    'GMnykI4BIAEoCzIRLmlhbS5Sb2xlTGFzdFVzZWRSDHJvbGVsYXN0dXNlZBIeCghyb2xlbmFtZR'
    'iz87zCASABKAlSCHJvbGVuYW1lEj0KDnJvbGVwb2xpY3lsaXN0GPzJ+f8BIAMoCzIRLmlhbS5Q'
    'b2xpY3lEZXRhaWxSDnJvbGVwb2xpY3lsaXN0EiAKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1'
    'IEdGFncw==');

@$core.Deprecated('Use roleLastUsedDescriptor instead')
const RoleLastUsed$json = {
  '1': 'RoleLastUsed',
  '2': [
    {'1': 'lastuseddate', '3': 70377689, '4': 1, '5': 9, '10': 'lastuseddate'},
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
  ],
};

/// Descriptor for `RoleLastUsed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List roleLastUsedDescriptor = $convert.base64Decode(
    'CgxSb2xlTGFzdFVzZWQSJQoMbGFzdHVzZWRkYXRlGNnBxyEgASgJUgxsYXN0dXNlZGRhdGUSGQ'
    'oGcmVnaW9uGJ7xuUkgASgJUgZyZWdpb24=');

@$core.Deprecated('Use roleUsageTypeDescriptor instead')
const RoleUsageType$json = {
  '1': 'RoleUsageType',
  '2': [
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
    {'1': 'resources', '3': 358918291, '4': 3, '5': 9, '10': 'resources'},
  ],
};

/// Descriptor for `RoleUsageType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List roleUsageTypeDescriptor = $convert.base64Decode(
    'Cg1Sb2xlVXNhZ2VUeXBlEhkKBnJlZ2lvbhie8blJIAEoCVIGcmVnaW9uEiAKCXJlc291cmNlcx'
    'iT0ZKrASADKAlSCXJlc291cmNlcw==');

@$core.Deprecated('Use sAMLPrivateKeyDescriptor instead')
const SAMLPrivateKey$json = {
  '1': 'SAMLPrivateKey',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
  ],
};

/// Descriptor for `SAMLPrivateKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sAMLPrivateKeyDescriptor = $convert.base64Decode(
    'Cg5TQU1MUHJpdmF0ZUtleRIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEh8KCXRpbWVzdGFtcB'
    'jEw7dNIAEoCVIJdGltZXN0YW1w');

@$core.Deprecated('Use sAMLProviderListEntryDescriptor instead')
const SAMLProviderListEntry$json = {
  '1': 'SAMLProviderListEntry',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {'1': 'validuntil', '3': 366644404, '4': 1, '5': 9, '10': 'validuntil'},
  ],
};

/// Descriptor for `SAMLProviderListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sAMLProviderListEntryDescriptor = $convert.base64Decode(
    'ChVTQU1MUHJvdmlkZXJMaXN0RW50cnkSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEiEKCmNyZWF0ZW'
    'RhdGUYkrn8ESABKAlSCmNyZWF0ZWRhdGUSIgoKdmFsaWR1bnRpbBi0mequASABKAlSCnZhbGlk'
    'dW50aWw=');

@$core.Deprecated('Use sSHPublicKeyDescriptor instead')
const SSHPublicKey$json = {
  '1': 'SSHPublicKey',
  '2': [
    {'1': 'fingerprint', '3': 502547484, '4': 1, '5': 9, '10': 'fingerprint'},
    {
      '1': 'sshpublickeybody',
      '3': 502753956,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeybody'
    },
    {
      '1': 'sshpublickeyid',
      '3': 154499415,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeyid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'uploaddate', '3': 89572857, '4': 1, '5': 9, '10': 'uploaddate'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `SSHPublicKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSHPublicKeyDescriptor = $convert.base64Decode(
    'CgxTU0hQdWJsaWNLZXkSJAoLZmluZ2VycHJpbnQYnIjR7wEgASgJUgtmaW5nZXJwcmludBIuCh'
    'Bzc2hwdWJsaWNrZXlib2R5GKTV3e8BIAEoCVIQc3NocHVibGlja2V5Ym9keRIpCg5zc2hwdWJs'
    'aWNrZXlpZBjX8tVJIAEoCVIOc3NocHVibGlja2V5aWQSKgoGc3RhdHVzGJDk+wIgASgOMg8uaW'
    'FtLnN0YXR1c1R5cGVSBnN0YXR1cxIhCgp1cGxvYWRkYXRlGPmL2yogASgJUgp1cGxvYWRkYXRl'
    'Eh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use sSHPublicKeyMetadataDescriptor instead')
const SSHPublicKeyMetadata$json = {
  '1': 'SSHPublicKeyMetadata',
  '2': [
    {
      '1': 'sshpublickeyid',
      '3': 154499415,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeyid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'uploaddate', '3': 89572857, '4': 1, '5': 9, '10': 'uploaddate'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `SSHPublicKeyMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSHPublicKeyMetadataDescriptor = $convert.base64Decode(
    'ChRTU0hQdWJsaWNLZXlNZXRhZGF0YRIpCg5zc2hwdWJsaWNrZXlpZBjX8tVJIAEoCVIOc3NocH'
    'VibGlja2V5aWQSKgoGc3RhdHVzGJDk+wIgASgOMg8uaWFtLnN0YXR1c1R5cGVSBnN0YXR1cxIh'
    'Cgp1cGxvYWRkYXRlGPmL2yogASgJUgp1cGxvYWRkYXRlEh4KCHVzZXJuYW1lGPrB1OEBIAEoCV'
    'IIdXNlcm5hbWU=');

@$core.Deprecated('Use sendDelegationTokenRequestDescriptor instead')
const SendDelegationTokenRequest$json = {
  '1': 'SendDelegationTokenRequest',
  '2': [
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
  ],
};

/// Descriptor for `SendDelegationTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendDelegationTokenRequestDescriptor =
    $convert.base64Decode(
        'ChpTZW5kRGVsZWdhdGlvblRva2VuUmVxdWVzdBI0ChNkZWxlZ2F0aW9ucmVxdWVzdGlkGJ63gI'
        'UBIAEoCVITZGVsZWdhdGlvbnJlcXVlc3RpZA==');

@$core.Deprecated('Use serverCertificateDescriptor instead')
const ServerCertificate$json = {
  '1': 'ServerCertificate',
  '2': [
    {
      '1': 'certificatebody',
      '3': 147144541,
      '4': 1,
      '5': 9,
      '10': 'certificatebody'
    },
    {
      '1': 'certificatechain',
      '3': 369292378,
      '4': 1,
      '5': 9,
      '10': 'certificatechain'
    },
    {
      '1': 'servercertificatemetadata',
      '3': 274199561,
      '4': 1,
      '5': 11,
      '6': '.iam.ServerCertificateMetadata',
      '10': 'servercertificatemetadata'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ServerCertificate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverCertificateDescriptor = $convert.base64Decode(
    'ChFTZXJ2ZXJDZXJ0aWZpY2F0ZRIrCg9jZXJ0aWZpY2F0ZWJvZHkY3f6URiABKAlSD2NlcnRpZm'
    'ljYXRlYm9keRIuChBjZXJ0aWZpY2F0ZWNoYWluGNroi7ABIAEoCVIQY2VydGlmaWNhdGVjaGFp'
    'bhJgChlzZXJ2ZXJjZXJ0aWZpY2F0ZW1ldGFkYXRhGIno34IBIAEoCzIeLmlhbS5TZXJ2ZXJDZX'
    'J0aWZpY2F0ZU1ldGFkYXRhUhlzZXJ2ZXJjZXJ0aWZpY2F0ZW1ldGFkYXRhEiAKBHRhZ3MYwcH2'
    'tQEgAygLMgguaWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use serverCertificateMetadataDescriptor instead')
const ServerCertificateMetadata$json = {
  '1': 'ServerCertificateMetadata',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'servercertificateid',
      '3': 528841161,
      '4': 1,
      '5': 9,
      '10': 'servercertificateid'
    },
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
    {'1': 'uploaddate', '3': 89572857, '4': 1, '5': 9, '10': 'uploaddate'},
  ],
};

/// Descriptor for `ServerCertificateMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverCertificateMetadataDescriptor = $convert.base64Decode(
    'ChlTZXJ2ZXJDZXJ0aWZpY2F0ZU1ldGFkYXRhEhQKA2Fybhidm+2/ASABKAlSA2FybhIhCgpleH'
    'BpcmF0aW9uGImpn3UgASgJUgpleHBpcmF0aW9uEhUKBHBhdGgY18ibWyABKAlSBHBhdGgSNAoT'
    'c2VydmVyY2VydGlmaWNhdGVpZBjJ85X8ASABKAlSE3NlcnZlcmNlcnRpZmljYXRlaWQSNwoVc2'
    'VydmVyY2VydGlmaWNhdGVuYW1lGP+Amx4gASgJUhVzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbWUSIQoK'
    'dXBsb2FkZGF0ZRj5i9sqIAEoCVIKdXBsb2FkZGF0ZQ==');

@$core.Deprecated('Use serviceAccessNotEnabledExceptionDescriptor instead')
const ServiceAccessNotEnabledException$json = {
  '1': 'ServiceAccessNotEnabledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceAccessNotEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceAccessNotEnabledExceptionDescriptor =
    $convert.base64Decode(
        'CiBTZXJ2aWNlQWNjZXNzTm90RW5hYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use serviceFailureExceptionDescriptor instead')
const ServiceFailureException$json = {
  '1': 'ServiceFailureException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceFailureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceFailureExceptionDescriptor =
    $convert.base64Decode(
        'ChdTZXJ2aWNlRmFpbHVyZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use serviceLastAccessedDescriptor instead')
const ServiceLastAccessed$json = {
  '1': 'ServiceLastAccessed',
  '2': [
    {
      '1': 'lastauthenticated',
      '3': 389149429,
      '4': 1,
      '5': 9,
      '10': 'lastauthenticated'
    },
    {
      '1': 'lastauthenticatedentity',
      '3': 530201520,
      '4': 1,
      '5': 9,
      '10': 'lastauthenticatedentity'
    },
    {
      '1': 'lastauthenticatedregion',
      '3': 91136313,
      '4': 1,
      '5': 9,
      '10': 'lastauthenticatedregion'
    },
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
    {
      '1': 'servicenamespace',
      '3': 432654764,
      '4': 1,
      '5': 9,
      '10': 'servicenamespace'
    },
    {
      '1': 'totalauthenticatedentities',
      '3': 289355788,
      '4': 1,
      '5': 5,
      '10': 'totalauthenticatedentities'
    },
    {
      '1': 'trackedactionslastaccessed',
      '3': 217251704,
      '4': 3,
      '5': 11,
      '6': '.iam.TrackedActionLastAccessed',
      '10': 'trackedactionslastaccessed'
    },
  ],
};

/// Descriptor for `ServiceLastAccessed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceLastAccessedDescriptor = $convert.base64Decode(
    'ChNTZXJ2aWNlTGFzdEFjY2Vzc2VkEjAKEWxhc3RhdXRoZW50aWNhdGVkGPXlx7kBIAEoCVIRbG'
    'FzdGF1dGhlbnRpY2F0ZWQSPAoXbGFzdGF1dGhlbnRpY2F0ZWRlbnRpdHkYsPfo/AEgASgJUhds'
    'YXN0YXV0aGVudGljYXRlZGVudGl0eRI7ChdsYXN0YXV0aGVudGljYXRlZHJlZ2lvbhi5wrorIA'
    'EoCVIXbGFzdGF1dGhlbnRpY2F0ZWRyZWdpb24SIwoLc2VydmljZW5hbWUY4KrbQSABKAlSC3Nl'
    'cnZpY2VuYW1lEi4KEHNlcnZpY2VuYW1lc3BhY2UYrJOnzgEgASgJUhBzZXJ2aWNlbmFtZXNwYW'
    'NlEkIKGnRvdGFsYXV0aGVudGljYXRlZGVudGl0aWVzGIzw/IkBIAEoBVIadG90YWxhdXRoZW50'
    'aWNhdGVkZW50aXRpZXMSYQoadHJhY2tlZGFjdGlvbnNsYXN0YWNjZXNzZWQY+P7LZyADKAsyHi'
    '5pYW0uVHJhY2tlZEFjdGlvbkxhc3RBY2Nlc3NlZFIadHJhY2tlZGFjdGlvbnNsYXN0YWNjZXNz'
    'ZWQ=');

@$core.Deprecated('Use serviceNotSupportedExceptionDescriptor instead')
const ServiceNotSupportedException$json = {
  '1': 'ServiceNotSupportedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceNotSupportedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceNotSupportedExceptionDescriptor =
    $convert.base64Decode(
        'ChxTZXJ2aWNlTm90U3VwcG9ydGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use serviceSpecificCredentialDescriptor instead')
const ServiceSpecificCredential$json = {
  '1': 'ServiceSpecificCredential',
  '2': [
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'expirationdate',
      '3': 17577725,
      '4': 1,
      '5': 9,
      '10': 'expirationdate'
    },
    {
      '1': 'servicecredentialalias',
      '3': 202396008,
      '4': 1,
      '5': 9,
      '10': 'servicecredentialalias'
    },
    {
      '1': 'servicecredentialsecret',
      '3': 3133058,
      '4': 1,
      '5': 9,
      '10': 'servicecredentialsecret'
    },
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
    {
      '1': 'servicepassword',
      '3': 296023718,
      '4': 1,
      '5': 9,
      '10': 'servicepassword'
    },
    {
      '1': 'servicespecificcredentialid',
      '3': 426936633,
      '4': 1,
      '5': 9,
      '10': 'servicespecificcredentialid'
    },
    {
      '1': 'serviceusername',
      '3': 328571397,
      '4': 1,
      '5': 9,
      '10': 'serviceusername'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ServiceSpecificCredential`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceSpecificCredentialDescriptor = $convert.base64Decode(
    'ChlTZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFsEiEKCmNyZWF0ZWRhdGUYkrn8ESABKAlSCmNyZW'
    'F0ZWRhdGUSKQoOZXhwaXJhdGlvbmRhdGUY/e2wCCABKAlSDmV4cGlyYXRpb25kYXRlEjkKFnNl'
    'cnZpY2VjcmVkZW50aWFsYWxpYXMY6KLBYCABKAlSFnNlcnZpY2VjcmVkZW50aWFsYWxpYXMSOw'
    'oXc2VydmljZWNyZWRlbnRpYWxzZWNyZXQYgp2/ASABKAlSF3NlcnZpY2VjcmVkZW50aWFsc2Vj'
    'cmV0EiMKC3NlcnZpY2VuYW1lGOCq20EgASgJUgtzZXJ2aWNlbmFtZRIsCg9zZXJ2aWNlcGFzc3'
    'dvcmQYpu2TjQEgASgJUg9zZXJ2aWNlcGFzc3dvcmQSRAobc2VydmljZXNwZWNpZmljY3JlZGVu'
    'dGlhbGlkGLmSyssBIAEoCVIbc2VydmljZXNwZWNpZmljY3JlZGVudGlhbGlkEiwKD3NlcnZpY2'
    'V1c2VybmFtZRiFtNacASABKAlSD3NlcnZpY2V1c2VybmFtZRIqCgZzdGF0dXMYkOT7AiABKA4y'
    'Dy5pYW0uc3RhdHVzVHlwZVIGc3RhdHVzEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbW'
    'U=');

@$core.Deprecated('Use serviceSpecificCredentialMetadataDescriptor instead')
const ServiceSpecificCredentialMetadata$json = {
  '1': 'ServiceSpecificCredentialMetadata',
  '2': [
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'expirationdate',
      '3': 17577725,
      '4': 1,
      '5': 9,
      '10': 'expirationdate'
    },
    {
      '1': 'servicecredentialalias',
      '3': 202396008,
      '4': 1,
      '5': 9,
      '10': 'servicecredentialalias'
    },
    {'1': 'servicename', '3': 137811296, '4': 1, '5': 9, '10': 'servicename'},
    {
      '1': 'servicespecificcredentialid',
      '3': 426936633,
      '4': 1,
      '5': 9,
      '10': 'servicespecificcredentialid'
    },
    {
      '1': 'serviceusername',
      '3': 328571397,
      '4': 1,
      '5': 9,
      '10': 'serviceusername'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `ServiceSpecificCredentialMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceSpecificCredentialMetadataDescriptor = $convert.base64Decode(
    'CiFTZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFsTWV0YWRhdGESIQoKY3JlYXRlZGF0ZRiSufwRIA'
    'EoCVIKY3JlYXRlZGF0ZRIpCg5leHBpcmF0aW9uZGF0ZRj97bAIIAEoCVIOZXhwaXJhdGlvbmRh'
    'dGUSOQoWc2VydmljZWNyZWRlbnRpYWxhbGlhcxjoosFgIAEoCVIWc2VydmljZWNyZWRlbnRpYW'
    'xhbGlhcxIjCgtzZXJ2aWNlbmFtZRjgqttBIAEoCVILc2VydmljZW5hbWUSRAobc2VydmljZXNw'
    'ZWNpZmljY3JlZGVudGlhbGlkGLmSyssBIAEoCVIbc2VydmljZXNwZWNpZmljY3JlZGVudGlhbG'
    'lkEiwKD3NlcnZpY2V1c2VybmFtZRiFtNacASABKAlSD3NlcnZpY2V1c2VybmFtZRIqCgZzdGF0'
    'dXMYkOT7AiABKA4yDy5pYW0uc3RhdHVzVHlwZVIGc3RhdHVzEh4KCHVzZXJuYW1lGPrB1OEBIA'
    'EoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use setDefaultPolicyVersionRequestDescriptor instead')
const SetDefaultPolicyVersionRequest$json = {
  '1': 'SetDefaultPolicyVersionRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `SetDefaultPolicyVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setDefaultPolicyVersionRequestDescriptor =
    $convert.base64Decode(
        'Ch5TZXREZWZhdWx0UG9saWN5VmVyc2lvblJlcXVlc3QSIAoJcG9saWN5YXJuGMPSuu0BIAEoCV'
        'IJcG9saWN5YXJuEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated(
    'Use setSecurityTokenServicePreferencesRequestDescriptor instead')
const SetSecurityTokenServicePreferencesRequest$json = {
  '1': 'SetSecurityTokenServicePreferencesRequest',
  '2': [
    {
      '1': 'globalendpointtokenversion',
      '3': 110518557,
      '4': 1,
      '5': 14,
      '6': '.iam.globalEndpointTokenVersion',
      '10': 'globalendpointtokenversion'
    },
  ],
};

/// Descriptor for `SetSecurityTokenServicePreferencesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    setSecurityTokenServicePreferencesRequestDescriptor = $convert.base64Decode(
        'CilTZXRTZWN1cml0eVRva2VuU2VydmljZVByZWZlcmVuY2VzUmVxdWVzdBJiChpnbG9iYWxlbm'
        'Rwb2ludHRva2VudmVyc2lvbhidwtk0IAEoDjIfLmlhbS5nbG9iYWxFbmRwb2ludFRva2VuVmVy'
        'c2lvblIaZ2xvYmFsZW5kcG9pbnR0b2tlbnZlcnNpb24=');

@$core.Deprecated('Use signingCertificateDescriptor instead')
const SigningCertificate$json = {
  '1': 'SigningCertificate',
  '2': [
    {
      '1': 'certificatebody',
      '3': 147144541,
      '4': 1,
      '5': 9,
      '10': 'certificatebody'
    },
    {
      '1': 'certificateid',
      '3': 149606126,
      '4': 1,
      '5': 9,
      '10': 'certificateid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'uploaddate', '3': 89572857, '4': 1, '5': 9, '10': 'uploaddate'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `SigningCertificate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signingCertificateDescriptor = $convert.base64Decode(
    'ChJTaWduaW5nQ2VydGlmaWNhdGUSKwoPY2VydGlmaWNhdGVib2R5GN3+lEYgASgJUg9jZXJ0aW'
    'ZpY2F0ZWJvZHkSJwoNY2VydGlmaWNhdGVpZBjunatHIAEoCVINY2VydGlmaWNhdGVpZBIqCgZz'
    'dGF0dXMYkOT7AiABKA4yDy5pYW0uc3RhdHVzVHlwZVIGc3RhdHVzEiEKCnVwbG9hZGRhdGUY+Y'
    'vbKiABKAlSCnVwbG9hZGRhdGUSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use simulateCustomPolicyRequestDescriptor instead')
const SimulateCustomPolicyRequest$json = {
  '1': 'SimulateCustomPolicyRequest',
  '2': [
    {'1': 'actionnames', '3': 106940510, '4': 3, '5': 9, '10': 'actionnames'},
    {'1': 'callerarn', '3': 249627822, '4': 1, '5': 9, '10': 'callerarn'},
    {
      '1': 'contextentries',
      '3': 361913797,
      '4': 3,
      '5': 11,
      '6': '.iam.ContextEntry',
      '10': 'contextentries'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'permissionsboundarypolicyinputlist',
      '3': 424163706,
      '4': 3,
      '5': 9,
      '10': 'permissionsboundarypolicyinputlist'
    },
    {
      '1': 'policyinputlist',
      '3': 320766346,
      '4': 3,
      '5': 9,
      '10': 'policyinputlist'
    },
    {'1': 'resourcearns', '3': 222677390, '4': 3, '5': 9, '10': 'resourcearns'},
    {
      '1': 'resourcehandlingoption',
      '3': 447283374,
      '4': 1,
      '5': 9,
      '10': 'resourcehandlingoption'
    },
    {
      '1': 'resourceowner',
      '3': 523737189,
      '4': 1,
      '5': 9,
      '10': 'resourceowner'
    },
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
  ],
};

/// Descriptor for `SimulateCustomPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List simulateCustomPolicyRequestDescriptor = $convert.base64Decode(
    'ChtTaW11bGF0ZUN1c3RvbVBvbGljeVJlcXVlc3QSIwoLYWN0aW9ubmFtZXMY3pD/MiADKAlSC2'
    'FjdGlvbm5hbWVzEh8KCWNhbGxlcmFybhiuiYR3IAEoCVIJY2FsbGVyYXJuEj0KDmNvbnRleHRl'
    'bnRyaWVzGMW7yawBIAMoCzIRLmlhbS5Db250ZXh0RW50cnlSDmNvbnRleHRlbnRyaWVzEhkKBm'
    '1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMS'
    'UgoicGVybWlzc2lvbnNib3VuZGFyeXBvbGljeWlucHV0bGlzdBj68qDKASADKAlSInBlcm1pc3'
    'Npb25zYm91bmRhcnlwb2xpY3lpbnB1dGxpc3QSLAoPcG9saWN5aW5wdXRsaXN0GIqD+pgBIAMo'
    'CVIPcG9saWN5aW5wdXRsaXN0EiUKDHJlc291cmNlYXJucxiOk5dqIAMoCVIMcmVzb3VyY2Vhcm'
    '5zEjoKFnJlc291cmNlaGFuZGxpbmdvcHRpb24YroGk1QEgASgJUhZyZXNvdXJjZWhhbmRsaW5n'
    'b3B0aW9uEigKDXJlc291cmNlb3duZXIY5bDe+QEgASgJUg1yZXNvdXJjZW93bmVyEikKDnJlc2'
    '91cmNlcG9saWN5GLCUwQcgASgJUg5yZXNvdXJjZXBvbGljeQ==');

@$core.Deprecated('Use simulatePolicyResponseDescriptor instead')
const SimulatePolicyResponse$json = {
  '1': 'SimulatePolicyResponse',
  '2': [
    {
      '1': 'evaluationresults',
      '3': 444233744,
      '4': 3,
      '5': 11,
      '6': '.iam.EvaluationResult',
      '10': 'evaluationresults'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `SimulatePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List simulatePolicyResponseDescriptor = $convert.base64Decode(
    'ChZTaW11bGF0ZVBvbGljeVJlc3BvbnNlEkcKEWV2YWx1YXRpb25yZXN1bHRzGJDw6dMBIAMoCz'
    'IVLmlhbS5FdmFsdWF0aW9uUmVzdWx0UhFldmFsdWF0aW9ucmVzdWx0cxIjCgtpc3RydW5jYXRl'
    'ZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXI=');

@$core.Deprecated('Use simulatePrincipalPolicyRequestDescriptor instead')
const SimulatePrincipalPolicyRequest$json = {
  '1': 'SimulatePrincipalPolicyRequest',
  '2': [
    {'1': 'actionnames', '3': 106940510, '4': 3, '5': 9, '10': 'actionnames'},
    {'1': 'callerarn', '3': 249627822, '4': 1, '5': 9, '10': 'callerarn'},
    {
      '1': 'contextentries',
      '3': 361913797,
      '4': 3,
      '5': 11,
      '6': '.iam.ContextEntry',
      '10': 'contextentries'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'permissionsboundarypolicyinputlist',
      '3': 424163706,
      '4': 3,
      '5': 9,
      '10': 'permissionsboundarypolicyinputlist'
    },
    {
      '1': 'policyinputlist',
      '3': 320766346,
      '4': 3,
      '5': 9,
      '10': 'policyinputlist'
    },
    {
      '1': 'policysourcearn',
      '3': 84285350,
      '4': 1,
      '5': 9,
      '10': 'policysourcearn'
    },
    {'1': 'resourcearns', '3': 222677390, '4': 3, '5': 9, '10': 'resourcearns'},
    {
      '1': 'resourcehandlingoption',
      '3': 447283374,
      '4': 1,
      '5': 9,
      '10': 'resourcehandlingoption'
    },
    {
      '1': 'resourceowner',
      '3': 523737189,
      '4': 1,
      '5': 9,
      '10': 'resourceowner'
    },
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
  ],
};

/// Descriptor for `SimulatePrincipalPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List simulatePrincipalPolicyRequestDescriptor = $convert.base64Decode(
    'Ch5TaW11bGF0ZVByaW5jaXBhbFBvbGljeVJlcXVlc3QSIwoLYWN0aW9ubmFtZXMY3pD/MiADKA'
    'lSC2FjdGlvbm5hbWVzEh8KCWNhbGxlcmFybhiuiYR3IAEoCVIJY2FsbGVyYXJuEj0KDmNvbnRl'
    'eHRlbnRyaWVzGMW7yawBIAMoCzIRLmlhbS5Db250ZXh0RW50cnlSDmNvbnRleHRlbnRyaWVzEh'
    'kKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRl'
    'bXMSUgoicGVybWlzc2lvbnNib3VuZGFyeXBvbGljeWlucHV0bGlzdBj68qDKASADKAlSInBlcm'
    '1pc3Npb25zYm91bmRhcnlwb2xpY3lpbnB1dGxpc3QSLAoPcG9saWN5aW5wdXRsaXN0GIqD+pgB'
    'IAMoCVIPcG9saWN5aW5wdXRsaXN0EisKD3BvbGljeXNvdXJjZWFybhimr5goIAEoCVIPcG9saW'
    'N5c291cmNlYXJuEiUKDHJlc291cmNlYXJucxiOk5dqIAMoCVIMcmVzb3VyY2Vhcm5zEjoKFnJl'
    'c291cmNlaGFuZGxpbmdvcHRpb24YroGk1QEgASgJUhZyZXNvdXJjZWhhbmRsaW5nb3B0aW9uEi'
    'gKDXJlc291cmNlb3duZXIY5bDe+QEgASgJUg1yZXNvdXJjZW93bmVyEikKDnJlc291cmNlcG9s'
    'aWN5GLCUwQcgASgJUg5yZXNvdXJjZXBvbGljeQ==');

@$core.Deprecated('Use statementDescriptor instead')
const Statement$json = {
  '1': 'Statement',
  '2': [
    {
      '1': 'endposition',
      '3': 52449718,
      '4': 1,
      '5': 11,
      '6': '.iam.Position',
      '10': 'endposition'
    },
    {
      '1': 'sourcepolicyid',
      '3': 364395760,
      '4': 1,
      '5': 9,
      '10': 'sourcepolicyid'
    },
    {
      '1': 'sourcepolicytype',
      '3': 13742643,
      '4': 1,
      '5': 14,
      '6': '.iam.PolicySourceType',
      '10': 'sourcepolicytype'
    },
    {
      '1': 'startposition',
      '3': 216244189,
      '4': 1,
      '5': 11,
      '6': '.iam.Position',
      '10': 'startposition'
    },
  ],
};

/// Descriptor for `Statement`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statementDescriptor = $convert.base64Decode(
    'CglTdGF0ZW1lbnQSMgoLZW5kcG9zaXRpb24YtqOBGSABKAsyDS5pYW0uUG9zaXRpb25SC2VuZH'
    'Bvc2l0aW9uEioKDnNvdXJjZXBvbGljeWlkGPD54K0BIAEoCVIOc291cmNlcG9saWN5aWQSRAoQ'
    'c291cmNlcG9saWN5dHlwZRiz5MYGIAEoDjIVLmlhbS5Qb2xpY3lTb3VyY2VUeXBlUhBzb3VyY2'
    'Vwb2xpY3l0eXBlEjYKDXN0YXJ0cG9zaXRpb24Y3b+OZyABKAsyDS5pYW0uUG9zaXRpb25SDXN0'
    'YXJ0cG9zaXRpb24=');

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

@$core.Deprecated('Use tagInstanceProfileRequestDescriptor instead')
const TagInstanceProfileRequest$json = {
  '1': 'TagInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagInstanceProfileRequestDescriptor = $convert.base64Decode(
    'ChlUYWdJbnN0YW5jZVByb2ZpbGVSZXF1ZXN0EjQKE2luc3RhbmNlcHJvZmlsZW5hbWUY7c682g'
    'EgASgJUhNpbnN0YW5jZXByb2ZpbGVuYW1lEiAKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1IE'
    'dGFncw==');

@$core.Deprecated('Use tagMFADeviceRequestDescriptor instead')
const TagMFADeviceRequest$json = {
  '1': 'TagMFADeviceRequest',
  '2': [
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagMFADeviceRequestDescriptor = $convert.base64Decode(
    'ChNUYWdNRkFEZXZpY2VSZXF1ZXN0EiYKDHNlcmlhbG51bWJlchjlurnHASABKAlSDHNlcmlhbG'
    '51bWJlchIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagOpenIDConnectProviderRequestDescriptor instead')
const TagOpenIDConnectProviderRequest$json = {
  '1': 'TagOpenIDConnectProviderRequest',
  '2': [
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagOpenIDConnectProviderRequestDescriptor =
    $convert.base64Decode(
        'Ch9UYWdPcGVuSURDb25uZWN0UHJvdmlkZXJSZXF1ZXN0Ej4KGG9wZW5pZGNvbm5lY3Rwcm92aW'
        'RlcmFybhjj84/pASABKAlSGG9wZW5pZGNvbm5lY3Rwcm92aWRlcmFybhIgCgR0YWdzGMHB9rUB'
        'IAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagPolicyRequestDescriptor instead')
const TagPolicyRequest$json = {
  '1': 'TagPolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagPolicyRequestDescriptor = $convert.base64Decode(
    'ChBUYWdQb2xpY3lSZXF1ZXN0EiAKCXBvbGljeWFybhjD0rrtASABKAlSCXBvbGljeWFybhIgCg'
    'R0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagRoleRequestDescriptor instead')
const TagRoleRequest$json = {
  '1': 'TagRoleRequest',
  '2': [
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagRoleRequestDescriptor = $convert.base64Decode(
    'Cg5UYWdSb2xlUmVxdWVzdBIeCghyb2xlbmFtZRiz87zCASABKAlSCHJvbGVuYW1lEiAKBHRhZ3'
    'MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use tagSAMLProviderRequestDescriptor instead')
const TagSAMLProviderRequest$json = {
  '1': 'TagSAMLProviderRequest',
  '2': [
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagSAMLProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagSAMLProviderRequestDescriptor =
    $convert.base64Decode(
        'ChZUYWdTQU1MUHJvdmlkZXJSZXF1ZXN0EiwKD3NhbWxwcm92aWRlcmFybhilzaiMASABKAlSD3'
        'NhbWxwcm92aWRlcmFybhIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagServerCertificateRequestDescriptor instead')
const TagServerCertificateRequest$json = {
  '1': 'TagServerCertificateRequest',
  '2': [
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagServerCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagServerCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChtUYWdTZXJ2ZXJDZXJ0aWZpY2F0ZVJlcXVlc3QSNwoVc2VydmVyY2VydGlmaWNhdGVuYW1lGP'
        '+Amx4gASgJUhVzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbWUSIAoEdGFncxjBwfa1ASADKAsyCC5pYW0u'
        'VGFnUgR0YWdz');

@$core.Deprecated('Use tagUserRequestDescriptor instead')
const TagUserRequest$json = {
  '1': 'TagUserRequest',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `TagUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagUserRequestDescriptor = $convert.base64Decode(
    'Cg5UYWdVc2VyUmVxdWVzdBIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3MSHgoIdX'
    'Nlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use trackedActionLastAccessedDescriptor instead')
const TrackedActionLastAccessed$json = {
  '1': 'TrackedActionLastAccessed',
  '2': [
    {'1': 'actionname', '3': 115541053, '4': 1, '5': 9, '10': 'actionname'},
    {
      '1': 'lastaccessedentity',
      '3': 442988598,
      '4': 1,
      '5': 9,
      '10': 'lastaccessedentity'
    },
    {
      '1': 'lastaccessedregion',
      '3': 213791971,
      '4': 1,
      '5': 9,
      '10': 'lastaccessedregion'
    },
    {
      '1': 'lastaccessedtime',
      '3': 407082528,
      '4': 1,
      '5': 9,
      '10': 'lastaccessedtime'
    },
  ],
};

/// Descriptor for `TrackedActionLastAccessed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trackedActionLastAccessedDescriptor = $convert.base64Decode(
    'ChlUcmFja2VkQWN0aW9uTGFzdEFjY2Vzc2VkEiEKCmFjdGlvbm5hbWUYvYiMNyABKAlSCmFjdG'
    'lvbm5hbWUSMgoSbGFzdGFjY2Vzc2VkZW50aXR5GLbwndMBIAEoCVISbGFzdGFjY2Vzc2VkZW50'
    'aXR5EjEKEmxhc3RhY2Nlc3NlZHJlZ2lvbhjj6fhlIAEoCVISbGFzdGFjY2Vzc2VkcmVnaW9uEi'
    '4KEGxhc3RhY2Nlc3NlZHRpbWUYoKyOwgEgASgJUhBsYXN0YWNjZXNzZWR0aW1l');

@$core.Deprecated('Use unmodifiableEntityExceptionDescriptor instead')
const UnmodifiableEntityException$json = {
  '1': 'UnmodifiableEntityException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnmodifiableEntityException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unmodifiableEntityExceptionDescriptor =
    $convert.base64Decode(
        'ChtVbm1vZGlmaWFibGVFbnRpdHlFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core
    .Deprecated('Use unrecognizedPublicKeyEncodingExceptionDescriptor instead')
const UnrecognizedPublicKeyEncodingException$json = {
  '1': 'UnrecognizedPublicKeyEncodingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnrecognizedPublicKeyEncodingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unrecognizedPublicKeyEncodingExceptionDescriptor =
    $convert.base64Decode(
        'CiZVbnJlY29nbml6ZWRQdWJsaWNLZXlFbmNvZGluZ0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyC'
        'cgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use untagInstanceProfileRequestDescriptor instead')
const UntagInstanceProfileRequest$json = {
  '1': 'UntagInstanceProfileRequest',
  '2': [
    {
      '1': 'instanceprofilename',
      '3': 458172269,
      '4': 1,
      '5': 9,
      '10': 'instanceprofilename'
    },
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagInstanceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagInstanceProfileRequestDescriptor =
    $convert.base64Decode(
        'ChtVbnRhZ0luc3RhbmNlUHJvZmlsZVJlcXVlc3QSNAoTaW5zdGFuY2Vwcm9maWxlbmFtZRjtzr'
        'zaASABKAlSE2luc3RhbmNlcHJvZmlsZW5hbWUSHAoHdGFna2V5cxj8w/OYASADKAlSB3RhZ2tl'
        'eXM=');

@$core.Deprecated('Use untagMFADeviceRequestDescriptor instead')
const UntagMFADeviceRequest$json = {
  '1': 'UntagMFADeviceRequest',
  '2': [
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagMFADeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagMFADeviceRequestDescriptor = $convert.base64Decode(
    'ChVVbnRhZ01GQURldmljZVJlcXVlc3QSJgoMc2VyaWFsbnVtYmVyGOW6uccBIAEoCVIMc2VyaW'
    'FsbnVtYmVyEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagOpenIDConnectProviderRequestDescriptor instead')
const UntagOpenIDConnectProviderRequest$json = {
  '1': 'UntagOpenIDConnectProviderRequest',
  '2': [
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagOpenIDConnectProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagOpenIDConnectProviderRequestDescriptor =
    $convert.base64Decode(
        'CiFVbnRhZ09wZW5JRENvbm5lY3RQcm92aWRlclJlcXVlc3QSPgoYb3BlbmlkY29ubmVjdHByb3'
        'ZpZGVyYXJuGOPzj+kBIAEoCVIYb3BlbmlkY29ubmVjdHByb3ZpZGVyYXJuEhwKB3RhZ2tleXMY'
        '/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagPolicyRequestDescriptor instead')
const UntagPolicyRequest$json = {
  '1': 'UntagPolicyRequest',
  '2': [
    {'1': 'policyarn', '3': 497985859, '4': 1, '5': 9, '10': 'policyarn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagPolicyRequestDescriptor = $convert.base64Decode(
    'ChJVbnRhZ1BvbGljeVJlcXVlc3QSIAoJcG9saWN5YXJuGMPSuu0BIAEoCVIJcG9saWN5YXJuEh'
    'wKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagRoleRequestDescriptor instead')
const UntagRoleRequest$json = {
  '1': 'UntagRoleRequest',
  '2': [
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagRoleRequestDescriptor = $convert.base64Decode(
    'ChBVbnRhZ1JvbGVSZXF1ZXN0Eh4KCHJvbGVuYW1lGLPzvMIBIAEoCVIIcm9sZW5hbWUSHAoHdG'
    'Fna2V5cxj8w/OYASADKAlSB3RhZ2tleXM=');

@$core.Deprecated('Use untagSAMLProviderRequestDescriptor instead')
const UntagSAMLProviderRequest$json = {
  '1': 'UntagSAMLProviderRequest',
  '2': [
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagSAMLProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagSAMLProviderRequestDescriptor =
    $convert.base64Decode(
        'ChhVbnRhZ1NBTUxQcm92aWRlclJlcXVlc3QSLAoPc2FtbHByb3ZpZGVyYXJuGKXNqIwBIAEoCV'
        'IPc2FtbHByb3ZpZGVyYXJuEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagServerCertificateRequestDescriptor instead')
const UntagServerCertificateRequest$json = {
  '1': 'UntagServerCertificateRequest',
  '2': [
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagServerCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagServerCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch1VbnRhZ1NlcnZlckNlcnRpZmljYXRlUmVxdWVzdBI3ChVzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbW'
        'UY/4CbHiABKAlSFXNlcnZlcmNlcnRpZmljYXRlbmFtZRIcCgd0YWdrZXlzGPzD85gBIAMoCVIH'
        'dGFna2V5cw==');

@$core.Deprecated('Use untagUserRequestDescriptor instead')
const UntagUserRequest$json = {
  '1': 'UntagUserRequest',
  '2': [
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UntagUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagUserRequestDescriptor = $convert.base64Decode(
    'ChBVbnRhZ1VzZXJSZXF1ZXN0EhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlzEh4KCHVzZX'
    'JuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use updateAccessKeyRequestDescriptor instead')
const UpdateAccessKeyRequest$json = {
  '1': 'UpdateAccessKeyRequest',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateAccessKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccessKeyRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVBY2Nlc3NLZXlSZXF1ZXN0EiQKC2FjY2Vzc2tleWlkGKC3t9gBIAEoCVILYWNjZX'
    'Nza2V5aWQSKgoGc3RhdHVzGJDk+wIgASgOMg8uaWFtLnN0YXR1c1R5cGVSBnN0YXR1cxIeCgh1'
    'c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use updateAccountPasswordPolicyRequestDescriptor instead')
const UpdateAccountPasswordPolicyRequest$json = {
  '1': 'UpdateAccountPasswordPolicyRequest',
  '2': [
    {
      '1': 'allowuserstochangepassword',
      '3': 286109075,
      '4': 1,
      '5': 8,
      '10': 'allowuserstochangepassword'
    },
    {'1': 'hardexpiry', '3': 422015840, '4': 1, '5': 8, '10': 'hardexpiry'},
    {
      '1': 'maxpasswordage',
      '3': 58745976,
      '4': 1,
      '5': 5,
      '10': 'maxpasswordage'
    },
    {
      '1': 'minimumpasswordlength',
      '3': 202098235,
      '4': 1,
      '5': 5,
      '10': 'minimumpasswordlength'
    },
    {
      '1': 'passwordreuseprevention',
      '3': 469673813,
      '4': 1,
      '5': 5,
      '10': 'passwordreuseprevention'
    },
    {
      '1': 'requirelowercasecharacters',
      '3': 83014890,
      '4': 1,
      '5': 8,
      '10': 'requirelowercasecharacters'
    },
    {
      '1': 'requirenumbers',
      '3': 385325811,
      '4': 1,
      '5': 8,
      '10': 'requirenumbers'
    },
    {
      '1': 'requiresymbols',
      '3': 98968906,
      '4': 1,
      '5': 8,
      '10': 'requiresymbols'
    },
    {
      '1': 'requireuppercasecharacters',
      '3': 187130985,
      '4': 1,
      '5': 8,
      '10': 'requireuppercasecharacters'
    },
  ],
};

/// Descriptor for `UpdateAccountPasswordPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccountPasswordPolicyRequestDescriptor = $convert.base64Decode(
    'CiJVcGRhdGVBY2NvdW50UGFzc3dvcmRQb2xpY3lSZXF1ZXN0EkIKGmFsbG93dXNlcnN0b2NoYW'
    '5nZXBhc3N3b3JkGJPbtogBIAEoCFIaYWxsb3d1c2Vyc3RvY2hhbmdlcGFzc3dvcmQSIgoKaGFy'
    'ZGV4cGlyeRjg5p3JASABKAhSCmhhcmRleHBpcnkSKQoObWF4cGFzc3dvcmRhZ2UY+MiBHCABKA'
    'VSDm1heHBhc3N3b3JkYWdlEjcKFW1pbmltdW1wYXNzd29yZGxlbmd0aBi7jK9gIAEoBVIVbWlu'
    'aW11bXBhc3N3b3JkbGVuZ3RoEjwKF3Bhc3N3b3JkcmV1c2VwcmV2ZW50aW9uGNXO+t8BIAEoBV'
    'IXcGFzc3dvcmRyZXVzZXByZXZlbnRpb24SQQoacmVxdWlyZWxvd2VyY2FzZWNoYXJhY3RlcnMY'
    '6unKJyABKAhSGnJlcXVpcmVsb3dlcmNhc2VjaGFyYWN0ZXJzEioKDnJlcXVpcmVudW1iZXJzGP'
    'O13rcBIAEoCFIOcmVxdWlyZW51bWJlcnMSKQoOcmVxdWlyZXN5bWJvbHMYysqYLyABKAhSDnJl'
    'cXVpcmVzeW1ib2xzEkEKGnJlcXVpcmV1cHBlcmNhc2VjaGFyYWN0ZXJzGOnInVkgASgIUhpyZX'
    'F1aXJldXBwZXJjYXNlY2hhcmFjdGVycw==');

@$core.Deprecated('Use updateAssumeRolePolicyRequestDescriptor instead')
const UpdateAssumeRolePolicyRequest$json = {
  '1': 'UpdateAssumeRolePolicyRequest',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `UpdateAssumeRolePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAssumeRolePolicyRequestDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVBc3N1bWVSb2xlUG9saWN5UmVxdWVzdBIpCg5wb2xpY3lkb2N1bWVudBjLrsFxIA'
        'EoCVIOcG9saWN5ZG9jdW1lbnQSHgoIcm9sZW5hbWUYs/O8wgEgASgJUghyb2xlbmFtZQ==');

@$core.Deprecated('Use updateDelegationRequestRequestDescriptor instead')
const UpdateDelegationRequestRequest$json = {
  '1': 'UpdateDelegationRequestRequest',
  '2': [
    {
      '1': 'delegationrequestid',
      '3': 278928286,
      '4': 1,
      '5': 9,
      '10': 'delegationrequestid'
    },
    {'1': 'notes', '3': 267393899, '4': 1, '5': 9, '10': 'notes'},
  ],
};

/// Descriptor for `UpdateDelegationRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDelegationRequestRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVEZWxlZ2F0aW9uUmVxdWVzdFJlcXVlc3QSNAoTZGVsZWdhdGlvbnJlcXVlc3RpZB'
        'iet4CFASABKAlSE2RlbGVnYXRpb25yZXF1ZXN0aWQSFwoFbm90ZXMY67bAfyABKAlSBW5vdGVz');

@$core.Deprecated('Use updateGroupRequestDescriptor instead')
const UpdateGroupRequest$json = {
  '1': 'UpdateGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'newgroupname', '3': 103615532, '4': 1, '5': 9, '10': 'newgroupname'},
    {'1': 'newpath', '3': 101274219, '4': 1, '5': 9, '10': 'newpath'},
  ],
};

/// Descriptor for `UpdateGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGroupRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVHcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEi'
    'UKDG5ld2dyb3VwbmFtZRismLQxIAEoCVIMbmV3Z3JvdXBuYW1lEhsKB25ld3BhdGgY66SlMCAB'
    'KAlSB25ld3BhdGg=');

@$core.Deprecated('Use updateLoginProfileRequestDescriptor instead')
const UpdateLoginProfileRequest$json = {
  '1': 'UpdateLoginProfileRequest',
  '2': [
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {
      '1': 'passwordresetrequired',
      '3': 330943683,
      '4': 1,
      '5': 8,
      '10': 'passwordresetrequired'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateLoginProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateLoginProfileRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVMb2dpblByb2ZpbGVSZXF1ZXN0Eh0KCHBhc3N3b3JkGLmQjGYgASgJUghwYXNzd2'
    '9yZBI4ChVwYXNzd29yZHJlc2V0cmVxdWlyZWQYw5nnnQEgASgIUhVwYXNzd29yZHJlc2V0cmVx'
    'dWlyZWQSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated(
    'Use updateOpenIDConnectProviderThumbprintRequestDescriptor instead')
const UpdateOpenIDConnectProviderThumbprintRequest$json = {
  '1': 'UpdateOpenIDConnectProviderThumbprintRequest',
  '2': [
    {
      '1': 'openidconnectproviderarn',
      '3': 488896995,
      '4': 1,
      '5': 9,
      '10': 'openidconnectproviderarn'
    },
    {
      '1': 'thumbprintlist',
      '3': 88528351,
      '4': 3,
      '5': 9,
      '10': 'thumbprintlist'
    },
  ],
};

/// Descriptor for `UpdateOpenIDConnectProviderThumbprintRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateOpenIDConnectProviderThumbprintRequestDescriptor =
    $convert.base64Decode(
        'CixVcGRhdGVPcGVuSURDb25uZWN0UHJvdmlkZXJUaHVtYnByaW50UmVxdWVzdBI+ChhvcGVuaW'
        'Rjb25uZWN0cHJvdmlkZXJhcm4Y4/OP6QEgASgJUhhvcGVuaWRjb25uZWN0cHJvdmlkZXJhcm4S'
        'KQoOdGh1bWJwcmludGxpc3QY36ubKiADKAlSDnRodW1icHJpbnRsaXN0');

@$core.Deprecated('Use updateRoleDescriptionRequestDescriptor instead')
const UpdateRoleDescriptionRequest$json = {
  '1': 'UpdateRoleDescriptionRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `UpdateRoleDescriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRoleDescriptionRequestDescriptor =
    $convert.base64Decode(
        'ChxVcGRhdGVSb2xlRGVzY3JpcHRpb25SZXF1ZXN0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUg'
        'tkZXNjcmlwdGlvbhIeCghyb2xlbmFtZRiz87zCASABKAlSCHJvbGVuYW1l');

@$core.Deprecated('Use updateRoleDescriptionResponseDescriptor instead')
const UpdateRoleDescriptionResponse$json = {
  '1': 'UpdateRoleDescriptionResponse',
  '2': [
    {
      '1': 'role',
      '3': 271285818,
      '4': 1,
      '5': 11,
      '6': '.iam.Role',
      '10': 'role'
    },
  ],
};

/// Descriptor for `UpdateRoleDescriptionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRoleDescriptionResponseDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVSb2xlRGVzY3JpcHRpb25SZXNwb25zZRIhCgRyb2xlGLr8rYEBIAEoCzIJLmlhbS'
        '5Sb2xlUgRyb2xl');

@$core.Deprecated('Use updateRoleRequestDescriptor instead')
const UpdateRoleRequest$json = {
  '1': 'UpdateRoleRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'maxsessionduration',
      '3': 391414912,
      '4': 1,
      '5': 5,
      '10': 'maxsessionduration'
    },
    {'1': 'rolename', '3': 407845299, '4': 1, '5': 9, '10': 'rolename'},
  ],
};

/// Descriptor for `UpdateRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRoleRequestDescriptor = $convert.base64Decode(
    'ChFVcGRhdGVSb2xlUmVxdWVzdBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb2'
    '4SMgoSbWF4c2Vzc2lvbmR1cmF0aW9uGICJ0roBIAEoBVISbWF4c2Vzc2lvbmR1cmF0aW9uEh4K'
    'CHJvbGVuYW1lGLPzvMIBIAEoCVIIcm9sZW5hbWU=');

@$core.Deprecated('Use updateRoleResponseDescriptor instead')
const UpdateRoleResponse$json = {
  '1': 'UpdateRoleResponse',
};

/// Descriptor for `UpdateRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRoleResponseDescriptor =
    $convert.base64Decode('ChJVcGRhdGVSb2xlUmVzcG9uc2U=');

@$core.Deprecated('Use updateSAMLProviderRequestDescriptor instead')
const UpdateSAMLProviderRequest$json = {
  '1': 'UpdateSAMLProviderRequest',
  '2': [
    {
      '1': 'addprivatekey',
      '3': 205767607,
      '4': 1,
      '5': 9,
      '10': 'addprivatekey'
    },
    {
      '1': 'assertionencryptionmode',
      '3': 474560298,
      '4': 1,
      '5': 14,
      '6': '.iam.assertionEncryptionModeType',
      '10': 'assertionencryptionmode'
    },
    {
      '1': 'removeprivatekey',
      '3': 506848260,
      '4': 1,
      '5': 9,
      '10': 'removeprivatekey'
    },
    {
      '1': 'samlmetadatadocument',
      '3': 282432645,
      '4': 1,
      '5': 9,
      '10': 'samlmetadatadocument'
    },
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
  ],
};

/// Descriptor for `UpdateSAMLProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSAMLProviderRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVTQU1MUHJvdmlkZXJSZXF1ZXN0EicKDWFkZHByaXZhdGVrZXkYt4ePYiABKAlSDW'
    'FkZHByaXZhdGVrZXkSXgoXYXNzZXJ0aW9uZW5jcnlwdGlvbm1vZGUYqu6k4gEgASgOMiAuaWFt'
    'LmFzc2VydGlvbkVuY3J5cHRpb25Nb2RlVHlwZVIXYXNzZXJ0aW9uZW5jcnlwdGlvbm1vZGUSLg'
    'oQcmVtb3ZlcHJpdmF0ZWtleRiEyNfxASABKAlSEHJlbW92ZXByaXZhdGVrZXkSNgoUc2FtbG1l'
    'dGFkYXRhZG9jdW1lbnQYhanWhgEgASgJUhRzYW1sbWV0YWRhdGFkb2N1bWVudBIsCg9zYW1scH'
    'JvdmlkZXJhcm4Ypc2ojAEgASgJUg9zYW1scHJvdmlkZXJhcm4=');

@$core.Deprecated('Use updateSAMLProviderResponseDescriptor instead')
const UpdateSAMLProviderResponse$json = {
  '1': 'UpdateSAMLProviderResponse',
  '2': [
    {
      '1': 'samlproviderarn',
      '3': 294266533,
      '4': 1,
      '5': 9,
      '10': 'samlproviderarn'
    },
  ],
};

/// Descriptor for `UpdateSAMLProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSAMLProviderResponseDescriptor =
    $convert.base64Decode(
        'ChpVcGRhdGVTQU1MUHJvdmlkZXJSZXNwb25zZRIsCg9zYW1scHJvdmlkZXJhcm4Ypc2ojAEgAS'
        'gJUg9zYW1scHJvdmlkZXJhcm4=');

@$core.Deprecated('Use updateSSHPublicKeyRequestDescriptor instead')
const UpdateSSHPublicKeyRequest$json = {
  '1': 'UpdateSSHPublicKeyRequest',
  '2': [
    {
      '1': 'sshpublickeyid',
      '3': 154499415,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeyid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateSSHPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSSHPublicKeyRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVTU0hQdWJsaWNLZXlSZXF1ZXN0EikKDnNzaHB1YmxpY2tleWlkGNfy1UkgASgJUg'
    '5zc2hwdWJsaWNrZXlpZBIqCgZzdGF0dXMYkOT7AiABKA4yDy5pYW0uc3RhdHVzVHlwZVIGc3Rh'
    'dHVzEh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use updateServerCertificateRequestDescriptor instead')
const UpdateServerCertificateRequest$json = {
  '1': 'UpdateServerCertificateRequest',
  '2': [
    {'1': 'newpath', '3': 101274219, '4': 1, '5': 9, '10': 'newpath'},
    {
      '1': 'newservercertificatename',
      '3': 68932667,
      '4': 1,
      '5': 9,
      '10': 'newservercertificatename'
    },
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
  ],
};

/// Descriptor for `UpdateServerCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateServerCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVTZXJ2ZXJDZXJ0aWZpY2F0ZVJlcXVlc3QSGwoHbmV3cGF0aBjrpKUwIAEoCVIHbm'
        'V3cGF0aBI9ChhuZXdzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbWUYu6jvICABKAlSGG5ld3NlcnZlcmNl'
        'cnRpZmljYXRlbmFtZRI3ChVzZXJ2ZXJjZXJ0aWZpY2F0ZW5hbWUY/4CbHiABKAlSFXNlcnZlcm'
        'NlcnRpZmljYXRlbmFtZQ==');

@$core
    .Deprecated('Use updateServiceSpecificCredentialRequestDescriptor instead')
const UpdateServiceSpecificCredentialRequest$json = {
  '1': 'UpdateServiceSpecificCredentialRequest',
  '2': [
    {
      '1': 'servicespecificcredentialid',
      '3': 426936633,
      '4': 1,
      '5': 9,
      '10': 'servicespecificcredentialid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateServiceSpecificCredentialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateServiceSpecificCredentialRequestDescriptor =
    $convert.base64Decode(
        'CiZVcGRhdGVTZXJ2aWNlU3BlY2lmaWNDcmVkZW50aWFsUmVxdWVzdBJEChtzZXJ2aWNlc3BlY2'
        'lmaWNjcmVkZW50aWFsaWQYuZLKywEgASgJUhtzZXJ2aWNlc3BlY2lmaWNjcmVkZW50aWFsaWQS'
        'KgoGc3RhdHVzGJDk+wIgASgOMg8uaWFtLnN0YXR1c1R5cGVSBnN0YXR1cxIeCgh1c2VybmFtZR'
        'j6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use updateSigningCertificateRequestDescriptor instead')
const UpdateSigningCertificateRequest$json = {
  '1': 'UpdateSigningCertificateRequest',
  '2': [
    {
      '1': 'certificateid',
      '3': 149606126,
      '4': 1,
      '5': 9,
      '10': 'certificateid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.iam.statusType',
      '10': 'status'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateSigningCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSigningCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVTaWduaW5nQ2VydGlmaWNhdGVSZXF1ZXN0EicKDWNlcnRpZmljYXRlaWQY7p2rRy'
        'ABKAlSDWNlcnRpZmljYXRlaWQSKgoGc3RhdHVzGJDk+wIgASgOMg8uaWFtLnN0YXR1c1R5cGVS'
        'BnN0YXR1cxIeCgh1c2VybmFtZRj6wdThASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use updateUserRequestDescriptor instead')
const UpdateUserRequest$json = {
  '1': 'UpdateUserRequest',
  '2': [
    {'1': 'newpath', '3': 101274219, '4': 1, '5': 9, '10': 'newpath'},
    {'1': 'newusername', '3': 339337854, '4': 1, '5': 9, '10': 'newusername'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserRequestDescriptor = $convert.base64Decode(
    'ChFVcGRhdGVVc2VyUmVxdWVzdBIbCgduZXdwYXRoGOukpTAgASgJUgduZXdwYXRoEiQKC25ld3'
    'VzZXJuYW1lGP7E56EBIAEoCVILbmV3dXNlcm5hbWUSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1'
    'c2VybmFtZQ==');

@$core.Deprecated('Use uploadSSHPublicKeyRequestDescriptor instead')
const UploadSSHPublicKeyRequest$json = {
  '1': 'UploadSSHPublicKeyRequest',
  '2': [
    {
      '1': 'sshpublickeybody',
      '3': 502753956,
      '4': 1,
      '5': 9,
      '10': 'sshpublickeybody'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UploadSSHPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadSSHPublicKeyRequestDescriptor =
    $convert.base64Decode(
        'ChlVcGxvYWRTU0hQdWJsaWNLZXlSZXF1ZXN0Ei4KEHNzaHB1YmxpY2tleWJvZHkYpNXd7wEgAS'
        'gJUhBzc2hwdWJsaWNrZXlib2R5Eh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use uploadSSHPublicKeyResponseDescriptor instead')
const UploadSSHPublicKeyResponse$json = {
  '1': 'UploadSSHPublicKeyResponse',
  '2': [
    {
      '1': 'sshpublickey',
      '3': 520385596,
      '4': 1,
      '5': 11,
      '6': '.iam.SSHPublicKey',
      '10': 'sshpublickey'
    },
  ],
};

/// Descriptor for `UploadSSHPublicKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadSSHPublicKeyResponseDescriptor =
    $convert.base64Decode(
        'ChpVcGxvYWRTU0hQdWJsaWNLZXlSZXNwb25zZRI5Cgxzc2hwdWJsaWNrZXkYvOiR+AEgASgLMh'
        'EuaWFtLlNTSFB1YmxpY0tleVIMc3NocHVibGlja2V5');

@$core.Deprecated('Use uploadServerCertificateRequestDescriptor instead')
const UploadServerCertificateRequest$json = {
  '1': 'UploadServerCertificateRequest',
  '2': [
    {
      '1': 'certificatebody',
      '3': 147144541,
      '4': 1,
      '5': 9,
      '10': 'certificatebody'
    },
    {
      '1': 'certificatechain',
      '3': 369292378,
      '4': 1,
      '5': 9,
      '10': 'certificatechain'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {'1': 'privatekey', '3': 173312700, '4': 1, '5': 9, '10': 'privatekey'},
    {
      '1': 'servercertificatename',
      '3': 63357055,
      '4': 1,
      '5': 9,
      '10': 'servercertificatename'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `UploadServerCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadServerCertificateRequestDescriptor = $convert.base64Decode(
    'Ch5VcGxvYWRTZXJ2ZXJDZXJ0aWZpY2F0ZVJlcXVlc3QSKwoPY2VydGlmaWNhdGVib2R5GN3+lE'
    'YgASgJUg9jZXJ0aWZpY2F0ZWJvZHkSLgoQY2VydGlmaWNhdGVjaGFpbhja6IuwASABKAlSEGNl'
    'cnRpZmljYXRlY2hhaW4SFQoEcGF0aBjXyJtbIAEoCVIEcGF0aBIhCgpwcml2YXRla2V5GLyV0l'
    'IgASgJUgpwcml2YXRla2V5EjcKFXNlcnZlcmNlcnRpZmljYXRlbmFtZRj/gJseIAEoCVIVc2Vy'
    'dmVyY2VydGlmaWNhdGVuYW1lEiAKBHRhZ3MYwcH2tQEgAygLMgguaWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use uploadServerCertificateResponseDescriptor instead')
const UploadServerCertificateResponse$json = {
  '1': 'UploadServerCertificateResponse',
  '2': [
    {
      '1': 'servercertificatemetadata',
      '3': 274199561,
      '4': 1,
      '5': 11,
      '6': '.iam.ServerCertificateMetadata',
      '10': 'servercertificatemetadata'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `UploadServerCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadServerCertificateResponseDescriptor =
    $convert.base64Decode(
        'Ch9VcGxvYWRTZXJ2ZXJDZXJ0aWZpY2F0ZVJlc3BvbnNlEmAKGXNlcnZlcmNlcnRpZmljYXRlbW'
        'V0YWRhdGEYiejfggEgASgLMh4uaWFtLlNlcnZlckNlcnRpZmljYXRlTWV0YWRhdGFSGXNlcnZl'
        'cmNlcnRpZmljYXRlbWV0YWRhdGESIAoEdGFncxjBwfa1ASADKAsyCC5pYW0uVGFnUgR0YWdz');

@$core.Deprecated('Use uploadSigningCertificateRequestDescriptor instead')
const UploadSigningCertificateRequest$json = {
  '1': 'UploadSigningCertificateRequest',
  '2': [
    {
      '1': 'certificatebody',
      '3': 147144541,
      '4': 1,
      '5': 9,
      '10': 'certificatebody'
    },
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UploadSigningCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadSigningCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch9VcGxvYWRTaWduaW5nQ2VydGlmaWNhdGVSZXF1ZXN0EisKD2NlcnRpZmljYXRlYm9keRjd/p'
        'RGIAEoCVIPY2VydGlmaWNhdGVib2R5Eh4KCHVzZXJuYW1lGPrB1OEBIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use uploadSigningCertificateResponseDescriptor instead')
const UploadSigningCertificateResponse$json = {
  '1': 'UploadSigningCertificateResponse',
  '2': [
    {
      '1': 'certificate',
      '3': 198060817,
      '4': 1,
      '5': 11,
      '6': '.iam.SigningCertificate',
      '10': 'certificate'
    },
  ],
};

/// Descriptor for `UploadSigningCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadSigningCertificateResponseDescriptor =
    $convert.base64Decode(
        'CiBVcGxvYWRTaWduaW5nQ2VydGlmaWNhdGVSZXNwb25zZRI8CgtjZXJ0aWZpY2F0ZRiR1rheIA'
        'EoCzIXLmlhbS5TaWduaW5nQ2VydGlmaWNhdGVSC2NlcnRpZmljYXRl');

@$core.Deprecated('Use userDescriptor instead')
const User$json = {
  '1': 'User',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'passwordlastused',
      '3': 177372702,
      '4': 1,
      '5': 9,
      '10': 'passwordlastused'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 11,
      '6': '.iam.AttachedPermissionsBoundary',
      '10': 'permissionsboundary'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {'1': 'userid', '3': 10274112, '4': 1, '5': 9, '10': 'userid'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `User`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userDescriptor = $convert.base64Decode(
    'CgRVc2VyEhQKA2Fybhidm+2/ASABKAlSA2FybhIhCgpjcmVhdGVkYXRlGJK5/BEgASgJUgpjcm'
    'VhdGVkYXRlEi0KEHBhc3N3b3JkbGFzdHVzZWQYnvzJVCABKAlSEHBhc3N3b3JkbGFzdHVzZWQS'
    'FQoEcGF0aBjXyJtbIAEoCVIEcGF0aBJVChNwZXJtaXNzaW9uc2JvdW5kYXJ5GLbDoAggASgLMi'
    'AuaWFtLkF0dGFjaGVkUGVybWlzc2lvbnNCb3VuZGFyeVITcGVybWlzc2lvbnNib3VuZGFyeRIg'
    'CgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3MSGQoGdXNlcmlkGMCK8wQgASgJUgZ1c2'
    'VyaWQSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use userDetailDescriptor instead')
const UserDetail$json = {
  '1': 'UserDetail',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'attachedmanagedpolicies',
      '3': 155564765,
      '4': 3,
      '5': 11,
      '6': '.iam.AttachedPolicy',
      '10': 'attachedmanagedpolicies'
    },
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {'1': 'grouplist', '3': 283781749, '4': 3, '5': 9, '10': 'grouplist'},
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {
      '1': 'permissionsboundary',
      '3': 17310134,
      '4': 1,
      '5': 11,
      '6': '.iam.AttachedPermissionsBoundary',
      '10': 'permissionsboundary'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {'1': 'userid', '3': 10274112, '4': 1, '5': 9, '10': 'userid'},
    {'1': 'username', '3': 473243898, '4': 1, '5': 9, '10': 'username'},
    {
      '1': 'userpolicylist',
      '3': 327234653,
      '4': 3,
      '5': 11,
      '6': '.iam.PolicyDetail',
      '10': 'userpolicylist'
    },
  ],
};

/// Descriptor for `UserDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userDetailDescriptor = $convert.base64Decode(
    'CgpVc2VyRGV0YWlsEhQKA2Fybhidm+2/ASABKAlSA2FybhJQChdhdHRhY2hlZG1hbmFnZWRwb2'
    'xpY2llcxjd9ZZKIAMoCzITLmlhbS5BdHRhY2hlZFBvbGljeVIXYXR0YWNoZWRtYW5hZ2VkcG9s'
    'aWNpZXMSIQoKY3JlYXRlZGF0ZRiSufwRIAEoCVIKY3JlYXRlZGF0ZRIgCglncm91cGxpc3QY9d'
    'SohwEgAygJUglncm91cGxpc3QSFQoEcGF0aBjXyJtbIAEoCVIEcGF0aBJVChNwZXJtaXNzaW9u'
    'c2JvdW5kYXJ5GLbDoAggASgLMiAuaWFtLkF0dGFjaGVkUGVybWlzc2lvbnNCb3VuZGFyeVITcG'
    'VybWlzc2lvbnNib3VuZGFyeRIgCgR0YWdzGMHB9rUBIAMoCzIILmlhbS5UYWdSBHRhZ3MSGQoG'
    'dXNlcmlkGMCK8wQgASgJUgZ1c2VyaWQSHgoIdXNlcm5hbWUY+sHU4QEgASgJUgh1c2VybmFtZR'
    'I9Cg51c2VycG9saWN5bGlzdBjd6IScASADKAsyES5pYW0uUG9saWN5RGV0YWlsUg51c2VycG9s'
    'aWN5bGlzdA==');

@$core.Deprecated('Use virtualMFADeviceDescriptor instead')
const VirtualMFADevice$json = {
  '1': 'VirtualMFADevice',
  '2': [
    {
      '1': 'base32stringseed',
      '3': 257240642,
      '4': 1,
      '5': 12,
      '10': 'base32stringseed'
    },
    {'1': 'enabledate', '3': 99725723, '4': 1, '5': 9, '10': 'enabledate'},
    {'1': 'qrcodepng', '3': 218976787, '4': 1, '5': 12, '10': 'qrcodepng'},
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.iam.Tag',
      '10': 'tags'
    },
    {
      '1': 'user',
      '3': 10894867,
      '4': 1,
      '5': 11,
      '6': '.iam.User',
      '10': 'user'
    },
  ],
};

/// Descriptor for `VirtualMFADevice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List virtualMFADeviceDescriptor = $convert.base64Decode(
    'ChBWaXJ0dWFsTUZBRGV2aWNlEi0KEGJhc2UzMnN0cmluZ3NlZWQYwtzUeiABKAxSEGJhc2UzMn'
    'N0cmluZ3NlZWQSIQoKZW5hYmxlZGF0ZRib48YvIAEoCVIKZW5hYmxlZGF0ZRIfCglxcmNvZGVw'
    'bmcYk6S1aCABKAxSCXFyY29kZXBuZxImCgxzZXJpYWxudW1iZXIY5bq5xwEgASgJUgxzZXJpYW'
    'xudW1iZXISIAoEdGFncxjBwfa1ASADKAsyCC5pYW0uVGFnUgR0YWdzEiAKBHVzZXIYk/yYBSAB'
    'KAsyCS5pYW0uVXNlclIEdXNlcg==');
